package slack

import (
	"net/http"
	"reflect"
	"testing"
)

func TestSlack_EndDND(t *testing.T) ***REMOVED***
	http.HandleFunc("/dnd.endDnd", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`***REMOVED*** "ok": true ***REMOVED***`))
	***REMOVED***)
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	err := api.EndDND()
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
***REMOVED***

func TestSlack_EndSnooze(t *testing.T) ***REMOVED***
	http.HandleFunc("/dnd.endSnooze", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`***REMOVED*** "ok": true,
                          "dnd_enabled": true,
                          "next_dnd_start_ts": 1450418400,
                          "next_dnd_end_ts": 1450454400,
                          "snooze_enabled": false ***REMOVED***`))
	***REMOVED***)
	state := DNDStatus***REMOVED***
		Enabled:            true,
		NextStartTimestamp: 1450418400,
		NextEndTimestamp:   1450454400,
		SnoozeInfo:         SnoozeInfo***REMOVED***SnoozeEnabled: false***REMOVED***,
	***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	snoozeState, err := api.EndSnooze()
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
	eq := reflect.DeepEqual(snoozeState, &state)
	if !eq ***REMOVED***
		t.Errorf("got %v; want %v", snoozeState, &state)
	***REMOVED***
***REMOVED***

func TestSlack_GetDNDInfo(t *testing.T) ***REMOVED***
	http.HandleFunc("/dnd.info", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`***REMOVED***
            "ok": true,
            "dnd_enabled": true,
            "next_dnd_start_ts": 1450416600,
            "next_dnd_end_ts": 1450452600,
            "snooze_enabled": true,
            "snooze_endtime": 1450416600,
            "snooze_remaining": 1196
    ***REMOVED***`))
	***REMOVED***)
	userDNDInfo := DNDStatus***REMOVED***
		Enabled:            true,
		NextStartTimestamp: 1450416600,
		NextEndTimestamp:   1450452600,
		SnoozeInfo: SnoozeInfo***REMOVED***
			SnoozeEnabled:   true,
			SnoozeEndTime:   1450416600,
			SnoozeRemaining: 1196,
		***REMOVED***,
	***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	userDNDInfoResponse, err := api.GetDNDInfo(nil)
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
	eq := reflect.DeepEqual(userDNDInfoResponse, &userDNDInfo)
	if !eq ***REMOVED***
		t.Errorf("got %v; want %v", userDNDInfoResponse, &userDNDInfo)
	***REMOVED***
***REMOVED***

func TestSlack_GetDNDTeamInfo(t *testing.T) ***REMOVED***
	http.HandleFunc("/dnd.teamInfo", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`***REMOVED***
            "ok": true,
            "users": ***REMOVED***
                "U023BECGF": ***REMOVED***
                    "dnd_enabled": true,
                    "next_dnd_start_ts": 1450387800,
                    "next_dnd_end_ts": 1450423800
            ***REMOVED***,
                "U058CJVAA": ***REMOVED***
                    "dnd_enabled": false,
                    "next_dnd_start_ts": 1,
                    "next_dnd_end_ts": 1
            ***REMOVED***
        ***REMOVED***
    ***REMOVED***`))
	***REMOVED***)
	usersDNDInfo := map[string]DNDStatus***REMOVED***
		"U023BECGF": DNDStatus***REMOVED***
			Enabled:            true,
			NextStartTimestamp: 1450387800,
			NextEndTimestamp:   1450423800,
		***REMOVED***,
		"U058CJVAA": DNDStatus***REMOVED***
			Enabled:            false,
			NextStartTimestamp: 1,
			NextEndTimestamp:   1,
		***REMOVED***,
	***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	usersDNDInfoResponse, err := api.GetDNDTeamInfo(nil)
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
	eq := reflect.DeepEqual(usersDNDInfoResponse, usersDNDInfo)
	if !eq ***REMOVED***
		t.Errorf("got %v; want %v", usersDNDInfoResponse, usersDNDInfo)
	***REMOVED***
***REMOVED***

func TestSlack_SetSnooze(t *testing.T) ***REMOVED***
	http.HandleFunc("/dnd.setSnooze", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`***REMOVED***
            "ok": true,
            "dnd_enabled": true,
            "snooze_endtime": 1450373897,
            "snooze_remaining": 60
    ***REMOVED***`))
	***REMOVED***)
	snooze := DNDStatus***REMOVED***
		Enabled: true,
		SnoozeInfo: SnoozeInfo***REMOVED***
			SnoozeEndTime:   1450373897,
			SnoozeRemaining: 60,
		***REMOVED***,
	***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")
	snoozeResponse, err := api.SetSnooze(60)
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %s", err)
	***REMOVED***
	eq := reflect.DeepEqual(snoozeResponse, &snooze)
	if !eq ***REMOVED***
		t.Errorf("got %v; want %v", snoozeResponse, &snooze)
	***REMOVED***
***REMOVED***
