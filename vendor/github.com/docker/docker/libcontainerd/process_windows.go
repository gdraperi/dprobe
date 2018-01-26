package libcontainerd

import (
	"io"
	"sync"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/docker/pkg/ioutils"
)

type autoClosingReader struct ***REMOVED***
	io.ReadCloser
	sync.Once
***REMOVED***

func (r *autoClosingReader) Read(b []byte) (n int, err error) ***REMOVED***
	n, err = r.ReadCloser.Read(b)
	if err != nil ***REMOVED***
		r.Once.Do(func() ***REMOVED*** r.ReadCloser.Close() ***REMOVED***)
	***REMOVED***
	return
***REMOVED***

func createStdInCloser(pipe io.WriteCloser, process hcsshim.Process) io.WriteCloser ***REMOVED***
	return ioutils.NewWriteCloserWrapper(pipe, func() error ***REMOVED***
		if err := pipe.Close(); err != nil ***REMOVED***
			return err
		***REMOVED***

		err := process.CloseStdin()
		if err != nil && !hcsshim.IsNotExist(err) && !hcsshim.IsAlreadyClosed(err) ***REMOVED***
			// This error will occur if the compute system is currently shutting down
			if perr, ok := err.(*hcsshim.ProcessError); ok && perr.Err != hcsshim.ErrVmcomputeOperationInvalidState ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***)
***REMOVED***

func (p *process) Cleanup() error ***REMOVED***
	return nil
***REMOVED***
