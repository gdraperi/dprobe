package vfs

import "github.com/docker/docker/daemon/graphdriver/copy"

func dirCopy(srcDir, dstDir string) error ***REMOVED***
	return copy.DirCopy(srcDir, dstDir, copy.Content, false)
***REMOVED***
