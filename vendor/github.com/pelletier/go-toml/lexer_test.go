package toml

import (
	"reflect"
	"testing"
)

func testFlow(t *testing.T, input string, expectedFlow []token) ***REMOVED***
	tokens := lexToml([]byte(input))
	if !reflect.DeepEqual(tokens, expectedFlow) ***REMOVED***
		t.Fatal("Different flows. Expected\n", expectedFlow, "\nGot:\n", tokens)
	***REMOVED***
***REMOVED***

func TestValidKeyGroup(t *testing.T) ***REMOVED***
	testFlow(t, "[hello world]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 2***REMOVED***, tokenKeyGroup, "hello world"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 13***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 14***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestNestedQuotedUnicodeKeyGroup(t *testing.T) ***REMOVED***
	testFlow(t, `[ j . "ʞ" . l ]`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 2***REMOVED***, tokenKeyGroup, ` j . "ʞ" . l `***REMOVED***,
		***REMOVED***Position***REMOVED***1, 15***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 16***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestUnclosedKeyGroup(t *testing.T) ***REMOVED***
	testFlow(t, "[hello world", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 2***REMOVED***, tokenError, "unclosed table key"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestComment(t *testing.T) ***REMOVED***
	testFlow(t, "# blahblah", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyGroupComment(t *testing.T) ***REMOVED***
	testFlow(t, "[hello world] # blahblah", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 2***REMOVED***, tokenKeyGroup, "hello world"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 13***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 25***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestMultipleKeyGroupsComment(t *testing.T) ***REMOVED***
	testFlow(t, "[hello world] # blahblah\n[test]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 2***REMOVED***, tokenKeyGroup, "hello world"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 13***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***2, 2***REMOVED***, tokenKeyGroup, "test"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 6***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 7***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestSimpleWindowsCRLF(t *testing.T) ***REMOVED***
	testFlow(t, "a=4\r\nb=2", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 2***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenInteger, "4"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenKey, "b"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 2***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 3***REMOVED***, tokenInteger, "2"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 4***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestBasicKey(t *testing.T) ***REMOVED***
	testFlow(t, "hello", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "hello"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestBasicKeyWithUnderscore(t *testing.T) ***REMOVED***
	testFlow(t, "hello_hello", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "hello_hello"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestBasicKeyWithDash(t *testing.T) ***REMOVED***
	testFlow(t, "hello-world", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "hello-world"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestBasicKeyWithUppercaseMix(t *testing.T) ***REMOVED***
	testFlow(t, "helloHELLOHello", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "helloHELLOHello"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 16***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestBasicKeyWithInternationalCharacters(t *testing.T) ***REMOVED***
	testFlow(t, "héllÖ", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "héllÖ"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestBasicKeyAndEqual(t *testing.T) ***REMOVED***
	testFlow(t, "hello =", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "hello"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyWithSharpAndEqual(t *testing.T) ***REMOVED***
	testFlow(t, "key#name = 5", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenError, "keys cannot contain # character"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyWithSymbolsAndEqual(t *testing.T) ***REMOVED***
	testFlow(t, "~!@$^&*()_+-`1234567890[]\\|/?><.,;:' = 5", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenError, "keys cannot contain ~ character"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualStringEscape(t *testing.T) ***REMOVED***
	testFlow(t, `foo = "hello\""`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "hello\""***REMOVED***,
		***REMOVED***Position***REMOVED***1, 16***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualStringUnfinished(t *testing.T) ***REMOVED***
	testFlow(t, `foo = "bar`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenError, "unclosed string"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualString(t *testing.T) ***REMOVED***
	testFlow(t, `foo = "bar"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "bar"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualTrue(t *testing.T) ***REMOVED***
	testFlow(t, "foo = true", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenTrue, "true"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualFalse(t *testing.T) ***REMOVED***
	testFlow(t, "foo = false", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenFalse, "false"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayNestedString(t *testing.T) ***REMOVED***
	testFlow(t, `a = [ ["hello", "world"] ]`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 9***REMOVED***, tokenString, "hello"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 15***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 18***REMOVED***, tokenString, "world"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 24***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 26***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 27***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayNestedInts(t *testing.T) ***REMOVED***
	testFlow(t, "a = [ [42, 21], [10] ]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenInteger, "42"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenInteger, "21"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 14***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 15***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 17***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 18***REMOVED***, tokenInteger, "10"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 20***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 22***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 23***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayInts(t *testing.T) ***REMOVED***
	testFlow(t, "a = [ 42, 21, 10, ]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "42"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 9***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenInteger, "21"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 13***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 15***REMOVED***, tokenInteger, "10"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 17***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 19***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 20***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestMultilineArrayComments(t *testing.T) ***REMOVED***
	testFlow(t, "a = [1, # wow\n2, # such items\n3, # so array\n]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenInteger, "1"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenInteger, "2"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 2***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***3, 1***REMOVED***, tokenInteger, "3"***REMOVED***,
		***REMOVED***Position***REMOVED***3, 2***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***4, 1***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***4, 2***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestNestedArraysComment(t *testing.T) ***REMOVED***
	toml := `
someArray = [
# does not work
["entry1"]
]`
	testFlow(t, toml, []token***REMOVED***
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenKey, "someArray"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 11***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 13***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***4, 1***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***4, 3***REMOVED***, tokenString, "entry1"***REMOVED***,
		***REMOVED***Position***REMOVED***4, 10***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***5, 1***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***5, 2***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualArrayBools(t *testing.T) ***REMOVED***
	testFlow(t, "foo = [true, false, true]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenTrue, "true"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 14***REMOVED***, tokenFalse, "false"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 19***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 21***REMOVED***, tokenTrue, "true"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 25***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 26***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualArrayBoolsWithComments(t *testing.T) ***REMOVED***
	testFlow(t, "foo = [true, false, true] # YEAH", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenTrue, "true"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 14***REMOVED***, tokenFalse, "false"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 19***REMOVED***, tokenComma, ","***REMOVED***,
		***REMOVED***Position***REMOVED***1, 21***REMOVED***, tokenTrue, "true"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 25***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 33***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestDateRegexp(t *testing.T) ***REMOVED***
	if dateRegexp.FindString("1979-05-27T07:32:00Z") == "" ***REMOVED***
		t.Error("basic lexing")
	***REMOVED***
	if dateRegexp.FindString("1979-05-27T00:32:00-07:00") == "" ***REMOVED***
		t.Error("offset lexing")
	***REMOVED***
	if dateRegexp.FindString("1979-05-27T00:32:00.999999-07:00") == "" ***REMOVED***
		t.Error("nano precision lexing")
	***REMOVED***
***REMOVED***

func TestKeyEqualDate(t *testing.T) ***REMOVED***
	testFlow(t, "foo = 1979-05-27T07:32:00Z", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenDate, "1979-05-27T07:32:00Z"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 27***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, "foo = 1979-05-27T00:32:00-07:00", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenDate, "1979-05-27T00:32:00-07:00"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 32***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, "foo = 1979-05-27T00:32:00.999999-07:00", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenDate, "1979-05-27T00:32:00.999999-07:00"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 39***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatEndingWithDot(t *testing.T) ***REMOVED***
	testFlow(t, "foo = 42.", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenError, "float cannot end with a dot"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatWithTwoDots(t *testing.T) ***REMOVED***
	testFlow(t, "foo = 4.2.", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenError, "cannot have two dots in one float"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatWithExponent1(t *testing.T) ***REMOVED***
	testFlow(t, "a = 5e+22", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenFloat, "5e+22"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatWithExponent2(t *testing.T) ***REMOVED***
	testFlow(t, "a = 5E+22", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenFloat, "5E+22"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatWithExponent3(t *testing.T) ***REMOVED***
	testFlow(t, "a = -5e+22", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenFloat, "-5e+22"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatWithExponent4(t *testing.T) ***REMOVED***
	testFlow(t, "a = -5e-22", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenFloat, "-5e-22"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatWithExponent5(t *testing.T) ***REMOVED***
	testFlow(t, "a = 6.626e-34", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenFloat, "6.626e-34"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 14***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestInvalidEsquapeSequence(t *testing.T) ***REMOVED***
	testFlow(t, `foo = "\x"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenError, "invalid escape sequence: \\x"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestNestedArrays(t *testing.T) ***REMOVED***
	testFlow(t, "foo = [[[]]]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 9***REMOVED***, tokenLeftBracket, "["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenRightBracket, "]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 13***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualNumber(t *testing.T) ***REMOVED***
	testFlow(t, "foo = 42", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "42"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 9***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = +42", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "+42"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = -42", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "-42"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = 4.2", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenFloat, "4.2"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = +4.2", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenFloat, "+4.2"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = -4.2", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenFloat, "-4.2"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = 1_000", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "1_000"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 12***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = 5_349_221", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "5_349_221"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 16***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = 1_2_3_4_5", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "1_2_3_4_5"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 16***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "flt8 = 9_224_617.445_991_228_313", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "flt8"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenFloat, "9_224_617.445_991_228_313"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 33***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = +", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenError, "no digit in that number"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestMultiline(t *testing.T) ***REMOVED***
	testFlow(t, "foo = 42\nbar=21", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenInteger, "42"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenKey, "bar"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 4***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 5***REMOVED***, tokenInteger, "21"***REMOVED***,
		***REMOVED***Position***REMOVED***2, 7***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualStringUnicodeEscape(t *testing.T) ***REMOVED***
	testFlow(t, `foo = "hello \u2665"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "hello ♥"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 21***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = "hello \U000003B4"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "hello δ"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 25***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = "\uabcd"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "\uabcd"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 15***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = "\uABCD"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "\uABCD"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 15***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = "\U000bcdef"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "\U000bcdef"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 19***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = "\U000BCDEF"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "\U000BCDEF"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 19***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = "\u2"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenError, "unfinished unicode escape"***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = "\U2"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenError, "unfinished unicode escape"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyEqualStringNoEscape(t *testing.T) ***REMOVED***
	testFlow(t, "foo = \"hello \u0002\"", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenError, "unescaped control character U+0002"***REMOVED***,
	***REMOVED***)
	testFlow(t, "foo = \"hello \u001F\"", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenError, "unescaped control character U+001F"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLiteralString(t *testing.T) ***REMOVED***
	testFlow(t, `foo = 'C:\Users\nodejs\templates'`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, `C:\Users\nodejs\templates`***REMOVED***,
		***REMOVED***Position***REMOVED***1, 34***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = '\\ServerX\admin$\system32\'`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, `\\ServerX\admin$\system32\`***REMOVED***,
		***REMOVED***Position***REMOVED***1, 35***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = 'Tom "Dubs" Preston-Werner'`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, `Tom "Dubs" Preston-Werner`***REMOVED***,
		***REMOVED***Position***REMOVED***1, 34***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = '<\i\c*\s*>'`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, `<\i\c*\s*>`***REMOVED***,
		***REMOVED***Position***REMOVED***1, 19***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, `foo = 'C:\Users\nodejs\unfinis`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenError, "unclosed string"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestMultilineLiteralString(t *testing.T) ***REMOVED***
	testFlow(t, `foo = '''hello 'literal' world'''`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenString, `hello 'literal' world`***REMOVED***,
		***REMOVED***Position***REMOVED***1, 34***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = '''\nhello\n'literal'\nworld'''", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenString, "hello\n'literal'\nworld"***REMOVED***,
		***REMOVED***Position***REMOVED***4, 9***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
	testFlow(t, "foo = '''\r\nhello\r\n'literal'\r\nworld'''", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenString, "hello\r\n'literal'\r\nworld"***REMOVED***,
		***REMOVED***Position***REMOVED***4, 9***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestMultilineString(t *testing.T) ***REMOVED***
	testFlow(t, `foo = """hello "literal" world"""`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenString, `hello "literal" world`***REMOVED***,
		***REMOVED***Position***REMOVED***1, 34***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = \"\"\"\r\nhello\\\r\n\"literal\"\\\nworld\"\"\"", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenString, "hello\"literal\"world"***REMOVED***,
		***REMOVED***Position***REMOVED***4, 9***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "foo = \"\"\"\\\n    \\\n    \\\n    hello\\\nmultiline\\\nworld\"\"\"", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 10***REMOVED***, tokenString, "hellomultilineworld"***REMOVED***,
		***REMOVED***Position***REMOVED***6, 9***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "key2 = \"\"\"\nThe quick brown \\\n\n\n  fox jumps over \\\n    the lazy dog.\"\"\"", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "key2"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenString, "The quick brown fox jumps over the lazy dog."***REMOVED***,
		***REMOVED***Position***REMOVED***6, 21***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "key2 = \"\"\"\\\n       The quick brown \\\n       fox jumps over \\\n       the lazy dog.\\\n       \"\"\"", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "key2"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenString, "The quick brown fox jumps over the lazy dog."***REMOVED***,
		***REMOVED***Position***REMOVED***5, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, `key2 = "Roses are red\nViolets are blue"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "key2"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 9***REMOVED***, tokenString, "Roses are red\nViolets are blue"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 41***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)

	testFlow(t, "key2 = \"\"\"\nRoses are red\nViolets are blue\"\"\"", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "key2"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***2, 1***REMOVED***, tokenString, "Roses are red\nViolets are blue"***REMOVED***,
		***REMOVED***Position***REMOVED***3, 20***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestUnicodeString(t *testing.T) ***REMOVED***
	testFlow(t, `foo = "hello ♥ world"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "hello ♥ world"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 22***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***
func TestEscapeInString(t *testing.T) ***REMOVED***
	testFlow(t, `foo = "\b\f\/"`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenString, "\b\f/"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 15***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyGroupArray(t *testing.T) ***REMOVED***
	testFlow(t, "[[foo]]", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenDoubleLeftBracket, "[["***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenKeyGroupArray, "foo"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 6***REMOVED***, tokenDoubleRightBracket, "]]"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 8***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestQuotedKey(t *testing.T) ***REMOVED***
	testFlow(t, "\"a b\" = 42", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a b"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 9***REMOVED***, tokenInteger, "42"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 11***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestKeyNewline(t *testing.T) ***REMOVED***
	testFlow(t, "a\n= 4", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenError, "keys cannot contain new lines"***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestInvalidFloat(t *testing.T) ***REMOVED***
	testFlow(t, "a=7e1_", []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 2***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenFloat, "7e1_"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 7***REMOVED***, tokenEOF, ""***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestLexUnknownRvalue(t *testing.T) ***REMOVED***
	testFlow(t, `a = !b`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenError, "no value can start with !"***REMOVED***,
	***REMOVED***)

	testFlow(t, `a = \b`, []token***REMOVED***
		***REMOVED***Position***REMOVED***1, 1***REMOVED***, tokenKey, "a"***REMOVED***,
		***REMOVED***Position***REMOVED***1, 3***REMOVED***, tokenEqual, "="***REMOVED***,
		***REMOVED***Position***REMOVED***1, 5***REMOVED***, tokenError, `no value can start with \`***REMOVED***,
	***REMOVED***)
***REMOVED***

func BenchmarkLexer(b *testing.B) ***REMOVED***
	sample := `title = "Hugo: A Fast and Flexible Website Generator"
baseurl = "http://gohugo.io/"
MetaDataFormat = "yaml"
pluralizeListTitles = false

[params]
  description = "Documentation of Hugo, a fast and flexible static site generator built with love by spf13, bep and friends in Go"
  author = "Steve Francia (spf13) and friends"
  release = "0.22-DEV"

[[menu.main]]
	name = "Download Hugo"
	pre = "<i class='fa fa-download'></i>"
	url = "https://github.com/spf13/hugo/releases"
	weight = -200
`
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		lexToml([]byte(sample))
	***REMOVED***
***REMOVED***
