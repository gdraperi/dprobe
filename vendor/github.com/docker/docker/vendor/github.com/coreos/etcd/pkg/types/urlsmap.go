// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"sort"
	"strings"
)

// URLsMap is a map from a name to its URLs.
type URLsMap map[string]URLs

// NewURLsMap returns a URLsMap instantiated from the given string,
// which consists of discovery-formatted names-to-URLs, like:
// mach0=http://1.1.1.1:2380,mach0=http://2.2.2.2::2380,mach1=http://3.3.3.3:2380,mach2=http://4.4.4.4:2380
func NewURLsMap(s string) (URLsMap, error) ***REMOVED***
	m := parse(s)

	cl := URLsMap***REMOVED******REMOVED***
	for name, urls := range m ***REMOVED***
		us, err := NewURLs(urls)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cl[name] = us
	***REMOVED***
	return cl, nil
***REMOVED***

// NewURLsMapFromStringMap takes a map of strings and returns a URLsMap. The
// string values in the map can be multiple values separated by the sep string.
func NewURLsMapFromStringMap(m map[string]string, sep string) (URLsMap, error) ***REMOVED***
	var err error
	um := URLsMap***REMOVED******REMOVED***
	for k, v := range m ***REMOVED***
		um[k], err = NewURLs(strings.Split(v, sep))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return um, nil
***REMOVED***

// String turns URLsMap into discovery-formatted name-to-URLs sorted by name.
func (c URLsMap) String() string ***REMOVED***
	var pairs []string
	for name, urls := range c ***REMOVED***
		for _, url := range urls ***REMOVED***
			pairs = append(pairs, fmt.Sprintf("%s=%s", name, url.String()))
		***REMOVED***
	***REMOVED***
	sort.Strings(pairs)
	return strings.Join(pairs, ",")
***REMOVED***

// URLs returns a list of all URLs.
// The returned list is sorted in ascending lexicographical order.
func (c URLsMap) URLs() []string ***REMOVED***
	var urls []string
	for _, us := range c ***REMOVED***
		for _, u := range us ***REMOVED***
			urls = append(urls, u.String())
		***REMOVED***
	***REMOVED***
	sort.Strings(urls)
	return urls
***REMOVED***

// Len returns the size of URLsMap.
func (c URLsMap) Len() int ***REMOVED***
	return len(c)
***REMOVED***

// parse parses the given string and returns a map listing the values specified for each key.
func parse(s string) map[string][]string ***REMOVED***
	m := make(map[string][]string)
	for s != "" ***REMOVED***
		key := s
		if i := strings.IndexAny(key, ","); i >= 0 ***REMOVED***
			key, s = key[:i], key[i+1:]
		***REMOVED*** else ***REMOVED***
			s = ""
		***REMOVED***
		if key == "" ***REMOVED***
			continue
		***REMOVED***
		value := ""
		if i := strings.Index(key, "="); i >= 0 ***REMOVED***
			key, value = key[:i], key[i+1:]
		***REMOVED***
		m[key] = append(m[key], value)
	***REMOVED***
	return m
***REMOVED***
