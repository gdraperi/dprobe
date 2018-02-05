package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

// ManageConnection can be called on a Slack RTM instance returned by the
// NewRTM method. It will connect to the slack RTM API and handle all incoming
// and outgoing events. If a connection fails then it will attempt to reconnect
// and will notify any listeners through an error event on the IncomingEvents
// channel.
//
// If the connection ends and the disconnect was unintentional then this will
// attempt to reconnect.
//
// This should only be called once per slack API! Otherwise expect undefined
// behavior.
//
// The defined error events are located in websocket_internals.go.
func (rtm *RTM) ManageConnection() ***REMOVED***
	var connectionCount int
	for ***REMOVED***
		connectionCount++
		// start trying to connect
		// the returned err is already passed onto the IncomingEvents channel
		info, conn, err := rtm.connect(connectionCount, rtm.useRTMStart)
		// if err != nil then the connection is sucessful - otherwise it is
		// fatal
		if err != nil ***REMOVED***
			rtm.Debugf("Failed to connect with RTM on try %d: %s", connectionCount, err)
			return
		***REMOVED***
		rtm.info = info
		rtm.IncomingEvents <- RTMEvent***REMOVED***"connected", &ConnectedEvent***REMOVED***
			ConnectionCount: connectionCount,
			Info:            info,
		***REMOVED******REMOVED***

		rtm.conn = conn
		rtm.isConnected = true

		rtm.Debugf("RTM connection succeeded on try %d", connectionCount)

		keepRunning := make(chan bool)
		// we're now connected (or have failed fatally) so we can set up
		// listeners
		go rtm.handleIncomingEvents(keepRunning)

		// this should be a blocking call until the connection has ended
		rtm.handleEvents(keepRunning, 30*time.Second)

		// after being disconnected we need to check if it was intentional
		// if not then we should try to reconnect
		if rtm.wasIntentional ***REMOVED***
			return
		***REMOVED***
		// else continue and run the loop again to connect
	***REMOVED***
***REMOVED***

// connect attempts to connect to the slack websocket API. It handles any
// errors that occur while connecting and will return once a connection
// has been successfully opened.
// If useRTMStart is false then it uses rtm.connect to create the connection,
// otherwise it uses rtm.start.
func (rtm *RTM) connect(connectionCount int, useRTMStart bool) (*Info, *websocket.Conn, error) ***REMOVED***
	// used to provide exponential backoff wait time with jitter before trying
	// to connect to slack again
	boff := &backoff***REMOVED***
		Min:    100 * time.Millisecond,
		Max:    5 * time.Minute,
		Factor: 2,
		Jitter: true,
	***REMOVED***

	for ***REMOVED***
		// send connecting event
		rtm.IncomingEvents <- RTMEvent***REMOVED***"connecting", &ConnectingEvent***REMOVED***
			Attempt:         boff.attempts + 1,
			ConnectionCount: connectionCount,
		***REMOVED******REMOVED***
		// attempt to start the connection
		info, conn, err := rtm.startRTMAndDial(useRTMStart)
		if err == nil ***REMOVED***
			return info, conn, nil
		***REMOVED***
		// check for fatal errors - currently only invalid_auth
		if sErr, ok := err.(*WebError); ok && (sErr.Error() == "invalid_auth" || sErr.Error() == "account_inactive") ***REMOVED***
			rtm.Debugf("Invalid auth when connecting with RTM: %s", err)
			rtm.IncomingEvents <- RTMEvent***REMOVED***"invalid_auth", &InvalidAuthEvent***REMOVED******REMOVED******REMOVED***
			return nil, nil, sErr
		***REMOVED***

		// any other errors are treated as recoverable and we try again after
		// sending the event along the IncomingEvents channel
		rtm.IncomingEvents <- RTMEvent***REMOVED***"connection_error", &ConnectionErrorEvent***REMOVED***
			Attempt:  boff.attempts,
			ErrorObj: err,
		***REMOVED******REMOVED***

		// check if Disconnect() has been invoked.
		select ***REMOVED***
		case _ = <-rtm.disconnected:
			rtm.IncomingEvents <- RTMEvent***REMOVED***"disconnected", &DisconnectedEvent***REMOVED***Intentional: true***REMOVED******REMOVED***
			return nil, nil, fmt.Errorf("disconnect received while trying to connect")
		default:
		***REMOVED***

		// get time we should wait before attempting to connect again
		dur := boff.Duration()
		rtm.Debugf("reconnection %d failed: %s", boff.attempts+1, err)
		rtm.Debugln(" -> reconnecting in", dur)
		time.Sleep(dur)
	***REMOVED***
