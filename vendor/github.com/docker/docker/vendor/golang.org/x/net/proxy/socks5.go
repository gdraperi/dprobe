// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"errors"
	"io"
	"net"
	"strconv"
)

// SOCKS5 returns a Dialer that makes SOCKSv5 connections to the given address
// with an optional username and password. See RFC 1928.
func SOCKS5(network, addr string, auth *Auth, forward Dialer) (Dialer, error) ***REMOVED***
	s := &socks5***REMOVED***
		network: network,
		addr:    addr,
		forward: forward,
	***REMOVED***
	if auth != nil ***REMOVED***
		s.user = auth.User
		s.password = auth.Password
	***REMOVED***

	return s, nil
***REMOVED***

type socks5 struct ***REMOVED***
	user, password string
	network, addr  string
	forward        Dialer
***REMOVED***

const socks5Version = 5

const (
	socks5AuthNone     = 0
	socks5AuthPassword = 2
)

const socks5Connect = 1

const (
	socks5IP4    = 1
	socks5Domain = 3
	socks5IP6    = 4
)

var socks5Errors = []string***REMOVED***
	"",
	"general failure",
	"connection forbidden",
	"network unreachable",
	"host unreachable",
	"connection refused",
	"TTL expired",
	"command not supported",
	"address type not supported",
***REMOVED***

// Dial connects to the address addr on the network net via the SOCKS5 proxy.
func (s *socks5) Dial(network, addr string) (net.Conn, error) ***REMOVED***
	switch network ***REMOVED***
	case "tcp", "tcp6", "tcp4":
	default:
		return nil, errors.New("proxy: no support for SOCKS5 proxy connections of type " + network)
	***REMOVED***

	conn, err := s.forward.Dial(s.network, s.addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := s.connect(conn, addr); err != nil ***REMOVED***
		conn.Close()
		return nil, err
	***REMOVED***
	return conn, nil
***REMOVED***

// connect takes an existing connection to a socks5 proxy server,
// and commands the server to extend that connection to target,
// which must be a canonical address with a host and port.
func (s *socks5) connect(conn net.Conn, target string) error ***REMOVED***
	host, portStr, err := net.SplitHostPort(target)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	port, err := strconv.Atoi(portStr)
	if err != nil ***REMOVED***
		return errors.New("proxy: failed to parse port number: " + portStr)
	***REMOVED***
	if port < 1 || port > 0xffff ***REMOVED***
		return errors.New("proxy: port number out of range: " + portStr)
	***REMOVED***

	// the size here is just an estimate
	buf := make([]byte, 0, 6+len(host))

	buf = append(buf, socks5Version)
	if len(s.user) > 0 && len(s.user) < 256 && len(s.password) < 256 ***REMOVED***
		buf = append(buf, 2 /* num auth methods */, socks5AuthNone, socks5AuthPassword)
	***REMOVED*** else ***REMOVED***
		buf = append(buf, 1 /* num auth methods */, socks5AuthNone)
	***REMOVED***

	if _, err := conn.Write(buf); err != nil ***REMOVED***
		return errors.New("proxy: failed to write greeting to SOCKS5 proxy at " + s.addr + ": " + err.Error())
	***REMOVED***

	if _, err := io.ReadFull(conn, buf[:2]); err != nil ***REMOVED***
		return errors.New("proxy: failed to read greeting from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	***REMOVED***
	if buf[0] != 5 ***REMOVED***
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " has unexpected version " + strconv.Itoa(int(buf[0])))
	***REMOVED***
	if buf[1] == 0xff ***REMOVED***
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " requires authentication")
	***REMOVED***

	if buf[1] == socks5AuthPassword ***REMOVED***
		buf = buf[:0]
		buf = append(buf, 1 /* password protocol version */)
		buf = append(buf, uint8(len(s.user)))
		buf = append(buf, s.user...)
		buf = append(buf, uint8(len(s.password)))
		buf = append(buf, s.password...)

		if _, err := conn.Write(buf); err != nil ***REMOVED***
			return errors.New("proxy: failed to write authentication request to SOCKS5 proxy at " + s.addr + ": " + err.Error())
		***REMOVED***

		if _, err := io.ReadFull(conn, buf[:2]); err != nil ***REMOVED***
			return errors.New("proxy: failed to read authentication reply from SOCKS5 proxy at " + s.addr + ": " + err.Error())
		***REMOVED***

		if buf[1] != 0 ***REMOVED***
			return errors.New("proxy: SOCKS5 proxy at " + s.addr + " rejected username/password")
		***REMOVED***
	***REMOVED***

	buf = buf[:0]
	buf = append(buf, socks5Version, socks5Connect, 0 /* reserved */)

	if ip := net.ParseIP(host); ip != nil ***REMOVED***
		if ip4 := ip.To4(); ip4 != nil ***REMOVED***
			buf = append(buf, socks5IP4)
			ip = ip4
		***REMOVED*** else ***REMOVED***
			buf = append(buf, socks5IP6)
		***REMOVED***
		buf = append(buf, ip...)
	***REMOVED*** else ***REMOVED***
		if len(host) > 255 ***REMOVED***
			return errors.New("proxy: destination hostname too long: " + host)
		***REMOVED***
		buf = append(buf, socks5Domain)
		buf = append(buf, byte(len(host)))
		buf = append(buf, host...)
	***REMOVED***
	buf = append(buf, byte(port>>8), byte(port))

	if _, err := conn.Write(buf); err != nil ***REMOVED***
		return errors.New("proxy: failed to write connect request to SOCKS5 proxy at " + s.addr + ": " + err.Error())
	***REMOVED***

	if _, err := io.ReadFull(conn, buf[:4]); err != nil ***REMOVED***
		return errors.New("proxy: failed to read connect reply from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	***REMOVED***

	failure := "unknown error"
	if int(buf[1]) < len(socks5Errors) ***REMOVED***
		failure = socks5Errors[buf[1]]
	***REMOVED***

	if len(failure) > 0 ***REMOVED***
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " failed to connect: " + failure)
	***REMOVED***

	bytesToDiscard := 0
	switch buf[3] ***REMOVED***
	case socks5IP4:
		bytesToDiscard = net.IPv4len
	case socks5IP6:
		bytesToDiscard = net.IPv6len
	case socks5Domain:
		_, err := io.ReadFull(conn, buf[:1])
		if err != nil ***REMOVED***
			return errors.New("proxy: failed to read domain length from SOCKS5 proxy at " + s.addr + ": " + err.Error())
		***REMOVED***
		bytesToDiscard = int(buf[0])
	default:
		return errors.New("proxy: got unknown address type " + strconv.Itoa(int(buf[3])) + " from SOCKS5 proxy at " + s.addr)
	***REMOVED***

	if cap(buf) < bytesToDiscard ***REMOVED***
		buf = make([]byte, bytesToDiscard)
	***REMOVED*** else ***REMOVED***
		buf = buf[:bytesToDiscard]
	***REMOVED***
	if _, err := io.ReadFull(conn, buf); err != nil ***REMOVED***
		return errors.New("proxy: failed to read address from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	***REMOVED***

	// Also need to discard the port number
	if _, err := io.ReadFull(conn, buf[:2]); err != nil ***REMOVED***
		return errors.New("proxy: failed to read port from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	***REMOVED***

	return nil
***REMOVED***
