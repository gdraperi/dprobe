package main

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/config"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/cli/build/fakestorage"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/integration-cli/environment"
	"github.com/docker/docker/integration-cli/fixtures/plugin"
	"github.com/docker/docker/integration-cli/registry"
	ienv "github.com/docker/docker/internal/test/environment"
	"github.com/docker/docker/pkg/reexec"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

const (
	// the private registry to use for tests
	privateRegistryURL = "127.0.0.1:5000"

	// path to containerd's ctr binary
	ctrBinary = "docker-containerd-ctr"

	// the docker daemon binary to use
	dockerdBinary = "dockerd"
)

var (
	testEnv *environment.Execution

	// the docker client binary to use
	dockerBinary = ""
)

func init() ***REMOVED***
	var err error

	reexec.Init() // This is required for external graphdriver tests

	testEnv, err = environment.New()
	if err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***
***REMOVED***

func TestMain(m *testing.M) ***REMOVED***
	dockerBinary = testEnv.DockerBinary()
	err := ienv.EnsureFrozenImagesLinux(&testEnv.Execution)
	if err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***

	testEnv.Print()
	os.Exit(m.Run())
***REMOVED***

func Test(t *testing.T) ***REMOVED***
	cli.SetTestEnvironment(testEnv)
	fakestorage.SetTestEnvironment(&testEnv.Execution)
	ienv.ProtectAll(t, &testEnv.Execution)
	check.TestingT(t)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerSuite***REMOVED******REMOVED***)
***REMOVED***

type DockerSuite struct ***REMOVED***
***REMOVED***

func (s *DockerSuite) OnTimeout(c *check.C) ***REMOVED***
	if !testEnv.IsLocalDaemon() ***REMOVED***
		return
	***REMOVED***
	path := filepath.Join(os.Getenv("DEST"), "docker.pid")
	b, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		c.Fatalf("Failed to get daemon PID from %s\n", path)
	***REMOVED***

	rawPid, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil ***REMOVED***
		c.Fatalf("Failed to parse pid from %s: %s\n", path, err)
	***REMOVED***

	daemonPid := int(rawPid)
	if daemonPid > 0 ***REMOVED***
		daemon.SignalDaemonDump(daemonPid)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TearDownTest(c *check.C) ***REMOVED***
	testEnv.Clean(c)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerRegistrySuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerRegistrySuite struct ***REMOVED***
	ds  *DockerSuite
	reg *registry.V2
	d   *daemon.Daemon
***REMOVED***

func (s *DockerRegistrySuite) OnTimeout(c *check.C) ***REMOVED***
	s.d.DumpStackAndQuit()
***REMOVED***

func (s *DockerRegistrySuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, registry.Hosting, SameHostDaemon)
	s.reg = setupRegistry(c, false, "", "")
	s.d = daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)
***REMOVED***

func (s *DockerRegistrySuite) TearDownTest(c *check.C) ***REMOVED***
	if s.reg != nil ***REMOVED***
		s.reg.Close()
	***REMOVED***
	if s.d != nil ***REMOVED***
		s.d.Stop(c)
	***REMOVED***
	s.ds.TearDownTest(c)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerSchema1RegistrySuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerSchema1RegistrySuite struct ***REMOVED***
	ds  *DockerSuite
	reg *registry.V2
	d   *daemon.Daemon
***REMOVED***

func (s *DockerSchema1RegistrySuite) OnTimeout(c *check.C) ***REMOVED***
	s.d.DumpStackAndQuit()
***REMOVED***

func (s *DockerSchema1RegistrySuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, registry.Hosting, NotArm64, SameHostDaemon)
	s.reg = setupRegistry(c, true, "", "")
	s.d = daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TearDownTest(c *check.C) ***REMOVED***
	if s.reg != nil ***REMOVED***
		s.reg.Close()
	***REMOVED***
	if s.d != nil ***REMOVED***
		s.d.Stop(c)
	***REMOVED***
	s.ds.TearDownTest(c)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerRegistryAuthHtpasswdSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerRegistryAuthHtpasswdSuite struct ***REMOVED***
	ds  *DockerSuite
	reg *registry.V2
	d   *daemon.Daemon
***REMOVED***

func (s *DockerRegistryAuthHtpasswdSuite) OnTimeout(c *check.C) ***REMOVED***
	s.d.DumpStackAndQuit()
***REMOVED***

func (s *DockerRegistryAuthHtpasswdSuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, registry.Hosting, SameHostDaemon)
	s.reg = setupRegistry(c, false, "htpasswd", "")
	s.d = daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)
***REMOVED***

