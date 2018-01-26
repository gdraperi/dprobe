package sysinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func findCgroupMountpoints() (map[string]string, error) ***REMOVED***
	cgMounts, err := cgroups.GetCgroupMounts(false)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Failed to parse cgroup information: %v", err)
	***REMOVED***
	mps := make(map[string]string)
	for _, m := range cgMounts ***REMOVED***
		for _, ss := range m.Subsystems ***REMOVED***
			mps[ss] = m.Mountpoint
		***REMOVED***
	***REMOVED***
	return mps, nil
***REMOVED***

// New returns a new SysInfo, using the filesystem to detect which features
// the kernel supports. If `quiet` is `false` warnings are printed in logs
// whenever an error occurs or misconfigurations are present.
func New(quiet bool) *SysInfo ***REMOVED***
	sysInfo := &SysInfo***REMOVED******REMOVED***
	cgMounts, err := findCgroupMountpoints()
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to parse cgroup information: %v", err)
	***REMOVED*** else ***REMOVED***
		sysInfo.cgroupMemInfo = checkCgroupMem(cgMounts, quiet)
		sysInfo.cgroupCPUInfo = checkCgroupCPU(cgMounts, quiet)
		sysInfo.cgroupBlkioInfo = checkCgroupBlkioInfo(cgMounts, quiet)
		sysInfo.cgroupCpusetInfo = checkCgroupCpusetInfo(cgMounts, quiet)
		sysInfo.cgroupPids = checkCgroupPids(quiet)
	***REMOVED***

	_, ok := cgMounts["devices"]
	sysInfo.CgroupDevicesEnabled = ok

	sysInfo.IPv4ForwardingDisabled = !readProcBool("/proc/sys/net/ipv4/ip_forward")
	sysInfo.BridgeNFCallIPTablesDisabled = !readProcBool("/proc/sys/net/bridge/bridge-nf-call-iptables")
	sysInfo.BridgeNFCallIP6TablesDisabled = !readProcBool("/proc/sys/net/bridge/bridge-nf-call-ip6tables")

	// Check if AppArmor is supported.
	if _, err := os.Stat("/sys/kernel/security/apparmor"); !os.IsNotExist(err) ***REMOVED***
		sysInfo.AppArmor = true
	***REMOVED***

	// Check if Seccomp is supported, via CONFIG_SECCOMP.
	if err := unix.Prctl(unix.PR_GET_SECCOMP, 0, 0, 0, 0); err != unix.EINVAL ***REMOVED***
		// Make sure the kernel has CONFIG_SECCOMP_FILTER.
		if err := unix.Prctl(unix.PR_SET_SECCOMP, unix.SECCOMP_MODE_FILTER, 0, 0, 0); err != unix.EINVAL ***REMOVED***
			sysInfo.Seccomp = true
		***REMOVED***
	***REMOVED***

	return sysInfo
***REMOVED***

// checkCgroupMem reads the memory information from the memory cgroup mount point.
func checkCgroupMem(cgMounts map[string]string, quiet bool) cgroupMemInfo ***REMOVED***
	mountPoint, ok := cgMounts["memory"]
	if !ok ***REMOVED***
		if !quiet ***REMOVED***
			logrus.Warn("Your kernel does not support cgroup memory limit")
		***REMOVED***
		return cgroupMemInfo***REMOVED******REMOVED***
	***REMOVED***

	swapLimit := cgroupEnabled(mountPoint, "memory.memsw.limit_in_bytes")
	if !quiet && !swapLimit ***REMOVED***
		logrus.Warn("Your kernel does not support swap memory limit")
	***REMOVED***
	memoryReservation := cgroupEnabled(mountPoint, "memory.soft_limit_in_bytes")
	if !quiet && !memoryReservation ***REMOVED***
		logrus.Warn("Your kernel does not support memory reservation")
	***REMOVED***
	oomKillDisable := cgroupEnabled(mountPoint, "memory.oom_control")
	if !quiet && !oomKillDisable ***REMOVED***
		logrus.Warn("Your kernel does not support oom control")
	***REMOVED***
	memorySwappiness := cgroupEnabled(mountPoint, "memory.swappiness")
	if !quiet && !memorySwappiness ***REMOVED***
		logrus.Warn("Your kernel does not support memory swappiness")
	***REMOVED***
	kernelMemory := cgroupEnabled(mountPoint, "memory.kmem.limit_in_bytes")
	if !quiet && !kernelMemory ***REMOVED***
		logrus.Warn("Your kernel does not support kernel memory limit")
	***REMOVED***

	return cgroupMemInfo***REMOVED***
		MemoryLimit:       true,
		SwapLimit:         swapLimit,
		MemoryReservation: memoryReservation,
		OomKillDisable:    oomKillDisable,
		MemorySwappiness:  memorySwappiness,
		KernelMemory:      kernelMemory,
	***REMOVED***
***REMOVED***

