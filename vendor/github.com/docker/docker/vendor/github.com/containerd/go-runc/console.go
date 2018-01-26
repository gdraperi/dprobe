package runc

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/containerd/console"
	"golang.org/x/sys/unix"
)

// NewConsoleSocket creates a new unix socket at the provided path to accept a
// pty master created by runc for use by the container
func NewConsoleSocket(path string) (*Socket, error) ***REMOVED***
	abs, err := filepath.Abs(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	addr, err := net.ResolveUnixAddr("unix", abs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	l, err := net.ListenUnix("unix", addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Socket***REMOVED***
		l:    l,
	***REMOVED***, nil
***REMOVED***

// NewTempConsoleSocket returns a temp console socket for use with a container
// On Close(), the socket is deleted
func NewTempConsoleSocket() (*Socket, error) ***REMOVED***
	dir, err := ioutil.TempDir("", "pty")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	abs, err := filepath.Abs(filepath.Join(dir, "pty.sock"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	addr, err := net.ResolveUnixAddr("unix", abs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	l, err := net.ListenUnix("unix", addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Socket***REMOVED***
		l:     l,
		rmdir: true,
	***REMOVED***, nil
***REMOVED***

// Socket is a unix socket that accepts the pty master created by runc
type Socket struct ***REMOVED***
	rmdir bool
	l     *net.UnixListener
***REMOVED***

// Path returns the path to the unix socket on disk
func (c *Socket) Path() string ***REMOVED***
	return c.l.Addr().String()
***REMOVED***

// recvFd waits for a file descriptor to be sent over the given AF_UNIX
// socket. The file name of the remote file descriptor will be recreated
// locally (it is sent as non-auxiliary data in the same payload).
func recvFd(socket *net.UnixConn) (*os.File, error) ***REMOVED***
	const MaxNameLen = 4096
	var oobSpace = unix.CmsgSpace(4)

	name := make([]byte, MaxNameLen)
	oob := make([]byte, oobSpace)

	n, oobn, _, _, err := socket.ReadMsgUnix(name, oob)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if n >= MaxNameLen || oobn != oobSpace ***REMOVED***
		return nil, fmt.Errorf("recvfd: incorrect number of bytes read (n=%d oobn=%d)", n, oobn)
	***REMOVED***

	// Truncate.
	name = name[:n]
	oob = oob[:oobn]

	scms, err := unix.ParseSocketControlMessage(oob)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(scms) != 1 ***REMOVED***
		return nil, fmt.Errorf("recvfd: number of SCMs is not 1: %d", len(scms))
	***REMOVED***
	scm := scms[0]

	fds, err := unix.ParseUnixRights(&scm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(fds) != 1 ***REMOVED***
		return nil, fmt.Errorf("recvfd: number of fds is not 1: %d", len(fds))
	***REMOVED***
	fd := uintptr(fds[0])

	return os.NewFile(fd, string(name)), nil
***REMOVED***

// ReceiveMaster blocks until the socket receives the pty master
func (c *Socket) ReceiveMaster() (console.Console, error) ***REMOVED***
	conn, err := c.l.Accept()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer conn.Close()
	uc, ok := conn.(*net.UnixConn)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("received connection which was not a unix socket")
	***REMOVED***
	f, err := recvFd(uc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return console.ConsoleFromFile(f)
***REMOVED***

// Close closes the unix socket
func (c *Socket) Close() error ***REMOVED***
	err := c.l.Close()
	if c.rmdir ***REMOVED***
		if rerr := os.RemoveAll(filepath.Dir(c.Path())); err == nil ***REMOVED***
			err = rerr
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***
