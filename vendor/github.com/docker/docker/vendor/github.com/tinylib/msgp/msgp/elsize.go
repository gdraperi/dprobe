package msgp

// size of every object on the wire,
// plus type information. gives us
// constant-time type information
// for traversing composite objects.
//
var sizes = [256]bytespec***REMOVED***
	mnil:      ***REMOVED***size: 1, extra: constsize, typ: NilType***REMOVED***,
	mfalse:    ***REMOVED***size: 1, extra: constsize, typ: BoolType***REMOVED***,
	mtrue:     ***REMOVED***size: 1, extra: constsize, typ: BoolType***REMOVED***,
	mbin8:     ***REMOVED***size: 2, extra: extra8, typ: BinType***REMOVED***,
	mbin16:    ***REMOVED***size: 3, extra: extra16, typ: BinType***REMOVED***,
	mbin32:    ***REMOVED***size: 5, extra: extra32, typ: BinType***REMOVED***,
	mext8:     ***REMOVED***size: 3, extra: extra8, typ: ExtensionType***REMOVED***,
	mext16:    ***REMOVED***size: 4, extra: extra16, typ: ExtensionType***REMOVED***,
	mext32:    ***REMOVED***size: 6, extra: extra32, typ: ExtensionType***REMOVED***,
	mfloat32:  ***REMOVED***size: 5, extra: constsize, typ: Float32Type***REMOVED***,
	mfloat64:  ***REMOVED***size: 9, extra: constsize, typ: Float64Type***REMOVED***,
	muint8:    ***REMOVED***size: 2, extra: constsize, typ: UintType***REMOVED***,
	muint16:   ***REMOVED***size: 3, extra: constsize, typ: UintType***REMOVED***,
	muint32:   ***REMOVED***size: 5, extra: constsize, typ: UintType***REMOVED***,
	muint64:   ***REMOVED***size: 9, extra: constsize, typ: UintType***REMOVED***,
	mint8:     ***REMOVED***size: 2, extra: constsize, typ: IntType***REMOVED***,
	mint16:    ***REMOVED***size: 3, extra: constsize, typ: IntType***REMOVED***,
	mint32:    ***REMOVED***size: 5, extra: constsize, typ: IntType***REMOVED***,
	mint64:    ***REMOVED***size: 9, extra: constsize, typ: IntType***REMOVED***,
	mfixext1:  ***REMOVED***size: 3, extra: constsize, typ: ExtensionType***REMOVED***,
	mfixext2:  ***REMOVED***size: 4, extra: constsize, typ: ExtensionType***REMOVED***,
	mfixext4:  ***REMOVED***size: 6, extra: constsize, typ: ExtensionType***REMOVED***,
	mfixext8:  ***REMOVED***size: 10, extra: constsize, typ: ExtensionType***REMOVED***,
	mfixext16: ***REMOVED***size: 18, extra: constsize, typ: ExtensionType***REMOVED***,
	mstr8:     ***REMOVED***size: 2, extra: extra8, typ: StrType***REMOVED***,
	mstr16:    ***REMOVED***size: 3, extra: extra16, typ: StrType***REMOVED***,
	mstr32:    ***REMOVED***size: 5, extra: extra32, typ: StrType***REMOVED***,
	marray16:  ***REMOVED***size: 3, extra: array16v, typ: ArrayType***REMOVED***,
	marray32:  ***REMOVED***size: 5, extra: array32v, typ: ArrayType***REMOVED***,
	mmap16:    ***REMOVED***size: 3, extra: map16v, typ: MapType***REMOVED***,
	mmap32:    ***REMOVED***size: 5, extra: map32v, typ: MapType***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	// set up fixed fields

	// fixint
	for i := mfixint; i < 0x80; i++ ***REMOVED***
		sizes[i] = bytespec***REMOVED***size: 1, extra: constsize, typ: IntType***REMOVED***
	***REMOVED***

	// nfixint
	for i := uint16(mnfixint); i < 0x100; i++ ***REMOVED***
		sizes[uint8(i)] = bytespec***REMOVED***size: 1, extra: constsize, typ: IntType***REMOVED***
	***REMOVED***

	// fixstr gets constsize,
	// since the prefix yields the size
	for i := mfixstr; i < 0xc0; i++ ***REMOVED***
		sizes[i] = bytespec***REMOVED***size: 1 + rfixstr(i), extra: constsize, typ: StrType***REMOVED***
	***REMOVED***

	// fixmap
	for i := mfixmap; i < 0x90; i++ ***REMOVED***
		sizes[i] = bytespec***REMOVED***size: 1, extra: varmode(2 * rfixmap(i)), typ: MapType***REMOVED***
	***REMOVED***

	// fixarray
	for i := mfixarray; i < 0xa0; i++ ***REMOVED***
		sizes[i] = bytespec***REMOVED***size: 1, extra: varmode(rfixarray(i)), typ: ArrayType***REMOVED***
	***REMOVED***
***REMOVED***

// a valid bytespsec has
// non-zero 'size' and
// non-zero 'typ'
type bytespec struct ***REMOVED***
	size  uint8   // prefix size information
	extra varmode // extra size information
	typ   Type    // type
	_     byte    // makes bytespec 4 bytes (yes, this matters)
***REMOVED***

// size mode
// if positive, # elements for composites
type varmode int8

const (
	constsize varmode = 0  // constant size (size bytes + uint8(varmode) objects)
	extra8            = -1 // has uint8(p[1]) extra bytes
	extra16           = -2 // has be16(p[1:]) extra bytes
	extra32           = -3 // has be32(p[1:]) extra bytes
	map16v            = -4 // use map16
	map32v            = -5 // use map32
	array16v          = -6 // use array16
	array32v          = -7 // use array32
)

func getType(v byte) Type ***REMOVED***
	return sizes[v].typ
***REMOVED***
