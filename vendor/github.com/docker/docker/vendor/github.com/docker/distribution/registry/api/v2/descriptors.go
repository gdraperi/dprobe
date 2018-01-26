package v2

import (
	"net/http"
	"regexp"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/opencontainers/go-digest"
)

var (
	nameParameterDescriptor = ParameterDescriptor***REMOVED***
		Name:        "name",
		Type:        "string",
		Format:      reference.NameRegexp.String(),
		Required:    true,
		Description: `Name of the target repository.`,
	***REMOVED***

	referenceParameterDescriptor = ParameterDescriptor***REMOVED***
		Name:        "reference",
		Type:        "string",
		Format:      reference.TagRegexp.String(),
		Required:    true,
		Description: `Tag or digest of the target manifest.`,
	***REMOVED***

	uuidParameterDescriptor = ParameterDescriptor***REMOVED***
		Name:        "uuid",
		Type:        "opaque",
		Required:    true,
		Description: "A uuid identifying the upload. This field can accept characters that match `[a-zA-Z0-9-_.=]+`.",
	***REMOVED***

	digestPathParameter = ParameterDescriptor***REMOVED***
		Name:        "digest",
		Type:        "path",
		Required:    true,
		Format:      digest.DigestRegexp.String(),
		Description: `Digest of desired blob.`,
	***REMOVED***

	hostHeader = ParameterDescriptor***REMOVED***
		Name:        "Host",
		Type:        "string",
		Description: "Standard HTTP Host Header. Should be set to the registry host.",
		Format:      "<registry host>",
		Examples:    []string***REMOVED***"registry-1.docker.io"***REMOVED***,
	***REMOVED***

	authHeader = ParameterDescriptor***REMOVED***
		Name:        "Authorization",
		Type:        "string",
		Description: "An RFC7235 compliant authorization header.",
		Format:      "<scheme> <token>",
		Examples:    []string***REMOVED***"Bearer dGhpcyBpcyBhIGZha2UgYmVhcmVyIHRva2VuIQ=="***REMOVED***,
	***REMOVED***

	authChallengeHeader = ParameterDescriptor***REMOVED***
		Name:        "WWW-Authenticate",
		Type:        "string",
		Description: "An RFC7235 compliant authentication challenge header.",
		Format:      `<scheme> realm="<realm>", ..."`,
		Examples: []string***REMOVED***
			`Bearer realm="https://auth.docker.com/", service="registry.docker.com", scopes="repository:library/ubuntu:pull"`,
		***REMOVED***,
	***REMOVED***

	contentLengthZeroHeader = ParameterDescriptor***REMOVED***
		Name:        "Content-Length",
		Description: "The `Content-Length` header must be zero and the body must be empty.",
		Type:        "integer",
		Format:      "0",
	***REMOVED***

	dockerUploadUUIDHeader = ParameterDescriptor***REMOVED***
		Name:        "Docker-Upload-UUID",
		Description: "Identifies the docker upload uuid for the current request.",
		Type:        "uuid",
		Format:      "<uuid>",
	***REMOVED***

	digestHeader = ParameterDescriptor***REMOVED***
		Name:        "Docker-Content-Digest",
		Description: "Digest of the targeted content for the request.",
		Type:        "digest",
		Format:      "<digest>",
	***REMOVED***

	linkHeader = ParameterDescriptor***REMOVED***
		Name:        "Link",
		Type:        "link",
		Description: "RFC5988 compliant rel='next' with URL to next result set, if available",
		Format:      `<<url>?n=<last n value>&last=<last entry from response>>; rel="next"`,
	***REMOVED***

	paginationParameters = []ParameterDescriptor***REMOVED***
		***REMOVED***
			Name:        "n",
			Type:        "integer",
			Description: "Limit the number of entries in each response. It not present, all entries will be returned.",
			Format:      "<integer>",
			Required:    false,
		***REMOVED***,
		***REMOVED***
			Name:        "last",
			Type:        "string",
			Description: "Result set will include values lexically after last.",
			Format:      "<integer>",
			Required:    false,
		***REMOVED***,
	***REMOVED***

	unauthorizedResponseDescriptor = ResponseDescriptor***REMOVED***
		Name:        "Authentication Required",
		StatusCode:  http.StatusUnauthorized,
		Description: "The client is not authenticated.",
		Headers: []ParameterDescriptor***REMOVED***
			authChallengeHeader,
			***REMOVED***
				Name:        "Content-Length",
				Type:        "integer",
				Description: "Length of the JSON response body.",
				Format:      "<length>",
			***REMOVED***,
		***REMOVED***,
		Body: BodyDescriptor***REMOVED***
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		***REMOVED***,
		ErrorCodes: []errcode.ErrorCode***REMOVED***
			errcode.ErrorCodeUnauthorized,
		***REMOVED***,
	***REMOVED***

	repositoryNotFoundResponseDescriptor = ResponseDescriptor***REMOVED***
		Name:        "No Such Repository Error",
		StatusCode:  http.StatusNotFound,
		Description: "The repository is not known to the registry.",
		Headers: []ParameterDescriptor***REMOVED***
			***REMOVED***
				Name:        "Content-Length",
				Type:        "integer",
				Description: "Length of the JSON response body.",
				Format:      "<length>",
			***REMOVED***,
		***REMOVED***,
		Body: BodyDescriptor***REMOVED***
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		***REMOVED***,
		ErrorCodes: []errcode.ErrorCode***REMOVED***
			ErrorCodeNameUnknown,
		***REMOVED***,
	***REMOVED***

	deniedResponseDescriptor = ResponseDescriptor***REMOVED***
		Name:        "Access Denied",
		StatusCode:  http.StatusForbidden,
		Description: "The client does not have required access to the repository.",
		Headers: []ParameterDescriptor***REMOVED***
			***REMOVED***
				Name:        "Content-Length",
				Type:        "integer",
				Description: "Length of the JSON response body.",
				Format:      "<length>",
			***REMOVED***,
		***REMOVED***,
		Body: BodyDescriptor***REMOVED***
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		***REMOVED***,
		ErrorCodes: []errcode.ErrorCode***REMOVED***
			errcode.ErrorCodeDenied,
		***REMOVED***,
	***REMOVED***

	tooManyRequestsDescriptor = ResponseDescriptor***REMOVED***
		Name:        "Too Many Requests",
		StatusCode:  http.StatusTooManyRequests,
		Description: "The client made too many requests within a time interval.",
		Headers: []ParameterDescriptor***REMOVED***
			***REMOVED***
				Name:        "Content-Length",
				Type:        "integer",
				Description: "Length of the JSON response body.",
				Format:      "<length>",
			***REMOVED***,
		***REMOVED***,
		Body: BodyDescriptor***REMOVED***
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		***REMOVED***,
		ErrorCodes: []errcode.ErrorCode***REMOVED***
			errcode.ErrorCodeTooManyRequests,
		***REMOVED***,
	***REMOVED***
)

