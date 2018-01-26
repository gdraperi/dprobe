package zk

/*
TODO:
* make sure a ping response comes back in a reasonable time

Possible watcher events:
* Event***REMOVED***Type: EventNotWatching, State: StateDisconnected, Path: path, Err: err***REMOVED***
*/

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNoServer = errors.New("zk: could not connect to a server")

const (
	bufferSize      = 1536 * 1024
	eventChanSize   = 6
	sendChanSize    = 16
	protectedPrefix = "_c_"
)

type watchType int

const (
	watchTypeData  = iota
	watchTypeExist = iota
	watchTypeChild = iota
)

type watchPathType struct ***REMOVED***
	path  string
	wType watchType
***REMOVED***

type Dialer func(network, address string, timeout time.Duration) (net.Conn, error)

type Conn struct ***REMOVED***
	lastZxid  int64
	sessionID int64
	state     State // must be 32-bit aligned
	xid       uint32
	timeout   int32 // session timeout in milliseconds
	passwd    []byte

	dialer          Dialer
	servers         []string
	serverIndex     int // remember last server that was tried during connect to round-robin attempts to servers
	lastServerIndex int // index of the last server that was successfully connected to and authenticated with
	conn            net.Conn
	eventChan       chan Event
	shouldQuit      chan struct***REMOVED******REMOVED***
	pingInterval    time.Duration
	recvTimeout     time.Duration
	connectTimeout  time.Duration

	sendChan     chan *request
	requests     map[int32]*request // Xid -> pending request
	requestsLock sync.Mutex
	watchers     map[watchPathType][]chan Event
	watchersLock sync.Mutex

	// Debug (used by unit tests)
	reconnectDelay time.Duration
***REMOVED***

type request struct ***REMOVED***
	xid        int32
	opcode     int32
	pkt        interface***REMOVED******REMOVED***
	recvStruct interface***REMOVED******REMOVED***
	recvChan   chan response

	// Because sending and receiving happen in separate go routines, there's
	// a possible race condition when creating watches from outside the read
	// loop. We must ensure that a watcher gets added to the list synchronously
	// with the response from the server on any request that creates a watch.
	// In order to not hard code the watch logic for each opcode in the recv
	// loop the caller can use recvFunc to insert some synchronously code
	// after a response.
	recvFunc func(*request, *responseHeader, error)
***REMOVED***

type response struct ***REMOVED***
	zxid int64
	err  error
***REMOVED***

type Event struct ***REMOVED***
	Type   EventType
	State  State
	Path   string // For non-session events, the path of the watched node.
	Err    error
	Server string // For connection events
***REMOVED***

// Connect establishes a new connection to a pool of zookeeper servers
// using the default net.Dialer. See ConnectWithDialer for further
// information about session timeout.
func Connect(servers []string, sessionTimeout time.Duration) (*Conn, <-chan Event, error) ***REMOVED***
	return ConnectWithDialer(servers, sessionTimeout, nil)
***REMOVED***

// ConnectWithDialer establishes a new connection to a pool of zookeeper
// servers. The provided session timeout sets the amount of time for which
// a session is considered valid after losing connection to a server. Within
// the session timeout it's possible to reestablish a connection to a different
// server and keep the same session. This is means any ephemeral nodes and
// watches are maintained.
func ConnectWithDialer(servers []string, sessionTimeout time.Duration, dialer Dialer) (*Conn, <-chan Event, error) ***REMOVED***
	if len(servers) == 0 ***REMOVED***
		return nil, nil, errors.New("zk: server list must not be empty")
	***REMOVED***

	recvTimeout := sessionTimeout * 2 / 3

	srvs := make([]string, len(servers))

	for i, addr := range servers ***REMOVED***
		if strings.Contains(addr, ":") ***REMOVED***
			srvs[i] = addr
		***REMOVED*** else ***REMOVED***
			srvs[i] = addr + ":" + strconv.Itoa(DefaultPort)
		***REMOVED***
	***REMOVED***

	// Randomize the order of the servers to avoid creating hotspots
	stringShuffle(srvs)

	ec := make(chan Event, eventChanSize)
	if dialer == nil ***REMOVED***
		dialer = net.DialTimeout
	***REMOVED***
	conn := Conn***REMOVED***
		dialer:          dialer,
		servers:         srvs,
		serverIndex:     0,
		lastServerIndex: -1,
		conn:            nil,
		state:           StateDisconnected,
		eventChan:       ec,
		shouldQuit:      make(chan struct***REMOVED******REMOVED***),
		recvTimeout:     recvTimeout,
		pingInterval:    recvTimeout / 2,
		connectTimeout:  1 * time.Second,
		sendChan:        make(chan *request, sendChanSize),
		requests:        make(map[int32]*request),
		watchers:        make(map[watchPathType][]chan Event),
		passwd:          emptyPassword,
		timeout:         int32(sessionTimeout.Nanoseconds() / 1e6),

		// Debug
		reconnectDelay: 0,
	***REMOVED***
	go func() ***REMOVED***
		conn.loop()
		conn.flushRequests(ErrClosing)
		conn.invalidateWatches(ErrClosing)
		close(conn.eventChan)
	***REMOVED***()
	return &conn, ec, nil
