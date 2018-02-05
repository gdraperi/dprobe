package scanner

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/json/token"
)

var f100 = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

type tokenPair struct ***REMOVED***
	tok  token.Type
	text string
***REMOVED***

var tokenLists = map[string][]tokenPair***REMOVED***
	"operator": []tokenPair***REMOVED***
		***REMOVED***token.LBRACK, "["***REMOVED***,
		***REMOVED***token.LBRACE, "***REMOVED***"***REMOVED***,
		***REMOVED***token.COMMA, ","***REMOVED***,
		***REMOVED***token.PERIOD, "."***REMOVED***,
		***REMOVED***token.RBRACK, "]"***REMOVED***,
		***REMOVED***token.RBRACE, "***REMOVED***"***REMOVED***,
	***REMOVED***,
	"bool": []tokenPair***REMOVED***
		***REMOVED***token.BOOL, "true"***REMOVED***,
		***REMOVED***token.BOOL, "false"***REMOVED***,
	***REMOVED***,
	"string": []tokenPair***REMOVED***
		***REMOVED***token.STRING, `" "`***REMOVED***,
		***REMOVED***token.STRING, `"a"`***REMOVED***,
		***REMOVED***token.STRING, `"本"`***REMOVED***,
		***REMOVED***token.STRING, `"$***REMOVED***file(\"foo\")***REMOVED***"`***REMOVED***,
		***REMOVED***token.STRING, `"\a"`***REMOVED***,
		***REMOVED***token.STRING, `"\b"`***REMOVED***,
		***REMOVED***token.STRING, `"\f"`***REMOVED***,
		***REMOVED***token.STRING, `"\n"`***REMOVED***,
		***REMOVED***token.STRING, `"\r"`***REMOVED***,
		***REMOVED***token.STRING, `"\t"`***REMOVED***,
		***REMOVED***token.STRING, `"\v"`***REMOVED***,
		***REMOVED***token.STRING, `"\""`***REMOVED***,
		***REMOVED***token.STRING, `"\000"`***REMOVED***,
		***REMOVED***token.STRING, `"\777"`***REMOVED***,
		***REMOVED***token.STRING, `"\x00"`***REMOVED***,
		***REMOVED***token.STRING, `"\xff"`***REMOVED***,
		***REMOVED***token.STRING, `"\u0000"`***REMOVED***,
		***REMOVED***token.STRING, `"\ufA16"`***REMOVED***,
		***REMOVED***token.STRING, `"\U00000000"`***REMOVED***,
		***REMOVED***token.STRING, `"\U0000ffAB"`***REMOVED***,
		***REMOVED***token.STRING, `"` + f100 + `"`***REMOVED***,
	***REMOVED***,
	"number": []tokenPair***REMOVED***
		***REMOVED***token.NUMBER, "0"***REMOVED***,
		***REMOVED***token.NUMBER, "1"***REMOVED***,
		***REMOVED***token.NUMBER, "9"***REMOVED***,
		***REMOVED***token.NUMBER, "42"***REMOVED***,
		***REMOVED***token.NUMBER, "1234567890"***REMOVED***,
		***REMOVED***token.NUMBER, "-0"***REMOVED***,
		***REMOVED***token.NUMBER, "-1"***REMOVED***,
		***REMOVED***token.NUMBER, "-9"***REMOVED***,
		***REMOVED***token.NUMBER, "-42"***REMOVED***,
		***REMOVED***token.NUMBER, "-1234567890"***REMOVED***,
	***REMOVED***,
	"float": []tokenPair***REMOVED***
		***REMOVED***token.FLOAT, "0."***REMOVED***,
		***REMOVED***token.FLOAT, "1."***REMOVED***,
		***REMOVED***token.FLOAT, "42."***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890."***REMOVED***,
		***REMOVED***token.FLOAT, ".0"***REMOVED***,
		***REMOVED***token.FLOAT, ".1"***REMOVED***,
		***REMOVED***token.FLOAT, ".42"***REMOVED***,
		***REMOVED***token.FLOAT, ".0123456789"***REMOVED***,
		***REMOVED***token.FLOAT, "0.0"***REMOVED***,
		***REMOVED***token.FLOAT, "1.0"***REMOVED***,
		***REMOVED***token.FLOAT, "42.0"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890.0"***REMOVED***,
		***REMOVED***token.FLOAT, "0e0"***REMOVED***,
		***REMOVED***token.FLOAT, "1e0"***REMOVED***,
		***REMOVED***token.FLOAT, "42e0"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890e0"***REMOVED***,
		***REMOVED***token.FLOAT, "0E0"***REMOVED***,
		***REMOVED***token.FLOAT, "1E0"***REMOVED***,
		***REMOVED***token.FLOAT, "42E0"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890E0"***REMOVED***,
		***REMOVED***token.FLOAT, "0e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "1e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "42e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "0E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "1E-10"***REMOVED***,
		***REMOVED***token.FLOAT, "42E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890E-10"***REMOVED***,
		***REMOVED***token.FLOAT, "01.8e0"***REMOVED***,
		***REMOVED***token.FLOAT, "1.4e0"***REMOVED***,
		***REMOVED***token.FLOAT, "42.2e0"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890.12e0"***REMOVED***,
		***REMOVED***token.FLOAT, "0.E0"***REMOVED***,
		***REMOVED***token.FLOAT, "1.12E0"***REMOVED***,
		***REMOVED***token.FLOAT, "42.123E0"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890.213E0"***REMOVED***,
		***REMOVED***token.FLOAT, "0.2e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "1.2e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "42.54e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890.98e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "0.1E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "1.1E-10"***REMOVED***,
		***REMOVED***token.FLOAT, "42.1E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "01234567890.1E-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-0.0"***REMOVED***,
		***REMOVED***token.FLOAT, "-1.0"***REMOVED***,
		***REMOVED***token.FLOAT, "-42.0"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890.0"***REMOVED***,
		***REMOVED***token.FLOAT, "-0e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-1e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-42e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-0E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-1E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-42E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-0e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-1e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-42e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-0E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-1E-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-42E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890E-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-01.8e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-1.4e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-42.2e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890.12e0"***REMOVED***,
		***REMOVED***token.FLOAT, "-0.E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-1.12E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-42.123E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890.213E0"***REMOVED***,
		***REMOVED***token.FLOAT, "-0.2e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-1.2e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-42.54e+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890.98e-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-0.1E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-1.1E-10"***REMOVED***,
		***REMOVED***token.FLOAT, "-42.1E+10"***REMOVED***,
		***REMOVED***token.FLOAT, "-01234567890.1E-10"***REMOVED***,
	***REMOVED***,
