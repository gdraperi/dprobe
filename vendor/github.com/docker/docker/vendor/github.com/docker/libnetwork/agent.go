package libnetwork

//go:generate protoc -I.:Godeps/_workspace/src/github.com/gogo/protobuf  --gogo_out=import_path=github.com/docker/libnetwork,Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto:. agent.proto

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/docker/go-events"
	"github.com/docker/libnetwork/cluster"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/networkdb"
	"github.com/docker/libnetwork/types"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
)

const (
	subsysGossip = "networking:gossip"
	subsysIPSec  = "networking:ipsec"
	keyringSize  = 3
)

// ByTime implements sort.Interface for []*types.EncryptionKey based on
// the LamportTime field.
type ByTime []*types.EncryptionKey

func (b ByTime) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b ByTime) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b ByTime) Less(i, j int) bool ***REMOVED*** return b[i].LamportTime < b[j].LamportTime ***REMOVED***

type agent struct ***REMOVED***
	networkDB         *networkdb.NetworkDB
	bindAddr          string
	advertiseAddr     string
	dataPathAddr      string
	coreCancelFuncs   []func()
	driverCancelFuncs map[string][]func()
	sync.Mutex
***REMOVED***

func (a *agent) dataPathAddress() string ***REMOVED***
	a.Lock()
	defer a.Unlock()
	if a.dataPathAddr != "" ***REMOVED***
		return a.dataPathAddr
	***REMOVED***
	return a.advertiseAddr
***REMOVED***

const libnetworkEPTable = "endpoint_table"

func getBindAddr(ifaceName string) (string, error) ***REMOVED***
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("failed to find interface %s: %v", ifaceName, err)
	***REMOVED***

	addrs, err := iface.Addrs()
	if err != nil ***REMOVED***
		return "", fmt.Errorf("failed to get interface addresses: %v", err)
	***REMOVED***

	for _, a := range addrs ***REMOVED***
		addr, ok := a.(*net.IPNet)
		if !ok ***REMOVED***
			continue
		***REMOVED***
		addrIP := addr.IP

		if addrIP.IsLinkLocalUnicast() ***REMOVED***
			continue
		***REMOVED***

		return addrIP.String(), nil
	***REMOVED***

	return "", fmt.Errorf("failed to get bind address")
***REMOVED***

func resolveAddr(addrOrInterface string) (string, error) ***REMOVED***
	// Try and see if this is a valid IP address
	if net.ParseIP(addrOrInterface) != nil ***REMOVED***
		return addrOrInterface, nil
	***REMOVED***

	addr, err := net.ResolveIPAddr("ip", addrOrInterface)
	if err != nil ***REMOVED***
		// If not a valid IP address, it should be a valid interface
		return getBindAddr(addrOrInterface)
	***REMOVED***
	return addr.String(), nil
***REMOVED***

