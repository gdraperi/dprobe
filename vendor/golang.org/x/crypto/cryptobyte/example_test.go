// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cryptobyte_test

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/cryptobyte/asn1"
)

func ExampleString_lengthPrefixed() ***REMOVED***
	// This is an example of parsing length-prefixed data (as found in, for
	// example, TLS). Imagine a 16-bit prefixed series of 8-bit prefixed
	// strings.

	input := cryptobyte.String([]byte***REMOVED***0, 12, 5, 'h', 'e', 'l', 'l', 'o', 5, 'w', 'o', 'r', 'l', 'd'***REMOVED***)
	var result []string

	var values cryptobyte.String
	if !input.ReadUint16LengthPrefixed(&values) ||
		!input.Empty() ***REMOVED***
		panic("bad format")
	***REMOVED***

	for !values.Empty() ***REMOVED***
		var value cryptobyte.String
		if !values.ReadUint8LengthPrefixed(&value) ***REMOVED***
			panic("bad format")
		***REMOVED***

		result = append(result, string(value))
	***REMOVED***

	// Output: []string***REMOVED***"hello", "world"***REMOVED***
	fmt.Printf("%#v\n", result)
***REMOVED***

func ExampleString_aSN1() ***REMOVED***
	// This is an example of parsing ASN.1 data that looks like:
	//    Foo ::= SEQUENCE ***REMOVED***
	//      version [6] INTEGER DEFAULT 0
	//      data OCTET STRING
	//***REMOVED***

	input := cryptobyte.String([]byte***REMOVED***0x30, 12, 0xa6, 3, 2, 1, 2, 4, 5, 'h', 'e', 'l', 'l', 'o'***REMOVED***)

	var (
		version                   int64
		data, inner, versionBytes cryptobyte.String
		haveVersion               bool
	)
	if !input.ReadASN1(&inner, asn1.SEQUENCE) ||
		!input.Empty() ||
		!inner.ReadOptionalASN1(&versionBytes, &haveVersion, asn1.Tag(6).Constructed().ContextSpecific()) ||
		(haveVersion && !versionBytes.ReadASN1Integer(&version)) ||
		(haveVersion && !versionBytes.Empty()) ||
		!inner.ReadASN1(&data, asn1.OCTET_STRING) ||
		!inner.Empty() ***REMOVED***
		panic("bad format")
	***REMOVED***

	// Output: haveVersion: true, version: 2, data: hello
	fmt.Printf("haveVersion: %t, version: %d, data: %s\n", haveVersion, version, string(data))
***REMOVED***

func ExampleBuilder_aSN1() ***REMOVED***
	// This is an example of building ASN.1 data that looks like:
	//    Foo ::= SEQUENCE ***REMOVED***
	//      version [6] INTEGER DEFAULT 0
	//      data OCTET STRING
	//***REMOVED***

	version := int64(2)
	data := []byte("hello")
	const defaultVersion = 0

	var b cryptobyte.Builder
	b.AddASN1(asn1.SEQUENCE, func(b *cryptobyte.Builder) ***REMOVED***
		if version != defaultVersion ***REMOVED***
			b.AddASN1(asn1.Tag(6).Constructed().ContextSpecific(), func(b *cryptobyte.Builder) ***REMOVED***
				b.AddASN1Int64(version)
			***REMOVED***)
		***REMOVED***
		b.AddASN1OctetString(data)
	***REMOVED***)

	result, err := b.Bytes()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	// Output: 300ca603020102040568656c6c6f
	fmt.Printf("%x\n", result)
***REMOVED***

func ExampleBuilder_lengthPrefixed() ***REMOVED***
	// This is an example of building length-prefixed data (as found in,
	// for example, TLS). Imagine a 16-bit prefixed series of 8-bit
	// prefixed strings.
	input := []string***REMOVED***"hello", "world"***REMOVED***

	var b cryptobyte.Builder
	b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) ***REMOVED***
		for _, value := range input ***REMOVED***
			b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) ***REMOVED***
				b.AddBytes([]byte(value))
			***REMOVED***)
		***REMOVED***
	***REMOVED***)

	result, err := b.Bytes()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	// Output: 000c0568656c6c6f05776f726c64
	fmt.Printf("%x\n", result)
***REMOVED***

func ExampleBuilder_lengthPrefixOverflow() ***REMOVED***
	// Writing more data that can be expressed by the length prefix results
	// in an error from Bytes().

	tooLarge := make([]byte, 256)

	var b cryptobyte.Builder
	b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) ***REMOVED***
		b.AddBytes(tooLarge)
	***REMOVED***)

	result, err := b.Bytes()
	fmt.Printf("len=%d err=%s\n", len(result), err)

	// Output: len=0 err=cryptobyte: pending child length 256 exceeds 1-byte length prefix
***REMOVED***

func ExampleBuilderContinuation_errorHandling() ***REMOVED***
	var b cryptobyte.Builder
	// Continuations that panic with a BuildError will cause Bytes to
	// return the inner error.
	b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) ***REMOVED***
		b.AddUint32(0)
		panic(cryptobyte.BuildError***REMOVED***Err: errors.New("example error")***REMOVED***)
	***REMOVED***)

	result, err := b.Bytes()
	fmt.Printf("len=%d err=%s\n", len(result), err)

	// Output: len=0 err=example error
***REMOVED***
