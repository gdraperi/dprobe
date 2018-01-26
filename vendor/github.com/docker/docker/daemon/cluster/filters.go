package cluster

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/filters"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	swarmapi "github.com/docker/swarmkit/api"
)

func newListNodesFilters(filter filters.Args) (*swarmapi.ListNodesRequest_Filters, error) ***REMOVED***
	accepted := map[string]bool***REMOVED***
		"name":       true,
		"id":         true,
		"label":      true,
		"role":       true,
		"membership": true,
	***REMOVED***
	if err := filter.Validate(accepted); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	f := &swarmapi.ListNodesRequest_Filters***REMOVED***
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
	***REMOVED***

	for _, r := range filter.Get("role") ***REMOVED***
		if role, ok := swarmapi.NodeRole_value[strings.ToUpper(r)]; ok ***REMOVED***
			f.Roles = append(f.Roles, swarmapi.NodeRole(role))
		***REMOVED*** else if r != "" ***REMOVED***
			return nil, fmt.Errorf("Invalid role filter: '%s'", r)
		***REMOVED***
	***REMOVED***

	for _, a := range filter.Get("membership") ***REMOVED***
		if membership, ok := swarmapi.NodeSpec_Membership_value[strings.ToUpper(a)]; ok ***REMOVED***
			f.Memberships = append(f.Memberships, swarmapi.NodeSpec_Membership(membership))
		***REMOVED*** else if a != "" ***REMOVED***
			return nil, fmt.Errorf("Invalid membership filter: '%s'", a)
		***REMOVED***
	***REMOVED***

	return f, nil
***REMOVED***

func newListTasksFilters(filter filters.Args, transformFunc func(filters.Args) error) (*swarmapi.ListTasksRequest_Filters, error) ***REMOVED***
	accepted := map[string]bool***REMOVED***
		"name":          true,
		"id":            true,
		"label":         true,
		"service":       true,
		"node":          true,
		"desired-state": true,
		// UpToDate is not meant to be exposed to users. It's for
		// internal use in checking create/update progress. Therefore,
		// we prefix it with a '_'.
		"_up-to-date": true,
		"runtime":     true,
	***REMOVED***
	if err := filter.Validate(accepted); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if transformFunc != nil ***REMOVED***
		if err := transformFunc(filter); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	f := &swarmapi.ListTasksRequest_Filters***REMOVED***
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
		ServiceIDs:   filter.Get("service"),
		NodeIDs:      filter.Get("node"),
		UpToDate:     len(filter.Get("_up-to-date")) != 0,
		Runtimes:     filter.Get("runtime"),
	***REMOVED***

	for _, s := range filter.Get("desired-state") ***REMOVED***
		if state, ok := swarmapi.TaskState_value[strings.ToUpper(s)]; ok ***REMOVED***
			f.DesiredStates = append(f.DesiredStates, swarmapi.TaskState(state))
		***REMOVED*** else if s != "" ***REMOVED***
			return nil, fmt.Errorf("Invalid desired-state filter: '%s'", s)
		***REMOVED***
	***REMOVED***

	return f, nil
***REMOVED***

func newListSecretsFilters(filter filters.Args) (*swarmapi.ListSecretsRequest_Filters, error) ***REMOVED***
	accepted := map[string]bool***REMOVED***
		"names": true,
		"name":  true,
		"id":    true,
		"label": true,
	***REMOVED***
	if err := filter.Validate(accepted); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &swarmapi.ListSecretsRequest_Filters***REMOVED***
		Names:        filter.Get("names"),
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
	***REMOVED***, nil
***REMOVED***

func newListConfigsFilters(filter filters.Args) (*swarmapi.ListConfigsRequest_Filters, error) ***REMOVED***
	accepted := map[string]bool***REMOVED***
		"name":  true,
		"id":    true,
		"label": true,
	***REMOVED***
	if err := filter.Validate(accepted); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &swarmapi.ListConfigsRequest_Filters***REMOVED***
		NamePrefixes: filter.Get("name"),
		IDPrefixes:   filter.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(filter.Get("label")),
	***REMOVED***, nil
***REMOVED***
