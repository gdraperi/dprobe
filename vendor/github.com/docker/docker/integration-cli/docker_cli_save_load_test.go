package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	digest "github.com/opencontainers/go-digest"
)

// save a repo using gz compression and try to load it using stdout
func (s *DockerSuite) TestSaveXzAndLoadRepoStdout(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	name := "test-save-xz-and-load-repo-stdout"
	dockerCmd(c, "run", "--name", name, "busybox", "true")

	repoName := "foobar-save-load-test-xz-gz"
	out, _ := dockerCmd(c, "commit", name, repoName)

	dockerCmd(c, "inspect", repoName)

	repoTarball, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("xz", "-c"),
		exec.Command("gzip", "-c"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo: %v %v", out, err))
	deleteImages(repoName)

	icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "load"***REMOVED***,
		Stdin:   strings.NewReader(repoTarball),
	***REMOVED***).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
	***REMOVED***)

	after, _, err := dockerCmdWithError("inspect", repoName)
	c.Assert(err, checker.NotNil, check.Commentf("the repo should not exist: %v", after))
***REMOVED***

// save a repo using xz+gz compression and try to load it using stdout
func (s *DockerSuite) TestSaveXzGzAndLoadRepoStdout(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	name := "test-save-xz-gz-and-load-repo-stdout"
	dockerCmd(c, "run", "--name", name, "busybox", "true")

	repoName := "foobar-save-load-test-xz-gz"
	dockerCmd(c, "commit", name, repoName)

	dockerCmd(c, "inspect", repoName)

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("xz", "-c"),
		exec.Command("gzip", "-c"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo: %v %v", out, err))

	deleteImages(repoName)

	icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "load"***REMOVED***,
		Stdin:   strings.NewReader(out),
	***REMOVED***).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
	***REMOVED***)

	after, _, err := dockerCmdWithError("inspect", repoName)
	c.Assert(err, checker.NotNil, check.Commentf("the repo should not exist: %v", after))
***REMOVED***

func (s *DockerSuite) TestSaveSingleTag(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	repoName := "foobar-save-single-tag-test"
	dockerCmd(c, "tag", "busybox:latest", fmt.Sprintf("%v:latest", repoName))

	out, _ := dockerCmd(c, "images", "-q", "--no-trunc", repoName)
	cleanedImageID := strings.TrimSpace(out)

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", fmt.Sprintf("%v:latest", repoName)),
		exec.Command("tar", "t"),
		exec.Command("grep", "-E", fmt.Sprintf("(^repositories$|%v)", cleanedImageID)))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo with image ID and 'repositories' file: %s, %v", out, err))
***REMOVED***

func (s *DockerSuite) TestSaveCheckTimes(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	repoName := "busybox:latest"
	out, _ := dockerCmd(c, "inspect", repoName)
	data := []struct ***REMOVED***
		ID      string
		Created time.Time
	***REMOVED******REMOVED******REMOVED***
	err := json.Unmarshal([]byte(out), &data)
	c.Assert(err, checker.IsNil, check.Commentf("failed to marshal from %q: err %v", repoName, err))
	c.Assert(len(data), checker.Not(checker.Equals), 0, check.Commentf("failed to marshal the data from %q", repoName))
	tarTvTimeFormat := "2006-01-02 15:04"
	out, err = RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("tar", "tv"),
		exec.Command("grep", "-E", fmt.Sprintf("%s %s", data[0].Created.Format(tarTvTimeFormat), digest.Digest(data[0].ID).Hex())))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo with image ID and 'repositories' file: %s, %v", out, err))
***REMOVED***

