package scanner

import (
	"bytes"
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/hcl/hcl/token"
)

var f100 = strings.Repeat("f", 100)

type tokenPair struct ***REMOVED***
	tok  token.Type
	text string
***REMOVED***

var tokenLists = map[string][]tokenPair***REMOVED***
	"comment": []tokenPair***REMOVED***
		***REMOVED***token.COMMENT, "//"***REMOVED***,
		***REMOVED***token.COMMENT, "////"***REMOVED***,
		***REMOVED***token.COMMENT, "// comment"***REMOVED***,
		***REMOVED***token.COMMENT, "// /* comment */"***REMOVED***,
		***REMOVED***token.COMMENT, "// // comment //"***REMOVED***,
		***REMOVED***token.COMMENT, "//" + f100***REMOVED***,
		***REMOVED***token.COMMENT, "#"***REMOVED***,
		***REMOVED***token.COMMENT, "##"***REMOVED***,
		***REMOVED***token.COMMENT, "# comment"***REMOVED***,
		***REMOVED***token.COMMENT, "# /* comment */"***REMOVED***,
		***REMOVED***token.COMMENT, "# # comment #"***REMOVED***,
		***REMOVED***token.COMMENT, "#" + f100***REMOVED***,
		***REMOVED***token.COMMENT, "/**/"***REMOVED***,
		***REMOVED***token.COMMENT, "/***/"***REMOVED***,
		***REMOVED***token.COMMENT, "/* comment */"***REMOVED***,
		***REMOVED***token.COMMENT, "/* // comment */"***REMOVED***,
		***REMOVED***token.COMMENT, "/* /* comment */"***REMOVED***,
		***REMOVED***token.COMMENT, "/*\n comment\n*/"***REMOVED***,
		***REMOVED***token.COMMENT, "/*" + f100 + "*/"***REMOVED***,
	***REMOVED***,
	"operator": []tokenPair***REMOVED***
		***REMOVED***token.LBRACK, "["***REMOVED***,
		***REMOVED***token.LBRACE, "***REMOVED***"***REMOVED***,
		***REMOVED***token.COMMA, ","***REMOVED***,
		***REMOVED***token.PERIOD, "."***REMOVED***,
		***REMOVED***token.RBRACK, "]"***REMOVED***,
		***REMOVED***token.RBRACE, "***REMOVED***"***REMOVED***,
		***REMOVED***token.ASSIGN, "="***REMOVED***,
		***REMOVED***token.ADD, "+"***REMOVED***,
		***REMOVED***token.SUB, "-"***REMOVED***,
	***REMOVED***,
	"bool": []tokenPair***REMOVED***
		***REMOVED***token.BOOL, "true"***REMOVED***,
		***REMOVED***token.BOOL, "false"***REMOVED***,
	***REMOVED***,
	"ident": []tokenPair***REMOVED***
		***REMOVED***token.IDENT, "a"***REMOVED***,
		***REMOVED***token.IDENT, "a0"***REMOVED***,
		***REMOVED***token.IDENT, "foobar"***REMOVED***,
		***REMOVED***token.IDENT, "foo-bar"***REMOVED***,
		***REMOVED***token.IDENT, "abc123"***REMOVED***,
		***REMOVED***token.IDENT, "LGTM"***REMOVED***,
		***REMOVED***token.IDENT, "_"***REMOVED***,
		***REMOVED***token.IDENT, "_abc123"***REMOVED***,
		***REMOVED***token.IDENT, "abc123_"***REMOVED***,
		***REMOVED***token.IDENT, "_abc_123_"***REMOVED***,
		***REMOVED***token.IDENT, "_äöü"***REMOVED***,
		***REMOVED***token.IDENT, "_本"***REMOVED***,
		***REMOVED***token.IDENT, "äöü"***REMOVED***,
		***REMOVED***token.IDENT, "本"***REMOVED***,
		***REMOVED***token.IDENT, "a۰۱۸"***REMOVED***,
		***REMOVED***token.IDENT, "foo६४"***REMOVED***,
		***REMOVED***token.IDENT, "bar９８７６"***REMOVED***,
	***REMOVED***,
	"heredoc": []tokenPair***REMOVED***
		***REMOVED***token.HEREDOC, "<<EOF\nhello\nworld\nEOF"***REMOVED***,
		***REMOVED***token.HEREDOC, "<<EOF123\nhello\nworld\nEOF123"***REMOVED***,
	***REMOVED***,
	"string": []tokenPair***REMOVED***
		***REMOVED***token.STRING, `" "`***REMOVED***,
		***REMOVED***token.STRING, `"a"`***REMOVED***,
		***REMOVED***token.STRING, `"本"`***REMOVED***,
		***REMOVED***token.STRING, `"$***REMOVED***file("foo")***REMOVED***"`***REMOVED***,
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
		***REMOVED***token.NUMBER, "00"***REMOVED***,
		***REMOVED***token.NUMBER, "01"***REMOVED***,
		***REMOVED***token.NUMBER, "07"***REMOVED***,
		***REMOVED***token.NUMBER, "042"***REMOVED***,
		***REMOVED***token.NUMBER, "01234567"***REMOVED***,
		***REMOVED***token.NUMBER, "0x0"***REMOVED***,
		***REMOVED***token.NUMBER, "0x1"***REMOVED***,
		***REMOVED***token.NUMBER, "0xf"***REMOVED***,
		***REMOVED***token.NUMBER, "0x42"***REMOVED***,
		***REMOVED***token.NUMBER, "0x123456789abcDEF"***REMOVED***,
		***REMOVED***token.NUMBER, "0x" + f100***REMOVED***,
		***REMOVED***token.NUMBER, "0X0"***REMOVED***,
		***REMOVED***token.NUMBER, "0X1"***REMOVED***,
		***REMOVED***token.NUMBER, "0XF"***REMOVED***,
		***REMOVED***token.NUMBER, "0X42"***REMOVED***,
		***REMOVED***token.NUMBER, "0X123456789abcDEF"***REMOVED***,
		***REMOVED***token.NUMBER, "0X" + f100***REMOVED***,
		***REMOVED***token.NUMBER, "-0"***REMOVED***,
		***REMOVED***token.NUMBER, "-1"***REMOVED***,
		***REMOVED***token.NUMBER, "-9"***REMOVED***,
		***REMOVED***token.NUMBER, "-42"***REMOVED***,
		***REMOVED***token.NUMBER, "-1234567890"***REMOVED***,
		***REMOVED***token.NUMBER, "-00"***REMOVED***,
		***REMOVED***token.NUMBER, "-01"***REMOVED***,
		***REMOVED***token.NUMBER, "-07"***REMOVED***,
		***REMOVED***token.NUMBER, "-29"***REMOVED***,
		***REMOVED***token.NUMBER, "-042"***REMOVED***,
		***REMOVED***token.NUMBER, "-01234567"***REMOVED***,
		***REMOVED***token.NUMBER, "-0x0"***REMOVED***,
		***REMOVED***token.NUMBER, "-0x1"***REMOVED***,
		***REMOVED***token.NUMBER, "-0xf"***REMOVED***,
		***REMOVED***token.NUMBER, "-0x42"***REMOVED***,
		***REMOVED***token.NUMBER, "-0x123456789abcDEF"***REMOVED***,
		***REMOVED***token.NUMBER, "-0x" + f100***REMOVED***,
		***REMOVED***token.NUMBER, "-0X0"***REMOVED***,
		***REMOVED***token.NUMBER, "-0X1"***REMOVED***,
		***REMOVED***token.NUMBER, "-0XF"***REMOVED***,
		***REMOVED***token.NUMBER, "-0X42"***REMOVED***,
		***REMOVED***token.NUMBER, "-0X123456789abcDEF"***REMOVED***,
		***REMOVED***token.NUMBER, "-0X" + f100***REMOVED***,
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
	"ident",
	"heredoc",
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
			s.Scan()
		***REMOVED***
	***REMOVED***
	// make sure there were no token-internal errors reported by scanner
	if s.ErrorCount != 0 ***REMOVED***
		t.Errorf("%d errors", s.ErrorCount)
	***REMOVED***
