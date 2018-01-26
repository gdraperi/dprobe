package ca

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/cloudflare/cfssl/helpers"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/equality"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/pkg/errors"
)

// IssuanceStateRotateMaxBatchSize is the maximum number of nodes we'll tell to rotate their certificates in any given update
const IssuanceStateRotateMaxBatchSize = 30

func hasIssuer(n *api.Node, info *IssuerInfo) bool ***REMOVED***
	if n.Description == nil || n.Description.TLSInfo == nil ***REMOVED***
		return false
	***REMOVED***
	return bytes.Equal(info.Subject, n.Description.TLSInfo.CertIssuerSubject) && bytes.Equal(info.PublicKey, n.Description.TLSInfo.CertIssuerPublicKey)
***REMOVED***

var errRootRotationChanged = errors.New("target root rotation has changed")

// rootRotationReconciler keeps track of all the nodes in the store so that we can determine which ones need reconciliation when nodes are updated
// or the root CA is updated.  This is meant to be used with watches on nodes and the cluster, and provides functions to be called when the
// cluster's RootCA has changed and when a node is added, updated, or removed.
type rootRotationReconciler struct ***REMOVED***
	mu                  sync.Mutex
	clusterID           string
	batchUpdateInterval time.Duration
	ctx                 context.Context
	store               *store.MemoryStore

	currentRootCA    *api.RootCA
	currentIssuer    IssuerInfo
	unconvergedNodes map[string]*api.Node

	wg     sync.WaitGroup
	cancel func()
***REMOVED***

// IssuerFromAPIRootCA returns the desired issuer given an API root CA object
func IssuerFromAPIRootCA(rootCA *api.RootCA) (*IssuerInfo, error) ***REMOVED***
	wantedIssuer := rootCA.CACert
	if rootCA.RootRotation != nil ***REMOVED***
		wantedIssuer = rootCA.RootRotation.CACert
	***REMOVED***
	issuerCerts, err := helpers.ParseCertificatesPEM(wantedIssuer)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "invalid certificate in cluster root CA object")
	***REMOVED***
	if len(issuerCerts) == 0 ***REMOVED***
		return nil, errors.New("invalid certificate in cluster root CA object")
	***REMOVED***
	return &IssuerInfo***REMOVED***
		Subject:   issuerCerts[0].RawSubject,
		PublicKey: issuerCerts[0].RawSubjectPublicKeyInfo,
	***REMOVED***, nil
***REMOVED***

