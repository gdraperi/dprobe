// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pipeline provides tools for creating translation pipelines.
//
// NOTE: UNDER DEVELOPMENT. API MAY CHANGE.
package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"go/parser"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/internal"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/tools/go/loader"
)

const (
	extractFile  = "extracted.gotext.json"
	outFile      = "out.gotext.json"
	gotextSuffix = "gotext.json"
)

// Config contains configuration for the translation pipeline.
type Config struct ***REMOVED***
	// Supported indicates the languages for which data should be generated.
	// The default is to support all locales for which there are matching
	// translation files.
	Supported []language.Tag

	// --- Extraction

	SourceLanguage language.Tag

	Packages []string

	// --- File structure

	// Dir is the root dir for all operations.
	Dir string

	// TranslationsPattern is a regular expression to match incoming translation
	// files. These files may appear in any directory rooted at Dir.
	// language for the translation files is determined as follows:
	//   1. From the Language field in the file.
	//   2. If not present, from a valid language tag in the filename, separated
	//      by dots (e.g. "en-US.json" or "incoming.pt_PT.xmb").
	//   3. If not present, from a the closest subdirectory in which the file
	//      is contained that parses as a valid language tag.
	TranslationsPattern string

	// OutPattern defines the location for translation files for a certain
	// language. The default is "***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/***REMOVED******REMOVED***.Language***REMOVED******REMOVED***/out.***REMOVED******REMOVED***.Ext***REMOVED******REMOVED***"
	OutPattern string

	// Format defines the file format for generated translation files.
	// The default is XMB. Alternatives are GetText, XLIFF, L20n, GoText.
	Format string

	Ext string

	// TODO:
	// Actions are additional actions to be performed after the initial extract
	// and merge.
	// Actions []struct ***REMOVED***
	// 	Name    string
	// 	Options map[string]string
	// ***REMOVED***

	// --- Generation

	// GenFile may be in a different package. It is not defined, it will
	// be written to stdout.
	GenFile string

	// GenPackage is the package or relative path into which to generate the
	// file. If not specified it is relative to the current directory.
	GenPackage string

	// DeclareVar defines a variable to which to assing the generated Catalog.
	DeclareVar string

	// SetDefault determines whether to assign the generated Catalog to
	// message.DefaultCatalog. The default for this is true if DeclareVar is
	// not defined, false otherwise.
	SetDefault bool

	// TODO:
	// - Printf-style configuration
	// - Template-style configuration
	// - Extraction options
	// - Rewrite options
	// - Generation options
***REMOVED***

// Operations:
// - extract:       get the strings
// - disambiguate:  find messages with the same key, but possible different meaning.
// - create out:    create a list of messages that need translations
// - load trans:    load the list of current translations
// - merge:         assign list of translations as done
// - (action)expand:    analyze features and create example sentences for each version.
// - (action)googletrans:   pre-populate messages with automatic translations.
// - (action)export:    send out messages somewhere non-standard
// - (action)import:    load messages from somewhere non-standard
// - vet program:   don't pass "foo" + var + "bar" strings. Not using funcs for translated strings.
// - vet trans:     coverage: all translations/ all features.
// - generate:      generate Go code

// State holds all accumulated information on translations during processing.
type State struct ***REMOVED***
	Config Config

	Package string
	program *loader.Program

	Extracted Messages `json:"messages"`

	// Messages includes all messages for which there need to be translations.
	// Duplicates may be eliminated. Generation will be done from these messages
	// (usually after merging).
	Messages []Messages

	// Translations are incoming translations for the application messages.
	Translations []Messages
***REMOVED***

func (s *State) dir() string ***REMOVED***
	if d := s.Config.Dir; d != "" ***REMOVED***
		return d
	***REMOVED***
	return "./locales"
***REMOVED***

func outPattern(s *State) (string, error) ***REMOVED***
	c := s.Config
	pat := c.OutPattern
	if pat == "" ***REMOVED***
		pat = "***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/***REMOVED******REMOVED***.Language***REMOVED******REMOVED***/out.***REMOVED******REMOVED***.Ext***REMOVED******REMOVED***"
	***REMOVED***

	ext := c.Ext
	if ext == "" ***REMOVED***
		ext = c.Format
	***REMOVED***
	if ext == "" ***REMOVED***
		ext = gotextSuffix
	***REMOVED***
	t, err := template.New("").Parse(pat)
	if err != nil ***REMOVED***
		return "", wrap(err, "error parsing template")
	***REMOVED***
	buf := bytes.Buffer***REMOVED******REMOVED***
	err = t.Execute(&buf, map[string]string***REMOVED***
		"Dir":      s.dir(),
		"Language": "%s",
		"Ext":      ext,
	***REMOVED***)
	return filepath.FromSlash(buf.String()), wrap(err, "incorrect OutPattern")
