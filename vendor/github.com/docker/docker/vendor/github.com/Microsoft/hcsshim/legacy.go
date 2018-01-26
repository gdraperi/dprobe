package hcsshim

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Microsoft/go-winio"
)

var errorIterationCanceled = errors.New("")

var mutatedUtilityVMFiles = map[string]bool***REMOVED***
	`EFI\Microsoft\Boot\BCD`:      true,
	`EFI\Microsoft\Boot\BCD.LOG`:  true,
	`EFI\Microsoft\Boot\BCD.LOG1`: true,
	`EFI\Microsoft\Boot\BCD.LOG2`: true,
***REMOVED***

const (
	filesPath          = `Files`
	hivesPath          = `Hives`
	utilityVMPath      = `UtilityVM`
	utilityVMFilesPath = `UtilityVM\Files`
)

func openFileOrDir(path string, mode uint32, createDisposition uint32) (file *os.File, err error) ***REMOVED***
	return winio.OpenForBackup(path, mode, syscall.FILE_SHARE_READ, createDisposition)
***REMOVED***

func makeLongAbsPath(path string) (string, error) ***REMOVED***
	if strings.HasPrefix(path, `\\?\`) || strings.HasPrefix(path, `\\.\`) ***REMOVED***
		return path, nil
	***REMOVED***
	if !filepath.IsAbs(path) ***REMOVED***
		absPath, err := filepath.Abs(path)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		path = absPath
	***REMOVED***
	if strings.HasPrefix(path, `\\`) ***REMOVED***
		return `\\?\UNC\` + path[2:], nil
	***REMOVED***
	return `\\?\` + path, nil
***REMOVED***

func hasPathPrefix(p, prefix string) bool ***REMOVED***
	return strings.HasPrefix(p, prefix) && len(p) > len(prefix) && p[len(prefix)] == '\\'
***REMOVED***

type fileEntry struct ***REMOVED***
	path string
	fi   os.FileInfo
	err  error
***REMOVED***

type legacyLayerReader struct ***REMOVED***
	root         string
	result       chan *fileEntry
	proceed      chan bool
	currentFile  *os.File
	backupReader *winio.BackupFileReader
***REMOVED***

// newLegacyLayerReader returns a new LayerReader that can read the Windows
// container layer transport format from disk.
func newLegacyLayerReader(root string) *legacyLayerReader ***REMOVED***
	r := &legacyLayerReader***REMOVED***
		root:    root,
		result:  make(chan *fileEntry),
		proceed: make(chan bool),
	***REMOVED***
	go r.walk()
	return r
***REMOVED***

func readTombstones(path string) (map[string]([]string), error) ***REMOVED***
	tf, err := os.Open(filepath.Join(path, "tombstones.txt"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer tf.Close()
	s := bufio.NewScanner(tf)
	if !s.Scan() || s.Text() != "\xef\xbb\xbfVersion 1.0" ***REMOVED***
		return nil, errors.New("Invalid tombstones file")
	***REMOVED***

	ts := make(map[string]([]string))
	for s.Scan() ***REMOVED***
		t := filepath.Join(filesPath, s.Text()[1:]) // skip leading `\`
		dir := filepath.Dir(t)
		ts[dir] = append(ts[dir], t)
	***REMOVED***
	if err = s.Err(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return ts, nil
***REMOVED***

func (r *legacyLayerReader) walkUntilCancelled() error ***REMOVED***
	root, err := makeLongAbsPath(r.root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.root = root
	ts, err := readTombstones(r.root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = filepath.Walk(r.root, func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Indirect fix for https://github.com/moby/moby/issues/32838#issuecomment-343610048.
		// Handle failure from what may be a golang bug in the conversion of
		// UTF16 to UTF8 in files which are left in the recycle bin. Os.Lstat
		// which is called by filepath.Walk will fail when a filename contains
		// unicode characters. Skip the recycle bin regardless which is goodness.
		if strings.HasPrefix(path, filepath.Join(r.root, `Files\$Recycle.Bin`)) ***REMOVED***
			return filepath.SkipDir
		***REMOVED***

		if path == r.root || path == filepath.Join(r.root, "tombstones.txt") || strings.HasSuffix(path, ".$wcidirs$") ***REMOVED***
			return nil
		***REMOVED***

		r.result <- &fileEntry***REMOVED***path, info, nil***REMOVED***
		if !<-r.proceed ***REMOVED***
			return errorIterationCanceled
		***REMOVED***

		// List all the tombstones.
		if info.IsDir() ***REMOVED***
			relPath, err := filepath.Rel(r.root, path)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if dts, ok := ts[relPath]; ok ***REMOVED***
				for _, t := range dts ***REMOVED***
					r.result <- &fileEntry***REMOVED***filepath.Join(r.root, t), nil, nil***REMOVED***
					if !<-r.proceed ***REMOVED***
						return errorIterationCanceled
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err == errorIterationCanceled ***REMOVED***
		return nil
	***REMOVED***
	if err == nil ***REMOVED***
		return io.EOF
	***REMOVED***
	return err
***REMOVED***

func (r *legacyLayerReader) walk() ***REMOVED***
	defer close(r.result)
	if !<-r.proceed ***REMOVED***
		return
	***REMOVED***

	err := r.walkUntilCancelled()
	if err != nil ***REMOVED***
		for ***REMOVED***
			r.result <- &fileEntry***REMOVED***err: err***REMOVED***
			if !<-r.proceed ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *legacyLayerReader) reset() ***REMOVED***
	if r.backupReader != nil ***REMOVED***
		r.backupReader.Close()
		r.backupReader = nil
	***REMOVED***
	if r.currentFile != nil ***REMOVED***
		r.currentFile.Close()
		r.currentFile = nil
	***REMOVED***
***REMOVED***

func findBackupStreamSize(r io.Reader) (int64, error) ***REMOVED***
	br := winio.NewBackupStreamReader(r)
	for ***REMOVED***
		hdr, err := br.Next()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				err = nil
			***REMOVED***
			return 0, err
		***REMOVED***
		if hdr.Id == winio.BackupData ***REMOVED***
			return hdr.Size, nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *legacyLayerReader) Next() (path string, size int64, fileInfo *winio.FileBasicInfo, err error) ***REMOVED***
	r.reset()
	r.proceed <- true
	fe := <-r.result
	if fe == nil ***REMOVED***
		err = errors.New("LegacyLayerReader closed")
		return
	***REMOVED***
	if fe.err != nil ***REMOVED***
		err = fe.err
		return
	***REMOVED***

	path, err = filepath.Rel(r.root, fe.path)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if fe.fi == nil ***REMOVED***
		// This is a tombstone. Return a nil fileInfo.
		return
	***REMOVED***

	if fe.fi.IsDir() && hasPathPrefix(path, filesPath) ***REMOVED***
		fe.path += ".$wcidirs$"
	***REMOVED***

	f, err := openFileOrDir(fe.path, syscall.GENERIC_READ, syscall.OPEN_EXISTING)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer func() ***REMOVED***
		if f != nil ***REMOVED***
			f.Close()
		***REMOVED***
	***REMOVED***()

	fileInfo, err = winio.GetFileBasicInfo(f)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if !hasPathPrefix(path, filesPath) ***REMOVED***
		size = fe.fi.Size()
		r.backupReader = winio.NewBackupFileReader(f, false)
		if path == hivesPath || path == filesPath ***REMOVED***
			// The Hives directory has a non-deterministic file time because of the
			// nature of the import process. Use the times from System_Delta.
			var g *os.File
			g, err = os.Open(filepath.Join(r.root, hivesPath, `System_Delta`))
			if err != nil ***REMOVED***
				return
			***REMOVED***
			attr := fileInfo.FileAttributes
			fileInfo, err = winio.GetFileBasicInfo(g)
			g.Close()
			if err != nil ***REMOVED***
				return
			***REMOVED***
			fileInfo.FileAttributes = attr
		***REMOVED***

		// The creation time and access time get reset for files outside of the Files path.
		fileInfo.CreationTime = fileInfo.LastWriteTime
		fileInfo.LastAccessTime = fileInfo.LastWriteTime

	***REMOVED*** else ***REMOVED***
		// The file attributes are written before the backup stream.
		var attr uint32
		err = binary.Read(f, binary.LittleEndian, &attr)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		fileInfo.FileAttributes = uintptr(attr)
		beginning := int64(4)

		// Find the accurate file size.
		if !fe.fi.IsDir() ***REMOVED***
			size, err = findBackupStreamSize(f)
			if err != nil ***REMOVED***
				err = &os.PathError***REMOVED***Op: "findBackupStreamSize", Path: fe.path, Err: err***REMOVED***
				return
			***REMOVED***
		***REMOVED***

		// Return back to the beginning of the backup stream.
		_, err = f.Seek(beginning, 0)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	r.currentFile = f
	f = nil
	return
***REMOVED***

func (r *legacyLayerReader) Read(b []byte) (int, error) ***REMOVED***
	if r.backupReader == nil ***REMOVED***
		if r.currentFile == nil ***REMOVED***
			return 0, io.EOF
		***REMOVED***
		return r.currentFile.Read(b)
	***REMOVED***
	return r.backupReader.Read(b)
***REMOVED***

func (r *legacyLayerReader) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	if r.backupReader == nil ***REMOVED***
		if r.currentFile == nil ***REMOVED***
			return 0, errors.New("no current file")
		***REMOVED***
		return r.currentFile.Seek(offset, whence)
	***REMOVED***
	return 0, errors.New("seek not supported on this stream")
***REMOVED***

func (r *legacyLayerReader) Close() error ***REMOVED***
	r.proceed <- false
	<-r.result
	r.reset()
	return nil
***REMOVED***

type pendingLink struct ***REMOVED***
	Path, Target string
***REMOVED***

type legacyLayerWriter struct ***REMOVED***
	root         string
	parentRoots  []string
	destRoot     string
	currentFile  *os.File
	backupWriter *winio.BackupFileWriter
	tombstones   []string
	pathFixed    bool
	HasUtilityVM bool
	uvmDi        []dirInfo
	addedFiles   map[string]bool
	PendingLinks []pendingLink
***REMOVED***

// newLegacyLayerWriter returns a LayerWriter that can write the contaler layer
// transport format to disk.
func newLegacyLayerWriter(root string, parentRoots []string, destRoot string) *legacyLayerWriter ***REMOVED***
	return &legacyLayerWriter***REMOVED***
		root:        root,
		parentRoots: parentRoots,
		destRoot:    destRoot,
		addedFiles:  make(map[string]bool),
	***REMOVED***
***REMOVED***

func (w *legacyLayerWriter) init() error ***REMOVED***
	if !w.pathFixed ***REMOVED***
		path, err := makeLongAbsPath(w.root)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for i, p := range w.parentRoots ***REMOVED***
			w.parentRoots[i], err = makeLongAbsPath(p)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		destPath, err := makeLongAbsPath(w.destRoot)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		w.root = path
		w.destRoot = destPath
		w.pathFixed = true
	***REMOVED***
	return nil
***REMOVED***

func (w *legacyLayerWriter) initUtilityVM() error ***REMOVED***
	if !w.HasUtilityVM ***REMOVED***
		err := os.Mkdir(filepath.Join(w.destRoot, utilityVMPath), 0)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Server 2016 does not support multiple layers for the utility VM, so
		// clone the utility VM from the parent layer into this layer. Use hard
		// links to avoid unnecessary copying, since most of the files are
		// immutable.
		err = cloneTree(filepath.Join(w.parentRoots[0], utilityVMFilesPath), filepath.Join(w.destRoot, utilityVMFilesPath), mutatedUtilityVMFiles)
		if err != nil ***REMOVED***
			return fmt.Errorf("cloning the parent utility VM image failed: %s", err)
		***REMOVED***
		w.HasUtilityVM = true
	***REMOVED***
	return nil
***REMOVED***

func (w *legacyLayerWriter) reset() ***REMOVED***
	if w.backupWriter != nil ***REMOVED***
		w.backupWriter.Close()
		w.backupWriter = nil
	***REMOVED***
	if w.currentFile != nil ***REMOVED***
		w.currentFile.Close()
		w.currentFile = nil
	***REMOVED***
***REMOVED***

// copyFileWithMetadata copies a file using the backup/restore APIs in order to preserve metadata
func copyFileWithMetadata(srcPath, destPath string, isDir bool) (fileInfo *winio.FileBasicInfo, err error) ***REMOVED***
	createDisposition := uint32(syscall.CREATE_NEW)
	if isDir ***REMOVED***
		err = os.Mkdir(destPath, 0)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		createDisposition = syscall.OPEN_EXISTING
	***REMOVED***

	src, err := openFileOrDir(srcPath, syscall.GENERIC_READ|winio.ACCESS_SYSTEM_SECURITY, syscall.OPEN_EXISTING)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer src.Close()
	srcr := winio.NewBackupFileReader(src, true)
	defer srcr.Close()

	fileInfo, err = winio.GetFileBasicInfo(src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	dest, err := openFileOrDir(destPath, syscall.GENERIC_READ|syscall.GENERIC_WRITE|winio.WRITE_DAC|winio.WRITE_OWNER|winio.ACCESS_SYSTEM_SECURITY, createDisposition)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer dest.Close()

	err = winio.SetFileBasicInfo(dest, fileInfo)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	destw := winio.NewBackupFileWriter(dest, true)
	defer func() ***REMOVED***
		cerr := destw.Close()
		if err == nil ***REMOVED***
			err = cerr
		***REMOVED***
	***REMOVED***()

	_, err = io.Copy(destw, srcr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return fileInfo, nil
***REMOVED***

// cloneTree clones a directory tree using hard links. It skips hard links for
// the file names in the provided map and just copies those files.
func cloneTree(srcPath, destPath string, mutatedFiles map[string]bool) error ***REMOVED***
	var di []dirInfo
	err := filepath.Walk(srcPath, func(srcFilePath string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		relPath, err := filepath.Rel(srcPath, srcFilePath)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		destFilePath := filepath.Join(destPath, relPath)

		fileAttributes := info.Sys().(*syscall.Win32FileAttributeData).FileAttributes
		// Directories, reparse points, and files that will be mutated during
		// utility VM import must be copied. All other files can be hard linked.
		isReparsePoint := fileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0
		// In go1.9, FileInfo.IsDir() returns false if the directory is also a symlink.
		// See: https://github.com/golang/go/commit/1989921aef60c83e6f9127a8448fb5ede10e9acc
		// Fixes the problem by checking syscall.FILE_ATTRIBUTE_DIRECTORY directly
		isDir := fileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0

		if isDir || isReparsePoint || mutatedFiles[relPath] ***REMOVED***
			fi, err := copyFileWithMetadata(srcFilePath, destFilePath, isDir)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if isDir && !isReparsePoint ***REMOVED***
				di = append(di, dirInfo***REMOVED***path: destFilePath, fileInfo: *fi***REMOVED***)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = os.Link(srcFilePath, destFilePath)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		// Don't recurse on reparse points in go1.8 and older. Filepath.Walk
		// handles this in go1.9 and newer.
		if isDir && isReparsePoint && shouldSkipDirectoryReparse ***REMOVED***
			return filepath.SkipDir
		***REMOVED***

		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return reapplyDirectoryTimes(di)
***REMOVED***

func (w *legacyLayerWriter) Add(name string, fileInfo *winio.FileBasicInfo) error ***REMOVED***
	w.reset()
	err := w.init()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if name == utilityVMPath ***REMOVED***
		return w.initUtilityVM()
	***REMOVED***

	if hasPathPrefix(name, utilityVMPath) ***REMOVED***
		if !w.HasUtilityVM ***REMOVED***
			return errors.New("missing UtilityVM directory")
		***REMOVED***
		if !hasPathPrefix(name, utilityVMFilesPath) && name != utilityVMFilesPath ***REMOVED***
			return errors.New("invalid UtilityVM layer")
		***REMOVED***
		path := filepath.Join(w.destRoot, name)
		createDisposition := uint32(syscall.OPEN_EXISTING)
		if (fileInfo.FileAttributes & syscall.FILE_ATTRIBUTE_DIRECTORY) != 0 ***REMOVED***
			st, err := os.Lstat(path)
			if err != nil && !os.IsNotExist(err) ***REMOVED***
				return err
			***REMOVED***
			if st != nil ***REMOVED***
				// Delete the existing file/directory if it is not the same type as this directory.
				existingAttr := st.Sys().(*syscall.Win32FileAttributeData).FileAttributes
				if (uint32(fileInfo.FileAttributes)^existingAttr)&(syscall.FILE_ATTRIBUTE_DIRECTORY|syscall.FILE_ATTRIBUTE_REPARSE_POINT) != 0 ***REMOVED***
					if err = os.RemoveAll(path); err != nil ***REMOVED***
						return err
					***REMOVED***
					st = nil
				***REMOVED***
			***REMOVED***
			if st == nil ***REMOVED***
				if err = os.Mkdir(path, 0); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			if fileInfo.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT == 0 ***REMOVED***
				w.uvmDi = append(w.uvmDi, dirInfo***REMOVED***path: path, fileInfo: *fileInfo***REMOVED***)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Overwrite any existing hard link.
			err = os.Remove(path)
			if err != nil && !os.IsNotExist(err) ***REMOVED***
				return err
			***REMOVED***
			createDisposition = syscall.CREATE_NEW
		***REMOVED***

		f, err := openFileOrDir(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE|winio.WRITE_DAC|winio.WRITE_OWNER|winio.ACCESS_SYSTEM_SECURITY, createDisposition)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer func() ***REMOVED***
			if f != nil ***REMOVED***
				f.Close()
				os.Remove(path)
			***REMOVED***
		***REMOVED***()

		err = winio.SetFileBasicInfo(f, fileInfo)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		w.backupWriter = winio.NewBackupFileWriter(f, true)
		w.currentFile = f
		w.addedFiles[name] = true
		f = nil
		return nil
	***REMOVED***

	path := filepath.Join(w.root, name)
	if (fileInfo.FileAttributes & syscall.FILE_ATTRIBUTE_DIRECTORY) != 0 ***REMOVED***
		err := os.Mkdir(path, 0)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		path += ".$wcidirs$"
	***REMOVED***

	f, err := openFileOrDir(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, syscall.CREATE_NEW)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if f != nil ***REMOVED***
			f.Close()
			os.Remove(path)
		***REMOVED***
	***REMOVED***()

	strippedFi := *fileInfo
	strippedFi.FileAttributes = 0
	err = winio.SetFileBasicInfo(f, &strippedFi)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if hasPathPrefix(name, hivesPath) ***REMOVED***
		w.backupWriter = winio.NewBackupFileWriter(f, false)
	***REMOVED*** else ***REMOVED***
		// The file attributes are written before the stream.
		err = binary.Write(f, binary.LittleEndian, uint32(fileInfo.FileAttributes))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	w.currentFile = f
	w.addedFiles[name] = true
	f = nil
	return nil
***REMOVED***

func (w *legacyLayerWriter) AddLink(name string, target string) error ***REMOVED***
	w.reset()
	err := w.init()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var roots []string
	if hasPathPrefix(target, filesPath) ***REMOVED***
		// Look for cross-layer hard link targets in the parent layers, since
		// nothing is in the destination path yet.
		roots = w.parentRoots
	***REMOVED*** else if hasPathPrefix(target, utilityVMFilesPath) ***REMOVED***
		// Since the utility VM is fully cloned into the destination path
		// already, look for cross-layer hard link targets directly in the
		// destination path.
		roots = []string***REMOVED***w.destRoot***REMOVED***
	***REMOVED***

	if roots == nil || (!hasPathPrefix(name, filesPath) && !hasPathPrefix(name, utilityVMFilesPath)) ***REMOVED***
		return errors.New("invalid hard link in layer")
	***REMOVED***

	// Find to try the target of the link in a previously added file. If that
	// fails, search in parent layers.
	var selectedRoot string
	if _, ok := w.addedFiles[target]; ok ***REMOVED***
		selectedRoot = w.destRoot
	***REMOVED*** else ***REMOVED***
		for _, r := range roots ***REMOVED***
			if _, err = os.Lstat(filepath.Join(r, target)); err != nil ***REMOVED***
				if !os.IsNotExist(err) ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				selectedRoot = r
				break
			***REMOVED***
		***REMOVED***
		if selectedRoot == "" ***REMOVED***
			return fmt.Errorf("failed to find link target for '%s' -> '%s'", name, target)
		***REMOVED***
	***REMOVED***
	// The link can't be written until after the ImportLayer call.
	w.PendingLinks = append(w.PendingLinks, pendingLink***REMOVED***
		Path:   filepath.Join(w.destRoot, name),
		Target: filepath.Join(selectedRoot, target),
	***REMOVED***)
	w.addedFiles[name] = true
	return nil
***REMOVED***

func (w *legacyLayerWriter) Remove(name string) error ***REMOVED***
	if hasPathPrefix(name, filesPath) ***REMOVED***
		w.tombstones = append(w.tombstones, name[len(filesPath)+1:])
	***REMOVED*** else if hasPathPrefix(name, utilityVMFilesPath) ***REMOVED***
		err := w.initUtilityVM()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Make sure the path exists; os.RemoveAll will not fail if the file is
		// already gone, and this needs to be a fatal error for diagnostics
		// purposes.
		path := filepath.Join(w.destRoot, name)
		if _, err := os.Lstat(path); err != nil ***REMOVED***
			return err
		***REMOVED***
		err = os.RemoveAll(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return fmt.Errorf("invalid tombstone %s", name)
	***REMOVED***

	return nil
***REMOVED***

func (w *legacyLayerWriter) Write(b []byte) (int, error) ***REMOVED***
	if w.backupWriter == nil ***REMOVED***
		if w.currentFile == nil ***REMOVED***
			return 0, errors.New("closed")
		***REMOVED***
		return w.currentFile.Write(b)
	***REMOVED***
	return w.backupWriter.Write(b)
***REMOVED***

func (w *legacyLayerWriter) Close() error ***REMOVED***
	w.reset()
	err := w.init()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	tf, err := os.Create(filepath.Join(w.root, "tombstones.txt"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer tf.Close()
	_, err = tf.Write([]byte("\xef\xbb\xbfVersion 1.0\n"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, t := range w.tombstones ***REMOVED***
		_, err = tf.Write([]byte(filepath.Join(`\`, t) + "\n"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if w.HasUtilityVM ***REMOVED***
		err = reapplyDirectoryTimes(w.uvmDi)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
