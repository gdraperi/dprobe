package zk

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// AuthACL produces an ACL list containing a single ACL which uses the
// provided permissions, with the scheme "auth", and ID "", which is used
// by ZooKeeper to represent any authenticated user.
func AuthACL(perms int32) []ACL ***REMOVED***
	return []ACL***REMOVED******REMOVED***perms, "auth", ""***REMOVED******REMOVED***
***REMOVED***

// WorldACL produces an ACL list containing a single ACL which uses the
// provided permissions, with the scheme "world", and ID "anyone", which
// is used by ZooKeeper to represent any user at all.
func WorldACL(perms int32) []ACL ***REMOVED***
	return []ACL***REMOVED******REMOVED***perms, "world", "anyone"***REMOVED******REMOVED***
***REMOVED***

func DigestACL(perms int32, user, password string) []ACL ***REMOVED***
	userPass := []byte(fmt.Sprintf("%s:%s", user, password))
	h := sha1.New()
	if n, err := h.Write(userPass); err != nil || n != len(userPass) ***REMOVED***
		panic("SHA1 failed")
	***REMOVED***
	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return []ACL***REMOVED******REMOVED***perms, "digest", fmt.Sprintf("%s:%s", user, digest)***REMOVED******REMOVED***
***REMOVED***

// FormatServers takes a slice of addresses, and makes sure they are in a format
// that resembles <addr>:<port>. If the server has no port provided, the
// DefaultPort constant is added to the end.
func FormatServers(servers []string) []string ***REMOVED***
	for i := range servers ***REMOVED***
		if !strings.Contains(servers[i], ":") ***REMOVED***
			servers[i] = servers[i] + ":" + strconv.Itoa(DefaultPort)
		***REMOVED***
	***REMOVED***
	return servers
***REMOVED***

// stringShuffle performs a Fisher-Yates shuffle on a slice of strings
func stringShuffle(s []string) ***REMOVED***
	for i := len(s) - 1; i > 0; i-- ***REMOVED***
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	***REMOVED***
***REMOVED***
