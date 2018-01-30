// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

// This program generates table.go and table_test.go based on the authoritative
// public suffix list at https://publicsuffix.org/list/effective_tld_names.dat
//
// The version is derived from
// https://api.github.com/repos/publicsuffix/list/commits?path=public_suffix_list.dat
// and a human-readable form is at
// https://github.com/publicsuffix/list/commits/master/public_suffix_list.dat
//
// To fetch a particular git revision, such as 5c70ccd250, pass
// -url "https://raw.githubusercontent.com/publicsuffix/list/5c70ccd250/public_suffix_list.dat"
// and -version "an explicit version string".

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/idna"
)

const (
	// These sum of these four values must be no greater than 32.
	nodesBitsChildren   = 10
	nodesBitsICANN      = 1
	nodesBitsTextOffset = 15
	nodesBitsTextLength = 6

	// These sum of these four values must be no greater than 32.
	childrenBitsWildcard = 1
	childrenBitsNodeType = 2
	childrenBitsHi       = 14
	childrenBitsLo       = 14
)

var (
	maxChildren   int
	maxTextOffset int
	maxTextLength int
	maxHi         uint32
	maxLo         uint32
)

func max(a, b int) int ***REMOVED***
	if a < b ***REMOVED***
		return b
	***REMOVED***
	return a
***REMOVED***

func u32max(a, b uint32) uint32 ***REMOVED***
	if a < b ***REMOVED***
		return b
	***REMOVED***
	return a
***REMOVED***

const (
	nodeTypeNormal     = 0
	nodeTypeException  = 1
	nodeTypeParentOnly = 2
	numNodeType        = 3
)

func nodeTypeStr(n int) string ***REMOVED***
	switch n ***REMOVED***
	case nodeTypeNormal:
		return "+"
	case nodeTypeException:
		return "!"
	case nodeTypeParentOnly:
		return "o"
	***REMOVED***
	panic("unreachable")
***REMOVED***

const (
	defaultURL   = "https://publicsuffix.org/list/effective_tld_names.dat"
	gitCommitURL = "https://api.github.com/repos/publicsuffix/list/commits?path=public_suffix_list.dat"
)

var (
	labelEncoding = map[string]uint32***REMOVED******REMOVED***
	labelsList    = []string***REMOVED******REMOVED***
	labelsMap     = map[string]bool***REMOVED******REMOVED***
	rules         = []string***REMOVED******REMOVED***

	// validSuffixRE is used to check that the entries in the public suffix
	// list are in canonical form (after Punycode encoding). Specifically,
	// capital letters are not allowed.
	validSuffixRE = regexp.MustCompile(`^[a-z0-9_\!\*\-\.]+$`)

	shaRE  = regexp.MustCompile(`"sha":"([^"]+)"`)
	dateRE = regexp.MustCompile(`"committer":***REMOVED***[^***REMOVED***]+"date":"([^"]+)"`)

	comments = flag.Bool("comments", false, "generate table.go comments, for debugging")
	subset   = flag.Bool("subset", false, "generate only a subset of the full table, for debugging")
	url      = flag.String("url", defaultURL, "URL of the publicsuffix.org list. If empty, stdin is read instead")
	v        = flag.Bool("v", false, "verbose output (to stderr)")
	version  = flag.String("version", "", "the effective_tld_names.dat version")
)

func main() ***REMOVED***
	if err := main1(); err != nil ***REMOVED***
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	***REMOVED***
***REMOVED***

