package archive

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"unsafe"

	"github.com/docker/docker/pkg/system"
	"golang.org/x/sys/unix"
)

// walker is used to implement collectFileInfoForChanges on linux. Where this
// method in general returns the entire contents of two directory trees, we
// optimize some FS calls out on linux. In particular, we take advantage of the
// fact that getdents(2) returns the inode of each file in the directory being
// walked, which, when walking two trees in parallel to generate a list of
// changes, can be used to prune subtrees without ever having to lstat(2) them
// directly. Eliminating stat calls in this way can save up to seconds on large
// images.
type walker struct ***REMOVED***
	dir1  string
	dir2  string
	root1 *FileInfo
	root2 *FileInfo
***REMOVED***

// collectFileInfoForChanges returns a complete representation of the trees
// rooted at dir1 and dir2, with one important exception: any subtree or
// leaf where the inode and device numbers are an exact match between dir1
// and dir2 will be pruned from the results. This method is *only* to be used
// to generating a list of changes between the two directories, as it does not
// reflect the full contents.
func collectFileInfoForChanges(dir1, dir2 string) (*FileInfo, *FileInfo, error) ***REMOVED***
	w := &walker***REMOVED***
		dir1:  dir1,
		dir2:  dir2,
		root1: newRootFileInfo(),
		root2: newRootFileInfo(),
	***REMOVED***

	i1, err := os.Lstat(w.dir1)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	i2, err := os.Lstat(w.dir2)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if err := w.walk("/", i1, i2); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return w.root1, w.root2, nil
***REMOVED***

// Given a FileInfo, its path info, and a reference to the root of the tree
// being constructed, register this file with the tree.
func walkchunk(path string, fi os.FileInfo, dir string, root *FileInfo) error ***REMOVED***
	if fi == nil ***REMOVED***
		return nil
	***REMOVED***
	parent := root.LookUp(filepath.Dir(path))
	if parent == nil ***REMOVED***
		return fmt.Errorf("walkchunk: Unexpectedly no parent for %s", path)
	***REMOVED***
	info := &FileInfo***REMOVED***
		name:     filepath.Base(path),
		children: make(map[string]*FileInfo),
		parent:   parent,
	***REMOVED***
	cpath := filepath.Join(dir, path)
	stat, err := system.FromStatT(fi.Sys().(*syscall.Stat_t))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	info.stat = stat
	info.capability, _ = system.Lgetxattr(cpath, "security.capability") // lgetxattr(2): fs access
	parent.children[info.name] = info
	return nil
***REMOVED***

// Walk a subtree rooted at the same path in both trees being iterated. For
// example, /docker/overlay/1234/a/b/c/d and /docker/overlay/8888/a/b/c/d
func (w *walker) walk(path string, i1, i2 os.FileInfo) (err error) ***REMOVED***
	// Register these nodes with the return trees, unless we're still at the
	// (already-created) roots:
	if path != "/" ***REMOVED***
		if err := walkchunk(path, i1, w.dir1, w.root1); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := walkchunk(path, i2, w.dir2, w.root2); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	is1Dir := i1 != nil && i1.IsDir()
	is2Dir := i2 != nil && i2.IsDir()

	sameDevice := false
	if i1 != nil && i2 != nil ***REMOVED***
		si1 := i1.Sys().(*syscall.Stat_t)
		si2 := i2.Sys().(*syscall.Stat_t)
		if si1.Dev == si2.Dev ***REMOVED***
			sameDevice = true
		***REMOVED***
	***REMOVED***

	// If these files are both non-existent, or leaves (non-dirs), we are done.
	if !is1Dir && !is2Dir ***REMOVED***
		return nil
	***REMOVED***

	// Fetch the names of all the files contained in both directories being walked:
	var names1, names2 []nameIno
	if is1Dir ***REMOVED***
		names1, err = readdirnames(filepath.Join(w.dir1, path)) // getdents(2): fs access
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if is2Dir ***REMOVED***
		names2, err = readdirnames(filepath.Join(w.dir2, path)) // getdents(2): fs access
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// We have lists of the files contained in both parallel directories, sorted
	// in the same order. Walk them in parallel, generating a unique merged list
	// of all items present in either or both directories.
	var names []string
	ix1 := 0
	ix2 := 0

	for ***REMOVED***
		if ix1 >= len(names1) ***REMOVED***
			break
		***REMOVED***
		if ix2 >= len(names2) ***REMOVED***
			break
		***REMOVED***

		ni1 := names1[ix1]
		ni2 := names2[ix2]

		switch bytes.Compare([]byte(ni1.name), []byte(ni2.name)) ***REMOVED***
		case -1: // ni1 < ni2 -- advance ni1
			// we will not encounter ni1 in names2
			names = append(names, ni1.name)
			ix1++
		case 0: // ni1 == ni2
			if ni1.ino != ni2.ino || !sameDevice ***REMOVED***
				names = append(names, ni1.name)
			***REMOVED***
			ix1++
			ix2++
		case 1: // ni1 > ni2 -- advance ni2
			// we will not encounter ni2 in names1
			names = append(names, ni2.name)
			ix2++
		***REMOVED***
	***REMOVED***
	for ix1 < len(names1) ***REMOVED***
		names = append(names, names1[ix1].name)
		ix1++
	***REMOVED***
	for ix2 < len(names2) ***REMOVED***
		names = append(names, names2[ix2].name)
		ix2++
	***REMOVED***

	// For each of the names present in either or both of the directories being
	// iterated, stat the name under each root, and recurse the pair of them:
	for _, name := range names ***REMOVED***
		fname := filepath.Join(path, name)
		var cInfo1, cInfo2 os.FileInfo
		if is1Dir ***REMOVED***
			cInfo1, err = os.Lstat(filepath.Join(w.dir1, fname)) // lstat(2): fs access
			if err != nil && !os.IsNotExist(err) ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if is2Dir ***REMOVED***
			cInfo2, err = os.Lstat(filepath.Join(w.dir2, fname)) // lstat(2): fs access
			if err != nil && !os.IsNotExist(err) ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if err = w.walk(fname, cInfo1, cInfo2); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// ***REMOVED***name,inode***REMOVED*** pairs used to support the early-pruning logic of the walker type
