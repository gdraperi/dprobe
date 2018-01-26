package daemon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	_ "github.com/docker/docker/pkg/discovery/memory"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/truncindex"
	"github.com/docker/docker/volume"
	volumedrivers "github.com/docker/docker/volume/drivers"
	"github.com/docker/docker/volume/local"
	"github.com/docker/docker/volume/store"
	"github.com/docker/go-connections/nat"
	"github.com/docker/libnetwork"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

//
// https://github.com/docker/docker/issues/8069
//

func TestGetContainer(t *testing.T) ***REMOVED***
	c1 := &container.Container***REMOVED***
		ID:   "5a4ff6a163ad4533d22d69a2b8960bf7fafdcba06e72d2febdba229008b0bf57",
		Name: "tender_bardeen",
	***REMOVED***

	c2 := &container.Container***REMOVED***
		ID:   "3cdbd1aa394fd68559fd1441d6eff2ab7c1e6363582c82febfaa8045df3bd8de",
		Name: "drunk_hawking",
	***REMOVED***

	c3 := &container.Container***REMOVED***
		ID:   "3cdbd1aa394fd68559fd1441d6eff2abfafdcba06e72d2febdba229008b0bf57",
		Name: "3cdbd1aa",
	***REMOVED***

	c4 := &container.Container***REMOVED***
		ID:   "75fb0b800922abdbef2d27e60abcdfaf7fb0698b2a96d22d3354da361a6ff4a5",
		Name: "5a4ff6a163ad4533d22d69a2b8960bf7fafdcba06e72d2febdba229008b0bf57",
	***REMOVED***

	c5 := &container.Container***REMOVED***
		ID:   "d22d69a2b8960bf7fafdcba06e72d2febdba960bf7fafdcba06e72d2f9008b060b",
		Name: "d22d69a2b896",
	***REMOVED***

	store := container.NewMemoryStore()
	store.Add(c1.ID, c1)
	store.Add(c2.ID, c2)
	store.Add(c3.ID, c3)
	store.Add(c4.ID, c4)
	store.Add(c5.ID, c5)

	index := truncindex.NewTruncIndex([]string***REMOVED******REMOVED***)
	index.Add(c1.ID)
	index.Add(c2.ID)
	index.Add(c3.ID)
	index.Add(c4.ID)
	index.Add(c5.ID)

	containersReplica, err := container.NewViewDB()
	if err != nil ***REMOVED***
		t.Fatalf("could not create ViewDB: %v", err)
	***REMOVED***

	daemon := &Daemon***REMOVED***
		containers:        store,
		containersReplica: containersReplica,
		idIndex:           index,
	***REMOVED***

	daemon.reserveName(c1.ID, c1.Name)
	daemon.reserveName(c2.ID, c2.Name)
	daemon.reserveName(c3.ID, c3.Name)
	daemon.reserveName(c4.ID, c4.Name)
	daemon.reserveName(c5.ID, c5.Name)

	if container, _ := daemon.GetContainer("3cdbd1aa394fd68559fd1441d6eff2ab7c1e6363582c82febfaa8045df3bd8de"); container != c2 ***REMOVED***
		t.Fatal("Should explicitly match full container IDs")
	***REMOVED***

	if container, _ := daemon.GetContainer("75fb0b8009"); container != c4 ***REMOVED***
		t.Fatal("Should match a partial ID")
	***REMOVED***

	if container, _ := daemon.GetContainer("drunk_hawking"); container != c2 ***REMOVED***
		t.Fatal("Should match a full name")
	***REMOVED***

	// c3.Name is a partial match for both c3.ID and c2.ID
	if c, _ := daemon.GetContainer("3cdbd1aa"); c != c3 ***REMOVED***
		t.Fatal("Should match a full name even though it collides with another container's ID")
	***REMOVED***

	if container, _ := daemon.GetContainer("d22d69a2b896"); container != c5 ***REMOVED***
		t.Fatal("Should match a container where the provided prefix is an exact match to the its name, and is also a prefix for its ID")
	***REMOVED***

	if _, err := daemon.GetContainer("3cdbd1"); err == nil ***REMOVED***
		t.Fatal("Should return an error when provided a prefix that partially matches multiple container ID's")
	***REMOVED***

	if _, err := daemon.GetContainer("nothing"); err == nil ***REMOVED***
		t.Fatal("Should return an error when provided a prefix that is neither a name or a partial match to an ID")
	***REMOVED***
***REMOVED***

