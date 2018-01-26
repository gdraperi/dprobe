package daemon

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/network"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/runconfig"
	"github.com/docker/go-connections/nat"
	"github.com/docker/libnetwork"
	netconst "github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/options"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

var (
	// ErrRootFSReadOnly is returned when a container
	// rootfs is marked readonly.
	ErrRootFSReadOnly = errors.New("container rootfs is marked read-only")
	getPortMapInfo    = container.GetSandboxPortMapInfo
)

func (daemon *Daemon) getDNSSearchSettings(container *container.Container) []string ***REMOVED***
	if len(container.HostConfig.DNSSearch) > 0 ***REMOVED***
		return container.HostConfig.DNSSearch
	***REMOVED***

	if len(daemon.configStore.DNSSearch) > 0 ***REMOVED***
		return daemon.configStore.DNSSearch
	***REMOVED***

	return nil
***REMOVED***

func (daemon *Daemon) buildSandboxOptions(container *container.Container) ([]libnetwork.SandboxOption, error) ***REMOVED***
	var (
		sboxOptions []libnetwork.SandboxOption
		err         error
		dns         []string
		dnsOptions  []string
		bindings    = make(nat.PortMap)
		pbList      []types.PortBinding
		exposeList  []types.TransportPort
	)

	defaultNetName := runconfig.DefaultDaemonNetworkMode().NetworkName()
	sboxOptions = append(sboxOptions, libnetwork.OptionHostname(container.Config.Hostname),
		libnetwork.OptionDomainname(container.Config.Domainname))

	if container.HostConfig.NetworkMode.IsHost() ***REMOVED***
		sboxOptions = append(sboxOptions, libnetwork.OptionUseDefaultSandbox())
		if len(container.HostConfig.ExtraHosts) == 0 ***REMOVED***
			sboxOptions = append(sboxOptions, libnetwork.OptionOriginHostsPath("/etc/hosts"))
		***REMOVED***
		if len(container.HostConfig.DNS) == 0 && len(daemon.configStore.DNS) == 0 &&
			len(container.HostConfig.DNSSearch) == 0 && len(daemon.configStore.DNSSearch) == 0 &&
			len(container.HostConfig.DNSOptions) == 0 && len(daemon.configStore.DNSOptions) == 0 ***REMOVED***
			sboxOptions = append(sboxOptions, libnetwork.OptionOriginResolvConfPath("/etc/resolv.conf"))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// OptionUseExternalKey is mandatory for userns support.
		// But optional for non-userns support
		sboxOptions = append(sboxOptions, libnetwork.OptionUseExternalKey())
	***REMOVED***

	if err = setupPathsAndSandboxOptions(container, &sboxOptions); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(container.HostConfig.DNS) > 0 ***REMOVED***
		dns = container.HostConfig.DNS
	***REMOVED*** else if len(daemon.configStore.DNS) > 0 ***REMOVED***
		dns = daemon.configStore.DNS
	***REMOVED***

	for _, d := range dns ***REMOVED***
		sboxOptions = append(sboxOptions, libnetwork.OptionDNS(d))
	***REMOVED***

	dnsSearch := daemon.getDNSSearchSettings(container)

	for _, ds := range dnsSearch ***REMOVED***
		sboxOptions = append(sboxOptions, libnetwork.OptionDNSSearch(ds))
	***REMOVED***

	if len(container.HostConfig.DNSOptions) > 0 ***REMOVED***
		dnsOptions = container.HostConfig.DNSOptions
	***REMOVED*** else if len(daemon.configStore.DNSOptions) > 0 ***REMOVED***
		dnsOptions = daemon.configStore.DNSOptions
	***REMOVED***

	for _, ds := range dnsOptions ***REMOVED***
		sboxOptions = append(sboxOptions, libnetwork.OptionDNSOptions(ds))
	***REMOVED***

	if container.NetworkSettings.SecondaryIPAddresses != nil ***REMOVED***
		name := container.Config.Hostname
		if container.Config.Domainname != "" ***REMOVED***
			name = name + "." + container.Config.Domainname
		***REMOVED***

		for _, a := range container.NetworkSettings.SecondaryIPAddresses ***REMOVED***
			sboxOptions = append(sboxOptions, libnetwork.OptionExtraHost(name, a.Addr))
		***REMOVED***
	***REMOVED***

	for _, extraHost := range container.HostConfig.ExtraHosts ***REMOVED***
		// allow IPv6 addresses in extra hosts; only split on first ":"
		if _, err := opts.ValidateExtraHost(extraHost); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		parts := strings.SplitN(extraHost, ":", 2)
		sboxOptions = append(sboxOptions, libnetwork.OptionExtraHost(parts[0], parts[1]))
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
				return nil, fmt.Errorf("Error parsing HostPort value(%s):%v", binding[i].HostPort, err)
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

	sboxOptions = append(sboxOptions,
		libnetwork.OptionPortMapping(pbList),
		libnetwork.OptionExposedPorts(exposeList))

	// Legacy Link feature is supported only for the default bridge network.
	// return if this call to build join options is not for default bridge network
	// Legacy Link is only supported by docker run --link
	bridgeSettings, ok := container.NetworkSettings.Networks[defaultNetName]
	if !ok || bridgeSettings.EndpointSettings == nil ***REMOVED***
		return sboxOptions, nil
	***REMOVED***

	if bridgeSettings.EndpointID == "" ***REMOVED***
		return sboxOptions, nil
	***REMOVED***

	var (
		childEndpoints, parentEndpoints []string
		cEndpointID                     string
	)

	children := daemon.children(container)
	for linkAlias, child := range children ***REMOVED***
		if !isLinkable(child) ***REMOVED***
			return nil, fmt.Errorf("Cannot link to %s, as it does not belong to the default network", child.Name)
		***REMOVED***
		_, alias := path.Split(linkAlias)
		// allow access to the linked container via the alias, real name, and container hostname
		aliasList := alias + " " + child.Config.Hostname
		// only add the name if alias isn't equal to the name
		if alias != child.Name[1:] ***REMOVED***
			aliasList = aliasList + " " + child.Name[1:]
		***REMOVED***
		sboxOptions = append(sboxOptions, libnetwork.OptionExtraHost(aliasList, child.NetworkSettings.Networks[defaultNetName].IPAddress))
		cEndpointID = child.NetworkSettings.Networks[defaultNetName].EndpointID
		if cEndpointID != "" ***REMOVED***
			childEndpoints = append(childEndpoints, cEndpointID)
		***REMOVED***
	***REMOVED***

	for alias, parent := range daemon.parents(container) ***REMOVED***
		if daemon.configStore.DisableBridge || !container.HostConfig.NetworkMode.IsPrivate() ***REMOVED***
			continue
		***REMOVED***

		_, alias = path.Split(alias)
		logrus.Debugf("Update /etc/hosts of %s for alias %s with ip %s", parent.ID, alias, bridgeSettings.IPAddress)
		sboxOptions = append(sboxOptions, libnetwork.OptionParentUpdate(
			parent.ID,
			alias,
			bridgeSettings.IPAddress,
		))
		if cEndpointID != "" ***REMOVED***
			parentEndpoints = append(parentEndpoints, cEndpointID)
		***REMOVED***
	***REMOVED***

	linkOptions := options.Generic***REMOVED***
		netlabel.GenericData: options.Generic***REMOVED***
			"ParentEndpoints": parentEndpoints,
			"ChildEndpoints":  childEndpoints,
		***REMOVED***,
	***REMOVED***

	sboxOptions = append(sboxOptions, libnetwork.OptionGeneric(linkOptions))
	return sboxOptions, nil
