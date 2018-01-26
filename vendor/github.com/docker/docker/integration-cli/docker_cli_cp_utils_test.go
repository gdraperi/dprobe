package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/pkg/archive"
	"github.com/go-check/check"
)

type fileType uint32

const (
	ftRegular fileType = iota
	ftDir
	ftSymlink
)

type fileData struct ***REMOVED***
	filetype fileType
	path     string
	contents string
	uid      int
	gid      int
	mode     int
***REMOVED***

func (fd fileData) creationCommand() string ***REMOVED***
	var command string

	switch fd.filetype ***REMOVED***
	case ftRegular:
		// Don't overwrite the file if it already exists!
		command = fmt.Sprintf("if [ ! -f %s ]; then echo %q > %s; fi", fd.path, fd.contents, fd.path)
	case ftDir:
		command = fmt.Sprintf("mkdir -p %s", fd.path)
	case ftSymlink:
		command = fmt.Sprintf("ln -fs %s %s", fd.contents, fd.path)
	***REMOVED***

	return command
***REMOVED***

func mkFilesCommand(fds []fileData) string ***REMOVED***
	commands := make([]string, len(fds))

	for i, fd := range fds ***REMOVED***
		commands[i] = fd.creationCommand()
	***REMOVED***

	return strings.Join(commands, " && ")
***REMOVED***

var defaultFileData = []fileData***REMOVED***
	***REMOVED***ftRegular, "file1", "file1", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "file2", "file2", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "file3", "file3", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "file4", "file4", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "file5", "file5", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "file6", "file6", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "file7", "file7", 0, 0, 0666***REMOVED***,
	***REMOVED***ftDir, "dir1", "", 0, 0, 0777***REMOVED***,
	***REMOVED***ftRegular, "dir1/file1-1", "file1-1", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "dir1/file1-2", "file1-2", 0, 0, 0666***REMOVED***,
	***REMOVED***ftDir, "dir2", "", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "dir2/file2-1", "file2-1", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "dir2/file2-2", "file2-2", 0, 0, 0666***REMOVED***,
	***REMOVED***ftDir, "dir3", "", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "dir3/file3-1", "file3-1", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "dir3/file3-2", "file3-2", 0, 0, 0666***REMOVED***,
	***REMOVED***ftDir, "dir4", "", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "dir4/file3-1", "file4-1", 0, 0, 0666***REMOVED***,
	***REMOVED***ftRegular, "dir4/file3-2", "file4-2", 0, 0, 0666***REMOVED***,
	***REMOVED***ftDir, "dir5", "", 0, 0, 0666***REMOVED***,
	***REMOVED***ftSymlink, "symlinkToFile1", "file1", 0, 0, 0666***REMOVED***,
	***REMOVED***ftSymlink, "symlinkToDir1", "dir1", 0, 0, 0666***REMOVED***,
	***REMOVED***ftSymlink, "brokenSymlinkToFileX", "fileX", 0, 0, 0666***REMOVED***,
	***REMOVED***ftSymlink, "brokenSymlinkToDirX", "dirX", 0, 0, 0666***REMOVED***,
	***REMOVED***ftSymlink, "symlinkToAbsDir", "/root", 0, 0, 0666***REMOVED***,
	***REMOVED***ftDir, "permdirtest", "", 2, 2, 0700***REMOVED***,
	***REMOVED***ftRegular, "permdirtest/permtest", "perm_test", 65534, 65534, 0400***REMOVED***,
***REMOVED***

func defaultMkContentCommand() string ***REMOVED***
	return mkFilesCommand(defaultFileData)
***REMOVED***

func makeTestContentInDir(c *check.C, dir string) ***REMOVED***
	for _, fd := range defaultFileData ***REMOVED***
		path := filepath.Join(dir, filepath.FromSlash(fd.path))
		switch fd.filetype ***REMOVED***
		case ftRegular:
			c.Assert(ioutil.WriteFile(path, []byte(fd.contents+"\n"), os.FileMode(fd.mode)), checker.IsNil)
		case ftDir:
			c.Assert(os.Mkdir(path, os.FileMode(fd.mode)), checker.IsNil)
		case ftSymlink:
			c.Assert(os.Symlink(fd.contents, path), checker.IsNil)
		***REMOVED***

		if fd.filetype != ftSymlink && runtime.GOOS != "windows" ***REMOVED***
			c.Assert(os.Chown(path, fd.uid, fd.gid), checker.IsNil)
		***REMOVED***
	***REMOVED***
***REMOVED***

type testContainerOptions struct ***REMOVED***
	addContent bool
	readOnly   bool
	volumes    []string
	workDir    string
	command    string
***REMOVED***

func makeTestContainer(c *check.C, options testContainerOptions) (containerID string) ***REMOVED***
	if options.addContent ***REMOVED***
		mkContentCmd := defaultMkContentCommand()
		if options.command == "" ***REMOVED***
			options.command = mkContentCmd
		***REMOVED*** else ***REMOVED***
			options.command = fmt.Sprintf("%s && %s", defaultMkContentCommand(), options.command)
		***REMOVED***
	***REMOVED***

	if options.command == "" ***REMOVED***
		options.command = "#(nop)"
	***REMOVED***

	args := []string***REMOVED***"run", "-d"***REMOVED***

	for _, volume := range options.volumes ***REMOVED***
		args = append(args, "-v", volume)
	***REMOVED***

	if options.workDir != "" ***REMOVED***
		args = append(args, "-w", options.workDir)
	***REMOVED***

	if options.readOnly ***REMOVED***
		args = append(args, "--read-only")
	***REMOVED***

	args = append(args, "busybox", "/bin/sh", "-c", options.command)

	out, _ := dockerCmd(c, args...)

	containerID = strings.TrimSpace(out)

	out, _ = dockerCmd(c, "wait", containerID)

	exitCode := strings.TrimSpace(out)
	if exitCode != "0" ***REMOVED***
		out, _ = dockerCmd(c, "logs", containerID)
	***REMOVED***
	c.Assert(exitCode, checker.Equals, "0", check.Commentf("failed to make test container: %s", out))

	return
