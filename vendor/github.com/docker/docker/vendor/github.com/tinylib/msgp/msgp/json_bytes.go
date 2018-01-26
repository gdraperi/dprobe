package msgp

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"time"
)

var unfuns [_maxtype]func(jsWriter, []byte, []byte) ([]byte, []byte, error)

func init() ***REMOVED***

	// NOTE(pmh): this is best expressed as a jump table,
	// but gc doesn't do that yet. revisit post-go1.5.
	unfuns = [_maxtype]func(jsWriter, []byte, []byte) ([]byte, []byte, error)***REMOVED***
		StrType:        rwStringBytes,
		BinType:        rwBytesBytes,
		MapType:        rwMapBytes,
		ArrayType:      rwArrayBytes,
		Float64Type:    rwFloat64Bytes,
		Float32Type:    rwFloat32Bytes,
		BoolType:       rwBoolBytes,
		IntType:        rwIntBytes,
		UintType:       rwUintBytes,
		NilType:        rwNullBytes,
		ExtensionType:  rwExtensionBytes,
		Complex64Type:  rwExtensionBytes,
		Complex128Type: rwExtensionBytes,
		TimeType:       rwTimeBytes,
	***REMOVED***
***REMOVED***

// UnmarshalAsJSON takes raw messagepack and writes
// it as JSON to 'w'. If an error is returned, the
// bytes not translated will also be returned. If
// no errors are encountered, the length of the returned
// slice will be zero.
func UnmarshalAsJSON(w io.Writer, msg []byte) ([]byte, error) ***REMOVED***
	var (
		scratch []byte
		cast    bool
		dst     jsWriter
		err     error
	)
	if jsw, ok := w.(jsWriter); ok ***REMOVED***
		dst = jsw
		cast = true
	***REMOVED*** else ***REMOVED***
		dst = bufio.NewWriterSize(w, 512)
	***REMOVED***
	for len(msg) > 0 && err == nil ***REMOVED***
		msg, scratch, err = writeNext(dst, msg, scratch)
	***REMOVED***
	if !cast && err == nil ***REMOVED***
		err = dst.(*bufio.Writer).Flush()
	***REMOVED***
	return msg, err
***REMOVED***

