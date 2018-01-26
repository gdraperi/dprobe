package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/snap"
	"github.com/coreos/etcd/wal"
	"github.com/coreos/etcd/wal/walpb"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/encryption"
	"github.com/pkg/errors"
)

// ErrNoWAL is returned if there are no WALs on disk
var ErrNoWAL = errors.New("no WAL present")

type walSnapDirs struct ***REMOVED***
	wal  string
	snap string
***REMOVED***

// the wal/snap directories in decreasing order of preference/version
var versionedWALSnapDirs = []walSnapDirs***REMOVED***
	***REMOVED***wal: "wal-v3-encrypted", snap: "snap-v3-encrypted"***REMOVED***,
	***REMOVED***wal: "wal-v3", snap: "snap-v3"***REMOVED***,
	***REMOVED***wal: "wal", snap: "snap"***REMOVED***,
***REMOVED***

// MultiDecrypter attempts to decrypt with a list of decrypters
type MultiDecrypter []encryption.Decrypter

// Decrypt tries to decrypt using all the decrypters
func (m MultiDecrypter) Decrypt(r api.MaybeEncryptedRecord) (result []byte, err error) ***REMOVED***
	for _, d := range m ***REMOVED***
		result, err = d.Decrypt(r)
		if err == nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// EncryptedRaftLogger saves raft data to disk
type EncryptedRaftLogger struct ***REMOVED***
	StateDir      string
	EncryptionKey []byte

	// mutex is locked for writing only when we need to replace the wal object and snapshotter
	// object, not when we're writing snapshots or wals (in which case it's locked for reading)
	encoderMu   sync.RWMutex
	wal         WAL
	snapshotter Snapshotter
***REMOVED***

// BootstrapFromDisk creates a new snapshotter and wal, and also reads the latest snapshot and WALs from disk
func (e *EncryptedRaftLogger) BootstrapFromDisk(ctx context.Context, oldEncryptionKeys ...[]byte) (*raftpb.Snapshot, WALData, error) ***REMOVED***
	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()

	walDir := e.walDir()
	snapDir := e.snapDir()

	encrypter, decrypter := encryption.Defaults(e.EncryptionKey)
	if oldEncryptionKeys != nil ***REMOVED***
		decrypters := []encryption.Decrypter***REMOVED***decrypter***REMOVED***
		for _, key := range oldEncryptionKeys ***REMOVED***
			_, d := encryption.Defaults(key)
			decrypters = append(decrypters, d)
		***REMOVED***
		decrypter = MultiDecrypter(decrypters)
	***REMOVED***

	snapFactory := NewSnapFactory(encrypter, decrypter)

	if !fileutil.Exist(snapDir) ***REMOVED***
		// If snapshots created by the etcd-v2 code exist, or by swarmkit development version,
		// read the latest snapshot and write it encoded to the new path.  The new path
		// prevents etc-v2 creating snapshots that are visible to us, but not encoded and
		// out of sync with our WALs, after a downgrade.
		for _, dirs := range versionedWALSnapDirs[1:] ***REMOVED***
			legacySnapDir := filepath.Join(e.StateDir, dirs.snap)
			if fileutil.Exist(legacySnapDir) ***REMOVED***
				if err := MigrateSnapshot(legacySnapDir, snapDir, OriginalSnap, snapFactory); err != nil ***REMOVED***
					return nil, WALData***REMOVED******REMOVED***, err
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// ensure the new directory exists
	if err := os.MkdirAll(snapDir, 0700); err != nil ***REMOVED***
		return nil, WALData***REMOVED******REMOVED***, errors.Wrap(err, "failed to create snapshot directory")
	***REMOVED***

	var (
		snapshotter Snapshotter
		walObj      WAL
		err         error
	)

	// Create a snapshotter and load snapshot data
	snapshotter = snapFactory.New(snapDir)
	snapshot, err := snapshotter.Load()
	if err != nil && err != snap.ErrNoSnapshot ***REMOVED***
		return nil, WALData***REMOVED******REMOVED***, err
	***REMOVED***

	walFactory := NewWALFactory(encrypter, decrypter)
	var walsnap walpb.Snapshot
	if snapshot != nil ***REMOVED***
		walsnap.Index = snapshot.Metadata.Index
		walsnap.Term = snapshot.Metadata.Term
	***REMOVED***

	if !wal.Exist(walDir) ***REMOVED***
		var walExists bool
		// If wals created by the etcd-v2 wal code exist, read the latest ones based
		// on this snapshot and encode them to wals in the new path to avoid adding
		// backwards-incompatible entries to those files.
		for _, dirs := range versionedWALSnapDirs[1:] ***REMOVED***
			legacyWALDir := filepath.Join(e.StateDir, dirs.wal)
			if !wal.Exist(legacyWALDir) ***REMOVED***
				continue
			***REMOVED***
			if err = MigrateWALs(ctx, legacyWALDir, walDir, OriginalWAL, walFactory, walsnap); err != nil ***REMOVED***
				return nil, WALData***REMOVED******REMOVED***, err
			***REMOVED***
			walExists = true
			break
		***REMOVED***
		if !walExists ***REMOVED***
			return nil, WALData***REMOVED******REMOVED***, ErrNoWAL
		***REMOVED***
	***REMOVED***

	walObj, waldata, err := ReadRepairWAL(ctx, walDir, walsnap, walFactory)
	if err != nil ***REMOVED***
		return nil, WALData***REMOVED******REMOVED***, err
	***REMOVED***

	e.snapshotter = snapshotter
	e.wal = walObj

	return snapshot, waldata, nil