func (c *controller) handleKeyChange(keys []*types.EncryptionKey) error ***REMOVED***
	drvEnc := discoverapi.DriverEncryptionUpdate***REMOVED******REMOVED***

	a := c.getAgent()
	if a == nil ***REMOVED***
		logrus.Debug("Skipping key change as agent is nil")
		return nil
	***REMOVED***

	// Find the deleted key. If the deleted key was the primary key,
	// a new primary key should be set before removing if from keyring.
	c.Lock()
	added := []byte***REMOVED******REMOVED***
	deleted := []byte***REMOVED******REMOVED***
	j := len(c.keys)
	for i := 0; i < j; ***REMOVED***
		same := false
		for _, key := range keys ***REMOVED***
			if same = key.LamportTime == c.keys[i].LamportTime; same ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if !same ***REMOVED***
			cKey := c.keys[i]
			if cKey.Subsystem == subsysGossip ***REMOVED***
				deleted = cKey.Key
			***REMOVED***

			if cKey.Subsystem == subsysIPSec ***REMOVED***
				drvEnc.Prune = cKey.Key
				drvEnc.PruneTag = cKey.LamportTime
			***REMOVED***
			c.keys[i], c.keys[j-1] = c.keys[j-1], c.keys[i]
			c.keys[j-1] = nil
			j--
		***REMOVED***
		i++
	***REMOVED***
	c.keys = c.keys[:j]

	// Find the new key and add it to the key ring
	for _, key := range keys ***REMOVED***
		same := false
		for _, cKey := range c.keys ***REMOVED***
			if same = cKey.LamportTime == key.LamportTime; same ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if !same ***REMOVED***
			c.keys = append(c.keys, key)
			if key.Subsystem == subsysGossip ***REMOVED***
				added = key.Key
			***REMOVED***

			if key.Subsystem == subsysIPSec ***REMOVED***
				drvEnc.Key = key.Key
				drvEnc.Tag = key.LamportTime
			***REMOVED***
		***REMOVED***
	***REMOVED***
	c.Unlock()

	if len(added) > 0 ***REMOVED***
		a.networkDB.SetKey(added)
	***REMOVED***

	key, _, err := c.getPrimaryKeyTag(subsysGossip)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.networkDB.SetPrimaryKey(key)

	key, tag, err := c.getPrimaryKeyTag(subsysIPSec)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	drvEnc.Primary = key
	drvEnc.PrimaryTag = tag

	if len(deleted) > 0 ***REMOVED***
		a.networkDB.RemoveKey(deleted)
	***REMOVED***

	c.drvRegistry.WalkDrivers(func(name string, driver driverapi.Driver, capability driverapi.Capability) bool ***REMOVED***
		err := driver.DiscoverNew(discoverapi.EncryptionKeysUpdate, drvEnc)
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to update datapath keys in driver %s: %v", name, err)
		***REMOVED***
		return false
	***REMOVED***)

	return nil
***REMOVED***