const (
	manifestBody = `***REMOVED***
   "name": <name>,
   "tag": <tag>,
   "fsLayers": [
      ***REMOVED***
         "blobSum": "<digest>"
  ***REMOVED***,
      ...
    ]
   ],
   "history": <v1 images>,
   "signature": <JWS>
***REMOVED***`

	errorsBody = `***REMOVED***
	"errors:" [
	    ***REMOVED***
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
    ***REMOVED***,
        ...
    ]
***REMOVED***`
)

// APIDescriptor exports descriptions of the layout of the v2 registry API.
var APIDescriptor = struct ***REMOVED***
	// RouteDescriptors provides a list of the routes available in the API.
	RouteDescriptors []RouteDescriptor
***REMOVED******REMOVED***
	RouteDescriptors: routeDescriptors,
***REMOVED***

// RouteDescriptor describes a route specified by name.
type RouteDescriptor struct ***REMOVED***
	// Name is the name of the route, as specified in RouteNameXXX exports.
	// These names a should be considered a unique reference for a route. If
	// the route is registered with gorilla, this is the name that will be
	// used.
	Name string

	// Path is a gorilla/mux-compatible regexp that can be used to match the
	// route. For any incoming method and path, only one route descriptor
	// should match.
	Path string

	// Entity should be a short, human-readalbe description of the object
	// targeted by the endpoint.
	Entity string

	// Description should provide an accurate overview of the functionality
	// provided by the route.
	Description string

	// Methods should describe the various HTTP methods that may be used on
	// this route, including request and response formats.
	Methods []MethodDescriptor
***REMOVED***

// MethodDescriptor provides a description of the requests that may be
// conducted with the target method.
type MethodDescriptor struct ***REMOVED***

	// Method is an HTTP method, such as GET, PUT or POST.
	Method string

	// Description should provide an overview of the functionality provided by
	// the covered method, suitable for use in documentation. Use of markdown
	// here is encouraged.
	Description string

	// Requests is a slice of request descriptors enumerating how this
	// endpoint may be used.
	Requests []RequestDescriptor
***REMOVED***

