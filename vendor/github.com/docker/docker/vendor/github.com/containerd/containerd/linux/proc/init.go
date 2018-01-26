// +build !windows

package proc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/containerd/console"
	"github.com/containerd/containerd/linux/runctypes"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/fifo"
	runc "github.com/containerd/go-runc"
	"github.com/containerd/typeurl"
	google_protobuf "github.com/gogo/protobuf/types"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// InitPidFile name of the file that contains the init pid
const InitPidFile = "init.pid"

// Init represents an initial process for a container
type Init struct ***REMOVED***
	wg sync.WaitGroup
	initState

	// mu is used to ensure that `Start()` and `Exited()` calls return in
	// the right order when invoked in separate go routines.
	// This is the case within the shim implementation as it makes use of
	// the reaper interface.
	mu sync.Mutex

	waitBlock chan struct***REMOVED******REMOVED***

	workDir string

	id       string
	bundle   string
	console  console.Console
	platform Platform
	io       runc.IO
	runtime  *runc.Runc
	status   int
	exited   time.Time
	pid      int
	closers  []io.Closer
	stdin    io.Closer
	stdio    Stdio
	rootfs   string
	IoUID    int
	IoGID    int
***REMOVED***

// NewRunc returns a new runc instance for a process
func NewRunc(root, path, namespace, runtime, criu string, systemd bool) *runc.Runc ***REMOVED***
	if root == "" ***REMOVED***
		root = RuncRoot
	***REMOVED***
	return &runc.Runc***REMOVED***
		Command:       runtime,
		Log:           filepath.Join(path, "log.json"),
		LogFormat:     runc.JSON,
		PdeathSignal:  syscall.SIGKILL,
		Root:          filepath.Join(root, namespace),
		Criu:          criu,
		SystemdCgroup: systemd,
	***REMOVED***
***REMOVED***

// New returns a new init process
func New(context context.Context, path, workDir, runtimeRoot, namespace, criu string, systemdCgroup bool, platform Platform, r *CreateConfig) (*Init, error) ***REMOVED***
	var success bool

	var options runctypes.CreateOptions
	if r.Options != nil ***REMOVED***
		v, err := typeurl.UnmarshalAny(r.Options)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		options = *v.(*runctypes.CreateOptions)
	***REMOVED***

	rootfs := filepath.Join(path, "rootfs")
	// count the number of successful mounts so we can undo
	// what was actually done rather than what should have been
	// done.
	defer func() ***REMOVED***
		if success ***REMOVED***
			return
		***REMOVED***
		if err2 := mount.UnmountAll(rootfs, 0); err2 != nil ***REMOVED***
			log.G(context).WithError(err2).Warn("Failed to cleanup rootfs mount")
		***REMOVED***
	***REMOVED***()
	for _, rm := range r.Rootfs ***REMOVED***
		m := &mount.Mount***REMOVED***
			Type:    rm.Type,
			Source:  rm.Source,
			Options: rm.Options,
		***REMOVED***
		if err := m.Mount(rootfs); err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to mount rootfs component %v", m)
		***REMOVED***
	***REMOVED***
	runtime := NewRunc(runtimeRoot, path, namespace, r.Runtime, criu, systemdCgroup)
	p := &Init***REMOVED***
		id:       r.ID,
		bundle:   r.Bundle,
		runtime:  runtime,
		platform: platform,
		stdio: Stdio***REMOVED***
			Stdin:    r.Stdin,
			Stdout:   r.Stdout,
			Stderr:   r.Stderr,
			Terminal: r.Terminal,
		***REMOVED***,
		rootfs:    rootfs,
		workDir:   workDir,
		status:    0,
		waitBlock: make(chan struct***REMOVED******REMOVED***),
		IoUID:     int(options.IoUid),
		IoGID:     int(options.IoGid),
	***REMOVED***
	p.initState = &createdState***REMOVED***p: p***REMOVED***
	var (
		err    error
		socket *runc.Socket
	)
	if r.Terminal ***REMOVED***
		if socket, err = runc.NewTempConsoleSocket(); err != nil ***REMOVED***
			return nil, errors.Wrap(err, "failed to create OCI runtime console socket")
		***REMOVED***
		defer socket.Close()
	***REMOVED*** else if hasNoIO(r) ***REMOVED***
		if p.io, err = runc.NewNullIO(); err != nil ***REMOVED***
			return nil, errors.Wrap(err, "creating new NULL IO")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if p.io, err = runc.NewPipeIO(int(options.IoUid), int(options.IoGid)); err != nil ***REMOVED***
			return nil, errors.Wrap(err, "failed to create OCI runtime io pipes")
		***REMOVED***
	***REMOVED***
	pidFile := filepath.Join(path, InitPidFile)
	if r.Checkpoint != "" ***REMOVED***
		opts := &runc.RestoreOpts***REMOVED***
			CheckpointOpts: runc.CheckpointOpts***REMOVED***
				ImagePath:  r.Checkpoint,
				WorkDir:    p.workDir,
				ParentPath: r.ParentCheckpoint,
			***REMOVED***,
			PidFile:     pidFile,
			IO:          p.io,
			NoPivot:     options.NoPivotRoot,
			Detach:      true,
			NoSubreaper: true,
		***REMOVED***
		p.initState = &createdCheckpointState***REMOVED***
			p:    p,
			opts: opts,
		***REMOVED***
		success = true
		return p, nil
	***REMOVED***
	opts := &runc.CreateOpts***REMOVED***
		PidFile:      pidFile,
		IO:           p.io,
		NoPivot:      options.NoPivotRoot,
		NoNewKeyring: options.NoNewKeyring,
	***REMOVED***
	if socket != nil ***REMOVED***
		opts.ConsoleSocket = socket
	***REMOVED***
	if err := p.runtime.Create(context, r.ID, r.Bundle, opts); err != nil ***REMOVED***
		return nil, p.runtimeError(err, "OCI runtime create failed")
	***REMOVED***
	if r.Stdin != "" ***REMOVED***
		sc, err := fifo.OpenFifo(context, r.Stdin, syscall.O_WRONLY|syscall.O_NONBLOCK, 0)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to open stdin fifo %s", r.Stdin)
		***REMOVED***
		p.stdin = sc
		p.closers = append(p.closers, sc)
	***REMOVED***
	var copyWaitGroup sync.WaitGroup
	if socket != nil ***REMOVED***
		console, err := socket.ReceiveMaster()
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "failed to retrieve console master")
		***REMOVED***
		console, err = platform.CopyConsole(context, console, r.Stdin, r.Stdout, r.Stderr, &p.wg, &copyWaitGroup)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "failed to start console copy")
		***REMOVED***
		p.console = console
	***REMOVED*** else if !hasNoIO(r) ***REMOVED***
		if err := copyPipes(context, p.io, r.Stdin, r.Stdout, r.Stderr, &p.wg, &copyWaitGroup); err != nil ***REMOVED***
			return nil, errors.Wrap(err, "failed to start io pipe copy")
		***REMOVED***
	***REMOVED***

	copyWaitGroup.Wait()
	pid, err := runc.ReadPidFile(pidFile)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to retrieve OCI runtime container pid")
	***REMOVED***
	p.pid = pid
	success = true
	return p, nil
