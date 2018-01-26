package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSuite) BenchmarkConcurrentContainerActions(c *check.C) ***REMOVED***
	maxConcurrency := runtime.GOMAXPROCS(0)
	numIterations := c.N
	outerGroup := &sync.WaitGroup***REMOVED******REMOVED***
	outerGroup.Add(maxConcurrency)
	chErr := make(chan error, numIterations*2*maxConcurrency)

	for i := 0; i < maxConcurrency; i++ ***REMOVED***
		go func() ***REMOVED***
			defer outerGroup.Done()
			innerGroup := &sync.WaitGroup***REMOVED******REMOVED***
			innerGroup.Add(2)

			go func() ***REMOVED***
				defer innerGroup.Done()
				for i := 0; i < numIterations; i++ ***REMOVED***
					args := []string***REMOVED***"run", "-d", defaultSleepImage***REMOVED***
					args = append(args, sleepCommandForDaemonPlatform()...)
					out, _, err := dockerCmdWithError(args...)
					if err != nil ***REMOVED***
						chErr <- fmt.Errorf(out)
						return
					***REMOVED***

					id := strings.TrimSpace(out)
					tmpDir, err := ioutil.TempDir("", "docker-concurrent-test-"+id)
					if err != nil ***REMOVED***
						chErr <- err
						return
					***REMOVED***
					defer os.RemoveAll(tmpDir)
					out, _, err = dockerCmdWithError("cp", id+":/tmp", tmpDir)
					if err != nil ***REMOVED***
						chErr <- fmt.Errorf(out)
						return
					***REMOVED***

					out, _, err = dockerCmdWithError("kill", id)
					if err != nil ***REMOVED***
						chErr <- fmt.Errorf(out)
					***REMOVED***

					out, _, err = dockerCmdWithError("start", id)
					if err != nil ***REMOVED***
						chErr <- fmt.Errorf(out)
					***REMOVED***

					out, _, err = dockerCmdWithError("kill", id)
					if err != nil ***REMOVED***
						chErr <- fmt.Errorf(out)
					***REMOVED***

					// don't do an rm -f here since it can potentially ignore errors from the graphdriver
					out, _, err = dockerCmdWithError("rm", id)
					if err != nil ***REMOVED***
						chErr <- fmt.Errorf(out)
					***REMOVED***
				***REMOVED***
			***REMOVED***()

			go func() ***REMOVED***
				defer innerGroup.Done()
				for i := 0; i < numIterations; i++ ***REMOVED***
					out, _, err := dockerCmdWithError("ps")
					if err != nil ***REMOVED***
						chErr <- fmt.Errorf(out)
					***REMOVED***
				***REMOVED***
			***REMOVED***()

			innerGroup.Wait()
		***REMOVED***()
	***REMOVED***

	outerGroup.Wait()
	close(chErr)

	for err := range chErr ***REMOVED***
		c.Assert(err, checker.IsNil)
	***REMOVED***
***REMOVED***
