// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/http2/hpack"
)

var stderrVerbose = flag.Bool("stderr_verbose", false, "Mirror verbosity to stderr, unbuffered")

func stderrv() io.Writer ***REMOVED***
	if *stderrVerbose ***REMOVED***
		return os.Stderr
	***REMOVED***

	return ioutil.Discard
***REMOVED***

type serverTester struct ***REMOVED***
	cc             net.Conn // client conn
	t              testing.TB
	ts             *httptest.Server
	fr             *Framer
	serverLogBuf   bytes.Buffer // logger for httptest.Server
	logFilter      []string     // substrings to filter out
	scMu           sync.Mutex   // guards sc
	sc             *serverConn
	hpackDec       *hpack.Decoder
	decodedHeaders [][2]string

	// If http2debug!=2, then we capture Frame debug logs that will be written
	// to t.Log after a test fails. The read and write logs use separate locks
	// and buffers so we don't accidentally introduce synchronization between
	// the read and write goroutines, which may hide data races.
	frameReadLogMu   sync.Mutex
	frameReadLogBuf  bytes.Buffer
	frameWriteLogMu  sync.Mutex
	frameWriteLogBuf bytes.Buffer

	// writing headers:
	headerBuf bytes.Buffer
	hpackEnc  *hpack.Encoder
***REMOVED***

func init() ***REMOVED***
	testHookOnPanicMu = new(sync.Mutex)
	goAwayTimeout = 25 * time.Millisecond
***REMOVED***

func resetHooks() ***REMOVED***
	testHookOnPanicMu.Lock()
	testHookOnPanic = nil
	testHookOnPanicMu.Unlock()
***REMOVED***

type serverTesterOpt string

var optOnlyServer = serverTesterOpt("only_server")
var optQuiet = serverTesterOpt("quiet_logging")
var optFramerReuseFrames = serverTesterOpt("frame_reuse_frames")

