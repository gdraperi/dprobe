package loggerutils

import (
	"testing"

	"github.com/docker/docker/daemon/logger"
)

func TestParseLogTagDefaultTag(t *testing.T) ***REMOVED***
	info := buildContext(map[string]string***REMOVED******REMOVED***)
	tag, e := ParseLogTag(info, "***REMOVED******REMOVED***.ID***REMOVED******REMOVED***")
	assertTag(t, e, tag, info.ID())
***REMOVED***

func TestParseLogTag(t *testing.T) ***REMOVED***
	info := buildContext(map[string]string***REMOVED***"tag": "***REMOVED******REMOVED***.ImageName***REMOVED******REMOVED***/***REMOVED******REMOVED***.Name***REMOVED******REMOVED***/***REMOVED******REMOVED***.ID***REMOVED******REMOVED***"***REMOVED***)
	tag, e := ParseLogTag(info, "***REMOVED******REMOVED***.ID***REMOVED******REMOVED***")
	assertTag(t, e, tag, "test-image/test-container/container-ab")
***REMOVED***

func TestParseLogTagEmptyTag(t *testing.T) ***REMOVED***
	info := buildContext(map[string]string***REMOVED******REMOVED***)
	tag, e := ParseLogTag(info, "***REMOVED******REMOVED***.DaemonName***REMOVED******REMOVED***/***REMOVED******REMOVED***.ID***REMOVED******REMOVED***")
	assertTag(t, e, tag, "test-dockerd/container-ab")
***REMOVED***

// Helpers

func buildContext(cfg map[string]string) logger.Info ***REMOVED***
	return logger.Info***REMOVED***
		ContainerID:        "container-abcdefghijklmnopqrstuvwxyz01234567890",
		ContainerName:      "/test-container",
		ContainerImageID:   "image-abcdefghijklmnopqrstuvwxyz01234567890",
		ContainerImageName: "test-image",
		Config:             cfg,
		DaemonName:         "test-dockerd",
	***REMOVED***
***REMOVED***

func assertTag(t *testing.T, e error, tag string, expected string) ***REMOVED***
	if e != nil ***REMOVED***
		t.Fatalf("Error generating tag: %q", e)
	***REMOVED***
	if tag != expected ***REMOVED***
		t.Fatalf("Wrong tag: %q, should be %q", tag, expected)
	***REMOVED***
***REMOVED***
