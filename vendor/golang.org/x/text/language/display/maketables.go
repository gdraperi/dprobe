// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Generator for display name tables.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var (
	test = flag.Bool("test", false,
		"test existing tables; can be used to compare web data with package data.")
	outputFile = flag.String("output", "tables.go", "output file")

	stats = flag.Bool("stats", false, "prints statistics to stderr")

	short = flag.Bool("short", false, `Use "short" alternatives, when available.`)
	draft = flag.String("draft",
		"contributed",
		`Minimal draft requirements (approved, contributed, provisional, unconfirmed).`)
	pkg = flag.String("package",
		"display",
		"the name of the package in which the generated file is to be included")

	tags = newTagSet("tags",
		[]language.Tag***REMOVED******REMOVED***,
		"space-separated list of tags to include or empty for all")
	dict = newTagSet("dict",
		dictTags(),
		"space-separated list or tags for which to include a Dictionary. "+
			`"" means the common list from go.text/language.`)
)

func dictTags() (tag []language.Tag) ***REMOVED***
	// TODO: replace with language.Common.Tags() once supported.
	const str = "af am ar ar-001 az bg bn ca cs da de el en en-US en-GB " +
		"es es-ES es-419 et fa fi fil fr fr-CA gu he hi hr hu hy id is it ja " +
		"ka kk km kn ko ky lo lt lv mk ml mn mr ms my ne nl no pa pl pt pt-BR " +
		"pt-PT ro ru si sk sl sq sr sr-Latn sv sw ta te th tr uk ur uz vi " +
		"zh zh-Hans zh-Hant zu"

	for _, s := range strings.Split(str, " ") ***REMOVED***
		tag = append(tag, language.MustParse(s))
	***REMOVED***
	return tag
***REMOVED***

func main() ***REMOVED***
	gen.Init()

	// Read the CLDR zip file.
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("main", "supplemental")
	d.SetSectionFilter("localeDisplayNames")
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		log.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	w := gen.NewCodeWriter()
	defer w.WriteGoFile(*outputFile, "display")

	gen.WriteCLDRVersion(w)

	b := builder***REMOVED***
		w:     w,
		data:  data,
		group: make(map[string]*group),
	***REMOVED***
	b.generate()
***REMOVED***

const tagForm = language.All

// tagSet is used to parse command line flags of tags. It implements the
// flag.Value interface.
type tagSet map[language.Tag]bool

func newTagSet(name string, tags []language.Tag, usage string) tagSet ***REMOVED***
	f := tagSet(make(map[language.Tag]bool))
	for _, t := range tags ***REMOVED***
		f[t] = true
	***REMOVED***
	flag.Var(f, name, usage)
	return f
***REMOVED***

// String implements the String method of the flag.Value interface.
func (f tagSet) String() string ***REMOVED***
	tags := []string***REMOVED******REMOVED***
	for t := range f ***REMOVED***
		tags = append(tags, t.String())
	***REMOVED***
	sort.Strings(tags)
	return strings.Join(tags, " ")
***REMOVED***

// Set implements Set from the flag.Value interface.
func (f tagSet) Set(s string) error ***REMOVED***
	if s != "" ***REMOVED***
		for _, s := range strings.Split(s, " ") ***REMOVED***
			if s != "" ***REMOVED***
				tag, err := tagForm.Parse(s)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				f[tag] = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (f tagSet) contains(t language.Tag) bool ***REMOVED***
	if len(f) == 0 ***REMOVED***
		return true
	***REMOVED***
	return f[t]
***REMOVED***

// builder is used to create all tables with display name information.
type builder struct ***REMOVED***
	w *gen.CodeWriter

	data *cldr.CLDR

	fromLocs []string

	// destination tags for the current locale.
	toTags     []string
	toTagIndex map[string]int

	// list of supported tags
	supported []language.Tag

	// key-value pairs per group
	group map[string]*group

	// statistics
	sizeIndex int // total size of all indexes of headers
	sizeData  int // total size of all data of headers
	totalSize int
***REMOVED***

type group struct ***REMOVED***
	// Maps from a given language to the Namer data for this language.
	lang    map[language.Tag]keyValues
	headers []header

	toTags        []string
	threeStart    int
	fourPlusStart int
***REMOVED***

