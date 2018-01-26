package ca

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	certForwardedKey = "forwarded_cert"
	certCNKey        = "forwarded_cert_cn"
	certOUKey        = "forwarded_cert_ou"
	certOrgKey       = "forwarded_cert_org"
	remoteAddrKey    = "remote_addr"
)

// forwardedTLSInfoFromContext obtains forwarded TLS CN/OU from the grpc.MD
// object in ctx.
func forwardedTLSInfoFromContext(ctx context.Context) (remoteAddr string, cn string, org string, ous []string) ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	if len(md[remoteAddrKey]) != 0 ***REMOVED***
		remoteAddr = md[remoteAddrKey][0]
	***REMOVED***
	if len(md[certCNKey]) != 0 ***REMOVED***
		cn = md[certCNKey][0]
	***REMOVED***
	if len(md[certOrgKey]) != 0 ***REMOVED***
		org = md[certOrgKey][0]
	***REMOVED***
	ous = md[certOUKey]
	return
***REMOVED***

func isForwardedRequest(ctx context.Context) bool ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	if len(md[certForwardedKey]) != 1 ***REMOVED***
		return false
	***REMOVED***
	return md[certForwardedKey][0] == "true"
***REMOVED***

// WithMetadataForwardTLSInfo reads certificate from context and returns context where
// ForwardCert is set based on original certificate.
func WithMetadataForwardTLSInfo(ctx context.Context) (context.Context, error) ***REMOVED***
	md, ok := metadata.FromContext(ctx)
	if !ok ***REMOVED***
		md = metadata.MD***REMOVED******REMOVED***
	***REMOVED***

	ous := []string***REMOVED******REMOVED***
	org := ""
	cn := ""

	certSubj, err := certSubjectFromContext(ctx)
	if err == nil ***REMOVED***
		cn = certSubj.CommonName
		ous = certSubj.OrganizationalUnit
		if len(certSubj.Organization) > 0 ***REMOVED***
			org = certSubj.Organization[0]
		***REMOVED***
	***REMOVED***

	// If there's no TLS cert, forward with blank TLS metadata.
	// Note that the presence of this blank metadata is extremely
	// important. Without it, it would look like manager is making
	// the request directly.
	md[certForwardedKey] = []string***REMOVED***"true"***REMOVED***
	md[certCNKey] = []string***REMOVED***cn***REMOVED***
	md[certOrgKey] = []string***REMOVED***org***REMOVED***
	md[certOUKey] = ous
	peer, ok := peer.FromContext(ctx)
	if ok ***REMOVED***
		md[remoteAddrKey] = []string***REMOVED***peer.Addr.String()***REMOVED***
	***REMOVED***

	return metadata.NewContext(ctx, md), nil
***REMOVED***