func newServerTester(t testing.TB, handler http.HandlerFunc, opts ...interface***REMOVED******REMOVED***) *serverTester ***REMOVED***
	resetHooks()

	ts := httptest.NewUnstartedServer(handler)

	tlsConfig := &tls.Config***REMOVED***
		InsecureSkipVerify: true,
		NextProtos:         []string***REMOVED***NextProtoTLS***REMOVED***,
	***REMOVED***

	var onlyServer, quiet, framerReuseFrames bool
	h2server := new(Server)
	for _, opt := range opts ***REMOVED***
		switch v := opt.(type) ***REMOVED***
		case func(*tls.Config):
			v(tlsConfig)
		case func(*httptest.Server):
			v(ts)
		case func(*Server):
			v(h2server)
		case serverTesterOpt:
			switch v ***REMOVED***
			case optOnlyServer:
				onlyServer = true
			case optQuiet:
				quiet = true
			case optFramerReuseFrames:
				framerReuseFrames = true
			***REMOVED***
		case func(net.Conn, http.ConnState):
			ts.Config.ConnState = v
		default:
			t.Fatalf("unknown newServerTester option type %T", v)
		***REMOVED***
	***REMOVED***

	ConfigureServer(ts.Config, h2server)

	st := &serverTester***REMOVED***
		t:  t,
		ts: ts,
	***REMOVED***
	st.hpackEnc = hpack.NewEncoder(&st.headerBuf)
	st.hpackDec = hpack.NewDecoder(initialHeaderTableSize, st.onHeaderField)

	ts.TLS = ts.Config.TLSConfig // the httptest.Server has its own copy of this TLS config
	if quiet ***REMOVED***
		ts.Config.ErrorLog = log.New(ioutil.Discard, "", 0)
	***REMOVED*** else ***REMOVED***
		ts.Config.ErrorLog = log.New(io.MultiWriter(stderrv(), twriter***REMOVED***t: t, st: st***REMOVED***, &st.serverLogBuf), "", log.LstdFlags)
	***REMOVED***
	ts.StartTLS()

	if VerboseLogs ***REMOVED***
		t.Logf("Running test server at: %s", ts.URL)
	***REMOVED***
	testHookGetServerConn = func(v *serverConn) ***REMOVED***
		st.scMu.Lock()
		defer st.scMu.Unlock()
		st.sc = v
	***REMOVED***
	log.SetOutput(io.MultiWriter(stderrv(), twriter***REMOVED***t: t, st: st***REMOVED***))
	if !onlyServer ***REMOVED***
		cc, err := tls.Dial("tcp", ts.Listener.Addr().String(), tlsConfig)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		st.cc = cc
		st.fr = NewFramer(cc, cc)
		if framerReuseFrames ***REMOVED***
			st.fr.SetReuseFrames()
		***REMOVED***
		if !logFrameReads && !logFrameWrites ***REMOVED***
			st.fr.debugReadLoggerf = func(m string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
				m = time.Now().Format("2006-01-02 15:04:05.999999999 ") + strings.TrimPrefix(m, "http2: ") + "\n"
				st.frameReadLogMu.Lock()
				fmt.Fprintf(&st.frameReadLogBuf, m, v...)
				st.frameReadLogMu.Unlock()
			***REMOVED***
			st.fr.debugWriteLoggerf = func(m string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
				m = time.Now().Format("2006-01-02 15:04:05.999999999 ") + strings.TrimPrefix(m, "http2: ") + "\n"
				st.frameWriteLogMu.Lock()
				fmt.Fprintf(&st.frameWriteLogBuf, m, v...)
				st.frameWriteLogMu.Unlock()
			***REMOVED***
			st.fr.logReads = true
			st.fr.logWrites = true
		***REMOVED***
	***REMOVED***
	return st
***REMOVED***

func (st *serverTester) closeConn() ***REMOVED***
	st.scMu.Lock()
	defer st.scMu.Unlock()
	st.sc.conn.Close()
***REMOVED***

func (st *serverTester) addLogFilter(phrase string) ***REMOVED***
	st.logFilter = append(st.logFilter, phrase)
***REMOVED***

func (st *serverTester) stream(id uint32) *stream ***REMOVED***
	ch := make(chan *stream, 1)
	st.sc.serveMsgCh <- func(int) ***REMOVED***
		ch <- st.sc.streams[id]
	***REMOVED***
	return <-ch
***REMOVED***

func (st *serverTester) streamState(id uint32) streamState ***REMOVED***
	ch := make(chan streamState, 1)
	st.sc.serveMsgCh <- func(int) ***REMOVED***
		state, _ := st.sc.state(id)
		ch <- state
	***REMOVED***
	return <-ch
***REMOVED***

// loopNum reports how many times this conn's select loop has gone around.
func (st *serverTester) loopNum() int ***REMOVED***
	lastc := make(chan int, 1)
	st.sc.serveMsgCh <- func(loopNum int) ***REMOVED***
		lastc <- loopNum
	***REMOVED***
	return <-lastc
***REMOVED***

// awaitIdle heuristically awaits for the server conn's select loop to be idle.
// The heuristic is that the server connection's serve loop must schedule
// 50 times in a row without any channel sends or receives occurring.
func (st *serverTester) awaitIdle() ***REMOVED***
	remain := 50
	last := st.loopNum()
	for remain > 0 ***REMOVED***
		n := st.loopNum()
		if n == last+1 ***REMOVED***
			remain--
		***REMOVED*** else ***REMOVED***
			remain = 50
		***REMOVED***
		last = n
	***REMOVED***
***REMOVED***

func (st *serverTester) Close() ***REMOVED***
	if st.t.Failed() ***REMOVED***
		st.frameReadLogMu.Lock()
		if st.frameReadLogBuf.Len() > 0 ***REMOVED***
			st.t.Logf("Framer read log:\n%s", st.frameReadLogBuf.String())
		***REMOVED***
		st.frameReadLogMu.Unlock()

		st.frameWriteLogMu.Lock()
		if st.frameWriteLogBuf.Len() > 0 ***REMOVED***
			st.t.Logf("Framer write log:\n%s", st.frameWriteLogBuf.String())
		***REMOVED***
		st.frameWriteLogMu.Unlock()

		// If we failed already (and are likely in a Fatal,
		// unwindowing), force close the connection, so the
		// httptest.Server doesn't wait forever for the conn
		// to close.
		if st.cc != nil ***REMOVED***
			st.cc.Close()
		***REMOVED***
	***REMOVED***
	st.ts.Close()
	if st.cc != nil ***REMOVED***
		st.cc.Close()
	***REMOVED***
	log.SetOutput(os.Stderr)
***REMOVED***

// greet initiates the client's HTTP/2 connection into a state where
// frames may be sent.
func (st *serverTester) greet() ***REMOVED***
	st.greetAndCheckSettings(func(Setting) error ***REMOVED*** return nil ***REMOVED***)
***REMOVED***

func (st *serverTester) greetAndCheckSettings(checkSetting func(s Setting) error) ***REMOVED***
	st.writePreface()
	st.writeInitialSettings()
	st.wantSettings().ForeachSetting(checkSetting)
	st.writeSettingsAck()

	// The initial WINDOW_UPDATE and SETTINGS ACK can come in any order.
	var gotSettingsAck bool
	var gotWindowUpdate bool

	for i := 0; i < 2; i++ ***REMOVED***
		f, err := st.readFrame()
		if err != nil ***REMOVED***
			st.t.Fatal(err)
		***REMOVED***
		switch f := f.(type) ***REMOVED***
		case *SettingsFrame:
			if !f.Header().Flags.Has(FlagSettingsAck) ***REMOVED***
				st.t.Fatal("Settings Frame didn't have ACK set")
			***REMOVED***
			gotSettingsAck = true

		case *WindowUpdateFrame:
			if f.FrameHeader.StreamID != 0 ***REMOVED***
				st.t.Fatalf("WindowUpdate StreamID = %d; want 0", f.FrameHeader.StreamID)
			***REMOVED***
			incr := uint32((&Server***REMOVED******REMOVED***).initialConnRecvWindowSize() - initialWindowSize)
			if f.Increment != incr ***REMOVED***
				st.t.Fatalf("WindowUpdate increment = %d; want %d", f.Increment, incr)
			***REMOVED***
			gotWindowUpdate = true

		default:
			st.t.Fatalf("Wanting a settings ACK or window update, received a %T", f)
		***REMOVED***
	***REMOVED***

	if !gotSettingsAck ***REMOVED***
		st.t.Fatalf("Didn't get a settings ACK")
	***REMOVED***
	if !gotWindowUpdate ***REMOVED***
		st.t.Fatalf("Didn't get a window update")
	***REMOVED***
***REMOVED***

func (st *serverTester) writePreface() ***REMOVED***
	n, err := st.cc.Write(clientPreface)
	if err != nil ***REMOVED***
		st.t.Fatalf("Error writing client preface: %v", err)
	***REMOVED***
	if n != len(clientPreface) ***REMOVED***
		st.t.Fatalf("Writing client preface, wrote %d bytes; want %d", n, len(clientPreface))
	***REMOVED***
***REMOVED***

func (st *serverTester) writeInitialSettings() ***REMOVED***
	if err := st.fr.WriteSettings(); err != nil ***REMOVED***
		st.t.Fatalf("Error writing initial SETTINGS frame from client to server: %v", err)
	***REMOVED***
***REMOVED***

func (st *serverTester) writeSettingsAck() ***REMOVED***
	if err := st.fr.WriteSettingsAck(); err != nil ***REMOVED***
		st.t.Fatalf("Error writing ACK of server's SETTINGS: %v", err)
	***REMOVED***
***REMOVED***

func (st *serverTester) writeHeaders(p HeadersFrameParam) ***REMOVED***
	if err := st.fr.WriteHeaders(p); err != nil ***REMOVED***
		st.t.Fatalf("Error writing HEADERS: %v", err)
	***REMOVED***
***REMOVED***

func (st *serverTester) writePriority(id uint32, p PriorityParam) ***REMOVED***
	if err := st.fr.WritePriority(id, p); err != nil ***REMOVED***
		st.t.Fatalf("Error writing PRIORITY: %v", err)
	***REMOVED***
***REMOVED***

func (st *serverTester) encodeHeaderField(k, v string) ***REMOVED***
	err := st.hpackEnc.WriteField(hpack.HeaderField***REMOVED***Name: k, Value: v***REMOVED***)
	if err != nil ***REMOVED***
		st.t.Fatalf("HPACK encoding error for %q/%q: %v", k, v, err)
	***REMOVED***
***REMOVED***

// encodeHeaderRaw is the magic-free version of encodeHeader.
// It takes 0 or more (k, v) pairs and encodes them.
func (st *serverTester) encodeHeaderRaw(headers ...string) []byte ***REMOVED***
	if len(headers)%2 == 1 ***REMOVED***
		panic("odd number of kv args")
	***REMOVED***
	st.headerBuf.Reset()
	for len(headers) > 0 ***REMOVED***
		k, v := headers[0], headers[1]
		st.encodeHeaderField(k, v)
		headers = headers[2:]
	***REMOVED***
	return st.headerBuf.Bytes()
***REMOVED***

// encodeHeader encodes headers and returns their HPACK bytes. headers
// must contain an even number of key/value pairs. There may be
// multiple pairs for keys (e.g. "cookie").  The :method, :path, and
// :scheme headers default to GET, / and https. The :authority header
// defaults to st.ts.Listener.Addr().
func (st *serverTester) encodeHeader(headers ...string) []byte ***REMOVED***
	if len(headers)%2 == 1 ***REMOVED***
		panic("odd number of kv args")
	***REMOVED***

	st.headerBuf.Reset()
	defaultAuthority := st.ts.Listener.Addr().String()

	if len(headers) == 0 ***REMOVED***
		// Fast path, mostly for benchmarks, so test code doesn't pollute
		// profiles when we're looking to improve server allocations.
		st.encodeHeaderField(":method", "GET")
		st.encodeHeaderField(":scheme", "https")
		st.encodeHeaderField(":authority", defaultAuthority)
		st.encodeHeaderField(":path", "/")
		return st.headerBuf.Bytes()
	***REMOVED***

	if len(headers) == 2 && headers[0] == ":method" ***REMOVED***
		// Another fast path for benchmarks.
		st.encodeHeaderField(":method", headers[1])
		st.encodeHeaderField(":scheme", "https")
		st.encodeHeaderField(":authority", defaultAuthority)
		st.encodeHeaderField(":path", "/")
		return st.headerBuf.Bytes()
	***REMOVED***

	pseudoCount := map[string]int***REMOVED******REMOVED***
	keys := []string***REMOVED***":method", ":scheme", ":authority", ":path"***REMOVED***
	vals := map[string][]string***REMOVED***
		":method":    ***REMOVED***"GET"***REMOVED***,
		":scheme":    ***REMOVED***"https"***REMOVED***,
		":authority": ***REMOVED***defaultAuthority***REMOVED***,
		":path":      ***REMOVED***"/"***REMOVED***,
	***REMOVED***
	for len(headers) > 0 ***REMOVED***
		k, v := headers[0], headers[1]
		headers = headers[2:]
		if _, ok := vals[k]; !ok ***REMOVED***
			keys = append(keys, k)
		***REMOVED***
		if strings.HasPrefix(k, ":") ***REMOVED***
			pseudoCount[k]++
			if pseudoCount[k] == 1 ***REMOVED***
				vals[k] = []string***REMOVED***v***REMOVED***
			***REMOVED*** else ***REMOVED***
				// Allows testing of invalid headers w/ dup pseudo fields.
				vals[k] = append(vals[k], v)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			vals[k] = append(vals[k], v)
		***REMOVED***
	***REMOVED***
	for _, k := range keys ***REMOVED***
		for _, v := range vals[k] ***REMOVED***
			st.encodeHeaderField(k, v)
		***REMOVED***
	***REMOVED***
	return st.headerBuf.Bytes()
***REMOVED***

// bodylessReq1 writes a HEADERS frames with StreamID 1 and EndStream and EndHeaders set.
func (st *serverTester) bodylessReq1(headers ...string) ***REMOVED***
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1, // clients send odd numbers
		BlockFragment: st.encodeHeader(headers...),
		EndStream:     true,
		EndHeaders:    true,
	***REMOVED***)
***REMOVED***

func (st *serverTester) writeData(streamID uint32, endStream bool, data []byte) ***REMOVED***
	if err := st.fr.WriteData(streamID, endStream, data); err != nil ***REMOVED***
		st.t.Fatalf("Error writing DATA: %v", err)
	***REMOVED***
***REMOVED***

func (st *serverTester) writeDataPadded(streamID uint32, endStream bool, data, pad []byte) ***REMOVED***
	if err := st.fr.WriteDataPadded(streamID, endStream, data, pad); err != nil ***REMOVED***
		st.t.Fatalf("Error writing DATA: %v", err)
	***REMOVED***
***REMOVED***

func readFrameTimeout(fr *Framer, wait time.Duration) (Frame, error) ***REMOVED***
	ch := make(chan interface***REMOVED******REMOVED***, 1)
	go func() ***REMOVED***
		fr, err := fr.ReadFrame()
		if err != nil ***REMOVED***
			ch <- err
		***REMOVED*** else ***REMOVED***
			ch <- fr
		***REMOVED***
	***REMOVED***()
	t := time.NewTimer(wait)
	select ***REMOVED***
	case v := <-ch:
		t.Stop()
		if fr, ok := v.(Frame); ok ***REMOVED***
			return fr, nil
		***REMOVED***
		return nil, v.(error)
	case <-t.C:
		return nil, errors.New("timeout waiting for frame")
	***REMOVED***
***REMOVED***

func (st *serverTester) readFrame() (Frame, error) ***REMOVED***
	return readFrameTimeout(st.fr, 2*time.Second)
***REMOVED***

func (st *serverTester) wantHeaders() *HeadersFrame ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting a HEADERS frame: %v", err)
	***REMOVED***
	hf, ok := f.(*HeadersFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *HeadersFrame", f)
	***REMOVED***
	return hf
***REMOVED***

func (st *serverTester) wantContinuation() *ContinuationFrame ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting a CONTINUATION frame: %v", err)
	***REMOVED***
	cf, ok := f.(*ContinuationFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *ContinuationFrame", f)
	***REMOVED***
	return cf
***REMOVED***

func (st *serverTester) wantData() *DataFrame ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting a DATA frame: %v", err)
	***REMOVED***
	df, ok := f.(*DataFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *DataFrame", f)
	***REMOVED***
	return df
***REMOVED***

func (st *serverTester) wantSettings() *SettingsFrame ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting a SETTINGS frame: %v", err)
	***REMOVED***
	sf, ok := f.(*SettingsFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *SettingsFrame", f)
	***REMOVED***
	return sf
***REMOVED***

func (st *serverTester) wantPing() *PingFrame ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting a PING frame: %v", err)
	***REMOVED***
	pf, ok := f.(*PingFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *PingFrame", f)
	***REMOVED***
	return pf
***REMOVED***

func (st *serverTester) wantGoAway() *GoAwayFrame ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting a GOAWAY frame: %v", err)
	***REMOVED***
	gf, ok := f.(*GoAwayFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *GoAwayFrame", f)
	***REMOVED***
	return gf
***REMOVED***

func (st *serverTester) wantRSTStream(streamID uint32, errCode ErrCode) ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting an RSTStream frame: %v", err)
	***REMOVED***
	rs, ok := f.(*RSTStreamFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *RSTStreamFrame", f)
	***REMOVED***
	if rs.FrameHeader.StreamID != streamID ***REMOVED***
		st.t.Fatalf("RSTStream StreamID = %d; want %d", rs.FrameHeader.StreamID, streamID)
	***REMOVED***
	if rs.ErrCode != errCode ***REMOVED***
		st.t.Fatalf("RSTStream ErrCode = %d (%s); want %d (%s)", rs.ErrCode, rs.ErrCode, errCode, errCode)
	***REMOVED***
***REMOVED***

func (st *serverTester) wantWindowUpdate(streamID, incr uint32) ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatalf("Error while expecting a WINDOW_UPDATE frame: %v", err)
	***REMOVED***
	wu, ok := f.(*WindowUpdateFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("got a %T; want *WindowUpdateFrame", f)
	***REMOVED***
	if wu.FrameHeader.StreamID != streamID ***REMOVED***
		st.t.Fatalf("WindowUpdate StreamID = %d; want %d", wu.FrameHeader.StreamID, streamID)
	***REMOVED***
	if wu.Increment != incr ***REMOVED***
		st.t.Fatalf("WindowUpdate increment = %d; want %d", wu.Increment, incr)
	***REMOVED***
***REMOVED***

func (st *serverTester) wantSettingsAck() ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatal(err)
	***REMOVED***
	sf, ok := f.(*SettingsFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("Wanting a settings ACK, received a %T", f)
	***REMOVED***
	if !sf.Header().Flags.Has(FlagSettingsAck) ***REMOVED***
		st.t.Fatal("Settings Frame didn't have ACK set")
	***REMOVED***
***REMOVED***

func (st *serverTester) wantPushPromise() *PushPromiseFrame ***REMOVED***
	f, err := st.readFrame()
	if err != nil ***REMOVED***
		st.t.Fatal(err)
	***REMOVED***
	ppf, ok := f.(*PushPromiseFrame)
	if !ok ***REMOVED***
		st.t.Fatalf("Wanted PushPromise, received %T", ppf)
	***REMOVED***
	return ppf
***REMOVED***

func TestServer(t *testing.T) ***REMOVED***
	gotReq := make(chan bool, 1)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Foo", "Bar")
		gotReq <- true
	***REMOVED***)
	defer st.Close()

	covers("3.5", `
		The server connection preface consists of a potentially empty
		SETTINGS frame ([SETTINGS]) that MUST be the first frame the
		server sends in the HTTP/2 connection.
	`)

	st.greet()
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1, // clients send odd numbers
		BlockFragment: st.encodeHeader(),
		EndStream:     true, // no DATA frames
		EndHeaders:    true,
	***REMOVED***)

	select ***REMOVED***
	case <-gotReq:
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for request")
	***REMOVED***
***REMOVED***

func TestServer_Request_Get(t *testing.T) ***REMOVED***
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader("foo-bar", "some-value"),
			EndStream:     true, // no DATA frames
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***, func(r *http.Request) ***REMOVED***
		if r.Method != "GET" ***REMOVED***
			t.Errorf("Method = %q; want GET", r.Method)
		***REMOVED***
		if r.URL.Path != "/" ***REMOVED***
			t.Errorf("URL.Path = %q; want /", r.URL.Path)
		***REMOVED***
		if r.ContentLength != 0 ***REMOVED***
			t.Errorf("ContentLength = %v; want 0", r.ContentLength)
		***REMOVED***
		if r.Close ***REMOVED***
			t.Error("Close = true; want false")
		***REMOVED***
		if !strings.Contains(r.RemoteAddr, ":") ***REMOVED***
			t.Errorf("RemoteAddr = %q; want something with a colon", r.RemoteAddr)
		***REMOVED***
		if r.Proto != "HTTP/2.0" || r.ProtoMajor != 2 || r.ProtoMinor != 0 ***REMOVED***
			t.Errorf("Proto = %q Major=%v,Minor=%v; want HTTP/2.0", r.Proto, r.ProtoMajor, r.ProtoMinor)
		***REMOVED***
		wantHeader := http.Header***REMOVED***
			"Foo-Bar": []string***REMOVED***"some-value"***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(r.Header, wantHeader) ***REMOVED***
			t.Errorf("Header = %#v; want %#v", r.Header, wantHeader)
		***REMOVED***
		if n, err := r.Body.Read([]byte(" ")); err != io.EOF || n != 0 ***REMOVED***
			t.Errorf("Read = %d, %v; want 0, EOF", n, err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Request_Get_PathSlashes(t *testing.T) ***REMOVED***
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":path", "/%2f/"),
			EndStream:     true, // no DATA frames
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***, func(r *http.Request) ***REMOVED***
		if r.RequestURI != "/%2f/" ***REMOVED***
			t.Errorf("RequestURI = %q; want /%%2f/", r.RequestURI)
		***REMOVED***
		if r.URL.Path != "///" ***REMOVED***
			t.Errorf("URL.Path = %q; want ///", r.URL.Path)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// TODO: add a test with EndStream=true on the HEADERS but setting a
// Content-Length anyway. Should we just omit it and force it to
// zero?

func TestServer_Request_Post_NoContentLength_EndStream(t *testing.T) ***REMOVED***
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":method", "POST"),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***, func(r *http.Request) ***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("Method = %q; want POST", r.Method)
		***REMOVED***
		if r.ContentLength != 0 ***REMOVED***
			t.Errorf("ContentLength = %v; want 0", r.ContentLength)
		***REMOVED***
		if n, err := r.Body.Read([]byte(" ")); err != io.EOF || n != 0 ***REMOVED***
			t.Errorf("Read = %d, %v; want 0, EOF", n, err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Request_Post_Body_ImmediateEOF(t *testing.T) ***REMOVED***
	testBodyContents(t, -1, "", func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":method", "POST"),
			EndStream:     false, // to say DATA frames are coming
			EndHeaders:    true,
		***REMOVED***)
		st.writeData(1, true, nil) // just kidding. empty body.
	***REMOVED***)
***REMOVED***

func TestServer_Request_Post_Body_OneData(t *testing.T) ***REMOVED***
	const content = "Some content"
	testBodyContents(t, -1, content, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":method", "POST"),
			EndStream:     false, // to say DATA frames are coming
			EndHeaders:    true,
		***REMOVED***)
		st.writeData(1, true, []byte(content))
	***REMOVED***)
***REMOVED***

func TestServer_Request_Post_Body_TwoData(t *testing.T) ***REMOVED***
	const content = "Some content"
	testBodyContents(t, -1, content, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":method", "POST"),
			EndStream:     false, // to say DATA frames are coming
			EndHeaders:    true,
		***REMOVED***)
		st.writeData(1, false, []byte(content[:5]))
		st.writeData(1, true, []byte(content[5:]))
	***REMOVED***)
***REMOVED***

func TestServer_Request_Post_Body_ContentLength_Correct(t *testing.T) ***REMOVED***
	const content = "Some content"
	testBodyContents(t, int64(len(content)), content, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID: 1, // clients send odd numbers
			BlockFragment: st.encodeHeader(
				":method", "POST",
				"content-length", strconv.Itoa(len(content)),
			),
			EndStream:  false, // to say DATA frames are coming
			EndHeaders: true,
		***REMOVED***)
		st.writeData(1, true, []byte(content))
	***REMOVED***)
***REMOVED***

func TestServer_Request_Post_Body_ContentLength_TooLarge(t *testing.T) ***REMOVED***
	testBodyContentsFail(t, 3, "request declared a Content-Length of 3 but only wrote 2 bytes",
		func(st *serverTester) ***REMOVED***
			st.writeHeaders(HeadersFrameParam***REMOVED***
				StreamID: 1, // clients send odd numbers
				BlockFragment: st.encodeHeader(
					":method", "POST",
					"content-length", "3",
				),
				EndStream:  false, // to say DATA frames are coming
				EndHeaders: true,
			***REMOVED***)
			st.writeData(1, true, []byte("12"))
		***REMOVED***)
***REMOVED***

func TestServer_Request_Post_Body_ContentLength_TooSmall(t *testing.T) ***REMOVED***
	testBodyContentsFail(t, 4, "sender tried to send more than declared Content-Length of 4 bytes",
		func(st *serverTester) ***REMOVED***
			st.writeHeaders(HeadersFrameParam***REMOVED***
				StreamID: 1, // clients send odd numbers
				BlockFragment: st.encodeHeader(
					":method", "POST",
					"content-length", "4",
				),
				EndStream:  false, // to say DATA frames are coming
				EndHeaders: true,
			***REMOVED***)
			st.writeData(1, true, []byte("12345"))
		***REMOVED***)
***REMOVED***

func testBodyContents(t *testing.T, wantContentLength int64, wantBody string, write func(st *serverTester)) ***REMOVED***
	testServerRequest(t, write, func(r *http.Request) ***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("Method = %q; want POST", r.Method)
		***REMOVED***
		if r.ContentLength != wantContentLength ***REMOVED***
			t.Errorf("ContentLength = %v; want %d", r.ContentLength, wantContentLength)
		***REMOVED***
		all, err := ioutil.ReadAll(r.Body)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if string(all) != wantBody ***REMOVED***
			t.Errorf("Read = %q; want %q", all, wantBody)
		***REMOVED***
		if err := r.Body.Close(); err != nil ***REMOVED***
			t.Fatalf("Close: %v", err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func testBodyContentsFail(t *testing.T, wantContentLength int64, wantReadError string, write func(st *serverTester)) ***REMOVED***
	testServerRequest(t, write, func(r *http.Request) ***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Errorf("Method = %q; want POST", r.Method)
		***REMOVED***
		if r.ContentLength != wantContentLength ***REMOVED***
			t.Errorf("ContentLength = %v; want %d", r.ContentLength, wantContentLength)
		***REMOVED***
		all, err := ioutil.ReadAll(r.Body)
		if err == nil ***REMOVED***
			t.Fatalf("expected an error (%q) reading from the body. Successfully read %q instead.",
				wantReadError, all)
		***REMOVED***
		if !strings.Contains(err.Error(), wantReadError) ***REMOVED***
			t.Fatalf("Body.Read = %v; want substring %q", err, wantReadError)
		***REMOVED***
		if err := r.Body.Close(); err != nil ***REMOVED***
			t.Fatalf("Close: %v", err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Using a Host header, instead of :authority
func TestServer_Request_Get_Host(t *testing.T) ***REMOVED***
	const host = "example.com"
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":authority", "", "host", host),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***, func(r *http.Request) ***REMOVED***
		if r.Host != host ***REMOVED***
			t.Errorf("Host = %q; want %q", r.Host, host)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Using an :authority pseudo-header, instead of Host
func TestServer_Request_Get_Authority(t *testing.T) ***REMOVED***
	const host = "example.com"
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":authority", host),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***, func(r *http.Request) ***REMOVED***
		if r.Host != host ***REMOVED***
			t.Errorf("Host = %q; want %q", r.Host, host)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Request_WithContinuation(t *testing.T) ***REMOVED***
	wantHeader := http.Header***REMOVED***
		"Foo-One":   []string***REMOVED***"value-one"***REMOVED***,
		"Foo-Two":   []string***REMOVED***"value-two"***REMOVED***,
		"Foo-Three": []string***REMOVED***"value-three"***REMOVED***,
	***REMOVED***
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		fullHeaders := st.encodeHeader(
			"foo-one", "value-one",
			"foo-two", "value-two",
			"foo-three", "value-three",
		)
		remain := fullHeaders
		chunks := 0
		for len(remain) > 0 ***REMOVED***
			const maxChunkSize = 5
			chunk := remain
			if len(chunk) > maxChunkSize ***REMOVED***
				chunk = chunk[:maxChunkSize]
			***REMOVED***
			remain = remain[len(chunk):]

			if chunks == 0 ***REMOVED***
				st.writeHeaders(HeadersFrameParam***REMOVED***
					StreamID:      1, // clients send odd numbers
					BlockFragment: chunk,
					EndStream:     true,  // no DATA frames
					EndHeaders:    false, // we'll have continuation frames
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				err := st.fr.WriteContinuation(1, len(remain) == 0, chunk)
				if err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
			***REMOVED***
			chunks++
		***REMOVED***
		if chunks < 2 ***REMOVED***
			t.Fatal("too few chunks")
		***REMOVED***
	***REMOVED***, func(r *http.Request) ***REMOVED***
		if !reflect.DeepEqual(r.Header, wantHeader) ***REMOVED***
			t.Errorf("Header = %#v; want %#v", r.Header, wantHeader)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Concatenated cookie headers. ("8.1.2.5 Compressing the Cookie Header Field")
func TestServer_Request_CookieConcat(t *testing.T) ***REMOVED***
	const host = "example.com"
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		st.bodylessReq1(
			":authority", host,
			"cookie", "a=b",
			"cookie", "c=d",
			"cookie", "e=f",
		)
	***REMOVED***, func(r *http.Request) ***REMOVED***
		const want = "a=b; c=d; e=f"
		if got := r.Header.Get("Cookie"); got != want ***REMOVED***
			t.Errorf("Cookie = %q; want %q", got, want)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_CapitalHeader(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1("UPPER", "v") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_HeaderFieldNameColon(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1("has:colon", "v") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_HeaderFieldNameNULL(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1("has\x00null", "v") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_HeaderFieldNameEmpty(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1("", "v") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_HeaderFieldValueNewline(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1("foo", "has\nnewline") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_HeaderFieldValueCR(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1("foo", "has\rcarriage") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_HeaderFieldValueDEL(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1("foo", "has\x7fdel") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_Pseudo_Missing_method(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1(":method", "") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_Pseudo_ExactlyOne(t *testing.T) ***REMOVED***
	// 8.1.2.3 Request Pseudo-Header Fields
	// "All HTTP/2 requests MUST include exactly one valid value" ...
	testRejectRequest(t, func(st *serverTester) ***REMOVED***
		st.addLogFilter("duplicate pseudo-header")
		st.bodylessReq1(":method", "GET", ":method", "POST")
	***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_Pseudo_AfterRegular(t *testing.T) ***REMOVED***
	// 8.1.2.3 Request Pseudo-Header Fields
	// "All pseudo-header fields MUST appear in the header block
	// before regular header fields. Any request or response that
	// contains a pseudo-header field that appears in a header
	// block after a regular header field MUST be treated as
	// malformed (Section 8.1.2.6)."
	testRejectRequest(t, func(st *serverTester) ***REMOVED***
		st.addLogFilter("pseudo-header after regular header")
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":method", Value: "GET"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: "regular", Value: "foobar"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":path", Value: "/"***REMOVED***)
		enc.WriteField(hpack.HeaderField***REMOVED***Name: ":scheme", Value: "https"***REMOVED***)
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: buf.Bytes(),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_Pseudo_Missing_path(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1(":path", "") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_Pseudo_Missing_scheme(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1(":scheme", "") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_Pseudo_scheme_invalid(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED*** st.bodylessReq1(":scheme", "bogus") ***REMOVED***)
***REMOVED***

func TestServer_Request_Reject_Pseudo_Unknown(t *testing.T) ***REMOVED***
	testRejectRequest(t, func(st *serverTester) ***REMOVED***
		st.addLogFilter(`invalid pseudo-header ":unknown_thing"`)
		st.bodylessReq1(":unknown_thing", "")
	***REMOVED***)
***REMOVED***

func testRejectRequest(t *testing.T, send func(*serverTester)) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		t.Error("server request made it to handler; should've been rejected")
	***REMOVED***)
	defer st.Close()

	st.greet()
	send(st)
	st.wantRSTStream(1, ErrCodeProtocol)
***REMOVED***

func testRejectRequestWithProtocolError(t *testing.T, send func(*serverTester)) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		t.Error("server request made it to handler; should've been rejected")
	***REMOVED***, optQuiet)
	defer st.Close()

	st.greet()
	send(st)
	gf := st.wantGoAway()
	if gf.ErrCode != ErrCodeProtocol ***REMOVED***
		t.Errorf("err code = %v; want %v", gf.ErrCode, ErrCodeProtocol)
	***REMOVED***
***REMOVED***

// Section 5.1, on idle connections: "Receiving any frame other than
// HEADERS or PRIORITY on a stream in this state MUST be treated as a
// connection error (Section 5.4.1) of type PROTOCOL_ERROR."
func TestRejectFrameOnIdle_WindowUpdate(t *testing.T) ***REMOVED***
	testRejectRequestWithProtocolError(t, func(st *serverTester) ***REMOVED***
		st.fr.WriteWindowUpdate(123, 456)
	***REMOVED***)
***REMOVED***
func TestRejectFrameOnIdle_Data(t *testing.T) ***REMOVED***
	testRejectRequestWithProtocolError(t, func(st *serverTester) ***REMOVED***
		st.fr.WriteData(123, true, nil)
	***REMOVED***)
***REMOVED***
func TestRejectFrameOnIdle_RSTStream(t *testing.T) ***REMOVED***
	testRejectRequestWithProtocolError(t, func(st *serverTester) ***REMOVED***
		st.fr.WriteRSTStream(123, ErrCodeCancel)
	***REMOVED***)
***REMOVED***

func TestServer_Request_Connect(t *testing.T) ***REMOVED***
	testServerRequest(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID: 1,
			BlockFragment: st.encodeHeaderRaw(
				":method", "CONNECT",
				":authority", "example.com:123",
			),
			EndStream:  true,
			EndHeaders: true,
		***REMOVED***)
	***REMOVED***, func(r *http.Request) ***REMOVED***
		if g, w := r.Method, "CONNECT"; g != w ***REMOVED***
			t.Errorf("Method = %q; want %q", g, w)
		***REMOVED***
		if g, w := r.RequestURI, "example.com:123"; g != w ***REMOVED***
			t.Errorf("RequestURI = %q; want %q", g, w)
		***REMOVED***
		if g, w := r.URL.Host, "example.com:123"; g != w ***REMOVED***
			t.Errorf("URL.Host = %q; want %q", g, w)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Request_Connect_InvalidPath(t *testing.T) ***REMOVED***
	testServerRejectsStream(t, ErrCodeProtocol, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID: 1,
			BlockFragment: st.encodeHeaderRaw(
				":method", "CONNECT",
				":authority", "example.com:123",
				":path", "/bogus",
			),
			EndStream:  true,
			EndHeaders: true,
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestServer_Request_Connect_InvalidScheme(t *testing.T) ***REMOVED***
	testServerRejectsStream(t, ErrCodeProtocol, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID: 1,
			BlockFragment: st.encodeHeaderRaw(
				":method", "CONNECT",
				":authority", "example.com:123",
				":scheme", "https",
			),
			EndStream:  true,
			EndHeaders: true,
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestServer_Ping(t *testing.T) ***REMOVED***
	st := newServerTester(t, nil)
	defer st.Close()
	st.greet()

	// Server should ignore this one, since it has ACK set.
	ackPingData := [8]byte***REMOVED***1, 2, 4, 8, 16, 32, 64, 128***REMOVED***
	if err := st.fr.WritePing(true, ackPingData); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// But the server should reply to this one, since ACK is false.
	pingData := [8]byte***REMOVED***1, 2, 3, 4, 5, 6, 7, 8***REMOVED***
	if err := st.fr.WritePing(false, pingData); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	pf := st.wantPing()
	if !pf.Flags.Has(FlagPingAck) ***REMOVED***
		t.Error("response ping doesn't have ACK set")
	***REMOVED***
	if pf.Data != pingData ***REMOVED***
		t.Errorf("response ping has data %q; want %q", pf.Data, pingData)
	***REMOVED***
***REMOVED***

func TestServer_RejectsLargeFrames(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("see golang.org/issue/13434")
	***REMOVED***

	st := newServerTester(t, nil)
	defer st.Close()
	st.greet()

	// Write too large of a frame (too large by one byte)
	// We ignore the return value because it's expected that the server
	// will only read the first 9 bytes (the headre) and then disconnect.
	st.fr.WriteRawFrame(0xff, 0, 0, make([]byte, defaultMaxReadFrameSize+1))

	gf := st.wantGoAway()
	if gf.ErrCode != ErrCodeFrameSize ***REMOVED***
		t.Errorf("GOAWAY err = %v; want %v", gf.ErrCode, ErrCodeFrameSize)
	***REMOVED***
	if st.serverLogBuf.Len() != 0 ***REMOVED***
		// Previously we spun here for a bit until the GOAWAY disconnect
		// timer fired, logging while we fired.
		t.Errorf("unexpected server output: %.500s\n", st.serverLogBuf.Bytes())
	***REMOVED***
***REMOVED***

func TestServer_Handler_Sends_WindowUpdate(t *testing.T) ***REMOVED***
	puppet := newHandlerPuppet()
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		puppet.act(w, r)
	***REMOVED***)
	defer st.Close()
	defer puppet.done()

	st.greet()

	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1, // clients send odd numbers
		BlockFragment: st.encodeHeader(":method", "POST"),
		EndStream:     false, // data coming
		EndHeaders:    true,
	***REMOVED***)
	st.writeData(1, false, []byte("abcdef"))
	puppet.do(readBodyHandler(t, "abc"))
	st.wantWindowUpdate(0, 3)
	st.wantWindowUpdate(1, 3)

	puppet.do(readBodyHandler(t, "def"))
	st.wantWindowUpdate(0, 3)
	st.wantWindowUpdate(1, 3)

	st.writeData(1, true, []byte("ghijkl")) // END_STREAM here
	puppet.do(readBodyHandler(t, "ghi"))
	puppet.do(readBodyHandler(t, "jkl"))
	st.wantWindowUpdate(0, 3)
	st.wantWindowUpdate(0, 3) // no more stream-level, since END_STREAM
***REMOVED***

// the version of the TestServer_Handler_Sends_WindowUpdate with padding.
// See golang.org/issue/16556
func TestServer_Handler_Sends_WindowUpdate_Padding(t *testing.T) ***REMOVED***
	puppet := newHandlerPuppet()
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		puppet.act(w, r)
	***REMOVED***)
	defer st.Close()
	defer puppet.done()

	st.greet()

	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: st.encodeHeader(":method", "POST"),
		EndStream:     false,
		EndHeaders:    true,
	***REMOVED***)
	st.writeDataPadded(1, false, []byte("abcdef"), []byte***REMOVED***0, 0, 0, 0***REMOVED***)

	// Expect to immediately get our 5 bytes of padding back for
	// both the connection and stream (4 bytes of padding + 1 byte of length)
	st.wantWindowUpdate(0, 5)
	st.wantWindowUpdate(1, 5)

	puppet.do(readBodyHandler(t, "abc"))
	st.wantWindowUpdate(0, 3)
	st.wantWindowUpdate(1, 3)

	puppet.do(readBodyHandler(t, "def"))
	st.wantWindowUpdate(0, 3)
	st.wantWindowUpdate(1, 3)
***REMOVED***

func TestServer_Send_GoAway_After_Bogus_WindowUpdate(t *testing.T) ***REMOVED***
	st := newServerTester(t, nil)
	defer st.Close()
	st.greet()
	if err := st.fr.WriteWindowUpdate(0, 1<<31-1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	gf := st.wantGoAway()
	if gf.ErrCode != ErrCodeFlowControl ***REMOVED***
		t.Errorf("GOAWAY err = %v; want %v", gf.ErrCode, ErrCodeFlowControl)
	***REMOVED***
	if gf.LastStreamID != 0 ***REMOVED***
		t.Errorf("GOAWAY last stream ID = %v; want %v", gf.LastStreamID, 0)
	***REMOVED***
***REMOVED***

func TestServer_Send_RstStream_After_Bogus_WindowUpdate(t *testing.T) ***REMOVED***
	inHandler := make(chan bool)
	blockHandler := make(chan bool)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		inHandler <- true
		<-blockHandler
	***REMOVED***)
	defer st.Close()
	defer close(blockHandler)
	st.greet()
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: st.encodeHeader(":method", "POST"),
		EndStream:     false, // keep it open
		EndHeaders:    true,
	***REMOVED***)
	<-inHandler
	// Send a bogus window update:
	if err := st.fr.WriteWindowUpdate(1, 1<<31-1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	st.wantRSTStream(1, ErrCodeFlowControl)
***REMOVED***

// testServerPostUnblock sends a hanging POST with unsent data to handler,
// then runs fn once in the handler, and verifies that the error returned from
// handler is acceptable. It fails if takes over 5 seconds for handler to exit.
func testServerPostUnblock(t *testing.T,
	handler func(http.ResponseWriter, *http.Request) error,
	fn func(*serverTester),
	checkErr func(error),
	otherHeaders ...string) ***REMOVED***
	inHandler := make(chan bool)
	errc := make(chan error, 1)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		inHandler <- true
		errc <- handler(w, r)
	***REMOVED***)
	defer st.Close()
	st.greet()
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: st.encodeHeader(append([]string***REMOVED***":method", "POST"***REMOVED***, otherHeaders...)...),
		EndStream:     false, // keep it open
		EndHeaders:    true,
	***REMOVED***)
	<-inHandler
	fn(st)
	select ***REMOVED***
	case err := <-errc:
		if checkErr != nil ***REMOVED***
			checkErr(err)
		***REMOVED***
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for Handler to return")
	***REMOVED***
***REMOVED***

func TestServer_RSTStream_Unblocks_Read(t *testing.T) ***REMOVED***
	testServerPostUnblock(t,
		func(w http.ResponseWriter, r *http.Request) (err error) ***REMOVED***
			_, err = r.Body.Read(make([]byte, 1))
			return
		***REMOVED***,
		func(st *serverTester) ***REMOVED***
			if err := st.fr.WriteRSTStream(1, ErrCodeCancel); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***,
		func(err error) ***REMOVED***
			want := StreamError***REMOVED***StreamID: 0x1, Code: 0x8***REMOVED***
			if !reflect.DeepEqual(err, want) ***REMOVED***
				t.Errorf("Read error = %v; want %v", err, want)
			***REMOVED***
		***REMOVED***,
	)
***REMOVED***

func TestServer_RSTStream_Unblocks_Header_Write(t *testing.T) ***REMOVED***
	// Run this test a bunch, because it doesn't always
	// deadlock. But with a bunch, it did.
	n := 50
	if testing.Short() ***REMOVED***
		n = 5
	***REMOVED***
	for i := 0; i < n; i++ ***REMOVED***
		testServer_RSTStream_Unblocks_Header_Write(t)
	***REMOVED***
***REMOVED***

func testServer_RSTStream_Unblocks_Header_Write(t *testing.T) ***REMOVED***
	inHandler := make(chan bool, 1)
	unblockHandler := make(chan bool, 1)
	headerWritten := make(chan bool, 1)
	wroteRST := make(chan bool, 1)

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		inHandler <- true
		<-wroteRST
		w.Header().Set("foo", "bar")
		w.WriteHeader(200)
		w.(http.Flusher).Flush()
		headerWritten <- true
		<-unblockHandler
	***REMOVED***)
	defer st.Close()

	st.greet()
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: st.encodeHeader(":method", "POST"),
		EndStream:     false, // keep it open
		EndHeaders:    true,
	***REMOVED***)
	<-inHandler
	if err := st.fr.WriteRSTStream(1, ErrCodeCancel); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	wroteRST <- true
	st.awaitIdle()
	select ***REMOVED***
	case <-headerWritten:
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for header write")
	***REMOVED***
	unblockHandler <- true
***REMOVED***

func TestServer_DeadConn_Unblocks_Read(t *testing.T) ***REMOVED***
	testServerPostUnblock(t,
		func(w http.ResponseWriter, r *http.Request) (err error) ***REMOVED***
			_, err = r.Body.Read(make([]byte, 1))
			return
		***REMOVED***,
		func(st *serverTester) ***REMOVED*** st.cc.Close() ***REMOVED***,
		func(err error) ***REMOVED***
			if err == nil ***REMOVED***
				t.Error("unexpected nil error from Request.Body.Read")
			***REMOVED***
		***REMOVED***,
	)
***REMOVED***

var blockUntilClosed = func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
	<-w.(http.CloseNotifier).CloseNotify()
	return nil
***REMOVED***

func TestServer_CloseNotify_After_RSTStream(t *testing.T) ***REMOVED***
	testServerPostUnblock(t, blockUntilClosed, func(st *serverTester) ***REMOVED***
		if err := st.fr.WriteRSTStream(1, ErrCodeCancel); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***, nil)
***REMOVED***

func TestServer_CloseNotify_After_ConnClose(t *testing.T) ***REMOVED***
	testServerPostUnblock(t, blockUntilClosed, func(st *serverTester) ***REMOVED*** st.cc.Close() ***REMOVED***, nil)
***REMOVED***

// that CloseNotify unblocks after a stream error due to the client's
// problem that's unrelated to them explicitly canceling it (which is
// TestServer_CloseNotify_After_RSTStream above)
func TestServer_CloseNotify_After_StreamError(t *testing.T) ***REMOVED***
	testServerPostUnblock(t, blockUntilClosed, func(st *serverTester) ***REMOVED***
		// data longer than declared Content-Length => stream error
		st.writeData(1, true, []byte("1234"))
	***REMOVED***, nil, "content-length", "3")
***REMOVED***

func TestServer_StateTransitions(t *testing.T) ***REMOVED***
	var st *serverTester
	inHandler := make(chan bool)
	writeData := make(chan bool)
	leaveHandler := make(chan bool)
	st = newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		inHandler <- true
		if st.stream(1) == nil ***REMOVED***
			t.Errorf("nil stream 1 in handler")
		***REMOVED***
		if got, want := st.streamState(1), stateOpen; got != want ***REMOVED***
			t.Errorf("in handler, state is %v; want %v", got, want)
		***REMOVED***
		writeData <- true
		if n, err := r.Body.Read(make([]byte, 1)); n != 0 || err != io.EOF ***REMOVED***
			t.Errorf("body read = %d, %v; want 0, EOF", n, err)
		***REMOVED***
		if got, want := st.streamState(1), stateHalfClosedRemote; got != want ***REMOVED***
			t.Errorf("in handler, state is %v; want %v", got, want)
		***REMOVED***

		<-leaveHandler
	***REMOVED***)
	st.greet()
	if st.stream(1) != nil ***REMOVED***
		t.Fatal("stream 1 should be empty")
	***REMOVED***
	if got := st.streamState(1); got != stateIdle ***REMOVED***
		t.Fatalf("stream 1 should be idle; got %v", got)
	***REMOVED***

	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: st.encodeHeader(":method", "POST"),
		EndStream:     false, // keep it open
		EndHeaders:    true,
	***REMOVED***)
	<-inHandler
	<-writeData
	st.writeData(1, true, nil)

	leaveHandler <- true
	hf := st.wantHeaders()
	if !hf.StreamEnded() ***REMOVED***
		t.Fatal("expected END_STREAM flag")
	***REMOVED***

	if got, want := st.streamState(1), stateClosed; got != want ***REMOVED***
		t.Errorf("at end, state is %v; want %v", got, want)
	***REMOVED***
	if st.stream(1) != nil ***REMOVED***
		t.Fatal("at end, stream 1 should be gone")
	***REMOVED***
***REMOVED***

// test HEADERS w/o EndHeaders + another HEADERS (should get rejected)
func TestServer_Rejects_HeadersNoEnd_Then_Headers(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    false,
		***REMOVED***)
		st.writeHeaders(HeadersFrameParam***REMOVED*** // Not a continuation.
			StreamID:      3, // different stream.
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***)
***REMOVED***

// test HEADERS w/o EndHeaders + PING (should get rejected)
func TestServer_Rejects_HeadersNoEnd_Then_Ping(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    false,
		***REMOVED***)
		if err := st.fr.WritePing(false, [8]byte***REMOVED******REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// test HEADERS w/ EndHeaders + a continuation HEADERS (should get rejected)
func TestServer_Rejects_HeadersEnd_Then_Continuation(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
		st.wantHeaders()
		if err := st.fr.WriteContinuation(1, true, encodeHeaderNoImplicit(t, "foo", "bar")); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// test HEADERS w/o EndHeaders + a continuation HEADERS on wrong stream ID
func TestServer_Rejects_HeadersNoEnd_Then_ContinuationWrongStream(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    false,
		***REMOVED***)
		if err := st.fr.WriteContinuation(3, true, encodeHeaderNoImplicit(t, "foo", "bar")); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// No HEADERS on stream 0.
func TestServer_Rejects_Headers0(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		st.fr.AllowIllegalWrites = true
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      0,
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***)
***REMOVED***

// No CONTINUATION on stream 0.
func TestServer_Rejects_Continuation0(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		st.fr.AllowIllegalWrites = true
		if err := st.fr.WriteContinuation(0, true, st.encodeHeader()); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// No PRIORITY on stream 0.
func TestServer_Rejects_Priority0(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		st.fr.AllowIllegalWrites = true
		st.writePriority(0, PriorityParam***REMOVED***StreamDep: 1***REMOVED***)
	***REMOVED***)
***REMOVED***

// No HEADERS frame with a self-dependence.
func TestServer_Rejects_HeadersSelfDependence(t *testing.T) ***REMOVED***
	testServerRejectsStream(t, ErrCodeProtocol, func(st *serverTester) ***REMOVED***
		st.fr.AllowIllegalWrites = true
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    true,
			Priority:      PriorityParam***REMOVED***StreamDep: 1***REMOVED***,
		***REMOVED***)
	***REMOVED***)
***REMOVED***

// No PRIORTY frame with a self-dependence.
func TestServer_Rejects_PrioritySelfDependence(t *testing.T) ***REMOVED***
	testServerRejectsStream(t, ErrCodeProtocol, func(st *serverTester) ***REMOVED***
		st.fr.AllowIllegalWrites = true
		st.writePriority(1, PriorityParam***REMOVED***StreamDep: 1***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestServer_Rejects_PushPromise(t *testing.T) ***REMOVED***
	testServerRejectsConn(t, func(st *serverTester) ***REMOVED***
		pp := PushPromiseParam***REMOVED***
			StreamID:  1,
			PromiseID: 3,
		***REMOVED***
		if err := st.fr.WritePushPromise(pp); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// testServerRejectsConn tests that the server hangs up with a GOAWAY
// frame and a server close after the client does something
// deserving a CONNECTION_ERROR.
func testServerRejectsConn(t *testing.T, writeReq func(*serverTester)) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED******REMOVED***)
	st.addLogFilter("connection error: PROTOCOL_ERROR")
	defer st.Close()
	st.greet()
	writeReq(st)

	st.wantGoAway()
	errc := make(chan error, 1)
	go func() ***REMOVED***
		fr, err := st.fr.ReadFrame()
		if err == nil ***REMOVED***
			err = fmt.Errorf("got frame of type %T", fr)
		***REMOVED***
		errc <- err
	***REMOVED***()
	select ***REMOVED***
	case err := <-errc:
		if err != io.EOF ***REMOVED***
			t.Errorf("ReadFrame = %v; want io.EOF", err)
		***REMOVED***
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for disconnect")
	***REMOVED***
***REMOVED***

// testServerRejectsStream tests that the server sends a RST_STREAM with the provided
// error code after a client sends a bogus request.
func testServerRejectsStream(t *testing.T, code ErrCode, writeReq func(*serverTester)) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED******REMOVED***)
	defer st.Close()
	st.greet()
	writeReq(st)
	st.wantRSTStream(1, code)
***REMOVED***

// testServerRequest sets up an idle HTTP/2 connection and lets you
// write a single request with writeReq, and then verify that the
// *http.Request is built correctly in checkReq.
func testServerRequest(t *testing.T, writeReq func(*serverTester), checkReq func(*http.Request)) ***REMOVED***
	gotReq := make(chan bool, 1)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Body == nil ***REMOVED***
			t.Fatal("nil Body")
		***REMOVED***
		checkReq(r)
		gotReq <- true
	***REMOVED***)
	defer st.Close()

	st.greet()
	writeReq(st)

	select ***REMOVED***
	case <-gotReq:
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for request")
	***REMOVED***
***REMOVED***

func getSlash(st *serverTester) ***REMOVED*** st.bodylessReq1() ***REMOVED***

func TestServer_Response_NoData(t *testing.T) ***REMOVED***
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		// Nothing.
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if !hf.StreamEnded() ***REMOVED***
			t.Fatal("want END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_NoData_Header_FooBar(t *testing.T) ***REMOVED***
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.Header().Set("Foo-Bar", "some-value")
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if !hf.StreamEnded() ***REMOVED***
			t.Fatal("want END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"foo-bar", "some-value"***REMOVED***,
			***REMOVED***"content-length", "0"***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_Data_Sniff_DoesntOverride(t *testing.T) ***REMOVED***
	const msg = "<html>this is HTML."
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.Header().Set("Content-Type", "foo/bar")
		io.WriteString(w, msg)
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("don't want END_STREAM, expecting data")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"content-type", "foo/bar"***REMOVED***,
			***REMOVED***"content-length", strconv.Itoa(len(msg))***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
		df := st.wantData()
		if !df.StreamEnded() ***REMOVED***
			t.Error("expected DATA to have END_STREAM flag")
		***REMOVED***
		if got := string(df.Data()); got != msg ***REMOVED***
			t.Errorf("got DATA %q; want %q", got, msg)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_TransferEncoding_chunked(t *testing.T) ***REMOVED***
	const msg = "hi"
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.Header().Set("Transfer-Encoding", "chunked") // should be stripped
		io.WriteString(w, msg)
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"content-type", "text/plain; charset=utf-8"***REMOVED***,
			***REMOVED***"content-length", strconv.Itoa(len(msg))***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Header accessed only after the initial write.
func TestServer_Response_Data_IgnoreHeaderAfterWrite_After(t *testing.T) ***REMOVED***
	const msg = "<html>this is HTML."
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		io.WriteString(w, msg)
		w.Header().Set("foo", "should be ignored")
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"content-type", "text/html; charset=utf-8"***REMOVED***,
			***REMOVED***"content-length", strconv.Itoa(len(msg))***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Header accessed before the initial write and later mutated.
func TestServer_Response_Data_IgnoreHeaderAfterWrite_Overwrite(t *testing.T) ***REMOVED***
	const msg = "<html>this is HTML."
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.Header().Set("foo", "proper value")
		io.WriteString(w, msg)
		w.Header().Set("foo", "should be ignored")
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"foo", "proper value"***REMOVED***,
			***REMOVED***"content-type", "text/html; charset=utf-8"***REMOVED***,
			***REMOVED***"content-length", strconv.Itoa(len(msg))***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_Data_SniffLenType(t *testing.T) ***REMOVED***
	const msg = "<html>this is HTML."
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		io.WriteString(w, msg)
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("don't want END_STREAM, expecting data")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"content-type", "text/html; charset=utf-8"***REMOVED***,
			***REMOVED***"content-length", strconv.Itoa(len(msg))***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
		df := st.wantData()
		if !df.StreamEnded() ***REMOVED***
			t.Error("expected DATA to have END_STREAM flag")
		***REMOVED***
		if got := string(df.Data()); got != msg ***REMOVED***
			t.Errorf("got DATA %q; want %q", got, msg)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_Header_Flush_MidWrite(t *testing.T) ***REMOVED***
	const msg = "<html>this is HTML"
	const msg2 = ", and this is the next chunk"
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		io.WriteString(w, msg)
		w.(http.Flusher).Flush()
		io.WriteString(w, msg2)
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"content-type", "text/html; charset=utf-8"***REMOVED***, // sniffed
			// and no content-length
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
		***REMOVED***
			df := st.wantData()
			if df.StreamEnded() ***REMOVED***
				t.Error("unexpected END_STREAM flag")
			***REMOVED***
			if got := string(df.Data()); got != msg ***REMOVED***
				t.Errorf("got DATA %q; want %q", got, msg)
			***REMOVED***
		***REMOVED***
		***REMOVED***
			df := st.wantData()
			if !df.StreamEnded() ***REMOVED***
				t.Error("wanted END_STREAM flag on last data chunk")
			***REMOVED***
			if got := string(df.Data()); got != msg2 ***REMOVED***
				t.Errorf("got DATA %q; want %q", got, msg2)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_LargeWrite(t *testing.T) ***REMOVED***
	const size = 1 << 20
	const maxFrameSize = 16 << 10
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		n, err := w.Write(bytes.Repeat([]byte("a"), size))
		if err != nil ***REMOVED***
			return fmt.Errorf("Write error: %v", err)
		***REMOVED***
		if n != size ***REMOVED***
			return fmt.Errorf("wrong size %d from Write", n)
		***REMOVED***
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		if err := st.fr.WriteSettings(
			Setting***REMOVED***SettingInitialWindowSize, 0***REMOVED***,
			Setting***REMOVED***SettingMaxFrameSize, maxFrameSize***REMOVED***,
		); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		st.wantSettingsAck()

		getSlash(st) // make the single request

		// Give the handler quota to write:
		if err := st.fr.WriteWindowUpdate(1, size); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		// Give the handler quota to write to connection-level
		// window as well
		if err := st.fr.WriteWindowUpdate(0, size); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"content-type", "text/plain; charset=utf-8"***REMOVED***, // sniffed
			// and no content-length
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***
		var bytes, frames int
		for ***REMOVED***
			df := st.wantData()
			bytes += len(df.Data())
			frames++
			for _, b := range df.Data() ***REMOVED***
				if b != 'a' ***REMOVED***
					t.Fatal("non-'a' byte seen in DATA")
				***REMOVED***
			***REMOVED***
			if df.StreamEnded() ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if bytes != size ***REMOVED***
			t.Errorf("Got %d bytes; want %d", bytes, size)
		***REMOVED***
		if want := int(size / maxFrameSize); frames < want || frames > want*2 ***REMOVED***
			t.Errorf("Got %d frames; want %d", frames, size)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Test that the handler can't write more than the client allows
func TestServer_Response_LargeWrite_FlowControlled(t *testing.T) ***REMOVED***
	// Make these reads. Before each read, the client adds exactly enough
	// flow-control to satisfy the read. Numbers chosen arbitrarily.
	reads := []int***REMOVED***123, 1, 13, 127***REMOVED***
	size := 0
	for _, n := range reads ***REMOVED***
		size += n
	***REMOVED***

	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.(http.Flusher).Flush()
		n, err := w.Write(bytes.Repeat([]byte("a"), size))
		if err != nil ***REMOVED***
			return fmt.Errorf("Write error: %v", err)
		***REMOVED***
		if n != size ***REMOVED***
			return fmt.Errorf("wrong size %d from Write", n)
		***REMOVED***
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		// Set the window size to something explicit for this test.
		// It's also how much initial data we expect.
		if err := st.fr.WriteSettings(Setting***REMOVED***SettingInitialWindowSize, uint32(reads[0])***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		st.wantSettingsAck()

		getSlash(st) // make the single request

		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***

		df := st.wantData()
		if got := len(df.Data()); got != reads[0] ***REMOVED***
			t.Fatalf("Initial window size = %d but got DATA with %d bytes", reads[0], got)
		***REMOVED***

		for _, quota := range reads[1:] ***REMOVED***
			if err := st.fr.WriteWindowUpdate(1, uint32(quota)); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			df := st.wantData()
			if int(quota) != len(df.Data()) ***REMOVED***
				t.Fatalf("read %d bytes after giving %d quota", len(df.Data()), quota)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Test that the handler blocked in a Write is unblocked if the server sends a RST_STREAM.
func TestServer_Response_RST_Unblocks_LargeWrite(t *testing.T) ***REMOVED***
	const size = 1 << 20
	const maxFrameSize = 16 << 10
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.(http.Flusher).Flush()
		errc := make(chan error, 1)
		go func() ***REMOVED***
			_, err := w.Write(bytes.Repeat([]byte("a"), size))
			errc <- err
		***REMOVED***()
		select ***REMOVED***
		case err := <-errc:
			if err == nil ***REMOVED***
				return errors.New("unexpected nil error from Write in handler")
			***REMOVED***
			return nil
		case <-time.After(2 * time.Second):
			return errors.New("timeout waiting for Write in handler")
		***REMOVED***
	***REMOVED***, func(st *serverTester) ***REMOVED***
		if err := st.fr.WriteSettings(
			Setting***REMOVED***SettingInitialWindowSize, 0***REMOVED***,
			Setting***REMOVED***SettingMaxFrameSize, maxFrameSize***REMOVED***,
		); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		st.wantSettingsAck()

		getSlash(st) // make the single request

		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***

		if err := st.fr.WriteRSTStream(1, ErrCodeCancel); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_Empty_Data_Not_FlowControlled(t *testing.T) ***REMOVED***
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.(http.Flusher).Flush()
		// Nothing; send empty DATA
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		// Handler gets no data quota:
		if err := st.fr.WriteSettings(Setting***REMOVED***SettingInitialWindowSize, 0***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		st.wantSettingsAck()

		getSlash(st) // make the single request

		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***

		df := st.wantData()
		if got := len(df.Data()); got != 0 ***REMOVED***
			t.Fatalf("unexpected %d DATA bytes; want 0", got)
		***REMOVED***
		if !df.StreamEnded() ***REMOVED***
			t.Fatal("DATA didn't have END_STREAM")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Response_Automatic100Continue(t *testing.T) ***REMOVED***
	const msg = "foo"
	const reply = "bar"
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		if v := r.Header.Get("Expect"); v != "" ***REMOVED***
			t.Errorf("Expect header = %q; want empty", v)
		***REMOVED***
		buf := make([]byte, len(msg))
		// This read should trigger the 100-continue being sent.
		if n, err := io.ReadFull(r.Body, buf); err != nil || n != len(msg) || string(buf) != msg ***REMOVED***
			return fmt.Errorf("ReadFull = %q, %v; want %q, nil", buf[:n], err, msg)
		***REMOVED***
		_, err := io.WriteString(w, reply)
		return err
	***REMOVED***, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader(":method", "POST", "expect", "100-continue"),
			EndStream:     false,
			EndHeaders:    true,
		***REMOVED***)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "100"***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Fatalf("Got headers %v; want %v", goth, wanth)
		***REMOVED***

		// Okay, they sent status 100, so we can send our
		// gigantic and/or sensitive "foo" payload now.
		st.writeData(1, true, []byte(msg))

		st.wantWindowUpdate(0, uint32(len(msg)))

		hf = st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("expected data to follow")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		goth = st.decodeHeader(hf.HeaderBlockFragment())
		wanth = [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"content-type", "text/plain; charset=utf-8"***REMOVED***,
			***REMOVED***"content-length", strconv.Itoa(len(reply))***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Got headers %v; want %v", goth, wanth)
		***REMOVED***

		df := st.wantData()
		if string(df.Data()) != reply ***REMOVED***
			t.Errorf("Client read %q; want %q", df.Data(), reply)
		***REMOVED***
		if !df.StreamEnded() ***REMOVED***
			t.Errorf("expect data stream end")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_HandlerWriteErrorOnDisconnect(t *testing.T) ***REMOVED***
	errc := make(chan error, 1)
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		p := []byte("some data.\n")
		for ***REMOVED***
			_, err := w.Write(p)
			if err != nil ***REMOVED***
				errc <- err
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: st.encodeHeader(),
			EndStream:     false,
			EndHeaders:    true,
		***REMOVED***)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("unexpected END_STREAM flag")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("want END_HEADERS flag")
		***REMOVED***
		// Close the connection and wait for the handler to (hopefully) notice.
		st.cc.Close()
		select ***REMOVED***
		case <-errc:
		case <-time.After(5 * time.Second):
			t.Error("timeout")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Rejects_Too_Many_Streams(t *testing.T) ***REMOVED***
	const testPath = "/some/path"

	inHandler := make(chan uint32)
	leaveHandler := make(chan bool)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		id := w.(*responseWriter).rws.stream.id
		inHandler <- id
		if id == 1+(defaultMaxStreams+1)*2 && r.URL.Path != testPath ***REMOVED***
			t.Errorf("decoded final path as %q; want %q", r.URL.Path, testPath)
		***REMOVED***
		<-leaveHandler
	***REMOVED***)
	defer st.Close()
	st.greet()
	nextStreamID := uint32(1)
	streamID := func() uint32 ***REMOVED***
		defer func() ***REMOVED*** nextStreamID += 2 ***REMOVED***()
		return nextStreamID
	***REMOVED***
	sendReq := func(id uint32, headers ...string) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      id,
			BlockFragment: st.encodeHeader(headers...),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
	***REMOVED***
	for i := 0; i < defaultMaxStreams; i++ ***REMOVED***
		sendReq(streamID())
		<-inHandler
	***REMOVED***
	defer func() ***REMOVED***
		for i := 0; i < defaultMaxStreams; i++ ***REMOVED***
			leaveHandler <- true
		***REMOVED***
	***REMOVED***()

	// And this one should cross the limit:
	// (It's also sent as a CONTINUATION, to verify we still track the decoder context,
	// even if we're rejecting it)
	rejectID := streamID()
	headerBlock := st.encodeHeader(":path", testPath)
	frag1, frag2 := headerBlock[:3], headerBlock[3:]
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      rejectID,
		BlockFragment: frag1,
		EndStream:     true,
		EndHeaders:    false, // CONTINUATION coming
	***REMOVED***)
	if err := st.fr.WriteContinuation(rejectID, true, frag2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	st.wantRSTStream(rejectID, ErrCodeProtocol)

	// But let a handler finish:
	leaveHandler <- true
	st.wantHeaders()

	// And now another stream should be able to start:
	goodID := streamID()
	sendReq(goodID, ":path", testPath)
	select ***REMOVED***
	case got := <-inHandler:
		if got != goodID ***REMOVED***
			t.Errorf("Got stream %d; want %d", got, goodID)
		***REMOVED***
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for handler")
	***REMOVED***
***REMOVED***

// So many response headers that the server needs to use CONTINUATION frames:
func TestServer_Response_ManyHeaders_With_Continuation(t *testing.T) ***REMOVED***
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		h := w.Header()
		for i := 0; i < 5000; i++ ***REMOVED***
			h.Set(fmt.Sprintf("x-header-%d", i), fmt.Sprintf("x-value-%d", i))
		***REMOVED***
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if hf.HeadersEnded() ***REMOVED***
			t.Fatal("got unwanted END_HEADERS flag")
		***REMOVED***
		n := 0
		for ***REMOVED***
			n++
			cf := st.wantContinuation()
			if cf.HeadersEnded() ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if n < 5 ***REMOVED***
			t.Errorf("Only got %d CONTINUATION frames; expected 5+ (currently 6)", n)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// This previously crashed (reported by Mathieu Lonjaret as observed
// while using Camlistore) because we got a DATA frame from the client
// after the handler exited and our logic at the time was wrong,
// keeping a stream in the map in stateClosed, which tickled an
// invariant check later when we tried to remove that stream (via
// defer sc.closeAllStreamsOnConnClose) when the serverConn serve loop
// ended.
func TestServer_NoCrash_HandlerClose_Then_ClientClose(t *testing.T) ***REMOVED***
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		// nothing
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: st.encodeHeader(),
			EndStream:     false, // DATA is coming
			EndHeaders:    true,
		***REMOVED***)
		hf := st.wantHeaders()
		if !hf.HeadersEnded() || !hf.StreamEnded() ***REMOVED***
			t.Fatalf("want END_HEADERS+END_STREAM, got %v", hf)
		***REMOVED***

		// Sent when the a Handler closes while a client has
		// indicated it's still sending DATA:
		st.wantRSTStream(1, ErrCodeNo)

		// Now the handler has ended, so it's ended its
		// stream, but the client hasn't closed its side
		// (stateClosedLocal).  So send more data and verify
		// it doesn't crash with an internal invariant panic, like
		// it did before.
		st.writeData(1, true, []byte("foo"))

		// Get our flow control bytes back, since the handler didn't get them.
		st.wantWindowUpdate(0, uint32(len("foo")))

		// Sent after a peer sends data anyway (admittedly the
		// previous RST_STREAM might've still been in-flight),
		// but they'll get the more friendly 'cancel' code
		// first.
		st.wantRSTStream(1, ErrCodeStreamClosed)

		// Set up a bunch of machinery to record the panic we saw
		// previously.
		var (
			panMu    sync.Mutex
			panicVal interface***REMOVED******REMOVED***
		)

		testHookOnPanicMu.Lock()
		testHookOnPanic = func(sc *serverConn, pv interface***REMOVED******REMOVED***) bool ***REMOVED***
			panMu.Lock()
			panicVal = pv
			panMu.Unlock()
			return true
		***REMOVED***
		testHookOnPanicMu.Unlock()

		// Now force the serve loop to end, via closing the connection.
		st.cc.Close()
		select ***REMOVED***
		case <-st.sc.doneServing:
			// Loop has exited.
			panMu.Lock()
			got := panicVal
			panMu.Unlock()
			if got != nil ***REMOVED***
				t.Errorf("Got panic: %v", got)
			***REMOVED***
		case <-time.After(5 * time.Second):
			t.Error("timeout")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestServer_Rejects_TLS10(t *testing.T) ***REMOVED*** testRejectTLS(t, tls.VersionTLS10) ***REMOVED***
func TestServer_Rejects_TLS11(t *testing.T) ***REMOVED*** testRejectTLS(t, tls.VersionTLS11) ***REMOVED***

func testRejectTLS(t *testing.T, max uint16) ***REMOVED***
	st := newServerTester(t, nil, func(c *tls.Config) ***REMOVED***
		c.MaxVersion = max
	***REMOVED***)
	defer st.Close()
	gf := st.wantGoAway()
	if got, want := gf.ErrCode, ErrCodeInadequateSecurity; got != want ***REMOVED***
		t.Errorf("Got error code %v; want %v", got, want)
	***REMOVED***
***REMOVED***

func TestServer_Rejects_TLSBadCipher(t *testing.T) ***REMOVED***
	st := newServerTester(t, nil, func(c *tls.Config) ***REMOVED***
		// Only list bad ones:
		c.CipherSuites = []uint16***REMOVED***
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			cipher_TLS_RSA_WITH_AES_128_CBC_SHA256,
		***REMOVED***
	***REMOVED***)
	defer st.Close()
	gf := st.wantGoAway()
	if got, want := gf.ErrCode, ErrCodeInadequateSecurity; got != want ***REMOVED***
		t.Errorf("Got error code %v; want %v", got, want)
	***REMOVED***
***REMOVED***

func TestServer_Advertises_Common_Cipher(t *testing.T) ***REMOVED***
	const requiredSuite = tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	st := newServerTester(t, nil, func(c *tls.Config) ***REMOVED***
		// Have the client only support the one required by the spec.
		c.CipherSuites = []uint16***REMOVED***requiredSuite***REMOVED***
	***REMOVED***, func(ts *httptest.Server) ***REMOVED***
		var srv *http.Server = ts.Config
		// Have the server configured with no specific cipher suites.
		// This tests that Go's defaults include the required one.
		srv.TLSConfig = nil
	***REMOVED***)
	defer st.Close()
	st.greet()
***REMOVED***

func (st *serverTester) onHeaderField(f hpack.HeaderField) ***REMOVED***
	if f.Name == "date" ***REMOVED***
		return
	***REMOVED***
	st.decodedHeaders = append(st.decodedHeaders, [2]string***REMOVED***f.Name, f.Value***REMOVED***)
***REMOVED***

func (st *serverTester) decodeHeader(headerBlock []byte) (pairs [][2]string) ***REMOVED***
	st.decodedHeaders = nil
	if _, err := st.hpackDec.Write(headerBlock); err != nil ***REMOVED***
		st.t.Fatalf("hpack decoding error: %v", err)
	***REMOVED***
	if err := st.hpackDec.Close(); err != nil ***REMOVED***
		st.t.Fatalf("hpack decoding error: %v", err)
	***REMOVED***
	return st.decodedHeaders
***REMOVED***

// testServerResponse sets up an idle HTTP/2 connection. The client function should
// write a single request that must be handled by the handler. This waits up to 5s
// for client to return, then up to an additional 2s for the handler to return.
func testServerResponse(t testing.TB,
	handler func(http.ResponseWriter, *http.Request) error,
	client func(*serverTester),
) ***REMOVED***
	errc := make(chan error, 1)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Body == nil ***REMOVED***
			t.Fatal("nil Body")
		***REMOVED***
		errc <- handler(w, r)
	***REMOVED***)
	defer st.Close()

	donec := make(chan bool)
	go func() ***REMOVED***
		defer close(donec)
		st.greet()
		client(st)
	***REMOVED***()

	select ***REMOVED***
	case <-donec:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout in client")
	***REMOVED***

	select ***REMOVED***
	case err := <-errc:
		if err != nil ***REMOVED***
			t.Fatalf("Error in handler: %v", err)
		***REMOVED***
	case <-time.After(2 * time.Second):
		t.Fatal("timeout in handler")
	***REMOVED***
***REMOVED***

// readBodyHandler returns an http Handler func that reads len(want)
// bytes from r.Body and fails t if the contents read were not
// the value of want.
func readBodyHandler(t *testing.T, want string) func(w http.ResponseWriter, r *http.Request) ***REMOVED***
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		buf := make([]byte, len(want))
		_, err := io.ReadFull(r.Body, buf)
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if string(buf) != want ***REMOVED***
			t.Errorf("read %q; want %q", buf, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestServerWithCurl currently fails, hence the LenientCipherSuites test. See:
//   https://github.com/tatsuhiro-t/nghttp2/issues/140 &
//   http://sourceforge.net/p/curl/bugs/1472/
func TestServerWithCurl(t *testing.T)                     ***REMOVED*** testServerWithCurl(t, false) ***REMOVED***
func TestServerWithCurl_LenientCipherSuites(t *testing.T) ***REMOVED*** testServerWithCurl(t, true) ***REMOVED***

func testServerWithCurl(t *testing.T, permitProhibitedCipherSuites bool) ***REMOVED***
	if runtime.GOOS != "linux" ***REMOVED***
		t.Skip("skipping Docker test when not on Linux; requires --net which won't work with boot2docker anyway")
	***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping curl test in short mode")
	***REMOVED***
	requireCurl(t)
	var gotConn int32
	testHookOnConn = func() ***REMOVED*** atomic.StoreInt32(&gotConn, 1) ***REMOVED***

	const msg = "Hello from curl!\n"
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Foo", "Bar")
		w.Header().Set("Client-Proto", r.Proto)
		io.WriteString(w, msg)
	***REMOVED***))
	ConfigureServer(ts.Config, &Server***REMOVED***
		PermitProhibitedCipherSuites: permitProhibitedCipherSuites,
	***REMOVED***)
	ts.TLS = ts.Config.TLSConfig // the httptest.Server has its own copy of this TLS config
	ts.StartTLS()
	defer ts.Close()

	t.Logf("Running test server for curl to hit at: %s", ts.URL)
	container := curl(t, "--silent", "--http2", "--insecure", "-v", ts.URL)
	defer kill(container)
	resc := make(chan interface***REMOVED******REMOVED***, 1)
	go func() ***REMOVED***
		res, err := dockerLogs(container)
		if err != nil ***REMOVED***
			resc <- err
		***REMOVED*** else ***REMOVED***
			resc <- res
		***REMOVED***
	***REMOVED***()
	select ***REMOVED***
	case res := <-resc:
		if err, ok := res.(error); ok ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		body := string(res.([]byte))
		// Search for both "key: value" and "key:value", since curl changed their format
		// Our Dockerfile contains the latest version (no space), but just in case people
		// didn't rebuild, check both.
		if !strings.Contains(body, "foo: Bar") && !strings.Contains(body, "foo:Bar") ***REMOVED***
			t.Errorf("didn't see foo: Bar header")
			t.Logf("Got: %s", body)
		***REMOVED***
		if !strings.Contains(body, "client-proto: HTTP/2") && !strings.Contains(body, "client-proto:HTTP/2") ***REMOVED***
			t.Errorf("didn't see client-proto: HTTP/2 header")
			t.Logf("Got: %s", res)
		***REMOVED***
		if !strings.Contains(string(res.([]byte)), msg) ***REMOVED***
			t.Errorf("didn't see %q content", msg)
			t.Logf("Got: %s", res)
		***REMOVED***
	case <-time.After(3 * time.Second):
		t.Errorf("timeout waiting for curl")
	***REMOVED***

	if atomic.LoadInt32(&gotConn) == 0 ***REMOVED***
		t.Error("never saw an http2 connection")
	***REMOVED***
***REMOVED***

var doh2load = flag.Bool("h2load", false, "Run h2load test")

func TestServerWithH2Load(t *testing.T) ***REMOVED***
	if !*doh2load ***REMOVED***
		t.Skip("Skipping without --h2load flag.")
	***REMOVED***
	if runtime.GOOS != "linux" ***REMOVED***
		t.Skip("skipping Docker test when not on Linux; requires --net which won't work with boot2docker anyway")
	***REMOVED***
	requireH2load(t)

	msg := strings.Repeat("Hello, h2load!\n", 5000)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, msg)
		w.(http.Flusher).Flush()
		io.WriteString(w, msg)
	***REMOVED***))
	ts.StartTLS()
	defer ts.Close()

	cmd := exec.Command("docker", "run", "--net=host", "--entrypoint=/usr/local/bin/h2load", "gohttp2/curl",
		"-n100000", "-c100", "-m100", ts.URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Issue 12843
func TestServerDoS_MaxHeaderListSize(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED******REMOVED***)
	defer st.Close()

	// shake hands
	frameSize := defaultMaxReadFrameSize
	var advHeaderListSize *uint32
	st.greetAndCheckSettings(func(s Setting) error ***REMOVED***
		switch s.ID ***REMOVED***
		case SettingMaxFrameSize:
			if s.Val < minMaxFrameSize ***REMOVED***
				frameSize = minMaxFrameSize
			***REMOVED*** else if s.Val > maxFrameSize ***REMOVED***
				frameSize = maxFrameSize
			***REMOVED*** else ***REMOVED***
				frameSize = int(s.Val)
			***REMOVED***
		case SettingMaxHeaderListSize:
			advHeaderListSize = &s.Val
		***REMOVED***
		return nil
	***REMOVED***)

	if advHeaderListSize == nil ***REMOVED***
		t.Errorf("server didn't advertise a max header list size")
	***REMOVED*** else if *advHeaderListSize == 0 ***REMOVED***
		t.Errorf("server advertised a max header list size of 0")
	***REMOVED***

	st.encodeHeaderField(":method", "GET")
	st.encodeHeaderField(":path", "/")
	st.encodeHeaderField(":scheme", "https")
	cookie := strings.Repeat("*", 4058)
	st.encodeHeaderField("cookie", cookie)
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: st.headerBuf.Bytes(),
		EndStream:     true,
		EndHeaders:    false,
	***REMOVED***)

	// Capture the short encoding of a duplicate ~4K cookie, now
	// that we've already sent it once.
	st.headerBuf.Reset()
	st.encodeHeaderField("cookie", cookie)

	// Now send 1MB of it.
	const size = 1 << 20
	b := bytes.Repeat(st.headerBuf.Bytes(), size/st.headerBuf.Len())
	for len(b) > 0 ***REMOVED***
		chunk := b
		if len(chunk) > frameSize ***REMOVED***
			chunk = chunk[:frameSize]
		***REMOVED***
		b = b[len(chunk):]
		st.fr.WriteContinuation(1, len(b) == 0, chunk)
	***REMOVED***

	h := st.wantHeaders()
	if !h.HeadersEnded() ***REMOVED***
		t.Fatalf("Got HEADERS without END_HEADERS set: %v", h)
	***REMOVED***
	headers := st.decodeHeader(h.HeaderBlockFragment())
	want := [][2]string***REMOVED***
		***REMOVED***":status", "431"***REMOVED***,
		***REMOVED***"content-type", "text/html; charset=utf-8"***REMOVED***,
		***REMOVED***"content-length", "63"***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(headers, want) ***REMOVED***
		t.Errorf("Headers mismatch.\n got: %q\nwant: %q\n", headers, want)
	***REMOVED***
***REMOVED***

func TestCompressionErrorOnWrite(t *testing.T) ***REMOVED***
	const maxStrLen = 8 << 10
	var serverConfig *http.Server
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// No response body.
	***REMOVED***, func(ts *httptest.Server) ***REMOVED***
		serverConfig = ts.Config
		serverConfig.MaxHeaderBytes = maxStrLen
	***REMOVED***)
	st.addLogFilter("connection error: COMPRESSION_ERROR")
	defer st.Close()
	st.greet()

	maxAllowed := st.sc.framer.maxHeaderStringLen()

	// Crank this up, now that we have a conn connected with the
	// hpack.Decoder's max string length set has been initialized
	// from the earlier low ~8K value. We want this higher so don't
	// hit the max header list size. We only want to test hitting
	// the max string size.
	serverConfig.MaxHeaderBytes = 1 << 20

	// First a request with a header that's exactly the max allowed size
	// for the hpack compression. It's still too long for the header list
	// size, so we'll get the 431 error, but that keeps the compression
	// context still valid.
	hbf := st.encodeHeader("foo", strings.Repeat("a", maxAllowed))

	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: hbf,
		EndStream:     true,
		EndHeaders:    true,
	***REMOVED***)
	h := st.wantHeaders()
	if !h.HeadersEnded() ***REMOVED***
		t.Fatalf("Got HEADERS without END_HEADERS set: %v", h)
	***REMOVED***
	headers := st.decodeHeader(h.HeaderBlockFragment())
	want := [][2]string***REMOVED***
		***REMOVED***":status", "431"***REMOVED***,
		***REMOVED***"content-type", "text/html; charset=utf-8"***REMOVED***,
		***REMOVED***"content-length", "63"***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(headers, want) ***REMOVED***
		t.Errorf("Headers mismatch.\n got: %q\nwant: %q\n", headers, want)
	***REMOVED***
	df := st.wantData()
	if !strings.Contains(string(df.Data()), "HTTP Error 431") ***REMOVED***
		t.Errorf("Unexpected data body: %q", df.Data())
	***REMOVED***
	if !df.StreamEnded() ***REMOVED***
		t.Fatalf("expect data stream end")
	***REMOVED***

	// And now send one that's just one byte too big.
	hbf = st.encodeHeader("bar", strings.Repeat("b", maxAllowed+1))
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      3,
		BlockFragment: hbf,
		EndStream:     true,
		EndHeaders:    true,
	***REMOVED***)
	ga := st.wantGoAway()
	if ga.ErrCode != ErrCodeCompression ***REMOVED***
		t.Errorf("GOAWAY err = %v; want ErrCodeCompression", ga.ErrCode)
	***REMOVED***
***REMOVED***

func TestCompressionErrorOnClose(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// No response body.
	***REMOVED***)
	st.addLogFilter("connection error: COMPRESSION_ERROR")
	defer st.Close()
	st.greet()

	hbf := st.encodeHeader("foo", "bar")
	hbf = hbf[:len(hbf)-1] // truncate one byte from the end, so hpack.Decoder.Close fails.
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: hbf,
		EndStream:     true,
		EndHeaders:    true,
	***REMOVED***)
	ga := st.wantGoAway()
	if ga.ErrCode != ErrCodeCompression ***REMOVED***
		t.Errorf("GOAWAY err = %v; want ErrCodeCompression", ga.ErrCode)
	***REMOVED***
