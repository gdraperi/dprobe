// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cryptobyte

import (
	"bytes"
	encoding_asn1 "encoding/asn1"
	"math/big"
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/cryptobyte/asn1"
)

type readASN1Test struct ***REMOVED***
	name string
	in   []byte
	tag  asn1.Tag
	ok   bool
	out  interface***REMOVED******REMOVED***
***REMOVED***

var readASN1TestData = []readASN1Test***REMOVED***
	***REMOVED***"valid", []byte***REMOVED***0x30, 2, 1, 2***REMOVED***, 0x30, true, []byte***REMOVED***1, 2***REMOVED******REMOVED***,
	***REMOVED***"truncated", []byte***REMOVED***0x30, 3, 1, 2***REMOVED***, 0x30, false, nil***REMOVED***,
	***REMOVED***"zero length of length", []byte***REMOVED***0x30, 0x80***REMOVED***, 0x30, false, nil***REMOVED***,
	***REMOVED***"invalid long form length", []byte***REMOVED***0x30, 0x81, 1, 1***REMOVED***, 0x30, false, nil***REMOVED***,
	***REMOVED***"non-minimal length", append([]byte***REMOVED***0x30, 0x82, 0, 0x80***REMOVED***, make([]byte, 0x80)...), 0x30, false, nil***REMOVED***,
	***REMOVED***"invalid tag", []byte***REMOVED***0xa1, 3, 0x4, 1, 1***REMOVED***, 31, false, nil***REMOVED***,
	***REMOVED***"high tag", []byte***REMOVED***0x1f, 0x81, 0x80, 0x01, 2, 1, 2***REMOVED***, 0xff /* actually 0x4001, but tag is uint8 */, false, nil***REMOVED***,
***REMOVED***

