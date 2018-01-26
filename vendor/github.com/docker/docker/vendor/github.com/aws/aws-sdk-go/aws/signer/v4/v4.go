// Package v4 implements signing for AWS V4 signer
//
// Provides request signing for request that need to be signed with
// AWS V4 Signatures.
//
// Standalone Signer
//
// Generally using the signer outside of the SDK should not require any additional
// logic when using Go v1.5 or higher. The signer does this by taking advantage
// of the URL.EscapedPath method. If your request URI requires additional escaping
// you many need to use the URL.Opaque to define what the raw URI should be sent
// to the service as.
//
// The signer will first check the URL.Opaque field, and use its value if set.
// The signer does require the URL.Opaque field to be set in the form of:
//
//     "//<hostname>/<path>"
//
//     // e.g.
//     "//example.com/some/path"
//
// The leading "//" and hostname are required or the URL.Opaque escaping will
// not work correctly.
//
// If URL.Opaque is not set the signer will fallback to the URL.EscapedPath()
// method and using the returned value. If you're using Go v1.4 you must set
// URL.Opaque if the URI path needs escaping. If URL.Opaque is not set with
// Go v1.5 the signer will fallback to URL.Path.
//
// AWS v4 signature validation requires that the canonical string's URI path
// element must be the URI escaped form of the HTTP request's path.
// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
//
// The Go HTTP client will perform escaping automatically on the request. Some
// of these escaping may cause signature validation errors because the HTTP
// request differs from the URI path or query that the signature was generated.
// https://golang.org/pkg/net/url/#URL.EscapedPath
//
// Because of this, it is recommended that when using the signer outside of the
// SDK that explicitly escaping the request prior to being signed is preferable,
// and will help prevent signature validation errors. This can be done by setting
// the URL.Opaque or URL.RawPath. The SDK will use URL.Opaque first and then
// call URL.EscapedPath() if Opaque is not set.
//
// If signing a request intended for HTTP2 server, and you're using Go 1.6.2
// through 1.7.4 you should use the URL.RawPath as the pre-escaped form of the
// request URL. https://github.com/golang/go/issues/16847 points to a bug in
// Go pre 1.8 that fails to make HTTP2 requests using absolute URL in the HTTP
// message. URL.Opaque generally will force Go to make requests with absolute URL.
// URL.RawPath does not do this, but RawPath must be a valid escaping of Path
// or url.EscapedPath will ignore the RawPath escaping.
//
// Test `TestStandaloneSign` provides a complete example of using the signer
// outside of the SDK and pre-escaping the URI path.
package v4

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
)

const (
	authHeaderPrefix = "AWS4-HMAC-SHA256"
	timeFormat       = "20060102T150405Z"
	shortTimeFormat  = "20060102"

	// emptyStringSHA256 is a SHA256 of an empty string
	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
)

var ignoredHeaders = rules***REMOVED***
	blacklist***REMOVED***
		mapRule***REMOVED***
			"Authorization":   struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"User-Agent":      struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amzn-Trace-Id": struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

// requiredSignedHeaders is a whitelist for build canonical headers.
var requiredSignedHeaders = rules***REMOVED***
	whitelist***REMOVED***
		mapRule***REMOVED***
			"Cache-Control":                                               struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"Content-Disposition":                                         struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"Content-Encoding":                                            struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"Content-Language":                                            struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"Content-Md5":                                                 struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"Content-Type":                                                struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"Expires":                                                     struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"If-Match":                                                    struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"If-Modified-Since":                                           struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"If-None-Match":                                               struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"If-Unmodified-Since":                                         struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"Range":                                                       struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Acl":                                                   struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source":                                           struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-If-Match":                                  struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-If-Modified-Since":                         struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-If-None-Match":                             struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-If-Unmodified-Since":                       struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-Range":                                     struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-Server-Side-Encryption-Customer-Algorithm": struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-Server-Side-Encryption-Customer-Key":       struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Copy-Source-Server-Side-Encryption-Customer-Key-Md5":   struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Grant-Full-control":                                    struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Grant-Read":                                            struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Grant-Read-Acp":                                        struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Grant-Write":                                           struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Grant-Write-Acp":                                       struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Metadata-Directive":                                    struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Mfa":                                                   struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Request-Payer":                                         struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Server-Side-Encryption":                                struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Server-Side-Encryption-Aws-Kms-Key-Id":                 struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Server-Side-Encryption-Customer-Algorithm":             struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Server-Side-Encryption-Customer-Key":                   struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Server-Side-Encryption-Customer-Key-Md5":               struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Storage-Class":                                         struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amz-Website-Redirect-Location":                             struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	patterns***REMOVED***"X-Amz-Meta-"***REMOVED***,
