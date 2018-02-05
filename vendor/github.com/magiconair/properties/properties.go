// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

// BUG(frank): Set() does not check for invalid unicode literals since this is currently handled by the lexer.
// BUG(frank): Write() does not allow to configure the newline character. Therefore, on Windows LF is used.

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// ErrorHandlerFunc defines the type of function which handles failures
// of the MustXXX() functions. An error handler function must exit
// the application after handling the error.
type ErrorHandlerFunc func(error)

// ErrorHandler is the function which handles failures of the MustXXX()
// functions. The default is LogFatalHandler.
var ErrorHandler ErrorHandlerFunc = LogFatalHandler

// LogHandlerFunc defines the function prototype for logging errors.
type LogHandlerFunc func(fmt string, args ...interface***REMOVED******REMOVED***)

// LogPrintf defines a log handler which uses log.Printf.
var LogPrintf LogHandlerFunc = log.Printf

// LogFatalHandler handles the error by logging a fatal error and exiting.
func LogFatalHandler(err error) ***REMOVED***
	log.Fatal(err)
***REMOVED***

// PanicHandler handles the error by panicking.
func PanicHandler(err error) ***REMOVED***
	panic(err)
***REMOVED***

// -----------------------------------------------------------------------------

// A Properties contains the key/value pairs from the properties input.
// All values are stored in unexpanded form and are expanded at runtime
type Properties struct ***REMOVED***
	// Pre-/Postfix for property expansion.
	Prefix  string
	Postfix string

	// DisableExpansion controls the expansion of properties on Get()
	// and the check for circular references on Set(). When set to
	// true Properties behaves like a simple key/value store and does
	// not check for circular references on Get() or on Set().
	DisableExpansion bool

	// Stores the key/value pairs
	m map[string]string

	// Stores the comments per key.
	c map[string][]string

	// Stores the keys in order of appearance.
	k []string
***REMOVED***

// NewProperties creates a new Properties struct with the default
// configuration for "$***REMOVED***key***REMOVED***" expressions.
func NewProperties() *Properties ***REMOVED***
	return &Properties***REMOVED***
		Prefix:  "$***REMOVED***",
		Postfix: "***REMOVED***",
		m:       map[string]string***REMOVED******REMOVED***,
		c:       map[string][]string***REMOVED******REMOVED***,
		k:       []string***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// Get returns the expanded value for the given key if exists.
// Otherwise, ok is false.
func (p *Properties) Get(key string) (value string, ok bool) ***REMOVED***
	v, ok := p.m[key]
	if p.DisableExpansion ***REMOVED***
		return v, ok
	***REMOVED***
	if !ok ***REMOVED***
		return "", false
	***REMOVED***

	expanded, err := p.expand(v)

	// we guarantee that the expanded value is free of
	// circular references and malformed expressions
	// so we panic if we still get an error here.
	if err != nil ***REMOVED***
		ErrorHandler(fmt.Errorf("%s in %q", err, key+" = "+v))
	***REMOVED***

	return expanded, true
***REMOVED***

// MustGet returns the expanded value for the given key if exists.
// Otherwise, it panics.
func (p *Properties) MustGet(key string) string ***REMOVED***
	if v, ok := p.Get(key); ok ***REMOVED***
		return v
	***REMOVED***
	ErrorHandler(invalidKeyError(key))
	panic("ErrorHandler should exit")
***REMOVED***

// ----------------------------------------------------------------------------

// ClearComments removes the comments for all keys.
func (p *Properties) ClearComments() ***REMOVED***
	p.c = map[string][]string***REMOVED******REMOVED***
***REMOVED***

// ----------------------------------------------------------------------------

// GetComment returns the last comment before the given key or an empty string.
func (p *Properties) GetComment(key string) string ***REMOVED***
	comments, ok := p.c[key]
	if !ok || len(comments) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return comments[len(comments)-1]
***REMOVED***

// ----------------------------------------------------------------------------

// GetComments returns all comments that appeared before the given key or nil.
func (p *Properties) GetComments(key string) []string ***REMOVED***
	if comments, ok := p.c[key]; ok ***REMOVED***
		return comments
	***REMOVED***
	return nil
***REMOVED***

