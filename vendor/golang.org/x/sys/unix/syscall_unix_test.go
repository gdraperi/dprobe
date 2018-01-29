// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package unix_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

// Tests that below functions, structures and constants are consistent
// on all Unix-like systems.
func _() ***REMOVED***
	// program scheduling priority functions and constants
	var (
		_ func(int, int, int) error   = unix.Setpriority
		_ func(int, int) (int, error) = unix.Getpriority
	)
	const (
		_ int = unix.PRIO_USER
		_ int = unix.PRIO_PROCESS
		_ int = unix.PRIO_PGRP
	)

	// termios constants
	const (
		_ int = unix.TCIFLUSH
		_ int = unix.TCIOFLUSH
		_ int = unix.TCOFLUSH
	)

	// fcntl file locking structure and constants
	var (
		_ = unix.Flock_t***REMOVED***
			Type:   int16(0),
			Whence: int16(0),
			Start:  int64(0),
			Len:    int64(0),
			Pid:    int32(0),
		***REMOVED***
	)
	const (
		_ = unix.F_GETLK
		_ = unix.F_SETLK
		_ = unix.F_SETLKW
	)
***REMOVED***

// TestFcntlFlock tests whether the file locking structure matches
// the calling convention of each kernel.
func TestFcntlFlock(t *testing.T) ***REMOVED***
	name := filepath.Join(os.TempDir(), "TestFcntlFlock")
	fd, err := unix.Open(name, unix.O_CREAT|unix.O_RDWR|unix.O_CLOEXEC, 0)
	if err != nil ***REMOVED***
		t.Fatalf("Open failed: %v", err)
	***REMOVED***
	defer unix.Unlink(name)
	defer unix.Close(fd)
	flock := unix.Flock_t***REMOVED***
		Type:  unix.F_RDLCK,
		Start: 0, Len: 0, Whence: 1,
	***REMOVED***
	if err := unix.FcntlFlock(uintptr(fd), unix.F_GETLK, &flock); err != nil ***REMOVED***
		t.Fatalf("FcntlFlock failed: %v", err)
	***REMOVED***
***REMOVED***

// TestPassFD tests passing a file descriptor over a Unix socket.
//
// This test involved both a parent and child process. The parent
// process is invoked as a normal test, with "go test", which then
// runs the child process by running the current test binary with args
// "-test.run=^TestPassFD$" and an environment variable used to signal
// that the test should become the child process instead.
func TestPassFD(t *testing.T) ***REMOVED***
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" ***REMOVED***
		passFDChild()
		return
	***REMOVED***

	tempDir, err := ioutil.TempDir("", "TestPassFD")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tempDir)

	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	if err != nil ***REMOVED***
		t.Fatalf("Socketpair: %v", err)
	***REMOVED***
	defer unix.Close(fds[0])
	defer unix.Close(fds[1])
	writeFile := os.NewFile(uintptr(fds[0]), "child-writes")
	readFile := os.NewFile(uintptr(fds[1]), "parent-reads")
	defer writeFile.Close()
	defer readFile.Close()

	cmd := exec.Command(os.Args[0], "-test.run=^TestPassFD$", "--", tempDir)
	cmd.Env = []string***REMOVED***"GO_WANT_HELPER_PROCESS=1"***REMOVED***
	if lp := os.Getenv("LD_LIBRARY_PATH"); lp != "" ***REMOVED***
		cmd.Env = append(cmd.Env, "LD_LIBRARY_PATH="+lp)
	***REMOVED***
	cmd.ExtraFiles = []*os.File***REMOVED***writeFile***REMOVED***

	out, err := cmd.CombinedOutput()
	if len(out) > 0 || err != nil ***REMOVED***
		t.Fatalf("child process: %q, %v", out, err)
	***REMOVED***

	c, err := net.FileConn(readFile)
	if err != nil ***REMOVED***
		t.Fatalf("FileConn: %v", err)
	***REMOVED***
	defer c.Close()

	uc, ok := c.(*net.UnixConn)
	if !ok ***REMOVED***
		t.Fatalf("unexpected FileConn type; expected UnixConn, got %T", c)
	***REMOVED***

	buf := make([]byte, 32) // expect 1 byte
	oob := make([]byte, 32) // expect 24 bytes
	closeUnix := time.AfterFunc(5*time.Second, func() ***REMOVED***
		t.Logf("timeout reading from unix socket")
		uc.Close()
	***REMOVED***)
	_, oobn, _, _, err := uc.ReadMsgUnix(buf, oob)
	if err != nil ***REMOVED***
		t.Fatalf("ReadMsgUnix: %v", err)
	***REMOVED***
	closeUnix.Stop()

	scms, err := unix.ParseSocketControlMessage(oob[:oobn])
	if err != nil ***REMOVED***
		t.Fatalf("ParseSocketControlMessage: %v", err)
	***REMOVED***
	if len(scms) != 1 ***REMOVED***
		t.Fatalf("expected 1 SocketControlMessage; got scms = %#v", scms)
	***REMOVED***
	scm := scms[0]
	gotFds, err := unix.ParseUnixRights(&scm)
	if err != nil ***REMOVED***
		t.Fatalf("unix.ParseUnixRights: %v", err)
	***REMOVED***
	if len(gotFds) != 1 ***REMOVED***
		t.Fatalf("wanted 1 fd; got %#v", gotFds)
	***REMOVED***

	f := os.NewFile(uintptr(gotFds[0]), "fd-from-child")
	defer f.Close()

	got, err := ioutil.ReadAll(f)
	want := "Hello from child process!\n"
	if string(got) != want ***REMOVED***
		t.Errorf("child process ReadAll: %q, %v; want %q", got, err, want)
	***REMOVED***
