// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// These are SSH message type numbers. They are scattered around several
// documents but many were taken from [SSH-PARAMETERS].
const (
	msgIgnore        = 2
	msgUnimplemented = 3
	msgDebug         = 4
	msgNewKeys       = 21
)

// SSH messages:
//
// These structures mirror the wire format of the corresponding SSH messages.
// They are marshaled using reflection with the marshal and unmarshal functions
// in this file. The only wrinkle is that a final member of type []byte with a
// ssh tag of "rest" receives the remainder of a packet when unmarshaling.

// See RFC 4253, section 11.1.
const msgDisconnect = 1

// disconnectMsg is the message that signals a disconnect. It is also
// the error type returned from mux.Wait()
type disconnectMsg struct ***REMOVED***
	Reason   uint32 `sshtype:"1"`
	Message  string
	Language string
***REMOVED***

func (d *disconnectMsg) Error() string ***REMOVED***
	return fmt.Sprintf("ssh: disconnect, reason %d: %s", d.Reason, d.Message)
***REMOVED***

// See RFC 4253, section 7.1.
const msgKexInit = 20

type kexInitMsg struct ***REMOVED***
	Cookie                  [16]byte `sshtype:"20"`
	KexAlgos                []string
	ServerHostKeyAlgos      []string
	CiphersClientServer     []string
	CiphersServerClient     []string
	MACsClientServer        []string
	MACsServerClient        []string
	CompressionClientServer []string
	CompressionServerClient []string
	LanguagesClientServer   []string
	LanguagesServerClient   []string
	FirstKexFollows         bool
	Reserved                uint32
***REMOVED***

// See RFC 4253, section 8.

// Diffie-Helman
const msgKexDHInit = 30

type kexDHInitMsg struct ***REMOVED***
	X *big.Int `sshtype:"30"`
***REMOVED***

const msgKexECDHInit = 30

type kexECDHInitMsg struct ***REMOVED***
	ClientPubKey []byte `sshtype:"30"`
***REMOVED***

const msgKexECDHReply = 31

type kexECDHReplyMsg struct ***REMOVED***
	HostKey         []byte `sshtype:"31"`
	EphemeralPubKey []byte
	Signature       []byte
***REMOVED***

const msgKexDHReply = 31

type kexDHReplyMsg struct ***REMOVED***
	HostKey   []byte `sshtype:"31"`
	Y         *big.Int
	Signature []byte
***REMOVED***

// See RFC 4253, section 10.
const msgServiceRequest = 5

type serviceRequestMsg struct ***REMOVED***
	Service string `sshtype:"5"`
***REMOVED***

// See RFC 4253, section 10.
const msgServiceAccept = 6

type serviceAcceptMsg struct ***REMOVED***
	Service string `sshtype:"6"`
***REMOVED***

// See RFC 4252, section 5.
const msgUserAuthRequest = 50

type userAuthRequestMsg struct ***REMOVED***
	User    string `sshtype:"50"`
	Service string
	Method  string
	Payload []byte `ssh:"rest"`
***REMOVED***

// Used for debug printouts of packets.
type userAuthSuccessMsg struct ***REMOVED***
***REMOVED***

// See RFC 4252, section 5.1
const msgUserAuthFailure = 51

type userAuthFailureMsg struct ***REMOVED***
	Methods        []string `sshtype:"51"`
	PartialSuccess bool
***REMOVED***

// See RFC 4252, section 5.1
const msgUserAuthSuccess = 52

// See RFC 4252, section 5.4
const msgUserAuthBanner = 53

type userAuthBannerMsg struct ***REMOVED***
	Message string `sshtype:"53"`
	// unused, but required to allow message parsing
	Language string
***REMOVED***

// See RFC 4256, section 3.2
const msgUserAuthInfoRequest = 60
const msgUserAuthInfoResponse = 61

type userAuthInfoRequestMsg struct ***REMOVED***
	User               string `sshtype:"60"`
	Instruction        string
	DeprecatedLanguage string
	NumPrompts         uint32
	Prompts            []byte `ssh:"rest"`
***REMOVED***

// See RFC 4254, section 5.1.
const msgChannelOpen = 90

