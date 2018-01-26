package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	"golang.org/x/net/context"
)

func (s *DockerSuite) TestVolumeCLICreate(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create")

	_, _, err := dockerCmdWithError("volume", "create", "-d", "nosuchdriver")
	c.Assert(err, check.NotNil)

	// test using hidden --name option
	out, _ := dockerCmd(c, "volume", "create", "--name=test")
	name := strings.TrimSpace(out)
	c.Assert(name, check.Equals, "test")

	out, _ = dockerCmd(c, "volume", "create", "test2")
	name = strings.TrimSpace(out)
	c.Assert(name, check.Equals, "test2")
***REMOVED***

func (s *DockerSuite) TestVolumeCLIInspect(c *check.C) ***REMOVED***
	c.Assert(
		exec.Command(dockerBinary, "volume", "inspect", "doesnotexist").Run(),
		check.Not(check.IsNil),
		check.Commentf("volume inspect should error on non-existent volume"),
	)

	out, _ := dockerCmd(c, "volume", "create")
	name := strings.TrimSpace(out)
	out, _ = dockerCmd(c, "volume", "inspect", "--format=***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***", name)
	c.Assert(strings.TrimSpace(out), check.Equals, name)

	dockerCmd(c, "volume", "create", "test")
	out, _ = dockerCmd(c, "volume", "inspect", "--format=***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***", "test")
	c.Assert(strings.TrimSpace(out), check.Equals, "test")
***REMOVED***

func (s *DockerSuite) TestVolumeCLIInspectMulti(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create", "test1")
	dockerCmd(c, "volume", "create", "test2")
	dockerCmd(c, "volume", "create", "test3")

	result := dockerCmdWithResult("volume", "inspect", "--format=***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***", "test1", "test2", "doesnotexist", "test3")
	result.Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Err:      "No such volume: doesnotexist",
	***REMOVED***)

	out := result.Stdout()
	c.Assert(out, checker.Contains, "test1")
	c.Assert(out, checker.Contains, "test2")
	c.Assert(out, checker.Contains, "test3")
***REMOVED***

func (s *DockerSuite) TestVolumeCLILs(c *check.C) ***REMOVED***
	prefix, _ := getPrefixAndSlashFromDaemonPlatform()
	dockerCmd(c, "volume", "create", "aaa")

	dockerCmd(c, "volume", "create", "test")

	dockerCmd(c, "volume", "create", "soo")
	dockerCmd(c, "run", "-v", "soo:"+prefix+"/foo", "busybox", "ls", "/")

	out, _ := dockerCmd(c, "volume", "ls", "-q")
	assertVolumesInList(c, out, []string***REMOVED***"aaa", "soo", "test"***REMOVED***)
***REMOVED***

func (s *DockerSuite) TestVolumeLsFormat(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create", "aaa")
	dockerCmd(c, "volume", "create", "test")
	dockerCmd(c, "volume", "create", "soo")

	out, _ := dockerCmd(c, "volume", "ls", "--format", "***REMOVED******REMOVED***.Name***REMOVED******REMOVED***")
	assertVolumesInList(c, out, []string***REMOVED***"aaa", "soo", "test"***REMOVED***)
***REMOVED***

func (s *DockerSuite) TestVolumeLsFormatDefaultFormat(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create", "aaa")
	dockerCmd(c, "volume", "create", "test")
	dockerCmd(c, "volume", "create", "soo")

	config := `***REMOVED***
		"volumesFormat": "***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED*** default"
***REMOVED***`
	d, err := ioutil.TempDir("", "integration-cli-")
	c.Assert(err, checker.IsNil)
	defer os.RemoveAll(d)

	err = ioutil.WriteFile(filepath.Join(d, "config.json"), []byte(config), 0644)
	c.Assert(err, checker.IsNil)

	out, _ := dockerCmd(c, "--config", d, "volume", "ls")
	assertVolumesInList(c, out, []string***REMOVED***"aaa default", "soo default", "test default"***REMOVED***)
***REMOVED***

// assertVolList checks volume retrieved with ls command
// equals to expected volume list
// note: out should be `volume ls [option]` result
func assertVolList(c *check.C, out string, expectVols []string) ***REMOVED***
	lines := strings.Split(out, "\n")
	var volList []string
	for _, line := range lines[1 : len(lines)-1] ***REMOVED***
		volFields := strings.Fields(line)
		// wrap all volume name in volList
		volList = append(volList, volFields[1])
	***REMOVED***

	// volume ls should contains all expected volumes
	c.Assert(volList, checker.DeepEquals, expectVols)
***REMOVED***

func assertVolumesInList(c *check.C, out string, expected []string) ***REMOVED***
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, expect := range expected ***REMOVED***
		found := false
		for _, v := range lines ***REMOVED***
			found = v == expect
			if found ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		c.Assert(found, checker.Equals, true, check.Commentf("Expected volume not found: %v, got: %v", expect, lines))
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestVolumeCLILsFilterDangling(c *check.C) ***REMOVED***
	prefix, _ := getPrefixAndSlashFromDaemonPlatform()
	dockerCmd(c, "volume", "create", "testnotinuse1")
	dockerCmd(c, "volume", "create", "testisinuse1")
	dockerCmd(c, "volume", "create", "testisinuse2")

	// Make sure both "created" (but not started), and started
	// containers are included in reference counting
	dockerCmd(c, "run", "--name", "volume-test1", "-v", "testisinuse1:"+prefix+"/foo", "busybox", "true")
	dockerCmd(c, "create", "--name", "volume-test2", "-v", "testisinuse2:"+prefix+"/foo", "busybox", "true")

	out, _ := dockerCmd(c, "volume", "ls")

	// No filter, all volumes should show
	c.Assert(out, checker.Contains, "testnotinuse1\n", check.Commentf("expected volume 'testnotinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse1\n", check.Commentf("expected volume 'testisinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse2\n", check.Commentf("expected volume 'testisinuse2' in output"))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "dangling=false")

	// Explicitly disabling dangling
	c.Assert(out, check.Not(checker.Contains), "testnotinuse1\n", check.Commentf("expected volume 'testnotinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse1\n", check.Commentf("expected volume 'testisinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse2\n", check.Commentf("expected volume 'testisinuse2' in output"))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "dangling=true")

	// Filter "dangling" volumes; only "dangling" (unused) volumes should be in the output
	c.Assert(out, checker.Contains, "testnotinuse1\n", check.Commentf("expected volume 'testnotinuse1' in output"))
	c.Assert(out, check.Not(checker.Contains), "testisinuse1\n", check.Commentf("volume 'testisinuse1' in output, but not expected"))
	c.Assert(out, check.Not(checker.Contains), "testisinuse2\n", check.Commentf("volume 'testisinuse2' in output, but not expected"))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "dangling=1")
	// Filter "dangling" volumes; only "dangling" (unused) volumes should be in the output, dangling also accept 1
	c.Assert(out, checker.Contains, "testnotinuse1\n", check.Commentf("expected volume 'testnotinuse1' in output"))
	c.Assert(out, check.Not(checker.Contains), "testisinuse1\n", check.Commentf("volume 'testisinuse1' in output, but not expected"))
	c.Assert(out, check.Not(checker.Contains), "testisinuse2\n", check.Commentf("volume 'testisinuse2' in output, but not expected"))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "dangling=0")
	// dangling=0 is same as dangling=false case
	c.Assert(out, check.Not(checker.Contains), "testnotinuse1\n", check.Commentf("expected volume 'testnotinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse1\n", check.Commentf("expected volume 'testisinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse2\n", check.Commentf("expected volume 'testisinuse2' in output"))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "name=testisin")
	c.Assert(out, check.Not(checker.Contains), "testnotinuse1\n", check.Commentf("expected volume 'testnotinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse1\n", check.Commentf("expected volume 'testisinuse1' in output"))
	c.Assert(out, checker.Contains, "testisinuse2\n", check.Commentf("expected volume 'testisinuse2' in output"))
***REMOVED***

func (s *DockerSuite) TestVolumeCLILsErrorWithInvalidFilterName(c *check.C) ***REMOVED***
	out, _, err := dockerCmdWithError("volume", "ls", "-f", "FOO=123")
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, "Invalid filter")
***REMOVED***

func (s *DockerSuite) TestVolumeCLILsWithIncorrectFilterValue(c *check.C) ***REMOVED***
	out, _, err := dockerCmdWithError("volume", "ls", "-f", "dangling=invalid")
	c.Assert(err, check.NotNil)
	c.Assert(out, checker.Contains, "Invalid filter")
***REMOVED***

func (s *DockerSuite) TestVolumeCLIRm(c *check.C) ***REMOVED***
	prefix, _ := getPrefixAndSlashFromDaemonPlatform()
	out, _ := dockerCmd(c, "volume", "create")
	id := strings.TrimSpace(out)

	dockerCmd(c, "volume", "create", "test")
	dockerCmd(c, "volume", "rm", id)
	dockerCmd(c, "volume", "rm", "test")

	volumeID := "testing"
	dockerCmd(c, "run", "-v", volumeID+":"+prefix+"/foo", "--name=test", "busybox", "sh", "-c", "echo hello > /foo/bar")

	icmd.RunCommand(dockerBinary, "volume", "rm", "testing").Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Error:    "exit status 1",
	***REMOVED***)

	out, _ = dockerCmd(c, "run", "--volumes-from=test", "--name=test2", "busybox", "sh", "-c", "cat /foo/bar")
	c.Assert(strings.TrimSpace(out), check.Equals, "hello")
	dockerCmd(c, "rm", "-fv", "test2")
	dockerCmd(c, "volume", "inspect", volumeID)
	dockerCmd(c, "rm", "-f", "test")

	out, _ = dockerCmd(c, "run", "--name=test2", "-v", volumeID+":"+prefix+"/foo", "busybox", "sh", "-c", "cat /foo/bar")
	c.Assert(strings.TrimSpace(out), check.Equals, "hello", check.Commentf("volume data was removed"))
	dockerCmd(c, "rm", "test2")

	dockerCmd(c, "volume", "rm", volumeID)
	c.Assert(
		exec.Command("volume", "rm", "doesnotexist").Run(),
		check.Not(check.IsNil),
		check.Commentf("volume rm should fail with non-existent volume"),
	)
***REMOVED***

// FIXME(vdemeester) should be a unit test in cli/command/volume package
func (s *DockerSuite) TestVolumeCLINoArgs(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "volume")
	// no args should produce the cmd usage output
	usage := "Usage:	docker volume COMMAND"
	c.Assert(out, checker.Contains, usage)

	// invalid arg should error and show the command usage on stderr
	icmd.RunCommand(dockerBinary, "volume", "somearg").Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Error:    "exit status 1",
		Err:      usage,
	***REMOVED***)

	// invalid flag should error and show the flag error and cmd usage
	result := icmd.RunCommand(dockerBinary, "volume", "--no-such-flag")
	result.Assert(c, icmd.Expected***REMOVED***
		ExitCode: 125,
		Error:    "exit status 125",
		Err:      usage,
	***REMOVED***)
	c.Assert(result.Stderr(), checker.Contains, "unknown flag: --no-such-flag")
***REMOVED***

func (s *DockerSuite) TestVolumeCLIInspectTmplError(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "volume", "create")
	name := strings.TrimSpace(out)

	out, exitCode, err := dockerCmdWithError("volume", "inspect", "--format='***REMOVED******REMOVED*** .FooBar ***REMOVED******REMOVED***'", name)
	c.Assert(err, checker.NotNil, check.Commentf("Output: %s", out))
	c.Assert(exitCode, checker.Equals, 1, check.Commentf("Output: %s", out))
	c.Assert(out, checker.Contains, "Template parsing error")
***REMOVED***

func (s *DockerSuite) TestVolumeCLICreateWithOpts(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	dockerCmd(c, "volume", "create", "-d", "local", "test", "--opt=type=tmpfs", "--opt=device=tmpfs", "--opt=o=size=1m,uid=1000")
	out, _ := dockerCmd(c, "run", "-v", "test:/foo", "busybox", "mount")

	mounts := strings.Split(out, "\n")
	var found bool
	for _, m := range mounts ***REMOVED***
		if strings.Contains(m, "/foo") ***REMOVED***
			found = true
			info := strings.Fields(m)
			// tmpfs on <path> type tmpfs (rw,relatime,size=1024k,uid=1000)
			c.Assert(info[0], checker.Equals, "tmpfs")
			c.Assert(info[2], checker.Equals, "/foo")
			c.Assert(info[4], checker.Equals, "tmpfs")
			c.Assert(info[5], checker.Contains, "uid=1000")
			c.Assert(info[5], checker.Contains, "size=1024k")
			break
		***REMOVED***
	***REMOVED***
	c.Assert(found, checker.Equals, true)
***REMOVED***

func (s *DockerSuite) TestVolumeCLICreateLabel(c *check.C) ***REMOVED***
	testVol := "testvolcreatelabel"
	testLabel := "foo"
	testValue := "bar"

	out, _, err := dockerCmdWithError("volume", "create", "--label", testLabel+"="+testValue, testVol)
	c.Assert(err, check.IsNil)

	out, _ = dockerCmd(c, "volume", "inspect", "--format=***REMOVED******REMOVED*** .Labels."+testLabel+" ***REMOVED******REMOVED***", testVol)
	c.Assert(strings.TrimSpace(out), check.Equals, testValue)
***REMOVED***

func (s *DockerSuite) TestVolumeCLICreateLabelMultiple(c *check.C) ***REMOVED***
	testVol := "testvolcreatelabel"

	testLabels := map[string]string***REMOVED***
		"foo": "bar",
		"baz": "foo",
	***REMOVED***

	args := []string***REMOVED***
		"volume",
		"create",
		testVol,
	***REMOVED***

	for k, v := range testLabels ***REMOVED***
		args = append(args, "--label", k+"="+v)
	***REMOVED***

	out, _, err := dockerCmdWithError(args...)
	c.Assert(err, check.IsNil)

	for k, v := range testLabels ***REMOVED***
		out, _ = dockerCmd(c, "volume", "inspect", "--format=***REMOVED******REMOVED*** .Labels."+k+" ***REMOVED******REMOVED***", testVol)
		c.Assert(strings.TrimSpace(out), check.Equals, v)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestVolumeCLILsFilterLabels(c *check.C) ***REMOVED***
	testVol1 := "testvolcreatelabel-1"
	out, _, err := dockerCmdWithError("volume", "create", "--label", "foo=bar1", testVol1)
	c.Assert(err, check.IsNil)

	testVol2 := "testvolcreatelabel-2"
	out, _, err = dockerCmdWithError("volume", "create", "--label", "foo=bar2", testVol2)
	c.Assert(err, check.IsNil)

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "label=foo")

	// filter with label=key
	c.Assert(out, checker.Contains, "testvolcreatelabel-1\n", check.Commentf("expected volume 'testvolcreatelabel-1' in output"))
	c.Assert(out, checker.Contains, "testvolcreatelabel-2\n", check.Commentf("expected volume 'testvolcreatelabel-2' in output"))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "label=foo=bar1")

	// filter with label=key=value
	c.Assert(out, checker.Contains, "testvolcreatelabel-1\n", check.Commentf("expected volume 'testvolcreatelabel-1' in output"))
	c.Assert(out, check.Not(checker.Contains), "testvolcreatelabel-2\n", check.Commentf("expected volume 'testvolcreatelabel-2 in output"))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "label=non-exist")
	outArr := strings.Split(strings.TrimSpace(out), "\n")
	c.Assert(len(outArr), check.Equals, 1, check.Commentf("\n%s", out))

	out, _ = dockerCmd(c, "volume", "ls", "--filter", "label=foo=non-exist")
	outArr = strings.Split(strings.TrimSpace(out), "\n")
	c.Assert(len(outArr), check.Equals, 1, check.Commentf("\n%s", out))
