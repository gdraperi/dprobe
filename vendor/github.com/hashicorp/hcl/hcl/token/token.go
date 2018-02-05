// Package token defines constants representing the lexical tokens for HCL
// (HashiCorp Configuration Language)
package token

import (
	"fmt"
	"strconv"
	"strings"

	hclstrconv "github.com/hashicorp/hcl/hcl/strconv"
)

// Token defines a single HCL token which can be obtained via the Scanner
type Token struct ***REMOVED***
	Type Type
	Pos  Pos
	Text string
	JSON bool
***REMOVED***

// Type is the set of lexical tokens of the HCL (HashiCorp Configuration Language)
type Type int

const (
	// Special tokens
	ILLEGAL Type = iota
	EOF
	COMMENT

	identifier_beg
	IDENT // literals
	literal_beg
	NUMBER  // 12345
	FLOAT   // 123.45
	BOOL    // true,false
	STRING  // "abc"
	HEREDOC // <<FOO\nbar\nFOO
	literal_end
	identifier_end

	operator_beg
	LBRACK // [
	LBRACE // ***REMOVED***
	COMMA  // ,
	PERIOD // .

	RBRACK // ]
	RBRACE // ***REMOVED***

	ASSIGN // =
	ADD    // +
	SUB    // -
	operator_end
)

var tokens = [...]string***REMOVED***
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",
	FLOAT:  "FLOAT",
	BOOL:   "BOOL",
	STRING: "STRING",

	LBRACK:  "LBRACK",
	LBRACE:  "LBRACE",
	COMMA:   "COMMA",
	PERIOD:  "PERIOD",
	HEREDOC: "HEREDOC",

	RBRACK: "RBRACK",
	RBRACE: "RBRACE",

	ASSIGN: "ASSIGN",
	ADD:    "ADD",
	SUB:    "SUB",
***REMOVED***

// String returns the string corresponding to the token tok.
func (t Type) String() string ***REMOVED***
	s := ""
	if 0 <= t && t < Type(len(tokens)) ***REMOVED***
		s = tokens[t]
	***REMOVED***
	if s == "" ***REMOVED***
		s = "token(" + strconv.Itoa(int(t)) + ")"
	***REMOVED***
	return s
***REMOVED***

// IsIdentifier returns true for tokens corresponding to identifiers and basic
// type literals; it returns false otherwise.
func (t Type) IsIdentifier() bool ***REMOVED*** return identifier_beg < t && t < identifier_end ***REMOVED***

// IsLiteral returns true for tokens corresponding to basic type literals; it
// returns false otherwise.
func (t Type) IsLiteral() bool ***REMOVED*** return literal_beg < t && t < literal_end ***REMOVED***

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (t Type) IsOperator() bool ***REMOVED*** return operator_beg < t && t < operator_end ***REMOVED***

// String returns the token's literal text. Note that this is only
// applicable for certain token types, such as token.IDENT,
// token.STRING, etc..
func (t Token) String() string ***REMOVED***
	return fmt.Sprintf("%s %s %s", t.Pos.String(), t.Type.String(), t.Text)
***REMOVED***

// Value returns the properly typed value for this token. The type of
// the returned interface***REMOVED******REMOVED*** is guaranteed based on the Type field.
//
// This can only be called for literal types. If it is called for any other
// type, this will panic.
func (t Token) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	switch t.Type ***REMOVED***
	case BOOL:
		if t.Text == "true" ***REMOVED***
			return true
		***REMOVED*** else if t.Text == "false" ***REMOVED***
			return false
		***REMOVED***

		panic("unknown bool value: " + t.Text)
	case FLOAT:
		v, err := strconv.ParseFloat(t.Text, 64)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		return float64(v)
	case NUMBER:
		v, err := strconv.ParseInt(t.Text, 0, 64)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		return int64(v)
	case IDENT:
		return t.Text
	case HEREDOC:
		return unindentHeredoc(t.Text)
	case STRING:
		// Determine the Unquote method to use. If it came from JSON,
		// then we need to use the built-in unquote since we have to
		// escape interpolations there.
		f := hclstrconv.Unquote
		if t.JSON ***REMOVED***
			f = strconv.Unquote
		***REMOVED***

		// This case occurs if json null is used
		if t.Text == "" ***REMOVED***
			return ""
		***REMOVED***

		v, err := f(t.Text)
		if err != nil ***REMOVED***
			panic(fmt.Sprintf("unquote %s err: %s", t.Text, err))
		***REMOVED***

		return v
	default:
		panic(fmt.Sprintf("unimplemented Value for type: %s", t.Type))
	***REMOVED***
***REMOVED***

// unindentHeredoc returns the string content of a HEREDOC if it is started with <<
// and the content of a HEREDOC with the hanging indent removed if it is started with
// a <<-, and the terminating line is at least as indented as the least indented line.
func unindentHeredoc(heredoc string) string ***REMOVED***
	// We need to find the end of the marker
	idx := strings.IndexByte(heredoc, '\n')
	if idx == -1 ***REMOVED***
		panic("heredoc doesn't contain newline")
	***REMOVED***

	unindent := heredoc[2] == '-'

	// We can optimize if the heredoc isn't marked for indentation
	if !unindent ***REMOVED***
		return string(heredoc[idx+1 : len(heredoc)-idx+1])
	***REMOVED***

	// We need to unindent each line based on the indentation level of the marker
	lines := strings.Split(string(heredoc[idx+1:len(heredoc)-idx+2]), "\n")
	whitespacePrefix := lines[len(lines)-1]

	isIndented := true
	for _, v := range lines ***REMOVED***
		if strings.HasPrefix(v, whitespacePrefix) ***REMOVED***
			continue
		***REMOVED***

		isIndented = false
		break
	***REMOVED***

	// If all lines are not at least as indented as the terminating mark, return the
	// heredoc as is, but trim the leading space from the marker on the final line.
	if !isIndented ***REMOVED***
		return strings.TrimRight(string(heredoc[idx+1:len(heredoc)-idx+1]), " \t")
	***REMOVED***

	unindentedLines := make([]string, len(lines))
	for k, v := range lines ***REMOVED***
		if k == len(lines)-1 ***REMOVED***
			unindentedLines[k] = ""
			break
		***REMOVED***

		unindentedLines[k] = strings.TrimPrefix(v, whitespacePrefix)
	***REMOVED***

	return strings.Join(unindentedLines, "\n")
***REMOVED***
