package assert

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// httpCode is a helper that returns HTTP code of the response. It returns -1
// if building a new request fails.
func httpCode(handler http.HandlerFunc, method, url string, values url.Values) int ***REMOVED***
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url+"?"+values.Encode(), nil)
	if err != nil ***REMOVED***
		return -1
	***REMOVED***
	handler(w, req)
	return w.Code
***REMOVED***

// HTTPSuccess asserts that a specified handler returns a success status code.
//
//  assert.HTTPSuccess(t, myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPSuccess(t TestingT, handler http.HandlerFunc, method, url string, values url.Values) bool ***REMOVED***
	code := httpCode(handler, method, url, values)
	if code == -1 ***REMOVED***
		return false
	***REMOVED***
	return code >= http.StatusOK && code <= http.StatusPartialContent
***REMOVED***

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
//  assert.HTTPRedirect(t, myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPRedirect(t TestingT, handler http.HandlerFunc, method, url string, values url.Values) bool ***REMOVED***
	code := httpCode(handler, method, url, values)
	if code == -1 ***REMOVED***
		return false
	***REMOVED***
	return code >= http.StatusMultipleChoices && code <= http.StatusTemporaryRedirect
***REMOVED***

// HTTPError asserts that a specified handler returns an error status code.
//
//  assert.HTTPError(t, myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPError(t TestingT, handler http.HandlerFunc, method, url string, values url.Values) bool ***REMOVED***
	code := httpCode(handler, method, url, values)
	if code == -1 ***REMOVED***
		return false
	***REMOVED***
	return code >= http.StatusBadRequest
***REMOVED***

// HTTPBody is a helper that returns HTTP body of the response. It returns
// empty string if building a new request fails.
func HTTPBody(handler http.HandlerFunc, method, url string, values url.Values) string ***REMOVED***
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url+"?"+values.Encode(), nil)
	if err != nil ***REMOVED***
		return ""
	***REMOVED***
	handler(w, req)
	return w.Body.String()
***REMOVED***

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
//  assert.HTTPBodyContains(t, myHandler, "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyContains(t TestingT, handler http.HandlerFunc, method, url string, values url.Values, str interface***REMOVED******REMOVED***) bool ***REMOVED***
	body := HTTPBody(handler, method, url, values)

	contains := strings.Contains(body, fmt.Sprint(str))
	if !contains ***REMOVED***
		Fail(t, fmt.Sprintf("Expected response body for \"%s\" to contain \"%s\" but found \"%s\"", url+"?"+values.Encode(), str, body))
	***REMOVED***

	return contains
***REMOVED***

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
//  assert.HTTPBodyNotContains(t, myHandler, "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyNotContains(t TestingT, handler http.HandlerFunc, method, url string, values url.Values, str interface***REMOVED******REMOVED***) bool ***REMOVED***
	body := HTTPBody(handler, method, url, values)

	contains := strings.Contains(body, fmt.Sprint(str))
	if contains ***REMOVED***
		Fail(t, fmt.Sprintf("Expected response body for \"%s\" to NOT contain \"%s\" but found \"%s\"", url+"?"+values.Encode(), str, body))
	***REMOVED***

	return !contains
***REMOVED***