type nameIno struct ***REMOVED***
	name string
	ino  uint64
***REMOVED***

type nameInoSlice []nameIno

func (s nameInoSlice) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s nameInoSlice) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s nameInoSlice) Less(i, j int) bool ***REMOVED*** return s[i].name < s[j].name ***REMOVED***

// readdirnames is a hacked-apart version of the Go stdlib code, exposing inode
// numbers further up the stack when reading directory contents. Unlike
// os.Readdirnames, which returns a list of filenames, this function returns a
// list of ***REMOVED***filename,inode***REMOVED*** pairs.
func readdirnames(dirname string) (names []nameIno, err error) ***REMOVED***
	var (
		size = 100
		buf  = make([]byte, 4096)
		nbuf int
		bufp int
		nb   int
	)

	f, err := os.Open(dirname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	names = make([]nameIno, 0, size) // Empty with room to grow.
	for ***REMOVED***
		// Refill the buffer if necessary
		if bufp >= nbuf ***REMOVED***
			bufp = 0
			nbuf, err = unix.ReadDirent(int(f.Fd()), buf) // getdents on linux
			if nbuf < 0 ***REMOVED***
				nbuf = 0
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, os.NewSyscallError("readdirent", err)
			***REMOVED***
			if nbuf <= 0 ***REMOVED***
				break // EOF
			***REMOVED***
		***REMOVED***

		// Drain the buffer
		nb, names = parseDirent(buf[bufp:nbuf], names)
		bufp += nb
	***REMOVED***

	sl := nameInoSlice(names)
	sort.Sort(sl)
	return sl, nil
***REMOVED***

// parseDirent is a minor modification of unix.ParseDirent (linux version)
// which returns ***REMOVED***name,inode***REMOVED*** pairs instead of just names.
func parseDirent(buf []byte, names []nameIno) (consumed int, newnames []nameIno) ***REMOVED***
	origlen := len(buf)
	for len(buf) > 0 ***REMOVED***
		dirent := (*unix.Dirent)(unsafe.Pointer(&buf[0]))
		buf = buf[dirent.Reclen:]
		if dirent.Ino == 0 ***REMOVED*** // File absent in directory.
			continue
		***REMOVED***
		bytes := (*[10000]byte)(unsafe.Pointer(&dirent.Name[0]))
		var name = string(bytes[0:clen(bytes[:])])
		if name == "." || name == ".." ***REMOVED*** // Useless names
			continue
		***REMOVED***
		names = append(names, nameIno***REMOVED***name, dirent.Ino***REMOVED***)
	***REMOVED***
	return origlen - len(buf), names
***REMOVED***

func clen(n []byte) int ***REMOVED***
	for i := 0; i < len(n); i++ ***REMOVED***
		if n[i] == 0 ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return len(n)
***REMOVED***

// OverlayChanges walks the path rw and determines changes for the files in the path,
// with respect to the parent layers
func OverlayChanges(layers []string, rw string) ([]Change, error) ***REMOVED***
	return changes(layers, rw, overlayDeletedFile, nil)
***REMOVED***

func overlayDeletedFile(root, path string, fi os.FileInfo) (string, error) ***REMOVED***
	if fi.Mode()&os.ModeCharDevice != 0 ***REMOVED***
		s := fi.Sys().(*syscall.Stat_t)
		if unix.Major(uint64(s.Rdev)) == 0 && unix.Minor(uint64(s.Rdev)) == 0 ***REMOVED*** // nolint: unconvert
			return path, nil
		***REMOVED***
	***REMOVED***
	if fi.Mode()&os.ModeDir != 0 ***REMOVED***
		opaque, err := system.Lgetxattr(filepath.Join(root, path), "trusted.overlay.opaque")
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if len(opaque) == 1 && opaque[0] == 'y' ***REMOVED***
			return path, nil
		***REMOVED***
	***REMOVED***

	return "", nil

***REMOVED***
