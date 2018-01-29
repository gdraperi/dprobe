// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package mgr can be used to manage Windows service programs.
// It can be used to install and remove them. It can also start,
// stop and pause them. The package can query / change current
// service state and config parameters.
//
package mgr

import (
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Mgr is used to manage Windows service.
type Mgr struct ***REMOVED***
	Handle windows.Handle
***REMOVED***

// Connect establishes a connection to the service control manager.
func Connect() (*Mgr, error) ***REMOVED***
	return ConnectRemote("")
***REMOVED***

// ConnectRemote establishes a connection to the
// service control manager on computer named host.
func ConnectRemote(host string) (*Mgr, error) ***REMOVED***
	var s *uint16
	if host != "" ***REMOVED***
		s = syscall.StringToUTF16Ptr(host)
	***REMOVED***
	h, err := windows.OpenSCManager(s, nil, windows.SC_MANAGER_ALL_ACCESS)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Mgr***REMOVED***Handle: h***REMOVED***, nil
***REMOVED***

// Disconnect closes connection to the service control manager m.
func (m *Mgr) Disconnect() error ***REMOVED***
	return windows.CloseServiceHandle(m.Handle)
***REMOVED***

func toPtr(s string) *uint16 ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return nil
	***REMOVED***
	return syscall.StringToUTF16Ptr(s)
***REMOVED***

// toStringBlock terminates strings in ss with 0, and then
// concatenates them together. It also adds extra 0 at the end.
func toStringBlock(ss []string) *uint16 ***REMOVED***
	if len(ss) == 0 ***REMOVED***
		return nil
	***REMOVED***
	t := ""
	for _, s := range ss ***REMOVED***
		if s != "" ***REMOVED***
			t += s + "\x00"
		***REMOVED***
	***REMOVED***
	if t == "" ***REMOVED***
		return nil
	***REMOVED***
	t += "\x00"
	return &utf16.Encode([]rune(t))[0]
***REMOVED***

// CreateService installs new service name on the system.
// The service will be executed by running exepath binary.
// Use config c to specify service parameters.
// Any args will be passed as command-line arguments when
// the service is started; these arguments are distinct from
// the arguments passed to Service.Start or via the "Start
// parameters" field in the service's Properties dialog box.
func (m *Mgr) CreateService(name, exepath string, c Config, args ...string) (*Service, error) ***REMOVED***
	if c.StartType == 0 ***REMOVED***
		c.StartType = StartManual
	***REMOVED***
	if c.ErrorControl == 0 ***REMOVED***
		c.ErrorControl = ErrorNormal
	***REMOVED***
	if c.ServiceType == 0 ***REMOVED***
		c.ServiceType = windows.SERVICE_WIN32_OWN_PROCESS
	***REMOVED***
	s := syscall.EscapeArg(exepath)
	for _, v := range args ***REMOVED***
		s += " " + syscall.EscapeArg(v)
	***REMOVED***
	h, err := windows.CreateService(m.Handle, toPtr(name), toPtr(c.DisplayName),
		windows.SERVICE_ALL_ACCESS, c.ServiceType,
		c.StartType, c.ErrorControl, toPtr(s), toPtr(c.LoadOrderGroup),
		nil, toStringBlock(c.Dependencies), toPtr(c.ServiceStartName), toPtr(c.Password))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if c.Description != "" ***REMOVED***
		err = updateDescription(h, c.Description)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &Service***REMOVED***Name: name, Handle: h***REMOVED***, nil
***REMOVED***

// OpenService retrieves access to service name, so it can
// be interrogated and controlled.
func (m *Mgr) OpenService(name string) (*Service, error) ***REMOVED***
	h, err := windows.OpenService(m.Handle, syscall.StringToUTF16Ptr(name), windows.SERVICE_ALL_ACCESS)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Service***REMOVED***Name: name, Handle: h***REMOVED***, nil
***REMOVED***

// ListServices enumerates services in the specified
// service control manager database m.
// If the caller does not have the SERVICE_QUERY_STATUS
// access right to a service, the service is silently
// omitted from the list of services returned.
func (m *Mgr) ListServices() ([]string, error) ***REMOVED***
	var err error
	var bytesNeeded, servicesReturned uint32
	var buf []byte
	for ***REMOVED***
		var p *byte
		if len(buf) > 0 ***REMOVED***
			p = &buf[0]
		***REMOVED***
		err = windows.EnumServicesStatusEx(m.Handle, windows.SC_ENUM_PROCESS_INFO,
			windows.SERVICE_WIN32, windows.SERVICE_STATE_ALL,
			p, uint32(len(buf)), &bytesNeeded, &servicesReturned, nil, nil)
		if err == nil ***REMOVED***
			break
		***REMOVED***
		if err != syscall.ERROR_MORE_DATA ***REMOVED***
			return nil, err
		***REMOVED***
		if bytesNeeded <= uint32(len(buf)) ***REMOVED***
			return nil, err
		***REMOVED***
		buf = make([]byte, bytesNeeded)
	***REMOVED***
	if servicesReturned == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	services := (*[1 << 20]windows.ENUM_SERVICE_STATUS_PROCESS)(unsafe.Pointer(&buf[0]))[:servicesReturned]
	var names []string
	for _, s := range services ***REMOVED***
		name := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(s.ServiceName))[:])
		names = append(names, name)
	***REMOVED***
	return names, nil
***REMOVED***
