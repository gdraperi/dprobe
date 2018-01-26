package user

import (
	"errors"
)

var (
	// The current operating system does not provide the required data for user lookups.
	ErrUnsupported = errors.New("user lookup: operating system does not provide passwd-formatted data")
	// No matching entries found in file.
	ErrNoPasswdEntries = errors.New("no matching entries in passwd file")
	ErrNoGroupEntries  = errors.New("no matching entries in group file")
)

func lookupUser(filter func(u User) bool) (User, error) ***REMOVED***
	// Get operating system-specific passwd reader-closer.
	passwd, err := GetPasswd()
	if err != nil ***REMOVED***
		return User***REMOVED******REMOVED***, err
	***REMOVED***
	defer passwd.Close()

	// Get the users.
	users, err := ParsePasswdFilter(passwd, filter)
	if err != nil ***REMOVED***
		return User***REMOVED******REMOVED***, err
	***REMOVED***

	// No user entries found.
	if len(users) == 0 ***REMOVED***
		return User***REMOVED******REMOVED***, ErrNoPasswdEntries
	***REMOVED***

	// Assume the first entry is the "correct" one.
	return users[0], nil
***REMOVED***

// LookupUser looks up a user by their username in /etc/passwd. If the user
// cannot be found (or there is no /etc/passwd file on the filesystem), then
// LookupUser returns an error.
func LookupUser(username string) (User, error) ***REMOVED***
	return lookupUser(func(u User) bool ***REMOVED***
		return u.Name == username
	***REMOVED***)
***REMOVED***

// LookupUid looks up a user by their user id in /etc/passwd. If the user cannot
// be found (or there is no /etc/passwd file on the filesystem), then LookupId
// returns an error.
func LookupUid(uid int) (User, error) ***REMOVED***
	return lookupUser(func(u User) bool ***REMOVED***
		return u.Uid == uid
	***REMOVED***)
***REMOVED***

func lookupGroup(filter func(g Group) bool) (Group, error) ***REMOVED***
	// Get operating system-specific group reader-closer.
	group, err := GetGroup()
	if err != nil ***REMOVED***
		return Group***REMOVED******REMOVED***, err
	***REMOVED***
	defer group.Close()

	// Get the users.
	groups, err := ParseGroupFilter(group, filter)
	if err != nil ***REMOVED***
		return Group***REMOVED******REMOVED***, err
	***REMOVED***

	// No user entries found.
	if len(groups) == 0 ***REMOVED***
		return Group***REMOVED******REMOVED***, ErrNoGroupEntries
	***REMOVED***

	// Assume the first entry is the "correct" one.
	return groups[0], nil
***REMOVED***

// LookupGroup looks up a group by its name in /etc/group. If the group cannot
// be found (or there is no /etc/group file on the filesystem), then LookupGroup
// returns an error.
func LookupGroup(groupname string) (Group, error) ***REMOVED***
	return lookupGroup(func(g Group) bool ***REMOVED***
		return g.Name == groupname
	***REMOVED***)
***REMOVED***

// LookupGid looks up a group by its group id in /etc/group. If the group cannot
// be found (or there is no /etc/group file on the filesystem), then LookupGid
// returns an error.
func LookupGid(gid int) (Group, error) ***REMOVED***
	return lookupGroup(func(g Group) bool ***REMOVED***
		return g.Gid == gid
	***REMOVED***)
***REMOVED***
