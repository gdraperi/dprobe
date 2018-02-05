// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build // import "golang.org/x/text/collate/build"

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// TODO: optimizations:
// - expandElem is currently 20K. By putting unique colElems in a separate
//   table and having a byte array of indexes into this table, we can reduce
//   the total size to about 7K. By also factoring out the length bytes, we
//   can reduce this to about 6K.
// - trie valueBlocks are currently 100K. There are a lot of sparse blocks
//   and many consecutive values with the same stride. This can be further
//   compacted.
// - Compress secondary weights into 8 bits.
// - Some LDML specs specify a context element. Currently we simply concatenate
//   those.  Context can be implemented using the contraction trie. If Builder
//   could analyze and detect when using a context makes sense, there is no
//   need to expose this construct in the API.

// A Builder builds a root collation table.  The user must specify the
// collation elements for each entry.  A common use will be to base the weights
// on those specified in the allkeys* file as provided by the UCA or CLDR.
type Builder struct ***REMOVED***
	index  *trieBuilder
	root   ordering
	locale []*Tailoring
	t      *table
	err    error
	built  bool

	minNonVar int // lowest primary recorded for a variable
	varTop    int // highest primary recorded for a non-variable

	// indexes used for reusing expansions and contractions
	expIndex map[string]int      // positions of expansions keyed by their string representation
	ctHandle map[string]ctHandle // contraction handles keyed by a concatenation of the suffixes
	ctElem   map[string]int      // contraction elements keyed by their string representation
***REMOVED***

// A Tailoring builds a collation table based on another collation table.
// The table is defined by specifying tailorings to the underlying table.
// See http://unicode.org/reports/tr35/ for an overview of tailoring
// collation tables.  The CLDR contains pre-defined tailorings for a variety
// of languages (See http://www.unicode.org/Public/cldr/<version>/core.zip.)
type Tailoring struct ***REMOVED***
	id      string
	builder *Builder
	index   *ordering

	anchor *entry
	before bool
***REMOVED***

// NewBuilder returns a new Builder.
func NewBuilder() *Builder ***REMOVED***
	return &Builder***REMOVED***
		index:    newTrieBuilder(),
		root:     makeRootOrdering(),
		expIndex: make(map[string]int),
		ctHandle: make(map[string]ctHandle),
		ctElem:   make(map[string]int),
	***REMOVED***
***REMOVED***

// Tailoring returns a Tailoring for the given locale.  One should
// have completed all calls to Add before calling Tailoring.
func (b *Builder) Tailoring(loc language.Tag) *Tailoring ***REMOVED***
	t := &Tailoring***REMOVED***
		id:      loc.String(),
		builder: b,
		index:   b.root.clone(),
	***REMOVED***
	t.index.id = t.id
	b.locale = append(b.locale, t)
	return t
***REMOVED***

