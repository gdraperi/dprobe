package slack

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

const (
	DEFAULT_MESSAGE_USERNAME         = ""
	DEFAULT_MESSAGE_REPLY_BROADCAST  = false
	DEFAULT_MESSAGE_ASUSER           = false
	DEFAULT_MESSAGE_PARSE            = ""
	DEFAULT_MESSAGE_THREAD_TIMESTAMP = ""
	DEFAULT_MESSAGE_LINK_NAMES       = 0
	DEFAULT_MESSAGE_UNFURL_LINKS     = false
	DEFAULT_MESSAGE_UNFURL_MEDIA     = true
	DEFAULT_MESSAGE_ICON_URL         = ""
	DEFAULT_MESSAGE_ICON_EMOJI       = ""
	DEFAULT_MESSAGE_MARKDOWN         = true
	DEFAULT_MESSAGE_ESCAPE_TEXT      = true
)

type chatResponseFull struct ***REMOVED***
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
	SlackResponse
***REMOVED***

// PostMessageParameters contains all the parameters necessary (including the optional ones) for a PostMessage() request
type PostMessageParameters struct ***REMOVED***
	Username        string       `json:"user_name"`
	AsUser          bool         `json:"as_user"`
	Parse           string       `json:"parse"`
	ThreadTimestamp string       `json:"thread_ts"`
	ReplyBroadcast  bool         `json:"reply_broadcast"`
	LinkNames       int          `json:"link_names"`
	Attachments     []Attachment `json:"attachments"`
	UnfurlLinks     bool         `json:"unfurl_links"`
	UnfurlMedia     bool         `json:"unfurl_media"`
	IconURL         string       `json:"icon_url"`
	IconEmoji       string       `json:"icon_emoji"`
	Markdown        bool         `json:"mrkdwn,omitempty"`
	EscapeText      bool         `json:"escape_text"`

	// chat.postEphemeral support
	Channel string `json:"channel"`
	User    string `json:"user"`
***REMOVED***

// NewPostMessageParameters provides an instance of PostMessageParameters with all the sane default values set
func NewPostMessageParameters() PostMessageParameters ***REMOVED***
	return PostMessageParameters***REMOVED***
		Username:        DEFAULT_MESSAGE_USERNAME,
		User:            DEFAULT_MESSAGE_USERNAME,
		AsUser:          DEFAULT_MESSAGE_ASUSER,
		Parse:           DEFAULT_MESSAGE_PARSE,
		ThreadTimestamp: DEFAULT_MESSAGE_THREAD_TIMESTAMP,
		LinkNames:       DEFAULT_MESSAGE_LINK_NAMES,
		Attachments:     nil,
		UnfurlLinks:     DEFAULT_MESSAGE_UNFURL_LINKS,
		UnfurlMedia:     DEFAULT_MESSAGE_UNFURL_MEDIA,
		IconURL:         DEFAULT_MESSAGE_ICON_URL,
		IconEmoji:       DEFAULT_MESSAGE_ICON_EMOJI,
		Markdown:        DEFAULT_MESSAGE_MARKDOWN,
		EscapeText:      DEFAULT_MESSAGE_ESCAPE_TEXT,
	***REMOVED***
***REMOVED***

// DeleteMessage deletes a message in a channel
func (api *Client) DeleteMessage(channel, messageTimestamp string) (string, string, error) ***REMOVED***
	respChannel, respTimestamp, _, err := api.SendMessageContext(context.Background(), channel, MsgOptionDelete(messageTimestamp))
	return respChannel, respTimestamp, err
***REMOVED***

// DeleteMessageContext deletes a message in a channel with a custom context
func (api *Client) DeleteMessageContext(ctx context.Context, channel, messageTimestamp string) (string, string, error) ***REMOVED***
	respChannel, respTimestamp, _, err := api.SendMessageContext(ctx, channel, MsgOptionDelete(messageTimestamp))
	return respChannel, respTimestamp, err
***REMOVED***

// PostMessage sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (api *Client) PostMessage(channel, text string, params PostMessageParameters) (string, string, error) ***REMOVED***
	respChannel, respTimestamp, _, err := api.SendMessageContext(
		context.Background(),
		channel,
		MsgOptionText(text, params.EscapeText),
		MsgOptionAttachments(params.Attachments...),
		MsgOptionPostMessageParameters(params),
	)
	return respChannel, respTimestamp, err
***REMOVED***

