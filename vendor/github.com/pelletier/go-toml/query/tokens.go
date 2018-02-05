package query

import (
	"fmt"
	"github.com/pelletier/go-toml"
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
	tokenKey
	tokenString
	tokenInteger
	tokenFloat
	tokenLeftBracket
	tokenRightBracket
	tokenLeftParen
	tokenRightParen
	tokenComma
	tokenColon
	tokenDollar
	tokenStar
	tokenQuestion
	tokenDot
	tokenDotDot
)

var tokenTypeNames = []string***REMOVED***
	"Error",
	"EOF",
	"Key",
	"String",
	"Integer",
	"Float",
	"[",
	"]",
	"(",
	")",
	",",
	":",
	"$",
	"*",
	"?",
	".",
	"..",
***REMOVED***

type token struct ***REMOVED***
	toml.Position
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

func isDigit(r rune) bool ***REMOVED***
	return unicode.IsNumber(r)
***REMOVED***

func isHexDigit(r rune) bool ***REMOVED***
	return isDigit(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
***REMOVED***
