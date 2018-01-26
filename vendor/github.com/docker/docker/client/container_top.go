package client

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types/container"
	"golang.org/x/net/context"
)

// ContainerTop shows process information from within a container.
func (cli *Client) ContainerTop(ctx context.Context, containerID string, arguments []string) (container.ContainerTopOKBody, error) ***REMOVED***
	var response container.ContainerTopOKBody
	query := url.Values***REMOVED******REMOVED***
	if len(arguments) > 0 ***REMOVED***
		query.Set("ps_args", strings.Join(arguments, " "))
	***REMOVED***

	resp, err := cli.get(ctx, "/containers/"+containerID+"/top", query, nil)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
***REMOVED***
