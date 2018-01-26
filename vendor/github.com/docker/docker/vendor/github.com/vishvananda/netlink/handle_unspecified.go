// +build !linux

package netlink

import (
	"net"
	"time"

	"github.com/vishvananda/netns"
)

type Handle struct***REMOVED******REMOVED***

func NewHandle(nlFamilies ...int) (*Handle, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func NewHandleAt(ns netns.NsHandle, nlFamilies ...int) (*Handle, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func NewHandleAtFrom(newNs, curNs netns.NsHandle) (*Handle, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) Delete() ***REMOVED******REMOVED***

func (h *Handle) SupportsNetlinkFamily(nlFamily int) bool ***REMOVED***
	return false
***REMOVED***

func (h *Handle) SetSocketTimeout(to time.Duration) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) SetPromiscOn(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) SetPromiscOff(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetUp(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetDown(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetMTU(link Link, mtu int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetName(link Link, name string) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetAlias(link Link, name string) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetHardwareAddr(link Link, hwaddr net.HardwareAddr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetVfHardwareAddr(link Link, vf int, hwaddr net.HardwareAddr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetVfVlan(link Link, vf, vlan int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetVfTxRate(link Link, vf, rate int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetMaster(link Link, master *Bridge) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetNoMaster(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetMasterByIndex(link Link, masterIndex int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetNsPid(link Link, nspid int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetNsFd(link Link, fd int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkAdd(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkDel(link Link) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkByName(name string) (Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) LinkByAlias(alias string) (Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) LinkByIndex(index int) (Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) LinkList() ([]Link, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetHairpin(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetGuard(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetFastLeave(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetLearning(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetRootBlock(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetFlood(link Link, mode bool) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) LinkSetTxQLen(link Link, qlen int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) setProtinfoAttr(link Link, mode bool, attr int) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) AddrAdd(link Link, addr *Addr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) AddrDel(link Link, addr *Addr) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) AddrList(link Link, family int) ([]Addr, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) ClassDel(class Class) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) ClassChange(class Class) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) ClassReplace(class Class) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) ClassAdd(class Class) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) ClassList(link Link, parent uint32) ([]Class, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) FilterDel(filter Filter) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) FilterAdd(filter Filter) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) FilterList(link Link, parent uint32) ([]Filter, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) NeighAdd(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) NeighSet(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) NeighAppend(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) NeighDel(neigh *Neigh) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (h *Handle) NeighList(linkIndex, family int) ([]Neigh, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) NeighProxyList(linkIndex, family int) ([]Neigh, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***