***REMOVED***

func (daemon *Daemon) updateNetworkSettings(container *container.Container, n libnetwork.Network, endpointConfig *networktypes.EndpointSettings) error ***REMOVED***
	if container.NetworkSettings == nil ***REMOVED***
		container.NetworkSettings = &network.Settings***REMOVED***Networks: make(map[string]*network.EndpointSettings)***REMOVED***
	***REMOVED***

	if !container.HostConfig.NetworkMode.IsHost() && containertypes.NetworkMode(n.Type()).IsHost() ***REMOVED***
		return runconfig.ErrConflictHostNetwork
	***REMOVED***

	for s, v := range container.NetworkSettings.Networks ***REMOVED***
		sn, err := daemon.FindNetwork(getNetworkID(s, v.EndpointSettings))
		if err != nil ***REMOVED***
			continue
		***REMOVED***

		if sn.Name() == n.Name() ***REMOVED***
			// If the network scope is swarm, then this
			// is an attachable network, which may not
			// be locally available previously.
			// So always update.
			if n.Info().Scope() == netconst.SwarmScope ***REMOVED***
				continue
			***REMOVED***
			// Avoid duplicate config
			return nil
		***REMOVED***
		if !containertypes.NetworkMode(sn.Type()).IsPrivate() ||
			!containertypes.NetworkMode(n.Type()).IsPrivate() ***REMOVED***
			return runconfig.ErrConflictSharedNetwork
		***REMOVED***
		if containertypes.NetworkMode(sn.Name()).IsNone() ||
			containertypes.NetworkMode(n.Name()).IsNone() ***REMOVED***
			return runconfig.ErrConflictNoNetwork
		***REMOVED***
	***REMOVED***

	container.NetworkSettings.Networks[n.Name()] = &network.EndpointSettings***REMOVED***
		EndpointSettings: endpointConfig,
	***REMOVED***

	return nil
***REMOVED***

