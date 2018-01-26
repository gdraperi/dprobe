package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/docker/docker/integration-cli/request"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/volume"
	"github.com/docker/go-connections/nat"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/poll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func (s *DockerSuite) TestContainerAPIGetAll(c *check.C) ***REMOVED***
	startCount := getContainerCount(c)
	name := "getall"
	dockerCmd(c, "run", "--name", name, "busybox", "true")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.ContainerListOptions***REMOVED***
		All: true,
	***REMOVED***
	containers, err := cli.ContainerList(context.Background(), options)
	c.Assert(err, checker.IsNil)
	c.Assert(containers, checker.HasLen, startCount+1)
	actual := containers[0].Names[0]
	c.Assert(actual, checker.Equals, "/"+name)
***REMOVED***

// regression test for empty json field being omitted #13691
func (s *DockerSuite) TestContainerAPIGetJSONNoFieldsOmitted(c *check.C) ***REMOVED***
	startCount := getContainerCount(c)
	dockerCmd(c, "run", "busybox", "true")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.ContainerListOptions***REMOVED***
		All: true,
	***REMOVED***
	containers, err := cli.ContainerList(context.Background(), options)
	c.Assert(err, checker.IsNil)
	c.Assert(containers, checker.HasLen, startCount+1)
	actual := fmt.Sprintf("%+v", containers[0])

	// empty Labels field triggered this bug, make sense to check for everything
	// cause even Ports for instance can trigger this bug
	// better safe than sorry..
	fields := []string***REMOVED***
		"ID",
		"Names",
		"Image",
		"Command",
		"Created",
		"Ports",
		"Labels",
		"Status",
		"NetworkSettings",
	***REMOVED***

	// decoding into types.Container do not work since it eventually unmarshal
	// and empty field to an empty go map, so we just check for a string
	for _, f := range fields ***REMOVED***
		if !strings.Contains(actual, f) ***REMOVED***
			c.Fatalf("Field %s is missing and it shouldn't", f)
		***REMOVED***
	***REMOVED***
***REMOVED***

type containerPs struct ***REMOVED***
	Names []string
	Ports []types.Port
***REMOVED***

// regression test for non-empty fields from #13901
func (s *DockerSuite) TestContainerAPIPsOmitFields(c *check.C) ***REMOVED***
	// Problematic for Windows porting due to networking not yet being passed back
	testRequires(c, DaemonIsLinux)
	name := "pstest"
	port := 80
	runSleepingContainer(c, "--name", name, "--expose", strconv.Itoa(port))

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.ContainerListOptions***REMOVED***
		All: true,
	***REMOVED***
	containers, err := cli.ContainerList(context.Background(), options)
	c.Assert(err, checker.IsNil)
	var foundContainer containerPs
	for _, c := range containers ***REMOVED***
		for _, testName := range c.Names ***REMOVED***
			if "/"+name == testName ***REMOVED***
				foundContainer.Names = c.Names
				foundContainer.Ports = c.Ports
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	c.Assert(foundContainer.Ports, checker.HasLen, 1)
	c.Assert(foundContainer.Ports[0].PrivatePort, checker.Equals, uint16(port))
	c.Assert(foundContainer.Ports[0].PublicPort, checker.NotNil)
	c.Assert(foundContainer.Ports[0].IP, checker.NotNil)
***REMOVED***

func (s *DockerSuite) TestContainerAPIGetExport(c *check.C) ***REMOVED***
	// Not supported on Windows as Windows does not support docker export
	testRequires(c, DaemonIsLinux)
	name := "exportcontainer"
	dockerCmd(c, "run", "--name", name, "busybox", "touch", "/test")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	body, err := cli.ContainerExport(context.Background(), name)
	c.Assert(err, checker.IsNil)
	defer body.Close()
	found := false
	for tarReader := tar.NewReader(body); ; ***REMOVED***
		h, err := tarReader.Next()
		if err != nil && err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if h.Name == "test" ***REMOVED***
			found = true
			break
		***REMOVED***
	***REMOVED***
	c.Assert(found, checker.True, check.Commentf("The created test file has not been found in the exported image"))
***REMOVED***

func (s *DockerSuite) TestContainerAPIGetChanges(c *check.C) ***REMOVED***
	// Not supported on Windows as Windows does not support docker diff (/containers/name/changes)
	testRequires(c, DaemonIsLinux)
	name := "changescontainer"
	dockerCmd(c, "run", "--name", name, "busybox", "rm", "/etc/passwd")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	changes, err := cli.ContainerDiff(context.Background(), name)
	c.Assert(err, checker.IsNil)

	// Check the changelog for removal of /etc/passwd
	success := false
	for _, elem := range changes ***REMOVED***
		if elem.Path == "/etc/passwd" && elem.Kind == 2 ***REMOVED***
			success = true
		***REMOVED***
	***REMOVED***
	c.Assert(success, checker.True, check.Commentf("/etc/passwd has been removed but is not present in the diff"))
***REMOVED***

func (s *DockerSuite) TestGetContainerStats(c *check.C) ***REMOVED***
	var (
		name = "statscontainer"
	)
	runSleepingContainer(c, "--name", name)

	type b struct ***REMOVED***
		stats types.ContainerStats
		err   error
	***REMOVED***

	bc := make(chan b, 1)
	go func() ***REMOVED***
		cli, err := client.NewEnvClient()
		c.Assert(err, checker.IsNil)
		defer cli.Close()

		stats, err := cli.ContainerStats(context.Background(), name, true)
		c.Assert(err, checker.IsNil)
		bc <- b***REMOVED***stats, err***REMOVED***
	***REMOVED***()

	// allow some time to stream the stats from the container
	time.Sleep(4 * time.Second)
	dockerCmd(c, "rm", "-f", name)

	// collect the results from the stats stream or timeout and fail
	// if the stream was not disconnected.
	select ***REMOVED***
	case <-time.After(2 * time.Second):
		c.Fatal("stream was not closed after container was removed")
	case sr := <-bc:
		dec := json.NewDecoder(sr.stats.Body)
		defer sr.stats.Body.Close()
		var s *types.Stats
		// decode only one object from the stream
		c.Assert(dec.Decode(&s), checker.IsNil)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestGetContainerStatsRmRunning(c *check.C) ***REMOVED***
	out := runSleepingContainer(c)
	id := strings.TrimSpace(out)

	buf := &ChannelBuffer***REMOVED***C: make(chan []byte, 1)***REMOVED***
	defer buf.Close()

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	stats, err := cli.ContainerStats(context.Background(), id, true)
	c.Assert(err, checker.IsNil)
	defer stats.Body.Close()

	chErr := make(chan error, 1)
	go func() ***REMOVED***
		_, err = io.Copy(buf, stats.Body)
		chErr <- err
	***REMOVED***()

	b := make([]byte, 32)
	// make sure we've got some stats
	_, err = buf.ReadTimeout(b, 2*time.Second)
	c.Assert(err, checker.IsNil)

	// Now remove without `-f` and make sure we are still pulling stats
	_, _, err = dockerCmdWithError("rm", id)
	c.Assert(err, checker.Not(checker.IsNil), check.Commentf("rm should have failed but didn't"))
	_, err = buf.ReadTimeout(b, 2*time.Second)
	c.Assert(err, checker.IsNil)

	dockerCmd(c, "rm", "-f", id)
	c.Assert(<-chErr, checker.IsNil)
***REMOVED***

// ChannelBuffer holds a chan of byte array that can be populate in a goroutine.
type ChannelBuffer struct ***REMOVED***
	C chan []byte
***REMOVED***

// Write implements Writer.
func (c *ChannelBuffer) Write(b []byte) (int, error) ***REMOVED***
	c.C <- b
	return len(b), nil
***REMOVED***

// Close closes the go channel.
func (c *ChannelBuffer) Close() error ***REMOVED***
	close(c.C)
	return nil
***REMOVED***

// ReadTimeout reads the content of the channel in the specified byte array with
// the specified duration as timeout.
func (c *ChannelBuffer) ReadTimeout(p []byte, n time.Duration) (int, error) ***REMOVED***
	select ***REMOVED***
	case b := <-c.C:
		return copy(p[0:], b), nil
	case <-time.After(n):
		return -1, fmt.Errorf("timeout reading from channel")
	***REMOVED***
***REMOVED***

