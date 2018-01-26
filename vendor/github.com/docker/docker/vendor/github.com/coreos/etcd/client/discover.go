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

import (
	"github.com/coreos/etcd/pkg/srv"
)

// Discoverer is an interface that wraps the Discover method.
type Discoverer interface ***REMOVED***
	// Discover looks up the etcd servers for the domain.
	Discover(domain string) ([]string, error)
***REMOVED***

type srvDiscover struct***REMOVED******REMOVED***

// NewSRVDiscover constructs a new Discoverer that uses the stdlib to lookup SRV records.
func NewSRVDiscover() Discoverer ***REMOVED***
	return &srvDiscover***REMOVED******REMOVED***
***REMOVED***

func (d *srvDiscover) Discover(domain string) ([]string, error) ***REMOVED***
	srvs, err := srv.GetClient("etcd-client", domain)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return srvs.Endpoints, nil
***REMOVED***
