package slack

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

type adminResponse struct ***REMOVED***
	OK    bool   `json:"ok"`
	Error string `json:"error"`
***REMOVED***

func adminRequest(ctx context.Context, client HTTPRequester, method string, teamName string, values url.Values, debug bool) (*adminResponse, error) ***REMOVED***
	adminResponse := &adminResponse***REMOVED******REMOVED***
	err := parseAdminResponse(ctx, client, method, teamName, values, adminResponse, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !adminResponse.OK ***REMOVED***
		return nil, errors.New(adminResponse.Error)
	***REMOVED***

	return adminResponse, nil
***REMOVED***

// DisableUser disabled a user account, given a user ID
func (api *Client) DisableUser(teamName string, uid string) error ***REMOVED***
	return api.DisableUserContext(context.Background(), teamName, uid)
***REMOVED***

// DisableUserContext disabled a user account, given a user ID with a custom context
func (api *Client) DisableUserContext(ctx context.Context, teamName string, uid string) error ***REMOVED***
	values := url.Values***REMOVED***
		"user":       ***REMOVED***uid***REMOVED***,
		"token":      ***REMOVED***api.token***REMOVED***,
		"set_active": ***REMOVED***"true"***REMOVED***,
		"_attempts":  ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "setInactive", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to disable user with id '%s': %s", uid, err)
	***REMOVED***

	return nil
***REMOVED***

// InviteGuest invites a user to Slack as a single-channel guest
func (api *Client) InviteGuest(teamName, channel, firstName, lastName, emailAddress string) error ***REMOVED***
	return api.InviteGuestContext(context.Background(), teamName, channel, firstName, lastName, emailAddress)
***REMOVED***

// InviteGuestContext invites a user to Slack as a single-channel guest with a custom context
func (api *Client) InviteGuestContext(ctx context.Context, teamName, channel, firstName, lastName, emailAddress string) error ***REMOVED***
	values := url.Values***REMOVED***
		"email":            ***REMOVED***emailAddress***REMOVED***,
		"channels":         ***REMOVED***channel***REMOVED***,
		"first_name":       ***REMOVED***firstName***REMOVED***,
		"last_name":        ***REMOVED***lastName***REMOVED***,
		"ultra_restricted": ***REMOVED***"1"***REMOVED***,
		"token":            ***REMOVED***api.token***REMOVED***,
		"set_active":       ***REMOVED***"true"***REMOVED***,
		"_attempts":        ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "invite", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to invite single-channel guest: %s", err)
	***REMOVED***

	return nil
***REMOVED***

// InviteRestricted invites a user to Slack as a restricted account
func (api *Client) InviteRestricted(teamName, channel, firstName, lastName, emailAddress string) error ***REMOVED***
	return api.InviteRestrictedContext(context.Background(), teamName, channel, firstName, lastName, emailAddress)
***REMOVED***

// InviteRestrictedContext invites a user to Slack as a restricted account with a custom context
func (api *Client) InviteRestrictedContext(ctx context.Context, teamName, channel, firstName, lastName, emailAddress string) error ***REMOVED***
	values := url.Values***REMOVED***
		"email":      ***REMOVED***emailAddress***REMOVED***,
		"channels":   ***REMOVED***channel***REMOVED***,
		"first_name": ***REMOVED***firstName***REMOVED***,
		"last_name":  ***REMOVED***lastName***REMOVED***,
		"restricted": ***REMOVED***"1"***REMOVED***,
		"token":      ***REMOVED***api.token***REMOVED***,
		"set_active": ***REMOVED***"true"***REMOVED***,
		"_attempts":  ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "invite", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to restricted account: %s", err)
	***REMOVED***

	return nil
***REMOVED***

// InviteToTeam invites a user to a Slack team
func (api *Client) InviteToTeam(teamName, firstName, lastName, emailAddress string) error ***REMOVED***
	return api.InviteToTeamContext(context.Background(), teamName, firstName, lastName, emailAddress)
***REMOVED***

// InviteToTeamContext invites a user to a Slack team with a custom context
func (api *Client) InviteToTeamContext(ctx context.Context, teamName, firstName, lastName, emailAddress string) error ***REMOVED***
	values := url.Values***REMOVED***
		"email":      ***REMOVED***emailAddress***REMOVED***,
		"first_name": ***REMOVED***firstName***REMOVED***,
		"last_name":  ***REMOVED***lastName***REMOVED***,
		"token":      ***REMOVED***api.token***REMOVED***,
		"set_active": ***REMOVED***"true"***REMOVED***,
		"_attempts":  ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "invite", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to invite to team: %s", err)
	***REMOVED***

	return nil
***REMOVED***

// SetRegular enables the specified user
func (api *Client) SetRegular(teamName, user string) error ***REMOVED***
	return api.SetRegularContext(context.Background(), teamName, user)
***REMOVED***

// SetRegularContext enables the specified user with a custom context
func (api *Client) SetRegularContext(ctx context.Context, teamName, user string) error ***REMOVED***
	values := url.Values***REMOVED***
		"user":       ***REMOVED***user***REMOVED***,
		"token":      ***REMOVED***api.token***REMOVED***,
		"set_active": ***REMOVED***"true"***REMOVED***,
		"_attempts":  ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "setRegular", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to change the user (%s) to a regular user: %s", user, err)
	***REMOVED***

	return nil
***REMOVED***

// SendSSOBindingEmail sends an SSO binding email to the specified user
func (api *Client) SendSSOBindingEmail(teamName, user string) error ***REMOVED***
	return api.SendSSOBindingEmailContext(context.Background(), teamName, user)
***REMOVED***

// SendSSOBindingEmailContext sends an SSO binding email to the specified user with a custom context
func (api *Client) SendSSOBindingEmailContext(ctx context.Context, teamName, user string) error ***REMOVED***
	values := url.Values***REMOVED***
		"user":       ***REMOVED***user***REMOVED***,
		"token":      ***REMOVED***api.token***REMOVED***,
		"set_active": ***REMOVED***"true"***REMOVED***,
		"_attempts":  ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "sendSSOBind", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to send SSO binding email for user (%s): %s", user, err)
	***REMOVED***

	return nil
***REMOVED***

// SetUltraRestricted converts a user into a single-channel guest
func (api *Client) SetUltraRestricted(teamName, uid, channel string) error ***REMOVED***
	return api.SetUltraRestrictedContext(context.Background(), teamName, uid, channel)
***REMOVED***

// SetUltraRestrictedContext converts a user into a single-channel guest with a custom context
func (api *Client) SetUltraRestrictedContext(ctx context.Context, teamName, uid, channel string) error ***REMOVED***
	values := url.Values***REMOVED***
		"user":       ***REMOVED***uid***REMOVED***,
		"channel":    ***REMOVED***channel***REMOVED***,
		"token":      ***REMOVED***api.token***REMOVED***,
		"set_active": ***REMOVED***"true"***REMOVED***,
		"_attempts":  ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "setUltraRestricted", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to ultra-restrict account: %s", err)
	***REMOVED***

	return nil
***REMOVED***

// SetRestricted converts a user into a restricted account
func (api *Client) SetRestricted(teamName, uid string) error ***REMOVED***
	return api.SetRestrictedContext(context.Background(), teamName, uid)
***REMOVED***

// SetRestrictedContext converts a user into a restricted account with a custom context
func (api *Client) SetRestrictedContext(ctx context.Context, teamName, uid string) error ***REMOVED***
	values := url.Values***REMOVED***
		"user":       ***REMOVED***uid***REMOVED***,
		"token":      ***REMOVED***api.token***REMOVED***,
		"set_active": ***REMOVED***"true"***REMOVED***,
		"_attempts":  ***REMOVED***"1"***REMOVED***,
	***REMOVED***

	_, err := adminRequest(ctx, api.httpclient, "setRestricted", teamName, values, api.debug)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to restrict account: %s", err)
	***REMOVED***

	return nil
***REMOVED***
