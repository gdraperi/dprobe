package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

// ChangeType represents the change type.
type ChangeType int

const (
	// ChangeModify represents the modify operation.
	ChangeModify = iota
	// ChangeAdd represents the add operation.
	ChangeAdd
	// ChangeDelete represents the delete operation.
	ChangeDelete
)

func (c ChangeType) String() string ***REMOVED***
	switch c ***REMOVED***
	case ChangeModify:
		return "C"
	case ChangeAdd:
		return "A"
	case ChangeDelete:
		return "D"
	***REMOVED***
	return ""
***REMOVED***

// Change represents a change, it wraps the change type and path.
// It describes changes of the files in the path respect to the
// parent layers. The change could be modify, add, delete.
// This is used for layer diff.
type Change struct ***REMOVED***
	Path string
	Kind ChangeType
***REMOVED***

func (change *Change) String() string ***REMOVED***
	return fmt.Sprintf("%s %s", change.Kind, change.Path)
***REMOVED***

// for sort.Sort
type changesByPath []Change

func (c changesByPath) Less(i, j int) bool ***REMOVED*** return c[i].Path < c[j].Path ***REMOVED***
func (c changesByPath) Len() int           ***REMOVED*** return len(c) ***REMOVED***
func (c changesByPath) Swap(i, j int)      ***REMOVED*** c[j], c[i] = c[i], c[j] ***REMOVED***

// Gnu tar and the go tar writer don't have sub-second mtime
// precision, which is problematic when we apply changes via tar
// files, we handle this by comparing for exact times, *or* same
// second count and either a or b having exactly 0 nanoseconds
func sameFsTime(a, b time.Time) bool ***REMOVED***
	return a == b ||
		(a.Unix() == b.Unix() &&
			(a.Nanosecond() == 0 || b.Nanosecond() == 0))
***REMOVED***

func sameFsTimeSpec(a, b syscall.Timespec) bool ***REMOVED***
	return a.Sec == b.Sec &&
		(a.Nsec == b.Nsec || a.Nsec == 0 || b.Nsec == 0)
***REMOVED***

// Changes walks the path rw and determines changes for the files in the path,
// with respect to the parent layers
func Changes(layers []string, rw string) ([]Change, error) ***REMOVED***
	return changes(layers, rw, aufsDeletedFile, aufsMetadataSkip)
***REMOVED***

func aufsMetadataSkip(path string) (skip bool, err error) ***REMOVED***
	skip, err = filepath.Match(string(os.PathSeparator)+WhiteoutMetaPrefix+"*", path)
	if err != nil ***REMOVED***
		skip = true
	***REMOVED***
	return
***REMOVED***

func aufsDeletedFile(root, path string, fi os.FileInfo) (string, error) ***REMOVED***
	f := filepath.Base(path)

	// If there is a whiteout, then the file was removed
	if strings.HasPrefix(f, WhiteoutPrefix) ***REMOVED***
		originalFile := f[len(WhiteoutPrefix):]
		return filepath.Join(filepath.Dir(path), originalFile), nil
	***REMOVED***

	return "", nil
***REMOVED***

type skipChange func(string) (bool, error)
type deleteChange func(string, string, os.FileInfo) (string, error)

