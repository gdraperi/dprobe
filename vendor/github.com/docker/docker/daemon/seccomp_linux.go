// +build linux,seccomp

package daemon

import (
	"fmt"

	"github.com/docker/docker/container"
	"github.com/docker/docker/profiles/seccomp"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

var supportsSeccomp = true

func setSeccomp(daemon *Daemon, rs *specs.Spec, c *container.Container) error ***REMOVED***
	var profile *specs.LinuxSeccomp
	var err error

	if c.HostConfig.Privileged ***REMOVED***
		return nil
	***REMOVED***

	if !daemon.seccompEnabled ***REMOVED***
		if c.SeccompProfile != "" && c.SeccompProfile != "unconfined" ***REMOVED***
			return fmt.Errorf("Seccomp is not enabled in your kernel, cannot run a custom seccomp profile.")
		***REMOVED***
		logrus.Warn("Seccomp is not enabled in your kernel, running container without default profile.")
		c.SeccompProfile = "unconfined"
	***REMOVED***
	if c.SeccompProfile == "unconfined" ***REMOVED***
		return nil
	***REMOVED***
	if c.SeccompProfile != "" ***REMOVED***
		profile, err = seccomp.LoadProfile(c.SeccompProfile, rs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if daemon.seccompProfile != nil ***REMOVED***
			profile, err = seccomp.LoadProfile(string(daemon.seccompProfile), rs)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			profile, err = seccomp.GetDefaultProfile(rs)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	rs.Linux.Seccomp = profile
	return nil
***REMOVED***
