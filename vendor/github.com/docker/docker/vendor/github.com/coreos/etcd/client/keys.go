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

//go:generate codecgen -d 1819 -r "Node|Response|Nodes" -o keys.generated.go keys.go

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/pkg/pathutil"
	"github.com/ugorji/go/codec"
	"golang.org/x/net/context"
)

const (
	ErrorCodeKeyNotFound  = 100
	ErrorCodeTestFailed   = 101
	ErrorCodeNotFile      = 102
	ErrorCodeNotDir       = 104
	ErrorCodeNodeExist    = 105
	ErrorCodeRootROnly    = 107
	ErrorCodeDirNotEmpty  = 108
	ErrorCodeUnauthorized = 110

	ErrorCodePrevValueRequired = 201
	ErrorCodeTTLNaN            = 202
	ErrorCodeIndexNaN          = 203
	ErrorCodeInvalidField      = 209
	ErrorCodeInvalidForm       = 210

	ErrorCodeRaftInternal = 300
	ErrorCodeLeaderElect  = 301

	ErrorCodeWatcherCleared    = 400
	ErrorCodeEventIndexCleared = 401
)

type Error struct ***REMOVED***
	Code    int    `json:"errorCode"`
	Message string `json:"message"`
	Cause   string `json:"cause"`
	Index   uint64 `json:"index"`
***REMOVED***

func (e Error) Error() string ***REMOVED***
	return fmt.Sprintf("%v: %v (%v) [%v]", e.Code, e.Message, e.Cause, e.Index)
***REMOVED***

var (
	ErrInvalidJSON = errors.New("client: response is invalid json. The endpoint is probably not valid etcd cluster endpoint.")
	ErrEmptyBody   = errors.New("client: response body is empty")
)

// PrevExistType is used to define an existence condition when setting
// or deleting Nodes.
type PrevExistType string

const (
	PrevIgnore  = PrevExistType("")
	PrevExist   = PrevExistType("true")
	PrevNoExist = PrevExistType("false")
)

var (
	defaultV2KeysPrefix = "/v2/keys"
)

// NewKeysAPI builds a KeysAPI that interacts with etcd's key-value
// API over HTTP.
func NewKeysAPI(c Client) KeysAPI ***REMOVED***
	return NewKeysAPIWithPrefix(c, defaultV2KeysPrefix)
***REMOVED***

// NewKeysAPIWithPrefix acts like NewKeysAPI, but allows the caller
// to provide a custom base URL path. This should only be used in
// very rare cases.
func NewKeysAPIWithPrefix(c Client, p string) KeysAPI ***REMOVED***
	return &httpKeysAPI***REMOVED***
		client: c,
		prefix: p,
	***REMOVED***
***REMOVED***

type KeysAPI interface ***REMOVED***
	// Get retrieves a set of Nodes from etcd
	Get(ctx context.Context, key string, opts *GetOptions) (*Response, error)

	// Set assigns a new value to a Node identified by a given key. The caller
	// may define a set of conditions in the SetOptions. If SetOptions.Dir=true
	// then value is ignored.
	Set(ctx context.Context, key, value string, opts *SetOptions) (*Response, error)

	// Delete removes a Node identified by the given key, optionally destroying
	// all of its children as well. The caller may define a set of required
	// conditions in an DeleteOptions object.
	Delete(ctx context.Context, key string, opts *DeleteOptions) (*Response, error)

	// Create is an alias for Set w/ PrevExist=false
	Create(ctx context.Context, key, value string) (*Response, error)

	// CreateInOrder is used to atomically create in-order keys within the given directory.
	CreateInOrder(ctx context.Context, dir, value string, opts *CreateInOrderOptions) (*Response, error)

	// Update is an alias for Set w/ PrevExist=true
	Update(ctx context.Context, key, value string) (*Response, error)

	// Watcher builds a new Watcher targeted at a specific Node identified
	// by the given key. The Watcher may be configured at creation time
	// through a WatcherOptions object. The returned Watcher is designed
	// to emit events that happen to a Node, and optionally to its children.
	Watcher(key string, opts *WatcherOptions) Watcher
***REMOVED***

type WatcherOptions struct ***REMOVED***
	// AfterIndex defines the index after-which the Watcher should
	// start emitting events. For example, if a value of 5 is
	// provided, the first event will have an index >= 6.
	//
	// Setting AfterIndex to 0 (default) means that the Watcher
	// should start watching for events starting at the current
	// index, whatever that may be.
	AfterIndex uint64

	// Recursive specifies whether or not the Watcher should emit
	// events that occur in children of the given keyspace. If set
	// to false (default), events will be limited to those that
	// occur for the exact key.
	Recursive bool
