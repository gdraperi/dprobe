package discovery

import "net"

// NewEntry creates a new entry.
func NewEntry(url string) (*Entry, error) ***REMOVED***
	host, port, err := net.SplitHostPort(url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Entry***REMOVED***host, port***REMOVED***, nil
***REMOVED***

// An Entry represents a host.
type Entry struct ***REMOVED***
	Host string
	Port string
***REMOVED***

// Equals returns true if cmp contains the same data.
func (e *Entry) Equals(cmp *Entry) bool ***REMOVED***
	return e.Host == cmp.Host && e.Port == cmp.Port
***REMOVED***

// String returns the string form of an entry.
func (e *Entry) String() string ***REMOVED***
	return net.JoinHostPort(e.Host, e.Port)
***REMOVED***

// Entries is a list of *Entry with some helpers.
type Entries []*Entry

// Equals returns true if cmp contains the same data.
func (e Entries) Equals(cmp Entries) bool ***REMOVED***
	// Check if the file has really changed.
	if len(e) != len(cmp) ***REMOVED***
		return false
	***REMOVED***
	for i := range e ***REMOVED***
		if !e[i].Equals(cmp[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Contains returns true if the Entries contain a given Entry.
func (e Entries) Contains(entry *Entry) bool ***REMOVED***
	for _, curr := range e ***REMOVED***
		if curr.Equals(entry) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Diff compares two entries and returns the added and removed entries.
func (e Entries) Diff(cmp Entries) (Entries, Entries) ***REMOVED***
	added := Entries***REMOVED******REMOVED***
	for _, entry := range cmp ***REMOVED***
		if !e.Contains(entry) ***REMOVED***
			added = append(added, entry)
		***REMOVED***
	***REMOVED***

	removed := Entries***REMOVED******REMOVED***
	for _, entry := range e ***REMOVED***
		if !cmp.Contains(entry) ***REMOVED***
			removed = append(removed, entry)
		***REMOVED***
	***REMOVED***

	return added, removed
***REMOVED***

// CreateEntries returns an array of entries based on the given addresses.
func CreateEntries(addrs []string) (Entries, error) ***REMOVED***
	entries := Entries***REMOVED******REMOVED***
	if addrs == nil ***REMOVED***
		return entries, nil
	***REMOVED***

	for _, addr := range addrs ***REMOVED***
		if len(addr) == 0 ***REMOVED***
			continue
		***REMOVED***
		entry, err := NewEntry(addr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		entries = append(entries, entry)
	***REMOVED***
	return entries, nil
***REMOVED***