func main1() error ***REMOVED***
	flag.Parse()
	if nodesBitsTextLength+nodesBitsTextOffset+nodesBitsICANN+nodesBitsChildren > 32 ***REMOVED***
		return fmt.Errorf("not enough bits to encode the nodes table")
	***REMOVED***
	if childrenBitsLo+childrenBitsHi+childrenBitsNodeType+childrenBitsWildcard > 32 ***REMOVED***
		return fmt.Errorf("not enough bits to encode the children table")
	***REMOVED***
	if *version == "" ***REMOVED***
		if *url != defaultURL ***REMOVED***
			return fmt.Errorf("-version was not specified, and the -url is not the default one")
		***REMOVED***
		sha, date, err := gitCommit()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*version = fmt.Sprintf("publicsuffix.org's public_suffix_list.dat, git revision %s (%s)", sha, date)
	***REMOVED***
	var r io.Reader = os.Stdin
	if *url != "" ***REMOVED***
		res, err := http.Get(*url)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if res.StatusCode != http.StatusOK ***REMOVED***
			return fmt.Errorf("bad GET status for %s: %d", *url, res.Status)
		***REMOVED***
		r = res.Body
		defer res.Body.Close()
	***REMOVED***

	var root node
	icann := false
	br := bufio.NewReader(r)
	for ***REMOVED***
		s, err := br.ReadString('\n')
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			return err
		***REMOVED***
		s = strings.TrimSpace(s)
		if strings.Contains(s, "BEGIN ICANN DOMAINS") ***REMOVED***
			icann = true
			continue
		***REMOVED***
		if strings.Contains(s, "END ICANN DOMAINS") ***REMOVED***
			icann = false
			continue
		***REMOVED***
		if s == "" || strings.HasPrefix(s, "//") ***REMOVED***
			continue
		***REMOVED***
		s, err = idna.ToASCII(s)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !validSuffixRE.MatchString(s) ***REMOVED***
			return fmt.Errorf("bad publicsuffix.org list data: %q", s)
		***REMOVED***

		if *subset ***REMOVED***
			switch ***REMOVED***
			case s == "ac.jp" || strings.HasSuffix(s, ".ac.jp"):
			case s == "ak.us" || strings.HasSuffix(s, ".ak.us"):
			case s == "ao" || strings.HasSuffix(s, ".ao"):
			case s == "ar" || strings.HasSuffix(s, ".ar"):
			case s == "arpa" || strings.HasSuffix(s, ".arpa"):
			case s == "cy" || strings.HasSuffix(s, ".cy"):
			case s == "dyndns.org" || strings.HasSuffix(s, ".dyndns.org"):
			case s == "jp":
			case s == "kobe.jp" || strings.HasSuffix(s, ".kobe.jp"):
			case s == "kyoto.jp" || strings.HasSuffix(s, ".kyoto.jp"):
			case s == "om" || strings.HasSuffix(s, ".om"):
			case s == "uk" || strings.HasSuffix(s, ".uk"):
			case s == "uk.com" || strings.HasSuffix(s, ".uk.com"):
			case s == "tw" || strings.HasSuffix(s, ".tw"):
			case s == "zw" || strings.HasSuffix(s, ".zw"):
			case s == "xn--p1ai" || strings.HasSuffix(s, ".xn--p1ai"):
				// xn--p1ai is Russian-Cyrillic "рф".
			default:
				continue
			***REMOVED***
		***REMOVED***

		rules = append(rules, s)

		nt, wildcard := nodeTypeNormal, false
		switch ***REMOVED***
		case strings.HasPrefix(s, "*."):
			s, nt = s[2:], nodeTypeParentOnly
			wildcard = true
		case strings.HasPrefix(s, "!"):
			s, nt = s[1:], nodeTypeException
		***REMOVED***
		labels := strings.Split(s, ".")
		for n, i := &root, len(labels)-1; i >= 0; i-- ***REMOVED***
			label := labels[i]
			n = n.child(label)
			if i == 0 ***REMOVED***
				if nt != nodeTypeParentOnly && n.nodeType == nodeTypeParentOnly ***REMOVED***
					n.nodeType = nt
				***REMOVED***
				n.icann = n.icann && icann
				n.wildcard = n.wildcard || wildcard
			***REMOVED***
			labelsMap[label] = true
		***REMOVED***
	***REMOVED***
	labelsList = make([]string, 0, len(labelsMap))
	for label := range labelsMap ***REMOVED***
		labelsList = append(labelsList, label)
	***REMOVED***
	sort.Strings(labelsList)

	if err := generate(printReal, &root, "table.go"); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := generate(printTest, &root, "table_test.go"); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func generate(p func(io.Writer, *node) error, root *node, filename string) error ***REMOVED***
	buf := new(bytes.Buffer)
	if err := p(buf, root); err != nil ***REMOVED***
		return err
	***REMOVED***
	b, err := format.Source(buf.Bytes())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(filename, b, 0644)
***REMOVED***

func gitCommit() (sha, date string, retErr error) ***REMOVED***
	res, err := http.Get(gitCommitURL)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	if res.StatusCode != http.StatusOK ***REMOVED***
		return "", "", fmt.Errorf("bad GET status for %s: %d", gitCommitURL, res.Status)
	***REMOVED***
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	if m := shaRE.FindSubmatch(b); m != nil ***REMOVED***
		sha = string(m[1])
	***REMOVED***
	if m := dateRE.FindSubmatch(b); m != nil ***REMOVED***
		date = string(m[1])
	***REMOVED***
	if sha == "" || date == "" ***REMOVED***
		retErr = fmt.Errorf("could not find commit SHA and date in %s", gitCommitURL)
	***REMOVED***
	return sha, date, retErr
