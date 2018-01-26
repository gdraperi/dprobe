// +build windows

package backuptar

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Microsoft/go-winio"
	"github.com/Microsoft/go-winio/archive/tar" // until archive/tar supports pax extensions in its interface
)

const (
	c_ISUID  = 04000   // Set uid
	c_ISGID  = 02000   // Set gid
	c_ISVTX  = 01000   // Save text (sticky bit)
	c_ISDIR  = 040000  // Directory
	c_ISFIFO = 010000  // FIFO
	c_ISREG  = 0100000 // Regular file
	c_ISLNK  = 0120000 // Symbolic link
	c_ISBLK  = 060000  // Block special file
	c_ISCHR  = 020000  // Character special file
	c_ISSOCK = 0140000 // Socket
)

const (
	hdrFileAttributes        = "fileattr"
	hdrSecurityDescriptor    = "sd"
	hdrRawSecurityDescriptor = "rawsd"
	hdrMountPoint            = "mountpoint"
	hdrEaPrefix              = "xattr."
)

func writeZeroes(w io.Writer, count int64) error ***REMOVED***
	buf := make([]byte, 8192)
	c := len(buf)
	for i := int64(0); i < count; i += int64(c) ***REMOVED***
		if int64(c) > count-i ***REMOVED***
			c = int(count - i)
		***REMOVED***
		_, err := w.Write(buf[:c])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func copySparse(t *tar.Writer, br *winio.BackupStreamReader) error ***REMOVED***
	curOffset := int64(0)
	for ***REMOVED***
		bhdr, err := br.Next()
		if err == io.EOF ***REMOVED***
			err = io.ErrUnexpectedEOF
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if bhdr.Id != winio.BackupSparseBlock ***REMOVED***
			return fmt.Errorf("unexpected stream %d", bhdr.Id)
		***REMOVED***

		// archive/tar does not support writing sparse files
		// so just write zeroes to catch up to the current offset.
		err = writeZeroes(t, bhdr.Offset-curOffset)
		if bhdr.Size == 0 ***REMOVED***
			break
		***REMOVED***
		n, err := io.Copy(t, br)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		curOffset = bhdr.Offset + n
	***REMOVED***
	return nil
***REMOVED***

// BasicInfoHeader creates a tar header from basic file information.
func BasicInfoHeader(name string, size int64, fileInfo *winio.FileBasicInfo) *tar.Header ***REMOVED***
	hdr := &tar.Header***REMOVED***
		Name:         filepath.ToSlash(name),
		Size:         size,
		Typeflag:     tar.TypeReg,
		ModTime:      time.Unix(0, fileInfo.LastWriteTime.Nanoseconds()),
		ChangeTime:   time.Unix(0, fileInfo.ChangeTime.Nanoseconds()),
		AccessTime:   time.Unix(0, fileInfo.LastAccessTime.Nanoseconds()),
		CreationTime: time.Unix(0, fileInfo.CreationTime.Nanoseconds()),
		Winheaders:   make(map[string]string),
	***REMOVED***
	hdr.Winheaders[hdrFileAttributes] = fmt.Sprintf("%d", fileInfo.FileAttributes)

	if (fileInfo.FileAttributes & syscall.FILE_ATTRIBUTE_DIRECTORY) != 0 ***REMOVED***
		hdr.Mode |= c_ISDIR
		hdr.Size = 0
		hdr.Typeflag = tar.TypeDir
	***REMOVED***
	return hdr
***REMOVED***

