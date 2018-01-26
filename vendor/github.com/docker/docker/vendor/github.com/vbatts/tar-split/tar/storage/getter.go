package storage

import (
	"bytes"
	"errors"
	"hash/crc64"
	"io"
	"os"
	"path/filepath"
)

// FileGetter is the interface for getting a stream of a file payload,
// addressed by name/filename. Presumably, the names will be scoped to relative
// file paths.
type FileGetter interface ***REMOVED***
	// Get returns a stream for the provided file path
	Get(filename string) (output io.ReadCloser, err error)
***REMOVED***

// FilePutter is the interface for storing a stream of a file payload,
// addressed by name/filename.
type FilePutter interface ***REMOVED***
	// Put returns the size of the stream received, and the crc64 checksum for
	// the provided stream
	Put(filename string, input io.Reader) (size int64, checksum []byte, err error)
***REMOVED***

// FileGetPutter is the interface that groups both Getting and Putting file
// payloads.
type FileGetPutter interface ***REMOVED***
	FileGetter
	FilePutter
***REMOVED***

// NewPathFileGetter returns a FileGetter that is for files relative to path
// relpath.
func NewPathFileGetter(relpath string) FileGetter ***REMOVED***
	return &pathFileGetter***REMOVED***root: relpath***REMOVED***
***REMOVED***

type pathFileGetter struct ***REMOVED***
	root string
***REMOVED***

func (pfg pathFileGetter) Get(filename string) (io.ReadCloser, error) ***REMOVED***
	return os.Open(filepath.Join(pfg.root, filename))
***REMOVED***

type bufferFileGetPutter struct ***REMOVED***
	files map[string][]byte
***REMOVED***

func (bfgp bufferFileGetPutter) Get(name string) (io.ReadCloser, error) ***REMOVED***
	if _, ok := bfgp.files[name]; !ok ***REMOVED***
		return nil, errors.New("no such file")
	***REMOVED***
	b := bytes.NewBuffer(bfgp.files[name])
	return &readCloserWrapper***REMOVED***b***REMOVED***, nil
***REMOVED***

func (bfgp *bufferFileGetPutter) Put(name string, r io.Reader) (int64, []byte, error) ***REMOVED***
	crc := crc64.New(CRCTable)
	buf := bytes.NewBuffer(nil)
	cw := io.MultiWriter(crc, buf)
	i, err := io.Copy(cw, r)
	if err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***
	bfgp.files[name] = buf.Bytes()
	return i, crc.Sum(nil), nil
***REMOVED***

type readCloserWrapper struct ***REMOVED***
	io.Reader
***REMOVED***

func (w *readCloserWrapper) Close() error ***REMOVED*** return nil ***REMOVED***

// NewBufferFileGetPutter is a simple in-memory FileGetPutter
//
// Implication is this is memory intensive...
// Probably best for testing or light weight cases.
func NewBufferFileGetPutter() FileGetPutter ***REMOVED***
	return &bufferFileGetPutter***REMOVED***
		files: map[string][]byte***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// NewDiscardFilePutter is a bit bucket FilePutter
func NewDiscardFilePutter() FilePutter ***REMOVED***
	return &bitBucketFilePutter***REMOVED******REMOVED***
***REMOVED***

type bitBucketFilePutter struct ***REMOVED***
***REMOVED***

func (bbfp *bitBucketFilePutter) Put(name string, r io.Reader) (int64, []byte, error) ***REMOVED***
	c := crc64.New(CRCTable)
	i, err := io.Copy(c, r)
	return i, c.Sum(nil), err
***REMOVED***

// CRCTable is the default table used for crc64 sum calculations
var CRCTable = crc64.MakeTable(crc64.ISO)
