// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package mgr

import (
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
)

// TODO(brainman): Use EnumDependentServices to enumerate dependent services.

// Service is used to access Windows service.
type Service struct ***REMOVED***
	Name   string
	Handle windows.Handle
***REMOVED***

// Delete marks service s for deletion from the service control manager database.
func (s *Service) Delete() error ***REMOVED***
	return windows.DeleteService(s.Handle)
***REMOVED***

// Close relinquish access to the service s.
func (s *Service) Close() error ***REMOVED***
	return windows.CloseServiceHandle(s.Handle)
***REMOVED***

// Start starts service s.
// args will be passed to svc.Handler.Execute.
func (s *Service) Start(args ...string) error ***REMOVED***
	var p **uint16
	if len(args) > 0 ***REMOVED***
		vs := make([]*uint16, len(args))
		for i := range vs ***REMOVED***
			vs[i] = syscall.StringToUTF16Ptr(args[i])
		***REMOVED***
		p = &vs[0]
	***REMOVED***
	return windows.StartService(s.Handle, uint32(len(args)), p)
***REMOVED***

// Control sends state change request c to the servce s.
func (s *Service) Control(c svc.Cmd) (svc.Status, error) ***REMOVED***
	var t windows.SERVICE_STATUS
	err := windows.ControlService(s.Handle, uint32(c), &t)
	if err != nil ***REMOVED***
		return svc.Status***REMOVED******REMOVED***, err
	***REMOVED***
	return svc.Status***REMOVED***
		State:   svc.State(t.CurrentState),
		Accepts: svc.Accepted(t.ControlsAccepted),
	***REMOVED***, nil
***REMOVED***

// Query returns current status of service s.
func (s *Service) Query() (svc.Status, error) ***REMOVED***
	var t windows.SERVICE_STATUS
	err := windows.QueryServiceStatus(s.Handle, &t)
	if err != nil ***REMOVED***
		return svc.Status***REMOVED******REMOVED***, err
	***REMOVED***
	return svc.Status***REMOVED***
		State:   svc.State(t.CurrentState),
		Accepts: svc.Accepted(t.ControlsAccepted),
	***REMOVED***, nil
***REMOVED***
