package logrus

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFormatting(t *testing.T) ***REMOVED***
	tf := &TextFormatter***REMOVED***DisableColors: true***REMOVED***

	testCases := []struct ***REMOVED***
		value    string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***`foo`, "time=\"0001-01-01T00:00:00Z\" level=panic test=foo\n"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		b, _ := tf.Format(WithField("test", tc.value))

		if string(b) != tc.expected ***REMOVED***
			t.Errorf("formatting expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestQuoting(t *testing.T) ***REMOVED***
	tf := &TextFormatter***REMOVED***DisableColors: true***REMOVED***

	checkQuoting := func(q bool, value interface***REMOVED******REMOVED***) ***REMOVED***
		b, _ := tf.Format(WithField("test", value))
		idx := bytes.Index(b, ([]byte)("test="))
		cont := bytes.Contains(b[idx+5:], []byte("\""))
		if cont != q ***REMOVED***
			if q ***REMOVED***
				t.Errorf("quoting expected for: %#v", value)
			***REMOVED*** else ***REMOVED***
				t.Errorf("quoting not expected for: %#v", value)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	checkQuoting(false, "")
	checkQuoting(false, "abcd")
	checkQuoting(false, "v1.0")
	checkQuoting(false, "1234567890")
	checkQuoting(false, "/foobar")
	checkQuoting(false, "foo_bar")
	checkQuoting(false, "foo@bar")
	checkQuoting(false, "foobar^")
	checkQuoting(false, "+/-_^@f.oobar")
	checkQuoting(true, "foobar$")
	checkQuoting(true, "&foobar")
	checkQuoting(true, "x y")
	checkQuoting(true, "x,y")
	checkQuoting(false, errors.New("invalid"))
	checkQuoting(true, errors.New("invalid argument"))

	// Test for quoting empty fields.
	tf.QuoteEmptyFields = true
	checkQuoting(true, "")
	checkQuoting(false, "abcd")
	checkQuoting(true, errors.New("invalid argument"))
***REMOVED***

func TestEscaping(t *testing.T) ***REMOVED***
	tf := &TextFormatter***REMOVED***DisableColors: true***REMOVED***

	testCases := []struct ***REMOVED***
		value    string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***`ba"r`, `ba\"r`***REMOVED***,
		***REMOVED***`ba'r`, `ba'r`***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		b, _ := tf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) ***REMOVED***
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEscaping_Interface(t *testing.T) ***REMOVED***
	tf := &TextFormatter***REMOVED***DisableColors: true***REMOVED***

	ts := time.Now()

	testCases := []struct ***REMOVED***
		value    interface***REMOVED******REMOVED***
		expected string
	***REMOVED******REMOVED***
		***REMOVED***ts, fmt.Sprintf("\"%s\"", ts.String())***REMOVED***,
		***REMOVED***errors.New("error: something went wrong"), "\"error: something went wrong\""***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		b, _ := tf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) ***REMOVED***
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTimestampFormat(t *testing.T) ***REMOVED***
	checkTimeStr := func(format string) ***REMOVED***
		customFormatter := &TextFormatter***REMOVED***DisableColors: true, TimestampFormat: format***REMOVED***
		customStr, _ := customFormatter.Format(WithField("test", "test"))
		timeStart := bytes.Index(customStr, ([]byte)("time="))
		timeEnd := bytes.Index(customStr, ([]byte)("level="))
		timeStr := customStr[timeStart+5+len("\"") : timeEnd-1-len("\"")]
		if format == "" ***REMOVED***
			format = time.RFC3339
		***REMOVED***
		_, e := time.Parse(format, (string)(timeStr))
		if e != nil ***REMOVED***
			t.Errorf("time string \"%s\" did not match provided time format \"%s\": %s", timeStr, format, e)
		***REMOVED***
	***REMOVED***

	checkTimeStr("2006-01-02T15:04:05.000000000Z07:00")
	checkTimeStr("Mon Jan _2 15:04:05 2006")
	checkTimeStr("")
***REMOVED***

func TestDisableTimestampWithColoredOutput(t *testing.T) ***REMOVED***
	tf := &TextFormatter***REMOVED***DisableTimestamp: true, ForceColors: true***REMOVED***

	b, _ := tf.Format(WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") ***REMOVED***
		t.Error("timestamp not expected when DisableTimestamp is true")
	***REMOVED***
***REMOVED***

// TODO add tests for sorting etc., this requires a parser for the text
// formatter output.
