// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package charset provides common text encodings for HTML documents.
//
// The mapping from encoding labels to encodings is defined at
// https://encoding.spec.whatwg.org/.
package charset // import "golang.org/x/net/html/charset"

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

// Lookup returns the encoding with the specified label, and its canonical
// name. It returns nil and the empty string if label is not one of the
// standard encodings for HTML. Matching is case-insensitive and ignores
// leading and trailing whitespace. Encoders will use HTML escape sequences for
// runes that are not supported by the character set.
func Lookup(label string) (e encoding.Encoding, name string) ***REMOVED***
	e, err := htmlindex.Get(label)
	if err != nil ***REMOVED***
		return nil, ""
	***REMOVED***
	name, _ = htmlindex.Name(e)
	return &htmlEncoding***REMOVED***e***REMOVED***, name
***REMOVED***

type htmlEncoding struct***REMOVED*** encoding.Encoding ***REMOVED***

func (h *htmlEncoding) NewEncoder() *encoding.Encoder ***REMOVED***
	// HTML requires a non-terminating legacy encoder. We use HTML escapes to
	// substitute unsupported code points.
	return encoding.HTMLEscapeUnsupported(h.Encoding.NewEncoder())
***REMOVED***

// DetermineEncoding determines the encoding of an HTML document by examining
// up to the first 1024 bytes of content and the declared Content-Type.
//
// See http://www.whatwg.org/specs/web-apps/current-work/multipage/parsing.html#determining-the-character-encoding
func DetermineEncoding(content []byte, contentType string) (e encoding.Encoding, name string, certain bool) ***REMOVED***
	if len(content) > 1024 ***REMOVED***
		content = content[:1024]
	***REMOVED***

	for _, b := range boms ***REMOVED***
		if bytes.HasPrefix(content, b.bom) ***REMOVED***
			e, name = Lookup(b.enc)
			return e, name, true
		***REMOVED***
	***REMOVED***

	if _, params, err := mime.ParseMediaType(contentType); err == nil ***REMOVED***
		if cs, ok := params["charset"]; ok ***REMOVED***
			if e, name = Lookup(cs); e != nil ***REMOVED***
				return e, name, true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(content) > 0 ***REMOVED***
		e, name = prescan(content)
		if e != nil ***REMOVED***
			return e, name, false
		***REMOVED***
	***REMOVED***

	// Try to detect UTF-8.
	// First eliminate any partial rune at the end.
	for i := len(content) - 1; i >= 0 && i > len(content)-4; i-- ***REMOVED***
		b := content[i]
		if b < 0x80 ***REMOVED***
			break
		***REMOVED***
		if utf8.RuneStart(b) ***REMOVED***
			content = content[:i]
			break
		***REMOVED***
	***REMOVED***
	hasHighBit := false
	for _, c := range content ***REMOVED***
		if c >= 0x80 ***REMOVED***
			hasHighBit = true
			break
		***REMOVED***
	***REMOVED***
	if hasHighBit && utf8.Valid(content) ***REMOVED***
		return encoding.Nop, "utf-8", false
	***REMOVED***

	// TODO: change default depending on user's locale?
	return charmap.Windows1252, "windows-1252", false
***REMOVED***

// NewReader returns an io.Reader that converts the content of r to UTF-8.
// It calls DetermineEncoding to find out what r's encoding is.
func NewReader(r io.Reader, contentType string) (io.Reader, error) ***REMOVED***
	preview := make([]byte, 1024)
	n, err := io.ReadFull(r, preview)
	switch ***REMOVED***
	case err == io.ErrUnexpectedEOF:
		preview = preview[:n]
		r = bytes.NewReader(preview)
	case err != nil:
		return nil, err
	default:
		r = io.MultiReader(bytes.NewReader(preview), r)
	***REMOVED***

	if e, _, _ := DetermineEncoding(preview, contentType); e != encoding.Nop ***REMOVED***
		r = transform.NewReader(r, e.NewDecoder())
	***REMOVED***
	return r, nil
***REMOVED***

