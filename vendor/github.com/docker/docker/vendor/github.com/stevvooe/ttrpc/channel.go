package ttrpc

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	messageHeaderLength = 10
	messageLengthMax    = 4 << 20
)

type messageType uint8

const (
	messageTypeRequest  messageType = 0x1
	messageTypeResponse messageType = 0x2
)

// messageHeader represents the fixed-length message header of 10 bytes sent
// with every request.
type messageHeader struct ***REMOVED***
	Length   uint32      // length excluding this header. b[:4]
	StreamID uint32      // identifies which request stream message is a part of. b[4:8]
	Type     messageType // message type b[8]
	Flags    uint8       // reserved          b[9]
***REMOVED***

func readMessageHeader(p []byte, r io.Reader) (messageHeader, error) ***REMOVED***
	_, err := io.ReadFull(r, p[:messageHeaderLength])
	if err != nil ***REMOVED***
		return messageHeader***REMOVED******REMOVED***, err
	***REMOVED***

	return messageHeader***REMOVED***
		Length:   binary.BigEndian.Uint32(p[:4]),
		StreamID: binary.BigEndian.Uint32(p[4:8]),
		Type:     messageType(p[8]),
		Flags:    p[9],
	***REMOVED***, nil
***REMOVED***

func writeMessageHeader(w io.Writer, p []byte, mh messageHeader) error ***REMOVED***
	binary.BigEndian.PutUint32(p[:4], mh.Length)
	binary.BigEndian.PutUint32(p[4:8], mh.StreamID)
	p[8] = byte(mh.Type)
	p[9] = mh.Flags

	_, err := w.Write(p[:])
	return err
***REMOVED***

var buffers sync.Pool

type channel struct ***REMOVED***
	bw    *bufio.Writer
	br    *bufio.Reader
	hrbuf [messageHeaderLength]byte // avoid alloc when reading header
	hwbuf [messageHeaderLength]byte
***REMOVED***

func newChannel(w io.Writer, r io.Reader) *channel ***REMOVED***
	return &channel***REMOVED***
		bw: bufio.NewWriter(w),
		br: bufio.NewReader(r),
	***REMOVED***
***REMOVED***

// recv a message from the channel. The returned buffer contains the message.
//
// If a valid grpc status is returned, the message header
// returned will be valid and caller should send that along to
// the correct consumer. The bytes on the underlying channel
// will be discarded.
func (ch *channel) recv(ctx context.Context) (messageHeader, []byte, error) ***REMOVED***
	mh, err := readMessageHeader(ch.hrbuf[:], ch.br)
	if err != nil ***REMOVED***
		return messageHeader***REMOVED******REMOVED***, nil, err
	***REMOVED***

	if mh.Length > uint32(messageLengthMax) ***REMOVED***
		if _, err := ch.br.Discard(int(mh.Length)); err != nil ***REMOVED***
			return mh, nil, errors.Wrapf(err, "failed to discard after receiving oversized message")
		***REMOVED***

		return mh, nil, status.Errorf(codes.ResourceExhausted, "message length %v exceed maximum message size of %v", mh.Length, messageLengthMax)
	***REMOVED***

	p := ch.getmbuf(int(mh.Length))
	if _, err := io.ReadFull(ch.br, p); err != nil ***REMOVED***
		return messageHeader***REMOVED******REMOVED***, nil, errors.Wrapf(err, "failed reading message")
	***REMOVED***

	return mh, p, nil
***REMOVED***

func (ch *channel) send(ctx context.Context, streamID uint32, t messageType, p []byte) error ***REMOVED***
	if err := writeMessageHeader(ch.bw, ch.hwbuf[:], messageHeader***REMOVED***Length: uint32(len(p)), StreamID: streamID, Type: t***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err := ch.bw.Write(p)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return ch.bw.Flush()
***REMOVED***

func (ch *channel) getmbuf(size int) []byte ***REMOVED***
	// we can't use the standard New method on pool because we want to allocate
	// based on size.
	b, ok := buffers.Get().(*[]byte)
	if !ok || cap(*b) < size ***REMOVED***
		// TODO(stevvooe): It may be better to allocate these in fixed length
		// buckets to reduce fragmentation but its not clear that would help
		// with performance. An ilogb approach or similar would work well.
		bb := make([]byte, size)
		b = &bb
	***REMOVED*** else ***REMOVED***
		*b = (*b)[:size]
	***REMOVED***
	return *b
***REMOVED***

func (ch *channel) putmbuf(p []byte) ***REMOVED***
	buffers.Put(&p)
***REMOVED***
