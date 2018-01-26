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

// Package pbutil defines interfaces for handling Protocol Buffer objects.
package pbutil

import "github.com/coreos/pkg/capnslog"

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/etcd", "pkg/pbutil")
)

type Marshaler interface ***REMOVED***
	Marshal() (data []byte, err error)
***REMOVED***

type Unmarshaler interface ***REMOVED***
	Unmarshal(data []byte) error
***REMOVED***

func MustMarshal(m Marshaler) []byte ***REMOVED***
	d, err := m.Marshal()
	if err != nil ***REMOVED***
		plog.Panicf("marshal should never fail (%v)", err)
	***REMOVED***
	return d
***REMOVED***

func MustUnmarshal(um Unmarshaler, data []byte) ***REMOVED***
	if err := um.Unmarshal(data); err != nil ***REMOVED***
		plog.Panicf("unmarshal should never fail (%v)", err)
	***REMOVED***
***REMOVED***

func MaybeUnmarshal(um Unmarshaler, data []byte) bool ***REMOVED***
	if err := um.Unmarshal(data); err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func GetBool(v *bool) (vv bool, set bool) ***REMOVED***
	if v == nil ***REMOVED***
		return false, false
	***REMOVED***
	return *v, true
***REMOVED***

func Boolp(b bool) *bool ***REMOVED*** return &b ***REMOVED***
