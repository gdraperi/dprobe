// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import "golang.org/x/net/internal/iana"

// BUG(mikio): On Windows, methods related to ICMPFilter are not
// implemented.

// An ICMPType represents a type of ICMP message.
type ICMPType int

func (typ ICMPType) String() string ***REMOVED***
	s, ok := icmpTypes[typ]
	if !ok ***REMOVED***
		return "<nil>"
	***REMOVED***
	return s
***REMOVED***

// Protocol returns the ICMPv6 protocol number.
func (typ ICMPType) Protocol() int ***REMOVED***
	return iana.ProtocolIPv6ICMP
***REMOVED***

// An ICMPFilter represents an ICMP message filter for incoming
// packets. The filter belongs to a packet delivery path on a host and
// it cannot interact with forwarding packets or tunnel-outer packets.
//
// Note: RFC 8200 defines a reasonable role model. A node means a
// device that implements IP. A router means a node that forwards IP
// packets not explicitly addressed to itself, and a host means a node
// that is not a router.
type ICMPFilter struct ***REMOVED***
	icmpv6Filter
***REMOVED***

// Accept accepts incoming ICMP packets including the type field value
// typ.
func (f *ICMPFilter) Accept(typ ICMPType) ***REMOVED***
	f.accept(typ)
***REMOVED***

// Block blocks incoming ICMP packets including the type field value
// typ.
func (f *ICMPFilter) Block(typ ICMPType) ***REMOVED***
	f.block(typ)
***REMOVED***

// SetAll sets the filter action to the filter.
func (f *ICMPFilter) SetAll(block bool) ***REMOVED***
	f.setAll(block)
***REMOVED***

// WillBlock reports whether the ICMP type will be blocked.
func (f *ICMPFilter) WillBlock(typ ICMPType) bool ***REMOVED***
	return f.willBlock(typ)
***REMOVED***