***REMOVED***

// test that a server handler can read trailers from a client
func TestServerReadsTrailers(t *testing.T) ***REMOVED***
	const testBody = "some test body"
	writeReq := func(st *serverTester) ***REMOVED***
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1, // clients send odd numbers
			BlockFragment: st.encodeHeader("trailer", "Foo, Bar", "trailer", "Baz"),
			EndStream:     false,
			EndHeaders:    true,
		***REMOVED***)
		st.writeData(1, false, []byte(testBody))
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID: 1, // clients send odd numbers
			BlockFragment: st.encodeHeaderRaw(
				"foo", "foov",
				"bar", "barv",
				"baz", "bazv",
				"surprise", "wasn't declared; shouldn't show up",
			),
			EndStream:  true,
			EndHeaders: true,
		***REMOVED***)
	***REMOVED***
	checkReq := func(r *http.Request) ***REMOVED***
		wantTrailer := http.Header***REMOVED***
			"Foo": nil,
			"Bar": nil,
			"Baz": nil,
		***REMOVED***
		if !reflect.DeepEqual(r.Trailer, wantTrailer) ***REMOVED***
			t.Errorf("initial Trailer = %v; want %v", r.Trailer, wantTrailer)
		***REMOVED***
		slurp, err := ioutil.ReadAll(r.Body)
		if string(slurp) != testBody ***REMOVED***
			t.Errorf("read body %q; want %q", slurp, testBody)
		***REMOVED***
		if err != nil ***REMOVED***
			t.Fatalf("Body slurp: %v", err)
		***REMOVED***
		wantTrailerAfter := http.Header***REMOVED***
			"Foo": ***REMOVED***"foov"***REMOVED***,
			"Bar": ***REMOVED***"barv"***REMOVED***,
			"Baz": ***REMOVED***"bazv"***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(r.Trailer, wantTrailerAfter) ***REMOVED***
			t.Errorf("final Trailer = %v; want %v", r.Trailer, wantTrailerAfter)
		***REMOVED***
	***REMOVED***
	testServerRequest(t, writeReq, checkReq)