// ----------------------------------------------------------------------------

// SetComment sets the comment for the key.
func (p *Properties) SetComment(key, comment string) ***REMOVED***
	p.c[key] = []string***REMOVED***comment***REMOVED***
***REMOVED***

// ----------------------------------------------------------------------------

// SetComments sets the comments for the key. If the comments are nil then
// all comments for this key are deleted.
func (p *Properties) SetComments(key string, comments []string) ***REMOVED***
	if comments == nil ***REMOVED***
		delete(p.c, key)
		return
	***REMOVED***
	p.c[key] = comments
***REMOVED***

// ----------------------------------------------------------------------------

// GetBool checks if the expanded value is one of '1', 'yes',
// 'true' or 'on' if the key exists. The comparison is case-insensitive.
// If the key does not exist the default value is returned.
func (p *Properties) GetBool(key string, def bool) bool ***REMOVED***
	v, err := p.getBool(key)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return v
***REMOVED***

// MustGetBool checks if the expanded value is one of '1', 'yes',
// 'true' or 'on' if the key exists. The comparison is case-insensitive.
// If the key does not exist the function panics.
func (p *Properties) MustGetBool(key string) bool ***REMOVED***
	v, err := p.getBool(key)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return v
***REMOVED***

func (p *Properties) getBool(key string) (value bool, err error) ***REMOVED***
	if v, ok := p.Get(key); ok ***REMOVED***
		return boolVal(v), nil
	***REMOVED***
	return false, invalidKeyError(key)
***REMOVED***

func boolVal(v string) bool ***REMOVED***
	v = strings.ToLower(v)
	return v == "1" || v == "true" || v == "yes" || v == "on"
***REMOVED***

// ----------------------------------------------------------------------------

// GetDuration parses the expanded value as an time.Duration (in ns) if the
// key exists. If key does not exist or the value cannot be parsed the default
// value is returned. In almost all cases you want to use GetParsedDuration().
func (p *Properties) GetDuration(key string, def time.Duration) time.Duration ***REMOVED***
	v, err := p.getInt64(key)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return time.Duration(v)
***REMOVED***

// MustGetDuration parses the expanded value as an time.Duration (in ns) if
// the key exists. If key does not exist or the value cannot be parsed the
// function panics. In almost all cases you want to use MustGetParsedDuration().
func (p *Properties) MustGetDuration(key string) time.Duration ***REMOVED***
	v, err := p.getInt64(key)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return time.Duration(v)
***REMOVED***

// ----------------------------------------------------------------------------

// GetParsedDuration parses the expanded value with time.ParseDuration() if the key exists.
// If key does not exist or the value cannot be parsed the default
// value is returned.
func (p *Properties) GetParsedDuration(key string, def time.Duration) time.Duration ***REMOVED***
	s, ok := p.Get(key)
	if !ok ***REMOVED***
		return def
	***REMOVED***
	v, err := time.ParseDuration(s)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return v
***REMOVED***

// MustGetParsedDuration parses the expanded value with time.ParseDuration() if the key exists.
// If key does not exist or the value cannot be parsed the function panics.
func (p *Properties) MustGetParsedDuration(key string) time.Duration ***REMOVED***
	s, ok := p.Get(key)
	if !ok ***REMOVED***
		ErrorHandler(invalidKeyError(key))
	***REMOVED***
	v, err := time.ParseDuration(s)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return v
***REMOVED***

// ----------------------------------------------------------------------------

// GetFloat64 parses the expanded value as a float64 if the key exists.
// If key does not exist or the value cannot be parsed the default
// value is returned.
func (p *Properties) GetFloat64(key string, def float64) float64 ***REMOVED***
	v, err := p.getFloat64(key)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return v
***REMOVED***

// MustGetFloat64 parses the expanded value as a float64 if the key exists.
// If key does not exist or the value cannot be parsed the function panics.
func (p *Properties) MustGetFloat64(key string) float64 ***REMOVED***
	v, err := p.getFloat64(key)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return v
***REMOVED***

func (p *Properties) getFloat64(key string) (value float64, err error) ***REMOVED***
	if v, ok := p.Get(key); ok ***REMOVED***
		value, err = strconv.ParseFloat(v, 64)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return value, nil
	***REMOVED***
	return 0, invalidKeyError(key)
