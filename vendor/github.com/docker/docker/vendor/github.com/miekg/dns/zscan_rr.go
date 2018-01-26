package dns

import (
	"encoding/base64"
	"net"
	"strconv"
	"strings"
)

type parserFunc struct ***REMOVED***
	// Func defines the function that parses the tokens and returns the RR
	// or an error. The last string contains any comments in the line as
	// they returned by the lexer as well.
	Func func(h RR_Header, c chan lex, origin string, file string) (RR, *ParseError, string)
	// Signals if the RR ending is of variable length, like TXT or records
	// that have Hexadecimal or Base64 as their last element in the Rdata. Records
	// that have a fixed ending or for instance A, AAAA, SOA and etc.
	Variable bool
***REMOVED***

// Parse the rdata of each rrtype.
// All data from the channel c is either zString or zBlank.
// After the rdata there may come a zBlank and then a zNewline
// or immediately a zNewline. If this is not the case we flag
// an *ParseError: garbage after rdata.
func setRR(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	parserfunc, ok := typeToparserFunc[h.Rrtype]
	if ok ***REMOVED***
		r, e, cm := parserfunc.Func(h, c, o, f)
		if parserfunc.Variable ***REMOVED***
			return r, e, cm
		***REMOVED***
		if e != nil ***REMOVED***
			return nil, e, ""
		***REMOVED***
		e, cm = slurpRemainder(c, f)
		if e != nil ***REMOVED***
			return nil, e, ""
		***REMOVED***
		return r, nil, cm
	***REMOVED***
	// RFC3957 RR (Unknown RR handling)
	return setRFC3597(h, c, o, f)
***REMOVED***

// A remainder of the rdata with embedded spaces, return the parsed string (sans the spaces)
// or an error
func endingToString(c chan lex, errstr, f string) (string, *ParseError, string) ***REMOVED***
	s := ""
	l := <-c // zString
	for l.value != zNewline && l.value != zEOF ***REMOVED***
		if l.err ***REMOVED***
			return s, &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
		***REMOVED***
		switch l.value ***REMOVED***
		case zString:
			s += l.token
		case zBlank: // Ok
		default:
			return "", &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
		***REMOVED***
		l = <-c
	***REMOVED***
	return s, nil, l.comment
***REMOVED***

// A remainder of the rdata with embedded spaces, return the parsed string slice (sans the spaces)
// or an error
func endingToTxtSlice(c chan lex, errstr, f string) ([]string, *ParseError, string) ***REMOVED***
	// Get the remaining data until we see a zNewline
	quote := false
	l := <-c
	var s []string
	if l.err ***REMOVED***
		return s, &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
	***REMOVED***
	switch l.value == zQuote ***REMOVED***
	case true: // A number of quoted string
		s = make([]string, 0)
		empty := true
		for l.value != zNewline && l.value != zEOF ***REMOVED***
			if l.err ***REMOVED***
				return nil, &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
			***REMOVED***
			switch l.value ***REMOVED***
			case zString:
				empty = false
				if len(l.token) > 255 ***REMOVED***
					// split up tokens that are larger than 255 into 255-chunks
					sx := []string***REMOVED******REMOVED***
					p, i := 0, 255
					for ***REMOVED***
						if i <= len(l.token) ***REMOVED***
							sx = append(sx, l.token[p:i])
						***REMOVED*** else ***REMOVED***
							sx = append(sx, l.token[p:])
							break

						***REMOVED***
						p, i = p+255, i+255
					***REMOVED***
					s = append(s, sx...)
					break
				***REMOVED***

				s = append(s, l.token)
			case zBlank:
				if quote ***REMOVED***
					// zBlank can only be seen in between txt parts.
					return nil, &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
				***REMOVED***
			case zQuote:
				if empty && quote ***REMOVED***
					s = append(s, "")
				***REMOVED***
				quote = !quote
				empty = true
			default:
				return nil, &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
			***REMOVED***
			l = <-c
		***REMOVED***
		if quote ***REMOVED***
			return nil, &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
		***REMOVED***
	case false: // Unquoted text record
		s = make([]string, 1)
		for l.value != zNewline && l.value != zEOF ***REMOVED***
			if l.err ***REMOVED***
				return s, &ParseError***REMOVED***f, errstr, l***REMOVED***, ""
			***REMOVED***
			s[0] += l.token
			l = <-c
		***REMOVED***
	***REMOVED***
	return s, nil, l.comment
***REMOVED***

func setA(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(A)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED*** // Dynamic updates.
		return rr, nil, ""
	***REMOVED***
	rr.A = net.ParseIP(l.token)
	if rr.A == nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad A A", l***REMOVED***, ""
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setAAAA(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(AAAA)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	rr.AAAA = net.ParseIP(l.token)
	if rr.AAAA == nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad AAAA AAAA", l***REMOVED***, ""
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setNS(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NS)
	rr.Hdr = h

	l := <-c
	rr.Ns = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Ns = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NS Ns", l***REMOVED***, ""
	***REMOVED***
	if rr.Ns[l.length-1] != '.' ***REMOVED***
		rr.Ns = appendOrigin(rr.Ns, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setPTR(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(PTR)
	rr.Hdr = h

	l := <-c
	rr.Ptr = l.token
	if l.length == 0 ***REMOVED*** // dynamic update rr.
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Ptr = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad PTR Ptr", l***REMOVED***, ""
	***REMOVED***
	if rr.Ptr[l.length-1] != '.' ***REMOVED***
		rr.Ptr = appendOrigin(rr.Ptr, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setNSAPPTR(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NSAPPTR)
	rr.Hdr = h

	l := <-c
	rr.Ptr = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Ptr = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSAP-PTR Ptr", l***REMOVED***, ""
	***REMOVED***
	if rr.Ptr[l.length-1] != '.' ***REMOVED***
		rr.Ptr = appendOrigin(rr.Ptr, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setRP(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(RP)
	rr.Hdr = h

	l := <-c
	rr.Mbox = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Mbox = o
	***REMOVED*** else ***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad RP Mbox", l***REMOVED***, ""
		***REMOVED***
		if rr.Mbox[l.length-1] != '.' ***REMOVED***
			rr.Mbox = appendOrigin(rr.Mbox, o)
		***REMOVED***
	***REMOVED***
	<-c // zBlank
	l = <-c
	rr.Txt = l.token
	if l.token == "@" ***REMOVED***
		rr.Txt = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RP Txt", l***REMOVED***, ""
	***REMOVED***
	if rr.Txt[l.length-1] != '.' ***REMOVED***
		rr.Txt = appendOrigin(rr.Txt, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setMR(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(MR)
	rr.Hdr = h

	l := <-c
	rr.Mr = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Mr = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MR Mr", l***REMOVED***, ""
	***REMOVED***
	if rr.Mr[l.length-1] != '.' ***REMOVED***
		rr.Mr = appendOrigin(rr.Mr, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setMB(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(MB)
	rr.Hdr = h

	l := <-c
	rr.Mb = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Mb = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MB Mb", l***REMOVED***, ""
	***REMOVED***
	if rr.Mb[l.length-1] != '.' ***REMOVED***
		rr.Mb = appendOrigin(rr.Mb, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setMG(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(MG)
	rr.Hdr = h

	l := <-c
	rr.Mg = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Mg = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MG Mg", l***REMOVED***, ""
	***REMOVED***
	if rr.Mg[l.length-1] != '.' ***REMOVED***
		rr.Mg = appendOrigin(rr.Mg, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setHINFO(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(HINFO)
	rr.Hdr = h

	chunks, e, c1 := endingToTxtSlice(c, "bad HINFO Fields", f)
	if e != nil ***REMOVED***
		return nil, e, c1
	***REMOVED***

	if ln := len(chunks); ln == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED*** else if ln == 1 ***REMOVED***
		// Can we split it?
		if out := strings.Fields(chunks[0]); len(out) > 1 ***REMOVED***
			chunks = out
		***REMOVED*** else ***REMOVED***
			chunks = append(chunks, "")
		***REMOVED***
	***REMOVED***

	rr.Cpu = chunks[0]
	rr.Os = strings.Join(chunks[1:], " ")

	return rr, nil, ""
***REMOVED***

func setMINFO(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(MINFO)
	rr.Hdr = h

	l := <-c
	rr.Rmail = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Rmail = o
	***REMOVED*** else ***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad MINFO Rmail", l***REMOVED***, ""
		***REMOVED***
		if rr.Rmail[l.length-1] != '.' ***REMOVED***
			rr.Rmail = appendOrigin(rr.Rmail, o)
		***REMOVED***
	***REMOVED***
	<-c // zBlank
	l = <-c
	rr.Email = l.token
	if l.token == "@" ***REMOVED***
		rr.Email = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MINFO Email", l***REMOVED***, ""
	***REMOVED***
	if rr.Email[l.length-1] != '.' ***REMOVED***
		rr.Email = appendOrigin(rr.Email, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setMF(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(MF)
	rr.Hdr = h

	l := <-c
	rr.Mf = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Mf = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MF Mf", l***REMOVED***, ""
	***REMOVED***
	if rr.Mf[l.length-1] != '.' ***REMOVED***
		rr.Mf = appendOrigin(rr.Mf, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setMD(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(MD)
	rr.Hdr = h

	l := <-c
	rr.Md = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Md = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MD Md", l***REMOVED***, ""
	***REMOVED***
	if rr.Md[l.length-1] != '.' ***REMOVED***
		rr.Md = appendOrigin(rr.Md, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setMX(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(MX)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MX Pref", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Mx = l.token
	if l.token == "@" ***REMOVED***
		rr.Mx = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad MX Mx", l***REMOVED***, ""
	***REMOVED***
	if rr.Mx[l.length-1] != '.' ***REMOVED***
		rr.Mx = appendOrigin(rr.Mx, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setRT(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(RT)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RT Preference", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Host = l.token
	if l.token == "@" ***REMOVED***
		rr.Host = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RT Host", l***REMOVED***, ""
	***REMOVED***
	if rr.Host[l.length-1] != '.' ***REMOVED***
		rr.Host = appendOrigin(rr.Host, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setAFSDB(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(AFSDB)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad AFSDB Subtype", l***REMOVED***, ""
	***REMOVED***
	rr.Subtype = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Hostname = l.token
	if l.token == "@" ***REMOVED***
		rr.Hostname = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad AFSDB Hostname", l***REMOVED***, ""
	***REMOVED***
	if rr.Hostname[l.length-1] != '.' ***REMOVED***
		rr.Hostname = appendOrigin(rr.Hostname, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setX25(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(X25)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad X25 PSDNAddress", l***REMOVED***, ""
	***REMOVED***
	rr.PSDNAddress = l.token
	return rr, nil, ""
***REMOVED***

func setKX(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(KX)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad KX Pref", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Exchanger = l.token
	if l.token == "@" ***REMOVED***
		rr.Exchanger = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad KX Exchanger", l***REMOVED***, ""
	***REMOVED***
	if rr.Exchanger[l.length-1] != '.' ***REMOVED***
		rr.Exchanger = appendOrigin(rr.Exchanger, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setCNAME(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(CNAME)
	rr.Hdr = h

	l := <-c
	rr.Target = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Target = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CNAME Target", l***REMOVED***, ""
	***REMOVED***
	if rr.Target[l.length-1] != '.' ***REMOVED***
		rr.Target = appendOrigin(rr.Target, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setDNAME(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(DNAME)
	rr.Hdr = h

	l := <-c
	rr.Target = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Target = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CNAME Target", l***REMOVED***, ""
	***REMOVED***
	if rr.Target[l.length-1] != '.' ***REMOVED***
		rr.Target = appendOrigin(rr.Target, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setSOA(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(SOA)
	rr.Hdr = h

	l := <-c
	rr.Ns = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	<-c // zBlank
	if l.token == "@" ***REMOVED***
		rr.Ns = o
	***REMOVED*** else ***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad SOA Ns", l***REMOVED***, ""
		***REMOVED***
		if rr.Ns[l.length-1] != '.' ***REMOVED***
			rr.Ns = appendOrigin(rr.Ns, o)
		***REMOVED***
	***REMOVED***

	l = <-c
	rr.Mbox = l.token
	if l.token == "@" ***REMOVED***
		rr.Mbox = o
	***REMOVED*** else ***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad SOA Mbox", l***REMOVED***, ""
		***REMOVED***
		if rr.Mbox[l.length-1] != '.' ***REMOVED***
			rr.Mbox = appendOrigin(rr.Mbox, o)
		***REMOVED***
	***REMOVED***
	<-c // zBlank

	var (
		v  uint32
		ok bool
	)
	for i := 0; i < 5; i++ ***REMOVED***
		l = <-c
		if l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad SOA zone parameter", l***REMOVED***, ""
		***REMOVED***
		if j, e := strconv.Atoi(l.token); e != nil ***REMOVED***
			if i == 0 ***REMOVED***
				// Serial should be a number
				return nil, &ParseError***REMOVED***f, "bad SOA zone parameter", l***REMOVED***, ""
			***REMOVED***
			if v, ok = stringToTtl(l.token); !ok ***REMOVED***
				return nil, &ParseError***REMOVED***f, "bad SOA zone parameter", l***REMOVED***, ""

			***REMOVED***
		***REMOVED*** else ***REMOVED***
			v = uint32(j)
		***REMOVED***
		switch i ***REMOVED***
		case 0:
			rr.Serial = v
			<-c // zBlank
		case 1:
			rr.Refresh = v
			<-c // zBlank
		case 2:
			rr.Retry = v
			<-c // zBlank
		case 3:
			rr.Expire = v
			<-c // zBlank
		case 4:
			rr.Minttl = v
		***REMOVED***
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setSRV(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(SRV)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad SRV Priority", l***REMOVED***, ""
	***REMOVED***
	rr.Priority = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad SRV Weight", l***REMOVED***, ""
	***REMOVED***
	rr.Weight = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad SRV Port", l***REMOVED***, ""
	***REMOVED***
	rr.Port = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Target = l.token
	if l.token == "@" ***REMOVED***
		rr.Target = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad SRV Target", l***REMOVED***, ""
	***REMOVED***
	if rr.Target[l.length-1] != '.' ***REMOVED***
		rr.Target = appendOrigin(rr.Target, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setNAPTR(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NAPTR)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Order", l***REMOVED***, ""
	***REMOVED***
	rr.Order = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Preference", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	// Flags
	<-c     // zBlank
	l = <-c // _QUOTE
	if l.value != zQuote ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Flags", l***REMOVED***, ""
	***REMOVED***
	l = <-c // Either String or Quote
	if l.value == zString ***REMOVED***
		rr.Flags = l.token
		l = <-c // _QUOTE
		if l.value != zQuote ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad NAPTR Flags", l***REMOVED***, ""
		***REMOVED***
	***REMOVED*** else if l.value == zQuote ***REMOVED***
		rr.Flags = ""
	***REMOVED*** else ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Flags", l***REMOVED***, ""
	***REMOVED***

	// Service
	<-c     // zBlank
	l = <-c // _QUOTE
	if l.value != zQuote ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Service", l***REMOVED***, ""
	***REMOVED***
	l = <-c // Either String or Quote
	if l.value == zString ***REMOVED***
		rr.Service = l.token
		l = <-c // _QUOTE
		if l.value != zQuote ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad NAPTR Service", l***REMOVED***, ""
		***REMOVED***
	***REMOVED*** else if l.value == zQuote ***REMOVED***
		rr.Service = ""
	***REMOVED*** else ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Service", l***REMOVED***, ""
	***REMOVED***

	// Regexp
	<-c     // zBlank
	l = <-c // _QUOTE
	if l.value != zQuote ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Regexp", l***REMOVED***, ""
	***REMOVED***
	l = <-c // Either String or Quote
	if l.value == zString ***REMOVED***
		rr.Regexp = l.token
		l = <-c // _QUOTE
		if l.value != zQuote ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad NAPTR Regexp", l***REMOVED***, ""
		***REMOVED***
	***REMOVED*** else if l.value == zQuote ***REMOVED***
		rr.Regexp = ""
	***REMOVED*** else ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Regexp", l***REMOVED***, ""
	***REMOVED***
	// After quote no space??
	<-c     // zBlank
	l = <-c // zString
	rr.Replacement = l.token
	if l.token == "@" ***REMOVED***
		rr.Replacement = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NAPTR Replacement", l***REMOVED***, ""
	***REMOVED***
	if rr.Replacement[l.length-1] != '.' ***REMOVED***
		rr.Replacement = appendOrigin(rr.Replacement, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setTALINK(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(TALINK)
	rr.Hdr = h

	l := <-c
	rr.PreviousName = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.PreviousName = o
	***REMOVED*** else ***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad TALINK PreviousName", l***REMOVED***, ""
		***REMOVED***
		if rr.PreviousName[l.length-1] != '.' ***REMOVED***
			rr.PreviousName = appendOrigin(rr.PreviousName, o)
		***REMOVED***
	***REMOVED***
	<-c // zBlank
	l = <-c
	rr.NextName = l.token
	if l.token == "@" ***REMOVED***
		rr.NextName = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad TALINK NextName", l***REMOVED***, ""
	***REMOVED***
	if rr.NextName[l.length-1] != '.' ***REMOVED***
		rr.NextName = appendOrigin(rr.NextName, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setLOC(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(LOC)
	rr.Hdr = h
	// Non zero defaults for LOC record, see RFC 1876, Section 3.
	rr.HorizPre = 165 // 10000
	rr.VertPre = 162  // 10
	rr.Size = 18      // 1
	ok := false
	// North
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Latitude", l***REMOVED***, ""
	***REMOVED***
	rr.Latitude = 1000 * 60 * 60 * uint32(i)

	<-c // zBlank
	// Either number, 'N' or 'S'
	l = <-c
	if rr.Latitude, ok = locCheckNorth(l.token, rr.Latitude); ok ***REMOVED***
		goto East
	***REMOVED***
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Latitude minutes", l***REMOVED***, ""
	***REMOVED***
	rr.Latitude += 1000 * 60 * uint32(i)

	<-c // zBlank
	l = <-c
	if i, e := strconv.ParseFloat(l.token, 32); e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Latitude seconds", l***REMOVED***, ""
	***REMOVED*** else ***REMOVED***
		rr.Latitude += uint32(1000 * i)
	***REMOVED***
	<-c // zBlank
	// Either number, 'N' or 'S'
	l = <-c
	if rr.Latitude, ok = locCheckNorth(l.token, rr.Latitude); ok ***REMOVED***
		goto East
	***REMOVED***
	// If still alive, flag an error
	return nil, &ParseError***REMOVED***f, "bad LOC Latitude North/South", l***REMOVED***, ""

East:
	// East
	<-c // zBlank
	l = <-c
	if i, e := strconv.Atoi(l.token); e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Longitude", l***REMOVED***, ""
	***REMOVED*** else ***REMOVED***
		rr.Longitude = 1000 * 60 * 60 * uint32(i)
	***REMOVED***
	<-c // zBlank
	// Either number, 'E' or 'W'
	l = <-c
	if rr.Longitude, ok = locCheckEast(l.token, rr.Longitude); ok ***REMOVED***
		goto Altitude
	***REMOVED***
	if i, e := strconv.Atoi(l.token); e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Longitude minutes", l***REMOVED***, ""
	***REMOVED*** else ***REMOVED***
		rr.Longitude += 1000 * 60 * uint32(i)
	***REMOVED***
	<-c // zBlank
	l = <-c
	if i, e := strconv.ParseFloat(l.token, 32); e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Longitude seconds", l***REMOVED***, ""
	***REMOVED*** else ***REMOVED***
		rr.Longitude += uint32(1000 * i)
	***REMOVED***
	<-c // zBlank
	// Either number, 'E' or 'W'
	l = <-c
	if rr.Longitude, ok = locCheckEast(l.token, rr.Longitude); ok ***REMOVED***
		goto Altitude
	***REMOVED***
	// If still alive, flag an error
	return nil, &ParseError***REMOVED***f, "bad LOC Longitude East/West", l***REMOVED***, ""

Altitude:
	<-c // zBlank
	l = <-c
	if l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Altitude", l***REMOVED***, ""
	***REMOVED***
	if l.token[len(l.token)-1] == 'M' || l.token[len(l.token)-1] == 'm' ***REMOVED***
		l.token = l.token[0 : len(l.token)-1]
	***REMOVED***
	if i, e := strconv.ParseFloat(l.token, 32); e != nil ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LOC Altitude", l***REMOVED***, ""
	***REMOVED*** else ***REMOVED***
		rr.Altitude = uint32(i*100.0 + 10000000.0 + 0.5)
	***REMOVED***

	// And now optionally the other values
	l = <-c
	count := 0
	for l.value != zNewline && l.value != zEOF ***REMOVED***
		switch l.value ***REMOVED***
		case zString:
			switch count ***REMOVED***
			case 0: // Size
				e, m, ok := stringToCm(l.token)
				if !ok ***REMOVED***
					return nil, &ParseError***REMOVED***f, "bad LOC Size", l***REMOVED***, ""
				***REMOVED***
				rr.Size = (e & 0x0f) | (m << 4 & 0xf0)
			case 1: // HorizPre
				e, m, ok := stringToCm(l.token)
				if !ok ***REMOVED***
					return nil, &ParseError***REMOVED***f, "bad LOC HorizPre", l***REMOVED***, ""
				***REMOVED***
				rr.HorizPre = (e & 0x0f) | (m << 4 & 0xf0)
			case 2: // VertPre
				e, m, ok := stringToCm(l.token)
				if !ok ***REMOVED***
					return nil, &ParseError***REMOVED***f, "bad LOC VertPre", l***REMOVED***, ""
				***REMOVED***
				rr.VertPre = (e & 0x0f) | (m << 4 & 0xf0)
			***REMOVED***
			count++
		case zBlank:
			// Ok
		default:
			return nil, &ParseError***REMOVED***f, "bad LOC Size, HorizPre or VertPre", l***REMOVED***, ""
		***REMOVED***
		l = <-c
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setHIP(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(HIP)
	rr.Hdr = h

	// HitLength is not represented
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad HIP PublicKeyAlgorithm", l***REMOVED***, ""
	***REMOVED***
	rr.PublicKeyAlgorithm = uint8(i)
	<-c     // zBlank
	l = <-c // zString
	if l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad HIP Hit", l***REMOVED***, ""
	***REMOVED***
	rr.Hit = l.token // This can not contain spaces, see RFC 5205 Section 6.
	rr.HitLength = uint8(len(rr.Hit)) / 2

	<-c     // zBlank
	l = <-c // zString
	if l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad HIP PublicKey", l***REMOVED***, ""
	***REMOVED***
	rr.PublicKey = l.token // This cannot contain spaces
	rr.PublicKeyLength = uint16(base64.StdEncoding.DecodedLen(len(rr.PublicKey)))

	// RendezvousServers (if any)
	l = <-c
	var xs []string
	for l.value != zNewline && l.value != zEOF ***REMOVED***
		switch l.value ***REMOVED***
		case zString:
			if l.token == "@" ***REMOVED***
				xs = append(xs, o)
				l = <-c
				continue
			***REMOVED***
			_, ok := IsDomainName(l.token)
			if !ok || l.length == 0 || l.err ***REMOVED***
				return nil, &ParseError***REMOVED***f, "bad HIP RendezvousServers", l***REMOVED***, ""
			***REMOVED***
			if l.token[l.length-1] != '.' ***REMOVED***
				l.token = appendOrigin(l.token, o)
			***REMOVED***
			xs = append(xs, l.token)
		case zBlank:
			// Ok
		default:
			return nil, &ParseError***REMOVED***f, "bad HIP RendezvousServers", l***REMOVED***, ""
		***REMOVED***
		l = <-c
	***REMOVED***
	rr.RendezvousServers = xs
	return rr, nil, l.comment
***REMOVED***

func setCERT(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(CERT)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	if v, ok := StringToCertType[l.token]; ok ***REMOVED***
		rr.Type = v
	***REMOVED*** else if i, e := strconv.Atoi(l.token); e != nil ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CERT Type", l***REMOVED***, ""
	***REMOVED*** else ***REMOVED***
		rr.Type = uint16(i)
	***REMOVED***
	<-c     // zBlank
	l = <-c // zString
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CERT KeyTag", l***REMOVED***, ""
	***REMOVED***
	rr.KeyTag = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	if v, ok := StringToAlgorithm[l.token]; ok ***REMOVED***
		rr.Algorithm = v
	***REMOVED*** else if i, e := strconv.Atoi(l.token); e != nil ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CERT Algorithm", l***REMOVED***, ""
	***REMOVED*** else ***REMOVED***
		rr.Algorithm = uint8(i)
	***REMOVED***
	s, e1, c1 := endingToString(c, "bad CERT Certificate", f)
	if e1 != nil ***REMOVED***
		return nil, e1, c1
	***REMOVED***
	rr.Certificate = s
	return rr, nil, c1
***REMOVED***

func setOPENPGPKEY(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(OPENPGPKEY)
	rr.Hdr = h

	s, e, c1 := endingToString(c, "bad OPENPGPKEY PublicKey", f)
	if e != nil ***REMOVED***
		return nil, e, c1
	***REMOVED***
	rr.PublicKey = s
	return rr, nil, c1
***REMOVED***

func setSIG(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	r, e, s := setRRSIG(h, c, o, f)
	if r != nil ***REMOVED***
		return &SIG***REMOVED****r.(*RRSIG)***REMOVED***, e, s
	***REMOVED***
	return nil, e, s
***REMOVED***

func setRRSIG(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(RRSIG)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	if t, ok := StringToType[l.tokenUpper]; !ok ***REMOVED***
		if strings.HasPrefix(l.tokenUpper, "TYPE") ***REMOVED***
			t, ok = typeToInt(l.tokenUpper)
			if !ok ***REMOVED***
				return nil, &ParseError***REMOVED***f, "bad RRSIG Typecovered", l***REMOVED***, ""
			***REMOVED***
			rr.TypeCovered = t
		***REMOVED*** else ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad RRSIG Typecovered", l***REMOVED***, ""
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rr.TypeCovered = t
	***REMOVED***
	<-c // zBlank
	l = <-c
	i, err := strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RRSIG Algorithm", l***REMOVED***, ""
	***REMOVED***
	rr.Algorithm = uint8(i)
	<-c // zBlank
	l = <-c
	i, err = strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RRSIG Labels", l***REMOVED***, ""
	***REMOVED***
	rr.Labels = uint8(i)
	<-c // zBlank
	l = <-c
	i, err = strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RRSIG OrigTtl", l***REMOVED***, ""
	***REMOVED***
	rr.OrigTtl = uint32(i)
	<-c // zBlank
	l = <-c
	if i, err := StringToTime(l.token); err != nil ***REMOVED***
		// Try to see if all numeric and use it as epoch
		if i, err := strconv.ParseInt(l.token, 10, 64); err == nil ***REMOVED***
			// TODO(miek): error out on > MAX_UINT32, same below
			rr.Expiration = uint32(i)
		***REMOVED*** else ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad RRSIG Expiration", l***REMOVED***, ""
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rr.Expiration = i
	***REMOVED***
	<-c // zBlank
	l = <-c
	if i, err := StringToTime(l.token); err != nil ***REMOVED***
		if i, err := strconv.ParseInt(l.token, 10, 64); err == nil ***REMOVED***
			rr.Inception = uint32(i)
		***REMOVED*** else ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad RRSIG Inception", l***REMOVED***, ""
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rr.Inception = i
	***REMOVED***
	<-c // zBlank
	l = <-c
	i, err = strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RRSIG KeyTag", l***REMOVED***, ""
	***REMOVED***
	rr.KeyTag = uint16(i)
	<-c // zBlank
	l = <-c
	rr.SignerName = l.token
	if l.token == "@" ***REMOVED***
		rr.SignerName = o
	***REMOVED*** else ***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad RRSIG SignerName", l***REMOVED***, ""
		***REMOVED***
		if rr.SignerName[l.length-1] != '.' ***REMOVED***
			rr.SignerName = appendOrigin(rr.SignerName, o)
		***REMOVED***
	***REMOVED***
	s, e, c1 := endingToString(c, "bad RRSIG Signature", f)
	if e != nil ***REMOVED***
		return nil, e, c1
	***REMOVED***
	rr.Signature = s
	return rr, nil, c1
***REMOVED***

func setNSEC(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NSEC)
	rr.Hdr = h

	l := <-c
	rr.NextDomain = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.NextDomain = o
	***REMOVED*** else ***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad NSEC NextDomain", l***REMOVED***, ""
		***REMOVED***
		if rr.NextDomain[l.length-1] != '.' ***REMOVED***
			rr.NextDomain = appendOrigin(rr.NextDomain, o)
		***REMOVED***
	***REMOVED***

	rr.TypeBitMap = make([]uint16, 0)
	var (
		k  uint16
		ok bool
	)
	l = <-c
	for l.value != zNewline && l.value != zEOF ***REMOVED***
		switch l.value ***REMOVED***
		case zBlank:
			// Ok
		case zString:
			if k, ok = StringToType[l.tokenUpper]; !ok ***REMOVED***
				if k, ok = typeToInt(l.tokenUpper); !ok ***REMOVED***
					return nil, &ParseError***REMOVED***f, "bad NSEC TypeBitMap", l***REMOVED***, ""
				***REMOVED***
			***REMOVED***
			rr.TypeBitMap = append(rr.TypeBitMap, k)
		default:
			return nil, &ParseError***REMOVED***f, "bad NSEC TypeBitMap", l***REMOVED***, ""
		***REMOVED***
		l = <-c
	***REMOVED***
	return rr, nil, l.comment
***REMOVED***

func setNSEC3(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NSEC3)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3 Hash", l***REMOVED***, ""
	***REMOVED***
	rr.Hash = uint8(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3 Flags", l***REMOVED***, ""
	***REMOVED***
	rr.Flags = uint8(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3 Iterations", l***REMOVED***, ""
	***REMOVED***
	rr.Iterations = uint16(i)
	<-c
	l = <-c
	if len(l.token) == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3 Salt", l***REMOVED***, ""
	***REMOVED***
	rr.SaltLength = uint8(len(l.token)) / 2
	rr.Salt = l.token

	<-c
	l = <-c
	if len(l.token) == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3 NextDomain", l***REMOVED***, ""
	***REMOVED***
	rr.HashLength = 20 // Fix for NSEC3 (sha1 160 bits)
	rr.NextDomain = l.token

	rr.TypeBitMap = make([]uint16, 0)
	var (
		k  uint16
		ok bool
	)
	l = <-c
	for l.value != zNewline && l.value != zEOF ***REMOVED***
		switch l.value ***REMOVED***
		case zBlank:
			// Ok
		case zString:
			if k, ok = StringToType[l.tokenUpper]; !ok ***REMOVED***
				if k, ok = typeToInt(l.tokenUpper); !ok ***REMOVED***
					return nil, &ParseError***REMOVED***f, "bad NSEC3 TypeBitMap", l***REMOVED***, ""
				***REMOVED***
			***REMOVED***
			rr.TypeBitMap = append(rr.TypeBitMap, k)
		default:
			return nil, &ParseError***REMOVED***f, "bad NSEC3 TypeBitMap", l***REMOVED***, ""
		***REMOVED***
		l = <-c
	***REMOVED***
	return rr, nil, l.comment
***REMOVED***

func setNSEC3PARAM(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NSEC3PARAM)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3PARAM Hash", l***REMOVED***, ""
	***REMOVED***
	rr.Hash = uint8(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3PARAM Flags", l***REMOVED***, ""
	***REMOVED***
	rr.Flags = uint8(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NSEC3PARAM Iterations", l***REMOVED***, ""
	***REMOVED***
	rr.Iterations = uint16(i)
	<-c
	l = <-c
	rr.SaltLength = uint8(len(l.token))
	rr.Salt = l.token
	return rr, nil, ""
***REMOVED***

func setEUI48(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(EUI48)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.length != 17 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad EUI48 Address", l***REMOVED***, ""
	***REMOVED***
	addr := make([]byte, 12)
	dash := 0
	for i := 0; i < 10; i += 2 ***REMOVED***
		addr[i] = l.token[i+dash]
		addr[i+1] = l.token[i+1+dash]
		dash++
		if l.token[i+1+dash] != '-' ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad EUI48 Address", l***REMOVED***, ""
		***REMOVED***
	***REMOVED***
	addr[10] = l.token[15]
	addr[11] = l.token[16]

	i, e := strconv.ParseUint(string(addr), 16, 48)
	if e != nil ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad EUI48 Address", l***REMOVED***, ""
	***REMOVED***
	rr.Address = i
	return rr, nil, ""
***REMOVED***

func setEUI64(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(EUI64)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.length != 23 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad EUI64 Address", l***REMOVED***, ""
	***REMOVED***
	addr := make([]byte, 16)
	dash := 0
	for i := 0; i < 14; i += 2 ***REMOVED***
		addr[i] = l.token[i+dash]
		addr[i+1] = l.token[i+1+dash]
		dash++
		if l.token[i+1+dash] != '-' ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad EUI64 Address", l***REMOVED***, ""
		***REMOVED***
	***REMOVED***
	addr[14] = l.token[21]
	addr[15] = l.token[22]

	i, e := strconv.ParseUint(string(addr), 16, 64)
	if e != nil ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad EUI68 Address", l***REMOVED***, ""
	***REMOVED***
	rr.Address = uint64(i)
	return rr, nil, ""
***REMOVED***

func setWKS(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(WKS)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	rr.Address = net.ParseIP(l.token)
	if rr.Address == nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad WKS Address", l***REMOVED***, ""
	***REMOVED***

	<-c // zBlank
	l = <-c
	proto := "tcp"
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad WKS Protocol", l***REMOVED***, ""
	***REMOVED***
	rr.Protocol = uint8(i)
	switch rr.Protocol ***REMOVED***
	case 17:
		proto = "udp"
	case 6:
		proto = "tcp"
	default:
		return nil, &ParseError***REMOVED***f, "bad WKS Protocol", l***REMOVED***, ""
	***REMOVED***

	<-c
	l = <-c
	rr.BitMap = make([]uint16, 0)
	var (
		k   int
		err error
	)
	for l.value != zNewline && l.value != zEOF ***REMOVED***
		switch l.value ***REMOVED***
		case zBlank:
			// Ok
		case zString:
			if k, err = net.LookupPort(proto, l.token); err != nil ***REMOVED***
				i, e := strconv.Atoi(l.token) // If a number use that
				if e != nil ***REMOVED***
					return nil, &ParseError***REMOVED***f, "bad WKS BitMap", l***REMOVED***, ""
				***REMOVED***
				rr.BitMap = append(rr.BitMap, uint16(i))
			***REMOVED***
			rr.BitMap = append(rr.BitMap, uint16(k))
		default:
			return nil, &ParseError***REMOVED***f, "bad WKS BitMap", l***REMOVED***, ""
		***REMOVED***
		l = <-c
	***REMOVED***
	return rr, nil, l.comment
***REMOVED***

func setSSHFP(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(SSHFP)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad SSHFP Algorithm", l***REMOVED***, ""
	***REMOVED***
	rr.Algorithm = uint8(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad SSHFP Type", l***REMOVED***, ""
	***REMOVED***
	rr.Type = uint8(i)
	<-c // zBlank
	s, e1, c1 := endingToString(c, "bad SSHFP Fingerprint", f)
	if e1 != nil ***REMOVED***
		return nil, e1, c1
	***REMOVED***
	rr.FingerPrint = s
	return rr, nil, ""
***REMOVED***

func setDNSKEYs(h RR_Header, c chan lex, o, f, typ string) (RR, *ParseError, string) ***REMOVED***
	rr := new(DNSKEY)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad " + typ + " Flags", l***REMOVED***, ""
	***REMOVED***
	rr.Flags = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad " + typ + " Protocol", l***REMOVED***, ""
	***REMOVED***
	rr.Protocol = uint8(i)
	<-c     // zBlank
	l = <-c // zString
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad " + typ + " Algorithm", l***REMOVED***, ""
	***REMOVED***
	rr.Algorithm = uint8(i)
	s, e1, c1 := endingToString(c, "bad "+typ+" PublicKey", f)
	if e1 != nil ***REMOVED***
		return nil, e1, c1
	***REMOVED***
	rr.PublicKey = s
	return rr, nil, c1
***REMOVED***

func setKEY(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	r, e, s := setDNSKEYs(h, c, o, f, "KEY")
	if r != nil ***REMOVED***
		return &KEY***REMOVED****r.(*DNSKEY)***REMOVED***, e, s
	***REMOVED***
	return nil, e, s
***REMOVED***

func setDNSKEY(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	r, e, s := setDNSKEYs(h, c, o, f, "DNSKEY")
	return r, e, s
***REMOVED***

func setCDNSKEY(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	r, e, s := setDNSKEYs(h, c, o, f, "CDNSKEY")
	if r != nil ***REMOVED***
		return &CDNSKEY***REMOVED****r.(*DNSKEY)***REMOVED***, e, s
	***REMOVED***
	return nil, e, s
***REMOVED***

func setRKEY(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(RKEY)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RKEY Flags", l***REMOVED***, ""
	***REMOVED***
	rr.Flags = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RKEY Protocol", l***REMOVED***, ""
	***REMOVED***
	rr.Protocol = uint8(i)
	<-c     // zBlank
	l = <-c // zString
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RKEY Algorithm", l***REMOVED***, ""
	***REMOVED***
	rr.Algorithm = uint8(i)
	s, e1, c1 := endingToString(c, "bad RKEY PublicKey", f)
	if e1 != nil ***REMOVED***
		return nil, e1, c1
	***REMOVED***
	rr.PublicKey = s
	return rr, nil, c1
***REMOVED***

func setEID(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(EID)
	rr.Hdr = h
	s, e, c1 := endingToString(c, "bad EID Endpoint", f)
	if e != nil ***REMOVED***
		return nil, e, c1
	***REMOVED***
	rr.Endpoint = s
	return rr, nil, c1
***REMOVED***

func setNIMLOC(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NIMLOC)
	rr.Hdr = h
	s, e, c1 := endingToString(c, "bad NIMLOC Locator", f)
	if e != nil ***REMOVED***
		return nil, e, c1
	***REMOVED***
	rr.Locator = s
	return rr, nil, c1
***REMOVED***

func setGPOS(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(GPOS)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	_, e := strconv.ParseFloat(l.token, 64)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad GPOS Longitude", l***REMOVED***, ""
	***REMOVED***
	rr.Longitude = l.token
	<-c // zBlank
	l = <-c
	_, e = strconv.ParseFloat(l.token, 64)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad GPOS Latitude", l***REMOVED***, ""
	***REMOVED***
	rr.Latitude = l.token
	<-c // zBlank
	l = <-c
	_, e = strconv.ParseFloat(l.token, 64)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad GPOS Altitude", l***REMOVED***, ""
	***REMOVED***
	rr.Altitude = l.token
	return rr, nil, ""
***REMOVED***

func setDSs(h RR_Header, c chan lex, o, f, typ string) (RR, *ParseError, string) ***REMOVED***
	rr := new(DS)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad " + typ + " KeyTag", l***REMOVED***, ""
	***REMOVED***
	rr.KeyTag = uint16(i)
	<-c // zBlank
	l = <-c
	if i, e := strconv.Atoi(l.token); e != nil ***REMOVED***
		i, ok := StringToAlgorithm[l.tokenUpper]
		if !ok || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad " + typ + " Algorithm", l***REMOVED***, ""
		***REMOVED***
		rr.Algorithm = i
	***REMOVED*** else ***REMOVED***
		rr.Algorithm = uint8(i)
	***REMOVED***
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad " + typ + " DigestType", l***REMOVED***, ""
	***REMOVED***
	rr.DigestType = uint8(i)
	s, e1, c1 := endingToString(c, "bad "+typ+" Digest", f)
	if e1 != nil ***REMOVED***
		return nil, e1, c1
	***REMOVED***
	rr.Digest = s
	return rr, nil, c1
***REMOVED***

func setDS(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	r, e, s := setDSs(h, c, o, f, "DS")
	return r, e, s
***REMOVED***

func setDLV(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	r, e, s := setDSs(h, c, o, f, "DLV")
	if r != nil ***REMOVED***
		return &DLV***REMOVED****r.(*DS)***REMOVED***, e, s
	***REMOVED***
	return nil, e, s
***REMOVED***

func setCDS(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	r, e, s := setDSs(h, c, o, f, "CDS")
	if r != nil ***REMOVED***
		return &CDS***REMOVED****r.(*DS)***REMOVED***, e, s
	***REMOVED***
	return nil, e, s
***REMOVED***

func setTA(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(TA)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad TA KeyTag", l***REMOVED***, ""
	***REMOVED***
	rr.KeyTag = uint16(i)
	<-c // zBlank
	l = <-c
	if i, e := strconv.Atoi(l.token); e != nil ***REMOVED***
		i, ok := StringToAlgorithm[l.tokenUpper]
		if !ok || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad TA Algorithm", l***REMOVED***, ""
		***REMOVED***
		rr.Algorithm = i
	***REMOVED*** else ***REMOVED***
		rr.Algorithm = uint8(i)
	***REMOVED***
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad TA DigestType", l***REMOVED***, ""
	***REMOVED***
	rr.DigestType = uint8(i)
	s, e, c1 := endingToString(c, "bad TA Digest", f)
	if e != nil ***REMOVED***
		return nil, e.(*ParseError), c1
	***REMOVED***
	rr.Digest = s
	return rr, nil, c1
***REMOVED***

func setTLSA(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(TLSA)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad TLSA Usage", l***REMOVED***, ""
	***REMOVED***
	rr.Usage = uint8(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad TLSA Selector", l***REMOVED***, ""
	***REMOVED***
	rr.Selector = uint8(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad TLSA MatchingType", l***REMOVED***, ""
	***REMOVED***
	rr.MatchingType = uint8(i)
	// So this needs be e2 (i.e. different than e), because...??t
	s, e2, c1 := endingToString(c, "bad TLSA Certificate", f)
	if e2 != nil ***REMOVED***
		return nil, e2, c1
	***REMOVED***
	rr.Certificate = s
	return rr, nil, c1
***REMOVED***

func setRFC3597(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(RFC3597)
	rr.Hdr = h
	l := <-c
	if l.token != "\\#" ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RFC3597 Rdata", l***REMOVED***, ""
	***REMOVED***
	<-c // zBlank
	l = <-c
	rdlength, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RFC3597 Rdata ", l***REMOVED***, ""
	***REMOVED***

	s, e1, c1 := endingToString(c, "bad RFC3597 Rdata", f)
	if e1 != nil ***REMOVED***
		return nil, e1, c1
	***REMOVED***
	if rdlength*2 != len(s) ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad RFC3597 Rdata", l***REMOVED***, ""
	***REMOVED***
	rr.Rdata = s
	return rr, nil, c1
***REMOVED***

func setSPF(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(SPF)
	rr.Hdr = h

	s, e, c1 := endingToTxtSlice(c, "bad SPF Txt", f)
	if e != nil ***REMOVED***
		return nil, e, ""
	***REMOVED***
	rr.Txt = s
	return rr, nil, c1
***REMOVED***

func setTXT(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(TXT)
	rr.Hdr = h

	// no zBlank reading here, because all this rdata is TXT
	s, e, c1 := endingToTxtSlice(c, "bad TXT Txt", f)
	if e != nil ***REMOVED***
		return nil, e, ""
	***REMOVED***
	rr.Txt = s
	return rr, nil, c1
***REMOVED***

// identical to setTXT
func setNINFO(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NINFO)
	rr.Hdr = h

	s, e, c1 := endingToTxtSlice(c, "bad NINFO ZSData", f)
	if e != nil ***REMOVED***
		return nil, e, ""
	***REMOVED***
	rr.ZSData = s
	return rr, nil, c1
***REMOVED***

func setURI(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(URI)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED*** // Dynamic updates.
		return rr, nil, ""
	***REMOVED***

	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad URI Priority", l***REMOVED***, ""
	***REMOVED***
	rr.Priority = uint16(i)
	<-c // zBlank
	l = <-c
	i, e = strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad URI Weight", l***REMOVED***, ""
	***REMOVED***
	rr.Weight = uint16(i)

	<-c // zBlank
	s, err, c1 := endingToTxtSlice(c, "bad URI Target", f)
	if err != nil ***REMOVED***
		return nil, err, ""
	***REMOVED***
	if len(s) > 1 ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad URI Target", l***REMOVED***, ""
	***REMOVED***
	rr.Target = s[0]
	return rr, nil, c1
***REMOVED***

func setDHCID(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	// awesome record to parse!
	rr := new(DHCID)
	rr.Hdr = h

	s, e, c1 := endingToString(c, "bad DHCID Digest", f)
	if e != nil ***REMOVED***
		return nil, e, c1
	***REMOVED***
	rr.Digest = s
	return rr, nil, c1
***REMOVED***

func setNID(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(NID)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad NID Preference", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	u, err := stringToNodeID(l)
	if err != nil || l.err ***REMOVED***
		return nil, err, ""
	***REMOVED***
	rr.NodeID = u
	return rr, nil, ""
***REMOVED***

func setL32(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(L32)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad L32 Preference", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Locator32 = net.ParseIP(l.token)
	if rr.Locator32 == nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad L32 Locator", l***REMOVED***, ""
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setLP(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(LP)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LP Preference", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Fqdn = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Fqdn = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad LP Fqdn", l***REMOVED***, ""
	***REMOVED***
	if rr.Fqdn[l.length-1] != '.' ***REMOVED***
		rr.Fqdn = appendOrigin(rr.Fqdn, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setL64(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(L64)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad L64 Preference", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	u, err := stringToNodeID(l)
	if err != nil || l.err ***REMOVED***
		return nil, err, ""
	***REMOVED***
	rr.Locator64 = u
	return rr, nil, ""
***REMOVED***

func setUID(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(UID)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad UID Uid", l***REMOVED***, ""
	***REMOVED***
	rr.Uid = uint32(i)
	return rr, nil, ""
***REMOVED***

func setGID(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(GID)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad GID Gid", l***REMOVED***, ""
	***REMOVED***
	rr.Gid = uint32(i)
	return rr, nil, ""
***REMOVED***

func setUINFO(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(UINFO)
	rr.Hdr = h
	s, e, c1 := endingToTxtSlice(c, "bad UINFO Uinfo", f)
	if e != nil ***REMOVED***
		return nil, e, ""
	***REMOVED***
	rr.Uinfo = s[0] // silently discard anything above
	return rr, nil, c1
***REMOVED***

func setPX(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(PX)
	rr.Hdr = h

	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	i, e := strconv.Atoi(l.token)
	if e != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad PX Preference", l***REMOVED***, ""
	***REMOVED***
	rr.Preference = uint16(i)
	<-c     // zBlank
	l = <-c // zString
	rr.Map822 = l.token
	if l.length == 0 ***REMOVED***
		return rr, nil, ""
	***REMOVED***
	if l.token == "@" ***REMOVED***
		rr.Map822 = o
		return rr, nil, ""
	***REMOVED***
	_, ok := IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad PX Map822", l***REMOVED***, ""
	***REMOVED***
	if rr.Map822[l.length-1] != '.' ***REMOVED***
		rr.Map822 = appendOrigin(rr.Map822, o)
	***REMOVED***
	<-c     // zBlank
	l = <-c // zString
	rr.Mapx400 = l.token
	if l.token == "@" ***REMOVED***
		rr.Mapx400 = o
		return rr, nil, ""
	***REMOVED***
	_, ok = IsDomainName(l.token)
	if !ok || l.length == 0 || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad PX Mapx400", l***REMOVED***, ""
	***REMOVED***
	if rr.Mapx400[l.length-1] != '.' ***REMOVED***
		rr.Mapx400 = appendOrigin(rr.Mapx400, o)
	***REMOVED***
	return rr, nil, ""
***REMOVED***

func setIPSECKEY(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(IPSECKEY)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, err := strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad IPSECKEY Precedence", l***REMOVED***, ""
	***REMOVED***
	rr.Precedence = uint8(i)
	<-c // zBlank
	l = <-c
	i, err = strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad IPSECKEY GatewayType", l***REMOVED***, ""
	***REMOVED***
	rr.GatewayType = uint8(i)
	<-c // zBlank
	l = <-c
	i, err = strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad IPSECKEY Algorithm", l***REMOVED***, ""
	***REMOVED***
	rr.Algorithm = uint8(i)

	// Now according to GatewayType we can have different elements here
	<-c // zBlank
	l = <-c
	switch rr.GatewayType ***REMOVED***
	case 0:
		fallthrough
	case 3:
		rr.GatewayName = l.token
		if l.token == "@" ***REMOVED***
			rr.GatewayName = o
		***REMOVED***
		_, ok := IsDomainName(l.token)
		if !ok || l.length == 0 || l.err ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad IPSECKEY GatewayName", l***REMOVED***, ""
		***REMOVED***
		if rr.GatewayName[l.length-1] != '.' ***REMOVED***
			rr.GatewayName = appendOrigin(rr.GatewayName, o)
		***REMOVED***
	case 1:
		rr.GatewayA = net.ParseIP(l.token)
		if rr.GatewayA == nil ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad IPSECKEY GatewayA", l***REMOVED***, ""
		***REMOVED***
	case 2:
		rr.GatewayAAAA = net.ParseIP(l.token)
		if rr.GatewayAAAA == nil ***REMOVED***
			return nil, &ParseError***REMOVED***f, "bad IPSECKEY GatewayAAAA", l***REMOVED***, ""
		***REMOVED***
	default:
		return nil, &ParseError***REMOVED***f, "bad IPSECKEY GatewayType", l***REMOVED***, ""
	***REMOVED***

	s, e, c1 := endingToString(c, "bad IPSECKEY PublicKey", f)
	if e != nil ***REMOVED***
		return nil, e, c1
	***REMOVED***
	rr.PublicKey = s
	return rr, nil, c1
***REMOVED***

func setCAA(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
	rr := new(CAA)
	rr.Hdr = h
	l := <-c
	if l.length == 0 ***REMOVED***
		return rr, nil, l.comment
	***REMOVED***
	i, err := strconv.Atoi(l.token)
	if err != nil || l.err ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CAA Flag", l***REMOVED***, ""
	***REMOVED***
	rr.Flag = uint8(i)

	<-c     // zBlank
	l = <-c // zString
	if l.value != zString ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CAA Tag", l***REMOVED***, ""
	***REMOVED***
	rr.Tag = l.token

	<-c // zBlank
	s, e, c1 := endingToTxtSlice(c, "bad CAA Value", f)
	if e != nil ***REMOVED***
		return nil, e, ""
	***REMOVED***
	if len(s) > 1 ***REMOVED***
		return nil, &ParseError***REMOVED***f, "bad CAA Value", l***REMOVED***, ""
	***REMOVED***
	rr.Value = s[0]
	return rr, nil, c1
***REMOVED***

var typeToparserFunc = map[uint16]parserFunc***REMOVED***
	TypeAAAA:       parserFunc***REMOVED***setAAAA, false***REMOVED***,
	TypeAFSDB:      parserFunc***REMOVED***setAFSDB, false***REMOVED***,
	TypeA:          parserFunc***REMOVED***setA, false***REMOVED***,
	TypeCAA:        parserFunc***REMOVED***setCAA, true***REMOVED***,
	TypeCDS:        parserFunc***REMOVED***setCDS, true***REMOVED***,
	TypeCDNSKEY:    parserFunc***REMOVED***setCDNSKEY, true***REMOVED***,
	TypeCERT:       parserFunc***REMOVED***setCERT, true***REMOVED***,
	TypeCNAME:      parserFunc***REMOVED***setCNAME, false***REMOVED***,
	TypeDHCID:      parserFunc***REMOVED***setDHCID, true***REMOVED***,
	TypeDLV:        parserFunc***REMOVED***setDLV, true***REMOVED***,
	TypeDNAME:      parserFunc***REMOVED***setDNAME, false***REMOVED***,
	TypeKEY:        parserFunc***REMOVED***setKEY, true***REMOVED***,
	TypeDNSKEY:     parserFunc***REMOVED***setDNSKEY, true***REMOVED***,
	TypeDS:         parserFunc***REMOVED***setDS, true***REMOVED***,
	TypeEID:        parserFunc***REMOVED***setEID, true***REMOVED***,
	TypeEUI48:      parserFunc***REMOVED***setEUI48, false***REMOVED***,
	TypeEUI64:      parserFunc***REMOVED***setEUI64, false***REMOVED***,
	TypeGID:        parserFunc***REMOVED***setGID, false***REMOVED***,
	TypeGPOS:       parserFunc***REMOVED***setGPOS, false***REMOVED***,
	TypeHINFO:      parserFunc***REMOVED***setHINFO, true***REMOVED***,
	TypeHIP:        parserFunc***REMOVED***setHIP, true***REMOVED***,
	TypeIPSECKEY:   parserFunc***REMOVED***setIPSECKEY, true***REMOVED***,
	TypeKX:         parserFunc***REMOVED***setKX, false***REMOVED***,
	TypeL32:        parserFunc***REMOVED***setL32, false***REMOVED***,
	TypeL64:        parserFunc***REMOVED***setL64, false***REMOVED***,
	TypeLOC:        parserFunc***REMOVED***setLOC, true***REMOVED***,
	TypeLP:         parserFunc***REMOVED***setLP, false***REMOVED***,
	TypeMB:         parserFunc***REMOVED***setMB, false***REMOVED***,
	TypeMD:         parserFunc***REMOVED***setMD, false***REMOVED***,
	TypeMF:         parserFunc***REMOVED***setMF, false***REMOVED***,
	TypeMG:         parserFunc***REMOVED***setMG, false***REMOVED***,
	TypeMINFO:      parserFunc***REMOVED***setMINFO, false***REMOVED***,
	TypeMR:         parserFunc***REMOVED***setMR, false***REMOVED***,
	TypeMX:         parserFunc***REMOVED***setMX, false***REMOVED***,
	TypeNAPTR:      parserFunc***REMOVED***setNAPTR, false***REMOVED***,
	TypeNID:        parserFunc***REMOVED***setNID, false***REMOVED***,
	TypeNIMLOC:     parserFunc***REMOVED***setNIMLOC, true***REMOVED***,
	TypeNINFO:      parserFunc***REMOVED***setNINFO, true***REMOVED***,
	TypeNSAPPTR:    parserFunc***REMOVED***setNSAPPTR, false***REMOVED***,
	TypeNSEC3PARAM: parserFunc***REMOVED***setNSEC3PARAM, false***REMOVED***,
	TypeNSEC3:      parserFunc***REMOVED***setNSEC3, true***REMOVED***,
	TypeNSEC:       parserFunc***REMOVED***setNSEC, true***REMOVED***,
	TypeNS:         parserFunc***REMOVED***setNS, false***REMOVED***,
	TypeOPENPGPKEY: parserFunc***REMOVED***setOPENPGPKEY, true***REMOVED***,
	TypePTR:        parserFunc***REMOVED***setPTR, false***REMOVED***,
	TypePX:         parserFunc***REMOVED***setPX, false***REMOVED***,
	TypeSIG:        parserFunc***REMOVED***setSIG, true***REMOVED***,
	TypeRKEY:       parserFunc***REMOVED***setRKEY, true***REMOVED***,
	TypeRP:         parserFunc***REMOVED***setRP, false***REMOVED***,
	TypeRRSIG:      parserFunc***REMOVED***setRRSIG, true***REMOVED***,
	TypeRT:         parserFunc***REMOVED***setRT, false***REMOVED***,
	TypeSOA:        parserFunc***REMOVED***setSOA, false***REMOVED***,
	TypeSPF:        parserFunc***REMOVED***setSPF, true***REMOVED***,
	TypeSRV:        parserFunc***REMOVED***setSRV, false***REMOVED***,
	TypeSSHFP:      parserFunc***REMOVED***setSSHFP, true***REMOVED***,
	TypeTALINK:     parserFunc***REMOVED***setTALINK, false***REMOVED***,
	TypeTA:         parserFunc***REMOVED***setTA, true***REMOVED***,
	TypeTLSA:       parserFunc***REMOVED***setTLSA, true***REMOVED***,
	TypeTXT:        parserFunc***REMOVED***setTXT, true***REMOVED***,
	TypeUID:        parserFunc***REMOVED***setUID, false***REMOVED***,
	TypeUINFO:      parserFunc***REMOVED***setUINFO, true***REMOVED***,
	TypeURI:        parserFunc***REMOVED***setURI, true***REMOVED***,
	TypeWKS:        parserFunc***REMOVED***setWKS, true***REMOVED***,
	TypeX25:        parserFunc***REMOVED***setX25, false***REMOVED***,
***REMOVED***
