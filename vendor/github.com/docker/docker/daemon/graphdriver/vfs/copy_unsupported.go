// +build !linux

package vfs

import "github.com/docker/docker/pkg/chrootarchive"

func dirCopy(srcDir, dstDir string) error ***REMOVED***
	return chrootarchive.NewArchiver(nil).CopyWithTar(srcDir, dstDir)
***REMOVED***
