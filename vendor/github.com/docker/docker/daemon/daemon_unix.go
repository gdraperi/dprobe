// +build linux freebsd

package daemon

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	containerd_cgroups "github.com/containerd/cgroups"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/blkiodev"
	pblkiodev "github.com/docker/docker/api/types/blkiodev"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/image"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/parsers"
	"github.com/docker/docker/pkg/parsers/kernel"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/docker/docker/runconfig"
	"github.com/docker/docker/volume"
	"github.com/docker/libnetwork"
	nwconfig "github.com/docker/libnetwork/config"
	"github.com/docker/libnetwork/drivers/bridge"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/netutils"
	"github.com/docker/libnetwork/options"
	lntypes "github.com/docker/libnetwork/types"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	rsystem "github.com/opencontainers/runc/libcontainer/system"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const (
	// DefaultShimBinary is the default shim to be used by containerd if none
	// is specified
	DefaultShimBinary = "docker-containerd-shim"

	// DefaultRuntimeBinary is the default runtime to be used by
	// containerd if none is specified
	DefaultRuntimeBinary = "docker-runc"

	// See https://git.kernel.org/cgit/linux/kernel/git/tip/tip.git/tree/kernel/sched/sched.h?id=8cd9234c64c584432f6992fe944ca9e46ca8ea76#n269
	linuxMinCPUShares = 2
	linuxMaxCPUShares = 262144
	platformSupported = true
	// It's not kernel limit, we want this 4M limit to supply a reasonable functional container
	linuxMinMemory = 4194304
	// constants for remapped root settings
	defaultIDSpecifier string = "default"
	defaultRemappedID  string = "dockremap"

	// constant for cgroup drivers
	cgroupFsDriver      = "cgroupfs"
	cgroupSystemdDriver = "systemd"

	// DefaultRuntimeName is the default runtime to be used by
	// containerd if none is specified
	DefaultRuntimeName = "docker-runc"
)

type containerGetter interface ***REMOVED***
	GetContainer(string) (*container.Container, error)
***REMOVED***

func getMemoryResources(config containertypes.Resources) *specs.LinuxMemory ***REMOVED***
	memory := specs.LinuxMemory***REMOVED******REMOVED***

	if config.Memory > 0 ***REMOVED***
		memory.Limit = &config.Memory
	***REMOVED***

	if config.MemoryReservation > 0 ***REMOVED***
		memory.Reservation = &config.MemoryReservation
	***REMOVED***

	if config.MemorySwap > 0 ***REMOVED***
		memory.Swap = &config.MemorySwap
	***REMOVED***

	if config.MemorySwappiness != nil ***REMOVED***
		swappiness := uint64(*config.MemorySwappiness)
		memory.Swappiness = &swappiness
	***REMOVED***

	if config.KernelMemory != 0 ***REMOVED***
		memory.Kernel = &config.KernelMemory
	***REMOVED***

	return &memory
***REMOVED***

func getCPUResources(config containertypes.Resources) (*specs.LinuxCPU, error) ***REMOVED***
	cpu := specs.LinuxCPU***REMOVED******REMOVED***

	if config.CPUShares < 0 ***REMOVED***
		return nil, fmt.Errorf("shares: invalid argument")
	***REMOVED***
	if config.CPUShares >= 0 ***REMOVED***
		shares := uint64(config.CPUShares)
		cpu.Shares = &shares
	***REMOVED***

	if config.CpusetCpus != "" ***REMOVED***
		cpu.Cpus = config.CpusetCpus
	***REMOVED***

	if config.CpusetMems != "" ***REMOVED***
		cpu.Mems = config.CpusetMems
	***REMOVED***

	if config.NanoCPUs > 0 ***REMOVED***
		// https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
		period := uint64(100 * time.Millisecond / time.Microsecond)
		quota := config.NanoCPUs * int64(period) / 1e9
		cpu.Period = &period
		cpu.Quota = &quota
	***REMOVED***

	if config.CPUPeriod != 0 ***REMOVED***
		period := uint64(config.CPUPeriod)
		cpu.Period = &period
	***REMOVED***

	if config.CPUQuota != 0 ***REMOVED***
		q := config.CPUQuota
		cpu.Quota = &q
	***REMOVED***

	if config.CPURealtimePeriod != 0 ***REMOVED***
		period := uint64(config.CPURealtimePeriod)
		cpu.RealtimePeriod = &period
	***REMOVED***

	if config.CPURealtimeRuntime != 0 ***REMOVED***
		c := config.CPURealtimeRuntime
		cpu.RealtimeRuntime = &c
	***REMOVED***

	return &cpu, nil
***REMOVED***

func getBlkioWeightDevices(config containertypes.Resources) ([]specs.LinuxWeightDevice, error) ***REMOVED***
	var stat unix.Stat_t
	var blkioWeightDevices []specs.LinuxWeightDevice

	for _, weightDevice := range config.BlkioWeightDevice ***REMOVED***
		if err := unix.Stat(weightDevice.Path, &stat); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		weight := weightDevice.Weight
		d := specs.LinuxWeightDevice***REMOVED***Weight: &weight***REMOVED***
		d.Major = int64(stat.Rdev / 256)
		d.Minor = int64(stat.Rdev % 256)
		blkioWeightDevices = append(blkioWeightDevices, d)
	***REMOVED***

	return blkioWeightDevices, nil
***REMOVED***

func (daemon *Daemon) parseSecurityOpt(container *container.Container, hostConfig *containertypes.HostConfig) error ***REMOVED***
	container.NoNewPrivileges = daemon.configStore.NoNewPrivileges
	return parseSecurityOpt(container, hostConfig)
***REMOVED***

func parseSecurityOpt(container *container.Container, config *containertypes.HostConfig) error ***REMOVED***
	var (
		labelOpts []string
		err       error
	)

	for _, opt := range config.SecurityOpt ***REMOVED***
		if opt == "no-new-privileges" ***REMOVED***
			container.NoNewPrivileges = true
			continue
		***REMOVED***
		if opt == "disable" ***REMOVED***
			labelOpts = append(labelOpts, "disable")
			continue
		***REMOVED***

		var con []string
		if strings.Contains(opt, "=") ***REMOVED***
			con = strings.SplitN(opt, "=", 2)
		***REMOVED*** else if strings.Contains(opt, ":") ***REMOVED***
			con = strings.SplitN(opt, ":", 2)
			logrus.Warn("Security options with `:` as a separator are deprecated and will be completely unsupported in 17.04, use `=` instead.")
		***REMOVED***
		if len(con) != 2 ***REMOVED***
			return fmt.Errorf("invalid --security-opt 1: %q", opt)
		***REMOVED***

		switch con[0] ***REMOVED***
		case "label":
			labelOpts = append(labelOpts, con[1])
		case "apparmor":
			container.AppArmorProfile = con[1]
		case "seccomp":
			container.SeccompProfile = con[1]
		case "no-new-privileges":
			noNewPrivileges, err := strconv.ParseBool(con[1])
			if err != nil ***REMOVED***
				return fmt.Errorf("invalid --security-opt 2: %q", opt)
			***REMOVED***
			container.NoNewPrivileges = noNewPrivileges
		default:
			return fmt.Errorf("invalid --security-opt 2: %q", opt)
		***REMOVED***
	***REMOVED***

	container.ProcessLabel, container.MountLabel, err = label.InitLabels(labelOpts)
	return err
***REMOVED***

func getBlkioThrottleDevices(devs []*blkiodev.ThrottleDevice) ([]specs.LinuxThrottleDevice, error) ***REMOVED***
	var throttleDevices []specs.LinuxThrottleDevice
	var stat unix.Stat_t

	for _, d := range devs ***REMOVED***
		if err := unix.Stat(d.Path, &stat); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		d := specs.LinuxThrottleDevice***REMOVED***Rate: d.Rate***REMOVED***
		d.Major = int64(stat.Rdev / 256)
		d.Minor = int64(stat.Rdev % 256)
		throttleDevices = append(throttleDevices, d)
	***REMOVED***

	return throttleDevices, nil
***REMOVED***

