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

	"golang.org/x/net/context"
)

type Role struct ***REMOVED***
	Role        string       `json:"role"`
	Permissions Permissions  `json:"permissions"`
	Grant       *Permissions `json:"grant,omitempty"`
	Revoke      *Permissions `json:"revoke,omitempty"`
***REMOVED***

type Permissions struct ***REMOVED***
	KV rwPermission `json:"kv"`
***REMOVED***

type rwPermission struct ***REMOVED***
	Read  []string `json:"read"`
	Write []string `json:"write"`
***REMOVED***

type PermissionType int

const (
	ReadPermission PermissionType = iota
	WritePermission
	ReadWritePermission
)

// NewAuthRoleAPI constructs a new AuthRoleAPI that uses HTTP to
// interact with etcd's role creation and modification features.
func NewAuthRoleAPI(c Client) AuthRoleAPI ***REMOVED***
	return &httpAuthRoleAPI***REMOVED***
		client: c,
	***REMOVED***
***REMOVED***

type AuthRoleAPI interface ***REMOVED***
	// AddRole adds a role.
	AddRole(ctx context.Context, role string) error

	// RemoveRole removes a role.
	RemoveRole(ctx context.Context, role string) error

	// GetRole retrieves role details.
	GetRole(ctx context.Context, role string) (*Role, error)

	// GrantRoleKV grants a role some permission prefixes for the KV store.
	GrantRoleKV(ctx context.Context, role string, prefixes []string, permType PermissionType) (*Role, error)

	// RevokeRoleKV revokes some permission prefixes for a role on the KV store.
	RevokeRoleKV(ctx context.Context, role string, prefixes []string, permType PermissionType) (*Role, error)

	// ListRoles lists roles.
	ListRoles(ctx context.Context) ([]string, error)
***REMOVED***

type httpAuthRoleAPI struct ***REMOVED***
	client httpClient
***REMOVED***

type authRoleAPIAction struct ***REMOVED***
	verb string
	name string
	role *Role
***REMOVED***

type authRoleAPIList struct***REMOVED******REMOVED***

func (list *authRoleAPIList) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2AuthURL(ep, "roles", "")
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	return req
***REMOVED***

func (l *authRoleAPIAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2AuthURL(ep, "roles", l.name)
	if l.role == nil ***REMOVED***
		req, _ := http.NewRequest(l.verb, u.String(), nil)
		return req
	***REMOVED***
	b, err := json.Marshal(l.role)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	body := bytes.NewReader(b)
	req, _ := http.NewRequest(l.verb, u.String(), body)
	req.Header.Set("Content-Type", "application/json")
	return req
***REMOVED***

func (r *httpAuthRoleAPI) ListRoles(ctx context.Context) ([]string, error) ***REMOVED***
	resp, body, err := r.client.Do(ctx, &authRoleAPIList***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = assertStatusCode(resp.StatusCode, http.StatusOK); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var roleList struct ***REMOVED***
		Roles []Role `json:"roles"`
	***REMOVED***
	if err = json.Unmarshal(body, &roleList); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ret := make([]string, 0, len(roleList.Roles))
	for _, r := range roleList.Roles ***REMOVED***
		ret = append(ret, r.Role)
	***REMOVED***
	return ret, nil
***REMOVED***

func (r *httpAuthRoleAPI) AddRole(ctx context.Context, rolename string) error ***REMOVED***
	role := &Role***REMOVED***
		Role: rolename,
	***REMOVED***
	return r.addRemoveRole(ctx, &authRoleAPIAction***REMOVED***
		verb: "PUT",
		name: rolename,
		role: role,
	***REMOVED***)
***REMOVED***

func (r *httpAuthRoleAPI) RemoveRole(ctx context.Context, rolename string) error ***REMOVED***
	return r.addRemoveRole(ctx, &authRoleAPIAction***REMOVED***
		verb: "DELETE",
		name: rolename,
	***REMOVED***)
***REMOVED***

func (r *httpAuthRoleAPI) addRemoveRole(ctx context.Context, req *authRoleAPIAction) error ***REMOVED***
	resp, body, err := r.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := assertStatusCode(resp.StatusCode, http.StatusOK, http.StatusCreated); err != nil ***REMOVED***
		var sec authError
		err := json.Unmarshal(body, &sec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return sec
	***REMOVED***
	return nil
***REMOVED***

func (r *httpAuthRoleAPI) GetRole(ctx context.Context, rolename string) (*Role, error) ***REMOVED***
	return r.modRole(ctx, &authRoleAPIAction***REMOVED***
		verb: "GET",
		name: rolename,
	***REMOVED***)
***REMOVED***

func buildRWPermission(prefixes []string, permType PermissionType) rwPermission ***REMOVED***
	var out rwPermission
	switch permType ***REMOVED***
	case ReadPermission:
		out.Read = prefixes
	case WritePermission:
		out.Write = prefixes
	case ReadWritePermission:
		out.Read = prefixes
		out.Write = prefixes
	***REMOVED***
	return out
***REMOVED***

func (r *httpAuthRoleAPI) GrantRoleKV(ctx context.Context, rolename string, prefixes []string, permType PermissionType) (*Role, error) ***REMOVED***
	rwp := buildRWPermission(prefixes, permType)
	role := &Role***REMOVED***
		Role: rolename,
		Grant: &Permissions***REMOVED***
			KV: rwp,
		***REMOVED***,
	***REMOVED***
	return r.modRole(ctx, &authRoleAPIAction***REMOVED***
		verb: "PUT",
		name: rolename,
		role: role,
	***REMOVED***)
***REMOVED***

func (r *httpAuthRoleAPI) RevokeRoleKV(ctx context.Context, rolename string, prefixes []string, permType PermissionType) (*Role, error) ***REMOVED***
	rwp := buildRWPermission(prefixes, permType)
	role := &Role***REMOVED***
		Role: rolename,
		Revoke: &Permissions***REMOVED***
			KV: rwp,
		***REMOVED***,
	***REMOVED***
	return r.modRole(ctx, &authRoleAPIAction***REMOVED***
		verb: "PUT",
		name: rolename,
		role: role,
	***REMOVED***)
***REMOVED***

func (r *httpAuthRoleAPI) modRole(ctx context.Context, req *authRoleAPIAction) (*Role, error) ***REMOVED***
	resp, body, err := r.client.Do(ctx, req)
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
	var role Role
	if err = json.Unmarshal(body, &role); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &role, nil
***REMOVED***
