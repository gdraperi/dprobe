package slack

import (
	"context"
	"errors"
	"net/url"
)

type OAuthResponseIncomingWebhook struct ***REMOVED***
	URL              string `json:"url"`
	Channel          string `json:"channel"`
	ChannelID        string `json:"channel_id,omitempty"`
	ConfigurationURL string `json:"configuration_url"`
***REMOVED***

type OAuthResponseBot struct ***REMOVED***
	BotUserID      string `json:"bot_user_id"`
	BotAccessToken string `json:"bot_access_token"`
***REMOVED***

type OAuthResponse struct ***REMOVED***
	AccessToken     string                       `json:"access_token"`
	Scope           string                       `json:"scope"`
	TeamName        string                       `json:"team_name"`
	TeamID          string                       `json:"team_id"`
	IncomingWebhook OAuthResponseIncomingWebhook `json:"incoming_webhook"`
	Bot             OAuthResponseBot             `json:"bot"`
	UserID          string                       `json:"user_id,omitempty"`
	SlackResponse
***REMOVED***

// GetOAuthToken retrieves an AccessToken
func GetOAuthToken(clientID, clientSecret, code, redirectURI string, debug bool) (accessToken string, scope string, err error) ***REMOVED***
	return GetOAuthTokenContext(context.Background(), clientID, clientSecret, code, redirectURI, debug)
***REMOVED***

// GetOAuthTokenContext retrieves an AccessToken with a custom context
func GetOAuthTokenContext(ctx context.Context, clientID, clientSecret, code, redirectURI string, debug bool) (accessToken string, scope string, err error) ***REMOVED***
	response, err := GetOAuthResponseContext(ctx, clientID, clientSecret, code, redirectURI, debug)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	return response.AccessToken, response.Scope, nil
***REMOVED***

func GetOAuthResponse(clientID, clientSecret, code, redirectURI string, debug bool) (resp *OAuthResponse, err error) ***REMOVED***
	return GetOAuthResponseContext(context.Background(), clientID, clientSecret, code, redirectURI, debug)
***REMOVED***

func GetOAuthResponseContext(ctx context.Context, clientID, clientSecret, code, redirectURI string, debug bool) (resp *OAuthResponse, err error) ***REMOVED***
	values := url.Values***REMOVED***
		"client_id":     ***REMOVED***clientID***REMOVED***,
		"client_secret": ***REMOVED***clientSecret***REMOVED***,
		"code":          ***REMOVED***code***REMOVED***,
		"redirect_uri":  ***REMOVED***redirectURI***REMOVED***,
	***REMOVED***
	response := &OAuthResponse***REMOVED******REMOVED***
	err = post(ctx, customHTTPClient, "oauth.access", values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***