func initDaemonWithVolumeStore(tmp string) (*Daemon, error) ***REMOVED***
	var err error
	daemon := &Daemon***REMOVED***
		repository: tmp,
		root:       tmp,
	***REMOVED***
	daemon.volumes, err = store.New(tmp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	volumesDriver, err := local.New(tmp, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	volumedrivers.Register(volumesDriver, volumesDriver.Name())

	return daemon, nil
***REMOVED***

func TestValidContainerNames(t *testing.T) ***REMOVED***
	invalidNames := []string***REMOVED***"-rm", "&sdfsfd", "safd%sd"***REMOVED***
	validNames := []string***REMOVED***"word-word", "word_word", "1weoid"***REMOVED***

	for _, name := range invalidNames ***REMOVED***
		if validContainerNamePattern.MatchString(name) ***REMOVED***
			t.Fatalf("%q is not a valid container name and was returned as valid.", name)
		***REMOVED***
	***REMOVED***

	for _, name := range validNames ***REMOVED***
		if !validContainerNamePattern.MatchString(name) ***REMOVED***
			t.Fatalf("%q is a valid container name and was returned as invalid.", name)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestContainerInitDNS(t *testing.T) ***REMOVED***
	tmp, err := ioutil.TempDir("", "docker-container-test-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)

	containerID := "d59df5276e7b219d510fe70565e0404bc06350e0d4b43fe961f22f339980170e"
	containerPath := filepath.Join(tmp, containerID)
	if err := os.MkdirAll(containerPath, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	config := `***REMOVED***"State":***REMOVED***"Running":true,"Paused":false,"Restarting":false,"OOMKilled":false,"Dead":false,"Pid":2464,"ExitCode":0,
"Error":"","StartedAt":"2015-05-26T16:48:53.869308965Z","FinishedAt":"0001-01-01T00:00:00Z"***REMOVED***,
"ID":"d59df5276e7b219d510fe70565e0404bc06350e0d4b43fe961f22f339980170e","Created":"2015-05-26T16:48:53.7987917Z","Path":"top",
"Args":[],"Config":***REMOVED***"Hostname":"d59df5276e7b","Domainname":"","User":"","Memory":0,"MemorySwap":0,"CpuShares":0,"Cpuset":"",
"AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"PortSpecs":null,"ExposedPorts":null,"Tty":true,"OpenStdin":true,
"StdinOnce":false,"Env":null,"Cmd":["top"],"Image":"ubuntu:latest","Volumes":null,"WorkingDir":"","Entrypoint":null,
"NetworkDisabled":false,"MacAddress":"","OnBuild":null,"Labels":***REMOVED******REMOVED******REMOVED***,"Image":"07f8e8c5e66084bef8f848877857537ffe1c47edd01a93af27e7161672ad0e95",
"NetworkSettings":***REMOVED***"IPAddress":"172.17.0.1","IPPrefixLen":16,"MacAddress":"02:42:ac:11:00:01","LinkLocalIPv6Address":"fe80::42:acff:fe11:1",
"LinkLocalIPv6PrefixLen":64,"GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"Gateway":"172.17.42.1","IPv6Gateway":"","Bridge":"docker0","Ports":***REMOVED******REMOVED******REMOVED***,
"ResolvConfPath":"/var/lib/docker/containers/d59df5276e7b219d510fe70565e0404bc06350e0d4b43fe961f22f339980170e/resolv.conf",
"HostnamePath":"/var/lib/docker/containers/d59df5276e7b219d510fe70565e0404bc06350e0d4b43fe961f22f339980170e/hostname",
"HostsPath":"/var/lib/docker/containers/d59df5276e7b219d510fe70565e0404bc06350e0d4b43fe961f22f339980170e/hosts",
"LogPath":"/var/lib/docker/containers/d59df5276e7b219d510fe70565e0404bc06350e0d4b43fe961f22f339980170e/d59df5276e7b219d510fe70565e0404bc06350e0d4b43fe961f22f339980170e-json.log",
"Name":"/ubuntu","Driver":"aufs","MountLabel":"","ProcessLabel":"","AppArmorProfile":"","RestartCount":0,
"UpdateDns":false,"Volumes":***REMOVED******REMOVED***,"VolumesRW":***REMOVED******REMOVED***,"AppliedVolumesFrom":null***REMOVED***`

	// Container struct only used to retrieve path to config file
	container := &container.Container***REMOVED***Root: containerPath***REMOVED***
	configPath, err := container.ConfigPath()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err = ioutil.WriteFile(configPath, []byte(config), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	hostConfig := `***REMOVED***"Binds":[],"ContainerIDFile":"","Memory":0,"MemorySwap":0,"CpuShares":0,"CpusetCpus":"",
"Privileged":false,"PortBindings":***REMOVED******REMOVED***,"Links":null,"PublishAllPorts":false,"Dns":null,"DnsOptions":null,"DnsSearch":null,"ExtraHosts":null,"VolumesFrom":null,
"Devices":[],"NetworkMode":"bridge","IpcMode":"","PidMode":"","CapAdd":null,"CapDrop":null,"RestartPolicy":***REMOVED***"Name":"no","MaximumRetryCount":0***REMOVED***,
"SecurityOpt":null,"ReadonlyRootfs":false,"Ulimits":null,"LogConfig":***REMOVED***"Type":"","Config":null***REMOVED***,"CgroupParent":""***REMOVED***`

	hostConfigPath, err := container.HostConfigPath()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err = ioutil.WriteFile(hostConfigPath, []byte(hostConfig), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	daemon, err := initDaemonWithVolumeStore(tmp)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer volumedrivers.Unregister(volume.DefaultDriverName)

	c, err := daemon.load(containerID)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if c.HostConfig.DNS == nil ***REMOVED***
		t.Fatal("Expected container DNS to not be nil")
	***REMOVED***

	if c.HostConfig.DNSSearch == nil ***REMOVED***
		t.Fatal("Expected container DNSSearch to not be nil")
	***REMOVED***

	if c.HostConfig.DNSOptions == nil ***REMOVED***
		t.Fatal("Expected container DNSOptions to not be nil")
	***REMOVED***
***REMOVED***

func newPortNoError(proto, port string) nat.Port ***REMOVED***
	p, _ := nat.NewPort(proto, port)
	return p
***REMOVED***

func TestMerge(t *testing.T) ***REMOVED***
	volumesImage := make(map[string]struct***REMOVED******REMOVED***)
	volumesImage["/test1"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	volumesImage["/test2"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	portsImage := make(nat.PortSet)
	portsImage[newPortNoError("tcp", "1111")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	portsImage[newPortNoError("tcp", "2222")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	configImage := &containertypes.Config***REMOVED***
		ExposedPorts: portsImage,
		Env:          []string***REMOVED***"VAR1=1", "VAR2=2"***REMOVED***,
		Volumes:      volumesImage,
	***REMOVED***

	portsUser := make(nat.PortSet)
	portsUser[newPortNoError("tcp", "2222")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	portsUser[newPortNoError("tcp", "3333")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	volumesUser := make(map[string]struct***REMOVED******REMOVED***)
	volumesUser["/test3"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	configUser := &containertypes.Config***REMOVED***
		ExposedPorts: portsUser,
		Env:          []string***REMOVED***"VAR2=3", "VAR3=3"***REMOVED***,
		Volumes:      volumesUser,
	***REMOVED***

	if err := merge(configUser, configImage); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	if len(configUser.ExposedPorts) != 3 ***REMOVED***
		t.Fatalf("Expected 3 ExposedPorts, 1111, 2222 and 3333, found %d", len(configUser.ExposedPorts))
	***REMOVED***
	for portSpecs := range configUser.ExposedPorts ***REMOVED***
		if portSpecs.Port() != "1111" && portSpecs.Port() != "2222" && portSpecs.Port() != "3333" ***REMOVED***
			t.Fatalf("Expected 1111 or 2222 or 3333, found %s", portSpecs)
		***REMOVED***
	***REMOVED***
	if len(configUser.Env) != 3 ***REMOVED***
		t.Fatalf("Expected 3 env var, VAR1=1, VAR2=3 and VAR3=3, found %d", len(configUser.Env))
	***REMOVED***
	for _, env := range configUser.Env ***REMOVED***
		if env != "VAR1=1" && env != "VAR2=3" && env != "VAR3=3" ***REMOVED***
			t.Fatalf("Expected VAR1=1 or VAR2=3 or VAR3=3, found %s", env)
		***REMOVED***
	***REMOVED***

	if len(configUser.Volumes) != 3 ***REMOVED***
		t.Fatalf("Expected 3 volumes, /test1, /test2 and /test3, found %d", len(configUser.Volumes))
	***REMOVED***
	for v := range configUser.Volumes ***REMOVED***
		if v != "/test1" && v != "/test2" && v != "/test3" ***REMOVED***
			t.Fatalf("Expected /test1 or /test2 or /test3, found %s", v)
		***REMOVED***
	***REMOVED***

	ports, _, err := nat.ParsePortSpecs([]string***REMOVED***"0000"***REMOVED***)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	configImage2 := &containertypes.Config***REMOVED***
		ExposedPorts: ports,
	***REMOVED***

	if err := merge(configUser, configImage2); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	if len(configUser.ExposedPorts) != 4 ***REMOVED***
		t.Fatalf("Expected 4 ExposedPorts, 0000, 1111, 2222 and 3333, found %d", len(configUser.ExposedPorts))
	***REMOVED***
	for portSpecs := range configUser.ExposedPorts ***REMOVED***
		if portSpecs.Port() != "0" && portSpecs.Port() != "1111" && portSpecs.Port() != "2222" && portSpecs.Port() != "3333" ***REMOVED***
			t.Fatalf("Expected %q or %q or %q or %q, found %s", 0, 1111, 2222, 3333, portSpecs)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestValidateContainerIsolation(t *testing.T) ***REMOVED***
	d := Daemon***REMOVED******REMOVED***

	_, err := d.verifyContainerSettings(runtime.GOOS, &containertypes.HostConfig***REMOVED***Isolation: containertypes.Isolation("invalid")***REMOVED***, nil, false)
	assert.EqualError(t, err, "invalid isolation 'invalid' on "+runtime.GOOS)
***REMOVED***

func TestFindNetworkErrorType(t *testing.T) ***REMOVED***
	d := Daemon***REMOVED******REMOVED***
	_, err := d.FindNetwork("fakeNet")
	_, ok := errors.Cause(err).(libnetwork.ErrNoSuchNetwork)
	if !errdefs.IsNotFound(err) || !ok ***REMOVED***
		assert.Fail(t, "The FindNetwork method MUST always return an error that implements the NotFound interface and is ErrNoSuchNetwork")
	***REMOVED***
***REMOVED***
