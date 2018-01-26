package hcsshim

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Microsoft/go-winio"
)

type baseLayerWriter struct ***REMOVED***
	root         string
	f            *os.File
	bw           *winio.BackupFileWriter
	err          error
	hasUtilityVM bool
	dirInfo      []dirInfo
***REMOVED***

type dirInfo struct ***REMOVED***
	path     string
	fileInfo winio.FileBasicInfo
***REMOVED***

// reapplyDirectoryTimes reapplies directory modification, creation, etc. times
// after processing of the directory tree has completed. The times are expected
// to be ordered such that parent directories come before child directories.
func reapplyDirectoryTimes(dis []dirInfo) error ***REMOVED***
	for i := range dis ***REMOVED***
		di := &dis[len(dis)-i-1] // reverse order: process child directories first
		f, err := winio.OpenForBackup(di.path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, syscall.FILE_SHARE_READ, syscall.OPEN_EXISTING)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = winio.SetFileBasicInfo(f, &di.fileInfo)
		f.Close()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (w *baseLayerWriter) closeCurrentFile() error ***REMOVED***
	if w.f != nil ***REMOVED***
		err := w.bw.Close()
		err2 := w.f.Close()
		w.f = nil
		w.bw = nil
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err2 != nil ***REMOVED***
			return err2
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (w *baseLayerWriter) Add(name string, fileInfo *winio.FileBasicInfo) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			w.err = err
		***REMOVED***
	***REMOVED***()

	err = w.closeCurrentFile()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if filepath.ToSlash(name) == `UtilityVM/Files` ***REMOVED***
		w.hasUtilityVM = true
	***REMOVED***

	path := filepath.Join(w.root, name)
	path, err = makeLongAbsPath(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var f *os.File
	defer func() ***REMOVED***
		if f != nil ***REMOVED***
			f.Close()
		***REMOVED***
	***REMOVED***()

	createmode := uint32(syscall.CREATE_NEW)
	if fileInfo.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 ***REMOVED***
		err := os.Mkdir(path, 0)
		if err != nil && !os.IsExist(err) ***REMOVED***
			return err
		***REMOVED***
		createmode = syscall.OPEN_EXISTING
		if fileInfo.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT == 0 ***REMOVED***
			w.dirInfo = append(w.dirInfo, dirInfo***REMOVED***path, *fileInfo***REMOVED***)
		***REMOVED***
	***REMOVED***

	mode := uint32(syscall.GENERIC_READ | syscall.GENERIC_WRITE | winio.WRITE_DAC | winio.WRITE_OWNER | winio.ACCESS_SYSTEM_SECURITY)
	f, err = winio.OpenForBackup(path, mode, syscall.FILE_SHARE_READ, createmode)
	if err != nil ***REMOVED***
		return makeError(err, "Failed to OpenForBackup", path)
	***REMOVED***

	err = winio.SetFileBasicInfo(f, fileInfo)
	if err != nil ***REMOVED***
		return makeError(err, "Failed to SetFileBasicInfo", path)
	***REMOVED***

	w.f = f
	w.bw = winio.NewBackupFileWriter(f, true)
	f = nil
	return nil
***REMOVED***

func (w *baseLayerWriter) AddLink(name string, target string) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			w.err = err
		***REMOVED***
	***REMOVED***()

	err = w.closeCurrentFile()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	linkpath, err := makeLongAbsPath(filepath.Join(w.root, name))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	linktarget, err := makeLongAbsPath(filepath.Join(w.root, target))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return os.Link(linktarget, linkpath)
***REMOVED***

func (w *baseLayerWriter) Remove(name string) error ***REMOVED***
	return errors.New("base layer cannot have tombstones")
***REMOVED***

func (w *baseLayerWriter) Write(b []byte) (int, error) ***REMOVED***
	n, err := w.bw.Write(b)
	if err != nil ***REMOVED***
		w.err = err
	***REMOVED***
	return n, err
***REMOVED***

func (w *baseLayerWriter) Close() error ***REMOVED***
	err := w.closeCurrentFile()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if w.err == nil ***REMOVED***
		// Restore the file times of all the directories, since they may have
		// been modified by creating child directories.
		err = reapplyDirectoryTimes(w.dirInfo)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = ProcessBaseLayer(w.root)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if w.hasUtilityVM ***REMOVED***
			err = ProcessUtilityVMImage(filepath.Join(w.root, "UtilityVM"))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return w.err
***REMOVED***
