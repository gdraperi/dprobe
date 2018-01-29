// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

var intLengthTests = []struct ***REMOVED***
	val, length int
***REMOVED******REMOVED***
	***REMOVED***0, 4 + 0***REMOVED***,
	***REMOVED***1, 4 + 1***REMOVED***,
	***REMOVED***127, 4 + 1***REMOVED***,
	***REMOVED***128, 4 + 2***REMOVED***,
	***REMOVED***-1, 4 + 1***REMOVED***,
***REMOVED***

func TestIntLength(t *testing.T) ***REMOVED***
	for _, test := range intLengthTests ***REMOVED***
		v := new(big.Int).SetInt64(int64(test.val))
		length := intLength(v)
		if length != test.length ***REMOVED***
			t.Errorf("For %d, got length %d but expected %d", test.val, length, test.length)
		***REMOVED***
	***REMOVED***
***REMOVED***

type msgAllTypes struct ***REMOVED***
	Bool    bool `sshtype:"21"`
	Array   [16]byte
	Uint64  uint64
	Uint32  uint32
	Uint8   uint8
	String  string
	Strings []string
	Bytes   []byte
	Int     *big.Int
	Rest    []byte `ssh:"rest"`
***REMOVED***

func (t *msgAllTypes) Generate(rand *rand.Rand, size int) reflect.Value ***REMOVED***
	m := &msgAllTypes***REMOVED******REMOVED***
	m.Bool = rand.Intn(2) == 1
	randomBytes(m.Array[:], rand)
	m.Uint64 = uint64(rand.Int63n(1<<63 - 1))
	m.Uint32 = uint32(rand.Intn((1 << 31) - 1))
	m.Uint8 = uint8(rand.Intn(1 << 8))
	m.String = string(m.Array[:])
	m.Strings = randomNameList(rand)
	m.Bytes = m.Array[:]
	m.Int = randomInt(rand)
	m.Rest = m.Array[:]
	return reflect.ValueOf(m)
***REMOVED***

func TestMarshalUnmarshal(t *testing.T) ***REMOVED***
	rand := rand.New(rand.NewSource(0))
	iface := &msgAllTypes***REMOVED******REMOVED***
	ty := reflect.ValueOf(iface).Type()

	n := 100
	if testing.Short() ***REMOVED***
		n = 5
	***REMOVED***
	for j := 0; j < n; j++ ***REMOVED***
		v, ok := quick.Value(ty, rand)
		if !ok ***REMOVED***
			t.Errorf("failed to create value")
			break
		***REMOVED***

		m1 := v.Elem().Interface()
		m2 := iface

		marshaled := Marshal(m1)
		if err := Unmarshal(marshaled, m2); err != nil ***REMOVED***
			t.Errorf("Unmarshal %#v: %s", m1, err)
			break
		***REMOVED***

		if !reflect.DeepEqual(v.Interface(), m2) ***REMOVED***
			t.Errorf("got: %#v\nwant:%#v\n%x", m2, m1, marshaled)
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUnmarshalEmptyPacket(t *testing.T) ***REMOVED***
	var b []byte
	var m channelRequestSuccessMsg
	if err := Unmarshal(b, &m); err == nil ***REMOVED***
		t.Fatalf("unmarshal of empty slice succeeded")
	***REMOVED***
***REMOVED***

