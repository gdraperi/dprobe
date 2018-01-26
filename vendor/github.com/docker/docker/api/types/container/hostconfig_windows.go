package container

// IsBridge indicates whether container uses the bridge network stack
// in windows it is given the name NAT
func (n NetworkMode) IsBridge() bool ***REMOVED***
	return n == "nat"
***REMOVED***

// IsHost indicates whether container uses the host network stack.
// returns false as this is not supported by windows
func (n NetworkMode) IsHost() bool ***REMOVED***
	return false
***REMOVED***

// IsUserDefined indicates user-created network
func (n NetworkMode) IsUserDefined() bool ***REMOVED***
	return !n.IsDefault() && !n.IsNone() && !n.IsBridge() && !n.IsContainer()
***REMOVED***

// IsValid indicates if an isolation technology is valid
func (i Isolation) IsValid() bool ***REMOVED***
	return i.IsDefault() || i.IsHyperV() || i.IsProcess()
***REMOVED***

// NetworkName returns the name of the network stack.
func (n NetworkMode) NetworkName() string ***REMOVED***
	if n.IsDefault() ***REMOVED***
		return "default"
	***REMOVED*** else if n.IsBridge() ***REMOVED***
		return "nat"
	***REMOVED*** else if n.IsNone() ***REMOVED***
		return "none"
	***REMOVED*** else if n.IsContainer() ***REMOVED***
		return "container"
	***REMOVED*** else if n.IsUserDefined() ***REMOVED***
		return n.UserDefined()
	***REMOVED***

	return ""
***REMOVED***
