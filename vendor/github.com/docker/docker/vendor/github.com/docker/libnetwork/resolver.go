package libnetwork

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/docker/libnetwork/types"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

// Resolver represents the embedded DNS server in Docker. It operates
// by listening on container's loopback interface for DNS queries.
type Resolver interface ***REMOVED***
	// Start starts the name server for the container
	Start() error
	// Stop stops the name server for the container. Stopped resolver
	// can be reused after running the SetupFunc again.
	Stop()
	// SetupFunc() provides the setup function that should be run
	// in the container's network namespace.
	SetupFunc(int) func()
	// NameServer() returns the IP of the DNS resolver for the
	// containers.
	NameServer() string
	// SetExtServers configures the external nameservers the resolver
	// should use to forward queries
	SetExtServers([]extDNSEntry)
	// ResolverOptions returns resolv.conf options that should be set
	ResolverOptions() []string
***REMOVED***

// DNSBackend represents a backend DNS resolver used for DNS name
// resolution. All the queries to the resolver are forwared to the
// backend resolver.
type DNSBackend interface ***REMOVED***
	// ResolveName resolves a service name to an IPv4 or IPv6 address by searching
	// the networks the sandbox is connected to. For IPv6 queries, second return
	// value will be true if the name exists in docker domain but doesn't have an
	// IPv6 address. Such queries shouldn't be forwarded to external nameservers.
	ResolveName(name string, iplen int) ([]net.IP, bool)
	// ResolveIP returns the service name for the passed in IP. IP is in reverse dotted
	// notation; the format used for DNS PTR records
	ResolveIP(name string) string
	// ResolveService returns all the backend details about the containers or hosts
	// backing a service. Its purpose is to satisfy an SRV query
	ResolveService(name string) ([]*net.SRV, []net.IP)
	// ExecFunc allows a function to be executed in the context of the backend
	// on behalf of the resolver.
	ExecFunc(f func()) error
	//NdotsSet queries the backends ndots dns option settings
	NdotsSet() bool
	// HandleQueryResp passes the name & IP from a response to the backend. backend
	// can use it to maintain any required state about the resolution
	HandleQueryResp(name string, ip net.IP)
***REMOVED***

const (
	dnsPort         = "53"
	ptrIPv4domain   = ".in-addr.arpa."
	ptrIPv6domain   = ".ip6.arpa."
	respTTL         = 600
	maxExtDNS       = 3 //max number of external servers to try
	extIOTimeout    = 4 * time.Second
	defaultRespSize = 512
	maxConcurrent   = 100
	logInterval     = 2 * time.Second
)

type extDNSEntry struct ***REMOVED***
	IPStr        string
	HostLoopback bool
***REMOVED***

// resolver implements the Resolver interface
type resolver struct ***REMOVED***
	backend       DNSBackend
	extDNSList    [maxExtDNS]extDNSEntry
	server        *dns.Server
	conn          *net.UDPConn
	tcpServer     *dns.Server
	tcpListen     *net.TCPListener
	err           error
	count         int32
	tStamp        time.Time
	queryLock     sync.Mutex
	listenAddress string
	proxyDNS      bool
	resolverKey   string
	startCh       chan struct***REMOVED******REMOVED***
***REMOVED***

func init() ***REMOVED***
	rand.Seed(time.Now().Unix())
***REMOVED***

// NewResolver creates a new instance of the Resolver
func NewResolver(address string, proxyDNS bool, resolverKey string, backend DNSBackend) Resolver ***REMOVED***
	return &resolver***REMOVED***
		backend:       backend,
		proxyDNS:      proxyDNS,
		listenAddress: address,
		resolverKey:   resolverKey,
		err:           fmt.Errorf("setup not done yet"),
		startCh:       make(chan struct***REMOVED******REMOVED***, 1),
	***REMOVED***
***REMOVED***