// regression test for gh13421
// previous test was just checking one stat entry so it didn't fail (stats with
// stream false always return one stat)
func (s *DockerSuite) TestGetContainerStatsStream(c *check.C) ***REMOVED***
	name := "statscontainer"
	runSleepingContainer(c, "--name", name)

	type b struct ***REMOVED***
		stats types.ContainerStats
		err   error
	***REMOVED***

	bc := make(chan b, 1)
	go func() ***REMOVED***
		cli, err := client.NewEnvClient()
		c.Assert(err, checker.IsNil)
		defer cli.Close()

		stats, err := cli.ContainerStats(context.Background(), name, true)
		c.Assert(err, checker.IsNil)
		bc <- b***REMOVED***stats, err***REMOVED***
	***REMOVED***()

	// allow some time to stream the stats from the container
	time.Sleep(4 * time.Second)
	dockerCmd(c, "rm", "-f", name)

	// collect the results from the stats stream or timeout and fail
	// if the stream was not disconnected.
	select ***REMOVED***
	case <-time.After(2 * time.Second):
		c.Fatal("stream was not closed after container was removed")
	case sr := <-bc:
		b, err := ioutil.ReadAll(sr.stats.Body)
		defer sr.stats.Body.Close()
		c.Assert(err, checker.IsNil)
		s := string(b)
		// count occurrences of "read" of types.Stats
		if l := strings.Count(s, "read"); l < 2 ***REMOVED***
			c.Fatalf("Expected more than one stat streamed, got %d", l)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestGetContainerStatsNoStream(c *check.C) ***REMOVED***
	name := "statscontainer"
	runSleepingContainer(c, "--name", name)

	type b struct ***REMOVED***
		stats types.ContainerStats
		err   error
	***REMOVED***

	bc := make(chan b, 1)

	go func() ***REMOVED***
		cli, err := client.NewEnvClient()
		c.Assert(err, checker.IsNil)
		defer cli.Close()

		stats, err := cli.ContainerStats(context.Background(), name, false)
		c.Assert(err, checker.IsNil)
		bc <- b***REMOVED***stats, err***REMOVED***
	***REMOVED***()

	// allow some time to stream the stats from the container
	time.Sleep(4 * time.Second)
	dockerCmd(c, "rm", "-f", name)

	// collect the results from the stats stream or timeout and fail
	// if the stream was not disconnected.
	select ***REMOVED***
	case <-time.After(2 * time.Second):
		c.Fatal("stream was not closed after container was removed")
	case sr := <-bc:
		b, err := ioutil.ReadAll(sr.stats.Body)
		defer sr.stats.Body.Close()
		c.Assert(err, checker.IsNil)
		s := string(b)
		// count occurrences of `"read"` of types.Stats
		c.Assert(strings.Count(s, `"read"`), checker.Equals, 1, check.Commentf("Expected only one stat streamed, got %d", strings.Count(s, `"read"`)))
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestGetStoppedContainerStats(c *check.C) ***REMOVED***
	name := "statscontainer"
	dockerCmd(c, "create", "--name", name, "busybox", "ps")

	chResp := make(chan error)

	// We expect an immediate response, but if it's not immediate, the test would hang, so put it in a goroutine
	// below we'll check this on a timeout.
	go func() ***REMOVED***
		cli, err := client.NewEnvClient()
		c.Assert(err, checker.IsNil)
		defer cli.Close()

		resp, err := cli.ContainerStats(context.Background(), name, false)
		defer resp.Body.Close()
		chResp <- err
	***REMOVED***()

	select ***REMOVED***
	case err := <-chResp:
		c.Assert(err, checker.IsNil)
	case <-time.After(10 * time.Second):
		c.Fatal("timeout waiting for stats response for stopped container")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestContainerAPIPause(c *check.C) ***REMOVED***
	// Problematic on Windows as Windows does not support pause
	testRequires(c, DaemonIsLinux)

	getPaused := func(c *check.C) []string ***REMOVED***
		return strings.Fields(cli.DockerCmd(c, "ps", "-f", "status=paused", "-q", "-a").Combined())
	***REMOVED***

	out := cli.DockerCmd(c, "run", "-d", "busybox", "sleep", "30").Combined()
	ContainerID := strings.TrimSpace(out)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerPause(context.Background(), ContainerID)
	c.Assert(err, checker.IsNil)

	pausedContainers := getPaused(c)

	if len(pausedContainers) != 1 || stringid.TruncateID(ContainerID) != pausedContainers[0] ***REMOVED***
		c.Fatalf("there should be one paused container and not %d", len(pausedContainers))
	***REMOVED***

	err = cli.ContainerUnpause(context.Background(), ContainerID)
	c.Assert(err, checker.IsNil)

	pausedContainers = getPaused(c)
	c.Assert(pausedContainers, checker.HasLen, 0, check.Commentf("There should be no paused container."))
***REMOVED***

func (s *DockerSuite) TestContainerAPITop(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "top")
	id := strings.TrimSpace(string(out))
	c.Assert(waitRun(id), checker.IsNil)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	top, err := cli.ContainerTop(context.Background(), id, []string***REMOVED***"aux"***REMOVED***)
	c.Assert(err, checker.IsNil)
	c.Assert(top.Titles, checker.HasLen, 11, check.Commentf("expected 11 titles, found %d: %v", len(top.Titles), top.Titles))

	if top.Titles[0] != "USER" || top.Titles[10] != "COMMAND" ***REMOVED***
		c.Fatalf("expected `USER` at `Titles[0]` and `COMMAND` at Titles[10]: %v", top.Titles)
	***REMOVED***
	c.Assert(top.Processes, checker.HasLen, 2, check.Commentf("expected 2 processes, found %d: %v", len(top.Processes), top.Processes))
	c.Assert(top.Processes[0][10], checker.Equals, "/bin/sh -c top")
	c.Assert(top.Processes[1][10], checker.Equals, "top")
***REMOVED***

func (s *DockerSuite) TestContainerAPITopWindows(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsWindows)
	out := runSleepingContainer(c, "-d")
	id := strings.TrimSpace(string(out))
	c.Assert(waitRun(id), checker.IsNil)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	top, err := cli.ContainerTop(context.Background(), id, nil)
	c.Assert(err, checker.IsNil)
	c.Assert(top.Titles, checker.HasLen, 4, check.Commentf("expected 4 titles, found %d: %v", len(top.Titles), top.Titles))

	if top.Titles[0] != "Name" || top.Titles[3] != "Private Working Set" ***REMOVED***
		c.Fatalf("expected `Name` at `Titles[0]` and `Private Working Set` at Titles[3]: %v", top.Titles)
	***REMOVED***
	c.Assert(len(top.Processes), checker.GreaterOrEqualThan, 2, check.Commentf("expected at least 2 processes, found %d: %v", len(top.Processes), top.Processes))

	foundProcess := false
	expectedProcess := "busybox.exe"
	for _, process := range top.Processes ***REMOVED***
		if process[0] == expectedProcess ***REMOVED***
			foundProcess = true
			break
		***REMOVED***
	***REMOVED***

	c.Assert(foundProcess, checker.Equals, true, check.Commentf("expected to find %s: %v", expectedProcess, top.Processes))
***REMOVED***

func (s *DockerSuite) TestContainerAPICommit(c *check.C) ***REMOVED***
	cName := "testapicommit"
	dockerCmd(c, "run", "--name="+cName, "busybox", "/bin/sh", "-c", "touch /test")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	options := types.ContainerCommitOptions***REMOVED***
		Reference: "testcontainerapicommit:testtag",
	***REMOVED***

	img, err := cli.ContainerCommit(context.Background(), cName, options)
	c.Assert(err, checker.IsNil)

	cmd := inspectField(c, img.ID, "Config.Cmd")
	c.Assert(cmd, checker.Equals, "[/bin/sh -c touch /test]", check.Commentf("got wrong Cmd from commit: %q", cmd))

	// sanity check, make sure the image is what we think it is
	dockerCmd(c, "run", img.ID, "ls", "/test")
***REMOVED***

func (s *DockerSuite) TestContainerAPICommitWithLabelInConfig(c *check.C) ***REMOVED***
	cName := "testapicommitwithconfig"
	dockerCmd(c, "run", "--name="+cName, "busybox", "/bin/sh", "-c", "touch /test")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	config := containertypes.Config***REMOVED***
		Labels: map[string]string***REMOVED***"key1": "value1", "key2": "value2"***REMOVED******REMOVED***

	options := types.ContainerCommitOptions***REMOVED***
		Reference: "testcontainerapicommitwithconfig",
		Config:    &config,
	***REMOVED***

	img, err := cli.ContainerCommit(context.Background(), cName, options)
	c.Assert(err, checker.IsNil)

	label1 := inspectFieldMap(c, img.ID, "Config.Labels", "key1")
	c.Assert(label1, checker.Equals, "value1")

	label2 := inspectFieldMap(c, img.ID, "Config.Labels", "key2")
	c.Assert(label2, checker.Equals, "value2")

	cmd := inspectField(c, img.ID, "Config.Cmd")
	c.Assert(cmd, checker.Equals, "[/bin/sh -c touch /test]", check.Commentf("got wrong Cmd from commit: %q", cmd))

	// sanity check, make sure the image is what we think it is
	dockerCmd(c, "run", img.ID, "ls", "/test")
***REMOVED***

func (s *DockerSuite) TestContainerAPIBadPort(c *check.C) ***REMOVED***
	// TODO Windows to Windows CI - Port this test
	testRequires(c, DaemonIsLinux)

	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"/bin/sh", "-c", "echo test"***REMOVED***,
	***REMOVED***

	hostConfig := containertypes.HostConfig***REMOVED***
		PortBindings: nat.PortMap***REMOVED***
			"8080/tcp": []nat.PortBinding***REMOVED***
				***REMOVED***
					HostIP:   "",
					HostPort: "aa80"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err.Error(), checker.Contains, `invalid port specification: "aa80"`)
