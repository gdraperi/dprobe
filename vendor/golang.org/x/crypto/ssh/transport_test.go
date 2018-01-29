// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"strings"
	"testing"
)

func TestReadVersion(t *testing.T) ***REMOVED***
	longVersion := strings.Repeat("SSH-2.0-bla", 50)[:253]
	multiLineVersion := strings.Repeat("ignored\r\n", 20) + "SSH-2.0-bla\r\n"
	cases := map[string]string***REMOVED***
		"SSH-2.0-bla\r\n":    "SSH-2.0-bla",
		"SSH-2.0-bla\n":      "SSH-2.0-bla",
		multiLineVersion:     "SSH-2.0-bla",
		longVersion + "\r\n": longVersion,
	***REMOVED***

	for in, want := range cases ***REMOVED***
		result, err := readVersion(bytes.NewBufferString(in))
		if err != nil ***REMOVED***
			t.Errorf("readVersion(%q): %s", in, err)
		***REMOVED***
		got := string(result)
		if got != want ***REMOVED***
			t.Errorf("got %q, want %q", got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReadVersionError(t *testing.T) ***REMOVED***
	longVersion := strings.Repeat("SSH-2.0-bla", 50)[:253]
	multiLineVersion := strings.Repeat("ignored\r\n", 50) + "SSH-2.0-bla\r\n"
	cases := []string***REMOVED***
		longVersion + "too-long\r\n",
		multiLineVersion,
	***REMOVED***
	for _, in := range cases ***REMOVED***
		if _, err := readVersion(bytes.NewBufferString(in)); err == nil ***REMOVED***
			t.Errorf("readVersion(%q) should have failed", in)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestExchangeVersionsBasic(t *testing.T) ***REMOVED***
	v := "SSH-2.0-bla"
	buf := bytes.NewBufferString(v + "\r\n")
	them, err := exchangeVersions(buf, []byte("xyz"))
	if err != nil ***REMOVED***
		t.Errorf("exchangeVersions: %v", err)
	***REMOVED***

	if want := "SSH-2.0-bla"; string(them) != want ***REMOVED***
		t.Errorf("got %q want %q for our version", them, want)
	***REMOVED***
***REMOVED***

func TestExchangeVersions(t *testing.T) ***REMOVED***
	cases := []string***REMOVED***
		"not\x000allowed",
		"not allowed\x01\r\n",
	***REMOVED***
	for _, c := range cases ***REMOVED***
		buf := bytes.NewBufferString("SSH-2.0-bla\r\n")
		if _, err := exchangeVersions(buf, []byte(c)); err == nil ***REMOVED***
			t.Errorf("exchangeVersions(%q): should have failed", c)
		***REMOVED***
	***REMOVED***
***REMOVED***

type closerBuffer struct ***REMOVED***
	bytes.Buffer
***REMOVED***

func (b *closerBuffer) Close() error ***REMOVED***
	return nil
***REMOVED***

func TestTransportMaxPacketWrite(t *testing.T) ***REMOVED***
	buf := &closerBuffer***REMOVED******REMOVED***
	tr := newTransport(buf, rand.Reader, true)
	huge := make([]byte, maxPacket+1)
	err := tr.writePacket(huge)
	if err == nil ***REMOVED***
		t.Errorf("transport accepted write for a huge packet.")
	***REMOVED***
***REMOVED***

func TestTransportMaxPacketReader(t *testing.T) ***REMOVED***
	var header [5]byte
	huge := make([]byte, maxPacket+128)
	binary.BigEndian.PutUint32(header[0:], uint32(len(huge)))
	// padding.
	header[4] = 0

	buf := &closerBuffer***REMOVED******REMOVED***
	buf.Write(header[:])
	buf.Write(huge)

	tr := newTransport(buf, rand.Reader, true)
	_, err := tr.readPacket()
	if err == nil ***REMOVED***
		t.Errorf("transport succeeded reading huge packet.")
	***REMOVED*** else if !strings.Contains(err.Error(), "large") ***REMOVED***
		t.Errorf("got %q, should mention %q", err.Error(), "large")
	***REMOVED***
***REMOVED***
