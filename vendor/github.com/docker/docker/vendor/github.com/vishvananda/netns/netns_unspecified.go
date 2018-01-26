// +build !linux

package netns

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

func Set(ns NsHandle) (err error) ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func New() (ns NsHandle, err error) ***REMOVED***
	return -1, ErrNotImplemented
***REMOVED***

func Get() (NsHandle, error) ***REMOVED***
	return -1, ErrNotImplemented
***REMOVED***

func GetFromName(name string) (NsHandle, error) ***REMOVED***
	return -1, ErrNotImplemented
***REMOVED***

func GetFromPid(pid int) (NsHandle, error) ***REMOVED***
	return -1, ErrNotImplemented
***REMOVED***

func GetFromDocker(id string) (NsHandle, error) ***REMOVED***
	return -1, ErrNotImplemented
***REMOVED***
