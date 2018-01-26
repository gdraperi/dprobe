//+build !windows,!solaris

package dbus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"syscall"
)

type oobReader struct ***REMOVED***
	conn *net.UnixConn
	oob  []byte
	buf  [4096]byte
***REMOVED***

func (o *oobReader) Read(b []byte) (n int, err error) ***REMOVED***
	n, oobn, flags, _, err := o.conn.ReadMsgUnix(b, o.buf[:])
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***
	if flags&syscall.MSG_CTRUNC != 0 ***REMOVED***
		return n, errors.New("dbus: control data truncated (too many fds received)")
	***REMOVED***
	o.oob = append(o.oob, o.buf[:oobn]...)
	return n, nil
***REMOVED***

type unixTransport struct ***REMOVED***
	*net.UnixConn
	hasUnixFDs bool
***REMOVED***

func newUnixTransport(keys string) (transport, error) ***REMOVED***
	var err error

	t := new(unixTransport)
	abstract := getKey(keys, "abstract")
	path := getKey(keys, "path")
	switch ***REMOVED***
	case abstract == "" && path == "":
		return nil, errors.New("dbus: invalid address (neither path nor abstract set)")
	case abstract != "" && path == "":
		t.UnixConn, err = net.DialUnix("unix", nil, &net.UnixAddr***REMOVED***Name: "@" + abstract, Net: "unix"***REMOVED***)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return t, nil
	case abstract == "" && path != "":
		t.UnixConn, err = net.DialUnix("unix", nil, &net.UnixAddr***REMOVED***Name: path, Net: "unix"***REMOVED***)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return t, nil
	default:
		return nil, errors.New("dbus: invalid address (both path and abstract set)")
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	transports["unix"] = newUnixTransport
***REMOVED***

func (t *unixTransport) EnableUnixFDs() ***REMOVED***
	t.hasUnixFDs = true
***REMOVED***

func (t *unixTransport) ReadMessage() (*Message, error) ***REMOVED***
	var (
		blen, hlen uint32
		csheader   [16]byte
		headers    []header
		order      binary.ByteOrder
		unixfds    uint32
	)
	// To be sure that all bytes of out-of-band data are read, we use a special
	// reader that uses ReadUnix on the underlying connection instead of Read
	// and gathers the out-of-band data in a buffer.
	rd := &oobReader***REMOVED***conn: t.UnixConn***REMOVED***
	// read the first 16 bytes (the part of the header that has a constant size),
	// from which we can figure out the length of the rest of the message
	if _, err := io.ReadFull(rd, csheader[:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch csheader[0] ***REMOVED***
	case 'l':
		order = binary.LittleEndian
	case 'B':
		order = binary.BigEndian
	default:
		return nil, InvalidMessageError("invalid byte order")
	***REMOVED***
	// csheader[4:8] -> length of message body, csheader[12:16] -> length of
	// header fields (without alignment)
	binary.Read(bytes.NewBuffer(csheader[4:8]), order, &blen)
	binary.Read(bytes.NewBuffer(csheader[12:]), order, &hlen)
	if hlen%8 != 0 ***REMOVED***
		hlen += 8 - (hlen % 8)
	***REMOVED***

	// decode headers and look for unix fds
	headerdata := make([]byte, hlen+4)
	copy(headerdata, csheader[12:])
	if _, err := io.ReadFull(t, headerdata[4:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	dec := newDecoder(bytes.NewBuffer(headerdata), order)
	dec.pos = 12
	vs, err := dec.Decode(Signature***REMOVED***"a(yv)"***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	Store(vs, &headers)
	for _, v := range headers ***REMOVED***
		if v.Field == byte(FieldUnixFDs) ***REMOVED***
			unixfds, _ = v.Variant.value.(uint32)
		***REMOVED***
	***REMOVED***
	all := make([]byte, 16+hlen+blen)
	copy(all, csheader[:])
	copy(all[16:], headerdata[4:])
	if _, err := io.ReadFull(rd, all[16+hlen:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if unixfds != 0 ***REMOVED***
		if !t.hasUnixFDs ***REMOVED***
			return nil, errors.New("dbus: got unix fds on unsupported transport")
		***REMOVED***
		// read the fds from the OOB data
		scms, err := syscall.ParseSocketControlMessage(rd.oob)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if len(scms) != 1 ***REMOVED***
			return nil, errors.New("dbus: received more than one socket control message")
		***REMOVED***
		fds, err := syscall.ParseUnixRights(&scms[0])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		msg, err := DecodeMessage(bytes.NewBuffer(all))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// substitute the values in the message body (which are indices for the
		// array receiver via OOB) with the actual values
		for i, v := range msg.Body ***REMOVED***
			if j, ok := v.(UnixFDIndex); ok ***REMOVED***
				if uint32(j) >= unixfds ***REMOVED***
					return nil, InvalidMessageError("invalid index for unix fd")
				***REMOVED***
				msg.Body[i] = UnixFD(fds[j])
			***REMOVED***
		***REMOVED***
		return msg, nil
	***REMOVED***
	return DecodeMessage(bytes.NewBuffer(all))
***REMOVED***

func (t *unixTransport) SendMessage(msg *Message) error ***REMOVED***
	fds := make([]int, 0)
	for i, v := range msg.Body ***REMOVED***
		if fd, ok := v.(UnixFD); ok ***REMOVED***
			msg.Body[i] = UnixFDIndex(len(fds))
			fds = append(fds, int(fd))
		***REMOVED***
	***REMOVED***
	if len(fds) != 0 ***REMOVED***
		if !t.hasUnixFDs ***REMOVED***
			return errors.New("dbus: unix fd passing not enabled")
		***REMOVED***
		msg.Headers[FieldUnixFDs] = MakeVariant(uint32(len(fds)))
		oob := syscall.UnixRights(fds...)
		buf := new(bytes.Buffer)
		msg.EncodeTo(buf, binary.LittleEndian)
		n, oobn, err := t.UnixConn.WriteMsgUnix(buf.Bytes(), oob, nil)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if n != buf.Len() || oobn != len(oob) ***REMOVED***
			return io.ErrShortWrite
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := msg.EncodeTo(t, binary.LittleEndian); err != nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (t *unixTransport) SupportsUnixFDs() bool ***REMOVED***
	return true
***REMOVED***
