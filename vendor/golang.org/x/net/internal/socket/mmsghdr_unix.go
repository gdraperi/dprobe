// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux netbsd

package socket

import "net"

type mmsghdrs []mmsghdr

func (hs mmsghdrs) pack(ms []Message, parseFn func([]byte, string) (net.Addr, error), marshalFn func(net.Addr) []byte) error ***REMOVED***
	for i := range hs ***REMOVED***
		vs := make([]iovec, len(ms[i].Buffers))
		var sa []byte
		if parseFn != nil ***REMOVED***
			sa = make([]byte, sizeofSockaddrInet6)
		***REMOVED***
		if marshalFn != nil ***REMOVED***
			sa = marshalFn(ms[i].Addr)
		***REMOVED***
		hs[i].Hdr.pack(vs, ms[i].Buffers, ms[i].OOB, sa)
	***REMOVED***
	return nil
***REMOVED***

func (hs mmsghdrs) unpack(ms []Message, parseFn func([]byte, string) (net.Addr, error), hint string) error ***REMOVED***
	for i := range hs ***REMOVED***
		ms[i].N = int(hs[i].Len)
		ms[i].NN = hs[i].Hdr.controllen()
		ms[i].Flags = hs[i].Hdr.flags()
		if parseFn != nil ***REMOVED***
			var err error
			ms[i].Addr, err = parseFn(hs[i].Hdr.name(), hint)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
