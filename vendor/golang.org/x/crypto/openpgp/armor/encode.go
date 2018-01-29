// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package armor

import (
	"encoding/base64"
	"io"
)

var armorHeaderSep = []byte(": ")
var blockEnd = []byte("\n=")
var newline = []byte("\n")
var armorEndOfLineOut = []byte("-----\n")

// writeSlices writes its arguments to the given Writer.
func writeSlices(out io.Writer, slices ...[]byte) (err error) ***REMOVED***
	for _, s := range slices ***REMOVED***
		_, err = out.Write(s)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// lineBreaker breaks data across several lines, all of the same byte length
// (except possibly the last). Lines are broken with a single '\n'.
type lineBreaker struct ***REMOVED***
	lineLength  int
	line        []byte
	used        int
	out         io.Writer
	haveWritten bool
***REMOVED***

func newLineBreaker(out io.Writer, lineLength int) *lineBreaker ***REMOVED***
	return &lineBreaker***REMOVED***
		lineLength: lineLength,
		line:       make([]byte, lineLength),
		used:       0,
		out:        out,
	***REMOVED***
***REMOVED***

func (l *lineBreaker) Write(b []byte) (n int, err error) ***REMOVED***
	n = len(b)

	if n == 0 ***REMOVED***
		return
	***REMOVED***

	if l.used == 0 && l.haveWritten ***REMOVED***
		_, err = l.out.Write([]byte***REMOVED***'\n'***REMOVED***)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if l.used+len(b) < l.lineLength ***REMOVED***
		l.used += copy(l.line[l.used:], b)
		return
	***REMOVED***

	l.haveWritten = true
	_, err = l.out.Write(l.line[0:l.used])
	if err != nil ***REMOVED***
		return
	***REMOVED***
	excess := l.lineLength - l.used
	l.used = 0

	_, err = l.out.Write(b[0:excess])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	_, err = l.Write(b[excess:])
	return
***REMOVED***

func (l *lineBreaker) Close() (err error) ***REMOVED***
	if l.used > 0 ***REMOVED***
		_, err = l.out.Write(l.line[0:l.used])
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

// encoding keeps track of a running CRC24 over the data which has been written
// to it and outputs a OpenPGP checksum when closed, followed by an armor
// trailer.
//
// It's built into a stack of io.Writers:
//    encoding -> base64 encoder -> lineBreaker -> out
type encoding struct ***REMOVED***
	out       io.Writer
	breaker   *lineBreaker
	b64       io.WriteCloser
	crc       uint32
	blockType []byte
***REMOVED***

func (e *encoding) Write(data []byte) (n int, err error) ***REMOVED***
	e.crc = crc24(e.crc, data)
	return e.b64.Write(data)
***REMOVED***

func (e *encoding) Close() (err error) ***REMOVED***
	err = e.b64.Close()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	e.breaker.Close()

	var checksumBytes [3]byte
	checksumBytes[0] = byte(e.crc >> 16)
	checksumBytes[1] = byte(e.crc >> 8)
	checksumBytes[2] = byte(e.crc)

	var b64ChecksumBytes [4]byte
	base64.StdEncoding.Encode(b64ChecksumBytes[:], checksumBytes[:])

	return writeSlices(e.out, blockEnd, b64ChecksumBytes[:], newline, armorEnd, e.blockType, armorEndOfLine)
***REMOVED***

// Encode returns a WriteCloser which will encode the data written to it in
// OpenPGP armor.
func Encode(out io.Writer, blockType string, headers map[string]string) (w io.WriteCloser, err error) ***REMOVED***
	bType := []byte(blockType)
	err = writeSlices(out, armorStart, bType, armorEndOfLineOut)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	for k, v := range headers ***REMOVED***
		err = writeSlices(out, []byte(k), armorHeaderSep, []byte(v), newline)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	_, err = out.Write(newline)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	e := &encoding***REMOVED***
		out:       out,
		breaker:   newLineBreaker(out, 64),
		crc:       crc24Init,
		blockType: bType,
	***REMOVED***
	e.b64 = base64.NewEncoder(base64.StdEncoding, e.breaker)
	return e, nil
***REMOVED***
