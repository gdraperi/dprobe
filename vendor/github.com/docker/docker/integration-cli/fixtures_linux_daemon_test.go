package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/fixtures/load"
	"github.com/go-check/check"
)

type testingT interface ***REMOVED***
	logT
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type logT interface ***REMOVED***
	Logf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

var ensureSyscallTestOnce sync.Once

func ensureSyscallTest(c *check.C) ***REMOVED***
	var doIt bool
	ensureSyscallTestOnce.Do(func() ***REMOVED***
		doIt = true
	***REMOVED***)
	if !doIt ***REMOVED***
		return
	***REMOVED***
	defer testEnv.ProtectImage(c, "syscall-test:latest")

	// if no match, must build in docker, which is significantly slower
	// (slower mostly because of the vfs graphdriver)
	if testEnv.OSType != runtime.GOOS ***REMOVED***
		ensureSyscallTestBuild(c)
		return
	***REMOVED***

	tmp, err := ioutil.TempDir("", "syscall-test-build")
	c.Assert(err, checker.IsNil, check.Commentf("couldn't create temp dir"))
	defer os.RemoveAll(tmp)

	gcc, err := exec.LookPath("gcc")
	c.Assert(err, checker.IsNil, check.Commentf("could not find gcc"))

	tests := []string***REMOVED***"userns", "ns", "acct", "setuid", "setgid", "socket", "raw"***REMOVED***
	for _, test := range tests ***REMOVED***
		out, err := exec.Command(gcc, "-g", "-Wall", "-static", fmt.Sprintf("../contrib/syscall-test/%s.c", test), "-o", fmt.Sprintf("%s/%s-test", tmp, test)).CombinedOutput()
		c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	***REMOVED***

	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" ***REMOVED***
		out, err := exec.Command(gcc, "-s", "-m32", "-nostdlib", "-static", "../contrib/syscall-test/exit32.s", "-o", tmp+"/"+"exit32-test").CombinedOutput()
		c.Assert(err, checker.IsNil, check.Commentf(string(out)))
	***REMOVED***

	dockerFile := filepath.Join(tmp, "Dockerfile")
	content := []byte(`
	FROM debian:jessie
	COPY . /usr/bin/
	`)
	err = ioutil.WriteFile(dockerFile, content, 600)
	c.Assert(err, checker.IsNil)

	var buildArgs []string
	if arg := os.Getenv("DOCKER_BUILD_ARGS"); strings.TrimSpace(arg) != "" ***REMOVED***
		buildArgs = strings.Split(arg, " ")
	***REMOVED***
	buildArgs = append(buildArgs, []string***REMOVED***"-q", "-t", "syscall-test", tmp***REMOVED***...)
	buildArgs = append([]string***REMOVED***"build"***REMOVED***, buildArgs...)
	dockerCmd(c, buildArgs...)
***REMOVED***

func ensureSyscallTestBuild(c *check.C) ***REMOVED***
	err := load.FrozenImagesLinux(testEnv.APIClient(), "buildpack-deps:jessie")
	c.Assert(err, checker.IsNil)

	var buildArgs []string
	if arg := os.Getenv("DOCKER_BUILD_ARGS"); strings.TrimSpace(arg) != "" ***REMOVED***
		buildArgs = strings.Split(arg, " ")
	***REMOVED***
	buildArgs = append(buildArgs, []string***REMOVED***"-q", "-t", "syscall-test", "../contrib/syscall-test"***REMOVED***...)
	buildArgs = append([]string***REMOVED***"build"***REMOVED***, buildArgs...)
	dockerCmd(c, buildArgs...)
***REMOVED***

func ensureNNPTest(c *check.C) ***REMOVED***
	defer testEnv.ProtectImage(c, "nnp-test:latest")
	if testEnv.OSType != runtime.GOOS ***REMOVED***
		ensureNNPTestBuild(c)
		return
	***REMOVED***

	tmp, err := ioutil.TempDir("", "docker-nnp-test")
	c.Assert(err, checker.IsNil)

	gcc, err := exec.LookPath("gcc")
	c.Assert(err, checker.IsNil, check.Commentf("could not find gcc"))

	out, err := exec.Command(gcc, "-g", "-Wall", "-static", "../contrib/nnp-test/nnp-test.c", "-o", filepath.Join(tmp, "nnp-test")).CombinedOutput()
	c.Assert(err, checker.IsNil, check.Commentf(string(out)))

	dockerfile := filepath.Join(tmp, "Dockerfile")
	content := `
	FROM debian:jessie
	COPY . /usr/bin
	RUN chmod +s /usr/bin/nnp-test
	`
	err = ioutil.WriteFile(dockerfile, []byte(content), 600)
	c.Assert(err, checker.IsNil, check.Commentf("could not write Dockerfile for nnp-test image"))

	var buildArgs []string
	if arg := os.Getenv("DOCKER_BUILD_ARGS"); strings.TrimSpace(arg) != "" ***REMOVED***
		buildArgs = strings.Split(arg, " ")
	***REMOVED***
	buildArgs = append(buildArgs, []string***REMOVED***"-q", "-t", "nnp-test", tmp***REMOVED***...)
	buildArgs = append([]string***REMOVED***"build"***REMOVED***, buildArgs...)
	dockerCmd(c, buildArgs...)
***REMOVED***

func ensureNNPTestBuild(c *check.C) ***REMOVED***
	err := load.FrozenImagesLinux(testEnv.APIClient(), "buildpack-deps:jessie")
	c.Assert(err, checker.IsNil)

	var buildArgs []string
	if arg := os.Getenv("DOCKER_BUILD_ARGS"); strings.TrimSpace(arg) != "" ***REMOVED***
		buildArgs = strings.Split(arg, " ")
	***REMOVED***
	buildArgs = append(buildArgs, []string***REMOVED***"-q", "-t", "npp-test", "../contrib/nnp-test"***REMOVED***...)
	buildArgs = append([]string***REMOVED***"build"***REMOVED***, buildArgs...)
	dockerCmd(c, buildArgs...)
***REMOVED***
