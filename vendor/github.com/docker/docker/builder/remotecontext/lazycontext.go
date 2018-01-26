package remotecontext

import (
	"encoding/hex"
	"os"
	"strings"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/pools"
	"github.com/pkg/errors"
)

// NewLazySource creates a new LazyContext. LazyContext defines a hashed build
// context based on a root directory. Individual files are hashed first time
// they are asked. It is not safe to call methods of LazyContext concurrently.
func NewLazySource(root containerfs.ContainerFS) (builder.Source, error) ***REMOVED***
	return &lazySource***REMOVED***
		root: root,
		sums: make(map[string]string),
	***REMOVED***, nil
***REMOVED***

type lazySource struct ***REMOVED***
	root containerfs.ContainerFS
	sums map[string]string
***REMOVED***

func (c *lazySource) Root() containerfs.ContainerFS ***REMOVED***
	return c.root
***REMOVED***

func (c *lazySource) Close() error ***REMOVED***
	return nil
***REMOVED***

func (c *lazySource) Hash(path string) (string, error) ***REMOVED***
	cleanPath, fullPath, err := normalize(path, c.root)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	relPath, err := Rel(c.root, fullPath)
	if err != nil ***REMOVED***
		return "", errors.WithStack(convertPathError(err, cleanPath))
	***REMOVED***

	fi, err := os.Lstat(fullPath)
	if err != nil ***REMOVED***
		// Backwards compatibility: a missing file returns a path as hash.
		// This is reached in the case of a broken symlink.
		return relPath, nil
	***REMOVED***

	sum, ok := c.sums[relPath]
	if !ok ***REMOVED***
		sum, err = c.prepareHash(relPath, fi)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***

	return sum, nil
***REMOVED***

func (c *lazySource) prepareHash(relPath string, fi os.FileInfo) (string, error) ***REMOVED***
	p := c.root.Join(c.root.Path(), relPath)
	h, err := NewFileHash(p, relPath, fi)
	if err != nil ***REMOVED***
		return "", errors.Wrapf(err, "failed to create hash for %s", relPath)
	***REMOVED***
	if fi.Mode().IsRegular() && fi.Size() > 0 ***REMOVED***
		f, err := c.root.Open(p)
		if err != nil ***REMOVED***
			return "", errors.Wrapf(err, "failed to open %s", relPath)
		***REMOVED***
		defer f.Close()
		if _, err := pools.Copy(h, f); err != nil ***REMOVED***
			return "", errors.Wrapf(err, "failed to copy file data for %s", relPath)
		***REMOVED***
	***REMOVED***
	sum := hex.EncodeToString(h.Sum(nil))
	c.sums[relPath] = sum
	return sum, nil
***REMOVED***

// Rel makes a path relative to base path. Same as `filepath.Rel` but can also
// handle UUID paths in windows.
func Rel(basepath containerfs.ContainerFS, targpath string) (string, error) ***REMOVED***
	// filepath.Rel can't handle UUID paths in windows
	if basepath.OS() == "windows" ***REMOVED***
		pfx := basepath.Path() + `\`
		if strings.HasPrefix(targpath, pfx) ***REMOVED***
			p := strings.TrimPrefix(targpath, pfx)
			if p == "" ***REMOVED***
				p = "."
			***REMOVED***
			return p, nil
		***REMOVED***
	***REMOVED***
	return basepath.Rel(basepath.Path(), targpath)
***REMOVED***
