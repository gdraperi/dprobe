// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dnsmessage

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func mustNewName(name string) Name ***REMOVED***
	n, err := NewName(name)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return n
***REMOVED***

func (m *Message) String() string ***REMOVED***
	s := fmt.Sprintf("Message: %#v\n", &m.Header)
	if len(m.Questions) > 0 ***REMOVED***
		s += "-- Questions\n"
		for _, q := range m.Questions ***REMOVED***
			s += fmt.Sprintf("%#v\n", q)
		***REMOVED***
	***REMOVED***
	if len(m.Answers) > 0 ***REMOVED***
		s += "-- Answers\n"
		for _, a := range m.Answers ***REMOVED***
			s += fmt.Sprintf("%#v\n", a)
		***REMOVED***
	***REMOVED***
	if len(m.Authorities) > 0 ***REMOVED***
		s += "-- Authorities\n"
		for _, ns := range m.Authorities ***REMOVED***
			s += fmt.Sprintf("%#v\n", ns)
		***REMOVED***
	***REMOVED***
	if len(m.Additionals) > 0 ***REMOVED***
		s += "-- Additionals\n"
		for _, e := range m.Additionals ***REMOVED***
			s += fmt.Sprintf("%#v\n", e)
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

func TestNameString(t *testing.T) ***REMOVED***
	want := "foo"
	name := mustNewName(want)
	if got := fmt.Sprint(name); got != want ***REMOVED***
		t.Errorf("got fmt.Sprint(%#v) = %s, want = %s", name, got, want)
	***REMOVED***
***REMOVED***

