package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/caps"
	"github.com/docker/docker/daemon/exec"
	"github.com/opencontainers/runc/libcontainer/apparmor"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func (daemon *Daemon) execSetPlatformOpt(c *container.Container, ec *exec.Config, p *specs.Process) error ***REMOVED***
	if len(ec.User) > 0 ***REMOVED***
		uid, gid, additionalGids, err := getUser(c, ec.User)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		p.User = specs.User***REMOVED***
			UID:            uid,
			GID:            gid,
			AdditionalGids: additionalGids,
		***REMOVED***
	***REMOVED***
	if ec.Privileged ***REMOVED***
		if p.Capabilities == nil ***REMOVED***
			p.Capabilities = &specs.LinuxCapabilities***REMOVED******REMOVED***
		***REMOVED***
		p.Capabilities.Bounding = caps.GetAllCapabilities()
		p.Capabilities.Permitted = p.Capabilities.Bounding
		p.Capabilities.Inheritable = p.Capabilities.Bounding
		p.Capabilities.Effective = p.Capabilities.Bounding
	***REMOVED***
	if apparmor.IsEnabled() ***REMOVED***
		var appArmorProfile string
		if c.AppArmorProfile != "" ***REMOVED***
			appArmorProfile = c.AppArmorProfile
		***REMOVED*** else if c.HostConfig.Privileged ***REMOVED***
			appArmorProfile = "unconfined"
		***REMOVED*** else ***REMOVED***
			appArmorProfile = "docker-default"
		***REMOVED***

		if appArmorProfile == "docker-default" ***REMOVED***
			// Unattended upgrades and other fun services can unload AppArmor
			// profiles inadvertently. Since we cannot store our profile in
			// /etc/apparmor.d, nor can we practically add other ways of
			// telling the system to keep our profile loaded, in order to make
			// sure that we keep the default profile enabled we dynamically
			// reload it if necessary.
			if err := ensureDefaultAppArmorProfile(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	daemon.setRlimits(&specs.Spec***REMOVED***Process: p***REMOVED***, c)
	return nil
***REMOVED***
