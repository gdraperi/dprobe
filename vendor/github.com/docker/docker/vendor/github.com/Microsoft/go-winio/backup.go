// +build windows

package winio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"syscall"
	"unicode/utf16"
)

//sys backupRead(h syscall.Handle, b []byte, bytesRead *uint32, abort bool, processSecurity bool, context *uintptr) (err error) = BackupRead
//sys backupWrite(h syscall.Handle, b []byte, bytesWritten *uint32, abort bool, processSecurity bool, context *uintptr) (err error) = BackupWrite

const (
	BackupData = uint32(iota + 1)
	BackupEaData
	BackupSecurity
	BackupAlternateData
	BackupLink
	BackupPropertyData
	BackupObjectId
	BackupReparseData
	BackupSparseBlock
	BackupTxfsData
)

const (
	StreamSparseAttributes = uint32(8)
)

const (
	WRITE_DAC              = 0x40000
	WRITE_OWNER            = 0x80000
	ACCESS_SYSTEM_SECURITY = 0x1000000
)

// BackupHeader represents a backup stream of a file.
type BackupHeader struct ***REMOVED***
	Id         uint32 // The backup stream ID
	Attributes uint32 // Stream attributes
	Size       int64  // The size of the stream in bytes
	Name       string // The name of the stream (for BackupAlternateData only).
	Offset     int64  // The offset of the stream in the file (for BackupSparseBlock only).
***REMOVED***

type win32StreamId struct ***REMOVED***
	StreamId   uint32
	Attributes uint32
	Size       uint64
	NameSize   uint32
***REMOVED***

// BackupStreamReader reads from a stream produced by the BackupRead Win32 API and produces a series
// of BackupHeader values.
type BackupStreamReader struct ***REMOVED***
	r         io.Reader
	bytesLeft int64
***REMOVED***

// NewBackupStreamReader produces a BackupStreamReader from any io.Reader.
func NewBackupStreamReader(r io.Reader) *BackupStreamReader ***REMOVED***
	return &BackupStreamReader***REMOVED***r, 0***REMOVED***
***REMOVED***

// Next returns the next backup stream and prepares for calls to Read(). It skips the remainder of the current stream if
// it was not completely read.
func (r *BackupStreamReader) Next() (*BackupHeader, error) ***REMOVED***
	if r.bytesLeft > 0 ***REMOVED***
		if s, ok := r.r.(io.Seeker); ok ***REMOVED***
			// Make sure Seek on io.SeekCurrent sometimes succeeds
			// before trying the actual seek.
			if _, err := s.Seek(0, io.SeekCurrent); err == nil ***REMOVED***
				if _, err = s.Seek(r.bytesLeft, io.SeekCurrent); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				r.bytesLeft = 0
			***REMOVED***
		***REMOVED***
		if _, err := io.Copy(ioutil.Discard, r); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	var wsi win32StreamId
	if err := binary.Read(r.r, binary.LittleEndian, &wsi); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	hdr := &BackupHeader***REMOVED***
		Id:         wsi.StreamId,
		Attributes: wsi.Attributes,
		Size:       int64(wsi.Size),
	***REMOVED***
	if wsi.NameSize != 0 ***REMOVED***
		name := make([]uint16, int(wsi.NameSize/2))
		if err := binary.Read(r.r, binary.LittleEndian, name); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		hdr.Name = syscall.UTF16ToString(name)
	***REMOVED***
	if wsi.StreamId == BackupSparseBlock ***REMOVED***
		if err := binary.Read(r.r, binary.LittleEndian, &hdr.Offset); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		hdr.Size -= 8
	***REMOVED***
	r.bytesLeft = hdr.Size
	return hdr, nil
***REMOVED***

// Read reads from the current backup stream.
func (r *BackupStreamReader) Read(b []byte) (int, error) ***REMOVED***
	if r.bytesLeft == 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	if int64(len(b)) > r.bytesLeft ***REMOVED***
		b = b[:r.bytesLeft]
	***REMOVED***
	n, err := r.r.Read(b)
	r.bytesLeft -= int64(n)
	if err == io.EOF ***REMOVED***
		err = io.ErrUnexpectedEOF
	***REMOVED*** else if r.bytesLeft == 0 && err == nil ***REMOVED***
		err = io.EOF
	***REMOVED***
	return n, err
***REMOVED***

// BackupStreamWriter writes a stream compatible with the BackupWrite Win32 API.
type BackupStreamWriter struct ***REMOVED***
	w         io.Writer
	bytesLeft int64
***REMOVED***

// NewBackupStreamWriter produces a BackupStreamWriter on top of an io.Writer.
func NewBackupStreamWriter(w io.Writer) *BackupStreamWriter ***REMOVED***
	return &BackupStreamWriter***REMOVED***w, 0***REMOVED***
***REMOVED***