***REMOVED***

type CreateInOrderOptions struct ***REMOVED***
	// TTL defines a period of time after-which the Node should
	// expire and no longer exist. Values <= 0 are ignored. Given
	// that the zero-value is ignored, TTL cannot be used to set
	// a TTL of 0.
	TTL time.Duration
***REMOVED***

type SetOptions struct ***REMOVED***
	// PrevValue specifies what the current value of the Node must
	// be in order for the Set operation to succeed.
	//
	// Leaving this field empty means that the caller wishes to
	// ignore the current value of the Node. This cannot be used
	// to compare the Node's current value to an empty string.
	//
	// PrevValue is ignored if Dir=true
	PrevValue string

	// PrevIndex indicates what the current ModifiedIndex of the
	// Node must be in order for the Set operation to succeed.
	//
	// If PrevIndex is set to 0 (default), no comparison is made.
	PrevIndex uint64

	// PrevExist specifies whether the Node must currently exist
	// (PrevExist) or not (PrevNoExist). If the caller does not
	// care about existence, set PrevExist to PrevIgnore, or simply
	// leave it unset.
	PrevExist PrevExistType

	// TTL defines a period of time after-which the Node should
	// expire and no longer exist. Values <= 0 are ignored. Given
	// that the zero-value is ignored, TTL cannot be used to set
	// a TTL of 0.
	TTL time.Duration

	// Refresh set to true means a TTL value can be updated
	// without firing a watch or changing the node value. A
	// value must not be provided when refreshing a key.
	Refresh bool

	// Dir specifies whether or not this Node should be created as a directory.
	Dir bool

	// NoValueOnSuccess specifies whether the response contains the current value of the Node.
	// If set, the response will only contain the current value when the request fails.
	NoValueOnSuccess bool
***REMOVED***

type GetOptions struct ***REMOVED***
	// Recursive defines whether or not all children of the Node
	// should be returned.
	Recursive bool

	// Sort instructs the server whether or not to sort the Nodes.
	// If true, the Nodes are sorted alphabetically by key in
	// ascending order (A to z). If false (default), the Nodes will
	// not be sorted and the ordering used should not be considered
	// predictable.
	Sort bool

	// Quorum specifies whether it gets the latest committed value that
	// has been applied in quorum of members, which ensures external
	// consistency (or linearizability).
	Quorum bool
***REMOVED***

type DeleteOptions struct ***REMOVED***
	// PrevValue specifies what the current value of the Node must
	// be in order for the Delete operation to succeed.
	//
	// Leaving this field empty means that the caller wishes to
	// ignore the current value of the Node. This cannot be used
	// to compare the Node's current value to an empty string.
	PrevValue string

	// PrevIndex indicates what the current ModifiedIndex of the
	// Node must be in order for the Delete operation to succeed.
	//
	// If PrevIndex is set to 0 (default), no comparison is made.
	PrevIndex uint64

	// Recursive defines whether or not all children of the Node
	// should be deleted. If set to true, all children of the Node
	// identified by the given key will be deleted. If left unset
	// or explicitly set to false, only a single Node will be
	// deleted.
	Recursive bool

	// Dir specifies whether or not this Node should be removed as a directory.
	Dir bool
***REMOVED***

type Watcher interface ***REMOVED***
	// Next blocks until an etcd event occurs, then returns a Response
	// representing that event. The behavior of Next depends on the
	// WatcherOptions used to construct the Watcher. Next is designed to
	// be called repeatedly, each time blocking until a subsequent event
	// is available.
	//
	// If the provided context is cancelled, Next will return a non-nil
	// error. Any other failures encountered while waiting for the next
	// event (connection issues, deserialization failures, etc) will
	// also result in a non-nil error.
	Next(context.Context) (*Response, error)
***REMOVED***

type Response struct ***REMOVED***
	// Action is the name of the operation that occurred. Possible values
	// include get, set, delete, update, create, compareAndSwap,
	// compareAndDelete and expire.
	Action string `json:"action"`

	// Node represents the state of the relevant etcd Node.
	Node *Node `json:"node"`

	// PrevNode represents the previous state of the Node. PrevNode is non-nil
	// only if the Node existed before the action occurred and the action
	// caused a change to the Node.
	PrevNode *Node `json:"prevNode"`

	// Index holds the cluster-level index at the time the Response was generated.
	// This index is not tied to the Node(s) contained in this Response.
	Index uint64 `json:"-"`

	// ClusterID holds the cluster-level ID reported by the server.  This
	// should be different for different etcd clusters.
	ClusterID string `json:"-"`
