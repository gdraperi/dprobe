// Package daemon exposes the functions that occur on the host server
// that the Docker daemon is running.
//
// In implementing the various functions of the daemon, there is often
// a method-specific struct for configuring the runtime behavior.
package daemon

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/daemon/discovery"
	"github.com/docker/docker/daemon/events"
	"github.com/docker/docker/daemon/exec"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/network"
	"github.com/docker/docker/errdefs"
	"github.com/sirupsen/logrus"
	// register graph drivers
	_ "github.com/docker/docker/daemon/graphdriver/register"
	"github.com/docker/docker/daemon/initlayer"
	"github.com/docker/docker/daemon/stats"
	dmetadata "github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/libcontainerd"
	"github.com/docker/docker/migrate/v1"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/pkg/truncindex"
	"github.com/docker/docker/plugin"
	pluginexec "github.com/docker/docker/plugin/executor/containerd"
	refstore "github.com/docker/docker/reference"
	"github.com/docker/docker/registry"
	"github.com/docker/docker/runconfig"
	volumedrivers "github.com/docker/docker/volume/drivers"
	"github.com/docker/docker/volume/local"
	"github.com/docker/docker/volume/store"
	"github.com/docker/libnetwork"
	"github.com/docker/libnetwork/cluster"
	nwconfig "github.com/docker/libnetwork/config"
	"github.com/docker/libtrust"
	"github.com/pkg/errors"
)

// ContainersNamespace is the name of the namespace used for users containers
const ContainersNamespace = "moby"

var (
	errSystemNotSupported = errors.New("the Docker daemon is not supported on this platform")
)

// Daemon holds information about the Docker daemon.
type Daemon struct ***REMOVED***
	ID                        string
	repository                string
	containers                container.Store
	containersReplica         container.ViewDB
	execCommands              *exec.Store
	downloadManager           *xfer.LayerDownloadManager
	uploadManager             *xfer.LayerUploadManager
	trustKey                  libtrust.PrivateKey
	idIndex                   *truncindex.TruncIndex
	configStore               *config.Config
	statsCollector            *stats.Collector
	defaultLogConfig          containertypes.LogConfig
	RegistryService           registry.Service
	EventsService             *events.Events
	netController             libnetwork.NetworkController
	volumes                   *store.VolumeStore
	discoveryWatcher          discovery.Reloader
	root                      string
	seccompEnabled            bool
	apparmorEnabled           bool
	shutdown                  bool
	idMappings                *idtools.IDMappings
	graphDrivers              map[string]string // By operating system
	referenceStore            refstore.Store
	imageStore                image.Store
	imageRoot                 string
	layerStores               map[string]layer.Store // By operating system
	distributionMetadataStore dmetadata.Store
	PluginStore               *plugin.Store // todo: remove
	pluginManager             *plugin.Manager
	linkIndex                 *linkIndex
	containerd                libcontainerd.Client
	containerdRemote          libcontainerd.Remote
	defaultIsolation          containertypes.Isolation // Default isolation mode on Windows
	clusterProvider           cluster.Provider
	cluster                   Cluster
	genericResources          []swarm.GenericResource
	metricsPluginListener     net.Listener

	machineMemory uint64

	seccompProfile     []byte
	seccompProfilePath string

	diskUsageRunning int32
	pruneRunning     int32
	hosts            map[string]bool // hosts stores the addresses the daemon is listening on
	startupDone      chan struct***REMOVED******REMOVED***

	attachmentStore network.AttachmentStore
***REMOVED***

// StoreHosts stores the addresses the daemon is listening on
func (daemon *Daemon) StoreHosts(hosts []string) ***REMOVED***
	if daemon.hosts == nil ***REMOVED***
		daemon.hosts = make(map[string]bool)
	***REMOVED***
	for _, h := range hosts ***REMOVED***
		daemon.hosts[h] = true
	***REMOVED***
***REMOVED***

// HasExperimental returns whether the experimental features of the daemon are enabled or not
func (daemon *Daemon) HasExperimental() bool ***REMOVED***
	return daemon.configStore != nil && daemon.configStore.Experimental
***REMOVED***

