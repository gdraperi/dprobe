package fifo

import (
	"io"
	"os"
	"runtime"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type fifo struct ***REMOVED***
	flag        int
	opened      chan struct***REMOVED******REMOVED***
	closed      chan struct***REMOVED******REMOVED***
	closing     chan struct***REMOVED******REMOVED***
	err         error
	file        *os.File
	closingOnce sync.Once // close has been called
	closedOnce  sync.Once // fifo is closed
	handle      *handle
***REMOVED***

var leakCheckWg *sync.WaitGroup

// OpenFifo opens a fifo. Returns io.ReadWriteCloser.
// Context can be used to cancel this function until open(2) has not returned.
// Accepted flags:
// - syscall.O_CREAT - create new fifo if one doesn't exist
// - syscall.O_RDONLY - open fifo only from reader side
// - syscall.O_WRONLY - open fifo only from writer side
// - syscall.O_RDWR - open fifo from both sides, never block on syscall level
// - syscall.O_NONBLOCK - return io.ReadWriteCloser even if other side of the
//     fifo isn't open. read/write will be connected after the actual fifo is
//     open or after fifo is closed.
func OpenFifo(ctx context.Context, fn string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) ***REMOVED***
	if _, err := os.Stat(fn); err != nil ***REMOVED***
		if os.IsNotExist(err) && flag&syscall.O_CREAT != 0 ***REMOVED***
			if err := mkfifo(fn, uint32(perm&os.ModePerm)); err != nil && !os.IsExist(err) ***REMOVED***
				return nil, errors.Wrapf(err, "error creating fifo %v", fn)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	block := flag&syscall.O_NONBLOCK == 0 || flag&syscall.O_RDWR != 0

	flag &= ^syscall.O_CREAT
	flag &= ^syscall.O_NONBLOCK

	h, err := getHandle(fn)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	f := &fifo***REMOVED***
		handle:  h,
		flag:    flag,
		opened:  make(chan struct***REMOVED******REMOVED***),
		closed:  make(chan struct***REMOVED******REMOVED***),
		closing: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	wg := leakCheckWg
	if wg != nil ***REMOVED***
		wg.Add(2)
	***REMOVED***

	go func() ***REMOVED***
		if wg != nil ***REMOVED***
			defer wg.Done()
		***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			select ***REMOVED***
			case <-f.opened:
			default:
				f.Close()
			***REMOVED***
		case <-f.opened:
		case <-f.closed:
		***REMOVED***
	***REMOVED***()
	go func() ***REMOVED***
		if wg != nil ***REMOVED***
			defer wg.Done()
		***REMOVED***
		var file *os.File
		fn, err := h.Path()
		if err == nil ***REMOVED***
			file, err = os.OpenFile(fn, flag, 0)
		***REMOVED***
		select ***REMOVED***
		case <-f.closing:
			if err == nil ***REMOVED***
				select ***REMOVED***
				case <-ctx.Done():
					err = ctx.Err()
				default:
					err = errors.Errorf("fifo %v was closed before opening", h.Name())
				***REMOVED***
				if file != nil ***REMOVED***
					file.Close()
				***REMOVED***
			***REMOVED***
		default:
		***REMOVED***
		if err != nil ***REMOVED***
			f.closedOnce.Do(func() ***REMOVED***
				f.err = err
				close(f.closed)
			***REMOVED***)
			return
		***REMOVED***
		f.file = file
		close(f.opened)
	***REMOVED***()
	if block ***REMOVED***
		select ***REMOVED***
		case <-f.opened:
		case <-f.closed:
			return nil, f.err
		***REMOVED***
	***REMOVED***
	return f, nil
***REMOVED***

// Read from a fifo to a byte array.
func (f *fifo) Read(b []byte) (int, error) ***REMOVED***
	if f.flag&syscall.O_WRONLY > 0 ***REMOVED***
		return 0, errors.New("reading from write-only fifo")
	***REMOVED***
	select ***REMOVED***
	case <-f.opened:
		return f.file.Read(b)
	default:
	***REMOVED***
	select ***REMOVED***
	case <-f.opened:
		return f.file.Read(b)
	case <-f.closed:
		return 0, errors.New("reading from a closed fifo")
	***REMOVED***
***REMOVED***

// Write from byte array to a fifo.
func (f *fifo) Write(b []byte) (int, error) ***REMOVED***
	if f.flag&(syscall.O_WRONLY|syscall.O_RDWR) == 0 ***REMOVED***
		return 0, errors.New("writing to read-only fifo")
	***REMOVED***
	select ***REMOVED***
	case <-f.opened:
		return f.file.Write(b)
	default:
	***REMOVED***
	select ***REMOVED***
	case <-f.opened:
		return f.file.Write(b)
	case <-f.closed:
		return 0, errors.New("writing to a closed fifo")
	***REMOVED***
***REMOVED***

// Close the fifo. Next reads/writes will error. This method can also be used
// before open(2) has returned and fifo was never opened.
func (f *fifo) Close() (retErr error) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-f.closed:
			f.handle.Close()
			return
		default:
			select ***REMOVED***
			case <-f.opened:
				f.closedOnce.Do(func() ***REMOVED***
					retErr = f.file.Close()
					f.err = retErr
					close(f.closed)
				***REMOVED***)
			default:
				if f.flag&syscall.O_RDWR != 0 ***REMOVED***
					runtime.Gosched()
					break
				***REMOVED***
				f.closingOnce.Do(func() ***REMOVED***
					close(f.closing)
				***REMOVED***)
				reverseMode := syscall.O_WRONLY
				if f.flag&syscall.O_WRONLY > 0 ***REMOVED***
					reverseMode = syscall.O_RDONLY
				***REMOVED***
				fn, err := f.handle.Path()
				// if Close() is called concurrently(shouldn't) it may cause error
				// because handle is closed
				select ***REMOVED***
				case <-f.closed:
				default:
					if err != nil ***REMOVED***
						// Path has become invalid. We will leak a goroutine.
						// This case should not happen in linux.
						f.closedOnce.Do(func() ***REMOVED***
							f.err = err
							close(f.closed)
						***REMOVED***)
						<-f.closed
						break
					***REMOVED***
					f, err := os.OpenFile(fn, reverseMode|syscall.O_NONBLOCK, 0)
					if err == nil ***REMOVED***
						f.Close()
					***REMOVED***
					runtime.Gosched()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
