// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

// TODO: add tests to check XML responses with the expected prefix path
func TestPrefix(t *testing.T) ***REMOVED***
	const dst, blah = "Destination", "blah blah blah"

	// createLockBody comes from the example in Section 9.10.7.
	const createLockBody = `<?xml version="1.0" encoding="utf-8" ?>
		<D:lockinfo xmlns:D='DAV:'>
			<D:lockscope><D:exclusive/></D:lockscope>
			<D:locktype><D:write/></D:locktype>
			<D:owner>
				<D:href>http://example.org/~ejw/contact.html</D:href>
			</D:owner>
		</D:lockinfo>
	`

	do := func(method, urlStr string, body string, wantStatusCode int, headers ...string) (http.Header, error) ***REMOVED***
		var bodyReader io.Reader
		if body != "" ***REMOVED***
			bodyReader = strings.NewReader(body)
		***REMOVED***
		req, err := http.NewRequest(method, urlStr, bodyReader)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for len(headers) >= 2 ***REMOVED***
			req.Header.Add(headers[0], headers[1])
			headers = headers[2:]
		***REMOVED***
		res, err := http.DefaultTransport.RoundTrip(req)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer res.Body.Close()
		if res.StatusCode != wantStatusCode ***REMOVED***
			return nil, fmt.Errorf("got status code %d, want %d", res.StatusCode, wantStatusCode)
		***REMOVED***
		return res.Header, nil
	***REMOVED***

	prefixes := []string***REMOVED***
		"/",
		"/a/",
		"/a/b/",
		"/a/b/c/",
	***REMOVED***
	ctx := context.Background()
	for _, prefix := range prefixes ***REMOVED***
		fs := NewMemFS()
		h := &Handler***REMOVED***
			FileSystem: fs,
			LockSystem: NewMemLS(),
		***REMOVED***
		mux := http.NewServeMux()
		if prefix != "/" ***REMOVED***
			h.Prefix = prefix
		***REMOVED***
		mux.Handle(prefix, h)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		// The script is:
		//	MKCOL /a
		//	MKCOL /a/b
		//	PUT   /a/b/c
		//	COPY  /a/b/c /a/b/d
		//	MKCOL /a/b/e
		//	MOVE  /a/b/d /a/b/e/f
		//	LOCK  /a/b/e/g
		//	PUT   /a/b/e/g
		// which should yield the (possibly stripped) filenames /a/b/c,
		// /a/b/e/f and /a/b/e/g, plus their parent directories.

		wantA := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusMovedPermanently,
			"/a/b/":   http.StatusNotFound,
			"/a/b/c/": http.StatusNotFound,
		***REMOVED***[prefix]
		if _, err := do("MKCOL", srv.URL+"/a", "", wantA); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q MKCOL /a: %v", prefix, err)
			continue
		***REMOVED***

		wantB := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusCreated,
			"/a/b/":   http.StatusMovedPermanently,
			"/a/b/c/": http.StatusNotFound,
		***REMOVED***[prefix]
		if _, err := do("MKCOL", srv.URL+"/a/b", "", wantB); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q MKCOL /a/b: %v", prefix, err)
			continue
		***REMOVED***

		wantC := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusCreated,
			"/a/b/":   http.StatusCreated,
			"/a/b/c/": http.StatusMovedPermanently,
		***REMOVED***[prefix]
		if _, err := do("PUT", srv.URL+"/a/b/c", blah, wantC); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q PUT /a/b/c: %v", prefix, err)
			continue
		***REMOVED***

		wantD := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusCreated,
			"/a/b/":   http.StatusCreated,
			"/a/b/c/": http.StatusMovedPermanently,
		***REMOVED***[prefix]
		if _, err := do("COPY", srv.URL+"/a/b/c", "", wantD, dst, srv.URL+"/a/b/d"); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q COPY /a/b/c /a/b/d: %v", prefix, err)
			continue
		***REMOVED***

		wantE := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusCreated,
			"/a/b/":   http.StatusCreated,
			"/a/b/c/": http.StatusNotFound,
		***REMOVED***[prefix]
		if _, err := do("MKCOL", srv.URL+"/a/b/e", "", wantE); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q MKCOL /a/b/e: %v", prefix, err)
			continue
		***REMOVED***

		wantF := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusCreated,
			"/a/b/":   http.StatusCreated,
			"/a/b/c/": http.StatusNotFound,
		***REMOVED***[prefix]
		if _, err := do("MOVE", srv.URL+"/a/b/d", "", wantF, dst, srv.URL+"/a/b/e/f"); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q MOVE /a/b/d /a/b/e/f: %v", prefix, err)
			continue
		***REMOVED***

		var lockToken string
		wantG := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusCreated,
			"/a/b/":   http.StatusCreated,
			"/a/b/c/": http.StatusNotFound,
		***REMOVED***[prefix]
		if h, err := do("LOCK", srv.URL+"/a/b/e/g", createLockBody, wantG); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q LOCK /a/b/e/g: %v", prefix, err)
			continue
		***REMOVED*** else ***REMOVED***
			lockToken = h.Get("Lock-Token")
		***REMOVED***

		ifHeader := fmt.Sprintf("<%s/a/b/e/g> (%s)", srv.URL, lockToken)
		wantH := map[string]int***REMOVED***
			"/":       http.StatusCreated,
			"/a/":     http.StatusCreated,
			"/a/b/":   http.StatusCreated,
			"/a/b/c/": http.StatusNotFound,
		***REMOVED***[prefix]
		if _, err := do("PUT", srv.URL+"/a/b/e/g", blah, wantH, "If", ifHeader); err != nil ***REMOVED***
			t.Errorf("prefix=%-9q PUT /a/b/e/g: %v", prefix, err)
			continue
		***REMOVED***

		got, err := find(ctx, nil, fs, "/")
		if err != nil ***REMOVED***
			t.Errorf("prefix=%-9q find: %v", prefix, err)
			continue
		***REMOVED***
		sort.Strings(got)
		want := map[string][]string***REMOVED***
			"/":       ***REMOVED***"/", "/a", "/a/b", "/a/b/c", "/a/b/e", "/a/b/e/f", "/a/b/e/g"***REMOVED***,
			"/a/":     ***REMOVED***"/", "/b", "/b/c", "/b/e", "/b/e/f", "/b/e/g"***REMOVED***,
			"/a/b/":   ***REMOVED***"/", "/c", "/e", "/e/f", "/e/g"***REMOVED***,
			"/a/b/c/": ***REMOVED***"/"***REMOVED***,
		***REMOVED***[prefix]
		if !reflect.DeepEqual(got, want) ***REMOVED***
			t.Errorf("prefix=%-9q find:\ngot  %v\nwant %v", prefix, got, want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEscapeXML(t *testing.T) ***REMOVED***
	// These test cases aren't exhaustive, and there is more than one way to
	// escape e.g. a quot (as "&#34;" or "&quot;") or an apos. We presume that
	// the encoding/xml package tests xml.EscapeText more thoroughly. This test
	// here is just a sanity check for this package's escapeXML function, and
	// its attempt to provide a fast path (and avoid a bytes.Buffer allocation)
	// when escaping filenames is obviously a no-op.
	testCases := map[string]string***REMOVED***
		"":              "",
		" ":             " ",
		"&":             "&amp;",
		"*":             "*",
		"+":             "+",
		",":             ",",
		"-":             "-",
		".":             ".",
		"/":             "/",
		"0":             "0",
		"9":             "9",
		":":             ":",
		"<":             "&lt;",
		">":             "&gt;",
		"A":             "A",
		"_":             "_",
		"a":             "a",
		"~":             "~",
		"\u0201":        "\u0201",
		"&amp;":         "&amp;amp;",
		"foo&<b/ar>baz": "foo&amp;&lt;b/ar&gt;baz",
	***REMOVED***

	for in, want := range testCases ***REMOVED***
		if got := escapeXML(in); got != want ***REMOVED***
			t.Errorf("in=%q: got %q, want %q", in, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFilenameEscape(t *testing.T) ***REMOVED***
	hrefRe := regexp.MustCompile(`<D:href>([^<]*)</D:href>`)
	displayNameRe := regexp.MustCompile(`<D:displayname>([^<]*)</D:displayname>`)
	do := func(method, urlStr string) (string, string, error) ***REMOVED***
		req, err := http.NewRequest(method, urlStr, nil)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		res, err := http.DefaultClient.Do(req)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		hrefMatch := hrefRe.FindStringSubmatch(string(b))
		if len(hrefMatch) != 2 ***REMOVED***
			return "", "", errors.New("D:href not found")
		***REMOVED***
		displayNameMatch := displayNameRe.FindStringSubmatch(string(b))
		if len(displayNameMatch) != 2 ***REMOVED***
			return "", "", errors.New("D:displayname not found")
		***REMOVED***

		return hrefMatch[1], displayNameMatch[1], nil
	***REMOVED***

	testCases := []struct ***REMOVED***
		name, wantHref, wantDisplayName string
	***REMOVED******REMOVED******REMOVED***
		name:            `/foo%bar`,
		wantHref:        `/foo%25bar`,
		wantDisplayName: `foo%bar`,
	***REMOVED***, ***REMOVED***
		name:            `/こんにちわ世界`,
		wantHref:        `/%E3%81%93%E3%82%93%E3%81%AB%E3%81%A1%E3%82%8F%E4%B8%96%E7%95%8C`,
		wantDisplayName: `こんにちわ世界`,
	***REMOVED***, ***REMOVED***
		name:            `/Program Files/`,
		wantHref:        `/Program%20Files`,
		wantDisplayName: `Program Files`,
	***REMOVED***, ***REMOVED***
		name:            `/go+lang`,
		wantHref:        `/go+lang`,
		wantDisplayName: `go+lang`,
	***REMOVED***, ***REMOVED***
		name:            `/go&lang`,
		wantHref:        `/go&amp;lang`,
		wantDisplayName: `go&amp;lang`,
	***REMOVED***, ***REMOVED***
		name:            `/go<lang`,
		wantHref:        `/go%3Clang`,
		wantDisplayName: `go&lt;lang`,
	***REMOVED******REMOVED***
	ctx := context.Background()
	fs := NewMemFS()
	for _, tc := range testCases ***REMOVED***
		if strings.HasSuffix(tc.name, "/") ***REMOVED***
			if err := fs.Mkdir(ctx, tc.name, 0755); err != nil ***REMOVED***
				t.Fatalf("name=%q: Mkdir: %v", tc.name, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			f, err := fs.OpenFile(ctx, tc.name, os.O_CREATE, 0644)
			if err != nil ***REMOVED***
				t.Fatalf("name=%q: OpenFile: %v", tc.name, err)
			***REMOVED***
			f.Close()
		***REMOVED***
	***REMOVED***

	srv := httptest.NewServer(&Handler***REMOVED***
		FileSystem: fs,
		LockSystem: NewMemLS(),
	***REMOVED***)
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		u.Path = tc.name
		gotHref, gotDisplayName, err := do("PROPFIND", u.String())
		if err != nil ***REMOVED***
			t.Errorf("name=%q: PROPFIND: %v", tc.name, err)
			continue
		***REMOVED***
		if gotHref != tc.wantHref ***REMOVED***
			t.Errorf("name=%q: got href %q, want %q", tc.name, gotHref, tc.wantHref)
		***REMOVED***
		if gotDisplayName != tc.wantDisplayName ***REMOVED***
			t.Errorf("name=%q: got dispayname %q, want %q", tc.name, gotDisplayName, tc.wantDisplayName)
		***REMOVED***
	***REMOVED***
***REMOVED***
