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

func TestImagePushReferenceError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			return nil, nil
		***REMOVED***),
	***REMOVED***
	// An empty reference is an invalid reference
	_, err := client.ImagePush(context.Background(), "", types.ImagePushOptions***REMOVED******REMOVED***)
	if err == nil || !strings.Contains(err.Error(), "invalid reference format") ***REMOVED***
		t.Fatalf("expected an error, got %v", err)
	***REMOVED***
	// An canonical reference cannot be pushed
	_, err = client.ImagePush(context.Background(), "repo@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", types.ImagePushOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "cannot push a digest reference" ***REMOVED***
		t.Fatalf("expected an error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePushAnyError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImagePush(context.Background(), "myimage", types.ImagePushOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePushStatusUnauthorizedError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	_, err := client.ImagePush(context.Background(), "myimage", types.ImagePushOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Unauthorized error" ***REMOVED***
		t.Fatalf("expected an Unauthorized Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePushWithUnauthorizedErrorAndPrivilegeFuncError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	privilegeFunc := func() (string, error) ***REMOVED***
		return "", fmt.Errorf("Error requesting privilege")
	***REMOVED***
	_, err := client.ImagePush(context.Background(), "myimage", types.ImagePushOptions***REMOVED***
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err == nil || err.Error() != "Error requesting privilege" ***REMOVED***
		t.Fatalf("expected an error requesting privilege, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePushWithUnauthorizedErrorAndAnotherUnauthorizedError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusUnauthorized, "Unauthorized error")),
	***REMOVED***
	privilegeFunc := func() (string, error) ***REMOVED***
		return "a-auth-header", nil
	***REMOVED***
	_, err := client.ImagePush(context.Background(), "myimage", types.ImagePushOptions***REMOVED***
		PrivilegeFunc: privilegeFunc,
	***REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Unauthorized error" ***REMOVED***
		t.Fatalf("expected an Unauthorized Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImagePushWithPrivilegedFuncNoError(t *testing.T) ***REMOVED***
	expectedURL := "/images/myimage/push"
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
			tag := query.Get("tag")
			if tag != "tag" ***REMOVED***
				return nil, fmt.Errorf("tag not set in URL query properly. Expected '%s', got %s", "tag", tag)
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
	resp, err := client.ImagePush(context.Background(), "myimage:tag", types.ImagePushOptions***REMOVED***
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

func TestImagePushWithoutErrors(t *testing.T) ***REMOVED***
	expectedOutput := "hello world"
	expectedURLFormat := "/images/%s/push"
	pullCases := []struct ***REMOVED***
		reference     string
		expectedImage string
		expectedTag   string
	***REMOVED******REMOVED***
		***REMOVED***
			reference:     "myimage",
			expectedImage: "myimage",
			expectedTag:   "",
		***REMOVED***,
		***REMOVED***
			reference:     "myimage:tag",
			expectedImage: "myimage",
			expectedTag:   "tag",
		***REMOVED***,
	***REMOVED***
	for _, pullCase := range pullCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				expectedURL := fmt.Sprintf(expectedURLFormat, pullCase.expectedImage)
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				query := req.URL.Query()
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
		resp, err := client.ImagePush(context.Background(), pullCase.reference, types.ImagePushOptions***REMOVED******REMOVED***)
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
