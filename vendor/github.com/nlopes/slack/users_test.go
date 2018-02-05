package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func getTestUserProfile() UserProfile ***REMOVED***
	return UserProfile***REMOVED***
		StatusText:            "testStatus",
		StatusEmoji:           ":construction:",
		RealName:              "Test Real Name",
		RealNameNormalized:    "Test Real Name Normalized",
		DisplayName:           "Test Display Name",
		DisplayNameNormalized: "Test Display Name Normalized",
		Email:    "test@test.com",
		Image24:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_24.jpg",
		Image32:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_32.jpg",
		Image48:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_48.jpg",
		Image72:  "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_72.jpg",
		Image192: "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_192.jpg",
	***REMOVED***
***REMOVED***

func getTestUser() User ***REMOVED***
	return User***REMOVED***
		ID:                "UXXXXXXXX",
		Name:              "Test User",
		Deleted:           false,
		Color:             "9f69e7",
		RealName:          "testuser",
		TZ:                "America/Los_Angeles",
		TZLabel:           "Pacific Daylight Time",
		TZOffset:          -25200,
		Profile:           getTestUserProfile(),
		IsBot:             false,
		IsAdmin:           false,
		IsOwner:           false,
		IsPrimaryOwner:    false,
		IsRestricted:      false,
		IsUltraRestricted: false,
		Has2FA:            false,
	***REMOVED***
***REMOVED***

func getUserIdentity(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`***REMOVED***
  "ok": true,
  "user": ***REMOVED***
    "id": "UXXXXXXXX",
    "name": "Test User",
    "email": "test@test.com",
    "image_24": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_24.jpg",
    "image_32": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_32.jpg",
    "image_48": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_48.jpg",
    "image_72": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_72.jpg",
    "image_192": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_192.jpg",
    "image_512": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_512.jpg"
  ***REMOVED***,
  "team": ***REMOVED***
    "id": "TXXXXXXXX",
    "name": "team-name",
    "domain": "team-domain",
    "image_34": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_34.jpg",
    "image_44": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_44.jpg",
    "image_68": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_68.jpg",
    "image_88": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_88.jpg",
    "image_102": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_102.jpg",
    "image_132": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_132.jpg",
    "image_230": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_230.jpg",
    "image_original": "https:\/\/s3-us-west-2.amazonaws.com\/slack-files2\/avatars\/2016-10-18\/92962080834_ef14c1469fc0741caea1_original.jpg"
  ***REMOVED***
***REMOVED***`)
	rw.Write(response)
***REMOVED***

func getUserByEmail(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct ***REMOVED***
		Ok   bool
		User User
	***REMOVED******REMOVED***
		Ok:   true,
		User: getTestUser(),
	***REMOVED***)
	rw.Write(response)
***REMOVED***

func httpTestErrReply(w http.ResponseWriter, clientErr bool, msg string) ***REMOVED***
	if clientErr ***REMOVED***
		w.WriteHeader(http.StatusBadRequest)
	***REMOVED*** else ***REMOVED***
		w.WriteHeader(http.StatusInternalServerError)
	***REMOVED***

	w.Header().Set("Content-Type", "application/json")

	body, _ := json.Marshal(&SlackResponse***REMOVED***
		Ok: false, Error: msg,
	***REMOVED***)

	w.Write(body)
***REMOVED***

func newProfileHandler(up *UserProfile) (setter func(http.ResponseWriter, *http.Request)) ***REMOVED***
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if up == nil ***REMOVED***
			httpTestErrReply(w, false, "err: UserProfile is nil")
			return
		***REMOVED***

		if err := r.ParseForm(); err != nil ***REMOVED***
			httpTestErrReply(w, true, fmt.Sprintf("err parsing form: %s", err.Error()))
			return
		***REMOVED***

		values := r.Form

		if len(values["profile"]) == 0 ***REMOVED***
			httpTestErrReply(w, true, `POST data must include a "profile" field`)
			return
		***REMOVED***

		profile := []byte(values["profile"][0])

		userProfile := UserProfile***REMOVED******REMOVED***

		if err := json.Unmarshal(profile, &userProfile); err != nil ***REMOVED***
			httpTestErrReply(w, true, fmt.Sprintf("err parsing JSON: %s\n\njson: `%s`", err.Error(), profile))
			return
		***REMOVED***

		*up = userProfile

		// TODO(theckman): enhance this to return a full User object
		fmt.Fprint(w, `***REMOVED***"ok":true***REMOVED***`)
	***REMOVED***
***REMOVED***

func TestGetUserIdentity(t *testing.T) ***REMOVED***
	http.HandleFunc("/users.identity", getUserIdentity)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	identity, err := api.GetUserIdentity()
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***

	// t.Fatal refers to -> t.Errorf & return
	if identity.User.ID != "UXXXXXXXX" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if identity.User.Name != "Test User" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if identity.User.Email != "test@test.com" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if identity.Team.ID != "TXXXXXXXX" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if identity.Team.Name != "team-name" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if identity.Team.Domain != "team-domain" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if identity.User.Image24 == "" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
	if identity.Team.Image34 == "" ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
***REMOVED***

func TestGetUserByEmail(t *testing.T) ***REMOVED***
	http.HandleFunc("/users.lookupByEmail", getUserByEmail)
	expectedUser := getTestUser()

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	user, err := api.GetUserByEmail("test@test.com")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
		return
	***REMOVED***
	if !reflect.DeepEqual(expectedUser, *user) ***REMOVED***
		t.Fatal(ErrIncorrectResponse)
	***REMOVED***
***REMOVED***

func TestUserCustomStatus(t *testing.T) ***REMOVED***
	up := &UserProfile***REMOVED******REMOVED***

	setUserProfile := newProfileHandler(up)

	http.HandleFunc("/users.profile.set", setUserProfile)

	once.Do(startServer)
	SLACK_API = "http://" + serverAddr + "/"
	api := New("testing-token")

	testSetUserCustomStatus(api, up, t)
	testUnsetUserCustomStatus(api, up, t)
***REMOVED***

func testSetUserCustomStatus(api *Client, up *UserProfile, t *testing.T) ***REMOVED***
	const (
		statusText  = "testStatus"
		statusEmoji = ":construction:"
	)

	if err := api.SetUserCustomStatus(statusText, statusEmoji); err != nil ***REMOVED***
		t.Fatalf(`SetUserCustomStatus(%q, %q) = %#v, want <nil>`, statusText, statusEmoji, err)
	***REMOVED***

	if up.StatusText != statusText ***REMOVED***
		t.Fatalf(`UserProfile.StatusText = %q, want %q`, up.StatusText, statusText)
	***REMOVED***

	if up.StatusEmoji != statusEmoji ***REMOVED***
		t.Fatalf(`UserProfile.StatusEmoji = %q, want %q`, up.StatusEmoji, statusEmoji)
	***REMOVED***
***REMOVED***

func testUnsetUserCustomStatus(api *Client, up *UserProfile, t *testing.T) ***REMOVED***
	if err := api.UnsetUserCustomStatus(); err != nil ***REMOVED***
		t.Fatalf(`UnsetUserCustomStatus() = %#v, want <nil>`, err)
	***REMOVED***

	if up.StatusText != "" ***REMOVED***
		t.Fatalf(`UserProfile.StatusText = %q, want %q`, up.StatusText, "")
	***REMOVED***

	if up.StatusEmoji != "" ***REMOVED***
		t.Fatalf(`UserProfile.StatusEmoji = %q, want %q`, up.StatusEmoji, "")
	***REMOVED***
***REMOVED***
