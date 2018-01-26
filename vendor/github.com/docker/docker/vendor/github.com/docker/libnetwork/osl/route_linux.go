package osl

import (
	"fmt"
	"net"

	"github.com/docker/libnetwork/types"
	"github.com/vishvananda/netlink"
)

func (n *networkNamespace) Gateway() net.IP ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.gw
***REMOVED***

func (n *networkNamespace) GatewayIPv6() net.IP ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.gwv6
***REMOVED***

func (n *networkNamespace) StaticRoutes() []*types.StaticRoute ***REMOVED***
	n.Lock()
	defer n.Unlock()

	routes := make([]*types.StaticRoute, len(n.staticRoutes))
	for i, route := range n.staticRoutes ***REMOVED***
		r := route.GetCopy()
		routes[i] = r
	***REMOVED***

	return routes
***REMOVED***

func (n *networkNamespace) setGateway(gw net.IP) ***REMOVED***
	n.Lock()
	n.gw = gw
	n.Unlock()
***REMOVED***

func (n *networkNamespace) setGatewayIPv6(gwv6 net.IP) ***REMOVED***
	n.Lock()
	n.gwv6 = gwv6
	n.Unlock()
***REMOVED***

func (n *networkNamespace) SetGateway(gw net.IP) error ***REMOVED***
	// Silently return if the gateway is empty
	if len(gw) == 0 ***REMOVED***
		return nil
	***REMOVED***

	err := n.programGateway(gw, true)
	if err == nil ***REMOVED***
		n.setGateway(gw)
	***REMOVED***

	return err
***REMOVED***

func (n *networkNamespace) UnsetGateway() error ***REMOVED***
	gw := n.Gateway()

	// Silently return if the gateway is empty
	if len(gw) == 0 ***REMOVED***
		return nil
	***REMOVED***

	err := n.programGateway(gw, false)
	if err == nil ***REMOVED***
		n.setGateway(net.IP***REMOVED******REMOVED***)
	***REMOVED***

	return err
***REMOVED***

func (n *networkNamespace) programGateway(gw net.IP, isAdd bool) error ***REMOVED***
	gwRoutes, err := n.nlHandle.RouteGet(gw)
	if err != nil ***REMOVED***
		return fmt.Errorf("route for the gateway %s could not be found: %v", gw, err)
	***REMOVED***

	var linkIndex int
	for _, gwRoute := range gwRoutes ***REMOVED***
		if gwRoute.Gw == nil ***REMOVED***
			linkIndex = gwRoute.LinkIndex
			break
		***REMOVED***
	***REMOVED***

	if linkIndex == 0 ***REMOVED***
		return fmt.Errorf("Direct route for the gateway %s could not be found", gw)
	***REMOVED***

	if isAdd ***REMOVED***
		return n.nlHandle.RouteAdd(&netlink.Route***REMOVED***
			Scope:     netlink.SCOPE_UNIVERSE,
			LinkIndex: linkIndex,
			Gw:        gw,
		***REMOVED***)
	***REMOVED***

	return n.nlHandle.RouteDel(&netlink.Route***REMOVED***
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: linkIndex,
		Gw:        gw,
	***REMOVED***)
***REMOVED***

// Program a route in to the namespace routing table.
func (n *networkNamespace) programRoute(path string, dest *net.IPNet, nh net.IP) error ***REMOVED***
	gwRoutes, err := n.nlHandle.RouteGet(nh)
	if err != nil ***REMOVED***
		return fmt.Errorf("route for the next hop %s could not be found: %v", nh, err)
	***REMOVED***

	return n.nlHandle.RouteAdd(&netlink.Route***REMOVED***
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: gwRoutes[0].LinkIndex,
		Gw:        nh,
		Dst:       dest,
	***REMOVED***)
***REMOVED***

// Delete a route from the namespace routing table.
func (n *networkNamespace) removeRoute(path string, dest *net.IPNet, nh net.IP) error ***REMOVED***
	gwRoutes, err := n.nlHandle.RouteGet(nh)
	if err != nil ***REMOVED***
		return fmt.Errorf("route for the next hop could not be found: %v", err)
	***REMOVED***

	return n.nlHandle.RouteDel(&netlink.Route***REMOVED***
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: gwRoutes[0].LinkIndex,
		Gw:        nh,
		Dst:       dest,
	***REMOVED***)
***REMOVED***

func (n *networkNamespace) SetGatewayIPv6(gwv6 net.IP) error ***REMOVED***
	// Silently return if the gateway is empty
	if len(gwv6) == 0 ***REMOVED***
		return nil
	***REMOVED***

	err := n.programGateway(gwv6, true)
	if err == nil ***REMOVED***
		n.setGatewayIPv6(gwv6)
	***REMOVED***

	return err
***REMOVED***

func (n *networkNamespace) UnsetGatewayIPv6() error ***REMOVED***
	gwv6 := n.GatewayIPv6()

	// Silently return if the gateway is empty
	if len(gwv6) == 0 ***REMOVED***
		return nil
	***REMOVED***

	err := n.programGateway(gwv6, false)
	if err == nil ***REMOVED***
		n.Lock()
		n.gwv6 = net.IP***REMOVED******REMOVED***
		n.Unlock()
	***REMOVED***

	return err
***REMOVED***

func (n *networkNamespace) AddStaticRoute(r *types.StaticRoute) error ***REMOVED***
	err := n.programRoute(n.nsPath(), r.Destination, r.NextHop)
	if err == nil ***REMOVED***
		n.Lock()
		n.staticRoutes = append(n.staticRoutes, r)
		n.Unlock()
	***REMOVED***
	return err
***REMOVED***

func (n *networkNamespace) RemoveStaticRoute(r *types.StaticRoute) error ***REMOVED***

	err := n.removeRoute(n.nsPath(), r.Destination, r.NextHop)
	if err == nil ***REMOVED***
		n.Lock()
		lastIndex := len(n.staticRoutes) - 1
		for i, v := range n.staticRoutes ***REMOVED***
			if v == r ***REMOVED***
				// Overwrite the route we're removing with the last element
				n.staticRoutes[i] = n.staticRoutes[lastIndex]
				// Shorten the slice to trim the extra element
				n.staticRoutes = n.staticRoutes[:lastIndex]
				break
			***REMOVED***
		***REMOVED***
		n.Unlock()
	***REMOVED***
	return err
***REMOVED***
