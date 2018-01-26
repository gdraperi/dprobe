package daemon

import (
	"errors"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	containertypes "github.com/docker/docker/api/types/container"
	timetypes "github.com/docker/docker/api/types/time"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/errdefs"
	"github.com/sirupsen/logrus"
)

// ContainerLogs copies the container's log channel to the channel provided in
// the config. If ContainerLogs returns an error, no messages have been copied.
// and the channel will be closed without data.
//
// if it returns nil, the config channel will be active and return log
// messages until it runs out or the context is canceled.
func (daemon *Daemon) ContainerLogs(ctx context.Context, containerName string, config *types.ContainerLogsOptions) (<-chan *backend.LogMessage, bool, error) ***REMOVED***
	lg := logrus.WithFields(logrus.Fields***REMOVED***
		"module":    "daemon",
		"method":    "(*Daemon).ContainerLogs",
		"container": containerName,
	***REMOVED***)

	if !(config.ShowStdout || config.ShowStderr) ***REMOVED***
		return nil, false, errdefs.InvalidParameter(errors.New("You must choose at least one stream"))
	***REMOVED***
	container, err := daemon.GetContainer(containerName)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***

	if container.RemovalInProgress || container.Dead ***REMOVED***
		return nil, false, errdefs.Conflict(errors.New("can not get logs from container which is dead or marked for removal"))
	***REMOVED***

	if container.HostConfig.LogConfig.Type == "none" ***REMOVED***
		return nil, false, logger.ErrReadLogsNotSupported***REMOVED******REMOVED***
	***REMOVED***

	cLog, cLogCreated, err := daemon.getLogger(container)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	if cLogCreated ***REMOVED***
		defer func() ***REMOVED***
			if err = cLog.Close(); err != nil ***REMOVED***
				logrus.Errorf("Error closing logger: %v", err)
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	logReader, ok := cLog.(logger.LogReader)
	if !ok ***REMOVED***
		return nil, false, logger.ErrReadLogsNotSupported***REMOVED******REMOVED***
	***REMOVED***

	follow := config.Follow && !cLogCreated
	tailLines, err := strconv.Atoi(config.Tail)
	if err != nil ***REMOVED***
		tailLines = -1
	***REMOVED***

	var since time.Time
	if config.Since != "" ***REMOVED***
		s, n, err := timetypes.ParseTimestamps(config.Since, 0)
		if err != nil ***REMOVED***
			return nil, false, err
		***REMOVED***
		since = time.Unix(s, n)
	***REMOVED***

	var until time.Time
	if config.Until != "" && config.Until != "0" ***REMOVED***
		s, n, err := timetypes.ParseTimestamps(config.Until, 0)
		if err != nil ***REMOVED***
			return nil, false, err
		***REMOVED***
		until = time.Unix(s, n)
	***REMOVED***

	readConfig := logger.ReadConfig***REMOVED***
		Since:  since,
		Until:  until,
		Tail:   tailLines,
		Follow: follow,
	***REMOVED***

	logs := logReader.ReadLogs(readConfig)

	// past this point, we can't possibly return any errors, so we can just
	// start a goroutine and return to tell the caller not to expect errors
	// (if the caller wants to give up on logs, they have to cancel the context)
	// this goroutine functions as a shim between the logger and the caller.
	messageChan := make(chan *backend.LogMessage, 1)
	go func() ***REMOVED***
		// set up some defers
		defer logs.Close()

		// close the messages channel. closing is the only way to signal above
		// that we're doing with logs (other than context cancel i guess).
		defer close(messageChan)

		lg.Debug("begin logs")
		for ***REMOVED***
			select ***REMOVED***
			// i do not believe as the system is currently designed any error
			// is possible, but we should be prepared to handle it anyway. if
			// we do get an error, copy only the error field to a new object so
			// we don't end up with partial data in the other fields
			case err := <-logs.Err:
				lg.Errorf("Error streaming logs: %v", err)
				select ***REMOVED***
				case <-ctx.Done():
				case messageChan <- &backend.LogMessage***REMOVED***Err: err***REMOVED***:
				***REMOVED***
				return
			case <-ctx.Done():
				lg.Debugf("logs: end stream, ctx is done: %v", ctx.Err())
				return
			case msg, ok := <-logs.Msg:
				// there is some kind of pool or ring buffer in the logger that
				// produces these messages, and a possible future optimization
				// might be to use that pool and reuse message objects
				if !ok ***REMOVED***
					lg.Debug("end logs")
					return
				***REMOVED***
				m := msg.AsLogMessage() // just a pointer conversion, does not copy data

				// there could be a case where the reader stops accepting
				// messages and the context is canceled. we need to check that
				// here, or otherwise we risk blocking forever on the message
				// send.
				select ***REMOVED***
				case <-ctx.Done():
					return
				case messageChan <- m:
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return messageChan, container.Config.Tty, nil
***REMOVED***

func (daemon *Daemon) getLogger(container *container.Container) (l logger.Logger, created bool, err error) ***REMOVED***
	container.Lock()
	if container.State.Running ***REMOVED***
		l = container.LogDriver
	***REMOVED***
	container.Unlock()
	if l == nil ***REMOVED***
		created = true
		l, err = container.StartLogger()
	***REMOVED***
	return
***REMOVED***

// mergeLogConfig merges the daemon log config to the container's log config if the container's log driver is not specified.
func (daemon *Daemon) mergeAndVerifyLogConfig(cfg *containertypes.LogConfig) error ***REMOVED***
	if cfg.Type == "" ***REMOVED***
		cfg.Type = daemon.defaultLogConfig.Type
	***REMOVED***

	if cfg.Config == nil ***REMOVED***
		cfg.Config = make(map[string]string)
	***REMOVED***

	if cfg.Type == daemon.defaultLogConfig.Type ***REMOVED***
		for k, v := range daemon.defaultLogConfig.Config ***REMOVED***
			if _, ok := cfg.Config[k]; !ok ***REMOVED***
				cfg.Config[k] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return logger.ValidateLogOpts(cfg.Type, cfg.Config)
***REMOVED***
