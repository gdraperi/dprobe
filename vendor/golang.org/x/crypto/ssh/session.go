// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

// Session implements an interactive session described in
// "RFC 4254, section 6".

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

type Signal string

// POSIX signals as listed in RFC 4254 Section 6.10.
const (
	SIGABRT Signal = "ABRT"
	SIGALRM Signal = "ALRM"
	SIGFPE  Signal = "FPE"
	SIGHUP  Signal = "HUP"
	SIGILL  Signal = "ILL"
	SIGINT  Signal = "INT"
	SIGKILL Signal = "KILL"
	SIGPIPE Signal = "PIPE"
	SIGQUIT Signal = "QUIT"
	SIGSEGV Signal = "SEGV"
	SIGTERM Signal = "TERM"
	SIGUSR1 Signal = "USR1"
	SIGUSR2 Signal = "USR2"
)

var signals = map[Signal]int***REMOVED***
	SIGABRT: 6,
	SIGALRM: 14,
	SIGFPE:  8,
	SIGHUP:  1,
	SIGILL:  4,
	SIGINT:  2,
	SIGKILL: 9,
	SIGPIPE: 13,
	SIGQUIT: 3,
	SIGSEGV: 11,
	SIGTERM: 15,
***REMOVED***

type TerminalModes map[uint8]uint32

// POSIX terminal mode flags as listed in RFC 4254 Section 8.
const (
	tty_OP_END    = 0
	VINTR         = 1
	VQUIT         = 2
	VERASE        = 3
	VKILL         = 4
	VEOF          = 5
	VEOL          = 6
	VEOL2         = 7
	VSTART        = 8
	VSTOP         = 9
	VSUSP         = 10
	VDSUSP        = 11
	VREPRINT      = 12
	VWERASE       = 13
	VLNEXT        = 14
	VFLUSH        = 15
	VSWTCH        = 16
	VSTATUS       = 17
	VDISCARD      = 18
	IGNPAR        = 30
	PARMRK        = 31
	INPCK         = 32
	ISTRIP        = 33
	INLCR         = 34
	IGNCR         = 35
	ICRNL         = 36
	IUCLC         = 37
	IXON          = 38
	IXANY         = 39
	IXOFF         = 40
	IMAXBEL       = 41
	ISIG          = 50
	ICANON        = 51
	XCASE         = 52
	ECHO          = 53
	ECHOE         = 54
	ECHOK         = 55
	ECHONL        = 56
	NOFLSH        = 57
	TOSTOP        = 58
	IEXTEN        = 59
	ECHOCTL       = 60
	ECHOKE        = 61
	PENDIN        = 62
	OPOST         = 70
	OLCUC         = 71
	ONLCR         = 72
	OCRNL         = 73
	ONOCR         = 74
	ONLRET        = 75
	CS7           = 90
	CS8           = 91
	PARENB        = 92
	PARODD        = 93
	TTY_OP_ISPEED = 128
	TTY_OP_OSPEED = 129
)

// A Session represents a connection to a remote command or shell.
type Session struct ***REMOVED***
	// Stdin specifies the remote process's standard input.
	// If Stdin is nil, the remote process reads from an empty
	// bytes.Buffer.
	Stdin io.Reader

	// Stdout and Stderr specify the remote process's standard
	// output and error.
	//
	// If either is nil, Run connects the corresponding file
	// descriptor to an instance of ioutil.Discard. There is a
	// fixed amount of buffering that is shared for the two streams.
	// If either blocks it may eventually cause the remote
	// command to block.
	Stdout io.Writer
	Stderr io.Writer

	ch        Channel // the channel backing this session
	started   bool    // true once Start, Run or Shell is invoked.
	copyFuncs []func() error
	errors    chan error // one send per copyFunc

	// true if pipe method is active
	stdinpipe, stdoutpipe, stderrpipe bool

	// stdinPipeWriter is non-nil if StdinPipe has not been called
	// and Stdin was specified by the user; it is the write end of
	// a pipe connecting Session.Stdin to the stdin channel.
	stdinPipeWriter io.WriteCloser

	exitStatus chan error
***REMOVED***