type channelOpenMsg struct ***REMOVED***
	ChanType         string `sshtype:"90"`
	PeersID          uint32
	PeersWindow      uint32
	MaxPacketSize    uint32
	TypeSpecificData []byte `ssh:"rest"`
***REMOVED***

const msgChannelExtendedData = 95
const msgChannelData = 94

// Used for debug print outs of packets.
type channelDataMsg struct ***REMOVED***
	PeersID uint32 `sshtype:"94"`
	Length  uint32
	Rest    []byte `ssh:"rest"`
***REMOVED***

// See RFC 4254, section 5.1.
const msgChannelOpenConfirm = 91

type channelOpenConfirmMsg struct ***REMOVED***
	PeersID          uint32 `sshtype:"91"`
	MyID             uint32
	MyWindow         uint32
	MaxPacketSize    uint32
	TypeSpecificData []byte `ssh:"rest"`
***REMOVED***

// See RFC 4254, section 5.1.
const msgChannelOpenFailure = 92

type channelOpenFailureMsg struct ***REMOVED***
	PeersID  uint32 `sshtype:"92"`
	Reason   RejectionReason
	Message  string
	Language string
***REMOVED***

const msgChannelRequest = 98

type channelRequestMsg struct ***REMOVED***
	PeersID             uint32 `sshtype:"98"`
	Request             string
	WantReply           bool
	RequestSpecificData []byte `ssh:"rest"`
***REMOVED***

// See RFC 4254, section 5.4.
const msgChannelSuccess = 99

type channelRequestSuccessMsg struct ***REMOVED***
	PeersID uint32 `sshtype:"99"`
***REMOVED***

// See RFC 4254, section 5.4.
const msgChannelFailure = 100

type channelRequestFailureMsg struct ***REMOVED***
	PeersID uint32 `sshtype:"100"`
***REMOVED***

// See RFC 4254, section 5.3
const msgChannelClose = 97

type channelCloseMsg struct ***REMOVED***
	PeersID uint32 `sshtype:"97"`
***REMOVED***

// See RFC 4254, section 5.3
const msgChannelEOF = 96

type channelEOFMsg struct ***REMOVED***
	PeersID uint32 `sshtype:"96"`
***REMOVED***

// See RFC 4254, section 4
const msgGlobalRequest = 80

type globalRequestMsg struct ***REMOVED***
	Type      string `sshtype:"80"`
	WantReply bool
	Data      []byte `ssh:"rest"`
***REMOVED***

// See RFC 4254, section 4
const msgRequestSuccess = 81

type globalRequestSuccessMsg struct ***REMOVED***
	Data []byte `ssh:"rest" sshtype:"81"`
***REMOVED***

// See RFC 4254, section 4
const msgRequestFailure = 82

type globalRequestFailureMsg struct ***REMOVED***
	Data []byte `ssh:"rest" sshtype:"82"`
***REMOVED***

// See RFC 4254, section 5.2
const msgChannelWindowAdjust = 93

type windowAdjustMsg struct ***REMOVED***
	PeersID         uint32 `sshtype:"93"`
	AdditionalBytes uint32
***REMOVED***

// See RFC 4252, section 7
const msgUserAuthPubKeyOk = 60

type userAuthPubKeyOkMsg struct ***REMOVED***
	Algo   string `sshtype:"60"`
	PubKey []byte
***REMOVED***

// typeTags returns the possible type bytes for the given reflect.Type, which
// should be a struct. The possible values are separated by a '|' character.
func typeTags(structType reflect.Type) (tags []byte) ***REMOVED***
	tagStr := structType.Field(0).Tag.Get("sshtype")

	for _, tag := range strings.Split(tagStr, "|") ***REMOVED***
		i, err := strconv.Atoi(tag)
		if err == nil ***REMOVED***
			tags = append(tags, byte(i))
		***REMOVED***
	***REMOVED***

	return tags
***REMOVED***

func fieldError(t reflect.Type, field int, problem string) error ***REMOVED***
	if problem != "" ***REMOVED***
		problem = ": " + problem
	***REMOVED***
	return fmt.Errorf("ssh: unmarshal error for field %s of type %s%s", t.Field(field).Name, t.Name(), problem)
