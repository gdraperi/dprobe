package distribution

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/progress"
	"github.com/opencontainers/go-digest"
)

func TestGetRepositoryMountCandidates(t *testing.T) ***REMOVED***
	for _, tc := range []struct ***REMOVED***
		name          string
		hmacKey       string
		targetRepo    string
		maxCandidates int
		metadata      []metadata.V2Metadata
		candidates    []metadata.V2Metadata
	***REMOVED******REMOVED***
		***REMOVED***
			name:          "empty metadata",
			targetRepo:    "busybox",
			maxCandidates: -1,
			metadata:      []metadata.V2Metadata***REMOVED******REMOVED***,
			candidates:    []metadata.V2Metadata***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:          "one item not matching",
			targetRepo:    "busybox",
			maxCandidates: -1,
			metadata:      []metadata.V2Metadata***REMOVED***taggedMetadata("key", "dgst", "127.0.0.1/repo")***REMOVED***,
			candidates:    []metadata.V2Metadata***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:          "one item matching",
			targetRepo:    "busybox",
			maxCandidates: -1,
			metadata:      []metadata.V2Metadata***REMOVED***taggedMetadata("hash", "1", "docker.io/library/hello-world")***REMOVED***,
			candidates:    []metadata.V2Metadata***REMOVED***taggedMetadata("hash", "1", "docker.io/library/hello-world")***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:          "allow missing SourceRepository",
			targetRepo:    "busybox",
			maxCandidates: -1,
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("1")***REMOVED***,
				***REMOVED***Digest: digest.Digest("3")***REMOVED***,
				***REMOVED***Digest: digest.Digest("2")***REMOVED***,
			***REMOVED***,
			candidates: []metadata.V2Metadata***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:          "handle docker.io",
			targetRepo:    "user/app",
			maxCandidates: -1,
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("1"), SourceRepository: "docker.io/user/foo"***REMOVED***,
				***REMOVED***Digest: digest.Digest("3"), SourceRepository: "docker.io/user/bar"***REMOVED***,
				***REMOVED***Digest: digest.Digest("2"), SourceRepository: "docker.io/library/app"***REMOVED***,
			***REMOVED***,
			candidates: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("3"), SourceRepository: "docker.io/user/bar"***REMOVED***,
				***REMOVED***Digest: digest.Digest("1"), SourceRepository: "docker.io/user/foo"***REMOVED***,
				***REMOVED***Digest: digest.Digest("2"), SourceRepository: "docker.io/library/app"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:          "sort more items",
			hmacKey:       "abcd",
			targetRepo:    "127.0.0.1/foo/bar",
			maxCandidates: -1,
			metadata: []metadata.V2Metadata***REMOVED***
				taggedMetadata("hash", "1", "docker.io/library/hello-world"),
				taggedMetadata("efgh", "2", "127.0.0.1/hello-world"),
				taggedMetadata("abcd", "3", "docker.io/library/busybox"),
				taggedMetadata("hash", "4", "docker.io/library/busybox"),
				taggedMetadata("hash", "5", "127.0.0.1/foo"),
				taggedMetadata("hash", "6", "127.0.0.1/bar"),
				taggedMetadata("efgh", "7", "127.0.0.1/foo/bar"),
				taggedMetadata("abcd", "8", "127.0.0.1/xyz"),
				taggedMetadata("hash", "9", "127.0.0.1/foo/app"),
			***REMOVED***,
			candidates: []metadata.V2Metadata***REMOVED***
				// first by matching hash
				taggedMetadata("abcd", "8", "127.0.0.1/xyz"),
				// then by longest matching prefix
				taggedMetadata("hash", "9", "127.0.0.1/foo/app"),
				taggedMetadata("hash", "5", "127.0.0.1/foo"),
				// sort the rest of the matching items in reversed order
				taggedMetadata("hash", "6", "127.0.0.1/bar"),
				taggedMetadata("efgh", "2", "127.0.0.1/hello-world"),
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:          "limit max candidates",
			hmacKey:       "abcd",
			targetRepo:    "user/app",
			maxCandidates: 3,
			metadata: []metadata.V2Metadata***REMOVED***
				taggedMetadata("abcd", "1", "docker.io/user/app1"),
				taggedMetadata("abcd", "2", "docker.io/user/app/base"),
				taggedMetadata("hash", "3", "docker.io/user/app"),
				taggedMetadata("abcd", "4", "127.0.0.1/user/app"),
				taggedMetadata("hash", "5", "docker.io/user/foo"),
				taggedMetadata("hash", "6", "docker.io/app/bar"),
			***REMOVED***,
			candidates: []metadata.V2Metadata***REMOVED***
				// first by matching hash
				taggedMetadata("abcd", "2", "docker.io/user/app/base"),
				taggedMetadata("abcd", "1", "docker.io/user/app1"),
				// then by longest matching prefix
				// "docker.io/usr/app" is excluded since candidates must
				// be from a different repository
				taggedMetadata("hash", "5", "docker.io/user/foo"),
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		repoInfo, err := reference.ParseNormalizedNamed(tc.targetRepo)
		if err != nil ***REMOVED***
			t.Fatalf("[%s] failed to parse reference name: %v", tc.name, err)
		***REMOVED***
		candidates := getRepositoryMountCandidates(repoInfo, []byte(tc.hmacKey), tc.maxCandidates, tc.metadata)
		if len(candidates) != len(tc.candidates) ***REMOVED***
			t.Errorf("[%s] got unexpected number of candidates: %d != %d", tc.name, len(candidates), len(tc.candidates))
		***REMOVED***
		for i := 0; i < len(candidates) && i < len(tc.candidates); i++ ***REMOVED***
			if !reflect.DeepEqual(candidates[i], tc.candidates[i]) ***REMOVED***
				t.Errorf("[%s] candidate %d does not match expected: %#+v != %#+v", tc.name, i, candidates[i], tc.candidates[i])
			***REMOVED***
		***REMOVED***
		for i := len(candidates); i < len(tc.candidates); i++ ***REMOVED***
			t.Errorf("[%s] missing expected candidate at position %d (%#+v)", tc.name, i, tc.candidates[i])
		***REMOVED***
		for i := len(tc.candidates); i < len(candidates); i++ ***REMOVED***
			t.Errorf("[%s] got unexpected candidate at position %d (%#+v)", tc.name, i, candidates[i])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestLayerAlreadyExists(t *testing.T) ***REMOVED***
	for _, tc := range []struct ***REMOVED***
		name                   string
		metadata               []metadata.V2Metadata
		targetRepo             string
		hmacKey                string
		maxExistenceChecks     int
		checkOtherRepositories bool
		remoteBlobs            map[digest.Digest]distribution.Descriptor
		remoteErrors           map[digest.Digest]error
		expectedDescriptor     distribution.Descriptor
		expectedExists         bool
		expectedError          error
		expectedRequests       []string
		expectedAdditions      []metadata.V2Metadata
		expectedRemovals       []metadata.V2Metadata
	***REMOVED******REMOVED***
		***REMOVED***
			name:                   "empty metadata",
			targetRepo:             "busybox",
			maxExistenceChecks:     3,
			checkOtherRepositories: true,
		***REMOVED***,
		***REMOVED***
			name:               "single not existent metadata",
			targetRepo:         "busybox",
			metadata:           []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("pear"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
			maxExistenceChecks: 3,
			expectedRequests:   []string***REMOVED***"pear"***REMOVED***,
			expectedRemovals:   []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("pear"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:               "access denied",
			targetRepo:         "busybox",
			maxExistenceChecks: 1,
			metadata:           []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
			remoteErrors:       map[digest.Digest]error***REMOVED***digest.Digest("apple"): distribution.ErrAccessDenied***REMOVED***,
			expectedError:      nil,
			expectedRequests:   []string***REMOVED***"apple"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:               "not matching repositories",
			targetRepo:         "busybox",
			maxExistenceChecks: 3,
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/library/hello-world"***REMOVED***,
				***REMOVED***Digest: digest.Digest("orange"), SourceRepository: "docker.io/library/busybox/subapp"***REMOVED***,
				***REMOVED***Digest: digest.Digest("pear"), SourceRepository: "docker.io/busybox"***REMOVED***,
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "busybox"***REMOVED***,
				***REMOVED***Digest: digest.Digest("banana"), SourceRepository: "127.0.0.1/busybox"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:                   "check other repositories",
			targetRepo:             "busybox",
			maxExistenceChecks:     10,
			checkOtherRepositories: true,
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/library/hello-world"***REMOVED***,
				***REMOVED***Digest: digest.Digest("orange"), SourceRepository: "docker.io/busybox/subapp"***REMOVED***,
				***REMOVED***Digest: digest.Digest("pear"), SourceRepository: "docker.io/busybox"***REMOVED***,
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "docker.io/library/busybox"***REMOVED***,
				***REMOVED***Digest: digest.Digest("banana"), SourceRepository: "127.0.0.1/busybox"***REMOVED***,
			***REMOVED***,
			expectedRequests: []string***REMOVED***"plum", "apple", "pear", "orange", "banana"***REMOVED***,
			expectedRemovals: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "docker.io/library/busybox"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:               "find existing blob",
			targetRepo:         "busybox",
			metadata:           []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
			maxExistenceChecks: 3,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("apple"): ***REMOVED***Digest: digest.Digest("apple")***REMOVED******REMOVED***,
			expectedDescriptor: distribution.Descriptor***REMOVED***Digest: digest.Digest("apple"), MediaType: schema2.MediaTypeLayer***REMOVED***,
			expectedExists:     true,
			expectedRequests:   []string***REMOVED***"apple"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:               "find existing blob with different hmac",
			targetRepo:         "busybox",
			metadata:           []metadata.V2Metadata***REMOVED******REMOVED***SourceRepository: "docker.io/library/busybox", Digest: digest.Digest("apple"), HMAC: "dummyhmac"***REMOVED******REMOVED***,
			maxExistenceChecks: 3,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("apple"): ***REMOVED***Digest: digest.Digest("apple")***REMOVED******REMOVED***,
			expectedDescriptor: distribution.Descriptor***REMOVED***Digest: digest.Digest("apple"), MediaType: schema2.MediaTypeLayer***REMOVED***,
			expectedExists:     true,
			expectedRequests:   []string***REMOVED***"apple"***REMOVED***,
			expectedAdditions:  []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:               "overwrite media types",
			targetRepo:         "busybox",
			metadata:           []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
			hmacKey:            "key",
			maxExistenceChecks: 3,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("apple"): ***REMOVED***Digest: digest.Digest("apple"), MediaType: "custom-media-type"***REMOVED******REMOVED***,
			expectedDescriptor: distribution.Descriptor***REMOVED***Digest: digest.Digest("apple"), MediaType: schema2.MediaTypeLayer***REMOVED***,
			expectedExists:     true,
			expectedRequests:   []string***REMOVED***"apple"***REMOVED***,
			expectedAdditions:  []metadata.V2Metadata***REMOVED***taggedMetadata("key", "apple", "docker.io/library/busybox")***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:       "find existing blob among many",
			targetRepo: "127.0.0.1/myapp",
			hmacKey:    "key",
			metadata: []metadata.V2Metadata***REMOVED***
				taggedMetadata("someotherkey", "pear", "127.0.0.1/myapp"),
				taggedMetadata("key", "apple", "127.0.0.1/myapp"),
				taggedMetadata("", "plum", "127.0.0.1/myapp"),
			***REMOVED***,
			maxExistenceChecks: 3,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("pear"): ***REMOVED***Digest: digest.Digest("pear")***REMOVED******REMOVED***,
			expectedDescriptor: distribution.Descriptor***REMOVED***Digest: digest.Digest("pear"), MediaType: schema2.MediaTypeLayer***REMOVED***,
			expectedExists:     true,
			expectedRequests:   []string***REMOVED***"apple", "plum", "pear"***REMOVED***,
			expectedAdditions:  []metadata.V2Metadata***REMOVED***taggedMetadata("key", "pear", "127.0.0.1/myapp")***REMOVED***,
			expectedRemovals: []metadata.V2Metadata***REMOVED***
				taggedMetadata("key", "apple", "127.0.0.1/myapp"),
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "127.0.0.1/myapp"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:       "reach maximum existence checks",
			targetRepo: "user/app",
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("pear"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("banana"), SourceRepository: "docker.io/user/app"***REMOVED***,
			***REMOVED***,
			maxExistenceChecks: 3,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("pear"): ***REMOVED***Digest: digest.Digest("pear")***REMOVED******REMOVED***,
			expectedExists:     false,
			expectedRequests:   []string***REMOVED***"banana", "plum", "apple"***REMOVED***,
			expectedRemovals: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("banana"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/user/app"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:       "zero allowed existence checks",
			targetRepo: "user/app",
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("pear"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("apple"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "docker.io/user/app"***REMOVED***,
				***REMOVED***Digest: digest.Digest("banana"), SourceRepository: "docker.io/user/app"***REMOVED***,
			***REMOVED***,
			maxExistenceChecks: 0,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("pear"): ***REMOVED***Digest: digest.Digest("pear")***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:       "stat single digest just once",
			targetRepo: "busybox",
			metadata: []metadata.V2Metadata***REMOVED***
				taggedMetadata("key1", "pear", "docker.io/library/busybox"),
				taggedMetadata("key2", "apple", "docker.io/library/busybox"),
				taggedMetadata("key3", "apple", "docker.io/library/busybox"),
			***REMOVED***,
			maxExistenceChecks: 3,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("pear"): ***REMOVED***Digest: digest.Digest("pear")***REMOVED******REMOVED***,
			expectedDescriptor: distribution.Descriptor***REMOVED***Digest: digest.Digest("pear"), MediaType: schema2.MediaTypeLayer***REMOVED***,
			expectedExists:     true,
			expectedRequests:   []string***REMOVED***"apple", "pear"***REMOVED***,
			expectedAdditions:  []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("pear"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
			expectedRemovals:   []metadata.V2Metadata***REMOVED***taggedMetadata("key3", "apple", "docker.io/library/busybox")***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:       "don't stop on first error",
			targetRepo: "user/app",
			hmacKey:    "key",
			metadata: []metadata.V2Metadata***REMOVED***
				taggedMetadata("key", "banana", "docker.io/user/app"),
				taggedMetadata("key", "orange", "docker.io/user/app"),
				taggedMetadata("key", "plum", "docker.io/user/app"),
			***REMOVED***,
			maxExistenceChecks: 3,
			remoteErrors:       map[digest.Digest]error***REMOVED***"orange": distribution.ErrAccessDenied***REMOVED***,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("apple"): ***REMOVED******REMOVED******REMOVED***,
			expectedError:      nil,
			expectedRequests:   []string***REMOVED***"plum", "orange", "banana"***REMOVED***,
			expectedRemovals: []metadata.V2Metadata***REMOVED***
				taggedMetadata("key", "plum", "docker.io/user/app"),
				taggedMetadata("key", "banana", "docker.io/user/app"),
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:       "remove outdated metadata",
			targetRepo: "docker.io/user/app",
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("plum"), SourceRepository: "docker.io/library/busybox"***REMOVED***,
				***REMOVED***Digest: digest.Digest("orange"), SourceRepository: "docker.io/user/app"***REMOVED***,
			***REMOVED***,
			maxExistenceChecks: 3,
			remoteErrors:       map[digest.Digest]error***REMOVED***"orange": distribution.ErrBlobUnknown***REMOVED***,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("plum"): ***REMOVED******REMOVED******REMOVED***,
			expectedExists:     false,
			expectedRequests:   []string***REMOVED***"orange"***REMOVED***,
			expectedRemovals:   []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("orange"), SourceRepository: "docker.io/user/app"***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			name:       "missing SourceRepository",
			targetRepo: "busybox",
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("1")***REMOVED***,
				***REMOVED***Digest: digest.Digest("3")***REMOVED***,
				***REMOVED***Digest: digest.Digest("2")***REMOVED***,
			***REMOVED***,
			maxExistenceChecks: 3,
			expectedExists:     false,
			expectedRequests:   []string***REMOVED***"2", "3", "1"***REMOVED***,
		***REMOVED***,

		***REMOVED***
			name:       "with and without SourceRepository",
			targetRepo: "busybox",
			metadata: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("1")***REMOVED***,
				***REMOVED***Digest: digest.Digest("2"), SourceRepository: "docker.io/library/busybox"***REMOVED***,
				***REMOVED***Digest: digest.Digest("3")***REMOVED***,
			***REMOVED***,
			remoteBlobs:        map[digest.Digest]distribution.Descriptor***REMOVED***digest.Digest("1"): ***REMOVED***Digest: digest.Digest("1")***REMOVED******REMOVED***,
			maxExistenceChecks: 3,
			expectedDescriptor: distribution.Descriptor***REMOVED***Digest: digest.Digest("1"), MediaType: schema2.MediaTypeLayer***REMOVED***,
			expectedExists:     true,
			expectedRequests:   []string***REMOVED***"2", "3", "1"***REMOVED***,
			expectedAdditions:  []metadata.V2Metadata***REMOVED******REMOVED***Digest: digest.Digest("1"), SourceRepository: "docker.io/library/busybox"***REMOVED******REMOVED***,
			expectedRemovals: []metadata.V2Metadata***REMOVED***
				***REMOVED***Digest: digest.Digest("2"), SourceRepository: "docker.io/library/busybox"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		repoInfo, err := reference.ParseNormalizedNamed(tc.targetRepo)
		if err != nil ***REMOVED***
			t.Fatalf("[%s] failed to parse reference name: %v", tc.name, err)
		***REMOVED***
		repo := &mockRepo***REMOVED***
			t:        t,
			errors:   tc.remoteErrors,
			blobs:    tc.remoteBlobs,
			requests: []string***REMOVED******REMOVED***,
		***REMOVED***
		ctx := context.Background()
		ms := &mockV2MetadataService***REMOVED******REMOVED***
		pd := &v2PushDescriptor***REMOVED***
			hmacKey:  []byte(tc.hmacKey),
			repoInfo: repoInfo,
			layer: &storeLayer***REMOVED***
				Layer: layer.EmptyLayer,
			***REMOVED***,
			repo:              repo,
			v2MetadataService: ms,
			pushState:         &pushState***REMOVED***remoteLayers: make(map[layer.DiffID]distribution.Descriptor)***REMOVED***,
			checkedDigests:    make(map[digest.Digest]struct***REMOVED******REMOVED***),
		***REMOVED***

		desc, exists, err := pd.layerAlreadyExists(ctx, &progressSink***REMOVED***t***REMOVED***, layer.EmptyLayer.DiffID(), tc.checkOtherRepositories, tc.maxExistenceChecks, tc.metadata)

		if !reflect.DeepEqual(desc, tc.expectedDescriptor) ***REMOVED***
			t.Errorf("[%s] got unexpected descriptor: %#+v != %#+v", tc.name, desc, tc.expectedDescriptor)
		***REMOVED***
		if exists != tc.expectedExists ***REMOVED***
			t.Errorf("[%s] got unexpected exists: %t != %t", tc.name, exists, tc.expectedExists)
		***REMOVED***
		if !reflect.DeepEqual(err, tc.expectedError) ***REMOVED***
			t.Errorf("[%s] got unexpected error: %#+v != %#+v", tc.name, err, tc.expectedError)
		***REMOVED***

		if len(repo.requests) != len(tc.expectedRequests) ***REMOVED***
			t.Errorf("[%s] got unexpected number of requests: %d != %d", tc.name, len(repo.requests), len(tc.expectedRequests))
		***REMOVED***
		for i := 0; i < len(repo.requests) && i < len(tc.expectedRequests); i++ ***REMOVED***
			if repo.requests[i] != tc.expectedRequests[i] ***REMOVED***
				t.Errorf("[%s] request %d does not match expected: %q != %q", tc.name, i, repo.requests[i], tc.expectedRequests[i])
			***REMOVED***
		***REMOVED***
		for i := len(repo.requests); i < len(tc.expectedRequests); i++ ***REMOVED***
			t.Errorf("[%s] missing expected request at position %d (%q)", tc.name, i, tc.expectedRequests[i])
		***REMOVED***
		for i := len(tc.expectedRequests); i < len(repo.requests); i++ ***REMOVED***
			t.Errorf("[%s] got unexpected request at position %d (%q)", tc.name, i, repo.requests[i])
		***REMOVED***

		if len(ms.added) != len(tc.expectedAdditions) ***REMOVED***
			t.Errorf("[%s] got unexpected number of additions: %d != %d", tc.name, len(ms.added), len(tc.expectedAdditions))
		***REMOVED***
		for i := 0; i < len(ms.added) && i < len(tc.expectedAdditions); i++ ***REMOVED***
			if ms.added[i] != tc.expectedAdditions[i] ***REMOVED***
				t.Errorf("[%s] added metadata at %d does not match expected: %q != %q", tc.name, i, ms.added[i], tc.expectedAdditions[i])
			***REMOVED***
		***REMOVED***
		for i := len(ms.added); i < len(tc.expectedAdditions); i++ ***REMOVED***
			t.Errorf("[%s] missing expected addition at position %d (%q)", tc.name, i, tc.expectedAdditions[i])
		***REMOVED***
		for i := len(tc.expectedAdditions); i < len(ms.added); i++ ***REMOVED***
			t.Errorf("[%s] unexpected metadata addition at position %d (%q)", tc.name, i, ms.added[i])
		***REMOVED***

		if len(ms.removed) != len(tc.expectedRemovals) ***REMOVED***
			t.Errorf("[%s] got unexpected number of removals: %d != %d", tc.name, len(ms.removed), len(tc.expectedRemovals))
		***REMOVED***
		for i := 0; i < len(ms.removed) && i < len(tc.expectedRemovals); i++ ***REMOVED***
			if ms.removed[i] != tc.expectedRemovals[i] ***REMOVED***
				t.Errorf("[%s] removed metadata at %d does not match expected: %q != %q", tc.name, i, ms.removed[i], tc.expectedRemovals[i])
			***REMOVED***
		***REMOVED***
		for i := len(ms.removed); i < len(tc.expectedRemovals); i++ ***REMOVED***
			t.Errorf("[%s] missing expected removal at position %d (%q)", tc.name, i, tc.expectedRemovals[i])
		***REMOVED***
		for i := len(tc.expectedRemovals); i < len(ms.removed); i++ ***REMOVED***
			t.Errorf("[%s] removed unexpected metadata at position %d (%q)", tc.name, i, ms.removed[i])
		***REMOVED***
	***REMOVED***
***REMOVED***

func taggedMetadata(key string, dgst string, sourceRepo string) metadata.V2Metadata ***REMOVED***
	meta := metadata.V2Metadata***REMOVED***
		Digest:           digest.Digest(dgst),
		SourceRepository: sourceRepo,
	***REMOVED***

	meta.HMAC = metadata.ComputeV2MetadataHMAC([]byte(key), &meta)
	return meta
***REMOVED***

type mockRepo struct ***REMOVED***
	t        *testing.T
	errors   map[digest.Digest]error
	blobs    map[digest.Digest]distribution.Descriptor
	requests []string
***REMOVED***

var _ distribution.Repository = &mockRepo***REMOVED******REMOVED***

func (m *mockRepo) Named() reference.Named ***REMOVED***
	m.t.Fatalf("Named() not implemented")
	return nil
***REMOVED***
func (m *mockRepo) Manifests(ctc context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) ***REMOVED***
	m.t.Fatalf("Manifests() not implemented")
	return nil, nil
***REMOVED***
func (m *mockRepo) Tags(ctc context.Context) distribution.TagService ***REMOVED***
	m.t.Fatalf("Tags() not implemented")
	return nil
***REMOVED***
func (m *mockRepo) Blobs(ctx context.Context) distribution.BlobStore ***REMOVED***
	return &mockBlobStore***REMOVED***
		repo: m,
	***REMOVED***
***REMOVED***

type mockBlobStore struct ***REMOVED***
	repo *mockRepo
***REMOVED***

var _ distribution.BlobStore = &mockBlobStore***REMOVED******REMOVED***

func (m *mockBlobStore) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) ***REMOVED***
	m.repo.requests = append(m.repo.requests, dgst.String())
	if err, exists := m.repo.errors[dgst]; exists ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	if desc, exists := m.repo.blobs[dgst]; exists ***REMOVED***
		return desc, nil
	***REMOVED***
	return distribution.Descriptor***REMOVED******REMOVED***, distribution.ErrBlobUnknown
***REMOVED***
func (m *mockBlobStore) Get(ctx context.Context, dgst digest.Digest) ([]byte, error) ***REMOVED***
	m.repo.t.Fatal("Get() not implemented")
	return nil, nil
***REMOVED***

func (m *mockBlobStore) Open(ctx context.Context, dgst digest.Digest) (distribution.ReadSeekCloser, error) ***REMOVED***
	m.repo.t.Fatal("Open() not implemented")
	return nil, nil
***REMOVED***

func (m *mockBlobStore) Put(ctx context.Context, mediaType string, p []byte) (distribution.Descriptor, error) ***REMOVED***
	m.repo.t.Fatal("Put() not implemented")
	return distribution.Descriptor***REMOVED******REMOVED***, nil
***REMOVED***

func (m *mockBlobStore) Create(ctx context.Context, options ...distribution.BlobCreateOption) (distribution.BlobWriter, error) ***REMOVED***
	m.repo.t.Fatal("Create() not implemented")
	return nil, nil
***REMOVED***
func (m *mockBlobStore) Resume(ctx context.Context, id string) (distribution.BlobWriter, error) ***REMOVED***
	m.repo.t.Fatal("Resume() not implemented")
	return nil, nil
***REMOVED***
func (m *mockBlobStore) Delete(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	m.repo.t.Fatal("Delete() not implemented")
	return nil
***REMOVED***
func (m *mockBlobStore) ServeBlob(ctx context.Context, w http.ResponseWriter, r *http.Request, dgst digest.Digest) error ***REMOVED***
	m.repo.t.Fatalf("ServeBlob() not implemented")
	return nil
***REMOVED***

type mockV2MetadataService struct ***REMOVED***
	added   []metadata.V2Metadata
	removed []metadata.V2Metadata
***REMOVED***

var _ metadata.V2MetadataService = &mockV2MetadataService***REMOVED******REMOVED***

func (*mockV2MetadataService) GetMetadata(diffID layer.DiffID) ([]metadata.V2Metadata, error) ***REMOVED***
	return nil, nil
***REMOVED***
func (*mockV2MetadataService) GetDiffID(dgst digest.Digest) (layer.DiffID, error) ***REMOVED***
	return "", nil
***REMOVED***
func (m *mockV2MetadataService) Add(diffID layer.DiffID, metadata metadata.V2Metadata) error ***REMOVED***
	m.added = append(m.added, metadata)
	return nil
***REMOVED***
func (m *mockV2MetadataService) TagAndAdd(diffID layer.DiffID, hmacKey []byte, meta metadata.V2Metadata) error ***REMOVED***
	meta.HMAC = metadata.ComputeV2MetadataHMAC(hmacKey, &meta)
	m.Add(diffID, meta)
	return nil
***REMOVED***
func (m *mockV2MetadataService) Remove(metadata metadata.V2Metadata) error ***REMOVED***
	m.removed = append(m.removed, metadata)
	return nil
***REMOVED***

type progressSink struct ***REMOVED***
	t *testing.T
***REMOVED***

func (s *progressSink) WriteProgress(p progress.Progress) error ***REMOVED***
	s.t.Logf("progress update: %#+v", p)
	return nil
***REMOVED***
