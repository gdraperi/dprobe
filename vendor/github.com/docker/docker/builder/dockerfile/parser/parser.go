// Package parser implements a parser and parse tree dumper for Dockerfiles.
package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/docker/docker/builder/dockerfile/command"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
)

// Node is a structure used to represent a parse tree.
//
// In the node there are three fields, Value, Next, and Children. Value is the
// current token's string value. Next is always the next non-child token, and
// children contains all the children. Here's an example:
//
// (value next (child child-next child-next-next) next-next)
//
// This data structure is frankly pretty lousy for handling complex languages,
// but lucky for us the Dockerfile isn't very complicated. This structure
// works a little more effectively than a "proper" parse tree for our needs.
//
type Node struct ***REMOVED***
	Value      string          // actual content
	Next       *Node           // the next item in the current sexp
	Children   []*Node         // the children of this sexp
	Attributes map[string]bool // special attributes for this node
	Original   string          // original line used before parsing
	Flags      []string        // only top Node should have this set
	StartLine  int             // the line in the original dockerfile where the node begins
	endLine    int             // the line in the original dockerfile where the node ends
***REMOVED***

// Dump dumps the AST defined by `node` as a list of sexps.
// Returns a string suitable for printing.
func (node *Node) Dump() string ***REMOVED***
	str := ""
	str += node.Value

	if len(node.Flags) > 0 ***REMOVED***
		str += fmt.Sprintf(" %q", node.Flags)
	***REMOVED***

	for _, n := range node.Children ***REMOVED***
		str += "(" + n.Dump() + ")\n"
	***REMOVED***

	for n := node.Next; n != nil; n = n.Next ***REMOVED***
		if len(n.Children) > 0 ***REMOVED***
			str += " " + n.Dump()
		***REMOVED*** else ***REMOVED***
			str += " " + strconv.Quote(n.Value)
		***REMOVED***
	***REMOVED***

	return strings.TrimSpace(str)
***REMOVED***

func (node *Node) lines(start, end int) ***REMOVED***
	node.StartLine = start
	node.endLine = end
***REMOVED***

// AddChild adds a new child node, and updates line information
func (node *Node) AddChild(child *Node, startLine, endLine int) ***REMOVED***
	child.lines(startLine, endLine)
	if node.StartLine < 0 ***REMOVED***
		node.StartLine = startLine
	***REMOVED***
	node.endLine = endLine
	node.Children = append(node.Children, child)
***REMOVED***

var (
	dispatch             map[string]func(string, *Directive) (*Node, map[string]bool, error)
	tokenWhitespace      = regexp.MustCompile(`[\t\v\f\r ]+`)
	tokenEscapeCommand   = regexp.MustCompile(`^#[ \t]*escape[ \t]*=[ \t]*(?P<escapechar>.).*$`)
	tokenPlatformCommand = regexp.MustCompile(`^#[ \t]*platform[ \t]*=[ \t]*(?P<platform>.*)$`)
	tokenComment         = regexp.MustCompile(`^#.*$`)
)

// DefaultEscapeToken is the default escape token
const DefaultEscapeToken = '\\'

// Directive is the structure used during a build run to hold the state of
// parsing directives.
type Directive struct ***REMOVED***
	escapeToken           rune           // Current escape token
	platformToken         string         // Current platform token
	lineContinuationRegex *regexp.Regexp // Current line continuation regex
	processingComplete    bool           // Whether we are done looking for directives
	escapeSeen            bool           // Whether the escape directive has been seen
	platformSeen          bool           // Whether the platform directive has been seen
***REMOVED***