// Add adds an entry to the collation element table, mapping
// a slice of runes to a sequence of collation elements.
// A collation element is specified as list of weights: []int***REMOVED***primary, secondary, ...***REMOVED***.
// The entries are typically obtained from a collation element table
// as defined in http://www.unicode.org/reports/tr10/#Data_Table_Format.
// Note that the collation elements specified by colelems are only used
// as a guide.  The actual weights generated by Builder may differ.
// The argument variables is a list of indices into colelems that should contain
// a value for each colelem that is a variable. (See the reference above.)
func (b *Builder) Add(runes []rune, colelems [][]int, variables []int) error ***REMOVED***
	str := string(runes)
	elems := make([]rawCE, len(colelems))
	for i, ce := range colelems ***REMOVED***
		if len(ce) == 0 ***REMOVED***
			break
		***REMOVED***
		elems[i] = makeRawCE(ce, 0)
		if len(ce) == 1 ***REMOVED***
			elems[i].w[1] = defaultSecondary
		***REMOVED***
		if len(ce) <= 2 ***REMOVED***
			elems[i].w[2] = defaultTertiary
		***REMOVED***
		if len(ce) <= 3 ***REMOVED***
			elems[i].w[3] = ce[0]
		***REMOVED***
	***REMOVED***
	for i, ce := range elems ***REMOVED***
		p := ce.w[0]
		isvar := false
		for _, j := range variables ***REMOVED***
			if i == j ***REMOVED***
				isvar = true
			***REMOVED***
		***REMOVED***
		if isvar ***REMOVED***
			if p >= b.minNonVar && b.minNonVar > 0 ***REMOVED***
				return fmt.Errorf("primary value %X of variable is larger than the smallest non-variable %X", p, b.minNonVar)
			***REMOVED***
			if p > b.varTop ***REMOVED***
				b.varTop = p
			***REMOVED***
		***REMOVED*** else if p > 1 ***REMOVED*** // 1 is a special primary value reserved for FFFE
			if p <= b.varTop ***REMOVED***
				return fmt.Errorf("primary value %X of non-variable is smaller than the highest variable %X", p, b.varTop)
			***REMOVED***
			if b.minNonVar == 0 || p < b.minNonVar ***REMOVED***
				b.minNonVar = p
			***REMOVED***
		***REMOVED***
	***REMOVED***
	elems, err := convertLargeWeights(elems)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cccs := []uint8***REMOVED******REMOVED***
	nfd := norm.NFD.String(str)
	for i := range nfd ***REMOVED***
		cccs = append(cccs, norm.NFD.PropertiesString(nfd[i:]).CCC())
	***REMOVED***
	if len(cccs) < len(elems) ***REMOVED***
		if len(cccs) > 2 ***REMOVED***
			return fmt.Errorf("number of decomposed characters should be greater or equal to the number of collation elements for len(colelems) > 3 (%d < %d)", len(cccs), len(elems))
		***REMOVED***
		p := len(elems) - 1
		for ; p > 0 && elems[p].w[0] == 0; p-- ***REMOVED***
			elems[p].ccc = cccs[len(cccs)-1]
		***REMOVED***
		for ; p >= 0; p-- ***REMOVED***
			elems[p].ccc = cccs[0]
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := range elems ***REMOVED***
			elems[i].ccc = cccs[i]
		***REMOVED***
	***REMOVED***
	// doNorm in collate.go assumes that the following conditions hold.
	if len(elems) > 1 && len(cccs) > 1 && cccs[0] != 0 && cccs[0] != cccs[len(cccs)-1] ***REMOVED***
		return fmt.Errorf("incompatible CCC values for expansion %X (%d)", runes, cccs)
	***REMOVED***
	b.root.newEntry(str, elems)
	return nil
***REMOVED***

func (t *Tailoring) setAnchor(anchor string) error ***REMOVED***
	anchor = norm.NFC.String(anchor)
	a := t.index.find(anchor)
	if a == nil ***REMOVED***
		a = t.index.newEntry(anchor, nil)
		a.implicit = true
		a.modified = true
		for _, r := range []rune(anchor) ***REMOVED***
			e := t.index.find(string(r))
			e.lock = true
		***REMOVED***
	***REMOVED***
	t.anchor = a
	return nil
***REMOVED***

// SetAnchor sets the point after which elements passed in subsequent calls to
// Insert will be inserted.  It is equivalent to the reset directive in an LDML
// specification.  See Insert for an example.
// SetAnchor supports the following logical reset positions:
// <first_tertiary_ignorable/>, <last_teriary_ignorable/>, <first_primary_ignorable/>,
// and <last_non_ignorable/>.
func (t *Tailoring) SetAnchor(anchor string) error ***REMOVED***
	if err := t.setAnchor(anchor); err != nil ***REMOVED***
		return err
	***REMOVED***
	t.before = false
	return nil
***REMOVED***

// SetAnchorBefore is similar to SetAnchor, except that subsequent calls to
// Insert will insert entries before the anchor.
func (t *Tailoring) SetAnchorBefore(anchor string) error ***REMOVED***
	if err := t.setAnchor(anchor); err != nil ***REMOVED***
		return err
	***REMOVED***
	t.before = true
	return nil
***REMOVED***

