package remotecontext

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/internal/testutil"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var binaryContext = []byte***REMOVED***0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00***REMOVED*** //xz magic

func TestSelectAcceptableMIME(t *testing.T) ***REMOVED***
	validMimeStrings := []string***REMOVED***
		"application/x-bzip2",
		"application/bzip2",
		"application/gzip",
		"application/x-gzip",
		"application/x-xz",
		"application/xz",
		"application/tar",
		"application/x-tar",
		"application/octet-stream",
		"text/plain",
	***REMOVED***

	invalidMimeStrings := []string***REMOVED***
		"",
		"application/octet",
		"application/json",
	***REMOVED***

	for _, m := range invalidMimeStrings ***REMOVED***
		if len(selectAcceptableMIME(m)) > 0 ***REMOVED***
			t.Fatalf("Should not have accepted %q", m)
		***REMOVED***
	***REMOVED***

	for _, m := range validMimeStrings ***REMOVED***
		if str := selectAcceptableMIME(m); str == "" ***REMOVED***
			t.Fatalf("Should have accepted %q", m)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestInspectEmptyResponse(t *testing.T) ***REMOVED***
	ct := "application/octet-stream"
	br := ioutil.NopCloser(bytes.NewReader([]byte("")))
	contentType, bReader, err := inspectResponse(ct, br, 0)
	if err == nil ***REMOVED***
		t.Fatal("Should have generated an error for an empty response")
	***REMOVED***
	if contentType != "application/octet-stream" ***REMOVED***
		t.Fatalf("Content type should be 'application/octet-stream' but is %q", contentType)
	***REMOVED***
	body, err := ioutil.ReadAll(bReader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(body) != 0 ***REMOVED***
		t.Fatal("response body should remain empty")
	***REMOVED***
***REMOVED***

func TestInspectResponseBinary(t *testing.T) ***REMOVED***
	ct := "application/octet-stream"
	br := ioutil.NopCloser(bytes.NewReader(binaryContext))
	contentType, bReader, err := inspectResponse(ct, br, int64(len(binaryContext)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if contentType != "application/octet-stream" ***REMOVED***
		t.Fatalf("Content type should be 'application/octet-stream' but is %q", contentType)
	***REMOVED***
	body, err := ioutil.ReadAll(bReader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(body) != len(binaryContext) ***REMOVED***
		t.Fatalf("Wrong response size %d, should be == len(binaryContext)", len(body))
	***REMOVED***
	for i := range body ***REMOVED***
		if body[i] != binaryContext[i] ***REMOVED***
			t.Fatalf("Corrupted response body at byte index %d", i)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestResponseUnsupportedContentType(t *testing.T) ***REMOVED***
	content := []byte(dockerfileContents)
	ct := "application/json"
	br := ioutil.NopCloser(bytes.NewReader(content))
	contentType, bReader, err := inspectResponse(ct, br, int64(len(dockerfileContents)))

	if err == nil ***REMOVED***
		t.Fatal("Should have returned an error on content-type 'application/json'")
	***REMOVED***
	if contentType != ct ***REMOVED***
		t.Fatalf("Should not have altered content-type: orig: %s, altered: %s", ct, contentType)
	***REMOVED***
	body, err := ioutil.ReadAll(bReader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(body) != dockerfileContents ***REMOVED***
		t.Fatalf("Corrupted response body %s", body)
	***REMOVED***
***REMOVED***

func TestInspectResponseTextSimple(t *testing.T) ***REMOVED***
	content := []byte(dockerfileContents)
	ct := "text/plain"
	br := ioutil.NopCloser(bytes.NewReader(content))
	contentType, bReader, err := inspectResponse(ct, br, int64(len(content)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if contentType != "text/plain" ***REMOVED***
		t.Fatalf("Content type should be 'text/plain' but is %q", contentType)
	***REMOVED***
	body, err := ioutil.ReadAll(bReader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(body) != dockerfileContents ***REMOVED***
		t.Fatalf("Corrupted response body %s", body)
	***REMOVED***
***REMOVED***

func TestInspectResponseEmptyContentType(t *testing.T) ***REMOVED***
	content := []byte(dockerfileContents)
	br := ioutil.NopCloser(bytes.NewReader(content))
	contentType, bodyReader, err := inspectResponse("", br, int64(len(content)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if contentType != "text/plain" ***REMOVED***
		t.Fatalf("Content type should be 'text/plain' but is %q", contentType)
	***REMOVED***
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(body) != dockerfileContents ***REMOVED***
		t.Fatalf("Corrupted response body %s", body)
	***REMOVED***
***REMOVED***

func TestUnknownContentLength(t *testing.T) ***REMOVED***
	content := []byte(dockerfileContents)
	ct := "text/plain"
	br := ioutil.NopCloser(bytes.NewReader(content))
	contentType, bReader, err := inspectResponse(ct, br, -1)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if contentType != "text/plain" ***REMOVED***
		t.Fatalf("Content type should be 'text/plain' but is %q", contentType)
	***REMOVED***
	body, err := ioutil.ReadAll(bReader)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(body) != dockerfileContents ***REMOVED***
		t.Fatalf("Corrupted response body %s", body)
	***REMOVED***
***REMOVED***

func TestDownloadRemote(t *testing.T) ***REMOVED***
	contextDir := fs.NewDir(t, "test-builder-download-remote",
		fs.WithFile(builder.DefaultDockerfileName, dockerfileContents))
	defer contextDir.Remove()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	serverURL, _ := url.Parse(server.URL)

	serverURL.Path = "/" + builder.DefaultDockerfileName
	remoteURL := serverURL.String()

	mux.Handle("/", http.FileServer(http.Dir(contextDir.Path())))

	contentType, content, err := downloadRemote(remoteURL)
	require.NoError(t, err)

	assert.Equal(t, mimeTypes.TextPlain, contentType)
	raw, err := ioutil.ReadAll(content)
	require.NoError(t, err)
	assert.Equal(t, dockerfileContents, string(raw))
***REMOVED***

func TestGetWithStatusError(t *testing.T) ***REMOVED***
	var testcases = []struct ***REMOVED***
		err          error
		statusCode   int
		expectedErr  string
		expectedBody string
	***REMOVED******REMOVED***
		***REMOVED***
			statusCode:   200,
			expectedBody: "THE BODY",
		***REMOVED***,
		***REMOVED***
			statusCode:   400,
			expectedErr:  "with status 400 Bad Request: broke",
			expectedBody: "broke",
		***REMOVED***,
	***REMOVED***
	for _, testcase := range testcases ***REMOVED***
		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
				buffer := bytes.NewBufferString(testcase.expectedBody)
				w.WriteHeader(testcase.statusCode)
				w.Write(buffer.Bytes())
			***REMOVED***),
		)
		defer ts.Close()
		response, err := GetWithStatusError(ts.URL)

		if testcase.expectedErr == "" ***REMOVED***
			require.NoError(t, err)

			body, err := readBody(response.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), testcase.expectedBody)
		***REMOVED*** else ***REMOVED***
			testutil.ErrorContains(t, err, testcase.expectedErr)
		***REMOVED***
	***REMOVED***
***REMOVED***

func readBody(b io.ReadCloser) ([]byte, error) ***REMOVED***
	defer b.Close()
	return ioutil.ReadAll(b)
***REMOVED***