***REMOVED***

func printTest(w io.Writer, n *node) error ***REMOVED***
	fmt.Fprintf(w, "// generated by go run gen.go; DO NOT EDIT\n\n")
	fmt.Fprintf(w, "package publicsuffix\n\nvar rules = [...]string***REMOVED***\n")
	for _, rule := range rules ***REMOVED***
		fmt.Fprintf(w, "%q,\n", rule)
	***REMOVED***
	fmt.Fprintf(w, "***REMOVED***\n\nvar nodeLabels = [...]string***REMOVED***\n")
	if err := n.walk(w, printNodeLabel); err != nil ***REMOVED***
		return err
	***REMOVED***
	fmt.Fprintf(w, "***REMOVED***\n")
	return nil
***REMOVED***

func printReal(w io.Writer, n *node) error ***REMOVED***
	const header = `// generated by go run gen.go; DO NOT EDIT

package publicsuffix

const version = %q

const (
	nodesBitsChildren   = %d
	nodesBitsICANN      = %d
	nodesBitsTextOffset = %d
	nodesBitsTextLength = %d

	childrenBitsWildcard = %d
	childrenBitsNodeType = %d
	childrenBitsHi       = %d
	childrenBitsLo       = %d
)

const (
	nodeTypeNormal     = %d
	nodeTypeException  = %d
	nodeTypeParentOnly = %d
)

// numTLD is the number of top level domains.
const numTLD = %d

`
	fmt.Fprintf(w, header, *version,
		nodesBitsChildren, nodesBitsICANN, nodesBitsTextOffset, nodesBitsTextLength,
		childrenBitsWildcard, childrenBitsNodeType, childrenBitsHi, childrenBitsLo,
		nodeTypeNormal, nodeTypeException, nodeTypeParentOnly, len(n.children))

	text := combineText(labelsList)
	if text == "" ***REMOVED***
		return fmt.Errorf("internal error: makeText returned no text")
	***REMOVED***
	for _, label := range labelsList ***REMOVED***
		offset, length := strings.Index(text, label), len(label)
		if offset < 0 ***REMOVED***
			return fmt.Errorf("internal error: could not find %q in text %q", label, text)
		***REMOVED***
		maxTextOffset, maxTextLength = max(maxTextOffset, offset), max(maxTextLength, length)
		if offset >= 1<<nodesBitsTextOffset ***REMOVED***
			return fmt.Errorf("text offset %d is too large, or nodeBitsTextOffset is too small", offset)
		***REMOVED***
		if length >= 1<<nodesBitsTextLength ***REMOVED***
			return fmt.Errorf("text length %d is too large, or nodeBitsTextLength is too small", length)
		***REMOVED***
		labelEncoding[label] = uint32(offset)<<nodesBitsTextLength | uint32(length)
	***REMOVED***
	fmt.Fprintf(w, "// Text is the combined text of all labels.\nconst text = ")
	for len(text) > 0 ***REMOVED***
		n, plus := len(text), ""
		if n > 64 ***REMOVED***
			n, plus = 64, " +"
		***REMOVED***
		fmt.Fprintf(w, "%q%s\n", text[:n], plus)
		text = text[n:]
	***REMOVED***

	if err := n.walk(w, assignIndexes); err != nil ***REMOVED***
		return err
	***REMOVED***

	fmt.Fprintf(w, `

// nodes is the list of nodes. Each node is represented as a uint32, which
// encodes the node's children, wildcard bit and node type (as an index into
// the children array), ICANN bit and text.
//
// If the table was generated with the -comments flag, there is a //-comment
// after each node's data. In it is the nodes-array indexes of the children,
// formatted as (n0x1234-n0x1256), with * denoting the wildcard bit. The
// nodeType is printed as + for normal, ! for exception, and o for parent-only
// nodes that have children but don't match a domain label in their own right.
// An I denotes an ICANN domain.
//
// The layout within the uint32, from MSB to LSB, is:
//	[%2d bits] unused
//	[%2d bits] children index
//	[%2d bits] ICANN bit
//	[%2d bits] text index
//	[%2d bits] text length
var nodes = [...]uint32***REMOVED***
`,
		32-nodesBitsChildren-nodesBitsICANN-nodesBitsTextOffset-nodesBitsTextLength,
		nodesBitsChildren, nodesBitsICANN, nodesBitsTextOffset, nodesBitsTextLength)
	if err := n.walk(w, printNode); err != nil ***REMOVED***
		return err
	***REMOVED***
	fmt.Fprintf(w, `***REMOVED***

// children is the list of nodes' children, the parent's wildcard bit and the
// parent's node type. If a node has no children then their children index
// will be in the range [0, 6), depending on the wildcard bit and node type.
//
// The layout within the uint32, from MSB to LSB, is:
//	[%2d bits] unused
//	[%2d bits] wildcard bit
//	[%2d bits] node type
//	[%2d bits] high nodes index (exclusive) of children
//	[%2d bits] low nodes index (inclusive) of children
var children=[...]uint32***REMOVED***
`,
		32-childrenBitsWildcard-childrenBitsNodeType-childrenBitsHi-childrenBitsLo,
		childrenBitsWildcard, childrenBitsNodeType, childrenBitsHi, childrenBitsLo)
	for i, c := range childrenEncoding ***REMOVED***
		s := "---------------"
		lo := c & (1<<childrenBitsLo - 1)
		hi := (c >> childrenBitsLo) & (1<<childrenBitsHi - 1)
		if lo != hi ***REMOVED***
			s = fmt.Sprintf("n0x%04x-n0x%04x", lo, hi)
		***REMOVED***
		nodeType := int(c>>(childrenBitsLo+childrenBitsHi)) & (1<<childrenBitsNodeType - 1)
		wildcard := c>>(childrenBitsLo+childrenBitsHi+childrenBitsNodeType) != 0
		if *comments ***REMOVED***
			fmt.Fprintf(w, "0x%08x, // c0x%04x (%s)%s %s\n",
				c, i, s, wildcardStr(wildcard), nodeTypeStr(nodeType))
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(w, "0x%x,\n", c)
		***REMOVED***
	***REMOVED***
	fmt.Fprintf(w, "***REMOVED***\n\n")
	fmt.Fprintf(w, "// max children %d (capacity %d)\n", maxChildren, 1<<nodesBitsChildren-1)
	fmt.Fprintf(w, "// max text offset %d (capacity %d)\n", maxTextOffset, 1<<nodesBitsTextOffset-1)
	fmt.Fprintf(w, "// max text length %d (capacity %d)\n", maxTextLength, 1<<nodesBitsTextLength-1)
	fmt.Fprintf(w, "// max hi %d (capacity %d)\n", maxHi, 1<<childrenBitsHi-1)
	fmt.Fprintf(w, "// max lo %d (capacity %d)\n", maxLo, 1<<childrenBitsLo-1)
	return nil