***REMOVED***

func (s *DockerSuite) TestVolumeCLILsFilterDrivers(c *check.C) ***REMOVED***
	// using default volume driver local to create volumes
	testVol1 := "testvol-1"
	out, _, err := dockerCmdWithError("volume", "create", testVol1)
	c.Assert(err, check.IsNil)

	testVol2 := "testvol-2"
	out, _, err = dockerCmdWithError("volume", "create", testVol2)
	c.Assert(err, check.IsNil)

	// filter with driver=local
	out, _ = dockerCmd(c, "volume", "ls", "--filter", "driver=local")
	c.Assert(out, checker.Contains, "testvol-1\n", check.Commentf("expected volume 'testvol-1' in output"))
	c.Assert(out, checker.Contains, "testvol-2\n", check.Commentf("expected volume 'testvol-2' in output"))

	// filter with driver=invaliddriver
	out, _ = dockerCmd(c, "volume", "ls", "--filter", "driver=invaliddriver")
	outArr := strings.Split(strings.TrimSpace(out), "\n")
	c.Assert(len(outArr), check.Equals, 1, check.Commentf("\n%s", out))

	// filter with driver=loca
	out, _ = dockerCmd(c, "volume", "ls", "--filter", "driver=loca")
	outArr = strings.Split(strings.TrimSpace(out), "\n")
	c.Assert(len(outArr), check.Equals, 1, check.Commentf("\n%s", out))

	// filter with driver=
	out, _ = dockerCmd(c, "volume", "ls", "--filter", "driver=")
	outArr = strings.Split(strings.TrimSpace(out), "\n")
	c.Assert(len(outArr), check.Equals, 1, check.Commentf("\n%s", out))
