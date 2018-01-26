package metadata

import (
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	"github.com/pkg/errors"
)

// V1IDService maps v1 IDs to layers on disk.
type V1IDService struct ***REMOVED***
	store Store
***REMOVED***

// NewV1IDService creates a new V1 ID mapping service.
func NewV1IDService(store Store) *V1IDService ***REMOVED***
	return &V1IDService***REMOVED***
		store: store,
	***REMOVED***
***REMOVED***

// namespace returns the namespace used by this service.
func (idserv *V1IDService) namespace() string ***REMOVED***
	return "v1id"
***REMOVED***

// Get finds a layer by its V1 ID.
func (idserv *V1IDService) Get(v1ID, registry string) (layer.DiffID, error) ***REMOVED***
	if idserv.store == nil ***REMOVED***
		return "", errors.New("no v1IDService storage")
	***REMOVED***
	if err := v1.ValidateID(v1ID); err != nil ***REMOVED***
		return layer.DiffID(""), err
	***REMOVED***

	idBytes, err := idserv.store.Get(idserv.namespace(), registry+","+v1ID)
	if err != nil ***REMOVED***
		return layer.DiffID(""), err
	***REMOVED***
	return layer.DiffID(idBytes), nil
***REMOVED***

// Set associates an image with a V1 ID.
func (idserv *V1IDService) Set(v1ID, registry string, id layer.DiffID) error ***REMOVED***
	if idserv.store == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := v1.ValidateID(v1ID); err != nil ***REMOVED***
		return err
	***REMOVED***
	return idserv.store.Set(idserv.namespace(), registry+","+v1ID, []byte(id))
***REMOVED***
