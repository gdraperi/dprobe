package daemon

import (
	"fmt"
	"net"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	clustertypes "github.com/docker/docker/daemon/cluster/provider"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/runconfig"
	"github.com/docker/libnetwork"
	lncluster "github.com/docker/libnetwork/cluster"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/ipamapi"
	networktypes "github.com/docker/libnetwork/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// NetworkControllerEnabled checks if the networking stack is enabled.
// This feature depends on OS primitives and it's disabled in systems like Windows.
func (daemon *Daemon) NetworkControllerEnabled() bool ***REMOVED***
	return daemon.netController != nil
***REMOVED***

// FindNetwork returns a network based on:
// 1. Full ID
// 2. Full Name
// 3. Partial ID
// as long as there is no ambiguity
func (daemon *Daemon) FindNetwork(term string) (libnetwork.Network, error) ***REMOVED***
	listByFullName := []libnetwork.Network***REMOVED******REMOVED***
	listByPartialID := []libnetwork.Network***REMOVED******REMOVED***
	for _, nw := range daemon.GetNetworks() ***REMOVED***
		if nw.ID() == term ***REMOVED***
			return nw, nil
		***REMOVED***
		if nw.Name() == term ***REMOVED***
			listByFullName = append(listByFullName, nw)
		***REMOVED***
		if strings.HasPrefix(nw.ID(), term) ***REMOVED***
			listByPartialID = append(listByPartialID, nw)
		***REMOVED***
	***REMOVED***
	switch ***REMOVED***
	case len(listByFullName) == 1:
		return listByFullName[0], nil
	case len(listByFullName) > 1:
		return nil, errdefs.InvalidParameter(errors.Errorf("network %s is ambiguous (%d matches found on name)", term, len(listByFullName)))
	case len(listByPartialID) == 1:
		return listByPartialID[0], nil
	case len(listByPartialID) > 1:
		return nil, errdefs.InvalidParameter(errors.Errorf("network %s is ambiguous (%d matches found based on ID prefix)", term, len(listByPartialID)))
	***REMOVED***

	// Be very careful to change the error type here, the
	// libnetwork.ErrNoSuchNetwork error is used by the controller
	// to retry the creation of the network as managed through the swarm manager
	return nil, errdefs.NotFound(libnetwork.ErrNoSuchNetwork(term))
***REMOVED***

// GetNetworkByID function returns a network whose ID matches the given ID.
// It fails with an error if no matching network is found.
func (daemon *Daemon) GetNetworkByID(id string) (libnetwork.Network, error) ***REMOVED***
	c := daemon.netController
	if c == nil ***REMOVED***
		return nil, libnetwork.ErrNoSuchNetwork(id)
	***REMOVED***
	return c.NetworkByID(id)
***REMOVED***

// GetNetworkByName function returns a network for a given network name.
// If no network name is given, the default network is returned.
func (daemon *Daemon) GetNetworkByName(name string) (libnetwork.Network, error) ***REMOVED***
	c := daemon.netController
	if c == nil ***REMOVED***
		return nil, libnetwork.ErrNoSuchNetwork(name)
	***REMOVED***
	if name == "" ***REMOVED***
		name = c.Config().Daemon.DefaultNetwork
	***REMOVED***
	return c.NetworkByName(name)
***REMOVED***

// GetNetworksByIDPrefix returns a list of networks whose ID partially matches zero or more networks
func (daemon *Daemon) GetNetworksByIDPrefix(partialID string) []libnetwork.Network ***REMOVED***
	c := daemon.netController
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	list := []libnetwork.Network***REMOVED******REMOVED***
	l := func(nw libnetwork.Network) bool ***REMOVED***
		if strings.HasPrefix(nw.ID(), partialID) ***REMOVED***
			list = append(list, nw)
		***REMOVED***
		return false
	***REMOVED***
	c.WalkNetworks(l)

	return list
***REMOVED***

// getAllNetworks returns a list containing all networks
func (daemon *Daemon) getAllNetworks() []libnetwork.Network ***REMOVED***
	c := daemon.netController
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	return c.Networks()
***REMOVED***

type ingressJob struct ***REMOVED***
	create  *clustertypes.NetworkCreateRequest
	ip      net.IP
	jobDone chan struct***REMOVED******REMOVED***
***REMOVED***

var (
	ingressWorkerOnce  sync.Once
	ingressJobsChannel chan *ingressJob
	ingressID          string
)

func (daemon *Daemon) startIngressWorker() ***REMOVED***
	ingressJobsChannel = make(chan *ingressJob, 100)
	go func() ***REMOVED***
		// nolint: gosimple
		for ***REMOVED***
			select ***REMOVED***
			case r := <-ingressJobsChannel:
				if r.create != nil ***REMOVED***
					daemon.setupIngress(r.create, r.ip, ingressID)
					ingressID = r.create.ID
				***REMOVED*** else ***REMOVED***
					daemon.releaseIngress(ingressID)
					ingressID = ""
				***REMOVED***
				close(r.jobDone)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

// enqueueIngressJob adds a ingress add/rm request to the worker queue.
// It guarantees the worker is started.
func (daemon *Daemon) enqueueIngressJob(job *ingressJob) ***REMOVED***
	ingressWorkerOnce.Do(daemon.startIngressWorker)
	ingressJobsChannel <- job
***REMOVED***

// SetupIngress setups ingress networking.
// The function returns a channel which will signal the caller when the programming is completed.
func (daemon *Daemon) SetupIngress(create clustertypes.NetworkCreateRequest, nodeIP string) (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***
	ip, _, err := net.ParseCIDR(nodeIP)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	done := make(chan struct***REMOVED******REMOVED***)
	daemon.enqueueIngressJob(&ingressJob***REMOVED***&create, ip, done***REMOVED***)
	return done, nil
***REMOVED***

// ReleaseIngress releases the ingress networking.
// The function returns a channel which will signal the caller when the programming is completed.
func (daemon *Daemon) ReleaseIngress() (<-chan struct***REMOVED******REMOVED***, error) ***REMOVED***
	done := make(chan struct***REMOVED******REMOVED***)
	daemon.enqueueIngressJob(&ingressJob***REMOVED***nil, nil, done***REMOVED***)
	return done, nil
***REMOVED***

func (daemon *Daemon) setupIngress(create *clustertypes.NetworkCreateRequest, ip net.IP, staleID string) ***REMOVED***
	controller := daemon.netController
	controller.AgentInitWait()

	if staleID != "" && staleID != create.ID ***REMOVED***
		daemon.releaseIngress(staleID)
	***REMOVED***

	if _, err := daemon.createNetwork(create.NetworkCreateRequest, create.ID, true); err != nil ***REMOVED***
		// If it is any other error other than already
		// exists error log error and return.
		if _, ok := err.(libnetwork.NetworkNameError); !ok ***REMOVED***
			logrus.Errorf("Failed creating ingress network: %v", err)
			return
		***REMOVED***
		// Otherwise continue down the call to create or recreate sandbox.
	***REMOVED***

	_, err := daemon.GetNetworkByID(create.ID)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed getting ingress network by id after creating: %v", err)
	***REMOVED***
***REMOVED***

func (daemon *Daemon) releaseIngress(id string) ***REMOVED***
	controller := daemon.netController

	if id == "" ***REMOVED***
		return
	***REMOVED***

	n, err := controller.NetworkByID(id)
	if err != nil ***REMOVED***
		logrus.Errorf("failed to retrieve ingress network %s: %v", id, err)
		return
	***REMOVED***

	if err := n.Delete(); err != nil ***REMOVED***
		logrus.Errorf("Failed to delete ingress network %s: %v", n.ID(), err)
		return
	***REMOVED***
***REMOVED***

// SetNetworkBootstrapKeys sets the bootstrap keys.
func (daemon *Daemon) SetNetworkBootstrapKeys(keys []*networktypes.EncryptionKey) error ***REMOVED***
	err := daemon.netController.SetKeys(keys)
	if err == nil ***REMOVED***
		// Upon successful key setting dispatch the keys available event
		daemon.cluster.SendClusterEvent(lncluster.EventNetworkKeysAvailable)
	***REMOVED***
	return err
***REMOVED***

// UpdateAttachment notifies the attacher about the attachment config.
func (daemon *Daemon) UpdateAttachment(networkName, networkID, containerID string, config *network.NetworkingConfig) error ***REMOVED***
	if daemon.clusterProvider == nil ***REMOVED***
		return fmt.Errorf("cluster provider is not initialized")
	***REMOVED***

	if err := daemon.clusterProvider.UpdateAttachment(networkName, containerID, config); err != nil ***REMOVED***
		return daemon.clusterProvider.UpdateAttachment(networkID, containerID, config)
	***REMOVED***

	return nil
***REMOVED***

// WaitForDetachment makes the cluster manager wait for detachment of
// the container from the network.
func (daemon *Daemon) WaitForDetachment(ctx context.Context, networkName, networkID, taskID, containerID string) error ***REMOVED***
	if daemon.clusterProvider == nil ***REMOVED***
		return fmt.Errorf("cluster provider is not initialized")
	***REMOVED***

	return daemon.clusterProvider.WaitForDetachment(ctx, networkName, networkID, taskID, containerID)
***REMOVED***

// CreateManagedNetwork creates an agent network.
func (daemon *Daemon) CreateManagedNetwork(create clustertypes.NetworkCreateRequest) error ***REMOVED***
	_, err := daemon.createNetwork(create.NetworkCreateRequest, create.ID, true)
	return err
***REMOVED***

// CreateNetwork creates a network with the given name, driver and other optional parameters
func (daemon *Daemon) CreateNetwork(create types.NetworkCreateRequest) (*types.NetworkCreateResponse, error) ***REMOVED***
	resp, err := daemon.createNetwork(create, "", false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp, err
***REMOVED***

func (daemon *Daemon) createNetwork(create types.NetworkCreateRequest, id string, agent bool) (*types.NetworkCreateResponse, error) ***REMOVED***
	if runconfig.IsPreDefinedNetwork(create.Name) && !agent ***REMOVED***
		err := fmt.Errorf("%s is a pre-defined network and cannot be created", create.Name)
		return nil, errdefs.Forbidden(err)
	***REMOVED***

	var warning string
	nw, err := daemon.GetNetworkByName(create.Name)
	if err != nil ***REMOVED***
		if _, ok := err.(libnetwork.ErrNoSuchNetwork); !ok ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if nw != nil ***REMOVED***
		// check if user defined CheckDuplicate, if set true, return err
		// otherwise prepare a warning message
		if create.CheckDuplicate ***REMOVED***
			if !agent || nw.Info().Dynamic() ***REMOVED***
				return nil, libnetwork.NetworkNameError(create.Name)
			***REMOVED***
		***REMOVED***
		warning = fmt.Sprintf("Network with name %s (id : %s) already exists", nw.Name(), nw.ID())
	***REMOVED***

	c := daemon.netController
	driver := create.Driver
	if driver == "" ***REMOVED***
		driver = c.Config().Daemon.DefaultDriver
	***REMOVED***

	nwOptions := []libnetwork.NetworkOption***REMOVED***
		libnetwork.NetworkOptionEnableIPv6(create.EnableIPv6),
		libnetwork.NetworkOptionDriverOpts(create.Options),
		libnetwork.NetworkOptionLabels(create.Labels),
		libnetwork.NetworkOptionAttachable(create.Attachable),
		libnetwork.NetworkOptionIngress(create.Ingress),
		libnetwork.NetworkOptionScope(create.Scope),
	***REMOVED***

	if create.ConfigOnly ***REMOVED***
		nwOptions = append(nwOptions, libnetwork.NetworkOptionConfigOnly())
	***REMOVED***

	if create.IPAM != nil ***REMOVED***
		ipam := create.IPAM
		v4Conf, v6Conf, err := getIpamConfig(ipam.Config)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		nwOptions = append(nwOptions, libnetwork.NetworkOptionIpam(ipam.Driver, "", v4Conf, v6Conf, ipam.Options))
	***REMOVED***

	if create.Internal ***REMOVED***
		nwOptions = append(nwOptions, libnetwork.NetworkOptionInternalNetwork())
	***REMOVED***
	if agent ***REMOVED***
		nwOptions = append(nwOptions, libnetwork.NetworkOptionDynamic())
		nwOptions = append(nwOptions, libnetwork.NetworkOptionPersist(false))
	***REMOVED***

	if create.ConfigFrom != nil ***REMOVED***
		nwOptions = append(nwOptions, libnetwork.NetworkOptionConfigFrom(create.ConfigFrom.Network))
	***REMOVED***

	if agent && driver == "overlay" && (create.Ingress || runtime.GOOS == "windows") ***REMOVED***
		nodeIP, exists := daemon.GetAttachmentStore().GetIPForNetwork(id)
		if !exists ***REMOVED***
			return nil, fmt.Errorf("Failed to find a load balancer IP to use for network: %v", id)
		***REMOVED***

		nwOptions = append(nwOptions, libnetwork.NetworkOptionLBEndpoint(nodeIP))
	***REMOVED***

	n, err := c.NewNetwork(driver, create.Name, id, nwOptions...)
	if err != nil ***REMOVED***
		if _, ok := err.(libnetwork.ErrDataStoreNotInitialized); ok ***REMOVED***
			// nolint: golint
			return nil, errors.New("This node is not a swarm manager. Use \"docker swarm init\" or \"docker swarm join\" to connect this node to swarm and try again.")
		***REMOVED***
		return nil, err
	***REMOVED***

	daemon.pluginRefCount(driver, driverapi.NetworkPluginEndpointType, plugingetter.Acquire)
	if create.IPAM != nil ***REMOVED***
		daemon.pluginRefCount(create.IPAM.Driver, ipamapi.PluginEndpointType, plugingetter.Acquire)
	***REMOVED***
	daemon.LogNetworkEvent(n, "create")

	return &types.NetworkCreateResponse***REMOVED***
		ID:      n.ID(),
		Warning: warning,
	***REMOVED***, nil
***REMOVED***

func (daemon *Daemon) pluginRefCount(driver, capability string, mode int) ***REMOVED***
	var builtinDrivers []string

	if capability == driverapi.NetworkPluginEndpointType ***REMOVED***
		builtinDrivers = daemon.netController.BuiltinDrivers()
	***REMOVED*** else if capability == ipamapi.PluginEndpointType ***REMOVED***
		builtinDrivers = daemon.netController.BuiltinIPAMDrivers()
	***REMOVED***

	for _, d := range builtinDrivers ***REMOVED***
		if d == driver ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if daemon.PluginStore != nil ***REMOVED***
		_, err := daemon.PluginStore.Get(driver, capability, mode)
		if err != nil ***REMOVED***
			logrus.WithError(err).WithFields(logrus.Fields***REMOVED***"mode": mode, "driver": driver***REMOVED***).Error("Error handling plugin refcount operation")
		***REMOVED***
	***REMOVED***
***REMOVED***

func getIpamConfig(data []network.IPAMConfig) ([]*libnetwork.IpamConf, []*libnetwork.IpamConf, error) ***REMOVED***
	ipamV4Cfg := []*libnetwork.IpamConf***REMOVED******REMOVED***
	ipamV6Cfg := []*libnetwork.IpamConf***REMOVED******REMOVED***
	for _, d := range data ***REMOVED***
		iCfg := libnetwork.IpamConf***REMOVED******REMOVED***
		iCfg.PreferredPool = d.Subnet
		iCfg.SubPool = d.IPRange
		iCfg.Gateway = d.Gateway
		iCfg.AuxAddresses = d.AuxAddress
		ip, _, err := net.ParseCIDR(d.Subnet)
		if err != nil ***REMOVED***
			return nil, nil, fmt.Errorf("Invalid subnet %s : %v", d.Subnet, err)
		***REMOVED***
		if ip.To4() != nil ***REMOVED***
			ipamV4Cfg = append(ipamV4Cfg, &iCfg)
		***REMOVED*** else ***REMOVED***
			ipamV6Cfg = append(ipamV6Cfg, &iCfg)
		***REMOVED***
	***REMOVED***
	return ipamV4Cfg, ipamV6Cfg, nil
***REMOVED***

// UpdateContainerServiceConfig updates a service configuration.
func (daemon *Daemon) UpdateContainerServiceConfig(containerName string, serviceConfig *clustertypes.ServiceConfig) error ***REMOVED***
	container, err := daemon.GetContainer(containerName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	container.NetworkSettings.Service = serviceConfig
	return nil
***REMOVED***

// ConnectContainerToNetwork connects the given container to the given
// network. If either cannot be found, an err is returned. If the
// network cannot be set up, an err is returned.
func (daemon *Daemon) ConnectContainerToNetwork(containerName, networkName string, endpointConfig *network.EndpointSettings) error ***REMOVED***
	container, err := daemon.GetContainer(containerName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return daemon.ConnectToNetwork(container, networkName, endpointConfig)
***REMOVED***

// DisconnectContainerFromNetwork disconnects the given container from
// the given network. If either cannot be found, an err is returned.
func (daemon *Daemon) DisconnectContainerFromNetwork(containerName string, networkName string, force bool) error ***REMOVED***
	container, err := daemon.GetContainer(containerName)
	if err != nil ***REMOVED***
		if force ***REMOVED***
			return daemon.ForceEndpointDelete(containerName, networkName)
		***REMOVED***
		return err
	***REMOVED***
	return daemon.DisconnectFromNetwork(container, networkName, force)
***REMOVED***

// GetNetworkDriverList returns the list of plugins drivers
// registered for network.
func (daemon *Daemon) GetNetworkDriverList() []string ***REMOVED***
	if !daemon.NetworkControllerEnabled() ***REMOVED***
		return nil
	***REMOVED***

	pluginList := daemon.netController.BuiltinDrivers()

	managedPlugins := daemon.PluginStore.GetAllManagedPluginsByCap(driverapi.NetworkPluginEndpointType)

	for _, plugin := range managedPlugins ***REMOVED***
		pluginList = append(pluginList, plugin.Name())
	***REMOVED***

	pluginMap := make(map[string]bool)
	for _, plugin := range pluginList ***REMOVED***
		pluginMap[plugin] = true
	***REMOVED***

	networks := daemon.netController.Networks()

	for _, network := range networks ***REMOVED***
		if !pluginMap[network.Type()] ***REMOVED***
			pluginList = append(pluginList, network.Type())
			pluginMap[network.Type()] = true
		***REMOVED***
	***REMOVED***

	sort.Strings(pluginList)

	return pluginList
***REMOVED***

// DeleteManagedNetwork deletes an agent network.
// The requirement of networkID is enforced.
func (daemon *Daemon) DeleteManagedNetwork(networkID string) error ***REMOVED***
	n, err := daemon.GetNetworkByID(networkID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return daemon.deleteNetwork(n, true)
***REMOVED***

// DeleteNetwork destroys a network unless it's one of docker's predefined networks.
func (daemon *Daemon) DeleteNetwork(networkID string) error ***REMOVED***
	n, err := daemon.GetNetworkByID(networkID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return daemon.deleteNetwork(n, false)
***REMOVED***

func (daemon *Daemon) deleteLoadBalancerSandbox(n libnetwork.Network) ***REMOVED***
	controller := daemon.netController

	//The only endpoint left should be the LB endpoint (nw.Name() + "-endpoint")
	endpoints := n.Endpoints()
	if len(endpoints) == 1 ***REMOVED***
		sandboxName := n.Name() + "-sbox"

		info := endpoints[0].Info()
		if info != nil ***REMOVED***
			sb := info.Sandbox()
			if sb != nil ***REMOVED***
				if err := sb.DisableService(); err != nil ***REMOVED***
					logrus.Warnf("Failed to disable service on sandbox %s: %v", sandboxName, err)
					//Ignore error and attempt to delete the load balancer endpoint
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err := endpoints[0].Delete(true); err != nil ***REMOVED***
			logrus.Warnf("Failed to delete endpoint %s (%s) in %s: %v", endpoints[0].Name(), endpoints[0].ID(), sandboxName, err)
			//Ignore error and attempt to delete the sandbox.
		***REMOVED***

		if err := controller.SandboxDestroy(sandboxName); err != nil ***REMOVED***
			logrus.Warnf("Failed to delete %s sandbox: %v", sandboxName, err)
			//Ignore error and attempt to delete the network.
		***REMOVED***
	***REMOVED***
***REMOVED***

func (daemon *Daemon) deleteNetwork(nw libnetwork.Network, dynamic bool) error ***REMOVED***
	if runconfig.IsPreDefinedNetwork(nw.Name()) && !dynamic ***REMOVED***
		err := fmt.Errorf("%s is a pre-defined network and cannot be removed", nw.Name())
		return errdefs.Forbidden(err)
	***REMOVED***

	if dynamic && !nw.Info().Dynamic() ***REMOVED***
		if runconfig.IsPreDefinedNetwork(nw.Name()) ***REMOVED***
			// Predefined networks now support swarm services. Make this
			// a no-op when cluster requests to remove the predefined network.
			return nil
		***REMOVED***
		err := fmt.Errorf("%s is not a dynamic network", nw.Name())
		return errdefs.Forbidden(err)
	***REMOVED***

	if err := nw.Delete(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// If this is not a configuration only network, we need to
	// update the corresponding remote drivers' reference counts
	if !nw.Info().ConfigOnly() ***REMOVED***
		daemon.pluginRefCount(nw.Type(), driverapi.NetworkPluginEndpointType, plugingetter.Release)
		ipamType, _, _, _ := nw.Info().IpamConfig()
		daemon.pluginRefCount(ipamType, ipamapi.PluginEndpointType, plugingetter.Release)
		daemon.LogNetworkEvent(nw, "destroy")
	***REMOVED***

	return nil
***REMOVED***

// GetNetworks returns a list of all networks
func (daemon *Daemon) GetNetworks() []libnetwork.Network ***REMOVED***
	return daemon.getAllNetworks()
***REMOVED***

// clearAttachableNetworks removes the attachable networks
// after disconnecting any connected container
func (daemon *Daemon) clearAttachableNetworks() ***REMOVED***
	for _, n := range daemon.GetNetworks() ***REMOVED***
		if !n.Info().Attachable() ***REMOVED***
			continue
		***REMOVED***
		for _, ep := range n.Endpoints() ***REMOVED***
			epInfo := ep.Info()
			if epInfo == nil ***REMOVED***
				continue
			***REMOVED***
			sb := epInfo.Sandbox()
			if sb == nil ***REMOVED***
				continue
			***REMOVED***
			containerID := sb.ContainerID()
			if err := daemon.DisconnectContainerFromNetwork(containerID, n.ID(), true); err != nil ***REMOVED***
				logrus.Warnf("Failed to disconnect container %s from swarm network %s on cluster leave: %v",
					containerID, n.Name(), err)
			***REMOVED***
		***REMOVED***
		if err := daemon.DeleteManagedNetwork(n.ID()); err != nil ***REMOVED***
			logrus.Warnf("Failed to remove swarm network %s on cluster leave: %v", n.Name(), err)
		***REMOVED***
	***REMOVED***
***REMOVED***