// SendRequest sends an out-of-band channel request on the SSH channel
// underlying the session.
func (s *Session) SendRequest(name string, wantReply bool, payload []byte) (bool, error) ***REMOVED***
	return s.ch.SendRequest(name, wantReply, payload)
***REMOVED***

func (s *Session) Close() error ***REMOVED***
	return s.ch.Close()
***REMOVED***

// RFC 4254 Section 6.4.
type setenvRequest struct ***REMOVED***
	Name  string
	Value string
***REMOVED***

// Setenv sets an environment variable that will be applied to any
// command executed by Shell or Run.
func (s *Session) Setenv(name, value string) error ***REMOVED***
	msg := setenvRequest***REMOVED***
		Name:  name,
		Value: value,
	***REMOVED***
	ok, err := s.ch.SendRequest("env", true, Marshal(&msg))
	if err == nil && !ok ***REMOVED***
		err = errors.New("ssh: setenv failed")
	***REMOVED***
	return err
***REMOVED***

// RFC 4254 Section 6.2.
type ptyRequestMsg struct ***REMOVED***
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
***REMOVED***

// RequestPty requests the association of a pty with the session on the remote host.
func (s *Session) RequestPty(term string, h, w int, termmodes TerminalModes) error ***REMOVED***
	var tm []byte
	for k, v := range termmodes ***REMOVED***
		kv := struct ***REMOVED***
			Key byte
			Val uint32
		***REMOVED******REMOVED***k, v***REMOVED***

		tm = append(tm, Marshal(&kv)...)
	***REMOVED***
	tm = append(tm, tty_OP_END)
	req := ptyRequestMsg***REMOVED***
		Term:     term,
		Columns:  uint32(w),
		Rows:     uint32(h),
		Width:    uint32(w * 8),
		Height:   uint32(h * 8),
		Modelist: string(tm),
	***REMOVED***
	ok, err := s.ch.SendRequest("pty-req", true, Marshal(&req))
	if err == nil && !ok ***REMOVED***
		err = errors.New("ssh: pty-req failed")
	***REMOVED***
	return err
***REMOVED***

// RFC 4254 Section 6.5.
type subsystemRequestMsg struct ***REMOVED***
	Subsystem string
***REMOVED***

// RequestSubsystem requests the association of a subsystem with the session on the remote host.
// A subsystem is a predefined command that runs in the background when the ssh session is initiated
func (s *Session) RequestSubsystem(subsystem string) error ***REMOVED***
	msg := subsystemRequestMsg***REMOVED***
		Subsystem: subsystem,
	***REMOVED***
	ok, err := s.ch.SendRequest("subsystem", true, Marshal(&msg))
	if err == nil && !ok ***REMOVED***
		err = errors.New("ssh: subsystem request failed")
	***REMOVED***
	return err
***REMOVED***

// RFC 4254 Section 6.7.
type ptyWindowChangeMsg struct ***REMOVED***
	Columns uint32
	Rows    uint32
	Width   uint32
	Height  uint32
***REMOVED***

// WindowChange informs the remote host about a terminal window dimension change to h rows and w columns.
func (s *Session) WindowChange(h, w int) error ***REMOVED***
	req := ptyWindowChangeMsg***REMOVED***
		Columns: uint32(w),
		Rows:    uint32(h),
		Width:   uint32(w * 8),
		Height:  uint32(h * 8),
	***REMOVED***
	_, err := s.ch.SendRequest("window-change", false, Marshal(&req))
	return err
***REMOVED***

// RFC 4254 Section 6.9.
type signalMsg struct ***REMOVED***
	Signal string
***REMOVED***

// Signal sends the given signal to the remote process.
// sig is one of the SIG* constants.
func (s *Session) Signal(sig Signal) error ***REMOVED***
	msg := signalMsg***REMOVED***
		Signal: string(sig),
	***REMOVED***

	_, err := s.ch.SendRequest("signal", false, Marshal(&msg))
	return err
***REMOVED***

// RFC 4254 Section 6.5.
type execMsg struct ***REMOVED***
	Command string
***REMOVED***

