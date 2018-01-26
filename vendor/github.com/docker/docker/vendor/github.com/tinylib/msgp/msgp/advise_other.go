// +build !linux appengine

package msgp

import (
	"os"
)

// TODO: darwin, BSD support

func adviseRead(mem []byte) ***REMOVED******REMOVED***

func adviseWrite(mem []byte) ***REMOVED******REMOVED***

func fallocate(f *os.File, sz int64) error ***REMOVED***
	return f.Truncate(sz)
***REMOVED***
