// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hpack

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
)

func (d *Decoder) mustAt(idx int) HeaderField ***REMOVED***
	if hf, ok := d.at(uint64(idx)); !ok ***REMOVED***
		panic(fmt.Sprintf("bogus index %d", idx))
	***REMOVED*** else ***REMOVED***
		return hf
	***REMOVED***
***REMOVED***

func TestDynamicTableAt(t *testing.T) ***REMOVED***
	d := NewDecoder(4096, nil)
	at := d.mustAt
	if got, want := at(2), (pair(":method", "GET")); got != want ***REMOVED***
		t.Errorf("at(2) = %v; want %v", got, want)
	***REMOVED***
	d.dynTab.add(pair("foo", "bar"))
	d.dynTab.add(pair("blake", "miz"))
	if got, want := at(staticTable.len()+1), (pair("blake", "miz")); got != want ***REMOVED***
		t.Errorf("at(dyn 1) = %v; want %v", got, want)
	***REMOVED***
	if got, want := at(staticTable.len()+2), (pair("foo", "bar")); got != want ***REMOVED***
		t.Errorf("at(dyn 2) = %v; want %v", got, want)
	***REMOVED***
	if got, want := at(3), (pair(":method", "POST")); got != want ***REMOVED***
		t.Errorf("at(3) = %v; want %v", got, want)
	***REMOVED***
***REMOVED***

func TestDynamicTableSizeEvict(t *testing.T) ***REMOVED***
	d := NewDecoder(4096, nil)
	if want := uint32(0); d.dynTab.size != want ***REMOVED***
		t.Fatalf("size = %d; want %d", d.dynTab.size, want)
	***REMOVED***
	add := d.dynTab.add
	add(pair("blake", "eats pizza"))
	if want := uint32(15 + 32); d.dynTab.size != want ***REMOVED***
		t.Fatalf("after pizza, size = %d; want %d", d.dynTab.size, want)
	***REMOVED***
	add(pair("foo", "bar"))
	if want := uint32(15 + 32 + 6 + 32); d.dynTab.size != want ***REMOVED***
		t.Fatalf("after foo bar, size = %d; want %d", d.dynTab.size, want)
	***REMOVED***
	d.dynTab.setMaxSize(15 + 32 + 1 /* slop */)
	if want := uint32(6 + 32); d.dynTab.size != want ***REMOVED***
		t.Fatalf("after setMaxSize, size = %d; want %d", d.dynTab.size, want)
	***REMOVED***
	if got, want := d.mustAt(staticTable.len()+1), (pair("foo", "bar")); got != want ***REMOVED***
		t.Errorf("at(dyn 1) = %v; want %v", got, want)
	***REMOVED***
	add(pair("long", strings.Repeat("x", 500)))
	if want := uint32(0); d.dynTab.size != want ***REMOVED***
		t.Fatalf("after big one, size = %d; want %d", d.dynTab.size, want)
	***REMOVED***
***REMOVED***

