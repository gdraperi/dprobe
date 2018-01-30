// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package lif

import (
	"fmt"
	"testing"
)

type addrFamily int

func (af addrFamily) String() string ***REMOVED***
	switch af ***REMOVED***
	case sysAF_UNSPEC:
		return "unspec"
	case sysAF_INET:
		return "inet4"
	case sysAF_INET6:
		return "inet6"
	default:
		return fmt.Sprintf("%d", af)
	***REMOVED***
***REMOVED***

const hexDigit = "0123456789abcdef"

type llAddr []byte

func (a llAddr) String() string ***REMOVED***
	if len(a) == 0 ***REMOVED***
		return ""
	***REMOVED***
	buf := make([]byte, 0, len(a)*3-1)
	for i, b := range a ***REMOVED***
		if i > 0 ***REMOVED***
			buf = append(buf, ':')
		***REMOVED***
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	***REMOVED***
	return string(buf)
***REMOVED***

type ipAddr []byte

func (a ipAddr) String() string ***REMOVED***
	if len(a) == 0 ***REMOVED***
		return "<nil>"
	***REMOVED***
	if len(a) == 4 ***REMOVED***
		return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
	***REMOVED***
	if len(a) == 16 ***REMOVED***
		return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x", a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14], a[15])
	***REMOVED***
	s := make([]byte, len(a)*2)
	for i, tn := range a ***REMOVED***
		s[i*2], s[i*2+1] = hexDigit[tn>>4], hexDigit[tn&0xf]
	***REMOVED***
	return string(s)
***REMOVED***

func (a *Inet4Addr) String() string ***REMOVED***
	return fmt.Sprintf("(%s %s %d)", addrFamily(a.Family()), ipAddr(a.IP[:]), a.PrefixLen)
***REMOVED***

func (a *Inet6Addr) String() string ***REMOVED***
	return fmt.Sprintf("(%s %s %d %d)", addrFamily(a.Family()), ipAddr(a.IP[:]), a.PrefixLen, a.ZoneID)
***REMOVED***

type addrPack struct ***REMOVED***
	af int
	as []Addr
***REMOVED***

func addrPacks() ([]addrPack, error) ***REMOVED***
	var lastErr error
	var aps []addrPack
	for _, af := range [...]int***REMOVED***sysAF_UNSPEC, sysAF_INET, sysAF_INET6***REMOVED*** ***REMOVED***
		as, err := Addrs(af, "")
		if err != nil ***REMOVED***
			lastErr = err
			continue
		***REMOVED***
		aps = append(aps, addrPack***REMOVED***af: af, as: as***REMOVED***)
	***REMOVED***
	return aps, lastErr
***REMOVED***

func TestAddrs(t *testing.T) ***REMOVED***
	aps, err := addrPacks()
	if len(aps) == 0 && err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	lps, err := linkPacks()
	if len(lps) == 0 && err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for _, lp := range lps ***REMOVED***
		n := 0
		for _, ll := range lp.lls ***REMOVED***
			as, err := Addrs(lp.af, ll.Name)
			if err != nil ***REMOVED***
				t.Fatal(lp.af, ll.Name, err)
			***REMOVED***
			t.Logf("af=%s name=%s %v", addrFamily(lp.af), ll.Name, as)
			n += len(as)
		***REMOVED***
		for _, ap := range aps ***REMOVED***
			if ap.af != lp.af ***REMOVED***
				continue
			***REMOVED***
			if n != len(ap.as) ***REMOVED***
				t.Errorf("af=%s got %d; want %d", addrFamily(lp.af), n, len(ap.as))
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
