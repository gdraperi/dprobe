package serf

import (
	"encoding/base64"
	"fmt"
	"sync"
)

// KeyManager encapsulates all functionality within Serf for handling
// encryption keyring changes across a cluster.
type KeyManager struct ***REMOVED***
	serf *Serf

	// Lock to protect read and write operations
	l sync.RWMutex
***REMOVED***

// keyRequest is used to contain input parameters which get broadcasted to all
// nodes as part of a key query operation.
type keyRequest struct ***REMOVED***
	Key []byte
***REMOVED***

// KeyResponse is used to relay a query for a list of all keys in use.
type KeyResponse struct ***REMOVED***
	Messages map[string]string // Map of node name to response message
	NumNodes int               // Total nodes memberlist knows of
	NumResp  int               // Total responses received
	NumErr   int               // Total errors from request

	// Keys is a mapping of the base64-encoded value of the key bytes to the
	// number of nodes that have the key installed.
	Keys map[string]int
***REMOVED***

// streamKeyResp takes care of reading responses from a channel and composing
// them into a KeyResponse. It will update a KeyResponse *in place* and
// therefore has nothing to return.
func (k *KeyManager) streamKeyResp(resp *KeyResponse, ch <-chan NodeResponse) ***REMOVED***
	for r := range ch ***REMOVED***
		var nodeResponse nodeKeyResponse

		resp.NumResp++

		// Decode the response
		if len(r.Payload) < 1 || messageType(r.Payload[0]) != messageKeyResponseType ***REMOVED***
			resp.Messages[r.From] = fmt.Sprintf(
				"Invalid key query response type: %v", r.Payload)
			resp.NumErr++
			goto NEXT
		***REMOVED***
		if err := decodeMessage(r.Payload[1:], &nodeResponse); err != nil ***REMOVED***
			resp.Messages[r.From] = fmt.Sprintf(
				"Failed to decode key query response: %v", r.Payload)
			resp.NumErr++
			goto NEXT
		***REMOVED***

		if !nodeResponse.Result ***REMOVED***
			resp.Messages[r.From] = nodeResponse.Message
			resp.NumErr++
		***REMOVED***

		// Currently only used for key list queries, this adds keys to a counter
		// and increments them for each node response which contains them.
		for _, key := range nodeResponse.Keys ***REMOVED***
			if _, ok := resp.Keys[key]; !ok ***REMOVED***
				resp.Keys[key] = 1
			***REMOVED*** else ***REMOVED***
				resp.Keys[key]++
			***REMOVED***
		***REMOVED***

	NEXT:
		// Return early if all nodes have responded. This allows us to avoid
		// waiting for the full timeout when there is nothing left to do.
		if resp.NumResp == resp.NumNodes ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// handleKeyRequest performs query broadcasting to all members for any type of
// key operation and manages gathering responses and packing them up into a
// KeyResponse for uniform response handling.
func (k *KeyManager) handleKeyRequest(key, query string) (*KeyResponse, error) ***REMOVED***
	resp := &KeyResponse***REMOVED***
		Messages: make(map[string]string),
		Keys:     make(map[string]int),
	***REMOVED***
	qName := internalQueryName(query)

	// Decode the new key into raw bytes
	rawKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil ***REMOVED***
		return resp, err
	***REMOVED***

	// Encode the query request
	req, err := encodeMessage(messageKeyRequestType, keyRequest***REMOVED***Key: rawKey***REMOVED***)
	if err != nil ***REMOVED***
		return resp, err
	***REMOVED***

	qParam := k.serf.DefaultQueryParams()
	queryResp, err := k.serf.Query(qName, req, qParam)
	if err != nil ***REMOVED***
		return resp, err
	***REMOVED***

	// Handle the response stream and populate the KeyResponse
	resp.NumNodes = k.serf.memberlist.NumMembers()
	k.streamKeyResp(resp, queryResp.respCh)

	// Check the response for any reported failure conditions
	if resp.NumErr != 0 ***REMOVED***
		return resp, fmt.Errorf("%d/%d nodes reported failure", resp.NumErr, resp.NumNodes)
	***REMOVED***
	if resp.NumResp != resp.NumNodes ***REMOVED***
		return resp, fmt.Errorf("%d/%d nodes reported success", resp.NumResp, resp.NumNodes)
	***REMOVED***

	return resp, nil
***REMOVED***

// InstallKey handles broadcasting a query to all members and gathering
// responses from each of them, returning a list of messages from each node
// and any applicable error conditions.
func (k *KeyManager) InstallKey(key string) (*KeyResponse, error) ***REMOVED***
	k.l.Lock()
	defer k.l.Unlock()

	return k.handleKeyRequest(key, installKeyQuery)
***REMOVED***

// UseKey handles broadcasting a primary key change to all members in the
// cluster, and gathering any response messages. If successful, there should
// be an empty KeyResponse returned.
func (k *KeyManager) UseKey(key string) (*KeyResponse, error) ***REMOVED***
	k.l.Lock()
	defer k.l.Unlock()

	return k.handleKeyRequest(key, useKeyQuery)
***REMOVED***

// RemoveKey handles broadcasting a key to the cluster for removal. Each member
// will receive this event, and if they have the key in their keyring, remove
// it. If any errors are encountered, RemoveKey will collect and relay them.
func (k *KeyManager) RemoveKey(key string) (*KeyResponse, error) ***REMOVED***
	k.l.Lock()
	defer k.l.Unlock()

	return k.handleKeyRequest(key, removeKeyQuery)
***REMOVED***

// ListKeys is used to collect installed keys from members in a Serf cluster
// and return an aggregated list of all installed keys. This is useful to
// operators to ensure that there are no lingering keys installed on any agents.
// Since having multiple keys installed can cause performance penalties in some
// cases, it's important to verify this information and remove unneeded keys.
func (k *KeyManager) ListKeys() (*KeyResponse, error) ***REMOVED***
	k.l.RLock()
	defer k.l.RUnlock()

	return k.handleKeyRequest("", listKeysQuery)
***REMOVED***
