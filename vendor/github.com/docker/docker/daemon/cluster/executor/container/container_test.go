package container

import (
	"testing"

	container "github.com/docker/docker/api/types/container"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/stretchr/testify/require"
)

func TestIsolationConversion(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		name string
		from swarmapi.ContainerSpec_Isolation
		to   container.Isolation
	***REMOVED******REMOVED***
		***REMOVED***name: "default", from: swarmapi.ContainerIsolationDefault, to: container.IsolationDefault***REMOVED***,
		***REMOVED***name: "process", from: swarmapi.ContainerIsolationProcess, to: container.IsolationProcess***REMOVED***,
		***REMOVED***name: "hyperv", from: swarmapi.ContainerIsolationHyperV, to: container.IsolationHyperV***REMOVED***,
	***REMOVED***
	for _, c := range cases ***REMOVED***
		t.Run(c.name, func(t *testing.T) ***REMOVED***
			task := swarmapi.Task***REMOVED***
				Spec: swarmapi.TaskSpec***REMOVED***
					Runtime: &swarmapi.TaskSpec_Container***REMOVED***
						Container: &swarmapi.ContainerSpec***REMOVED***
							Image:     "alpine:latest",
							Isolation: c.from,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***
			config := containerConfig***REMOVED***task: &task***REMOVED***
			require.Equal(t, c.to, config.hostConfig().Isolation)
		***REMOVED***)
	***REMOVED***
***REMOVED***