***REMOVED***

var orderedTokenLists = []string***REMOVED***
	"comment",
	"operator",
	"bool",
	"string",
	"number",
	"float",
***REMOVED***

func TestPosition(t *testing.T) ***REMOVED***
	// create artifical source code
	buf := new(bytes.Buffer)

	for _, listName := range orderedTokenLists ***REMOVED***
		for _, ident := range tokenLists[listName] ***REMOVED***
			fmt.Fprintf(buf, "\t\t\t\t%s\n", ident.text)
		***REMOVED***
	***REMOVED***

	s := New(buf.Bytes())

	pos := token.Pos***REMOVED***"", 4, 1, 5***REMOVED***
	s.Scan()
	for _, listName := range orderedTokenLists ***REMOVED***

		for _, k := range tokenLists[listName] ***REMOVED***
			curPos := s.tokPos
			// fmt.Printf("[%q] s = %+v:%+v\n", k.text, curPos.Offset, curPos.Column)

			if curPos.Offset != pos.Offset ***REMOVED***
				t.Fatalf("offset = %d, want %d for %q", curPos.Offset, pos.Offset, k.text)
			***REMOVED***
			if curPos.Line != pos.Line ***REMOVED***
				t.Fatalf("line = %d, want %d for %q", curPos.Line, pos.Line, k.text)
			***REMOVED***
			if curPos.Column != pos.Column ***REMOVED***
				t.Fatalf("column = %d, want %d for %q", curPos.Column, pos.Column, k.text)
			***REMOVED***
			pos.Offset += 4 + len(k.text) + 1     // 4 tabs + token bytes + newline
			pos.Line += countNewlines(k.text) + 1 // each token is on a new line

			s.Error = func(pos token.Pos, msg string) ***REMOVED***
				t.Errorf("error %q for %q", msg, k.text)
			***REMOVED***

			s.Scan()
		***REMOVED***
	***REMOVED***
	// make sure there were no token-internal errors reported by scanner
	if s.ErrorCount != 0 ***REMOVED***
		t.Errorf("%d errors", s.ErrorCount)
	***REMOVED***
***REMOVED***

func TestComment(t *testing.T) ***REMOVED***
	testTokenList(t, tokenLists["comment"])
***REMOVED***

func TestOperator(t *testing.T) ***REMOVED***
	testTokenList(t, tokenLists["operator"])
***REMOVED***

func TestBool(t *testing.T) ***REMOVED***
	testTokenList(t, tokenLists["bool"])
***REMOVED***

func TestIdent(t *testing.T) ***REMOVED***
	testTokenList(t, tokenLists["ident"])
***REMOVED***

func TestString(t *testing.T) ***REMOVED***
	testTokenList(t, tokenLists["string"])
***REMOVED***

func TestNumber(t *testing.T) ***REMOVED***
	testTokenList(t, tokenLists["number"])
***REMOVED***

func TestFloat(t *testing.T) ***REMOVED***
	testTokenList(t, tokenLists["float"])
***REMOVED***

