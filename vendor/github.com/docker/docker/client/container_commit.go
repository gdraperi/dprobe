package client

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerCommit applies changes into a container and creates a new tagged image.
func (cli *Client) ContainerCommit(ctx context.Context, container string, options types.ContainerCommitOptions) (types.IDResponse, error) ***REMOVED***
	var repository, tag string
	if options.Reference != "" ***REMOVED***
		ref, err := reference.ParseNormalizedNamed(options.Reference)
		if err != nil ***REMOVED***
			return types.IDResponse***REMOVED******REMOVED***, err
		***REMOVED***

		if _, isCanonical := ref.(reference.Canonical); isCanonical ***REMOVED***
			return types.IDResponse***REMOVED******REMOVED***, errors.New("refusing to create a tag with a digest reference")
		***REMOVED***
		ref = reference.TagNameOnly(ref)

		if tagged, ok := ref.(reference.Tagged); ok ***REMOVED***
			tag = tagged.Tag()
		***REMOVED***
		repository = reference.FamiliarName(ref)
	***REMOVED***

	query := url.Values***REMOVED******REMOVED***
	query.Set("container", container)
	query.Set("repo", repository)
	query.Set("tag", tag)
	query.Set("comment", options.Comment)
	query.Set("author", options.Author)
	for _, change := range options.Changes ***REMOVED***
		query.Add("changes", change)
	***REMOVED***
	if !options.Pause ***REMOVED***
		query.Set("pause", "0")
	***REMOVED***

	var response types.IDResponse
	resp, err := cli.post(ctx, "/commit", query, options.Config, nil)
	if err != nil ***REMOVED***
		return response, err
	***REMOVED***

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
***REMOVED***