// Insert sets the ordering of str relative to the entry set by the previous
// call to SetAnchor or Insert.  The argument extend corresponds
// to the extend elements as defined in LDML.  A non-empty value for extend
// will cause the collation elements corresponding to extend to be appended
// to the collation elements generated for the entry added by Insert.
// This has the same net effect as sorting str after the string anchor+extend.
// See http://www.unicode.org/reports/tr10/#Tailoring_Example for details
// on parametric tailoring and http://unicode.org/reports/tr35/#Collation_Elements
// for full details on LDML.
//
// Examples: create a tailoring for Swedish, where "ä" is ordered after "z"
// at the primary sorting level:
//      t := b.Tailoring("se")
// 		t.SetAnchor("z")
// 		t.Insert(colltab.Primary, "ä", "")
// Order "ü" after "ue" at the secondary sorting level:
//		t.SetAnchor("ue")
//		t.Insert(colltab.Secondary, "ü","")
// or
//		t.SetAnchor("u")
//		t.Insert(colltab.Secondary, "ü", "e")
// Order "q" afer "ab" at the secondary level and "Q" after "q"
// at the tertiary level:
// 		t.SetAnchor("ab")
// 		t.Insert(colltab.Secondary, "q", "")
// 		t.Insert(colltab.Tertiary, "Q", "")
// Order "b" before "a":
//      t.SetAnchorBefore("a")
//      t.Insert(colltab.Primary, "b", "")
// Order "0" after the last primary ignorable:
//      t.SetAnchor("<last_primary_ignorable/>")
//      t.Insert(colltab.Primary, "0", "")
func (t *Tailoring) Insert(level colltab.Level, str, extend string) error ***REMOVED***
	if t.anchor == nil ***REMOVED***
		return fmt.Errorf("%s:Insert: no anchor point set for tailoring of %s", t.id, str)
	***REMOVED***
	str = norm.NFC.String(str)
	e := t.index.find(str)
	if e == nil ***REMOVED***
		e = t.index.newEntry(str, nil)
	***REMOVED*** else if e.logical != noAnchor ***REMOVED***
		return fmt.Errorf("%s:Insert: cannot reinsert logical reset position %q", t.id, e.str)
	***REMOVED***
	if e.lock ***REMOVED***
		return fmt.Errorf("%s:Insert: cannot reinsert element %q", t.id, e.str)
	***REMOVED***
	a := t.anchor
	// Find the first element after the anchor which differs at a level smaller or
	// equal to the given level.  Then insert at this position.
	// See http://unicode.org/reports/tr35/#Collation_Elements, Section 5.14.5 for details.
	e.before = t.before
	if t.before ***REMOVED***
		t.before = false
		if a.prev == nil ***REMOVED***
			a.insertBefore(e)
		***REMOVED*** else ***REMOVED***
			for a = a.prev; a.level > level; a = a.prev ***REMOVED***
			***REMOVED***
			a.insertAfter(e)
		***REMOVED***
		e.level = level
	***REMOVED*** else ***REMOVED***
		for ; a.level > level; a = a.next ***REMOVED***
		***REMOVED***
		e.level = a.level
		if a != e ***REMOVED***
			a.insertAfter(e)
			a.level = level
		***REMOVED*** else ***REMOVED***
			// We don't set a to prev itself. This has the effect of the entry
			// getting new collation elements that are an increment of itself.
			// This is intentional.
			a.prev.level = level
		***REMOVED***
	***REMOVED***
	e.extend = norm.NFD.String(extend)
	e.exclude = false
	e.modified = true
	e.elems = nil
	t.anchor = e
	return nil
***REMOVED***

func (o *ordering) getWeight(e *entry) []rawCE ***REMOVED***
	if len(e.elems) == 0 && e.logical == noAnchor ***REMOVED***
		if e.implicit ***REMOVED***
			for _, r := range e.runes ***REMOVED***
				e.elems = append(e.elems, o.getWeight(o.find(string(r)))...)
			***REMOVED***
		***REMOVED*** else if e.before ***REMOVED***
			count := [colltab.Identity + 1]int***REMOVED******REMOVED***
			a := e
			for ; a.elems == nil && !a.implicit; a = a.next ***REMOVED***
				count[a.level]++
			***REMOVED***
			e.elems = []rawCE***REMOVED***makeRawCE(a.elems[0].w, a.elems[0].ccc)***REMOVED***
			for i := colltab.Primary; i < colltab.Quaternary; i++ ***REMOVED***
				if count[i] != 0 ***REMOVED***
					e.elems[0].w[i] -= count[i]
					break
				***REMOVED***
			***REMOVED***
			if e.prev != nil ***REMOVED***
				o.verifyWeights(e.prev, e, e.prev.level)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			prev := e.prev
			e.elems = nextWeight(prev.level, o.getWeight(prev))
			o.verifyWeights(e, e.next, e.level)
		***REMOVED***
	***REMOVED***
	return e.elems
***REMOVED***

func (o *ordering) addExtension(e *entry) ***REMOVED***
	if ex := o.find(e.extend); ex != nil ***REMOVED***
		e.elems = append(e.elems, ex.elems...)
	***REMOVED*** else ***REMOVED***
		for _, r := range []rune(e.extend) ***REMOVED***
			e.elems = append(e.elems, o.find(string(r)).elems...)
		***REMOVED***
	***REMOVED***
	e.extend = ""