***REMOVED***

var errShortRead = errors.New("ssh: short read")

// Unmarshal parses data in SSH wire format into a structure. The out
// argument should be a pointer to struct. If the first member of the
// struct has the "sshtype" tag set to a '|'-separated set of numbers
// in decimal, the packet must start with one of those numbers. In
// case of error, Unmarshal returns a ParseError or
// UnexpectedMessageError.
func Unmarshal(data []byte, out interface***REMOVED******REMOVED***) error ***REMOVED***
	v := reflect.ValueOf(out).Elem()
	structType := v.Type()
	expectedTypes := typeTags(structType)

	var expectedType byte
	if len(expectedTypes) > 0 ***REMOVED***
		expectedType = expectedTypes[0]
	***REMOVED***

	if len(data) == 0 ***REMOVED***
		return parseError(expectedType)
	***REMOVED***

	if len(expectedTypes) > 0 ***REMOVED***
		goodType := false
		for _, e := range expectedTypes ***REMOVED***
			if e > 0 && data[0] == e ***REMOVED***
				goodType = true
				break
			***REMOVED***
		***REMOVED***
		if !goodType ***REMOVED***
			return fmt.Errorf("ssh: unexpected message type %d (expected one of %v)", data[0], expectedTypes)
		***REMOVED***
		data = data[1:]
	***REMOVED***

	var ok bool
	for i := 0; i < v.NumField(); i++ ***REMOVED***
		field := v.Field(i)
		t := field.Type()
		switch t.Kind() ***REMOVED***
		case reflect.Bool:
			if len(data) < 1 ***REMOVED***
				return errShortRead
			***REMOVED***
			field.SetBool(data[0] != 0)
			data = data[1:]
		case reflect.Array:
			if t.Elem().Kind() != reflect.Uint8 ***REMOVED***
				return fieldError(structType, i, "array of unsupported type")
			***REMOVED***
			if len(data) < t.Len() ***REMOVED***
				return errShortRead
			***REMOVED***
			for j, n := 0, t.Len(); j < n; j++ ***REMOVED***
				field.Index(j).Set(reflect.ValueOf(data[j]))
			***REMOVED***
			data = data[t.Len():]
		case reflect.Uint64:
			var u64 uint64
			if u64, data, ok = parseUint64(data); !ok ***REMOVED***
				return errShortRead
			***REMOVED***
			field.SetUint(u64)
		case reflect.Uint32:
			var u32 uint32
			if u32, data, ok = parseUint32(data); !ok ***REMOVED***
				return errShortRead
			***REMOVED***
			field.SetUint(uint64(u32))
		case reflect.Uint8:
			if len(data) < 1 ***REMOVED***
				return errShortRead
			***REMOVED***
			field.SetUint(uint64(data[0]))
			data = data[1:]
		case reflect.String:
			var s []byte
			if s, data, ok = parseString(data); !ok ***REMOVED***
				return fieldError(structType, i, "")
			***REMOVED***
			field.SetString(string(s))
		case reflect.Slice:
			switch t.Elem().Kind() ***REMOVED***
			case reflect.Uint8:
				if structType.Field(i).Tag.Get("ssh") == "rest" ***REMOVED***
					field.Set(reflect.ValueOf(data))
					data = nil
				***REMOVED*** else ***REMOVED***
					var s []byte
					if s, data, ok = parseString(data); !ok ***REMOVED***
						return errShortRead
					***REMOVED***
					field.Set(reflect.ValueOf(s))
				***REMOVED***
			case reflect.String:
				var nl []string
				if nl, data, ok = parseNameList(data); !ok ***REMOVED***
					return errShortRead
				***REMOVED***
				field.Set(reflect.ValueOf(nl))
			default:
				return fieldError(structType, i, "slice of unsupported type")
			***REMOVED***
		case reflect.Ptr:
			if t == bigIntType ***REMOVED***
				var n *big.Int
				if n, data, ok = parseInt(data); !ok ***REMOVED***
					return errShortRead
				***REMOVED***
				field.Set(reflect.ValueOf(n))
			***REMOVED*** else ***REMOVED***
				return fieldError(structType, i, "pointer to unsupported type")
			***REMOVED***
		default:
			return fieldError(structType, i, fmt.Sprintf("unsupported type: %v", t))
		***REMOVED***
	***REMOVED***

	if len(data) != 0 ***REMOVED***
		return parseError(expectedType)
	***REMOVED***

	return nil
