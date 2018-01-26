package ca

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type localRequestKeyType struct***REMOVED******REMOVED***

// LocalRequestKey is a context key to mark a request that originating on the
// local node. The associated value is a RemoteNodeInfo structure describing the
// local node.
var LocalRequestKey = localRequestKeyType***REMOVED******REMOVED***

// LogTLSState logs information about the TLS connection and remote peers
func LogTLSState(ctx context.Context, tlsState *tls.ConnectionState) ***REMOVED***
	if tlsState == nil ***REMOVED***
		log.G(ctx).Debugf("no TLS Chains found")
		return
	***REMOVED***

	peerCerts := []string***REMOVED******REMOVED***
	verifiedChain := []string***REMOVED******REMOVED***
	for _, cert := range tlsState.PeerCertificates ***REMOVED***
		peerCerts = append(peerCerts, cert.Subject.CommonName)
	***REMOVED***
	for _, chain := range tlsState.VerifiedChains ***REMOVED***
		subjects := []string***REMOVED******REMOVED***
		for _, cert := range chain ***REMOVED***
			subjects = append(subjects, cert.Subject.CommonName)
		***REMOVED***
		verifiedChain = append(verifiedChain, strings.Join(subjects, ","))
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"peer.peerCert": peerCerts,
		// "peer.verifiedChain": verifiedChain***REMOVED***,
	***REMOVED***).Debugf("")
***REMOVED***

// getCertificateSubject extracts the subject from a verified client certificate
func getCertificateSubject(tlsState *tls.ConnectionState) (pkix.Name, error) ***REMOVED***
	if tlsState == nil ***REMOVED***
		return pkix.Name***REMOVED******REMOVED***, status.Errorf(codes.PermissionDenied, "request is not using TLS")
	***REMOVED***
	if len(tlsState.PeerCertificates) == 0 ***REMOVED***
		return pkix.Name***REMOVED******REMOVED***, status.Errorf(codes.PermissionDenied, "no client certificates in request")
	***REMOVED***
	if len(tlsState.VerifiedChains) == 0 ***REMOVED***
		return pkix.Name***REMOVED******REMOVED***, status.Errorf(codes.PermissionDenied, "no verified chains for remote certificate")
	***REMOVED***

	return tlsState.VerifiedChains[0][0].Subject, nil
***REMOVED***

func tlsConnStateFromContext(ctx context.Context) (*tls.ConnectionState, error) ***REMOVED***
	peer, ok := peer.FromContext(ctx)
	if !ok ***REMOVED***
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied: no peer info")
	***REMOVED***
	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok ***REMOVED***
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied: peer didn't not present valid peer certificate")
	***REMOVED***
	return &tlsInfo.State, nil
***REMOVED***

// certSubjectFromContext extracts pkix.Name from context.
func certSubjectFromContext(ctx context.Context) (pkix.Name, error) ***REMOVED***
	connState, err := tlsConnStateFromContext(ctx)
	if err != nil ***REMOVED***
		return pkix.Name***REMOVED******REMOVED***, err
	***REMOVED***
	return getCertificateSubject(connState)
***REMOVED***

// AuthorizeOrgAndRole takes in a context and a list of roles, and returns
// the Node ID of the node.
func AuthorizeOrgAndRole(ctx context.Context, org string, blacklistedCerts map[string]*api.BlacklistedCertificate, ou ...string) (string, error) ***REMOVED***
	certSubj, err := certSubjectFromContext(ctx)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// Check if the current certificate has an OU that authorizes
	// access to this method
	if intersectArrays(certSubj.OrganizationalUnit, ou) ***REMOVED***
		return authorizeOrg(certSubj, org, blacklistedCerts)
	***REMOVED***

	return "", status.Errorf(codes.PermissionDenied, "Permission denied: remote certificate not part of OUs: %v", ou)
***REMOVED***

// authorizeOrg takes in a certificate subject and an organization, and returns
// the Node ID of the node.
func authorizeOrg(certSubj pkix.Name, org string, blacklistedCerts map[string]*api.BlacklistedCertificate) (string, error) ***REMOVED***
	if _, ok := blacklistedCerts[certSubj.CommonName]; ok ***REMOVED***
		return "", status.Errorf(codes.PermissionDenied, "Permission denied: node %s was removed from swarm", certSubj.CommonName)
	***REMOVED***

	if len(certSubj.Organization) > 0 && certSubj.Organization[0] == org ***REMOVED***
		return certSubj.CommonName, nil
	***REMOVED***

	return "", status.Errorf(codes.PermissionDenied, "Permission denied: remote certificate not part of organization: %s", org)
***REMOVED***

