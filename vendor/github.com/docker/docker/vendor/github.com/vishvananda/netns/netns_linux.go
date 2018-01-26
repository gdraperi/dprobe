// +build linux

package netns

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	// These constants belong in the syscall library but have not been
	// added yet.
	CLONE_NEWUTS  = 0x04000000 /* New utsname group? */
	CLONE_NEWIPC  = 0x08000000 /* New ipcs */
	CLONE_NEWUSER = 0x10000000 /* New user namespace */
	CLONE_NEWPID  = 0x20000000 /* New pid namespace */
	CLONE_NEWNET  = 0x40000000 /* New network namespace */
	CLONE_IO      = 0x80000000 /* Get io context */
)

// Setns sets namespace using syscall. Note that this should be a method
// in syscall but it has not been added.
func Setns(ns NsHandle, nstype int) (err error) ***REMOVED***
	_, _, e1 := syscall.Syscall(SYS_SETNS, uintptr(ns), uintptr(nstype), 0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

// Set sets the current network namespace to the namespace represented
// by NsHandle.
func Set(ns NsHandle) (err error) ***REMOVED***
	return Setns(ns, CLONE_NEWNET)
***REMOVED***

// New creates a new network namespace and returns a handle to it.
func New() (ns NsHandle, err error) ***REMOVED***
	if err := syscall.Unshare(CLONE_NEWNET); err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return Get()
***REMOVED***

// Get gets a handle to the current threads network namespace.
func Get() (NsHandle, error) ***REMOVED***
	return GetFromThread(os.Getpid(), syscall.Gettid())
***REMOVED***

// GetFromPath gets a handle to a network namespace
// identified by the path
func GetFromPath(path string) (NsHandle, error) ***REMOVED***
	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return NsHandle(fd), nil
***REMOVED***

// GetFromName gets a handle to a named network namespace such as one
// created by `ip netns add`.
func GetFromName(name string) (NsHandle, error) ***REMOVED***
	return GetFromPath(fmt.Sprintf("/var/run/netns/%s", name))
***REMOVED***

// GetFromPid gets a handle to the network namespace of a given pid.
func GetFromPid(pid int) (NsHandle, error) ***REMOVED***
	return GetFromPath(fmt.Sprintf("/proc/%d/ns/net", pid))
***REMOVED***

// GetFromThread gets a handle to the network namespace of a given pid and tid.
func GetFromThread(pid, tid int) (NsHandle, error) ***REMOVED***
	return GetFromPath(fmt.Sprintf("/proc/%d/task/%d/ns/net", pid, tid))
***REMOVED***

// GetFromDocker gets a handle to the network namespace of a docker container.
// Id is prefixed matched against the running docker containers, so a short
// identifier can be used as long as it isn't ambiguous.
func GetFromDocker(id string) (NsHandle, error) ***REMOVED***
	pid, err := getPidForContainer(id)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return GetFromPid(pid)
***REMOVED***

// borrowed from docker/utils/utils.go
func findCgroupMountpoint(cgroupType string) (string, error) ***REMOVED***
	output, err := ioutil.ReadFile("/proc/mounts")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// /proc/mounts has 6 fields per line, one mount per line, e.g.
	// cgroup /sys/fs/cgroup/devices cgroup rw,relatime,devices 0 0
	for _, line := range strings.Split(string(output), "\n") ***REMOVED***
		parts := strings.Split(line, " ")
		if len(parts) == 6 && parts[2] == "cgroup" ***REMOVED***
			for _, opt := range strings.Split(parts[3], ",") ***REMOVED***
				if opt == cgroupType ***REMOVED***
					return parts[1], nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return "", fmt.Errorf("cgroup mountpoint not found for %s", cgroupType)
***REMOVED***

// Returns the relative path to the cgroup docker is running in.
// borrowed from docker/utils/utils.go
// modified to get the docker pid instead of using /proc/self
func getThisCgroup(cgroupType string) (string, error) ***REMOVED***
	dockerpid, err := ioutil.ReadFile("/var/run/docker.pid")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	result := strings.Split(string(dockerpid), "\n")
	if len(result) == 0 || len(result[0]) == 0 ***REMOVED***
		return "", fmt.Errorf("docker pid not found in /var/run/docker.pid")
	***REMOVED***
	pid, err := strconv.Atoi(result[0])

	output, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cgroup", pid))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	for _, line := range strings.Split(string(output), "\n") ***REMOVED***
		parts := strings.Split(line, ":")
		// any type used by docker should work
		if parts[1] == cgroupType ***REMOVED***
			return parts[2], nil
		***REMOVED***
	***REMOVED***
	return "", fmt.Errorf("cgroup '%s' not found in /proc/%d/cgroup", cgroupType, pid)
***REMOVED***

// Returns the first pid in a container.
// borrowed from docker/utils/utils.go
// modified to only return the first pid
// modified to glob with id
// modified to search for newer docker containers
func getPidForContainer(id string) (int, error) ***REMOVED***
	pid := 0

	// memory is chosen randomly, any cgroup used by docker works
	cgroupType := "memory"

	cgroupRoot, err := findCgroupMountpoint(cgroupType)
	if err != nil ***REMOVED***
		return pid, err
	***REMOVED***

	cgroupThis, err := getThisCgroup(cgroupType)
	if err != nil ***REMOVED***
		return pid, err
	***REMOVED***

	id += "*"

	attempts := []string***REMOVED***
		filepath.Join(cgroupRoot, cgroupThis, id, "tasks"),
		// With more recent lxc versions use, cgroup will be in lxc/
		filepath.Join(cgroupRoot, cgroupThis, "lxc", id, "tasks"),
		// With more recent dockee, cgroup will be in docker/
		filepath.Join(cgroupRoot, cgroupThis, "docker", id, "tasks"),
	***REMOVED***

	var filename string
	for _, attempt := range attempts ***REMOVED***
		filenames, _ := filepath.Glob(attempt)
		if len(filenames) > 1 ***REMOVED***
			return pid, fmt.Errorf("Ambiguous id supplied: %v", filenames)
		***REMOVED*** else if len(filenames) == 1 ***REMOVED***
			filename = filenames[0]
			break
		***REMOVED***
	***REMOVED***

	if filename == "" ***REMOVED***
		return pid, fmt.Errorf("Unable to find container: %v", id[:len(id)-1])
	***REMOVED***

	output, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return pid, err
	***REMOVED***

	result := strings.Split(string(output), "\n")
	if len(result) == 0 || len(result[0]) == 0 ***REMOVED***
		return pid, fmt.Errorf("No pid found for container")
	***REMOVED***

	pid, err = strconv.Atoi(result[0])
	if err != nil ***REMOVED***
		return pid, fmt.Errorf("Invalid pid '%s': %s", result[0], err)
	***REMOVED***

	return pid, nil
***REMOVED***