***REMOVED***

func (s *DockerSuite) TestVolumeCLIRmForceUsage(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "volume", "create")
	id := strings.TrimSpace(out)

	dockerCmd(c, "volume", "rm", "-f", id)
	dockerCmd(c, "volume", "rm", "--force", "nonexist")
***REMOVED***

func (s *DockerSuite) TestVolumeCLIRmForce(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon, DaemonIsLinux)

	name := "test"
	out, _ := dockerCmd(c, "volume", "create", name)
	id := strings.TrimSpace(out)
	c.Assert(id, checker.Equals, name)

	out, _ = dockerCmd(c, "volume", "inspect", "--format", "***REMOVED******REMOVED***.Mountpoint***REMOVED******REMOVED***", name)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")
	// Mountpoint is in the form of "/var/lib/docker/volumes/.../_data", removing `/_data`
	path := strings.TrimSuffix(strings.TrimSpace(out), "/_data")
	icmd.RunCommand("rm", "-rf", path).Assert(c, icmd.Success)

	dockerCmd(c, "volume", "rm", "-f", name)
	out, _ = dockerCmd(c, "volume", "ls")
	c.Assert(out, checker.Not(checker.Contains), name)
	dockerCmd(c, "volume", "create", name)
	out, _ = dockerCmd(c, "volume", "ls")
	c.Assert(out, checker.Contains, name)
