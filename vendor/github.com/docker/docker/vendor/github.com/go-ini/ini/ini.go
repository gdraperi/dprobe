// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package ini provides INI file read and write functionality in Go.
package ini

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// Name for default section. You can use this constant or the string literal.
	// In most of cases, an empty string is all you need to access the section.
	DEFAULT_SECTION = "DEFAULT"

	// Maximum allowed depth when recursively substituing variable names.
	_DEPTH_VALUES = 99
	_VERSION      = "1.25.4"
)

// Version returns current package version literal.
func Version() string ***REMOVED***
	return _VERSION
***REMOVED***

var (
	// Delimiter to determine or compose a new line.
	// This variable will be changed to "\r\n" automatically on Windows
	// at package init time.
	LineBreak = "\n"

	// Variable regexp pattern: %(variable)s
	varPattern = regexp.MustCompile(`%\(([^\)]+)\)s`)

	// Indicate whether to align "=" sign with spaces to produce pretty output
	// or reduce all possible spaces for compact format.
	PrettyFormat = true

	// Explicitly write DEFAULT section header
	DefaultHeader = false
)

func init() ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		LineBreak = "\r\n"
	***REMOVED***
***REMOVED***

func inSlice(str string, s []string) bool ***REMOVED***
	for _, v := range s ***REMOVED***
		if str == v ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// dataSource is an interface that returns object which can be read and closed.
type dataSource interface ***REMOVED***
	ReadCloser() (io.ReadCloser, error)
***REMOVED***

// sourceFile represents an object that contains content on the local file system.
type sourceFile struct ***REMOVED***
	name string
***REMOVED***

func (s sourceFile) ReadCloser() (_ io.ReadCloser, err error) ***REMOVED***
	return os.Open(s.name)
***REMOVED***

type bytesReadCloser struct ***REMOVED***
	reader io.Reader
***REMOVED***

func (rc *bytesReadCloser) Read(p []byte) (n int, err error) ***REMOVED***
	return rc.reader.Read(p)
***REMOVED***

func (rc *bytesReadCloser) Close() error ***REMOVED***
	return nil
***REMOVED***

// sourceData represents an object that contains content in memory.
type sourceData struct ***REMOVED***
	data []byte
***REMOVED***

func (s *sourceData) ReadCloser() (io.ReadCloser, error) ***REMOVED***
	return ioutil.NopCloser(bytes.NewReader(s.data)), nil
***REMOVED***

// sourceReadCloser represents an input stream with Close method.
type sourceReadCloser struct ***REMOVED***
	reader io.ReadCloser
***REMOVED***

func (s *sourceReadCloser) ReadCloser() (io.ReadCloser, error) ***REMOVED***
	return s.reader, nil
***REMOVED***

// File represents a combination of a or more INI file(s) in memory.
type File struct ***REMOVED***
	// Should make things safe, but sometimes doesn't matter.
	BlockMode bool
	// Make sure data is safe in multiple goroutines.
	lock sync.RWMutex

	// Allow combination of multiple data sources.
	dataSources []dataSource
	// Actual data is stored here.
	sections map[string]*Section

	// To keep data in order.
	sectionList []string

	options LoadOptions

	NameMapper
	ValueMapper
***REMOVED***

// newFile initializes File object with given data sources.
func newFile(dataSources []dataSource, opts LoadOptions) *File ***REMOVED***
	return &File***REMOVED***
		BlockMode:   true,
		dataSources: dataSources,
		sections:    make(map[string]*Section),
		sectionList: make([]string, 0, 10),
		options:     opts,
	***REMOVED***
***REMOVED***

func parseDataSource(source interface***REMOVED******REMOVED***) (dataSource, error) ***REMOVED***
	switch s := source.(type) ***REMOVED***
	case string:
		return sourceFile***REMOVED***s***REMOVED***, nil
	case []byte:
		return &sourceData***REMOVED***s***REMOVED***, nil
	case io.ReadCloser:
		return &sourceReadCloser***REMOVED***s***REMOVED***, nil
	default:
		return nil, fmt.Errorf("error parsing data source: unknown type '%s'", s)
	***REMOVED***