func (r *resolver) SetupFunc(port int) func() ***REMOVED***
	return (func() ***REMOVED***
		var err error

		// DNS operates primarily on UDP
		addr := &net.UDPAddr***REMOVED***
			IP:   net.ParseIP(r.listenAddress),
			Port: port,
		***REMOVED***

		r.conn, err = net.ListenUDP("udp", addr)
		if err != nil ***REMOVED***
			r.err = fmt.Errorf("error in opening name server socket %v", err)
			return
		***REMOVED***

		// Listen on a TCP as well
		tcpaddr := &net.TCPAddr***REMOVED***
			IP:   net.ParseIP(r.listenAddress),
			Port: port,
		***REMOVED***

		r.tcpListen, err = net.ListenTCP("tcp", tcpaddr)
		if err != nil ***REMOVED***
			r.err = fmt.Errorf("error in opening name TCP server socket %v", err)
			return
		***REMOVED***
		r.err = nil
	***REMOVED***)
***REMOVED***

func (r *resolver) Start() error ***REMOVED***
	r.startCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	defer func() ***REMOVED*** <-r.startCh ***REMOVED***()

	// make sure the resolver has been setup before starting
	if r.err != nil ***REMOVED***
		return r.err
	***REMOVED***

	if err := r.setupIPTable(); err != nil ***REMOVED***
		return fmt.Errorf("setting up IP table rules failed: %v", err)
	***REMOVED***

	s := &dns.Server***REMOVED***Handler: r, PacketConn: r.conn***REMOVED***
	r.server = s
	go func() ***REMOVED***
		s.ActivateAndServe()
	***REMOVED***()

	tcpServer := &dns.Server***REMOVED***Handler: r, Listener: r.tcpListen***REMOVED***
	r.tcpServer = tcpServer
	go func() ***REMOVED***
		tcpServer.ActivateAndServe()
	***REMOVED***()
	return nil
***REMOVED***

func (r *resolver) Stop() ***REMOVED***
	r.startCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	defer func() ***REMOVED*** <-r.startCh ***REMOVED***()

	if r.server != nil ***REMOVED***
		r.server.Shutdown()
	***REMOVED***
	if r.tcpServer != nil ***REMOVED***
		r.tcpServer.Shutdown()
	***REMOVED***
	r.conn = nil
	r.tcpServer = nil
	r.err = fmt.Errorf("setup not done yet")
	r.tStamp = time.Time***REMOVED******REMOVED***
	r.count = 0
	r.queryLock = sync.Mutex***REMOVED******REMOVED***
***REMOVED***

func (r *resolver) SetExtServers(extDNS []extDNSEntry) ***REMOVED***
	l := len(extDNS)
	if l > maxExtDNS ***REMOVED***
		l = maxExtDNS
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		r.extDNSList[i] = extDNS[i]
	***REMOVED***
***REMOVED***

func (r *resolver) NameServer() string ***REMOVED***
	return r.listenAddress
***REMOVED***

func (r *resolver) ResolverOptions() []string ***REMOVED***
	return []string***REMOVED***"ndots:0"***REMOVED***
***REMOVED***

func setCommonFlags(msg *dns.Msg) ***REMOVED***
	msg.RecursionAvailable = true
***REMOVED***

func shuffleAddr(addr []net.IP) []net.IP ***REMOVED***
	for i := len(addr) - 1; i > 0; i-- ***REMOVED***
		r := rand.Intn(i + 1)
		addr[i], addr[r] = addr[r], addr[i]
	***REMOVED***
	return addr
***REMOVED***

func createRespMsg(query *dns.Msg) *dns.Msg ***REMOVED***
	resp := new(dns.Msg)
	resp.SetReply(query)
	setCommonFlags(resp)

	return resp
***REMOVED***

