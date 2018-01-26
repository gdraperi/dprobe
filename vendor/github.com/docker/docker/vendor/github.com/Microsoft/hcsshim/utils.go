package hcsshim

import (
	"io"
	"syscall"

	"github.com/Microsoft/go-winio"
)

// makeOpenFiles calls winio.MakeOpenFile for each handle in a slice but closes all the handles
// if there is an error.
func makeOpenFiles(hs []syscall.Handle) (_ []io.ReadWriteCloser, err error) ***REMOVED***
	fs := make([]io.ReadWriteCloser, len(hs))
	for i, h := range hs ***REMOVED***
		if h != syscall.Handle(0) ***REMOVED***
			if err == nil ***REMOVED***
				fs[i], err = winio.MakeOpenFile(h)
			***REMOVED***
			if err != nil ***REMOVED***
				syscall.Close(h)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		for _, f := range fs ***REMOVED***
			if f != nil ***REMOVED***
				f.Close()
			***REMOVED***
		***REMOVED***
		return nil, err
	***REMOVED***
	return fs, nil
***REMOVED***