// PostMessageContext sends a message to a channel with a custom context
// For more details, see PostMessage documentation
func (api *Client) PostMessageContext(ctx context.Context, channel, text string, params PostMessageParameters) (string, string, error) ***REMOVED***
	respChannel, respTimestamp, _, err := api.SendMessageContext(
		ctx,
		channel,
		MsgOptionText(text, params.EscapeText),
		MsgOptionAttachments(params.Attachments...),
		MsgOptionPostMessageParameters(params),
	)
	return respChannel, respTimestamp, err
***REMOVED***

// PostEphemeral sends an ephemeral message to a user in a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (api *Client) PostEphemeral(channel, userID string, options ...MsgOption) (string, error) ***REMOVED***
	options = append(options, MsgOptionPostEphemeral())
	return api.PostEphemeralContext(
		context.Background(),
		channel,
		userID,
		options...,
	)
***REMOVED***

// PostEphemeralContext sends an ephemeal message to a user in a channel with a custom context
// For more details, see PostEphemeral documentation
func (api *Client) PostEphemeralContext(ctx context.Context, channel, userID string, options ...MsgOption) (string, error) ***REMOVED***
	path, values, err := ApplyMsgOptions(api.token, channel, options...)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	values.Add("user", userID)

	response, err := chatRequest(ctx, api.httpclient, path, values, api.debug)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return response.Timestamp, nil
***REMOVED***

// UpdateMessage updates a message in a channel
func (api *Client) UpdateMessage(channel, timestamp, text string) (string, string, string, error) ***REMOVED***
	return api.UpdateMessageContext(context.Background(), channel, timestamp, text)
***REMOVED***

// UpdateMessageContext updates a message in a channel
func (api *Client) UpdateMessageContext(ctx context.Context, channel, timestamp, text string) (string, string, string, error) ***REMOVED***
	return api.SendMessageContext(ctx, channel, MsgOptionUpdate(timestamp), MsgOptionText(text, true))
***REMOVED***

// SendMessage more flexible method for configuring messages.
func (api *Client) SendMessage(channel string, options ...MsgOption) (string, string, string, error) ***REMOVED***
	return api.SendMessageContext(context.Background(), channel, options...)
***REMOVED***

// SendMessageContext more flexible method for configuring messages with a custom context.
func (api *Client) SendMessageContext(ctx context.Context, channel string, options ...MsgOption) (string, string, string, error) ***REMOVED***
	channel, values, err := ApplyMsgOptions(api.token, channel, options...)
	if err != nil ***REMOVED***
		return "", "", "", err
	***REMOVED***

	response, err := chatRequest(ctx, api.httpclient, channel, values, api.debug)
	if err != nil ***REMOVED***
		return "", "", "", err
	***REMOVED***

	return response.Channel, response.Timestamp, response.Text, nil
***REMOVED***

// ApplyMsgOptions utility function for debugging/testing chat requests.
func ApplyMsgOptions(token, channel string, options ...MsgOption) (string, url.Values, error) ***REMOVED***
	config := sendConfig***REMOVED***
		mode: chatPostMessage,
		values: url.Values***REMOVED***
			"token":   ***REMOVED***token***REMOVED***,
			"channel": ***REMOVED***channel***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, opt := range options ***REMOVED***
		if err := opt(&config); err != nil ***REMOVED***
			return string(config.mode), config.values, err
		***REMOVED***
	***REMOVED***

	return string(config.mode), config.values, nil
***REMOVED***

func escapeMessage(message string) string ***REMOVED***
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(message)
***REMOVED***

func chatRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*chatResponseFull, error) ***REMOVED***
	response := &chatResponseFull***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

type sendMode string

const (
	chatUpdate        sendMode = "chat.update"
	chatPostMessage   sendMode = "chat.postMessage"
	chatDelete        sendMode = "chat.delete"
	chatPostEphemeral sendMode = "chat.postEphemeral"
)

type sendConfig struct ***REMOVED***
	mode   sendMode
	values url.Values
***REMOVED***

// MsgOption option provided when sending a message.
type MsgOption func(*sendConfig) error

