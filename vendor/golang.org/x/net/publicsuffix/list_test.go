// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package publicsuffix

import (
	"sort"
	"strings"
	"testing"
)

func TestNodeLabel(t *testing.T) ***REMOVED***
	for i, want := range nodeLabels ***REMOVED***
		got := nodeLabel(uint32(i))
		if got != want ***REMOVED***
			t.Errorf("%d: got %q, want %q", i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFind(t *testing.T) ***REMOVED***
	testCases := []string***REMOVED***
		"",
		"a",
		"a0",
		"aaaa",
		"ao",
		"ap",
		"ar",
		"aro",
		"arp",
		"arpa",
		"arpaa",
		"arpb",
		"az",
		"b",
		"b0",
		"ba",
		"z",
		"zu",
		"zv",
		"zw",
		"zx",
		"zy",
		"zz",
		"zzzz",
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		got := find(tc, 0, numTLD)
		want := notFound
		for i := uint32(0); i < numTLD; i++ ***REMOVED***
			if tc == nodeLabel(i) ***REMOVED***
				want = i
				break
			***REMOVED***
		***REMOVED***
		if got != want ***REMOVED***
			t.Errorf("%q: got %d, want %d", tc, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestICANN(t *testing.T) ***REMOVED***
	testCases := map[string]bool***REMOVED***
		"foo.org":            true,
		"foo.co.uk":          true,
		"foo.dyndns.org":     false,
		"foo.go.dyndns.org":  false,
		"foo.blogspot.co.uk": false,
		"foo.intranet":       false,
	***REMOVED***
	for domain, want := range testCases ***REMOVED***
		_, got := PublicSuffix(domain)
		if got != want ***REMOVED***
			t.Errorf("%q: got %v, want %v", domain, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

var publicSuffixTestCases = []struct ***REMOVED***
	domain, want string
***REMOVED******REMOVED***
	// Empty string.
	***REMOVED***"", ""***REMOVED***,

	// The .ao rules are:
	// ao
	// ed.ao
	// gv.ao
	// og.ao
	// co.ao
	// pb.ao
	// it.ao
	***REMOVED***"ao", "ao"***REMOVED***,
	***REMOVED***"www.ao", "ao"***REMOVED***,
	***REMOVED***"pb.ao", "pb.ao"***REMOVED***,
	***REMOVED***"www.pb.ao", "pb.ao"***REMOVED***,
	***REMOVED***"www.xxx.yyy.zzz.pb.ao", "pb.ao"***REMOVED***,

	// The .ar rules are:
	// ar
	// com.ar
	// edu.ar
	// gob.ar
	// gov.ar
	// int.ar
	// mil.ar
	// net.ar
	// org.ar
	// tur.ar
	// blogspot.com.ar
	***REMOVED***"ar", "ar"***REMOVED***,
	***REMOVED***"www.ar", "ar"***REMOVED***,
	***REMOVED***"nic.ar", "ar"***REMOVED***,
	***REMOVED***"www.nic.ar", "ar"***REMOVED***,
	***REMOVED***"com.ar", "com.ar"***REMOVED***,
	***REMOVED***"www.com.ar", "com.ar"***REMOVED***,
	***REMOVED***"blogspot.com.ar", "blogspot.com.ar"***REMOVED***,
	***REMOVED***"www.blogspot.com.ar", "blogspot.com.ar"***REMOVED***,
	***REMOVED***"www.xxx.yyy.zzz.blogspot.com.ar", "blogspot.com.ar"***REMOVED***,
	***REMOVED***"logspot.com.ar", "com.ar"***REMOVED***,
	***REMOVED***"zlogspot.com.ar", "com.ar"***REMOVED***,
	***REMOVED***"zblogspot.com.ar", "com.ar"***REMOVED***,

	// The .arpa rules are:
	// arpa
	// e164.arpa
	// in-addr.arpa
	// ip6.arpa
	// iris.arpa
	// uri.arpa
	// urn.arpa
	***REMOVED***"arpa", "arpa"***REMOVED***,
	***REMOVED***"www.arpa", "arpa"***REMOVED***,
	***REMOVED***"urn.arpa", "urn.arpa"***REMOVED***,
	***REMOVED***"www.urn.arpa", "urn.arpa"***REMOVED***,
	***REMOVED***"www.xxx.yyy.zzz.urn.arpa", "urn.arpa"***REMOVED***,

	// The relevant ***REMOVED***kobe,kyoto***REMOVED***.jp rules are:
	// jp
	// *.kobe.jp
	// !city.kobe.jp
	// kyoto.jp
	// ide.kyoto.jp
	***REMOVED***"jp", "jp"***REMOVED***,
	***REMOVED***"kobe.jp", "jp"***REMOVED***,
	***REMOVED***"c.kobe.jp", "c.kobe.jp"***REMOVED***,
	***REMOVED***"b.c.kobe.jp", "c.kobe.jp"***REMOVED***,
	***REMOVED***"a.b.c.kobe.jp", "c.kobe.jp"***REMOVED***,
	***REMOVED***"city.kobe.jp", "kobe.jp"***REMOVED***,
	***REMOVED***"www.city.kobe.jp", "kobe.jp"***REMOVED***,
	***REMOVED***"kyoto.jp", "kyoto.jp"***REMOVED***,
	***REMOVED***"test.kyoto.jp", "kyoto.jp"***REMOVED***,
	***REMOVED***"ide.kyoto.jp", "ide.kyoto.jp"***REMOVED***,
	***REMOVED***"b.ide.kyoto.jp", "ide.kyoto.jp"***REMOVED***,
	***REMOVED***"a.b.ide.kyoto.jp", "ide.kyoto.jp"***REMOVED***,

	// The .tw rules are:
	// tw
	// edu.tw
	// gov.tw
	// mil.tw
	// com.tw
	// net.tw
	// org.tw
	// idv.tw
	// game.tw
	// ebiz.tw
	// club.tw
	// 網路.tw (xn--zf0ao64a.tw)
	// 組織.tw (xn--uc0atv.tw)
	// 商業.tw (xn--czrw28b.tw)
	// blogspot.tw
	***REMOVED***"tw", "tw"***REMOVED***,
	***REMOVED***"aaa.tw", "tw"***REMOVED***,
	***REMOVED***"www.aaa.tw", "tw"***REMOVED***,
	***REMOVED***"xn--czrw28b.aaa.tw", "tw"***REMOVED***,
	***REMOVED***"edu.tw", "edu.tw"***REMOVED***,
	***REMOVED***"www.edu.tw", "edu.tw"***REMOVED***,
	***REMOVED***"xn--czrw28b.edu.tw", "edu.tw"***REMOVED***,
	***REMOVED***"xn--czrw28b.tw", "xn--czrw28b.tw"***REMOVED***,
	***REMOVED***"www.xn--czrw28b.tw", "xn--czrw28b.tw"***REMOVED***,
	***REMOVED***"xn--uc0atv.xn--czrw28b.tw", "xn--czrw28b.tw"***REMOVED***,
	***REMOVED***"xn--kpry57d.tw", "tw"***REMOVED***,

	// The .uk rules are:
	// uk
	// ac.uk
	// co.uk
	// gov.uk
	// ltd.uk
	// me.uk
	// net.uk
	// nhs.uk
	// org.uk
	// plc.uk
	// police.uk
	// *.sch.uk
	// blogspot.co.uk
	***REMOVED***"uk", "uk"***REMOVED***,
	***REMOVED***"aaa.uk", "uk"***REMOVED***,
	***REMOVED***"www.aaa.uk", "uk"***REMOVED***,
	***REMOVED***"mod.uk", "uk"***REMOVED***,
	***REMOVED***"www.mod.uk", "uk"***REMOVED***,
	***REMOVED***"sch.uk", "uk"***REMOVED***,
	***REMOVED***"mod.sch.uk", "mod.sch.uk"***REMOVED***,
	***REMOVED***"www.sch.uk", "www.sch.uk"***REMOVED***,
	***REMOVED***"blogspot.co.uk", "blogspot.co.uk"***REMOVED***,
	***REMOVED***"blogspot.nic.uk", "uk"***REMOVED***,
	***REMOVED***"blogspot.sch.uk", "blogspot.sch.uk"***REMOVED***,

	// The .рф rules are
	// рф (xn--p1ai)
	***REMOVED***"xn--p1ai", "xn--p1ai"***REMOVED***,
	***REMOVED***"aaa.xn--p1ai", "xn--p1ai"***REMOVED***,
	***REMOVED***"www.xxx.yyy.xn--p1ai", "xn--p1ai"***REMOVED***,

	// The .bd rules are:
	// *.bd
	***REMOVED***"bd", "bd"***REMOVED***,
	***REMOVED***"www.bd", "www.bd"***REMOVED***,
	***REMOVED***"zzz.bd", "zzz.bd"***REMOVED***,
	***REMOVED***"www.zzz.bd", "zzz.bd"***REMOVED***,
	***REMOVED***"www.xxx.yyy.zzz.bd", "zzz.bd"***REMOVED***,

	// There are no .nosuchtld rules.
	***REMOVED***"nosuchtld", "nosuchtld"***REMOVED***,
	***REMOVED***"foo.nosuchtld", "nosuchtld"***REMOVED***,
	***REMOVED***"bar.foo.nosuchtld", "nosuchtld"***REMOVED***,
***REMOVED***

func BenchmarkPublicSuffix(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, tc := range publicSuffixTestCases ***REMOVED***
			List.PublicSuffix(tc.domain)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPublicSuffix(t *testing.T) ***REMOVED***
	for _, tc := range publicSuffixTestCases ***REMOVED***
		got := List.PublicSuffix(tc.domain)
		if got != tc.want ***REMOVED***
			t.Errorf("%q: got %q, want %q", tc.domain, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSlowPublicSuffix(t *testing.T) ***REMOVED***
	for _, tc := range publicSuffixTestCases ***REMOVED***
		got := slowPublicSuffix(tc.domain)
		if got != tc.want ***REMOVED***
			t.Errorf("%q: got %q, want %q", tc.domain, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// slowPublicSuffix implements the canonical (but O(number of rules)) public
// suffix algorithm described at http://publicsuffix.org/list/.
//
// 1. Match domain against all rules and take note of the matching ones.
// 2. If no rules match, the prevailing rule is "*".
// 3. If more than one rule matches, the prevailing rule is the one which is an exception rule.
// 4. If there is no matching exception rule, the prevailing rule is the one with the most labels.
// 5. If the prevailing rule is a exception rule, modify it by removing the leftmost label.
// 6. The public suffix is the set of labels from the domain which directly match the labels of the prevailing rule (joined by dots).
// 7. The registered or registrable domain is the public suffix plus one additional label.
//
// This function returns the public suffix, not the registrable domain, and so
// it stops after step 6.
func slowPublicSuffix(domain string) string ***REMOVED***
	match := func(rulePart, domainPart string) bool ***REMOVED***
		switch rulePart[0] ***REMOVED***
		case '*':
			return true
		case '!':
			return rulePart[1:] == domainPart
		***REMOVED***
		return rulePart == domainPart
	***REMOVED***

	domainParts := strings.Split(domain, ".")
	var matchingRules [][]string

loop:
	for _, rule := range rules ***REMOVED***
		ruleParts := strings.Split(rule, ".")
		if len(domainParts) < len(ruleParts) ***REMOVED***
			continue
		***REMOVED***
		for i := range ruleParts ***REMOVED***
			rulePart := ruleParts[len(ruleParts)-1-i]
			domainPart := domainParts[len(domainParts)-1-i]
			if !match(rulePart, domainPart) ***REMOVED***
				continue loop
			***REMOVED***
		***REMOVED***
		matchingRules = append(matchingRules, ruleParts)
	***REMOVED***
	if len(matchingRules) == 0 ***REMOVED***
		matchingRules = append(matchingRules, []string***REMOVED***"*"***REMOVED***)
	***REMOVED*** else ***REMOVED***
		sort.Sort(byPriority(matchingRules))
	***REMOVED***
	prevailing := matchingRules[0]
	if prevailing[0][0] == '!' ***REMOVED***
		prevailing = prevailing[1:]
	***REMOVED***
	if prevailing[0][0] == '*' ***REMOVED***
		replaced := domainParts[len(domainParts)-len(prevailing)]
		prevailing = append([]string***REMOVED***replaced***REMOVED***, prevailing[1:]...)
	***REMOVED***
	return strings.Join(prevailing, ".")
***REMOVED***

type byPriority [][]string

func (b byPriority) Len() int      ***REMOVED*** return len(b) ***REMOVED***
func (b byPriority) Swap(i, j int) ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b byPriority) Less(i, j int) bool ***REMOVED***
	if b[i][0][0] == '!' ***REMOVED***
		return true
	***REMOVED***
	if b[j][0][0] == '!' ***REMOVED***
		return false
	***REMOVED***
	return len(b[i]) > len(b[j])
***REMOVED***

// eTLDPlusOneTestCases come from
// https://github.com/publicsuffix/list/blob/master/tests/test_psl.txt
var eTLDPlusOneTestCases = []struct ***REMOVED***
	domain, want string
***REMOVED******REMOVED***
	// Empty input.
	***REMOVED***"", ""***REMOVED***,
	// Unlisted TLD.
	***REMOVED***"example", ""***REMOVED***,
	***REMOVED***"example.example", "example.example"***REMOVED***,
	***REMOVED***"b.example.example", "example.example"***REMOVED***,
	***REMOVED***"a.b.example.example", "example.example"***REMOVED***,
	// TLD with only 1 rule.
	***REMOVED***"biz", ""***REMOVED***,
	***REMOVED***"domain.biz", "domain.biz"***REMOVED***,
	***REMOVED***"b.domain.biz", "domain.biz"***REMOVED***,
	***REMOVED***"a.b.domain.biz", "domain.biz"***REMOVED***,
	// TLD with some 2-level rules.
	***REMOVED***"com", ""***REMOVED***,
	***REMOVED***"example.com", "example.com"***REMOVED***,
	***REMOVED***"b.example.com", "example.com"***REMOVED***,
	***REMOVED***"a.b.example.com", "example.com"***REMOVED***,
	***REMOVED***"uk.com", ""***REMOVED***,
	***REMOVED***"example.uk.com", "example.uk.com"***REMOVED***,
	***REMOVED***"b.example.uk.com", "example.uk.com"***REMOVED***,
	***REMOVED***"a.b.example.uk.com", "example.uk.com"***REMOVED***,
	***REMOVED***"test.ac", "test.ac"***REMOVED***,
	// TLD with only 1 (wildcard) rule.
	***REMOVED***"mm", ""***REMOVED***,
	***REMOVED***"c.mm", ""***REMOVED***,
	***REMOVED***"b.c.mm", "b.c.mm"***REMOVED***,
	***REMOVED***"a.b.c.mm", "b.c.mm"***REMOVED***,
	// More complex TLD.
	***REMOVED***"jp", ""***REMOVED***,
	***REMOVED***"test.jp", "test.jp"***REMOVED***,
	***REMOVED***"www.test.jp", "test.jp"***REMOVED***,
	***REMOVED***"ac.jp", ""***REMOVED***,
	***REMOVED***"test.ac.jp", "test.ac.jp"***REMOVED***,
	***REMOVED***"www.test.ac.jp", "test.ac.jp"***REMOVED***,
	***REMOVED***"kyoto.jp", ""***REMOVED***,
	***REMOVED***"test.kyoto.jp", "test.kyoto.jp"***REMOVED***,
	***REMOVED***"ide.kyoto.jp", ""***REMOVED***,
	***REMOVED***"b.ide.kyoto.jp", "b.ide.kyoto.jp"***REMOVED***,
	***REMOVED***"a.b.ide.kyoto.jp", "b.ide.kyoto.jp"***REMOVED***,
	***REMOVED***"c.kobe.jp", ""***REMOVED***,
	***REMOVED***"b.c.kobe.jp", "b.c.kobe.jp"***REMOVED***,
	***REMOVED***"a.b.c.kobe.jp", "b.c.kobe.jp"***REMOVED***,
	***REMOVED***"city.kobe.jp", "city.kobe.jp"***REMOVED***,
	***REMOVED***"www.city.kobe.jp", "city.kobe.jp"***REMOVED***,
	// TLD with a wildcard rule and exceptions.
	***REMOVED***"ck", ""***REMOVED***,
	***REMOVED***"test.ck", ""***REMOVED***,
	***REMOVED***"b.test.ck", "b.test.ck"***REMOVED***,
	***REMOVED***"a.b.test.ck", "b.test.ck"***REMOVED***,
	***REMOVED***"www.ck", "www.ck"***REMOVED***,
	***REMOVED***"www.www.ck", "www.ck"***REMOVED***,
	// US K12.
	***REMOVED***"us", ""***REMOVED***,
	***REMOVED***"test.us", "test.us"***REMOVED***,
	***REMOVED***"www.test.us", "test.us"***REMOVED***,
	***REMOVED***"ak.us", ""***REMOVED***,
	***REMOVED***"test.ak.us", "test.ak.us"***REMOVED***,
	***REMOVED***"www.test.ak.us", "test.ak.us"***REMOVED***,
	***REMOVED***"k12.ak.us", ""***REMOVED***,
	***REMOVED***"test.k12.ak.us", "test.k12.ak.us"***REMOVED***,
	***REMOVED***"www.test.k12.ak.us", "test.k12.ak.us"***REMOVED***,
	// Punycoded IDN labels
	***REMOVED***"xn--85x722f.com.cn", "xn--85x722f.com.cn"***REMOVED***,
	***REMOVED***"xn--85x722f.xn--55qx5d.cn", "xn--85x722f.xn--55qx5d.cn"***REMOVED***,
	***REMOVED***"www.xn--85x722f.xn--55qx5d.cn", "xn--85x722f.xn--55qx5d.cn"***REMOVED***,
	***REMOVED***"shishi.xn--55qx5d.cn", "shishi.xn--55qx5d.cn"***REMOVED***,
	***REMOVED***"xn--55qx5d.cn", ""***REMOVED***,
	***REMOVED***"xn--85x722f.xn--fiqs8s", "xn--85x722f.xn--fiqs8s"***REMOVED***,
	***REMOVED***"www.xn--85x722f.xn--fiqs8s", "xn--85x722f.xn--fiqs8s"***REMOVED***,
	***REMOVED***"shishi.xn--fiqs8s", "shishi.xn--fiqs8s"***REMOVED***,
	***REMOVED***"xn--fiqs8s", ""***REMOVED***,
***REMOVED***

func TestEffectiveTLDPlusOne(t *testing.T) ***REMOVED***
	for _, tc := range eTLDPlusOneTestCases ***REMOVED***
		got, _ := EffectiveTLDPlusOne(tc.domain)
		if got != tc.want ***REMOVED***
			t.Errorf("%q: got %q, want %q", tc.domain, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