***REMOVED***

type node struct ***REMOVED***
	label    string
	nodeType int
	icann    bool
	wildcard bool
	// nodesIndex and childrenIndex are the index of this node in the nodes
	// and the index of its children offset/length in the children arrays.
	nodesIndex, childrenIndex int
	// firstChild is the index of this node's first child, or zero if this
	// node has no children.
	firstChild int
	// children are the node's children, in strictly increasing node label order.
	children []*node
***REMOVED***

func (n *node) walk(w io.Writer, f func(w1 io.Writer, n1 *node) error) error ***REMOVED***
	if err := f(w, n); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, c := range n.children ***REMOVED***
		if err := c.walk(w, f); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// child returns the child of n with the given label. The child is created if
// it did not exist beforehand.
func (n *node) child(label string) *node ***REMOVED***
	for _, c := range n.children ***REMOVED***
		if c.label == label ***REMOVED***
			return c
		***REMOVED***
	***REMOVED***
	c := &node***REMOVED***
		label:    label,
		nodeType: nodeTypeParentOnly,
		icann:    true,
	***REMOVED***
	n.children = append(n.children, c)
	sort.Sort(byLabel(n.children))
	return c
***REMOVED***

type byLabel []*node

func (b byLabel) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b byLabel) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b byLabel) Less(i, j int) bool ***REMOVED*** return b[i].label < b[j].label ***REMOVED***

var nextNodesIndex int