// WriteHeader writes the next backup stream header and prepares for calls to Write().
func (w *BackupStreamWriter) WriteHeader(hdr *BackupHeader) error ***REMOVED***
	if w.bytesLeft != 0 ***REMOVED***
		return fmt.Errorf("missing %d bytes", w.bytesLeft)
	***REMOVED***
	name := utf16.Encode([]rune(hdr.Name))
	wsi := win32StreamId***REMOVED***
		StreamId:   hdr.Id,
		Attributes: hdr.Attributes,
		Size:       uint64(hdr.Size),
		NameSize:   uint32(len(name) * 2),
	***REMOVED***
	if hdr.Id == BackupSparseBlock ***REMOVED***
		// Include space for the int64 block offset
		wsi.Size += 8
	***REMOVED***
	if err := binary.Write(w.w, binary.LittleEndian, &wsi); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(name) != 0 ***REMOVED***
		if err := binary.Write(w.w, binary.LittleEndian, name); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if hdr.Id == BackupSparseBlock ***REMOVED***
		if err := binary.Write(w.w, binary.LittleEndian, hdr.Offset); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	w.bytesLeft = hdr.Size
	return nil
***REMOVED***

// Write writes to the current backup stream.
func (w *BackupStreamWriter) Write(b []byte) (int, error) ***REMOVED***
	if w.bytesLeft < int64(len(b)) ***REMOVED***
		return 0, fmt.Errorf("too many bytes by %d", int64(len(b))-w.bytesLeft)
	***REMOVED***
	n, err := w.w.Write(b)
	w.bytesLeft -= int64(n)
	return n, err
***REMOVED***

// BackupFileReader provides an io.ReadCloser interface on top of the BackupRead Win32 API.
type BackupFileReader struct ***REMOVED***
	f               *os.File
	includeSecurity bool
	ctx             uintptr
***REMOVED***

// NewBackupFileReader returns a new BackupFileReader from a file handle. If includeSecurity is true,
// Read will attempt to read the security descriptor of the file.
func NewBackupFileReader(f *os.File, includeSecurity bool) *BackupFileReader ***REMOVED***
	r := &BackupFileReader***REMOVED***f, includeSecurity, 0***REMOVED***
	return r
***REMOVED***

// Read reads a backup stream from the file by calling the Win32 API BackupRead().
func (r *BackupFileReader) Read(b []byte) (int, error) ***REMOVED***
	var bytesRead uint32
	err := backupRead(syscall.Handle(r.f.Fd()), b, &bytesRead, false, r.includeSecurity, &r.ctx)
	if err != nil ***REMOVED***
		return 0, &os.PathError***REMOVED***"BackupRead", r.f.Name(), err***REMOVED***
	***REMOVED***
	runtime.KeepAlive(r.f)
	if bytesRead == 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	return int(bytesRead), nil
***REMOVED***

// Close frees Win32 resources associated with the BackupFileReader. It does not close
// the underlying file.
func (r *BackupFileReader) Close() error ***REMOVED***
	if r.ctx != 0 ***REMOVED***
		backupRead(syscall.Handle(r.f.Fd()), nil, nil, true, false, &r.ctx)
		runtime.KeepAlive(r.f)
		r.ctx = 0
	***REMOVED***
	return nil
***REMOVED***

// BackupFileWriter provides an io.WriteCloser interface on top of the BackupWrite Win32 API.
type BackupFileWriter struct ***REMOVED***
	f               *os.File
	includeSecurity bool
	ctx             uintptr
***REMOVED***

// NewBackupFileWriter returns a new BackupFileWriter from a file handle. If includeSecurity is true,
// Write() will attempt to restore the security descriptor from the stream.
func NewBackupFileWriter(f *os.File, includeSecurity bool) *BackupFileWriter ***REMOVED***
	w := &BackupFileWriter***REMOVED***f, includeSecurity, 0***REMOVED***
	return w
***REMOVED***

// Write restores a portion of the file using the provided backup stream.
func (w *BackupFileWriter) Write(b []byte) (int, error) ***REMOVED***
	var bytesWritten uint32
	err := backupWrite(syscall.Handle(w.f.Fd()), b, &bytesWritten, false, w.includeSecurity, &w.ctx)
	if err != nil ***REMOVED***
		return 0, &os.PathError***REMOVED***"BackupWrite", w.f.Name(), err***REMOVED***
	***REMOVED***
	runtime.KeepAlive(w.f)
	if int(bytesWritten) != len(b) ***REMOVED***
		return int(bytesWritten), errors.New("not all bytes could be written")
	***REMOVED***
	return len(b), nil
***REMOVED***

// Close frees Win32 resources associated with the BackupFileWriter. It does not
// close the underlying file.
func (w *BackupFileWriter) Close() error ***REMOVED***
	if w.ctx != 0 ***REMOVED***
		backupWrite(syscall.Handle(w.f.Fd()), nil, nil, true, false, &w.ctx)
		runtime.KeepAlive(w.f)
		w.ctx = 0
	***REMOVED***
	return nil
***REMOVED***

// OpenForBackup opens a file or directory, potentially skipping access checks if the backup
// or restore privileges have been acquired.
//
// If the file opened was a directory, it cannot be used with Readdir().
func OpenForBackup(path string, access uint32, share uint32, createmode uint32) (*os.File, error) ***REMOVED***
	winPath, err := syscall.UTF16FromString(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	h, err := syscall.CreateFile(&winPath[0], access, share, nil, createmode, syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil ***REMOVED***
		err = &os.PathError***REMOVED***Op: "open", Path: path, Err: err***REMOVED***
		return nil, err
	***REMOVED***
	return os.NewFile(uintptr(h), path), nil
***REMOVED***
