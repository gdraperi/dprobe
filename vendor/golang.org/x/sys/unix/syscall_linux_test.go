// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package unix_test

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestFchmodat(t *testing.T) ***REMOVED***
	defer chtmpdir(t)()

	touch(t, "file1")
	err := os.Symlink("file1", "symlink1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = unix.Fchmodat(unix.AT_FDCWD, "symlink1", 0444, 0)
	if err != nil ***REMOVED***
		t.Fatalf("Fchmodat: unexpected error: %v", err)
	***REMOVED***

	fi, err := os.Stat("file1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if fi.Mode() != 0444 ***REMOVED***
		t.Errorf("Fchmodat: failed to change mode: expected %v, got %v", 0444, fi.Mode())
	***REMOVED***

	err = unix.Fchmodat(unix.AT_FDCWD, "symlink1", 0444, unix.AT_SYMLINK_NOFOLLOW)
	if err != unix.EOPNOTSUPP ***REMOVED***
		t.Fatalf("Fchmodat: unexpected error: %v, expected EOPNOTSUPP", err)
	***REMOVED***
***REMOVED***

func TestIoctlGetInt(t *testing.T) ***REMOVED***
	f, err := os.Open("/dev/random")
	if err != nil ***REMOVED***
		t.Fatalf("failed to open device: %v", err)
	***REMOVED***
	defer f.Close()

	v, err := unix.IoctlGetInt(int(f.Fd()), unix.RNDGETENTCNT)
	if err != nil ***REMOVED***
		t.Fatalf("failed to perform ioctl: %v", err)
	***REMOVED***

	t.Logf("%d bits of entropy available", v)
***REMOVED***

func TestPpoll(t *testing.T) ***REMOVED***
	f, cleanup := mktmpfifo(t)
	defer cleanup()

	const timeout = 100 * time.Millisecond

	ok := make(chan bool, 1)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-time.After(10 * timeout):
			t.Errorf("Ppoll: failed to timeout after %d", 10*timeout)
		case <-ok:
		***REMOVED***
	***REMOVED***()

	fds := []unix.PollFd***REMOVED******REMOVED***Fd: int32(f.Fd()), Events: unix.POLLIN***REMOVED******REMOVED***
	timeoutTs := unix.NsecToTimespec(int64(timeout))
	n, err := unix.Ppoll(fds, &timeoutTs, nil)
	ok <- true
	if err != nil ***REMOVED***
		t.Errorf("Ppoll: unexpected error: %v", err)
		return
	***REMOVED***
	if n != 0 ***REMOVED***
		t.Errorf("Ppoll: wrong number of events: got %v, expected %v", n, 0)
		return
	***REMOVED***
***REMOVED***

func TestTime(t *testing.T) ***REMOVED***
	var ut unix.Time_t
	ut2, err := unix.Time(&ut)
	if err != nil ***REMOVED***
		t.Fatalf("Time: %v", err)
	***REMOVED***
	if ut != ut2 ***REMOVED***
		t.Errorf("Time: return value %v should be equal to argument %v", ut2, ut)
	***REMOVED***

	var now time.Time

	for i := 0; i < 10; i++ ***REMOVED***
		ut, err = unix.Time(nil)
		if err != nil ***REMOVED***
			t.Fatalf("Time: %v", err)
		***REMOVED***

		now = time.Now()

		if int64(ut) == now.Unix() ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	t.Errorf("Time: return value %v should be nearly equal to time.Now().Unix() %v", ut, now.Unix())
***REMOVED***

func TestUtime(t *testing.T) ***REMOVED***
	defer chtmpdir(t)()

	touch(t, "file1")

	buf := &unix.Utimbuf***REMOVED***
		Modtime: 12345,
	***REMOVED***

	err := unix.Utime("file1", buf)
	if err != nil ***REMOVED***
		t.Fatalf("Utime: %v", err)
	***REMOVED***

	fi, err := os.Stat("file1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if fi.ModTime().Unix() != 12345 ***REMOVED***
		t.Errorf("Utime: failed to change modtime: expected %v, got %v", 12345, fi.ModTime().Unix())
	***REMOVED***
***REMOVED***

func TestUtimesNanoAt(t *testing.T) ***REMOVED***
	defer chtmpdir(t)()

	symlink := "symlink1"
	os.Remove(symlink)
	err := os.Symlink("nonexisting", symlink)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ts := []unix.Timespec***REMOVED***
		***REMOVED***Sec: 1111, Nsec: 2222***REMOVED***,
		***REMOVED***Sec: 3333, Nsec: 4444***REMOVED***,
	***REMOVED***
	err = unix.UtimesNanoAt(unix.AT_FDCWD, symlink, ts, unix.AT_SYMLINK_NOFOLLOW)
	if err != nil ***REMOVED***
		t.Fatalf("UtimesNanoAt: %v", err)
	***REMOVED***

	var st unix.Stat_t
	err = unix.Lstat(symlink, &st)
	if err != nil ***REMOVED***
		t.Fatalf("Lstat: %v", err)
	***REMOVED***
	if st.Atim != ts[0] ***REMOVED***
		t.Errorf("UtimesNanoAt: wrong atime: %v", st.Atim)
	***REMOVED***
	if st.Mtim != ts[1] ***REMOVED***
		t.Errorf("UtimesNanoAt: wrong mtime: %v", st.Mtim)
	***REMOVED***
***REMOVED***

func TestGetrlimit(t *testing.T) ***REMOVED***
	var rlim unix.Rlimit
	err := unix.Getrlimit(unix.RLIMIT_AS, &rlim)
	if err != nil ***REMOVED***
		t.Fatalf("Getrlimit: %v", err)
	***REMOVED***
***REMOVED***

func TestSelect(t *testing.T) ***REMOVED***
	_, err := unix.Select(0, nil, nil, nil, &unix.Timeval***REMOVED***Sec: 0, Usec: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Select: %v", err)
	***REMOVED***

	dur := 150 * time.Millisecond
	tv := unix.NsecToTimeval(int64(dur))
	start := time.Now()
	_, err = unix.Select(0, nil, nil, nil, &tv)
	took := time.Since(start)
	if err != nil ***REMOVED***
		t.Fatalf("Select: %v", err)
	***REMOVED***

	if took < dur ***REMOVED***
		t.Errorf("Select: timeout should have been at least %v, got %v", dur, took)
	***REMOVED***
***REMOVED***

func TestPselect(t *testing.T) ***REMOVED***
	_, err := unix.Pselect(0, nil, nil, nil, &unix.Timespec***REMOVED***Sec: 0, Nsec: 0***REMOVED***, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Pselect: %v", err)
	***REMOVED***

	dur := 2500 * time.Microsecond
	ts := unix.NsecToTimespec(int64(dur))
	start := time.Now()
	_, err = unix.Pselect(0, nil, nil, nil, &ts, nil)
	took := time.Since(start)
	if err != nil ***REMOVED***
		t.Fatalf("Pselect: %v", err)
	***REMOVED***

	if took < dur ***REMOVED***
		t.Errorf("Pselect: timeout should have been at least %v, got %v", dur, took)
	***REMOVED***
***REMOVED***

func TestFstatat(t *testing.T) ***REMOVED***
	defer chtmpdir(t)()

	touch(t, "file1")

	var st1 unix.Stat_t
	err := unix.Stat("file1", &st1)
	if err != nil ***REMOVED***
		t.Fatalf("Stat: %v", err)
	***REMOVED***

	var st2 unix.Stat_t
	err = unix.Fstatat(unix.AT_FDCWD, "file1", &st2, 0)
	if err != nil ***REMOVED***
		t.Fatalf("Fstatat: %v", err)
	***REMOVED***

	if st1 != st2 ***REMOVED***
		t.Errorf("Fstatat: returned stat does not match Stat")
	***REMOVED***

	err = os.Symlink("file1", "symlink1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = unix.Lstat("symlink1", &st1)
	if err != nil ***REMOVED***
		t.Fatalf("Lstat: %v", err)
	***REMOVED***

	err = unix.Fstatat(unix.AT_FDCWD, "symlink1", &st2, unix.AT_SYMLINK_NOFOLLOW)
	if err != nil ***REMOVED***
		t.Fatalf("Fstatat: %v", err)
	***REMOVED***

	if st1 != st2 ***REMOVED***
		t.Errorf("Fstatat: returned stat does not match Lstat")
	***REMOVED***
***REMOVED***

func TestSchedSetaffinity(t *testing.T) ***REMOVED***
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var oldMask unix.CPUSet
	err := unix.SchedGetaffinity(0, &oldMask)
	if err != nil ***REMOVED***
		t.Fatalf("SchedGetaffinity: %v", err)
	***REMOVED***

	var newMask unix.CPUSet
	newMask.Zero()
	if newMask.Count() != 0 ***REMOVED***
		t.Errorf("CpuZero: didn't zero CPU set: %v", newMask)
	***REMOVED***
	cpu := 1
	newMask.Set(cpu)
	if newMask.Count() != 1 || !newMask.IsSet(cpu) ***REMOVED***
		t.Errorf("CpuSet: didn't set CPU %d in set: %v", cpu, newMask)
	***REMOVED***
	cpu = 5
	newMask.Set(cpu)
	if newMask.Count() != 2 || !newMask.IsSet(cpu) ***REMOVED***
		t.Errorf("CpuSet: didn't set CPU %d in set: %v", cpu, newMask)
	***REMOVED***
	newMask.Clear(cpu)
	if newMask.Count() != 1 || newMask.IsSet(cpu) ***REMOVED***
		t.Errorf("CpuClr: didn't clear CPU %d in set: %v", cpu, newMask)
	***REMOVED***

	err = unix.SchedSetaffinity(0, &newMask)
	if err != nil ***REMOVED***
		t.Fatalf("SchedSetaffinity: %v", err)
	***REMOVED***

	var gotMask unix.CPUSet
	err = unix.SchedGetaffinity(0, &gotMask)
	if err != nil ***REMOVED***
		t.Fatalf("SchedGetaffinity: %v", err)
	***REMOVED***

	if gotMask != newMask ***REMOVED***
		t.Errorf("SchedSetaffinity: returned affinity mask does not match set affinity mask")
	***REMOVED***

	// Restore old mask so it doesn't affect successive tests
	err = unix.SchedSetaffinity(0, &oldMask)
	if err != nil ***REMOVED***
		t.Fatalf("SchedSetaffinity: %v", err)
	***REMOVED***
***REMOVED***

func TestStatx(t *testing.T) ***REMOVED***
	var stx unix.Statx_t
	err := unix.Statx(unix.AT_FDCWD, ".", 0, 0, &stx)
	if err == unix.ENOSYS ***REMOVED***
		t.Skip("statx syscall is not available, skipping test")
	***REMOVED*** else if err != nil ***REMOVED***
		t.Fatalf("Statx: %v", err)
	***REMOVED***

	defer chtmpdir(t)()
	touch(t, "file1")

	var st unix.Stat_t
	err = unix.Stat("file1", &st)
	if err != nil ***REMOVED***
		t.Fatalf("Stat: %v", err)
	***REMOVED***

	flags := unix.AT_STATX_SYNC_AS_STAT
	err = unix.Statx(unix.AT_FDCWD, "file1", flags, unix.STATX_ALL, &stx)
	if err != nil ***REMOVED***
		t.Fatalf("Statx: %v", err)
	***REMOVED***

	if uint32(stx.Mode) != st.Mode ***REMOVED***
		t.Errorf("Statx: returned stat mode does not match Stat")
	***REMOVED***

	atime := unix.StatxTimestamp***REMOVED***Sec: int64(st.Atim.Sec), Nsec: uint32(st.Atim.Nsec)***REMOVED***
	ctime := unix.StatxTimestamp***REMOVED***Sec: int64(st.Ctim.Sec), Nsec: uint32(st.Ctim.Nsec)***REMOVED***
	mtime := unix.StatxTimestamp***REMOVED***Sec: int64(st.Mtim.Sec), Nsec: uint32(st.Mtim.Nsec)***REMOVED***

	if stx.Atime != atime ***REMOVED***
		t.Errorf("Statx: returned stat atime does not match Stat")
	***REMOVED***
	if stx.Ctime != ctime ***REMOVED***
		t.Errorf("Statx: returned stat ctime does not match Stat")
	***REMOVED***
	if stx.Mtime != mtime ***REMOVED***
		t.Errorf("Statx: returned stat mtime does not match Stat")
	***REMOVED***

	err = os.Symlink("file1", "symlink1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = unix.Lstat("symlink1", &st)
	if err != nil ***REMOVED***
		t.Fatalf("Lstat: %v", err)
	***REMOVED***

	err = unix.Statx(unix.AT_FDCWD, "symlink1", flags, unix.STATX_BASIC_STATS, &stx)
	if err != nil ***REMOVED***
		t.Fatalf("Statx: %v", err)
	***REMOVED***

	// follow symlink, expect a regulat file
	if stx.Mode&unix.S_IFREG == 0 ***REMOVED***
		t.Errorf("Statx: didn't follow symlink")
	***REMOVED***

	err = unix.Statx(unix.AT_FDCWD, "symlink1", flags|unix.AT_SYMLINK_NOFOLLOW, unix.STATX_ALL, &stx)
	if err != nil ***REMOVED***
		t.Fatalf("Statx: %v", err)
	***REMOVED***

	// follow symlink, expect a symlink
	if stx.Mode&unix.S_IFLNK == 0 ***REMOVED***
		t.Errorf("Statx: unexpectedly followed symlink")
	***REMOVED***
	if uint32(stx.Mode) != st.Mode ***REMOVED***
		t.Errorf("Statx: returned stat mode does not match Lstat")
	***REMOVED***

	atime = unix.StatxTimestamp***REMOVED***Sec: int64(st.Atim.Sec), Nsec: uint32(st.Atim.Nsec)***REMOVED***
	ctime = unix.StatxTimestamp***REMOVED***Sec: int64(st.Ctim.Sec), Nsec: uint32(st.Ctim.Nsec)***REMOVED***
	mtime = unix.StatxTimestamp***REMOVED***Sec: int64(st.Mtim.Sec), Nsec: uint32(st.Mtim.Nsec)***REMOVED***

	if stx.Atime != atime ***REMOVED***
		t.Errorf("Statx: returned stat atime does not match Lstat")
	***REMOVED***
	if stx.Ctime != ctime ***REMOVED***
		t.Errorf("Statx: returned stat ctime does not match Lstat")
	***REMOVED***
	if stx.Mtime != mtime ***REMOVED***
		t.Errorf("Statx: returned stat mtime does not match Lstat")
	***REMOVED***
***REMOVED***

// utilities taken from os/os_test.go

func touch(t *testing.T, name string) ***REMOVED***
	f, err := os.Create(name)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := f.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// chtmpdir changes the working directory to a new temporary directory and
// provides a cleanup function. Used when PWD is read-only.
func chtmpdir(t *testing.T) func() ***REMOVED***
	oldwd, err := os.Getwd()
	if err != nil ***REMOVED***
		t.Fatalf("chtmpdir: %v", err)
	***REMOVED***
	d, err := ioutil.TempDir("", "test")
	if err != nil ***REMOVED***
		t.Fatalf("chtmpdir: %v", err)
	***REMOVED***
	if err := os.Chdir(d); err != nil ***REMOVED***
		t.Fatalf("chtmpdir: %v", err)
	***REMOVED***
	return func() ***REMOVED***
		if err := os.Chdir(oldwd); err != nil ***REMOVED***
			t.Fatalf("chtmpdir: %v", err)
		***REMOVED***
		os.RemoveAll(d)
	***REMOVED***
***REMOVED***