***REMOVED***

// Marshal serializes the message in msg to SSH wire format.  The msg
// argument should be a struct or pointer to struct. If the first
// member has the "sshtype" tag set to a number in decimal, that
// number is prepended to the result. If the last of member has the
// "ssh" tag set to "rest", its contents are appended to the output.
func Marshal(msg interface***REMOVED******REMOVED***) []byte ***REMOVED***
	out := make([]byte, 0, 64)
	return marshalStruct(out, msg)
***REMOVED***

func marshalStruct(out []byte, msg interface***REMOVED******REMOVED***) []byte ***REMOVED***
	v := reflect.Indirect(reflect.ValueOf(msg))
	msgTypes := typeTags(v.Type())
	if len(msgTypes) > 0 ***REMOVED***
		out = append(out, msgTypes[0])
	***REMOVED***

	for i, n := 0, v.NumField(); i < n; i++ ***REMOVED***
		field := v.Field(i)
		switch t := field.Type(); t.Kind() ***REMOVED***
		case reflect.Bool:
			var v uint8
			if field.Bool() ***REMOVED***
				v = 1
			***REMOVED***
			out = append(out, v)
		case reflect.Array:
			if t.Elem().Kind() != reflect.Uint8 ***REMOVED***
				panic(fmt.Sprintf("array of non-uint8 in field %d: %T", i, field.Interface()))
			***REMOVED***
			for j, l := 0, t.Len(); j < l; j++ ***REMOVED***
				out = append(out, uint8(field.Index(j).Uint()))
			***REMOVED***
		case reflect.Uint32:
			out = appendU32(out, uint32(field.Uint()))
		case reflect.Uint64:
			out = appendU64(out, uint64(field.Uint()))
		case reflect.Uint8:
			out = append(out, uint8(field.Uint()))
		case reflect.String:
			s := field.String()
			out = appendInt(out, len(s))
			out = append(out, s...)
		case reflect.Slice:
			switch t.Elem().Kind() ***REMOVED***
			case reflect.Uint8:
				if v.Type().Field(i).Tag.Get("ssh") != "rest" ***REMOVED***
					out = appendInt(out, field.Len())
				***REMOVED***
				out = append(out, field.Bytes()...)
			case reflect.String:
				offset := len(out)
				out = appendU32(out, 0)
				if n := field.Len(); n > 0 ***REMOVED***
					for j := 0; j < n; j++ ***REMOVED***
						f := field.Index(j)
						if j != 0 ***REMOVED***
							out = append(out, ',')
						***REMOVED***
						out = append(out, f.String()...)
					***REMOVED***
					// overwrite length value
					binary.BigEndian.PutUint32(out[offset:], uint32(len(out)-offset-4))
				***REMOVED***
			default:
				panic(fmt.Sprintf("slice of unknown type in field %d: %T", i, field.Interface()))
			***REMOVED***
		case reflect.Ptr:
			if t == bigIntType ***REMOVED***
				var n *big.Int
				nValue := reflect.ValueOf(&n)
				nValue.Elem().Set(field)
				needed := intLength(n)
				oldLength := len(out)

				if cap(out)-len(out) < needed ***REMOVED***
					newOut := make([]byte, len(out), 2*(len(out)+needed))
					copy(newOut, out)
					out = newOut
				***REMOVED***
				out = out[:oldLength+needed]
				marshalInt(out[oldLength:], n)
			***REMOVED*** else ***REMOVED***
				panic(fmt.Sprintf("pointer to unknown type in field %d: %T", i, field.Interface()))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return out
***REMOVED***

var bigOne = big.NewInt(1)

func parseString(in []byte) (out, rest []byte, ok bool) ***REMOVED***
	if len(in) < 4 ***REMOVED***
		return
	***REMOVED***
	length := binary.BigEndian.Uint32(in)
	in = in[4:]
	if uint32(len(in)) < length ***REMOVED***
		return
	***REMOVED***
	out = in[:length]
	rest = in[length:]
	ok = true
	return
