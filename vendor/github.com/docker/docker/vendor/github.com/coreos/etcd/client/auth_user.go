// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/net/context"
)

var (
	defaultV2AuthPrefix = "/v2/auth"
)

type User struct ***REMOVED***
	User     string   `json:"user"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles"`
	Grant    []string `json:"grant,omitempty"`
	Revoke   []string `json:"revoke,omitempty"`
***REMOVED***

// userListEntry is the user representation given by the server for ListUsers
type userListEntry struct ***REMOVED***
	User  string `json:"user"`
	Roles []Role `json:"roles"`
***REMOVED***

type UserRoles struct ***REMOVED***
	User  string `json:"user"`
	Roles []Role `json:"roles"`
***REMOVED***

func v2AuthURL(ep url.URL, action string, name string) *url.URL ***REMOVED***
	if name != "" ***REMOVED***
		ep.Path = path.Join(ep.Path, defaultV2AuthPrefix, action, name)
		return &ep
	***REMOVED***
	ep.Path = path.Join(ep.Path, defaultV2AuthPrefix, action)
	return &ep
***REMOVED***

// NewAuthAPI constructs a new AuthAPI that uses HTTP to
// interact with etcd's general auth features.
func NewAuthAPI(c Client) AuthAPI ***REMOVED***
	return &httpAuthAPI***REMOVED***
		client: c,
	***REMOVED***
***REMOVED***

type AuthAPI interface ***REMOVED***
	// Enable auth.
	Enable(ctx context.Context) error

	// Disable auth.
	Disable(ctx context.Context) error
***REMOVED***

type httpAuthAPI struct ***REMOVED***
	client httpClient
***REMOVED***

func (s *httpAuthAPI) Enable(ctx context.Context) error ***REMOVED***
	return s.enableDisable(ctx, &authAPIAction***REMOVED***"PUT"***REMOVED***)
***REMOVED***

func (s *httpAuthAPI) Disable(ctx context.Context) error ***REMOVED***
	return s.enableDisable(ctx, &authAPIAction***REMOVED***"DELETE"***REMOVED***)
***REMOVED***