***REMOVED***

var transRE = regexp.MustCompile(`.*\.` + gotextSuffix)

// Import loads existing translation files.
func (s *State) Import() error ***REMOVED***
	outPattern, err := outPattern(s)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	re := transRE
	if pat := s.Config.TranslationsPattern; pat != "" ***REMOVED***
		if re, err = regexp.Compile(pat); err != nil ***REMOVED***
			return wrapf(err, "error parsing regexp %q", s.Config.TranslationsPattern)
		***REMOVED***
	***REMOVED***
	x := importer***REMOVED***s, outPattern, re***REMOVED***
	return x.walkImport(s.dir(), s.Config.SourceLanguage)
***REMOVED***

type importer struct ***REMOVED***
	state      *State
	outPattern string
	transFile  *regexp.Regexp
***REMOVED***

func (i *importer) walkImport(path string, tag language.Tag) error ***REMOVED***
	files, err := ioutil.ReadDir(path)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	for _, f := range files ***REMOVED***
		name := f.Name()
		tag := tag
		if f.IsDir() ***REMOVED***
			if t, err := language.Parse(name); err == nil ***REMOVED***
				tag = t
			***REMOVED***
			// We ignore errors
			if err := i.walkImport(filepath.Join(path, name), tag); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***
		for _, l := range strings.Split(name, ".") ***REMOVED***
			if t, err := language.Parse(l); err == nil ***REMOVED***
				tag = t
			***REMOVED***
		***REMOVED***
		file := filepath.Join(path, name)
		// TODO: Should we skip files that match output files?
		if fmt.Sprintf(i.outPattern, tag) == file ***REMOVED***
			continue
		***REMOVED***
		// TODO: handle different file formats.
		if !i.transFile.MatchString(name) ***REMOVED***
			continue
		***REMOVED***
		b, err := ioutil.ReadFile(file)
		if err != nil ***REMOVED***
			return wrap(err, "read file failed")
		***REMOVED***
		var translations Messages
		if err := json.Unmarshal(b, &translations); err != nil ***REMOVED***
			return wrap(err, "parsing translation file failed")
		***REMOVED***
		i.state.Translations = append(i.state.Translations, translations)
	***REMOVED***
	return nil
***REMOVED***