***REMOVED***

// ----------------------------------------------------------------------------

// GetInt parses the expanded value as an int if the key exists.
// If key does not exist or the value cannot be parsed the default
// value is returned. If the value does not fit into an int the
// function panics with an out of range error.
func (p *Properties) GetInt(key string, def int) int ***REMOVED***
	v, err := p.getInt64(key)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return intRangeCheck(key, v)
***REMOVED***

// MustGetInt parses the expanded value as an int if the key exists.
// If key does not exist or the value cannot be parsed the function panics.
// If the value does not fit into an int the function panics with
// an out of range error.
func (p *Properties) MustGetInt(key string) int ***REMOVED***
	v, err := p.getInt64(key)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return intRangeCheck(key, v)
***REMOVED***

// ----------------------------------------------------------------------------

// GetInt64 parses the expanded value as an int64 if the key exists.
// If key does not exist or the value cannot be parsed the default
// value is returned.
func (p *Properties) GetInt64(key string, def int64) int64 ***REMOVED***
	v, err := p.getInt64(key)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return v
***REMOVED***

// MustGetInt64 parses the expanded value as an int if the key exists.
// If key does not exist or the value cannot be parsed the function panics.
func (p *Properties) MustGetInt64(key string) int64 ***REMOVED***
	v, err := p.getInt64(key)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return v
***REMOVED***

func (p *Properties) getInt64(key string) (value int64, err error) ***REMOVED***
	if v, ok := p.Get(key); ok ***REMOVED***
		value, err = strconv.ParseInt(v, 10, 64)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return value, nil
	***REMOVED***
	return 0, invalidKeyError(key)
***REMOVED***

// ----------------------------------------------------------------------------

// GetUint parses the expanded value as an uint if the key exists.
// If key does not exist or the value cannot be parsed the default
// value is returned. If the value does not fit into an int the
// function panics with an out of range error.
func (p *Properties) GetUint(key string, def uint) uint ***REMOVED***
	v, err := p.getUint64(key)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return uintRangeCheck(key, v)
***REMOVED***

// MustGetUint parses the expanded value as an int if the key exists.
// If key does not exist or the value cannot be parsed the function panics.
// If the value does not fit into an int the function panics with
// an out of range error.
func (p *Properties) MustGetUint(key string) uint ***REMOVED***
	v, err := p.getUint64(key)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return uintRangeCheck(key, v)
***REMOVED***

// ----------------------------------------------------------------------------

// GetUint64 parses the expanded value as an uint64 if the key exists.
// If key does not exist or the value cannot be parsed the default
// value is returned.
func (p *Properties) GetUint64(key string, def uint64) uint64 ***REMOVED***
	v, err := p.getUint64(key)
	if err != nil ***REMOVED***
		return def
	***REMOVED***
	return v
***REMOVED***

// MustGetUint64 parses the expanded value as an int if the key exists.
// If key does not exist or the value cannot be parsed the function panics.
func (p *Properties) MustGetUint64(key string) uint64 ***REMOVED***
	v, err := p.getUint64(key)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return v
***REMOVED***

func (p *Properties) getUint64(key string) (value uint64, err error) ***REMOVED***
	if v, ok := p.Get(key); ok ***REMOVED***
		value, err = strconv.ParseUint(v, 10, 64)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return value, nil
	***REMOVED***
	return 0, invalidKeyError(key)
***REMOVED***

// ----------------------------------------------------------------------------

// GetString returns the expanded value for the given key if exists or
// the default value otherwise.
func (p *Properties) GetString(key, def string) string ***REMOVED***
	if v, ok := p.Get(key); ok ***REMOVED***
		return v
	***REMOVED***
	return def
***REMOVED***

// MustGetString returns the expanded value for the given key if exists or
// panics otherwise.
func (p *Properties) MustGetString(key string) string ***REMOVED***
	if v, ok := p.Get(key); ok ***REMOVED***
		return v
	***REMOVED***
	ErrorHandler(invalidKeyError(key))
	panic("ErrorHandler should exit")
***REMOVED***

