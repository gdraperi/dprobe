// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// slashClean is equivalent to but slightly more efficient than
// path.Clean("/" + name).
func slashClean(name string) string ***REMOVED***
	if name == "" || name[0] != '/' ***REMOVED***
		name = "/" + name
	***REMOVED***
	return path.Clean(name)
***REMOVED***

// A FileSystem implements access to a collection of named files. The elements
// in a file path are separated by slash ('/', U+002F) characters, regardless
// of host operating system convention.
//
// Each method has the same semantics as the os package's function of the same
// name.
//
// Note that the os.Rename documentation says that "OS-specific restrictions
// might apply". In particular, whether or not renaming a file or directory
// overwriting another existing file or directory is an error is OS-dependent.
type FileSystem interface ***REMOVED***
	Mkdir(ctx context.Context, name string, perm os.FileMode) error
	OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (File, error)
	RemoveAll(ctx context.Context, name string) error
	Rename(ctx context.Context, oldName, newName string) error
	Stat(ctx context.Context, name string) (os.FileInfo, error)
***REMOVED***

// A File is returned by a FileSystem's OpenFile method and can be served by a
// Handler.
//
// A File may optionally implement the DeadPropsHolder interface, if it can
// load and save dead properties.
type File interface ***REMOVED***
	http.File
	io.Writer
***REMOVED***

// A Dir implements FileSystem using the native file system restricted to a
// specific directory tree.
//
// While the FileSystem.OpenFile method takes '/'-separated paths, a Dir's
// string value is a filename on the native file system, not a URL, so it is
// separated by filepath.Separator, which isn't necessarily '/'.
//
// An empty Dir is treated as ".".
type Dir string

func (d Dir) resolve(name string) string ***REMOVED***
	// This implementation is based on Dir.Open's code in the standard net/http package.
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") ***REMOVED***
		return ""
	***REMOVED***
	dir := string(d)
	if dir == "" ***REMOVED***
		dir = "."
	***REMOVED***
	return filepath.Join(dir, filepath.FromSlash(slashClean(name)))
***REMOVED***

func (d Dir) Mkdir(ctx context.Context, name string, perm os.FileMode) error ***REMOVED***
	if name = d.resolve(name); name == "" ***REMOVED***
		return os.ErrNotExist
	***REMOVED***
	return os.Mkdir(name, perm)
***REMOVED***

