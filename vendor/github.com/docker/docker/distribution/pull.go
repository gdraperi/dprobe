package distribution

import (
	"fmt"
	"runtime"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/pkg/progress"
	refstore "github.com/docker/docker/reference"
	"github.com/docker/docker/registry"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Puller is an interface that abstracts pulling for different API versions.
type Puller interface ***REMOVED***
	// Pull tries to pull the image referenced by `tag`
	// Pull returns an error if any, as well as a boolean that determines whether to retry Pull on the next configured endpoint.
	//
	Pull(ctx context.Context, ref reference.Named, os string) error
***REMOVED***

// newPuller returns a Puller interface that will pull from either a v1 or v2
// registry. The endpoint argument contains a Version field that determines
// whether a v1 or v2 puller will be created. The other parameters are passed
// through to the underlying puller implementation for use during the actual
// pull operation.
func newPuller(endpoint registry.APIEndpoint, repoInfo *registry.RepositoryInfo, imagePullConfig *ImagePullConfig) (Puller, error) ***REMOVED***
	switch endpoint.Version ***REMOVED***
	case registry.APIVersion2:
		return &v2Puller***REMOVED***
			V2MetadataService: metadata.NewV2MetadataService(imagePullConfig.MetadataStore),
			endpoint:          endpoint,
			config:            imagePullConfig,
			repoInfo:          repoInfo,
		***REMOVED***, nil
	case registry.APIVersion1:
		return &v1Puller***REMOVED***
			v1IDService: metadata.NewV1IDService(imagePullConfig.MetadataStore),
			endpoint:    endpoint,
			config:      imagePullConfig,
			repoInfo:    repoInfo,
		***REMOVED***, nil
	***REMOVED***
	return nil, fmt.Errorf("unknown version %d for registry %s", endpoint.Version, endpoint.URL)
***REMOVED***

// Pull initiates a pull operation. image is the repository name to pull, and
// tag may be either empty, or indicate a specific tag to pull.
func Pull(ctx context.Context, ref reference.Named, imagePullConfig *ImagePullConfig) error ***REMOVED***
	// Resolve the Repository name from fqn to RepositoryInfo
	repoInfo, err := imagePullConfig.RegistryService.ResolveRepository(ref)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// makes sure name is not `scratch`
	if err := ValidateRepoName(repoInfo.Name); err != nil ***REMOVED***
		return err
	***REMOVED***

	endpoints, err := imagePullConfig.RegistryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var (
		lastErr error

		// discardNoSupportErrors is used to track whether an endpoint encountered an error of type registry.ErrNoSupport
		// By default it is false, which means that if an ErrNoSupport error is encountered, it will be saved in lastErr.
		// As soon as another kind of error is encountered, discardNoSupportErrors is set to true, avoiding the saving of
		// any subsequent ErrNoSupport errors in lastErr.
		// It's needed for pull-by-digest on v1 endpoints: if there are only v1 endpoints configured, the error should be
		// returned and displayed, but if there was a v2 endpoint which supports pull-by-digest, then the last relevant
		// error is the ones from v2 endpoints not v1.
		discardNoSupportErrors bool

		// confirmedV2 is set to true if a pull attempt managed to
		// confirm that it was talking to a v2 registry. This will
		// prevent fallback to the v1 protocol.
		confirmedV2 bool

		// confirmedTLSRegistries is a map indicating which registries
		// are known to be using TLS. There should never be a plaintext
		// retry for any of these.
		confirmedTLSRegistries = make(map[string]struct***REMOVED******REMOVED***)
	)
	for _, endpoint := range endpoints ***REMOVED***
		if imagePullConfig.RequireSchema2 && endpoint.Version == registry.APIVersion1 ***REMOVED***
			continue
		***REMOVED***

		if confirmedV2 && endpoint.Version == registry.APIVersion1 ***REMOVED***
			logrus.Debugf("Skipping v1 endpoint %s because v2 registry was detected", endpoint.URL)
			continue
		***REMOVED***

		if endpoint.URL.Scheme != "https" ***REMOVED***
			if _, confirmedTLS := confirmedTLSRegistries[endpoint.URL.Host]; confirmedTLS ***REMOVED***
				logrus.Debugf("Skipping non-TLS endpoint %s for host/port that appears to use TLS", endpoint.URL)
				continue
			***REMOVED***
		***REMOVED***

		logrus.Debugf("Trying to pull %s from %s %s", reference.FamiliarName(repoInfo.Name), endpoint.URL, endpoint.Version)

		puller, err := newPuller(endpoint, repoInfo, imagePullConfig)
		if err != nil ***REMOVED***
			lastErr = err
			continue
		***REMOVED***

		// Make sure we default the OS if it hasn't been supplied
		if imagePullConfig.OS == "" ***REMOVED***
			imagePullConfig.OS = runtime.GOOS
		***REMOVED***

		if err := puller.Pull(ctx, ref, imagePullConfig.OS); err != nil ***REMOVED***
			// Was this pull cancelled? If so, don't try to fall
			// back.
			fallback := false
			select ***REMOVED***
			case <-ctx.Done():
			default:
				if fallbackErr, ok := err.(fallbackError); ok ***REMOVED***
					fallback = true
					confirmedV2 = confirmedV2 || fallbackErr.confirmedV2
					if fallbackErr.transportOK && endpoint.URL.Scheme == "https" ***REMOVED***
						confirmedTLSRegistries[endpoint.URL.Host] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					***REMOVED***
					err = fallbackErr.err
				***REMOVED***
			***REMOVED***
			if fallback ***REMOVED***
				if _, ok := err.(ErrNoSupport); !ok ***REMOVED***
					// Because we found an error that's not ErrNoSupport, discard all subsequent ErrNoSupport errors.
					discardNoSupportErrors = true
					// append subsequent errors
					lastErr = err
				***REMOVED*** else if !discardNoSupportErrors ***REMOVED***
					// Save the ErrNoSupport error, because it's either the first error or all encountered errors
					// were also ErrNoSupport errors.
					// append subsequent errors
					lastErr = err
				***REMOVED***
				logrus.Infof("Attempting next endpoint for pull after error: %v", err)
				continue
			***REMOVED***
			logrus.Errorf("Not continuing with pull after error: %v", err)
			return TranslatePullError(err, ref)
		***REMOVED***

		imagePullConfig.ImageEventLogger(reference.FamiliarString(ref), reference.FamiliarName(repoInfo.Name), "pull")
		return nil
	***REMOVED***

	if lastErr == nil ***REMOVED***
		lastErr = fmt.Errorf("no endpoints found for %s", reference.FamiliarString(ref))
	***REMOVED***

	return TranslatePullError(lastErr, ref)
***REMOVED***

// writeStatus writes a status message to out. If layersDownloaded is true, the
// status message indicates that a newer image was downloaded. Otherwise, it
// indicates that the image is up to date. requestedTag is the tag the message
// will refer to.
func writeStatus(requestedTag string, out progress.Output, layersDownloaded bool) ***REMOVED***
	if layersDownloaded ***REMOVED***
		progress.Message(out, "", "Status: Downloaded newer image for "+requestedTag)
	***REMOVED*** else ***REMOVED***
		progress.Message(out, "", "Status: Image is up to date for "+requestedTag)
	***REMOVED***
***REMOVED***

// ValidateRepoName validates the name of a repository.
func ValidateRepoName(name reference.Named) error ***REMOVED***
	if reference.FamiliarName(name) == api.NoBaseImageSpecifier ***REMOVED***
		return errors.WithStack(reservedNameError(api.NoBaseImageSpecifier))
	***REMOVED***
	return nil
***REMOVED***

func addDigestReference(store refstore.Store, ref reference.Named, dgst digest.Digest, id digest.Digest) error ***REMOVED***
	dgstRef, err := reference.WithDigest(reference.TrimNamed(ref), dgst)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if oldTagID, err := store.Get(dgstRef); err == nil ***REMOVED***
		if oldTagID != id ***REMOVED***
			// Updating digests not supported by reference store
			logrus.Errorf("Image ID for digest %s changed from %s to %s, cannot update", dgst.String(), oldTagID, id)
		***REMOVED***
		return nil
	***REMOVED*** else if err != refstore.ErrDoesNotExist ***REMOVED***
		return err
	***REMOVED***

	return store.AddDigest(dgstRef, id, true)
***REMOVED***
