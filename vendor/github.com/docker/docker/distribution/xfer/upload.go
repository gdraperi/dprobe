package xfer

import (
	"errors"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/progress"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const maxUploadAttempts = 5

// LayerUploadManager provides task management and progress reporting for
// uploads.
type LayerUploadManager struct ***REMOVED***
	tm           TransferManager
	waitDuration time.Duration
***REMOVED***

// SetConcurrency sets the max concurrent uploads for each push
func (lum *LayerUploadManager) SetConcurrency(concurrency int) ***REMOVED***
	lum.tm.SetConcurrency(concurrency)
***REMOVED***

// NewLayerUploadManager returns a new LayerUploadManager.
func NewLayerUploadManager(concurrencyLimit int, options ...func(*LayerUploadManager)) *LayerUploadManager ***REMOVED***
	manager := LayerUploadManager***REMOVED***
		tm:           NewTransferManager(concurrencyLimit),
		waitDuration: time.Second,
	***REMOVED***
	for _, option := range options ***REMOVED***
		option(&manager)
	***REMOVED***
	return &manager
***REMOVED***

type uploadTransfer struct ***REMOVED***
	Transfer

	remoteDescriptor distribution.Descriptor
	err              error
***REMOVED***

// An UploadDescriptor references a layer that may need to be uploaded.
type UploadDescriptor interface ***REMOVED***
	// Key returns the key used to deduplicate uploads.
	Key() string
	// ID returns the ID for display purposes.
	ID() string
	// DiffID should return the DiffID for this layer.
	DiffID() layer.DiffID
	// Upload is called to perform the Upload.
	Upload(ctx context.Context, progressOutput progress.Output) (distribution.Descriptor, error)
	// SetRemoteDescriptor provides the distribution.Descriptor that was
	// returned by Upload. This descriptor is not to be confused with
	// the UploadDescriptor interface, which is used for internally
	// identifying layers that are being uploaded.
	SetRemoteDescriptor(descriptor distribution.Descriptor)
***REMOVED***

// Upload is a blocking function which ensures the listed layers are present on
// the remote registry. It uses the string returned by the Key method to
// deduplicate uploads.
func (lum *LayerUploadManager) Upload(ctx context.Context, layers []UploadDescriptor, progressOutput progress.Output) error ***REMOVED***
	var (
		uploads          []*uploadTransfer
		dedupDescriptors = make(map[string]*uploadTransfer)
	)

	for _, descriptor := range layers ***REMOVED***
		progress.Update(progressOutput, descriptor.ID(), "Preparing")

		key := descriptor.Key()
		if _, present := dedupDescriptors[key]; present ***REMOVED***
			continue
		***REMOVED***

		xferFunc := lum.makeUploadFunc(descriptor)
		upload, watcher := lum.tm.Transfer(descriptor.Key(), xferFunc, progressOutput)
		defer upload.Release(watcher)
		uploads = append(uploads, upload.(*uploadTransfer))
		dedupDescriptors[key] = upload.(*uploadTransfer)
	***REMOVED***

	for _, upload := range uploads ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-upload.Transfer.Done():
			if upload.err != nil ***REMOVED***
				return upload.err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, l := range layers ***REMOVED***
		l.SetRemoteDescriptor(dedupDescriptors[l.Key()].remoteDescriptor)
	***REMOVED***

	return nil
***REMOVED***

func (lum *LayerUploadManager) makeUploadFunc(descriptor UploadDescriptor) DoFunc ***REMOVED***
	return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
		u := &uploadTransfer***REMOVED***
			Transfer: NewTransfer(),
		***REMOVED***

		go func() ***REMOVED***
			defer func() ***REMOVED***
				close(progressChan)
			***REMOVED***()

			progressOutput := progress.ChanOutput(progressChan)

			select ***REMOVED***
			case <-start:
			default:
				progress.Update(progressOutput, descriptor.ID(), "Waiting")
				<-start
			***REMOVED***

			retries := 0
			for ***REMOVED***
				remoteDescriptor, err := descriptor.Upload(u.Transfer.Context(), progressOutput)
				if err == nil ***REMOVED***
					u.remoteDescriptor = remoteDescriptor
					break
				***REMOVED***

				// If an error was returned because the context
				// was cancelled, we shouldn't retry.
				select ***REMOVED***
				case <-u.Transfer.Context().Done():
					u.err = err
					return
				default:
				***REMOVED***

				retries++
				if _, isDNR := err.(DoNotRetry); isDNR || retries == maxUploadAttempts ***REMOVED***
					logrus.Errorf("Upload failed: %v", err)
					u.err = err
					return
				***REMOVED***

				logrus.Errorf("Upload failed, retrying: %v", err)
				delay := retries * 5
				ticker := time.NewTicker(lum.waitDuration)

			selectLoop:
				for ***REMOVED***
					progress.Updatef(progressOutput, descriptor.ID(), "Retrying in %d second%s", delay, (map[bool]string***REMOVED***true: "s"***REMOVED***)[delay != 1])
					select ***REMOVED***
					case <-ticker.C:
						delay--
						if delay == 0 ***REMOVED***
							ticker.Stop()
							break selectLoop
						***REMOVED***
					case <-u.Transfer.Context().Done():
						ticker.Stop()
						u.err = errors.New("upload cancelled during retry delay")
						return
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***()

		return u
	***REMOVED***
***REMOVED***
