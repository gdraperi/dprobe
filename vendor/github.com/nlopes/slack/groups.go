package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

// Group contains all the information for a group
type Group struct ***REMOVED***
	groupConversation
	IsGroup bool `json:"is_group"`
***REMOVED***

type groupResponseFull struct ***REMOVED***
	Group          Group   `json:"group"`
	Groups         []Group `json:"groups"`
	Purpose        string  `json:"purpose"`
	Topic          string  `json:"topic"`
	NotInGroup     bool    `json:"not_in_group"`
	NoOp           bool    `json:"no_op"`
	AlreadyClosed  bool    `json:"already_closed"`
	AlreadyOpen    bool    `json:"already_open"`
	AlreadyInGroup bool    `json:"already_in_group"`
	Channel        Channel `json:"channel"`
	History
	SlackResponse
***REMOVED***

func groupRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*groupResponseFull, error) ***REMOVED***
	response := &groupResponseFull***REMOVED******REMOVED***
	err := postForm(ctx, client, SLACK_API+path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// ArchiveGroup archives a private group
func (api *Client) ArchiveGroup(group string) error ***REMOVED***
	return api.ArchiveGroupContext(context.Background(), group)
***REMOVED***

// ArchiveGroupContext archives a private group
func (api *Client) ArchiveGroupContext(ctx context.Context, group string) error ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
	***REMOVED***

	_, err := groupRequest(ctx, api.httpclient, "groups.archive", values, api.debug)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return err
***REMOVED***

// UnarchiveGroup unarchives a private group
func (api *Client) UnarchiveGroup(group string) error ***REMOVED***
	return api.UnarchiveGroupContext(context.Background(), group)
***REMOVED***

// UnarchiveGroupContext unarchives a private group
func (api *Client) UnarchiveGroupContext(ctx context.Context, group string) error ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
	***REMOVED***

	_, err := groupRequest(ctx, api.httpclient, "groups.unarchive", values, api.debug)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// CreateGroup creates a private group
func (api *Client) CreateGroup(group string) (*Group, error) ***REMOVED***
	return api.CreateGroupContext(context.Background(), group)
***REMOVED***

// CreateGroupContext creates a private group
func (api *Client) CreateGroupContext(ctx context.Context, group string) (*Group, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"name":  ***REMOVED***group***REMOVED***,
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.create", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Group, nil
***REMOVED***

// CreateChildGroup creates a new private group archiving the old one
// This method takes an existing private group and performs the following steps:
//   1. Renames the existing group (from "example" to "example-archived").
//   2. Archives the existing group.
//   3. Creates a new group with the name of the existing group.
//   4. Adds all members of the existing group to the new group.
func (api *Client) CreateChildGroup(group string) (*Group, error) ***REMOVED***
	return api.CreateChildGroupContext(context.Background(), group)
***REMOVED***

// CreateChildGroupContext creates a new private group archiving the old one with a custom context
// For more information see CreateChildGroup
func (api *Client) CreateChildGroupContext(ctx context.Context, group string) (*Group, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.createChild", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Group, nil
***REMOVED***

// CloseGroup closes a private group
func (api *Client) CloseGroup(group string) (bool, bool, error) ***REMOVED***
	return api.CloseGroupContext(context.Background(), group)
***REMOVED***

// CloseGroupContext closes a private group with a custom context
func (api *Client) CloseGroupContext(ctx context.Context, group string) (bool, bool, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
	***REMOVED***

	response, err := imRequest(ctx, api.httpclient, "groups.close", values, api.debug)
	if err != nil ***REMOVED***
		return false, false, err
	***REMOVED***
	return response.NoOp, response.AlreadyClosed, nil
***REMOVED***

// GetGroupHistory fetches all the history for a private group
func (api *Client) GetGroupHistory(group string, params HistoryParameters) (*History, error) ***REMOVED***
	return api.GetGroupHistoryContext(context.Background(), group, params)
***REMOVED***

// GetGroupHistoryContext fetches all the history for a private group with a custom context
func (api *Client) GetGroupHistoryContext(ctx context.Context, group string, params HistoryParameters) (*History, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
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

	response, err := groupRequest(ctx, api.httpclient, "groups.history", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.History, nil
***REMOVED***

// InviteUserToGroup invites a specific user to a private group
func (api *Client) InviteUserToGroup(group, user string) (*Group, bool, error) ***REMOVED***
	return api.InviteUserToGroupContext(context.Background(), group, user)
***REMOVED***

// InviteUserToGroupContext invites a specific user to a private group with a custom context
func (api *Client) InviteUserToGroupContext(ctx context.Context, group, user string) (*Group, bool, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
		"user":    ***REMOVED***user***REMOVED***,
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.invite", values, api.debug)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	return &response.Group, response.AlreadyInGroup, nil
***REMOVED***

// LeaveGroup makes authenticated user leave the group
func (api *Client) LeaveGroup(group string) error ***REMOVED***
	return api.LeaveGroupContext(context.Background(), group)
***REMOVED***

// LeaveGroupContext makes authenticated user leave the group with a custom context
func (api *Client) LeaveGroupContext(ctx context.Context, group string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
	***REMOVED***

	if _, err = groupRequest(ctx, api.httpclient, "groups.leave", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// KickUserFromGroup kicks a user from a group
func (api *Client) KickUserFromGroup(group, user string) error ***REMOVED***
	return api.KickUserFromGroupContext(context.Background(), group, user)
***REMOVED***

// KickUserFromGroupContext kicks a user from a group with a custom context
func (api *Client) KickUserFromGroupContext(ctx context.Context, group, user string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
		"user":    ***REMOVED***user***REMOVED***,
	***REMOVED***

	if _, err = groupRequest(ctx, api.httpclient, "groups.kick", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// GetGroups retrieves all groups
func (api *Client) GetGroups(excludeArchived bool) ([]Group, error) ***REMOVED***
	return api.GetGroupsContext(context.Background(), excludeArchived)
***REMOVED***

// GetGroupsContext retrieves all groups with a custom context
func (api *Client) GetGroupsContext(ctx context.Context, excludeArchived bool) ([]Group, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if excludeArchived ***REMOVED***
		values.Add("exclude_archived", "1")
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.list", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return response.Groups, nil
***REMOVED***

// GetGroupInfo retrieves the given group
func (api *Client) GetGroupInfo(group string) (*Group, error) ***REMOVED***
	return api.GetGroupInfoContext(context.Background(), group)
***REMOVED***

// GetGroupInfoContext retrieves the given group with a custom context
func (api *Client) GetGroupInfoContext(ctx context.Context, group string) (*Group, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.info", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Group, nil
***REMOVED***

// SetGroupReadMark sets the read mark on a private group
// Clients should try to avoid making this call too often. When needing to mark a read position, a client should set a
// timer before making the call. In this way, any further updates needed during the timeout will not generate extra
// calls (just one per channel). This is useful for when reading scroll-back history, or following a busy live
// channel. A timeout of 5 seconds is a good starting point. Be sure to flush these calls on shutdown/logout.
func (api *Client) SetGroupReadMark(group, ts string) error ***REMOVED***
	return api.SetGroupReadMarkContext(context.Background(), group, ts)
***REMOVED***

// SetGroupReadMarkContext sets the read mark on a private group with a custom context
// For more details see SetGroupReadMark
func (api *Client) SetGroupReadMarkContext(ctx context.Context, group, ts string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
		"ts":      ***REMOVED***ts***REMOVED***,
	***REMOVED***

	if _, err = groupRequest(ctx, api.httpclient, "groups.mark", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// OpenGroup opens a private group
func (api *Client) OpenGroup(group string) (bool, bool, error) ***REMOVED***
	return api.OpenGroupContext(context.Background(), group)
***REMOVED***

// OpenGroupContext opens a private group with a custom context
func (api *Client) OpenGroupContext(ctx context.Context, group string) (bool, bool, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.open", values, api.debug)
	if err != nil ***REMOVED***
		return false, false, err
	***REMOVED***
	return response.NoOp, response.AlreadyOpen, nil
***REMOVED***

// RenameGroup renames a group
// XXX: They return a channel, not a group. What is this crap? :(
// Inconsistent api it seems.
func (api *Client) RenameGroup(group, name string) (*Channel, error) ***REMOVED***
	return api.RenameGroupContext(context.Background(), group, name)
***REMOVED***

// RenameGroupContext renames a group with a custom context
func (api *Client) RenameGroupContext(ctx context.Context, group, name string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
		"name":    ***REMOVED***name***REMOVED***,
	***REMOVED***

	// XXX: the created entry in this call returns a string instead of a number
	// so I may have to do some workaround to solve it.
	response, err := groupRequest(ctx, api.httpclient, "groups.rename", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// SetGroupPurpose sets the group purpose
func (api *Client) SetGroupPurpose(group, purpose string) (string, error) ***REMOVED***
	return api.SetGroupPurposeContext(context.Background(), group, purpose)
***REMOVED***

// SetGroupPurposeContext sets the group purpose with a custom context
func (api *Client) SetGroupPurposeContext(ctx context.Context, group, purpose string) (string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
		"purpose": ***REMOVED***purpose***REMOVED***,
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.setPurpose", values, api.debug)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return response.Purpose, nil
***REMOVED***

// SetGroupTopic sets the group topic
func (api *Client) SetGroupTopic(group, topic string) (string, error) ***REMOVED***
	return api.SetGroupTopicContext(context.Background(), group, topic)
***REMOVED***

// SetGroupTopicContext sets the group topic with a custom context
func (api *Client) SetGroupTopicContext(ctx context.Context, group, topic string) (string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***group***REMOVED***,
		"topic":   ***REMOVED***topic***REMOVED***,
	***REMOVED***

	response, err := groupRequest(ctx, api.httpclient, "groups.setTopic", values, api.debug)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return response.Topic, nil
***REMOVED***
