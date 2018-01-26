// +build !linux

package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/pkg/system"
)

func collectFileInfoForChanges(oldDir, newDir string) (*FileInfo, *FileInfo, error) ***REMOVED***
	var (
		oldRoot, newRoot *FileInfo
		err1, err2       error
		errs             = make(chan error, 2)
	)
	go func() ***REMOVED***
		oldRoot, err1 = collectFileInfo(oldDir)
		errs <- err1
	***REMOVED***()
	go func() ***REMOVED***
		newRoot, err2 = collectFileInfo(newDir)
		errs <- err2
	***REMOVED***()

	// block until both routines have returned
	for i := 0; i < 2; i++ ***REMOVED***
		if err := <-errs; err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***

	return oldRoot, newRoot, nil
***REMOVED***

func collectFileInfo(sourceDir string) (*FileInfo, error) ***REMOVED***
	root := newRootFileInfo()

	err := filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Rebase path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// As this runs on the daemon side, file paths are OS specific.
		relPath = filepath.Join(string(os.PathSeparator), relPath)

		// See https://github.com/golang/go/issues/9168 - bug in filepath.Join.
		// Temporary workaround. If the returned path starts with two backslashes,
		// trim it down to a single backslash. Only relevant on Windows.
		if runtime.GOOS == "windows" ***REMOVED***
			if strings.HasPrefix(relPath, `\\`) ***REMOVED***
				relPath = relPath[1:]
			***REMOVED***
		***REMOVED***

		if relPath == string(os.PathSeparator) ***REMOVED***
			return nil
		***REMOVED***

		parent := root.LookUp(filepath.Dir(relPath))
		if parent == nil ***REMOVED***
			return fmt.Errorf("collectFileInfo: Unexpectedly no parent for %s", relPath)
		***REMOVED***

		info := &FileInfo***REMOVED***
			name:     filepath.Base(relPath),
			children: make(map[string]*FileInfo),
			parent:   parent,
		***REMOVED***

		s, err := system.Lstat(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		info.stat = s

		info.capability, _ = system.Lgetxattr(path, "security.capability")

		parent.children[info.name] = info

		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return root, nil
***REMOVED***
