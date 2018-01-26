package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tmp string

func init() ***REMOVED***
	tmp = "/tmp/"
	if runtime.GOOS == "windows" ***REMOVED***
		tmp = os.Getenv("TEMP") + `\`
	***REMOVED***
***REMOVED***

var defaultArchiver = NewDefaultArchiver()

func defaultTarUntar(src, dst string) error ***REMOVED***
	return defaultArchiver.TarUntar(src, dst)
***REMOVED***

func defaultUntarPath(src, dst string) error ***REMOVED***
	return defaultArchiver.UntarPath(src, dst)
***REMOVED***

func defaultCopyFileWithTar(src, dst string) (err error) ***REMOVED***
	return defaultArchiver.CopyFileWithTar(src, dst)
***REMOVED***

func defaultCopyWithTar(src, dst string) error ***REMOVED***
	return defaultArchiver.CopyWithTar(src, dst)
***REMOVED***

func TestIsArchivePathDir(t *testing.T) ***REMOVED***
	cmd := exec.Command("sh", "-c", "mkdir -p /tmp/archivedir")
	output, err := cmd.CombinedOutput()
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create an archive file for test : %s.", output)
	***REMOVED***
	if IsArchivePath(tmp + "archivedir") ***REMOVED***
		t.Fatalf("Incorrectly recognised directory as an archive")
	***REMOVED***
***REMOVED***

func TestIsArchivePathInvalidFile(t *testing.T) ***REMOVED***
	cmd := exec.Command("sh", "-c", "dd if=/dev/zero bs=1024 count=1 of=/tmp/archive && gzip --stdout /tmp/archive > /tmp/archive.gz")
	output, err := cmd.CombinedOutput()
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create an archive file for test : %s.", output)
	***REMOVED***
	if IsArchivePath(tmp + "archive") ***REMOVED***
		t.Fatalf("Incorrectly recognised invalid tar path as archive")
	***REMOVED***
	if IsArchivePath(tmp + "archive.gz") ***REMOVED***
		t.Fatalf("Incorrectly recognised invalid compressed tar path as archive")
	***REMOVED***
***REMOVED***

func TestIsArchivePathTar(t *testing.T) ***REMOVED***
	whichTar := "tar"
	cmdStr := fmt.Sprintf("touch /tmp/archivedata && %s -cf /tmp/archive /tmp/archivedata && gzip --stdout /tmp/archive > /tmp/archive.gz", whichTar)
	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create an archive file for test : %s.", output)
	***REMOVED***
	if !IsArchivePath(tmp + "/archive") ***REMOVED***
		t.Fatalf("Did not recognise valid tar path as archive")
	***REMOVED***
	if !IsArchivePath(tmp + "archive.gz") ***REMOVED***
		t.Fatalf("Did not recognise valid compressed tar path as archive")
	***REMOVED***
***REMOVED***

func testDecompressStream(t *testing.T, ext, compressCommand string) io.Reader ***REMOVED***
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("touch /tmp/archive && %s /tmp/archive", compressCommand))
	output, err := cmd.CombinedOutput()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create an archive file for test : %s.", output)
	***REMOVED***
	filename := "archive." + ext
	archive, err := os.Open(tmp + filename)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to open file %s: %v", filename, err)
	***REMOVED***
	defer archive.Close()

	r, err := DecompressStream(archive)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to decompress %s: %v", filename, err)
	***REMOVED***
	if _, err = ioutil.ReadAll(r); err != nil ***REMOVED***
		t.Fatalf("Failed to read the decompressed stream: %v ", err)
	***REMOVED***
	if err = r.Close(); err != nil ***REMOVED***
		t.Fatalf("Failed to close the decompressed stream: %v ", err)
	***REMOVED***

	return r
***REMOVED***

func TestDecompressStreamGzip(t *testing.T) ***REMOVED***
	testDecompressStream(t, "gz", "gzip -f")
***REMOVED***

func TestDecompressStreamBzip2(t *testing.T) ***REMOVED***
	testDecompressStream(t, "bz2", "bzip2 -f")
***REMOVED***

func TestDecompressStreamXz(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Xz not present in msys2")
	***REMOVED***
	testDecompressStream(t, "xz", "xz -f")
***REMOVED***

func TestCompressStreamXzUnsupported(t *testing.T) ***REMOVED***
	dest, err := os.Create(tmp + "dest")
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create the destination file")
	***REMOVED***
	defer dest.Close()

	_, err = CompressStream(dest, Xz)
	if err == nil ***REMOVED***
		t.Fatalf("Should fail as xz is unsupported for compression format.")
	***REMOVED***
***REMOVED***

func TestCompressStreamBzip2Unsupported(t *testing.T) ***REMOVED***
	dest, err := os.Create(tmp + "dest")
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create the destination file")
	***REMOVED***
	defer dest.Close()

	_, err = CompressStream(dest, Xz)
	if err == nil ***REMOVED***
		t.Fatalf("Should fail as xz is unsupported for compression format.")
	***REMOVED***
***REMOVED***

func TestCompressStreamInvalid(t *testing.T) ***REMOVED***
	dest, err := os.Create(tmp + "dest")
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create the destination file")
	***REMOVED***
	defer dest.Close()

	_, err = CompressStream(dest, -1)
	if err == nil ***REMOVED***
		t.Fatalf("Should fail as xz is unsupported for compression format.")
	***REMOVED***
***REMOVED***

func TestExtensionInvalid(t *testing.T) ***REMOVED***
	compression := Compression(-1)
	output := compression.Extension()
	if output != "" ***REMOVED***
		t.Fatalf("The extension of an invalid compression should be an empty string.")
	***REMOVED***
***REMOVED***

func TestExtensionUncompressed(t *testing.T) ***REMOVED***
	compression := Uncompressed
	output := compression.Extension()
	if output != "tar" ***REMOVED***
		t.Fatalf("The extension of an uncompressed archive should be 'tar'.")
	***REMOVED***
***REMOVED***
func TestExtensionBzip2(t *testing.T) ***REMOVED***
	compression := Bzip2
	output := compression.Extension()
	if output != "tar.bz2" ***REMOVED***
		t.Fatalf("The extension of a bzip2 archive should be 'tar.bz2'")
	***REMOVED***
***REMOVED***
func TestExtensionGzip(t *testing.T) ***REMOVED***
	compression := Gzip
	output := compression.Extension()
	if output != "tar.gz" ***REMOVED***
		t.Fatalf("The extension of a bzip2 archive should be 'tar.gz'")
	***REMOVED***
***REMOVED***
func TestExtensionXz(t *testing.T) ***REMOVED***
	compression := Xz
	output := compression.Extension()
	if output != "tar.xz" ***REMOVED***
		t.Fatalf("The extension of a bzip2 archive should be 'tar.xz'")
	***REMOVED***
***REMOVED***

func TestCmdStreamLargeStderr(t *testing.T) ***REMOVED***
	cmd := exec.Command("sh", "-c", "dd if=/dev/zero bs=1k count=1000 of=/dev/stderr; echo hello")
	out, err := cmdStream(cmd, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to start command: %s", err)
	***REMOVED***
	errCh := make(chan error)
	go func() ***REMOVED***
		_, err := io.Copy(ioutil.Discard, out)
		errCh <- err
	***REMOVED***()
	select ***REMOVED***
	case err := <-errCh:
		if err != nil ***REMOVED***
			t.Fatalf("Command should not have failed (err=%.100s...)", err)
		***REMOVED***
	case <-time.After(5 * time.Second):
		t.Fatalf("Command did not complete in 5 seconds; probable deadlock")
	***REMOVED***
***REMOVED***

func TestCmdStreamBad(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing in CI but not locally
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows CI machines")
	***REMOVED***
	badCmd := exec.Command("sh", "-c", "echo hello; echo >&2 error couldn\\'t reverse the phase pulser; exit 1")
	out, err := cmdStream(badCmd, nil)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to start command: %s", err)
	***REMOVED***
	if output, err := ioutil.ReadAll(out); err == nil ***REMOVED***
		t.Fatalf("Command should have failed")
	***REMOVED*** else if err.Error() != "exit status 1: error couldn't reverse the phase pulser\n" ***REMOVED***
		t.Fatalf("Wrong error value (%s)", err)
	***REMOVED*** else if s := string(output); s != "hello\n" ***REMOVED***
		t.Fatalf("Command output should be '%s', not '%s'", "hello\\n", output)
	***REMOVED***
***REMOVED***

func TestCmdStreamGood(t *testing.T) ***REMOVED***
	cmd := exec.Command("sh", "-c", "echo hello; exit 0")
	out, err := cmdStream(cmd, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if output, err := ioutil.ReadAll(out); err != nil ***REMOVED***
		t.Fatalf("Command should not have failed (err=%s)", err)
	***REMOVED*** else if s := string(output); s != "hello\n" ***REMOVED***
		t.Fatalf("Command output should be '%s', not '%s'", "hello\\n", output)
	***REMOVED***
***REMOVED***

func TestUntarPathWithInvalidDest(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-archive-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempFolder)
	invalidDestFolder := filepath.Join(tempFolder, "invalidDest")
	// Create a src file
	srcFile := filepath.Join(tempFolder, "src")
	tarFile := filepath.Join(tempFolder, "src.tar")
	os.Create(srcFile)
	os.Create(invalidDestFolder) // being a file (not dir) should cause an error

	// Translate back to Unix semantics as next exec.Command is run under sh
	srcFileU := srcFile
	tarFileU := tarFile
	if runtime.GOOS == "windows" ***REMOVED***
		tarFileU = "/tmp/" + filepath.Base(filepath.Dir(tarFile)) + "/src.tar"
		srcFileU = "/tmp/" + filepath.Base(filepath.Dir(srcFile)) + "/src"
	***REMOVED***

	cmd := exec.Command("sh", "-c", "tar cf "+tarFileU+" "+srcFileU)
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	err = defaultUntarPath(tarFile, invalidDestFolder)
	if err == nil ***REMOVED***
		t.Fatalf("UntarPath with invalid destination path should throw an error.")
	***REMOVED***
***REMOVED***

func TestUntarPathWithInvalidSrc(t *testing.T) ***REMOVED***
	dest, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create the destination file")
	***REMOVED***
	defer os.RemoveAll(dest)
	err = defaultUntarPath("/invalid/path", dest)
	if err == nil ***REMOVED***
		t.Fatalf("UntarPath with invalid src path should throw an error.")
	***REMOVED***
***REMOVED***

func TestUntarPath(t *testing.T) ***REMOVED***
	tmpFolder, err := ioutil.TempDir("", "docker-archive-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpFolder)
	srcFile := filepath.Join(tmpFolder, "src")
	tarFile := filepath.Join(tmpFolder, "src.tar")
	os.Create(filepath.Join(tmpFolder, "src"))

	destFolder := filepath.Join(tmpFolder, "dest")
	err = os.MkdirAll(destFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create the destination file")
	***REMOVED***

	// Translate back to Unix semantics as next exec.Command is run under sh
	srcFileU := srcFile
	tarFileU := tarFile
	if runtime.GOOS == "windows" ***REMOVED***
		tarFileU = "/tmp/" + filepath.Base(filepath.Dir(tarFile)) + "/src.tar"
		srcFileU = "/tmp/" + filepath.Base(filepath.Dir(srcFile)) + "/src"
	***REMOVED***
	cmd := exec.Command("sh", "-c", "tar cf "+tarFileU+" "+srcFileU)
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	err = defaultUntarPath(tarFile, destFolder)
	if err != nil ***REMOVED***
		t.Fatalf("UntarPath shouldn't throw an error, %s.", err)
	***REMOVED***
	expectedFile := filepath.Join(destFolder, srcFileU)
	_, err = os.Stat(expectedFile)
	if err != nil ***REMOVED***
		t.Fatalf("Destination folder should contain the source file but did not.")
	***REMOVED***
***REMOVED***

// Do the same test as above but with the destination as file, it should fail
func TestUntarPathWithDestinationFile(t *testing.T) ***REMOVED***
	tmpFolder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpFolder)
	srcFile := filepath.Join(tmpFolder, "src")
	tarFile := filepath.Join(tmpFolder, "src.tar")
	os.Create(filepath.Join(tmpFolder, "src"))

	// Translate back to Unix semantics as next exec.Command is run under sh
	srcFileU := srcFile
	tarFileU := tarFile
	if runtime.GOOS == "windows" ***REMOVED***
		tarFileU = "/tmp/" + filepath.Base(filepath.Dir(tarFile)) + "/src.tar"
		srcFileU = "/tmp/" + filepath.Base(filepath.Dir(srcFile)) + "/src"
	***REMOVED***
	cmd := exec.Command("sh", "-c", "tar cf "+tarFileU+" "+srcFileU)
	_, err = cmd.CombinedOutput()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	destFile := filepath.Join(tmpFolder, "dest")
	_, err = os.Create(destFile)
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create the destination file")
	***REMOVED***
	err = defaultUntarPath(tarFile, destFile)
	if err == nil ***REMOVED***
		t.Fatalf("UntarPath should throw an error if the destination if a file")
	***REMOVED***
***REMOVED***

// Do the same test as above but with the destination folder already exists
// and the destination file is a directory
// It's working, see https://github.com/docker/docker/issues/10040
func TestUntarPathWithDestinationSrcFileAsFolder(t *testing.T) ***REMOVED***
	tmpFolder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpFolder)
	srcFile := filepath.Join(tmpFolder, "src")
	tarFile := filepath.Join(tmpFolder, "src.tar")
	os.Create(srcFile)

	// Translate back to Unix semantics as next exec.Command is run under sh
	srcFileU := srcFile
	tarFileU := tarFile
	if runtime.GOOS == "windows" ***REMOVED***
		tarFileU = "/tmp/" + filepath.Base(filepath.Dir(tarFile)) + "/src.tar"
		srcFileU = "/tmp/" + filepath.Base(filepath.Dir(srcFile)) + "/src"
	***REMOVED***

	cmd := exec.Command("sh", "-c", "tar cf "+tarFileU+" "+srcFileU)
	_, err = cmd.CombinedOutput()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	destFolder := filepath.Join(tmpFolder, "dest")
	err = os.MkdirAll(destFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatalf("Fail to create the destination folder")
	***REMOVED***
	// Let's create a folder that will has the same path as the extracted file (from tar)
	destSrcFileAsFolder := filepath.Join(destFolder, srcFileU)
	err = os.MkdirAll(destSrcFileAsFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = defaultUntarPath(tarFile, destFolder)
	if err != nil ***REMOVED***
		t.Fatalf("UntarPath should throw not throw an error if the extracted file already exists and is a folder")
	***REMOVED***
***REMOVED***

func TestCopyWithTarInvalidSrc(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(nil)
	***REMOVED***
	destFolder := filepath.Join(tempFolder, "dest")
	invalidSrc := filepath.Join(tempFolder, "doesnotexists")
	err = os.MkdirAll(destFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = defaultCopyWithTar(invalidSrc, destFolder)
	if err == nil ***REMOVED***
		t.Fatalf("archiver.CopyWithTar with invalid src path should throw an error.")
	***REMOVED***
***REMOVED***

func TestCopyWithTarInexistentDestWillCreateIt(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(nil)
	***REMOVED***
	srcFolder := filepath.Join(tempFolder, "src")
	inexistentDestFolder := filepath.Join(tempFolder, "doesnotexists")
	err = os.MkdirAll(srcFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = defaultCopyWithTar(srcFolder, inexistentDestFolder)
	if err != nil ***REMOVED***
		t.Fatalf("CopyWithTar with an inexistent folder shouldn't fail.")
	***REMOVED***
	_, err = os.Stat(inexistentDestFolder)
	if err != nil ***REMOVED***
		t.Fatalf("CopyWithTar with an inexistent folder should create it.")
	***REMOVED***
***REMOVED***

// Test CopyWithTar with a file as src
func TestCopyWithTarSrcFile(t *testing.T) ***REMOVED***
	folder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(folder)
	dest := filepath.Join(folder, "dest")
	srcFolder := filepath.Join(folder, "src")
	src := filepath.Join(folder, filepath.Join("src", "src"))
	err = os.MkdirAll(srcFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = os.MkdirAll(dest, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ioutil.WriteFile(src, []byte("content"), 0777)
	err = defaultCopyWithTar(src, dest)
	if err != nil ***REMOVED***
		t.Fatalf("archiver.CopyWithTar shouldn't throw an error, %s.", err)
	***REMOVED***
	_, err = os.Stat(dest)
	// FIXME Check the content
	if err != nil ***REMOVED***
		t.Fatalf("Destination file should be the same as the source.")
	***REMOVED***
***REMOVED***

// Test CopyWithTar with a folder as src
func TestCopyWithTarSrcFolder(t *testing.T) ***REMOVED***
	folder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(folder)
	dest := filepath.Join(folder, "dest")
	src := filepath.Join(folder, filepath.Join("src", "folder"))
	err = os.MkdirAll(src, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = os.MkdirAll(dest, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ioutil.WriteFile(filepath.Join(src, "file"), []byte("content"), 0777)
	err = defaultCopyWithTar(src, dest)
	if err != nil ***REMOVED***
		t.Fatalf("archiver.CopyWithTar shouldn't throw an error, %s.", err)
	***REMOVED***
	_, err = os.Stat(dest)
	// FIXME Check the content (the file inside)
	if err != nil ***REMOVED***
		t.Fatalf("Destination folder should contain the source file but did not.")
	***REMOVED***
***REMOVED***

func TestCopyFileWithTarInvalidSrc(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tempFolder)
	destFolder := filepath.Join(tempFolder, "dest")
	err = os.MkdirAll(destFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	invalidFile := filepath.Join(tempFolder, "doesnotexists")
	err = defaultCopyFileWithTar(invalidFile, destFolder)
	if err == nil ***REMOVED***
		t.Fatalf("archiver.CopyWithTar with invalid src path should throw an error.")
	***REMOVED***
***REMOVED***

func TestCopyFileWithTarInexistentDestWillCreateIt(t *testing.T) ***REMOVED***
	tempFolder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(nil)
	***REMOVED***
	defer os.RemoveAll(tempFolder)
	srcFile := filepath.Join(tempFolder, "src")
	inexistentDestFolder := filepath.Join(tempFolder, "doesnotexists")
	_, err = os.Create(srcFile)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = defaultCopyFileWithTar(srcFile, inexistentDestFolder)
	if err != nil ***REMOVED***
		t.Fatalf("CopyWithTar with an inexistent folder shouldn't fail.")
	***REMOVED***
	_, err = os.Stat(inexistentDestFolder)
	if err != nil ***REMOVED***
		t.Fatalf("CopyWithTar with an inexistent folder should create it.")
	***REMOVED***
	// FIXME Test the src file and content
***REMOVED***

func TestCopyFileWithTarSrcFolder(t *testing.T) ***REMOVED***
	folder, err := ioutil.TempDir("", "docker-archive-copyfilewithtar-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(folder)
	dest := filepath.Join(folder, "dest")
	src := filepath.Join(folder, "srcfolder")
	err = os.MkdirAll(src, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = os.MkdirAll(dest, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = defaultCopyFileWithTar(src, dest)
	if err == nil ***REMOVED***
		t.Fatalf("CopyFileWithTar should throw an error with a folder.")
	***REMOVED***
***REMOVED***

func TestCopyFileWithTarSrcFile(t *testing.T) ***REMOVED***
	folder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(folder)
	dest := filepath.Join(folder, "dest")
	srcFolder := filepath.Join(folder, "src")
	src := filepath.Join(folder, filepath.Join("src", "src"))
	err = os.MkdirAll(srcFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = os.MkdirAll(dest, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ioutil.WriteFile(src, []byte("content"), 0777)
	err = defaultCopyWithTar(src, dest+"/")
	if err != nil ***REMOVED***
		t.Fatalf("archiver.CopyFileWithTar shouldn't throw an error, %s.", err)
	***REMOVED***
	_, err = os.Stat(dest)
	if err != nil ***REMOVED***
		t.Fatalf("Destination folder should contain the source file but did not.")
	***REMOVED***
***REMOVED***

func TestTarFiles(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out how to port this test.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	// try without hardlinks
	if err := checkNoChanges(1000, false); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// try with hardlinks
	if err := checkNoChanges(1000, true); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func checkNoChanges(fileNum int, hardlinks bool) error ***REMOVED***
	srcDir, err := ioutil.TempDir("", "docker-test-srcDir")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(srcDir)

	destDir, err := ioutil.TempDir("", "docker-test-destDir")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(destDir)

	_, err = prepareUntarSourceDirectory(fileNum, srcDir, hardlinks)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = defaultTarUntar(srcDir, destDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	changes, err := ChangesDirs(destDir, srcDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(changes) > 0 ***REMOVED***
		return fmt.Errorf("with %d files and %v hardlinks: expected 0 changes, got %d", fileNum, hardlinks, len(changes))
	***REMOVED***
	return nil
***REMOVED***

func tarUntar(t *testing.T, origin string, options *TarOptions) ([]Change, error) ***REMOVED***
	archive, err := TarWithOptions(origin, options)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer archive.Close()

	buf := make([]byte, 10)
	if _, err := archive.Read(buf); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	wrap := io.MultiReader(bytes.NewReader(buf), archive)

	detectedCompression := DetectCompression(buf)
	compression := options.Compression
	if detectedCompression.Extension() != compression.Extension() ***REMOVED***
		return nil, fmt.Errorf("Wrong compression detected. Actual compression: %s, found %s", compression.Extension(), detectedCompression.Extension())
	***REMOVED***

	tmp, err := ioutil.TempDir("", "docker-test-untar")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer os.RemoveAll(tmp)
	if err := Untar(wrap, tmp, nil); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := os.Stat(tmp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return ChangesDirs(origin, tmp)
***REMOVED***

func TestTarUntar(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out how to fix this test.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-untar-origin")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(origin)
	if err := ioutil.WriteFile(filepath.Join(origin, "1"), []byte("hello world"), 0700); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(origin, "2"), []byte("welcome!"), 0700); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(origin, "3"), []byte("will be ignored"), 0700); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for _, c := range []Compression***REMOVED***
		Uncompressed,
		Gzip,
	***REMOVED*** ***REMOVED***
		changes, err := tarUntar(t, origin, &TarOptions***REMOVED***
			Compression:     c,
			ExcludePatterns: []string***REMOVED***"3"***REMOVED***,
		***REMOVED***)

		if err != nil ***REMOVED***
			t.Fatalf("Error tar/untar for compression %s: %s", c.Extension(), err)
		***REMOVED***

		if len(changes) != 1 || changes[0].Path != "/3" ***REMOVED***
			t.Fatalf("Unexpected differences after tarUntar: %v", changes)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTarWithOptionsChownOptsAlwaysOverridesIdPair(t *testing.T) ***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-tar-chown-opt")
	require.NoError(t, err)

	defer os.RemoveAll(origin)
	filePath := filepath.Join(origin, "1")
	err = ioutil.WriteFile(filePath, []byte("hello world"), 0700)
	require.NoError(t, err)

	idMaps := []idtools.IDMap***REMOVED***
		0: ***REMOVED***
			ContainerID: 0,
			HostID:      0,
			Size:        65536,
		***REMOVED***,
		1: ***REMOVED***
			ContainerID: 0,
			HostID:      100000,
			Size:        65536,
		***REMOVED***,
	***REMOVED***

	cases := []struct ***REMOVED***
		opts        *TarOptions
		expectedUID int
		expectedGID int
	***REMOVED******REMOVED***
		***REMOVED***&TarOptions***REMOVED***ChownOpts: &idtools.IDPair***REMOVED***UID: 1337, GID: 42***REMOVED******REMOVED***, 1337, 42***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***ChownOpts: &idtools.IDPair***REMOVED***UID: 100001, GID: 100001***REMOVED***, UIDMaps: idMaps, GIDMaps: idMaps***REMOVED***, 100001, 100001***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***ChownOpts: &idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***, NoLchown: false***REMOVED***, 0, 0***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***ChownOpts: &idtools.IDPair***REMOVED***UID: 1, GID: 1***REMOVED***, NoLchown: true***REMOVED***, 1, 1***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***ChownOpts: &idtools.IDPair***REMOVED***UID: 1000, GID: 1000***REMOVED***, NoLchown: true***REMOVED***, 1000, 1000***REMOVED***,
	***REMOVED***
	for _, testCase := range cases ***REMOVED***
		reader, err := TarWithOptions(filePath, testCase.opts)
		require.NoError(t, err)
		tr := tar.NewReader(reader)
		defer reader.Close()
		for ***REMOVED***
			hdr, err := tr.Next()
			if err == io.EOF ***REMOVED***
				// end of tar archive
				break
			***REMOVED***
			require.NoError(t, err)
			assert.Equal(t, hdr.Uid, testCase.expectedUID, "Uid equals expected value")
			assert.Equal(t, hdr.Gid, testCase.expectedGID, "Gid equals expected value")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTarWithOptions(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out how to fix this test.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-untar-origin")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := ioutil.TempDir(origin, "folder"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(origin)
	if err := ioutil.WriteFile(filepath.Join(origin, "1"), []byte("hello world"), 0700); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(origin, "2"), []byte("welcome!"), 0700); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	cases := []struct ***REMOVED***
		opts       *TarOptions
		numChanges int
	***REMOVED******REMOVED***
		***REMOVED***&TarOptions***REMOVED***IncludeFiles: []string***REMOVED***"1"***REMOVED******REMOVED***, 2***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***ExcludePatterns: []string***REMOVED***"2"***REMOVED******REMOVED***, 1***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***ExcludePatterns: []string***REMOVED***"1", "folder*"***REMOVED******REMOVED***, 2***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***IncludeFiles: []string***REMOVED***"1", "1"***REMOVED******REMOVED***, 2***REMOVED***,
		***REMOVED***&TarOptions***REMOVED***IncludeFiles: []string***REMOVED***"1"***REMOVED***, RebaseNames: map[string]string***REMOVED***"1": "test"***REMOVED******REMOVED***, 4***REMOVED***,
	***REMOVED***
	for _, testCase := range cases ***REMOVED***
		changes, err := tarUntar(t, origin, testCase.opts)
		if err != nil ***REMOVED***
			t.Fatalf("Error tar/untar when testing inclusion/exclusion: %s", err)
		***REMOVED***
		if len(changes) != testCase.numChanges ***REMOVED***
			t.Errorf("Expected %d changes, got %d for %+v:",
				testCase.numChanges, len(changes), testCase.opts)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Some tar archives such as http://haproxy.1wt.eu/download/1.5/src/devel/haproxy-1.5-dev21.tar.gz
// use PAX Global Extended Headers.
// Failing prevents the archives from being uncompressed during ADD
func TestTypeXGlobalHeaderDoesNotFail(t *testing.T) ***REMOVED***
	hdr := tar.Header***REMOVED***Typeflag: tar.TypeXGlobalHeader***REMOVED***
	tmpDir, err := ioutil.TempDir("", "docker-test-archive-pax-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpDir)
	err = createTarFile(filepath.Join(tmpDir, "pax_global_header"), tmpDir, &hdr, nil, true, nil, false)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Some tar have both GNU specific (huge uid) and Ustar specific (long name) things.
// Not supposed to happen (should use PAX instead of Ustar for long name) but it does and it should still work.
func TestUntarUstarGnuConflict(t *testing.T) ***REMOVED***
	f, err := os.Open("testdata/broken.tar")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer f.Close()

	found := false
	tr := tar.NewReader(f)
	// Iterate through the files in the archive.
	for ***REMOVED***
		hdr, err := tr.Next()
		if err == io.EOF ***REMOVED***
			// end of tar archive
			break
		***REMOVED***
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if hdr.Name == "root/.cpanm/work/1395823785.24209/Plack-1.0030/blib/man3/Plack::Middleware::LighttpdScriptNameFix.3pm" ***REMOVED***
			found = true
			break
		***REMOVED***
	***REMOVED***
	if !found ***REMOVED***
		t.Fatalf("%s not found in the archive", "root/.cpanm/work/1395823785.24209/Plack-1.0030/blib/man3/Plack::Middleware::LighttpdScriptNameFix.3pm")
	***REMOVED***
***REMOVED***

func prepareUntarSourceDirectory(numberOfFiles int, targetPath string, makeLinks bool) (int, error) ***REMOVED***
	fileData := []byte("fooo")
	for n := 0; n < numberOfFiles; n++ ***REMOVED***
		fileName := fmt.Sprintf("file-%d", n)
		if err := ioutil.WriteFile(filepath.Join(targetPath, fileName), fileData, 0700); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if makeLinks ***REMOVED***
			if err := os.Link(filepath.Join(targetPath, fileName), filepath.Join(targetPath, fileName+"-link")); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	totalSize := numberOfFiles * len(fileData)
	return totalSize, nil
***REMOVED***

func BenchmarkTarUntar(b *testing.B) ***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-untar-origin")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	tempDir, err := ioutil.TempDir("", "docker-test-untar-destination")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	target := filepath.Join(tempDir, "dest")
	n, err := prepareUntarSourceDirectory(100, origin, false)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(origin)
	defer os.RemoveAll(tempDir)

	b.ResetTimer()
	b.SetBytes(int64(n))
	for n := 0; n < b.N; n++ ***REMOVED***
		err := defaultTarUntar(origin, target)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		os.RemoveAll(target)
	***REMOVED***
***REMOVED***

func BenchmarkTarUntarWithLinks(b *testing.B) ***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-untar-origin")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	tempDir, err := ioutil.TempDir("", "docker-test-untar-destination")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	target := filepath.Join(tempDir, "dest")
	n, err := prepareUntarSourceDirectory(100, origin, true)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(origin)
	defer os.RemoveAll(tempDir)

	b.ResetTimer()
	b.SetBytes(int64(n))
	for n := 0; n < b.N; n++ ***REMOVED***
		err := defaultTarUntar(origin, target)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		os.RemoveAll(target)
	***REMOVED***
***REMOVED***

func TestUntarInvalidFilenames(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out how to fix this test.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Passes but hits breakoutError: platform and architecture is not supported")
	***REMOVED***
	for i, headers := range [][]*tar.Header***REMOVED***
		***REMOVED***
			***REMOVED***
				Name:     "../victim/dotdot",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			***REMOVED***
				// Note the leading slash
				Name:     "/../victim/slash-dotdot",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := testBreakout("untar", "docker-TestUntarInvalidFilenames", headers); err != nil ***REMOVED***
			t.Fatalf("i=%d. %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUntarHardlinkToSymlink(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("hardlinks on Windows")
	***REMOVED***
	for i, headers := range [][]*tar.Header***REMOVED***
		***REMOVED***
			***REMOVED***
				Name:     "symlink1",
				Typeflag: tar.TypeSymlink,
				Linkname: "regfile",
				Mode:     0644,
			***REMOVED***,
			***REMOVED***
				Name:     "symlink2",
				Typeflag: tar.TypeLink,
				Linkname: "symlink1",
				Mode:     0644,
			***REMOVED***,
			***REMOVED***
				Name:     "regfile",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := testBreakout("untar", "docker-TestUntarHardlinkToSymlink", headers); err != nil ***REMOVED***
			t.Fatalf("i=%d. %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUntarInvalidHardlink(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("hardlinks on Windows")
	***REMOVED***
	for i, headers := range [][]*tar.Header***REMOVED***
		***REMOVED*** // try reading victim/hello (../)
			***REMOVED***
				Name:     "dotdot",
				Typeflag: tar.TypeLink,
				Linkname: "../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (/../)
			***REMOVED***
				Name:     "slash-dotdot",
				Typeflag: tar.TypeLink,
				// Note the leading slash
				Linkname: "/../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try writing victim/file
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim/file",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (hardlink, symlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "symlink",
				Typeflag: tar.TypeSymlink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // Try reading victim/hello (hardlink, hardlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "hardlink",
				Typeflag: tar.TypeLink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // Try removing victim directory (hardlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := testBreakout("untar", "docker-TestUntarInvalidHardlink", headers); err != nil ***REMOVED***
			t.Fatalf("i=%d. %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUntarInvalidSymlink(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("hardlinks on Windows")
	***REMOVED***
	for i, headers := range [][]*tar.Header***REMOVED***
		***REMOVED*** // try reading victim/hello (../)
			***REMOVED***
				Name:     "dotdot",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (/../)
			***REMOVED***
				Name:     "slash-dotdot",
				Typeflag: tar.TypeSymlink,
				// Note the leading slash
				Linkname: "/../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try writing victim/file
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim/file",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (symlink, symlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "symlink",
				Typeflag: tar.TypeSymlink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (symlink, hardlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "hardlink",
				Typeflag: tar.TypeLink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try removing victim directory (symlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try writing to victim/newdir/newfile with a symlink in the path
			***REMOVED***
				// this header needs to be before the next one, or else there is an error
				Name:     "dir/loophole",
				Typeflag: tar.TypeSymlink,
				Linkname: "../../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "dir/loophole/newdir/newfile",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := testBreakout("untar", "docker-TestUntarInvalidSymlink", headers); err != nil ***REMOVED***
			t.Fatalf("i=%d. %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTempArchiveCloseMultipleTimes(t *testing.T) ***REMOVED***
	reader := ioutil.NopCloser(strings.NewReader("hello"))
	tempArchive, err := NewTempArchive(reader, "")
	require.NoError(t, err)
	buf := make([]byte, 10)
	n, err := tempArchive.Read(buf)
	require.NoError(t, err)
	if n != 5 ***REMOVED***
		t.Fatalf("Expected to read 5 bytes. Read %d instead", n)
	***REMOVED***
	for i := 0; i < 3; i++ ***REMOVED***
		if err = tempArchive.Close(); err != nil ***REMOVED***
			t.Fatalf("i=%d. Unexpected error closing temp archive: %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReplaceFileTarWrapper(t *testing.T) ***REMOVED***
	filesInArchive := 20
	testcases := []struct ***REMOVED***
		doc       string
		filename  string
		modifier  TarModifierFunc
		expected  string
		fileCount int
	***REMOVED******REMOVED***
		***REMOVED***
			doc:       "Modifier creates a new file",
			filename:  "newfile",
			modifier:  createModifier(t),
			expected:  "the new content",
			fileCount: filesInArchive + 1,
		***REMOVED***,
		***REMOVED***
			doc:       "Modifier replaces a file",
			filename:  "file-2",
			modifier:  createOrReplaceModifier,
			expected:  "the new content",
			fileCount: filesInArchive,
		***REMOVED***,
		***REMOVED***
			doc:       "Modifier replaces the last file",
			filename:  fmt.Sprintf("file-%d", filesInArchive-1),
			modifier:  createOrReplaceModifier,
			expected:  "the new content",
			fileCount: filesInArchive,
		***REMOVED***,
		***REMOVED***
			doc:       "Modifier appends to a file",
			filename:  "file-3",
			modifier:  appendModifier,
			expected:  "fooo\nnext line",
			fileCount: filesInArchive,
		***REMOVED***,
	***REMOVED***

	for _, testcase := range testcases ***REMOVED***
		sourceArchive, cleanup := buildSourceArchive(t, filesInArchive)
		defer cleanup()

		resultArchive := ReplaceFileTarWrapper(
			sourceArchive,
			map[string]TarModifierFunc***REMOVED***testcase.filename: testcase.modifier***REMOVED***)

		actual := readFileFromArchive(t, resultArchive, testcase.filename, testcase.fileCount, testcase.doc)
		assert.Equal(t, testcase.expected, actual, testcase.doc)
	***REMOVED***
***REMOVED***

// TestPrefixHeaderReadable tests that files that could be created with the
// version of this package that was built with <=go17 are still readable.
func TestPrefixHeaderReadable(t *testing.T) ***REMOVED***
	// https://gist.github.com/stevvooe/e2a790ad4e97425896206c0816e1a882#file-out-go
	var testFile = []byte("\x1f\x8b\x08\x08\x44\x21\x68\x59\x00\x03\x74\x2e\x74\x61\x72\x00\x4b\xcb\xcf\x67\xa0\x35\x30\x80\x00\x86\x06\x10\x47\x01\xc1\x37\x40\x00\x54\xb6\xb1\xa1\xa9\x99\x09\x48\x25\x1d\x40\x69\x71\x49\x62\x91\x02\xe5\x76\xa1\x79\x84\x21\x91\xd6\x80\x72\xaf\x8f\x82\x51\x30\x0a\x46\x36\x00\x00\xf0\x1c\x1e\x95\x00\x06\x00\x00")

	tmpDir, err := ioutil.TempDir("", "prefix-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	err = Untar(bytes.NewReader(testFile), tmpDir, nil)
	require.NoError(t, err)

	baseName := "foo"
	pth := strings.Repeat("a", 100-len(baseName)) + "/" + baseName

	_, err = os.Lstat(filepath.Join(tmpDir, pth))
	require.NoError(t, err)
***REMOVED***

func buildSourceArchive(t *testing.T, numberOfFiles int) (io.ReadCloser, func()) ***REMOVED***
	srcDir, err := ioutil.TempDir("", "docker-test-srcDir")
	require.NoError(t, err)

	_, err = prepareUntarSourceDirectory(numberOfFiles, srcDir, false)
	require.NoError(t, err)

	sourceArchive, err := TarWithOptions(srcDir, &TarOptions***REMOVED******REMOVED***)
	require.NoError(t, err)
	return sourceArchive, func() ***REMOVED***
		os.RemoveAll(srcDir)
		sourceArchive.Close()
	***REMOVED***
***REMOVED***

func createOrReplaceModifier(path string, header *tar.Header, content io.Reader) (*tar.Header, []byte, error) ***REMOVED***
	return &tar.Header***REMOVED***
		Mode:     0600,
		Typeflag: tar.TypeReg,
	***REMOVED***, []byte("the new content"), nil
***REMOVED***

func createModifier(t *testing.T) TarModifierFunc ***REMOVED***
	return func(path string, header *tar.Header, content io.Reader) (*tar.Header, []byte, error) ***REMOVED***
		assert.Nil(t, content)
		return createOrReplaceModifier(path, header, content)
	***REMOVED***
***REMOVED***

func appendModifier(path string, header *tar.Header, content io.Reader) (*tar.Header, []byte, error) ***REMOVED***
	buffer := bytes.Buffer***REMOVED******REMOVED***
	if content != nil ***REMOVED***
		if _, err := buffer.ReadFrom(content); err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***
	buffer.WriteString("\nnext line")
	return &tar.Header***REMOVED***Mode: 0600, Typeflag: tar.TypeReg***REMOVED***, buffer.Bytes(), nil
***REMOVED***

func readFileFromArchive(t *testing.T, archive io.ReadCloser, name string, expectedCount int, doc string) string ***REMOVED***
	destDir, err := ioutil.TempDir("", "docker-test-destDir")
	require.NoError(t, err)
	defer os.RemoveAll(destDir)

	err = Untar(archive, destDir, nil)
	require.NoError(t, err)

	files, _ := ioutil.ReadDir(destDir)
	assert.Len(t, files, expectedCount, doc)

	content, err := ioutil.ReadFile(filepath.Join(destDir, name))
	assert.NoError(t, err)
	return string(content)
***REMOVED***

func TestDisablePigz(t *testing.T) ***REMOVED***
	_, err := exec.LookPath("unpigz")
	if err != nil ***REMOVED***
		t.Log("Test will not check full path when Pigz not installed")
	***REMOVED***

	os.Setenv("MOBY_DISABLE_PIGZ", "true")
	defer os.Unsetenv("MOBY_DISABLE_PIGZ")

	r := testDecompressStream(t, "gz", "gzip -f")
	// For the bufio pool
	outsideReaderCloserWrapper := r.(*ioutils.ReadCloserWrapper)
	// For the context canceller
	contextReaderCloserWrapper := outsideReaderCloserWrapper.Reader.(*ioutils.ReadCloserWrapper)

	assert.IsType(t, &gzip.Reader***REMOVED******REMOVED***, contextReaderCloserWrapper.Reader)
***REMOVED***

func TestPigz(t *testing.T) ***REMOVED***
	r := testDecompressStream(t, "gz", "gzip -f")
	// For the bufio pool
	outsideReaderCloserWrapper := r.(*ioutils.ReadCloserWrapper)
	// For the context canceller
	contextReaderCloserWrapper := outsideReaderCloserWrapper.Reader.(*ioutils.ReadCloserWrapper)

	_, err := exec.LookPath("unpigz")
	if err == nil ***REMOVED***
		t.Log("Tested whether Pigz is used, as it installed")
		assert.IsType(t, &io.PipeReader***REMOVED******REMOVED***, contextReaderCloserWrapper.Reader)
	***REMOVED*** else ***REMOVED***
		t.Log("Tested whether Pigz is not used, as it not installed")
		assert.IsType(t, &gzip.Reader***REMOVED******REMOVED***, contextReaderCloserWrapper.Reader)
	***REMOVED***
***REMOVED***