***REMOVED***

// startRTMAndDial attempts to connect to the slack websocket. If useRTMStart is true,
// then it returns the  full information returned by the "rtm.start" method on the
// slack API. Else it uses the "rtm.connect" method to connect
func (rtm *RTM) startRTMAndDial(useRTMStart bool) (*Info, *websocket.Conn, error) ***REMOVED***
	var info *Info
	var url string
	var err error

	if useRTMStart ***REMOVED***
		rtm.Debugf("Starting RTM")
		info, url, err = rtm.StartRTM()
	***REMOVED*** else ***REMOVED***
		rtm.Debugf("Connecting to RTM")
		info, url, err = rtm.ConnectRTM()
	***REMOVED***
	if err != nil ***REMOVED***
		rtm.Debugf("Failed to start or connect to RTM: %s", err)
		return nil, nil, err
	***REMOVED***

	rtm.Debugf("Dialing to websocket on url %s", url)
	// Only use HTTPS for connections to prevent MITM attacks on the connection.
	upgradeHeader := http.Header***REMOVED******REMOVED***
	upgradeHeader.Add("Origin", "https://api.slack.com")
	conn, _, err := websocket.DefaultDialer.Dial(url, upgradeHeader)
	if err != nil ***REMOVED***
		rtm.Debugf("Failed to dial to the websocket: %s", err)
		return nil, nil, err
	***REMOVED***
	return info, conn, err
***REMOVED***

// killConnection stops the websocket connection and signals to all goroutines
// that they should cease listening to the connection for events.
//
// This should not be called directly! Instead a boolean value (true for
// intentional, false otherwise) should be sent to the killChannel on the RTM.
func (rtm *RTM) killConnection(keepRunning chan bool, intentional bool) error ***REMOVED***
	rtm.Debugln("killing connection")
	if rtm.isConnected ***REMOVED***
		close(keepRunning)
	***REMOVED***
	rtm.isConnected = false
	rtm.wasIntentional = intentional
	err := rtm.conn.Close()
	rtm.IncomingEvents <- RTMEvent***REMOVED***"disconnected", &DisconnectedEvent***REMOVED***intentional***REMOVED******REMOVED***
	return err
***REMOVED***

// handleEvents is a blocking function that handles all events. This sends
// pings when asked to (on rtm.forcePing) and upon every given elapsed
// interval. This also sends outgoing messages that are received from the RTM's
// outgoingMessages channel. This also handles incoming raw events from the RTM
// rawEvents channel.
func (rtm *RTM) handleEvents(keepRunning chan bool, interval time.Duration) ***REMOVED***
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for ***REMOVED***
		select ***REMOVED***
		// catch "stop" signal on channel close
		case intentional := <-rtm.killChannel:
			_ = rtm.killConnection(keepRunning, intentional)
			return
			// send pings on ticker interval
		case <-ticker.C:
			err := rtm.ping()
			if err != nil ***REMOVED***
				_ = rtm.killConnection(keepRunning, false)
				return
			***REMOVED***
		case <-rtm.forcePing:
			err := rtm.ping()
			if err != nil ***REMOVED***
				_ = rtm.killConnection(keepRunning, false)
				return
			***REMOVED***
		// listen for messages that need to be sent
		case msg := <-rtm.outgoingMessages:
			rtm.sendOutgoingMessage(msg)
		// listen for incoming messages that need to be parsed
		case rawEvent := <-rtm.rawEvents:
			rtm.handleRawEvent(rawEvent)
		***REMOVED***
	***REMOVED***
***REMOVED***

// handleIncomingEvents monitors the RTM's opened websocket for any incoming
// events. It pushes the raw events onto the RTM channel rawEvents.
//
// This will stop executing once the RTM's keepRunning channel has been closed
// or has anything sent to it.
func (rtm *RTM) handleIncomingEvents(keepRunning <-chan bool) ***REMOVED***
	for ***REMOVED***
		// non-blocking listen to see if channel is closed
		select ***REMOVED***
		// catch "stop" signal on channel close
		case <-keepRunning:
			return
		default:
			rtm.receiveIncomingEvent()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (rtm *RTM) sendWithDeadline(msg interface***REMOVED******REMOVED***) error ***REMOVED***
	// set a write deadline on the connection
	if err := rtm.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := rtm.conn.WriteJSON(msg); err != nil ***REMOVED***
		return err
	***REMOVED***
	// remove write deadline
	return rtm.conn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
