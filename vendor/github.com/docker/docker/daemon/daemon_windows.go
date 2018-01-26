package daemon

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/image"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/parsers"
	"github.com/docker/docker/pkg/platform"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/runconfig"
	"github.com/docker/libnetwork"
	nwconfig "github.com/docker/libnetwork/config"
	"github.com/docker/libnetwork/datastore"
	winlibnetwork "github.com/docker/libnetwork/drivers/windows"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/options"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	defaultNetworkSpace  = "172.16.0.0/12"
	platformSupported    = true
	windowsMinCPUShares  = 1
	windowsMaxCPUShares  = 10000
	windowsMinCPUPercent = 1
	windowsMaxCPUPercent = 100
)

// Windows has no concept of an execution state directory. So use config.Root here.
func getPluginExecRoot(root string) string ***REMOVED***
	return filepath.Join(root, "plugins")
***REMOVED***

func (daemon *Daemon) parseSecurityOpt(container *container.Container, hostConfig *containertypes.HostConfig) error ***REMOVED***
	return parseSecurityOpt(container, hostConfig)
***REMOVED***

func parseSecurityOpt(container *container.Container, config *containertypes.HostConfig) error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) getLayerInit() func(containerfs.ContainerFS) error ***REMOVED***
	return nil
***REMOVED***

func checkKernel() error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) getCgroupDriver() string ***REMOVED***
	return ""
***REMOVED***

// adaptContainerSettings is called during container creation to modify any
// settings necessary in the HostConfig structure.
func (daemon *Daemon) adaptContainerSettings(hostConfig *containertypes.HostConfig, adjustCPUShares bool) error ***REMOVED***
	if hostConfig == nil ***REMOVED***
		return nil
	***REMOVED***

	return nil
***REMOVED***

