// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!solaris

/*
The h2i command is an interactive HTTP/2 console.

Usage:
  $ h2i [flags] <hostname>

Interactive commands in the console: (all parts case-insensitive)

  ping [data]
  settings ack
  settings FOO=n BAR=z
  headers      (open a new stream by typing HTTP/1.1)
*/
package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

// Flags
var (
	flagNextProto = flag.String("nextproto", "h2,h2-14", "Comma-separated list of NPN/ALPN protocol names to negotiate.")
	flagInsecure  = flag.Bool("insecure", false, "Whether to skip TLS cert validation")
	flagSettings  = flag.String("settings", "empty", "comma-separated list of KEY=value settings for the initial SETTINGS frame. The magic value 'empty' sends an empty initial settings frame, and the magic value 'omit' causes no initial settings frame to be sent.")
	flagDial      = flag.String("dial", "", "optional ip:port to dial, to connect to a host:port but use a different SNI name (including a SNI name without DNS)")
)

type command struct ***REMOVED***
	run func(*h2i, []string) error // required

	// complete optionally specifies tokens (case-insensitive) which are
	// valid for this subcommand.
	complete func() []string
***REMOVED***

var commands = map[string]command***REMOVED***
	"ping": ***REMOVED***run: (*h2i).cmdPing***REMOVED***,
	"settings": ***REMOVED***
		run: (*h2i).cmdSettings,
		complete: func() []string ***REMOVED***
			return []string***REMOVED***
				"ACK",
				http2.SettingHeaderTableSize.String(),
				http2.SettingEnablePush.String(),
				http2.SettingMaxConcurrentStreams.String(),
				http2.SettingInitialWindowSize.String(),
				http2.SettingMaxFrameSize.String(),
				http2.SettingMaxHeaderListSize.String(),
			***REMOVED***
		***REMOVED***,
	***REMOVED***,
	"quit":    ***REMOVED***run: (*h2i).cmdQuit***REMOVED***,
	"headers": ***REMOVED***run: (*h2i).cmdHeaders***REMOVED***,
***REMOVED***

func usage() ***REMOVED***
	fmt.Fprintf(os.Stderr, "Usage: h2i <hostname>\n\n")
	flag.PrintDefaults()
***REMOVED***

// withPort adds ":443" if another port isn't already present.
func withPort(host string) string ***REMOVED***
	if _, _, err := net.SplitHostPort(host); err != nil ***REMOVED***
		return net.JoinHostPort(host, "443")
	***REMOVED***
	return host
***REMOVED***

// withoutPort strips the port from addr if present.
func withoutPort(addr string) string ***REMOVED***
	if h, _, err := net.SplitHostPort(addr); err == nil ***REMOVED***
		return h
	***REMOVED***
	return addr
***REMOVED***

// h2i is the app's state.
type h2i struct ***REMOVED***
	host   string
	tc     *tls.Conn
	framer *http2.Framer
	term   *terminal.Terminal

	// owned by the command loop:
	streamID uint32
	hbuf     bytes.Buffer
	henc     *hpack.Encoder

	// owned by the readFrames loop:
	peerSetting map[http2.SettingID]uint32
	hdec        *hpack.Decoder
***REMOVED***

func main() ***REMOVED***
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 ***REMOVED***
		usage()
		os.Exit(2)
	***REMOVED***
	log.SetFlags(0)

	host := flag.Arg(0)
	app := &h2i***REMOVED***
		host:        host,
		peerSetting: make(map[http2.SettingID]uint32),
	***REMOVED***
	app.henc = hpack.NewEncoder(&app.hbuf)

	if err := app.Main(); err != nil ***REMOVED***
		if app.term != nil ***REMOVED***
			app.logf("%v\n", err)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(os.Stderr, "%v\n", err)
		***REMOVED***
		os.Exit(1)
	***REMOVED***
	fmt.Fprintf(os.Stdout, "\n")
***REMOVED***

