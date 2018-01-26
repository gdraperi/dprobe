package daemon

import (
	"github.com/Microsoft/opengcs/client"
	"github.com/docker/docker/container"
)

func (daemon *Daemon) getLibcontainerdCreateOptions(container *container.Container) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// LCOW options.
	if container.OS == "linux" ***REMOVED***
		config := &client.Config***REMOVED******REMOVED***
		if err := config.GenerateDefault(daemon.configStore.GraphOptions); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// Override from user-supplied options.
		for k, v := range container.HostConfig.StorageOpt ***REMOVED***
			switch k ***REMOVED***
			case "lcow.kirdpath":
				config.KirdPath = v
			case "lcow.kernel":
				config.KernelFile = v
			case "lcow.initrd":
				config.InitrdFile = v
			case "lcow.vhdx":
				config.Vhdx = v
			case "lcow.bootparameters":
				config.BootParameters = v
			***REMOVED***
		***REMOVED***
		if err := config.Validate(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return config, nil
	***REMOVED***

	return nil, nil
***REMOVED***
