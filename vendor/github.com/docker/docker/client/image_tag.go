package client

import (
	"net/url"

	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// ImageTag tags an image in the docker host
func (cli *Client) ImageTag(ctx context.Context, source, target string) error ***REMOVED***
	if _, err := reference.ParseAnyReference(source); err != nil ***REMOVED***
		return errors.Wrapf(err, "Error parsing reference: %q is not a valid repository/tag", source)
	***REMOVED***

	ref, err := reference.ParseNormalizedNamed(target)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "Error parsing reference: %q is not a valid repository/tag", target)
	***REMOVED***

	if _, isCanonical := ref.(reference.Canonical); isCanonical ***REMOVED***
		return errors.New("refusing to create a tag with a digest reference")
	***REMOVED***

	ref = reference.TagNameOnly(ref)

	query := url.Values***REMOVED******REMOVED***
	query.Set("repo", reference.FamiliarName(ref))
	if tagged, ok := ref.(reference.Tagged); ok ***REMOVED***
		query.Set("tag", tagged.Tag())
	***REMOVED***

	resp, err := cli.post(ctx, "/images/"+source+"/tag", query, nil, nil)
	ensureReaderClosed(resp)
	return err
***REMOVED***