***REMOVED***

// test that a server handler can send trailers
func TestServerWritesTrailers_WithFlush(t *testing.T)    ***REMOVED*** testServerWritesTrailers(t, true) ***REMOVED***
func TestServerWritesTrailers_WithoutFlush(t *testing.T) ***REMOVED*** testServerWritesTrailers(t, false) ***REMOVED***

func testServerWritesTrailers(t *testing.T, withFlush bool) ***REMOVED***
	// See https://httpwg.github.io/specs/rfc7540.html#rfc.section.8.1.3
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.Header().Set("Trailer", "Server-Trailer-A, Server-Trailer-B")
		w.Header().Add("Trailer", "Server-Trailer-C")
		w.Header().Add("Trailer", "Transfer-Encoding, Content-Length, Trailer") // filtered

		// Regular headers:
		w.Header().Set("Foo", "Bar")
		w.Header().Set("Content-Length", "5") // len("Hello")

		io.WriteString(w, "Hello")
		if withFlush ***REMOVED***
			w.(http.Flusher).Flush()
		***REMOVED***
		w.Header().Set("Server-Trailer-A", "valuea")
		w.Header().Set("Server-Trailer-C", "valuec") // skipping B
		// After a flush, random keys like Server-Surprise shouldn't show up:
		w.Header().Set("Server-Surpise", "surprise! this isn't predeclared!")
		// But we do permit promoting keys to trailers after a
		// flush if they start with the magic
		// otherwise-invalid "Trailer:" prefix:
		w.Header().Set("Trailer:Post-Header-Trailer", "hi1")
		w.Header().Set("Trailer:post-header-trailer2", "hi2")
		w.Header().Set("Trailer:Range", "invalid")
		w.Header().Set("Trailer:Foo\x01Bogus", "invalid")
		w.Header().Set("Transfer-Encoding", "should not be included; Forbidden by RFC 2616 14.40")
		w.Header().Set("Content-Length", "should not be included; Forbidden by RFC 2616 14.40")
		w.Header().Set("Trailer", "should not be included; Forbidden by RFC 2616 14.40")
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if hf.StreamEnded() ***REMOVED***
			t.Fatal("response HEADERS had END_STREAM")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("response HEADERS didn't have END_HEADERS")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"foo", "Bar"***REMOVED***,
			***REMOVED***"trailer", "Server-Trailer-A, Server-Trailer-B"***REMOVED***,
			***REMOVED***"trailer", "Server-Trailer-C"***REMOVED***,
			***REMOVED***"trailer", "Transfer-Encoding, Content-Length, Trailer"***REMOVED***,
			***REMOVED***"content-type", "text/plain; charset=utf-8"***REMOVED***,
			***REMOVED***"content-length", "5"***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Header mismatch.\n got: %v\nwant: %v", goth, wanth)
		***REMOVED***
		df := st.wantData()
		if string(df.Data()) != "Hello" ***REMOVED***
			t.Fatalf("Client read %q; want Hello", df.Data())
		***REMOVED***
		if df.StreamEnded() ***REMOVED***
			t.Fatalf("data frame had STREAM_ENDED")
		***REMOVED***
		tf := st.wantHeaders() // for the trailers
		if !tf.StreamEnded() ***REMOVED***
			t.Fatalf("trailers HEADERS lacked END_STREAM")
		***REMOVED***
		if !tf.HeadersEnded() ***REMOVED***
			t.Fatalf("trailers HEADERS lacked END_HEADERS")
		***REMOVED***
		wanth = [][2]string***REMOVED***
			***REMOVED***"post-header-trailer", "hi1"***REMOVED***,
			***REMOVED***"post-header-trailer2", "hi2"***REMOVED***,
			***REMOVED***"server-trailer-a", "valuea"***REMOVED***,
			***REMOVED***"server-trailer-c", "valuec"***REMOVED***,
		***REMOVED***
		goth = st.decodeHeader(tf.HeaderBlockFragment())
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Header mismatch.\n got: %v\nwant: %v", goth, wanth)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// validate transmitted header field names & values
// golang.org/issue/14048
func TestServerDoesntWriteInvalidHeaders(t *testing.T) ***REMOVED***
	testServerResponse(t, func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		w.Header().Add("OK1", "x")
		w.Header().Add("Bad:Colon", "x") // colon (non-token byte) in key
		w.Header().Add("Bad1\x00", "x")  // null in key
		w.Header().Add("Bad2", "x\x00y") // null in value
		return nil
	***REMOVED***, func(st *serverTester) ***REMOVED***
		getSlash(st)
		hf := st.wantHeaders()
		if !hf.StreamEnded() ***REMOVED***
			t.Error("response HEADERS lacked END_STREAM")
		***REMOVED***
		if !hf.HeadersEnded() ***REMOVED***
			t.Fatal("response HEADERS didn't have END_HEADERS")
		***REMOVED***
		goth := st.decodeHeader(hf.HeaderBlockFragment())
		wanth := [][2]string***REMOVED***
			***REMOVED***":status", "200"***REMOVED***,
			***REMOVED***"ok1", "x"***REMOVED***,
			***REMOVED***"content-length", "0"***REMOVED***,
		***REMOVED***
		if !reflect.DeepEqual(goth, wanth) ***REMOVED***
			t.Errorf("Header mismatch.\n got: %v\nwant: %v", goth, wanth)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkServerGets(b *testing.B) ***REMOVED***
	defer disableGoroutineTracking()()
	b.ReportAllocs()

	const msg = "Hello, world"
	st := newServerTester(b, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, msg)
	***REMOVED***)
	defer st.Close()
	st.greet()

	// Give the server quota to reply. (plus it has the the 64KB)
	if err := st.fr.WriteWindowUpdate(0, uint32(b.N*len(msg))); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	for i := 0; i < b.N; i++ ***REMOVED***
		id := 1 + uint32(i)*2
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      id,
			BlockFragment: st.encodeHeader(),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
		st.wantHeaders()
		df := st.wantData()
		if !df.StreamEnded() ***REMOVED***
			b.Fatalf("DATA didn't have END_STREAM; got %v", df)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkServerPosts(b *testing.B) ***REMOVED***
	defer disableGoroutineTracking()()
	b.ReportAllocs()

	const msg = "Hello, world"
	st := newServerTester(b, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// Consume the (empty) body from th peer before replying, otherwise
		// the server will sometimes (depending on scheduling) send the peer a
		// a RST_STREAM with the CANCEL error code.
		if n, err := io.Copy(ioutil.Discard, r.Body); n != 0 || err != nil ***REMOVED***
			b.Errorf("Copy error; got %v, %v; want 0, nil", n, err)
		***REMOVED***
		io.WriteString(w, msg)
	***REMOVED***)
	defer st.Close()
	st.greet()

	// Give the server quota to reply. (plus it has the the 64KB)
	if err := st.fr.WriteWindowUpdate(0, uint32(b.N*len(msg))); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	for i := 0; i < b.N; i++ ***REMOVED***
		id := 1 + uint32(i)*2
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      id,
			BlockFragment: st.encodeHeader(":method", "POST"),
			EndStream:     false,
			EndHeaders:    true,
		***REMOVED***)
		st.writeData(id, true, nil)
		st.wantHeaders()
		df := st.wantData()
		if !df.StreamEnded() ***REMOVED***
			b.Fatalf("DATA didn't have END_STREAM; got %v", df)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Send a stream of messages from server to client in separate data frames.
