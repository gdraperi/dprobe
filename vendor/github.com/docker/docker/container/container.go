package container

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/containerd/containerd/cio"
	containertypes "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	networktypes "github.com/docker/docker/api/types/network"
	swarmtypes "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/container/stream"
	"github.com/docker/docker/daemon/exec"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/jsonfilelog"
	"github.com/docker/docker/daemon/network"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/symlink"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/restartmanager"
	"github.com/docker/docker/runconfig"
	"github.com/docker/docker/volume"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/docker/libnetwork"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/options"
	"github.com/docker/libnetwork/types"
	agentexec "github.com/docker/swarmkit/agent/exec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const configFileName = "config.v2.json"

var (
	errInvalidEndpoint = errors.New("invalid endpoint while building port map info")
	errInvalidNetwork  = errors.New("invalid network settings while building port map info")
)

// ExitStatus provides exit reasons for a container.
type ExitStatus struct ***REMOVED***
	// The exit code with which the container exited.
	ExitCode int

	// Whether the container encountered an OOM.
	OOMKilled bool

	// Time at which the container died
	ExitedAt time.Time
***REMOVED***

// Container holds the structure defining a container object.
type Container struct ***REMOVED***
	StreamConfig *stream.Config
	// embed for Container to support states directly.
	*State          `json:"State"`          // Needed for Engine API version <= 1.11
	Root            string                  `json:"-"` // Path to the "home" of the container, including metadata.
	BaseFS          containerfs.ContainerFS `json:"-"` // interface containing graphdriver mount
	RWLayer         layer.RWLayer           `json:"-"`
	ID              string
	Created         time.Time
	Managed         bool
	Path            string
	Args            []string
	Config          *containertypes.Config
	ImageID         image.ID `json:"Image"`
	NetworkSettings *network.Settings
	LogPath         string
	Name            string
	Driver          string
	OS              string
	// MountLabel contains the options for the 'mount' command
	MountLabel             string
	ProcessLabel           string
	RestartCount           int
	HasBeenStartedBefore   bool
	HasBeenManuallyStopped bool // used for unless-stopped restart policy
	MountPoints            map[string]*volume.MountPoint
	HostConfig             *containertypes.HostConfig `json:"-"` // do not serialize the host config in the json, otherwise we'll make the container unportable
	ExecCommands           *exec.Store                `json:"-"`
	DependencyStore        agentexec.DependencyGetter `json:"-"`
	SecretReferences       []*swarmtypes.SecretReference
	ConfigReferences       []*swarmtypes.ConfigReference
	// logDriver for closing
	LogDriver      logger.Logger  `json:"-"`
	LogCopier      *logger.Copier `json:"-"`
	restartManager restartmanager.RestartManager
	attachContext  *attachContext

	// Fields here are specific to Unix platforms
	AppArmorProfile string
	HostnamePath    string
	HostsPath       string
	ShmPath         string
	ResolvConfPath  string
	SeccompProfile  string
	NoNewPrivileges bool

	// Fields here are specific to Windows
	NetworkSharedContainerID string   `json:"-"`
	SharedEndpointList       []string `json:"-"`
***REMOVED***

// NewBaseContainer creates a new container with its
// basic configuration.
func NewBaseContainer(id, root string) *Container ***REMOVED***
	return &Container***REMOVED***
		ID:            id,
		State:         NewState(),
		ExecCommands:  exec.NewStore(),
		Root:          root,
		MountPoints:   make(map[string]*volume.MountPoint),
		StreamConfig:  stream.NewConfig(),
		attachContext: &attachContext***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// FromDisk loads the container configuration stored in the host.
func (container *Container) FromDisk() error ***REMOVED***
	pth, err := container.ConfigPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	jsonSource, err := os.Open(pth)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer jsonSource.Close()

	dec := json.NewDecoder(jsonSource)

	// Load container settings
	if err := dec.Decode(container); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Ensure the operating system is set if blank. Assume it is the OS of the
	// host OS if not, to ensure containers created before multiple-OS
	// support are migrated
	if container.OS == "" ***REMOVED***
		container.OS = runtime.GOOS
	***REMOVED***

	return container.readHostConfig()
***REMOVED***

// toDisk saves the container configuration on disk and returns a deep copy.
func (container *Container) toDisk() (*Container, error) ***REMOVED***
	var (
		buf      bytes.Buffer
		deepCopy Container
	)
	pth, err := container.ConfigPath()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Save container settings
	f, err := ioutils.NewAtomicFileWriter(pth, 0600)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	w := io.MultiWriter(&buf, f)
	if err := json.NewEncoder(w).Encode(container); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := json.NewDecoder(&buf).Decode(&deepCopy); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	deepCopy.HostConfig, err = container.WriteHostConfig()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &deepCopy, nil
***REMOVED***

// CheckpointTo makes the Container's current state visible to queries, and persists state.
// Callers must hold a Container lock.
func (container *Container) CheckpointTo(store ViewDB) error ***REMOVED***
	deepCopy, err := container.toDisk()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return store.Save(deepCopy)