// childrenEncoding are the encoded entries in the generated children array.
// All these pre-defined entries have no children.
var childrenEncoding = []uint32***REMOVED***
	0 << (childrenBitsLo + childrenBitsHi), // Without wildcard bit, nodeTypeNormal.
	1 << (childrenBitsLo + childrenBitsHi), // Without wildcard bit, nodeTypeException.
	2 << (childrenBitsLo + childrenBitsHi), // Without wildcard bit, nodeTypeParentOnly.
	4 << (childrenBitsLo + childrenBitsHi), // With wildcard bit, nodeTypeNormal.
	5 << (childrenBitsLo + childrenBitsHi), // With wildcard bit, nodeTypeException.
	6 << (childrenBitsLo + childrenBitsHi), // With wildcard bit, nodeTypeParentOnly.
***REMOVED***

var firstCallToAssignIndexes = true

func assignIndexes(w io.Writer, n *node) error ***REMOVED***
	if len(n.children) != 0 ***REMOVED***
		// Assign nodesIndex.
		n.firstChild = nextNodesIndex
		for _, c := range n.children ***REMOVED***
			c.nodesIndex = nextNodesIndex
			nextNodesIndex++
		***REMOVED***

		// The root node's children is implicit.
		if firstCallToAssignIndexes ***REMOVED***
			firstCallToAssignIndexes = false
			return nil
		***REMOVED***

		// Assign childrenIndex.
		maxChildren = max(maxChildren, len(childrenEncoding))
		if len(childrenEncoding) >= 1<<nodesBitsChildren ***REMOVED***
			return fmt.Errorf("children table size %d is too large, or nodeBitsChildren is too small", len(childrenEncoding))
		***REMOVED***
		n.childrenIndex = len(childrenEncoding)
		lo := uint32(n.firstChild)
		hi := lo + uint32(len(n.children))
		maxLo, maxHi = u32max(maxLo, lo), u32max(maxHi, hi)
		if lo >= 1<<childrenBitsLo ***REMOVED***
			return fmt.Errorf("children lo %d is too large, or childrenBitsLo is too small", lo)
		***REMOVED***
		if hi >= 1<<childrenBitsHi ***REMOVED***
			return fmt.Errorf("children hi %d is too large, or childrenBitsHi is too small", hi)
		***REMOVED***
		enc := hi<<childrenBitsLo | lo
		enc |= uint32(n.nodeType) << (childrenBitsLo + childrenBitsHi)
		if n.wildcard ***REMOVED***
			enc |= 1 << (childrenBitsLo + childrenBitsHi + childrenBitsNodeType)
		***REMOVED***
		childrenEncoding = append(childrenEncoding, enc)
	***REMOVED*** else ***REMOVED***
		n.childrenIndex = n.nodeType
		if n.wildcard ***REMOVED***
			n.childrenIndex += numNodeType
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func printNode(w io.Writer, n *node) error ***REMOVED***
	for _, c := range n.children ***REMOVED***
		s := "---------------"
		if len(c.children) != 0 ***REMOVED***
			s = fmt.Sprintf("n0x%04x-n0x%04x", c.firstChild, c.firstChild+len(c.children))
		***REMOVED***
		encoding := labelEncoding[c.label]
		if c.icann ***REMOVED***
			encoding |= 1 << (nodesBitsTextLength + nodesBitsTextOffset)
		***REMOVED***
		encoding |= uint32(c.childrenIndex) << (nodesBitsTextLength + nodesBitsTextOffset + nodesBitsICANN)
		if *comments ***REMOVED***
			fmt.Fprintf(w, "0x%08x, // n0x%04x c0x%04x (%s)%s %s %s %s\n",
				encoding, c.nodesIndex, c.childrenIndex, s, wildcardStr(c.wildcard),
				nodeTypeStr(c.nodeType), icannStr(c.icann), c.label,
			)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(w, "0x%x,\n", encoding)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func printNodeLabel(w io.Writer, n *node) error ***REMOVED***
	for _, c := range n.children ***REMOVED***
		fmt.Fprintf(w, "%q,\n", c.label)
	***REMOVED***
	return nil
***REMOVED***

func icannStr(icann bool) string ***REMOVED***
	if icann ***REMOVED***
		return "I"
	***REMOVED***
	return " "
***REMOVED***

func wildcardStr(wildcard bool) string ***REMOVED***
	if wildcard ***REMOVED***
		return "*"
	***REMOVED***
	return " "
***REMOVED***

// combineText combines all the strings in labelsList to form one giant string.
// Overlapping strings will be merged: "arpa" and "parliament" could yield
// "arparliament".
func combineText(labelsList []string) string ***REMOVED***
	beforeLength := 0
	for _, s := range labelsList ***REMOVED***
		beforeLength += len(s)
	***REMOVED***

	text := crush(removeSubstrings(labelsList))
	if *v ***REMOVED***
		fmt.Fprintf(os.Stderr, "crushed %d bytes to become %d bytes\n", beforeLength, len(text))
	***REMOVED***
	return text