***REMOVED***

// BootstrapNew creates a new snapshotter and WAL writer, expecting that there is nothing on disk
func (e *EncryptedRaftLogger) BootstrapNew(metadata []byte) error ***REMOVED***
	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()
	encrypter, decrypter := encryption.Defaults(e.EncryptionKey)
	walFactory := NewWALFactory(encrypter, decrypter)

	for _, dirpath := range []string***REMOVED***filepath.Dir(e.walDir()), e.snapDir()***REMOVED*** ***REMOVED***
		if err := os.MkdirAll(dirpath, 0700); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to create %s", dirpath)
		***REMOVED***
	***REMOVED***
	var err error
	// the wal directory must not already exist upon creation
	e.wal, err = walFactory.Create(e.walDir(), metadata)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to create WAL")
	***REMOVED***

	e.snapshotter = NewSnapFactory(encrypter, decrypter).New(e.snapDir())
	return nil
***REMOVED***

func (e *EncryptedRaftLogger) walDir() string ***REMOVED***
	return filepath.Join(e.StateDir, versionedWALSnapDirs[0].wal)
***REMOVED***

func (e *EncryptedRaftLogger) snapDir() string ***REMOVED***
	return filepath.Join(e.StateDir, versionedWALSnapDirs[0].snap)
***REMOVED***

// RotateEncryptionKey swaps out the encoders and decoders used by the wal and snapshotter
func (e *EncryptedRaftLogger) RotateEncryptionKey(newKey []byte) ***REMOVED***
	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()

	if e.wal != nil ***REMOVED*** // if the wal exists, the snapshotter exists
		// We don't want to have to close the WAL, because we can't open a new one.
		// We need to know the previous snapshot, because when you open a WAL you
		// have to read out all the entries from a particular snapshot, or you can't
		// write.  So just rotate the encoders out from under it.  We already
		// have a lock on writing to snapshots and WALs.
		wrapped, ok := e.wal.(*wrappedWAL)
		if !ok ***REMOVED***
			panic(fmt.Errorf("EncryptedRaftLogger's WAL is not a wrappedWAL"))
		***REMOVED***

		wrapped.encrypter, wrapped.decrypter = encryption.Defaults(newKey)

		e.snapshotter = NewSnapFactory(wrapped.encrypter, wrapped.decrypter).New(e.snapDir())
	***REMOVED***
	e.EncryptionKey = newKey
***REMOVED***

// SaveSnapshot actually saves a given snapshot to both the WAL and the snapshot.
func (e *EncryptedRaftLogger) SaveSnapshot(snapshot raftpb.Snapshot) error ***REMOVED***

	walsnap := walpb.Snapshot***REMOVED***
		Index: snapshot.Metadata.Index,
		Term:  snapshot.Metadata.Term,
	***REMOVED***

	e.encoderMu.RLock()
	if err := e.wal.SaveSnapshot(walsnap); err != nil ***REMOVED***
		e.encoderMu.RUnlock()
		return err
	***REMOVED***

	snapshotter := e.snapshotter
	e.encoderMu.RUnlock()

	if err := snapshotter.SaveSnap(snapshot); err != nil ***REMOVED***
		return err
	***REMOVED***
	return e.wal.ReleaseLockTo(snapshot.Metadata.Index)
