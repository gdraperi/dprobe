package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

func (s *DockerSuite) TestLogsAPIWithStdout(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "-t", "busybox", "/bin/sh", "-c", "while true; do echo hello; sleep 1; done")
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	type logOut struct ***REMOVED***
		out string
		err error
	***REMOVED***

	chLog := make(chan logOut)
	res, body, err := request.Get(fmt.Sprintf("/containers/%s/logs?follow=1&stdout=1&timestamps=1", id))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	go func() ***REMOVED***
		defer body.Close()
		out, err := bufio.NewReader(body).ReadString('\n')
		if err != nil ***REMOVED***
			chLog <- logOut***REMOVED***"", err***REMOVED***
			return
		***REMOVED***
		chLog <- logOut***REMOVED***strings.TrimSpace(out), err***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case l := <-chLog:
		c.Assert(l.err, checker.IsNil)
		if !strings.HasSuffix(l.out, "hello") ***REMOVED***
			c.Fatalf("expected log output to container 'hello', but it does not")
		***REMOVED***
	case <-time.After(30 * time.Second):
		c.Fatal("timeout waiting for logs to exit")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestLogsAPINoStdoutNorStderr(c *check.C) ***REMOVED***
	name := "logs_test"
	dockerCmd(c, "run", "-d", "-t", "--name", name, "busybox", "/bin/sh")
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	_, err = cli.ContainerLogs(context.Background(), name, types.ContainerLogsOptions***REMOVED******REMOVED***)
	expected := "Bad parameters: you must choose at least one stream"
	c.Assert(err.Error(), checker.Contains, expected)
***REMOVED***

// Regression test for #12704
func (s *DockerSuite) TestLogsAPIFollowEmptyOutput(c *check.C) ***REMOVED***
	name := "logs_test"
	t0 := time.Now()
	dockerCmd(c, "run", "-d", "-t", "--name", name, "busybox", "sleep", "10")

	_, body, err := request.Get(fmt.Sprintf("/containers/%s/logs?follow=1&stdout=1&stderr=1&tail=all", name))
	t1 := time.Now()
	c.Assert(err, checker.IsNil)
	body.Close()
	elapsed := t1.Sub(t0).Seconds()
	if elapsed > 20.0 ***REMOVED***
		c.Fatalf("HTTP response was not immediate (elapsed %.1fs)", elapsed)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestLogsAPIContainerNotFound(c *check.C) ***REMOVED***
	name := "nonExistentContainer"
	resp, _, err := request.Get(fmt.Sprintf("/containers/%s/logs?follow=1&stdout=1&stderr=1&tail=all", name))
	c.Assert(err, checker.IsNil)
	c.Assert(resp.StatusCode, checker.Equals, http.StatusNotFound)
***REMOVED***

