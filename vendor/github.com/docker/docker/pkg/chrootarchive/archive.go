package chrootarchive

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
)

// NewArchiver returns a new Archiver which uses chrootarchive.Untar
func NewArchiver(idMappings *idtools.IDMappings) *archive.Archiver ***REMOVED***
	if idMappings == nil ***REMOVED***
		idMappings = &idtools.IDMappings***REMOVED******REMOVED***
	***REMOVED***
	return &archive.Archiver***REMOVED***
		Untar:         Untar,
		IDMappingsVar: idMappings,
	***REMOVED***
***REMOVED***

// Untar reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive may be compressed with one of the following algorithms:
//  identity (uncompressed), gzip, bzip2, xz.
func Untar(tarArchive io.Reader, dest string, options *archive.TarOptions) error ***REMOVED***
	return untarHandler(tarArchive, dest, options, true)
***REMOVED***

// UntarUncompressed reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive must be an uncompressed stream.
func UntarUncompressed(tarArchive io.Reader, dest string, options *archive.TarOptions) error ***REMOVED***
	return untarHandler(tarArchive, dest, options, false)
***REMOVED***

// Handler for teasing out the automatic decompression
func untarHandler(tarArchive io.Reader, dest string, options *archive.TarOptions, decompress bool) error ***REMOVED***
	if tarArchive == nil ***REMOVED***
		return fmt.Errorf("Empty archive")
	***REMOVED***
	if options == nil ***REMOVED***
		options = &archive.TarOptions***REMOVED******REMOVED***
	***REMOVED***
	if options.ExcludePatterns == nil ***REMOVED***
		options.ExcludePatterns = []string***REMOVED******REMOVED***
	***REMOVED***

	idMappings := idtools.NewIDMappingsFromMaps(options.UIDMaps, options.GIDMaps)
	rootIDs := idMappings.RootPair()

	dest = filepath.Clean(dest)
	if _, err := os.Stat(dest); os.IsNotExist(err) ***REMOVED***
		if err := idtools.MkdirAllAndChownNew(dest, 0755, rootIDs); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	r := ioutil.NopCloser(tarArchive)
	if decompress ***REMOVED***
		decompressedArchive, err := archive.DecompressStream(tarArchive)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer decompressedArchive.Close()
		r = decompressedArchive
	***REMOVED***

	return invokeUnpack(r, dest, options)
***REMOVED***
