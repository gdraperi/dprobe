package agent

import (
	"github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

type resourceAllocator struct ***REMOVED***
	agent *Agent
***REMOVED***

// ResourceAllocator is an interface to allocate resource such as
// network attachments from a worker node.
type ResourceAllocator interface ***REMOVED***
	// AttachNetwork creates a network attachment in the manager
	// given a target network and a unique ID representing the
	// connecting entity and optionally a list of ipv4/ipv6
	// addresses to be assigned to the attachment. AttachNetwork
	// returns a unique ID for the attachment if successful or an
	// error in case of failure.
	AttachNetwork(ctx context.Context, id, target string, addresses []string) (string, error)

	// DetachNetworks deletes a network attachment for the passed
	// attachment ID. The attachment ID is obtained from a
	// previous AttachNetwork call.
	DetachNetwork(ctx context.Context, aID string) error
***REMOVED***

// AttachNetwork creates a network attachment.
func (r *resourceAllocator) AttachNetwork(ctx context.Context, id, target string, addresses []string) (string, error) ***REMOVED***
	var taskID string
	if err := r.agent.withSession(ctx, func(session *session) error ***REMOVED***
		client := api.NewResourceAllocatorClient(session.conn.ClientConn)
		r, err := client.AttachNetwork(ctx, &api.AttachNetworkRequest***REMOVED***
			Config: &api.NetworkAttachmentConfig***REMOVED***
				Target:    target,
				Addresses: addresses,
			***REMOVED***,
			ContainerID: id,
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		taskID = r.AttachmentID
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return taskID, nil
***REMOVED***

// DetachNetwork deletes a network attachment.
func (r *resourceAllocator) DetachNetwork(ctx context.Context, aID string) error ***REMOVED***
	return r.agent.withSession(ctx, func(session *session) error ***REMOVED***
		client := api.NewResourceAllocatorClient(session.conn.ClientConn)
		_, err := client.DetachNetwork(ctx, &api.DetachNetworkRequest***REMOVED***
			AttachmentID: aID,
		***REMOVED***)

		return err
	***REMOVED***)
***REMOVED***

// ResourceAllocator provides an interface to access resource
// allocation methods such as AttachNetwork and DetachNetwork.
func (a *Agent) ResourceAllocator() ResourceAllocator ***REMOVED***
	return &resourceAllocator***REMOVED***agent: a***REMOVED***
***REMOVED***