func verifyContainerResources(resources *containertypes.Resources, isHyperv bool) ([]string, error) ***REMOVED***
	warnings := []string***REMOVED******REMOVED***
	fixMemorySwappiness(resources)
	if !isHyperv ***REMOVED***
		// The processor resource controls are mutually exclusive on
		// Windows Server Containers, the order of precedence is
		// CPUCount first, then CPUShares, and CPUPercent last.
		if resources.CPUCount > 0 ***REMOVED***
			if resources.CPUShares > 0 ***REMOVED***
				warnings = append(warnings, "Conflicting options: CPU count takes priority over CPU shares on Windows Server Containers. CPU shares discarded")
				logrus.Warn("Conflicting options: CPU count takes priority over CPU shares on Windows Server Containers. CPU shares discarded")
				resources.CPUShares = 0
			***REMOVED***
			if resources.CPUPercent > 0 ***REMOVED***
				warnings = append(warnings, "Conflicting options: CPU count takes priority over CPU percent on Windows Server Containers. CPU percent discarded")
				logrus.Warn("Conflicting options: CPU count takes priority over CPU percent on Windows Server Containers. CPU percent discarded")
				resources.CPUPercent = 0
			***REMOVED***
		***REMOVED*** else if resources.CPUShares > 0 ***REMOVED***
			if resources.CPUPercent > 0 ***REMOVED***
				warnings = append(warnings, "Conflicting options: CPU shares takes priority over CPU percent on Windows Server Containers. CPU percent discarded")
				logrus.Warn("Conflicting options: CPU shares takes priority over CPU percent on Windows Server Containers. CPU percent discarded")
				resources.CPUPercent = 0
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if resources.CPUShares < 0 || resources.CPUShares > windowsMaxCPUShares ***REMOVED***
		return warnings, fmt.Errorf("range of CPUShares is from %d to %d", windowsMinCPUShares, windowsMaxCPUShares)
	***REMOVED***
	if resources.CPUPercent < 0 || resources.CPUPercent > windowsMaxCPUPercent ***REMOVED***
		return warnings, fmt.Errorf("range of CPUPercent is from %d to %d", windowsMinCPUPercent, windowsMaxCPUPercent)
	***REMOVED***
	if resources.CPUCount < 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid CPUCount: CPUCount cannot be negative")
	***REMOVED***

	if resources.NanoCPUs > 0 && resources.CPUPercent > 0 ***REMOVED***
		return warnings, fmt.Errorf("conflicting options: Nano CPUs and CPU Percent cannot both be set")
	***REMOVED***
	if resources.NanoCPUs > 0 && resources.CPUShares > 0 ***REMOVED***
		return warnings, fmt.Errorf("conflicting options: Nano CPUs and CPU Shares cannot both be set")
	***REMOVED***
	// The precision we could get is 0.01, because on Windows we have to convert to CPUPercent.
	// We don't set the lower limit here and it is up to the underlying platform (e.g., Windows) to return an error.
	if resources.NanoCPUs < 0 || resources.NanoCPUs > int64(sysinfo.NumCPU())*1e9 ***REMOVED***
		return warnings, fmt.Errorf("range of CPUs is from 0.01 to %d.00, as there are only %d CPUs available", sysinfo.NumCPU(), sysinfo.NumCPU())
	***REMOVED***

	osv := system.GetOSVersion()
	if resources.NanoCPUs > 0 && isHyperv && osv.Build < 16175 ***REMOVED***
		leftoverNanoCPUs := resources.NanoCPUs % 1e9
		if leftoverNanoCPUs != 0 && resources.NanoCPUs > 1e9 ***REMOVED***
			resources.NanoCPUs = ((resources.NanoCPUs + 1e9/2) / 1e9) * 1e9
			warningString := fmt.Sprintf("Your current OS version does not support Hyper-V containers with NanoCPUs greater than 1000000000 but not divisible by 1000000000. NanoCPUs rounded to %d", resources.NanoCPUs)
			warnings = append(warnings, warningString)
			logrus.Warn(warningString)
		***REMOVED***
	***REMOVED***

	if len(resources.BlkioDeviceReadBps) > 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support BlkioDeviceReadBps")
	***REMOVED***
	if len(resources.BlkioDeviceReadIOps) > 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support BlkioDeviceReadIOps")
	***REMOVED***
	if len(resources.BlkioDeviceWriteBps) > 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support BlkioDeviceWriteBps")
	***REMOVED***
	if len(resources.BlkioDeviceWriteIOps) > 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support BlkioDeviceWriteIOps")
	***REMOVED***
	if resources.BlkioWeight > 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support BlkioWeight")
	***REMOVED***
	if len(resources.BlkioWeightDevice) > 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support BlkioWeightDevice")
	***REMOVED***
	if resources.CgroupParent != "" ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support CgroupParent")
	***REMOVED***
	if resources.CPUPeriod != 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support CPUPeriod")
	***REMOVED***
	if resources.CpusetCpus != "" ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support CpusetCpus")
	***REMOVED***
	if resources.CpusetMems != "" ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support CpusetMems")
	***REMOVED***
	if resources.KernelMemory != 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support KernelMemory")
	***REMOVED***
	if resources.MemoryReservation != 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support MemoryReservation")
	***REMOVED***
	if resources.MemorySwap != 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support MemorySwap")
	***REMOVED***
	if resources.MemorySwappiness != nil ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support MemorySwappiness")
	***REMOVED***
	if resources.OomKillDisable != nil && *resources.OomKillDisable ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support OomKillDisable")
	***REMOVED***
	if resources.PidsLimit != 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support PidsLimit")
	***REMOVED***
	if len(resources.Ulimits) != 0 ***REMOVED***
		return warnings, fmt.Errorf("invalid option: Windows does not support Ulimits")
	***REMOVED***
	return warnings, nil
***REMOVED***

// verifyPlatformContainerSettings performs platform-specific validation of the
// hostconfig and config structures.
func verifyPlatformContainerSettings(daemon *Daemon, hostConfig *containertypes.HostConfig, config *containertypes.Config, update bool) ([]string, error) ***REMOVED***
	warnings := []string***REMOVED******REMOVED***

	hyperv := daemon.runAsHyperVContainer(hostConfig)
	if !hyperv && system.IsWindowsClient() && !system.IsIoTCore() ***REMOVED***
		// @engine maintainers. This block should not be removed. It partially enforces licensing
		// restrictions on Windows. Ping @jhowardmsft if there are concerns or PRs to change this.
		return warnings, fmt.Errorf("Windows client operating systems only support Hyper-V containers")
	***REMOVED***

	w, err := verifyContainerResources(&hostConfig.Resources, hyperv)
	warnings = append(warnings, w...)
	return warnings, err
