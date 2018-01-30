// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!netbsd

package socket

import "net"

type mmsghdr struct***REMOVED******REMOVED***

type mmsghdrs []mmsghdr

func (hs mmsghdrs) pack(ms []Message, parseFn func([]byte, string) (net.Addr, error), marshalFn func(net.Addr) []byte) error ***REMOVED***
	return nil
***REMOVED***

func (hs mmsghdrs) unpack(ms []Message, parseFn func([]byte, string) (net.Addr, error), hint string) error ***REMOVED***
	return nil
***REMOVED***