***REMOVED***

func (s *DockerSuite) TestContainerAPICreate(c *check.C) ***REMOVED***
	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"/bin/sh", "-c", "touch /test && ls /test"***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, checker.IsNil)

	out, _ := dockerCmd(c, "start", "-a", container.ID)
	c.Assert(strings.TrimSpace(out), checker.Equals, "/test")
***REMOVED***

func (s *DockerSuite) TestContainerAPICreateEmptyConfig(c *check.C) ***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &containertypes.Config***REMOVED******REMOVED***, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")

	expected := "No command specified"
	c.Assert(err.Error(), checker.Contains, expected)
***REMOVED***

func (s *DockerSuite) TestContainerAPICreateMultipleNetworksConfig(c *check.C) ***REMOVED***
	// Container creation must fail if client specified configurations for more than one network
	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***

	networkingConfig := networktypes.NetworkingConfig***REMOVED***
		EndpointsConfig: map[string]*networktypes.EndpointSettings***REMOVED***
			"net1": ***REMOVED******REMOVED***,
			"net2": ***REMOVED******REMOVED***,
			"net3": ***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networkingConfig, "")
	msg := err.Error()
	// network name order in error message is not deterministic
	c.Assert(msg, checker.Contains, "Container cannot be connected to network endpoints")
	c.Assert(msg, checker.Contains, "net1")
	c.Assert(msg, checker.Contains, "net2")
	c.Assert(msg, checker.Contains, "net3")
***REMOVED***

func (s *DockerSuite) TestContainerAPICreateWithHostName(c *check.C) ***REMOVED***
	domainName := "test-domain"
	hostName := "test-hostname"
	config := containertypes.Config***REMOVED***
		Image:      "busybox",
		Hostname:   hostName,
		Domainname: domainName,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, checker.IsNil)

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	c.Assert(err, checker.IsNil)

	c.Assert(containerJSON.Config.Hostname, checker.Equals, hostName, check.Commentf("Mismatched Hostname"))
	c.Assert(containerJSON.Config.Domainname, checker.Equals, domainName, check.Commentf("Mismatched Domainname"))
***REMOVED***

func (s *DockerSuite) TestContainerAPICreateBridgeNetworkMode(c *check.C) ***REMOVED***
	// Windows does not support bridge
	testRequires(c, DaemonIsLinux)
	UtilCreateNetworkMode(c, "bridge")
***REMOVED***

func (s *DockerSuite) TestContainerAPICreateOtherNetworkModes(c *check.C) ***REMOVED***
	// Windows does not support these network modes
	testRequires(c, DaemonIsLinux, NotUserNamespace)
	UtilCreateNetworkMode(c, "host")
	UtilCreateNetworkMode(c, "container:web1")
***REMOVED***

func UtilCreateNetworkMode(c *check.C, networkMode containertypes.NetworkMode) ***REMOVED***
	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***

	hostConfig := containertypes.HostConfig***REMOVED***
		NetworkMode: networkMode,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, checker.IsNil)

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	c.Assert(err, checker.IsNil)

	c.Assert(containerJSON.HostConfig.NetworkMode, checker.Equals, containertypes.NetworkMode(networkMode), check.Commentf("Mismatched NetworkMode"))
***REMOVED***

func (s *DockerSuite) TestContainerAPICreateWithCpuSharesCpuset(c *check.C) ***REMOVED***
	// TODO Windows to Windows CI. The CpuShares part could be ported.
	testRequires(c, DaemonIsLinux)
	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***

	hostConfig := containertypes.HostConfig***REMOVED***
		Resources: containertypes.Resources***REMOVED***
			CPUShares:  512,
			CpusetCpus: "0",
		***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, checker.IsNil)

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	c.Assert(err, checker.IsNil)

	out := inspectField(c, containerJSON.ID, "HostConfig.CpuShares")
	c.Assert(out, checker.Equals, "512")

	outCpuset := inspectField(c, containerJSON.ID, "HostConfig.CpusetCpus")
	c.Assert(outCpuset, checker.Equals, "0")
***REMOVED***

func (s *DockerSuite) TestContainerAPIVerifyHeader(c *check.C) ***REMOVED***
	config := map[string]interface***REMOVED******REMOVED******REMOVED***
		"Image": "busybox",
	***REMOVED***

	create := func(ct string) (*http.Response, io.ReadCloser, error) ***REMOVED***
		jsonData := bytes.NewBuffer(nil)
		c.Assert(json.NewEncoder(jsonData).Encode(config), checker.IsNil)
		return request.Post("/containers/create", request.RawContent(ioutil.NopCloser(jsonData)), request.ContentType(ct))
	***REMOVED***

	// Try with no content-type
	res, body, err := create("")
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)
	body.Close()

	// Try with wrong content-type
	res, body, err = create("application/xml")
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)
	body.Close()

	// now application/json
	res, body, err = create("application/json")
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusCreated)
	body.Close()
***REMOVED***

//Issue 14230. daemon should return 500 for invalid port syntax
func (s *DockerSuite) TestContainerAPIInvalidPortSyntax(c *check.C) ***REMOVED***
	config := `***REMOVED***
				  "Image": "busybox",
				  "HostConfig": ***REMOVED***
					"NetworkMode": "default",
					"PortBindings": ***REMOVED***
					  "19039;1230": [
						***REMOVED******REMOVED***
					  ]
					***REMOVED***
				  ***REMOVED***
				***REMOVED***`

	res, body, err := request.Post("/containers/create", request.RawString(config), request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)

	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(string(b[:]), checker.Contains, "invalid port")
***REMOVED***

func (s *DockerSuite) TestContainerAPIRestartPolicyInvalidPolicyName(c *check.C) ***REMOVED***
	config := `***REMOVED***
		"Image": "busybox",
		"HostConfig": ***REMOVED***
			"RestartPolicy": ***REMOVED***
				"Name": "something",
				"MaximumRetryCount": 0
			***REMOVED***
		***REMOVED***
	***REMOVED***`

	res, body, err := request.Post("/containers/create", request.RawString(config), request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)

	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(string(b[:]), checker.Contains, "invalid restart policy")
***REMOVED***

func (s *DockerSuite) TestContainerAPIRestartPolicyRetryMismatch(c *check.C) ***REMOVED***
	config := `***REMOVED***
		"Image": "busybox",
		"HostConfig": ***REMOVED***
			"RestartPolicy": ***REMOVED***
				"Name": "always",
				"MaximumRetryCount": 2
			***REMOVED***
		***REMOVED***
	***REMOVED***`

	res, body, err := request.Post("/containers/create", request.RawString(config), request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)

	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(string(b[:]), checker.Contains, "maximum retry count cannot be used with restart policy")
***REMOVED***

func (s *DockerSuite) TestContainerAPIRestartPolicyNegativeRetryCount(c *check.C) ***REMOVED***
	config := `***REMOVED***
		"Image": "busybox",
		"HostConfig": ***REMOVED***
			"RestartPolicy": ***REMOVED***
				"Name": "on-failure",
				"MaximumRetryCount": -2
			***REMOVED***
		***REMOVED***
	***REMOVED***`

	res, body, err := request.Post("/containers/create", request.RawString(config), request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)

	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(string(b[:]), checker.Contains, "maximum retry count cannot be negative")
***REMOVED***

func (s *DockerSuite) TestContainerAPIRestartPolicyDefaultRetryCount(c *check.C) ***REMOVED***
	config := `***REMOVED***
		"Image": "busybox",
		"HostConfig": ***REMOVED***
			"RestartPolicy": ***REMOVED***
				"Name": "on-failure",
				"MaximumRetryCount": 0
			***REMOVED***
		***REMOVED***
	***REMOVED***`

	res, _, err := request.Post("/containers/create", request.RawString(config), request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusCreated)
***REMOVED***

// Issue 7941 - test to make sure a "null" in JSON is just ignored.
// W/o this fix a null in JSON would be parsed into a string var as "null"
func (s *DockerSuite) TestContainerAPIPostCreateNull(c *check.C) ***REMOVED***
	config := `***REMOVED***
		"Hostname":"",
		"Domainname":"",
		"Memory":0,
		"MemorySwap":0,
		"CpuShares":0,
		"Cpuset":null,
		"AttachStdin":true,
		"AttachStdout":true,
		"AttachStderr":true,
		"ExposedPorts":***REMOVED******REMOVED***,
		"Tty":true,
		"OpenStdin":true,
		"StdinOnce":true,
		"Env":[],
		"Cmd":"ls",
		"Image":"busybox",
		"Volumes":***REMOVED******REMOVED***,
		"WorkingDir":"",
		"Entrypoint":null,
		"NetworkDisabled":false,
		"OnBuild":null***REMOVED***`

	res, body, err := request.Post("/containers/create", request.RawString(config), request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusCreated)

	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	type createResp struct ***REMOVED***
		ID string
	***REMOVED***
	var container createResp
	c.Assert(json.Unmarshal(b, &container), checker.IsNil)
	out := inspectField(c, container.ID, "HostConfig.CpusetCpus")
	c.Assert(out, checker.Equals, "")

	outMemory := inspectField(c, container.ID, "HostConfig.Memory")
	c.Assert(outMemory, checker.Equals, "0")
	outMemorySwap := inspectField(c, container.ID, "HostConfig.MemorySwap")
	c.Assert(outMemorySwap, checker.Equals, "0")