func (s *httpAuthAPI) enableDisable(ctx context.Context, req httpAction) error ***REMOVED***
	resp, body, err := s.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = assertStatusCode(resp.StatusCode, http.StatusOK, http.StatusCreated); err != nil ***REMOVED***
		var sec authError
		err = json.Unmarshal(body, &sec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return sec
	***REMOVED***
	return nil
***REMOVED***

type authAPIAction struct ***REMOVED***
	verb string
***REMOVED***

func (l *authAPIAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2AuthURL(ep, "enable", "")
	req, _ := http.NewRequest(l.verb, u.String(), nil)
	return req
***REMOVED***

type authError struct ***REMOVED***
	Message string `json:"message"`
	Code    int    `json:"-"`
***REMOVED***

func (e authError) Error() string ***REMOVED***
	return e.Message
***REMOVED***

// NewAuthUserAPI constructs a new AuthUserAPI that uses HTTP to
// interact with etcd's user creation and modification features.
func NewAuthUserAPI(c Client) AuthUserAPI ***REMOVED***
	return &httpAuthUserAPI***REMOVED***
		client: c,
	***REMOVED***
***REMOVED***

type AuthUserAPI interface ***REMOVED***
	// AddUser adds a user.
	AddUser(ctx context.Context, username string, password string) error

	// RemoveUser removes a user.
	RemoveUser(ctx context.Context, username string) error

	// GetUser retrieves user details.
	GetUser(ctx context.Context, username string) (*User, error)

	// GrantUser grants a user some permission roles.
	GrantUser(ctx context.Context, username string, roles []string) (*User, error)

	// RevokeUser revokes some permission roles from a user.
	RevokeUser(ctx context.Context, username string, roles []string) (*User, error)

	// ChangePassword changes the user's password.
	ChangePassword(ctx context.Context, username string, password string) (*User, error)

	// ListUsers lists the users.
	ListUsers(ctx context.Context) ([]string, error)
***REMOVED***

type httpAuthUserAPI struct ***REMOVED***
	client httpClient
***REMOVED***

type authUserAPIAction struct ***REMOVED***
	verb     string
	username string
	user     *User
***REMOVED***

type authUserAPIList struct***REMOVED******REMOVED***

func (list *authUserAPIList) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2AuthURL(ep, "users", "")
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	return req
***REMOVED***

func (l *authUserAPIAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2AuthURL(ep, "users", l.username)
	if l.user == nil ***REMOVED***
		req, _ := http.NewRequest(l.verb, u.String(), nil)
		return req
	***REMOVED***
	b, err := json.Marshal(l.user)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	body := bytes.NewReader(b)
	req, _ := http.NewRequest(l.verb, u.String(), body)
	req.Header.Set("Content-Type", "application/json")
	return req
***REMOVED***

func (u *httpAuthUserAPI) ListUsers(ctx context.Context) ([]string, error) ***REMOVED***
	resp, body, err := u.client.Do(ctx, &authUserAPIList***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = assertStatusCode(resp.StatusCode, http.StatusOK); err != nil ***REMOVED***
		var sec authError
		err = json.Unmarshal(body, &sec)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, sec
	***REMOVED***

	var userList struct ***REMOVED***
		Users []userListEntry `json:"users"`
	***REMOVED***

	if err = json.Unmarshal(body, &userList); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ret := make([]string, 0, len(userList.Users))
	for _, u := range userList.Users ***REMOVED***
		ret = append(ret, u.User)
	***REMOVED***
	return ret, nil
***REMOVED***

func (u *httpAuthUserAPI) AddUser(ctx context.Context, username string, password string) error ***REMOVED***
	user := &User***REMOVED***
		User:     username,
		Password: password,
	***REMOVED***
	return u.addRemoveUser(ctx, &authUserAPIAction***REMOVED***
		verb:     "PUT",
		username: username,
		user:     user,
	***REMOVED***)
***REMOVED***

func (u *httpAuthUserAPI) RemoveUser(ctx context.Context, username string) error ***REMOVED***
	return u.addRemoveUser(ctx, &authUserAPIAction***REMOVED***
		verb:     "DELETE",
		username: username,
	***REMOVED***)
***REMOVED***

func (u *httpAuthUserAPI) addRemoveUser(ctx context.Context, req *authUserAPIAction) error ***REMOVED***
	resp, body, err := u.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = assertStatusCode(resp.StatusCode, http.StatusOK, http.StatusCreated); err != nil ***REMOVED***
		var sec authError
		err = json.Unmarshal(body, &sec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return sec
	***REMOVED***
	return nil
***REMOVED***

func (u *httpAuthUserAPI) GetUser(ctx context.Context, username string) (*User, error) ***REMOVED***
	return u.modUser(ctx, &authUserAPIAction***REMOVED***
		verb:     "GET",
		username: username,
	***REMOVED***)
***REMOVED***

func (u *httpAuthUserAPI) GrantUser(ctx context.Context, username string, roles []string) (*User, error) ***REMOVED***
	user := &User***REMOVED***
		User:  username,
		Grant: roles,
	***REMOVED***
	return u.modUser(ctx, &authUserAPIAction***REMOVED***
		verb:     "PUT",
		username: username,
		user:     user,
	***REMOVED***)
***REMOVED***

func (u *httpAuthUserAPI) RevokeUser(ctx context.Context, username string, roles []string) (*User, error) ***REMOVED***
	user := &User***REMOVED***
		User:   username,
		Revoke: roles,
	***REMOVED***
	return u.modUser(ctx, &authUserAPIAction***REMOVED***
		verb:     "PUT",
		username: username,
		user:     user,
	***REMOVED***)
***REMOVED***

func (u *httpAuthUserAPI) ChangePassword(ctx context.Context, username string, password string) (*User, error) ***REMOVED***
	user := &User***REMOVED***
		User:     username,
		Password: password,
	***REMOVED***
	return u.modUser(ctx, &authUserAPIAction***REMOVED***
		verb:     "PUT",
		username: username,
		user:     user,
	***REMOVED***)
***REMOVED***

func (u *httpAuthUserAPI) modUser(ctx context.Context, req *authUserAPIAction) (*User, error) ***REMOVED***
	resp, body, err := u.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = assertStatusCode(resp.StatusCode, http.StatusOK); err != nil ***REMOVED***
		var sec authError
		err = json.Unmarshal(body, &sec)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, sec
	***REMOVED***
	var user User
	if err = json.Unmarshal(body, &user); err != nil ***REMOVED***
		var userR UserRoles
		if urerr := json.Unmarshal(body, &userR); urerr != nil ***REMOVED***
			return nil, err
		***REMOVED***
		user.User = userR.User
		for _, r := range userR.Roles ***REMOVED***
			user.Roles = append(user.Roles, r.Role)
		***REMOVED***
	***REMOVED***
	return &user, nil
***REMOVED***