***REMOVED***

// readHostConfig reads the host configuration from disk for the container.
func (container *Container) readHostConfig() error ***REMOVED***
	container.HostConfig = &containertypes.HostConfig***REMOVED******REMOVED***
	// If the hostconfig file does not exist, do not read it.
	// (We still have to initialize container.HostConfig,
	// but that's OK, since we just did that above.)
	pth, err := container.HostConfigPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	f, err := os.Open(pth)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&container.HostConfig); err != nil ***REMOVED***
		return err
	***REMOVED***

	container.InitDNSHostConfig()

	return nil
***REMOVED***

// WriteHostConfig saves the host configuration on disk for the container,
// and returns a deep copy of the saved object. Callers must hold a Container lock.
func (container *Container) WriteHostConfig() (*containertypes.HostConfig, error) ***REMOVED***
	var (
		buf      bytes.Buffer
		deepCopy containertypes.HostConfig
	)

	pth, err := container.HostConfigPath()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	f, err := ioutils.NewAtomicFileWriter(pth, 0644)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	w := io.MultiWriter(&buf, f)
	if err := json.NewEncoder(w).Encode(&container.HostConfig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := json.NewDecoder(&buf).Decode(&deepCopy); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &deepCopy, nil
***REMOVED***

// SetupWorkingDirectory sets up the container's working directory as set in container.Config.WorkingDir
func (container *Container) SetupWorkingDirectory(rootIDs idtools.IDPair) error ***REMOVED***
	// TODO @jhowardmsft, @gupta-ak LCOW Support. This will need revisiting.
	// We will need to do remote filesystem operations here.
	if container.OS != runtime.GOOS ***REMOVED***
		return nil
	***REMOVED***

	if container.Config.WorkingDir == "" ***REMOVED***
		return nil
	***REMOVED***

	container.Config.WorkingDir = filepath.Clean(container.Config.WorkingDir)
	pth, err := container.GetResourcePath(container.Config.WorkingDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := idtools.MkdirAllAndChownNew(pth, 0755, rootIDs); err != nil ***REMOVED***
		pthInfo, err2 := os.Stat(pth)
		if err2 == nil && pthInfo != nil && !pthInfo.IsDir() ***REMOVED***
			return errors.Errorf("Cannot mkdir: %s is not a directory", container.Config.WorkingDir)
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

// GetResourcePath evaluates `path` in the scope of the container's BaseFS, with proper path
// sanitisation. Symlinks are all scoped to the BaseFS of the container, as
// though the container's BaseFS was `/`.
//
// The BaseFS of a container is the host-facing path which is bind-mounted as
// `/` inside the container. This method is essentially used to access a
// particular path inside the container as though you were a process in that
// container.
//
// NOTE: The returned path is *only* safely scoped inside the container's BaseFS
//       if no component of the returned path changes (such as a component
//       symlinking to a different path) between using this method and using the
//       path. See symlink.FollowSymlinkInScope for more details.
func (container *Container) GetResourcePath(path string) (string, error) ***REMOVED***
	// IMPORTANT - These are paths on the OS where the daemon is running, hence
	// any filepath operations must be done in an OS agnostic way.
	r, e := container.BaseFS.ResolveScopedPath(path, false)

	// Log this here on the daemon side as there's otherwise no indication apart
	// from the error being propagated all the way back to the client. This makes
	// debugging significantly easier and clearly indicates the error comes from the daemon.
	if e != nil ***REMOVED***
		logrus.Errorf("Failed to ResolveScopedPath BaseFS %s path %s %s\n", container.BaseFS.Path(), path, e)
	***REMOVED***
	return r, e
***REMOVED***

// GetRootResourcePath evaluates `path` in the scope of the container's root, with proper path
// sanitisation. Symlinks are all scoped to the root of the container, as
// though the container's root was `/`.
//
// The root of a container is the host-facing configuration metadata directory.
// Only use this method to safely access the container's `container.json` or
// other metadata files. If in doubt, use container.GetResourcePath.
//
// NOTE: The returned path is *only* safely scoped inside the container's root
//       if no component of the returned path changes (such as a component
//       symlinking to a different path) between using this method and using the
//       path. See symlink.FollowSymlinkInScope for more details.
func (container *Container) GetRootResourcePath(path string) (string, error) ***REMOVED***
	// IMPORTANT - These are paths on the OS where the daemon is running, hence
	// any filepath operations must be done in an OS agnostic way.
	cleanPath := filepath.Join(string(os.PathSeparator), path)
	return symlink.FollowSymlinkInScope(filepath.Join(container.Root, cleanPath), container.Root)
***REMOVED***

// ExitOnNext signals to the monitor that it should not restart the container
// after we send the kill signal.
func (container *Container) ExitOnNext() ***REMOVED***
	container.RestartManager().Cancel()
***REMOVED***

// HostConfigPath returns the path to the container's JSON hostconfig
func (container *Container) HostConfigPath() (string, error) ***REMOVED***
	return container.GetRootResourcePath("hostconfig.json")
***REMOVED***

// ConfigPath returns the path to the container's JSON config
func (container *Container) ConfigPath() (string, error) ***REMOVED***
	return container.GetRootResourcePath(configFileName)
***REMOVED***

// CheckpointDir returns the directory checkpoints are stored in
func (container *Container) CheckpointDir() string ***REMOVED***
	return filepath.Join(container.Root, "checkpoints")
***REMOVED***

// StartLogger starts a new logger driver for the container.
func (container *Container) StartLogger() (logger.Logger, error) ***REMOVED***
	cfg := container.HostConfig.LogConfig
	initDriver, err := logger.GetLogDriver(cfg.Type)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to get logging factory")
	***REMOVED***
	info := logger.Info***REMOVED***
		Config:              cfg.Config,
		ContainerID:         container.ID,
		ContainerName:       container.Name,
		ContainerEntrypoint: container.Path,
		ContainerArgs:       container.Args,
		ContainerImageID:    container.ImageID.String(),
		ContainerImageName:  container.Config.Image,
		ContainerCreated:    container.Created,
		ContainerEnv:        container.Config.Env,
		ContainerLabels:     container.Config.Labels,
		DaemonName:          "docker",
	***REMOVED***

	// Set logging file for "json-logger"
	if cfg.Type == jsonfilelog.Name ***REMOVED***
		info.LogPath, err = container.GetRootResourcePath(fmt.Sprintf("%s-json.log", container.ID))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	l, err := initDriver(info)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if containertypes.LogMode(cfg.Config["mode"]) == containertypes.LogModeNonBlock ***REMOVED***
		bufferSize := int64(-1)
		if s, exists := cfg.Config["max-buffer-size"]; exists ***REMOVED***
			bufferSize, err = units.RAMInBytes(s)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		l = logger.NewRingLogger(l, info, bufferSize)
	***REMOVED***
	return l, nil
***REMOVED***

// GetProcessLabel returns the process label for the container.
func (container *Container) GetProcessLabel() string ***REMOVED***
	// even if we have a process label return "" if we are running
	// in privileged mode
	if container.HostConfig.Privileged ***REMOVED***
		return ""
	***REMOVED***
	return container.ProcessLabel
***REMOVED***

// GetMountLabel returns the mounting label for the container.
// This label is empty if the container is privileged.
func (container *Container) GetMountLabel() string ***REMOVED***
	return container.MountLabel
***REMOVED***

// GetExecIDs returns the list of exec commands running on the container.
func (container *Container) GetExecIDs() []string ***REMOVED***
	return container.ExecCommands.List()
***REMOVED***

// ShouldRestart decides whether the daemon should restart the container or not.
// This is based on the container's restart policy.
func (container *Container) ShouldRestart() bool ***REMOVED***
	shouldRestart, _, _ := container.RestartManager().ShouldRestart(uint32(container.ExitCode()), container.HasBeenManuallyStopped, container.FinishedAt.Sub(container.StartedAt))
	return shouldRestart
***REMOVED***

// AddMountPointWithVolume adds a new mount point configured with a volume to the container.
func (container *Container) AddMountPointWithVolume(destination string, vol volume.Volume, rw bool) ***REMOVED***
	operatingSystem := container.OS
	if operatingSystem == "" ***REMOVED***
		operatingSystem = runtime.GOOS
	***REMOVED***
	volumeParser := volume.NewParser(operatingSystem)
	container.MountPoints[destination] = &volume.MountPoint***REMOVED***
		Type:        mounttypes.TypeVolume,
		Name:        vol.Name(),
		Driver:      vol.DriverName(),
		Destination: destination,
		RW:          rw,
		Volume:      vol,
		CopyData:    volumeParser.DefaultCopyMode(),
	***REMOVED***
***REMOVED***

// UnmountVolumes unmounts all volumes
func (container *Container) UnmountVolumes(volumeEventLog func(name, action string, attributes map[string]string)) error ***REMOVED***
	var errors []string
	for _, volumeMount := range container.MountPoints ***REMOVED***
		if volumeMount.Volume == nil ***REMOVED***
			continue
		***REMOVED***

		if err := volumeMount.Cleanup(); err != nil ***REMOVED***
			errors = append(errors, err.Error())
			continue
		***REMOVED***

		attributes := map[string]string***REMOVED***
			"driver":    volumeMount.Volume.DriverName(),
			"container": container.ID,
		***REMOVED***
		volumeEventLog(volumeMount.Volume.Name(), "unmount", attributes)
	***REMOVED***
	if len(errors) > 0 ***REMOVED***
		return fmt.Errorf("error while unmounting volumes for container %s: %s", container.ID, strings.Join(errors, "; "))
	***REMOVED***
	return nil
***REMOVED***

// IsDestinationMounted checks whether a path is mounted on the container or not.
func (container *Container) IsDestinationMounted(destination string) bool ***REMOVED***
	return container.MountPoints[destination] != nil
***REMOVED***

// StopSignal returns the signal used to stop the container.
func (container *Container) StopSignal() int ***REMOVED***
	var stopSignal syscall.Signal
	if container.Config.StopSignal != "" ***REMOVED***
		stopSignal, _ = signal.ParseSignal(container.Config.StopSignal)
	***REMOVED***

	if int(stopSignal) == 0 ***REMOVED***
		stopSignal, _ = signal.ParseSignal(signal.DefaultStopSignal)
	***REMOVED***
	return int(stopSignal)
***REMOVED***

// StopTimeout returns the timeout (in seconds) used to stop the container.
func (container *Container) StopTimeout() int ***REMOVED***
	if container.Config.StopTimeout != nil ***REMOVED***
		return *container.Config.StopTimeout
	***REMOVED***
	return DefaultStopTimeout
***REMOVED***

// InitDNSHostConfig ensures that the dns fields are never nil.
// New containers don't ever have those fields nil,
// but pre created containers can still have those nil values.
// The non-recommended host configuration in the start api can
// make these fields nil again, this corrects that issue until
// we remove that behavior for good.
// See https://github.com/docker/docker/pull/17779
// for a more detailed explanation on why we don't want that.
func (container *Container) InitDNSHostConfig() ***REMOVED***
	container.Lock()
	defer container.Unlock()
	if container.HostConfig.DNS == nil ***REMOVED***
		container.HostConfig.DNS = make([]string, 0)
	***REMOVED***

	if container.HostConfig.DNSSearch == nil ***REMOVED***
		container.HostConfig.DNSSearch = make([]string, 0)
	***REMOVED***

	if container.HostConfig.DNSOptions == nil ***REMOVED***
		container.HostConfig.DNSOptions = make([]string, 0)
	***REMOVED***
***REMOVED***

// GetEndpointInNetwork returns the container's endpoint to the provided network.
func (container *Container) GetEndpointInNetwork(n libnetwork.Network) (libnetwork.Endpoint, error) ***REMOVED***
	endpointName := strings.TrimPrefix(container.Name, "/")
	return n.EndpointByName(endpointName)
***REMOVED***

func (container *Container) buildPortMapInfo(ep libnetwork.Endpoint) error ***REMOVED***
	if ep == nil ***REMOVED***
		return errInvalidEndpoint
	***REMOVED***

	networkSettings := container.NetworkSettings
	if networkSettings == nil ***REMOVED***
		return errInvalidNetwork
	***REMOVED***

	if len(networkSettings.Ports) == 0 ***REMOVED***
		pm, err := getEndpointPortMapInfo(ep)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		networkSettings.Ports = pm
	***REMOVED***
	return nil
***REMOVED***

func getEndpointPortMapInfo(ep libnetwork.Endpoint) (nat.PortMap, error) ***REMOVED***
	pm := nat.PortMap***REMOVED******REMOVED***
	driverInfo, err := ep.DriverInfo()
	if err != nil ***REMOVED***
		return pm, err
	***REMOVED***

	if driverInfo == nil ***REMOVED***
		// It is not an error for epInfo to be nil
		return pm, nil
	***REMOVED***

	if expData, ok := driverInfo[netlabel.ExposedPorts]; ok ***REMOVED***
		if exposedPorts, ok := expData.([]types.TransportPort); ok ***REMOVED***
			for _, tp := range exposedPorts ***REMOVED***
				natPort, err := nat.NewPort(tp.Proto.String(), strconv.Itoa(int(tp.Port)))
				if err != nil ***REMOVED***
					return pm, fmt.Errorf("Error parsing Port value(%v):%v", tp.Port, err)
				***REMOVED***
				pm[natPort] = nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	mapData, ok := driverInfo[netlabel.PortMap]
	if !ok ***REMOVED***
		return pm, nil
	***REMOVED***

	if portMapping, ok := mapData.([]types.PortBinding); ok ***REMOVED***
		for _, pp := range portMapping ***REMOVED***
			natPort, err := nat.NewPort(pp.Proto.String(), strconv.Itoa(int(pp.Port)))
			if err != nil ***REMOVED***
				return pm, err
			***REMOVED***
			natBndg := nat.PortBinding***REMOVED***HostIP: pp.HostIP.String(), HostPort: strconv.Itoa(int(pp.HostPort))***REMOVED***
			pm[natPort] = append(pm[natPort], natBndg)
		***REMOVED***
	***REMOVED***

	return pm, nil
***REMOVED***

// GetSandboxPortMapInfo retrieves the current port-mapping programmed for the given sandbox
func GetSandboxPortMapInfo(sb libnetwork.Sandbox) nat.PortMap ***REMOVED***
	pm := nat.PortMap***REMOVED******REMOVED***
	if sb == nil ***REMOVED***
		return pm
	***REMOVED***

	for _, ep := range sb.Endpoints() ***REMOVED***
		pm, _ = getEndpointPortMapInfo(ep)
		if len(pm) > 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return pm
***REMOVED***

// BuildEndpointInfo sets endpoint-related fields on container.NetworkSettings based on the provided network and endpoint.
func (container *Container) BuildEndpointInfo(n libnetwork.Network, ep libnetwork.Endpoint) error ***REMOVED***
	if ep == nil ***REMOVED***
		return errInvalidEndpoint
	***REMOVED***

	networkSettings := container.NetworkSettings
	if networkSettings == nil ***REMOVED***
		return errInvalidNetwork
	***REMOVED***

	epInfo := ep.Info()
	if epInfo == nil ***REMOVED***
		// It is not an error to get an empty endpoint info
		return nil
	***REMOVED***

	if _, ok := networkSettings.Networks[n.Name()]; !ok ***REMOVED***
		networkSettings.Networks[n.Name()] = &network.EndpointSettings***REMOVED***
			EndpointSettings: &networktypes.EndpointSettings***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED***
	networkSettings.Networks[n.Name()].NetworkID = n.ID()
	networkSettings.Networks[n.Name()].EndpointID = ep.ID()

	iface := epInfo.Iface()
	if iface == nil ***REMOVED***
		return nil
	***REMOVED***

	if iface.MacAddress() != nil ***REMOVED***
		networkSettings.Networks[n.Name()].MacAddress = iface.MacAddress().String()
	***REMOVED***

	if iface.Address() != nil ***REMOVED***
		ones, _ := iface.Address().Mask.Size()
		networkSettings.Networks[n.Name()].IPAddress = iface.Address().IP.String()
		networkSettings.Networks[n.Name()].IPPrefixLen = ones
	***REMOVED***

	if iface.AddressIPv6() != nil && iface.AddressIPv6().IP.To16() != nil ***REMOVED***
		onesv6, _ := iface.AddressIPv6().Mask.Size()
		networkSettings.Networks[n.Name()].GlobalIPv6Address = iface.AddressIPv6().IP.String()
		networkSettings.Networks[n.Name()].GlobalIPv6PrefixLen = onesv6
	***REMOVED***

	return nil
***REMOVED***

type named interface ***REMOVED***
	Name() string
***REMOVED***

// UpdateJoinInfo updates network settings when container joins network n with endpoint ep.
func (container *Container) UpdateJoinInfo(n named, ep libnetwork.Endpoint) error ***REMOVED***
	if err := container.buildPortMapInfo(ep); err != nil ***REMOVED***
		return err
	***REMOVED***

	epInfo := ep.Info()
	if epInfo == nil ***REMOVED***
		// It is not an error to get an empty endpoint info
		return nil
	***REMOVED***
	if epInfo.Gateway() != nil ***REMOVED***
		container.NetworkSettings.Networks[n.Name()].Gateway = epInfo.Gateway().String()
	***REMOVED***
	if epInfo.GatewayIPv6().To16() != nil ***REMOVED***
		container.NetworkSettings.Networks[n.Name()].IPv6Gateway = epInfo.GatewayIPv6().String()
	***REMOVED***

	return nil
***REMOVED***

// UpdateSandboxNetworkSettings updates the sandbox ID and Key.
func (container *Container) UpdateSandboxNetworkSettings(sb libnetwork.Sandbox) error ***REMOVED***
	container.NetworkSettings.SandboxID = sb.ID()
	container.NetworkSettings.SandboxKey = sb.Key()
	return nil
***REMOVED***

// BuildJoinOptions builds endpoint Join options from a given network.
func (container *Container) BuildJoinOptions(n named) ([]libnetwork.EndpointOption, error) ***REMOVED***
	var joinOptions []libnetwork.EndpointOption
	if epConfig, ok := container.NetworkSettings.Networks[n.Name()]; ok ***REMOVED***
		for _, str := range epConfig.Links ***REMOVED***
			name, alias, err := opts.ParseLink(str)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			joinOptions = append(joinOptions, libnetwork.CreateOptionAlias(name, alias))
		***REMOVED***
		for k, v := range epConfig.DriverOpts ***REMOVED***
			joinOptions = append(joinOptions, libnetwork.EndpointOptionGeneric(options.Generic***REMOVED***k: v***REMOVED***))
		***REMOVED***
	***REMOVED***

	return joinOptions, nil
***REMOVED***

// BuildCreateEndpointOptions builds endpoint options from a given network.
func (container *Container) BuildCreateEndpointOptions(n libnetwork.Network, epConfig *networktypes.EndpointSettings, sb libnetwork.Sandbox, daemonDNS []string) ([]libnetwork.EndpointOption, error) ***REMOVED***
	var (
		bindings      = make(nat.PortMap)
		pbList        []types.PortBinding
		exposeList    []types.TransportPort
		createOptions []libnetwork.EndpointOption
	)

	defaultNetName := runconfig.DefaultDaemonNetworkMode().NetworkName()

	if (!container.EnableServiceDiscoveryOnDefaultNetwork() && n.Name() == defaultNetName) ||
		container.NetworkSettings.IsAnonymousEndpoint ***REMOVED***
		createOptions = append(createOptions, libnetwork.CreateOptionAnonymous())
	***REMOVED***

	if epConfig != nil ***REMOVED***
		ipam := epConfig.IPAMConfig

		if ipam != nil ***REMOVED***
			var (
				ipList          []net.IP
				ip, ip6, linkip net.IP
			)

			for _, ips := range ipam.LinkLocalIPs ***REMOVED***
				if linkip = net.ParseIP(ips); linkip == nil && ips != "" ***REMOVED***
					return nil, errors.Errorf("Invalid link-local IP address: %s", ipam.LinkLocalIPs)
				***REMOVED***
				ipList = append(ipList, linkip)

			***REMOVED***

			if ip = net.ParseIP(ipam.IPv4Address); ip == nil && ipam.IPv4Address != "" ***REMOVED***
				return nil, errors.Errorf("Invalid IPv4 address: %s)", ipam.IPv4Address)
			***REMOVED***

			if ip6 = net.ParseIP(ipam.IPv6Address); ip6 == nil && ipam.IPv6Address != "" ***REMOVED***
				return nil, errors.Errorf("Invalid IPv6 address: %s)", ipam.IPv6Address)
			***REMOVED***

			createOptions = append(createOptions,
				libnetwork.CreateOptionIpam(ip, ip6, ipList, nil))

		***REMOVED***

		for _, alias := range epConfig.Aliases ***REMOVED***
			createOptions = append(createOptions, libnetwork.CreateOptionMyAlias(alias))
		***REMOVED***
		for k, v := range epConfig.DriverOpts ***REMOVED***
			createOptions = append(createOptions, libnetwork.EndpointOptionGeneric(options.Generic***REMOVED***k: v***REMOVED***))
		***REMOVED***
	***REMOVED***

	if container.NetworkSettings.Service != nil ***REMOVED***
		svcCfg := container.NetworkSettings.Service

		var vip string
		if svcCfg.VirtualAddresses[n.ID()] != nil ***REMOVED***
			vip = svcCfg.VirtualAddresses[n.ID()].IPv4
		***REMOVED***

		var portConfigs []*libnetwork.PortConfig
		for _, portConfig := range svcCfg.ExposedPorts ***REMOVED***
			portConfigs = append(portConfigs, &libnetwork.PortConfig***REMOVED***
				Name:          portConfig.Name,
				Protocol:      libnetwork.PortConfig_Protocol(portConfig.Protocol),
				TargetPort:    portConfig.TargetPort,
				PublishedPort: portConfig.PublishedPort,
			***REMOVED***)
		***REMOVED***

		createOptions = append(createOptions, libnetwork.CreateOptionService(svcCfg.Name, svcCfg.ID, net.ParseIP(vip), portConfigs, svcCfg.Aliases[n.ID()]))
	***REMOVED***

	if !containertypes.NetworkMode(n.Name()).IsUserDefined() ***REMOVED***
		createOptions = append(createOptions, libnetwork.CreateOptionDisableResolution())
	***REMOVED***

	// configs that are applicable only for the endpoint in the network
	// to which container was connected to on docker run.
	// Ideally all these network-specific endpoint configurations must be moved under
	// container.NetworkSettings.Networks[n.Name()]
	if n.Name() == container.HostConfig.NetworkMode.NetworkName() ||
		(n.Name() == defaultNetName && container.HostConfig.NetworkMode.IsDefault()) ***REMOVED***
		if container.Config.MacAddress != "" ***REMOVED***
			mac, err := net.ParseMAC(container.Config.MacAddress)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			genericOption := options.Generic***REMOVED***
				netlabel.MacAddress: mac,
			***REMOVED***

			createOptions = append(createOptions, libnetwork.EndpointOptionGeneric(genericOption))
		***REMOVED***

	***REMOVED***

	// Port-mapping rules belong to the container & applicable only to non-internal networks
	portmaps := GetSandboxPortMapInfo(sb)
	if n.Info().Internal() || len(portmaps) > 0 ***REMOVED***
		return createOptions, nil
	***REMOVED***

	if container.HostConfig.PortBindings != nil ***REMOVED***
		for p, b := range container.HostConfig.PortBindings ***REMOVED***
			bindings[p] = []nat.PortBinding***REMOVED******REMOVED***
			for _, bb := range b ***REMOVED***
				bindings[p] = append(bindings[p], nat.PortBinding***REMOVED***
					HostIP:   bb.HostIP,
					HostPort: bb.HostPort,
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	portSpecs := container.Config.ExposedPorts
	ports := make([]nat.Port, len(portSpecs))
	var i int
	for p := range portSpecs ***REMOVED***
		ports[i] = p
		i++
	***REMOVED***
	nat.SortPortMap(ports, bindings)
	for _, port := range ports ***REMOVED***
		expose := types.TransportPort***REMOVED******REMOVED***
		expose.Proto = types.ParseProtocol(port.Proto())
		expose.Port = uint16(port.Int())
		exposeList = append(exposeList, expose)

		pb := types.PortBinding***REMOVED***Port: expose.Port, Proto: expose.Proto***REMOVED***
		binding := bindings[port]
		for i := 0; i < len(binding); i++ ***REMOVED***
			pbCopy := pb.GetCopy()
			newP, err := nat.NewPort(nat.SplitProtoPort(binding[i].HostPort))
			var portStart, portEnd int
			if err == nil ***REMOVED***
				portStart, portEnd, err = newP.Range()
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, errors.Wrapf(err, "Error parsing HostPort value (%s)", binding[i].HostPort)
			***REMOVED***
			pbCopy.HostPort = uint16(portStart)
			pbCopy.HostPortEnd = uint16(portEnd)
			pbCopy.HostIP = net.ParseIP(binding[i].HostIP)
			pbList = append(pbList, pbCopy)
		***REMOVED***

		if container.HostConfig.PublishAllPorts && len(binding) == 0 ***REMOVED***
			pbList = append(pbList, pb)
		***REMOVED***
	***REMOVED***

	var dns []string

	if len(container.HostConfig.DNS) > 0 ***REMOVED***
		dns = container.HostConfig.DNS
	***REMOVED*** else if len(daemonDNS) > 0 ***REMOVED***
		dns = daemonDNS
	***REMOVED***

	if len(dns) > 0 ***REMOVED***
		createOptions = append(createOptions,
			libnetwork.CreateOptionDNS(dns))
	***REMOVED***

	createOptions = append(createOptions,
		libnetwork.CreateOptionPortMapping(pbList),
		libnetwork.CreateOptionExposedPorts(exposeList))

	return createOptions, nil
***REMOVED***

// UpdateMonitor updates monitor configure for running container
func (container *Container) UpdateMonitor(restartPolicy containertypes.RestartPolicy) ***REMOVED***
	type policySetter interface ***REMOVED***
		SetPolicy(containertypes.RestartPolicy)
	***REMOVED***

	if rm, ok := container.RestartManager().(policySetter); ok ***REMOVED***
		rm.SetPolicy(restartPolicy)
	***REMOVED***
***REMOVED***

// FullHostname returns hostname and optional domain appended to it.
func (container *Container) FullHostname() string ***REMOVED***
	fullHostname := container.Config.Hostname
	if container.Config.Domainname != "" ***REMOVED***
		fullHostname = fmt.Sprintf("%s.%s", fullHostname, container.Config.Domainname)
	***REMOVED***
	return fullHostname
***REMOVED***

// RestartManager returns the current restartmanager instance connected to container.
func (container *Container) RestartManager() restartmanager.RestartManager ***REMOVED***
	if container.restartManager == nil ***REMOVED***
		container.restartManager = restartmanager.New(container.HostConfig.RestartPolicy, container.RestartCount)
	***REMOVED***
	return container.restartManager
***REMOVED***

// ResetRestartManager initializes new restartmanager based on container config
func (container *Container) ResetRestartManager(resetCount bool) ***REMOVED***
	if container.restartManager != nil ***REMOVED***
		container.restartManager.Cancel()
	***REMOVED***
	if resetCount ***REMOVED***
		container.RestartCount = 0
	***REMOVED***
	container.restartManager = nil
***REMOVED***

type attachContext struct ***REMOVED***
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
***REMOVED***

// InitAttachContext initializes or returns existing context for attach calls to
// track container liveness.
func (container *Container) InitAttachContext() context.Context ***REMOVED***
	container.attachContext.mu.Lock()
	defer container.attachContext.mu.Unlock()
	if container.attachContext.ctx == nil ***REMOVED***
		container.attachContext.ctx, container.attachContext.cancel = context.WithCancel(context.Background())
	***REMOVED***
	return container.attachContext.ctx
***REMOVED***

// CancelAttachContext cancels attach context. All attach calls should detach
// after this call.
func (container *Container) CancelAttachContext() ***REMOVED***
	container.attachContext.mu.Lock()
	if container.attachContext.ctx != nil ***REMOVED***
		container.attachContext.cancel()
		container.attachContext.ctx = nil
	***REMOVED***
	container.attachContext.mu.Unlock()
***REMOVED***

func (container *Container) startLogging() error ***REMOVED***
	if container.HostConfig.LogConfig.Type == "none" ***REMOVED***
		return nil // do not start logging routines
	***REMOVED***

	l, err := container.StartLogger()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to initialize logging driver: %v", err)
	***REMOVED***

	copier := logger.NewCopier(map[string]io.Reader***REMOVED***"stdout": container.StdoutPipe(), "stderr": container.StderrPipe()***REMOVED***, l)
	container.LogCopier = copier
	copier.Run()
	container.LogDriver = l

	// set LogPath field only for json-file logdriver
	if jl, ok := l.(*jsonfilelog.JSONFileLogger); ok ***REMOVED***
		container.LogPath = jl.LogPath()
	***REMOVED***

	return nil
***REMOVED***

// StdinPipe gets the stdin stream of the container
func (container *Container) StdinPipe() io.WriteCloser ***REMOVED***
	return container.StreamConfig.StdinPipe()
***REMOVED***

// StdoutPipe gets the stdout stream of the container
func (container *Container) StdoutPipe() io.ReadCloser ***REMOVED***
	return container.StreamConfig.StdoutPipe()
***REMOVED***

// StderrPipe gets the stderr stream of the container
func (container *Container) StderrPipe() io.ReadCloser ***REMOVED***
	return container.StreamConfig.StderrPipe()
***REMOVED***

// CloseStreams closes the container's stdio streams
func (container *Container) CloseStreams() error ***REMOVED***
	return container.StreamConfig.CloseStreams()
***REMOVED***

// InitializeStdio is called by libcontainerd to connect the stdio.
func (container *Container) InitializeStdio(iop *cio.DirectIO) (cio.IO, error) ***REMOVED***
	if err := container.startLogging(); err != nil ***REMOVED***
		container.Reset(false)
		return nil, err
	***REMOVED***

	container.StreamConfig.CopyToPipe(iop)

	if container.StreamConfig.Stdin() == nil && !container.Config.Tty ***REMOVED***
		if iop.Stdin != nil ***REMOVED***
			if err := iop.Stdin.Close(); err != nil ***REMOVED***
				logrus.Warnf("error closing stdin: %+v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return &rio***REMOVED***IO: iop, sc: container.StreamConfig***REMOVED***, nil
***REMOVED***

// MountsResourcePath returns the path where mounts are stored for the given mount
func (container *Container) MountsResourcePath(mount string) (string, error) ***REMOVED***
	return container.GetRootResourcePath(filepath.Join("mounts", mount))
***REMOVED***

// SecretMountPath returns the path of the secret mount for the container
func (container *Container) SecretMountPath() (string, error) ***REMOVED***
	return container.MountsResourcePath("secrets")
***REMOVED***

// SecretFilePath returns the path to the location of a secret on the host.
func (container *Container) SecretFilePath(secretRef swarmtypes.SecretReference) (string, error) ***REMOVED***
	secrets, err := container.SecretMountPath()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return filepath.Join(secrets, secretRef.SecretID), nil
***REMOVED***

func getSecretTargetPath(r *swarmtypes.SecretReference) string ***REMOVED***
	if filepath.IsAbs(r.File.Name) ***REMOVED***
		return r.File.Name
	***REMOVED***

	return filepath.Join(containerSecretMountPath, r.File.Name)
***REMOVED***

// ConfigsDirPath returns the path to the directory where configs are stored on
// disk.
func (container *Container) ConfigsDirPath() (string, error) ***REMOVED***
	return container.GetRootResourcePath("configs")
***REMOVED***

// ConfigFilePath returns the path to the on-disk location of a config.
func (container *Container) ConfigFilePath(configRef swarmtypes.ConfigReference) (string, error) ***REMOVED***
	configs, err := container.ConfigsDirPath()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return filepath.Join(configs, configRef.ConfigID), nil
***REMOVED***

// CreateDaemonEnvironment creates a new environment variable slice for this container.
func (container *Container) CreateDaemonEnvironment(tty bool, linkedEnv []string) []string ***REMOVED***
	// Setup environment
	os := container.OS
	if os == "" ***REMOVED***
		os = runtime.GOOS
	***REMOVED***
	env := []string***REMOVED******REMOVED***
	if runtime.GOOS != "windows" || (runtime.GOOS == "windows" && os == "linux") ***REMOVED***
		env = []string***REMOVED***
			"PATH=" + system.DefaultPathEnv(os),
			"HOSTNAME=" + container.Config.Hostname,
		***REMOVED***
		if tty ***REMOVED***
			env = append(env, "TERM=xterm")
		***REMOVED***
		env = append(env, linkedEnv...)
	***REMOVED***

	// because the env on the container can override certain default values
	// we need to replace the 'env' keys where they match and append anything
	// else.
	env = ReplaceOrAppendEnvValues(env, container.Config.Env)
	return env
***REMOVED***

type rio struct ***REMOVED***
	cio.IO

	sc *stream.Config
***REMOVED***

func (i *rio) Close() error ***REMOVED***
	i.IO.Close()

	return i.sc.CloseStreams()
***REMOVED***

func (i *rio) Wait() ***REMOVED***
	i.sc.Wait()

	i.IO.Wait()
***REMOVED***
