// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"errors"
	"google.golang.org/grpc/naming"
)

// PoolResolver provides a fixed list of addresses to load balance between
// and does not provide further updates.
type PoolResolver struct ***REMOVED***
	poolSize int
	dialOpt  *DialSettings
	ch       chan []*naming.Update
***REMOVED***

// NewPoolResolver returns a PoolResolver
// This is an EXPERIMENTAL API and may be changed or removed in the future.
func NewPoolResolver(size int, o *DialSettings) *PoolResolver ***REMOVED***
	return &PoolResolver***REMOVED***poolSize: size, dialOpt: o***REMOVED***
***REMOVED***

// Resolve returns a Watcher for the endpoint defined by the DialSettings
// provided to NewPoolResolver.
func (r *PoolResolver) Resolve(target string) (naming.Watcher, error) ***REMOVED***
	if r.dialOpt.Endpoint == "" ***REMOVED***
		return nil, errors.New("No endpoint configured")
	***REMOVED***
	addrs := make([]*naming.Update, 0, r.poolSize)
	for i := 0; i < r.poolSize; i++ ***REMOVED***
		addrs = append(addrs, &naming.Update***REMOVED***Op: naming.Add, Addr: r.dialOpt.Endpoint, Metadata: i***REMOVED***)
	***REMOVED***
	r.ch = make(chan []*naming.Update, 1)
	r.ch <- addrs
	return r, nil
***REMOVED***

// Next returns a static list of updates on the first call,
// and blocks indefinitely until Close is called on subsequent calls.
func (r *PoolResolver) Next() ([]*naming.Update, error) ***REMOVED***
	return <-r.ch, nil
***REMOVED***

func (r *PoolResolver) Close() ***REMOVED***
	close(r.ch)
***REMOVED***
