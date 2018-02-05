package slack

import (
	"encoding/json"
	"fmt"
)

// AckMessage is used for messages received in reply to other messages
type AckMessage struct ***REMOVED***
	ReplyTo   int    `json:"reply_to"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
	RTMResponse
***REMOVED***

// RTMResponse encapsulates response details as returned by the Slack API
type RTMResponse struct ***REMOVED***
	Ok    bool      `json:"ok"`
	Error *RTMError `json:"error"`
***REMOVED***

// RTMError encapsulates error information as returned by the Slack API
type RTMError struct ***REMOVED***
	Code int
	Msg  string
***REMOVED***

func (s RTMError) Error() string ***REMOVED***
	return fmt.Sprintf("Code %d - %s", s.Code, s.Msg)
***REMOVED***

// MessageEvent represents a Slack Message (used as the event type for an incoming message)
type MessageEvent Message

// RTMEvent is the main wrapper. You will find all the other messages attached
type RTMEvent struct ***REMOVED***
	Type string
	Data interface***REMOVED******REMOVED***
***REMOVED***

// HelloEvent represents the hello event
type HelloEvent struct***REMOVED******REMOVED***

// PresenceChangeEvent represents the presence change event
type PresenceChangeEvent struct ***REMOVED***
	Type     string `json:"type"`
	Presence string `json:"presence"`
	User     string `json:"user"`
***REMOVED***

// UserTypingEvent represents the user typing event
type UserTypingEvent struct ***REMOVED***
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel string `json:"channel"`
***REMOVED***

// PrefChangeEvent represents a user preferences change event
type PrefChangeEvent struct ***REMOVED***
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
***REMOVED***

// ManualPresenceChangeEvent represents the manual presence change event
type ManualPresenceChangeEvent struct ***REMOVED***
	Type     string `json:"type"`
	Presence string `json:"presence"`
***REMOVED***

// UserChangeEvent represents the user change event
type UserChangeEvent struct ***REMOVED***
	Type string `json:"type"`
	User User   `json:"user"`
***REMOVED***

// EmojiChangedEvent represents the emoji changed event
type EmojiChangedEvent struct ***REMOVED***
	Type           string   `json:"type"`
	SubType        string   `json:"subtype"`
	Name           string   `json:"name"`
	Names          []string `json:"names"`
	Value          string   `json:"value"` 
	EventTimestamp string   `json:"event_ts"`
***REMOVED***

// CommandsChangedEvent represents the commands changed event
type CommandsChangedEvent struct ***REMOVED***
	Type           string `json:"type"`
	EventTimestamp string `json:"event_ts"`
***REMOVED***

// EmailDomainChangedEvent represents the email domain changed event
type EmailDomainChangedEvent struct ***REMOVED***
	Type           string `json:"type"`
	EventTimestamp string `json:"event_ts"`
	EmailDomain    string `json:"email_domain"`
***REMOVED***

// BotAddedEvent represents the bot added event
type BotAddedEvent struct ***REMOVED***
	Type string `json:"type"`
	Bot  Bot    `json:"bot"`
***REMOVED***

// BotChangedEvent represents the bot changed event
type BotChangedEvent struct ***REMOVED***
	Type string `json:"type"`
	Bot  Bot    `json:"bot"`
***REMOVED***

// AccountsChangedEvent represents the accounts changed event
type AccountsChangedEvent struct ***REMOVED***
	Type string `json:"type"`
***REMOVED***

// ReconnectUrlEvent represents the receiving reconnect url event
type ReconnectUrlEvent struct ***REMOVED***
	Type string `json:"type"`
	URL  string `json:"url"`
***REMOVED***
