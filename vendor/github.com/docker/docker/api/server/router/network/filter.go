package network

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/runconfig"
)

func filterNetworkByType(nws []types.NetworkResource, netType string) ([]types.NetworkResource, error) ***REMOVED***
	retNws := []types.NetworkResource***REMOVED******REMOVED***
	switch netType ***REMOVED***
	case "builtin":
		for _, nw := range nws ***REMOVED***
			if runconfig.IsPreDefinedNetwork(nw.Name) ***REMOVED***
				retNws = append(retNws, nw)
			***REMOVED***
		***REMOVED***
	case "custom":
		for _, nw := range nws ***REMOVED***
			if !runconfig.IsPreDefinedNetwork(nw.Name) ***REMOVED***
				retNws = append(retNws, nw)
			***REMOVED***
		***REMOVED***
	default:
		return nil, invalidFilter(netType)
	***REMOVED***
	return retNws, nil
***REMOVED***

type invalidFilter string

func (e invalidFilter) Error() string ***REMOVED***
	return "Invalid filter: 'type'='" + string(e) + "'"
***REMOVED***

func (e invalidFilter) InvalidParameter() ***REMOVED******REMOVED***

// filterNetworks filters network list according to user specified filter
// and returns user chosen networks
func filterNetworks(nws []types.NetworkResource, filter filters.Args) ([]types.NetworkResource, error) ***REMOVED***
	// if filter is empty, return original network list
	if filter.Len() == 0 ***REMOVED***
		return nws, nil
	***REMOVED***

	displayNet := []types.NetworkResource***REMOVED******REMOVED***
	for _, nw := range nws ***REMOVED***
		if filter.Contains("driver") ***REMOVED***
			if !filter.ExactMatch("driver", nw.Driver) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if filter.Contains("name") ***REMOVED***
			if !filter.Match("name", nw.Name) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if filter.Contains("id") ***REMOVED***
			if !filter.Match("id", nw.ID) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if filter.Contains("label") ***REMOVED***
			if !filter.MatchKVList("label", nw.Labels) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if filter.Contains("scope") ***REMOVED***
			if !filter.ExactMatch("scope", nw.Scope) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		displayNet = append(displayNet, nw)
	***REMOVED***

	if filter.Contains("type") ***REMOVED***
		typeNet := []types.NetworkResource***REMOVED******REMOVED***
		errFilter := filter.WalkValues("type", func(fval string) error ***REMOVED***
			passList, err := filterNetworkByType(displayNet, fval)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			typeNet = append(typeNet, passList...)
			return nil
		***REMOVED***)
		if errFilter != nil ***REMOVED***
			return nil, errFilter
		***REMOVED***
		displayNet = typeNet
	***REMOVED***

	return displayNet, nil
***REMOVED***
