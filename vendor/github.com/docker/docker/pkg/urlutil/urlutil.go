// Package urlutil provides helper function to check urls kind.
// It supports http urls, git urls and transport url (tcp://, …)
package urlutil

import (
	"regexp"
	"strings"
)

var (
	validPrefixes = map[string][]string***REMOVED***
		"url":       ***REMOVED***"http://", "https://"***REMOVED***,
		"git":       ***REMOVED***"git://", "github.com/", "git@"***REMOVED***,
		"transport": ***REMOVED***"tcp://", "tcp+tls://", "udp://", "unix://", "unixgram://"***REMOVED***,
	***REMOVED***
	urlPathWithFragmentSuffix = regexp.MustCompile(".git(?:#.+)?$")
)

// IsURL returns true if the provided str is an HTTP(S) URL.
func IsURL(str string) bool ***REMOVED***
	return checkURL(str, "url")
***REMOVED***

// IsGitURL returns true if the provided str is a git repository URL.
func IsGitURL(str string) bool ***REMOVED***
	if IsURL(str) && urlPathWithFragmentSuffix.MatchString(str) ***REMOVED***
		return true
	***REMOVED***
	return checkURL(str, "git")
***REMOVED***

// IsTransportURL returns true if the provided str is a transport (tcp, tcp+tls, udp, unix) URL.
func IsTransportURL(str string) bool ***REMOVED***
	return checkURL(str, "transport")
***REMOVED***

func checkURL(str, kind string) bool ***REMOVED***
	for _, prefix := range validPrefixes[kind] ***REMOVED***
		if strings.HasPrefix(str, prefix) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
