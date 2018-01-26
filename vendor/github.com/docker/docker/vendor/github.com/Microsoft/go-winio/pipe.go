// +build windows

package winio

import (
	"errors"
	"io"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)

//sys connectNamedPipe(pipe syscall.Handle, o *syscall.Overlapped) (err error) = ConnectNamedPipe
//sys createNamedPipe(name string, flags uint32, pipeMode uint32, maxInstances uint32, outSize uint32, inSize uint32, defaultTimeout uint32, sa *syscall.SecurityAttributes) (handle syscall.Handle, err error)  [failretval==syscall.InvalidHandle] = CreateNamedPipeW
//sys createFile(name string, access uint32, mode uint32, sa *syscall.SecurityAttributes, createmode uint32, attrs uint32, templatefile syscall.Handle) (handle syscall.Handle, err error) [failretval==syscall.InvalidHandle] = CreateFileW
//sys waitNamedPipe(name string, timeout uint32) (err error) = WaitNamedPipeW
//sys getNamedPipeInfo(pipe syscall.Handle, flags *uint32, outSize *uint32, inSize *uint32, maxInstances *uint32) (err error) = GetNamedPipeInfo
//sys getNamedPipeHandleState(pipe syscall.Handle, state *uint32, curInstances *uint32, maxCollectionCount *uint32, collectDataTimeout *uint32, userName *uint16, maxUserNameSize uint32) (err error) = GetNamedPipeHandleStateW
//sys localAlloc(uFlags uint32, length uint32) (ptr uintptr) = LocalAlloc

const (
	cERROR_PIPE_BUSY      = syscall.Errno(231)
	cERROR_NO_DATA        = syscall.Errno(232)
	cERROR_PIPE_CONNECTED = syscall.Errno(535)
	cERROR_SEM_TIMEOUT    = syscall.Errno(121)

	cPIPE_ACCESS_DUPLEX            = 0x3
	cFILE_FLAG_FIRST_PIPE_INSTANCE = 0x80000
	cSECURITY_SQOS_PRESENT         = 0x100000
	cSECURITY_ANONYMOUS            = 0

	cPIPE_REJECT_REMOTE_CLIENTS = 0x8

	cPIPE_UNLIMITED_INSTANCES = 255

	cNMPWAIT_USE_DEFAULT_WAIT = 0
	cNMPWAIT_NOWAIT           = 1

	cPIPE_TYPE_MESSAGE = 4

	cPIPE_READMODE_MESSAGE = 2
)

var (
	// ErrPipeListenerClosed is returned for pipe operations on listeners that have been closed.
	// This error should match net.errClosing since docker takes a dependency on its text.
	ErrPipeListenerClosed = errors.New("use of closed network connection")

	errPipeWriteClosed = errors.New("pipe has been closed for write")
)

type win32Pipe struct ***REMOVED***
	*win32File
	path string
***REMOVED***

type win32MessageBytePipe struct ***REMOVED***
	win32Pipe
	writeClosed bool
	readEOF     bool
***REMOVED***

type pipeAddress string

func (f *win32Pipe) LocalAddr() net.Addr ***REMOVED***
	return pipeAddress(f.path)
***REMOVED***

func (f *win32Pipe) RemoteAddr() net.Addr ***REMOVED***
	return pipeAddress(f.path)
***REMOVED***

func (f *win32Pipe) SetDeadline(t time.Time) error ***REMOVED***
	f.SetReadDeadline(t)
	f.SetWriteDeadline(t)
	return nil
***REMOVED***