***REMOVED***

type Node struct ***REMOVED***
	// Key represents the unique location of this Node (e.g. "/foo/bar").
	Key string `json:"key"`

	// Dir reports whether node describes a directory.
	Dir bool `json:"dir,omitempty"`

	// Value is the current data stored on this Node. If this Node
	// is a directory, Value will be empty.
	Value string `json:"value"`

	// Nodes holds the children of this Node, only if this Node is a directory.
	// This slice of will be arbitrarily deep (children, grandchildren, great-
	// grandchildren, etc.) if a recursive Get or Watch request were made.
	Nodes Nodes `json:"nodes"`

	// CreatedIndex is the etcd index at-which this Node was created.
	CreatedIndex uint64 `json:"createdIndex"`

	// ModifiedIndex is the etcd index at-which this Node was last modified.
	ModifiedIndex uint64 `json:"modifiedIndex"`

	// Expiration is the server side expiration time of the key.
	Expiration *time.Time `json:"expiration,omitempty"`

	// TTL is the time to live of the key in second.
	TTL int64 `json:"ttl,omitempty"`
***REMOVED***

func (n *Node) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***Key: %s, CreatedIndex: %d, ModifiedIndex: %d, TTL: %d***REMOVED***", n.Key, n.CreatedIndex, n.ModifiedIndex, n.TTL)
***REMOVED***

// TTLDuration returns the Node's TTL as a time.Duration object
func (n *Node) TTLDuration() time.Duration ***REMOVED***
	return time.Duration(n.TTL) * time.Second
***REMOVED***

type Nodes []*Node

// interfaces for sorting

func (ns Nodes) Len() int           ***REMOVED*** return len(ns) ***REMOVED***
func (ns Nodes) Less(i, j int) bool ***REMOVED*** return ns[i].Key < ns[j].Key ***REMOVED***
func (ns Nodes) Swap(i, j int)      ***REMOVED*** ns[i], ns[j] = ns[j], ns[i] ***REMOVED***

type httpKeysAPI struct ***REMOVED***
	client httpClient
	prefix string
***REMOVED***

func (k *httpKeysAPI) Set(ctx context.Context, key, val string, opts *SetOptions) (*Response, error) ***REMOVED***
	act := &setAction***REMOVED***
		Prefix: k.prefix,
		Key:    key,
		Value:  val,
	***REMOVED***

	if opts != nil ***REMOVED***
		act.PrevValue = opts.PrevValue
		act.PrevIndex = opts.PrevIndex
		act.PrevExist = opts.PrevExist
		act.TTL = opts.TTL
		act.Refresh = opts.Refresh
		act.Dir = opts.Dir
		act.NoValueOnSuccess = opts.NoValueOnSuccess
	***REMOVED***

	doCtx := ctx
	if act.PrevExist == PrevNoExist ***REMOVED***
		doCtx = context.WithValue(doCtx, &oneShotCtxValue, &oneShotCtxValue)
	***REMOVED***
	resp, body, err := k.client.Do(doCtx, act)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return unmarshalHTTPResponse(resp.StatusCode, resp.Header, body)
***REMOVED***

func (k *httpKeysAPI) Create(ctx context.Context, key, val string) (*Response, error) ***REMOVED***
	return k.Set(ctx, key, val, &SetOptions***REMOVED***PrevExist: PrevNoExist***REMOVED***)
***REMOVED***

func (k *httpKeysAPI) CreateInOrder(ctx context.Context, dir, val string, opts *CreateInOrderOptions) (*Response, error) ***REMOVED***
	act := &createInOrderAction***REMOVED***
		Prefix: k.prefix,
		Dir:    dir,
		Value:  val,
	***REMOVED***

	if opts != nil ***REMOVED***
		act.TTL = opts.TTL
	***REMOVED***

	resp, body, err := k.client.Do(ctx, act)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return unmarshalHTTPResponse(resp.StatusCode, resp.Header, body)
***REMOVED***

func (k *httpKeysAPI) Update(ctx context.Context, key, val string) (*Response, error) ***REMOVED***
	return k.Set(ctx, key, val, &SetOptions***REMOVED***PrevExist: PrevExist***REMOVED***)
***REMOVED***

