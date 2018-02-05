// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Encoding specifies encoding of the input data.
type Encoding uint

const (
	// UTF8 interprets the input data as UTF-8.
	UTF8 Encoding = 1 << iota

	// ISO_8859_1 interprets the input data as ISO-8859-1.
	ISO_8859_1
)

// Load reads a buffer into a Properties struct.
func Load(buf []byte, enc Encoding) (*Properties, error) ***REMOVED***
	return loadBuf(buf, enc)
***REMOVED***

// LoadString reads an UTF8 string into a properties struct.
func LoadString(s string) (*Properties, error) ***REMOVED***
	return loadBuf([]byte(s), UTF8)
***REMOVED***

// LoadMap creates a new Properties struct from a string map.
func LoadMap(m map[string]string) *Properties ***REMOVED***
	p := NewProperties()
	for k, v := range m ***REMOVED***
		p.Set(k, v)
	***REMOVED***
	return p
***REMOVED***

// LoadFile reads a file into a Properties struct.
func LoadFile(filename string, enc Encoding) (*Properties, error) ***REMOVED***
	return loadAll([]string***REMOVED***filename***REMOVED***, enc, false)
***REMOVED***

// LoadFiles reads multiple files in the given order into
// a Properties struct. If 'ignoreMissing' is true then
// non-existent files will not be reported as error.
func LoadFiles(filenames []string, enc Encoding, ignoreMissing bool) (*Properties, error) ***REMOVED***
	return loadAll(filenames, enc, ignoreMissing)
***REMOVED***

// LoadURL reads the content of the URL into a Properties struct.
//
// The encoding is determined via the Content-Type header which
// should be set to 'text/plain'. If the 'charset' parameter is
// missing, 'iso-8859-1' or 'latin1' the encoding is set to
// ISO-8859-1. If the 'charset' parameter is set to 'utf-8' the
// encoding is set to UTF-8. A missing content type header is
// interpreted as 'text/plain; charset=utf-8'.
func LoadURL(url string) (*Properties, error) ***REMOVED***
	return loadAll([]string***REMOVED***url***REMOVED***, UTF8, false)
***REMOVED***

// LoadURLs reads the content of multiple URLs in the given order into a
// Properties struct. If 'ignoreMissing' is true then a 404 status code will
// not be reported as error. See LoadURL for the Content-Type header
// and the encoding.
func LoadURLs(urls []string, ignoreMissing bool) (*Properties, error) ***REMOVED***
	return loadAll(urls, UTF8, ignoreMissing)
***REMOVED***

// LoadAll reads the content of multiple URLs or files in the given order into a
// Properties struct. If 'ignoreMissing' is true then a 404 status code or missing file will
// not be reported as error. Encoding sets the encoding for files. For the URLs please see
// LoadURL for the Content-Type header and the encoding.
func LoadAll(names []string, enc Encoding, ignoreMissing bool) (*Properties, error) ***REMOVED***
	return loadAll(names, enc, ignoreMissing)
***REMOVED***

// MustLoadString reads an UTF8 string into a Properties struct and
// panics on error.
func MustLoadString(s string) *Properties ***REMOVED***
	return must(LoadString(s))
***REMOVED***

// MustLoadFile reads a file into a Properties struct and
// panics on error.
func MustLoadFile(filename string, enc Encoding) *Properties ***REMOVED***
	return must(LoadFile(filename, enc))
***REMOVED***

// MustLoadFiles reads multiple files in the given order into
// a Properties struct and panics on error. If 'ignoreMissing'
// is true then non-existent files will not be reported as error.
func MustLoadFiles(filenames []string, enc Encoding, ignoreMissing bool) *Properties ***REMOVED***
	return must(LoadFiles(filenames, enc, ignoreMissing))
***REMOVED***

// MustLoadURL reads the content of a URL into a Properties struct and
// panics on error.
func MustLoadURL(url string) *Properties ***REMOVED***
	return must(LoadURL(url))
***REMOVED***

// MustLoadURLs reads the content of multiple URLs in the given order into a
// Properties struct and panics on error. If 'ignoreMissing' is true then a 404
// status code will not be reported as error.
func MustLoadURLs(urls []string, ignoreMissing bool) *Properties ***REMOVED***
	return must(LoadURLs(urls, ignoreMissing))
