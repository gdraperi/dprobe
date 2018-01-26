package main

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/go-check/check"
)

func init() ***REMOVED***
	// FIXME. Temporarily turning this off for Windows as GH16039 was breaking
	// Windows to Linux CI @icecrime
	if runtime.GOOS != "windows" ***REMOVED***
		check.Suite(newDockerHubPullSuite())
	***REMOVED***
***REMOVED***

// DockerHubPullSuite provides an isolated daemon that doesn't have all the
// images that are baked into our 'global' test environment daemon (e.g.,
// busybox, httpserver, ...).
//
// We use it for push/pull tests where we want to start fresh, and measure the
// relative impact of each individual operation. As part of this suite, all
// images are removed after each test.
type DockerHubPullSuite struct ***REMOVED***
	d  *daemon.Daemon
	ds *DockerSuite
***REMOVED***

// newDockerHubPullSuite returns a new instance of a DockerHubPullSuite.
func newDockerHubPullSuite() *DockerHubPullSuite ***REMOVED***
	return &DockerHubPullSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// SetUpSuite starts the suite daemon.
func (s *DockerHubPullSuite) SetUpSuite(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon)
	s.d = daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)
	s.d.Start(c)
***REMOVED***

// TearDownSuite stops the suite daemon.
func (s *DockerHubPullSuite) TearDownSuite(c *check.C) ***REMOVED***
	if s.d != nil ***REMOVED***
		s.d.Stop(c)
	***REMOVED***
***REMOVED***

// SetUpTest declares that all tests of this suite require network.
func (s *DockerHubPullSuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, Network)
***REMOVED***

// TearDownTest removes all images from the suite daemon.
func (s *DockerHubPullSuite) TearDownTest(c *check.C) ***REMOVED***
	out := s.Cmd(c, "images", "-aq")
	images := strings.Split(out, "\n")
	images = append([]string***REMOVED***"rmi", "-f"***REMOVED***, images...)
	s.d.Cmd(images...)
	s.ds.TearDownTest(c)
***REMOVED***

// Cmd executes a command against the suite daemon and returns the combined
// output. The function fails the test when the command returns an error.
func (s *DockerHubPullSuite) Cmd(c *check.C, name string, arg ...string) string ***REMOVED***
	out, err := s.CmdWithError(name, arg...)
	c.Assert(err, checker.IsNil, check.Commentf("%q failed with errors: %s, %v", strings.Join(arg, " "), out, err))
	return out
***REMOVED***

// CmdWithError executes a command against the suite daemon and returns the
// combined output as well as any error.
func (s *DockerHubPullSuite) CmdWithError(name string, arg ...string) (string, error) ***REMOVED***
	c := s.MakeCmd(name, arg...)
	b, err := c.CombinedOutput()
	return string(b), err
***REMOVED***

// MakeCmd returns an exec.Cmd command to run against the suite daemon.
func (s *DockerHubPullSuite) MakeCmd(name string, arg ...string) *exec.Cmd ***REMOVED***
	args := []string***REMOVED***"--host", s.d.Sock(), name***REMOVED***
	args = append(args, arg...)
	return exec.Command(dockerBinary, args...)
***REMOVED***