***REMOVED***

var (
	comma         = []byte***REMOVED***','***REMOVED***
	emptyNameList = []string***REMOVED******REMOVED***
)

func parseNameList(in []byte) (out []string, rest []byte, ok bool) ***REMOVED***
	contents, rest, ok := parseString(in)
	if !ok ***REMOVED***
		return
	***REMOVED***
	if len(contents) == 0 ***REMOVED***
		out = emptyNameList
		return
	***REMOVED***
	parts := bytes.Split(contents, comma)
	out = make([]string, len(parts))
	for i, part := range parts ***REMOVED***
		out[i] = string(part)
	***REMOVED***
	return
***REMOVED***

func parseInt(in []byte) (out *big.Int, rest []byte, ok bool) ***REMOVED***
	contents, rest, ok := parseString(in)
	if !ok ***REMOVED***
		return
	***REMOVED***
	out = new(big.Int)

	if len(contents) > 0 && contents[0]&0x80 == 0x80 ***REMOVED***
		// This is a negative number
		notBytes := make([]byte, len(contents))
		for i := range notBytes ***REMOVED***
			notBytes[i] = ^contents[i]
		***REMOVED***
		out.SetBytes(notBytes)
		out.Add(out, bigOne)
		out.Neg(out)
	***REMOVED*** else ***REMOVED***
		// Positive number
		out.SetBytes(contents)
	***REMOVED***
	ok = true
	return
***REMOVED***

func parseUint32(in []byte) (uint32, []byte, bool) ***REMOVED***
	if len(in) < 4 ***REMOVED***
		return 0, nil, false
	***REMOVED***
	return binary.BigEndian.Uint32(in), in[4:], true
***REMOVED***

func parseUint64(in []byte) (uint64, []byte, bool) ***REMOVED***
	if len(in) < 8 ***REMOVED***
		return 0, nil, false
	***REMOVED***
	return binary.BigEndian.Uint64(in), in[8:], true
***REMOVED***

func intLength(n *big.Int) int ***REMOVED***
	length := 4 /* length bytes */
	if n.Sign() < 0 ***REMOVED***
		nMinus1 := new(big.Int).Neg(n)
		nMinus1.Sub(nMinus1, bigOne)
		bitLen := nMinus1.BitLen()
		if bitLen%8 == 0 ***REMOVED***
			// The number will need 0xff padding
			length++
		***REMOVED***
		length += (bitLen + 7) / 8
	***REMOVED*** else if n.Sign() == 0 ***REMOVED***
		// A zero is the zero length string
	***REMOVED*** else ***REMOVED***
		bitLen := n.BitLen()
		if bitLen%8 == 0 ***REMOVED***
			// The number will need 0x00 padding
			length++
		***REMOVED***
		length += (bitLen + 7) / 8
	***REMOVED***

	return length
***REMOVED***

func marshalUint32(to []byte, n uint32) []byte ***REMOVED***
	binary.BigEndian.PutUint32(to, n)
	return to[4:]
***REMOVED***

func marshalUint64(to []byte, n uint64) []byte ***REMOVED***
	binary.BigEndian.PutUint64(to, n)
	return to[8:]
***REMOVED***

