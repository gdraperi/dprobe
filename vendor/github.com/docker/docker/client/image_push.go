package client

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
)

// ImagePush requests the docker host to push an image to a remote registry.
// It executes the privileged function if the operation is unauthorized
// and it tries one more time.
// It's up to the caller to handle the io.ReadCloser and close it properly.
func (cli *Client) ImagePush(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error) ***REMOVED***
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if _, isCanonical := ref.(reference.Canonical); isCanonical ***REMOVED***
		return nil, errors.New("cannot push a digest reference")
	***REMOVED***

	tag := ""
	name := reference.FamiliarName(ref)

	if nameTaggedRef, isNamedTagged := ref.(reference.NamedTagged); isNamedTagged ***REMOVED***
		tag = nameTaggedRef.Tag()
	***REMOVED***

	query := url.Values***REMOVED******REMOVED***
	query.Set("tag", tag)

	resp, err := cli.tryImagePush(ctx, name, query, options.RegistryAuth)
	if resp.statusCode == http.StatusUnauthorized && options.PrivilegeFunc != nil ***REMOVED***
		newAuthHeader, privilegeErr := options.PrivilegeFunc()
		if privilegeErr != nil ***REMOVED***
			return nil, privilegeErr
		***REMOVED***
		resp, err = cli.tryImagePush(ctx, name, query, newAuthHeader)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp.body, nil
***REMOVED***

func (cli *Client) tryImagePush(ctx context.Context, imageID string, query url.Values, registryAuth string) (serverResponse, error) ***REMOVED***
	headers := map[string][]string***REMOVED***"X-Registry-Auth": ***REMOVED***registryAuth***REMOVED******REMOVED***
	return cli.post(ctx, "/images/"+imageID+"/push", query, nil, headers)
***REMOVED***
