package client

import (
	"encoding/json"
	"net/url"

	registrytypes "github.com/docker/docker/api/types/registry"
	"golang.org/x/net/context"
)

// DistributionInspect returns the image digest with full Manifest
func (cli *Client) DistributionInspect(ctx context.Context, image, encodedRegistryAuth string) (registrytypes.DistributionInspect, error) ***REMOVED***
	// Contact the registry to retrieve digest and platform information
	var distributionInspect registrytypes.DistributionInspect

	if err := cli.NewVersionError("1.30", "distribution inspect"); err != nil ***REMOVED***
		return distributionInspect, err
	***REMOVED***
	var headers map[string][]string

	if encodedRegistryAuth != "" ***REMOVED***
		headers = map[string][]string***REMOVED***
			"X-Registry-Auth": ***REMOVED***encodedRegistryAuth***REMOVED***,
		***REMOVED***
	***REMOVED***

	resp, err := cli.get(ctx, "/distribution/"+image+"/json", url.Values***REMOVED******REMOVED***, headers)
	if err != nil ***REMOVED***
		return distributionInspect, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&distributionInspect)
	ensureReaderClosed(resp)
	return distributionInspect, err
***REMOVED***
