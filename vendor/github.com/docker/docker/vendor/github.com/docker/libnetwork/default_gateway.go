package libnetwork

import (
	"fmt"
	"strings"

	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

const (
	gwEPlen = 12
)

var procGwNetwork = make(chan (bool), 1)

/*
   libnetwork creates a bridge network "docker_gw_bridge" for providing
   default gateway for the containers if none of the container's endpoints
   have GW set by the driver. ICC is set to false for the GW_bridge network.

   If a driver can't provide external connectivity it can choose to not set
   the GW IP for the endpoint.

   endpoint on the GW_bridge network is managed dynamically by libnetwork.
   ie:
   - its created when an endpoint without GW joins the container
   - its deleted when an endpoint with GW joins the container
*/

func (sb *sandbox) setupDefaultGW() error ***REMOVED***

	// check if the container already has a GW endpoint
	if ep := sb.getEndpointInGWNetwork(); ep != nil ***REMOVED***
		return nil
	***REMOVED***

	c := sb.controller

	// Look for default gw network. In case of error (includes not found),
	// retry and create it if needed in a serialized execution.
	n, err := c.NetworkByName(libnGWNetwork)
	if err != nil ***REMOVED***
		if n, err = c.defaultGwNetwork(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	createOptions := []EndpointOption***REMOVED***CreateOptionAnonymous()***REMOVED***

	eplen := gwEPlen
	if len(sb.containerID) < gwEPlen ***REMOVED***
		eplen = len(sb.containerID)
	***REMOVED***

	sbLabels := sb.Labels()

	if sbLabels[netlabel.PortMap] != nil ***REMOVED***
		createOptions = append(createOptions, CreateOptionPortMapping(sbLabels[netlabel.PortMap].([]types.PortBinding)))
	***REMOVED***

	if sbLabels[netlabel.ExposedPorts] != nil ***REMOVED***
		createOptions = append(createOptions, CreateOptionExposedPorts(sbLabels[netlabel.ExposedPorts].([]types.TransportPort)))
	***REMOVED***

	epOption := getPlatformOption()
	if epOption != nil ***REMOVED***
		createOptions = append(createOptions, epOption)
	***REMOVED***

	newEp, err := n.CreateEndpoint("gateway_"+sb.containerID[0:eplen], createOptions...)
	if err != nil ***REMOVED***
		return fmt.Errorf("container %s: endpoint create on GW Network failed: %v", sb.containerID, err)
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if err2 := newEp.Delete(true); err2 != nil ***REMOVED***
				logrus.Warnf("Failed to remove gw endpoint for container %s after failing to join the gateway network: %v",
					sb.containerID, err2)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	epLocal := newEp.(*endpoint)

	if err = epLocal.sbJoin(sb); err != nil ***REMOVED***
		return fmt.Errorf("container %s: endpoint join on GW Network failed: %v", sb.containerID, err)
	***REMOVED***

	return nil
***REMOVED***

// If present, detach and remove the endpoint connecting the sandbox to the default gw network.
func (sb *sandbox) clearDefaultGW() error ***REMOVED***
	var ep *endpoint

	if ep = sb.getEndpointInGWNetwork(); ep == nil ***REMOVED***
		return nil
	***REMOVED***
	if err := ep.sbLeave(sb, false); err != nil ***REMOVED***
		return fmt.Errorf("container %s: endpoint leaving GW Network failed: %v", sb.containerID, err)
	***REMOVED***
	if err := ep.Delete(false); err != nil ***REMOVED***
		return fmt.Errorf("container %s: deleting endpoint on GW Network failed: %v", sb.containerID, err)
	***REMOVED***
	return nil
***REMOVED***

// Evaluate whether the sandbox requires a default gateway based
// on the endpoints to which it is connected. It does not account
// for the default gateway network endpoint.

func (sb *sandbox) needDefaultGW() bool ***REMOVED***
	var needGW bool

	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if ep.endpointInGWNetwork() ***REMOVED***
			continue
		***REMOVED***
		if ep.getNetwork().Type() == "null" || ep.getNetwork().Type() == "host" ***REMOVED***
			continue
		***REMOVED***
		if ep.getNetwork().Internal() ***REMOVED***
			continue
		***REMOVED***
		// During stale sandbox cleanup, joinInfo may be nil
		if ep.joinInfo != nil && ep.joinInfo.disableGatewayService ***REMOVED***
			continue
		***REMOVED***
		// TODO v6 needs to be handled.
		if len(ep.Gateway()) > 0 ***REMOVED***
			return false
		***REMOVED***
		for _, r := range ep.StaticRoutes() ***REMOVED***
			if r.Destination != nil && r.Destination.String() == "0.0.0.0/0" ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		needGW = true
	***REMOVED***

	return needGW
***REMOVED***

func (sb *sandbox) getEndpointInGWNetwork() *endpoint ***REMOVED***
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if ep.getNetwork().name == libnGWNetwork && strings.HasPrefix(ep.Name(), "gateway_") ***REMOVED***
			return ep
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (ep *endpoint) endpointInGWNetwork() bool ***REMOVED***
	if ep.getNetwork().name == libnGWNetwork && strings.HasPrefix(ep.Name(), "gateway_") ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (sb *sandbox) getEPwithoutGateway() *endpoint ***REMOVED***
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if ep.getNetwork().Type() == "null" || ep.getNetwork().Type() == "host" ***REMOVED***
			continue
		***REMOVED***
		if len(ep.Gateway()) == 0 ***REMOVED***
			return ep
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Looks for the default gw network and creates it if not there.
// Parallel executions are serialized.
func (c *controller) defaultGwNetwork() (Network, error) ***REMOVED***
	procGwNetwork <- true
	defer func() ***REMOVED*** <-procGwNetwork ***REMOVED***()

	n, err := c.NetworkByName(libnGWNetwork)
	if err != nil ***REMOVED***
		if _, ok := err.(types.NotFoundError); ok ***REMOVED***
			n, err = c.createGWNetwork()
		***REMOVED***
	***REMOVED***
	return n, err
***REMOVED***

// Returns the endpoint which is providing external connectivity to the sandbox
func (sb *sandbox) getGatewayEndpoint() *endpoint ***REMOVED***
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		if ep.getNetwork().Type() == "null" || ep.getNetwork().Type() == "host" ***REMOVED***
			continue
		***REMOVED***
		if len(ep.Gateway()) != 0 ***REMOVED***
			return ep
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