***REMOVED***

func (s *DockerSuite) TestCreateWithTooLowMemoryLimit(c *check.C) ***REMOVED***
	// TODO Windows: Port once memory is supported
	testRequires(c, DaemonIsLinux)
	config := `***REMOVED***
		"Image":     "busybox",
		"Cmd":       "ls",
		"OpenStdin": true,
		"CpuShares": 100,
		"Memory":    524287
	***REMOVED***`

	res, body, err := request.Post("/containers/create", request.RawString(config), request.JSON)
	c.Assert(err, checker.IsNil)
	b, err2 := request.ReadBody(body)
	c.Assert(err2, checker.IsNil)

	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)
	c.Assert(string(b), checker.Contains, "Minimum memory limit allowed is 4MB")
***REMOVED***

func (s *DockerSuite) TestContainerAPIRename(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "--name", "TestContainerAPIRename", "-d", "busybox", "sh")

	containerID := strings.TrimSpace(out)
	newName := "TestContainerAPIRenameNew"

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRename(context.Background(), containerID, newName)
	c.Assert(err, checker.IsNil)

	name := inspectField(c, containerID, "Name")
	c.Assert(name, checker.Equals, "/"+newName, check.Commentf("Failed to rename container"))
***REMOVED***

func (s *DockerSuite) TestContainerAPIKill(c *check.C) ***REMOVED***
	name := "test-api-kill"
	runSleepingContainer(c, "-i", "--name", name)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerKill(context.Background(), name, "SIGKILL")
	c.Assert(err, checker.IsNil)

	state := inspectField(c, name, "State.Running")
	c.Assert(state, checker.Equals, "false", check.Commentf("got wrong State from container %s: %q", name, state))
***REMOVED***

func (s *DockerSuite) TestContainerAPIRestart(c *check.C) ***REMOVED***
	name := "test-api-restart"
	runSleepingContainer(c, "-di", "--name", name)
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	timeout := 1 * time.Second
	err = cli.ContainerRestart(context.Background(), name, &timeout)
	c.Assert(err, checker.IsNil)

	c.Assert(waitInspect(name, "***REMOVED******REMOVED*** .State.Restarting  ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .State.Running  ***REMOVED******REMOVED***", "false true", 15*time.Second), checker.IsNil)
***REMOVED***

func (s *DockerSuite) TestContainerAPIRestartNotimeoutParam(c *check.C) ***REMOVED***
	name := "test-api-restart-no-timeout-param"
	out := runSleepingContainer(c, "-di", "--name", name)
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRestart(context.Background(), name, nil)
	c.Assert(err, checker.IsNil)

	c.Assert(waitInspect(name, "***REMOVED******REMOVED*** .State.Restarting  ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .State.Running  ***REMOVED******REMOVED***", "false true", 15*time.Second), checker.IsNil)
***REMOVED***

func (s *DockerSuite) TestContainerAPIStart(c *check.C) ***REMOVED***
	name := "testing-start"
	config := containertypes.Config***REMOVED***
		Image:     "busybox",
		Cmd:       append([]string***REMOVED***"/bin/sh", "-c"***REMOVED***, sleepCommandForDaemonPlatform()...),
		OpenStdin: true,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, name)
	c.Assert(err, checker.IsNil)

	err = cli.ContainerStart(context.Background(), name, types.ContainerStartOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)

	// second call to start should give 304
	// maybe add ContainerStartWithRaw to test it
	err = cli.ContainerStart(context.Background(), name, types.ContainerStartOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)

	// TODO(tibor): figure out why this doesn't work on windows
***REMOVED***

func (s *DockerSuite) TestContainerAPIStop(c *check.C) ***REMOVED***
	name := "test-api-stop"
	runSleepingContainer(c, "-i", "--name", name)
	timeout := 30 * time.Second

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerStop(context.Background(), name, &timeout)
	c.Assert(err, checker.IsNil)
	c.Assert(waitInspect(name, "***REMOVED******REMOVED*** .State.Running  ***REMOVED******REMOVED***", "false", 60*time.Second), checker.IsNil)

	// second call to start should give 304
	// maybe add ContainerStartWithRaw to test it
	err = cli.ContainerStop(context.Background(), name, &timeout)
	c.Assert(err, checker.IsNil)
***REMOVED***

func (s *DockerSuite) TestContainerAPIWait(c *check.C) ***REMOVED***
	name := "test-api-wait"

	sleepCmd := "/bin/sleep"
	if testEnv.OSType == "windows" ***REMOVED***
		sleepCmd = "sleep"
	***REMOVED***
	dockerCmd(c, "run", "--name", name, "busybox", sleepCmd, "2")

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	waitresC, errC := cli.ContainerWait(context.Background(), name, "")

	select ***REMOVED***
	case err = <-errC:
		c.Assert(err, checker.IsNil)
	case waitres := <-waitresC:
		c.Assert(waitres.StatusCode, checker.Equals, int64(0))
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestContainerAPICopyNotExistsAnyMore(c *check.C) ***REMOVED***
	name := "test-container-api-copy"
	dockerCmd(c, "run", "--name", name, "busybox", "touch", "/test.txt")

	postData := types.CopyConfig***REMOVED***
		Resource: "/test.txt",
	***REMOVED***
	// no copy in client/
	res, _, err := request.Post("/containers/"+name+"/copy", request.JSONBody(postData))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusNotFound)
***REMOVED***

func (s *DockerSuite) TestContainerAPICopyPre124(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // Windows only supports 1.25 or later
	name := "test-container-api-copy"
	dockerCmd(c, "run", "--name", name, "busybox", "touch", "/test.txt")

	postData := types.CopyConfig***REMOVED***
		Resource: "/test.txt",
	***REMOVED***

	res, body, err := request.Post("/v1.23/containers/"+name+"/copy", request.JSONBody(postData))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	found := false
	for tarReader := tar.NewReader(body); ; ***REMOVED***
		h, err := tarReader.Next()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			c.Fatal(err)
		***REMOVED***
		if h.Name == "test.txt" ***REMOVED***
			found = true
			break
		***REMOVED***
	***REMOVED***
	c.Assert(found, checker.True)
***REMOVED***

func (s *DockerSuite) TestContainerAPICopyResourcePathEmptyPre124(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // Windows only supports 1.25 or later
	name := "test-container-api-copy-resource-empty"
	dockerCmd(c, "run", "--name", name, "busybox", "touch", "/test.txt")

	postData := types.CopyConfig***REMOVED***
		Resource: "",
	***REMOVED***

	res, body, err := request.Post("/v1.23/containers/"+name+"/copy", request.JSONBody(postData))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)
	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(string(b), checker.Matches, "Path cannot be empty\n")
***REMOVED***

func (s *DockerSuite) TestContainerAPICopyResourcePathNotFoundPre124(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // Windows only supports 1.25 or later
	name := "test-container-api-copy-resource-not-found"
	dockerCmd(c, "run", "--name", name, "busybox")

	postData := types.CopyConfig***REMOVED***
		Resource: "/notexist",
	***REMOVED***

	res, body, err := request.Post("/v1.23/containers/"+name+"/copy", request.JSONBody(postData))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusNotFound)

	b, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(string(b), checker.Matches, "Could not find the file /notexist in container "+name+"\n")
***REMOVED***

func (s *DockerSuite) TestContainerAPICopyContainerNotFoundPr124(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // Windows only supports 1.25 or later
	postData := types.CopyConfig***REMOVED***
		Resource: "/something",
	***REMOVED***

	res, _, err := request.Post("/v1.23/containers/notexists/copy", request.JSONBody(postData))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusNotFound)
***REMOVED***

func (s *DockerSuite) TestContainerAPIDelete(c *check.C) ***REMOVED***
	out := runSleepingContainer(c)

	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	dockerCmd(c, "stop", id)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
***REMOVED***

func (s *DockerSuite) TestContainerAPIDeleteNotExist(c *check.C) ***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRemove(context.Background(), "doesnotexist", types.ContainerRemoveOptions***REMOVED******REMOVED***)
	c.Assert(err.Error(), checker.Contains, "No such container: doesnotexist")
***REMOVED***

func (s *DockerSuite) TestContainerAPIDeleteForce(c *check.C) ***REMOVED***
	out := runSleepingContainer(c)
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	removeOptions := types.ContainerRemoveOptions***REMOVED***
		Force: true,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRemove(context.Background(), id, removeOptions)
	c.Assert(err, checker.IsNil)
***REMOVED***

func (s *DockerSuite) TestContainerAPIDeleteRemoveLinks(c *check.C) ***REMOVED***
	// Windows does not support links
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "--name", "tlink1", "busybox", "top")

	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	out, _ = dockerCmd(c, "run", "--link", "tlink1:tlink1", "--name", "tlink2", "-d", "busybox", "top")

	id2 := strings.TrimSpace(out)
	c.Assert(waitRun(id2), checker.IsNil)

	links := inspectFieldJSON(c, id2, "HostConfig.Links")
	c.Assert(links, checker.Equals, "[\"/tlink1:/tlink2/tlink1\"]", check.Commentf("expected to have links between containers"))

	removeOptions := types.ContainerRemoveOptions***REMOVED***
		RemoveLinks: true,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRemove(context.Background(), "tlink2/tlink1", removeOptions)
	c.Assert(err, check.IsNil)

	linksPostRm := inspectFieldJSON(c, id2, "HostConfig.Links")
	c.Assert(linksPostRm, checker.Equals, "null", check.Commentf("call to api deleteContainer links should have removed the specified links"))
