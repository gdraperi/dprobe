package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"golang.org/x/net/context"
)

// ImageSearch makes the docker host to search by a term in a remote registry.
// The list of results is not sorted in any fashion.
func (cli *Client) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) ***REMOVED***
	var results []registry.SearchResult
	query := url.Values***REMOVED******REMOVED***
	query.Set("term", term)
	query.Set("limit", fmt.Sprintf("%d", options.Limit))

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToJSON(options.Filters)
		if err != nil ***REMOVED***
			return results, err
		***REMOVED***
		query.Set("filters", filterJSON)
	***REMOVED***

	resp, err := cli.tryImageSearch(ctx, query, options.RegistryAuth)
	if resp.statusCode == http.StatusUnauthorized && options.PrivilegeFunc != nil ***REMOVED***
		newAuthHeader, privilegeErr := options.PrivilegeFunc()
		if privilegeErr != nil ***REMOVED***
			return results, privilegeErr
		***REMOVED***
		resp, err = cli.tryImageSearch(ctx, query, newAuthHeader)
	***REMOVED***
	if err != nil ***REMOVED***
		return results, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&results)
	ensureReaderClosed(resp)
	return results, err
***REMOVED***

func (cli *Client) tryImageSearch(ctx context.Context, query url.Values, registryAuth string) (serverResponse, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"X-Registry-Auth": ***REMOVED***registryAuth***REMOVED******REMOVED***
	return cli.get(ctx, "/images/search", query, headers)
***REMOVED***
