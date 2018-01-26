package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/integration-cli/registry"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	"golang.org/x/net/context"
)

// Deprecated
func daemonHost() string ***REMOVED***
	return request.DaemonHost()
***REMOVED***

func deleteImages(images ...string) error ***REMOVED***
	args := []string***REMOVED***dockerBinary, "rmi", "-f"***REMOVED***
	return icmd.RunCmd(icmd.Cmd***REMOVED***Command: append(args, images...)***REMOVED***).Error
***REMOVED***

// Deprecated: use cli.Docker or cli.DockerCmd
func dockerCmdWithError(args ...string) (string, int, error) ***REMOVED***
	result := cli.Docker(cli.Args(args...))
	if result.Error != nil ***REMOVED***
		return result.Combined(), result.ExitCode, result.Compare(icmd.Success)
	***REMOVED***
	return result.Combined(), result.ExitCode, result.Error
***REMOVED***

// Deprecated: use cli.Docker or cli.DockerCmd
func dockerCmd(c *check.C, args ...string) (string, int) ***REMOVED***
	result := cli.DockerCmd(c, args...)
	return result.Combined(), result.ExitCode
***REMOVED***

// Deprecated: use cli.Docker or cli.DockerCmd
func dockerCmdWithResult(args ...string) *icmd.Result ***REMOVED***
	return cli.Docker(cli.Args(args...))
***REMOVED***

func findContainerIP(c *check.C, id string, network string) string ***REMOVED***
	out, _ := dockerCmd(c, "inspect", fmt.Sprintf("--format='***REMOVED******REMOVED*** .NetworkSettings.Networks.%s.IPAddress ***REMOVED******REMOVED***'", network), id)
	return strings.Trim(out, " \r\n'")
***REMOVED***

func getContainerCount(c *check.C) int ***REMOVED***
	const containers = "Containers:"

	result := icmd.RunCommand(dockerBinary, "info")
	result.Assert(c, icmd.Success)

	lines := strings.Split(result.Combined(), "\n")
	for _, line := range lines ***REMOVED***
		if strings.Contains(line, containers) ***REMOVED***
			output := strings.TrimSpace(line)
			output = strings.TrimLeft(output, containers)
			output = strings.Trim(output, " ")
			containerCount, err := strconv.Atoi(output)
			c.Assert(err, checker.IsNil)
			return containerCount
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func inspectFieldAndUnmarshall(c *check.C, name, field string, output interface***REMOVED******REMOVED***) ***REMOVED***
	str := inspectFieldJSON(c, name, field)
	err := json.Unmarshal([]byte(str), output)
	if c != nil ***REMOVED***
		c.Assert(err, check.IsNil, check.Commentf("failed to unmarshal: %v", err))
	***REMOVED***
***REMOVED***

// Deprecated: use cli.Inspect
func inspectFilter(name, filter string) (string, error) ***REMOVED***
	format := fmt.Sprintf("***REMOVED******REMOVED***%s***REMOVED******REMOVED***", filter)
	result := icmd.RunCommand(dockerBinary, "inspect", "-f", format, name)
	if result.Error != nil || result.ExitCode != 0 ***REMOVED***
		return "", fmt.Errorf("failed to inspect %s: %s", name, result.Combined())
	***REMOVED***
	return strings.TrimSpace(result.Combined()), nil
***REMOVED***

// Deprecated: use cli.Inspect
func inspectFieldWithError(name, field string) (string, error) ***REMOVED***
	return inspectFilter(name, fmt.Sprintf(".%s", field))
***REMOVED***

// Deprecated: use cli.Inspect
func inspectField(c *check.C, name, field string) string ***REMOVED***
	out, err := inspectFilter(name, fmt.Sprintf(".%s", field))
	if c != nil ***REMOVED***
		c.Assert(err, check.IsNil)
	***REMOVED***
	return out
***REMOVED***

// Deprecated: use cli.Inspect
func inspectFieldJSON(c *check.C, name, field string) string ***REMOVED***
	out, err := inspectFilter(name, fmt.Sprintf("json .%s", field))
	if c != nil ***REMOVED***
		c.Assert(err, check.IsNil)
	***REMOVED***
	return out
***REMOVED***

// Deprecated: use cli.Inspect
func inspectFieldMap(c *check.C, name, path, field string) string ***REMOVED***
	out, err := inspectFilter(name, fmt.Sprintf("index .%s %q", path, field))
	if c != nil ***REMOVED***
		c.Assert(err, check.IsNil)
	***REMOVED***
	return out