// Brings up performance issues seen in long streams.
// Created to show problem in go issue #18502
func BenchmarkServerToClientStreamDefaultOptions(b *testing.B) ***REMOVED***
	benchmarkServerToClientStream(b)
***REMOVED***

// Justification for Change-Id: Iad93420ef6c3918f54249d867098f1dadfa324d8
// Expect to see memory/alloc reduction by opting in to Frame reuse with the Framer.
func BenchmarkServerToClientStreamReuseFrames(b *testing.B) ***REMOVED***
	benchmarkServerToClientStream(b, optFramerReuseFrames)
***REMOVED***

func benchmarkServerToClientStream(b *testing.B, newServerOpts ...interface***REMOVED******REMOVED***) ***REMOVED***
	defer disableGoroutineTracking()()
	b.ReportAllocs()
	const msgLen = 1
	// default window size
	const windowSize = 1<<16 - 1

	// next message to send from the server and for the client to expect
	nextMsg := func(i int) []byte ***REMOVED***
		msg := make([]byte, msgLen)
		msg[0] = byte(i)
		if len(msg) != msgLen ***REMOVED***
			panic("invalid test setup msg length")
		***REMOVED***
		return msg
	***REMOVED***

	st := newServerTester(b, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// Consume the (empty) body from th peer before replying, otherwise
		// the server will sometimes (depending on scheduling) send the peer a
		// a RST_STREAM with the CANCEL error code.
		if n, err := io.Copy(ioutil.Discard, r.Body); n != 0 || err != nil ***REMOVED***
			b.Errorf("Copy error; got %v, %v; want 0, nil", n, err)
		***REMOVED***
		for i := 0; i < b.N; i += 1 ***REMOVED***
			w.Write(nextMsg(i))
			w.(http.Flusher).Flush()
		***REMOVED***
	***REMOVED***, newServerOpts...)
	defer st.Close()
	st.greet()

	const id = uint32(1)

	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      id,
		BlockFragment: st.encodeHeader(":method", "POST"),
		EndStream:     false,
		EndHeaders:    true,
	***REMOVED***)

	st.writeData(id, true, nil)
	st.wantHeaders()

	var pendingWindowUpdate = uint32(0)

	for i := 0; i < b.N; i += 1 ***REMOVED***
		expected := nextMsg(i)
		df := st.wantData()
		if bytes.Compare(expected, df.data) != 0 ***REMOVED***
			b.Fatalf("Bad message received; want %v; got %v", expected, df.data)
		***REMOVED***
		// try to send infrequent but large window updates so they don't overwhelm the test
		pendingWindowUpdate += uint32(len(df.data))
		if pendingWindowUpdate >= windowSize/2 ***REMOVED***
			if err := st.fr.WriteWindowUpdate(0, pendingWindowUpdate); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
			if err := st.fr.WriteWindowUpdate(id, pendingWindowUpdate); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
			pendingWindowUpdate = 0
		***REMOVED***
	***REMOVED***
	df := st.wantData()
	if !df.StreamEnded() ***REMOVED***
		b.Fatalf("DATA didn't have END_STREAM; got %v", df)
	***REMOVED***
