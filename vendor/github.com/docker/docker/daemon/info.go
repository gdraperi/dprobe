package daemon

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/debug"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/parsers/kernel"
	"github.com/docker/docker/pkg/parsers/operatingsystem"
	"github.com/docker/docker/pkg/platform"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/registry"
	"github.com/docker/docker/volume/drivers"
	"github.com/docker/go-connections/sockets"
	"github.com/sirupsen/logrus"
)

// SystemInfo returns information about the host server the daemon is running on.
func (daemon *Daemon) SystemInfo() (*types.Info, error) ***REMOVED***
	kernelVersion := "<unknown>"
	if kv, err := kernel.GetKernelVersion(); err != nil ***REMOVED***
		logrus.Warnf("Could not get kernel version: %v", err)
	***REMOVED*** else ***REMOVED***
		kernelVersion = kv.String()
	***REMOVED***

	operatingSystem := "<unknown>"
	if s, err := operatingsystem.GetOperatingSystem(); err != nil ***REMOVED***
		logrus.Warnf("Could not get operating system name: %v", err)
	***REMOVED*** else ***REMOVED***
		operatingSystem = s
	***REMOVED***

	// Don't do containerized check on Windows
	if runtime.GOOS != "windows" ***REMOVED***
		if inContainer, err := operatingsystem.IsContainerized(); err != nil ***REMOVED***
			logrus.Errorf("Could not determine if daemon is containerized: %v", err)
			operatingSystem += " (error determining if containerized)"
		***REMOVED*** else if inContainer ***REMOVED***
			operatingSystem += " (containerized)"
		***REMOVED***
	***REMOVED***

	meminfo, err := system.ReadMemInfo()
	if err != nil ***REMOVED***
		logrus.Errorf("Could not read system memory info: %v", err)
		meminfo = &system.MemInfo***REMOVED******REMOVED***
	***REMOVED***

	sysInfo := sysinfo.New(true)
	cRunning, cPaused, cStopped := stateCtr.get()

	securityOptions := []string***REMOVED******REMOVED***
	if sysInfo.AppArmor ***REMOVED***
		securityOptions = append(securityOptions, "name=apparmor")
	***REMOVED***
	if sysInfo.Seccomp && supportsSeccomp ***REMOVED***
		profile := daemon.seccompProfilePath
		if profile == "" ***REMOVED***
			profile = "default"
		***REMOVED***
		securityOptions = append(securityOptions, fmt.Sprintf("name=seccomp,profile=%s", profile))
	***REMOVED***
	if selinuxEnabled() ***REMOVED***
		securityOptions = append(securityOptions, "name=selinux")
	***REMOVED***
	rootIDs := daemon.idMappings.RootPair()
	if rootIDs.UID != 0 || rootIDs.GID != 0 ***REMOVED***
		securityOptions = append(securityOptions, "name=userns")
	***REMOVED***

	var ds [][2]string
	drivers := ""
	for os, gd := range daemon.graphDrivers ***REMOVED***
		ds = append(ds, daemon.layerStores[os].DriverStatus()...)
		drivers += gd
		if len(daemon.graphDrivers) > 1 ***REMOVED***
			drivers += fmt.Sprintf(" (%s) ", os)
		***REMOVED***
	***REMOVED***
	drivers = strings.TrimSpace(drivers)

	v := &types.Info***REMOVED***
		ID:                 daemon.ID,
		Containers:         cRunning + cPaused + cStopped,
		ContainersRunning:  cRunning,
		ContainersPaused:   cPaused,
		ContainersStopped:  cStopped,
		Images:             len(daemon.imageStore.Map()),
		Driver:             drivers,
		DriverStatus:       ds,
		Plugins:            daemon.showPluginsInfo(),
		IPv4Forwarding:     !sysInfo.IPv4ForwardingDisabled,
		BridgeNfIptables:   !sysInfo.BridgeNFCallIPTablesDisabled,
		BridgeNfIP6tables:  !sysInfo.BridgeNFCallIP6TablesDisabled,
		Debug:              debug.IsEnabled(),
		NFd:                fileutils.GetTotalUsedFds(),
		NGoroutines:        runtime.NumGoroutine(),
		SystemTime:         time.Now().Format(time.RFC3339Nano),
		LoggingDriver:      daemon.defaultLogConfig.Type,
		CgroupDriver:       daemon.getCgroupDriver(),
		NEventsListener:    daemon.EventsService.SubscribersCount(),
		KernelVersion:      kernelVersion,
		OperatingSystem:    operatingSystem,
		IndexServerAddress: registry.IndexServer,
		OSType:             platform.OSType,
		Architecture:       platform.Architecture,
		RegistryConfig:     daemon.RegistryService.ServiceConfig(),
		NCPU:               sysinfo.NumCPU(),
		MemTotal:           meminfo.MemTotal,
		GenericResources:   daemon.genericResources,
		DockerRootDir:      daemon.configStore.Root,
		Labels:             daemon.configStore.Labels,
		ExperimentalBuild:  daemon.configStore.Experimental,
		ServerVersion:      dockerversion.Version,
		ClusterStore:       daemon.configStore.ClusterStore,
		ClusterAdvertise:   daemon.configStore.ClusterAdvertise,
		HTTPProxy:          sockets.GetProxyEnv("http_proxy"),
		HTTPSProxy:         sockets.GetProxyEnv("https_proxy"),
		NoProxy:            sockets.GetProxyEnv("no_proxy"),
		LiveRestoreEnabled: daemon.configStore.LiveRestoreEnabled,
		SecurityOptions:    securityOptions,
		Isolation:          daemon.defaultIsolation,
	***REMOVED***

	// Retrieve platform specific info
	daemon.FillPlatformInfo(v, sysInfo)

	hostname := ""
	if hn, err := os.Hostname(); err != nil ***REMOVED***
		logrus.Warnf("Could not get hostname: %v", err)
	***REMOVED*** else ***REMOVED***
		hostname = hn
	***REMOVED***
	v.Name = hostname

	return v, nil
