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
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/coreos/etcd/pkg/pbutil"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/wal/walpb"

	"github.com/coreos/pkg/capnslog"
)

const (
	metadataType int64 = iota + 1
	entryType
	stateType
	crcType
	snapshotType

	// warnSyncDuration is the amount of time allotted to an fsync before
	// logging a warning
	warnSyncDuration = time.Second
)

var (
	// SegmentSizeBytes is the preallocated size of each wal segment file.
	// The actual size might be larger than this. In general, the default
	// value should be used, but this is defined as an exported variable
	// so that tests can set a different segment size.
	SegmentSizeBytes int64 = 64 * 1000 * 1000 // 64MB

	plog = capnslog.NewPackageLogger("github.com/coreos/etcd", "wal")

	ErrMetadataConflict = errors.New("wal: conflicting metadata found")
	ErrFileNotFound     = errors.New("wal: file not found")
	ErrCRCMismatch      = errors.New("wal: crc mismatch")
	ErrSnapshotMismatch = errors.New("wal: snapshot mismatch")
	ErrSnapshotNotFound = errors.New("wal: snapshot not found")
	crcTable            = crc32.MakeTable(crc32.Castagnoli)
)

// WAL is a logical representation of the stable storage.
// WAL is either in read mode or append mode but not both.
// A newly created WAL is in append mode, and ready for appending records.
// A just opened WAL is in read mode, and ready for reading records.
// The WAL will be ready for appending after reading out all the previous records.
type WAL struct ***REMOVED***
	dir string // the living directory of the underlay files

	// dirFile is a fd for the wal directory for syncing on Rename
	dirFile *os.File

	metadata []byte           // metadata recorded at the head of each WAL
	state    raftpb.HardState // hardstate recorded at the head of WAL

	start     walpb.Snapshot // snapshot to start reading
	decoder   *decoder       // decoder to decode records
	readClose func() error   // closer for decode reader

	mu      sync.Mutex
	enti    uint64   // index of the last entry saved to the wal
	encoder *encoder // encoder to encode records

	locks []*fileutil.LockedFile // the locked files the WAL holds (the name is increasing)
	fp    *filePipeline
***REMOVED***

