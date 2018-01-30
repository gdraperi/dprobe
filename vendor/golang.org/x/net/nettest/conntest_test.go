// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package nettest

import (
	"net"
	"os"
	"runtime"
	"testing"

	"golang.org/x/net/internal/nettest"
)

func TestTestConn(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** name, network string ***REMOVED******REMOVED***
		***REMOVED***"TCP", "tcp"***REMOVED***,
		***REMOVED***"UnixPipe", "unix"***REMOVED***,
		***REMOVED***"UnixPacketPipe", "unixpacket"***REMOVED***,
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			if !nettest.TestableNetwork(tt.network) ***REMOVED***
				t.Skipf("not supported on %s", runtime.GOOS)
			***REMOVED***

			mp := func() (c1, c2 net.Conn, stop func(), err error) ***REMOVED***
				ln, err := nettest.NewLocalListener(tt.network)
				if err != nil ***REMOVED***
					return nil, nil, nil, err
				***REMOVED***

				// Start a connection between two endpoints.
				var err1, err2 error
				done := make(chan bool)
				go func() ***REMOVED***
					c2, err2 = ln.Accept()
					close(done)
				***REMOVED***()
				c1, err1 = net.Dial(ln.Addr().Network(), ln.Addr().String())
				<-done

				stop = func() ***REMOVED***
					if err1 == nil ***REMOVED***
						c1.Close()
					***REMOVED***
					if err2 == nil ***REMOVED***
						c2.Close()
					***REMOVED***
					ln.Close()
					switch tt.network ***REMOVED***
					case "unix", "unixpacket":
						os.Remove(ln.Addr().String())
					***REMOVED***
				***REMOVED***

				switch ***REMOVED***
				case err1 != nil:
					stop()
					return nil, nil, nil, err1
				case err2 != nil:
					stop()
					return nil, nil, nil, err2
				default:
					return c1, c2, stop, nil
				***REMOVED***
			***REMOVED***

			TestConn(t, mp)
		***REMOVED***)
	***REMOVED***
***REMOVED***