***REMOVED***

// go-fuzz bug, originally reported at https://github.com/bradfitz/http2/issues/53
// Verify we don't hang.
func TestIssue53(t *testing.T) ***REMOVED***
	const data = "PRI * HTTP/2.0\r\n\r\nSM" +
		"\r\n\r\n\x00\x00\x00\x01\ainfinfin\ad"
	s := &http.Server***REMOVED***
		ErrorLog: log.New(io.MultiWriter(stderrv(), twriter***REMOVED***t: t***REMOVED***), "", log.LstdFlags),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
			w.Write([]byte("hello"))
		***REMOVED***),
	***REMOVED***
	s2 := &Server***REMOVED***
		MaxReadFrameSize:             1 << 16,
		PermitProhibitedCipherSuites: true,
	***REMOVED***
	c := &issue53Conn***REMOVED***[]byte(data), false, false***REMOVED***
	s2.ServeConn(c, &ServeConnOpts***REMOVED***BaseConfig: s***REMOVED***)
	if !c.closed ***REMOVED***
		t.Fatal("connection is not closed")
	***REMOVED***
***REMOVED***

type issue53Conn struct ***REMOVED***
	data    []byte
	closed  bool
	written bool
***REMOVED***

func (c *issue53Conn) Read(b []byte) (n int, err error) ***REMOVED***
	if len(c.data) == 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	n = copy(b, c.data)
	c.data = c.data[n:]
	return
