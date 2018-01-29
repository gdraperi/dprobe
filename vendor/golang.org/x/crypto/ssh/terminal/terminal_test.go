// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"bytes"
	"io"
	"os"
	"testing"
)

type MockTerminal struct ***REMOVED***
	toSend       []byte
	bytesPerRead int
	received     []byte
***REMOVED***

func (c *MockTerminal) Read(data []byte) (n int, err error) ***REMOVED***
	n = len(data)
	if n == 0 ***REMOVED***
		return
	***REMOVED***
	if n > len(c.toSend) ***REMOVED***
		n = len(c.toSend)
	***REMOVED***
	if n == 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	if c.bytesPerRead > 0 && n > c.bytesPerRead ***REMOVED***
		n = c.bytesPerRead
	***REMOVED***
	copy(data, c.toSend[:n])
	c.toSend = c.toSend[n:]
	return
***REMOVED***

func (c *MockTerminal) Write(data []byte) (n int, err error) ***REMOVED***
	c.received = append(c.received, data...)
	return len(data), nil
***REMOVED***

func TestClose(t *testing.T) ***REMOVED***
	c := &MockTerminal***REMOVED******REMOVED***
	ss := NewTerminal(c, "> ")
	line, err := ss.ReadLine()
	if line != "" ***REMOVED***
		t.Errorf("Expected empty line but got: %s", line)
	***REMOVED***
	if err != io.EOF ***REMOVED***
		t.Errorf("Error should have been EOF but got: %s", err)
	***REMOVED***
***REMOVED***

var keyPressTests = []struct ***REMOVED***
	in             string
	line           string
	err            error
	throwAwayLines int
***REMOVED******REMOVED***
	***REMOVED***
		err: io.EOF,
	***REMOVED***,
	***REMOVED***
		in:   "\r",
		line: "",
	***REMOVED***,
	***REMOVED***
		in:   "foo\r",
		line: "foo",
	***REMOVED***,
	***REMOVED***
		in:   "a\x1b[Cb\r", // right
		line: "ab",
	***REMOVED***,
	***REMOVED***
		in:   "a\x1b[Db\r", // left
		line: "ba",
	***REMOVED***,
	***REMOVED***
		in:   "a\177b\r", // backspace
		line: "b",
	***REMOVED***,
	***REMOVED***
		in: "\x1b[A\r", // up
	***REMOVED***,
	***REMOVED***
		in: "\x1b[B\r", // down
	***REMOVED***,
	***REMOVED***
		in:   "line\x1b[A\x1b[B\r", // up then down
		line: "line",
	***REMOVED***,
	***REMOVED***
		in:             "line1\rline2\x1b[A\r", // recall previous line.
		line:           "line1",
		throwAwayLines: 1,
	***REMOVED***,
	***REMOVED***
		// recall two previous lines and append.
		in:             "line1\rline2\rline3\x1b[A\x1b[Axxx\r",
		line:           "line1xxx",
		throwAwayLines: 2,
	***REMOVED***,
	***REMOVED***
		// Ctrl-A to move to beginning of line followed by ^K to kill
		// line.
		in:   "a b \001\013\r",
		line: "",
	***REMOVED***,
	***REMOVED***
		// Ctrl-A to move to beginning of line, Ctrl-E to move to end,
		// finally ^K to kill nothing.
		in:   "a b \001\005\013\r",
		line: "a b ",
	***REMOVED***,
	***REMOVED***
		in:   "\027\r",
		line: "",
	***REMOVED***,
	***REMOVED***
		in:   "a\027\r",
		line: "",
	***REMOVED***,
	***REMOVED***
		in:   "a \027\r",
		line: "",
	***REMOVED***,
	***REMOVED***
		in:   "a b\027\r",
		line: "a ",
	***REMOVED***,
	***REMOVED***
		in:   "a b \027\r",
		line: "a ",
	***REMOVED***,
	***REMOVED***
		in:   "one two thr\x1b[D\027\r",
		line: "one two r",
	***REMOVED***,
	***REMOVED***
		in:   "\013\r",
		line: "",
	***REMOVED***,
	***REMOVED***
		in:   "a\013\r",
		line: "a",
	***REMOVED***,
	***REMOVED***
		in:   "ab\x1b[D\013\r",
		line: "a",
	***REMOVED***,
	***REMOVED***
		in:   "Ξεσκεπάζω\r",
		line: "Ξεσκεπάζω",
	***REMOVED***,
	***REMOVED***
		in:             "£\r\x1b[A\177\r", // non-ASCII char, enter, up, backspace.
		line:           "",
		throwAwayLines: 1,
	***REMOVED***,
	***REMOVED***
		in:             "£\r££\x1b[A\x1b[B\177\r", // non-ASCII char, enter, 2x non-ASCII, up, down, backspace, enter.
		line:           "£",
		throwAwayLines: 1,
	***REMOVED***,
	***REMOVED***
		// Ctrl-D at the end of the line should be ignored.
		in:   "a\004\r",
		line: "a",
	***REMOVED***,
	***REMOVED***
		// a, b, left, Ctrl-D should erase the b.
		in:   "ab\x1b[D\004\r",
		line: "a",
	***REMOVED***,
	***REMOVED***
		// a, b, c, d, left, left, ^U should erase to the beginning of
		// the line.
		in:   "abcd\x1b[D\x1b[D\025\r",
		line: "cd",
	***REMOVED***,
	***REMOVED***
		// Bracketed paste mode: control sequences should be returned
		// verbatim in paste mode.
		in:   "abc\x1b[200~de\177f\x1b[201~\177\r",
		line: "abcde\177",
	***REMOVED***,
	***REMOVED***
		// Enter in bracketed paste mode should still work.
		in:             "abc\x1b[200~d\refg\x1b[201~h\r",
		line:           "efgh",
		throwAwayLines: 1,
	***REMOVED***,
	***REMOVED***
		// Lines consisting entirely of pasted data should be indicated as such.
		in:   "\x1b[200~a\r",
		line: "a",
		err:  ErrPasteIndicator,
	***REMOVED***,
