package toml

import (
	"fmt"
	"strconv"
	"unicode"
)

// Define tokens
type tokenType int

const (
	eof = -(iota + 1)
)

const (
	tokenError tokenType = iota
	tokenEOF
	tokenComment
	tokenKey
	tokenString
	tokenInteger
	tokenTrue
	tokenFalse
	tokenFloat
	tokenInf
	tokenNan
	tokenEqual
	tokenLeftBracket
	tokenRightBracket
	tokenLeftCurlyBrace
	tokenRightCurlyBrace
	tokenLeftParen
	tokenRightParen
	tokenDoubleLeftBracket
	tokenDoubleRightBracket
	tokenDate
	tokenKeyGroup
	tokenKeyGroupArray
	tokenComma
	tokenColon
	tokenDollar
	tokenStar
	tokenQuestion
	tokenDot
	tokenDotDot
	tokenEOL
)

var tokenTypeNames = []string***REMOVED***
	"Error",
	"EOF",
	"Comment",
	"Key",
	"String",
	"Integer",
	"True",
	"False",
	"Float",
	"Inf",
	"NaN",
	"=",
	"[",
	"]",
	"***REMOVED***",
	"***REMOVED***",
	"(",
	")",
	"]]",
	"[[",
	"Date",
	"KeyGroup",
	"KeyGroupArray",
	",",
	":",
	"$",
	"*",
	"?",
	".",
	"..",
	"EOL",
***REMOVED***

type token struct ***REMOVED***
	Position
	typ tokenType
	val string
***REMOVED***

func (tt tokenType) String() string ***REMOVED***
	idx := int(tt)
	if idx < len(tokenTypeNames) ***REMOVED***
		return tokenTypeNames[idx]
	***REMOVED***
	return "Unknown"
***REMOVED***

func (t token) Int() int ***REMOVED***
	if result, err := strconv.Atoi(t.val); err != nil ***REMOVED***
		panic(err)
	***REMOVED*** else ***REMOVED***
		return result
	***REMOVED***
***REMOVED***

func (t token) String() string ***REMOVED***
	switch t.typ ***REMOVED***
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	***REMOVED***

	return fmt.Sprintf("%q", t.val)
***REMOVED***

func isSpace(r rune) bool ***REMOVED***
	return r == ' ' || r == '\t'
***REMOVED***

func isAlphanumeric(r rune) bool ***REMOVED***
	return unicode.IsLetter(r) || r == '_'
***REMOVED***

func isKeyChar(r rune) bool ***REMOVED***
	// Keys start with the first character that isn't whitespace or [ and end
	// with the last non-whitespace character before the equals sign. Keys
	// cannot contain a # character."
	return !(r == '\r' || r == '\n' || r == eof || r == '=')
***REMOVED***

func isKeyStartChar(r rune) bool ***REMOVED***
	return !(isSpace(r) || r == '\r' || r == '\n' || r == eof || r == '[')
***REMOVED***

func isDigit(r rune) bool ***REMOVED***
	return unicode.IsNumber(r)
***REMOVED***

func isHexDigit(r rune) bool ***REMOVED***
	return isDigit(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
***REMOVED***