func (daemon *Daemon) restore() error ***REMOVED***
	containers := make(map[string]*container.Container)

	logrus.Info("Loading containers: start.")

	dir, err := ioutil.ReadDir(daemon.repository)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, v := range dir ***REMOVED***
		id := v.Name()
		container, err := daemon.load(id)
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to load container %v: %v", id, err)
			continue
		***REMOVED***
		if !system.IsOSSupported(container.OS) ***REMOVED***
			logrus.Errorf("Failed to load container %v: %s (%q)", id, system.ErrNotSupportedOperatingSystem, container.OS)
			continue
		***REMOVED***
		// Ignore the container if it does not support the current driver being used by the graph
		currentDriverForContainerOS := daemon.graphDrivers[container.OS]
		if (container.Driver == "" && currentDriverForContainerOS == "aufs") || container.Driver == currentDriverForContainerOS ***REMOVED***
			rwlayer, err := daemon.layerStores[container.OS].GetRWLayer(container.ID)
			if err != nil ***REMOVED***
				logrus.Errorf("Failed to load container mount %v: %v", id, err)
				continue
			***REMOVED***
			container.RWLayer = rwlayer
			logrus.Debugf("Loaded container %v, isRunning: %v", container.ID, container.IsRunning())

			containers[container.ID] = container
		***REMOVED*** else ***REMOVED***
			logrus.Debugf("Cannot load container %s because it was created with another graph driver.", container.ID)
		***REMOVED***
	***REMOVED***

	removeContainers := make(map[string]*container.Container)
	restartContainers := make(map[*container.Container]chan struct***REMOVED******REMOVED***)
	activeSandboxes := make(map[string]interface***REMOVED******REMOVED***)
	for id, c := range containers ***REMOVED***
		if err := daemon.registerName(c); err != nil ***REMOVED***
			logrus.Errorf("Failed to register container name %s: %s", c.ID, err)
			delete(containers, id)
			continue
		***REMOVED***
		// verify that all volumes valid and have been migrated from the pre-1.7 layout
		if err := daemon.verifyVolumesInfo(c); err != nil ***REMOVED***
			// don't skip the container due to error
			logrus.Errorf("Failed to verify volumes for container '%s': %v", c.ID, err)
		***REMOVED***
		if err := daemon.Register(c); err != nil ***REMOVED***
			logrus.Errorf("Failed to register container %s: %s", c.ID, err)
			delete(containers, id)
			continue
		***REMOVED***

		// The LogConfig.Type is empty if the container was created before docker 1.12 with default log driver.
		// We should rewrite it to use the daemon defaults.
		// Fixes https://github.com/docker/docker/issues/22536
		if c.HostConfig.LogConfig.Type == "" ***REMOVED***
			if err := daemon.mergeAndVerifyLogConfig(&c.HostConfig.LogConfig); err != nil ***REMOVED***
				logrus.Errorf("Failed to verify log config for container %s: %q", c.ID, err)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var (
		wg      sync.WaitGroup
		mapLock sync.Mutex
	)
	for _, c := range containers ***REMOVED***
		wg.Add(1)
		go func(c *container.Container) ***REMOVED***
			defer wg.Done()
			daemon.backportMountSpec(c)
			if err := daemon.checkpointAndSave(c); err != nil ***REMOVED***
				logrus.WithError(err).WithField("container", c.ID).Error("error saving backported mountspec to disk")
			***REMOVED***

			daemon.setStateCounter(c)

			logrus.WithFields(logrus.Fields***REMOVED***
				"container": c.ID,
				"running":   c.IsRunning(),
				"paused":    c.IsPaused(),
			***REMOVED***).Debug("restoring container")

			var (
				err      error
				alive    bool
				ec       uint32
				exitedAt time.Time
			)

			alive, _, err = daemon.containerd.Restore(context.Background(), c.ID, c.InitializeStdio)
			if err != nil && !errdefs.IsNotFound(err) ***REMOVED***
				logrus.Errorf("Failed to restore container %s with containerd: %s", c.ID, err)
				return
			***REMOVED***
			if !alive ***REMOVED***
				ec, exitedAt, err = daemon.containerd.DeleteTask(context.Background(), c.ID)
				if err != nil && !errdefs.IsNotFound(err) ***REMOVED***
					logrus.WithError(err).Errorf("Failed to delete container %s from containerd", c.ID)
					return
				***REMOVED***
			***REMOVED*** else if !daemon.configStore.LiveRestoreEnabled ***REMOVED***
				if err := daemon.kill(c, c.StopSignal()); err != nil && !errdefs.IsNotFound(err) ***REMOVED***
					logrus.WithError(err).WithField("container", c.ID).Error("error shutting down container")
					return
				***REMOVED***
			***REMOVED***

			if c.IsRunning() || c.IsPaused() ***REMOVED***
				c.RestartManager().Cancel() // manually start containers because some need to wait for swarm networking

				if c.IsPaused() && alive ***REMOVED***
					s, err := daemon.containerd.Status(context.Background(), c.ID)
					if err != nil ***REMOVED***
						logrus.WithError(err).WithField("container", c.ID).
							Errorf("Failed to get container status")
					***REMOVED*** else ***REMOVED***
						logrus.WithField("container", c.ID).WithField("state", s).
							Info("restored container paused")
						switch s ***REMOVED***
						case libcontainerd.StatusPaused, libcontainerd.StatusPausing:
							// nothing to do
						case libcontainerd.StatusStopped:
							alive = false
						case libcontainerd.StatusUnknown:
							logrus.WithField("container", c.ID).
								Error("Unknown status for container during restore")
						default:
							// running
							c.Lock()
							c.Paused = false
							daemon.setStateCounter(c)
							if err := c.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
								logrus.WithError(err).WithField("container", c.ID).
									Error("Failed to update stopped container state")
							***REMOVED***
							c.Unlock()
						***REMOVED***
					***REMOVED***
				***REMOVED***

				if !alive ***REMOVED***
					c.Lock()
					c.SetStopped(&container.ExitStatus***REMOVED***ExitCode: int(ec), ExitedAt: exitedAt***REMOVED***)
					daemon.Cleanup(c)
					if err := c.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
						logrus.Errorf("Failed to update stopped container %s state: %v", c.ID, err)
					***REMOVED***
					c.Unlock()
				***REMOVED***

				// we call Mount and then Unmount to get BaseFs of the container
				if err := daemon.Mount(c); err != nil ***REMOVED***
					// The mount is unlikely to fail. However, in case mount fails
					// the container should be allowed to restore here. Some functionalities
					// (like docker exec -u user) might be missing but container is able to be
					// stopped/restarted/removed.
					// See #29365 for related information.
					// The error is only logged here.
					logrus.Warnf("Failed to mount container on getting BaseFs path %v: %v", c.ID, err)
				***REMOVED*** else ***REMOVED***
					if err := daemon.Unmount(c); err != nil ***REMOVED***
						logrus.Warnf("Failed to umount container on getting BaseFs path %v: %v", c.ID, err)
					***REMOVED***
				***REMOVED***

				c.ResetRestartManager(false)
				if !c.HostConfig.NetworkMode.IsContainer() && c.IsRunning() ***REMOVED***
					options, err := daemon.buildSandboxOptions(c)
					if err != nil ***REMOVED***
						logrus.Warnf("Failed build sandbox option to restore container %s: %v", c.ID, err)
					***REMOVED***
					mapLock.Lock()
					activeSandboxes[c.NetworkSettings.SandboxID] = options
					mapLock.Unlock()
				***REMOVED***
			***REMOVED***

			// get list of containers we need to restart

			// Do not autostart containers which
			// has endpoints in a swarm scope
			// network yet since the cluster is
			// not initialized yet. We will start
			// it after the cluster is
			// initialized.
			if daemon.configStore.AutoRestart && c.ShouldRestart() && !c.NetworkSettings.HasSwarmEndpoint ***REMOVED***
				mapLock.Lock()
				restartContainers[c] = make(chan struct***REMOVED******REMOVED***)
				mapLock.Unlock()
			***REMOVED*** else if c.HostConfig != nil && c.HostConfig.AutoRemove ***REMOVED***
				mapLock.Lock()
				removeContainers[c.ID] = c
				mapLock.Unlock()
			***REMOVED***

			c.Lock()
			if c.RemovalInProgress ***REMOVED***
				// We probably crashed in the middle of a removal, reset
				// the flag.
				//
				// We DO NOT remove the container here as we do not
				// know if the user had requested for either the
				// associated volumes, network links or both to also
				// be removed. So we put the container in the "dead"
				// state and leave further processing up to them.
				logrus.Debugf("Resetting RemovalInProgress flag from %v", c.ID)
				c.RemovalInProgress = false
				c.Dead = true
				if err := c.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
					logrus.Errorf("Failed to update RemovalInProgress container %s state: %v", c.ID, err)
				***REMOVED***
			***REMOVED***
			c.Unlock()
		***REMOVED***(c)
	***REMOVED***
	wg.Wait()
	daemon.netController, err = daemon.initNetworkController(daemon.configStore, activeSandboxes)
	if err != nil ***REMOVED***
		return fmt.Errorf("Error initializing network controller: %v", err)
	***REMOVED***

	// Now that all the containers are registered, register the links
	for _, c := range containers ***REMOVED***
		if err := daemon.registerLinks(c, c.HostConfig); err != nil ***REMOVED***
			logrus.Errorf("failed to register link for container %s: %v", c.ID, err)
		***REMOVED***
	***REMOVED***

	group := sync.WaitGroup***REMOVED******REMOVED***
	for c, notifier := range restartContainers ***REMOVED***
		group.Add(1)

		go func(c *container.Container, chNotify chan struct***REMOVED******REMOVED***) ***REMOVED***
			defer group.Done()

			logrus.Debugf("Starting container %s", c.ID)

			// ignore errors here as this is a best effort to wait for children to be
			//   running before we try to start the container
			children := daemon.children(c)
			timeout := time.After(5 * time.Second)
			for _, child := range children ***REMOVED***
				if notifier, exists := restartContainers[child]; exists ***REMOVED***
					select ***REMOVED***
					case <-notifier:
					case <-timeout:
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Make sure networks are available before starting
			daemon.waitForNetworks(c)
			if err := daemon.containerStart(c, "", "", true); err != nil ***REMOVED***
				logrus.Errorf("Failed to start container %s: %s", c.ID, err)
			***REMOVED***
			close(chNotify)
		***REMOVED***(c, notifier)

	***REMOVED***
	group.Wait()

	removeGroup := sync.WaitGroup***REMOVED******REMOVED***
	for id := range removeContainers ***REMOVED***
		removeGroup.Add(1)
		go func(cid string) ***REMOVED***
			if err := daemon.ContainerRm(cid, &types.ContainerRmConfig***REMOVED***ForceRemove: true, RemoveVolume: true***REMOVED***); err != nil ***REMOVED***
				logrus.Errorf("Failed to remove container %s: %s", cid, err)
			***REMOVED***
			removeGroup.Done()
		***REMOVED***(id)
	***REMOVED***
	removeGroup.Wait()

	// any containers that were started above would already have had this done,
	// however we need to now prepare the mountpoints for the rest of the containers as well.
	// This shouldn't cause any issue running on the containers that already had this run.
	// This must be run after any containers with a restart policy so that containerized plugins
	// can have a chance to be running before we try to initialize them.
	for _, c := range containers ***REMOVED***
		// if the container has restart policy, do not
		// prepare the mountpoints since it has been done on restarting.
		// This is to speed up the daemon start when a restart container
		// has a volume and the volume driver is not available.
		if _, ok := restartContainers[c]; ok ***REMOVED***
			continue
		***REMOVED*** else if _, ok := removeContainers[c.ID]; ok ***REMOVED***
			// container is automatically removed, skip it.
			continue
		***REMOVED***

		group.Add(1)
		go func(c *container.Container) ***REMOVED***
			defer group.Done()
			if err := daemon.prepareMountPoints(c); err != nil ***REMOVED***
				logrus.Error(err)
			***REMOVED***
		***REMOVED***(c)
	***REMOVED***

	group.Wait()

	logrus.Info("Loading containers: done.")

	return nil
