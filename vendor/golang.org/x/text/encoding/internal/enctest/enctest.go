// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enctest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)

// Encoder or Decoder
type Transcoder interface ***REMOVED***
	transform.Transformer
	Bytes([]byte) ([]byte, error)
	String(string) (string, error)
***REMOVED***

func TestEncoding(t *testing.T, e encoding.Encoding, encoded, utf8, prefix, suffix string) ***REMOVED***
	for _, direction := range []string***REMOVED***"Decode", "Encode"***REMOVED*** ***REMOVED***
		t.Run(fmt.Sprintf("%v/%s", e, direction), func(t *testing.T) ***REMOVED***

			var coder Transcoder
			var want, src, wPrefix, sPrefix, wSuffix, sSuffix string
			if direction == "Decode" ***REMOVED***
				coder, want, src = e.NewDecoder(), utf8, encoded
				wPrefix, sPrefix, wSuffix, sSuffix = "", prefix, "", suffix
			***REMOVED*** else ***REMOVED***
				coder, want, src = e.NewEncoder(), encoded, utf8
				wPrefix, sPrefix, wSuffix, sSuffix = prefix, "", suffix, ""
			***REMOVED***

			dst := make([]byte, len(wPrefix)+len(want)+len(wSuffix))
			nDst, nSrc, err := coder.Transform(dst, []byte(sPrefix+src+sSuffix), true)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if nDst != len(wPrefix)+len(want)+len(wSuffix) ***REMOVED***
				t.Fatalf("nDst got %d, want %d",
					nDst, len(wPrefix)+len(want)+len(wSuffix))
			***REMOVED***
			if nSrc != len(sPrefix)+len(src)+len(sSuffix) ***REMOVED***
				t.Fatalf("nSrc got %d, want %d",
					nSrc, len(sPrefix)+len(src)+len(sSuffix))
			***REMOVED***
			if got := string(dst); got != wPrefix+want+wSuffix ***REMOVED***
				t.Fatalf("\ngot  %q\nwant %q", got, wPrefix+want+wSuffix)
			***REMOVED***

			for _, n := range []int***REMOVED***0, 1, 2, 10, 123, 4567***REMOVED*** ***REMOVED***
				input := sPrefix + strings.Repeat(src, n) + sSuffix
				g, err := coder.String(input)
				if err != nil ***REMOVED***
					t.Fatalf("Bytes: n=%d: %v", n, err)
				***REMOVED***
				if len(g) == 0 && len(input) == 0 ***REMOVED***
					// If the input is empty then the output can be empty,
					// regardless of whatever wPrefix is.
					continue
				***REMOVED***
				got1, want1 := string(g), wPrefix+strings.Repeat(want, n)+wSuffix
				if got1 != want1 ***REMOVED***
					t.Fatalf("ReadAll: n=%d\ngot  %q\nwant %q",
						n, trim(got1), trim(want1))
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestFile(t *testing.T, e encoding.Encoding) ***REMOVED***
	for _, dir := range []string***REMOVED***"Decode", "Encode"***REMOVED*** ***REMOVED***
		t.Run(fmt.Sprintf("%s/%s", e, dir), func(t *testing.T) ***REMOVED***
			dst, src, transformer, err := load(dir, e)
			if err != nil ***REMOVED***
				t.Fatalf("load: %v", err)
			***REMOVED***
			buf, err := transformer.Bytes(src)
			if err != nil ***REMOVED***
				t.Fatalf("transform: %v", err)
			***REMOVED***
			if !bytes.Equal(buf, dst) ***REMOVED***
				t.Error("transformed bytes did not match golden file")
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func Benchmark(b *testing.B, enc encoding.Encoding) ***REMOVED***
	for _, direction := range []string***REMOVED***"Decode", "Encode"***REMOVED*** ***REMOVED***
		b.Run(fmt.Sprintf("%s/%s", enc, direction), func(b *testing.B) ***REMOVED***
			_, src, transformer, err := load(direction, enc)
			if err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
			b.SetBytes(int64(len(src)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				r := transform.NewReader(bytes.NewReader(src), transformer)
				io.Copy(ioutil.Discard, r)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// testdataFiles are files in testdata/*.txt.
var testdataFiles = []struct ***REMOVED***
	mib           identifier.MIB
	basename, ext string
***REMOVED******REMOVED***
	***REMOVED***identifier.Windows1252, "candide", "windows-1252"***REMOVED***,
	***REMOVED***identifier.EUCPkdFmtJapanese, "rashomon", "euc-jp"***REMOVED***,
	***REMOVED***identifier.ISO2022JP, "rashomon", "iso-2022-jp"***REMOVED***,
	***REMOVED***identifier.ShiftJIS, "rashomon", "shift-jis"***REMOVED***,
	***REMOVED***identifier.EUCKR, "unsu-joh-eun-nal", "euc-kr"***REMOVED***,
	***REMOVED***identifier.GBK, "sunzi-bingfa-simplified", "gbk"***REMOVED***,
	***REMOVED***identifier.HZGB2312, "sunzi-bingfa-gb-levels-1-and-2", "hz-gb2312"***REMOVED***,
	***REMOVED***identifier.Big5, "sunzi-bingfa-traditional", "big5"***REMOVED***,
	***REMOVED***identifier.UTF16LE, "candide", "utf-16le"***REMOVED***,
	***REMOVED***identifier.UTF8, "candide", "utf-8"***REMOVED***,
	***REMOVED***identifier.UTF32BE, "candide", "utf-32be"***REMOVED***,

	// GB18030 is a superset of GBK and is nominally a Simplified Chinese
	// encoding, but it can also represent the entire Basic Multilingual
	// Plane, including codepoints like 'â' that aren't encodable by GBK.
	// GB18030 on Simplified Chinese should perform similarly to GBK on
	// Simplified Chinese. GB18030 on "candide" is more interesting.
	***REMOVED***identifier.GB18030, "candide", "gb18030"***REMOVED***,
***REMOVED***

func load(direction string, enc encoding.Encoding) ([]byte, []byte, Transcoder, error) ***REMOVED***
	basename, ext, count := "", "", 0
	for _, tf := range testdataFiles ***REMOVED***
		if mib, _ := enc.(identifier.Interface).ID(); tf.mib == mib ***REMOVED***
			basename, ext = tf.basename, tf.ext
			count++
		***REMOVED***
	***REMOVED***
	if count != 1 ***REMOVED***
		if count == 0 ***REMOVED***
			return nil, nil, nil, fmt.Errorf("no testdataFiles for %s", enc)
		***REMOVED***
		return nil, nil, nil, fmt.Errorf("too many testdataFiles for %s", enc)
	***REMOVED***
	dstFile := fmt.Sprintf("../testdata/%s-%s.txt", basename, ext)
	srcFile := fmt.Sprintf("../testdata/%s-utf-8.txt", basename)
	var coder Transcoder = encoding.ReplaceUnsupported(enc.NewEncoder())
	if direction == "Decode" ***REMOVED***
		dstFile, srcFile = srcFile, dstFile
		coder = enc.NewDecoder()
	***REMOVED***
	dst, err := ioutil.ReadFile(dstFile)
	if err != nil ***REMOVED***
		if dst, err = ioutil.ReadFile("../" + dstFile); err != nil ***REMOVED***
			return nil, nil, nil, err
		***REMOVED***
	***REMOVED***
	src, err := ioutil.ReadFile(srcFile)
	if err != nil ***REMOVED***
		if src, err = ioutil.ReadFile("../" + srcFile); err != nil ***REMOVED***
			return nil, nil, nil, err
		***REMOVED***
	***REMOVED***
	return dst, src, coder, nil
***REMOVED***

func trim(s string) string ***REMOVED***
	if len(s) < 120 ***REMOVED***
		return s
	***REMOVED***
	return s[:50] + "..." + s[len(s)-50:]
***REMOVED***
