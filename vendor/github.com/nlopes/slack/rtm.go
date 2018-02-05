package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const (
	websocketDefaultTimeout = 10 * time.Second
)

// StartRTM calls the "rtm.start" endpoint and returns the provided URL and the full Info block.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) StartRTM() (info *Info, websocketURL string, err error) ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), websocketDefaultTimeout)
	defer cancel()

	return api.StartRTMContext(ctx)
***REMOVED***

// StartRTMContext calls the "rtm.start" endpoint and returns the provided URL and the full Info block with a custom context.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) StartRTMContext(ctx context.Context) (info *Info, websocketURL string, err error) ***REMOVED***
	response := &infoResponseFull***REMOVED******REMOVED***
	err = post(ctx, api.httpclient, "rtm.start", url.Values***REMOVED***"token": ***REMOVED***api.token***REMOVED******REMOVED***, response, api.debug)
	if err != nil ***REMOVED***
		return nil, "", fmt.Errorf("post: %s", err)
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, "", response.Error
	***REMOVED***
	api.Debugln("Using URL:", response.Info.URL)
	return &response.Info, response.Info.URL, nil
***REMOVED***

// ConnectRTM calls the "rtm.connect" endpoint and returns the provided URL and the compact Info block.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) ConnectRTM() (info *Info, websocketURL string, err error) ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), websocketDefaultTimeout)
	defer cancel()

	return api.ConnectRTMContext(ctx)
***REMOVED***

// ConnectRTM calls the "rtm.connect" endpoint and returns the provided URL and the compact Info block with a custom context.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) ConnectRTMContext(ctx context.Context) (info *Info, websocketURL string, err error) ***REMOVED***
	response := &infoResponseFull***REMOVED******REMOVED***
	err = post(ctx, api.httpclient, "rtm.connect", url.Values***REMOVED***"token": ***REMOVED***api.token***REMOVED******REMOVED***, response, api.debug)
	if err != nil ***REMOVED***
		api.Debugf("Failed to connect to RTM: %s", err)
		return nil, "", fmt.Errorf("post: %s", err)
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, "", response.Error
	***REMOVED***
	api.Debugln("Using URL:", response.Info.URL)
	return &response.Info, response.Info.URL, nil
***REMOVED***

// NewRTM returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
func (api *Client) NewRTM() *RTM ***REMOVED***
	return api.NewRTMWithOptions(nil)
***REMOVED***

// NewRTMWithOptions returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
// This also allows to configure various options available for RTM API.
func (api *Client) NewRTMWithOptions(options *RTMOptions) *RTM ***REMOVED***
	result := &RTM***REMOVED***
		Client:           *api,
		IncomingEvents:   make(chan RTMEvent, 50),
		outgoingMessages: make(chan OutgoingMessage, 20),
		pings:            make(map[int]time.Time),
		isConnected:      false,
		wasIntentional:   true,
		killChannel:      make(chan bool),
		disconnected:     make(chan struct***REMOVED******REMOVED***),
		forcePing:        make(chan bool),
		rawEvents:        make(chan json.RawMessage),
		idGen:            NewSafeID(1),
	***REMOVED***

	if options != nil ***REMOVED***
		result.useRTMStart = options.UseRTMStart
	***REMOVED*** else ***REMOVED***
		result.useRTMStart = true
	***REMOVED***

	return result
***REMOVED***