func (daemon *Daemon) updateEndpointNetworkSettings(container *container.Container, n libnetwork.Network, ep libnetwork.Endpoint) error ***REMOVED***
	if err := container.BuildEndpointInfo(n, ep); err != nil ***REMOVED***
		return err
	***REMOVED***

	if container.HostConfig.NetworkMode == runconfig.DefaultDaemonNetworkMode() ***REMOVED***
		container.NetworkSettings.Bridge = daemon.configStore.BridgeConfig.Iface
	***REMOVED***

	return nil
***REMOVED***

// UpdateNetwork is used to update the container's network (e.g. when linked containers
// get removed/unlinked).
func (daemon *Daemon) updateNetwork(container *container.Container) error ***REMOVED***
	var (
		start = time.Now()
		ctrl  = daemon.netController
		sid   = container.NetworkSettings.SandboxID
	)

	sb, err := ctrl.SandboxByID(sid)
	if err != nil ***REMOVED***
		return fmt.Errorf("error locating sandbox id %s: %v", sid, err)
	***REMOVED***

	// Find if container is connected to the default bridge network
	var n libnetwork.Network
	for name, v := range container.NetworkSettings.Networks ***REMOVED***
		sn, err := daemon.FindNetwork(getNetworkID(name, v.EndpointSettings))
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		if sn.Name() == runconfig.DefaultDaemonNetworkMode().NetworkName() ***REMOVED***
			n = sn
			break
		***REMOVED***
	***REMOVED***

	if n == nil ***REMOVED***
		// Not connected to the default bridge network; Nothing to do
		return nil
	***REMOVED***

	options, err := daemon.buildSandboxOptions(container)
	if err != nil ***REMOVED***
		return fmt.Errorf("Update network failed: %v", err)
	***REMOVED***

	if err := sb.Refresh(options...); err != nil ***REMOVED***
		return fmt.Errorf("Update network failed: Failure in refresh sandbox %s: %v", sid, err)
	***REMOVED***

	networkActions.WithValues("update").UpdateSince(start)

	return nil
***REMOVED***

func (daemon *Daemon) findAndAttachNetwork(container *container.Container, idOrName string, epConfig *networktypes.EndpointSettings) (libnetwork.Network, *networktypes.NetworkingConfig, error) ***REMOVED***
	n, err := daemon.FindNetwork(getNetworkID(idOrName, epConfig))
	if err != nil ***REMOVED***
		// We should always be able to find the network for a
		// managed container.
		if container.Managed ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***

	// If we found a network and if it is not dynamically created
	// we should never attempt to attach to that network here.
	if n != nil ***REMOVED***
		if container.Managed || !n.Info().Dynamic() ***REMOVED***
			return n, nil, nil
		***REMOVED***
	***REMOVED***

	var addresses []string
	if epConfig != nil && epConfig.IPAMConfig != nil ***REMOVED***
		if epConfig.IPAMConfig.IPv4Address != "" ***REMOVED***
			addresses = append(addresses, epConfig.IPAMConfig.IPv4Address)
		***REMOVED***

		if epConfig.IPAMConfig.IPv6Address != "" ***REMOVED***
			addresses = append(addresses, epConfig.IPAMConfig.IPv6Address)
		***REMOVED***
	***REMOVED***

	var (
		config     *networktypes.NetworkingConfig
		retryCount int
	)

	for ***REMOVED***
		// In all other cases, attempt to attach to the network to
		// trigger attachment in the swarm cluster manager.
		if daemon.clusterProvider != nil ***REMOVED***
			var err error
			config, err = daemon.clusterProvider.AttachNetwork(getNetworkID(idOrName, epConfig), container.ID, addresses)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		***REMOVED***

		n, err = daemon.FindNetwork(getNetworkID(idOrName, epConfig))
		if err != nil ***REMOVED***
			if daemon.clusterProvider != nil ***REMOVED***
				if err := daemon.clusterProvider.DetachNetwork(getNetworkID(idOrName, epConfig), container.ID); err != nil ***REMOVED***
					logrus.Warnf("Could not rollback attachment for container %s to network %s: %v", container.ID, idOrName, err)
				***REMOVED***
			***REMOVED***

			// Retry network attach again if we failed to
			// find the network after successful
			// attachment because the only reason that
			// would happen is if some other container
			// attached to the swarm scope network went down
			// and removed the network while we were in
			// the process of attaching.
			if config != nil ***REMOVED***
				if _, ok := err.(libnetwork.ErrNoSuchNetwork); ok ***REMOVED***
					if retryCount >= 5 ***REMOVED***
						return nil, nil, fmt.Errorf("could not find network %s after successful attachment", idOrName)
					***REMOVED***
					retryCount++
					continue
				***REMOVED***
			***REMOVED***

			return nil, nil, err
		***REMOVED***

		break
	***REMOVED***

	// This container has attachment to a swarm scope
	// network. Update the container network settings accordingly.
	container.NetworkSettings.HasSwarmEndpoint = true
	return n, config, nil
***REMOVED***

