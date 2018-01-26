package discovery

import (
	"testing"

	"github.com/go-check/check"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) ***REMOVED*** check.TestingT(t) ***REMOVED***

type DiscoverySuite struct***REMOVED******REMOVED***

var _ = check.Suite(&DiscoverySuite***REMOVED******REMOVED***)

func (s *DiscoverySuite) TestNewEntry(c *check.C) ***REMOVED***
	entry, err := NewEntry("127.0.0.1:2375")
	c.Assert(err, check.IsNil)
	c.Assert(entry.Equals(&Entry***REMOVED***Host: "127.0.0.1", Port: "2375"***REMOVED***), check.Equals, true)
	c.Assert(entry.String(), check.Equals, "127.0.0.1:2375")

	entry, err = NewEntry("[2001:db8:0:f101::2]:2375")
	c.Assert(err, check.IsNil)
	c.Assert(entry.Equals(&Entry***REMOVED***Host: "2001:db8:0:f101::2", Port: "2375"***REMOVED***), check.Equals, true)
	c.Assert(entry.String(), check.Equals, "[2001:db8:0:f101::2]:2375")

	_, err = NewEntry("127.0.0.1")
	c.Assert(err, check.NotNil)
***REMOVED***

func (s *DiscoverySuite) TestParse(c *check.C) ***REMOVED***
	scheme, uri := parse("127.0.0.1:2375")
	c.Assert(scheme, check.Equals, "nodes")
	c.Assert(uri, check.Equals, "127.0.0.1:2375")

	scheme, uri = parse("localhost:2375")
	c.Assert(scheme, check.Equals, "nodes")
	c.Assert(uri, check.Equals, "localhost:2375")

	scheme, uri = parse("scheme://127.0.0.1:2375")
	c.Assert(scheme, check.Equals, "scheme")
	c.Assert(uri, check.Equals, "127.0.0.1:2375")

	scheme, uri = parse("scheme://localhost:2375")
	c.Assert(scheme, check.Equals, "scheme")
	c.Assert(uri, check.Equals, "localhost:2375")

	scheme, uri = parse("")
	c.Assert(scheme, check.Equals, "nodes")
	c.Assert(uri, check.Equals, "")
***REMOVED***

func (s *DiscoverySuite) TestCreateEntries(c *check.C) ***REMOVED***
	entries, err := CreateEntries(nil)
	c.Assert(entries, check.DeepEquals, Entries***REMOVED******REMOVED***)
	c.Assert(err, check.IsNil)

	entries, err = CreateEntries([]string***REMOVED***"127.0.0.1:2375", "127.0.0.2:2375", "[2001:db8:0:f101::2]:2375", ""***REMOVED***)
	c.Assert(err, check.IsNil)
	expected := Entries***REMOVED***
		&Entry***REMOVED***Host: "127.0.0.1", Port: "2375"***REMOVED***,
		&Entry***REMOVED***Host: "127.0.0.2", Port: "2375"***REMOVED***,
		&Entry***REMOVED***Host: "2001:db8:0:f101::2", Port: "2375"***REMOVED***,
	***REMOVED***
	c.Assert(entries.Equals(expected), check.Equals, true)

	_, err = CreateEntries([]string***REMOVED***"127.0.0.1", "127.0.0.2"***REMOVED***)
	c.Assert(err, check.NotNil)
***REMOVED***

func (s *DiscoverySuite) TestContainsEntry(c *check.C) ***REMOVED***
	entries, err := CreateEntries([]string***REMOVED***"127.0.0.1:2375", "127.0.0.2:2375", ""***REMOVED***)
	c.Assert(err, check.IsNil)
	c.Assert(entries.Contains(&Entry***REMOVED***Host: "127.0.0.1", Port: "2375"***REMOVED***), check.Equals, true)
	c.Assert(entries.Contains(&Entry***REMOVED***Host: "127.0.0.3", Port: "2375"***REMOVED***), check.Equals, false)
***REMOVED***

func (s *DiscoverySuite) TestEntriesEquality(c *check.C) ***REMOVED***
	entries := Entries***REMOVED***
		&Entry***REMOVED***Host: "127.0.0.1", Port: "2375"***REMOVED***,
		&Entry***REMOVED***Host: "127.0.0.2", Port: "2375"***REMOVED***,
	***REMOVED***

	// Same
	c.Assert(entries.Equals(Entries***REMOVED***
		&Entry***REMOVED***Host: "127.0.0.1", Port: "2375"***REMOVED***,
		&Entry***REMOVED***Host: "127.0.0.2", Port: "2375"***REMOVED***,
	***REMOVED***), check.
		Equals, true)

	// Different size
	c.Assert(entries.Equals(Entries***REMOVED***
		&Entry***REMOVED***Host: "127.0.0.1", Port: "2375"***REMOVED***,
		&Entry***REMOVED***Host: "127.0.0.2", Port: "2375"***REMOVED***,
		&Entry***REMOVED***Host: "127.0.0.3", Port: "2375"***REMOVED***,
	***REMOVED***), check.
		Equals, false)

	// Different content
	c.Assert(entries.Equals(Entries***REMOVED***
		&Entry***REMOVED***Host: "127.0.0.1", Port: "2375"***REMOVED***,
		&Entry***REMOVED***Host: "127.0.0.42", Port: "2375"***REMOVED***,
	***REMOVED***), check.
		Equals, false)

***REMOVED***

func (s *DiscoverySuite) TestEntriesDiff(c *check.C) ***REMOVED***
	entry1 := &Entry***REMOVED***Host: "1.1.1.1", Port: "1111"***REMOVED***
	entry2 := &Entry***REMOVED***Host: "2.2.2.2", Port: "2222"***REMOVED***
	entry3 := &Entry***REMOVED***Host: "3.3.3.3", Port: "3333"***REMOVED***
	entries := Entries***REMOVED***entry1, entry2***REMOVED***

	// No diff
	added, removed := entries.Diff(Entries***REMOVED***entry2, entry1***REMOVED***)
	c.Assert(added, check.HasLen, 0)
	c.Assert(removed, check.HasLen, 0)

	// Add
	added, removed = entries.Diff(Entries***REMOVED***entry2, entry3, entry1***REMOVED***)
	c.Assert(added, check.HasLen, 1)
	c.Assert(added.Contains(entry3), check.Equals, true)
	c.Assert(removed, check.HasLen, 0)

	// Remove
	added, removed = entries.Diff(Entries***REMOVED***entry2***REMOVED***)
	c.Assert(added, check.HasLen, 0)
	c.Assert(removed, check.HasLen, 1)
	c.Assert(removed.Contains(entry1), check.Equals, true)

	// Add and remove
	added, removed = entries.Diff(Entries***REMOVED***entry1, entry3***REMOVED***)
	c.Assert(added, check.HasLen, 1)
	c.Assert(added.Contains(entry3), check.Equals, true)
	c.Assert(removed, check.HasLen, 1)
	c.Assert(removed.Contains(entry2), check.Equals, true)
***REMOVED***
