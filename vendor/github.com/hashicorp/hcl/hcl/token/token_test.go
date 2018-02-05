package token

import (
	"reflect"
	"testing"
)

func TestTypeString(t *testing.T) ***REMOVED***
	var tokens = []struct ***REMOVED***
		tt  Type
		str string
	***REMOVED******REMOVED***
		***REMOVED***ILLEGAL, "ILLEGAL"***REMOVED***,
		***REMOVED***EOF, "EOF"***REMOVED***,
		***REMOVED***COMMENT, "COMMENT"***REMOVED***,
		***REMOVED***IDENT, "IDENT"***REMOVED***,
		***REMOVED***NUMBER, "NUMBER"***REMOVED***,
		***REMOVED***FLOAT, "FLOAT"***REMOVED***,
		***REMOVED***BOOL, "BOOL"***REMOVED***,
		***REMOVED***STRING, "STRING"***REMOVED***,
		***REMOVED***HEREDOC, "HEREDOC"***REMOVED***,
		***REMOVED***LBRACK, "LBRACK"***REMOVED***,
		***REMOVED***LBRACE, "LBRACE"***REMOVED***,
		***REMOVED***COMMA, "COMMA"***REMOVED***,
		***REMOVED***PERIOD, "PERIOD"***REMOVED***,
		***REMOVED***RBRACK, "RBRACK"***REMOVED***,
		***REMOVED***RBRACE, "RBRACE"***REMOVED***,
		***REMOVED***ASSIGN, "ASSIGN"***REMOVED***,
		***REMOVED***ADD, "ADD"***REMOVED***,
		***REMOVED***SUB, "SUB"***REMOVED***,
	***REMOVED***

	for _, token := range tokens ***REMOVED***
		if token.tt.String() != token.str ***REMOVED***
			t.Errorf("want: %q got:%q\n", token.str, token.tt)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestTokenValue(t *testing.T) ***REMOVED***
	var tokens = []struct ***REMOVED***
		tt Token
		v  interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED***Token***REMOVED***Type: BOOL, Text: `true`***REMOVED***, true***REMOVED***,
		***REMOVED***Token***REMOVED***Type: BOOL, Text: `false`***REMOVED***, false***REMOVED***,
		***REMOVED***Token***REMOVED***Type: FLOAT, Text: `3.14`***REMOVED***, float64(3.14)***REMOVED***,
		***REMOVED***Token***REMOVED***Type: NUMBER, Text: `42`***REMOVED***, int64(42)***REMOVED***,
		***REMOVED***Token***REMOVED***Type: IDENT, Text: `foo`***REMOVED***, "foo"***REMOVED***,
		***REMOVED***Token***REMOVED***Type: STRING, Text: `"foo"`***REMOVED***, "foo"***REMOVED***,
		***REMOVED***Token***REMOVED***Type: STRING, Text: `"foo\nbar"`***REMOVED***, "foo\nbar"***REMOVED***,
		***REMOVED***Token***REMOVED***Type: STRING, Text: `"$***REMOVED***file("foo")***REMOVED***"`***REMOVED***, `$***REMOVED***file("foo")***REMOVED***`***REMOVED***,
		***REMOVED***
			Token***REMOVED***
				Type: STRING,
				Text: `"$***REMOVED***replace("foo", ".", "\\.")***REMOVED***"`,
			***REMOVED***,
			`$***REMOVED***replace("foo", ".", "\\.")***REMOVED***`***REMOVED***,
		***REMOVED***Token***REMOVED***Type: HEREDOC, Text: "<<EOF\nfoo\nbar\nEOF"***REMOVED***, "foo\nbar"***REMOVED***,
	***REMOVED***

	for _, token := range tokens ***REMOVED***
		if val := token.tt.Value(); !reflect.DeepEqual(val, token.v) ***REMOVED***
			t.Errorf("want: %v got:%v\n", token.v, val)
		***REMOVED***
	***REMOVED***

***REMOVED***
