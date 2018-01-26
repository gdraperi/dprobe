package netlink

import (
	"errors"
	"fmt"
	"net"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

const (
	sizeofSocketID      = 0x30
	sizeofSocketRequest = sizeofSocketID + 0x8
	sizeofSocket        = sizeofSocketID + 0x18
)

type socketRequest struct ***REMOVED***
	Family   uint8
	Protocol uint8
	Ext      uint8
	pad      uint8
	States   uint32
	ID       SocketID
***REMOVED***

type writeBuffer struct ***REMOVED***
	Bytes []byte
	pos   int
***REMOVED***

func (b *writeBuffer) Write(c byte) ***REMOVED***
	b.Bytes[b.pos] = c
	b.pos++
***REMOVED***

func (b *writeBuffer) Next(n int) []byte ***REMOVED***
	s := b.Bytes[b.pos : b.pos+n]
	b.pos += n
	return s
***REMOVED***

func (r *socketRequest) Serialize() []byte ***REMOVED***
	b := writeBuffer***REMOVED***Bytes: make([]byte, sizeofSocketRequest)***REMOVED***
	b.Write(r.Family)
	b.Write(r.Protocol)
	b.Write(r.Ext)
	b.Write(r.pad)
	native.PutUint32(b.Next(4), r.States)
	networkOrder.PutUint16(b.Next(2), r.ID.SourcePort)
	networkOrder.PutUint16(b.Next(2), r.ID.DestinationPort)
	copy(b.Next(4), r.ID.Source.To4())
	b.Next(12)
	copy(b.Next(4), r.ID.Destination.To4())
	b.Next(12)
	native.PutUint32(b.Next(4), r.ID.Interface)
	native.PutUint32(b.Next(4), r.ID.Cookie[0])
	native.PutUint32(b.Next(4), r.ID.Cookie[1])
	return b.Bytes
***REMOVED***

func (r *socketRequest) Len() int ***REMOVED*** return sizeofSocketRequest ***REMOVED***

type readBuffer struct ***REMOVED***
	Bytes []byte
	pos   int
***REMOVED***

func (b *readBuffer) Read() byte ***REMOVED***
	c := b.Bytes[b.pos]
	b.pos++
	return c
***REMOVED***

func (b *readBuffer) Next(n int) []byte ***REMOVED***
	s := b.Bytes[b.pos : b.pos+n]
	b.pos += n
	return s
***REMOVED***

func (s *Socket) deserialize(b []byte) error ***REMOVED***
	if len(b) < sizeofSocket ***REMOVED***
		return fmt.Errorf("socket data short read (%d); want %d", len(b), sizeofSocket)
	***REMOVED***
	rb := readBuffer***REMOVED***Bytes: b***REMOVED***
	s.Family = rb.Read()
	s.State = rb.Read()
	s.Timer = rb.Read()
	s.Retrans = rb.Read()
	s.ID.SourcePort = networkOrder.Uint16(rb.Next(2))
	s.ID.DestinationPort = networkOrder.Uint16(rb.Next(2))
	s.ID.Source = net.IPv4(rb.Read(), rb.Read(), rb.Read(), rb.Read())
	rb.Next(12)
	s.ID.Destination = net.IPv4(rb.Read(), rb.Read(), rb.Read(), rb.Read())
	rb.Next(12)
	s.ID.Interface = native.Uint32(rb.Next(4))
	s.ID.Cookie[0] = native.Uint32(rb.Next(4))
	s.ID.Cookie[1] = native.Uint32(rb.Next(4))
	s.Expires = native.Uint32(rb.Next(4))
	s.RQueue = native.Uint32(rb.Next(4))
	s.WQueue = native.Uint32(rb.Next(4))
	s.UID = native.Uint32(rb.Next(4))
	s.INode = native.Uint32(rb.Next(4))
	return nil
***REMOVED***

// SocketGet returns the Socket identified by its local and remote addresses.
func SocketGet(local, remote net.Addr) (*Socket, error) ***REMOVED***
	localTCP, ok := local.(*net.TCPAddr)
	if !ok ***REMOVED***
		return nil, ErrNotImplemented
	***REMOVED***
	remoteTCP, ok := remote.(*net.TCPAddr)
	if !ok ***REMOVED***
		return nil, ErrNotImplemented
	***REMOVED***
	localIP := localTCP.IP.To4()
	if localIP == nil ***REMOVED***
		return nil, ErrNotImplemented
	***REMOVED***
	remoteIP := remoteTCP.IP.To4()
	if remoteIP == nil ***REMOVED***
		return nil, ErrNotImplemented
	***REMOVED***

	s, err := nl.Subscribe(syscall.NETLINK_INET_DIAG)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer s.Close()
	req := nl.NewNetlinkRequest(nl.SOCK_DIAG_BY_FAMILY, 0)
	req.AddData(&socketRequest***REMOVED***
		Family:   syscall.AF_INET,
		Protocol: syscall.IPPROTO_TCP,
		ID: SocketID***REMOVED***
			SourcePort:      uint16(localTCP.Port),
			DestinationPort: uint16(remoteTCP.Port),
			Source:          localIP,
			Destination:     remoteIP,
			Cookie:          [2]uint32***REMOVED***nl.TCPDIAG_NOCOOKIE, nl.TCPDIAG_NOCOOKIE***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	s.Send(req)
	msgs, err := s.Receive()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(msgs) == 0 ***REMOVED***
		return nil, errors.New("no message nor error from netlink")
	***REMOVED***
	if len(msgs) > 2 ***REMOVED***
		return nil, fmt.Errorf("multiple (%d) matching sockets", len(msgs))
	***REMOVED***
	sock := &Socket***REMOVED******REMOVED***
	if err := sock.deserialize(msgs[0].Data); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return sock, nil
***REMOVED***