***REMOVED***

// MustLoadAll reads the content of multiple URLs or files in the given order into a
// Properties struct. If 'ignoreMissing' is true then a 404 status code or missing file will
// not be reported as error. Encoding sets the encoding for files. For the URLs please see
// LoadURL for the Content-Type header and the encoding. It panics on error.
func MustLoadAll(names []string, enc Encoding, ignoreMissing bool) *Properties ***REMOVED***
	return must(LoadAll(names, enc, ignoreMissing))
***REMOVED***

func loadBuf(buf []byte, enc Encoding) (*Properties, error) ***REMOVED***
	p, err := parse(convert(buf, enc))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p, p.check()
***REMOVED***

func loadAll(names []string, enc Encoding, ignoreMissing bool) (*Properties, error) ***REMOVED***
	result := NewProperties()
	for _, name := range names ***REMOVED***
		n, err := expandName(name)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var p *Properties
		if strings.HasPrefix(n, "http://") || strings.HasPrefix(n, "https://") ***REMOVED***
			p, err = loadURL(n, ignoreMissing)
		***REMOVED*** else ***REMOVED***
			p, err = loadFile(n, enc, ignoreMissing)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		result.Merge(p)

	***REMOVED***
	return result, result.check()
***REMOVED***

func loadFile(filename string, enc Encoding, ignoreMissing bool) (*Properties, error) ***REMOVED***
	data, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		if ignoreMissing && os.IsNotExist(err) ***REMOVED***
			LogPrintf("properties: %s not found. skipping", filename)
			return NewProperties(), nil
		***REMOVED***
		return nil, err
	***REMOVED***
	p, err := parse(convert(data, enc))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p, nil
***REMOVED***

func loadURL(url string, ignoreMissing bool) (*Properties, error) ***REMOVED***
	resp, err := http.Get(url)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("properties: error fetching %q. %s", url, err)
	***REMOVED***
	if resp.StatusCode == 404 && ignoreMissing ***REMOVED***
		LogPrintf("properties: %s returned %d. skipping", url, resp.StatusCode)
		return NewProperties(), nil
	***REMOVED***
	if resp.StatusCode != 200 ***REMOVED***
		return nil, fmt.Errorf("properties: %s returned %d", url, resp.StatusCode)
	***REMOVED***
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("properties: %s error reading response. %s", url, err)
	***REMOVED***
	if err = resp.Body.Close(); err != nil ***REMOVED***
		return nil, fmt.Errorf("properties: %s error reading response. %s", url, err)
	***REMOVED***

	ct := resp.Header.Get("Content-Type")
	var enc Encoding
	switch strings.ToLower(ct) ***REMOVED***
	case "text/plain", "text/plain; charset=iso-8859-1", "text/plain; charset=latin1":
		enc = ISO_8859_1
	case "", "text/plain; charset=utf-8":
		enc = UTF8
	default:
		return nil, fmt.Errorf("properties: invalid content type %s", ct)
	***REMOVED***

	p, err := parse(convert(body, enc))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p, nil
***REMOVED***

func must(p *Properties, err error) *Properties ***REMOVED***
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return p
***REMOVED***

// expandName expands $***REMOVED***ENV_VAR***REMOVED*** expressions in a name.
// If the environment variable does not exist then it will be replaced
// with an empty string. Malformed expressions like "$***REMOVED***ENV_VAR" will
// be reported as error.
func expandName(name string) (string, error) ***REMOVED***
	return expand(name, make(map[string]bool), "$***REMOVED***", "***REMOVED***", make(map[string]string))
***REMOVED***

// Interprets a byte buffer either as an ISO-8859-1 or UTF-8 encoded string.
// For ISO-8859-1 we can convert each byte straight into a rune since the
// first 256 unicode code points cover ISO-8859-1.
func convert(buf []byte, enc Encoding) string ***REMOVED***
	switch enc ***REMOVED***
	case UTF8:
		return string(buf)
	case ISO_8859_1:
		runes := make([]rune, len(buf))
		for i, b := range buf ***REMOVED***
			runes[i] = rune(b)
		***REMOVED***
		return string(runes)
	default:
		ErrorHandler(fmt.Errorf("unsupported encoding %v", enc))
	***REMOVED***
	panic("ErrorHandler should exit")
***REMOVED***