func checkKernel() error ***REMOVED***
	// Check for unsupported kernel versions
	// FIXME: it would be cleaner to not test for specific versions, but rather
	// test for specific functionalities.
	// Unfortunately we can't test for the feature "does not cause a kernel panic"
	// without actually causing a kernel panic, so we need this workaround until
	// the circumstances of pre-3.10 crashes are clearer.
	// For details see https://github.com/docker/docker/issues/407
	// Docker 1.11 and above doesn't actually run on kernels older than 3.4,
	// due to containerd-shim usage of PR_SET_CHILD_SUBREAPER (introduced in 3.4).
	if !kernel.CheckKernelVersion(3, 10, 0) ***REMOVED***
		v, _ := kernel.GetKernelVersion()
		if os.Getenv("DOCKER_NOWARN_KERNEL_VERSION") == "" ***REMOVED***
			logrus.Fatalf("Your Linux kernel version %s is not supported for running docker. Please upgrade your kernel to 3.10.0 or newer.", v.String())
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// adaptContainerSettings is called during container creation to modify any
// settings necessary in the HostConfig structure.
func (daemon *Daemon) adaptContainerSettings(hostConfig *containertypes.HostConfig, adjustCPUShares bool) error ***REMOVED***
	if adjustCPUShares && hostConfig.CPUShares > 0 ***REMOVED***
		// Handle unsupported CPUShares
		if hostConfig.CPUShares < linuxMinCPUShares ***REMOVED***
			logrus.Warnf("Changing requested CPUShares of %d to minimum allowed of %d", hostConfig.CPUShares, linuxMinCPUShares)
			hostConfig.CPUShares = linuxMinCPUShares
		***REMOVED*** else if hostConfig.CPUShares > linuxMaxCPUShares ***REMOVED***
			logrus.Warnf("Changing requested CPUShares of %d to maximum allowed of %d", hostConfig.CPUShares, linuxMaxCPUShares)
			hostConfig.CPUShares = linuxMaxCPUShares
		***REMOVED***
	***REMOVED***
	if hostConfig.Memory > 0 && hostConfig.MemorySwap == 0 ***REMOVED***
		// By default, MemorySwap is set to twice the size of Memory.
		hostConfig.MemorySwap = hostConfig.Memory * 2
	***REMOVED***
	if hostConfig.ShmSize == 0 ***REMOVED***
		hostConfig.ShmSize = config.DefaultShmSize
		if daemon.configStore != nil ***REMOVED***
			hostConfig.ShmSize = int64(daemon.configStore.ShmSize)
		***REMOVED***
	***REMOVED***
	// Set default IPC mode, if unset for container
	if hostConfig.IpcMode.IsEmpty() ***REMOVED***
		m := config.DefaultIpcMode
		if daemon.configStore != nil ***REMOVED***
			m = daemon.configStore.IpcMode
		***REMOVED***
		hostConfig.IpcMode = containertypes.IpcMode(m)
	***REMOVED***

	adaptSharedNamespaceContainer(daemon, hostConfig)

	var err error
	opts, err := daemon.generateSecurityOpt(hostConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	hostConfig.SecurityOpt = append(hostConfig.SecurityOpt, opts...)
	if hostConfig.OomKillDisable == nil ***REMOVED***
		defaultOomKillDisable := false
		hostConfig.OomKillDisable = &defaultOomKillDisable
	***REMOVED***

	return nil
***REMOVED***

// adaptSharedNamespaceContainer replaces container name with its ID in hostConfig.
// To be more precisely, it modifies `container:name` to `container:ID` of PidMode, IpcMode
// and NetworkMode.
//
// When a container shares its namespace with another container, use ID can keep the namespace
// sharing connection between the two containers even the another container is renamed.
func adaptSharedNamespaceContainer(daemon containerGetter, hostConfig *containertypes.HostConfig) ***REMOVED***
	containerPrefix := "container:"
	if hostConfig.PidMode.IsContainer() ***REMOVED***
		pidContainer := hostConfig.PidMode.Container()
		// if there is any error returned here, we just ignore it and leave it to be
		// handled in the following logic
		if c, err := daemon.GetContainer(pidContainer); err == nil ***REMOVED***
			hostConfig.PidMode = containertypes.PidMode(containerPrefix + c.ID)
		***REMOVED***
	***REMOVED***
	if hostConfig.IpcMode.IsContainer() ***REMOVED***
		ipcContainer := hostConfig.IpcMode.Container()
		if c, err := daemon.GetContainer(ipcContainer); err == nil ***REMOVED***
			hostConfig.IpcMode = containertypes.IpcMode(containerPrefix + c.ID)
		***REMOVED***
	***REMOVED***
	if hostConfig.NetworkMode.IsContainer() ***REMOVED***
		netContainer := hostConfig.NetworkMode.ConnectedContainer()
		if c, err := daemon.GetContainer(netContainer); err == nil ***REMOVED***
			hostConfig.NetworkMode = containertypes.NetworkMode(containerPrefix + c.ID)
		***REMOVED***
	***REMOVED***
***REMOVED***

func verifyContainerResources(resources *containertypes.Resources, sysInfo *sysinfo.SysInfo, update bool) ([]string, error) ***REMOVED***
	warnings := []string***REMOVED******REMOVED***
	fixMemorySwappiness(resources)

	// memory subsystem checks and adjustments
	if resources.Memory != 0 && resources.Memory < linuxMinMemory ***REMOVED***
		return warnings, fmt.Errorf("Minimum memory limit allowed is 4MB")
	***REMOVED***
	if resources.Memory > 0 && !sysInfo.MemoryLimit ***REMOVED***
		warnings = append(warnings, "Your kernel does not support memory limit capabilities or the cgroup is not mounted. Limitation discarded.")
		logrus.Warn("Your kernel does not support memory limit capabilities or the cgroup is not mounted. Limitation discarded.")
		resources.Memory = 0
		resources.MemorySwap = -1
	***REMOVED***
	if resources.Memory > 0 && resources.MemorySwap != -1 && !sysInfo.SwapLimit ***REMOVED***
		warnings = append(warnings, "Your kernel does not support swap limit capabilities or the cgroup is not mounted. Memory limited without swap.")
		logrus.Warn("Your kernel does not support swap limit capabilities,or the cgroup is not mounted. Memory limited without swap.")
		resources.MemorySwap = -1
	***REMOVED***
	if resources.Memory > 0 && resources.MemorySwap > 0 && resources.MemorySwap < resources.Memory ***REMOVED***
		return warnings, fmt.Errorf("Minimum memoryswap limit should be larger than memory limit, see usage")
	***REMOVED***
	if resources.Memory == 0 && resources.MemorySwap > 0 && !update ***REMOVED***
		return warnings, fmt.Errorf("You should always set the Memory limit when using Memoryswap limit, see usage")
	***REMOVED***
	if resources.MemorySwappiness != nil && !sysInfo.MemorySwappiness ***REMOVED***
		warnings = append(warnings, "Your kernel does not support memory swappiness capabilities or the cgroup is not mounted. Memory swappiness discarded.")
		logrus.Warn("Your kernel does not support memory swappiness capabilities, or the cgroup is not mounted. Memory swappiness discarded.")
		resources.MemorySwappiness = nil
	***REMOVED***
	if resources.MemorySwappiness != nil ***REMOVED***
		swappiness := *resources.MemorySwappiness
		if swappiness < 0 || swappiness > 100 ***REMOVED***
			return warnings, fmt.Errorf("Invalid value: %v, valid memory swappiness range is 0-100", swappiness)
		***REMOVED***
	***REMOVED***
	if resources.MemoryReservation > 0 && !sysInfo.MemoryReservation ***REMOVED***
		warnings = append(warnings, "Your kernel does not support memory soft limit capabilities or the cgroup is not mounted. Limitation discarded.")
		logrus.Warn("Your kernel does not support memory soft limit capabilities or the cgroup is not mounted. Limitation discarded.")
		resources.MemoryReservation = 0
	***REMOVED***
	if resources.MemoryReservation > 0 && resources.MemoryReservation < linuxMinMemory ***REMOVED***
		return warnings, fmt.Errorf("Minimum memory reservation allowed is 4MB")
	***REMOVED***
	if resources.Memory > 0 && resources.MemoryReservation > 0 && resources.Memory < resources.MemoryReservation ***REMOVED***
		return warnings, fmt.Errorf("Minimum memory limit can not be less than memory reservation limit, see usage")
	***REMOVED***
	if resources.KernelMemory > 0 && !sysInfo.KernelMemory ***REMOVED***
		warnings = append(warnings, "Your kernel does not support kernel memory limit capabilities or the cgroup is not mounted. Limitation discarded.")
		logrus.Warn("Your kernel does not support kernel memory limit capabilities or the cgroup is not mounted. Limitation discarded.")
		resources.KernelMemory = 0
	***REMOVED***
	if resources.KernelMemory > 0 && resources.KernelMemory < linuxMinMemory ***REMOVED***
		return warnings, fmt.Errorf("Minimum kernel memory limit allowed is 4MB")
	***REMOVED***
	if resources.KernelMemory > 0 && !kernel.CheckKernelVersion(4, 0, 0) ***REMOVED***
		warnings = append(warnings, "You specified a kernel memory limit on a kernel older than 4.0. Kernel memory limits are experimental on older kernels, it won't work as expected and can cause your system to be unstable.")
		logrus.Warn("You specified a kernel memory limit on a kernel older than 4.0. Kernel memory limits are experimental on older kernels, it won't work as expected and can cause your system to be unstable.")
	***REMOVED***
	if resources.OomKillDisable != nil && !sysInfo.OomKillDisable ***REMOVED***
		// only produce warnings if the setting wasn't to *disable* the OOM Kill; no point
		// warning the caller if they already wanted the feature to be off
		if *resources.OomKillDisable ***REMOVED***
			warnings = append(warnings, "Your kernel does not support OomKillDisable. OomKillDisable discarded.")
			logrus.Warn("Your kernel does not support OomKillDisable. OomKillDisable discarded.")
		***REMOVED***
		resources.OomKillDisable = nil
	***REMOVED***

	if resources.PidsLimit != 0 && !sysInfo.PidsLimit ***REMOVED***
		warnings = append(warnings, "Your kernel does not support pids limit capabilities or the cgroup is not mounted. PIDs limit discarded.")
		logrus.Warn("Your kernel does not support pids limit capabilities or the cgroup is not mounted. PIDs limit discarded.")
		resources.PidsLimit = 0
	***REMOVED***

	// cpu subsystem checks and adjustments
	if resources.NanoCPUs > 0 && resources.CPUPeriod > 0 ***REMOVED***
		return warnings, fmt.Errorf("Conflicting options: Nano CPUs and CPU Period cannot both be set")
	***REMOVED***
	if resources.NanoCPUs > 0 && resources.CPUQuota > 0 ***REMOVED***
		return warnings, fmt.Errorf("Conflicting options: Nano CPUs and CPU Quota cannot both be set")
	***REMOVED***
	if resources.NanoCPUs > 0 && (!sysInfo.CPUCfsPeriod || !sysInfo.CPUCfsQuota) ***REMOVED***
		return warnings, fmt.Errorf("NanoCPUs can not be set, as your kernel does not support CPU cfs period/quota or the cgroup is not mounted")
	***REMOVED***
	// The highest precision we could get on Linux is 0.001, by setting
	//   cpu.cfs_period_us=1000ms
	//   cpu.cfs_quota=1ms
	// See the following link for details:
	// https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
	// Here we don't set the lower limit and it is up to the underlying platform (e.g., Linux) to return an error.
	// The error message is 0.01 so that this is consistent with Windows
	if resources.NanoCPUs < 0 || resources.NanoCPUs > int64(sysinfo.NumCPU())*1e9 ***REMOVED***
		return warnings, fmt.Errorf("Range of CPUs is from 0.01 to %d.00, as there are only %d CPUs available", sysinfo.NumCPU(), sysinfo.NumCPU())
	***REMOVED***

	if resources.CPUShares > 0 && !sysInfo.CPUShares ***REMOVED***
		warnings = append(warnings, "Your kernel does not support CPU shares or the cgroup is not mounted. Shares discarded.")
		logrus.Warn("Your kernel does not support CPU shares or the cgroup is not mounted. Shares discarded.")
		resources.CPUShares = 0
	***REMOVED***
	if resources.CPUPeriod > 0 && !sysInfo.CPUCfsPeriod ***REMOVED***
		warnings = append(warnings, "Your kernel does not support CPU cfs period or the cgroup is not mounted. Period discarded.")
		logrus.Warn("Your kernel does not support CPU cfs period or the cgroup is not mounted. Period discarded.")
		resources.CPUPeriod = 0
	***REMOVED***
	if resources.CPUPeriod != 0 && (resources.CPUPeriod < 1000 || resources.CPUPeriod > 1000000) ***REMOVED***
		return warnings, fmt.Errorf("CPU cfs period can not be less than 1ms (i.e. 1000) or larger than 1s (i.e. 1000000)")
	***REMOVED***
	if resources.CPUQuota > 0 && !sysInfo.CPUCfsQuota ***REMOVED***
		warnings = append(warnings, "Your kernel does not support CPU cfs quota or the cgroup is not mounted. Quota discarded.")
		logrus.Warn("Your kernel does not support CPU cfs quota or the cgroup is not mounted. Quota discarded.")
		resources.CPUQuota = 0
	***REMOVED***
	if resources.CPUQuota > 0 && resources.CPUQuota < 1000 ***REMOVED***
		return warnings, fmt.Errorf("CPU cfs quota can not be less than 1ms (i.e. 1000)")
	***REMOVED***
	if resources.CPUPercent > 0 ***REMOVED***
		warnings = append(warnings, fmt.Sprintf("%s does not support CPU percent. Percent discarded.", runtime.GOOS))
		logrus.Warnf("%s does not support CPU percent. Percent discarded.", runtime.GOOS)
		resources.CPUPercent = 0
	***REMOVED***

	// cpuset subsystem checks and adjustments
	if (resources.CpusetCpus != "" || resources.CpusetMems != "") && !sysInfo.Cpuset ***REMOVED***
		warnings = append(warnings, "Your kernel does not support cpuset or the cgroup is not mounted. Cpuset discarded.")
		logrus.Warn("Your kernel does not support cpuset or the cgroup is not mounted. Cpuset discarded.")
		resources.CpusetCpus = ""
		resources.CpusetMems = ""
	***REMOVED***
	cpusAvailable, err := sysInfo.IsCpusetCpusAvailable(resources.CpusetCpus)
	if err != nil ***REMOVED***
		return warnings, fmt.Errorf("Invalid value %s for cpuset cpus", resources.CpusetCpus)
	***REMOVED***
	if !cpusAvailable ***REMOVED***
		return warnings, fmt.Errorf("Requested CPUs are not available - requested %s, available: %s", resources.CpusetCpus, sysInfo.Cpus)
	***REMOVED***
	memsAvailable, err := sysInfo.IsCpusetMemsAvailable(resources.CpusetMems)
	if err != nil ***REMOVED***
		return warnings, fmt.Errorf("Invalid value %s for cpuset mems", resources.CpusetMems)
	***REMOVED***
	if !memsAvailable ***REMOVED***
		return warnings, fmt.Errorf("Requested memory nodes are not available - requested %s, available: %s", resources.CpusetMems, sysInfo.Mems)
	***REMOVED***

	// blkio subsystem checks and adjustments
	if resources.BlkioWeight > 0 && !sysInfo.BlkioWeight ***REMOVED***
		warnings = append(warnings, "Your kernel does not support Block I/O weight or the cgroup is not mounted. Weight discarded.")
		logrus.Warn("Your kernel does not support Block I/O weight or the cgroup is not mounted. Weight discarded.")
		resources.BlkioWeight = 0
	***REMOVED***
	if resources.BlkioWeight > 0 && (resources.BlkioWeight < 10 || resources.BlkioWeight > 1000) ***REMOVED***
		return warnings, fmt.Errorf("Range of blkio weight is from 10 to 1000")
	***REMOVED***
	if resources.IOMaximumBandwidth != 0 || resources.IOMaximumIOps != 0 ***REMOVED***
		return warnings, fmt.Errorf("Invalid QoS settings: %s does not support Maximum IO Bandwidth or Maximum IO IOps", runtime.GOOS)
	***REMOVED***
	if len(resources.BlkioWeightDevice) > 0 && !sysInfo.BlkioWeightDevice ***REMOVED***
		warnings = append(warnings, "Your kernel does not support Block I/O weight_device or the cgroup is not mounted. Weight-device discarded.")
		logrus.Warn("Your kernel does not support Block I/O weight_device or the cgroup is not mounted. Weight-device discarded.")
		resources.BlkioWeightDevice = []*pblkiodev.WeightDevice***REMOVED******REMOVED***
	***REMOVED***
	if len(resources.BlkioDeviceReadBps) > 0 && !sysInfo.BlkioReadBpsDevice ***REMOVED***
		warnings = append(warnings, "Your kernel does not support BPS Block I/O read limit or the cgroup is not mounted. Block I/O BPS read limit discarded.")
		logrus.Warn("Your kernel does not support BPS Block I/O read limit or the cgroup is not mounted. Block I/O BPS read limit discarded")
		resources.BlkioDeviceReadBps = []*pblkiodev.ThrottleDevice***REMOVED******REMOVED***
	***REMOVED***
	if len(resources.BlkioDeviceWriteBps) > 0 && !sysInfo.BlkioWriteBpsDevice ***REMOVED***
		warnings = append(warnings, "Your kernel does not support BPS Block I/O write limit or the cgroup is not mounted. Block I/O BPS write limit discarded.")
		logrus.Warn("Your kernel does not support BPS Block I/O write limit or the cgroup is not mounted. Block I/O BPS write limit discarded.")
		resources.BlkioDeviceWriteBps = []*pblkiodev.ThrottleDevice***REMOVED******REMOVED***

	***REMOVED***
	if len(resources.BlkioDeviceReadIOps) > 0 && !sysInfo.BlkioReadIOpsDevice ***REMOVED***
		warnings = append(warnings, "Your kernel does not support IOPS Block read limit or the cgroup is not mounted. Block I/O IOPS read limit discarded.")
		logrus.Warn("Your kernel does not support IOPS Block I/O read limit in IO or the cgroup is not mounted. Block I/O IOPS read limit discarded.")
		resources.BlkioDeviceReadIOps = []*pblkiodev.ThrottleDevice***REMOVED******REMOVED***
	***REMOVED***
	if len(resources.BlkioDeviceWriteIOps) > 0 && !sysInfo.BlkioWriteIOpsDevice ***REMOVED***
		warnings = append(warnings, "Your kernel does not support IOPS Block write limit or the cgroup is not mounted. Block I/O IOPS write limit discarded.")
		logrus.Warn("Your kernel does not support IOPS Block I/O write limit or the cgroup is not mounted. Block I/O IOPS write limit discarded.")
		resources.BlkioDeviceWriteIOps = []*pblkiodev.ThrottleDevice***REMOVED******REMOVED***
	***REMOVED***

	return warnings, nil
***REMOVED***

func (daemon *Daemon) getCgroupDriver() string ***REMOVED***
	cgroupDriver := cgroupFsDriver

	if UsingSystemd(daemon.configStore) ***REMOVED***
		cgroupDriver = cgroupSystemdDriver
	***REMOVED***
	return cgroupDriver
***REMOVED***

// getCD gets the raw value of the native.cgroupdriver option, if set.
func getCD(config *config.Config) string ***REMOVED***
	for _, option := range config.ExecOptions ***REMOVED***
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil || !strings.EqualFold(key, "native.cgroupdriver") ***REMOVED***
			continue
		***REMOVED***
		return val
	***REMOVED***
	return ""
***REMOVED***

// VerifyCgroupDriver validates native.cgroupdriver
func VerifyCgroupDriver(config *config.Config) error ***REMOVED***
	cd := getCD(config)
	if cd == "" || cd == cgroupFsDriver || cd == cgroupSystemdDriver ***REMOVED***
		return nil
	***REMOVED***
	return fmt.Errorf("native.cgroupdriver option %s not supported", cd)
***REMOVED***

// UsingSystemd returns true if cli option includes native.cgroupdriver=systemd
func UsingSystemd(config *config.Config) bool ***REMOVED***
	return getCD(config) == cgroupSystemdDriver
***REMOVED***

// verifyPlatformContainerSettings performs platform-specific validation of the
// hostconfig and config structures.
func verifyPlatformContainerSettings(daemon *Daemon, hostConfig *containertypes.HostConfig, config *containertypes.Config, update bool) ([]string, error) ***REMOVED***
	var warnings []string
	sysInfo := sysinfo.New(true)

	w, err := verifyContainerResources(&hostConfig.Resources, sysInfo, update)

	// no matter err is nil or not, w could have data in itself.
	warnings = append(warnings, w...)

	if err != nil ***REMOVED***
		return warnings, err
	***REMOVED***

	if hostConfig.ShmSize < 0 ***REMOVED***
		return warnings, fmt.Errorf("SHM size can not be less than 0")
	***REMOVED***

	if hostConfig.OomScoreAdj < -1000 || hostConfig.OomScoreAdj > 1000 ***REMOVED***
		return warnings, fmt.Errorf("Invalid value %d, range for oom score adj is [-1000, 1000]", hostConfig.OomScoreAdj)
	***REMOVED***

	// ip-forwarding does not affect container with '--net=host' (or '--net=none')
	if sysInfo.IPv4ForwardingDisabled && !(hostConfig.NetworkMode.IsHost() || hostConfig.NetworkMode.IsNone()) ***REMOVED***
		warnings = append(warnings, "IPv4 forwarding is disabled. Networking will not work.")
		logrus.Warn("IPv4 forwarding is disabled. Networking will not work")
	***REMOVED***
	// check for various conflicting options with user namespaces
	if daemon.configStore.RemappedRoot != "" && hostConfig.UsernsMode.IsPrivate() ***REMOVED***
		if hostConfig.Privileged ***REMOVED***
			return warnings, fmt.Errorf("privileged mode is incompatible with user namespaces.  You must run the container in the host namespace when running privileged mode")
		***REMOVED***
		if hostConfig.NetworkMode.IsHost() && !hostConfig.UsernsMode.IsHost() ***REMOVED***
			return warnings, fmt.Errorf("cannot share the host's network namespace when user namespaces are enabled")
		***REMOVED***
		if hostConfig.PidMode.IsHost() && !hostConfig.UsernsMode.IsHost() ***REMOVED***
			return warnings, fmt.Errorf("cannot share the host PID namespace when user namespaces are enabled")
		***REMOVED***
	***REMOVED***
	if hostConfig.CgroupParent != "" && UsingSystemd(daemon.configStore) ***REMOVED***
		// CgroupParent for systemd cgroup should be named as "xxx.slice"
		if len(hostConfig.CgroupParent) <= 6 || !strings.HasSuffix(hostConfig.CgroupParent, ".slice") ***REMOVED***
			return warnings, fmt.Errorf("cgroup-parent for systemd cgroup should be a valid slice named as \"xxx.slice\"")
		***REMOVED***
	***REMOVED***
	if hostConfig.Runtime == "" ***REMOVED***
		hostConfig.Runtime = daemon.configStore.GetDefaultRuntimeName()
	***REMOVED***

	if rt := daemon.configStore.GetRuntime(hostConfig.Runtime); rt == nil ***REMOVED***
		return warnings, fmt.Errorf("Unknown runtime specified %s", hostConfig.Runtime)
	***REMOVED***

	parser := volume.NewParser(runtime.GOOS)
	for dest := range hostConfig.Tmpfs ***REMOVED***
		if err := parser.ValidateTmpfsMountDestination(dest); err != nil ***REMOVED***
			return warnings, err
		***REMOVED***
	***REMOVED***

	return warnings, nil
***REMOVED***

func (daemon *Daemon) loadRuntimes() error ***REMOVED***
	return daemon.initRuntimes(daemon.configStore.Runtimes)
***REMOVED***

func (daemon *Daemon) initRuntimes(runtimes map[string]types.Runtime) (err error) ***REMOVED***
	runtimeDir := filepath.Join(daemon.configStore.Root, "runtimes")
	// Remove old temp directory if any
	os.RemoveAll(runtimeDir + "-old")
	tmpDir, err := ioutils.TempDir(daemon.configStore.Root, "gen-runtimes")
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to get temp dir to generate runtime scripts")
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if err1 := os.RemoveAll(tmpDir); err1 != nil ***REMOVED***
				logrus.WithError(err1).WithField("dir", tmpDir).
					Warnf("failed to remove tmp dir")
			***REMOVED***
			return
		***REMOVED***

		if err = os.Rename(runtimeDir, runtimeDir+"-old"); err != nil ***REMOVED***
			return
		***REMOVED***
		if err = os.Rename(tmpDir, runtimeDir); err != nil ***REMOVED***
			err = errors.Wrapf(err, "failed to setup runtimes dir, new containers may not start")
			return
		***REMOVED***
		if err = os.RemoveAll(runtimeDir + "-old"); err != nil ***REMOVED***
			logrus.WithError(err).WithField("dir", tmpDir).
				Warnf("failed to remove old runtimes dir")
		***REMOVED***
	***REMOVED***()

	for name, rt := range runtimes ***REMOVED***
		if len(rt.Args) == 0 ***REMOVED***
			continue
		***REMOVED***

		script := filepath.Join(tmpDir, name)
		content := fmt.Sprintf("#!/bin/sh\n%s %s $@\n", rt.Path, strings.Join(rt.Args, " "))
		if err := ioutil.WriteFile(script, []byte(content), 0700); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// verifyDaemonSettings performs validation of daemon config struct
func verifyDaemonSettings(conf *config.Config) error ***REMOVED***
	// Check for mutually incompatible config options
	if conf.BridgeConfig.Iface != "" && conf.BridgeConfig.IP != "" ***REMOVED***
		return fmt.Errorf("You specified -b & --bip, mutually exclusive options. Please specify only one")
	***REMOVED***
	if !conf.BridgeConfig.EnableIPTables && !conf.BridgeConfig.InterContainerCommunication ***REMOVED***
		return fmt.Errorf("You specified --iptables=false with --icc=false. ICC=false uses iptables to function. Please set --icc or --iptables to true")
	***REMOVED***
	if !conf.BridgeConfig.EnableIPTables && conf.BridgeConfig.EnableIPMasq ***REMOVED***
		conf.BridgeConfig.EnableIPMasq = false
	***REMOVED***
	if err := VerifyCgroupDriver(conf); err != nil ***REMOVED***
		return err
	***REMOVED***
	if conf.CgroupParent != "" && UsingSystemd(conf) ***REMOVED***
		if len(conf.CgroupParent) <= 6 || !strings.HasSuffix(conf.CgroupParent, ".slice") ***REMOVED***
			return fmt.Errorf("cgroup-parent for systemd cgroup should be a valid slice named as \"xxx.slice\"")
		***REMOVED***
	***REMOVED***

	if conf.DefaultRuntime == "" ***REMOVED***
		conf.DefaultRuntime = config.StockRuntimeName
	***REMOVED***
	if conf.Runtimes == nil ***REMOVED***
		conf.Runtimes = make(map[string]types.Runtime)
	***REMOVED***
	conf.Runtimes[config.StockRuntimeName] = types.Runtime***REMOVED***Path: DefaultRuntimeName***REMOVED***

	return nil
***REMOVED***

// checkSystem validates platform-specific requirements
func checkSystem() error ***REMOVED***
	if os.Geteuid() != 0 ***REMOVED***
		return fmt.Errorf("The Docker daemon needs to be run as root")
	***REMOVED***
	return checkKernel()
***REMOVED***

// configureMaxThreads sets the Go runtime max threads threshold
// which is 90% of the kernel setting from /proc/sys/kernel/threads-max
func configureMaxThreads(config *config.Config) error ***REMOVED***
	mt, err := ioutil.ReadFile("/proc/sys/kernel/threads-max")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mtint, err := strconv.Atoi(strings.TrimSpace(string(mt)))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	maxThreads := (mtint / 100) * 90
	debug.SetMaxThreads(maxThreads)
	logrus.Debugf("Golang's threads limit set to %d", maxThreads)
	return nil
***REMOVED***

func overlaySupportsSelinux() (bool, error) ***REMOVED***
	f, err := os.Open("/proc/kallsyms")
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return false, nil
		***REMOVED***
		return false, err
	***REMOVED***
	defer f.Close()

	var symAddr, symType, symName, text string

	s := bufio.NewScanner(f)
	for s.Scan() ***REMOVED***
		if err := s.Err(); err != nil ***REMOVED***
			return false, err
		***REMOVED***

		text = s.Text()
		if _, err := fmt.Sscanf(text, "%s %s %s", &symAddr, &symType, &symName); err != nil ***REMOVED***
			return false, fmt.Errorf("Scanning '%s' failed: %s", text, err)
		***REMOVED***

		// Check for presence of symbol security_inode_copy_up.
		if symName == "security_inode_copy_up" ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***
	return false, nil
***REMOVED***

// configureKernelSecuritySupport configures and validates security support for the kernel
func configureKernelSecuritySupport(config *config.Config, driverName string) error ***REMOVED***
	if config.EnableSelinuxSupport ***REMOVED***
		if !selinuxEnabled() ***REMOVED***
			logrus.Warn("Docker could not enable SELinux on the host system")
			return nil
		***REMOVED***

		if driverName == "overlay" || driverName == "overlay2" ***REMOVED***
			// If driver is overlay or overlay2, make sure kernel
			// supports selinux with overlay.
			supported, err := overlaySupportsSelinux()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if !supported ***REMOVED***
				logrus.Warnf("SELinux is not supported with the %v graph driver on this kernel", driverName)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		selinuxSetDisabled()
	***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) initNetworkController(config *config.Config, activeSandboxes map[string]interface***REMOVED******REMOVED***) (libnetwork.NetworkController, error) ***REMOVED***
	netOptions, err := daemon.networkOptions(config, daemon.PluginStore, activeSandboxes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	controller, err := libnetwork.New(netOptions...)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error obtaining controller instance: %v", err)
	***REMOVED***

	if len(activeSandboxes) > 0 ***REMOVED***
		logrus.Info("There are old running containers, the network config will not take affect")
		return controller, nil
	***REMOVED***

	// Initialize default network on "null"
	if n, _ := controller.NetworkByName("none"); n == nil ***REMOVED***
		if _, err := controller.NewNetwork("null", "none", "", libnetwork.NetworkOptionPersist(true)); err != nil ***REMOVED***
			return nil, fmt.Errorf("Error creating default \"null\" network: %v", err)
		***REMOVED***
	***REMOVED***

	// Initialize default network on "host"
	if n, _ := controller.NetworkByName("host"); n == nil ***REMOVED***
		if _, err := controller.NewNetwork("host", "host", "", libnetwork.NetworkOptionPersist(true)); err != nil ***REMOVED***
			return nil, fmt.Errorf("Error creating default \"host\" network: %v", err)
		***REMOVED***
	***REMOVED***

	// Clear stale bridge network
	if n, err := controller.NetworkByName("bridge"); err == nil ***REMOVED***
		if err = n.Delete(); err != nil ***REMOVED***
			return nil, fmt.Errorf("could not delete the default bridge network: %v", err)
		***REMOVED***
	***REMOVED***

	if !config.DisableBridge ***REMOVED***
		// Initialize default driver "bridge"
		if err := initBridgeDriver(controller, config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		removeDefaultBridgeInterface()
	***REMOVED***

	return controller, nil
***REMOVED***

func driverOptions(config *config.Config) []nwconfig.Option ***REMOVED***
	bridgeConfig := options.Generic***REMOVED***
		"EnableIPForwarding":  config.BridgeConfig.EnableIPForward,
		"EnableIPTables":      config.BridgeConfig.EnableIPTables,
		"EnableUserlandProxy": config.BridgeConfig.EnableUserlandProxy,
		"UserlandProxyPath":   config.BridgeConfig.UserlandProxyPath***REMOVED***
	bridgeOption := options.Generic***REMOVED***netlabel.GenericData: bridgeConfig***REMOVED***

	dOptions := []nwconfig.Option***REMOVED******REMOVED***
	dOptions = append(dOptions, nwconfig.OptionDriverConfig("bridge", bridgeOption))
	return dOptions
***REMOVED***

func initBridgeDriver(controller libnetwork.NetworkController, config *config.Config) error ***REMOVED***
	bridgeName := bridge.DefaultBridgeName
	if config.BridgeConfig.Iface != "" ***REMOVED***
		bridgeName = config.BridgeConfig.Iface
	***REMOVED***
	netOption := map[string]string***REMOVED***
		bridge.BridgeName:         bridgeName,
		bridge.DefaultBridge:      strconv.FormatBool(true),
		netlabel.DriverMTU:        strconv.Itoa(config.Mtu),
		bridge.EnableIPMasquerade: strconv.FormatBool(config.BridgeConfig.EnableIPMasq),
		bridge.EnableICC:          strconv.FormatBool(config.BridgeConfig.InterContainerCommunication),
	***REMOVED***

	// --ip processing
	if config.BridgeConfig.DefaultIP != nil ***REMOVED***
		netOption[bridge.DefaultBindingIP] = config.BridgeConfig.DefaultIP.String()
	***REMOVED***

	var (
		ipamV4Conf *libnetwork.IpamConf
		ipamV6Conf *libnetwork.IpamConf
	)

	ipamV4Conf = &libnetwork.IpamConf***REMOVED***AuxAddresses: make(map[string]string)***REMOVED***

	nwList, nw6List, err := netutils.ElectInterfaceAddresses(bridgeName)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "list bridge addresses failed")
	***REMOVED***

	nw := nwList[0]
	if len(nwList) > 1 && config.BridgeConfig.FixedCIDR != "" ***REMOVED***
		_, fCIDR, err := net.ParseCIDR(config.BridgeConfig.FixedCIDR)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "parse CIDR failed")
		***REMOVED***
		// Iterate through in case there are multiple addresses for the bridge
		for _, entry := range nwList ***REMOVED***
			if fCIDR.Contains(entry.IP) ***REMOVED***
				nw = entry
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	ipamV4Conf.PreferredPool = lntypes.GetIPNetCanonical(nw).String()
	hip, _ := lntypes.GetHostPartIP(nw.IP, nw.Mask)
	if hip.IsGlobalUnicast() ***REMOVED***
		ipamV4Conf.Gateway = nw.IP.String()
	***REMOVED***

	if config.BridgeConfig.IP != "" ***REMOVED***
		ipamV4Conf.PreferredPool = config.BridgeConfig.IP
		ip, _, err := net.ParseCIDR(config.BridgeConfig.IP)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ipamV4Conf.Gateway = ip.String()
	***REMOVED*** else if bridgeName == bridge.DefaultBridgeName && ipamV4Conf.PreferredPool != "" ***REMOVED***
		logrus.Infof("Default bridge (%s) is assigned with an IP address %s. Daemon option --bip can be used to set a preferred IP address", bridgeName, ipamV4Conf.PreferredPool)
	***REMOVED***

	if config.BridgeConfig.FixedCIDR != "" ***REMOVED***
		_, fCIDR, err := net.ParseCIDR(config.BridgeConfig.FixedCIDR)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		ipamV4Conf.SubPool = fCIDR.String()
	***REMOVED***

	if config.BridgeConfig.DefaultGatewayIPv4 != nil ***REMOVED***
		ipamV4Conf.AuxAddresses["DefaultGatewayIPv4"] = config.BridgeConfig.DefaultGatewayIPv4.String()
	***REMOVED***

	var deferIPv6Alloc bool
	if config.BridgeConfig.FixedCIDRv6 != "" ***REMOVED***
		_, fCIDRv6, err := net.ParseCIDR(config.BridgeConfig.FixedCIDRv6)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// In case user has specified the daemon flag --fixed-cidr-v6 and the passed network has
		// at least 48 host bits, we need to guarantee the current behavior where the containers'
		// IPv6 addresses will be constructed based on the containers' interface MAC address.
		// We do so by telling libnetwork to defer the IPv6 address allocation for the endpoints
		// on this network until after the driver has created the endpoint and returned the
		// constructed address. Libnetwork will then reserve this address with the ipam driver.
		ones, _ := fCIDRv6.Mask.Size()
		deferIPv6Alloc = ones <= 80

		if ipamV6Conf == nil ***REMOVED***
			ipamV6Conf = &libnetwork.IpamConf***REMOVED***AuxAddresses: make(map[string]string)***REMOVED***
		***REMOVED***
		ipamV6Conf.PreferredPool = fCIDRv6.String()

		// In case the --fixed-cidr-v6 is specified and the current docker0 bridge IPv6
		// address belongs to the same network, we need to inform libnetwork about it, so
		// that it can be reserved with IPAM and it will not be given away to somebody else
		for _, nw6 := range nw6List ***REMOVED***
			if fCIDRv6.Contains(nw6.IP) ***REMOVED***
				ipamV6Conf.Gateway = nw6.IP.String()
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if config.BridgeConfig.DefaultGatewayIPv6 != nil ***REMOVED***
		if ipamV6Conf == nil ***REMOVED***
			ipamV6Conf = &libnetwork.IpamConf***REMOVED***AuxAddresses: make(map[string]string)***REMOVED***
		***REMOVED***
		ipamV6Conf.AuxAddresses["DefaultGatewayIPv6"] = config.BridgeConfig.DefaultGatewayIPv6.String()
	***REMOVED***

	v4Conf := []*libnetwork.IpamConf***REMOVED***ipamV4Conf***REMOVED***
	v6Conf := []*libnetwork.IpamConf***REMOVED******REMOVED***
	if ipamV6Conf != nil ***REMOVED***
		v6Conf = append(v6Conf, ipamV6Conf)
	***REMOVED***
	// Initialize default network on "bridge" with the same name
	_, err = controller.NewNetwork("bridge", "bridge", "",
		libnetwork.NetworkOptionEnableIPv6(config.BridgeConfig.EnableIPv6),
		libnetwork.NetworkOptionDriverOpts(netOption),
		libnetwork.NetworkOptionIpam("default", "", v4Conf, v6Conf, nil),
		libnetwork.NetworkOptionDeferIPv6Alloc(deferIPv6Alloc))
	if err != nil ***REMOVED***
		return fmt.Errorf("Error creating default \"bridge\" network: %v", err)
	***REMOVED***
	return nil