func (app *h2i) Main() error ***REMOVED***
	cfg := &tls.Config***REMOVED***
		ServerName:         withoutPort(app.host),
		NextProtos:         strings.Split(*flagNextProto, ","),
		InsecureSkipVerify: *flagInsecure,
	***REMOVED***

	hostAndPort := *flagDial
	if hostAndPort == "" ***REMOVED***
		hostAndPort = withPort(app.host)
	***REMOVED***
	log.Printf("Connecting to %s ...", hostAndPort)
	tc, err := tls.Dial("tcp", hostAndPort, cfg)
	if err != nil ***REMOVED***
		return fmt.Errorf("Error dialing %s: %v", hostAndPort, err)
	***REMOVED***
	log.Printf("Connected to %v", tc.RemoteAddr())
	defer tc.Close()

	if err := tc.Handshake(); err != nil ***REMOVED***
		return fmt.Errorf("TLS handshake: %v", err)
	***REMOVED***
	if !*flagInsecure ***REMOVED***
		if err := tc.VerifyHostname(app.host); err != nil ***REMOVED***
			return fmt.Errorf("VerifyHostname: %v", err)
		***REMOVED***
	***REMOVED***
	state := tc.ConnectionState()
	log.Printf("Negotiated protocol %q", state.NegotiatedProtocol)
	if !state.NegotiatedProtocolIsMutual || state.NegotiatedProtocol == "" ***REMOVED***
		return fmt.Errorf("Could not negotiate protocol mutually")
	***REMOVED***

	if _, err := io.WriteString(tc, http2.ClientPreface); err != nil ***REMOVED***
		return err
	***REMOVED***

	app.framer = http2.NewFramer(tc, tc)

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer terminal.Restore(0, oldState)

	var screen = struct ***REMOVED***
		io.Reader
		io.Writer
	***REMOVED******REMOVED***os.Stdin, os.Stdout***REMOVED***

	app.term = terminal.NewTerminal(screen, "h2i> ")
	lastWord := regexp.MustCompile(`.+\W(\w+)$`)
	app.term.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) ***REMOVED***
		if key != '\t' ***REMOVED***
			return
		***REMOVED***
		if pos != len(line) ***REMOVED***
			// TODO: we're being lazy for now, only supporting tab completion at the end.
			return
		***REMOVED***
		// Auto-complete for the command itself.
		if !strings.Contains(line, " ") ***REMOVED***
			var name string
			name, _, ok = lookupCommand(line)
			if !ok ***REMOVED***
				return
			***REMOVED***
			return name, len(name), true
		***REMOVED***
		_, c, ok := lookupCommand(line[:strings.IndexByte(line, ' ')])
		if !ok || c.complete == nil ***REMOVED***
			return
		***REMOVED***
		if strings.HasSuffix(line, " ") ***REMOVED***
			app.logf("%s", strings.Join(c.complete(), " "))
			return line, pos, true
		***REMOVED***
		m := lastWord.FindStringSubmatch(line)
		if m == nil ***REMOVED***
			return line, len(line), true
		***REMOVED***
		soFar := m[1]
		var match []string
		for _, cand := range c.complete() ***REMOVED***
			if len(soFar) > len(cand) || !strings.EqualFold(cand[:len(soFar)], soFar) ***REMOVED***
				continue
			***REMOVED***
			match = append(match, cand)
		***REMOVED***
		if len(match) == 0 ***REMOVED***
			return
		***REMOVED***
		if len(match) > 1 ***REMOVED***
			// TODO: auto-complete any common prefix
			app.logf("%s", strings.Join(match, " "))
			return line, pos, true
		***REMOVED***
		newLine = line[:len(line)-len(soFar)] + match[0]
		return newLine, len(newLine), true

	***REMOVED***

	errc := make(chan error, 2)
	go func() ***REMOVED*** errc <- app.readFrames() ***REMOVED***()
	go func() ***REMOVED*** errc <- app.readConsole() ***REMOVED***()
	return <-errc
***REMOVED***

func (app *h2i) logf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprintf(app.term, format+"\r\n", args...)
***REMOVED***

func (app *h2i) readConsole() error ***REMOVED***
	if s := *flagSettings; s != "omit" ***REMOVED***
		var args []string
		if s != "empty" ***REMOVED***
			args = strings.Split(s, ",")
		***REMOVED***
		_, c, ok := lookupCommand("settings")
		if !ok ***REMOVED***
			panic("settings command not found")
		***REMOVED***
		c.run(app, args)
	***REMOVED***

	for ***REMOVED***
		line, err := app.term.ReadLine()
		if err == io.EOF ***REMOVED***
			return nil
		***REMOVED***
		if err != nil ***REMOVED***
			return fmt.Errorf("terminal.ReadLine: %v", err)
		***REMOVED***
		f := strings.Fields(line)
		if len(f) == 0 ***REMOVED***
			continue
		***REMOVED***
		cmd, args := f[0], f[1:]
		if _, c, ok := lookupCommand(cmd); ok ***REMOVED***
			err = c.run(app, args)
		***REMOVED*** else ***REMOVED***
			app.logf("Unknown command %q", line)
		***REMOVED***
		if err == errExitApp ***REMOVED***
			return nil
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***