func (d Dir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	if name = d.resolve(name); name == "" ***REMOVED***
		return nil, os.ErrNotExist
	***REMOVED***
	f, err := os.OpenFile(name, flag, perm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return f, nil
***REMOVED***

func (d Dir) RemoveAll(ctx context.Context, name string) error ***REMOVED***
	if name = d.resolve(name); name == "" ***REMOVED***
		return os.ErrNotExist
	***REMOVED***
	if name == filepath.Clean(string(d)) ***REMOVED***
		// Prohibit removing the virtual root directory.
		return os.ErrInvalid
	***REMOVED***
	return os.RemoveAll(name)
***REMOVED***

func (d Dir) Rename(ctx context.Context, oldName, newName string) error ***REMOVED***
	if oldName = d.resolve(oldName); oldName == "" ***REMOVED***
		return os.ErrNotExist
	***REMOVED***
	if newName = d.resolve(newName); newName == "" ***REMOVED***
		return os.ErrNotExist
	***REMOVED***
	if root := filepath.Clean(string(d)); root == oldName || root == newName ***REMOVED***
		// Prohibit renaming from or to the virtual root directory.
		return os.ErrInvalid
	***REMOVED***
	return os.Rename(oldName, newName)
***REMOVED***

func (d Dir) Stat(ctx context.Context, name string) (os.FileInfo, error) ***REMOVED***
	if name = d.resolve(name); name == "" ***REMOVED***
		return nil, os.ErrNotExist
	***REMOVED***
	return os.Stat(name)
***REMOVED***

// NewMemFS returns a new in-memory FileSystem implementation.
func NewMemFS() FileSystem ***REMOVED***
	return &memFS***REMOVED***
		root: memFSNode***REMOVED***
			children: make(map[string]*memFSNode),
			mode:     0660 | os.ModeDir,
			modTime:  time.Now(),
		***REMOVED***,
	***REMOVED***
***REMOVED***

// A memFS implements FileSystem, storing all metadata and actual file data
// in-memory. No limits on filesystem size are used, so it is not recommended
// this be used where the clients are untrusted.
//
// Concurrent access is permitted. The tree structure is protected by a mutex,
// and each node's contents and metadata are protected by a per-node mutex.
//
// TODO: Enforce file permissions.
type memFS struct ***REMOVED***
	mu   sync.Mutex
	root memFSNode
***REMOVED***

// TODO: clean up and rationalize the walk/find code.

// walk walks the directory tree for the fullname, calling f at each step. If f
// returns an error, the walk will be aborted and return that same error.
//
// dir is the directory at that step, frag is the name fragment, and final is
// whether it is the final step. For example, walking "/foo/bar/x" will result
// in 3 calls to f:
//   - "/", "foo", false
//   - "/foo/", "bar", false
//   - "/foo/bar/", "x", true
// The frag argument will be empty only if dir is the root node and the walk
// ends at that root node.
func (fs *memFS) walk(op, fullname string, f func(dir *memFSNode, frag string, final bool) error) error ***REMOVED***
	original := fullname
	fullname = slashClean(fullname)

	// Strip any leading "/"s to make fullname a relative path, as the walk
	// starts at fs.root.
	if fullname[0] == '/' ***REMOVED***
		fullname = fullname[1:]
	***REMOVED***
	dir := &fs.root

	for ***REMOVED***
		frag, remaining := fullname, ""
		i := strings.IndexRune(fullname, '/')
		final := i < 0
		if !final ***REMOVED***
			frag, remaining = fullname[:i], fullname[i+1:]
		***REMOVED***
		if frag == "" && dir != &fs.root ***REMOVED***
			panic("webdav: empty path fragment for a clean path")
		***REMOVED***
		if err := f(dir, frag, final); err != nil ***REMOVED***
			return &os.PathError***REMOVED***
				Op:   op,
				Path: original,
				Err:  err,
			***REMOVED***
		***REMOVED***
		if final ***REMOVED***
			break
		***REMOVED***
		child := dir.children[frag]
		if child == nil ***REMOVED***
			return &os.PathError***REMOVED***
				Op:   op,
				Path: original,
				Err:  os.ErrNotExist,
			***REMOVED***
		***REMOVED***
		if !child.mode.IsDir() ***REMOVED***
			return &os.PathError***REMOVED***
				Op:   op,
				Path: original,
				Err:  os.ErrInvalid,
			***REMOVED***
		***REMOVED***
		dir, fullname = child, remaining
	***REMOVED***
	return nil
***REMOVED***

// find returns the parent of the named node and the relative name fragment
// from the parent to the child. For example, if finding "/foo/bar/baz" then
// parent will be the node for "/foo/bar" and frag will be "baz".
//
// If the fullname names the root node, then parent, frag and err will be zero.
//
// find returns an error if the parent does not already exist or the parent
// isn't a directory, but it will not return an error per se if the child does
// not already exist. The error returned is either nil or an *os.PathError
// whose Op is op.
func (fs *memFS) find(op, fullname string) (parent *memFSNode, frag string, err error) ***REMOVED***
	err = fs.walk(op, fullname, func(parent0 *memFSNode, frag0 string, final bool) error ***REMOVED***
		if !final ***REMOVED***
			return nil
		***REMOVED***
		if frag0 != "" ***REMOVED***
			parent, frag = parent0, frag0
		***REMOVED***
		return nil
	***REMOVED***)
	return parent, frag, err
***REMOVED***

func (fs *memFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("mkdir", name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if dir == nil ***REMOVED***
		// We can't create the root.
		return os.ErrInvalid
	***REMOVED***
	if _, ok := dir.children[frag]; ok ***REMOVED***
		return os.ErrExist
	***REMOVED***
	dir.children[frag] = &memFSNode***REMOVED***
		children: make(map[string]*memFSNode),
		mode:     perm.Perm() | os.ModeDir,
		modTime:  time.Now(),
	***REMOVED***
	return nil
***REMOVED***

func (fs *memFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("open", name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var n *memFSNode
	if dir == nil ***REMOVED***
		// We're opening the root.
		if flag&(os.O_WRONLY|os.O_RDWR) != 0 ***REMOVED***
			return nil, os.ErrPermission
		***REMOVED***
		n, frag = &fs.root, "/"

	***REMOVED*** else ***REMOVED***
		n = dir.children[frag]
		if flag&(os.O_SYNC|os.O_APPEND) != 0 ***REMOVED***
			// memFile doesn't support these flags yet.
			return nil, os.ErrInvalid
		***REMOVED***
		if flag&os.O_CREATE != 0 ***REMOVED***
			if flag&os.O_EXCL != 0 && n != nil ***REMOVED***
				return nil, os.ErrExist
			***REMOVED***
			if n == nil ***REMOVED***
				n = &memFSNode***REMOVED***
					mode: perm.Perm(),
				***REMOVED***
				dir.children[frag] = n
			***REMOVED***
		***REMOVED***
		if n == nil ***REMOVED***
			return nil, os.ErrNotExist
		***REMOVED***
		if flag&(os.O_WRONLY|os.O_RDWR) != 0 && flag&os.O_TRUNC != 0 ***REMOVED***
			n.mu.Lock()
			n.data = nil
			n.mu.Unlock()
		***REMOVED***
	***REMOVED***

	children := make([]os.FileInfo, 0, len(n.children))
	for cName, c := range n.children ***REMOVED***
		children = append(children, c.stat(cName))
	***REMOVED***
	return &memFile***REMOVED***
		n:                n,
		nameSnapshot:     frag,
		childrenSnapshot: children,
	***REMOVED***, nil
***REMOVED***

func (fs *memFS) RemoveAll(ctx context.Context, name string) error ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("remove", name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if dir == nil ***REMOVED***
		// We can't remove the root.
		return os.ErrInvalid
	***REMOVED***
	delete(dir.children, frag)
	return nil
***REMOVED***

func (fs *memFS) Rename(ctx context.Context, oldName, newName string) error ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()

	oldName = slashClean(oldName)
	newName = slashClean(newName)
	if oldName == newName ***REMOVED***
		return nil
	***REMOVED***
	if strings.HasPrefix(newName, oldName+"/") ***REMOVED***
		// We can't rename oldName to be a sub-directory of itself.
		return os.ErrInvalid
	***REMOVED***

	oDir, oFrag, err := fs.find("rename", oldName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if oDir == nil ***REMOVED***
		// We can't rename from the root.
		return os.ErrInvalid
	***REMOVED***

	nDir, nFrag, err := fs.find("rename", newName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if nDir == nil ***REMOVED***
		// We can't rename to the root.
		return os.ErrInvalid
	***REMOVED***

	oNode, ok := oDir.children[oFrag]
	if !ok ***REMOVED***
		return os.ErrNotExist
	***REMOVED***
	if oNode.children != nil ***REMOVED***
		if nNode, ok := nDir.children[nFrag]; ok ***REMOVED***
			if nNode.children == nil ***REMOVED***
				return errNotADirectory
			***REMOVED***
			if len(nNode.children) != 0 ***REMOVED***
				return errDirectoryNotEmpty
			***REMOVED***
		***REMOVED***
	***REMOVED***
	delete(oDir.children, oFrag)
	nDir.children[nFrag] = oNode
	return nil
***REMOVED***

func (fs *memFS) Stat(ctx context.Context, name string) (os.FileInfo, error) ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("stat", name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if dir == nil ***REMOVED***
		// We're stat'ting the root.
		return fs.root.stat("/"), nil
	***REMOVED***
	if n, ok := dir.children[frag]; ok ***REMOVED***
		return n.stat(path.Base(name)), nil
	***REMOVED***
	return nil, os.ErrNotExist
***REMOVED***

// A memFSNode represents a single entry in the in-memory filesystem and also
// implements os.FileInfo.
type memFSNode struct ***REMOVED***
	// children is protected by memFS.mu.
	children map[string]*memFSNode

	mu        sync.Mutex
	data      []byte
	mode      os.FileMode
	modTime   time.Time
	deadProps map[xml.Name]Property
***REMOVED***

func (n *memFSNode) stat(name string) *memFileInfo ***REMOVED***
	n.mu.Lock()
	defer n.mu.Unlock()
	return &memFileInfo***REMOVED***
		name:    name,
		size:    int64(len(n.data)),
		mode:    n.mode,
		modTime: n.modTime,
	***REMOVED***
***REMOVED***

func (n *memFSNode) DeadProps() (map[xml.Name]Property, error) ***REMOVED***
	n.mu.Lock()
	defer n.mu.Unlock()
	if len(n.deadProps) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	ret := make(map[xml.Name]Property, len(n.deadProps))
	for k, v := range n.deadProps ***REMOVED***
		ret[k] = v
	***REMOVED***
	return ret, nil
***REMOVED***

func (n *memFSNode) Patch(patches []Proppatch) ([]Propstat, error) ***REMOVED***
	n.mu.Lock()
	defer n.mu.Unlock()
	pstat := Propstat***REMOVED***Status: http.StatusOK***REMOVED***
	for _, patch := range patches ***REMOVED***
		for _, p := range patch.Props ***REMOVED***
			pstat.Props = append(pstat.Props, Property***REMOVED***XMLName: p.XMLName***REMOVED***)
			if patch.Remove ***REMOVED***
				delete(n.deadProps, p.XMLName)
				continue
			***REMOVED***
			if n.deadProps == nil ***REMOVED***
				n.deadProps = map[xml.Name]Property***REMOVED******REMOVED***
			***REMOVED***
			n.deadProps[p.XMLName] = p
		***REMOVED***
	***REMOVED***
	return []Propstat***REMOVED***pstat***REMOVED***, nil
***REMOVED***

type memFileInfo struct ***REMOVED***
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
***REMOVED***

func (f *memFileInfo) Name() string       ***REMOVED*** return f.name ***REMOVED***
func (f *memFileInfo) Size() int64        ***REMOVED*** return f.size ***REMOVED***
func (f *memFileInfo) Mode() os.FileMode  ***REMOVED*** return f.mode ***REMOVED***
func (f *memFileInfo) ModTime() time.Time ***REMOVED*** return f.modTime ***REMOVED***
func (f *memFileInfo) IsDir() bool        ***REMOVED*** return f.mode.IsDir() ***REMOVED***
func (f *memFileInfo) Sys() interface***REMOVED******REMOVED***   ***REMOVED*** return nil ***REMOVED***

// A memFile is a File implementation for a memFSNode. It is a per-file (not
// per-node) read/write position, and a snapshot of the memFS' tree structure
// (a node's name and children) for that node.
type memFile struct ***REMOVED***
	n                *memFSNode
	nameSnapshot     string
	childrenSnapshot []os.FileInfo
	// pos is protected by n.mu.
	pos int
***REMOVED***

// A *memFile implements the optional DeadPropsHolder interface.
var _ DeadPropsHolder = (*memFile)(nil)

func (f *memFile) DeadProps() (map[xml.Name]Property, error)     ***REMOVED*** return f.n.DeadProps() ***REMOVED***
func (f *memFile) Patch(patches []Proppatch) ([]Propstat, error) ***REMOVED*** return f.n.Patch(patches) ***REMOVED***

func (f *memFile) Close() error ***REMOVED***
	return nil
***REMOVED***

func (f *memFile) Read(p []byte) (int, error) ***REMOVED***
	f.n.mu.Lock()
	defer f.n.mu.Unlock()
	if f.n.mode.IsDir() ***REMOVED***
		return 0, os.ErrInvalid
	***REMOVED***
	if f.pos >= len(f.n.data) ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	n := copy(p, f.n.data[f.pos:])
	f.pos += n
	return n, nil
***REMOVED***

func (f *memFile) Readdir(count int) ([]os.FileInfo, error) ***REMOVED***
	f.n.mu.Lock()
	defer f.n.mu.Unlock()
	if !f.n.mode.IsDir() ***REMOVED***
		return nil, os.ErrInvalid
	***REMOVED***
	old := f.pos
	if old >= len(f.childrenSnapshot) ***REMOVED***
		// The os.File Readdir docs say that at the end of a directory,
		// the error is io.EOF if count > 0 and nil if count <= 0.
		if count > 0 ***REMOVED***
			return nil, io.EOF
		***REMOVED***
		return nil, nil
	***REMOVED***
	if count > 0 ***REMOVED***
		f.pos += count
		if f.pos > len(f.childrenSnapshot) ***REMOVED***
			f.pos = len(f.childrenSnapshot)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		f.pos = len(f.childrenSnapshot)
		old = 0
	***REMOVED***
	return f.childrenSnapshot[old:f.pos], nil
***REMOVED***

func (f *memFile) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	f.n.mu.Lock()
	defer f.n.mu.Unlock()
	npos := f.pos
	// TODO: How to handle offsets greater than the size of system int?
	switch whence ***REMOVED***
	case os.SEEK_SET:
		npos = int(offset)
	case os.SEEK_CUR:
		npos += int(offset)
	case os.SEEK_END:
		npos = len(f.n.data) + int(offset)
	default:
		npos = -1
	***REMOVED***
	if npos < 0 ***REMOVED***
		return 0, os.ErrInvalid
	***REMOVED***
	f.pos = npos
	return int64(f.pos), nil
***REMOVED***

func (f *memFile) Stat() (os.FileInfo, error) ***REMOVED***
	return f.n.stat(f.nameSnapshot), nil
***REMOVED***

func (f *memFile) Write(p []byte) (int, error) ***REMOVED***
	lenp := len(p)
	f.n.mu.Lock()
	defer f.n.mu.Unlock()

	if f.n.mode.IsDir() ***REMOVED***
		return 0, os.ErrInvalid
	***REMOVED***
	if f.pos < len(f.n.data) ***REMOVED***
		n := copy(f.n.data[f.pos:], p)
		f.pos += n
		p = p[n:]
	***REMOVED*** else if f.pos > len(f.n.data) ***REMOVED***
		// Write permits the creation of holes, if we've seek'ed past the
		// existing end of file.
		if f.pos <= cap(f.n.data) ***REMOVED***
			oldLen := len(f.n.data)
			f.n.data = f.n.data[:f.pos]
			hole := f.n.data[oldLen:]
			for i := range hole ***REMOVED***
				hole[i] = 0
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			d := make([]byte, f.pos, f.pos+len(p))
			copy(d, f.n.data)
			f.n.data = d
		***REMOVED***
	***REMOVED***

	if len(p) > 0 ***REMOVED***
		// We should only get here if f.pos == len(f.n.data).
		f.n.data = append(f.n.data, p...)
		f.pos = len(f.n.data)
	***REMOVED***
	f.n.modTime = time.Now()
	return lenp, nil
***REMOVED***

// moveFiles moves files and/or directories from src to dst.
//
// See section 9.9.4 for when various HTTP status codes apply.
func moveFiles(ctx context.Context, fs FileSystem, src, dst string, overwrite bool) (status int, err error) ***REMOVED***
	created := false
	if _, err := fs.Stat(ctx, dst); err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return http.StatusForbidden, err
		***REMOVED***
		created = true
	***REMOVED*** else if overwrite ***REMOVED***
		// Section 9.9.3 says that "If a resource exists at the destination
		// and the Overwrite header is "T", then prior to performing the move,
		// the server must perform a DELETE with "Depth: infinity" on the
		// destination resource.
		if err := fs.RemoveAll(ctx, dst); err != nil ***REMOVED***
			return http.StatusForbidden, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return http.StatusPreconditionFailed, os.ErrExist
	***REMOVED***
	if err := fs.Rename(ctx, src, dst); err != nil ***REMOVED***
		return http.StatusForbidden, err
	***REMOVED***
	if created ***REMOVED***
		return http.StatusCreated, nil
	***REMOVED***
	return http.StatusNoContent, nil
***REMOVED***

func copyProps(dst, src File) error ***REMOVED***
	d, ok := dst.(DeadPropsHolder)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	s, ok := src.(DeadPropsHolder)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	m, err := s.DeadProps()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	props := make([]Property, 0, len(m))
	for _, prop := range m ***REMOVED***
		props = append(props, prop)
	***REMOVED***
	_, err = d.Patch([]Proppatch***REMOVED******REMOVED***Props: props***REMOVED******REMOVED***)
	return err
***REMOVED***

// copyFiles copies files and/or directories from src to dst.
//
// See section 9.8.5 for when various HTTP status codes apply.
func copyFiles(ctx context.Context, fs FileSystem, src, dst string, overwrite bool, depth int, recursion int) (status int, err error) ***REMOVED***
	if recursion == 1000 ***REMOVED***
		return http.StatusInternalServerError, errRecursionTooDeep
	***REMOVED***
	recursion++

	// TODO: section 9.8.3 says that "Note that an infinite-depth COPY of /A/
	// into /A/B/ could lead to infinite recursion if not handled correctly."

	srcFile, err := fs.OpenFile(ctx, src, os.O_RDONLY, 0)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return http.StatusNotFound, err
		***REMOVED***
		return http.StatusInternalServerError, err
	***REMOVED***
	defer srcFile.Close()
	srcStat, err := srcFile.Stat()
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return http.StatusNotFound, err
		***REMOVED***
		return http.StatusInternalServerError, err
	***REMOVED***
	srcPerm := srcStat.Mode() & os.ModePerm

	created := false
	if _, err := fs.Stat(ctx, dst); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			created = true
		***REMOVED*** else ***REMOVED***
			return http.StatusForbidden, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !overwrite ***REMOVED***
			return http.StatusPreconditionFailed, os.ErrExist
		***REMOVED***
		if err := fs.RemoveAll(ctx, dst); err != nil && !os.IsNotExist(err) ***REMOVED***
			return http.StatusForbidden, err
		***REMOVED***
	***REMOVED***

	if srcStat.IsDir() ***REMOVED***
		if err := fs.Mkdir(ctx, dst, srcPerm); err != nil ***REMOVED***
			return http.StatusForbidden, err
		***REMOVED***
		if depth == infiniteDepth ***REMOVED***
			children, err := srcFile.Readdir(-1)
			if err != nil ***REMOVED***
				return http.StatusForbidden, err
			***REMOVED***
			for _, c := range children ***REMOVED***
				name := c.Name()
				s := path.Join(src, name)
				d := path.Join(dst, name)
				cStatus, cErr := copyFiles(ctx, fs, s, d, overwrite, depth, recursion)
				if cErr != nil ***REMOVED***
					// TODO: MultiStatus.
					return cStatus, cErr
				***REMOVED***
			***REMOVED***
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		dstFile, err := fs.OpenFile(ctx, dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcPerm)
		if err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				return http.StatusConflict, err
			***REMOVED***
			return http.StatusForbidden, err

		***REMOVED***
		_, copyErr := io.Copy(dstFile, srcFile)
		propsErr := copyProps(dstFile, srcFile)
		closeErr := dstFile.Close()
		if copyErr != nil ***REMOVED***
			return http.StatusInternalServerError, copyErr
		***REMOVED***
		if propsErr != nil ***REMOVED***
			return http.StatusInternalServerError, propsErr
		***REMOVED***
		if closeErr != nil ***REMOVED***
			return http.StatusInternalServerError, closeErr
		***REMOVED***
	***REMOVED***

	if created ***REMOVED***
		return http.StatusCreated, nil
	***REMOVED***
	return http.StatusNoContent, nil
***REMOVED***

// walkFS traverses filesystem fs starting at name up to depth levels.
//
// Allowed values for depth are 0, 1 or infiniteDepth. For each visited node,
// walkFS calls walkFn. If a visited file system node is a directory and
// walkFn returns filepath.SkipDir, walkFS will skip traversal of this node.
func walkFS(ctx context.Context, fs FileSystem, depth int, name string, info os.FileInfo, walkFn filepath.WalkFunc) error ***REMOVED***
	// This implementation is based on Walk's code in the standard path/filepath package.
	err := walkFn(name, info, nil)
	if err != nil ***REMOVED***
		if info.IsDir() && err == filepath.SkipDir ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
	if !info.IsDir() || depth == 0 ***REMOVED***
		return nil
	***REMOVED***
	if depth == 1 ***REMOVED***
		depth = 0
	***REMOVED***

	// Read directory names.
	f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return walkFn(name, info, err)
	***REMOVED***
	fileInfos, err := f.Readdir(0)
	f.Close()
	if err != nil ***REMOVED***
		return walkFn(name, info, err)
	***REMOVED***

	for _, fileInfo := range fileInfos ***REMOVED***
		filename := path.Join(name, fileInfo.Name())
		fileInfo, err := fs.Stat(ctx, filename)
		if err != nil ***REMOVED***
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = walkFS(ctx, fs, depth, filename, fileInfo, walkFn)
			if err != nil ***REMOVED***
				if !fileInfo.IsDir() || err != filepath.SkipDir ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
