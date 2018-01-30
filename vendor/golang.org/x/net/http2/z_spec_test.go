// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var coverSpec = flag.Bool("coverspec", false, "Run spec coverage tests")

// The global map of sentence coverage for the http2 spec.
var defaultSpecCoverage specCoverage

var loadSpecOnce sync.Once

func loadSpec() ***REMOVED***
	if f, err := os.Open("testdata/draft-ietf-httpbis-http2.xml"); err != nil ***REMOVED***
		panic(err)
	***REMOVED*** else ***REMOVED***
		defaultSpecCoverage = readSpecCov(f)
		f.Close()
	***REMOVED***
***REMOVED***

// covers marks all sentences for section sec in defaultSpecCoverage. Sentences not
// "covered" will be included in report outputted by TestSpecCoverage.
func covers(sec, sentences string) ***REMOVED***
	loadSpecOnce.Do(loadSpec)
	defaultSpecCoverage.cover(sec, sentences)
***REMOVED***

type specPart struct ***REMOVED***
	section  string
	sentence string
***REMOVED***

func (ss specPart) Less(oo specPart) bool ***REMOVED***
	atoi := func(s string) int ***REMOVED***
		n, err := strconv.Atoi(s)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		return n
	***REMOVED***
	a := strings.Split(ss.section, ".")
	b := strings.Split(oo.section, ".")
	for len(a) > 0 ***REMOVED***
		if len(b) == 0 ***REMOVED***
			return false
		***REMOVED***
		x, y := atoi(a[0]), atoi(b[0])
		if x == y ***REMOVED***
			a, b = a[1:], b[1:]
			continue
		***REMOVED***
		return x < y
	***REMOVED***
	if len(b) > 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

type bySpecSection []specPart

func (a bySpecSection) Len() int           ***REMOVED*** return len(a) ***REMOVED***
func (a bySpecSection) Less(i, j int) bool ***REMOVED*** return a[i].Less(a[j]) ***REMOVED***
func (a bySpecSection) Swap(i, j int)      ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***

type specCoverage struct ***REMOVED***
	coverage map[specPart]bool
	d        *xml.Decoder
***REMOVED***

func joinSection(sec []int) string ***REMOVED***
	s := fmt.Sprintf("%d", sec[0])
	for _, n := range sec[1:] ***REMOVED***
		s = fmt.Sprintf("%s.%d", s, n)
	***REMOVED***
	return s
***REMOVED***

func (sc specCoverage) readSection(sec []int) ***REMOVED***
	var (
		buf = new(bytes.Buffer)
		sub = 0
	)
	for ***REMOVED***
		tk, err := sc.d.Token()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				return
			***REMOVED***
			panic(err)
		***REMOVED***
		switch v := tk.(type) ***REMOVED***
		case xml.StartElement:
			if skipElement(v) ***REMOVED***
				if err := sc.d.Skip(); err != nil ***REMOVED***
					panic(err)
				***REMOVED***
				if v.Name.Local == "section" ***REMOVED***
					sub++
				***REMOVED***
				break
			***REMOVED***
			switch v.Name.Local ***REMOVED***
			case "section":
				sub++
				sc.readSection(append(sec, sub))
			case "xref":
				buf.Write(sc.readXRef(v))
			***REMOVED***
		case xml.CharData:
			if len(sec) == 0 ***REMOVED***
				break
			***REMOVED***
			buf.Write(v)
		case xml.EndElement:
			if v.Name.Local == "section" ***REMOVED***
				sc.addSentences(joinSection(sec), buf.String())
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sc specCoverage) readXRef(se xml.StartElement) []byte ***REMOVED***
	var b []byte
	for ***REMOVED***
		tk, err := sc.d.Token()
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		switch v := tk.(type) ***REMOVED***
		case xml.CharData:
			if b != nil ***REMOVED***
				panic("unexpected CharData")
			***REMOVED***
			b = []byte(string(v))
		case xml.EndElement:
			if v.Name.Local != "xref" ***REMOVED***
				panic("expected </xref>")
			***REMOVED***
			if b != nil ***REMOVED***
				return b
			***REMOVED***
			sig := attrSig(se)
			switch sig ***REMOVED***
			case "target":
				return []byte(fmt.Sprintf("[%s]", attrValue(se, "target")))
			case "fmt-of,rel,target", "fmt-,,rel,target":
				return []byte(fmt.Sprintf("[%s, %s]", attrValue(se, "target"), attrValue(se, "rel")))
			case "fmt-of,sec,target", "fmt-,,sec,target":
				return []byte(fmt.Sprintf("[section %s of %s]", attrValue(se, "sec"), attrValue(se, "target")))
			case "fmt-of,rel,sec,target":
				return []byte(fmt.Sprintf("[section %s of %s, %s]", attrValue(se, "sec"), attrValue(se, "target"), attrValue(se, "rel")))
			default:
				panic(fmt.Sprintf("unknown attribute signature %q in %#v", sig, fmt.Sprintf("%#v", se)))
			***REMOVED***
		default:
			panic(fmt.Sprintf("unexpected tag %q", v))
		***REMOVED***
	***REMOVED***
***REMOVED***

var skipAnchor = map[string]bool***REMOVED***
	"intro":    true,
	"Overview": true,
***REMOVED***

var skipTitle = map[string]bool***REMOVED***
	"Acknowledgements":            true,
	"Change Log":                  true,
	"Document Organization":       true,
	"Conventions and Terminology": true,
***REMOVED***

