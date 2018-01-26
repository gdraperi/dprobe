package idtools

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// IDMap contains a single entry for user namespace range remapping. An array
// of IDMap entries represents the structure that will be provided to the Linux
// kernel for creating a user namespace.
type IDMap struct ***REMOVED***
	ContainerID int `json:"container_id"`
	HostID      int `json:"host_id"`
	Size        int `json:"size"`
***REMOVED***

type subIDRange struct ***REMOVED***
	Start  int
	Length int
***REMOVED***

type ranges []subIDRange

func (e ranges) Len() int           ***REMOVED*** return len(e) ***REMOVED***
func (e ranges) Swap(i, j int)      ***REMOVED*** e[i], e[j] = e[j], e[i] ***REMOVED***
func (e ranges) Less(i, j int) bool ***REMOVED*** return e[i].Start < e[j].Start ***REMOVED***

const (
	subuidFileName string = "/etc/subuid"
	subgidFileName string = "/etc/subgid"
)

// MkdirAllAndChown creates a directory (include any along the path) and then modifies
// ownership to the requested uid/gid.  If the directory already exists, this
// function will still change ownership to the requested uid/gid pair.
func MkdirAllAndChown(path string, mode os.FileMode, owner IDPair) error ***REMOVED***
	return mkdirAs(path, mode, owner.UID, owner.GID, true, true)
***REMOVED***

// MkdirAndChown creates a directory and then modifies ownership to the requested uid/gid.
// If the directory already exists, this function still changes ownership.
// Note that unlike os.Mkdir(), this function does not return IsExist error
// in case path already exists.
func MkdirAndChown(path string, mode os.FileMode, owner IDPair) error ***REMOVED***
	return mkdirAs(path, mode, owner.UID, owner.GID, false, true)
***REMOVED***

// MkdirAllAndChownNew creates a directory (include any along the path) and then modifies
// ownership ONLY of newly created directories to the requested uid/gid. If the
// directories along the path exist, no change of ownership will be performed
func MkdirAllAndChownNew(path string, mode os.FileMode, owner IDPair) error ***REMOVED***
	return mkdirAs(path, mode, owner.UID, owner.GID, true, false)
***REMOVED***

// GetRootUIDGID retrieves the remapped root uid/gid pair from the set of maps.
// If the maps are empty, then the root uid/gid will default to "real" 0/0
func GetRootUIDGID(uidMap, gidMap []IDMap) (int, int, error) ***REMOVED***
	uid, err := toHost(0, uidMap)
	if err != nil ***REMOVED***
		return -1, -1, err
	***REMOVED***
	gid, err := toHost(0, gidMap)
	if err != nil ***REMOVED***
		return -1, -1, err
	***REMOVED***
	return uid, gid, nil
***REMOVED***

// toContainer takes an id mapping, and uses it to translate a
// host ID to the remapped ID. If no map is provided, then the translation
// assumes a 1-to-1 mapping and returns the passed in id
func toContainer(hostID int, idMap []IDMap) (int, error) ***REMOVED***
	if idMap == nil ***REMOVED***
		return hostID, nil
	***REMOVED***
	for _, m := range idMap ***REMOVED***
		if (hostID >= m.HostID) && (hostID <= (m.HostID + m.Size - 1)) ***REMOVED***
			contID := m.ContainerID + (hostID - m.HostID)
			return contID, nil
		***REMOVED***
	***REMOVED***
	return -1, fmt.Errorf("Host ID %d cannot be mapped to a container ID", hostID)
***REMOVED***

// toHost takes an id mapping and a remapped ID, and translates the
// ID to the mapped host ID. If no map is provided, then the translation
// assumes a 1-to-1 mapping and returns the passed in id #
func toHost(contID int, idMap []IDMap) (int, error) ***REMOVED***
	if idMap == nil ***REMOVED***
		return contID, nil
	***REMOVED***
	for _, m := range idMap ***REMOVED***
		if (contID >= m.ContainerID) && (contID <= (m.ContainerID + m.Size - 1)) ***REMOVED***
			hostID := m.HostID + (contID - m.ContainerID)
			return hostID, nil
		***REMOVED***
	***REMOVED***
	return -1, fmt.Errorf("Container ID %d cannot be mapped to a host ID", contID)
***REMOVED***

// IDPair is a UID and GID pair
type IDPair struct ***REMOVED***
	UID int
	GID int
***REMOVED***

// IDMappings contains a mappings of UIDs and GIDs
type IDMappings struct ***REMOVED***
	uids []IDMap
	gids []IDMap
***REMOVED***

