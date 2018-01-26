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
	"fmt"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/pkg/types"
)

var (
	defaultV2MembersPrefix = "/v2/members"
	defaultLeaderSuffix    = "/leader"
)

type Member struct ***REMOVED***
	// ID is the unique identifier of this Member.
	ID string `json:"id"`

	// Name is a human-readable, non-unique identifier of this Member.
	Name string `json:"name"`

	// PeerURLs represents the HTTP(S) endpoints this Member uses to
	// participate in etcd's consensus protocol.
	PeerURLs []string `json:"peerURLs"`

	// ClientURLs represents the HTTP(S) endpoints on which this Member
	// serves it's client-facing APIs.
	ClientURLs []string `json:"clientURLs"`
***REMOVED***

type memberCollection []Member

func (c *memberCollection) UnmarshalJSON(data []byte) error ***REMOVED***
	d := struct ***REMOVED***
		Members []Member
	***REMOVED******REMOVED******REMOVED***

	if err := json.Unmarshal(data, &d); err != nil ***REMOVED***
		return err
	***REMOVED***

	if d.Members == nil ***REMOVED***
		*c = make([]Member, 0)
		return nil
	***REMOVED***

	*c = d.Members
	return nil
***REMOVED***

type memberCreateOrUpdateRequest struct ***REMOVED***
	PeerURLs types.URLs
***REMOVED***

func (m *memberCreateOrUpdateRequest) MarshalJSON() ([]byte, error) ***REMOVED***
	s := struct ***REMOVED***
		PeerURLs []string `json:"peerURLs"`
	***REMOVED******REMOVED***
		PeerURLs: make([]string, len(m.PeerURLs)),
	***REMOVED***

	for i, u := range m.PeerURLs ***REMOVED***
		s.PeerURLs[i] = u.String()
	***REMOVED***

	return json.Marshal(&s)
***REMOVED***

// NewMembersAPI constructs a new MembersAPI that uses HTTP to
// interact with etcd's membership API.
func NewMembersAPI(c Client) MembersAPI ***REMOVED***
	return &httpMembersAPI***REMOVED***
		client: c,
	***REMOVED***
***REMOVED***

type MembersAPI interface ***REMOVED***
	// List enumerates the current cluster membership.
	List(ctx context.Context) ([]Member, error)

	// Add instructs etcd to accept a new Member into the cluster.
	Add(ctx context.Context, peerURL string) (*Member, error)

	// Remove demotes an existing Member out of the cluster.
	Remove(ctx context.Context, mID string) error

	// Update instructs etcd to update an existing Member in the cluster.
	Update(ctx context.Context, mID string, peerURLs []string) error

	// Leader gets current leader of the cluster
	Leader(ctx context.Context) (*Member, error)
***REMOVED***

type httpMembersAPI struct ***REMOVED***
	client httpClient
***REMOVED***

func (m *httpMembersAPI) List(ctx context.Context) ([]Member, error) ***REMOVED***
	req := &membersAPIActionList***REMOVED******REMOVED***
	resp, body, err := m.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := assertStatusCode(resp.StatusCode, http.StatusOK); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var mCollection memberCollection
	if err := json.Unmarshal(body, &mCollection); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return []Member(mCollection), nil
***REMOVED***