***REMOVED***

func (s *DockerSuite) TestContainerAPIDeleteConflict(c *check.C) ***REMOVED***
	out := runSleepingContainer(c)

	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions***REMOVED******REMOVED***)
	expected := "cannot remove a running container"
	c.Assert(err.Error(), checker.Contains, expected)
***REMOVED***

func (s *DockerSuite) TestContainerAPIDeleteRemoveVolume(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon)

	vol := "/testvolume"
	if testEnv.OSType == "windows" ***REMOVED***
		vol = `c:\testvolume`
	***REMOVED***

	out := runSleepingContainer(c, "-v", vol)

	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	source, err := inspectMountSourceField(id, vol)
	_, err = os.Stat(source)
	c.Assert(err, checker.IsNil)

	removeOptions := types.ContainerRemoveOptions***REMOVED***
		Force:         true,
		RemoveVolumes: true,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRemove(context.Background(), id, removeOptions)
	c.Assert(err, check.IsNil)

	_, err = os.Stat(source)
	c.Assert(os.IsNotExist(err), checker.True, check.Commentf("expected to get ErrNotExist error, got %v", err))
***REMOVED***

// Regression test for https://github.com/docker/docker/issues/6231
func (s *DockerSuite) TestContainerAPIChunkedEncoding(c *check.C) ***REMOVED***

	config := map[string]interface***REMOVED******REMOVED******REMOVED***
		"Image":     "busybox",
		"Cmd":       append([]string***REMOVED***"/bin/sh", "-c"***REMOVED***, sleepCommandForDaemonPlatform()...),
		"OpenStdin": true,
	***REMOVED***

	resp, _, err := request.Post("/containers/create", request.JSONBody(config), func(req *http.Request) error ***REMOVED***
		// This is a cheat to make the http request do chunked encoding
		// Otherwise (just setting the Content-Encoding to chunked) net/http will overwrite
		// https://golang.org/src/pkg/net/http/request.go?s=11980:12172
		req.ContentLength = -1
		return nil
	***REMOVED***)
	c.Assert(err, checker.IsNil, check.Commentf("error creating container with chunked encoding"))
	defer resp.Body.Close()
	c.Assert(resp.StatusCode, checker.Equals, http.StatusCreated)
***REMOVED***

func (s *DockerSuite) TestContainerAPIPostContainerStop(c *check.C) ***REMOVED***
	out := runSleepingContainer(c)

	containerID := strings.TrimSpace(out)
	c.Assert(waitRun(containerID), checker.IsNil)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerStop(context.Background(), containerID, nil)
	c.Assert(err, checker.IsNil)
	c.Assert(waitInspect(containerID, "***REMOVED******REMOVED*** .State.Running  ***REMOVED******REMOVED***", "false", 60*time.Second), checker.IsNil)
***REMOVED***

// #14170
func (s *DockerSuite) TestPostContainerAPICreateWithStringOrSliceEntrypoint(c *check.C) ***REMOVED***
	config := containertypes.Config***REMOVED***
		Image:      "busybox",
		Entrypoint: []string***REMOVED***"echo"***REMOVED***,
		Cmd:        []string***REMOVED***"hello", "world"***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "echotest")
	c.Assert(err, checker.IsNil)
	out, _ := dockerCmd(c, "start", "-a", "echotest")
	c.Assert(strings.TrimSpace(out), checker.Equals, "hello world")

	config2 := struct ***REMOVED***
		Image      string
		Entrypoint string
		Cmd        []string
	***REMOVED******REMOVED***"busybox", "echo", []string***REMOVED***"hello", "world"***REMOVED******REMOVED***
	_, _, err = request.Post("/containers/create?name=echotest2", request.JSONBody(config2))
	c.Assert(err, checker.IsNil)
	out, _ = dockerCmd(c, "start", "-a", "echotest2")
	c.Assert(strings.TrimSpace(out), checker.Equals, "hello world")
***REMOVED***

// #14170
func (s *DockerSuite) TestPostContainersCreateWithStringOrSliceCmd(c *check.C) ***REMOVED***
	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"echo", "hello", "world"***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "echotest")
	c.Assert(err, checker.IsNil)
	out, _ := dockerCmd(c, "start", "-a", "echotest")
	c.Assert(strings.TrimSpace(out), checker.Equals, "hello world")

	config2 := struct ***REMOVED***
		Image      string
		Entrypoint string
		Cmd        string
	***REMOVED******REMOVED***"busybox", "echo", "hello world"***REMOVED***
	_, _, err = request.Post("/containers/create?name=echotest2", request.JSONBody(config2))
	c.Assert(err, checker.IsNil)
	out, _ = dockerCmd(c, "start", "-a", "echotest2")
	c.Assert(strings.TrimSpace(out), checker.Equals, "hello world")
***REMOVED***

// regression #14318
func (s *DockerSuite) TestPostContainersCreateWithStringOrSliceCapAddDrop(c *check.C) ***REMOVED***
	// Windows doesn't support CapAdd/CapDrop
	testRequires(c, DaemonIsLinux)
	config := struct ***REMOVED***
		Image   string
		CapAdd  string
		CapDrop string
	***REMOVED******REMOVED***"busybox", "NET_ADMIN", "SYS_ADMIN"***REMOVED***
	res, _, err := request.Post("/containers/create?name=capaddtest0", request.JSONBody(config))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusCreated)

	config2 := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***
	hostConfig := containertypes.HostConfig***REMOVED***
		CapAdd:  []string***REMOVED***"NET_ADMIN", "SYS_ADMIN"***REMOVED***,
		CapDrop: []string***REMOVED***"SETGID"***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config2, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "capaddtest1")
	c.Assert(err, checker.IsNil)
***REMOVED***

// #14915
func (s *DockerSuite) TestContainerAPICreateNoHostConfig118(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux) // Windows only support 1.25 or later
	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***

	cli, err := request.NewEnvClientWithVersion("v1.18")

	_, err = cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, checker.IsNil)
***REMOVED***

// Ensure an error occurs when you have a container read-only rootfs but you
// extract an archive to a symlink in a writable volume which points to a
// directory outside of the volume.
func (s *DockerSuite) TestPutContainerArchiveErrSymlinkInVolumeToReadOnlyRootfs(c *check.C) ***REMOVED***
	// Windows does not support read-only rootfs
	// Requires local volume mount bind.
	// --read-only + userns has remount issues
	testRequires(c, SameHostDaemon, NotUserNamespace, DaemonIsLinux)

	testVol := getTestDir(c, "test-put-container-archive-err-symlink-in-volume-to-read-only-rootfs-")
	defer os.RemoveAll(testVol)

	makeTestContentInDir(c, testVol)

	cID := makeTestContainer(c, testContainerOptions***REMOVED***
		readOnly: true,
		volumes:  defaultVolumes(testVol), // Our bind mount is at /vol2
	***REMOVED***)

	// Attempt to extract to a symlink in the volume which points to a
	// directory outside the volume. This should cause an error because the
	// rootfs is read-only.
	var httpClient *http.Client
	cli, err := client.NewClient(daemonHost(), "v1.20", httpClient, map[string]string***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)

	err = cli.CopyToContainer(context.Background(), cID, "/vol2/symlinkToAbsDir", nil, types.CopyToContainerOptions***REMOVED******REMOVED***)
	c.Assert(err.Error(), checker.Contains, "container rootfs is marked read-only")
***REMOVED***

func (s *DockerSuite) TestPostContainersCreateWithWrongCpusetValues(c *check.C) ***REMOVED***
	// Not supported on Windows
	testRequires(c, DaemonIsLinux)

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***
	hostConfig1 := containertypes.HostConfig***REMOVED***
		Resources: containertypes.Resources***REMOVED***
			CpusetCpus: "1-42,,",
		***REMOVED***,
	***REMOVED***
	name := "wrong-cpuset-cpus"

	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig1, &networktypes.NetworkingConfig***REMOVED******REMOVED***, name)
	expected := "Invalid value 1-42,, for cpuset cpus"
	c.Assert(err.Error(), checker.Contains, expected)

	hostConfig2 := containertypes.HostConfig***REMOVED***
		Resources: containertypes.Resources***REMOVED***
			CpusetMems: "42-3,1--",
		***REMOVED***,
	***REMOVED***
	name = "wrong-cpuset-mems"
	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig2, &networktypes.NetworkingConfig***REMOVED******REMOVED***, name)
	expected = "Invalid value 42-3,1-- for cpuset mems"
	c.Assert(err.Error(), checker.Contains, expected)
