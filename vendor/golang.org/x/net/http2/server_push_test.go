// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package http2

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestServer_Push_Success(t *testing.T) ***REMOVED***
	const (
		mainBody   = "<html>index page</html>"
		pushedBody = "<html>pushed page</html>"
		userAgent  = "testagent"
		cookie     = "testcookie"
	)

	var stURL string
	checkPromisedReq := func(r *http.Request, wantMethod string, wantH http.Header) error ***REMOVED***
		if got, want := r.Method, wantMethod; got != want ***REMOVED***
			return fmt.Errorf("promised Req.Method=%q, want %q", got, want)
		***REMOVED***
		if got, want := r.Header, wantH; !reflect.DeepEqual(got, want) ***REMOVED***
			return fmt.Errorf("promised Req.Header=%q, want %q", got, want)
		***REMOVED***
		if got, want := "https://"+r.Host, stURL; got != want ***REMOVED***
			return fmt.Errorf("promised Req.Host=%q, want %q", got, want)
		***REMOVED***
		if r.Body == nil ***REMOVED***
			return fmt.Errorf("nil Body")
		***REMOVED***
		if buf, err := ioutil.ReadAll(r.Body); err != nil || len(buf) != 0 ***REMOVED***
			return fmt.Errorf("ReadAll(Body)=%q,%v, want '',nil", buf, err)
		***REMOVED***
		return nil
	***REMOVED***

	errc := make(chan error, 3)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		switch r.URL.RequestURI() ***REMOVED***
		case "/":
			// Push "/pushed?get" as a GET request, using an absolute URL.
			opt := &http.PushOptions***REMOVED***
				Header: http.Header***REMOVED***
					"User-Agent": ***REMOVED***userAgent***REMOVED***,
				***REMOVED***,
			***REMOVED***
			if err := w.(http.Pusher).Push(stURL+"/pushed?get", opt); err != nil ***REMOVED***
				errc <- fmt.Errorf("error pushing /pushed?get: %v", err)
				return
			***REMOVED***
			// Push "/pushed?head" as a HEAD request, using a path.
			opt = &http.PushOptions***REMOVED***
				Method: "HEAD",
				Header: http.Header***REMOVED***
					"User-Agent": ***REMOVED***userAgent***REMOVED***,
					"Cookie":     ***REMOVED***cookie***REMOVED***,
				***REMOVED***,
			***REMOVED***
			if err := w.(http.Pusher).Push("/pushed?head", opt); err != nil ***REMOVED***
				errc <- fmt.Errorf("error pushing /pushed?head: %v", err)
				return
			***REMOVED***
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(mainBody)))
			w.WriteHeader(200)
			io.WriteString(w, mainBody)
			errc <- nil

		case "/pushed?get":
			wantH := http.Header***REMOVED******REMOVED***
			wantH.Set("User-Agent", userAgent)
			if err := checkPromisedReq(r, "GET", wantH); err != nil ***REMOVED***
				errc <- fmt.Errorf("/pushed?get: %v", err)
				return
			***REMOVED***
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(pushedBody)))
			w.WriteHeader(200)
			io.WriteString(w, pushedBody)
			errc <- nil

		case "/pushed?head":
			wantH := http.Header***REMOVED******REMOVED***
			wantH.Set("User-Agent", userAgent)
			wantH.Set("Cookie", cookie)
			if err := checkPromisedReq(r, "HEAD", wantH); err != nil ***REMOVED***
				errc <- fmt.Errorf("/pushed?head: %v", err)
				return
			***REMOVED***
			w.WriteHeader(204)
			errc <- nil

		default:
			errc <- fmt.Errorf("unknown RequestURL %q", r.URL.RequestURI())
		***REMOVED***
	***REMOVED***)
	stURL = st.ts.URL

	// Send one request, which should push two responses.
	st.greet()
	getSlash(st)
	for k := 0; k < 3; k++ ***REMOVED***
		select ***REMOVED***
		case <-time.After(2 * time.Second):
			t.Errorf("timeout waiting for handler %d to finish", k)
		case err := <-errc:
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	checkPushPromise := func(f Frame, promiseID uint32, wantH [][2]string) error ***REMOVED***
		pp, ok := f.(*PushPromiseFrame)
		if !ok ***REMOVED***
			return fmt.Errorf("got a %T; want *PushPromiseFrame", f)
		***REMOVED***
		if !pp.HeadersEnded() ***REMOVED***
			return fmt.Errorf("want END_HEADERS flag in PushPromiseFrame")
		***REMOVED***
		if got, want := pp.PromiseID, promiseID; got != want ***REMOVED***
			return fmt.Errorf("got PromiseID %v; want %v", got, want)
		***REMOVED***
		gotH := st.decodeHeader(pp.HeaderBlockFragment())
		if !reflect.DeepEqual(gotH, wantH) ***REMOVED***
			return fmt.Errorf("got promised headers %v; want %v", gotH, wantH)
		***REMOVED***
		return nil
	***REMOVED***
	checkHeaders := func(f Frame, wantH [][2]string) error ***REMOVED***
		hf, ok := f.(*HeadersFrame)
		if !ok ***REMOVED***
			return fmt.Errorf("got a %T; want *HeadersFrame", f)
		***REMOVED***
		gotH := st.decodeHeader(hf.HeaderBlockFragment())
		if !reflect.DeepEqual(gotH, wantH) ***REMOVED***
			return fmt.Errorf("got response headers %v; want %v", gotH, wantH)
		***REMOVED***
		return nil
	***REMOVED***
	checkData := func(f Frame, wantData string) error ***REMOVED***
		df, ok := f.(*DataFrame)
		if !ok ***REMOVED***
			return fmt.Errorf("got a %T; want *DataFrame", f)
		***REMOVED***
		if gotData := string(df.Data()); gotData != wantData ***REMOVED***
			return fmt.Errorf("got response data %q; want %q", gotData, wantData)
		***REMOVED***
		return nil
	***REMOVED***

	// Stream 1 has 2 PUSH_PROMISE + HEADERS + DATA
	// Stream 2 has HEADERS + DATA
	// Stream 4 has HEADERS
	expected := map[uint32][]func(Frame) error***REMOVED***
		1: ***REMOVED***
			func(f Frame) error ***REMOVED***
				return checkPushPromise(f, 2, [][2]string***REMOVED***
					***REMOVED***":method", "GET"***REMOVED***,
					***REMOVED***":scheme", "https"***REMOVED***,
					***REMOVED***":authority", st.ts.Listener.Addr().String()***REMOVED***,
					***REMOVED***":path", "/pushed?get"***REMOVED***,
					***REMOVED***"user-agent", userAgent***REMOVED***,
				***REMOVED***)
			***REMOVED***,
			func(f Frame) error ***REMOVED***
				return checkPushPromise(f, 4, [][2]string***REMOVED***
					***REMOVED***":method", "HEAD"***REMOVED***,
					***REMOVED***":scheme", "https"***REMOVED***,
					***REMOVED***":authority", st.ts.Listener.Addr().String()***REMOVED***,
					***REMOVED***":path", "/pushed?head"***REMOVED***,
					***REMOVED***"cookie", cookie***REMOVED***,
					***REMOVED***"user-agent", userAgent***REMOVED***,
				***REMOVED***)
			***REMOVED***,
			func(f Frame) error ***REMOVED***
				return checkHeaders(f, [][2]string***REMOVED***
					***REMOVED***":status", "200"***REMOVED***,
					***REMOVED***"content-type", "text/html"***REMOVED***,
					***REMOVED***"content-length", strconv.Itoa(len(mainBody))***REMOVED***,
				***REMOVED***)
			***REMOVED***,
			func(f Frame) error ***REMOVED***
				return checkData(f, mainBody)
			***REMOVED***,
		***REMOVED***,
		2: ***REMOVED***
			func(f Frame) error ***REMOVED***
				return checkHeaders(f, [][2]string***REMOVED***
					***REMOVED***":status", "200"***REMOVED***,
					***REMOVED***"content-type", "text/html"***REMOVED***,
					***REMOVED***"content-length", strconv.Itoa(len(pushedBody))***REMOVED***,
				***REMOVED***)
			***REMOVED***,
			func(f Frame) error ***REMOVED***
				return checkData(f, pushedBody)
			***REMOVED***,
		***REMOVED***,
		4: ***REMOVED***
			func(f Frame) error ***REMOVED***
				return checkHeaders(f, [][2]string***REMOVED***
					***REMOVED***":status", "204"***REMOVED***,
				***REMOVED***)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	consumed := map[uint32]int***REMOVED******REMOVED***
	for k := 0; len(expected) > 0; k++ ***REMOVED***
		f, err := st.readFrame()
		if err != nil ***REMOVED***
			for id, left := range expected ***REMOVED***
				t.Errorf("stream %d: missing %d frames", id, len(left))
			***REMOVED***
			t.Fatalf("readFrame %d: %v", k, err)
		***REMOVED***
		id := f.Header().StreamID
		label := fmt.Sprintf("stream %d, frame %d", id, consumed[id])
		if len(expected[id]) == 0 ***REMOVED***
			t.Fatalf("%s: unexpected frame %#+v", label, f)
		***REMOVED***
		check := expected[id][0]
		expected[id] = expected[id][1:]
		if len(expected[id]) == 0 ***REMOVED***
			delete(expected, id)
		***REMOVED***
		if err := check(f); err != nil ***REMOVED***
			t.Fatalf("%s: %v", label, err)
		***REMOVED***
		consumed[id]++
	***REMOVED***
