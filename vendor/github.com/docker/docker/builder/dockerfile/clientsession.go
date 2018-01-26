package dockerfile

import (
	"time"

	"github.com/docker/docker/builder/fscache"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/filesync"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const sessionConnectTimeout = 5 * time.Second

// ClientSessionTransport is a transport for copying files from docker client
// to the daemon.
type ClientSessionTransport struct***REMOVED******REMOVED***

// NewClientSessionTransport returns new ClientSessionTransport instance
func NewClientSessionTransport() *ClientSessionTransport ***REMOVED***
	return &ClientSessionTransport***REMOVED******REMOVED***
***REMOVED***

// Copy data from a remote to a destination directory.
func (cst *ClientSessionTransport) Copy(ctx context.Context, id fscache.RemoteIdentifier, dest string, cu filesync.CacheUpdater) error ***REMOVED***
	csi, ok := id.(*ClientSessionSourceIdentifier)
	if !ok ***REMOVED***
		return errors.New("invalid identifier for client session")
	***REMOVED***

	return filesync.FSSync(ctx, csi.caller, filesync.FSSendRequestOpt***REMOVED***
		IncludePatterns: csi.includePatterns,
		DestDir:         dest,
		CacheUpdater:    cu,
	***REMOVED***)
***REMOVED***

// ClientSessionSourceIdentifier is an identifier that can be used for requesting
// files from remote client
type ClientSessionSourceIdentifier struct ***REMOVED***
	includePatterns []string
	caller          session.Caller
	uuid            string
***REMOVED***

// NewClientSessionSourceIdentifier returns new ClientSessionSourceIdentifier instance
func NewClientSessionSourceIdentifier(ctx context.Context, sg SessionGetter, uuid string) (*ClientSessionSourceIdentifier, error) ***REMOVED***
	csi := &ClientSessionSourceIdentifier***REMOVED***
		uuid: uuid,
	***REMOVED***
	caller, err := sg.Get(ctx, uuid)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to get session for %s", uuid)
	***REMOVED***

	csi.caller = caller
	return csi, nil
***REMOVED***

// Transport returns transport identifier for remote identifier
func (csi *ClientSessionSourceIdentifier) Transport() string ***REMOVED***
	return remotecontext.ClientSessionRemote
***REMOVED***

// SharedKey returns shared key for remote identifier. Shared key is used
// for finding the base for a repeated transfer.
func (csi *ClientSessionSourceIdentifier) SharedKey() string ***REMOVED***
	return csi.caller.SharedKey()
***REMOVED***

// Key returns unique key for remote identifier. Requests with same key return
// same data.
func (csi *ClientSessionSourceIdentifier) Key() string ***REMOVED***
	return csi.uuid
***REMOVED***
