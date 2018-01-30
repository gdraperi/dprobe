// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build h2demo

package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"hash/crc32"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"go4.org/syncutil/singleflight"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

var (
	prod = flag.Bool("prod", false, "Whether to configure itself to be the production http2.golang.org server.")

	httpsAddr = flag.String("https_addr", "localhost:4430", "TLS address to listen on ('host:port' or ':port'). Required.")
	httpAddr  = flag.String("http_addr", "", "Plain HTTP address to listen on ('host:port', or ':port'). Empty means no HTTP.")

	hostHTTP  = flag.String("http_host", "", "Optional host or host:port to use for http:// links to this service. By default, this is implied from -http_addr.")
	hostHTTPS = flag.String("https_host", "", "Optional host or host:port to use for http:// links to this service. By default, this is implied from -https_addr.")
)

func homeOldHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	io.WriteString(w, `<html>
<body>
<h1>Go + HTTP/2</h1>
<p>Welcome to <a href="https://golang.org/">the Go language</a>'s <a href="https://http2.github.io/">HTTP/2</a> demo & interop server.</p>
<p>Unfortunately, you're <b>not</b> using HTTP/2 right now. To do so:</p>
<ul>
   <li>Use Firefox Nightly or go to <b>about:config</b> and enable "network.http.spdy.enabled.http2draft"</li>
   <li>Use Google Chrome Canary and/or go to <b>chrome://flags/#enable-spdy4</b> to <i>Enable SPDY/4</i> (Chrome's name for HTTP/2)</li>
</ul>
<p>See code & instructions for connecting at <a href="https://github.com/golang/net/tree/master/http2">https://github.com/golang/net/tree/master/http2</a>.</p>

</body></html>`)
***REMOVED***

func home(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.URL.Path != "/" ***REMOVED***
		http.NotFound(w, r)
		return
	***REMOVED***
	io.WriteString(w, `<html>
<body>
<h1>Go + HTTP/2</h1>

<p>Welcome to <a href="https://golang.org/">the Go language</a>'s <a
href="https://http2.github.io/">HTTP/2</a> demo & interop server.</p>

<p>Congratulations, <b>you're using HTTP/2 right now</b>.</p>

<p>This server exists for others in the HTTP/2 community to test their HTTP/2 client implementations and point out flaws in our server.</p>

<p>
The code is at <a href="https://golang.org/x/net/http2">golang.org/x/net/http2</a> and
is used transparently by the Go standard library from Go 1.6 and later.
</p>

<p>Contact info: <i>bradfitz@golang.org</i>, or <a
href="https://golang.org/s/http2bug">file a bug</a>.</p>

<h2>Handlers for testing</h2>
<ul>
  <li>GET <a href="/reqinfo">/reqinfo</a> to dump the request + headers received</li>
  <li>GET <a href="/clockstream">/clockstream</a> streams the current time every second</li>
  <li>GET <a href="/gophertiles">/gophertiles</a> to see a page with a bunch of images</li>
  <li>GET <a href="/serverpush">/serverpush</a> to see a page with server push</li>
  <li>GET <a href="/file/gopher.png">/file/gopher.png</a> for a small file (does If-Modified-Since, Content-Range, etc)</li>
  <li>GET <a href="/file/go.src.tar.gz">/file/go.src.tar.gz</a> for a larger file (~10 MB)</li>
  <li>GET <a href="/redirect">/redirect</a> to redirect back to / (this page)</li>
  <li>GET <a href="/goroutines">/goroutines</a> to see all active goroutines in this server</li>
  <li>PUT something to <a href="/crc32">/crc32</a> to get a count of number of bytes and its CRC-32</li>
  <li>PUT something to <a href="/ECHO">/ECHO</a> and it will be streamed back to you capitalized</li>
</ul>

</body></html>`)
***REMOVED***

func reqInfoHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "RequestURI: %q\n", r.RequestURI)
	fmt.Fprintf(w, "URL: %#v\n", r.URL)
	fmt.Fprintf(w, "Body.ContentLength: %d (-1 means unknown)\n", r.ContentLength)
	fmt.Fprintf(w, "Close: %v (relevant for HTTP/1 only)\n", r.Close)
	fmt.Fprintf(w, "TLS: %#v\n", r.TLS)
	fmt.Fprintf(w, "\nHeaders:\n")
	r.Header.Write(w)
