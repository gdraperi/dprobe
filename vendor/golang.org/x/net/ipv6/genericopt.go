// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import "syscall"

// TrafficClass returns the traffic class field value for outgoing
// packets.
func (c *genericOpt) TrafficClass() (int, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoTrafficClass]
	if !ok ***REMOVED***
		return 0, errOpNoSupport
	***REMOVED***
	return so.GetInt(c.Conn)
***REMOVED***

// SetTrafficClass sets the traffic class field value for future
// outgoing packets.
func (c *genericOpt) SetTrafficClass(tclass int) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoTrafficClass]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.SetInt(c.Conn, tclass)
***REMOVED***

// HopLimit returns the hop limit field value for outgoing packets.
func (c *genericOpt) HopLimit() (int, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoHopLimit]
	if !ok ***REMOVED***
		return 0, errOpNoSupport
	***REMOVED***
	return so.GetInt(c.Conn)
***REMOVED***

// SetHopLimit sets the hop limit field value for future outgoing
// packets.
func (c *genericOpt) SetHopLimit(hoplim int) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoHopLimit]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.SetInt(c.Conn, hoplim)
***REMOVED***