func lookupCommand(prefix string) (name string, c command, ok bool) ***REMOVED***
	prefix = strings.ToLower(prefix)
	if c, ok = commands[prefix]; ok ***REMOVED***
		return prefix, c, ok
	***REMOVED***

	for full, candidate := range commands ***REMOVED***
		if strings.HasPrefix(full, prefix) ***REMOVED***
			if c.run != nil ***REMOVED***
				return "", command***REMOVED******REMOVED***, false // ambiguous
			***REMOVED***
			c = candidate
			name = full
		***REMOVED***
	***REMOVED***
	return name, c, c.run != nil
***REMOVED***

var errExitApp = errors.New("internal sentinel error value to quit the console reading loop")

func (a *h2i) cmdQuit(args []string) error ***REMOVED***
	if len(args) > 0 ***REMOVED***
		a.logf("the QUIT command takes no argument")
		return nil
	***REMOVED***
	return errExitApp
***REMOVED***

func (a *h2i) cmdSettings(args []string) error ***REMOVED***
	if len(args) == 1 && strings.EqualFold(args[0], "ACK") ***REMOVED***
		return a.framer.WriteSettingsAck()
	***REMOVED***
	var settings []http2.Setting
	for _, arg := range args ***REMOVED***
		if strings.EqualFold(arg, "ACK") ***REMOVED***
			a.logf("Error: ACK must be only argument with the SETTINGS command")
			return nil
		***REMOVED***
		eq := strings.Index(arg, "=")
		if eq == -1 ***REMOVED***
			a.logf("Error: invalid argument %q (expected SETTING_NAME=nnnn)", arg)
			return nil
		***REMOVED***
		sid, ok := settingByName(arg[:eq])
		if !ok ***REMOVED***
			a.logf("Error: unknown setting name %q", arg[:eq])
			return nil
		***REMOVED***
		val, err := strconv.ParseUint(arg[eq+1:], 10, 32)
		if err != nil ***REMOVED***
			a.logf("Error: invalid argument %q (expected SETTING_NAME=nnnn)", arg)
			return nil
		***REMOVED***
		settings = append(settings, http2.Setting***REMOVED***
			ID:  sid,
			Val: uint32(val),
		***REMOVED***)
	***REMOVED***
	a.logf("Sending: %v", settings)
	return a.framer.WriteSettings(settings...)
***REMOVED***

func settingByName(name string) (http2.SettingID, bool) ***REMOVED***
	for _, sid := range [...]http2.SettingID***REMOVED***
		http2.SettingHeaderTableSize,
		http2.SettingEnablePush,
		http2.SettingMaxConcurrentStreams,
		http2.SettingInitialWindowSize,
		http2.SettingMaxFrameSize,
		http2.SettingMaxHeaderListSize,
	***REMOVED*** ***REMOVED***
		if strings.EqualFold(sid.String(), name) ***REMOVED***
			return sid, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

func (app *h2i) cmdPing(args []string) error ***REMOVED***
	if len(args) > 1 ***REMOVED***
		app.logf("invalid PING usage: only accepts 0 or 1 args")
		return nil // nil means don't end the program
	***REMOVED***
	var data [8]byte
	if len(args) == 1 ***REMOVED***
		copy(data[:], args[0])
	***REMOVED*** else ***REMOVED***
		copy(data[:], "h2i_ping")
	***REMOVED***
	return app.framer.WritePing(false, data)
***REMOVED***

func (app *h2i) cmdHeaders(args []string) error ***REMOVED***
	if len(args) > 0 ***REMOVED***
		app.logf("Error: HEADERS doesn't yet take arguments.")
		// TODO: flags for restricting window size, to force CONTINUATION
		// frames.
		return nil
	***REMOVED***
	var h1req bytes.Buffer
	app.term.SetPrompt("(as HTTP/1.1)> ")
	defer app.term.SetPrompt("h2i> ")
	for ***REMOVED***
		line, err := app.term.ReadLine()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		h1req.WriteString(line)
		h1req.WriteString("\r\n")
		if line == "" ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	req, err := http.ReadRequest(bufio.NewReader(&h1req))
	if err != nil ***REMOVED***
		app.logf("Invalid HTTP/1.1 request: %v", err)
		return nil
	***REMOVED***
	if app.streamID == 0 ***REMOVED***
		app.streamID = 1
	***REMOVED*** else ***REMOVED***
		app.streamID += 2
	***REMOVED***
	app.logf("Opening Stream-ID %d:", app.streamID)
	hbf := app.encodeHeaders(req)
	if len(hbf) > 16<<10 ***REMOVED***
		app.logf("TODO: h2i doesn't yet write CONTINUATION frames. Copy it from transport.go")
		return nil
	***REMOVED***
	return app.framer.WriteHeaders(http2.HeadersFrameParam***REMOVED***
		StreamID:      app.streamID,
		BlockFragment: hbf,
		EndStream:     req.Method == "GET" || req.Method == "HEAD", // good enough for now
		EndHeaders:    true,                                        // for now
	***REMOVED***)