***REMOVED***

// RestartSwarmContainers restarts any autostart container which has a
// swarm endpoint.
func (daemon *Daemon) RestartSwarmContainers() ***REMOVED***
	group := sync.WaitGroup***REMOVED******REMOVED***
	for _, c := range daemon.List() ***REMOVED***
		if !c.IsRunning() && !c.IsPaused() ***REMOVED***
			// Autostart all the containers which has a
			// swarm endpoint now that the cluster is
			// initialized.
			if daemon.configStore.AutoRestart && c.ShouldRestart() && c.NetworkSettings.HasSwarmEndpoint ***REMOVED***
				group.Add(1)
				go func(c *container.Container) ***REMOVED***
					defer group.Done()
					if err := daemon.containerStart(c, "", "", true); err != nil ***REMOVED***
						logrus.Error(err)
					***REMOVED***
				***REMOVED***(c)
			***REMOVED***
		***REMOVED***

	***REMOVED***
	group.Wait()
***REMOVED***

// waitForNetworks is used during daemon initialization when starting up containers
// It ensures that all of a container's networks are available before the daemon tries to start the container.
// In practice it just makes sure the discovery service is available for containers which use a network that require discovery.
func (daemon *Daemon) waitForNetworks(c *container.Container) ***REMOVED***
	if daemon.discoveryWatcher == nil ***REMOVED***
		return
	***REMOVED***
	// Make sure if the container has a network that requires discovery that the discovery service is available before starting
	for netName := range c.NetworkSettings.Networks ***REMOVED***
		// If we get `ErrNoSuchNetwork` here, we can assume that it is due to discovery not being ready
		// Most likely this is because the K/V store used for discovery is in a container and needs to be started
		if _, err := daemon.netController.NetworkByName(netName); err != nil ***REMOVED***
			if _, ok := err.(libnetwork.ErrNoSuchNetwork); !ok ***REMOVED***
				continue
			***REMOVED***
			// use a longish timeout here due to some slowdowns in libnetwork if the k/v store is on anything other than --net=host
			// FIXME: why is this slow???
			logrus.Debugf("Container %s waiting for network to be ready", c.Name)
			select ***REMOVED***
			case <-daemon.discoveryWatcher.ReadyCh():
			case <-time.After(60 * time.Second):
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (daemon *Daemon) children(c *container.Container) map[string]*container.Container ***REMOVED***
	return daemon.linkIndex.children(c)
***REMOVED***

// parents returns the names of the parent containers of the container
// with the given name.
func (daemon *Daemon) parents(c *container.Container) map[string]*container.Container ***REMOVED***
	return daemon.linkIndex.parents(c)
