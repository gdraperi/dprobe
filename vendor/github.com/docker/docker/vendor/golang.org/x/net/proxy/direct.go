// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net"
)

type direct struct***REMOVED******REMOVED***

// Direct is a direct proxy: one that makes network connections directly.
var Direct = direct***REMOVED******REMOVED***

func (direct) Dial(network, addr string) (net.Conn, error) ***REMOVED***
	return net.Dial(network, addr)
***REMOVED***
