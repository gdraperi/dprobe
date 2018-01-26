// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// +build !debug

package diff

var debug debugger

type debugger struct***REMOVED******REMOVED***

func (debugger) Begin(_, _ int, f EqualFunc, _, _ *EditScript) EqualFunc ***REMOVED***
	return f
***REMOVED***
func (debugger) Update() ***REMOVED******REMOVED***
func (debugger) Finish() ***REMOVED******REMOVED***
