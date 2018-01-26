package genericresource

import (
	"strconv"
	"strings"

	"github.com/docker/swarmkit/api"
)

func discreteToString(d *api.GenericResource_DiscreteResourceSpec) string ***REMOVED***
	return strconv.FormatInt(d.DiscreteResourceSpec.Value, 10)
***REMOVED***

// Kind returns the kind key as a string
func Kind(res *api.GenericResource) string ***REMOVED***
	switch r := res.Resource.(type) ***REMOVED***
	case *api.GenericResource_DiscreteResourceSpec:
		return r.DiscreteResourceSpec.Kind
	case *api.GenericResource_NamedResourceSpec:
		return r.NamedResourceSpec.Kind
	***REMOVED***

	return ""
***REMOVED***

// Value returns the value key as a string
func Value(res *api.GenericResource) string ***REMOVED***
	switch res := res.Resource.(type) ***REMOVED***
	case *api.GenericResource_DiscreteResourceSpec:
		return discreteToString(res)
	case *api.GenericResource_NamedResourceSpec:
		return res.NamedResourceSpec.Value
	***REMOVED***

	return ""
***REMOVED***

// EnvFormat returns the environment string version of the resource
func EnvFormat(res []*api.GenericResource, prefix string) []string ***REMOVED***
	envs := make(map[string][]string)
	for _, v := range res ***REMOVED***
		key := Kind(v)
		val := Value(v)
		envs[key] = append(envs[key], val)
	***REMOVED***

	env := make([]string, 0, len(res))
	for k, v := range envs ***REMOVED***
		k = strings.ToUpper(prefix + "_" + k)
		env = append(env, k+"="+strings.Join(v, ","))
	***REMOVED***

	return env
***REMOVED***