// NewIDMappings takes a requested user and group name and
// using the data from /etc/sub***REMOVED***uid,gid***REMOVED*** ranges, creates the
// proper uid and gid remapping ranges for that user/group pair
func NewIDMappings(username, groupname string) (*IDMappings, error) ***REMOVED***
	subuidRanges, err := parseSubuid(username)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	subgidRanges, err := parseSubgid(groupname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(subuidRanges) == 0 ***REMOVED***
		return nil, fmt.Errorf("No subuid ranges found for user %q", username)
	***REMOVED***
	if len(subgidRanges) == 0 ***REMOVED***
		return nil, fmt.Errorf("No subgid ranges found for group %q", groupname)
	***REMOVED***

	return &IDMappings***REMOVED***
		uids: createIDMap(subuidRanges),
		gids: createIDMap(subgidRanges),
	***REMOVED***, nil
***REMOVED***

// NewIDMappingsFromMaps creates a new mapping from two slices
// Deprecated: this is a temporary shim while transitioning to IDMapping
func NewIDMappingsFromMaps(uids []IDMap, gids []IDMap) *IDMappings ***REMOVED***
	return &IDMappings***REMOVED***uids: uids, gids: gids***REMOVED***
***REMOVED***

// RootPair returns a uid and gid pair for the root user. The error is ignored
// because a root user always exists, and the defaults are correct when the uid
// and gid maps are empty.
func (i *IDMappings) RootPair() IDPair ***REMOVED***
	uid, gid, _ := GetRootUIDGID(i.uids, i.gids)
	return IDPair***REMOVED***UID: uid, GID: gid***REMOVED***
***REMOVED***

// ToHost returns the host UID and GID for the container uid, gid.
// Remapping is only performed if the ids aren't already the remapped root ids
func (i *IDMappings) ToHost(pair IDPair) (IDPair, error) ***REMOVED***
	var err error
	target := i.RootPair()

	if pair.UID != target.UID ***REMOVED***
		target.UID, err = toHost(pair.UID, i.uids)
		if err != nil ***REMOVED***
			return target, err
		***REMOVED***
	***REMOVED***

	if pair.GID != target.GID ***REMOVED***
		target.GID, err = toHost(pair.GID, i.gids)
	***REMOVED***
	return target, err
***REMOVED***

// ToContainer returns the container UID and GID for the host uid and gid
func (i *IDMappings) ToContainer(pair IDPair) (int, int, error) ***REMOVED***
	uid, err := toContainer(pair.UID, i.uids)
	if err != nil ***REMOVED***
		return -1, -1, err
	***REMOVED***
	gid, err := toContainer(pair.GID, i.gids)
	return uid, gid, err
***REMOVED***

// Empty returns true if there are no id mappings
func (i *IDMappings) Empty() bool ***REMOVED***
	return len(i.uids) == 0 && len(i.gids) == 0
***REMOVED***

// UIDs return the UID mapping
// TODO: remove this once everything has been refactored to use pairs
func (i *IDMappings) UIDs() []IDMap ***REMOVED***
	return i.uids
***REMOVED***

// GIDs return the UID mapping
// TODO: remove this once everything has been refactored to use pairs
func (i *IDMappings) GIDs() []IDMap ***REMOVED***
	return i.gids
***REMOVED***

func createIDMap(subidRanges ranges) []IDMap ***REMOVED***
	idMap := []IDMap***REMOVED******REMOVED***

	// sort the ranges by lowest ID first
	sort.Sort(subidRanges)
	containerID := 0
	for _, idrange := range subidRanges ***REMOVED***
		idMap = append(idMap, IDMap***REMOVED***
			ContainerID: containerID,
			HostID:      idrange.Start,
			Size:        idrange.Length,
		***REMOVED***)
		containerID = containerID + idrange.Length
	***REMOVED***
	return idMap
***REMOVED***

func parseSubuid(username string) (ranges, error) ***REMOVED***
	return parseSubidFile(subuidFileName, username)
***REMOVED***

func parseSubgid(username string) (ranges, error) ***REMOVED***
	return parseSubidFile(subgidFileName, username)
***REMOVED***

// parseSubidFile will read the appropriate file (/etc/subuid or /etc/subgid)
// and return all found ranges for a specified username. If the special value
// "ALL" is supplied for username, then all ranges in the file will be returned
func parseSubidFile(path, username string) (ranges, error) ***REMOVED***
	var rangeList ranges

	subidFile, err := os.Open(path)
	if err != nil ***REMOVED***
		return rangeList, err
	***REMOVED***
	defer subidFile.Close()

	s := bufio.NewScanner(subidFile)
	for s.Scan() ***REMOVED***
		if err := s.Err(); err != nil ***REMOVED***
			return rangeList, err
		***REMOVED***

		text := strings.TrimSpace(s.Text())
		if text == "" || strings.HasPrefix(text, "#") ***REMOVED***
			continue
		***REMOVED***
		parts := strings.Split(text, ":")
		if len(parts) != 3 ***REMOVED***
			return rangeList, fmt.Errorf("Cannot parse subuid/gid information: Format not correct for %s file", path)
		***REMOVED***
		if parts[0] == username || username == "ALL" ***REMOVED***
			startid, err := strconv.Atoi(parts[1])
			if err != nil ***REMOVED***
				return rangeList, fmt.Errorf("String to int conversion failed during subuid/gid parsing of %s: %v", path, err)
			***REMOVED***
			length, err := strconv.Atoi(parts[2])
			if err != nil ***REMOVED***
				return rangeList, fmt.Errorf("String to int conversion failed during subuid/gid parsing of %s: %v", path, err)
			***REMOVED***
			rangeList = append(rangeList, subIDRange***REMOVED***startid, length***REMOVED***)
		***REMOVED***
	***REMOVED***
	return rangeList, nil
***REMOVED***
