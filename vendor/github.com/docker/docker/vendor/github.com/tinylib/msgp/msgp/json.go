package msgp

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"unicode/utf8"
)

var (
	null = []byte("null")
	hex  = []byte("0123456789abcdef")
)

var defuns [_maxtype]func(jsWriter, *Reader) (int, error)

// note: there is an initialization loop if
// this isn't set up during init()
func init() ***REMOVED***
	// since none of these functions are inline-able,
	// there is not much of a penalty to the indirect
	// call. however, this is best expressed as a jump-table...
	defuns = [_maxtype]func(jsWriter, *Reader) (int, error)***REMOVED***
		StrType:        rwString,
		BinType:        rwBytes,
		MapType:        rwMap,
		ArrayType:      rwArray,
		Float64Type:    rwFloat64,
		Float32Type:    rwFloat32,
		BoolType:       rwBool,
		IntType:        rwInt,
		UintType:       rwUint,
		NilType:        rwNil,
		ExtensionType:  rwExtension,
		Complex64Type:  rwExtension,
		Complex128Type: rwExtension,
		TimeType:       rwTime,
	***REMOVED***
***REMOVED***

// this is the interface
// used to write json
type jsWriter interface ***REMOVED***
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
***REMOVED***

// CopyToJSON reads MessagePack from 'src' and copies it
// as JSON to 'dst' until EOF.
func CopyToJSON(dst io.Writer, src io.Reader) (n int64, err error) ***REMOVED***
	r := NewReader(src)
	n, err = r.WriteToJSON(dst)
	freeR(r)
	return
***REMOVED***

// WriteToJSON translates MessagePack from 'r' and writes it as
// JSON to 'w' until the underlying reader returns io.EOF. It returns
// the number of bytes written, and an error if it stopped before EOF.
func (r *Reader) WriteToJSON(w io.Writer) (n int64, err error) ***REMOVED***
	var j jsWriter
	var bf *bufio.Writer
	if jsw, ok := w.(jsWriter); ok ***REMOVED***
		j = jsw
	***REMOVED*** else ***REMOVED***
		bf = bufio.NewWriter(w)
		j = bf
	***REMOVED***
	var nn int
	for err == nil ***REMOVED***
		nn, err = rwNext(j, r)
		n += int64(nn)
	***REMOVED***
	if err != io.EOF ***REMOVED***
		if bf != nil ***REMOVED***
			bf.Flush()
		***REMOVED***
		return
	***REMOVED***
	err = nil
	if bf != nil ***REMOVED***
		err = bf.Flush()
	***REMOVED***
	return
***REMOVED***

func rwNext(w jsWriter, src *Reader) (int, error) ***REMOVED***
	t, err := src.NextType()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return defuns[t](w, src)
***REMOVED***

func rwMap(dst jsWriter, src *Reader) (n int, err error) ***REMOVED***
	var comma bool
	var sz uint32
	var field []byte

	sz, err = src.ReadMapHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if sz == 0 ***REMOVED***
		return dst.WriteString("***REMOVED******REMOVED***")
	***REMOVED***

	err = dst.WriteByte('***REMOVED***')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++
	var nn int
	for i := uint32(0); i < sz; i++ ***REMOVED***
		if comma ***REMOVED***
			err = dst.WriteByte(',')
			if err != nil ***REMOVED***
				return
			***REMOVED***
			n++
		***REMOVED***

		field, err = src.ReadMapKeyPtr()
		if err != nil ***REMOVED***
			return
		***REMOVED***
		nn, err = rwquoted(dst, field)
		n += nn
		if err != nil ***REMOVED***
			return
		***REMOVED***

		err = dst.WriteByte(':')
		if err != nil ***REMOVED***
			return
		***REMOVED***
		n++
		nn, err = rwNext(dst, src)
		n += nn
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if !comma ***REMOVED***
			comma = true
		***REMOVED***
	***REMOVED***

	err = dst.WriteByte('***REMOVED***')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++
	return
