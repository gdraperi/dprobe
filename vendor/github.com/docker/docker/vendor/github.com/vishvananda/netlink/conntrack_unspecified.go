// +build !linux

package netlink

// ConntrackTableType Conntrack table for the netlink operation
type ConntrackTableType uint8

// InetFamily Family type
type InetFamily uint8

// ConntrackFlow placeholder
type ConntrackFlow struct***REMOVED******REMOVED***

// ConntrackFilter placeholder
type ConntrackFilter struct***REMOVED******REMOVED***

// ConntrackTableList returns the flow list of a table of a specific family
// conntrack -L [table] [options]          List conntrack or expectation table
func ConntrackTableList(table ConntrackTableType, family InetFamily) ([]*ConntrackFlow, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

// ConntrackTableFlush flushes all the flows of a specified table
// conntrack -F [table]            Flush table
// The flush operation applies to all the family types
func ConntrackTableFlush(table ConntrackTableType) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

// ConntrackDeleteFilter deletes entries on the specified table on the base of the filter
// conntrack -D [table] parameters         Delete conntrack or expectation
func ConntrackDeleteFilter(table ConntrackTableType, family InetFamily, filter *ConntrackFilter) (uint, error) ***REMOVED***
	return 0, ErrNotImplemented
***REMOVED***

// ConntrackTableList returns the flow list of a table of a specific family using the netlink handle passed
// conntrack -L [table] [options]          List conntrack or expectation table
func (h *Handle) ConntrackTableList(table ConntrackTableType, family InetFamily) ([]*ConntrackFlow, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

// ConntrackTableFlush flushes all the flows of a specified table using the netlink handle passed
// conntrack -F [table]            Flush table
// The flush operation applies to all the family types
func (h *Handle) ConntrackTableFlush(table ConntrackTableType) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

// ConntrackDeleteFilter deletes entries on the specified table on the base of the filter using the netlink handle passed
// conntrack -D [table] parameters         Delete conntrack or expectation
func (h *Handle) ConntrackDeleteFilter(table ConntrackTableType, family InetFamily, filter *ConntrackFilter) (uint, error) ***REMOVED***
	return 0, ErrNotImplemented
***REMOVED***