***REMOVED***

// allowedHoisting is a whitelist for build query headers. The boolean value
// represents whether or not it is a pattern.
var allowedQueryHoisting = inclusiveRules***REMOVED***
	blacklist***REMOVED***requiredSignedHeaders***REMOVED***,
	patterns***REMOVED***"X-Amz-"***REMOVED***,
***REMOVED***

// Signer applies AWS v4 signing to given request. Use this to sign requests
// that need to be signed with AWS V4 Signatures.
type Signer struct ***REMOVED***
	// The authentication credentials the request will be signed against.
	// This value must be set to sign requests.
	Credentials *credentials.Credentials

	// Sets the log level the signer should use when reporting information to
	// the logger. If the logger is nil nothing will be logged. See
	// aws.LogLevelType for more information on available logging levels
	//
	// By default nothing will be logged.
	Debug aws.LogLevelType

	// The logger loging information will be written to. If there the logger
	// is nil, nothing will be logged.
	Logger aws.Logger

	// Disables the Signer's moving HTTP header key/value pairs from the HTTP
	// request header to the request's query string. This is most commonly used
	// with pre-signed requests preventing headers from being added to the
	// request's query string.
	DisableHeaderHoisting bool

	// Disables the automatic escaping of the URI path of the request for the
	// siganture's canonical string's path. For services that do not need additional
	// escaping then use this to disable the signer escaping the path.
	//
	// S3 is an example of a service that does not need additional escaping.
	//
	// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
	DisableURIPathEscaping bool

	// Disales the automatical setting of the HTTP request's Body field with the
	// io.ReadSeeker passed in to the signer. This is useful if you're using a
	// custom wrapper around the body for the io.ReadSeeker and want to preserve
	// the Body value on the Request.Body.
	//
	// This does run the risk of signing a request with a body that will not be
	// sent in the request. Need to ensure that the underlying data of the Body
	// values are the same.
	DisableRequestBodyOverwrite bool

	// currentTimeFn returns the time value which represents the current time.
	// This value should only be used for testing. If it is nil the default
	// time.Now will be used.
	currentTimeFn func() time.Time

	// UnsignedPayload will prevent signing of the payload. This will only
	// work for services that have support for this.
	UnsignedPayload bool
***REMOVED***

// NewSigner returns a Signer pointer configured with the credentials and optional
// option values provided. If not options are provided the Signer will use its
// default configuration.
func NewSigner(credentials *credentials.Credentials, options ...func(*Signer)) *Signer ***REMOVED***
	v4 := &Signer***REMOVED***
		Credentials: credentials,
	***REMOVED***

	for _, option := range options ***REMOVED***
		option(v4)
	***REMOVED***

	return v4
***REMOVED***

type signingCtx struct ***REMOVED***
	ServiceName      string
	Region           string
	Request          *http.Request
	Body             io.ReadSeeker
	Query            url.Values
	Time             time.Time
	ExpireTime       time.Duration
	SignedHeaderVals http.Header

	DisableURIPathEscaping bool

	credValues         credentials.Value
	isPresign          bool
	formattedTime      string
	formattedShortTime string
	unsignedPayload    bool

	bodyDigest       string
	signedHeaders    string
	canonicalHeaders string
	canonicalString  string
	credentialString string
	stringToSign     string
	signature        string
	authorization    string
***REMOVED***

