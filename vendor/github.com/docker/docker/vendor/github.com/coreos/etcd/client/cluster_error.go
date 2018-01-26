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

package client

import "fmt"

type ClusterError struct ***REMOVED***
	Errors []error
***REMOVED***

func (ce *ClusterError) Error() string ***REMOVED***
	s := ErrClusterUnavailable.Error()
	for i, e := range ce.Errors ***REMOVED***
		s += fmt.Sprintf("; error #%d: %s\n", i, e)
	***REMOVED***
	return s
***REMOVED***

func (ce *ClusterError) Detail() string ***REMOVED***
	s := ""
	for i, e := range ce.Errors ***REMOVED***
		s += fmt.Sprintf("error #%d: %s\n", i, e)
	***REMOVED***
	return s
***REMOVED***