***REMOVED***

// Deprecated: use cli.Inspect
func inspectMountSourceField(name, destination string) (string, error) ***REMOVED***
	m, err := inspectMountPoint(name, destination)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return m.Source, nil
***REMOVED***

// Deprecated: use cli.Inspect
func inspectMountPoint(name, destination string) (types.MountPoint, error) ***REMOVED***
	out, err := inspectFilter(name, "json .Mounts")
	if err != nil ***REMOVED***
		return types.MountPoint***REMOVED******REMOVED***, err
	***REMOVED***

	return inspectMountPointJSON(out, destination)
***REMOVED***

var errMountNotFound = errors.New("mount point not found")

// Deprecated: use cli.Inspect
func inspectMountPointJSON(j, destination string) (types.MountPoint, error) ***REMOVED***
	var mp []types.MountPoint
	if err := json.Unmarshal([]byte(j), &mp); err != nil ***REMOVED***
		return types.MountPoint***REMOVED******REMOVED***, err
	***REMOVED***

	var m *types.MountPoint
	for _, c := range mp ***REMOVED***
		if c.Destination == destination ***REMOVED***
			m = &c
			break
		***REMOVED***
	***REMOVED***

	if m == nil ***REMOVED***
		return types.MountPoint***REMOVED******REMOVED***, errMountNotFound
	***REMOVED***

	return *m, nil
***REMOVED***

// Deprecated: use cli.Inspect
func inspectImage(c *check.C, name, filter string) string ***REMOVED***
	args := []string***REMOVED***"inspect", "--type", "image"***REMOVED***
	if filter != "" ***REMOVED***
		format := fmt.Sprintf("***REMOVED******REMOVED***%s***REMOVED******REMOVED***", filter)
		args = append(args, "-f", format)
	***REMOVED***
	args = append(args, name)
	result := icmd.RunCommand(dockerBinary, args...)
	result.Assert(c, icmd.Success)
	return strings.TrimSpace(result.Combined())
***REMOVED***

func getIDByName(c *check.C, name string) string ***REMOVED***
	id, err := inspectFieldWithError(name, "Id")
	c.Assert(err, checker.IsNil)
	return id
***REMOVED***

// Deprecated: use cli.Build
func buildImageSuccessfully(c *check.C, name string, cmdOperators ...cli.CmdOperator) ***REMOVED***
	buildImage(name, cmdOperators...).Assert(c, icmd.Success)
***REMOVED***

// Deprecated: use cli.Build
func buildImage(name string, cmdOperators ...cli.CmdOperator) *icmd.Result ***REMOVED***
	return cli.Docker(cli.Build(name), cmdOperators...)
***REMOVED***

// Deprecated: use trustedcmd
func trustedBuild(cmd *icmd.Cmd) func() ***REMOVED***
	trustedCmd(cmd)
	return nil
***REMOVED***

// Write `content` to the file at path `dst`, creating it if necessary,
// as well as any missing directories.
// The file is truncated if it already exists.
// Fail the test when error occurs.
func writeFile(dst, content string, c *check.C) ***REMOVED***
	// Create subdirectories if necessary
	c.Assert(os.MkdirAll(path.Dir(dst), 0700), check.IsNil)
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0700)
	c.Assert(err, check.IsNil)
	defer f.Close()
	// Write content (truncate if it exists)
	_, err = io.Copy(f, strings.NewReader(content))
	c.Assert(err, check.IsNil)
***REMOVED***

// Return the contents of file at path `src`.
// Fail the test when error occurs.
func readFile(src string, c *check.C) (content string) ***REMOVED***
	data, err := ioutil.ReadFile(src)
	c.Assert(err, check.IsNil)

	return string(data)
***REMOVED***

func containerStorageFile(containerID, basename string) string ***REMOVED***
	return filepath.Join(testEnv.PlatformDefaults.ContainerStoragePath, containerID, basename)
***REMOVED***

// docker commands that use this function must be run with the '-d' switch.
func runCommandAndReadContainerFile(c *check.C, filename string, command string, args ...string) []byte ***REMOVED***
	result := icmd.RunCommand(command, args...)
	result.Assert(c, icmd.Success)
	contID := strings.TrimSpace(result.Combined())
	if err := waitRun(contID); err != nil ***REMOVED***
		c.Fatalf("%v: %q", contID, err)
	***REMOVED***
	return readContainerFile(c, contID, filename)
