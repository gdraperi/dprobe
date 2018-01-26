package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// serverResponse is a wrapper for http API responses.
type serverResponse struct ***REMOVED***
	body       io.ReadCloser
	header     http.Header
	statusCode int
	reqURL     *url.URL
***REMOVED***

// head sends an http request to the docker API using the method HEAD.
func (cli *Client) head(ctx context.Context, path string, query url.Values, headers map[string][]string) (serverResponse, error) ***REMOVED***
	return cli.sendRequest(ctx, "HEAD", path, query, nil, headers)
***REMOVED***

// get sends an http request to the docker API using the method GET with a specific Go context.
func (cli *Client) get(ctx context.Context, path string, query url.Values, headers map[string][]string) (serverResponse, error) ***REMOVED***
	return cli.sendRequest(ctx, "GET", path, query, nil, headers)
***REMOVED***

// post sends an http request to the docker API using the method POST with a specific Go context.
func (cli *Client) post(ctx context.Context, path string, query url.Values, obj interface***REMOVED******REMOVED***, headers map[string][]string) (serverResponse, error) ***REMOVED***
	body, headers, err := encodeBody(obj, headers)
	if err != nil ***REMOVED***
		return serverResponse***REMOVED******REMOVED***, err
	***REMOVED***
	return cli.sendRequest(ctx, "POST", path, query, body, headers)
***REMOVED***

func (cli *Client) postRaw(ctx context.Context, path string, query url.Values, body io.Reader, headers map[string][]string) (serverResponse, error) ***REMOVED***
	return cli.sendRequest(ctx, "POST", path, query, body, headers)
***REMOVED***

// put sends an http request to the docker API using the method PUT.
func (cli *Client) put(ctx context.Context, path string, query url.Values, obj interface***REMOVED******REMOVED***, headers map[string][]string) (serverResponse, error) ***REMOVED***
	body, headers, err := encodeBody(obj, headers)
	if err != nil ***REMOVED***
		return serverResponse***REMOVED******REMOVED***, err
	***REMOVED***
	return cli.sendRequest(ctx, "PUT", path, query, body, headers)
***REMOVED***

// putRaw sends an http request to the docker API using the method PUT.
func (cli *Client) putRaw(ctx context.Context, path string, query url.Values, body io.Reader, headers map[string][]string) (serverResponse, error) ***REMOVED***
	return cli.sendRequest(ctx, "PUT", path, query, body, headers)
***REMOVED***

// delete sends an http request to the docker API using the method DELETE.
func (cli *Client) delete(ctx context.Context, path string, query url.Values, headers map[string][]string) (serverResponse, error) ***REMOVED***
	return cli.sendRequest(ctx, "DELETE", path, query, nil, headers)
***REMOVED***

type headers map[string][]string

func encodeBody(obj interface***REMOVED******REMOVED***, headers headers) (io.Reader, headers, error) ***REMOVED***
	if obj == nil ***REMOVED***
		return nil, headers, nil
	***REMOVED***

	body, err := encodeData(obj)
	if err != nil ***REMOVED***
		return nil, headers, err
	***REMOVED***
	if headers == nil ***REMOVED***
		headers = make(map[string][]string)
	***REMOVED***
	headers["Content-Type"] = []string***REMOVED***"application/json"***REMOVED***
	return body, headers, nil
***REMOVED***

