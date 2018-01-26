package plugin

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/plugin"
	"github.com/docker/docker/registry"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// CreateOpt is is passed used to change the default plugin config before
// creating it
type CreateOpt func(*Config)

// Config wraps types.PluginConfig to provide some extra state for options
// extra customizations on the plugin details, such as using a custom binary to
// create the plugin with.
type Config struct ***REMOVED***
	*types.PluginConfig
	binPath string
***REMOVED***

// WithBinary is a CreateOpt to set an custom binary to create the plugin with.
// This binary must be statically compiled.
func WithBinary(bin string) CreateOpt ***REMOVED***
	return func(cfg *Config) ***REMOVED***
		cfg.binPath = bin
	***REMOVED***
***REMOVED***

// CreateClient is the interface used for `BuildPlugin` to interact with the
// daemon.
type CreateClient interface ***REMOVED***
	PluginCreate(context.Context, io.Reader, types.PluginCreateOptions) error
***REMOVED***

// Create creates a new plugin with the specified name
func Create(ctx context.Context, c CreateClient, name string, opts ...CreateOpt) error ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "create-test-plugin")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	tar, err := makePluginBundle(tmpDir, opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer tar.Close()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return c.PluginCreate(ctx, tar, types.PluginCreateOptions***REMOVED***RepoName: name***REMOVED***)
***REMOVED***

// CreateInRegistry makes a plugin (locally) and pushes it to a registry.
// This does not use a dockerd instance to create or push the plugin.
// If you just want to create a plugin in some daemon, use `Create`.
//
// This can be useful when testing plugins on swarm where you don't really want
// the plugin to exist on any of the daemons (immediately) and there needs to be
// some way to distribute the plugin.
func CreateInRegistry(ctx context.Context, repo string, auth *types.AuthConfig, opts ...CreateOpt) error ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "create-test-plugin-local")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	inPath := filepath.Join(tmpDir, "plugin")
	if err := os.MkdirAll(inPath, 0755); err != nil ***REMOVED***
		return errors.Wrap(err, "error creating plugin root")
	***REMOVED***

	tar, err := makePluginBundle(inPath, opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer tar.Close()

	dummyExec := func(m *plugin.Manager) (plugin.Executor, error) ***REMOVED***
		return nil, nil
	***REMOVED***

	regService, err := registry.NewService(registry.ServiceOptions***REMOVED***V2Only: true***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	managerConfig := plugin.ManagerConfig***REMOVED***
		Store:           plugin.NewStore(),
		RegistryService: regService,
		Root:            filepath.Join(tmpDir, "root"),
		ExecRoot:        "/run/docker", // manager init fails if not set
		CreateExecutor:  dummyExec,
		LogPluginEvent:  func(id, name, action string) ***REMOVED******REMOVED***, // panics when not set
	***REMOVED***
	manager, err := plugin.NewManager(managerConfig)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error creating plugin manager")
	***REMOVED***

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := manager.CreateFromContext(ctx, tar, &types.PluginCreateOptions***REMOVED***RepoName: repo***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	if auth == nil ***REMOVED***
		auth = &types.AuthConfig***REMOVED******REMOVED***
	***REMOVED***
	err = manager.Push(ctx, repo, nil, auth, ioutil.Discard)
	return errors.Wrap(err, "error pushing plugin")
***REMOVED***

func makePluginBundle(inPath string, opts ...CreateOpt) (io.ReadCloser, error) ***REMOVED***
	p := &types.PluginConfig***REMOVED***
		Interface: types.PluginConfigInterface***REMOVED***
			Socket: "basic.sock",
			Types:  []types.PluginInterfaceType***REMOVED******REMOVED***Capability: "docker.dummy/1.0"***REMOVED******REMOVED***,
		***REMOVED***,
		Entrypoint: []string***REMOVED***"/basic"***REMOVED***,
	***REMOVED***
	cfg := &Config***REMOVED***
		PluginConfig: p,
	***REMOVED***
	for _, o := range opts ***REMOVED***
		o(cfg)
	***REMOVED***
	if cfg.binPath == "" ***REMOVED***
		binPath, err := ensureBasicPluginBin()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cfg.binPath = binPath
	***REMOVED***

	configJSON, err := json.Marshal(p)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(inPath, "config.json"), configJSON, 0644); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := os.MkdirAll(filepath.Join(inPath, "rootfs", filepath.Dir(p.Entrypoint[0])), 0755); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error creating plugin rootfs dir")
	***REMOVED***
	if err := archive.NewDefaultArchiver().CopyFileWithTar(cfg.binPath, filepath.Join(inPath, "rootfs", p.Entrypoint[0])); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error copying plugin binary to rootfs path")
	***REMOVED***
	tar, err := archive.Tar(inPath, archive.Uncompressed)
	return tar, errors.Wrap(err, "error making plugin archive")
***REMOVED***

func ensureBasicPluginBin() (string, error) ***REMOVED***
	name := "docker-basic-plugin"
	p, err := exec.LookPath(name)
	if err == nil ***REMOVED***
		return p, nil
	***REMOVED***

	goBin, err := exec.LookPath("go")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	installPath := filepath.Join(os.Getenv("GOPATH"), "bin", name)
	cmd := exec.Command(goBin, "build", "-o", installPath, "./"+filepath.Join("fixtures", "plugin", "basic"))
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "error building basic plugin bin: %s", string(out))
	***REMOVED***
	return installPath, nil
***REMOVED***
