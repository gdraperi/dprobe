//+build !windows

package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
)

func (daemon *Daemon) saveApparmorConfig(container *container.Container) error ***REMOVED***
	container.AppArmorProfile = "" //we don't care about the previous value.

	if !daemon.apparmorEnabled ***REMOVED***
		return nil // if apparmor is disabled there is nothing to do here.
	***REMOVED***

	if err := parseSecurityOpt(container, container.HostConfig); err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	if !container.HostConfig.Privileged ***REMOVED***
		if container.AppArmorProfile == "" ***REMOVED***
			container.AppArmorProfile = defaultApparmorProfile
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		container.AppArmorProfile = "unconfined"
	***REMOVED***
	return nil
***REMOVED***