***REMOVED***

func readContainerFile(c *check.C, containerID, filename string) []byte ***REMOVED***
	f, err := os.Open(containerStorageFile(containerID, filename))
	c.Assert(err, checker.IsNil)
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	c.Assert(err, checker.IsNil)
	return content
***REMOVED***

func readContainerFileWithExec(c *check.C, containerID, filename string) []byte ***REMOVED***
	result := icmd.RunCommand(dockerBinary, "exec", containerID, "cat", filename)
	result.Assert(c, icmd.Success)
	return []byte(result.Combined())
***REMOVED***

// daemonTime provides the current time on the daemon host
func daemonTime(c *check.C) time.Time ***REMOVED***
	if testEnv.IsLocalDaemon() ***REMOVED***
		return time.Now()
	***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, check.IsNil)
	defer cli.Close()

	info, err := cli.Info(context.Background())
	c.Assert(err, check.IsNil)

	dt, err := time.Parse(time.RFC3339Nano, info.SystemTime)
	c.Assert(err, check.IsNil, check.Commentf("invalid time format in GET /info response"))
	return dt
***REMOVED***

// daemonUnixTime returns the current time on the daemon host with nanoseconds precision.
// It return the time formatted how the client sends timestamps to the server.
func daemonUnixTime(c *check.C) string ***REMOVED***
	return parseEventTime(daemonTime(c))
***REMOVED***

func parseEventTime(t time.Time) string ***REMOVED***
	return fmt.Sprintf("%d.%09d", t.Unix(), int64(t.Nanosecond()))
***REMOVED***

func setupRegistry(c *check.C, schema1 bool, auth, tokenURL string) *registry.V2 ***REMOVED***
	reg, err := registry.NewV2(schema1, auth, tokenURL, privateRegistryURL)
	c.Assert(err, check.IsNil)

	// Wait for registry to be ready to serve requests.
	for i := 0; i != 50; i++ ***REMOVED***
		if err = reg.Ping(); err == nil ***REMOVED***
			break
		***REMOVED***
		time.Sleep(100 * time.Millisecond)
	***REMOVED***

	c.Assert(err, check.IsNil, check.Commentf("Timeout waiting for test registry to become available: %v", err))
	return reg
***REMOVED***

func setupNotary(c *check.C) *testNotary ***REMOVED***
	ts, err := newTestNotary(c)
	c.Assert(err, check.IsNil)

	return ts
***REMOVED***

// appendBaseEnv appends the minimum set of environment variables to exec the
// docker cli binary for testing with correct configuration to the given env
// list.
func appendBaseEnv(isTLS bool, env ...string) []string ***REMOVED***
	preserveList := []string***REMOVED***
		// preserve remote test host
		"DOCKER_HOST",

		// windows: requires preserving SystemRoot, otherwise dial tcp fails
		// with "GetAddrInfoW: A non-recoverable error occurred during a database lookup."
		"SystemRoot",

		// testing help text requires the $PATH to dockerd is set
		"PATH",
	***REMOVED***
	if isTLS ***REMOVED***
		preserveList = append(preserveList, "DOCKER_TLS_VERIFY", "DOCKER_CERT_PATH")
	***REMOVED***

	for _, key := range preserveList ***REMOVED***
		if val := os.Getenv(key); val != "" ***REMOVED***
			env = append(env, fmt.Sprintf("%s=%s", key, val))
		***REMOVED***
	***REMOVED***
	return env
***REMOVED***

func createTmpFile(c *check.C, content string) string ***REMOVED***
	f, err := ioutil.TempFile("", "testfile")
	c.Assert(err, check.IsNil)

	filename := f.Name()

	err = ioutil.WriteFile(filename, []byte(content), 0644)
	c.Assert(err, check.IsNil)

	return filename
***REMOVED***

// waitRun will wait for the specified container to be running, maximum 5 seconds.
// Deprecated: use cli.WaitFor
func waitRun(contID string) error ***REMOVED***
	return waitInspect(contID, "***REMOVED******REMOVED***.State.Running***REMOVED******REMOVED***", "true", 5*time.Second)
***REMOVED***

// waitInspect will wait for the specified container to have the specified string
// in the inspect output. It will wait until the specified timeout (in seconds)
// is reached.
// Deprecated: use cli.WaitFor
func waitInspect(name, expr, expected string, timeout time.Duration) error ***REMOVED***
	return waitInspectWithArgs(name, expr, expected, timeout)
