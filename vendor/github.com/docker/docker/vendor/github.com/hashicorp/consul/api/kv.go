package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// KVPair is used to represent a single K/V entry
type KVPair struct ***REMOVED***
	Key         string
	CreateIndex uint64
	ModifyIndex uint64
	LockIndex   uint64
	Flags       uint64
	Value       []byte
	Session     string
***REMOVED***

// KVPairs is a list of KVPair objects
type KVPairs []*KVPair

// KV is used to manipulate the K/V API
type KV struct ***REMOVED***
	c *Client
***REMOVED***

// KV is used to return a handle to the K/V apis
func (c *Client) KV() *KV ***REMOVED***
	return &KV***REMOVED***c***REMOVED***
***REMOVED***

// Get is used to lookup a single key
func (k *KV) Get(key string, q *QueryOptions) (*KVPair, *QueryMeta, error) ***REMOVED***
	resp, qm, err := k.getInternal(key, nil, q)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if resp == nil ***REMOVED***
		return nil, qm, nil
	***REMOVED***
	defer resp.Body.Close()

	var entries []*KVPair
	if err := decodeBody(resp, &entries); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if len(entries) > 0 ***REMOVED***
		return entries[0], qm, nil
	***REMOVED***
	return nil, qm, nil
***REMOVED***

// List is used to lookup all keys under a prefix
func (k *KV) List(prefix string, q *QueryOptions) (KVPairs, *QueryMeta, error) ***REMOVED***
	resp, qm, err := k.getInternal(prefix, map[string]string***REMOVED***"recurse": ""***REMOVED***, q)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if resp == nil ***REMOVED***
		return nil, qm, nil
	***REMOVED***
	defer resp.Body.Close()

	var entries []*KVPair
	if err := decodeBody(resp, &entries); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return entries, qm, nil
***REMOVED***

// Keys is used to list all the keys under a prefix. Optionally,
// a separator can be used to limit the responses.
func (k *KV) Keys(prefix, separator string, q *QueryOptions) ([]string, *QueryMeta, error) ***REMOVED***
	params := map[string]string***REMOVED***"keys": ""***REMOVED***
	if separator != "" ***REMOVED***
		params["separator"] = separator
	***REMOVED***
	resp, qm, err := k.getInternal(prefix, params, q)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if resp == nil ***REMOVED***
		return nil, qm, nil
	***REMOVED***
	defer resp.Body.Close()

	var entries []string
	if err := decodeBody(resp, &entries); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return entries, qm, nil
***REMOVED***

func (k *KV) getInternal(key string, params map[string]string, q *QueryOptions) (*http.Response, *QueryMeta, error) ***REMOVED***
	r := k.c.newRequest("GET", "/v1/kv/"+key)
	r.setQueryOptions(q)
	for param, val := range params ***REMOVED***
		r.params.Set(param, val)
	***REMOVED***
	rtt, resp, err := k.c.doRequest(r)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	qm := &QueryMeta***REMOVED******REMOVED***
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	if resp.StatusCode == 404 ***REMOVED***
		resp.Body.Close()
		return nil, qm, nil
	***REMOVED*** else if resp.StatusCode != 200 ***REMOVED***
		resp.Body.Close()
		return nil, nil, fmt.Errorf("Unexpected response code: %d", resp.StatusCode)
	***REMOVED***
	return resp, qm, nil
***REMOVED***

