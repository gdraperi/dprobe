package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/snap"
	"github.com/docker/swarmkit/manager/encryption"
	"github.com/pkg/errors"
)

// This package wraps the github.com/coreos/etcd/snap package, and encrypts
// the bytes of whatever snapshot is passed to it, and decrypts the bytes of
// whatever snapshot it reads.

// Snapshotter is the interface presented by github.com/coreos/etcd/snap.Snapshotter that we depend upon
type Snapshotter interface ***REMOVED***
	SaveSnap(snapshot raftpb.Snapshot) error
	Load() (*raftpb.Snapshot, error)
***REMOVED***

// SnapFactory provides an interface for the different ways to get a Snapshotter object.
// For instance, the etcd/snap package itself provides this
type SnapFactory interface ***REMOVED***
	New(dirpath string) Snapshotter
***REMOVED***

var _ Snapshotter = &wrappedSnap***REMOVED******REMOVED***
var _ Snapshotter = &snap.Snapshotter***REMOVED******REMOVED***
var _ SnapFactory = snapCryptor***REMOVED******REMOVED***

// wrappedSnap wraps a github.com/coreos/etcd/snap.Snapshotter, and handles
// encrypting/decrypting.
type wrappedSnap struct ***REMOVED***
	*snap.Snapshotter
	encrypter encryption.Encrypter
	decrypter encryption.Decrypter
***REMOVED***

// SaveSnap encrypts the snapshot data (if an encrypter is exists) before passing it onto the
// wrapped snap.Snapshotter's SaveSnap function.
func (s *wrappedSnap) SaveSnap(snapshot raftpb.Snapshot) error ***REMOVED***
	toWrite := snapshot
	var err error
	toWrite.Data, err = encryption.Encrypt(snapshot.Data, s.encrypter)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.Snapshotter.SaveSnap(toWrite)
***REMOVED***

// Load decrypts the snapshot data (if a decrypter is exists) after reading it using the
// wrapped snap.Snapshotter's Load function.
func (s *wrappedSnap) Load() (*raftpb.Snapshot, error) ***REMOVED***
	snapshot, err := s.Snapshotter.Load()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	snapshot.Data, err = encryption.Decrypt(snapshot.Data, s.decrypter)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return snapshot, nil
***REMOVED***

// snapCryptor is an object that provides the same functions as `etcd/wal`
// and `etcd/snap` that we need to open a WAL object or Snapshotter object
type snapCryptor struct ***REMOVED***
	encrypter encryption.Encrypter
	decrypter encryption.Decrypter
***REMOVED***

// NewSnapFactory returns a new object that can read from and write to encrypted
// snapshots on disk
func NewSnapFactory(encrypter encryption.Encrypter, decrypter encryption.Decrypter) SnapFactory ***REMOVED***
	return snapCryptor***REMOVED***
		encrypter: encrypter,
		decrypter: decrypter,
	***REMOVED***
***REMOVED***

// NewSnapshotter returns a new Snapshotter with the given encrypters and decrypters
func (sc snapCryptor) New(dirpath string) Snapshotter ***REMOVED***
	return &wrappedSnap***REMOVED***
		Snapshotter: snap.New(dirpath),
		encrypter:   sc.encrypter,
		decrypter:   sc.decrypter,
	***REMOVED***
***REMOVED***

type originalSnap struct***REMOVED******REMOVED***

func (o originalSnap) New(dirpath string) Snapshotter ***REMOVED***
	return snap.New(dirpath)
***REMOVED***

// OriginalSnap is the original `snap` package as an implementation of the SnapFactory interface
var OriginalSnap SnapFactory = originalSnap***REMOVED******REMOVED***

// MigrateSnapshot reads the latest existing snapshot from one directory, encoded one way, and writes
// it to a new directory, encoded a different way
func MigrateSnapshot(oldDir, newDir string, oldFactory, newFactory SnapFactory) error ***REMOVED***
	// use temporary snapshot directory so initialization appears atomic
	oldSnapshotter := oldFactory.New(oldDir)
	snapshot, err := oldSnapshotter.Load()
	switch err ***REMOVED***
	case snap.ErrNoSnapshot: // if there's no snapshot, the migration succeeded
		return nil
	case nil:
		break
	default:
		return err
	***REMOVED***

	tmpdirpath := filepath.Clean(newDir) + ".tmp"
	if fileutil.Exist(tmpdirpath) ***REMOVED***
		if err := os.RemoveAll(tmpdirpath); err != nil ***REMOVED***
			return errors.Wrap(err, "could not remove temporary snapshot directory")
		***REMOVED***
	***REMOVED***
	if err := fileutil.CreateDirAll(tmpdirpath); err != nil ***REMOVED***
		return errors.Wrap(err, "could not create temporary snapshot directory")
	***REMOVED***
	tmpSnapshotter := newFactory.New(tmpdirpath)

	// write the new snapshot to the temporary location
	if err = tmpSnapshotter.SaveSnap(*snapshot); err != nil ***REMOVED***
		return err
	***REMOVED***

	return os.Rename(tmpdirpath, newDir)
***REMOVED***

// ListSnapshots lists all the snapshot files in a particular directory and returns
// the snapshot files in reverse lexical order (newest first)
func ListSnapshots(dirpath string) ([]string, error) ***REMOVED***
	dirents, err := ioutil.ReadDir(dirpath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var snapshots []string
	for _, dirent := range dirents ***REMOVED***
		if strings.HasSuffix(dirent.Name(), ".snap") ***REMOVED***
			snapshots = append(snapshots, dirent.Name())
		***REMOVED***
	***REMOVED***

	// Sort snapshot filenames in reverse lexical order
	sort.Sort(sort.Reverse(sort.StringSlice(snapshots)))
	return snapshots, nil
***REMOVED***