***REMOVED***

// GC garbage collects snapshots and wals older than the provided index and term
func (e *EncryptedRaftLogger) GC(index uint64, term uint64, keepOldSnapshots uint64) error ***REMOVED***
	// Delete any older snapshots
	curSnapshot := fmt.Sprintf("%016x-%016x%s", term, index, ".snap")

	snapshots, err := ListSnapshots(e.snapDir())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Ignore any snapshots that are older than the current snapshot.
	// Delete the others. Rather than doing lexical comparisons, we look
	// at what exists before/after the current snapshot in the slice.
	// This means that if the current snapshot doesn't appear in the
	// directory for some strange reason, we won't delete anything, which
	// is the safe behavior.
	curSnapshotIdx := -1
	var (
		removeErr      error
		oldestSnapshot string
	)

	for i, snapFile := range snapshots ***REMOVED***
		if curSnapshotIdx >= 0 && i > curSnapshotIdx ***REMOVED***
			if uint64(i-curSnapshotIdx) > keepOldSnapshots ***REMOVED***
				err := os.Remove(filepath.Join(e.snapDir(), snapFile))
				if err != nil && removeErr == nil ***REMOVED***
					removeErr = err
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED*** else if snapFile == curSnapshot ***REMOVED***
			curSnapshotIdx = i
		***REMOVED***
		oldestSnapshot = snapFile
	***REMOVED***

	if removeErr != nil ***REMOVED***
		return removeErr
	***REMOVED***

	// Remove any WAL files that only contain data from before the oldest
	// remaining snapshot.

	if oldestSnapshot == "" ***REMOVED***
		return nil
	***REMOVED***

	// Parse index out of oldest snapshot's filename
	var snapTerm, snapIndex uint64
	_, err = fmt.Sscanf(oldestSnapshot, "%016x-%016x.snap", &snapTerm, &snapIndex)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "malformed snapshot filename %s", oldestSnapshot)
	***REMOVED***

	wals, err := ListWALs(e.walDir())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	found := false
	deleteUntil := -1

	for i, walName := range wals ***REMOVED***
		var walSeq, walIndex uint64
		_, err = fmt.Sscanf(walName, "%016x-%016x.wal", &walSeq, &walIndex)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "could not parse WAL name %s", walName)
		***REMOVED***

		if walIndex >= snapIndex ***REMOVED***
			deleteUntil = i - 1
			found = true
			break
		***REMOVED***
	***REMOVED***

	// If all WAL files started with indices below the oldest snapshot's
	// index, we can delete all but the newest WAL file.
	if !found && len(wals) != 0 ***REMOVED***
		deleteUntil = len(wals) - 1
	***REMOVED***

	for i := 0; i < deleteUntil; i++ ***REMOVED***
		walPath := filepath.Join(e.walDir(), wals[i])
		l, err := fileutil.TryLockFile(walPath, os.O_WRONLY, fileutil.PrivateFileMode)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "could not lock old WAL file %s for removal", wals[i])
		***REMOVED***
		err = os.Remove(walPath)
		l.Close()
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "error removing old WAL file %s", wals[i])
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// SaveEntries saves only entries to disk
func (e *EncryptedRaftLogger) SaveEntries(st raftpb.HardState, entries []raftpb.Entry) error ***REMOVED***
	e.encoderMu.RLock()
	defer e.encoderMu.RUnlock()

	if e.wal == nil ***REMOVED***
		return fmt.Errorf("raft WAL has either been closed or has never been created")
	***REMOVED***
	return e.wal.Save(st, entries)
***REMOVED***

// Close closes the logger - it will have to be bootstrapped again to start writing
func (e *EncryptedRaftLogger) Close(ctx context.Context) ***REMOVED***
	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()

	if e.wal != nil ***REMOVED***
		if err := e.wal.Close(); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("error closing raft WAL")
		***REMOVED***
	***REMOVED***

	e.wal = nil
	e.snapshotter = nil
***REMOVED***

// Clear closes the existing WAL and removes the WAL and snapshot.
func (e *EncryptedRaftLogger) Clear(ctx context.Context) error ***REMOVED***
	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()

	if e.wal != nil ***REMOVED***
		if err := e.wal.Close(); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("error closing raft WAL")
		***REMOVED***
	***REMOVED***
	e.snapshotter = nil

	os.RemoveAll(e.walDir())
	os.RemoveAll(e.snapDir())
	return nil
***REMOVED***
