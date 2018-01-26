package manager

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/manager/encryption"
	"github.com/docker/swarmkit/manager/state/raft"
)

const (
	// the raft DEK (data encryption key) is stored in the TLS key as a header
	// these are the header values
	pemHeaderRaftDEK              = "raft-dek"
	pemHeaderRaftPendingDEK       = "raft-dek-pending"
	pemHeaderRaftDEKNeedsRotation = "raft-dek-needs-rotation"
)

// RaftDEKData contains all the data stored in TLS pem headers
type RaftDEKData struct ***REMOVED***
	raft.EncryptionKeys
	NeedsRotation bool
***REMOVED***

// UnmarshalHeaders loads the state of the DEK manager given the current TLS headers
func (r RaftDEKData) UnmarshalHeaders(headers map[string]string, kekData ca.KEKData) (ca.PEMKeyHeaders, error) ***REMOVED***
	var (
		currentDEK, pendingDEK []byte
		err                    error
	)

	if currentDEKStr, ok := headers[pemHeaderRaftDEK]; ok ***REMOVED***
		currentDEK, err = decodePEMHeaderValue(currentDEKStr, kekData.KEK)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if pendingDEKStr, ok := headers[pemHeaderRaftPendingDEK]; ok ***REMOVED***
		pendingDEK, err = decodePEMHeaderValue(pendingDEKStr, kekData.KEK)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if pendingDEK != nil && currentDEK == nil ***REMOVED***
		return nil, fmt.Errorf("there is a pending DEK, but no current DEK")
	***REMOVED***

	_, ok := headers[pemHeaderRaftDEKNeedsRotation]
	return RaftDEKData***REMOVED***
		NeedsRotation: ok,
		EncryptionKeys: raft.EncryptionKeys***REMOVED***
			CurrentDEK: currentDEK,
			PendingDEK: pendingDEK,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

// MarshalHeaders returns new headers given the current KEK
func (r RaftDEKData) MarshalHeaders(kekData ca.KEKData) (map[string]string, error) ***REMOVED***
	headers := make(map[string]string)
	for headerKey, contents := range map[string][]byte***REMOVED***
		pemHeaderRaftDEK:        r.CurrentDEK,
		pemHeaderRaftPendingDEK: r.PendingDEK,
	***REMOVED*** ***REMOVED***
		if contents != nil ***REMOVED***
			dekStr, err := encodePEMHeaderValue(contents, kekData.KEK)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			headers[headerKey] = dekStr
		***REMOVED***
	***REMOVED***

	if r.NeedsRotation ***REMOVED***
		headers[pemHeaderRaftDEKNeedsRotation] = "true"
	***REMOVED***

	// return a function that updates the dek data on write success
	return headers, nil
***REMOVED***

// UpdateKEK optionally sets NeedRotation to true if we go from unlocked to locked
func (r RaftDEKData) UpdateKEK(oldKEK, candidateKEK ca.KEKData) ca.PEMKeyHeaders ***REMOVED***
	if _, unlockedToLocked, err := compareKEKs(oldKEK, candidateKEK); err == nil && unlockedToLocked ***REMOVED***
		return RaftDEKData***REMOVED***
			EncryptionKeys: r.EncryptionKeys,
			NeedsRotation:  true,
		***REMOVED***
	***REMOVED***
	return r
***REMOVED***

// Returns whether the old KEK should be replaced with the new KEK, whether we went from
// unlocked to locked, and whether there was an error (the versions are the same, but the
// keks are different)
func compareKEKs(oldKEK, candidateKEK ca.KEKData) (bool, bool, error) ***REMOVED***
	keksEqual := subtle.ConstantTimeCompare(oldKEK.KEK, candidateKEK.KEK) == 1
	switch ***REMOVED***
	case oldKEK.Version == candidateKEK.Version && !keksEqual:
		return false, false, fmt.Errorf("candidate KEK has the same version as the current KEK, but a different KEK value")
	case oldKEK.Version >= candidateKEK.Version || keksEqual:
		return false, false, nil
	default:
		return true, oldKEK.KEK == nil, nil
	***REMOVED***
***REMOVED***

// RaftDEKManager manages the raft DEK keys using TLS headers
type RaftDEKManager struct ***REMOVED***
	kw         ca.KeyWriter
	rotationCh chan struct***REMOVED******REMOVED***
***REMOVED***

var errNoUpdateNeeded = fmt.Errorf("don't need to rotate or update")

// this error is returned if the KeyReadWriter's PEMKeyHeaders object is no longer a RaftDEKData object -
// this can happen if the node is no longer a manager, for example
var errNotUsingRaftDEKData = fmt.Errorf("RaftDEKManager can no longer store and manage TLS key headers")

// NewRaftDEKManager returns a RaftDEKManager that uses the current key writer
// and header manager
func NewRaftDEKManager(kw ca.KeyWriter) (*RaftDEKManager, error) ***REMOVED***
	// If there is no current DEK, generate one and write it to disk
	err := kw.ViewAndUpdateHeaders(func(h ca.PEMKeyHeaders) (ca.PEMKeyHeaders, error) ***REMOVED***
		dekData, ok := h.(RaftDEKData)
		// it wasn't a raft DEK manager before - just replace it
		if !ok || dekData.CurrentDEK == nil ***REMOVED***
			return RaftDEKData***REMOVED***
				EncryptionKeys: raft.EncryptionKeys***REMOVED***
					CurrentDEK: encryption.GenerateSecretKey(),
				***REMOVED***,
			***REMOVED***, nil
		***REMOVED***
		return nil, errNoUpdateNeeded
	***REMOVED***)
	if err != nil && err != errNoUpdateNeeded ***REMOVED***
		return nil, err
	***REMOVED***
	return &RaftDEKManager***REMOVED***
		kw:         kw,
		rotationCh: make(chan struct***REMOVED******REMOVED***, 1),
	***REMOVED***, nil