***REMOVED***

func TestKeyPresses(t *testing.T) ***REMOVED***
	for i, test := range keyPressTests ***REMOVED***
		for j := 1; j < len(test.in); j++ ***REMOVED***
			c := &MockTerminal***REMOVED***
				toSend:       []byte(test.in),
				bytesPerRead: j,
			***REMOVED***
			ss := NewTerminal(c, "> ")
			for k := 0; k < test.throwAwayLines; k++ ***REMOVED***
				_, err := ss.ReadLine()
				if err != nil ***REMOVED***
					t.Errorf("Throwaway line %d from test %d resulted in error: %s", k, i, err)
				***REMOVED***
			***REMOVED***
			line, err := ss.ReadLine()
			if line != test.line ***REMOVED***
				t.Errorf("Line resulting from test %d (%d bytes per read) was '%s', expected '%s'", i, j, line, test.line)
				break
			***REMOVED***
			if err != test.err ***REMOVED***
				t.Errorf("Error resulting from test %d (%d bytes per read) was '%v', expected '%v'", i, j, err, test.err)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPasswordNotSaved(t *testing.T) ***REMOVED***
	c := &MockTerminal***REMOVED***
		toSend:       []byte("password\r\x1b[A\r"),
		bytesPerRead: 1,
	***REMOVED***
	ss := NewTerminal(c, "> ")
	pw, _ := ss.ReadPassword("> ")
	if pw != "password" ***REMOVED***
		t.Fatalf("failed to read password, got %s", pw)
	***REMOVED***
	line, _ := ss.ReadLine()
	if len(line) > 0 ***REMOVED***
		t.Fatalf("password was saved in history")
	***REMOVED***
***REMOVED***