func (s *DockerSuite) TestSaveImageId(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	repoName := "foobar-save-image-id-test"
	dockerCmd(c, "tag", "emptyfs:latest", fmt.Sprintf("%v:latest", repoName))

	out, _ := dockerCmd(c, "images", "-q", "--no-trunc", repoName)
	cleanedLongImageID := strings.TrimPrefix(strings.TrimSpace(out), "sha256:")

	out, _ = dockerCmd(c, "images", "-q", repoName)
	cleanedShortImageID := strings.TrimSpace(out)

	// Make sure IDs are not empty
	c.Assert(cleanedLongImageID, checker.Not(check.Equals), "", check.Commentf("Id should not be empty."))
	c.Assert(cleanedShortImageID, checker.Not(check.Equals), "", check.Commentf("Id should not be empty."))

	saveCmd := exec.Command(dockerBinary, "save", cleanedShortImageID)
	tarCmd := exec.Command("tar", "t")

	var err error
	tarCmd.Stdin, err = saveCmd.StdoutPipe()
	c.Assert(err, checker.IsNil, check.Commentf("cannot set stdout pipe for tar: %v", err))
	grepCmd := exec.Command("grep", cleanedLongImageID)
	grepCmd.Stdin, err = tarCmd.StdoutPipe()
	c.Assert(err, checker.IsNil, check.Commentf("cannot set stdout pipe for grep: %v", err))

	c.Assert(tarCmd.Start(), checker.IsNil, check.Commentf("tar failed with error: %v", err))
	c.Assert(saveCmd.Start(), checker.IsNil, check.Commentf("docker save failed with error: %v", err))
	defer func() ***REMOVED***
		saveCmd.Wait()
		tarCmd.Wait()
		dockerCmd(c, "rmi", repoName)
	***REMOVED***()

	out, _, err = runCommandWithOutput(grepCmd)

	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo with image ID: %s, %v", out, err))
***REMOVED***

// save a repo and try to load it using flags
func (s *DockerSuite) TestSaveAndLoadRepoFlags(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	name := "test-save-and-load-repo-flags"
	dockerCmd(c, "run", "--name", name, "busybox", "true")

	repoName := "foobar-save-load-test"

	deleteImages(repoName)
	dockerCmd(c, "commit", name, repoName)

	before, _ := dockerCmd(c, "inspect", repoName)

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command(dockerBinary, "load"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save and load repo: %s, %v", out, err))

	after, _ := dockerCmd(c, "inspect", repoName)
	c.Assert(before, checker.Equals, after, check.Commentf("inspect is not the same after a save / load"))
***REMOVED***

func (s *DockerSuite) TestSaveWithNoExistImage(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	imgName := "foobar-non-existing-image"

	out, _, err := dockerCmdWithError("save", "-o", "test-img.tar", imgName)
	c.Assert(err, checker.NotNil, check.Commentf("save image should fail for non-existing image"))
	c.Assert(out, checker.Contains, fmt.Sprintf("No such image: %s", imgName))
***REMOVED***

func (s *DockerSuite) TestSaveMultipleNames(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	repoName := "foobar-save-multi-name-test"

	// Make one image
	dockerCmd(c, "tag", "emptyfs:latest", fmt.Sprintf("%v-one:latest", repoName))

	// Make two images
	dockerCmd(c, "tag", "emptyfs:latest", fmt.Sprintf("%v-two:latest", repoName))

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", fmt.Sprintf("%v-one", repoName), fmt.Sprintf("%v-two:latest", repoName)),
		exec.Command("tar", "xO", "repositories"),
		exec.Command("grep", "-q", "-E", "(-one|-two)"),
	)
	c.Assert(err, checker.IsNil, check.Commentf("failed to save multiple repos: %s, %v", out, err))
***REMOVED***

func (s *DockerSuite) TestSaveRepoWithMultipleImages(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	makeImage := func(from string, tag string) string ***REMOVED***
		var (
			out string
		)
		out, _ = dockerCmd(c, "run", "-d", from, "true")
		cleanedContainerID := strings.TrimSpace(out)

		out, _ = dockerCmd(c, "commit", cleanedContainerID, tag)
		imageID := strings.TrimSpace(out)
		return imageID
	***REMOVED***

	repoName := "foobar-save-multi-images-test"
	tagFoo := repoName + ":foo"
	tagBar := repoName + ":bar"

	idFoo := makeImage("busybox:latest", tagFoo)
	idBar := makeImage("busybox:latest", tagBar)

	deleteImages(repoName)

	// create the archive
	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName, "busybox:latest"),
		exec.Command("tar", "t"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save multiple images: %s, %v", out, err))

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var actual []string
	for _, l := range lines ***REMOVED***
		if regexp.MustCompile("^[a-f0-9]***REMOVED***64***REMOVED***\\.json$").Match([]byte(l)) ***REMOVED***
			actual = append(actual, strings.TrimSuffix(l, ".json"))
		***REMOVED***
	***REMOVED***

	// make the list of expected layers
	out = inspectField(c, "busybox:latest", "Id")
	expected := []string***REMOVED***strings.TrimSpace(out), idFoo, idBar***REMOVED***

	// prefixes are not in tar
	for i := range expected ***REMOVED***
		expected[i] = digest.Digest(expected[i]).Hex()
	***REMOVED***

	sort.Strings(actual)
	sort.Strings(expected)
	c.Assert(actual, checker.DeepEquals, expected, check.Commentf("archive does not contains the right layers: got %v, expected %v, output: %q", actual, expected, out))
***REMOVED***

// Issue #6722 #5892 ensure directories are included in changes
func (s *DockerSuite) TestSaveDirectoryPermissions(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	layerEntries := []string***REMOVED***"opt/", "opt/a/", "opt/a/b/", "opt/a/b/c"***REMOVED***
	layerEntriesAUFS := []string***REMOVED***"./", ".wh..wh.aufs", ".wh..wh.orph/", ".wh..wh.plnk/", "opt/", "opt/a/", "opt/a/b/", "opt/a/b/c"***REMOVED***

	name := "save-directory-permissions"
	tmpDir, err := ioutil.TempDir("", "save-layers-with-directories")
	c.Assert(err, checker.IsNil, check.Commentf("failed to create temporary directory: %s", err))
	extractionDirectory := filepath.Join(tmpDir, "image-extraction-dir")
	os.Mkdir(extractionDirectory, 0777)

	defer os.RemoveAll(tmpDir)
	buildImageSuccessfully(c, name, build.WithDockerfile(`FROM busybox
	RUN adduser -D user && mkdir -p /opt/a/b && chown -R user:user /opt/a
	RUN touch /opt/a/b/c && chown user:user /opt/a/b/c`))

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", name),
		exec.Command("tar", "-xf", "-", "-C", extractionDirectory),
	)
	c.Assert(err, checker.IsNil, check.Commentf("failed to save and extract image: %s", out))

	dirs, err := ioutil.ReadDir(extractionDirectory)
	c.Assert(err, checker.IsNil, check.Commentf("failed to get a listing of the layer directories: %s", err))

	found := false
	for _, entry := range dirs ***REMOVED***
		var entriesSansDev []string
		if entry.IsDir() ***REMOVED***
			layerPath := filepath.Join(extractionDirectory, entry.Name(), "layer.tar")

			f, err := os.Open(layerPath)
			c.Assert(err, checker.IsNil, check.Commentf("failed to open %s: %s", layerPath, err))
			defer f.Close()

			entries, err := listTar(f)
			for _, e := range entries ***REMOVED***
				if !strings.Contains(e, "dev/") ***REMOVED***
					entriesSansDev = append(entriesSansDev, e)
				***REMOVED***
			***REMOVED***
			c.Assert(err, checker.IsNil, check.Commentf("encountered error while listing tar entries: %s", err))

			if reflect.DeepEqual(entriesSansDev, layerEntries) || reflect.DeepEqual(entriesSansDev, layerEntriesAUFS) ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	c.Assert(found, checker.Equals, true, check.Commentf("failed to find the layer with the right content listing"))

***REMOVED***

func listTar(f io.Reader) ([]string, error) ***REMOVED***
	tr := tar.NewReader(f)
	var entries []string

	for ***REMOVED***
		th, err := tr.Next()
		if err == io.EOF ***REMOVED***
			// end of tar archive
			return entries, nil
		***REMOVED***
		if err != nil ***REMOVED***
			return entries, err
		***REMOVED***
		entries = append(entries, th.Name)
	***REMOVED***
***REMOVED***

// Test loading a weird image where one of the layers is of zero size.
// The layer.tar file is actually zero bytes, no padding or anything else.
// See issue: 18170
func (s *DockerSuite) TestLoadZeroSizeLayer(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	dockerCmd(c, "load", "-i", "fixtures/load/emptyLayer.tar")
***REMOVED***

func (s *DockerSuite) TestSaveLoadParents(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	makeImage := func(from string, addfile string) string ***REMOVED***
		var (
			out string
		)
		out, _ = dockerCmd(c, "run", "-d", from, "touch", addfile)
		cleanedContainerID := strings.TrimSpace(out)

		out, _ = dockerCmd(c, "commit", cleanedContainerID)
		imageID := strings.TrimSpace(out)

		dockerCmd(c, "rm", "-f", cleanedContainerID)
		return imageID
	***REMOVED***

	idFoo := makeImage("busybox", "foo")
	idBar := makeImage(idFoo, "bar")

	tmpDir, err := ioutil.TempDir("", "save-load-parents")
	c.Assert(err, checker.IsNil)
	defer os.RemoveAll(tmpDir)

	c.Log("tmpdir", tmpDir)

	outfile := filepath.Join(tmpDir, "out.tar")

	dockerCmd(c, "save", "-o", outfile, idBar, idFoo)
	dockerCmd(c, "rmi", idBar)
	dockerCmd(c, "load", "-i", outfile)

	inspectOut := inspectField(c, idBar, "Parent")
	c.Assert(inspectOut, checker.Equals, idFoo)

	inspectOut = inspectField(c, idFoo, "Parent")
	c.Assert(inspectOut, checker.Equals, "")
***REMOVED***

func (s *DockerSuite) TestSaveLoadNoTag(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	name := "saveloadnotag"

	buildImageSuccessfully(c, name, build.WithDockerfile("FROM busybox\nENV foo=bar"))
	id := inspectField(c, name, "Id")

	// Test to make sure that save w/o name just shows imageID during load
	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", id),
		exec.Command(dockerBinary, "load"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save and load repo: %s, %v", out, err))

	// Should not show 'name' but should show the image ID during the load
	c.Assert(out, checker.Not(checker.Contains), "Loaded image: ")
	c.Assert(out, checker.Contains, "Loaded image ID:")
	c.Assert(out, checker.Contains, id)

	// Test to make sure that save by name shows that name during load
	out, err = RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", name),
		exec.Command(dockerBinary, "load"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save and load repo: %s, %v", out, err))
	c.Assert(out, checker.Contains, "Loaded image: "+name+":latest")
	c.Assert(out, checker.Not(checker.Contains), "Loaded image ID:")
***REMOVED***
