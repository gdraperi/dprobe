package token

import (
	"fmt"
	"strconv"

	hcltoken "github.com/hashicorp/hcl/hcl/token"
)

// Token defines a single HCL token which can be obtained via the Scanner
type Token struct ***REMOVED***
	Type Type
	Pos  Pos
	Text string
***REMOVED***

// Type is the set of lexical tokens of the HCL (HashiCorp Configuration Language)
type Type int

const (
	// Special tokens
	ILLEGAL Type = iota
	EOF

	identifier_beg
	literal_beg
	NUMBER // 12345
	FLOAT  // 123.45
	BOOL   // true,false
	STRING // "abc"
	NULL   // null
	literal_end
	identifier_end

	operator_beg
	LBRACK // [
	LBRACE // ***REMOVED***
	COMMA  // ,
	PERIOD // .
	COLON  // :

	RBRACK // ]
	RBRACE // ***REMOVED***

	operator_end
)

var tokens = [...]string***REMOVED***
	ILLEGAL: "ILLEGAL",

	EOF: "EOF",

	NUMBER: "NUMBER",
	FLOAT:  "FLOAT",
	BOOL:   "BOOL",
	STRING: "STRING",
	NULL:   "NULL",

	LBRACK: "LBRACK",
	LBRACE: "LBRACE",
	COMMA:  "COMMA",
	PERIOD: "PERIOD",
	COLON:  "COLON",

	RBRACK: "RBRACK",
	RBRACE: "RBRACE",
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

// HCLToken converts this token to an HCL token.
//
// The token type must be a literal type or this will panic.
func (t Token) HCLToken() hcltoken.Token ***REMOVED***
	switch t.Type ***REMOVED***
	case BOOL:
		return hcltoken.Token***REMOVED***Type: hcltoken.BOOL, Text: t.Text***REMOVED***
	case FLOAT:
		return hcltoken.Token***REMOVED***Type: hcltoken.FLOAT, Text: t.Text***REMOVED***
	case NULL:
		return hcltoken.Token***REMOVED***Type: hcltoken.STRING, Text: ""***REMOVED***
	case NUMBER:
		return hcltoken.Token***REMOVED***Type: hcltoken.NUMBER, Text: t.Text***REMOVED***
	case STRING:
		return hcltoken.Token***REMOVED***Type: hcltoken.STRING, Text: t.Text, JSON: true***REMOVED***
	default:
		panic(fmt.Sprintf("unimplemented HCLToken for type: %s", t.Type))
	***REMOVED***
***REMOVED***
