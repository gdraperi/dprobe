package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/docker/swarmkit/agent/configs"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/agent/secrets"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/naming"
	"github.com/pkg/errors"
)

// Platform holds information about the underlying platform of the node
type Platform struct ***REMOVED***
	Architecture string
	OS           string
***REMOVED***

// Context defines the strict set of values that can be injected into a
// template expression in SwarmKit data structure.
// NOTE: Be very careful adding any fields to this structure with types
// that have methods defined on them. The template would be able to
// invoke those methods.
type Context struct ***REMOVED***
	Service struct ***REMOVED***
		ID     string
		Name   string
		Labels map[string]string
	***REMOVED***

	Node struct ***REMOVED***
		ID       string
		Hostname string
		Platform Platform
	***REMOVED***

	Task struct ***REMOVED***
		ID   string
		Name string
		Slot string

		// NOTE(stevvooe): Why no labels here? Tasks don't actually have labels
		// (from a user perspective). The labels are part of the container! If
		// one wants to use labels for templating, use service labels!
	***REMOVED***
***REMOVED***

// NewContext returns a new template context from the data available in the
// task and the node where it is scheduled to run.
// The provided context can then be used to populate runtime values in a
// ContainerSpec.
func NewContext(n *api.NodeDescription, t *api.Task) (ctx Context) ***REMOVED***
	ctx.Service.ID = t.ServiceID
	ctx.Service.Name = t.ServiceAnnotations.Name
	ctx.Service.Labels = t.ServiceAnnotations.Labels

	ctx.Node.ID = t.NodeID

	// Add node information to context only if we have them available
	if n != nil ***REMOVED***
		ctx.Node.Hostname = n.Hostname
		ctx.Node.Platform = Platform***REMOVED***
			Architecture: n.Platform.Architecture,
			OS:           n.Platform.OS,
		***REMOVED***
	***REMOVED***
	ctx.Task.ID = t.ID
	ctx.Task.Name = naming.Task(t)

	if t.Slot != 0 ***REMOVED***
		ctx.Task.Slot = fmt.Sprint(t.Slot)
	***REMOVED*** else ***REMOVED***
		// fall back to node id for slot when there is no slot
		ctx.Task.Slot = t.NodeID
	***REMOVED***

	return
***REMOVED***

// Expand treats the string s as a template and populates it with values from
// the context.
func (ctx *Context) Expand(s string) (string, error) ***REMOVED***
	tmpl, err := newTemplate(s, nil)
	if err != nil ***REMOVED***
		return s, err
	***REMOVED***

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil ***REMOVED***
		return s, err
	***REMOVED***

	return buf.String(), nil
***REMOVED***

// PayloadContext provides a context for expanding a config or secret payload.
// NOTE: Be very careful adding any fields to this structure with types
// that have methods defined on them. The template would be able to
// invoke those methods.
type PayloadContext struct ***REMOVED***
	Context

	t                 *api.Task
	restrictedSecrets exec.SecretGetter
	restrictedConfigs exec.ConfigGetter
	sensitive         bool
***REMOVED***

func (ctx *PayloadContext) secretGetter(target string) (string, error) ***REMOVED***
	if ctx.restrictedSecrets == nil ***REMOVED***
		return "", errors.New("secrets unavailable")
	***REMOVED***

	container := ctx.t.Spec.GetContainer()
	if container == nil ***REMOVED***
		return "", errors.New("task is not a container")
	***REMOVED***

	for _, secretRef := range container.Secrets ***REMOVED***
		file := secretRef.GetFile()
		if file != nil && file.Name == target ***REMOVED***
			secret, err := ctx.restrictedSecrets.Get(secretRef.SecretID)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			ctx.sensitive = true
			return string(secret.Spec.Data), nil
		***REMOVED***
	***REMOVED***

	return "", errors.Errorf("secret target %s not found", target)
***REMOVED***

func (ctx *PayloadContext) configGetter(target string) (string, error) ***REMOVED***
	if ctx.restrictedConfigs == nil ***REMOVED***
		return "", errors.New("configs unavailable")
	***REMOVED***

	container := ctx.t.Spec.GetContainer()
	if container == nil ***REMOVED***
		return "", errors.New("task is not a container")
	***REMOVED***

	for _, configRef := range container.Configs ***REMOVED***
		file := configRef.GetFile()
		if file != nil && file.Name == target ***REMOVED***
			config, err := ctx.restrictedConfigs.Get(configRef.ConfigID)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			return string(config.Spec.Data), nil
		***REMOVED***
	***REMOVED***

	return "", errors.Errorf("config target %s not found", target)
***REMOVED***

func (ctx *PayloadContext) envGetter(variable string) (string, error) ***REMOVED***
	container := ctx.t.Spec.GetContainer()
	if container == nil ***REMOVED***
		return "", errors.New("task is not a container")
	***REMOVED***

	for _, env := range container.Env ***REMOVED***
		parts := strings.SplitN(env, "=", 2)

		if len(parts) > 1 && parts[0] == variable ***REMOVED***
			return parts[1], nil
		***REMOVED***
	***REMOVED***
	return "", nil
***REMOVED***

// NewPayloadContextFromTask returns a new template context from the data
// available in the task and the node where it is scheduled to run.
// This context also provides access to the configs
// and secrets that the task has access to. The provided context can then
// be used to populate runtime values in a templated config or secret.
func NewPayloadContextFromTask(node *api.NodeDescription, t *api.Task, dependencies exec.DependencyGetter) (ctx PayloadContext) ***REMOVED***
	return PayloadContext***REMOVED***
		Context:           NewContext(node, t),
		t:                 t,
		restrictedSecrets: secrets.Restrict(dependencies.Secrets(), t),
		restrictedConfigs: configs.Restrict(dependencies.Configs(), t),
	***REMOVED***
***REMOVED***

// Expand treats the string s as a template and populates it with values from
// the context.
func (ctx *PayloadContext) Expand(s string) (string, error) ***REMOVED***
	funcMap := template.FuncMap***REMOVED***
		"secret": ctx.secretGetter,
		"config": ctx.configGetter,
		"env":    ctx.envGetter,
	***REMOVED***

	tmpl, err := newTemplate(s, funcMap)
	if err != nil ***REMOVED***
		return s, err
	***REMOVED***

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil ***REMOVED***
		return s, err
	***REMOVED***

	return buf.String(), nil
***REMOVED***
