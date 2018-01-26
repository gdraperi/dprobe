package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/pkg/discovery"

	"github.com/go-check/check"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) ***REMOVED*** check.TestingT(t) ***REMOVED***

type DiscoverySuite struct***REMOVED******REMOVED***

var _ = check.Suite(&DiscoverySuite***REMOVED******REMOVED***)

func (s *DiscoverySuite) TestInitialize(c *check.C) ***REMOVED***
	d := &Discovery***REMOVED******REMOVED***
	d.Initialize("/path/to/file", 1000, 0, nil)
	c.Assert(d.path, check.Equals, "/path/to/file")
***REMOVED***

func (s *DiscoverySuite) TestNew(c *check.C) ***REMOVED***
	d, err := discovery.New("file:///path/to/file", 0, 0, nil)
	c.Assert(err, check.IsNil)
	c.Assert(d.(*Discovery).path, check.Equals, "/path/to/file")
***REMOVED***

func (s *DiscoverySuite) TestContent(c *check.C) ***REMOVED***
	data := `
1.1.1.[1:2]:1111
2.2.2.[2:4]:2222
`
	ips := parseFileContent([]byte(data))
	c.Assert(ips, check.HasLen, 5)
	c.Assert(ips[0], check.Equals, "1.1.1.1:1111")
	c.Assert(ips[1], check.Equals, "1.1.1.2:1111")
	c.Assert(ips[2], check.Equals, "2.2.2.2:2222")
	c.Assert(ips[3], check.Equals, "2.2.2.3:2222")
	c.Assert(ips[4], check.Equals, "2.2.2.4:2222")
***REMOVED***

func (s *DiscoverySuite) TestRegister(c *check.C) ***REMOVED***
	discovery := &Discovery***REMOVED***path: "/path/to/file"***REMOVED***
	c.Assert(discovery.Register("0.0.0.0"), check.NotNil)
***REMOVED***

func (s *DiscoverySuite) TestParsingContentsWithComments(c *check.C) ***REMOVED***
	data := `
### test ###
1.1.1.1:1111 # inline comment
# 2.2.2.2:2222
      ### empty line with comment
    3.3.3.3:3333
### test ###
`
	ips := parseFileContent([]byte(data))
	c.Assert(ips, check.HasLen, 2)
	c.Assert("1.1.1.1:1111", check.Equals, ips[0])
	c.Assert("3.3.3.3:3333", check.Equals, ips[1])
***REMOVED***

func (s *DiscoverySuite) TestWatch(c *check.C) ***REMOVED***
	data := `
1.1.1.1:1111
2.2.2.2:2222
`
	expected := discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "1.1.1.1", Port: "1111"***REMOVED***,
		&discovery.Entry***REMOVED***Host: "2.2.2.2", Port: "2222"***REMOVED***,
	***REMOVED***

	// Create a temporary file and remove it.
	tmp, err := ioutil.TempFile(os.TempDir(), "discovery-file-test")
	c.Assert(err, check.IsNil)
	c.Assert(tmp.Close(), check.IsNil)
	c.Assert(os.Remove(tmp.Name()), check.IsNil)

	// Set up file discovery.
	d := &Discovery***REMOVED******REMOVED***
	d.Initialize(tmp.Name(), 1000, 0, nil)
	stopCh := make(chan struct***REMOVED******REMOVED***)
	ch, errCh := d.Watch(stopCh)

	// Make sure it fires errors since the file doesn't exist.
	c.Assert(<-errCh, check.NotNil)
	// We have to drain the error channel otherwise Watch will get stuck.
	go func() ***REMOVED***
		for range errCh ***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Write the file and make sure we get the expected value back.
	c.Assert(ioutil.WriteFile(tmp.Name(), []byte(data), 0600), check.IsNil)
	c.Assert(<-ch, check.DeepEquals, expected)

	// Add a new entry and look it up.
	expected = append(expected, &discovery.Entry***REMOVED***Host: "3.3.3.3", Port: "3333"***REMOVED***)
	f, err := os.OpenFile(tmp.Name(), os.O_APPEND|os.O_WRONLY, 0600)
	c.Assert(err, check.IsNil)
	c.Assert(f, check.NotNil)
	_, err = f.WriteString("\n3.3.3.3:3333\n")
	c.Assert(err, check.IsNil)
	f.Close()
	c.Assert(<-ch, check.DeepEquals, expected)

	// Stop and make sure it closes all channels.
	close(stopCh)
	c.Assert(<-ch, check.IsNil)
	c.Assert(<-errCh, check.IsNil)
***REMOVED***
