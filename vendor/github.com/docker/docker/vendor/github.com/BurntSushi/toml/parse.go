package toml

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type parser struct ***REMOVED***
	mapping map[string]interface***REMOVED******REMOVED***
	types   map[string]tomlType
	lx      *lexer

	// A list of keys in the order that they appear in the TOML data.
	ordered []Key

	// the full key for the current hash in scope
	context Key

	// the base key name for everything except hashes
	currentKey string

	// rough approximation of line number
	approxLine int

	// A map of 'key.group.names' to whether they were created implicitly.
	implicits map[string]bool
***REMOVED***

type parseError string

func (pe parseError) Error() string ***REMOVED***
	return string(pe)
***REMOVED***

func parse(data string) (p *parser, err error) ***REMOVED***
	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			var ok bool
			if err, ok = r.(parseError); ok ***REMOVED***
				return
			***REMOVED***
			panic(r)
		***REMOVED***
	***REMOVED***()

	p = &parser***REMOVED***
		mapping:   make(map[string]interface***REMOVED******REMOVED***),
		types:     make(map[string]tomlType),
		lx:        lex(data),
		ordered:   make([]Key, 0),
		implicits: make(map[string]bool),
	***REMOVED***
	for ***REMOVED***
		item := p.next()
		if item.typ == itemEOF ***REMOVED***
			break
		***REMOVED***
		p.topLevel(item)
	***REMOVED***

	return p, nil
***REMOVED***

func (p *parser) panicf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	msg := fmt.Sprintf("Near line %d (last key parsed '%s'): %s",
		p.approxLine, p.current(), fmt.Sprintf(format, v...))
	panic(parseError(msg))
***REMOVED***

func (p *parser) next() item ***REMOVED***
	it := p.lx.nextItem()
	if it.typ == itemError ***REMOVED***
		p.panicf("%s", it.val)
	***REMOVED***
	return it
***REMOVED***