// Merge merges the extracted messages with the existing translations.
func (s *State) Merge() error ***REMOVED***
	if s.Messages != nil ***REMOVED***
		panic("already merged")
	***REMOVED***
	// Create an index for each unique message.
	// Duplicates are okay as long as the substitution arguments are okay as
	// well.
	// Top-level messages are okay to appear in multiple substitution points.

	// Collect key equivalence.
	msgs := []*Message***REMOVED******REMOVED***
	keyToIDs := map[string]*Message***REMOVED******REMOVED***
	for _, m := range s.Extracted.Messages ***REMOVED***
		m := m
		if prev, ok := keyToIDs[m.Key]; ok ***REMOVED***
			if err := checkEquivalence(&m, prev); err != nil ***REMOVED***
				warnf("Key %q matches conflicting messages: %v and %v", m.Key, prev.ID, m.ID)
				// TODO: track enough information so that the rewriter can
				// suggest/disambiguate messages.
			***REMOVED***
			// TODO: add position to message.
			continue
		***REMOVED***
		i := len(msgs)
		msgs = append(msgs, &m)
		keyToIDs[m.Key] = msgs[i]
	***REMOVED***

	// Messages with different keys may still refer to the same translated
	// message (e.g. different whitespace). Filter these.
	idMap := map[string]bool***REMOVED******REMOVED***
	filtered := []*Message***REMOVED******REMOVED***
	for _, m := range msgs ***REMOVED***
		found := false
		for _, id := range m.ID ***REMOVED***
			found = found || idMap[id]
		***REMOVED***
		if !found ***REMOVED***
			filtered = append(filtered, m)
		***REMOVED***
		for _, id := range m.ID ***REMOVED***
			idMap[id] = true
		***REMOVED***
	***REMOVED***

	// Build index of translations.
	translations := map[language.Tag]map[string]Message***REMOVED******REMOVED***
	languages := append([]language.Tag***REMOVED******REMOVED***, s.Config.Supported...)

	for _, t := range s.Translations ***REMOVED***
		tag := t.Language
		if _, ok := translations[tag]; !ok ***REMOVED***
			translations[tag] = map[string]Message***REMOVED******REMOVED***
			languages = append(languages, tag)
		***REMOVED***
		for _, m := range t.Messages ***REMOVED***
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
	languages = internal.UniqueTags(languages)

	for _, tag := range languages ***REMOVED***
		ms := Messages***REMOVED***Language: tag***REMOVED***
		for _, orig := range filtered ***REMOVED***
			m := *orig
			m.Key = ""
			m.Position = ""

			for _, id := range m.ID ***REMOVED***
				if t, ok := translations[tag][id]; ok ***REMOVED***
					m.Translation = t.Translation
					if t.TranslatorComment != "" ***REMOVED***
						m.TranslatorComment = t.TranslatorComment
						m.Fuzzy = t.Fuzzy
					***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			if tag == s.Config.SourceLanguage && m.Translation.IsEmpty() ***REMOVED***
				m.Translation = m.Message
				if m.TranslatorComment == "" ***REMOVED***
					m.TranslatorComment = "Copied from source."
					m.Fuzzy = true
				***REMOVED***
			***REMOVED***
			// TODO: if translation is empty: pre-expand based on available
			// linguistic features. This may also be done as a plugin.
			ms.Messages = append(ms.Messages, m)
		***REMOVED***
		s.Messages = append(s.Messages, ms)
	***REMOVED***
	return nil
***REMOVED***

// Export writes out the messages to translation out files.
func (s *State) Export() error ***REMOVED***
	path, err := outPattern(s)
	if err != nil ***REMOVED***
		return wrap(err, "export failed")
	***REMOVED***
	for _, out := range s.Messages ***REMOVED***
		// TODO: inject translations from existing files to avoid retranslation.
		data, err := json.MarshalIndent(out, "", "    ")
		if err != nil ***REMOVED***
			return wrap(err, "JSON marshal failed")
		***REMOVED***
		file := fmt.Sprintf(path, out.Language)
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil ***REMOVED***
			return wrap(err, "dir create failed")
		***REMOVED***
		if err := ioutil.WriteFile(file, data, 0644); err != nil ***REMOVED***
			return wrap(err, "write failed")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var (
	ws    = runes.In(unicode.White_Space).Contains
	notWS = runes.NotIn(unicode.White_Space).Contains
)

func trimWS(s string) (trimmed, leadWS, trailWS string) ***REMOVED***
	trimmed = strings.TrimRightFunc(s, ws)
	trailWS = s[len(trimmed):]
	if i := strings.IndexFunc(trimmed, notWS); i > 0 ***REMOVED***
		leadWS = trimmed[:i]
		trimmed = trimmed[i:]
	***REMOVED***
	return trimmed, leadWS, trailWS
***REMOVED***

// NOTE: The command line tool already prefixes with "gotext:".
var (
	wrap = func(err error, msg string) error ***REMOVED***
		if err == nil ***REMOVED***
			return nil
		***REMOVED***
		return fmt.Errorf("%s: %v", msg, err)
	***REMOVED***
	wrapf = func(err error, msg string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
		if err == nil ***REMOVED***
			return nil
		***REMOVED***
		return wrap(err, fmt.Sprintf(msg, args...))
	***REMOVED***
	errorf = fmt.Errorf
)

func warnf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	// TODO: don't log.
	log.Printf(format, args...)
***REMOVED***

func loadPackages(conf *loader.Config, args []string) (*loader.Program, error) ***REMOVED***
	if len(args) == 0 ***REMOVED***
		args = []string***REMOVED***"."***REMOVED***
	***REMOVED***

	conf.Build = &build.Default
	conf.ParserMode = parser.ParseComments

	// Use the initial packages from the command line.
	args, err := conf.FromArgs(args, false)
	if err != nil ***REMOVED***
		return nil, wrap(err, "loading packages failed")
	***REMOVED***

	// Load, parse and type-check the whole program.
	return conf.Load()
***REMOVED***
