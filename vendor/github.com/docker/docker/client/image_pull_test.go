package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
)

func TestImagePullReferenceParseError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			return nil, nil
		***REMOVED***),
	***REMOVED***
	// An empty reference is an invalid reference
	_, err := client.ImagePull(context.Background(), "", types.ImagePullOptions***REMOVED******REMOVED***)
	if err == nil || !strings.Contains(err.Error(), "invalid reference format") ***REMOVED***
		t.Fatalf("expected an error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePullAnyError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImagePull(context.Background(), "myimage", types.ImagePullOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePullStatusUnauthorizedError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	_, err := client.ImagePull(context.Background(), "myimage", types.ImagePullOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Unauthorized error" ***REMOVED***
		t.Fatalf("expected an Unauthorized Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePullWithUnauthorizedErrorAndPrivilegeFuncError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	privilegeFunc := func() (string, error) ***REMOVED***
		return "", fmt.Errorf("Error requesting privilege")
	***REMOVED***
	_, err := client.ImagePull(context.Background(), "myimage", types.ImagePullOptions***REMOVED***
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err == nil || err.Error() != "Error requesting privilege" ***REMOVED***
		t.Fatalf("expected an error requesting privilege, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePullWithUnauthorizedErrorAndAnotherUnauthorizedError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	privilegeFunc := func() (string, error) ***REMOVED***
		return "a-auth-header", nil
	***REMOVED***
	_, err := client.ImagePull(context.Background(), "myimage", types.ImagePullOptions***REMOVED***
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Unauthorized error" ***REMOVED***
		t.Fatalf("expected an Unauthorized Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePullWithPrivilegedFuncNoError(t *testing.T) ***REMOVED***
	expectedURL := "/images/create"
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
				return nil, fmt.Errorf("Invalid auth header : expected %s, got %s", "IAmValid", auth)
			***REMOVED***
			query := req.URL.Query()
			fromImage := query.Get("fromImage")
			if fromImage != "myimage" ***REMOVED***
				return nil, fmt.Errorf("fromimage not set in URL query properly. Expected '%s', got %s", "myimage", fromImage)
			***REMOVED***
			tag := query.Get("tag")
			if tag != "latest" ***REMOVED***
				return nil, fmt.Errorf("tag not set in URL query properly. Expected '%s', got %s", "latest", tag)
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("hello world"))),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***
	privilegeFunc := func() (string, error) ***REMOVED***
		return "IAmValid", nil
	***REMOVED***
	resp, err := client.ImagePull(context.Background(), "myimage", types.ImagePullOptions***REMOVED***
		RegistryAuth:  "NotValid",
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	body, err := ioutil.ReadAll(resp)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(body) != "hello world" ***REMOVED***
		t.Fatalf("expected 'hello world', got %s", string(body))
	***REMOVED***
***REMOVED***

func TestImagePullWithoutErrors(t *testing.T) ***REMOVED***
	expectedURL := "/images/create"
	expectedOutput := "hello world"
	pullCases := []struct ***REMOVED***
		all           bool
		reference     string
		expectedImage string
		expectedTag   string
	***REMOVED******REMOVED***
		***REMOVED***
			all:           false,
			reference:     "myimage",
			expectedImage: "myimage",
			expectedTag:   "latest",
		***REMOVED***,
		***REMOVED***
			all:           false,
			reference:     "myimage:tag",
			expectedImage: "myimage",
			expectedTag:   "tag",
		***REMOVED***,
		***REMOVED***
			all:           true,
			reference:     "myimage",
			expectedImage: "myimage",
			expectedTag:   "",
		***REMOVED***,
		***REMOVED***
			all:           true,
			reference:     "myimage:anything",
			expectedImage: "myimage",
			expectedTag:   "",
		***REMOVED***,
	***REMOVED***
	for _, pullCase := range pullCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				query := req.URL.Query()
				fromImage := query.Get("fromImage")
				if fromImage != pullCase.expectedImage ***REMOVED***
					return nil, fmt.Errorf("fromimage not set in URL query properly. Expected '%s', got %s", pullCase.expectedImage, fromImage)
				***REMOVED***
				tag := query.Get("tag")
				if tag != pullCase.expectedTag ***REMOVED***
					return nil, fmt.Errorf("tag not set in URL query properly. Expected '%s', got %s", pullCase.expectedTag, tag)
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(expectedOutput))),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***
		resp, err := client.ImagePull(context.Background(), pullCase.reference, types.ImagePullOptions***REMOVED***
			All: pullCase.all,
		***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		body, err := ioutil.ReadAll(resp)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if string(body) != expectedOutput ***REMOVED***
			t.Fatalf("expected '%s', got %s", expectedOutput, string(body))
		***REMOVED***
	***REMOVED***
***REMOVED***