func skipElement(s xml.StartElement) bool ***REMOVED***
	switch s.Name.Local ***REMOVED***
	case "artwork":
		return true
	case "section":
		for _, attr := range s.Attr ***REMOVED***
			switch attr.Name.Local ***REMOVED***
			case "anchor":
				if skipAnchor[attr.Value] || strings.HasPrefix(attr.Value, "changes.since.") ***REMOVED***
					return true
				***REMOVED***
			case "title":
				if skipTitle[attr.Value] ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func readSpecCov(r io.Reader) specCoverage ***REMOVED***
	sc := specCoverage***REMOVED***
		coverage: map[specPart]bool***REMOVED******REMOVED***,
		d:        xml.NewDecoder(r)***REMOVED***
	sc.readSection(nil)
	return sc
***REMOVED***

func (sc specCoverage) addSentences(sec string, sentence string) ***REMOVED***
	for _, s := range parseSentences(sentence) ***REMOVED***
		sc.coverage[specPart***REMOVED***sec, s***REMOVED***] = false
	***REMOVED***
***REMOVED***

func (sc specCoverage) cover(sec string, sentence string) ***REMOVED***
	for _, s := range parseSentences(sentence) ***REMOVED***
		p := specPart***REMOVED***sec, s***REMOVED***
		if _, ok := sc.coverage[p]; !ok ***REMOVED***
			panic(fmt.Sprintf("Not found in spec: %q, %q", sec, s))
		***REMOVED***
		sc.coverage[specPart***REMOVED***sec, s***REMOVED***] = true
	***REMOVED***

***REMOVED***

var whitespaceRx = regexp.MustCompile(`\s+`)

func parseSentences(sens string) []string ***REMOVED***
	sens = strings.TrimSpace(sens)
	if sens == "" ***REMOVED***
		return nil
	***REMOVED***
	ss := strings.Split(whitespaceRx.ReplaceAllString(sens, " "), ". ")
	for i, s := range ss ***REMOVED***
		s = strings.TrimSpace(s)
		if !strings.HasSuffix(s, ".") ***REMOVED***
			s += "."
		***REMOVED***
		ss[i] = s
	***REMOVED***
	return ss
***REMOVED***

func TestSpecParseSentences(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		ss   string
		want []string
	***REMOVED******REMOVED***
		***REMOVED***"Sentence 1. Sentence 2.",
			[]string***REMOVED***
				"Sentence 1.",
				"Sentence 2.",
			***REMOVED******REMOVED***,
		***REMOVED***"Sentence 1.  \nSentence 2.\tSentence 3.",
			[]string***REMOVED***
				"Sentence 1.",
				"Sentence 2.",
				"Sentence 3.",
			***REMOVED******REMOVED***,
	***REMOVED***

	for i, tt := range tests ***REMOVED***
		got := parseSentences(tt.ss)
		if !reflect.DeepEqual(got, tt.want) ***REMOVED***
			t.Errorf("%d: got = %q, want %q", i, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSpecCoverage(t *testing.T) ***REMOVED***
	if !*coverSpec ***REMOVED***
		t.Skip()
	***REMOVED***

	loadSpecOnce.Do(loadSpec)

	var (
		list     []specPart
		cv       = defaultSpecCoverage.coverage
		total    = len(cv)
		complete = 0
	)

	for sp, touched := range defaultSpecCoverage.coverage ***REMOVED***
		if touched ***REMOVED***
			complete++
		***REMOVED*** else ***REMOVED***
			list = append(list, sp)
		***REMOVED***
	***REMOVED***
	sort.Stable(bySpecSection(list))

	if testing.Short() && len(list) > 5 ***REMOVED***
		list = list[:5]
	***REMOVED***

	for _, p := range list ***REMOVED***
		t.Errorf("\tSECTION %s: %s", p.section, p.sentence)
	***REMOVED***

	t.Logf("%d/%d (%d%%) sentences covered", complete, total, (complete/total)*100)
***REMOVED***

func attrSig(se xml.StartElement) string ***REMOVED***
	var names []string
	for _, attr := range se.Attr ***REMOVED***
		if attr.Name.Local == "fmt" ***REMOVED***
			names = append(names, "fmt-"+attr.Value)
		***REMOVED*** else ***REMOVED***
			names = append(names, attr.Name.Local)
		***REMOVED***
	***REMOVED***
	sort.Strings(names)
	return strings.Join(names, ",")
***REMOVED***

func attrValue(se xml.StartElement, attr string) string ***REMOVED***
	for _, a := range se.Attr ***REMOVED***
		if a.Name.Local == attr ***REMOVED***
			return a.Value
		***REMOVED***
	***REMOVED***
	panic("unknown attribute " + attr)
***REMOVED***

func TestSpecPartLess(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		sec1, sec2 string
		want       bool
	***REMOVED******REMOVED***
		***REMOVED***"6.2.1", "6.2", false***REMOVED***,
		***REMOVED***"6.2", "6.2.1", true***REMOVED***,
		***REMOVED***"6.10", "6.10.1", true***REMOVED***,
		***REMOVED***"6.10", "6.1.1", false***REMOVED***, // 10, not 1
		***REMOVED***"6.1", "6.1", false***REMOVED***,    // equal, so not less
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		got := (specPart***REMOVED***tt.sec1, "foo"***REMOVED***).Less(specPart***REMOVED***tt.sec2, "foo"***REMOVED***)
		if got != tt.want ***REMOVED***
			t.Errorf("Less(%q, %q) = %v; want %v", tt.sec1, tt.sec2, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