// NewReaderLabel returns a reader that converts from the specified charset to
// UTF-8. It uses Lookup to find the encoding that corresponds to label, and
// returns an error if Lookup returns nil. It is suitable for use as
// encoding/xml.Decoder's CharsetReader function.
func NewReaderLabel(label string, input io.Reader) (io.Reader, error) ***REMOVED***
	e, _ := Lookup(label)
	if e == nil ***REMOVED***
		return nil, fmt.Errorf("unsupported charset: %q", label)
	***REMOVED***
	return transform.NewReader(input, e.NewDecoder()), nil
***REMOVED***

func prescan(content []byte) (e encoding.Encoding, name string) ***REMOVED***
	z := html.NewTokenizer(bytes.NewReader(content))
	for ***REMOVED***
		switch z.Next() ***REMOVED***
		case html.ErrorToken:
			return nil, ""

		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, hasAttr := z.TagName()
			if !bytes.Equal(tagName, []byte("meta")) ***REMOVED***
				continue
			***REMOVED***
			attrList := make(map[string]bool)
			gotPragma := false

			const (
				dontKnow = iota
				doNeedPragma
				doNotNeedPragma
			)
			needPragma := dontKnow

			name = ""
			e = nil
			for hasAttr ***REMOVED***
				var key, val []byte
				key, val, hasAttr = z.TagAttr()
				ks := string(key)
				if attrList[ks] ***REMOVED***
					continue
				***REMOVED***
				attrList[ks] = true
				for i, c := range val ***REMOVED***
					if 'A' <= c && c <= 'Z' ***REMOVED***
						val[i] = c + 0x20
					***REMOVED***
				***REMOVED***

				switch ks ***REMOVED***
				case "http-equiv":
					if bytes.Equal(val, []byte("content-type")) ***REMOVED***
						gotPragma = true
					***REMOVED***

				case "content":
					if e == nil ***REMOVED***
						name = fromMetaElement(string(val))
						if name != "" ***REMOVED***
							e, name = Lookup(name)
							if e != nil ***REMOVED***
								needPragma = doNeedPragma
							***REMOVED***
						***REMOVED***
					***REMOVED***

				case "charset":
					e, name = Lookup(string(val))
					needPragma = doNotNeedPragma
				***REMOVED***
			***REMOVED***

			if needPragma == dontKnow || needPragma == doNeedPragma && !gotPragma ***REMOVED***
				continue
			***REMOVED***

			if strings.HasPrefix(name, "utf-16") ***REMOVED***
				name = "utf-8"
				e = encoding.Nop
			***REMOVED***

			if e != nil ***REMOVED***
				return e, name
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func fromMetaElement(s string) string ***REMOVED***
	for s != "" ***REMOVED***
		csLoc := strings.Index(s, "charset")
		if csLoc == -1 ***REMOVED***
			return ""
		***REMOVED***
		s = s[csLoc+len("charset"):]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if !strings.HasPrefix(s, "=") ***REMOVED***
			continue
		***REMOVED***
		s = s[1:]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if s == "" ***REMOVED***
			return ""
		***REMOVED***
		if q := s[0]; q == '"' || q == '\'' ***REMOVED***
			s = s[1:]
			closeQuote := strings.IndexRune(s, rune(q))
			if closeQuote == -1 ***REMOVED***
				return ""
			***REMOVED***
			return s[:closeQuote]
		***REMOVED***

		end := strings.IndexAny(s, "; \t\n\f\r")
		if end == -1 ***REMOVED***
			end = len(s)
		***REMOVED***
		return s[:end]
	***REMOVED***
	return ""
***REMOVED***

var boms = []struct ***REMOVED***
	bom []byte
	enc string
***REMOVED******REMOVED***
	***REMOVED***[]byte***REMOVED***0xfe, 0xff***REMOVED***, "utf-16be"***REMOVED***,
	***REMOVED***[]byte***REMOVED***0xff, 0xfe***REMOVED***, "utf-16le"***REMOVED***,
	***REMOVED***[]byte***REMOVED***0xef, 0xbb, 0xbf***REMOVED***, "utf-8"***REMOVED***,
***REMOVED***
