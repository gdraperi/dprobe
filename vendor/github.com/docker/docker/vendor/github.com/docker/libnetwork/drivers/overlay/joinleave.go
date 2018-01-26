package overlay

import (
	"fmt"
	"net"
	"syscall"

	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return fmt.Errorf("could not find network with id %s", nid)
	***REMOVED***

	ep := n.endpoint(eid)
	if ep == nil ***REMOVED***
		return fmt.Errorf("could not find endpoint with id %s", eid)
	***REMOVED***

	if n.secure && len(d.keys) == 0 ***REMOVED***
		return fmt.Errorf("cannot join secure network: encryption keys not present")
	***REMOVED***

	nlh := ns.NlHandle()

	if n.secure && !nlh.SupportsNetlinkFamily(syscall.NETLINK_XFRM) ***REMOVED***
		return fmt.Errorf("cannot join secure network: required modules to install IPSEC rules are missing on host")
	***REMOVED***

	s := n.getSubnetforIP(ep.addr)
	if s == nil ***REMOVED***
		return fmt.Errorf("could not find subnet for endpoint %s", eid)
	***REMOVED***

	if err := n.obtainVxlanID(s); err != nil ***REMOVED***
		return fmt.Errorf("couldn't get vxlan id for %q: %v", s.subnetIP.String(), err)
	***REMOVED***

	if err := n.joinSandbox(false); err != nil ***REMOVED***
		return fmt.Errorf("network sandbox join failed: %v", err)
	***REMOVED***

	if err := n.joinSubnetSandbox(s, false); err != nil ***REMOVED***
		return fmt.Errorf("subnet sandbox join failed for %q: %v", s.subnetIP.String(), err)
	***REMOVED***

	// joinSubnetSandbox gets called when an endpoint comes up on a new subnet in the
	// overlay network. Hence the Endpoint count should be updated outside joinSubnetSandbox
	n.incEndpointCount()

	sbox := n.sandbox()

	overlayIfName, containerIfName, err := createVethPair()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ep.ifName = containerIfName

	if err = d.writeEndpointToStore(ep); err != nil ***REMOVED***
		return fmt.Errorf("failed to update overlay endpoint %s to local data store: %v", ep.id[0:7], err)
	***REMOVED***

	// Set the container interface and its peer MTU to 1450 to allow
	// for 50 bytes vxlan encap (inner eth header(14) + outer IP(20) +
	// outer UDP(8) + vxlan header(8))
	mtu := n.maxMTU()

	veth, err := nlh.LinkByName(overlayIfName)
	if err != nil ***REMOVED***
		return fmt.Errorf("cound not find link by name %s: %v", overlayIfName, err)
	***REMOVED***
	err = nlh.LinkSetMTU(veth, mtu)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = sbox.AddInterface(overlayIfName, "veth",
		sbox.InterfaceOptions().Master(s.brName)); err != nil ***REMOVED***
		return fmt.Errorf("could not add veth pair inside the network sandbox: %v", err)
	***REMOVED***

	veth, err = nlh.LinkByName(containerIfName)
	if err != nil ***REMOVED***
		return fmt.Errorf("could not find link by name %s: %v", containerIfName, err)
	***REMOVED***
	err = nlh.LinkSetMTU(veth, mtu)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = nlh.LinkSetHardwareAddr(veth, ep.mac); err != nil ***REMOVED***
		return fmt.Errorf("could not set mac address (%v) to the container interface: %v", ep.mac, err)
	***REMOVED***

	for _, sub := range n.subnets ***REMOVED***
		if sub == s ***REMOVED***
			continue
		***REMOVED***
		if err = jinfo.AddStaticRoute(sub.subnetIP, types.NEXTHOP, s.gwIP.IP); err != nil ***REMOVED***
			logrus.Errorf("Adding subnet %s static route in network %q failed\n", s.subnetIP, n.id)
		***REMOVED***
	***REMOVED***

	if iNames := jinfo.InterfaceName(); iNames != nil ***REMOVED***
		err = iNames.SetNames(containerIfName, "eth")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	d.peerAdd(nid, eid, ep.addr.IP, ep.addr.Mask, ep.mac, net.ParseIP(d.advertiseAddress), false, false, true)

	if err = d.checkEncryption(nid, nil, n.vxlanID(s), true, true); err != nil ***REMOVED***
		logrus.Warn(err)
	***REMOVED***

	buf, err := proto.Marshal(&PeerRecord***REMOVED***
		EndpointIP:       ep.addr.String(),
		EndpointMAC:      ep.mac.String(),
		TunnelEndpointIP: d.advertiseAddress,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := jinfo.AddTableEntry(ovPeerTable, eid, buf); err != nil ***REMOVED***
		logrus.Errorf("overlay: Failed adding table entry to joininfo: %v", err)
	***REMOVED***

	d.pushLocalEndpointEvent("join", nid, eid)

	return nil
***REMOVED***

func (d *driver) DecodeTableEntry(tablename string, key string, value []byte) (string, map[string]string) ***REMOVED***
	if tablename != ovPeerTable ***REMOVED***
		logrus.Errorf("DecodeTableEntry: unexpected table name %s", tablename)
		return "", nil
	***REMOVED***

	var peer PeerRecord
	if err := proto.Unmarshal(value, &peer); err != nil ***REMOVED***
		logrus.Errorf("DecodeTableEntry: failed to unmarshal peer record for key %s: %v", key, err)
		return "", nil
	***REMOVED***

	return key, map[string]string***REMOVED***
		"Host IP": peer.TunnelEndpointIP,
	***REMOVED***
***REMOVED***

func (d *driver) EventNotify(etype driverapi.EventType, nid, tableName, key string, value []byte) ***REMOVED***
	if tableName != ovPeerTable ***REMOVED***
		logrus.Errorf("Unexpected table notification for table %s received", tableName)
		return
	***REMOVED***

	eid := key

	var peer PeerRecord
	if err := proto.Unmarshal(value, &peer); err != nil ***REMOVED***
		logrus.Errorf("Failed to unmarshal peer record: %v", err)
		return
	***REMOVED***

	// Ignore local peers. We already know about them and they
	// should not be added to vxlan fdb.
	if peer.TunnelEndpointIP == d.advertiseAddress ***REMOVED***
		return
	***REMOVED***

	addr, err := types.ParseCIDR(peer.EndpointIP)
	if err != nil ***REMOVED***
		logrus.Errorf("Invalid peer IP %s received in event notify", peer.EndpointIP)
		return
	***REMOVED***

	mac, err := net.ParseMAC(peer.EndpointMAC)
	if err != nil ***REMOVED***
		logrus.Errorf("Invalid mac %s received in event notify", peer.EndpointMAC)
		return
	***REMOVED***

	vtep := net.ParseIP(peer.TunnelEndpointIP)
	if vtep == nil ***REMOVED***
		logrus.Errorf("Invalid VTEP %s received in event notify", peer.TunnelEndpointIP)
		return
	***REMOVED***

	if etype == driverapi.Delete ***REMOVED***
		d.peerDelete(nid, eid, addr.IP, addr.Mask, mac, vtep, false)
		return
	***REMOVED***

	d.peerAdd(nid, eid, addr.IP, addr.Mask, mac, vtep, false, false, false)
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return fmt.Errorf("could not find network with id %s", nid)
	***REMOVED***

	ep := n.endpoint(eid)

	if ep == nil ***REMOVED***
		return types.InternalMaskableErrorf("could not find endpoint with id %s", eid)
	***REMOVED***

	if d.notifyCh != nil ***REMOVED***
		d.notifyCh <- ovNotify***REMOVED***
			action: "leave",
			nw:     n,
			ep:     ep,
		***REMOVED***
	***REMOVED***

	d.peerDelete(nid, eid, ep.addr.IP, ep.addr.Mask, ep.mac, net.ParseIP(d.advertiseAddress), true)

	n.leaveSandbox()

	return nil
***REMOVED***
