// +build linux

package fsutils

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"golang.org/x/sys/unix"
)

func testSupportsDType(t *testing.T, expected bool, mkfsCommand string, mkfsArg ...string) ***REMOVED***
	// check whether mkfs is installed
	if _, err := exec.LookPath(mkfsCommand); err != nil ***REMOVED***
		t.Skipf("%s not installed: %v", mkfsCommand, err)
	***REMOVED***

	// create a sparse image
	imageSize := int64(32 * 1024 * 1024)
	imageFile, err := ioutil.TempFile("", "fsutils-image")
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

	// create a mountpoint
	mountpoint, err := ioutil.TempDir("", "fsutils-mountpoint")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(mountpoint)

	// format the image
	args := append(mkfsArg, imageFileName)
	t.Logf("Executing `%s %v`", mkfsCommand, args)
	out, err := exec.Command(mkfsCommand, args...).CombinedOutput()
	if len(out) > 0 ***REMOVED***
		t.Log(string(out))
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// loopback-mount the image.
	// for ease of setting up loopback device, we use os/exec rather than unix.Mount
	out, err = exec.Command("mount", "-o", "loop", imageFileName, mountpoint).CombinedOutput()
	if len(out) > 0 ***REMOVED***
		t.Log(string(out))
	***REMOVED***
	if err != nil ***REMOVED***
		t.Skip("skipping the test because mount failed")
	***REMOVED***
	defer func() ***REMOVED***
		if err := unix.Unmount(mountpoint, 0); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***()

	// check whether it supports d_type
	result, err := SupportsDType(mountpoint)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	t.Logf("Supports d_type: %v", result)
	if result != expected ***REMOVED***
		t.Fatalf("expected %v, got %v", expected, result)
	***REMOVED***
***REMOVED***

func TestSupportsDTypeWithFType0XFS(t *testing.T) ***REMOVED***
	testSupportsDType(t, false, "mkfs.xfs", "-m", "crc=0", "-n", "ftype=0")
***REMOVED***

func TestSupportsDTypeWithFType1XFS(t *testing.T) ***REMOVED***
	testSupportsDType(t, true, "mkfs.xfs", "-m", "crc=0", "-n", "ftype=1")
***REMOVED***

func TestSupportsDTypeWithExt4(t *testing.T) ***REMOVED***
	testSupportsDType(t, true, "mkfs.ext4")
***REMOVED***
