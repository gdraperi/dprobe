package overlay

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
)

type ovNotify struct ***REMOVED***
	action string
	ep     *endpoint
	nw     *network
***REMOVED***

type logWriter struct***REMOVED******REMOVED***

func (l *logWriter) Write(p []byte) (int, error) ***REMOVED***
	str := string(p)

	switch ***REMOVED***
	case strings.Contains(str, "[WARN]"):
		logrus.Warn(str)
	case strings.Contains(str, "[DEBUG]"):
		logrus.Debug(str)
	case strings.Contains(str, "[INFO]"):
		logrus.Info(str)
	case strings.Contains(str, "[ERR]"):
		logrus.Error(str)
	***REMOVED***

	return len(p), nil
***REMOVED***

func (d *driver) serfInit() error ***REMOVED***
	var err error

	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = d.advertiseAddress

	d.eventCh = make(chan serf.Event, 4)
	config.EventCh = d.eventCh
	config.UserCoalescePeriod = 1 * time.Second
	config.UserQuiescentPeriod = 50 * time.Millisecond

	config.LogOutput = &logWriter***REMOVED******REMOVED***
	config.MemberlistConfig.LogOutput = config.LogOutput

	s, err := serf.Create(config)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to create cluster node: %v", err)
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			s.Shutdown()
		***REMOVED***
	***REMOVED***()

	d.serfInstance = s

	d.notifyCh = make(chan ovNotify)
	d.exitCh = make(chan chan struct***REMOVED******REMOVED***)

	go d.startSerfLoop(d.eventCh, d.notifyCh, d.exitCh)
	return nil
***REMOVED***

func (d *driver) serfJoin(neighIP string) error ***REMOVED***
	if neighIP == "" ***REMOVED***
		return fmt.Errorf("no neighbor to join")
	***REMOVED***
	if _, err := d.serfInstance.Join([]string***REMOVED***neighIP***REMOVED***, true); err != nil ***REMOVED***
		return fmt.Errorf("Failed to join the cluster at neigh IP %s: %v",
			neighIP, err)
	***REMOVED***
	return nil
***REMOVED***

func (d *driver) notifyEvent(event ovNotify) ***REMOVED***
	ep := event.ep

	ePayload := fmt.Sprintf("%s %s %s %s", event.action, ep.addr.IP.String(),
		net.IP(ep.addr.Mask).String(), ep.mac.String())
	eName := fmt.Sprintf("jl %s %s %s", d.serfInstance.LocalMember().Addr.String(),
		event.nw.id, ep.id)

	if err := d.serfInstance.UserEvent(eName, []byte(ePayload), true); err != nil ***REMOVED***
		logrus.Errorf("Sending user event failed: %v\n", err)
	***REMOVED***
***REMOVED***

func (d *driver) processEvent(u serf.UserEvent) ***REMOVED***
	logrus.Debugf("Received user event name:%s, payload:%s LTime:%d \n", u.Name,
		string(u.Payload), uint64(u.LTime))

	var dummy, action, vtepStr, nid, eid, ipStr, maskStr, macStr string
	if _, err := fmt.Sscan(u.Name, &dummy, &vtepStr, &nid, &eid); err != nil ***REMOVED***
		fmt.Printf("Failed to scan name string: %v\n", err)
	***REMOVED***

	if _, err := fmt.Sscan(string(u.Payload), &action,
		&ipStr, &maskStr, &macStr); err != nil ***REMOVED***
		fmt.Printf("Failed to scan value string: %v\n", err)
	***REMOVED***

	logrus.Debugf("Parsed data = %s/%s/%s/%s/%s/%s\n", nid, eid, vtepStr, ipStr, maskStr, macStr)

	mac, err := net.ParseMAC(macStr)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to parse mac: %v\n", err)
	***REMOVED***

	if d.serfInstance.LocalMember().Addr.String() == vtepStr ***REMOVED***
		return
	***REMOVED***

	switch action ***REMOVED***
	case "join":
		d.peerAdd(nid, eid, net.ParseIP(ipStr), net.IPMask(net.ParseIP(maskStr).To4()), mac, net.ParseIP(vtepStr), false, false, false)
	case "leave":
		d.peerDelete(nid, eid, net.ParseIP(ipStr), net.IPMask(net.ParseIP(maskStr).To4()), mac, net.ParseIP(vtepStr), false)
	***REMOVED***