***REMOVED***

func crcHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.Method != "PUT" ***REMOVED***
		http.Error(w, "PUT required.", 400)
		return
	***REMOVED***
	crc := crc32.NewIEEE()
	n, err := io.Copy(crc, r.Body)
	if err == nil ***REMOVED***
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "bytes=%d, CRC32=%x", n, crc.Sum(nil))
	***REMOVED***
***REMOVED***

type capitalizeReader struct ***REMOVED***
	r io.Reader
***REMOVED***

func (cr capitalizeReader) Read(p []byte) (n int, err error) ***REMOVED***
	n, err = cr.r.Read(p)
	for i, b := range p[:n] ***REMOVED***
		if b >= 'a' && b <= 'z' ***REMOVED***
			p[i] = b - ('a' - 'A')
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

type flushWriter struct ***REMOVED***
	w io.Writer
***REMOVED***

func (fw flushWriter) Write(p []byte) (n int, err error) ***REMOVED***
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok ***REMOVED***
		f.Flush()
	***REMOVED***
	return
***REMOVED***

func echoCapitalHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.Method != "PUT" ***REMOVED***
		http.Error(w, "PUT required.", 400)
		return
	***REMOVED***
	io.Copy(flushWriter***REMOVED***w***REMOVED***, capitalizeReader***REMOVED***r.Body***REMOVED***)
***REMOVED***

var (
	fsGrp   singleflight.Group
	fsMu    sync.Mutex // guards fsCache
	fsCache = map[string]http.Handler***REMOVED******REMOVED***
)

// fileServer returns a file-serving handler that proxies URL.
// It lazily fetches URL on the first access and caches its contents forever.
func fileServer(url string, latency time.Duration) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if latency > 0 ***REMOVED***
			time.Sleep(latency)
		***REMOVED***
		hi, err := fsGrp.Do(url, func() (interface***REMOVED******REMOVED***, error) ***REMOVED***
			fsMu.Lock()
			if h, ok := fsCache[url]; ok ***REMOVED***
				fsMu.Unlock()
				return h, nil
			***REMOVED***
			fsMu.Unlock()

			res, err := http.Get(url)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			defer res.Body.Close()
			slurp, err := ioutil.ReadAll(res.Body)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			modTime := time.Now()
			var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
				http.ServeContent(w, r, path.Base(url), modTime, bytes.NewReader(slurp))
			***REMOVED***)
			fsMu.Lock()
			fsCache[url] = h
			fsMu.Unlock()
			return h, nil
		***REMOVED***)
		if err != nil ***REMOVED***
			http.Error(w, err.Error(), 500)
			return
		***REMOVED***
		hi.(http.Handler).ServeHTTP(w, r)
	***REMOVED***)
***REMOVED***

func clockStreamHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	clientGone := w.(http.CloseNotifier).CloseNotify()
	w.Header().Set("Content-Type", "text/plain")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Fprintf(w, "# ~1KB of junk to force browsers to start rendering immediately: \n")
	io.WriteString(w, strings.Repeat("# xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n", 13))

	for ***REMOVED***
		fmt.Fprintf(w, "%v\n", time.Now())
		w.(http.Flusher).Flush()
		select ***REMOVED***
		case <-ticker.C:
		case <-clientGone:
			log.Printf("Client %v disconnected from the clock", r.RemoteAddr)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func registerHandlers() ***REMOVED***
	tiles := newGopherTilesHandler()
	push := newPushHandler()

	mux2 := http.NewServeMux()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		switch ***REMOVED***
		case r.URL.Path == "/gophertiles":
			tiles.ServeHTTP(w, r) // allow HTTP/2 + HTTP/1.x
			return
		case strings.HasPrefix(r.URL.Path, "/serverpush"):
			push.ServeHTTP(w, r) // allow HTTP/2 + HTTP/1.x
			return
		case r.TLS == nil: // do not allow HTTP/1.x for anything else
			http.Redirect(w, r, "https://"+httpsHost()+"/", http.StatusFound)
			return
		***REMOVED***
		if r.ProtoMajor == 1 ***REMOVED***
			if r.URL.Path == "/reqinfo" ***REMOVED***
				reqInfoHandler(w, r)
				return
			***REMOVED***
			homeOldHTTP(w, r)
			return
		***REMOVED***
		mux2.ServeHTTP(w, r)
	***REMOVED***)
	mux2.HandleFunc("/", home)
	mux2.Handle("/file/gopher.png", fileServer("https://golang.org/doc/gopher/frontpage.png", 0))
	mux2.Handle("/file/go.src.tar.gz", fileServer("https://storage.googleapis.com/golang/go1.4.1.src.tar.gz", 0))
	mux2.HandleFunc("/reqinfo", reqInfoHandler)
	mux2.HandleFunc("/crc32", crcHandler)
	mux2.HandleFunc("/ECHO", echoCapitalHandler)
	mux2.HandleFunc("/clockstream", clockStreamHandler)
	mux2.Handle("/gophertiles", tiles)
	mux2.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		http.Redirect(w, r, "/", http.StatusFound)
	***REMOVED***)
	stripHomedir := regexp.MustCompile(`/(Users|home)/\w+`)
	mux2.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		buf := make([]byte, 2<<20)
		w.Write(stripHomedir.ReplaceAll(buf[:runtime.Stack(buf, true)], nil))
	***REMOVED***)
***REMOVED***

var pushResources = map[string]http.Handler***REMOVED***
	"/serverpush/static/jquery.min.js": fileServer("https://ajax.googleapis.com/ajax/libs/jquery/1.8.2/jquery.min.js", 100*time.Millisecond),
	"/serverpush/static/godocs.js":     fileServer("https://golang.org/lib/godoc/godocs.js", 100*time.Millisecond),
	"/serverpush/static/playground.js": fileServer("https://golang.org/lib/godoc/playground.js", 100*time.Millisecond),
	"/serverpush/static/style.css":     fileServer("https://golang.org/lib/godoc/style.css", 100*time.Millisecond),
***REMOVED***