***REMOVED***

func TestServer_Push_SuccessNoRace(t *testing.T) ***REMOVED***
	// Regression test for issue #18326. Ensure the request handler can mutate
	// pushed request headers without racing with the PUSH_PROMISE write.
	errc := make(chan error, 2)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		switch r.URL.RequestURI() ***REMOVED***
		case "/":
			opt := &http.PushOptions***REMOVED***
				Header: http.Header***REMOVED***"User-Agent": ***REMOVED***"testagent"***REMOVED******REMOVED***,
			***REMOVED***
			if err := w.(http.Pusher).Push("/pushed", opt); err != nil ***REMOVED***
				errc <- fmt.Errorf("error pushing: %v", err)
				return
			***REMOVED***
			w.WriteHeader(200)
			errc <- nil

		case "/pushed":
			// Update request header, ensure there is no race.
			r.Header.Set("User-Agent", "newagent")
			r.Header.Set("Cookie", "cookie")
			w.WriteHeader(200)
			errc <- nil

		default:
			errc <- fmt.Errorf("unknown RequestURL %q", r.URL.RequestURI())
		***REMOVED***
	***REMOVED***)

	// Send one request, which should push one response.
	st.greet()
	getSlash(st)
	for k := 0; k < 2; k++ ***REMOVED***
		select ***REMOVED***
		case <-time.After(2 * time.Second):
			t.Errorf("timeout waiting for handler %d to finish", k)
		case err := <-errc:
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestServer_Push_RejectRecursivePush(t *testing.T) ***REMOVED***
	// Expect two requests, but might get three if there's a bug and the second push succeeds.
	errc := make(chan error, 3)
	handler := func(w http.ResponseWriter, r *http.Request) error ***REMOVED***
		baseURL := "https://" + r.Host
		switch r.URL.Path ***REMOVED***
		case "/":
			if err := w.(http.Pusher).Push(baseURL+"/push1", nil); err != nil ***REMOVED***
				return fmt.Errorf("first Push()=%v, want nil", err)
			***REMOVED***
			return nil

		case "/push1":
			if got, want := w.(http.Pusher).Push(baseURL+"/push2", nil), ErrRecursivePush; got != want ***REMOVED***
				return fmt.Errorf("Push()=%v, want %v", got, want)
			***REMOVED***
			return nil

		default:
			return fmt.Errorf("unexpected path: %q", r.URL.Path)
		***REMOVED***
	***REMOVED***
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		errc <- handler(w, r)
	***REMOVED***)
	defer st.Close()
	st.greet()
	getSlash(st)
	if err := <-errc; err != nil ***REMOVED***
		t.Errorf("First request failed: %v", err)
	***REMOVED***
	if err := <-errc; err != nil ***REMOVED***
		t.Errorf("Second request failed: %v", err)
	***REMOVED***