// CloseWrite closes the write side of a message pipe in byte mode.
func (f *win32MessageBytePipe) CloseWrite() error ***REMOVED***
	if f.writeClosed ***REMOVED***
		return errPipeWriteClosed
	***REMOVED***
	err := f.win32File.Flush()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = f.win32File.Write(nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.writeClosed = true
	return nil
***REMOVED***

// Write writes bytes to a message pipe in byte mode. Zero-byte writes are ignored, since
// they are used to implement CloseWrite().
func (f *win32MessageBytePipe) Write(b []byte) (int, error) ***REMOVED***
	if f.writeClosed ***REMOVED***
		return 0, errPipeWriteClosed
	***REMOVED***
	if len(b) == 0 ***REMOVED***
		return 0, nil
	***REMOVED***
	return f.win32File.Write(b)
***REMOVED***

// Read reads bytes from a message pipe in byte mode. A read of a zero-byte message on a message
// mode pipe will return io.EOF, as will all subsequent reads.
func (f *win32MessageBytePipe) Read(b []byte) (int, error) ***REMOVED***
	if f.readEOF ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	n, err := f.win32File.Read(b)
	if err == io.EOF ***REMOVED***
		// If this was the result of a zero-byte read, then
		// it is possible that the read was due to a zero-size
		// message. Since we are simulating CloseWrite with a
		// zero-byte message, ensure that all future Read() calls
		// also return EOF.
		f.readEOF = true
	***REMOVED***
	return n, err
***REMOVED***

func (s pipeAddress) Network() string ***REMOVED***
	return "pipe"
***REMOVED***

func (s pipeAddress) String() string ***REMOVED***
	return string(s)
***REMOVED***

// DialPipe connects to a named pipe by path, timing out if the connection
// takes longer than the specified duration. If timeout is nil, then the timeout
// is the default timeout established by the pipe server.
func DialPipe(path string, timeout *time.Duration) (net.Conn, error) ***REMOVED***
	var absTimeout time.Time
	if timeout != nil ***REMOVED***
		absTimeout = time.Now().Add(*timeout)
	***REMOVED***
	var err error
	var h syscall.Handle
	for ***REMOVED***
		h, err = createFile(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, syscall.FILE_FLAG_OVERLAPPED|cSECURITY_SQOS_PRESENT|cSECURITY_ANONYMOUS, 0)
		if err != cERROR_PIPE_BUSY ***REMOVED***
			break
		***REMOVED***
		now := time.Now()
		var ms uint32
		if absTimeout.IsZero() ***REMOVED***
			ms = cNMPWAIT_USE_DEFAULT_WAIT
		***REMOVED*** else if now.After(absTimeout) ***REMOVED***
			ms = cNMPWAIT_NOWAIT
		***REMOVED*** else ***REMOVED***
			ms = uint32(absTimeout.Sub(now).Nanoseconds() / 1000 / 1000)
		***REMOVED***
		err = waitNamedPipe(path, ms)
		if err != nil ***REMOVED***
			if err == cERROR_SEM_TIMEOUT ***REMOVED***
				return nil, ErrTimeout
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "open", Path: path, Err: err***REMOVED***
	***REMOVED***

	var flags uint32
	err = getNamedPipeInfo(h, &flags, nil, nil, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var state uint32
	err = getNamedPipeHandleState(h, &state, nil, nil, nil, nil, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if state&cPIPE_READMODE_MESSAGE != 0 ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "open", Path: path, Err: errors.New("message readmode pipes not supported")***REMOVED***
	***REMOVED***

	f, err := makeWin32File(h)
	if err != nil ***REMOVED***
		syscall.Close(h)
		return nil, err
	***REMOVED***

	// If the pipe is in message mode, return a message byte pipe, which
	// supports CloseWrite().
	if flags&cPIPE_TYPE_MESSAGE != 0 ***REMOVED***
		return &win32MessageBytePipe***REMOVED***
			win32Pipe: win32Pipe***REMOVED***win32File: f, path: path***REMOVED***,
		***REMOVED***, nil
	***REMOVED***
	return &win32Pipe***REMOVED***win32File: f, path: path***REMOVED***, nil
***REMOVED***

type acceptResponse struct ***REMOVED***
	f   *win32File
	err error
***REMOVED***

type win32PipeListener struct ***REMOVED***
	firstHandle        syscall.Handle
	path               string
	securityDescriptor []byte
	config             PipeConfig
	acceptCh           chan (chan acceptResponse)
	closeCh            chan int
	doneCh             chan int
***REMOVED***

func makeServerPipeHandle(path string, securityDescriptor []byte, c *PipeConfig, first bool) (syscall.Handle, error) ***REMOVED***
	var flags uint32 = cPIPE_ACCESS_DUPLEX | syscall.FILE_FLAG_OVERLAPPED
	if first ***REMOVED***
		flags |= cFILE_FLAG_FIRST_PIPE_INSTANCE
	***REMOVED***

	var mode uint32 = cPIPE_REJECT_REMOTE_CLIENTS
	if c.MessageMode ***REMOVED***
		mode |= cPIPE_TYPE_MESSAGE
	***REMOVED***

	sa := &syscall.SecurityAttributes***REMOVED******REMOVED***
	sa.Length = uint32(unsafe.Sizeof(*sa))
	if securityDescriptor != nil ***REMOVED***
		len := uint32(len(securityDescriptor))
		sa.SecurityDescriptor = localAlloc(0, len)
		defer localFree(sa.SecurityDescriptor)
		copy((*[0xffff]byte)(unsafe.Pointer(sa.SecurityDescriptor))[:], securityDescriptor)
	***REMOVED***
	h, err := createNamedPipe(path, flags, mode, cPIPE_UNLIMITED_INSTANCES, uint32(c.OutputBufferSize), uint32(c.InputBufferSize), 0, sa)
	if err != nil ***REMOVED***
		return 0, &os.PathError***REMOVED***Op: "open", Path: path, Err: err***REMOVED***
	***REMOVED***
	return h, nil
***REMOVED***

func (l *win32PipeListener) makeServerPipe() (*win32File, error) ***REMOVED***
	h, err := makeServerPipeHandle(l.path, l.securityDescriptor, &l.config, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	f, err := makeWin32File(h)
	if err != nil ***REMOVED***
		syscall.Close(h)
		return nil, err
	***REMOVED***
	return f, nil
***REMOVED***

func (l *win32PipeListener) makeConnectedServerPipe() (*win32File, error) ***REMOVED***
	p, err := l.makeServerPipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Wait for the client to connect.
	ch := make(chan error)
	go func(p *win32File) ***REMOVED***
		ch <- connectPipe(p)
	***REMOVED***(p)

	select ***REMOVED***
	case err = <-ch:
		if err != nil ***REMOVED***
			p.Close()
			p = nil
		***REMOVED***
	case <-l.closeCh:
		// Abort the connect request by closing the handle.
		p.Close()
		p = nil
		err = <-ch
		if err == nil || err == ErrFileClosed ***REMOVED***
			err = ErrPipeListenerClosed
		***REMOVED***
	***REMOVED***
	return p, err
***REMOVED***

func (l *win32PipeListener) listenerRoutine() ***REMOVED***
	closed := false
	for !closed ***REMOVED***
		select ***REMOVED***
		case <-l.closeCh:
			closed = true
		case responseCh := <-l.acceptCh:
			var (
				p   *win32File
				err error
			)
			for ***REMOVED***
				p, err = l.makeConnectedServerPipe()
				// If the connection was immediately closed by the client, try
				// again.
				if err != cERROR_NO_DATA ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			responseCh <- acceptResponse***REMOVED***p, err***REMOVED***
			closed = err == ErrPipeListenerClosed
		***REMOVED***
	***REMOVED***
	syscall.Close(l.firstHandle)
	l.firstHandle = 0
	// Notify Close() and Accept() callers that the handle has been closed.
	close(l.doneCh)
***REMOVED***

// PipeConfig contain configuration for the pipe listener.
type PipeConfig struct ***REMOVED***
	// SecurityDescriptor contains a Windows security descriptor in SDDL format.
	SecurityDescriptor string

	// MessageMode determines whether the pipe is in byte or message mode. In either
	// case the pipe is read in byte mode by default. The only practical difference in
	// this implementation is that CloseWrite() is only supported for message mode pipes;
	// CloseWrite() is implemented as a zero-byte write, but zero-byte writes are only
	// transferred to the reader (and returned as io.EOF in this implementation)
	// when the pipe is in message mode.
	MessageMode bool

	// InputBufferSize specifies the size the input buffer, in bytes.
	InputBufferSize int32

	// OutputBufferSize specifies the size the input buffer, in bytes.
	OutputBufferSize int32
***REMOVED***

// ListenPipe creates a listener on a Windows named pipe path, e.g. \\.\pipe\mypipe.
// The pipe must not already exist.
func ListenPipe(path string, c *PipeConfig) (net.Listener, error) ***REMOVED***
	var (
		sd  []byte
		err error
	)
	if c == nil ***REMOVED***
		c = &PipeConfig***REMOVED******REMOVED***
	***REMOVED***
	if c.SecurityDescriptor != "" ***REMOVED***
		sd, err = SddlToSecurityDescriptor(c.SecurityDescriptor)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	h, err := makeServerPipeHandle(path, sd, c, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Immediately open and then close a client handle so that the named pipe is
	// created but not currently accepting connections.
	h2, err := createFile(path, 0, 0, nil, syscall.OPEN_EXISTING, cSECURITY_SQOS_PRESENT|cSECURITY_ANONYMOUS, 0)
	if err != nil ***REMOVED***
		syscall.Close(h)
		return nil, err
	***REMOVED***
	syscall.Close(h2)
	l := &win32PipeListener***REMOVED***
		firstHandle:        h,
		path:               path,
		securityDescriptor: sd,
		config:             *c,
		acceptCh:           make(chan (chan acceptResponse)),
		closeCh:            make(chan int),
		doneCh:             make(chan int),
	***REMOVED***
	go l.listenerRoutine()
	return l, nil
***REMOVED***

func connectPipe(p *win32File) error ***REMOVED***
	c, err := p.prepareIo()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer p.wg.Done()

	err = connectNamedPipe(p.handle, &c.o)
	_, err = p.asyncIo(c, nil, 0, err)
	if err != nil && err != cERROR_PIPE_CONNECTED ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (l *win32PipeListener) Accept() (net.Conn, error) ***REMOVED***
	ch := make(chan acceptResponse)
	select ***REMOVED***
	case l.acceptCh <- ch:
		response := <-ch
		err := response.err
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if l.config.MessageMode ***REMOVED***
			return &win32MessageBytePipe***REMOVED***
				win32Pipe: win32Pipe***REMOVED***win32File: response.f, path: l.path***REMOVED***,
			***REMOVED***, nil
		***REMOVED***
		return &win32Pipe***REMOVED***win32File: response.f, path: l.path***REMOVED***, nil
	case <-l.doneCh:
		return nil, ErrPipeListenerClosed
	***REMOVED***
***REMOVED***

func (l *win32PipeListener) Close() error ***REMOVED***
	select ***REMOVED***
	case l.closeCh <- 1:
		<-l.doneCh
	case <-l.doneCh:
	***REMOVED***
	return nil
***REMOVED***

func (l *win32PipeListener) Addr() net.Addr ***REMOVED***
	return pipeAddress(l.path)
***REMOVED***
