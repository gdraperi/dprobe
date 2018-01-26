package daemon

import (
	"context"

	"github.com/docker/docker/container"
	"github.com/docker/docker/libcontainerd"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// postRunProcessing starts a servicing container if required
func (daemon *Daemon) postRunProcessing(c *container.Container, ei libcontainerd.EventInfo) error ***REMOVED***
	if ei.ExitCode == 0 && ei.UpdatePending ***REMOVED***
		spec, err := daemon.createSpec(c)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Turn on servicing
		spec.Windows.Servicing = true

		copts, err := daemon.getLibcontainerdCreateOptions(c)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Create a new servicing container, which will start, complete the
		// update, and merge back the results if it succeeded, all as part of
		// the below function call.
		ctx := context.Background()
		svcID := c.ID + "_servicing"
		logger := logrus.WithField("container", svcID)
		if err := daemon.containerd.Create(ctx, svcID, spec, copts); err != nil ***REMOVED***
			c.SetExitCode(-1)
			return errors.Wrap(err, "post-run update servicing failed")
		***REMOVED***
		_, err = daemon.containerd.Start(ctx, svcID, "", false, nil)
		if err != nil ***REMOVED***
			logger.WithError(err).Warn("failed to run servicing container")
			if err := daemon.containerd.Delete(ctx, svcID); err != nil ***REMOVED***
				logger.WithError(err).Warn("failed to delete servicing container")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if _, _, err := daemon.containerd.DeleteTask(ctx, svcID); err != nil ***REMOVED***
				logger.WithError(err).Warn("failed to delete servicing container task")
			***REMOVED***
			if err := daemon.containerd.Delete(ctx, svcID); err != nil ***REMOVED***
				logger.WithError(err).Warn("failed to delete servicing container")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
