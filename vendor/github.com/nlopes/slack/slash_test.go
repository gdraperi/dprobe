package slack

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestSlash_ServeHTTP(t *testing.T) ***REMOVED***
	once.Do(startServer)
	serverURL := fmt.Sprintf("http://%s/slash", serverAddr)

	tests := []struct ***REMOVED***
		body           url.Values
		wantParams     SlashCommand
		wantStatusCode int
	***REMOVED******REMOVED***
		***REMOVED***
			body: url.Values***REMOVED***
				"command":      []string***REMOVED***"/command"***REMOVED***,
				"team_domain":  []string***REMOVED***"team"***REMOVED***,
				"channel_id":   []string***REMOVED***"C1234ABCD"***REMOVED***,
				"text":         []string***REMOVED***"text"***REMOVED***,
				"team_id":      []string***REMOVED***"T1234ABCD"***REMOVED***,
				"user_id":      []string***REMOVED***"U1234ABCD"***REMOVED***,
				"user_name":    []string***REMOVED***"username"***REMOVED***,
				"response_url": []string***REMOVED***"https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY"***REMOVED***,
				"token":        []string***REMOVED***"valid"***REMOVED***,
				"channel_name": []string***REMOVED***"channel"***REMOVED***,
				"trigger_id":   []string***REMOVED***"0000000000.1111111111.222222222222aaaaaaaaaaaaaa"***REMOVED***,
			***REMOVED***,
			wantParams: SlashCommand***REMOVED***
				Command:     "/command",
				TeamDomain:  "team",
				ChannelID:   "C1234ABCD",
				Text:        "text",
				TeamID:      "T1234ABCD",
				UserID:      "U1234ABCD",
				UserName:    "username",
				ResponseURL: "https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY",
				Token:       "valid",
				ChannelName: "channel",
				TriggerID:   "0000000000.1111111111.222222222222aaaaaaaaaaaaaa",
			***REMOVED***,
			wantStatusCode: http.StatusOK,
		***REMOVED***,
		***REMOVED***
			body: url.Values***REMOVED***
				"token": []string***REMOVED***"invalid"***REMOVED***,
			***REMOVED***,
			wantParams: SlashCommand***REMOVED***
				Token: "invalid",
			***REMOVED***,
			wantStatusCode: http.StatusUnauthorized,
		***REMOVED***,
	***REMOVED***

	var slashCommand SlashCommand
	client := &http.Client***REMOVED******REMOVED***
	http.HandleFunc("/slash", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		var err error
		slashCommand, err = SlashCommandParse(r)
		if err != nil ***REMOVED***
			w.WriteHeader(http.StatusInternalServerError)
		***REMOVED***
		acceptableTokens := []string***REMOVED***"valid", "valid2"***REMOVED***
		if !slashCommand.ValidateToken(acceptableTokens...) ***REMOVED***
			w.WriteHeader(http.StatusUnauthorized)
		***REMOVED***
	***REMOVED***)

	for i, test := range tests ***REMOVED***
		req, err := http.NewRequest(http.MethodPost, serverURL, strings.NewReader(test.body.Encode()))
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***

		if resp.StatusCode != test.wantStatusCode ***REMOVED***
			t.Errorf("%d: Got status code %d, want %d", i, resp.StatusCode, test.wantStatusCode)
		***REMOVED***
		if !reflect.DeepEqual(slashCommand, test.wantParams) ***REMOVED***
			t.Errorf("%d: Got params %#v, want %#v", i, slashCommand, test.wantParams)
		***REMOVED***
		resp.Body.Close()
	***REMOVED***
***REMOVED***