***REMOVED***

// passFDChild is the child process used by TestPassFD.
func passFDChild() ***REMOVED***
	defer os.Exit(0)

	// Look for our fd. It should be fd 3, but we work around an fd leak
	// bug here (http://golang.org/issue/2603) to let it be elsewhere.
	var uc *net.UnixConn
	for fd := uintptr(3); fd <= 10; fd++ ***REMOVED***
		f := os.NewFile(fd, "unix-conn")
		var ok bool
		netc, _ := net.FileConn(f)
		uc, ok = netc.(*net.UnixConn)
		if ok ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if uc == nil ***REMOVED***
		fmt.Println("failed to find unix fd")
		return
	***REMOVED***

	// Make a file f to send to our parent process on uc.
	// We make it in tempDir, which our parent will clean up.
	flag.Parse()
	tempDir := flag.Arg(0)
	f, err := ioutil.TempFile(tempDir, "")
	if err != nil ***REMOVED***
		fmt.Printf("TempFile: %v", err)
		return
	***REMOVED***

	f.Write([]byte("Hello from child process!\n"))
	f.Seek(0, 0)

	rights := unix.UnixRights(int(f.Fd()))
	dummyByte := []byte("x")
	n, oobn, err := uc.WriteMsgUnix(dummyByte, rights, nil)
	if err != nil ***REMOVED***
		fmt.Printf("WriteMsgUnix: %v", err)
		return
	***REMOVED***
	if n != 1 || oobn != len(rights) ***REMOVED***
		fmt.Printf("WriteMsgUnix = %d, %d; want 1, %d", n, oobn, len(rights))
		return
	***REMOVED***
***REMOVED***

