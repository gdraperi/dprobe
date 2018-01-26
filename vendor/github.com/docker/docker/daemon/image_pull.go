package daemon

import (
	"io"
	"runtime"
	"strings"

	dist "github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/distribution"
	progressutils "github.com/docker/docker/distribution/utils"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/registry"
	"github.com/opencontainers/go-digest"
	"golang.org/x/net/context"
)

// PullImage initiates a pull operation. image is the repository name to pull, and
// tag may be either empty, or indicate a specific tag to pull.
func (daemon *Daemon) PullImage(ctx context.Context, image, tag, os string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error ***REMOVED***
	// Special case: "pull -a" may send an image name with a
	// trailing :. This is ugly, but let's not break API
	// compatibility.
	image = strings.TrimSuffix(image, ":")

	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	if tag != "" ***REMOVED***
		// The "tag" could actually be a digest.
		var dgst digest.Digest
		dgst, err = digest.Parse(tag)
		if err == nil ***REMOVED***
			ref, err = reference.WithDigest(reference.TrimNamed(ref), dgst)
		***REMOVED*** else ***REMOVED***
			ref, err = reference.WithTag(ref, tag)
		***REMOVED***
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***
	***REMOVED***

	return daemon.pullImageWithReference(ctx, ref, os, metaHeaders, authConfig, outStream)
***REMOVED***

func (daemon *Daemon) pullImageWithReference(ctx context.Context, ref reference.Named, os string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error ***REMOVED***
	// Include a buffer so that slow client connections don't affect
	// transfer performance.
	progressChan := make(chan progress.Progress, 100)

	writesDone := make(chan struct***REMOVED******REMOVED***)

	ctx, cancelFunc := context.WithCancel(ctx)

	go func() ***REMOVED***
		progressutils.WriteDistributionProgress(cancelFunc, outStream, progressChan)
		close(writesDone)
	***REMOVED***()

	// Default to the host OS platform in case it hasn't been populated with an explicit value.
	if os == "" ***REMOVED***
		os = runtime.GOOS
	***REMOVED***

	imagePullConfig := &distribution.ImagePullConfig***REMOVED***
		Config: distribution.Config***REMOVED***
			MetaHeaders:      metaHeaders,
			AuthConfig:       authConfig,
			ProgressOutput:   progress.ChanOutput(progressChan),
			RegistryService:  daemon.RegistryService,
			ImageEventLogger: daemon.LogImageEvent,
			MetadataStore:    daemon.distributionMetadataStore,
			ImageStore:       distribution.NewImageConfigStoreFromStore(daemon.imageStore),
			ReferenceStore:   daemon.referenceStore,
		***REMOVED***,
		DownloadManager: daemon.downloadManager,
		Schema2Types:    distribution.ImageTypes,
		OS:              os,
	***REMOVED***

	err := distribution.Pull(ctx, ref, imagePullConfig)
	close(progressChan)
	<-writesDone
	return err
***REMOVED***

// GetRepository returns a repository from the registry.
func (daemon *Daemon) GetRepository(ctx context.Context, ref reference.Named, authConfig *types.AuthConfig) (dist.Repository, bool, error) ***REMOVED***
	// get repository info
	repoInfo, err := daemon.RegistryService.ResolveRepository(ref)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	// makes sure name is not empty or `scratch`
	if err := distribution.ValidateRepoName(repoInfo.Name); err != nil ***REMOVED***
		return nil, false, errdefs.InvalidParameter(err)
	***REMOVED***

	// get endpoints
	endpoints, err := daemon.RegistryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***

	// retrieve repository
	var (
		confirmedV2 bool
		repository  dist.Repository
		lastError   error
	)

	for _, endpoint := range endpoints ***REMOVED***
		if endpoint.Version == registry.APIVersion1 ***REMOVED***
			continue
		***REMOVED***

		repository, confirmedV2, lastError = distribution.NewV2Repository(ctx, repoInfo, endpoint, nil, authConfig, "pull")
		if lastError == nil && confirmedV2 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return repository, confirmedV2, lastError
***REMOVED***
