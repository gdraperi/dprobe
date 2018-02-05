package strconv

import "testing"

type quoteTest struct ***REMOVED***
	in    string
	out   string
	ascii string
***REMOVED***

var quotetests = []quoteTest***REMOVED***
	***REMOVED***"\a\b\f\r\n\t\v", `"\a\b\f\r\n\t\v"`, `"\a\b\f\r\n\t\v"`***REMOVED***,
	***REMOVED***"\\", `"\\"`, `"\\"`***REMOVED***,
	***REMOVED***"abc\xffdef", `"abc\xffdef"`, `"abc\xffdef"`***REMOVED***,
	***REMOVED***"\u263a", `"☺"`, `"\u263a"`***REMOVED***,
	***REMOVED***"\U0010ffff", `"\U0010ffff"`, `"\U0010ffff"`***REMOVED***,
	***REMOVED***"\x04", `"\x04"`, `"\x04"`***REMOVED***,
***REMOVED***

type unQuoteTest struct ***REMOVED***
	in  string
	out string
***REMOVED***

var unquotetests = []unQuoteTest***REMOVED***
	***REMOVED***`""`, ""***REMOVED***,
	***REMOVED***`"a"`, "a"***REMOVED***,
	***REMOVED***`"abc"`, "abc"***REMOVED***,
	***REMOVED***`"☺"`, "☺"***REMOVED***,
	***REMOVED***`"hello world"`, "hello world"***REMOVED***,
	***REMOVED***`"\xFF"`, "\xFF"***REMOVED***,
	***REMOVED***`"\377"`, "\377"***REMOVED***,
	***REMOVED***`"\u1234"`, "\u1234"***REMOVED***,
	***REMOVED***`"\U00010111"`, "\U00010111"***REMOVED***,
	***REMOVED***`"\U0001011111"`, "\U0001011111"***REMOVED***,
	***REMOVED***`"\a\b\f\n\r\t\v\\\""`, "\a\b\f\n\r\t\v\\\""***REMOVED***,
	***REMOVED***`"'"`, "'"***REMOVED***,
	***REMOVED***`"$***REMOVED***file("foo")***REMOVED***"`, `$***REMOVED***file("foo")***REMOVED***`***REMOVED***,
	***REMOVED***`"$***REMOVED***file("\"foo\"")***REMOVED***"`, `$***REMOVED***file("\"foo\"")***REMOVED***`***REMOVED***,
	***REMOVED***`"echo $***REMOVED***var.region***REMOVED***$***REMOVED***element(split(",",var.zones),0)***REMOVED***"`,
		`echo $***REMOVED***var.region***REMOVED***$***REMOVED***element(split(",",var.zones),0)***REMOVED***`***REMOVED***,
	***REMOVED***`"$***REMOVED***HH\\:mm\\:ss***REMOVED***"`, `$***REMOVED***HH\\:mm\\:ss***REMOVED***`***REMOVED***,
	***REMOVED***`"$***REMOVED***\n***REMOVED***"`, `$***REMOVED***\n***REMOVED***`***REMOVED***,
***REMOVED***

var misquoted = []string***REMOVED***
	``,
	`"`,
	`"a`,
	`"'`,
	`b"`,
	`"\"`,
	`"\9"`,
	`"\19"`,
	`"\129"`,
	`'\'`,
	`'\9'`,
	`'\19'`,
	`'\129'`,
	`'ab'`,
	`"\x1!"`,
	`"\U12345678"`,
	`"\z"`,
	"`",
	"`xxx",
	"`\"",
	`"\'"`,
	`'\"'`,
	"\"\n\"",
	"\"\\n\n\"",
	"'\n'",
	`"$***REMOVED***"`,
	`"$***REMOVED***foo***REMOVED******REMOVED***"`,
	"\"$***REMOVED***foo***REMOVED***\n\"",
***REMOVED***

func TestUnquote(t *testing.T) ***REMOVED***
	for _, tt := range unquotetests ***REMOVED***
		if out, err := Unquote(tt.in); err != nil || out != tt.out ***REMOVED***
			t.Errorf("Unquote(%#q) = %q, %v want %q, nil", tt.in, out, err, tt.out)
		***REMOVED***
	***REMOVED***

	// run the quote tests too, backward
	for _, tt := range quotetests ***REMOVED***
		if in, err := Unquote(tt.out); in != tt.in ***REMOVED***
			t.Errorf("Unquote(%#q) = %q, %v, want %q, nil", tt.out, in, err, tt.in)
		***REMOVED***
	***REMOVED***

	for _, s := range misquoted ***REMOVED***
		if out, err := Unquote(s); out != "" || err != ErrSyntax ***REMOVED***
			t.Errorf("Unquote(%#q) = %q, %v want %q, %v", s, out, err, "", ErrSyntax)
		***REMOVED***
	***REMOVED***
***REMOVED***