***REMOVED***

// sendOutgoingMessage sends the given OutgoingMessage to the slack websocket.
//
// It does not currently detect if a outgoing message fails due to a disconnect
// and instead lets a future failed 'PING' detect the failed connection.
func (rtm *RTM) sendOutgoingMessage(msg OutgoingMessage) ***REMOVED***
	rtm.Debugln("Sending message:", msg)
	if len(msg.Text) > MaxMessageTextLength ***REMOVED***
		rtm.IncomingEvents <- RTMEvent***REMOVED***"outgoing_error", &MessageTooLongEvent***REMOVED***
			Message:   msg,
			MaxLength: MaxMessageTextLength,
		***REMOVED******REMOVED***
		return
	***REMOVED***

	if err := rtm.sendWithDeadline(msg); err != nil ***REMOVED***
		rtm.IncomingEvents <- RTMEvent***REMOVED***"outgoing_error", &OutgoingErrorEvent***REMOVED***
			Message:  msg,
			ErrorObj: err,
		***REMOVED******REMOVED***
		// TODO force ping?
	***REMOVED***
***REMOVED***

// ping sends a 'PING' message to the RTM's websocket. If the 'PING' message
// fails to send then this returns an error signifying that the connection
// should be considered disconnected.
//
// This does not handle incoming 'PONG' responses but does store the time of
// each successful 'PING' send so latency can be detected upon a 'PONG'
// response.
func (rtm *RTM) ping() error ***REMOVED***
	id := rtm.idGen.Next()
	rtm.Debugln("Sending PING ", id)
	rtm.pings[id] = time.Now()

	msg := &Ping***REMOVED***ID: id, Type: "ping"***REMOVED***

	if err := rtm.sendWithDeadline(msg); err != nil ***REMOVED***
		rtm.Debugf("RTM Error sending 'PING %d': %s", id, err.Error())
		return err
	***REMOVED***
	return nil
***REMOVED***

// receiveIncomingEvent attempts to receive an event from the RTM's websocket.
// This will block until a frame is available from the websocket.
func (rtm *RTM) receiveIncomingEvent() ***REMOVED***
	event := json.RawMessage***REMOVED******REMOVED***
	err := rtm.conn.ReadJSON(&event)
	if err == io.EOF ***REMOVED***
		// EOF's don't seem to signify a failed connection so instead we ignore
		// them here and detect a failed connection upon attempting to send a
		// 'PING' message

		// trigger a 'PING' to detect pontential websocket disconnect
		rtm.forcePing <- true
		return
	***REMOVED*** else if err != nil ***REMOVED***
		rtm.IncomingEvents <- RTMEvent***REMOVED***"incoming_error", &IncomingEventError***REMOVED***
			ErrorObj: err,
		***REMOVED******REMOVED***
		// force a ping here too?
		return
	***REMOVED*** else if len(event) == 0 ***REMOVED***
		rtm.Debugln("Received empty event")
		return
	***REMOVED***
	rtm.Debugln("Incoming Event:", string(event[:]))
	rtm.rawEvents <- event
***REMOVED***

// handleRawEvent takes a raw JSON message received from the slack websocket
// and handles the encoded event.
func (rtm *RTM) handleRawEvent(rawEvent json.RawMessage) ***REMOVED***
	event := &Event***REMOVED******REMOVED***
	err := json.Unmarshal(rawEvent, event)
	if err != nil ***REMOVED***
		rtm.IncomingEvents <- RTMEvent***REMOVED***"unmarshalling_error", &UnmarshallingErrorEvent***REMOVED***err***REMOVED******REMOVED***
		return
	***REMOVED***
	switch event.Type ***REMOVED***
	case "":
		rtm.handleAck(rawEvent)
	case "hello":
		rtm.IncomingEvents <- RTMEvent***REMOVED***"hello", &HelloEvent***REMOVED******REMOVED******REMOVED***
	case "pong":
		rtm.handlePong(rawEvent)
	case "desktop_notification":
		rtm.Debugln("Received desktop notification, ignoring")
	default:
		rtm.handleEvent(event.Type, rawEvent)
	***REMOVED***
***REMOVED***