// Sign signs AWS v4 requests with the provided body, service name, region the
// request is made to, and time the request is signed at. The signTime allows
// you to specify that a request is signed for the future, and cannot be
// used until then.
//
// Returns a list of HTTP headers that were included in the signature or an
// error if signing the request failed. Generally for signed requests this value
// is not needed as the full request context will be captured by the http.Request
// value. It is included for reference though.
//
// Sign will set the request's Body to be the `body` parameter passed in. If
// the body is not already an io.ReadCloser, it will be wrapped within one. If
// a `nil` body parameter passed to Sign, the request's Body field will be
// also set to nil. Its important to note that this functionality will not
// change the request's ContentLength of the request.
//
// Sign differs from Presign in that it will sign the request using HTTP
// header values. This type of signing is intended for http.Request values that
// will not be shared, or are shared in a way the header values on the request
// will not be lost.
//
// The requests body is an io.ReadSeeker so the SHA256 of the body can be
// generated. To bypass the signer computing the hash you can set the
// "X-Amz-Content-Sha256" header with a precomputed value. The signer will
// only compute the hash if the request header value is empty.
func (v4 Signer) Sign(r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error) ***REMOVED***
	return v4.signWithBody(r, body, service, region, 0, false, signTime)
***REMOVED***

// Presign signs AWS v4 requests with the provided body, service name, region
// the request is made to, and time the request is signed at. The signTime
// allows you to specify that a request is signed for the future, and cannot
// be used until then.
//
// Returns a list of HTTP headers that were included in the signature or an
// error if signing the request failed. For presigned requests these headers
// and their values must be included on the HTTP request when it is made. This
// is helpful to know what header values need to be shared with the party the
// presigned request will be distributed to.
//
// Presign differs from Sign in that it will sign the request using query string
// instead of header values. This allows you to share the Presigned Request's
// URL with third parties, or distribute it throughout your system with minimal
// dependencies.
//
// Presign also takes an exp value which is the duration the
// signed request will be valid after the signing time. This is allows you to
// set when the request will expire.
//
// The requests body is an io.ReadSeeker so the SHA256 of the body can be
// generated. To bypass the signer computing the hash you can set the
// "X-Amz-Content-Sha256" header with a precomputed value. The signer will
// only compute the hash if the request header value is empty.
//
// Presigning a S3 request will not compute the body's SHA256 hash by default.
// This is done due to the general use case for S3 presigned URLs is to share
// PUT/GET capabilities. If you would like to include the body's SHA256 in the
// presigned request's signature you can set the "X-Amz-Content-Sha256"
// HTTP header and that will be included in the request's signature.
func (v4 Signer) Presign(r *http.Request, body io.ReadSeeker, service, region string, exp time.Duration, signTime time.Time) (http.Header, error) ***REMOVED***
	return v4.signWithBody(r, body, service, region, exp, true, signTime)
***REMOVED***

func (v4 Signer) signWithBody(r *http.Request, body io.ReadSeeker, service, region string, exp time.Duration, isPresign bool, signTime time.Time) (http.Header, error) ***REMOVED***
	currentTimeFn := v4.currentTimeFn
	if currentTimeFn == nil ***REMOVED***
		currentTimeFn = time.Now
	***REMOVED***

	ctx := &signingCtx***REMOVED***
		Request:                r,
		Body:                   body,
		Query:                  r.URL.Query(),
		Time:                   signTime,
		ExpireTime:             exp,
		isPresign:              isPresign,
		ServiceName:            service,
		Region:                 region,
		DisableURIPathEscaping: v4.DisableURIPathEscaping,
		unsignedPayload:        v4.UnsignedPayload,
	***REMOVED***

	for key := range ctx.Query ***REMOVED***
		sort.Strings(ctx.Query[key])
	***REMOVED***

	if ctx.isRequestSigned() ***REMOVED***
		ctx.Time = currentTimeFn()
		ctx.handlePresignRemoval()
	***REMOVED***

	var err error
	ctx.credValues, err = v4.Credentials.Get()
	if err != nil ***REMOVED***
		return http.Header***REMOVED******REMOVED***, err
	***REMOVED***

	ctx.sanitizeHostForHeader()
	ctx.assignAmzQueryValues()
	ctx.build(v4.DisableHeaderHoisting)

	// If the request is not presigned the body should be attached to it. This
	// prevents the confusion of wanting to send a signed request without
	// the body the request was signed for attached.
	if !(v4.DisableRequestBodyOverwrite || ctx.isPresign) ***REMOVED***
		var reader io.ReadCloser
		if body != nil ***REMOVED***
			var ok bool
			if reader, ok = body.(io.ReadCloser); !ok ***REMOVED***
				reader = ioutil.NopCloser(body)
			***REMOVED***
		***REMOVED***
		r.Body = reader
	***REMOVED***

	if v4.Debug.Matches(aws.LogDebugWithSigning) ***REMOVED***
		v4.logSigningInfo(ctx)
	***REMOVED***

	return ctx.SignedHeaderVals, nil
