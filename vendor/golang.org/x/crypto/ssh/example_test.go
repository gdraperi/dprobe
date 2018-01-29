// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func ExampleNewServerConn() ***REMOVED***
	// Public key authentication is done by comparing
	// the public key of a received connection
	// with the entries in the authorized_keys file.
	authorizedKeysBytes, err := ioutil.ReadFile("authorized_keys")
	if err != nil ***REMOVED***
		log.Fatalf("Failed to load authorized_keys, err: %v", err)
	***REMOVED***

	authorizedKeysMap := map[string]bool***REMOVED******REMOVED***
	for len(authorizedKeysBytes) > 0 ***REMOVED***
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***

		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeysBytes = rest
	***REMOVED***

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig***REMOVED***
		// Remove to disable password auth.
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) ***REMOVED***
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			if c.User() == "testuser" && string(pass) == "tiger" ***REMOVED***
				return nil, nil
			***REMOVED***
			return nil, fmt.Errorf("password rejected for %q", c.User())
		***REMOVED***,

		// Remove to disable public key auth.
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) ***REMOVED***
			if authorizedKeysMap[string(pubKey.Marshal())] ***REMOVED***
				return &ssh.Permissions***REMOVED***
					// Record the public key used for authentication.
					Extensions: map[string]string***REMOVED***
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					***REMOVED***,
				***REMOVED***, nil
			***REMOVED***
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		***REMOVED***,
	***REMOVED***

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil ***REMOVED***
		log.Fatal("Failed to load private key: ", err)
	***REMOVED***

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil ***REMOVED***
		log.Fatal("Failed to parse private key: ", err)
	***REMOVED***

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil ***REMOVED***
		log.Fatal("failed to listen for connection: ", err)
	***REMOVED***
	nConn, err := listener.Accept()
	if err != nil ***REMOVED***
		log.Fatal("failed to accept incoming connection: ", err)
	***REMOVED***

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil ***REMOVED***
		log.Fatal("failed to handshake: ", err)
	***REMOVED***
	log.Printf("logged in with key %s", conn.Permissions.Extensions["pubkey-fp"])

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans ***REMOVED***
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" ***REMOVED***
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		***REMOVED***
		channel, requests, err := newChannel.Accept()
		if err != nil ***REMOVED***
			log.Fatalf("Could not accept channel: %v", err)
		***REMOVED***

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "shell" request.
		go func(in <-chan *ssh.Request) ***REMOVED***
			for req := range in ***REMOVED***
				req.Reply(req.Type == "shell", nil)
			***REMOVED***
		***REMOVED***(requests)

		term := terminal.NewTerminal(channel, "> ")

		go func() ***REMOVED***
			defer channel.Close()
			for ***REMOVED***
				line, err := term.ReadLine()
				if err != nil ***REMOVED***
					break
				***REMOVED***
				fmt.Println(line)
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

func ExampleHostKeyCheck() ***REMOVED***
	// Every client must provide a host key check.  Here is a
	// simple-minded parse of OpenSSH's known_hosts file
	host := "hostname"
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() ***REMOVED***
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 ***REMOVED***
			continue
		***REMOVED***
		if strings.Contains(fields[0], host) ***REMOVED***
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil ***REMOVED***
				log.Fatalf("error parsing %q: %v", fields[2], err)
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if hostKey == nil ***REMOVED***
		log.Fatalf("no hostkey for %s", host)
	***REMOVED***

	config := ssh.ClientConfig***REMOVED***
		User:            os.Getenv("USER"),
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	***REMOVED***

	_, err = ssh.Dial("tcp", host+":22", &config)
	log.Println(err)
***REMOVED***

func ExampleDial() ***REMOVED***
	var hostKey ssh.PublicKey
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig***REMOVED***
		User: "username",
		Auth: []ssh.AuthMethod***REMOVED***
			ssh.Password("yourpassword"),
		***REMOVED***,
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	***REMOVED***
	client, err := ssh.Dial("tcp", "yourserver.com:22", config)
	if err != nil ***REMOVED***
		log.Fatal("Failed to dial: ", err)
	***REMOVED***

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil ***REMOVED***
		log.Fatal("Failed to create session: ", err)
	***REMOVED***
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil ***REMOVED***
		log.Fatal("Failed to run: " + err.Error())
	***REMOVED***
	fmt.Println(b.String())
***REMOVED***

func ExamplePublicKeys() ***REMOVED***
	var hostKey ssh.PublicKey
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile("/home/user/.ssh/id_rsa")
	if err != nil ***REMOVED***
		log.Fatalf("unable to read private key: %v", err)
	***REMOVED***

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil ***REMOVED***
		log.Fatalf("unable to parse private key: %v", err)
	***REMOVED***

	config := &ssh.ClientConfig***REMOVED***
		User: "user",
		Auth: []ssh.AuthMethod***REMOVED***
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		***REMOVED***,
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	***REMOVED***

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", "host.com:22", config)
	if err != nil ***REMOVED***
		log.Fatalf("unable to connect: %v", err)
	***REMOVED***
	defer client.Close()
***REMOVED***

func ExampleClient_Listen() ***REMOVED***
	var hostKey ssh.PublicKey
	config := &ssh.ClientConfig***REMOVED***
		User: "username",
		Auth: []ssh.AuthMethod***REMOVED***
			ssh.Password("password"),
		***REMOVED***,
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	***REMOVED***
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil ***REMOVED***
		log.Fatal("unable to connect: ", err)
	***REMOVED***
	defer conn.Close()

	// Request the remote side to open port 8080 on all interfaces.
	l, err := conn.Listen("tcp", "0.0.0.0:8080")
	if err != nil ***REMOVED***
		log.Fatal("unable to register tcp forward: ", err)
	***REMOVED***
	defer l.Close()

	// Serve HTTP with your SSH server acting as a reverse proxy.
	http.Serve(l, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) ***REMOVED***
		fmt.Fprintf(resp, "Hello world!\n")
	***REMOVED***))
***REMOVED***

func ExampleSession_RequestPty() ***REMOVED***
	var hostKey ssh.PublicKey
	// Create client config
	config := &ssh.ClientConfig***REMOVED***
		User: "username",
		Auth: []ssh.AuthMethod***REMOVED***
			ssh.Password("password"),
		***REMOVED***,
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	***REMOVED***
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil ***REMOVED***
		log.Fatal("unable to connect: ", err)
	***REMOVED***
	defer conn.Close()
	// Create a session
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		log.Fatal("unable to create session: ", err)
	***REMOVED***
	defer session.Close()
	// Set up terminal modes
	modes := ssh.TerminalModes***REMOVED***
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	***REMOVED***
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil ***REMOVED***
		log.Fatal("request for pseudo terminal failed: ", err)
	***REMOVED***
	// Start remote shell
	if err := session.Shell(); err != nil ***REMOVED***
		log.Fatal("failed to start shell: ", err)
	***REMOVED***
***REMOVED***
