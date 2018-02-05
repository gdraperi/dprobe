// This file is generated with "go test -tags generate". DO NOT EDIT!
// +build !generate

package triegen_test

// lookup returns the trie value for the first UTF-8 encoding in s and
// the width in bytes of this encoding. The size will be 0 if s does not
// hold enough bytes to complete the encoding. len(s) must be greater than 0.
func (t *randTrie) lookup(s []byte) (v uint8, sz int) ***REMOVED***
	c0 := s[0]
	switch ***REMOVED***
	case c0 < 0x80: // is ASCII
		return randValues[c0], 1
	case c0 < 0xC2:
		return 0, 1 // Illegal UTF-8: not a starter, not ASCII.
	case c0 < 0xE0: // 2-byte UTF-8
		if len(s) < 2 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := randIndex[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c1), 2
	case c0 < 0xF0: // 3-byte UTF-8
		if len(s) < 3 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := randIndex[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = randIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c2), 3
	case c0 < 0xF8: // 4-byte UTF-8
		if len(s) < 4 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := randIndex[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = randIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o = uint32(i)<<6 + uint32(c2)
		i = randIndex[o]
		c3 := s[3]
		if c3 < 0x80 || 0xC0 <= c3 ***REMOVED***
			return 0, 3 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c3), 4
	***REMOVED***
	// Illegal rune
	return 0, 1
***REMOVED***

// lookupUnsafe returns the trie value for the first UTF-8 encoding in s.
// s must start with a full and valid UTF-8 encoded rune.
func (t *randTrie) lookupUnsafe(s []byte) uint8 ***REMOVED***
	c0 := s[0]
	if c0 < 0x80 ***REMOVED*** // is ASCII
		return randValues[c0]
	***REMOVED***
	i := randIndex[c0]
	if c0 < 0xE0 ***REMOVED*** // 2-byte UTF-8
		return t.lookupValue(uint32(i), s[1])
	***REMOVED***
	i = randIndex[uint32(i)<<6+uint32(s[1])]
	if c0 < 0xF0 ***REMOVED*** // 3-byte UTF-8
		return t.lookupValue(uint32(i), s[2])
	***REMOVED***
	i = randIndex[uint32(i)<<6+uint32(s[2])]
	if c0 < 0xF8 ***REMOVED*** // 4-byte UTF-8
		return t.lookupValue(uint32(i), s[3])
	***REMOVED***
	return 0
***REMOVED***

// lookupString returns the trie value for the first UTF-8 encoding in s and
// the width in bytes of this encoding. The size will be 0 if s does not
// hold enough bytes to complete the encoding. len(s) must be greater than 0.
func (t *randTrie) lookupString(s string) (v uint8, sz int) ***REMOVED***
	c0 := s[0]
	switch ***REMOVED***
	case c0 < 0x80: // is ASCII
		return randValues[c0], 1
	case c0 < 0xC2:
		return 0, 1 // Illegal UTF-8: not a starter, not ASCII.
	case c0 < 0xE0: // 2-byte UTF-8
		if len(s) < 2 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := randIndex[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c1), 2
	case c0 < 0xF0: // 3-byte UTF-8
		if len(s) < 3 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := randIndex[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = randIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c2), 3
	case c0 < 0xF8: // 4-byte UTF-8
		if len(s) < 4 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := randIndex[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = randIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o = uint32(i)<<6 + uint32(c2)
		i = randIndex[o]
		c3 := s[3]
		if c3 < 0x80 || 0xC0 <= c3 ***REMOVED***
			return 0, 3 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c3), 4
	***REMOVED***
	// Illegal rune
	return 0, 1
***REMOVED***

// lookupStringUnsafe returns the trie value for the first UTF-8 encoding in s.
// s must start with a full and valid UTF-8 encoded rune.
func (t *randTrie) lookupStringUnsafe(s string) uint8 ***REMOVED***
	c0 := s[0]
	if c0 < 0x80 ***REMOVED*** // is ASCII
		return randValues[c0]
	***REMOVED***
	i := randIndex[c0]
	if c0 < 0xE0 ***REMOVED*** // 2-byte UTF-8
		return t.lookupValue(uint32(i), s[1])
	***REMOVED***
	i = randIndex[uint32(i)<<6+uint32(s[1])]
	if c0 < 0xF0 ***REMOVED*** // 3-byte UTF-8
		return t.lookupValue(uint32(i), s[2])
	***REMOVED***
	i = randIndex[uint32(i)<<6+uint32(s[2])]
	if c0 < 0xF8 ***REMOVED*** // 4-byte UTF-8
		return t.lookupValue(uint32(i), s[3])
	***REMOVED***
	return 0
