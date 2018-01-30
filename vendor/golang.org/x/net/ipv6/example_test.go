// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
)

func ExampleConn_markingTCP() ***REMOVED***
	ln, err := net.Listen("tcp", "[::]:1024")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer ln.Close()

	for ***REMOVED***
		c, err := ln.Accept()
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		go func(c net.Conn) ***REMOVED***
			defer c.Close()
			if c.RemoteAddr().(*net.TCPAddr).IP.To16() != nil && c.RemoteAddr().(*net.TCPAddr).IP.To4() == nil ***REMOVED***
				p := ipv6.NewConn(c)
				if err := p.SetTrafficClass(0x28); err != nil ***REMOVED*** // DSCP AF11
					log.Fatal(err)
				***REMOVED***
				if err := p.SetHopLimit(128); err != nil ***REMOVED***
					log.Fatal(err)
				***REMOVED***
			***REMOVED***
			if _, err := c.Write([]byte("HELLO-R-U-THERE-ACK")); err != nil ***REMOVED***
				log.Fatal(err)
			***REMOVED***
		***REMOVED***(c)
	***REMOVED***
***REMOVED***

func ExamplePacketConn_servingOneShotMulticastDNS() ***REMOVED***
	c, err := net.ListenPacket("udp6", "[::]:5353") // mDNS over UDP
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv6.NewPacketConn(c)

	en0, err := net.InterfaceByName("en0")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	mDNSLinkLocal := net.UDPAddr***REMOVED***IP: net.ParseIP("ff02::fb")***REMOVED***
	if err := p.JoinGroup(en0, &mDNSLinkLocal); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer p.LeaveGroup(en0, &mDNSLinkLocal)
	if err := p.SetControlMessage(ipv6.FlagDst|ipv6.FlagInterface, true); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	var wcm ipv6.ControlMessage
	b := make([]byte, 1500)
	for ***REMOVED***
		_, rcm, peer, err := p.ReadFrom(b)
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		if !rcm.Dst.IsMulticast() || !rcm.Dst.Equal(mDNSLinkLocal.IP) ***REMOVED***
			continue
		***REMOVED***
		wcm.IfIndex = rcm.IfIndex
		answers := []byte("FAKE-MDNS-ANSWERS") // fake mDNS answers, you need to implement this
		if _, err := p.WriteTo(answers, &wcm, peer); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func ExamplePacketConn_tracingIPPacketRoute() ***REMOVED***
	// Tracing an IP packet route to www.google.com.

	const host = "www.google.com"
	ips, err := net.LookupIP(host)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	var dst net.IPAddr
	for _, ip := range ips ***REMOVED***
		if ip.To16() != nil && ip.To4() == nil ***REMOVED***
			dst.IP = ip
			fmt.Printf("using %v for tracing an IP packet route to %s\n", dst.IP, host)
			break
		***REMOVED***
	***REMOVED***
	if dst.IP == nil ***REMOVED***
		log.Fatal("no AAAA record found")
	***REMOVED***

	c, err := net.ListenPacket("ip6:58", "::") // ICMP for IPv6
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv6.NewPacketConn(c)

	if err := p.SetControlMessage(ipv6.FlagHopLimit|ipv6.FlagSrc|ipv6.FlagDst|ipv6.FlagInterface, true); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	wm := icmp.Message***REMOVED***
		Type: ipv6.ICMPTypeEchoRequest, Code: 0,
		Body: &icmp.Echo***REMOVED***
			ID:   os.Getpid() & 0xffff,
			Data: []byte("HELLO-R-U-THERE"),
		***REMOVED***,
	***REMOVED***
	var f ipv6.ICMPFilter
	f.SetAll(true)
	f.Accept(ipv6.ICMPTypeTimeExceeded)
	f.Accept(ipv6.ICMPTypeEchoReply)
	if err := p.SetICMPFilter(&f); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	var wcm ipv6.ControlMessage
	rb := make([]byte, 1500)
	for i := 1; i <= 64; i++ ***REMOVED*** // up to 64 hops
		wm.Body.(*icmp.Echo).Seq = i
		wb, err := wm.Marshal(nil)
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***

		// In the real world usually there are several
		// multiple traffic-engineered paths for each hop.
		// You may need to probe a few times to each hop.
		begin := time.Now()
		wcm.HopLimit = i
		if _, err := p.WriteTo(wb, &wcm, &dst); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		n, rcm, peer, err := p.ReadFrom(rb)
		if err != nil ***REMOVED***
			if err, ok := err.(net.Error); ok && err.Timeout() ***REMOVED***
				fmt.Printf("%v\t*\n", i)
				continue
			***REMOVED***
			log.Fatal(err)
		***REMOVED***
		rm, err := icmp.ParseMessage(58, rb[:n])
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		rtt := time.Since(begin)

		// In the real world you need to determine whether the
		// received message is yours using ControlMessage.Src,
		// ControlMesage.Dst, icmp.Echo.ID and icmp.Echo.Seq.
		switch rm.Type ***REMOVED***
		case ipv6.ICMPTypeTimeExceeded:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, rcm)
		case ipv6.ICMPTypeEchoReply:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, rcm)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func ExamplePacketConn_advertisingOSPFHello() ***REMOVED***
	c, err := net.ListenPacket("ip6:89", "::") // OSPF for IPv6
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv6.NewPacketConn(c)

	en0, err := net.InterfaceByName("en0")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	allSPFRouters := net.IPAddr***REMOVED***IP: net.ParseIP("ff02::5")***REMOVED***
	if err := p.JoinGroup(en0, &allSPFRouters); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer p.LeaveGroup(en0, &allSPFRouters)

	hello := make([]byte, 24) // fake hello data, you need to implement this
	ospf := make([]byte, 16)  // fake ospf header, you need to implement this
	ospf[0] = 3               // version 3
	ospf[1] = 1               // hello packet
	ospf = append(ospf, hello...)
	if err := p.SetChecksum(true, 12); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	cm := ipv6.ControlMessage***REMOVED***
		TrafficClass: 0xc0, // DSCP CS6
		HopLimit:     1,
		IfIndex:      en0.Index,
	***REMOVED***
	if _, err := p.WriteTo(ospf, &cm, &allSPFRouters); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
***REMOVED***
