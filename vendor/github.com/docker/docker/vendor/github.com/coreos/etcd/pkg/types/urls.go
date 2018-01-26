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
	"errors"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
)

type URLs []url.URL

func NewURLs(strs []string) (URLs, error) ***REMOVED***
	all := make([]url.URL, len(strs))
	if len(all) == 0 ***REMOVED***
		return nil, errors.New("no valid URLs given")
	***REMOVED***
	for i, in := range strs ***REMOVED***
		in = strings.TrimSpace(in)
		u, err := url.Parse(in)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "unix" && u.Scheme != "unixs" ***REMOVED***
			return nil, fmt.Errorf("URL scheme must be http, https, unix, or unixs: %s", in)
		***REMOVED***
		if _, _, err := net.SplitHostPort(u.Host); err != nil ***REMOVED***
			return nil, fmt.Errorf(`URL address does not have the form "host:port": %s`, in)
		***REMOVED***
		if u.Path != "" ***REMOVED***
			return nil, fmt.Errorf("URL must not contain a path: %s", in)
		***REMOVED***
		all[i] = *u
	***REMOVED***
	us := URLs(all)
	us.Sort()

	return us, nil
***REMOVED***

func MustNewURLs(strs []string) URLs ***REMOVED***
	urls, err := NewURLs(strs)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return urls
***REMOVED***

func (us URLs) String() string ***REMOVED***
	return strings.Join(us.StringSlice(), ",")
***REMOVED***

func (us *URLs) Sort() ***REMOVED***
	sort.Sort(us)
***REMOVED***
func (us URLs) Len() int           ***REMOVED*** return len(us) ***REMOVED***
func (us URLs) Less(i, j int) bool ***REMOVED*** return us[i].String() < us[j].String() ***REMOVED***
func (us URLs) Swap(i, j int)      ***REMOVED*** us[i], us[j] = us[j], us[i] ***REMOVED***

func (us URLs) StringSlice() []string ***REMOVED***
	out := make([]string, len(us))
	for i := range us ***REMOVED***
		out[i] = us[i].String()
	***REMOVED***

	return out
***REMOVED***
