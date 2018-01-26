package dbus

import (
	"encoding/hex"
)

// AuthExternal returns an Auth that authenticates as the given user with the
// EXTERNAL mechanism.
func AuthExternal(user string) Auth ***REMOVED***
	return authExternal***REMOVED***user***REMOVED***
***REMOVED***

// AuthExternal implements the EXTERNAL authentication mechanism.
type authExternal struct ***REMOVED***
	user string
***REMOVED***

func (a authExternal) FirstData() ([]byte, []byte, AuthStatus) ***REMOVED***
	b := make([]byte, 2*len(a.user))
	hex.Encode(b, []byte(a.user))
	return []byte("EXTERNAL"), b, AuthOk
***REMOVED***

func (a authExternal) HandleData(b []byte) ([]byte, AuthStatus) ***REMOVED***
	return nil, AuthError
***REMOVED***
