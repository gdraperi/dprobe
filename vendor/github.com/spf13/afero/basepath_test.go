package afero

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBasePath(t *testing.T) ***REMOVED***
	baseFs := &MemMapFs***REMOVED******REMOVED***
	baseFs.MkdirAll("/base/path/tmp", 0777)
	bp := NewBasePathFs(baseFs, "/base/path")

	if _, err := bp.Create("/tmp/foo"); err != nil ***REMOVED***
		t.Errorf("Failed to set real path")
	***REMOVED***

	if fh, err := bp.Create("../tmp/bar"); err == nil ***REMOVED***
		t.Errorf("succeeded in creating %s ...", fh.Name())
	***REMOVED***
***REMOVED***

func TestBasePathRoot(t *testing.T) ***REMOVED***
	baseFs := &MemMapFs***REMOVED******REMOVED***
	baseFs.MkdirAll("/base/path/foo/baz", 0777)
	baseFs.MkdirAll("/base/path/boo/", 0777)
	bp := NewBasePathFs(baseFs, "/base/path")

	rd, err := ReadDir(bp, string(os.PathSeparator))

	if len(rd) != 2 ***REMOVED***
		t.Errorf("base path doesn't respect root")
	***REMOVED***

	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestRealPath(t *testing.T) ***REMOVED***
	fs := NewOsFs()
	baseDir, err := TempDir(fs, "", "base")
	if err != nil ***REMOVED***
		t.Fatal("error creating tempDir", err)
	***REMOVED***
	defer fs.RemoveAll(baseDir)
	anotherDir, err := TempDir(fs, "", "another")
	if err != nil ***REMOVED***
		t.Fatal("error creating tempDir", err)
	***REMOVED***
	defer fs.RemoveAll(anotherDir)

	bp := NewBasePathFs(fs, baseDir).(*BasePathFs)

	subDir := filepath.Join(baseDir, "s1")

	realPath, err := bp.RealPath("/s1")

	if err != nil ***REMOVED***
		t.Errorf("Got error %s", err)
	***REMOVED***

	if realPath != subDir ***REMOVED***
		t.Errorf("Expected \n%s got \n%s", subDir, realPath)
	***REMOVED***

	if runtime.GOOS == "windows" ***REMOVED***
		_, err = bp.RealPath(anotherDir)

		if err == nil ***REMOVED***
			t.Errorf("Expected error")
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		// on *nix we have no way of just looking at the path and tell that anotherDir
		// is not inside the base file system.
		// The user will receive an os.ErrNotExist later.
		surrealPath, err := bp.RealPath(anotherDir)

		if err != nil ***REMOVED***
			t.Errorf("Got error %s", err)
		***REMOVED***

		excpected := filepath.Join(baseDir, anotherDir)

		if surrealPath != excpected ***REMOVED***
			t.Errorf("Expected \n%s got \n%s", excpected, surrealPath)
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestNestedBasePaths(t *testing.T) ***REMOVED***
	type dirSpec struct ***REMOVED***
		Dir1, Dir2, Dir3 string
	***REMOVED***
	dirSpecs := []dirSpec***REMOVED***
		dirSpec***REMOVED***Dir1: "/", Dir2: "/", Dir3: "/"***REMOVED***,
		dirSpec***REMOVED***Dir1: "/", Dir2: "/path2", Dir3: "/"***REMOVED***,
		dirSpec***REMOVED***Dir1: "/path1/dir", Dir2: "/path2/dir/", Dir3: "/path3/dir"***REMOVED***,
		dirSpec***REMOVED***Dir1: "C:/path1", Dir2: "path2/dir", Dir3: "/path3/dir/"***REMOVED***,
	***REMOVED***

	for _, ds := range dirSpecs ***REMOVED***
		memFs := NewMemMapFs()
		level1Fs := NewBasePathFs(memFs, ds.Dir1)
		level2Fs := NewBasePathFs(level1Fs, ds.Dir2)
		level3Fs := NewBasePathFs(level2Fs, ds.Dir3)

		type spec struct ***REMOVED***
			BaseFs   Fs
			FileName string
		***REMOVED***
		specs := []spec***REMOVED***
			spec***REMOVED***BaseFs: level3Fs, FileName: "f.txt"***REMOVED***,
			spec***REMOVED***BaseFs: level2Fs, FileName: "f.txt"***REMOVED***,
			spec***REMOVED***BaseFs: level1Fs, FileName: "f.txt"***REMOVED***,
		***REMOVED***

		for _, s := range specs ***REMOVED***
			if err := s.BaseFs.MkdirAll(s.FileName, 0755); err != nil ***REMOVED***
				t.Errorf("Got error %s", err.Error())
			***REMOVED***
			if _, err := s.BaseFs.Stat(s.FileName); err != nil ***REMOVED***
				t.Errorf("Got error %s", err.Error())
			***REMOVED***

			if s.BaseFs == level3Fs ***REMOVED***
				pathToExist := filepath.Join(ds.Dir3, s.FileName)
				if _, err := level2Fs.Stat(pathToExist); err != nil ***REMOVED***
					t.Errorf("Got error %s (path %s)", err.Error(), pathToExist)
				***REMOVED***
			***REMOVED*** else if s.BaseFs == level2Fs ***REMOVED***
				pathToExist := filepath.Join(ds.Dir2, ds.Dir3, s.FileName)
				if _, err := level1Fs.Stat(pathToExist); err != nil ***REMOVED***
					t.Errorf("Got error %s (path %s)", err.Error(), pathToExist)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
