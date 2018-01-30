// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/net/http2/hpack"
)

var knownFailing = flag.Bool("known_failing", false, "Run known-failing tests.")

func condSkipFailingTest(t *testing.T) ***REMOVED***
	if !*knownFailing ***REMOVED***
		t.Skip("Skipping known-failing test without --known_failing")
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	inTests = true
	DebugGoroutines = true
	flag.BoolVar(&VerboseLogs, "verboseh2", VerboseLogs, "Verbose HTTP/2 debug logging")
***REMOVED***

func TestSettingString(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		s    Setting
		want string
	***REMOVED******REMOVED***
		***REMOVED***Setting***REMOVED***SettingMaxFrameSize, 123***REMOVED***, "[MAX_FRAME_SIZE = 123]"***REMOVED***,
		***REMOVED***Setting***REMOVED***1<<16 - 1, 123***REMOVED***, "[UNKNOWN_SETTING_65535 = 123]"***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		got := fmt.Sprint(tt.s)
		if got != tt.want ***REMOVED***
			t.Errorf("%d. for %#v, string = %q; want %q", i, tt.s, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

type twriter struct ***REMOVED***
	t  testing.TB
	st *serverTester // optional
***REMOVED***

func (w twriter) Write(p []byte) (n int, err error) ***REMOVED***
	if w.st != nil ***REMOVED***
		ps := string(p)
		for _, phrase := range w.st.logFilter ***REMOVED***
			if strings.Contains(ps, phrase) ***REMOVED***
				return len(p), nil // no logging
			***REMOVED***
		***REMOVED***
	***REMOVED***
	w.t.Logf("%s", p)
	return len(p), nil
***REMOVED***

// like encodeHeader, but don't add implicit pseudo headers.
func encodeHeaderNoImplicit(t *testing.T, headers ...string) []byte ***REMOVED***
	var buf bytes.Buffer
	enc := hpack.NewEncoder(&buf)
	for len(headers) > 0 ***REMOVED***
		k, v := headers[0], headers[1]
		headers = headers[2:]
		if err := enc.WriteField(hpack.HeaderField***REMOVED***Name: k, Value: v***REMOVED***); err != nil ***REMOVED***
			t.Fatalf("HPACK encoding error for %q/%q: %v", k, v, err)
		***REMOVED***
	***REMOVED***
	return buf.Bytes()
***REMOVED***

// Verify that curl has http2.
func requireCurl(t *testing.T) ***REMOVED***
	out, err := dockerLogs(curl(t, "--version"))
	if err != nil ***REMOVED***
		t.Skipf("failed to determine curl features; skipping test")
	***REMOVED***
	if !strings.Contains(string(out), "HTTP2") ***REMOVED***
		t.Skip("curl doesn't support HTTP2; skipping test")
	***REMOVED***
***REMOVED***

func curl(t *testing.T, args ...string) (container string) ***REMOVED***
	out, err := exec.Command("docker", append([]string***REMOVED***"run", "-d", "--net=host", "gohttp2/curl"***REMOVED***, args...)...).Output()
	if err != nil ***REMOVED***
		t.Skipf("Failed to run curl in docker: %v, %s", err, out)
	***REMOVED***
	return strings.TrimSpace(string(out))
***REMOVED***

// Verify that h2load exists.
func requireH2load(t *testing.T) ***REMOVED***
	out, err := dockerLogs(h2load(t, "--version"))
	if err != nil ***REMOVED***
		t.Skipf("failed to probe h2load; skipping test: %s", out)
	***REMOVED***
	if !strings.Contains(string(out), "h2load nghttp2/") ***REMOVED***
		t.Skipf("h2load not present; skipping test. (Output=%q)", out)
	***REMOVED***
***REMOVED***

func h2load(t *testing.T, args ...string) (container string) ***REMOVED***
	out, err := exec.Command("docker", append([]string***REMOVED***"run", "-d", "--net=host", "--entrypoint=/usr/local/bin/h2load", "gohttp2/curl"***REMOVED***, args...)...).Output()
	if err != nil ***REMOVED***
		t.Skipf("Failed to run h2load in docker: %v, %s", err, out)
	***REMOVED***
	return strings.TrimSpace(string(out))
***REMOVED***

type puppetCommand struct ***REMOVED***
	fn   func(w http.ResponseWriter, r *http.Request)
	done chan<- bool
***REMOVED***

type handlerPuppet struct ***REMOVED***
	ch chan puppetCommand
***REMOVED***

func newHandlerPuppet() *handlerPuppet ***REMOVED***
	return &handlerPuppet***REMOVED***
		ch: make(chan puppetCommand),
	***REMOVED***
***REMOVED***

func (p *handlerPuppet) act(w http.ResponseWriter, r *http.Request) ***REMOVED***
	for cmd := range p.ch ***REMOVED***
		cmd.fn(w, r)
		cmd.done <- true
	***REMOVED***
***REMOVED***

func (p *handlerPuppet) done() ***REMOVED*** close(p.ch) ***REMOVED***
func (p *handlerPuppet) do(fn func(http.ResponseWriter, *http.Request)) ***REMOVED***
	done := make(chan bool)
	p.ch <- puppetCommand***REMOVED***fn, done***REMOVED***
	<-done
***REMOVED***
func dockerLogs(container string) ([]byte, error) ***REMOVED***
	out, err := exec.Command("docker", "wait", container).CombinedOutput()
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	exitStatus, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil ***REMOVED***
		return out, errors.New("unexpected exit status from docker wait")
	***REMOVED***
	out, err = exec.Command("docker", "logs", container).CombinedOutput()
	exec.Command("docker", "rm", container).Run()
	if err == nil && exitStatus != 0 ***REMOVED***
		err = fmt.Errorf("exit status %d: %s", exitStatus, out)
	***REMOVED***
	return out, err
***REMOVED***

func kill(container string) ***REMOVED***
	exec.Command("docker", "kill", container).Run()
	exec.Command("docker", "rm", container).Run()
***REMOVED***

func cleanDate(res *http.Response) ***REMOVED***
	if d := res.Header["Date"]; len(d) == 1 ***REMOVED***
		d[0] = "XXX"
	***REMOVED***
***REMOVED***

func TestSorterPoolAllocs(t *testing.T) ***REMOVED***
	ss := []string***REMOVED***"a", "b", "c"***REMOVED***
	h := http.Header***REMOVED***
		"a": nil,
		"b": nil,
		"c": nil,
	***REMOVED***
	sorter := new(sorter)

	if allocs := testing.AllocsPerRun(100, func() ***REMOVED***
		sorter.SortStrings(ss)
	***REMOVED***); allocs >= 1 ***REMOVED***
		t.Logf("SortStrings allocs = %v; want <1", allocs)
	***REMOVED***

	if allocs := testing.AllocsPerRun(5, func() ***REMOVED***
		if len(sorter.Keys(h)) != 3 ***REMOVED***
			t.Fatal("wrong result")
		***REMOVED***
	***REMOVED***); allocs > 0 ***REMOVED***
		t.Logf("Keys allocs = %v; want <1", allocs)
	***REMOVED***
***REMOVED***
