// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"golang.org/x/net/http2/hpack"
)

func testFramer() (*Framer, *bytes.Buffer) ***REMOVED***
	buf := new(bytes.Buffer)
	return NewFramer(buf, buf), buf
***REMOVED***

func TestFrameSizes(t *testing.T) ***REMOVED***
	// Catch people rearranging the FrameHeader fields.
	if got, want := int(unsafe.Sizeof(FrameHeader***REMOVED******REMOVED***)), 12; got != want ***REMOVED***
		t.Errorf("FrameHeader size = %d; want %d", got, want)
	***REMOVED***
***REMOVED***

func TestFrameTypeString(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		ft   FrameType
		want string
	***REMOVED******REMOVED***
		***REMOVED***FrameData, "DATA"***REMOVED***,
		***REMOVED***FramePing, "PING"***REMOVED***,
		***REMOVED***FrameGoAway, "GOAWAY"***REMOVED***,
		***REMOVED***0xf, "UNKNOWN_FRAME_TYPE_15"***REMOVED***,
	***REMOVED***

	for i, tt := range tests ***REMOVED***
		got := tt.ft.String()
		if got != tt.want ***REMOVED***
			t.Errorf("%d. String(FrameType %d) = %q; want %q", i, int(tt.ft), got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWriteRST(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	var streamID uint32 = 1<<24 + 2<<16 + 3<<8 + 4
	var errCode uint32 = 7<<24 + 6<<16 + 5<<8 + 4
	fr.WriteRSTStream(streamID, ErrCode(errCode))
	const wantEnc = "\x00\x00\x04\x03\x00\x01\x02\x03\x04\x07\x06\x05\x04"
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	want := &RSTStreamFrame***REMOVED***
		FrameHeader: FrameHeader***REMOVED***
			valid:    true,
			Type:     0x3,
			Flags:    0x0,
			Length:   0x4,
			StreamID: 0x1020304,
		***REMOVED***,
		ErrCode: 0x7060504,
	***REMOVED***
	if !reflect.DeepEqual(f, want) ***REMOVED***
		t.Errorf("parsed back %#v; want %#v", f, want)
	***REMOVED***
***REMOVED***

func TestWriteData(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	var streamID uint32 = 1<<24 + 2<<16 + 3<<8 + 4
	data := []byte("ABC")
	fr.WriteData(streamID, true, data)
	const wantEnc = "\x00\x00\x03\x00\x01\x01\x02\x03\x04ABC"
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	df, ok := f.(*DataFrame)
	if !ok ***REMOVED***
		t.Fatalf("got %T; want *DataFrame", f)
	***REMOVED***
	if !bytes.Equal(df.Data(), data) ***REMOVED***
		t.Errorf("got %q; want %q", df.Data(), data)
	***REMOVED***
	if f.Header().Flags&1 == 0 ***REMOVED***
		t.Errorf("didn't see END_STREAM flag")
	***REMOVED***
***REMOVED***

func TestWriteDataPadded(t *testing.T) ***REMOVED***
	tests := [...]struct ***REMOVED***
		streamID   uint32
		endStream  bool
		data       []byte
		pad        []byte
		wantHeader FrameHeader
	***REMOVED******REMOVED***
		// Unpadded:
		0: ***REMOVED***
			streamID:  1,
			endStream: true,
			data:      []byte("foo"),
			pad:       nil,
			wantHeader: FrameHeader***REMOVED***
				Type:     FrameData,
				Flags:    FlagDataEndStream,
				Length:   3,
				StreamID: 1,
			***REMOVED***,
		***REMOVED***,

		// Padded bit set, but no padding:
		1: ***REMOVED***
			streamID:  1,
			endStream: true,
			data:      []byte("foo"),
			pad:       []byte***REMOVED******REMOVED***,
			wantHeader: FrameHeader***REMOVED***
				Type:     FrameData,
				Flags:    FlagDataEndStream | FlagDataPadded,
				Length:   4,
				StreamID: 1,
			***REMOVED***,
		***REMOVED***,

		// Padded bit set, with padding:
		2: ***REMOVED***
			streamID:  1,
			endStream: false,
			data:      []byte("foo"),
			pad:       []byte***REMOVED***0, 0, 0***REMOVED***,
			wantHeader: FrameHeader***REMOVED***
				Type:     FrameData,
				Flags:    FlagDataPadded,
				Length:   7,
				StreamID: 1,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		fr, _ := testFramer()
		fr.WriteDataPadded(tt.streamID, tt.endStream, tt.data, tt.pad)
		f, err := fr.ReadFrame()
		if err != nil ***REMOVED***
			t.Errorf("%d. ReadFrame: %v", i, err)
			continue
		***REMOVED***
		got := f.Header()
		tt.wantHeader.valid = true
		if got != tt.wantHeader ***REMOVED***
			t.Errorf("%d. read %+v; want %+v", i, got, tt.wantHeader)
			continue
		***REMOVED***
		df := f.(*DataFrame)
		if !bytes.Equal(df.Data(), tt.data) ***REMOVED***
			t.Errorf("%d. got %q; want %q", i, df.Data(), tt.data)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWriteHeaders(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		name      string
		p         HeadersFrameParam
		wantEnc   string
		wantFrame *HeadersFrame
	***REMOVED******REMOVED***
		***REMOVED***
			"basic",
			HeadersFrameParam***REMOVED***
				StreamID:      42,
				BlockFragment: []byte("abc"),
				Priority:      PriorityParam***REMOVED******REMOVED***,
			***REMOVED***,
			"\x00\x00\x03\x01\x00\x00\x00\x00*abc",
			&HeadersFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					valid:    true,
					StreamID: 42,
					Type:     FrameHeaders,
					Length:   uint32(len("abc")),
				***REMOVED***,
				Priority:      PriorityParam***REMOVED******REMOVED***,
				headerFragBuf: []byte("abc"),
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"basic + end flags",
			HeadersFrameParam***REMOVED***
				StreamID:      42,
				BlockFragment: []byte("abc"),
				EndStream:     true,
				EndHeaders:    true,
				Priority:      PriorityParam***REMOVED******REMOVED***,
			***REMOVED***,
			"\x00\x00\x03\x01\x05\x00\x00\x00*abc",
			&HeadersFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					valid:    true,
					StreamID: 42,
					Type:     FrameHeaders,
					Flags:    FlagHeadersEndStream | FlagHeadersEndHeaders,
					Length:   uint32(len("abc")),
				***REMOVED***,
				Priority:      PriorityParam***REMOVED******REMOVED***,
				headerFragBuf: []byte("abc"),
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"with padding",
			HeadersFrameParam***REMOVED***
				StreamID:      42,
				BlockFragment: []byte("abc"),
				EndStream:     true,
				EndHeaders:    true,
				PadLength:     5,
				Priority:      PriorityParam***REMOVED******REMOVED***,
			***REMOVED***,
			"\x00\x00\t\x01\r\x00\x00\x00*\x05abc\x00\x00\x00\x00\x00",
			&HeadersFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					valid:    true,
					StreamID: 42,
					Type:     FrameHeaders,
					Flags:    FlagHeadersEndStream | FlagHeadersEndHeaders | FlagHeadersPadded,
					Length:   uint32(1 + len("abc") + 5), // pad length + contents + padding
				***REMOVED***,
				Priority:      PriorityParam***REMOVED******REMOVED***,
				headerFragBuf: []byte("abc"),
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"with priority",
			HeadersFrameParam***REMOVED***
				StreamID:      42,
				BlockFragment: []byte("abc"),
				EndStream:     true,
				EndHeaders:    true,
				PadLength:     2,
				Priority: PriorityParam***REMOVED***
					StreamDep: 15,
					Exclusive: true,
					Weight:    127,
				***REMOVED***,
			***REMOVED***,
			"\x00\x00\v\x01-\x00\x00\x00*\x02\x80\x00\x00\x0f\u007fabc\x00\x00",
			&HeadersFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					valid:    true,
					StreamID: 42,
					Type:     FrameHeaders,
					Flags:    FlagHeadersEndStream | FlagHeadersEndHeaders | FlagHeadersPadded | FlagHeadersPriority,
					Length:   uint32(1 + 5 + len("abc") + 2), // pad length + priority + contents + padding
				***REMOVED***,
				Priority: PriorityParam***REMOVED***
					StreamDep: 15,
					Exclusive: true,
					Weight:    127,
				***REMOVED***,
				headerFragBuf: []byte("abc"),
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"with priority stream dep zero", // golang.org/issue/15444
			HeadersFrameParam***REMOVED***
				StreamID:      42,
				BlockFragment: []byte("abc"),
				EndStream:     true,
				EndHeaders:    true,
				PadLength:     2,
				Priority: PriorityParam***REMOVED***
					StreamDep: 0,
					Exclusive: true,
					Weight:    127,
				***REMOVED***,
			***REMOVED***,
			"\x00\x00\v\x01-\x00\x00\x00*\x02\x80\x00\x00\x00\u007fabc\x00\x00",
			&HeadersFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					valid:    true,
					StreamID: 42,
					Type:     FrameHeaders,
					Flags:    FlagHeadersEndStream | FlagHeadersEndHeaders | FlagHeadersPadded | FlagHeadersPriority,
					Length:   uint32(1 + 5 + len("abc") + 2), // pad length + priority + contents + padding
				***REMOVED***,
				Priority: PriorityParam***REMOVED***
					StreamDep: 0,
					Exclusive: true,
					Weight:    127,
				***REMOVED***,
				headerFragBuf: []byte("abc"),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		fr, buf := testFramer()
		if err := fr.WriteHeaders(tt.p); err != nil ***REMOVED***
			t.Errorf("test %q: %v", tt.name, err)
			continue
		***REMOVED***
		if buf.String() != tt.wantEnc ***REMOVED***
			t.Errorf("test %q: encoded %q; want %q", tt.name, buf.Bytes(), tt.wantEnc)
		***REMOVED***
		f, err := fr.ReadFrame()
		if err != nil ***REMOVED***
			t.Errorf("test %q: failed to read the frame back: %v", tt.name, err)
			continue
		***REMOVED***
		if !reflect.DeepEqual(f, tt.wantFrame) ***REMOVED***
			t.Errorf("test %q: mismatch.\n got: %#v\nwant: %#v\n", tt.name, f, tt.wantFrame)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWriteInvalidStreamDep(t *testing.T) ***REMOVED***
	fr, _ := testFramer()
	err := fr.WriteHeaders(HeadersFrameParam***REMOVED***
		StreamID: 42,
		Priority: PriorityParam***REMOVED***
			StreamDep: 1 << 31,
		***REMOVED***,
	***REMOVED***)
	if err != errDepStreamID ***REMOVED***
		t.Errorf("header error = %v; want %q", err, errDepStreamID)
	***REMOVED***

	err = fr.WritePriority(2, PriorityParam***REMOVED***StreamDep: 1 << 31***REMOVED***)
	if err != errDepStreamID ***REMOVED***
		t.Errorf("priority error = %v; want %q", err, errDepStreamID)
	***REMOVED***
***REMOVED***

func TestWriteContinuation(t *testing.T) ***REMOVED***
	const streamID = 42
	tests := []struct ***REMOVED***
		name string
		end  bool
		frag []byte

		wantFrame *ContinuationFrame
	***REMOVED******REMOVED***
		***REMOVED***
			"not end",
			false,
			[]byte("abc"),
			&ContinuationFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					valid:    true,
					StreamID: streamID,
					Type:     FrameContinuation,
					Length:   uint32(len("abc")),
				***REMOVED***,
				headerFragBuf: []byte("abc"),
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"end",
			true,
			[]byte("def"),
			&ContinuationFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					valid:    true,
					StreamID: streamID,
					Type:     FrameContinuation,
					Flags:    FlagContinuationEndHeaders,
					Length:   uint32(len("def")),
				***REMOVED***,
				headerFragBuf: []byte("def"),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		fr, _ := testFramer()
		if err := fr.WriteContinuation(streamID, tt.end, tt.frag); err != nil ***REMOVED***
			t.Errorf("test %q: %v", tt.name, err)
			continue
		***REMOVED***
		fr.AllowIllegalReads = true
		f, err := fr.ReadFrame()
		if err != nil ***REMOVED***
			t.Errorf("test %q: failed to read the frame back: %v", tt.name, err)
			continue
		***REMOVED***
		if !reflect.DeepEqual(f, tt.wantFrame) ***REMOVED***
			t.Errorf("test %q: mismatch.\n got: %#v\nwant: %#v\n", tt.name, f, tt.wantFrame)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWritePriority(t *testing.T) ***REMOVED***
	const streamID = 42
	tests := []struct ***REMOVED***
		name      string
		priority  PriorityParam
		wantFrame *PriorityFrame
	***REMOVED******REMOVED***
		***REMOVED***
			"not exclusive",
			PriorityParam***REMOVED***
				StreamDep: 2,
				Exclusive: false,
				Weight:    127,
			***REMOVED***,
			&PriorityFrame***REMOVED***
				FrameHeader***REMOVED***
					valid:    true,
					StreamID: streamID,
					Type:     FramePriority,
					Length:   5,
				***REMOVED***,
				PriorityParam***REMOVED***
					StreamDep: 2,
					Exclusive: false,
					Weight:    127,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			"exclusive",
			PriorityParam***REMOVED***
				StreamDep: 3,
				Exclusive: true,
				Weight:    77,
			***REMOVED***,
			&PriorityFrame***REMOVED***
				FrameHeader***REMOVED***
					valid:    true,
					StreamID: streamID,
					Type:     FramePriority,
					Length:   5,
				***REMOVED***,
				PriorityParam***REMOVED***
					StreamDep: 3,
					Exclusive: true,
					Weight:    77,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		fr, _ := testFramer()
		if err := fr.WritePriority(streamID, tt.priority); err != nil ***REMOVED***
			t.Errorf("test %q: %v", tt.name, err)
			continue
		***REMOVED***
		f, err := fr.ReadFrame()
		if err != nil ***REMOVED***
			t.Errorf("test %q: failed to read the frame back: %v", tt.name, err)
			continue
		***REMOVED***
		if !reflect.DeepEqual(f, tt.wantFrame) ***REMOVED***
			t.Errorf("test %q: mismatch.\n got: %#v\nwant: %#v\n", tt.name, f, tt.wantFrame)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWriteSettings(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	settings := []Setting***REMOVED******REMOVED***1, 2***REMOVED***, ***REMOVED***3, 4***REMOVED******REMOVED***
	fr.WriteSettings(settings...)
	const wantEnc = "\x00\x00\f\x04\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x02\x00\x03\x00\x00\x00\x04"
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	sf, ok := f.(*SettingsFrame)
	if !ok ***REMOVED***
		t.Fatalf("Got a %T; want a SettingsFrame", f)
	***REMOVED***
	var got []Setting
	sf.ForeachSetting(func(s Setting) error ***REMOVED***
		got = append(got, s)
		valBack, ok := sf.Value(s.ID)
		if !ok || valBack != s.Val ***REMOVED***
			t.Errorf("Value(%d) = %v, %v; want %v, true", s.ID, valBack, ok, s.Val)
		***REMOVED***
		return nil
	***REMOVED***)
	if !reflect.DeepEqual(settings, got) ***REMOVED***
		t.Errorf("Read settings %+v != written settings %+v", got, settings)
	***REMOVED***
***REMOVED***

func TestWriteSettingsAck(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	fr.WriteSettingsAck()
	const wantEnc = "\x00\x00\x00\x04\x01\x00\x00\x00\x00"
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
***REMOVED***

func TestWriteWindowUpdate(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	const streamID = 1<<24 + 2<<16 + 3<<8 + 4
	const incr = 7<<24 + 6<<16 + 5<<8 + 4
	if err := fr.WriteWindowUpdate(streamID, incr); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	const wantEnc = "\x00\x00\x04\x08\x00\x01\x02\x03\x04\x07\x06\x05\x04"
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	want := &WindowUpdateFrame***REMOVED***
		FrameHeader: FrameHeader***REMOVED***
			valid:    true,
			Type:     0x8,
			Flags:    0x0,
			Length:   0x4,
			StreamID: 0x1020304,
		***REMOVED***,
		Increment: 0x7060504,
	***REMOVED***
	if !reflect.DeepEqual(f, want) ***REMOVED***
		t.Errorf("parsed back %#v; want %#v", f, want)
	***REMOVED***
***REMOVED***

func TestWritePing(t *testing.T)    ***REMOVED*** testWritePing(t, false) ***REMOVED***
func TestWritePingAck(t *testing.T) ***REMOVED*** testWritePing(t, true) ***REMOVED***

func testWritePing(t *testing.T, ack bool) ***REMOVED***
	fr, buf := testFramer()
	if err := fr.WritePing(ack, [8]byte***REMOVED***1, 2, 3, 4, 5, 6, 7, 8***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	var wantFlags Flags
	if ack ***REMOVED***
		wantFlags = FlagPingAck
	***REMOVED***
	var wantEnc = "\x00\x00\x08\x06" + string(wantFlags) + "\x00\x00\x00\x00" + "\x01\x02\x03\x04\x05\x06\x07\x08"
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***

	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	want := &PingFrame***REMOVED***
		FrameHeader: FrameHeader***REMOVED***
			valid:    true,
			Type:     0x6,
			Flags:    wantFlags,
			Length:   0x8,
			StreamID: 0,
		***REMOVED***,
		Data: [8]byte***REMOVED***1, 2, 3, 4, 5, 6, 7, 8***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(f, want) ***REMOVED***
		t.Errorf("parsed back %#v; want %#v", f, want)
	***REMOVED***
***REMOVED***

func TestReadFrameHeader(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in   string
		want FrameHeader
	***REMOVED******REMOVED***
		***REMOVED***in: "\x00\x00\x00" + "\x00" + "\x00" + "\x00\x00\x00\x00", want: FrameHeader***REMOVED******REMOVED******REMOVED***,
		***REMOVED***in: "\x01\x02\x03" + "\x04" + "\x05" + "\x06\x07\x08\x09", want: FrameHeader***REMOVED***
			Length: 66051, Type: 4, Flags: 5, StreamID: 101124105,
		***REMOVED******REMOVED***,
		// Ignore high bit:
		***REMOVED***in: "\xff\xff\xff" + "\xff" + "\xff" + "\xff\xff\xff\xff", want: FrameHeader***REMOVED***
			Length: 16777215, Type: 255, Flags: 255, StreamID: 2147483647***REMOVED******REMOVED***,
		***REMOVED***in: "\xff\xff\xff" + "\xff" + "\xff" + "\x7f\xff\xff\xff", want: FrameHeader***REMOVED***
			Length: 16777215, Type: 255, Flags: 255, StreamID: 2147483647***REMOVED******REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		got, err := readFrameHeader(make([]byte, 9), strings.NewReader(tt.in))
		if err != nil ***REMOVED***
			t.Errorf("%d. readFrameHeader(%q) = %v", i, tt.in, err)
			continue
		***REMOVED***
		tt.want.valid = true
		if got != tt.want ***REMOVED***
			t.Errorf("%d. readFrameHeader(%q) = %+v; want %+v", i, tt.in, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReadWriteFrameHeader(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		len      uint32
		typ      FrameType
		flags    Flags
		streamID uint32
	***REMOVED******REMOVED***
		***REMOVED***len: 0, typ: 255, flags: 1, streamID: 0***REMOVED***,
		***REMOVED***len: 0, typ: 255, flags: 1, streamID: 1***REMOVED***,
		***REMOVED***len: 0, typ: 255, flags: 1, streamID: 255***REMOVED***,
		***REMOVED***len: 0, typ: 255, flags: 1, streamID: 256***REMOVED***,
		***REMOVED***len: 0, typ: 255, flags: 1, streamID: 65535***REMOVED***,
		***REMOVED***len: 0, typ: 255, flags: 1, streamID: 65536***REMOVED***,

		***REMOVED***len: 0, typ: 1, flags: 255, streamID: 1***REMOVED***,
		***REMOVED***len: 255, typ: 1, flags: 255, streamID: 1***REMOVED***,
		***REMOVED***len: 256, typ: 1, flags: 255, streamID: 1***REMOVED***,
		***REMOVED***len: 65535, typ: 1, flags: 255, streamID: 1***REMOVED***,
		***REMOVED***len: 65536, typ: 1, flags: 255, streamID: 1***REMOVED***,
		***REMOVED***len: 16777215, typ: 1, flags: 255, streamID: 1***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		fr, buf := testFramer()
		fr.startWrite(tt.typ, tt.flags, tt.streamID)
		fr.writeBytes(make([]byte, tt.len))
		fr.endWrite()
		fh, err := ReadFrameHeader(buf)
		if err != nil ***REMOVED***
			t.Errorf("ReadFrameHeader(%+v) = %v", tt, err)
			continue
		***REMOVED***
		if fh.Type != tt.typ || fh.Flags != tt.flags || fh.Length != tt.len || fh.StreamID != tt.streamID ***REMOVED***
			t.Errorf("ReadFrameHeader(%+v) = %+v; mismatch", tt, fh)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestWriteTooLargeFrame(t *testing.T) ***REMOVED***
	fr, _ := testFramer()
	fr.startWrite(0, 1, 1)
	fr.writeBytes(make([]byte, 1<<24))
	err := fr.endWrite()
	if err != ErrFrameTooLarge ***REMOVED***
		t.Errorf("endWrite = %v; want errFrameTooLarge", err)
	***REMOVED***
***REMOVED***

func TestWriteGoAway(t *testing.T) ***REMOVED***
	const debug = "foo"
	fr, buf := testFramer()
	if err := fr.WriteGoAway(0x01020304, 0x05060708, []byte(debug)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	const wantEnc = "\x00\x00\v\a\x00\x00\x00\x00\x00\x01\x02\x03\x04\x05\x06\x07\x08" + debug
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	want := &GoAwayFrame***REMOVED***
		FrameHeader: FrameHeader***REMOVED***
			valid:    true,
			Type:     0x7,
			Flags:    0,
			Length:   uint32(4 + 4 + len(debug)),
			StreamID: 0,
		***REMOVED***,
		LastStreamID: 0x01020304,
		ErrCode:      0x05060708,
		debugData:    []byte(debug),
	***REMOVED***
	if !reflect.DeepEqual(f, want) ***REMOVED***
		t.Fatalf("parsed back:\n%#v\nwant:\n%#v", f, want)
	***REMOVED***
	if got := string(f.(*GoAwayFrame).DebugData()); got != debug ***REMOVED***
		t.Errorf("debug data = %q; want %q", got, debug)
	***REMOVED***
***REMOVED***

func TestWritePushPromise(t *testing.T) ***REMOVED***
	pp := PushPromiseParam***REMOVED***
		StreamID:      42,
		PromiseID:     42,
		BlockFragment: []byte("abc"),
	***REMOVED***
	fr, buf := testFramer()
	if err := fr.WritePushPromise(pp); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	const wantEnc = "\x00\x00\x07\x05\x00\x00\x00\x00*\x00\x00\x00*abc"
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	_, ok := f.(*PushPromiseFrame)
	if !ok ***REMOVED***
		t.Fatalf("got %T; want *PushPromiseFrame", f)
	***REMOVED***
	want := &PushPromiseFrame***REMOVED***
		FrameHeader: FrameHeader***REMOVED***
			valid:    true,
			Type:     0x5,
			Flags:    0x0,
			Length:   0x7,
			StreamID: 42,
		***REMOVED***,
		PromiseID:     42,
		headerFragBuf: []byte("abc"),
	***REMOVED***
	if !reflect.DeepEqual(f, want) ***REMOVED***
		t.Fatalf("parsed back:\n%#v\nwant:\n%#v", f, want)
	***REMOVED***
***REMOVED***

// test checkFrameOrder and that HEADERS and CONTINUATION frames can't be intermingled.
func TestReadFrameOrder(t *testing.T) ***REMOVED***
	head := func(f *Framer, id uint32, end bool) ***REMOVED***
		f.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      id,
			BlockFragment: []byte("foo"), // unused, but non-empty
			EndHeaders:    end,
		***REMOVED***)
	***REMOVED***
	cont := func(f *Framer, id uint32, end bool) ***REMOVED***
		f.WriteContinuation(id, end, []byte("foo"))
	***REMOVED***

	tests := [...]struct ***REMOVED***
		name    string
		w       func(*Framer)
		atLeast int
		wantErr string
	***REMOVED******REMOVED***
		0: ***REMOVED***
			w: func(f *Framer) ***REMOVED***
				head(f, 1, true)
			***REMOVED***,
		***REMOVED***,
		1: ***REMOVED***
			w: func(f *Framer) ***REMOVED***
				head(f, 1, true)
				head(f, 2, true)
			***REMOVED***,
		***REMOVED***,
		2: ***REMOVED***
			wantErr: "got HEADERS for stream 2; expected CONTINUATION following HEADERS for stream 1",
			w: func(f *Framer) ***REMOVED***
				head(f, 1, false)
				head(f, 2, true)
			***REMOVED***,
		***REMOVED***,
		3: ***REMOVED***
			wantErr: "got DATA for stream 1; expected CONTINUATION following HEADERS for stream 1",
			w: func(f *Framer) ***REMOVED***
				head(f, 1, false)
			***REMOVED***,
		***REMOVED***,
		4: ***REMOVED***
			w: func(f *Framer) ***REMOVED***
				head(f, 1, false)
				cont(f, 1, true)
				head(f, 2, true)
			***REMOVED***,
		***REMOVED***,
		5: ***REMOVED***
			wantErr: "got CONTINUATION for stream 2; expected stream 1",
			w: func(f *Framer) ***REMOVED***
				head(f, 1, false)
				cont(f, 2, true)
				head(f, 2, true)
			***REMOVED***,
		***REMOVED***,
		6: ***REMOVED***
			wantErr: "unexpected CONTINUATION for stream 1",
			w: func(f *Framer) ***REMOVED***
				cont(f, 1, true)
			***REMOVED***,
		***REMOVED***,
		7: ***REMOVED***
			wantErr: "unexpected CONTINUATION for stream 1",
			w: func(f *Framer) ***REMOVED***
				cont(f, 1, false)
			***REMOVED***,
		***REMOVED***,
		8: ***REMOVED***
			wantErr: "HEADERS frame with stream ID 0",
			w: func(f *Framer) ***REMOVED***
				head(f, 0, true)
			***REMOVED***,
		***REMOVED***,
		9: ***REMOVED***
			wantErr: "CONTINUATION frame with stream ID 0",
			w: func(f *Framer) ***REMOVED***
				cont(f, 0, true)
			***REMOVED***,
		***REMOVED***,
		10: ***REMOVED***
			wantErr: "unexpected CONTINUATION for stream 1",
			atLeast: 5,
			w: func(f *Framer) ***REMOVED***
				head(f, 1, false)
				cont(f, 1, false)
				cont(f, 1, false)
				cont(f, 1, false)
				cont(f, 1, true)
				cont(f, 1, false)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		buf := new(bytes.Buffer)
		f := NewFramer(buf, buf)
		f.AllowIllegalWrites = true
		tt.w(f)
		f.WriteData(1, true, nil) // to test transition away from last step

		var err error
		n := 0
		var log bytes.Buffer
		for ***REMOVED***
			var got Frame
			got, err = f.ReadFrame()
			fmt.Fprintf(&log, "  read %v, %v\n", got, err)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			n++
		***REMOVED***
		if err == io.EOF ***REMOVED***
			err = nil
		***REMOVED***
		ok := tt.wantErr == ""
		if ok && err != nil ***REMOVED***
			t.Errorf("%d. after %d good frames, ReadFrame = %v; want success\n%s", i, n, err, log.Bytes())
			continue
		***REMOVED***
		if !ok && err != ConnectionError(ErrCodeProtocol) ***REMOVED***
			t.Errorf("%d. after %d good frames, ReadFrame = %v; want ConnectionError(ErrCodeProtocol)\n%s", i, n, err, log.Bytes())
			continue
		***REMOVED***
		if !((f.errDetail == nil && tt.wantErr == "") || (fmt.Sprint(f.errDetail) == tt.wantErr)) ***REMOVED***
			t.Errorf("%d. framer eror = %q; want %q\n%s", i, f.errDetail, tt.wantErr, log.Bytes())
		***REMOVED***
		if n < tt.atLeast ***REMOVED***
			t.Errorf("%d. framer only read %d frames; want at least %d\n%s", i, n, tt.atLeast, log.Bytes())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMetaFrameHeader(t *testing.T) ***REMOVED***
	write := func(f *Framer, frags ...[]byte) ***REMOVED***
		for i, frag := range frags ***REMOVED***
			end := (i == len(frags)-1)
			if i == 0 ***REMOVED***
				f.WriteHeaders(HeadersFrameParam***REMOVED***
					StreamID:      1,
					BlockFragment: frag,
					EndHeaders:    end,
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				f.WriteContinuation(1, end, frag)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	want := func(flags Flags, length uint32, pairs ...string) *MetaHeadersFrame ***REMOVED***
		mh := &MetaHeadersFrame***REMOVED***
			HeadersFrame: &HeadersFrame***REMOVED***
				FrameHeader: FrameHeader***REMOVED***
					Type:     FrameHeaders,
					Flags:    flags,
					Length:   length,
					StreamID: 1,
				***REMOVED***,
			***REMOVED***,
			Fields: []hpack.HeaderField(nil),
		***REMOVED***
		for len(pairs) > 0 ***REMOVED***
			mh.Fields = append(mh.Fields, hpack.HeaderField***REMOVED***
				Name:  pairs[0],
				Value: pairs[1],
			***REMOVED***)
			pairs = pairs[2:]
		***REMOVED***
		return mh
	***REMOVED***
	truncated := func(mh *MetaHeadersFrame) *MetaHeadersFrame ***REMOVED***
		mh.Truncated = true
		return mh
	***REMOVED***

	const noFlags Flags = 0

	oneKBString := strings.Repeat("a", 1<<10)

	tests := [...]struct ***REMOVED***
		name              string
		w                 func(*Framer)
		want              interface***REMOVED******REMOVED*** // *MetaHeaderFrame or error
		wantErrReason     string
		maxHeaderListSize uint32
	***REMOVED******REMOVED***
		0: ***REMOVED***
			name: "single_headers",
			w: func(f *Framer) ***REMOVED***
				var he hpackEncoder
				all := he.encodeHeaderRaw(t, ":method", "GET", ":path", "/")
				write(f, all)
			***REMOVED***,
			want: want(FlagHeadersEndHeaders, 2, ":method", "GET", ":path", "/"),
		***REMOVED***,
		1: ***REMOVED***
			name: "with_continuation",
			w: func(f *Framer) ***REMOVED***
				var he hpackEncoder
				all := he.encodeHeaderRaw(t, ":method", "GET", ":path", "/", "foo", "bar")
				write(f, all[:1], all[1:])
			***REMOVED***,
			want: want(noFlags, 1, ":method", "GET", ":path", "/", "foo", "bar"),
		***REMOVED***,
		2: ***REMOVED***
			name: "with_two_continuation",
			w: func(f *Framer) ***REMOVED***
				var he hpackEncoder
				all := he.encodeHeaderRaw(t, ":method", "GET", ":path", "/", "foo", "bar")
				write(f, all[:2], all[2:4], all[4:])
			***REMOVED***,
			want: want(noFlags, 2, ":method", "GET", ":path", "/", "foo", "bar"),
		***REMOVED***,
		3: ***REMOVED***
			name: "big_string_okay",
			w: func(f *Framer) ***REMOVED***
				var he hpackEncoder
				all := he.encodeHeaderRaw(t, ":method", "GET", ":path", "/", "foo", oneKBString)
				write(f, all[:2], all[2:])
			***REMOVED***,
			want: want(noFlags, 2, ":method", "GET", ":path", "/", "foo", oneKBString),
		***REMOVED***,
		4: ***REMOVED***
			name: "big_string_error",
			w: func(f *Framer) ***REMOVED***
				var he hpackEncoder
				all := he.encodeHeaderRaw(t, ":method", "GET", ":path", "/", "foo", oneKBString)
				write(f, all[:2], all[2:])
			***REMOVED***,
			maxHeaderListSize: (1 << 10) / 2,
			want:              ConnectionError(ErrCodeCompression),
		***REMOVED***,
		5: ***REMOVED***
			name: "max_header_list_truncated",
			w: func(f *Framer) ***REMOVED***
				var he hpackEncoder
				var pairs = []string***REMOVED***":method", "GET", ":path", "/"***REMOVED***
				for i := 0; i < 100; i++ ***REMOVED***
					pairs = append(pairs, "foo", "bar")
				***REMOVED***
				all := he.encodeHeaderRaw(t, pairs...)
				write(f, all[:2], all[2:])
			***REMOVED***,
			maxHeaderListSize: (1 << 10) / 2,
			want: truncated(want(noFlags, 2,
				":method", "GET",
				":path", "/",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar",
				"foo", "bar", // 11
			)),
		***REMOVED***,
		6: ***REMOVED***
			name: "pseudo_order",
			w: func(f *Framer) ***REMOVED***
				write(f, encodeHeaderRaw(t,
					":method", "GET",
					"foo", "bar",
					":path", "/", // bogus
				))
			***REMOVED***,
			want:          streamError(1, ErrCodeProtocol),
			wantErrReason: "pseudo header field after regular",
		***REMOVED***,
		7: ***REMOVED***
			name: "pseudo_unknown",
			w: func(f *Framer) ***REMOVED***
				write(f, encodeHeaderRaw(t,
					":unknown", "foo", // bogus
					"foo", "bar",
				))
			***REMOVED***,
			want:          streamError(1, ErrCodeProtocol),
			wantErrReason: "invalid pseudo-header \":unknown\"",
		***REMOVED***,
		8: ***REMOVED***
			name: "pseudo_mix_request_response",
			w: func(f *Framer) ***REMOVED***
				write(f, encodeHeaderRaw(t,
					":method", "GET",
					":status", "100",
				))
			***REMOVED***,
			want:          streamError(1, ErrCodeProtocol),
			wantErrReason: "mix of request and response pseudo headers",
		***REMOVED***,
		9: ***REMOVED***
			name: "pseudo_dup",
			w: func(f *Framer) ***REMOVED***
				write(f, encodeHeaderRaw(t,
					":method", "GET",
					":method", "POST",
				))
			***REMOVED***,
			want:          streamError(1, ErrCodeProtocol),
			wantErrReason: "duplicate pseudo-header \":method\"",
		***REMOVED***,
		10: ***REMOVED***
			name: "trailer_okay_no_pseudo",
			w:    func(f *Framer) ***REMOVED*** write(f, encodeHeaderRaw(t, "foo", "bar")) ***REMOVED***,
			want: want(FlagHeadersEndHeaders, 8, "foo", "bar"),
		***REMOVED***,
		11: ***REMOVED***
			name:          "invalid_field_name",
			w:             func(f *Framer) ***REMOVED*** write(f, encodeHeaderRaw(t, "CapitalBad", "x")) ***REMOVED***,
			want:          streamError(1, ErrCodeProtocol),
			wantErrReason: "invalid header field name \"CapitalBad\"",
		***REMOVED***,
		12: ***REMOVED***
			name:          "invalid_field_value",
			w:             func(f *Framer) ***REMOVED*** write(f, encodeHeaderRaw(t, "key", "bad_null\x00")) ***REMOVED***,
			want:          streamError(1, ErrCodeProtocol),
			wantErrReason: "invalid header field value \"bad_null\\x00\"",
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		buf := new(bytes.Buffer)
		f := NewFramer(buf, buf)
		f.ReadMetaHeaders = hpack.NewDecoder(initialHeaderTableSize, nil)
		f.MaxHeaderListSize = tt.maxHeaderListSize
		tt.w(f)

		name := tt.name
		if name == "" ***REMOVED***
			name = fmt.Sprintf("test index %d", i)
		***REMOVED***

		var got interface***REMOVED******REMOVED***
		var err error
		got, err = f.ReadFrame()
		if err != nil ***REMOVED***
			got = err

			// Ignore the StreamError.Cause field, if it matches the wantErrReason.
			// The test table above predates the Cause field.
			if se, ok := err.(StreamError); ok && se.Cause != nil && se.Cause.Error() == tt.wantErrReason ***REMOVED***
				se.Cause = nil
				got = se
			***REMOVED***
		***REMOVED***
		if !reflect.DeepEqual(got, tt.want) ***REMOVED***
			if mhg, ok := got.(*MetaHeadersFrame); ok ***REMOVED***
				if mhw, ok := tt.want.(*MetaHeadersFrame); ok ***REMOVED***
					hg := mhg.HeadersFrame
					hw := mhw.HeadersFrame
					if hg != nil && hw != nil && !reflect.DeepEqual(*hg, *hw) ***REMOVED***
						t.Errorf("%s: headers differ:\n got: %+v\nwant: %+v\n", name, *hg, *hw)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			str := func(v interface***REMOVED******REMOVED***) string ***REMOVED***
				if _, ok := v.(error); ok ***REMOVED***
					return fmt.Sprintf("error %v", v)
				***REMOVED*** else ***REMOVED***
					return fmt.Sprintf("value %#v", v)
				***REMOVED***
			***REMOVED***
			t.Errorf("%s:\n got: %v\nwant: %s", name, str(got), str(tt.want))
		***REMOVED***
		if tt.wantErrReason != "" && tt.wantErrReason != fmt.Sprint(f.errDetail) ***REMOVED***
			t.Errorf("%s: got error reason %q; want %q", name, f.errDetail, tt.wantErrReason)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSetReuseFrames(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	fr.SetReuseFrames()

	// Check that DataFrames are reused. Note that
	// SetReuseFrames only currently implements reuse of DataFrames.
	firstDf := readAndVerifyDataFrame("ABC", 3, fr, buf, t)

	for i := 0; i < 10; i++ ***REMOVED***
		df := readAndVerifyDataFrame("XYZ", 3, fr, buf, t)
		if df != firstDf ***REMOVED***
			t.Errorf("Expected Framer to return references to the same DataFrame. Have %v and %v", &df, &firstDf)
		***REMOVED***
	***REMOVED***

	for i := 0; i < 10; i++ ***REMOVED***
		df := readAndVerifyDataFrame("", 0, fr, buf, t)
		if df != firstDf ***REMOVED***
			t.Errorf("Expected Framer to return references to the same DataFrame. Have %v and %v", &df, &firstDf)
		***REMOVED***
	***REMOVED***

	for i := 0; i < 10; i++ ***REMOVED***
		df := readAndVerifyDataFrame("HHH", 3, fr, buf, t)
		if df != firstDf ***REMOVED***
			t.Errorf("Expected Framer to return references to the same DataFrame. Have %v and %v", &df, &firstDf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSetReuseFramesMoreThanOnce(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	fr.SetReuseFrames()

	firstDf := readAndVerifyDataFrame("ABC", 3, fr, buf, t)
	fr.SetReuseFrames()

	for i := 0; i < 10; i++ ***REMOVED***
		df := readAndVerifyDataFrame("XYZ", 3, fr, buf, t)
		// SetReuseFrames should be idempotent
		fr.SetReuseFrames()
		if df != firstDf ***REMOVED***
			t.Errorf("Expected Framer to return references to the same DataFrame. Have %v and %v", &df, &firstDf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNoSetReuseFrames(t *testing.T) ***REMOVED***
	fr, buf := testFramer()
	const numNewDataFrames = 10
	dfSoFar := make([]interface***REMOVED******REMOVED***, numNewDataFrames)

	// Check that DataFrames are not reused if SetReuseFrames wasn't called.
	// SetReuseFrames only currently implements reuse of DataFrames.
	for i := 0; i < numNewDataFrames; i++ ***REMOVED***
		df := readAndVerifyDataFrame("XYZ", 3, fr, buf, t)
		for _, item := range dfSoFar ***REMOVED***
			if df == item ***REMOVED***
				t.Errorf("Expected Framer to return new DataFrames since SetNoReuseFrames not set.")
			***REMOVED***
		***REMOVED***
		dfSoFar[i] = df
	***REMOVED***
***REMOVED***

func readAndVerifyDataFrame(data string, length byte, fr *Framer, buf *bytes.Buffer, t *testing.T) *DataFrame ***REMOVED***
	var streamID uint32 = 1<<24 + 2<<16 + 3<<8 + 4
	fr.WriteData(streamID, true, []byte(data))
	wantEnc := "\x00\x00" + string(length) + "\x00\x01\x01\x02\x03\x04" + data
	if buf.String() != wantEnc ***REMOVED***
		t.Errorf("encoded as %q; want %q", buf.Bytes(), wantEnc)
	***REMOVED***
	f, err := fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	df, ok := f.(*DataFrame)
	if !ok ***REMOVED***
		t.Fatalf("got %T; want *DataFrame", f)
	***REMOVED***
	if !bytes.Equal(df.Data(), []byte(data)) ***REMOVED***
		t.Errorf("got %q; want %q", df.Data(), []byte(data))
	***REMOVED***
	if f.Header().Flags&1 == 0 ***REMOVED***
		t.Errorf("didn't see END_STREAM flag")
	***REMOVED***
	return df
***REMOVED***

func encodeHeaderRaw(t *testing.T, pairs ...string) []byte ***REMOVED***
	var he hpackEncoder
	return he.encodeHeaderRaw(t, pairs...)
***REMOVED***