// ----------------------------------------------------------------------------

// Filter returns a new properties object which contains all properties
// for which the key matches the pattern.
func (p *Properties) Filter(pattern string) (*Properties, error) ***REMOVED***
	re, err := regexp.Compile(pattern)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return p.FilterRegexp(re), nil
***REMOVED***

// FilterRegexp returns a new properties object which contains all properties
// for which the key matches the regular expression.
func (p *Properties) FilterRegexp(re *regexp.Regexp) *Properties ***REMOVED***
	pp := NewProperties()
	for _, k := range p.k ***REMOVED***
		if re.MatchString(k) ***REMOVED***
			// TODO(fs): we are ignoring the error which flags a circular reference.
			// TODO(fs): since we are just copying a subset of keys this cannot happen (fingers crossed)
			pp.Set(k, p.m[k])
		***REMOVED***
	***REMOVED***
	return pp
***REMOVED***

// FilterPrefix returns a new properties object with a subset of all keys
// with the given prefix.
func (p *Properties) FilterPrefix(prefix string) *Properties ***REMOVED***
	pp := NewProperties()
	for _, k := range p.k ***REMOVED***
		if strings.HasPrefix(k, prefix) ***REMOVED***
			// TODO(fs): we are ignoring the error which flags a circular reference.
			// TODO(fs): since we are just copying a subset of keys this cannot happen (fingers crossed)
			pp.Set(k, p.m[k])
		***REMOVED***
	***REMOVED***
	return pp
***REMOVED***

// FilterStripPrefix returns a new properties object with a subset of all keys
// with the given prefix and the prefix removed from the keys.
func (p *Properties) FilterStripPrefix(prefix string) *Properties ***REMOVED***
	pp := NewProperties()
	n := len(prefix)
	for _, k := range p.k ***REMOVED***
		if len(k) > len(prefix) && strings.HasPrefix(k, prefix) ***REMOVED***
			// TODO(fs): we are ignoring the error which flags a circular reference.
			// TODO(fs): since we are modifying keys I am not entirely sure whether we can create a circular reference
			// TODO(fs): this function should probably return an error but the signature is fixed
			pp.Set(k[n:], p.m[k])
		***REMOVED***
	***REMOVED***
	return pp
***REMOVED***

// Len returns the number of keys.
func (p *Properties) Len() int ***REMOVED***
	return len(p.m)
***REMOVED***

// Keys returns all keys in the same order as in the input.
func (p *Properties) Keys() []string ***REMOVED***
	keys := make([]string, len(p.k))
	copy(keys, p.k)
	return keys
***REMOVED***

// Set sets the property key to the corresponding value.
// If a value for key existed before then ok is true and prev
// contains the previous value. If the value contains a
// circular reference or a malformed expression then
// an error is returned.
// An empty key is silently ignored.
func (p *Properties) Set(key, value string) (prev string, ok bool, err error) ***REMOVED***
	if key == "" ***REMOVED***
		return "", false, nil
	***REMOVED***

	// if expansion is disabled we allow circular references
	if p.DisableExpansion ***REMOVED***
		prev, ok = p.Get(key)
		p.m[key] = value
		if !ok ***REMOVED***
			p.k = append(p.k, key)
		***REMOVED***
		return prev, ok, nil
	***REMOVED***

	// to check for a circular reference we temporarily need
	// to set the new value. If there is an error then revert
	// to the previous state. Only if all tests are successful
	// then we add the key to the p.k list.
	prev, ok = p.Get(key)
	p.m[key] = value

	// now check for a circular reference
	_, err = p.expand(value)
	if err != nil ***REMOVED***

		// revert to the previous state
		if ok ***REMOVED***
			p.m[key] = prev
		***REMOVED*** else ***REMOVED***
			delete(p.m, key)
		***REMOVED***

		return "", false, err
	***REMOVED***

	if !ok ***REMOVED***
		p.k = append(p.k, key)
	***REMOVED***

	return prev, ok, nil
***REMOVED***

// SetValue sets property key to the default string value
// as defined by fmt.Sprintf("%v").
func (p *Properties) SetValue(key string, value interface***REMOVED******REMOVED***) error ***REMOVED***
	_, _, err := p.Set(key, fmt.Sprintf("%v", value))
	return err
