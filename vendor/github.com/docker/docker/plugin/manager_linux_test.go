package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/plugin/v2"
)

func TestManagerWithPluginMounts(t *testing.T) ***REMOVED***
	root, err := ioutil.TempDir("", "test-store-with-plugin-mounts")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer system.EnsureRemoveAll(root)

	s := NewStore()
	managerRoot := filepath.Join(root, "manager")
	p1 := newTestPlugin(t, "test1", "testcap", managerRoot)

	p2 := newTestPlugin(t, "test2", "testcap", managerRoot)
	p2.PluginObj.Enabled = true

	m, err := NewManager(
		ManagerConfig***REMOVED***
			Store:          s,
			Root:           managerRoot,
			ExecRoot:       filepath.Join(root, "exec"),
			CreateExecutor: func(*Manager) (Executor, error) ***REMOVED*** return nil, nil ***REMOVED***,
			LogPluginEvent: func(_, _, _ string) ***REMOVED******REMOVED***,
		***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := s.Add(p1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := s.Add(p2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Create a mount to simulate a plugin that has created it's own mounts
	p2Mount := filepath.Join(p2.Rootfs, "testmount")
	if err := os.MkdirAll(p2Mount, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := mount.Mount("tmpfs", p2Mount, "tmpfs", ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := m.Remove(p1.Name(), &types.PluginRmConfig***REMOVED***ForceRemove: true***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if mounted, err := mount.Mounted(p2Mount); !mounted || err != nil ***REMOVED***
		t.Fatalf("expected %s to be mounted, err: %v", p2Mount, err)
	***REMOVED***
***REMOVED***

func newTestPlugin(t *testing.T, name, cap, root string) *v2.Plugin ***REMOVED***
	rootfs := filepath.Join(root, name)
	if err := os.MkdirAll(rootfs, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	p := v2.Plugin***REMOVED***PluginObj: types.Plugin***REMOVED***Name: name***REMOVED******REMOVED***
	p.Rootfs = rootfs
	iType := types.PluginInterfaceType***REMOVED***Capability: cap, Prefix: "docker", Version: "1.0"***REMOVED***
	i := types.PluginConfigInterface***REMOVED***Socket: "plugins.sock", Types: []types.PluginInterfaceType***REMOVED***iType***REMOVED******REMOVED***
	p.PluginObj.Config.Interface = i
	p.PluginObj.ID = name

	return &p
***REMOVED***
