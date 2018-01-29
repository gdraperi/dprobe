// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!gccgo,!appengine

package argon2

func init() ***REMOVED***
	useSSE4 = supportsSSE4()
***REMOVED***

//go:noescape
func supportsSSE4() bool

//go:noescape
func mixBlocksSSE2(out, a, b, c *block)

//go:noescape
func xorBlocksSSE2(out, a, b, c *block)

//go:noescape
func blamkaSSE4(b *block)

func processBlockSSE(out, in1, in2 *block, xor bool) ***REMOVED***
	var t block
	mixBlocksSSE2(&t, in1, in2, &t)
	if useSSE4 ***REMOVED***
		blamkaSSE4(&t)
	***REMOVED*** else ***REMOVED***
		for i := 0; i < blockLength; i += 16 ***REMOVED***
			blamkaGeneric(
				&t[i+0], &t[i+1], &t[i+2], &t[i+3],
				&t[i+4], &t[i+5], &t[i+6], &t[i+7],
				&t[i+8], &t[i+9], &t[i+10], &t[i+11],
				&t[i+12], &t[i+13], &t[i+14], &t[i+15],
			)
		***REMOVED***
		for i := 0; i < blockLength/8; i += 2 ***REMOVED***
			blamkaGeneric(
				&t[i], &t[i+1], &t[16+i], &t[16+i+1],
				&t[32+i], &t[32+i+1], &t[48+i], &t[48+i+1],
				&t[64+i], &t[64+i+1], &t[80+i], &t[80+i+1],
				&t[96+i], &t[96+i+1], &t[112+i], &t[112+i+1],
			)
		***REMOVED***
	***REMOVED***
	if xor ***REMOVED***
		xorBlocksSSE2(out, in1, in2, &t)
	***REMOVED*** else ***REMOVED***
		mixBlocksSSE2(out, in1, in2, &t)
	***REMOVED***
***REMOVED***

func processBlock(out, in1, in2 *block) ***REMOVED***
	processBlockSSE(out, in1, in2, false)
***REMOVED***

func processBlockXOR(out, in1, in2 *block) ***REMOVED***
	processBlockSSE(out, in1, in2, true)
***REMOVED***