***REMOVED***

func (ctx *signingCtx) sanitizeHostForHeader() ***REMOVED***
	request.SanitizeHostForHeader(ctx.Request)
***REMOVED***

func (ctx *signingCtx) handlePresignRemoval() ***REMOVED***
	if !ctx.isPresign ***REMOVED***
		return
	***REMOVED***

	// The credentials have expired for this request. The current signing
	// is invalid, and needs to be request because the request will fail.
	ctx.removePresign()

	// Update the request's query string to ensure the values stays in
	// sync in the case retrieving the new credentials fails.
	ctx.Request.URL.RawQuery = ctx.Query.Encode()
***REMOVED***

func (ctx *signingCtx) assignAmzQueryValues() ***REMOVED***
	if ctx.isPresign ***REMOVED***
		ctx.Query.Set("X-Amz-Algorithm", authHeaderPrefix)
		if ctx.credValues.SessionToken != "" ***REMOVED***
			ctx.Query.Set("X-Amz-Security-Token", ctx.credValues.SessionToken)
		***REMOVED*** else ***REMOVED***
			ctx.Query.Del("X-Amz-Security-Token")
		***REMOVED***

		return
	***REMOVED***

	if ctx.credValues.SessionToken != "" ***REMOVED***
		ctx.Request.Header.Set("X-Amz-Security-Token", ctx.credValues.SessionToken)
	***REMOVED***
***REMOVED***

// SignRequestHandler is a named request handler the SDK will use to sign
// service client request with using the V4 signature.
var SignRequestHandler = request.NamedHandler***REMOVED***
	Name: "v4.SignRequestHandler", Fn: SignSDKRequest,
***REMOVED***

// SignSDKRequest signs an AWS request with the V4 signature. This
// request handler should only be used with the SDK's built in service client's
// API operation requests.
//
// This function should not be used on its on its own, but in conjunction with
// an AWS service client's API operation call. To sign a standalone request
// not created by a service client's API operation method use the "Sign" or
// "Presign" functions of the "Signer" type.
//
// If the credentials of the request's config are set to
// credentials.AnonymousCredentials the request will not be signed.
func SignSDKRequest(req *request.Request) ***REMOVED***
	signSDKRequestWithCurrTime(req, time.Now)
***REMOVED***

// BuildNamedHandler will build a generic handler for signing.
func BuildNamedHandler(name string, opts ...func(*Signer)) request.NamedHandler ***REMOVED***
	return request.NamedHandler***REMOVED***
		Name: name,
		Fn: func(req *request.Request) ***REMOVED***
			signSDKRequestWithCurrTime(req, time.Now, opts...)
		***REMOVED***,
	***REMOVED***
***REMOVED***

