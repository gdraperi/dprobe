package tarexport

import (
	"github.com/docker/distribution"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	refstore "github.com/docker/docker/reference"
)

const (
	manifestFileName           = "manifest.json"
	legacyLayerFileName        = "layer.tar"
	legacyConfigFileName       = "json"
	legacyVersionFileName      = "VERSION"
	legacyRepositoriesFileName = "repositories"
)

type manifestItem struct ***REMOVED***
	Config       string
	RepoTags     []string
	Layers       []string
	Parent       image.ID                                 `json:",omitempty"`
	LayerSources map[layer.DiffID]distribution.Descriptor `json:",omitempty"`
***REMOVED***

type tarexporter struct ***REMOVED***
	is             image.Store
	lss            map[string]layer.Store
	rs             refstore.Store
	loggerImgEvent LogImageEvent
***REMOVED***

// LogImageEvent defines interface for event generation related to image tar(load and save) operations
type LogImageEvent interface ***REMOVED***
	//LogImageEvent generates an event related to an image operation
	LogImageEvent(imageID, refName, action string)
***REMOVED***

// NewTarExporter returns new Exporter for tar packages
func NewTarExporter(is image.Store, lss map[string]layer.Store, rs refstore.Store, loggerImgEvent LogImageEvent) image.Exporter ***REMOVED***
	return &tarexporter***REMOVED***
		is:             is,
		lss:            lss,
		rs:             rs,
		loggerImgEvent: loggerImgEvent,
	***REMOVED***
***REMOVED***
