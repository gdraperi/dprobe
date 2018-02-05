package slack

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type starsHandler struct ***REMOVED***
	gotParams map[string]string
	response  string
***REMOVED***

func newStarsHandler() *starsHandler ***REMOVED***
	return &starsHandler***REMOVED***
		gotParams: make(map[string]string),
		response:  `***REMOVED*** "ok": true ***REMOVED***`,
	***REMOVED***
***REMOVED***

func (sh *starsHandler) accumulateFormValue(k string, r *http.Request) ***REMOVED***
	if v := r.FormValue(k); v != "" ***REMOVED***
		sh.gotParams[k] = v
	***REMOVED***
***REMOVED***

func (sh *starsHandler) handler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	sh.accumulateFormValue("user", r)
	sh.accumulateFormValue("count", r)
	sh.accumulateFormValue("channel", r)
	sh.accumulateFormValue("file", r)
	sh.accumulateFormValue("file_comment", r)
	sh.accumulateFormValue("page", r)
	sh.accumulateFormValue("timestamp", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(sh.response))
***REMOVED***

func TestSlack_AddStar(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	tests := []struct ***REMOVED***
		channel    string
		ref        ItemRef
		wantParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			"ChannelID",
			NewRefToMessage("ChannelID", "123"),
			map[string]string***REMOVED***
				"channel":   "ChannelID",
				"timestamp": "123",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"ChannelID",
			NewRefToFile("FileID"),
			map[string]string***REMOVED***
				"channel": "ChannelID",
				"file":    "FileID",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"ChannelID",
			NewRefToComment("FileCommentID"),
			map[string]string***REMOVED***
				"channel":      "ChannelID",
				"file_comment": "FileCommentID",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	var rh *starsHandler
	http.HandleFunc("/stars.add", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)
	for i, test := range tests ***REMOVED***
		rh = newStarsHandler()
		err := api.AddStar(test.channel, test.ref)
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) ***REMOVED***
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSlack_RemoveStar(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	tests := []struct ***REMOVED***
		channel    string
		ref        ItemRef
		wantParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			"ChannelID",
			NewRefToMessage("ChannelID", "123"),
			map[string]string***REMOVED***
				"channel":   "ChannelID",
				"timestamp": "123",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"ChannelID",
			NewRefToFile("FileID"),
			map[string]string***REMOVED***
				"channel": "ChannelID",
				"file":    "FileID",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"ChannelID",
			NewRefToComment("FileCommentID"),
			map[string]string***REMOVED***
				"channel":      "ChannelID",
				"file_comment": "FileCommentID",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	var rh *starsHandler
	http.HandleFunc("/stars.remove", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)
	for i, test := range tests ***REMOVED***
		rh = newStarsHandler()
		err := api.RemoveStar(test.channel, test.ref)
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) ***REMOVED***
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSlack_ListStars(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	rh := newStarsHandler()
	http.HandleFunc("/stars.list", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)
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
	want := []Item***REMOVED***
		NewMessageItem("C1", &Message***REMOVED***Msg: Msg***REMOVED***
			Text: "hello",
			Reactions: []ItemReaction***REMOVED***
				ItemReaction***REMOVED***Name: "astonished", Count: 3, Users: []string***REMOVED***"U1", "U2", "U3"***REMOVED******REMOVED***,
				ItemReaction***REMOVED***Name: "clock1", Count: 3, Users: []string***REMOVED***"U1", "U2"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED******REMOVED***),
		NewFileItem(&File***REMOVED***Name: "toy"***REMOVED***),
		NewFileCommentItem(&File***REMOVED***Name: "toy"***REMOVED***, &Comment***REMOVED***Comment: "cool toy"***REMOVED***),
	***REMOVED***
	wantStarred := make([]StarredItem, len(want))
	for i, item := range want ***REMOVED***
		wantStarred[i] = StarredItem(item)
	***REMOVED***
	wantParams := map[string]string***REMOVED***
		"count": "200",
		"page":  "2",
	***REMOVED***
	params := NewStarsParameters()
	params.Count = 200
	params.Page = 2
	got, paging, err := api.ListStars(params)
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
	if !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("Got Stars %#v, want %#v", got, want)
		for i, item := range got ***REMOVED***
			fmt.Printf("Item %d, Type: %s\n", i, item.Type)
			fmt.Printf("Message  %#v\n", item.Message)
			fmt.Printf("File     %#v\n", item.File)
			fmt.Printf("Comment  %#v\n", item.Comment)
		***REMOVED***
	***REMOVED***
	if !reflect.DeepEqual(rh.gotParams, wantParams) ***REMOVED***
		t.Errorf("Got params %#v, want %#v", rh.gotParams, wantParams)
	***REMOVED***
	if reflect.DeepEqual(paging, Paging***REMOVED******REMOVED***) ***REMOVED***
		t.Errorf("Want paging data, got empty struct")
	***REMOVED***
	// Test GetStarred
	gotStarred, paging, err := api.GetStarred(params)
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
	if !reflect.DeepEqual(gotStarred, wantStarred) ***REMOVED***
		t.Errorf("Got Stars %#v, want %#v", gotStarred, wantStarred)
		for i, item := range got ***REMOVED***
			fmt.Printf("Item %d, Type: %s\n", i, item.Type)
			fmt.Printf("Message  %#v\n", item.Message)
			fmt.Printf("File     %#v\n", item.File)
			fmt.Printf("Comment  %#v\n", item.Comment)
		***REMOVED***
	***REMOVED***
	if !reflect.DeepEqual(rh.gotParams, wantParams) ***REMOVED***
		t.Errorf("Got params %#v, want %#v", rh.gotParams, wantParams)
	***REMOVED***
	if reflect.DeepEqual(paging, Paging***REMOVED******REMOVED***) ***REMOVED***
		t.Errorf("Want paging data, got empty struct")
	***REMOVED***
***REMOVED***
