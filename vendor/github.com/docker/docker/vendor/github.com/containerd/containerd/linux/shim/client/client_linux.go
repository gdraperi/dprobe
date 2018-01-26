// +build linux

package client

import (
	"os/exec"
	"syscall"

	"github.com/containerd/cgroups"
	"github.com/pkg/errors"
)

func getSysProcAttr() *syscall.SysProcAttr ***REMOVED***
	return &syscall.SysProcAttr***REMOVED***
		Setpgid: true,
	***REMOVED***
***REMOVED***

func setCgroup(cgroupPath string, cmd *exec.Cmd) error ***REMOVED***
	cg, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(cgroupPath))
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to load cgroup %s", cgroupPath)
	***REMOVED***
	if err := cg.Add(cgroups.Process***REMOVED***
		Pid: cmd.Process.Pid,
	***REMOVED***); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to join cgroup %s", cgroupPath)
	***REMOVED***
	return nil
***REMOVED***