// Start runs cmd on the remote host. Typically, the remote
// server passes cmd to the shell for interpretation.
// A Session only accepts one call to Run, Start or Shell.
func (s *Session) Start(cmd string) error ***REMOVED***
	if s.started ***REMOVED***
		return errors.New("ssh: session already started")
	***REMOVED***
	req := execMsg***REMOVED***
		Command: cmd,
	***REMOVED***

	ok, err := s.ch.SendRequest("exec", true, Marshal(&req))
	if err == nil && !ok ***REMOVED***
		err = fmt.Errorf("ssh: command %v failed", cmd)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.start()
***REMOVED***

// Run runs cmd on the remote host. Typically, the remote
// server passes cmd to the shell for interpretation.
// A Session only accepts one call to Run, Start, Shell, Output,
// or CombinedOutput.
//
// The returned error is nil if the command runs, has no problems
// copying stdin, stdout, and stderr, and exits with a zero exit
// status.
//
// If the remote server does not send an exit status, an error of type
// *ExitMissingError is returned. If the command completes
// unsuccessfully or is interrupted by a signal, the error is of type
// *ExitError. Other error types may be returned for I/O problems.
func (s *Session) Run(cmd string) error ***REMOVED***
	err := s.Start(cmd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.Wait()
***REMOVED***

// Output runs cmd on the remote host and returns its standard output.
func (s *Session) Output(cmd string) ([]byte, error) ***REMOVED***
	if s.Stdout != nil ***REMOVED***
		return nil, errors.New("ssh: Stdout already set")
	***REMOVED***
	var b bytes.Buffer
	s.Stdout = &b
	err := s.Run(cmd)
	return b.Bytes(), err
***REMOVED***

type singleWriter struct ***REMOVED***
	b  bytes.Buffer
	mu sync.Mutex
***REMOVED***

func (w *singleWriter) Write(p []byte) (int, error) ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
***REMOVED***

// CombinedOutput runs cmd on the remote host and returns its combined
// standard output and standard error.
func (s *Session) CombinedOutput(cmd string) ([]byte, error) ***REMOVED***
	if s.Stdout != nil ***REMOVED***
		return nil, errors.New("ssh: Stdout already set")
	***REMOVED***
	if s.Stderr != nil ***REMOVED***
		return nil, errors.New("ssh: Stderr already set")
	***REMOVED***
	var b singleWriter
	s.Stdout = &b
	s.Stderr = &b
	err := s.Run(cmd)
	return b.b.Bytes(), err
***REMOVED***

// Shell starts a login shell on the remote host. A Session only
// accepts one call to Run, Start, Shell, Output, or CombinedOutput.
func (s *Session) Shell() error ***REMOVED***
	if s.started ***REMOVED***
		return errors.New("ssh: session already started")
	***REMOVED***

	ok, err := s.ch.SendRequest("shell", true, nil)
	if err == nil && !ok ***REMOVED***
		return errors.New("ssh: could not start shell")
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.start()
***REMOVED***

func (s *Session) start() error ***REMOVED***
	s.started = true

	type F func(*Session)
	for _, setupFd := range []F***REMOVED***(*Session).stdin, (*Session).stdout, (*Session).stderr***REMOVED*** ***REMOVED***
		setupFd(s)
	***REMOVED***

	s.errors = make(chan error, len(s.copyFuncs))
	for _, fn := range s.copyFuncs ***REMOVED***
		go func(fn func() error) ***REMOVED***
			s.errors <- fn()
		***REMOVED***(fn)
	***REMOVED***
	return nil
***REMOVED***

// Wait waits for the remote command to exit.
//
// The returned error is nil if the command runs, has no problems
// copying stdin, stdout, and stderr, and exits with a zero exit
// status.
//
// If the remote server does not send an exit status, an error of type
// *ExitMissingError is returned. If the command completes
// unsuccessfully or is interrupted by a signal, the error is of type
// *ExitError. Other error types may be returned for I/O problems.
func (s *Session) Wait() error ***REMOVED***
	if !s.started ***REMOVED***
		return errors.New("ssh: session not started")
	***REMOVED***
	waitErr := <-s.exitStatus

	if s.stdinPipeWriter != nil ***REMOVED***
		s.stdinPipeWriter.Close()
	***REMOVED***
	var copyError error
	for range s.copyFuncs ***REMOVED***
		if err := <-s.errors; err != nil && copyError == nil ***REMOVED***
			copyError = err
		***REMOVED***
	***REMOVED***
	if waitErr != nil ***REMOVED***
		return waitErr
	***REMOVED***
	return copyError
***REMOVED***

func (s *Session) wait(reqs <-chan *Request) error ***REMOVED***
	wm := Waitmsg***REMOVED***status: -1***REMOVED***
	// Wait for msg channel to be closed before returning.
	for msg := range reqs ***REMOVED***
		switch msg.Type ***REMOVED***
		case "exit-status":
			wm.status = int(binary.BigEndian.Uint32(msg.Payload))
		case "exit-signal":
			var sigval struct ***REMOVED***
				Signal     string
				CoreDumped bool
				Error      string
				Lang       string
			***REMOVED***
			if err := Unmarshal(msg.Payload, &sigval); err != nil ***REMOVED***
				return err
			***REMOVED***

			// Must sanitize strings?
			wm.signal = sigval.Signal
			wm.msg = sigval.Error
			wm.lang = sigval.Lang
		default:
			// This handles keepalives and matches
			// OpenSSH's behaviour.
			if msg.WantReply ***REMOVED***
				msg.Reply(false, nil)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if wm.status == 0 ***REMOVED***
		return nil
	***REMOVED***
	if wm.status == -1 ***REMOVED***
		// exit-status was never sent from server
		if wm.signal == "" ***REMOVED***
			// signal was not sent either.  RFC 4254
			// section 6.10 recommends against this
			// behavior, but it is allowed, so we let
			// clients handle it.
			return &ExitMissingError***REMOVED******REMOVED***
		***REMOVED***
		wm.status = 128
		if _, ok := signals[Signal(wm.signal)]; ok ***REMOVED***
			wm.status += signals[Signal(wm.signal)]
		***REMOVED***
	***REMOVED***

	return &ExitError***REMOVED***wm***REMOVED***
***REMOVED***

// ExitMissingError is returned if a session is torn down cleanly, but
// the server sends no confirmation of the exit status.
type ExitMissingError struct***REMOVED******REMOVED***

func (e *ExitMissingError) Error() string ***REMOVED***
	return "wait: remote command exited without exit status or exit signal"
***REMOVED***

func (s *Session) stdin() ***REMOVED***
	if s.stdinpipe ***REMOVED***
		return
	***REMOVED***
	var stdin io.Reader
	if s.Stdin == nil ***REMOVED***
		stdin = new(bytes.Buffer)
	***REMOVED*** else ***REMOVED***
		r, w := io.Pipe()
		go func() ***REMOVED***
			_, err := io.Copy(w, s.Stdin)
			w.CloseWithError(err)
		***REMOVED***()
		stdin, s.stdinPipeWriter = r, w
	***REMOVED***
	s.copyFuncs = append(s.copyFuncs, func() error ***REMOVED***
		_, err := io.Copy(s.ch, stdin)
		if err1 := s.ch.CloseWrite(); err == nil && err1 != io.EOF ***REMOVED***
			err = err1
		***REMOVED***
		return err
	***REMOVED***)
***REMOVED***

func (s *Session) stdout() ***REMOVED***
	if s.stdoutpipe ***REMOVED***
		return
	***REMOVED***
	if s.Stdout == nil ***REMOVED***
		s.Stdout = ioutil.Discard
	***REMOVED***
	s.copyFuncs = append(s.copyFuncs, func() error ***REMOVED***
		_, err := io.Copy(s.Stdout, s.ch)
		return err
	***REMOVED***)
***REMOVED***

func (s *Session) stderr() ***REMOVED***
	if s.stderrpipe ***REMOVED***
		return
	***REMOVED***
	if s.Stderr == nil ***REMOVED***
		s.Stderr = ioutil.Discard
	***REMOVED***
	s.copyFuncs = append(s.copyFuncs, func() error ***REMOVED***
		_, err := io.Copy(s.Stderr, s.ch.Stderr())
		return err
	***REMOVED***)
***REMOVED***

// sessionStdin reroutes Close to CloseWrite.
type sessionStdin struct ***REMOVED***
	io.Writer
	ch Channel
***REMOVED***

func (s *sessionStdin) Close() error ***REMOVED***
	return s.ch.CloseWrite()
***REMOVED***

// StdinPipe returns a pipe that will be connected to the
// remote command's standard input when the command starts.
func (s *Session) StdinPipe() (io.WriteCloser, error) ***REMOVED***
	if s.Stdin != nil ***REMOVED***
		return nil, errors.New("ssh: Stdin already set")
	***REMOVED***
	if s.started ***REMOVED***
		return nil, errors.New("ssh: StdinPipe after process started")
	***REMOVED***
	s.stdinpipe = true
	return &sessionStdin***REMOVED***s.ch, s.ch***REMOVED***, nil
***REMOVED***

// StdoutPipe returns a pipe that will be connected to the
// remote command's standard output when the command starts.
// There is a fixed amount of buffering that is shared between
// stdout and stderr streams. If the StdoutPipe reader is
// not serviced fast enough it may eventually cause the
// remote command to block.
func (s *Session) StdoutPipe() (io.Reader, error) ***REMOVED***
	if s.Stdout != nil ***REMOVED***
		return nil, errors.New("ssh: Stdout already set")
	***REMOVED***
	if s.started ***REMOVED***
		return nil, errors.New("ssh: StdoutPipe after process started")
	***REMOVED***
	s.stdoutpipe = true
	return s.ch, nil
***REMOVED***

// StderrPipe returns a pipe that will be connected to the
// remote command's standard error when the command starts.
// There is a fixed amount of buffering that is shared between
// stdout and stderr streams. If the StderrPipe reader is
// not serviced fast enough it may eventually cause the
// remote command to block.
func (s *Session) StderrPipe() (io.Reader, error) ***REMOVED***
	if s.Stderr != nil ***REMOVED***
		return nil, errors.New("ssh: Stderr already set")
	***REMOVED***
	if s.started ***REMOVED***
		return nil, errors.New("ssh: StderrPipe after process started")
	***REMOVED***
	s.stderrpipe = true
	return s.ch.Stderr(), nil
***REMOVED***

// newSession returns a new interactive session on the remote host.
func newSession(ch Channel, reqs <-chan *Request) (*Session, error) ***REMOVED***
	s := &Session***REMOVED***
		ch: ch,
	***REMOVED***
	s.exitStatus = make(chan error, 1)
	go func() ***REMOVED***
		s.exitStatus <- s.wait(reqs)
	***REMOVED***()

	return s, nil
***REMOVED***

// An ExitError reports unsuccessful completion of a remote command.
type ExitError struct ***REMOVED***
	Waitmsg
***REMOVED***

func (e *ExitError) Error() string ***REMOVED***
	return e.Waitmsg.String()
***REMOVED***

// Waitmsg stores the information about an exited remote command
// as reported by Wait.
type Waitmsg struct ***REMOVED***
	status int
	signal string
	msg    string
	lang   string
***REMOVED***

// ExitStatus returns the exit status of the remote command.
func (w Waitmsg) ExitStatus() int ***REMOVED***
	return w.status
***REMOVED***

// Signal returns the exit signal of the remote command if
// it was terminated violently.
func (w Waitmsg) Signal() string ***REMOVED***
	return w.signal
***REMOVED***

// Msg returns the exit message given by the remote command
func (w Waitmsg) Msg() string ***REMOVED***
	return w.msg
***REMOVED***

// Lang returns the language tag. See RFC 3066
func (w Waitmsg) Lang() string ***REMOVED***
	return w.lang
***REMOVED***

func (w Waitmsg) String() string ***REMOVED***
	str := fmt.Sprintf("Process exited with status %v", w.status)
	if w.signal != "" ***REMOVED***
		str += fmt.Sprintf(" from signal %v", w.signal)
	***REMOVED***
	if w.msg != "" ***REMOVED***
		str += fmt.Sprintf(". Reason was: %v", w.msg)
	***REMOVED***
	return str
***REMOVED***