***REMOVED***

func (daemon *Daemon) registerLink(parent, child *container.Container, alias string) error ***REMOVED***
	fullName := path.Join(parent.Name, alias)
	if err := daemon.containersReplica.ReserveName(fullName, child.ID); err != nil ***REMOVED***
		if err == container.ErrNameReserved ***REMOVED***
			logrus.Warnf("error registering link for %s, to %s, as alias %s, ignoring: %v", parent.ID, child.ID, alias, err)
			return nil
		***REMOVED***
		return err
	***REMOVED***
	daemon.linkIndex.link(parent, child, fullName)
	return nil
***REMOVED***

// DaemonJoinsCluster informs the daemon has joined the cluster and provides
// the handler to query the cluster component
func (daemon *Daemon) DaemonJoinsCluster(clusterProvider cluster.Provider) ***REMOVED***
	daemon.setClusterProvider(clusterProvider)
***REMOVED***

// DaemonLeavesCluster informs the daemon has left the cluster
func (daemon *Daemon) DaemonLeavesCluster() ***REMOVED***
	// Daemon is in charge of removing the attachable networks with
	// connected containers when the node leaves the swarm
	daemon.clearAttachableNetworks()
	// We no longer need the cluster provider, stop it now so that
	// the network agent will stop listening to cluster events.
	daemon.setClusterProvider(nil)
	// Wait for the networking cluster agent to stop
	daemon.netController.AgentStopWait()
	// Daemon is in charge of removing the ingress network when the
	// node leaves the swarm. Wait for job to be done or timeout.
	// This is called also on graceful daemon shutdown. We need to
	// wait, because the ingress release has to happen before the
	// network controller is stopped.
	if done, err := daemon.ReleaseIngress(); err == nil ***REMOVED***
		select ***REMOVED***
		case <-done:
		case <-time.After(5 * time.Second):
			logrus.Warnf("timeout while waiting for ingress network removal")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		logrus.Warnf("failed to initiate ingress network removal: %v", err)
	***REMOVED***

	daemon.attachmentStore.ClearAttachments()
***REMOVED***

// setClusterProvider sets a component for querying the current cluster state.
func (daemon *Daemon) setClusterProvider(clusterProvider cluster.Provider) ***REMOVED***
	daemon.clusterProvider = clusterProvider
	daemon.netController.SetClusterProvider(clusterProvider)
***REMOVED***

// IsSwarmCompatible verifies if the current daemon
// configuration is compatible with the swarm mode
func (daemon *Daemon) IsSwarmCompatible() error ***REMOVED***
	if daemon.configStore == nil ***REMOVED***
		return nil
	***REMOVED***
	return daemon.configStore.IsSwarmCompatible()
***REMOVED***