***REMOVED***

func (c *issue53Conn) Write(b []byte) (n int, err error) ***REMOVED***
	c.written = true
	return len(b), nil
***REMOVED***

func (c *issue53Conn) Close() error ***REMOVED***
	c.closed = true
	return nil
***REMOVED***

func (c *issue53Conn) LocalAddr() net.Addr ***REMOVED***
	return &net.TCPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1), Port: 49706***REMOVED***
***REMOVED***
func (c *issue53Conn) RemoteAddr() net.Addr ***REMOVED***
	return &net.TCPAddr***REMOVED***IP: net.IPv4(127, 0, 0, 1), Port: 49706***REMOVED***
***REMOVED***
func (c *issue53Conn) SetDeadline(t time.Time) error      ***REMOVED*** return nil ***REMOVED***
func (c *issue53Conn) SetReadDeadline(t time.Time) error  ***REMOVED*** return nil ***REMOVED***
func (c *issue53Conn) SetWriteDeadline(t time.Time) error ***REMOVED*** return nil ***REMOVED***

// golang.org/issue/12895
func TestConfigureServer(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		name      string
		tlsConfig *tls.Config
		wantErr   string
	***REMOVED******REMOVED***
		***REMOVED***
			name: "empty server",
		***REMOVED***,
		***REMOVED***
			name: "just the required cipher suite",
			tlsConfig: &tls.Config***REMOVED***
				CipherSuites: []uint16***REMOVED***tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "just the alternative required cipher suite",
			tlsConfig: &tls.Config***REMOVED***
				CipherSuites: []uint16***REMOVED***tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "missing required cipher suite",
			tlsConfig: &tls.Config***REMOVED***
				CipherSuites: []uint16***REMOVED***tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384***REMOVED***,
			***REMOVED***,
			wantErr: "is missing an HTTP/2-required AES_128_GCM_SHA256 cipher.",
		***REMOVED***,
		***REMOVED***
			name: "required after bad",
			tlsConfig: &tls.Config***REMOVED***
				CipherSuites: []uint16***REMOVED***tls.TLS_RSA_WITH_RC4_128_SHA, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256***REMOVED***,
			***REMOVED***,
			wantErr: "contains an HTTP/2-approved cipher suite (0xc02f), but it comes after",
		***REMOVED***,
		***REMOVED***
			name: "bad after required",
			tlsConfig: &tls.Config***REMOVED***
				CipherSuites: []uint16***REMOVED***tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_RSA_WITH_RC4_128_SHA***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		srv := &http.Server***REMOVED***TLSConfig: tt.tlsConfig***REMOVED***
		err := ConfigureServer(srv, nil)
		if (err != nil) != (tt.wantErr != "") ***REMOVED***
			if tt.wantErr != "" ***REMOVED***
				t.Errorf("%s: success, but want error", tt.name)
			***REMOVED*** else ***REMOVED***
				t.Errorf("%s: unexpected error: %v", tt.name, err)
			***REMOVED***
		***REMOVED***
		if err != nil && tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) ***REMOVED***
			t.Errorf("%s: err = %v; want substring %q", tt.name, err, tt.wantErr)
		***REMOVED***
		if err == nil && !srv.TLSConfig.PreferServerCipherSuites ***REMOVED***
			t.Errorf("%s: PreferServerCipherSuite is false; want true", tt.name)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestServerRejectHeadWithBody(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// No response body.
	***REMOVED***)
	defer st.Close()
	st.greet()
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1, // clients send odd numbers
		BlockFragment: st.encodeHeader(":method", "HEAD"),
		EndStream:     false, // what we're testing, a bogus HEAD request with body
		EndHeaders:    true,
	***REMOVED***)
	st.wantRSTStream(1, ErrCodeProtocol)
***REMOVED***

func TestServerNoAutoContentLengthOnHead(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// No response body. (or smaller than one frame)
	***REMOVED***)
	defer st.Close()
	st.greet()
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1, // clients send odd numbers
		BlockFragment: st.encodeHeader(":method", "HEAD"),
		EndStream:     true,
		EndHeaders:    true,
	***REMOVED***)
	h := st.wantHeaders()
	headers := st.decodeHeader(h.HeaderBlockFragment())
	want := [][2]string***REMOVED***
		***REMOVED***":status", "200"***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(headers, want) ***REMOVED***
		t.Errorf("Headers mismatch.\n got: %q\nwant: %q\n", headers, want)
	***REMOVED***
