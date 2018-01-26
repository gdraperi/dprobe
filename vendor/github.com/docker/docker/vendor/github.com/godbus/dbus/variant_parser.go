package dbus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type varParser struct ***REMOVED***
	tokens []varToken
	i      int
***REMOVED***

func (p *varParser) backup() ***REMOVED***
	p.i--
***REMOVED***

func (p *varParser) next() varToken ***REMOVED***
	if p.i < len(p.tokens) ***REMOVED***
		t := p.tokens[p.i]
		p.i++
		return t
	***REMOVED***
	return varToken***REMOVED***typ: tokEOF***REMOVED***
***REMOVED***

type varNode interface ***REMOVED***
	Infer() (Signature, error)
	String() string
	Sigs() sigSet
	Value(Signature) (interface***REMOVED******REMOVED***, error)
***REMOVED***

func varMakeNode(p *varParser) (varNode, error) ***REMOVED***
	var sig Signature

	for ***REMOVED***
		t := p.next()
		switch t.typ ***REMOVED***
		case tokEOF:
			return nil, io.ErrUnexpectedEOF
		case tokError:
			return nil, errors.New(t.val)
		case tokNumber:
			return varMakeNumNode(t, sig)
		case tokString:
			return varMakeStringNode(t, sig)
		case tokBool:
			if sig.str != "" && sig.str != "b" ***REMOVED***
				return nil, varTypeError***REMOVED***t.val, sig***REMOVED***
			***REMOVED***
			b, err := strconv.ParseBool(t.val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return boolNode(b), nil
		case tokArrayStart:
			return varMakeArrayNode(p, sig)
		case tokVariantStart:
			return varMakeVariantNode(p, sig)
		case tokDictStart:
			return varMakeDictNode(p, sig)
		case tokType:
			if sig.str != "" ***REMOVED***
				return nil, errors.New("unexpected type annotation")
			***REMOVED***
			if t.val[0] == '@' ***REMOVED***
				sig.str = t.val[1:]
			***REMOVED*** else ***REMOVED***
				sig.str = varTypeMap[t.val]
			***REMOVED***
		case tokByteString:
			if sig.str != "" && sig.str != "ay" ***REMOVED***
				return nil, varTypeError***REMOVED***t.val, sig***REMOVED***
			***REMOVED***
			b, err := varParseByteString(t.val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return byteStringNode(b), nil
		default:
			return nil, fmt.Errorf("unexpected %q", t.val)
		***REMOVED***
	***REMOVED***
***REMOVED***

type varTypeError struct ***REMOVED***
	val string
	sig Signature
***REMOVED***

func (e varTypeError) Error() string ***REMOVED***
	return fmt.Sprintf("dbus: can't parse %q as type %q", e.val, e.sig.str)
***REMOVED***

type sigSet map[Signature]bool

func (s sigSet) Empty() bool ***REMOVED***
	return len(s) == 0
***REMOVED***

func (s sigSet) Intersect(s2 sigSet) sigSet ***REMOVED***
	r := make(sigSet)
	for k := range s ***REMOVED***
		if s2[k] ***REMOVED***
			r[k] = true
		***REMOVED***
	***REMOVED***
	return r
***REMOVED***

func (s sigSet) Single() (Signature, bool) ***REMOVED***
	if len(s) == 1 ***REMOVED***
		for k := range s ***REMOVED***
			return k, true
		***REMOVED***
	***REMOVED***
	return Signature***REMOVED******REMOVED***, false
***REMOVED***

func (s sigSet) ToArray() sigSet ***REMOVED***
	r := make(sigSet, len(s))
	for k := range s ***REMOVED***
		r[Signature***REMOVED***"a" + k.str***REMOVED***] = true
	***REMOVED***
	return r
***REMOVED***

type numNode struct ***REMOVED***
	sig Signature
	str string
	val interface***REMOVED******REMOVED***
***REMOVED***

var numSigSet = sigSet***REMOVED***
	Signature***REMOVED***"y"***REMOVED***: true,
	Signature***REMOVED***"n"***REMOVED***: true,
	Signature***REMOVED***"q"***REMOVED***: true,
	Signature***REMOVED***"i"***REMOVED***: true,
	Signature***REMOVED***"u"***REMOVED***: true,
	Signature***REMOVED***"x"***REMOVED***: true,
	Signature***REMOVED***"t"***REMOVED***: true,
	Signature***REMOVED***"d"***REMOVED***: true,
***REMOVED***

func (n numNode) Infer() (Signature, error) ***REMOVED***
	if strings.ContainsAny(n.str, ".e") ***REMOVED***
		return Signature***REMOVED***"d"***REMOVED***, nil
	***REMOVED***
	return Signature***REMOVED***"i"***REMOVED***, nil
***REMOVED***

func (n numNode) String() string ***REMOVED***
	return n.str
***REMOVED***

func (n numNode) Sigs() sigSet ***REMOVED***
	if n.sig.str != "" ***REMOVED***
		return sigSet***REMOVED***n.sig: true***REMOVED***
	***REMOVED***
	if strings.ContainsAny(n.str, ".e") ***REMOVED***
		return sigSet***REMOVED***Signature***REMOVED***"d"***REMOVED***: true***REMOVED***
	***REMOVED***
	return numSigSet
***REMOVED***

func (n numNode) Value(sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if n.sig.str != "" && n.sig != sig ***REMOVED***
		return nil, varTypeError***REMOVED***n.str, sig***REMOVED***
	***REMOVED***
	if n.val != nil ***REMOVED***
		return n.val, nil
	***REMOVED***
	return varNumAs(n.str, sig)
***REMOVED***

func varMakeNumNode(tok varToken, sig Signature) (varNode, error) ***REMOVED***
	if sig.str == "" ***REMOVED***
		return numNode***REMOVED***str: tok.val***REMOVED***, nil
	***REMOVED***
	num, err := varNumAs(tok.val, sig)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return numNode***REMOVED***sig: sig, val: num***REMOVED***, nil
***REMOVED***

func varNumAs(s string, sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	isUnsigned := false
	size := 32
	switch sig.str ***REMOVED***
	case "n":
		size = 16
	case "i":
	case "x":
		size = 64
	case "y":
		size = 8
		isUnsigned = true
	case "q":
		size = 16
		isUnsigned = true
	case "u":
		isUnsigned = true
	case "t":
		size = 64
		isUnsigned = true
	case "d":
		d, err := strconv.ParseFloat(s, 64)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return d, nil
	default:
		return nil, varTypeError***REMOVED***s, sig***REMOVED***
	***REMOVED***
	base := 10
	if strings.HasPrefix(s, "0x") ***REMOVED***
		base = 16
		s = s[2:]
	***REMOVED***
	if strings.HasPrefix(s, "0") && len(s) != 1 ***REMOVED***
		base = 8
		s = s[1:]
	***REMOVED***
	if isUnsigned ***REMOVED***
		i, err := strconv.ParseUint(s, base, size)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var v interface***REMOVED******REMOVED*** = i
		switch sig.str ***REMOVED***
		case "y":
			v = byte(i)
		case "q":
			v = uint16(i)
		case "u":
			v = uint32(i)
		***REMOVED***
		return v, nil
	***REMOVED***
	i, err := strconv.ParseInt(s, base, size)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var v interface***REMOVED******REMOVED*** = i
	switch sig.str ***REMOVED***
	case "n":
		v = int16(i)
	case "i":
		v = int32(i)
	***REMOVED***
	return v, nil
***REMOVED***

type stringNode struct ***REMOVED***
	sig Signature
	str string      // parsed
	val interface***REMOVED******REMOVED*** // has correct type
***REMOVED***

var stringSigSet = sigSet***REMOVED***
	Signature***REMOVED***"s"***REMOVED***: true,
	Signature***REMOVED***"g"***REMOVED***: true,
	Signature***REMOVED***"o"***REMOVED***: true,
***REMOVED***

func (n stringNode) Infer() (Signature, error) ***REMOVED***
	return Signature***REMOVED***"s"***REMOVED***, nil
***REMOVED***

func (n stringNode) String() string ***REMOVED***
	return n.str
***REMOVED***

func (n stringNode) Sigs() sigSet ***REMOVED***
	if n.sig.str != "" ***REMOVED***
		return sigSet***REMOVED***n.sig: true***REMOVED***
	***REMOVED***
	return stringSigSet
***REMOVED***

func (n stringNode) Value(sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if n.sig.str != "" && n.sig != sig ***REMOVED***
		return nil, varTypeError***REMOVED***n.str, sig***REMOVED***
	***REMOVED***
	if n.val != nil ***REMOVED***
		return n.val, nil
	***REMOVED***
	switch ***REMOVED***
	case sig.str == "g":
		return Signature***REMOVED***n.str***REMOVED***, nil
	case sig.str == "o":
		return ObjectPath(n.str), nil
	case sig.str == "s":
		return n.str, nil
	default:
		return nil, varTypeError***REMOVED***n.str, sig***REMOVED***
	***REMOVED***
***REMOVED***

func varMakeStringNode(tok varToken, sig Signature) (varNode, error) ***REMOVED***
	if sig.str != "" && sig.str != "s" && sig.str != "g" && sig.str != "o" ***REMOVED***
		return nil, fmt.Errorf("invalid type %q for string", sig.str)
	***REMOVED***
	s, err := varParseString(tok.val)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	n := stringNode***REMOVED***str: s***REMOVED***
	if sig.str == "" ***REMOVED***
		return stringNode***REMOVED***str: s***REMOVED***, nil
	***REMOVED***
	n.sig = sig
	switch sig.str ***REMOVED***
	case "o":
		n.val = ObjectPath(s)
	case "g":
		n.val = Signature***REMOVED***s***REMOVED***
	case "s":
		n.val = s
	***REMOVED***
	return n, nil
***REMOVED***

func varParseString(s string) (string, error) ***REMOVED***
	// quotes are guaranteed to be there
	s = s[1 : len(s)-1]
	buf := new(bytes.Buffer)
	for len(s) != 0 ***REMOVED***
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 ***REMOVED***
			return "", errors.New("invalid UTF-8")
		***REMOVED***
		s = s[size:]
		if r != '\\' ***REMOVED***
			buf.WriteRune(r)
			continue
		***REMOVED***
		r, size = utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 ***REMOVED***
			return "", errors.New("invalid UTF-8")
		***REMOVED***
		s = s[size:]
		switch r ***REMOVED***
		case 'a':
			buf.WriteRune(0x7)
		case 'b':
			buf.WriteRune(0x8)
		case 'f':
			buf.WriteRune(0xc)
		case 'n':
			buf.WriteRune('\n')
		case 'r':
			buf.WriteRune('\r')
		case 't':
			buf.WriteRune('\t')
		case '\n':
		case 'u':
			if len(s) < 4 ***REMOVED***
				return "", errors.New("short unicode escape")
			***REMOVED***
			r, err := strconv.ParseUint(s[:4], 16, 32)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			buf.WriteRune(rune(r))
			s = s[4:]
		case 'U':
			if len(s) < 8 ***REMOVED***
				return "", errors.New("short unicode escape")
			***REMOVED***
			r, err := strconv.ParseUint(s[:8], 16, 32)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			buf.WriteRune(rune(r))
			s = s[8:]
		default:
			buf.WriteRune(r)
		***REMOVED***
	***REMOVED***
	return buf.String(), nil
***REMOVED***

var boolSigSet = sigSet***REMOVED***Signature***REMOVED***"b"***REMOVED***: true***REMOVED***

type boolNode bool

func (boolNode) Infer() (Signature, error) ***REMOVED***
	return Signature***REMOVED***"b"***REMOVED***, nil
***REMOVED***

func (b boolNode) String() string ***REMOVED***
	if b ***REMOVED***
		return "true"
	***REMOVED***
	return "false"
***REMOVED***

func (boolNode) Sigs() sigSet ***REMOVED***
	return boolSigSet
***REMOVED***

func (b boolNode) Value(sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if sig.str != "b" ***REMOVED***
		return nil, varTypeError***REMOVED***b.String(), sig***REMOVED***
	***REMOVED***
	return bool(b), nil
***REMOVED***

type arrayNode struct ***REMOVED***
	set      sigSet
	children []varNode
	val      interface***REMOVED******REMOVED***
***REMOVED***

func (n arrayNode) Infer() (Signature, error) ***REMOVED***
	for _, v := range n.children ***REMOVED***
		csig, err := varInfer(v)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		return Signature***REMOVED***"a" + csig.str***REMOVED***, nil
	***REMOVED***
	return Signature***REMOVED******REMOVED***, fmt.Errorf("can't infer type for %q", n.String())
***REMOVED***

func (n arrayNode) String() string ***REMOVED***
	s := "["
	for i, v := range n.children ***REMOVED***
		s += v.String()
		if i != len(n.children)-1 ***REMOVED***
			s += ", "
		***REMOVED***
	***REMOVED***
	return s + "]"
***REMOVED***

func (n arrayNode) Sigs() sigSet ***REMOVED***
	return n.set
***REMOVED***

func (n arrayNode) Value(sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if n.set.Empty() ***REMOVED***
		// no type information whatsoever, so this must be an empty slice
		return reflect.MakeSlice(typeFor(sig.str), 0, 0).Interface(), nil
	***REMOVED***
	if !n.set[sig] ***REMOVED***
		return nil, varTypeError***REMOVED***n.String(), sig***REMOVED***
	***REMOVED***
	s := reflect.MakeSlice(typeFor(sig.str), len(n.children), len(n.children))
	for i, v := range n.children ***REMOVED***
		rv, err := v.Value(Signature***REMOVED***sig.str[1:]***REMOVED***)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		s.Index(i).Set(reflect.ValueOf(rv))
	***REMOVED***
	return s.Interface(), nil
***REMOVED***

func varMakeArrayNode(p *varParser, sig Signature) (varNode, error) ***REMOVED***
	var n arrayNode
	if sig.str != "" ***REMOVED***
		n.set = sigSet***REMOVED***sig: true***REMOVED***
	***REMOVED***
	if t := p.next(); t.typ == tokArrayEnd ***REMOVED***
		return n, nil
	***REMOVED*** else ***REMOVED***
		p.backup()
	***REMOVED***
Loop:
	for ***REMOVED***
		t := p.next()
		switch t.typ ***REMOVED***
		case tokEOF:
			return nil, io.ErrUnexpectedEOF
		case tokError:
			return nil, errors.New(t.val)
		***REMOVED***
		p.backup()
		cn, err := varMakeNode(p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if cset := cn.Sigs(); !cset.Empty() ***REMOVED***
			if n.set.Empty() ***REMOVED***
				n.set = cset.ToArray()
			***REMOVED*** else ***REMOVED***
				nset := cset.ToArray().Intersect(n.set)
				if nset.Empty() ***REMOVED***
					return nil, fmt.Errorf("can't parse %q with given type information", cn.String())
				***REMOVED***
				n.set = nset
			***REMOVED***
		***REMOVED***
		n.children = append(n.children, cn)
		switch t := p.next(); t.typ ***REMOVED***
		case tokEOF:
			return nil, io.ErrUnexpectedEOF
		case tokError:
			return nil, errors.New(t.val)
		case tokArrayEnd:
			break Loop
		case tokComma:
			continue
		default:
			return nil, fmt.Errorf("unexpected %q", t.val)
		***REMOVED***
	***REMOVED***
	return n, nil
***REMOVED***

type variantNode struct ***REMOVED***
	n varNode
***REMOVED***

var variantSet = sigSet***REMOVED***
	Signature***REMOVED***"v"***REMOVED***: true,
***REMOVED***

func (variantNode) Infer() (Signature, error) ***REMOVED***
	return Signature***REMOVED***"v"***REMOVED***, nil
***REMOVED***

func (n variantNode) String() string ***REMOVED***
	return "<" + n.n.String() + ">"
***REMOVED***

func (variantNode) Sigs() sigSet ***REMOVED***
	return variantSet
***REMOVED***

func (n variantNode) Value(sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if sig.str != "v" ***REMOVED***
		return nil, varTypeError***REMOVED***n.String(), sig***REMOVED***
	***REMOVED***
	sig, err := varInfer(n.n)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	v, err := n.n.Value(sig)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return MakeVariant(v), nil
***REMOVED***

func varMakeVariantNode(p *varParser, sig Signature) (varNode, error) ***REMOVED***
	n, err := varMakeNode(p)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if t := p.next(); t.typ != tokVariantEnd ***REMOVED***
		return nil, fmt.Errorf("unexpected %q", t.val)
	***REMOVED***
	vn := variantNode***REMOVED***n***REMOVED***
	if sig.str != "" && sig.str != "v" ***REMOVED***
		return nil, varTypeError***REMOVED***vn.String(), sig***REMOVED***
	***REMOVED***
	return variantNode***REMOVED***n***REMOVED***, nil
***REMOVED***

type dictEntry struct ***REMOVED***
	key, val varNode
***REMOVED***

type dictNode struct ***REMOVED***
	kset, vset sigSet
	children   []dictEntry
	val        interface***REMOVED******REMOVED***
***REMOVED***

func (n dictNode) Infer() (Signature, error) ***REMOVED***
	for _, v := range n.children ***REMOVED***
		ksig, err := varInfer(v.key)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		vsig, err := varInfer(v.val)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		return Signature***REMOVED***"a***REMOVED***" + ksig.str + vsig.str + "***REMOVED***"***REMOVED***, nil
	***REMOVED***
	return Signature***REMOVED******REMOVED***, fmt.Errorf("can't infer type for %q", n.String())
***REMOVED***

func (n dictNode) String() string ***REMOVED***
	s := "***REMOVED***"
	for i, v := range n.children ***REMOVED***
		s += v.key.String() + ": " + v.val.String()
		if i != len(n.children)-1 ***REMOVED***
			s += ", "
		***REMOVED***
	***REMOVED***
	return s + "***REMOVED***"
***REMOVED***

func (n dictNode) Sigs() sigSet ***REMOVED***
	r := sigSet***REMOVED******REMOVED***
	for k := range n.kset ***REMOVED***
		for v := range n.vset ***REMOVED***
			sig := "a***REMOVED***" + k.str + v.str + "***REMOVED***"
			r[Signature***REMOVED***sig***REMOVED***] = true
		***REMOVED***
	***REMOVED***
	return r
***REMOVED***

func (n dictNode) Value(sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	set := n.Sigs()
	if set.Empty() ***REMOVED***
		// no type information -> empty dict
		return reflect.MakeMap(typeFor(sig.str)).Interface(), nil
	***REMOVED***
	if !set[sig] ***REMOVED***
		return nil, varTypeError***REMOVED***n.String(), sig***REMOVED***
	***REMOVED***
	m := reflect.MakeMap(typeFor(sig.str))
	ksig := Signature***REMOVED***sig.str[2:3]***REMOVED***
	vsig := Signature***REMOVED***sig.str[3 : len(sig.str)-1]***REMOVED***
	for _, v := range n.children ***REMOVED***
		kv, err := v.key.Value(ksig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		vv, err := v.val.Value(vsig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m.SetMapIndex(reflect.ValueOf(kv), reflect.ValueOf(vv))
	***REMOVED***
	return m.Interface(), nil
***REMOVED***

func varMakeDictNode(p *varParser, sig Signature) (varNode, error) ***REMOVED***
	var n dictNode

	if sig.str != "" ***REMOVED***
		if len(sig.str) < 5 ***REMOVED***
			return nil, fmt.Errorf("invalid signature %q for dict type", sig)
		***REMOVED***
		ksig := Signature***REMOVED***string(sig.str[2])***REMOVED***
		vsig := Signature***REMOVED***sig.str[3 : len(sig.str)-1]***REMOVED***
		n.kset = sigSet***REMOVED***ksig: true***REMOVED***
		n.vset = sigSet***REMOVED***vsig: true***REMOVED***
	***REMOVED***
	if t := p.next(); t.typ == tokDictEnd ***REMOVED***
		return n, nil
	***REMOVED*** else ***REMOVED***
		p.backup()
	***REMOVED***
Loop:
	for ***REMOVED***
		t := p.next()
		switch t.typ ***REMOVED***
		case tokEOF:
			return nil, io.ErrUnexpectedEOF
		case tokError:
			return nil, errors.New(t.val)
		***REMOVED***
		p.backup()
		kn, err := varMakeNode(p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if kset := kn.Sigs(); !kset.Empty() ***REMOVED***
			if n.kset.Empty() ***REMOVED***
				n.kset = kset
			***REMOVED*** else ***REMOVED***
				n.kset = kset.Intersect(n.kset)
				if n.kset.Empty() ***REMOVED***
					return nil, fmt.Errorf("can't parse %q with given type information", kn.String())
				***REMOVED***
			***REMOVED***
		***REMOVED***
		t = p.next()
		switch t.typ ***REMOVED***
		case tokEOF:
			return nil, io.ErrUnexpectedEOF
		case tokError:
			return nil, errors.New(t.val)
		case tokColon:
		default:
			return nil, fmt.Errorf("unexpected %q", t.val)
		***REMOVED***
		t = p.next()
		switch t.typ ***REMOVED***
		case tokEOF:
			return nil, io.ErrUnexpectedEOF
		case tokError:
			return nil, errors.New(t.val)
		***REMOVED***
		p.backup()
		vn, err := varMakeNode(p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if vset := vn.Sigs(); !vset.Empty() ***REMOVED***
			if n.vset.Empty() ***REMOVED***
				n.vset = vset
			***REMOVED*** else ***REMOVED***
				n.vset = n.vset.Intersect(vset)
				if n.vset.Empty() ***REMOVED***
					return nil, fmt.Errorf("can't parse %q with given type information", vn.String())
				***REMOVED***
			***REMOVED***
		***REMOVED***
		n.children = append(n.children, dictEntry***REMOVED***kn, vn***REMOVED***)
		t = p.next()
		switch t.typ ***REMOVED***
		case tokEOF:
			return nil, io.ErrUnexpectedEOF
		case tokError:
			return nil, errors.New(t.val)
		case tokDictEnd:
			break Loop
		case tokComma:
			continue
		default:
			return nil, fmt.Errorf("unexpected %q", t.val)
		***REMOVED***
	***REMOVED***
	return n, nil
***REMOVED***

type byteStringNode []byte

var byteStringSet = sigSet***REMOVED***
	Signature***REMOVED***"ay"***REMOVED***: true,
***REMOVED***

func (byteStringNode) Infer() (Signature, error) ***REMOVED***
	return Signature***REMOVED***"ay"***REMOVED***, nil
***REMOVED***

func (b byteStringNode) String() string ***REMOVED***
	return string(b)
***REMOVED***

func (b byteStringNode) Sigs() sigSet ***REMOVED***
	return byteStringSet
***REMOVED***

func (b byteStringNode) Value(sig Signature) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if sig.str != "ay" ***REMOVED***
		return nil, varTypeError***REMOVED***b.String(), sig***REMOVED***
	***REMOVED***
	return []byte(b), nil
***REMOVED***

func varParseByteString(s string) ([]byte, error) ***REMOVED***
	// quotes and b at start are guaranteed to be there
	b := make([]byte, 0, 1)
	s = s[2 : len(s)-1]
	for len(s) != 0 ***REMOVED***
		c := s[0]
		s = s[1:]
		if c != '\\' ***REMOVED***
			b = append(b, c)
			continue
		***REMOVED***
		c = s[0]
		s = s[1:]
		switch c ***REMOVED***
		case 'a':
			b = append(b, 0x7)
		case 'b':
			b = append(b, 0x8)
		case 'f':
			b = append(b, 0xc)
		case 'n':
			b = append(b, '\n')
		case 'r':
			b = append(b, '\r')
		case 't':
			b = append(b, '\t')
		case 'x':
			if len(s) < 2 ***REMOVED***
				return nil, errors.New("short escape")
			***REMOVED***
			n, err := strconv.ParseUint(s[:2], 16, 8)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			b = append(b, byte(n))
			s = s[2:]
		case '0':
			if len(s) < 3 ***REMOVED***
				return nil, errors.New("short escape")
			***REMOVED***
			n, err := strconv.ParseUint(s[:3], 8, 8)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			b = append(b, byte(n))
			s = s[3:]
		default:
			b = append(b, c)
		***REMOVED***
	***REMOVED***
	return append(b, 0), nil
***REMOVED***

func varInfer(n varNode) (Signature, error) ***REMOVED***
	if sig, ok := n.Sigs().Single(); ok ***REMOVED***
		return sig, nil
	***REMOVED***
	return n.Infer()
***REMOVED***
