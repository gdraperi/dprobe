package distribution

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/registry"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Pusher is an interface that abstracts pushing for different API versions.
type Pusher interface ***REMOVED***
	// Push tries to push the image configured at the creation of Pusher.
	// Push returns an error if any, as well as a boolean that determines whether to retry Push on the next configured endpoint.
	//
	// TODO(tiborvass): have Push() take a reference to repository + tag, so that the pusher itself is repository-agnostic.
	Push(ctx context.Context) error
***REMOVED***

const compressionBufSize = 32768

// NewPusher creates a new Pusher interface that will push to either a v1 or v2
// registry. The endpoint argument contains a Version field that determines
// whether a v1 or v2 pusher will be created. The other parameters are passed
// through to the underlying pusher implementation for use during the actual
// push operation.
func NewPusher(ref reference.Named, endpoint registry.APIEndpoint, repoInfo *registry.RepositoryInfo, imagePushConfig *ImagePushConfig) (Pusher, error) ***REMOVED***
	switch endpoint.Version ***REMOVED***
	case registry.APIVersion2:
		return &v2Pusher***REMOVED***
			v2MetadataService: metadata.NewV2MetadataService(imagePushConfig.MetadataStore),
			ref:               ref,
			endpoint:          endpoint,
			repoInfo:          repoInfo,
			config:            imagePushConfig,
		***REMOVED***, nil
	case registry.APIVersion1:
		return &v1Pusher***REMOVED***
			v1IDService: metadata.NewV1IDService(imagePushConfig.MetadataStore),
			ref:         ref,
			endpoint:    endpoint,
			repoInfo:    repoInfo,
			config:      imagePushConfig,
		***REMOVED***, nil
	***REMOVED***
	return nil, fmt.Errorf("unknown version %d for registry %s", endpoint.Version, endpoint.URL)
***REMOVED***

// Push initiates a push operation on ref.
// ref is the specific variant of the image to be pushed.
// If no tag is provided, all tags will be pushed.
func Push(ctx context.Context, ref reference.Named, imagePushConfig *ImagePushConfig) error ***REMOVED***
	// FIXME: Allow to interrupt current push when new push of same image is done.

	// Resolve the Repository name from fqn to RepositoryInfo
	repoInfo, err := imagePushConfig.RegistryService.ResolveRepository(ref)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	endpoints, err := imagePushConfig.RegistryService.LookupPushEndpoints(reference.Domain(repoInfo.Name))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	progress.Messagef(imagePushConfig.ProgressOutput, "", "The push refers to repository [%s]", repoInfo.Name.Name())

	associations := imagePushConfig.ReferenceStore.ReferencesByName(repoInfo.Name)
	if len(associations) == 0 ***REMOVED***
		return fmt.Errorf("An image does not exist locally with the tag: %s", reference.FamiliarName(repoInfo.Name))
	***REMOVED***

	var (
		lastErr error

		// confirmedV2 is set to true if a push attempt managed to
		// confirm that it was talking to a v2 registry. This will
		// prevent fallback to the v1 protocol.
		confirmedV2 bool

		// confirmedTLSRegistries is a map indicating which registries
		// are known to be using TLS. There should never be a plaintext
		// retry for any of these.
		confirmedTLSRegistries = make(map[string]struct***REMOVED******REMOVED***)
	)

	for _, endpoint := range endpoints ***REMOVED***
		if imagePushConfig.RequireSchema2 && endpoint.Version == registry.APIVersion1 ***REMOVED***
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

		logrus.Debugf("Trying to push %s to %s %s", repoInfo.Name.Name(), endpoint.URL, endpoint.Version)

		pusher, err := NewPusher(ref, endpoint, repoInfo, imagePushConfig)
		if err != nil ***REMOVED***
			lastErr = err
			continue
		***REMOVED***
		if err := pusher.Push(ctx); err != nil ***REMOVED***
			// Was this push cancelled? If so, don't try to fall
			// back.
			select ***REMOVED***
			case <-ctx.Done():
			default:
				if fallbackErr, ok := err.(fallbackError); ok ***REMOVED***
					confirmedV2 = confirmedV2 || fallbackErr.confirmedV2
					if fallbackErr.transportOK && endpoint.URL.Scheme == "https" ***REMOVED***
						confirmedTLSRegistries[endpoint.URL.Host] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					***REMOVED***
					err = fallbackErr.err
					lastErr = err
					logrus.Infof("Attempting next endpoint for push after error: %v", err)
					continue
				***REMOVED***
			***REMOVED***

			logrus.Errorf("Not continuing with push after error: %v", err)
			return err
		***REMOVED***

		imagePushConfig.ImageEventLogger(reference.FamiliarString(ref), reference.FamiliarName(repoInfo.Name), "push")
		return nil
	***REMOVED***

	if lastErr == nil ***REMOVED***
		lastErr = fmt.Errorf("no endpoints found for %s", repoInfo.Name.Name())
	***REMOVED***
	return lastErr
***REMOVED***

// compress returns an io.ReadCloser which will supply a compressed version of
// the provided Reader. The caller must close the ReadCloser after reading the
// compressed data.
//
// Note that this function returns a reader instead of taking a writer as an
// argument so that it can be used with httpBlobWriter's ReadFrom method.
// Using httpBlobWriter's Write method would send a PATCH request for every
// Write call.
//
// The second return value is a channel that gets closed when the goroutine
// is finished. This allows the caller to make sure the goroutine finishes
// before it releases any resources connected with the reader that was
// passed in.
func compress(in io.Reader) (io.ReadCloser, chan struct***REMOVED******REMOVED***) ***REMOVED***
	compressionDone := make(chan struct***REMOVED******REMOVED***)

	pipeReader, pipeWriter := io.Pipe()
	// Use a bufio.Writer to avoid excessive chunking in HTTP request.
	bufWriter := bufio.NewWriterSize(pipeWriter, compressionBufSize)
	compressor := gzip.NewWriter(bufWriter)

	go func() ***REMOVED***
		_, err := io.Copy(compressor, in)
		if err == nil ***REMOVED***
			err = compressor.Close()
		***REMOVED***
		if err == nil ***REMOVED***
			err = bufWriter.Flush()
		***REMOVED***
		if err != nil ***REMOVED***
			pipeWriter.CloseWithError(err)
		***REMOVED*** else ***REMOVED***
			pipeWriter.Close()
		***REMOVED***
		close(compressionDone)
	***REMOVED***()

	return pipeReader, compressionDone
***REMOVED***