***REMOVED***

func (s *DockerSuite) TestPostContainersCreateShmSizeNegative(c *check.C) ***REMOVED***
	// ShmSize is not supported on Windows
	testRequires(c, DaemonIsLinux)
	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***
	hostConfig := containertypes.HostConfig***REMOVED***
		ShmSize: -1,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err.Error(), checker.Contains, "SHM size can not be less than 0")
***REMOVED***

func (s *DockerSuite) TestPostContainersCreateShmSizeHostConfigOmitted(c *check.C) ***REMOVED***
	// ShmSize is not supported on Windows
	testRequires(c, DaemonIsLinux)
	var defaultSHMSize int64 = 67108864
	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"mount"***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, check.IsNil)

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	c.Assert(err, check.IsNil)

	c.Assert(containerJSON.HostConfig.ShmSize, check.Equals, defaultSHMSize)

	out, _ := dockerCmd(c, "start", "-i", containerJSON.ID)
	shmRegexp := regexp.MustCompile(`shm on /dev/shm type tmpfs(.*)size=65536k`)
	if !shmRegexp.MatchString(out) ***REMOVED***
		c.Fatalf("Expected shm of 64MB in mount command, got %v", out)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestPostContainersCreateShmSizeOmitted(c *check.C) ***REMOVED***
	// ShmSize is not supported on Windows
	testRequires(c, DaemonIsLinux)
	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"mount"***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, check.IsNil)

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	c.Assert(err, check.IsNil)

	c.Assert(containerJSON.HostConfig.ShmSize, check.Equals, int64(67108864))

	out, _ := dockerCmd(c, "start", "-i", containerJSON.ID)
	shmRegexp := regexp.MustCompile(`shm on /dev/shm type tmpfs(.*)size=65536k`)
	if !shmRegexp.MatchString(out) ***REMOVED***
		c.Fatalf("Expected shm of 64MB in mount command, got %v", out)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestPostContainersCreateWithShmSize(c *check.C) ***REMOVED***
	// ShmSize is not supported on Windows
	testRequires(c, DaemonIsLinux)
	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"mount"***REMOVED***,
	***REMOVED***

	hostConfig := containertypes.HostConfig***REMOVED***
		ShmSize: 1073741824,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, check.IsNil)

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	c.Assert(err, check.IsNil)

	c.Assert(containerJSON.HostConfig.ShmSize, check.Equals, int64(1073741824))

	out, _ := dockerCmd(c, "start", "-i", containerJSON.ID)
	shmRegex := regexp.MustCompile(`shm on /dev/shm type tmpfs(.*)size=1048576k`)
	if !shmRegex.MatchString(out) ***REMOVED***
		c.Fatalf("Expected shm of 1GB in mount command, got %v", out)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestPostContainersCreateMemorySwappinessHostConfigOmitted(c *check.C) ***REMOVED***
	// Swappiness is not supported on Windows
	testRequires(c, DaemonIsLinux)
	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	container, err := cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
	c.Assert(err, check.IsNil)

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	c.Assert(err, check.IsNil)

	c.Assert(containerJSON.HostConfig.MemorySwappiness, check.IsNil)
***REMOVED***

// check validation is done daemon side and not only in cli
func (s *DockerSuite) TestPostContainersCreateWithOomScoreAdjInvalidRange(c *check.C) ***REMOVED***
	// OomScoreAdj is not supported on Windows
	testRequires(c, DaemonIsLinux)

	config := containertypes.Config***REMOVED***
		Image: "busybox",
	***REMOVED***

	hostConfig := containertypes.HostConfig***REMOVED***
		OomScoreAdj: 1001,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	name := "oomscoreadj-over"
	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, name)

	expected := "Invalid value 1001, range for oom score adj is [-1000, 1000]"
	c.Assert(err.Error(), checker.Contains, expected)

	hostConfig = containertypes.HostConfig***REMOVED***
		OomScoreAdj: -1001,
	***REMOVED***

	name = "oomscoreadj-low"
	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, name)

	expected = "Invalid value -1001, range for oom score adj is [-1000, 1000]"
	c.Assert(err.Error(), checker.Contains, expected)
***REMOVED***

// test case for #22210 where an empty container name caused panic.
func (s *DockerSuite) TestContainerAPIDeleteWithEmptyName(c *check.C) ***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	err = cli.ContainerRemove(context.Background(), "", types.ContainerRemoveOptions***REMOVED******REMOVED***)
	c.Assert(err.Error(), checker.Contains, "No such container")
***REMOVED***

