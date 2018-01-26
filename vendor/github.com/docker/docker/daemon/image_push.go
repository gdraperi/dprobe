package daemon

import (
	"io"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/distribution"
	progressutils "github.com/docker/docker/distribution/utils"
	"github.com/docker/docker/pkg/progress"
	"golang.org/x/net/context"
)

// PushImage initiates a push operation on the repository named localName.
func (daemon *Daemon) PushImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error ***REMOVED***
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tag != "" ***REMOVED***
		// Push by digest is not supported, so only tags are supported.
		ref, err = reference.WithTag(ref, tag)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Include a buffer so that slow client connections don't affect
	// transfer performance.
	progressChan := make(chan progress.Progress, 100)

	writesDone := make(chan struct***REMOVED******REMOVED***)

	ctx, cancelFunc := context.WithCancel(ctx)

	go func() ***REMOVED***
		progressutils.WriteDistributionProgress(cancelFunc, outStream, progressChan)
		close(writesDone)
	***REMOVED***()

	imagePushConfig := &distribution.ImagePushConfig***REMOVED***
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
		ConfigMediaType: schema2.MediaTypeImageConfig,
		LayerStores:     distribution.NewLayerProvidersFromStores(daemon.layerStores),
		TrustKey:        daemon.trustKey,
		UploadManager:   daemon.uploadManager,
	***REMOVED***

	err = distribution.Push(ctx, ref, imagePushConfig)
	close(progressChan)
	<-writesDone
	return err
***REMOVED***
