// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
	"fmt"
	"go/build"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/text/collate"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/internal"
	"golang.org/x/text/internal/catmsg"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/loader"
)

var transRe = regexp.MustCompile(`messages\.(.*)\.json`)

// Generate writes a Go file that defines a Catalog with translated messages.
// Translations are retrieved from s.Messages, not s.Translations, so it
// is assumed Merge has been called.
func (s *State) Generate() error ***REMOVED***
	path := s.Config.GenPackage
	if path == "" ***REMOVED***
		path = "."
	***REMOVED***
	isDir := path[0] == '.'
	prog, err := loadPackages(&loader.Config***REMOVED******REMOVED***, []string***REMOVED***path***REMOVED***)
	if err != nil ***REMOVED***
		return wrap(err, "could not load package")
	***REMOVED***
	pkgs := prog.InitialPackages()
	if len(pkgs) != 1 ***REMOVED***
		return errorf("more than one package selected: %v", pkgs)
	***REMOVED***
	pkg := pkgs[0].Pkg.Name()

	cw, err := s.generate()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !isDir ***REMOVED***
		gopath := build.Default.GOPATH
		path = filepath.Join(gopath, filepath.FromSlash(pkgs[0].Pkg.Path()))
	***REMOVED***
	path = filepath.Join(path, s.Config.GenFile)
	cw.WriteGoFile(path, pkg) // TODO: WriteGoFile should return error.
	return err
***REMOVED***

// WriteGen writes a Go file with the given package name to w that defines a
// Catalog with translated messages. Translations are retrieved from s.Messages,
// not s.Translations, so it is assumed Merge has been called.
func (s *State) WriteGen(w io.Writer, pkg string) error ***REMOVED***
	cw, err := s.generate()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = cw.WriteGo(w, pkg, "")
	return err
***REMOVED***

// Generate is deprecated; use (*State).Generate().
func Generate(w io.Writer, pkg string, extracted *Messages, trans ...Messages) (n int, err error) ***REMOVED***
	s := State***REMOVED***
		Extracted:    *extracted,
		Translations: trans,
	***REMOVED***
	cw, err := s.generate()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return cw.WriteGo(w, pkg, "")
***REMOVED***

func (s *State) generate() (*gen.CodeWriter, error) ***REMOVED***
	// Build up index of translations and original messages.
	translations := map[language.Tag]map[string]Message***REMOVED******REMOVED***
	languages := []language.Tag***REMOVED******REMOVED***
	usedKeys := map[string]int***REMOVED******REMOVED***

	for _, loc := range s.Messages ***REMOVED***
		tag := loc.Language
		if _, ok := translations[tag]; !ok ***REMOVED***
			translations[tag] = map[string]Message***REMOVED******REMOVED***
			languages = append(languages, tag)
		***REMOVED***
		for _, m := range loc.Messages ***REMOVED***
			if !m.Translation.IsEmpty() ***REMOVED***
				for _, id := range m.ID ***REMOVED***
					if _, ok := translations[tag][id]; ok ***REMOVED***
						warnf("Duplicate translation in locale %q for message %q", tag, id)
					***REMOVED***
					translations[tag][id] = m
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Verify completeness and register keys.
	internal.SortTags(languages)

	langVars := []string***REMOVED******REMOVED***
	for _, tag := range languages ***REMOVED***
		langVars = append(langVars, strings.Replace(tag.String(), "-", "_", -1))
		dict := translations[tag]
		for _, msg := range s.Extracted.Messages ***REMOVED***
			for _, id := range msg.ID ***REMOVED***
				if trans, ok := dict[id]; ok && !trans.Translation.IsEmpty() ***REMOVED***
					if _, ok := usedKeys[msg.Key]; !ok ***REMOVED***
						usedKeys[msg.Key] = len(usedKeys)
					***REMOVED***
					break
				***REMOVED***
				// TODO: log missing entry.
				warnf("%s: Missing entry for %q.", tag, id)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	cw := gen.NewCodeWriter()

	x := &struct ***REMOVED***
		Fallback  language.Tag
		Languages []string
	***REMOVED******REMOVED***
		Fallback:  s.Extracted.Language,
		Languages: langVars,
	***REMOVED***

	if err := lookup.Execute(cw, x); err != nil ***REMOVED***
		return nil, wrap(err, "error")
	***REMOVED***

	keyToIndex := []string***REMOVED******REMOVED***
	for k := range usedKeys ***REMOVED***
		keyToIndex = append(keyToIndex, k)
	***REMOVED***
	sort.Strings(keyToIndex)
	fmt.Fprint(cw, "var messageKeyToIndex = map[string]int***REMOVED***\n")
	for _, k := range keyToIndex ***REMOVED***
		fmt.Fprintf(cw, "%q: %d,\n", k, usedKeys[k])
	***REMOVED***
	fmt.Fprint(cw, "***REMOVED***\n\n")

	for i, tag := range languages ***REMOVED***
		dict := translations[tag]
		a := make([]string, len(usedKeys))
		for _, msg := range s.Extracted.Messages ***REMOVED***
			for _, id := range msg.ID ***REMOVED***
				if trans, ok := dict[id]; ok && !trans.Translation.IsEmpty() ***REMOVED***
					m, err := assemble(&msg, &trans.Translation)
					if err != nil ***REMOVED***
						return nil, wrap(err, "error")
					***REMOVED***
					_, leadWS, trailWS := trimWS(msg.Key)
					if leadWS != "" || trailWS != "" ***REMOVED***
						m = catmsg.Affix***REMOVED***
							Message: m,
							Prefix:  leadWS,
							Suffix:  trailWS,
						***REMOVED***
					***REMOVED***
					// TODO: support macros.
					data, err := catmsg.Compile(tag, nil, m)
					if err != nil ***REMOVED***
						return nil, wrap(err, "error")
					***REMOVED***
					key := usedKeys[msg.Key]
					if d := a[key]; d != "" && d != data ***REMOVED***
						warnf("Duplicate non-consistent translation for key %q, picking the one for message %q", msg.Key, id)
					***REMOVED***
					a[key] = string(data)
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		index := []uint32***REMOVED***0***REMOVED***
		p := 0
		for _, s := range a ***REMOVED***
			p += len(s)
			index = append(index, uint32(p))
		***REMOVED***

		cw.WriteVar(langVars[i]+"Index", index)
		cw.WriteConst(langVars[i]+"Data", strings.Join(a, ""))
	***REMOVED***
	return cw, nil