// MsgOptionPost posts a messages, this is the default.
func MsgOptionPost() MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.mode = chatPostMessage
		config.values.Del("ts")
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionPostEphemeral posts an ephemeral message
func MsgOptionPostEphemeral() MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.mode = chatPostEphemeral
		config.values.Del("ts")
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionUpdate updates a message based on the timestamp.
func MsgOptionUpdate(timestamp string) MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.mode = chatUpdate
		config.values.Add("ts", timestamp)
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionDelete deletes a message based on the timestamp.
func MsgOptionDelete(timestamp string) MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.mode = chatDelete
		config.values.Add("ts", timestamp)
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionAsUser whether or not to send the message as the user.
func MsgOptionAsUser(b bool) MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		if b != DEFAULT_MESSAGE_ASUSER ***REMOVED***
			config.values.Set("as_user", "true")
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionText provide the text for the message, optionally escape the provided
// text.
func MsgOptionText(text string, escape bool) MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		if escape ***REMOVED***
			text = escapeMessage(text)
		***REMOVED***
		config.values.Add("text", text)
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionAttachments provide attachments for the message.
func MsgOptionAttachments(attachments ...Attachment) MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		if attachments == nil ***REMOVED***
			return nil
		***REMOVED***

		attachments, err := json.Marshal(attachments)
		if err == nil ***REMOVED***
			config.values.Set("attachments", string(attachments))
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

// MsgOptionEnableLinkUnfurl enables link unfurling
func MsgOptionEnableLinkUnfurl() MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.values.Set("unfurl_links", "true")
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionDisableLinkUnfurl disables link unfurling
func MsgOptionDisableLinkUnfurl() MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.values.Set("unfurl_links", "false")
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionDisableMediaUnfurl disables media unfurling.
func MsgOptionDisableMediaUnfurl() MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.values.Set("unfurl_media", "false")
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionDisableMarkdown disables markdown.
func MsgOptionDisableMarkdown() MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		config.values.Set("mrkdwn", "false")
		return nil
	***REMOVED***
***REMOVED***

// MsgOptionPostMessageParameters maintain backwards compatibility.
func MsgOptionPostMessageParameters(params PostMessageParameters) MsgOption ***REMOVED***
	return func(config *sendConfig) error ***REMOVED***
		if params.Username != DEFAULT_MESSAGE_USERNAME ***REMOVED***
			config.values.Set("username", string(params.Username))
		***REMOVED***

		// chat.postEphemeral support
		if params.User != DEFAULT_MESSAGE_USERNAME ***REMOVED***
			config.values.Set("user", params.User)
		***REMOVED***

		// never generates an error.
		MsgOptionAsUser(params.AsUser)(config)

		if params.Parse != DEFAULT_MESSAGE_PARSE ***REMOVED***
			config.values.Set("parse", string(params.Parse))
		***REMOVED***
		if params.LinkNames != DEFAULT_MESSAGE_LINK_NAMES ***REMOVED***
			config.values.Set("link_names", "1")
		***REMOVED***

		if params.UnfurlLinks != DEFAULT_MESSAGE_UNFURL_LINKS ***REMOVED***
			config.values.Set("unfurl_links", "true")
		***REMOVED***

		// I want to send a message with explicit `as_user` `true` and `unfurl_links` `false` in request.
		// Because setting `as_user` to `true` will change the default value for `unfurl_links` to `true` on Slack API side.
		if params.AsUser != DEFAULT_MESSAGE_ASUSER && params.UnfurlLinks == DEFAULT_MESSAGE_UNFURL_LINKS ***REMOVED***
			config.values.Set("unfurl_links", "false")
		***REMOVED***
		if params.UnfurlMedia != DEFAULT_MESSAGE_UNFURL_MEDIA ***REMOVED***
			config.values.Set("unfurl_media", "false")
		***REMOVED***
		if params.IconURL != DEFAULT_MESSAGE_ICON_URL ***REMOVED***
			config.values.Set("icon_url", params.IconURL)
		***REMOVED***
		if params.IconEmoji != DEFAULT_MESSAGE_ICON_EMOJI ***REMOVED***
			config.values.Set("icon_emoji", params.IconEmoji)
		***REMOVED***
		if params.Markdown != DEFAULT_MESSAGE_MARKDOWN ***REMOVED***
			config.values.Set("mrkdwn", "false")
		***REMOVED***

		if params.ThreadTimestamp != DEFAULT_MESSAGE_THREAD_TIMESTAMP ***REMOVED***
			config.values.Set("thread_ts", params.ThreadTimestamp)
		***REMOVED***
		if params.ReplyBroadcast != DEFAULT_MESSAGE_REPLY_BROADCAST ***REMOVED***
			config.values.Set("reply_broadcast", "true")
		***REMOVED***

		return nil
	***REMOVED***
***REMOVED***