***REMOVED***

// Wait for the process to exit
func (p *Init) Wait() ***REMOVED***
	<-p.waitBlock
***REMOVED***

// ID of the process
func (p *Init) ID() string ***REMOVED***
	return p.id
***REMOVED***

// Pid of the process
func (p *Init) Pid() int ***REMOVED***
	return p.pid
***REMOVED***

// ExitStatus of the process
func (p *Init) ExitStatus() int ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.status
***REMOVED***

// ExitedAt at time when the process exited
func (p *Init) ExitedAt() time.Time ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.exited
***REMOVED***

// Status of the process
func (p *Init) Status(ctx context.Context) (string, error) ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	c, err := p.runtime.State(ctx, p.id)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return "stopped", nil
		***REMOVED***
		return "", p.runtimeError(err, "OCI runtime state failed")
	***REMOVED***
	return c.Status, nil
***REMOVED***

func (p *Init) start(context context.Context) error ***REMOVED***
	err := p.runtime.Start(context, p.id)
	return p.runtimeError(err, "OCI runtime start failed")
***REMOVED***

func (p *Init) setExited(status int) ***REMOVED***
	p.exited = time.Now()
	p.status = status
	p.platform.ShutdownConsole(context.Background(), p.console)
	close(p.waitBlock)
***REMOVED***

func (p *Init) delete(context context.Context) error ***REMOVED***
	p.KillAll(context)
	p.wg.Wait()
	err := p.runtime.Delete(context, p.id, nil)
	// ignore errors if a runtime has already deleted the process
	// but we still hold metadata and pipes
	//
	// this is common during a checkpoint, runc will delete the container state
	// after a checkpoint and the container will no longer exist within runc
	if err != nil ***REMOVED***
		if strings.Contains(err.Error(), "does not exist") ***REMOVED***
			err = nil
		***REMOVED*** else ***REMOVED***
			err = p.runtimeError(err, "failed to delete task")
		***REMOVED***
	***REMOVED***
	if p.io != nil ***REMOVED***
		for _, c := range p.closers ***REMOVED***
			c.Close()
		***REMOVED***
		p.io.Close()
	***REMOVED***
	if err2 := mount.UnmountAll(p.rootfs, 0); err2 != nil ***REMOVED***
		log.G(context).WithError(err2).Warn("failed to cleanup rootfs mount")
		if err == nil ***REMOVED***
			err = errors.Wrap(err2, "failed rootfs umount")
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (p *Init) resize(ws console.WinSize) error ***REMOVED***
	if p.console == nil ***REMOVED***
		return nil
	***REMOVED***
	return p.console.Resize(ws)
