package daemon

import (
	"fmt"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
)

// errImageDoesNotExist is error returned when no image can be found for a reference.
type errImageDoesNotExist struct ***REMOVED***
	ref reference.Reference
***REMOVED***

func (e errImageDoesNotExist) Error() string ***REMOVED***
	ref := e.ref
	if named, ok := ref.(reference.Named); ok ***REMOVED***
		ref = reference.TagNameOnly(named)
	***REMOVED***
	return fmt.Sprintf("No such image: %s", reference.FamiliarString(ref))
***REMOVED***

func (e errImageDoesNotExist) NotFound() ***REMOVED******REMOVED***

// GetImageIDAndOS returns an image ID and operating system corresponding to the image referred to by
// refOrID.
func (daemon *Daemon) GetImageIDAndOS(refOrID string) (image.ID, string, error) ***REMOVED***
	ref, err := reference.ParseAnyReference(refOrID)
	if err != nil ***REMOVED***
		return "", "", errdefs.InvalidParameter(err)
	***REMOVED***
	namedRef, ok := ref.(reference.Named)
	if !ok ***REMOVED***
		digested, ok := ref.(reference.Digested)
		if !ok ***REMOVED***
			return "", "", errImageDoesNotExist***REMOVED***ref***REMOVED***
		***REMOVED***
		id := image.IDFromDigest(digested.Digest())
		if img, err := daemon.imageStore.Get(id); err == nil ***REMOVED***
			return id, img.OperatingSystem(), nil
		***REMOVED***
		return "", "", errImageDoesNotExist***REMOVED***ref***REMOVED***
	***REMOVED***

	if digest, err := daemon.referenceStore.Get(namedRef); err == nil ***REMOVED***
		// Search the image stores to get the operating system, defaulting to host OS.
		id := image.IDFromDigest(digest)
		if img, err := daemon.imageStore.Get(id); err == nil ***REMOVED***
			return id, img.OperatingSystem(), nil
		***REMOVED***
	***REMOVED***

	// Search based on ID
	if id, err := daemon.imageStore.Search(refOrID); err == nil ***REMOVED***
		img, err := daemon.imageStore.Get(id)
		if err != nil ***REMOVED***
			return "", "", errImageDoesNotExist***REMOVED***ref***REMOVED***
		***REMOVED***
		return id, img.OperatingSystem(), nil
	***REMOVED***

	return "", "", errImageDoesNotExist***REMOVED***ref***REMOVED***
***REMOVED***

// GetImage returns an image corresponding to the image referred to by refOrID.
func (daemon *Daemon) GetImage(refOrID string) (*image.Image, error) ***REMOVED***
	imgID, _, err := daemon.GetImageIDAndOS(refOrID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return daemon.imageStore.Get(imgID)
***REMOVED***