// set sets the typ to the name for locale loc.
func (g *group) set(t language.Tag, typ, name string) ***REMOVED***
	kv := g.lang[t]
	if kv == nil ***REMOVED***
		kv = make(keyValues)
		g.lang[t] = kv
	***REMOVED***
	if kv[typ] == "" ***REMOVED***
		kv[typ] = name
	***REMOVED***
***REMOVED***

type keyValues map[string]string

type header struct ***REMOVED***
	tag   language.Tag
	data  string
	index []uint16
***REMOVED***

var versionInfo = `// Version is deprecated. Use CLDRVersion.
const Version = %#v

`

var self = language.MustParse("mul")

// generate builds and writes all tables.
func (b *builder) generate() ***REMOVED***
	fmt.Fprintf(b.w, versionInfo, cldr.Version)

	b.filter()
	b.setData("lang", func(g *group, loc language.Tag, ldn *cldr.LocaleDisplayNames) ***REMOVED***
		if ldn.Languages != nil ***REMOVED***
			for _, v := range ldn.Languages.Language ***REMOVED***
				lang := v.Type
				if lang == "root" ***REMOVED***
					// We prefer the data from "und"
					// TODO: allow both the data for root and und somehow.
					continue
				***REMOVED***
				tag := tagForm.MustParse(lang)
				if tags.contains(tag) ***REMOVED***
					g.set(loc, tag.String(), v.Data())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	b.setData("script", func(g *group, loc language.Tag, ldn *cldr.LocaleDisplayNames) ***REMOVED***
		if ldn.Scripts != nil ***REMOVED***
			for _, v := range ldn.Scripts.Script ***REMOVED***
				code := language.MustParseScript(v.Type)
				if code.IsPrivateUse() ***REMOVED*** // Qaaa..Qabx
					// TODO: data currently appears to be very meager.
					// Reconsider if we have data for English.
					if loc == language.English ***REMOVED***
						log.Fatal("Consider including data for private use scripts.")
					***REMOVED***
					continue
				***REMOVED***
				g.set(loc, code.String(), v.Data())
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	b.setData("region", func(g *group, loc language.Tag, ldn *cldr.LocaleDisplayNames) ***REMOVED***
		if ldn.Territories != nil ***REMOVED***
			for _, v := range ldn.Territories.Territory ***REMOVED***
				g.set(loc, language.MustParseRegion(v.Type).String(), v.Data())
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	b.makeSupported()

	b.writeParents()

	b.writeGroup("lang")
	b.writeGroup("script")
	b.writeGroup("region")

	b.w.WriteConst("numSupported", len(b.supported))
	buf := bytes.Buffer***REMOVED******REMOVED***
	for _, tag := range b.supported ***REMOVED***
		fmt.Fprint(&buf, tag.String(), "|")
	***REMOVED***
	b.w.WriteConst("supported", buf.String())

	b.writeDictionaries()

	b.supported = []language.Tag***REMOVED***self***REMOVED***

	// Compute the names of locales in their own language. Some of these names
	// may be specified in their parent locales. We iterate the maximum depth
	// of the parent three times to match successive parents of tags until a
	// possible match is found.
	for i := 0; i < 4; i++ ***REMOVED***
		b.setData("self", func(g *group, tag language.Tag, ldn *cldr.LocaleDisplayNames) ***REMOVED***
			parent := tag
			if b, s, r := tag.Raw(); i > 0 && (s != language.Script***REMOVED******REMOVED*** && r == language.Region***REMOVED******REMOVED***) ***REMOVED***
				parent, _ = language.Raw.Compose(b)
			***REMOVED***
			if ldn.Languages != nil ***REMOVED***
				for _, v := range ldn.Languages.Language ***REMOVED***
					key := tagForm.MustParse(v.Type)
					saved := key
					if key == parent ***REMOVED***
						g.set(self, tag.String(), v.Data())
					***REMOVED***
					for k := 0; k < i; k++ ***REMOVED***
						key = key.Parent()
					***REMOVED***
					if key == tag ***REMOVED***
						g.set(self, saved.String(), v.Data()) // set does not overwrite a value.
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	b.writeGroup("self")
***REMOVED***

func (b *builder) setData(name string, f func(*group, language.Tag, *cldr.LocaleDisplayNames)) ***REMOVED***
	b.sizeIndex = 0
	b.sizeData = 0
	b.toTags = nil
	b.fromLocs = nil
	b.toTagIndex = make(map[string]int)

	g := b.group[name]
	if g == nil ***REMOVED***
		g = &group***REMOVED***lang: make(map[language.Tag]keyValues)***REMOVED***
		b.group[name] = g
	***REMOVED***
	for _, loc := range b.data.Locales() ***REMOVED***
		// We use RawLDML instead of LDML as we are managing our own inheritance
		// in this implementation.
		ldml := b.data.RawLDML(loc)

		// We do not support the POSIX variant (it is not a supported BCP 47
		// variant). This locale also doesn't happen to contain any data, so
		// we'll skip it by checking for this.
		tag, err := tagForm.Parse(loc)
		if err != nil ***REMOVED***
			if ldml.LocaleDisplayNames != nil ***REMOVED***
				log.Fatalf("setData: %v", err)
			***REMOVED***
			continue
		***REMOVED***
		if ldml.LocaleDisplayNames != nil && tags.contains(tag) ***REMOVED***
			f(g, tag, ldml.LocaleDisplayNames)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *builder) filter() ***REMOVED***
	filter := func(s *cldr.Slice) ***REMOVED***
		if *short ***REMOVED***
			s.SelectOnePerGroup("alt", []string***REMOVED***"short", ""***REMOVED***)
		***REMOVED*** else ***REMOVED***
			s.SelectOnePerGroup("alt", []string***REMOVED***"stand-alone", ""***REMOVED***)
		***REMOVED***
		d, err := cldr.ParseDraft(*draft)
		if err != nil ***REMOVED***
			log.Fatalf("filter: %v", err)
		***REMOVED***
		s.SelectDraft(d)
	***REMOVED***
	for _, loc := range b.data.Locales() ***REMOVED***
		if ldn := b.data.RawLDML(loc).LocaleDisplayNames; ldn != nil ***REMOVED***
			if ldn.Languages != nil ***REMOVED***
				s := cldr.MakeSlice(&ldn.Languages.Language)
				if filter(&s); len(ldn.Languages.Language) == 0 ***REMOVED***
					ldn.Languages = nil
				***REMOVED***
			***REMOVED***
			if ldn.Scripts != nil ***REMOVED***
				s := cldr.MakeSlice(&ldn.Scripts.Script)
				if filter(&s); len(ldn.Scripts.Script) == 0 ***REMOVED***
					ldn.Scripts = nil
				***REMOVED***
			***REMOVED***
			if ldn.Territories != nil ***REMOVED***
				s := cldr.MakeSlice(&ldn.Territories.Territory)
				if filter(&s); len(ldn.Territories.Territory) == 0 ***REMOVED***
					ldn.Territories = nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// makeSupported creates a list of all supported locales.
func (b *builder) makeSupported() ***REMOVED***
	// tags across groups
	for _, g := range b.group ***REMOVED***
		for t, _ := range g.lang ***REMOVED***
			b.supported = append(b.supported, t)
		***REMOVED***
	***REMOVED***
	b.supported = b.supported[:unique(tagsSorter(b.supported))]

***REMOVED***

type tagsSorter []language.Tag

func (a tagsSorter) Len() int           ***REMOVED*** return len(a) ***REMOVED***
func (a tagsSorter) Swap(i, j int)      ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***
func (a tagsSorter) Less(i, j int) bool ***REMOVED*** return a[i].String() < a[j].String() ***REMOVED***

func (b *builder) writeGroup(name string) ***REMOVED***
	g := b.group[name]

	for _, kv := range g.lang ***REMOVED***
		for t, _ := range kv ***REMOVED***
			g.toTags = append(g.toTags, t)
		***REMOVED***
	***REMOVED***
	g.toTags = g.toTags[:unique(tagsBySize(g.toTags))]

	// Allocate header per supported value.
	g.headers = make([]header, len(b.supported))
	for i, sup := range b.supported ***REMOVED***
		kv, ok := g.lang[sup]
		if !ok ***REMOVED***
			g.headers[i].tag = sup
			continue
		***REMOVED***
		data := []byte***REMOVED******REMOVED***
		index := make([]uint16, len(g.toTags), len(g.toTags)+1)
		for j, t := range g.toTags ***REMOVED***
			index[j] = uint16(len(data))
			data = append(data, kv[t]...)
		***REMOVED***
		index = append(index, uint16(len(data)))

		// Trim the tail of the index.
		// TODO: indexes can be reduced in size quite a bit more.
		n := len(index)
		for ; n >= 2 && index[n-2] == index[n-1]; n-- ***REMOVED***
		***REMOVED***
		index = index[:n]

		// Workaround for a bug in CLDR 26.
		// See http://unicode.org/cldr/trac/ticket/8042.
		if cldr.Version == "26" && sup.String() == "hsb" ***REMOVED***
			data = bytes.Replace(data, []byte***REMOVED***'"'***REMOVED***, nil, 1)
		***REMOVED***
		g.headers[i] = header***REMOVED***sup, string(data), index***REMOVED***
	***REMOVED***
	g.writeTable(b.w, name)
***REMOVED***

type tagsBySize []string

func (l tagsBySize) Len() int      ***REMOVED*** return len(l) ***REMOVED***
func (l tagsBySize) Swap(i, j int) ***REMOVED*** l[i], l[j] = l[j], l[i] ***REMOVED***
func (l tagsBySize) Less(i, j int) bool ***REMOVED***
	a, b := l[i], l[j]
	// Sort single-tag entries based on size first. Otherwise alphabetic.
	if len(a) != len(b) && (len(a) <= 4 || len(b) <= 4) ***REMOVED***
		return len(a) < len(b)
	***REMOVED***
	return a < b
***REMOVED***

// parentIndices returns slice a of len(tags) where tags[a[i]] is the parent
// of tags[i].
func parentIndices(tags []language.Tag) []int16 ***REMOVED***
	index := make(map[language.Tag]int16)
	for i, t := range tags ***REMOVED***
		index[t] = int16(i)
	***REMOVED***

	// Construct default parents.
	parents := make([]int16, len(tags))
	for i, t := range tags ***REMOVED***
		parents[i] = -1
		for t = t.Parent(); t != language.Und; t = t.Parent() ***REMOVED***
			if j, ok := index[t]; ok ***REMOVED***
				parents[i] = j
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return parents
***REMOVED***

func (b *builder) writeParents() ***REMOVED***
	parents := parentIndices(b.supported)
	fmt.Fprintf(b.w, "var parents = ")
	b.w.WriteArray(parents)
***REMOVED***

// writeKeys writes keys to a special index used by the display package.
// tags are assumed to be sorted by length.
func writeKeys(w *gen.CodeWriter, name string, keys []string) ***REMOVED***
	w.Size += int(3 * reflect.TypeOf("").Size())
	w.WriteComment("Number of keys: %d", len(keys))
	fmt.Fprintf(w, "var (\n\t%sIndex = tagIndex***REMOVED***\n", name)
	for i := 2; i <= 4; i++ ***REMOVED***
		sub := []string***REMOVED******REMOVED***
		for _, t := range keys ***REMOVED***
			if len(t) != i ***REMOVED***
				break
			***REMOVED***
			sub = append(sub, t)
		***REMOVED***
		s := strings.Join(sub, "")
		w.WriteString(s)
		fmt.Fprintf(w, ",\n")
		keys = keys[len(sub):]
	***REMOVED***
	fmt.Fprintln(w, "\t***REMOVED***")
	if len(keys) > 0 ***REMOVED***
		w.Size += int(reflect.TypeOf([]string***REMOVED******REMOVED***).Size())
		fmt.Fprintf(w, "\t%sTagsLong = ", name)
		w.WriteSlice(keys)
	***REMOVED***
	fmt.Fprintln(w, ")\n")
***REMOVED***

// identifier creates an identifier from the given tag.
func identifier(t language.Tag) string ***REMOVED***
	return strings.Replace(t.String(), "-", "", -1)
***REMOVED***

func (h *header) writeEntry(w *gen.CodeWriter, name string) ***REMOVED***
	if len(dict) > 0 && dict.contains(h.tag) ***REMOVED***
		fmt.Fprintf(w, "\t***REMOVED*** // %s\n", h.tag)
		fmt.Fprintf(w, "\t\t%[1]s%[2]sStr,\n\t\t%[1]s%[2]sIdx,\n", identifier(h.tag), name)
		fmt.Fprintln(w, "\t***REMOVED***,")
	***REMOVED*** else if len(h.data) == 0 ***REMOVED***
		fmt.Fprintln(w, "\t\t***REMOVED******REMOVED***, //", h.tag)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(w, "\t***REMOVED*** // %s\n", h.tag)
		w.WriteString(h.data)
		fmt.Fprintln(w, ",")
		w.WriteSlice(h.index)
		fmt.Fprintln(w, ",\n\t***REMOVED***,")
	***REMOVED***
***REMOVED***

// write the data for the given header as single entries. The size for this data
// was already accounted for in writeEntry.
func (h *header) writeSingle(w *gen.CodeWriter, name string) ***REMOVED***
	if len(dict) > 0 && dict.contains(h.tag) ***REMOVED***
		tag := identifier(h.tag)
		w.WriteConst(tag+name+"Str", h.data)

		// Note that we create a slice instead of an array. If we use an array
		// we need to refer to it as a[:] in other tables, which will cause the
		// array to always be included by the linker. See Issue 7651.
		w.WriteVar(tag+name+"Idx", h.index)
	***REMOVED***
***REMOVED***

// WriteTable writes an entry for a single Namer.
func (g *group) writeTable(w *gen.CodeWriter, name string) ***REMOVED***
	start := w.Size
	writeKeys(w, name, g.toTags)
	w.Size += len(g.headers) * int(reflect.ValueOf(g.headers[0]).Type().Size())

	fmt.Fprintf(w, "var %sHeaders = [%d]header***REMOVED***\n", name, len(g.headers))

	title := strings.Title(name)
	for _, h := range g.headers ***REMOVED***
		h.writeEntry(w, title)
	***REMOVED***
	fmt.Fprintln(w, "***REMOVED***\n")

	for _, h := range g.headers ***REMOVED***
		h.writeSingle(w, title)
	***REMOVED***
	n := w.Size - start
	fmt.Fprintf(w, "// Total size for %s: %d bytes (%d KB)\n\n", name, n, n/1000)
***REMOVED***

func (b *builder) writeDictionaries() ***REMOVED***
	fmt.Fprintln(b.w, "// Dictionary entries of frequent languages")
	fmt.Fprintln(b.w, "var (")
	parents := parentIndices(b.supported)

	for i, t := range b.supported ***REMOVED***
		if dict.contains(t) ***REMOVED***
			ident := identifier(t)
			fmt.Fprintf(b.w, "\t%s = Dictionary***REMOVED*** // %s\n", ident, t)
			if p := parents[i]; p == -1 ***REMOVED***
				fmt.Fprintln(b.w, "\t\tnil,")
			***REMOVED*** else ***REMOVED***
				fmt.Fprintf(b.w, "\t\t&%s,\n", identifier(b.supported[p]))
			***REMOVED***
			fmt.Fprintf(b.w, "\t\theader***REMOVED***%[1]sLangStr, %[1]sLangIdx***REMOVED***,\n", ident)
			fmt.Fprintf(b.w, "\t\theader***REMOVED***%[1]sScriptStr, %[1]sScriptIdx***REMOVED***,\n", ident)
			fmt.Fprintf(b.w, "\t\theader***REMOVED***%[1]sRegionStr, %[1]sRegionIdx***REMOVED***,\n", ident)
			fmt.Fprintln(b.w, "\t***REMOVED***")
		***REMOVED***
	***REMOVED***
	fmt.Fprintln(b.w, ")")

	var s string
	var a []uint16
	sz := reflect.TypeOf(s).Size()
	sz += reflect.TypeOf(a).Size()
	sz *= 3
	sz += reflect.TypeOf(&a).Size()
	n := int(sz) * len(dict)
	fmt.Fprintf(b.w, "// Total size for %d entries: %d bytes (%d KB)\n\n", len(dict), n, n/1000)

	b.w.Size += n
***REMOVED***

// unique sorts the given lists and removes duplicate entries by swapping them
// past position k, where k is the number of unique values. It returns k.
func unique(a sort.Interface) int ***REMOVED***
	if a.Len() == 0 ***REMOVED***
		return 0
	***REMOVED***
	sort.Sort(a)
	k := 1
	for i := 1; i < a.Len(); i++ ***REMOVED***
		if a.Less(k-1, i) ***REMOVED***
			if k != i ***REMOVED***
				a.Swap(k, i)
			***REMOVED***
			k++
		***REMOVED***
	***REMOVED***
	return k
***REMOVED***