***REMOVED***

func (p *Init) pause(context context.Context) error ***REMOVED***
	err := p.runtime.Pause(context, p.id)
	return p.runtimeError(err, "OCI runtime pause failed")
***REMOVED***

func (p *Init) resume(context context.Context) error ***REMOVED***
	err := p.runtime.Resume(context, p.id)
	return p.runtimeError(err, "OCI runtime resume failed")
***REMOVED***

func (p *Init) kill(context context.Context, signal uint32, all bool) error ***REMOVED***
	err := p.runtime.Kill(context, p.id, int(signal), &runc.KillOpts***REMOVED***
		All: all,
	***REMOVED***)
	return checkKillError(err)
***REMOVED***

// KillAll processes belonging to the init process
func (p *Init) KillAll(context context.Context) error ***REMOVED***
	err := p.runtime.Kill(context, p.id, int(syscall.SIGKILL), &runc.KillOpts***REMOVED***
		All: true,
	***REMOVED***)
	return p.runtimeError(err, "OCI runtime killall failed")
***REMOVED***

// Stdin of the process
func (p *Init) Stdin() io.Closer ***REMOVED***
	return p.stdin
***REMOVED***

// Runtime returns the OCI runtime configured for the init process
func (p *Init) Runtime() *runc.Runc ***REMOVED***
	return p.runtime
***REMOVED***

// Exec returns a new exec'd process
func (p *Init) Exec(context context.Context, path string, r *ExecConfig) (Process, error) ***REMOVED***
	// process exec request
	var spec specs.Process
	if err := json.Unmarshal(r.Spec.Value, &spec); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	spec.Terminal = r.Terminal

	e := &execProcess***REMOVED***
		id:     r.ID,
		path:   path,
		parent: p,
		spec:   spec,
		stdio: Stdio***REMOVED***
			Stdin:    r.Stdin,
			Stdout:   r.Stdout,
			Stderr:   r.Stderr,
			Terminal: r.Terminal,
		***REMOVED***,
		waitBlock: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	e.State = &execCreatedState***REMOVED***p: e***REMOVED***
	return e, nil
***REMOVED***

func (p *Init) checkpoint(context context.Context, r *CheckpointConfig) error ***REMOVED***
	var options runctypes.CheckpointOptions
	if r.Options != nil ***REMOVED***
		v, err := typeurl.UnmarshalAny(r.Options)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		options = *v.(*runctypes.CheckpointOptions)
	***REMOVED***
	var actions []runc.CheckpointAction
	if !options.Exit ***REMOVED***
		actions = append(actions, runc.LeaveRunning)
	***REMOVED***
	work := filepath.Join(p.workDir, "criu-work")
	defer os.RemoveAll(work)
	if err := p.runtime.Checkpoint(context, p.id, &runc.CheckpointOpts***REMOVED***
		WorkDir:                  work,
		ImagePath:                r.Path,
		AllowOpenTCP:             options.OpenTcp,
		AllowExternalUnixSockets: options.ExternalUnixSockets,
		AllowTerminal:            options.Terminal,
		FileLocks:                options.FileLocks,
		EmptyNamespaces:          options.EmptyNamespaces,
	***REMOVED***, actions...); err != nil ***REMOVED***
		dumpLog := filepath.Join(p.bundle, "criu-dump.log")
		if cerr := copyFile(dumpLog, filepath.Join(work, "dump.log")); cerr != nil ***REMOVED***
			log.G(context).Error(err)
		***REMOVED***
		return fmt.Errorf("%s path= %s", criuError(err), dumpLog)
	***REMOVED***
	return nil
***REMOVED***

func (p *Init) update(context context.Context, r *google_protobuf.Any) error ***REMOVED***
	var resources specs.LinuxResources
	if err := json.Unmarshal(r.Value, &resources); err != nil ***REMOVED***
		return err
	***REMOVED***
	return p.runtime.Update(context, p.id, &resources)
***REMOVED***

// Stdio of the process
func (p *Init) Stdio() Stdio ***REMOVED***
	return p.stdio
***REMOVED***

func (p *Init) runtimeError(rErr error, msg string) error ***REMOVED***
	if rErr == nil ***REMOVED***
		return nil
	***REMOVED***

	rMsg, err := getLastRuntimeError(p.runtime)
	switch ***REMOVED***
	case err != nil:
		return errors.Wrapf(rErr, "%s: %s (%s)", msg, "unable to retrieve OCI runtime error", err.Error())
	case rMsg == "":
		return errors.Wrap(rErr, msg)
	default:
		return errors.Errorf("%s: %s", msg, rMsg)
	***REMOVED***
***REMOVED***
