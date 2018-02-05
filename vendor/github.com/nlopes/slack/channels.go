package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

type channelResponseFull struct ***REMOVED***
	Channel      Channel   `json:"channel"`
	Channels     []Channel `json:"channels"`
	Purpose      string    `json:"purpose"`
	Topic        string    `json:"topic"`
	NotInChannel bool      `json:"not_in_channel"`
	History
	SlackResponse
***REMOVED***

// Channel contains information about the channel
type Channel struct ***REMOVED***
	groupConversation
	IsChannel bool   `json:"is_channel"`
	IsGeneral bool   `json:"is_general"`
	IsMember  bool   `json:"is_member"`
	Locale    string `json:"locale"`
***REMOVED***

func channelRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*channelResponseFull, error) ***REMOVED***
	response := &channelResponseFull***REMOVED******REMOVED***
	err := postForm(ctx, client, SLACK_API+path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// ArchiveChannel archives the given channel
// see https://api.slack.com/methods/channels.archive
func (api *Client) ArchiveChannel(channelID string) error ***REMOVED***
	return api.ArchiveChannelContext(context.Background(), channelID)
***REMOVED***

// ArchiveChannelContext archives the given channel with a custom context
// see https://api.slack.com/methods/channels.archive
func (api *Client) ArchiveChannelContext(ctx context.Context, channelID string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***

	if _, err = channelRequest(ctx, api.httpclient, "channels.archive", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// UnarchiveChannel unarchives the given channel
// see https://api.slack.com/methods/channels.unarchive
func (api *Client) UnarchiveChannel(channelID string) error ***REMOVED***
	return api.UnarchiveChannelContext(context.Background(), channelID)
***REMOVED***

// UnarchiveChannelContext unarchives the given channel with a custom context
// see https://api.slack.com/methods/channels.unarchive
func (api *Client) UnarchiveChannelContext(ctx context.Context, channelID string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***

	if _, err = channelRequest(ctx, api.httpclient, "channels.unarchive", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// CreateChannel creates a channel with the given name and returns a *Channel
// see https://api.slack.com/methods/channels.create
func (api *Client) CreateChannel(channelName string) (*Channel, error) ***REMOVED***
	return api.CreateChannelContext(context.Background(), channelName)
***REMOVED***

// CreateChannelContext creates a channel with the given name and returns a *Channel with a custom context
// see https://api.slack.com/methods/channels.create
func (api *Client) CreateChannelContext(ctx context.Context, channelName string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"name":  ***REMOVED***channelName***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.create", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// GetChannelHistory retrieves the channel history
// see https://api.slack.com/methods/channels.history
func (api *Client) GetChannelHistory(channelID string, params HistoryParameters) (*History, error) ***REMOVED***
	return api.GetChannelHistoryContext(context.Background(), channelID, params)
***REMOVED***

// GetChannelHistoryContext retrieves the channel history with a custom context
// see https://api.slack.com/methods/channels.history
func (api *Client) GetChannelHistoryContext(ctx context.Context, channelID string, params HistoryParameters) (*History, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
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

	response, err := channelRequest(ctx, api.httpclient, "channels.history", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.History, nil
***REMOVED***

// GetChannelInfo retrieves the given channel
// see https://api.slack.com/methods/channels.info
func (api *Client) GetChannelInfo(channelID string) (*Channel, error) ***REMOVED***
	return api.GetChannelInfoContext(context.Background(), channelID)
***REMOVED***

// GetChannelInfoContext retrieves the given channel with a custom context
// see https://api.slack.com/methods/channels.info
func (api *Client) GetChannelInfoContext(ctx context.Context, channelID string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.info", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// InviteUserToChannel invites a user to a given channel and returns a *Channel
// see https://api.slack.com/methods/channels.invite
func (api *Client) InviteUserToChannel(channelID, user string) (*Channel, error) ***REMOVED***
	return api.InviteUserToChannelContext(context.Background(), channelID, user)
***REMOVED***

// InviteUserToChannelCustom invites a user to a given channel and returns a *Channel with a custom context
// see https://api.slack.com/methods/channels.invite
func (api *Client) InviteUserToChannelContext(ctx context.Context, channelID, user string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"user":    ***REMOVED***user***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.invite", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// JoinChannel joins the currently authenticated user to a channel
// see https://api.slack.com/methods/channels.join
func (api *Client) JoinChannel(channelName string) (*Channel, error) ***REMOVED***
	return api.JoinChannelContext(context.Background(), channelName)
***REMOVED***

// JoinChannelContext joins the currently authenticated user to a channel with a custom context
// see https://api.slack.com/methods/channels.join
func (api *Client) JoinChannelContext(ctx context.Context, channelName string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"name":  ***REMOVED***channelName***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.join", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// LeaveChannel makes the authenticated user leave the given channel
// see https://api.slack.com/methods/channels.leave
func (api *Client) LeaveChannel(channelID string) (bool, error) ***REMOVED***
	return api.LeaveChannelContext(context.Background(), channelID)
***REMOVED***

// LeaveChannelContext makes the authenticated user leave the given channel with a custom context
// see https://api.slack.com/methods/channels.leave
func (api *Client) LeaveChannelContext(ctx context.Context, channelID string) (bool, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.leave", values, api.debug)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	return response.NotInChannel, nil
***REMOVED***

// KickUserFromChannel kicks a user from a given channel
// see https://api.slack.com/methods/channels.kick
func (api *Client) KickUserFromChannel(channelID, user string) error ***REMOVED***
	return api.KickUserFromChannelContext(context.Background(), channelID, user)
***REMOVED***

// KickUserFromChannelContext kicks a user from a given channel with a custom context
// see https://api.slack.com/methods/channels.kick
func (api *Client) KickUserFromChannelContext(ctx context.Context, channelID, user string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"user":    ***REMOVED***user***REMOVED***,
	***REMOVED***

	if _, err = channelRequest(ctx, api.httpclient, "channels.kick", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// GetChannels retrieves all the channels
// see https://api.slack.com/methods/channels.list
func (api *Client) GetChannels(excludeArchived bool) ([]Channel, error) ***REMOVED***
	return api.GetChannelsContext(context.Background(), excludeArchived)
***REMOVED***

// GetChannelsContext retrieves all the channels with a custom context
// see https://api.slack.com/methods/channels.list
func (api *Client) GetChannelsContext(ctx context.Context, excludeArchived bool) ([]Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if excludeArchived ***REMOVED***
		values.Add("exclude_archived", "1")
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.list", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return response.Channels, nil
***REMOVED***

// SetChannelReadMark sets the read mark of a given channel to a specific point
// Clients should try to avoid making this call too often. When needing to mark a read position, a client should set a
// timer before making the call. In this way, any further updates needed during the timeout will not generate extra calls
// (just one per channel). This is useful for when reading scroll-back history, or following a busy live channel. A
// timeout of 5 seconds is a good starting point. Be sure to flush these calls on shutdown/logout.
// see https://api.slack.com/methods/channels.mark
func (api *Client) SetChannelReadMark(channelID, ts string) error ***REMOVED***
	return api.SetChannelReadMarkContext(context.Background(), channelID, ts)
***REMOVED***

// SetChannelReadMarkContext sets the read mark of a given channel to a specific point with a custom context
// For more details see SetChannelReadMark documentation
// see https://api.slack.com/methods/channels.mark
func (api *Client) SetChannelReadMarkContext(ctx context.Context, channelID, ts string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"ts":      ***REMOVED***ts***REMOVED***,
	***REMOVED***

	if _, err = channelRequest(ctx, api.httpclient, "channels.mark", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// RenameChannel renames a given channel
// see https://api.slack.com/methods/channels.rename
func (api *Client) RenameChannel(channelID, name string) (*Channel, error) ***REMOVED***
	return api.RenameChannelContext(context.Background(), channelID, name)
***REMOVED***

// RenameChannelContext renames a given channel with a custom context
// see https://api.slack.com/methods/channels.rename
func (api *Client) RenameChannelContext(ctx context.Context, channelID, name string) (*Channel, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"name":    ***REMOVED***name***REMOVED***,
	***REMOVED***

	// XXX: the created entry in this call returns a string instead of a number
	// so I may have to do some workaround to solve it.
	response, err := channelRequest(ctx, api.httpclient, "channels.rename", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Channel, nil
***REMOVED***

// SetChannelPurpose sets the channel purpose and returns the purpose that was successfully set
// see https://api.slack.com/methods/channels.setPurpose
func (api *Client) SetChannelPurpose(channelID, purpose string) (string, error) ***REMOVED***
	return api.SetChannelPurposeContext(context.Background(), channelID, purpose)
***REMOVED***

// SetChannelPurposeContext sets the channel purpose and returns the purpose that was successfully set with a custom context
// see https://api.slack.com/methods/channels.setPurpose
func (api *Client) SetChannelPurposeContext(ctx context.Context, channelID, purpose string) (string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"purpose": ***REMOVED***purpose***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.setPurpose", values, api.debug)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return response.Purpose, nil
***REMOVED***

// SetChannelTopic sets the channel topic and returns the topic that was successfully set
// see https://api.slack.com/methods/channels.setTopic
func (api *Client) SetChannelTopic(channelID, topic string) (string, error) ***REMOVED***
	return api.SetChannelTopicContext(context.Background(), channelID, topic)
***REMOVED***

// SetChannelTopicContext sets the channel topic and returns the topic that was successfully set with a custom context
// see https://api.slack.com/methods/channels.setTopic
func (api *Client) SetChannelTopicContext(ctx context.Context, channelID, topic string) (string, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":   ***REMOVED***api.token***REMOVED***,
		"channel": ***REMOVED***channelID***REMOVED***,
		"topic":   ***REMOVED***topic***REMOVED***,
	***REMOVED***

	response, err := channelRequest(ctx, api.httpclient, "channels.setTopic", values, api.debug)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return response.Topic, nil
***REMOVED***

// GetChannelReplies gets an entire thread (a message plus all the messages in reply to it).
// see https://api.slack.com/methods/channels.replies
func (api *Client) GetChannelReplies(channelID, thread_ts string) ([]Message, error) ***REMOVED***
	return api.GetChannelRepliesContext(context.Background(), channelID, thread_ts)
***REMOVED***

// GetChannelRepliesContext gets an entire thread (a message plus all the messages in reply to it) with a custom context
// see https://api.slack.com/methods/channels.replies
func (api *Client) GetChannelRepliesContext(ctx context.Context, channelID, thread_ts string) ([]Message, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token":     ***REMOVED***api.token***REMOVED***,
		"channel":   ***REMOVED***channelID***REMOVED***,
		"thread_ts": ***REMOVED***thread_ts***REMOVED***,
	***REMOVED***
	response, err := channelRequest(ctx, api.httpclient, "channels.replies", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return response.History.Messages, nil
***REMOVED***
