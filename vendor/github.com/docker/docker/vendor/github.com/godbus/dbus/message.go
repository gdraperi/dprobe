package dbus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strconv"
)

const protoVersion byte = 1

// Flags represents the possible flags of a D-Bus message.
type Flags byte

const (
	// FlagNoReplyExpected signals that the message is not expected to generate
	// a reply. If this flag is set on outgoing messages, any possible reply
	// will be discarded.
	FlagNoReplyExpected Flags = 1 << iota
	// FlagNoAutoStart signals that the message bus should not automatically
	// start an application when handling this message.
	FlagNoAutoStart
	// FlagAllowInteractiveAuthorization may be set on a method call
	// message to inform the receiving side that the caller is prepared
	// to wait for interactive authorization, which might take a
	// considerable time to complete. For instance, if this flag is set,
	// it would be appropriate to query the user for passwords or
	// confirmation via Polkit or a similar framework.
	FlagAllowInteractiveAuthorization
)

// Type represents the possible types of a D-Bus message.
type Type byte

const (
	TypeMethodCall Type = 1 + iota
	TypeMethodReply
	TypeError
	TypeSignal
	typeMax
)

func (t Type) String() string ***REMOVED***
	switch t ***REMOVED***
	case TypeMethodCall:
		return "method call"
	case TypeMethodReply:
		return "reply"
	case TypeError:
		return "error"
	case TypeSignal:
		return "signal"
	***REMOVED***
	return "invalid"
***REMOVED***

// HeaderField represents the possible byte codes for the headers
// of a D-Bus message.
type HeaderField byte

const (
	FieldPath HeaderField = 1 + iota
	FieldInterface
	FieldMember
	FieldErrorName
	FieldReplySerial
	FieldDestination
	FieldSender
	FieldSignature
	FieldUnixFDs
	fieldMax
)

// An InvalidMessageError describes the reason why a D-Bus message is regarded as
// invalid.
type InvalidMessageError string

func (e InvalidMessageError) Error() string ***REMOVED***
	return "dbus: invalid message: " + string(e)
***REMOVED***

// fieldType are the types of the various header fields.
var fieldTypes = [fieldMax]reflect.Type***REMOVED***
	FieldPath:        objectPathType,
	FieldInterface:   stringType,
	FieldMember:      stringType,
	FieldErrorName:   stringType,
	FieldReplySerial: uint32Type,
	FieldDestination: stringType,
	FieldSender:      stringType,
	FieldSignature:   signatureType,
	FieldUnixFDs:     uint32Type,
***REMOVED***

// requiredFields lists the header fields that are required by the different
// message types.
var requiredFields = [typeMax][]HeaderField***REMOVED***
	TypeMethodCall:  ***REMOVED***FieldPath, FieldMember***REMOVED***,
	TypeMethodReply: ***REMOVED***FieldReplySerial***REMOVED***,
	TypeError:       ***REMOVED***FieldErrorName, FieldReplySerial***REMOVED***,
	TypeSignal:      ***REMOVED***FieldPath, FieldInterface, FieldMember***REMOVED***,
***REMOVED***

// Message represents a single D-Bus message.
type Message struct ***REMOVED***
	Type
	Flags
	Headers map[HeaderField]Variant
	Body    []interface***REMOVED******REMOVED***

	serial uint32
***REMOVED***

type header struct ***REMOVED***
	Field byte
	Variant
***REMOVED***

