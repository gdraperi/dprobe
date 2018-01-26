package container

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	loggerCloseTimeout = 10 * time.Second
)

// Reset puts a container into a state where it can be restarted again.
func (container *Container) Reset(lock bool) ***REMOVED***
	if lock ***REMOVED***
		container.Lock()
		defer container.Unlock()
	***REMOVED***

	if err := container.CloseStreams(); err != nil ***REMOVED***
		logrus.Errorf("%s: %s", container.ID, err)
	***REMOVED***

	// Re-create a brand new stdin pipe once the container exited
	if container.Config.OpenStdin ***REMOVED***
		container.StreamConfig.NewInputPipes()
	***REMOVED***

	if container.LogDriver != nil ***REMOVED***
		if container.LogCopier != nil ***REMOVED***
			exit := make(chan struct***REMOVED******REMOVED***)
			go func() ***REMOVED***
				container.LogCopier.Wait()
				close(exit)
			***REMOVED***()
			select ***REMOVED***
			case <-time.After(loggerCloseTimeout):
				logrus.Warn("Logger didn't exit in time: logs may be truncated")
			case <-exit:
			***REMOVED***
		***REMOVED***
		container.LogDriver.Close()
		container.LogCopier = nil
		container.LogDriver = nil
	***REMOVED***
***REMOVED***
