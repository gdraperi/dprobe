// +build !windows

package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/docker/docker/pkg/parsers/kernel"
	"github.com/docker/docker/pkg/sysinfo"
)

var (
	// SysInfo stores information about which features a kernel supports.
	SysInfo *sysinfo.SysInfo
)

func cpuCfsPeriod() bool ***REMOVED***
	return testEnv.DaemonInfo.CPUCfsPeriod
***REMOVED***

func cpuCfsQuota() bool ***REMOVED***
	return testEnv.DaemonInfo.CPUCfsQuota
***REMOVED***

func cpuShare() bool ***REMOVED***
	return testEnv.DaemonInfo.CPUShares
***REMOVED***

func oomControl() bool ***REMOVED***
	return testEnv.DaemonInfo.OomKillDisable
***REMOVED***

func pidsLimit() bool ***REMOVED***
	return SysInfo.PidsLimit
***REMOVED***

func kernelMemorySupport() bool ***REMOVED***
	return testEnv.DaemonInfo.KernelMemory
***REMOVED***

func memoryLimitSupport() bool ***REMOVED***
	return testEnv.DaemonInfo.MemoryLimit
***REMOVED***

func memoryReservationSupport() bool ***REMOVED***
	return SysInfo.MemoryReservation
***REMOVED***

func swapMemorySupport() bool ***REMOVED***
	return testEnv.DaemonInfo.SwapLimit
***REMOVED***

func memorySwappinessSupport() bool ***REMOVED***
	return SameHostDaemon() && SysInfo.MemorySwappiness
***REMOVED***

func blkioWeight() bool ***REMOVED***
	return SameHostDaemon() && SysInfo.BlkioWeight
***REMOVED***

func cgroupCpuset() bool ***REMOVED***
	return testEnv.DaemonInfo.CPUSet
***REMOVED***

func seccompEnabled() bool ***REMOVED***
	return supportsSeccomp && SysInfo.Seccomp
***REMOVED***

func bridgeNfIptables() bool ***REMOVED***
	return !SysInfo.BridgeNFCallIPTablesDisabled
***REMOVED***

func bridgeNfIP6tables() bool ***REMOVED***
	return !SysInfo.BridgeNFCallIP6TablesDisabled
***REMOVED***

func unprivilegedUsernsClone() bool ***REMOVED***
	content, err := ioutil.ReadFile("/proc/sys/kernel/unprivileged_userns_clone")
	return err != nil || !strings.Contains(string(content), "0")
***REMOVED***

func ambientCapabilities() bool ***REMOVED***
	content, err := ioutil.ReadFile("/proc/self/status")
	return err != nil || strings.Contains(string(content), "CapAmb:")
***REMOVED***

func overlayFSSupported() bool ***REMOVED***
	cmd := exec.Command(dockerBinary, "run", "--rm", "busybox", "/bin/sh", "-c", "cat /proc/filesystems")
	out, err := cmd.CombinedOutput()
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return bytes.Contains(out, []byte("overlay\n"))
***REMOVED***

func overlay2Supported() bool ***REMOVED***
	if !overlayFSSupported() ***REMOVED***
		return false
	***REMOVED***

	daemonV, err := kernel.ParseRelease(testEnv.DaemonInfo.KernelVersion)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	requiredV := kernel.VersionInfo***REMOVED***Kernel: 4***REMOVED***
	return kernel.CompareKernelVersion(*daemonV, requiredV) > -1

***REMOVED***

func init() ***REMOVED***
	if SameHostDaemon() ***REMOVED***
		SysInfo = sysinfo.New(true)
	***REMOVED***
***REMOVED***
