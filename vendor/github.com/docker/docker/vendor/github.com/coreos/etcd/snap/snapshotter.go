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

// Package snap stores raft nodes' states with snapshots.
package snap

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	pioutil "github.com/coreos/etcd/pkg/ioutil"
	"github.com/coreos/etcd/pkg/pbutil"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/snap/snappb"

	"github.com/coreos/pkg/capnslog"
)

const (
	snapSuffix = ".snap"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/etcd", "snap")

	ErrNoSnapshot    = errors.New("snap: no available snapshot")
	ErrEmptySnapshot = errors.New("snap: empty snapshot")
	ErrCRCMismatch   = errors.New("snap: crc mismatch")
	crcTable         = crc32.MakeTable(crc32.Castagnoli)

	// A map of valid files that can be present in the snap folder.
	validFiles = map[string]bool***REMOVED***
		"db": true,
	***REMOVED***
)

type Snapshotter struct ***REMOVED***
	dir string
***REMOVED***

func New(dir string) *Snapshotter ***REMOVED***
	return &Snapshotter***REMOVED***
		dir: dir,
	***REMOVED***
***REMOVED***

func (s *Snapshotter) SaveSnap(snapshot raftpb.Snapshot) error ***REMOVED***
	if raft.IsEmptySnap(snapshot) ***REMOVED***
		return nil
	***REMOVED***
	return s.save(&snapshot)
***REMOVED***

func (s *Snapshotter) save(snapshot *raftpb.Snapshot) error ***REMOVED***
	start := time.Now()

	fname := fmt.Sprintf("%016x-%016x%s", snapshot.Metadata.Term, snapshot.Metadata.Index, snapSuffix)
	b := pbutil.MustMarshal(snapshot)
	crc := crc32.Update(0, crcTable, b)
	snap := snappb.Snapshot***REMOVED***Crc: crc, Data: b***REMOVED***
	d, err := snap.Marshal()
	if err != nil ***REMOVED***
		return err
	***REMOVED*** else ***REMOVED***
		marshallingDurations.Observe(float64(time.Since(start)) / float64(time.Second))
	***REMOVED***

	err = pioutil.WriteAndSyncFile(filepath.Join(s.dir, fname), d, 0666)
	if err == nil ***REMOVED***
		saveDurations.Observe(float64(time.Since(start)) / float64(time.Second))
	***REMOVED*** else ***REMOVED***
		err1 := os.Remove(filepath.Join(s.dir, fname))
		if err1 != nil ***REMOVED***
			plog.Errorf("failed to remove broken snapshot file %s", filepath.Join(s.dir, fname))
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (s *Snapshotter) Load() (*raftpb.Snapshot, error) ***REMOVED***
	names, err := s.snapNames()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var snap *raftpb.Snapshot
	for _, name := range names ***REMOVED***
		if snap, err = loadSnap(s.dir, name); err == nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, ErrNoSnapshot
	***REMOVED***
	return snap, nil
***REMOVED***

func loadSnap(dir, name string) (*raftpb.Snapshot, error) ***REMOVED***
	fpath := filepath.Join(dir, name)
	snap, err := Read(fpath)
	if err != nil ***REMOVED***
		renameBroken(fpath)
	***REMOVED***
	return snap, err
***REMOVED***

// Read reads the snapshot named by snapname and returns the snapshot.
func Read(snapname string) (*raftpb.Snapshot, error) ***REMOVED***
	b, err := ioutil.ReadFile(snapname)
	if err != nil ***REMOVED***
		plog.Errorf("cannot read file %v: %v", snapname, err)
		return nil, err
	***REMOVED***

	if len(b) == 0 ***REMOVED***
		plog.Errorf("unexpected empty snapshot")
		return nil, ErrEmptySnapshot
	***REMOVED***

	var serializedSnap snappb.Snapshot
	if err = serializedSnap.Unmarshal(b); err != nil ***REMOVED***
		plog.Errorf("corrupted snapshot file %v: %v", snapname, err)
		return nil, err
	***REMOVED***

	if len(serializedSnap.Data) == 0 || serializedSnap.Crc == 0 ***REMOVED***
		plog.Errorf("unexpected empty snapshot")
		return nil, ErrEmptySnapshot
	***REMOVED***

	crc := crc32.Update(0, crcTable, serializedSnap.Data)
	if crc != serializedSnap.Crc ***REMOVED***
		plog.Errorf("corrupted snapshot file %v: crc mismatch", snapname)
		return nil, ErrCRCMismatch
	***REMOVED***

	var snap raftpb.Snapshot
	if err = snap.Unmarshal(serializedSnap.Data); err != nil ***REMOVED***
		plog.Errorf("corrupted snapshot file %v: %v", snapname, err)
		return nil, err
	***REMOVED***
	return &snap, nil
***REMOVED***

// snapNames returns the filename of the snapshots in logical time order (from newest to oldest).
// If there is no available snapshots, an ErrNoSnapshot will be returned.
func (s *Snapshotter) snapNames() ([]string, error) ***REMOVED***
	dir, err := os.Open(s.dir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	snaps := checkSuffix(names)
	if len(snaps) == 0 ***REMOVED***
		return nil, ErrNoSnapshot
	***REMOVED***
	sort.Sort(sort.Reverse(sort.StringSlice(snaps)))
	return snaps, nil
***REMOVED***

func checkSuffix(names []string) []string ***REMOVED***
	snaps := []string***REMOVED******REMOVED***
	for i := range names ***REMOVED***
		if strings.HasSuffix(names[i], snapSuffix) ***REMOVED***
			snaps = append(snaps, names[i])
		***REMOVED*** else ***REMOVED***
			// If we find a file which is not a snapshot then check if it's
			// a vaild file. If not throw out a warning.
			if _, ok := validFiles[names[i]]; !ok ***REMOVED***
				plog.Warningf("skipped unexpected non snapshot file %v", names[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return snaps
***REMOVED***

func renameBroken(path string) ***REMOVED***
	brokenPath := path + ".broken"
	if err := os.Rename(path, brokenPath); err != nil ***REMOVED***
		plog.Warningf("cannot rename broken snapshot file %v to %v: %v", path, brokenPath, err)
	***REMOVED***
***REMOVED***
