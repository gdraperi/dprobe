// +build linux

package quota

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// 10MB
const testQuotaSize = 10 * 1024 * 1024
const imageSize = 64 * 1024 * 1024

func TestBlockDev(t *testing.T) ***REMOVED***
	mkfs, err := exec.LookPath("mkfs.xfs")
	if err != nil ***REMOVED***
		t.Skip("mkfs.xfs not found in PATH")
	***REMOVED***

	// create a sparse image
	imageFile, err := ioutil.TempFile("", "xfs-image")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	imageFileName := imageFile.Name()
	defer os.Remove(imageFileName)
	if _, err = imageFile.Seek(imageSize-1, 0); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err = imageFile.Write([]byte***REMOVED***0***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err = imageFile.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// The reason for disabling these options is sometimes people run with a newer userspace
	// than kernelspace
	out, err := exec.Command(mkfs, "-m", "crc=0,finobt=0", imageFileName).CombinedOutput()
	if len(out) > 0 ***REMOVED***
		t.Log(string(out))
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	t.Run("testBlockDevQuotaDisabled", wrapMountTest(imageFileName, false, testBlockDevQuotaDisabled))
	t.Run("testBlockDevQuotaEnabled", wrapMountTest(imageFileName, true, testBlockDevQuotaEnabled))
	t.Run("testSmallerThanQuota", wrapMountTest(imageFileName, true, wrapQuotaTest(testSmallerThanQuota)))
	t.Run("testBiggerThanQuota", wrapMountTest(imageFileName, true, wrapQuotaTest(testBiggerThanQuota)))
	t.Run("testRetrieveQuota", wrapMountTest(imageFileName, true, wrapQuotaTest(testRetrieveQuota)))
***REMOVED***

func wrapMountTest(imageFileName string, enableQuota bool, testFunc func(t *testing.T, mountPoint, backingFsDev string)) func(*testing.T) ***REMOVED***
	return func(t *testing.T) ***REMOVED***
		mountOptions := "loop"

		if enableQuota ***REMOVED***
			mountOptions = mountOptions + ",prjquota"
		***REMOVED***

		mountPointDir := fs.NewDir(t, "xfs-mountPoint")
		defer mountPointDir.Remove()
		mountPoint := mountPointDir.Path()

		out, err := exec.Command("mount", "-o", mountOptions, imageFileName, mountPoint).CombinedOutput()
		if err != nil ***REMOVED***
			_, err := os.Stat("/proc/fs/xfs")
			if os.IsNotExist(err) ***REMOVED***
				t.Skip("no /proc/fs/xfs")
			***REMOVED***
		***REMOVED***

		require.NoError(t, err, "mount failed: %s", out)

		defer func() ***REMOVED***
			require.NoError(t, unix.Unmount(mountPoint, 0))
		***REMOVED***()

		backingFsDev, err := makeBackingFsDev(mountPoint)
		require.NoError(t, err)

		testFunc(t, mountPoint, backingFsDev)
	***REMOVED***
***REMOVED***

func testBlockDevQuotaDisabled(t *testing.T, mountPoint, backingFsDev string) ***REMOVED***
	hasSupport, err := hasQuotaSupport(backingFsDev)
	require.NoError(t, err)
	assert.False(t, hasSupport)
***REMOVED***

func testBlockDevQuotaEnabled(t *testing.T, mountPoint, backingFsDev string) ***REMOVED***
	hasSupport, err := hasQuotaSupport(backingFsDev)
	require.NoError(t, err)
	assert.True(t, hasSupport)
***REMOVED***

func wrapQuotaTest(testFunc func(t *testing.T, ctrl *Control, mountPoint, testDir, testSubDir string)) func(t *testing.T, mountPoint, backingFsDev string) ***REMOVED***
	return func(t *testing.T, mountPoint, backingFsDev string) ***REMOVED***
		testDir, err := ioutil.TempDir(mountPoint, "per-test")
		require.NoError(t, err)
		defer os.RemoveAll(testDir)

		ctrl, err := NewControl(testDir)
		require.NoError(t, err)

		testSubDir, err := ioutil.TempDir(testDir, "quota-test")
		require.NoError(t, err)
		testFunc(t, ctrl, mountPoint, testDir, testSubDir)
	***REMOVED***

***REMOVED***

func testSmallerThanQuota(t *testing.T, ctrl *Control, homeDir, testDir, testSubDir string) ***REMOVED***
	require.NoError(t, ctrl.SetQuota(testSubDir, Quota***REMOVED***testQuotaSize***REMOVED***))
	smallerThanQuotaFile := filepath.Join(testSubDir, "smaller-than-quota")
	require.NoError(t, ioutil.WriteFile(smallerThanQuotaFile, make([]byte, testQuotaSize/2), 0644))
	require.NoError(t, os.Remove(smallerThanQuotaFile))
***REMOVED***

func testBiggerThanQuota(t *testing.T, ctrl *Control, homeDir, testDir, testSubDir string) ***REMOVED***
	// Make sure the quota is being enforced
	// TODO: When we implement this under EXT4, we need to shed CAP_SYS_RESOURCE, otherwise
	// we're able to violate quota without issue
	require.NoError(t, ctrl.SetQuota(testSubDir, Quota***REMOVED***testQuotaSize***REMOVED***))

	biggerThanQuotaFile := filepath.Join(testSubDir, "bigger-than-quota")
	err := ioutil.WriteFile(biggerThanQuotaFile, make([]byte, testQuotaSize+1), 0644)
	require.Error(t, err)
	if err == io.ErrShortWrite ***REMOVED***
		require.NoError(t, os.Remove(biggerThanQuotaFile))
	***REMOVED***
***REMOVED***

func testRetrieveQuota(t *testing.T, ctrl *Control, homeDir, testDir, testSubDir string) ***REMOVED***
	// Validate that we can retrieve quota
	require.NoError(t, ctrl.SetQuota(testSubDir, Quota***REMOVED***testQuotaSize***REMOVED***))

	var q Quota
	require.NoError(t, ctrl.GetQuota(testSubDir, &q))
	assert.EqualValues(t, testQuotaSize, q.Size)
***REMOVED***
