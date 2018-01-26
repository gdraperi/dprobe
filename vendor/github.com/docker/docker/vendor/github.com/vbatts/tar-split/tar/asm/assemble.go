package asm

import (
	"bytes"
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"sync"

	"github.com/vbatts/tar-split/tar/storage"
)

// NewOutputTarStream returns an io.ReadCloser that is an assembled tar archive
// stream.
//
// It takes a storage.FileGetter, for mapping the file payloads that are to be read in,
// and a storage.Unpacker, which has access to the rawbytes and file order
// metadata. With the combination of these two items, a precise assembled Tar
// archive is possible.
func NewOutputTarStream(fg storage.FileGetter, up storage.Unpacker) io.ReadCloser ***REMOVED***
	// ... Since these are interfaces, this is possible, so let's not have a nil pointer
	if fg == nil || up == nil ***REMOVED***
		return nil
	***REMOVED***
	pr, pw := io.Pipe()
	go func() ***REMOVED***
		err := WriteOutputTarStream(fg, up, pw)
		if err != nil ***REMOVED***
			pw.CloseWithError(err)
		***REMOVED*** else ***REMOVED***
			pw.Close()
		***REMOVED***
	***REMOVED***()
	return pr
***REMOVED***

// WriteOutputTarStream writes assembled tar archive to a writer.
func WriteOutputTarStream(fg storage.FileGetter, up storage.Unpacker, w io.Writer) error ***REMOVED***
	// ... Since these are interfaces, this is possible, so let's not have a nil pointer
	if fg == nil || up == nil ***REMOVED***
		return nil
	***REMOVED***
	var copyBuffer []byte
	var crcHash hash.Hash
	var crcSum []byte
	var multiWriter io.Writer
	for ***REMOVED***
		entry, err := up.Next()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				return nil
			***REMOVED***
			return err
		***REMOVED***
		switch entry.Type ***REMOVED***
		case storage.SegmentType:
			if _, err := w.Write(entry.Payload); err != nil ***REMOVED***
				return err
			***REMOVED***
		case storage.FileType:
			if entry.Size == 0 ***REMOVED***
				continue
			***REMOVED***
			fh, err := fg.Get(entry.GetName())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if crcHash == nil ***REMOVED***
				crcHash = crc64.New(storage.CRCTable)
				crcSum = make([]byte, 8)
				multiWriter = io.MultiWriter(w, crcHash)
				copyBuffer = byteBufferPool.Get().([]byte)
				defer byteBufferPool.Put(copyBuffer)
			***REMOVED*** else ***REMOVED***
				crcHash.Reset()
			***REMOVED***

			if _, err := copyWithBuffer(multiWriter, fh, copyBuffer); err != nil ***REMOVED***
				fh.Close()
				return err
			***REMOVED***

			if !bytes.Equal(crcHash.Sum(crcSum[:0]), entry.Payload) ***REMOVED***
				// I would rather this be a comparable ErrInvalidChecksum or such,
				// but since it's coming through the PipeReader, the context of
				// _which_ file would be lost...
				fh.Close()
				return fmt.Errorf("file integrity checksum failed for %q", entry.GetName())
			***REMOVED***
			fh.Close()
		***REMOVED***
	***REMOVED***
***REMOVED***

var byteBufferPool = &sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return make([]byte, 32*1024)
	***REMOVED***,
***REMOVED***

// copyWithBuffer is taken from stdlib io.Copy implementation
// https://github.com/golang/go/blob/go1.5.1/src/io/io.go#L367
func copyWithBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) ***REMOVED***
	for ***REMOVED***
		nr, er := src.Read(buf)
		if nr > 0 ***REMOVED***
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 ***REMOVED***
				written += int64(nw)
			***REMOVED***
			if ew != nil ***REMOVED***
				err = ew
				break
			***REMOVED***
			if nr != nw ***REMOVED***
				err = io.ErrShortWrite
				break
			***REMOVED***
		***REMOVED***
		if er == io.EOF ***REMOVED***
			break
		***REMOVED***
		if er != nil ***REMOVED***
			err = er
			break
		***REMOVED***
	***REMOVED***
	return written, err
***REMOVED***