***REMOVED***

// randTrie. Total size: 9280 bytes (9.06 KiB). Checksum: 6debd324a8debb8f.
type randTrie struct***REMOVED******REMOVED***

func newRandTrie(i int) *randTrie ***REMOVED***
	return &randTrie***REMOVED******REMOVED***
***REMOVED***

// lookupValue determines the type of block n and looks up the value for b.
func (t *randTrie) lookupValue(n uint32, b byte) uint8 ***REMOVED***
	switch ***REMOVED***
	default:
		return uint8(randValues[n<<6+uint32(b)])
	***REMOVED***
***REMOVED***

// randValues: 56 blocks, 3584 entries, 3584 bytes
// The third block is the zero block.
var randValues = [3584]uint8***REMOVED***
	// Block 0x0, offset 0x0
	// Block 0x1, offset 0x40
	// Block 0x2, offset 0x80
	// Block 0x3, offset 0xc0
	0xc9: 0x0001,
	// Block 0x4, offset 0x100
	0x100: 0x0001,
	// Block 0x5, offset 0x140
	0x155: 0x0001,
	// Block 0x6, offset 0x180
	0x196: 0x0001,
	// Block 0x7, offset 0x1c0
	0x1ef: 0x0001,
	// Block 0x8, offset 0x200
	0x206: 0x0001,
	// Block 0x9, offset 0x240
	0x258: 0x0001,
	// Block 0xa, offset 0x280
	0x288: 0x0001,
	// Block 0xb, offset 0x2c0
	0x2f2: 0x0001,
	// Block 0xc, offset 0x300
	0x304: 0x0001,
	// Block 0xd, offset 0x340
	0x34b: 0x0001,
	// Block 0xe, offset 0x380
	0x3ba: 0x0001,
	// Block 0xf, offset 0x3c0
	0x3f5: 0x0001,
	// Block 0x10, offset 0x400
	0x41d: 0x0001,
	// Block 0x11, offset 0x440
	0x442: 0x0001,
	// Block 0x12, offset 0x480
	0x4bb: 0x0001,
	// Block 0x13, offset 0x4c0
	0x4e9: 0x0001,
	// Block 0x14, offset 0x500
	0x53e: 0x0001,
	// Block 0x15, offset 0x540
	0x55f: 0x0001,
	// Block 0x16, offset 0x580
	0x5b7: 0x0001,
	// Block 0x17, offset 0x5c0
	0x5d9: 0x0001,
	// Block 0x18, offset 0x600
	0x60e: 0x0001,
	// Block 0x19, offset 0x640
	0x652: 0x0001,
	// Block 0x1a, offset 0x680
	0x68f: 0x0001,
	// Block 0x1b, offset 0x6c0
	0x6dc: 0x0001,
	// Block 0x1c, offset 0x700
	0x703: 0x0001,
	// Block 0x1d, offset 0x740
	0x741: 0x0001,
	// Block 0x1e, offset 0x780
	0x79b: 0x0001,
	// Block 0x1f, offset 0x7c0
	0x7f1: 0x0001,
	// Block 0x20, offset 0x800
	0x833: 0x0001,
	// Block 0x21, offset 0x840
	0x853: 0x0001,
	// Block 0x22, offset 0x880
	0x8a2: 0x0001,
	// Block 0x23, offset 0x8c0
	0x8f8: 0x0001,
	// Block 0x24, offset 0x900
	0x917: 0x0001,
	// Block 0x25, offset 0x940
	0x945: 0x0001,
	// Block 0x26, offset 0x980
	0x99e: 0x0001,
	// Block 0x27, offset 0x9c0
	0x9fd: 0x0001,
	// Block 0x28, offset 0xa00
	0xa0d: 0x0001,
	// Block 0x29, offset 0xa40
	0xa66: 0x0001,
	// Block 0x2a, offset 0xa80
	0xaab: 0x0001,
	// Block 0x2b, offset 0xac0
	0xaea: 0x0001,
	// Block 0x2c, offset 0xb00
	0xb2d: 0x0001,
	// Block 0x2d, offset 0xb40
	0xb54: 0x0001,
	// Block 0x2e, offset 0xb80
	0xb90: 0x0001,
	// Block 0x2f, offset 0xbc0
	0xbe5: 0x0001,
	// Block 0x30, offset 0xc00
	0xc28: 0x0001,
	// Block 0x31, offset 0xc40
	0xc7c: 0x0001,
	// Block 0x32, offset 0xc80
	0xcbf: 0x0001,
	// Block 0x33, offset 0xcc0
	0xcc7: 0x0001,
	// Block 0x34, offset 0xd00
	0xd34: 0x0001,
	// Block 0x35, offset 0xd40
	0xd61: 0x0001,
	// Block 0x36, offset 0xd80
	0xdb9: 0x0001,
	// Block 0x37, offset 0xdc0
	0xdda: 0x0001,