func (cli *Client) buildRequest(method, path string, body io.Reader, headers headers) (*http.Request, error) ***REMOVED***
	expectedPayload := (method == "POST" || method == "PUT")
	if expectedPayload && body == nil ***REMOVED***
		body = bytes.NewReader([]byte***REMOVED******REMOVED***)
	***REMOVED***

	req, err := http.NewRequest(method, path, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req = cli.addHeaders(req, headers)

	if cli.proto == "unix" || cli.proto == "npipe" ***REMOVED***
		// For local communications, it doesn't matter what the host is. We just
		// need a valid and meaningful host name. (See #189)
		req.Host = "docker"
	***REMOVED***

	req.URL.Host = cli.addr
	req.URL.Scheme = cli.scheme

	if expectedPayload && req.Header.Get("Content-Type") == "" ***REMOVED***
		req.Header.Set("Content-Type", "text/plain")
	***REMOVED***
	return req, nil
***REMOVED***

func (cli *Client) sendRequest(ctx context.Context, method, path string, query url.Values, body io.Reader, headers headers) (serverResponse, error) ***REMOVED***
	req, err := cli.buildRequest(method, cli.getAPIPath(path, query), body, headers)
	if err != nil ***REMOVED***
		return serverResponse***REMOVED******REMOVED***, err
	***REMOVED***
	resp, err := cli.doRequest(ctx, req)
	if err != nil ***REMOVED***
		return resp, err
	***REMOVED***
	if err := cli.checkResponseErr(resp); err != nil ***REMOVED***
		return resp, err
	***REMOVED***
	return resp, nil
***REMOVED***

func (cli *Client) doRequest(ctx context.Context, req *http.Request) (serverResponse, error) ***REMOVED***
	serverResp := serverResponse***REMOVED***statusCode: -1, reqURL: req.URL***REMOVED***

	resp, err := ctxhttp.Do(ctx, cli.client, req)
	if err != nil ***REMOVED***
		if cli.scheme != "https" && strings.Contains(err.Error(), "malformed HTTP response") ***REMOVED***
			return serverResp, fmt.Errorf("%v.\n* Are you trying to connect to a TLS-enabled daemon without TLS?", err)
		***REMOVED***

		if cli.scheme == "https" && strings.Contains(err.Error(), "bad certificate") ***REMOVED***
			return serverResp, fmt.Errorf("The server probably has client authentication (--tlsverify) enabled. Please check your TLS client certification settings: %v", err)
		***REMOVED***

		// Don't decorate context sentinel errors; users may be comparing to
		// them directly.
		switch err ***REMOVED***
		case context.Canceled, context.DeadlineExceeded:
			return serverResp, err
		***REMOVED***

		if nErr, ok := err.(*url.Error); ok ***REMOVED***
			if nErr, ok := nErr.Err.(*net.OpError); ok ***REMOVED***
				if os.IsPermission(nErr.Err) ***REMOVED***
					return serverResp, errors.Wrapf(err, "Got permission denied while trying to connect to the Docker daemon socket at %v", cli.host)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err, ok := err.(net.Error); ok ***REMOVED***
			if err.Timeout() ***REMOVED***
				return serverResp, ErrorConnectionFailed(cli.host)
			***REMOVED***
			if !err.Temporary() ***REMOVED***
				if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "dial unix") ***REMOVED***
					return serverResp, ErrorConnectionFailed(cli.host)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Although there's not a strongly typed error for this in go-winio,
		// lots of people are using the default configuration for the docker
		// daemon on Windows where the daemon is listening on a named pipe
		// `//./pipe/docker_engine, and the client must be running elevated.
		// Give users a clue rather than the not-overly useful message
		// such as `error during connect: Get http://%2F%2F.%2Fpipe%2Fdocker_engine/v1.26/info:
		// open //./pipe/docker_engine: The system cannot find the file specified.`.
		// Note we can't string compare "The system cannot find the file specified" as
		// this is localised - for example in French the error would be
		// `open //./pipe/docker_engine: Le fichier spécifié est introuvable.`
		if strings.Contains(err.Error(), `open //./pipe/docker_engine`) ***REMOVED***
			err = errors.New(err.Error() + " In the default daemon configuration on Windows, the docker client must be run elevated to connect. This error may also indicate that the docker daemon is not running.")
		***REMOVED***

		return serverResp, errors.Wrap(err, "error during connect")
	***REMOVED***

	if resp != nil ***REMOVED***
		serverResp.statusCode = resp.StatusCode
		serverResp.body = resp.Body
		serverResp.header = resp.Header
	***REMOVED***
	return serverResp, nil
***REMOVED***

func (cli *Client) checkResponseErr(serverResp serverResponse) error ***REMOVED***
	if serverResp.statusCode >= 200 && serverResp.statusCode < 400 ***REMOVED***
		return nil
	***REMOVED***

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(body) == 0 ***REMOVED***
		return fmt.Errorf("request returned %s for API route and version %s, check if the server supports the requested API version", http.StatusText(serverResp.statusCode), serverResp.reqURL)
	***REMOVED***

	var ct string
	if serverResp.header != nil ***REMOVED***
		ct = serverResp.header.Get("Content-Type")
	***REMOVED***

	var errorMessage string
	if (cli.version == "" || versions.GreaterThan(cli.version, "1.23")) && ct == "application/json" ***REMOVED***
		var errorResponse types.ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil ***REMOVED***
			return fmt.Errorf("Error reading JSON: %v", err)
		***REMOVED***
		errorMessage = errorResponse.Message
	***REMOVED*** else ***REMOVED***
		errorMessage = string(body)
	***REMOVED***

	return fmt.Errorf("Error response from daemon: %s", strings.TrimSpace(errorMessage))
***REMOVED***

func (cli *Client) addHeaders(req *http.Request, headers headers) *http.Request ***REMOVED***
	// Add CLI Config's HTTP Headers BEFORE we set the Docker headers
	// then the user can't change OUR headers
	for k, v := range cli.customHTTPHeaders ***REMOVED***
		if versions.LessThan(cli.version, "1.25") && k == "User-Agent" ***REMOVED***
			continue
		***REMOVED***
		req.Header.Set(k, v)
	***REMOVED***

	if headers != nil ***REMOVED***
		for k, v := range headers ***REMOVED***
			req.Header[k] = v
		***REMOVED***
	***REMOVED***
	return req
***REMOVED***

func encodeData(data interface***REMOVED******REMOVED***) (*bytes.Buffer, error) ***REMOVED***
	params := bytes.NewBuffer(nil)
	if data != nil ***REMOVED***
		if err := json.NewEncoder(params).Encode(data); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return params, nil
***REMOVED***

func ensureReaderClosed(response serverResponse) ***REMOVED***
	if response.body != nil ***REMOVED***
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, response.body, 512)
		response.body.Close()
	***REMOVED***
***REMOVED***