***REMOVED***

// NeedsRotation returns a boolean about whether we should do a rotation
func (r *RaftDEKManager) NeedsRotation() bool ***REMOVED***
	h, _ := r.kw.GetCurrentState()
	data, ok := h.(RaftDEKData)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return data.NeedsRotation || data.EncryptionKeys.PendingDEK != nil
***REMOVED***

// GetKeys returns the current set of DEKs.  If NeedsRotation is true, and there
// is no existing PendingDEK, it will try to create one.  If there are any errors
// doing so, just return the original.
func (r *RaftDEKManager) GetKeys() raft.EncryptionKeys ***REMOVED***
	var newKeys, originalKeys raft.EncryptionKeys
	err := r.kw.ViewAndUpdateHeaders(func(h ca.PEMKeyHeaders) (ca.PEMKeyHeaders, error) ***REMOVED***
		data, ok := h.(RaftDEKData)
		if !ok ***REMOVED***
			return nil, errNotUsingRaftDEKData
		***REMOVED***
		originalKeys = data.EncryptionKeys
		if !data.NeedsRotation || data.PendingDEK != nil ***REMOVED***
			return nil, errNoUpdateNeeded
		***REMOVED***
		newKeys = raft.EncryptionKeys***REMOVED***
			CurrentDEK: data.CurrentDEK,
			PendingDEK: encryption.GenerateSecretKey(),
		***REMOVED***
		return RaftDEKData***REMOVED***EncryptionKeys: newKeys***REMOVED***, nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return originalKeys
	***REMOVED***
	return newKeys
***REMOVED***

// RotationNotify the channel used to notify subscribers as to whether there
// should be a rotation done
func (r *RaftDEKManager) RotationNotify() chan struct***REMOVED******REMOVED*** ***REMOVED***
	return r.rotationCh
***REMOVED***

// UpdateKeys will set the updated encryption keys in the headers.  This finishes
// a rotation, and is expected to set the CurrentDEK to the previous PendingDEK.
func (r *RaftDEKManager) UpdateKeys(newKeys raft.EncryptionKeys) error ***REMOVED***
	return r.kw.ViewAndUpdateHeaders(func(h ca.PEMKeyHeaders) (ca.PEMKeyHeaders, error) ***REMOVED***
		data, ok := h.(RaftDEKData)
		if !ok ***REMOVED***
			return nil, errNotUsingRaftDEKData
		***REMOVED***
		// If there is no current DEK, we are basically wiping out all DEKs (no header object)
		if newKeys.CurrentDEK == nil ***REMOVED***
			return nil, nil
		***REMOVED***
		return RaftDEKData***REMOVED***
			EncryptionKeys: newKeys,
			NeedsRotation:  data.NeedsRotation,
		***REMOVED***, nil
	***REMOVED***)
***REMOVED***

// MaybeUpdateKEK does a KEK rotation if one is required.  Returns whether
// the kek was updated, whether it went from unlocked to locked, and any errors.
func (r *RaftDEKManager) MaybeUpdateKEK(candidateKEK ca.KEKData) (bool, bool, error) ***REMOVED***
	var updated, unlockedToLocked bool
	err := r.kw.ViewAndRotateKEK(func(currentKEK ca.KEKData, h ca.PEMKeyHeaders) (ca.KEKData, ca.PEMKeyHeaders, error) ***REMOVED***
		var err error
		updated, unlockedToLocked, err = compareKEKs(currentKEK, candidateKEK)
		if err == nil && !updated ***REMOVED*** // if we don't need to rotate the KEK, don't bother updating
			err = errNoUpdateNeeded
		***REMOVED***
		if err != nil ***REMOVED***
			return ca.KEKData***REMOVED******REMOVED***, nil, err
		***REMOVED***

		data, ok := h.(RaftDEKData)
		if !ok ***REMOVED***
			return ca.KEKData***REMOVED******REMOVED***, nil, errNotUsingRaftDEKData
		***REMOVED***

		if unlockedToLocked ***REMOVED***
			data.NeedsRotation = true
		***REMOVED***
		return candidateKEK, data, nil
	***REMOVED***)
	if err == errNoUpdateNeeded ***REMOVED***
		err = nil
	***REMOVED***

	if err == nil && unlockedToLocked ***REMOVED***
		r.rotationCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return updated, unlockedToLocked, err
***REMOVED***

func decodePEMHeaderValue(headerValue string, kek []byte) ([]byte, error) ***REMOVED***
	var decrypter encryption.Decrypter = encryption.NoopCrypter
	if kek != nil ***REMOVED***
		_, decrypter = encryption.Defaults(kek)
	***REMOVED***
	valueBytes, err := base64.StdEncoding.DecodeString(headerValue)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	result, err := encryption.Decrypt(valueBytes, decrypter)
	if err != nil ***REMOVED***
		return nil, ca.ErrInvalidKEK***REMOVED***Wrapped: err***REMOVED***
	***REMOVED***
	return result, nil
***REMOVED***

func encodePEMHeaderValue(headerValue []byte, kek []byte) (string, error) ***REMOVED***
	var encrypter encryption.Encrypter = encryption.NoopCrypter
	if kek != nil ***REMOVED***
		encrypter, _ = encryption.Defaults(kek)
	***REMOVED***
	encrypted, err := encryption.Encrypt(headerValue, encrypter)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return base64.StdEncoding.EncodeToString(encrypted), nil
***REMOVED***