// checkCgroupCPU reads the cpu information from the cpu cgroup mount point.
func checkCgroupCPU(cgMounts map[string]string, quiet bool) cgroupCPUInfo ***REMOVED***
	mountPoint, ok := cgMounts["cpu"]
	if !ok ***REMOVED***
		if !quiet ***REMOVED***
			logrus.Warn("Unable to find cpu cgroup in mounts")
		***REMOVED***
		return cgroupCPUInfo***REMOVED******REMOVED***
	***REMOVED***

	cpuShares := cgroupEnabled(mountPoint, "cpu.shares")
	if !quiet && !cpuShares ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup cpu shares")
	***REMOVED***

	cpuCfsPeriod := cgroupEnabled(mountPoint, "cpu.cfs_period_us")
	if !quiet && !cpuCfsPeriod ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup cfs period")
	***REMOVED***

	cpuCfsQuota := cgroupEnabled(mountPoint, "cpu.cfs_quota_us")
	if !quiet && !cpuCfsQuota ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup cfs quotas")
	***REMOVED***

	cpuRealtimePeriod := cgroupEnabled(mountPoint, "cpu.rt_period_us")
	if !quiet && !cpuRealtimePeriod ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup rt period")
	***REMOVED***

	cpuRealtimeRuntime := cgroupEnabled(mountPoint, "cpu.rt_runtime_us")
	if !quiet && !cpuRealtimeRuntime ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup rt runtime")
	***REMOVED***

	return cgroupCPUInfo***REMOVED***
		CPUShares:          cpuShares,
		CPUCfsPeriod:       cpuCfsPeriod,
		CPUCfsQuota:        cpuCfsQuota,
		CPURealtimePeriod:  cpuRealtimePeriod,
		CPURealtimeRuntime: cpuRealtimeRuntime,
	***REMOVED***
***REMOVED***

// checkCgroupBlkioInfo reads the blkio information from the blkio cgroup mount point.
func checkCgroupBlkioInfo(cgMounts map[string]string, quiet bool) cgroupBlkioInfo ***REMOVED***
	mountPoint, ok := cgMounts["blkio"]
	if !ok ***REMOVED***
		if !quiet ***REMOVED***
			logrus.Warn("Unable to find blkio cgroup in mounts")
		***REMOVED***
		return cgroupBlkioInfo***REMOVED******REMOVED***
	***REMOVED***

	weight := cgroupEnabled(mountPoint, "blkio.weight")
	if !quiet && !weight ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup blkio weight")
	***REMOVED***

	weightDevice := cgroupEnabled(mountPoint, "blkio.weight_device")
	if !quiet && !weightDevice ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup blkio weight_device")
	***REMOVED***

	readBpsDevice := cgroupEnabled(mountPoint, "blkio.throttle.read_bps_device")
	if !quiet && !readBpsDevice ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup blkio throttle.read_bps_device")
	***REMOVED***

	writeBpsDevice := cgroupEnabled(mountPoint, "blkio.throttle.write_bps_device")
	if !quiet && !writeBpsDevice ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup blkio throttle.write_bps_device")
	***REMOVED***
	readIOpsDevice := cgroupEnabled(mountPoint, "blkio.throttle.read_iops_device")
	if !quiet && !readIOpsDevice ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup blkio throttle.read_iops_device")
	***REMOVED***

	writeIOpsDevice := cgroupEnabled(mountPoint, "blkio.throttle.write_iops_device")
	if !quiet && !writeIOpsDevice ***REMOVED***
		logrus.Warn("Your kernel does not support cgroup blkio throttle.write_iops_device")
	***REMOVED***
	return cgroupBlkioInfo***REMOVED***
		BlkioWeight:          weight,
		BlkioWeightDevice:    weightDevice,
		BlkioReadBpsDevice:   readBpsDevice,
		BlkioWriteBpsDevice:  writeBpsDevice,
		BlkioReadIOpsDevice:  readIOpsDevice,
		BlkioWriteIOpsDevice: writeIOpsDevice,
	***REMOVED***
***REMOVED***

// checkCgroupCpusetInfo reads the cpuset information from the cpuset cgroup mount point.
func checkCgroupCpusetInfo(cgMounts map[string]string, quiet bool) cgroupCpusetInfo ***REMOVED***
	mountPoint, ok := cgMounts["cpuset"]
	if !ok ***REMOVED***
		if !quiet ***REMOVED***
			logrus.Warn("Unable to find cpuset cgroup in mounts")
		***REMOVED***
		return cgroupCpusetInfo***REMOVED******REMOVED***
	***REMOVED***

	cpus, err := ioutil.ReadFile(path.Join(mountPoint, "cpuset.cpus"))
	if err != nil ***REMOVED***
		return cgroupCpusetInfo***REMOVED******REMOVED***
	***REMOVED***

	mems, err := ioutil.ReadFile(path.Join(mountPoint, "cpuset.mems"))
	if err != nil ***REMOVED***
		return cgroupCpusetInfo***REMOVED******REMOVED***
	***REMOVED***

	return cgroupCpusetInfo***REMOVED***
		Cpuset: true,
		Cpus:   strings.TrimSpace(string(cpus)),
		Mems:   strings.TrimSpace(string(mems)),
	***REMOVED***
***REMOVED***

// checkCgroupPids reads the pids information from the pids cgroup mount point.
func checkCgroupPids(quiet bool) cgroupPids ***REMOVED***
	_, err := cgroups.FindCgroupMountpoint("pids")
	if err != nil ***REMOVED***
		if !quiet ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
		return cgroupPids***REMOVED******REMOVED***
	***REMOVED***

	return cgroupPids***REMOVED***
		PidsLimit: true,
	***REMOVED***
***REMOVED***

func cgroupEnabled(mountPoint, name string) bool ***REMOVED***
	_, err := os.Stat(path.Join(mountPoint, name))
	return err == nil
***REMOVED***

func readProcBool(path string) bool ***REMOVED***
	val, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return strings.TrimSpace(string(val)) == "1"
***REMOVED***