***REMOVED***

// TestVolumeCLIRmForceInUse verifies that repeated `docker volume rm -f` calls does not remove a volume
// if it is in use. Test case for https://github.com/docker/docker/issues/31446
func (s *DockerSuite) TestVolumeCLIRmForceInUse(c *check.C) ***REMOVED***
	name := "testvolume"
	out, _ := dockerCmd(c, "volume", "create", name)
	id := strings.TrimSpace(out)
	c.Assert(id, checker.Equals, name)

	prefix, slash := getPrefixAndSlashFromDaemonPlatform()
	out, e := dockerCmd(c, "create", "-v", "testvolume:"+prefix+slash+"foo", "busybox")
	cid := strings.TrimSpace(out)

	_, _, err := dockerCmdWithError("volume", "rm", "-f", name)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), checker.Contains, "volume is in use")
	out, _ = dockerCmd(c, "volume", "ls")
	c.Assert(out, checker.Contains, name)

	// The original issue did not _remove_ the volume from the list
	// the first time. But a second call to `volume rm` removed it.
	// Calling `volume rm` a second time to confirm it's not removed
	// when calling twice.
	_, _, err = dockerCmdWithError("volume", "rm", "-f", name)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), checker.Contains, "volume is in use")
	out, _ = dockerCmd(c, "volume", "ls")
	c.Assert(out, checker.Contains, name)

	// Verify removing the volume after the container is removed works
	_, e = dockerCmd(c, "rm", cid)
	c.Assert(e, check.Equals, 0)

	_, e = dockerCmd(c, "volume", "rm", "-f", name)
	c.Assert(e, check.Equals, 0)

	out, e = dockerCmd(c, "volume", "ls")
	c.Assert(e, check.Equals, 0)
	c.Assert(out, checker.Not(checker.Contains), name)
