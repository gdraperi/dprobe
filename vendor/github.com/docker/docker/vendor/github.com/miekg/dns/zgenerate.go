package dns

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Parse the $GENERATE statement as used in BIND9 zones.
// See http://www.zytrax.com/books/dns/ch8/generate.html for instance.
// We are called after '$GENERATE '. After which we expect:
// * the range (12-24/2)
// * lhs (ownername)
// * [[ttl][class]]
// * type
// * rhs (rdata)
// But we are lazy here, only the range is parsed *all* occurences
// of $ after that are interpreted.
// Any error are returned as a string value, the empty string signals
// "no error".
func generate(l lex, c chan lex, t chan *Token, o string) string ***REMOVED***
	step := 1
	if i := strings.IndexAny(l.token, "/"); i != -1 ***REMOVED***
		if i+1 == len(l.token) ***REMOVED***
			return "bad step in $GENERATE range"
		***REMOVED***
		if s, e := strconv.Atoi(l.token[i+1:]); e == nil ***REMOVED***
			if s < 0 ***REMOVED***
				return "bad step in $GENERATE range"
			***REMOVED***
			step = s
		***REMOVED*** else ***REMOVED***
			return "bad step in $GENERATE range"
		***REMOVED***
		l.token = l.token[:i]
	***REMOVED***
	sx := strings.SplitN(l.token, "-", 2)
	if len(sx) != 2 ***REMOVED***
		return "bad start-stop in $GENERATE range"
	***REMOVED***
	start, err := strconv.Atoi(sx[0])
	if err != nil ***REMOVED***
		return "bad start in $GENERATE range"
	***REMOVED***
	end, err := strconv.Atoi(sx[1])
	if err != nil ***REMOVED***
		return "bad stop in $GENERATE range"
	***REMOVED***
	if end < 0 || start < 0 || end < start ***REMOVED***
		return "bad range in $GENERATE range"
	***REMOVED***

	<-c // _BLANK
	// Create a complete new string, which we then parse again.
	s := ""
BuildRR:
	l = <-c
	if l.value != zNewline && l.value != zEOF ***REMOVED***
		s += l.token
		goto BuildRR
	***REMOVED***
	for i := start; i <= end; i += step ***REMOVED***
		var (
			escape bool
			dom    bytes.Buffer
			mod    string
			err    string
			offset int
		)

		for j := 0; j < len(s); j++ ***REMOVED*** // No 'range' because we need to jump around
			switch s[j] ***REMOVED***
			case '\\':
				if escape ***REMOVED***
					dom.WriteByte('\\')
					escape = false
					continue
				***REMOVED***
				escape = true
			case '$':
				mod = "%d"
				offset = 0
				if escape ***REMOVED***
					dom.WriteByte('$')
					escape = false
					continue
				***REMOVED***
				escape = false
				if j+1 >= len(s) ***REMOVED*** // End of the string
					dom.WriteString(fmt.Sprintf(mod, i+offset))
					continue
				***REMOVED*** else ***REMOVED***
					if s[j+1] == '$' ***REMOVED***
						dom.WriteByte('$')
						j++
						continue
					***REMOVED***
				***REMOVED***
				// Search for ***REMOVED*** and ***REMOVED***
				if s[j+1] == '***REMOVED***' ***REMOVED*** // Modifier block
					sep := strings.Index(s[j+2:], "***REMOVED***")
					if sep == -1 ***REMOVED***
						return "bad modifier in $GENERATE"
					***REMOVED***
					mod, offset, err = modToPrintf(s[j+2 : j+2+sep])
					if err != "" ***REMOVED***
						return err
					***REMOVED***
					j += 2 + sep // Jump to it
				***REMOVED***
				dom.WriteString(fmt.Sprintf(mod, i+offset))
			default:
				if escape ***REMOVED*** // Pretty useless here
					escape = false
					continue
				***REMOVED***
				dom.WriteByte(s[j])
			***REMOVED***
		***REMOVED***
		// Re-parse the RR and send it on the current channel t
		rx, e := NewRR("$ORIGIN " + o + "\n" + dom.String())
		if e != nil ***REMOVED***
			return e.(*ParseError).err
		***REMOVED***
		t <- &Token***REMOVED***RR: rx***REMOVED***
		// Its more efficient to first built the rrlist and then parse it in
		// one go! But is this a problem?
	***REMOVED***
	return ""
***REMOVED***

// Convert a $GENERATE modifier 0,0,d to something Printf can deal with.
func modToPrintf(s string) (string, int, string) ***REMOVED***
	xs := strings.SplitN(s, ",", 3)
	if len(xs) != 3 ***REMOVED***
		return "", 0, "bad modifier in $GENERATE"
	***REMOVED***
	// xs[0] is offset, xs[1] is width, xs[2] is base
	if xs[2] != "o" && xs[2] != "d" && xs[2] != "x" && xs[2] != "X" ***REMOVED***
		return "", 0, "bad base in $GENERATE"
	***REMOVED***
	offset, err := strconv.Atoi(xs[0])
	if err != nil || offset > 255 ***REMOVED***
		return "", 0, "bad offset in $GENERATE"
	***REMOVED***
	width, err := strconv.Atoi(xs[1])
	if err != nil || width > 255 ***REMOVED***
		return "", offset, "bad width in $GENERATE"
	***REMOVED***
	switch ***REMOVED***
	case width < 0:
		return "", offset, "bad width in $GENERATE"
	case width == 0:
		return "%" + xs[1] + xs[2], offset, ""
	***REMOVED***
	return "%0" + xs[1] + xs[2], offset, ""
***REMOVED***
