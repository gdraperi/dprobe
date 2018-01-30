// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"net"
	"syscall"

	"golang.org/x/net/bpf"
)

// MulticastTTL returns the time-to-live field value for outgoing
// multicast packets.
func (c *dgramOpt) MulticastTTL() (int, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoMulticastTTL]
	if !ok ***REMOVED***
		return 0, errOpNoSupport
	***REMOVED***
	return so.GetInt(c.Conn)
***REMOVED***

// SetMulticastTTL sets the time-to-live field value for future
// outgoing multicast packets.
func (c *dgramOpt) SetMulticastTTL(ttl int) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoMulticastTTL]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.SetInt(c.Conn, ttl)
***REMOVED***

// MulticastInterface returns the default interface for multicast
// packet transmissions.
func (c *dgramOpt) MulticastInterface() (*net.Interface, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return nil, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoMulticastInterface]
	if !ok ***REMOVED***
		return nil, errOpNoSupport
	***REMOVED***
	return so.getMulticastInterface(c.Conn)
***REMOVED***

// SetMulticastInterface sets the default interface for future
// multicast packet transmissions.
func (c *dgramOpt) SetMulticastInterface(ifi *net.Interface) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoMulticastInterface]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.setMulticastInterface(c.Conn, ifi)
***REMOVED***

// MulticastLoopback reports whether transmitted multicast packets
// should be copied and send back to the originator.
func (c *dgramOpt) MulticastLoopback() (bool, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return false, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoMulticastLoopback]
	if !ok ***REMOVED***
		return false, errOpNoSupport
	***REMOVED***
	on, err := so.GetInt(c.Conn)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return on == 1, nil
***REMOVED***

// SetMulticastLoopback sets whether transmitted multicast packets
// should be copied and send back to the originator.
func (c *dgramOpt) SetMulticastLoopback(on bool) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoMulticastLoopback]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.SetInt(c.Conn, boolint(on))
***REMOVED***

// JoinGroup joins the group address group on the interface ifi.
// By default all sources that can cast data to group are accepted.
// It's possible to mute and unmute data transmission from a specific
// source by using ExcludeSourceSpecificGroup and
// IncludeSourceSpecificGroup.
// JoinGroup uses the system assigned multicast interface when ifi is
// nil, although this is not recommended because the assignment
// depends on platforms and sometimes it might require routing
// configuration.
func (c *dgramOpt) JoinGroup(ifi *net.Interface, group net.Addr) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoJoinGroup]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	grp := netAddrToIP4(group)
	if grp == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	return so.setGroup(c.Conn, ifi, grp)
***REMOVED***

// LeaveGroup leaves the group address group on the interface ifi
// regardless of whether the group is any-source group or
// source-specific group.
func (c *dgramOpt) LeaveGroup(ifi *net.Interface, group net.Addr) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoLeaveGroup]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	grp := netAddrToIP4(group)
	if grp == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	return so.setGroup(c.Conn, ifi, grp)
***REMOVED***

// JoinSourceSpecificGroup joins the source-specific group comprising
// group and source on the interface ifi.
// JoinSourceSpecificGroup uses the system assigned multicast
// interface when ifi is nil, although this is not recommended because
// the assignment depends on platforms and sometimes it might require
// routing configuration.
func (c *dgramOpt) JoinSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoJoinSourceGroup]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	grp := netAddrToIP4(group)
	if grp == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	src := netAddrToIP4(source)
	if src == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	return so.setSourceGroup(c.Conn, ifi, grp, src)
***REMOVED***

// LeaveSourceSpecificGroup leaves the source-specific group on the
// interface ifi.
func (c *dgramOpt) LeaveSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoLeaveSourceGroup]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	grp := netAddrToIP4(group)
	if grp == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	src := netAddrToIP4(source)
	if src == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	return so.setSourceGroup(c.Conn, ifi, grp, src)
***REMOVED***

// ExcludeSourceSpecificGroup excludes the source-specific group from
// the already joined any-source groups by JoinGroup on the interface
// ifi.
func (c *dgramOpt) ExcludeSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoBlockSourceGroup]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	grp := netAddrToIP4(group)
	if grp == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	src := netAddrToIP4(source)
	if src == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	return so.setSourceGroup(c.Conn, ifi, grp, src)
***REMOVED***

// IncludeSourceSpecificGroup includes the excluded source-specific
// group by ExcludeSourceSpecificGroup again on the interface ifi.
func (c *dgramOpt) IncludeSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoUnblockSourceGroup]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	grp := netAddrToIP4(group)
	if grp == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	src := netAddrToIP4(source)
	if src == nil ***REMOVED***
		return errMissingAddress
	***REMOVED***
	return so.setSourceGroup(c.Conn, ifi, grp, src)
***REMOVED***

// ICMPFilter returns an ICMP filter.
// Currently only Linux supports this.
func (c *dgramOpt) ICMPFilter() (*ICMPFilter, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return nil, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoICMPFilter]
	if !ok ***REMOVED***
		return nil, errOpNoSupport
	***REMOVED***
	return so.getICMPFilter(c.Conn)
***REMOVED***

// SetICMPFilter deploys the ICMP filter.
// Currently only Linux supports this.
func (c *dgramOpt) SetICMPFilter(f *ICMPFilter) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoICMPFilter]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.setICMPFilter(c.Conn, f)
***REMOVED***

// SetBPF attaches a BPF program to the connection.
//
// Only supported on Linux.
func (c *dgramOpt) SetBPF(filter []bpf.RawInstruction) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoAttachFilter]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.setBPF(c.Conn, filter)
***REMOVED***