***REMOVED***

// MustSet sets the property key to the corresponding value.
// If a value for key existed before then ok is true and prev
// contains the previous value. An empty key is silently ignored.
func (p *Properties) MustSet(key, value string) (prev string, ok bool) ***REMOVED***
	prev, ok, err := p.Set(key, value)
	if err != nil ***REMOVED***
		ErrorHandler(err)
	***REMOVED***
	return prev, ok
***REMOVED***

// String returns a string of all expanded 'key = value' pairs.
func (p *Properties) String() string ***REMOVED***
	var s string
	for _, key := range p.k ***REMOVED***
		value, _ := p.Get(key)
		s = fmt.Sprintf("%s%s = %s\n", s, key, value)
	***REMOVED***
	return s
***REMOVED***

// Write writes all unexpanded 'key = value' pairs to the given writer.
// Write returns the number of bytes written and any write error encountered.
func (p *Properties) Write(w io.Writer, enc Encoding) (n int, err error) ***REMOVED***
	return p.WriteComment(w, "", enc)
***REMOVED***

// WriteComment writes all unexpanced 'key = value' pairs to the given writer.
// If prefix is not empty then comments are written with a blank line and the
// given prefix. The prefix should be either "# " or "! " to be compatible with
// the properties file format. Otherwise, the properties parser will not be
// able to read the file back in. It returns the number of bytes written and
// any write error encountered.
func (p *Properties) WriteComment(w io.Writer, prefix string, enc Encoding) (n int, err error) ***REMOVED***
	var x int

	for _, key := range p.k ***REMOVED***
		value := p.m[key]

		if prefix != "" ***REMOVED***
			if comments, ok := p.c[key]; ok ***REMOVED***
				// don't print comments if they are all empty
				allEmpty := true
				for _, c := range comments ***REMOVED***
					if c != "" ***REMOVED***
						allEmpty = false
						break
					***REMOVED***
				***REMOVED***

				if !allEmpty ***REMOVED***
					// add a blank line between entries but not at the top
					if len(comments) > 0 && n > 0 ***REMOVED***
						x, err = fmt.Fprintln(w)
						if err != nil ***REMOVED***
							return
						***REMOVED***
						n += x
					***REMOVED***

					for _, c := range comments ***REMOVED***
						x, err = fmt.Fprintf(w, "%s%s\n", prefix, encode(c, "", enc))
						if err != nil ***REMOVED***
							return
						***REMOVED***
						n += x
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		x, err = fmt.Fprintf(w, "%s = %s\n", encode(key, " :", enc), encode(value, "", enc))
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n += x
	***REMOVED***
	return
***REMOVED***

// Map returns a copy of the properties as a map.
func (p *Properties) Map() map[string]string ***REMOVED***
	m := make(map[string]string)
	for k, v := range p.m ***REMOVED***
		m[k] = v
	***REMOVED***
	return m
***REMOVED***

// FilterFunc returns a copy of the properties which includes the values which passed all filters.
func (p *Properties) FilterFunc(filters ...func(k, v string) bool) *Properties ***REMOVED***
	pp := NewProperties()
outer:
	for k, v := range p.m ***REMOVED***
		for _, f := range filters ***REMOVED***
			if !f(k, v) ***REMOVED***
				continue outer
			***REMOVED***
			pp.Set(k, v)
		***REMOVED***
	***REMOVED***
	return pp
***REMOVED***

// ----------------------------------------------------------------------------

// Delete removes the key and its comments.
func (p *Properties) Delete(key string) ***REMOVED***
	delete(p.m, key)
	delete(p.c, key)
	newKeys := []string***REMOVED******REMOVED***
	for _, k := range p.k ***REMOVED***
		if k != key ***REMOVED***
			newKeys = append(newKeys, k)
		***REMOVED***
	***REMOVED***
	p.k = newKeys
***REMOVED***

// Merge merges properties, comments and keys from other *Properties into p
func (p *Properties) Merge(other *Properties) ***REMOVED***
	for k, v := range other.m ***REMOVED***
		p.m[k] = v
	***REMOVED***
	for k, v := range other.c ***REMOVED***
		p.c[k] = v
	***REMOVED***

