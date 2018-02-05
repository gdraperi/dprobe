package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

// ItemReaction is the reactions that have happened on an item.
type ItemReaction struct ***REMOVED***
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Users []string `json:"users"`
***REMOVED***

// ReactedItem is an item that was reacted to, and the details of the
// reactions.
type ReactedItem struct ***REMOVED***
	Item
	Reactions []ItemReaction
***REMOVED***

// GetReactionsParameters is the inputs to get reactions to an item.
type GetReactionsParameters struct ***REMOVED***
	Full bool
***REMOVED***

// NewGetReactionsParameters initializes the inputs to get reactions to an item.
func NewGetReactionsParameters() GetReactionsParameters ***REMOVED***
	return GetReactionsParameters***REMOVED***
		Full: false,
	***REMOVED***
***REMOVED***

type getReactionsResponseFull struct ***REMOVED***
	Type string
	M    struct ***REMOVED***
		Reactions []ItemReaction
	***REMOVED*** `json:"message"`
	F struct ***REMOVED***
		Reactions []ItemReaction
	***REMOVED*** `json:"file"`
	FC struct ***REMOVED***
		Reactions []ItemReaction
	***REMOVED*** `json:"comment"`
	SlackResponse
***REMOVED***

func (res getReactionsResponseFull) extractReactions() []ItemReaction ***REMOVED***
	switch res.Type ***REMOVED***
	case "message":
		return res.M.Reactions
	case "file":
		return res.F.Reactions
	case "file_comment":
		return res.FC.Reactions
	***REMOVED***
	return []ItemReaction***REMOVED******REMOVED***
***REMOVED***

const (
	DEFAULT_REACTIONS_USER  = ""
	DEFAULT_REACTIONS_COUNT = 100
	DEFAULT_REACTIONS_PAGE  = 1
	DEFAULT_REACTIONS_FULL  = false
)

// ListReactionsParameters is the inputs to find all reactions by a user.
type ListReactionsParameters struct ***REMOVED***
	User  string
	Count int
	Page  int
	Full  bool
***REMOVED***

// NewListReactionsParameters initializes the inputs to find all reactions
// performed by a user.
func NewListReactionsParameters() ListReactionsParameters ***REMOVED***
	return ListReactionsParameters***REMOVED***
		User:  DEFAULT_REACTIONS_USER,
		Count: DEFAULT_REACTIONS_COUNT,
		Page:  DEFAULT_REACTIONS_PAGE,
		Full:  DEFAULT_REACTIONS_FULL,
	***REMOVED***
***REMOVED***

type listReactionsResponseFull struct ***REMOVED***
	Items []struct ***REMOVED***
		Type    string
		Channel string
		M       struct ***REMOVED***
			*Message
		***REMOVED*** `json:"message"`
		F struct ***REMOVED***
			*File
			Reactions []ItemReaction
		***REMOVED*** `json:"file"`
		FC struct ***REMOVED***
			*Comment
			Reactions []ItemReaction
		***REMOVED*** `json:"comment"`
	***REMOVED***
	Paging `json:"paging"`
	SlackResponse
***REMOVED***

func (res listReactionsResponseFull) extractReactedItems() []ReactedItem ***REMOVED***
	items := make([]ReactedItem, len(res.Items))
	for i, input := range res.Items ***REMOVED***
		item := ReactedItem***REMOVED******REMOVED***
		item.Type = input.Type
		switch input.Type ***REMOVED***
		case "message":
			item.Channel = input.Channel
			item.Message = input.M.Message
			item.Reactions = input.M.Reactions
		case "file":
			item.File = input.F.File
			item.Reactions = input.F.Reactions
		case "file_comment":
			item.File = input.F.File
			item.Comment = input.FC.Comment
			item.Reactions = input.FC.Reactions
		***REMOVED***
		items[i] = item
	***REMOVED***
	return items
***REMOVED***

// AddReaction adds a reaction emoji to a message, file or file comment.
func (api *Client) AddReaction(name string, item ItemRef) error ***REMOVED***
	return api.AddReactionContext(context.Background(), name, item)
***REMOVED***

// AddReactionContext adds a reaction emoji to a message, file or file comment with a custom context.
func (api *Client) AddReactionContext(ctx context.Context, name string, item ItemRef) error ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if name != "" ***REMOVED***
		values.Set("name", name)
	***REMOVED***
	if item.Channel != "" ***REMOVED***
		values.Set("channel", string(item.Channel))
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
	if err := post(ctx, api.httpclient, "reactions.add", values, response, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// RemoveReaction removes a reaction emoji from a message, file or file comment.
func (api *Client) RemoveReaction(name string, item ItemRef) error ***REMOVED***
	return api.RemoveReactionContext(context.Background(), name, item)
***REMOVED***

// RemoveReactionContext removes a reaction emoji from a message, file or file comment with a custom context.
func (api *Client) RemoveReactionContext(ctx context.Context, name string, item ItemRef) error ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if name != "" ***REMOVED***
		values.Set("name", name)
	***REMOVED***
	if item.Channel != "" ***REMOVED***
		values.Set("channel", string(item.Channel))
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
	if err := post(ctx, api.httpclient, "reactions.remove", values, response, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// GetReactions returns details about the reactions on an item.
func (api *Client) GetReactions(item ItemRef, params GetReactionsParameters) ([]ItemReaction, error) ***REMOVED***
	return api.GetReactionsContext(context.Background(), item, params)
***REMOVED***

// GetReactionsContext returns details about the reactions on an item with a custom context
func (api *Client) GetReactionsContext(ctx context.Context, item ItemRef, params GetReactionsParameters) ([]ItemReaction, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if item.Channel != "" ***REMOVED***
		values.Set("channel", string(item.Channel))
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
	if params.Full != DEFAULT_REACTIONS_FULL ***REMOVED***
		values.Set("full", strconv.FormatBool(params.Full))
	***REMOVED***

	response := &getReactionsResponseFull***REMOVED******REMOVED***
	if err := post(ctx, api.httpclient, "reactions.get", values, response, api.debug); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response.extractReactions(), nil
***REMOVED***

// ListReactions returns information about the items a user reacted to.
func (api *Client) ListReactions(params ListReactionsParameters) ([]ReactedItem, *Paging, error) ***REMOVED***
	return api.ListReactionsContext(context.Background(), params)
***REMOVED***

// ListReactionsContext returns information about the items a user reacted to with a custom context.
func (api *Client) ListReactionsContext(ctx context.Context, params ListReactionsParameters) ([]ReactedItem, *Paging, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if params.User != DEFAULT_REACTIONS_USER ***REMOVED***
		values.Add("user", params.User)
	***REMOVED***
	if params.Count != DEFAULT_REACTIONS_COUNT ***REMOVED***
		values.Add("count", strconv.Itoa(params.Count))
	***REMOVED***
	if params.Page != DEFAULT_REACTIONS_PAGE ***REMOVED***
		values.Add("page", strconv.Itoa(params.Page))
	***REMOVED***
	if params.Full != DEFAULT_REACTIONS_FULL ***REMOVED***
		values.Add("full", strconv.FormatBool(params.Full))
	***REMOVED***

	response := &listReactionsResponseFull***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "reactions.list", values, response, api.debug)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, nil, errors.New(response.Error)
	***REMOVED***
	return response.extractReactedItems(), &response.Paging, nil
***REMOVED***