***REMOVED***

func (app *h2i) readFrames() error ***REMOVED***
	for ***REMOVED***
		f, err := app.framer.ReadFrame()
		if err != nil ***REMOVED***
			return fmt.Errorf("ReadFrame: %v", err)
		***REMOVED***
		app.logf("%v", f)
		switch f := f.(type) ***REMOVED***
		case *http2.PingFrame:
			app.logf("  Data = %q", f.Data)
		case *http2.SettingsFrame:
			f.ForeachSetting(func(s http2.Setting) error ***REMOVED***
				app.logf("  %v", s)
				app.peerSetting[s.ID] = s.Val
				return nil
			***REMOVED***)
		case *http2.WindowUpdateFrame:
			app.logf("  Window-Increment = %v", f.Increment)
		case *http2.GoAwayFrame:
			app.logf("  Last-Stream-ID = %d; Error-Code = %v (%d)", f.LastStreamID, f.ErrCode, f.ErrCode)
		case *http2.DataFrame:
			app.logf("  %q", f.Data())
		case *http2.HeadersFrame:
			if f.HasPriority() ***REMOVED***
				app.logf("  PRIORITY = %v", f.Priority)
			***REMOVED***
			if app.hdec == nil ***REMOVED***
				// TODO: if the user uses h2i to send a SETTINGS frame advertising
				// something larger, we'll need to respect SETTINGS_HEADER_TABLE_SIZE
				// and stuff here instead of using the 4k default. But for now:
				tableSize := uint32(4 << 10)
				app.hdec = hpack.NewDecoder(tableSize, app.onNewHeaderField)
			***REMOVED***
			app.hdec.Write(f.HeaderBlockFragment())
		case *http2.PushPromiseFrame:
			if app.hdec == nil ***REMOVED***
				// TODO: if the user uses h2i to send a SETTINGS frame advertising
				// something larger, we'll need to respect SETTINGS_HEADER_TABLE_SIZE
				// and stuff here instead of using the 4k default. But for now:
				tableSize := uint32(4 << 10)
				app.hdec = hpack.NewDecoder(tableSize, app.onNewHeaderField)
			***REMOVED***
			app.hdec.Write(f.HeaderBlockFragment())
		***REMOVED***
	***REMOVED***
***REMOVED***

// called from readLoop
func (app *h2i) onNewHeaderField(f hpack.HeaderField) ***REMOVED***
	if f.Sensitive ***REMOVED***
		app.logf("  %s = %q (SENSITIVE)", f.Name, f.Value)
	***REMOVED***
	app.logf("  %s = %q", f.Name, f.Value)
***REMOVED***

func (app *h2i) encodeHeaders(req *http.Request) []byte ***REMOVED***
	app.hbuf.Reset()

	// TODO(bradfitz): figure out :authority-vs-Host stuff between http2 and Go
	host := req.Host
	if host == "" ***REMOVED***
		host = req.URL.Host
	***REMOVED***

	path := req.RequestURI
	if path == "" ***REMOVED***
		path = "/"
	***REMOVED***

	app.writeHeader(":authority", host) // probably not right for all sites
	app.writeHeader(":method", req.Method)
	app.writeHeader(":path", path)
	app.writeHeader(":scheme", "https")

	for k, vv := range req.Header ***REMOVED***
		lowKey := strings.ToLower(k)
		if lowKey == "host" ***REMOVED***
			continue
		***REMOVED***
		for _, v := range vv ***REMOVED***
			app.writeHeader(lowKey, v)
		***REMOVED***
	***REMOVED***
	return app.hbuf.Bytes()
***REMOVED***

func (app *h2i) writeHeader(name, value string) ***REMOVED***
	app.henc.WriteField(hpack.HeaderField***REMOVED***Name: name, Value: value***REMOVED***)
	app.logf(" %s = %s", name, value)
***REMOVED***
