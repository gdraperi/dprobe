package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

type imChannel struct ***REMOVED***
	ID string `json:"id"`
***REMOVED***

type imResponseFull struct ***REMOVED***
	NoOp          bool      `json:"no_op"`
	AlreadyClosed bool      `json:"already_closed"`
	AlreadyOpen   bool      `json:"already_open"`
	Channel       imChannel `json:"channel"`
	IMs           []IM      `json:"ims"`
	History
	SlackResponse
***REMOVED***

// IM contains information related to the Direct Message channel
type IM struct ***REMOVED***
	conversation
	IsIM          bool   `json:"is_im"`
	User          string `json:"user"`
	IsUserDeleted bool   `json:"is_user_deleted"`
***REMOVED***

func imRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*imResponseFull, error) ***REMOVED***
	response := &imResponseFull***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// CloseIMChannel closes the direct message channel
func (api *Client) CloseIMChannel(channel string) (bool, bool, error) ***REMOVED***
	return api.CloseIMChannelContext(context.Background(), channel)
***REMOVED***

// CloseIMChannelContext closes the direct message channel with a custom context
func (api *Client) CloseIMChannelContext(ctx context.Context, channel string) (bool, bool, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channel***REMOVED***,
	***REMOVED***

	response, err := imRequest(ctx, api.httpclient, "im.close", values, api.debug)
	if err != nil ***REMOVED***
		return false, false, err
	***REMOVED***
	return response.NoOp, response.AlreadyClosed, nil
***REMOVED***

// OpenIMChannel opens a direct message channel to the user provided as argument
// Returns some status and the channel ID
func (api *Client) OpenIMChannel(user string) (bool, bool, string, error) ***REMOVED***
	return api.OpenIMChannelContext(context.Background(), user)
***REMOVED***

// OpenIMChannelContext opens a direct message channel to the user provided as argument with a custom context
// Returns some status and the channel ID
func (api *Client) OpenIMChannelContext(ctx context.Context, user string) (bool, bool, string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"user":  ***REMOVED***user***REMOVED***,
	***REMOVED***

	response, err := imRequest(ctx, api.httpclient, "im.open", values, api.debug)
	if err != nil ***REMOVED***
		return false, false, "", err
	***REMOVED***
	return response.NoOp, response.AlreadyOpen, response.Channel.ID, nil
***REMOVED***

// MarkIMChannel sets the read mark of a direct message channel to a specific point
func (api *Client) MarkIMChannel(channel, ts string) (err error) ***REMOVED***
	return api.MarkIMChannelContext(context.Background(), channel, ts)
***REMOVED***

// MarkIMChannelContext sets the read mark of a direct message channel to a specific point with a custom context
func (api *Client) MarkIMChannelContext(ctx context.Context, channel, ts string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channel***REMOVED***,
		"ts":      ***REMOVED***ts***REMOVED***,
	***REMOVED***

	_, err = imRequest(ctx, api.httpclient, "im.mark", values, api.debug)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return
***REMOVED***

// GetIMHistory retrieves the direct message channel history
func (api *Client) GetIMHistory(channel string, params HistoryParameters) (*History, error) ***REMOVED***
	return api.GetIMHistoryContext(context.Background(), channel, params)
***REMOVED***

// GetIMHistoryContext retrieves the direct message channel history with a custom context
func (api *Client) GetIMHistoryContext(ctx context.Context, channel string, params HistoryParameters) (*History, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channel***REMOVED***,
	***REMOVED***
	if params.Latest != DEFAULT_HISTORY_LATEST ***REMOVED***
		values.Add("latest", params.Latest)
	***REMOVED***
	if params.Oldest != DEFAULT_HISTORY_OLDEST ***REMOVED***
		values.Add("oldest", params.Oldest)
	***REMOVED***
	if params.Count != DEFAULT_HISTORY_COUNT ***REMOVED***
		values.Add("count", strconv.Itoa(params.Count))
	***REMOVED***
	if params.Inclusive != DEFAULT_HISTORY_INCLUSIVE ***REMOVED***
		if params.Inclusive ***REMOVED***
			values.Add("inclusive", "1")
		***REMOVED*** else ***REMOVED***
			values.Add("inclusive", "0")
		***REMOVED***
	***REMOVED***
	if params.Unreads != DEFAULT_HISTORY_UNREADS ***REMOVED***
		if params.Unreads ***REMOVED***
			values.Add("unreads", "1")
		***REMOVED*** else ***REMOVED***
			values.Add("unreads", "0")
		***REMOVED***
	***REMOVED***

	response, err := imRequest(ctx, api.httpclient, "im.history", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.History, nil
***REMOVED***

// GetIMChannels returns the list of direct message channels
func (api *Client) GetIMChannels() ([]IM, error) ***REMOVED***
	return api.GetIMChannelsContext(context.Background())
***REMOVED***

// GetIMChannelsContext returns the list of direct message channels with a custom context
func (api *Client) GetIMChannelsContext(ctx context.Context) ([]IM, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***

	response, err := imRequest(ctx, api.httpclient, "im.list", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return response.IMs, nil
***REMOVED***
