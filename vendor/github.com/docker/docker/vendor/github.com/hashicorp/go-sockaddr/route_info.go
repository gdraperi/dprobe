package sockaddr

// RouteInterface specifies an interface for obtaining memoized route table and
// network information from a given OS.
type RouteInterface interface ***REMOVED***
	// GetDefaultInterfaceName returns the name of the interface that has a
	// default route or an error and an empty string if a problem was
	// encountered.
	GetDefaultInterfaceName() (string, error)
***REMOVED***

// VisitCommands visits each command used by the platform-specific RouteInfo
// implementation.
func (ri routeInfo) VisitCommands(fn func(name string, cmd []string)) ***REMOVED***
	for k, v := range ri.cmds ***REMOVED***
		cmds := append([]string(nil), v...)
		fn(k, cmds)
	***REMOVED***
***REMOVED***