func signSDKRequestWithCurrTime(req *request.Request, curTimeFn func() time.Time, opts ...func(*Signer)) ***REMOVED***
	// If the request does not need to be signed ignore the signing of the
	// request if the AnonymousCredentials object is used.
	if req.Config.Credentials == credentials.AnonymousCredentials ***REMOVED***
		return
	***REMOVED***

	region := req.ClientInfo.SigningRegion
	if region == "" ***REMOVED***
		region = aws.StringValue(req.Config.Region)
	***REMOVED***

	name := req.ClientInfo.SigningName
	if name == "" ***REMOVED***
		name = req.ClientInfo.ServiceName
	***REMOVED***

	v4 := NewSigner(req.Config.Credentials, func(v4 *Signer) ***REMOVED***
		v4.Debug = req.Config.LogLevel.Value()
		v4.Logger = req.Config.Logger
		v4.DisableHeaderHoisting = req.NotHoist
		v4.currentTimeFn = curTimeFn
		if name == "s3" ***REMOVED***
			// S3 service should not have any escaping applied
			v4.DisableURIPathEscaping = true
		***REMOVED***
		// Prevents setting the HTTPRequest's Body. Since the Body could be
		// wrapped in a custom io.Closer that we do not want to be stompped
		// on top of by the signer.
		v4.DisableRequestBodyOverwrite = true
	***REMOVED***)

	for _, opt := range opts ***REMOVED***
		opt(v4)
	***REMOVED***

	signingTime := req.Time
	if !req.LastSignedAt.IsZero() ***REMOVED***
		signingTime = req.LastSignedAt
	***REMOVED***

	signedHeaders, err := v4.signWithBody(req.HTTPRequest, req.GetBody(),
		name, region, req.ExpireTime, req.ExpireTime > 0, signingTime,
	)
	if err != nil ***REMOVED***
		req.Error = err
		req.SignedHeaderVals = nil
		return
	***REMOVED***

	req.SignedHeaderVals = signedHeaders
	req.LastSignedAt = curTimeFn()
***REMOVED***

const logSignInfoMsg = `DEBUG: Request Signature:
---[ CANONICAL STRING  ]-----------------------------
%s
---[ STRING TO SIGN ]--------------------------------
%s%s
-----------------------------------------------------`
const logSignedURLMsg = `
---[ SIGNED URL ]------------------------------------
%s`

func (v4 *Signer) logSigningInfo(ctx *signingCtx) ***REMOVED***
	signedURLMsg := ""
	if ctx.isPresign ***REMOVED***
		signedURLMsg = fmt.Sprintf(logSignedURLMsg, ctx.Request.URL.String())
	***REMOVED***
	msg := fmt.Sprintf(logSignInfoMsg, ctx.canonicalString, ctx.stringToSign, signedURLMsg)
	v4.Logger.Log(msg)
***REMOVED***

func (ctx *signingCtx) build(disableHeaderHoisting bool) ***REMOVED***
	ctx.buildTime()             // no depends
	ctx.buildCredentialString() // no depends

	ctx.buildBodyDigest()

	unsignedHeaders := ctx.Request.Header
	if ctx.isPresign ***REMOVED***
		if !disableHeaderHoisting ***REMOVED***
			urlValues := url.Values***REMOVED******REMOVED***
			urlValues, unsignedHeaders = buildQuery(allowedQueryHoisting, unsignedHeaders) // no depends
			for k := range urlValues ***REMOVED***
				ctx.Query[k] = urlValues[k]
			***REMOVED***
		***REMOVED***
	***REMOVED***

	ctx.buildCanonicalHeaders(ignoredHeaders, unsignedHeaders)
	ctx.buildCanonicalString() // depends on canon headers / signed headers
	ctx.buildStringToSign()    // depends on canon string
	ctx.buildSignature()       // depends on string to sign

	if ctx.isPresign ***REMOVED***
		ctx.Request.URL.RawQuery += "&X-Amz-Signature=" + ctx.signature
	***REMOVED*** else ***REMOVED***
		parts := []string***REMOVED***
			authHeaderPrefix + " Credential=" + ctx.credValues.AccessKeyID + "/" + ctx.credentialString,
			"SignedHeaders=" + ctx.signedHeaders,
			"Signature=" + ctx.signature,
		***REMOVED***
		ctx.Request.Header.Set("Authorization", strings.Join(parts, ", "))
	***REMOVED***
***REMOVED***

func (ctx *signingCtx) buildTime() ***REMOVED***
	ctx.formattedTime = ctx.Time.UTC().Format(timeFormat)
	ctx.formattedShortTime = ctx.Time.UTC().Format(shortTimeFormat)

	if ctx.isPresign ***REMOVED***
		duration := int64(ctx.ExpireTime / time.Second)
		ctx.Query.Set("X-Amz-Date", ctx.formattedTime)
		ctx.Query.Set("X-Amz-Expires", strconv.FormatInt(duration, 10))
	***REMOVED*** else ***REMOVED***
		ctx.Request.Header.Set("X-Amz-Date", ctx.formattedTime)
	***REMOVED***
