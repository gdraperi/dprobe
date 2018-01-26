// +build !apparmor !linux

package apparmor

import (
	"errors"
)

var ErrApparmorNotEnabled = errors.New("apparmor: config provided but apparmor not supported")

func IsEnabled() bool ***REMOVED***
	return false
***REMOVED***

func ApplyProfile(name string) error ***REMOVED***
	if name != "" ***REMOVED***
		return ErrApparmorNotEnabled
	***REMOVED***
	return nil
***REMOVED***
