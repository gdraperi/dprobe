package client

import (
	"io"
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// PluginCreate creates a plugin
func (cli *Client) PluginCreate(ctx context.Context, createContext io.Reader, createOptions types.PluginCreateOptions) error ***REMOVED***
	headers := http.Header(make(map[string][]string))
	headers.Set("Content-Type", "application/x-tar")

	query := url.Values***REMOVED******REMOVED***
	query.Set("name", createOptions.RepoName)

	resp, err := cli.postRaw(ctx, "/plugins/create", query, createContext, headers)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	ensureReaderClosed(resp)
	return err
***REMOVED***
