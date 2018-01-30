// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build nacl plan9

package nettest

import (
	"fmt"
	"runtime"
)

func maxOpenFiles() int ***REMOVED***
	return defaultMaxOpenFiles
***REMOVED***

func supportsRawIPSocket() (string, bool) ***REMOVED***
	return fmt.Sprintf("not supported on %s", runtime.GOOS), false
***REMOVED***

func supportsIPv6MulticastDeliveryOnLoopback() bool ***REMOVED***
	return false
***REMOVED***

func causesIPv6Crash() bool ***REMOVED***
	return false
***REMOVED***

func protocolNotSupported(err error) bool ***REMOVED***
	return false
***REMOVED***
