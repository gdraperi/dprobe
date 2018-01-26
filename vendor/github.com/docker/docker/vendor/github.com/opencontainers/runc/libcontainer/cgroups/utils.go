// +build linux

package cgroups

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/go-units"
)

const (
	cgroupNamePrefix = "name="
	CgroupProcesses  = "cgroup.procs"
)

// https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt
func FindCgroupMountpoint(subsystem string) (string, error) ***REMOVED***
	mnt, _, err := FindCgroupMountpointAndRoot(subsystem)
	return mnt, err
***REMOVED***

func FindCgroupMountpointAndRoot(subsystem string) (string, string, error) ***REMOVED***
	// We are not using mount.GetMounts() because it's super-inefficient,
	// parsing it directly sped up x10 times because of not using Sscanf.
	// It was one of two major performance drawbacks in container start.
	if !isSubsystemAvailable(subsystem) ***REMOVED***
		return "", "", NewNotFoundError(subsystem)
	***REMOVED***
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() ***REMOVED***
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") ***REMOVED***
			if opt == subsystem ***REMOVED***
				return fields[4], fields[3], nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if err := scanner.Err(); err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	return "", "", NewNotFoundError(subsystem)
***REMOVED***

func isSubsystemAvailable(subsystem string) bool ***REMOVED***
	cgroups, err := ParseCgroupFile("/proc/self/cgroup")
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	_, avail := cgroups[subsystem]
	return avail
***REMOVED***

func GetClosestMountpointAncestor(dir, mountinfo string) string ***REMOVED***
	deepestMountPoint := ""
	for _, mountInfoEntry := range strings.Split(mountinfo, "\n") ***REMOVED***
		mountInfoParts := strings.Fields(mountInfoEntry)
		if len(mountInfoParts) < 5 ***REMOVED***
			continue
		***REMOVED***
		mountPoint := mountInfoParts[4]
		if strings.HasPrefix(mountPoint, deepestMountPoint) && strings.HasPrefix(dir, mountPoint) ***REMOVED***
			deepestMountPoint = mountPoint
		***REMOVED***
	***REMOVED***
	return deepestMountPoint
***REMOVED***

func FindCgroupMountpointDir() (string, error) ***REMOVED***
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() ***REMOVED***
		text := scanner.Text()
		fields := strings.Split(text, " ")
		// Safe as mountinfo encodes mountpoints with spaces as \040.
		index := strings.Index(text, " - ")
		postSeparatorFields := strings.Fields(text[index+3:])
		numPostFields := len(postSeparatorFields)

		// This is an error as we can't detect if the mount is for "cgroup"
		if numPostFields == 0 ***REMOVED***
			return "", fmt.Errorf("Found no fields post '-' in %q", text)
		***REMOVED***

		if postSeparatorFields[0] == "cgroup" ***REMOVED***
			// Check that the mount is properly formated.
			if numPostFields < 3 ***REMOVED***
				return "", fmt.Errorf("Error found less than 3 fields post '-' in %q", text)
			***REMOVED***

			return filepath.Dir(fields[4]), nil
		***REMOVED***
	***REMOVED***
	if err := scanner.Err(); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return "", NewNotFoundError("cgroup")
***REMOVED***

type Mount struct ***REMOVED***
	Mountpoint string
	Root       string
	Subsystems []string
***REMOVED***

func (m Mount) GetOwnCgroup(cgroups map[string]string) (string, error) ***REMOVED***
	if len(m.Subsystems) == 0 ***REMOVED***
		return "", fmt.Errorf("no subsystem for mount")
	***REMOVED***

	return getControllerPath(m.Subsystems[0], cgroups)
***REMOVED***