***REMOVED***

// SystemVersion returns version information about the daemon.
func (daemon *Daemon) SystemVersion() types.Version ***REMOVED***
	kernelVersion := "<unknown>"
	if kv, err := kernel.GetKernelVersion(); err != nil ***REMOVED***
		logrus.Warnf("Could not get kernel version: %v", err)
	***REMOVED*** else ***REMOVED***
		kernelVersion = kv.String()
	***REMOVED***

	v := types.Version***REMOVED***
		Components: []types.ComponentVersion***REMOVED***
			***REMOVED***
				Name:    "Engine",
				Version: dockerversion.Version,
				Details: map[string]string***REMOVED***
					"GitCommit":     dockerversion.GitCommit,
					"ApiVersion":    api.DefaultVersion,
					"MinAPIVersion": api.MinVersion,
					"GoVersion":     runtime.Version(),
					"Os":            runtime.GOOS,
					"Arch":          runtime.GOARCH,
					"BuildTime":     dockerversion.BuildTime,
					"KernelVersion": kernelVersion,
					"Experimental":  fmt.Sprintf("%t", daemon.configStore.Experimental),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		// Populate deprecated fields for older clients
		Version:       dockerversion.Version,
		GitCommit:     dockerversion.GitCommit,
		APIVersion:    api.DefaultVersion,
		MinAPIVersion: api.MinVersion,
		GoVersion:     runtime.Version(),
		Os:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		BuildTime:     dockerversion.BuildTime,
		KernelVersion: kernelVersion,
		Experimental:  daemon.configStore.Experimental,
	***REMOVED***

	v.Platform.Name = dockerversion.PlatformName

	return v
***REMOVED***

func (daemon *Daemon) showPluginsInfo() types.PluginsInfo ***REMOVED***
	var pluginsInfo types.PluginsInfo

	pluginsInfo.Volume = volumedrivers.GetDriverList()
	pluginsInfo.Network = daemon.GetNetworkDriverList()
	// The authorization plugins are returned in the order they are
	// used as they constitute a request/response modification chain.
	pluginsInfo.Authorization = daemon.configStore.AuthorizationPlugins
	pluginsInfo.Log = logger.ListDrivers()

	return pluginsInfo
***REMOVED***
