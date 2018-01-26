package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func TestContainerInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, err := client.ContainerInspect(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerInspectContainerNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***

	_, err := client.ContainerInspect(context.Background(), "unknown")
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a containerNotFound error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerInspect(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/json"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(types.ContainerJSON***REMOVED***
				ContainerJSONBase: &types.ContainerJSONBase***REMOVED***
					ID:    "container_id",
					Image: "image",
					Name:  "name",
				***REMOVED***,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	r, err := client.ContainerInspect(context.Background(), "container_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.ID != "container_id" ***REMOVED***
		t.Fatalf("expected `container_id`, got %s", r.ID)
	***REMOVED***
	if r.Image != "image" ***REMOVED***
		t.Fatalf("expected `image`, got %s", r.Image)
	***REMOVED***
	if r.Name != "name" ***REMOVED***
		t.Fatalf("expected `name`, got %s", r.Name)
	***REMOVED***
***REMOVED***

func TestContainerInspectNode(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			content, err := json.Marshal(types.ContainerJSON***REMOVED***
				ContainerJSONBase: &types.ContainerJSONBase***REMOVED***
					ID:    "container_id",
					Image: "image",
					Name:  "name",
					Node: &types.ContainerNode***REMOVED***
						ID:     "container_node_id",
						Addr:   "container_node",
						Labels: map[string]string***REMOVED***"foo": "bar"***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	r, err := client.ContainerInspect(context.Background(), "container_id")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.ID != "container_id" ***REMOVED***
		t.Fatalf("expected `container_id`, got %s", r.ID)
	***REMOVED***
	if r.Image != "image" ***REMOVED***
		t.Fatalf("expected `image`, got %s", r.Image)
	***REMOVED***
	if r.Name != "name" ***REMOVED***
		t.Fatalf("expected `name`, got %s", r.Name)
	***REMOVED***
	if r.Node.ID != "container_node_id" ***REMOVED***
		t.Fatalf("expected `container_node_id`, got %s", r.Node.ID)
	***REMOVED***
	if r.Node.Addr != "container_node" ***REMOVED***
		t.Fatalf("expected `container_node`, got %s", r.Node.Addr)
	***REMOVED***
	foo, ok := r.Node.Labels["foo"]
	if foo != "bar" || !ok ***REMOVED***
		t.Fatalf("expected `bar` for label `foo`")
	***REMOVED***
***REMOVED***
