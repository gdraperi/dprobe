// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4_test

import (
	"testing"

	"golang.org/x/net/ipv4"
)

func TestControlMessageParseWithFuzz(t *testing.T) ***REMOVED***
	var cm ipv4.ControlMessage
	for _, fuzz := range []string***REMOVED***
		"\f\x00\x00\x00\x00\x00\x00\x00\x14\x00\x00\x00",
		"\f\x00\x00\x00\x00\x00\x00\x00\x1a\x00\x00\x00",
	***REMOVED*** ***REMOVED***
		cm.Parse([]byte(fuzz))
	***REMOVED***
***REMOVED***
