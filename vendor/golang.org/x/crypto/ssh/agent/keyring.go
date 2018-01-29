// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type privKey struct ***REMOVED***
	signer  ssh.Signer
	comment string
	expire  *time.Time
***REMOVED***

type keyring struct ***REMOVED***
	mu   sync.Mutex
	keys []privKey

	locked     bool
	passphrase []byte
***REMOVED***

var errLocked = errors.New("agent: locked")

// NewKeyring returns an Agent that holds keys in memory.  It is safe
// for concurrent use by multiple goroutines.
func NewKeyring() Agent ***REMOVED***
	return &keyring***REMOVED******REMOVED***
***REMOVED***

// RemoveAll removes all identities.
func (r *keyring) RemoveAll() error ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.locked ***REMOVED***
		return errLocked
	***REMOVED***

	r.keys = nil
	return nil
***REMOVED***

// removeLocked does the actual key removal. The caller must already be holding the
// keyring mutex.
func (r *keyring) removeLocked(want []byte) error ***REMOVED***
	found := false
	for i := 0; i < len(r.keys); ***REMOVED***
		if bytes.Equal(r.keys[i].signer.PublicKey().Marshal(), want) ***REMOVED***
			found = true
			r.keys[i] = r.keys[len(r.keys)-1]
			r.keys = r.keys[:len(r.keys)-1]
			continue
		***REMOVED*** else ***REMOVED***
			i++
		***REMOVED***
	***REMOVED***

	if !found ***REMOVED***
		return errors.New("agent: key not found")
	***REMOVED***
	return nil
***REMOVED***

// Remove removes all identities with the given public key.
func (r *keyring) Remove(key ssh.PublicKey) error ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.locked ***REMOVED***
		return errLocked
	***REMOVED***

	return r.removeLocked(key.Marshal())
***REMOVED***

// Lock locks the agent. Sign and Remove will fail, and List will return an empty list.
func (r *keyring) Lock(passphrase []byte) error ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.locked ***REMOVED***
		return errLocked
	***REMOVED***

	r.locked = true
	r.passphrase = passphrase
	return nil
***REMOVED***

// Unlock undoes the effect of Lock
func (r *keyring) Unlock(passphrase []byte) error ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.locked ***REMOVED***
		return errors.New("agent: not locked")
	***REMOVED***
	if len(passphrase) != len(r.passphrase) || 1 != subtle.ConstantTimeCompare(passphrase, r.passphrase) ***REMOVED***
		return fmt.Errorf("agent: incorrect passphrase")
	***REMOVED***

	r.locked = false
	r.passphrase = nil
	return nil
***REMOVED***

// expireKeysLocked removes expired keys from the keyring. If a key was added
// with a lifetimesecs contraint and seconds >= lifetimesecs seconds have
// ellapsed, it is removed. The caller *must* be holding the keyring mutex.
func (r *keyring) expireKeysLocked() ***REMOVED***
	for _, k := range r.keys ***REMOVED***
		if k.expire != nil && time.Now().After(*k.expire) ***REMOVED***
			r.removeLocked(k.signer.PublicKey().Marshal())
		***REMOVED***
	***REMOVED***
***REMOVED***

// List returns the identities known to the agent.
func (r *keyring) List() ([]*Key, error) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.locked ***REMOVED***
		// section 2.7: locked agents return empty.
		return nil, nil
	***REMOVED***

	r.expireKeysLocked()
	var ids []*Key
	for _, k := range r.keys ***REMOVED***
		pub := k.signer.PublicKey()
		ids = append(ids, &Key***REMOVED***
			Format:  pub.Type(),
			Blob:    pub.Marshal(),
			Comment: k.comment***REMOVED***)
	***REMOVED***
	return ids, nil
***REMOVED***

// Insert adds a private key to the keyring. If a certificate
// is given, that certificate is added as public key. Note that
// any constraints given are ignored.
func (r *keyring) Add(key AddedKey) error ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.locked ***REMOVED***
		return errLocked
	***REMOVED***
	signer, err := ssh.NewSignerFromKey(key.PrivateKey)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if cert := key.Certificate; cert != nil ***REMOVED***
		signer, err = ssh.NewCertSigner(cert, signer)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	p := privKey***REMOVED***
		signer:  signer,
		comment: key.Comment,
	***REMOVED***

	if key.LifetimeSecs > 0 ***REMOVED***
		t := time.Now().Add(time.Duration(key.LifetimeSecs) * time.Second)
		p.expire = &t
	***REMOVED***

	r.keys = append(r.keys, p)

	return nil
***REMOVED***

// Sign returns a signature for the data.
func (r *keyring) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.locked ***REMOVED***
		return nil, errLocked
	***REMOVED***

	r.expireKeysLocked()
	wanted := key.Marshal()
	for _, k := range r.keys ***REMOVED***
		if bytes.Equal(k.signer.PublicKey().Marshal(), wanted) ***REMOVED***
			return k.signer.Sign(rand.Reader, data)
		***REMOVED***
	***REMOVED***
	return nil, errors.New("not found")
***REMOVED***

// Signers returns signers for all the known keys.
func (r *keyring) Signers() ([]ssh.Signer, error) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.locked ***REMOVED***
		return nil, errLocked
	***REMOVED***

	r.expireKeysLocked()
	s := make([]ssh.Signer, 0, len(r.keys))
	for _, k := range r.keys ***REMOVED***
		s = append(s, k.signer)
	***REMOVED***
	return s, nil
***REMOVED***
