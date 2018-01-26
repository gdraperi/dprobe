package osl

import (
	"fmt"
	"net"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// IfaceOption is a function option type to set interface options
type IfaceOption func(i *nwIface)

type nwIface struct ***REMOVED***
	srcName     string
	dstName     string
	master      string
	dstMaster   string
	mac         net.HardwareAddr
	address     *net.IPNet
	addressIPv6 *net.IPNet
	llAddrs     []*net.IPNet
	routes      []*net.IPNet
	bridge      bool
	ns          *networkNamespace
	sync.Mutex
***REMOVED***

func (i *nwIface) SrcName() string ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return i.srcName
***REMOVED***

func (i *nwIface) DstName() string ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return i.dstName
***REMOVED***

func (i *nwIface) DstMaster() string ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return i.dstMaster
***REMOVED***

func (i *nwIface) Bridge() bool ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return i.bridge
***REMOVED***

func (i *nwIface) Master() string ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return i.master
***REMOVED***

func (i *nwIface) MacAddress() net.HardwareAddr ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return types.GetMacCopy(i.mac)
***REMOVED***

func (i *nwIface) Address() *net.IPNet ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return types.GetIPNetCopy(i.address)
***REMOVED***

func (i *nwIface) AddressIPv6() *net.IPNet ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return types.GetIPNetCopy(i.addressIPv6)
***REMOVED***

func (i *nwIface) LinkLocalAddresses() []*net.IPNet ***REMOVED***
	i.Lock()
	defer i.Unlock()

	return i.llAddrs
***REMOVED***

func (i *nwIface) Routes() []*net.IPNet ***REMOVED***
	i.Lock()
	defer i.Unlock()

	routes := make([]*net.IPNet, len(i.routes))
	for index, route := range i.routes ***REMOVED***
		r := types.GetIPNetCopy(route)
		routes[index] = r
	***REMOVED***

	return routes
***REMOVED***

func (n *networkNamespace) Interfaces() []Interface ***REMOVED***
	n.Lock()
	defer n.Unlock()

	ifaces := make([]Interface, len(n.iFaces))

	for i, iface := range n.iFaces ***REMOVED***
		ifaces[i] = iface
	***REMOVED***

	return ifaces
***REMOVED***

func (i *nwIface) Remove() error ***REMOVED***
	i.Lock()
	n := i.ns
	i.Unlock()

	n.Lock()
	isDefault := n.isDefault
	nlh := n.nlHandle
	n.Unlock()

	// Find the network interface identified by the DstName attribute.
	iface, err := nlh.LinkByName(i.DstName())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Down the interface before configuring
	if err := nlh.LinkSetDown(iface); err != nil ***REMOVED***
		return err
	***REMOVED***

	err = nlh.LinkSetName(iface, i.SrcName())
	if err != nil ***REMOVED***
		logrus.Debugf("LinkSetName failed for interface %s: %v", i.SrcName(), err)
		return err
	***REMOVED***

	// if it is a bridge just delete it.
	if i.Bridge() ***REMOVED***
		if err := nlh.LinkDel(iface); err != nil ***REMOVED***
			return fmt.Errorf("failed deleting bridge %q: %v", i.SrcName(), err)
		***REMOVED***
	***REMOVED*** else if !isDefault ***REMOVED***
		// Move the network interface to caller namespace.
		if err := nlh.LinkSetNsFd(iface, ns.ParseHandlerInt()); err != nil ***REMOVED***
			logrus.Debugf("LinkSetNsPid failed for interface %s: %v", i.SrcName(), err)
			return err
		***REMOVED***
	***REMOVED***

	n.Lock()
	for index, intf := range n.iFaces ***REMOVED***
		if intf == i ***REMOVED***
			n.iFaces = append(n.iFaces[:index], n.iFaces[index+1:]...)
			break
		***REMOVED***
	***REMOVED***
	n.Unlock()

	n.checkLoV6()

	return nil
***REMOVED***

// Returns the sandbox's side veth interface statistics
func (i *nwIface) Statistics() (*types.InterfaceStatistics, error) ***REMOVED***
	i.Lock()
	n := i.ns
	i.Unlock()

	l, err := n.nlHandle.LinkByName(i.DstName())
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to retrieve the statistics for %s in netns %s: %v", i.DstName(), n.path, err)
	***REMOVED***

	stats := l.Attrs().Statistics
	if stats == nil ***REMOVED***
		return nil, fmt.Errorf("no statistics were returned")
	***REMOVED***

	return &types.InterfaceStatistics***REMOVED***
		RxBytes:   uint64(stats.RxBytes),
		TxBytes:   uint64(stats.TxBytes),
		RxPackets: uint64(stats.RxPackets),
		TxPackets: uint64(stats.TxPackets),
		RxDropped: uint64(stats.RxDropped),
		TxDropped: uint64(stats.TxDropped),
	***REMOVED***, nil
***REMOVED***

func (n *networkNamespace) findDst(srcName string, isBridge bool) string ***REMOVED***
	n.Lock()
	defer n.Unlock()

	for _, i := range n.iFaces ***REMOVED***
		// The master should match the srcname of the interface and the
		// master interface should be of type bridge, if searching for a bridge type
		if i.SrcName() == srcName && (!isBridge || i.Bridge()) ***REMOVED***
			return i.DstName()
		***REMOVED***
	***REMOVED***

	return ""
***REMOVED***

func (n *networkNamespace) AddInterface(srcName, dstPrefix string, options ...IfaceOption) error ***REMOVED***
	i := &nwIface***REMOVED***srcName: srcName, dstName: dstPrefix, ns: n***REMOVED***
	i.processInterfaceOptions(options...)

	if i.master != "" ***REMOVED***
		i.dstMaster = n.findDst(i.master, true)
		if i.dstMaster == "" ***REMOVED***
			return fmt.Errorf("could not find an appropriate master %q for %q",
				i.master, i.srcName)
		***REMOVED***
	***REMOVED***

	n.Lock()
	if n.isDefault ***REMOVED***
		i.dstName = i.srcName
	***REMOVED*** else ***REMOVED***
		i.dstName = fmt.Sprintf("%s%d", dstPrefix, n.nextIfIndex[dstPrefix])
		n.nextIfIndex[dstPrefix]++
	***REMOVED***

	path := n.path
	isDefault := n.isDefault
	nlh := n.nlHandle
	nlhHost := ns.NlHandle()
	n.Unlock()

	// If it is a bridge interface we have to create the bridge inside
	// the namespace so don't try to lookup the interface using srcName
	if i.bridge ***REMOVED***
		link := &netlink.Bridge***REMOVED***
			LinkAttrs: netlink.LinkAttrs***REMOVED***
				Name: i.srcName,
			***REMOVED***,
		***REMOVED***
		if err := nlh.LinkAdd(link); err != nil ***REMOVED***
			return fmt.Errorf("failed to create bridge %q: %v", i.srcName, err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Find the network interface identified by the SrcName attribute.
		iface, err := nlhHost.LinkByName(i.srcName)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to get link by name %q: %v", i.srcName, err)
		***REMOVED***

		// Move the network interface to the destination
		// namespace only if the namespace is not a default
		// type
		if !isDefault ***REMOVED***
			newNs, err := netns.GetFromPath(path)
			if err != nil ***REMOVED***
				return fmt.Errorf("failed get network namespace %q: %v", path, err)
			***REMOVED***
			defer newNs.Close()
			if err := nlhHost.LinkSetNsFd(iface, int(newNs)); err != nil ***REMOVED***
				return fmt.Errorf("failed to set namespace on link %q: %v", i.srcName, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Find the network interface identified by the SrcName attribute.
	iface, err := nlh.LinkByName(i.srcName)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get link by name %q: %v", i.srcName, err)
	***REMOVED***

	// Down the interface before configuring
	if err := nlh.LinkSetDown(iface); err != nil ***REMOVED***
		return fmt.Errorf("failed to set link down: %v", err)
	***REMOVED***

	// Configure the interface now this is moved in the proper namespace.
	if err := configureInterface(nlh, iface, i); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Up the interface.
	cnt := 0
	for err = nlh.LinkSetUp(iface); err != nil && cnt < 3; cnt++ ***REMOVED***
		logrus.Debugf("retrying link setup because of: %v", err)
		time.Sleep(10 * time.Millisecond)
		err = nlh.LinkSetUp(iface)
	***REMOVED***
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to set link up: %v", err)
	***REMOVED***

	// Set the routes on the interface. This can only be done when the interface is up.
	if err := setInterfaceRoutes(nlh, iface, i); err != nil ***REMOVED***
		return fmt.Errorf("error setting interface %q routes to %q: %v", iface.Attrs().Name, i.Routes(), err)
	***REMOVED***

	n.Lock()
	n.iFaces = append(n.iFaces, i)
	n.Unlock()

	n.checkLoV6()

	return nil
***REMOVED***

func configureInterface(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	ifaceName := iface.Attrs().Name
	ifaceConfigurators := []struct ***REMOVED***
		Fn         func(*netlink.Handle, netlink.Link, *nwIface) error
		ErrMessage string
	***REMOVED******REMOVED***
		***REMOVED***setInterfaceName, fmt.Sprintf("error renaming interface %q to %q", ifaceName, i.DstName())***REMOVED***,
		***REMOVED***setInterfaceMAC, fmt.Sprintf("error setting interface %q MAC to %q", ifaceName, i.MacAddress())***REMOVED***,
		***REMOVED***setInterfaceIP, fmt.Sprintf("error setting interface %q IP to %v", ifaceName, i.Address())***REMOVED***,
		***REMOVED***setInterfaceIPv6, fmt.Sprintf("error setting interface %q IPv6 to %v", ifaceName, i.AddressIPv6())***REMOVED***,
		***REMOVED***setInterfaceMaster, fmt.Sprintf("error setting interface %q master to %q", ifaceName, i.DstMaster())***REMOVED***,
		***REMOVED***setInterfaceLinkLocalIPs, fmt.Sprintf("error setting interface %q link local IPs to %v", ifaceName, i.LinkLocalAddresses())***REMOVED***,
	***REMOVED***

	for _, config := range ifaceConfigurators ***REMOVED***
		if err := config.Fn(nlh, iface, i); err != nil ***REMOVED***
			return fmt.Errorf("%s: %v", config.ErrMessage, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func setInterfaceMaster(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	if i.DstMaster() == "" ***REMOVED***
		return nil
	***REMOVED***

	return nlh.LinkSetMaster(iface, &netlink.Bridge***REMOVED***
		LinkAttrs: netlink.LinkAttrs***REMOVED***Name: i.DstMaster()***REMOVED******REMOVED***)
***REMOVED***

func setInterfaceMAC(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	if i.MacAddress() == nil ***REMOVED***
		return nil
	***REMOVED***
	return nlh.LinkSetHardwareAddr(iface, i.MacAddress())
***REMOVED***

func setInterfaceIP(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	if i.Address() == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := checkRouteConflict(nlh, i.Address(), netlink.FAMILY_V4); err != nil ***REMOVED***
		return err
	***REMOVED***
	ipAddr := &netlink.Addr***REMOVED***IPNet: i.Address(), Label: ""***REMOVED***
	return nlh.AddrAdd(iface, ipAddr)
***REMOVED***

func setInterfaceIPv6(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	if i.AddressIPv6() == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := checkRouteConflict(nlh, i.AddressIPv6(), netlink.FAMILY_V6); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := setIPv6(i.ns.path, i.DstName(), true); err != nil ***REMOVED***
		return fmt.Errorf("failed to enable ipv6: %v", err)
	***REMOVED***
	ipAddr := &netlink.Addr***REMOVED***IPNet: i.AddressIPv6(), Label: "", Flags: syscall.IFA_F_NODAD***REMOVED***
	return nlh.AddrAdd(iface, ipAddr)
***REMOVED***

func setInterfaceLinkLocalIPs(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	for _, llIP := range i.LinkLocalAddresses() ***REMOVED***
		ipAddr := &netlink.Addr***REMOVED***IPNet: llIP***REMOVED***
		if err := nlh.AddrAdd(iface, ipAddr); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func setInterfaceName(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	return nlh.LinkSetName(iface, i.DstName())
***REMOVED***

func setInterfaceRoutes(nlh *netlink.Handle, iface netlink.Link, i *nwIface) error ***REMOVED***
	for _, route := range i.Routes() ***REMOVED***
		err := nlh.RouteAdd(&netlink.Route***REMOVED***
			Scope:     netlink.SCOPE_LINK,
			LinkIndex: iface.Attrs().Index,
			Dst:       route,
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// In older kernels (like the one in Centos 6.6 distro) sysctl does not have netns support. Therefore
// we cannot gather the statistics from /sys/class/net/<dev>/statistics/<counter> files. Per-netns stats
// are naturally found in /proc/net/dev in kernels which support netns (ifconfig relies on that).
const (
	netStatsFile = "/proc/net/dev"
	base         = "[ ]*%s:([ ]+[0-9]+)***REMOVED***16***REMOVED***"
)

func scanInterfaceStats(data, ifName string, i *types.InterfaceStatistics) error ***REMOVED***
	var (
		bktStr string
		bkt    uint64
	)

	regex := fmt.Sprintf(base, ifName)
	re := regexp.MustCompile(regex)
	line := re.FindString(data)

	_, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
		&bktStr, &i.RxBytes, &i.RxPackets, &i.RxErrors, &i.RxDropped, &bkt, &bkt, &bkt,
		&bkt, &i.TxBytes, &i.TxPackets, &i.TxErrors, &i.TxDropped, &bkt, &bkt, &bkt, &bkt)

	return err
***REMOVED***

func checkRouteConflict(nlh *netlink.Handle, address *net.IPNet, family int) error ***REMOVED***
	routes, err := nlh.RouteList(nil, family)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, route := range routes ***REMOVED***
		if route.Dst != nil ***REMOVED***
			if route.Dst.Contains(address.IP) || address.Contains(route.Dst.IP) ***REMOVED***
				return fmt.Errorf("cannot program address %v in sandbox interface because it conflicts with existing route %s",
					address, route)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