func (r *resolver) handleMXQuery(name string, query *dns.Msg) (*dns.Msg, error) ***REMOVED***
	addrv4, _ := r.backend.ResolveName(name, types.IPv4)
	addrv6, _ := r.backend.ResolveName(name, types.IPv6)

	if addrv4 == nil && addrv6 == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	// We were able to resolve the name. Respond with an empty list with
	// RcodeSuccess/NOERROR so that email clients can treat it as "implicit MX"
	// [RFC 5321 Section-5.1] and issue a Type A/AAAA query for the name.

	resp := createRespMsg(query)
	return resp, nil
***REMOVED***

func (r *resolver) handleIPQuery(name string, query *dns.Msg, ipType int) (*dns.Msg, error) ***REMOVED***
	var addr []net.IP
	var ipv6Miss bool
	addr, ipv6Miss = r.backend.ResolveName(name, ipType)

	if addr == nil && ipv6Miss ***REMOVED***
		// Send a reply without any Answer sections
		logrus.Debugf("[resolver] lookup name %s present without IPv6 address", name)
		resp := createRespMsg(query)
		return resp, nil
	***REMOVED***
	if addr == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	logrus.Debugf("[resolver] lookup for %s: IP %v", name, addr)

	resp := createRespMsg(query)
	if len(addr) > 1 ***REMOVED***
		addr = shuffleAddr(addr)
	***REMOVED***
	if ipType == types.IPv4 ***REMOVED***
		for _, ip := range addr ***REMOVED***
			rr := new(dns.A)
			rr.Hdr = dns.RR_Header***REMOVED***Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: respTTL***REMOVED***
			rr.A = ip
			resp.Answer = append(resp.Answer, rr)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, ip := range addr ***REMOVED***
			rr := new(dns.AAAA)
			rr.Hdr = dns.RR_Header***REMOVED***Name: name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: respTTL***REMOVED***
			rr.AAAA = ip
			resp.Answer = append(resp.Answer, rr)
		***REMOVED***
	***REMOVED***
	return resp, nil
***REMOVED***

func (r *resolver) handlePTRQuery(ptr string, query *dns.Msg) (*dns.Msg, error) ***REMOVED***
	parts := []string***REMOVED******REMOVED***

	if strings.HasSuffix(ptr, ptrIPv4domain) ***REMOVED***
		parts = strings.Split(ptr, ptrIPv4domain)
	***REMOVED*** else if strings.HasSuffix(ptr, ptrIPv6domain) ***REMOVED***
		parts = strings.Split(ptr, ptrIPv6domain)
	***REMOVED*** else ***REMOVED***
		return nil, fmt.Errorf("invalid PTR query, %v", ptr)
	***REMOVED***

	host := r.backend.ResolveIP(parts[0])

	if len(host) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	logrus.Debugf("[resolver] lookup for IP %s: name %s", parts[0], host)
	fqdn := dns.Fqdn(host)

	resp := new(dns.Msg)
	resp.SetReply(query)
	setCommonFlags(resp)

	rr := new(dns.PTR)
	rr.Hdr = dns.RR_Header***REMOVED***Name: ptr, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: respTTL***REMOVED***
	rr.Ptr = fqdn
	resp.Answer = append(resp.Answer, rr)
	return resp, nil
***REMOVED***

func (r *resolver) handleSRVQuery(svc string, query *dns.Msg) (*dns.Msg, error) ***REMOVED***

	srv, ip := r.backend.ResolveService(svc)

	if len(srv) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	if len(srv) != len(ip) ***REMOVED***
		return nil, fmt.Errorf("invalid reply for SRV query %s", svc)
	***REMOVED***

	resp := createRespMsg(query)

	for i, r := range srv ***REMOVED***
		rr := new(dns.SRV)
		rr.Hdr = dns.RR_Header***REMOVED***Name: svc, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: respTTL***REMOVED***
		rr.Port = r.Port
		rr.Target = r.Target
		resp.Answer = append(resp.Answer, rr)

		rr1 := new(dns.A)
		rr1.Hdr = dns.RR_Header***REMOVED***Name: r.Target, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: respTTL***REMOVED***
		rr1.A = ip[i]
		resp.Extra = append(resp.Extra, rr1)
	***REMOVED***
	return resp, nil