***REMOVED***

func (s *DockerSuite) TestVolumeCliInspectWithVolumeOpts(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	// Without options
	name := "test1"
	dockerCmd(c, "volume", "create", "-d", "local", name)
	out, _ := dockerCmd(c, "volume", "inspect", "--format=***REMOVED******REMOVED*** .Options ***REMOVED******REMOVED***", name)
	c.Assert(strings.TrimSpace(out), checker.Contains, "map[]")

	// With options
	name = "test2"
	k1, v1 := "type", "tmpfs"
	k2, v2 := "device", "tmpfs"
	k3, v3 := "o", "size=1m,uid=1000"
	dockerCmd(c, "volume", "create", "-d", "local", name, "--opt", fmt.Sprintf("%s=%s", k1, v1), "--opt", fmt.Sprintf("%s=%s", k2, v2), "--opt", fmt.Sprintf("%s=%s", k3, v3))
	out, _ = dockerCmd(c, "volume", "inspect", "--format=***REMOVED******REMOVED*** .Options ***REMOVED******REMOVED***", name)
	c.Assert(strings.TrimSpace(out), checker.Contains, fmt.Sprintf("%s:%s", k1, v1))
	c.Assert(strings.TrimSpace(out), checker.Contains, fmt.Sprintf("%s:%s", k2, v2))
	c.Assert(strings.TrimSpace(out), checker.Contains, fmt.Sprintf("%s:%s", k3, v3))
***REMOVED***

// Test case (1) for 21845: duplicate targets for --volumes-from
func (s *DockerSuite) TestDuplicateMountpointsForVolumesFrom(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	image := "vimage"
	buildImageSuccessfully(c, image, build.WithDockerfile(`
		FROM busybox
		VOLUME ["/tmp/data"]`))

	dockerCmd(c, "run", "--name=data1", image, "true")
	dockerCmd(c, "run", "--name=data2", image, "true")

	out, _ := dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "data1")
	data1 := strings.TrimSpace(out)
	c.Assert(data1, checker.Not(checker.Equals), "")

	out, _ = dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "data2")
	data2 := strings.TrimSpace(out)
	c.Assert(data2, checker.Not(checker.Equals), "")

	// Both volume should exist
	out, _ = dockerCmd(c, "volume", "ls", "-q")
	c.Assert(strings.TrimSpace(out), checker.Contains, data1)
	c.Assert(strings.TrimSpace(out), checker.Contains, data2)

	out, _, err := dockerCmdWithError("run", "--name=app", "--volumes-from=data1", "--volumes-from=data2", "-d", "busybox", "top")
	c.Assert(err, checker.IsNil, check.Commentf("Out: %s", out))

	// Only the second volume will be referenced, this is backward compatible
	out, _ = dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "app")
	c.Assert(strings.TrimSpace(out), checker.Equals, data2)

	dockerCmd(c, "rm", "-f", "-v", "app")
	dockerCmd(c, "rm", "-f", "-v", "data1")
	dockerCmd(c, "rm", "-f", "-v", "data2")

	// Both volume should not exist
	out, _ = dockerCmd(c, "volume", "ls", "-q")
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data1)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data2)
***REMOVED***