***REMOVED***

func (ctx *signingCtx) buildCredentialString() ***REMOVED***
	ctx.credentialString = strings.Join([]string***REMOVED***
		ctx.formattedShortTime,
		ctx.Region,
		ctx.ServiceName,
		"aws4_request",
	***REMOVED***, "/")

	if ctx.isPresign ***REMOVED***
		ctx.Query.Set("X-Amz-Credential", ctx.credValues.AccessKeyID+"/"+ctx.credentialString)
	***REMOVED***
***REMOVED***

func buildQuery(r rule, header http.Header) (url.Values, http.Header) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	unsignedHeaders := http.Header***REMOVED******REMOVED***
	for k, h := range header ***REMOVED***
		if r.IsValid(k) ***REMOVED***
			query[k] = h
		***REMOVED*** else ***REMOVED***
			unsignedHeaders[k] = h
		***REMOVED***
	***REMOVED***

	return query, unsignedHeaders
***REMOVED***
func (ctx *signingCtx) buildCanonicalHeaders(r rule, header http.Header) ***REMOVED***
	var headers []string
	headers = append(headers, "host")
	for k, v := range header ***REMOVED***
		canonicalKey := http.CanonicalHeaderKey(k)
		if !r.IsValid(canonicalKey) ***REMOVED***
			continue // ignored header
		***REMOVED***
		if ctx.SignedHeaderVals == nil ***REMOVED***
			ctx.SignedHeaderVals = make(http.Header)
		***REMOVED***

		lowerCaseKey := strings.ToLower(k)
		if _, ok := ctx.SignedHeaderVals[lowerCaseKey]; ok ***REMOVED***
			// include additional values
			ctx.SignedHeaderVals[lowerCaseKey] = append(ctx.SignedHeaderVals[lowerCaseKey], v...)
			continue
		***REMOVED***

		headers = append(headers, lowerCaseKey)
		ctx.SignedHeaderVals[lowerCaseKey] = v
	***REMOVED***
	sort.Strings(headers)

	ctx.signedHeaders = strings.Join(headers, ";")

	if ctx.isPresign ***REMOVED***
		ctx.Query.Set("X-Amz-SignedHeaders", ctx.signedHeaders)
	***REMOVED***

	headerValues := make([]string, len(headers))
	for i, k := range headers ***REMOVED***
		if k == "host" ***REMOVED***
			if ctx.Request.Host != "" ***REMOVED***
				headerValues[i] = "host:" + ctx.Request.Host
			***REMOVED*** else ***REMOVED***
				headerValues[i] = "host:" + ctx.Request.URL.Host
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			headerValues[i] = k + ":" +
				strings.Join(ctx.SignedHeaderVals[k], ",")
		***REMOVED***
	***REMOVED***
	stripExcessSpaces(headerValues)
	ctx.canonicalHeaders = strings.Join(headerValues, "\n")
***REMOVED***

func (ctx *signingCtx) buildCanonicalString() ***REMOVED***
	ctx.Request.URL.RawQuery = strings.Replace(ctx.Query.Encode(), "+", "%20", -1)

	uri := getURIPath(ctx.Request.URL)

	if !ctx.DisableURIPathEscaping ***REMOVED***
		uri = rest.EscapePath(uri, false)
	***REMOVED***

	ctx.canonicalString = strings.Join([]string***REMOVED***
		ctx.Request.Method,
		uri,
		ctx.Request.URL.RawQuery,
		ctx.canonicalHeaders + "\n",
		ctx.signedHeaders,
		ctx.bodyDigest,
	***REMOVED***, "\n")
***REMOVED***

func (ctx *signingCtx) buildStringToSign() ***REMOVED***
	ctx.stringToSign = strings.Join([]string***REMOVED***
		authHeaderPrefix,
		ctx.formattedTime,
		ctx.credentialString,
		hex.EncodeToString(makeSha256([]byte(ctx.canonicalString))),
	***REMOVED***, "\n")
***REMOVED***

