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

package wal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/etcd/pkg/fileutil"
)

var (
	badWalName = errors.New("bad wal name")
)

func Exist(dirpath string) bool ***REMOVED***
	names, err := fileutil.ReadDir(dirpath)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return len(names) != 0
***REMOVED***

// searchIndex returns the last array index of names whose raft index section is
// equal to or smaller than the given index.
// The given names MUST be sorted.
func searchIndex(names []string, index uint64) (int, bool) ***REMOVED***
	for i := len(names) - 1; i >= 0; i-- ***REMOVED***
		name := names[i]
		_, curIndex, err := parseWalName(name)
		if err != nil ***REMOVED***
			plog.Panicf("parse correct name should never fail: %v", err)
		***REMOVED***
		if index >= curIndex ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***
	return -1, false
***REMOVED***

// names should have been sorted based on sequence number.
// isValidSeq checks whether seq increases continuously.
func isValidSeq(names []string) bool ***REMOVED***
	var lastSeq uint64
	for _, name := range names ***REMOVED***
		curSeq, _, err := parseWalName(name)
		if err != nil ***REMOVED***
			plog.Panicf("parse correct name should never fail: %v", err)
		***REMOVED***
		if lastSeq != 0 && lastSeq != curSeq-1 ***REMOVED***
			return false
		***REMOVED***
		lastSeq = curSeq
	***REMOVED***
	return true
***REMOVED***
func readWalNames(dirpath string) ([]string, error) ***REMOVED***
	names, err := fileutil.ReadDir(dirpath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	wnames := checkWalNames(names)
	if len(wnames) == 0 ***REMOVED***
		return nil, ErrFileNotFound
	***REMOVED***
	return wnames, nil
***REMOVED***

func checkWalNames(names []string) []string ***REMOVED***
	wnames := make([]string, 0)
	for _, name := range names ***REMOVED***
		if _, _, err := parseWalName(name); err != nil ***REMOVED***
			// don't complain about left over tmp files
			if !strings.HasSuffix(name, ".tmp") ***REMOVED***
				plog.Warningf("ignored file %v in wal", name)
			***REMOVED***
			continue
		***REMOVED***
		wnames = append(wnames, name)
	***REMOVED***
	return wnames
***REMOVED***

func parseWalName(str string) (seq, index uint64, err error) ***REMOVED***
	if !strings.HasSuffix(str, ".wal") ***REMOVED***
		return 0, 0, badWalName
	***REMOVED***
	_, err = fmt.Sscanf(str, "%016x-%016x.wal", &seq, &index)
	return seq, index, err
***REMOVED***

func walName(seq, index uint64) string ***REMOVED***
	return fmt.Sprintf("%016x-%016x.wal", seq, index)
***REMOVED***
