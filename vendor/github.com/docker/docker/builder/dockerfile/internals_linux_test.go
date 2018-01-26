package dockerfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/idtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChownFlagParsing(t *testing.T) ***REMOVED***
	testFiles := map[string]string***REMOVED***
		"passwd": `root:x:0:0::/bin:/bin/false
bin:x:1:1::/bin:/bin/false
wwwwww:x:21:33::/bin:/bin/false
unicorn:x:1001:1002::/bin:/bin/false
		`,
		"group": `root:x:0:
bin:x:1:
wwwwww:x:33:
unicorn:x:1002:
somegrp:x:5555:
othergrp:x:6666:
		`,
	***REMOVED***
	// test mappings for validating use of maps
	idMaps := []idtools.IDMap***REMOVED***
		***REMOVED***
			ContainerID: 0,
			HostID:      100000,
			Size:        65536,
		***REMOVED***,
	***REMOVED***
	remapped := idtools.NewIDMappingsFromMaps(idMaps, idMaps)
	unmapped := &idtools.IDMappings***REMOVED******REMOVED***

	contextDir, cleanup := createTestTempDir(t, "", "builder-chown-parse-test")
	defer cleanup()

	if err := os.Mkdir(filepath.Join(contextDir, "etc"), 0755); err != nil ***REMOVED***
		t.Fatalf("error creating test directory: %v", err)
	***REMOVED***

	for filename, content := range testFiles ***REMOVED***
		createTestTempFile(t, filepath.Join(contextDir, "etc"), filename, content, 0644)
	***REMOVED***

	// positive tests
	for _, testcase := range []struct ***REMOVED***
		name      string
		chownStr  string
		idMapping *idtools.IDMappings
		expected  idtools.IDPair
	***REMOVED******REMOVED***
		***REMOVED***
			name:      "UIDNoMap",
			chownStr:  "1",
			idMapping: unmapped,
			expected:  idtools.IDPair***REMOVED***UID: 1, GID: 1***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:      "UIDGIDNoMap",
			chownStr:  "0:1",
			idMapping: unmapped,
			expected:  idtools.IDPair***REMOVED***UID: 0, GID: 1***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:      "UIDWithMap",
			chownStr:  "0",
			idMapping: remapped,
			expected:  idtools.IDPair***REMOVED***UID: 100000, GID: 100000***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:      "UIDGIDWithMap",
			chownStr:  "1:33",
			idMapping: remapped,
			expected:  idtools.IDPair***REMOVED***UID: 100001, GID: 100033***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:      "UserNoMap",
			chownStr:  "bin:5555",
			idMapping: unmapped,
			expected:  idtools.IDPair***REMOVED***UID: 1, GID: 5555***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:      "GroupWithMap",
			chownStr:  "0:unicorn",
			idMapping: remapped,
			expected:  idtools.IDPair***REMOVED***UID: 100000, GID: 101002***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:      "UserOnlyWithMap",
			chownStr:  "unicorn",
			idMapping: remapped,
			expected:  idtools.IDPair***REMOVED***UID: 101001, GID: 101002***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		t.Run(testcase.name, func(t *testing.T) ***REMOVED***
			idPair, err := parseChownFlag(testcase.chownStr, contextDir, testcase.idMapping)
			require.NoError(t, err, "Failed to parse chown flag: %q", testcase.chownStr)
			assert.Equal(t, testcase.expected, idPair, "chown flag mapping failure")
		***REMOVED***)
	***REMOVED***

	// error tests
	for _, testcase := range []struct ***REMOVED***
		name      string
		chownStr  string
		idMapping *idtools.IDMappings
		descr     string
	***REMOVED******REMOVED***
		***REMOVED***
			name:      "BadChownFlagFormat",
			chownStr:  "bob:1:555",
			idMapping: unmapped,
			descr:     "invalid chown string format: bob:1:555",
		***REMOVED***,
		***REMOVED***
			name:      "UserNoExist",
			chownStr:  "bob",
			idMapping: unmapped,
			descr:     "can't find uid for user bob: no such user: bob",
		***REMOVED***,
		***REMOVED***
			name:      "GroupNoExist",
			chownStr:  "root:bob",
			idMapping: unmapped,
			descr:     "can't find gid for group bob: no such group: bob",
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		t.Run(testcase.name, func(t *testing.T) ***REMOVED***
			_, err := parseChownFlag(testcase.chownStr, contextDir, testcase.idMapping)
			assert.EqualError(t, err, testcase.descr, "Expected error string doesn't match")
		***REMOVED***)
	***REMOVED***
***REMOVED***
