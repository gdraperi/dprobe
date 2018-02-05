package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// Conversation is the foundation for IM and BaseGroupConversation
type conversation struct ***REMOVED***
	ID                 string   `json:"id"`
	Created            JSONTime `json:"created"`
	IsOpen             bool     `json:"is_open"`
	LastRead           string   `json:"last_read,omitempty"`
	Latest             *Message `json:"latest,omitempty"`
	UnreadCount        int      `json:"unread_count,omitempty"`
	UnreadCountDisplay int      `json:"unread_count_display,omitempty"`
	IsGroup            bool     `json:"is_group"`
	IsShared           bool     `json:"is_shared"`
	IsIM               bool     `json:"is_im"`
	IsExtShared        bool     `json:"is_ext_shared"`
	IsOrgShared        bool     `json:"is_org_shared"`
	IsPendingExtShared bool     `json:"is_pending_ext_shared"`
	IsPrivate          bool     `json:"is_private"`
	IsMpIM             bool     `json:"is_mpim"`
	Unlinked           int      `json:"unlinked"`
	NameNormalized     string   `json:"name_normalized"`
	NumMembers         int      `json:"num_members"`
	Priority           float64  `json:"priority"`
	// TODO support pending_shared
	// TODO support previous_names
***REMOVED***

// GroupConversation is the foundation for Group and Channel
type groupConversation struct ***REMOVED***
	conversation
	Name       string   `json:"name"`
	Creator    string   `json:"creator"`
	IsArchived bool     `json:"is_archived"`
	Members    []string `json:"members"`
	Topic      Topic    `json:"topic"`
	Purpose    Purpose  `json:"purpose"`
***REMOVED***

// Topic contains information about the topic
type Topic struct ***REMOVED***
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
***REMOVED***

// Purpose contains information about the purpose
type Purpose struct ***REMOVED***
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
***REMOVED***

type GetUsersInConversationParameters struct ***REMOVED***
	ChannelID string
	Cursor    string
	Limit     int
***REMOVED***

type responseMetaData struct ***REMOVED***
	NextCursor string `json:"next_cursor"`
***REMOVED***

// GetUsersInConversation returns the list of users in a conversation
func (api *Client) GetUsersInConversation(params *GetUsersInConversationParameters) ([]string, string, error) ***REMOVED***
	return api.GetUsersInConversationContext(context.Background(), params)
***REMOVED***

