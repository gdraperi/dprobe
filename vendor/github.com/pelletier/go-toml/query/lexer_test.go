package query

import (
	"github.com/pelletier/go-toml"
	"testing"
)

func testQLFlow(t *testing.T, input string, expectedFlow []token) ***REMOVED***
	ch := lexQuery(input)
	for idx, expected := range expectedFlow ***REMOVED***
		token := <-ch
		if token != expected ***REMOVED***
			t.Log("While testing #", idx, ":", input)
			t.Log("compared (got)", token, "to (expected)", expected)
			t.Log("\tvalue:", token.val, "<->", expected.val)
			t.Log("\tvalue as bytes:", []byte(token.val), "<->", []byte(expected.val))
			t.Log("\ttype:", token.typ.String(), "<->", expected.typ.String())
			t.Log("\tline:", token.Line, "<->", expected.Line)
			t.Log("\tcolumn:", token.Col, "<->", expected.Col)
			t.Log("compared", token, "to", expected)
			t.FailNow()
		***REMOVED***
	***REMOVED***

	tok, ok := <-ch
	if ok ***REMOVED***
		t.Log("channel is not closed!")
		t.Log(len(ch)+1, "tokens remaining:")

		t.Log("token ->", tok)
		for token := range ch ***REMOVED***
			t.Log("token ->", token)
		***REMOVED***
		t.FailNow()
	***REMOVED***
***REMOVED***

func TestLexSpecialChars(t *testing.T) ***REMOVED***
	testQLFlow(t, " .$[]..()?*", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenDot, "."***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 3***REMOVED***, tokenDollar, "$"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 4***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 5***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 6***REMOVED***, tokenDotDot, ".."***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 8***REMOVED***, tokenLeftParen, "("***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 9***REMOVED***, tokenRightParen, ")"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 10***REMOVED***, tokenQuestion, "?"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 11***REMOVED***, tokenStar, "*"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 12***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexString(t *testing.T) ***REMOVED***
	testQLFlow(t, "'foo\n'", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenString, "foo\n"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***2, 2***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexDoubleString(t *testing.T) ***REMOVED***
	testQLFlow(t, `"bar"`, []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenString, "bar"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 6***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexStringEscapes(t *testing.T) ***REMOVED***
	testQLFlow(t, `"foo \" \' \b \f \/ \t \r \\ \u03A9 \U00012345 \n bar"`, []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenString, "foo \" ' \b \f / \t \r \\ \u03A9 \U00012345 \n bar"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 55***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexStringUnfinishedUnicode4(t *testing.T) ***REMOVED***
	testQLFlow(t, `"\u000"`, []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenError, "unfinished unicode escape"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexStringUnfinishedUnicode8(t *testing.T) ***REMOVED***
	testQLFlow(t, `"\U0000"`, []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenError, "unfinished unicode escape"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexStringInvalidEscape(t *testing.T) ***REMOVED***
	testQLFlow(t, `"\x"`, []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenError, "invalid escape sequence: \\x"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexStringUnfinished(t *testing.T) ***REMOVED***
	testQLFlow(t, `"bar`, []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenError, "unclosed string"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexKey(t *testing.T) ***REMOVED***
	testQLFlow(t, "foo", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 4***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexRecurse(t *testing.T) ***REMOVED***
	testQLFlow(t, "$..*", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenDollar, "$"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenDotDot, ".."***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 4***REMOVED***, tokenStar, "*"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 5***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexBracketKey(t *testing.T) ***REMOVED***
	testQLFlow(t, "$[foo]", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenDollar, "$"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 2***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 3***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 6***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 7***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexSpace(t *testing.T) ***REMOVED***
	testQLFlow(t, "foo bar baz", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 5***REMOVED***, tokenKey, "bar"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 9***REMOVED***, tokenKey, "baz"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 12***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexInteger(t *testing.T) ***REMOVED***
	testQLFlow(t, "100 +200 -300", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenInteger, "100"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 5***REMOVED***, tokenInteger, "+200"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 10***REMOVED***, tokenInteger, "-300"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 14***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexFloat(t *testing.T) ***REMOVED***
	testQLFlow(t, "100.0 +200.0 -300.0", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenFloat, "100.0"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 7***REMOVED***, tokenFloat, "+200.0"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 14***REMOVED***, tokenFloat, "-300.0"***REMOVED***,
		***REMOVED***toml.Position***REMOVED***1, 20***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexFloatWithMultipleDots(t *testing.T) ***REMOVED***
	testQLFlow(t, "4.2.", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenError, "cannot have two dots in one float"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexFloatLeadingDot(t *testing.T) ***REMOVED***
	testQLFlow(t, "+.1", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenError, "cannot start float with a dot"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexFloatWithTrailingDot(t *testing.T) ***REMOVED***
	testQLFlow(t, "42.", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenError, "float cannot end with a dot"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexNumberWithoutDigit(t *testing.T) ***REMOVED***
	testQLFlow(t, "+", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenError, "no digit in that number"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexUnknown(t *testing.T) ***REMOVED***
	testQLFlow(t, "^", []token***REMOVED***
		***REMOVED***toml.Position***REMOVED***1, 1***REMOVED***, tokenError, "unexpected char: '94'"***REMOVED***,
	***REMOVED***)
***REMOVED***