***REMOVED***

func (d *driver) processQuery(q *serf.Query) ***REMOVED***
	logrus.Debugf("Received query name:%s, payload:%s\n", q.Name,
		string(q.Payload))

	var nid, ipStr string
	if _, err := fmt.Sscan(string(q.Payload), &nid, &ipStr); err != nil ***REMOVED***
		fmt.Printf("Failed to scan query payload string: %v\n", err)
	***REMOVED***

	pKey, pEntry, err := d.peerDbSearch(nid, net.ParseIP(ipStr))
	if err != nil ***REMOVED***
		return
	***REMOVED***

	logrus.Debugf("Sending peer query resp mac %v, mask %s, vtep %s", pKey.peerMac, net.IP(pEntry.peerIPMask).String(), pEntry.vtep)
	q.Respond([]byte(fmt.Sprintf("%s %s %s", pKey.peerMac.String(), net.IP(pEntry.peerIPMask).String(), pEntry.vtep.String())))
***REMOVED***

func (d *driver) resolvePeer(nid string, peerIP net.IP) (net.HardwareAddr, net.IPMask, net.IP, error) ***REMOVED***
	if d.serfInstance == nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("could not resolve peer: serf instance not initialized")
	***REMOVED***

	qPayload := fmt.Sprintf("%s %s", string(nid), peerIP.String())
	resp, err := d.serfInstance.Query("peerlookup", []byte(qPayload), nil)
	if err != nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("resolving peer by querying the cluster failed: %v", err)
	***REMOVED***

	respCh := resp.ResponseCh()
	select ***REMOVED***
	case r := <-respCh:
		var macStr, maskStr, vtepStr string
		if _, err := fmt.Sscan(string(r.Payload), &macStr, &maskStr, &vtepStr); err != nil ***REMOVED***
			return nil, nil, nil, fmt.Errorf("bad response %q for the resolve query: %v", string(r.Payload), err)
		***REMOVED***

		mac, err := net.ParseMAC(macStr)
		if err != nil ***REMOVED***
			return nil, nil, nil, fmt.Errorf("failed to parse mac: %v", err)
		***REMOVED***

		logrus.Debugf("Received peer query response, mac %s, vtep %s, mask %s", macStr, vtepStr, maskStr)
		return mac, net.IPMask(net.ParseIP(maskStr).To4()), net.ParseIP(vtepStr), nil

	case <-time.After(time.Second):
		return nil, nil, nil, fmt.Errorf("timed out resolving peer by querying the cluster")
	***REMOVED***
***REMOVED***

func (d *driver) startSerfLoop(eventCh chan serf.Event, notifyCh chan ovNotify,
	exitCh chan chan struct***REMOVED******REMOVED***) ***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case notify, ok := <-notifyCh:
			if !ok ***REMOVED***
				break
			***REMOVED***

			d.notifyEvent(notify)
		case ch, ok := <-exitCh:
			if !ok ***REMOVED***
				break
			***REMOVED***

			if err := d.serfInstance.Leave(); err != nil ***REMOVED***
				logrus.Errorf("failed leaving the cluster: %v\n", err)
			***REMOVED***

			d.serfInstance.Shutdown()
			close(ch)
			return
		case e, ok := <-eventCh:
			if !ok ***REMOVED***
				break
			***REMOVED***

			if e.EventType() == serf.EventQuery ***REMOVED***
				d.processQuery(e.(*serf.Query))
				break
			***REMOVED***

			u, ok := e.(serf.UserEvent)
			if !ok ***REMOVED***
				break
			***REMOVED***
			d.processEvent(u)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *driver) isSerfAlive() bool ***REMOVED***
	d.Lock()
	serfInstance := d.serfInstance
	d.Unlock()
	if serfInstance == nil || serfInstance.State() != serf.SerfAlive ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***
