package fluent

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Entry) DecodeMsg(dc *msgp.Reader) (err error) ***REMOVED***
	var zxvk uint32
	zxvk, err = dc.ReadArrayHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zxvk != 2 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 2, Got: zxvk***REMOVED***
		return
	***REMOVED***
	z.Time, err = dc.ReadInt64()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Record, err = dc.ReadIntf()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// EncodeMsg implements msgp.Encodable
func (z Entry) EncodeMsg(en *msgp.Writer) (err error) ***REMOVED***
	// array header, size 2
	err = en.Append(0x92)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = en.WriteInt64(z.Time)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteIntf(z.Record)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// MarshalMsg implements msgp.Marshaler
func (z Entry) MarshalMsg(b []byte) (o []byte, err error) ***REMOVED***
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendInt64(o, z.Time)
	o, err = msgp.AppendIntf(o, z.Record)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Entry) UnmarshalMsg(bts []byte) (o []byte, err error) ***REMOVED***
	var zbzg uint32
	zbzg, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zbzg != 2 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 2, Got: zbzg***REMOVED***
		return
	***REMOVED***
	z.Time, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Record, bts, err = msgp.ReadIntfBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	o = bts
	return
***REMOVED***

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Entry) Msgsize() (s int) ***REMOVED***
	s = 1 + msgp.Int64Size + msgp.GuessSize(z.Record)
	return
***REMOVED***

// DecodeMsg implements msgp.Decodable
func (z *Forward) DecodeMsg(dc *msgp.Reader) (err error) ***REMOVED***
	var zcmr uint32
	zcmr, err = dc.ReadArrayHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zcmr != 3 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 3, Got: zcmr***REMOVED***
		return
	***REMOVED***
	z.Tag, err = dc.ReadString()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var zajw uint32
	zajw, err = dc.ReadArrayHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if cap(z.Entries) >= int(zajw) ***REMOVED***
		z.Entries = (z.Entries)[:zajw]
	***REMOVED*** else ***REMOVED***
		z.Entries = make([]Entry, zajw)
	***REMOVED***
	for zbai := range z.Entries ***REMOVED***
		var zwht uint32
		zwht, err = dc.ReadArrayHeader()
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if zwht != 2 ***REMOVED***
			err = msgp.ArrayError***REMOVED***Wanted: 2, Got: zwht***REMOVED***
			return
		***REMOVED***
		z.Entries[zbai].Time, err = dc.ReadInt64()
		if err != nil ***REMOVED***
			return
		***REMOVED***
		z.Entries[zbai].Record, err = dc.ReadIntf()
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	z.Option, err = dc.ReadIntf()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// EncodeMsg implements msgp.Encodable
func (z *Forward) EncodeMsg(en *msgp.Writer) (err error) ***REMOVED***
	// array header, size 3
	err = en.Append(0x93)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = en.WriteString(z.Tag)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteArrayHeader(uint32(len(z.Entries)))
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for zbai := range z.Entries ***REMOVED***
		// array header, size 2
		err = en.Append(0x92)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = en.WriteInt64(z.Entries[zbai].Time)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		err = en.WriteIntf(z.Entries[zbai].Record)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	err = en.WriteIntf(z.Option)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// MarshalMsg implements msgp.Marshaler
func (z *Forward) MarshalMsg(b []byte) (o []byte, err error) ***REMOVED***
	o = msgp.Require(b, z.Msgsize())
	// array header, size 3
	o = append(o, 0x93)
	o = msgp.AppendString(o, z.Tag)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Entries)))
	for zbai := range z.Entries ***REMOVED***
		// array header, size 2
		o = append(o, 0x92)
		o = msgp.AppendInt64(o, z.Entries[zbai].Time)
		o, err = msgp.AppendIntf(o, z.Entries[zbai].Record)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	o, err = msgp.AppendIntf(o, z.Option)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Forward) UnmarshalMsg(bts []byte) (o []byte, err error) ***REMOVED***
	var zhct uint32
	zhct, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zhct != 3 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 3, Got: zhct***REMOVED***
		return
	***REMOVED***
	z.Tag, bts, err = msgp.ReadStringBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var zcua uint32
	zcua, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if cap(z.Entries) >= int(zcua) ***REMOVED***
		z.Entries = (z.Entries)[:zcua]
	***REMOVED*** else ***REMOVED***
		z.Entries = make([]Entry, zcua)
	***REMOVED***
	for zbai := range z.Entries ***REMOVED***
		var zxhx uint32
		zxhx, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if zxhx != 2 ***REMOVED***
			err = msgp.ArrayError***REMOVED***Wanted: 2, Got: zxhx***REMOVED***
			return
		***REMOVED***
		z.Entries[zbai].Time, bts, err = msgp.ReadInt64Bytes(bts)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		z.Entries[zbai].Record, bts, err = msgp.ReadIntfBytes(bts)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	z.Option, bts, err = msgp.ReadIntfBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	o = bts
	return
