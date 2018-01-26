package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
)

func TestContainerExecCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerExecCreate(context.Background(), "container_id", types.ExecConfig***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerExecCreate(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/exec"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***
			// FIXME validate the content is the given ExecConfig ?
			if err := req.ParseForm(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			execConfig := &types.ExecConfig***REMOVED******REMOVED***
			if err := json.NewDecoder(req.Body).Decode(execConfig); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if execConfig.User != "user" ***REMOVED***
				return nil, fmt.Errorf("expected an execConfig with User == 'user', got %v", execConfig)
			***REMOVED***
			b, err := json.Marshal(types.IDResponse***REMOVED***
				ID: "exec_id",
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	r, err := client.ContainerExecCreate(context.Background(), "container_id", types.ExecConfig***REMOVED***
		User: "user",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.ID != "exec_id" ***REMOVED***
		t.Fatalf("expected `exec_id`, got %s", r.ID)
	***REMOVED***
***REMOVED***

func TestContainerExecStartError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.ContainerExecStart(context.Background(), "nothing", types.ExecStartCheck***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerExecStart(t *testing.T) ***REMOVED***
	expectedURL := "/exec/exec_id/start"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if err := req.ParseForm(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			execStartCheck := &types.ExecStartCheck***REMOVED******REMOVED***
			if err := json.NewDecoder(req.Body).Decode(execStartCheck); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if execStartCheck.Tty || !execStartCheck.Detach ***REMOVED***
				return nil, fmt.Errorf("expected execStartCheck***REMOVED***Detach:true,Tty:false***REMOVED***, got %v", execStartCheck)
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	err := client.ContainerExecStart(context.Background(), "exec_id", types.ExecStartCheck***REMOVED***
		Detach: true,
		Tty:    false,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestContainerExecInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerExecInspect(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerExecInspect(t *testing.T) ***REMOVED***
	expectedURL := "/exec/exec_id/json"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			b, err := json.Marshal(types.ContainerExecInspect***REMOVED***
				ExecID:      "exec_id",
				ContainerID: "container_id",
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	inspect, err := client.ContainerExecInspect(context.Background(), "exec_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if inspect.ExecID != "exec_id" ***REMOVED***
		t.Fatalf("expected ExecID to be `exec_id`, got %s", inspect.ExecID)
	***REMOVED***
	if inspect.ContainerID != "container_id" ***REMOVED***
		t.Fatalf("expected ContainerID `container_id`, got %s", inspect.ContainerID)
	***REMOVED***
***REMOVED***