func (s *DockerSuite) TestContainerAPIStatsWithNetworkDisabled(c *check.C) ***REMOVED***
	// Problematic on Windows as Windows does not support stats
	testRequires(c, DaemonIsLinux)

	name := "testing-network-disabled"

	config := containertypes.Config***REMOVED***
		Image:           "busybox",
		Cmd:             []string***REMOVED***"top"***REMOVED***,
		NetworkDisabled: true,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &containertypes.HostConfig***REMOVED******REMOVED***, &networktypes.NetworkingConfig***REMOVED******REMOVED***, name)
	c.Assert(err, checker.IsNil)

	err = cli.ContainerStart(context.Background(), name, types.ContainerStartOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)

	c.Assert(waitRun(name), check.IsNil)

	type b struct ***REMOVED***
		stats types.ContainerStats
		err   error
	***REMOVED***
	bc := make(chan b, 1)
	go func() ***REMOVED***
		stats, err := cli.ContainerStats(context.Background(), name, false)
		bc <- b***REMOVED***stats, err***REMOVED***
	***REMOVED***()

	// allow some time to stream the stats from the container
	time.Sleep(4 * time.Second)
	dockerCmd(c, "rm", "-f", name)

	// collect the results from the stats stream or timeout and fail
	// if the stream was not disconnected.
	select ***REMOVED***
	case <-time.After(2 * time.Second):
		c.Fatal("stream was not closed after container was removed")
	case sr := <-bc:
		c.Assert(sr.err, checker.IsNil)
		sr.stats.Body.Close()
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestContainersAPICreateMountsValidation(c *check.C) ***REMOVED***
	type testCase struct ***REMOVED***
		config     containertypes.Config
		hostConfig containertypes.HostConfig
		msg        string
	***REMOVED***

	prefix, slash := getPrefixAndSlashFromDaemonPlatform()
	destPath := prefix + slash + "foo"
	notExistPath := prefix + slash + "notexist"

	cases := []testCase***REMOVED***
		***REMOVED***
			config: containertypes.Config***REMOVED***
				Image: "busybox",
			***REMOVED***,
			hostConfig: containertypes.HostConfig***REMOVED***
				Mounts: []mounttypes.Mount***REMOVED******REMOVED***
					Type:   "notreal",
					Target: destPath,
				***REMOVED***,
				***REMOVED***,
			***REMOVED***,

			msg: "mount type unknown",
		***REMOVED***,
		***REMOVED***
			config: containertypes.Config***REMOVED***
				Image: "busybox",
			***REMOVED***,
			hostConfig: containertypes.HostConfig***REMOVED***
				Mounts: []mounttypes.Mount***REMOVED******REMOVED***
					Type: "bind"***REMOVED******REMOVED******REMOVED***,
			msg: "Target must not be empty",
		***REMOVED***,
		***REMOVED***
			config: containertypes.Config***REMOVED***
				Image: "busybox",
			***REMOVED***,
			hostConfig: containertypes.HostConfig***REMOVED***
				Mounts: []mounttypes.Mount***REMOVED******REMOVED***
					Type:   "bind",
					Target: destPath***REMOVED******REMOVED******REMOVED***,
			msg: "Source must not be empty",
		***REMOVED***,
		***REMOVED***
			config: containertypes.Config***REMOVED***
				Image: "busybox",
			***REMOVED***,
			hostConfig: containertypes.HostConfig***REMOVED***
				Mounts: []mounttypes.Mount***REMOVED******REMOVED***
					Type:   "bind",
					Source: notExistPath,
					Target: destPath***REMOVED******REMOVED******REMOVED***,
			msg: "bind source path does not exist",
		***REMOVED***,
		***REMOVED***
			config: containertypes.Config***REMOVED***
				Image: "busybox",
			***REMOVED***,
			hostConfig: containertypes.HostConfig***REMOVED***
				Mounts: []mounttypes.Mount***REMOVED******REMOVED***
					Type: "volume"***REMOVED******REMOVED******REMOVED***,
			msg: "Target must not be empty",
		***REMOVED***,
		***REMOVED***
			config: containertypes.Config***REMOVED***
				Image: "busybox",
			***REMOVED***,
			hostConfig: containertypes.HostConfig***REMOVED***
				Mounts: []mounttypes.Mount***REMOVED******REMOVED***
					Type:   "volume",
					Source: "hello",
					Target: destPath***REMOVED******REMOVED******REMOVED***,
			msg: "",
		***REMOVED***,
		***REMOVED***
			config: containertypes.Config***REMOVED***
				Image: "busybox",
			***REMOVED***,
			hostConfig: containertypes.HostConfig***REMOVED***
				Mounts: []mounttypes.Mount***REMOVED******REMOVED***
					Type:   "volume",
					Source: "hello2",
					Target: destPath,
					VolumeOptions: &mounttypes.VolumeOptions***REMOVED***
						DriverConfig: &mounttypes.Driver***REMOVED***
							Name: "local"***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
			msg: "",
		***REMOVED***,
	***REMOVED***

	if SameHostDaemon() ***REMOVED***
		tmpDir, err := ioutils.TempDir("", "test-mounts-api")
		c.Assert(err, checker.IsNil)
		defer os.RemoveAll(tmpDir)
		cases = append(cases, []testCase***REMOVED***
			***REMOVED***
				config: containertypes.Config***REMOVED***
					Image: "busybox",
				***REMOVED***,
				hostConfig: containertypes.HostConfig***REMOVED***
					Mounts: []mounttypes.Mount***REMOVED******REMOVED***
						Type:   "bind",
						Source: tmpDir,
						Target: destPath***REMOVED******REMOVED******REMOVED***,
				msg: "",
			***REMOVED***,
			***REMOVED***
				config: containertypes.Config***REMOVED***
					Image: "busybox",
				***REMOVED***,
				hostConfig: containertypes.HostConfig***REMOVED***
					Mounts: []mounttypes.Mount***REMOVED******REMOVED***
						Type:          "bind",
						Source:        tmpDir,
						Target:        destPath,
						VolumeOptions: &mounttypes.VolumeOptions***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
				msg: "VolumeOptions must not be specified",
			***REMOVED***,
		***REMOVED***...)
	***REMOVED***

	if DaemonIsLinux() ***REMOVED***
		cases = append(cases, []testCase***REMOVED***
			***REMOVED***
				config: containertypes.Config***REMOVED***
					Image: "busybox",
				***REMOVED***,
				hostConfig: containertypes.HostConfig***REMOVED***
					Mounts: []mounttypes.Mount***REMOVED******REMOVED***
						Type:   "volume",
						Source: "hello3",
						Target: destPath,
						VolumeOptions: &mounttypes.VolumeOptions***REMOVED***
							DriverConfig: &mounttypes.Driver***REMOVED***
								Name:    "local",
								Options: map[string]string***REMOVED***"o": "size=1"***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
				msg: "",
			***REMOVED***,
			***REMOVED***
				config: containertypes.Config***REMOVED***
					Image: "busybox",
				***REMOVED***,
				hostConfig: containertypes.HostConfig***REMOVED***
					Mounts: []mounttypes.Mount***REMOVED******REMOVED***
						Type:   "tmpfs",
						Target: destPath***REMOVED******REMOVED******REMOVED***,
				msg: "",
			***REMOVED***,
			***REMOVED***
				config: containertypes.Config***REMOVED***
					Image: "busybox",
				***REMOVED***,
				hostConfig: containertypes.HostConfig***REMOVED***
					Mounts: []mounttypes.Mount***REMOVED******REMOVED***
						Type:   "tmpfs",
						Target: destPath,
						TmpfsOptions: &mounttypes.TmpfsOptions***REMOVED***
							SizeBytes: 4096 * 1024,
							Mode:      0700,
						***REMOVED******REMOVED******REMOVED******REMOVED***,
				msg: "",
			***REMOVED***,

			***REMOVED***
				config: containertypes.Config***REMOVED***
					Image: "busybox",
				***REMOVED***,
				hostConfig: containertypes.HostConfig***REMOVED***
					Mounts: []mounttypes.Mount***REMOVED******REMOVED***
						Type:   "tmpfs",
						Source: "/shouldnotbespecified",
						Target: destPath***REMOVED******REMOVED******REMOVED***,
				msg: "Source must not be specified",
			***REMOVED***,
		***REMOVED***...)

	***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	for i, x := range cases ***REMOVED***
		c.Logf("case %d", i)
		_, err = cli.ContainerCreate(context.Background(), &x.config, &x.hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "")
		if len(x.msg) > 0 ***REMOVED***
			c.Assert(err.Error(), checker.Contains, x.msg, check.Commentf("%v", cases[i].config))
		***REMOVED*** else ***REMOVED***
			c.Assert(err, checker.IsNil)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestContainerAPICreateMountsBindRead(c *check.C) ***REMOVED***
	testRequires(c, NotUserNamespace, SameHostDaemon)
	// also with data in the host side
	prefix, slash := getPrefixAndSlashFromDaemonPlatform()
	destPath := prefix + slash + "foo"
	tmpDir, err := ioutil.TempDir("", "test-mounts-api-bind")
	c.Assert(err, checker.IsNil)
	defer os.RemoveAll(tmpDir)
	err = ioutil.WriteFile(filepath.Join(tmpDir, "bar"), []byte("hello"), 666)
	c.Assert(err, checker.IsNil)
	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"/bin/sh", "-c", "cat /foo/bar"***REMOVED***,
	***REMOVED***
	hostConfig := containertypes.HostConfig***REMOVED***
		Mounts: []mounttypes.Mount***REMOVED***
			***REMOVED***Type: "bind", Source: tmpDir, Target: destPath***REMOVED***,
		***REMOVED***,
	***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, "test")
	c.Assert(err, checker.IsNil)

	out, _ := dockerCmd(c, "start", "-a", "test")
	c.Assert(out, checker.Equals, "hello")
***REMOVED***

