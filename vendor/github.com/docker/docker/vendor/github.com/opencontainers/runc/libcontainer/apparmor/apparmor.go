// +build apparmor,linux

package apparmor

import (
	"fmt"
	"io/ioutil"
	"os"
)

// IsEnabled returns true if apparmor is enabled for the host.
func IsEnabled() bool ***REMOVED***
	if _, err := os.Stat("/sys/kernel/security/apparmor"); err == nil && os.Getenv("container") == "" ***REMOVED***
		if _, err = os.Stat("/sbin/apparmor_parser"); err == nil ***REMOVED***
			buf, err := ioutil.ReadFile("/sys/module/apparmor/parameters/enabled")
			return err == nil && len(buf) > 1 && buf[0] == 'Y'
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func setprocattr(attr, value string) error ***REMOVED***
	// Under AppArmor you can only change your own attr, so use /proc/self/
	// instead of /proc/<tid>/ like libapparmor does
	path := fmt.Sprintf("/proc/self/attr/%s", attr)

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s", value)
	return err
***REMOVED***

// changeOnExec reimplements aa_change_onexec from libapparmor in Go
func changeOnExec(name string) error ***REMOVED***
	value := "exec " + name
	if err := setprocattr("exec", value); err != nil ***REMOVED***
		return fmt.Errorf("apparmor failed to apply profile: %s", err)
	***REMOVED***
	return nil
***REMOVED***

// ApplyProfile will apply the profile with the specified name to the process after
// the next exec.
func ApplyProfile(name string) error ***REMOVED***
	if name == "" ***REMOVED***
		return nil
	***REMOVED***

	return changeOnExec(name)
***REMOVED***
