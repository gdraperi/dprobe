// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

// An InterfaceMessage represents an interface message.
type InterfaceMessage struct ***REMOVED***
	Version int    // message version
	Type    int    // message type
	Flags   int    // interface flags
	Index   int    // interface index
	Name    string // interface name
	Addrs   []Addr // addresses

	extOff int    // offset of header extension
	raw    []byte // raw message
***REMOVED***

// An InterfaceAddrMessage represents an interface address message.
type InterfaceAddrMessage struct ***REMOVED***
	Version int    // message version
	Type    int    // message type
	Flags   int    // interface flags
	Index   int    // interface index
	Addrs   []Addr // addresses

	raw []byte // raw message
***REMOVED***

// Sys implements the Sys method of Message interface.
func (m *InterfaceAddrMessage) Sys() []Sys ***REMOVED*** return nil ***REMOVED***

// An InterfaceMulticastAddrMessage represents an interface multicast
// address message.
type InterfaceMulticastAddrMessage struct ***REMOVED***
	Version int    // message version
	Type    int    // messsage type
	Flags   int    // interface flags
	Index   int    // interface index
	Addrs   []Addr // addresses

	raw []byte // raw message
***REMOVED***

// Sys implements the Sys method of Message interface.
func (m *InterfaceMulticastAddrMessage) Sys() []Sys ***REMOVED*** return nil ***REMOVED***

// An InterfaceAnnounceMessage represents an interface announcement
// message.
type InterfaceAnnounceMessage struct ***REMOVED***
	Version int    // message version
	Type    int    // message type
	Index   int    // interface index
	Name    string // interface name
	What    int    // what type of announcement

	raw []byte // raw message
***REMOVED***

// Sys implements the Sys method of Message interface.
func (m *InterfaceAnnounceMessage) Sys() []Sys ***REMOVED*** return nil ***REMOVED***
