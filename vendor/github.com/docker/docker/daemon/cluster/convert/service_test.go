package convert

import (
	"testing"

	containertypes "github.com/docker/docker/api/types/container"
	swarmtypes "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/swarm/runtime"
	swarmapi "github.com/docker/swarmkit/api"
	google_protobuf3 "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestServiceConvertFromGRPCRuntimeContainer(t *testing.T) ***REMOVED***
	gs := swarmapi.Service***REMOVED***
		Meta: swarmapi.Meta***REMOVED***
			Version: swarmapi.Version***REMOVED***
				Index: 1,
			***REMOVED***,
			CreatedAt: nil,
			UpdatedAt: nil,
		***REMOVED***,
		SpecVersion: &swarmapi.Version***REMOVED***
			Index: 1,
		***REMOVED***,
		Spec: swarmapi.ServiceSpec***REMOVED***
			Task: swarmapi.TaskSpec***REMOVED***
				Runtime: &swarmapi.TaskSpec_Container***REMOVED***
					Container: &swarmapi.ContainerSpec***REMOVED***
						Image: "alpine:latest",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	svc, err := ServiceFromGRPC(gs)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if svc.Spec.TaskTemplate.Runtime != swarmtypes.RuntimeContainer ***REMOVED***
		t.Fatalf("expected type %s; received %T", swarmtypes.RuntimeContainer, svc.Spec.TaskTemplate.Runtime)
	***REMOVED***
***REMOVED***

func TestServiceConvertFromGRPCGenericRuntimePlugin(t *testing.T) ***REMOVED***
	kind := string(swarmtypes.RuntimePlugin)
	url := swarmtypes.RuntimeURLPlugin
	gs := swarmapi.Service***REMOVED***
		Meta: swarmapi.Meta***REMOVED***
			Version: swarmapi.Version***REMOVED***
				Index: 1,
			***REMOVED***,
			CreatedAt: nil,
			UpdatedAt: nil,
		***REMOVED***,
		SpecVersion: &swarmapi.Version***REMOVED***
			Index: 1,
		***REMOVED***,
		Spec: swarmapi.ServiceSpec***REMOVED***
			Task: swarmapi.TaskSpec***REMOVED***
				Runtime: &swarmapi.TaskSpec_Generic***REMOVED***
					Generic: &swarmapi.GenericRuntimeSpec***REMOVED***
						Kind: kind,
						Payload: &google_protobuf3.Any***REMOVED***
							TypeUrl: string(url),
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	svc, err := ServiceFromGRPC(gs)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if svc.Spec.TaskTemplate.Runtime != swarmtypes.RuntimePlugin ***REMOVED***
		t.Fatalf("expected type %s; received %T", swarmtypes.RuntimePlugin, svc.Spec.TaskTemplate.Runtime)
	***REMOVED***
***REMOVED***

func TestServiceConvertToGRPCGenericRuntimePlugin(t *testing.T) ***REMOVED***
	s := swarmtypes.ServiceSpec***REMOVED***
		TaskTemplate: swarmtypes.TaskSpec***REMOVED***
			Runtime:    swarmtypes.RuntimePlugin,
			PluginSpec: &runtime.PluginSpec***REMOVED******REMOVED***,
		***REMOVED***,
		Mode: swarmtypes.ServiceMode***REMOVED***
			Global: &swarmtypes.GlobalService***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***

	svc, err := ServiceSpecToGRPC(s)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	v, ok := svc.Task.Runtime.(*swarmapi.TaskSpec_Generic)
	if !ok ***REMOVED***
		t.Fatal("expected type swarmapi.TaskSpec_Generic")
	***REMOVED***

	if v.Generic.Payload.TypeUrl != string(swarmtypes.RuntimeURLPlugin) ***REMOVED***
		t.Fatalf("expected url %s; received %s", swarmtypes.RuntimeURLPlugin, v.Generic.Payload.TypeUrl)
	***REMOVED***
***REMOVED***

func TestServiceConvertToGRPCContainerRuntime(t *testing.T) ***REMOVED***
	image := "alpine:latest"
	s := swarmtypes.ServiceSpec***REMOVED***
		TaskTemplate: swarmtypes.TaskSpec***REMOVED***
			ContainerSpec: &swarmtypes.ContainerSpec***REMOVED***
				Image: image,
			***REMOVED***,
		***REMOVED***,
		Mode: swarmtypes.ServiceMode***REMOVED***
			Global: &swarmtypes.GlobalService***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***

	svc, err := ServiceSpecToGRPC(s)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	v, ok := svc.Task.Runtime.(*swarmapi.TaskSpec_Container)
	if !ok ***REMOVED***
		t.Fatal("expected type swarmapi.TaskSpec_Container")
	***REMOVED***

	if v.Container.Image != image ***REMOVED***
		t.Fatalf("expected image %s; received %s", image, v.Container.Image)
	***REMOVED***
***REMOVED***

func TestServiceConvertToGRPCGenericRuntimeCustom(t *testing.T) ***REMOVED***
	s := swarmtypes.ServiceSpec***REMOVED***
		TaskTemplate: swarmtypes.TaskSpec***REMOVED***
			Runtime: "customruntime",
		***REMOVED***,
		Mode: swarmtypes.ServiceMode***REMOVED***
			Global: &swarmtypes.GlobalService***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***

	if _, err := ServiceSpecToGRPC(s); err != ErrUnsupportedRuntime ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestServiceConvertToGRPCIsolation(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		name string
		from containertypes.Isolation
		to   swarmapi.ContainerSpec_Isolation
	***REMOVED******REMOVED***
		***REMOVED***name: "empty", from: containertypes.IsolationEmpty, to: swarmapi.ContainerIsolationDefault***REMOVED***,
		***REMOVED***name: "default", from: containertypes.IsolationDefault, to: swarmapi.ContainerIsolationDefault***REMOVED***,
		***REMOVED***name: "process", from: containertypes.IsolationProcess, to: swarmapi.ContainerIsolationProcess***REMOVED***,
		***REMOVED***name: "hyperv", from: containertypes.IsolationHyperV, to: swarmapi.ContainerIsolationHyperV***REMOVED***,
		***REMOVED***name: "proCess", from: containertypes.Isolation("proCess"), to: swarmapi.ContainerIsolationProcess***REMOVED***,
		***REMOVED***name: "hypErv", from: containertypes.Isolation("hypErv"), to: swarmapi.ContainerIsolationHyperV***REMOVED***,
	***REMOVED***
	for _, c := range cases ***REMOVED***
		t.Run(c.name, func(t *testing.T) ***REMOVED***
			s := swarmtypes.ServiceSpec***REMOVED***
				TaskTemplate: swarmtypes.TaskSpec***REMOVED***
					ContainerSpec: &swarmtypes.ContainerSpec***REMOVED***
						Image:     "alpine:latest",
						Isolation: c.from,
					***REMOVED***,
				***REMOVED***,
				Mode: swarmtypes.ServiceMode***REMOVED***
					Global: &swarmtypes.GlobalService***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED***
			res, err := ServiceSpecToGRPC(s)
			require.NoError(t, err)
			v, ok := res.Task.Runtime.(*swarmapi.TaskSpec_Container)
			if !ok ***REMOVED***
				t.Fatal("expected type swarmapi.TaskSpec_Container")
			***REMOVED***
			require.Equal(t, c.to, v.Container.Isolation)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestServiceConvertFromGRPCIsolation(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		name string
		from swarmapi.ContainerSpec_Isolation
		to   containertypes.Isolation
	***REMOVED******REMOVED***
		***REMOVED***name: "default", to: containertypes.IsolationDefault, from: swarmapi.ContainerIsolationDefault***REMOVED***,
		***REMOVED***name: "process", to: containertypes.IsolationProcess, from: swarmapi.ContainerIsolationProcess***REMOVED***,
		***REMOVED***name: "hyperv", to: containertypes.IsolationHyperV, from: swarmapi.ContainerIsolationHyperV***REMOVED***,
	***REMOVED***
	for _, c := range cases ***REMOVED***
		t.Run(c.name, func(t *testing.T) ***REMOVED***
			gs := swarmapi.Service***REMOVED***
				Meta: swarmapi.Meta***REMOVED***
					Version: swarmapi.Version***REMOVED***
						Index: 1,
					***REMOVED***,
					CreatedAt: nil,
					UpdatedAt: nil,
				***REMOVED***,
				SpecVersion: &swarmapi.Version***REMOVED***
					Index: 1,
				***REMOVED***,
				Spec: swarmapi.ServiceSpec***REMOVED***
					Task: swarmapi.TaskSpec***REMOVED***
						Runtime: &swarmapi.TaskSpec_Container***REMOVED***
							Container: &swarmapi.ContainerSpec***REMOVED***
								Image:     "alpine:latest",
								Isolation: c.from,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***

			svc, err := ServiceFromGRPC(gs)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***

			require.Equal(t, c.to, svc.Spec.TaskTemplate.ContainerSpec.Isolation)
		***REMOVED***)
	***REMOVED***
***REMOVED***