***REMOVED***

func (o *ordering) verifyWeights(a, b *entry, level colltab.Level) error ***REMOVED***
	if level == colltab.Identity || b == nil || b.elems == nil || a.elems == nil ***REMOVED***
		return nil
	***REMOVED***
	for i := colltab.Primary; i < level; i++ ***REMOVED***
		if a.elems[0].w[i] < b.elems[0].w[i] ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	if a.elems[0].w[level] >= b.elems[0].w[level] ***REMOVED***
		err := fmt.Errorf("%s:overflow: collation elements of %q (%X) overflows those of %q (%X) at level %d (%X >= %X)", o.id, a.str, a.runes, b.str, b.runes, level, a.elems, b.elems)
		log.Println(err)
		// TODO: return the error instead, or better, fix the conflicting entry by making room.
	***REMOVED***
	return nil
***REMOVED***

func (b *Builder) error(e error) ***REMOVED***
	if e != nil ***REMOVED***
		b.err = e
	***REMOVED***
***REMOVED***

func (b *Builder) errorID(locale string, e error) ***REMOVED***
	if e != nil ***REMOVED***
		b.err = fmt.Errorf("%s:%v", locale, e)
	***REMOVED***
***REMOVED***

// patchNorm ensures that NFC and NFD counterparts are consistent.
func (o *ordering) patchNorm() ***REMOVED***
	// Insert the NFD counterparts, if necessary.
	for _, e := range o.ordered ***REMOVED***
		nfd := norm.NFD.String(e.str)
		if nfd != e.str ***REMOVED***
			if e0 := o.find(nfd); e0 != nil && !e0.modified ***REMOVED***
				e0.elems = e.elems
			***REMOVED*** else if e.modified && !equalCEArrays(o.genColElems(nfd), e.elems) ***REMOVED***
				e := o.newEntry(nfd, e.elems)
				e.modified = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Update unchanged composed forms if one of their parts changed.
	for _, e := range o.ordered ***REMOVED***
		nfd := norm.NFD.String(e.str)
		if e.modified || nfd == e.str ***REMOVED***
			continue
		***REMOVED***
		if e0 := o.find(nfd); e0 != nil ***REMOVED***
			e.elems = e0.elems
		***REMOVED*** else ***REMOVED***
			e.elems = o.genColElems(nfd)
			if norm.NFD.LastBoundary([]byte(nfd)) == 0 ***REMOVED***
				r := []rune(nfd)
				head := string(r[0])
				tail := ""
				for i := 1; i < len(r); i++ ***REMOVED***
					s := norm.NFC.String(head + string(r[i]))
					if e0 := o.find(s); e0 != nil && e0.modified ***REMOVED***
						head = s
					***REMOVED*** else ***REMOVED***
						tail += string(r[i])
					***REMOVED***
				***REMOVED***
				e.elems = append(o.genColElems(head), o.genColElems(tail)...)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Exclude entries for which the individual runes generate the same collation elements.
	for _, e := range o.ordered ***REMOVED***
		if len(e.runes) > 1 && equalCEArrays(o.genColElems(e.str), e.elems) ***REMOVED***
			e.exclude = true
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *Builder) buildOrdering(o *ordering) ***REMOVED***
	for _, e := range o.ordered ***REMOVED***
		o.getWeight(e)
	***REMOVED***
	for _, e := range o.ordered ***REMOVED***
		o.addExtension(e)
	***REMOVED***
	o.patchNorm()
	o.sort()
	simplify(o)
	b.processExpansions(o)   // requires simplify
	b.processContractions(o) // requires simplify

	t := newNode()
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		if !e.skip() ***REMOVED***
			ce, err := e.encode()
			b.errorID(o.id, err)
			t.insert(e.runes[0], ce)
		***REMOVED***
	***REMOVED***
	o.handle = b.index.addTrie(t)
***REMOVED***

func (b *Builder) build() (*table, error) ***REMOVED***
	if b.built ***REMOVED***
		return b.t, b.err
	***REMOVED***
	b.built = true
	b.t = &table***REMOVED***
		Table: colltab.Table***REMOVED***
			MaxContractLen: utf8.UTFMax,
			VariableTop:    uint32(b.varTop),
		***REMOVED***,
	***REMOVED***

	b.buildOrdering(&b.root)
	b.t.root = b.root.handle
	for _, t := range b.locale ***REMOVED***
		b.buildOrdering(t.index)
		if b.err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	i, err := b.index.generate()
	b.t.trie = *i
	b.t.Index = colltab.Trie***REMOVED***
		Index:   i.index,
		Values:  i.values,
		Index0:  i.index[blockSize*b.t.root.lookupStart:],
		Values0: i.values[blockSize*b.t.root.valueStart:],
	***REMOVED***
	b.error(err)
	return b.t, b.err
