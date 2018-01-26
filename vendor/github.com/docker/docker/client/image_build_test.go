package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-units"
)

func TestImageBuildError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ImageBuild(context.Background(), nil, types.ImageBuildOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestImageBuild(t *testing.T) ***REMOVED***
	v1 := "value1"
	v2 := "value2"
	emptyRegistryConfig := "bnVsbA=="
	buildCases := []struct ***REMOVED***
		buildOptions           types.ImageBuildOptions
		expectedQueryParams    map[string]string
		expectedTags           []string
		expectedRegistryConfig string
	***REMOVED******REMOVED***
		***REMOVED***
			buildOptions: types.ImageBuildOptions***REMOVED***
				SuppressOutput: true,
				NoCache:        true,
				Remove:         true,
				ForceRemove:    true,
				PullParent:     true,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"q":       "1",
				"nocache": "1",
				"rm":      "1",
				"forcerm": "1",
				"pull":    "1",
			***REMOVED***,
			expectedTags:           []string***REMOVED******REMOVED***,
			expectedRegistryConfig: emptyRegistryConfig,
		***REMOVED***,
		***REMOVED***
			buildOptions: types.ImageBuildOptions***REMOVED***
				SuppressOutput: false,
				NoCache:        false,
				Remove:         false,
				ForceRemove:    false,
				PullParent:     false,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"q":       "",
				"nocache": "",
				"rm":      "0",
				"forcerm": "",
				"pull":    "",
			***REMOVED***,
			expectedTags:           []string***REMOVED******REMOVED***,
			expectedRegistryConfig: emptyRegistryConfig,
		***REMOVED***,
		***REMOVED***
			buildOptions: types.ImageBuildOptions***REMOVED***
				RemoteContext: "remoteContext",
				Isolation:     container.Isolation("isolation"),
				CPUSetCPUs:    "2",
				CPUSetMems:    "12",
				CPUShares:     20,
				CPUQuota:      10,
				CPUPeriod:     30,
				Memory:        256,
				MemorySwap:    512,
				ShmSize:       10,
				CgroupParent:  "cgroup_parent",
				Dockerfile:    "Dockerfile",
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"remote":       "remoteContext",
				"isolation":    "isolation",
				"cpusetcpus":   "2",
				"cpusetmems":   "12",
				"cpushares":    "20",
				"cpuquota":     "10",
				"cpuperiod":    "30",
				"memory":       "256",
				"memswap":      "512",
				"shmsize":      "10",
				"cgroupparent": "cgroup_parent",
				"dockerfile":   "Dockerfile",
				"rm":           "0",
			***REMOVED***,
			expectedTags:           []string***REMOVED******REMOVED***,
			expectedRegistryConfig: emptyRegistryConfig,
		***REMOVED***,
		***REMOVED***
			buildOptions: types.ImageBuildOptions***REMOVED***
				BuildArgs: map[string]*string***REMOVED***
					"ARG1": &v1,
					"ARG2": &v2,
					"ARG3": nil,
				***REMOVED***,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"buildargs": `***REMOVED***"ARG1":"value1","ARG2":"value2","ARG3":null***REMOVED***`,
				"rm":        "0",
			***REMOVED***,
			expectedTags:           []string***REMOVED******REMOVED***,
			expectedRegistryConfig: emptyRegistryConfig,
		***REMOVED***,
		***REMOVED***
			buildOptions: types.ImageBuildOptions***REMOVED***
				Ulimits: []*units.Ulimit***REMOVED***
					***REMOVED***
						Name: "nproc",
						Hard: 65557,
						Soft: 65557,
					***REMOVED***,
					***REMOVED***
						Name: "nofile",
						Hard: 20000,
						Soft: 40000,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"ulimits": `[***REMOVED***"Name":"nproc","Hard":65557,"Soft":65557***REMOVED***,***REMOVED***"Name":"nofile","Hard":20000,"Soft":40000***REMOVED***]`,
				"rm":      "0",
			***REMOVED***,
			expectedTags:           []string***REMOVED******REMOVED***,
			expectedRegistryConfig: emptyRegistryConfig,
		***REMOVED***,
		***REMOVED***
			buildOptions: types.ImageBuildOptions***REMOVED***
				AuthConfigs: map[string]types.AuthConfig***REMOVED***
					"https://index.docker.io/v1/": ***REMOVED***
						Auth: "dG90bwo=",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"rm": "0",
			***REMOVED***,
			expectedTags:           []string***REMOVED******REMOVED***,
			expectedRegistryConfig: "eyJodHRwczovL2luZGV4LmRvY2tlci5pby92MS8iOnsiYXV0aCI6ImRHOTBid289In19",
		***REMOVED***,
	***REMOVED***
	for _, buildCase := range buildCases ***REMOVED***
		expectedURL := "/build"
		client := &Client***REMOVED***
			client: newMockClient(func(r *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(r.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, r.URL)
				***REMOVED***
				// Check request headers
				registryConfig := r.Header.Get("X-Registry-Config")
				if registryConfig != buildCase.expectedRegistryConfig ***REMOVED***
					return nil, fmt.Errorf("X-Registry-Config header not properly set in the request. Expected '%s', got %s", buildCase.expectedRegistryConfig, registryConfig)
				***REMOVED***
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/x-tar" ***REMOVED***
					return nil, fmt.Errorf("Content-type header not properly set in the request. Expected 'application/x-tar', got %s", contentType)
				***REMOVED***

				// Check query parameters
				query := r.URL.Query()
				for key, expected := range buildCase.expectedQueryParams ***REMOVED***
					actual := query.Get(key)
					if actual != expected ***REMOVED***
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					***REMOVED***
				***REMOVED***

				// Check tags
				if len(buildCase.expectedTags) > 0 ***REMOVED***
					tags := query["t"]
					if !reflect.DeepEqual(tags, buildCase.expectedTags) ***REMOVED***
						return nil, fmt.Errorf("t (tags) not set in URL query properly. Expected '%s', got %s", buildCase.expectedTags, tags)
					***REMOVED***
				***REMOVED***

				headers := http.Header***REMOVED******REMOVED***
				headers.Add("Server", "Docker/v1.23 (MyOS)")
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("body"))),
					Header:     headers,
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***
		buildResponse, err := client.ImageBuild(context.Background(), nil, buildCase.buildOptions)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if buildResponse.OSType != "MyOS" ***REMOVED***
			t.Fatalf("expected OSType to be 'MyOS', got %s", buildResponse.OSType)
		***REMOVED***
		response, err := ioutil.ReadAll(buildResponse.Body)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		buildResponse.Body.Close()
		if string(response) != "body" ***REMOVED***
			t.Fatalf("expected Body to contain 'body' string, got %s", response)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGetDockerOS(t *testing.T) ***REMOVED***
	cases := map[string]string***REMOVED***
		"Docker/v1.22 (linux)":   "linux",
		"Docker/v1.22 (windows)": "windows",
		"Foo/v1.22 (bar)":        "",
	***REMOVED***
	for header, os := range cases ***REMOVED***
		g := getDockerOS(header)
		if g != os ***REMOVED***
			t.Fatalf("Expected %s, got %s", os, g)
		***REMOVED***
	***REMOVED***
***REMOVED***