// updateContainerNetworkSettings updates the network settings
func (daemon *Daemon) updateContainerNetworkSettings(container *container.Container, endpointsConfig map[string]*networktypes.EndpointSettings) ***REMOVED***
	var n libnetwork.Network

	mode := container.HostConfig.NetworkMode
	if container.Config.NetworkDisabled || mode.IsContainer() ***REMOVED***
		return
	***REMOVED***

	networkName := mode.NetworkName()
	if mode.IsDefault() ***REMOVED***
		networkName = daemon.netController.Config().Daemon.DefaultNetwork
	***REMOVED***

	if mode.IsUserDefined() ***REMOVED***
		var err error

		n, err = daemon.FindNetwork(networkName)
		if err == nil ***REMOVED***
			networkName = n.Name()
		***REMOVED***
	***REMOVED***

	if container.NetworkSettings == nil ***REMOVED***
		container.NetworkSettings = &network.Settings***REMOVED******REMOVED***
	***REMOVED***

	if len(endpointsConfig) > 0 ***REMOVED***
		if container.NetworkSettings.Networks == nil ***REMOVED***
			container.NetworkSettings.Networks = make(map[string]*network.EndpointSettings)
		***REMOVED***

		for name, epConfig := range endpointsConfig ***REMOVED***
			container.NetworkSettings.Networks[name] = &network.EndpointSettings***REMOVED***
				EndpointSettings: epConfig,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if container.NetworkSettings.Networks == nil ***REMOVED***
		container.NetworkSettings.Networks = make(map[string]*network.EndpointSettings)
		container.NetworkSettings.Networks[networkName] = &network.EndpointSettings***REMOVED***
			EndpointSettings: &networktypes.EndpointSettings***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED***

	// Convert any settings added by client in default name to
	// engine's default network name key
	if mode.IsDefault() ***REMOVED***
		if nConf, ok := container.NetworkSettings.Networks[mode.NetworkName()]; ok ***REMOVED***
			container.NetworkSettings.Networks[networkName] = nConf
			delete(container.NetworkSettings.Networks, mode.NetworkName())
		***REMOVED***
	***REMOVED***

	if !mode.IsUserDefined() ***REMOVED***
		return
	***REMOVED***
	// Make sure to internally store the per network endpoint config by network name
	if _, ok := container.NetworkSettings.Networks[networkName]; ok ***REMOVED***
		return
	***REMOVED***

	if n != nil ***REMOVED***
		if nwConfig, ok := container.NetworkSettings.Networks[n.ID()]; ok ***REMOVED***
			container.NetworkSettings.Networks[networkName] = nwConfig
			delete(container.NetworkSettings.Networks, n.ID())
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (daemon *Daemon) allocateNetwork(container *container.Container) error ***REMOVED***
	start := time.Now()
	controller := daemon.netController

	if daemon.netController == nil ***REMOVED***
		return nil
	***REMOVED***

	// Cleanup any stale sandbox left over due to ungraceful daemon shutdown
	if err := controller.SandboxDestroy(container.ID); err != nil ***REMOVED***
		logrus.Errorf("failed to cleanup up stale network sandbox for container %s", container.ID)
	***REMOVED***

	if container.Config.NetworkDisabled || container.HostConfig.NetworkMode.IsContainer() ***REMOVED***
		return nil
	***REMOVED***

	updateSettings := false

	if len(container.NetworkSettings.Networks) == 0 ***REMOVED***
		daemon.updateContainerNetworkSettings(container, nil)
		updateSettings = true
	***REMOVED***

	// always connect default network first since only default
	// network mode support link and we need do some setting
	// on sandbox initialize for link, but the sandbox only be initialized
	// on first network connecting.
	defaultNetName := runconfig.DefaultDaemonNetworkMode().NetworkName()
	if nConf, ok := container.NetworkSettings.Networks[defaultNetName]; ok ***REMOVED***
		cleanOperationalData(nConf)
		if err := daemon.connectToNetwork(container, defaultNetName, nConf.EndpointSettings, updateSettings); err != nil ***REMOVED***
			return err
		***REMOVED***

	***REMOVED***

	// the intermediate map is necessary because "connectToNetwork" modifies "container.NetworkSettings.Networks"
	networks := make(map[string]*network.EndpointSettings)
	for n, epConf := range container.NetworkSettings.Networks ***REMOVED***
		if n == defaultNetName ***REMOVED***
			continue
		***REMOVED***

		networks[n] = epConf
	***REMOVED***

	for netName, epConf := range networks ***REMOVED***
		cleanOperationalData(epConf)
		if err := daemon.connectToNetwork(container, netName, epConf.EndpointSettings, updateSettings); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// If the container is not to be connected to any network,
	// create its network sandbox now if not present
	if len(networks) == 0 ***REMOVED***
		if nil == daemon.getNetworkSandbox(container) ***REMOVED***
			options, err := daemon.buildSandboxOptions(container)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			sb, err := daemon.netController.NewSandbox(container.ID, options...)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			container.UpdateSandboxNetworkSettings(sb)
			defer func() ***REMOVED***
				if err != nil ***REMOVED***
					sb.Delete()
				***REMOVED***
			***REMOVED***()
		***REMOVED***

	***REMOVED***

	if _, err := container.WriteHostConfig(); err != nil ***REMOVED***
		return err
	***REMOVED***
	networkActions.WithValues("allocate").UpdateSince(start)
	return nil