func (ctx *signingCtx) buildSignature() ***REMOVED***
	secret := ctx.credValues.SecretAccessKey
	date := makeHmac([]byte("AWS4"+secret), []byte(ctx.formattedShortTime))
	region := makeHmac(date, []byte(ctx.Region))
	service := makeHmac(region, []byte(ctx.ServiceName))
	credentials := makeHmac(service, []byte("aws4_request"))
	signature := makeHmac(credentials, []byte(ctx.stringToSign))
	ctx.signature = hex.EncodeToString(signature)
***REMOVED***

func (ctx *signingCtx) buildBodyDigest() ***REMOVED***
	hash := ctx.Request.Header.Get("X-Amz-Content-Sha256")
	if hash == "" ***REMOVED***
		if ctx.unsignedPayload || (ctx.isPresign && ctx.ServiceName == "s3") ***REMOVED***
			hash = "UNSIGNED-PAYLOAD"
		***REMOVED*** else if ctx.Body == nil ***REMOVED***
			hash = emptyStringSHA256
		***REMOVED*** else ***REMOVED***
			hash = hex.EncodeToString(makeSha256Reader(ctx.Body))
		***REMOVED***
		if ctx.unsignedPayload || ctx.ServiceName == "s3" || ctx.ServiceName == "glacier" ***REMOVED***
			ctx.Request.Header.Set("X-Amz-Content-Sha256", hash)
		***REMOVED***
	***REMOVED***
	ctx.bodyDigest = hash
***REMOVED***

// isRequestSigned returns if the request is currently signed or presigned
func (ctx *signingCtx) isRequestSigned() bool ***REMOVED***
	if ctx.isPresign && ctx.Query.Get("X-Amz-Signature") != "" ***REMOVED***
		return true
	***REMOVED***
	if ctx.Request.Header.Get("Authorization") != "" ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// unsign removes signing flags for both signed and presigned requests.
func (ctx *signingCtx) removePresign() ***REMOVED***
	ctx.Query.Del("X-Amz-Algorithm")
	ctx.Query.Del("X-Amz-Signature")
	ctx.Query.Del("X-Amz-Security-Token")
	ctx.Query.Del("X-Amz-Date")
	ctx.Query.Del("X-Amz-Expires")
	ctx.Query.Del("X-Amz-Credential")
	ctx.Query.Del("X-Amz-SignedHeaders")
***REMOVED***

func makeHmac(key []byte, data []byte) []byte ***REMOVED***
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
***REMOVED***

func makeSha256(data []byte) []byte ***REMOVED***
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
***REMOVED***

func makeSha256Reader(reader io.ReadSeeker) []byte ***REMOVED***
	hash := sha256.New()
	start, _ := reader.Seek(0, 1)
	defer reader.Seek(start, 0)

	io.Copy(hash, reader)
	return hash.Sum(nil)
***REMOVED***

const doubleSpace = "  "

// stripExcessSpaces will rewrite the passed in slice's string values to not
// contain muliple side-by-side spaces.
func stripExcessSpaces(vals []string) ***REMOVED***
	var j, k, l, m, spaces int
	for i, str := range vals ***REMOVED***
		// Trim trailing spaces
		for j = len(str) - 1; j >= 0 && str[j] == ' '; j-- ***REMOVED***
		***REMOVED***

		// Trim leading spaces
		for k = 0; k < j && str[k] == ' '; k++ ***REMOVED***
		***REMOVED***
		str = str[k : j+1]

		// Strip multiple spaces.
		j = strings.Index(str, doubleSpace)
		if j < 0 ***REMOVED***
			vals[i] = str
			continue
		***REMOVED***

		buf := []byte(str)
		for k, m, l = j, j, len(buf); k < l; k++ ***REMOVED***
			if buf[k] == ' ' ***REMOVED***
				if spaces == 0 ***REMOVED***
					// First space.
					buf[m] = buf[k]
					m++
				***REMOVED***
				spaces++
			***REMOVED*** else ***REMOVED***
				// End of multiple spaces.
				spaces = 0
				buf[m] = buf[k]
				m++
			***REMOVED***
		***REMOVED***

		vals[i] = string(buf[:m])
	***REMOVED***
***REMOVED***