***REMOVED***

// Deprecated: use cli.WaitFor
func waitInspectWithArgs(name, expr, expected string, timeout time.Duration, arg ...string) error ***REMOVED***
	return daemon.WaitInspectWithArgs(dockerBinary, name, expr, expected, timeout, arg...)
***REMOVED***

func getInspectBody(c *check.C, version, id string) []byte ***REMOVED***
	cli, err := request.NewEnvClientWithVersion(version)
	c.Assert(err, check.IsNil)
	defer cli.Close()
	_, body, err := cli.ContainerInspectWithRaw(context.Background(), id, false)
	c.Assert(err, check.IsNil)
	return body
***REMOVED***

// Run a long running idle task in a background container using the
// system-specific default image and command.
func runSleepingContainer(c *check.C, extraArgs ...string) string ***REMOVED***
	return runSleepingContainerInImage(c, defaultSleepImage, extraArgs...)
***REMOVED***

// Run a long running idle task in a background container using the specified
// image and the system-specific command.
func runSleepingContainerInImage(c *check.C, image string, extraArgs ...string) string ***REMOVED***
	args := []string***REMOVED***"run", "-d"***REMOVED***
	args = append(args, extraArgs...)
	args = append(args, image)
	args = append(args, sleepCommandForDaemonPlatform()...)
	return strings.TrimSpace(cli.DockerCmd(c, args...).Combined())
***REMOVED***

// minimalBaseImage returns the name of the minimal base image for the current
// daemon platform.
func minimalBaseImage() string ***REMOVED***
	return testEnv.PlatformDefaults.BaseImage
***REMOVED***

func getGoroutineNumber() (int, error) ***REMOVED***
	cli, err := client.NewEnvClient()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer cli.Close()

	info, err := cli.Info(context.Background())
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return info.NGoroutines, nil
***REMOVED***

func waitForGoroutines(expected int) error ***REMOVED***
	t := time.After(30 * time.Second)
	for ***REMOVED***
		select ***REMOVED***
		case <-t:
			n, err := getGoroutineNumber()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if n > expected ***REMOVED***
				return fmt.Errorf("leaked goroutines: expected less than or equal to %d, got: %d", expected, n)
			***REMOVED***
		default:
			n, err := getGoroutineNumber()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if n <= expected ***REMOVED***
				return nil
			***REMOVED***
			time.Sleep(200 * time.Millisecond)
		***REMOVED***
	***REMOVED***
***REMOVED***

// getErrorMessage returns the error message from an error API response
func getErrorMessage(c *check.C, body []byte) string ***REMOVED***
	var resp types.ErrorResponse
	c.Assert(json.Unmarshal(body, &resp), check.IsNil)
	return strings.TrimSpace(resp.Message)
***REMOVED***

func waitAndAssert(c *check.C, timeout time.Duration, f checkF, checker check.Checker, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	after := time.After(timeout)
	for ***REMOVED***
		v, comment := f(c)
		assert, _ := checker.Check(append([]interface***REMOVED******REMOVED******REMOVED***v***REMOVED***, args...), checker.Info().Params)
		select ***REMOVED***
		case <-after:
			assert = true
		default:
		***REMOVED***
		if assert ***REMOVED***
			if comment != nil ***REMOVED***
				args = append(args, comment)
			***REMOVED***
			c.Assert(v, checker, args...)
			return
		***REMOVED***
		time.Sleep(100 * time.Millisecond)
	***REMOVED***
***REMOVED***

type checkF func(*check.C) (interface***REMOVED******REMOVED***, check.CommentInterface)
type reducer func(...interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED***

func reducedCheck(r reducer, funcs ...checkF) checkF ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		var values []interface***REMOVED******REMOVED***
		var comments []string
		for _, f := range funcs ***REMOVED***
			v, comment := f(c)
			values = append(values, v)
			if comment != nil ***REMOVED***
				comments = append(comments, comment.CheckCommentString())
			***REMOVED***
		***REMOVED***
		return r(values...), check.Commentf("%v", strings.Join(comments, ", "))
	***REMOVED***
***REMOVED***

func sumAsIntegers(vals ...interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	var s int
	for _, v := range vals ***REMOVED***
		s += v.(int)
	***REMOVED***
	return s
***REMOVED***