***REMOVED***

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Forward) Msgsize() (s int) ***REMOVED***
	s = 1 + msgp.StringPrefixSize + len(z.Tag) + msgp.ArrayHeaderSize
	for zbai := range z.Entries ***REMOVED***
		s += 1 + msgp.Int64Size + msgp.GuessSize(z.Entries[zbai].Record)
	***REMOVED***
	s += msgp.GuessSize(z.Option)
	return
***REMOVED***

// DecodeMsg implements msgp.Decodable
func (z *Message) DecodeMsg(dc *msgp.Reader) (err error) ***REMOVED***
	var zlqf uint32
	zlqf, err = dc.ReadArrayHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zlqf != 4 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 4, Got: zlqf***REMOVED***
		return
	***REMOVED***
	z.Tag, err = dc.ReadString()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Time, err = dc.ReadInt64()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Record, err = dc.ReadIntf()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Option, err = dc.ReadIntf()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// EncodeMsg implements msgp.Encodable
func (z *Message) EncodeMsg(en *msgp.Writer) (err error) ***REMOVED***
	// array header, size 4
	err = en.Append(0x94)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = en.WriteString(z.Tag)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteInt64(z.Time)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteIntf(z.Record)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteIntf(z.Option)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// MarshalMsg implements msgp.Marshaler
func (z *Message) MarshalMsg(b []byte) (o []byte, err error) ***REMOVED***
	o = msgp.Require(b, z.Msgsize())
	// array header, size 4
	o = append(o, 0x94)
	o = msgp.AppendString(o, z.Tag)
	o = msgp.AppendInt64(o, z.Time)
	o, err = msgp.AppendIntf(o, z.Record)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	o, err = msgp.AppendIntf(o, z.Option)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Message) UnmarshalMsg(bts []byte) (o []byte, err error) ***REMOVED***
	var zdaf uint32
	zdaf, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zdaf != 4 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 4, Got: zdaf***REMOVED***
		return
	***REMOVED***
	z.Tag, bts, err = msgp.ReadStringBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Time, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Record, bts, err = msgp.ReadIntfBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Option, bts, err = msgp.ReadIntfBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	o = bts
	return
***REMOVED***

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Message) Msgsize() (s int) ***REMOVED***
	s = 1 + msgp.StringPrefixSize + len(z.Tag) + msgp.Int64Size + msgp.GuessSize(z.Record) + msgp.GuessSize(z.Option)
	return
***REMOVED***

// DecodeMsg implements msgp.Decodable
func (z *MessageExt) DecodeMsg(dc *msgp.Reader) (err error) ***REMOVED***
	var zpks uint32
	zpks, err = dc.ReadArrayHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zpks != 4 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 4, Got: zpks***REMOVED***
		return
	***REMOVED***
	z.Tag, err = dc.ReadString()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = dc.ReadExtension(&z.Time)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Record, err = dc.ReadIntf()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Option, err = dc.ReadIntf()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// EncodeMsg implements msgp.Encodable
func (z *MessageExt) EncodeMsg(en *msgp.Writer) (err error) ***REMOVED***
	// array header, size 4
	err = en.Append(0x94)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = en.WriteString(z.Tag)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteExtension(&z.Time)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteIntf(z.Record)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = en.WriteIntf(z.Option)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// MarshalMsg implements msgp.Marshaler
func (z *MessageExt) MarshalMsg(b []byte) (o []byte, err error) ***REMOVED***
	o = msgp.Require(b, z.Msgsize())
	// array header, size 4
	o = append(o, 0x94)
	o = msgp.AppendString(o, z.Tag)
	o, err = msgp.AppendExtension(o, &z.Time)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	o, err = msgp.AppendIntf(o, z.Record)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	o, err = msgp.AppendIntf(o, z.Option)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MessageExt) UnmarshalMsg(bts []byte) (o []byte, err error) ***REMOVED***
	var zjfb uint32
	zjfb, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if zjfb != 4 ***REMOVED***
		err = msgp.ArrayError***REMOVED***Wanted: 4, Got: zjfb***REMOVED***
		return
	***REMOVED***
	z.Tag, bts, err = msgp.ReadStringBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	bts, err = msgp.ReadExtensionBytes(bts, &z.Time)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Record, bts, err = msgp.ReadIntfBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	z.Option, bts, err = msgp.ReadIntfBytes(bts)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	o = bts
	return
***REMOVED***

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *MessageExt) Msgsize() (s int) ***REMOVED***
	s = 1 + msgp.StringPrefixSize + len(z.Tag) + msgp.ExtensionPrefixSize + z.Time.Len() + msgp.GuessSize(z.Record) + msgp.GuessSize(z.Option)
	return
***REMOVED***
