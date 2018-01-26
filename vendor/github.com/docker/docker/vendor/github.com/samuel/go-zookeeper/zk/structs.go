package zk

import (
	"encoding/binary"
	"errors"
	"reflect"
	"runtime"
	"time"
)

var (
	ErrUnhandledFieldType = errors.New("zk: unhandled field type")
	ErrPtrExpected        = errors.New("zk: encode/decode expect a non-nil pointer to struct")
	ErrShortBuffer        = errors.New("zk: buffer too small")
)

type ACL struct ***REMOVED***
	Perms  int32
	Scheme string
	ID     string
***REMOVED***

type Stat struct ***REMOVED***
	Czxid          int64 // The zxid of the change that caused this znode to be created.
	Mzxid          int64 // The zxid of the change that last modified this znode.
	Ctime          int64 // The time in milliseconds from epoch when this znode was created.
	Mtime          int64 // The time in milliseconds from epoch when this znode was last modified.
	Version        int32 // The number of changes to the data of this znode.
	Cversion       int32 // The number of changes to the children of this znode.
	Aversion       int32 // The number of changes to the ACL of this znode.
	EphemeralOwner int64 // The session id of the owner of this znode if the znode is an ephemeral node. If it is not an ephemeral node, it will be zero.
	DataLength     int32 // The length of the data field of this znode.
	NumChildren    int32 // The number of children of this znode.
	Pzxid          int64 // last modified children
***REMOVED***

// ServerClient is the information for a single Zookeeper client and its session.
// This is used to parse/extract the output fo the `cons` command.
type ServerClient struct ***REMOVED***
	Queued        int64
	Received      int64
	Sent          int64
	SessionID     int64
	Lcxid         int64
	Lzxid         int64
	Timeout       int32
	LastLatency   int32
	MinLatency    int32
	AvgLatency    int32
	MaxLatency    int32
	Established   time.Time
	LastResponse  time.Time
	Addr          string
	LastOperation string // maybe?
	Error         error
***REMOVED***

// ServerClients is a struct for the FLWCons() function. It's used to provide
// the list of Clients.
//
// This is needed because FLWCons() takes multiple servers.
type ServerClients struct ***REMOVED***
	Clients []*ServerClient
	Error   error
***REMOVED***

// ServerStats is the information pulled from the Zookeeper `stat` command.
type ServerStats struct ***REMOVED***
	Sent        int64
	Received    int64
	NodeCount   int64
	MinLatency  int64
	AvgLatency  int64
	MaxLatency  int64
	Connections int64
	Outstanding int64
	Epoch       int32
	Counter     int32
	BuildTime   time.Time
	Mode        Mode
	Version     string
	Error       error
***REMOVED***

type requestHeader struct ***REMOVED***
	Xid    int32
	Opcode int32
***REMOVED***

type responseHeader struct ***REMOVED***
	Xid  int32
	Zxid int64
	Err  ErrCode
***REMOVED***

type multiHeader struct ***REMOVED***
	Type int32
	Done bool
	Err  ErrCode
***REMOVED***

type auth struct ***REMOVED***
	Type   int32
	Scheme string
	Auth   []byte
***REMOVED***

// Generic request structs

type pathRequest struct ***REMOVED***
	Path string
***REMOVED***

type PathVersionRequest struct ***REMOVED***
	Path    string
	Version int32
***REMOVED***

type pathWatchRequest struct ***REMOVED***
	Path  string
	Watch bool
***REMOVED***

type pathResponse struct ***REMOVED***
	Path string
***REMOVED***

type statResponse struct ***REMOVED***
	Stat Stat
***REMOVED***

//

type CheckVersionRequest PathVersionRequest
type closeRequest struct***REMOVED******REMOVED***
type closeResponse struct***REMOVED******REMOVED***

type connectRequest struct ***REMOVED***
	ProtocolVersion int32
	LastZxidSeen    int64
	TimeOut         int32
	SessionID       int64
	Passwd          []byte
***REMOVED***

type connectResponse struct ***REMOVED***
	ProtocolVersion int32
	TimeOut         int32
	SessionID       int64
	Passwd          []byte
***REMOVED***

type CreateRequest struct ***REMOVED***
	Path  string
	Data  []byte
	Acl   []ACL
	Flags int32
***REMOVED***

type createResponse pathResponse
type DeleteRequest PathVersionRequest
type deleteResponse struct***REMOVED******REMOVED***

type errorResponse struct ***REMOVED***
	Err int32
***REMOVED***

type existsRequest pathWatchRequest
type existsResponse statResponse
type getAclRequest pathRequest

type getAclResponse struct ***REMOVED***
	Acl  []ACL
	Stat Stat
***REMOVED***

type getChildrenRequest pathRequest

type getChildrenResponse struct ***REMOVED***
	Children []string
