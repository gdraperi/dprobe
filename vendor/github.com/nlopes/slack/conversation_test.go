package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Channel
var simpleChannel = `***REMOVED***
    "id": "C024BE91L",
    "name": "fun",
    "is_channel": true,
    "created": 1360782804,
    "creator": "U024BE7LH",
    "is_archived": false,
    "is_general": false,
    "members": [
        "U024BE7LH"
    ],
    "topic": ***REMOVED***
        "value": "Fun times",
        "creator": "U024BE7LV",
        "last_set": 1369677212
***REMOVED***,
    "purpose": ***REMOVED***
        "value": "This channel is for fun",
        "creator": "U024BE7LH",
        "last_set": 1360782804
***REMOVED***,
    "is_member": true,
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0
***REMOVED***`

func unmarshalChannel(j string) (*Channel, error) ***REMOVED***
	channel := &Channel***REMOVED******REMOVED***
	if err := json.Unmarshal([]byte(j), &channel); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return channel, nil
***REMOVED***

func TestSimpleChannel(t *testing.T) ***REMOVED***
	channel, err := unmarshalChannel(simpleChannel)
	assert.Nil(t, err)
	assertSimpleChannel(t, channel)
***REMOVED***

func assertSimpleChannel(t *testing.T, channel *Channel) ***REMOVED***
	assert.NotNil(t, channel)
	assert.Equal(t, "C024BE91L", channel.ID)
	assert.Equal(t, "fun", channel.Name)
	assert.Equal(t, true, channel.IsChannel)
	assert.Equal(t, JSONTime(1360782804), channel.Created)
	assert.Equal(t, "U024BE7LH", channel.Creator)
	assert.Equal(t, false, channel.IsArchived)
	assert.Equal(t, false, channel.IsGeneral)
	assert.Equal(t, true, channel.IsMember)
	assert.Equal(t, "1401383885.000061", channel.LastRead)
	assert.Equal(t, 0, channel.UnreadCount)
	assert.Equal(t, 0, channel.UnreadCountDisplay)
***REMOVED***

func TestCreateSimpleChannel(t *testing.T) ***REMOVED***
	channel := &Channel***REMOVED******REMOVED***
	channel.ID = "C024BE91L"
	channel.Name = "fun"
	channel.IsChannel = true
	channel.Created = JSONTime(1360782804)
	channel.Creator = "U024BE7LH"
	channel.IsArchived = false
	channel.IsGeneral = false
	channel.IsMember = true
	channel.LastRead = "1401383885.000061"
	channel.UnreadCount = 0
	channel.UnreadCountDisplay = 0
	assertSimpleChannel(t, channel)
***REMOVED***

// Group
var simpleGroup = `***REMOVED***
    "id": "G024BE91L",
    "name": "secretplans",
    "is_group": true,
    "created": 1360782804,
    "creator": "U024BE7LH",
    "is_archived": false,
    "members": [
        "U024BE7LH"
    ],
    "topic": ***REMOVED***
        "value": "Secret plans on hold",
        "creator": "U024BE7LV",
        "last_set": 1369677212
***REMOVED***,
    "purpose": ***REMOVED***
        "value": "Discuss secret plans that no-one else should know",
        "creator": "U024BE7LH",
        "last_set": 1360782804
***REMOVED***,
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0
***REMOVED***`

func unmarshalGroup(j string) (*Group, error) ***REMOVED***
	group := &Group***REMOVED******REMOVED***
	if err := json.Unmarshal([]byte(j), &group); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return group, nil
***REMOVED***

func TestSimpleGroup(t *testing.T) ***REMOVED***
	group, err := unmarshalGroup(simpleGroup)
	assert.Nil(t, err)
	assertSimpleGroup(t, group)
***REMOVED***

func assertSimpleGroup(t *testing.T, group *Group) ***REMOVED***
	assert.NotNil(t, group)
	assert.Equal(t, "G024BE91L", group.ID)
	assert.Equal(t, "secretplans", group.Name)
	assert.Equal(t, true, group.IsGroup)
	assert.Equal(t, JSONTime(1360782804), group.Created)
	assert.Equal(t, "U024BE7LH", group.Creator)
	assert.Equal(t, false, group.IsArchived)
	assert.Equal(t, "1401383885.000061", group.LastRead)
	assert.Equal(t, 0, group.UnreadCount)
	assert.Equal(t, 0, group.UnreadCountDisplay)
***REMOVED***

func TestCreateSimpleGroup(t *testing.T) ***REMOVED***
	group := &Group***REMOVED******REMOVED***
	group.ID = "G024BE91L"
	group.Name = "secretplans"
	group.IsGroup = true
	group.Created = JSONTime(1360782804)
	group.Creator = "U024BE7LH"
	group.IsArchived = false
	group.LastRead = "1401383885.000061"
	group.UnreadCount = 0
	group.UnreadCountDisplay = 0
	assertSimpleGroup(t, group)
