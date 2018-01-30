// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package nettest provides utilities for network testing.
package nettest // import "golang.org/x/net/internal/nettest"

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"
)

var (
	supportsIPv4 bool
	supportsIPv6 bool
)

func init() ***REMOVED***
	if ln, err := net.Listen("tcp4", "127.0.0.1:0"); err == nil ***REMOVED***
		ln.Close()
		supportsIPv4 = true
	***REMOVED***
	if ln, err := net.Listen("tcp6", "[::1]:0"); err == nil ***REMOVED***
		ln.Close()
		supportsIPv6 = true
	***REMOVED***
***REMOVED***

// SupportsIPv4 reports whether the platform supports IPv4 networking
// functionality.
func SupportsIPv4() bool ***REMOVED*** return supportsIPv4 ***REMOVED***

// SupportsIPv6 reports whether the platform supports IPv6 networking
// functionality.
func SupportsIPv6() bool ***REMOVED*** return supportsIPv6 ***REMOVED***

// SupportsRawIPSocket reports whether the platform supports raw IP
// sockets.
func SupportsRawIPSocket() (string, bool) ***REMOVED***
	return supportsRawIPSocket()
***REMOVED***

// SupportsIPv6MulticastDeliveryOnLoopback reports whether the
// platform supports IPv6 multicast packet delivery on software
// loopback interface.
func SupportsIPv6MulticastDeliveryOnLoopback() bool ***REMOVED***
	return supportsIPv6MulticastDeliveryOnLoopback()
***REMOVED***

// ProtocolNotSupported reports whether err is a protocol not
// supported error.
func ProtocolNotSupported(err error) bool ***REMOVED***
	return protocolNotSupported(err)
***REMOVED***

// TestableNetwork reports whether network is testable on the current
// platform configuration.
func TestableNetwork(network string) bool ***REMOVED***
	// This is based on logic from standard library's
	// net/platform_test.go.
	switch network ***REMOVED***
	case "unix", "unixgram":
		switch runtime.GOOS ***REMOVED***
		case "android", "nacl", "plan9", "windows":
			return false
		***REMOVED***
		if runtime.GOOS == "darwin" && (runtime.GOARCH == "arm" || runtime.GOARCH == "arm64") ***REMOVED***
			return false
		***REMOVED***
	case "unixpacket":
		switch runtime.GOOS ***REMOVED***
		case "android", "darwin", "freebsd", "nacl", "plan9", "windows":
			return false
		case "netbsd":
			// It passes on amd64 at least. 386 fails (Issue 22927). arm is unknown.
			if runtime.GOARCH == "386" ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// NewLocalListener returns a listener which listens to a loopback IP
// address or local file system path.
// Network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
func NewLocalListener(network string) (net.Listener, error) ***REMOVED***
	switch network ***REMOVED***
	case "tcp":
		if supportsIPv4 ***REMOVED***
			if ln, err := net.Listen("tcp4", "127.0.0.1:0"); err == nil ***REMOVED***
				return ln, nil
			***REMOVED***
		***REMOVED***
		if supportsIPv6 ***REMOVED***
			return net.Listen("tcp6", "[::1]:0")
		***REMOVED***
	case "tcp4":
		if supportsIPv4 ***REMOVED***
			return net.Listen("tcp4", "127.0.0.1:0")
		***REMOVED***
	case "tcp6":
		if supportsIPv6 ***REMOVED***
			return net.Listen("tcp6", "[::1]:0")
		***REMOVED***
	case "unix", "unixpacket":
		return net.Listen(network, localPath())
	***REMOVED***
	return nil, fmt.Errorf("%s is not supported", network)
***REMOVED***

// NewLocalPacketListener returns a packet listener which listens to a
// loopback IP address or local file system path.
// Network must be "udp", "udp4", "udp6" or "unixgram".
func NewLocalPacketListener(network string) (net.PacketConn, error) ***REMOVED***
	switch network ***REMOVED***
	case "udp":
		if supportsIPv4 ***REMOVED***
			if c, err := net.ListenPacket("udp4", "127.0.0.1:0"); err == nil ***REMOVED***
				return c, nil
			***REMOVED***
		***REMOVED***
		if supportsIPv6 ***REMOVED***
			return net.ListenPacket("udp6", "[::1]:0")
		***REMOVED***
	case "udp4":
		if supportsIPv4 ***REMOVED***
			return net.ListenPacket("udp4", "127.0.0.1:0")
		***REMOVED***
	case "udp6":
		if supportsIPv6 ***REMOVED***
			return net.ListenPacket("udp6", "[::1]:0")
		***REMOVED***
	case "unixgram":
		return net.ListenPacket(network, localPath())
	***REMOVED***
	return nil, fmt.Errorf("%s is not supported", network)
***REMOVED***

func localPath() string ***REMOVED***
	f, err := ioutil.TempFile("", "nettest")
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	path := f.Name()
	f.Close()
	os.Remove(path)
	return path
***REMOVED***