***REMOVED***

type getChildren2Request pathWatchRequest

type getChildren2Response struct ***REMOVED***
	Children []string
	Stat     Stat
***REMOVED***

type getDataRequest pathWatchRequest

type getDataResponse struct ***REMOVED***
	Data []byte
	Stat Stat
***REMOVED***

type getMaxChildrenRequest pathRequest

type getMaxChildrenResponse struct ***REMOVED***
	Max int32
***REMOVED***

type getSaslRequest struct ***REMOVED***
	Token []byte
***REMOVED***

type pingRequest struct***REMOVED******REMOVED***
type pingResponse struct***REMOVED******REMOVED***

type setAclRequest struct ***REMOVED***
	Path    string
	Acl     []ACL
	Version int32
***REMOVED***

type setAclResponse statResponse

type SetDataRequest struct ***REMOVED***
	Path    string
	Data    []byte
	Version int32
***REMOVED***

type setDataResponse statResponse

type setMaxChildren struct ***REMOVED***
	Path string
	Max  int32
***REMOVED***

type setSaslRequest struct ***REMOVED***
	Token string
***REMOVED***

type setSaslResponse struct ***REMOVED***
	Token string
***REMOVED***

type setWatchesRequest struct ***REMOVED***
	RelativeZxid int64
	DataWatches  []string
	ExistWatches []string
	ChildWatches []string
***REMOVED***

type setWatchesResponse struct***REMOVED******REMOVED***

type syncRequest pathRequest
type syncResponse pathResponse

type setAuthRequest auth
type setAuthResponse struct***REMOVED******REMOVED***

type multiRequestOp struct ***REMOVED***
	Header multiHeader
	Op     interface***REMOVED******REMOVED***
***REMOVED***
type multiRequest struct ***REMOVED***
	Ops        []multiRequestOp
	DoneHeader multiHeader
***REMOVED***
type multiResponseOp struct ***REMOVED***
	Header multiHeader
	String string
	Stat   *Stat
***REMOVED***
type multiResponse struct ***REMOVED***
	Ops        []multiResponseOp
	DoneHeader multiHeader
***REMOVED***

func (r *multiRequest) Encode(buf []byte) (int, error) ***REMOVED***
	total := 0
	for _, op := range r.Ops ***REMOVED***
		op.Header.Done = false
		n, err := encodePacketValue(buf[total:], reflect.ValueOf(op))
		if err != nil ***REMOVED***
			return total, err
		***REMOVED***
		total += n
	***REMOVED***
	r.DoneHeader.Done = true
	n, err := encodePacketValue(buf[total:], reflect.ValueOf(r.DoneHeader))
	if err != nil ***REMOVED***
		return total, err
	***REMOVED***
	total += n

	return total, nil
***REMOVED***

func (r *multiRequest) Decode(buf []byte) (int, error) ***REMOVED***
	r.Ops = make([]multiRequestOp, 0)
	r.DoneHeader = multiHeader***REMOVED***-1, true, -1***REMOVED***
	total := 0
	for ***REMOVED***
		header := &multiHeader***REMOVED******REMOVED***
		n, err := decodePacketValue(buf[total:], reflect.ValueOf(header))
		if err != nil ***REMOVED***
			return total, err
		***REMOVED***
		total += n
		if header.Done ***REMOVED***
			r.DoneHeader = *header
			break
		***REMOVED***

		req := requestStructForOp(header.Type)
		if req == nil ***REMOVED***
			return total, ErrAPIError
		***REMOVED***
		n, err = decodePacketValue(buf[total:], reflect.ValueOf(req))
		if err != nil ***REMOVED***
			return total, err
		***REMOVED***
		total += n
		r.Ops = append(r.Ops, multiRequestOp***REMOVED****header, req***REMOVED***)
	***REMOVED***
	return total, nil
***REMOVED***

func (r *multiResponse) Decode(buf []byte) (int, error) ***REMOVED***
	r.Ops = make([]multiResponseOp, 0)
	r.DoneHeader = multiHeader***REMOVED***-1, true, -1***REMOVED***
	total := 0
	for ***REMOVED***
		header := &multiHeader***REMOVED******REMOVED***
		n, err := decodePacketValue(buf[total:], reflect.ValueOf(header))
		if err != nil ***REMOVED***
			return total, err
		***REMOVED***
		total += n
		if header.Done ***REMOVED***
			r.DoneHeader = *header
			break
		***REMOVED***

		res := multiResponseOp***REMOVED***Header: *header***REMOVED***
		var w reflect.Value
		switch header.Type ***REMOVED***
		default:
			return total, ErrAPIError
		case opCreate:
			w = reflect.ValueOf(&res.String)
		case opSetData:
			res.Stat = new(Stat)
			w = reflect.ValueOf(res.Stat)
		case opCheck, opDelete:
		***REMOVED***
		if w.IsValid() ***REMOVED***
			n, err := decodePacketValue(buf[total:], w)
			if err != nil ***REMOVED***
				return total, err
			***REMOVED***
			total += n
		***REMOVED***
		r.Ops = append(r.Ops, res)
	***REMOVED***
	return total, nil