// Test case (2) for 21845: duplicate targets for --volumes-from and -v (bind)
func (s *DockerSuite) TestDuplicateMountpointsForVolumesFromAndBind(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	image := "vimage"
	buildImageSuccessfully(c, image, build.WithDockerfile(`
                FROM busybox
                VOLUME ["/tmp/data"]`))

	dockerCmd(c, "run", "--name=data1", image, "true")
	dockerCmd(c, "run", "--name=data2", image, "true")

	out, _ := dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "data1")
	data1 := strings.TrimSpace(out)
	c.Assert(data1, checker.Not(checker.Equals), "")

	out, _ = dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "data2")
	data2 := strings.TrimSpace(out)
	c.Assert(data2, checker.Not(checker.Equals), "")

	// Both volume should exist
	out, _ = dockerCmd(c, "volume", "ls", "-q")
	c.Assert(strings.TrimSpace(out), checker.Contains, data1)
	c.Assert(strings.TrimSpace(out), checker.Contains, data2)

	// /tmp/data is automatically created, because we are not using the modern mount API here
	out, _, err := dockerCmdWithError("run", "--name=app", "--volumes-from=data1", "--volumes-from=data2", "-v", "/tmp/data:/tmp/data", "-d", "busybox", "top")
	c.Assert(err, checker.IsNil, check.Commentf("Out: %s", out))

	// No volume will be referenced (mount is /tmp/data), this is backward compatible
	out, _ = dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "app")
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data1)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data2)

	dockerCmd(c, "rm", "-f", "-v", "app")
	dockerCmd(c, "rm", "-f", "-v", "data1")
	dockerCmd(c, "rm", "-f", "-v", "data2")

	// Both volume should not exist
	out, _ = dockerCmd(c, "volume", "ls", "-q")
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data1)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data2)
***REMOVED***