***REMOVED***

// Remove default bridge interface if present (--bridge=none use case)
func removeDefaultBridgeInterface() ***REMOVED***
	if lnk, err := netlink.LinkByName(bridge.DefaultBridgeName); err == nil ***REMOVED***
		if err := netlink.LinkDel(lnk); err != nil ***REMOVED***
			logrus.Warnf("Failed to remove bridge interface (%s): %v", bridge.DefaultBridgeName, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (daemon *Daemon) getLayerInit() func(containerfs.ContainerFS) error ***REMOVED***
	return daemon.setupInitLayer
***REMOVED***

// Parse the remapped root (user namespace) option, which can be one of:
//   username            - valid username from /etc/passwd
//   username:groupname  - valid username; valid groupname from /etc/group
//   uid                 - 32-bit unsigned int valid Linux UID value
//   uid:gid             - uid value; 32-bit unsigned int Linux GID value
//
//  If no groupname is specified, and a username is specified, an attempt
//  will be made to lookup a gid for that username as a groupname
//
//  If names are used, they are verified to exist in passwd/group
func parseRemappedRoot(usergrp string) (string, string, error) ***REMOVED***

	var (
		userID, groupID     int
		username, groupname string
	)

	idparts := strings.Split(usergrp, ":")
	if len(idparts) > 2 ***REMOVED***
		return "", "", fmt.Errorf("Invalid user/group specification in --userns-remap: %q", usergrp)
	***REMOVED***

	if uid, err := strconv.ParseInt(idparts[0], 10, 32); err == nil ***REMOVED***
		// must be a uid; take it as valid
		userID = int(uid)
		luser, err := idtools.LookupUID(userID)
		if err != nil ***REMOVED***
			return "", "", fmt.Errorf("Uid %d has no entry in /etc/passwd: %v", userID, err)
		***REMOVED***
		username = luser.Name
		if len(idparts) == 1 ***REMOVED***
			// if the uid was numeric and no gid was specified, take the uid as the gid
			groupID = userID
			lgrp, err := idtools.LookupGID(groupID)
			if err != nil ***REMOVED***
				return "", "", fmt.Errorf("Gid %d has no entry in /etc/group: %v", groupID, err)
			***REMOVED***
			groupname = lgrp.Name
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		lookupName := idparts[0]
		// special case: if the user specified "default", they want Docker to create or
		// use (after creation) the "dockremap" user/group for root remapping
		if lookupName == defaultIDSpecifier ***REMOVED***
			lookupName = defaultRemappedID
		***REMOVED***
		luser, err := idtools.LookupUser(lookupName)
		if err != nil && idparts[0] != defaultIDSpecifier ***REMOVED***
			// error if the name requested isn't the special "dockremap" ID
			return "", "", fmt.Errorf("Error during uid lookup for %q: %v", lookupName, err)
		***REMOVED*** else if err != nil ***REMOVED***
			// special case-- if the username == "default", then we have been asked
			// to create a new entry pair in /etc/***REMOVED***passwd,group***REMOVED*** for which the /etc/sub***REMOVED***uid,gid***REMOVED***
			// ranges will be used for the user and group mappings in user namespaced containers
			_, _, err := idtools.AddNamespaceRangesUser(defaultRemappedID)
			if err == nil ***REMOVED***
				return defaultRemappedID, defaultRemappedID, nil
			***REMOVED***
			return "", "", fmt.Errorf("Error during %q user creation: %v", defaultRemappedID, err)
		***REMOVED***
		username = luser.Name
		if len(idparts) == 1 ***REMOVED***
			// we only have a string username, and no group specified; look up gid from username as group
			group, err := idtools.LookupGroup(lookupName)
			if err != nil ***REMOVED***
				return "", "", fmt.Errorf("Error during gid lookup for %q: %v", lookupName, err)
			***REMOVED***
			groupname = group.Name
		***REMOVED***
	***REMOVED***

	if len(idparts) == 2 ***REMOVED***
		// groupname or gid is separately specified and must be resolved
		// to an unsigned 32-bit gid
		if gid, err := strconv.ParseInt(idparts[1], 10, 32); err == nil ***REMOVED***
			// must be a gid, take it as valid
			groupID = int(gid)
			lgrp, err := idtools.LookupGID(groupID)
			if err != nil ***REMOVED***
				return "", "", fmt.Errorf("Gid %d has no entry in /etc/passwd: %v", groupID, err)
			***REMOVED***
			groupname = lgrp.Name
		***REMOVED*** else ***REMOVED***
			// not a number; attempt a lookup
			if _, err := idtools.LookupGroup(idparts[1]); err != nil ***REMOVED***
				return "", "", fmt.Errorf("Error during groupname lookup for %q: %v", idparts[1], err)
			***REMOVED***
			groupname = idparts[1]
		***REMOVED***
	***REMOVED***
	return username, groupname, nil
***REMOVED***

func setupRemappedRoot(config *config.Config) (*idtools.IDMappings, error) ***REMOVED***
	if runtime.GOOS != "linux" && config.RemappedRoot != "" ***REMOVED***
		return nil, fmt.Errorf("User namespaces are only supported on Linux")
	***REMOVED***

	// if the daemon was started with remapped root option, parse
	// the config option to the int uid,gid values
	if config.RemappedRoot != "" ***REMOVED***
		username, groupname, err := parseRemappedRoot(config.RemappedRoot)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if username == "root" ***REMOVED***
			// Cannot setup user namespaces with a 1-to-1 mapping; "--root=0:0" is a no-op
			// effectively
			logrus.Warn("User namespaces: root cannot be remapped with itself; user namespaces are OFF")
			return &idtools.IDMappings***REMOVED******REMOVED***, nil
		***REMOVED***
		logrus.Infof("User namespaces: ID ranges will be mapped to subuid/subgid ranges of: %s:%s", username, groupname)
		// update remapped root setting now that we have resolved them to actual names
		config.RemappedRoot = fmt.Sprintf("%s:%s", username, groupname)

		mappings, err := idtools.NewIDMappings(username, groupname)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "Can't create ID mappings: %v")
		***REMOVED***
		return mappings, nil
	***REMOVED***
	return &idtools.IDMappings***REMOVED******REMOVED***, nil