// assumption:  UpdateRootCA will never be called with a `nil` root CA because the caller will be acting in response to
// a store update event
func (r *rootRotationReconciler) UpdateRootCA(newRootCA *api.RootCA) ***REMOVED***
	issuerInfo, err := IssuerFromAPIRootCA(newRootCA)
	if err != nil ***REMOVED***
		log.G(r.ctx).WithError(err).Error("unable to update process the current root CA")
		return
	***REMOVED***

	var (
		shouldStartNewLoop, waitForPrevLoop bool
		loopCtx                             context.Context
	)
	r.mu.Lock()
	defer func() ***REMOVED***
		r.mu.Unlock()
		if shouldStartNewLoop ***REMOVED***
			if waitForPrevLoop ***REMOVED***
				r.wg.Wait()
			***REMOVED***
			r.wg.Add(1)
			go r.runReconcilerLoop(loopCtx, newRootCA)
		***REMOVED***
	***REMOVED***()

	// check if the issuer has changed, first
	if reflect.DeepEqual(&r.currentIssuer, issuerInfo) ***REMOVED***
		r.currentRootCA = newRootCA
		return
	***REMOVED***
	// If the issuer has changed, iterate through all the nodes to figure out which ones need rotation
	if newRootCA.RootRotation != nil ***REMOVED***
		var nodes []*api.Node
		r.store.View(func(tx store.ReadTx) ***REMOVED***
			nodes, err = store.FindNodes(tx, store.All)
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(r.ctx).WithError(err).Error("unable to list nodes, so unable to process the current root CA")
			return
		***REMOVED***

		// from here on out, there will be no more errors that cause us to have to abandon updating the Root CA,
		// so we can start making changes to r's fields
		r.unconvergedNodes = make(map[string]*api.Node)
		for _, n := range nodes ***REMOVED***
			if !hasIssuer(n, issuerInfo) ***REMOVED***
				r.unconvergedNodes[n.ID] = n
			***REMOVED***
		***REMOVED***
		shouldStartNewLoop = true
		if r.cancel != nil ***REMOVED*** // there's already a loop going, so cancel it
			r.cancel()
			waitForPrevLoop = true
		***REMOVED***
		loopCtx, r.cancel = context.WithCancel(r.ctx)
	***REMOVED*** else ***REMOVED***
		r.unconvergedNodes = nil
	***REMOVED***
	r.currentRootCA = newRootCA
	r.currentIssuer = *issuerInfo
***REMOVED***

// assumption:  UpdateNode will never be called with a `nil` node because the caller will be acting in response to
// a store update event
func (r *rootRotationReconciler) UpdateNode(node *api.Node) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	// if we're not in the middle of a root rotation ignore the update
	if r.currentRootCA == nil || r.currentRootCA.RootRotation == nil ***REMOVED***
		return
	***REMOVED***
	if hasIssuer(node, &r.currentIssuer) ***REMOVED***
		delete(r.unconvergedNodes, node.ID)
	***REMOVED*** else ***REMOVED***
		r.unconvergedNodes[node.ID] = node
	***REMOVED***
***REMOVED***

// assumption:  DeleteNode will never be called with a `nil` node because the caller will be acting in response to
// a store update event
func (r *rootRotationReconciler) DeleteNode(node *api.Node) ***REMOVED***
	r.mu.Lock()
	delete(r.unconvergedNodes, node.ID)
	r.mu.Unlock()
***REMOVED***

func (r *rootRotationReconciler) runReconcilerLoop(ctx context.Context, loopRootCA *api.RootCA) ***REMOVED***
	defer r.wg.Done()
	for ***REMOVED***
		r.mu.Lock()
		if len(r.unconvergedNodes) == 0 ***REMOVED***
			r.mu.Unlock()

			err := r.store.Update(func(tx store.Tx) error ***REMOVED***
				return r.finishRootRotation(tx, loopRootCA)
			***REMOVED***)
			if err == nil ***REMOVED***
				log.G(r.ctx).Info("completed root rotation")
				return
			***REMOVED***
			log.G(r.ctx).WithError(err).Error("could not complete root rotation")
			if err == errRootRotationChanged ***REMOVED***
				// if the root rotation has changed, this loop will be cancelled anyway, so may as well abort early
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			var toUpdate []*api.Node
			for _, n := range r.unconvergedNodes ***REMOVED***
				iState := n.Certificate.Status.State
				if iState != api.IssuanceStateRenew && iState != api.IssuanceStatePending && iState != api.IssuanceStateRotate ***REMOVED***
					n = n.Copy()
					n.Certificate.Status.State = api.IssuanceStateRotate
					toUpdate = append(toUpdate, n)
					if len(toUpdate) >= IssuanceStateRotateMaxBatchSize ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
			r.mu.Unlock()

			if err := r.batchUpdateNodes(toUpdate); err != nil ***REMOVED***
				log.G(r.ctx).WithError(err).Errorf("store error when trying to batch update %d nodes to request certificate rotation", len(toUpdate))
			***REMOVED***
		***REMOVED***

		select ***REMOVED***
		case <-ctx.Done():
			return
		case <-time.After(r.batchUpdateInterval):
		***REMOVED***
	***REMOVED***
***REMOVED***

// This function assumes that the expected root CA has root rotation.  This is intended to be used by
// `reconcileNodeRootsAndCerts`, which uses the root CA from the `lastSeenClusterRootCA`, and checks
// that it has a root rotation before calling this function.
func (r *rootRotationReconciler) finishRootRotation(tx store.Tx, expectedRootCA *api.RootCA) error ***REMOVED***
	cluster := store.GetCluster(tx, r.clusterID)
	if cluster == nil ***REMOVED***
		return fmt.Errorf("unable to get cluster %s", r.clusterID)
	***REMOVED***

	// If the RootCA object has changed (because another root rotation was started or because some other node
	// had finished the root rotation), we cannot finish the root rotation that we were working on.
	if !equality.RootCAEqualStable(expectedRootCA, &cluster.RootCA) ***REMOVED***
		return errRootRotationChanged
	***REMOVED***

	var signerCert []byte
	if len(cluster.RootCA.RootRotation.CAKey) > 0 ***REMOVED***
		signerCert = cluster.RootCA.RootRotation.CACert
	***REMOVED***
	// we don't actually have to parse out the default node expiration from the cluster - we are just using
	// the ca.RootCA object to generate new tokens and the digest
	updatedRootCA, err := NewRootCA(cluster.RootCA.RootRotation.CACert, signerCert, cluster.RootCA.RootRotation.CAKey,
		DefaultNodeCertExpiration, nil)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "invalid cluster root rotation object")
	***REMOVED***
	cluster.RootCA = api.RootCA***REMOVED***
		CACert:     cluster.RootCA.RootRotation.CACert,
		CAKey:      cluster.RootCA.RootRotation.CAKey,
		CACertHash: updatedRootCA.Digest.String(),
		JoinTokens: api.JoinTokens***REMOVED***
			Worker:  GenerateJoinToken(&updatedRootCA),
			Manager: GenerateJoinToken(&updatedRootCA),
		***REMOVED***,
		LastForcedRotation: cluster.RootCA.LastForcedRotation,
	***REMOVED***
	return store.UpdateCluster(tx, cluster)
***REMOVED***

func (r *rootRotationReconciler) batchUpdateNodes(toUpdate []*api.Node) error ***REMOVED***
	if len(toUpdate) == 0 ***REMOVED***
		return nil
	***REMOVED***
	err := r.store.Batch(func(batch *store.Batch) error ***REMOVED***
		// Directly update the nodes rather than get + update, and ignore version errors.  Since
		// `rootRotationReconciler` should be hooked up to all node update/delete/create events, we should have
		// close to the latest versions of all the nodes.  If not, the node will updated later and the
		// next batch of updates should catch it.
		for _, n := range toUpdate ***REMOVED***
			if err := batch.Update(func(tx store.Tx) error ***REMOVED***
				return store.UpdateNode(tx, n)
			***REMOVED***); err != nil && err != store.ErrSequenceConflict ***REMOVED***
				log.G(r.ctx).WithError(err).Errorf("unable to update node %s to request a certificate rotation", n.ID)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	return err
***REMOVED***
