// +build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall" // only for exec
	"unsafe"

	"golang.org/x/sys/unix"
)

// If arg2 is nonzero, set the "child subreaper" attribute of the
// calling process; if arg2 is zero, unset the attribute.  When a
// process is marked as a child subreaper, all of the children
// that it creates, and their descendants, will be marked as
// having a subreaper.  In effect, a subreaper fulfills the role
// of init(1) for its descendant processes.  Upon termination of
// a process that is orphaned (i.e., its immediate parent has
// already terminated) and marked as having a subreaper, the
// nearest still living ancestor subreaper will receive a SIGCHLD
// signal and be able to wait(2) on the process to discover its
// termination status.
const PR_SET_CHILD_SUBREAPER = 36

type ParentDeathSignal int

func (p ParentDeathSignal) Restore() error ***REMOVED***
	if p == 0 ***REMOVED***
		return nil
	***REMOVED***
	current, err := GetParentDeathSignal()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if p == current ***REMOVED***
		return nil
	***REMOVED***
	return p.Set()
***REMOVED***

func (p ParentDeathSignal) Set() error ***REMOVED***
	return SetParentDeathSignal(uintptr(p))
***REMOVED***

func Execv(cmd string, args []string, env []string) error ***REMOVED***
	name, err := exec.LookPath(cmd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return syscall.Exec(name, args, env)
***REMOVED***

func Prlimit(pid, resource int, limit unix.Rlimit) error ***REMOVED***
	_, _, err := unix.RawSyscall6(unix.SYS_PRLIMIT64, uintptr(pid), uintptr(resource), uintptr(unsafe.Pointer(&limit)), uintptr(unsafe.Pointer(&limit)), 0, 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func SetParentDeathSignal(sig uintptr) error ***REMOVED***
	if err := unix.Prctl(unix.PR_SET_PDEATHSIG, sig, 0, 0, 0); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func GetParentDeathSignal() (ParentDeathSignal, error) ***REMOVED***
	var sig int
	if err := unix.Prctl(unix.PR_GET_PDEATHSIG, uintptr(unsafe.Pointer(&sig)), 0, 0, 0); err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return ParentDeathSignal(sig), nil
***REMOVED***

func SetKeepCaps() error ***REMOVED***
	if err := unix.Prctl(unix.PR_SET_KEEPCAPS, 1, 0, 0, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func ClearKeepCaps() error ***REMOVED***
	if err := unix.Prctl(unix.PR_SET_KEEPCAPS, 0, 0, 0, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func Setctty() error ***REMOVED***
	if err := unix.IoctlSetInt(0, unix.TIOCSCTTY, 0); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// RunningInUserNS detects whether we are currently running in a user namespace.
// Copied from github.com/lxc/lxd/shared/util.go
func RunningInUserNS() bool ***REMOVED***
	file, err := os.Open("/proc/self/uid_map")
	if err != nil ***REMOVED***
		// This kernel-provided file only exists if user namespaces are supported
		return false
	***REMOVED***
	defer file.Close()

	buf := bufio.NewReader(file)
	l, _, err := buf.ReadLine()
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	line := string(l)
	var a, b, c int64
	fmt.Sscanf(line, "%d %d %d", &a, &b, &c)
	/*
	 * We assume we are in the initial user namespace if we have a full
	 * range - 4294967295 uids starting at uid 0.
	 */
	if a == 0 && b == 0 && c == 4294967295 ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// SetSubreaper sets the value i as the subreaper setting for the calling process
func SetSubreaper(i int) error ***REMOVED***
	return unix.Prctl(PR_SET_CHILD_SUBREAPER, uintptr(i), 0, 0, 0)
***REMOVED***

// GetSubreaper returns the subreaper setting for the calling process
func GetSubreaper() (int, error) ***REMOVED***
	var i uintptr

	if err := unix.Prctl(unix.PR_GET_CHILD_SUBREAPER, uintptr(unsafe.Pointer(&i)), 0, 0, 0); err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	return int(i), nil
***REMOVED***
