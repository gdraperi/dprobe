package image

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DigestWalkFunc is function called by StoreBackend.Walk
type DigestWalkFunc func(id digest.Digest) error

// StoreBackend provides interface for image.Store persistence
type StoreBackend interface ***REMOVED***
	Walk(f DigestWalkFunc) error
	Get(id digest.Digest) ([]byte, error)
	Set(data []byte) (digest.Digest, error)
	Delete(id digest.Digest) error
	SetMetadata(id digest.Digest, key string, data []byte) error
	GetMetadata(id digest.Digest, key string) ([]byte, error)
	DeleteMetadata(id digest.Digest, key string) error
***REMOVED***

// fs implements StoreBackend using the filesystem.
type fs struct ***REMOVED***
	sync.RWMutex
	root string
***REMOVED***

const (
	contentDirName  = "content"
	metadataDirName = "metadata"
)

// NewFSStoreBackend returns new filesystem based backend for image.Store
func NewFSStoreBackend(root string) (StoreBackend, error) ***REMOVED***
	return newFSStore(root)
***REMOVED***

func newFSStore(root string) (*fs, error) ***REMOVED***
	s := &fs***REMOVED***
		root: root,
	***REMOVED***
	if err := os.MkdirAll(filepath.Join(root, contentDirName, string(digest.Canonical)), 0700); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to create storage backend")
	***REMOVED***
	if err := os.MkdirAll(filepath.Join(root, metadataDirName, string(digest.Canonical)), 0700); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to create storage backend")
	***REMOVED***
	return s, nil
***REMOVED***

func (s *fs) contentFile(dgst digest.Digest) string ***REMOVED***
	return filepath.Join(s.root, contentDirName, string(dgst.Algorithm()), dgst.Hex())
***REMOVED***

func (s *fs) metadataDir(dgst digest.Digest) string ***REMOVED***
	return filepath.Join(s.root, metadataDirName, string(dgst.Algorithm()), dgst.Hex())
***REMOVED***

// Walk calls the supplied callback for each image ID in the storage backend.
func (s *fs) Walk(f DigestWalkFunc) error ***REMOVED***
	// Only Canonical digest (sha256) is currently supported
	s.RLock()
	dir, err := ioutil.ReadDir(filepath.Join(s.root, contentDirName, string(digest.Canonical)))
	s.RUnlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, v := range dir ***REMOVED***
		dgst := digest.NewDigestFromHex(string(digest.Canonical), v.Name())
		if err := dgst.Validate(); err != nil ***REMOVED***
			logrus.Debugf("skipping invalid digest %s: %s", dgst, err)
			continue
		***REMOVED***
		if err := f(dgst); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Get returns the content stored under a given digest.
func (s *fs) Get(dgst digest.Digest) ([]byte, error) ***REMOVED***
	s.RLock()
	defer s.RUnlock()

	return s.get(dgst)
***REMOVED***

func (s *fs) get(dgst digest.Digest) ([]byte, error) ***REMOVED***
	content, err := ioutil.ReadFile(s.contentFile(dgst))
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to get digest %s", dgst)
	***REMOVED***

	// todo: maybe optional
	if digest.FromBytes(content) != dgst ***REMOVED***
		return nil, fmt.Errorf("failed to verify: %v", dgst)
	***REMOVED***

	return content, nil
***REMOVED***

// Set stores content by checksum.
func (s *fs) Set(data []byte) (digest.Digest, error) ***REMOVED***
	s.Lock()
	defer s.Unlock()

	if len(data) == 0 ***REMOVED***
		return "", fmt.Errorf("invalid empty data")
	***REMOVED***

	dgst := digest.FromBytes(data)
	if err := ioutils.AtomicWriteFile(s.contentFile(dgst), data, 0600); err != nil ***REMOVED***
		return "", errors.Wrap(err, "failed to write digest data")
	***REMOVED***

	return dgst, nil
***REMOVED***

// Delete removes content and metadata files associated with the digest.
func (s *fs) Delete(dgst digest.Digest) error ***REMOVED***
	s.Lock()
	defer s.Unlock()

	if err := os.RemoveAll(s.metadataDir(dgst)); err != nil ***REMOVED***
		return err
	***REMOVED***
	return os.Remove(s.contentFile(dgst))
***REMOVED***

// SetMetadata sets metadata for a given ID. It fails if there's no base file.
func (s *fs) SetMetadata(dgst digest.Digest, key string, data []byte) error ***REMOVED***
	s.Lock()
	defer s.Unlock()
	if _, err := s.get(dgst); err != nil ***REMOVED***
		return err
	***REMOVED***

	baseDir := filepath.Join(s.metadataDir(dgst))
	if err := os.MkdirAll(baseDir, 0700); err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutils.AtomicWriteFile(filepath.Join(s.metadataDir(dgst), key), data, 0600)
***REMOVED***

// GetMetadata returns metadata for a given digest.
func (s *fs) GetMetadata(dgst digest.Digest, key string) ([]byte, error) ***REMOVED***
	s.RLock()
	defer s.RUnlock()

	if _, err := s.get(dgst); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	bytes, err := ioutil.ReadFile(filepath.Join(s.metadataDir(dgst), key))
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to read metadata")
	***REMOVED***
	return bytes, nil
***REMOVED***

// DeleteMetadata removes the metadata associated with a digest.
func (s *fs) DeleteMetadata(dgst digest.Digest, key string) error ***REMOVED***
	s.Lock()
	defer s.Unlock()

	return os.RemoveAll(filepath.Join(s.metadataDir(dgst), key))
***REMOVED***
