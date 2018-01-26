package daemon

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
)

type conflictType int

const (
	conflictDependentChild conflictType = (1 << iota)
	conflictRunningContainer
	conflictActiveReference
	conflictStoppedContainer
	conflictHard = conflictDependentChild | conflictRunningContainer
	conflictSoft = conflictActiveReference | conflictStoppedContainer
)

// ImageDelete deletes the image referenced by the given imageRef from this
// daemon. The given imageRef can be an image ID, ID prefix, or a repository
// reference (with an optional tag or digest, defaulting to the tag name
// "latest"). There is differing behavior depending on whether the given
// imageRef is a repository reference or not.
//
// If the given imageRef is a repository reference then that repository
// reference will be removed. However, if there exists any containers which
// were created using the same image reference then the repository reference
// cannot be removed unless either there are other repository references to the
// same image or force is true. Following removal of the repository reference,
// the referenced image itself will attempt to be deleted as described below
// but quietly, meaning any image delete conflicts will cause the image to not
// be deleted and the conflict will not be reported.
//
// There may be conflicts preventing deletion of an image and these conflicts
// are divided into two categories grouped by their severity:
//
// Hard Conflict:
// 	- a pull or build using the image.
// 	- any descendant image.
// 	- any running container using the image.
//
// Soft Conflict:
// 	- any stopped container using the image.
// 	- any repository tag or digest references to the image.
//
// The image cannot be removed if there are any hard conflicts and can be
// removed if there are soft conflicts only if force is true.
//
// If prune is true, ancestor images will each attempt to be deleted quietly,
// meaning any delete conflicts will cause the image to not be deleted and the
// conflict will not be reported.
//
// FIXME: remove ImageDelete's dependency on Daemon, then move to the graph
// package. This would require that we no longer need the daemon to determine
// whether images are being used by a stopped or running container.
func (daemon *Daemon) ImageDelete(imageRef string, force, prune bool) ([]types.ImageDeleteResponseItem, error) ***REMOVED***
	start := time.Now()
	records := []types.ImageDeleteResponseItem***REMOVED******REMOVED***

	imgID, operatingSystem, err := daemon.GetImageIDAndOS(imageRef)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !system.IsOSSupported(operatingSystem) ***REMOVED***
		return nil, errors.Errorf("unable to delete image: %q", system.ErrNotSupportedOperatingSystem)
	***REMOVED***

	repoRefs := daemon.referenceStore.References(imgID.Digest())

	var removedRepositoryRef bool
	if !isImageIDPrefix(imgID.String(), imageRef) ***REMOVED***
		// A repository reference was given and should be removed
		// first. We can only remove this reference if either force is
		// true, there are multiple repository references to this
		// image, or there are no containers using the given reference.
		if !force && isSingleReference(repoRefs) ***REMOVED***
			if container := daemon.getContainerUsingImage(imgID); container != nil ***REMOVED***
				// If we removed the repository reference then
				// this image would remain "dangling" and since
				// we really want to avoid that the client must
				// explicitly force its removal.
				err := errors.Errorf("conflict: unable to remove repository reference %q (must force) - container %s is using its referenced image %s", imageRef, stringid.TruncateID(container.ID), stringid.TruncateID(imgID.String()))
				return nil, errdefs.Conflict(err)
			***REMOVED***
		***REMOVED***

		parsedRef, err := reference.ParseNormalizedNamed(imageRef)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		parsedRef, err = daemon.removeImageRef(parsedRef)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		untaggedRecord := types.ImageDeleteResponseItem***REMOVED***Untagged: reference.FamiliarString(parsedRef)***REMOVED***

		daemon.LogImageEvent(imgID.String(), imgID.String(), "untag")
		records = append(records, untaggedRecord)

		repoRefs = daemon.referenceStore.References(imgID.Digest())

		// If a tag reference was removed and the only remaining
		// references to the same repository are digest references,
		// then clean up those digest references.
		if _, isCanonical := parsedRef.(reference.Canonical); !isCanonical ***REMOVED***
			foundRepoTagRef := false
			for _, repoRef := range repoRefs ***REMOVED***
				if _, repoRefIsCanonical := repoRef.(reference.Canonical); !repoRefIsCanonical && parsedRef.Name() == repoRef.Name() ***REMOVED***
					foundRepoTagRef = true
					break
				***REMOVED***
			***REMOVED***
			if !foundRepoTagRef ***REMOVED***
				// Remove canonical references from same repository
				remainingRefs := []reference.Named***REMOVED******REMOVED***
				for _, repoRef := range repoRefs ***REMOVED***
					if _, repoRefIsCanonical := repoRef.(reference.Canonical); repoRefIsCanonical && parsedRef.Name() == repoRef.Name() ***REMOVED***
						if _, err := daemon.removeImageRef(repoRef); err != nil ***REMOVED***
							return records, err
						***REMOVED***

						untaggedRecord := types.ImageDeleteResponseItem***REMOVED***Untagged: reference.FamiliarString(repoRef)***REMOVED***
						records = append(records, untaggedRecord)
					***REMOVED*** else ***REMOVED***
						remainingRefs = append(remainingRefs, repoRef)

					***REMOVED***
				***REMOVED***
				repoRefs = remainingRefs
			***REMOVED***
		***REMOVED***

		// If it has remaining references then the untag finished the remove
		if len(repoRefs) > 0 ***REMOVED***
			return records, nil
		***REMOVED***

		removedRepositoryRef = true
	***REMOVED*** else ***REMOVED***
		// If an ID reference was given AND there is at most one tag
		// reference to the image AND all references are within one
		// repository, then remove all references.
		if isSingleReference(repoRefs) ***REMOVED***
			c := conflictHard
			if !force ***REMOVED***
				c |= conflictSoft &^ conflictActiveReference
			***REMOVED***
			if conflict := daemon.checkImageDeleteConflict(imgID, c); conflict != nil ***REMOVED***
				return nil, conflict
			***REMOVED***

			for _, repoRef := range repoRefs ***REMOVED***
				parsedRef, err := daemon.removeImageRef(repoRef)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				untaggedRecord := types.ImageDeleteResponseItem***REMOVED***Untagged: reference.FamiliarString(parsedRef)***REMOVED***

				daemon.LogImageEvent(imgID.String(), imgID.String(), "untag")
				records = append(records, untaggedRecord)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := daemon.imageDeleteHelper(imgID, &records, force, prune, removedRepositoryRef); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	imageActions.WithValues("delete").UpdateSince(start)

	return records, nil
