package memory

import (
	"testing"

	"github.com/docker/docker/pkg/discovery"
	"github.com/go-check/check"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) ***REMOVED*** check.TestingT(t) ***REMOVED***

type discoverySuite struct***REMOVED******REMOVED***

var _ = check.Suite(&discoverySuite***REMOVED******REMOVED***)

func (s *discoverySuite) TestWatch(c *check.C) ***REMOVED***
	d := &Discovery***REMOVED******REMOVED***
	d.Initialize("foo", 1000, 0, nil)
	stopCh := make(chan struct***REMOVED******REMOVED***)
	ch, errCh := d.Watch(stopCh)

	// We have to drain the error channel otherwise Watch will get stuck.
	go func() ***REMOVED***
		for range errCh ***REMOVED***
		***REMOVED***
	***REMOVED***()

	expected := discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "1.1.1.1", Port: "1111"***REMOVED***,
	***REMOVED***

	c.Assert(d.Register("1.1.1.1:1111"), check.IsNil)
	c.Assert(<-ch, check.DeepEquals, expected)

	expected = discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "1.1.1.1", Port: "1111"***REMOVED***,
		&discovery.Entry***REMOVED***Host: "2.2.2.2", Port: "2222"***REMOVED***,
	***REMOVED***

	c.Assert(d.Register("2.2.2.2:2222"), check.IsNil)
	c.Assert(<-ch, check.DeepEquals, expected)

	// Stop and make sure it closes all channels.
	close(stopCh)
	c.Assert(<-ch, check.IsNil)
	c.Assert(<-errCh, check.IsNil)
***REMOVED***