***REMOVED***

func testServer_Push_RejectSingleRequest(t *testing.T, doPush func(http.Pusher, *http.Request) error, settings ...Setting) ***REMOVED***
	// Expect one request, but might get two if there's a bug and the push succeeds.
	errc := make(chan error, 2)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		errc <- doPush(w.(http.Pusher), r)
	***REMOVED***)
	defer st.Close()
	st.greet()
	if err := st.fr.WriteSettings(settings...); err != nil ***REMOVED***
		st.t.Fatalf("WriteSettings: %v", err)
	***REMOVED***
	st.wantSettingsAck()
	getSlash(st)
	if err := <-errc; err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	// Should not get a PUSH_PROMISE frame.
	hf := st.wantHeaders()
	if !hf.StreamEnded() ***REMOVED***
		t.Error("stream should end after headers")
	***REMOVED***
***REMOVED***

func TestServer_Push_RejectIfDisabled(t *testing.T) ***REMOVED***
	testServer_Push_RejectSingleRequest(t,
		func(p http.Pusher, r *http.Request) error ***REMOVED***
			if got, want := p.Push("https://"+r.Host+"/pushed", nil), http.ErrNotSupported; got != want ***REMOVED***
				return fmt.Errorf("Push()=%v, want %v", got, want)
			***REMOVED***
			return nil
		***REMOVED***,
		Setting***REMOVED***SettingEnablePush, 0***REMOVED***)
