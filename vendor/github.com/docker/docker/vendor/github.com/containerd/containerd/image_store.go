package containerd

import (
	"context"

	imagesapi "github.com/containerd/containerd/api/services/images/v1"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	ptypes "github.com/gogo/protobuf/types"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type remoteImages struct ***REMOVED***
	client imagesapi.ImagesClient
***REMOVED***

// NewImageStoreFromClient returns a new image store client
func NewImageStoreFromClient(client imagesapi.ImagesClient) images.Store ***REMOVED***
	return &remoteImages***REMOVED***
		client: client,
	***REMOVED***
***REMOVED***

func (s *remoteImages) Get(ctx context.Context, name string) (images.Image, error) ***REMOVED***
	resp, err := s.client.Get(ctx, &imagesapi.GetImageRequest***REMOVED***
		Name: name,
	***REMOVED***)
	if err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	return imageFromProto(resp.Image), nil
***REMOVED***

func (s *remoteImages) List(ctx context.Context, filters ...string) ([]images.Image, error) ***REMOVED***
	resp, err := s.client.List(ctx, &imagesapi.ListImagesRequest***REMOVED***
		Filters: filters,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***

	return imagesFromProto(resp.Images), nil
***REMOVED***

func (s *remoteImages) Create(ctx context.Context, image images.Image) (images.Image, error) ***REMOVED***
	created, err := s.client.Create(ctx, &imagesapi.CreateImageRequest***REMOVED***
		Image: imageToProto(&image),
	***REMOVED***)
	if err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	return imageFromProto(&created.Image), nil
***REMOVED***

func (s *remoteImages) Update(ctx context.Context, image images.Image, fieldpaths ...string) (images.Image, error) ***REMOVED***
	var updateMask *ptypes.FieldMask
	if len(fieldpaths) > 0 ***REMOVED***
		updateMask = &ptypes.FieldMask***REMOVED***
			Paths: fieldpaths,
		***REMOVED***
	***REMOVED***

	updated, err := s.client.Update(ctx, &imagesapi.UpdateImageRequest***REMOVED***
		Image:      imageToProto(&image),
		UpdateMask: updateMask,
	***REMOVED***)
	if err != nil ***REMOVED***
		return images.Image***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***

	return imageFromProto(&updated.Image), nil
***REMOVED***

func (s *remoteImages) Delete(ctx context.Context, name string, opts ...images.DeleteOpt) error ***REMOVED***
	var do images.DeleteOptions
	for _, opt := range opts ***REMOVED***
		if err := opt(ctx, &do); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	_, err := s.client.Delete(ctx, &imagesapi.DeleteImageRequest***REMOVED***
		Name: name,
		Sync: do.Synchronous,
	***REMOVED***)

	return errdefs.FromGRPC(err)
***REMOVED***

func imageToProto(image *images.Image) imagesapi.Image ***REMOVED***
	return imagesapi.Image***REMOVED***
		Name:      image.Name,
		Labels:    image.Labels,
		Target:    descToProto(&image.Target),
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	***REMOVED***
***REMOVED***

func imageFromProto(imagepb *imagesapi.Image) images.Image ***REMOVED***
	return images.Image***REMOVED***
		Name:      imagepb.Name,
		Labels:    imagepb.Labels,
		Target:    descFromProto(&imagepb.Target),
		CreatedAt: imagepb.CreatedAt,
		UpdatedAt: imagepb.UpdatedAt,
	***REMOVED***
***REMOVED***

func imagesFromProto(imagespb []imagesapi.Image) []images.Image ***REMOVED***
	var images []images.Image

	for _, image := range imagespb ***REMOVED***
		images = append(images, imageFromProto(&image))
	***REMOVED***

	return images
***REMOVED***

func descFromProto(desc *types.Descriptor) ocispec.Descriptor ***REMOVED***
	return ocispec.Descriptor***REMOVED***
		MediaType: desc.MediaType,
		Size:      desc.Size_,
		Digest:    desc.Digest,
	***REMOVED***
***REMOVED***

func descToProto(desc *ocispec.Descriptor) types.Descriptor ***REMOVED***
	return types.Descriptor***REMOVED***
		MediaType: desc.MediaType,
		Size_:     desc.Size,
		Digest:    desc.Digest,
	***REMOVED***
***REMOVED***