func marshalInt(to []byte, n *big.Int) []byte ***REMOVED***
	lengthBytes := to
	to = to[4:]
	length := 0

	if n.Sign() < 0 ***REMOVED***
		// A negative number has to be converted to two's-complement
		// form. So we'll subtract 1 and invert. If the
		// most-significant-bit isn't set then we'll need to pad the
		// beginning with 0xff in order to keep the number negative.
		nMinus1 := new(big.Int).Neg(n)
		nMinus1.Sub(nMinus1, bigOne)
		bytes := nMinus1.Bytes()
		for i := range bytes ***REMOVED***
			bytes[i] ^= 0xff
		***REMOVED***
		if len(bytes) == 0 || bytes[0]&0x80 == 0 ***REMOVED***
			to[0] = 0xff
			to = to[1:]
			length++
		***REMOVED***
		nBytes := copy(to, bytes)
		to = to[nBytes:]
		length += nBytes
	***REMOVED*** else if n.Sign() == 0 ***REMOVED***
		// A zero is the zero length string
	***REMOVED*** else ***REMOVED***
		bytes := n.Bytes()
		if len(bytes) > 0 && bytes[0]&0x80 != 0 ***REMOVED***
			// We'll have to pad this with a 0x00 in order to
			// stop it looking like a negative number.
			to[0] = 0
			to = to[1:]
			length++
		***REMOVED***
		nBytes := copy(to, bytes)
		to = to[nBytes:]
		length += nBytes
	***REMOVED***

	lengthBytes[0] = byte(length >> 24)
	lengthBytes[1] = byte(length >> 16)
	lengthBytes[2] = byte(length >> 8)
	lengthBytes[3] = byte(length)
	return to
***REMOVED***

func writeInt(w io.Writer, n *big.Int) ***REMOVED***
	length := intLength(n)
	buf := make([]byte, length)
	marshalInt(buf, n)
	w.Write(buf)
***REMOVED***

func writeString(w io.Writer, s []byte) ***REMOVED***
	var lengthBytes [4]byte
	lengthBytes[0] = byte(len(s) >> 24)
	lengthBytes[1] = byte(len(s) >> 16)
	lengthBytes[2] = byte(len(s) >> 8)
	lengthBytes[3] = byte(len(s))
	w.Write(lengthBytes[:])
	w.Write(s)
***REMOVED***

func stringLength(n int) int ***REMOVED***
	return 4 + n
***REMOVED***

func marshalString(to []byte, s []byte) []byte ***REMOVED***
	to[0] = byte(len(s) >> 24)
	to[1] = byte(len(s) >> 16)
	to[2] = byte(len(s) >> 8)
	to[3] = byte(len(s))
	to = to[4:]
	copy(to, s)
	return to[len(s):]
***REMOVED***

var bigIntType = reflect.TypeOf((*big.Int)(nil))

// Decode a packet into its corresponding message.
func decode(packet []byte) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	var msg interface***REMOVED******REMOVED***
	switch packet[0] ***REMOVED***
	case msgDisconnect:
		msg = new(disconnectMsg)
	case msgServiceRequest:
		msg = new(serviceRequestMsg)
	case msgServiceAccept:
		msg = new(serviceAcceptMsg)
	case msgKexInit:
		msg = new(kexInitMsg)
	case msgKexDHInit:
		msg = new(kexDHInitMsg)
	case msgKexDHReply:
		msg = new(kexDHReplyMsg)
	case msgUserAuthRequest:
		msg = new(userAuthRequestMsg)
	case msgUserAuthSuccess:
		return new(userAuthSuccessMsg), nil
	case msgUserAuthFailure:
		msg = new(userAuthFailureMsg)
	case msgUserAuthPubKeyOk:
		msg = new(userAuthPubKeyOkMsg)
	case msgGlobalRequest:
		msg = new(globalRequestMsg)
	case msgRequestSuccess:
		msg = new(globalRequestSuccessMsg)
	case msgRequestFailure:
		msg = new(globalRequestFailureMsg)
	case msgChannelOpen:
		msg = new(channelOpenMsg)
	case msgChannelData:
		msg = new(channelDataMsg)
	case msgChannelOpenConfirm:
		msg = new(channelOpenConfirmMsg)
	case msgChannelOpenFailure:
		msg = new(channelOpenFailureMsg)
	case msgChannelWindowAdjust:
		msg = new(windowAdjustMsg)
	case msgChannelEOF:
		msg = new(channelEOFMsg)
	case msgChannelClose:
		msg = new(channelCloseMsg)
	case msgChannelRequest:
		msg = new(channelRequestMsg)
	case msgChannelSuccess:
		msg = new(channelRequestSuccessMsg)
	case msgChannelFailure:
		msg = new(channelRequestFailureMsg)
	default:
		return nil, unexpectedMessageError(0, packet[0])
	***REMOVED***
	if err := Unmarshal(packet, msg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return msg, nil
***REMOVED***