***REMOVED***

func setupDaemonRoot(config *config.Config, rootDir string, rootIDs idtools.IDPair) error ***REMOVED***
	config.Root = rootDir
	// the docker root metadata directory needs to have execute permissions for all users (g+x,o+x)
	// so that syscalls executing as non-root, operating on subdirectories of the graph root
	// (e.g. mounted layers of a container) can traverse this path.
	// The user namespace support will create subdirectories for the remapped root host uid:gid
	// pair owned by that same uid:gid pair for proper write access to those needed metadata and
	// layer content subtrees.
	if _, err := os.Stat(rootDir); err == nil ***REMOVED***
		// root current exists; verify the access bits are correct by setting them
		if err = os.Chmod(rootDir, 0711); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if os.IsNotExist(err) ***REMOVED***
		// no root exists yet, create it 0711 with root:root ownership
		if err := os.MkdirAll(rootDir, 0711); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// if user namespaces are enabled we will create a subtree underneath the specified root
	// with any/all specified remapped root uid/gid options on the daemon creating
	// a new subdirectory with ownership set to the remapped uid/gid (so as to allow
	// `chdir()` to work for containers namespaced to that uid/gid)
	if config.RemappedRoot != "" ***REMOVED***
		config.Root = filepath.Join(rootDir, fmt.Sprintf("%d.%d", rootIDs.UID, rootIDs.GID))
		logrus.Debugf("Creating user namespaced daemon root: %s", config.Root)
		// Create the root directory if it doesn't exist
		if err := idtools.MkdirAllAndChown(config.Root, 0700, rootIDs); err != nil ***REMOVED***
			return fmt.Errorf("Cannot create daemon root: %s: %v", config.Root, err)
		***REMOVED***
		// we also need to verify that any pre-existing directories in the path to
		// the graphroot won't block access to remapped root--if any pre-existing directory
		// has strict permissions that don't allow "x", container start will fail, so
		// better to warn and fail now
		dirPath := config.Root
		for ***REMOVED***
			dirPath = filepath.Dir(dirPath)
			if dirPath == "/" ***REMOVED***
				break
			***REMOVED***
			if !idtools.CanAccess(dirPath, rootIDs) ***REMOVED***
				return fmt.Errorf("a subdirectory in your graphroot path (%s) restricts access to the remapped root uid/gid; please fix by allowing 'o+x' permissions on existing directories", config.Root)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := ensureSharedOrSlave(config.Root); err != nil ***REMOVED***
		if err := mount.MakeShared(config.Root); err != nil ***REMOVED***
			logrus.WithError(err).WithField("dir", config.Root).Warn("Could not set daemon root propagation to shared, this is not generally critical but may cause some functionality to not work or fallback to less desirable behavior")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// registerLinks writes the links to a file.
