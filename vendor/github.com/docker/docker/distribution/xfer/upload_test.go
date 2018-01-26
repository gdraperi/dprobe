package xfer

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/progress"
	"golang.org/x/net/context"
)

const maxUploadConcurrency = 3

type mockUploadDescriptor struct ***REMOVED***
	currentUploads  *int32
	diffID          layer.DiffID
	simulateRetries int
***REMOVED***

// Key returns the key used to deduplicate downloads.
func (u *mockUploadDescriptor) Key() string ***REMOVED***
	return u.diffID.String()
***REMOVED***

// ID returns the ID for display purposes.
func (u *mockUploadDescriptor) ID() string ***REMOVED***
	return u.diffID.String()
***REMOVED***

// DiffID should return the DiffID for this layer.
func (u *mockUploadDescriptor) DiffID() layer.DiffID ***REMOVED***
	return u.diffID
***REMOVED***

// SetRemoteDescriptor is not used in the mock.
func (u *mockUploadDescriptor) SetRemoteDescriptor(remoteDescriptor distribution.Descriptor) ***REMOVED***
***REMOVED***

// Upload is called to perform the upload.
func (u *mockUploadDescriptor) Upload(ctx context.Context, progressOutput progress.Output) (distribution.Descriptor, error) ***REMOVED***
	if u.currentUploads != nil ***REMOVED***
		defer atomic.AddInt32(u.currentUploads, -1)

		if atomic.AddInt32(u.currentUploads, 1) > maxUploadConcurrency ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, errors.New("concurrency limit exceeded")
		***REMOVED***
	***REMOVED***

	// Sleep a bit to simulate a time-consuming upload.
	for i := int64(0); i <= 10; i++ ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return distribution.Descriptor***REMOVED******REMOVED***, ctx.Err()
		case <-time.After(10 * time.Millisecond):
			progressOutput.WriteProgress(progress.Progress***REMOVED***ID: u.ID(), Current: i, Total: 10***REMOVED***)
		***REMOVED***
	***REMOVED***

	if u.simulateRetries != 0 ***REMOVED***
		u.simulateRetries--
		return distribution.Descriptor***REMOVED******REMOVED***, errors.New("simulating retry")
	***REMOVED***

	return distribution.Descriptor***REMOVED******REMOVED***, nil
***REMOVED***

func uploadDescriptors(currentUploads *int32) []UploadDescriptor ***REMOVED***
	return []UploadDescriptor***REMOVED***
		&mockUploadDescriptor***REMOVED***currentUploads, layer.DiffID("sha256:cbbf2f9a99b47fc460d422812b6a5adff7dfee951d8fa2e4a98caa0382cfbdbf"), 0***REMOVED***,
		&mockUploadDescriptor***REMOVED***currentUploads, layer.DiffID("sha256:1515325234325236634634608943609283523908626098235490238423902343"), 0***REMOVED***,
		&mockUploadDescriptor***REMOVED***currentUploads, layer.DiffID("sha256:6929356290463485374960346430698374523437683470934634534953453453"), 0***REMOVED***,
		&mockUploadDescriptor***REMOVED***currentUploads, layer.DiffID("sha256:cbbf2f9a99b47fc460d422812b6a5adff7dfee951d8fa2e4a98caa0382cfbdbf"), 0***REMOVED***,
		&mockUploadDescriptor***REMOVED***currentUploads, layer.DiffID("sha256:8159352387436803946235346346368745389534789534897538734598734987"), 1***REMOVED***,
		&mockUploadDescriptor***REMOVED***currentUploads, layer.DiffID("sha256:4637863963478346897346987346987346789346789364879364897364987346"), 0***REMOVED***,
	***REMOVED***
***REMOVED***

func TestSuccessfulUpload(t *testing.T) ***REMOVED***
	lum := NewLayerUploadManager(maxUploadConcurrency, func(m *LayerUploadManager) ***REMOVED*** m.waitDuration = time.Millisecond ***REMOVED***)

	progressChan := make(chan progress.Progress)
	progressDone := make(chan struct***REMOVED******REMOVED***)
	receivedProgress := make(map[string]int64)

	go func() ***REMOVED***
		for p := range progressChan ***REMOVED***
			receivedProgress[p.ID] = p.Current
		***REMOVED***
		close(progressDone)
	***REMOVED***()

	var currentUploads int32
	descriptors := uploadDescriptors(&currentUploads)

	err := lum.Upload(context.Background(), descriptors, progress.ChanOutput(progressChan))
	if err != nil ***REMOVED***
		t.Fatalf("upload error: %v", err)
	***REMOVED***

	close(progressChan)
	<-progressDone
***REMOVED***

func TestCancelledUpload(t *testing.T) ***REMOVED***
	lum := NewLayerUploadManager(maxUploadConcurrency, func(m *LayerUploadManager) ***REMOVED*** m.waitDuration = time.Millisecond ***REMOVED***)

	progressChan := make(chan progress.Progress)
	progressDone := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		for range progressChan ***REMOVED***
		***REMOVED***
		close(progressDone)
	***REMOVED***()

	ctx, cancel := context.WithCancel(context.Background())

	go func() ***REMOVED***
		<-time.After(time.Millisecond)
		cancel()
	***REMOVED***()

	descriptors := uploadDescriptors(nil)
	err := lum.Upload(ctx, descriptors, progress.ChanOutput(progressChan))
	if err != context.Canceled ***REMOVED***
		t.Fatal("expected upload to be cancelled")
	***REMOVED***

	close(progressChan)
	<-progressDone
***REMOVED***
