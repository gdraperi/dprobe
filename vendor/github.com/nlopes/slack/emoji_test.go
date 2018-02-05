package slack

import (
	"net/http"
	"reflect"
	"testing"
)

func getEmojiHandler(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`***REMOVED***"ok": true, "emoji": ***REMOVED***
			"bowtie": "https://my.slack.com/emoji/bowtie/46ec6f2bb0.png",
			"squirrel": "https://my.slack.com/emoji/squirrel/f35f40c0e0.png",
			"shipit": "alias:squirrel"
		***REMOVED******REMOVED***`)
	rw.Write(response)
***REMOVED***

func TestGetEmoji(t *testing.T) ***REMOVED***
	http.HandleFunc("/emoji.list", getEmojiHandler)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	emojisResponse := map[string]string***REMOVED***
		"bowtie":   "https://my.slack.com/emoji/bowtie/46ec6f2bb0.png",
		"squirrel": "https://my.slack.com/emoji/squirrel/f35f40c0e0.png",
		"shipit":   "alias:squirrel",
	***REMOVED***

	emojis, err := api.GetEmoji()
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	eq := reflect.DeepEqual(emojis, emojisResponse)
	if !eq ***REMOVED***
		t.Errorf("got %v; want %v", emojis, emojisResponse)
	***REMOVED***
***REMOVED***