// NewDaemon sets up everything for the daemon to be able to service
// requests from the webserver.
func NewDaemon(config *config.Config, registryService registry.Service, containerdRemote libcontainerd.Remote, pluginStore *plugin.Store) (daemon *Daemon, err error) ***REMOVED***
	setDefaultMtu(config)

	// Ensure that we have a correct root key limit for launching containers.
	if err := ModifyRootKeyLimit(); err != nil ***REMOVED***
		logrus.Warnf("unable to modify root key limit, number of containers could be limited by this quota: %v", err)
	***REMOVED***

	// Ensure we have compatible and valid configuration options
	if err := verifyDaemonSettings(config); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Do we have a disabled network?
	config.DisableBridge = isBridgeNetworkDisabled(config)

	// Verify the platform is supported as a daemon
	if !platformSupported ***REMOVED***
		return nil, errSystemNotSupported
	***REMOVED***

	// Validate platform-specific requirements
	if err := checkSystem(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	idMappings, err := setupRemappedRoot(config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	rootIDs := idMappings.RootPair()
	if err := setupDaemonProcess(config); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// set up the tmpDir to use a canonical path
	tmp, err := prepareTempDir(config.Root, rootIDs)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Unable to get the TempDir under %s: %s", config.Root, err)
	***REMOVED***
	realTmp, err := getRealPath(tmp)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Unable to get the full path to the TempDir (%s): %s", tmp, err)
	***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		if _, err := os.Stat(realTmp); err != nil && os.IsNotExist(err) ***REMOVED***
			if err := system.MkdirAll(realTmp, 0700, ""); err != nil ***REMOVED***
				return nil, fmt.Errorf("Unable to create the TempDir (%s): %s", realTmp, err)
			***REMOVED***
		***REMOVED***
		os.Setenv("TEMP", realTmp)
		os.Setenv("TMP", realTmp)
	***REMOVED*** else ***REMOVED***
		os.Setenv("TMPDIR", realTmp)
	***REMOVED***

	d := &Daemon***REMOVED***
		configStore: config,
		PluginStore: pluginStore,
		startupDone: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	// Ensure the daemon is properly shutdown if there is a failure during
	// initialization
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if err := d.Shutdown(); err != nil ***REMOVED***
				logrus.Error(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err := d.setGenericResources(config); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// set up SIGUSR1 handler on Unix-like systems, or a Win32 global event
	// on Windows to dump Go routine stacks
	stackDumpDir := config.Root
	if execRoot := config.GetExecRoot(); execRoot != "" ***REMOVED***
		stackDumpDir = execRoot
	***REMOVED***
	d.setupDumpStackTrap(stackDumpDir)

	if err := d.setupSeccompProfile(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Set the default isolation mode (only applicable on Windows)
	if err := d.setDefaultIsolation(); err != nil ***REMOVED***
		return nil, fmt.Errorf("error setting default isolation mode: %v", err)
	***REMOVED***

	logrus.Debugf("Using default logging driver %s", config.LogConfig.Type)

	if err := configureMaxThreads(config); err != nil ***REMOVED***
		logrus.Warnf("Failed to configure golang's threads limit: %v", err)
	***REMOVED***

	if err := ensureDefaultAppArmorProfile(); err != nil ***REMOVED***
		logrus.Errorf(err.Error())
	***REMOVED***

	daemonRepo := filepath.Join(config.Root, "containers")
	if err := idtools.MkdirAllAndChown(daemonRepo, 0700, rootIDs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Create the directory where we'll store the runtime scripts (i.e. in
	// order to support runtimeArgs)
	daemonRuntimes := filepath.Join(config.Root, "runtimes")
	if err := system.MkdirAll(daemonRuntimes, 0700, ""); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := d.loadRuntimes(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if runtime.GOOS == "windows" ***REMOVED***
		if err := system.MkdirAll(filepath.Join(config.Root, "credentialspecs"), 0, ""); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// On Windows we don't support the environment variable, or a user supplied graphdriver
	// as Windows has no choice in terms of which graphdrivers to use. It's a case of
	// running Windows containers on Windows - windowsfilter, running Linux containers on Windows,
	// lcow. Unix platforms however run a single graphdriver for all containers, and it can
	// be set through an environment variable, a daemon start parameter, or chosen through
	// initialization of the layerstore through driver priority order for example.
	d.graphDrivers = make(map[string]string)
	d.layerStores = make(map[string]layer.Store)
	if runtime.GOOS == "windows" ***REMOVED***
		d.graphDrivers[runtime.GOOS] = "windowsfilter"
		if system.LCOWSupported() ***REMOVED***
			d.graphDrivers["linux"] = "lcow"
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		driverName := os.Getenv("DOCKER_DRIVER")
		if driverName == "" ***REMOVED***
			driverName = config.GraphDriver
		***REMOVED*** else ***REMOVED***
			logrus.Infof("Setting the storage driver from the $DOCKER_DRIVER environment variable (%s)", driverName)
		***REMOVED***
		d.graphDrivers[runtime.GOOS] = driverName // May still be empty. Layerstore init determines instead.
	***REMOVED***

	d.RegistryService = registryService
	logger.RegisterPluginGetter(d.PluginStore)

	metricsSockPath, err := d.listenMetricsSock()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	registerMetricsPluginCallback(d.PluginStore, metricsSockPath)

	createPluginExec := func(m *plugin.Manager) (plugin.Executor, error) ***REMOVED***
		return pluginexec.New(getPluginExecRoot(config.Root), containerdRemote, m)
	***REMOVED***

	// Plugin system initialization should happen before restore. Do not change order.
	d.pluginManager, err = plugin.NewManager(plugin.ManagerConfig***REMOVED***
		Root:               filepath.Join(config.Root, "plugins"),
		ExecRoot:           getPluginExecRoot(config.Root),
		Store:              d.PluginStore,
		CreateExecutor:     createPluginExec,
		RegistryService:    registryService,
		LiveRestoreEnabled: config.LiveRestoreEnabled,
		LogPluginEvent:     d.LogPluginEvent, // todo: make private
		AuthzMiddleware:    config.AuthzMiddleware,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "couldn't create plugin manager")
	***REMOVED***

	for operatingSystem, gd := range d.graphDrivers ***REMOVED***
		d.layerStores[operatingSystem], err = layer.NewStoreFromOptions(layer.StoreOptions***REMOVED***
			Root: config.Root,
			MetadataStorePathTemplate: filepath.Join(config.Root, "image", "%s", "layerdb"),
			GraphDriver:               gd,
			GraphDriverOptions:        config.GraphOptions,
			IDMappings:                idMappings,
			PluginGetter:              d.PluginStore,
			ExperimentalEnabled:       config.Experimental,
			OS:                        operatingSystem,
		***REMOVED***)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// As layerstore initialization may set the driver
	for os := range d.graphDrivers ***REMOVED***
		d.graphDrivers[os] = d.layerStores[os].DriverName()
	***REMOVED***

	// Configure and validate the kernels security support. Note this is a Linux/FreeBSD
	// operation only, so it is safe to pass *just* the runtime OS graphdriver.
	if err := configureKernelSecuritySupport(config, d.graphDrivers[runtime.GOOS]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	logrus.Debugf("Max Concurrent Downloads: %d", *config.MaxConcurrentDownloads)
	d.downloadManager = xfer.NewLayerDownloadManager(d.layerStores, *config.MaxConcurrentDownloads)
	logrus.Debugf("Max Concurrent Uploads: %d", *config.MaxConcurrentUploads)
	d.uploadManager = xfer.NewLayerUploadManager(*config.MaxConcurrentUploads)

	d.imageRoot = filepath.Join(config.Root, "image", d.graphDrivers[runtime.GOOS])
	ifs, err := image.NewFSStoreBackend(filepath.Join(d.imageRoot, "imagedb"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	lgrMap := make(map[string]image.LayerGetReleaser)
	for os, ls := range d.layerStores ***REMOVED***
		lgrMap[os] = ls
	***REMOVED***
	d.imageStore, err = image.NewImageStore(ifs, lgrMap)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Configure the volumes driver
	volStore, err := d.configureVolumes(rootIDs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	trustKey, err := loadOrCreateTrustKey(config.TrustKeyPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	trustDir := filepath.Join(config.Root, "trust")

	if err := system.MkdirAll(trustDir, 0700, ""); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	eventsService := events.New()

	// We have a single tag/reference store for the daemon globally. However, it's
	// stored under the graphdriver. On host platforms which only support a single
	// container OS, but multiple selectable graphdrivers, this means depending on which
	// graphdriver is chosen, the global reference store is under there. For
	// platforms which support multiple container operating systems, this is slightly
	// more problematic as where does the global ref store get located? Fortunately,
	// for Windows, which is currently the only daemon supporting multiple container
	// operating systems, the list of graphdrivers available isn't user configurable.
	// For backwards compatibility, we just put it under the windowsfilter
	// directory regardless.
	refStoreLocation := filepath.Join(d.imageRoot, `repositories.json`)
	rs, err := refstore.NewReferenceStore(refStoreLocation)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Couldn't create reference store repository: %s", err)
	***REMOVED***
	d.referenceStore = rs

	d.distributionMetadataStore, err = dmetadata.NewFSMetadataStore(filepath.Join(d.imageRoot, "distribution"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// No content-addressability migration on Windows as it never supported pre-CA
	if runtime.GOOS != "windows" ***REMOVED***
		migrationStart := time.Now()
		if err := v1.Migrate(config.Root, d.graphDrivers[runtime.GOOS], d.layerStores[runtime.GOOS], d.imageStore, rs, d.distributionMetadataStore); err != nil ***REMOVED***
			logrus.Errorf("Graph migration failed: %q. Your old graph data was found to be too inconsistent for upgrading to content-addressable storage. Some of the old data was probably not upgraded. We recommend starting over with a clean storage directory if possible.", err)
		***REMOVED***
		logrus.Infof("Graph migration to content-addressability took %.2f seconds", time.Since(migrationStart).Seconds())
	***REMOVED***

	// Discovery is only enabled when the daemon is launched with an address to advertise.  When
	// initialized, the daemon is registered and we can store the discovery backend as it's read-only
	if err := d.initDiscovery(config); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sysInfo := sysinfo.New(false)
	// Check if Devices cgroup is mounted, it is hard requirement for container security,
	// on Linux.
	if runtime.GOOS == "linux" && !sysInfo.CgroupDevicesEnabled ***REMOVED***
		return nil, errors.New("Devices cgroup isn't mounted")
	***REMOVED***

	d.ID = trustKey.PublicKey().KeyID()
	d.repository = daemonRepo
	d.containers = container.NewMemoryStore()
	if d.containersReplica, err = container.NewViewDB(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	d.execCommands = exec.NewStore()
	d.trustKey = trustKey
	d.idIndex = truncindex.NewTruncIndex([]string***REMOVED******REMOVED***)
	d.statsCollector = d.newStatsCollector(1 * time.Second)
	d.defaultLogConfig = containertypes.LogConfig***REMOVED***
		Type:   config.LogConfig.Type,
		Config: config.LogConfig.Config,
	***REMOVED***
	d.EventsService = eventsService
	d.volumes = volStore
	d.root = config.Root
	d.idMappings = idMappings
	d.seccompEnabled = sysInfo.Seccomp
	d.apparmorEnabled = sysInfo.AppArmor
	d.containerdRemote = containerdRemote

	d.linkIndex = newLinkIndex()

	go d.execCommandGC()

	d.containerd, err = containerdRemote.NewClient(ContainersNamespace, d)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := d.restore(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	close(d.startupDone)

	// FIXME: this method never returns an error
	info, _ := d.SystemInfo()

	engineInfo.WithValues(
		dockerversion.Version,
		dockerversion.GitCommit,
		info.Architecture,
		info.Driver,
		info.KernelVersion,
		info.OperatingSystem,
		info.OSType,
		info.ID,
	).Set(1)
	engineCpus.Set(float64(info.NCPU))
	engineMemory.Set(float64(info.MemTotal))

	gd := ""
	for os, driver := range d.graphDrivers ***REMOVED***
		if len(gd) > 0 ***REMOVED***
			gd += ", "
		***REMOVED***
		gd += driver
		if len(d.graphDrivers) > 1 ***REMOVED***
			gd = fmt.Sprintf("%s (%s)", gd, os)
		***REMOVED***
	***REMOVED***
	logrus.WithFields(logrus.Fields***REMOVED***
		"version":        dockerversion.Version,
		"commit":         dockerversion.GitCommit,
		"graphdriver(s)": gd,
	***REMOVED***).Info("Docker daemon")

	return d, nil
***REMOVED***

func (daemon *Daemon) waitForStartupDone() ***REMOVED***
	<-daemon.startupDone
***REMOVED***

func (daemon *Daemon) shutdownContainer(c *container.Container) error ***REMOVED***
	stopTimeout := c.StopTimeout()

	// If container failed to exit in stopTimeout seconds of SIGTERM, then using the force
	if err := daemon.containerStop(c, stopTimeout); err != nil ***REMOVED***
		return fmt.Errorf("Failed to stop container %s with error: %v", c.ID, err)
	***REMOVED***

	// Wait without timeout for the container to exit.
	// Ignore the result.
	<-c.Wait(context.Background(), container.WaitConditionNotRunning)
	return nil
***REMOVED***

// ShutdownTimeout returns the shutdown timeout based on the max stopTimeout of the containers,
// and is limited by daemon's ShutdownTimeout.
func (daemon *Daemon) ShutdownTimeout() int ***REMOVED***
	// By default we use daemon's ShutdownTimeout.
	shutdownTimeout := daemon.configStore.ShutdownTimeout

	graceTimeout := 5
	if daemon.containers != nil ***REMOVED***
		for _, c := range daemon.containers.List() ***REMOVED***
			if shutdownTimeout >= 0 ***REMOVED***
				stopTimeout := c.StopTimeout()
				if stopTimeout < 0 ***REMOVED***
					shutdownTimeout = -1
				***REMOVED*** else ***REMOVED***
					if stopTimeout+graceTimeout > shutdownTimeout ***REMOVED***
						shutdownTimeout = stopTimeout + graceTimeout
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return shutdownTimeout
***REMOVED***

// Shutdown stops the daemon.
func (daemon *Daemon) Shutdown() error ***REMOVED***
	daemon.shutdown = true
	// Keep mounts and networking running on daemon shutdown if
	// we are to keep containers running and restore them.

	if daemon.configStore.LiveRestoreEnabled && daemon.containers != nil ***REMOVED***
		// check if there are any running containers, if none we should do some cleanup
		if ls, err := daemon.Containers(&types.ContainerListOptions***REMOVED******REMOVED***); len(ls) != 0 || err != nil ***REMOVED***
			// metrics plugins still need some cleanup
			daemon.cleanupMetricsPlugins()
			return nil
		***REMOVED***
	***REMOVED***

	if daemon.containers != nil ***REMOVED***
		logrus.Debugf("daemon configured with a %d seconds minimum shutdown timeout", daemon.configStore.ShutdownTimeout)
		logrus.Debugf("start clean shutdown of all containers with a %d seconds timeout...", daemon.ShutdownTimeout())
		daemon.containers.ApplyAll(func(c *container.Container) ***REMOVED***
			if !c.IsRunning() ***REMOVED***
				return
			***REMOVED***
			logrus.Debugf("stopping %s", c.ID)
			if err := daemon.shutdownContainer(c); err != nil ***REMOVED***
				logrus.Errorf("Stop container error: %v", err)
				return
			***REMOVED***
			if mountid, err := daemon.layerStores[c.OS].GetMountID(c.ID); err == nil ***REMOVED***
				daemon.cleanupMountsByID(mountid)
			***REMOVED***
			logrus.Debugf("container stopped %s", c.ID)
		***REMOVED***)
	***REMOVED***

	if daemon.volumes != nil ***REMOVED***
		if err := daemon.volumes.Shutdown(); err != nil ***REMOVED***
			logrus.Errorf("Error shutting down volume store: %v", err)
		***REMOVED***
	***REMOVED***

	for os, ls := range daemon.layerStores ***REMOVED***
		if ls != nil ***REMOVED***
			if err := ls.Cleanup(); err != nil ***REMOVED***
				logrus.Errorf("Error during layer Store.Cleanup(): %v %s", err, os)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// If we are part of a cluster, clean up cluster's stuff
	if daemon.clusterProvider != nil ***REMOVED***
		logrus.Debugf("start clean shutdown of cluster resources...")
		daemon.DaemonLeavesCluster()
	***REMOVED***

	daemon.cleanupMetricsPlugins()

	// Shutdown plugins after containers and layerstore. Don't change the order.
	daemon.pluginShutdown()

	// trigger libnetwork Stop only if it's initialized
	if daemon.netController != nil ***REMOVED***
		daemon.netController.Stop()
	***REMOVED***

	return daemon.cleanupMounts()
***REMOVED***

// Mount sets container.BaseFS
// (is it not set coming in? why is it unset?)
func (daemon *Daemon) Mount(container *container.Container) error ***REMOVED***
	dir, err := container.RWLayer.Mount(container.GetMountLabel())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	logrus.Debugf("container mounted via layerStore: %v", dir)

	if container.BaseFS != nil && container.BaseFS.Path() != dir.Path() ***REMOVED***
		// The mount path reported by the graph driver should always be trusted on Windows, since the
		// volume path for a given mounted layer may change over time.  This should only be an error
		// on non-Windows operating systems.
		if runtime.GOOS != "windows" ***REMOVED***
			daemon.Unmount(container)
			return fmt.Errorf("Error: driver %s is returning inconsistent paths for container %s ('%s' then '%s')",
				daemon.GraphDriverName(container.OS), container.ID, container.BaseFS, dir)
		***REMOVED***
	***REMOVED***
	container.BaseFS = dir // TODO: combine these fields
	return nil
***REMOVED***

// Unmount unsets the container base filesystem
func (daemon *Daemon) Unmount(container *container.Container) error ***REMOVED***
	if err := container.RWLayer.Unmount(); err != nil ***REMOVED***
		logrus.Errorf("Error unmounting container %s: %s", container.ID, err)
		return err
	***REMOVED***

	return nil
***REMOVED***

// Subnets return the IPv4 and IPv6 subnets of networks that are manager by Docker.
func (daemon *Daemon) Subnets() ([]net.IPNet, []net.IPNet) ***REMOVED***
	var v4Subnets []net.IPNet
	var v6Subnets []net.IPNet

	managedNetworks := daemon.netController.Networks()

	for _, managedNetwork := range managedNetworks ***REMOVED***
		v4infos, v6infos := managedNetwork.Info().IpamInfo()
		for _, info := range v4infos ***REMOVED***
			if info.IPAMData.Pool != nil ***REMOVED***
				v4Subnets = append(v4Subnets, *info.IPAMData.Pool)
			***REMOVED***
		***REMOVED***
		for _, info := range v6infos ***REMOVED***
			if info.IPAMData.Pool != nil ***REMOVED***
				v6Subnets = append(v6Subnets, *info.IPAMData.Pool)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return v4Subnets, v6Subnets
***REMOVED***

// GraphDriverName returns the name of the graph driver used by the layer.Store
func (daemon *Daemon) GraphDriverName(os string) string ***REMOVED***
	return daemon.layerStores[os].DriverName()
***REMOVED***

// prepareTempDir prepares and returns the default directory to use
// for temporary files.
// If it doesn't exist, it is created. If it exists, its content is removed.
func prepareTempDir(rootDir string, rootIDs idtools.IDPair) (string, error) ***REMOVED***
	var tmpDir string
	if tmpDir = os.Getenv("DOCKER_TMPDIR"); tmpDir == "" ***REMOVED***
		tmpDir = filepath.Join(rootDir, "tmp")
		newName := tmpDir + "-old"
		if err := os.Rename(tmpDir, newName); err == nil ***REMOVED***
			go func() ***REMOVED***
				if err := os.RemoveAll(newName); err != nil ***REMOVED***
					logrus.Warnf("failed to delete old tmp directory: %s", newName)
				***REMOVED***
			***REMOVED***()
		***REMOVED*** else if !os.IsNotExist(err) ***REMOVED***
			logrus.Warnf("failed to rename %s for background deletion: %s. Deleting synchronously", tmpDir, err)
			if err := os.RemoveAll(tmpDir); err != nil ***REMOVED***
				logrus.Warnf("failed to delete old tmp directory: %s", tmpDir)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// We don't remove the content of tmpdir if it's not the default,
	// it may hold things that do not belong to us.
	return tmpDir, idtools.MkdirAllAndChown(tmpDir, 0700, rootIDs)
***REMOVED***

func (daemon *Daemon) setupInitLayer(initPath containerfs.ContainerFS) error ***REMOVED***
	rootIDs := daemon.idMappings.RootPair()
	return initlayer.Setup(initPath, rootIDs)
***REMOVED***

func (daemon *Daemon) setGenericResources(conf *config.Config) error ***REMOVED***
	genericResources, err := config.ParseGenericResources(conf.NodeGenericResources)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	daemon.genericResources = genericResources

	return nil
***REMOVED***

func setDefaultMtu(conf *config.Config) ***REMOVED***
	// do nothing if the config does not have the default 0 value.
	if conf.Mtu != 0 ***REMOVED***
		return
	***REMOVED***
	conf.Mtu = config.DefaultNetworkMtu
***REMOVED***

func (daemon *Daemon) configureVolumes(rootIDs idtools.IDPair) (*store.VolumeStore, error) ***REMOVED***
	volumesDriver, err := local.New(daemon.configStore.Root, rootIDs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	volumedrivers.RegisterPluginGetter(daemon.PluginStore)

	if !volumedrivers.Register(volumesDriver, volumesDriver.Name()) ***REMOVED***
		return nil, errors.New("local volume driver could not be registered")
	***REMOVED***
	return store.New(daemon.configStore.Root)
***REMOVED***

// IsShuttingDown tells whether the daemon is shutting down or not
func (daemon *Daemon) IsShuttingDown() bool ***REMOVED***
	return daemon.shutdown
***REMOVED***

// initDiscovery initializes the discovery watcher for this daemon.
func (daemon *Daemon) initDiscovery(conf *config.Config) error ***REMOVED***
	advertise, err := config.ParseClusterAdvertiseSettings(conf.ClusterStore, conf.ClusterAdvertise)
	if err != nil ***REMOVED***
		if err == discovery.ErrDiscoveryDisabled ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	conf.ClusterAdvertise = advertise
	discoveryWatcher, err := discovery.Init(conf.ClusterStore, conf.ClusterAdvertise, conf.ClusterOpts)
	if err != nil ***REMOVED***
		return fmt.Errorf("discovery initialization failed (%v)", err)
	***REMOVED***

	daemon.discoveryWatcher = discoveryWatcher
	return nil
***REMOVED***

func isBridgeNetworkDisabled(conf *config.Config) bool ***REMOVED***
	return conf.BridgeConfig.Iface == config.DisableNetworkBridge
***REMOVED***

func (daemon *Daemon) networkOptions(dconfig *config.Config, pg plugingetter.PluginGetter, activeSandboxes map[string]interface***REMOVED******REMOVED***) ([]nwconfig.Option, error) ***REMOVED***
	options := []nwconfig.Option***REMOVED******REMOVED***
	if dconfig == nil ***REMOVED***
		return options, nil
	***REMOVED***

	options = append(options, nwconfig.OptionExperimental(dconfig.Experimental))
	options = append(options, nwconfig.OptionDataDir(dconfig.Root))
	options = append(options, nwconfig.OptionExecRoot(dconfig.GetExecRoot()))

	dd := runconfig.DefaultDaemonNetworkMode()
	dn := runconfig.DefaultDaemonNetworkMode().NetworkName()
	options = append(options, nwconfig.OptionDefaultDriver(string(dd)))
	options = append(options, nwconfig.OptionDefaultNetwork(dn))

	if strings.TrimSpace(dconfig.ClusterStore) != "" ***REMOVED***
		kv := strings.Split(dconfig.ClusterStore, "://")
		if len(kv) != 2 ***REMOVED***
			return nil, errors.New("kv store daemon config must be of the form KV-PROVIDER://KV-URL")
		***REMOVED***
		options = append(options, nwconfig.OptionKVProvider(kv[0]))
		options = append(options, nwconfig.OptionKVProviderURL(kv[1]))
	***REMOVED***
	if len(dconfig.ClusterOpts) > 0 ***REMOVED***
		options = append(options, nwconfig.OptionKVOpts(dconfig.ClusterOpts))
	***REMOVED***

	if daemon.discoveryWatcher != nil ***REMOVED***
		options = append(options, nwconfig.OptionDiscoveryWatcher(daemon.discoveryWatcher))
	***REMOVED***

	if dconfig.ClusterAdvertise != "" ***REMOVED***
		options = append(options, nwconfig.OptionDiscoveryAddress(dconfig.ClusterAdvertise))
	***REMOVED***

	options = append(options, nwconfig.OptionLabels(dconfig.Labels))
	options = append(options, driverOptions(dconfig)...)

	if daemon.configStore != nil && daemon.configStore.LiveRestoreEnabled && len(activeSandboxes) != 0 ***REMOVED***
		options = append(options, nwconfig.OptionActiveSandboxes(activeSandboxes))
	***REMOVED***

	if pg != nil ***REMOVED***
		options = append(options, nwconfig.OptionPluginGetter(pg))
	***REMOVED***

	options = append(options, nwconfig.OptionNetworkControlPlaneMTU(dconfig.NetworkControlPlaneMTU))

	return options, nil
***REMOVED***

// GetCluster returns the cluster
func (daemon *Daemon) GetCluster() Cluster ***REMOVED***
	return daemon.cluster
***REMOVED***

// SetCluster sets the cluster
func (daemon *Daemon) SetCluster(cluster Cluster) ***REMOVED***
	daemon.cluster = cluster
***REMOVED***

func (daemon *Daemon) pluginShutdown() ***REMOVED***
	manager := daemon.pluginManager
	// Check for a valid manager object. In error conditions, daemon init can fail
	// and shutdown called, before plugin manager is initialized.
	if manager != nil ***REMOVED***
		manager.Shutdown()
	***REMOVED***
***REMOVED***

// PluginManager returns current pluginManager associated with the daemon
func (daemon *Daemon) PluginManager() *plugin.Manager ***REMOVED*** // set up before daemon to avoid this method
	return daemon.pluginManager
***REMOVED***

// PluginGetter returns current pluginStore associated with the daemon
func (daemon *Daemon) PluginGetter() *plugin.Store ***REMOVED***
	return daemon.PluginStore
***REMOVED***

// CreateDaemonRoot creates the root for the daemon
func CreateDaemonRoot(config *config.Config) error ***REMOVED***
	// get the canonical path to the Docker root directory
	var realRoot string
	if _, err := os.Stat(config.Root); err != nil && os.IsNotExist(err) ***REMOVED***
		realRoot = config.Root
	***REMOVED*** else ***REMOVED***
		realRoot, err = getRealPath(config.Root)
		if err != nil ***REMOVED***
			return fmt.Errorf("Unable to get the full path to root (%s): %s", config.Root, err)
		***REMOVED***
	***REMOVED***

	idMappings, err := setupRemappedRoot(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return setupDaemonRoot(config, realRoot, idMappings.RootPair())
***REMOVED***

// checkpointAndSave grabs a container lock to safely call container.CheckpointTo
func (daemon *Daemon) checkpointAndSave(container *container.Container) error ***REMOVED***
	container.Lock()
	defer container.Unlock()
	if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
		return fmt.Errorf("Error saving container state: %v", err)
	***REMOVED***
	return nil
***REMOVED***

// because the CLI sends a -1 when it wants to unset the swappiness value
// we need to clear it on the server side
func fixMemorySwappiness(resources *containertypes.Resources) ***REMOVED***
	if resources.MemorySwappiness != nil && *resources.MemorySwappiness == -1 ***REMOVED***
		resources.MemorySwappiness = nil
	***REMOVED***
***REMOVED***

// GetAttachmentStore returns current attachment store associated with the daemon
func (daemon *Daemon) GetAttachmentStore() *network.AttachmentStore ***REMOVED***
	return &daemon.attachmentStore
***REMOVED***
