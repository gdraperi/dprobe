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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"golang.org/x/net/context"
)

func TestImageSearchAnyError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImageSearch(context.Background(), "some-image", types.ImageSearchOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageSearchStatusUnauthorizedError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	_, err := client.ImageSearch(context.Background(), "some-image", types.ImageSearchOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Unauthorized error" ***REMOVED***
		t.Fatalf("expected an Unauthorized Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageSearchWithUnauthorizedErrorAndPrivilegeFuncError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	privilegeFunc := func() (string, error) ***REMOVED***
		return "", fmt.Errorf("Error requesting privilege")
	***REMOVED***
	_, err := client.ImageSearch(context.Background(), "some-image", types.ImageSearchOptions***REMOVED***
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err == nil || err.Error() != "Error requesting privilege" ***REMOVED***
		t.Fatalf("expected an error requesting privilege, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageSearchWithUnauthorizedErrorAndAnotherUnauthorizedError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	privilegeFunc := func() (string, error) ***REMOVED***
		return "a-auth-header", nil
	***REMOVED***
	_, err := client.ImageSearch(context.Background(), "some-image", types.ImageSearchOptions***REMOVED***
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Unauthorized error" ***REMOVED***
		t.Fatalf("expected an Unauthorized Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageSearchWithPrivilegedFuncNoError(t *testing.T) ***REMOVED***
	expectedURL := "/images/search"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			auth := req.Header.Get("X-Registry-Auth")
			if auth == "NotValid" ***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusUnauthorized,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("Invalid credentials"))),
				***REMOVED***, nil
			***REMOVED***
			if auth != "IAmValid" ***REMOVED***
				return nil, fmt.Errorf("Invalid auth header : expected 'IAmValid', got %s", auth)
			***REMOVED***
			query := req.URL.Query()
			term := query.Get("term")
			if term != "some-image" ***REMOVED***
				return nil, fmt.Errorf("term not set in URL query properly. Expected 'some-image', got %s", term)
			***REMOVED***
			content, err := json.Marshal([]registry.SearchResult***REMOVED***
				***REMOVED***
					Name: "anything",
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
	privilegeFunc := func() (string, error) ***REMOVED***
		return "IAmValid", nil
	***REMOVED***
	results, err := client.ImageSearch(context.Background(), "some-image", types.ImageSearchOptions***REMOVED***
		RegistryAuth:  "NotValid",
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(results) != 1 ***REMOVED***
		t.Fatalf("expected 1 result, got %v", results)
	***REMOVED***
***REMOVED***

func TestImageSearchWithoutErrors(t *testing.T) ***REMOVED***
	expectedURL := "/images/search"
	filterArgs := filters.NewArgs()
	filterArgs.Add("is-automated", "true")
	filterArgs.Add("stars", "3")

	expectedFilters := `***REMOVED***"is-automated":***REMOVED***"true":true***REMOVED***,"stars":***REMOVED***"3":true***REMOVED******REMOVED***`

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			query := req.URL.Query()
			term := query.Get("term")
			if term != "some-image" ***REMOVED***
				return nil, fmt.Errorf("term not set in URL query properly. Expected 'some-image', got %s", term)
			***REMOVED***
			filters := query.Get("filters")
			if filters != expectedFilters ***REMOVED***
				return nil, fmt.Errorf("filters not set in URL query properly. Expected '%s', got %s", expectedFilters, filters)
			***REMOVED***
			content, err := json.Marshal([]registry.SearchResult***REMOVED***
				***REMOVED***
					Name: "anything",
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
	results, err := client.ImageSearch(context.Background(), "some-image", types.ImageSearchOptions***REMOVED***
		Filters: filterArgs,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(results) != 1 ***REMOVED***
		t.Fatalf("expected a result, got %v", results)
	***REMOVED***
***REMOVED***
