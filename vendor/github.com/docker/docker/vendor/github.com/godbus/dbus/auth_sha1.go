package dbus

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"os"
)

// AuthCookieSha1 returns an Auth that authenticates as the given user with the
// DBUS_COOKIE_SHA1 mechanism. The home parameter should specify the home
// directory of the user.
func AuthCookieSha1(user, home string) Auth ***REMOVED***
	return authCookieSha1***REMOVED***user, home***REMOVED***
***REMOVED***

type authCookieSha1 struct ***REMOVED***
	user, home string
***REMOVED***

func (a authCookieSha1) FirstData() ([]byte, []byte, AuthStatus) ***REMOVED***
	b := make([]byte, 2*len(a.user))
	hex.Encode(b, []byte(a.user))
	return []byte("DBUS_COOKIE_SHA1"), b, AuthContinue
***REMOVED***

func (a authCookieSha1) HandleData(data []byte) ([]byte, AuthStatus) ***REMOVED***
	challenge := make([]byte, len(data)/2)
	_, err := hex.Decode(challenge, data)
	if err != nil ***REMOVED***
		return nil, AuthError
	***REMOVED***
	b := bytes.Split(challenge, []byte***REMOVED***' '***REMOVED***)
	if len(b) != 3 ***REMOVED***
		return nil, AuthError
	***REMOVED***
	context := b[0]
	id := b[1]
	svchallenge := b[2]
	cookie := a.getCookie(context, id)
	if cookie == nil ***REMOVED***
		return nil, AuthError
	***REMOVED***
	clchallenge := a.generateChallenge()
	if clchallenge == nil ***REMOVED***
		return nil, AuthError
	***REMOVED***
	hash := sha1.New()
	hash.Write(bytes.Join([][]byte***REMOVED***svchallenge, clchallenge, cookie***REMOVED***, []byte***REMOVED***':'***REMOVED***))
	hexhash := make([]byte, 2*hash.Size())
	hex.Encode(hexhash, hash.Sum(nil))
	data = append(clchallenge, ' ')
	data = append(data, hexhash...)
	resp := make([]byte, 2*len(data))
	hex.Encode(resp, data)
	return resp, AuthOk
***REMOVED***

// getCookie searches for the cookie identified by id in context and returns
// the cookie content or nil. (Since HandleData can't return a specific error,
// but only whether an error occured, this function also doesn't bother to
// return an error.)
func (a authCookieSha1) getCookie(context, id []byte) []byte ***REMOVED***
	file, err := os.Open(a.home + "/.dbus-keyrings/" + string(context))
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	defer file.Close()
	rd := bufio.NewReader(file)
	for ***REMOVED***
		line, err := rd.ReadBytes('\n')
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		line = line[:len(line)-1]
		b := bytes.Split(line, []byte***REMOVED***' '***REMOVED***)
		if len(b) != 3 ***REMOVED***
			return nil
		***REMOVED***
		if bytes.Equal(b[0], id) ***REMOVED***
			return b[2]
		***REMOVED***
	***REMOVED***
***REMOVED***

// generateChallenge returns a random, hex-encoded challenge, or nil on error
// (see above).
func (a authCookieSha1) generateChallenge() []byte ***REMOVED***
	b := make([]byte, 16)
	n, err := rand.Read(b)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	if n != 16 ***REMOVED***
		return nil
	***REMOVED***
	enc := make([]byte, 32)
	hex.Encode(enc, b)
	return enc
***REMOVED***
