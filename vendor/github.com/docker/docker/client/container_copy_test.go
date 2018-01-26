package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
)

func TestContainerStatPathError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerStatPath(context.Background(), "container_id", "path")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerStatPathNotFoundError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Not found")),
	***REMOVED***
	_, err := client.ContainerStatPath(context.Background(), "container_id", "path")
	if !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a not found error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerStatPathNoHeaderError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	_, err := client.ContainerStatPath(context.Background(), "container_id", "path/to/file")
	if err == nil ***REMOVED***
		t.Fatalf("expected an error, got nothing")
	***REMOVED***
***REMOVED***

func TestContainerStatPath(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/archive"
	expectedPath := "path/to/file"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "HEAD" ***REMOVED***
				return nil, fmt.Errorf("expected HEAD method, got %s", req.Method)
			***REMOVED***
			query := req.URL.Query()
			path := query.Get("path")
			if path != expectedPath ***REMOVED***
				return nil, fmt.Errorf("path not set in URL query properly")
			***REMOVED***
			content, err := json.Marshal(types.ContainerPathStat***REMOVED***
				Name: "name",
				Mode: 0700,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			base64PathStat := base64.StdEncoding.EncodeToString(content)
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				Header: http.Header***REMOVED***
					"X-Docker-Container-Path-Stat": []string***REMOVED***base64PathStat***REMOVED***,
				***REMOVED***,
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	stat, err := client.ContainerStatPath(context.Background(), "container_id", expectedPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if stat.Name != "name" ***REMOVED***
		t.Fatalf("expected container path stat name to be 'name', got '%s'", stat.Name)
	***REMOVED***
	if stat.Mode != 0700 ***REMOVED***
		t.Fatalf("expected container path stat mode to be 0700, got '%v'", stat.Mode)
	***REMOVED***
***REMOVED***

func TestCopyToContainerError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	err := client.CopyToContainer(context.Background(), "container_id", "path/to/file", bytes.NewReader([]byte("")), types.CopyToContainerOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCopyToContainerNotFoundError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Not found")),
	***REMOVED***
	err := client.CopyToContainer(context.Background(), "container_id", "path/to/file", bytes.NewReader([]byte("")), types.CopyToContainerOptions***REMOVED******REMOVED***)
	if !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a not found error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCopyToContainerNotStatusOKError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNoContent, "No content")),
	***REMOVED***
	err := client.CopyToContainer(context.Background(), "container_id", "path/to/file", bytes.NewReader([]byte("")), types.CopyToContainerOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "unexpected status code from daemon: 204" ***REMOVED***
		t.Fatalf("expected an unexpected status code error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCopyToContainer(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/archive"
	expectedPath := "path/to/file"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "PUT" ***REMOVED***
				return nil, fmt.Errorf("expected PUT method, got %s", req.Method)
			***REMOVED***
			query := req.URL.Query()
			path := query.Get("path")
			if path != expectedPath ***REMOVED***
				return nil, fmt.Errorf("path not set in URL query properly, expected '%s', got %s", expectedPath, path)
			***REMOVED***
			noOverwriteDirNonDir := query.Get("noOverwriteDirNonDir")
			if noOverwriteDirNonDir != "true" ***REMOVED***
				return nil, fmt.Errorf("noOverwriteDirNonDir not set in URL query properly, expected true, got %s", noOverwriteDirNonDir)
			***REMOVED***

			content, err := ioutil.ReadAll(req.Body)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if err := req.Body.Close(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if string(content) != "content" ***REMOVED***
				return nil, fmt.Errorf("expected content to be 'content', got %s", string(content))
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	err := client.CopyToContainer(context.Background(), "container_id", expectedPath, bytes.NewReader([]byte("content")), types.CopyToContainerOptions***REMOVED***
		AllowOverwriteDirWithFile: false,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestCopyFromContainerError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, _, err := client.CopyFromContainer(context.Background(), "container_id", "path/to/file")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCopyFromContainerNotFoundError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNotFound, "Not found")),
	***REMOVED***
	_, _, err := client.CopyFromContainer(context.Background(), "container_id", "path/to/file")
	if !IsErrNotFound(err) ***REMOVED***
		t.Fatalf("expected a not found error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCopyFromContainerNotStatusOKError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusNoContent, "No content")),
	***REMOVED***
	_, _, err := client.CopyFromContainer(context.Background(), "container_id", "path/to/file")
	if err == nil || err.Error() != "unexpected status code from daemon: 204" ***REMOVED***
		t.Fatalf("expected an unexpected status code error, got %v", err)
	***REMOVED***
***REMOVED***

func TestCopyFromContainerNoHeaderError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	_, _, err := client.CopyFromContainer(context.Background(), "container_id", "path/to/file")
	if err == nil ***REMOVED***
		t.Fatalf("expected an error, got nothing")
	***REMOVED***
***REMOVED***

func TestCopyFromContainer(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/archive"
	expectedPath := "path/to/file"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "GET" ***REMOVED***
				return nil, fmt.Errorf("expected GET method, got %s", req.Method)
			***REMOVED***
			query := req.URL.Query()
			path := query.Get("path")
			if path != expectedPath ***REMOVED***
				return nil, fmt.Errorf("path not set in URL query properly, expected '%s', got %s", expectedPath, path)
			***REMOVED***

			headercontent, err := json.Marshal(types.ContainerPathStat***REMOVED***
				Name: "name",
				Mode: 0700,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			base64PathStat := base64.StdEncoding.EncodeToString(headercontent)

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("content"))),
				Header: http.Header***REMOVED***
					"X-Docker-Container-Path-Stat": []string***REMOVED***base64PathStat***REMOVED***,
				***REMOVED***,
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	r, stat, err := client.CopyFromContainer(context.Background(), "container_id", expectedPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if stat.Name != "name" ***REMOVED***
		t.Fatalf("expected container path stat name to be 'name', got '%s'", stat.Name)
	***REMOVED***
	if stat.Mode != 0700 ***REMOVED***
		t.Fatalf("expected container path stat mode to be 0700, got '%v'", stat.Mode)
	***REMOVED***
	content, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := r.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(content) != "content" ***REMOVED***
		t.Fatalf("expected content to be 'content', got %s", string(content))
	***REMOVED***
***REMOVED***