***REMOVED***

func TestServer_Push_RejectWhenNoConcurrentStreams(t *testing.T) ***REMOVED***
	testServer_Push_RejectSingleRequest(t,
		func(p http.Pusher, r *http.Request) error ***REMOVED***
			if got, want := p.Push("https://"+r.Host+"/pushed", nil), ErrPushLimitReached; got != want ***REMOVED***
				return fmt.Errorf("Push()=%v, want %v", got, want)
			***REMOVED***
			return nil
		***REMOVED***,
		Setting***REMOVED***SettingMaxConcurrentStreams, 0***REMOVED***)
***REMOVED***

func TestServer_Push_RejectWrongScheme(t *testing.T) ***REMOVED***
	testServer_Push_RejectSingleRequest(t,
		func(p http.Pusher, r *http.Request) error ***REMOVED***
			if err := p.Push("http://"+r.Host+"/pushed", nil); err == nil ***REMOVED***
				return errors.New("Push() should have failed (push target URL is http)")
			***REMOVED***
			return nil
		***REMOVED***)
***REMOVED***

func TestServer_Push_RejectMissingHost(t *testing.T) ***REMOVED***
	testServer_Push_RejectSingleRequest(t,
		func(p http.Pusher, r *http.Request) error ***REMOVED***
			if err := p.Push("https:pushed", nil); err == nil ***REMOVED***
				return errors.New("Push() should have failed (push target URL missing host)")
			***REMOVED***
			return nil
		***REMOVED***)
***REMOVED***

func TestServer_Push_RejectRelativePath(t *testing.T) ***REMOVED***
	testServer_Push_RejectSingleRequest(t,
		func(p http.Pusher, r *http.Request) error ***REMOVED***
			if err := p.Push("../test", nil); err == nil ***REMOVED***
				return errors.New("Push() should have failed (push target is a relative path)")
			***REMOVED***
			return nil
		***REMOVED***)
***REMOVED***

func TestServer_Push_RejectForbiddenMethod(t *testing.T) ***REMOVED***
	testServer_Push_RejectSingleRequest(t,
		func(p http.Pusher, r *http.Request) error ***REMOVED***
			if err := p.Push("https://"+r.Host+"/pushed", &http.PushOptions***REMOVED***Method: "POST"***REMOVED***); err == nil ***REMOVED***
				return errors.New("Push() should have failed (cannot promise a POST)")
			***REMOVED***
			return nil
		***REMOVED***)
***REMOVED***

