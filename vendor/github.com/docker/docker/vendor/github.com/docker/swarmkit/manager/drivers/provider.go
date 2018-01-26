package drivers

import (
	"fmt"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/swarmkit/api"
)

// DriverProvider provides external drivers
type DriverProvider struct ***REMOVED***
	pluginGetter plugingetter.PluginGetter
***REMOVED***

// New returns a new driver provider
func New(pluginGetter plugingetter.PluginGetter) *DriverProvider ***REMOVED***
	return &DriverProvider***REMOVED***pluginGetter: pluginGetter***REMOVED***
***REMOVED***

// NewSecretDriver creates a new driver for fetching secrets
func (m *DriverProvider) NewSecretDriver(driver *api.Driver) (*SecretDriver, error) ***REMOVED***
	if m.pluginGetter == nil ***REMOVED***
		return nil, fmt.Errorf("plugin getter is nil")
	***REMOVED***
	if driver == nil && driver.Name == "" ***REMOVED***
		return nil, fmt.Errorf("driver specification is nil")
	***REMOVED***
	// Search for the specified plugin
	plugin, err := m.pluginGetter.Get(driver.Name, SecretsProviderCapability, plugingetter.Lookup)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewSecretDriver(plugin), nil
***REMOVED***
