package sockaddr

import (
	"fmt"
	"strings"
)

type UnixSock struct ***REMOVED***
	SockAddr
	path string
***REMOVED***
type UnixSocks []*UnixSock

// unixAttrMap is a map of the UnixSockAddr type-specific attributes.
var unixAttrMap map[AttrName]func(UnixSock) string
var unixAttrs []AttrName

func init() ***REMOVED***
	unixAttrInit()
***REMOVED***

// NewUnixSock creates an UnixSock from a string path.  String can be in the
// form of either URI-based string (e.g. `file:///etc/passwd`), an absolute
// path (e.g. `/etc/passwd`), or a relative path (e.g. `./foo`).
func NewUnixSock(s string) (ret UnixSock, err error) ***REMOVED***
	ret.path = s
	return ret, nil
***REMOVED***

// CmpAddress follows the Cmp() standard protocol and returns:
//
// - -1 If the receiver should sort first because its name lexically sorts before arg
// - 0 if the SockAddr arg is not a UnixSock, or is a UnixSock with the same path.
// - 1 If the argument should sort first.
func (us UnixSock) CmpAddress(sa SockAddr) int ***REMOVED***
	usb, ok := sa.(UnixSock)
	if !ok ***REMOVED***
		return sortDeferDecision
	***REMOVED***

	return strings.Compare(us.Path(), usb.Path())
***REMOVED***

// DialPacketArgs returns the arguments required to be passed to net.DialUnix()
// with the `unixgram` network type.
func (us UnixSock) DialPacketArgs() (network, dialArgs string) ***REMOVED***
	return "unixgram", us.path
***REMOVED***

// DialStreamArgs returns the arguments required to be passed to net.DialUnix()
// with the `unix` network type.
func (us UnixSock) DialStreamArgs() (network, dialArgs string) ***REMOVED***
	return "unix", us.path
***REMOVED***

// Equal returns true if a SockAddr is equal to the receiving UnixSock.
func (us UnixSock) Equal(sa SockAddr) bool ***REMOVED***
	usb, ok := sa.(UnixSock)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	if us.Path() != usb.Path() ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

// ListenPacketArgs returns the arguments required to be passed to
// net.ListenUnixgram() with the `unixgram` network type.
func (us UnixSock) ListenPacketArgs() (network, dialArgs string) ***REMOVED***
	return "unixgram", us.path
***REMOVED***

// ListenStreamArgs returns the arguments required to be passed to
// net.ListenUnix() with the `unix` network type.
func (us UnixSock) ListenStreamArgs() (network, dialArgs string) ***REMOVED***
	return "unix", us.path
***REMOVED***

// MustUnixSock is a helper method that must return an UnixSock or panic on
// invalid input.
func MustUnixSock(addr string) UnixSock ***REMOVED***
	us, err := NewUnixSock(addr)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("Unable to create a UnixSock from %+q: %v", addr, err))
	***REMOVED***
	return us
***REMOVED***

// Path returns the given path of the UnixSock
func (us UnixSock) Path() string ***REMOVED***
	return us.path
***REMOVED***

// String returns the path of the UnixSock
func (us UnixSock) String() string ***REMOVED***
	return fmt.Sprintf("%+q", us.path)
***REMOVED***

// Type is used as a type switch and returns TypeUnix
func (UnixSock) Type() SockAddrType ***REMOVED***
	return TypeUnix
***REMOVED***

// UnixSockAttrs returns a list of attributes supported by the UnixSockAddr type
func UnixSockAttrs() []AttrName ***REMOVED***
	return unixAttrs
***REMOVED***

// UnixSockAttr returns a string representation of an attribute for the given
// UnixSock.
func UnixSockAttr(us UnixSock, attrName AttrName) string ***REMOVED***
	fn, found := unixAttrMap[attrName]
	if !found ***REMOVED***
		return ""
	***REMOVED***

	return fn(us)
***REMOVED***

// unixAttrInit is called once at init()
func unixAttrInit() ***REMOVED***
	// Sorted for human readability
	unixAttrs = []AttrName***REMOVED***
		"path",
	***REMOVED***

	unixAttrMap = map[AttrName]func(us UnixSock) string***REMOVED***
		"path": func(us UnixSock) string ***REMOVED***
			return us.Path()
		***REMOVED***,
	***REMOVED***
***REMOVED***