***REMOVED***

type LoadOptions struct ***REMOVED***
	// Loose indicates whether the parser should ignore nonexistent files or return error.
	Loose bool
	// Insensitive indicates whether the parser forces all section and key names to lowercase.
	Insensitive bool
	// IgnoreContinuation indicates whether to ignore continuation lines while parsing.
	IgnoreContinuation bool
	// AllowBooleanKeys indicates whether to allow boolean type keys or treat as value is missing.
	// This type of keys are mostly used in my.cnf.
	AllowBooleanKeys bool
	// AllowShadows indicates whether to keep track of keys with same name under same section.
	AllowShadows bool
	// Some INI formats allow group blocks that store a block of raw content that doesn't otherwise
	// conform to key/value pairs. Specify the names of those blocks here.
	UnparseableSections []string
***REMOVED***

func LoadSources(opts LoadOptions, source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) (_ *File, err error) ***REMOVED***
	sources := make([]dataSource, len(others)+1)
	sources[0], err = parseDataSource(source)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for i := range others ***REMOVED***
		sources[i+1], err = parseDataSource(others[i])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	f := newFile(sources, opts)
	if err = f.Reload(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return f, nil
***REMOVED***

// Load loads and parses from INI data sources.
// Arguments can be mixed of file name with string type, or raw data in []byte.
// It will return error if list contains nonexistent files.
func Load(source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) (*File, error) ***REMOVED***
	return LoadSources(LoadOptions***REMOVED******REMOVED***, source, others...)
***REMOVED***

// LooseLoad has exactly same functionality as Load function
// except it ignores nonexistent files instead of returning error.
func LooseLoad(source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) (*File, error) ***REMOVED***
	return LoadSources(LoadOptions***REMOVED***Loose: true***REMOVED***, source, others...)
***REMOVED***

// InsensitiveLoad has exactly same functionality as Load function
// except it forces all section and key names to be lowercased.
func InsensitiveLoad(source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) (*File, error) ***REMOVED***
	return LoadSources(LoadOptions***REMOVED***Insensitive: true***REMOVED***, source, others...)
***REMOVED***

// InsensitiveLoad has exactly same functionality as Load function
// except it allows have shadow keys.
func ShadowLoad(source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) (*File, error) ***REMOVED***
	return LoadSources(LoadOptions***REMOVED***AllowShadows: true***REMOVED***, source, others...)
***REMOVED***

// Empty returns an empty file object.
func Empty() *File ***REMOVED***
	// Ignore error here, we sure our data is good.
	f, _ := Load([]byte(""))
	return f
***REMOVED***

// NewSection creates a new section.
func (f *File) NewSection(name string) (*Section, error) ***REMOVED***
	if len(name) == 0 ***REMOVED***
		return nil, errors.New("error creating new section: empty section name")
	***REMOVED*** else if f.options.Insensitive && name != DEFAULT_SECTION ***REMOVED***
		name = strings.ToLower(name)
	***REMOVED***

	if f.BlockMode ***REMOVED***
		f.lock.Lock()
		defer f.lock.Unlock()
	***REMOVED***

	if inSlice(name, f.sectionList) ***REMOVED***
		return f.sections[name], nil
	***REMOVED***

	f.sectionList = append(f.sectionList, name)
	f.sections[name] = newSection(f, name)
	return f.sections[name], nil
***REMOVED***

// NewRawSection creates a new section with an unparseable body.
func (f *File) NewRawSection(name, body string) (*Section, error) ***REMOVED***
	section, err := f.NewSection(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	section.isRawSection = true
	section.rawBody = body
	return section, nil
***REMOVED***

// NewSections creates a list of sections.
func (f *File) NewSections(names ...string) (err error) ***REMOVED***
	for _, name := range names ***REMOVED***
		if _, err = f.NewSection(name); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetSection returns section by given name.
func (f *File) GetSection(name string) (*Section, error) ***REMOVED***
	if len(name) == 0 ***REMOVED***
		name = DEFAULT_SECTION
	***REMOVED*** else if f.options.Insensitive ***REMOVED***
		name = strings.ToLower(name)
	***REMOVED***

	if f.BlockMode ***REMOVED***
		f.lock.RLock()
		defer f.lock.RUnlock()
	***REMOVED***

	sec := f.sections[name]
	if sec == nil ***REMOVED***
		return nil, fmt.Errorf("section '%s' does not exist", name)
	***REMOVED***
	return sec, nil
***REMOVED***

// Section assumes named section exists and returns a zero-value when not.
func (f *File) Section(name string) *Section ***REMOVED***
	sec, err := f.GetSection(name)
	if err != nil ***REMOVED***
		// Note: It's OK here because the only possible error is empty section name,
		// but if it's empty, this piece of code won't be executed.
		sec, _ = f.NewSection(name)
		return sec
	***REMOVED***
	return sec
***REMOVED***

// Section returns list of Section.
func (f *File) Sections() []*Section ***REMOVED***
	sections := make([]*Section, len(f.sectionList))
	for i := range f.sectionList ***REMOVED***
		sections[i] = f.Section(f.sectionList[i])
	***REMOVED***
	return sections
***REMOVED***

// SectionStrings returns list of section names.
func (f *File) SectionStrings() []string ***REMOVED***
	list := make([]string, len(f.sectionList))
	copy(list, f.sectionList)
	return list
***REMOVED***

// DeleteSection deletes a section.
func (f *File) DeleteSection(name string) ***REMOVED***
	if f.BlockMode ***REMOVED***
		f.lock.Lock()
		defer f.lock.Unlock()
	***REMOVED***

	if len(name) == 0 ***REMOVED***
		name = DEFAULT_SECTION
	***REMOVED***

	for i, s := range f.sectionList ***REMOVED***
		if s == name ***REMOVED***
			f.sectionList = append(f.sectionList[:i], f.sectionList[i+1:]...)
			delete(f.sections, name)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *File) reload(s dataSource) error ***REMOVED***
	r, err := s.ReadCloser()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer r.Close()

	return f.parse(r)
***REMOVED***

// Reload reloads and parses all data sources.
func (f *File) Reload() (err error) ***REMOVED***
	for _, s := range f.dataSources ***REMOVED***
		if err = f.reload(s); err != nil ***REMOVED***
			// In loose mode, we create an empty default section for nonexistent files.
			if os.IsNotExist(err) && f.options.Loose ***REMOVED***
				f.parse(bytes.NewBuffer(nil))
				continue
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Append appends one or more data sources and reloads automatically.
func (f *File) Append(source interface***REMOVED******REMOVED***, others ...interface***REMOVED******REMOVED***) error ***REMOVED***
	ds, err := parseDataSource(source)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.dataSources = append(f.dataSources, ds)
	for _, s := range others ***REMOVED***
		ds, err = parseDataSource(s)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		f.dataSources = append(f.dataSources, ds)
	***REMOVED***
	return f.Reload()
***REMOVED***

// WriteToIndent writes content into io.Writer with given indention.
// If PrettyFormat has been set to be true,
// it will align "=" sign with spaces under each section.
func (f *File) WriteToIndent(w io.Writer, indent string) (n int64, err error) ***REMOVED***
	equalSign := "="
	if PrettyFormat ***REMOVED***
		equalSign = " = "
	***REMOVED***

	// Use buffer to make sure target is safe until finish encoding.
	buf := bytes.NewBuffer(nil)
	for i, sname := range f.sectionList ***REMOVED***
		sec := f.Section(sname)
		if len(sec.Comment) > 0 ***REMOVED***
			if sec.Comment[0] != '#' && sec.Comment[0] != ';' ***REMOVED***
				sec.Comment = "; " + sec.Comment
			***REMOVED***
			if _, err = buf.WriteString(sec.Comment + LineBreak); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
		***REMOVED***

		if i > 0 || DefaultHeader ***REMOVED***
			if _, err = buf.WriteString("[" + sname + "]" + LineBreak); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Write nothing if default section is empty
			if len(sec.keyList) == 0 ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if sec.isRawSection ***REMOVED***
			if _, err = buf.WriteString(sec.rawBody); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			continue
		***REMOVED***

		// Count and generate alignment length and buffer spaces using the
		// longest key. Keys may be modifed if they contain certain characters so
		// we need to take that into account in our calculation.
		alignLength := 0
		if PrettyFormat ***REMOVED***
			for _, kname := range sec.keyList ***REMOVED***
				keyLength := len(kname)
				// First case will surround key by ` and second by """
				if strings.ContainsAny(kname, "\"=:") ***REMOVED***
					keyLength += 2
				***REMOVED*** else if strings.Contains(kname, "`") ***REMOVED***
					keyLength += 6
				***REMOVED***

				if keyLength > alignLength ***REMOVED***
					alignLength = keyLength
				***REMOVED***
			***REMOVED***
		***REMOVED***
		alignSpaces := bytes.Repeat([]byte(" "), alignLength)

	KEY_LIST:
		for _, kname := range sec.keyList ***REMOVED***
			key := sec.Key(kname)
			if len(key.Comment) > 0 ***REMOVED***
				if len(indent) > 0 && sname != DEFAULT_SECTION ***REMOVED***
					buf.WriteString(indent)
				***REMOVED***
				if key.Comment[0] != '#' && key.Comment[0] != ';' ***REMOVED***
					key.Comment = "; " + key.Comment
				***REMOVED***
				if _, err = buf.WriteString(key.Comment + LineBreak); err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED***

			if len(indent) > 0 && sname != DEFAULT_SECTION ***REMOVED***
				buf.WriteString(indent)
			***REMOVED***

			switch ***REMOVED***
			case key.isAutoIncrement:
				kname = "-"
			case strings.ContainsAny(kname, "\"=:"):
				kname = "`" + kname + "`"
			case strings.Contains(kname, "`"):
				kname = `"""` + kname + `"""`
			***REMOVED***

			for _, val := range key.ValueWithShadows() ***REMOVED***
				if _, err = buf.WriteString(kname); err != nil ***REMOVED***
					return 0, err
				***REMOVED***

				if key.isBooleanType ***REMOVED***
					if kname != sec.keyList[len(sec.keyList)-1] ***REMOVED***
						buf.WriteString(LineBreak)
					***REMOVED***
					continue KEY_LIST
				***REMOVED***

				// Write out alignment spaces before "=" sign
				if PrettyFormat ***REMOVED***
					buf.Write(alignSpaces[:alignLength-len(kname)])
				***REMOVED***

				// In case key value contains "\n", "`", "\"", "#" or ";"
				if strings.ContainsAny(val, "\n`") ***REMOVED***
					val = `"""` + val + `"""`
				***REMOVED*** else if strings.ContainsAny(val, "#;") ***REMOVED***
					val = "`" + val + "`"
				***REMOVED***
				if _, err = buf.WriteString(equalSign + val + LineBreak); err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Put a line between sections
		if _, err = buf.WriteString(LineBreak); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***

	return buf.WriteTo(w)
***REMOVED***

// WriteTo writes file content into io.Writer.
func (f *File) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	return f.WriteToIndent(w, "")
***REMOVED***

// SaveToIndent writes content to file system with given value indention.
func (f *File) SaveToIndent(filename, indent string) error ***REMOVED***
	// Note: Because we are truncating with os.Create,
	// 	so it's safer to save to a temporary file location and rename afte done.
	tmpPath := filename + "." + strconv.Itoa(time.Now().Nanosecond()) + ".tmp"
	defer os.Remove(tmpPath)

	fw, err := os.Create(tmpPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err = f.WriteToIndent(fw, indent); err != nil ***REMOVED***
		fw.Close()
		return err
	***REMOVED***
	fw.Close()

	// Remove old file and rename the new one.
	os.Remove(filename)
	return os.Rename(tmpPath, filename)
***REMOVED***

// SaveTo writes content to file system.
func (f *File) SaveTo(filename string) error ***REMOVED***
	return f.SaveToIndent(filename, "")
***REMOVED***
