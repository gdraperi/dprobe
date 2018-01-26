package template

import (
	"fmt"
	"strings"

	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/pkg/errors"
)

// ExpandContainerSpec expands templated fields in the runtime using the task
// state and the node where it is scheduled to run.
// Templating is all evaluated on the agent-side, before execution.
//
// Note that these are projected only on runtime values, since active task
// values are typically manipulated in the manager.
func ExpandContainerSpec(n *api.NodeDescription, t *api.Task) (*api.ContainerSpec, error) ***REMOVED***
	container := t.Spec.GetContainer()
	if container == nil ***REMOVED***
		return nil, errors.Errorf("task missing ContainerSpec to expand")
	***REMOVED***

	container = container.Copy()
	ctx := NewContext(n, t)

	var err error
	container.Env, err = expandEnv(ctx, container.Env)
	if err != nil ***REMOVED***
		return container, errors.Wrap(err, "expanding env failed")
	***REMOVED***

	// For now, we only allow templating of string-based mount fields
	container.Mounts, err = expandMounts(ctx, container.Mounts)
	if err != nil ***REMOVED***
		return container, errors.Wrap(err, "expanding mounts failed")
	***REMOVED***

	container.Hostname, err = ctx.Expand(container.Hostname)
	return container, errors.Wrap(err, "expanding hostname failed")
***REMOVED***

func expandMounts(ctx Context, mounts []api.Mount) ([]api.Mount, error) ***REMOVED***
	if len(mounts) == 0 ***REMOVED***
		return mounts, nil
	***REMOVED***

	expanded := make([]api.Mount, len(mounts))
	for i, mount := range mounts ***REMOVED***
		var err error
		mount.Source, err = ctx.Expand(mount.Source)
		if err != nil ***REMOVED***
			return mounts, errors.Wrapf(err, "expanding mount source %q", mount.Source)
		***REMOVED***

		mount.Target, err = ctx.Expand(mount.Target)
		if err != nil ***REMOVED***
			return mounts, errors.Wrapf(err, "expanding mount target %q", mount.Target)
		***REMOVED***

		if mount.VolumeOptions != nil ***REMOVED***
			mount.VolumeOptions.Labels, err = expandMap(ctx, mount.VolumeOptions.Labels)
			if err != nil ***REMOVED***
				return mounts, errors.Wrap(err, "expanding volume labels")
			***REMOVED***

			if mount.VolumeOptions.DriverConfig != nil ***REMOVED***
				mount.VolumeOptions.DriverConfig.Options, err = expandMap(ctx, mount.VolumeOptions.DriverConfig.Options)
				if err != nil ***REMOVED***
					return mounts, errors.Wrap(err, "expanding volume driver config")
				***REMOVED***
			***REMOVED***
		***REMOVED***

		expanded[i] = mount
	***REMOVED***

	return expanded, nil
***REMOVED***

func expandMap(ctx Context, m map[string]string) (map[string]string, error) ***REMOVED***
	var (
		n   = make(map[string]string, len(m))
		err error
	)

	for k, v := range m ***REMOVED***
		v, err = ctx.Expand(v)
		if err != nil ***REMOVED***
			return m, errors.Wrapf(err, "expanding map entry %q=%q", k, v)
		***REMOVED***

		n[k] = v
	***REMOVED***

	return n, nil
***REMOVED***

func expandEnv(ctx Context, values []string) ([]string, error) ***REMOVED***
	var result []string
	for _, value := range values ***REMOVED***
		var (
			parts = strings.SplitN(value, "=", 2)
			entry = parts[0]
		)

		if len(parts) > 1 ***REMOVED***
			expanded, err := ctx.Expand(parts[1])
			if err != nil ***REMOVED***
				return values, errors.Wrapf(err, "expanding env %q", value)
			***REMOVED***

			entry = fmt.Sprintf("%s=%s", entry, expanded)
		***REMOVED***

		result = append(result, entry)
	***REMOVED***

	return result, nil
***REMOVED***

func expandPayload(ctx *PayloadContext, payload []byte) ([]byte, error) ***REMOVED***
	result, err := ctx.Expand(string(payload))
	if err != nil ***REMOVED***
		return payload, err
	***REMOVED***
	return []byte(result), nil
***REMOVED***

// ExpandSecretSpec expands the template inside the secret payload, if any.
// Templating is evaluated on the agent-side.
func ExpandSecretSpec(s *api.Secret, node *api.NodeDescription, t *api.Task, dependencies exec.DependencyGetter) (*api.SecretSpec, error) ***REMOVED***
	if s.Spec.Templating == nil ***REMOVED***
		return &s.Spec, nil
	***REMOVED***
	if s.Spec.Templating.Name == "golang" ***REMOVED***
		ctx := NewPayloadContextFromTask(node, t, dependencies)
		secretSpec := s.Spec.Copy()

		var err error
		secretSpec.Data, err = expandPayload(&ctx, secretSpec.Data)
		return secretSpec, err
	***REMOVED***
	return &s.Spec, errors.New("unrecognized template type")
***REMOVED***

// ExpandConfigSpec expands the template inside the config payload, if any.
// Templating is evaluated on the agent-side.
func ExpandConfigSpec(c *api.Config, node *api.NodeDescription, t *api.Task, dependencies exec.DependencyGetter) (*api.ConfigSpec, bool, error) ***REMOVED***
	if c.Spec.Templating == nil ***REMOVED***
		return &c.Spec, false, nil
	***REMOVED***
	if c.Spec.Templating.Name == "golang" ***REMOVED***
		ctx := NewPayloadContextFromTask(node, t, dependencies)
		configSpec := c.Spec.Copy()

		var err error
		configSpec.Data, err = expandPayload(&ctx, configSpec.Data)
		return configSpec, ctx.sensitive, err
	***REMOVED***
	return &c.Spec, false, errors.New("unrecognized template type")
***REMOVED***
