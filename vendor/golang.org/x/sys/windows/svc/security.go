// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package svc

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func allocSid(subAuth0 uint32) (*windows.SID, error) ***REMOVED***
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(&windows.SECURITY_NT_AUTHORITY,
		1, subAuth0, 0, 0, 0, 0, 0, 0, 0, &sid)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return sid, nil
***REMOVED***

// IsAnInteractiveSession determines if calling process is running interactively.
// It queries the process token for membership in the Interactive group.
// http://stackoverflow.com/questions/2668851/how-do-i-detect-that-my-application-is-running-as-service-or-in-an-interactive-s
func IsAnInteractiveSession() (bool, error) ***REMOVED***
	interSid, err := allocSid(windows.SECURITY_INTERACTIVE_RID)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer windows.FreeSid(interSid)

	serviceSid, err := allocSid(windows.SECURITY_SERVICE_RID)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer windows.FreeSid(serviceSid)

	t, err := windows.OpenCurrentProcessToken()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer t.Close()

	gs, err := t.GetTokenGroups()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	p := unsafe.Pointer(&gs.Groups[0])
	groups := (*[2 << 20]windows.SIDAndAttributes)(p)[:gs.GroupCount]
	for _, g := range groups ***REMOVED***
		if windows.EqualSid(g.Sid, interSid) ***REMOVED***
			return true, nil
		***REMOVED***
		if windows.EqualSid(g.Sid, serviceSid) ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***
	return false, nil
***REMOVED***