***REMOVED***

func rwArray(dst jsWriter, src *Reader) (n int, err error) ***REMOVED***
	err = dst.WriteByte('[')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var sz uint32
	var nn int
	sz, err = src.ReadArrayHeader()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	comma := false
	for i := uint32(0); i < sz; i++ ***REMOVED***
		if comma ***REMOVED***
			err = dst.WriteByte(',')
			if err != nil ***REMOVED***
				return
			***REMOVED***
			n++
		***REMOVED***
		nn, err = rwNext(dst, src)
		n += nn
		if err != nil ***REMOVED***
			return
		***REMOVED***
		comma = true
	***REMOVED***

	err = dst.WriteByte(']')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++
	return
***REMOVED***

func rwNil(dst jsWriter, src *Reader) (int, error) ***REMOVED***
	err := src.ReadNil()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return dst.Write(null)
***REMOVED***

func rwFloat32(dst jsWriter, src *Reader) (int, error) ***REMOVED***
	f, err := src.ReadFloat32()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	src.scratch = strconv.AppendFloat(src.scratch[:0], float64(f), 'f', -1, 64)
	return dst.Write(src.scratch)
***REMOVED***

func rwFloat64(dst jsWriter, src *Reader) (int, error) ***REMOVED***
	f, err := src.ReadFloat64()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	src.scratch = strconv.AppendFloat(src.scratch[:0], f, 'f', -1, 32)
	return dst.Write(src.scratch)
***REMOVED***

func rwInt(dst jsWriter, src *Reader) (int, error) ***REMOVED***
	i, err := src.ReadInt64()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	src.scratch = strconv.AppendInt(src.scratch[:0], i, 10)
	return dst.Write(src.scratch)
***REMOVED***

func rwUint(dst jsWriter, src *Reader) (int, error) ***REMOVED***
	u, err := src.ReadUint64()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	src.scratch = strconv.AppendUint(src.scratch[:0], u, 10)
	return dst.Write(src.scratch)
***REMOVED***

func rwBool(dst jsWriter, src *Reader) (int, error) ***REMOVED***
	b, err := src.ReadBool()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if b ***REMOVED***
		return dst.WriteString("true")
	***REMOVED***
	return dst.WriteString("false")
***REMOVED***

func rwTime(dst jsWriter, src *Reader) (int, error) ***REMOVED***
	t, err := src.ReadTime()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	bts, err := t.MarshalJSON()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return dst.Write(bts)
***REMOVED***

func rwExtension(dst jsWriter, src *Reader) (n int, err error) ***REMOVED***
	et, err := src.peekExtensionType()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	// registered extensions can override
	// the JSON encoding
	if j, ok := extensionReg[et]; ok ***REMOVED***
		var bts []byte
		e := j()
		err = src.ReadExtension(e)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		bts, err = json.Marshal(e)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		return dst.Write(bts)
	***REMOVED***

	e := RawExtension***REMOVED******REMOVED***
	e.Type = et
	err = src.ReadExtension(&e)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	var nn int
	err = dst.WriteByte('***REMOVED***')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++

	nn, err = dst.WriteString(`"type:"`)
	n += nn
	if err != nil ***REMOVED***
		return
	***REMOVED***

	src.scratch = strconv.AppendInt(src.scratch[0:0], int64(e.Type), 10)
	nn, err = dst.Write(src.scratch)
	n += nn
	if err != nil ***REMOVED***
		return
	***REMOVED***

	nn, err = dst.WriteString(`,"data":"`)
	n += nn
	if err != nil ***REMOVED***
		return
	***REMOVED***

	enc := base64.NewEncoder(base64.StdEncoding, dst)

	nn, err = enc.Write(e.Data)
	n += nn
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = enc.Close()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	nn, err = dst.WriteString(`"***REMOVED***`)
	n += nn
	return
***REMOVED***

