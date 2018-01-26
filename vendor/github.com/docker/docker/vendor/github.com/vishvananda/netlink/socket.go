package netlink

import "net"

// SocketID identifies a single socket.
type SocketID struct ***REMOVED***
	SourcePort      uint16
	DestinationPort uint16
	Source          net.IP
	Destination     net.IP
	Interface       uint32
	Cookie          [2]uint32
***REMOVED***

// Socket represents a netlink socket.
type Socket struct ***REMOVED***
	Family  uint8
	State   uint8
	Timer   uint8
	Retrans uint8
	ID      SocketID
	Expires uint32
	RQueue  uint32
	WQueue  uint32
	UID     uint32
	INode   uint32
***REMOVED***
