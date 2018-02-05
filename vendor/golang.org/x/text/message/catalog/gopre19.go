// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.9

package catalog

import "golang.org/x/text/internal/catmsg"

// A Message holds a collection of translations for the same phrase that may
// vary based on the values of substitution arguments.
type Message interface ***REMOVED***
	catmsg.Message
***REMOVED***

func firstInSequence(m []Message) catmsg.Message ***REMOVED***
	a := []catmsg.Message***REMOVED******REMOVED***
	for _, m := range m ***REMOVED***
		a = append(a, m)
	***REMOVED***
	return catmsg.FirstOf(a)
***REMOVED***
