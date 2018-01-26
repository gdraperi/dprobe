package user

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	minId = 0
	maxId = 1<<31 - 1 //for 32-bit systems compatibility
)

var (
	ErrRange = fmt.Errorf("uids and gids must be in range %d-%d", minId, maxId)
)

type User struct ***REMOVED***
	Name  string
	Pass  string
	Uid   int
	Gid   int
	Gecos string
	Home  string
	Shell string
***REMOVED***

type Group struct ***REMOVED***
	Name string
	Pass string
	Gid  int
	List []string
***REMOVED***

func parseLine(line string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	if line == "" ***REMOVED***
		return
	***REMOVED***

	parts := strings.Split(line, ":")
	for i, p := range parts ***REMOVED***
		// Ignore cases where we don't have enough fields to populate the arguments.
		// Some configuration files like to misbehave.
		if len(v) <= i ***REMOVED***
			break
		***REMOVED***

		// Use the type of the argument to figure out how to parse it, scanf() style.
		// This is legit.
		switch e := v[i].(type) ***REMOVED***
		case *string:
			*e = p
		case *int:
			// "numbers", with conversion errors ignored because of some misbehaving configuration files.
			*e, _ = strconv.Atoi(p)
		case *[]string:
			// Comma-separated lists.
			if p != "" ***REMOVED***
				*e = strings.Split(p, ",")
			***REMOVED*** else ***REMOVED***
				*e = []string***REMOVED******REMOVED***
			***REMOVED***
		default:
			// Someone goof'd when writing code using this function. Scream so they can hear us.
			panic(fmt.Sprintf("parseLine only accepts ***REMOVED****string, *int, *[]string***REMOVED*** as arguments! %#v is not a pointer!", e))
		***REMOVED***
	***REMOVED***
***REMOVED***

func ParsePasswdFile(path string) ([]User, error) ***REMOVED***
	passwd, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer passwd.Close()
	return ParsePasswd(passwd)
***REMOVED***

func ParsePasswd(passwd io.Reader) ([]User, error) ***REMOVED***
	return ParsePasswdFilter(passwd, nil)
***REMOVED***

func ParsePasswdFileFilter(path string, filter func(User) bool) ([]User, error) ***REMOVED***
	passwd, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer passwd.Close()
	return ParsePasswdFilter(passwd, filter)
***REMOVED***

func ParsePasswdFilter(r io.Reader, filter func(User) bool) ([]User, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, fmt.Errorf("nil source for passwd-formatted data")
	***REMOVED***

	var (
		s   = bufio.NewScanner(r)
		out = []User***REMOVED******REMOVED***
	)

	for s.Scan() ***REMOVED***
		if err := s.Err(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		line := strings.TrimSpace(s.Text())
		if line == "" ***REMOVED***
			continue
		***REMOVED***

		// see: man 5 passwd
		//  name:password:UID:GID:GECOS:directory:shell
		// Name:Pass:Uid:Gid:Gecos:Home:Shell
		//  root:x:0:0:root:/root:/bin/bash
		//  adm:x:3:4:adm:/var/adm:/bin/false
		p := User***REMOVED******REMOVED***
		parseLine(line, &p.Name, &p.Pass, &p.Uid, &p.Gid, &p.Gecos, &p.Home, &p.Shell)

		if filter == nil || filter(p) ***REMOVED***
			out = append(out, p)
		***REMOVED***
	***REMOVED***

	return out, nil
***REMOVED***

func ParseGroupFile(path string) ([]Group, error) ***REMOVED***
	group, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer group.Close()
	return ParseGroup(group)
***REMOVED***

func ParseGroup(group io.Reader) ([]Group, error) ***REMOVED***
	return ParseGroupFilter(group, nil)
***REMOVED***

func ParseGroupFileFilter(path string, filter func(Group) bool) ([]Group, error) ***REMOVED***
	group, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer group.Close()
	return ParseGroupFilter(group, filter)
***REMOVED***

func ParseGroupFilter(r io.Reader, filter func(Group) bool) ([]Group, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, fmt.Errorf("nil source for group-formatted data")
	***REMOVED***

	var (
		s   = bufio.NewScanner(r)
		out = []Group***REMOVED******REMOVED***
	)

	for s.Scan() ***REMOVED***
		if err := s.Err(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		text := s.Text()
		if text == "" ***REMOVED***
			continue
		***REMOVED***

		// see: man 5 group
		//  group_name:password:GID:user_list
		// Name:Pass:Gid:List
		//  root:x:0:root
		//  adm:x:4:root,adm,daemon
		p := Group***REMOVED******REMOVED***
		parseLine(text, &p.Name, &p.Pass, &p.Gid, &p.List)

		if filter == nil || filter(p) ***REMOVED***
			out = append(out, p)
		***REMOVED***
	***REMOVED***

	return out, nil
***REMOVED***

type ExecUser struct ***REMOVED***
	Uid   int
	Gid   int
	Sgids []int
	Home  string
***REMOVED***

// GetExecUserPath is a wrapper for GetExecUser. It reads data from each of the
// given file paths and uses that data as the arguments to GetExecUser. If the
// files cannot be opened for any reason, the error is ignored and a nil
// io.Reader is passed instead.
func GetExecUserPath(userSpec string, defaults *ExecUser, passwdPath, groupPath string) (*ExecUser, error) ***REMOVED***
	var passwd, group io.Reader

	if passwdFile, err := os.Open(passwdPath); err == nil ***REMOVED***
		passwd = passwdFile
		defer passwdFile.Close()
	***REMOVED***

	if groupFile, err := os.Open(groupPath); err == nil ***REMOVED***
		group = groupFile
		defer groupFile.Close()
	***REMOVED***

	return GetExecUser(userSpec, defaults, passwd, group)
***REMOVED***

// GetExecUser parses a user specification string (using the passwd and group
// readers as sources for /etc/passwd and /etc/group data, respectively). In
// the case of blank fields or missing data from the sources, the values in
// defaults is used.
//
// GetExecUser will return an error if a user or group literal could not be
// found in any entry in passwd and group respectively.
//
// Examples of valid user specifications are:
//     * ""
//     * "user"
//     * "uid"
//     * "user:group"
//     * "uid:gid
//     * "user:gid"
//     * "uid:group"
//
// It should be noted that if you specify a numeric user or group id, they will
// not be evaluated as usernames (only the metadata will be filled). So attempting
// to parse a user with user.Name = "1337" will produce the user with a UID of
// 1337.
func GetExecUser(userSpec string, defaults *ExecUser, passwd, group io.Reader) (*ExecUser, error) ***REMOVED***
	if defaults == nil ***REMOVED***
		defaults = new(ExecUser)
	***REMOVED***

	// Copy over defaults.
	user := &ExecUser***REMOVED***
		Uid:   defaults.Uid,
		Gid:   defaults.Gid,
		Sgids: defaults.Sgids,
		Home:  defaults.Home,
	***REMOVED***

	// Sgids slice *cannot* be nil.
	if user.Sgids == nil ***REMOVED***
		user.Sgids = []int***REMOVED******REMOVED***
	***REMOVED***

	// Allow for userArg to have either "user" syntax, or optionally "user:group" syntax
	var userArg, groupArg string
	parseLine(userSpec, &userArg, &groupArg)

	// Convert userArg and groupArg to be numeric, so we don't have to execute
	// Atoi *twice* for each iteration over lines.
	uidArg, uidErr := strconv.Atoi(userArg)
	gidArg, gidErr := strconv.Atoi(groupArg)

	// Find the matching user.
	users, err := ParsePasswdFilter(passwd, func(u User) bool ***REMOVED***
		if userArg == "" ***REMOVED***
			// Default to current state of the user.
			return u.Uid == user.Uid
		***REMOVED***

		if uidErr == nil ***REMOVED***
			// If the userArg is numeric, always treat it as a UID.
			return uidArg == u.Uid
		***REMOVED***

		return u.Name == userArg
	***REMOVED***)

	// If we can't find the user, we have to bail.
	if err != nil && passwd != nil ***REMOVED***
		if userArg == "" ***REMOVED***
			userArg = strconv.Itoa(user.Uid)
		***REMOVED***
		return nil, fmt.Errorf("unable to find user %s: %v", userArg, err)
	***REMOVED***

	var matchedUserName string
	if len(users) > 0 ***REMOVED***
		// First match wins, even if there's more than one matching entry.
		matchedUserName = users[0].Name
		user.Uid = users[0].Uid
		user.Gid = users[0].Gid
		user.Home = users[0].Home
	***REMOVED*** else if userArg != "" ***REMOVED***
		// If we can't find a user with the given username, the only other valid
		// option is if it's a numeric username with no associated entry in passwd.

		if uidErr != nil ***REMOVED***
			// Not numeric.
			return nil, fmt.Errorf("unable to find user %s: %v", userArg, ErrNoPasswdEntries)
		***REMOVED***
		user.Uid = uidArg

		// Must be inside valid uid range.
		if user.Uid < minId || user.Uid > maxId ***REMOVED***
			return nil, ErrRange
		***REMOVED***

		// Okay, so it's numeric. We can just roll with this.
	***REMOVED***

	// On to the groups. If we matched a username, we need to do this because of
	// the supplementary group IDs.
	if groupArg != "" || matchedUserName != "" ***REMOVED***
		groups, err := ParseGroupFilter(group, func(g Group) bool ***REMOVED***
			// If the group argument isn't explicit, we'll just search for it.
			if groupArg == "" ***REMOVED***
				// Check if user is a member of this group.
				for _, u := range g.List ***REMOVED***
					if u == matchedUserName ***REMOVED***
						return true
					***REMOVED***
				***REMOVED***
				return false
			***REMOVED***

			if gidErr == nil ***REMOVED***
				// If the groupArg is numeric, always treat it as a GID.
				return gidArg == g.Gid
			***REMOVED***

			return g.Name == groupArg
		***REMOVED***)
		if err != nil && group != nil ***REMOVED***
			return nil, fmt.Errorf("unable to find groups for spec %v: %v", matchedUserName, err)
		***REMOVED***

		// Only start modifying user.Gid if it is in explicit form.
		if groupArg != "" ***REMOVED***
			if len(groups) > 0 ***REMOVED***
				// First match wins, even if there's more than one matching entry.
				user.Gid = groups[0].Gid
			***REMOVED*** else ***REMOVED***
				// If we can't find a group with the given name, the only other valid
				// option is if it's a numeric group name with no associated entry in group.

				if gidErr != nil ***REMOVED***
					// Not numeric.
					return nil, fmt.Errorf("unable to find group %s: %v", groupArg, ErrNoGroupEntries)
				***REMOVED***
				user.Gid = gidArg

				// Must be inside valid gid range.
				if user.Gid < minId || user.Gid > maxId ***REMOVED***
					return nil, ErrRange
				***REMOVED***

				// Okay, so it's numeric. We can just roll with this.
			***REMOVED***
		***REMOVED*** else if len(groups) > 0 ***REMOVED***
			// Supplementary group ids only make sense if in the implicit form.
			user.Sgids = make([]int, len(groups))
			for i, group := range groups ***REMOVED***
				user.Sgids[i] = group.Gid
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return user, nil
***REMOVED***

// GetAdditionalGroups looks up a list of groups by name or group id
// against the given /etc/group formatted data. If a group name cannot
// be found, an error will be returned. If a group id cannot be found,
// or the given group data is nil, the id will be returned as-is
// provided it is in the legal range.
func GetAdditionalGroups(additionalGroups []string, group io.Reader) ([]int, error) ***REMOVED***
	var groups = []Group***REMOVED******REMOVED***
	if group != nil ***REMOVED***
		var err error
		groups, err = ParseGroupFilter(group, func(g Group) bool ***REMOVED***
			for _, ag := range additionalGroups ***REMOVED***
				if g.Name == ag || strconv.Itoa(g.Gid) == ag ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
			return false
		***REMOVED***)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Unable to find additional groups %v: %v", additionalGroups, err)
		***REMOVED***
	***REMOVED***

	gidMap := make(map[int]struct***REMOVED******REMOVED***)
	for _, ag := range additionalGroups ***REMOVED***
		var found bool
		for _, g := range groups ***REMOVED***
			// if we found a matched group either by name or gid, take the
			// first matched as correct
			if g.Name == ag || strconv.Itoa(g.Gid) == ag ***REMOVED***
				if _, ok := gidMap[g.Gid]; !ok ***REMOVED***
					gidMap[g.Gid] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					found = true
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// we asked for a group but didn't find it. let's check to see
		// if we wanted a numeric group
		if !found ***REMOVED***
			gid, err := strconv.Atoi(ag)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("Unable to find group %s", ag)
			***REMOVED***
			// Ensure gid is inside gid range.
			if gid < minId || gid > maxId ***REMOVED***
				return nil, ErrRange
			***REMOVED***
			gidMap[gid] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	gids := []int***REMOVED******REMOVED***
	for gid := range gidMap ***REMOVED***
		gids = append(gids, gid)
	***REMOVED***
	return gids, nil
***REMOVED***

// GetAdditionalGroupsPath is a wrapper around GetAdditionalGroups
// that opens the groupPath given and gives it as an argument to
// GetAdditionalGroups.
func GetAdditionalGroupsPath(additionalGroups []string, groupPath string) ([]int, error) ***REMOVED***
	var group io.Reader

	if groupFile, err := os.Open(groupPath); err == nil ***REMOVED***
		group = groupFile
		defer groupFile.Close()
	***REMOVED***
	return GetAdditionalGroups(additionalGroups, group)
***REMOVED***