func (s *DockerSuite) TestLogsAPIUntilFutureFollow(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	name := "logsuntilfuturefollow"
	dockerCmd(c, "run", "-d", "--name", name, "busybox", "/bin/sh", "-c", "while true; do date +%s; sleep 1; done")
	c.Assert(waitRun(name), checker.IsNil)

	untilSecs := 5
	untilDur, err := time.ParseDuration(fmt.Sprintf("%ds", untilSecs))
	c.Assert(err, checker.IsNil)
	until := daemonTime(c).Add(untilDur)

	client, err := request.NewClient()
	if err != nil ***REMOVED***
		c.Fatal(err)
	***REMOVED***

	cfg := types.ContainerLogsOptions***REMOVED***Until: until.Format(time.RFC3339Nano), Follow: true, ShowStdout: true, Timestamps: true***REMOVED***
	reader, err := client.ContainerLogs(context.Background(), name, cfg)
	c.Assert(err, checker.IsNil)

	type logOut struct ***REMOVED***
		out string
		err error
	***REMOVED***

	chLog := make(chan logOut)

	go func() ***REMOVED***
		bufReader := bufio.NewReader(reader)
		defer reader.Close()
		for i := 0; i < untilSecs; i++ ***REMOVED***
			out, _, err := bufReader.ReadLine()
			if err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					return
				***REMOVED***
				chLog <- logOut***REMOVED***"", err***REMOVED***
				return
			***REMOVED***

			chLog <- logOut***REMOVED***strings.TrimSpace(string(out)), err***REMOVED***
		***REMOVED***
	***REMOVED***()

	for i := 0; i < untilSecs; i++ ***REMOVED***
		select ***REMOVED***
		case l := <-chLog:
			c.Assert(l.err, checker.IsNil)
			i, err := strconv.ParseInt(strings.Split(l.out, " ")[1], 10, 64)
			c.Assert(err, checker.IsNil)
			c.Assert(time.Unix(i, 0).UnixNano(), checker.LessOrEqualThan, until.UnixNano())
		case <-time.After(20 * time.Second):
			c.Fatal("timeout waiting for logs to exit")
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestLogsAPIUntil(c *check.C) ***REMOVED***
	name := "logsuntil"
	dockerCmd(c, "run", "--name", name, "busybox", "/bin/sh", "-c", "for i in $(seq 1 3); do echo log$i; sleep 1; done")

	client, err := request.NewClient()
	if err != nil ***REMOVED***
		c.Fatal(err)
	***REMOVED***

	extractBody := func(c *check.C, cfg types.ContainerLogsOptions) []string ***REMOVED***
		reader, err := client.ContainerLogs(context.Background(), name, cfg)
		c.Assert(err, checker.IsNil)

		actualStdout := new(bytes.Buffer)
		actualStderr := ioutil.Discard
		_, err = stdcopy.StdCopy(actualStdout, actualStderr, reader)
		c.Assert(err, checker.IsNil)

		return strings.Split(actualStdout.String(), "\n")
	***REMOVED***

	// Get timestamp of second log line
	allLogs := extractBody(c, types.ContainerLogsOptions***REMOVED***Timestamps: true, ShowStdout: true***REMOVED***)
	c.Assert(len(allLogs), checker.GreaterOrEqualThan, 3)

	t, err := time.Parse(time.RFC3339Nano, strings.Split(allLogs[1], " ")[0])
	c.Assert(err, checker.IsNil)
	until := t.Format(time.RFC3339Nano)

	// Get logs until the timestamp of second line, i.e. first two lines
	logs := extractBody(c, types.ContainerLogsOptions***REMOVED***Timestamps: true, ShowStdout: true, Until: until***REMOVED***)

	// Ensure log lines after cut-off are excluded
	logsString := strings.Join(logs, "\n")
	c.Assert(logsString, checker.Not(checker.Contains), "log3", check.Commentf("unexpected log message returned, until=%v", until))
***REMOVED***

func (s *DockerSuite) TestLogsAPIUntilDefaultValue(c *check.C) ***REMOVED***
	name := "logsuntildefaultval"
	dockerCmd(c, "run", "--name", name, "busybox", "/bin/sh", "-c", "for i in $(seq 1 3); do echo log$i; done")

	client, err := request.NewClient()
	if err != nil ***REMOVED***
		c.Fatal(err)
	***REMOVED***

	extractBody := func(c *check.C, cfg types.ContainerLogsOptions) []string ***REMOVED***
		reader, err := client.ContainerLogs(context.Background(), name, cfg)
		c.Assert(err, checker.IsNil)

		actualStdout := new(bytes.Buffer)
		actualStderr := ioutil.Discard
		_, err = stdcopy.StdCopy(actualStdout, actualStderr, reader)
		c.Assert(err, checker.IsNil)

		return strings.Split(actualStdout.String(), "\n")
	***REMOVED***

	// Get timestamp of second log line
	allLogs := extractBody(c, types.ContainerLogsOptions***REMOVED***Timestamps: true, ShowStdout: true***REMOVED***)

	// Test with default value specified and parameter omitted
	defaultLogs := extractBody(c, types.ContainerLogsOptions***REMOVED***Timestamps: true, ShowStdout: true, Until: "0"***REMOVED***)
	c.Assert(defaultLogs, checker.DeepEquals, allLogs)
***REMOVED***
