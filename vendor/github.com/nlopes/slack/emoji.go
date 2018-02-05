package slack

import (
	"context"
	"errors"
	"net/url"
)

type emojiResponseFull struct ***REMOVED***
	Emoji map[string]string `json:"emoji"`
	SlackResponse
***REMOVED***

// GetEmoji retrieves all the emojis
func (api *Client) GetEmoji() (map[string]string, error) ***REMOVED***
	return api.GetEmojiContext(context.Background())
***REMOVED***

// GetEmojiContext retrieves all the emojis with a custom context
func (api *Client) GetEmojiContext(ctx context.Context) (map[string]string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	response := &emojiResponseFull***REMOVED******REMOVED***

	err := post(ctx, api.httpclient, "emoji.list", values, response, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response.Emoji, nil
***REMOVED***
