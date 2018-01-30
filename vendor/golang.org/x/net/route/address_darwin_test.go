// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"reflect"
	"testing"
)

type parseAddrsOnDarwinTest struct ***REMOVED***
	attrs uint
	fn    func(int, []byte) (int, Addr, error)
	b     []byte
	as    []Addr
***REMOVED***

var parseAddrsOnDarwinLittleEndianTests = []parseAddrsOnDarwinTest***REMOVED***
	***REMOVED***
		sysRTA_DST | sysRTA_GATEWAY | sysRTA_NETMASK,
		parseKernelInetAddr,
		[]byte***REMOVED***
			0x10, 0x2, 0x0, 0x0, 0xc0, 0xa8, 0x56, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,

			0x14, 0x12, 0x4, 0x0, 0x6, 0x0, 0x0, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x0, 0x0, 0x0, 0x0,

			0x7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		***REMOVED***,
		[]Addr***REMOVED***
			&Inet4Addr***REMOVED***IP: [4]byte***REMOVED***192, 168, 86, 0***REMOVED******REMOVED***,
			&LinkAddr***REMOVED***Index: 4***REMOVED***,
			&Inet4Addr***REMOVED***IP: [4]byte***REMOVED***255, 255, 255, 255***REMOVED******REMOVED***,
			nil,
			nil,
			nil,
			nil,
			nil,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestParseAddrsOnDarwin(t *testing.T) ***REMOVED***
	tests := parseAddrsOnDarwinLittleEndianTests
	if nativeEndian != littleEndian ***REMOVED***
		t.Skip("no test for non-little endian machine yet")
	***REMOVED***

	for i, tt := range tests ***REMOVED***
		as, err := parseAddrs(tt.attrs, tt.fn, tt.b)
		if err != nil ***REMOVED***
			t.Error(i, err)
			continue
		***REMOVED***
		if !reflect.DeepEqual(as, tt.as) ***REMOVED***
			t.Errorf("#%d: got %+v; want %+v", i, as, tt.as)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***
