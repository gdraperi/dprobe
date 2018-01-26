package remotecontext

import (
	"archive/tar"
	"crypto/sha256"
	"hash"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/tarsum"
)

// NewFileHash returns new hash that is used for the builder cache keys
func NewFileHash(path, name string, fi os.FileInfo) (hash.Hash, error) ***REMOVED***
	var link string
	if fi.Mode()&os.ModeSymlink != 0 ***REMOVED***
		var err error
		link, err = os.Readlink(path)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	hdr, err := archive.FileInfoHeader(name, fi, link)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := archive.ReadSecurityXattrToTarHeader(path, hdr); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tsh := &tarsumHash***REMOVED***hdr: hdr, Hash: sha256.New()***REMOVED***
	tsh.Reset() // initialize header
	return tsh, nil
***REMOVED***

type tarsumHash struct ***REMOVED***
	hash.Hash
	hdr *tar.Header
***REMOVED***

// Reset resets the Hash to its initial state.
func (tsh *tarsumHash) Reset() ***REMOVED***
	// comply with hash.Hash and reset to the state hash had before any writes
	tsh.Hash.Reset()
	tarsum.WriteV1Header(tsh.hdr, tsh.Hash)
***REMOVED***
