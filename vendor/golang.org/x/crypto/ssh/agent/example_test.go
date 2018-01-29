// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent_test

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func ExampleClientAgent() ***REMOVED***
	// ssh-agent has a UNIX socket under $SSH_AUTH_SOCK
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil ***REMOVED***
		log.Fatalf("net.Dial: %v", err)
	***REMOVED***
	agentClient := agent.NewClient(conn)
	config := &ssh.ClientConfig***REMOVED***
		User: "username",
		Auth: []ssh.AuthMethod***REMOVED***
			// Use a callback rather than PublicKeys
			// so we only consult the agent once the remote server
			// wants it.
			ssh.PublicKeysCallback(agentClient.Signers),
		***REMOVED***,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	***REMOVED***

	sshc, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil ***REMOVED***
		log.Fatalf("Dial: %v", err)
	***REMOVED***
	// .. use sshc
	sshc.Close()
***REMOVED***
