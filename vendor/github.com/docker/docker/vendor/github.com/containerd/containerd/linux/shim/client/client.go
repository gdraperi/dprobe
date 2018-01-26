// +build !windows

package client

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stevvooe/ttrpc"

	"github.com/containerd/containerd/events"
	"github.com/containerd/containerd/linux/shim"
	shimapi "github.com/containerd/containerd/linux/shim/v1"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/reaper"
	"github.com/containerd/containerd/sys"
	ptypes "github.com/gogo/protobuf/types"
)

var empty = &ptypes.Empty***REMOVED******REMOVED***

// Opt is an option for a shim client configuration
type Opt func(context.Context, shim.Config) (shimapi.ShimService, io.Closer, error)

// WithStart executes a new shim process
func WithStart(binary, address, daemonAddress, cgroup string, debug bool, exitHandler func()) Opt ***REMOVED***
	return func(ctx context.Context, config shim.Config) (_ shimapi.ShimService, _ io.Closer, err error) ***REMOVED***
		socket, err := newSocket(address)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		defer socket.Close()
		f, err := socket.File()
		if err != nil ***REMOVED***
			return nil, nil, errors.Wrapf(err, "failed to get fd for socket %s", address)
		***REMOVED***
		defer f.Close()

		cmd, err := newCommand(binary, daemonAddress, debug, config, f)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		ec, err := reaper.Default.Start(cmd)
		if err != nil ***REMOVED***
			return nil, nil, errors.Wrapf(err, "failed to start shim")
		***REMOVED***
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				cmd.Process.Kill()
			***REMOVED***
		***REMOVED***()
		go func() ***REMOVED***
			reaper.Default.Wait(cmd, ec)
			exitHandler()
		***REMOVED***()
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"pid":     cmd.Process.Pid,
			"address": address,
			"debug":   debug,
		***REMOVED***).Infof("shim %s started", binary)
		// set shim in cgroup if it is provided
		if cgroup != "" ***REMOVED***
			if err := setCgroup(cgroup, cmd); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			log.G(ctx).WithFields(logrus.Fields***REMOVED***
				"pid":     cmd.Process.Pid,
				"address": address,
			***REMOVED***).Infof("shim placed in cgroup %s", cgroup)
		***REMOVED***
		if err = sys.SetOOMScore(cmd.Process.Pid, sys.OOMScoreMaxKillable); err != nil ***REMOVED***
			return nil, nil, errors.Wrap(err, "failed to set OOM Score on shim")
		***REMOVED***
		c, clo, err := WithConnect(address)(ctx, config)
		if err != nil ***REMOVED***
			return nil, nil, errors.Wrap(err, "failed to connect")
		***REMOVED***
		return c, clo, nil
	***REMOVED***
***REMOVED***

func newCommand(binary, daemonAddress string, debug bool, config shim.Config, socket *os.File) (*exec.Cmd, error) ***REMOVED***
	selfExe, err := os.Executable()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	args := []string***REMOVED***
		"-namespace", config.Namespace,
		"-workdir", config.WorkDir,
		"-address", daemonAddress,
		"-containerd-binary", selfExe,
	***REMOVED***

	if config.Criu != "" ***REMOVED***
		args = append(args, "-criu-path", config.Criu)
	***REMOVED***
	if config.RuntimeRoot != "" ***REMOVED***
		args = append(args, "-runtime-root", config.RuntimeRoot)
	***REMOVED***
	if config.SystemdCgroup ***REMOVED***
		args = append(args, "-systemd-cgroup")
	***REMOVED***
	if debug ***REMOVED***
		args = append(args, "-debug")
	***REMOVED***

	cmd := exec.Command(binary, args...)
	cmd.Dir = config.Path
	// make sure the shim can be re-parented to system init
	// and is cloned in a new mount namespace because the overlay/filesystems
	// will be mounted by the shim
	cmd.SysProcAttr = getSysProcAttr()
	cmd.ExtraFiles = append(cmd.ExtraFiles, socket)
	if debug ***REMOVED***
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	***REMOVED***
	return cmd, nil
***REMOVED***

