package macvlan

import (
	"fmt"

	"github.com/docker/libnetwork/osl"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

func (d *driver) network(nid string) *network ***REMOVED***
	d.Lock()
	n, ok := d.networks[nid]
	d.Unlock()
	if !ok ***REMOVED***
		logrus.Errorf("network id %s not found", nid)
	***REMOVED***

	return n
***REMOVED***

func (d *driver) addNetwork(n *network) ***REMOVED***
	d.Lock()
	d.networks[n.id] = n
	d.Unlock()
***REMOVED***

func (d *driver) deleteNetwork(nid string) ***REMOVED***
	d.Lock()
	delete(d.networks, nid)
	d.Unlock()
***REMOVED***

// getNetworks Safely returns a slice of existing networks
func (d *driver) getNetworks() []*network ***REMOVED***
	d.Lock()
	defer d.Unlock()

	ls := make([]*network, 0, len(d.networks))
	for _, nw := range d.networks ***REMOVED***
		ls = append(ls, nw)
	***REMOVED***

	return ls
***REMOVED***

func (n *network) endpoint(eid string) *endpoint ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.endpoints[eid]
***REMOVED***

func (n *network) addEndpoint(ep *endpoint) ***REMOVED***
	n.Lock()
	n.endpoints[ep.id] = ep
	n.Unlock()
***REMOVED***

func (n *network) deleteEndpoint(eid string) ***REMOVED***
	n.Lock()
	delete(n.endpoints, eid)
	n.Unlock()
***REMOVED***

func (n *network) getEndpoint(eid string) (*endpoint, error) ***REMOVED***
	n.Lock()
	defer n.Unlock()
	if eid == "" ***REMOVED***
		return nil, fmt.Errorf("endpoint id %s not found", eid)
	***REMOVED***
	if ep, ok := n.endpoints[eid]; ok ***REMOVED***
		return ep, nil
	***REMOVED***

	return nil, nil
***REMOVED***

func validateID(nid, eid string) error ***REMOVED***
	if nid == "" ***REMOVED***
		return fmt.Errorf("invalid network id")
	***REMOVED***
	if eid == "" ***REMOVED***
		return fmt.Errorf("invalid endpoint id")
	***REMOVED***
	return nil
***REMOVED***

func (n *network) sandbox() osl.Sandbox ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.sbox
***REMOVED***

func (n *network) setSandbox(sbox osl.Sandbox) ***REMOVED***
	n.Lock()
	n.sbox = sbox
	n.Unlock()
***REMOVED***

func (d *driver) getNetwork(id string) (*network, error) ***REMOVED***
	d.Lock()
	defer d.Unlock()
	if id == "" ***REMOVED***
		return nil, types.BadRequestErrorf("invalid network id: %s", id)
	***REMOVED***
	if nw, ok := d.networks[id]; ok ***REMOVED***
		return nw, nil
	***REMOVED***

	return nil, types.NotFoundErrorf("network not found: %s", id)
***REMOVED***
