package ttrpc

import (
	"context"
	"net"
	"os"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type UnixCredentialsFunc func(*unix.Ucred) error

func (fn UnixCredentialsFunc) Handshake(ctx context.Context, conn net.Conn) (net.Conn, interface***REMOVED******REMOVED***, error) ***REMOVED***
	uc, err := requireUnixSocket(conn)
	if err != nil ***REMOVED***
		return nil, nil, errors.Wrap(err, "ttrpc.UnixCredentialsFunc: require unix socket")
	***REMOVED***

	rs, err := uc.SyscallConn()
	if err != nil ***REMOVED***
		return nil, nil, errors.Wrap(err, "ttrpc.UnixCredentialsFunc: (net.UnixConn).SyscallConn failed")
	***REMOVED***
	var (
		ucred    *unix.Ucred
		ucredErr error
	)
	if err := rs.Control(func(fd uintptr) ***REMOVED***
		ucred, ucredErr = unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
	***REMOVED***); err != nil ***REMOVED***
		return nil, nil, errors.Wrapf(err, "ttrpc.UnixCredentialsFunc: (*syscall.RawConn).Control failed")
	***REMOVED***

	if ucredErr != nil ***REMOVED***
		return nil, nil, errors.Wrapf(err, "ttrpc.UnixCredentialsFunc: failed to retrieve socket peer credentials")
	***REMOVED***

	if err := fn(ucred); err != nil ***REMOVED***
		return nil, nil, errors.Wrapf(err, "ttrpc.UnixCredentialsFunc: credential check failed")
	***REMOVED***

	return uc, ucred, nil
***REMOVED***

// UnixSocketRequireUidGid requires specific *effective* UID/GID, rather than the real UID/GID.
//
// For example, if a daemon binary is owned by the root (UID 0) with SUID bit but running as an
// unprivileged user (UID 1001), the effective UID becomes 0, and the real UID becomes 1001.
// So calling this function with uid=0 allows a connection from effective UID 0 but rejects
// a connection from effective UID 1001.
//
// See socket(7), SO_PEERCRED: "The returned credentials are those that were in effect at the time of the call to connect(2) or socketpair(2)."
func UnixSocketRequireUidGid(uid, gid int) UnixCredentialsFunc ***REMOVED***
	return func(ucred *unix.Ucred) error ***REMOVED***
		return requireUidGid(ucred, uid, gid)
	***REMOVED***
***REMOVED***

func UnixSocketRequireRoot() UnixCredentialsFunc ***REMOVED***
	return UnixSocketRequireUidGid(0, 0)
***REMOVED***

// UnixSocketRequireSameUser resolves the current effective unix user and returns a
// UnixCredentialsFunc that will validate incoming unix connections against the
// current credentials.
//
// This is useful when using abstract sockets that are accessible by all users.
func UnixSocketRequireSameUser() UnixCredentialsFunc ***REMOVED***
	euid, egid := os.Geteuid(), os.Getegid()
	return UnixSocketRequireUidGid(euid, egid)
***REMOVED***

func requireRoot(ucred *unix.Ucred) error ***REMOVED***
	return requireUidGid(ucred, 0, 0)
***REMOVED***

func requireUidGid(ucred *unix.Ucred, uid, gid int) error ***REMOVED***
	if (uid != -1 && uint32(uid) != ucred.Uid) || (gid != -1 && uint32(gid) != ucred.Gid) ***REMOVED***
		return errors.Wrap(syscall.EPERM, "ttrpc: invalid credentials")
	***REMOVED***
	return nil
***REMOVED***

func requireUnixSocket(conn net.Conn) (*net.UnixConn, error) ***REMOVED***
	uc, ok := conn.(*net.UnixConn)
	if !ok ***REMOVED***
		return nil, errors.New("a unix socket connection is required")
	***REMOVED***

	return uc, nil
***REMOVED***
