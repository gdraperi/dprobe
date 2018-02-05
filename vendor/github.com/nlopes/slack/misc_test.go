package slack

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

var (
	parseResponseOnce sync.Once
)

func parseResponseHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	token := r.FormValue("token")
	log.Println(token)
	if token == "" ***REMOVED***
		rw.Write([]byte(`***REMOVED***"ok":false,"error":"not_authed"***REMOVED***`))
		return
	***REMOVED***
	if token != validToken ***REMOVED***
		rw.Write([]byte(`***REMOVED***"ok":false,"error":"invalid_auth"***REMOVED***`))
		return
	***REMOVED***
	response := []byte(`***REMOVED***"ok": true***REMOVED***`)
	rw.Write(response)
***REMOVED***

func setParseResponseHandler() ***REMOVED***
	http.HandleFunc("/parseResponse", parseResponseHandler)
***REMOVED***

func TestParseResponse(t *testing.T) ***REMOVED***
	parseResponseOnce.Do(setParseResponseHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	values := url.Values***REMOVED***
		"token": ***REMOVED***validToken***REMOVED***,
	***REMOVED***
	responsePartial := &SlackResponse***REMOVED******REMOVED***
	err := post(context.Background(), http.DefaultClient, "parseResponse", values, responsePartial, false)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
	***REMOVED***
***REMOVED***

func TestParseResponseNoToken(t *testing.T) ***REMOVED***
	parseResponseOnce.Do(setParseResponseHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	values := url.Values***REMOVED******REMOVED***
	responsePartial := &SlackResponse***REMOVED******REMOVED***
	err := post(context.Background(), http.DefaultClient, "parseResponse", values, responsePartial, false)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if responsePartial.Ok == true ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
	***REMOVED*** else if responsePartial.Error != "not_authed" ***REMOVED***
		t.Errorf("got %v; want %v", responsePartial.Error, "not_authed")
	***REMOVED***
***REMOVED***

func TestParseResponseInvalidToken(t *testing.T) ***REMOVED***
	parseResponseOnce.Do(setParseResponseHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	values := url.Values***REMOVED***
		"token": ***REMOVED***"whatever"***REMOVED***,
	***REMOVED***
	responsePartial := &SlackResponse***REMOVED******REMOVED***
	err := post(context.Background(), http.DefaultClient, "parseResponse", values, responsePartial, false)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if responsePartial.Ok == true ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
	***REMOVED*** else if responsePartial.Error != "invalid_auth" ***REMOVED***
		t.Errorf("got %v; want %v", responsePartial.Error, "invalid_auth")
	***REMOVED***
***REMOVED***
