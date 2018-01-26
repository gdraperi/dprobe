// +build !windows

package container

// IsValid indicates if an isolation technology is valid
func (i Isolation) IsValid() bool ***REMOVED***
	return i.IsDefault()
***REMOVED***

// NetworkName returns the name of the network stack.
func (n NetworkMode) NetworkName() string ***REMOVED***
	if n.IsBridge() ***REMOVED***
		return "bridge"
	***REMOVED*** else if n.IsHost() ***REMOVED***
		return "host"
	***REMOVED*** else if n.IsContainer() ***REMOVED***
		return "container"
	***REMOVED*** else if n.IsNone() ***REMOVED***
		return "none"
	***REMOVED*** else if n.IsDefault() ***REMOVED***
		return "default"
	***REMOVED*** else if n.IsUserDefined() ***REMOVED***
		return n.UserDefined()
	***REMOVED***
	return ""
***REMOVED***

// IsBridge indicates whether container uses the bridge network stack
func (n NetworkMode) IsBridge() bool ***REMOVED***
	return n == "bridge"
***REMOVED***

// IsHost indicates whether container uses the host network stack.
func (n NetworkMode) IsHost() bool ***REMOVED***
	return n == "host"
***REMOVED***

// IsUserDefined indicates user-created network
func (n NetworkMode) IsUserDefined() bool ***REMOVED***
	return !n.IsDefault() && !n.IsBridge() && !n.IsHost() && !n.IsNone() && !n.IsContainer()
***REMOVED***
