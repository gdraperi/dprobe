package osl

import "net"

func (nh *neigh) processNeighOptions(options ...NeighOption) ***REMOVED***
	for _, opt := range options ***REMOVED***
		if opt != nil ***REMOVED***
			opt(nh)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *networkNamespace) LinkName(name string) NeighOption ***REMOVED***
	return func(nh *neigh) ***REMOVED***
		nh.linkName = name
	***REMOVED***
***REMOVED***

func (n *networkNamespace) Family(family int) NeighOption ***REMOVED***
	return func(nh *neigh) ***REMOVED***
		nh.family = family
	***REMOVED***
***REMOVED***

func (i *nwIface) processInterfaceOptions(options ...IfaceOption) ***REMOVED***
	for _, opt := range options ***REMOVED***
		if opt != nil ***REMOVED***
			opt(i)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *networkNamespace) Bridge(isBridge bool) IfaceOption ***REMOVED***
	return func(i *nwIface) ***REMOVED***
		i.bridge = isBridge
	***REMOVED***
***REMOVED***

func (n *networkNamespace) Master(name string) IfaceOption ***REMOVED***
	return func(i *nwIface) ***REMOVED***
		i.master = name
	***REMOVED***
***REMOVED***

func (n *networkNamespace) MacAddress(mac net.HardwareAddr) IfaceOption ***REMOVED***
	return func(i *nwIface) ***REMOVED***
		i.mac = mac
	***REMOVED***
***REMOVED***

func (n *networkNamespace) Address(addr *net.IPNet) IfaceOption ***REMOVED***
	return func(i *nwIface) ***REMOVED***
		i.address = addr
	***REMOVED***
***REMOVED***

func (n *networkNamespace) AddressIPv6(addr *net.IPNet) IfaceOption ***REMOVED***
	return func(i *nwIface) ***REMOVED***
		i.addressIPv6 = addr
	***REMOVED***
***REMOVED***

func (n *networkNamespace) LinkLocalAddresses(list []*net.IPNet) IfaceOption ***REMOVED***
	return func(i *nwIface) ***REMOVED***
		i.llAddrs = list
	***REMOVED***
***REMOVED***

func (n *networkNamespace) Routes(routes []*net.IPNet) IfaceOption ***REMOVED***
	return func(i *nwIface) ***REMOVED***
		i.routes = routes
	***REMOVED***
***REMOVED***
