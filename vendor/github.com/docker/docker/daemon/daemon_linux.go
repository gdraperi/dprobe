package daemon

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/mount"
	"github.com/sirupsen/logrus"
)

// On Linux, plugins use a static path for storing execution state,
// instead of deriving path from daemon's exec-root. This is because
// plugin socket files are created here and they cannot exceed max
// path length of 108 bytes.
func getPluginExecRoot(root string) string ***REMOVED***
	return "/run/docker/plugins"
***REMOVED***

func (daemon *Daemon) cleanupMountsByID(id string) error ***REMOVED***
	logrus.Debugf("Cleaning up old mountid %s: start.", id)
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	return daemon.cleanupMountsFromReaderByID(f, id, mount.Unmount)
***REMOVED***

func (daemon *Daemon) cleanupMountsFromReaderByID(reader io.Reader, id string, unmount func(target string) error) error ***REMOVED***
	if daemon.root == "" ***REMOVED***
		return nil
	***REMOVED***
	var errors []string

	regexps := getCleanPatterns(id)
	sc := bufio.NewScanner(reader)
	for sc.Scan() ***REMOVED***
		if fields := strings.Fields(sc.Text()); len(fields) >= 4 ***REMOVED***
			if mnt := fields[4]; strings.HasPrefix(mnt, daemon.root) ***REMOVED***
				for _, p := range regexps ***REMOVED***
					if p.MatchString(mnt) ***REMOVED***
						if err := unmount(mnt); err != nil ***REMOVED***
							logrus.Error(err)
							errors = append(errors, err.Error())
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := sc.Err(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		return fmt.Errorf("Error cleaning up mounts:\n%v", strings.Join(errors, "\n"))
	***REMOVED***

	logrus.Debugf("Cleaning up old mountid %v: done.", id)
	return nil
***REMOVED***

// cleanupMounts umounts shm/mqueue mounts for old containers
func (daemon *Daemon) cleanupMounts() error ***REMOVED***
	return daemon.cleanupMountsByID("")
***REMOVED***

func getCleanPatterns(id string) (regexps []*regexp.Regexp) ***REMOVED***
	var patterns []string
	if id == "" ***REMOVED***
		id = "[0-9a-f]***REMOVED***64***REMOVED***"
		patterns = append(patterns, "containers/"+id+"/shm")
	***REMOVED***
	patterns = append(patterns, "aufs/mnt/"+id+"$", "overlay/"+id+"/merged$", "zfs/graph/"+id+"$")
	for _, p := range patterns ***REMOVED***
		r, err := regexp.Compile(p)
		if err == nil ***REMOVED***
			regexps = append(regexps, r)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func getRealPath(path string) (string, error) ***REMOVED***
	return fileutils.ReadSymlinkedDirectory(path)
***REMOVED***
