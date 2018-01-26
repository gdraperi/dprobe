package daemon

import (
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/image"
)

// TagImage creates the tag specified by newTag, pointing to the image named
// imageName (alternatively, imageName can also be an image ID).
func (daemon *Daemon) TagImage(imageName, repository, tag string) error ***REMOVED***
	imageID, _, err := daemon.GetImageIDAndOS(imageName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	newTag, err := reference.ParseNormalizedNamed(repository)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tag != "" ***REMOVED***
		if newTag, err = reference.WithTag(reference.TrimNamed(newTag), tag); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return daemon.TagImageWithReference(imageID, newTag)
***REMOVED***

// TagImageWithReference adds the given reference to the image ID provided.
func (daemon *Daemon) TagImageWithReference(imageID image.ID, newTag reference.Named) error ***REMOVED***
	if err := daemon.referenceStore.AddTag(newTag, imageID.Digest(), true); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := daemon.imageStore.SetLastUpdated(imageID); err != nil ***REMOVED***
		return err
	***REMOVED***
	daemon.LogImageEvent(imageID.String(), reference.FamiliarString(newTag), "tag")
	return nil
***REMOVED***
