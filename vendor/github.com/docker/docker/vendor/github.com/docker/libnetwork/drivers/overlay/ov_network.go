package overlay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/netutils"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/resolvconf"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

var (
	hostMode    bool
	networkOnce sync.Once
	networkMu   sync.Mutex
	vniTbl      = make(map[uint32]string)
)

type networkTable map[string]*network

type subnet struct ***REMOVED***
	once      *sync.Once
	vxlanName string
	brName    string
	vni       uint32
	initErr   error
	subnetIP  *net.IPNet
	gwIP      *net.IPNet
***REMOVED***

type subnetJSON struct ***REMOVED***
	SubnetIP string
	GwIP     string
	Vni      uint32
***REMOVED***

type network struct ***REMOVED***
	id        string
	dbIndex   uint64
	dbExists  bool
	sbox      osl.Sandbox
	nlSocket  *nl.NetlinkSocket
	endpoints endpointTable
	driver    *driver
	joinCnt   int
	once      *sync.Once
	initEpoch int
	initErr   error
	subnets   []*subnet
	secure    bool
	mtu       int
	sync.Mutex
***REMOVED***

func init() ***REMOVED***
	reexec.Register("set-default-vlan", setDefaultVlan)
***REMOVED***

func setDefaultVlan() ***REMOVED***
	if len(os.Args) < 3 ***REMOVED***
		logrus.Error("insufficient number of arguments")
		os.Exit(1)
	***REMOVED***

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	nsPath := os.Args[1]
	ns, err := netns.GetFromPath(nsPath)
	if err != nil ***REMOVED***
		logrus.Errorf("overlay namespace get failed, %v", err)
		os.Exit(1)
	***REMOVED***
	if err = netns.Set(ns); err != nil ***REMOVED***
		logrus.Errorf("setting into overlay namespace failed, %v", err)
		os.Exit(1)
	***REMOVED***

	// make sure the sysfs mount doesn't propagate back
	if err = syscall.Unshare(syscall.CLONE_NEWNS); err != nil ***REMOVED***
		logrus.Errorf("unshare failed, %v", err)
		os.Exit(1)
	***REMOVED***

	flag := syscall.MS_PRIVATE | syscall.MS_REC
	if err = syscall.Mount("", "/", "", uintptr(flag), ""); err != nil ***REMOVED***
		logrus.Errorf("root mount failed, %v", err)
		os.Exit(1)
	***REMOVED***

	if err = syscall.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil ***REMOVED***
		logrus.Errorf("mounting sysfs failed, %v", err)
		os.Exit(1)
	***REMOVED***

	brName := os.Args[2]
	path := filepath.Join("/sys/class/net", brName, "bridge/default_pvid")
	data := []byte***REMOVED***'0', '\n'***REMOVED***

	if err = ioutil.WriteFile(path, data, 0644); err != nil ***REMOVED***
		logrus.Errorf("enabling default vlan on bridge %s failed %v", brName, err)
		os.Exit(1)
	***REMOVED***
	os.Exit(0)
***REMOVED***

func (d *driver) NetworkAllocate(id string, option map[string]string, ipV4Data, ipV6Data []driverapi.IPAMData) (map[string]string, error) ***REMOVED***
	return nil, types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) NetworkFree(id string) error ***REMOVED***
	return types.NotImplementedErrorf("not implemented")
***REMOVED***