// Create creates a WAL ready for appending records. The given metadata is
// recorded at the head of each WAL file, and can be retrieved with ReadAll.
func Create(dirpath string, metadata []byte) (*WAL, error) ***REMOVED***
	if Exist(dirpath) ***REMOVED***
		return nil, os.ErrExist
	***REMOVED***

	// keep temporary wal directory so WAL initialization appears atomic
	tmpdirpath := filepath.Clean(dirpath) + ".tmp"
	if fileutil.Exist(tmpdirpath) ***REMOVED***
		if err := os.RemoveAll(tmpdirpath); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if err := fileutil.CreateDirAll(tmpdirpath); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p := filepath.Join(tmpdirpath, walName(0, 0))
	f, err := fileutil.LockFile(p, os.O_WRONLY|os.O_CREATE, fileutil.PrivateFileMode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err = f.Seek(0, io.SeekEnd); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = fileutil.Preallocate(f.File, SegmentSizeBytes, true); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	w := &WAL***REMOVED***
		dir:      dirpath,
		metadata: metadata,
	***REMOVED***
	w.encoder, err = newFileEncoder(f.File, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	w.locks = append(w.locks, f)
	if err = w.saveCrc(0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = w.encoder.encode(&walpb.Record***REMOVED***Type: metadataType, Data: metadata***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = w.SaveSnapshot(walpb.Snapshot***REMOVED******REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if w, err = w.renameWal(tmpdirpath); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// directory was renamed; sync parent dir to persist rename
	pdir, perr := fileutil.OpenDir(filepath.Dir(w.dir))
	if perr != nil ***REMOVED***
		return nil, perr
	***REMOVED***
	if perr = fileutil.Fsync(pdir); perr != nil ***REMOVED***
		return nil, perr
	***REMOVED***
	if perr = pdir.Close(); err != nil ***REMOVED***
		return nil, perr
	***REMOVED***

	return w, nil
***REMOVED***

// Open opens the WAL at the given snap.
// The snap SHOULD have been previously saved to the WAL, or the following
// ReadAll will fail.
// The returned WAL is ready to read and the first record will be the one after
// the given snap. The WAL cannot be appended to before reading out all of its
// previous records.
func Open(dirpath string, snap walpb.Snapshot) (*WAL, error) ***REMOVED***
	w, err := openAtIndex(dirpath, snap, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if w.dirFile, err = fileutil.OpenDir(w.dir); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return w, nil
***REMOVED***

// OpenForRead only opens the wal files for read.
// Write on a read only wal panics.
func OpenForRead(dirpath string, snap walpb.Snapshot) (*WAL, error) ***REMOVED***
	return openAtIndex(dirpath, snap, false)
***REMOVED***

func openAtIndex(dirpath string, snap walpb.Snapshot, write bool) (*WAL, error) ***REMOVED***
	names, err := readWalNames(dirpath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	nameIndex, ok := searchIndex(names, snap.Index)
	if !ok || !isValidSeq(names[nameIndex:]) ***REMOVED***
		return nil, ErrFileNotFound
	***REMOVED***

	// open the wal files
	rcs := make([]io.ReadCloser, 0)
	rs := make([]io.Reader, 0)
	ls := make([]*fileutil.LockedFile, 0)
	for _, name := range names[nameIndex:] ***REMOVED***
		p := filepath.Join(dirpath, name)
		if write ***REMOVED***
			l, err := fileutil.TryLockFile(p, os.O_RDWR, fileutil.PrivateFileMode)
			if err != nil ***REMOVED***
				closeAll(rcs...)
				return nil, err
			***REMOVED***
			ls = append(ls, l)
			rcs = append(rcs, l)
		***REMOVED*** else ***REMOVED***
			rf, err := os.OpenFile(p, os.O_RDONLY, fileutil.PrivateFileMode)
			if err != nil ***REMOVED***
				closeAll(rcs...)
				return nil, err
			***REMOVED***
			ls = append(ls, nil)
			rcs = append(rcs, rf)
		***REMOVED***
		rs = append(rs, rcs[len(rcs)-1])
	***REMOVED***

	closer := func() error ***REMOVED*** return closeAll(rcs...) ***REMOVED***

	// create a WAL ready for reading
	w := &WAL***REMOVED***
		dir:       dirpath,
		start:     snap,
		decoder:   newDecoder(rs...),
		readClose: closer,
		locks:     ls,
	***REMOVED***

	if write ***REMOVED***
		// write reuses the file descriptors from read; don't close so
		// WAL can append without dropping the file lock
		w.readClose = nil
		if _, _, err := parseWalName(filepath.Base(w.tail().Name())); err != nil ***REMOVED***
			closer()
			return nil, err
		***REMOVED***
		w.fp = newFilePipeline(w.dir, SegmentSizeBytes)
	***REMOVED***

	return w, nil
***REMOVED***

// ReadAll reads out records of the current WAL.
// If opened in write mode, it must read out all records until EOF. Or an error
// will be returned.
// If opened in read mode, it will try to read all records if possible.
// If it cannot read out the expected snap, it will return ErrSnapshotNotFound.
// If loaded snap doesn't match with the expected one, it will return
// all the records and error ErrSnapshotMismatch.
// TODO: detect not-last-snap error.
// TODO: maybe loose the checking of match.
// After ReadAll, the WAL will be ready for appending new records.
func (w *WAL) ReadAll() (metadata []byte, state raftpb.HardState, ents []raftpb.Entry, err error) ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	rec := &walpb.Record***REMOVED******REMOVED***
	decoder := w.decoder

	var match bool
	for err = decoder.decode(rec); err == nil; err = decoder.decode(rec) ***REMOVED***
		switch rec.Type ***REMOVED***
		case entryType:
			e := mustUnmarshalEntry(rec.Data)
			if e.Index > w.start.Index ***REMOVED***
				ents = append(ents[:e.Index-w.start.Index-1], e)
			***REMOVED***
			w.enti = e.Index
		case stateType:
			state = mustUnmarshalState(rec.Data)
		case metadataType:
			if metadata != nil && !bytes.Equal(metadata, rec.Data) ***REMOVED***
				state.Reset()
				return nil, state, nil, ErrMetadataConflict
			***REMOVED***
			metadata = rec.Data
		case crcType:
			crc := decoder.crc.Sum32()
			// current crc of decoder must match the crc of the record.
			// do no need to match 0 crc, since the decoder is a new one at this case.
			if crc != 0 && rec.Validate(crc) != nil ***REMOVED***
				state.Reset()
				return nil, state, nil, ErrCRCMismatch
			***REMOVED***
			decoder.updateCRC(rec.Crc)
		case snapshotType:
			var snap walpb.Snapshot
			pbutil.MustUnmarshal(&snap, rec.Data)
			if snap.Index == w.start.Index ***REMOVED***
				if snap.Term != w.start.Term ***REMOVED***
					state.Reset()
					return nil, state, nil, ErrSnapshotMismatch
				***REMOVED***
				match = true
			***REMOVED***
		default:
			state.Reset()
			return nil, state, nil, fmt.Errorf("unexpected block type %d", rec.Type)
		***REMOVED***
	***REMOVED***

	switch w.tail() ***REMOVED***
	case nil:
		// We do not have to read out all entries in read mode.
		// The last record maybe a partial written one, so
		// ErrunexpectedEOF might be returned.
		if err != io.EOF && err != io.ErrUnexpectedEOF ***REMOVED***
			state.Reset()
			return nil, state, nil, err
		***REMOVED***
	default:
		// We must read all of the entries if WAL is opened in write mode.
		if err != io.EOF ***REMOVED***
			state.Reset()
			return nil, state, nil, err
		***REMOVED***
		// decodeRecord() will return io.EOF if it detects a zero record,
		// but this zero record may be followed by non-zero records from
		// a torn write. Overwriting some of these non-zero records, but
		// not all, will cause CRC errors on WAL open. Since the records
		// were never fully synced to disk in the first place, it's safe
		// to zero them out to avoid any CRC errors from new writes.
		if _, err = w.tail().Seek(w.decoder.lastOffset(), io.SeekStart); err != nil ***REMOVED***
			return nil, state, nil, err
		***REMOVED***
		if err = fileutil.ZeroToEnd(w.tail().File); err != nil ***REMOVED***
			return nil, state, nil, err
		***REMOVED***
	***REMOVED***

	err = nil
	if !match ***REMOVED***
		err = ErrSnapshotNotFound
	***REMOVED***

	// close decoder, disable reading
	if w.readClose != nil ***REMOVED***
		w.readClose()
		w.readClose = nil
	***REMOVED***
	w.start = walpb.Snapshot***REMOVED******REMOVED***

	w.metadata = metadata

	if w.tail() != nil ***REMOVED***
		// create encoder (chain crc with the decoder), enable appending
		w.encoder, err = newFileEncoder(w.tail().File, w.decoder.lastCRC())
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	w.decoder = nil

	return metadata, state, ents, err
***REMOVED***

// cut closes current file written and creates a new one ready to append.
// cut first creates a temp wal file and writes necessary headers into it.
// Then cut atomically rename temp wal file to a wal file.
func (w *WAL) cut() error ***REMOVED***
	// close old wal file; truncate to avoid wasting space if an early cut
	off, serr := w.tail().Seek(0, io.SeekCurrent)
	if serr != nil ***REMOVED***
		return serr
	***REMOVED***
	if err := w.tail().Truncate(off); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := w.sync(); err != nil ***REMOVED***
		return err
	***REMOVED***

	fpath := filepath.Join(w.dir, walName(w.seq()+1, w.enti+1))

	// create a temp wal file with name sequence + 1, or truncate the existing one
	newTail, err := w.fp.Open()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// update writer and save the previous crc
	w.locks = append(w.locks, newTail)
	prevCrc := w.encoder.crc.Sum32()
	w.encoder, err = newFileEncoder(w.tail().File, prevCrc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = w.saveCrc(prevCrc); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = w.encoder.encode(&walpb.Record***REMOVED***Type: metadataType, Data: w.metadata***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = w.saveState(&w.state); err != nil ***REMOVED***
		return err
	***REMOVED***
	// atomically move temp wal file to wal file
	if err = w.sync(); err != nil ***REMOVED***
		return err
	***REMOVED***

	off, err = w.tail().Seek(0, io.SeekCurrent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = os.Rename(newTail.Name(), fpath); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = fileutil.Fsync(w.dirFile); err != nil ***REMOVED***
		return err
	***REMOVED***

	newTail.Close()

	if newTail, err = fileutil.LockFile(fpath, os.O_WRONLY, fileutil.PrivateFileMode); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err = newTail.Seek(off, io.SeekStart); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.locks[len(w.locks)-1] = newTail

	prevCrc = w.encoder.crc.Sum32()
	w.encoder, err = newFileEncoder(w.tail().File, prevCrc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	plog.Infof("segmented wal file %v is created", fpath)
	return nil
***REMOVED***

func (w *WAL) sync() error ***REMOVED***
	if w.encoder != nil ***REMOVED***
		if err := w.encoder.flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	start := time.Now()
	err := fileutil.Fdatasync(w.tail().File)

	duration := time.Since(start)
	if duration > warnSyncDuration ***REMOVED***
		plog.Warningf("sync duration of %v, expected less than %v", duration, warnSyncDuration)
	***REMOVED***
	syncDurations.Observe(duration.Seconds())

	return err
***REMOVED***

// ReleaseLockTo releases the locks, which has smaller index than the given index
// except the largest one among them.
// For example, if WAL is holding lock 1,2,3,4,5,6, ReleaseLockTo(4) will release
// lock 1,2 but keep 3. ReleaseLockTo(5) will release 1,2,3 but keep 4.
func (w *WAL) ReleaseLockTo(index uint64) error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	var smaller int
	found := false

	for i, l := range w.locks ***REMOVED***
		_, lockIndex, err := parseWalName(filepath.Base(l.Name()))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if lockIndex >= index ***REMOVED***
			smaller = i - 1
			found = true
			break
		***REMOVED***
	***REMOVED***

	// if no lock index is greater than the release index, we can
	// release lock up to the last one(excluding).
	if !found && len(w.locks) != 0 ***REMOVED***
		smaller = len(w.locks) - 1
	***REMOVED***

	if smaller <= 0 ***REMOVED***
		return nil
	***REMOVED***

	for i := 0; i < smaller; i++ ***REMOVED***
		if w.locks[i] == nil ***REMOVED***
			continue
		***REMOVED***
		w.locks[i].Close()
	***REMOVED***
	w.locks = w.locks[smaller:]

	return nil
***REMOVED***

func (w *WAL) Close() error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.fp != nil ***REMOVED***
		w.fp.Close()
		w.fp = nil
	***REMOVED***

	if w.tail() != nil ***REMOVED***
		if err := w.sync(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, l := range w.locks ***REMOVED***
		if l == nil ***REMOVED***
			continue
		***REMOVED***
		if err := l.Close(); err != nil ***REMOVED***
			plog.Errorf("failed to unlock during closing wal: %s", err)
		***REMOVED***
	***REMOVED***

	return w.dirFile.Close()
***REMOVED***

func (w *WAL) saveEntry(e *raftpb.Entry) error ***REMOVED***
	// TODO: add MustMarshalTo to reduce one allocation.
	b := pbutil.MustMarshal(e)
	rec := &walpb.Record***REMOVED***Type: entryType, Data: b***REMOVED***
	if err := w.encoder.encode(rec); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.enti = e.Index
	return nil
***REMOVED***

func (w *WAL) saveState(s *raftpb.HardState) error ***REMOVED***
	if raft.IsEmptyHardState(*s) ***REMOVED***
		return nil
	***REMOVED***
	w.state = *s
	b := pbutil.MustMarshal(s)
	rec := &walpb.Record***REMOVED***Type: stateType, Data: b***REMOVED***
	return w.encoder.encode(rec)
***REMOVED***

func (w *WAL) Save(st raftpb.HardState, ents []raftpb.Entry) error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	// short cut, do not call sync
	if raft.IsEmptyHardState(st) && len(ents) == 0 ***REMOVED***
		return nil
	***REMOVED***

	mustSync := raft.MustSync(st, w.state, len(ents))

	// TODO(xiangli): no more reference operator
	for i := range ents ***REMOVED***
		if err := w.saveEntry(&ents[i]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := w.saveState(&st); err != nil ***REMOVED***
		return err
	***REMOVED***

	curOff, err := w.tail().Seek(0, io.SeekCurrent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if curOff < SegmentSizeBytes ***REMOVED***
		if mustSync ***REMOVED***
			return w.sync()
		***REMOVED***
		return nil
	***REMOVED***

	return w.cut()
***REMOVED***

func (w *WAL) SaveSnapshot(e walpb.Snapshot) error ***REMOVED***
	b := pbutil.MustMarshal(&e)

	w.mu.Lock()
	defer w.mu.Unlock()

	rec := &walpb.Record***REMOVED***Type: snapshotType, Data: b***REMOVED***
	if err := w.encoder.encode(rec); err != nil ***REMOVED***
		return err
	***REMOVED***
	// update enti only when snapshot is ahead of last index
	if w.enti < e.Index ***REMOVED***
		w.enti = e.Index
	***REMOVED***
	return w.sync()
***REMOVED***

func (w *WAL) saveCrc(prevCrc uint32) error ***REMOVED***
	return w.encoder.encode(&walpb.Record***REMOVED***Type: crcType, Crc: prevCrc***REMOVED***)
***REMOVED***

func (w *WAL) tail() *fileutil.LockedFile ***REMOVED***
	if len(w.locks) > 0 ***REMOVED***
		return w.locks[len(w.locks)-1]
	***REMOVED***
	return nil
***REMOVED***

func (w *WAL) seq() uint64 ***REMOVED***
	t := w.tail()
	if t == nil ***REMOVED***
		return 0
	***REMOVED***
	seq, _, err := parseWalName(filepath.Base(t.Name()))
	if err != nil ***REMOVED***
		plog.Fatalf("bad wal name %s (%v)", t.Name(), err)
	***REMOVED***
	return seq
***REMOVED***

func closeAll(rcs ...io.ReadCloser) error ***REMOVED***
	for _, f := range rcs ***REMOVED***
		if err := f.Close(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
