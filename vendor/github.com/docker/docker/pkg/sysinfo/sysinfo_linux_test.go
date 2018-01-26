package sysinfo

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestReadProcBool(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "test-sysinfo-proc")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	procFile := filepath.Join(tmpDir, "read-proc-bool")
	err = ioutil.WriteFile(procFile, []byte("1"), 0644)
	require.NoError(t, err)

	if !readProcBool(procFile) ***REMOVED***
		t.Fatal("expected proc bool to be true, got false")
	***REMOVED***

	if err := ioutil.WriteFile(procFile, []byte("0"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if readProcBool(procFile) ***REMOVED***
		t.Fatal("expected proc bool to be false, got true")
	***REMOVED***

	if readProcBool(path.Join(tmpDir, "no-exist")) ***REMOVED***
		t.Fatal("should be false for non-existent entry")
	***REMOVED***

***REMOVED***

func TestCgroupEnabled(t *testing.T) ***REMOVED***
	cgroupDir, err := ioutil.TempDir("", "cgroup-test")
	require.NoError(t, err)
	defer os.RemoveAll(cgroupDir)

	if cgroupEnabled(cgroupDir, "test") ***REMOVED***
		t.Fatal("cgroupEnabled should be false")
	***REMOVED***

	err = ioutil.WriteFile(path.Join(cgroupDir, "test"), []byte***REMOVED******REMOVED***, 0644)
	require.NoError(t, err)

	if !cgroupEnabled(cgroupDir, "test") ***REMOVED***
		t.Fatal("cgroupEnabled should be true")
	***REMOVED***
***REMOVED***

func TestNew(t *testing.T) ***REMOVED***
	sysInfo := New(false)
	require.NotNil(t, sysInfo)
	checkSysInfo(t, sysInfo)

	sysInfo = New(true)
	require.NotNil(t, sysInfo)
	checkSysInfo(t, sysInfo)
***REMOVED***

func checkSysInfo(t *testing.T, sysInfo *SysInfo) ***REMOVED***
	// Check if Seccomp is supported, via CONFIG_SECCOMP.then sysInfo.Seccomp must be TRUE , else FALSE
	if err := unix.Prctl(unix.PR_GET_SECCOMP, 0, 0, 0, 0); err != unix.EINVAL ***REMOVED***
		// Make sure the kernel has CONFIG_SECCOMP_FILTER.
		if err := unix.Prctl(unix.PR_SET_SECCOMP, unix.SECCOMP_MODE_FILTER, 0, 0, 0); err != unix.EINVAL ***REMOVED***
			require.True(t, sysInfo.Seccomp)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		require.False(t, sysInfo.Seccomp)
	***REMOVED***
***REMOVED***

func TestNewAppArmorEnabled(t *testing.T) ***REMOVED***
	// Check if AppArmor is supported. then it must be TRUE , else FALSE
	if _, err := os.Stat("/sys/kernel/security/apparmor"); err != nil ***REMOVED***
		t.Skip("App Armor Must be Enabled")
	***REMOVED***

	sysInfo := New(true)
	require.True(t, sysInfo.AppArmor)
***REMOVED***

func TestNewAppArmorDisabled(t *testing.T) ***REMOVED***
	// Check if AppArmor is supported. then it must be TRUE , else FALSE
	if _, err := os.Stat("/sys/kernel/security/apparmor"); !os.IsNotExist(err) ***REMOVED***
		t.Skip("App Armor Must be Disabled")
	***REMOVED***

	sysInfo := New(true)
	require.False(t, sysInfo.AppArmor)
***REMOVED***

func TestNumCPU(t *testing.T) ***REMOVED***
	cpuNumbers := NumCPU()
	if cpuNumbers <= 0 ***REMOVED***
		t.Fatal("CPU returned must be greater than zero")
	***REMOVED***
***REMOVED***
