package urlutil

import "testing"

var (
	gitUrls = []string***REMOVED***
		"git://github.com/docker/docker",
		"git@github.com:docker/docker.git",
		"git@bitbucket.org:atlassianlabs/atlassian-docker.git",
		"https://github.com/docker/docker.git",
		"http://github.com/docker/docker.git",
		"http://github.com/docker/docker.git#branch",
		"http://github.com/docker/docker.git#:dir",
	***REMOVED***
	incompleteGitUrls = []string***REMOVED***
		"github.com/docker/docker",
	***REMOVED***
	invalidGitUrls = []string***REMOVED***
		"http://github.com/docker/docker.git:#branch",
	***REMOVED***
	transportUrls = []string***REMOVED***
		"tcp://example.com",
		"tcp+tls://example.com",
		"udp://example.com",
		"unix:///example",
		"unixgram:///example",
	***REMOVED***
)

func TestIsGIT(t *testing.T) ***REMOVED***
	for _, url := range gitUrls ***REMOVED***
		if !IsGitURL(url) ***REMOVED***
			t.Fatalf("%q should be detected as valid Git url", url)
		***REMOVED***
	***REMOVED***

	for _, url := range incompleteGitUrls ***REMOVED***
		if !IsGitURL(url) ***REMOVED***
			t.Fatalf("%q should be detected as valid Git url", url)
		***REMOVED***
	***REMOVED***

	for _, url := range invalidGitUrls ***REMOVED***
		if IsGitURL(url) ***REMOVED***
			t.Fatalf("%q should not be detected as valid Git prefix", url)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsTransport(t *testing.T) ***REMOVED***
	for _, url := range transportUrls ***REMOVED***
		if !IsTransportURL(url) ***REMOVED***
			t.Fatalf("%q should be detected as valid Transport url", url)
		***REMOVED***
	***REMOVED***
***REMOVED***