// RequestDescriptor covers a particular set of headers and parameters that
// can be carried out with the parent method. Its most helpful to have one
// RequestDescriptor per API use case.
type RequestDescriptor struct ***REMOVED***
	// Name provides a short identifier for the request, usable as a title or
	// to provide quick context for the particular request.
	Name string

	// Description should cover the requests purpose, covering any details for
	// this particular use case.
	Description string

	// Headers describes headers that must be used with the HTTP request.
	Headers []ParameterDescriptor

	// PathParameters enumerate the parameterized path components for the
	// given request, as defined in the route's regular expression.
	PathParameters []ParameterDescriptor

	// QueryParameters provides a list of query parameters for the given
	// request.
	QueryParameters []ParameterDescriptor

	// Body describes the format of the request body.
	Body BodyDescriptor

	// Successes enumerates the possible responses that are considered to be
	// the result of a successful request.
	Successes []ResponseDescriptor

	// Failures covers the possible failures from this particular request.
	Failures []ResponseDescriptor
***REMOVED***

// ResponseDescriptor describes the components of an API response.
type ResponseDescriptor struct ***REMOVED***
	// Name provides a short identifier for the response, usable as a title or
	// to provide quick context for the particular response.
	Name string

	// Description should provide a brief overview of the role of the
	// response.
	Description string

	// StatusCode specifies the status received by this particular response.
	StatusCode int

	// Headers covers any headers that may be returned from the response.
	Headers []ParameterDescriptor

	// Fields describes any fields that may be present in the response.
	Fields []ParameterDescriptor

	// ErrorCodes enumerates the error codes that may be returned along with
	// the response.
	ErrorCodes []errcode.ErrorCode

	// Body describes the body of the response, if any.
	Body BodyDescriptor
***REMOVED***

// BodyDescriptor describes a request body and its expected content type. For
// the most  part, it should be example json or some placeholder for body
// data in documentation.
type BodyDescriptor struct ***REMOVED***
	ContentType string
	Format      string
***REMOVED***

// ParameterDescriptor describes the format of a request parameter, which may
// be a header, path parameter or query parameter.
type ParameterDescriptor struct ***REMOVED***
	// Name is the name of the parameter, either of the path component or
	// query parameter.
	Name string

	// Type specifies the type of the parameter, such as string, integer, etc.
	Type string

	// Description provides a human-readable description of the parameter.
	Description string

	// Required means the field is required when set.
	Required bool

	// Format is a specifying the string format accepted by this parameter.
	Format string

	// Regexp is a compiled regular expression that can be used to validate
	// the contents of the parameter.
	Regexp *regexp.Regexp

	// Examples provides multiple examples for the values that might be valid
	// for this parameter.
	Examples []string
***REMOVED***