// handleAck handles an incoming 'ACK' message.
func (rtm *RTM) handleAck(event json.RawMessage) ***REMOVED***
	ack := &AckMessage***REMOVED******REMOVED***
	if err := json.Unmarshal(event, ack); err != nil ***REMOVED***
		rtm.Debugln("RTM Error unmarshalling 'ack' event:", err)
		rtm.Debugln(" -> Erroneous 'ack' event:", string(event))
		return
	***REMOVED***

	if ack.Ok ***REMOVED***
		rtm.IncomingEvents <- RTMEvent***REMOVED***"ack", ack***REMOVED***
	***REMOVED*** else if ack.RTMResponse.Error != nil ***REMOVED***
		// As there is no documentation for RTM error-codes, this
		// identification of a rate-limit warning is very brittle.
		if ack.RTMResponse.Error.Code == -1 && ack.RTMResponse.Error.Msg == "slow down, too many messages..." ***REMOVED***
			rtm.IncomingEvents <- RTMEvent***REMOVED***"ack_error", &RateLimitEvent***REMOVED******REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			rtm.IncomingEvents <- RTMEvent***REMOVED***"ack_error", &AckErrorEvent***REMOVED***ack.Error***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rtm.IncomingEvents <- RTMEvent***REMOVED***"ack_error", &AckErrorEvent***REMOVED***fmt.Errorf("ack decode failure")***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// handlePong handles an incoming 'PONG' message which should be in response to
// a previously sent 'PING' message. This is then used to compute the
// connection's latency.
func (rtm *RTM) handlePong(event json.RawMessage) ***REMOVED***
	pong := &Pong***REMOVED******REMOVED***
	if err := json.Unmarshal(event, pong); err != nil ***REMOVED***
		rtm.Debugln("RTM Error unmarshalling 'pong' event:", err)
		rtm.Debugln(" -> Erroneous 'ping' event:", string(event))
		return
	***REMOVED***
	if pingTime, exists := rtm.pings[pong.ReplyTo]; exists ***REMOVED***
		latency := time.Since(pingTime)
		rtm.IncomingEvents <- RTMEvent***REMOVED***"latency_report", &LatencyReport***REMOVED***Value: latency***REMOVED******REMOVED***
		delete(rtm.pings, pong.ReplyTo)
	***REMOVED*** else ***REMOVED***
		rtm.Debugln("RTM Error - unmatched 'pong' event:", string(event))
	***REMOVED***
***REMOVED***

// handleEvent is the "default" response to an event that does not have a
// special case. It matches the command's name to a mapping of defined events
// and then sends the corresponding event struct to the IncomingEvents channel.
// If the event type is not found or the event cannot be unmarshalled into the
// correct struct then this sends an UnmarshallingErrorEvent to the
// IncomingEvents channel.
func (rtm *RTM) handleEvent(typeStr string, event json.RawMessage) ***REMOVED***
	v, exists := eventMapping[typeStr]
	if !exists ***REMOVED***
		rtm.Debugf("RTM Error, received unmapped event %q: %s\n", typeStr, string(event))
		err := fmt.Errorf("RTM Error: Received unmapped event %q: %s\n", typeStr, string(event))
		rtm.IncomingEvents <- RTMEvent***REMOVED***"unmarshalling_error", &UnmarshallingErrorEvent***REMOVED***err***REMOVED******REMOVED***
		return
	***REMOVED***
	t := reflect.TypeOf(v)
	recvEvent := reflect.New(t).Interface()
	err := json.Unmarshal(event, recvEvent)
	if err != nil ***REMOVED***
		rtm.Debugf("RTM Error, could not unmarshall event %q: %s\n", typeStr, string(event))
		err := fmt.Errorf("RTM Error: Could not unmarshall event %q: %s\n", typeStr, string(event))
		rtm.IncomingEvents <- RTMEvent***REMOVED***"unmarshalling_error", &UnmarshallingErrorEvent***REMOVED***err***REMOVED******REMOVED***
		return
	***REMOVED***
	rtm.IncomingEvents <- RTMEvent***REMOVED***typeStr, recvEvent***REMOVED***
***REMOVED***