// TestUnixRightsRoundtrip tests that UnixRights, ParseSocketControlMessage,
// and ParseUnixRights are able to successfully round-trip lists of file descriptors.
func TestUnixRightsRoundtrip(t *testing.T) ***REMOVED***
	testCases := [...][][]int***REMOVED***
		***REMOVED******REMOVED***42***REMOVED******REMOVED***,
		***REMOVED******REMOVED***1, 2***REMOVED******REMOVED***,
		***REMOVED******REMOVED***3, 4, 5***REMOVED******REMOVED***,
		***REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED******REMOVED***1, 2***REMOVED***, ***REMOVED***3, 4, 5***REMOVED***, ***REMOVED******REMOVED***, ***REMOVED***7***REMOVED******REMOVED***,
	***REMOVED***
	for _, testCase := range testCases ***REMOVED***
		b := []byte***REMOVED******REMOVED***
		var n int
		for _, fds := range testCase ***REMOVED***
			// Last assignment to n wins
			n = len(b) + unix.CmsgLen(4*len(fds))
			b = append(b, unix.UnixRights(fds...)...)
		***REMOVED***
		// Truncate b
		b = b[:n]

		scms, err := unix.ParseSocketControlMessage(b)
		if err != nil ***REMOVED***
			t.Fatalf("ParseSocketControlMessage: %v", err)
		***REMOVED***
		if len(scms) != len(testCase) ***REMOVED***
			t.Fatalf("expected %v SocketControlMessage; got scms = %#v", len(testCase), scms)
		***REMOVED***
		for i, scm := range scms ***REMOVED***
			gotFds, err := unix.ParseUnixRights(&scm)
			if err != nil ***REMOVED***
				t.Fatalf("ParseUnixRights: %v", err)
			***REMOVED***
			wantFds := testCase[i]
			if len(gotFds) != len(wantFds) ***REMOVED***
				t.Fatalf("expected %v fds, got %#v", len(wantFds), gotFds)
			***REMOVED***
			for j, fd := range gotFds ***REMOVED***
				if fd != wantFds[j] ***REMOVED***
					t.Fatalf("expected fd %v, got %v", wantFds[j], fd)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRlimit(t *testing.T) ***REMOVED***
	var rlimit, zero unix.Rlimit
	err := unix.Getrlimit(unix.RLIMIT_NOFILE, &rlimit)
	if err != nil ***REMOVED***
		t.Fatalf("Getrlimit: save failed: %v", err)
	***REMOVED***
	if zero == rlimit ***REMOVED***
		t.Fatalf("Getrlimit: save failed: got zero value %#v", rlimit)
	***REMOVED***
	set := rlimit
	set.Cur = set.Max - 1
	err = unix.Setrlimit(unix.RLIMIT_NOFILE, &set)
	if err != nil ***REMOVED***
		t.Fatalf("Setrlimit: set failed: %#v %v", set, err)
	***REMOVED***
	var get unix.Rlimit
	err = unix.Getrlimit(unix.RLIMIT_NOFILE, &get)
	if err != nil ***REMOVED***
		t.Fatalf("Getrlimit: get failed: %v", err)
	***REMOVED***
	set = rlimit
	set.Cur = set.Max - 1
	if set != get ***REMOVED***
		// Seems like Darwin requires some privilege to
		// increase the soft limit of rlimit sandbox, though
		// Setrlimit never reports an error.
		switch runtime.GOOS ***REMOVED***
		case "darwin":
		default:
			t.Fatalf("Rlimit: change failed: wanted %#v got %#v", set, get)
		***REMOVED***
	***REMOVED***
	err = unix.Setrlimit(unix.RLIMIT_NOFILE, &rlimit)
	if err != nil ***REMOVED***
		t.Fatalf("Setrlimit: restore failed: %#v %v", rlimit, err)
	***REMOVED***
***REMOVED***

func TestSeekFailure(t *testing.T) ***REMOVED***
	_, err := unix.Seek(-1, 0, 0)
	if err == nil ***REMOVED***
		t.Fatalf("Seek(-1, 0, 0) did not fail")
	***REMOVED***
	str := err.Error() // used to crash on Linux
	t.Logf("Seek: %v", str)
	if str == "" ***REMOVED***
		t.Fatalf("Seek(-1, 0, 0) return error with empty message")
	***REMOVED***
***REMOVED***

func TestDup(t *testing.T) ***REMOVED***
	file, err := ioutil.TempFile("", "TestDup")
	if err != nil ***REMOVED***
		t.Fatalf("Tempfile failed: %v", err)
	***REMOVED***
	defer os.Remove(file.Name())
	defer file.Close()
	f := int(file.Fd())

	newFd, err := unix.Dup(f)
	if err != nil ***REMOVED***
		t.Fatalf("Dup: %v", err)
	***REMOVED***

	err = unix.Dup2(newFd, newFd+1)
	if err != nil ***REMOVED***
		t.Fatalf("Dup2: %v", err)
	***REMOVED***

	b1 := []byte("Test123")
	b2 := make([]byte, 7)
	_, err = unix.Write(newFd+1, b1)
	if err != nil ***REMOVED***
		t.Fatalf("Write to dup2 fd failed: %v", err)
	***REMOVED***
	_, err = unix.Seek(f, 0, 0)
	if err != nil ***REMOVED***
		t.Fatalf("Seek failed: %v", err)
	***REMOVED***
	_, err = unix.Read(f, b2)
	if err != nil ***REMOVED***
		t.Fatalf("Read back failed: %v", err)
	***REMOVED***
	if string(b1) != string(b2) ***REMOVED***
		t.Errorf("Dup: stdout write not in file, expected %v, got %v", string(b1), string(b2))
	***REMOVED***