***REMOVED***

func truncateResp(resp *dns.Msg, maxSize int, isTCP bool) ***REMOVED***
	if !isTCP ***REMOVED***
		resp.Truncated = true
	***REMOVED***

	srv := resp.Question[0].Qtype == dns.TypeSRV
	// trim the Answer RRs one by one till the whole message fits
	// within the reply size
	for resp.Len() > maxSize ***REMOVED***
		resp.Answer = resp.Answer[:len(resp.Answer)-1]

		if srv && len(resp.Extra) > 0 ***REMOVED***
			resp.Extra = resp.Extra[:len(resp.Extra)-1]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *resolver) ServeDNS(w dns.ResponseWriter, query *dns.Msg) ***REMOVED***
	var (
		extConn net.Conn
		resp    *dns.Msg
		err     error
	)

	if query == nil || len(query.Question) == 0 ***REMOVED***
		return
	***REMOVED***
	name := query.Question[0].Name

	switch query.Question[0].Qtype ***REMOVED***
	case dns.TypeA:
		resp, err = r.handleIPQuery(name, query, types.IPv4)
	case dns.TypeAAAA:
		resp, err = r.handleIPQuery(name, query, types.IPv6)
	case dns.TypeMX:
		resp, err = r.handleMXQuery(name, query)
	case dns.TypePTR:
		resp, err = r.handlePTRQuery(name, query)
	case dns.TypeSRV:
		resp, err = r.handleSRVQuery(name, query)
	***REMOVED***

	if err != nil ***REMOVED***
		logrus.Error(err)
		return
	***REMOVED***

	if resp == nil ***REMOVED***
		// If the backend doesn't support proxying dns request
		// fail the response
		if !r.proxyDNS ***REMOVED***
			resp = new(dns.Msg)
			resp.SetRcode(query, dns.RcodeServerFailure)
			w.WriteMsg(resp)
			return
		***REMOVED***

		// If the user sets ndots > 0 explicitly and the query is
		// in the root domain don't forward it out. We will return
		// failure and let the client retry with the search domain
		// attached
		switch query.Question[0].Qtype ***REMOVED***
		case dns.TypeA:
			fallthrough
		case dns.TypeAAAA:
			if r.backend.NdotsSet() && !strings.Contains(strings.TrimSuffix(name, "."), ".") ***REMOVED***
				resp = createRespMsg(query)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	proto := w.LocalAddr().Network()
	maxSize := 0
	if proto == "tcp" ***REMOVED***
		maxSize = dns.MaxMsgSize - 1
	***REMOVED*** else if proto == "udp" ***REMOVED***
		optRR := query.IsEdns0()
		if optRR != nil ***REMOVED***
			maxSize = int(optRR.UDPSize())
		***REMOVED***
		if maxSize < defaultRespSize ***REMOVED***
			maxSize = defaultRespSize
		***REMOVED***
	***REMOVED***

	if resp != nil ***REMOVED***
		if resp.Len() > maxSize ***REMOVED***
			truncateResp(resp, maxSize, proto == "tcp")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := 0; i < maxExtDNS; i++ ***REMOVED***
			extDNS := &r.extDNSList[i]
			if extDNS.IPStr == "" ***REMOVED***
				break
			***REMOVED***
			extConnect := func() ***REMOVED***
				addr := fmt.Sprintf("%s:%d", extDNS.IPStr, 53)
				extConn, err = net.DialTimeout(proto, addr, extIOTimeout)
			***REMOVED***

			if extDNS.HostLoopback ***REMOVED***
				extConnect()
			***REMOVED*** else ***REMOVED***
				execErr := r.backend.ExecFunc(extConnect)
				if execErr != nil ***REMOVED***
					logrus.Warn(execErr)
					continue
				***REMOVED***
			***REMOVED***
			if err != nil ***REMOVED***
				logrus.Warnf("[resolver] connect failed: %s", err)
				continue
			***REMOVED***

			queryType := dns.TypeToString[query.Question[0].Qtype]
			logrus.Debugf("[resolver] query %s (%s) from %s, forwarding to %s:%s", name, queryType,
				extConn.LocalAddr().String(), proto, extDNS.IPStr)

			// Timeout has to be set for every IO operation.
			extConn.SetDeadline(time.Now().Add(extIOTimeout))
			co := &dns.Conn***REMOVED***
				Conn:    extConn,
				UDPSize: uint16(maxSize),
			***REMOVED***
			defer co.Close()

			// limits the number of outstanding concurrent queries.
			if !r.forwardQueryStart() ***REMOVED***
				old := r.tStamp
				r.tStamp = time.Now()
				if r.tStamp.Sub(old) > logInterval ***REMOVED***
					logrus.Errorf("[resolver] more than %v concurrent queries from %s", maxConcurrent, extConn.LocalAddr().String())
				***REMOVED***
				continue
			***REMOVED***

			err = co.WriteMsg(query)
			if err != nil ***REMOVED***
				r.forwardQueryEnd()
				logrus.Debugf("[resolver] send to DNS server failed, %s", err)
				continue
			***REMOVED***

			resp, err = co.ReadMsg()
			// Truncated DNS replies should be sent to the client so that the
			// client can retry over TCP
			if err != nil && err != dns.ErrTruncated ***REMOVED***
				r.forwardQueryEnd()
				logrus.Debugf("[resolver] read from DNS server failed, %s", err)
				continue
			***REMOVED***
			r.forwardQueryEnd()
			if resp != nil ***REMOVED***
				answers := 0
				for _, rr := range resp.Answer ***REMOVED***
					h := rr.Header()
					switch h.Rrtype ***REMOVED***
					case dns.TypeA:
						answers++
						ip := rr.(*dns.A).A
						logrus.Debugf("[resolver] received A record %q for %q from %s:%s", ip, h.Name, proto, extDNS.IPStr)
						r.backend.HandleQueryResp(h.Name, ip)
					case dns.TypeAAAA:
						answers++
						ip := rr.(*dns.AAAA).AAAA
						logrus.Debugf("[resolver] received AAAA record %q for %q from %s:%s", ip, h.Name, proto, extDNS.IPStr)
						r.backend.HandleQueryResp(h.Name, ip)
					***REMOVED***
				***REMOVED***
				if resp.Answer == nil || answers == 0 ***REMOVED***
					logrus.Debugf("[resolver] external DNS %s:%s did not return any %s records for %q", proto, extDNS.IPStr, queryType, name)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				logrus.Debugf("[resolver] external DNS %s:%s returned empty response for %q", proto, extDNS.IPStr, name)
			***REMOVED***
			resp.Compress = true
			break
		***REMOVED***
		if resp == nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if err = w.WriteMsg(resp); err != nil ***REMOVED***
		logrus.Errorf("[resolver] error writing resolver resp, %s", err)
	***REMOVED***
***REMOVED***

func (r *resolver) forwardQueryStart() bool ***REMOVED***
	r.queryLock.Lock()
	defer r.queryLock.Unlock()

	if r.count == maxConcurrent ***REMOVED***
		return false
	***REMOVED***
	r.count++

	return true
***REMOVED***

func (r *resolver) forwardQueryEnd() ***REMOVED***
	r.queryLock.Lock()
	defer r.queryLock.Unlock()

	if r.count == 0 ***REMOVED***
		logrus.Error("[resolver] invalid concurrent query count")
	***REMOVED*** else ***REMOVED***
		r.count--
	***REMOVED***
***REMOVED***
