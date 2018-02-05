package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

const (
	DEFAULT_STARS_USER  = ""
	DEFAULT_STARS_COUNT = 100
	DEFAULT_STARS_PAGE  = 1
)

type StarsParameters struct ***REMOVED***
	User  string
	Count int
	Page  int
***REMOVED***

type StarredItem Item

type listResponseFull struct ***REMOVED***
	Items  []Item `json:"items"`
	Paging `json:"paging"`
	SlackResponse
***REMOVED***

// NewStarsParameters initialises StarsParameters with default values
func NewStarsParameters() StarsParameters ***REMOVED***
	return StarsParameters***REMOVED***
		User:  DEFAULT_STARS_USER,
		Count: DEFAULT_STARS_COUNT,
		Page:  DEFAULT_STARS_PAGE,
	***REMOVED***
***REMOVED***

// AddStar stars an item in a channel
func (api *Client) AddStar(channel string, item ItemRef) error ***REMOVED***
	return api.AddStarContext(context.Background(), channel, item)
***REMOVED***

// AddStarContext stars an item in a channel with a custom context
func (api *Client) AddStarContext(ctx context.Context, channel string, item ItemRef) error ***REMOVED***
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
	if err := post(ctx, api.httpclient, "stars.add", values, response, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// RemoveStar removes a starred item from a channel
func (api *Client) RemoveStar(channel string, item ItemRef) error ***REMOVED***
	return api.RemoveStarContext(context.Background(), channel, item)
***REMOVED***

// RemoveStarContext removes a starred item from a channel with a custom context
func (api *Client) RemoveStarContext(ctx context.Context, channel string, item ItemRef) error ***REMOVED***
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
	if err := post(ctx, api.httpclient, "stars.remove", values, response, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// ListStars returns information about the stars a user added
func (api *Client) ListStars(params StarsParameters) ([]Item, *Paging, error) ***REMOVED***
	return api.ListStarsContext(context.Background(), params)
***REMOVED***

// ListStarsContext returns information about the stars a user added with a custom context
func (api *Client) ListStarsContext(ctx context.Context, params StarsParameters) ([]Item, *Paging, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if params.User != DEFAULT_STARS_USER ***REMOVED***
		values.Add("user", params.User)
	***REMOVED***
	if params.Count != DEFAULT_STARS_COUNT ***REMOVED***
		values.Add("count", strconv.Itoa(params.Count))
	***REMOVED***
	if params.Page != DEFAULT_STARS_PAGE ***REMOVED***
		values.Add("page", strconv.Itoa(params.Page))
	***REMOVED***

	response := &listResponseFull***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "stars.list", values, response, api.debug)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, nil, errors.New(response.Error)
	***REMOVED***
	return response.Items, &response.Paging, nil
***REMOVED***

// GetStarred returns a list of StarredItem items.
//
// The user then has to iterate over them and figure out what they should
// be looking at according to what is in the Type.
//    for _, item := range items ***REMOVED***
//        switch c.Type ***REMOVED***
//        case "file_comment":
//            log.Println(c.Comment)
//        case "file":
//             ...
//
//***REMOVED***
// This function still exists to maintain backwards compatibility.
// I exposed it as returning []StarredItem, so it shall stay as StarredItem
func (api *Client) GetStarred(params StarsParameters) ([]StarredItem, *Paging, error) ***REMOVED***
	return api.GetStarredContext(context.Background(), params)
***REMOVED***

// GetStarredContext returns a list of StarredItem items with a custom context
//
// For more details see GetStarred
func (api *Client) GetStarredContext(ctx context.Context, params StarsParameters) ([]StarredItem, *Paging, error) ***REMOVED***
	items, paging, err := api.ListStarsContext(ctx, params)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	starredItems := make([]StarredItem, len(items))
	for i, item := range items ***REMOVED***
		starredItems[i] = StarredItem(item)
	***REMOVED***
	return starredItems, paging, nil
***REMOVED***