***REMOVED***

func (c *Conn) Close() ***REMOVED***
	close(c.shouldQuit)

	select ***REMOVED***
	case <-c.queueRequest(opClose, &closeRequest***REMOVED******REMOVED***, &closeResponse***REMOVED******REMOVED***, nil):
	case <-time.After(time.Second):
	***REMOVED***
***REMOVED***

func (c *Conn) State() State ***REMOVED***
	return State(atomic.LoadInt32((*int32)(&c.state)))
***REMOVED***

func (c *Conn) setState(state State) ***REMOVED***
	atomic.StoreInt32((*int32)(&c.state), int32(state))
	select ***REMOVED***
	case c.eventChan <- Event***REMOVED***Type: EventSession, State: state, Server: c.servers[c.serverIndex]***REMOVED***:
	default:
		// panic("zk: event channel full - it must be monitored and never allowed to be full")
	***REMOVED***
***REMOVED***

func (c *Conn) connect() error ***REMOVED***
	c.setState(StateConnecting)
	for ***REMOVED***
		c.serverIndex = (c.serverIndex + 1) % len(c.servers)
		if c.serverIndex == c.lastServerIndex ***REMOVED***
			c.flushUnsentRequests(ErrNoServer)
			select ***REMOVED***
			case <-time.After(time.Second):
				// pass
			case <-c.shouldQuit:
				c.setState(StateDisconnected)
				c.flushUnsentRequests(ErrClosing)
				return ErrClosing
			***REMOVED***
		***REMOVED*** else if c.lastServerIndex < 0 ***REMOVED***
			// lastServerIndex defaults to -1 to avoid a delay on the initial connect
			c.lastServerIndex = 0
		***REMOVED***

		zkConn, err := c.dialer("tcp", c.servers[c.serverIndex], c.connectTimeout)
		if err == nil ***REMOVED***
			c.conn = zkConn
			c.setState(StateConnected)
			return nil
		***REMOVED***

		log.Printf("Failed to connect to %s: %+v", c.servers[c.serverIndex], err)
	***REMOVED***
***REMOVED***

func (c *Conn) loop() ***REMOVED***
	for ***REMOVED***
		if err := c.connect(); err != nil ***REMOVED***
			// c.Close() was called
			return
		***REMOVED***

		err := c.authenticate()
		switch ***REMOVED***
		case err == ErrSessionExpired:
			c.invalidateWatches(err)
		case err != nil && c.conn != nil:
			c.conn.Close()
		case err == nil:
			c.lastServerIndex = c.serverIndex
			closeChan := make(chan struct***REMOVED******REMOVED***) // channel to tell send loop stop
			var wg sync.WaitGroup

			wg.Add(1)
			go func() ***REMOVED***
				c.sendLoop(c.conn, closeChan)
				c.conn.Close() // causes recv loop to EOF/exit
				wg.Done()
			***REMOVED***()

			wg.Add(1)
			go func() ***REMOVED***
				err = c.recvLoop(c.conn)
				if err == nil ***REMOVED***
					panic("zk: recvLoop should never return nil error")
				***REMOVED***
				close(closeChan) // tell send loop to exit
				wg.Done()
			***REMOVED***()

			wg.Wait()
		***REMOVED***

		c.setState(StateDisconnected)

		// Yeesh
		if err != io.EOF && err != ErrSessionExpired && !strings.Contains(err.Error(), "use of closed network connection") ***REMOVED***
			log.Println(err)
		***REMOVED***

		select ***REMOVED***
		case <-c.shouldQuit:
			c.flushRequests(ErrClosing)
			return
		default:
		***REMOVED***

		if err != ErrSessionExpired ***REMOVED***
			err = ErrConnectionClosed
		***REMOVED***
		c.flushRequests(err)

		if c.reconnectDelay > 0 ***REMOVED***
			select ***REMOVED***
			case <-c.shouldQuit:
				return
			case <-time.After(c.reconnectDelay):
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Conn) flushUnsentRequests(err error) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		default:
			return
		case req := <-c.sendChan:
			req.recvChan <- response***REMOVED***-1, err***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Send error to all pending requests and clear request map