***REMOVED***

func TestNullChar(t *testing.T) ***REMOVED***
	s := New([]byte("\"\\0"))
	s.Scan() // Used to panic
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

func TestWindowsLineEndings(t *testing.T) ***REMOVED***
	hcl := `// This should have Windows line endings
resource "aws_instance" "foo" ***REMOVED***
    user_data=<<HEREDOC
    test script
HEREDOC
***REMOVED***`
	hclWindowsEndings := strings.Replace(hcl, "\n", "\r\n", -1)

	literals := []struct ***REMOVED***
		tokenType token.Type
		literal   string
	***REMOVED******REMOVED***
		***REMOVED***token.COMMENT, "// This should have Windows line endings\r"***REMOVED***,
		***REMOVED***token.IDENT, `resource`***REMOVED***,
		***REMOVED***token.STRING, `"aws_instance"`***REMOVED***,
		***REMOVED***token.STRING, `"foo"`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `user_data`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.HEREDOC, "<<HEREDOC\r\n    test script\r\nHEREDOC\r\n"***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
	***REMOVED***

	s := New([]byte(hclWindowsEndings))
	for _, l := range literals ***REMOVED***
		tok := s.Scan()

		if l.tokenType != tok.Type ***REMOVED***
			t.Errorf("got: %s want %s for %s\n", tok, l.tokenType, tok.String())
		***REMOVED***

		if l.literal != tok.Text ***REMOVED***
			t.Errorf("got:\n%v\nwant:\n%v\n", []byte(tok.Text), []byte(l.literal))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRealExample(t *testing.T) ***REMOVED***
	complexHCL := `// This comes from Terraform, as a test
	variable "foo" ***REMOVED***
	    default = "bar"
	    description = "bar"
	***REMOVED***

	provider "aws" ***REMOVED***
	  access_key = "foo"
	  secret_key = "$***REMOVED***replace(var.foo, ".", "\\.")***REMOVED***"
	***REMOVED***

	resource "aws_security_group" "firewall" ***REMOVED***
	    count = 5
	***REMOVED***

	resource aws_instance "web" ***REMOVED***
	    ami = "$***REMOVED***var.foo***REMOVED***"
	    security_groups = [
	        "foo",
	        "$***REMOVED***aws_security_group.firewall.foo***REMOVED***"
	    ]

	    network_interface ***REMOVED***
	        device_index = 0
	        description = <<EOF
Main interface
EOF
	***REMOVED***

		network_interface ***REMOVED***
	        device_index = 1
	        description = <<-EOF
			Outer text
				Indented text
			EOF
		***REMOVED***
	***REMOVED***`

	literals := []struct ***REMOVED***
		tokenType token.Type
		literal   string
	***REMOVED******REMOVED***
		***REMOVED***token.COMMENT, `// This comes from Terraform, as a test`***REMOVED***,
		***REMOVED***token.IDENT, `variable`***REMOVED***,
		***REMOVED***token.STRING, `"foo"`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `default`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.STRING, `"bar"`***REMOVED***,
		***REMOVED***token.IDENT, `description`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.STRING, `"bar"`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `provider`***REMOVED***,
		***REMOVED***token.STRING, `"aws"`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `access_key`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.STRING, `"foo"`***REMOVED***,
		***REMOVED***token.IDENT, `secret_key`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.STRING, `"$***REMOVED***replace(var.foo, ".", "\\.")***REMOVED***"`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `resource`***REMOVED***,
		***REMOVED***token.STRING, `"aws_security_group"`***REMOVED***,
		***REMOVED***token.STRING, `"firewall"`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `count`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.NUMBER, `5`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `resource`***REMOVED***,
		***REMOVED***token.IDENT, `aws_instance`***REMOVED***,
		***REMOVED***token.STRING, `"web"`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `ami`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.STRING, `"$***REMOVED***var.foo***REMOVED***"`***REMOVED***,
		***REMOVED***token.IDENT, `security_groups`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.LBRACK, `[`***REMOVED***,
		***REMOVED***token.STRING, `"foo"`***REMOVED***,
		***REMOVED***token.COMMA, `,`***REMOVED***,
		***REMOVED***token.STRING, `"$***REMOVED***aws_security_group.firewall.foo***REMOVED***"`***REMOVED***,
		***REMOVED***token.RBRACK, `]`***REMOVED***,
		***REMOVED***token.IDENT, `network_interface`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `device_index`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.NUMBER, `0`***REMOVED***,
		***REMOVED***token.IDENT, `description`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.HEREDOC, "<<EOF\nMain interface\nEOF\n"***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `network_interface`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `device_index`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.NUMBER, `1`***REMOVED***,
		***REMOVED***token.IDENT, `description`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.HEREDOC, "<<-EOF\n\t\t\tOuter text\n\t\t\t\tIndented text\n\t\t\tEOF\n"***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.EOF, ``***REMOVED***,
	***REMOVED***

	s := New([]byte(complexHCL))
	for _, l := range literals ***REMOVED***
		tok := s.Scan()
		if l.tokenType != tok.Type ***REMOVED***
			t.Errorf("got: %s want %s for %s\n", tok, l.tokenType, tok.String())
		***REMOVED***

		if l.literal != tok.Text ***REMOVED***
			t.Errorf("got:\n%+v\n%s\n want:\n%+v\n%s\n", []byte(tok.String()), tok, []byte(l.literal), l.literal)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestScan_crlf(t *testing.T) ***REMOVED***
	complexHCL := "foo ***REMOVED***\r\n  bar = \"baz\"\r\n***REMOVED***\r\n"

	literals := []struct ***REMOVED***
		tokenType token.Type
		literal   string
	***REMOVED******REMOVED***
		***REMOVED***token.IDENT, `foo`***REMOVED***,
		***REMOVED***token.LBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.IDENT, `bar`***REMOVED***,
		***REMOVED***token.ASSIGN, `=`***REMOVED***,
		***REMOVED***token.STRING, `"baz"`***REMOVED***,
		***REMOVED***token.RBRACE, `***REMOVED***`***REMOVED***,
		***REMOVED***token.EOF, ``***REMOVED***,
	***REMOVED***

	s := New([]byte(complexHCL))
	for _, l := range literals ***REMOVED***
		tok := s.Scan()
		if l.tokenType != tok.Type ***REMOVED***
			t.Errorf("got: %s want %s for %s\n", tok, l.tokenType, tok.String())
		***REMOVED***

		if l.literal != tok.Text ***REMOVED***
			t.Errorf("got:\n%+v\n%s\n want:\n%+v\n%s\n", []byte(tok.String()), tok, []byte(l.literal), l.literal)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestError(t *testing.T) ***REMOVED***
	testError(t, "\x80", "1:1", "illegal UTF-8 encoding", token.ILLEGAL)
	testError(t, "\xff", "1:1", "illegal UTF-8 encoding", token.ILLEGAL)

	testError(t, "ab\x80", "1:3", "illegal UTF-8 encoding", token.IDENT)
	testError(t, "abc\xff", "1:4", "illegal UTF-8 encoding", token.IDENT)

	testError(t, `"ab`+"\x80", "1:4", "illegal UTF-8 encoding", token.STRING)
	testError(t, `"abc`+"\xff", "1:5", "illegal UTF-8 encoding", token.STRING)

	testError(t, `01238`, "1:6", "illegal octal number", token.NUMBER)
	testError(t, `01238123`, "1:9", "illegal octal number", token.NUMBER)
	testError(t, `0x`, "1:3", "illegal hexadecimal number", token.NUMBER)
	testError(t, `0xg`, "1:3", "illegal hexadecimal number", token.NUMBER)
	testError(t, `'aa'`, "1:1", "illegal char", token.ILLEGAL)

	testError(t, `"`, "1:2", "literal not terminated", token.STRING)
	testError(t, `"abc`, "1:5", "literal not terminated", token.STRING)
	testError(t, `"abc`+"\n", "1:5", "literal not terminated", token.STRING)
	testError(t, `"$***REMOVED***abc`+"\n", "2:1", "literal not terminated", token.STRING)
	testError(t, `/*/`, "1:4", "comment not terminated", token.COMMENT)
	testError(t, `/foo`, "1:1", "expected '/' for comment", token.COMMENT)
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