***REMOVED***

// isSingleReference returns true when all references are from one repository
// and there is at most one tag. Returns false for empty input.
func isSingleReference(repoRefs []reference.Named) bool ***REMOVED***
	if len(repoRefs) <= 1 ***REMOVED***
		return len(repoRefs) == 1
	***REMOVED***
	var singleRef reference.Named
	canonicalRefs := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, repoRef := range repoRefs ***REMOVED***
		if _, isCanonical := repoRef.(reference.Canonical); isCanonical ***REMOVED***
			canonicalRefs[repoRef.Name()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED*** else if singleRef == nil ***REMOVED***
			singleRef = repoRef
		***REMOVED*** else ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if singleRef == nil ***REMOVED***
		// Just use first canonical ref
		singleRef = repoRefs[0]
	***REMOVED***
	_, ok := canonicalRefs[singleRef.Name()]
	return len(canonicalRefs) == 1 && ok
***REMOVED***

// isImageIDPrefix returns whether the given possiblePrefix is a prefix of the
// given imageID.
func isImageIDPrefix(imageID, possiblePrefix string) bool ***REMOVED***
	if strings.HasPrefix(imageID, possiblePrefix) ***REMOVED***
		return true
	***REMOVED***

	if i := strings.IndexRune(imageID, ':'); i >= 0 ***REMOVED***
		return strings.HasPrefix(imageID[i+1:], possiblePrefix)
	***REMOVED***

	return false
***REMOVED***

// getContainerUsingImage returns a container that was created using the given
// imageID. Returns nil if there is no such container.
func (daemon *Daemon) getContainerUsingImage(imageID image.ID) *container.Container ***REMOVED***
	return daemon.containers.First(func(c *container.Container) bool ***REMOVED***
		return c.ImageID == imageID
	***REMOVED***)
***REMOVED***

// removeImageRef attempts to parse and remove the given image reference from
// this daemon's store of repository tag/digest references. The given
// repositoryRef must not be an image ID but a repository name followed by an
// optional tag or digest reference. If tag or digest is omitted, the default
// tag is used. Returns the resolved image reference and an error.
func (daemon *Daemon) removeImageRef(ref reference.Named) (reference.Named, error) ***REMOVED***
	ref = reference.TagNameOnly(ref)

	// Ignore the boolean value returned, as far as we're concerned, this
	// is an idempotent operation and it's okay if the reference didn't
	// exist in the first place.
	_, err := daemon.referenceStore.Delete(ref)

	return ref, err
***REMOVED***

// removeAllReferencesToImageID attempts to remove every reference to the given
// imgID from this daemon's store of repository tag/digest references. Returns
// on the first encountered error. Removed references are logged to this
// daemon's event service. An "Untagged" types.ImageDeleteResponseItem is added to the
// given list of records.
func (daemon *Daemon) removeAllReferencesToImageID(imgID image.ID, records *[]types.ImageDeleteResponseItem) error ***REMOVED***
	imageRefs := daemon.referenceStore.References(imgID.Digest())

	for _, imageRef := range imageRefs ***REMOVED***
		parsedRef, err := daemon.removeImageRef(imageRef)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		untaggedRecord := types.ImageDeleteResponseItem***REMOVED***Untagged: reference.FamiliarString(parsedRef)***REMOVED***

		daemon.LogImageEvent(imgID.String(), imgID.String(), "untag")
		*records = append(*records, untaggedRecord)
	***REMOVED***

	return nil
***REMOVED***

// ImageDeleteConflict holds a soft or hard conflict and an associated error.
// Implements the error interface.
type imageDeleteConflict struct ***REMOVED***
	hard    bool
	used    bool
	imgID   image.ID
	message string
***REMOVED***

func (idc *imageDeleteConflict) Error() string ***REMOVED***
	var forceMsg string
	if idc.hard ***REMOVED***
		forceMsg = "cannot be forced"
	***REMOVED*** else ***REMOVED***
		forceMsg = "must be forced"
	***REMOVED***

	return fmt.Sprintf("conflict: unable to delete %s (%s) - %s", stringid.TruncateID(idc.imgID.String()), forceMsg, idc.message)
***REMOVED***

func (idc *imageDeleteConflict) Conflict() ***REMOVED******REMOVED***

// imageDeleteHelper attempts to delete the given image from this daemon. If
// the image has any hard delete conflicts (child images or running containers
// using the image) then it cannot be deleted. If the image has any soft delete
// conflicts (any tags/digests referencing the image or any stopped container
// using the image) then it can only be deleted if force is true. If the delete
// succeeds and prune is true, the parent images are also deleted if they do
// not have any soft or hard delete conflicts themselves. Any deleted images
// and untagged references are appended to the given records. If any error or
// conflict is encountered, it will be returned immediately without deleting
// the image. If quiet is true, any encountered conflicts will be ignored and
// the function will return nil immediately without deleting the image.
func (daemon *Daemon) imageDeleteHelper(imgID image.ID, records *[]types.ImageDeleteResponseItem, force, prune, quiet bool) error ***REMOVED***
	// First, determine if this image has any conflicts. Ignore soft conflicts
	// if force is true.
	c := conflictHard
	if !force ***REMOVED***
		c |= conflictSoft
	***REMOVED***
	if conflict := daemon.checkImageDeleteConflict(imgID, c); conflict != nil ***REMOVED***
		if quiet && (!daemon.imageIsDangling(imgID) || conflict.used) ***REMOVED***
			// Ignore conflicts UNLESS the image is "dangling" or not being used in
			// which case we want the user to know.
			return nil
		***REMOVED***

		// There was a conflict and it's either a hard conflict OR we are not
		// forcing deletion on soft conflicts.
		return conflict
	***REMOVED***

	parent, err := daemon.imageStore.GetParent(imgID)
	if err != nil ***REMOVED***
		// There may be no parent
		parent = ""
	***REMOVED***

	// Delete all repository tag/digest references to this image.
	if err := daemon.removeAllReferencesToImageID(imgID, records); err != nil ***REMOVED***
		return err
	***REMOVED***

	removedLayers, err := daemon.imageStore.Delete(imgID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	daemon.LogImageEvent(imgID.String(), imgID.String(), "delete")
	*records = append(*records, types.ImageDeleteResponseItem***REMOVED***Deleted: imgID.String()***REMOVED***)
	for _, removedLayer := range removedLayers ***REMOVED***
		*records = append(*records, types.ImageDeleteResponseItem***REMOVED***Deleted: removedLayer.ChainID.String()***REMOVED***)
	***REMOVED***

	if !prune || parent == "" ***REMOVED***
		return nil
	***REMOVED***

	// We need to prune the parent image. This means delete it if there are
	// no tags/digests referencing it and there are no containers using it (
	// either running or stopped).
	// Do not force prunings, but do so quietly (stopping on any encountered
	// conflicts).
	return daemon.imageDeleteHelper(parent, records, false, true, true)
***REMOVED***

// checkImageDeleteConflict determines whether there are any conflicts
// preventing deletion of the given image from this daemon. A hard conflict is
// any image which has the given image as a parent or any running container
// using the image. A soft conflict is any tags/digest referencing the given
// image or any stopped container using the image. If ignoreSoftConflicts is
// true, this function will not check for soft conflict conditions.
func (daemon *Daemon) checkImageDeleteConflict(imgID image.ID, mask conflictType) *imageDeleteConflict ***REMOVED***
	// Check if the image has any descendant images.
	if mask&conflictDependentChild != 0 && len(daemon.imageStore.Children(imgID)) > 0 ***REMOVED***
		return &imageDeleteConflict***REMOVED***
			hard:    true,
			imgID:   imgID,
			message: "image has dependent child images",
		***REMOVED***
	***REMOVED***

	if mask&conflictRunningContainer != 0 ***REMOVED***
		// Check if any running container is using the image.
		running := func(c *container.Container) bool ***REMOVED***
			return c.IsRunning() && c.ImageID == imgID
		***REMOVED***
		if container := daemon.containers.First(running); container != nil ***REMOVED***
			return &imageDeleteConflict***REMOVED***
				imgID:   imgID,
				hard:    true,
				used:    true,
				message: fmt.Sprintf("image is being used by running container %s", stringid.TruncateID(container.ID)),
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check if any repository tags/digest reference this image.
	if mask&conflictActiveReference != 0 && len(daemon.referenceStore.References(imgID.Digest())) > 0 ***REMOVED***
		return &imageDeleteConflict***REMOVED***
			imgID:   imgID,
			message: "image is referenced in multiple repositories",
		***REMOVED***
	***REMOVED***

	if mask&conflictStoppedContainer != 0 ***REMOVED***
		// Check if any stopped containers reference this image.
		stopped := func(c *container.Container) bool ***REMOVED***
			return !c.IsRunning() && c.ImageID == imgID
		***REMOVED***
		if container := daemon.containers.First(stopped); container != nil ***REMOVED***
			return &imageDeleteConflict***REMOVED***
				imgID:   imgID,
				used:    true,
				message: fmt.Sprintf("image is being used by stopped container %s", stringid.TruncateID(container.ID)),
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// imageIsDangling returns whether the given image is "dangling" which means
// that there are no repository references to the given image and it has no
// child images.
func (daemon *Daemon) imageIsDangling(imgID image.ID) bool ***REMOVED***
	return !(len(daemon.referenceStore.References(imgID.Digest())) > 0 || len(daemon.imageStore.Children(imgID)) > 0)
***REMOVED***