func (c *Conn) flushRequests(err error) ***REMOVED***
	c.requestsLock.Lock()
	for _, req := range c.requests ***REMOVED***
		req.recvChan <- response***REMOVED***-1, err***REMOVED***
	***REMOVED***
	c.requests = make(map[int32]*request)
	c.requestsLock.Unlock()
***REMOVED***

// Send error to all watchers and clear watchers map
func (c *Conn) invalidateWatches(err error) ***REMOVED***
	c.watchersLock.Lock()
	defer c.watchersLock.Unlock()

	if len(c.watchers) >= 0 ***REMOVED***
		for pathType, watchers := range c.watchers ***REMOVED***
			ev := Event***REMOVED***Type: EventNotWatching, State: StateDisconnected, Path: pathType.path, Err: err***REMOVED***
			for _, ch := range watchers ***REMOVED***
				ch <- ev
				close(ch)
			***REMOVED***
		***REMOVED***
		c.watchers = make(map[watchPathType][]chan Event)
	***REMOVED***
***REMOVED***

func (c *Conn) sendSetWatches() ***REMOVED***
	c.watchersLock.Lock()
	defer c.watchersLock.Unlock()

	if len(c.watchers) == 0 ***REMOVED***
		return
	***REMOVED***

	req := &setWatchesRequest***REMOVED***
		RelativeZxid: c.lastZxid,
		DataWatches:  make([]string, 0),
		ExistWatches: make([]string, 0),
		ChildWatches: make([]string, 0),
	***REMOVED***
	n := 0
	for pathType, watchers := range c.watchers ***REMOVED***
		if len(watchers) == 0 ***REMOVED***
			continue
		***REMOVED***
		switch pathType.wType ***REMOVED***
		case watchTypeData:
			req.DataWatches = append(req.DataWatches, pathType.path)
		case watchTypeExist:
			req.ExistWatches = append(req.ExistWatches, pathType.path)
		case watchTypeChild:
			req.ChildWatches = append(req.ChildWatches, pathType.path)
		***REMOVED***
		n++
	***REMOVED***
	if n == 0 ***REMOVED***
		return
	***REMOVED***

	go func() ***REMOVED***
		res := &setWatchesResponse***REMOVED******REMOVED***
		_, err := c.request(opSetWatches, req, res, nil)
		if err != nil ***REMOVED***
			log.Printf("Failed to set previous watches: %s", err.Error())
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (c *Conn) authenticate() error ***REMOVED***
	buf := make([]byte, 256)

	// connect request

	n, err := encodePacket(buf[4:], &connectRequest***REMOVED***
		ProtocolVersion: protocolVersion,
		LastZxidSeen:    c.lastZxid,
		TimeOut:         c.timeout,
		SessionID:       c.sessionID,
		Passwd:          c.passwd,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	binary.BigEndian.PutUint32(buf[:4], uint32(n))

	c.conn.SetWriteDeadline(time.Now().Add(c.recvTimeout * 10))
	_, err = c.conn.Write(buf[:n+4])
	c.conn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.sendSetWatches()

	// connect response

	// package length
	c.conn.SetReadDeadline(time.Now().Add(c.recvTimeout * 10))
	_, err = io.ReadFull(c.conn, buf[:4])
	c.conn.SetReadDeadline(time.Time***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		// Sometimes zookeeper just drops connection on invalid session data,
		// we prefer to drop session and start from scratch when that event
		// occurs instead of dropping into loop of connect/disconnect attempts
		c.sessionID = 0
		c.passwd = emptyPassword
		c.lastZxid = 0
		c.setState(StateExpired)
		return ErrSessionExpired
	***REMOVED***

	blen := int(binary.BigEndian.Uint32(buf[:4]))
	if cap(buf) < blen ***REMOVED***
		buf = make([]byte, blen)
	***REMOVED***

	_, err = io.ReadFull(c.conn, buf[:blen])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r := connectResponse***REMOVED******REMOVED***
	_, err = decodePacket(buf[:blen], &r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.SessionID == 0 ***REMOVED***
		c.sessionID = 0
		c.passwd = emptyPassword
		c.lastZxid = 0
		c.setState(StateExpired)
		return ErrSessionExpired
	***REMOVED***

	if c.sessionID != r.SessionID ***REMOVED***
		atomic.StoreUint32(&c.xid, 0)
	***REMOVED***
	c.timeout = r.TimeOut
	c.sessionID = r.SessionID
	c.passwd = r.Passwd
	c.setState(StateHasSession)

	return nil
***REMOVED***

func (c *Conn) sendLoop(conn net.Conn, closeChan <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	pingTicker := time.NewTicker(c.pingInterval)
	defer pingTicker.Stop()

	buf := make([]byte, bufferSize)
	for ***REMOVED***
		select ***REMOVED***
		case req := <-c.sendChan:
			header := &requestHeader***REMOVED***req.xid, req.opcode***REMOVED***
			n, err := encodePacket(buf[4:], header)
			if err != nil ***REMOVED***
				req.recvChan <- response***REMOVED***-1, err***REMOVED***
				continue
			***REMOVED***

			n2, err := encodePacket(buf[4+n:], req.pkt)
			if err != nil ***REMOVED***
				req.recvChan <- response***REMOVED***-1, err***REMOVED***
				continue
			***REMOVED***

			n += n2

			binary.BigEndian.PutUint32(buf[:4], uint32(n))

			c.requestsLock.Lock()
			select ***REMOVED***
			case <-closeChan:
				req.recvChan <- response***REMOVED***-1, ErrConnectionClosed***REMOVED***
				c.requestsLock.Unlock()
				return ErrConnectionClosed
			default:
			***REMOVED***
			c.requests[req.xid] = req
			c.requestsLock.Unlock()

			conn.SetWriteDeadline(time.Now().Add(c.recvTimeout))
			_, err = conn.Write(buf[:n+4])
			conn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				req.recvChan <- response***REMOVED***-1, err***REMOVED***
				conn.Close()
				return err
			***REMOVED***
		case <-pingTicker.C:
			n, err := encodePacket(buf[4:], &requestHeader***REMOVED***Xid: -2, Opcode: opPing***REMOVED***)
			if err != nil ***REMOVED***
				panic("zk: opPing should never fail to serialize")
			***REMOVED***

			binary.BigEndian.PutUint32(buf[:4], uint32(n))

			conn.SetWriteDeadline(time.Now().Add(c.recvTimeout))
			_, err = conn.Write(buf[:n+4])
			conn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				conn.Close()
				return err
			***REMOVED***
		case <-closeChan:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Conn) recvLoop(conn net.Conn) error ***REMOVED***
	buf := make([]byte, bufferSize)
	for ***REMOVED***
		// package length
		conn.SetReadDeadline(time.Now().Add(c.recvTimeout))
		_, err := io.ReadFull(conn, buf[:4])
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		blen := int(binary.BigEndian.Uint32(buf[:4]))
		if cap(buf) < blen ***REMOVED***
			buf = make([]byte, blen)
		***REMOVED***

		_, err = io.ReadFull(conn, buf[:blen])
		conn.SetReadDeadline(time.Time***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		res := responseHeader***REMOVED******REMOVED***
		_, err = decodePacket(buf[:16], &res)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if res.Xid == -1 ***REMOVED***
			res := &watcherEvent***REMOVED******REMOVED***
			_, err := decodePacket(buf[16:16+blen], res)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			ev := Event***REMOVED***
				Type:  res.Type,
				State: res.State,
				Path:  res.Path,
				Err:   nil,
			***REMOVED***
			select ***REMOVED***
			case c.eventChan <- ev:
			default:
			***REMOVED***
			wTypes := make([]watchType, 0, 2)
			switch res.Type ***REMOVED***
			case EventNodeCreated:
				wTypes = append(wTypes, watchTypeExist)
			case EventNodeDeleted, EventNodeDataChanged:
				wTypes = append(wTypes, watchTypeExist, watchTypeData, watchTypeChild)
			case EventNodeChildrenChanged:
				wTypes = append(wTypes, watchTypeChild)
			***REMOVED***
			c.watchersLock.Lock()
			for _, t := range wTypes ***REMOVED***
				wpt := watchPathType***REMOVED***res.Path, t***REMOVED***
				if watchers := c.watchers[wpt]; watchers != nil && len(watchers) > 0 ***REMOVED***
					for _, ch := range watchers ***REMOVED***
						ch <- ev
						close(ch)
					***REMOVED***
					delete(c.watchers, wpt)
				***REMOVED***
			***REMOVED***
			c.watchersLock.Unlock()
		***REMOVED*** else if res.Xid == -2 ***REMOVED***
			// Ping response. Ignore.
		***REMOVED*** else if res.Xid < 0 ***REMOVED***
			log.Printf("Xid < 0 (%d) but not ping or watcher event", res.Xid)
		***REMOVED*** else ***REMOVED***
			if res.Zxid > 0 ***REMOVED***
				c.lastZxid = res.Zxid
			***REMOVED***

			c.requestsLock.Lock()
			req, ok := c.requests[res.Xid]
			if ok ***REMOVED***
				delete(c.requests, res.Xid)
			***REMOVED***
			c.requestsLock.Unlock()

			if !ok ***REMOVED***
				log.Printf("Response for unknown request with xid %d", res.Xid)
			***REMOVED*** else ***REMOVED***
				if res.Err != 0 ***REMOVED***
					err = res.Err.toError()
				***REMOVED*** else ***REMOVED***
					_, err = decodePacket(buf[16:16+blen], req.recvStruct)
				***REMOVED***
				if req.recvFunc != nil ***REMOVED***
					req.recvFunc(req, &res, err)
				***REMOVED***
				req.recvChan <- response***REMOVED***res.Zxid, err***REMOVED***
				if req.opcode == opClose ***REMOVED***
					return io.EOF
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Conn) nextXid() int32 ***REMOVED***
	return int32(atomic.AddUint32(&c.xid, 1) & 0x7fffffff)
***REMOVED***

func (c *Conn) addWatcher(path string, watchType watchType) <-chan Event ***REMOVED***
	c.watchersLock.Lock()
	defer c.watchersLock.Unlock()

	ch := make(chan Event, 1)
	wpt := watchPathType***REMOVED***path, watchType***REMOVED***
	c.watchers[wpt] = append(c.watchers[wpt], ch)
	return ch
***REMOVED***

func (c *Conn) queueRequest(opcode int32, req interface***REMOVED******REMOVED***, res interface***REMOVED******REMOVED***, recvFunc func(*request, *responseHeader, error)) <-chan response ***REMOVED***
	rq := &request***REMOVED***
		xid:        c.nextXid(),
		opcode:     opcode,
		pkt:        req,
		recvStruct: res,
		recvChan:   make(chan response, 1),
		recvFunc:   recvFunc,
	***REMOVED***
	c.sendChan <- rq
	return rq.recvChan
***REMOVED***

func (c *Conn) request(opcode int32, req interface***REMOVED******REMOVED***, res interface***REMOVED******REMOVED***, recvFunc func(*request, *responseHeader, error)) (int64, error) ***REMOVED***
	r := <-c.queueRequest(opcode, req, res, recvFunc)
	return r.zxid, r.err
***REMOVED***

func (c *Conn) AddAuth(scheme string, auth []byte) error ***REMOVED***
	_, err := c.request(opSetAuth, &setAuthRequest***REMOVED***Type: 0, Scheme: scheme, Auth: auth***REMOVED***, &setAuthResponse***REMOVED******REMOVED***, nil)
	return err
***REMOVED***

func (c *Conn) Children(path string) ([]string, *Stat, error) ***REMOVED***
	res := &getChildren2Response***REMOVED******REMOVED***
	_, err := c.request(opGetChildren2, &getChildren2Request***REMOVED***Path: path, Watch: false***REMOVED***, res, nil)
	return res.Children, &res.Stat, err
***REMOVED***

func (c *Conn) ChildrenW(path string) ([]string, *Stat, <-chan Event, error) ***REMOVED***
	var ech <-chan Event
	res := &getChildren2Response***REMOVED******REMOVED***
	_, err := c.request(opGetChildren2, &getChildren2Request***REMOVED***Path: path, Watch: true***REMOVED***, res, func(req *request, res *responseHeader, err error) ***REMOVED***
		if err == nil ***REMOVED***
			ech = c.addWatcher(path, watchTypeChild)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***
	return res.Children, &res.Stat, ech, err
***REMOVED***

func (c *Conn) Get(path string) ([]byte, *Stat, error) ***REMOVED***
	res := &getDataResponse***REMOVED******REMOVED***
	_, err := c.request(opGetData, &getDataRequest***REMOVED***Path: path, Watch: false***REMOVED***, res, nil)
	return res.Data, &res.Stat, err
***REMOVED***

// GetW returns the contents of a znode and sets a watch
func (c *Conn) GetW(path string) ([]byte, *Stat, <-chan Event, error) ***REMOVED***
	var ech <-chan Event
	res := &getDataResponse***REMOVED******REMOVED***
	_, err := c.request(opGetData, &getDataRequest***REMOVED***Path: path, Watch: true***REMOVED***, res, func(req *request, res *responseHeader, err error) ***REMOVED***
		if err == nil ***REMOVED***
			ech = c.addWatcher(path, watchTypeData)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***
	return res.Data, &res.Stat, ech, err
***REMOVED***

func (c *Conn) Set(path string, data []byte, version int32) (*Stat, error) ***REMOVED***
	res := &setDataResponse***REMOVED******REMOVED***
	_, err := c.request(opSetData, &SetDataRequest***REMOVED***path, data, version***REMOVED***, res, nil)
	return &res.Stat, err
***REMOVED***

func (c *Conn) Create(path string, data []byte, flags int32, acl []ACL) (string, error) ***REMOVED***
	res := &createResponse***REMOVED******REMOVED***
	_, err := c.request(opCreate, &CreateRequest***REMOVED***path, data, acl, flags***REMOVED***, res, nil)
	return res.Path, err
***REMOVED***

// CreateProtectedEphemeralSequential fixes a race condition if the server crashes
// after it creates the node. On reconnect the session may still be valid so the
// ephemeral node still exists. Therefore, on reconnect we need to check if a node
// with a GUID generated on create exists.
func (c *Conn) CreateProtectedEphemeralSequential(path string, data []byte, acl []ACL) (string, error) ***REMOVED***
	var guid [16]byte
	_, err := io.ReadFull(rand.Reader, guid[:16])
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	guidStr := fmt.Sprintf("%x", guid)

	parts := strings.Split(path, "/")
	parts[len(parts)-1] = fmt.Sprintf("%s%s-%s", protectedPrefix, guidStr, parts[len(parts)-1])
	rootPath := strings.Join(parts[:len(parts)-1], "/")
	protectedPath := strings.Join(parts, "/")

	var newPath string
	for i := 0; i < 3; i++ ***REMOVED***
		newPath, err = c.Create(protectedPath, data, FlagEphemeral|FlagSequence, acl)
		switch err ***REMOVED***
		case ErrSessionExpired:
			// No need to search for the node since it can't exist. Just try again.
		case ErrConnectionClosed:
			children, _, err := c.Children(rootPath)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			for _, p := range children ***REMOVED***
				parts := strings.Split(p, "/")
				if pth := parts[len(parts)-1]; strings.HasPrefix(pth, protectedPrefix) ***REMOVED***
					if g := pth[len(protectedPrefix) : len(protectedPrefix)+32]; g == guidStr ***REMOVED***
						return rootPath + "/" + p, nil
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case nil:
			return newPath, nil
		default:
			return "", err
		***REMOVED***
	***REMOVED***
	return "", err
***REMOVED***

func (c *Conn) Delete(path string, version int32) error ***REMOVED***
	_, err := c.request(opDelete, &DeleteRequest***REMOVED***path, version***REMOVED***, &deleteResponse***REMOVED******REMOVED***, nil)
	return err
***REMOVED***

func (c *Conn) Exists(path string) (bool, *Stat, error) ***REMOVED***
	res := &existsResponse***REMOVED******REMOVED***
	_, err := c.request(opExists, &existsRequest***REMOVED***Path: path, Watch: false***REMOVED***, res, nil)
	exists := true
	if err == ErrNoNode ***REMOVED***
		exists = false
		err = nil
	***REMOVED***
	return exists, &res.Stat, err
***REMOVED***

func (c *Conn) ExistsW(path string) (bool, *Stat, <-chan Event, error) ***REMOVED***
	var ech <-chan Event
	res := &existsResponse***REMOVED******REMOVED***
	_, err := c.request(opExists, &existsRequest***REMOVED***Path: path, Watch: true***REMOVED***, res, func(req *request, res *responseHeader, err error) ***REMOVED***
		if err == nil ***REMOVED***
			ech = c.addWatcher(path, watchTypeData)
		***REMOVED*** else if err == ErrNoNode ***REMOVED***
			ech = c.addWatcher(path, watchTypeExist)
		***REMOVED***
	***REMOVED***)
	exists := true
	if err == ErrNoNode ***REMOVED***
		exists = false
		err = nil
	***REMOVED***
	if err != nil ***REMOVED***
		return false, nil, nil, err
	***REMOVED***
	return exists, &res.Stat, ech, err
***REMOVED***

func (c *Conn) GetACL(path string) ([]ACL, *Stat, error) ***REMOVED***
	res := &getAclResponse***REMOVED******REMOVED***
	_, err := c.request(opGetAcl, &getAclRequest***REMOVED***Path: path***REMOVED***, res, nil)
	return res.Acl, &res.Stat, err
***REMOVED***

func (c *Conn) SetACL(path string, acl []ACL, version int32) (*Stat, error) ***REMOVED***
	res := &setAclResponse***REMOVED******REMOVED***
	_, err := c.request(opSetAcl, &setAclRequest***REMOVED***Path: path, Acl: acl, Version: version***REMOVED***, res, nil)
	return &res.Stat, err
***REMOVED***

func (c *Conn) Sync(path string) (string, error) ***REMOVED***
	res := &syncResponse***REMOVED******REMOVED***
	_, err := c.request(opSync, &syncRequest***REMOVED***Path: path***REMOVED***, res, nil)
	return res.Path, err
***REMOVED***

type MultiResponse struct ***REMOVED***
	Stat   *Stat
	String string
***REMOVED***

// Multi executes multiple ZooKeeper operations or none of them. The provided
// ops must be one of *CreateRequest, *DeleteRequest, *SetDataRequest, or
// *CheckVersionRequest.
func (c *Conn) Multi(ops ...interface***REMOVED******REMOVED***) ([]MultiResponse, error) ***REMOVED***
	req := &multiRequest***REMOVED***
		Ops:        make([]multiRequestOp, 0, len(ops)),
		DoneHeader: multiHeader***REMOVED***Type: -1, Done: true, Err: -1***REMOVED***,
	***REMOVED***
	for _, op := range ops ***REMOVED***
		var opCode int32
		switch op.(type) ***REMOVED***
		case *CreateRequest:
			opCode = opCreate
		case *SetDataRequest:
			opCode = opSetData
		case *DeleteRequest:
			opCode = opDelete
		case *CheckVersionRequest:
			opCode = opCheck
		default:
			return nil, fmt.Errorf("uknown operation type %T", op)
		***REMOVED***
		req.Ops = append(req.Ops, multiRequestOp***REMOVED***multiHeader***REMOVED***opCode, false, -1***REMOVED***, op***REMOVED***)
	***REMOVED***
	res := &multiResponse***REMOVED******REMOVED***
	_, err := c.request(opMulti, req, res, nil)
	mr := make([]MultiResponse, len(res.Ops))
	for i, op := range res.Ops ***REMOVED***
		mr[i] = MultiResponse***REMOVED***Stat: op.Stat, String: op.String***REMOVED***
	***REMOVED***
	return mr, err
***REMOVED***
