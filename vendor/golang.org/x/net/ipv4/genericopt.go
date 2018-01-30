// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import "syscall"

// TOS returns the type-of-service field value for outgoing packets.
func (c *genericOpt) TOS() (int, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoTOS]
	if !ok ***REMOVED***
		return 0, errOpNoSupport
	***REMOVED***
	return so.GetInt(c.Conn)
***REMOVED***

// SetTOS sets the type-of-service field value for future outgoing
// packets.
func (c *genericOpt) SetTOS(tos int) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoTOS]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.SetInt(c.Conn, tos)
***REMOVED***

// TTL returns the time-to-live field value for outgoing packets.
func (c *genericOpt) TTL() (int, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoTTL]
	if !ok ***REMOVED***
		return 0, errOpNoSupport
	***REMOVED***
	return so.GetInt(c.Conn)
***REMOVED***

// SetTTL sets the time-to-live field value for future outgoing
// packets.
func (c *genericOpt) SetTTL(ttl int) error ***REMOVED***
	if !c.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoTTL]
	if !ok ***REMOVED***
		return errOpNoSupport
	***REMOVED***
	return so.SetInt(c.Conn, ttl)
***REMOVED***
