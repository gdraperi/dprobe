// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Viper is a application configuration system.
// It believes that applications can be configured a variety of ways
// via flags, ENVIRONMENT variables, configuration files retrieved
// from the file system, or a remote key/value store.

package viper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/spf13/afero"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

// ConfigParseError denotes failing to parse configuration file.
type ConfigParseError struct ***REMOVED***
	err error
***REMOVED***

// Error returns the formatted configuration error.
func (pe ConfigParseError) Error() string ***REMOVED***
	return fmt.Sprintf("While parsing config: %s", pe.err.Error())
***REMOVED***

// toCaseInsensitiveValue checks if the value is a  map;
// if so, create a copy and lower-case the keys recursively.
func toCaseInsensitiveValue(value interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	switch v := value.(type) ***REMOVED***
	case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		value = copyAndInsensitiviseMap(cast.ToStringMap(v))
	case map[string]interface***REMOVED******REMOVED***:
		value = copyAndInsensitiviseMap(v)
	***REMOVED***

	return value
***REMOVED***

// copyAndInsensitiviseMap behaves like insensitiviseMap, but creates a copy of
// any map it makes case insensitive.
func copyAndInsensitiviseMap(m map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	nm := make(map[string]interface***REMOVED******REMOVED***)

	for key, val := range m ***REMOVED***
		lkey := strings.ToLower(key)
		switch v := val.(type) ***REMOVED***
		case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
			nm[lkey] = copyAndInsensitiviseMap(cast.ToStringMap(v))
		case map[string]interface***REMOVED******REMOVED***:
			nm[lkey] = copyAndInsensitiviseMap(v)
		default:
			nm[lkey] = v
		***REMOVED***
	***REMOVED***

	return nm
***REMOVED***

func insensitiviseMap(m map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	for key, val := range m ***REMOVED***
		switch val.(type) ***REMOVED***
		case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
			// nested map: cast and recursively insensitivise
			val = cast.ToStringMap(val)
			insensitiviseMap(val.(map[string]interface***REMOVED******REMOVED***))
		case map[string]interface***REMOVED******REMOVED***:
			// nested map: recursively insensitivise
			insensitiviseMap(val.(map[string]interface***REMOVED******REMOVED***))
		***REMOVED***

		lower := strings.ToLower(key)
		if key != lower ***REMOVED***
			// remove old key (not lower-cased)
			delete(m, key)
		***REMOVED***
		// update map
		m[lower] = val
	***REMOVED***
***REMOVED***

func absPathify(inPath string) string ***REMOVED***
	jww.INFO.Println("Trying to resolve absolute path to", inPath)

	if strings.HasPrefix(inPath, "$HOME") ***REMOVED***
		inPath = userHomeDir() + inPath[5:]
	***REMOVED***

	if strings.HasPrefix(inPath, "$") ***REMOVED***
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	***REMOVED***

	if filepath.IsAbs(inPath) ***REMOVED***
		return filepath.Clean(inPath)
	***REMOVED***

	p, err := filepath.Abs(inPath)
	if err == nil ***REMOVED***
		return filepath.Clean(p)
	***REMOVED***

	jww.ERROR.Println("Couldn't discover absolute path")
	jww.ERROR.Println(err)
	return ""
***REMOVED***

// Check if File / Directory Exists
func exists(fs afero.Fs, path string) (bool, error) ***REMOVED***
	_, err := fs.Stat(path)
	if err == nil ***REMOVED***
		return true, nil
	***REMOVED***
	if os.IsNotExist(err) ***REMOVED***
		return false, nil
	***REMOVED***
	return false, err
***REMOVED***

func stringInSlice(a string, list []string) bool ***REMOVED***
	for _, b := range list ***REMOVED***
		if b == a ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func userHomeDir() string ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" ***REMOVED***
			home = os.Getenv("USERPROFILE")
		***REMOVED***
		return home
	***REMOVED***
	return os.Getenv("HOME")
***REMOVED***

func safeMul(a, b uint) uint ***REMOVED***
	c := a * b
	if a > 1 && b > 1 && c/b != a ***REMOVED***
		return 0
	***REMOVED***
	return c
***REMOVED***

// parseSizeInBytes converts strings like 1GB or 12 mb into an unsigned integer number of bytes
func parseSizeInBytes(sizeStr string) uint ***REMOVED***
	sizeStr = strings.TrimSpace(sizeStr)
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 ***REMOVED***
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' ***REMOVED***
			if lastChar > 1 ***REMOVED***
				switch unicode.ToLower(rune(sizeStr[lastChar-1])) ***REMOVED***
				case 'k':
					multiplier = 1 << 10
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << 20
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << 30
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	size := cast.ToInt(sizeStr)
	if size < 0 ***REMOVED***
		size = 0
	***REMOVED***

	return safeMul(uint(size), multiplier)
***REMOVED***

// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
func deepSearch(m map[string]interface***REMOVED******REMOVED***, path []string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	for _, k := range path ***REMOVED***
		m2, ok := m[k]
		if !ok ***REMOVED***
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface***REMOVED******REMOVED***)
			m[k] = m3
			m = m3
			continue
		***REMOVED***
		m3, ok := m2.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface***REMOVED******REMOVED***)
			m[k] = m3
		***REMOVED***
		// continue search from here
		m = m3
	***REMOVED***
	return m
***REMOVED***