func TestRealExample(t *testing.T) ***REMOVED***
	complexReal := `
***REMOVED***
    "variable": ***REMOVED***
        "foo": ***REMOVED***
            "default": "bar",
            "description": "bar",
            "depends_on": ["something"]
    ***REMOVED***
***REMOVED***
***REMOVED***`

	literals := []struct ***REMOVED***
		tokenType token.Type
		literal   string
	***REMOVED******REMOVED***
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.STRING, `"variable"`***REMOVED***,
		***REMOVED***token.COLON, `:`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.STRING, `"foo"`***REMOVED***,
		***REMOVED***token.COLON, `:`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.STRING, `"default"`***REMOVED***,
		***REMOVED***token.COLON, `:`***REMOVED***,
		***REMOVED***token.STRING, `"bar"`***REMOVED***,
		***REMOVED***token.COMMA, `,`***REMOVED***,
		***REMOVED***token.STRING, `"description"`***REMOVED***,
		***REMOVED***token.COLON, `:`***REMOVED***,
		***REMOVED***token.STRING, `"bar"`***REMOVED***,
		***REMOVED***token.COMMA, `,`***REMOVED***,
		***REMOVED***token.STRING, `"depends_on"`***REMOVED***,
		***REMOVED***token.COLON, `:`***REMOVED***,
		***REMOVED***token.LBRACK, `[`***REMOVED***,
		***REMOVED***token.STRING, `"something"`***REMOVED***,
		***REMOVED***token.RBRACK, `]`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.EOF, ``***REMOVED***,
	***REMOVED***

	s := New([]byte(complexReal))
	for _, l := range literals ***REMOVED***
		tok := s.Scan()
		if l.tokenType != tok.Type ***REMOVED***
			t.Errorf("got: %s want %s for %s\n", tok, l.tokenType, tok.String())
		***REMOVED***

		if l.literal != tok.Text ***REMOVED***
			t.Errorf("got: %s want %s\n", tok, l.literal)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestError(t *testing.T) ***REMOVED***
	testError(t, "\x80", "1:1", "illegal UTF-8 encoding", token.ILLEGAL)
	testError(t, "\xff", "1:1", "illegal UTF-8 encoding", token.ILLEGAL)

	testError(t, `"ab`+"\x80", "1:4", "illegal UTF-8 encoding", token.STRING)
	testError(t, `"abc`+"\xff", "1:5", "illegal UTF-8 encoding", token.STRING)

	testError(t, `01238`, "1:7", "numbers cannot start with 0", token.NUMBER)
	testError(t, `01238123`, "1:10", "numbers cannot start with 0", token.NUMBER)
	testError(t, `'aa'`, "1:1", "illegal char: '", token.ILLEGAL)

	testError(t, `"`, "1:2", "literal not terminated", token.STRING)
	testError(t, `"abc`, "1:5", "literal not terminated", token.STRING)
	testError(t, `"abc`+"\n", "1:5", "literal not terminated", token.STRING)
***REMOVED***

func testError(t *testing.T, src, pos, msg string, tok token.Type) ***REMOVED***
	s := New([]byte(src))

	errorCalled := false
	s.Error = func(p token.Pos, m string) ***REMOVED***
		if !errorCalled ***REMOVED***
			if pos != p.String() ***REMOVED***
				t.Errorf("pos = %q, want %q for %q", p, pos, src)
			***REMOVED***

			if m != msg ***REMOVED***
				t.Errorf("msg = %q, want %q for %q", m, msg, src)
			***REMOVED***
			errorCalled = true
		***REMOVED***
	***REMOVED***

	tk := s.Scan()
	if tk.Type != tok ***REMOVED***
		t.Errorf("tok = %s, want %s for %q", tk, tok, src)
	***REMOVED***
	if !errorCalled ***REMOVED***
		t.Errorf("error handler not called for %q", src)
	***REMOVED***
	if s.ErrorCount == 0 ***REMOVED***
		t.Errorf("count = %d, want > 0 for %q", s.ErrorCount, src)
	***REMOVED***
***REMOVED***

func testTokenList(t *testing.T, tokenList []tokenPair) ***REMOVED***
	// create artifical source code
	buf := new(bytes.Buffer)
	for _, ident := range tokenList ***REMOVED***
		fmt.Fprintf(buf, "%s\n", ident.text)
	***REMOVED***

	s := New(buf.Bytes())
	for _, ident := range tokenList ***REMOVED***
		tok := s.Scan()
		if tok.Type != ident.tok ***REMOVED***
			t.Errorf("tok = %q want %q for %q\n", tok, ident.tok, ident.text)
		***REMOVED***

		if tok.Text != ident.text ***REMOVED***
			t.Errorf("text = %q want %q", tok.String(), ident.text)
		***REMOVED***

	***REMOVED***
***REMOVED***

func countNewlines(s string) int ***REMOVED***
	n := 0
	for _, ch := range s ***REMOVED***
		if ch == '\n' ***REMOVED***
			n++
		***REMOVED***
	***REMOVED***
	return n
***REMOVED***
