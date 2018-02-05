package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type WebResponse struct ***REMOVED***
	Ok    bool      `json:"ok"`
	Error *WebError `json:"error"`
***REMOVED***

type WebError string

func (s WebError) Error() string ***REMOVED***
	return string(s)
***REMOVED***

type RateLimitedError struct ***REMOVED***
	RetryAfter time.Duration
***REMOVED***

func (e *RateLimitedError) Error() string ***REMOVED***
	return fmt.Sprintf("Slack rate limit exceeded, retry after %s", e.RetryAfter)
***REMOVED***

func fileUploadReq(ctx context.Context, path, fieldname, filename string, values url.Values, r io.Reader) (*http.Request, error) ***REMOVED***
	body := &bytes.Buffer***REMOVED******REMOVED***
	wr := multipart.NewWriter(body)

	ioWriter, err := wr.CreateFormFile(fieldname, filename)
	if err != nil ***REMOVED***
		wr.Close()
		return nil, err
	***REMOVED***
	_, err = io.Copy(ioWriter, r)
	if err != nil ***REMOVED***
		wr.Close()
		return nil, err
	***REMOVED***
	// Close the multipart writer or the footer won't be written
	wr.Close()
	req, err := http.NewRequest("POST", path, body)
	req = req.WithContext(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Add("Content-Type", wr.FormDataContentType())
	req.URL.RawQuery = (values).Encode()
	return req, nil
***REMOVED***

func parseResponseBody(body io.ReadCloser, intf *interface***REMOVED******REMOVED***, debug bool) error ***REMOVED***
	response, err := ioutil.ReadAll(body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// FIXME: will be api.Debugf
	if debug ***REMOVED***
		logger.Printf("parseResponseBody: %s\n", string(response))
	***REMOVED***

	return json.Unmarshal(response, &intf)
***REMOVED***

func postLocalWithMultipartResponse(ctx context.Context, client HTTPRequester, path, fpath, fieldname string, values url.Values, intf interface***REMOVED******REMOVED***, debug bool) error ***REMOVED***
	fullpath, err := filepath.Abs(fpath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	file, err := os.Open(fullpath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer file.Close()
	return postWithMultipartResponse(ctx, client, path, filepath.Base(fpath), fieldname, values, file, intf, debug)
***REMOVED***

func postWithMultipartResponse(ctx context.Context, client HTTPRequester, path, name, fieldname string, values url.Values, r io.Reader, intf interface***REMOVED******REMOVED***, debug bool) error ***REMOVED***
	req, err := fileUploadReq(ctx, SLACK_API+path, fieldname, name, values, r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests ***REMOVED***
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return &RateLimitedError***REMOVED***time.Duration(retry) * time.Second***REMOVED***
	***REMOVED***

	// Slack seems to send an HTML body along with 5xx error codes. Don't parse it.
	if resp.StatusCode != http.StatusOK ***REMOVED***
		logResponse(resp, debug)
		return fmt.Errorf("Slack server error: %s.", resp.Status)
	***REMOVED***

	return parseResponseBody(resp.Body, &intf, debug)
***REMOVED***

func postForm(ctx context.Context, client HTTPRequester, endpoint string, values url.Values, intf interface***REMOVED******REMOVED***, debug bool) error ***REMOVED***
	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests ***REMOVED***
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return &RateLimitedError***REMOVED***time.Duration(retry) * time.Second***REMOVED***
	***REMOVED***

	// Slack seems to send an HTML body along with 5xx error codes. Don't parse it.
	if resp.StatusCode != http.StatusOK ***REMOVED***
		logResponse(resp, debug)
		return fmt.Errorf("Slack server error: %s.", resp.Status)
	***REMOVED***

	return parseResponseBody(resp.Body, &intf, debug)
***REMOVED***

func post(ctx context.Context, client HTTPRequester, path string, values url.Values, intf interface***REMOVED******REMOVED***, debug bool) error ***REMOVED***
	return postForm(ctx, client, SLACK_API+path, values, intf, debug)
***REMOVED***

func parseAdminResponse(ctx context.Context, client HTTPRequester, method string, teamName string, values url.Values, intf interface***REMOVED******REMOVED***, debug bool) error ***REMOVED***
	endpoint := fmt.Sprintf(SLACK_WEB_API_FORMAT, teamName, method, time.Now().Unix())
	return postForm(ctx, client, endpoint, values, intf, debug)
***REMOVED***

func logResponse(resp *http.Response, debug bool) error ***REMOVED***
	if debug ***REMOVED***
		text, err := httputil.DumpResponse(resp, true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		logger.Print(string(text))
	***REMOVED***

	return nil
***REMOVED***

func okJsonHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(SlackResponse***REMOVED***
		Ok: true,
	***REMOVED***)
	rw.Write(response)
***REMOVED***
