package slack

import (
	"context"
	"errors"
	"net/url"
)

type listPinsResponseFull struct ***REMOVED***
	Items  []Item
	Paging `json:"paging"`
	SlackResponse
***REMOVED***

// AddPin pins an item in a channel
func (api *Client) AddPin(channel string, item ItemRef) error ***REMOVED***
	return api.AddPinContext(context.Background(), channel, item)
***REMOVED***

// AddPinContext pins an item in a channel with a custom context
func (api *Client) AddPinContext(ctx context.Context, channel string, item ItemRef) error ***REMOVED***
	values := url.Values***REMOVED***
		"channel": ***REMOVED***channel***REMOVED***,
		"token":   ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if item.Timestamp != "" ***REMOVED***
		values.Set("timestamp", string(item.Timestamp))
	***REMOVED***
	if item.File != "" ***REMOVED***
		values.Set("file", string(item.File))
	***REMOVED***
	if item.Comment != "" ***REMOVED***
		values.Set("file_comment", string(item.Comment))
	***REMOVED***

	response := &SlackResponse***REMOVED******REMOVED***
	if err := post(ctx, api.httpclient, "pins.add", values, response, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// RemovePin un-pins an item from a channel
func (api *Client) RemovePin(channel string, item ItemRef) error ***REMOVED***
	return api.RemovePinContext(context.Background(), channel, item)
***REMOVED***

// RemovePinContext un-pins an item from a channel with a custom context
func (api *Client) RemovePinContext(ctx context.Context, channel string, item ItemRef) error ***REMOVED***
	values := url.Values***REMOVED***
		"channel": ***REMOVED***channel***REMOVED***,
		"token":   ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if item.Timestamp != "" ***REMOVED***
		values.Set("timestamp", string(item.Timestamp))
	***REMOVED***
	if item.File != "" ***REMOVED***
		values.Set("file", string(item.File))
	***REMOVED***
	if item.Comment != "" ***REMOVED***
		values.Set("file_comment", string(item.Comment))
	***REMOVED***

	response := &SlackResponse***REMOVED******REMOVED***
	if err := post(ctx, api.httpclient, "pins.remove", values, response, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// ListPins returns information about the items a user reacted to.
func (api *Client) ListPins(channel string) ([]Item, *Paging, error) ***REMOVED***
	return api.ListPinsContext(context.Background(), channel)
***REMOVED***

// ListPinsContext returns information about the items a user reacted to with a custom context.
func (api *Client) ListPinsContext(ctx context.Context, channel string) ([]Item, *Paging, error) ***REMOVED***
	values := url.Values***REMOVED***
		"channel": ***REMOVED***channel***REMOVED***,
		"token":   ***REMOVED***api.token***REMOVED***,
	***REMOVED***

	response := &listPinsResponseFull***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "pins.list", values, response, api.debug)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, nil, errors.New(response.Error)
	***REMOVED***
	return response.Items, &response.Paging, nil
***REMOVED***
