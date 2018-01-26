/*
 *
 * Copyright 2014, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package transport

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

const (
	// http2MaxFrameLen specifies the max length of a HTTP2 frame.
	http2MaxFrameLen = 16384 // 16KB frame
	// http://http2.github.io/http2-spec/#SettingValues
	http2InitHeaderTableSize = 4096
	// http2IOBufSize specifies the buffer size for sending frames.
	http2IOBufSize = 32 * 1024
)

var (
	clientPreface   = []byte(http2.ClientPreface)
	http2ErrConvTab = map[http2.ErrCode]codes.Code***REMOVED***
		http2.ErrCodeNo:                 codes.Internal,
		http2.ErrCodeProtocol:           codes.Internal,
		http2.ErrCodeInternal:           codes.Internal,
		http2.ErrCodeFlowControl:        codes.ResourceExhausted,
		http2.ErrCodeSettingsTimeout:    codes.Internal,
		http2.ErrCodeStreamClosed:       codes.Internal,
		http2.ErrCodeFrameSize:          codes.Internal,
		http2.ErrCodeRefusedStream:      codes.Unavailable,
		http2.ErrCodeCancel:             codes.Canceled,
		http2.ErrCodeCompression:        codes.Internal,
		http2.ErrCodeConnect:            codes.Internal,
		http2.ErrCodeEnhanceYourCalm:    codes.ResourceExhausted,
		http2.ErrCodeInadequateSecurity: codes.PermissionDenied,
		http2.ErrCodeHTTP11Required:     codes.FailedPrecondition,
	***REMOVED***
	statusCodeConvTab = map[codes.Code]http2.ErrCode***REMOVED***
		codes.Internal:          http2.ErrCodeInternal,
		codes.Canceled:          http2.ErrCodeCancel,
		codes.Unavailable:       http2.ErrCodeRefusedStream,
		codes.ResourceExhausted: http2.ErrCodeEnhanceYourCalm,
		codes.PermissionDenied:  http2.ErrCodeInadequateSecurity,
	***REMOVED***
)

// Records the states during HPACK decoding. Must be reset once the
// decoding of the entire headers are finished.
type decodeState struct ***REMOVED***
	encoding string
	// statusGen caches the stream status received from the trailer the server
	// sent.  Client side only.  Do not access directly.  After all trailers are
	// parsed, use the status method to retrieve the status.
	statusGen *status.Status
	// rawStatusCode and rawStatusMsg are set from the raw trailer fields and are not
	// intended for direct access outside of parsing.
	rawStatusCode int32
	rawStatusMsg  string
	// Server side only fields.
	timeoutSet bool
	timeout    time.Duration
	method     string
	// key-value metadata map from the peer.
	mdata map[string][]string
***REMOVED***

// isReservedHeader checks whether hdr belongs to HTTP2 headers
// reserved by gRPC protocol. Any other headers are classified as the
// user-specified metadata.
func isReservedHeader(hdr string) bool ***REMOVED***
	if hdr != "" && hdr[0] == ':' ***REMOVED***
		return true
	***REMOVED***
	switch hdr ***REMOVED***
	case "content-type",
		"grpc-message-type",
		"grpc-encoding",
		"grpc-message",
		"grpc-status",
		"grpc-timeout",
		"grpc-status-details-bin",
		"te":
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// isWhitelistedPseudoHeader checks whether hdr belongs to HTTP2 pseudoheaders
// that should be propagated into metadata visible to users.
func isWhitelistedPseudoHeader(hdr string) bool ***REMOVED***
	switch hdr ***REMOVED***
	case ":authority":
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func validContentType(t string) bool ***REMOVED***
	e := "application/grpc"
	if !strings.HasPrefix(t, e) ***REMOVED***
		return false
	***REMOVED***
	// Support variations on the content-type
	// (e.g. "application/grpc+blah", "application/grpc;blah").
	if len(t) > len(e) && t[len(e)] != '+' && t[len(e)] != ';' ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func (d *decodeState) status() *status.Status ***REMOVED***
	if d.statusGen == nil ***REMOVED***
		// No status-details were provided; generate status using code/msg.
		d.statusGen = status.New(codes.Code(d.rawStatusCode), d.rawStatusMsg)
	***REMOVED***
	return d.statusGen
***REMOVED***

const binHdrSuffix = "-bin"

func encodeBinHeader(v []byte) string ***REMOVED***
	return base64.RawStdEncoding.EncodeToString(v)
***REMOVED***

func decodeBinHeader(v string) ([]byte, error) ***REMOVED***
	if len(v)%4 == 0 ***REMOVED***
		// Input was padded, or padding was not necessary.
		return base64.StdEncoding.DecodeString(v)
	***REMOVED***
	return base64.RawStdEncoding.DecodeString(v)
***REMOVED***

func encodeMetadataHeader(k, v string) string ***REMOVED***
	if strings.HasSuffix(k, binHdrSuffix) ***REMOVED***
		return encodeBinHeader(([]byte)(v))
	***REMOVED***
	return v
***REMOVED***

func decodeMetadataHeader(k, v string) (string, error) ***REMOVED***
	if strings.HasSuffix(k, binHdrSuffix) ***REMOVED***
		b, err := decodeBinHeader(v)
		return string(b), err
	***REMOVED***
	return v, nil
***REMOVED***

func (d *decodeState) processHeaderField(f hpack.HeaderField) error ***REMOVED***
	switch f.Name ***REMOVED***
	case "content-type":
		if !validContentType(f.Value) ***REMOVED***
			return streamErrorf(codes.FailedPrecondition, "transport: received the unexpected content-type %q", f.Value)
		***REMOVED***
	case "grpc-encoding":
		d.encoding = f.Value
	case "grpc-status":
		code, err := strconv.Atoi(f.Value)
		if err != nil ***REMOVED***
			return streamErrorf(codes.Internal, "transport: malformed grpc-status: %v", err)
		***REMOVED***
		d.rawStatusCode = int32(code)
	case "grpc-message":
		d.rawStatusMsg = decodeGrpcMessage(f.Value)
	case "grpc-status-details-bin":
		v, err := decodeBinHeader(f.Value)
		if err != nil ***REMOVED***
			return streamErrorf(codes.Internal, "transport: malformed grpc-status-details-bin: %v", err)
		***REMOVED***
		s := &spb.Status***REMOVED******REMOVED***
		if err := proto.Unmarshal(v, s); err != nil ***REMOVED***
			return streamErrorf(codes.Internal, "transport: malformed grpc-status-details-bin: %v", err)
		***REMOVED***
		d.statusGen = status.FromProto(s)
	case "grpc-timeout":
		d.timeoutSet = true
		var err error
		if d.timeout, err = decodeTimeout(f.Value); err != nil ***REMOVED***
			return streamErrorf(codes.Internal, "transport: malformed time-out: %v", err)
		***REMOVED***
	case ":path":
		d.method = f.Value
	default:
		if !isReservedHeader(f.Name) || isWhitelistedPseudoHeader(f.Name) ***REMOVED***
			if d.mdata == nil ***REMOVED***
				d.mdata = make(map[string][]string)
			***REMOVED***
			v, err := decodeMetadataHeader(f.Name, f.Value)
			if err != nil ***REMOVED***
				grpclog.Printf("Failed to decode (%q, %q): %v", f.Name, f.Value, err)
				return nil
			***REMOVED***
			d.mdata[f.Name] = append(d.mdata[f.Name], v)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type timeoutUnit uint8

const (
	hour        timeoutUnit = 'H'
	minute      timeoutUnit = 'M'
	second      timeoutUnit = 'S'
	millisecond timeoutUnit = 'm'
	microsecond timeoutUnit = 'u'
	nanosecond  timeoutUnit = 'n'
)

func timeoutUnitToDuration(u timeoutUnit) (d time.Duration, ok bool) ***REMOVED***
	switch u ***REMOVED***
	case hour:
		return time.Hour, true
	case minute:
		return time.Minute, true
	case second:
		return time.Second, true
	case millisecond:
		return time.Millisecond, true
	case microsecond:
		return time.Microsecond, true
	case nanosecond:
		return time.Nanosecond, true
	default:
	***REMOVED***
	return
***REMOVED***

const maxTimeoutValue int64 = 100000000 - 1

// div does integer division and round-up the result. Note that this is
// equivalent to (d+r-1)/r but has less chance to overflow.
func div(d, r time.Duration) int64 ***REMOVED***
	if m := d % r; m > 0 ***REMOVED***
		return int64(d/r + 1)
	***REMOVED***
	return int64(d / r)
***REMOVED***

// TODO(zhaoq): It is the simplistic and not bandwidth efficient. Improve it.
func encodeTimeout(t time.Duration) string ***REMOVED***
	if t <= 0 ***REMOVED***
		return "0n"
	***REMOVED***
	if d := div(t, time.Nanosecond); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "n"
	***REMOVED***
	if d := div(t, time.Microsecond); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "u"
	***REMOVED***
	if d := div(t, time.Millisecond); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "m"
	***REMOVED***
	if d := div(t, time.Second); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "S"
	***REMOVED***
	if d := div(t, time.Minute); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "M"
	***REMOVED***
	// Note that maxTimeoutValue * time.Hour > MaxInt64.
	return strconv.FormatInt(div(t, time.Hour), 10) + "H"
***REMOVED***

func decodeTimeout(s string) (time.Duration, error) ***REMOVED***
	size := len(s)
	if size < 2 ***REMOVED***
		return 0, fmt.Errorf("transport: timeout string is too short: %q", s)
	***REMOVED***
	unit := timeoutUnit(s[size-1])
	d, ok := timeoutUnitToDuration(unit)
	if !ok ***REMOVED***
		return 0, fmt.Errorf("transport: timeout unit is not recognized: %q", s)
	***REMOVED***
	t, err := strconv.ParseInt(s[:size-1], 10, 64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return d * time.Duration(t), nil
***REMOVED***

const (
	spaceByte   = ' '
	tildaByte   = '~'
	percentByte = '%'
)

// encodeGrpcMessage is used to encode status code in header field
// "grpc-message".
// It checks to see if each individual byte in msg is an
// allowable byte, and then either percent encoding or passing it through.
// When percent encoding, the byte is converted into hexadecimal notation
// with a '%' prepended.
func encodeGrpcMessage(msg string) string ***REMOVED***
	if msg == "" ***REMOVED***
		return ""
	***REMOVED***
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ ***REMOVED***
		c := msg[i]
		if !(c >= spaceByte && c < tildaByte && c != percentByte) ***REMOVED***
			return encodeGrpcMessageUnchecked(msg)
		***REMOVED***
	***REMOVED***
	return msg
***REMOVED***

func encodeGrpcMessageUnchecked(msg string) string ***REMOVED***
	var buf bytes.Buffer
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ ***REMOVED***
		c := msg[i]
		if c >= spaceByte && c < tildaByte && c != percentByte ***REMOVED***
			buf.WriteByte(c)
		***REMOVED*** else ***REMOVED***
			buf.WriteString(fmt.Sprintf("%%%02X", c))
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***

// decodeGrpcMessage decodes the msg encoded by encodeGrpcMessage.
func decodeGrpcMessage(msg string) string ***REMOVED***
	if msg == "" ***REMOVED***
		return ""
	***REMOVED***
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ ***REMOVED***
		if msg[i] == percentByte && i+2 < lenMsg ***REMOVED***
			return decodeGrpcMessageUnchecked(msg)
		***REMOVED***
	***REMOVED***
	return msg
***REMOVED***

func decodeGrpcMessageUnchecked(msg string) string ***REMOVED***
	var buf bytes.Buffer
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ ***REMOVED***
		c := msg[i]
		if c == percentByte && i+2 < lenMsg ***REMOVED***
			parsed, err := strconv.ParseUint(msg[i+1:i+3], 16, 8)
			if err != nil ***REMOVED***
				buf.WriteByte(c)
			***REMOVED*** else ***REMOVED***
				buf.WriteByte(byte(parsed))
				i += 2
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			buf.WriteByte(c)
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***

type framer struct ***REMOVED***
	numWriters int32
	reader     io.Reader
	writer     *bufio.Writer
	fr         *http2.Framer
***REMOVED***

func newFramer(conn net.Conn) *framer ***REMOVED***
	f := &framer***REMOVED***
		reader: bufio.NewReaderSize(conn, http2IOBufSize),
		writer: bufio.NewWriterSize(conn, http2IOBufSize),
	***REMOVED***
	f.fr = http2.NewFramer(f.writer, f.reader)
	// Opt-in to Frame reuse API on framer to reduce garbage.
	// Frames aren't safe to read from after a subsequent call to ReadFrame.
	f.fr.SetReuseFrames()
	f.fr.ReadMetaHeaders = hpack.NewDecoder(http2InitHeaderTableSize, nil)
	return f
***REMOVED***

func (f *framer) adjustNumWriters(i int32) int32 ***REMOVED***
	return atomic.AddInt32(&f.numWriters, i)
***REMOVED***

// The following writeXXX functions can only be called when the caller gets
// unblocked from writableChan channel (i.e., owns the privilege to write).

func (f *framer) writeContinuation(forceFlush bool, streamID uint32, endHeaders bool, headerBlockFragment []byte) error ***REMOVED***
	if err := f.fr.WriteContinuation(streamID, endHeaders, headerBlockFragment); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writeData(forceFlush bool, streamID uint32, endStream bool, data []byte) error ***REMOVED***
	if err := f.fr.WriteData(streamID, endStream, data); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writeGoAway(forceFlush bool, maxStreamID uint32, code http2.ErrCode, debugData []byte) error ***REMOVED***
	if err := f.fr.WriteGoAway(maxStreamID, code, debugData); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writeHeaders(forceFlush bool, p http2.HeadersFrameParam) error ***REMOVED***
	if err := f.fr.WriteHeaders(p); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writePing(forceFlush, ack bool, data [8]byte) error ***REMOVED***
	if err := f.fr.WritePing(ack, data); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writePriority(forceFlush bool, streamID uint32, p http2.PriorityParam) error ***REMOVED***
	if err := f.fr.WritePriority(streamID, p); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writePushPromise(forceFlush bool, p http2.PushPromiseParam) error ***REMOVED***
	if err := f.fr.WritePushPromise(p); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writeRSTStream(forceFlush bool, streamID uint32, code http2.ErrCode) error ***REMOVED***
	if err := f.fr.WriteRSTStream(streamID, code); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writeSettings(forceFlush bool, settings ...http2.Setting) error ***REMOVED***
	if err := f.fr.WriteSettings(settings...); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writeSettingsAck(forceFlush bool) error ***REMOVED***
	if err := f.fr.WriteSettingsAck(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) writeWindowUpdate(forceFlush bool, streamID, incr uint32) error ***REMOVED***
	if err := f.fr.WriteWindowUpdate(streamID, incr); err != nil ***REMOVED***
		return err
	***REMOVED***
	if forceFlush ***REMOVED***
		return f.writer.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (f *framer) flushWrite() error ***REMOVED***
	return f.writer.Flush()
***REMOVED***

func (f *framer) readFrame() (http2.Frame, error) ***REMOVED***
	return f.fr.ReadFrame()
***REMOVED***

func (f *framer) errorDetail() error ***REMOVED***
	return f.fr.ErrorDetail()
***REMOVED***