func TestQuestionPackUnpack(t *testing.T) ***REMOVED***
	want := Question***REMOVED***
		Name:  mustNewName("."),
		Type:  TypeA,
		Class: ClassINET,
	***REMOVED***
	buf, err := want.pack(make([]byte, 1, 50), map[string]int***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("Packing failed:", err)
	***REMOVED***
	var p Parser
	p.msg = buf
	p.header.questions = 1
	p.section = sectionQuestions
	p.off = 1
	got, err := p.Question()
	if err != nil ***REMOVED***
		t.Fatalf("Unpacking failed: %v\n%s", err, string(buf[1:]))
	***REMOVED***
	if p.off != len(buf) ***REMOVED***
		t.Errorf("Unpacked different amount than packed: got n = %d, want = %d", p.off, len(buf))
	***REMOVED***
	if !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("Got = %+v, want = %+v", got, want)
	***REMOVED***
***REMOVED***

func TestName(t *testing.T) ***REMOVED***
	tests := []string***REMOVED***
		"",
		".",
		"google..com",
		"google.com",
		"google..com.",
		"google.com.",
		".google.com.",
		"www..google.com.",
		"www.google.com.",
	***REMOVED***

	for _, test := range tests ***REMOVED***
		n, err := NewName(test)
		if err != nil ***REMOVED***
			t.Errorf("Creating name for %q: %v", test, err)
			continue
		***REMOVED***
		if ns := n.String(); ns != test ***REMOVED***
			t.Errorf("Got %#v.String() = %q, want = %q", n, ns, test)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNamePackUnpack(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in   string
		want string
		err  error
	***REMOVED******REMOVED***
		***REMOVED***"", "", errNonCanonicalName***REMOVED***,
		***REMOVED***".", ".", nil***REMOVED***,
		***REMOVED***"google..com", "", errNonCanonicalName***REMOVED***,
		***REMOVED***"google.com", "", errNonCanonicalName***REMOVED***,
		***REMOVED***"google..com.", "", errZeroSegLen***REMOVED***,
		***REMOVED***"google.com.", "google.com.", nil***REMOVED***,
		***REMOVED***".google.com.", "", errZeroSegLen***REMOVED***,
		***REMOVED***"www..google.com.", "", errZeroSegLen***REMOVED***,
		***REMOVED***"www.google.com.", "www.google.com.", nil***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		in := mustNewName(test.in)
		want := mustNewName(test.want)
		buf, err := in.pack(make([]byte, 0, 30), map[string]int***REMOVED******REMOVED***)
		if err != test.err ***REMOVED***
			t.Errorf("Packing of %q: got err = %v, want err = %v", test.in, err, test.err)
			continue
		***REMOVED***
		if test.err != nil ***REMOVED***
			continue
		***REMOVED***
		var got Name
		n, err := got.unpack(buf, 0)
		if err != nil ***REMOVED***
			t.Errorf("Unpacking for %q failed: %v", test.in, err)
			continue
		***REMOVED***
		if n != len(buf) ***REMOVED***
			t.Errorf(
				"Unpacked different amount than packed for %q: got n = %d, want = %d",
				test.in,
				n,
				len(buf),
			)
		***REMOVED***
		if got != want ***REMOVED***
			t.Errorf("Unpacking packing of %q: got = %#v, want = %#v", test.in, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func checkErrorPrefix(err error, prefix string) bool ***REMOVED***
	e, ok := err.(*nestedError)
	return ok && e.s == prefix
***REMOVED***

func TestHeaderUnpackError(t *testing.T) ***REMOVED***
	wants := []string***REMOVED***
		"id",
		"bits",
		"questions",
		"answers",
		"authorities",
		"additionals",
	***REMOVED***
	var buf []byte
	var h header
	for _, want := range wants ***REMOVED***
		n, err := h.unpack(buf, 0)
		if n != 0 || !checkErrorPrefix(err, want) ***REMOVED***
			t.Errorf("got h.unpack([%d]byte, 0) = %d, %v, want = 0, %s", len(buf), n, err, want)
		***REMOVED***
		buf = append(buf, 0, 0)
	***REMOVED***
***REMOVED***

func TestParserStart(t *testing.T) ***REMOVED***
	const want = "unpacking header"
	var p Parser
	for i := 0; i <= 1; i++ ***REMOVED***
		_, err := p.Start([]byte***REMOVED******REMOVED***)
		if !checkErrorPrefix(err, want) ***REMOVED***
			t.Errorf("got p.Start(nil) = _, %v, want = _, %s", err, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestResourceNotStarted(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		name string
		fn   func(*Parser) error
	***REMOVED******REMOVED***
		***REMOVED***"CNAMEResource", func(p *Parser) error ***REMOVED*** _, err := p.CNAMEResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"MXResource", func(p *Parser) error ***REMOVED*** _, err := p.MXResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"NSResource", func(p *Parser) error ***REMOVED*** _, err := p.NSResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"PTRResource", func(p *Parser) error ***REMOVED*** _, err := p.PTRResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"SOAResource", func(p *Parser) error ***REMOVED*** _, err := p.SOAResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"TXTResource", func(p *Parser) error ***REMOVED*** _, err := p.TXTResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"SRVResource", func(p *Parser) error ***REMOVED*** _, err := p.SRVResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"AResource", func(p *Parser) error ***REMOVED*** _, err := p.AResource(); return err ***REMOVED******REMOVED***,
		***REMOVED***"AAAAResource", func(p *Parser) error ***REMOVED*** _, err := p.AAAAResource(); return err ***REMOVED******REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		if err := test.fn(&Parser***REMOVED******REMOVED***); err != ErrNotStarted ***REMOVED***
			t.Errorf("got _, %v = p.%s(), want = _, %v", err, test.name, ErrNotStarted)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDNSPackUnpack(t *testing.T) ***REMOVED***
	wants := []Message***REMOVED***
		***REMOVED***
			Questions: []Question***REMOVED***
				***REMOVED***
					Name:  mustNewName("."),
					Type:  TypeAAAA,
					Class: ClassINET,
				***REMOVED***,
			***REMOVED***,
			Answers:     []Resource***REMOVED******REMOVED***,
			Authorities: []Resource***REMOVED******REMOVED***,
			Additionals: []Resource***REMOVED******REMOVED***,
		***REMOVED***,
		largeTestMsg(),
	***REMOVED***
	for i, want := range wants ***REMOVED***
		b, err := want.Pack()
		if err != nil ***REMOVED***
			t.Fatalf("%d: packing failed: %v", i, err)
		***REMOVED***
		var got Message
		err = got.Unpack(b)
		if err != nil ***REMOVED***
			t.Fatalf("%d: unpacking failed: %v", i, err)
		***REMOVED***
		if !reflect.DeepEqual(got, want) ***REMOVED***
			t.Errorf("%d: got = %+v, want = %+v", i, &got, &want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSkipAll(t *testing.T) ***REMOVED***
	msg := largeTestMsg()
	buf, err := msg.Pack()
	if err != nil ***REMOVED***
		t.Fatal("Packing large test message:", err)
	***REMOVED***
	var p Parser
	if _, err := p.Start(buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tests := []struct ***REMOVED***
		name string
		f    func() error
	***REMOVED******REMOVED***
		***REMOVED***"SkipAllQuestions", p.SkipAllQuestions***REMOVED***,
		***REMOVED***"SkipAllAnswers", p.SkipAllAnswers***REMOVED***,
		***REMOVED***"SkipAllAuthorities", p.SkipAllAuthorities***REMOVED***,
		***REMOVED***"SkipAllAdditionals", p.SkipAllAdditionals***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		for i := 1; i <= 3; i++ ***REMOVED***
			if err := test.f(); err != nil ***REMOVED***
				t.Errorf("Call #%d to %s(): %v", i, test.name, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSkipEach(t *testing.T) ***REMOVED***
	msg := smallTestMsg()

	buf, err := msg.Pack()
	if err != nil ***REMOVED***
		t.Fatal("Packing test message:", err)
	***REMOVED***
	var p Parser
	if _, err := p.Start(buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tests := []struct ***REMOVED***
		name string
		f    func() error
	***REMOVED******REMOVED***
		***REMOVED***"SkipQuestion", p.SkipQuestion***REMOVED***,
		***REMOVED***"SkipAnswer", p.SkipAnswer***REMOVED***,
		***REMOVED***"SkipAuthority", p.SkipAuthority***REMOVED***,
		***REMOVED***"SkipAdditional", p.SkipAdditional***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		if err := test.f(); err != nil ***REMOVED***
			t.Errorf("First call: got %s() = %v, want = %v", test.name, err, nil)
		***REMOVED***
		if err := test.f(); err != ErrSectionDone ***REMOVED***
			t.Errorf("Second call: got %s() = %v, want = %v", test.name, err, ErrSectionDone)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSkipAfterRead(t *testing.T) ***REMOVED***
	msg := smallTestMsg()

	buf, err := msg.Pack()
	if err != nil ***REMOVED***
		t.Fatal("Packing test message:", err)
	***REMOVED***
	var p Parser
	if _, err := p.Start(buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tests := []struct ***REMOVED***
		name string
		skip func() error
		read func() error
	***REMOVED******REMOVED***
		***REMOVED***"Question", p.SkipQuestion, func() error ***REMOVED*** _, err := p.Question(); return err ***REMOVED******REMOVED***,
		***REMOVED***"Answer", p.SkipAnswer, func() error ***REMOVED*** _, err := p.Answer(); return err ***REMOVED******REMOVED***,
		***REMOVED***"Authority", p.SkipAuthority, func() error ***REMOVED*** _, err := p.Authority(); return err ***REMOVED******REMOVED***,
		***REMOVED***"Additional", p.SkipAdditional, func() error ***REMOVED*** _, err := p.Additional(); return err ***REMOVED******REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		if err := test.read(); err != nil ***REMOVED***
			t.Errorf("Got %s() = _, %v, want = _, %v", test.name, err, nil)
		***REMOVED***
		if err := test.skip(); err != ErrSectionDone ***REMOVED***
			t.Errorf("Got Skip%s() = %v, want = %v", test.name, err, ErrSectionDone)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSkipNotStarted(t *testing.T) ***REMOVED***
	var p Parser

	tests := []struct ***REMOVED***
		name string
		f    func() error
	***REMOVED******REMOVED***
		***REMOVED***"SkipAllQuestions", p.SkipAllQuestions***REMOVED***,
		***REMOVED***"SkipAllAnswers", p.SkipAllAnswers***REMOVED***,
		***REMOVED***"SkipAllAuthorities", p.SkipAllAuthorities***REMOVED***,
		***REMOVED***"SkipAllAdditionals", p.SkipAllAdditionals***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		if err := test.f(); err != ErrNotStarted ***REMOVED***
			t.Errorf("Got %s() = %v, want = %v", test.name, err, ErrNotStarted)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTooManyRecords(t *testing.T) ***REMOVED***
	const recs = int(^uint16(0)) + 1
	tests := []struct ***REMOVED***
		name string
		msg  Message
		want error
	***REMOVED******REMOVED***
		***REMOVED***
			"Questions",
			Message***REMOVED***
				Questions: make([]Question, recs),
			***REMOVED***,
			errTooManyQuestions,
		***REMOVED***,
		***REMOVED***
			"Answers",
			Message***REMOVED***
				Answers: make([]Resource, recs),
			***REMOVED***,
			errTooManyAnswers,
		***REMOVED***,
		***REMOVED***
			"Authorities",
			Message***REMOVED***
				Authorities: make([]Resource, recs),
			***REMOVED***,
			errTooManyAuthorities,
		***REMOVED***,
		***REMOVED***
			"Additionals",
			Message***REMOVED***
				Additionals: make([]Resource, recs),
			***REMOVED***,
			errTooManyAdditionals,
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		if _, got := test.msg.Pack(); got != test.want ***REMOVED***
			t.Errorf("Packing %d %s: got = %v, want = %v", recs, test.name, got, test.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestVeryLongTxt(t *testing.T) ***REMOVED***
	want := Resource***REMOVED***
		ResourceHeader***REMOVED***
			Name:  mustNewName("foo.bar.example.com."),
			Type:  TypeTXT,
			Class: ClassINET,
		***REMOVED***,
		&TXTResource***REMOVED***loremIpsum***REMOVED***,
	***REMOVED***
	buf, err := want.pack(make([]byte, 0, 8000), map[string]int***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("Packing failed:", err)
	***REMOVED***
	var got Resource
	off, err := got.Header.unpack(buf, 0)
	if err != nil ***REMOVED***
		t.Fatal("Unpacking ResourceHeader failed:", err)
	***REMOVED***
	body, n, err := unpackResourceBody(buf, off, got.Header)
	if err != nil ***REMOVED***
		t.Fatal("Unpacking failed:", err)
	***REMOVED***
	got.Body = body
	if n != len(buf) ***REMOVED***
		t.Errorf("Unpacked different amount than packed: got n = %d, want = %d", n, len(buf))
	***REMOVED***
	if !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("Got = %#v, want = %#v", got, want)
	***REMOVED***
***REMOVED***

func TestStartError(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		name string
		fn   func(*Builder) error
	***REMOVED******REMOVED***
		***REMOVED***"Questions", func(b *Builder) error ***REMOVED*** return b.StartQuestions() ***REMOVED******REMOVED***,
		***REMOVED***"Answers", func(b *Builder) error ***REMOVED*** return b.StartAnswers() ***REMOVED******REMOVED***,
		***REMOVED***"Authorities", func(b *Builder) error ***REMOVED*** return b.StartAuthorities() ***REMOVED******REMOVED***,
		***REMOVED***"Additionals", func(b *Builder) error ***REMOVED*** return b.StartAdditionals() ***REMOVED******REMOVED***,
	***REMOVED***

	envs := []struct ***REMOVED***
		name string
		fn   func() *Builder
		want error
	***REMOVED******REMOVED***
		***REMOVED***"sectionNotStarted", func() *Builder ***REMOVED*** return &Builder***REMOVED***section: sectionNotStarted***REMOVED*** ***REMOVED***, ErrNotStarted***REMOVED***,
		***REMOVED***"sectionDone", func() *Builder ***REMOVED*** return &Builder***REMOVED***section: sectionDone***REMOVED*** ***REMOVED***, ErrSectionDone***REMOVED***,
	***REMOVED***

	for _, env := range envs ***REMOVED***
		for _, test := range tests ***REMOVED***
			if got := test.fn(env.fn()); got != env.want ***REMOVED***
				t.Errorf("got Builder***REMOVED***%s***REMOVED***.Start%s = %v, want = %v", env.name, test.name, got, env.want)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBuilderResourceError(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		name string
		fn   func(*Builder) error
	***REMOVED******REMOVED***
		***REMOVED***"CNAMEResource", func(b *Builder) error ***REMOVED*** return b.CNAMEResource(ResourceHeader***REMOVED******REMOVED***, CNAMEResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"MXResource", func(b *Builder) error ***REMOVED*** return b.MXResource(ResourceHeader***REMOVED******REMOVED***, MXResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"NSResource", func(b *Builder) error ***REMOVED*** return b.NSResource(ResourceHeader***REMOVED******REMOVED***, NSResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"PTRResource", func(b *Builder) error ***REMOVED*** return b.PTRResource(ResourceHeader***REMOVED******REMOVED***, PTRResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"SOAResource", func(b *Builder) error ***REMOVED*** return b.SOAResource(ResourceHeader***REMOVED******REMOVED***, SOAResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"TXTResource", func(b *Builder) error ***REMOVED*** return b.TXTResource(ResourceHeader***REMOVED******REMOVED***, TXTResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"SRVResource", func(b *Builder) error ***REMOVED*** return b.SRVResource(ResourceHeader***REMOVED******REMOVED***, SRVResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"AResource", func(b *Builder) error ***REMOVED*** return b.AResource(ResourceHeader***REMOVED******REMOVED***, AResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
		***REMOVED***"AAAAResource", func(b *Builder) error ***REMOVED*** return b.AAAAResource(ResourceHeader***REMOVED******REMOVED***, AAAAResource***REMOVED******REMOVED***) ***REMOVED******REMOVED***,
	***REMOVED***

	envs := []struct ***REMOVED***
		name string
		fn   func() *Builder
		want error
	***REMOVED******REMOVED***
		***REMOVED***"sectionNotStarted", func() *Builder ***REMOVED*** return &Builder***REMOVED***section: sectionNotStarted***REMOVED*** ***REMOVED***, ErrNotStarted***REMOVED***,
		***REMOVED***"sectionHeader", func() *Builder ***REMOVED*** return &Builder***REMOVED***section: sectionHeader***REMOVED*** ***REMOVED***, ErrNotStarted***REMOVED***,
		***REMOVED***"sectionQuestions", func() *Builder ***REMOVED*** return &Builder***REMOVED***section: sectionQuestions***REMOVED*** ***REMOVED***, ErrNotStarted***REMOVED***,
		***REMOVED***"sectionDone", func() *Builder ***REMOVED*** return &Builder***REMOVED***section: sectionDone***REMOVED*** ***REMOVED***, ErrSectionDone***REMOVED***,
	***REMOVED***

	for _, env := range envs ***REMOVED***
		for _, test := range tests ***REMOVED***
			if got := test.fn(env.fn()); got != env.want ***REMOVED***
				t.Errorf("got Builder***REMOVED***%s***REMOVED***.%s = %v, want = %v", env.name, test.name, got, env.want)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFinishError(t *testing.T) ***REMOVED***
	var b Builder
	want := ErrNotStarted
	if _, got := b.Finish(); got != want ***REMOVED***
		t.Errorf("got Builder***REMOVED******REMOVED***.Finish() = %v, want = %v", got, want)
	***REMOVED***
***REMOVED***

func TestBuilder(t *testing.T) ***REMOVED***
	msg := largeTestMsg()
	want, err := msg.Pack()
	if err != nil ***REMOVED***
		t.Fatal("Packing without builder:", err)
	***REMOVED***

	var b Builder
	b.Start(nil, msg.Header)

	if err := b.StartQuestions(); err != nil ***REMOVED***
		t.Fatal("b.StartQuestions():", err)
	***REMOVED***
	for _, q := range msg.Questions ***REMOVED***
		if err := b.Question(q); err != nil ***REMOVED***
			t.Fatalf("b.Question(%#v): %v", q, err)
		***REMOVED***
	***REMOVED***

	if err := b.StartAnswers(); err != nil ***REMOVED***
		t.Fatal("b.StartAnswers():", err)
	***REMOVED***
	for _, a := range msg.Answers ***REMOVED***
		switch a.Header.Type ***REMOVED***
		case TypeA:
			if err := b.AResource(a.Header, *a.Body.(*AResource)); err != nil ***REMOVED***
				t.Fatalf("b.AResource(%#v): %v", a, err)
			***REMOVED***
		case TypeNS:
			if err := b.NSResource(a.Header, *a.Body.(*NSResource)); err != nil ***REMOVED***
				t.Fatalf("b.NSResource(%#v): %v", a, err)
			***REMOVED***
		case TypeCNAME:
			if err := b.CNAMEResource(a.Header, *a.Body.(*CNAMEResource)); err != nil ***REMOVED***
				t.Fatalf("b.CNAMEResource(%#v): %v", a, err)
			***REMOVED***
		case TypeSOA:
			if err := b.SOAResource(a.Header, *a.Body.(*SOAResource)); err != nil ***REMOVED***
				t.Fatalf("b.SOAResource(%#v): %v", a, err)
			***REMOVED***
		case TypePTR:
			if err := b.PTRResource(a.Header, *a.Body.(*PTRResource)); err != nil ***REMOVED***
				t.Fatalf("b.PTRResource(%#v): %v", a, err)
			***REMOVED***
		case TypeMX:
			if err := b.MXResource(a.Header, *a.Body.(*MXResource)); err != nil ***REMOVED***
				t.Fatalf("b.MXResource(%#v): %v", a, err)
			***REMOVED***
		case TypeTXT:
			if err := b.TXTResource(a.Header, *a.Body.(*TXTResource)); err != nil ***REMOVED***
				t.Fatalf("b.TXTResource(%#v): %v", a, err)
			***REMOVED***
		case TypeAAAA:
			if err := b.AAAAResource(a.Header, *a.Body.(*AAAAResource)); err != nil ***REMOVED***
				t.Fatalf("b.AAAAResource(%#v): %v", a, err)
			***REMOVED***
		case TypeSRV:
			if err := b.SRVResource(a.Header, *a.Body.(*SRVResource)); err != nil ***REMOVED***
				t.Fatalf("b.SRVResource(%#v): %v", a, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := b.StartAuthorities(); err != nil ***REMOVED***
		t.Fatal("b.StartAuthorities():", err)
	***REMOVED***
	for _, a := range msg.Authorities ***REMOVED***
		if err := b.NSResource(a.Header, *a.Body.(*NSResource)); err != nil ***REMOVED***
			t.Fatalf("b.NSResource(%#v): %v", a, err)
		***REMOVED***
	***REMOVED***

	if err := b.StartAdditionals(); err != nil ***REMOVED***
		t.Fatal("b.StartAdditionals():", err)
	***REMOVED***
	for _, a := range msg.Additionals ***REMOVED***
		if err := b.TXTResource(a.Header, *a.Body.(*TXTResource)); err != nil ***REMOVED***
			t.Fatalf("b.TXTResource(%#v): %v", a, err)
		***REMOVED***
	***REMOVED***

	got, err := b.Finish()
	if err != nil ***REMOVED***
		t.Fatal("b.Finish():", err)
	***REMOVED***
	if !bytes.Equal(got, want) ***REMOVED***
		t.Fatalf("Got from Builder: %#v\nwant = %#v", got, want)
	***REMOVED***
***REMOVED***

func TestResourcePack(t *testing.T) ***REMOVED***
	for _, tt := range []struct ***REMOVED***
		m   Message
		err error
	***REMOVED******REMOVED***
		***REMOVED***
			Message***REMOVED***
				Questions: []Question***REMOVED***
					***REMOVED***
						Name:  mustNewName("."),
						Type:  TypeAAAA,
						Class: ClassINET,
					***REMOVED***,
				***REMOVED***,
				Answers: []Resource***REMOVED******REMOVED***ResourceHeader***REMOVED******REMOVED***, nil***REMOVED******REMOVED***,
			***REMOVED***,
			&nestedError***REMOVED***"packing Answer", errNilResouceBody***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Message***REMOVED***
				Questions: []Question***REMOVED***
					***REMOVED***
						Name:  mustNewName("."),
						Type:  TypeAAAA,
						Class: ClassINET,
					***REMOVED***,
				***REMOVED***,
				Authorities: []Resource***REMOVED******REMOVED***ResourceHeader***REMOVED******REMOVED***, (*NSResource)(nil)***REMOVED******REMOVED***,
			***REMOVED***,
			&nestedError***REMOVED***"packing Authority",
				&nestedError***REMOVED***"ResourceHeader",
					&nestedError***REMOVED***"Name", errNonCanonicalName***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Message***REMOVED***
				Questions: []Question***REMOVED***
					***REMOVED***
						Name:  mustNewName("."),
						Type:  TypeA,
						Class: ClassINET,
					***REMOVED***,
				***REMOVED***,
				Additionals: []Resource***REMOVED******REMOVED***ResourceHeader***REMOVED******REMOVED***, nil***REMOVED******REMOVED***,
			***REMOVED***,
			&nestedError***REMOVED***"packing Additional", errNilResouceBody***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		_, err := tt.m.Pack()
		if !reflect.DeepEqual(err, tt.err) ***REMOVED***
			t.Errorf("got %v for %v; want %v", err, tt.m, tt.err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkParsing(b *testing.B) ***REMOVED***
	b.ReportAllocs()

	name := mustNewName("foo.bar.example.com.")
	msg := Message***REMOVED***
		Header: Header***REMOVED***Response: true, Authoritative: true***REMOVED***,
		Questions: []Question***REMOVED***
			***REMOVED***
				Name:  name,
				Type:  TypeA,
				Class: ClassINET,
			***REMOVED***,
		***REMOVED***,
		Answers: []Resource***REMOVED***
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Class: ClassINET,
				***REMOVED***,
				&AResource***REMOVED***[4]byte***REMOVED******REMOVED******REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Class: ClassINET,
				***REMOVED***,
				&AAAAResource***REMOVED***[16]byte***REMOVED******REMOVED******REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Class: ClassINET,
				***REMOVED***,
				&CNAMEResource***REMOVED***name***REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Class: ClassINET,
				***REMOVED***,
				&NSResource***REMOVED***name***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	buf, err := msg.Pack()
	if err != nil ***REMOVED***
		b.Fatal("msg.Pack():", err)
	***REMOVED***

	for i := 0; i < b.N; i++ ***REMOVED***
		var p Parser
		if _, err := p.Start(buf); err != nil ***REMOVED***
			b.Fatal("p.Start(buf):", err)
		***REMOVED***

		for ***REMOVED***
			_, err := p.Question()
			if err == ErrSectionDone ***REMOVED***
				break
			***REMOVED***
			if err != nil ***REMOVED***
				b.Fatal("p.Question():", err)
			***REMOVED***
		***REMOVED***

		for ***REMOVED***
			h, err := p.AnswerHeader()
			if err == ErrSectionDone ***REMOVED***
				break
			***REMOVED***
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***

			switch h.Type ***REMOVED***
			case TypeA:
				if _, err := p.AResource(); err != nil ***REMOVED***
					b.Fatal("p.AResource():", err)
				***REMOVED***
			case TypeAAAA:
				if _, err := p.AAAAResource(); err != nil ***REMOVED***
					b.Fatal("p.AAAAResource():", err)
				***REMOVED***
			case TypeCNAME:
				if _, err := p.CNAMEResource(); err != nil ***REMOVED***
					b.Fatal("p.CNAMEResource():", err)
				***REMOVED***
			case TypeNS:
				if _, err := p.NSResource(); err != nil ***REMOVED***
					b.Fatal("p.NSResource():", err)
				***REMOVED***
			default:
				b.Fatalf("unknown type: %T", h)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkBuilding(b *testing.B) ***REMOVED***
	b.ReportAllocs()

	name := mustNewName("foo.bar.example.com.")
	buf := make([]byte, 0, packStartingCap)

	for i := 0; i < b.N; i++ ***REMOVED***
		var bld Builder
		bld.StartWithoutCompression(buf, Header***REMOVED***Response: true, Authoritative: true***REMOVED***)

		if err := bld.StartQuestions(); err != nil ***REMOVED***
			b.Fatal("bld.StartQuestions():", err)
		***REMOVED***
		q := Question***REMOVED***
			Name:  name,
			Type:  TypeA,
			Class: ClassINET,
		***REMOVED***
		if err := bld.Question(q); err != nil ***REMOVED***
			b.Fatalf("bld.Question(%+v): %v", q, err)
		***REMOVED***

		hdr := ResourceHeader***REMOVED***
			Name:  name,
			Class: ClassINET,
		***REMOVED***
		if err := bld.StartAnswers(); err != nil ***REMOVED***
			b.Fatal("bld.StartQuestions():", err)
		***REMOVED***

		ar := AResource***REMOVED***[4]byte***REMOVED******REMOVED******REMOVED***
		if err := bld.AResource(hdr, ar); err != nil ***REMOVED***
			b.Fatalf("bld.AResource(%+v, %+v): %v", hdr, ar, err)
		***REMOVED***

		aaar := AAAAResource***REMOVED***[16]byte***REMOVED******REMOVED******REMOVED***
		if err := bld.AAAAResource(hdr, aaar); err != nil ***REMOVED***
			b.Fatalf("bld.AAAAResource(%+v, %+v): %v", hdr, aaar, err)
		***REMOVED***

		cnr := CNAMEResource***REMOVED***name***REMOVED***
		if err := bld.CNAMEResource(hdr, cnr); err != nil ***REMOVED***
			b.Fatalf("bld.CNAMEResource(%+v, %+v): %v", hdr, cnr, err)
		***REMOVED***

		nsr := NSResource***REMOVED***name***REMOVED***
		if err := bld.NSResource(hdr, nsr); err != nil ***REMOVED***
			b.Fatalf("bld.NSResource(%+v, %+v): %v", hdr, nsr, err)
		***REMOVED***

		if _, err := bld.Finish(); err != nil ***REMOVED***
			b.Fatal("bld.Finish():", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func smallTestMsg() Message ***REMOVED***
	name := mustNewName("example.com.")
	return Message***REMOVED***
		Header: Header***REMOVED***Response: true, Authoritative: true***REMOVED***,
		Questions: []Question***REMOVED***
			***REMOVED***
				Name:  name,
				Type:  TypeA,
				Class: ClassINET,
			***REMOVED***,
		***REMOVED***,
		Answers: []Resource***REMOVED***
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeA,
					Class: ClassINET,
				***REMOVED***,
				&AResource***REMOVED***[4]byte***REMOVED***127, 0, 0, 1***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Authorities: []Resource***REMOVED***
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeA,
					Class: ClassINET,
				***REMOVED***,
				&AResource***REMOVED***[4]byte***REMOVED***127, 0, 0, 1***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Additionals: []Resource***REMOVED***
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeA,
					Class: ClassINET,
				***REMOVED***,
				&AResource***REMOVED***[4]byte***REMOVED***127, 0, 0, 1***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func BenchmarkPack(b *testing.B) ***REMOVED***
	msg := largeTestMsg()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := msg.Pack(); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkAppendPack(b *testing.B) ***REMOVED***
	msg := largeTestMsg()
	buf := make([]byte, 0, packStartingCap)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := msg.AppendPack(buf[:0]); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func largeTestMsg() Message ***REMOVED***
	name := mustNewName("foo.bar.example.com.")
	return Message***REMOVED***
		Header: Header***REMOVED***Response: true, Authoritative: true***REMOVED***,
		Questions: []Question***REMOVED***
			***REMOVED***
				Name:  name,
				Type:  TypeA,
				Class: ClassINET,
			***REMOVED***,
		***REMOVED***,
		Answers: []Resource***REMOVED***
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeA,
					Class: ClassINET,
				***REMOVED***,
				&AResource***REMOVED***[4]byte***REMOVED***127, 0, 0, 1***REMOVED******REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeA,
					Class: ClassINET,
				***REMOVED***,
				&AResource***REMOVED***[4]byte***REMOVED***127, 0, 0, 2***REMOVED******REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeAAAA,
					Class: ClassINET,
				***REMOVED***,
				&AAAAResource***REMOVED***[16]byte***REMOVED***1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16***REMOVED******REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeCNAME,
					Class: ClassINET,
				***REMOVED***,
				&CNAMEResource***REMOVED***mustNewName("alias.example.com.")***REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeSOA,
					Class: ClassINET,
				***REMOVED***,
				&SOAResource***REMOVED***
					NS:      mustNewName("ns1.example.com."),
					MBox:    mustNewName("mb.example.com."),
					Serial:  1,
					Refresh: 2,
					Retry:   3,
					Expire:  4,
					MinTTL:  5,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypePTR,
					Class: ClassINET,
				***REMOVED***,
				&PTRResource***REMOVED***mustNewName("ptr.example.com.")***REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeMX,
					Class: ClassINET,
				***REMOVED***,
				&MXResource***REMOVED***
					7,
					mustNewName("mx.example.com."),
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeSRV,
					Class: ClassINET,
				***REMOVED***,
				&SRVResource***REMOVED***
					8,
					9,
					11,
					mustNewName("srv.example.com."),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Authorities: []Resource***REMOVED***
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeNS,
					Class: ClassINET,
				***REMOVED***,
				&NSResource***REMOVED***mustNewName("ns1.example.com.")***REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeNS,
					Class: ClassINET,
				***REMOVED***,
				&NSResource***REMOVED***mustNewName("ns2.example.com.")***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Additionals: []Resource***REMOVED***
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeTXT,
					Class: ClassINET,
				***REMOVED***,
				&TXTResource***REMOVED***"So Long, and Thanks for All the Fish"***REMOVED***,
			***REMOVED***,
			***REMOVED***
				ResourceHeader***REMOVED***
					Name:  name,
					Type:  TypeTXT,
					Class: ClassINET,
				***REMOVED***,
				&TXTResource***REMOVED***"Hamster Huey and the Gooey Kablooie"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

const loremIpsum = `
Lorem ipsum dolor sit amet, nec enim antiopam id, an ullum choro
nonumes qui, pro eu debet honestatis mediocritatem. No alia enim eos,
magna signiferumque ex vis. Mei no aperiri dissentias, cu vel quas
regione. Malorum quaeque vim ut, eum cu semper aliquid invidunt, ei
nam ipsum assentior.

Nostrum appellantur usu no, vis ex probatus adipiscing. Cu usu illum
facilis eleifend. Iusto conceptam complectitur vim id. Tale omnesque
no usu, ei oblique sadipscing vim. At nullam voluptua usu, mei laudem
reformidans et. Qui ei eros porro reformidans, ius suas veritus
torquatos ex. Mea te facer alterum consequat.

Soleat torquatos democritum sed et, no mea congue appareat, facer
aliquam nec in. Has te ipsum tritani. At justo dicta option nec, movet
phaedrum ad nam. Ea detracto verterem liberavisse has, delectus
suscipiantur in mei. Ex nam meliore complectitur. Ut nam omnis
honestatis quaerendum, ea mea nihil affert detracto, ad vix rebum
mollis.

Ut epicurei praesent neglegentur pri, prima fuisset intellegebat ad
vim. An habemus comprehensam usu, at enim dignissim pro. Eam reque
vivendum adipisci ea. Vel ne odio choro minimum. Sea admodum
dissentiet ex. Mundi tamquam evertitur ius cu. Homero postea iisque ut
pro, vel ne saepe senserit consetetur.

Nulla utamur facilisis ius ea, in viderer diceret pertinax eum. Mei no
enim quodsi facilisi, ex sed aeterno appareat mediocritatem, eum
sententiae deterruisset ut. At suas timeam euismod cum, offendit
appareat interpretaris ne vix. Vel ea civibus albucius, ex vim quidam
accusata intellegebat, noluisse instructior sea id. Nec te nonumes
habemus appellantur, quis dignissim vituperata eu nam.

At vix apeirian patrioque vituperatoribus, an usu agam assum. Debet
iisque an mea. Per eu dicant ponderum accommodare. Pri alienum
placerat senserit an, ne eum ferri abhorreant vituperatoribus. Ut mea
eligendi disputationi. Ius no tation everti impedit, ei magna quidam
mediocritatem pri.

Legendos perpetua iracundia ne usu, no ius ullum epicurei intellegam,
ad modus epicuri lucilius eam. In unum quaerendum usu. Ne diam paulo
has, ea veri virtute sed. Alia honestatis conclusionemque mea eu, ut
iudico albucius his.

Usu essent probatus eu, sed omnis dolor delicatissimi ex. No qui augue
dissentias dissentiet. Laudem recteque no usu, vel an velit noluisse,
an sed utinam eirmod appetere. Ne mea fuisset inimicus ocurreret. At
vis dicant abhorreant, utinam forensibus nec ne, mei te docendi
consequat. Brute inermis persecuti cum id. Ut ipsum munere propriae
usu, dicit graeco disputando id has.

Eros dolore quaerendum nam ei. Timeam ornatus inciderint pro id. Nec
torquatos sadipscing ei, ancillae molestie per in. Malis principes duo
ea, usu liber postulant ei.

Graece timeam voluptatibus eu eam. Alia probatus quo no, ea scripta
feugiat duo. Congue option meliore ex qui, noster invenire appellantur
ea vel. Eu exerci legendos vel. Consetetur repudiandae vim ut. Vix an
probo minimum, et nam illud falli tempor.

Cum dico signiferumque eu. Sed ut regione maiorum, id veritus insolens
tacimates vix. Eu mel sint tamquam lucilius, duo no oporteat
tacimates. Atqui augue concludaturque vix ei, id mel utroque menandri.

Ad oratio blandit aliquando pro. Vis et dolorum rationibus
philosophia, ad cum nulla molestie. Hinc fuisset adversarium eum et,
ne qui nisl verear saperet, vel te quaestio forensibus. Per odio
option delenit an. Alii placerat has no, in pri nihil platonem
cotidieque. Est ut elit copiosae scaevola, debet tollit maluisset sea
an.

Te sea hinc debet pericula, liber ridens fabulas cu sed, quem mutat
accusam mea et. Elitr labitur albucius et pri, an labore feugait mel.
Velit zril melius usu ea. Ad stet putent interpretaris qui. Mel no
error volumus scripserit. In pro paulo iudico, quo ei dolorem
verterem, affert fabellas dissentiet ea vix.

Vis quot deserunt te. Error aliquid detraxit eu usu, vis alia eruditi
salutatus cu. Est nostrud bonorum an, ei usu alii salutatus. Vel at
nisl primis, eum ex aperiri noluisse reformidans. Ad veri velit
utroque vis, ex equidem detraxit temporibus has.

Inermis appareat usu ne. Eros placerat periculis mea ad, in dictas
pericula pro. Errem postulant at usu, ea nec amet ornatus mentitum. Ad
mazim graeco eum, vel ex percipit volutpat iudicabit, sit ne delicata
interesset. Mel sapientem prodesset abhorreant et, oblique suscipit
eam id.

An maluisset disputando mea, vidit mnesarchum pri et. Malis insolens
inciderint no sea. Ea persius maluisset vix, ne vim appellantur
instructior, consul quidam definiebas pri id. Cum integre feugiat
pericula in, ex sed persius similique, mel ne natum dicit percipitur.

Primis discere ne pri, errem putent definitionem at vis. Ei mel dolore
neglegentur, mei tincidunt percipitur ei. Pro ad simul integre
rationibus. Eu vel alii honestatis definitiones, mea no nonumy
reprehendunt.

Dicta appareat legendos est cu. Eu vel congue dicunt omittam, no vix
adhuc minimum constituam, quot noluisse id mel. Eu quot sale mutat
duo, ex nisl munere invenire duo. Ne nec ullum utamur. Pro alterum
debitis nostrum no, ut vel aliquid vivendo.

Aliquip fierent praesent quo ne, id sit audiam recusabo delicatissimi.
Usu postulant incorrupte cu. At pro dicit tibique intellegam, cibo
dolore impedit id eam, et aeque feugait assentior has. Quando sensibus
nec ex. Possit sensibus pri ad, unum mutat periculis cu vix.

Mundi tibique vix te, duo simul partiendo qualisque id, est at vidit
sonet tempor. No per solet aeterno deseruisse. Petentium salutandi
definiebas pri cu. Munere vivendum est in. Ei justo congue eligendi
vis, modus offendit omittantur te mel.

Integre voluptaria in qui, sit habemus tractatos constituam no. Utinam
melius conceptam est ne, quo in minimum apeirian delicata, ut ius
porro recusabo. Dicant expetenda vix no, ludus scripserit sed ex, eu
his modo nostro. Ut etiam sonet his, quodsi inciderint philosophia te
per. Nullam lobortis eu cum, vix an sonet efficiendi repudiandae. Vis
ad idque fabellas intellegebat.

Eum commodo senserit conclusionemque ex. Sed forensibus sadipscing ut,
mei in facer delicata periculis, sea ne hinc putent cetero. Nec ne
alia corpora invenire, alia prima soleat te cum. Eleifend posidonium
nam at.

Dolorum indoctum cu quo, ex dolor legendos recteque eam, cu pri zril
discere. Nec civibus officiis dissentiunt ex, est te liber ludus
elaboraret. Cum ea fabellas invenire. Ex vim nostrud eripuit
comprehensam, nam te inermis delectus, saepe inermis senserit.
`
