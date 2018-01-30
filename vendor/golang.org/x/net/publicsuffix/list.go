// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go

// Package publicsuffix provides a public suffix list based on data from
// http://publicsuffix.org/. A public suffix is one under which Internet users
// can directly register names.
package publicsuffix // import "golang.org/x/net/publicsuffix"

// TODO: specify case sensitivity and leading/trailing dot behavior for
// func PublicSuffix and func EffectiveTLDPlusOne.

import (
	"fmt"
	"net/http/cookiejar"
	"strings"
)

// List implements the cookiejar.PublicSuffixList interface by calling the
// PublicSuffix function.
var List cookiejar.PublicSuffixList = list***REMOVED******REMOVED***

type list struct***REMOVED******REMOVED***

func (list) PublicSuffix(domain string) string ***REMOVED***
	ps, _ := PublicSuffix(domain)
	return ps
***REMOVED***

func (list) String() string ***REMOVED***
	return version
***REMOVED***

// PublicSuffix returns the public suffix of the domain using a copy of the
// publicsuffix.org database compiled into the library.
//
// icann is whether the public suffix is managed by the Internet Corporation
// for Assigned Names and Numbers. If not, the public suffix is privately
// managed. For example, foo.org and foo.co.uk are ICANN domains,
// foo.dyndns.org and foo.blogspot.co.uk are private domains.
//
// Use cases for distinguishing ICANN domains like foo.com from private
// domains like foo.appspot.com can be found at
// https://wiki.mozilla.org/Public_Suffix_List/Use_Cases
func PublicSuffix(domain string) (publicSuffix string, icann bool) ***REMOVED***
	lo, hi := uint32(0), uint32(numTLD)
	s, suffix, wildcard := domain, len(domain), false
loop:
	for ***REMOVED***
		dot := strings.LastIndex(s, ".")
		if wildcard ***REMOVED***
			suffix = 1 + dot
		***REMOVED***
		if lo == hi ***REMOVED***
			break
		***REMOVED***
		f := find(s[1+dot:], lo, hi)
		if f == notFound ***REMOVED***
			break
		***REMOVED***

		u := nodes[f] >> (nodesBitsTextOffset + nodesBitsTextLength)
		icann = u&(1<<nodesBitsICANN-1) != 0
		u >>= nodesBitsICANN
		u = children[u&(1<<nodesBitsChildren-1)]
		lo = u & (1<<childrenBitsLo - 1)
		u >>= childrenBitsLo
		hi = u & (1<<childrenBitsHi - 1)
		u >>= childrenBitsHi
		switch u & (1<<childrenBitsNodeType - 1) ***REMOVED***
		case nodeTypeNormal:
			suffix = 1 + dot
		case nodeTypeException:
			suffix = 1 + len(s)
			break loop
		***REMOVED***
		u >>= childrenBitsNodeType
		wildcard = u&(1<<childrenBitsWildcard-1) != 0

		if dot == -1 ***REMOVED***
			break
		***REMOVED***
		s = s[:dot]
	***REMOVED***
	if suffix == len(domain) ***REMOVED***
		// If no rules match, the prevailing rule is "*".
		return domain[1+strings.LastIndex(domain, "."):], icann
	***REMOVED***
	return domain[suffix:], icann
***REMOVED***

const notFound uint32 = 1<<32 - 1

// find returns the index of the node in the range [lo, hi) whose label equals
// label, or notFound if there is no such node. The range is assumed to be in
// strictly increasing node label order.
func find(label string, lo, hi uint32) uint32 ***REMOVED***
	for lo < hi ***REMOVED***
		mid := lo + (hi-lo)/2
		s := nodeLabel(mid)
		if s < label ***REMOVED***
			lo = mid + 1
		***REMOVED*** else if s == label ***REMOVED***
			return mid
		***REMOVED*** else ***REMOVED***
			hi = mid
		***REMOVED***
	***REMOVED***
	return notFound
***REMOVED***

// nodeLabel returns the label for the i'th node.
func nodeLabel(i uint32) string ***REMOVED***
	x := nodes[i]
	length := x & (1<<nodesBitsTextLength - 1)
	x >>= nodesBitsTextLength
	offset := x & (1<<nodesBitsTextOffset - 1)
	return text[offset : offset+length]
***REMOVED***

// EffectiveTLDPlusOne returns the effective top level domain plus one more
// label. For example, the eTLD+1 for "foo.bar.golang.org" is "golang.org".
func EffectiveTLDPlusOne(domain string) (string, error) ***REMOVED***
	suffix, _ := PublicSuffix(domain)
	if len(domain) <= len(suffix) ***REMOVED***
		return "", fmt.Errorf("publicsuffix: cannot derive eTLD+1 for domain %q", domain)
	***REMOVED***
	i := len(domain) - len(suffix) - 1
	if domain[i] != '.' ***REMOVED***
		return "", fmt.Errorf("publicsuffix: invalid public suffix %q for domain %q", suffix, domain)
	***REMOVED***
	return domain[1+strings.LastIndex(domain[:i], "."):], nil
***REMOVED***