***REMOVED***

func (daemon *Daemon) getNetworkSandbox(container *container.Container) libnetwork.Sandbox ***REMOVED***
	var sb libnetwork.Sandbox
	daemon.netController.WalkSandboxes(func(s libnetwork.Sandbox) bool ***REMOVED***
		if s.ContainerID() == container.ID ***REMOVED***
			sb = s
			return true
		***REMOVED***
		return false
	***REMOVED***)
	return sb
***REMOVED***

// hasUserDefinedIPAddress returns whether the passed endpoint configuration contains IP address configuration
func hasUserDefinedIPAddress(epConfig *networktypes.EndpointSettings) bool ***REMOVED***
	return epConfig != nil && epConfig.IPAMConfig != nil && (len(epConfig.IPAMConfig.IPv4Address) > 0 || len(epConfig.IPAMConfig.IPv6Address) > 0)
***REMOVED***

// User specified ip address is acceptable only for networks with user specified subnets.
func validateNetworkingConfig(n libnetwork.Network, epConfig *networktypes.EndpointSettings) error ***REMOVED***
	if n == nil || epConfig == nil ***REMOVED***
		return nil
	***REMOVED***
	if !hasUserDefinedIPAddress(epConfig) ***REMOVED***
		return nil
	***REMOVED***
	_, _, nwIPv4Configs, nwIPv6Configs := n.Info().IpamConfig()
	for _, s := range []struct ***REMOVED***
		ipConfigured  bool
		subnetConfigs []*libnetwork.IpamConf
	***REMOVED******REMOVED***
		***REMOVED***
			ipConfigured:  len(epConfig.IPAMConfig.IPv4Address) > 0,
			subnetConfigs: nwIPv4Configs,
		***REMOVED***,
		***REMOVED***
			ipConfigured:  len(epConfig.IPAMConfig.IPv6Address) > 0,
			subnetConfigs: nwIPv6Configs,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if s.ipConfigured ***REMOVED***
			foundSubnet := false
			for _, cfg := range s.subnetConfigs ***REMOVED***
				if len(cfg.PreferredPool) > 0 ***REMOVED***
					foundSubnet = true
					break
				***REMOVED***
			***REMOVED***
			if !foundSubnet ***REMOVED***
				return runconfig.ErrUnsupportedNetworkNoSubnetAndIP
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// cleanOperationalData resets the operational data from the passed endpoint settings
func cleanOperationalData(es *network.EndpointSettings) ***REMOVED***
	es.EndpointID = ""
	es.Gateway = ""
	es.IPAddress = ""
	es.IPPrefixLen = 0
	es.IPv6Gateway = ""
	es.GlobalIPv6Address = ""
	es.GlobalIPv6PrefixLen = 0
	es.MacAddress = ""
	if es.IPAMOperational ***REMOVED***
		es.IPAMConfig = nil
	***REMOVED***
***REMOVED***

