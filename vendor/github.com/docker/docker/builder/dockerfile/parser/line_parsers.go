package parser

// line parsers are dispatch calls that parse a single unit of text into a
// Node object which contains the whole statement. Dockerfiles have varied
// (but not usually unique, see ONBUILD for a unique example) parsing rules
// per-command, and these unify the processing in a way that makes it
// manageable.

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/docker/docker/builder/dockerfile/command"
)

var (
	errDockerfileNotStringArray = errors.New("when using JSON array syntax, arrays must be comprised of strings only")
)

const (
	commandLabel = "LABEL"
)

// ignore the current argument. This will still leave a command parsed, but
// will not incorporate the arguments into the ast.
func parseIgnore(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	return &Node***REMOVED******REMOVED***, nil, nil
***REMOVED***

// used for onbuild. Could potentially be used for anything that represents a
// statement with sub-statements.
//
// ONBUILD RUN foo bar -> (onbuild (run foo bar))
//
func parseSubCommand(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	if rest == "" ***REMOVED***
		return nil, nil, nil
	***REMOVED***

	child, err := newNodeFromLine(rest, d)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return &Node***REMOVED***Children: []*Node***REMOVED***child***REMOVED******REMOVED***, nil, nil
***REMOVED***

// helper to parse words (i.e space delimited or quoted strings) in a statement.
// The quotes are preserved as part of this function and they are stripped later
// as part of processWords().
func parseWords(rest string, d *Directive) []string ***REMOVED***
	const (
		inSpaces = iota // looking for start of a word
		inWord
		inQuote
	)

	words := []string***REMOVED******REMOVED***
	phase := inSpaces
	word := ""
	quote := '\000'
	blankOK := false
	var ch rune
	var chWidth int

	for pos := 0; pos <= len(rest); pos += chWidth ***REMOVED***
		if pos != len(rest) ***REMOVED***
			ch, chWidth = utf8.DecodeRuneInString(rest[pos:])
		***REMOVED***

		if phase == inSpaces ***REMOVED*** // Looking for start of word
			if pos == len(rest) ***REMOVED*** // end of input
				break
			***REMOVED***
			if unicode.IsSpace(ch) ***REMOVED*** // skip spaces
				continue
			***REMOVED***
			phase = inWord // found it, fall through
		***REMOVED***
		if (phase == inWord || phase == inQuote) && (pos == len(rest)) ***REMOVED***
			if blankOK || len(word) > 0 ***REMOVED***
				words = append(words, word)
			***REMOVED***
			break
		***REMOVED***
		if phase == inWord ***REMOVED***
			if unicode.IsSpace(ch) ***REMOVED***
				phase = inSpaces
				if blankOK || len(word) > 0 ***REMOVED***
					words = append(words, word)
				***REMOVED***
				word = ""
				blankOK = false
				continue
			***REMOVED***
			if ch == '\'' || ch == '"' ***REMOVED***
				quote = ch
				blankOK = true
				phase = inQuote
			***REMOVED***
			if ch == d.escapeToken ***REMOVED***
				if pos+chWidth == len(rest) ***REMOVED***
					continue // just skip an escape token at end of line
				***REMOVED***
				// If we're not quoted and we see an escape token, then always just
				// add the escape token plus the char to the word, even if the char
				// is a quote.
				word += string(ch)
				pos += chWidth
				ch, chWidth = utf8.DecodeRuneInString(rest[pos:])
			***REMOVED***
			word += string(ch)
			continue
		***REMOVED***
		if phase == inQuote ***REMOVED***
			if ch == quote ***REMOVED***
				phase = inWord
			***REMOVED***
			// The escape token is special except for ' quotes - can't escape anything for '
			if ch == d.escapeToken && quote != '\'' ***REMOVED***
				if pos+chWidth == len(rest) ***REMOVED***
					phase = inWord
					continue // just skip the escape token at end
				***REMOVED***
				pos += chWidth
				word += string(ch)
				ch, chWidth = utf8.DecodeRuneInString(rest[pos:])
			***REMOVED***
			word += string(ch)
		***REMOVED***
	***REMOVED***

	return words
***REMOVED***

// parse environment like statements. Note that this does *not* handle
// variable interpolation, which will be handled in the evaluator.
func parseNameVal(rest string, key string, d *Directive) (*Node, error) ***REMOVED***
	// This is kind of tricky because we need to support the old
	// variant:   KEY name value
	// as well as the new one:    KEY name=value ...
	// The trigger to know which one is being used will be whether we hit
	// a space or = first.  space ==> old, "=" ==> new

	words := parseWords(rest, d)
	if len(words) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	// Old format (KEY name value)
	if !strings.Contains(words[0], "=") ***REMOVED***
		parts := tokenWhitespace.Split(rest, 2)
		if len(parts) < 2 ***REMOVED***
			return nil, fmt.Errorf(key + " must have two arguments")
		***REMOVED***
		return newKeyValueNode(parts[0], parts[1]), nil
	***REMOVED***

	var rootNode *Node
	var prevNode *Node
	for _, word := range words ***REMOVED***
		if !strings.Contains(word, "=") ***REMOVED***
			return nil, fmt.Errorf("Syntax error - can't find = in %q. Must be of the form: name=value", word)
		***REMOVED***

		parts := strings.SplitN(word, "=", 2)
		node := newKeyValueNode(parts[0], parts[1])
		rootNode, prevNode = appendKeyValueNode(node, rootNode, prevNode)
	***REMOVED***

	return rootNode, nil
***REMOVED***

func newKeyValueNode(key, value string) *Node ***REMOVED***
	return &Node***REMOVED***
		Value: key,
		Next:  &Node***REMOVED***Value: value***REMOVED***,
	***REMOVED***
***REMOVED***

func appendKeyValueNode(node, rootNode, prevNode *Node) (*Node, *Node) ***REMOVED***
	if rootNode == nil ***REMOVED***
		rootNode = node
	***REMOVED***
	if prevNode != nil ***REMOVED***
		prevNode.Next = node
	***REMOVED***

	prevNode = node.Next
	return rootNode, prevNode
***REMOVED***

func parseEnv(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	node, err := parseNameVal(rest, "ENV", d)
	return node, nil, err
***REMOVED***

func parseLabel(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	node, err := parseNameVal(rest, commandLabel, d)
	return node, nil, err
***REMOVED***

// NodeFromLabels returns a Node for the injected labels
func NodeFromLabels(labels map[string]string) *Node ***REMOVED***
	keys := []string***REMOVED******REMOVED***
	for key := range labels ***REMOVED***
		keys = append(keys, key)
	***REMOVED***
	// Sort the label to have a repeatable order
	sort.Strings(keys)

	labelPairs := []string***REMOVED******REMOVED***
	var rootNode *Node
	var prevNode *Node
	for _, key := range keys ***REMOVED***
		value := labels[key]
		labelPairs = append(labelPairs, fmt.Sprintf("%q='%s'", key, value))
		// Value must be single quoted to prevent env variable expansion
		// See https://github.com/docker/docker/issues/26027
		node := newKeyValueNode(key, "'"+value+"'")
		rootNode, prevNode = appendKeyValueNode(node, rootNode, prevNode)
	***REMOVED***

	return &Node***REMOVED***
		Value:    command.Label,
		Original: commandLabel + " " + strings.Join(labelPairs, " "),
		Next:     rootNode,
	***REMOVED***
***REMOVED***

// parses a statement containing one or more keyword definition(s) and/or
// value assignments, like `name1 name2= name3="" name4=value`.
// Note that this is a stricter format than the old format of assignment,
// allowed by parseNameVal(), in a way that this only allows assignment of the
// form `keyword=[<value>]` like  `name2=`, `name3=""`, and `name4=value` above.
// In addition, a keyword definition alone is of the form `keyword` like `name1`
// above. And the assignments `name2=` and `name3=""` are equivalent and
// assign an empty value to the respective keywords.
func parseNameOrNameVal(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	words := parseWords(rest, d)
	if len(words) == 0 ***REMOVED***
		return nil, nil, nil
	***REMOVED***

	var (
		rootnode *Node
		prevNode *Node
	)
	for i, word := range words ***REMOVED***
		node := &Node***REMOVED******REMOVED***
		node.Value = word
		if i == 0 ***REMOVED***
			rootnode = node
		***REMOVED*** else ***REMOVED***
			prevNode.Next = node
		***REMOVED***
		prevNode = node
	***REMOVED***

	return rootnode, nil, nil
***REMOVED***

// parses a whitespace-delimited set of arguments. The result is effectively a
// linked list of string arguments.
func parseStringsWhitespaceDelimited(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	if rest == "" ***REMOVED***
		return nil, nil, nil
	***REMOVED***

	node := &Node***REMOVED******REMOVED***
	rootnode := node
	prevnode := node
	for _, str := range tokenWhitespace.Split(rest, -1) ***REMOVED*** // use regexp
		prevnode = node
		node.Value = str
		node.Next = &Node***REMOVED******REMOVED***
		node = node.Next
	***REMOVED***

	// XXX to get around regexp.Split *always* providing an empty string at the
	// end due to how our loop is constructed, nil out the last node in the
	// chain.
	prevnode.Next = nil

	return rootnode, nil, nil
***REMOVED***

// parseString just wraps the string in quotes and returns a working node.
func parseString(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	if rest == "" ***REMOVED***
		return nil, nil, nil
	***REMOVED***
	n := &Node***REMOVED******REMOVED***
	n.Value = rest
	return n, nil, nil
***REMOVED***

// parseJSON converts JSON arrays to an AST.
func parseJSON(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	rest = strings.TrimLeftFunc(rest, unicode.IsSpace)
	if !strings.HasPrefix(rest, "[") ***REMOVED***
		return nil, nil, fmt.Errorf(`Error parsing "%s" as a JSON array`, rest)
	***REMOVED***

	var myJSON []interface***REMOVED******REMOVED***
	if err := json.NewDecoder(strings.NewReader(rest)).Decode(&myJSON); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	var top, prev *Node
	for _, str := range myJSON ***REMOVED***
		s, ok := str.(string)
		if !ok ***REMOVED***
			return nil, nil, errDockerfileNotStringArray
		***REMOVED***

		node := &Node***REMOVED***Value: s***REMOVED***
		if prev == nil ***REMOVED***
			top = node
		***REMOVED*** else ***REMOVED***
			prev.Next = node
		***REMOVED***
		prev = node
	***REMOVED***

	return top, map[string]bool***REMOVED***"json": true***REMOVED***, nil
***REMOVED***

// parseMaybeJSON determines if the argument appears to be a JSON array. If
// so, passes to parseJSON; if not, quotes the result and returns a single
// node.
func parseMaybeJSON(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	if rest == "" ***REMOVED***
		return nil, nil, nil
	***REMOVED***

	node, attrs, err := parseJSON(rest, d)

	if err == nil ***REMOVED***
		return node, attrs, nil
	***REMOVED***
	if err == errDockerfileNotStringArray ***REMOVED***
		return nil, nil, err
	***REMOVED***

	node = &Node***REMOVED******REMOVED***
	node.Value = rest
	return node, nil, nil
***REMOVED***

// parseMaybeJSONToList determines if the argument appears to be a JSON array. If
// so, passes to parseJSON; if not, attempts to parse it as a whitespace
// delimited string.
func parseMaybeJSONToList(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	node, attrs, err := parseJSON(rest, d)

	if err == nil ***REMOVED***
		return node, attrs, nil
	***REMOVED***
	if err == errDockerfileNotStringArray ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return parseStringsWhitespaceDelimited(rest, d)
***REMOVED***

// The HEALTHCHECK command is like parseMaybeJSON, but has an extra type argument.
func parseHealthConfig(rest string, d *Directive) (*Node, map[string]bool, error) ***REMOVED***
	// Find end of first argument
	var sep int
	for ; sep < len(rest); sep++ ***REMOVED***
		if unicode.IsSpace(rune(rest[sep])) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	next := sep
	for ; next < len(rest); next++ ***REMOVED***
		if !unicode.IsSpace(rune(rest[next])) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if sep == 0 ***REMOVED***
		return nil, nil, nil
	***REMOVED***

	typ := rest[:sep]
	cmd, attrs, err := parseMaybeJSON(rest[next:], d)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return &Node***REMOVED***Value: typ, Next: cmd***REMOVED***, attrs, err
***REMOVED***
