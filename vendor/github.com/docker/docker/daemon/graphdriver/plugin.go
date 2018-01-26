package graphdriver

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/plugin/v2"
)

func lookupPlugin(name string, pg plugingetter.PluginGetter, config Options) (Driver, error) ***REMOVED***
	if !config.ExperimentalEnabled ***REMOVED***
		return nil, fmt.Errorf("graphdriver plugins are only supported with experimental mode")
	***REMOVED***
	pl, err := pg.Get(name, "GraphDriver", plugingetter.Acquire)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Error looking up graphdriver plugin %s: %v", name, err)
	***REMOVED***
	return newPluginDriver(name, pl, config)
***REMOVED***

func newPluginDriver(name string, pl plugingetter.CompatPlugin, config Options) (Driver, error) ***REMOVED***
	home := config.Root
	if !pl.IsV1() ***REMOVED***
		if p, ok := pl.(*v2.Plugin); ok ***REMOVED***
			if p.PropagatedMount != "" ***REMOVED***
				home = p.PluginObj.Config.PropagatedMount
			***REMOVED***
		***REMOVED***
	***REMOVED***
	proxy := &graphDriverProxy***REMOVED***name, pl, Capabilities***REMOVED******REMOVED******REMOVED***
	return proxy, proxy.Init(filepath.Join(home, name), config.DriverOptions, config.UIDMaps, config.GIDMaps)
***REMOVED***
