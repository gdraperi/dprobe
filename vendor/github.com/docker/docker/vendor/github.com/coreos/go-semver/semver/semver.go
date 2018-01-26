// Copyright 2013-2015 CoreOS, Inc.
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

// Semantic Versions http://semver.org
package semver

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct ***REMOVED***
	Major      int64
	Minor      int64
	Patch      int64
	PreRelease PreRelease
	Metadata   string
***REMOVED***

type PreRelease string

func splitOff(input *string, delim string) (val string) ***REMOVED***
	parts := strings.SplitN(*input, delim, 2)

	if len(parts) == 2 ***REMOVED***
		*input = parts[0]
		val = parts[1]
	***REMOVED***

	return val
***REMOVED***

func New(version string) *Version ***REMOVED***
	return Must(NewVersion(version))
***REMOVED***

func NewVersion(version string) (*Version, error) ***REMOVED***
	v := Version***REMOVED******REMOVED***

	if err := v.Set(version); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &v, nil
***REMOVED***

// Must is a helper for wrapping NewVersion and will panic if err is not nil.
func Must(v *Version, err error) *Version ***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return v
***REMOVED***

// Set parses and updates v from the given version string. Implements flag.Value
func (v *Version) Set(version string) error ***REMOVED***
	metadata := splitOff(&version, "+")
	preRelease := PreRelease(splitOff(&version, "-"))
	dotParts := strings.SplitN(version, ".", 3)

	if len(dotParts) != 3 ***REMOVED***
		return fmt.Errorf("%s is not in dotted-tri format", version)
	***REMOVED***

	parsed := make([]int64, 3, 3)

	for i, v := range dotParts[:3] ***REMOVED***
		val, err := strconv.ParseInt(v, 10, 64)
		parsed[i] = val
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	v.Metadata = metadata
	v.PreRelease = preRelease
	v.Major = parsed[0]
	v.Minor = parsed[1]
	v.Patch = parsed[2]
	return nil
***REMOVED***

func (v Version) String() string ***REMOVED***
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.PreRelease != "" ***REMOVED***
		fmt.Fprintf(&buffer, "-%s", v.PreRelease)
	***REMOVED***

	if v.Metadata != "" ***REMOVED***
		fmt.Fprintf(&buffer, "+%s", v.Metadata)
	***REMOVED***

	return buffer.String()
***REMOVED***

func (v *Version) UnmarshalYAML(unmarshal func(interface***REMOVED******REMOVED***) error) error ***REMOVED***
	var data string
	if err := unmarshal(&data); err != nil ***REMOVED***
		return err
	***REMOVED***
	return v.Set(data)
***REMOVED***

func (v Version) MarshalJSON() ([]byte, error) ***REMOVED***
	return []byte(`"` + v.String() + `"`), nil
***REMOVED***

func (v *Version) UnmarshalJSON(data []byte) error ***REMOVED***
	l := len(data)
	if l == 0 || string(data) == `""` ***REMOVED***
		return nil
	***REMOVED***
	if l < 2 || data[0] != '"' || data[l-1] != '"' ***REMOVED***
		return errors.New("invalid semver string")
	***REMOVED***
	return v.Set(string(data[1 : l-1]))
***REMOVED***

// Compare tests if v is less than, equal to, or greater than versionB,
// returning -1, 0, or +1 respectively.
func (v Version) Compare(versionB Version) int ***REMOVED***
	if cmp := recursiveCompare(v.Slice(), versionB.Slice()); cmp != 0 ***REMOVED***
		return cmp
	***REMOVED***
	return preReleaseCompare(v, versionB)
***REMOVED***

// Equal tests if v is equal to versionB.
func (v Version) Equal(versionB Version) bool ***REMOVED***
	return v.Compare(versionB) == 0
***REMOVED***

// LessThan tests if v is less than versionB.
func (v Version) LessThan(versionB Version) bool ***REMOVED***
	return v.Compare(versionB) < 0
***REMOVED***

// Slice converts the comparable parts of the semver into a slice of integers.
func (v Version) Slice() []int64 ***REMOVED***
	return []int64***REMOVED***v.Major, v.Minor, v.Patch***REMOVED***
***REMOVED***

func (p PreRelease) Slice() []string ***REMOVED***
	preRelease := string(p)
	return strings.Split(preRelease, ".")
***REMOVED***

func preReleaseCompare(versionA Version, versionB Version) int ***REMOVED***
	a := versionA.PreRelease
	b := versionB.PreRelease

	/* Handle the case where if two versions are otherwise equal it is the
	 * one without a PreRelease that is greater */
	if len(a) == 0 && (len(b) > 0) ***REMOVED***
		return 1
	***REMOVED*** else if len(b) == 0 && (len(a) > 0) ***REMOVED***
		return -1
	***REMOVED***

	// If there is a prerelease, check and compare each part.
	return recursivePreReleaseCompare(a.Slice(), b.Slice())
***REMOVED***

func recursiveCompare(versionA []int64, versionB []int64) int ***REMOVED***
	if len(versionA) == 0 ***REMOVED***
		return 0
	***REMOVED***

	a := versionA[0]
	b := versionB[0]

	if a > b ***REMOVED***
		return 1
	***REMOVED*** else if a < b ***REMOVED***
		return -1
	***REMOVED***

	return recursiveCompare(versionA[1:], versionB[1:])
***REMOVED***

func recursivePreReleaseCompare(versionA []string, versionB []string) int ***REMOVED***
	// A larger set of pre-release fields has a higher precedence than a smaller set,
	// if all of the preceding identifiers are equal.
	if len(versionA) == 0 ***REMOVED***
		if len(versionB) > 0 ***REMOVED***
			return -1
		***REMOVED***
		return 0
	***REMOVED*** else if len(versionB) == 0 ***REMOVED***
		// We're longer than versionB so return 1.
		return 1
	***REMOVED***

	a := versionA[0]
	b := versionB[0]

	aInt := false
	bInt := false

	aI, err := strconv.Atoi(versionA[0])
	if err == nil ***REMOVED***
		aInt = true
	***REMOVED***

	bI, err := strconv.Atoi(versionB[0])
	if err == nil ***REMOVED***
		bInt = true
	***REMOVED***

	// Handle Integer Comparison
	if aInt && bInt ***REMOVED***
		if aI > bI ***REMOVED***
			return 1
		***REMOVED*** else if aI < bI ***REMOVED***
			return -1
		***REMOVED***
	***REMOVED***

	// Handle String Comparison
	if a > b ***REMOVED***
		return 1
	***REMOVED*** else if a < b ***REMOVED***
		return -1
	***REMOVED***

	return recursivePreReleaseCompare(versionA[1:], versionB[1:])
***REMOVED***

// BumpMajor increments the Major field by 1 and resets all other fields to their default values
func (v *Version) BumpMajor() ***REMOVED***
	v.Major += 1
	v.Minor = 0
	v.Patch = 0
	v.PreRelease = PreRelease("")
	v.Metadata = ""
***REMOVED***

// BumpMinor increments the Minor field by 1 and resets all other fields to their default values
func (v *Version) BumpMinor() ***REMOVED***
	v.Minor += 1
	v.Patch = 0
	v.PreRelease = PreRelease("")
	v.Metadata = ""
***REMOVED***

// BumpPatch increments the Patch field by 1 and resets all other fields to their default values
func (v *Version) BumpPatch() ***REMOVED***
	v.Patch += 1
	v.PreRelease = PreRelease("")
	v.Metadata = ""
***REMOVED***
