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

// Package version implements etcd version parsing and contains latest version
// information.
package version

import (
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

var (
	// MinClusterVersion is the min cluster version this etcd binary is compatible with.
	MinClusterVersion = "3.0.0"
	Version           = "3.2.1"
	APIVersion        = "unknown"

	// Git SHA Value will be set during build
	GitSHA = "Not provided (use ./build instead of go build)"
)

func init() ***REMOVED***
	ver, err := semver.NewVersion(Version)
	if err == nil ***REMOVED***
		APIVersion = fmt.Sprintf("%d.%d", ver.Major, ver.Minor)
	***REMOVED***
***REMOVED***

type Versions struct ***REMOVED***
	Server  string `json:"etcdserver"`
	Cluster string `json:"etcdcluster"`
	// TODO: raft state machine version
***REMOVED***

// Cluster only keeps the major.minor.
func Cluster(v string) string ***REMOVED***
	vs := strings.Split(v, ".")
	if len(vs) <= 2 ***REMOVED***
		return v
	***REMOVED***
	return fmt.Sprintf("%s.%s", vs[0], vs[1])
***REMOVED***
