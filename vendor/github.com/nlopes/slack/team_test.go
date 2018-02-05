package slack

import (
	"errors"
	"net/http"
	"testing"
	"strings"
)

var (
	ErrIncorrectResponse = errors.New("Response is incorrect")
)

func getTeamInfo(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`***REMOVED***"ok": true, "team": ***REMOVED***
			"id": "F0UWHUX",
			"name": "notalar",
			"domain": "notalar",
			"icon": ***REMOVED***
              "image_34": "https://slack.global.ssl.fastly.net/66f9/img/avatars-teams/ava_0002-34.png",
              "image_44": "https://slack.global.ssl.fastly.net/66f9/img/avatars-teams/ava_0002-44.png",
              "image_55": "https://slack.global.ssl.fastly.net/66f9/img/avatars-teams/ava_0002-55.png",
              "image_default": true
      ***REMOVED***
		***REMOVED******REMOVED***`)
	rw.Write(response)
***REMOVED***

func TestGetTeamInfo(t *testing.T) ***REMOVED***
	http.HandleFunc("/team.info", getTeamInfo)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	teamInfo, err := api.GetTeamInfo()
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***

	// t.Fatal refers to -> t.Errorf & return
	if teamInfo.ID != "F0UWHUX" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if teamInfo.Domain != "notalar" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if teamInfo.Name != "notalar" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if teamInfo.Icon == nil ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
***REMOVED***

func getTeamAccessLogs(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`***REMOVED***"ok": true, "logins": [***REMOVED***
			"user_id": "F0UWHUX",
			"username": "notalar",
			"date_first": 1475684477,
			"date_last": 1475684645,
			"count": 8,
			"ip": "127.0.0.1",
			"user_agent": "SlackWeb/3abb0ae2380d48a9ae20c58cc624ebcd Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Slack/1.2.6 Chrome/45.0.2454.85 AtomShell/0.34.3 Safari/537.36 Slack_SSB/1.2.6",
			"isp": "AT&T U-verse",
                        "country": "US",
                        "region": "IN"
                    ***REMOVED***,
                        ***REMOVED***
                        "user_id": "XUHWU0F",
			"username": "ralaton",
			"date_first": 1447395893,
			"date_last": 1447395965,
			"count": 5,
			"ip": "192.168.0.1",
			"user_agent": "com.tinyspeck.chatlyio/2.60 (iPhone; iOS 9.1; Scale/3.00)",
			"isp": null,
                        "country": null,
                        "region": null
                    ***REMOVED***],
                        "paging": ***REMOVED***
    			"count": 2,
    			"total": 2,
    			"page": 1,
    			"pages": 1
    			***REMOVED***
  ***REMOVED***`)
	rw.Write(response)
***REMOVED***

func TestGetAccessLogs(t *testing.T) ***REMOVED***
	http.HandleFunc("/team.accessLogs", getTeamAccessLogs)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	logins, paging, err := api.GetAccessLogs(NewAccessLogParameters())
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***

	if len(logins) != 2 ***REMOVED***
		t.Fatal("Should have been 2 logins")
	***REMOVED***

	// test the first login
	login1 := logins[0]
	login2 := logins[1]

	if (login1.UserID != "F0UWHUX") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.Username != "notalar") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.DateFirst != 1475684477) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.DateLast != 1475684645) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.Count != 8) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.IP != "127.0.0.1") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (!strings.HasPrefix(login1.UserAgent, "SlackWeb")) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.ISP != "AT&T U-verse") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.Country != "US") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login1.Region != "IN") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***

	// test that the null values from login2 are coming across correctly
	if (login2.ISP != "") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login2.Country != "") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (login2.Region != "") ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***

	// test the paging
	if (paging.Count != 2) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (paging.Total != 2) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (paging.Page != 1) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if (paging.Pages != 1) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
***REMOVED***

