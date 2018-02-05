package slack

import (
	"net/http"
	"reflect"
	"testing"
)

type userGroupsHandler struct ***REMOVED***
	gotParams map[string]string
	response  string
***REMOVED***

func newUserGroupsHandler() *userGroupsHandler ***REMOVED***
	return &userGroupsHandler***REMOVED***
		gotParams: make(map[string]string),
		response: `***REMOVED***
    "ok": true,
    "usergroup": ***REMOVED***
        "id": "S0615G0KT",
        "team_id": "T060RNRCH",
        "is_usergroup": true,
        "name": "Marketing Team",
        "description": "Marketing gurus, PR experts and product advocates.",
        "handle": "marketing-team",
        "is_external": false,
        "date_create": 1446746793,
        "date_update": 1446746793,
        "date_delete": 0,
        "auto_type": null,
        "created_by": "U060RNRCZ",
        "updated_by": "U060RNRCZ",
        "deleted_by": null,
        "prefs": ***REMOVED***
            "channels": [

            ],
            "groups": [

            ]
    ***REMOVED***,
        "user_count": 0
***REMOVED***
***REMOVED***`,
	***REMOVED***
***REMOVED***

func (ugh *userGroupsHandler) accumulateFormValue(k string, r *http.Request) ***REMOVED***
	if v := r.FormValue(k); v != "" ***REMOVED***
		ugh.gotParams[k] = v
	***REMOVED***
***REMOVED***

func (ugh *userGroupsHandler) handler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	ugh.accumulateFormValue("name", r)
	ugh.accumulateFormValue("description", r)
	ugh.accumulateFormValue("handle", r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ugh.response))
***REMOVED***

func TestCreateUserGroup(t *testing.T) ***REMOVED***
	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	tests := []struct ***REMOVED***
		userGroup  UserGroup
		wantParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			UserGroup***REMOVED***
				Name:        "Marketing Team",
				Description: "Marketing gurus, PR experts and product advocates.",
				Handle:      "marketing-team"***REMOVED***,
			map[string]string***REMOVED***
				"name":        "Marketing Team",
				"description": "Marketing gurus, PR experts and product advocates.",
				"handle":      "marketing-team",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	var rh *userGroupsHandler
	http.HandleFunc("/usergroups.create", func(w http.ResponseWriter, r *http.Request) ***REMOVED*** rh.handler(w, r) ***REMOVED***)

	for i, test := range tests ***REMOVED***
		rh = newUserGroupsHandler()
		_, err := api.CreateUserGroup(test.userGroup)
		if err != nil ***REMOVED***
			t.Fatalf("%d: Unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(rh.gotParams, test.wantParams) ***REMOVED***
			t.Errorf("%d: Got params %#v, want %#v", i, rh.gotParams, test.wantParams)
		***REMOVED***
	***REMOVED***
***REMOVED***

func getUserGroups(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`***REMOVED***
    "ok": true,
    "usergroups": [
        ***REMOVED***
            "id": "S0614TZR7",
            "team_id": "T060RNRCH",
            "is_usergroup": true,
            "name": "Team Admins",
            "description": "A group of all Administrators on your team.",
            "handle": "admins",
            "is_external": false,
            "date_create": 1446598059,
            "date_update": 1446670362,
            "date_delete": 0,
            "auto_type": "admin",
            "created_by": "USLACKBOT",
            "updated_by": "U060RNRCZ",
            "deleted_by": null,
            "prefs": ***REMOVED***
                "channels": [
                  "channel1",
                  "channel2"
                ],
                "groups": [
                  "group1",
                  "group2",
                  "group3"
                ]
        ***REMOVED***,
            "user_count": 2
    ***REMOVED***
    ]
***REMOVED***`)
	rw.Write(response)
***REMOVED***

func TestGetUserGroups(t *testing.T) ***REMOVED***
	http.HandleFunc("/usergroups.list", getUserGroups)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	userGroups, err := api.GetUserGroups()
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***

	// t.Fatal refers to -> t.Errorf & return
	if len(userGroups) != 1 ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***

	S0614TZR7 := UserGroup***REMOVED***
		ID:          "S0614TZR7",
		TeamID:      "T060RNRCH",
		IsUserGroup: true,
		Name:        "Team Admins",
		Description: "A group of all Administrators on your team.",
		Handle:      "admins",
		IsExternal:  false,
		DateCreate:  1446598059,
		DateUpdate:  1446670362,
		DateDelete:  0,
		AutoType:    "admin",
		CreatedBy:   "USLACKBOT",
		UpdatedBy:   "U060RNRCZ",
		DeletedBy:   "",
		Prefs: UserGroupPrefs***REMOVED***
			Channels: []string***REMOVED***"channel1", "channel2"***REMOVED***,
			Groups:   []string***REMOVED***"group1", "group2", "group3"***REMOVED***,
		***REMOVED***,
		UserCount: 2,
	***REMOVED***

	if !reflect.DeepEqual(userGroups[0], S0614TZR7) ***REMOVED***
		t.Errorf("Got %#v, want %#v", userGroups[0], S0614TZR7)
	***REMOVED***
***REMOVED***