***REMOVED***

// IM
var simpleIM = `***REMOVED***
    "id": "D024BFF1M",
    "is_im": true,
    "user": "U024BE7LH",
    "created": 1360782804,
    "is_user_deleted": false,
    "is_open": true,
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0
***REMOVED***`

func unmarshalIM(j string) (*IM, error) ***REMOVED***
	im := &IM***REMOVED******REMOVED***
	if err := json.Unmarshal([]byte(j), &im); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return im, nil
***REMOVED***

func TestSimpleIM(t *testing.T) ***REMOVED***
	im, err := unmarshalIM(simpleIM)
	assert.Nil(t, err)
	assertSimpleIM(t, im)
***REMOVED***

func assertSimpleIM(t *testing.T, im *IM) ***REMOVED***
	assert.NotNil(t, im)
	assert.Equal(t, "D024BFF1M", im.ID)
	assert.Equal(t, true, im.IsIM)
	assert.Equal(t, JSONTime(1360782804), im.Created)
	assert.Equal(t, false, im.IsUserDeleted)
	assert.Equal(t, true, im.IsOpen)
	assert.Equal(t, "1401383885.000061", im.LastRead)
	assert.Equal(t, 0, im.UnreadCount)
	assert.Equal(t, 0, im.UnreadCountDisplay)
***REMOVED***

func TestCreateSimpleIM(t *testing.T) ***REMOVED***
	im := &IM***REMOVED******REMOVED***
	im.ID = "D024BFF1M"
	im.IsIM = true
	im.Created = JSONTime(1360782804)
	im.IsUserDeleted = false
	im.IsOpen = true
	im.LastRead = "1401383885.000061"
	im.UnreadCount = 0
	im.UnreadCountDisplay = 0
	assertSimpleIM(t, im)
***REMOVED***

func getTestMembers() []string ***REMOVED***
	return []string***REMOVED***"test"***REMOVED***
***REMOVED***

func getUsersInConversation(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		SlackResponse
		Members          []string         `json:"members"`
		ResponseMetaData responseMetaData `json:"response_metadata"`
	***REMOVED******REMOVED***
		SlackResponse:    SlackResponse***REMOVED***Ok: true***REMOVED***,
		Members:          getTestMembers(),
		ResponseMetaData: responseMetaData***REMOVED***NextCursor: ""***REMOVED***,
	***REMOVED***)
	rw.Write(response)
***REMOVED***

func TestGetUsersInConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.members", getUsersInConversation)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := GetUsersInConversationParameters***REMOVED***
		ChannelID: "CXXXXXXXX",
	***REMOVED***

	expectedMembers := getTestMembers()

	members, _, err := api.GetUsersInConversation(&params)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if !reflect.DeepEqual(expectedMembers, members) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
***REMOVED***

func TestArchiveConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.archive", okJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	err := api.ArchiveConversation("CXXXXXXXX")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func TestUnArchiveConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.unarchive", okJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	err := api.UnArchiveConversation("CXXXXXXXX")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func getTestChannel() *Channel ***REMOVED***
	return &Channel***REMOVED***
		groupConversation: groupConversation***REMOVED***
			Topic: Topic***REMOVED***
				Value: "response topic",
			***REMOVED***,
			Purpose: Purpose***REMOVED***
				Value: "response purpose",
			***REMOVED***,
		***REMOVED******REMOVED***
***REMOVED***

func okChannelJsonHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		SlackResponse
		Channel *Channel `json:"channel"`
	***REMOVED******REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED***,
		Channel:       getTestChannel(),
	***REMOVED***)
	rw.Write(response)
***REMOVED***

func TestSetTopicOfConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.setTopic", okChannelJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	inputChannel := getTestChannel()
	channel, err := api.SetTopicOfConversation("CXXXXXXXX", inputChannel.Topic.Value)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if channel.Topic.Value != inputChannel.Topic.Value ***REMOVED***
		t.Fatalf(`topic = '%s', want '%s'`, channel.Topic.Value, inputChannel.Topic.Value)
	***REMOVED***
***REMOVED***

func TestSetPurposeOfConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.setPurpose", okChannelJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	inputChannel := getTestChannel()
	channel, err := api.SetPurposeOfConversation("CXXXXXXXX", inputChannel.Purpose.Value)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if channel.Purpose.Value != inputChannel.Purpose.Value ***REMOVED***
		t.Fatalf(`purpose = '%s', want '%s'`, channel.Purpose.Value, inputChannel.Purpose.Value)
	***REMOVED***
***REMOVED***

func TestRenameConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.rename", okChannelJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	inputChannel := getTestChannel()
	channel, err := api.RenameConversation("CXXXXXXXX", inputChannel.Name)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if channel.Name != inputChannel.Name ***REMOVED***
		t.Fatalf(`channelName = '%s', want '%s'`, channel.Name, inputChannel.Name)
	***REMOVED***
***REMOVED***

func TestInviteUsersToConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.invite", okChannelJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	users := []string***REMOVED***"UXXXXXXX1", "UXXXXXXX2"***REMOVED***
	channel, err := api.InviteUsersToConversation("CXXXXXXXX", users...)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if channel == nil ***REMOVED***
		t.Error("channel should not be nil")
		return
	***REMOVED***
***REMOVED***

func TestKickUserFromConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.kick", okJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	err := api.KickUserFromConversation("CXXXXXXXX", "UXXXXXXXX")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func closeConversationHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		SlackResponse
		NoOp          bool `json:"no_op"`
		AlreadyClosed bool `json:"already_closed"`
	***REMOVED******REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED******REMOVED***)
	rw.Write(response)
***REMOVED***

func TestCloseConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.close", closeConversationHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	_, _, err := api.CloseConversation("CXXXXXXXX")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func TestCreateConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.create", okChannelJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	channel, err := api.CreateConversation("CXXXXXXXX", false)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if channel == nil ***REMOVED***
		t.Error("channel should not be nil")
		return
	***REMOVED***
***REMOVED***

func TestGetConversationInfo(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.info", okChannelJsonHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	channel, err := api.GetConversationInfo("CXXXXXXXX", false)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if channel == nil ***REMOVED***
		t.Error("channel should not be nil")
		return
	***REMOVED***
***REMOVED***

func leaveConversationHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		SlackResponse
		NotInChannel bool `json:"not_in_channel"`
	***REMOVED******REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED******REMOVED***)
	rw.Write(response)
***REMOVED***

func TestLeaveConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.leave", leaveConversationHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	_, err := api.LeaveConversation("CXXXXXXXX")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func getConversationRepliesHander(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		SlackResponse
		HasMore          bool `json:"has_more"`
		ResponseMetaData struct ***REMOVED***
			NextCursor string `json:"next_cursor"`
		***REMOVED*** `json:"response_metadata"`
		Messages []Message `json:"messages"`
	***REMOVED******REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED***,
		Messages:      []Message***REMOVED******REMOVED******REMOVED***)
	rw.Write(response)
***REMOVED***

func TestGetConversationReplies(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.replies", getConversationRepliesHander)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := GetConversationRepliesParameters***REMOVED***
		ChannelID: "CXXXXXXXX",
		Timestamp: "1234567890.123456",
	***REMOVED***
	_, _, _, err := api.GetConversationReplies(&params)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func getConversationsHander(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		SlackResponse
		ResponseMetaData struct ***REMOVED***
			NextCursor string `json:"next_cursor"`
		***REMOVED*** `json:"response_metadata"`
		Channels []Channel `json:"channels"`
	***REMOVED******REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED***,
		Channels:      []Channel***REMOVED******REMOVED******REMOVED***)
	rw.Write(response)
***REMOVED***

func TestGetConversations(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.list", getConversationsHander)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := GetConversationsParameters***REMOVED******REMOVED***
	_, _, err := api.GetConversations(&params)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func openConversationHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		SlackResponse
		NoOp        bool     `json:"no_op"`
		AlreadyOpen bool     `json:"already_open"`
		Channel     *Channel `json:"channel"`
	***REMOVED******REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED******REMOVED***)
	rw.Write(response)
***REMOVED***

func TestOpenConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.open", openConversationHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := OpenConversationParameters***REMOVED***ChannelID: "CXXXXXXXX"***REMOVED***
	_, _, _, err := api.OpenConversation(&params)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func joinConversationHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		Channel          *Channel `json:"channel"`
		Warning          string   `json:"warning"`
		ResponseMetaData *struct ***REMOVED***
			Warnings []string `json:"warnings"`
		***REMOVED*** `json:"response_metadata"`
		SlackResponse
	***REMOVED******REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED******REMOVED***)
	rw.Write(response)
***REMOVED***

func TestJoinConversation(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.join", joinConversationHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	_, _, _, err := api.JoinConversation("CXXXXXXXX")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***

func getConversationHistoryHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(GetConversationHistoryResponse***REMOVED***
		SlackResponse: SlackResponse***REMOVED***Ok: true***REMOVED******REMOVED***)
	rw.Write(response)
***REMOVED***

func TestGetConversationHistory(t *testing.T) ***REMOVED***
	http.HandleFunc("/conversations.history", getConversationHistoryHandler)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	params := GetConversationHistoryParameters***REMOVED***ChannelID: "CXXXXXXXX"***REMOVED***
	_, err := api.GetConversationHistory(&params)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
***REMOVED***