***REMOVED***

// Build builds the root Collator.
func (b *Builder) Build() (colltab.Weighter, error) ***REMOVED***
	table, err := b.build()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return table, nil
***REMOVED***

// Build builds a Collator for Tailoring t.
func (t *Tailoring) Build() (colltab.Weighter, error) ***REMOVED***
	// TODO: implement.
	return nil, nil
***REMOVED***

// Print prints the tables for b and all its Tailorings as a Go file
// that can be included in the Collate package.
func (b *Builder) Print(w io.Writer) (n int, err error) ***REMOVED***
	p := func(nn int, e error) ***REMOVED***
		n += nn
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	t, err := b.build()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	p(fmt.Fprintf(w, `var availableLocales = "und`))
	for _, loc := range b.locale ***REMOVED***
		if loc.id != "und" ***REMOVED***
			p(fmt.Fprintf(w, ",%s", loc.id))
		***REMOVED***
	***REMOVED***
	p(fmt.Fprint(w, "\"\n\n"))
	p(fmt.Fprintf(w, "const varTop = 0x%x\n\n", b.varTop))
	p(fmt.Fprintln(w, "var locales = [...]tableIndex***REMOVED***"))
	for _, loc := range b.locale ***REMOVED***
		if loc.id == "und" ***REMOVED***
			p(t.fprintIndex(w, loc.index.handle, loc.id))
		***REMOVED***
	***REMOVED***
	for _, loc := range b.locale ***REMOVED***
		if loc.id != "und" ***REMOVED***
			p(t.fprintIndex(w, loc.index.handle, loc.id))
		***REMOVED***
	***REMOVED***
	p(fmt.Fprint(w, "***REMOVED***\n\n"))
	n, _, err = t.fprint(w, "main")
	return
***REMOVED***

// reproducibleFromNFKD checks whether the given expansion could be generated
// from an NFKD expansion.
func reproducibleFromNFKD(e *entry, exp, nfkd []rawCE) bool ***REMOVED***
	// Length must be equal.
	if len(exp) != len(nfkd) ***REMOVED***
		return false
	***REMOVED***
	for i, ce := range exp ***REMOVED***
		// Primary and secondary values should be equal.
		if ce.w[0] != nfkd[i].w[0] || ce.w[1] != nfkd[i].w[1] ***REMOVED***
			return false
		***REMOVED***
		// Tertiary values should be equal to maxTertiary for third element onwards.
		// TODO: there seem to be a lot of cases in CLDR (e.g. ㏭ in zh.xml) that can
		// simply be dropped.  Try this out by dropping the following code.
		if i >= 2 && ce.w[2] != maxTertiary ***REMOVED***
			return false
		***REMOVED***
		if _, err := makeCE(ce); err != nil ***REMOVED***
			// Simply return false. The error will be caught elsewhere.
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func simplify(o *ordering) ***REMOVED***
	// Runes that are a starter of a contraction should not be removed.
	// (To date, there is only Kannada character 0CCA.)
	keep := make(map[rune]bool)
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		if len(e.runes) > 1 ***REMOVED***
			keep[e.runes[0]] = true
		***REMOVED***
	***REMOVED***
	// Tag entries for which the runes NFKD decompose to identical values.
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		s := e.str
		nfkd := norm.NFKD.String(s)
		nfd := norm.NFD.String(s)
		if e.decompose || len(e.runes) > 1 || len(e.elems) == 1 || keep[e.runes[0]] || nfkd == nfd ***REMOVED***
			continue
		***REMOVED***
		if reproducibleFromNFKD(e, e.elems, o.genColElems(nfkd)) ***REMOVED***
			e.decompose = true
		***REMOVED***
	***REMOVED***
***REMOVED***