***REMOVED***

// golang.org/issue/13495
func TestServerNoDuplicateContentType(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header()["Content-Type"] = []string***REMOVED***""***REMOVED***
		fmt.Fprintf(w, "<html><head></head><body>hi</body></html>")
	***REMOVED***)
	defer st.Close()
	st.greet()
	st.writeHeaders(HeadersFrameParam***REMOVED***
		StreamID:      1,
		BlockFragment: st.encodeHeader(),
		EndStream:     true,
		EndHeaders:    true,
	***REMOVED***)
	h := st.wantHeaders()
	headers := st.decodeHeader(h.HeaderBlockFragment())
	want := [][2]string***REMOVED***
		***REMOVED***":status", "200"***REMOVED***,
		***REMOVED***"content-type", ""***REMOVED***,
		***REMOVED***"content-length", "41"***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(headers, want) ***REMOVED***
		t.Errorf("Headers mismatch.\n got: %q\nwant: %q\n", headers, want)
	***REMOVED***
***REMOVED***

func disableGoroutineTracking() (restore func()) ***REMOVED***
	old := DebugGoroutines
	DebugGoroutines = false
	return func() ***REMOVED*** DebugGoroutines = old ***REMOVED***
***REMOVED***

func BenchmarkServer_GetRequest(b *testing.B) ***REMOVED***
	defer disableGoroutineTracking()()
	b.ReportAllocs()
	const msg = "Hello, world."
	st := newServerTester(b, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		n, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil || n > 0 ***REMOVED***
			b.Errorf("Read %d bytes, error %v; want 0 bytes.", n, err)
		***REMOVED***
		io.WriteString(w, msg)
	***REMOVED***)
	defer st.Close()

	st.greet()
	// Give the server quota to reply. (plus it has the the 64KB)
	if err := st.fr.WriteWindowUpdate(0, uint32(b.N*len(msg))); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	hbf := st.encodeHeader(":method", "GET")
	for i := 0; i < b.N; i++ ***REMOVED***
		streamID := uint32(1 + 2*i)
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      streamID,
			BlockFragment: hbf,
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
		st.wantHeaders()
		st.wantData()
	***REMOVED***
***REMOVED***

func BenchmarkServer_PostRequest(b *testing.B) ***REMOVED***
	defer disableGoroutineTracking()()
	b.ReportAllocs()
	const msg = "Hello, world."
	st := newServerTester(b, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		n, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil || n > 0 ***REMOVED***
			b.Errorf("Read %d bytes, error %v; want 0 bytes.", n, err)
		***REMOVED***
		io.WriteString(w, msg)
	***REMOVED***)
	defer st.Close()
	st.greet()
	// Give the server quota to reply. (plus it has the the 64KB)
	if err := st.fr.WriteWindowUpdate(0, uint32(b.N*len(msg))); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	hbf := st.encodeHeader(":method", "POST")
	for i := 0; i < b.N; i++ ***REMOVED***
		streamID := uint32(1 + 2*i)
		st.writeHeaders(HeadersFrameParam***REMOVED***
			StreamID:      streamID,
			BlockFragment: hbf,
			EndStream:     false,
			EndHeaders:    true,
		***REMOVED***)
		st.writeData(streamID, true, nil)
		st.wantHeaders()
		st.wantData()
	***REMOVED***
***REMOVED***

type connStateConn struct ***REMOVED***
	net.Conn
	cs tls.ConnectionState
***REMOVED***

func (c connStateConn) ConnectionState() tls.ConnectionState ***REMOVED*** return c.cs ***REMOVED***

// golang.org/issue/12737 -- handle any net.Conn, not just
// *tls.Conn.
func TestServerHandleCustomConn(t *testing.T) ***REMOVED***
	var s Server
	c1, c2 := net.Pipe()
	clientDone := make(chan struct***REMOVED******REMOVED***)
	handlerDone := make(chan struct***REMOVED******REMOVED***)
	var req *http.Request
	go func() ***REMOVED***
		defer close(clientDone)
		defer c2.Close()
		fr := NewFramer(c2, c2)
		io.WriteString(c2, ClientPreface)
		fr.WriteSettings()
		fr.WriteSettingsAck()
		f, err := fr.ReadFrame()
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if sf, ok := f.(*SettingsFrame); !ok || sf.IsAck() ***REMOVED***
			t.Errorf("Got %v; want non-ACK SettingsFrame", summarizeFrame(f))
			return
		***REMOVED***
		f, err = fr.ReadFrame()
		if err != nil ***REMOVED***
			t.Error(err)
			return
		***REMOVED***
		if sf, ok := f.(*SettingsFrame); !ok || !sf.IsAck() ***REMOVED***
			t.Errorf("Got %v; want ACK SettingsFrame", summarizeFrame(f))
			return
		***REMOVED***
		var henc hpackEncoder
		fr.WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      1,
			BlockFragment: henc.encodeHeaderRaw(t, ":method", "GET", ":path", "/", ":scheme", "https", ":authority", "foo.com"),
			EndStream:     true,
			EndHeaders:    true,
		***REMOVED***)
		go io.Copy(ioutil.Discard, c2)
		<-handlerDone
	***REMOVED***()
	const testString = "my custom ConnectionState"
	fakeConnState := tls.ConnectionState***REMOVED***
		ServerName:  testString,
		Version:     tls.VersionTLS12,
		CipherSuite: cipher_TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	***REMOVED***
	go s.ServeConn(connStateConn***REMOVED***c1, fakeConnState***REMOVED***, &ServeConnOpts***REMOVED***
		BaseConfig: &http.Server***REMOVED***
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
				defer close(handlerDone)
				req = r
			***REMOVED***),
		***REMOVED******REMOVED***)
	select ***REMOVED***
	case <-clientDone:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for handler")
	***REMOVED***
	if req.TLS == nil ***REMOVED***
		t.Fatalf("Request.TLS is nil. Got: %#v", req)
	***REMOVED***
	if req.TLS.ServerName != testString ***REMOVED***
		t.Fatalf("Request.TLS = %+v; want ServerName of %q", req.TLS, testString)
	***REMOVED***
***REMOVED***

// golang.org/issue/14214
func TestServer_Rejects_ConnHeaders(t *testing.T) ***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		t.Error("should not get to Handler")
	***REMOVED***)
	defer st.Close()
	st.greet()
	st.bodylessReq1("connection", "foo")
	hf := st.wantHeaders()
	goth := st.decodeHeader(hf.HeaderBlockFragment())
	wanth := [][2]string***REMOVED***
		***REMOVED***":status", "400"***REMOVED***,
		***REMOVED***"content-type", "text/plain; charset=utf-8"***REMOVED***,
		***REMOVED***"x-content-type-options", "nosniff"***REMOVED***,
		***REMOVED***"content-length", "51"***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(goth, wanth) ***REMOVED***
		t.Errorf("Got headers %v; want %v", goth, wanth)
	***REMOVED***
***REMOVED***

type hpackEncoder struct ***REMOVED***
	enc *hpack.Encoder
	buf bytes.Buffer
***REMOVED***

func (he *hpackEncoder) encodeHeaderRaw(t *testing.T, headers ...string) []byte ***REMOVED***
	if len(headers)%2 == 1 ***REMOVED***
		panic("odd number of kv args")
	***REMOVED***
	he.buf.Reset()
	if he.enc == nil ***REMOVED***
		he.enc = hpack.NewEncoder(&he.buf)
	***REMOVED***
	for len(headers) > 0 ***REMOVED***
		k, v := headers[0], headers[1]
		err := he.enc.WriteField(hpack.HeaderField***REMOVED***Name: k, Value: v***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatalf("HPACK encoding error for %q/%q: %v", k, v, err)
		***REMOVED***
		headers = headers[2:]
	***REMOVED***
	return he.buf.Bytes()
***REMOVED***

func TestCheckValidHTTP2Request(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		h    http.Header
		want error
	***REMOVED******REMOVED***
		***REMOVED***
			h:    http.Header***REMOVED***"Te": ***REMOVED***"trailers"***REMOVED******REMOVED***,
			want: nil,
		***REMOVED***,
		***REMOVED***
			h:    http.Header***REMOVED***"Te": ***REMOVED***"trailers", "bogus"***REMOVED******REMOVED***,
			want: errors.New(`request header "TE" may only be "trailers" in HTTP/2`),
		***REMOVED***,
		***REMOVED***
			h:    http.Header***REMOVED***"Foo": ***REMOVED***""***REMOVED******REMOVED***,
			want: nil,
		***REMOVED***,
		***REMOVED***
			h:    http.Header***REMOVED***"Connection": ***REMOVED***""***REMOVED******REMOVED***,
			want: errors.New(`request header "Connection" is not valid in HTTP/2`),
		***REMOVED***,
		***REMOVED***
			h:    http.Header***REMOVED***"Proxy-Connection": ***REMOVED***""***REMOVED******REMOVED***,
			want: errors.New(`request header "Proxy-Connection" is not valid in HTTP/2`),
		***REMOVED***,
		***REMOVED***
			h:    http.Header***REMOVED***"Keep-Alive": ***REMOVED***""***REMOVED******REMOVED***,
			want: errors.New(`request header "Keep-Alive" is not valid in HTTP/2`),
		***REMOVED***,
		***REMOVED***
			h:    http.Header***REMOVED***"Upgrade": ***REMOVED***""***REMOVED******REMOVED***,
			want: errors.New(`request header "Upgrade" is not valid in HTTP/2`),
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		got := checkValidHTTP2RequestHeaders(tt.h)
		if !reflect.DeepEqual(got, tt.want) ***REMOVED***
			t.Errorf("%d. checkValidHTTP2Request = %v; want %v", i, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// golang.org/issue/14030
func TestExpect100ContinueAfterHandlerWrites(t *testing.T) ***REMOVED***
	const msg = "Hello"
	const msg2 = "World"

	doRead := make(chan bool, 1)
	defer close(doRead) // fallback cleanup

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, msg)
		w.(http.Flusher).Flush()

		// Do a read, which might force a 100-continue status to be sent.
		<-doRead
		r.Body.Read(make([]byte, 10))

		io.WriteString(w, msg2)

	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	req, _ := http.NewRequest("POST", st.ts.URL, io.LimitReader(neverEnding('A'), 2<<20))
	req.Header.Set("Expect", "100-continue")

	res, err := tr.RoundTrip(req)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()

	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(res.Body, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(buf) != msg ***REMOVED***
		t.Fatalf("msg = %q; want %q", buf, msg)
	***REMOVED***

	doRead <- true

	if _, err := io.ReadFull(res.Body, buf); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(buf) != msg2 ***REMOVED***
		t.Fatalf("second msg = %q; want %q", buf, msg2)
	***REMOVED***
***REMOVED***

type funcReader func([]byte) (n int, err error)

func (f funcReader) Read(p []byte) (n int, err error) ***REMOVED*** return f(p) ***REMOVED***

// golang.org/issue/16481 -- return flow control when streams close with unread data.
// (The Server version of the bug. See also TestUnreadFlowControlReturned_Transport)
func TestUnreadFlowControlReturned_Server(t *testing.T) ***REMOVED***
	unblock := make(chan bool, 1)
	defer close(unblock)

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// Don't read the 16KB request body. Wait until the client's
		// done sending it and then return. This should cause the Server
		// to then return those 16KB of flow control to the client.
		<-unblock
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()

	// This previously hung on the 4th iteration.
	for i := 0; i < 6; i++ ***REMOVED***
		body := io.MultiReader(
			io.LimitReader(neverEnding('A'), 16<<10),
			funcReader(func([]byte) (n int, err error) ***REMOVED***
				unblock <- true
				return 0, io.EOF
			***REMOVED***),
		)
		req, _ := http.NewRequest("POST", st.ts.URL, body)
		res, err := tr.RoundTrip(req)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		res.Body.Close()
	***REMOVED***

***REMOVED***

func TestServerIdleTimeout(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping in short mode")
	***REMOVED***

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
	***REMOVED***, func(h2s *Server) ***REMOVED***
		h2s.IdleTimeout = 500 * time.Millisecond
	***REMOVED***)
	defer st.Close()

	st.greet()
	ga := st.wantGoAway()
	if ga.ErrCode != ErrCodeNo ***REMOVED***
		t.Errorf("GOAWAY error = %v; want ErrCodeNo", ga.ErrCode)
	***REMOVED***
***REMOVED***

func TestServerIdleTimeout_AfterRequest(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping in short mode")
	***REMOVED***
	const timeout = 250 * time.Millisecond

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		time.Sleep(timeout * 2)
	***REMOVED***, func(h2s *Server) ***REMOVED***
		h2s.IdleTimeout = timeout
	***REMOVED***)
	defer st.Close()

	st.greet()

	// Send a request which takes twice the timeout. Verifies the
	// idle timeout doesn't fire while we're in a request:
	st.bodylessReq1()
	st.wantHeaders()

	// But the idle timeout should be rearmed after the request
	// is done:
	ga := st.wantGoAway()
	if ga.ErrCode != ErrCodeNo ***REMOVED***
		t.Errorf("GOAWAY error = %v; want ErrCodeNo", ga.ErrCode)
	***REMOVED***
***REMOVED***

// grpc-go closes the Request.Body currently with a Read.
// Verify that it doesn't race.
// See https://github.com/grpc/grpc-go/pull/938
func TestRequestBodyReadCloseRace(t *testing.T) ***REMOVED***
	for i := 0; i < 100; i++ ***REMOVED***
		body := &requestBody***REMOVED***
			pipe: &pipe***REMOVED***
				b: new(bytes.Buffer),
			***REMOVED***,
		***REMOVED***
		body.pipe.CloseWithError(io.EOF)

		done := make(chan bool, 1)
		buf := make([]byte, 10)
		go func() ***REMOVED***
			time.Sleep(1 * time.Millisecond)
			body.Close()
			done <- true
		***REMOVED***()
		body.Read(buf)
		<-done
	***REMOVED***
***REMOVED***

func TestIssue20704Race(t *testing.T) ***REMOVED***
	if testing.Short() && os.Getenv("GO_BUILDER_NAME") == "" ***REMOVED***
		t.Skip("skipping in short mode")
	***REMOVED***
	const (
		itemSize  = 1 << 10
		itemCount = 100
	)

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		for i := 0; i < itemCount; i++ ***REMOVED***
			_, err := w.Write(make([]byte, itemSize))
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***, optOnlyServer)
	defer st.Close()

	tr := &Transport***REMOVED***TLSClientConfig: tlsConfigInsecure***REMOVED***
	defer tr.CloseIdleConnections()
	cl := &http.Client***REMOVED***Transport: tr***REMOVED***

	for i := 0; i < 1000; i++ ***REMOVED***
		resp, err := cl.Get(st.ts.URL)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		// Force a RST stream to the server by closing without
		// reading the body:
		resp.Body.Close()
	***REMOVED***
***REMOVED***
