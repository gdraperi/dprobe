package client

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
)

// ImagePull requests the docker host to pull an image from a remote registry.
// It executes the privileged function if the operation is unauthorized
// and it tries one more time.
// It's up to the caller to handle the io.ReadCloser and close it properly.
//
// FIXME(vdemeester): there is currently used in a few way in docker/docker
// - if not in trusted content, ref is used to pass the whole reference, and tag is empty
// - if in trusted content, ref is used to pass the reference name, and tag for the digest
func (cli *Client) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) ***REMOVED***
	ref, err := reference.ParseNormalizedNamed(refStr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	query := url.Values***REMOVED******REMOVED***
	query.Set("fromImage", reference.FamiliarName(ref))
	if !options.All ***REMOVED***
		query.Set("tag", getAPITagFromNamedRef(ref))
	***REMOVED***
	if options.Platform != "" ***REMOVED***
		query.Set("platform", strings.ToLower(options.Platform))
	***REMOVED***

	resp, err := cli.tryImageCreate(ctx, query, options.RegistryAuth)
	if resp.statusCode == http.StatusUnauthorized && options.PrivilegeFunc != nil ***REMOVED***
		newAuthHeader, privilegeErr := options.PrivilegeFunc()
		if privilegeErr != nil ***REMOVED***
			return nil, privilegeErr
		***REMOVED***
		resp, err = cli.tryImageCreate(ctx, query, newAuthHeader)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp.body, nil
***REMOVED***

// getAPITagFromNamedRef returns a tag from the specified reference.
// This function is necessary as long as the docker "server" api expects
// digests to be sent as tags and makes a distinction between the name
// and tag/digest part of a reference.
func getAPITagFromNamedRef(ref reference.Named) string ***REMOVED***
	if digested, ok := ref.(reference.Digested); ok ***REMOVED***
		return digested.Digest().String()
	***REMOVED***
	ref = reference.TagNameOnly(ref)
	if tagged, ok := ref.(reference.Tagged); ok ***REMOVED***
		return tagged.Tag()
	***REMOVED***
	return ""
***REMOVED***