func getCgroupMountsHelper(ss map[string]bool, mi io.Reader, all bool) ([]Mount, error) ***REMOVED***
	res := make([]Mount, 0, len(ss))
	scanner := bufio.NewScanner(mi)
	numFound := 0
	for scanner.Scan() && numFound < len(ss) ***REMOVED***
		txt := scanner.Text()
		sepIdx := strings.Index(txt, " - ")
		if sepIdx == -1 ***REMOVED***
			return nil, fmt.Errorf("invalid mountinfo format")
		***REMOVED***
		if txt[sepIdx+3:sepIdx+10] == "cgroup2" || txt[sepIdx+3:sepIdx+9] != "cgroup" ***REMOVED***
			continue
		***REMOVED***
		fields := strings.Split(txt, " ")
		m := Mount***REMOVED***
			Mountpoint: fields[4],
			Root:       fields[3],
		***REMOVED***
		for _, opt := range strings.Split(fields[len(fields)-1], ",") ***REMOVED***
			if !ss[opt] ***REMOVED***
				continue
			***REMOVED***
			if strings.HasPrefix(opt, cgroupNamePrefix) ***REMOVED***
				m.Subsystems = append(m.Subsystems, opt[len(cgroupNamePrefix):])
			***REMOVED*** else ***REMOVED***
				m.Subsystems = append(m.Subsystems, opt)
			***REMOVED***
			if !all ***REMOVED***
				numFound++
			***REMOVED***
		***REMOVED***
		res = append(res, m)
	***REMOVED***
	if err := scanner.Err(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return res, nil
***REMOVED***

// GetCgroupMounts returns the mounts for the cgroup subsystems.
// all indicates whether to return just the first instance or all the mounts.
func GetCgroupMounts(all bool) ([]Mount, error) ***REMOVED***
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	allSubsystems, err := ParseCgroupFile("/proc/self/cgroup")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	allMap := make(map[string]bool)
	for s := range allSubsystems ***REMOVED***
		allMap[s] = true
	***REMOVED***
	return getCgroupMountsHelper(allMap, f, all)
***REMOVED***

// GetAllSubsystems returns all the cgroup subsystems supported by the kernel
func GetAllSubsystems() ([]string, error) ***REMOVED***
	f, err := os.Open("/proc/cgroups")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	subsystems := []string***REMOVED******REMOVED***

	s := bufio.NewScanner(f)
	for s.Scan() ***REMOVED***
		text := s.Text()
		if text[0] != '#' ***REMOVED***
			parts := strings.Fields(text)
			if len(parts) >= 4 && parts[3] != "0" ***REMOVED***
				subsystems = append(subsystems, parts[0])
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if err := s.Err(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return subsystems, nil
***REMOVED***

// GetOwnCgroup returns the relative path to the cgroup docker is running in.
func GetOwnCgroup(subsystem string) (string, error) ***REMOVED***
	cgroups, err := ParseCgroupFile("/proc/self/cgroup")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return getControllerPath(subsystem, cgroups)
***REMOVED***

func GetOwnCgroupPath(subsystem string) (string, error) ***REMOVED***
	cgroup, err := GetOwnCgroup(subsystem)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return getCgroupPathHelper(subsystem, cgroup)
***REMOVED***

func GetInitCgroup(subsystem string) (string, error) ***REMOVED***
	cgroups, err := ParseCgroupFile("/proc/1/cgroup")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return getControllerPath(subsystem, cgroups)
***REMOVED***

func GetInitCgroupPath(subsystem string) (string, error) ***REMOVED***
	cgroup, err := GetInitCgroup(subsystem)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return getCgroupPathHelper(subsystem, cgroup)
***REMOVED***

func getCgroupPathHelper(subsystem, cgroup string) (string, error) ***REMOVED***
	mnt, root, err := FindCgroupMountpointAndRoot(subsystem)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// This is needed for nested containers, because in /proc/self/cgroup we
	// see pathes from host, which don't exist in container.
	relCgroup, err := filepath.Rel(root, cgroup)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return filepath.Join(mnt, relCgroup), nil
***REMOVED***

func readProcsFile(dir string) ([]int, error) ***REMOVED***
	f, err := os.Open(filepath.Join(dir, CgroupProcesses))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	var (
		s   = bufio.NewScanner(f)
		out = []int***REMOVED******REMOVED***
	)

	for s.Scan() ***REMOVED***
		if t := s.Text(); t != "" ***REMOVED***
			pid, err := strconv.Atoi(t)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			out = append(out, pid)
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***

// ParseCgroupFile parses the given cgroup file, typically from
// /proc/<pid>/cgroup, into a map of subgroups to cgroup names.
func ParseCgroupFile(path string) (map[string]string, error) ***REMOVED***
	f, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	return parseCgroupFromReader(f)
***REMOVED***

// helper function for ParseCgroupFile to make testing easier
func parseCgroupFromReader(r io.Reader) (map[string]string, error) ***REMOVED***
	s := bufio.NewScanner(r)
	cgroups := make(map[string]string)

	for s.Scan() ***REMOVED***
		text := s.Text()
		// from cgroups(7):
		// /proc/[pid]/cgroup
		// ...
		// For each cgroup hierarchy ... there is one entry
		// containing three colon-separated fields of the form:
		//     hierarchy-ID:subsystem-list:cgroup-path
		parts := strings.SplitN(text, ":", 3)
		if len(parts) < 3 ***REMOVED***
			return nil, fmt.Errorf("invalid cgroup entry: must contain at least two colons: %v", text)
		***REMOVED***

		for _, subs := range strings.Split(parts[1], ",") ***REMOVED***
			cgroups[subs] = parts[2]
		***REMOVED***
	***REMOVED***
	if err := s.Err(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return cgroups, nil
***REMOVED***

func getControllerPath(subsystem string, cgroups map[string]string) (string, error) ***REMOVED***

	if p, ok := cgroups[subsystem]; ok ***REMOVED***
		return p, nil
	***REMOVED***

	if p, ok := cgroups[cgroupNamePrefix+subsystem]; ok ***REMOVED***
		return p, nil
	***REMOVED***

	return "", NewNotFoundError(subsystem)
***REMOVED***

func PathExists(path string) bool ***REMOVED***
	if _, err := os.Stat(path); err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func EnterPid(cgroupPaths map[string]string, pid int) error ***REMOVED***
	for _, path := range cgroupPaths ***REMOVED***
		if PathExists(path) ***REMOVED***
			if err := WriteCgroupProc(path, pid); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// RemovePaths iterates over the provided paths removing them.
// We trying to remove all paths five times with increasing delay between tries.
// If after all there are not removed cgroups - appropriate error will be
// returned.
func RemovePaths(paths map[string]string) (err error) ***REMOVED***
	delay := 10 * time.Millisecond
	for i := 0; i < 5; i++ ***REMOVED***
		if i != 0 ***REMOVED***
			time.Sleep(delay)
			delay *= 2
		***REMOVED***
		for s, p := range paths ***REMOVED***
			os.RemoveAll(p)
			// TODO: here probably should be logging
			_, err := os.Stat(p)
			// We need this strange way of checking cgroups existence because
			// RemoveAll almost always returns error, even on already removed
			// cgroups
			if os.IsNotExist(err) ***REMOVED***
				delete(paths, s)
			***REMOVED***
		***REMOVED***
		if len(paths) == 0 ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("Failed to remove paths: %v", paths)
***REMOVED***

func GetHugePageSize() ([]string, error) ***REMOVED***
	var pageSizes []string
	sizeList := []string***REMOVED***"B", "kB", "MB", "GB", "TB", "PB"***REMOVED***
	files, err := ioutil.ReadDir("/sys/kernel/mm/hugepages")
	if err != nil ***REMOVED***
		return pageSizes, err
	***REMOVED***
	for _, st := range files ***REMOVED***
		nameArray := strings.Split(st.Name(), "-")
		pageSize, err := units.RAMInBytes(nameArray[1])
		if err != nil ***REMOVED***
			return []string***REMOVED******REMOVED***, err
		***REMOVED***
		sizeString := units.CustomSize("%g%s", float64(pageSize), 1024.0, sizeList)
		pageSizes = append(pageSizes, sizeString)
	***REMOVED***

	return pageSizes, nil
***REMOVED***

// GetPids returns all pids, that were added to cgroup at path.
func GetPids(path string) ([]int, error) ***REMOVED***
	return readProcsFile(path)
***REMOVED***

// GetAllPids returns all pids, that were added to cgroup at path and to all its
// subcgroups.
func GetAllPids(path string) ([]int, error) ***REMOVED***
	var pids []int
	// collect pids from all sub-cgroups
	err := filepath.Walk(path, func(p string, info os.FileInfo, iErr error) error ***REMOVED***
		dir, file := filepath.Split(p)
		if file != CgroupProcesses ***REMOVED***
			return nil
		***REMOVED***
		if iErr != nil ***REMOVED***
			return iErr
		***REMOVED***
		cPids, err := readProcsFile(dir)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		pids = append(pids, cPids...)
		return nil
	***REMOVED***)
	return pids, err
***REMOVED***

// WriteCgroupProc writes the specified pid into the cgroup's cgroup.procs file
func WriteCgroupProc(dir string, pid int) error ***REMOVED***
	// Normally dir should not be empty, one case is that cgroup subsystem
	// is not mounted, we will get empty dir, and we want it fail here.
	if dir == "" ***REMOVED***
		return fmt.Errorf("no such directory for %s", CgroupProcesses)
	***REMOVED***

	// Dont attach any pid to the cgroup if -1 is specified as a pid
	if pid != -1 ***REMOVED***
		if err := ioutil.WriteFile(filepath.Join(dir, CgroupProcesses), []byte(strconv.Itoa(pid)), 0700); err != nil ***REMOVED***
			return fmt.Errorf("failed to write %v to %v: %v", pid, CgroupProcesses, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
