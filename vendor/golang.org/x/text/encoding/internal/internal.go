// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal contains code that is shared among encoding implementations.
package internal

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)

// Encoding is an implementation of the Encoding interface that adds the String
// and ID methods to an existing encoding.
type Encoding struct ***REMOVED***
	encoding.Encoding
	Name string
	MIB  identifier.MIB
***REMOVED***

// _ verifies that Encoding implements identifier.Interface.
var _ identifier.Interface = (*Encoding)(nil)

func (e *Encoding) String() string ***REMOVED***
	return e.Name
***REMOVED***

func (e *Encoding) ID() (mib identifier.MIB, other string) ***REMOVED***
	return e.MIB, ""
***REMOVED***

// SimpleEncoding is an Encoding that combines two Transformers.
type SimpleEncoding struct ***REMOVED***
	Decoder transform.Transformer
	Encoder transform.Transformer
***REMOVED***

func (e *SimpleEncoding) NewDecoder() *encoding.Decoder ***REMOVED***
	return &encoding.Decoder***REMOVED***Transformer: e.Decoder***REMOVED***
***REMOVED***

func (e *SimpleEncoding) NewEncoder() *encoding.Encoder ***REMOVED***
	return &encoding.Encoder***REMOVED***Transformer: e.Encoder***REMOVED***
***REMOVED***

// FuncEncoding is an Encoding that combines two functions returning a new
// Transformer.
type FuncEncoding struct ***REMOVED***
	Decoder func() transform.Transformer
	Encoder func() transform.Transformer
***REMOVED***

func (e FuncEncoding) NewDecoder() *encoding.Decoder ***REMOVED***
	return &encoding.Decoder***REMOVED***Transformer: e.Decoder()***REMOVED***
***REMOVED***

func (e FuncEncoding) NewEncoder() *encoding.Encoder ***REMOVED***
	return &encoding.Encoder***REMOVED***Transformer: e.Encoder()***REMOVED***
***REMOVED***

// A RepertoireError indicates a rune is not in the repertoire of a destination
// encoding. It is associated with an encoding-specific suggested replacement
// byte.
type RepertoireError byte

// Error implements the error interrface.
func (r RepertoireError) Error() string ***REMOVED***
	return "encoding: rune not supported by encoding."
***REMOVED***

// Replacement returns the replacement string associated with this error.
func (r RepertoireError) Replacement() byte ***REMOVED*** return byte(r) ***REMOVED***

var ErrASCIIReplacement = RepertoireError(encoding.ASCIISub)
