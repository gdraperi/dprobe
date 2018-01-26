// +build !windows

package hack

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"
)

type bufConn struct ***REMOVED***
	net.Conn
	buf *bytes.Buffer
***REMOVED***

func (bc *bufConn) Read(b []byte) (int, error) ***REMOVED***
	return bc.buf.Read(b)
***REMOVED***

func TestHeaderOverrideHack(t *testing.T) ***REMOVED***
	tests := [][2][]byte***REMOVED***
		***REMOVED***
			[]byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\r\n\r\n"),
			[]byte("GET /foo\nHost: \r\nConnection: close\r\nUser-Agent: Docker\r\n\r\n"),
		***REMOVED***,
		***REMOVED***
			[]byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\nFoo: Bar\r\n"),
			[]byte("GET /foo\nHost: \r\nConnection: close\r\nUser-Agent: Docker\nFoo: Bar\r\n"),
		***REMOVED***,
		***REMOVED***
			[]byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\r\n\r\ntest something!"),
			[]byte("GET /foo\nHost: \r\nConnection: close\r\nUser-Agent: Docker\r\n\r\ntest something!"),
		***REMOVED***,
		***REMOVED***
			[]byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\r\n\r\ntest something! " + strings.Repeat("test", 15000)),
			[]byte("GET /foo\nHost: \r\nConnection: close\r\nUser-Agent: Docker\r\n\r\ntest something! " + strings.Repeat("test", 15000)),
		***REMOVED***,
		***REMOVED***
			[]byte("GET /foo\nFoo: Bar\nHost: /var/run/docker.sock\nUser-Agent: Docker\r\n\r\n"),
			[]byte("GET /foo\nFoo: Bar\nHost: /var/run/docker.sock\nUser-Agent: Docker\r\n\r\n"),
		***REMOVED***,
	***REMOVED***

	// Test for https://github.com/docker/docker/issues/23045
	h0 := "GET /foo\nUser-Agent: Docker\r\n\r\n"
	h0 = h0 + strings.Repeat("a", 4096-len(h0)-1) + "\n"
	tests = append(tests, [2][]byte***REMOVED***[]byte(h0), []byte(h0)***REMOVED***)

	for _, pair := range tests ***REMOVED***
		read := make([]byte, 4096)
		client := &bufConn***REMOVED***
			buf: bytes.NewBuffer(pair[0]),
		***REMOVED***
		l := MalformedHostHeaderOverrideConn***REMOVED***client, true***REMOVED***

		n, err := l.Read(read)
		if err != nil && err != io.EOF ***REMOVED***
			t.Fatalf("read: %d - %d, err: %v\n%s", n, len(pair[0]), err, string(read[:n]))
		***REMOVED***
		if !bytes.Equal(read[:n], pair[1][:n]) ***REMOVED***
			t.Fatalf("\n%s\n%s\n", read[:n], pair[1][:n])
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkWithHack(b *testing.B) ***REMOVED***
	client, srv := net.Pipe()
	done := make(chan struct***REMOVED******REMOVED***)
	req := []byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\n")
	read := make([]byte, 4096)
	b.SetBytes(int64(len(req) * 30))

	l := MalformedHostHeaderOverrideConn***REMOVED***client, true***REMOVED***
	go func() ***REMOVED***
		for ***REMOVED***
			if _, err := srv.Write(req); err != nil ***REMOVED***
				srv.Close()
				break
			***REMOVED***
			l.first = true // make sure each subsequent run uses the hack parsing
		***REMOVED***
		close(done)
	***REMOVED***()

	for i := 0; i < b.N; i++ ***REMOVED***
		for i := 0; i < 30; i++ ***REMOVED***
			if n, err := l.Read(read); err != nil && err != io.EOF ***REMOVED***
				b.Fatalf("read: %d - %d, err: %v\n%s", n, len(req), err, string(read[:n]))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	l.Close()
	<-done
***REMOVED***

func BenchmarkNoHack(b *testing.B) ***REMOVED***
	client, srv := net.Pipe()
	done := make(chan struct***REMOVED******REMOVED***)
	req := []byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\n")
	read := make([]byte, 4096)
	b.SetBytes(int64(len(req) * 30))

	go func() ***REMOVED***
		for ***REMOVED***
			if _, err := srv.Write(req); err != nil ***REMOVED***
				srv.Close()
				break
			***REMOVED***
		***REMOVED***
		close(done)
	***REMOVED***()

	for i := 0; i < b.N; i++ ***REMOVED***
		for i := 0; i < 30; i++ ***REMOVED***
			if _, err := client.Read(read); err != nil && err != io.EOF ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	client.Close()
	<-done
***REMOVED***
