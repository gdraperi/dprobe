package slack

import (
	"net/http"
	"testing"
)

func getBotInfo(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`***REMOVED***"ok": true, "bot": ***REMOVED***
			"id":"B02875YLA",
			"deleted":false,
			"name":"github",
			"icons": ***REMOVED***
              "image_36":"https:\/\/a.slack-edge.com\/2fac\/plugins\/github\/assets\/service_36.png",
              "image_48":"https:\/\/a.slack-edge.com\/2fac\/plugins\/github\/assets\/service_48.png",
              "image_72":"https:\/\/a.slack-edge.com\/2fac\/plugins\/github\/assets\/service_72.png"
        ***REMOVED***
    ***REMOVED******REMOVED***`)
	rw.Write(response)
***REMOVED***

func TestGetBotInfo(t *testing.T) ***REMOVED***
	http.HandleFunc("/bots.info", getBotInfo)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	bot, err := api.GetBotInfo("B02875YLA")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***

	if bot.ID != "B02875YLA" ***REMOVED***
		t.Fatal("Incorrect ID")
	***REMOVED***
	if bot.Name != "github" ***REMOVED***
		t.Fatal("Incorrect Name")
	***REMOVED***
	if len(bot.Icons.Image36) == 0 ***REMOVED***
		t.Fatal("Missing Image36")
	***REMOVED***
	if len(bot.Icons.Image48) == 0 ***REMOVED***
		t.Fatal("Missing Image38")
	***REMOVED***
	if len(bot.Icons.Image72) == 0 ***REMOVED***
		t.Fatal("Missing Image72")
	***REMOVED***
***REMOVED***
