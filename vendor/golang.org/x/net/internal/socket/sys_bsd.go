// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd openbsd

package socket

import "errors"

func recvmmsg(s uintptr, hs []mmsghdr, flags int) (int, error) ***REMOVED***
	return 0, errors.New("not implemented")
***REMOVED***

func sendmmsg(s uintptr, hs []mmsghdr, flags int) (int, error) ***REMOVED***
	return 0, errors.New("not implemented")
***REMOVED***
