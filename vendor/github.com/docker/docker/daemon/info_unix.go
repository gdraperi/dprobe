// +build !windows

package daemon

import (
	"context"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// FillPlatformInfo fills the platform related info.
func (daemon *Daemon) FillPlatformInfo(v *types.Info, sysInfo *sysinfo.SysInfo) ***REMOVED***
	v.MemoryLimit = sysInfo.MemoryLimit
	v.SwapLimit = sysInfo.SwapLimit
	v.KernelMemory = sysInfo.KernelMemory
	v.OomKillDisable = sysInfo.OomKillDisable
	v.CPUCfsPeriod = sysInfo.CPUCfsPeriod
	v.CPUCfsQuota = sysInfo.CPUCfsQuota
	v.CPUShares = sysInfo.CPUShares
	v.CPUSet = sysInfo.Cpuset
	v.Runtimes = daemon.configStore.GetAllRuntimes()
	v.DefaultRuntime = daemon.configStore.GetDefaultRuntimeName()
	v.InitBinary = daemon.configStore.GetInitPath()

	v.RuncCommit.Expected = dockerversion.RuncCommitID
	defaultRuntimeBinary := daemon.configStore.GetRuntime(v.DefaultRuntime).Path
	if rv, err := exec.Command(defaultRuntimeBinary, "--version").Output(); err == nil ***REMOVED***
		parts := strings.Split(strings.TrimSpace(string(rv)), "\n")
		if len(parts) == 3 ***REMOVED***
			parts = strings.Split(parts[1], ": ")
			if len(parts) == 2 ***REMOVED***
				v.RuncCommit.ID = strings.TrimSpace(parts[1])
			***REMOVED***
		***REMOVED***

		if v.RuncCommit.ID == "" ***REMOVED***
			logrus.Warnf("failed to retrieve %s version: unknown output format: %s", defaultRuntimeBinary, string(rv))
			v.RuncCommit.ID = "N/A"
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		logrus.Warnf("failed to retrieve %s version: %v", defaultRuntimeBinary, err)
		v.RuncCommit.ID = "N/A"
	***REMOVED***

	v.ContainerdCommit.Expected = dockerversion.ContainerdCommitID
	if rv, err := daemon.containerd.Version(context.Background()); err == nil ***REMOVED***
		v.ContainerdCommit.ID = rv.Revision
	***REMOVED*** else ***REMOVED***
		logrus.Warnf("failed to retrieve containerd version: %v", err)
		v.ContainerdCommit.ID = "N/A"
	***REMOVED***

	defaultInitBinary := daemon.configStore.GetInitPath()
	if rv, err := exec.Command(defaultInitBinary, "--version").Output(); err == nil ***REMOVED***
		ver, err := parseInitVersion(string(rv))

		if err != nil ***REMOVED***
			logrus.Warnf("failed to retrieve %s version: %s", defaultInitBinary, err)
		***REMOVED***
		v.InitCommit = ver
	***REMOVED*** else ***REMOVED***
		logrus.Warnf("failed to retrieve %s version: %s", defaultInitBinary, err)
		v.InitCommit.ID = "N/A"
	***REMOVED***
***REMOVED***

// parseInitVersion parses a Tini version string, and extracts the version.
func parseInitVersion(v string) (types.Commit, error) ***REMOVED***
	version := types.Commit***REMOVED***ID: "", Expected: dockerversion.InitCommitID***REMOVED***
	parts := strings.Split(strings.TrimSpace(v), " - ")

	if len(parts) >= 2 ***REMOVED***
		gitParts := strings.Split(parts[1], ".")
		if len(gitParts) == 2 && gitParts[0] == "git" ***REMOVED***
			version.ID = gitParts[1]
			version.Expected = dockerversion.InitCommitID[0:len(version.ID)]
		***REMOVED***
	***REMOVED***
	if version.ID == "" && strings.HasPrefix(parts[0], "tini version ") ***REMOVED***
		version.ID = "v" + strings.TrimPrefix(parts[0], "tini version ")
	***REMOVED***
	if version.ID == "" ***REMOVED***
		version.ID = "N/A"
		return version, errors.Errorf("unknown output format: %s", v)
	***REMOVED***
	return version, nil
***REMOVED***
