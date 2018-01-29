// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"strings"
	"testing"
)

func TestClientVersion(t *testing.T) ***REMOVED***
	for _, tt := range []struct ***REMOVED***
		name      string
		version   string
		multiLine string
		wantErr   bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:    "default version",
			version: packageVersion,
		***REMOVED***,
		***REMOVED***
			name:    "custom version",
			version: "SSH-2.0-CustomClientVersionString",
		***REMOVED***,
		***REMOVED***
			name:      "good multi line version",
			version:   packageVersion,
			multiLine: strings.Repeat("ignored\r\n", 20),
		***REMOVED***,
		***REMOVED***
			name:      "bad multi line version",
			version:   packageVersion,
			multiLine: "bad multi line version",
			wantErr:   true,
		***REMOVED***,
		***REMOVED***
			name:      "long multi line version",
			version:   packageVersion,
			multiLine: strings.Repeat("long multi line version\r\n", 50)[:256],
			wantErr:   true,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			c1, c2, err := netPipe()
			if err != nil ***REMOVED***
				t.Fatalf("netPipe: %v", err)
			***REMOVED***
			defer c1.Close()
			defer c2.Close()
			go func() ***REMOVED***
				if tt.multiLine != "" ***REMOVED***
					c1.Write([]byte(tt.multiLine))
				***REMOVED***
				NewClientConn(c1, "", &ClientConfig***REMOVED***
					ClientVersion:   tt.version,
					HostKeyCallback: InsecureIgnoreHostKey(),
				***REMOVED***)
				c1.Close()
			***REMOVED***()
			conf := &ServerConfig***REMOVED***NoClientAuth: true***REMOVED***
			conf.AddHostKey(testSigners["rsa"])
			conn, _, _, err := NewServerConn(c2, conf)
			if err == nil == tt.wantErr ***REMOVED***
				t.Fatalf("got err %v; wantErr %t", err, tt.wantErr)
			***REMOVED***
			if tt.wantErr ***REMOVED***
				// Don't verify the version on an expected error.
				return
			***REMOVED***
			if got := string(conn.ClientVersion()); got != tt.version ***REMOVED***
				t.Fatalf("got %q; want %q", got, tt.version)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestHostKeyCheck(t *testing.T) ***REMOVED***
	for _, tt := range []struct ***REMOVED***
		name      string
		wantError string
		key       PublicKey
	***REMOVED******REMOVED***
		***REMOVED***"no callback", "must specify HostKeyCallback", nil***REMOVED***,
		***REMOVED***"correct key", "", testSigners["rsa"].PublicKey()***REMOVED***,
		***REMOVED***"mismatch", "mismatch", testSigners["ecdsa"].PublicKey()***REMOVED***,
	***REMOVED*** ***REMOVED***
		c1, c2, err := netPipe()
		if err != nil ***REMOVED***
			t.Fatalf("netPipe: %v", err)
		***REMOVED***
		defer c1.Close()
		defer c2.Close()
		serverConf := &ServerConfig***REMOVED***
			NoClientAuth: true,
		***REMOVED***
		serverConf.AddHostKey(testSigners["rsa"])

		go NewServerConn(c1, serverConf)
		clientConf := ClientConfig***REMOVED***
			User: "user",
		***REMOVED***
		if tt.key != nil ***REMOVED***
			clientConf.HostKeyCallback = FixedHostKey(tt.key)
		***REMOVED***

		_, _, _, err = NewClientConn(c2, "", &clientConf)
		if err != nil ***REMOVED***
			if tt.wantError == "" || !strings.Contains(err.Error(), tt.wantError) ***REMOVED***
				t.Errorf("%s: got error %q, missing %q", tt.name, err.Error(), tt.wantError)
			***REMOVED***
		***REMOVED*** else if tt.wantError != "" ***REMOVED***
			t.Errorf("%s: succeeded, but want error string %q", tt.name, tt.wantError)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBannerCallback(t *testing.T) ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	serverConf := &ServerConfig***REMOVED***
		PasswordCallback: func(conn ConnMetadata, password []byte) (*Permissions, error) ***REMOVED***
			return &Permissions***REMOVED******REMOVED***, nil
		***REMOVED***,
		BannerCallback: func(conn ConnMetadata) string ***REMOVED***
			return "Hello World"
		***REMOVED***,
	***REMOVED***
	serverConf.AddHostKey(testSigners["rsa"])
	go NewServerConn(c1, serverConf)

	var receivedBanner string
	var bannerCount int
	clientConf := ClientConfig***REMOVED***
		Auth: []AuthMethod***REMOVED***
			Password("123"),
		***REMOVED***,
		User:            "user",
		HostKeyCallback: InsecureIgnoreHostKey(),
		BannerCallback: func(message string) error ***REMOVED***
			bannerCount++
			receivedBanner = message
			return nil
		***REMOVED***,
	***REMOVED***

	_, _, _, err = NewClientConn(c2, "", &clientConf)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if bannerCount != 1 ***REMOVED***
		t.Errorf("got %d banners; want 1", bannerCount)
	***REMOVED***

	expected := "Hello World"
	if receivedBanner != expected ***REMOVED***
		t.Fatalf("got %s; want %s", receivedBanner, expected)
	***REMOVED***
***REMOVED***