func (c *controller) agentSetup(clusterProvider cluster.Provider) error ***REMOVED***
	agent := c.getAgent()

	// If the agent is already present there is no need to try to initilize it again
	if agent != nil ***REMOVED***
		return nil
	***REMOVED***

	bindAddr := clusterProvider.GetLocalAddress()
	advAddr := clusterProvider.GetAdvertiseAddress()
	dataAddr := clusterProvider.GetDataPathAddress()
	remoteList := clusterProvider.GetRemoteAddressList()
	remoteAddrList := make([]string, 0, len(remoteList))
	for _, remote := range remoteList ***REMOVED***
		addr, _, _ := net.SplitHostPort(remote)
		remoteAddrList = append(remoteAddrList, addr)
	***REMOVED***

	listen := clusterProvider.GetListenAddress()
	listenAddr, _, _ := net.SplitHostPort(listen)

	logrus.Infof("Initializing Libnetwork Agent Listen-Addr=%s Local-addr=%s Adv-addr=%s Data-addr=%s Remote-addr-list=%v MTU=%d",
		listenAddr, bindAddr, advAddr, dataAddr, remoteAddrList, c.Config().Daemon.NetworkControlPlaneMTU)
	if advAddr != "" && agent == nil ***REMOVED***
		if err := c.agentInit(listenAddr, bindAddr, advAddr, dataAddr); err != nil ***REMOVED***
			logrus.Errorf("error in agentInit: %v", err)
			return err
		***REMOVED***
		c.drvRegistry.WalkDrivers(func(name string, driver driverapi.Driver, capability driverapi.Capability) bool ***REMOVED***
			if capability.ConnectivityScope == datastore.GlobalScope ***REMOVED***
				c.agentDriverNotify(driver)
			***REMOVED***
			return false
		***REMOVED***)
	***REMOVED***

	if len(remoteAddrList) > 0 ***REMOVED***
		if err := c.agentJoin(remoteAddrList); err != nil ***REMOVED***
			logrus.Errorf("Error in joining gossip cluster : %v(join will be retried in background)", err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// For a given subsystem getKeys sorts the keys by lamport time and returns
// slice of keys and lamport time which can used as a unique tag for the keys
func (c *controller) getKeys(subsys string) ([][]byte, []uint64) ***REMOVED***
	c.Lock()
	defer c.Unlock()

	sort.Sort(ByTime(c.keys))

	keys := [][]byte***REMOVED******REMOVED***
	tags := []uint64***REMOVED******REMOVED***
	for _, key := range c.keys ***REMOVED***
		if key.Subsystem == subsys ***REMOVED***
			keys = append(keys, key.Key)
			tags = append(tags, key.LamportTime)
		***REMOVED***
	***REMOVED***

	keys[0], keys[1] = keys[1], keys[0]
	tags[0], tags[1] = tags[1], tags[0]
	return keys, tags
***REMOVED***

// getPrimaryKeyTag returns the primary key for a given subsystem from the
// list of sorted key and the associated tag
func (c *controller) getPrimaryKeyTag(subsys string) ([]byte, uint64, error) ***REMOVED***
	c.Lock()
	defer c.Unlock()
	sort.Sort(ByTime(c.keys))
	keys := []*types.EncryptionKey***REMOVED******REMOVED***
	for _, key := range c.keys ***REMOVED***
		if key.Subsystem == subsys ***REMOVED***
			keys = append(keys, key)
		***REMOVED***
	***REMOVED***
	return keys[1].Key, keys[1].LamportTime, nil
***REMOVED***

func (c *controller) agentInit(listenAddr, bindAddrOrInterface, advertiseAddr, dataPathAddr string) error ***REMOVED***
	bindAddr, err := resolveAddr(bindAddrOrInterface)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	keys, _ := c.getKeys(subsysGossip)

	netDBConf := networkdb.DefaultConfig()
	netDBConf.BindAddr = listenAddr
	netDBConf.AdvertiseAddr = advertiseAddr
	netDBConf.Keys = keys
	if c.Config().Daemon.NetworkControlPlaneMTU != 0 ***REMOVED***
		// Consider the MTU remove the IP hdr (IPv4 or IPv6) and the TCP/UDP hdr.
		// To be on the safe side let's cut 100 bytes
		netDBConf.PacketBufferSize = (c.Config().Daemon.NetworkControlPlaneMTU - 100)
		logrus.Debugf("Control plane MTU: %d will initialize NetworkDB with: %d",
			c.Config().Daemon.NetworkControlPlaneMTU, netDBConf.PacketBufferSize)
	***REMOVED***
	nDB, err := networkdb.New(netDBConf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Register the diagnose handlers
	c.DiagnoseServer.RegisterHandler(nDB, networkdb.NetDbPaths2Func)

	var cancelList []func()
	ch, cancel := nDB.Watch(libnetworkEPTable, "", "")
	cancelList = append(cancelList, cancel)
	nodeCh, cancel := nDB.Watch(networkdb.NodeTable, "", "")
	cancelList = append(cancelList, cancel)

	c.Lock()
	c.agent = &agent***REMOVED***
		networkDB:         nDB,
		bindAddr:          bindAddr,
		advertiseAddr:     advertiseAddr,
		dataPathAddr:      dataPathAddr,
		coreCancelFuncs:   cancelList,
		driverCancelFuncs: make(map[string][]func()),
	***REMOVED***
	c.Unlock()

	go c.handleTableEvents(ch, c.handleEpTableEvent)
	go c.handleTableEvents(nodeCh, c.handleNodeTableEvent)

	drvEnc := discoverapi.DriverEncryptionConfig***REMOVED******REMOVED***
	keys, tags := c.getKeys(subsysIPSec)
	drvEnc.Keys = keys
	drvEnc.Tags = tags

	c.drvRegistry.WalkDrivers(func(name string, driver driverapi.Driver, capability driverapi.Capability) bool ***REMOVED***
		err := driver.DiscoverNew(discoverapi.EncryptionKeysConfig, drvEnc)
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to set datapath keys in driver %s: %v", name, err)
		***REMOVED***
		return false
	***REMOVED***)

	c.WalkNetworks(joinCluster)

	return nil
***REMOVED***

func (c *controller) agentJoin(remoteAddrList []string) error ***REMOVED***
	agent := c.getAgent()
	if agent == nil ***REMOVED***
		return nil
	***REMOVED***
	return agent.networkDB.Join(remoteAddrList)
***REMOVED***

func (c *controller) agentDriverNotify(d driverapi.Driver) ***REMOVED***
	agent := c.getAgent()
	if agent == nil ***REMOVED***
		return
	***REMOVED***

	if err := d.DiscoverNew(discoverapi.NodeDiscovery, discoverapi.NodeDiscoveryData***REMOVED***
		Address:     agent.dataPathAddress(),
		BindAddress: agent.bindAddr,
		Self:        true,
	***REMOVED***); err != nil ***REMOVED***
		logrus.Warnf("Failed the node discovery in driver: %v", err)
	***REMOVED***

	drvEnc := discoverapi.DriverEncryptionConfig***REMOVED******REMOVED***
	keys, tags := c.getKeys(subsysIPSec)
	drvEnc.Keys = keys
	drvEnc.Tags = tags

	if err := d.DiscoverNew(discoverapi.EncryptionKeysConfig, drvEnc); err != nil ***REMOVED***
		logrus.Warnf("Failed to set datapath keys in driver: %v", err)
	***REMOVED***
***REMOVED***

func (c *controller) agentClose() ***REMOVED***
	// Acquire current agent instance and reset its pointer
	// then run closing functions
	c.Lock()
	agent := c.agent
	c.agent = nil
	c.Unlock()

	if agent == nil ***REMOVED***
		return
	***REMOVED***

	var cancelList []func()

	agent.Lock()
	for _, cancelFuncs := range agent.driverCancelFuncs ***REMOVED***
		cancelList = append(cancelList, cancelFuncs...)
	***REMOVED***

	// Add also the cancel functions for the network db
	cancelList = append(cancelList, agent.coreCancelFuncs...)
	agent.Unlock()

	for _, cancel := range cancelList ***REMOVED***
		cancel()
	***REMOVED***

	agent.networkDB.Close()
***REMOVED***

// Task has the backend container details
type Task struct ***REMOVED***
	Name       string
	EndpointID string
	EndpointIP string
	Info       map[string]string
***REMOVED***

// ServiceInfo has service specific details along with the list of backend tasks
type ServiceInfo struct ***REMOVED***
	VIP          string
	LocalLBIndex int
	Tasks        []Task
	Ports        []string
***REMOVED***

type epRecord struct ***REMOVED***
	ep      EndpointRecord
	info    map[string]string
	lbIndex int
***REMOVED***

func (n *network) Services() map[string]ServiceInfo ***REMOVED***
	eps := make(map[string]epRecord)

	if !n.isClusterEligible() ***REMOVED***
		return nil
	***REMOVED***
	agent := n.getController().getAgent()
	if agent == nil ***REMOVED***
		return nil
	***REMOVED***

	// Walk through libnetworkEPTable and fetch the driver agnostic endpoint info
	entries := agent.networkDB.GetTableByNetwork(libnetworkEPTable, n.id)
	for eid, value := range entries ***REMOVED***
		var epRec EndpointRecord
		nid := n.ID()
		if err := proto.Unmarshal(value.Value, &epRec); err != nil ***REMOVED***
			logrus.Errorf("Unmarshal of libnetworkEPTable failed for endpoint %s in network %s, %v", eid, nid, err)
			continue
		***REMOVED***
		i := n.getController().getLBIndex(epRec.ServiceID, nid, epRec.IngressPorts)
		eps[eid] = epRecord***REMOVED***
			ep:      epRec,
			lbIndex: i,
		***REMOVED***
	***REMOVED***

	// Walk through the driver's tables, have the driver decode the entries
	// and return the tuple ***REMOVED***ep ID, value***REMOVED***. value is a string that coveys
	// relevant info about the endpoint.
	d, err := n.driver(true)
	if err != nil ***REMOVED***
		logrus.Errorf("Could not resolve driver for network %s/%s while fetching services: %v", n.networkType, n.ID(), err)
		return nil
	***REMOVED***
	for _, table := range n.driverTables ***REMOVED***
		if table.objType != driverapi.EndpointObject ***REMOVED***
			continue
		***REMOVED***
		entries := agent.networkDB.GetTableByNetwork(table.name, n.id)
		for key, value := range entries ***REMOVED***
			epID, info := d.DecodeTableEntry(table.name, key, value.Value)
			if ep, ok := eps[epID]; !ok ***REMOVED***
				logrus.Errorf("Inconsistent driver and libnetwork state for endpoint %s", epID)
			***REMOVED*** else ***REMOVED***
				ep.info = info
				eps[epID] = ep
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// group the endpoints into a map keyed by the service name
	sinfo := make(map[string]ServiceInfo)
	for ep, epr := range eps ***REMOVED***
		var (
			s  ServiceInfo
			ok bool
		)
		if s, ok = sinfo[epr.ep.ServiceName]; !ok ***REMOVED***
			s = ServiceInfo***REMOVED***
				VIP:          epr.ep.VirtualIP,
				LocalLBIndex: epr.lbIndex,
			***REMOVED***
		***REMOVED***
		ports := []string***REMOVED******REMOVED***
		if s.Ports == nil ***REMOVED***
			for _, port := range epr.ep.IngressPorts ***REMOVED***
				p := fmt.Sprintf("Target: %d, Publish: %d", port.TargetPort, port.PublishedPort)
				ports = append(ports, p)
			***REMOVED***
			s.Ports = ports
		***REMOVED***
		s.Tasks = append(s.Tasks, Task***REMOVED***
			Name:       epr.ep.Name,
			EndpointID: ep,
			EndpointIP: epr.ep.EndpointIP,
			Info:       epr.info,
		***REMOVED***)
		sinfo[epr.ep.ServiceName] = s
	***REMOVED***
	return sinfo
***REMOVED***

func (n *network) isClusterEligible() bool ***REMOVED***
	if n.scope != datastore.SwarmScope || !n.driverIsMultihost() ***REMOVED***
		return false
	***REMOVED***
	return n.getController().getAgent() != nil
***REMOVED***

func (n *network) joinCluster() error ***REMOVED***
	if !n.isClusterEligible() ***REMOVED***
		return nil
	***REMOVED***

	agent := n.getController().getAgent()
	if agent == nil ***REMOVED***
		return nil
	***REMOVED***

	return agent.networkDB.JoinNetwork(n.ID())
***REMOVED***

func (n *network) leaveCluster() error ***REMOVED***
	if !n.isClusterEligible() ***REMOVED***
		return nil
	***REMOVED***

	agent := n.getController().getAgent()
	if agent == nil ***REMOVED***
		return nil
	***REMOVED***

	return agent.networkDB.LeaveNetwork(n.ID())
***REMOVED***

func (ep *endpoint) addDriverInfoToCluster() error ***REMOVED***
	n := ep.getNetwork()
	if !n.isClusterEligible() ***REMOVED***
		return nil
	***REMOVED***
	if ep.joinInfo == nil ***REMOVED***
		return nil
	***REMOVED***

	agent := n.getController().getAgent()
	if agent == nil ***REMOVED***
		return nil
	***REMOVED***

	for _, te := range ep.joinInfo.driverTableEntries ***REMOVED***
		if err := agent.networkDB.CreateEntry(te.tableName, n.ID(), te.key, te.value); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (ep *endpoint) deleteDriverInfoFromCluster() error ***REMOVED***
	n := ep.getNetwork()
	if !n.isClusterEligible() ***REMOVED***
		return nil
	***REMOVED***
	if ep.joinInfo == nil ***REMOVED***
		return nil
	***REMOVED***

	agent := n.getController().getAgent()
	if agent == nil ***REMOVED***
		return nil
	***REMOVED***

	for _, te := range ep.joinInfo.driverTableEntries ***REMOVED***
		if err := agent.networkDB.DeleteEntry(te.tableName, n.ID(), te.key); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (ep *endpoint) addServiceInfoToCluster(sb *sandbox) error ***REMOVED***
	if ep.isAnonymous() && len(ep.myAliases) == 0 || ep.Iface().Address() == nil ***REMOVED***
		return nil
	***REMOVED***

	n := ep.getNetwork()
	if !n.isClusterEligible() ***REMOVED***
		return nil
	***REMOVED***

	sb.Service.Lock()
	defer sb.Service.Unlock()
	logrus.Debugf("addServiceInfoToCluster START for %s %s", ep.svcName, ep.ID())

	// Check that the endpoint is still present on the sandbox before adding it to the service discovery.
	// This is to handle a race between the EnableService and the sbLeave
	// It is possible that the EnableService starts, fetches the list of the endpoints and
	// by the time the addServiceInfoToCluster is called the endpoint got removed from the sandbox
	// The risk is that the deleteServiceInfoToCluster happens before the addServiceInfoToCluster.
	// This check under the Service lock of the sandbox ensure the correct behavior.
	// If the addServiceInfoToCluster arrives first may find or not the endpoint and will proceed or exit
	// but in any case the deleteServiceInfoToCluster will follow doing the cleanup if needed.
	// In case the deleteServiceInfoToCluster arrives first, this one is happening after the endpoint is
	// removed from the list, in this situation the delete will bail out not finding any data to cleanup
	// and the add will bail out not finding the endpoint on the sandbox.
	if e := sb.getEndpoint(ep.ID()); e == nil ***REMOVED***
		logrus.Warnf("addServiceInfoToCluster suppressing service resolution ep is not anymore in the sandbox %s", ep.ID())
		return nil
	***REMOVED***

	c := n.getController()
	agent := c.getAgent()

	name := ep.Name()
	if ep.isAnonymous() ***REMOVED***
		name = ep.MyAliases()[0]
	***REMOVED***

	var ingressPorts []*PortConfig
	if ep.svcID != "" ***REMOVED***
		// This is a task part of a service
		// Gossip ingress ports only in ingress network.
		if n.ingress ***REMOVED***
			ingressPorts = ep.ingressPorts
		***REMOVED***
		if err := c.addServiceBinding(ep.svcName, ep.svcID, n.ID(), ep.ID(), name, ep.virtualIP, ingressPorts, ep.svcAliases, ep.myAliases, ep.Iface().Address().IP, "addServiceInfoToCluster"); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// This is a container simply attached to an attachable network
		if err := c.addContainerNameResolution(n.ID(), ep.ID(), name, ep.myAliases, ep.Iface().Address().IP, "addServiceInfoToCluster"); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	buf, err := proto.Marshal(&EndpointRecord***REMOVED***
		Name:         name,
		ServiceName:  ep.svcName,
		ServiceID:    ep.svcID,
		VirtualIP:    ep.virtualIP.String(),
		IngressPorts: ingressPorts,
		Aliases:      ep.svcAliases,
		TaskAliases:  ep.myAliases,
		EndpointIP:   ep.Iface().Address().IP.String(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if agent != nil ***REMOVED***
		if err := agent.networkDB.CreateEntry(libnetworkEPTable, n.ID(), ep.ID(), buf); err != nil ***REMOVED***
			logrus.Warnf("addServiceInfoToCluster NetworkDB CreateEntry failed for %s %s err:%s", ep.id, n.id, err)
			return err
		***REMOVED***
	***REMOVED***

	logrus.Debugf("addServiceInfoToCluster END for %s %s", ep.svcName, ep.ID())

	return nil
***REMOVED***

func (ep *endpoint) deleteServiceInfoFromCluster(sb *sandbox, method string) error ***REMOVED***
	if ep.isAnonymous() && len(ep.myAliases) == 0 ***REMOVED***
		return nil
	***REMOVED***

	n := ep.getNetwork()
	if !n.isClusterEligible() ***REMOVED***
		return nil
	***REMOVED***

	sb.Service.Lock()
	defer sb.Service.Unlock()
	logrus.Debugf("deleteServiceInfoFromCluster from %s START for %s %s", method, ep.svcName, ep.ID())

	c := n.getController()
	agent := c.getAgent()

	name := ep.Name()
	if ep.isAnonymous() ***REMOVED***
		name = ep.MyAliases()[0]
	***REMOVED***

	if agent != nil ***REMOVED***
		// First delete from networkDB then locally
		if err := agent.networkDB.DeleteEntry(libnetworkEPTable, n.ID(), ep.ID()); err != nil ***REMOVED***
			logrus.Warnf("deleteServiceInfoFromCluster NetworkDB DeleteEntry failed for %s %s err:%s", ep.id, n.id, err)
		***REMOVED***
	***REMOVED***

	if ep.Iface().Address() != nil ***REMOVED***
		if ep.svcID != "" ***REMOVED***
			// This is a task part of a service
			var ingressPorts []*PortConfig
			if n.ingress ***REMOVED***
				ingressPorts = ep.ingressPorts
			***REMOVED***
			if err := c.rmServiceBinding(ep.svcName, ep.svcID, n.ID(), ep.ID(), name, ep.virtualIP, ingressPorts, ep.svcAliases, ep.myAliases, ep.Iface().Address().IP, "deleteServiceInfoFromCluster", true); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// This is a container simply attached to an attachable network
			if err := c.delContainerNameResolution(n.ID(), ep.ID(), name, ep.myAliases, ep.Iface().Address().IP, "deleteServiceInfoFromCluster"); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	logrus.Debugf("deleteServiceInfoFromCluster from %s END for %s %s", method, ep.svcName, ep.ID())

	return nil
***REMOVED***

func (n *network) addDriverWatches() ***REMOVED***
	if !n.isClusterEligible() ***REMOVED***
		return
	***REMOVED***

	c := n.getController()
	agent := c.getAgent()
	if agent == nil ***REMOVED***
		return
	***REMOVED***
	for _, table := range n.driverTables ***REMOVED***
		ch, cancel := agent.networkDB.Watch(table.name, n.ID(), "")
		agent.Lock()
		agent.driverCancelFuncs[n.ID()] = append(agent.driverCancelFuncs[n.ID()], cancel)
		agent.Unlock()
		go c.handleTableEvents(ch, n.handleDriverTableEvent)
		d, err := n.driver(false)
		if err != nil ***REMOVED***
			logrus.Errorf("Could not resolve driver %s while walking driver tabl: %v", n.networkType, err)
			return
		***REMOVED***

		agent.networkDB.WalkTable(table.name, func(nid, key string, value []byte, deleted bool) bool ***REMOVED***
			// skip the entries that are mark for deletion, this is safe because this function is
			// called at initialization time so there is no state to delete
			if nid == n.ID() && !deleted ***REMOVED***
				d.EventNotify(driverapi.Create, nid, table.name, key, value)
			***REMOVED***
			return false
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (n *network) cancelDriverWatches() ***REMOVED***
	if !n.isClusterEligible() ***REMOVED***
		return
	***REMOVED***

	agent := n.getController().getAgent()
	if agent == nil ***REMOVED***
		return
	***REMOVED***

	agent.Lock()
	cancelFuncs := agent.driverCancelFuncs[n.ID()]
	delete(agent.driverCancelFuncs, n.ID())
	agent.Unlock()

	for _, cancel := range cancelFuncs ***REMOVED***
		cancel()
	***REMOVED***
***REMOVED***

func (c *controller) handleTableEvents(ch *events.Channel, fn func(events.Event)) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case ev := <-ch.C:
			fn(ev)
		case <-ch.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *network) handleDriverTableEvent(ev events.Event) ***REMOVED***
	d, err := n.driver(false)
	if err != nil ***REMOVED***
		logrus.Errorf("Could not resolve driver %s while handling driver table event: %v", n.networkType, err)
		return
	***REMOVED***

	var (
		etype driverapi.EventType
		tname string
		key   string
		value []byte
	)

	switch event := ev.(type) ***REMOVED***
	case networkdb.CreateEvent:
		tname = event.Table
		key = event.Key
		value = event.Value
		etype = driverapi.Create
	case networkdb.DeleteEvent:
		tname = event.Table
		key = event.Key
		value = event.Value
		etype = driverapi.Delete
	case networkdb.UpdateEvent:
		tname = event.Table
		key = event.Key
		value = event.Value
		etype = driverapi.Delete
	***REMOVED***

	d.EventNotify(etype, n.ID(), tname, key, value)
***REMOVED***

func (c *controller) handleNodeTableEvent(ev events.Event) ***REMOVED***
	var (
		value    []byte
		isAdd    bool
		nodeAddr networkdb.NodeAddr
	)
	switch event := ev.(type) ***REMOVED***
	case networkdb.CreateEvent:
		value = event.Value
		isAdd = true
	case networkdb.DeleteEvent:
		value = event.Value
	case networkdb.UpdateEvent:
		logrus.Errorf("Unexpected update node table event = %#v", event)
	***REMOVED***

	err := json.Unmarshal(value, &nodeAddr)
	if err != nil ***REMOVED***
		logrus.Errorf("Error unmarshalling node table event %v", err)
		return
	***REMOVED***
	c.processNodeDiscovery([]net.IP***REMOVED***nodeAddr.Addr***REMOVED***, isAdd)

***REMOVED***

func (c *controller) handleEpTableEvent(ev events.Event) ***REMOVED***
	var (
		nid   string
		eid   string
		value []byte
		isAdd bool
		epRec EndpointRecord
	)

	switch event := ev.(type) ***REMOVED***
	case networkdb.CreateEvent:
		nid = event.NetworkID
		eid = event.Key
		value = event.Value
		isAdd = true
	case networkdb.DeleteEvent:
		nid = event.NetworkID
		eid = event.Key
		value = event.Value
	case networkdb.UpdateEvent:
		logrus.Errorf("Unexpected update service table event = %#v", event)
		return
	***REMOVED***

	err := proto.Unmarshal(value, &epRec)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to unmarshal service table value: %v", err)
		return
	***REMOVED***

	containerName := epRec.Name
	svcName := epRec.ServiceName
	svcID := epRec.ServiceID
	vip := net.ParseIP(epRec.VirtualIP)
	ip := net.ParseIP(epRec.EndpointIP)
	ingressPorts := epRec.IngressPorts
	serviceAliases := epRec.Aliases
	taskAliases := epRec.TaskAliases

	if containerName == "" || ip == nil ***REMOVED***
		logrus.Errorf("Invalid endpoint name/ip received while handling service table event %s", value)
		return
	***REMOVED***

	if isAdd ***REMOVED***
		logrus.Debugf("handleEpTableEvent ADD %s R:%v", eid, epRec)
		if svcID != "" ***REMOVED***
			// This is a remote task part of a service
			if err := c.addServiceBinding(svcName, svcID, nid, eid, containerName, vip, ingressPorts, serviceAliases, taskAliases, ip, "handleEpTableEvent"); err != nil ***REMOVED***
				logrus.Errorf("failed adding service binding for %s epRec:%v err:%v", eid, epRec, err)
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// This is a remote container simply attached to an attachable network
			if err := c.addContainerNameResolution(nid, eid, containerName, taskAliases, ip, "handleEpTableEvent"); err != nil ***REMOVED***
				logrus.Errorf("failed adding container name resolution for %s epRec:%v err:%v", eid, epRec, err)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		logrus.Debugf("handleEpTableEvent DEL %s R:%v", eid, epRec)
		if svcID != "" ***REMOVED***
			// This is a remote task part of a service
			if err := c.rmServiceBinding(svcName, svcID, nid, eid, containerName, vip, ingressPorts, serviceAliases, taskAliases, ip, "handleEpTableEvent", true); err != nil ***REMOVED***
				logrus.Errorf("failed removing service binding for %s epRec:%v err:%v", eid, epRec, err)
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// This is a remote container simply attached to an attachable network
			if err := c.delContainerNameResolution(nid, eid, containerName, taskAliases, ip, "handleEpTableEvent"); err != nil ***REMOVED***
				logrus.Errorf("failed removing container name resolution for %s epRec:%v err:%v", eid, epRec, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