// setEscapeToken sets the default token for escaping characters in a Dockerfile.
func (d *Directive) setEscapeToken(s string) error ***REMOVED***
	if s != "`" && s != "\\" ***REMOVED***
		return fmt.Errorf("invalid ESCAPE '%s'. Must be ` or \\", s)
	***REMOVED***
	d.escapeToken = rune(s[0])
	d.lineContinuationRegex = regexp.MustCompile(`\` + s + `[ \t]*$`)
	return nil
***REMOVED***

// setPlatformToken sets the default platform for pulling images in a Dockerfile.
func (d *Directive) setPlatformToken(s string) error ***REMOVED***
	s = strings.ToLower(s)
	valid := []string***REMOVED***runtime.GOOS***REMOVED***
	if system.LCOWSupported() ***REMOVED***
		valid = append(valid, "linux")
	***REMOVED***
	for _, item := range valid ***REMOVED***
		if s == item ***REMOVED***
			d.platformToken = s
			return nil
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("invalid PLATFORM '%s'. Must be one of %v", s, valid)
***REMOVED***

// possibleParserDirective looks for one or more parser directives '# escapeToken=<char>' and
// '# platform=<string>'. Parser directives must precede any builder instruction
// or other comments, and cannot be repeated.
func (d *Directive) possibleParserDirective(line string) error ***REMOVED***
	if d.processingComplete ***REMOVED***
		return nil
	***REMOVED***

	tecMatch := tokenEscapeCommand.FindStringSubmatch(strings.ToLower(line))
	if len(tecMatch) != 0 ***REMOVED***
		for i, n := range tokenEscapeCommand.SubexpNames() ***REMOVED***
			if n == "escapechar" ***REMOVED***
				if d.escapeSeen ***REMOVED***
					return errors.New("only one escape parser directive can be used")
				***REMOVED***
				d.escapeSeen = true
				return d.setEscapeToken(tecMatch[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Only recognise a platform token if LCOW is supported
	if system.LCOWSupported() ***REMOVED***
		tpcMatch := tokenPlatformCommand.FindStringSubmatch(strings.ToLower(line))
		if len(tpcMatch) != 0 ***REMOVED***
			for i, n := range tokenPlatformCommand.SubexpNames() ***REMOVED***
				if n == "platform" ***REMOVED***
					if d.platformSeen ***REMOVED***
						return errors.New("only one platform parser directive can be used")
					***REMOVED***
					d.platformSeen = true
					return d.setPlatformToken(tpcMatch[i])
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	d.processingComplete = true
	return nil
***REMOVED***

// NewDefaultDirective returns a new Directive with the default escapeToken token
func NewDefaultDirective() *Directive ***REMOVED***
	directive := Directive***REMOVED******REMOVED***
	directive.setEscapeToken(string(DefaultEscapeToken))
	return &directive
***REMOVED***

func init() ***REMOVED***
	// Dispatch Table. see line_parsers.go for the parse functions.
	// The command is parsed and mapped to the line parser. The line parser
	// receives the arguments but not the command, and returns an AST after
	// reformulating the arguments according to the rules in the parser
	// functions. Errors are propagated up by Parse() and the resulting AST can
	// be incorporated directly into the existing AST as a next.
	dispatch = map[string]func(string, *Directive) (*Node, map[string]bool, error)***REMOVED***
		command.Add:         parseMaybeJSONToList,
		command.Arg:         parseNameOrNameVal,
		command.Cmd:         parseMaybeJSON,
		command.Copy:        parseMaybeJSONToList,
		command.Entrypoint:  parseMaybeJSON,
		command.Env:         parseEnv,
		command.Expose:      parseStringsWhitespaceDelimited,
		command.From:        parseStringsWhitespaceDelimited,
		command.Healthcheck: parseHealthConfig,
		command.Label:       parseLabel,
		command.Maintainer:  parseString,
		command.Onbuild:     parseSubCommand,
		command.Run:         parseMaybeJSON,
		command.Shell:       parseMaybeJSON,
		command.StopSignal:  parseString,
		command.User:        parseString,
		command.Volume:      parseMaybeJSONToList,
		command.Workdir:     parseString,
	***REMOVED***
***REMOVED***

// newNodeFromLine splits the line into parts, and dispatches to a function
// based on the command and command arguments. A Node is created from the
// result of the dispatch.
func newNodeFromLine(line string, directive *Directive) (*Node, error) ***REMOVED***
	cmd, flags, args, err := splitCommand(line)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	fn := dispatch[cmd]
	// Ignore invalid Dockerfile instructions
	if fn == nil ***REMOVED***
		fn = parseIgnore
	***REMOVED***
	next, attrs, err := fn(args, directive)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Node***REMOVED***
		Value:      cmd,
		Original:   line,
		Flags:      flags,
		Next:       next,
		Attributes: attrs,
	***REMOVED***, nil
***REMOVED***

// Result is the result of parsing a Dockerfile
type Result struct ***REMOVED***
	AST         *Node
	EscapeToken rune
	// TODO @jhowardmsft - see https://github.com/moby/moby/issues/34617
	// This next field will be removed in a future update for LCOW support.
	OS       string
	Warnings []string
***REMOVED***

// PrintWarnings to the writer
func (r *Result) PrintWarnings(out io.Writer) ***REMOVED***
	if len(r.Warnings) == 0 ***REMOVED***
		return
	***REMOVED***
	fmt.Fprintf(out, strings.Join(r.Warnings, "\n")+"\n")
***REMOVED***

// Parse reads lines from a Reader, parses the lines into an AST and returns
// the AST and escape token
func Parse(rwc io.Reader) (*Result, error) ***REMOVED***
	d := NewDefaultDirective()
	currentLine := 0
	root := &Node***REMOVED***StartLine: -1***REMOVED***
	scanner := bufio.NewScanner(rwc)
	warnings := []string***REMOVED******REMOVED***

	var err error
	for scanner.Scan() ***REMOVED***
		bytesRead := scanner.Bytes()
		if currentLine == 0 ***REMOVED***
			// First line, strip the byte-order-marker if present
			bytesRead = bytes.TrimPrefix(bytesRead, utf8bom)
		***REMOVED***
		bytesRead, err = processLine(d, bytesRead, true)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		currentLine++

		startLine := currentLine
		line, isEndOfLine := trimContinuationCharacter(string(bytesRead), d)
		if isEndOfLine && line == "" ***REMOVED***
			continue
		***REMOVED***

		var hasEmptyContinuationLine bool
		for !isEndOfLine && scanner.Scan() ***REMOVED***
			bytesRead, err := processLine(d, scanner.Bytes(), false)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			currentLine++

			if isComment(scanner.Bytes()) ***REMOVED***
				// original line was a comment (processLine strips comments)
				continue
			***REMOVED***
			if isEmptyContinuationLine(bytesRead) ***REMOVED***
				hasEmptyContinuationLine = true
				continue
			***REMOVED***

			continuationLine := string(bytesRead)
			continuationLine, isEndOfLine = trimContinuationCharacter(continuationLine, d)
			line += continuationLine
		***REMOVED***

		if hasEmptyContinuationLine ***REMOVED***
			warning := "[WARNING]: Empty continuation line found in:\n    " + line
			warnings = append(warnings, warning)
		***REMOVED***

		child, err := newNodeFromLine(line, d)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		root.AddChild(child, startLine, currentLine)
	***REMOVED***

	if len(warnings) > 0 ***REMOVED***
		warnings = append(warnings, "[WARNING]: Empty continuation lines will become errors in a future release.")
	***REMOVED***
	return &Result***REMOVED***
		AST:         root,
		Warnings:    warnings,
		EscapeToken: d.escapeToken,
		OS:          d.platformToken,
	***REMOVED***, handleScannerError(scanner.Err())
***REMOVED***

func trimComments(src []byte) []byte ***REMOVED***
	return tokenComment.ReplaceAll(src, []byte***REMOVED******REMOVED***)
***REMOVED***

func trimWhitespace(src []byte) []byte ***REMOVED***
	return bytes.TrimLeftFunc(src, unicode.IsSpace)
***REMOVED***

func isComment(line []byte) bool ***REMOVED***
	return tokenComment.Match(trimWhitespace(line))
***REMOVED***

func isEmptyContinuationLine(line []byte) bool ***REMOVED***
	return len(trimWhitespace(line)) == 0
***REMOVED***

var utf8bom = []byte***REMOVED***0xEF, 0xBB, 0xBF***REMOVED***

func trimContinuationCharacter(line string, d *Directive) (string, bool) ***REMOVED***
	if d.lineContinuationRegex.MatchString(line) ***REMOVED***
		line = d.lineContinuationRegex.ReplaceAllString(line, "")
		return line, false
	***REMOVED***
	return line, true
***REMOVED***

// TODO: remove stripLeftWhitespace after deprecation period. It seems silly
// to preserve whitespace on continuation lines. Why is that done?
func processLine(d *Directive, token []byte, stripLeftWhitespace bool) ([]byte, error) ***REMOVED***
	if stripLeftWhitespace ***REMOVED***
		token = trimWhitespace(token)
	***REMOVED***
	return trimComments(token), d.possibleParserDirective(string(token))
***REMOVED***

func handleScannerError(err error) error ***REMOVED***
	switch err ***REMOVED***
	case bufio.ErrTooLong:
		return errors.Errorf("dockerfile line greater than max allowed size of %d", bufio.MaxScanTokenSize-1)
	default:
		return err
	***REMOVED***
***REMOVED***