func newPushHandler() http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		for path, handler := range pushResources ***REMOVED***
			if r.URL.Path == path ***REMOVED***
				handler.ServeHTTP(w, r)
				return
			***REMOVED***
		***REMOVED***

		cacheBust := time.Now().UnixNano()
		if pusher, ok := w.(http.Pusher); ok ***REMOVED***
			for path := range pushResources ***REMOVED***
				url := fmt.Sprintf("%s?%d", path, cacheBust)
				if err := pusher.Push(url, nil); err != nil ***REMOVED***
					log.Printf("Failed to push %v: %v", path, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		time.Sleep(100 * time.Millisecond) // fake network latency + parsing time
		if err := pushTmpl.Execute(w, struct ***REMOVED***
			CacheBust int64
			HTTPSHost string
			HTTPHost  string
		***REMOVED******REMOVED***
			CacheBust: cacheBust,
			HTTPSHost: httpsHost(),
			HTTPHost:  httpHost(),
		***REMOVED***); err != nil ***REMOVED***
			log.Printf("Executing server push template: %v", err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func newGopherTilesHandler() http.Handler ***REMOVED***
	const gopherURL = "https://blog.golang.org/go-programming-language-turns-two_gophers.jpg"
	res, err := http.Get(gopherURL)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	if res.StatusCode != 200 ***REMOVED***
		log.Fatalf("Error fetching %s: %v", gopherURL, res.Status)
	***REMOVED***
	slurp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	im, err := jpeg.Decode(bytes.NewReader(slurp))
	if err != nil ***REMOVED***
		if len(slurp) > 1024 ***REMOVED***
			slurp = slurp[:1024]
		***REMOVED***
		log.Fatalf("Failed to decode gopher image: %v (got %q)", err, slurp)
	***REMOVED***

	type subImager interface ***REMOVED***
		SubImage(image.Rectangle) image.Image
	***REMOVED***
	const tileSize = 32
	xt := im.Bounds().Max.X / tileSize
	yt := im.Bounds().Max.Y / tileSize
	var tile [][][]byte // y -> x -> jpeg bytes
	for yi := 0; yi < yt; yi++ ***REMOVED***
		var row [][]byte
		for xi := 0; xi < xt; xi++ ***REMOVED***
			si := im.(subImager).SubImage(image.Rectangle***REMOVED***
				Min: image.Point***REMOVED***xi * tileSize, yi * tileSize***REMOVED***,
				Max: image.Point***REMOVED***(xi + 1) * tileSize, (yi + 1) * tileSize***REMOVED***,
			***REMOVED***)
			buf := new(bytes.Buffer)
			if err := jpeg.Encode(buf, si, &jpeg.Options***REMOVED***Quality: 90***REMOVED***); err != nil ***REMOVED***
				log.Fatal(err)
			***REMOVED***
			row = append(row, buf.Bytes())
		***REMOVED***
		tile = append(tile, row)
	***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		ms, _ := strconv.Atoi(r.FormValue("latency"))
		const nanosPerMilli = 1e6
		if r.FormValue("x") != "" ***REMOVED***
			x, _ := strconv.Atoi(r.FormValue("x"))
			y, _ := strconv.Atoi(r.FormValue("y"))
			if ms <= 1000 ***REMOVED***
				time.Sleep(time.Duration(ms) * nanosPerMilli)
			***REMOVED***
			if x >= 0 && x < xt && y >= 0 && y < yt ***REMOVED***
				http.ServeContent(w, r, "", time.Time***REMOVED******REMOVED***, bytes.NewReader(tile[y][x]))
				return
			***REMOVED***
		***REMOVED***
		io.WriteString(w, "<html><body onload='showtimes()'>")
		fmt.Fprintf(w, "A grid of %d tiled images is below. Compare:<p>", xt*yt)
		for _, ms := range []int***REMOVED***0, 30, 200, 1000***REMOVED*** ***REMOVED***
			d := time.Duration(ms) * nanosPerMilli
			fmt.Fprintf(w, "[<a href='https://%s/gophertiles?latency=%d'>HTTP/2, %v latency</a>] [<a href='http://%s/gophertiles?latency=%d'>HTTP/1, %v latency</a>]<br>\n",
				httpsHost(), ms, d,
				httpHost(), ms, d,
			)
		***REMOVED***
		io.WriteString(w, "<p>\n")
		cacheBust := time.Now().UnixNano()
		for y := 0; y < yt; y++ ***REMOVED***
			for x := 0; x < xt; x++ ***REMOVED***
				fmt.Fprintf(w, "<img width=%d height=%d src='/gophertiles?x=%d&y=%d&cachebust=%d&latency=%d'>",
					tileSize, tileSize, x, y, cacheBust, ms)
			***REMOVED***
			io.WriteString(w, "<br/>\n")
		***REMOVED***
		io.WriteString(w, `<p><div id='loadtimes'></div></p>
<script>
function showtimes() ***REMOVED***
	var times = 'Times from connection start:<br>'
	times += 'DOM loaded: ' + (window.performance.timing.domContentLoadedEventEnd - window.performance.timing.connectStart) + 'ms<br>'
	times += 'DOM complete (images loaded): ' + (window.performance.timing.domComplete - window.performance.timing.connectStart) + 'ms<br>'
	document.getElementById('loadtimes').innerHTML = times
***REMOVED***
</script>
<hr><a href='/'>&lt;&lt Back to Go HTTP/2 demo server</a></body></html>`)
	***REMOVED***)
***REMOVED***

func httpsHost() string ***REMOVED***
	if *hostHTTPS != "" ***REMOVED***
		return *hostHTTPS
	***REMOVED***
	if v := *httpsAddr; strings.HasPrefix(v, ":") ***REMOVED***
		return "localhost" + v
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

func httpHost() string ***REMOVED***
	if *hostHTTP != "" ***REMOVED***
		return *hostHTTP
	***REMOVED***
	if v := *httpAddr; strings.HasPrefix(v, ":") ***REMOVED***
		return "localhost" + v
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

func serveProdTLS() error ***REMOVED***
	const cacheDir = "/var/cache/autocert"
	if err := os.MkdirAll(cacheDir, 0700); err != nil ***REMOVED***
		return err
	***REMOVED***
	m := autocert.Manager***REMOVED***
		Cache:      autocert.DirCache(cacheDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("http2.golang.org"),
	***REMOVED***
	srv := &http.Server***REMOVED***
		TLSConfig: &tls.Config***REMOVED***
			GetCertificate: m.GetCertificate,
		***REMOVED***,
	***REMOVED***
	http2.ConfigureServer(srv, &http2.Server***REMOVED***
		NewWriteScheduler: func() http2.WriteScheduler ***REMOVED***
			return http2.NewPriorityWriteScheduler(nil)
		***REMOVED***,
	***REMOVED***)
	ln, err := net.Listen("tcp", ":443")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return srv.Serve(tls.NewListener(tcpKeepAliveListener***REMOVED***ln.(*net.TCPListener)***REMOVED***, srv.TLSConfig))
***REMOVED***

type tcpKeepAliveListener struct ***REMOVED***
	*net.TCPListener
***REMOVED***

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) ***REMOVED***
	tc, err := ln.AcceptTCP()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
***REMOVED***

func serveProd() error ***REMOVED***
	errc := make(chan error, 2)
	go func() ***REMOVED*** errc <- http.ListenAndServe(":80", nil) ***REMOVED***()
	go func() ***REMOVED*** errc <- serveProdTLS() ***REMOVED***()
	return <-errc
***REMOVED***

const idleTimeout = 5 * time.Minute
const activeTimeout = 10 * time.Minute

// TODO: put this into the standard library and actually send
// PING frames and GOAWAY, etc: golang.org/issue/14204
func idleTimeoutHook() func(net.Conn, http.ConnState) ***REMOVED***
	var mu sync.Mutex
	m := map[net.Conn]*time.Timer***REMOVED******REMOVED***
	return func(c net.Conn, cs http.ConnState) ***REMOVED***
		mu.Lock()
		defer mu.Unlock()
		if t, ok := m[c]; ok ***REMOVED***
			delete(m, c)
			t.Stop()
		***REMOVED***
		var d time.Duration
		switch cs ***REMOVED***
		case http.StateNew, http.StateIdle:
			d = idleTimeout
		case http.StateActive:
			d = activeTimeout
		default:
			return
		***REMOVED***
		m[c] = time.AfterFunc(d, func() ***REMOVED***
			log.Printf("closing idle conn %v after %v", c.RemoteAddr(), d)
			go c.Close()
		***REMOVED***)
	***REMOVED***
***REMOVED***

func main() ***REMOVED***
	var srv http.Server
	flag.BoolVar(&http2.VerboseLogs, "verbose", false, "Verbose HTTP/2 debugging.")
	flag.Parse()
	srv.Addr = *httpsAddr
	srv.ConnState = idleTimeoutHook()

	registerHandlers()

	if *prod ***REMOVED***
		*hostHTTP = "http2.golang.org"
		*hostHTTPS = "http2.golang.org"
		log.Fatal(serveProd())
	***REMOVED***

	url := "https://" + httpsHost() + "/"
	log.Printf("Listening on " + url)
	http2.ConfigureServer(&srv, &http2.Server***REMOVED******REMOVED***)

	if *httpAddr != "" ***REMOVED***
		go func() ***REMOVED***
			log.Printf("Listening on http://" + httpHost() + "/ (for unencrypted HTTP/1)")
			log.Fatal(http.ListenAndServe(*httpAddr, nil))
		***REMOVED***()
	***REMOVED***

	go func() ***REMOVED***
		log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
	***REMOVED***()
	select ***REMOVED******REMOVED***
***REMOVED***