***REMOVED***

func makeCatFileCommand(path string) string ***REMOVED***
	return fmt.Sprintf("if [ -f %s ]; then cat %s; fi", path, path)
***REMOVED***

func cpPath(pathElements ...string) string ***REMOVED***
	localizedPathElements := make([]string, len(pathElements))
	for i, path := range pathElements ***REMOVED***
		localizedPathElements[i] = filepath.FromSlash(path)
	***REMOVED***
	return strings.Join(localizedPathElements, string(filepath.Separator))
***REMOVED***

func cpPathTrailingSep(pathElements ...string) string ***REMOVED***
	return fmt.Sprintf("%s%c", cpPath(pathElements...), filepath.Separator)
***REMOVED***

func containerCpPath(containerID string, pathElements ...string) string ***REMOVED***
	joined := strings.Join(pathElements, "/")
	return fmt.Sprintf("%s:%s", containerID, joined)
***REMOVED***

func containerCpPathTrailingSep(containerID string, pathElements ...string) string ***REMOVED***
	return fmt.Sprintf("%s/", containerCpPath(containerID, pathElements...))
***REMOVED***

func runDockerCp(c *check.C, src, dst string, params []string) (err error) ***REMOVED***
	c.Logf("running `docker cp %s %s %s`", strings.Join(params, " "), src, dst)

	args := []string***REMOVED***"cp"***REMOVED***

	args = append(args, params...)

	args = append(args, src, dst)

	out, _, err := runCommandWithOutput(exec.Command(dockerBinary, args...))
	if err != nil ***REMOVED***
		err = fmt.Errorf("error executing `docker cp` command: %s: %s", err, out)
	***REMOVED***

	return
***REMOVED***

func startContainerGetOutput(c *check.C, containerID string) (out string, err error) ***REMOVED***
	c.Logf("running `docker start -a %s`", containerID)

	args := []string***REMOVED***"start", "-a", containerID***REMOVED***

	out, _, err = runCommandWithOutput(exec.Command(dockerBinary, args...))
	if err != nil ***REMOVED***
		err = fmt.Errorf("error executing `docker start` command: %s: %s", err, out)
	***REMOVED***

	return
***REMOVED***

func getTestDir(c *check.C, label string) (tmpDir string) ***REMOVED***
	var err error

	tmpDir, err = ioutil.TempDir("", label)
	// unable to make temporary directory
	c.Assert(err, checker.IsNil)

	return
***REMOVED***

func isCpNotExist(err error) bool ***REMOVED***
	return strings.Contains(strings.ToLower(err.Error()), "could not find the file")
***REMOVED***

func isCpDirNotExist(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), archive.ErrDirNotExists.Error())
***REMOVED***

func isCpNotDir(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), archive.ErrNotDirectory.Error()) || strings.Contains(err.Error(), "filename, directory name, or volume label syntax is incorrect")
***REMOVED***

func isCpCannotCopyDir(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), archive.ErrCannotCopyDir.Error())
***REMOVED***

func isCpCannotCopyReadOnly(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), "marked read-only")
***REMOVED***

func isCannotOverwriteNonDirWithDir(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), "cannot overwrite non-directory")
***REMOVED***

func fileContentEquals(c *check.C, filename, contents string) (err error) ***REMOVED***
	c.Logf("checking that file %q contains %q\n", filename, contents)

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	expectedBytes, err := ioutil.ReadAll(strings.NewReader(contents))
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if !bytes.Equal(fileBytes, expectedBytes) ***REMOVED***
		err = fmt.Errorf("file content not equal - expected %q, got %q", string(expectedBytes), string(fileBytes))
	***REMOVED***

	return
***REMOVED***

func symlinkTargetEquals(c *check.C, symlink, expectedTarget string) (err error) ***REMOVED***
	c.Logf("checking that the symlink %q points to %q\n", symlink, expectedTarget)

	actualTarget, err := os.Readlink(symlink)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if actualTarget != expectedTarget ***REMOVED***
		err = fmt.Errorf("symlink target points to %q not %q", actualTarget, expectedTarget)
	***REMOVED***

	return
***REMOVED***

func containerStartOutputEquals(c *check.C, containerID, contents string) (err error) ***REMOVED***
	c.Logf("checking that container %q start output contains %q\n", containerID, contents)

	out, err := startContainerGetOutput(c, containerID)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if out != contents ***REMOVED***
		err = fmt.Errorf("output contents not equal - expected %q, got %q", contents, out)
	***REMOVED***

	return
***REMOVED***

func defaultVolumes(tmpDir string) []string ***REMOVED***
	if SameHostDaemon() ***REMOVED***
		return []string***REMOVED***
			"/vol1",
			fmt.Sprintf("%s:/vol2", tmpDir),
			fmt.Sprintf("%s:/vol3", filepath.Join(tmpDir, "vol3")),
			fmt.Sprintf("%s:/vol_ro:ro", filepath.Join(tmpDir, "vol_ro")),
		***REMOVED***
	***REMOVED***

	// Can't bind-mount volumes with separate host daemon.
	return []string***REMOVED***"/vol1", "/vol2", "/vol3", "/vol_ro:/vol_ro:ro"***REMOVED***
***REMOVED***