// Test case (3) for 21845: duplicate targets for --volumes-from and `Mounts` (API only)
func (s *DockerSuite) TestDuplicateMountpointsForVolumesFromAndMounts(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon, DaemonIsLinux)

	image := "vimage"
	buildImageSuccessfully(c, image, build.WithDockerfile(`
                FROM busybox
                VOLUME ["/tmp/data"]`))

	dockerCmd(c, "run", "--name=data1", image, "true")
	dockerCmd(c, "run", "--name=data2", image, "true")

	out, _ := dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "data1")
	data1 := strings.TrimSpace(out)
	c.Assert(data1, checker.Not(checker.Equals), "")

	out, _ = dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "data2")
	data2 := strings.TrimSpace(out)
	c.Assert(data2, checker.Not(checker.Equals), "")

	// Both volume should exist
	out, _ = dockerCmd(c, "volume", "ls", "-q")
	c.Assert(strings.TrimSpace(out), checker.Contains, data1)
	c.Assert(strings.TrimSpace(out), checker.Contains, data2)

	err := os.MkdirAll("/tmp/data", 0755)
	c.Assert(err, checker.IsNil)
	// Mounts is available in API
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	config := container.Config***REMOVED***
		Cmd:   []string***REMOVED***"top"***REMOVED***,
		Image: "busybox",
	***REMOVED***

	hostConfig := container.HostConfig***REMOVED***
		VolumesFrom: []string***REMOVED***"data1", "data2"***REMOVED***,
		Mounts: []mount.Mount***REMOVED***
			***REMOVED***
				Type:   "bind",
				Source: "/tmp/data",
				Target: "/tmp/data",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &network.NetworkingConfig***REMOVED******REMOVED***, "app")

	c.Assert(err, checker.IsNil)

	// No volume will be referenced (mount is /tmp/data), this is backward compatible
	out, _ = dockerCmd(c, "inspect", "--format", "***REMOVED******REMOVED***(index .Mounts 0).Name***REMOVED******REMOVED***", "app")
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data1)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data2)

	dockerCmd(c, "rm", "-f", "-v", "app")
	dockerCmd(c, "rm", "-f", "-v", "data1")
	dockerCmd(c, "rm", "-f", "-v", "data2")

	// Both volume should not exist
	out, _ = dockerCmd(c, "volume", "ls", "-q")
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data1)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Contains), data2)
***REMOVED***
