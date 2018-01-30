// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp_test

import (
	"log"
	"net"
	"os"
	"runtime"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
)

func ExamplePacketConn_nonPrivilegedPing() ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "darwin":
	case "linux":
		log.Println("you may need to adjust the net.ipv4.ping_group_range kernel state")
	default:
		log.Println("not supported on", runtime.GOOS)
		return
	***REMOVED***

	c, err := icmp.ListenPacket("udp6", "fe80::1%en0")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer c.Close()

	wm := icmp.Message***REMOVED***
		Type: ipv6.ICMPTypeEchoRequest, Code: 0,
		Body: &icmp.Echo***REMOVED***
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		***REMOVED***,
	***REMOVED***
	wb, err := wm.Marshal(nil)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	if _, err := c.WriteTo(wb, &net.UDPAddr***REMOVED***IP: net.ParseIP("ff02::1"), Zone: "en0"***REMOVED***); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	rb := make([]byte, 1500)
	n, peer, err := c.ReadFrom(rb)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	rm, err := icmp.ParseMessage(58, rb[:n])
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	switch rm.Type ***REMOVED***
	case ipv6.ICMPTypeEchoReply:
		log.Printf("got reflection from %v", peer)
	default:
		log.Printf("got %+v; want echo reply", rm)
	***REMOVED***
***REMOVED***
