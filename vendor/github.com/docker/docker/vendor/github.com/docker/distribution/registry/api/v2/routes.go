package v2

import "github.com/gorilla/mux"

// The following are definitions of the name under which all V2 routes are
// registered. These symbols can be used to look up a route based on the name.
const (
	RouteNameBase            = "base"
	RouteNameManifest        = "manifest"
	RouteNameTags            = "tags"
	RouteNameBlob            = "blob"
	RouteNameBlobUpload      = "blob-upload"
	RouteNameBlobUploadChunk = "blob-upload-chunk"
	RouteNameCatalog         = "catalog"
)

var allEndpoints = []string***REMOVED***
	RouteNameManifest,
	RouteNameCatalog,
	RouteNameTags,
	RouteNameBlob,
	RouteNameBlobUpload,
	RouteNameBlobUploadChunk,
***REMOVED***

// Router builds a gorilla router with named routes for the various API
// methods. This can be used directly by both server implementations and
// clients.
func Router() *mux.Router ***REMOVED***
	return RouterWithPrefix("")
***REMOVED***

// RouterWithPrefix builds a gorilla router with a configured prefix
// on all routes.
func RouterWithPrefix(prefix string) *mux.Router ***REMOVED***
	rootRouter := mux.NewRouter()
	router := rootRouter
	if prefix != "" ***REMOVED***
		router = router.PathPrefix(prefix).Subrouter()
	***REMOVED***

	router.StrictSlash(true)

	for _, descriptor := range routeDescriptors ***REMOVED***
		router.Path(descriptor.Path).Name(descriptor.Name)
	***REMOVED***

	return rootRouter
***REMOVED***
