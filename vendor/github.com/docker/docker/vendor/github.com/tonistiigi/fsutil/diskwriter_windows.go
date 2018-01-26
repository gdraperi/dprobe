// +build windows

package fsutil

import (
	"os"
	"time"

	"github.com/pkg/errors"
)

func rewriteMetadata(p string, stat *Stat) error ***REMOVED***
	return chtimes(p, stat.ModTime)
***REMOVED***

func chtimes(path string, un int64) error ***REMOVED***
	mtime := time.Unix(0, un)
	return os.Chtimes(path, mtime, mtime)
***REMOVED***

// handleTarTypeBlockCharFifo is an OS-specific helper function used by
// createTarFile to handle the following types of header: Block; Char; Fifo
func handleTarTypeBlockCharFifo(path string, stat *Stat) error ***REMOVED***
	return errors.New("Not implemented on windows")
***REMOVED***
