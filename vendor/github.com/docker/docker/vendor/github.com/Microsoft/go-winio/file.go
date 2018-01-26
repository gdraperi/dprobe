// +build windows

package winio

import (
	"errors"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

//sys cancelIoEx(file syscall.Handle, o *syscall.Overlapped) (err error) = CancelIoEx
//sys createIoCompletionPort(file syscall.Handle, port syscall.Handle, key uintptr, threadCount uint32) (newport syscall.Handle, err error) = CreateIoCompletionPort
//sys getQueuedCompletionStatus(port syscall.Handle, bytes *uint32, key *uintptr, o **ioOperation, timeout uint32) (err error) = GetQueuedCompletionStatus
//sys setFileCompletionNotificationModes(h syscall.Handle, flags uint8) (err error) = SetFileCompletionNotificationModes
//sys timeBeginPeriod(period uint32) (n int32) = winmm.timeBeginPeriod

type atomicBool int32

func (b *atomicBool) isSet() bool ***REMOVED*** return atomic.LoadInt32((*int32)(b)) != 0 ***REMOVED***
func (b *atomicBool) setFalse()   ***REMOVED*** atomic.StoreInt32((*int32)(b), 0) ***REMOVED***
func (b *atomicBool) setTrue()    ***REMOVED*** atomic.StoreInt32((*int32)(b), 1) ***REMOVED***
func (b *atomicBool) swap(new bool) bool ***REMOVED***
	var newInt int32
	if new ***REMOVED***
		newInt = 1
	***REMOVED***
	return atomic.SwapInt32((*int32)(b), newInt) == 1
***REMOVED***

const (
	cFILE_SKIP_COMPLETION_PORT_ON_SUCCESS = 1
	cFILE_SKIP_SET_EVENT_ON_HANDLE        = 2
)

var (
	ErrFileClosed = errors.New("file has already been closed")
	ErrTimeout    = &timeoutError***REMOVED******REMOVED***
)

type timeoutError struct***REMOVED******REMOVED***

func (e *timeoutError) Error() string   ***REMOVED*** return "i/o timeout" ***REMOVED***
func (e *timeoutError) Timeout() bool   ***REMOVED*** return true ***REMOVED***
func (e *timeoutError) Temporary() bool ***REMOVED*** return true ***REMOVED***

type timeoutChan chan struct***REMOVED******REMOVED***

var ioInitOnce sync.Once
var ioCompletionPort syscall.Handle

// ioResult contains the result of an asynchronous IO operation
type ioResult struct ***REMOVED***
	bytes uint32
	err   error
***REMOVED***

// ioOperation represents an outstanding asynchronous Win32 IO
type ioOperation struct ***REMOVED***
	o  syscall.Overlapped
	ch chan ioResult
***REMOVED***

func initIo() ***REMOVED***
	h, err := createIoCompletionPort(syscall.InvalidHandle, 0, 0, 0xffffffff)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	ioCompletionPort = h
	go ioCompletionProcessor(h)
***REMOVED***

// win32File implements Reader, Writer, and Closer on a Win32 handle without blocking in a syscall.
// It takes ownership of this handle and will close it if it is garbage collected.
type win32File struct ***REMOVED***
	handle        syscall.Handle
	wg            sync.WaitGroup
	wgLock        sync.RWMutex
	closing       atomicBool
	readDeadline  deadlineHandler
	writeDeadline deadlineHandler
***REMOVED***

type deadlineHandler struct ***REMOVED***
	setLock     sync.Mutex
	channel     timeoutChan
	channelLock sync.RWMutex
	timer       *time.Timer
	timedout    atomicBool
***REMOVED***

// makeWin32File makes a new win32File from an existing file handle
func makeWin32File(h syscall.Handle) (*win32File, error) ***REMOVED***
	f := &win32File***REMOVED***handle: h***REMOVED***
	ioInitOnce.Do(initIo)
	_, err := createIoCompletionPort(h, ioCompletionPort, 0, 0xffffffff)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = setFileCompletionNotificationModes(h, cFILE_SKIP_COMPLETION_PORT_ON_SUCCESS|cFILE_SKIP_SET_EVENT_ON_HANDLE)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	f.readDeadline.channel = make(timeoutChan)
	f.writeDeadline.channel = make(timeoutChan)
	return f, nil
***REMOVED***

func MakeOpenFile(h syscall.Handle) (io.ReadWriteCloser, error) ***REMOVED***
	return makeWin32File(h)
***REMOVED***

// closeHandle closes the resources associated with a Win32 handle
func (f *win32File) closeHandle() ***REMOVED***
	f.wgLock.Lock()
	// Atomically set that we are closing, releasing the resources only once.
	if !f.closing.swap(true) ***REMOVED***
		f.wgLock.Unlock()
		// cancel all IO and wait for it to complete
		cancelIoEx(f.handle, nil)
		f.wg.Wait()
		// at this point, no new IO can start
		syscall.Close(f.handle)
		f.handle = 0
	***REMOVED*** else ***REMOVED***
		f.wgLock.Unlock()
	***REMOVED***
***REMOVED***

// Close closes a win32File.
func (f *win32File) Close() error ***REMOVED***
	f.closeHandle()
	return nil
***REMOVED***

// prepareIo prepares for a new IO operation.
// The caller must call f.wg.Done() when the IO is finished, prior to Close() returning.
func (f *win32File) prepareIo() (*ioOperation, error) ***REMOVED***
	f.wgLock.RLock()
	if f.closing.isSet() ***REMOVED***
		f.wgLock.RUnlock()
		return nil, ErrFileClosed
	***REMOVED***
	f.wg.Add(1)
	f.wgLock.RUnlock()
	c := &ioOperation***REMOVED******REMOVED***
	c.ch = make(chan ioResult)
	return c, nil
***REMOVED***

// ioCompletionProcessor processes completed async IOs forever
func ioCompletionProcessor(h syscall.Handle) ***REMOVED***
	// Set the timer resolution to 1. This fixes a performance regression in golang 1.6.
	timeBeginPeriod(1)
	for ***REMOVED***
		var bytes uint32
		var key uintptr
		var op *ioOperation
		err := getQueuedCompletionStatus(h, &bytes, &key, &op, syscall.INFINITE)
		if op == nil ***REMOVED***
			panic(err)
		***REMOVED***
		op.ch <- ioResult***REMOVED***bytes, err***REMOVED***
	***REMOVED***
***REMOVED***

// asyncIo processes the return value from ReadFile or WriteFile, blocking until
// the operation has actually completed.
func (f *win32File) asyncIo(c *ioOperation, d *deadlineHandler, bytes uint32, err error) (int, error) ***REMOVED***
	if err != syscall.ERROR_IO_PENDING ***REMOVED***
		return int(bytes), err
	***REMOVED***

	if f.closing.isSet() ***REMOVED***
		cancelIoEx(f.handle, &c.o)
	***REMOVED***

	var timeout timeoutChan
	if d != nil ***REMOVED***
		d.channelLock.Lock()
		timeout = d.channel
		d.channelLock.Unlock()
	***REMOVED***

	var r ioResult
	select ***REMOVED***
	case r = <-c.ch:
		err = r.err
		if err == syscall.ERROR_OPERATION_ABORTED ***REMOVED***
			if f.closing.isSet() ***REMOVED***
				err = ErrFileClosed
			***REMOVED***
		***REMOVED***
	case <-timeout:
		cancelIoEx(f.handle, &c.o)
		r = <-c.ch
		err = r.err
		if err == syscall.ERROR_OPERATION_ABORTED ***REMOVED***
			err = ErrTimeout
		***REMOVED***
	***REMOVED***

	// runtime.KeepAlive is needed, as c is passed via native
	// code to ioCompletionProcessor, c must remain alive
	// until the channel read is complete.
	runtime.KeepAlive(c)
	return int(r.bytes), err
***REMOVED***

// Read reads from a file handle.
func (f *win32File) Read(b []byte) (int, error) ***REMOVED***
	c, err := f.prepareIo()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer f.wg.Done()

	if f.readDeadline.timedout.isSet() ***REMOVED***
		return 0, ErrTimeout
	***REMOVED***

	var bytes uint32
	err = syscall.ReadFile(f.handle, b, &bytes, &c.o)
	n, err := f.asyncIo(c, &f.readDeadline, bytes, err)
	runtime.KeepAlive(b)

	// Handle EOF conditions.
	if err == nil && n == 0 && len(b) != 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED*** else if err == syscall.ERROR_BROKEN_PIPE ***REMOVED***
		return 0, io.EOF
	***REMOVED*** else ***REMOVED***
		return n, err
	***REMOVED***
***REMOVED***

// Write writes to a file handle.
func (f *win32File) Write(b []byte) (int, error) ***REMOVED***
	c, err := f.prepareIo()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer f.wg.Done()

	if f.writeDeadline.timedout.isSet() ***REMOVED***
		return 0, ErrTimeout
	***REMOVED***

	var bytes uint32
	err = syscall.WriteFile(f.handle, b, &bytes, &c.o)
	n, err := f.asyncIo(c, &f.writeDeadline, bytes, err)
	runtime.KeepAlive(b)
	return n, err
***REMOVED***

func (f *win32File) SetReadDeadline(deadline time.Time) error ***REMOVED***
	return f.readDeadline.set(deadline)
***REMOVED***

func (f *win32File) SetWriteDeadline(deadline time.Time) error ***REMOVED***
	return f.writeDeadline.set(deadline)
***REMOVED***

func (f *win32File) Flush() error ***REMOVED***
	return syscall.FlushFileBuffers(f.handle)
***REMOVED***

func (d *deadlineHandler) set(deadline time.Time) error ***REMOVED***
	d.setLock.Lock()
	defer d.setLock.Unlock()

	if d.timer != nil ***REMOVED***
		if !d.timer.Stop() ***REMOVED***
			<-d.channel
		***REMOVED***
		d.timer = nil
	***REMOVED***
	d.timedout.setFalse()

	select ***REMOVED***
	case <-d.channel:
		d.channelLock.Lock()
		d.channel = make(chan struct***REMOVED******REMOVED***)
		d.channelLock.Unlock()
	default:
	***REMOVED***

	if deadline.IsZero() ***REMOVED***
		return nil
	***REMOVED***

	timeoutIO := func() ***REMOVED***
		d.timedout.setTrue()
		close(d.channel)
	***REMOVED***

	now := time.Now()
	duration := deadline.Sub(now)
	if deadline.After(now) ***REMOVED***
		// Deadline is in the future, set a timer to wait
		d.timer = time.AfterFunc(duration, timeoutIO)
	***REMOVED*** else ***REMOVED***
		// Deadline is in the past. Cancel all pending IO now.
		timeoutIO()
	***REMOVED***
	return nil
***REMOVED***
