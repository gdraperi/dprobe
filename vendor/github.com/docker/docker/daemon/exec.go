package daemon

import (
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/container"
	"github.com/docker/docker/container/stream"
	"github.com/docker/docker/daemon/exec"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Seconds to wait after sending TERM before trying KILL
const termProcessTimeout = 10

func (d *Daemon) registerExecCommand(container *container.Container, config *exec.Config) ***REMOVED***
	// Storing execs in container in order to kill them gracefully whenever the container is stopped or removed.
	container.ExecCommands.Add(config.ID, config)
	// Storing execs in daemon for easy access via Engine API.
	d.execCommands.Add(config.ID, config)
***REMOVED***

// ExecExists looks up the exec instance and returns a bool if it exists or not.
// It will also return the error produced by `getConfig`
func (d *Daemon) ExecExists(name string) (bool, error) ***REMOVED***
	if _, err := d.getExecConfig(name); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return true, nil
***REMOVED***

// getExecConfig looks up the exec instance by name. If the container associated
// with the exec instance is stopped or paused, it will return an error.
func (d *Daemon) getExecConfig(name string) (*exec.Config, error) ***REMOVED***
	ec := d.execCommands.Get(name)

	// If the exec is found but its container is not in the daemon's list of
	// containers then it must have been deleted, in which case instead of
	// saying the container isn't running, we should return a 404 so that
	// the user sees the same error now that they will after the
	// 5 minute clean-up loop is run which erases old/dead execs.

	if ec != nil ***REMOVED***
		if container := d.containers.Get(ec.ContainerID); container != nil ***REMOVED***
			if !container.IsRunning() ***REMOVED***
				return nil, fmt.Errorf("Container %s is not running: %s", container.ID, container.State.String())
			***REMOVED***
			if container.IsPaused() ***REMOVED***
				return nil, errExecPaused(container.ID)
			***REMOVED***
			if container.IsRestarting() ***REMOVED***
				return nil, errContainerIsRestarting(container.ID)
			***REMOVED***
			return ec, nil
		***REMOVED***
	***REMOVED***

	return nil, errExecNotFound(name)
***REMOVED***

func (d *Daemon) unregisterExecCommand(container *container.Container, execConfig *exec.Config) ***REMOVED***
	container.ExecCommands.Delete(execConfig.ID, execConfig.Pid)
	d.execCommands.Delete(execConfig.ID, execConfig.Pid)
***REMOVED***

func (d *Daemon) getActiveContainer(name string) (*container.Container, error) ***REMOVED***
	container, err := d.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !container.IsRunning() ***REMOVED***
		return nil, errNotRunning(container.ID)
	***REMOVED***
	if container.IsPaused() ***REMOVED***
		return nil, errExecPaused(name)
	***REMOVED***
	if container.IsRestarting() ***REMOVED***
		return nil, errContainerIsRestarting(container.ID)
	***REMOVED***
	return container, nil
***REMOVED***

// ContainerExecCreate sets up an exec in a running container.
func (d *Daemon) ContainerExecCreate(name string, config *types.ExecConfig) (string, error) ***REMOVED***
	cntr, err := d.getActiveContainer(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	cmd := strslice.StrSlice(config.Cmd)
	entrypoint, args := d.getEntrypointAndArgs(strslice.StrSlice***REMOVED******REMOVED***, cmd)

	keys := []byte***REMOVED******REMOVED***
	if config.DetachKeys != "" ***REMOVED***
		keys, err = term.ToBytes(config.DetachKeys)
		if err != nil ***REMOVED***
			err = fmt.Errorf("Invalid escape keys (%s) provided", config.DetachKeys)
			return "", err
		***REMOVED***
	***REMOVED***

	execConfig := exec.NewConfig()
	execConfig.OpenStdin = config.AttachStdin
	execConfig.OpenStdout = config.AttachStdout
	execConfig.OpenStderr = config.AttachStderr
	execConfig.ContainerID = cntr.ID
	execConfig.DetachKeys = keys
	execConfig.Entrypoint = entrypoint
	execConfig.Args = args
	execConfig.Tty = config.Tty
	execConfig.Privileged = config.Privileged
	execConfig.User = config.User
	execConfig.WorkingDir = config.WorkingDir

	linkedEnv, err := d.setupLinkedContainers(cntr)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	execConfig.Env = container.ReplaceOrAppendEnvValues(cntr.CreateDaemonEnvironment(config.Tty, linkedEnv), config.Env)
	if len(execConfig.User) == 0 ***REMOVED***
		execConfig.User = cntr.Config.User
	***REMOVED***
	if len(execConfig.WorkingDir) == 0 ***REMOVED***
		execConfig.WorkingDir = cntr.Config.WorkingDir
	***REMOVED***

	d.registerExecCommand(cntr, execConfig)

	attributes := map[string]string***REMOVED***
		"execID": execConfig.ID,
	***REMOVED***
	d.LogContainerEventWithAttributes(cntr, "exec_create: "+execConfig.Entrypoint+" "+strings.Join(execConfig.Args, " "), attributes)

	return execConfig.ID, nil
***REMOVED***

// ContainerExecStart starts a previously set up exec instance. The
// std streams are set up.
// If ctx is cancelled, the process is terminated.
func (d *Daemon) ContainerExecStart(ctx context.Context, name string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (err error) ***REMOVED***
	var (
		cStdin           io.ReadCloser
		cStdout, cStderr io.Writer
	)

	ec, err := d.getExecConfig(name)
	if err != nil ***REMOVED***
		return errExecNotFound(name)
	***REMOVED***

	ec.Lock()
	if ec.ExitCode != nil ***REMOVED***
		ec.Unlock()
		err := fmt.Errorf("Error: Exec command %s has already run", ec.ID)
		return errdefs.Conflict(err)
	***REMOVED***

	if ec.Running ***REMOVED***
		ec.Unlock()
		return errdefs.Conflict(fmt.Errorf("Error: Exec command %s is already running", ec.ID))
	***REMOVED***
	ec.Running = true
	ec.Unlock()

	c := d.containers.Get(ec.ContainerID)
	logrus.Debugf("starting exec command %s in container %s", ec.ID, c.ID)
	attributes := map[string]string***REMOVED***
		"execID": ec.ID,
	***REMOVED***
	d.LogContainerEventWithAttributes(c, "exec_start: "+ec.Entrypoint+" "+strings.Join(ec.Args, " "), attributes)

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			ec.Lock()
			ec.Running = false
			exitCode := 126
			ec.ExitCode = &exitCode
			if err := ec.CloseStreams(); err != nil ***REMOVED***
				logrus.Errorf("failed to cleanup exec %s streams: %s", c.ID, err)
			***REMOVED***
			ec.Unlock()
			c.ExecCommands.Delete(ec.ID, ec.Pid)
		***REMOVED***
	***REMOVED***()

	if ec.OpenStdin && stdin != nil ***REMOVED***
		r, w := io.Pipe()
		go func() ***REMOVED***
			defer w.Close()
			defer logrus.Debug("Closing buffered stdin pipe")
			pools.Copy(w, stdin)
		***REMOVED***()
		cStdin = r
	***REMOVED***
	if ec.OpenStdout ***REMOVED***
		cStdout = stdout
	***REMOVED***
	if ec.OpenStderr ***REMOVED***
		cStderr = stderr
	***REMOVED***

	if ec.OpenStdin ***REMOVED***
		ec.StreamConfig.NewInputPipes()
	***REMOVED*** else ***REMOVED***
		ec.StreamConfig.NewNopInputPipe()
	***REMOVED***

	p := &specs.Process***REMOVED***
		Args:     append([]string***REMOVED***ec.Entrypoint***REMOVED***, ec.Args...),
		Env:      ec.Env,
		Terminal: ec.Tty,
		Cwd:      ec.WorkingDir,
	***REMOVED***
	if p.Cwd == "" ***REMOVED***
		p.Cwd = "/"
	***REMOVED***

	if err := d.execSetPlatformOpt(c, ec, p); err != nil ***REMOVED***
		return err
	***REMOVED***

	attachConfig := stream.AttachConfig***REMOVED***
		TTY:        ec.Tty,
		UseStdin:   cStdin != nil,
		UseStdout:  cStdout != nil,
		UseStderr:  cStderr != nil,
		Stdin:      cStdin,
		Stdout:     cStdout,
		Stderr:     cStderr,
		DetachKeys: ec.DetachKeys,
		CloseStdin: true,
	***REMOVED***
	ec.StreamConfig.AttachStreams(&attachConfig)
	attachErr := ec.StreamConfig.CopyStreams(ctx, &attachConfig)

	// Synchronize with libcontainerd event loop
	ec.Lock()
	c.ExecCommands.Lock()
	systemPid, err := d.containerd.Exec(ctx, c.ID, ec.ID, p, cStdin != nil, ec.InitializeStdio)
	if err != nil ***REMOVED***
		c.ExecCommands.Unlock()
		ec.Unlock()
		return translateContainerdStartErr(ec.Entrypoint, ec.SetExitCode, err)
	***REMOVED***
	ec.Pid = systemPid
	c.ExecCommands.Unlock()
	ec.Unlock()

	select ***REMOVED***
	case <-ctx.Done():
		logrus.Debugf("Sending TERM signal to process %v in container %v", name, c.ID)
		d.containerd.SignalProcess(ctx, c.ID, name, int(signal.SignalMap["TERM"]))
		select ***REMOVED***
		case <-time.After(termProcessTimeout * time.Second):
			logrus.Infof("Container %v, process %v failed to exit within %d seconds of signal TERM - using the force", c.ID, name, termProcessTimeout)
			d.containerd.SignalProcess(ctx, c.ID, name, int(signal.SignalMap["KILL"]))
		case <-attachErr:
			// TERM signal worked
		***REMOVED***
		return fmt.Errorf("context cancelled")
	case err := <-attachErr:
		if err != nil ***REMOVED***
			if _, ok := err.(term.EscapeError); !ok ***REMOVED***
				return errdefs.System(errors.Wrap(err, "exec attach failed"))
			***REMOVED***
			attributes := map[string]string***REMOVED***
				"execID": ec.ID,
			***REMOVED***
			d.LogContainerEventWithAttributes(c, "exec_detach", attributes)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// execCommandGC runs a ticker to clean up the daemon references
// of exec configs that are no longer part of the container.
func (d *Daemon) execCommandGC() ***REMOVED***
	for range time.Tick(5 * time.Minute) ***REMOVED***
		var (
			cleaned          int
			liveExecCommands = d.containerExecIds()
		)
		for id, config := range d.execCommands.Commands() ***REMOVED***
			if config.CanRemove ***REMOVED***
				cleaned++
				d.execCommands.Delete(id, config.Pid)
			***REMOVED*** else ***REMOVED***
				if _, exists := liveExecCommands[id]; !exists ***REMOVED***
					config.CanRemove = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if cleaned > 0 ***REMOVED***
			logrus.Debugf("clean %d unused exec commands", cleaned)
		***REMOVED***
	***REMOVED***
***REMOVED***

// containerExecIds returns a list of all the current exec ids that are in use
// and running inside a container.
func (d *Daemon) containerExecIds() map[string]struct***REMOVED******REMOVED*** ***REMOVED***
	ids := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, c := range d.containers.List() ***REMOVED***
		for _, id := range c.ExecCommands.List() ***REMOVED***
			ids[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	return ids
***REMOVED***
