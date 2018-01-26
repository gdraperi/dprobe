package discovery

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryOptsErrors(t *testing.T) ***REMOVED***
	var testcases = []struct ***REMOVED***
		doc  string
		opts map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			doc:  "discovery.ttl < discovery.heartbeat",
			opts: map[string]string***REMOVED***"discovery.heartbeat": "10", "discovery.ttl": "5"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			doc:  "discovery.ttl == discovery.heartbeat",
			opts: map[string]string***REMOVED***"discovery.heartbeat": "10", "discovery.ttl": "10"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			doc:  "negative discovery.heartbeat",
			opts: map[string]string***REMOVED***"discovery.heartbeat": "-10", "discovery.ttl": "10"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			doc:  "negative discovery.ttl",
			opts: map[string]string***REMOVED***"discovery.heartbeat": "10", "discovery.ttl": "-10"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			doc:  "invalid discovery.heartbeat",
			opts: map[string]string***REMOVED***"discovery.heartbeat": "invalid"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			doc:  "invalid discovery.ttl",
			opts: map[string]string***REMOVED***"discovery.ttl": "invalid"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, testcase := range testcases ***REMOVED***
		_, _, err := discoveryOpts(testcase.opts)
		assert.Error(t, err, testcase.doc)
	***REMOVED***
***REMOVED***

func TestDiscoveryOpts(t *testing.T) ***REMOVED***
	clusterOpts := map[string]string***REMOVED***"discovery.heartbeat": "10", "discovery.ttl": "20"***REMOVED***
	heartbeat, ttl, err := discoveryOpts(clusterOpts)
	require.NoError(t, err)
	assert.Equal(t, 10*time.Second, heartbeat)
	assert.Equal(t, 20*time.Second, ttl)

	clusterOpts = map[string]string***REMOVED***"discovery.heartbeat": "10"***REMOVED***
	heartbeat, ttl, err = discoveryOpts(clusterOpts)
	require.NoError(t, err)
	assert.Equal(t, 10*time.Second, heartbeat)
	assert.Equal(t, 10*defaultDiscoveryTTLFactor*time.Second, ttl)

	clusterOpts = map[string]string***REMOVED***"discovery.ttl": "30"***REMOVED***
	heartbeat, ttl, err = discoveryOpts(clusterOpts)
	require.NoError(t, err)

	if ttl != 30*time.Second ***REMOVED***
		t.Fatalf("TTL - Expected : %v, Actual : %v", 30*time.Second, ttl)
	***REMOVED***

	expected := 30 * time.Second / defaultDiscoveryTTLFactor
	if heartbeat != expected ***REMOVED***
		t.Fatalf("Heartbeat - Expected : %v, Actual : %v", expected, heartbeat)
	***REMOVED***

	discoveryTTL := fmt.Sprintf("%d", defaultDiscoveryTTLFactor-1)
	clusterOpts = map[string]string***REMOVED***"discovery.ttl": discoveryTTL***REMOVED***
	heartbeat, _, err = discoveryOpts(clusterOpts)
	if err == nil && heartbeat == 0 ***REMOVED***
		t.Fatal("discovery.heartbeat must be positive")
	***REMOVED***

	clusterOpts = map[string]string***REMOVED******REMOVED***
	heartbeat, ttl, err = discoveryOpts(clusterOpts)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if heartbeat != defaultDiscoveryHeartbeat ***REMOVED***
		t.Fatalf("Heartbeat - Expected : %v, Actual : %v", defaultDiscoveryHeartbeat, heartbeat)
	***REMOVED***

	expected = defaultDiscoveryHeartbeat * defaultDiscoveryTTLFactor
	if ttl != expected ***REMOVED***
		t.Fatalf("TTL - Expected : %v, Actual : %v", expected, ttl)
	***REMOVED***
***REMOVED***