// eventMapping holds a mapping of event names to their corresponding struct
// implementations. The structs should be instances of the unmarshalling
// target for the matching event type.
var eventMapping = map[string]interface***REMOVED******REMOVED******REMOVED***
	"message":         MessageEvent***REMOVED******REMOVED***,
	"presence_change": PresenceChangeEvent***REMOVED******REMOVED***,
	"user_typing":     UserTypingEvent***REMOVED******REMOVED***,

	"channel_marked":          ChannelMarkedEvent***REMOVED******REMOVED***,
	"channel_created":         ChannelCreatedEvent***REMOVED******REMOVED***,
	"channel_joined":          ChannelJoinedEvent***REMOVED******REMOVED***,
	"channel_left":            ChannelLeftEvent***REMOVED******REMOVED***,
	"channel_deleted":         ChannelDeletedEvent***REMOVED******REMOVED***,
	"channel_rename":          ChannelRenameEvent***REMOVED******REMOVED***,
	"channel_archive":         ChannelArchiveEvent***REMOVED******REMOVED***,
	"channel_unarchive":       ChannelUnarchiveEvent***REMOVED******REMOVED***,
	"channel_history_changed": ChannelHistoryChangedEvent***REMOVED******REMOVED***,

	"dnd_updated":      DNDUpdatedEvent***REMOVED******REMOVED***,
	"dnd_updated_user": DNDUpdatedEvent***REMOVED******REMOVED***,

	"im_created":         IMCreatedEvent***REMOVED******REMOVED***,
	"im_open":            IMOpenEvent***REMOVED******REMOVED***,
	"im_close":           IMCloseEvent***REMOVED******REMOVED***,
	"im_marked":          IMMarkedEvent***REMOVED******REMOVED***,
	"im_history_changed": IMHistoryChangedEvent***REMOVED******REMOVED***,

	"group_marked":          GroupMarkedEvent***REMOVED******REMOVED***,
	"group_open":            GroupOpenEvent***REMOVED******REMOVED***,
	"group_joined":          GroupJoinedEvent***REMOVED******REMOVED***,
	"group_left":            GroupLeftEvent***REMOVED******REMOVED***,
	"group_close":           GroupCloseEvent***REMOVED******REMOVED***,
	"group_rename":          GroupRenameEvent***REMOVED******REMOVED***,
	"group_archive":         GroupArchiveEvent***REMOVED******REMOVED***,
	"group_unarchive":       GroupUnarchiveEvent***REMOVED******REMOVED***,
	"group_history_changed": GroupHistoryChangedEvent***REMOVED******REMOVED***,

	"file_created":         FileCreatedEvent***REMOVED******REMOVED***,
	"file_shared":          FileSharedEvent***REMOVED******REMOVED***,
	"file_unshared":        FileUnsharedEvent***REMOVED******REMOVED***,
	"file_public":          FilePublicEvent***REMOVED******REMOVED***,
	"file_private":         FilePrivateEvent***REMOVED******REMOVED***,
	"file_change":          FileChangeEvent***REMOVED******REMOVED***,
	"file_deleted":         FileDeletedEvent***REMOVED******REMOVED***,
	"file_comment_added":   FileCommentAddedEvent***REMOVED******REMOVED***,
	"file_comment_edited":  FileCommentEditedEvent***REMOVED******REMOVED***,
	"file_comment_deleted": FileCommentDeletedEvent***REMOVED******REMOVED***,

	"pin_added":   PinAddedEvent***REMOVED******REMOVED***,
	"pin_removed": PinRemovedEvent***REMOVED******REMOVED***,

	"star_added":   StarAddedEvent***REMOVED******REMOVED***,
	"star_removed": StarRemovedEvent***REMOVED******REMOVED***,

	"reaction_added":   ReactionAddedEvent***REMOVED******REMOVED***,
	"reaction_removed": ReactionRemovedEvent***REMOVED******REMOVED***,

	"pref_change": PrefChangeEvent***REMOVED******REMOVED***,

	"team_join":              TeamJoinEvent***REMOVED******REMOVED***,
	"team_rename":            TeamRenameEvent***REMOVED******REMOVED***,
	"team_pref_change":       TeamPrefChangeEvent***REMOVED******REMOVED***,
	"team_domain_change":     TeamDomainChangeEvent***REMOVED******REMOVED***,
	"team_migration_started": TeamMigrationStartedEvent***REMOVED******REMOVED***,

	"manual_presence_change": ManualPresenceChangeEvent***REMOVED******REMOVED***,

	"user_change": UserChangeEvent***REMOVED******REMOVED***,

	"emoji_changed": EmojiChangedEvent***REMOVED******REMOVED***,

	"commands_changed": CommandsChangedEvent***REMOVED******REMOVED***,

	"email_domain_changed": EmailDomainChangedEvent***REMOVED******REMOVED***,

	"bot_added":   BotAddedEvent***REMOVED******REMOVED***,
	"bot_changed": BotChangedEvent***REMOVED******REMOVED***,

	"accounts_changed": AccountsChangedEvent***REMOVED******REMOVED***,

	"reconnect_url": ReconnectUrlEvent***REMOVED******REMOVED***,
***REMOVED***