func (daemon *Daemon) updateNetworkConfig(container *container.Container, n libnetwork.Network, endpointConfig *networktypes.EndpointSettings, updateSettings bool) error ***REMOVED***

	if !containertypes.NetworkMode(n.Name()).IsUserDefined() ***REMOVED***
		if hasUserDefinedIPAddress(endpointConfig) && !enableIPOnPredefinedNetwork() ***REMOVED***
			return runconfig.ErrUnsupportedNetworkAndIP
		***REMOVED***
		if endpointConfig != nil && len(endpointConfig.Aliases) > 0 && !container.EnableServiceDiscoveryOnDefaultNetwork() ***REMOVED***
			return runconfig.ErrUnsupportedNetworkAndAlias
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		addShortID := true
		shortID := stringid.TruncateID(container.ID)
		for _, alias := range endpointConfig.Aliases ***REMOVED***
			if alias == shortID ***REMOVED***
				addShortID = false
				break
			***REMOVED***
		***REMOVED***
		if addShortID ***REMOVED***
			endpointConfig.Aliases = append(endpointConfig.Aliases, shortID)
		***REMOVED***
	***REMOVED***

	if err := validateNetworkingConfig(n, endpointConfig); err != nil ***REMOVED***
		return err
	***REMOVED***

	if updateSettings ***REMOVED***
		if err := daemon.updateNetworkSettings(container, n, endpointConfig); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) connectToNetwork(container *container.Container, idOrName string, endpointConfig *networktypes.EndpointSettings, updateSettings bool) (err error) ***REMOVED***
	start := time.Now()
	if container.HostConfig.NetworkMode.IsContainer() ***REMOVED***
		return runconfig.ErrConflictSharedNetwork
	***REMOVED***
	if containertypes.NetworkMode(idOrName).IsBridge() &&
		daemon.configStore.DisableBridge ***REMOVED***
		container.Config.NetworkDisabled = true
		return nil
	***REMOVED***
	if endpointConfig == nil ***REMOVED***
		endpointConfig = &networktypes.EndpointSettings***REMOVED******REMOVED***
	***REMOVED***

	n, config, err := daemon.findAndAttachNetwork(container, idOrName, endpointConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n == nil ***REMOVED***
		return nil
	***REMOVED***

	var operIPAM bool
	if config != nil ***REMOVED***
		if epConfig, ok := config.EndpointsConfig[n.Name()]; ok ***REMOVED***
			if endpointConfig.IPAMConfig == nil ||
				(endpointConfig.IPAMConfig.IPv4Address == "" &&
					endpointConfig.IPAMConfig.IPv6Address == "" &&
					len(endpointConfig.IPAMConfig.LinkLocalIPs) == 0) ***REMOVED***
				operIPAM = true
			***REMOVED***

			// copy IPAMConfig and NetworkID from epConfig via AttachNetwork
			endpointConfig.IPAMConfig = epConfig.IPAMConfig
			endpointConfig.NetworkID = epConfig.NetworkID
		***REMOVED***
	***REMOVED***

	err = daemon.updateNetworkConfig(container, n, endpointConfig, updateSettings)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	controller := daemon.netController
	sb := daemon.getNetworkSandbox(container)
	createOptions, err := container.BuildCreateEndpointOptions(n, endpointConfig, sb, daemon.configStore.DNS)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	endpointName := strings.TrimPrefix(container.Name, "/")
	ep, err := n.CreateEndpoint(endpointName, createOptions...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if e := ep.Delete(false); e != nil ***REMOVED***
				logrus.Warnf("Could not rollback container connection to network %s", idOrName)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	container.NetworkSettings.Networks[n.Name()] = &network.EndpointSettings***REMOVED***
		EndpointSettings: endpointConfig,
		IPAMOperational:  operIPAM,
	***REMOVED***
	if _, ok := container.NetworkSettings.Networks[n.ID()]; ok ***REMOVED***
		delete(container.NetworkSettings.Networks, n.ID())
	***REMOVED***

	if err := daemon.updateEndpointNetworkSettings(container, n, ep); err != nil ***REMOVED***
		return err
	***REMOVED***

	if sb == nil ***REMOVED***
		options, err := daemon.buildSandboxOptions(container)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		sb, err = controller.NewSandbox(container.ID, options...)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		container.UpdateSandboxNetworkSettings(sb)
	***REMOVED***

	joinOptions, err := container.BuildJoinOptions(n)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := ep.Join(sb, joinOptions...); err != nil ***REMOVED***
		return err
	***REMOVED***

	if !container.Managed ***REMOVED***
		// add container name/alias to DNS
		if err := daemon.ActivateContainerServiceBinding(container.Name); err != nil ***REMOVED***
			return fmt.Errorf("Activate container service binding for %s failed: %v", container.Name, err)
		***REMOVED***
	***REMOVED***

	if err := container.UpdateJoinInfo(n, ep); err != nil ***REMOVED***
		return fmt.Errorf("Updating join info failed: %v", err)
	***REMOVED***

	container.NetworkSettings.Ports = getPortMapInfo(sb)

	daemon.LogNetworkEventWithAttributes(n, "connect", map[string]string***REMOVED***"container": container.ID***REMOVED***)
	networkActions.WithValues("connect").UpdateSince(start)
	return nil
***REMOVED***

// ForceEndpointDelete deletes an endpoint from a network forcefully
func (daemon *Daemon) ForceEndpointDelete(name string, networkName string) error ***REMOVED***
	n, err := daemon.FindNetwork(networkName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ep, err := n.EndpointByName(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ep.Delete(true)
***REMOVED***

func (daemon *Daemon) disconnectFromNetwork(container *container.Container, n libnetwork.Network, force bool) error ***REMOVED***
	var (
		ep   libnetwork.Endpoint
		sbox libnetwork.Sandbox
	)

	s := func(current libnetwork.Endpoint) bool ***REMOVED***
		epInfo := current.Info()
		if epInfo == nil ***REMOVED***
			return false
		***REMOVED***
		if sb := epInfo.Sandbox(); sb != nil ***REMOVED***
			if sb.ContainerID() == container.ID ***REMOVED***
				ep = current
				sbox = sb
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
	n.WalkEndpoints(s)

	if ep == nil && force ***REMOVED***
		epName := strings.TrimPrefix(container.Name, "/")
		ep, err := n.EndpointByName(epName)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return ep.Delete(force)
	***REMOVED***

	if ep == nil ***REMOVED***
		return fmt.Errorf("container %s is not connected to network %s", container.ID, n.Name())
	***REMOVED***

	if err := ep.Leave(sbox); err != nil ***REMOVED***
		return fmt.Errorf("container %s failed to leave network %s: %v", container.ID, n.Name(), err)
	***REMOVED***

	container.NetworkSettings.Ports = getPortMapInfo(sbox)

	if err := ep.Delete(false); err != nil ***REMOVED***
		return fmt.Errorf("endpoint delete failed for container %s on network %s: %v", container.ID, n.Name(), err)
	***REMOVED***

	delete(container.NetworkSettings.Networks, n.Name())

	daemon.tryDetachContainerFromClusterNetwork(n, container)

	return nil
***REMOVED***

func (daemon *Daemon) tryDetachContainerFromClusterNetwork(network libnetwork.Network, container *container.Container) ***REMOVED***
	if daemon.clusterProvider != nil && network.Info().Dynamic() && !container.Managed ***REMOVED***
		if err := daemon.clusterProvider.DetachNetwork(network.Name(), container.ID); err != nil ***REMOVED***
			logrus.Warnf("error detaching from network %s: %v", network.Name(), err)
			if err := daemon.clusterProvider.DetachNetwork(network.ID(), container.ID); err != nil ***REMOVED***
				logrus.Warnf("error detaching from network %s: %v", network.ID(), err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	attributes := map[string]string***REMOVED***
		"container": container.ID,
	***REMOVED***
	daemon.LogNetworkEventWithAttributes(network, "disconnect", attributes)
***REMOVED***

func (daemon *Daemon) initializeNetworking(container *container.Container) error ***REMOVED***
	var err error

	if container.HostConfig.NetworkMode.IsContainer() ***REMOVED***
		// we need to get the hosts files from the container to join
		nc, err := daemon.getNetworkedContainer(container.ID, container.HostConfig.NetworkMode.ConnectedContainer())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = daemon.initializeNetworkingPaths(container, nc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		container.Config.Hostname = nc.Config.Hostname
		container.Config.Domainname = nc.Config.Domainname
		return nil
	***REMOVED***

	if container.HostConfig.NetworkMode.IsHost() ***REMOVED***
		if container.Config.Hostname == "" ***REMOVED***
			container.Config.Hostname, err = os.Hostname()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := daemon.allocateNetwork(container); err != nil ***REMOVED***
		return err
	***REMOVED***

	return container.BuildHostnameFile()
***REMOVED***

func (daemon *Daemon) getNetworkedContainer(containerID, connectedContainerID string) (*container.Container, error) ***REMOVED***
	nc, err := daemon.GetContainer(connectedContainerID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if containerID == nc.ID ***REMOVED***
		return nil, fmt.Errorf("cannot join own network")
	***REMOVED***
	if !nc.IsRunning() ***REMOVED***
		err := fmt.Errorf("cannot join network of a non running container: %s", connectedContainerID)
		return nil, errdefs.Conflict(err)
	***REMOVED***
	if nc.IsRestarting() ***REMOVED***
		return nil, errContainerIsRestarting(connectedContainerID)
	***REMOVED***
	return nc, nil
***REMOVED***

func (daemon *Daemon) releaseNetwork(container *container.Container) ***REMOVED***
	start := time.Now()
	if daemon.netController == nil ***REMOVED***
		return
	***REMOVED***
	if container.HostConfig.NetworkMode.IsContainer() || container.Config.NetworkDisabled ***REMOVED***
		return
	***REMOVED***

	sid := container.NetworkSettings.SandboxID
	settings := container.NetworkSettings.Networks
	container.NetworkSettings.Ports = nil

	if sid == "" ***REMOVED***
		return
	***REMOVED***

	var networks []libnetwork.Network
	for n, epSettings := range settings ***REMOVED***
		if nw, err := daemon.FindNetwork(getNetworkID(n, epSettings.EndpointSettings)); err == nil ***REMOVED***
			networks = append(networks, nw)
		***REMOVED***

		if epSettings.EndpointSettings == nil ***REMOVED***
			continue
		***REMOVED***

		cleanOperationalData(epSettings)
	***REMOVED***

	sb, err := daemon.netController.SandboxByID(sid)
	if err != nil ***REMOVED***
		logrus.Warnf("error locating sandbox id %s: %v", sid, err)
		return
	***REMOVED***
	if err := sb.DisableService(); err != nil ***REMOVED***
		logrus.WithFields(logrus.Fields***REMOVED***"container": container.ID, "sandbox": sid***REMOVED***).WithError(err).Error("Error removing service from sandbox")
	***REMOVED***

	if err := sb.Delete(); err != nil ***REMOVED***
		logrus.Errorf("Error deleting sandbox id %s for container %s: %v", sid, container.ID, err)
	***REMOVED***

	for _, nw := range networks ***REMOVED***
		daemon.tryDetachContainerFromClusterNetwork(nw, container)
	***REMOVED***
	networkActions.WithValues("release").UpdateSince(start)
***REMOVED***

func errRemovalContainer(containerID string) error ***REMOVED***
	return fmt.Errorf("Container %s is marked for removal and cannot be connected or disconnected to the network", containerID)
***REMOVED***

// ConnectToNetwork connects a container to a network
func (daemon *Daemon) ConnectToNetwork(container *container.Container, idOrName string, endpointConfig *networktypes.EndpointSettings) error ***REMOVED***
	if endpointConfig == nil ***REMOVED***
		endpointConfig = &networktypes.EndpointSettings***REMOVED******REMOVED***
	***REMOVED***
	container.Lock()
	defer container.Unlock()

	if !container.Running ***REMOVED***
		if container.RemovalInProgress || container.Dead ***REMOVED***
			return errRemovalContainer(container.ID)
		***REMOVED***

		n, err := daemon.FindNetwork(idOrName)
		if err == nil && n != nil ***REMOVED***
			if err := daemon.updateNetworkConfig(container, n, endpointConfig, true); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			container.NetworkSettings.Networks[idOrName] = &network.EndpointSettings***REMOVED***
				EndpointSettings: endpointConfig,
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if !daemon.isNetworkHotPluggable() ***REMOVED***
		return fmt.Errorf(runtime.GOOS + " does not support connecting a running container to a network")
	***REMOVED*** else ***REMOVED***
		if err := daemon.connectToNetwork(container, idOrName, endpointConfig, true); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return container.CheckpointTo(daemon.containersReplica)
***REMOVED***

// DisconnectFromNetwork disconnects container from network n.
func (daemon *Daemon) DisconnectFromNetwork(container *container.Container, networkName string, force bool) error ***REMOVED***
	n, err := daemon.FindNetwork(networkName)
	container.Lock()
	defer container.Unlock()

	if !container.Running || (err != nil && force) ***REMOVED***
		if container.RemovalInProgress || container.Dead ***REMOVED***
			return errRemovalContainer(container.ID)
		***REMOVED***
		// In case networkName is resolved we will use n.Name()
		// this will cover the case where network id is passed.
		if n != nil ***REMOVED***
			networkName = n.Name()
		***REMOVED***
		if _, ok := container.NetworkSettings.Networks[networkName]; !ok ***REMOVED***
			return fmt.Errorf("container %s is not connected to the network %s", container.ID, networkName)
		***REMOVED***
		delete(container.NetworkSettings.Networks, networkName)
	***REMOVED*** else if err == nil && !daemon.isNetworkHotPluggable() ***REMOVED***
		return fmt.Errorf(runtime.GOOS + " does not support connecting a running container to a network")
	***REMOVED*** else if err == nil ***REMOVED***
		if container.HostConfig.NetworkMode.IsHost() && containertypes.NetworkMode(n.Type()).IsHost() ***REMOVED***
			return runconfig.ErrConflictHostNetwork
		***REMOVED***

		if err := daemon.disconnectFromNetwork(container, n, false); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return err
	***REMOVED***

	if err := container.CheckpointTo(daemon.containersReplica); err != nil ***REMOVED***
		return err
	***REMOVED***

	if n != nil ***REMOVED***
		daemon.LogNetworkEventWithAttributes(n, "disconnect", map[string]string***REMOVED***
			"container": container.ID,
		***REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

// ActivateContainerServiceBinding puts this container into load balancer active rotation and DNS response
func (daemon *Daemon) ActivateContainerServiceBinding(containerName string) error ***REMOVED***
	container, err := daemon.GetContainer(containerName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	sb := daemon.getNetworkSandbox(container)
	if sb == nil ***REMOVED***
		return fmt.Errorf("network sandbox does not exist for container %s", containerName)
	***REMOVED***
	return sb.EnableService()
***REMOVED***

// DeactivateContainerServiceBinding removes this container from load balancer active rotation, and DNS response
func (daemon *Daemon) DeactivateContainerServiceBinding(containerName string) error ***REMOVED***
	container, err := daemon.GetContainer(containerName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	sb := daemon.getNetworkSandbox(container)
	if sb == nil ***REMOVED***
		// If the network sandbox is not found, then there is nothing to deactivate
		logrus.Debugf("Could not find network sandbox for container %s on service binding deactivation request", containerName)
		return nil
	***REMOVED***
	return sb.DisableService()
***REMOVED***

func getNetworkID(name string, endpointSettings *networktypes.EndpointSettings) string ***REMOVED***
	// We only want to prefer NetworkID for user defined networks.
	// For systems like bridge, none, etc. the name is preferred (otherwise restart may cause issues)
	if containertypes.NetworkMode(name).IsUserDefined() && endpointSettings != nil && endpointSettings.NetworkID != "" ***REMOVED***
		return endpointSettings.NetworkID
	***REMOVED***
	return name
***REMOVED***