func (daemon *Daemon) registerLinks(container *container.Container, hostConfig *containertypes.HostConfig) error ***REMOVED***
	if hostConfig == nil || hostConfig.NetworkMode.IsUserDefined() ***REMOVED***
		return nil
	***REMOVED***

	for _, l := range hostConfig.Links ***REMOVED***
		name, alias, err := opts.ParseLink(l)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		child, err := daemon.GetContainer(name)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "could not get container for %s", name)
		***REMOVED***
		for child.HostConfig.NetworkMode.IsContainer() ***REMOVED***
			parts := strings.SplitN(string(child.HostConfig.NetworkMode), ":", 2)
			child, err = daemon.GetContainer(parts[1])
			if err != nil ***REMOVED***
				return errors.Wrapf(err, "Could not get container for %s", parts[1])
			***REMOVED***
		***REMOVED***
		if child.HostConfig.NetworkMode.IsHost() ***REMOVED***
			return runconfig.ErrConflictHostNetworkAndLinks
		***REMOVED***
		if err := daemon.registerLink(container, child, alias); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// After we load all the links into the daemon
	// set them to nil on the hostconfig
	_, err := container.WriteHostConfig()
	return err
***REMOVED***

// conditionalMountOnStart is a platform specific helper function during the
// container start to call mount.
func (daemon *Daemon) conditionalMountOnStart(container *container.Container) error ***REMOVED***
	return daemon.Mount(container)
