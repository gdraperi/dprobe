// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import (
	"errors"
	"io"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

// RequestAgentForwarding sets up agent forwarding for the session.
// ForwardToAgent or ForwardToRemote should be called to route
// the authentication requests.
func RequestAgentForwarding(session *ssh.Session) error ***REMOVED***
	ok, err := session.SendRequest("auth-agent-req@openssh.com", true, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !ok ***REMOVED***
		return errors.New("forwarding request denied")
	***REMOVED***
	return nil
***REMOVED***

// ForwardToAgent routes authentication requests to the given keyring.
func ForwardToAgent(client *ssh.Client, keyring Agent) error ***REMOVED***
	channels := client.HandleChannelOpen(channelType)
	if channels == nil ***REMOVED***
		return errors.New("agent: already have handler for " + channelType)
	***REMOVED***

	go func() ***REMOVED***
		for ch := range channels ***REMOVED***
			channel, reqs, err := ch.Accept()
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			go ssh.DiscardRequests(reqs)
			go func() ***REMOVED***
				ServeAgent(keyring, channel)
				channel.Close()
			***REMOVED***()
		***REMOVED***
	***REMOVED***()
	return nil
***REMOVED***

const channelType = "auth-agent@openssh.com"

// ForwardToRemote routes authentication requests to the ssh-agent
// process serving on the given unix socket.
func ForwardToRemote(client *ssh.Client, addr string) error ***REMOVED***
	channels := client.HandleChannelOpen(channelType)
	if channels == nil ***REMOVED***
		return errors.New("agent: already have handler for " + channelType)
	***REMOVED***
	conn, err := net.Dial("unix", addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	conn.Close()

	go func() ***REMOVED***
		for ch := range channels ***REMOVED***
			channel, reqs, err := ch.Accept()
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			go ssh.DiscardRequests(reqs)
			go forwardUnixSocket(channel, addr)
		***REMOVED***
	***REMOVED***()
	return nil
***REMOVED***

func forwardUnixSocket(channel ssh.Channel, addr string) ***REMOVED***
	conn, err := net.Dial("unix", addr)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	var wg sync.WaitGroup
	wg.Add(2)
	go func() ***REMOVED***
		io.Copy(conn, channel)
		conn.(*net.UnixConn).CloseWrite()
		wg.Done()
	***REMOVED***()
	go func() ***REMOVED***
		io.Copy(channel, conn)
		channel.CloseWrite()
		wg.Done()
	***REMOVED***()

	wg.Wait()
	conn.Close()
	channel.Close()
***REMOVED***
