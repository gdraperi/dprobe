package daemon

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/api/types/versions/v1p20"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/network"
	"github.com/docker/docker/errdefs"
	volumestore "github.com/docker/docker/volume/store"
	"github.com/docker/go-connections/nat"
)

// ContainerInspect returns low-level information about a
// container. Returns an error if the container cannot be found, or if
// there is an error getting the data.
func (daemon *Daemon) ContainerInspect(name string, size bool, version string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch ***REMOVED***
	case versions.LessThan(version, "1.20"):
		return daemon.containerInspectPre120(name)
	case versions.Equal(version, "1.20"):
		return daemon.containerInspect120(name)
	***REMOVED***
	return daemon.ContainerInspectCurrent(name, size)
***REMOVED***

// ContainerInspectCurrent returns low-level information about a
// container in a most recent api version.
func (daemon *Daemon) ContainerInspectCurrent(name string, size bool) (*types.ContainerJSON, error) ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	container.Lock()

	base, err := daemon.getInspectData(container)
	if err != nil ***REMOVED***
		container.Unlock()
		return nil, err
	***REMOVED***

	apiNetworks := make(map[string]*networktypes.EndpointSettings)
	for name, epConf := range container.NetworkSettings.Networks ***REMOVED***
		if epConf.EndpointSettings != nil ***REMOVED***
			// We must make a copy of this pointer object otherwise it can race with other operations
			apiNetworks[name] = epConf.EndpointSettings.Copy()
		***REMOVED***
	***REMOVED***

	mountPoints := container.GetMountPoints()
	networkSettings := &types.NetworkSettings***REMOVED***
		NetworkSettingsBase: types.NetworkSettingsBase***REMOVED***
			Bridge:                 container.NetworkSettings.Bridge,
			SandboxID:              container.NetworkSettings.SandboxID,
			HairpinMode:            container.NetworkSettings.HairpinMode,
			LinkLocalIPv6Address:   container.NetworkSettings.LinkLocalIPv6Address,
			LinkLocalIPv6PrefixLen: container.NetworkSettings.LinkLocalIPv6PrefixLen,
			SandboxKey:             container.NetworkSettings.SandboxKey,
			SecondaryIPAddresses:   container.NetworkSettings.SecondaryIPAddresses,
			SecondaryIPv6Addresses: container.NetworkSettings.SecondaryIPv6Addresses,
		***REMOVED***,
		DefaultNetworkSettings: daemon.getDefaultNetworkSettings(container.NetworkSettings.Networks),
		Networks:               apiNetworks,
	***REMOVED***

	ports := make(nat.PortMap, len(container.NetworkSettings.Ports))
	for k, pm := range container.NetworkSettings.Ports ***REMOVED***
		ports[k] = pm
	***REMOVED***
	networkSettings.NetworkSettingsBase.Ports = ports

	container.Unlock()

	if size ***REMOVED***
		sizeRw, sizeRootFs := daemon.getSize(base.ID)
		base.SizeRw = &sizeRw
		base.SizeRootFs = &sizeRootFs
	***REMOVED***

	return &types.ContainerJSON***REMOVED***
		ContainerJSONBase: base,
		Mounts:            mountPoints,
		Config:            container.Config,
		NetworkSettings:   networkSettings,
	***REMOVED***, nil
***REMOVED***

// containerInspect120 serializes the master version of a container into a json type.
func (daemon *Daemon) containerInspect120(name string) (*v1p20.ContainerJSON, error) ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	container.Lock()
	defer container.Unlock()

	base, err := daemon.getInspectData(container)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	mountPoints := container.GetMountPoints()
	config := &v1p20.ContainerConfig***REMOVED***
		Config:          container.Config,
		MacAddress:      container.Config.MacAddress,
		NetworkDisabled: container.Config.NetworkDisabled,
		ExposedPorts:    container.Config.ExposedPorts,
		VolumeDriver:    container.HostConfig.VolumeDriver,
	***REMOVED***
	networkSettings := daemon.getBackwardsCompatibleNetworkSettings(container.NetworkSettings)

	return &v1p20.ContainerJSON***REMOVED***
		ContainerJSONBase: base,
		Mounts:            mountPoints,
		Config:            config,
		NetworkSettings:   networkSettings,
	***REMOVED***, nil
***REMOVED***

func (daemon *Daemon) getInspectData(container *container.Container) (*types.ContainerJSONBase, error) ***REMOVED***
	// make a copy to play with
	hostConfig := *container.HostConfig

	children := daemon.children(container)
	hostConfig.Links = nil // do not expose the internal structure
	for linkAlias, child := range children ***REMOVED***
		hostConfig.Links = append(hostConfig.Links, fmt.Sprintf("%s:%s", child.Name, linkAlias))
	***REMOVED***

	// We merge the Ulimits from hostConfig with daemon default
	daemon.mergeUlimits(&hostConfig)

	var containerHealth *types.Health
	if container.State.Health != nil ***REMOVED***
		containerHealth = &types.Health***REMOVED***
			Status:        container.State.Health.Status(),
			FailingStreak: container.State.Health.FailingStreak,
			Log:           append([]*types.HealthcheckResult***REMOVED******REMOVED***, container.State.Health.Log...),
		***REMOVED***
	***REMOVED***

	containerState := &types.ContainerState***REMOVED***
		Status:     container.State.StateString(),
		Running:    container.State.Running,
		Paused:     container.State.Paused,
		Restarting: container.State.Restarting,
		OOMKilled:  container.State.OOMKilled,
		Dead:       container.State.Dead,
		Pid:        container.State.Pid,
		ExitCode:   container.State.ExitCode(),
		Error:      container.State.ErrorMsg,
		StartedAt:  container.State.StartedAt.Format(time.RFC3339Nano),
		FinishedAt: container.State.FinishedAt.Format(time.RFC3339Nano),
		Health:     containerHealth,
	***REMOVED***

	contJSONBase := &types.ContainerJSONBase***REMOVED***
		ID:           container.ID,
		Created:      container.Created.Format(time.RFC3339Nano),
		Path:         container.Path,
		Args:         container.Args,
		State:        containerState,
		Image:        container.ImageID.String(),
		LogPath:      container.LogPath,
		Name:         container.Name,
		RestartCount: container.RestartCount,
		Driver:       container.Driver,
		Platform:     container.OS,
		MountLabel:   container.MountLabel,
		ProcessLabel: container.ProcessLabel,
		ExecIDs:      container.GetExecIDs(),
		HostConfig:   &hostConfig,
	***REMOVED***

	// Now set any platform-specific fields
	contJSONBase = setPlatformSpecificContainerFields(container, contJSONBase)

	contJSONBase.GraphDriver.Name = container.Driver

	graphDriverData, err := container.RWLayer.Metadata()
	// If container is marked as Dead, the container's graphdriver metadata
	// could have been removed, it will cause error if we try to get the metadata,
	// we can ignore the error if the container is dead.
	if err != nil && !container.Dead ***REMOVED***
		return nil, errdefs.System(err)
	***REMOVED***
	contJSONBase.GraphDriver.Data = graphDriverData

	return contJSONBase, nil
***REMOVED***

// ContainerExecInspect returns low-level information about the exec
// command. An error is returned if the exec cannot be found.
func (daemon *Daemon) ContainerExecInspect(id string) (*backend.ExecInspect, error) ***REMOVED***
	e := daemon.execCommands.Get(id)
	if e == nil ***REMOVED***
		return nil, errExecNotFound(id)
	***REMOVED***

	if container := daemon.containers.Get(e.ContainerID); container == nil ***REMOVED***
		return nil, errExecNotFound(id)
	***REMOVED***

	pc := inspectExecProcessConfig(e)

	return &backend.ExecInspect***REMOVED***
		ID:            e.ID,
		Running:       e.Running,
		ExitCode:      e.ExitCode,
		ProcessConfig: pc,
		OpenStdin:     e.OpenStdin,
		OpenStdout:    e.OpenStdout,
		OpenStderr:    e.OpenStderr,
		CanRemove:     e.CanRemove,
		ContainerID:   e.ContainerID,
		DetachKeys:    e.DetachKeys,
		Pid:           e.Pid,
	***REMOVED***, nil
***REMOVED***

// VolumeInspect looks up a volume by name. An error is returned if
// the volume cannot be found.
func (daemon *Daemon) VolumeInspect(name string) (*types.Volume, error) ***REMOVED***
	v, err := daemon.volumes.Get(name)
	if err != nil ***REMOVED***
		if volumestore.IsNotExist(err) ***REMOVED***
			return nil, volumeNotFound(name)
		***REMOVED***
		return nil, errdefs.System(err)
	***REMOVED***
	apiV := volumeToAPIType(v)
	apiV.Mountpoint = v.Path()
	apiV.Status = v.Status()
	return apiV, nil
***REMOVED***

func (daemon *Daemon) getBackwardsCompatibleNetworkSettings(settings *network.Settings) *v1p20.NetworkSettings ***REMOVED***
	result := &v1p20.NetworkSettings***REMOVED***
		NetworkSettingsBase: types.NetworkSettingsBase***REMOVED***
			Bridge:                 settings.Bridge,
			SandboxID:              settings.SandboxID,
			HairpinMode:            settings.HairpinMode,
			LinkLocalIPv6Address:   settings.LinkLocalIPv6Address,
			LinkLocalIPv6PrefixLen: settings.LinkLocalIPv6PrefixLen,
			Ports:                  settings.Ports,
			SandboxKey:             settings.SandboxKey,
			SecondaryIPAddresses:   settings.SecondaryIPAddresses,
			SecondaryIPv6Addresses: settings.SecondaryIPv6Addresses,
		***REMOVED***,
		DefaultNetworkSettings: daemon.getDefaultNetworkSettings(settings.Networks),
	***REMOVED***

	return result
***REMOVED***

// getDefaultNetworkSettings creates the deprecated structure that holds the information
// about the bridge network for a container.
func (daemon *Daemon) getDefaultNetworkSettings(networks map[string]*network.EndpointSettings) types.DefaultNetworkSettings ***REMOVED***
	var settings types.DefaultNetworkSettings

	if defaultNetwork, ok := networks["bridge"]; ok && defaultNetwork.EndpointSettings != nil ***REMOVED***
		settings.EndpointID = defaultNetwork.EndpointID
		settings.Gateway = defaultNetwork.Gateway
		settings.GlobalIPv6Address = defaultNetwork.GlobalIPv6Address
		settings.GlobalIPv6PrefixLen = defaultNetwork.GlobalIPv6PrefixLen
		settings.IPAddress = defaultNetwork.IPAddress
		settings.IPPrefixLen = defaultNetwork.IPPrefixLen
		settings.IPv6Gateway = defaultNetwork.IPv6Gateway
		settings.MacAddress = defaultNetwork.MacAddress
	***REMOVED***
	return settings
***REMOVED***
