// +build !windows

package daemon

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/volume"
	"github.com/docker/docker/volume/drivers"
	"github.com/docker/docker/volume/local"
	"github.com/docker/docker/volume/store"
	"github.com/stretchr/testify/require"
)

type fakeContainerGetter struct ***REMOVED***
	containers map[string]*container.Container
***REMOVED***

func (f *fakeContainerGetter) GetContainer(cid string) (*container.Container, error) ***REMOVED***
	container, ok := f.containers[cid]
	if !ok ***REMOVED***
		return nil, errors.New("container not found")
	***REMOVED***
	return container, nil
***REMOVED***

// Unix test as uses settings which are not available on Windows
func TestAdjustSharedNamespaceContainerName(t *testing.T) ***REMOVED***
	fakeID := "abcdef1234567890"
	hostConfig := &containertypes.HostConfig***REMOVED***
		IpcMode:     containertypes.IpcMode("container:base"),
		PidMode:     containertypes.PidMode("container:base"),
		NetworkMode: containertypes.NetworkMode("container:base"),
	***REMOVED***
	containerStore := &fakeContainerGetter***REMOVED******REMOVED***
	containerStore.containers = make(map[string]*container.Container)
	containerStore.containers["base"] = &container.Container***REMOVED***
		ID: fakeID,
	***REMOVED***

	adaptSharedNamespaceContainer(containerStore, hostConfig)
	if hostConfig.IpcMode != containertypes.IpcMode("container:"+fakeID) ***REMOVED***
		t.Errorf("Expected IpcMode to be container:%s", fakeID)
	***REMOVED***
	if hostConfig.PidMode != containertypes.PidMode("container:"+fakeID) ***REMOVED***
		t.Errorf("Expected PidMode to be container:%s", fakeID)
	***REMOVED***
	if hostConfig.NetworkMode != containertypes.NetworkMode("container:"+fakeID) ***REMOVED***
		t.Errorf("Expected NetworkMode to be container:%s", fakeID)
	***REMOVED***
***REMOVED***

