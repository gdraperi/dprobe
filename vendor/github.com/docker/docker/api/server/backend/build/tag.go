package build

import (
	"fmt"
	"io"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/image"
	"github.com/pkg/errors"
)

// Tagger is responsible for tagging an image created by a builder
type Tagger struct ***REMOVED***
	imageComponent ImageComponent
	stdout         io.Writer
	repoAndTags    []reference.Named
***REMOVED***

// NewTagger returns a new Tagger for tagging the images of a build.
// If any of the names are invalid tags an error is returned.
func NewTagger(backend ImageComponent, stdout io.Writer, names []string) (*Tagger, error) ***REMOVED***
	reposAndTags, err := sanitizeRepoAndTags(names)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Tagger***REMOVED***
		imageComponent: backend,
		stdout:         stdout,
		repoAndTags:    reposAndTags,
	***REMOVED***, nil
***REMOVED***

// TagImages creates image tags for the imageID
func (bt *Tagger) TagImages(imageID image.ID) error ***REMOVED***
	for _, rt := range bt.repoAndTags ***REMOVED***
		if err := bt.imageComponent.TagImageWithReference(imageID, rt); err != nil ***REMOVED***
			return err
		***REMOVED***
		fmt.Fprintf(bt.stdout, "Successfully tagged %s\n", reference.FamiliarString(rt))
	***REMOVED***
	return nil
***REMOVED***

// sanitizeRepoAndTags parses the raw "t" parameter received from the client
// to a slice of repoAndTag.
// It also validates each repoName and tag.
func sanitizeRepoAndTags(names []string) ([]reference.Named, error) ***REMOVED***
	var (
		repoAndTags []reference.Named
		// This map is used for deduplicating the "-t" parameter.
		uniqNames = make(map[string]struct***REMOVED******REMOVED***)
	)
	for _, repo := range names ***REMOVED***
		if repo == "" ***REMOVED***
			continue
		***REMOVED***

		ref, err := reference.ParseNormalizedNamed(repo)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if _, isCanonical := ref.(reference.Canonical); isCanonical ***REMOVED***
			return nil, errors.New("build tag cannot contain a digest")
		***REMOVED***

		ref = reference.TagNameOnly(ref)

		nameWithTag := ref.String()

		if _, exists := uniqNames[nameWithTag]; !exists ***REMOVED***
			uniqNames[nameWithTag] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			repoAndTags = append(repoAndTags, ref)
		***REMOVED***
	***REMOVED***
	return repoAndTags, nil
***REMOVED***
