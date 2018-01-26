package networkdb

import "github.com/gogo/protobuf/proto"

const (
	// Compound message header overhead 1 byte(message type) + 4
	// bytes (num messages)
	compoundHeaderOverhead = 5

	// Overhead for each embedded message in a compound message 4
	// bytes (len of embedded message)
	compoundOverhead = 4
)

func encodeRawMessage(t MessageType, raw []byte) ([]byte, error) ***REMOVED***
	gMsg := GossipMessage***REMOVED***
		Type: t,
		Data: raw,
	***REMOVED***

	buf, err := proto.Marshal(&gMsg)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return buf, nil
***REMOVED***

func encodeMessage(t MessageType, msg interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	buf, err := proto.Marshal(msg.(proto.Message))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	buf, err = encodeRawMessage(t, buf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return buf, nil
***REMOVED***

func decodeMessage(buf []byte) (MessageType, []byte, error) ***REMOVED***
	var gMsg GossipMessage

	err := proto.Unmarshal(buf, &gMsg)
	if err != nil ***REMOVED***
		return MessageTypeInvalid, nil, err
	***REMOVED***

	return gMsg.Type, gMsg.Data, nil
***REMOVED***

// makeCompoundMessage takes a list of messages and generates
// a single compound message containing all of them
func makeCompoundMessage(msgs [][]byte) []byte ***REMOVED***
	cMsg := CompoundMessage***REMOVED******REMOVED***

	cMsg.Messages = make([]*CompoundMessage_SimpleMessage, 0, len(msgs))
	for _, m := range msgs ***REMOVED***
		cMsg.Messages = append(cMsg.Messages, &CompoundMessage_SimpleMessage***REMOVED***
			Payload: m,
		***REMOVED***)
	***REMOVED***

	buf, err := proto.Marshal(&cMsg)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	gMsg := GossipMessage***REMOVED***
		Type: MessageTypeCompound,
		Data: buf,
	***REMOVED***

	buf, err = proto.Marshal(&gMsg)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	return buf
***REMOVED***

// decodeCompoundMessage splits a compound message and returns
// the slices of individual messages. Returns any potential error.
func decodeCompoundMessage(buf []byte) ([][]byte, error) ***REMOVED***
	var cMsg CompoundMessage
	if err := proto.Unmarshal(buf, &cMsg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	parts := make([][]byte, 0, len(cMsg.Messages))
	for _, m := range cMsg.Messages ***REMOVED***
		parts = append(parts, m.Payload)
	***REMOVED***

	return parts, nil
***REMOVED***
