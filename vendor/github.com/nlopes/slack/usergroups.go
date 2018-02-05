package slack

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

// UserGroup contains all the information of a user group
type UserGroup struct ***REMOVED***
	ID          string         `json:"id"`
	TeamID      string         `json:"team_id"`
	IsUserGroup bool           `json:"is_usergroup"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Handle      string         `json:"handle"`
	IsExternal  bool           `json:"is_external"`
	DateCreate  JSONTime       `json:"date_create"`
	DateUpdate  JSONTime       `json:"date_update"`
	DateDelete  JSONTime       `json:"date_delete"`
	AutoType    string         `json:"auto_type"`
	CreatedBy   string         `json:"created_by"`
	UpdatedBy   string         `json:"updated_by"`
	DeletedBy   string         `json:"deleted_by"`
	Prefs       UserGroupPrefs `json:"prefs"`
	UserCount   int            `json:"user_count"`
***REMOVED***

// UserGroupPrefs contains default channels and groups (private channels)
type UserGroupPrefs struct ***REMOVED***
	Channels []string `json:"channels"`
	Groups   []string `json:"groups"`
***REMOVED***

type userGroupResponseFull struct ***REMOVED***
	UserGroups []UserGroup `json:"usergroups"`
	UserGroup  UserGroup   `json:"usergroup"`
	Users      []string    `json:"users"`
	SlackResponse
***REMOVED***

func userGroupRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*userGroupResponseFull, error) ***REMOVED***
	response := &userGroupResponseFull***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// CreateUserGroup creates a new user group
func (api *Client) CreateUserGroup(userGroup UserGroup) (UserGroup, error) ***REMOVED***
	return api.CreateUserGroupContext(context.Background(), userGroup)
***REMOVED***

// CreateUserGroupContext creates a new user group with a custom context
func (api *Client) CreateUserGroupContext(ctx context.Context, userGroup UserGroup) (UserGroup, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"name":  ***REMOVED***userGroup.Name***REMOVED***,
	***REMOVED***

	if userGroup.Handle != "" ***REMOVED***
		values["handle"] = []string***REMOVED***userGroup.Handle***REMOVED***
	***REMOVED***

	if userGroup.Description != "" ***REMOVED***
		values["description"] = []string***REMOVED***userGroup.Description***REMOVED***
	***REMOVED***

	if len(userGroup.Prefs.Channels) > 0 ***REMOVED***
		values["channels"] = []string***REMOVED***strings.Join(userGroup.Prefs.Channels, ",")***REMOVED***
	***REMOVED***

	response, err := userGroupRequest(ctx, api.httpclient, "usergroups.create", values, api.debug)
	if err != nil ***REMOVED***
		return UserGroup***REMOVED******REMOVED***, err
	***REMOVED***
	return response.UserGroup, nil
***REMOVED***

// DisableUserGroup disables an existing user group
func (api *Client) DisableUserGroup(userGroup string) (UserGroup, error) ***REMOVED***
	return api.DisableUserGroupContext(context.Background(), userGroup)
***REMOVED***

// DisableUserGroupContext disables an existing user group with a custom context
func (api *Client) DisableUserGroupContext(ctx context.Context, userGroup string) (UserGroup, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":     ***REMOVED***api.token***REMOVED***,
		"usergroup": ***REMOVED***userGroup***REMOVED***,
	***REMOVED***

	response, err := userGroupRequest(ctx, api.httpclient, "usergroups.disable", values, api.debug)
	if err != nil ***REMOVED***
		return UserGroup***REMOVED******REMOVED***, err
	***REMOVED***
	return response.UserGroup, nil
***REMOVED***

// EnableUserGroup enables an existing user group
func (api *Client) EnableUserGroup(userGroup string) (UserGroup, error) ***REMOVED***
	return api.EnableUserGroupContext(context.Background(), userGroup)
***REMOVED***

// EnableUserGroupContext enables an existing user group with a custom context
func (api *Client) EnableUserGroupContext(ctx context.Context, userGroup string) (UserGroup, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":     ***REMOVED***api.token***REMOVED***,
		"usergroup": ***REMOVED***userGroup***REMOVED***,
	***REMOVED***

	response, err := userGroupRequest(ctx, api.httpclient, "usergroups.enable", values, api.debug)
	if err != nil ***REMOVED***
		return UserGroup***REMOVED******REMOVED***, err
	***REMOVED***
	return response.UserGroup, nil
***REMOVED***

// GetUserGroups returns a list of user groups for the team
func (api *Client) GetUserGroups() ([]UserGroup, error) ***REMOVED***
	return api.GetUserGroupsContext(context.Background())
***REMOVED***

// GetUserGroupsContext returns a list of user groups for the team with a custom context
func (api *Client) GetUserGroupsContext(ctx context.Context) ([]UserGroup, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***

	response, err := userGroupRequest(ctx, api.httpclient, "usergroups.list", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return response.UserGroups, nil
***REMOVED***

// UpdateUserGroup will update an existing user group
func (api *Client) UpdateUserGroup(userGroup UserGroup) (UserGroup, error) ***REMOVED***
	return api.UpdateUserGroupContext(context.Background(), userGroup)
***REMOVED***

// UpdateUserGroupContext will update an existing user group with a custom context
func (api *Client) UpdateUserGroupContext(ctx context.Context, userGroup UserGroup) (UserGroup, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":     ***REMOVED***api.token***REMOVED***,
		"usergroup": ***REMOVED***userGroup.ID***REMOVED***,
	***REMOVED***

	if userGroup.Name != "" ***REMOVED***
		values["name"] = []string***REMOVED***userGroup.Name***REMOVED***
	***REMOVED***

	if userGroup.Handle != "" ***REMOVED***
		values["handle"] = []string***REMOVED***userGroup.Handle***REMOVED***
	***REMOVED***

	if userGroup.Description != "" ***REMOVED***
		values["description"] = []string***REMOVED***userGroup.Description***REMOVED***
	***REMOVED***

	response, err := userGroupRequest(ctx, api.httpclient, "usergroups.update", values, api.debug)
	if err != nil ***REMOVED***
		return UserGroup***REMOVED******REMOVED***, err
	***REMOVED***
	return response.UserGroup, nil
***REMOVED***

// GetUserGroupMembers will retrieve the current list of users in a group
func (api *Client) GetUserGroupMembers(userGroup string) ([]string, error) ***REMOVED***
	return api.GetUserGroupMembersContext(context.Background(), userGroup)
***REMOVED***

// GetUserGroupMembersContext will retrieve the current list of users in a group with a custom context
func (api *Client) GetUserGroupMembersContext(ctx context.Context, userGroup string) ([]string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":     ***REMOVED***api.token***REMOVED***,
		"usergroup": ***REMOVED***userGroup***REMOVED***,
	***REMOVED***

	response, err := userGroupRequest(ctx, api.httpclient, "usergroups.users.list", values, api.debug)
	if err != nil ***REMOVED***
		return []string***REMOVED******REMOVED***, err
	***REMOVED***
	return response.Users, nil
***REMOVED***

// UpdateUserGroupMembers will update the members of an existing user group
func (api *Client) UpdateUserGroupMembers(userGroup string, members string) (UserGroup, error) ***REMOVED***
	return api.UpdateUserGroupMembersContext(context.Background(), userGroup, members)
***REMOVED***

// UpdateUserGroupMembersContext will update the members of an existing user group with a custom context
func (api *Client) UpdateUserGroupMembersContext(ctx context.Context, userGroup string, members string) (UserGroup, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":     ***REMOVED***api.token***REMOVED***,
		"usergroup": ***REMOVED***userGroup***REMOVED***,
		"users":     ***REMOVED***members***REMOVED***,
	***REMOVED***

	response, err := userGroupRequest(ctx, api.httpclient, "usergroups.users.update", values, api.debug)
	if err != nil ***REMOVED***
		return UserGroup***REMOVED******REMOVED***, err
	***REMOVED***
	return response.UserGroup, nil
***REMOVED***
