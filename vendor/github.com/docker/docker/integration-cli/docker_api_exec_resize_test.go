package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestExecResizeAPIHeightWidthNoInt(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "top")
	cleanedContainerID := strings.TrimSpace(out)

	endpoint := "/exec/" + cleanedContainerID + "/resize?h=foo&w=bar"
	res, _, err := request.Post(endpoint)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusBadRequest)
***REMOVED***

// Part of #14845
func (s *DockerSuite) TestExecResizeImmediatelyAfterExecStart(c *check.C) ***REMOVED***
	name := "exec_resize_test"
	dockerCmd(c, "run", "-d", "-i", "-t", "--name", name, "--restart", "always", "busybox", "/bin/sh")

	testExecResize := func() error ***REMOVED***
		data := map[string]interface***REMOVED******REMOVED******REMOVED***
			"AttachStdin": true,
			"Cmd":         []string***REMOVED***"/bin/sh"***REMOVED***,
		***REMOVED***
		uri := fmt.Sprintf("/containers/%s/exec", name)
		res, body, err := request.Post(uri, request.JSONBody(data))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if res.StatusCode != http.StatusCreated ***REMOVED***
			return fmt.Errorf("POST %s is expected to return %d, got %d", uri, http.StatusCreated, res.StatusCode)
		***REMOVED***

		buf, err := request.ReadBody(body)
		c.Assert(err, checker.IsNil)

		out := map[string]string***REMOVED******REMOVED***
		err = json.Unmarshal(buf, &out)
		if err != nil ***REMOVED***
			return fmt.Errorf("ExecCreate returned invalid json. Error: %q", err.Error())
		***REMOVED***

		execID := out["Id"]
		if len(execID) < 1 ***REMOVED***
			return fmt.Errorf("ExecCreate got invalid execID")
		***REMOVED***

		payload := bytes.NewBufferString(`***REMOVED***"Tty":true***REMOVED***`)
		conn, _, err := request.SockRequestHijack("POST", fmt.Sprintf("/exec/%s/start", execID), payload, "application/json", daemonHost())
		if err != nil ***REMOVED***
			return fmt.Errorf("Failed to start the exec: %q", err.Error())
		***REMOVED***
		defer conn.Close()

		_, rc, err := request.Post(fmt.Sprintf("/exec/%s/resize?h=24&w=80", execID), request.ContentType("text/plain"))
		// It's probably a panic of the daemon if io.ErrUnexpectedEOF is returned.
		if err == io.ErrUnexpectedEOF ***REMOVED***
			return fmt.Errorf("The daemon might have crashed.")
		***REMOVED***

		if err == nil ***REMOVED***
			rc.Close()
		***REMOVED***

		// We only interested in the io.ErrUnexpectedEOF error, so we return nil otherwise.
		return nil
	***REMOVED***

	// The panic happens when daemon.ContainerExecStart is called but the
	// container.Exec is not called.
	// Because the panic is not 100% reproducible, we send the requests concurrently
	// to increase the probability that the problem is triggered.
	var (
		n  = 10
		ch = make(chan error, n)
		wg sync.WaitGroup
	)
	for i := 0; i < n; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			if err := testExecResize(); err != nil ***REMOVED***
				ch <- err
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	wg.Wait()
	select ***REMOVED***
	case err := <-ch:
		c.Fatal(err.Error())
	default:
	***REMOVED***
***REMOVED***
