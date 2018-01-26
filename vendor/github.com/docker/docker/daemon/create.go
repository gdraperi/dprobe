package daemon

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/runconfig"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/sirupsen/logrus"
)

// CreateManagedContainer creates a container that is managed by a Service
func (daemon *Daemon) CreateManagedContainer(params types.ContainerCreateConfig) (containertypes.ContainerCreateCreatedBody, error) ***REMOVED***
	return daemon.containerCreate(params, true)
***REMOVED***

// ContainerCreate creates a regular container
func (daemon *Daemon) ContainerCreate(params types.ContainerCreateConfig) (containertypes.ContainerCreateCreatedBody, error) ***REMOVED***
	return daemon.containerCreate(params, false)
***REMOVED***

func (daemon *Daemon) containerCreate(params types.ContainerCreateConfig, managed bool) (containertypes.ContainerCreateCreatedBody, error) ***REMOVED***
	start := time.Now()
	if params.Config == nil ***REMOVED***
		return containertypes.ContainerCreateCreatedBody***REMOVED******REMOVED***, errdefs.InvalidParameter(errors.New("Config cannot be empty in order to create a container"))
	***REMOVED***

	os := runtime.GOOS
	if params.Config.Image != "" ***REMOVED***
		img, err := daemon.GetImage(params.Config.Image)
		if err == nil ***REMOVED***
			os = img.OS
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// This mean scratch. On Windows, we can safely assume that this is a linux
		// container. On other platforms, it's the host OS (which it already is)
		if runtime.GOOS == "windows" && system.LCOWSupported() ***REMOVED***
			os = "linux"
		***REMOVED***
	***REMOVED***

	warnings, err := daemon.verifyContainerSettings(os, params.HostConfig, params.Config, false)
	if err != nil ***REMOVED***
		return containertypes.ContainerCreateCreatedBody***REMOVED***Warnings: warnings***REMOVED***, errdefs.InvalidParameter(err)
	***REMOVED***

	err = verifyNetworkingConfig(params.NetworkingConfig)
	if err != nil ***REMOVED***
		return containertypes.ContainerCreateCreatedBody***REMOVED***Warnings: warnings***REMOVED***, errdefs.InvalidParameter(err)
	***REMOVED***

	if params.HostConfig == nil ***REMOVED***
		params.HostConfig = &containertypes.HostConfig***REMOVED******REMOVED***
	***REMOVED***
	err = daemon.adaptContainerSettings(params.HostConfig, params.AdjustCPUShares)
	if err != nil ***REMOVED***
		return containertypes.ContainerCreateCreatedBody***REMOVED***Warnings: warnings***REMOVED***, errdefs.InvalidParameter(err)
	***REMOVED***

	container, err := daemon.create(params, managed)
	if err != nil ***REMOVED***
		return containertypes.ContainerCreateCreatedBody***REMOVED***Warnings: warnings***REMOVED***, err
	***REMOVED***
	containerActions.WithValues("create").UpdateSince(start)

	return containertypes.ContainerCreateCreatedBody***REMOVED***ID: container.ID, Warnings: warnings***REMOVED***, nil
***REMOVED***