***REMOVED***

func assemble(m *Message, t *Text) (msg catmsg.Message, err error) ***REMOVED***
	keys := []string***REMOVED******REMOVED***
	for k := range t.Var ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	sort.Strings(keys)
	var a []catmsg.Message
	for _, k := range keys ***REMOVED***
		t := t.Var[k]
		m, err := assemble(m, &t)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		a = append(a, &catmsg.Var***REMOVED***Name: k, Message: m***REMOVED***)
	***REMOVED***
	if t.Select != nil ***REMOVED***
		s, err := assembleSelect(m, t.Select)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		a = append(a, s)
	***REMOVED***
	if t.Msg != "" ***REMOVED***
		sub, err := m.Substitute(t.Msg)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		a = append(a, catmsg.String(sub))
	***REMOVED***
	switch len(a) ***REMOVED***
	case 0:
		return nil, errorf("generate: empty message")
	case 1:
		return a[0], nil
	default:
		return catmsg.FirstOf(a), nil

	***REMOVED***
***REMOVED***

func assembleSelect(m *Message, s *Select) (msg catmsg.Message, err error) ***REMOVED***
	cases := []string***REMOVED******REMOVED***
	for c := range s.Cases ***REMOVED***
		cases = append(cases, c)
	***REMOVED***
	sortCases(cases)

	caseMsg := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, c := range cases ***REMOVED***
		cm := s.Cases[c]
		m, err := assemble(m, &cm)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		caseMsg = append(caseMsg, c, m)
	***REMOVED***

	ph := m.Placeholder(s.Arg)

	switch s.Feature ***REMOVED***
	case "plural":
		// TODO: only printf-style selects are supported as of yet.
		return plural.Selectf(ph.ArgNum, ph.String, caseMsg...), nil
	***REMOVED***
	return nil, errorf("unknown feature type %q", s.Feature)
***REMOVED***

func sortCases(cases []string) ***REMOVED***
	// TODO: implement full interface.
	sort.Slice(cases, func(i, j int) bool ***REMOVED***
		if cases[j] == "other" && cases[i] != "other" ***REMOVED***
			return true
		***REMOVED***
		// the following code relies on '<' < '=' < any letter.
		return cmpNumeric(cases[i], cases[j]) == -1
	***REMOVED***)
***REMOVED***

var cmpNumeric = collate.New(language.Und, collate.Numeric).CompareString

var lookup = template.Must(template.New("gen").Parse(`
import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct ***REMOVED***
	index []uint32
	data  string
***REMOVED***

func (d *dictionary) Lookup(key string) (data string, ok bool) ***REMOVED***
	p := messageKeyToIndex[key]
	start, end := d.index[p], d.index[p+1]
	if start == end ***REMOVED***
		return "", false
	***REMOVED***
	return d.data[start:end], true
***REMOVED***

func init() ***REMOVED***
	dict := map[string]catalog.Dictionary***REMOVED***
		***REMOVED******REMOVED***range .Languages***REMOVED******REMOVED***"***REMOVED******REMOVED***.***REMOVED******REMOVED***": &dictionary***REMOVED***index: ***REMOVED******REMOVED***.***REMOVED******REMOVED***Index, data: ***REMOVED******REMOVED***.***REMOVED******REMOVED***Data ***REMOVED***,
		***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED***
	fallback := language.MustParse("***REMOVED******REMOVED***.Fallback***REMOVED******REMOVED***")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	message.DefaultCatalog = cat
***REMOVED***

`))