func writeNext(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	if len(msg) < 1 ***REMOVED***
		return msg, scratch, ErrShortBytes
	***REMOVED***
	t := getType(msg[0])
	if t == InvalidType ***REMOVED***
		return msg, scratch, InvalidPrefixError(msg[0])
	***REMOVED***
	if t == ExtensionType ***REMOVED***
		et, err := peekExtension(msg)
		if err != nil ***REMOVED***
			return nil, scratch, err
		***REMOVED***
		if et == TimeExtension ***REMOVED***
			t = TimeType
		***REMOVED***
	***REMOVED***
	return unfuns[t](w, msg, scratch)
***REMOVED***

func rwArrayBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	sz, msg, err := ReadArrayHeaderBytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	err = w.WriteByte('[')
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	for i := uint32(0); i < sz; i++ ***REMOVED***
		if i != 0 ***REMOVED***
			err = w.WriteByte(',')
			if err != nil ***REMOVED***
				return msg, scratch, err
			***REMOVED***
		***REMOVED***
		msg, scratch, err = writeNext(w, msg, scratch)
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
	***REMOVED***
	err = w.WriteByte(']')
	return msg, scratch, err
***REMOVED***

func rwMapBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	sz, msg, err := ReadMapHeaderBytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	err = w.WriteByte('***REMOVED***')
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	for i := uint32(0); i < sz; i++ ***REMOVED***
		if i != 0 ***REMOVED***
			err = w.WriteByte(',')
			if err != nil ***REMOVED***
				return msg, scratch, err
			***REMOVED***
		***REMOVED***
		msg, scratch, err = rwMapKeyBytes(w, msg, scratch)
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
		err = w.WriteByte(':')
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
		msg, scratch, err = writeNext(w, msg, scratch)
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
	***REMOVED***
	err = w.WriteByte('***REMOVED***')
	return msg, scratch, err
***REMOVED***

func rwMapKeyBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	msg, scratch, err := rwStringBytes(w, msg, scratch)
	if err != nil ***REMOVED***
		if tperr, ok := err.(TypeError); ok && tperr.Encoded == BinType ***REMOVED***
			return rwBytesBytes(w, msg, scratch)
		***REMOVED***
	***REMOVED***
	return msg, scratch, err
***REMOVED***

func rwStringBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	str, msg, err := ReadStringZC(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	_, err = rwquoted(w, str)
	return msg, scratch, err
***REMOVED***

func rwBytesBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	bts, msg, err := ReadBytesZC(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	l := base64.StdEncoding.EncodedLen(len(bts))
	if cap(scratch) >= l ***REMOVED***
		scratch = scratch[0:l]
	***REMOVED*** else ***REMOVED***
		scratch = make([]byte, l)
	***REMOVED***
	base64.StdEncoding.Encode(scratch, bts)
	err = w.WriteByte('"')
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	_, err = w.Write(scratch)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	err = w.WriteByte('"')
	return msg, scratch, err
***REMOVED***

func rwNullBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	msg, err := ReadNilBytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	_, err = w.Write(null)
	return msg, scratch, err
***REMOVED***

func rwBoolBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	b, msg, err := ReadBoolBytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	if b ***REMOVED***
		_, err = w.WriteString("true")
		return msg, scratch, err
	***REMOVED***
	_, err = w.WriteString("false")
	return msg, scratch, err
***REMOVED***

func rwIntBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	i, msg, err := ReadInt64Bytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	scratch = strconv.AppendInt(scratch[0:0], i, 10)
	_, err = w.Write(scratch)
	return msg, scratch, err
***REMOVED***

func rwUintBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	u, msg, err := ReadUint64Bytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	scratch = strconv.AppendUint(scratch[0:0], u, 10)
	_, err = w.Write(scratch)
	return msg, scratch, err
***REMOVED***

func rwFloatBytes(w jsWriter, msg []byte, f64 bool, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	var f float64
	var err error
	var sz int
	if f64 ***REMOVED***
		sz = 64
		f, msg, err = ReadFloat64Bytes(msg)
	***REMOVED*** else ***REMOVED***
		sz = 32
		var v float32
		v, msg, err = ReadFloat32Bytes(msg)
		f = float64(v)
	***REMOVED***
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	scratch = strconv.AppendFloat(scratch, f, 'f', -1, sz)
	_, err = w.Write(scratch)
	return msg, scratch, err
***REMOVED***

func rwFloat32Bytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	var f float32
	var err error
	f, msg, err = ReadFloat32Bytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	scratch = strconv.AppendFloat(scratch[:0], float64(f), 'f', -1, 32)
	_, err = w.Write(scratch)
	return msg, scratch, err
***REMOVED***

func rwFloat64Bytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	var f float64
	var err error
	f, msg, err = ReadFloat64Bytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	scratch = strconv.AppendFloat(scratch[:0], f, 'f', -1, 64)
	_, err = w.Write(scratch)
	return msg, scratch, err
***REMOVED***

func rwTimeBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	var t time.Time
	var err error
	t, msg, err = ReadTimeBytes(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	bts, err := t.MarshalJSON()
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	_, err = w.Write(bts)
	return msg, scratch, err
***REMOVED***

func rwExtensionBytes(w jsWriter, msg []byte, scratch []byte) ([]byte, []byte, error) ***REMOVED***
	var err error
	var et int8
	et, err = peekExtension(msg)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***

	// if it's time.Time
	if et == TimeExtension ***REMOVED***
		var tm time.Time
		tm, msg, err = ReadTimeBytes(msg)
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
		bts, err := tm.MarshalJSON()
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
		_, err = w.Write(bts)
		return msg, scratch, err
	***REMOVED***

	// if the extension is registered,
	// use its canonical JSON form
	if f, ok := extensionReg[et]; ok ***REMOVED***
		e := f()
		msg, err = ReadExtensionBytes(msg, e)
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
		bts, err := json.Marshal(e)
		if err != nil ***REMOVED***
			return msg, scratch, err
		***REMOVED***
		_, err = w.Write(bts)
		return msg, scratch, err
	***REMOVED***

	// otherwise, write `***REMOVED***"type": <num>, "data": "<base64data>"***REMOVED***`
	r := RawExtension***REMOVED******REMOVED***
	r.Type = et
	msg, err = ReadExtensionBytes(msg, &r)
	if err != nil ***REMOVED***
		return msg, scratch, err
	***REMOVED***
	scratch, err = writeExt(w, r, scratch)
	return msg, scratch, err
***REMOVED***

func writeExt(w jsWriter, r RawExtension, scratch []byte) ([]byte, error) ***REMOVED***
	_, err := w.WriteString(`***REMOVED***"type":`)
	if err != nil ***REMOVED***
		return scratch, err
	***REMOVED***
	scratch = strconv.AppendInt(scratch[0:0], int64(r.Type), 10)
	_, err = w.Write(scratch)
	if err != nil ***REMOVED***
		return scratch, err
	***REMOVED***
	_, err = w.WriteString(`,"data":"`)
	if err != nil ***REMOVED***
		return scratch, err
	***REMOVED***
	l := base64.StdEncoding.EncodedLen(len(r.Data))
	if cap(scratch) >= l ***REMOVED***
		scratch = scratch[0:l]
	***REMOVED*** else ***REMOVED***
		scratch = make([]byte, l)
	***REMOVED***
	base64.StdEncoding.Encode(scratch, r.Data)
	_, err = w.Write(scratch)
	if err != nil ***REMOVED***
		return scratch, err
	***REMOVED***
	_, err = w.WriteString(`"***REMOVED***`)
	return scratch, err
***REMOVED***