// appendExpansion converts the given collation sequence to
// collation elements and adds them to the expansion table.
// It returns an index to the expansion table.
func (b *Builder) appendExpansion(e *entry) int ***REMOVED***
	t := b.t
	i := len(t.ExpandElem)
	ce := uint32(len(e.elems))
	t.ExpandElem = append(t.ExpandElem, ce)
	for _, w := range e.elems ***REMOVED***
		ce, err := makeCE(w)
		if err != nil ***REMOVED***
			b.error(err)
			return -1
		***REMOVED***
		t.ExpandElem = append(t.ExpandElem, ce)
	***REMOVED***
	return i
***REMOVED***

// processExpansions extracts data necessary to generate
// the extraction tables.
func (b *Builder) processExpansions(o *ordering) ***REMOVED***
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		if !e.expansion() ***REMOVED***
			continue
		***REMOVED***
		key := fmt.Sprintf("%v", e.elems)
		i, ok := b.expIndex[key]
		if !ok ***REMOVED***
			i = b.appendExpansion(e)
			b.expIndex[key] = i
		***REMOVED***
		e.expansionIndex = i
	***REMOVED***
***REMOVED***

func (b *Builder) processContractions(o *ordering) ***REMOVED***
	// Collate contractions per starter rune.
	starters := []rune***REMOVED******REMOVED***
	cm := make(map[rune][]*entry)
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		if e.contraction() ***REMOVED***
			if len(e.str) > b.t.MaxContractLen ***REMOVED***
				b.t.MaxContractLen = len(e.str)
			***REMOVED***
			r := e.runes[0]
			if _, ok := cm[r]; !ok ***REMOVED***
				starters = append(starters, r)
			***REMOVED***
			cm[r] = append(cm[r], e)
		***REMOVED***
	***REMOVED***
	// Add entries of single runes that are at a start of a contraction.
	for e := o.front(); e != nil; e, _ = e.nextIndexed() ***REMOVED***
		if !e.contraction() ***REMOVED***
			r := e.runes[0]
			if _, ok := cm[r]; ok ***REMOVED***
				cm[r] = append(cm[r], e)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Build the tries for the contractions.
	t := b.t
	for _, r := range starters ***REMOVED***
		l := cm[r]
		// Compute suffix strings. There are 31 different contraction suffix
		// sets for 715 contractions and 82 contraction starter runes as of
		// version 6.0.0.
		sufx := []string***REMOVED******REMOVED***
		hasSingle := false
		for _, e := range l ***REMOVED***
			if len(e.runes) > 1 ***REMOVED***
				sufx = append(sufx, string(e.runes[1:]))
			***REMOVED*** else ***REMOVED***
				hasSingle = true
			***REMOVED***
		***REMOVED***
		if !hasSingle ***REMOVED***
			b.error(fmt.Errorf("no single entry for starter rune %U found", r))
			continue
		***REMOVED***
		// Unique the suffix set.
		sort.Strings(sufx)
		key := strings.Join(sufx, "\n")
		handle, ok := b.ctHandle[key]
		if !ok ***REMOVED***
			var err error
			handle, err = appendTrie(&t.ContractTries, sufx)
			if err != nil ***REMOVED***
				b.error(err)
			***REMOVED***
			b.ctHandle[key] = handle
		***REMOVED***
		// Bucket sort entries in index order.
		es := make([]*entry, len(l))
		for _, e := range l ***REMOVED***
			var p, sn int
			if len(e.runes) > 1 ***REMOVED***
				str := []byte(string(e.runes[1:]))
				p, sn = lookup(&t.ContractTries, handle, str)
				if sn != len(str) ***REMOVED***
					log.Fatalf("%s: processContractions: unexpected length for '%X'; len=%d; want %d", o.id, e.runes, sn, len(str))
				***REMOVED***
			***REMOVED***
			if es[p] != nil ***REMOVED***
				log.Fatalf("%s: multiple contractions for position %d for rune %U", o.id, p, e.runes[0])
			***REMOVED***
			es[p] = e
		***REMOVED***
		// Create collation elements for contractions.
		elems := []uint32***REMOVED******REMOVED***
		for _, e := range es ***REMOVED***
			ce, err := e.encodeBase()
			b.errorID(o.id, err)
			elems = append(elems, ce)
		***REMOVED***
		key = fmt.Sprintf("%v", elems)
		i, ok := b.ctElem[key]
		if !ok ***REMOVED***
			i = len(t.ContractElem)
			b.ctElem[key] = i
			t.ContractElem = append(t.ContractElem, elems...)
		***REMOVED***
		// Store info in entry for starter rune.
		es[0].contractionIndex = i
		es[0].contractionHandle = handle
	***REMOVED***
***REMOVED***