func TestReadASN1(t *testing.T) ***REMOVED***
	for _, test := range readASN1TestData ***REMOVED***
		t.Run(test.name, func(t *testing.T) ***REMOVED***
			var in, out String = test.in, nil
			ok := in.ReadASN1(&out, test.tag)
			if ok != test.ok || ok && !bytes.Equal(out, test.out.([]byte)) ***REMOVED***
				t.Errorf("in.ReadASN1() = %v, want %v; out = %v, want %v", ok, test.ok, out, test.out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestReadASN1Optional(t *testing.T) ***REMOVED***
	var empty String
	var present bool
	ok := empty.ReadOptionalASN1(nil, &present, 0xa0)
	if !ok || present ***REMOVED***
		t.Errorf("empty.ReadOptionalASN1() = %v, want true; present = %v want false", ok, present)
	***REMOVED***

	var in, out String = []byte***REMOVED***0xa1, 3, 0x4, 1, 1***REMOVED***, nil
	ok = in.ReadOptionalASN1(&out, &present, 0xa0)
	if !ok || present ***REMOVED***
		t.Errorf("in.ReadOptionalASN1() = %v, want true, present = %v, want false", ok, present)
	***REMOVED***
	ok = in.ReadOptionalASN1(&out, &present, 0xa1)
	wantBytes := []byte***REMOVED***4, 1, 1***REMOVED***
	if !ok || !present || !bytes.Equal(out, wantBytes) ***REMOVED***
		t.Errorf("in.ReadOptionalASN1() = %v, want true; present = %v, want true; out = %v, want = %v", ok, present, out, wantBytes)
	***REMOVED***
***REMOVED***

var optionalOctetStringTestData = []struct ***REMOVED***
	readASN1Test
	present bool
***REMOVED******REMOVED***
	***REMOVED***readASN1Test***REMOVED***"empty", []byte***REMOVED******REMOVED***, 0xa0, true, []byte***REMOVED******REMOVED******REMOVED***, false***REMOVED***,
	***REMOVED***readASN1Test***REMOVED***"invalid", []byte***REMOVED***0xa1, 3, 0x4, 2, 1***REMOVED***, 0xa1, false, []byte***REMOVED******REMOVED******REMOVED***, true***REMOVED***,
	***REMOVED***readASN1Test***REMOVED***"missing", []byte***REMOVED***0xa1, 3, 0x4, 1, 1***REMOVED***, 0xa0, true, []byte***REMOVED******REMOVED******REMOVED***, false***REMOVED***,
	***REMOVED***readASN1Test***REMOVED***"present", []byte***REMOVED***0xa1, 3, 0x4, 1, 1***REMOVED***, 0xa1, true, []byte***REMOVED***1***REMOVED******REMOVED***, true***REMOVED***,
***REMOVED***

func TestReadASN1OptionalOctetString(t *testing.T) ***REMOVED***
	for _, test := range optionalOctetStringTestData ***REMOVED***
		t.Run(test.name, func(t *testing.T) ***REMOVED***
			in := String(test.in)
			var out []byte
			var present bool
			ok := in.ReadOptionalASN1OctetString(&out, &present, test.tag)
			if ok != test.ok || present != test.present || !bytes.Equal(out, test.out.([]byte)) ***REMOVED***
				t.Errorf("in.ReadOptionalASN1OctetString() = %v, want %v; present = %v want %v; out = %v, want %v", ok, test.ok, present, test.present, out, test.out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

const defaultInt = -1

var optionalIntTestData = []readASN1Test***REMOVED***
	***REMOVED***"empty", []byte***REMOVED******REMOVED***, 0xa0, true, defaultInt***REMOVED***,
	***REMOVED***"invalid", []byte***REMOVED***0xa1, 3, 0x2, 2, 127***REMOVED***, 0xa1, false, 0***REMOVED***,
	***REMOVED***"missing", []byte***REMOVED***0xa1, 3, 0x2, 1, 127***REMOVED***, 0xa0, true, defaultInt***REMOVED***,
	***REMOVED***"present", []byte***REMOVED***0xa1, 3, 0x2, 1, 42***REMOVED***, 0xa1, true, 42***REMOVED***,
***REMOVED***

func TestReadASN1OptionalInteger(t *testing.T) ***REMOVED***
	for _, test := range optionalIntTestData ***REMOVED***
		t.Run(test.name, func(t *testing.T) ***REMOVED***
			in := String(test.in)
			var out int
			ok := in.ReadOptionalASN1Integer(&out, test.tag, defaultInt)
			if ok != test.ok || ok && out != test.out.(int) ***REMOVED***
				t.Errorf("in.ReadOptionalASN1Integer() = %v, want %v; out = %v, want %v", ok, test.ok, out, test.out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestReadASN1IntegerSigned(t *testing.T) ***REMOVED***
	testData64 := []struct ***REMOVED***
		in  []byte
		out int64
	***REMOVED******REMOVED***
		***REMOVED***[]byte***REMOVED***2, 3, 128, 0, 0***REMOVED***, -0x800000***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 2, 255, 0***REMOVED***, -256***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 2, 255, 127***REMOVED***, -129***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 128***REMOVED***, -128***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 255***REMOVED***, -1***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 0***REMOVED***, 0***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 1***REMOVED***, 1***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 2***REMOVED***, 2***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 127***REMOVED***, 127***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 2, 0, 128***REMOVED***, 128***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 2, 1, 0***REMOVED***, 256***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 4, 0, 128, 0, 0***REMOVED***, 0x800000***REMOVED***,
	***REMOVED***
	for i, test := range testData64 ***REMOVED***
		in := String(test.in)
		var out int64
		ok := in.ReadASN1Integer(&out)
		if !ok || out != test.out ***REMOVED***
			t.Errorf("#%d: in.ReadASN1Integer() = %v, want true; out = %d, want %d", i, ok, out, test.out)
		***REMOVED***
	***REMOVED***

	// Repeat the same cases, reading into a big.Int.
	t.Run("big.Int", func(t *testing.T) ***REMOVED***
		for i, test := range testData64 ***REMOVED***
			in := String(test.in)
			var out big.Int
			ok := in.ReadASN1Integer(&out)
			if !ok || out.Int64() != test.out ***REMOVED***
				t.Errorf("#%d: in.ReadASN1Integer() = %v, want true; out = %d, want %d", i, ok, out.Int64(), test.out)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestReadASN1IntegerUnsigned(t *testing.T) ***REMOVED***
	testData := []struct ***REMOVED***
		in  []byte
		out uint64
	***REMOVED******REMOVED***
		***REMOVED***[]byte***REMOVED***2, 1, 0***REMOVED***, 0***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 1***REMOVED***, 1***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 2***REMOVED***, 2***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 1, 127***REMOVED***, 127***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 2, 0, 128***REMOVED***, 128***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 2, 1, 0***REMOVED***, 256***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 4, 0, 128, 0, 0***REMOVED***, 0x800000***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 8, 127, 255, 255, 255, 255, 255, 255, 255***REMOVED***, 0x7fffffffffffffff***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 9, 0, 128, 0, 0, 0, 0, 0, 0, 0***REMOVED***, 0x8000000000000000***REMOVED***,
		***REMOVED***[]byte***REMOVED***2, 9, 0, 255, 255, 255, 255, 255, 255, 255, 255***REMOVED***, 0xffffffffffffffff***REMOVED***,
	***REMOVED***
	for i, test := range testData ***REMOVED***
		in := String(test.in)
		var out uint64
		ok := in.ReadASN1Integer(&out)
		if !ok || out != test.out ***REMOVED***
			t.Errorf("#%d: in.ReadASN1Integer() = %v, want true; out = %d, want %d", i, ok, out, test.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReadASN1IntegerInvalid(t *testing.T) ***REMOVED***
	testData := []String***REMOVED***
		[]byte***REMOVED***3, 1, 0***REMOVED***, // invalid tag
		// truncated
		[]byte***REMOVED***2, 1***REMOVED***,
		[]byte***REMOVED***2, 2, 0***REMOVED***,
		// not minimally encoded
		[]byte***REMOVED***2, 2, 0, 1***REMOVED***,
		[]byte***REMOVED***2, 2, 0xff, 0xff***REMOVED***,
	***REMOVED***

	for i, test := range testData ***REMOVED***
		var out int64
		if test.ReadASN1Integer(&out) ***REMOVED***
			t.Errorf("#%d: in.ReadASN1Integer() = true, want false (out = %d)", i, out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestASN1ObjectIdentifier(t *testing.T) ***REMOVED***
	testData := []struct ***REMOVED***
		in  []byte
		ok  bool
		out []int
	***REMOVED******REMOVED***
		***REMOVED***[]byte***REMOVED******REMOVED***, false, []int***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***6, 0***REMOVED***, false, []int***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***5, 1, 85***REMOVED***, false, []int***REMOVED***2, 5***REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***6, 1, 85***REMOVED***, true, []int***REMOVED***2, 5***REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***6, 2, 85, 0x02***REMOVED***, true, []int***REMOVED***2, 5, 2***REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***6, 4, 85, 0x02, 0xc0, 0x00***REMOVED***, true, []int***REMOVED***2, 5, 2, 0x2000***REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***6, 3, 0x81, 0x34, 0x03***REMOVED***, true, []int***REMOVED***2, 100, 3***REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***6, 7, 85, 0x02, 0xc0, 0x80, 0x80, 0x80, 0x80***REMOVED***, false, []int***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for i, test := range testData ***REMOVED***
		in := String(test.in)
		var out encoding_asn1.ObjectIdentifier
		ok := in.ReadASN1ObjectIdentifier(&out)
		if ok != test.ok || ok && !out.Equal(test.out) ***REMOVED***
			t.Errorf("#%d: in.ReadASN1ObjectIdentifier() = %v, want %v; out = %v, want %v", i, ok, test.ok, out, test.out)
			continue
		***REMOVED***

		var b Builder
		b.AddASN1ObjectIdentifier(out)
		result, err := b.Bytes()
		if builderOk := err == nil; test.ok != builderOk ***REMOVED***
			t.Errorf("#%d: error from Builder.Bytes: %s", i, err)
			continue
		***REMOVED***
		if test.ok && !bytes.Equal(result, test.in) ***REMOVED***
			t.Errorf("#%d: reserialisation didn't match, got %x, want %x", i, result, test.in)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReadASN1GeneralizedTime(t *testing.T) ***REMOVED***
	testData := []struct ***REMOVED***
		in  string
		ok  bool
		out time.Time
	***REMOVED******REMOVED***
		***REMOVED***"20100102030405Z", true, time.Date(2010, 01, 02, 03, 04, 05, 0, time.UTC)***REMOVED***,
		***REMOVED***"20100102030405", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100102030405+0607", true, time.Date(2010, 01, 02, 03, 04, 05, 0, time.FixedZone("", 6*60*60+7*60))***REMOVED***,
		***REMOVED***"20100102030405-0607", true, time.Date(2010, 01, 02, 03, 04, 05, 0, time.FixedZone("", -6*60*60-7*60))***REMOVED***,
		/* These are invalid times. However, the time package normalises times
		 * and they were accepted in some versions. See #11134. */
		***REMOVED***"00000100000000Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20101302030405Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100002030405Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100100030405Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100132030405Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100231030405Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100102240405Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100102036005Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100102030460Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"-20100102030410Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"2010-0102030410Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"2010-0002030410Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"201001-02030410Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"20100102-030410Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"2010010203-0410Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
		***REMOVED***"201001020304-10Z", false, time.Time***REMOVED******REMOVED******REMOVED***,
	***REMOVED***
	for i, test := range testData ***REMOVED***
		in := String(append([]byte***REMOVED***byte(asn1.GeneralizedTime), byte(len(test.in))***REMOVED***, test.in...))
		var out time.Time
		ok := in.ReadASN1GeneralizedTime(&out)
		if ok != test.ok || ok && !reflect.DeepEqual(out, test.out) ***REMOVED***
			t.Errorf("#%d: in.ReadASN1GeneralizedTime() = %v, want %v; out = %q, want %q", i, ok, test.ok, out, test.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReadASN1BitString(t *testing.T) ***REMOVED***
	testData := []struct ***REMOVED***
		in  []byte
		ok  bool
		out encoding_asn1.BitString
	***REMOVED******REMOVED***
		***REMOVED***[]byte***REMOVED******REMOVED***, false, encoding_asn1.BitString***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***0x00***REMOVED***, true, encoding_asn1.BitString***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***0x07, 0x00***REMOVED***, true, encoding_asn1.BitString***REMOVED***Bytes: []byte***REMOVED***0***REMOVED***, BitLength: 1***REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***0x07, 0x01***REMOVED***, false, encoding_asn1.BitString***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***0x07, 0x40***REMOVED***, false, encoding_asn1.BitString***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***0x08, 0x00***REMOVED***, false, encoding_asn1.BitString***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***0xff***REMOVED***, false, encoding_asn1.BitString***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte***REMOVED***0xfe, 0x00***REMOVED***, false, encoding_asn1.BitString***REMOVED******REMOVED******REMOVED***,
	***REMOVED***
	for i, test := range testData ***REMOVED***
		in := String(append([]byte***REMOVED***3, byte(len(test.in))***REMOVED***, test.in...))
		var out encoding_asn1.BitString
		ok := in.ReadASN1BitString(&out)
		if ok != test.ok || ok && (!bytes.Equal(out.Bytes, test.out.Bytes) || out.BitLength != test.out.BitLength) ***REMOVED***
			t.Errorf("#%d: in.ReadASN1BitString() = %v, want %v; out = %v, want %v", i, ok, test.ok, out, test.out)
		***REMOVED***
	***REMOVED***
***REMOVED***