var setSizeTests = []struct ***REMOVED***
	width, height int
***REMOVED******REMOVED***
	***REMOVED***40, 13***REMOVED***,
	***REMOVED***80, 24***REMOVED***,
	***REMOVED***132, 43***REMOVED***,
***REMOVED***

func TestTerminalSetSize(t *testing.T) ***REMOVED***
	for _, setSize := range setSizeTests ***REMOVED***
		c := &MockTerminal***REMOVED***
			toSend:       []byte("password\r\x1b[A\r"),
			bytesPerRead: 1,
		***REMOVED***
		ss := NewTerminal(c, "> ")
		ss.SetSize(setSize.width, setSize.height)
		pw, _ := ss.ReadPassword("Password: ")
		if pw != "password" ***REMOVED***
			t.Fatalf("failed to read password, got %s", pw)
		***REMOVED***
		if string(c.received) != "Password: \r\n" ***REMOVED***
			t.Errorf("failed to set the temporary prompt expected %q, got %q", "Password: ", c.received)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReadPasswordLineEnd(t *testing.T) ***REMOVED***
	var tests = []struct ***REMOVED***
		input string
		want  string
	***REMOVED******REMOVED***
		***REMOVED***"\n", ""***REMOVED***,
		***REMOVED***"\r\n", ""***REMOVED***,
		***REMOVED***"test\r\n", "test"***REMOVED***,
		***REMOVED***"testtesttesttes\n", "testtesttesttes"***REMOVED***,
		***REMOVED***"testtesttesttes\r\n", "testtesttesttes"***REMOVED***,
		***REMOVED***"testtesttesttesttest\n", "testtesttesttesttest"***REMOVED***,
		***REMOVED***"testtesttesttesttest\r\n", "testtesttesttesttest"***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		buf := new(bytes.Buffer)
		if _, err := buf.WriteString(test.input); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		have, err := readPasswordLine(buf)
		if err != nil ***REMOVED***
			t.Errorf("readPasswordLine(%q) failed: %v", test.input, err)
			continue
		***REMOVED***
		if string(have) != test.want ***REMOVED***
			t.Errorf("readPasswordLine(%q) returns %q, but %q is expected", test.input, string(have), test.want)
			continue
		***REMOVED***

		if _, err = buf.WriteString(test.input); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		have, err = readPasswordLine(buf)
		if err != nil ***REMOVED***
			t.Errorf("readPasswordLine(%q) failed: %v", test.input, err)
			continue
		***REMOVED***
		if string(have) != test.want ***REMOVED***
			t.Errorf("readPasswordLine(%q) returns %q, but %q is expected", test.input, string(have), test.want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMakeRawState(t *testing.T) ***REMOVED***
	fd := int(os.Stdout.Fd())
	if !IsTerminal(fd) ***REMOVED***
		t.Skip("stdout is not a terminal; skipping test")
	***REMOVED***

	st, err := GetState(fd)
	if err != nil ***REMOVED***
		t.Fatalf("failed to get terminal state from GetState: %s", err)
	***REMOVED***
	defer Restore(fd, st)
	raw, err := MakeRaw(fd)
	if err != nil ***REMOVED***
		t.Fatalf("failed to get terminal state from MakeRaw: %s", err)
	***REMOVED***

	if *st != *raw ***REMOVED***
		t.Errorf("states do not match; was %v, expected %v", raw, st)
	***REMOVED***
***REMOVED***

func TestOutputNewlines(t *testing.T) ***REMOVED***
	// \n should be changed to \r\n in terminal output.
	buf := new(bytes.Buffer)
	term := NewTerminal(buf, ">")

	term.Write([]byte("1\n2\n"))
	output := string(buf.Bytes())
	const expected = "1\r\n2\r\n"

	if output != expected ***REMOVED***
		t.Errorf("incorrect output: was %q, expected %q", output, expected)
	***REMOVED***
***REMOVED***