// Put is used to write a new value. Only the
// Key, Flags and Value is respected.
func (k *KV) Put(p *KVPair, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	params := make(map[string]string, 1)
	if p.Flags != 0 ***REMOVED***
		params["flags"] = strconv.FormatUint(p.Flags, 10)
	***REMOVED***
	_, wm, err := k.put(p.Key, params, p.Value, q)
	return wm, err
***REMOVED***

// CAS is used for a Check-And-Set operation. The Key,
// ModifyIndex, Flags and Value are respected. Returns true
// on success or false on failures.
func (k *KV) CAS(p *KVPair, q *WriteOptions) (bool, *WriteMeta, error) ***REMOVED***
	params := make(map[string]string, 2)
	if p.Flags != 0 ***REMOVED***
		params["flags"] = strconv.FormatUint(p.Flags, 10)
	***REMOVED***
	params["cas"] = strconv.FormatUint(p.ModifyIndex, 10)
	return k.put(p.Key, params, p.Value, q)
***REMOVED***

// Acquire is used for a lock acquisiiton operation. The Key,
// Flags, Value and Session are respected. Returns true
// on success or false on failures.
func (k *KV) Acquire(p *KVPair, q *WriteOptions) (bool, *WriteMeta, error) ***REMOVED***
	params := make(map[string]string, 2)
	if p.Flags != 0 ***REMOVED***
		params["flags"] = strconv.FormatUint(p.Flags, 10)
	***REMOVED***
	params["acquire"] = p.Session
	return k.put(p.Key, params, p.Value, q)
***REMOVED***

// Release is used for a lock release operation. The Key,
// Flags, Value and Session are respected. Returns true
// on success or false on failures.
func (k *KV) Release(p *KVPair, q *WriteOptions) (bool, *WriteMeta, error) ***REMOVED***
	params := make(map[string]string, 2)
	if p.Flags != 0 ***REMOVED***
		params["flags"] = strconv.FormatUint(p.Flags, 10)
	***REMOVED***
	params["release"] = p.Session
	return k.put(p.Key, params, p.Value, q)
***REMOVED***

func (k *KV) put(key string, params map[string]string, body []byte, q *WriteOptions) (bool, *WriteMeta, error) ***REMOVED***
	r := k.c.newRequest("PUT", "/v1/kv/"+key)
	r.setWriteOptions(q)
	for param, val := range params ***REMOVED***
		r.params.Set(param, val)
	***REMOVED***
	r.body = bytes.NewReader(body)
	rtt, resp, err := requireOK(k.c.doRequest(r))
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &WriteMeta***REMOVED******REMOVED***
	qm.RequestTime = rtt

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil ***REMOVED***
		return false, nil, fmt.Errorf("Failed to read response: %v", err)
	***REMOVED***
	res := strings.Contains(string(buf.Bytes()), "true")
	return res, qm, nil
***REMOVED***

// Delete is used to delete a single key
func (k *KV) Delete(key string, w *WriteOptions) (*WriteMeta, error) ***REMOVED***
	_, qm, err := k.deleteInternal(key, nil, w)
	return qm, err
***REMOVED***

// DeleteCAS is used for a Delete Check-And-Set operation. The Key
// and ModifyIndex are respected. Returns true on success or false on failures.
func (k *KV) DeleteCAS(p *KVPair, q *WriteOptions) (bool, *WriteMeta, error) ***REMOVED***
	params := map[string]string***REMOVED***
		"cas": strconv.FormatUint(p.ModifyIndex, 10),
	***REMOVED***
	return k.deleteInternal(p.Key, params, q)
***REMOVED***

// DeleteTree is used to delete all keys under a prefix
func (k *KV) DeleteTree(prefix string, w *WriteOptions) (*WriteMeta, error) ***REMOVED***
	_, qm, err := k.deleteInternal(prefix, map[string]string***REMOVED***"recurse": ""***REMOVED***, w)
	return qm, err
***REMOVED***

func (k *KV) deleteInternal(key string, params map[string]string, q *WriteOptions) (bool, *WriteMeta, error) ***REMOVED***
	r := k.c.newRequest("DELETE", "/v1/kv/"+key)
	r.setWriteOptions(q)
	for param, val := range params ***REMOVED***
		r.params.Set(param, val)
	***REMOVED***
	rtt, resp, err := requireOK(k.c.doRequest(r))
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***
	defer resp.Body.Close()

	qm := &WriteMeta***REMOVED******REMOVED***
	qm.RequestTime = rtt

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil ***REMOVED***
		return false, nil, fmt.Errorf("Failed to read response: %v", err)
	***REMOVED***
	res := strings.Contains(string(buf.Bytes()), "true")
	return res, qm, nil
***REMOVED***