func (p *parser) bug(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.Fatalf("BUG: %s\n\n", fmt.Sprintf(format, v...))
***REMOVED***

func (p *parser) expect(typ itemType) item ***REMOVED***
	it := p.next()
	p.assertEqual(typ, it.typ)
	return it
***REMOVED***

func (p *parser) assertEqual(expected, got itemType) ***REMOVED***
	if expected != got ***REMOVED***
		p.bug("Expected '%s' but got '%s'.", expected, got)
	***REMOVED***
***REMOVED***

func (p *parser) topLevel(item item) ***REMOVED***
	switch item.typ ***REMOVED***
	case itemCommentStart:
		p.approxLine = item.line
		p.expect(itemText)
	case itemTableStart:
		kg := p.next()
		p.approxLine = kg.line

		var key Key
		for ; kg.typ != itemTableEnd && kg.typ != itemEOF; kg = p.next() ***REMOVED***
			key = append(key, p.keyString(kg))
		***REMOVED***
		p.assertEqual(itemTableEnd, kg.typ)

		p.establishContext(key, false)
		p.setType("", tomlHash)
		p.ordered = append(p.ordered, key)
	case itemArrayTableStart:
		kg := p.next()
		p.approxLine = kg.line

		var key Key
		for ; kg.typ != itemArrayTableEnd && kg.typ != itemEOF; kg = p.next() ***REMOVED***
			key = append(key, p.keyString(kg))
		***REMOVED***
		p.assertEqual(itemArrayTableEnd, kg.typ)

		p.establishContext(key, true)
		p.setType("", tomlArrayHash)
		p.ordered = append(p.ordered, key)
	case itemKeyStart:
		kname := p.next()
		p.approxLine = kname.line
		p.currentKey = p.keyString(kname)

		val, typ := p.value(p.next())
		p.setValue(p.currentKey, val)
		p.setType(p.currentKey, typ)
		p.ordered = append(p.ordered, p.context.add(p.currentKey))
		p.currentKey = ""
	default:
		p.bug("Unexpected type at top level: %s", item.typ)
	***REMOVED***
***REMOVED***

// Gets a string for a key (or part of a key in a table name).
func (p *parser) keyString(it item) string ***REMOVED***
	switch it.typ ***REMOVED***
	case itemText:
		return it.val
	case itemString, itemMultilineString,
		itemRawString, itemRawMultilineString:
		s, _ := p.value(it)
		return s.(string)
	default:
		p.bug("Unexpected key type: %s", it.typ)
		panic("unreachable")
	***REMOVED***
***REMOVED***

// value translates an expected value from the lexer into a Go value wrapped
// as an empty interface.
func (p *parser) value(it item) (interface***REMOVED******REMOVED***, tomlType) ***REMOVED***
	switch it.typ ***REMOVED***
	case itemString:
		return p.replaceEscapes(it.val), p.typeOfPrimitive(it)
	case itemMultilineString:
		trimmed := stripFirstNewline(stripEscapedWhitespace(it.val))
		return p.replaceEscapes(trimmed), p.typeOfPrimitive(it)
	case itemRawString:
		return it.val, p.typeOfPrimitive(it)
	case itemRawMultilineString:
		return stripFirstNewline(it.val), p.typeOfPrimitive(it)
	case itemBool:
		switch it.val ***REMOVED***
		case "true":
			return true, p.typeOfPrimitive(it)
		case "false":
			return false, p.typeOfPrimitive(it)
		***REMOVED***
		p.bug("Expected boolean value, but got '%s'.", it.val)
	case itemInteger:
		num, err := strconv.ParseInt(it.val, 10, 64)
		if err != nil ***REMOVED***
			// See comment below for floats describing why we make a
			// distinction between a bug and a user error.
			if e, ok := err.(*strconv.NumError); ok &&
				e.Err == strconv.ErrRange ***REMOVED***

				p.panicf("Integer '%s' is out of the range of 64-bit "+
					"signed integers.", it.val)
			***REMOVED*** else ***REMOVED***
				p.bug("Expected integer value, but got '%s'.", it.val)
			***REMOVED***
		***REMOVED***
		return num, p.typeOfPrimitive(it)
	case itemFloat:
		num, err := strconv.ParseFloat(it.val, 64)
		if err != nil ***REMOVED***
			// Distinguish float values. Normally, it'd be a bug if the lexer
			// provides an invalid float, but it's possible that the float is
			// out of range of valid values (which the lexer cannot determine).
			// So mark the former as a bug but the latter as a legitimate user
			// error.
			//
			// This is also true for integers.
			if e, ok := err.(*strconv.NumError); ok &&
				e.Err == strconv.ErrRange ***REMOVED***

				p.panicf("Float '%s' is out of the range of 64-bit "+
					"IEEE-754 floating-point numbers.", it.val)
			***REMOVED*** else ***REMOVED***
				p.bug("Expected float value, but got '%s'.", it.val)
			***REMOVED***
		***REMOVED***
		return num, p.typeOfPrimitive(it)
	case itemDatetime:
		t, err := time.Parse("2006-01-02T15:04:05Z", it.val)
		if err != nil ***REMOVED***
			p.bug("Expected Zulu formatted DateTime, but got '%s'.", it.val)
		***REMOVED***
		return t, p.typeOfPrimitive(it)
	case itemArray:
		array := make([]interface***REMOVED******REMOVED***, 0)
		types := make([]tomlType, 0)

		for it = p.next(); it.typ != itemArrayEnd; it = p.next() ***REMOVED***
			if it.typ == itemCommentStart ***REMOVED***
				p.expect(itemText)
				continue
			***REMOVED***

			val, typ := p.value(it)
			array = append(array, val)
			types = append(types, typ)
		***REMOVED***
		return array, p.typeOfArray(types)
	***REMOVED***
	p.bug("Unexpected value type: %s", it.typ)
	panic("unreachable")
***REMOVED***

// establishContext sets the current context of the parser,
// where the context is either a hash or an array of hashes. Which one is
// set depends on the value of the `array` parameter.
//
// Establishing the context also makes sure that the key isn't a duplicate, and
// will create implicit hashes automatically.
func (p *parser) establishContext(key Key, array bool) ***REMOVED***
	var ok bool

	// Always start at the top level and drill down for our context.
	hashContext := p.mapping
	keyContext := make(Key, 0)

	// We only need implicit hashes for key[0:-1]
	for _, k := range key[0 : len(key)-1] ***REMOVED***
		_, ok = hashContext[k]
		keyContext = append(keyContext, k)

		// No key? Make an implicit hash and move on.
		if !ok ***REMOVED***
			p.addImplicit(keyContext)
			hashContext[k] = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***

		// If the hash context is actually an array of tables, then set
		// the hash context to the last element in that array.
		//
		// Otherwise, it better be a table, since this MUST be a key group (by
		// virtue of it not being the last element in a key).
		switch t := hashContext[k].(type) ***REMOVED***
		case []map[string]interface***REMOVED******REMOVED***:
			hashContext = t[len(t)-1]
		case map[string]interface***REMOVED******REMOVED***:
			hashContext = t
		default:
			p.panicf("Key '%s' was already created as a hash.", keyContext)
		***REMOVED***
	***REMOVED***

	p.context = keyContext
	if array ***REMOVED***
		// If this is the first element for this array, then allocate a new
		// list of tables for it.
		k := key[len(key)-1]
		if _, ok := hashContext[k]; !ok ***REMOVED***
			hashContext[k] = make([]map[string]interface***REMOVED******REMOVED***, 0, 5)
		***REMOVED***

		// Add a new table. But make sure the key hasn't already been used
		// for something else.
		if hash, ok := hashContext[k].([]map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
			hashContext[k] = append(hash, make(map[string]interface***REMOVED******REMOVED***))
		***REMOVED*** else ***REMOVED***
			p.panicf("Key '%s' was already created and cannot be used as "+
				"an array.", keyContext)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		p.setValue(key[len(key)-1], make(map[string]interface***REMOVED******REMOVED***))
	***REMOVED***
	p.context = append(p.context, key[len(key)-1])
***REMOVED***

// setValue sets the given key to the given value in the current context.
// It will make sure that the key hasn't already been defined, account for
// implicit key groups.
func (p *parser) setValue(key string, value interface***REMOVED******REMOVED***) ***REMOVED***
	var tmpHash interface***REMOVED******REMOVED***
	var ok bool

	hash := p.mapping
	keyContext := make(Key, 0)
	for _, k := range p.context ***REMOVED***
		keyContext = append(keyContext, k)
		if tmpHash, ok = hash[k]; !ok ***REMOVED***
			p.bug("Context for key '%s' has not been established.", keyContext)
		***REMOVED***
		switch t := tmpHash.(type) ***REMOVED***
		case []map[string]interface***REMOVED******REMOVED***:
			// The context is a table of hashes. Pick the most recent table
			// defined as the current hash.
			hash = t[len(t)-1]
		case map[string]interface***REMOVED******REMOVED***:
			hash = t
		default:
			p.bug("Expected hash to have type 'map[string]interface***REMOVED******REMOVED***', but "+
				"it has '%T' instead.", tmpHash)
		***REMOVED***
	***REMOVED***
	keyContext = append(keyContext, key)

	if _, ok := hash[key]; ok ***REMOVED***
		// Typically, if the given key has already been set, then we have
		// to raise an error since duplicate keys are disallowed. However,
		// it's possible that a key was previously defined implicitly. In this
		// case, it is allowed to be redefined concretely. (See the
		// `tests/valid/implicit-and-explicit-after.toml` test in `toml-test`.)
		//
		// But we have to make sure to stop marking it as an implicit. (So that
		// another redefinition provokes an error.)
		//
		// Note that since it has already been defined (as a hash), we don't
		// want to overwrite it. So our business is done.
		if p.isImplicit(keyContext) ***REMOVED***
			p.removeImplicit(keyContext)
			return
		***REMOVED***

		// Otherwise, we have a concrete key trying to override a previous
		// key, which is *always* wrong.
		p.panicf("Key '%s' has already been defined.", keyContext)
	***REMOVED***
	hash[key] = value
***REMOVED***

// setType sets the type of a particular value at a given key.
// It should be called immediately AFTER setValue.
//
// Note that if `key` is empty, then the type given will be applied to the
// current context (which is either a table or an array of tables).
func (p *parser) setType(key string, typ tomlType) ***REMOVED***
	keyContext := make(Key, 0, len(p.context)+1)
	for _, k := range p.context ***REMOVED***
		keyContext = append(keyContext, k)
	***REMOVED***
	if len(key) > 0 ***REMOVED*** // allow type setting for hashes
		keyContext = append(keyContext, key)
	***REMOVED***
	p.types[keyContext.String()] = typ
***REMOVED***

// addImplicit sets the given Key as having been created implicitly.
func (p *parser) addImplicit(key Key) ***REMOVED***
	p.implicits[key.String()] = true
***REMOVED***

// removeImplicit stops tagging the given key as having been implicitly
// created.
func (p *parser) removeImplicit(key Key) ***REMOVED***
	p.implicits[key.String()] = false
***REMOVED***

// isImplicit returns true if the key group pointed to by the key was created
// implicitly.
func (p *parser) isImplicit(key Key) bool ***REMOVED***
	return p.implicits[key.String()]
***REMOVED***

// current returns the full key name of the current context.
func (p *parser) current() string ***REMOVED***
	if len(p.currentKey) == 0 ***REMOVED***
		return p.context.String()
	***REMOVED***
	if len(p.context) == 0 ***REMOVED***
		return p.currentKey
	***REMOVED***
	return fmt.Sprintf("%s.%s", p.context, p.currentKey)
***REMOVED***

func stripFirstNewline(s string) string ***REMOVED***
	if len(s) == 0 || s[0] != '\n' ***REMOVED***
		return s
	***REMOVED***
	return s[1:len(s)]
***REMOVED***

func stripEscapedWhitespace(s string) string ***REMOVED***
	esc := strings.Split(s, "\\\n")
	if len(esc) > 1 ***REMOVED***
		for i := 1; i < len(esc); i++ ***REMOVED***
			esc[i] = strings.TrimLeftFunc(esc[i], unicode.IsSpace)
		***REMOVED***
	***REMOVED***
	return strings.Join(esc, "")
***REMOVED***

func (p *parser) replaceEscapes(str string) string ***REMOVED***
	var replaced []rune
	s := []byte(str)
	r := 0
	for r < len(s) ***REMOVED***
		if s[r] != '\\' ***REMOVED***
			c, size := utf8.DecodeRune(s[r:])
			r += size
			replaced = append(replaced, c)
			continue
		***REMOVED***
		r += 1
		if r >= len(s) ***REMOVED***
			p.bug("Escape sequence at end of string.")
			return ""
		***REMOVED***
		switch s[r] ***REMOVED***
		default:
			p.bug("Expected valid escape code after \\, but got %q.", s[r])
			return ""
		case 'b':
			replaced = append(replaced, rune(0x0008))
			r += 1
		case 't':
			replaced = append(replaced, rune(0x0009))
			r += 1
		case 'n':
			replaced = append(replaced, rune(0x000A))
			r += 1
		case 'f':
			replaced = append(replaced, rune(0x000C))
			r += 1
		case 'r':
			replaced = append(replaced, rune(0x000D))
			r += 1
		case '"':
			replaced = append(replaced, rune(0x0022))
			r += 1
		case '\\':
			replaced = append(replaced, rune(0x005C))
			r += 1
		case 'u':
			// At this point, we know we have a Unicode escape of the form
			// `uXXXX` at [r, r+5). (Because the lexer guarantees this
			// for us.)
			escaped := p.asciiEscapeToUnicode(s[r+1 : r+5])
			replaced = append(replaced, escaped)
			r += 5
		case 'U':
			// At this point, we know we have a Unicode escape of the form
			// `uXXXX` at [r, r+9). (Because the lexer guarantees this
			// for us.)
			escaped := p.asciiEscapeToUnicode(s[r+1 : r+9])
			replaced = append(replaced, escaped)
			r += 9
		***REMOVED***
	***REMOVED***
	return string(replaced)
***REMOVED***

func (p *parser) asciiEscapeToUnicode(bs []byte) rune ***REMOVED***
	s := string(bs)
	hex, err := strconv.ParseUint(strings.ToLower(s), 16, 32)
	if err != nil ***REMOVED***
		p.bug("Could not parse '%s' as a hexadecimal number, but the "+
			"lexer claims it's OK: %s", s, err)
	***REMOVED***

	// BUG(burntsushi)
	// I honestly don't understand how this works. I can't seem
	// to find a way to make this fail. I figured this would fail on invalid
	// UTF-8 characters like U+DCFF, but it doesn't.
	if !utf8.ValidString(string(rune(hex))) ***REMOVED***
		p.panicf("Escaped character '\\u%s' is not valid UTF-8.", s)
	***REMOVED***
	return rune(hex)
***REMOVED***

func isStringType(ty itemType) bool ***REMOVED***
	return ty == itemString || ty == itemMultilineString ||
		ty == itemRawString || ty == itemRawMultilineString
***REMOVED***
