// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package catmsg

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/text/language"
)

type renderer struct ***REMOVED***
	args   []int
	result string
***REMOVED***

func (r *renderer) Arg(i int) interface***REMOVED******REMOVED*** ***REMOVED***
	if i >= len(r.args) ***REMOVED***
		return nil
	***REMOVED***
	return r.args[i]
***REMOVED***

func (r *renderer) Render(s string) ***REMOVED***
	if r.result != "" ***REMOVED***
		r.result += "|"
	***REMOVED***
	r.result += s
***REMOVED***

func TestCodec(t *testing.T) ***REMOVED***
	type test struct ***REMOVED***
		args   []int
		out    string
		decErr string
	***REMOVED***
	single := func(out, err string) []test ***REMOVED*** return []test***REMOVED******REMOVED***out: out, decErr: err***REMOVED******REMOVED*** ***REMOVED***
	testCases := []struct ***REMOVED***
		desc   string
		m      Message
		enc    string
		encErr string
		tests  []test
	***REMOVED******REMOVED******REMOVED***
		desc:   "unused variable",
		m:      &Var***REMOVED***"name", String("foo")***REMOVED***,
		encErr: errIsVar.Error(),
		tests:  single("", ""),
	***REMOVED***, ***REMOVED***
		desc:  "empty",
		m:     empty***REMOVED******REMOVED***,
		tests: single("", ""),
	***REMOVED***, ***REMOVED***
		desc:  "sequence with empty",
		m:     seq***REMOVED***empty***REMOVED******REMOVED******REMOVED***,
		tests: single("", ""),
	***REMOVED***, ***REMOVED***
		desc:  "raw string",
		m:     Raw("foo"),
		tests: single("foo", ""),
	***REMOVED***, ***REMOVED***
		desc:  "raw string no sub",
		m:     Raw("$***REMOVED***foo***REMOVED***"),
		enc:   "\x02$***REMOVED***foo***REMOVED***",
		tests: single("$***REMOVED***foo***REMOVED***", ""),
	***REMOVED***, ***REMOVED***
		desc:  "simple string",
		m:     String("foo"),
		tests: single("foo", ""),
	***REMOVED***, ***REMOVED***
		desc:  "affix",
		m:     &Affix***REMOVED***String("foo"), "\t", "\n"***REMOVED***,
		tests: single("\t|foo|\n", ""),
	***REMOVED***, ***REMOVED***
		desc:   "missing var",
		m:      String("foo$***REMOVED***bar***REMOVED***"),
		enc:    "\x03\x03foo\x02\x03bar",
		encErr: `unknown var "bar"`,
		tests:  single("foo|bar", ""),
	***REMOVED***, ***REMOVED***
		desc: "empty var",
		m: seq***REMOVED***
			&Var***REMOVED***"bar", seq***REMOVED******REMOVED******REMOVED***,
			String("foo$***REMOVED***bar***REMOVED***"),
		***REMOVED***,
		enc: "\x00\x05\x04\x02bar\x03\x03foo\x00\x00",
		// TODO: recognize that it is cheaper to substitute bar.
		tests: single("foo|bar", ""),
	***REMOVED***, ***REMOVED***
		desc: "var after value",
		m: seq***REMOVED***
			String("foo$***REMOVED***bar***REMOVED***"),
			&Var***REMOVED***"bar", String("baz")***REMOVED***,
		***REMOVED***,
		encErr: errIsVar.Error(),
		tests:  single("foo|bar", ""),
	***REMOVED***, ***REMOVED***
		desc: "substitution",
		m: seq***REMOVED***
			&Var***REMOVED***"bar", String("baz")***REMOVED***,
			String("foo$***REMOVED***bar***REMOVED***"),
		***REMOVED***,
		tests: single("foo|baz", ""),
	***REMOVED***, ***REMOVED***
		desc: "affix with substitution",
		m: &Affix***REMOVED***seq***REMOVED***
			&Var***REMOVED***"bar", String("baz")***REMOVED***,
			String("foo$***REMOVED***bar***REMOVED***"),
		***REMOVED***, "\t", "\n"***REMOVED***,
		tests: single("\t|foo|baz|\n", ""),
	***REMOVED***, ***REMOVED***
		desc: "shadowed variable",
		m: seq***REMOVED***
			&Var***REMOVED***"bar", String("baz")***REMOVED***,
			seq***REMOVED***
				&Var***REMOVED***"bar", String("BAZ")***REMOVED***,
				String("foo$***REMOVED***bar***REMOVED***"),
			***REMOVED***,
		***REMOVED***,
		tests: single("foo|BAZ", ""),
	***REMOVED***, ***REMOVED***
		desc:  "nested value",
		m:     nestedLang***REMOVED***nestedLang***REMOVED***empty***REMOVED******REMOVED******REMOVED******REMOVED***,
		tests: single("nl|nl", ""),
	***REMOVED***, ***REMOVED***
		desc: "not shadowed variable",
		m: seq***REMOVED***
			&Var***REMOVED***"bar", String("baz")***REMOVED***,
			seq***REMOVED***
				String("foo$***REMOVED***bar***REMOVED***"),
				&Var***REMOVED***"bar", String("BAZ")***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		encErr: errIsVar.Error(),
		tests:  single("foo|baz", ""),
	***REMOVED***, ***REMOVED***
		desc: "duplicate variable",
		m: seq***REMOVED***
			&Var***REMOVED***"bar", String("baz")***REMOVED***,
			&Var***REMOVED***"bar", String("BAZ")***REMOVED***,
			String("$***REMOVED***bar***REMOVED***"),
		***REMOVED***,
		encErr: "catmsg: duplicate variable \"bar\"",
		tests:  single("baz", ""),
	***REMOVED***, ***REMOVED***
		desc: "complete incomplete variable",
		m: seq***REMOVED***
			&Var***REMOVED***"bar", incomplete***REMOVED******REMOVED******REMOVED***,
			String("$***REMOVED***bar***REMOVED***"),
		***REMOVED***,
		enc: "\x00\t\b\x01\x01\x14\x04\x02bar\x03\x00\x00\x00",
		// TODO: recognize that it is cheaper to substitute bar.
		tests: single("bar", ""),
	***REMOVED***, ***REMOVED***
		desc: "incomplete sequence",
		m: seq***REMOVED***
			incomplete***REMOVED******REMOVED***,
			incomplete***REMOVED******REMOVED***,
		***REMOVED***,
		encErr: ErrIncomplete.Error(),
		tests:  single("", ErrNoMatch.Error()),
	***REMOVED***, ***REMOVED***
		desc: "compile error variable",
		m: seq***REMOVED***
			&Var***REMOVED***"bar", errorCompileMsg***REMOVED******REMOVED******REMOVED***,
			String("$***REMOVED***bar***REMOVED***"),
		***REMOVED***,
		encErr: errCompileTest.Error(),
		tests:  single("bar", ""),
	***REMOVED***, ***REMOVED***
		desc:   "compile error message",
		m:      errorCompileMsg***REMOVED******REMOVED***,
		encErr: errCompileTest.Error(),
		tests:  single("", ""),
	***REMOVED***, ***REMOVED***
		desc: "compile error sequence",
		m: seq***REMOVED***
			errorCompileMsg***REMOVED******REMOVED***,
			errorCompileMsg***REMOVED******REMOVED***,
		***REMOVED***,
		encErr: errCompileTest.Error(),
		tests:  single("", ""),
	***REMOVED***, ***REMOVED***
		desc:  "macro",
		m:     String("$***REMOVED***exists(1)***REMOVED***"),
		tests: single("you betya!", ""),
	***REMOVED***, ***REMOVED***
		desc:  "macro incomplete",
		m:     String("$***REMOVED***incomplete(1)***REMOVED***"),
		enc:   "\x03\x00\x01\nincomplete\x01",
		tests: single("incomplete", ""),
	***REMOVED***, ***REMOVED***
		desc:  "macro undefined at end",
		m:     String("$***REMOVED***undefined(1)***REMOVED***"),
		enc:   "\x03\x00\x01\tundefined\x01",
		tests: single("undefined", "catmsg: undefined macro \"undefined\""),
	***REMOVED***, ***REMOVED***
		desc:  "macro undefined with more text following",
		m:     String("$***REMOVED***undefined(1)***REMOVED***."),
		enc:   "\x03\x00\x01\tundefined\x01\x01.",
		tests: single("undefined|.", "catmsg: undefined macro \"undefined\""),
	***REMOVED***, ***REMOVED***
		desc:   "macro missing paren",
		m:      String("$***REMOVED***missing(1***REMOVED***"),
		encErr: "catmsg: missing ')'",
		tests:  single("$!(MISSINGPAREN)", ""),
	***REMOVED***, ***REMOVED***
		desc:   "macro bad num",
		m:      String("aa$***REMOVED***bad(a)***REMOVED***"),
		encErr: "catmsg: invalid number \"a\"",
		tests:  single("aa$!(BADNUM)", ""),
	***REMOVED***, ***REMOVED***
		desc:   "var missing brace",
		m:      String("a$***REMOVED***missing"),
		encErr: "catmsg: missing '***REMOVED***'",
		tests:  single("a$!(MISSINGBRACE)", ""),
	***REMOVED******REMOVED***
	r := &renderer***REMOVED******REMOVED***
	dec := NewDecoder(language.Und, r, macros)
	for _, tc := range testCases ***REMOVED***
		t.Run(tc.desc, func(t *testing.T) ***REMOVED***
			// Use a language other than Und so that we can test
			// passing the language to nested values.
			data, err := Compile(language.Dutch, macros, tc.m)
			if failErr(err, tc.encErr) ***REMOVED***
				t.Errorf("encoding error: got %+q; want %+q", err, tc.encErr)
			***REMOVED***
			if tc.enc != "" && data != tc.enc ***REMOVED***
				t.Errorf("encoding: got %+q; want %+q", data, tc.enc)
			***REMOVED***
			for _, st := range tc.tests ***REMOVED***
				t.Run("", func(t *testing.T) ***REMOVED***
					*r = renderer***REMOVED***args: st.args***REMOVED***
					if err = dec.Execute(data); failErr(err, st.decErr) ***REMOVED***
						t.Errorf("decoding error: got %+q; want %+q", err, st.decErr)
					***REMOVED***
					if r.result != st.out ***REMOVED***
						t.Errorf("decode: got %+q; want %+q", r.result, st.out)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func failErr(got error, want string) bool ***REMOVED***
	if got == nil ***REMOVED***
		return want != ""
	***REMOVED***
	return want == "" || !strings.Contains(got.Error(), want)
***REMOVED***

type seq []Message

func (s seq) Compile(e *Encoder) (err error) ***REMOVED***
	err = ErrIncomplete
	e.EncodeMessageType(msgFirst)
	for _, m := range s ***REMOVED***
		// Pass only the last error, but allow erroneous or complete messages
		// here to allow testing different scenarios.
		err = e.EncodeMessage(m)
	***REMOVED***
	return err
***REMOVED***

type empty struct***REMOVED******REMOVED***

func (empty) Compile(e *Encoder) (err error) ***REMOVED*** return nil ***REMOVED***

var msgIncomplete = Register(
	"golang.org/x/text/internal/catmsg.incomplete",
	func(d *Decoder) bool ***REMOVED*** return false ***REMOVED***)

type incomplete struct***REMOVED******REMOVED***

func (incomplete) Compile(e *Encoder) (err error) ***REMOVED***
	e.EncodeMessageType(msgIncomplete)
	return ErrIncomplete
***REMOVED***

var msgNested = Register(
	"golang.org/x/text/internal/catmsg.nested",
	func(d *Decoder) bool ***REMOVED***
		d.Render(d.DecodeString())
		d.ExecuteMessage()
		return true
	***REMOVED***)

type nestedLang struct***REMOVED*** Message ***REMOVED***

func (n nestedLang) Compile(e *Encoder) (err error) ***REMOVED***
	e.EncodeMessageType(msgNested)
	e.EncodeString(e.Language().String())
	e.EncodeMessage(n.Message)
	return nil
***REMOVED***

type errorCompileMsg struct***REMOVED******REMOVED***

var errCompileTest = errors.New("catmsg: compile error test")

func (errorCompileMsg) Compile(e *Encoder) (err error) ***REMOVED***
	return errCompileTest
***REMOVED***

type dictionary struct***REMOVED******REMOVED***

var (
	macros       = dictionary***REMOVED******REMOVED***
	dictMessages = map[string]string***REMOVED***
		"exists":     compile(String("you betya!")),
		"incomplete": compile(incomplete***REMOVED******REMOVED***),
	***REMOVED***
)

func (d dictionary) Lookup(key string) (data string, ok bool) ***REMOVED***
	data, ok = dictMessages[key]
	return
***REMOVED***

func compile(m Message) (data string) ***REMOVED***
	data, _ = Compile(language.Und, macros, m)
	return data
***REMOVED***
