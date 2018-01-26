package distribution

import (
	"fmt"
	"sync"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/image"
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/registry"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type v1Pusher struct ***REMOVED***
	v1IDService *metadata.V1IDService
	endpoint    registry.APIEndpoint
	ref         reference.Named
	repoInfo    *registry.RepositoryInfo
	config      *ImagePushConfig
	session     *registry.Session
***REMOVED***

func (p *v1Pusher) Push(ctx context.Context) error ***REMOVED***
	tlsConfig, err := p.config.RegistryService.TLSConfig(p.repoInfo.Index.Name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Adds Docker-specific headers as well as user-specified headers (metaHeaders)
	tr := transport.NewTransport(
		// TODO(tiborvass): was NoTimeout
		registry.NewTransport(tlsConfig),
		registry.Headers(dockerversion.DockerUserAgent(ctx), p.config.MetaHeaders)...,
	)
	client := registry.HTTPClient(tr)
	v1Endpoint := p.endpoint.ToV1Endpoint(dockerversion.DockerUserAgent(ctx), p.config.MetaHeaders)
	p.session, err = registry.NewSession(client, p.config.AuthConfig, v1Endpoint)
	if err != nil ***REMOVED***
		// TODO(dmcgowan): Check if should fallback
		return fallbackError***REMOVED***err: err***REMOVED***
	***REMOVED***
	if err := p.pushRepository(ctx); err != nil ***REMOVED***
		// TODO(dmcgowan): Check if should fallback
		return err
	***REMOVED***
	return nil
***REMOVED***

// v1Image exposes the configuration, filesystem layer ID, and a v1 ID for an
// image being pushed to a v1 registry.
type v1Image interface ***REMOVED***
	Config() []byte
	Layer() layer.Layer
	V1ID() string
***REMOVED***

type v1ImageCommon struct ***REMOVED***
	layer  layer.Layer
	config []byte
	v1ID   string
***REMOVED***

func (common *v1ImageCommon) Config() []byte ***REMOVED***
	return common.config
***REMOVED***

func (common *v1ImageCommon) V1ID() string ***REMOVED***
	return common.v1ID
***REMOVED***

func (common *v1ImageCommon) Layer() layer.Layer ***REMOVED***
	return common.layer
***REMOVED***

// v1TopImage defines a runnable (top layer) image being pushed to a v1
// registry.
type v1TopImage struct ***REMOVED***
	v1ImageCommon
	imageID image.ID
***REMOVED***

func newV1TopImage(imageID image.ID, img *image.Image, l layer.Layer, parent *v1DependencyImage) (*v1TopImage, error) ***REMOVED***
	v1ID := imageID.Digest().Hex()
	parentV1ID := ""
	if parent != nil ***REMOVED***
		parentV1ID = parent.V1ID()
	***REMOVED***

	config, err := v1.MakeV1ConfigFromConfig(img, v1ID, parentV1ID, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &v1TopImage***REMOVED***
		v1ImageCommon: v1ImageCommon***REMOVED***
			v1ID:   v1ID,
			config: config,
			layer:  l,
		***REMOVED***,
		imageID: imageID,
	***REMOVED***, nil
***REMOVED***

// v1DependencyImage defines a dependency layer being pushed to a v1 registry.
type v1DependencyImage struct ***REMOVED***
	v1ImageCommon
***REMOVED***

func newV1DependencyImage(l layer.Layer, parent *v1DependencyImage) *v1DependencyImage ***REMOVED***
	v1ID := digest.Digest(l.ChainID()).Hex()

	var config string
	if parent != nil ***REMOVED***
		config = fmt.Sprintf(`***REMOVED***"id":"%s","parent":"%s"***REMOVED***`, v1ID, parent.V1ID())
	***REMOVED*** else ***REMOVED***
		config = fmt.Sprintf(`***REMOVED***"id":"%s"***REMOVED***`, v1ID)
	***REMOVED***
	return &v1DependencyImage***REMOVED***
		v1ImageCommon: v1ImageCommon***REMOVED***
			v1ID:   v1ID,
			config: []byte(config),
			layer:  l,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Retrieve the all the images to be uploaded in the correct order
func (p *v1Pusher) getImageList() (imageList []v1Image, tagsByImage map[image.ID][]string, referencedLayers []PushLayer, err error) ***REMOVED***
	tagsByImage = make(map[image.ID][]string)

	// Ignore digest references
	if _, isCanonical := p.ref.(reference.Canonical); isCanonical ***REMOVED***
		return
	***REMOVED***

	tagged, isTagged := p.ref.(reference.NamedTagged)
	if isTagged ***REMOVED***
		// Push a specific tag
		var imgID image.ID
		var dgst digest.Digest
		dgst, err = p.config.ReferenceStore.Get(p.ref)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		imgID = image.IDFromDigest(dgst)

		imageList, err = p.imageListForTag(imgID, nil, &referencedLayers)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		tagsByImage[imgID] = []string***REMOVED***tagged.Tag()***REMOVED***

		return
	***REMOVED***

	imagesSeen := make(map[digest.Digest]struct***REMOVED******REMOVED***)
	dependenciesSeen := make(map[layer.ChainID]*v1DependencyImage)

	associations := p.config.ReferenceStore.ReferencesByName(p.ref)
	for _, association := range associations ***REMOVED***
		if tagged, isTagged = association.Ref.(reference.NamedTagged); !isTagged ***REMOVED***
			// Ignore digest references.
			continue
		***REMOVED***

		imgID := image.IDFromDigest(association.ID)
		tagsByImage[imgID] = append(tagsByImage[imgID], tagged.Tag())

		if _, present := imagesSeen[association.ID]; present ***REMOVED***
			// Skip generating image list for already-seen image
			continue
		***REMOVED***
		imagesSeen[association.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

		imageListForThisTag, err := p.imageListForTag(imgID, dependenciesSeen, &referencedLayers)
		if err != nil ***REMOVED***
			return nil, nil, nil, err
		***REMOVED***

		// append to main image list
		imageList = append(imageList, imageListForThisTag...)
	***REMOVED***
	if len(imageList) == 0 ***REMOVED***
		return nil, nil, nil, fmt.Errorf("No images found for the requested repository / tag")
	***REMOVED***
	logrus.Debugf("Image list: %v", imageList)
	logrus.Debugf("Tags by image: %v", tagsByImage)

	return
***REMOVED***

func (p *v1Pusher) imageListForTag(imgID image.ID, dependenciesSeen map[layer.ChainID]*v1DependencyImage, referencedLayers *[]PushLayer) (imageListForThisTag []v1Image, err error) ***REMOVED***
	ics, ok := p.config.ImageStore.(*imageConfigStore)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("only image store images supported for v1 push")
	***REMOVED***
	img, err := ics.Store.Get(imgID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	topLayerID := img.RootFS.ChainID()

	if !system.IsOSSupported(img.OperatingSystem()) ***REMOVED***
		return nil, system.ErrNotSupportedOperatingSystem
	***REMOVED***
	pl, err := p.config.LayerStores[img.OperatingSystem()].Get(topLayerID)
	*referencedLayers = append(*referencedLayers, pl)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get top layer from image: %v", err)
	***REMOVED***

	// V1 push is deprecated, only support existing layerstore layers
	lsl, ok := pl.(*storeLayer)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("only layer store layers supported for v1 push")
	***REMOVED***
	l := lsl.Layer

	dependencyImages, parent := generateDependencyImages(l.Parent(), dependenciesSeen)

	topImage, err := newV1TopImage(imgID, img, l, parent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	imageListForThisTag = append(dependencyImages, topImage)

	return
***REMOVED***

func generateDependencyImages(l layer.Layer, dependenciesSeen map[layer.ChainID]*v1DependencyImage) (imageListForThisTag []v1Image, parent *v1DependencyImage) ***REMOVED***
	if l == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	imageListForThisTag, parent = generateDependencyImages(l.Parent(), dependenciesSeen)

	if dependenciesSeen != nil ***REMOVED***
		if dependencyImage, present := dependenciesSeen[l.ChainID()]; present ***REMOVED***
			// This layer is already on the list, we can ignore it
			// and all its parents.
			return imageListForThisTag, dependencyImage
		***REMOVED***
	***REMOVED***

	dependencyImage := newV1DependencyImage(l, parent)
	imageListForThisTag = append(imageListForThisTag, dependencyImage)

	if dependenciesSeen != nil ***REMOVED***
		dependenciesSeen[l.ChainID()] = dependencyImage
	***REMOVED***

	return imageListForThisTag, dependencyImage
***REMOVED***

// createImageIndex returns an index of an image's layer IDs and tags.
func createImageIndex(images []v1Image, tags map[image.ID][]string) []*registry.ImgData ***REMOVED***
	var imageIndex []*registry.ImgData
	for _, img := range images ***REMOVED***
		v1ID := img.V1ID()

		if topImage, isTopImage := img.(*v1TopImage); isTopImage ***REMOVED***
			if tags, hasTags := tags[topImage.imageID]; hasTags ***REMOVED***
				// If an image has tags you must add an entry in the image index
				// for each tag
				for _, tag := range tags ***REMOVED***
					imageIndex = append(imageIndex, &registry.ImgData***REMOVED***
						ID:  v1ID,
						Tag: tag,
					***REMOVED***)
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		// If the image does not have a tag it still needs to be sent to the
		// registry with an empty tag so that it is associated with the repository
		imageIndex = append(imageIndex, &registry.ImgData***REMOVED***
			ID:  v1ID,
			Tag: "",
		***REMOVED***)
	***REMOVED***
	return imageIndex
***REMOVED***

// lookupImageOnEndpoint checks the specified endpoint to see if an image exists
// and if it is absent then it sends the image id to the channel to be pushed.
func (p *v1Pusher) lookupImageOnEndpoint(wg *sync.WaitGroup, endpoint string, images chan v1Image, imagesToPush chan string) ***REMOVED***
	defer wg.Done()
	for image := range images ***REMOVED***
		v1ID := image.V1ID()
		truncID := stringid.TruncateID(image.Layer().DiffID().String())
		if err := p.session.LookupRemoteImage(v1ID, endpoint); err != nil ***REMOVED***
			logrus.Errorf("Error in LookupRemoteImage: %s", err)
			imagesToPush <- v1ID
			progress.Update(p.config.ProgressOutput, truncID, "Waiting")
		***REMOVED*** else ***REMOVED***
			progress.Update(p.config.ProgressOutput, truncID, "Already exists")
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *v1Pusher) pushImageToEndpoint(ctx context.Context, endpoint string, imageList []v1Image, tags map[image.ID][]string, repo *registry.RepositoryData) error ***REMOVED***
	workerCount := len(imageList)
	// start a maximum of 5 workers to check if images exist on the specified endpoint.
	if workerCount > 5 ***REMOVED***
		workerCount = 5
	***REMOVED***
	var (
		wg           = &sync.WaitGroup***REMOVED******REMOVED***
		imageData    = make(chan v1Image, workerCount*2)
		imagesToPush = make(chan string, workerCount*2)
		pushes       = make(chan map[string]struct***REMOVED******REMOVED***, 1)
	)
	for i := 0; i < workerCount; i++ ***REMOVED***
		wg.Add(1)
		go p.lookupImageOnEndpoint(wg, endpoint, imageData, imagesToPush)
	***REMOVED***
	// start a go routine that consumes the images to push
	go func() ***REMOVED***
		shouldPush := make(map[string]struct***REMOVED******REMOVED***)
		for id := range imagesToPush ***REMOVED***
			shouldPush[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
		pushes <- shouldPush
	***REMOVED***()
	for _, v1Image := range imageList ***REMOVED***
		imageData <- v1Image
	***REMOVED***
	// close the channel to notify the workers that there will be no more images to check.
	close(imageData)
	wg.Wait()
	close(imagesToPush)
	// wait for all the images that require pushes to be collected into a consumable map.
	shouldPush := <-pushes
	// finish by pushing any images and tags to the endpoint.  The order that the images are pushed
	// is very important that is why we are still iterating over the ordered list of imageIDs.
	for _, img := range imageList ***REMOVED***
		v1ID := img.V1ID()
		if _, push := shouldPush[v1ID]; push ***REMOVED***
			if _, err := p.pushImage(ctx, img, endpoint); err != nil ***REMOVED***
				// FIXME: Continue on error?
				return err
			***REMOVED***
		***REMOVED***
		if topImage, isTopImage := img.(*v1TopImage); isTopImage ***REMOVED***
			for _, tag := range tags[topImage.imageID] ***REMOVED***
				progress.Messagef(p.config.ProgressOutput, "", "Pushing tag for rev [%s] on ***REMOVED***%s***REMOVED***", stringid.TruncateID(v1ID), endpoint+"repositories/"+reference.Path(p.repoInfo.Name)+"/tags/"+tag)
				if err := p.session.PushRegistryTag(p.repoInfo.Name, v1ID, tag, endpoint); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// pushRepository pushes layers that do not already exist on the registry.
func (p *v1Pusher) pushRepository(ctx context.Context) error ***REMOVED***
	imgList, tags, referencedLayers, err := p.getImageList()
	defer func() ***REMOVED***
		for _, l := range referencedLayers ***REMOVED***
			l.Release()
		***REMOVED***
	***REMOVED***()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	imageIndex := createImageIndex(imgList, tags)
	for _, data := range imageIndex ***REMOVED***
		logrus.Debugf("Pushing ID: %s with Tag: %s", data.ID, data.Tag)
	***REMOVED***

	// Register all the images in a repository with the registry
	// If an image is not in this list it will not be associated with the repository
	repoData, err := p.session.PushImageJSONIndex(p.repoInfo.Name, imageIndex, false, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// push the repository to each of the endpoints only if it does not exist.
	for _, endpoint := range repoData.Endpoints ***REMOVED***
		if err := p.pushImageToEndpoint(ctx, endpoint, imgList, tags, repoData); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	_, err = p.session.PushImageJSONIndex(p.repoInfo.Name, imageIndex, true, repoData.Endpoints)
	return err
***REMOVED***

func (p *v1Pusher) pushImage(ctx context.Context, v1Image v1Image, ep string) (checksum string, err error) ***REMOVED***
	l := v1Image.Layer()
	v1ID := v1Image.V1ID()
	truncID := stringid.TruncateID(l.DiffID().String())

	jsonRaw := v1Image.Config()
	progress.Update(p.config.ProgressOutput, truncID, "Pushing")

	// General rule is to use ID for graph accesses and compatibilityID for
	// calls to session.registry()
	imgData := &registry.ImgData***REMOVED***
		ID: v1ID,
	***REMOVED***

	// Send the json
	if err := p.session.PushImageJSONRegistry(imgData, jsonRaw, ep); err != nil ***REMOVED***
		if err == registry.ErrAlreadyExists ***REMOVED***
			progress.Update(p.config.ProgressOutput, truncID, "Image already pushed, skipping")
			return "", nil
		***REMOVED***
		return "", err
	***REMOVED***

	arch, err := l.TarStream()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer arch.Close()

	// don't care if this fails; best effort
	size, _ := l.DiffSize()

	// Send the layer
	logrus.Debugf("rendered layer for %s of [%d] size", v1ID, size)

	reader := progress.NewProgressReader(ioutils.NewCancelReadCloser(ctx, arch), p.config.ProgressOutput, size, truncID, "Pushing")
	defer reader.Close()

	checksum, checksumPayload, err := p.session.PushImageLayerRegistry(v1ID, reader, ep, jsonRaw)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	imgData.Checksum = checksum
	imgData.ChecksumPayload = checksumPayload
	// Send the checksum
	if err := p.session.PushImageChecksumRegistry(imgData, ep); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err := p.v1IDService.Set(v1ID, p.repoInfo.Index.Name, l.DiffID()); err != nil ***REMOVED***
		logrus.Warnf("Could not set v1 ID mapping: %v", err)
	***REMOVED***

	progress.Update(p.config.ProgressOutput, truncID, "Image successfully pushed")
	return imgData.Checksum, nil
***REMOVED***
