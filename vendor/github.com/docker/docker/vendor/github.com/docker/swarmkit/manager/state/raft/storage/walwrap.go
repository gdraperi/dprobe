package storage

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/wal"
	"github.com/coreos/etcd/wal/walpb"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/encryption"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// This package wraps the github.com/coreos/etcd/wal package, and encrypts
// the bytes of whatever entry is passed to it, and decrypts the bytes of
// whatever entry it reads.

// WAL is the interface presented by github.com/coreos/etcd/wal.WAL that we depend upon
type WAL interface ***REMOVED***
	ReadAll() ([]byte, raftpb.HardState, []raftpb.Entry, error)
	ReleaseLockTo(index uint64) error
	Close() error
	Save(st raftpb.HardState, ents []raftpb.Entry) error
	SaveSnapshot(e walpb.Snapshot) error
***REMOVED***

// WALFactory provides an interface for the different ways to get a WAL object.
// For instance, the etcd/wal package itself provides this
type WALFactory interface ***REMOVED***
	Create(dirpath string, metadata []byte) (WAL, error)
	Open(dirpath string, walsnap walpb.Snapshot) (WAL, error)
***REMOVED***

var _ WAL = &wrappedWAL***REMOVED******REMOVED***
var _ WAL = &wal.WAL***REMOVED******REMOVED***
var _ WALFactory = walCryptor***REMOVED******REMOVED***

// wrappedWAL wraps a github.com/coreos/etcd/wal.WAL, and handles encrypting/decrypting
type wrappedWAL struct ***REMOVED***
	*wal.WAL
	encrypter encryption.Encrypter
	decrypter encryption.Decrypter
***REMOVED***

// ReadAll wraps the wal.WAL.ReadAll() function, but it first checks to see if the
// metadata indicates that the entries are encryptd, and if so, decrypts them.
func (w *wrappedWAL) ReadAll() ([]byte, raftpb.HardState, []raftpb.Entry, error) ***REMOVED***
	metadata, state, ents, err := w.WAL.ReadAll()
	if err != nil ***REMOVED***
		return metadata, state, ents, err
	***REMOVED***
	for i, ent := range ents ***REMOVED***
		ents[i].Data, err = encryption.Decrypt(ent.Data, w.decrypter)
		if err != nil ***REMOVED***
			return nil, raftpb.HardState***REMOVED******REMOVED***, nil, err
		***REMOVED***
	***REMOVED***

	return metadata, state, ents, nil
***REMOVED***

// Save encrypts the entry data (if an encrypter is exists) before passing it onto the
// wrapped wal.WAL's Save function.
func (w *wrappedWAL) Save(st raftpb.HardState, ents []raftpb.Entry) error ***REMOVED***
	var writeEnts []raftpb.Entry
	for _, ent := range ents ***REMOVED***
		data, err := encryption.Encrypt(ent.Data, w.encrypter)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		writeEnts = append(writeEnts, raftpb.Entry***REMOVED***
			Index: ent.Index,
			Term:  ent.Term,
			Type:  ent.Type,
			Data:  data,
		***REMOVED***)
	***REMOVED***

	return w.WAL.Save(st, writeEnts)
***REMOVED***

// walCryptor is an object that provides the same functions as `etcd/wal`
// and `etcd/snap` that we need to open a WAL object or Snapshotter object
type walCryptor struct ***REMOVED***
	encrypter encryption.Encrypter
	decrypter encryption.Decrypter
***REMOVED***

// NewWALFactory returns an object that can be used to produce objects that
// will read from and write to encrypted WALs on disk.
func NewWALFactory(encrypter encryption.Encrypter, decrypter encryption.Decrypter) WALFactory ***REMOVED***
	return walCryptor***REMOVED***
		encrypter: encrypter,
		decrypter: decrypter,
	***REMOVED***
***REMOVED***

// Create returns a new WAL object with the given encrypters and decrypters.
func (wc walCryptor) Create(dirpath string, metadata []byte) (WAL, error) ***REMOVED***
	w, err := wal.Create(dirpath, metadata)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &wrappedWAL***REMOVED***
		WAL:       w,
		encrypter: wc.encrypter,
		decrypter: wc.decrypter,
	***REMOVED***, nil
***REMOVED***

// Open returns a new WAL object with the given encrypters and decrypters.
func (wc walCryptor) Open(dirpath string, snap walpb.Snapshot) (WAL, error) ***REMOVED***
	w, err := wal.Open(dirpath, snap)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &wrappedWAL***REMOVED***
		WAL:       w,
		encrypter: wc.encrypter,
		decrypter: wc.decrypter,
	***REMOVED***, nil
***REMOVED***

type originalWAL struct***REMOVED******REMOVED***

func (o originalWAL) Create(dirpath string, metadata []byte) (WAL, error) ***REMOVED***
	return wal.Create(dirpath, metadata)