***REMOVED***

type byLength []string

func (s byLength) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s byLength) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s byLength) Less(i, j int) bool ***REMOVED*** return len(s[i]) < len(s[j]) ***REMOVED***

// removeSubstrings returns a copy of its input with any strings removed
// that are substrings of other provided strings.
func removeSubstrings(input []string) []string ***REMOVED***
	// Make a copy of input.
	ss := append(make([]string, 0, len(input)), input...)
	sort.Sort(byLength(ss))

	for i, shortString := range ss ***REMOVED***
		// For each string, only consider strings higher than it in sort order, i.e.
		// of equal length or greater.
		for _, longString := range ss[i+1:] ***REMOVED***
			if strings.Contains(longString, shortString) ***REMOVED***
				ss[i] = ""
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Remove the empty strings.
	sort.Strings(ss)
	for len(ss) > 0 && ss[0] == "" ***REMOVED***
		ss = ss[1:]
	***REMOVED***
	return ss
***REMOVED***

// crush combines a list of strings, taking advantage of overlaps. It returns a
// single string that contains each input string as a substring.
func crush(ss []string) string ***REMOVED***
	maxLabelLen := 0
	for _, s := range ss ***REMOVED***
		if maxLabelLen < len(s) ***REMOVED***
			maxLabelLen = len(s)
		***REMOVED***
	***REMOVED***

	for prefixLen := maxLabelLen; prefixLen > 0; prefixLen-- ***REMOVED***
		prefixes := makePrefixMap(ss, prefixLen)
		for i, s := range ss ***REMOVED***
			if len(s) <= prefixLen ***REMOVED***
				continue
			***REMOVED***
			mergeLabel(ss, i, prefixLen, prefixes)
		***REMOVED***
	***REMOVED***

	return strings.Join(ss, "")
***REMOVED***

// mergeLabel merges the label at ss[i] with the first available matching label
// in prefixMap, where the last "prefixLen" characters in ss[i] match the first
// "prefixLen" characters in the matching label.
// It will merge ss[i] repeatedly until no more matches are available.
// All matching labels merged into ss[i] are replaced by "".
func mergeLabel(ss []string, i, prefixLen int, prefixes prefixMap) ***REMOVED***
	s := ss[i]
	suffix := s[len(s)-prefixLen:]
	for _, j := range prefixes[suffix] ***REMOVED***
		// Empty strings mean "already used." Also avoid merging with self.
		if ss[j] == "" || i == j ***REMOVED***
			continue
		***REMOVED***
		if *v ***REMOVED***
			fmt.Fprintf(os.Stderr, "%d-length overlap at (%4d,%4d): %q and %q share %q\n",
				prefixLen, i, j, ss[i], ss[j], suffix)
		***REMOVED***
		ss[i] += ss[j][prefixLen:]
		ss[j] = ""
		// ss[i] has a new suffix, so merge again if possible.
		// Note: we only have to merge again at the same prefix length. Shorter
		// prefix lengths will be handled in the next iteration of crush's for loop.
		// Can there be matches for longer prefix lengths, introduced by the merge?
		// I believe that any such matches would by necessity have been eliminated
		// during substring removal or merged at a higher prefix length. For
		// instance, in crush("abc", "cde", "bcdef"), combining "abc" and "cde"
		// would yield "abcde", which could be merged with "bcdef." However, in
		// practice "cde" would already have been elimintated by removeSubstrings.
		mergeLabel(ss, i, prefixLen, prefixes)
		return
	***REMOVED***
***REMOVED***

// prefixMap maps from a prefix to a list of strings containing that prefix. The
// list of strings is represented as indexes into a slice of strings stored
// elsewhere.
type prefixMap map[string][]int

// makePrefixMap constructs a prefixMap from a slice of strings.
func makePrefixMap(ss []string, prefixLen int) prefixMap ***REMOVED***
	prefixes := make(prefixMap)
	for i, s := range ss ***REMOVED***
		// We use < rather than <= because if a label matches on a prefix equal to
		// its full length, that's actually a substring match handled by
		// removeSubstrings.
		if prefixLen < len(s) ***REMOVED***
			prefix := s[:prefixLen]
			prefixes[prefix] = append(prefixes[prefix], i)
		***REMOVED***
	***REMOVED***

	return prefixes
***REMOVED***