// WriteTarFileFromBackupStream writes a file to a tar writer using data from a Win32 backup stream.
//
// This encodes Win32 metadata as tar pax vendor extensions starting with MSWINDOWS.
//
// The additional Win32 metadata is:
//
// MSWINDOWS.fileattr: The Win32 file attributes, as a decimal value
//
// MSWINDOWS.rawsd: The Win32 security descriptor, in raw binary format
//
// MSWINDOWS.mountpoint: If present, this is a mount point and not a symlink, even though the type is '2' (symlink)
func WriteTarFileFromBackupStream(t *tar.Writer, r io.Reader, name string, size int64, fileInfo *winio.FileBasicInfo) error ***REMOVED***
	name = filepath.ToSlash(name)
	hdr := BasicInfoHeader(name, size, fileInfo)

	// If r can be seeked, then this function is two-pass: pass 1 collects the
	// tar header data, and pass 2 copies the data stream. If r cannot be
	// seeked, then some header data (in particular EAs) will be silently lost.
	var (
		restartPos int64
		err        error
	)
	sr, readTwice := r.(io.Seeker)
	if readTwice ***REMOVED***
		if restartPos, err = sr.Seek(0, io.SeekCurrent); err != nil ***REMOVED***
			readTwice = false
		***REMOVED***
	***REMOVED***

	br := winio.NewBackupStreamReader(r)
	var dataHdr *winio.BackupHeader
	for dataHdr == nil ***REMOVED***
		bhdr, err := br.Next()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch bhdr.Id ***REMOVED***
		case winio.BackupData:
			hdr.Mode |= c_ISREG
			if !readTwice ***REMOVED***
				dataHdr = bhdr
			***REMOVED***
		case winio.BackupSecurity:
			sd, err := ioutil.ReadAll(br)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hdr.Winheaders[hdrRawSecurityDescriptor] = base64.StdEncoding.EncodeToString(sd)

		case winio.BackupReparseData:
			hdr.Mode |= c_ISLNK
			hdr.Typeflag = tar.TypeSymlink
			reparseBuffer, err := ioutil.ReadAll(br)
			rp, err := winio.DecodeReparsePoint(reparseBuffer)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if rp.IsMountPoint ***REMOVED***
				hdr.Winheaders[hdrMountPoint] = "1"
			***REMOVED***
			hdr.Linkname = rp.Target

		case winio.BackupEaData:
			eab, err := ioutil.ReadAll(br)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			eas, err := winio.DecodeExtendedAttributes(eab)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			for _, ea := range eas ***REMOVED***
				// Use base64 encoding for the binary value. Note that there
				// is no way to encode the EA's flags, since their use doesn't
				// make any sense for persisted EAs.
				hdr.Winheaders[hdrEaPrefix+ea.Name] = base64.StdEncoding.EncodeToString(ea.Value)
			***REMOVED***

		case winio.BackupAlternateData, winio.BackupLink, winio.BackupPropertyData, winio.BackupObjectId, winio.BackupTxfsData:
			// ignore these streams
		default:
			return fmt.Errorf("%s: unknown stream ID %d", name, bhdr.Id)
		***REMOVED***
	***REMOVED***

	err = t.WriteHeader(hdr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if readTwice ***REMOVED***
		// Get back to the data stream.
		if _, err = sr.Seek(restartPos, io.SeekStart); err != nil ***REMOVED***
			return err
		***REMOVED***
		for dataHdr == nil ***REMOVED***
			bhdr, err := br.Next()
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if bhdr.Id == winio.BackupData ***REMOVED***
				dataHdr = bhdr
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if dataHdr != nil ***REMOVED***
		// A data stream was found. Copy the data.
		if (dataHdr.Attributes & winio.StreamSparseAttributes) == 0 ***REMOVED***
			if size != dataHdr.Size ***REMOVED***
				return fmt.Errorf("%s: mismatch between file size %d and header size %d", name, size, dataHdr.Size)
			***REMOVED***
			_, err = io.Copy(t, br)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = copySparse(t, br)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Look for streams after the data stream. The only ones we handle are alternate data streams.
	// Other streams may have metadata that could be serialized, but the tar header has already
	// been written. In practice, this means that we don't get EA or TXF metadata.
	for ***REMOVED***
		bhdr, err := br.Next()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch bhdr.Id ***REMOVED***
		case winio.BackupAlternateData:
			altName := bhdr.Name
			if strings.HasSuffix(altName, ":$DATA") ***REMOVED***
				altName = altName[:len(altName)-len(":$DATA")]
			***REMOVED***
			if (bhdr.Attributes & winio.StreamSparseAttributes) == 0 ***REMOVED***
				hdr = &tar.Header***REMOVED***
					Name:       name + altName,
					Mode:       hdr.Mode,
					Typeflag:   tar.TypeReg,
					Size:       bhdr.Size,
					ModTime:    hdr.ModTime,
					AccessTime: hdr.AccessTime,
					ChangeTime: hdr.ChangeTime,
				***REMOVED***
				err = t.WriteHeader(hdr)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				_, err = io.Copy(t, br)
				if err != nil ***REMOVED***
					return err
				***REMOVED***

			***REMOVED*** else ***REMOVED***
				// Unsupported for now, since the size of the alternate stream is not present
				// in the backup stream until after the data has been read.
				return errors.New("tar of sparse alternate data streams is unsupported")
			***REMOVED***
		case winio.BackupEaData, winio.BackupLink, winio.BackupPropertyData, winio.BackupObjectId, winio.BackupTxfsData:
			// ignore these streams
		default:
			return fmt.Errorf("%s: unknown stream ID %d after data", name, bhdr.Id)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// FileInfoFromHeader retrieves basic Win32 file information from a tar header, using the additional metadata written by
// WriteTarFileFromBackupStream.
func FileInfoFromHeader(hdr *tar.Header) (name string, size int64, fileInfo *winio.FileBasicInfo, err error) ***REMOVED***
	name = hdr.Name
	if hdr.Typeflag == tar.TypeReg || hdr.Typeflag == tar.TypeRegA ***REMOVED***
		size = hdr.Size
	***REMOVED***
	fileInfo = &winio.FileBasicInfo***REMOVED***
		LastAccessTime: syscall.NsecToFiletime(hdr.AccessTime.UnixNano()),
		LastWriteTime:  syscall.NsecToFiletime(hdr.ModTime.UnixNano()),
		ChangeTime:     syscall.NsecToFiletime(hdr.ChangeTime.UnixNano()),
		CreationTime:   syscall.NsecToFiletime(hdr.CreationTime.UnixNano()),
	***REMOVED***
	if attrStr, ok := hdr.Winheaders[hdrFileAttributes]; ok ***REMOVED***
		attr, err := strconv.ParseUint(attrStr, 10, 32)
		if err != nil ***REMOVED***
			return "", 0, nil, err
		***REMOVED***
		fileInfo.FileAttributes = uintptr(attr)
	***REMOVED*** else ***REMOVED***
		if hdr.Typeflag == tar.TypeDir ***REMOVED***
			fileInfo.FileAttributes |= syscall.FILE_ATTRIBUTE_DIRECTORY
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// WriteBackupStreamFromTarFile writes a Win32 backup stream from the current tar file. Since this function may process multiple
// tar file entries in order to collect all the alternate data streams for the file, it returns the next
// tar file that was not processed, or io.EOF is there are no more.
func WriteBackupStreamFromTarFile(w io.Writer, t *tar.Reader, hdr *tar.Header) (*tar.Header, error) ***REMOVED***
	bw := winio.NewBackupStreamWriter(w)
	var sd []byte
	var err error
	// Maintaining old SDDL-based behavior for backward compatibility.  All new tar headers written
	// by this library will have raw binary for the security descriptor.
	if sddl, ok := hdr.Winheaders[hdrSecurityDescriptor]; ok ***REMOVED***
		sd, err = winio.SddlToSecurityDescriptor(sddl)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if sdraw, ok := hdr.Winheaders[hdrRawSecurityDescriptor]; ok ***REMOVED***
		sd, err = base64.StdEncoding.DecodeString(sdraw)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if len(sd) != 0 ***REMOVED***
		bhdr := winio.BackupHeader***REMOVED***
			Id:   winio.BackupSecurity,
			Size: int64(len(sd)),
		***REMOVED***
		err := bw.WriteHeader(&bhdr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = bw.Write(sd)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	var eas []winio.ExtendedAttribute
	for k, v := range hdr.Winheaders ***REMOVED***
		if !strings.HasPrefix(k, hdrEaPrefix) ***REMOVED***
			continue
		***REMOVED***
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		eas = append(eas, winio.ExtendedAttribute***REMOVED***
			Name:  k[len(hdrEaPrefix):],
			Value: data,
		***REMOVED***)
	***REMOVED***
	if len(eas) != 0 ***REMOVED***
		eadata, err := winio.EncodeExtendedAttributes(eas)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		bhdr := winio.BackupHeader***REMOVED***
			Id:   winio.BackupEaData,
			Size: int64(len(eadata)),
		***REMOVED***
		err = bw.WriteHeader(&bhdr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = bw.Write(eadata)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if hdr.Typeflag == tar.TypeSymlink ***REMOVED***
		_, isMountPoint := hdr.Winheaders[hdrMountPoint]
		rp := winio.ReparsePoint***REMOVED***
			Target:       filepath.FromSlash(hdr.Linkname),
			IsMountPoint: isMountPoint,
		***REMOVED***
		reparse := winio.EncodeReparsePoint(&rp)
		bhdr := winio.BackupHeader***REMOVED***
			Id:   winio.BackupReparseData,
			Size: int64(len(reparse)),
		***REMOVED***
		err := bw.WriteHeader(&bhdr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = bw.Write(reparse)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if hdr.Typeflag == tar.TypeReg || hdr.Typeflag == tar.TypeRegA ***REMOVED***
		bhdr := winio.BackupHeader***REMOVED***
			Id:   winio.BackupData,
			Size: hdr.Size,
		***REMOVED***
		err := bw.WriteHeader(&bhdr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = io.Copy(bw, t)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	// Copy all the alternate data streams and return the next non-ADS header.
	for ***REMOVED***
		ahdr, err := t.Next()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if ahdr.Typeflag != tar.TypeReg || !strings.HasPrefix(ahdr.Name, hdr.Name+":") ***REMOVED***
			return ahdr, nil
		***REMOVED***
		bhdr := winio.BackupHeader***REMOVED***
			Id:   winio.BackupAlternateData,
			Size: ahdr.Size,
			Name: ahdr.Name[len(hdr.Name):] + ":$DATA",
		***REMOVED***
		err = bw.WriteHeader(&bhdr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = io.Copy(bw, t)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
***REMOVED***