***REMOVED***

type watcherEvent struct ***REMOVED***
	Type  EventType
	State State
	Path  string
***REMOVED***

type decoder interface ***REMOVED***
	Decode(buf []byte) (int, error)
***REMOVED***

type encoder interface ***REMOVED***
	Encode(buf []byte) (int, error)
***REMOVED***

func decodePacket(buf []byte, st interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			if e, ok := r.(runtime.Error); ok && e.Error() == "runtime error: slice bounds out of range" ***REMOVED***
				err = ErrShortBuffer
			***REMOVED*** else ***REMOVED***
				panic(r)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.IsNil() ***REMOVED***
		return 0, ErrPtrExpected
	***REMOVED***
	return decodePacketValue(buf, v)
***REMOVED***

func decodePacketValue(buf []byte, v reflect.Value) (int, error) ***REMOVED***
	rv := v
	kind := v.Kind()
	if kind == reflect.Ptr ***REMOVED***
		if v.IsNil() ***REMOVED***
			v.Set(reflect.New(v.Type().Elem()))
		***REMOVED***
		v = v.Elem()
		kind = v.Kind()
	***REMOVED***

	n := 0
	switch kind ***REMOVED***
	default:
		return n, ErrUnhandledFieldType
	case reflect.Struct:
		if de, ok := rv.Interface().(decoder); ok ***REMOVED***
			return de.Decode(buf)
		***REMOVED*** else if de, ok := v.Interface().(decoder); ok ***REMOVED***
			return de.Decode(buf)
		***REMOVED*** else ***REMOVED***
			for i := 0; i < v.NumField(); i++ ***REMOVED***
				field := v.Field(i)
				n2, err := decodePacketValue(buf[n:], field)
				n += n2
				if err != nil ***REMOVED***
					return n, err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case reflect.Bool:
		v.SetBool(buf[n] != 0)
		n++
	case reflect.Int32:
		v.SetInt(int64(binary.BigEndian.Uint32(buf[n : n+4])))
		n += 4
	case reflect.Int64:
		v.SetInt(int64(binary.BigEndian.Uint64(buf[n : n+8])))
		n += 8
	case reflect.String:
		ln := int(binary.BigEndian.Uint32(buf[n : n+4]))
		v.SetString(string(buf[n+4 : n+4+ln]))
		n += 4 + ln
	case reflect.Slice:
		switch v.Type().Elem().Kind() ***REMOVED***
		default:
			count := int(binary.BigEndian.Uint32(buf[n : n+4]))
			n += 4
			values := reflect.MakeSlice(v.Type(), count, count)
			v.Set(values)
			for i := 0; i < count; i++ ***REMOVED***
				n2, err := decodePacketValue(buf[n:], values.Index(i))
				n += n2
				if err != nil ***REMOVED***
					return n, err
				***REMOVED***
			***REMOVED***
		case reflect.Uint8:
			ln := int(int32(binary.BigEndian.Uint32(buf[n : n+4])))
			if ln < 0 ***REMOVED***
				n += 4
				v.SetBytes(nil)
			***REMOVED*** else ***REMOVED***
				bytes := make([]byte, ln)
				copy(bytes, buf[n+4:n+4+ln])
				v.SetBytes(bytes)
				n += 4 + ln
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return n, nil
***REMOVED***

func encodePacket(buf []byte, st interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			if e, ok := r.(runtime.Error); ok && e.Error() == "runtime error: slice bounds out of range" ***REMOVED***
				err = ErrShortBuffer
			***REMOVED*** else ***REMOVED***
				panic(r)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr || v.IsNil() ***REMOVED***
		return 0, ErrPtrExpected
	***REMOVED***
	return encodePacketValue(buf, v)
***REMOVED***

func encodePacketValue(buf []byte, v reflect.Value) (int, error) ***REMOVED***
	rv := v
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface ***REMOVED***
		v = v.Elem()
	***REMOVED***

	n := 0
	switch v.Kind() ***REMOVED***
	default:
		return n, ErrUnhandledFieldType
	case reflect.Struct:
		if en, ok := rv.Interface().(encoder); ok ***REMOVED***
			return en.Encode(buf)
		***REMOVED*** else if en, ok := v.Interface().(encoder); ok ***REMOVED***
			return en.Encode(buf)
		***REMOVED*** else ***REMOVED***
			for i := 0; i < v.NumField(); i++ ***REMOVED***
				field := v.Field(i)
				n2, err := encodePacketValue(buf[n:], field)
				n += n2
				if err != nil ***REMOVED***
					return n, err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case reflect.Bool:
		if v.Bool() ***REMOVED***
			buf[n] = 1
		***REMOVED*** else ***REMOVED***
			buf[n] = 0
		***REMOVED***
		n++
	case reflect.Int32:
		binary.BigEndian.PutUint32(buf[n:n+4], uint32(v.Int()))
		n += 4
	case reflect.Int64:
		binary.BigEndian.PutUint64(buf[n:n+8], uint64(v.Int()))
		n += 8
	case reflect.String:
		str := v.String()
		binary.BigEndian.PutUint32(buf[n:n+4], uint32(len(str)))
		copy(buf[n+4:n+4+len(str)], []byte(str))
		n += 4 + len(str)
	case reflect.Slice:
		switch v.Type().Elem().Kind() ***REMOVED***
		default:
			count := v.Len()
			startN := n
			n += 4
			for i := 0; i < count; i++ ***REMOVED***
				n2, err := encodePacketValue(buf[n:], v.Index(i))
				n += n2
				if err != nil ***REMOVED***
					return n, err
				***REMOVED***
			***REMOVED***
			binary.BigEndian.PutUint32(buf[startN:startN+4], uint32(count))
		case reflect.Uint8:
			if v.IsNil() ***REMOVED***
				binary.BigEndian.PutUint32(buf[n:n+4], uint32(0xffffffff))
				n += 4
			***REMOVED*** else ***REMOVED***
				bytes := v.Bytes()
				binary.BigEndian.PutUint32(buf[n:n+4], uint32(len(bytes)))
				copy(buf[n+4:n+4+len(bytes)], bytes)
				n += 4 + len(bytes)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return n, nil
***REMOVED***

func requestStructForOp(op int32) interface***REMOVED******REMOVED*** ***REMOVED***
	switch op ***REMOVED***
	case opClose:
		return &closeRequest***REMOVED******REMOVED***
	case opCreate:
		return &CreateRequest***REMOVED******REMOVED***
	case opDelete:
		return &DeleteRequest***REMOVED******REMOVED***
	case opExists:
		return &existsRequest***REMOVED******REMOVED***
	case opGetAcl:
		return &getAclRequest***REMOVED******REMOVED***
	case opGetChildren:
		return &getChildrenRequest***REMOVED******REMOVED***
	case opGetChildren2:
		return &getChildren2Request***REMOVED******REMOVED***
	case opGetData:
		return &getDataRequest***REMOVED******REMOVED***
	case opPing:
		return &pingRequest***REMOVED******REMOVED***
	case opSetAcl:
		return &setAclRequest***REMOVED******REMOVED***
	case opSetData:
		return &SetDataRequest***REMOVED******REMOVED***
	case opSetWatches:
		return &setWatchesRequest***REMOVED******REMOVED***
	case opSync:
		return &syncRequest***REMOVED******REMOVED***
	case opSetAuth:
		return &setAuthRequest***REMOVED******REMOVED***
	case opCheck:
		return &CheckVersionRequest***REMOVED******REMOVED***
	case opMulti:
		return &multiRequest***REMOVED******REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func responseStructForOp(op int32) interface***REMOVED******REMOVED*** ***REMOVED***
	switch op ***REMOVED***
	case opClose:
		return &closeResponse***REMOVED******REMOVED***
	case opCreate:
		return &createResponse***REMOVED******REMOVED***
	case opDelete:
		return &deleteResponse***REMOVED******REMOVED***
	case opExists:
		return &existsResponse***REMOVED******REMOVED***
	case opGetAcl:
		return &getAclResponse***REMOVED******REMOVED***
	case opGetChildren:
		return &getChildrenResponse***REMOVED******REMOVED***
	case opGetChildren2:
		return &getChildren2Response***REMOVED******REMOVED***
	case opGetData:
		return &getDataResponse***REMOVED******REMOVED***
	case opPing:
		return &pingResponse***REMOVED******REMOVED***
	case opSetAcl:
		return &setAclResponse***REMOVED******REMOVED***
	case opSetData:
		return &setDataResponse***REMOVED******REMOVED***
	case opSetWatches:
		return &setWatchesResponse***REMOVED******REMOVED***
	case opSync:
		return &syncResponse***REMOVED******REMOVED***
	case opWatcherEvent:
		return &watcherEvent***REMOVED******REMOVED***
	case opSetAuth:
		return &setAuthResponse***REMOVED******REMOVED***
	// case opCheck:
	// 	return &checkVersionResponse***REMOVED******REMOVED***
	case opMulti:
		return &multiResponse***REMOVED******REMOVED***
	***REMOVED***
	return nil
***REMOVED***
