// +build !windows

package mount

import (
	"os"
	"path"
	"testing"
)

func TestMountOptionsParsing(t *testing.T) ***REMOVED***
	options := "noatime,ro,size=10k"

	flag, data := parseOptions(options)

	if data != "size=10k" ***REMOVED***
		t.Fatalf("Expected size=10 got %s", data)
	***REMOVED***

	expectedFlag := NOATIME | RDONLY

	if flag != expectedFlag ***REMOVED***
		t.Fatalf("Expected %d got %d", expectedFlag, flag)
	***REMOVED***
***REMOVED***

func TestMounted(t *testing.T) ***REMOVED***
	tmp := path.Join(os.TempDir(), "mount-tests")
	if err := os.MkdirAll(tmp, 0777); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)

	var (
		sourceDir  = path.Join(tmp, "source")
		targetDir  = path.Join(tmp, "target")
		sourcePath = path.Join(sourceDir, "file.txt")
		targetPath = path.Join(targetDir, "file.txt")
	)

	os.Mkdir(sourceDir, 0777)
	os.Mkdir(targetDir, 0777)

	f, err := os.Create(sourcePath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("hello")
	f.Close()

	f, err = os.Create(targetPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.Close()

	if err := Mount(sourceDir, targetDir, "none", "bind,rw"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer func() ***REMOVED***
		if err := Unmount(targetDir); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***()

	mounted, err := Mounted(targetDir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !mounted ***REMOVED***
		t.Fatalf("Expected %s to be mounted", targetDir)
	***REMOVED***
	if _, err := os.Stat(targetDir); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestMountReadonly(t *testing.T) ***REMOVED***
	tmp := path.Join(os.TempDir(), "mount-tests")
	if err := os.MkdirAll(tmp, 0777); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)

	var (
		sourceDir  = path.Join(tmp, "source")
		targetDir  = path.Join(tmp, "target")
		sourcePath = path.Join(sourceDir, "file.txt")
		targetPath = path.Join(targetDir, "file.txt")
	)

	os.Mkdir(sourceDir, 0777)
	os.Mkdir(targetDir, 0777)

	f, err := os.Create(sourcePath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("hello")
	f.Close()

	f, err = os.Create(targetPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.Close()

	if err := Mount(sourceDir, targetDir, "none", "bind,ro"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer func() ***REMOVED***
		if err := Unmount(targetDir); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***()

	f, err = os.OpenFile(targetPath, os.O_RDWR, 0777)
	if err == nil ***REMOVED***
		t.Fatal("Should not be able to open a ro file as rw")
	***REMOVED***
***REMOVED***

func TestGetMounts(t *testing.T) ***REMOVED***
	mounts, err := GetMounts()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	root := false
	for _, entry := range mounts ***REMOVED***
		if entry.Mountpoint == "/" ***REMOVED***
			root = true
		***REMOVED***
	***REMOVED***

	if !root ***REMOVED***
		t.Fatal("/ should be mounted at least")
	***REMOVED***
***REMOVED***

func TestMergeTmpfsOptions(t *testing.T) ***REMOVED***
	options := []string***REMOVED***"noatime", "ro", "size=10k", "defaults", "atime", "defaults", "rw", "rprivate", "size=1024k", "slave"***REMOVED***
	expected := []string***REMOVED***"atime", "rw", "size=1024k", "slave"***REMOVED***
	merged, err := MergeTmpfsOptions(options)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(expected) != len(merged) ***REMOVED***
		t.Fatalf("Expected %s got %s", expected, merged)
	***REMOVED***
	for index := range merged ***REMOVED***
		if merged[index] != expected[index] ***REMOVED***
			t.Fatalf("Expected %s for the %dth option, got %s", expected, index, merged)
		***REMOVED***
	***REMOVED***

	options = []string***REMOVED***"noatime", "ro", "size=10k", "atime", "rw", "rprivate", "size=1024k", "slave", "size"***REMOVED***
	_, err = MergeTmpfsOptions(options)
	if err == nil ***REMOVED***
		t.Fatal("Expected error got nil")
	***REMOVED***
***REMOVED***