func (s *DockerRegistryAuthHtpasswdSuite) TearDownTest(c *check.C) ***REMOVED***
	if s.reg != nil ***REMOVED***
		out, err := s.d.Cmd("logout", privateRegistryURL)
		c.Assert(err, check.IsNil, check.Commentf(out))
		s.reg.Close()
	***REMOVED***
	if s.d != nil ***REMOVED***
		s.d.Stop(c)
	***REMOVED***
	s.ds.TearDownTest(c)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerRegistryAuthTokenSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerRegistryAuthTokenSuite struct ***REMOVED***
	ds  *DockerSuite
	reg *registry.V2
	d   *daemon.Daemon
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) OnTimeout(c *check.C) ***REMOVED***
	s.d.DumpStackAndQuit()
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, registry.Hosting, SameHostDaemon)
	s.d = daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) TearDownTest(c *check.C) ***REMOVED***
	if s.reg != nil ***REMOVED***
		out, err := s.d.Cmd("logout", privateRegistryURL)
		c.Assert(err, check.IsNil, check.Commentf(out))
		s.reg.Close()
	***REMOVED***
	if s.d != nil ***REMOVED***
		s.d.Stop(c)
	***REMOVED***
	s.ds.TearDownTest(c)
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) setupRegistryWithTokenService(c *check.C, tokenURL string) ***REMOVED***
	if s == nil ***REMOVED***
		c.Fatal("registry suite isn't initialized")
	***REMOVED***
	s.reg = setupRegistry(c, false, "token", tokenURL)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerDaemonSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerDaemonSuite struct ***REMOVED***
	ds *DockerSuite
	d  *daemon.Daemon
***REMOVED***

func (s *DockerDaemonSuite) OnTimeout(c *check.C) ***REMOVED***
	s.d.DumpStackAndQuit()
***REMOVED***

func (s *DockerDaemonSuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon)
	s.d = daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)
***REMOVED***

func (s *DockerDaemonSuite) TearDownTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon)
	if s.d != nil ***REMOVED***
		s.d.Stop(c)
	***REMOVED***
	s.ds.TearDownTest(c)
***REMOVED***