var routeDescriptors = []RouteDescriptor***REMOVED***
	***REMOVED***
		Name:        RouteNameBase,
		Path:        "/v2/",
		Entity:      "Base",
		Description: `Base V2 API route. Typically, this can be used for lightweight version checks and to validate registry authentication.`,
		Methods: []MethodDescriptor***REMOVED***
			***REMOVED***
				Method:      "GET",
				Description: "Check that the endpoint implements Docker Registry API V2.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The API implements V2 protocol and is accessible.",
								StatusCode:  http.StatusOK,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The registry does not implement the V2 API.",
								StatusCode:  http.StatusNotFound,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Name:        RouteNameTags,
		Path:        "/v2/***REMOVED***name:" + reference.NameRegexp.String() + "***REMOVED***/tags/list",
		Entity:      "Tags",
		Description: "Retrieve information about tags.",
		Methods: []MethodDescriptor***REMOVED***
			***REMOVED***
				Method:      "GET",
				Description: "Fetch the tags under the repository identified by `name`.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Name:        "Tags",
						Description: "Return all tags for the repository",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								StatusCode:  http.StatusOK,
								Description: "A list of tags for the named repository.",
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Content-Length",
										Type:        "integer",
										Description: "Length of the JSON response body.",
										Format:      "<length>",
									***REMOVED***,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format: `***REMOVED***
    "name": <name>,
    "tags": [
        <tag>,
        ...
    ]
***REMOVED***`,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
					***REMOVED***
						Name:            "Tags Paginated",
						Description:     "Return a portion of the tags for the specified repository.",
						PathParameters:  []ParameterDescriptor***REMOVED***nameParameterDescriptor***REMOVED***,
						QueryParameters: paginationParameters,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								StatusCode:  http.StatusOK,
								Description: "A list of tags for the named repository.",
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Content-Length",
										Type:        "integer",
										Description: "Length of the JSON response body.",
										Format:      "<length>",
									***REMOVED***,
									linkHeader,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format: `***REMOVED***
    "name": <name>,
    "tags": [
        <tag>,
        ...
    ],
***REMOVED***`,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Name:        RouteNameManifest,
		Path:        "/v2/***REMOVED***name:" + reference.NameRegexp.String() + "***REMOVED***/manifests/***REMOVED***reference:" + reference.TagRegexp.String() + "|" + digest.DigestRegexp.String() + "***REMOVED***",
		Entity:      "Manifest",
		Description: "Create, update, delete and retrieve manifests.",
		Methods: []MethodDescriptor***REMOVED***
			***REMOVED***
				Method:      "GET",
				Description: "Fetch the manifest identified by `name` and `reference` where `reference` can be a tag or digest. A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							referenceParameterDescriptor,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The manifest identified by `name` and `reference`. The contents can be used to identify and resolve resources required to run the specified image.",
								StatusCode:  http.StatusOK,
								Headers: []ParameterDescriptor***REMOVED***
									digestHeader,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "<media type of manifest>",
									Format:      manifestBody,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The name or reference was invalid.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameInvalid,
									ErrorCodeTagInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Method:      "PUT",
				Description: "Put the manifest identified by `name` and `reference` where `reference` can be a tag or digest.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							referenceParameterDescriptor,
						***REMOVED***,
						Body: BodyDescriptor***REMOVED***
							ContentType: "<media type of manifest>",
							Format:      manifestBody,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The manifest has been accepted by the registry and is stored under the specified `name` and `tag`.",
								StatusCode:  http.StatusCreated,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Location",
										Type:        "url",
										Description: "The canonical location url of the uploaded manifest.",
										Format:      "<url>",
									***REMOVED***,
									contentLengthZeroHeader,
									digestHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:        "Invalid Manifest",
								Description: "The received manifest was invalid in some way, as described by the error codes. The client should resolve the issue and retry the request.",
								StatusCode:  http.StatusBadRequest,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameInvalid,
									ErrorCodeTagInvalid,
									ErrorCodeManifestInvalid,
									ErrorCodeManifestUnverified,
									ErrorCodeBlobUnknown,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
							***REMOVED***
								Name:        "Missing Layer(s)",
								Description: "One or more layers may be missing during a manifest upload. If so, the missing layers will be enumerated in the error response.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeBlobUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format: `***REMOVED***
    "errors:" [***REMOVED***
            "code": "BLOB_UNKNOWN",
            "message": "blob unknown to registry",
            "detail": ***REMOVED***
                "digest": "<digest>"
        ***REMOVED***
    ***REMOVED***,
        ...
    ]
