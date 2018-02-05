package slack

import (
	"context"
	"errors"
	"net/url"
)

// Bot contains information about a bot
type Bot struct ***REMOVED***
	ID      string `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
	Icons   Icons  `json:"icons"`
***REMOVED***

type botResponseFull struct ***REMOVED***
	Bot `json:"bot,omitempty"` // GetBotInfo
	SlackResponse
***REMOVED***

func botRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*botResponseFull, error) ***REMOVED***
	response := &botResponseFull***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// GetBotInfo will retrieve the complete bot information
func (api *Client) GetBotInfo(bot string) (*Bot, error) ***REMOVED***
	return api.GetBotInfoContext(context.Background(), bot)
***REMOVED***

// GetBotInfoContext will retrieve the complete bot information using a custom context
func (api *Client) GetBotInfoContext(ctx context.Context, bot string) (*Bot, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"bot":   ***REMOVED***bot***REMOVED***,
	***REMOVED***

	response, err := botRequest(ctx, api.httpclient, "bots.info", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Bot, nil
***REMOVED***
