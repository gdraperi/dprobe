package layer

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// DigestSHA256EmptyTar is the canonical sha256 digest of empty tar file -
// (1024 NULL bytes)
const DigestSHA256EmptyTar = DiffID("sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef")

type emptyLayer struct***REMOVED******REMOVED***

// EmptyLayer is a layer that corresponds to empty tar.
var EmptyLayer = &emptyLayer***REMOVED******REMOVED***

func (el *emptyLayer) TarStream() (io.ReadCloser, error) ***REMOVED***
	buf := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)
	tarWriter.Close()
	return ioutil.NopCloser(buf), nil
***REMOVED***

func (el *emptyLayer) TarStreamFrom(p ChainID) (io.ReadCloser, error) ***REMOVED***
	if p == "" ***REMOVED***
		return el.TarStream()
	***REMOVED***
	return nil, fmt.Errorf("can't get parent tar stream of an empty layer")
***REMOVED***

func (el *emptyLayer) ChainID() ChainID ***REMOVED***
	return ChainID(DigestSHA256EmptyTar)
***REMOVED***

func (el *emptyLayer) DiffID() DiffID ***REMOVED***
	return DigestSHA256EmptyTar
***REMOVED***

func (el *emptyLayer) Parent() Layer ***REMOVED***
	return nil
***REMOVED***

func (el *emptyLayer) Size() (size int64, err error) ***REMOVED***
	return 0, nil
***REMOVED***

func (el *emptyLayer) DiffSize() (size int64, err error) ***REMOVED***
	return 0, nil
***REMOVED***

func (el *emptyLayer) Metadata() (map[string]string, error) ***REMOVED***
	return make(map[string]string), nil
***REMOVED***

// IsEmpty returns true if the layer is an EmptyLayer
func IsEmpty(diffID DiffID) bool ***REMOVED***
	return diffID == DigestSHA256EmptyTar
***REMOVED***