***REMOVED***

func TestPoll(t *testing.T) ***REMOVED***
	f, cleanup := mktmpfifo(t)
	defer cleanup()

	const timeout = 100

	ok := make(chan bool, 1)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-time.After(10 * timeout * time.Millisecond):
			t.Errorf("Poll: failed to timeout after %d milliseconds", 10*timeout)
		case <-ok:
		***REMOVED***
	***REMOVED***()

	fds := []unix.PollFd***REMOVED******REMOVED***Fd: int32(f.Fd()), Events: unix.POLLIN***REMOVED******REMOVED***
	n, err := unix.Poll(fds, timeout)
	ok <- true
	if err != nil ***REMOVED***
		t.Errorf("Poll: unexpected error: %v", err)
		return
	***REMOVED***
	if n != 0 ***REMOVED***
		t.Errorf("Poll: wrong number of events: got %v, expected %v", n, 0)
		return
	***REMOVED***
***REMOVED***

func TestGetwd(t *testing.T) ***REMOVED***
	fd, err := os.Open(".")
	if err != nil ***REMOVED***
		t.Fatalf("Open .: %s", err)
	***REMOVED***
	defer fd.Close()
	// These are chosen carefully not to be symlinks on a Mac
	// (unlike, say, /var, /etc)
	dirs := []string***REMOVED***"/", "/usr/bin"***REMOVED***
	if runtime.GOOS == "darwin" ***REMOVED***
		switch runtime.GOARCH ***REMOVED***
		case "arm", "arm64":
			d1, err := ioutil.TempDir("", "d1")
			if err != nil ***REMOVED***
				t.Fatalf("TempDir: %v", err)
			***REMOVED***
			d2, err := ioutil.TempDir("", "d2")
			if err != nil ***REMOVED***
				t.Fatalf("TempDir: %v", err)
			***REMOVED***
			dirs = []string***REMOVED***d1, d2***REMOVED***
		***REMOVED***
	***REMOVED***
	oldwd := os.Getenv("PWD")
	for _, d := range dirs ***REMOVED***
		err = os.Chdir(d)
		if err != nil ***REMOVED***
			t.Fatalf("Chdir: %v", err)
		***REMOVED***
		pwd, err := unix.Getwd()
		if err != nil ***REMOVED***
			t.Fatalf("Getwd in %s: %s", d, err)
		***REMOVED***
		os.Setenv("PWD", oldwd)
		err = fd.Chdir()
		if err != nil ***REMOVED***
			// We changed the current directory and cannot go back.
			// Don't let the tests continue; they'll scribble
			// all over some other directory.
			fmt.Fprintf(os.Stderr, "fchdir back to dot failed: %s\n", err)
			os.Exit(1)
		***REMOVED***
		if pwd != d ***REMOVED***
			t.Fatalf("Getwd returned %q want %q", pwd, d)
		***REMOVED***
	***REMOVED***
***REMOVED***

// mktmpfifo creates a temporary FIFO and provides a cleanup function.
func mktmpfifo(t *testing.T) (*os.File, func()) ***REMOVED***
	err := unix.Mkfifo("fifo", 0666)
	if err != nil ***REMOVED***
		t.Fatalf("mktmpfifo: failed to create FIFO: %v", err)
	***REMOVED***

	f, err := os.OpenFile("fifo", os.O_RDWR, 0666)
	if err != nil ***REMOVED***
		os.Remove("fifo")
		t.Fatalf("mktmpfifo: failed to open FIFO: %v", err)
	***REMOVED***

	return f, func() ***REMOVED***
		f.Close()
		os.Remove("fifo")
	***REMOVED***
***REMOVED***