// GetUsersInConversationContext returns the list of users in a conversation with a custom context
func (api *Client) GetUsersInConversationContext(ctx context.Context, params *GetUsersInConversationParameters) ([]string, string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***params.ChannelID***REMOVED***,
	***REMOVED***
	if params.Cursor != "" ***REMOVED***
		values.Add("cursor", params.Cursor)
	***REMOVED***
	if params.Limit != 0 ***REMOVED***
		values.Add("limit", string(params.Limit))
	***REMOVED***
	response := struct ***REMOVED***
		Members          []string         `json:"members"`
		ResponseMetaData responseMetaData `json:"response_metadata"`
		SlackResponse
	***REMOVED******REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.members", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, "", errors.New(response.Error)
	***REMOVED***
	return response.Members, response.ResponseMetaData.NextCursor, nil
***REMOVED***

// ArchiveConversation archives a conversation
func (api *Client) ArchiveConversation(channelID string) error ***REMOVED***
	return api.ArchiveConversationContext(context.Background(), channelID)
***REMOVED***

// ArchiveConversationContext archives a conversation with a custom context
func (api *Client) ArchiveConversationContext(ctx context.Context, channelID string) error ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***
	response := SlackResponse***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.archive", values, &response, api.debug)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// UnArchiveConversation reverses conversation archival
func (api *Client) UnArchiveConversation(channelID string) error ***REMOVED***
	return api.UnArchiveConversationContext(context.Background(), channelID)
***REMOVED***

// UnArchiveConversationContext reverses conversation archival with a custom context
func (api *Client) UnArchiveConversationContext(ctx context.Context, channelID string) error ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***
	response := SlackResponse***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.unarchive", values, &response, api.debug)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// SetTopicOfConversation sets the topic for a conversation
func (api *Client) SetTopicOfConversation(channelID, topic string) (*Channel, error) ***REMOVED***
	return api.SetTopicOfConversationContext(context.Background(), channelID, topic)
***REMOVED***

// SetTopicOfConversationContext sets the topic for a conversation with a custom context
func (api *Client) SetTopicOfConversationContext(ctx context.Context, channelID, topic string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"topic":   ***REMOVED***topic***REMOVED***,
	***REMOVED***
	response := struct ***REMOVED***
		SlackResponse
		Channel *Channel `json:"channel"`
	***REMOVED******REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.setTopic", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response.Channel, nil
***REMOVED***

// SetPurposeOfConversation sets the purpose for a conversation
func (api *Client) SetPurposeOfConversation(channelID, purpose string) (*Channel, error) ***REMOVED***
	return api.SetPurposeOfConversationContext(context.Background(), channelID, purpose)
***REMOVED***

// SetPurposeOfConversationContext sets the purpose for a conversation with a custom context
func (api *Client) SetPurposeOfConversationContext(ctx context.Context, channelID, purpose string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"purpose": ***REMOVED***purpose***REMOVED***,
	***REMOVED***
	response := struct ***REMOVED***
		SlackResponse
		Channel *Channel `json:"channel"`
	***REMOVED******REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.setPurpose", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response.Channel, nil
***REMOVED***

// RenameConversation renames a conversation
func (api *Client) RenameConversation(channelID, channelName string) (*Channel, error) ***REMOVED***
	return api.RenameConversationContext(context.Background(), channelID, channelName)
***REMOVED***

// RenameConversationContext renames a conversation with a custom context
func (api *Client) RenameConversationContext(ctx context.Context, channelID, channelName string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"name":    ***REMOVED***channelName***REMOVED***,
	***REMOVED***
	response := struct ***REMOVED***
		SlackResponse
		Channel *Channel `json:"channel"`
	***REMOVED******REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.rename", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response.Channel, nil
***REMOVED***

// InviteUsersToConversation invites users to a channel
func (api *Client) InviteUsersToConversation(channelID string, users ...string) (*Channel, error) ***REMOVED***
	return api.InviteUsersToConversationContext(context.Background(), channelID, users...)
***REMOVED***

// InviteUsersToConversationContext invites users to a channel with a custom context
func (api *Client) InviteUsersToConversationContext(ctx context.Context, channelID string, users ...string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"users":   ***REMOVED***strings.Join(users, ",")***REMOVED***,
	***REMOVED***
	response := struct ***REMOVED***
		SlackResponse
		Channel *Channel `json:"channel"`
	***REMOVED******REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.invite", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response.Channel, nil
***REMOVED***

// KickUserFromConversation removes a user from a conversation
func (api *Client) KickUserFromConversation(channelID string, user string) error ***REMOVED***
	return api.KickUserFromConversationContext(context.Background(), channelID, user)
***REMOVED***

// KickUserFromConversationContext removes a user from a conversation with a custom context
func (api *Client) KickUserFromConversationContext(ctx context.Context, channelID string, user string) error ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"user":    ***REMOVED***user***REMOVED***,
	***REMOVED***
	response := SlackResponse***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.kick", values, &response, api.debug)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return errors.New(response.Error)
	***REMOVED***
	return nil
***REMOVED***

// CloseConversation closes a direct message or multi-person direct message
func (api *Client) CloseConversation(channelID string) (noOp bool, alreadyClosed bool, err error) ***REMOVED***
	return api.CloseConversationContext(context.Background(), channelID)
***REMOVED***

// CloseConversationContext closes a direct message or multi-person direct message with a custom context
func (api *Client) CloseConversationContext(ctx context.Context, channelID string) (noOp bool, alreadyClosed bool, err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***
	response := struct ***REMOVED***
		SlackResponse
		NoOp          bool `json:"no_op"`
		AlreadyClosed bool `json:"already_closed"`
	***REMOVED******REMOVED******REMOVED***

	err = post(ctx, api.httpclient, "conversations.close", values, &response, api.debug)
	if err != nil ***REMOVED***
		return false, false, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return false, false, errors.New(response.Error)
	***REMOVED***
	return response.NoOp, response.AlreadyClosed, nil
***REMOVED***

// CreateConversation initiates a public or private channel-based conversation
func (api *Client) CreateConversation(channelName string, isPrivate bool) (*Channel, error) ***REMOVED***
	return api.CreateConversationContext(context.Background(), channelName, isPrivate)
***REMOVED***

// CreateConversationContext initiates a public or private channel-based conversation with a custom context
func (api *Client) CreateConversationContext(ctx context.Context, channelName string, isPrivate bool) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":      ***REMOVED***api.token***REMOVED***,
		"name":       ***REMOVED***channelName***REMOVED***,
		"is_private": ***REMOVED***strconv.FormatBool(isPrivate)***REMOVED***,
	***REMOVED***
	response, err := channelRequest(
		ctx, api.httpclient, "conversations.create", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// GetConversationInfo retrieves information about a conversation
func (api *Client) GetConversationInfo(channelID string, includeLocale bool) (*Channel, error) ***REMOVED***
	return api.GetConversationInfoContext(context.Background(), channelID, includeLocale)
***REMOVED***

// GetConversationInfoContext retrieves information about a conversation with a custom context
func (api *Client) GetConversationInfoContext(ctx context.Context, channelID string, includeLocale bool) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":          ***REMOVED***api.token***REMOVED***,
		"channel":        ***REMOVED***channelID***REMOVED***,
		"include_locale": ***REMOVED***strconv.FormatBool(includeLocale)***REMOVED***,
	***REMOVED***
	response, err := channelRequest(
		ctx, api.httpclient, "conversations.info", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// LeaveConversation leaves a conversation
func (api *Client) LeaveConversation(channelID string) (bool, error) ***REMOVED***
	return api.LeaveConversationContext(context.Background(), channelID)
***REMOVED***

// LeaveConversationContext leaves a conversation with a custom context
func (api *Client) LeaveConversationContext(ctx context.Context, channelID string) (bool, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "conversations.leave", values, api.debug)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	return response.NotInChannel, nil
***REMOVED***

type GetConversationRepliesParameters struct ***REMOVED***
	ChannelID string
	Timestamp string
	Cursor    string
	Inclusive bool
	Latest    string
	Limit     int
	Oldest    string
***REMOVED***

// GetConversationReplies retrieves a thread of messages posted to a conversation
func (api *Client) GetConversationReplies(params *GetConversationRepliesParameters) (msgs []Message, hasMore bool, nextCursor string, err error) ***REMOVED***
	return api.GetConversationRepliesContext(context.Background(), params)
***REMOVED***

// GetConversationRepliesContext retrieves a thread of messages posted to a conversation with a custom context
func (api *Client) GetConversationRepliesContext(ctx context.Context, params *GetConversationRepliesParameters) (msgs []Message, hasMore bool, nextCursor string, err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***params.ChannelID***REMOVED***,
		"ts":      ***REMOVED***params.Timestamp***REMOVED***,
	***REMOVED***
	if params.Cursor != "" ***REMOVED***
		values.Add("cursor", params.Cursor)
	***REMOVED***
	if params.Latest != "" ***REMOVED***
		values.Add("latest", params.Latest)
	***REMOVED***
	if params.Limit != 0 ***REMOVED***
		values.Add("limit", string(params.Limit))
	***REMOVED***
	if params.Oldest != "" ***REMOVED***
		values.Add("oldest", params.Oldest)
	***REMOVED***
	if params.Inclusive ***REMOVED***
		values.Add("inclusive", "1")
	***REMOVED*** else ***REMOVED***
		values.Add("inclusive", "0")
	***REMOVED***
	response := struct ***REMOVED***
		SlackResponse
		HasMore          bool `json:"has_more"`
		ResponseMetaData struct ***REMOVED***
			NextCursor string `json:"next_cursor"`
		***REMOVED*** `json:"response_metadata"`
		Messages []Message `json:"messages"`
	***REMOVED******REMOVED******REMOVED***

	err = post(ctx, api.httpclient, "conversations.replies", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, false, "", err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, false, "", errors.New(response.Error)
	***REMOVED***
	return response.Messages, response.HasMore, response.ResponseMetaData.NextCursor, nil
***REMOVED***

type GetConversationsParameters struct ***REMOVED***
	Cursor          string
	ExcludeArchived string
	Limit           int
	Types           []string
***REMOVED***

// GetConversations returns the list of channels in a Slack team
func (api *Client) GetConversations(params *GetConversationsParameters) (channels []Channel, nextCursor string, err error) ***REMOVED***
	return api.GetConversationsContext(context.Background(), params)
***REMOVED***

// GetConversationsContext returns the list of channels in a Slack team with a custom context
func (api *Client) GetConversationsContext(ctx context.Context, params *GetConversationsParameters) (channels []Channel, nextCursor string, err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":            ***REMOVED***api.token***REMOVED***,
		"exclude_archived": ***REMOVED***params.ExcludeArchived***REMOVED***,
	***REMOVED***
	if params.Cursor != "" ***REMOVED***
		values.Add("cursor", params.Cursor)
	***REMOVED***
	if params.Limit != 0 ***REMOVED***
		values.Add("limit", string(params.Limit))
	***REMOVED***
	if params.Types != nil ***REMOVED***
		values.Add("types", strings.Join(params.Types, ","))
	***REMOVED***
	response := struct ***REMOVED***
		Channels         []Channel        `json:"channels"`
		ResponseMetaData responseMetaData `json:"response_metadata"`
		SlackResponse
	***REMOVED******REMOVED******REMOVED***
	err = post(ctx, api.httpclient, "conversations.list", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, "", errors.New(response.Error)
	***REMOVED***
	return response.Channels, response.ResponseMetaData.NextCursor, nil
***REMOVED***

type OpenConversationParameters struct ***REMOVED***
	ChannelID string
	ReturnIM  bool
	Users     []string
***REMOVED***

// OpenConversation opens or resumes a direct message or multi-person direct message
func (api *Client) OpenConversation(params *OpenConversationParameters) (*Channel, bool, bool, error) ***REMOVED***
	return api.OpenConversationContext(context.Background(), params)
***REMOVED***

// OpenConversationContext opens or resumes a direct message or multi-person direct message with a custom context
func (api *Client) OpenConversationContext(ctx context.Context, params *OpenConversationParameters) (*Channel, bool, bool, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":     ***REMOVED***api.token***REMOVED***,
		"return_im": ***REMOVED***strconv.FormatBool(params.ReturnIM)***REMOVED***,
	***REMOVED***
	if params.ChannelID != "" ***REMOVED***
		values.Add("channel", params.ChannelID)
	***REMOVED***
	if params.Users != nil ***REMOVED***
		values.Add("users", strings.Join(params.Users, ","))
	***REMOVED***
	response := struct ***REMOVED***
		Channel     *Channel `json:"channel"`
		NoOp        bool     `json:"no_op"`
		AlreadyOpen bool     `json:"already_open"`
		SlackResponse
	***REMOVED******REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.open", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, false, false, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, false, false, errors.New(response.Error)
	***REMOVED***
	return response.Channel, response.NoOp, response.AlreadyOpen, nil
***REMOVED***

// JoinConversation joins an existing conversation
func (api *Client) JoinConversation(channelID string) (*Channel, string, []string, error) ***REMOVED***
	return api.JoinConversationContext(context.Background(), channelID)
***REMOVED***

// JoinConversationContext joins an existing conversation with a custom context
func (api *Client) JoinConversationContext(ctx context.Context, channelID string) (*Channel, string, []string, error) ***REMOVED***
	values := url.Values***REMOVED***"token": ***REMOVED***api.token***REMOVED***, "channel": ***REMOVED***channelID***REMOVED******REMOVED***
	response := struct ***REMOVED***
		Channel          *Channel `json:"channel"`
		Warning          string   `json:"warning"`
		ResponseMetaData *struct ***REMOVED***
			Warnings []string `json:"warnings"`
		***REMOVED*** `json:"response_metadata"`
		SlackResponse
	***REMOVED******REMOVED******REMOVED***
	err := post(ctx, api.httpclient, "conversations.join", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, "", nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, "", nil, errors.New(response.Error)
	***REMOVED***
	var warnings []string
	if response.ResponseMetaData != nil ***REMOVED***
		warnings = response.ResponseMetaData.Warnings
	***REMOVED***
	return response.Channel, response.Warning, warnings, nil
***REMOVED***

type GetConversationHistoryParameters struct ***REMOVED***
	ChannelID string
	Cursor    string
	Inclusive bool
	Latest    string
	Limit     int
	Oldest    string
***REMOVED***

type GetConversationHistoryResponse struct ***REMOVED***
	SlackResponse
	HasMore          bool   `json:"has_more"`
	PinCount         int    `json:"pin_count"`
	Latest           string `json:"latest"`
	ResponseMetaData struct ***REMOVED***
		NextCursor string `json:"next_cursor"`
	***REMOVED*** `json:"response_metadata"`
	Messages []Message `json:"messages"`
***REMOVED***

// GetConversationHistory joins an existing conversation
func (api *Client) GetConversationHistory(params *GetConversationHistoryParameters) (*GetConversationHistoryResponse, error) ***REMOVED***
	return api.GetConversationHistoryContext(context.Background(), params)
***REMOVED***

// GetConversationHistoryContext joins an existing conversation with a custom context
func (api *Client) GetConversationHistoryContext(ctx context.Context, params *GetConversationHistoryParameters) (*GetConversationHistoryResponse, error) ***REMOVED***
	values := url.Values***REMOVED***"token": ***REMOVED***api.token***REMOVED***, "channel": ***REMOVED***params.ChannelID***REMOVED******REMOVED***
	if params.Cursor != "" ***REMOVED***
		values.Add("cursor", params.Cursor)
	***REMOVED***
	if params.Inclusive ***REMOVED***
		values.Add("inclusive", "1")
	***REMOVED*** else ***REMOVED***
		values.Add("inclusive", "0")
	***REMOVED***
	if params.Latest != "" ***REMOVED***
		values.Add("latest", params.Latest)
	***REMOVED***
	if params.Limit != 0 ***REMOVED***
		values.Add("limit", string(params.Limit))
	***REMOVED***
	if params.Oldest != "" ***REMOVED***
		values.Add("oldest", params.Oldest)
	***REMOVED***

	response := GetConversationHistoryResponse***REMOVED******REMOVED***

	err := post(ctx, api.httpclient, "conversations.history", values, &response, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return &response, nil
***REMOVED***