***REMOVED***

// verifyDaemonSettings performs validation of daemon config struct
func verifyDaemonSettings(config *config.Config) error ***REMOVED***
	return nil
***REMOVED***

// checkSystem validates platform-specific requirements
func checkSystem() error ***REMOVED***
	// Validate the OS version. Note that docker.exe must be manifested for this
	// call to return the correct version.
	osv := system.GetOSVersion()
	if osv.MajorVersion < 10 ***REMOVED***
		return fmt.Errorf("This version of Windows does not support the docker daemon")
	***REMOVED***
	if osv.Build < 14393 ***REMOVED***
		return fmt.Errorf("The docker daemon requires build 14393 or later of Windows Server 2016 or Windows 10")
	***REMOVED***

	vmcompute := windows.NewLazySystemDLL("vmcompute.dll")
	if vmcompute.Load() != nil ***REMOVED***
		return fmt.Errorf("failed to load vmcompute.dll, ensure that the Containers feature is installed")
	***REMOVED***

	// Ensure that the required Host Network Service and vmcompute services
	// are running. Docker will fail in unexpected ways if this is not present.
	var requiredServices = []string***REMOVED***"hns", "vmcompute"***REMOVED***
	if err := ensureServicesInstalled(requiredServices); err != nil ***REMOVED***
		return errors.Wrap(err, "a required service is not installed, ensure the Containers feature is installed")
	***REMOVED***

	return nil
***REMOVED***

func ensureServicesInstalled(services []string) error ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer m.Disconnect()
	for _, service := range services ***REMOVED***
		s, err := m.OpenService(service)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to open service %s", service)
		***REMOVED***
		s.Close()
	***REMOVED***
	return nil
***REMOVED***

// configureKernelSecuritySupport configures and validate security support for the kernel
func configureKernelSecuritySupport(config *config.Config, driverName string) error ***REMOVED***
	return nil
***REMOVED***