func (m *httpMembersAPI) Add(ctx context.Context, peerURL string) (*Member, error) ***REMOVED***
	urls, err := types.NewURLs([]string***REMOVED***peerURL***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req := &membersAPIActionAdd***REMOVED***peerURLs: urls***REMOVED***
	resp, body, err := m.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := assertStatusCode(resp.StatusCode, http.StatusCreated, http.StatusConflict); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if resp.StatusCode != http.StatusCreated ***REMOVED***
		var merr membersError
		if err := json.Unmarshal(body, &merr); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, merr
	***REMOVED***

	var memb Member
	if err := json.Unmarshal(body, &memb); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &memb, nil
***REMOVED***

func (m *httpMembersAPI) Update(ctx context.Context, memberID string, peerURLs []string) error ***REMOVED***
	urls, err := types.NewURLs(peerURLs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	req := &membersAPIActionUpdate***REMOVED***peerURLs: urls, memberID: memberID***REMOVED***
	resp, body, err := m.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := assertStatusCode(resp.StatusCode, http.StatusNoContent, http.StatusNotFound, http.StatusConflict); err != nil ***REMOVED***
		return err
	***REMOVED***

	if resp.StatusCode != http.StatusNoContent ***REMOVED***
		var merr membersError
		if err := json.Unmarshal(body, &merr); err != nil ***REMOVED***
			return err
		***REMOVED***
		return merr
	***REMOVED***

	return nil
***REMOVED***

func (m *httpMembersAPI) Remove(ctx context.Context, memberID string) error ***REMOVED***
	req := &membersAPIActionRemove***REMOVED***memberID: memberID***REMOVED***
	resp, _, err := m.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return assertStatusCode(resp.StatusCode, http.StatusNoContent, http.StatusGone)
***REMOVED***

func (m *httpMembersAPI) Leader(ctx context.Context) (*Member, error) ***REMOVED***
	req := &membersAPIActionLeader***REMOVED******REMOVED***
	resp, body, err := m.client.Do(ctx, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := assertStatusCode(resp.StatusCode, http.StatusOK); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var leader Member
	if err := json.Unmarshal(body, &leader); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &leader, nil
***REMOVED***

type membersAPIActionList struct***REMOVED******REMOVED***

func (l *membersAPIActionList) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2MembersURL(ep)
	req, _ := http.NewRequest("GET", u.String(), nil)
	return req
***REMOVED***

type membersAPIActionRemove struct ***REMOVED***
	memberID string
***REMOVED***

func (d *membersAPIActionRemove) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2MembersURL(ep)
	u.Path = path.Join(u.Path, d.memberID)
	req, _ := http.NewRequest("DELETE", u.String(), nil)
	return req
***REMOVED***

type membersAPIActionAdd struct ***REMOVED***
	peerURLs types.URLs
***REMOVED***

func (a *membersAPIActionAdd) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2MembersURL(ep)
	m := memberCreateOrUpdateRequest***REMOVED***PeerURLs: a.peerURLs***REMOVED***
	b, _ := json.Marshal(&m)
	req, _ := http.NewRequest("POST", u.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return req
***REMOVED***

type membersAPIActionUpdate struct ***REMOVED***
	memberID string
	peerURLs types.URLs
***REMOVED***

func (a *membersAPIActionUpdate) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2MembersURL(ep)
	m := memberCreateOrUpdateRequest***REMOVED***PeerURLs: a.peerURLs***REMOVED***
	u.Path = path.Join(u.Path, a.memberID)
	b, _ := json.Marshal(&m)
	req, _ := http.NewRequest("PUT", u.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return req
***REMOVED***

func assertStatusCode(got int, want ...int) (err error) ***REMOVED***
	for _, w := range want ***REMOVED***
		if w == got ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("unexpected status code %d", got)
***REMOVED***

type membersAPIActionLeader struct***REMOVED******REMOVED***

func (l *membersAPIActionLeader) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2MembersURL(ep)
	u.Path = path.Join(u.Path, defaultLeaderSuffix)
	req, _ := http.NewRequest("GET", u.String(), nil)
	return req
***REMOVED***

// v2MembersURL add the necessary path to the provided endpoint
// to route requests to the default v2 members API.
func v2MembersURL(ep url.URL) *url.URL ***REMOVED***
	ep.Path = path.Join(ep.Path, defaultV2MembersPrefix)
	return &ep
***REMOVED***

type membersError struct ***REMOVED***
	Message string `json:"message"`
	Code    int    `json:"-"`
***REMOVED***

func (e membersError) Error() string ***REMOVED***
	return e.Message
***REMOVED***
