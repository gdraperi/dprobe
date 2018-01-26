package remotefs

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/docker/docker/pkg/archive"
)

// ReadError is an utility function that reads a serialized error from the given reader
// and deserializes it.
func ReadError(in io.Reader) (*ExportedError, error) ***REMOVED***
	b, err := ioutil.ReadAll(in)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// No error
	if len(b) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	var exportedErr ExportedError
	if err := json.Unmarshal(b, &exportedErr); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &exportedErr, nil
***REMOVED***

// ExportedToError will convert a ExportedError to an error. It will try to match
// the error to any existing known error like os.ErrNotExist. Otherwise, it will just
// return an implementation of the error interface.
func ExportedToError(ee *ExportedError) error ***REMOVED***
	if ee.Error() == os.ErrNotExist.Error() ***REMOVED***
		return os.ErrNotExist
	***REMOVED*** else if ee.Error() == os.ErrExist.Error() ***REMOVED***
		return os.ErrExist
	***REMOVED*** else if ee.Error() == os.ErrPermission.Error() ***REMOVED***
		return os.ErrPermission
	***REMOVED*** else if ee.Error() == io.EOF.Error() ***REMOVED***
		return io.EOF
	***REMOVED***
	return ee
***REMOVED***

// WriteError is an utility function that serializes the error
// and writes it to the output writer.
func WriteError(err error, out io.Writer) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	err = fixOSError(err)

	var errno int
	switch typedError := err.(type) ***REMOVED***
	case *os.PathError:
		if se, ok := typedError.Err.(syscall.Errno); ok ***REMOVED***
			errno = int(se)
		***REMOVED***
	case *os.LinkError:
		if se, ok := typedError.Err.(syscall.Errno); ok ***REMOVED***
			errno = int(se)
		***REMOVED***
	case *os.SyscallError:
		if se, ok := typedError.Err.(syscall.Errno); ok ***REMOVED***
			errno = int(se)
		***REMOVED***
	***REMOVED***

	exportedError := &ExportedError***REMOVED***
		ErrString: err.Error(),
		ErrNum:    errno,
	***REMOVED***

	b, err1 := json.Marshal(exportedError)
	if err1 != nil ***REMOVED***
		return err1
	***REMOVED***

	_, err1 = out.Write(b)
	if err1 != nil ***REMOVED***
		return err1
	***REMOVED***
	return nil
***REMOVED***

// fixOSError converts possible platform dependent error into the portable errors in the
// Go os package if possible.
func fixOSError(err error) error ***REMOVED***
	// The os.IsExist, os.IsNotExist, and os.IsPermissions functions are platform
	// dependent, so sending the raw error might break those functions on a different OS.
	// Go defines portable errors for these.
	if os.IsExist(err) ***REMOVED***
		return os.ErrExist
	***REMOVED*** else if os.IsNotExist(err) ***REMOVED***
		return os.ErrNotExist
	***REMOVED*** else if os.IsPermission(err) ***REMOVED***
		return os.ErrPermission
	***REMOVED***
	return err
***REMOVED***

// ReadTarOptions reads from the specified reader and deserializes an archive.TarOptions struct.
func ReadTarOptions(r io.Reader) (*archive.TarOptions, error) ***REMOVED***
	var size uint64
	if err := binary.Read(r, binary.BigEndian, &size); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rawJSON := make([]byte, size)
	if _, err := io.ReadFull(r, rawJSON); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var opts archive.TarOptions
	if err := json.Unmarshal(rawJSON, &opts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &opts, nil
***REMOVED***

// WriteTarOptions serializes a archive.TarOptions struct and writes it to the writer.
func WriteTarOptions(w io.Writer, opts *archive.TarOptions) error ***REMOVED***
	optsBuf, err := json.Marshal(opts)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	optsSize := uint64(len(optsBuf))
	optsSizeBuf := &bytes.Buffer***REMOVED******REMOVED***
	if err := binary.Write(optsSizeBuf, binary.BigEndian, optsSize); err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err := optsSizeBuf.WriteTo(w); err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err := w.Write(optsBuf); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// ReadFileHeader reads from r and returns a deserialized FileHeader
func ReadFileHeader(r io.Reader) (*FileHeader, error) ***REMOVED***
	hdr := &FileHeader***REMOVED******REMOVED***
	if err := binary.Read(r, binary.BigEndian, hdr); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return hdr, nil
***REMOVED***

// WriteFileHeader serializes a FileHeader and writes it to w, along with any extra data
func WriteFileHeader(w io.Writer, hdr *FileHeader, extraData []byte) error ***REMOVED***
	if err := binary.Write(w, binary.BigEndian, hdr); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := w.Write(extraData); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