// Test Mounts comes out as expected for the MountPoint
func (s *DockerSuite) TestContainersAPICreateMountsCreate(c *check.C) ***REMOVED***
	prefix, slash := getPrefixAndSlashFromDaemonPlatform()
	destPath := prefix + slash + "foo"

	var (
		testImg string
	)
	if testEnv.OSType != "windows" ***REMOVED***
		testImg = "test-mount-config"
		buildImageSuccessfully(c, testImg, build.WithDockerfile(`
	FROM busybox
	RUN mkdir `+destPath+` && touch `+destPath+slash+`bar
	CMD cat `+destPath+slash+`bar
	`))
	***REMOVED*** else ***REMOVED***
		testImg = "busybox"
	***REMOVED***

	type testCase struct ***REMOVED***
		spec     mounttypes.Mount
		expected types.MountPoint
	***REMOVED***

	var selinuxSharedLabel string
	if runtime.GOOS == "linux" ***REMOVED***
		selinuxSharedLabel = "z"
	***REMOVED***

	cases := []testCase***REMOVED***
		// use literal strings here for `Type` instead of the defined constants in the volume package to keep this honest
		// Validation of the actual `Mount` struct is done in another test is not needed here
		***REMOVED***
			spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath***REMOVED***,
			expected: types.MountPoint***REMOVED***Driver: volume.DefaultDriverName, Type: "volume", RW: true, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
		***REMOVED***,
		***REMOVED***
			spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath + slash***REMOVED***,
			expected: types.MountPoint***REMOVED***Driver: volume.DefaultDriverName, Type: "volume", RW: true, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
		***REMOVED***,
		***REMOVED***
			spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath, Source: "test1"***REMOVED***,
			expected: types.MountPoint***REMOVED***Type: "volume", Name: "test1", RW: true, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
		***REMOVED***,
		***REMOVED***
			spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath, ReadOnly: true, Source: "test2"***REMOVED***,
			expected: types.MountPoint***REMOVED***Type: "volume", Name: "test2", RW: false, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
		***REMOVED***,
		***REMOVED***
			spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath, Source: "test3", VolumeOptions: &mounttypes.VolumeOptions***REMOVED***DriverConfig: &mounttypes.Driver***REMOVED***Name: volume.DefaultDriverName***REMOVED******REMOVED******REMOVED***,
			expected: types.MountPoint***REMOVED***Driver: volume.DefaultDriverName, Type: "volume", Name: "test3", RW: true, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
		***REMOVED***,
	***REMOVED***

	if SameHostDaemon() ***REMOVED***
		// setup temp dir for testing binds
		tmpDir1, err := ioutil.TempDir("", "test-mounts-api-1")
		c.Assert(err, checker.IsNil)
		defer os.RemoveAll(tmpDir1)
		cases = append(cases, []testCase***REMOVED***
			***REMOVED***
				spec: mounttypes.Mount***REMOVED***
					Type:   "bind",
					Source: tmpDir1,
					Target: destPath,
				***REMOVED***,
				expected: types.MountPoint***REMOVED***
					Type:        "bind",
					RW:          true,
					Destination: destPath,
					Source:      tmpDir1,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				spec:     mounttypes.Mount***REMOVED***Type: "bind", Source: tmpDir1, Target: destPath, ReadOnly: true***REMOVED***,
				expected: types.MountPoint***REMOVED***Type: "bind", RW: false, Destination: destPath, Source: tmpDir1***REMOVED***,
			***REMOVED***,
		***REMOVED***...)

		// for modes only supported on Linux
		if DaemonIsLinux() ***REMOVED***
			tmpDir3, err := ioutils.TempDir("", "test-mounts-api-3")
			c.Assert(err, checker.IsNil)
			defer os.RemoveAll(tmpDir3)

			c.Assert(mount.Mount(tmpDir3, tmpDir3, "none", "bind,rw"), checker.IsNil)
			c.Assert(mount.ForceMount("", tmpDir3, "none", "shared"), checker.IsNil)

			cases = append(cases, []testCase***REMOVED***
				***REMOVED***
					spec:     mounttypes.Mount***REMOVED***Type: "bind", Source: tmpDir3, Target: destPath***REMOVED***,
					expected: types.MountPoint***REMOVED***Type: "bind", RW: true, Destination: destPath, Source: tmpDir3***REMOVED***,
				***REMOVED***,
				***REMOVED***
					spec:     mounttypes.Mount***REMOVED***Type: "bind", Source: tmpDir3, Target: destPath, ReadOnly: true***REMOVED***,
					expected: types.MountPoint***REMOVED***Type: "bind", RW: false, Destination: destPath, Source: tmpDir3***REMOVED***,
				***REMOVED***,
				***REMOVED***
					spec:     mounttypes.Mount***REMOVED***Type: "bind", Source: tmpDir3, Target: destPath, ReadOnly: true, BindOptions: &mounttypes.BindOptions***REMOVED***Propagation: "shared"***REMOVED******REMOVED***,
					expected: types.MountPoint***REMOVED***Type: "bind", RW: false, Destination: destPath, Source: tmpDir3, Propagation: "shared"***REMOVED***,
				***REMOVED***,
			***REMOVED***...)
		***REMOVED***
	***REMOVED***

	if testEnv.OSType != "windows" ***REMOVED*** // Windows does not support volume populate
		cases = append(cases, []testCase***REMOVED***
			***REMOVED***
				spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath, VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED******REMOVED***,
				expected: types.MountPoint***REMOVED***Driver: volume.DefaultDriverName, Type: "volume", RW: true, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
			***REMOVED***,
			***REMOVED***
				spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath + slash, VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED******REMOVED***,
				expected: types.MountPoint***REMOVED***Driver: volume.DefaultDriverName, Type: "volume", RW: true, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
			***REMOVED***,
			***REMOVED***
				spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath, Source: "test4", VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED******REMOVED***,
				expected: types.MountPoint***REMOVED***Type: "volume", Name: "test4", RW: true, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
			***REMOVED***,
			***REMOVED***
				spec:     mounttypes.Mount***REMOVED***Type: "volume", Target: destPath, Source: "test5", ReadOnly: true, VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED******REMOVED***,
				expected: types.MountPoint***REMOVED***Type: "volume", Name: "test5", RW: false, Destination: destPath, Mode: selinuxSharedLabel***REMOVED***,
			***REMOVED***,
		***REMOVED***...)
	***REMOVED***

	type wrapper struct ***REMOVED***
		containertypes.Config
		HostConfig containertypes.HostConfig
	***REMOVED***
	type createResp struct ***REMOVED***
		ID string `json:"Id"`
	***REMOVED***

	ctx := context.Background()
	apiclient := testEnv.APIClient()
	for i, x := range cases ***REMOVED***
		c.Logf("case %d - config: %v", i, x.spec)
		container, err := apiclient.ContainerCreate(
			ctx,
			&containertypes.Config***REMOVED***Image: testImg***REMOVED***,
			&containertypes.HostConfig***REMOVED***Mounts: []mounttypes.Mount***REMOVED***x.spec***REMOVED******REMOVED***,
			&networktypes.NetworkingConfig***REMOVED******REMOVED***,
			"")
		require.NoError(c, err)

		containerInspect, err := apiclient.ContainerInspect(ctx, container.ID)
		require.NoError(c, err)
		mps := containerInspect.Mounts
		require.Len(c, mps, 1)
		mountPoint := mps[0]

		if x.expected.Source != "" ***REMOVED***
			assert.Equal(c, x.expected.Source, mountPoint.Source)
		***REMOVED***
		if x.expected.Name != "" ***REMOVED***
			assert.Equal(c, x.expected.Name, mountPoint.Name)
		***REMOVED***
		if x.expected.Driver != "" ***REMOVED***
			assert.Equal(c, x.expected.Driver, mountPoint.Driver)
		***REMOVED***
		if x.expected.Propagation != "" ***REMOVED***
			assert.Equal(c, x.expected.Propagation, mountPoint.Propagation)
		***REMOVED***
		assert.Equal(c, x.expected.RW, mountPoint.RW)
		assert.Equal(c, x.expected.Type, mountPoint.Type)
		assert.Equal(c, x.expected.Mode, mountPoint.Mode)
		assert.Equal(c, x.expected.Destination, mountPoint.Destination)

		err = apiclient.ContainerStart(ctx, container.ID, types.ContainerStartOptions***REMOVED******REMOVED***)
		require.NoError(c, err)
		poll.WaitOn(c, containerExit(apiclient, container.ID), poll.WithDelay(time.Second))

		err = apiclient.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions***REMOVED***
			RemoveVolumes: true,
			Force:         true,
		***REMOVED***)
		require.NoError(c, err)

		switch ***REMOVED***

		// Named volumes still exist after the container is removed
		case x.spec.Type == "volume" && len(x.spec.Source) > 0:
			_, err := apiclient.VolumeInspect(ctx, mountPoint.Name)
			require.NoError(c, err)

		// Bind mounts are never removed with the container
		case x.spec.Type == "bind":

		// anonymous volumes are removed
		default:
			_, err := apiclient.VolumeInspect(ctx, mountPoint.Name)
			assert.True(c, client.IsErrNotFound(err))
		***REMOVED***
	***REMOVED***
***REMOVED***

func containerExit(apiclient client.APIClient, name string) func(poll.LogT) poll.Result ***REMOVED***
	return func(logT poll.LogT) poll.Result ***REMOVED***
		container, err := apiclient.ContainerInspect(context.Background(), name)
		if err != nil ***REMOVED***
			return poll.Error(err)
		***REMOVED***
		switch container.State.Status ***REMOVED***
		case "created", "running":
			return poll.Continue("container %s is %s, waiting for exit", name, container.State.Status)
		***REMOVED***
		return poll.Success()
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestContainersAPICreateMountsTmpfs(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	type testCase struct ***REMOVED***
		cfg             mounttypes.Mount
		expectedOptions []string
	***REMOVED***
	target := "/foo"
	cases := []testCase***REMOVED***
		***REMOVED***
			cfg: mounttypes.Mount***REMOVED***
				Type:   "tmpfs",
				Target: target***REMOVED***,
			expectedOptions: []string***REMOVED***"rw", "nosuid", "nodev", "noexec", "relatime"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			cfg: mounttypes.Mount***REMOVED***
				Type:   "tmpfs",
				Target: target,
				TmpfsOptions: &mounttypes.TmpfsOptions***REMOVED***
					SizeBytes: 4096 * 1024, Mode: 0700***REMOVED******REMOVED***,
			expectedOptions: []string***REMOVED***"rw", "nosuid", "nodev", "noexec", "relatime", "size=4096k", "mode=700"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	config := containertypes.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"/bin/sh", "-c", fmt.Sprintf("mount | grep 'tmpfs on %s'", target)***REMOVED***,
	***REMOVED***
	for i, x := range cases ***REMOVED***
		cName := fmt.Sprintf("test-tmpfs-%d", i)
		hostConfig := containertypes.HostConfig***REMOVED***
			Mounts: []mounttypes.Mount***REMOVED***x.cfg***REMOVED***,
		***REMOVED***

		_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networktypes.NetworkingConfig***REMOVED******REMOVED***, cName)
		c.Assert(err, checker.IsNil)
		out, _ := dockerCmd(c, "start", "-a", cName)
		for _, option := range x.expectedOptions ***REMOVED***
			c.Assert(out, checker.Contains, option)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Regression test for #33334
// Makes sure that when a container which has a custom stop signal + restart=always
// gets killed (with SIGKILL) by the kill API, that the restart policy is cancelled.
func (s *DockerSuite) TestContainerKillCustomStopSignal(c *check.C) ***REMOVED***
	id := strings.TrimSpace(runSleepingContainer(c, "--stop-signal=SIGTERM", "--restart=always"))
	res, _, err := request.Post("/containers/" + id + "/kill")
	c.Assert(err, checker.IsNil)
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	c.Assert(res.StatusCode, checker.Equals, http.StatusNoContent, check.Commentf(string(b)))
	err = waitInspect(id, "***REMOVED******REMOVED***.State.Running***REMOVED******REMOVED*** ***REMOVED******REMOVED***.State.Restarting***REMOVED******REMOVED***", "false false", 30*time.Second)
	c.Assert(err, checker.IsNil)
***REMOVED***