***REMOVED***

// conditionalUnmountOnCleanup is a platform specific helper function called
// during the cleanup of a container to unmount.
func (daemon *Daemon) conditionalUnmountOnCleanup(container *container.Container) error ***REMOVED***
	return daemon.Unmount(container)
***REMOVED***

func copyBlkioEntry(entries []*containerd_cgroups.BlkIOEntry) []types.BlkioStatEntry ***REMOVED***
	out := make([]types.BlkioStatEntry, len(entries))
	for i, re := range entries ***REMOVED***
		out[i] = types.BlkioStatEntry***REMOVED***
			Major: re.Major,
			Minor: re.Minor,
			Op:    re.Op,
			Value: re.Value,
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***

func (daemon *Daemon) stats(c *container.Container) (*types.StatsJSON, error) ***REMOVED***
	if !c.IsRunning() ***REMOVED***
		return nil, errNotRunning(c.ID)
	***REMOVED***
	cs, err := daemon.containerd.Stats(context.Background(), c.ID)
	if err != nil ***REMOVED***
		if strings.Contains(err.Error(), "container not found") ***REMOVED***
			return nil, containerNotFound(c.ID)
		***REMOVED***
		return nil, err
	***REMOVED***
	s := &types.StatsJSON***REMOVED******REMOVED***
	s.Read = cs.Read
	stats := cs.Metrics
	if stats.Blkio != nil ***REMOVED***
		s.BlkioStats = types.BlkioStats***REMOVED***
			IoServiceBytesRecursive: copyBlkioEntry(stats.Blkio.IoServiceBytesRecursive),
			IoServicedRecursive:     copyBlkioEntry(stats.Blkio.IoServicedRecursive),
			IoQueuedRecursive:       copyBlkioEntry(stats.Blkio.IoQueuedRecursive),
			IoServiceTimeRecursive:  copyBlkioEntry(stats.Blkio.IoServiceTimeRecursive),
			IoWaitTimeRecursive:     copyBlkioEntry(stats.Blkio.IoWaitTimeRecursive),
			IoMergedRecursive:       copyBlkioEntry(stats.Blkio.IoMergedRecursive),
			IoTimeRecursive:         copyBlkioEntry(stats.Blkio.IoTimeRecursive),
			SectorsRecursive:        copyBlkioEntry(stats.Blkio.SectorsRecursive),
		***REMOVED***
	***REMOVED***
	if stats.CPU != nil ***REMOVED***
		s.CPUStats = types.CPUStats***REMOVED***
			CPUUsage: types.CPUUsage***REMOVED***
				TotalUsage:        stats.CPU.Usage.Total,
				PercpuUsage:       stats.CPU.Usage.PerCPU,
				UsageInKernelmode: stats.CPU.Usage.Kernel,
				UsageInUsermode:   stats.CPU.Usage.User,
			***REMOVED***,
			ThrottlingData: types.ThrottlingData***REMOVED***
				Periods:          stats.CPU.Throttling.Periods,
				ThrottledPeriods: stats.CPU.Throttling.ThrottledPeriods,
				ThrottledTime:    stats.CPU.Throttling.ThrottledTime,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	if stats.Memory != nil ***REMOVED***
		raw := make(map[string]uint64)
		raw["cache"] = stats.Memory.Cache
		raw["rss"] = stats.Memory.RSS
		raw["rss_huge"] = stats.Memory.RSSHuge
		raw["mapped_file"] = stats.Memory.MappedFile
		raw["dirty"] = stats.Memory.Dirty
		raw["writeback"] = stats.Memory.Writeback
		raw["pgpgin"] = stats.Memory.PgPgIn
		raw["pgpgout"] = stats.Memory.PgPgOut
		raw["pgfault"] = stats.Memory.PgFault
		raw["pgmajfault"] = stats.Memory.PgMajFault
		raw["inactive_anon"] = stats.Memory.InactiveAnon
		raw["active_anon"] = stats.Memory.ActiveAnon
		raw["inactive_file"] = stats.Memory.InactiveFile
		raw["active_file"] = stats.Memory.ActiveFile
		raw["unevictable"] = stats.Memory.Unevictable
		raw["hierarchical_memory_limit"] = stats.Memory.HierarchicalMemoryLimit
		raw["hierarchical_memsw_limit"] = stats.Memory.HierarchicalSwapLimit
		raw["total_cache"] = stats.Memory.TotalCache
		raw["total_rss"] = stats.Memory.TotalRSS
		raw["total_rss_huge"] = stats.Memory.TotalRSSHuge
		raw["total_mapped_file"] = stats.Memory.TotalMappedFile
		raw["total_dirty"] = stats.Memory.TotalDirty
		raw["total_writeback"] = stats.Memory.TotalWriteback
		raw["total_pgpgin"] = stats.Memory.TotalPgPgIn
		raw["total_pgpgout"] = stats.Memory.TotalPgPgOut
		raw["total_pgfault"] = stats.Memory.TotalPgFault
		raw["total_pgmajfault"] = stats.Memory.TotalPgMajFault
		raw["total_inactive_anon"] = stats.Memory.TotalInactiveAnon
		raw["total_active_anon"] = stats.Memory.TotalActiveAnon
		raw["total_inactive_file"] = stats.Memory.TotalInactiveFile
		raw["total_active_file"] = stats.Memory.TotalActiveFile
		raw["total_unevictable"] = stats.Memory.TotalUnevictable

		if stats.Memory.Usage != nil ***REMOVED***
			s.MemoryStats = types.MemoryStats***REMOVED***
				Stats:    raw,
				Usage:    stats.Memory.Usage.Usage,
				MaxUsage: stats.Memory.Usage.Max,
				Limit:    stats.Memory.Usage.Limit,
				Failcnt:  stats.Memory.Usage.Failcnt,
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			s.MemoryStats = types.MemoryStats***REMOVED***
				Stats: raw,
			***REMOVED***
		***REMOVED***

		// if the container does not set memory limit, use the machineMemory
		if s.MemoryStats.Limit > daemon.machineMemory && daemon.machineMemory > 0 ***REMOVED***
			s.MemoryStats.Limit = daemon.machineMemory
		***REMOVED***
	***REMOVED***

	if stats.Pids != nil ***REMOVED***
		s.PidsStats = types.PidsStats***REMOVED***
			Current: stats.Pids.Current,
			Limit:   stats.Pids.Limit,
		***REMOVED***
	***REMOVED***

	return s, nil