// Unix test as uses settings which are not available on Windows
func TestAdjustCPUShares(t *testing.T) ***REMOVED***
	tmp, err := ioutil.TempDir("", "docker-daemon-unix-test-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)
	daemon := &Daemon***REMOVED***
		repository: tmp,
		root:       tmp,
	***REMOVED***

	hostConfig := &containertypes.HostConfig***REMOVED***
		Resources: containertypes.Resources***REMOVED***CPUShares: linuxMinCPUShares - 1***REMOVED***,
	***REMOVED***
	daemon.adaptContainerSettings(hostConfig, true)
	if hostConfig.CPUShares != linuxMinCPUShares ***REMOVED***
		t.Errorf("Expected CPUShares to be %d", linuxMinCPUShares)
	***REMOVED***

	hostConfig.CPUShares = linuxMaxCPUShares + 1
	daemon.adaptContainerSettings(hostConfig, true)
	if hostConfig.CPUShares != linuxMaxCPUShares ***REMOVED***
		t.Errorf("Expected CPUShares to be %d", linuxMaxCPUShares)
	***REMOVED***

	hostConfig.CPUShares = 0
	daemon.adaptContainerSettings(hostConfig, true)
	if hostConfig.CPUShares != 0 ***REMOVED***
		t.Error("Expected CPUShares to be unchanged")
	***REMOVED***

	hostConfig.CPUShares = 1024
	daemon.adaptContainerSettings(hostConfig, true)
	if hostConfig.CPUShares != 1024 ***REMOVED***
		t.Error("Expected CPUShares to be unchanged")
	***REMOVED***
***REMOVED***

// Unix test as uses settings which are not available on Windows
func TestAdjustCPUSharesNoAdjustment(t *testing.T) ***REMOVED***
	tmp, err := ioutil.TempDir("", "docker-daemon-unix-test-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)
	daemon := &Daemon***REMOVED***
		repository: tmp,
		root:       tmp,
	***REMOVED***

	hostConfig := &containertypes.HostConfig***REMOVED***
		Resources: containertypes.Resources***REMOVED***CPUShares: linuxMinCPUShares - 1***REMOVED***,
	***REMOVED***
	daemon.adaptContainerSettings(hostConfig, false)
	if hostConfig.CPUShares != linuxMinCPUShares-1 ***REMOVED***
		t.Errorf("Expected CPUShares to be %d", linuxMinCPUShares-1)
	***REMOVED***

	hostConfig.CPUShares = linuxMaxCPUShares + 1
	daemon.adaptContainerSettings(hostConfig, false)
	if hostConfig.CPUShares != linuxMaxCPUShares+1 ***REMOVED***
		t.Errorf("Expected CPUShares to be %d", linuxMaxCPUShares+1)
	***REMOVED***

	hostConfig.CPUShares = 0
	daemon.adaptContainerSettings(hostConfig, false)
	if hostConfig.CPUShares != 0 ***REMOVED***
		t.Error("Expected CPUShares to be unchanged")
	***REMOVED***

	hostConfig.CPUShares = 1024
	daemon.adaptContainerSettings(hostConfig, false)
	if hostConfig.CPUShares != 1024 ***REMOVED***
		t.Error("Expected CPUShares to be unchanged")
	***REMOVED***
***REMOVED***

// Unix test as uses settings which are not available on Windows
func TestParseSecurityOptWithDeprecatedColon(t *testing.T) ***REMOVED***
	container := &container.Container***REMOVED******REMOVED***
	config := &containertypes.HostConfig***REMOVED******REMOVED***

	// test apparmor
	config.SecurityOpt = []string***REMOVED***"apparmor=test_profile"***REMOVED***
	if err := parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	***REMOVED***
	if container.AppArmorProfile != "test_profile" ***REMOVED***
		t.Fatalf("Unexpected AppArmorProfile, expected: \"test_profile\", got %q", container.AppArmorProfile)
	***REMOVED***

	// test seccomp
	sp := "/path/to/seccomp_test.json"
	config.SecurityOpt = []string***REMOVED***"seccomp=" + sp***REMOVED***
	if err := parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	***REMOVED***
	if container.SeccompProfile != sp ***REMOVED***
		t.Fatalf("Unexpected AppArmorProfile, expected: %q, got %q", sp, container.SeccompProfile)
	***REMOVED***

	// test valid label
	config.SecurityOpt = []string***REMOVED***"label=user:USER"***REMOVED***
	if err := parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	***REMOVED***

	// test invalid label
	config.SecurityOpt = []string***REMOVED***"label"***REMOVED***
	if err := parseSecurityOpt(container, config); err == nil ***REMOVED***
		t.Fatal("Expected parseSecurityOpt error, got nil")
	***REMOVED***

	// test invalid opt
	config.SecurityOpt = []string***REMOVED***"test"***REMOVED***
	if err := parseSecurityOpt(container, config); err == nil ***REMOVED***
		t.Fatal("Expected parseSecurityOpt error, got nil")
	***REMOVED***
***REMOVED***

func TestParseSecurityOpt(t *testing.T) ***REMOVED***
	container := &container.Container***REMOVED******REMOVED***
	config := &containertypes.HostConfig***REMOVED******REMOVED***

	// test apparmor
	config.SecurityOpt = []string***REMOVED***"apparmor=test_profile"***REMOVED***
	if err := parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	***REMOVED***
	if container.AppArmorProfile != "test_profile" ***REMOVED***
		t.Fatalf("Unexpected AppArmorProfile, expected: \"test_profile\", got %q", container.AppArmorProfile)
	***REMOVED***

	// test seccomp
	sp := "/path/to/seccomp_test.json"
	config.SecurityOpt = []string***REMOVED***"seccomp=" + sp***REMOVED***
	if err := parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	***REMOVED***
	if container.SeccompProfile != sp ***REMOVED***
		t.Fatalf("Unexpected SeccompProfile, expected: %q, got %q", sp, container.SeccompProfile)
	***REMOVED***

	// test valid label
	config.SecurityOpt = []string***REMOVED***"label=user:USER"***REMOVED***
	if err := parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	***REMOVED***

	// test invalid label
	config.SecurityOpt = []string***REMOVED***"label"***REMOVED***
	if err := parseSecurityOpt(container, config); err == nil ***REMOVED***
		t.Fatal("Expected parseSecurityOpt error, got nil")
	***REMOVED***

	// test invalid opt
	config.SecurityOpt = []string***REMOVED***"test"***REMOVED***
	if err := parseSecurityOpt(container, config); err == nil ***REMOVED***
		t.Fatal("Expected parseSecurityOpt error, got nil")
	***REMOVED***
***REMOVED***

func TestParseNNPSecurityOptions(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED***
		configStore: &config.Config***REMOVED***NoNewPrivileges: true***REMOVED***,
	***REMOVED***
	container := &container.Container***REMOVED******REMOVED***
	config := &containertypes.HostConfig***REMOVED******REMOVED***

	// test NNP when "daemon:true" and "no-new-privileges=false""
	config.SecurityOpt = []string***REMOVED***"no-new-privileges=false"***REMOVED***

	if err := daemon.parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected daemon.parseSecurityOpt error: %v", err)
	***REMOVED***
	if container.NoNewPrivileges ***REMOVED***
		t.Fatalf("container.NoNewPrivileges should be FALSE: %v", container.NoNewPrivileges)
	***REMOVED***

	// test NNP when "daemon:false" and "no-new-privileges=true""
	daemon.configStore.NoNewPrivileges = false
	config.SecurityOpt = []string***REMOVED***"no-new-privileges=true"***REMOVED***

	if err := daemon.parseSecurityOpt(container, config); err != nil ***REMOVED***
		t.Fatalf("Unexpected daemon.parseSecurityOpt error: %v", err)
	***REMOVED***
	if !container.NoNewPrivileges ***REMOVED***
		t.Fatalf("container.NoNewPrivileges should be TRUE: %v", container.NoNewPrivileges)
	***REMOVED***
***REMOVED***

func TestNetworkOptions(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	dconfigCorrect := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ClusterStore:     "consul://localhost:8500",
			ClusterAdvertise: "192.168.0.1:8000",
		***REMOVED***,
	***REMOVED***

	if _, err := daemon.networkOptions(dconfigCorrect, nil, nil); err != nil ***REMOVED***
		t.Fatalf("Expect networkOptions success, got error: %v", err)
	***REMOVED***

	dconfigWrong := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ClusterStore: "consul://localhost:8500://test://bbb",
		***REMOVED***,
	***REMOVED***

	if _, err := daemon.networkOptions(dconfigWrong, nil, nil); err == nil ***REMOVED***
		t.Fatal("Expected networkOptions error, got nil")
	***REMOVED***
***REMOVED***

func TestMigratePre17Volumes(t *testing.T) ***REMOVED***
	rootDir, err := ioutil.TempDir("", "test-daemon-volumes")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(rootDir)

	volumeRoot := filepath.Join(rootDir, "volumes")
	err = os.MkdirAll(volumeRoot, 0755)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	containerRoot := filepath.Join(rootDir, "containers")
	cid := "1234"
	err = os.MkdirAll(filepath.Join(containerRoot, cid), 0755)
	require.NoError(t, err)

	vid := "5678"
	vfsPath := filepath.Join(rootDir, "vfs", "dir", vid)
	err = os.MkdirAll(vfsPath, 0755)
	require.NoError(t, err)

	config := []byte(`
		***REMOVED***
			"ID": "` + cid + `",
			"Volumes": ***REMOVED***
				"/foo": "` + vfsPath + `",
				"/bar": "/foo",
				"/quux": "/quux"
			***REMOVED***,
			"VolumesRW": ***REMOVED***
				"/foo": true,
				"/bar": true,
				"/quux": false
			***REMOVED***
		***REMOVED***
	`)

	volStore, err := store.New(volumeRoot)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	drv, err := local.New(volumeRoot, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	volumedrivers.Register(drv, volume.DefaultDriverName)

	daemon := &Daemon***REMOVED***
		root:       rootDir,
		repository: containerRoot,
		volumes:    volStore,
	***REMOVED***
	err = ioutil.WriteFile(filepath.Join(containerRoot, cid, "config.v2.json"), config, 600)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	c, err := daemon.load(cid)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := daemon.verifyVolumesInfo(c); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected := map[string]volume.MountPoint***REMOVED***
		"/foo":  ***REMOVED***Destination: "/foo", RW: true, Name: vid***REMOVED***,
		"/bar":  ***REMOVED***Source: "/foo", Destination: "/bar", RW: true***REMOVED***,
		"/quux": ***REMOVED***Source: "/quux", Destination: "/quux", RW: false***REMOVED***,
	***REMOVED***
	for id, mp := range c.MountPoints ***REMOVED***
		x, exists := expected[id]
		if !exists ***REMOVED***
			t.Fatal("volume not migrated")
		***REMOVED***
		if mp.Source != x.Source || mp.Destination != x.Destination || mp.RW != x.RW || mp.Name != x.Name ***REMOVED***
			t.Fatalf("got unexpected mountpoint, expected: %+v, got: %+v", x, mp)
		***REMOVED***
	***REMOVED***
***REMOVED***