// DecodeMessage tries to decode a single message in the D-Bus wire format
// from the given reader. The byte order is figured out from the first byte.
// The possibly returned error can be an error of the underlying reader, an
// InvalidMessageError or a FormatError.
func DecodeMessage(rd io.Reader) (msg *Message, err error) ***REMOVED***
	var order binary.ByteOrder
	var hlength, length uint32
	var typ, flags, proto byte
	var headers []header

	b := make([]byte, 1)
	_, err = rd.Read(b)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	switch b[0] ***REMOVED***
	case 'l':
		order = binary.LittleEndian
	case 'B':
		order = binary.BigEndian
	default:
		return nil, InvalidMessageError("invalid byte order")
	***REMOVED***

	dec := newDecoder(rd, order)
	dec.pos = 1

	msg = new(Message)
	vs, err := dec.Decode(Signature***REMOVED***"yyyuu"***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = Store(vs, &typ, &flags, &proto, &length, &msg.serial); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	msg.Type = Type(typ)
	msg.Flags = Flags(flags)

	// get the header length separately because we need it later
	b = make([]byte, 4)
	_, err = io.ReadFull(rd, b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	binary.Read(bytes.NewBuffer(b), order, &hlength)
	if hlength+length+16 > 1<<27 ***REMOVED***
		return nil, InvalidMessageError("message is too long")
	***REMOVED***
	dec = newDecoder(io.MultiReader(bytes.NewBuffer(b), rd), order)
	dec.pos = 12
	vs, err = dec.Decode(Signature***REMOVED***"a(yv)"***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = Store(vs, &headers); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	msg.Headers = make(map[HeaderField]Variant)
	for _, v := range headers ***REMOVED***
		msg.Headers[HeaderField(v.Field)] = v.Variant
	***REMOVED***

	dec.align(8)
	body := make([]byte, int(length))
	if length != 0 ***REMOVED***
		_, err := io.ReadFull(rd, body)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if err = msg.IsValid(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	sig, _ := msg.Headers[FieldSignature].value.(Signature)
	if sig.str != "" ***REMOVED***
		buf := bytes.NewBuffer(body)
		dec = newDecoder(buf, order)
		vs, err := dec.Decode(sig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		msg.Body = vs
	***REMOVED***

	return
***REMOVED***

// EncodeTo encodes and sends a message to the given writer. The byte order must
// be either binary.LittleEndian or binary.BigEndian. If the message is not
// valid or an error occurs when writing, an error is returned.
func (msg *Message) EncodeTo(out io.Writer, order binary.ByteOrder) error ***REMOVED***
	if err := msg.IsValid(); err != nil ***REMOVED***
		return err
	***REMOVED***
	var vs [7]interface***REMOVED******REMOVED***
	switch order ***REMOVED***
	case binary.LittleEndian:
		vs[0] = byte('l')
	case binary.BigEndian:
		vs[0] = byte('B')
	default:
		return errors.New("dbus: invalid byte order")
	***REMOVED***
	body := new(bytes.Buffer)
	enc := newEncoder(body, order)
	if len(msg.Body) != 0 ***REMOVED***
		enc.Encode(msg.Body...)
	***REMOVED***
	vs[1] = msg.Type
	vs[2] = msg.Flags
	vs[3] = protoVersion
	vs[4] = uint32(len(body.Bytes()))
	vs[5] = msg.serial
	headers := make([]header, 0, len(msg.Headers))
	for k, v := range msg.Headers ***REMOVED***
		headers = append(headers, header***REMOVED***byte(k), v***REMOVED***)
	***REMOVED***
	vs[6] = headers
	var buf bytes.Buffer
	enc = newEncoder(&buf, order)
	enc.Encode(vs[:]...)
	enc.align(8)
	body.WriteTo(&buf)
	if buf.Len() > 1<<27 ***REMOVED***
		return InvalidMessageError("message is too long")
	***REMOVED***
	if _, err := buf.WriteTo(out); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// IsValid checks whether msg is a valid message and returns an
// InvalidMessageError if it is not.
func (msg *Message) IsValid() error ***REMOVED***
	if msg.Flags & ^(FlagNoAutoStart|FlagNoReplyExpected|FlagAllowInteractiveAuthorization) != 0 ***REMOVED***
		return InvalidMessageError("invalid flags")
	***REMOVED***
	if msg.Type == 0 || msg.Type >= typeMax ***REMOVED***
		return InvalidMessageError("invalid message type")
	***REMOVED***
	for k, v := range msg.Headers ***REMOVED***
		if k == 0 || k >= fieldMax ***REMOVED***
			return InvalidMessageError("invalid header")
		***REMOVED***
		if reflect.TypeOf(v.value) != fieldTypes[k] ***REMOVED***
			return InvalidMessageError("invalid type of header field")
		***REMOVED***
	***REMOVED***
	for _, v := range requiredFields[msg.Type] ***REMOVED***
		if _, ok := msg.Headers[v]; !ok ***REMOVED***
			return InvalidMessageError("missing required header")
		***REMOVED***
	***REMOVED***
	if path, ok := msg.Headers[FieldPath]; ok ***REMOVED***
		if !path.value.(ObjectPath).IsValid() ***REMOVED***
			return InvalidMessageError("invalid path name")
		***REMOVED***
	***REMOVED***
	if iface, ok := msg.Headers[FieldInterface]; ok ***REMOVED***
		if !isValidInterface(iface.value.(string)) ***REMOVED***
			return InvalidMessageError("invalid interface name")
		***REMOVED***
	***REMOVED***
	if member, ok := msg.Headers[FieldMember]; ok ***REMOVED***
		if !isValidMember(member.value.(string)) ***REMOVED***
			return InvalidMessageError("invalid member name")
		***REMOVED***
	***REMOVED***
	if errname, ok := msg.Headers[FieldErrorName]; ok ***REMOVED***
		if !isValidInterface(errname.value.(string)) ***REMOVED***
			return InvalidMessageError("invalid error name")
		***REMOVED***
	***REMOVED***
	if len(msg.Body) != 0 ***REMOVED***
		if _, ok := msg.Headers[FieldSignature]; !ok ***REMOVED***
			return InvalidMessageError("missing signature")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Serial returns the message's serial number. The returned value is only valid
// for messages received by eavesdropping.
func (msg *Message) Serial() uint32 ***REMOVED***
	return msg.serial
***REMOVED***

// String returns a string representation of a message similar to the format of
// dbus-monitor.
func (msg *Message) String() string ***REMOVED***
	if err := msg.IsValid(); err != nil ***REMOVED***
		return "<invalid>"
	***REMOVED***
	s := msg.Type.String()
	if v, ok := msg.Headers[FieldSender]; ok ***REMOVED***
		s += " from " + v.value.(string)
	***REMOVED***
	if v, ok := msg.Headers[FieldDestination]; ok ***REMOVED***
		s += " to " + v.value.(string)
	***REMOVED***
	s += " serial " + strconv.FormatUint(uint64(msg.serial), 10)
	if v, ok := msg.Headers[FieldReplySerial]; ok ***REMOVED***
		s += " reply_serial " + strconv.FormatUint(uint64(v.value.(uint32)), 10)
	***REMOVED***
	if v, ok := msg.Headers[FieldUnixFDs]; ok ***REMOVED***
		s += " unixfds " + strconv.FormatUint(uint64(v.value.(uint32)), 10)
	***REMOVED***
	if v, ok := msg.Headers[FieldPath]; ok ***REMOVED***
		s += " path " + string(v.value.(ObjectPath))
	***REMOVED***
	if v, ok := msg.Headers[FieldInterface]; ok ***REMOVED***
		s += " interface " + v.value.(string)
	***REMOVED***
	if v, ok := msg.Headers[FieldErrorName]; ok ***REMOVED***
		s += " error " + v.value.(string)
	***REMOVED***
	if v, ok := msg.Headers[FieldMember]; ok ***REMOVED***
		s += " member " + v.value.(string)
	***REMOVED***
	if len(msg.Body) != 0 ***REMOVED***
		s += "\n"
	***REMOVED***
	for i, v := range msg.Body ***REMOVED***
		s += "  " + MakeVariant(v).String()
		if i != len(msg.Body)-1 ***REMOVED***
			s += "\n"
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***
