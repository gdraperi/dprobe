// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package sortkeys

import (
	"sort"
)

func Strings(l []string) ***REMOVED***
	sort.Strings(l)
***REMOVED***

func Float64s(l []float64) ***REMOVED***
	sort.Float64s(l)
***REMOVED***

func Float32s(l []float32) ***REMOVED***
	sort.Sort(Float32Slice(l))
***REMOVED***

func Int64s(l []int64) ***REMOVED***
	sort.Sort(Int64Slice(l))
***REMOVED***

func Int32s(l []int32) ***REMOVED***
	sort.Sort(Int32Slice(l))
***REMOVED***

func Uint64s(l []uint64) ***REMOVED***
	sort.Sort(Uint64Slice(l))
***REMOVED***

func Uint32s(l []uint32) ***REMOVED***
	sort.Sort(Uint32Slice(l))
***REMOVED***

func Bools(l []bool) ***REMOVED***
	sort.Sort(BoolSlice(l))
***REMOVED***

type BoolSlice []bool

func (p BoolSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p BoolSlice) Less(i, j int) bool ***REMOVED*** return p[j] ***REMOVED***
func (p BoolSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

type Int64Slice []int64

func (p Int64Slice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p Int64Slice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p Int64Slice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

type Int32Slice []int32

func (p Int32Slice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p Int32Slice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p Int32Slice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p Uint64Slice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p Uint64Slice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

type Uint32Slice []uint32

func (p Uint32Slice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p Uint32Slice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p Uint32Slice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

type Float32Slice []float32

func (p Float32Slice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p Float32Slice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p Float32Slice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***