***REMOVED***

// setDefaultIsolation determines the default isolation mode for the
// daemon to run in. This is only applicable on Windows
func (daemon *Daemon) setDefaultIsolation() error ***REMOVED***
	return nil
***REMOVED***

func rootFSToAPIType(rootfs *image.RootFS) types.RootFS ***REMOVED***
	var layers []string
	for _, l := range rootfs.DiffIDs ***REMOVED***
		layers = append(layers, l.String())
	***REMOVED***
	return types.RootFS***REMOVED***
		Type:   rootfs.Type,
		Layers: layers,
	***REMOVED***
***REMOVED***

// setupDaemonProcess sets various settings for the daemon's process
func setupDaemonProcess(config *config.Config) error ***REMOVED***
	// setup the daemons oom_score_adj
	if err := setupOOMScoreAdj(config.OOMScoreAdjust); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := setMayDetachMounts(); err != nil ***REMOVED***
		logrus.WithError(err).Warn("Could not set may_detach_mounts kernel parameter")
	***REMOVED***
	return nil
***REMOVED***

// This is used to allow removal of mountpoints that may be mounted in other
// namespaces on RHEL based kernels starting from RHEL 7.4.
// Without this setting, removals on these RHEL based kernels may fail with
// "device or resource busy".
// This setting is not available in upstream kernels as it is not configurable,
// but has been in the upstream kernels since 3.15.
func setMayDetachMounts() error ***REMOVED***
	f, err := os.OpenFile("/proc/sys/fs/may_detach_mounts", os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return errors.Wrap(err, "error opening may_detach_mounts kernel config file")
	***REMOVED***
	defer f.Close()

	_, err = f.WriteString("1")
	if os.IsPermission(err) ***REMOVED***
		// Setting may_detach_mounts does not work in an
		// unprivileged container. Ignore the error, but log
		// it if we appear not to be in that situation.
		if !rsystem.RunningInUserNS() ***REMOVED***
			logrus.Debugf("Permission denied writing %q to /proc/sys/fs/may_detach_mounts", "1")
		***REMOVED***
		return nil
	***REMOVED***
	return err
***REMOVED***

func setupOOMScoreAdj(score int) error ***REMOVED***
	f, err := os.OpenFile("/proc/self/oom_score_adj", os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()
	stringScore := strconv.Itoa(score)
	_, err = f.WriteString(stringScore)
	if os.IsPermission(err) ***REMOVED***
		// Setting oom_score_adj does not work in an
		// unprivileged container. Ignore the error, but log
		// it if we appear not to be in that situation.
		if !rsystem.RunningInUserNS() ***REMOVED***
			logrus.Debugf("Permission denied writing %q to /proc/self/oom_score_adj", stringScore)
		***REMOVED***
		return nil
	***REMOVED***

	return err
***REMOVED***

func (daemon *Daemon) initCgroupsPath(path string) error ***REMOVED***
	if path == "/" || path == "." ***REMOVED***
		return nil
	***REMOVED***

	if daemon.configStore.CPURealtimePeriod == 0 && daemon.configStore.CPURealtimeRuntime == 0 ***REMOVED***
		return nil
	***REMOVED***

	// Recursively create cgroup to ensure that the system and all parent cgroups have values set
	// for the period and runtime as this limits what the children can be set to.
	daemon.initCgroupsPath(filepath.Dir(path))

	mnt, root, err := cgroups.FindCgroupMountpointAndRoot("cpu")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// When docker is run inside docker, the root is based of the host cgroup.
	// Should this be handled in runc/libcontainer/cgroups ?
	if strings.HasPrefix(root, "/docker/") ***REMOVED***
		root = "/"
	***REMOVED***

	path = filepath.Join(mnt, root, path)
	sysinfo := sysinfo.New(true)
	if err := maybeCreateCPURealTimeFile(sysinfo.CPURealtimePeriod, daemon.configStore.CPURealtimePeriod, "cpu.rt_period_us", path); err != nil ***REMOVED***
		return err
	***REMOVED***
	return maybeCreateCPURealTimeFile(sysinfo.CPURealtimeRuntime, daemon.configStore.CPURealtimeRuntime, "cpu.rt_runtime_us", path)
***REMOVED***

func maybeCreateCPURealTimeFile(sysinfoPresent bool, configValue int64, file string, path string) error ***REMOVED***
	if sysinfoPresent && configValue != 0 ***REMOVED***
		if err := os.MkdirAll(path, 0755); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := ioutil.WriteFile(filepath.Join(path, file), []byte(strconv.FormatInt(configValue, 10)), 0700); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) setupSeccompProfile() error ***REMOVED***
	if daemon.configStore.SeccompProfile != "" ***REMOVED***
		daemon.seccompProfilePath = daemon.configStore.SeccompProfile
		b, err := ioutil.ReadFile(daemon.configStore.SeccompProfile)
		if err != nil ***REMOVED***
			return fmt.Errorf("opening seccomp profile (%s) failed: %v", daemon.configStore.SeccompProfile, err)
		***REMOVED***
		daemon.seccompProfile = b
	***REMOVED***
	return nil
***REMOVED***
