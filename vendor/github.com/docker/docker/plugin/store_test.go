package plugin

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/plugin/v2"
)

func TestFilterByCapNeg(t *testing.T) ***REMOVED***
	p := v2.Plugin***REMOVED***PluginObj: types.Plugin***REMOVED***Name: "test:latest"***REMOVED******REMOVED***
	iType := types.PluginInterfaceType***REMOVED***Capability: "volumedriver", Prefix: "docker", Version: "1.0"***REMOVED***
	i := types.PluginConfigInterface***REMOVED***Socket: "plugins.sock", Types: []types.PluginInterfaceType***REMOVED***iType***REMOVED******REMOVED***
	p.PluginObj.Config.Interface = i

	_, err := p.FilterByCap("foobar")
	if err == nil ***REMOVED***
		t.Fatalf("expected inadequate error, got %v", err)
	***REMOVED***
***REMOVED***

func TestFilterByCapPos(t *testing.T) ***REMOVED***
	p := v2.Plugin***REMOVED***PluginObj: types.Plugin***REMOVED***Name: "test:latest"***REMOVED******REMOVED***

	iType := types.PluginInterfaceType***REMOVED***Capability: "volumedriver", Prefix: "docker", Version: "1.0"***REMOVED***
	i := types.PluginConfigInterface***REMOVED***Socket: "plugins.sock", Types: []types.PluginInterfaceType***REMOVED***iType***REMOVED******REMOVED***
	p.PluginObj.Config.Interface = i

	_, err := p.FilterByCap("volumedriver")
	if err != nil ***REMOVED***
		t.Fatalf("expected no error, got %v", err)
	***REMOVED***
***REMOVED***

func TestStoreGetPluginNotMatchCapRefs(t *testing.T) ***REMOVED***
	s := NewStore()
	p := v2.Plugin***REMOVED***PluginObj: types.Plugin***REMOVED***Name: "test:latest"***REMOVED******REMOVED***

	iType := types.PluginInterfaceType***REMOVED***Capability: "whatever", Prefix: "docker", Version: "1.0"***REMOVED***
	i := types.PluginConfigInterface***REMOVED***Socket: "plugins.sock", Types: []types.PluginInterfaceType***REMOVED***iType***REMOVED******REMOVED***
	p.PluginObj.Config.Interface = i

	if err := s.Add(&p); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := s.Get("test", "volumedriver", plugingetter.Acquire); err == nil ***REMOVED***
		t.Fatal("exepcted error when getting plugin that doesn't match the passed in capability")
	***REMOVED***

	if refs := p.GetRefCount(); refs != 0 ***REMOVED***
		t.Fatalf("reference count should be 0, got: %d", refs)
	***REMOVED***

	p.PluginObj.Enabled = true
	if _, err := s.Get("test", "volumedriver", plugingetter.Acquire); err == nil ***REMOVED***
		t.Fatal("exepcted error when getting plugin that doesn't match the passed in capability")
	***REMOVED***

	if refs := p.GetRefCount(); refs != 0 ***REMOVED***
		t.Fatalf("reference count should be 0, got: %d", refs)
	***REMOVED***
***REMOVED***