***REMOVED***
func (o originalWAL) Open(dirpath string, walsnap walpb.Snapshot) (WAL, error) ***REMOVED***
	return wal.Open(dirpath, walsnap)
***REMOVED***

// OriginalWAL is the original `wal` package as an implementation of the WALFactory interface
var OriginalWAL WALFactory = originalWAL***REMOVED******REMOVED***

// WALData contains all the data returned by a WAL's ReadAll() function
// (metadata, hardwate, and entries)
type WALData struct ***REMOVED***
	Metadata  []byte
	HardState raftpb.HardState
	Entries   []raftpb.Entry
***REMOVED***

// ReadRepairWAL opens a WAL for reading, and attempts to read it.  If we can't read it, attempts to repair
// and read again.
func ReadRepairWAL(
	ctx context.Context,
	walDir string,
	walsnap walpb.Snapshot,
	factory WALFactory,
) (WAL, WALData, error) ***REMOVED***
	var (
		reader   WAL
		metadata []byte
		st       raftpb.HardState
		ents     []raftpb.Entry
		err      error
	)
	repaired := false
	for ***REMOVED***
		if reader, err = factory.Open(walDir, walsnap); err != nil ***REMOVED***
			return nil, WALData***REMOVED******REMOVED***, errors.Wrap(err, "failed to open WAL")
		***REMOVED***
		if metadata, st, ents, err = reader.ReadAll(); err != nil ***REMOVED***
			if closeErr := reader.Close(); closeErr != nil ***REMOVED***
				return nil, WALData***REMOVED******REMOVED***, closeErr
			***REMOVED***
			if _, ok := err.(encryption.ErrCannotDecrypt); ok ***REMOVED***
				return nil, WALData***REMOVED******REMOVED***, errors.Wrap(err, "failed to decrypt WAL")
			***REMOVED***
			// we can only repair ErrUnexpectedEOF and we never repair twice.
			if repaired || err != io.ErrUnexpectedEOF ***REMOVED***
				return nil, WALData***REMOVED******REMOVED***, errors.Wrap(err, "irreparable WAL error")
			***REMOVED***
			if !wal.Repair(walDir) ***REMOVED***
				return nil, WALData***REMOVED******REMOVED***, errors.Wrap(err, "WAL error cannot be repaired")
			***REMOVED***
			log.G(ctx).WithError(err).Info("repaired WAL error")
			repaired = true
			continue
		***REMOVED***
		break
	***REMOVED***
	return reader, WALData***REMOVED***
		Metadata:  metadata,
		HardState: st,
		Entries:   ents,
	***REMOVED***, nil
***REMOVED***

// MigrateWALs reads existing WALs (from a particular snapshot and beyond) from one directory, encoded one way,
// and writes them to a new directory, encoded a different way
func MigrateWALs(ctx context.Context, oldDir, newDir string, oldFactory, newFactory WALFactory, snapshot walpb.Snapshot) error ***REMOVED***
	oldReader, waldata, err := ReadRepairWAL(ctx, oldDir, snapshot, oldFactory)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	oldReader.Close()

	if err := os.MkdirAll(filepath.Dir(newDir), 0700); err != nil ***REMOVED***
		return errors.Wrap(err, "could not create parent directory")
	***REMOVED***

	// keep temporary wal directory so WAL initialization appears atomic
	tmpdirpath := filepath.Clean(newDir) + ".tmp"
	if err := os.RemoveAll(tmpdirpath); err != nil ***REMOVED***
		return errors.Wrap(err, "could not remove temporary WAL directory")
	***REMOVED***
	defer os.RemoveAll(tmpdirpath)

	tmpWAL, err := newFactory.Create(tmpdirpath, waldata.Metadata)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "could not create new WAL in temporary WAL directory")
	***REMOVED***
	defer tmpWAL.Close()

	if err := tmpWAL.SaveSnapshot(snapshot); err != nil ***REMOVED***
		return errors.Wrap(err, "could not write WAL snapshot in temporary directory")
	***REMOVED***

	if err := tmpWAL.Save(waldata.HardState, waldata.Entries); err != nil ***REMOVED***
		return errors.Wrap(err, "could not migrate WALs to temporary directory")
	***REMOVED***
	if err := tmpWAL.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return os.Rename(tmpdirpath, newDir)
***REMOVED***

// ListWALs lists all the wals in a directory and returns the list in lexical
// order (oldest first)
func ListWALs(dirpath string) ([]string, error) ***REMOVED***
	dirents, err := ioutil.ReadDir(dirpath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var wals []string
	for _, dirent := range dirents ***REMOVED***
		if strings.HasSuffix(dirent.Name(), ".wal") ***REMOVED***
			wals = append(wals, dirent.Name())
		***REMOVED***
	***REMOVED***

	// Sort WAL filenames in lexical order
	sort.Sort(sort.StringSlice(wals))
	return wals, nil
***REMOVED***