func (k *httpKeysAPI) Delete(ctx context.Context, key string, opts *DeleteOptions) (*Response, error) ***REMOVED***
	act := &deleteAction***REMOVED***
		Prefix: k.prefix,
		Key:    key,
	***REMOVED***

	if opts != nil ***REMOVED***
		act.PrevValue = opts.PrevValue
		act.PrevIndex = opts.PrevIndex
		act.Dir = opts.Dir
		act.Recursive = opts.Recursive
	***REMOVED***

	doCtx := context.WithValue(ctx, &oneShotCtxValue, &oneShotCtxValue)
	resp, body, err := k.client.Do(doCtx, act)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return unmarshalHTTPResponse(resp.StatusCode, resp.Header, body)
***REMOVED***

func (k *httpKeysAPI) Get(ctx context.Context, key string, opts *GetOptions) (*Response, error) ***REMOVED***
	act := &getAction***REMOVED***
		Prefix: k.prefix,
		Key:    key,
	***REMOVED***

	if opts != nil ***REMOVED***
		act.Recursive = opts.Recursive
		act.Sorted = opts.Sort
		act.Quorum = opts.Quorum
	***REMOVED***

	resp, body, err := k.client.Do(ctx, act)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return unmarshalHTTPResponse(resp.StatusCode, resp.Header, body)
***REMOVED***

func (k *httpKeysAPI) Watcher(key string, opts *WatcherOptions) Watcher ***REMOVED***
	act := waitAction***REMOVED***
		Prefix: k.prefix,
		Key:    key,
	***REMOVED***

	if opts != nil ***REMOVED***
		act.Recursive = opts.Recursive
		if opts.AfterIndex > 0 ***REMOVED***
			act.WaitIndex = opts.AfterIndex + 1
		***REMOVED***
	***REMOVED***

	return &httpWatcher***REMOVED***
		client:   k.client,
		nextWait: act,
	***REMOVED***
***REMOVED***

type httpWatcher struct ***REMOVED***
	client   httpClient
	nextWait waitAction
***REMOVED***

func (hw *httpWatcher) Next(ctx context.Context) (*Response, error) ***REMOVED***
	for ***REMOVED***
		httpresp, body, err := hw.client.Do(ctx, &hw.nextWait)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		resp, err := unmarshalHTTPResponse(httpresp.StatusCode, httpresp.Header, body)
		if err != nil ***REMOVED***
			if err == ErrEmptyBody ***REMOVED***
				continue
			***REMOVED***
			return nil, err
		***REMOVED***

		hw.nextWait.WaitIndex = resp.Node.ModifiedIndex + 1
		return resp, nil
	***REMOVED***
***REMOVED***

// v2KeysURL forms a URL representing the location of a key.
// The endpoint argument represents the base URL of an etcd
// server. The prefix is the path needed to route from the
// provided endpoint's path to the root of the keys API
// (typically "/v2/keys").
func v2KeysURL(ep url.URL, prefix, key string) *url.URL ***REMOVED***
	// We concatenate all parts together manually. We cannot use
	// path.Join because it does not reserve trailing slash.
	// We call CanonicalURLPath to further cleanup the path.
	if prefix != "" && prefix[0] != '/' ***REMOVED***
		prefix = "/" + prefix
	***REMOVED***
	if key != "" && key[0] != '/' ***REMOVED***
		key = "/" + key
	***REMOVED***
	ep.Path = pathutil.CanonicalURLPath(ep.Path + prefix + key)
	return &ep
***REMOVED***

type getAction struct ***REMOVED***
	Prefix    string
	Key       string
	Recursive bool
	Sorted    bool
	Quorum    bool
***REMOVED***

