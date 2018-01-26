package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"golang.org/x/net/context"
)

func TestContainerCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerCreate(context.Background(), nil, nil, nil, "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error while testing StatusInternalServerError, got %v", err)
	***REMOVED***

	// 404 doesn't automatically means an unknown image
	client = &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Server error")),
	***REMOVED***
	_, err = client.ContainerCreate(context.Background(), nil, nil, nil, "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error while testing StatusNotFound, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerCreateImageNotFound(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "No such image")),
	***REMOVED***
	_, err := client.ContainerCreate(context.Background(), &container.Config***REMOVED***Image: "unknown_image"***REMOVED***, nil, nil, "unknown")
	if err == nil || !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected an imageNotFound error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerCreateWithName(t *testing.T) ***REMOVED***
	expectedURL := "/containers/create"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			name := req.URL.Query().Get("name")
			if name != "container_name" ***REMOVED***
				return nil, fmt.Errorf("container name not set in URL query properly. Expected `container_name`, got %s", name)
			***REMOVED***
			b, err := json.Marshal(container.ContainerCreateCreatedBody***REMOVED***
				ID: "container_id",
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

	r, err := client.ContainerCreate(context.Background(), nil, nil, nil, "container_name")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.ID != "container_id" ***REMOVED***
		t.Fatalf("expected `container_id`, got %s", r.ID)
	***REMOVED***
***REMOVED***

// TestContainerCreateAutoRemove validates that a client using API 1.24 always disables AutoRemove. When using API 1.25
// or up, AutoRemove should not be disabled.
func TestContainerCreateAutoRemove(t *testing.T) ***REMOVED***
	autoRemoveValidator := func(expectedValue bool) func(req *http.Request) (*http.Response, error) ***REMOVED***
		return func(req *http.Request) (*http.Response, error) ***REMOVED***
			var config configWrapper

			if err := json.NewDecoder(req.Body).Decode(&config); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if config.HostConfig.AutoRemove != expectedValue ***REMOVED***
				return nil, fmt.Errorf("expected AutoRemove to be %v, got %v", expectedValue, config.HostConfig.AutoRemove)
			***REMOVED***
			b, err := json.Marshal(container.ContainerCreateCreatedBody***REMOVED***
				ID: "container_id",
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***
	***REMOVED***

	client := &Client***REMOVED***
		client:  newMockClient(autoRemoveValidator(false)),
		version: "1.24",
	***REMOVED***
	if _, err := client.ContainerCreate(context.Background(), nil, &container.HostConfig***REMOVED***AutoRemove: true***REMOVED***, nil, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	client = &Client***REMOVED***
		client:  newMockClient(autoRemoveValidator(true)),
		version: "1.25",
	***REMOVED***
	if _, err := client.ContainerCreate(context.Background(), nil, &container.HostConfig***REMOVED***AutoRemove: true***REMOVED***, nil, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