func newSocket(address string) (*net.UnixListener, error) ***REMOVED***
	if len(address) > 106 ***REMOVED***
		return nil, errors.Errorf("%q: unix socket path too long (limit 106)", address)
	***REMOVED***
	l, err := net.Listen("unix", "\x00"+address)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to listen to abstract unix socket %q", address)
	***REMOVED***

	return l.(*net.UnixListener), nil
***REMOVED***

func connect(address string, d func(string, time.Duration) (net.Conn, error)) (net.Conn, error) ***REMOVED***
	return d(address, 100*time.Second)
***REMOVED***

func annonDialer(address string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	address = strings.TrimPrefix(address, "unix://")
	return net.DialTimeout("unix", "\x00"+address, timeout)
***REMOVED***

// WithConnect connects to an existing shim
func WithConnect(address string) Opt ***REMOVED***
	return func(ctx context.Context, config shim.Config) (shimapi.ShimService, io.Closer, error) ***REMOVED***
		conn, err := connect(address, annonDialer)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		return shimapi.NewShimClient(ttrpc.NewClient(conn)), conn, nil
	***REMOVED***
***REMOVED***

// WithLocal uses an in process shim
func WithLocal(publisher events.Publisher) func(context.Context, shim.Config) (shimapi.ShimService, io.Closer, error) ***REMOVED***
	return func(ctx context.Context, config shim.Config) (shimapi.ShimService, io.Closer, error) ***REMOVED***
		service, err := shim.NewService(config, publisher)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		return shim.NewLocal(service), nil, nil
	***REMOVED***
***REMOVED***

// New returns a new shim client
func New(ctx context.Context, config shim.Config, opt Opt) (*Client, error) ***REMOVED***
	s, c, err := opt(ctx, config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Client***REMOVED***
		ShimService: s,
		c:           c,
		exitCh:      make(chan struct***REMOVED******REMOVED***),
	***REMOVED***, nil
***REMOVED***

// Client is a shim client containing the connection to a shim
type Client struct ***REMOVED***
	shimapi.ShimService

	c        io.Closer
	exitCh   chan struct***REMOVED******REMOVED***
	exitOnce sync.Once
***REMOVED***

// IsAlive returns true if the shim can be contacted.
// NOTE: a negative answer doesn't mean that the process is gone.
func (c *Client) IsAlive(ctx context.Context) (bool, error) ***REMOVED***
	_, err := c.ShimInfo(ctx, empty)
	if err != nil ***REMOVED***
		// TODO(stevvooe): There are some error conditions that need to be
		// handle with unix sockets existence to give the right answer here.
		return false, err
	***REMOVED***
	return true, nil
***REMOVED***

// StopShim signals the shim to exit and wait for the process to disappear
func (c *Client) StopShim(ctx context.Context) error ***REMOVED***
	return c.signalShim(ctx, unix.SIGTERM)
***REMOVED***

// KillShim kills the shim forcefully and wait for the process to disappear
func (c *Client) KillShim(ctx context.Context) error ***REMOVED***
	return c.signalShim(ctx, unix.SIGKILL)
***REMOVED***

// Close the cient connection
func (c *Client) Close() error ***REMOVED***
	if c.c == nil ***REMOVED***
		return nil
	***REMOVED***
	return c.c.Close()
***REMOVED***

func (c *Client) signalShim(ctx context.Context, sig syscall.Signal) error ***REMOVED***
	info, err := c.ShimInfo(ctx, empty)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pid := int(info.ShimPid)
	// make sure we don't kill ourselves if we are running a local shim
	if os.Getpid() == pid ***REMOVED***
		return nil
	***REMOVED***
	if err := unix.Kill(pid, sig); err != nil && err != unix.ESRCH ***REMOVED***
		return err
	***REMOVED***
	// wait for shim to die after being signaled
	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	case <-c.waitForExit(pid):
		return nil
	***REMOVED***
***REMOVED***

func (c *Client) waitForExit(pid int) <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	c.exitOnce.Do(func() ***REMOVED***
		for ***REMOVED***
			// use kill(pid, 0) here because the shim could have been reparented
			// and we are no longer able to waitpid(pid, ...) on the shim
			if err := unix.Kill(pid, 0); err == unix.ESRCH ***REMOVED***
				close(c.exitCh)
				return
			***REMOVED***
			time.Sleep(10 * time.Millisecond)
		***REMOVED***
	***REMOVED***)
	return c.exitCh
***REMOVED***