***REMOVED***

// randIndex: 89 blocks, 5696 entries, 5696 bytes
// Block 0 is the zero block.
var randIndex = [5696]uint8***REMOVED***
	// Block 0x0, offset 0x0
	// Block 0x1, offset 0x40
	// Block 0x2, offset 0x80
	// Block 0x3, offset 0xc0
	0xe1: 0x02, 0xe3: 0x03, 0xe4: 0x04,
	0xea: 0x05, 0xeb: 0x06, 0xec: 0x07,
	0xf0: 0x10, 0xf1: 0x24, 0xf2: 0x3d, 0xf3: 0x4f, 0xf4: 0x56,
	// Block 0x4, offset 0x100
	0x107: 0x01,
	// Block 0x5, offset 0x140
	0x16c: 0x02,
	// Block 0x6, offset 0x180
	0x19c: 0x03,
	0x1ae: 0x04,
	// Block 0x7, offset 0x1c0
	0x1d8: 0x05,
	0x1f7: 0x06,
	// Block 0x8, offset 0x200
	0x20c: 0x07,
	// Block 0x9, offset 0x240
	0x24a: 0x08,
	// Block 0xa, offset 0x280
	0x2b6: 0x09,
	// Block 0xb, offset 0x2c0
	0x2d5: 0x0a,
	// Block 0xc, offset 0x300
	0x31a: 0x0b,
	// Block 0xd, offset 0x340
	0x373: 0x0c,
	// Block 0xe, offset 0x380
	0x38b: 0x0d,
	// Block 0xf, offset 0x3c0
	0x3f0: 0x0e,
	// Block 0x10, offset 0x400
	0x433: 0x0f,
	// Block 0x11, offset 0x440
	0x45d: 0x10,
	// Block 0x12, offset 0x480
	0x491: 0x08, 0x494: 0x09, 0x497: 0x0a,
	0x49b: 0x0b, 0x49c: 0x0c,
	0x4a1: 0x0d,
	0x4ad: 0x0e,
	0x4ba: 0x0f,
	// Block 0x13, offset 0x4c0
	0x4c1: 0x11,
	// Block 0x14, offset 0x500
	0x531: 0x12,
	// Block 0x15, offset 0x540
	0x546: 0x13,
	// Block 0x16, offset 0x580
	0x5ab: 0x14,
	// Block 0x17, offset 0x5c0
	0x5d4: 0x11,
	0x5fe: 0x11,
	// Block 0x18, offset 0x600
	0x618: 0x0a,
	// Block 0x19, offset 0x640
	0x65b: 0x15,
	// Block 0x1a, offset 0x680
	0x6a0: 0x16,
	// Block 0x1b, offset 0x6c0
	0x6d2: 0x17,
	0x6f6: 0x18,
	// Block 0x1c, offset 0x700
	0x711: 0x19,
	// Block 0x1d, offset 0x740
	0x768: 0x1a,
	// Block 0x1e, offset 0x780
	0x783: 0x1b,
	// Block 0x1f, offset 0x7c0
	0x7f9: 0x1c,
	// Block 0x20, offset 0x800
	0x831: 0x1d,
	// Block 0x21, offset 0x840
	0x85e: 0x1e,
	// Block 0x22, offset 0x880
	0x898: 0x1f,
	// Block 0x23, offset 0x8c0
	0x8c7: 0x18,
	0x8d5: 0x14,
	0x8f7: 0x20,
	0x8fe: 0x1f,
	// Block 0x24, offset 0x900
	0x905: 0x21,
	// Block 0x25, offset 0x940
	0x966: 0x03,
	// Block 0x26, offset 0x980
	0x981: 0x07, 0x983: 0x11,
	0x989: 0x12, 0x98a: 0x13, 0x98e: 0x14, 0x98f: 0x15,
	0x992: 0x16, 0x995: 0x17, 0x996: 0x18,
	0x998: 0x19, 0x999: 0x1a, 0x99b: 0x1b, 0x99f: 0x1c,
	0x9a3: 0x1d,
	0x9ad: 0x1e, 0x9af: 0x1f,
	0x9b0: 0x20, 0x9b1: 0x21,
	0x9b8: 0x22, 0x9bd: 0x23,
	// Block 0x27, offset 0x9c0
	0x9cd: 0x22,
	// Block 0x28, offset 0xa00
	0xa0c: 0x08,
	// Block 0x29, offset 0xa40
	0xa6f: 0x1c,
	// Block 0x2a, offset 0xa80
	0xa90: 0x1a,
	0xaaf: 0x23,
	// Block 0x2b, offset 0xac0
	0xae3: 0x19,
	0xae8: 0x24,
	0xafc: 0x25,
	// Block 0x2c, offset 0xb00
	0xb13: 0x26,
	// Block 0x2d, offset 0xb40
	0xb67: 0x1c,
	// Block 0x2e, offset 0xb80
	0xb8f: 0x0b,
	// Block 0x2f, offset 0xbc0
	0xbcb: 0x27,
	0xbe7: 0x26,
	// Block 0x30, offset 0xc00
	0xc34: 0x16,
	// Block 0x31, offset 0xc40
	0xc62: 0x03,
	// Block 0x32, offset 0xc80
	0xcbb: 0x12,
	// Block 0x33, offset 0xcc0
	0xcdf: 0x09,
	// Block 0x34, offset 0xd00
	0xd34: 0x0a,
	// Block 0x35, offset 0xd40
	0xd41: 0x1e,
	// Block 0x36, offset 0xd80
	0xd83: 0x28,
	// Block 0x37, offset 0xdc0
	0xdc0: 0x15,
	// Block 0x38, offset 0xe00
	0xe1a: 0x15,
	// Block 0x39, offset 0xe40
	0xe65: 0x29,
	// Block 0x3a, offset 0xe80
	0xe86: 0x1f,
	// Block 0x3b, offset 0xec0
	0xeec: 0x18,
	// Block 0x3c, offset 0xf00
	0xf28: 0x2a,
	// Block 0x3d, offset 0xf40
	0xf53: 0x08,
	// Block 0x3e, offset 0xf80
	0xfa2: 0x2b,
	0xfaa: 0x17,
	// Block 0x3f, offset 0xfc0
	0xfc0: 0x25, 0xfc2: 0x26,
	0xfc9: 0x27, 0xfcd: 0x28, 0xfce: 0x29,
	0xfd5: 0x2a,
	0xfd8: 0x2b, 0xfd9: 0x2c, 0xfdf: 0x2d,
	0xfe1: 0x2e, 0xfe2: 0x2f, 0xfe3: 0x30, 0xfe6: 0x31,
	0xfe9: 0x32, 0xfec: 0x33, 0xfed: 0x34, 0xfef: 0x35,
	0xff1: 0x36, 0xff2: 0x37, 0xff3: 0x38, 0xff4: 0x39,
	0xffa: 0x3a, 0xffc: 0x3b, 0xffe: 0x3c,
	// Block 0x40, offset 0x1000
	0x102c: 0x2c,
	// Block 0x41, offset 0x1040
	0x1074: 0x2c,
	// Block 0x42, offset 0x1080
	0x108c: 0x08,
	0x10a0: 0x2d,
	// Block 0x43, offset 0x10c0
	0x10e8: 0x10,
	// Block 0x44, offset 0x1100
	0x110f: 0x13,
	// Block 0x45, offset 0x1140
	0x114b: 0x2e,
	// Block 0x46, offset 0x1180
	0x118b: 0x23,
	0x119d: 0x0c,
	// Block 0x47, offset 0x11c0
	0x11c3: 0x12,
	0x11f9: 0x0f,
	// Block 0x48, offset 0x1200
	0x121e: 0x1b,
	// Block 0x49, offset 0x1240
	0x1270: 0x2f,
	// Block 0x4a, offset 0x1280
	0x128a: 0x1b,
	0x12a7: 0x02,
	// Block 0x4b, offset 0x12c0
	0x12fb: 0x14,
	// Block 0x4c, offset 0x1300
	0x1333: 0x30,
	// Block 0x4d, offset 0x1340
	0x134d: 0x31,
	// Block 0x4e, offset 0x1380
	0x138e: 0x15,
	// Block 0x4f, offset 0x13c0
	0x13f4: 0x32,
	// Block 0x50, offset 0x1400
	0x141b: 0x33,
	// Block 0x51, offset 0x1440
	0x1448: 0x3e, 0x1449: 0x3f, 0x144a: 0x40, 0x144f: 0x41,
	0x1459: 0x42, 0x145c: 0x43, 0x145e: 0x44, 0x145f: 0x45,
	0x1468: 0x46, 0x1469: 0x47, 0x146c: 0x48, 0x146d: 0x49, 0x146e: 0x4a,
	0x1472: 0x4b, 0x1473: 0x4c,
	0x1479: 0x4d, 0x147b: 0x4e,
	// Block 0x52, offset 0x1480
	0x1480: 0x34,
	0x1499: 0x11,
	0x14b6: 0x2c,
	// Block 0x53, offset 0x14c0
	0x14e4: 0x0d,
	// Block 0x54, offset 0x1500
	0x1527: 0x08,
	// Block 0x55, offset 0x1540
	0x1555: 0x2b,
	// Block 0x56, offset 0x1580
	0x15b2: 0x35,
	// Block 0x57, offset 0x15c0
	0x15f2: 0x1c, 0x15f4: 0x29,
	// Block 0x58, offset 0x1600
	0x1600: 0x50, 0x1603: 0x51,
	0x1608: 0x52, 0x160a: 0x53, 0x160d: 0x54, 0x160e: 0x55,
