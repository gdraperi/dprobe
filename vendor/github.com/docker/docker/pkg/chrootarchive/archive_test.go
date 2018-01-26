package chrootarchive

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/pkg/system"
)

func init() ***REMOVED***
	reexec.Init()
***REMOVED***

var chrootArchiver = NewArchiver(nil)

func TarUntar(src, dst string) error ***REMOVED***
	return chrootArchiver.TarUntar(src, dst)
***REMOVED***

func CopyFileWithTar(src, dst string) (err error) ***REMOVED***
	return chrootArchiver.CopyFileWithTar(src, dst)
***REMOVED***

func UntarPath(src, dst string) error ***REMOVED***
	return chrootArchiver.UntarPath(src, dst)
***REMOVED***

func CopyWithTar(src, dst string) error ***REMOVED***
	return chrootArchiver.CopyWithTar(src, dst)
***REMOVED***

func TestChrootTarUntar(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootTarUntar")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	src := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(src, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(src, "toto"), []byte("hello toto"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(src, "lolo"), []byte("hello lolo"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	stream, err := archive.Tar(src, archive.Uncompressed)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	dest := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(dest, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := Untar(stream, dest, &archive.TarOptions***REMOVED***ExcludePatterns: []string***REMOVED***"lolo"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// gh#10426: Verify the fix for having a huge excludes list (like on `docker load` with large # of
// local images)
func TestChrootUntarWithHugeExcludesList(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootUntarHugeExcludes")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	src := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(src, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(src, "toto"), []byte("hello toto"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	stream, err := archive.Tar(src, archive.Uncompressed)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	dest := filepath.Join(tmpdir, "dest")
	if err := system.MkdirAll(dest, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	options := &archive.TarOptions***REMOVED******REMOVED***
	//65534 entries of 64-byte strings ~= 4MB of environment space which should overflow
	//on most systems when passed via environment or command line arguments
	excludes := make([]string, 65534)
	for i := 0; i < 65534; i++ ***REMOVED***
		excludes[i] = strings.Repeat(string(i), 64)
	***REMOVED***
	options.ExcludePatterns = excludes
	if err := Untar(stream, dest, options); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestChrootUntarEmptyArchive(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootUntarEmptyArchive")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	if err := Untar(nil, tmpdir, nil); err == nil ***REMOVED***
		t.Fatal("expected error on empty archive")
	***REMOVED***
***REMOVED***

func prepareSourceDirectory(numberOfFiles int, targetPath string, makeSymLinks bool) (int, error) ***REMOVED***
	fileData := []byte("fooo")
	for n := 0; n < numberOfFiles; n++ ***REMOVED***
		fileName := fmt.Sprintf("file-%d", n)
		if err := ioutil.WriteFile(filepath.Join(targetPath, fileName), fileData, 0700); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if makeSymLinks ***REMOVED***
			if err := os.Symlink(filepath.Join(targetPath, fileName), filepath.Join(targetPath, fileName+"-link")); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	totalSize := numberOfFiles * len(fileData)
	return totalSize, nil
***REMOVED***

func getHash(filename string) (uint32, error) ***REMOVED***
	stream, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	hash := crc32.NewIEEE()
	hash.Write(stream)
	return hash.Sum32(), nil
***REMOVED***

func compareDirectories(src string, dest string) error ***REMOVED***
	changes, err := archive.ChangesDirs(dest, src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(changes) > 0 ***REMOVED***
		return fmt.Errorf("Unexpected differences after untar: %v", changes)
	***REMOVED***
	return nil
***REMOVED***

func compareFiles(src string, dest string) error ***REMOVED***
	srcHash, err := getHash(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	destHash, err := getHash(dest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if srcHash != destHash ***REMOVED***
		return fmt.Errorf("%s is different from %s", src, dest)
	***REMOVED***
	return nil
***REMOVED***

func TestChrootTarUntarWithSymlink(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootTarUntarWithSymlink")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	src := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(src, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := prepareSourceDirectory(10, src, false); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	dest := filepath.Join(tmpdir, "dest")
	if err := TarUntar(src, dest); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareDirectories(src, dest); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestChrootCopyWithTar(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootCopyWithTar")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	src := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(src, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := prepareSourceDirectory(10, src, true); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Copy directory
	dest := filepath.Join(tmpdir, "dest")
	if err := CopyWithTar(src, dest); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareDirectories(src, dest); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Copy file
	srcfile := filepath.Join(src, "file-1")
	dest = filepath.Join(tmpdir, "destFile")
	destfile := filepath.Join(dest, "file-1")
	if err := CopyWithTar(srcfile, destfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareFiles(srcfile, destfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Copy symbolic link
	srcLinkfile := filepath.Join(src, "file-1-link")
	dest = filepath.Join(tmpdir, "destSymlink")
	destLinkfile := filepath.Join(dest, "file-1-link")
	if err := CopyWithTar(srcLinkfile, destLinkfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareFiles(srcLinkfile, destLinkfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestChrootCopyFileWithTar(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootCopyFileWithTar")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	src := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(src, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := prepareSourceDirectory(10, src, true); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Copy directory
	dest := filepath.Join(tmpdir, "dest")
	if err := CopyFileWithTar(src, dest); err == nil ***REMOVED***
		t.Fatal("Expected error on copying directory")
	***REMOVED***

	// Copy file
	srcfile := filepath.Join(src, "file-1")
	dest = filepath.Join(tmpdir, "destFile")
	destfile := filepath.Join(dest, "file-1")
	if err := CopyFileWithTar(srcfile, destfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareFiles(srcfile, destfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Copy symbolic link
	srcLinkfile := filepath.Join(src, "file-1-link")
	dest = filepath.Join(tmpdir, "destSymlink")
	destLinkfile := filepath.Join(dest, "file-1-link")
	if err := CopyFileWithTar(srcLinkfile, destLinkfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareFiles(srcLinkfile, destLinkfile); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestChrootUntarPath(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootUntarPath")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	src := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(src, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := prepareSourceDirectory(10, src, false); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	dest := filepath.Join(tmpdir, "dest")
	// Untar a directory
	if err := UntarPath(src, dest); err == nil ***REMOVED***
		t.Fatal("Expected error on untaring a directory")
	***REMOVED***

	// Untar a tar file
	stream, err := archive.Tar(src, archive.Uncompressed)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	tarfile := filepath.Join(tmpdir, "src.tar")
	if err := ioutil.WriteFile(tarfile, buf.Bytes(), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := UntarPath(tarfile, dest); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareDirectories(src, dest); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

type slowEmptyTarReader struct ***REMOVED***
	size      int
	offset    int
	chunkSize int
***REMOVED***

// Read is a slow reader of an empty tar (like the output of "tar c --files-from /dev/null")
func (s *slowEmptyTarReader) Read(p []byte) (int, error) ***REMOVED***
	time.Sleep(100 * time.Millisecond)
	count := s.chunkSize
	if len(p) < s.chunkSize ***REMOVED***
		count = len(p)
	***REMOVED***
	for i := 0; i < count; i++ ***REMOVED***
		p[i] = 0
	***REMOVED***
	s.offset += count
	if s.offset > s.size ***REMOVED***
		return count, io.EOF
	***REMOVED***
	return count, nil
***REMOVED***

func TestChrootUntarEmptyArchiveFromSlowReader(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootUntarEmptyArchiveFromSlowReader")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	dest := filepath.Join(tmpdir, "dest")
	if err := system.MkdirAll(dest, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	stream := &slowEmptyTarReader***REMOVED***size: 10240, chunkSize: 1024***REMOVED***
	if err := Untar(stream, dest, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestChrootApplyEmptyArchiveFromSlowReader(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootApplyEmptyArchiveFromSlowReader")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	dest := filepath.Join(tmpdir, "dest")
	if err := system.MkdirAll(dest, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	stream := &slowEmptyTarReader***REMOVED***size: 10240, chunkSize: 1024***REMOVED***
	if _, err := ApplyLayer(dest, stream); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestChrootApplyDotDotFile(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-TestChrootApplyDotDotFile")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	src := filepath.Join(tmpdir, "src")
	if err := system.MkdirAll(src, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(src, "..gitme"), []byte(""), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	stream, err := archive.Tar(src, archive.Uncompressed)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	dest := filepath.Join(tmpdir, "dest")
	if err := system.MkdirAll(dest, 0700, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := ApplyLayer(dest, stream); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
