// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dnsmessage_test

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

func mustNewName(name string) dnsmessage.Name ***REMOVED***
	n, err := dnsmessage.NewName(name)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return n
***REMOVED***

func ExampleParser() ***REMOVED***
	msg := dnsmessage.Message***REMOVED***
		Header: dnsmessage.Header***REMOVED***Response: true, Authoritative: true***REMOVED***,
		Questions: []dnsmessage.Question***REMOVED***
			***REMOVED***
				Name:  mustNewName("foo.bar.example.com."),
				Type:  dnsmessage.TypeA,
				Class: dnsmessage.ClassINET,
			***REMOVED***,
			***REMOVED***
				Name:  mustNewName("bar.example.com."),
				Type:  dnsmessage.TypeA,
				Class: dnsmessage.ClassINET,
			***REMOVED***,
		***REMOVED***,
		Answers: []dnsmessage.Resource***REMOVED***
			***REMOVED***
				dnsmessage.ResourceHeader***REMOVED***
					Name:  mustNewName("foo.bar.example.com."),
					Type:  dnsmessage.TypeA,
					Class: dnsmessage.ClassINET,
				***REMOVED***,
				&dnsmessage.AResource***REMOVED***[4]byte***REMOVED***127, 0, 0, 1***REMOVED******REMOVED***,
			***REMOVED***,
			***REMOVED***
				dnsmessage.ResourceHeader***REMOVED***
					Name:  mustNewName("bar.example.com."),
					Type:  dnsmessage.TypeA,
					Class: dnsmessage.ClassINET,
				***REMOVED***,
				&dnsmessage.AResource***REMOVED***[4]byte***REMOVED***127, 0, 0, 2***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	buf, err := msg.Pack()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	wantName := "bar.example.com."

	var p dnsmessage.Parser
	if _, err := p.Start(buf); err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	for ***REMOVED***
		q, err := p.Question()
		if err == dnsmessage.ErrSectionDone ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		if q.Name.String() != wantName ***REMOVED***
			continue
		***REMOVED***

		fmt.Println("Found question for name", wantName)
		if err := p.SkipAllQuestions(); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		break
	***REMOVED***

	var gotIPs []net.IP
	for ***REMOVED***
		h, err := p.AnswerHeader()
		if err == dnsmessage.ErrSectionDone ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		if (h.Type != dnsmessage.TypeA && h.Type != dnsmessage.TypeAAAA) || h.Class != dnsmessage.ClassINET ***REMOVED***
			continue
		***REMOVED***

		if !strings.EqualFold(h.Name.String(), wantName) ***REMOVED***
			if err := p.SkipAnswer(); err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			continue
		***REMOVED***

		switch h.Type ***REMOVED***
		case dnsmessage.TypeA:
			r, err := p.AResource()
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			gotIPs = append(gotIPs, r.A[:])
		case dnsmessage.TypeAAAA:
			r, err := p.AAAAResource()
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			gotIPs = append(gotIPs, r.AAAA[:])
		***REMOVED***
	***REMOVED***

	fmt.Printf("Found A/AAAA records for name %s: %v\n", wantName, gotIPs)

	// Output:
	// Found question for name bar.example.com.
	// Found A/AAAA records for name bar.example.com.: [127.0.0.2]
***REMOVED***
