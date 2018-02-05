package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type SnoozeDebug struct ***REMOVED***
	SnoozeEndDate string `json:"snooze_end_date"`
***REMOVED***

type SnoozeInfo struct ***REMOVED***
	SnoozeEnabled   bool        `json:"snooze_enabled,omitempty"`
	SnoozeEndTime   int         `json:"snooze_endtime,omitempty"`
	SnoozeRemaining int         `json:"snooze_remaining,omitempty"`
	SnoozeDebug     SnoozeDebug `json:"snooze_debug,omitempty"`
***REMOVED***

type DNDStatus struct ***REMOVED***
	Enabled            bool `json:"dnd_enabled"`
	NextStartTimestamp int  `json:"next_dnd_start_ts"`
	NextEndTimestamp   int  `json:"next_dnd_end_ts"`
	SnoozeInfo
***REMOVED***

type dndResponseFull struct ***REMOVED***
	DNDStatus
	SlackResponse
***REMOVED***

type dndTeamInfoResponse struct ***REMOVED***
	Users map[string]DNDStatus `json:"users"`
	SlackResponse
***REMOVED***

func dndRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*dndResponseFull, error) ***REMOVED***
	response := &dndResponseFull***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// EndDND ends the user's scheduled Do Not Disturb session
func (api *Client) EndDND() error ***REMOVED***
	return api.EndDNDContext(context.Background())
***REMOVED***

// EndDNDContext ends the user's scheduled Do Not Disturb session with a custom context
func (api *Client) EndDNDContext(ctx context.Context) error ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***

	response := &SlackResponse***REMOVED******REMOVED***

	if err := post(ctx, api.httpclient, "dnd.endDnd", values, response, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// EndSnooze ends the current user's snooze mode
func (api *Client) EndSnooze() (*DNDStatus, error) ***REMOVED***
	return api.EndSnoozeContext(context.Background())
***REMOVED***

// EndSnoozeContext ends the current user's snooze mode with a custom context
func (api *Client) EndSnoozeContext(ctx context.Context) (*DNDStatus, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***

	response, err := dndRequest(ctx, api.httpclient, "dnd.endSnooze", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.DNDStatus, nil
***REMOVED***

// GetDNDInfo provides information about a user's current Do Not Disturb settings.
func (api *Client) GetDNDInfo(user *string) (*DNDStatus, error) ***REMOVED***
	return api.GetDNDInfoContext(context.Background(), user)
***REMOVED***

// GetDNDInfoContext provides information about a user's current Do Not Disturb settings with a custom context.
func (api *Client) GetDNDInfoContext(ctx context.Context, user *string) (*DNDStatus, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if user != nil ***REMOVED***
		values.Set("user", *user)
	***REMOVED***

	response, err := dndRequest(ctx, api.httpclient, "dnd.info", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.DNDStatus, nil
***REMOVED***

// GetDNDTeamInfo provides information about a user's current Do Not Disturb settings.
func (api *Client) GetDNDTeamInfo(users []string) (map[string]DNDStatus, error) ***REMOVED***
	return api.GetDNDTeamInfoContext(context.Background(), users)
***REMOVED***

// GetDNDTeamInfoContext provides information about a user's current Do Not Disturb settings with a custom context.
func (api *Client) GetDNDTeamInfoContext(ctx context.Context, users []string) (map[string]DNDStatus, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"users": ***REMOVED***strings.Join(users, ",")***REMOVED***,
	***REMOVED***
	response := &dndTeamInfoResponse***REMOVED******REMOVED***

	if err := post(ctx, api.httpclient, "dnd.teamInfo", values, response, api.debug); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response.Users, nil
***REMOVED***

// SetSnooze adjusts the snooze duration for a user's Do Not Disturb
// settings. If a snooze session is not already active for the user, invoking
// this method will begin one for the specified duration.
func (api *Client) SetSnooze(minutes int) (*DNDStatus, error) ***REMOVED***
	return api.SetSnoozeContext(context.Background(), minutes)
***REMOVED***

// SetSnooze adjusts the snooze duration for a user's Do Not Disturb settings with a custom context.
// For more information see the SetSnooze docs
func (api *Client) SetSnoozeContext(ctx context.Context, minutes int) (*DNDStatus, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":       ***REMOVED***api.token***REMOVED***,
		"num_minutes": ***REMOVED***strconv.Itoa(minutes)***REMOVED***,
	***REMOVED***

	response, err := dndRequest(ctx, api.httpclient, "dnd.setSnooze", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.DNDStatus, nil
***REMOVED***