// Create creates a new container from the given configuration with a given name.
func (daemon *Daemon) create(params types.ContainerCreateConfig, managed bool) (retC *container.Container, retErr error) ***REMOVED***
	var (
		container *container.Container
		img       *image.Image
		imgID     image.ID
		err       error
	)

	os := runtime.GOOS
	if params.Config.Image != "" ***REMOVED***
		img, err = daemon.GetImage(params.Config.Image)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if img.OS != "" ***REMOVED***
			os = img.OS
		***REMOVED*** else ***REMOVED***
			// default to the host OS except on Windows with LCOW
			if runtime.GOOS == "windows" && system.LCOWSupported() ***REMOVED***
				os = "linux"
			***REMOVED***
		***REMOVED***
		imgID = img.ID()

		if runtime.GOOS == "windows" && img.OS == "linux" && !system.LCOWSupported() ***REMOVED***
			return nil, errors.New("operating system on which parent image was created is not Windows")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if runtime.GOOS == "windows" ***REMOVED***
			os = "linux" // 'scratch' case.
		***REMOVED***
	***REMOVED***

	if err := daemon.mergeAndVerifyConfig(params.Config, img); err != nil ***REMOVED***
		return nil, errdefs.InvalidParameter(err)
	***REMOVED***

	if err := daemon.mergeAndVerifyLogConfig(&params.HostConfig.LogConfig); err != nil ***REMOVED***
		return nil, errdefs.InvalidParameter(err)
	***REMOVED***

	if container, err = daemon.newContainer(params.Name, os, params.Config, params.HostConfig, imgID, managed); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if retErr != nil ***REMOVED***
			if err := daemon.cleanupContainer(container, true, true); err != nil ***REMOVED***
				logrus.Errorf("failed to cleanup container on create error: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err := daemon.setSecurityOptions(container, params.HostConfig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	container.HostConfig.StorageOpt = params.HostConfig.StorageOpt

	// Fixes: https://github.com/moby/moby/issues/34074 and
	// https://github.com/docker/for-win/issues/999.
	// Merge the daemon's storage options if they aren't already present. We only
	// do this on Windows as there's no effective sandbox size limit other than
	// physical on Linux.
	if runtime.GOOS == "windows" ***REMOVED***
		if container.HostConfig.StorageOpt == nil ***REMOVED***
			container.HostConfig.StorageOpt = make(map[string]string)
		***REMOVED***
		for _, v := range daemon.configStore.GraphOptions ***REMOVED***
			opt := strings.SplitN(v, "=", 2)
			if _, ok := container.HostConfig.StorageOpt[opt[0]]; !ok ***REMOVED***
				container.HostConfig.StorageOpt[opt[0]] = opt[1]
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Set RWLayer for container after mount labels have been set
	if err := daemon.setRWLayer(container); err != nil ***REMOVED***
		return nil, errdefs.System(err)
	***REMOVED***

	rootIDs := daemon.idMappings.RootPair()
	if err := idtools.MkdirAndChown(container.Root, 0700, rootIDs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := idtools.MkdirAndChown(container.CheckpointDir(), 0700, rootIDs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := daemon.setHostConfig(container, params.HostConfig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := daemon.createContainerOSSpecificSettings(container, params.Config, params.HostConfig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var endpointsConfigs map[string]*networktypes.EndpointSettings
	if params.NetworkingConfig != nil ***REMOVED***
		endpointsConfigs = params.NetworkingConfig.EndpointsConfig
	***REMOVED***
	// Make sure NetworkMode has an acceptable value. We do this to ensure
	// backwards API compatibility.
	runconfig.SetDefaultNetModeIfBlank(container.HostConfig)

	daemon.updateContainerNetworkSettings(container, endpointsConfigs)
	if err := daemon.Register(container); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	stateCtr.set(container.ID, "stopped")
	daemon.LogContainerEvent(container, "create")
	return container, nil
***REMOVED***

func toHostConfigSelinuxLabels(labels []string) []string ***REMOVED***
	for i, l := range labels ***REMOVED***
		labels[i] = "label=" + l
	***REMOVED***
	return labels
***REMOVED***

func (daemon *Daemon) generateSecurityOpt(hostConfig *containertypes.HostConfig) ([]string, error) ***REMOVED***
	for _, opt := range hostConfig.SecurityOpt ***REMOVED***
		con := strings.Split(opt, "=")
		if con[0] == "label" ***REMOVED***
			// Caller overrode SecurityOpts
			return nil, nil
		***REMOVED***
	***REMOVED***
	ipcMode := hostConfig.IpcMode
	pidMode := hostConfig.PidMode
	privileged := hostConfig.Privileged
	if ipcMode.IsHost() || pidMode.IsHost() || privileged ***REMOVED***
		return toHostConfigSelinuxLabels(label.DisableSecOpt()), nil
	***REMOVED***

	var ipcLabel []string
	var pidLabel []string
	ipcContainer := ipcMode.Container()
	pidContainer := pidMode.Container()
	if ipcContainer != "" ***REMOVED***
		c, err := daemon.GetContainer(ipcContainer)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ipcLabel = label.DupSecOpt(c.ProcessLabel)
		if pidContainer == "" ***REMOVED***
			return toHostConfigSelinuxLabels(ipcLabel), err
		***REMOVED***
	***REMOVED***
	if pidContainer != "" ***REMOVED***
		c, err := daemon.GetContainer(pidContainer)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		pidLabel = label.DupSecOpt(c.ProcessLabel)
		if ipcContainer == "" ***REMOVED***
			return toHostConfigSelinuxLabels(pidLabel), err
		***REMOVED***
	***REMOVED***

	if pidLabel != nil && ipcLabel != nil ***REMOVED***
		for i := 0; i < len(pidLabel); i++ ***REMOVED***
			if pidLabel[i] != ipcLabel[i] ***REMOVED***
				return nil, fmt.Errorf("--ipc and --pid containers SELinux labels aren't the same")
			***REMOVED***
		***REMOVED***
		return toHostConfigSelinuxLabels(pidLabel), nil
	***REMOVED***
	return nil, nil
***REMOVED***

func (daemon *Daemon) setRWLayer(container *container.Container) error ***REMOVED***
	var layerID layer.ChainID
	if container.ImageID != "" ***REMOVED***
		img, err := daemon.imageStore.Get(container.ImageID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		layerID = img.RootFS.ChainID()
	***REMOVED***

	rwLayerOpts := &layer.CreateRWLayerOpts***REMOVED***
		MountLabel: container.MountLabel,
		InitFunc:   daemon.getLayerInit(),
		StorageOpt: container.HostConfig.StorageOpt,
	***REMOVED***

	// Indexing by OS is safe here as validation of OS has already been performed in create() (the only
	// caller), and guaranteed non-nil
	rwLayer, err := daemon.layerStores[container.OS].CreateRWLayer(container.ID, layerID, rwLayerOpts)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	container.RWLayer = rwLayer

	return nil
***REMOVED***

// VolumeCreate creates a volume with the specified name, driver, and opts
// This is called directly from the Engine API
func (daemon *Daemon) VolumeCreate(name, driverName string, opts, labels map[string]string) (*types.Volume, error) ***REMOVED***
	if name == "" ***REMOVED***
		name = stringid.GenerateNonCryptoID()
	***REMOVED***

	v, err := daemon.volumes.Create(name, driverName, opts, labels)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	daemon.LogVolumeEvent(v.Name(), "create", map[string]string***REMOVED***"driver": v.DriverName()***REMOVED***)
	apiV := volumeToAPIType(v)
	apiV.Mountpoint = v.Path()
	return apiV, nil
***REMOVED***

func (daemon *Daemon) mergeAndVerifyConfig(config *containertypes.Config, img *image.Image) error ***REMOVED***
	if img != nil && img.Config != nil ***REMOVED***
		if err := merge(config, img.Config); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// Reset the Entrypoint if it is [""]
	if len(config.Entrypoint) == 1 && config.Entrypoint[0] == "" ***REMOVED***
		config.Entrypoint = nil
	***REMOVED***
	if len(config.Entrypoint) == 0 && len(config.Cmd) == 0 ***REMOVED***
		return fmt.Errorf("No command specified")
	***REMOVED***
	return nil
***REMOVED***

// Checks if the client set configurations for more than one network while creating a container
// Also checks if the IPAMConfig is valid
func verifyNetworkingConfig(nwConfig *networktypes.NetworkingConfig) error ***REMOVED***
	if nwConfig == nil || len(nwConfig.EndpointsConfig) == 0 ***REMOVED***
		return nil
	***REMOVED***
	if len(nwConfig.EndpointsConfig) == 1 ***REMOVED***
		for k, v := range nwConfig.EndpointsConfig ***REMOVED***
			if v == nil ***REMOVED***
				return errdefs.InvalidParameter(errors.Errorf("no EndpointSettings for %s", k))
			***REMOVED***
			if v.IPAMConfig != nil ***REMOVED***
				if v.IPAMConfig.IPv4Address != "" && net.ParseIP(v.IPAMConfig.IPv4Address).To4() == nil ***REMOVED***
					return errors.Errorf("invalid IPv4 address: %s", v.IPAMConfig.IPv4Address)
				***REMOVED***
				if v.IPAMConfig.IPv6Address != "" ***REMOVED***
					n := net.ParseIP(v.IPAMConfig.IPv6Address)
					// if the address is an invalid network address (ParseIP == nil) or if it is
					// an IPv4 address (To4() != nil), then it is an invalid IPv6 address
					if n == nil || n.To4() != nil ***REMOVED***
						return errors.Errorf("invalid IPv6 address: %s", v.IPAMConfig.IPv6Address)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	l := make([]string, 0, len(nwConfig.EndpointsConfig))
	for k := range nwConfig.EndpointsConfig ***REMOVED***
		l = append(l, k)
	***REMOVED***
	return errors.Errorf("Container cannot be connected to network endpoints: %s", strings.Join(l, ", "))
***REMOVED***