func TestUnmarshalUnexpectedPacket(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		I uint32 `sshtype:"43"`
		S string
		B bool
	***REMOVED***

	s := S***REMOVED***11, "hello", true***REMOVED***
	packet := Marshal(s)
	packet[0] = 42
	roundtrip := S***REMOVED******REMOVED***
	err := Unmarshal(packet, &roundtrip)
	if err == nil ***REMOVED***
		t.Fatal("expected error, not nil")
	***REMOVED***
***REMOVED***

func TestMarshalPtr(t *testing.T) ***REMOVED***
	s := struct ***REMOVED***
		S string
	***REMOVED******REMOVED***"hello"***REMOVED***

	m1 := Marshal(s)
	m2 := Marshal(&s)
	if !bytes.Equal(m1, m2) ***REMOVED***
		t.Errorf("got %q, want %q for marshaled pointer", m2, m1)
	***REMOVED***
***REMOVED***

func TestBareMarshalUnmarshal(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		I uint32
		S string
		B bool
	***REMOVED***

	s := S***REMOVED***42, "hello", true***REMOVED***
	packet := Marshal(s)
	roundtrip := S***REMOVED******REMOVED***
	Unmarshal(packet, &roundtrip)

	if !reflect.DeepEqual(s, roundtrip) ***REMOVED***
		t.Errorf("got %#v, want %#v", roundtrip, s)
	***REMOVED***
***REMOVED***

func TestBareMarshal(t *testing.T) ***REMOVED***
	type S2 struct ***REMOVED***
		I uint32
	***REMOVED***
	s := S2***REMOVED***42***REMOVED***
	packet := Marshal(s)
	i, rest, ok := parseUint32(packet)
	if len(rest) > 0 || !ok ***REMOVED***
		t.Errorf("parseInt(%q): parse error", packet)
	***REMOVED***
	if i != s.I ***REMOVED***
		t.Errorf("got %d, want %d", i, s.I)
	***REMOVED***
***REMOVED***

func TestUnmarshalShortKexInitPacket(t *testing.T) ***REMOVED***
	// This used to panic.
	// Issue 11348
	packet := []byte***REMOVED***0x14, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0xff, 0xff, 0xff, 0xff***REMOVED***
	kim := &kexInitMsg***REMOVED******REMOVED***
	if err := Unmarshal(packet, kim); err == nil ***REMOVED***
		t.Error("truncated packet unmarshaled without error")
	***REMOVED***
***REMOVED***

func TestMarshalMultiTag(t *testing.T) ***REMOVED***
	var res struct ***REMOVED***
		A uint32 `sshtype:"1|2"`
	***REMOVED***

	good1 := struct ***REMOVED***
		A uint32 `sshtype:"1"`
	***REMOVED******REMOVED***
		1,
	***REMOVED***
	good2 := struct ***REMOVED***
		A uint32 `sshtype:"2"`
	***REMOVED******REMOVED***
		1,
	***REMOVED***

	if e := Unmarshal(Marshal(good1), &res); e != nil ***REMOVED***
		t.Errorf("error unmarshaling multipart tag: %v", e)
	***REMOVED***

	if e := Unmarshal(Marshal(good2), &res); e != nil ***REMOVED***
		t.Errorf("error unmarshaling multipart tag: %v", e)
	***REMOVED***

	bad1 := struct ***REMOVED***
		A uint32 `sshtype:"3"`
	***REMOVED******REMOVED***
		1,
	***REMOVED***
	if e := Unmarshal(Marshal(bad1), &res); e == nil ***REMOVED***
		t.Errorf("bad struct unmarshaled without error")
	***REMOVED***
***REMOVED***

func randomBytes(out []byte, rand *rand.Rand) ***REMOVED***
	for i := 0; i < len(out); i++ ***REMOVED***
		out[i] = byte(rand.Int31())
	***REMOVED***
***REMOVED***

func randomNameList(rand *rand.Rand) []string ***REMOVED***
	ret := make([]string, rand.Int31()&15)
	for i := range ret ***REMOVED***
		s := make([]byte, 1+(rand.Int31()&15))
		for j := range s ***REMOVED***
			s[j] = 'a' + uint8(rand.Int31()&15)
		***REMOVED***
		ret[i] = string(s)
	***REMOVED***
	return ret
***REMOVED***

func randomInt(rand *rand.Rand) *big.Int ***REMOVED***
	return new(big.Int).SetInt64(int64(int32(rand.Uint32())))
***REMOVED***

func (*kexInitMsg) Generate(rand *rand.Rand, size int) reflect.Value ***REMOVED***
	ki := &kexInitMsg***REMOVED******REMOVED***
	randomBytes(ki.Cookie[:], rand)
	ki.KexAlgos = randomNameList(rand)
	ki.ServerHostKeyAlgos = randomNameList(rand)
	ki.CiphersClientServer = randomNameList(rand)
	ki.CiphersServerClient = randomNameList(rand)
	ki.MACsClientServer = randomNameList(rand)
	ki.MACsServerClient = randomNameList(rand)
	ki.CompressionClientServer = randomNameList(rand)
	ki.CompressionServerClient = randomNameList(rand)
	ki.LanguagesClientServer = randomNameList(rand)
	ki.LanguagesServerClient = randomNameList(rand)
	if rand.Int31()&1 == 1 ***REMOVED***
		ki.FirstKexFollows = true
	***REMOVED***
	return reflect.ValueOf(ki)
***REMOVED***

func (*kexDHInitMsg) Generate(rand *rand.Rand, size int) reflect.Value ***REMOVED***
	dhi := &kexDHInitMsg***REMOVED******REMOVED***
	dhi.X = randomInt(rand)
	return reflect.ValueOf(dhi)
***REMOVED***

var (
	_kexInitMsg   = new(kexInitMsg).Generate(rand.New(rand.NewSource(0)), 10).Elem().Interface()
	_kexDHInitMsg = new(kexDHInitMsg).Generate(rand.New(rand.NewSource(0)), 10).Elem().Interface()

	_kexInit   = Marshal(_kexInitMsg)
	_kexDHInit = Marshal(_kexDHInitMsg)
)

func BenchmarkMarshalKexInitMsg(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		Marshal(_kexInitMsg)
	***REMOVED***
***REMOVED***

func BenchmarkUnmarshalKexInitMsg(b *testing.B) ***REMOVED***
	m := new(kexInitMsg)
	for i := 0; i < b.N; i++ ***REMOVED***
		Unmarshal(_kexInit, m)
	***REMOVED***
***REMOVED***

func BenchmarkMarshalKexDHInitMsg(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		Marshal(_kexDHInitMsg)
	***REMOVED***
***REMOVED***

func BenchmarkUnmarshalKexDHInitMsg(b *testing.B) ***REMOVED***
	m := new(kexDHInitMsg)
	for i := 0; i < b.N; i++ ***REMOVED***
		Unmarshal(_kexDHInit, m)
	***REMOVED***
***REMOVED***
