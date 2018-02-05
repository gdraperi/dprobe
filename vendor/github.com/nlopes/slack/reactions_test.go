package slack

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type reactionsHandler struct ***REMOVED***
	gotParams map[string]string
	response  string
***REMOVED***

func newReactionsHandler() *reactionsHandler ***REMOVED***
	return &reactionsHandler***REMOVED***
		gotParams: make(map[string]string),
		response:  `***REMOVED*** "ok": true ***REMOVED***`,
	***REMOVED***
***REMOVED***

func (rh *reactionsHandler) accumulateFormValue(k string, r *http.Request) ***REMOVED***
	if v := r.FormValue(k); v != "" ***REMOVED***
		rh.gotParams[k] = v
	***REMOVED***
***REMOVED***

func (rh *reactionsHandler) handler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	rh.accumulateFormValue("channel", r)
	rh.accumulateFormValue("count", r)
	rh.accumulateFormValue("file", r)
	rh.accumulateFormValue("file_comment", r)
	rh.accumulateFormValue("full", r)
	rh.accumulateFormValue("name", r)
	rh.accumulateFormValue("page", r)
	rh.accumulateFormValue("timestamp", r)
	rh.accumulateFormValue("user", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rh.response))
***REMOVED***

func TestSlack_AddReaction(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	tests := []struct ***REMOVED***
		name       string
		ref        ItemRef
		wantParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			"thumbsup",
			NewRefToMessage("ChannelID", "123"),
			map[string]string***REMOVED***
				"name":      "thumbsup",
				"channel":   "ChannelID",
				"timestamp": "123",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"thumbsup",
			NewRefToFile("FileID"),
			map[string]string***REMOVED***
				"name": "thumbsup",
				"file": "FileID",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"thumbsup",
			NewRefToComment("FileCommentID"),
			map[string]string***REMOVED***
				"name":         "thumbsup",
				"file_comment": "FileCommentID",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	var rh *reactionsHandler
	http.HandleFunc("/reactions.add", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)
	for i, test := range tests ***REMOVED***
		rh = newReactionsHandler()
		err := api.AddReaction(test.name, test.ref)
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) ***REMOVED***
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSlack_RemoveReaction(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	tests := []struct ***REMOVED***
		name       string
		ref        ItemRef
		wantParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			"thumbsup",
			NewRefToMessage("ChannelID", "123"),
			map[string]string***REMOVED***
				"name":      "thumbsup",
				"channel":   "ChannelID",
				"timestamp": "123",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"thumbsup",
			NewRefToFile("FileID"),
			map[string]string***REMOVED***
				"name": "thumbsup",
				"file": "FileID",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"thumbsup",
			NewRefToComment("FileCommentID"),
			map[string]string***REMOVED***
				"name":         "thumbsup",
				"file_comment": "FileCommentID",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	var rh *reactionsHandler
	http.HandleFunc("/reactions.remove", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)
	for i, test := range tests ***REMOVED***
		rh = newReactionsHandler()
		err := api.RemoveReaction(test.name, test.ref)
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) ***REMOVED***
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSlack_GetReactions(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	tests := []struct ***REMOVED***
		ref           ItemRef
		params        GetReactionsParameters
		wantParams    map[string]string
		json          string
		wantReactions []ItemReaction
	***REMOVED******REMOVED***
		***REMOVED***
			NewRefToMessage("ChannelID", "123"),
			GetReactionsParameters***REMOVED******REMOVED***,
			map[string]string***REMOVED***
				"channel":   "ChannelID",
				"timestamp": "123",
			***REMOVED***,
			`***REMOVED***"ok": true,
    "type": "message",
    "message": ***REMOVED***
        "reactions": [
            ***REMOVED***
                "name": "astonished",
                "count": 3,
                "users": [ "U1", "U2", "U3" ]
        ***REMOVED***,
            ***REMOVED***
                "name": "clock1",
                "count": 3,
                "users": [ "U1", "U2" ]
        ***REMOVED***
        ]
***REMOVED******REMOVED***`,
			[]ItemReaction***REMOVED***
				ItemReaction***REMOVED***Name: "astonished", Count: 3, Users: []string***REMOVED***"U1", "U2", "U3"***REMOVED******REMOVED***,
				ItemReaction***REMOVED***Name: "clock1", Count: 3, Users: []string***REMOVED***"U1", "U2"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			NewRefToFile("FileID"),
			GetReactionsParameters***REMOVED***Full: true***REMOVED***,
			map[string]string***REMOVED***
				"file": "FileID",
				"full": "true",
			***REMOVED***,
			`***REMOVED***"ok": true,
    "type": "file",
    "file": ***REMOVED***
        "reactions": [
            ***REMOVED***
                "name": "astonished",
                "count": 3,
                "users": [ "U1", "U2", "U3" ]
        ***REMOVED***,
            ***REMOVED***
                "name": "clock1",
                "count": 3,
                "users": [ "U1", "U2" ]
        ***REMOVED***
        ]
***REMOVED******REMOVED***`,
			[]ItemReaction***REMOVED***
				ItemReaction***REMOVED***Name: "astonished", Count: 3, Users: []string***REMOVED***"U1", "U2", "U3"***REMOVED******REMOVED***,
				ItemReaction***REMOVED***Name: "clock1", Count: 3, Users: []string***REMOVED***"U1", "U2"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***

			NewRefToComment("FileCommentID"),
			GetReactionsParameters***REMOVED******REMOVED***,
			map[string]string***REMOVED***
				"file_comment": "FileCommentID",
			***REMOVED***,
			`***REMOVED***"ok": true,
    "type": "file_comment",
    "file": ***REMOVED******REMOVED***,
    "comment": ***REMOVED***
        "reactions": [
            ***REMOVED***
                "name": "astonished",
                "count": 3,
                "users": [ "U1", "U2", "U3" ]
        ***REMOVED***,
            ***REMOVED***
                "name": "clock1",
                "count": 3,
                "users": [ "U1", "U2" ]
        ***REMOVED***
        ]
***REMOVED******REMOVED***`,
			[]ItemReaction***REMOVED***
				ItemReaction***REMOVED***Name: "astonished", Count: 3, Users: []string***REMOVED***"U1", "U2", "U3"***REMOVED******REMOVED***,
				ItemReaction***REMOVED***Name: "clock1", Count: 3, Users: []string***REMOVED***"U1", "U2"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	var rh *reactionsHandler
	http.HandleFunc("/reactions.get", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)
	for i, test := range tests ***REMOVED***
		rh = newReactionsHandler()
		rh.response = test.json
		got, err := api.GetReactions(test.ref, test.params)
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(got, test.wantReactions) ***REMOVED***
			t.Errorf("%d: Got reaction %#v, want %#v", i, got, test.wantReactions)
		***REMOVED***
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) ***REMOVED***
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSlack_ListReactions(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	rh := newReactionsHandler()
	http.HandleFunc("/reactions.list", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)
	rh.response = `***REMOVED***"ok": true,
    "items": [
        ***REMOVED***
            "type": "message",
            "channel": "C1",
            "message": ***REMOVED***
                "text": "hello",
                "reactions": [
                    ***REMOVED***
                        "name": "astonished",
                        "count": 3,
                        "users": [ "U1", "U2", "U3" ]
                ***REMOVED***,
                    ***REMOVED***
                        "name": "clock1",
                        "count": 3,
                        "users": [ "U1", "U2" ]
                ***REMOVED***
                ]
        ***REMOVED***
    ***REMOVED***,
        ***REMOVED***
            "type": "file",
            "file": ***REMOVED***
                "name": "toy",
                "reactions": [
                    ***REMOVED***
                        "name": "clock1",
                        "count": 3,
                        "users": [ "U1", "U2" ]
                ***REMOVED***
                ]
        ***REMOVED***
    ***REMOVED***,
        ***REMOVED***
            "type": "file_comment",
            "file": ***REMOVED***
                "name": "toy"
        ***REMOVED***,
            "comment": ***REMOVED***
                "comment": "cool toy",
                "reactions": [
                    ***REMOVED***
                        "name": "astonished",
                        "count": 3,
                        "users": [ "U1", "U2", "U3" ]
                ***REMOVED***
                ]
        ***REMOVED***
    ***REMOVED***
    ],
    "paging": ***REMOVED***
        "count": 100,
        "total": 4,
        "page": 1,
        "pages": 1
***REMOVED******REMOVED***`
	want := []ReactedItem***REMOVED***
		ReactedItem***REMOVED***
			Item: NewMessageItem("C1", &Message***REMOVED***Msg: Msg***REMOVED***
				Text: "hello",
				Reactions: []ItemReaction***REMOVED***
					ItemReaction***REMOVED***Name: "astonished", Count: 3, Users: []string***REMOVED***"U1", "U2", "U3"***REMOVED******REMOVED***,
					ItemReaction***REMOVED***Name: "clock1", Count: 3, Users: []string***REMOVED***"U1", "U2"***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED******REMOVED***),
			Reactions: []ItemReaction***REMOVED***
				ItemReaction***REMOVED***Name: "astonished", Count: 3, Users: []string***REMOVED***"U1", "U2", "U3"***REMOVED******REMOVED***,
				ItemReaction***REMOVED***Name: "clock1", Count: 3, Users: []string***REMOVED***"U1", "U2"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		ReactedItem***REMOVED***
			Item: NewFileItem(&File***REMOVED***Name: "toy"***REMOVED***),
			Reactions: []ItemReaction***REMOVED***
				ItemReaction***REMOVED***Name: "clock1", Count: 3, Users: []string***REMOVED***"U1", "U2"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		ReactedItem***REMOVED***
			Item: NewFileCommentItem(&File***REMOVED***Name: "toy"***REMOVED***, &Comment***REMOVED***Comment: "cool toy"***REMOVED***),
			Reactions: []ItemReaction***REMOVED***
				ItemReaction***REMOVED***Name: "astonished", Count: 3, Users: []string***REMOVED***"U1", "U2", "U3"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	wantParams := map[string]string***REMOVED***
		"user":  "User",
		"count": "200",
		"page":  "2",
		"full":  "true",
	***REMOVED***
	params := NewListReactionsParameters()
	params.User = "User"
	params.Count = 200
	params.Page = 2
	params.Full = true
	got, paging, err := api.ListReactions(params)
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
	if !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("Got reaction %#v, want %#v", got, want)
		for i, item := range got ***REMOVED***
			fmt.Printf("Item %d, Type: %s\n", i, item.Type)
			fmt.Printf("Message  %#v\n", item.Message)
			fmt.Printf("File     %#v\n", item.File)
			fmt.Printf("Comment  %#v\n", item.Comment)
			fmt.Printf("Reactions %#v\n", item.Reactions)
		***REMOVED***
	***REMOVED***
	if !reflect.DeepEqual(rh.gotParams, wantParams) ***REMOVED***
		t.Errorf("Got params %#v, want %#v", rh.gotParams, wantParams)
	***REMOVED***
	if reflect.DeepEqual(paging, Paging***REMOVED******REMOVED***) ***REMOVED***
		t.Errorf("Want paging data, got empty struct")
	***REMOVED***
***REMOVED***