***REMOVED***

// lookup returns the trie value for the first UTF-8 encoding in s and
// the width in bytes of this encoding. The size will be 0 if s does not
// hold enough bytes to complete the encoding. len(s) must be greater than 0.
func (t *multiTrie) lookup(s []byte) (v uint64, sz int) ***REMOVED***
	c0 := s[0]
	switch ***REMOVED***
	case c0 < 0x80: // is ASCII
		return t.ascii[c0], 1
	case c0 < 0xC2:
		return 0, 1 // Illegal UTF-8: not a starter, not ASCII.
	case c0 < 0xE0: // 2-byte UTF-8
		if len(s) < 2 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := t.utf8Start[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c1), 2
	case c0 < 0xF0: // 3-byte UTF-8
		if len(s) < 3 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := t.utf8Start[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = multiIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c2), 3
	case c0 < 0xF8: // 4-byte UTF-8
		if len(s) < 4 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := t.utf8Start[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = multiIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o = uint32(i)<<6 + uint32(c2)
		i = multiIndex[o]
		c3 := s[3]
		if c3 < 0x80 || 0xC0 <= c3 ***REMOVED***
			return 0, 3 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c3), 4
	***REMOVED***
	// Illegal rune
	return 0, 1
***REMOVED***

// lookupUnsafe returns the trie value for the first UTF-8 encoding in s.
// s must start with a full and valid UTF-8 encoded rune.
func (t *multiTrie) lookupUnsafe(s []byte) uint64 ***REMOVED***
	c0 := s[0]
	if c0 < 0x80 ***REMOVED*** // is ASCII
		return t.ascii[c0]
	***REMOVED***
	i := t.utf8Start[c0]
	if c0 < 0xE0 ***REMOVED*** // 2-byte UTF-8
		return t.lookupValue(uint32(i), s[1])
	***REMOVED***
	i = multiIndex[uint32(i)<<6+uint32(s[1])]
	if c0 < 0xF0 ***REMOVED*** // 3-byte UTF-8
		return t.lookupValue(uint32(i), s[2])
	***REMOVED***
	i = multiIndex[uint32(i)<<6+uint32(s[2])]
	if c0 < 0xF8 ***REMOVED*** // 4-byte UTF-8
		return t.lookupValue(uint32(i), s[3])
	***REMOVED***
	return 0
