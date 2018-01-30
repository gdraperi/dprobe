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
)

func TestEncoderTableSizeUpdate(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		size1, size2 uint32
		wantHex      string
	***REMOVED******REMOVED***
		// Should emit 2 table size updates (2048 and 4096)
		***REMOVED***2048, 4096, "3fe10f 3fe11f 82"***REMOVED***,

		// Should emit 1 table size update (2048)
		***REMOVED***16384, 2048, "3fe10f 82"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		var buf bytes.Buffer
		e := NewEncoder(&buf)
		e.SetMaxDynamicTableSize(tt.size1)
		e.SetMaxDynamicTableSize(tt.size2)
		if err := e.WriteField(pair(":method", "GET")); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		want := removeSpace(tt.wantHex)
		if got := hex.EncodeToString(buf.Bytes()); got != want ***REMOVED***
			t.Errorf("e.SetDynamicTableSize %v, %v = %q; want %q", tt.size1, tt.size2, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEncoderWriteField(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	var got []HeaderField
	d := NewDecoder(4<<10, func(f HeaderField) ***REMOVED***
		got = append(got, f)
	***REMOVED***)

	tests := []struct ***REMOVED***
		hdrs []HeaderField
	***REMOVED******REMOVED***
		***REMOVED***[]HeaderField***REMOVED***
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", "www.example.com"),
		***REMOVED******REMOVED***,
		***REMOVED***[]HeaderField***REMOVED***
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", "www.example.com"),
			pair("cache-control", "no-cache"),
		***REMOVED******REMOVED***,
		***REMOVED***[]HeaderField***REMOVED***
			pair(":method", "GET"),
			pair(":scheme", "https"),
			pair(":path", "/index.html"),
			pair(":authority", "www.example.com"),
			pair("custom-key", "custom-value"),
		***REMOVED******REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		buf.Reset()
		got = got[:0]
		for _, hf := range tt.hdrs ***REMOVED***
			if err := e.WriteField(hf); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
		_, err := d.Write(buf.Bytes())
		if err != nil ***REMOVED***
			t.Errorf("%d. Decoder Write = %v", i, err)
		***REMOVED***
		if !reflect.DeepEqual(got, tt.hdrs) ***REMOVED***
			t.Errorf("%d. Decoded %+v; want %+v", i, got, tt.hdrs)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEncoderSearchTable(t *testing.T) ***REMOVED***
	e := NewEncoder(nil)

	e.dynTab.add(pair("foo", "bar"))
	e.dynTab.add(pair("blake", "miz"))
	e.dynTab.add(pair(":method", "GET"))

	tests := []struct ***REMOVED***
		hf        HeaderField
		wantI     uint64
		wantMatch bool
	***REMOVED******REMOVED***
		// Name and Value match
		***REMOVED***pair("foo", "bar"), uint64(staticTable.len()) + 3, true***REMOVED***,
		***REMOVED***pair("blake", "miz"), uint64(staticTable.len()) + 2, true***REMOVED***,
		***REMOVED***pair(":method", "GET"), 2, true***REMOVED***,

		// Only name match because Sensitive == true. This is allowed to match
		// any ":method" entry. The current implementation uses the last entry
		// added in newStaticTable.
		***REMOVED***HeaderField***REMOVED***":method", "GET", true***REMOVED***, 3, false***REMOVED***,

		// Only Name matches
		***REMOVED***pair("foo", "..."), uint64(staticTable.len()) + 3, false***REMOVED***,
		***REMOVED***pair("blake", "..."), uint64(staticTable.len()) + 2, false***REMOVED***,
		// As before, this is allowed to match any ":method" entry.
		***REMOVED***pair(":method", "..."), 3, false***REMOVED***,

		// None match
		***REMOVED***pair("foo-", "bar"), 0, false***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		if gotI, gotMatch := e.searchTable(tt.hf); gotI != tt.wantI || gotMatch != tt.wantMatch ***REMOVED***
			t.Errorf("d.search(%+v) = %v, %v; want %v, %v", tt.hf, gotI, gotMatch, tt.wantI, tt.wantMatch)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAppendVarInt(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		n    byte
		i    uint64
		want []byte
	***REMOVED******REMOVED***
		// Fits in a byte:
		***REMOVED***1, 0, []byte***REMOVED***0***REMOVED******REMOVED***,
		***REMOVED***2, 2, []byte***REMOVED***2***REMOVED******REMOVED***,
		***REMOVED***3, 6, []byte***REMOVED***6***REMOVED******REMOVED***,
		***REMOVED***4, 14, []byte***REMOVED***14***REMOVED******REMOVED***,
		***REMOVED***5, 30, []byte***REMOVED***30***REMOVED******REMOVED***,
		***REMOVED***6, 62, []byte***REMOVED***62***REMOVED******REMOVED***,
		***REMOVED***7, 126, []byte***REMOVED***126***REMOVED******REMOVED***,
		***REMOVED***8, 254, []byte***REMOVED***254***REMOVED******REMOVED***,

		// Multiple bytes:
		***REMOVED***5, 1337, []byte***REMOVED***31, 154, 10***REMOVED******REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		got := appendVarInt(nil, tt.n, tt.i)
		if !bytes.Equal(got, tt.want) ***REMOVED***
			t.Errorf("appendVarInt(nil, %v, %v) = %v; want %v", tt.n, tt.i, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAppendHpackString(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		s, wantHex string
	***REMOVED******REMOVED***
		// Huffman encoded
		***REMOVED***"www.example.com", "8c f1e3 c2e5 f23a 6ba0 ab90 f4ff"***REMOVED***,

		// Not Huffman encoded
		***REMOVED***"a", "01 61"***REMOVED***,

		// zero length
		***REMOVED***"", "00"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		want := removeSpace(tt.wantHex)
		buf := appendHpackString(nil, tt.s)
		if got := hex.EncodeToString(buf); want != got ***REMOVED***
			t.Errorf("appendHpackString(nil, %q) = %q; want %q", tt.s, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAppendIndexed(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		i       uint64
		wantHex string
	***REMOVED******REMOVED***
		// 1 byte
		***REMOVED***1, "81"***REMOVED***,
		***REMOVED***126, "fe"***REMOVED***,

		// 2 bytes
		***REMOVED***127, "ff00"***REMOVED***,
		***REMOVED***128, "ff01"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		want := removeSpace(tt.wantHex)
		buf := appendIndexed(nil, tt.i)
		if got := hex.EncodeToString(buf); want != got ***REMOVED***
			t.Errorf("appendIndex(nil, %v) = %q; want %q", tt.i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAppendNewName(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		f        HeaderField
		indexing bool
		wantHex  string
	***REMOVED******REMOVED***
		// Incremental indexing
		***REMOVED***HeaderField***REMOVED***"custom-key", "custom-value", false***REMOVED***, true, "40 88 25a8 49e9 5ba9 7d7f 89 25a8 49e9 5bb8 e8b4 bf"***REMOVED***,

		// Without indexing
		***REMOVED***HeaderField***REMOVED***"custom-key", "custom-value", false***REMOVED***, false, "00 88 25a8 49e9 5ba9 7d7f 89 25a8 49e9 5bb8 e8b4 bf"***REMOVED***,

		// Never indexed
		***REMOVED***HeaderField***REMOVED***"custom-key", "custom-value", true***REMOVED***, true, "10 88 25a8 49e9 5ba9 7d7f 89 25a8 49e9 5bb8 e8b4 bf"***REMOVED***,
		***REMOVED***HeaderField***REMOVED***"custom-key", "custom-value", true***REMOVED***, false, "10 88 25a8 49e9 5ba9 7d7f 89 25a8 49e9 5bb8 e8b4 bf"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		want := removeSpace(tt.wantHex)
		buf := appendNewName(nil, tt.f, tt.indexing)
		if got := hex.EncodeToString(buf); want != got ***REMOVED***
			t.Errorf("appendNewName(nil, %+v, %v) = %q; want %q", tt.f, tt.indexing, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAppendIndexedName(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		f        HeaderField
		i        uint64
		indexing bool
		wantHex  string
	***REMOVED******REMOVED***
		// Incremental indexing
		***REMOVED***HeaderField***REMOVED***":status", "302", false***REMOVED***, 8, true, "48 82 6402"***REMOVED***,

		// Without indexing
		***REMOVED***HeaderField***REMOVED***":status", "302", false***REMOVED***, 8, false, "08 82 6402"***REMOVED***,

		// Never indexed
		***REMOVED***HeaderField***REMOVED***":status", "302", true***REMOVED***, 8, true, "18 82 6402"***REMOVED***,
		***REMOVED***HeaderField***REMOVED***":status", "302", true***REMOVED***, 8, false, "18 82 6402"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		want := removeSpace(tt.wantHex)
		buf := appendIndexedName(nil, tt.f, tt.i, tt.indexing)
		if got := hex.EncodeToString(buf); want != got ***REMOVED***
			t.Errorf("appendIndexedName(nil, %+v, %v) = %q; want %q", tt.f, tt.indexing, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAppendTableSize(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		i       uint32
		wantHex string
	***REMOVED******REMOVED***
		// Fits into 1 byte
		***REMOVED***30, "3e"***REMOVED***,

		// Extra byte
		***REMOVED***31, "3f00"***REMOVED***,
		***REMOVED***32, "3f01"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		want := removeSpace(tt.wantHex)
		buf := appendTableSize(nil, tt.i)
		if got := hex.EncodeToString(buf); want != got ***REMOVED***
			t.Errorf("appendTableSize(nil, %v) = %q; want %q", tt.i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEncoderSetMaxDynamicTableSize(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	tests := []struct ***REMOVED***
		v           uint32
		wantUpdate  bool
		wantMinSize uint32
		wantMaxSize uint32
	***REMOVED******REMOVED***
		// Set new table size to 2048
		***REMOVED***2048, true, 2048, 2048***REMOVED***,

		// Set new table size to 16384, but still limited to
		// 4096
		***REMOVED***16384, true, 2048, 4096***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		e.SetMaxDynamicTableSize(tt.v)
		if got := e.tableSizeUpdate; tt.wantUpdate != got ***REMOVED***
			t.Errorf("e.tableSizeUpdate = %v; want %v", got, tt.wantUpdate)
		***REMOVED***
		if got := e.minSize; tt.wantMinSize != got ***REMOVED***
			t.Errorf("e.minSize = %v; want %v", got, tt.wantMinSize)
		***REMOVED***
		if got := e.dynTab.maxSize; tt.wantMaxSize != got ***REMOVED***
			t.Errorf("e.maxSize = %v; want %v", got, tt.wantMaxSize)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEncoderSetMaxDynamicTableSizeLimit(t *testing.T) ***REMOVED***
	e := NewEncoder(nil)
	// 4095 < initialHeaderTableSize means maxSize is truncated to
	// 4095.
	e.SetMaxDynamicTableSizeLimit(4095)
	if got, want := e.dynTab.maxSize, uint32(4095); got != want ***REMOVED***
		t.Errorf("e.dynTab.maxSize = %v; want %v", got, want)
	***REMOVED***
	if got, want := e.maxSizeLimit, uint32(4095); got != want ***REMOVED***
		t.Errorf("e.maxSizeLimit = %v; want %v", got, want)
	***REMOVED***
	if got, want := e.tableSizeUpdate, true; got != want ***REMOVED***
		t.Errorf("e.tableSizeUpdate = %v; want %v", got, want)
	***REMOVED***
	// maxSize will be truncated to maxSizeLimit
	e.SetMaxDynamicTableSize(16384)
	if got, want := e.dynTab.maxSize, uint32(4095); got != want ***REMOVED***
		t.Errorf("e.dynTab.maxSize = %v; want %v", got, want)
	***REMOVED***
	// 8192 > current maxSizeLimit, so maxSize does not change.
	e.SetMaxDynamicTableSizeLimit(8192)
	if got, want := e.dynTab.maxSize, uint32(4095); got != want ***REMOVED***
		t.Errorf("e.dynTab.maxSize = %v; want %v", got, want)
	***REMOVED***
	if got, want := e.maxSizeLimit, uint32(8192); got != want ***REMOVED***
		t.Errorf("e.maxSizeLimit = %v; want %v", got, want)
	***REMOVED***
***REMOVED***

func removeSpace(s string) string ***REMOVED***
	return strings.Replace(s, " ", "", -1)
***REMOVED***

func BenchmarkEncoderSearchTable(b *testing.B) ***REMOVED***
	e := NewEncoder(nil)

	// A sample of possible header fields.
	// This is not based on any actual data from HTTP/2 traces.
	var possible []HeaderField
	for _, f := range staticTable.ents ***REMOVED***
		if f.Value == "" ***REMOVED***
			possible = append(possible, f)
			continue
		***REMOVED***
		// Generate 5 random values, except for cookie and set-cookie,
		// which we know can have many values in practice.
		num := 5
		if f.Name == "cookie" || f.Name == "set-cookie" ***REMOVED***
			num = 25
		***REMOVED***
		for i := 0; i < num; i++ ***REMOVED***
			f.Value = fmt.Sprintf("%s-%d", f.Name, i)
			possible = append(possible, f)
		***REMOVED***
	***REMOVED***
	for k := 0; k < 10; k++ ***REMOVED***
		f := HeaderField***REMOVED***
			Name:      fmt.Sprintf("x-header-%d", k),
			Sensitive: rand.Int()%2 == 0,
		***REMOVED***
		for i := 0; i < 5; i++ ***REMOVED***
			f.Value = fmt.Sprintf("%s-%d", f.Name, i)
			possible = append(possible, f)
		***REMOVED***
	***REMOVED***

	// Add a random sample to the dynamic table. This very loosely simulates
	// a history of 100 requests with 20 header fields per request.
	for r := 0; r < 100*20; r++ ***REMOVED***
		f := possible[rand.Int31n(int32(len(possible)))]
		// Skip if this is in the staticTable verbatim.
		if _, has := staticTable.search(f); !has ***REMOVED***
			e.dynTab.add(f)
		***REMOVED***
	***REMOVED***

	b.ResetTimer()
	for n := 0; n < b.N; n++ ***REMOVED***
		for _, f := range possible ***REMOVED***
			e.searchTable(f)
		***REMOVED***
	***REMOVED***
***REMOVED***
