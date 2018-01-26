// +build !windows

package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

type logMessage struct ***REMOVED***
	err  error
	data []byte
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogs(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	// we have multiple services here for detecting the goroutine issue #28915
	services := map[string]string***REMOVED***
		"TestServiceLogs1": "hello1",
		"TestServiceLogs2": "hello2",
	***REMOVED***

	for name, message := range services ***REMOVED***
		out, err := d.Cmd("service", "create", "--detach", "--no-resolve-image", "--name", name, "busybox",
			"sh", "-c", fmt.Sprintf("echo %s; tail -f /dev/null", message))
		c.Assert(err, checker.IsNil)
		c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")
	***REMOVED***

	// make sure task has been deployed.
	waitAndAssert(c, defaultReconciliationTimeout,
		d.CheckRunningTaskImages, checker.DeepEquals,
		map[string]int***REMOVED***"busybox:latest": len(services)***REMOVED***)

	for name, message := range services ***REMOVED***
		out, err := d.Cmd("service", "logs", name)
		c.Assert(err, checker.IsNil)
		c.Logf("log for %q: %q", name, out)
		c.Assert(out, checker.Contains, message)
	***REMOVED***
***REMOVED***

// countLogLines returns a closure that can be used with waitAndAssert to
// verify that a minimum number of expected container log messages have been
// output.
func countLogLines(d *daemon.Swarm, name string) func(*check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	return func(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
		result := icmd.RunCmd(d.Command("service", "logs", "-t", "--raw", name))
		result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
		// if this returns an emptystring, trying to split it later will return
		// an array containing emptystring. a valid log line will NEVER be
		// emptystring because we ask for the timestamp.
		if result.Stdout() == "" ***REMOVED***
			return 0, check.Commentf("Empty stdout")
		***REMOVED***
		lines := strings.Split(strings.TrimSpace(result.Stdout()), "\n")
		return len(lines), check.Commentf("output, %q", string(result.Stdout()))
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsCompleteness(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "TestServiceLogsCompleteness"

	// make a service that prints 6 lines
	out, err := d.Cmd("service", "create", "--detach", "--no-resolve-image", "--name", name, "busybox", "sh", "-c", "for line in $(seq 0 5); do echo log test $line; done; sleep 100000")
	c.Assert(err, checker.IsNil)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")

	// make sure task has been deployed.
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)
	// and make sure we have all the log lines
	waitAndAssert(c, defaultReconciliationTimeout, countLogLines(d, name), checker.Equals, 6)

	out, err = d.Cmd("service", "logs", name)
	c.Assert(err, checker.IsNil)
	lines := strings.Split(strings.TrimSpace(out), "\n")

	// i have heard anecdotal reports that logs may come back from the engine
	// mis-ordered. if this tests fails, consider the possibility that that
	// might be occurring
	for i, line := range lines ***REMOVED***
		c.Assert(line, checker.Contains, fmt.Sprintf("log test %v", i))
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsTail(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "TestServiceLogsTail"

	// make a service that prints 6 lines
	out, err := d.Cmd("service", "create", "--detach", "--no-resolve-image", "--name", name, "busybox", "sh", "-c", "for line in $(seq 1 6); do echo log test $line; done; sleep 100000")
	c.Assert(err, checker.IsNil)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")

	// make sure task has been deployed.
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)
	waitAndAssert(c, defaultReconciliationTimeout, countLogLines(d, name), checker.Equals, 6)

	out, err = d.Cmd("service", "logs", "--tail=2", name)
	c.Assert(err, checker.IsNil)
	lines := strings.Split(strings.TrimSpace(out), "\n")

	for i, line := range lines ***REMOVED***
		// doing i+5 is hacky but not too fragile, it's good enough. if it flakes something else is wrong
		c.Assert(line, checker.Contains, fmt.Sprintf("log test %v", i+5))
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsSince(c *check.C) ***REMOVED***
	// See DockerSuite.TestLogsSince, which is where this comes from
	d := s.AddDaemon(c, true, true)

	name := "TestServiceLogsSince"

	out, err := d.Cmd("service", "create", "--detach", "--no-resolve-image", "--name", name, "busybox", "sh", "-c", "for i in $(seq 1 3); do sleep .1; echo log$i; done; sleep 10000000")
	c.Assert(err, checker.IsNil)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)
	// wait a sec for the logs to come in
	waitAndAssert(c, defaultReconciliationTimeout, countLogLines(d, name), checker.Equals, 3)

	out, err = d.Cmd("service", "logs", "-t", name)
	c.Assert(err, checker.IsNil)

	log2Line := strings.Split(strings.Split(out, "\n")[1], " ")
	t, err := time.Parse(time.RFC3339Nano, log2Line[0]) // timestamp log2 is written
	c.Assert(err, checker.IsNil)
	u := t.Add(50 * time.Millisecond) // add .05s so log1 & log2 don't show up
	since := u.Format(time.RFC3339Nano)

	out, err = d.Cmd("service", "logs", "-t", fmt.Sprintf("--since=%v", since), name)
	c.Assert(err, checker.IsNil)

	unexpected := []string***REMOVED***"log1", "log2"***REMOVED***
	expected := []string***REMOVED***"log3"***REMOVED***
	for _, v := range unexpected ***REMOVED***
		c.Assert(out, checker.Not(checker.Contains), v, check.Commentf("unexpected log message returned, since=%v", u))
	***REMOVED***
	for _, v := range expected ***REMOVED***
		c.Assert(out, checker.Contains, v, check.Commentf("expected log message %v, was not present, since=%v", u))
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsFollow(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "TestServiceLogsFollow"

	out, err := d.Cmd("service", "create", "--detach", "--no-resolve-image", "--name", name, "busybox", "sh", "-c", "while true; do echo log test; sleep 0.1; done")
	c.Assert(err, checker.IsNil)
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")

	// make sure task has been deployed.
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)

	args := []string***REMOVED***"service", "logs", "-f", name***REMOVED***
	cmd := exec.Command(dockerBinary, d.PrependHostArg(args)...)
	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	c.Assert(cmd.Start(), checker.IsNil)
	go cmd.Wait()

	// Make sure pipe is written to
	ch := make(chan *logMessage)
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		reader := bufio.NewReader(r)
		for ***REMOVED***
			msg := &logMessage***REMOVED******REMOVED***
			msg.data, _, msg.err = reader.ReadLine()
			select ***REMOVED***
			case ch <- msg:
			case <-done:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	for i := 0; i < 3; i++ ***REMOVED***
		msg := <-ch
		c.Assert(msg.err, checker.IsNil)
		c.Assert(string(msg.data), checker.Contains, "log test")
	***REMOVED***
	close(done)

	c.Assert(cmd.Process.Kill(), checker.IsNil)
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsTaskLogs(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "TestServicelogsTaskLogs"
	replicas := 2

	result := icmd.RunCmd(d.Command(
		// create a service with the name
		"service", "create", "--detach", "--no-resolve-image", "--name", name,
		// which has some number of replicas
		fmt.Sprintf("--replicas=%v", replicas),
		// which has this the task id as an environment variable templated in
		"--env", "TASK=***REMOVED******REMOVED***.Task.ID***REMOVED******REMOVED***",
		// and runs this command to print exactly 6 logs lines
		"busybox", "sh", "-c", "for line in $(seq 0 5); do echo $TASK log test $line; done; sleep 100000",
	))
	result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
	// ^^ verify that we get no error
	// then verify that we have an id in stdout
	id := strings.TrimSpace(result.Stdout())
	c.Assert(id, checker.Not(checker.Equals), "")
	// so, right here, we're basically inspecting by id and returning only
	// the ID. if they don't match, the service doesn't exist.
	result = icmd.RunCmd(d.Command("service", "inspect", "--format=\"***REMOVED******REMOVED***.ID***REMOVED******REMOVED***\"", id))
	result.Assert(c, icmd.Expected***REMOVED***Out: id***REMOVED***)

	// make sure task has been deployed.
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, replicas)
	waitAndAssert(c, defaultReconciliationTimeout, countLogLines(d, name), checker.Equals, 6*replicas)

	// get the task ids
	result = icmd.RunCmd(d.Command("service", "ps", "-q", name))
	result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
	// make sure we have two tasks
	taskIDs := strings.Split(strings.TrimSpace(result.Stdout()), "\n")
	c.Assert(taskIDs, checker.HasLen, replicas)

	for _, taskID := range taskIDs ***REMOVED***
		c.Logf("checking task %v", taskID)
		result := icmd.RunCmd(d.Command("service", "logs", taskID))
		result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
		lines := strings.Split(strings.TrimSpace(result.Stdout()), "\n")

		c.Logf("checking messages for %v", taskID)
		for i, line := range lines ***REMOVED***
			// make sure the message is in order
			c.Assert(line, checker.Contains, fmt.Sprintf("log test %v", i))
			// make sure it contains the task id
			c.Assert(line, checker.Contains, taskID)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsTTY(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "TestServiceLogsTTY"

	result := icmd.RunCmd(d.Command(
		// create a service
		"service", "create", "--detach", "--no-resolve-image",
		// name it $name
		"--name", name,
		// use a TTY
		"-t",
		// busybox image, shell string
		"busybox", "sh", "-c",
		// echo to stdout and stderr
		"echo out; (echo err 1>&2); sleep 10000",
	))

	result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
	id := strings.TrimSpace(result.Stdout())
	c.Assert(id, checker.Not(checker.Equals), "")
	// so, right here, we're basically inspecting by id and returning only
	// the ID. if they don't match, the service doesn't exist.
	result = icmd.RunCmd(d.Command("service", "inspect", "--format=\"***REMOVED******REMOVED***.ID***REMOVED******REMOVED***\"", id))
	result.Assert(c, icmd.Expected***REMOVED***Out: id***REMOVED***)

	// make sure task has been deployed.
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)
	// and make sure we have all the log lines
	waitAndAssert(c, defaultReconciliationTimeout, countLogLines(d, name), checker.Equals, 2)

	cmd := d.Command("service", "logs", "--raw", name)
	result = icmd.RunCmd(cmd)
	// for some reason there is carriage return in the output. i think this is
	// just expected.
	result.Assert(c, icmd.Expected***REMOVED***Out: "out\r\nerr\r\n"***REMOVED***)
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsNoHangDeletedContainer(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "TestServiceLogsNoHangDeletedContainer"

	result := icmd.RunCmd(d.Command(
		// create a service
		"service", "create", "--detach", "--no-resolve-image",
		// name it $name
		"--name", name,
		// busybox image, shell string
		"busybox", "sh", "-c",
		// echo to stdout and stderr
		"while true; do echo line; sleep 2; done",
	))

	// confirm that the command succeeded
	result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
	// get the service id
	id := strings.TrimSpace(result.Stdout())
	c.Assert(id, checker.Not(checker.Equals), "")

	// make sure task has been deployed.
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)
	// and make sure we have all the log lines
	waitAndAssert(c, defaultReconciliationTimeout, countLogLines(d, name), checker.Equals, 2)

	// now find and nuke the container
	result = icmd.RunCmd(d.Command("ps", "-q"))
	containerID := strings.TrimSpace(result.Stdout())
	c.Assert(containerID, checker.Not(checker.Equals), "")
	result = icmd.RunCmd(d.Command("stop", containerID))
	result.Assert(c, icmd.Expected***REMOVED***Out: containerID***REMOVED***)
	result = icmd.RunCmd(d.Command("rm", containerID))
	result.Assert(c, icmd.Expected***REMOVED***Out: containerID***REMOVED***)

	// run logs. use tail 2 to make sure we don't try to get a bunch of logs
	// somehow and slow down execution time
	cmd := d.Command("service", "logs", "--tail", "2", id)
	// start the command and then wait for it to finish with a 3 second timeout
	result = icmd.StartCmd(cmd)
	result = icmd.WaitOnCmd(3*time.Second, result)

	// then, assert that the result matches expected. if the command timed out,
	// if the command is timed out, result.Timeout will be true, but the
	// Expected defaults to false
	result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
***REMOVED***

func (s *DockerSwarmSuite) TestServiceLogsDetails(c *check.C) ***REMOVED***
	d := s.AddDaemon(c, true, true)

	name := "TestServiceLogsDetails"

	result := icmd.RunCmd(d.Command(
		// create a service
		"service", "create", "--detach", "--no-resolve-image",
		// name it $name
		"--name", name,
		// add an environment variable
		"--env", "asdf=test1",
		// add a log driver (without explicitly setting a driver, log-opt doesn't work)
		"--log-driver", "json-file",
		// add a log option to print the environment variable
		"--log-opt", "env=asdf",
		// busybox image, shell string
		"busybox", "sh", "-c",
		// make a log line
		"echo LogLine; while true; do sleep 1; done;",
	))

	result.Assert(c, icmd.Expected***REMOVED******REMOVED***)
	id := strings.TrimSpace(result.Stdout())
	c.Assert(id, checker.Not(checker.Equals), "")

	// make sure task has been deployed
	waitAndAssert(c, defaultReconciliationTimeout, d.CheckActiveContainerCount, checker.Equals, 1)
	// and make sure we have all the log lines
	waitAndAssert(c, defaultReconciliationTimeout, countLogLines(d, name), checker.Equals, 1)

	// First, test without pretty printing
	// call service logs with details. set raw to skip pretty printing
	result = icmd.RunCmd(d.Command("service", "logs", "--raw", "--details", name))
	// in this case, we should get details and we should get log message, but
	// there will also be context as details (which will fall after the detail
	// we inserted in alphabetical order
	result.Assert(c, icmd.Expected***REMOVED***Out: "asdf=test1"***REMOVED***)
	result.Assert(c, icmd.Expected***REMOVED***Out: "LogLine"***REMOVED***)

	// call service logs with details. this time, don't pass raw
	result = icmd.RunCmd(d.Command("service", "logs", "--details", id))
	// in this case, we should get details space logmessage as well. the context
	// is part of the pretty part of the logline
	result.Assert(c, icmd.Expected***REMOVED***Out: "asdf=test1 LogLine"***REMOVED***)
***REMOVED***
