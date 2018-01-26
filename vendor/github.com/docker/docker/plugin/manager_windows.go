package plugin

import (
	"fmt"

	"github.com/docker/docker/plugin/v2"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func (pm *Manager) enable(p *v2.Plugin, c *controller, force bool) error ***REMOVED***
	return fmt.Errorf("Not implemented")
***REMOVED***

func (pm *Manager) initSpec(p *v2.Plugin) (*specs.Spec, error) ***REMOVED***
	return nil, fmt.Errorf("Not implemented")
***REMOVED***

func (pm *Manager) disable(p *v2.Plugin, c *controller) error ***REMOVED***
	return fmt.Errorf("Not implemented")
***REMOVED***

func (pm *Manager) restore(p *v2.Plugin) error ***REMOVED***
	return fmt.Errorf("Not implemented")
***REMOVED***

// Shutdown plugins
func (pm *Manager) Shutdown() ***REMOVED***
***REMOVED***

func setupRoot(root string) error ***REMOVED*** return nil ***REMOVED***