func (s *DockerDaemonSuite) TearDownSuite(c *check.C) ***REMOVED***
	filepath.Walk(daemon.SockRoot, func(path string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			// ignore errors here
			// not cleaning up sockets is not really an error
			return nil
		***REMOVED***
		if fi.Mode() == os.ModeSocket ***REMOVED***
			syscall.Unlink(path)
		***REMOVED***
		return nil
	***REMOVED***)
	os.RemoveAll(daemon.SockRoot)
***REMOVED***

const defaultSwarmPort = 2477

func init() ***REMOVED***
	check.Suite(&DockerSwarmSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerSwarmSuite struct ***REMOVED***
	server      *httptest.Server
	ds          *DockerSuite
	daemons     []*daemon.Swarm
	daemonsLock sync.Mutex // protect access to daemons
	portIndex   int
***REMOVED***

func (s *DockerSwarmSuite) OnTimeout(c *check.C) ***REMOVED***
	s.daemonsLock.Lock()
	defer s.daemonsLock.Unlock()
	for _, d := range s.daemons ***REMOVED***
		d.DumpStackAndQuit()
	***REMOVED***
***REMOVED***

func (s *DockerSwarmSuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon)
***REMOVED***

func (s *DockerSwarmSuite) AddDaemon(c *check.C, joinSwarm, manager bool) *daemon.Swarm ***REMOVED***
	d := &daemon.Swarm***REMOVED***
		Daemon: daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
			Experimental: testEnv.DaemonInfo.ExperimentalBuild,
		***REMOVED***),
		Port: defaultSwarmPort + s.portIndex,
	***REMOVED***
	d.ListenAddr = fmt.Sprintf("0.0.0.0:%d", d.Port)
	args := []string***REMOVED***"--iptables=false", "--swarm-default-advertise-addr=lo"***REMOVED*** // avoid networking conflicts
	d.StartWithBusybox(c, args...)

	if joinSwarm ***REMOVED***
		if len(s.daemons) > 0 ***REMOVED***
			tokens := s.daemons[0].JoinTokens(c)
			token := tokens.Worker
			if manager ***REMOVED***
				token = tokens.Manager
			***REMOVED***
			c.Assert(d.Join(swarm.JoinRequest***REMOVED***
				RemoteAddrs: []string***REMOVED***s.daemons[0].ListenAddr***REMOVED***,
				JoinToken:   token,
			***REMOVED***), check.IsNil)
		***REMOVED*** else ***REMOVED***
			c.Assert(d.Init(swarm.InitRequest***REMOVED******REMOVED***), check.IsNil)
		***REMOVED***
	***REMOVED***

	s.portIndex++
	s.daemonsLock.Lock()
	s.daemons = append(s.daemons, d)
	s.daemonsLock.Unlock()

	return d
***REMOVED***

func (s *DockerSwarmSuite) TearDownTest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	s.daemonsLock.Lock()
	for _, d := range s.daemons ***REMOVED***
		if d != nil ***REMOVED***
			d.Stop(c)
			// FIXME(vdemeester) should be handled by SwarmDaemon ?
			// raft state file is quite big (64MB) so remove it after every test
			walDir := filepath.Join(d.Root, "swarm/raft/wal")
			if err := os.RemoveAll(walDir); err != nil ***REMOVED***
				c.Logf("error removing %v: %v", walDir, err)
			***REMOVED***

			d.CleanupExecRoot(c)
		***REMOVED***
	***REMOVED***
	s.daemons = nil
	s.daemonsLock.Unlock()

	s.portIndex = 0
	s.ds.TearDownTest(c)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerTrustSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerTrustSuite struct ***REMOVED***
	ds  *DockerSuite
	reg *registry.V2
	not *testNotary
***REMOVED***

func (s *DockerTrustSuite) OnTimeout(c *check.C) ***REMOVED***
	s.ds.OnTimeout(c)
***REMOVED***

func (s *DockerTrustSuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, registry.Hosting, NotaryServerHosting)
	s.reg = setupRegistry(c, false, "", "")
	s.not = setupNotary(c)
***REMOVED***

func (s *DockerTrustSuite) TearDownTest(c *check.C) ***REMOVED***
	if s.reg != nil ***REMOVED***
		s.reg.Close()
	***REMOVED***
	if s.not != nil ***REMOVED***
		s.not.Close()
	***REMOVED***

	// Remove trusted keys and metadata after test
	os.RemoveAll(filepath.Join(config.Dir(), "trust"))
	s.ds.TearDownTest(c)
***REMOVED***

func init() ***REMOVED***
	ds := &DockerSuite***REMOVED******REMOVED***
	check.Suite(&DockerTrustedSwarmSuite***REMOVED***
		trustSuite: DockerTrustSuite***REMOVED***
			ds: ds,
		***REMOVED***,
		swarmSuite: DockerSwarmSuite***REMOVED***
			ds: ds,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerTrustedSwarmSuite struct ***REMOVED***
	swarmSuite DockerSwarmSuite
	trustSuite DockerTrustSuite
	reg        *registry.V2
	not        *testNotary
***REMOVED***

func (s *DockerTrustedSwarmSuite) SetUpTest(c *check.C) ***REMOVED***
	s.swarmSuite.SetUpTest(c)
	s.trustSuite.SetUpTest(c)
***REMOVED***

func (s *DockerTrustedSwarmSuite) TearDownTest(c *check.C) ***REMOVED***
	s.trustSuite.TearDownTest(c)
	s.swarmSuite.TearDownTest(c)
***REMOVED***

func (s *DockerTrustedSwarmSuite) OnTimeout(c *check.C) ***REMOVED***
	s.swarmSuite.OnTimeout(c)
***REMOVED***

func init() ***REMOVED***
	check.Suite(&DockerPluginSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type DockerPluginSuite struct ***REMOVED***
	ds       *DockerSuite
	registry *registry.V2
***REMOVED***

func (ps *DockerPluginSuite) registryHost() string ***REMOVED***
	return privateRegistryURL
***REMOVED***

func (ps *DockerPluginSuite) getPluginRepo() string ***REMOVED***
	return path.Join(ps.registryHost(), "plugin", "basic")
***REMOVED***
func (ps *DockerPluginSuite) getPluginRepoWithTag() string ***REMOVED***
	return ps.getPluginRepo() + ":" + "latest"
***REMOVED***

func (ps *DockerPluginSuite) SetUpSuite(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, registry.Hosting)
	ps.registry = setupRegistry(c, false, "", "")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := plugin.CreateInRegistry(ctx, ps.getPluginRepo(), nil)
	c.Assert(err, checker.IsNil, check.Commentf("failed to create plugin"))
***REMOVED***

func (ps *DockerPluginSuite) TearDownSuite(c *check.C) ***REMOVED***
	if ps.registry != nil ***REMOVED***
		ps.registry.Close()
	***REMOVED***
***REMOVED***

func (ps *DockerPluginSuite) TearDownTest(c *check.C) ***REMOVED***
	ps.ds.TearDownTest(c)
***REMOVED***

func (ps *DockerPluginSuite) OnTimeout(c *check.C) ***REMOVED***
	ps.ds.OnTimeout(c)
***REMOVED***
