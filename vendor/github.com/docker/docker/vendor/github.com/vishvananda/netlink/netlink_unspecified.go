// +build !linux

package netlink

import "net"

func LinkSetUp(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetDown(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetMTU(link Link, mtu int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetMaster(link Link, master *Bridge) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetNsPid(link Link, nspid int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetNsFd(link Link, fd int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetName(link Link, name string) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetAlias(link Link, name string) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetHardwareAddr(link Link, hwaddr net.HardwareAddr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetVfHardwareAddr(link Link, vf int, hwaddr net.HardwareAddr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetVfVlan(link Link, vf, vlan int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetVfTxRate(link Link, vf, rate int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetNoMaster(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetMasterByIndex(link Link, masterIndex int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetXdpFd(link Link, fd int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetARPOff(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetARPOn(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkByName(name string) (Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func LinkByAlias(alias string) (Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func LinkByIndex(index int) (Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func LinkSetHairpin(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetGuard(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetFastLeave(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetLearning(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetRootBlock(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetFlood(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkSetTxQLen(link Link, qlen int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkAdd(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkDel(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func SetHairpin(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func SetGuard(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func SetFastLeave(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func SetLearning(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func SetRootBlock(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func SetFlood(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func LinkList() ([]Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func AddrAdd(link Link, addr *Addr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func AddrDel(link Link, addr *Addr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func AddrList(link Link, family int) ([]Addr, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func RouteAdd(route *Route) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func RouteDel(route *Route) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func RouteList(link Link, family int) ([]Route, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func XfrmPolicyAdd(policy *XfrmPolicy) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func XfrmPolicyDel(policy *XfrmPolicy) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func XfrmPolicyList(family int) ([]XfrmPolicy, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func XfrmStateAdd(policy *XfrmState) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func XfrmStateDel(policy *XfrmState) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func XfrmStateList(family int) ([]XfrmState, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func NeighAdd(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func NeighSet(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func NeighAppend(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func NeighDel(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func NeighList(linkIndex, family int) ([]Neigh, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func NeighDeserialize(m []byte) (*Neigh, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func SocketGet(local, remote net.Addr) (*Socket, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***