***REMOVED***

// lookupString returns the trie value for the first UTF-8 encoding in s and
// the width in bytes of this encoding. The size will be 0 if s does not
// hold enough bytes to complete the encoding. len(s) must be greater than 0.
func (t *multiTrie) lookupString(s string) (v uint64, sz int) ***REMOVED***
	c0 := s[0]
	switch ***REMOVED***
	case c0 < 0x80: // is ASCII
		return t.ascii[c0], 1
	case c0 < 0xC2:
		return 0, 1 // Illegal UTF-8: not a starter, not ASCII.
	case c0 < 0xE0: // 2-byte UTF-8
		if len(s) < 2 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := t.utf8Start[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c1), 2
	case c0 < 0xF0: // 3-byte UTF-8
		if len(s) < 3 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := t.utf8Start[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = multiIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c2), 3
	case c0 < 0xF8: // 4-byte UTF-8
		if len(s) < 4 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := t.utf8Start[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = multiIndex[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o = uint32(i)<<6 + uint32(c2)
		i = multiIndex[o]
		c3 := s[3]
		if c3 < 0x80 || 0xC0 <= c3 ***REMOVED***
			return 0, 3 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c3), 4
	***REMOVED***
	// Illegal rune
	return 0, 1
***REMOVED***

// lookupStringUnsafe returns the trie value for the first UTF-8 encoding in s.
// s must start with a full and valid UTF-8 encoded rune.
func (t *multiTrie) lookupStringUnsafe(s string) uint64 ***REMOVED***
	c0 := s[0]
	if c0 < 0x80 ***REMOVED*** // is ASCII
		return t.ascii[c0]
	***REMOVED***
	i := t.utf8Start[c0]
	if c0 < 0xE0 ***REMOVED*** // 2-byte UTF-8
		return t.lookupValue(uint32(i), s[1])
	***REMOVED***
	i = multiIndex[uint32(i)<<6+uint32(s[1])]
	if c0 < 0xF0 ***REMOVED*** // 3-byte UTF-8
		return t.lookupValue(uint32(i), s[2])
	***REMOVED***
	i = multiIndex[uint32(i)<<6+uint32(s[2])]
	if c0 < 0xF8 ***REMOVED*** // 4-byte UTF-8
		return t.lookupValue(uint32(i), s[3])
	***REMOVED***
	return 0
***REMOVED***

// multiTrie. Total size: 18250 bytes (17.82 KiB). Checksum: a69a609d8696aa5e.
type multiTrie struct ***REMOVED***
	ascii     []uint64 // index for ASCII bytes
	utf8Start []uint8  // index for UTF-8 bytes >= 0xC0
***REMOVED***

func newMultiTrie(i int) *multiTrie ***REMOVED***
	h := multiTrieHandles[i]
	return &multiTrie***REMOVED***multiValues[uint32(h.ascii)<<6:], multiIndex[uint32(h.multi)<<6:]***REMOVED***
***REMOVED***

type multiTrieHandle struct ***REMOVED***
	ascii, multi uint8
***REMOVED***

// multiTrieHandles: 5 handles, 10 bytes
var multiTrieHandles = [5]multiTrieHandle***REMOVED***
	***REMOVED***0, 0***REMOVED***,   // 8c1e77823143d35c: all
	***REMOVED***0, 23***REMOVED***,  // 8fb58ff8243b45b0: ASCII only
	***REMOVED***0, 23***REMOVED***,  // 8fb58ff8243b45b0: ASCII only 2
	***REMOVED***0, 24***REMOVED***,  // 2ccc43994f11046f: BMP only
	***REMOVED***30, 25***REMOVED***, // ce448591bdcb4733: No BMP
***REMOVED***

// lookupValue determines the type of block n and looks up the value for b.
func (t *multiTrie) lookupValue(n uint32, b byte) uint64 ***REMOVED***
	switch ***REMOVED***
	default:
		return uint64(multiValues[n<<6+uint32(b)])
	***REMOVED***
***REMOVED***

// multiValues: 32 blocks, 2048 entries, 16384 bytes
// The third block is the zero block.
var multiValues = [2048]uint64***REMOVED***
	// Block 0x0, offset 0x0
	0x03: 0x6e361699800b9fb8, 0x04: 0x52d3935a34f6f0b, 0x05: 0x2948319393e7ef10,
	0x07: 0x20f03b006704f663, 0x08: 0x6c15c0732bb2495f, 0x09: 0xe54e2c59d953551,
	0x0f: 0x33d8a825807d8037, 0x10: 0x6ecd93cb12168b92, 0x11: 0x6a81c9c0ce86e884,
	0x1f: 0xa03e77aac8be79b, 0x20: 0x28591d0e7e486efa, 0x21: 0x716fa3bc398dec8,
	0x3f: 0x4fd3bcfa72bce8b0,
	// Block 0x1, offset 0x40
	0x40: 0x3cbaef3db8ba5f12, 0x41: 0x2d262347c1f56357,
	0x7f: 0x782caa2d25a418a9,
	// Block 0x2, offset 0x80
	// Block 0x3, offset 0xc0
	0xc0: 0x6bbd1f937b1ff5d2, 0xc1: 0x732e23088d2eb8a4,
	// Block 0x4, offset 0x100
	0x13f: 0x56f8c4c82f5962dc,
	// Block 0x5, offset 0x140
	0x140: 0x57dc4544729a5da2, 0x141: 0x2f62f9cd307ffa0d,
	// Block 0x6, offset 0x180
	0x1bf: 0x7bf4d0ebf302a088,
	// Block 0x7, offset 0x1c0
	0x1c0: 0x1f0d67f249e59931, 0x1c1: 0x3011def73aa550c7,
	// Block 0x8, offset 0x200
	0x23f: 0x5de81c1dff6bf29d,
	// Block 0x9, offset 0x240
	0x240: 0x752c035737b825e8, 0x241: 0x1e793399081e3bb3,
	// Block 0xa, offset 0x280
	0x2bf: 0x6a28f01979cbf059,
	// Block 0xb, offset 0x2c0
	0x2c0: 0x373a4b0f2cbd4c74, 0x2c1: 0x4fd2c288683b767c,
	// Block 0xc, offset 0x300
	0x33f: 0x5a10ffa9e29184fb,
	// Block 0xd, offset 0x340
	0x340: 0x700f9bdb53fff6a5, 0x341: 0xcde93df0427eb79,
	// Block 0xe, offset 0x380
	0x3bf: 0x74071288fff39c76,
	// Block 0xf, offset 0x3c0
	0x3c0: 0x481fc2f510e5268a, 0x3c1: 0x7565c28164204849,
	// Block 0x10, offset 0x400
	0x43f: 0x5676a62fd49c6bec,
	// Block 0x11, offset 0x440
	0x440: 0x2f2d15776cbafc6b, 0x441: 0x4c55e8dc0ff11a3f,
	// Block 0x12, offset 0x480
	0x4bf: 0x69d6f0fe711fafc9,
	// Block 0x13, offset 0x4c0
	0x4c0: 0x33181de28cfb062d, 0x4c1: 0x2ef3adc6bb2f2d02,
	// Block 0x14, offset 0x500
	0x53f: 0xe03b31814c95f8b,
	// Block 0x15, offset 0x540
	0x540: 0x3bf6dc9a1c115603, 0x541: 0x6984ec9b7f51f7fc,
	// Block 0x16, offset 0x580
	0x5bf: 0x3c02ea92fb168559,
	// Block 0x17, offset 0x5c0
	0x5c0: 0x1badfe42e7629494, 0x5c1: 0x6dc4a554005f7645,
	// Block 0x18, offset 0x600
	0x63f: 0x3bb2ed2a72748f4b,
	// Block 0x19, offset 0x640
	0x640: 0x291354cd6767ec10, 0x641: 0x2c3a4715e3c070d6,
	// Block 0x1a, offset 0x680
	0x6bf: 0x352711cfb7236418,
	// Block 0x1b, offset 0x6c0
	0x6c0: 0x3a59d34fb8bceda, 0x6c1: 0x5e90d8ebedd64fa1,
	// Block 0x1c, offset 0x700
	0x73f: 0x7191a77b28d23110,
	// Block 0x1d, offset 0x740
	0x740: 0x4ca7f0c1623423d8, 0x741: 0x4f7156d996e2d0de,
	// Block 0x1e, offset 0x780
	// Block 0x1f, offset 0x7c0
***REMOVED***

// multiIndex: 29 blocks, 1856 entries, 1856 bytes
// Block 0 is the zero block.
var multiIndex = [1856]uint8***REMOVED***
	// Block 0x0, offset 0x0
	// Block 0x1, offset 0x40
	// Block 0x2, offset 0x80
	// Block 0x3, offset 0xc0
	0xc2: 0x01, 0xc3: 0x02, 0xc4: 0x03, 0xc7: 0x04,
	0xc8: 0x05, 0xcf: 0x06,
	0xd0: 0x07,
	0xdf: 0x08,
	0xe0: 0x02, 0xe1: 0x03, 0xe2: 0x04, 0xe3: 0x05, 0xe4: 0x06, 0xe7: 0x07,
	0xe8: 0x08, 0xef: 0x09,
	0xf0: 0x0e, 0xf1: 0x11, 0xf2: 0x13, 0xf3: 0x15, 0xf4: 0x17,
	// Block 0x4, offset 0x100
	0x120: 0x09,
	0x13f: 0x0a,
	// Block 0x5, offset 0x140
	0x140: 0x0b,
	0x17f: 0x0c,
	// Block 0x6, offset 0x180
	0x180: 0x0d,
	// Block 0x7, offset 0x1c0
	0x1ff: 0x0e,
	// Block 0x8, offset 0x200
	0x200: 0x0f,
	// Block 0x9, offset 0x240
	0x27f: 0x10,
	// Block 0xa, offset 0x280
	0x280: 0x11,
	// Block 0xb, offset 0x2c0
	0x2ff: 0x12,
	// Block 0xc, offset 0x300
	0x300: 0x13,
	// Block 0xd, offset 0x340
	0x37f: 0x14,
	// Block 0xe, offset 0x380
	0x380: 0x15,
	// Block 0xf, offset 0x3c0
	0x3ff: 0x16,
	// Block 0x10, offset 0x400
	0x410: 0x0a,
	0x41f: 0x0b,
	0x420: 0x0c,
	0x43f: 0x0d,
	// Block 0x11, offset 0x440
	0x440: 0x17,
	// Block 0x12, offset 0x480
	0x4bf: 0x18,
	// Block 0x13, offset 0x4c0
	0x4c0: 0x0f,
	0x4ff: 0x10,
	// Block 0x14, offset 0x500
	0x500: 0x19,
	// Block 0x15, offset 0x540
	0x540: 0x12,
	// Block 0x16, offset 0x580
	0x5bf: 0x1a,
	// Block 0x17, offset 0x5c0
	0x5ff: 0x14,
	// Block 0x18, offset 0x600
	0x600: 0x1b,
	// Block 0x19, offset 0x640
	0x640: 0x16,
	// Block 0x1a, offset 0x680
	// Block 0x1b, offset 0x6c0
	0x6c2: 0x01, 0x6c3: 0x02, 0x6c4: 0x03, 0x6c7: 0x04,
	0x6c8: 0x05, 0x6cf: 0x06,
	0x6d0: 0x07,
	0x6df: 0x08,
	0x6e0: 0x02, 0x6e1: 0x03, 0x6e2: 0x04, 0x6e3: 0x05, 0x6e4: 0x06, 0x6e7: 0x07,
	0x6e8: 0x08, 0x6ef: 0x09,
	// Block 0x1c, offset 0x700
	0x730: 0x0e, 0x731: 0x11, 0x732: 0x13, 0x733: 0x15, 0x734: 0x17,
***REMOVED***
