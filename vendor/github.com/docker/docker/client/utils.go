package client

import (
	"net/url"
	"regexp"

	"github.com/docker/docker/api/types/filters"
)

var headerRegexp = regexp.MustCompile(`\ADocker/.+\s\((.+)\)\z`)

// getDockerOS returns the operating system based on the server header from the daemon.
func getDockerOS(serverHeader string) string ***REMOVED***
	var osType string
	matches := headerRegexp.FindStringSubmatch(serverHeader)
	if len(matches) > 0 ***REMOVED***
		osType = matches[1]
	***REMOVED***
	return osType
***REMOVED***

// getFiltersQuery returns a url query with "filters" query term, based on the
// filters provided.
func getFiltersQuery(f filters.Args) (url.Values, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	if f.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToJSON(f)
		if err != nil ***REMOVED***
			return query, err
		***REMOVED***
		query.Set("filters", filterJSON)
	***REMOVED***
	return query, nil
***REMOVED***
