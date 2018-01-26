package overlay

import (
	"fmt"
	"net"

	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/types"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// Join method is invoked when a Sandbox is attached to an endpoint.
func (d *driver) Join(nid, eid string, sboxKey string, jinfo driverapi.JoinInfo, options map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return fmt.Errorf("could not find network with id %s", nid)
	***REMOVED***

	ep := n.endpoint(eid)
	if ep == nil ***REMOVED***
		return fmt.Errorf("could not find endpoint with id %s", eid)
	***REMOVED***

	buf, err := proto.Marshal(&PeerRecord***REMOVED***
		EndpointIP:       ep.addr.String(),
		EndpointMAC:      ep.mac.String(),
		TunnelEndpointIP: n.providerAddress,
	***REMOVED***)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := jinfo.AddTableEntry(ovPeerTable, eid, buf); err != nil ***REMOVED***
		logrus.Errorf("overlay: Failed adding table entry to joininfo: %v", err)
	***REMOVED***

	if ep.disablegateway ***REMOVED***
		jinfo.DisableGatewayService()
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) EventNotify(etype driverapi.EventType, nid, tableName, key string, value []byte) ***REMOVED***
	if tableName != ovPeerTable ***REMOVED***
		logrus.Errorf("Unexpected table notification for table %s received", tableName)
		return
	***REMOVED***

	eid := key

	var peer PeerRecord
	if err := proto.Unmarshal(value, &peer); err != nil ***REMOVED***
		logrus.Errorf("Failed to unmarshal peer record: %v", err)
		return
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return
	***REMOVED***

	// Ignore local peers. We already know about them and they
	// should not be added to vxlan fdb.
	if peer.TunnelEndpointIP == n.providerAddress ***REMOVED***
		return
	***REMOVED***

	addr, err := types.ParseCIDR(peer.EndpointIP)
	if err != nil ***REMOVED***
		logrus.Errorf("Invalid peer IP %s received in event notify", peer.EndpointIP)
		return
	***REMOVED***

	mac, err := net.ParseMAC(peer.EndpointMAC)
	if err != nil ***REMOVED***
		logrus.Errorf("Invalid mac %s received in event notify", peer.EndpointMAC)
		return
	***REMOVED***

	vtep := net.ParseIP(peer.TunnelEndpointIP)
	if vtep == nil ***REMOVED***
		logrus.Errorf("Invalid VTEP %s received in event notify", peer.TunnelEndpointIP)
		return
	***REMOVED***

	if etype == driverapi.Delete ***REMOVED***
		d.peerDelete(nid, eid, addr.IP, addr.Mask, mac, vtep, true)
		return
	***REMOVED***

	d.peerAdd(nid, eid, addr.IP, addr.Mask, mac, vtep, true)
***REMOVED***

func (d *driver) DecodeTableEntry(tablename string, key string, value []byte) (string, map[string]string) ***REMOVED***
	return "", nil
***REMOVED***

// Leave method is invoked when a Sandbox detaches from an endpoint.
func (d *driver) Leave(nid, eid string) error ***REMOVED***
	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
