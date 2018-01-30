// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package lif

import (
	"fmt"
	"testing"
)

func (ll *Link) String() string ***REMOVED***
	return fmt.Sprintf("name=%s index=%d type=%d flags=%#x mtu=%d addr=%v", ll.Name, ll.Index, ll.Type, ll.Flags, ll.MTU, llAddr(ll.Addr))
***REMOVED***

type linkPack struct ***REMOVED***
	af  int
	lls []Link
***REMOVED***

func linkPacks() ([]linkPack, error) ***REMOVED***
	var lastErr error
	var lps []linkPack
	for _, af := range [...]int***REMOVED***sysAF_UNSPEC, sysAF_INET, sysAF_INET6***REMOVED*** ***REMOVED***
		lls, err := Links(af, "")
		if err != nil ***REMOVED***
			lastErr = err
			continue
		***REMOVED***
		lps = append(lps, linkPack***REMOVED***af: af, lls: lls***REMOVED***)
	***REMOVED***
	return lps, lastErr
***REMOVED***

func TestLinks(t *testing.T) ***REMOVED***
	lps, err := linkPacks()
	if len(lps) == 0 && err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for _, lp := range lps ***REMOVED***
		n := 0
		for _, sll := range lp.lls ***REMOVED***
			lls, err := Links(lp.af, sll.Name)
			if err != nil ***REMOVED***
				t.Fatal(lp.af, sll.Name, err)
			***REMOVED***
			for _, ll := range lls ***REMOVED***
				if ll.Name != sll.Name || ll.Index != sll.Index ***REMOVED***
					t.Errorf("af=%s got %v; want %v", addrFamily(lp.af), &ll, &sll)
					continue
				***REMOVED***
				t.Logf("af=%s name=%s %v", addrFamily(lp.af), sll.Name, &ll)
				n++
			***REMOVED***
		***REMOVED***
		if n != len(lp.lls) ***REMOVED***
			t.Errorf("af=%s got %d; want %d", addrFamily(lp.af), n, len(lp.lls))
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***
