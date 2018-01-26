package layer

import (
	"errors"
)

// Getter is an interface to get the path to a layer on the host.
type Getter interface ***REMOVED***
	// GetLayerPath gets the path for the layer. This is different from Get()
	// since that returns an interface to account for umountable layers.
	GetLayerPath(id string) (string, error)
***REMOVED***

// GetLayerPath returns the path to a layer
func GetLayerPath(s Store, layer ChainID) (string, error) ***REMOVED***
	ls, ok := s.(*layerStore)
	if !ok ***REMOVED***
		return "", errors.New("unsupported layer store")
	***REMOVED***
	ls.layerL.Lock()
	defer ls.layerL.Unlock()

	rl, ok := ls.layerMap[layer]
	if !ok ***REMOVED***
		return "", ErrLayerDoesNotExist
	***REMOVED***

	if layerGetter, ok := ls.driver.(Getter); ok ***REMOVED***
		return layerGetter.GetLayerPath(rl.cacheID)
	***REMOVED***
	path, err := ls.driver.Get(rl.cacheID, "")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err := ls.driver.Put(rl.cacheID); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return path.Path(), nil
***REMOVED***

func (ls *layerStore) mountID(name string) string ***REMOVED***
	// windows has issues if container ID doesn't match mount ID
	return name
***REMOVED***
