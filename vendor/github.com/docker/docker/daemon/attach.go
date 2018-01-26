package daemon

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/container"
	"github.com/docker/docker/container/stream"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/term"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ContainerAttach attaches to logs according to the config passed in. See ContainerAttachConfig.
func (daemon *Daemon) ContainerAttach(prefixOrName string, c *backend.ContainerAttachConfig) error ***REMOVED***
	keys := []byte***REMOVED******REMOVED***
	var err error
	if c.DetachKeys != "" ***REMOVED***
		keys, err = term.ToBytes(c.DetachKeys)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(errors.Errorf("Invalid detach keys (%s) provided", c.DetachKeys))
		***REMOVED***
	***REMOVED***

	container, err := daemon.GetContainer(prefixOrName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if container.IsPaused() ***REMOVED***
		err := fmt.Errorf("container %s is paused, unpause the container before attach", prefixOrName)
		return errdefs.Conflict(err)
	***REMOVED***
	if container.IsRestarting() ***REMOVED***
		err := fmt.Errorf("container %s is restarting, wait until the container is running", prefixOrName)
		return errdefs.Conflict(err)
	***REMOVED***

	cfg := stream.AttachConfig***REMOVED***
		UseStdin:   c.UseStdin,
		UseStdout:  c.UseStdout,
		UseStderr:  c.UseStderr,
		TTY:        container.Config.Tty,
		CloseStdin: container.Config.StdinOnce,
		DetachKeys: keys,
	***REMOVED***
	container.StreamConfig.AttachStreams(&cfg)

	inStream, outStream, errStream, err := c.GetStreams()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer inStream.Close()

	if !container.Config.Tty && c.MuxStreams ***REMOVED***
		errStream = stdcopy.NewStdWriter(errStream, stdcopy.Stderr)
		outStream = stdcopy.NewStdWriter(outStream, stdcopy.Stdout)
	***REMOVED***

	if cfg.UseStdin ***REMOVED***
		cfg.Stdin = inStream
	***REMOVED***
	if cfg.UseStdout ***REMOVED***
		cfg.Stdout = outStream
	***REMOVED***
	if cfg.UseStderr ***REMOVED***
		cfg.Stderr = errStream
	***REMOVED***

	if err := daemon.containerAttach(container, &cfg, c.Logs, c.Stream); err != nil ***REMOVED***
		fmt.Fprintf(outStream, "Error attaching: %s\n", err)
	***REMOVED***
	return nil
***REMOVED***

// ContainerAttachRaw attaches the provided streams to the container's stdio
func (daemon *Daemon) ContainerAttachRaw(prefixOrName string, stdin io.ReadCloser, stdout, stderr io.Writer, doStream bool, attached chan struct***REMOVED******REMOVED***) error ***REMOVED***
	container, err := daemon.GetContainer(prefixOrName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cfg := stream.AttachConfig***REMOVED***
		UseStdin:   stdin != nil,
		UseStdout:  stdout != nil,
		UseStderr:  stderr != nil,
		TTY:        container.Config.Tty,
		CloseStdin: container.Config.StdinOnce,
	***REMOVED***
	container.StreamConfig.AttachStreams(&cfg)
	close(attached)
	if cfg.UseStdin ***REMOVED***
		cfg.Stdin = stdin
	***REMOVED***
	if cfg.UseStdout ***REMOVED***
		cfg.Stdout = stdout
	***REMOVED***
	if cfg.UseStderr ***REMOVED***
		cfg.Stderr = stderr
	***REMOVED***

	return daemon.containerAttach(container, &cfg, false, doStream)
***REMOVED***

func (daemon *Daemon) containerAttach(c *container.Container, cfg *stream.AttachConfig, logs, doStream bool) error ***REMOVED***
	if logs ***REMOVED***
		logDriver, logCreated, err := daemon.getLogger(c)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if logCreated ***REMOVED***
			defer func() ***REMOVED***
				if err = logDriver.Close(); err != nil ***REMOVED***
					logrus.Errorf("Error closing logger: %v", err)
				***REMOVED***
			***REMOVED***()
		***REMOVED***
		cLog, ok := logDriver.(logger.LogReader)
		if !ok ***REMOVED***
			return logger.ErrReadLogsNotSupported***REMOVED******REMOVED***
		***REMOVED***
		logs := cLog.ReadLogs(logger.ReadConfig***REMOVED***Tail: -1***REMOVED***)
		defer logs.Close()

	LogLoop:
		for ***REMOVED***
			select ***REMOVED***
			case msg, ok := <-logs.Msg:
				if !ok ***REMOVED***
					break LogLoop
				***REMOVED***
				if msg.Source == "stdout" && cfg.Stdout != nil ***REMOVED***
					cfg.Stdout.Write(msg.Line)
				***REMOVED***
				if msg.Source == "stderr" && cfg.Stderr != nil ***REMOVED***
					cfg.Stderr.Write(msg.Line)
				***REMOVED***
			case err := <-logs.Err:
				logrus.Errorf("Error streaming logs: %v", err)
				break LogLoop
			***REMOVED***
		***REMOVED***
	***REMOVED***

	daemon.LogContainerEvent(c, "attach")

	if !doStream ***REMOVED***
		return nil
	***REMOVED***

	if cfg.Stdin != nil ***REMOVED***
		r, w := io.Pipe()
		go func(stdin io.ReadCloser) ***REMOVED***
			defer w.Close()
			defer logrus.Debug("Closing buffered stdin pipe")
			io.Copy(w, stdin)
		***REMOVED***(cfg.Stdin)
		cfg.Stdin = r
	***REMOVED***

	if !c.Config.OpenStdin ***REMOVED***
		cfg.Stdin = nil
	***REMOVED***

	if c.Config.StdinOnce && !c.Config.Tty ***REMOVED***
		// Wait for the container to stop before returning.
		waitChan := c.Wait(context.Background(), container.WaitConditionNotRunning)
		defer func() ***REMOVED***
			<-waitChan // Ignore returned exit code.
		***REMOVED***()
	***REMOVED***

	ctx := c.InitAttachContext()
	err := <-c.StreamConfig.CopyStreams(ctx, cfg)
	if err != nil ***REMOVED***
		if _, ok := err.(term.EscapeError); ok ***REMOVED***
			daemon.LogContainerEvent(c, "detach")
		***REMOVED*** else ***REMOVED***
			logrus.Errorf("attach failed with error: %v", err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
