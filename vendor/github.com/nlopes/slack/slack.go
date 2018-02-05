package slack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Added as a var so that we can change this for testing purposes
var SLACK_API string = "https://slack.com/api/"
var SLACK_WEB_API_FORMAT string = "https://%s.slack.com/api/users.admin.%s?t=%s"

// HTTPClient sets a custom http.Client
// deprecated: in favor of SetHTTPClient()
var HTTPClient = &http.Client***REMOVED******REMOVED***

var customHTTPClient HTTPRequester = HTTPClient

// HTTPRequester defines the minimal interface needed for an http.Client to be implemented.
//
// Use it in conjunction with the SetHTTPClient function to allow for other capabilities
// like a tracing http.Client
type HTTPRequester interface ***REMOVED***
	Do(*http.Request) (*http.Response, error)
***REMOVED***

// SetHTTPClient allows you to specify a custom http.Client
// Use this instead of the package level HTTPClient variable if you want to use a custom client like the
// Stackdriver Trace HTTPClient https://godoc.org/cloud.google.com/go/trace#HTTPClient
func SetHTTPClient(client HTTPRequester) ***REMOVED***
	customHTTPClient = client
***REMOVED***

type SlackResponse struct ***REMOVED***
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
***REMOVED***

type AuthTestResponse struct ***REMOVED***
	URL    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
***REMOVED***

type authTestResponseFull struct ***REMOVED***
	SlackResponse
	AuthTestResponse
***REMOVED***

type Client struct ***REMOVED***
	token      string
	info       Info
	debug      bool
	httpclient HTTPRequester
***REMOVED***

// Option defines an option for a Client
type Option func(*Client)

// OptionHTTPClient - provide a custom http client to the slack client.
func OptionHTTPClient(c HTTPRequester) func(*Client) ***REMOVED***
	return func(s *Client) ***REMOVED***
		s.httpclient = c
	***REMOVED***
***REMOVED***

// New builds a slack client from the provided token and options.
func New(token string, options ...Option) *Client ***REMOVED***
	s := &Client***REMOVED***
		token:      token,
		httpclient: customHTTPClient,
	***REMOVED***

	for _, opt := range options ***REMOVED***
		opt(s)
	***REMOVED***

	return s
***REMOVED***

// AuthTest tests if the user is able to do authenticated requests or not
func (api *Client) AuthTest() (response *AuthTestResponse, error error) ***REMOVED***
	return api.AuthTestContext(context.Background())
***REMOVED***

// AuthTestContext tests if the user is able to do authenticated requests or not with a custom context
func (api *Client) AuthTestContext(ctx context.Context) (response *AuthTestResponse, error error) ***REMOVED***
	api.Debugf("Challenging auth...")
	responseFull := &authTestResponseFull***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "auth.test", url.Values***REMOVED***"token": ***REMOVED***api.token***REMOVED******REMOVED***, responseFull, api.debug)
	if err != nil ***REMOVED***
		api.Debugf("failed to test for auth: %s", err)
		return nil, err
	***REMOVED***
	if !responseFull.Ok ***REMOVED***
		api.Debugf("auth response was not Ok: %s", responseFull.Error)
		return nil, errors.New(responseFull.Error)
	***REMOVED***

	api.Debugf("Auth challenge was successful with response %+v", responseFull.AuthTestResponse)
	return &responseFull.AuthTestResponse, nil
***REMOVED***

// SetDebug switches the api into debug mode
// When in debug mode, it logs various info about what its doing
// If you ever use this in production, don't call SetDebug(true)
func (api *Client) SetDebug(debug bool) ***REMOVED***
	api.debug = debug
	if debug && logger == nil ***REMOVED***
		SetLogger(log.New(os.Stdout, "nlopes/slack", log.LstdFlags|log.Lshortfile))
	***REMOVED***
***REMOVED***

// Debugf print a formatted debug line.
func (api *Client) Debugf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	if api.debug ***REMOVED***
		logger.Output(2, fmt.Sprintf(format, v...))
	***REMOVED***
***REMOVED***

// Debugln print a debug line.
func (api *Client) Debugln(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	if api.debug ***REMOVED***
		logger.Output(2, fmt.Sprintln(v...))
	***REMOVED***
***REMOVED***