***REMOVED***`,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Name:        "Not allowed",
								Description: "Manifest put is not allowed because the registry is configured as a pull-through cache or for some other reason",
								StatusCode:  http.StatusMethodNotAllowed,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									errcode.ErrorCodeUnsupported,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Method:      "DELETE",
				Description: "Delete the manifest identified by `name` and `reference`. Note that a manifest can _only_ be deleted by `digest`.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							referenceParameterDescriptor,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								StatusCode: http.StatusAccepted,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:        "Invalid Name or Reference",
								Description: "The specified `name` or `reference` were invalid and the delete was unable to proceed.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameInvalid,
									ErrorCodeTagInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
							***REMOVED***
								Name:        "Unknown Manifest",
								Description: "The specified `name` or `reference` are unknown to the registry and the delete was unable to proceed. Clients can assume the manifest was already deleted if this response is returned.",
								StatusCode:  http.StatusNotFound,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameUnknown,
									ErrorCodeManifestUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Name:        "Not allowed",
								Description: "Manifest delete is not allowed because the registry is configured as a pull-through cache or `delete` has been disabled.",
								StatusCode:  http.StatusMethodNotAllowed,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									errcode.ErrorCodeUnsupported,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,

	***REMOVED***
		Name:        RouteNameBlob,
		Path:        "/v2/***REMOVED***name:" + reference.NameRegexp.String() + "***REMOVED***/blobs/***REMOVED***digest:" + digest.DigestRegexp.String() + "***REMOVED***",
		Entity:      "Blob",
		Description: "Operations on blobs identified by `name` and `digest`. Used to fetch or delete layers by digest.",
		Methods: []MethodDescriptor***REMOVED***
			***REMOVED***
				Method:      "GET",
				Description: "Retrieve the blob from the registry identified by `digest`. A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Name: "Fetch Blob",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							digestPathParameter,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The blob identified by `digest` is available. The blob content will be present in the body of the request.",
								StatusCode:  http.StatusOK,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Content-Length",
										Type:        "integer",
										Description: "The length of the requested blob content.",
										Format:      "<length>",
									***REMOVED***,
									digestHeader,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/octet-stream",
									Format:      "<blob binary data>",
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The blob identified by `digest` is available at the provided location.",
								StatusCode:  http.StatusTemporaryRedirect,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Location",
										Type:        "url",
										Description: "The location where the layer should be accessible.",
										Format:      "<blob location>",
									***REMOVED***,
									digestHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "There was a problem with the request that needs to be addressed by the client, such as an invalid `name` or `tag`.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameInvalid,
									ErrorCodeDigestInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The blob, identified by `name` and `digest`, is unknown to the registry.",
								StatusCode:  http.StatusNotFound,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameUnknown,
									ErrorCodeBlobUnknown,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
					***REMOVED***
						Name:        "Fetch Blob Part",
						Description: "This endpoint may also support RFC7233 compliant range requests. Support can be detected by issuing a HEAD request. If the header `Accept-Range: bytes` is returned, range requests can be used to fetch partial content.",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
							***REMOVED***
								Name:        "Range",
								Type:        "string",
								Description: "HTTP Range header specifying blob chunk.",
								Format:      "bytes=<start>-<end>",
							***REMOVED***,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							digestPathParameter,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The blob identified by `digest` is available. The specified chunk of blob content will be present in the body of the request.",
								StatusCode:  http.StatusPartialContent,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Content-Length",
										Type:        "integer",
										Description: "The length of the requested blob chunk.",
										Format:      "<length>",
									***REMOVED***,
									***REMOVED***
										Name:        "Content-Range",
										Type:        "byte range",
										Description: "Content range of blob chunk.",
										Format:      "bytes <start>-<end>/<size>",
									***REMOVED***,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/octet-stream",
									Format:      "<blob binary data>",
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "There was a problem with the request that needs to be addressed by the client, such as an invalid `name` or `tag`.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameInvalid,
									ErrorCodeDigestInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								StatusCode: http.StatusNotFound,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameUnknown,
									ErrorCodeBlobUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The range specification cannot be satisfied for the requested content. This can happen when the range is not formatted correctly or if the range is outside of the valid size of the content.",
								StatusCode:  http.StatusRequestedRangeNotSatisfiable,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Method:      "DELETE",
				Description: "Delete the blob identified by `name` and `digest`",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							digestPathParameter,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								StatusCode: http.StatusAccepted,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Content-Length",
										Type:        "integer",
										Description: "0",
										Format:      "0",
									***REMOVED***,
									digestHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:       "Invalid Name or Digest",
								StatusCode: http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The blob, identified by `name` and `digest`, is unknown to the registry.",
								StatusCode:  http.StatusNotFound,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameUnknown,
									ErrorCodeBlobUnknown,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "Blob delete is not allowed because the registry is configured as a pull-through cache or `delete` has been disabled",
								StatusCode:  http.StatusMethodNotAllowed,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									errcode.ErrorCodeUnsupported,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,

			// TODO(stevvooe): We may want to add a PUT request here to
			// kickoff an upload of a blob, integrated with the blob upload
			// API.
		***REMOVED***,
	***REMOVED***,

	***REMOVED***
		Name:        RouteNameBlobUpload,
		Path:        "/v2/***REMOVED***name:" + reference.NameRegexp.String() + "***REMOVED***/blobs/uploads/",
		Entity:      "Initiate Blob Upload",
		Description: "Initiate a blob upload. This endpoint can be used to create resumable uploads or monolithic uploads.",
		Methods: []MethodDescriptor***REMOVED***
			***REMOVED***
				Method:      "POST",
				Description: "Initiate a resumable blob upload. If successful, an upload location will be provided to complete the upload. Optionally, if the `digest` parameter is present, the request body will be used to complete the upload in a single request.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Name:        "Initiate Monolithic Blob Upload",
						Description: "Upload a blob identified by the `digest` parameter in single request. This upload will not be resumable unless a recoverable error is returned.",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
							***REMOVED***
								Name:   "Content-Length",
								Type:   "integer",
								Format: "<length of blob>",
							***REMOVED***,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
						***REMOVED***,
						QueryParameters: []ParameterDescriptor***REMOVED***
							***REMOVED***
								Name:        "digest",
								Type:        "query",
								Format:      "<digest>",
								Regexp:      digest.DigestRegexp,
								Description: `Digest of uploaded blob. If present, the upload will be completed, in a single request, with contents of the request body as the resulting blob.`,
							***REMOVED***,
						***REMOVED***,
						Body: BodyDescriptor***REMOVED***
							ContentType: "application/octect-stream",
							Format:      "<binary data>",
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The blob has been created in the registry and is available at the provided location.",
								StatusCode:  http.StatusCreated,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:   "Location",
										Type:   "url",
										Format: "<blob location>",
									***REMOVED***,
									contentLengthZeroHeader,
									dockerUploadUUIDHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:       "Invalid Name or Digest",
								StatusCode: http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Name:        "Not allowed",
								Description: "Blob upload is not allowed because the registry is configured as a pull-through cache or for some other reason",
								StatusCode:  http.StatusMethodNotAllowed,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									errcode.ErrorCodeUnsupported,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
					***REMOVED***
						Name:        "Initiate Resumable Blob Upload",
						Description: "Initiate a resumable blob upload with an empty request body.",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
							contentLengthZeroHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The upload has been created. The `Location` header must be used to complete the upload. The response should be identical to a `GET` request on the contents of the returned `Location` header.",
								StatusCode:  http.StatusAccepted,
								Headers: []ParameterDescriptor***REMOVED***
									contentLengthZeroHeader,
									***REMOVED***
										Name:        "Location",
										Type:        "url",
										Format:      "/v2/<name>/blobs/uploads/<uuid>",
										Description: "The location of the created upload. Clients should use the contents verbatim to complete the upload, adding parameters where required.",
									***REMOVED***,
									***REMOVED***
										Name:        "Range",
										Format:      "0-0",
										Description: "Range header indicating the progress of the upload. When starting an upload, it will return an empty range, since no content has been received.",
									***REMOVED***,
									dockerUploadUUIDHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:       "Invalid Name or Digest",
								StatusCode: http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
					***REMOVED***
						Name:        "Mount Blob",
						Description: "Mount a blob identified by the `mount` parameter from another repository.",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
							contentLengthZeroHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
						***REMOVED***,
						QueryParameters: []ParameterDescriptor***REMOVED***
							***REMOVED***
								Name:        "mount",
								Type:        "query",
								Format:      "<digest>",
								Regexp:      digest.DigestRegexp,
								Description: `Digest of blob to mount from the source repository.`,
							***REMOVED***,
							***REMOVED***
								Name:        "from",
								Type:        "query",
								Format:      "<repository name>",
								Regexp:      reference.NameRegexp,
								Description: `Name of the source repository.`,
							***REMOVED***,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "The blob has been mounted in the repository and is available at the provided location.",
								StatusCode:  http.StatusCreated,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:   "Location",
										Type:   "url",
										Format: "<blob location>",
									***REMOVED***,
									contentLengthZeroHeader,
									dockerUploadUUIDHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:       "Invalid Name or Digest",
								StatusCode: http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Name:        "Not allowed",
								Description: "Blob mount is not allowed because the registry is configured as a pull-through cache or for some other reason",
								StatusCode:  http.StatusMethodNotAllowed,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									errcode.ErrorCodeUnsupported,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,

	***REMOVED***
		Name:        RouteNameBlobUploadChunk,
		Path:        "/v2/***REMOVED***name:" + reference.NameRegexp.String() + "***REMOVED***/blobs/uploads/***REMOVED***uuid:[a-zA-Z0-9-_.=]+***REMOVED***",
		Entity:      "Blob Upload",
		Description: "Interact with blob uploads. Clients should never assemble URLs for this endpoint and should only take it through the `Location` header on related API requests. The `Location` header and its parameters should be preserved by clients, using the latest value returned via upload related API calls.",
		Methods: []MethodDescriptor***REMOVED***
			***REMOVED***
				Method:      "GET",
				Description: "Retrieve status of upload identified by `uuid`. The primary purpose of this endpoint is to resolve the current status of a resumable upload.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Description: "Retrieve the progress of the current upload, as reported by the `Range` header.",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							uuidParameterDescriptor,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:        "Upload Progress",
								Description: "The upload is known and in progress. The last received offset is available in the `Range` header.",
								StatusCode:  http.StatusNoContent,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Range",
										Type:        "header",
										Format:      "0-<offset>",
										Description: "Range indicating the current progress of the upload.",
									***REMOVED***,
									contentLengthZeroHeader,
									dockerUploadUUIDHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "There was an error processing the upload and it must be restarted.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
									ErrorCodeBlobUploadInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The upload is unknown to the registry. The upload must be restarted.",
								StatusCode:  http.StatusNotFound,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeBlobUploadUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Method:      "PATCH",
				Description: "Upload a chunk of data for the specified upload.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Name:        "Stream upload",
						Description: "Upload a stream of data to upload without completing the upload.",
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							uuidParameterDescriptor,
						***REMOVED***,
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
						***REMOVED***,
						Body: BodyDescriptor***REMOVED***
							ContentType: "application/octet-stream",
							Format:      "<binary data>",
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:        "Data Accepted",
								Description: "The stream of data has been accepted and the current progress is available in the range header. The updated upload location is available in the `Location` header.",
								StatusCode:  http.StatusNoContent,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Location",
										Type:        "url",
										Format:      "/v2/<name>/blobs/uploads/<uuid>",
										Description: "The location of the upload. Clients should assume this changes after each request. Clients should use the contents verbatim to complete the upload, adding parameters where required.",
									***REMOVED***,
									***REMOVED***
										Name:        "Range",
										Type:        "header",
										Format:      "0-<offset>",
										Description: "Range indicating the current progress of the upload.",
									***REMOVED***,
									contentLengthZeroHeader,
									dockerUploadUUIDHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "There was an error processing the upload and it must be restarted.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
									ErrorCodeBlobUploadInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The upload is unknown to the registry. The upload must be restarted.",
								StatusCode:  http.StatusNotFound,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeBlobUploadUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
					***REMOVED***
						Name:        "Chunked upload",
						Description: "Upload a chunk of data to specified upload without completing the upload. The data will be uploaded to the specified Content Range.",
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							uuidParameterDescriptor,
						***REMOVED***,
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
							***REMOVED***
								Name:        "Content-Range",
								Type:        "header",
								Format:      "<start of range>-<end of range, inclusive>",
								Required:    true,
								Description: "Range of bytes identifying the desired block of content represented by the body. Start must the end offset retrieved via status check plus one. Note that this is a non-standard use of the `Content-Range` header.",
							***REMOVED***,
							***REMOVED***
								Name:        "Content-Length",
								Type:        "integer",
								Format:      "<length of chunk>",
								Description: "Length of the chunk being uploaded, corresponding the length of the request body.",
							***REMOVED***,
						***REMOVED***,
						Body: BodyDescriptor***REMOVED***
							ContentType: "application/octet-stream",
							Format:      "<binary chunk>",
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:        "Chunk Accepted",
								Description: "The chunk of data has been accepted and the current progress is available in the range header. The updated upload location is available in the `Location` header.",
								StatusCode:  http.StatusNoContent,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Location",
										Type:        "url",
										Format:      "/v2/<name>/blobs/uploads/<uuid>",
										Description: "The location of the upload. Clients should assume this changes after each request. Clients should use the contents verbatim to complete the upload, adding parameters where required.",
									***REMOVED***,
									***REMOVED***
										Name:        "Range",
										Type:        "header",
										Format:      "0-<offset>",
										Description: "Range indicating the current progress of the upload.",
									***REMOVED***,
									contentLengthZeroHeader,
									dockerUploadUUIDHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "There was an error processing the upload and it must be restarted.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
									ErrorCodeBlobUploadInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The upload is unknown to the registry. The upload must be restarted.",
								StatusCode:  http.StatusNotFound,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeBlobUploadUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The `Content-Range` specification cannot be accepted, either because it does not overlap with the current progress or it is invalid.",
								StatusCode:  http.StatusRequestedRangeNotSatisfiable,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Method:      "PUT",
				Description: "Complete the upload specified by `uuid`, optionally appending the body as the final chunk.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Description: "Complete the upload, providing all the data in the body, if necessary. A request without a body will just complete the upload with previously uploaded content.",
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
							***REMOVED***
								Name:        "Content-Length",
								Type:        "integer",
								Format:      "<length of data>",
								Description: "Length of the data being uploaded, corresponding to the length of the request body. May be zero if no data is provided.",
							***REMOVED***,
						***REMOVED***,
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							uuidParameterDescriptor,
						***REMOVED***,
						QueryParameters: []ParameterDescriptor***REMOVED***
							***REMOVED***
								Name:        "digest",
								Type:        "string",
								Format:      "<digest>",
								Regexp:      digest.DigestRegexp,
								Required:    true,
								Description: `Digest of uploaded blob.`,
							***REMOVED***,
						***REMOVED***,
						Body: BodyDescriptor***REMOVED***
							ContentType: "application/octet-stream",
							Format:      "<binary data>",
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:        "Upload Complete",
								Description: "The upload has been completed and accepted by the registry. The canonical location will be available in the `Location` header.",
								StatusCode:  http.StatusNoContent,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Location",
										Type:        "url",
										Format:      "<blob location>",
										Description: "The canonical location of the blob for retrieval",
									***REMOVED***,
									***REMOVED***
										Name:        "Content-Range",
										Type:        "header",
										Format:      "<start of range>-<end of range, inclusive>",
										Description: "Range of bytes identifying the desired block of content represented by the body. Start must match the end of offset retrieved via status check. Note that this is a non-standard use of the `Content-Range` header.",
									***REMOVED***,
									contentLengthZeroHeader,
									digestHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "There was an error processing the upload and it must be restarted.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeDigestInvalid,
									ErrorCodeNameInvalid,
									ErrorCodeBlobUploadInvalid,
									errcode.ErrorCodeUnsupported,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The upload is unknown to the registry. The upload must be restarted.",
								StatusCode:  http.StatusNotFound,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeBlobUploadUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Method:      "DELETE",
				Description: "Cancel outstanding upload processes, releasing associated resources. If this is not called, the unfinished uploads will eventually timeout.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Description: "Cancel the upload specified by `uuid`.",
						PathParameters: []ParameterDescriptor***REMOVED***
							nameParameterDescriptor,
							uuidParameterDescriptor,
						***REMOVED***,
						Headers: []ParameterDescriptor***REMOVED***
							hostHeader,
							authHeader,
							contentLengthZeroHeader,
						***REMOVED***,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Name:        "Upload Deleted",
								Description: "The upload has been successfully deleted.",
								StatusCode:  http.StatusNoContent,
								Headers: []ParameterDescriptor***REMOVED***
									contentLengthZeroHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
						Failures: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "An error was encountered processing the delete. The client may ignore this error.",
								StatusCode:  http.StatusBadRequest,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeNameInvalid,
									ErrorCodeBlobUploadInvalid,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							***REMOVED***
								Description: "The upload is unknown to the registry. The client may ignore this error and assume the upload has been deleted.",
								StatusCode:  http.StatusNotFound,
								ErrorCodes: []errcode.ErrorCode***REMOVED***
									ErrorCodeBlobUploadUnknown,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								***REMOVED***,
							***REMOVED***,
							unauthorizedResponseDescriptor,
							repositoryNotFoundResponseDescriptor,
							deniedResponseDescriptor,
							tooManyRequestsDescriptor,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Name:        RouteNameCatalog,
		Path:        "/v2/_catalog",
		Entity:      "Catalog",
		Description: "List a set of available repositories in the local registry cluster. Does not provide any indication of what may be available upstream. Applications can only determine if a repository is available but not if it is not available.",
		Methods: []MethodDescriptor***REMOVED***
			***REMOVED***
				Method:      "GET",
				Description: "Retrieve a sorted, json list of repositories available in the registry.",
				Requests: []RequestDescriptor***REMOVED***
					***REMOVED***
						Name:        "Catalog Fetch",
						Description: "Request an unabridged list of repositories available.  The implementation may impose a maximum limit and return a partial set with pagination links.",
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								Description: "Returns the unabridged list of repositories as a json response.",
								StatusCode:  http.StatusOK,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Content-Length",
										Type:        "integer",
										Description: "Length of the JSON response body.",
										Format:      "<length>",
									***REMOVED***,
								***REMOVED***,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format: `***REMOVED***
	"repositories": [
		<name>,
		...
	]
***REMOVED***`,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
					***REMOVED***
						Name:            "Catalog Fetch Paginated",
						Description:     "Return the specified portion of repositories.",
						QueryParameters: paginationParameters,
						Successes: []ResponseDescriptor***REMOVED***
							***REMOVED***
								StatusCode: http.StatusOK,
								Body: BodyDescriptor***REMOVED***
									ContentType: "application/json; charset=utf-8",
									Format: `***REMOVED***
	"repositories": [
		<name>,
		...
	]
	"next": "<url>?last=<name>&n=<last value of n>"
***REMOVED***`,
								***REMOVED***,
								Headers: []ParameterDescriptor***REMOVED***
									***REMOVED***
										Name:        "Content-Length",
										Type:        "integer",
										Description: "Length of the JSON response body.",
										Format:      "<length>",
									***REMOVED***,
									linkHeader,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

var routeDescriptorsMap map[string]RouteDescriptor

func init() ***REMOVED***
	routeDescriptorsMap = make(map[string]RouteDescriptor, len(routeDescriptors))

	for _, descriptor := range routeDescriptors ***REMOVED***
		routeDescriptorsMap[descriptor.Name] = descriptor
	***REMOVED***
***REMOVED***