outer:
	for _, otherKey := range other.k ***REMOVED***
		for _, key := range p.k ***REMOVED***
			if otherKey == key ***REMOVED***
				continue outer
			***REMOVED***
		***REMOVED***
		p.k = append(p.k, otherKey)
	***REMOVED***
***REMOVED***

// ----------------------------------------------------------------------------

// check expands all values and returns an error if a circular reference or
// a malformed expression was found.
func (p *Properties) check() error ***REMOVED***
	for _, value := range p.m ***REMOVED***
		if _, err := p.expand(value); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (p *Properties) expand(input string) (string, error) ***REMOVED***
	// no pre/postfix -> nothing to expand
	if p.Prefix == "" && p.Postfix == "" ***REMOVED***
		return input, nil
	***REMOVED***

	return expand(input, make(map[string]bool), p.Prefix, p.Postfix, p.m)
***REMOVED***

// expand recursively expands expressions of '(prefix)key(postfix)' to their corresponding values.
// The function keeps track of the keys that were already expanded and stops if it
// detects a circular reference or a malformed expression of the form '(prefix)key'.
func expand(s string, keys map[string]bool, prefix, postfix string, values map[string]string) (string, error) ***REMOVED***
	start := strings.Index(s, prefix)
	if start == -1 ***REMOVED***
		return s, nil
	***REMOVED***

	keyStart := start + len(prefix)
	keyLen := strings.Index(s[keyStart:], postfix)
	if keyLen == -1 ***REMOVED***
		return "", fmt.Errorf("malformed expression")
	***REMOVED***

	end := keyStart + keyLen + len(postfix) - 1
	key := s[keyStart : keyStart+keyLen]

	// fmt.Printf("s:%q pp:%q start:%d end:%d keyStart:%d keyLen:%d key:%q\n", s, prefix + "..." + postfix, start, end, keyStart, keyLen, key)

	if _, ok := keys[key]; ok ***REMOVED***
		return "", fmt.Errorf("circular reference")
	***REMOVED***

	val, ok := values[key]
	if !ok ***REMOVED***
		val = os.Getenv(key)
	***REMOVED***

	// remember that we've seen the key
	keys[key] = true

	return expand(s[:start]+val+s[end+1:], keys, prefix, postfix, values)
***REMOVED***

// encode encodes a UTF-8 string to ISO-8859-1 and escapes some characters.
func encode(s string, special string, enc Encoding) string ***REMOVED***
	switch enc ***REMOVED***
	case UTF8:
		return encodeUtf8(s, special)
	case ISO_8859_1:
		return encodeIso(s, special)
	default:
		panic(fmt.Sprintf("unsupported encoding %v", enc))
	***REMOVED***
***REMOVED***

func encodeUtf8(s string, special string) string ***REMOVED***
	v := ""
	for pos := 0; pos < len(s); ***REMOVED***
		r, w := utf8.DecodeRuneInString(s[pos:])
		pos += w
		v += escape(r, special)
	***REMOVED***
	return v
***REMOVED***

func encodeIso(s string, special string) string ***REMOVED***
	var r rune
	var w int
	var v string
	for pos := 0; pos < len(s); ***REMOVED***
		switch r, w = utf8.DecodeRuneInString(s[pos:]); ***REMOVED***
		case r < 1<<8: // single byte rune -> escape special chars only
			v += escape(r, special)
		case r < 1<<16: // two byte rune -> unicode literal
			v += fmt.Sprintf("\\u%04x", r)
		default: // more than two bytes per rune -> can't encode
			v += "?"
		***REMOVED***
		pos += w
	***REMOVED***
	return v
***REMOVED***

func escape(r rune, special string) string ***REMOVED***
	switch r ***REMOVED***
	case '\f':
		return "\\f"
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	case '\t':
		return "\\t"
	default:
		if strings.ContainsRune(special, r) ***REMOVED***
			return "\\" + string(r)
		***REMOVED***
		return string(r)
	***REMOVED***
***REMOVED***

func invalidKeyError(key string) error ***REMOVED***
	return fmt.Errorf("unknown property: %s", key)
***REMOVED***