func changes(layers []string, rw string, dc deleteChange, sc skipChange) ([]Change, error) ***REMOVED***
	var (
		changes     []Change
		changedDirs = make(map[string]struct***REMOVED******REMOVED***)
	)

	err := filepath.Walk(rw, func(path string, f os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Rebase path
		path, err = filepath.Rel(rw, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// As this runs on the daemon side, file paths are OS specific.
		path = filepath.Join(string(os.PathSeparator), path)

		// Skip root
		if path == string(os.PathSeparator) ***REMOVED***
			return nil
		***REMOVED***

		if sc != nil ***REMOVED***
			if skip, err := sc(path); skip ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		change := Change***REMOVED***
			Path: path,
		***REMOVED***

		deletedFile, err := dc(rw, path, f)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Find out what kind of modification happened
		if deletedFile != "" ***REMOVED***
			change.Path = deletedFile
			change.Kind = ChangeDelete
		***REMOVED*** else ***REMOVED***
			// Otherwise, the file was added
			change.Kind = ChangeAdd

			// ...Unless it already existed in a top layer, in which case, it's a modification
			for _, layer := range layers ***REMOVED***
				stat, err := os.Stat(filepath.Join(layer, path))
				if err != nil && !os.IsNotExist(err) ***REMOVED***
					return err
				***REMOVED***
				if err == nil ***REMOVED***
					// The file existed in the top layer, so that's a modification

					// However, if it's a directory, maybe it wasn't actually modified.
					// If you modify /foo/bar/baz, then /foo will be part of the changed files only because it's the parent of bar
					if stat.IsDir() && f.IsDir() ***REMOVED***
						if f.Size() == stat.Size() && f.Mode() == stat.Mode() && sameFsTime(f.ModTime(), stat.ModTime()) ***REMOVED***
							// Both directories are the same, don't record the change
							return nil
						***REMOVED***
					***REMOVED***
					change.Kind = ChangeModify
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// If /foo/bar/file.txt is modified, then /foo/bar must be part of the changed files.
		// This block is here to ensure the change is recorded even if the
		// modify time, mode and size of the parent directory in the rw and ro layers are all equal.
		// Check https://github.com/docker/docker/pull/13590 for details.
		if f.IsDir() ***REMOVED***
			changedDirs[path] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
		if change.Kind == ChangeAdd || change.Kind == ChangeDelete ***REMOVED***
			parent := filepath.Dir(path)
			if _, ok := changedDirs[parent]; !ok && parent != "/" ***REMOVED***
				changes = append(changes, Change***REMOVED***Path: parent, Kind: ChangeModify***REMOVED***)
				changedDirs[parent] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***

		// Record change
		changes = append(changes, change)
		return nil
	***REMOVED***)
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED***
	return changes, nil
***REMOVED***

// FileInfo describes the information of a file.
type FileInfo struct ***REMOVED***
	parent     *FileInfo
	name       string
	stat       *system.StatT
	children   map[string]*FileInfo
	capability []byte
	added      bool
***REMOVED***

// LookUp looks up the file information of a file.
func (info *FileInfo) LookUp(path string) *FileInfo ***REMOVED***
	// As this runs on the daemon side, file paths are OS specific.
	parent := info
	if path == string(os.PathSeparator) ***REMOVED***
		return info
	***REMOVED***

	pathElements := strings.Split(path, string(os.PathSeparator))
	for _, elem := range pathElements ***REMOVED***
		if elem != "" ***REMOVED***
			child := parent.children[elem]
			if child == nil ***REMOVED***
				return nil
			***REMOVED***
			parent = child
		***REMOVED***
	***REMOVED***
	return parent
***REMOVED***

func (info *FileInfo) path() string ***REMOVED***
	if info.parent == nil ***REMOVED***
		// As this runs on the daemon side, file paths are OS specific.
		return string(os.PathSeparator)
	***REMOVED***
	return filepath.Join(info.parent.path(), info.name)
***REMOVED***

func (info *FileInfo) addChanges(oldInfo *FileInfo, changes *[]Change) ***REMOVED***

	sizeAtEntry := len(*changes)

	if oldInfo == nil ***REMOVED***
		// add
		change := Change***REMOVED***
			Path: info.path(),
			Kind: ChangeAdd,
		***REMOVED***
		*changes = append(*changes, change)
		info.added = true
	***REMOVED***

	// We make a copy so we can modify it to detect additions
	// also, we only recurse on the old dir if the new info is a directory
	// otherwise any previous delete/change is considered recursive
	oldChildren := make(map[string]*FileInfo)
	if oldInfo != nil && info.isDir() ***REMOVED***
		for k, v := range oldInfo.children ***REMOVED***
			oldChildren[k] = v
		***REMOVED***
	***REMOVED***

	for name, newChild := range info.children ***REMOVED***
		oldChild := oldChildren[name]
		if oldChild != nil ***REMOVED***
			// change?
			oldStat := oldChild.stat
			newStat := newChild.stat
			// Note: We can't compare inode or ctime or blocksize here, because these change
			// when copying a file into a container. However, that is not generally a problem
			// because any content change will change mtime, and any status change should
			// be visible when actually comparing the stat fields. The only time this
			// breaks down is if some code intentionally hides a change by setting
			// back mtime
			if statDifferent(oldStat, newStat) ||
				!bytes.Equal(oldChild.capability, newChild.capability) ***REMOVED***
				change := Change***REMOVED***
					Path: newChild.path(),
					Kind: ChangeModify,
				***REMOVED***
				*changes = append(*changes, change)
				newChild.added = true
			***REMOVED***

			// Remove from copy so we can detect deletions
			delete(oldChildren, name)
		***REMOVED***

		newChild.addChanges(oldChild, changes)
	***REMOVED***
	for _, oldChild := range oldChildren ***REMOVED***
		// delete
		change := Change***REMOVED***
			Path: oldChild.path(),
			Kind: ChangeDelete,
		***REMOVED***
		*changes = append(*changes, change)
	***REMOVED***

	// If there were changes inside this directory, we need to add it, even if the directory
	// itself wasn't changed. This is needed to properly save and restore filesystem permissions.
	// As this runs on the daemon side, file paths are OS specific.
	if len(*changes) > sizeAtEntry && info.isDir() && !info.added && info.path() != string(os.PathSeparator) ***REMOVED***
		change := Change***REMOVED***
			Path: info.path(),
			Kind: ChangeModify,
		***REMOVED***
		// Let's insert the directory entry before the recently added entries located inside this dir
		*changes = append(*changes, change) // just to resize the slice, will be overwritten
		copy((*changes)[sizeAtEntry+1:], (*changes)[sizeAtEntry:])
		(*changes)[sizeAtEntry] = change
	***REMOVED***

***REMOVED***

// Changes add changes to file information.
func (info *FileInfo) Changes(oldInfo *FileInfo) []Change ***REMOVED***
	var changes []Change

	info.addChanges(oldInfo, &changes)

	return changes
***REMOVED***

func newRootFileInfo() *FileInfo ***REMOVED***
	// As this runs on the daemon side, file paths are OS specific.
	root := &FileInfo***REMOVED***
		name:     string(os.PathSeparator),
		children: make(map[string]*FileInfo),
	***REMOVED***
	return root
***REMOVED***

// ChangesDirs compares two directories and generates an array of Change objects describing the changes.
// If oldDir is "", then all files in newDir will be Add-Changes.
func ChangesDirs(newDir, oldDir string) ([]Change, error) ***REMOVED***
	var (
		oldRoot, newRoot *FileInfo
	)
	if oldDir == "" ***REMOVED***
		emptyDir, err := ioutil.TempDir("", "empty")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer os.Remove(emptyDir)
		oldDir = emptyDir
	***REMOVED***
	oldRoot, newRoot, err := collectFileInfoForChanges(oldDir, newDir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return newRoot.Changes(oldRoot), nil
***REMOVED***

// ChangesSize calculates the size in bytes of the provided changes, based on newDir.
func ChangesSize(newDir string, changes []Change) int64 ***REMOVED***
	var (
		size int64
		sf   = make(map[uint64]struct***REMOVED******REMOVED***)
	)
	for _, change := range changes ***REMOVED***
		if change.Kind == ChangeModify || change.Kind == ChangeAdd ***REMOVED***
			file := filepath.Join(newDir, change.Path)
			fileInfo, err := os.Lstat(file)
			if err != nil ***REMOVED***
				logrus.Errorf("Can not stat %q: %s", file, err)
				continue
			***REMOVED***

			if fileInfo != nil && !fileInfo.IsDir() ***REMOVED***
				if hasHardlinks(fileInfo) ***REMOVED***
					inode := getIno(fileInfo)
					if _, ok := sf[inode]; !ok ***REMOVED***
						size += fileInfo.Size()
						sf[inode] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					size += fileInfo.Size()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return size
***REMOVED***

// ExportChanges produces an Archive from the provided changes, relative to dir.
func ExportChanges(dir string, changes []Change, uidMaps, gidMaps []idtools.IDMap) (io.ReadCloser, error) ***REMOVED***
	reader, writer := io.Pipe()
	go func() ***REMOVED***
		ta := newTarAppender(idtools.NewIDMappingsFromMaps(uidMaps, gidMaps), writer, nil)

		// this buffer is needed for the duration of this piped stream
		defer pools.BufioWriter32KPool.Put(ta.Buffer)

		sort.Sort(changesByPath(changes))

		// In general we log errors here but ignore them because
		// during e.g. a diff operation the container can continue
		// mutating the filesystem and we can see transient errors
		// from this
		for _, change := range changes ***REMOVED***
			if change.Kind == ChangeDelete ***REMOVED***
				whiteOutDir := filepath.Dir(change.Path)
				whiteOutBase := filepath.Base(change.Path)
				whiteOut := filepath.Join(whiteOutDir, WhiteoutPrefix+whiteOutBase)
				timestamp := time.Now()
				hdr := &tar.Header***REMOVED***
					Name:       whiteOut[1:],
					Size:       0,
					ModTime:    timestamp,
					AccessTime: timestamp,
					ChangeTime: timestamp,
				***REMOVED***
				if err := ta.TarWriter.WriteHeader(hdr); err != nil ***REMOVED***
					logrus.Debugf("Can't write whiteout header: %s", err)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				path := filepath.Join(dir, change.Path)
				if err := ta.addTarFile(path, change.Path[1:]); err != nil ***REMOVED***
					logrus.Debugf("Can't add file %s to tar: %s", path, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Make sure to check the error on Close.
		if err := ta.TarWriter.Close(); err != nil ***REMOVED***
			logrus.Debugf("Can't close layer: %s", err)
		***REMOVED***
		if err := writer.Close(); err != nil ***REMOVED***
			logrus.Debugf("failed close Changes writer: %s", err)
		***REMOVED***
	***REMOVED***()
	return reader, nil
***REMOVED***
