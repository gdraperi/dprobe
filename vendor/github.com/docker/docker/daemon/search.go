package daemon

import (
	"strconv"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/dockerversion"
)

var acceptedSearchFilterTags = map[string]bool***REMOVED***
	"is-automated": true,
	"is-official":  true,
	"stars":        true,
***REMOVED***

// SearchRegistryForImages queries the registry for images matching
// term. authConfig is used to login.
func (daemon *Daemon) SearchRegistryForImages(ctx context.Context, filtersArgs string, term string, limit int,
	authConfig *types.AuthConfig,
	headers map[string][]string) (*registrytypes.SearchResults, error) ***REMOVED***

	searchFilters, err := filters.FromJSON(filtersArgs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := searchFilters.Validate(acceptedSearchFilterTags); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var isAutomated, isOfficial bool
	var hasStarFilter = 0
	if searchFilters.Contains("is-automated") ***REMOVED***
		if searchFilters.UniqueExactMatch("is-automated", "true") ***REMOVED***
			isAutomated = true
		***REMOVED*** else if !searchFilters.UniqueExactMatch("is-automated", "false") ***REMOVED***
			return nil, invalidFilter***REMOVED***"is-automated", searchFilters.Get("is-automated")***REMOVED***
		***REMOVED***
	***REMOVED***
	if searchFilters.Contains("is-official") ***REMOVED***
		if searchFilters.UniqueExactMatch("is-official", "true") ***REMOVED***
			isOfficial = true
		***REMOVED*** else if !searchFilters.UniqueExactMatch("is-official", "false") ***REMOVED***
			return nil, invalidFilter***REMOVED***"is-official", searchFilters.Get("is-official")***REMOVED***
		***REMOVED***
	***REMOVED***
	if searchFilters.Contains("stars") ***REMOVED***
		hasStars := searchFilters.Get("stars")
		for _, hasStar := range hasStars ***REMOVED***
			iHasStar, err := strconv.Atoi(hasStar)
			if err != nil ***REMOVED***
				return nil, invalidFilter***REMOVED***"stars", hasStar***REMOVED***
			***REMOVED***
			if iHasStar > hasStarFilter ***REMOVED***
				hasStarFilter = iHasStar
			***REMOVED***
		***REMOVED***
	***REMOVED***

	unfilteredResult, err := daemon.RegistryService.Search(ctx, term, limit, authConfig, dockerversion.DockerUserAgent(ctx), headers)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	filteredResults := []registrytypes.SearchResult***REMOVED******REMOVED***
	for _, result := range unfilteredResult.Results ***REMOVED***
		if searchFilters.Contains("is-automated") ***REMOVED***
			if isAutomated != result.IsAutomated ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if searchFilters.Contains("is-official") ***REMOVED***
			if isOfficial != result.IsOfficial ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		if searchFilters.Contains("stars") ***REMOVED***
			if result.StarCount < hasStarFilter ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		filteredResults = append(filteredResults, result)
	***REMOVED***

	return &registrytypes.SearchResults***REMOVED***
		Query:      unfilteredResult.Query,
		NumResults: len(filteredResults),
		Results:    filteredResults,
	***REMOVED***, nil
***REMOVED***
