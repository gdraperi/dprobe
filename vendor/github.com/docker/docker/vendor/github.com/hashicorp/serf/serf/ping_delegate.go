package serf

import (
	"bytes"
	"log"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/coordinate"
)

// pingDelegate is notified when memberlist successfully completes a direct ping
// of a peer node. We use this to update our estimated network coordinate, as
// well as cache the coordinate of the peer.
type pingDelegate struct ***REMOVED***
	serf *Serf
***REMOVED***

const (
	// PingVersion is an internal version for the ping message, above the normal
	// versioning we get from the protocol version. This enables small updates
	// to the ping message without a full protocol bump.
	PingVersion = 1
)

// AckPayload is called to produce a payload to send back in response to a ping
// request.
func (p *pingDelegate) AckPayload() []byte ***REMOVED***
	var buf bytes.Buffer

	// The first byte is the version number, forming a simple header.
	version := []byte***REMOVED***PingVersion***REMOVED***
	buf.Write(version)

	// The rest of the message is the serialized coordinate.
	enc := codec.NewEncoder(&buf, &codec.MsgpackHandle***REMOVED******REMOVED***)
	if err := enc.Encode(p.serf.coordClient.GetCoordinate()); err != nil ***REMOVED***
		log.Printf("[ERR] serf: Failed to encode coordinate: %v\n", err)
	***REMOVED***
	return buf.Bytes()
***REMOVED***

// NotifyPingComplete is called when this node successfully completes a direct ping
// of a peer node.
func (p *pingDelegate) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) ***REMOVED***
	if payload == nil || len(payload) == 0 ***REMOVED***
		return
	***REMOVED***

	// Verify ping version in the header.
	version := payload[0]
	if version != PingVersion ***REMOVED***
		log.Printf("[ERR] serf: Unsupported ping version: %v", version)
		return
	***REMOVED***

	// Process the remainder of the message as a coordinate.
	r := bytes.NewReader(payload[1:])
	dec := codec.NewDecoder(r, &codec.MsgpackHandle***REMOVED******REMOVED***)
	var coord coordinate.Coordinate
	if err := dec.Decode(&coord); err != nil ***REMOVED***
		log.Printf("[ERR] serf: Failed to decode coordinate from ping: %v", err)
	***REMOVED***

	// Apply the update. Since this is a coordinate coming from some place
	// else we harden this and look for dimensionality problems proactively.
	before := p.serf.coordClient.GetCoordinate()
	if before.IsCompatibleWith(&coord) ***REMOVED***
		after := p.serf.coordClient.Update(other.Name, &coord, rtt)

		// Publish some metrics to give us an idea of how much we are
		// adjusting each time we update.
		d := float32(before.DistanceTo(after).Seconds() * 1.0e3)
		metrics.AddSample([]string***REMOVED***"serf", "coordinate", "adjustment-ms"***REMOVED***, d)

		// Cache the coordinate for the other node, and add our own
		// to the cache as well since it just got updated. This lets
		// users call GetCachedCoordinate with our node name, which is
		// more friendly.
		p.serf.coordCacheLock.Lock()
		p.serf.coordCache[other.Name] = &coord
		p.serf.coordCache[p.serf.config.NodeName] = p.serf.coordClient.GetCoordinate()
		p.serf.coordCacheLock.Unlock()
	***REMOVED*** else ***REMOVED***
		log.Printf("[ERR] serf: Rejected bad coordinate: %v\n", coord)
	***REMOVED***
***REMOVED***