// AuthorizeForwardedRoleAndOrg checks for proper roles and organization of caller. The RPC may have
// been proxied by a manager, in which case the manager is authenticated and
// so is the certificate information that it forwarded. It returns the node ID
// of the original client.
func AuthorizeForwardedRoleAndOrg(ctx context.Context, authorizedRoles, forwarderRoles []string, org string, blacklistedCerts map[string]*api.BlacklistedCertificate) (string, error) ***REMOVED***
	if isForwardedRequest(ctx) ***REMOVED***
		_, err := AuthorizeOrgAndRole(ctx, org, blacklistedCerts, forwarderRoles...)
		if err != nil ***REMOVED***
			return "", status.Errorf(codes.PermissionDenied, "Permission denied: unauthorized forwarder role: %v", err)
		***REMOVED***

		// This was a forwarded request. Authorize the forwarder, and
		// check if the forwarded role matches one of the authorized
		// roles.
		_, forwardedID, forwardedOrg, forwardedOUs := forwardedTLSInfoFromContext(ctx)

		if len(forwardedOUs) == 0 || forwardedID == "" || forwardedOrg == "" ***REMOVED***
			return "", status.Errorf(codes.PermissionDenied, "Permission denied: missing information in forwarded request")
		***REMOVED***

		if !intersectArrays(forwardedOUs, authorizedRoles) ***REMOVED***
			return "", status.Errorf(codes.PermissionDenied, "Permission denied: unauthorized forwarded role, expecting: %v", authorizedRoles)
		***REMOVED***

		if forwardedOrg != org ***REMOVED***
			return "", status.Errorf(codes.PermissionDenied, "Permission denied: organization mismatch, expecting: %s", org)
		***REMOVED***

		return forwardedID, nil
	***REMOVED***

	// There wasn't any node being forwarded, check if this is a direct call by the expected role
	nodeID, err := AuthorizeOrgAndRole(ctx, org, blacklistedCerts, authorizedRoles...)
	if err == nil ***REMOVED***
		return nodeID, nil
	***REMOVED***

	return "", status.Errorf(codes.PermissionDenied, "Permission denied: unauthorized peer role: %v", err)
***REMOVED***

// intersectArrays returns true when there is at least one element in common
// between the two arrays
func intersectArrays(orig, tgt []string) bool ***REMOVED***
	for _, i := range orig ***REMOVED***
		for _, x := range tgt ***REMOVED***
			if i == x ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// RemoteNodeInfo describes a node sending an RPC request.
type RemoteNodeInfo struct ***REMOVED***
	// Roles is a list of roles contained in the node's certificate
	// (or forwarded by a trusted node).
	Roles []string

	// Organization is the organization contained in the node's certificate
	// (or forwarded by a trusted node).
	Organization string

	// NodeID is the node's ID, from the CN field in its certificate
	// (or forwarded by a trusted node).
	NodeID string

	// ForwardedBy contains information for the node that forwarded this
	// request. It is set to nil if the request was received directly.
	ForwardedBy *RemoteNodeInfo

	// RemoteAddr is the address that this node is connecting to the cluster
	// from.
	RemoteAddr string
***REMOVED***

// RemoteNode returns the node ID and role from the client's TLS certificate.
// If the RPC was forwarded, the original client's ID and role is returned, as
// well as the forwarder's ID. This function does not do authorization checks -
// it only looks up the node ID.
func RemoteNode(ctx context.Context) (RemoteNodeInfo, error) ***REMOVED***
	// If we have a value on the context that marks this as a local
	// request, we return the node info from the context.
	localNodeInfo := ctx.Value(LocalRequestKey)

	if localNodeInfo != nil ***REMOVED***
		nodeInfo, ok := localNodeInfo.(RemoteNodeInfo)
		if ok ***REMOVED***
			return nodeInfo, nil
		***REMOVED***
	***REMOVED***

	certSubj, err := certSubjectFromContext(ctx)
	if err != nil ***REMOVED***
		return RemoteNodeInfo***REMOVED******REMOVED***, err
	***REMOVED***

	org := ""
	if len(certSubj.Organization) > 0 ***REMOVED***
		org = certSubj.Organization[0]
	***REMOVED***

	peer, ok := peer.FromContext(ctx)
	if !ok ***REMOVED***
		return RemoteNodeInfo***REMOVED******REMOVED***, status.Errorf(codes.PermissionDenied, "Permission denied: no peer info")
	***REMOVED***

	directInfo := RemoteNodeInfo***REMOVED***
		Roles:        certSubj.OrganizationalUnit,
		NodeID:       certSubj.CommonName,
		Organization: org,
		RemoteAddr:   peer.Addr.String(),
	***REMOVED***

	if isForwardedRequest(ctx) ***REMOVED***
		remoteAddr, cn, org, ous := forwardedTLSInfoFromContext(ctx)
		if len(ous) == 0 || cn == "" || org == "" ***REMOVED***
			return RemoteNodeInfo***REMOVED******REMOVED***, status.Errorf(codes.PermissionDenied, "Permission denied: missing information in forwarded request")
		***REMOVED***
		return RemoteNodeInfo***REMOVED***
			Roles:        ous,
			NodeID:       cn,
			Organization: org,
			ForwardedBy:  &directInfo,
			RemoteAddr:   remoteAddr,
		***REMOVED***, nil
	***REMOVED***

	return directInfo, nil
***REMOVED***