func TestDecoderDecode(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		name       string
		in         []byte
		want       []HeaderField
		wantDynTab []HeaderField // newest entry first
	***REMOVED******REMOVED***
		// C.2.1 Literal Header Field with Indexing
		// http://http2.github.io/http2-spec/compression.html#rfc.section.C.2.1
		***REMOVED***"C.2.1", dehex("400a 6375 7374 6f6d 2d6b 6579 0d63 7573 746f 6d2d 6865 6164 6572"),
			[]HeaderField***REMOVED***pair("custom-key", "custom-header")***REMOVED***,
			[]HeaderField***REMOVED***pair("custom-key", "custom-header")***REMOVED***,
		***REMOVED***,

		// C.2.2 Literal Header Field without Indexing
		// http://http2.github.io/http2-spec/compression.html#rfc.section.C.2.2
		***REMOVED***"C.2.2", dehex("040c 2f73 616d 706c 652f 7061 7468"),
			[]HeaderField***REMOVED***pair(":path", "/sample/path")***REMOVED***,
			[]HeaderField***REMOVED******REMOVED******REMOVED***,

		// C.2.3 Literal Header Field never Indexed
		// http://http2.github.io/http2-spec/compression.html#rfc.section.C.2.3
		***REMOVED***"C.2.3", dehex("1008 7061 7373 776f 7264 0673 6563 7265 74"),
			[]HeaderField***REMOVED******REMOVED***"password", "secret", true***REMOVED******REMOVED***,
			[]HeaderField***REMOVED******REMOVED******REMOVED***,

		// C.2.4 Indexed Header Field
		// http://http2.github.io/http2-spec/compression.html#rfc.section.C.2.4
		***REMOVED***"C.2.4", []byte("\x82"),
			[]HeaderField***REMOVED***pair(":method", "GET")***REMOVED***,
			[]HeaderField***REMOVED******REMOVED******REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		d := NewDecoder(4096, nil)
		hf, err := d.DecodeFull(tt.in)
		if err != nil ***REMOVED***
			t.Errorf("%s: %v", tt.name, err)
			continue
		***REMOVED***
		if !reflect.DeepEqual(hf, tt.want) ***REMOVED***
			t.Errorf("%s: Got %v; want %v", tt.name, hf, tt.want)
		***REMOVED***
		gotDynTab := d.dynTab.reverseCopy()
		if !reflect.DeepEqual(gotDynTab, tt.wantDynTab) ***REMOVED***
			t.Errorf("%s: dynamic table after = %v; want %v", tt.name, gotDynTab, tt.wantDynTab)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (dt *dynamicTable) reverseCopy() (hf []HeaderField) ***REMOVED***
	hf = make([]HeaderField, len(dt.table.ents))
	for i := range hf ***REMOVED***
		hf[i] = dt.table.ents[len(dt.table.ents)-1-i]
	***REMOVED***
	return
***REMOVED***

type encAndWant struct ***REMOVED***
	enc         []byte
	want        []HeaderField
	wantDynTab  []HeaderField
	wantDynSize uint32
***REMOVED***

// C.3 Request Examples without Huffman Coding
// http://http2.github.io/http2-spec/compression.html#rfc.section.C.3
func TestDecodeC3_NoHuffman(t *testing.T) ***REMOVED***
	testDecodeSeries(t, 4096, []encAndWant***REMOVED***
		***REMOVED***dehex("8286 8441 0f77 7777 2e65 7861 6d70 6c65 2e63 6f6d"),
			[]HeaderField***REMOVED***
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", "www.example.com"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair(":authority", "www.example.com"),
			***REMOVED***,
			57,
		***REMOVED***,
		***REMOVED***dehex("8286 84be 5808 6e6f 2d63 6163 6865"),
			[]HeaderField***REMOVED***
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", "www.example.com"),
				pair("cache-control", "no-cache"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("cache-control", "no-cache"),
				pair(":authority", "www.example.com"),
			***REMOVED***,
			110,
		***REMOVED***,
		***REMOVED***dehex("8287 85bf 400a 6375 7374 6f6d 2d6b 6579 0c63 7573 746f 6d2d 7661 6c75 65"),
			[]HeaderField***REMOVED***
				pair(":method", "GET"),
				pair(":scheme", "https"),
				pair(":path", "/index.html"),
				pair(":authority", "www.example.com"),
				pair("custom-key", "custom-value"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("custom-key", "custom-value"),
				pair("cache-control", "no-cache"),
				pair(":authority", "www.example.com"),
			***REMOVED***,
			164,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// C.4 Request Examples with Huffman Coding
// http://http2.github.io/http2-spec/compression.html#rfc.section.C.4
func TestDecodeC4_Huffman(t *testing.T) ***REMOVED***
	testDecodeSeries(t, 4096, []encAndWant***REMOVED***
		***REMOVED***dehex("8286 8441 8cf1 e3c2 e5f2 3a6b a0ab 90f4 ff"),
			[]HeaderField***REMOVED***
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", "www.example.com"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair(":authority", "www.example.com"),
			***REMOVED***,
			57,
		***REMOVED***,
		***REMOVED***dehex("8286 84be 5886 a8eb 1064 9cbf"),
			[]HeaderField***REMOVED***
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", "www.example.com"),
				pair("cache-control", "no-cache"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("cache-control", "no-cache"),
				pair(":authority", "www.example.com"),
			***REMOVED***,
			110,
		***REMOVED***,
		***REMOVED***dehex("8287 85bf 4088 25a8 49e9 5ba9 7d7f 8925 a849 e95b b8e8 b4bf"),
			[]HeaderField***REMOVED***
				pair(":method", "GET"),
				pair(":scheme", "https"),
				pair(":path", "/index.html"),
				pair(":authority", "www.example.com"),
				pair("custom-key", "custom-value"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("custom-key", "custom-value"),
				pair("cache-control", "no-cache"),
				pair(":authority", "www.example.com"),
			***REMOVED***,
			164,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// http://http2.github.io/http2-spec/compression.html#rfc.section.C.5
// "This section shows several consecutive header lists, corresponding
// to HTTP responses, on the same connection. The HTTP/2 setting
// parameter SETTINGS_HEADER_TABLE_SIZE is set to the value of 256
// octets, causing some evictions to occur."
func TestDecodeC5_ResponsesNoHuff(t *testing.T) ***REMOVED***
	testDecodeSeries(t, 256, []encAndWant***REMOVED***
		***REMOVED***dehex(`
4803 3330 3258 0770 7269 7661 7465 611d
4d6f 6e2c 2032 3120 4f63 7420 3230 3133
2032 303a 3133 3a32 3120 474d 546e 1768
7474 7073 3a2f 2f77 7777 2e65 7861 6d70
6c65 2e63 6f6d
`),
			[]HeaderField***REMOVED***
				pair(":status", "302"),
				pair("cache-control", "private"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("location", "https://www.example.com"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("location", "https://www.example.com"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("cache-control", "private"),
				pair(":status", "302"),
			***REMOVED***,
			222,
		***REMOVED***,
		***REMOVED***dehex("4803 3330 37c1 c0bf"),
			[]HeaderField***REMOVED***
				pair(":status", "307"),
				pair("cache-control", "private"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("location", "https://www.example.com"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair(":status", "307"),
				pair("location", "https://www.example.com"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("cache-control", "private"),
			***REMOVED***,
			222,
		***REMOVED***,
		***REMOVED***dehex(`
88c1 611d 4d6f 6e2c 2032 3120 4f63 7420
3230 3133 2032 303a 3133 3a32 3220 474d
54c0 5a04 677a 6970 7738 666f 6f3d 4153
444a 4b48 514b 425a 584f 5157 454f 5049
5541 5851 5745 4f49 553b 206d 6178 2d61
6765 3d33 3630 303b 2076 6572 7369 6f6e
3d31
`),
			[]HeaderField***REMOVED***
				pair(":status", "200"),
				pair("cache-control", "private"),
				pair("date", "Mon, 21 Oct 2013 20:13:22 GMT"),
				pair("location", "https://www.example.com"),
				pair("content-encoding", "gzip"),
				pair("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"),
				pair("content-encoding", "gzip"),
				pair("date", "Mon, 21 Oct 2013 20:13:22 GMT"),
			***REMOVED***,
			215,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// http://http2.github.io/http2-spec/compression.html#rfc.section.C.6
// "This section shows the same examples as the previous section, but
// using Huffman encoding for the literal values. The HTTP/2 setting
// parameter SETTINGS_HEADER_TABLE_SIZE is set to the value of 256
// octets, causing some evictions to occur. The eviction mechanism
// uses the length of the decoded literal values, so the same
// evictions occurs as in the previous section."
func TestDecodeC6_ResponsesHuffman(t *testing.T) ***REMOVED***
	testDecodeSeries(t, 256, []encAndWant***REMOVED***
		***REMOVED***dehex(`
4882 6402 5885 aec3 771a 4b61 96d0 7abe
9410 54d4 44a8 2005 9504 0b81 66e0 82a6
2d1b ff6e 919d 29ad 1718 63c7 8f0b 97c8
e9ae 82ae 43d3
`),
			[]HeaderField***REMOVED***
				pair(":status", "302"),
				pair("cache-control", "private"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("location", "https://www.example.com"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("location", "https://www.example.com"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("cache-control", "private"),
				pair(":status", "302"),
			***REMOVED***,
			222,
		***REMOVED***,
		***REMOVED***dehex("4883 640e ffc1 c0bf"),
			[]HeaderField***REMOVED***
				pair(":status", "307"),
				pair("cache-control", "private"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("location", "https://www.example.com"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair(":status", "307"),
				pair("location", "https://www.example.com"),
				pair("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
				pair("cache-control", "private"),
			***REMOVED***,
			222,
		***REMOVED***,
		***REMOVED***dehex(`
88c1 6196 d07a be94 1054 d444 a820 0595
040b 8166 e084 a62d 1bff c05a 839b d9ab
77ad 94e7 821d d7f2 e6c7 b335 dfdf cd5b
3960 d5af 2708 7f36 72c1 ab27 0fb5 291f
9587 3160 65c0 03ed 4ee5 b106 3d50 07
`),
			[]HeaderField***REMOVED***
				pair(":status", "200"),
				pair("cache-control", "private"),
				pair("date", "Mon, 21 Oct 2013 20:13:22 GMT"),
				pair("location", "https://www.example.com"),
				pair("content-encoding", "gzip"),
				pair("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"),
			***REMOVED***,
			[]HeaderField***REMOVED***
				pair("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"),
				pair("content-encoding", "gzip"),
				pair("date", "Mon, 21 Oct 2013 20:13:22 GMT"),
			***REMOVED***,
			215,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func testDecodeSeries(t *testing.T, size uint32, steps []encAndWant) ***REMOVED***
	d := NewDecoder(size, nil)
	for i, step := range steps ***REMOVED***
		hf, err := d.DecodeFull(step.enc)
		if err != nil ***REMOVED***
			t.Fatalf("Error at step index %d: %v", i, err)
		***REMOVED***
		if !reflect.DeepEqual(hf, step.want) ***REMOVED***
			t.Fatalf("At step index %d: Got headers %v; want %v", i, hf, step.want)
		***REMOVED***
		gotDynTab := d.dynTab.reverseCopy()
		if !reflect.DeepEqual(gotDynTab, step.wantDynTab) ***REMOVED***
			t.Errorf("After step index %d, dynamic table = %v; want %v", i, gotDynTab, step.wantDynTab)
		***REMOVED***
		if d.dynTab.size != step.wantDynSize ***REMOVED***
			t.Errorf("After step index %d, dynamic table size = %v; want %v", i, d.dynTab.size, step.wantDynSize)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHuffmanDecodeExcessPadding(t *testing.T) ***REMOVED***
	tests := [][]byte***REMOVED***
		***REMOVED***0xff***REMOVED***,                                   // Padding Exceeds 7 bits
		***REMOVED***0x1f, 0xff***REMOVED***,                             // ***REMOVED***"a", 1 byte excess padding***REMOVED***
		***REMOVED***0x1f, 0xff, 0xff***REMOVED***,                       // ***REMOVED***"a", 2 byte excess padding***REMOVED***
		***REMOVED***0x1f, 0xff, 0xff, 0xff***REMOVED***,                 // ***REMOVED***"a", 3 byte excess padding***REMOVED***
		***REMOVED***0xff, 0x9f, 0xff, 0xff, 0xff***REMOVED***,           // ***REMOVED***"a", 29 bit excess padding***REMOVED***
		***REMOVED***'R', 0xbc, '0', 0xff, 0xff, 0xff, 0xff***REMOVED***, // Padding ends on partial symbol.
	***REMOVED***
	for i, in := range tests ***REMOVED***
		var buf bytes.Buffer
		if _, err := HuffmanDecode(&buf, in); err != ErrInvalidHuffman ***REMOVED***
			t.Errorf("test-%d: decode(%q) = %v; want ErrInvalidHuffman", i, in, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHuffmanDecodeEOS(t *testing.T) ***REMOVED***
	in := []byte***REMOVED***0xff, 0xff, 0xff, 0xff, 0xfc***REMOVED*** // ***REMOVED***EOS, "?"***REMOVED***
	var buf bytes.Buffer
	if _, err := HuffmanDecode(&buf, in); err != ErrInvalidHuffman ***REMOVED***
		t.Errorf("error = %v; want ErrInvalidHuffman", err)
	***REMOVED***
***REMOVED***

func TestHuffmanDecodeMaxLengthOnTrailingByte(t *testing.T) ***REMOVED***
	in := []byte***REMOVED***0x00, 0x01***REMOVED*** // ***REMOVED***"0", "0", "0"***REMOVED***
	var buf bytes.Buffer
	if err := huffmanDecode(&buf, 2, in); err != ErrStringLength ***REMOVED***
		t.Errorf("error = %v; want ErrStringLength", err)
	***REMOVED***
***REMOVED***

func TestHuffmanDecodeCorruptPadding(t *testing.T) ***REMOVED***
	in := []byte***REMOVED***0x00***REMOVED***
	var buf bytes.Buffer
	if _, err := HuffmanDecode(&buf, in); err != ErrInvalidHuffman ***REMOVED***
		t.Errorf("error = %v; want ErrInvalidHuffman", err)
	***REMOVED***
***REMOVED***

func TestHuffmanDecode(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		inHex, want string
	***REMOVED******REMOVED***
		***REMOVED***"f1e3 c2e5 f23a 6ba0 ab90 f4ff", "www.example.com"***REMOVED***,
		***REMOVED***"a8eb 1064 9cbf", "no-cache"***REMOVED***,
		***REMOVED***"25a8 49e9 5ba9 7d7f", "custom-key"***REMOVED***,
		***REMOVED***"25a8 49e9 5bb8 e8b4 bf", "custom-value"***REMOVED***,
		***REMOVED***"6402", "302"***REMOVED***,
		***REMOVED***"aec3 771a 4b", "private"***REMOVED***,
		***REMOVED***"d07a be94 1054 d444 a820 0595 040b 8166 e082 a62d 1bff", "Mon, 21 Oct 2013 20:13:21 GMT"***REMOVED***,
		***REMOVED***"9d29 ad17 1863 c78f 0b97 c8e9 ae82 ae43 d3", "https://www.example.com"***REMOVED***,
		***REMOVED***"9bd9 ab", "gzip"***REMOVED***,
		***REMOVED***"94e7 821d d7f2 e6c7 b335 dfdf cd5b 3960 d5af 2708 7f36 72c1 ab27 0fb5 291f 9587 3160 65c0 03ed 4ee5 b106 3d50 07",
			"foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		var buf bytes.Buffer
		in, err := hex.DecodeString(strings.Replace(tt.inHex, " ", "", -1))
		if err != nil ***REMOVED***
			t.Errorf("%d. hex input error: %v", i, err)
			continue
		***REMOVED***
		if _, err := HuffmanDecode(&buf, in); err != nil ***REMOVED***
			t.Errorf("%d. decode error: %v", i, err)
			continue
		***REMOVED***
		if got := buf.String(); tt.want != got ***REMOVED***
			t.Errorf("%d. decode = %q; want %q", i, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAppendHuffmanString(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in, want string
	***REMOVED******REMOVED***
		***REMOVED***"www.example.com", "f1e3 c2e5 f23a 6ba0 ab90 f4ff"***REMOVED***,
		***REMOVED***"no-cache", "a8eb 1064 9cbf"***REMOVED***,
		***REMOVED***"custom-key", "25a8 49e9 5ba9 7d7f"***REMOVED***,
		***REMOVED***"custom-value", "25a8 49e9 5bb8 e8b4 bf"***REMOVED***,
		***REMOVED***"302", "6402"***REMOVED***,
		***REMOVED***"private", "aec3 771a 4b"***REMOVED***,
		***REMOVED***"Mon, 21 Oct 2013 20:13:21 GMT", "d07a be94 1054 d444 a820 0595 040b 8166 e082 a62d 1bff"***REMOVED***,
		***REMOVED***"https://www.example.com", "9d29 ad17 1863 c78f 0b97 c8e9 ae82 ae43 d3"***REMOVED***,
		***REMOVED***"gzip", "9bd9 ab"***REMOVED***,
		***REMOVED***"foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1",
			"94e7 821d d7f2 e6c7 b335 dfdf cd5b 3960 d5af 2708 7f36 72c1 ab27 0fb5 291f 9587 3160 65c0 03ed 4ee5 b106 3d50 07"***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		buf := []byte***REMOVED******REMOVED***
		want := strings.Replace(tt.want, " ", "", -1)
		buf = AppendHuffmanString(buf, tt.in)
		if got := hex.EncodeToString(buf); want != got ***REMOVED***
			t.Errorf("%d. encode = %q; want %q", i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHuffmanMaxStrLen(t *testing.T) ***REMOVED***
	const msg = "Some string"
	huff := AppendHuffmanString(nil, msg)

	testGood := func(max int) ***REMOVED***
		var out bytes.Buffer
		if err := huffmanDecode(&out, max, huff); err != nil ***REMOVED***
			t.Errorf("For maxLen=%d, unexpected error: %v", max, err)
		***REMOVED***
		if out.String() != msg ***REMOVED***
			t.Errorf("For maxLen=%d, out = %q; want %q", max, out.String(), msg)
		***REMOVED***
	***REMOVED***
	testGood(0)
	testGood(len(msg))
	testGood(len(msg) + 1)

	var out bytes.Buffer
	if err := huffmanDecode(&out, len(msg)-1, huff); err != ErrStringLength ***REMOVED***
		t.Errorf("err = %v; want ErrStringLength", err)
	***REMOVED***
***REMOVED***

func TestHuffmanRoundtripStress(t *testing.T) ***REMOVED***
	const Len = 50 // of uncompressed string
	input := make([]byte, Len)
	var output bytes.Buffer
	var huff []byte

	n := 5000
	if testing.Short() ***REMOVED***
		n = 100
	***REMOVED***
	seed := time.Now().UnixNano()
	t.Logf("Seed = %v", seed)
	src := rand.New(rand.NewSource(seed))
	var encSize int64
	for i := 0; i < n; i++ ***REMOVED***
		for l := range input ***REMOVED***
			input[l] = byte(src.Intn(256))
		***REMOVED***
		huff = AppendHuffmanString(huff[:0], string(input))
		encSize += int64(len(huff))
		output.Reset()
		if err := huffmanDecode(&output, 0, huff); err != nil ***REMOVED***
			t.Errorf("Failed to decode %q -> %q -> error %v", input, huff, err)
			continue
		***REMOVED***
		if !bytes.Equal(output.Bytes(), input) ***REMOVED***
			t.Errorf("Roundtrip failure on %q -> %q -> %q", input, huff, output.Bytes())
		***REMOVED***
	***REMOVED***
	t.Logf("Compressed size of original: %0.02f%% (%v -> %v)", 100*(float64(encSize)/(Len*float64(n))), Len*n, encSize)
***REMOVED***

func TestHuffmanDecodeFuzz(t *testing.T) ***REMOVED***
	const Len = 50 // of compressed
	var buf, zbuf bytes.Buffer

	n := 5000
	if testing.Short() ***REMOVED***
		n = 100
	***REMOVED***
	seed := time.Now().UnixNano()
	t.Logf("Seed = %v", seed)
	src := rand.New(rand.NewSource(seed))
	numFail := 0
	for i := 0; i < n; i++ ***REMOVED***
		zbuf.Reset()
		if i == 0 ***REMOVED***
			// Start with at least one invalid one.
			zbuf.WriteString("00\x91\xff\xff\xff\xff\xc8")
		***REMOVED*** else ***REMOVED***
			for l := 0; l < Len; l++ ***REMOVED***
				zbuf.WriteByte(byte(src.Intn(256)))
			***REMOVED***
		***REMOVED***

		buf.Reset()
		if err := huffmanDecode(&buf, 0, zbuf.Bytes()); err != nil ***REMOVED***
			if err == ErrInvalidHuffman ***REMOVED***
				numFail++
				continue
			***REMOVED***
			t.Errorf("Failed to decode %q: %v", zbuf.Bytes(), err)
			continue
		***REMOVED***
	***REMOVED***
	t.Logf("%0.02f%% are invalid (%d / %d)", 100*float64(numFail)/float64(n), numFail, n)
	if numFail < 1 ***REMOVED***
		t.Error("expected at least one invalid huffman encoding (test starts with one)")
	***REMOVED***
***REMOVED***

func TestReadVarInt(t *testing.T) ***REMOVED***
	type res struct ***REMOVED***
		i        uint64
		consumed int
		err      error
	***REMOVED***
	tests := []struct ***REMOVED***
		n    byte
		p    []byte
		want res
	***REMOVED******REMOVED***
		// Fits in a byte:
		***REMOVED***1, []byte***REMOVED***0***REMOVED***, res***REMOVED***0, 1, nil***REMOVED******REMOVED***,
		***REMOVED***2, []byte***REMOVED***2***REMOVED***, res***REMOVED***2, 1, nil***REMOVED******REMOVED***,
		***REMOVED***3, []byte***REMOVED***6***REMOVED***, res***REMOVED***6, 1, nil***REMOVED******REMOVED***,
		***REMOVED***4, []byte***REMOVED***14***REMOVED***, res***REMOVED***14, 1, nil***REMOVED******REMOVED***,
		***REMOVED***5, []byte***REMOVED***30***REMOVED***, res***REMOVED***30, 1, nil***REMOVED******REMOVED***,
		***REMOVED***6, []byte***REMOVED***62***REMOVED***, res***REMOVED***62, 1, nil***REMOVED******REMOVED***,
		***REMOVED***7, []byte***REMOVED***126***REMOVED***, res***REMOVED***126, 1, nil***REMOVED******REMOVED***,
		***REMOVED***8, []byte***REMOVED***254***REMOVED***, res***REMOVED***254, 1, nil***REMOVED******REMOVED***,

		// Doesn't fit in a byte:
		***REMOVED***1, []byte***REMOVED***1***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,
		***REMOVED***2, []byte***REMOVED***3***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,
		***REMOVED***3, []byte***REMOVED***7***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,
		***REMOVED***4, []byte***REMOVED***15***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,
		***REMOVED***5, []byte***REMOVED***31***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,
		***REMOVED***6, []byte***REMOVED***63***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,
		***REMOVED***7, []byte***REMOVED***127***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,
		***REMOVED***8, []byte***REMOVED***255***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,

		// Ignoring top bits:
		***REMOVED***5, []byte***REMOVED***255, 154, 10***REMOVED***, res***REMOVED***1337, 3, nil***REMOVED******REMOVED***, // high dummy three bits: 111
		***REMOVED***5, []byte***REMOVED***159, 154, 10***REMOVED***, res***REMOVED***1337, 3, nil***REMOVED******REMOVED***, // high dummy three bits: 100
		***REMOVED***5, []byte***REMOVED***191, 154, 10***REMOVED***, res***REMOVED***1337, 3, nil***REMOVED******REMOVED***, // high dummy three bits: 101

		// Extra byte:
		***REMOVED***5, []byte***REMOVED***191, 154, 10, 2***REMOVED***, res***REMOVED***1337, 3, nil***REMOVED******REMOVED***, // extra byte

		// Short a byte:
		***REMOVED***5, []byte***REMOVED***191, 154***REMOVED***, res***REMOVED***0, 0, errNeedMore***REMOVED******REMOVED***,

		// integer overflow:
		***REMOVED***1, []byte***REMOVED***255, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128***REMOVED***, res***REMOVED***0, 0, errVarintOverflow***REMOVED******REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		i, remain, err := readVarInt(tt.n, tt.p)
		consumed := len(tt.p) - len(remain)
		got := res***REMOVED***i, consumed, err***REMOVED***
		if got != tt.want ***REMOVED***
			t.Errorf("readVarInt(%d, %v ~ %x) = %+v; want %+v", tt.n, tt.p, tt.p, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Fuzz crash, originally reported at https://github.com/bradfitz/http2/issues/56
func TestHuffmanFuzzCrash(t *testing.T) ***REMOVED***
	got, err := HuffmanDecodeToString([]byte("00\x91\xff\xff\xff\xff\xc8"))
	if got != "" ***REMOVED***
		t.Errorf("Got %q; want empty string", got)
	***REMOVED***
	if err != ErrInvalidHuffman ***REMOVED***
		t.Errorf("Err = %v; want ErrInvalidHuffman", err)
	***REMOVED***
***REMOVED***

func pair(name, value string) HeaderField ***REMOVED***
	return HeaderField***REMOVED***Name: name, Value: value***REMOVED***
***REMOVED***

func dehex(s string) []byte ***REMOVED***
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	b, err := hex.DecodeString(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return b
***REMOVED***

func TestEmitEnabled(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.WriteField(HeaderField***REMOVED***Name: "foo", Value: "bar"***REMOVED***)
	enc.WriteField(HeaderField***REMOVED***Name: "foo", Value: "bar"***REMOVED***)

	numCallback := 0
	var dec *Decoder
	dec = NewDecoder(8<<20, func(HeaderField) ***REMOVED***
		numCallback++
		dec.SetEmitEnabled(false)
	***REMOVED***)
	if !dec.EmitEnabled() ***REMOVED***
		t.Errorf("initial emit enabled = false; want true")
	***REMOVED***
	if _, err := dec.Write(buf.Bytes()); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	if numCallback != 1 ***REMOVED***
		t.Errorf("num callbacks = %d; want 1", numCallback)
	***REMOVED***
	if dec.EmitEnabled() ***REMOVED***
		t.Errorf("emit enabled = true; want false")
	***REMOVED***
***REMOVED***

func TestSaveBufLimit(t *testing.T) ***REMOVED***
	const maxStr = 1 << 10
	var got []HeaderField
	dec := NewDecoder(initialHeaderTableSize, func(hf HeaderField) ***REMOVED***
		got = append(got, hf)
	***REMOVED***)
	dec.SetMaxStringLength(maxStr)
	var frag []byte
	frag = append(frag[:0], encodeTypeByte(false, false))
	frag = appendVarInt(frag, 7, 3)
	frag = append(frag, "foo"...)
	frag = appendVarInt(frag, 7, 3)
	frag = append(frag, "bar"...)

	if _, err := dec.Write(frag); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	want := []HeaderField***REMOVED******REMOVED***Name: "foo", Value: "bar"***REMOVED******REMOVED***
	if !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("After small writes, got %v; want %v", got, want)
	***REMOVED***

	frag = append(frag[:0], encodeTypeByte(false, false))
	frag = appendVarInt(frag, 7, maxStr*3)
	frag = append(frag, make([]byte, maxStr*3)...)

	_, err := dec.Write(frag)
	if err != ErrStringLength ***REMOVED***
		t.Fatalf("Write error = %v; want ErrStringLength", err)
	***REMOVED***
***REMOVED***