func (g *getAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2KeysURL(ep, g.Prefix, g.Key)

	params := u.Query()
	params.Set("recursive", strconv.FormatBool(g.Recursive))
	params.Set("sorted", strconv.FormatBool(g.Sorted))
	params.Set("quorum", strconv.FormatBool(g.Quorum))
	u.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	return req
***REMOVED***

type waitAction struct ***REMOVED***
	Prefix    string
	Key       string
	WaitIndex uint64
	Recursive bool
***REMOVED***

func (w *waitAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2KeysURL(ep, w.Prefix, w.Key)

	params := u.Query()
	params.Set("wait", "true")
	params.Set("waitIndex", strconv.FormatUint(w.WaitIndex, 10))
	params.Set("recursive", strconv.FormatBool(w.Recursive))
	u.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	return req
***REMOVED***

type setAction struct ***REMOVED***
	Prefix           string
	Key              string
	Value            string
	PrevValue        string
	PrevIndex        uint64
	PrevExist        PrevExistType
	TTL              time.Duration
	Refresh          bool
	Dir              bool
	NoValueOnSuccess bool
***REMOVED***

func (a *setAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2KeysURL(ep, a.Prefix, a.Key)

	params := u.Query()
	form := url.Values***REMOVED******REMOVED***

	// we're either creating a directory or setting a key
	if a.Dir ***REMOVED***
		params.Set("dir", strconv.FormatBool(a.Dir))
	***REMOVED*** else ***REMOVED***
		// These options are only valid for setting a key
		if a.PrevValue != "" ***REMOVED***
			params.Set("prevValue", a.PrevValue)
		***REMOVED***
		form.Add("value", a.Value)
	***REMOVED***

	// Options which apply to both setting a key and creating a dir
	if a.PrevIndex != 0 ***REMOVED***
		params.Set("prevIndex", strconv.FormatUint(a.PrevIndex, 10))
	***REMOVED***
	if a.PrevExist != PrevIgnore ***REMOVED***
		params.Set("prevExist", string(a.PrevExist))
	***REMOVED***
	if a.TTL > 0 ***REMOVED***
		form.Add("ttl", strconv.FormatUint(uint64(a.TTL.Seconds()), 10))
	***REMOVED***

	if a.Refresh ***REMOVED***
		form.Add("refresh", "true")
	***REMOVED***
	if a.NoValueOnSuccess ***REMOVED***
		params.Set("noValueOnSuccess", strconv.FormatBool(a.NoValueOnSuccess))
	***REMOVED***

	u.RawQuery = params.Encode()
	body := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("PUT", u.String(), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
***REMOVED***

type deleteAction struct ***REMOVED***
	Prefix    string
	Key       string
	PrevValue string
	PrevIndex uint64
	Dir       bool
	Recursive bool
***REMOVED***

func (a *deleteAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2KeysURL(ep, a.Prefix, a.Key)

	params := u.Query()
	if a.PrevValue != "" ***REMOVED***
		params.Set("prevValue", a.PrevValue)
	***REMOVED***
	if a.PrevIndex != 0 ***REMOVED***
		params.Set("prevIndex", strconv.FormatUint(a.PrevIndex, 10))
	***REMOVED***
	if a.Dir ***REMOVED***
		params.Set("dir", "true")
	***REMOVED***
	if a.Recursive ***REMOVED***
		params.Set("recursive", "true")
	***REMOVED***
	u.RawQuery = params.Encode()

	req, _ := http.NewRequest("DELETE", u.String(), nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
***REMOVED***

type createInOrderAction struct ***REMOVED***
	Prefix string
	Dir    string
	Value  string
	TTL    time.Duration
***REMOVED***

func (a *createInOrderAction) HTTPRequest(ep url.URL) *http.Request ***REMOVED***
	u := v2KeysURL(ep, a.Prefix, a.Dir)

	form := url.Values***REMOVED******REMOVED***
	form.Add("value", a.Value)
	if a.TTL > 0 ***REMOVED***
		form.Add("ttl", strconv.FormatUint(uint64(a.TTL.Seconds()), 10))
	***REMOVED***
	body := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("POST", u.String(), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
***REMOVED***

func unmarshalHTTPResponse(code int, header http.Header, body []byte) (res *Response, err error) ***REMOVED***
	switch code ***REMOVED***
	case http.StatusOK, http.StatusCreated:
		if len(body) == 0 ***REMOVED***
			return nil, ErrEmptyBody
		***REMOVED***
		res, err = unmarshalSuccessfulKeysResponse(header, body)
	default:
		err = unmarshalFailedKeysResponse(body)
	***REMOVED***

	return
***REMOVED***

func unmarshalSuccessfulKeysResponse(header http.Header, body []byte) (*Response, error) ***REMOVED***
	var res Response
	err := codec.NewDecoderBytes(body, new(codec.JsonHandle)).Decode(&res)
	if err != nil ***REMOVED***
		return nil, ErrInvalidJSON
	***REMOVED***
	if header.Get("X-Etcd-Index") != "" ***REMOVED***
		res.Index, err = strconv.ParseUint(header.Get("X-Etcd-Index"), 10, 64)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	res.ClusterID = header.Get("X-Etcd-Cluster-ID")
	return &res, nil
***REMOVED***

func unmarshalFailedKeysResponse(body []byte) error ***REMOVED***
	var etcdErr Error
	if err := json.Unmarshal(body, &etcdErr); err != nil ***REMOVED***
		return ErrInvalidJSON
	***REMOVED***
	return etcdErr
***REMOVED***