// configureMaxThreads sets the Go runtime max threads threshold
func configureMaxThreads(config *config.Config) error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) initNetworkController(config *config.Config, activeSandboxes map[string]interface***REMOVED******REMOVED***) (libnetwork.NetworkController, error) ***REMOVED***
	netOptions, err := daemon.networkOptions(config, nil, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	controller, err := libnetwork.New(netOptions...)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error obtaining controller instance: %v", err)
	***REMOVED***

	hnsresponse, err := hcsshim.HNSListNetworkRequest("GET", "", "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Remove networks not present in HNS
	for _, v := range controller.Networks() ***REMOVED***
		options := v.Info().DriverOptions()
		hnsid := options[winlibnetwork.HNSID]
		found := false

		for _, v := range hnsresponse ***REMOVED***
			if v.Id == hnsid ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***

		if !found ***REMOVED***
			// global networks should not be deleted by local HNS
			if v.Info().Scope() != datastore.GlobalScope ***REMOVED***
				err = v.Delete()
				if err != nil ***REMOVED***
					logrus.Errorf("Error occurred when removing network %v", err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	_, err = controller.NewNetwork("null", "none", "", libnetwork.NetworkOptionPersist(false))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defaultNetworkExists := false

	if network, err := controller.NetworkByName(runconfig.DefaultDaemonNetworkMode().NetworkName()); err == nil ***REMOVED***
		options := network.Info().DriverOptions()
		for _, v := range hnsresponse ***REMOVED***
			if options[winlibnetwork.HNSID] == v.Id ***REMOVED***
				defaultNetworkExists = true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// discover and add HNS networks to windows
	// network that exist are removed and added again
	for _, v := range hnsresponse ***REMOVED***
		if strings.ToLower(v.Type) == "private" ***REMOVED***
			continue // workaround for HNS reporting unsupported networks
		***REMOVED***
		var n libnetwork.Network
		s := func(current libnetwork.Network) bool ***REMOVED***
			options := current.Info().DriverOptions()
			if options[winlibnetwork.HNSID] == v.Id ***REMOVED***
				n = current
				return true
			***REMOVED***
			return false
		***REMOVED***

		controller.WalkNetworks(s)

		drvOptions := make(map[string]string)

		if n != nil ***REMOVED***
			// global networks should not be deleted by local HNS
			if n.Info().Scope() == datastore.GlobalScope ***REMOVED***
				continue
			***REMOVED***
			v.Name = n.Name()
			// This will not cause network delete from HNS as the network
			// is not yet populated in the libnetwork windows driver

			// restore option if it existed before
			drvOptions = n.Info().DriverOptions()
			n.Delete()
		***REMOVED***
		netOption := map[string]string***REMOVED***
			winlibnetwork.NetworkName: v.Name,
			winlibnetwork.HNSID:       v.Id,
		***REMOVED***

		// add persisted driver options
		for k, v := range drvOptions ***REMOVED***
			if k != winlibnetwork.NetworkName && k != winlibnetwork.HNSID ***REMOVED***
				netOption[k] = v
			***REMOVED***
		***REMOVED***

		v4Conf := []*libnetwork.IpamConf***REMOVED******REMOVED***
		for _, subnet := range v.Subnets ***REMOVED***
			ipamV4Conf := libnetwork.IpamConf***REMOVED******REMOVED***
			ipamV4Conf.PreferredPool = subnet.AddressPrefix
			ipamV4Conf.Gateway = subnet.GatewayAddress
			v4Conf = append(v4Conf, &ipamV4Conf)
		***REMOVED***

		name := v.Name

		// If there is no nat network create one from the first NAT network
		// encountered if it doesn't already exist
		if !defaultNetworkExists &&
			runconfig.DefaultDaemonNetworkMode() == containertypes.NetworkMode(strings.ToLower(v.Type)) &&
			n == nil ***REMOVED***
			name = runconfig.DefaultDaemonNetworkMode().NetworkName()
			defaultNetworkExists = true
		***REMOVED***

		v6Conf := []*libnetwork.IpamConf***REMOVED******REMOVED***
		_, err := controller.NewNetwork(strings.ToLower(v.Type), name, "",
			libnetwork.NetworkOptionGeneric(options.Generic***REMOVED***
				netlabel.GenericData: netOption,
			***REMOVED***),
			libnetwork.NetworkOptionIpam("default", "", v4Conf, v6Conf, nil),
		)

		if err != nil ***REMOVED***
			logrus.Errorf("Error occurred when creating network %v", err)
		***REMOVED***
	***REMOVED***

	if !config.DisableBridge ***REMOVED***
		// Initialize default driver "bridge"
		if err := initBridgeDriver(controller, config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return controller, nil
***REMOVED***

func initBridgeDriver(controller libnetwork.NetworkController, config *config.Config) error ***REMOVED***
	if _, err := controller.NetworkByName(runconfig.DefaultDaemonNetworkMode().NetworkName()); err == nil ***REMOVED***
		return nil
	***REMOVED***

	netOption := map[string]string***REMOVED***
		winlibnetwork.NetworkName: runconfig.DefaultDaemonNetworkMode().NetworkName(),
	***REMOVED***

	var ipamOption libnetwork.NetworkOption
	var subnetPrefix string

	if config.BridgeConfig.FixedCIDR != "" ***REMOVED***
		subnetPrefix = config.BridgeConfig.FixedCIDR
	***REMOVED*** else ***REMOVED***
		// TP5 doesn't support properly detecting subnet
		osv := system.GetOSVersion()
		if osv.Build < 14360 ***REMOVED***
			subnetPrefix = defaultNetworkSpace
		***REMOVED***
	***REMOVED***

	if subnetPrefix != "" ***REMOVED***
		ipamV4Conf := libnetwork.IpamConf***REMOVED******REMOVED***
		ipamV4Conf.PreferredPool = subnetPrefix
		v4Conf := []*libnetwork.IpamConf***REMOVED***&ipamV4Conf***REMOVED***
		v6Conf := []*libnetwork.IpamConf***REMOVED******REMOVED***
		ipamOption = libnetwork.NetworkOptionIpam("default", "", v4Conf, v6Conf, nil)
	***REMOVED***

	_, err := controller.NewNetwork(string(runconfig.DefaultDaemonNetworkMode()), runconfig.DefaultDaemonNetworkMode().NetworkName(), "",
		libnetwork.NetworkOptionGeneric(options.Generic***REMOVED***
			netlabel.GenericData: netOption,
		***REMOVED***),
		ipamOption,
	)

	if err != nil ***REMOVED***
		return fmt.Errorf("Error creating default network: %v", err)
	***REMOVED***

	return nil
***REMOVED***

// registerLinks sets up links between containers and writes the
// configuration out for persistence. As of Windows TP4, links are not supported.
func (daemon *Daemon) registerLinks(container *container.Container, hostConfig *containertypes.HostConfig) error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) cleanupMountsByID(in string) error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) cleanupMounts() error ***REMOVED***
	return nil
***REMOVED***

func setupRemappedRoot(config *config.Config) (*idtools.IDMappings, error) ***REMOVED***
	return &idtools.IDMappings***REMOVED******REMOVED***, nil
***REMOVED***

func setupDaemonRoot(config *config.Config, rootDir string, rootIDs idtools.IDPair) error ***REMOVED***
	config.Root = rootDir
	// Create the root directory if it doesn't exists
	if err := system.MkdirAllWithACL(config.Root, 0, system.SddlAdministratorsLocalSystem); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// runasHyperVContainer returns true if we are going to run as a Hyper-V container
func (daemon *Daemon) runAsHyperVContainer(hostConfig *containertypes.HostConfig) bool ***REMOVED***
	if hostConfig.Isolation.IsDefault() ***REMOVED***
		// Container is set to use the default, so take the default from the daemon configuration
		return daemon.defaultIsolation.IsHyperV()
	***REMOVED***

	// Container is requesting an isolation mode. Honour it.
	return hostConfig.Isolation.IsHyperV()

***REMOVED***

// conditionalMountOnStart is a platform specific helper function during the
// container start to call mount.
func (daemon *Daemon) conditionalMountOnStart(container *container.Container) error ***REMOVED***
	// Bail out now for Linux containers. We cannot mount the containers filesystem on the
	// host as it is a non-Windows filesystem.
	if system.LCOWSupported() && container.OS != "windows" ***REMOVED***
		return nil
	***REMOVED***

	// We do not mount if a Hyper-V container as it needs to be mounted inside the
	// utility VM, not the host.
	if !daemon.runAsHyperVContainer(container.HostConfig) ***REMOVED***
		return daemon.Mount(container)
	***REMOVED***
	return nil
***REMOVED***

// conditionalUnmountOnCleanup is a platform specific helper function called
// during the cleanup of a container to unmount.
func (daemon *Daemon) conditionalUnmountOnCleanup(container *container.Container) error ***REMOVED***
	// Bail out now for Linux containers
	if system.LCOWSupported() && container.OS != "windows" ***REMOVED***
		return nil
	***REMOVED***

	// We do not unmount if a Hyper-V container
	if !daemon.runAsHyperVContainer(container.HostConfig) ***REMOVED***
		return daemon.Unmount(container)
	***REMOVED***
	return nil
***REMOVED***

func driverOptions(config *config.Config) []nwconfig.Option ***REMOVED***
	return []nwconfig.Option***REMOVED******REMOVED***
***REMOVED***

func (daemon *Daemon) stats(c *container.Container) (*types.StatsJSON, error) ***REMOVED***
	if !c.IsRunning() ***REMOVED***
		return nil, errNotRunning(c.ID)
	***REMOVED***

	// Obtain the stats from HCS via libcontainerd
	stats, err := daemon.containerd.Stats(context.Background(), c.ID)
	if err != nil ***REMOVED***
		if strings.Contains(err.Error(), "container not found") ***REMOVED***
			return nil, containerNotFound(c.ID)
		***REMOVED***
		return nil, err
	***REMOVED***

	// Start with an empty structure
	s := &types.StatsJSON***REMOVED******REMOVED***
	s.Stats.Read = stats.Read
	s.Stats.NumProcs = platform.NumProcs()

	if stats.HCSStats != nil ***REMOVED***
		hcss := stats.HCSStats
		// Populate the CPU/processor statistics
		s.CPUStats = types.CPUStats***REMOVED***
			CPUUsage: types.CPUUsage***REMOVED***
				TotalUsage:        hcss.Processor.TotalRuntime100ns,
				UsageInKernelmode: hcss.Processor.RuntimeKernel100ns,
				UsageInUsermode:   hcss.Processor.RuntimeKernel100ns,
			***REMOVED***,
		***REMOVED***

		// Populate the memory statistics
		s.MemoryStats = types.MemoryStats***REMOVED***
			Commit:            hcss.Memory.UsageCommitBytes,
			CommitPeak:        hcss.Memory.UsageCommitPeakBytes,
			PrivateWorkingSet: hcss.Memory.UsagePrivateWorkingSetBytes,
		***REMOVED***

		// Populate the storage statistics
		s.StorageStats = types.StorageStats***REMOVED***
			ReadCountNormalized:  hcss.Storage.ReadCountNormalized,
			ReadSizeBytes:        hcss.Storage.ReadSizeBytes,
			WriteCountNormalized: hcss.Storage.WriteCountNormalized,
			WriteSizeBytes:       hcss.Storage.WriteSizeBytes,
		***REMOVED***

		// Populate the network statistics
		s.Networks = make(map[string]types.NetworkStats)
		for _, nstats := range hcss.Network ***REMOVED***
			s.Networks[nstats.EndpointId] = types.NetworkStats***REMOVED***
				RxBytes:   nstats.BytesReceived,
				RxPackets: nstats.PacketsReceived,
				RxDropped: nstats.DroppedPacketsIncoming,
				TxBytes:   nstats.BytesSent,
				TxPackets: nstats.PacketsSent,
				TxDropped: nstats.DroppedPacketsOutgoing,
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return s, nil
***REMOVED***

// setDefaultIsolation determine the default isolation mode for the
// daemon to run in. This is only applicable on Windows
func (daemon *Daemon) setDefaultIsolation() error ***REMOVED***
	daemon.defaultIsolation = containertypes.Isolation("process")
	// On client SKUs, default to Hyper-V. Note that IoT reports as a client SKU
	// but it should not be treated as such.
	if system.IsWindowsClient() && !system.IsIoTCore() ***REMOVED***
		daemon.defaultIsolation = containertypes.Isolation("hyperv")
	***REMOVED***
	for _, option := range daemon.configStore.ExecOptions ***REMOVED***
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		key = strings.ToLower(key)
		switch key ***REMOVED***

		case "isolation":
			if !containertypes.Isolation(val).IsValid() ***REMOVED***
				return fmt.Errorf("Invalid exec-opt value for 'isolation':'%s'", val)
			***REMOVED***
			if containertypes.Isolation(val).IsHyperV() ***REMOVED***
				daemon.defaultIsolation = containertypes.Isolation("hyperv")
			***REMOVED***
			if containertypes.Isolation(val).IsProcess() ***REMOVED***
				if system.IsWindowsClient() && !system.IsIoTCore() ***REMOVED***
					// @engine maintainers. This block should not be removed. It partially enforces licensing
					// restrictions on Windows. Ping @jhowardmsft if there are concerns or PRs to change this.
					return fmt.Errorf("Windows client operating systems only support Hyper-V containers")
				***REMOVED***
				daemon.defaultIsolation = containertypes.Isolation("process")
			***REMOVED***
		default:
			return fmt.Errorf("Unrecognised exec-opt '%s'\n", key)
		***REMOVED***
	***REMOVED***

	logrus.Infof("Windows default isolation mode: %s", daemon.defaultIsolation)
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

func setupDaemonProcess(config *config.Config) error ***REMOVED***
	return nil
***REMOVED***

// verifyVolumesInfo is a no-op on windows.
// This is called during daemon initialization to migrate volumes from pre-1.7.
// volumes were not supported on windows pre-1.7
func (daemon *Daemon) verifyVolumesInfo(container *container.Container) error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) setupSeccompProfile() error ***REMOVED***
	return nil
***REMOVED***

func getRealPath(path string) (string, error) ***REMOVED***
	if system.IsIoTCore() ***REMOVED***
		// Due to https://github.com/golang/go/issues/20506, path expansion
		// does not work correctly on the default IoT Core configuration.
		// TODO @darrenstahlmsft remove this once golang/go/20506 is fixed
		return path, nil
	***REMOVED***
	return fileutils.ReadSymlinkedDirectory(path)
***REMOVED***

func (daemon *Daemon) loadRuntimes() error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) initRuntimes(_ map[string]types.Runtime) error ***REMOVED***
	return nil
***REMOVED***
