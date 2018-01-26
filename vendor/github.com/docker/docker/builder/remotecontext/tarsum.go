package remotecontext

import (
	"os"
	"sync"

	"github.com/docker/docker/pkg/containerfs"
	iradix "github.com/hashicorp/go-immutable-radix"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fsutil"
)

type hashed interface ***REMOVED***
	Digest() digest.Digest
***REMOVED***

// CachableSource is a source that contains cache records for its contents
type CachableSource struct ***REMOVED***
	mu   sync.Mutex
	root containerfs.ContainerFS
	tree *iradix.Tree
	txn  *iradix.Txn
***REMOVED***

// NewCachableSource creates new CachableSource
func NewCachableSource(root string) *CachableSource ***REMOVED***
	ts := &CachableSource***REMOVED***
		tree: iradix.New(),
		root: containerfs.NewLocalContainerFS(root),
	***REMOVED***
	return ts
***REMOVED***

// MarshalBinary marshals current cache information to a byte array
func (cs *CachableSource) MarshalBinary() ([]byte, error) ***REMOVED***
	b := TarsumBackup***REMOVED***Hashes: make(map[string]string)***REMOVED***
	root := cs.getRoot()
	root.Walk(func(k []byte, v interface***REMOVED******REMOVED***) bool ***REMOVED***
		b.Hashes[string(k)] = v.(*fileInfo).sum
		return false
	***REMOVED***)
	return b.Marshal()
***REMOVED***

// UnmarshalBinary decodes cache information for presented byte array
func (cs *CachableSource) UnmarshalBinary(data []byte) error ***REMOVED***
	var b TarsumBackup
	if err := b.Unmarshal(data); err != nil ***REMOVED***
		return err
	***REMOVED***
	txn := iradix.New().Txn()
	for p, v := range b.Hashes ***REMOVED***
		txn.Insert([]byte(p), &fileInfo***REMOVED***sum: v***REMOVED***)
	***REMOVED***
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.tree = txn.Commit()
	return nil
***REMOVED***

// Scan rescans the cache information from the file system
func (cs *CachableSource) Scan() error ***REMOVED***
	lc, err := NewLazySource(cs.root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	txn := iradix.New().Txn()
	err = cs.root.Walk(cs.root.Path(), func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to walk %s", path)
		***REMOVED***
		rel, err := Rel(cs.root, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		h, err := lc.Hash(rel)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		txn.Insert([]byte(rel), &fileInfo***REMOVED***sum: h***REMOVED***)
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.tree = txn.Commit()
	return nil
***REMOVED***

// HandleChange notifies the source about a modification operation
func (cs *CachableSource) HandleChange(kind fsutil.ChangeKind, p string, fi os.FileInfo, err error) (retErr error) ***REMOVED***
	cs.mu.Lock()
	if cs.txn == nil ***REMOVED***
		cs.txn = cs.tree.Txn()
	***REMOVED***
	if kind == fsutil.ChangeKindDelete ***REMOVED***
		cs.txn.Delete([]byte(p))
		cs.mu.Unlock()
		return
	***REMOVED***

	h, ok := fi.(hashed)
	if !ok ***REMOVED***
		cs.mu.Unlock()
		return errors.Errorf("invalid fileinfo: %s", p)
	***REMOVED***

	hfi := &fileInfo***REMOVED***
		sum: h.Digest().Hex(),
	***REMOVED***
	cs.txn.Insert([]byte(p), hfi)
	cs.mu.Unlock()
	return nil
***REMOVED***

func (cs *CachableSource) getRoot() *iradix.Node ***REMOVED***
	cs.mu.Lock()
	if cs.txn != nil ***REMOVED***
		cs.tree = cs.txn.Commit()
		cs.txn = nil
	***REMOVED***
	t := cs.tree
	cs.mu.Unlock()
	return t.Root()
***REMOVED***

// Close closes the source
func (cs *CachableSource) Close() error ***REMOVED***
	return nil
***REMOVED***

// Hash returns a hash for a single file in the source
func (cs *CachableSource) Hash(path string) (string, error) ***REMOVED***
	n := cs.getRoot()
	// TODO: check this for symlinks
	v, ok := n.Get([]byte(path))
	if !ok ***REMOVED***
		return path, nil
	***REMOVED***
	return v.(*fileInfo).sum, nil
***REMOVED***

// Root returns a root directory for the source
func (cs *CachableSource) Root() containerfs.ContainerFS ***REMOVED***
	return cs.root
***REMOVED***

type fileInfo struct ***REMOVED***
	sum string
***REMOVED***

func (fi *fileInfo) Hash() string ***REMOVED***
	return fi.sum
***REMOVED***
