// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net"
	"strings"
)

// A PerHost directs connections to a default Dialer unless the hostname
// requested matches one of a number of exceptions.
type PerHost struct ***REMOVED***
	def, bypass Dialer

	bypassNetworks []*net.IPNet
	bypassIPs      []net.IP
	bypassZones    []string
	bypassHosts    []string
***REMOVED***

// NewPerHost returns a PerHost Dialer that directs connections to either
// defaultDialer or bypass, depending on whether the connection matches one of
// the configured rules.
func NewPerHost(defaultDialer, bypass Dialer) *PerHost ***REMOVED***
	return &PerHost***REMOVED***
		def:    defaultDialer,
		bypass: bypass,
	***REMOVED***
***REMOVED***

// Dial connects to the address addr on the given network through either
// defaultDialer or bypass.
func (p *PerHost) Dial(network, addr string) (c net.Conn, err error) ***REMOVED***
	host, _, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return p.dialerForRequest(host).Dial(network, addr)
***REMOVED***

func (p *PerHost) dialerForRequest(host string) Dialer ***REMOVED***
	if ip := net.ParseIP(host); ip != nil ***REMOVED***
		for _, net := range p.bypassNetworks ***REMOVED***
			if net.Contains(ip) ***REMOVED***
				return p.bypass
			***REMOVED***
		***REMOVED***
		for _, bypassIP := range p.bypassIPs ***REMOVED***
			if bypassIP.Equal(ip) ***REMOVED***
				return p.bypass
			***REMOVED***
		***REMOVED***
		return p.def
	***REMOVED***

	for _, zone := range p.bypassZones ***REMOVED***
		if strings.HasSuffix(host, zone) ***REMOVED***
			return p.bypass
		***REMOVED***
		if host == zone[1:] ***REMOVED***
			// For a zone "example.com", we match "example.com"
			// too.
			return p.bypass
		***REMOVED***
	***REMOVED***
	for _, bypassHost := range p.bypassHosts ***REMOVED***
		if bypassHost == host ***REMOVED***
			return p.bypass
		***REMOVED***
	***REMOVED***
	return p.def
***REMOVED***

// AddFromString parses a string that contains comma-separated values
// specifying hosts that should use the bypass proxy. Each value is either an
// IP address, a CIDR range, a zone (*.example.com) or a hostname
// (localhost). A best effort is made to parse the string and errors are
// ignored.
func (p *PerHost) AddFromString(s string) ***REMOVED***
	hosts := strings.Split(s, ",")
	for _, host := range hosts ***REMOVED***
		host = strings.TrimSpace(host)
		if len(host) == 0 ***REMOVED***
			continue
		***REMOVED***
		if strings.Contains(host, "/") ***REMOVED***
			// We assume that it's a CIDR address like 127.0.0.0/8
			if _, net, err := net.ParseCIDR(host); err == nil ***REMOVED***
				p.AddNetwork(net)
			***REMOVED***
			continue
		***REMOVED***
		if ip := net.ParseIP(host); ip != nil ***REMOVED***
			p.AddIP(ip)
			continue
		***REMOVED***
		if strings.HasPrefix(host, "*.") ***REMOVED***
			p.AddZone(host[1:])
			continue
		***REMOVED***
		p.AddHost(host)
	***REMOVED***
***REMOVED***

// AddIP specifies an IP address that will use the bypass proxy. Note that
// this will only take effect if a literal IP address is dialed. A connection
// to a named host will never match an IP.
func (p *PerHost) AddIP(ip net.IP) ***REMOVED***
	p.bypassIPs = append(p.bypassIPs, ip)
***REMOVED***

// AddNetwork specifies an IP range that will use the bypass proxy. Note that
// this will only take effect if a literal IP address is dialed. A connection
// to a named host will never match.
func (p *PerHost) AddNetwork(net *net.IPNet) ***REMOVED***
	p.bypassNetworks = append(p.bypassNetworks, net)
***REMOVED***

// AddZone specifies a DNS suffix that will use the bypass proxy. A zone of
// "example.com" matches "example.com" and all of its subdomains.
func (p *PerHost) AddZone(zone string) ***REMOVED***
	if strings.HasSuffix(zone, ".") ***REMOVED***
		zone = zone[:len(zone)-1]
	***REMOVED***
	if !strings.HasPrefix(zone, ".") ***REMOVED***
		zone = "." + zone
	***REMOVED***
	p.bypassZones = append(p.bypassZones, zone)
***REMOVED***

// AddHost specifies a hostname that will use the bypass proxy.
func (p *PerHost) AddHost(host string) ***REMOVED***
	if strings.HasSuffix(host, ".") ***REMOVED***
		host = host[:len(host)-1]
	***REMOVED***
	p.bypassHosts = append(p.bypassHosts, host)
***REMOVED***
