package client

import (
	"bytes"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvClient(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("skipping unix only test for windows")
	***REMOVED***
	cases := []struct ***REMOVED***
		envs            map[string]string
		expectedError   string
		expectedVersion string
	***REMOVED******REMOVED***
		***REMOVED***
			envs:            map[string]string***REMOVED******REMOVED***,
			expectedVersion: api.DefaultVersion,
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_CERT_PATH": "invalid/path",
			***REMOVED***,
			expectedError: "Could not load X509 key pair: open invalid/path/cert.pem: no such file or directory",
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_CERT_PATH": "testdata/",
			***REMOVED***,
			expectedVersion: api.DefaultVersion,
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_CERT_PATH":  "testdata/",
				"DOCKER_TLS_VERIFY": "1",
			***REMOVED***,
			expectedVersion: api.DefaultVersion,
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_CERT_PATH": "testdata/",
				"DOCKER_HOST":      "https://notaunixsocket",
			***REMOVED***,
			expectedVersion: api.DefaultVersion,
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_HOST": "host",
			***REMOVED***,
			expectedError: "unable to parse docker host `host`",
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_HOST": "invalid://url",
			***REMOVED***,
			expectedVersion: api.DefaultVersion,
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_API_VERSION": "anything",
			***REMOVED***,
			expectedVersion: "anything",
		***REMOVED***,
		***REMOVED***
			envs: map[string]string***REMOVED***
				"DOCKER_API_VERSION": "1.22",
			***REMOVED***,
			expectedVersion: "1.22",
		***REMOVED***,
	***REMOVED***

	env := envToMap()
	defer mapToEnv(env)
	for _, c := range cases ***REMOVED***
		mapToEnv(env)
		mapToEnv(c.envs)
		apiclient, err := NewEnvClient()
		if c.expectedError != "" ***REMOVED***
			assert.Error(t, err)
			assert.Equal(t, c.expectedError, err.Error())
		***REMOVED*** else ***REMOVED***
			assert.NoError(t, err)
			version := apiclient.ClientVersion()
			assert.Equal(t, c.expectedVersion, version)
		***REMOVED***

		if c.envs["DOCKER_TLS_VERIFY"] != "" ***REMOVED***
			// pedantic checking that this is handled correctly
			tr := apiclient.client.Transport.(*http.Transport)
			assert.NotNil(t, tr.TLSClientConfig)
			assert.Equal(t, tr.TLSClientConfig.InsecureSkipVerify, false)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGetAPIPath(t *testing.T) ***REMOVED***
	testcases := []struct ***REMOVED***
		version  string
		path     string
		query    url.Values
		expected string
	***REMOVED******REMOVED***
		***REMOVED***"", "/containers/json", nil, "/containers/json"***REMOVED***,
		***REMOVED***"", "/containers/json", url.Values***REMOVED******REMOVED***, "/containers/json"***REMOVED***,
		***REMOVED***"", "/containers/json", url.Values***REMOVED***"s": []string***REMOVED***"c"***REMOVED******REMOVED***, "/containers/json?s=c"***REMOVED***,
		***REMOVED***"1.22", "/containers/json", nil, "/v1.22/containers/json"***REMOVED***,
		***REMOVED***"1.22", "/containers/json", url.Values***REMOVED******REMOVED***, "/v1.22/containers/json"***REMOVED***,
		***REMOVED***"1.22", "/containers/json", url.Values***REMOVED***"s": []string***REMOVED***"c"***REMOVED******REMOVED***, "/v1.22/containers/json?s=c"***REMOVED***,
		***REMOVED***"v1.22", "/containers/json", nil, "/v1.22/containers/json"***REMOVED***,
		***REMOVED***"v1.22", "/containers/json", url.Values***REMOVED******REMOVED***, "/v1.22/containers/json"***REMOVED***,
		***REMOVED***"v1.22", "/containers/json", url.Values***REMOVED***"s": []string***REMOVED***"c"***REMOVED******REMOVED***, "/v1.22/containers/json?s=c"***REMOVED***,
		***REMOVED***"v1.22", "/networks/kiwl$%^", nil, "/v1.22/networks/kiwl$%25%5E"***REMOVED***,
	***REMOVED***

	for _, testcase := range testcases ***REMOVED***
		c := Client***REMOVED***version: testcase.version, basePath: "/"***REMOVED***
		actual := c.getAPIPath(testcase.path, testcase.query)
		assert.Equal(t, actual, testcase.expected)
	***REMOVED***
***REMOVED***

func TestParseHost(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		host  string
		proto string
		addr  string
		base  string
		err   bool
	***REMOVED******REMOVED***
		***REMOVED***"", "", "", "", true***REMOVED***,
		***REMOVED***"foobar", "", "", "", true***REMOVED***,
		***REMOVED***"foo://bar", "foo", "bar", "", false***REMOVED***,
		***REMOVED***"tcp://localhost:2476", "tcp", "localhost:2476", "", false***REMOVED***,
		***REMOVED***"tcp://localhost:2476/path", "tcp", "localhost:2476", "/path", false***REMOVED***,
	***REMOVED***

	for _, cs := range cases ***REMOVED***
		p, a, b, e := ParseHost(cs.host)
		if cs.err ***REMOVED***
			assert.Error(t, e)
		***REMOVED***
		assert.Equal(t, cs.proto, p)
		assert.Equal(t, cs.addr, a)
		assert.Equal(t, cs.base, b)
	***REMOVED***
***REMOVED***

func TestParseHostURL(t *testing.T) ***REMOVED***
	testcases := []struct ***REMOVED***
		host        string
		expected    *url.URL
		expectedErr string
	***REMOVED******REMOVED***
		***REMOVED***
			host:        "",
			expectedErr: "unable to parse docker host",
		***REMOVED***,
		***REMOVED***
			host:        "foobar",
			expectedErr: "unable to parse docker host",
		***REMOVED***,
		***REMOVED***
			host:     "foo://bar",
			expected: &url.URL***REMOVED***Scheme: "foo", Host: "bar"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			host:     "tcp://localhost:2476",
			expected: &url.URL***REMOVED***Scheme: "tcp", Host: "localhost:2476"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			host:     "tcp://localhost:2476/path",
			expected: &url.URL***REMOVED***Scheme: "tcp", Host: "localhost:2476", Path: "/path"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, testcase := range testcases ***REMOVED***
		actual, err := ParseHostURL(testcase.host)
		if testcase.expectedErr != "" ***REMOVED***
			testutil.ErrorContains(t, err, testcase.expectedErr)
		***REMOVED***
		assert.Equal(t, testcase.expected, actual)
	***REMOVED***
***REMOVED***

func TestNewEnvClientSetsDefaultVersion(t *testing.T) ***REMOVED***
	env := envToMap()
	defer mapToEnv(env)

	envMap := map[string]string***REMOVED***
		"DOCKER_HOST":        "",
		"DOCKER_API_VERSION": "",
		"DOCKER_TLS_VERIFY":  "",
		"DOCKER_CERT_PATH":   "",
	***REMOVED***
	mapToEnv(envMap)

	client, err := NewEnvClient()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, client.version, api.DefaultVersion)

	expected := "1.22"
	os.Setenv("DOCKER_API_VERSION", expected)
	client, err = NewEnvClient()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, expected, client.version)
***REMOVED***

// TestNegotiateAPIVersionEmpty asserts that client.Client can
// negotiate a compatible APIVersion when omitted
func TestNegotiateAPIVersionEmpty(t *testing.T) ***REMOVED***
	env := envToMap()
	defer mapToEnv(env)

	envMap := map[string]string***REMOVED***
		"DOCKER_API_VERSION": "",
	***REMOVED***
	mapToEnv(envMap)

	client, err := NewEnvClient()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ping := types.Ping***REMOVED***
		APIVersion:   "",
		OSType:       "linux",
		Experimental: false,
	***REMOVED***

	// set our version to something new
	client.version = "1.25"

	// if no version from server, expect the earliest
	// version before APIVersion was implemented
	expected := "1.24"

	// test downgrade
	client.NegotiateAPIVersionPing(ping)
	assert.Equal(t, expected, client.version)
***REMOVED***

// TestNegotiateAPIVersion asserts that client.Client can
// negotiate a compatible APIVersion with the server
func TestNegotiateAPIVersion(t *testing.T) ***REMOVED***
	client, err := NewEnvClient()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected := "1.21"

	ping := types.Ping***REMOVED***
		APIVersion:   expected,
		OSType:       "linux",
		Experimental: false,
	***REMOVED***

	// set our version to something new
	client.version = "1.22"

	// test downgrade
	client.NegotiateAPIVersionPing(ping)
	assert.Equal(t, expected, client.version)

	// set the client version to something older, and verify that we keep the
	// original setting.
	expected = "1.20"
	client.version = expected
	client.NegotiateAPIVersionPing(ping)
	assert.Equal(t, expected, client.version)

***REMOVED***

// TestNegotiateAPIVersionOverride asserts that we honor
// the environment variable DOCKER_API_VERSION when negotianing versions
func TestNegotiateAPVersionOverride(t *testing.T) ***REMOVED***
	env := envToMap()
	defer mapToEnv(env)

	envMap := map[string]string***REMOVED***
		"DOCKER_API_VERSION": "9.99",
	***REMOVED***
	mapToEnv(envMap)

	client, err := NewEnvClient()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ping := types.Ping***REMOVED***
		APIVersion:   "1.24",
		OSType:       "linux",
		Experimental: false,
	***REMOVED***

	expected := envMap["DOCKER_API_VERSION"]

	// test that we honored the env var
	client.NegotiateAPIVersionPing(ping)
	assert.Equal(t, expected, client.version)
***REMOVED***

// mapToEnv takes a map of environment variables and sets them
func mapToEnv(env map[string]string) ***REMOVED***
	for k, v := range env ***REMOVED***
		os.Setenv(k, v)
	***REMOVED***
***REMOVED***

// envToMap returns a map of environment variables
func envToMap() map[string]string ***REMOVED***
	env := make(map[string]string)
	for _, e := range os.Environ() ***REMOVED***
		kv := strings.SplitAfterN(e, "=", 2)
		env[kv[0]] = kv[1]
	***REMOVED***

	return env
***REMOVED***

type roundTripFunc func(*http.Request) (*http.Response, error)

func (rtf roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	return rtf(req)
***REMOVED***

type bytesBufferClose struct ***REMOVED***
	*bytes.Buffer
***REMOVED***

func (bbc bytesBufferClose) Close() error ***REMOVED***
	return nil
***REMOVED***

func TestClientRedirect(t *testing.T) ***REMOVED***
	client := &http.Client***REMOVED***
		CheckRedirect: CheckRedirect,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if req.URL.String() == "/bla" ***REMOVED***
				return &http.Response***REMOVED***StatusCode: 404***REMOVED***, nil
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: 301,
				Header:     map[string][]string***REMOVED***"Location": ***REMOVED***"/bla"***REMOVED******REMOVED***,
				Body:       bytesBufferClose***REMOVED***bytes.NewBuffer(nil)***REMOVED***,
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	cases := []struct ***REMOVED***
		httpMethod  string
		expectedErr error
		statusCode  int
	***REMOVED******REMOVED***
		***REMOVED***http.MethodGet, nil, 301***REMOVED***,
		***REMOVED***http.MethodPost, &url.Error***REMOVED***Op: "Post", URL: "/bla", Err: ErrRedirect***REMOVED***, 301***REMOVED***,
		***REMOVED***http.MethodPut, &url.Error***REMOVED***Op: "Put", URL: "/bla", Err: ErrRedirect***REMOVED***, 301***REMOVED***,
		***REMOVED***http.MethodDelete, &url.Error***REMOVED***Op: "Delete", URL: "/bla", Err: ErrRedirect***REMOVED***, 301***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		req, err := http.NewRequest(tc.httpMethod, "/redirectme", nil)
		assert.NoError(t, err)
		resp, err := client.Do(req)
		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.statusCode, resp.StatusCode)
	***REMOVED***
***REMOVED***
