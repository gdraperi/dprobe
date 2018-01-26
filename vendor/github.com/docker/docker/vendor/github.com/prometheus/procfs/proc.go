package procfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Proc provides information about a running process.
type Proc struct ***REMOVED***
	// The process ID.
	PID int

	fs FS
***REMOVED***

// Procs represents a list of Proc structs.
type Procs []Proc

func (p Procs) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p Procs) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***
func (p Procs) Less(i, j int) bool ***REMOVED*** return p[i].PID < p[j].PID ***REMOVED***

// Self returns a process for the current process read via /proc/self.
func Self() (Proc, error) ***REMOVED***
	fs, err := NewFS(DefaultMountPoint)
	if err != nil ***REMOVED***
		return Proc***REMOVED******REMOVED***, err
	***REMOVED***
	return fs.Self()
***REMOVED***

// NewProc returns a process for the given pid under /proc.
func NewProc(pid int) (Proc, error) ***REMOVED***
	fs, err := NewFS(DefaultMountPoint)
	if err != nil ***REMOVED***
		return Proc***REMOVED******REMOVED***, err
	***REMOVED***
	return fs.NewProc(pid)
***REMOVED***

// AllProcs returns a list of all currently available processes under /proc.
func AllProcs() (Procs, error) ***REMOVED***
	fs, err := NewFS(DefaultMountPoint)
	if err != nil ***REMOVED***
		return Procs***REMOVED******REMOVED***, err
	***REMOVED***
	return fs.AllProcs()
***REMOVED***

// Self returns a process for the current process.
func (fs FS) Self() (Proc, error) ***REMOVED***
	p, err := os.Readlink(fs.Path("self"))
	if err != nil ***REMOVED***
		return Proc***REMOVED******REMOVED***, err
	***REMOVED***
	pid, err := strconv.Atoi(strings.Replace(p, string(fs), "", -1))
	if err != nil ***REMOVED***
		return Proc***REMOVED******REMOVED***, err
	***REMOVED***
	return fs.NewProc(pid)
***REMOVED***

// NewProc returns a process for the given pid.
func (fs FS) NewProc(pid int) (Proc, error) ***REMOVED***
	if _, err := os.Stat(fs.Path(strconv.Itoa(pid))); err != nil ***REMOVED***
		return Proc***REMOVED******REMOVED***, err
	***REMOVED***
	return Proc***REMOVED***PID: pid, fs: fs***REMOVED***, nil
***REMOVED***

// AllProcs returns a list of all currently available processes.
func (fs FS) AllProcs() (Procs, error) ***REMOVED***
	d, err := os.Open(fs.Path())
	if err != nil ***REMOVED***
		return Procs***REMOVED******REMOVED***, err
	***REMOVED***
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil ***REMOVED***
		return Procs***REMOVED******REMOVED***, fmt.Errorf("could not read %s: %s", d.Name(), err)
	***REMOVED***

	p := Procs***REMOVED******REMOVED***
	for _, n := range names ***REMOVED***
		pid, err := strconv.ParseInt(n, 10, 64)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		p = append(p, Proc***REMOVED***PID: int(pid), fs: fs***REMOVED***)
	***REMOVED***

	return p, nil
***REMOVED***

// CmdLine returns the command line of a process.
func (p Proc) CmdLine() ([]string, error) ***REMOVED***
	f, err := os.Open(p.path("cmdline"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(data) < 1 ***REMOVED***
		return []string***REMOVED******REMOVED***, nil
	***REMOVED***

	return strings.Split(string(data[:len(data)-1]), string(byte(0))), nil
***REMOVED***

// Comm returns the command name of a process.
func (p Proc) Comm() (string, error) ***REMOVED***
	f, err := os.Open(p.path("comm"))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return strings.TrimSpace(string(data)), nil
***REMOVED***

// Executable returns the absolute path of the executable command of a process.
func (p Proc) Executable() (string, error) ***REMOVED***
	exe, err := os.Readlink(p.path("exe"))
	if os.IsNotExist(err) ***REMOVED***
		return "", nil
	***REMOVED***

	return exe, err
***REMOVED***

// FileDescriptors returns the currently open file descriptors of a process.
func (p Proc) FileDescriptors() ([]uintptr, error) ***REMOVED***
	names, err := p.fileDescriptors()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	fds := make([]uintptr, len(names))
	for i, n := range names ***REMOVED***
		fd, err := strconv.ParseInt(n, 10, 32)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("could not parse fd %s: %s", n, err)
		***REMOVED***
		fds[i] = uintptr(fd)
	***REMOVED***

	return fds, nil
***REMOVED***

// FileDescriptorTargets returns the targets of all file descriptors of a process.
// If a file descriptor is not a symlink to a file (like a socket), that value will be the empty string.
func (p Proc) FileDescriptorTargets() ([]string, error) ***REMOVED***
	names, err := p.fileDescriptors()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	targets := make([]string, len(names))

	for i, name := range names ***REMOVED***
		target, err := os.Readlink(p.path("fd", name))
		if err == nil ***REMOVED***
			targets[i] = target
		***REMOVED***
	***REMOVED***

	return targets, nil
***REMOVED***

// FileDescriptorsLen returns the number of currently open file descriptors of
// a process.
func (p Proc) FileDescriptorsLen() (int, error) ***REMOVED***
	fds, err := p.fileDescriptors()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return len(fds), nil
***REMOVED***

func (p Proc) fileDescriptors() ([]string, error) ***REMOVED***
	d, err := os.Open(p.path("fd"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("could not read %s: %s", d.Name(), err)
	***REMOVED***

	return names, nil
***REMOVED***

func (p Proc) path(pa ...string) string ***REMOVED***
	return p.fs.Path(append([]string***REMOVED***strconv.Itoa(p.PID)***REMOVED***, pa...)...)
***REMOVED***
