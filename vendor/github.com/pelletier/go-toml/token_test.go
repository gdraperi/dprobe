package toml

import "testing"

func TestTokenStringer(t *testing.T) ***REMOVED***
	var tests = []struct ***REMOVED***
		tt     tokenType
		expect string
	***REMOVED******REMOVED***
		***REMOVED***tokenError, "Error"***REMOVED***,
		***REMOVED***tokenEOF, "EOF"***REMOVED***,
		***REMOVED***tokenComment, "Comment"***REMOVED***,
		***REMOVED***tokenKey, "Key"***REMOVED***,
		***REMOVED***tokenString, "String"***REMOVED***,
		***REMOVED***tokenInteger, "Integer"***REMOVED***,
		***REMOVED***tokenTrue, "True"***REMOVED***,
		***REMOVED***tokenFalse, "False"***REMOVED***,
		***REMOVED***tokenFloat, "Float"***REMOVED***,
		***REMOVED***tokenEqual, "="***REMOVED***,
		***REMOVED***tokenLeftBracket, "["***REMOVED***,
		***REMOVED***tokenRightBracket, "]"***REMOVED***,
		***REMOVED***tokenLeftCurlyBrace, "***REMOVED***"***REMOVED***,
		***REMOVED***tokenRightCurlyBrace, "***REMOVED***"***REMOVED***,
		***REMOVED***tokenLeftParen, "("***REMOVED***,
		***REMOVED***tokenRightParen, ")"***REMOVED***,
		***REMOVED***tokenDoubleLeftBracket, "]]"***REMOVED***,
		***REMOVED***tokenDoubleRightBracket, "[["***REMOVED***,
		***REMOVED***tokenDate, "Date"***REMOVED***,
		***REMOVED***tokenKeyGroup, "KeyGroup"***REMOVED***,
		***REMOVED***tokenKeyGroupArray, "KeyGroupArray"***REMOVED***,
		***REMOVED***tokenComma, ","***REMOVED***,
		***REMOVED***tokenColon, ":"***REMOVED***,
		***REMOVED***tokenDollar, "$"***REMOVED***,
		***REMOVED***tokenStar, "*"***REMOVED***,
		***REMOVED***tokenQuestion, "?"***REMOVED***,
		***REMOVED***tokenDot, "."***REMOVED***,
		***REMOVED***tokenDotDot, ".."***REMOVED***,
		***REMOVED***tokenEOL, "EOL"***REMOVED***,
		***REMOVED***tokenEOL + 1, "Unknown"***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		got := test.tt.String()
		if got != test.expect ***REMOVED***
			t.Errorf("[%d] invalid string of token type; got %q, expected %q", i, got, test.expect)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTokenString(t *testing.T) ***REMOVED***
	var tests = []struct ***REMOVED***
		tok    token
		expect string
	***REMOVED******REMOVED***
		***REMOVED***token***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenEOF, ""***REMOVED***, "EOF"***REMOVED***,
		***REMOVED***token***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenError, "Δt"***REMOVED***, "Δt"***REMOVED***,
		***REMOVED***token***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenString, "bar"***REMOVED***, `"bar"`***REMOVED***,
		***REMOVED***token***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenString, "123456789012345"***REMOVED***, `"123456789012345"`***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		got := test.tok.String()
		if got != test.expect ***REMOVED***
			t.Errorf("[%d] invalid of string token; got %q, expected %q", i, got, test.expect)
		***REMOVED***
	***REMOVED***
***REMOVED***
