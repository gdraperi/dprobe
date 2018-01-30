// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4_test

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func ExampleConn_markingTCP() ***REMOVED***
	ln, err := net.Listen("tcp", "0.0.0.0:1024")
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
			if c.RemoteAddr().(*net.TCPAddr).IP.To4() != nil ***REMOVED***
				p := ipv4.NewConn(c)
				if err := p.SetTOS(0x28); err != nil ***REMOVED*** // DSCP AF11
					log.Fatal(err)
				***REMOVED***
				if err := p.SetTTL(128); err != nil ***REMOVED***
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
	c, err := net.ListenPacket("udp4", "0.0.0.0:5353") // mDNS over UDP
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv4.NewPacketConn(c)

	en0, err := net.InterfaceByName("en0")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	mDNSLinkLocal := net.UDPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 251)***REMOVED***
	if err := p.JoinGroup(en0, &mDNSLinkLocal); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer p.LeaveGroup(en0, &mDNSLinkLocal)
	if err := p.SetControlMessage(ipv4.FlagDst, true); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	b := make([]byte, 1500)
	for ***REMOVED***
		_, cm, peer, err := p.ReadFrom(b)
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		if !cm.Dst.IsMulticast() || !cm.Dst.Equal(mDNSLinkLocal.IP) ***REMOVED***
			continue
		***REMOVED***
		answers := []byte("FAKE-MDNS-ANSWERS") // fake mDNS answers, you need to implement this
		if _, err := p.WriteTo(answers, nil, peer); err != nil ***REMOVED***
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
		if ip.To4() != nil ***REMOVED***
			dst.IP = ip
			fmt.Printf("using %v for tracing an IP packet route to %s\n", dst.IP, host)
			break
		***REMOVED***
	***REMOVED***
	if dst.IP == nil ***REMOVED***
		log.Fatal("no A record found")
	***REMOVED***

	c, err := net.ListenPacket("ip4:1", "0.0.0.0") // ICMP for IPv4
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer c.Close()
	p := ipv4.NewPacketConn(c)

	if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	wm := icmp.Message***REMOVED***
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo***REMOVED***
			ID:   os.Getpid() & 0xffff,
			Data: []byte("HELLO-R-U-THERE"),
		***REMOVED***,
	***REMOVED***

	rb := make([]byte, 1500)
	for i := 1; i <= 64; i++ ***REMOVED*** // up to 64 hops
		wm.Body.(*icmp.Echo).Seq = i
		wb, err := wm.Marshal(nil)
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		if err := p.SetTTL(i); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***

		// In the real world usually there are several
		// multiple traffic-engineered paths for each hop.
		// You may need to probe a few times to each hop.
		begin := time.Now()
		if _, err := p.WriteTo(wb, nil, &dst); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		n, cm, peer, err := p.ReadFrom(rb)
		if err != nil ***REMOVED***
			if err, ok := err.(net.Error); ok && err.Timeout() ***REMOVED***
				fmt.Printf("%v\t*\n", i)
				continue
			***REMOVED***
			log.Fatal(err)
		***REMOVED***
		rm, err := icmp.ParseMessage(1, rb[:n])
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		rtt := time.Since(begin)

		// In the real world you need to determine whether the
		// received message is yours using ControlMessage.Src,
		// ControlMessage.Dst, icmp.Echo.ID and icmp.Echo.Seq.
		switch rm.Type ***REMOVED***
		case ipv4.ICMPTypeTimeExceeded:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, cm)
		case ipv4.ICMPTypeEchoReply:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, cm)
			return
		default:
			log.Printf("unknown ICMP message: %+v\n", rm)
		***REMOVED***
	***REMOVED***
***REMOVED***

func ExampleRawConn_advertisingOSPFHello() ***REMOVED***
	c, err := net.ListenPacket("ip4:89", "0.0.0.0") // OSPF for IPv4
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer c.Close()
	r, err := ipv4.NewRawConn(c)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	en0, err := net.InterfaceByName("en0")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	allSPFRouters := net.IPAddr***REMOVED***IP: net.IPv4(224, 0, 0, 5)***REMOVED***
	if err := r.JoinGroup(en0, &allSPFRouters); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer r.LeaveGroup(en0, &allSPFRouters)

	hello := make([]byte, 24) // fake hello data, you need to implement this
	ospf := make([]byte, 24)  // fake ospf header, you need to implement this
	ospf[0] = 2               // version 2
	ospf[1] = 1               // hello packet
	ospf = append(ospf, hello...)
	iph := &ipv4.Header***REMOVED***
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      0xc0, // DSCP CS6
		TotalLen: ipv4.HeaderLen + len(ospf),
		TTL:      1,
		Protocol: 89,
		Dst:      allSPFRouters.IP.To4(),
	***REMOVED***

	var cm *ipv4.ControlMessage
	switch runtime.GOOS ***REMOVED***
	case "darwin", "linux":
		cm = &ipv4.ControlMessage***REMOVED***IfIndex: en0.Index***REMOVED***
	default:
		if err := r.SetMulticastInterface(en0); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
	***REMOVED***
	if err := r.WriteTo(iph, ospf, cm); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
***REMOVED***
