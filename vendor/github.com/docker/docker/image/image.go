package image

import (
	"encoding/json"
	"errors"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/layer"
	"github.com/opencontainers/go-digest"
)

// ID is the content-addressable ID of an image.
type ID digest.Digest

func (id ID) String() string ***REMOVED***
	return id.Digest().String()
***REMOVED***

// Digest converts ID into a digest
func (id ID) Digest() digest.Digest ***REMOVED***
	return digest.Digest(id)
***REMOVED***

// IDFromDigest creates an ID from a digest
func IDFromDigest(digest digest.Digest) ID ***REMOVED***
	return ID(digest)
***REMOVED***

// V1Image stores the V1 image configuration.
type V1Image struct ***REMOVED***
	// ID is a unique 64 character identifier of the image
	ID string `json:"id,omitempty"`
	// Parent is the ID of the parent image
	Parent string `json:"parent,omitempty"`
	// Comment is the commit message that was set when committing the image
	Comment string `json:"comment,omitempty"`
	// Created is the timestamp at which the image was created
	Created time.Time `json:"created"`
	// Container is the id of the container used to commit
	Container string `json:"container,omitempty"`
	// ContainerConfig is the configuration of the container that is committed into the image
	ContainerConfig container.Config `json:"container_config,omitempty"`
	// DockerVersion specifies the version of Docker that was used to build the image
	DockerVersion string `json:"docker_version,omitempty"`
	// Author is the name of the author that was specified when committing the image
	Author string `json:"author,omitempty"`
	// Config is the configuration of the container received from the client
	Config *container.Config `json:"config,omitempty"`
	// Architecture is the hardware that the image is built and runs on
	Architecture string `json:"architecture,omitempty"`
	// OS is the operating system used to build and run the image
	OS string `json:"os,omitempty"`
	// Size is the total size of the image including all layers it is composed of
	Size int64 `json:",omitempty"`
***REMOVED***

// Image stores the image configuration
type Image struct ***REMOVED***
	V1Image
	Parent     ID        `json:"parent,omitempty"`
	RootFS     *RootFS   `json:"rootfs,omitempty"`
	History    []History `json:"history,omitempty"`
	OSVersion  string    `json:"os.version,omitempty"`
	OSFeatures []string  `json:"os.features,omitempty"`

	// rawJSON caches the immutable JSON associated with this image.
	rawJSON []byte

	// computedID is the ID computed from the hash of the image config.
	// Not to be confused with the legacy V1 ID in V1Image.
	computedID ID
***REMOVED***

// RawJSON returns the immutable JSON associated with the image.
func (img *Image) RawJSON() []byte ***REMOVED***
	return img.rawJSON
***REMOVED***

// ID returns the image's content-addressable ID.
func (img *Image) ID() ID ***REMOVED***
	return img.computedID
***REMOVED***

// ImageID stringifies ID.
func (img *Image) ImageID() string ***REMOVED***
	return img.ID().String()
***REMOVED***

// RunConfig returns the image's container config.
func (img *Image) RunConfig() *container.Config ***REMOVED***
	return img.Config
***REMOVED***

// OperatingSystem returns the image's operating system. If not populated, defaults to the host runtime OS.
func (img *Image) OperatingSystem() string ***REMOVED***
	os := img.OS
	if os == "" ***REMOVED***
		os = runtime.GOOS
	***REMOVED***
	return os
***REMOVED***

// MarshalJSON serializes the image to JSON. It sorts the top-level keys so
// that JSON that's been manipulated by a push/pull cycle with a legacy
// registry won't end up with a different key order.
func (img *Image) MarshalJSON() ([]byte, error) ***REMOVED***
	type MarshalImage Image

	pass1, err := json.Marshal(MarshalImage(*img))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var c map[string]*json.RawMessage
	if err := json.Unmarshal(pass1, &c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return json.Marshal(c)
***REMOVED***

// ChildConfig is the configuration to apply to an Image to create a new
// Child image. Other properties of the image are copied from the parent.
type ChildConfig struct ***REMOVED***
	ContainerID     string
	Author          string
	Comment         string
	DiffID          layer.DiffID
	ContainerConfig *container.Config
	Config          *container.Config
***REMOVED***

// NewChildImage creates a new Image as a child of this image.
func NewChildImage(img *Image, child ChildConfig, platform string) *Image ***REMOVED***
	isEmptyLayer := layer.IsEmpty(child.DiffID)
	var rootFS *RootFS
	if img.RootFS != nil ***REMOVED***
		rootFS = img.RootFS.Clone()
	***REMOVED*** else ***REMOVED***
		rootFS = NewRootFS()
	***REMOVED***

	if !isEmptyLayer ***REMOVED***
		rootFS.Append(child.DiffID)
	***REMOVED***
	imgHistory := NewHistory(
		child.Author,
		child.Comment,
		strings.Join(child.ContainerConfig.Cmd, " "),
		isEmptyLayer)

	return &Image***REMOVED***
		V1Image: V1Image***REMOVED***
			DockerVersion:   dockerversion.Version,
			Config:          child.Config,
			Architecture:    runtime.GOARCH,
			OS:              platform,
			Container:       child.ContainerID,
			ContainerConfig: *child.ContainerConfig,
			Author:          child.Author,
			Created:         imgHistory.Created,
		***REMOVED***,
		RootFS:     rootFS,
		History:    append(img.History, imgHistory),
		OSFeatures: img.OSFeatures,
		OSVersion:  img.OSVersion,
	***REMOVED***
***REMOVED***

// History stores build commands that were used to create an image
type History struct ***REMOVED***
	// Created is the timestamp at which the image was created
	Created time.Time `json:"created"`
	// Author is the name of the author that was specified when committing the image
	Author string `json:"author,omitempty"`
	// CreatedBy keeps the Dockerfile command used while building the image
	CreatedBy string `json:"created_by,omitempty"`
	// Comment is the commit message that was set when committing the image
	Comment string `json:"comment,omitempty"`
	// EmptyLayer is set to true if this history item did not generate a
	// layer. Otherwise, the history item is associated with the next
	// layer in the RootFS section.
	EmptyLayer bool `json:"empty_layer,omitempty"`
***REMOVED***

// NewHistory creates a new history struct from arguments, and sets the created
// time to the current time in UTC
func NewHistory(author, comment, createdBy string, isEmptyLayer bool) History ***REMOVED***
	return History***REMOVED***
		Author:     author,
		Created:    time.Now().UTC(),
		CreatedBy:  createdBy,
		Comment:    comment,
		EmptyLayer: isEmptyLayer,
	***REMOVED***
***REMOVED***

// Exporter provides interface for loading and saving images
type Exporter interface ***REMOVED***
	Load(io.ReadCloser, io.Writer, bool) error
	// TODO: Load(net.Context, io.ReadCloser, <- chan StatusMessage) error
	Save([]string, io.Writer) error
***REMOVED***

// NewFromJSON creates an Image configuration from json.
func NewFromJSON(src []byte) (*Image, error) ***REMOVED***
	img := &Image***REMOVED******REMOVED***

	if err := json.Unmarshal(src, img); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if img.RootFS == nil ***REMOVED***
		return nil, errors.New("invalid image JSON, no RootFS key")
	***REMOVED***

	img.rawJSON = src

	return img, nil
***REMOVED***