func (d *driver) CreateNetwork(id string, option map[string]interface***REMOVED******REMOVED***, nInfo driverapi.NetworkInfo, ipV4Data, ipV6Data []driverapi.IPAMData) error ***REMOVED***
	if id == "" ***REMOVED***
		return fmt.Errorf("invalid network id")
	***REMOVED***
	if len(ipV4Data) == 0 || ipV4Data[0].Pool.String() == "0.0.0.0/0" ***REMOVED***
		return types.BadRequestErrorf("ipv4 pool is empty")
	***REMOVED***

	// Since we perform lazy configuration make sure we try
	// configuring the driver when we enter CreateNetwork
	if err := d.configure(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := &network***REMOVED***
		id:        id,
		driver:    d,
		endpoints: endpointTable***REMOVED******REMOVED***,
		once:      &sync.Once***REMOVED******REMOVED***,
		subnets:   []*subnet***REMOVED******REMOVED***,
	***REMOVED***

	vnis := make([]uint32, 0, len(ipV4Data))
	if gval, ok := option[netlabel.GenericData]; ok ***REMOVED***
		optMap := gval.(map[string]string)
		if val, ok := optMap[netlabel.OverlayVxlanIDList]; ok ***REMOVED***
			logrus.Debugf("overlay: Received vxlan IDs: %s", val)
			vniStrings := strings.Split(val, ",")
			for _, vniStr := range vniStrings ***REMOVED***
				vni, err := strconv.Atoi(vniStr)
				if err != nil ***REMOVED***
					return fmt.Errorf("invalid vxlan id value %q passed", vniStr)
				***REMOVED***

				vnis = append(vnis, uint32(vni))
			***REMOVED***
		***REMOVED***
		if _, ok := optMap[secureOption]; ok ***REMOVED***
			n.secure = true
		***REMOVED***
		if val, ok := optMap[netlabel.DriverMTU]; ok ***REMOVED***
			var err error
			if n.mtu, err = strconv.Atoi(val); err != nil ***REMOVED***
				return fmt.Errorf("failed to parse %v: %v", val, err)
			***REMOVED***
			if n.mtu < 0 ***REMOVED***
				return fmt.Errorf("invalid MTU value: %v", n.mtu)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// If we are getting vnis from libnetwork, either we get for
	// all subnets or none.
	if len(vnis) != 0 && len(vnis) < len(ipV4Data) ***REMOVED***
		return fmt.Errorf("insufficient vnis(%d) passed to overlay", len(vnis))
	***REMOVED***

	for i, ipd := range ipV4Data ***REMOVED***
		s := &subnet***REMOVED***
			subnetIP: ipd.Pool,
			gwIP:     ipd.Gateway,
			once:     &sync.Once***REMOVED******REMOVED***,
		***REMOVED***

		if len(vnis) != 0 ***REMOVED***
			s.vni = vnis[i]
		***REMOVED***

		n.subnets = append(n.subnets, s)
	***REMOVED***

	if err := n.writeToStore(); err != nil ***REMOVED***
		return fmt.Errorf("failed to update data store for network %v: %v", n.id, err)
	***REMOVED***

	// Make sure no rule is on the way from any stale secure network
	if !n.secure ***REMOVED***
		for _, vni := range vnis ***REMOVED***
			programMangle(vni, false)
			programInput(vni, false)
		***REMOVED***
	***REMOVED***

	if nInfo != nil ***REMOVED***
		if err := nInfo.TableEventRegister(ovPeerTable, driverapi.EndpointObject); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	d.addNetwork(n)
	return nil
***REMOVED***

func (d *driver) DeleteNetwork(nid string) error ***REMOVED***
	if nid == "" ***REMOVED***
		return fmt.Errorf("invalid network id")
	***REMOVED***

	// Make sure driver resources are initialized before proceeding
	if err := d.configure(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return fmt.Errorf("could not find network with id %s", nid)
	***REMOVED***

	for _, ep := range n.endpoints ***REMOVED***
		if ep.ifName != "" ***REMOVED***
			if link, err := ns.NlHandle().LinkByName(ep.ifName); err != nil ***REMOVED***
				ns.NlHandle().LinkDel(link)
			***REMOVED***
		***REMOVED***

		if err := d.deleteEndpointFromStore(ep); err != nil ***REMOVED***
			logrus.Warnf("Failed to delete overlay endpoint %s from local store: %v", ep.id[0:7], err)
		***REMOVED***
	***REMOVED***
	// flush the peerDB entries
	d.peerFlush(nid)
	d.deleteNetwork(nid)

	vnis, err := n.releaseVxlanID()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if n.secure ***REMOVED***
		for _, vni := range vnis ***REMOVED***
			programMangle(vni, false)
			programInput(vni, false)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) ProgramExternalConnectivity(nid, eid string, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (d *driver) RevokeExternalConnectivity(nid, eid string) error ***REMOVED***
	return nil
***REMOVED***

func (n *network) incEndpointCount() ***REMOVED***
	n.Lock()
	defer n.Unlock()
	n.joinCnt++
***REMOVED***

func (n *network) joinSandbox(restore bool) error ***REMOVED***
	// If there is a race between two go routines here only one will win
	// the other will wait.
	n.once.Do(func() ***REMOVED***
		// save the error status of initSandbox in n.initErr so that
		// all the racing go routines are able to know the status.
		n.initErr = n.initSandbox(restore)
	***REMOVED***)

	return n.initErr
***REMOVED***

func (n *network) joinSubnetSandbox(s *subnet, restore bool) error ***REMOVED***
	s.once.Do(func() ***REMOVED***
		s.initErr = n.initSubnetSandbox(s, restore)
	***REMOVED***)
	return s.initErr
***REMOVED***

func (n *network) leaveSandbox() ***REMOVED***
	n.Lock()
	defer n.Unlock()
	n.joinCnt--
	if n.joinCnt != 0 ***REMOVED***
		return
	***REMOVED***

	// We are about to destroy sandbox since the container is leaving the network
	// Reinitialize the once variable so that we will be able to trigger one time
	// sandbox initialization(again) when another container joins subsequently.
	n.once = &sync.Once***REMOVED******REMOVED***
	for _, s := range n.subnets ***REMOVED***
		s.once = &sync.Once***REMOVED******REMOVED***
	***REMOVED***

	n.destroySandbox()
***REMOVED***

// to be called while holding network lock
func (n *network) destroySandbox() ***REMOVED***
	if n.sbox != nil ***REMOVED***
		for _, iface := range n.sbox.Info().Interfaces() ***REMOVED***
			if err := iface.Remove(); err != nil ***REMOVED***
				logrus.Debugf("Remove interface %s failed: %v", iface.SrcName(), err)
			***REMOVED***
		***REMOVED***

		for _, s := range n.subnets ***REMOVED***
			if hostMode ***REMOVED***
				if err := removeFilters(n.id[:12], s.brName); err != nil ***REMOVED***
					logrus.Warnf("Could not remove overlay filters: %v", err)
				***REMOVED***
			***REMOVED***

			if s.vxlanName != "" ***REMOVED***
				err := deleteInterface(s.vxlanName)
				if err != nil ***REMOVED***
					logrus.Warnf("could not cleanup sandbox properly: %v", err)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if hostMode ***REMOVED***
			if err := removeNetworkChain(n.id[:12]); err != nil ***REMOVED***
				logrus.Warnf("could not remove network chain: %v", err)
			***REMOVED***
		***REMOVED***

		// Close the netlink socket, this will also release the watchMiss goroutine that is using it
		if n.nlSocket != nil ***REMOVED***
			n.nlSocket.Close()
			n.nlSocket = nil
		***REMOVED***

		n.sbox.Destroy()
		n.sbox = nil
	***REMOVED***
***REMOVED***

func populateVNITbl() ***REMOVED***
	filepath.Walk(filepath.Dir(osl.GenerateKey("walk")),
		func(path string, info os.FileInfo, err error) error ***REMOVED***
			_, fname := filepath.Split(path)

			if len(strings.Split(fname, "-")) <= 1 ***REMOVED***
				return nil
			***REMOVED***

			ns, err := netns.GetFromPath(path)
			if err != nil ***REMOVED***
				logrus.Errorf("Could not open namespace path %s during vni population: %v", path, err)
				return nil
			***REMOVED***
			defer ns.Close()

			nlh, err := netlink.NewHandleAt(ns, syscall.NETLINK_ROUTE)
			if err != nil ***REMOVED***
				logrus.Errorf("Could not open netlink handle during vni population for ns %s: %v", path, err)
				return nil
			***REMOVED***
			defer nlh.Delete()

			err = nlh.SetSocketTimeout(soTimeout)
			if err != nil ***REMOVED***
				logrus.Warnf("Failed to set the timeout on the netlink handle sockets for vni table population: %v", err)
			***REMOVED***

			links, err := nlh.LinkList()
			if err != nil ***REMOVED***
				logrus.Errorf("Failed to list interfaces during vni population for ns %s: %v", path, err)
				return nil
			***REMOVED***

			for _, l := range links ***REMOVED***
				if l.Type() == "vxlan" ***REMOVED***
					vniTbl[uint32(l.(*netlink.Vxlan).VxlanId)] = path
				***REMOVED***
			***REMOVED***

			return nil
		***REMOVED***)
***REMOVED***

func networkOnceInit() ***REMOVED***
	populateVNITbl()

	if os.Getenv("_OVERLAY_HOST_MODE") != "" ***REMOVED***
		hostMode = true
		return
	***REMOVED***

	err := createVxlan("testvxlan", 1, 0)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to create testvxlan interface: %v", err)
		return
	***REMOVED***

	defer deleteInterface("testvxlan")

	path := "/proc/self/ns/net"
	hNs, err := netns.GetFromPath(path)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to get network namespace from path %s while setting host mode: %v", path, err)
		return
	***REMOVED***
	defer hNs.Close()

	nlh := ns.NlHandle()

	iface, err := nlh.LinkByName("testvxlan")
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to get link testvxlan while setting host mode: %v", err)
		return
	***REMOVED***

	// If we are not able to move the vxlan interface to a namespace
	// then fallback to host mode
	if err := nlh.LinkSetNsFd(iface, int(hNs)); err != nil ***REMOVED***
		hostMode = true
	***REMOVED***
***REMOVED***

func (n *network) generateVxlanName(s *subnet) string ***REMOVED***
	id := n.id
	if len(n.id) > 5 ***REMOVED***
		id = n.id[:5]
	***REMOVED***

	return "vx-" + fmt.Sprintf("%06x", n.vxlanID(s)) + "-" + id
***REMOVED***

func (n *network) generateBridgeName(s *subnet) string ***REMOVED***
	id := n.id
	if len(n.id) > 5 ***REMOVED***
		id = n.id[:5]
	***REMOVED***

	return n.getBridgeNamePrefix(s) + "-" + id
***REMOVED***

func (n *network) getBridgeNamePrefix(s *subnet) string ***REMOVED***
	return "ov-" + fmt.Sprintf("%06x", n.vxlanID(s))
***REMOVED***

func checkOverlap(nw *net.IPNet) error ***REMOVED***
	var nameservers []string

	if rc, err := resolvconf.Get(); err == nil ***REMOVED***
		nameservers = resolvconf.GetNameserversAsCIDR(rc.Content)
	***REMOVED***

	if err := netutils.CheckNameserverOverlaps(nameservers, nw); err != nil ***REMOVED***
		return fmt.Errorf("overlay subnet %s failed check with nameserver: %v: %v", nw.String(), nameservers, err)
	***REMOVED***

	if err := netutils.CheckRouteOverlaps(nw); err != nil ***REMOVED***
		return fmt.Errorf("overlay subnet %s failed check with host route table: %v", nw.String(), err)
	***REMOVED***

	return nil
***REMOVED***

func (n *network) restoreSubnetSandbox(s *subnet, brName, vxlanName string) error ***REMOVED***
	sbox := n.sandbox()

	// restore overlay osl sandbox
	Ifaces := make(map[string][]osl.IfaceOption)
	brIfaceOption := make([]osl.IfaceOption, 2)
	brIfaceOption = append(brIfaceOption, sbox.InterfaceOptions().Address(s.gwIP))
	brIfaceOption = append(brIfaceOption, sbox.InterfaceOptions().Bridge(true))
	Ifaces[brName+"+br"] = brIfaceOption

	err := sbox.Restore(Ifaces, nil, nil, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	Ifaces = make(map[string][]osl.IfaceOption)
	vxlanIfaceOption := make([]osl.IfaceOption, 1)
	vxlanIfaceOption = append(vxlanIfaceOption, sbox.InterfaceOptions().Master(brName))
	Ifaces[vxlanName+"+vxlan"] = vxlanIfaceOption
	return sbox.Restore(Ifaces, nil, nil, nil)
***REMOVED***

func (n *network) setupSubnetSandbox(s *subnet, brName, vxlanName string) error ***REMOVED***

	if hostMode ***REMOVED***
		// Try to delete stale bridge interface if it exists
		if err := deleteInterface(brName); err != nil ***REMOVED***
			deleteInterfaceBySubnet(n.getBridgeNamePrefix(s), s)
		***REMOVED***
		// Try to delete the vxlan interface by vni if already present
		deleteVxlanByVNI("", n.vxlanID(s))

		if err := checkOverlap(s.subnetIP); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if !hostMode ***REMOVED***
		// Try to find this subnet's vni is being used in some
		// other namespace by looking at vniTbl that we just
		// populated in the once init. If a hit is found then
		// it must a stale namespace from previous
		// life. Destroy it completely and reclaim resourced.
		networkMu.Lock()
		path, ok := vniTbl[n.vxlanID(s)]
		networkMu.Unlock()

		if ok ***REMOVED***
			deleteVxlanByVNI(path, n.vxlanID(s))
			if err := syscall.Unmount(path, syscall.MNT_FORCE); err != nil ***REMOVED***
				logrus.Errorf("unmount of %s failed: %v", path, err)
			***REMOVED***
			os.Remove(path)

			networkMu.Lock()
			delete(vniTbl, n.vxlanID(s))
			networkMu.Unlock()
		***REMOVED***
	***REMOVED***

	// create a bridge and vxlan device for this subnet and move it to the sandbox
	sbox := n.sandbox()

	if err := sbox.AddInterface(brName, "br",
		sbox.InterfaceOptions().Address(s.gwIP),
		sbox.InterfaceOptions().Bridge(true)); err != nil ***REMOVED***
		return fmt.Errorf("bridge creation in sandbox failed for subnet %q: %v", s.subnetIP.String(), err)
	***REMOVED***

	err := createVxlan(vxlanName, n.vxlanID(s), n.maxMTU())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := sbox.AddInterface(vxlanName, "vxlan",
		sbox.InterfaceOptions().Master(brName)); err != nil ***REMOVED***
		return fmt.Errorf("vxlan interface creation failed for subnet %q: %v", s.subnetIP.String(), err)
	***REMOVED***

	if !hostMode ***REMOVED***
		var name string
		for _, i := range sbox.Info().Interfaces() ***REMOVED***
			if i.Bridge() ***REMOVED***
				name = i.DstName()
			***REMOVED***
		***REMOVED***
		cmd := &exec.Cmd***REMOVED***
			Path:   reexec.Self(),
			Args:   []string***REMOVED***"set-default-vlan", sbox.Key(), name***REMOVED***,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		***REMOVED***
		if err := cmd.Run(); err != nil ***REMOVED***
			// not a fatal error
			logrus.Errorf("reexec to set bridge default vlan failed %v", err)
		***REMOVED***
	***REMOVED***

	if hostMode ***REMOVED***
		if err := addFilters(n.id[:12], brName); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (n *network) initSubnetSandbox(s *subnet, restore bool) error ***REMOVED***
	brName := n.generateBridgeName(s)
	vxlanName := n.generateVxlanName(s)

	if restore ***REMOVED***
		if err := n.restoreSubnetSandbox(s, brName, vxlanName); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := n.setupSubnetSandbox(s, brName, vxlanName); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	n.Lock()
	s.vxlanName = vxlanName
	s.brName = brName
	n.Unlock()

	return nil
***REMOVED***

func (n *network) cleanupStaleSandboxes() ***REMOVED***
	filepath.Walk(filepath.Dir(osl.GenerateKey("walk")),
		func(path string, info os.FileInfo, err error) error ***REMOVED***
			_, fname := filepath.Split(path)

			pList := strings.Split(fname, "-")
			if len(pList) <= 1 ***REMOVED***
				return nil
			***REMOVED***

			pattern := pList[1]
			if strings.Contains(n.id, pattern) ***REMOVED***
				// Delete all vnis
				deleteVxlanByVNI(path, 0)
				syscall.Unmount(path, syscall.MNT_DETACH)
				os.Remove(path)

				// Now that we have destroyed this
				// sandbox, remove all references to
				// it in vniTbl so that we don't
				// inadvertently destroy the sandbox
				// created in this life.
				networkMu.Lock()
				for vni, tblPath := range vniTbl ***REMOVED***
					if tblPath == path ***REMOVED***
						delete(vniTbl, vni)
					***REMOVED***
				***REMOVED***
				networkMu.Unlock()
			***REMOVED***

			return nil
		***REMOVED***)
***REMOVED***

func (n *network) initSandbox(restore bool) error ***REMOVED***
	n.Lock()
	n.initEpoch++
	n.Unlock()

	networkOnce.Do(networkOnceInit)

	if !restore ***REMOVED***
		if hostMode ***REMOVED***
			if err := addNetworkChain(n.id[:12]); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		// If there are any stale sandboxes related to this network
		// from previous daemon life clean it up here
		n.cleanupStaleSandboxes()
	***REMOVED***

	// In the restore case network sandbox already exist; but we don't know
	// what epoch number it was created with. It has to be retrieved by
	// searching the net namespaces.
	var key string
	if restore ***REMOVED***
		key = osl.GenerateKey("-" + n.id)
	***REMOVED*** else ***REMOVED***
		key = osl.GenerateKey(fmt.Sprintf("%d-", n.initEpoch) + n.id)
	***REMOVED***

	sbox, err := osl.NewSandbox(key, !hostMode, restore)
	if err != nil ***REMOVED***
		return fmt.Errorf("could not get network sandbox (oper %t): %v", restore, err)
	***REMOVED***

	// this is needed to let the peerAdd configure the sandbox
	n.setSandbox(sbox)

	if !restore ***REMOVED***
		// Initialize the sandbox with all the peers previously received from networkdb
		n.driver.initSandboxPeerDB(n.id)
	***REMOVED***

	// If we are in swarm mode, we don't need anymore the watchMiss routine.
	// This will save 1 thread and 1 netlink socket per network
	if !n.driver.isSerfAlive() ***REMOVED***
		return nil
	***REMOVED***

	var nlSock *nl.NetlinkSocket
	sbox.InvokeFunc(func() ***REMOVED***
		nlSock, err = nl.Subscribe(syscall.NETLINK_ROUTE, syscall.RTNLGRP_NEIGH)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		// set the receive timeout to not remain stuck on the RecvFrom if the fd gets closed
		tv := syscall.NsecToTimeval(soTimeout.Nanoseconds())
		err = nlSock.SetReceiveTimeout(&tv)
	***REMOVED***)
	n.setNetlinkSocket(nlSock)

	if err == nil ***REMOVED***
		go n.watchMiss(nlSock, key)
	***REMOVED*** else ***REMOVED***
		logrus.Errorf("failed to subscribe to neighbor group netlink messages for overlay network %s in sbox %s: %v",
			n.id, sbox.Key(), err)
	***REMOVED***

	return nil
***REMOVED***

func (n *network) watchMiss(nlSock *nl.NetlinkSocket, nsPath string) ***REMOVED***
	// With the new version of the netlink library the deserialize function makes
	// requests about the interface of the netlink message. This can succeed only
	// if this go routine is in the target namespace. For this reason following we
	// lock the thread on that namespace
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	newNs, err := netns.GetFromPath(nsPath)
	if err != nil ***REMOVED***
		logrus.WithError(err).Errorf("failed to get the namespace %s", nsPath)
		return
	***REMOVED***
	defer newNs.Close()
	if err = netns.Set(newNs); err != nil ***REMOVED***
		logrus.WithError(err).Errorf("failed to enter the namespace %s", nsPath)
		return
	***REMOVED***
	for ***REMOVED***
		msgs, err := nlSock.Receive()
		if err != nil ***REMOVED***
			n.Lock()
			nlFd := nlSock.GetFd()
			n.Unlock()
			if nlFd == -1 ***REMOVED***
				// The netlink socket got closed, simply exit to not leak this goroutine
				return
			***REMOVED***
			// When the receive timeout expires the receive will return EAGAIN
			if err == syscall.EAGAIN ***REMOVED***
				// we continue here to avoid spam for timeouts
				continue
			***REMOVED***
			logrus.Errorf("Failed to receive from netlink: %v ", err)
			continue
		***REMOVED***

		for _, msg := range msgs ***REMOVED***
			if msg.Header.Type != syscall.RTM_GETNEIGH && msg.Header.Type != syscall.RTM_NEWNEIGH ***REMOVED***
				continue
			***REMOVED***

			neigh, err := netlink.NeighDeserialize(msg.Data)
			if err != nil ***REMOVED***
				logrus.Errorf("Failed to deserialize netlink ndmsg: %v", err)
				continue
			***REMOVED***

			var (
				ip             net.IP
				mac            net.HardwareAddr
				l2Miss, l3Miss bool
			)
			if neigh.IP.To4() != nil ***REMOVED***
				ip = neigh.IP
				l3Miss = true
			***REMOVED*** else if neigh.HardwareAddr != nil ***REMOVED***
				mac = []byte(neigh.HardwareAddr)
				ip = net.IP(mac[2:])
				l2Miss = true
			***REMOVED*** else ***REMOVED***
				continue
			***REMOVED***

			// Not any of the network's subnets. Ignore.
			if !n.contains(ip) ***REMOVED***
				continue
			***REMOVED***

			if neigh.State&(netlink.NUD_STALE|netlink.NUD_INCOMPLETE) == 0 ***REMOVED***
				continue
			***REMOVED***

			logrus.Debugf("miss notification: dest IP %v, dest MAC %v", ip, mac)
			mac, IPmask, vtep, err := n.driver.resolvePeer(n.id, ip)
			if err != nil ***REMOVED***
				logrus.Errorf("could not resolve peer %q: %v", ip, err)
				continue
			***REMOVED***
			n.driver.peerAdd(n.id, "dummy", ip, IPmask, mac, vtep, l2Miss, l3Miss, false)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *driver) addNetwork(n *network) ***REMOVED***
	d.Lock()
	d.networks[n.id] = n
	d.Unlock()
***REMOVED***

func (d *driver) deleteNetwork(nid string) ***REMOVED***
	d.Lock()
	delete(d.networks, nid)
	d.Unlock()
***REMOVED***

func (d *driver) network(nid string) *network ***REMOVED***
	d.Lock()
	n, ok := d.networks[nid]
	d.Unlock()
	if !ok ***REMOVED***
		n = d.getNetworkFromStore(nid)
		if n != nil ***REMOVED***
			n.driver = d
			n.endpoints = endpointTable***REMOVED******REMOVED***
			n.once = &sync.Once***REMOVED******REMOVED***
			d.Lock()
			d.networks[nid] = n
			d.Unlock()
		***REMOVED***
	***REMOVED***

	return n
***REMOVED***

func (d *driver) getNetworkFromStore(nid string) *network ***REMOVED***
	if d.store == nil ***REMOVED***
		return nil
	***REMOVED***

	n := &network***REMOVED***id: nid***REMOVED***
	if err := d.store.GetObject(datastore.Key(n.Key()...), n); err != nil ***REMOVED***
		return nil
	***REMOVED***

	return n
***REMOVED***

func (n *network) sandbox() osl.Sandbox ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.sbox
***REMOVED***

func (n *network) setSandbox(sbox osl.Sandbox) ***REMOVED***
	n.Lock()
	n.sbox = sbox
	n.Unlock()
***REMOVED***

func (n *network) setNetlinkSocket(nlSk *nl.NetlinkSocket) ***REMOVED***
	n.Lock()
	n.nlSocket = nlSk
	n.Unlock()
***REMOVED***

func (n *network) vxlanID(s *subnet) uint32 ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return s.vni
***REMOVED***

func (n *network) setVxlanID(s *subnet, vni uint32) ***REMOVED***
	n.Lock()
	s.vni = vni
	n.Unlock()
***REMOVED***

func (n *network) Key() []string ***REMOVED***
	return []string***REMOVED***"overlay", "network", n.id***REMOVED***
***REMOVED***

func (n *network) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***"overlay", "network"***REMOVED***
***REMOVED***

func (n *network) Value() []byte ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***

	netJSON := []*subnetJSON***REMOVED******REMOVED***

	for _, s := range n.subnets ***REMOVED***
		sj := &subnetJSON***REMOVED***
			SubnetIP: s.subnetIP.String(),
			GwIP:     s.gwIP.String(),
			Vni:      s.vni,
		***REMOVED***
		netJSON = append(netJSON, sj)
	***REMOVED***

	m["secure"] = n.secure
	m["subnets"] = netJSON
	m["mtu"] = n.mtu
	b, err := json.Marshal(m)
	if err != nil ***REMOVED***
		return []byte***REMOVED******REMOVED***
	***REMOVED***

	return b
***REMOVED***

func (n *network) Index() uint64 ***REMOVED***
	return n.dbIndex
***REMOVED***

func (n *network) SetIndex(index uint64) ***REMOVED***
	n.dbIndex = index
	n.dbExists = true
***REMOVED***

func (n *network) Exists() bool ***REMOVED***
	return n.dbExists
***REMOVED***

func (n *network) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (n *network) SetValue(value []byte) error ***REMOVED***
	var (
		m       map[string]interface***REMOVED******REMOVED***
		newNet  bool
		isMap   = true
		netJSON = []*subnetJSON***REMOVED******REMOVED***
	)

	if err := json.Unmarshal(value, &m); err != nil ***REMOVED***
		err := json.Unmarshal(value, &netJSON)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		isMap = false
	***REMOVED***

	if len(n.subnets) == 0 ***REMOVED***
		newNet = true
	***REMOVED***

	if isMap ***REMOVED***
		if val, ok := m["secure"]; ok ***REMOVED***
			n.secure = val.(bool)
		***REMOVED***
		if val, ok := m["mtu"]; ok ***REMOVED***
			n.mtu = int(val.(float64))
		***REMOVED***
		bytes, err := json.Marshal(m["subnets"])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := json.Unmarshal(bytes, &netJSON); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	for _, sj := range netJSON ***REMOVED***
		subnetIPstr := sj.SubnetIP
		gwIPstr := sj.GwIP
		vni := sj.Vni

		subnetIP, _ := types.ParseCIDR(subnetIPstr)
		gwIP, _ := types.ParseCIDR(gwIPstr)

		if newNet ***REMOVED***
			s := &subnet***REMOVED***
				subnetIP: subnetIP,
				gwIP:     gwIP,
				vni:      vni,
				once:     &sync.Once***REMOVED******REMOVED***,
			***REMOVED***
			n.subnets = append(n.subnets, s)
		***REMOVED*** else ***REMOVED***
			sNet := n.getMatchingSubnet(subnetIP)
			if sNet != nil ***REMOVED***
				sNet.vni = vni
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (n *network) DataScope() string ***REMOVED***
	return datastore.GlobalScope
***REMOVED***

func (n *network) writeToStore() error ***REMOVED***
	if n.driver.store == nil ***REMOVED***
		return nil
	***REMOVED***

	return n.driver.store.PutObjectAtomic(n)
***REMOVED***

func (n *network) releaseVxlanID() ([]uint32, error) ***REMOVED***
	if len(n.subnets) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	if n.driver.store != nil ***REMOVED***
		if err := n.driver.store.DeleteObjectAtomic(n); err != nil ***REMOVED***
			if err == datastore.ErrKeyModified || err == datastore.ErrKeyNotFound ***REMOVED***
				// In both the above cases we can safely assume that the key has been removed by some other
				// instance and so simply get out of here
				return nil, nil
			***REMOVED***

			return nil, fmt.Errorf("failed to delete network to vxlan id map: %v", err)
		***REMOVED***
	***REMOVED***
	var vnis []uint32
	for _, s := range n.subnets ***REMOVED***
		if n.driver.vxlanIdm != nil ***REMOVED***
			vni := n.vxlanID(s)
			vnis = append(vnis, vni)
			n.driver.vxlanIdm.Release(uint64(vni))
		***REMOVED***

		n.setVxlanID(s, 0)
	***REMOVED***

	return vnis, nil
***REMOVED***

func (n *network) obtainVxlanID(s *subnet) error ***REMOVED***
	//return if the subnet already has a vxlan id assigned
	if s.vni != 0 ***REMOVED***
		return nil
	***REMOVED***

	if n.driver.store == nil ***REMOVED***
		return fmt.Errorf("no valid vxlan id and no datastore configured, cannot obtain vxlan id")
	***REMOVED***

	for ***REMOVED***
		if err := n.driver.store.GetObject(datastore.Key(n.Key()...), n); err != nil ***REMOVED***
			return fmt.Errorf("getting network %q from datastore failed %v", n.id, err)
		***REMOVED***

		if s.vni == 0 ***REMOVED***
			vxlanID, err := n.driver.vxlanIdm.GetID(true)
			if err != nil ***REMOVED***
				return fmt.Errorf("failed to allocate vxlan id: %v", err)
			***REMOVED***

			n.setVxlanID(s, uint32(vxlanID))
			if err := n.writeToStore(); err != nil ***REMOVED***
				n.driver.vxlanIdm.Release(uint64(n.vxlanID(s)))
				n.setVxlanID(s, 0)
				if err == datastore.ErrKeyModified ***REMOVED***
					continue
				***REMOVED***
				return fmt.Errorf("network %q failed to update data store: %v", n.id, err)
			***REMOVED***
			return nil
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// contains return true if the passed ip belongs to one the network's
// subnets
func (n *network) contains(ip net.IP) bool ***REMOVED***
	for _, s := range n.subnets ***REMOVED***
		if s.subnetIP.Contains(ip) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// getSubnetforIP returns the subnet to which the given IP belongs
func (n *network) getSubnetforIP(ip *net.IPNet) *subnet ***REMOVED***
	for _, s := range n.subnets ***REMOVED***
		// first check if the mask lengths are the same
		i, _ := s.subnetIP.Mask.Size()
		j, _ := ip.Mask.Size()
		if i != j ***REMOVED***
			continue
		***REMOVED***
		if s.subnetIP.Contains(ip.IP) ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// getMatchingSubnet return the network's subnet that matches the input
func (n *network) getMatchingSubnet(ip *net.IPNet) *subnet ***REMOVED***
	if ip == nil ***REMOVED***
		return nil
	***REMOVED***
	for _, s := range n.subnets ***REMOVED***
		// first check if the mask lengths are the same
		i, _ := s.subnetIP.Mask.Size()
		j, _ := ip.Mask.Size()
		if i != j ***REMOVED***
			continue
		***REMOVED***
		if s.subnetIP.IP.Equal(ip.IP) ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
