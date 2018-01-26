package runconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type f struct ***REMOVED***
	file       string
	entrypoint strslice.StrSlice
***REMOVED***

func TestDecodeContainerConfig(t *testing.T) ***REMOVED***

	var (
		fixtures []f
		image    string
	)

	if runtime.GOOS != "windows" ***REMOVED***
		image = "ubuntu"
		fixtures = []f***REMOVED***
			***REMOVED***"fixtures/unix/container_config_1_14.json", strslice.StrSlice***REMOVED******REMOVED******REMOVED***,
			***REMOVED***"fixtures/unix/container_config_1_17.json", strslice.StrSlice***REMOVED***"bash"***REMOVED******REMOVED***,
			***REMOVED***"fixtures/unix/container_config_1_19.json", strslice.StrSlice***REMOVED***"bash"***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		image = "windows"
		fixtures = []f***REMOVED***
			***REMOVED***"fixtures/windows/container_config_1_19.json", strslice.StrSlice***REMOVED***"cmd"***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED***

	for _, f := range fixtures ***REMOVED***
		b, err := ioutil.ReadFile(f.file)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		c, h, _, err := decodeContainerConfig(bytes.NewReader(b))
		if err != nil ***REMOVED***
			t.Fatal(fmt.Errorf("Error parsing %s: %v", f, err))
		***REMOVED***

		if c.Image != image ***REMOVED***
			t.Fatalf("Expected %s image, found %s\n", image, c.Image)
		***REMOVED***

		if len(c.Entrypoint) != len(f.entrypoint) ***REMOVED***
			t.Fatalf("Expected %v, found %v\n", f.entrypoint, c.Entrypoint)
		***REMOVED***

		if h != nil && h.Memory != 1000 ***REMOVED***
			t.Fatalf("Expected memory to be 1000, found %d\n", h.Memory)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestDecodeContainerConfigIsolation validates isolation passed
// to the daemon in the hostConfig structure. Note this is platform specific
// as to what level of container isolation is supported.
func TestDecodeContainerConfigIsolation(t *testing.T) ***REMOVED***

	// An Invalid isolation level
	if _, _, _, err := callDecodeContainerConfigIsolation("invalid"); err != nil ***REMOVED***
		if !strings.Contains(err.Error(), `Invalid isolation: "invalid"`) ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	// Blank isolation (== default)
	if _, _, _, err := callDecodeContainerConfigIsolation(""); err != nil ***REMOVED***
		t.Fatal("Blank isolation should have succeeded")
	***REMOVED***

	// Default isolation
	if _, _, _, err := callDecodeContainerConfigIsolation("default"); err != nil ***REMOVED***
		t.Fatal("default isolation should have succeeded")
	***REMOVED***

	// Process isolation (Valid on Windows only)
	if runtime.GOOS == "windows" ***REMOVED***
		if _, _, _, err := callDecodeContainerConfigIsolation("process"); err != nil ***REMOVED***
			t.Fatal("process isolation should have succeeded")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if _, _, _, err := callDecodeContainerConfigIsolation("process"); err != nil ***REMOVED***
			if !strings.Contains(err.Error(), `Invalid isolation: "process"`) ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Hyper-V Containers isolation (Valid on Windows only)
	if runtime.GOOS == "windows" ***REMOVED***
		if _, _, _, err := callDecodeContainerConfigIsolation("hyperv"); err != nil ***REMOVED***
			t.Fatal("hyperv isolation should have succeeded")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if _, _, _, err := callDecodeContainerConfigIsolation("hyperv"); err != nil ***REMOVED***
			if !strings.Contains(err.Error(), `Invalid isolation: "hyperv"`) ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// callDecodeContainerConfigIsolation is a utility function to call
// DecodeContainerConfig for validating isolation
func callDecodeContainerConfigIsolation(isolation string) (*container.Config, *container.HostConfig, *networktypes.NetworkingConfig, error) ***REMOVED***
	var (
		b   []byte
		err error
	)
	w := ContainerConfigWrapper***REMOVED***
		Config: &container.Config***REMOVED******REMOVED***,
		HostConfig: &container.HostConfig***REMOVED***
			NetworkMode: "none",
			Isolation:   container.Isolation(isolation)***REMOVED***,
	***REMOVED***
	if b, err = json.Marshal(w); err != nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("Error on marshal %s", err.Error())
	***REMOVED***
	return decodeContainerConfig(bytes.NewReader(b))
***REMOVED***

type decodeConfigTestcase struct ***REMOVED***
	doc                string
	wrapper            ContainerConfigWrapper
	expectedErr        string
	expectedConfig     *container.Config
	expectedHostConfig *container.HostConfig
	goos               string
***REMOVED***

func runDecodeContainerConfigTestCase(testcase decodeConfigTestcase) func(t *testing.T) ***REMOVED***
	return func(t *testing.T) ***REMOVED***
		raw := marshal(t, testcase.wrapper, testcase.doc)
		config, hostConfig, _, err := decodeContainerConfig(bytes.NewReader(raw))
		if testcase.expectedErr != "" ***REMOVED***
			if !assert.Error(t, err) ***REMOVED***
				return
			***REMOVED***
			assert.Contains(t, err.Error(), testcase.expectedErr)
			return
		***REMOVED***
		assert.NoError(t, err)
		assert.Equal(t, testcase.expectedConfig, config)
		assert.Equal(t, testcase.expectedHostConfig, hostConfig)
	***REMOVED***
***REMOVED***

func marshal(t *testing.T, w ContainerConfigWrapper, doc string) []byte ***REMOVED***
	b, err := json.Marshal(w)
	require.NoError(t, err, "%s: failed to encode config wrapper", doc)
	return b
***REMOVED***

func containerWrapperWithVolume(volume string) ContainerConfigWrapper ***REMOVED***
	return ContainerConfigWrapper***REMOVED***
		Config: &container.Config***REMOVED***
			Volumes: map[string]struct***REMOVED******REMOVED******REMOVED***
				volume: ***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		HostConfig: &container.HostConfig***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func containerWrapperWithBind(bind string) ContainerConfigWrapper ***REMOVED***
	return ContainerConfigWrapper***REMOVED***
		Config: &container.Config***REMOVED***
			Volumes: map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		HostConfig: &container.HostConfig***REMOVED***
			Binds: []string***REMOVED***bind***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***