func TestServer_Push_RejectForbiddenHeader(t *testing.T) ***REMOVED***
	testServer_Push_RejectSingleRequest(t,
		func(p http.Pusher, r *http.Request) error ***REMOVED***
			header := http.Header***REMOVED***
				"Content-Length":   ***REMOVED***"10"***REMOVED***,
				"Content-Encoding": ***REMOVED***"gzip"***REMOVED***,
				"Trailer":          ***REMOVED***"Foo"***REMOVED***,
				"Te":               ***REMOVED***"trailers"***REMOVED***,
				"Host":             ***REMOVED***"test.com"***REMOVED***,
				":authority":       ***REMOVED***"test.com"***REMOVED***,
			***REMOVED***
			if err := p.Push("https://"+r.Host+"/pushed", &http.PushOptions***REMOVED***Header: header***REMOVED***); err == nil ***REMOVED***
				return errors.New("Push() should have failed (forbidden headers)")
			***REMOVED***
			return nil
		***REMOVED***)
***REMOVED***

func TestServer_Push_StateTransitions(t *testing.T) ***REMOVED***
	const body = "foo"

	gotPromise := make(chan bool)
	finishedPush := make(chan bool)

	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		switch r.URL.RequestURI() ***REMOVED***
		case "/":
			if err := w.(http.Pusher).Push("/pushed", nil); err != nil ***REMOVED***
				t.Errorf("Push error: %v", err)
			***REMOVED***
			// Don't finish this request until the push finishes so we don't
			// nondeterministically interleave output frames with the push.
			<-finishedPush
		case "/pushed":
			<-gotPromise
		***REMOVED***
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(200)
		io.WriteString(w, body)
	***REMOVED***)
	defer st.Close()

	st.greet()
	if st.stream(2) != nil ***REMOVED***
		t.Fatal("stream 2 should be empty")
	***REMOVED***
	if got, want := st.streamState(2), stateIdle; got != want ***REMOVED***
		t.Fatalf("streamState(2)=%v, want %v", got, want)
	***REMOVED***
	getSlash(st)
	// After the PUSH_PROMISE is sent, the stream should be stateHalfClosedRemote.
	st.wantPushPromise()
	if got, want := st.streamState(2), stateHalfClosedRemote; got != want ***REMOVED***
		t.Fatalf("streamState(2)=%v, want %v", got, want)
	***REMOVED***
	// We stall the HTTP handler for "/pushed" until the above check. If we don't
	// stall the handler, then the handler might write HEADERS and DATA and finish
	// the stream before we check st.streamState(2) -- should that happen, we'll
	// see stateClosed and fail the above check.
	close(gotPromise)
	st.wantHeaders()
	if df := st.wantData(); !df.StreamEnded() ***REMOVED***
		t.Fatal("expected END_STREAM flag on DATA")
	***REMOVED***
	if got, want := st.streamState(2), stateClosed; got != want ***REMOVED***
		t.Fatalf("streamState(2)=%v, want %v", got, want)
	***REMOVED***
	close(finishedPush)
***REMOVED***

func TestServer_Push_RejectAfterGoAway(t *testing.T) ***REMOVED***
	var readyOnce sync.Once
	ready := make(chan struct***REMOVED******REMOVED***)
	errc := make(chan error, 2)
	st := newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		select ***REMOVED***
		case <-ready:
		case <-time.After(5 * time.Second):
			errc <- fmt.Errorf("timeout waiting for GOAWAY to be processed")
		***REMOVED***
		if got, want := w.(http.Pusher).Push("https://"+r.Host+"/pushed", nil), http.ErrNotSupported; got != want ***REMOVED***
			errc <- fmt.Errorf("Push()=%v, want %v", got, want)
		***REMOVED***
		errc <- nil
	***REMOVED***)
	defer st.Close()
	st.greet()
	getSlash(st)

	// Send GOAWAY and wait for it to be processed.
	st.fr.WriteGoAway(1, ErrCodeNo, nil)
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-ready:
				return
			default:
			***REMOVED***
			st.sc.serveMsgCh <- func(loopNum int) ***REMOVED***
				if !st.sc.pushEnabled ***REMOVED***
					readyOnce.Do(func() ***REMOVED*** close(ready) ***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	if err := <-errc; err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***