func rwString(dst jsWriter, src *Reader) (n int, err error) ***REMOVED***
	var p []byte
	p, err = src.R.Peek(1)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	lead := p[0]
	var read int

	if isfixstr(lead) ***REMOVED***
		read = int(rfixstr(lead))
		src.R.Skip(1)
		goto write
	***REMOVED***

	switch lead ***REMOVED***
	case mstr8:
		p, err = src.R.Next(2)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int(uint8(p[1]))
	case mstr16:
		p, err = src.R.Next(3)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int(big.Uint16(p[1:]))
	case mstr32:
		p, err = src.R.Next(5)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		read = int(big.Uint32(p[1:]))
	default:
		err = badPrefix(StrType, lead)
		return
	***REMOVED***
write:
	p, err = src.R.Next(read)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n, err = rwquoted(dst, p)
	return
***REMOVED***

func rwBytes(dst jsWriter, src *Reader) (n int, err error) ***REMOVED***
	var nn int
	err = dst.WriteByte('"')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++
	src.scratch, err = src.ReadBytes(src.scratch[:0])
	if err != nil ***REMOVED***
		return
	***REMOVED***
	enc := base64.NewEncoder(base64.StdEncoding, dst)
	nn, err = enc.Write(src.scratch)
	n += nn
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = enc.Close()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = dst.WriteByte('"')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++
	return
***REMOVED***

// Below (c) The Go Authors, 2009-2014
// Subject to the BSD-style license found at http://golang.org
//
// see: encoding/json/encode.go:(*encodeState).stringbytes()
func rwquoted(dst jsWriter, s []byte) (n int, err error) ***REMOVED***
	var nn int
	err = dst.WriteByte('"')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++
	start := 0
	for i := 0; i < len(s); ***REMOVED***
		if b := s[i]; b < utf8.RuneSelf ***REMOVED***
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' ***REMOVED***
				i++
				continue
			***REMOVED***
			if start < i ***REMOVED***
				nn, err = dst.Write(s[start:i])
				n += nn
				if err != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
			switch b ***REMOVED***
			case '\\', '"':
				err = dst.WriteByte('\\')
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
				err = dst.WriteByte(b)
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
			case '\n':
				err = dst.WriteByte('\\')
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
				err = dst.WriteByte('n')
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
			case '\r':
				err = dst.WriteByte('\\')
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
				err = dst.WriteByte('r')
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
			default:
				nn, err = dst.WriteString(`\u00`)
				n += nn
				if err != nil ***REMOVED***
					return
				***REMOVED***
				err = dst.WriteByte(hex[b>>4])
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
				err = dst.WriteByte(hex[b&0xF])
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
			***REMOVED***
			i++
			start = i
			continue
		***REMOVED***
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 ***REMOVED***
			if start < i ***REMOVED***
				nn, err = dst.Write(s[start:i])
				n += nn
				if err != nil ***REMOVED***
					return
				***REMOVED***
				nn, err = dst.WriteString(`\ufffd`)
				n += nn
				if err != nil ***REMOVED***
					return
				***REMOVED***
				i += size
				start = i
				continue
			***REMOVED***
		***REMOVED***
		if c == '\u2028' || c == '\u2029' ***REMOVED***
			if start < i ***REMOVED***
				nn, err = dst.Write(s[start:i])
				n += nn
				if err != nil ***REMOVED***
					return
				***REMOVED***
				nn, err = dst.WriteString(`\u202`)
				n += nn
				if err != nil ***REMOVED***
					return
				***REMOVED***
				err = dst.WriteByte(hex[c&0xF])
				if err != nil ***REMOVED***
					return
				***REMOVED***
				n++
			***REMOVED***
		***REMOVED***
		i += size
	***REMOVED***
	if start < len(s) ***REMOVED***
		nn, err = dst.Write(s[start:])
		n += nn
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	err = dst.WriteByte('"')
	if err != nil ***REMOVED***
		return
	***REMOVED***
	n++
	return
***REMOVED***
