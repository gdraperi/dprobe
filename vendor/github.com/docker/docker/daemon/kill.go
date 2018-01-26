package daemon

import (
	"context"
	"fmt"
	"runtime"
	"syscall"
	"time"

	containerpkg "github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/libcontainerd"
	"github.com/docker/docker/pkg/signal"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type errNoSuchProcess struct ***REMOVED***
	pid    int
	signal int
***REMOVED***

func (e errNoSuchProcess) Error() string ***REMOVED***
	return fmt.Sprintf("Cannot kill process (pid=%d) with signal %d: no such process.", e.pid, e.signal)
***REMOVED***

func (errNoSuchProcess) NotFound() ***REMOVED******REMOVED***

// isErrNoSuchProcess returns true if the error
// is an instance of errNoSuchProcess.
func isErrNoSuchProcess(err error) bool ***REMOVED***
	_, ok := err.(errNoSuchProcess)
	return ok
***REMOVED***

// ContainerKill sends signal to the container
// If no signal is given (sig 0), then Kill with SIGKILL and wait
// for the container to exit.
// If a signal is given, then just send it to the container and return.
func (daemon *Daemon) ContainerKill(name string, sig uint64) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if sig != 0 && !signal.ValidSignalForPlatform(syscall.Signal(sig)) ***REMOVED***
		return fmt.Errorf("The %s daemon does not support signal %d", runtime.GOOS, sig)
	***REMOVED***

	// If no signal is passed, or SIGKILL, perform regular Kill (SIGKILL + wait())
	if sig == 0 || syscall.Signal(sig) == syscall.SIGKILL ***REMOVED***
		return daemon.Kill(container)
	***REMOVED***
	return daemon.killWithSignal(container, int(sig))
***REMOVED***

// killWithSignal sends the container the given signal. This wrapper for the
// host specific kill command prepares the container before attempting
// to send the signal. An error is returned if the container is paused
// or not running, or if there is a problem returned from the
// underlying kill command.
func (daemon *Daemon) killWithSignal(container *containerpkg.Container, sig int) error ***REMOVED***
	logrus.Debugf("Sending kill signal %d to container %s", sig, container.ID)
	container.Lock()
	defer container.Unlock()

	daemon.stopHealthchecks(container)

	if !container.Running ***REMOVED***
		return errNotRunning(container.ID)
	***REMOVED***

	var unpause bool
	if container.Config.StopSignal != "" && syscall.Signal(sig) != syscall.SIGKILL ***REMOVED***
		containerStopSignal, err := signal.ParseSignal(container.Config.StopSignal)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if containerStopSignal == syscall.Signal(sig) ***REMOVED***
			container.ExitOnNext()
			unpause = container.Paused
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		container.ExitOnNext()
		unpause = container.Paused
	***REMOVED***

	if !daemon.IsShuttingDown() ***REMOVED***
		container.HasBeenManuallyStopped = true
	***REMOVED***

	// if the container is currently restarting we do not need to send the signal
	// to the process. Telling the monitor that it should exit on its next event
	// loop is enough
	if container.Restarting ***REMOVED***
		return nil
	***REMOVED***

	if err := daemon.kill(container, sig); err != nil ***REMOVED***
		if errdefs.IsNotFound(err) ***REMOVED***
			unpause = false
			logrus.WithError(err).WithField("container", container.ID).WithField("action", "kill").Debug("container kill failed because of 'container not found' or 'no such process'")
		***REMOVED*** else ***REMOVED***
			return errors.Wrapf(err, "Cannot kill container %s", container.ID)
		***REMOVED***
	***REMOVED***

	if unpause ***REMOVED***
		// above kill signal will be sent once resume is finished
		if err := daemon.containerd.Resume(context.Background(), container.ID); err != nil ***REMOVED***
			logrus.Warn("Cannot unpause container %s: %s", container.ID, err)
		***REMOVED***
	***REMOVED***

	attributes := map[string]string***REMOVED***
		"signal": fmt.Sprintf("%d", sig),
	***REMOVED***
	daemon.LogContainerEventWithAttributes(container, "kill", attributes)
	return nil
***REMOVED***

// Kill forcefully terminates a container.
func (daemon *Daemon) Kill(container *containerpkg.Container) error ***REMOVED***
	if !container.IsRunning() ***REMOVED***
		return errNotRunning(container.ID)
	***REMOVED***

	// 1. Send SIGKILL
	if err := daemon.killPossiblyDeadProcess(container, int(syscall.SIGKILL)); err != nil ***REMOVED***
		// While normally we might "return err" here we're not going to
		// because if we can't stop the container by this point then
		// it's probably because it's already stopped. Meaning, between
		// the time of the IsRunning() call above and now it stopped.
		// Also, since the err return will be environment specific we can't
		// look for any particular (common) error that would indicate
		// that the process is already dead vs something else going wrong.
		// So, instead we'll give it up to 2 more seconds to complete and if
		// by that time the container is still running, then the error
		// we got is probably valid and so we return it to the caller.
		if isErrNoSuchProcess(err) ***REMOVED***
			return nil
		***REMOVED***

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if status := <-container.Wait(ctx, containerpkg.WaitConditionNotRunning); status.Err() != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// 2. Wait for the process to die, in last resort, try to kill the process directly
	if err := killProcessDirectly(container); err != nil ***REMOVED***
		if isErrNoSuchProcess(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	// Wait for exit with no timeout.
	// Ignore returned status.
	<-container.Wait(context.Background(), containerpkg.WaitConditionNotRunning)

	return nil
***REMOVED***

// killPossibleDeadProcess is a wrapper around killSig() suppressing "no such process" error.
func (daemon *Daemon) killPossiblyDeadProcess(container *containerpkg.Container, sig int) error ***REMOVED***
	err := daemon.killWithSignal(container, sig)
	if errdefs.IsNotFound(err) ***REMOVED***
		e := errNoSuchProcess***REMOVED***container.GetPID(), sig***REMOVED***
		logrus.Debug(e)
		return e
	***REMOVED***
	return err
***REMOVED***

func (daemon *Daemon) kill(c *containerpkg.Container, sig int) error ***REMOVED***
	return daemon.containerd.SignalProcess(context.Background(), c.ID, libcontainerd.InitProcessName, sig)
***REMOVED***
