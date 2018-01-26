package asm

import (
	"io"

	"github.com/vbatts/tar-split/archive/tar"
	"github.com/vbatts/tar-split/tar/storage"
)

// NewInputTarStream wraps the Reader stream of a tar archive and provides a
// Reader stream of the same.
//
// In the middle it will pack the segments and file metadata to storage.Packer
// `p`.
//
// The the storage.FilePutter is where payload of files in the stream are
// stashed. If this stashing is not needed, you can provide a nil
// storage.FilePutter. Since the checksumming is still needed, then a default
// of NewDiscardFilePutter will be used internally
func NewInputTarStream(r io.Reader, p storage.Packer, fp storage.FilePutter) (io.Reader, error) ***REMOVED***
	// What to do here... folks will want their own access to the Reader that is
	// their tar archive stream, but we'll need that same stream to use our
	// forked 'archive/tar'.
	// Perhaps do an io.TeeReader that hands back an io.Reader for them to read
	// from, and we'll MITM the stream to store metadata.
	// We'll need a storage.FilePutter too ...

	// Another concern, whether to do any storage.FilePutter operations, such that we
	// don't extract any amount of the archive. But then again, we're not making
	// files/directories, hardlinks, etc. Just writing the io to the storage.FilePutter.
	// Perhaps we have a DiscardFilePutter that is a bit bucket.

	// we'll return the pipe reader, since TeeReader does not buffer and will
	// only read what the outputRdr Read's. Since Tar archives have padding on
	// the end, we want to be the one reading the padding, even if the user's
	// `archive/tar` doesn't care.
	pR, pW := io.Pipe()
	outputRdr := io.TeeReader(r, pW)

	// we need a putter that will generate the crc64 sums of file payloads
	if fp == nil ***REMOVED***
		fp = storage.NewDiscardFilePutter()
	***REMOVED***

	go func() ***REMOVED***
		tr := tar.NewReader(outputRdr)
		tr.RawAccounting = true
		for ***REMOVED***
			hdr, err := tr.Next()
			if err != nil ***REMOVED***
				if err != io.EOF ***REMOVED***
					pW.CloseWithError(err)
					return
				***REMOVED***
				// even when an EOF is reached, there is often 1024 null bytes on
				// the end of an archive. Collect them too.
				if b := tr.RawBytes(); len(b) > 0 ***REMOVED***
					_, err := p.AddEntry(storage.Entry***REMOVED***
						Type:    storage.SegmentType,
						Payload: b,
					***REMOVED***)
					if err != nil ***REMOVED***
						pW.CloseWithError(err)
						return
					***REMOVED***
				***REMOVED***
				break // not return. We need the end of the reader.
			***REMOVED***
			if hdr == nil ***REMOVED***
				break // not return. We need the end of the reader.
			***REMOVED***

			if b := tr.RawBytes(); len(b) > 0 ***REMOVED***
				_, err := p.AddEntry(storage.Entry***REMOVED***
					Type:    storage.SegmentType,
					Payload: b,
				***REMOVED***)
				if err != nil ***REMOVED***
					pW.CloseWithError(err)
					return
				***REMOVED***
			***REMOVED***

			var csum []byte
			if hdr.Size > 0 ***REMOVED***
				var err error
				_, csum, err = fp.Put(hdr.Name, tr)
				if err != nil ***REMOVED***
					pW.CloseWithError(err)
					return
				***REMOVED***
			***REMOVED***

			entry := storage.Entry***REMOVED***
				Type:    storage.FileType,
				Size:    hdr.Size,
				Payload: csum,
			***REMOVED***
			// For proper marshalling of non-utf8 characters
			entry.SetName(hdr.Name)

			// File entries added, regardless of size
			_, err = p.AddEntry(entry)
			if err != nil ***REMOVED***
				pW.CloseWithError(err)
				return
			***REMOVED***

			if b := tr.RawBytes(); len(b) > 0 ***REMOVED***
				_, err = p.AddEntry(storage.Entry***REMOVED***
					Type:    storage.SegmentType,
					Payload: b,
				***REMOVED***)
				if err != nil ***REMOVED***
					pW.CloseWithError(err)
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// It is allowable, and not uncommon that there is further padding on
		// the end of an archive, apart from the expected 1024 null bytes. We
		// do this in chunks rather than in one go to avoid cases where a
		// maliciously crafted tar file tries to trick us into reading many GBs
		// into memory.
		const paddingChunkSize = 1024 * 1024
		var paddingChunk [paddingChunkSize]byte
		for ***REMOVED***
			var isEOF bool
			n, err := outputRdr.Read(paddingChunk[:])
			if err != nil ***REMOVED***
				if err != io.EOF ***REMOVED***
					pW.CloseWithError(err)
					return
				***REMOVED***
				isEOF = true
			***REMOVED***
			_, err = p.AddEntry(storage.Entry***REMOVED***
				Type:    storage.SegmentType,
				Payload: paddingChunk[:n],
			***REMOVED***)
			if err != nil ***REMOVED***
				pW.CloseWithError(err)
				return
			***REMOVED***
			if isEOF ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		pW.Close()
	***REMOVED***()

	return pR, nil
***REMOVED***
