package bridge

import (
	"fmt"
	"math/rand"
	"net"
	"syscall"
	"time"
	"unsafe"

	"github.com/docker/libnetwork/netutils"
)

const (
	ifNameSize   = 16
	ioctlBrAdd   = 0x89a0
	ioctlBrAddIf = 0x89a2
)

type ifreqIndex struct ***REMOVED***
	IfrnName  [ifNameSize]byte
	IfruIndex int32
***REMOVED***

type ifreqHwaddr struct ***REMOVED***
	IfrnName   [ifNameSize]byte
	IfruHwaddr syscall.RawSockaddr
***REMOVED***

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// THIS CODE DOES NOT COMMUNICATE WITH KERNEL VIA RTNETLINK INTERFACE
// IT IS HERE FOR BACKWARDS COMPATIBILITY WITH OLDER LINUX KERNELS
// WHICH SHIP WITH OLDER NOT ENTIRELY FUNCTIONAL VERSION OF NETLINK
func getIfSocket() (fd int, err error) ***REMOVED***
	for _, socket := range []int***REMOVED***
		syscall.AF_INET,
		syscall.AF_PACKET,
		syscall.AF_INET6,
	***REMOVED*** ***REMOVED***
		if fd, err = syscall.Socket(socket, syscall.SOCK_DGRAM, 0); err == nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if err == nil ***REMOVED***
		return fd, nil
	***REMOVED***
	return -1, err
***REMOVED***

func ifIoctBridge(iface, master *net.Interface, op uintptr) error ***REMOVED***
	if len(master.Name) >= ifNameSize ***REMOVED***
		return fmt.Errorf("Interface name %s too long", master.Name)
	***REMOVED***

	s, err := getIfSocket()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer syscall.Close(s)

	ifr := ifreqIndex***REMOVED******REMOVED***
	copy(ifr.IfrnName[:len(ifr.IfrnName)-1], master.Name)
	ifr.IfruIndex = int32(iface.Index)

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s), op, uintptr(unsafe.Pointer(&ifr))); err != 0 ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// Add a slave to a bridge device.  This is more backward-compatible than
// netlink.NetworkSetMaster and works on RHEL 6.
func ioctlAddToBridge(iface, master *net.Interface) error ***REMOVED***
	return ifIoctBridge(iface, master, ioctlBrAddIf)
***REMOVED***

func ioctlSetMacAddress(name, addr string) error ***REMOVED***
	if len(name) >= ifNameSize ***REMOVED***
		return fmt.Errorf("Interface name %s too long", name)
	***REMOVED***

	hw, err := net.ParseMAC(addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	s, err := getIfSocket()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer syscall.Close(s)

	ifr := ifreqHwaddr***REMOVED******REMOVED***
	ifr.IfruHwaddr.Family = syscall.ARPHRD_ETHER
	copy(ifr.IfrnName[:len(ifr.IfrnName)-1], name)

	for i := 0; i < 6; i++ ***REMOVED***
		ifr.IfruHwaddr.Data[i] = ifrDataByte(hw[i])
	***REMOVED***

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s), syscall.SIOCSIFHWADDR, uintptr(unsafe.Pointer(&ifr))); err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func ioctlCreateBridge(name string, setMacAddr bool) error ***REMOVED***
	if len(name) >= ifNameSize ***REMOVED***
		return fmt.Errorf("Interface name %s too long", name)
	***REMOVED***

	s, err := getIfSocket()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer syscall.Close(s)

	nameBytePtr, err := syscall.BytePtrFromString(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s), ioctlBrAdd, uintptr(unsafe.Pointer(nameBytePtr))); err != 0 ***REMOVED***
		return err
	***REMOVED***
	if setMacAddr ***REMOVED***
		return ioctlSetMacAddress(name, netutils.GenerateRandomMAC().String())
	***REMOVED***
	return nil
***REMOVED***
