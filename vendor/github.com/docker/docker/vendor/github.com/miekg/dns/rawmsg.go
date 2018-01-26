package dns

// These raw* functions do not use reflection, they directly set the values
// in the buffer. There are faster than their reflection counterparts.

// RawSetId sets the message id in buf.
func rawSetId(msg []byte, i uint16) bool ***REMOVED***
	if len(msg) < 2 ***REMOVED***
		return false
	***REMOVED***
	msg[0], msg[1] = packUint16(i)
	return true
***REMOVED***

// rawSetQuestionLen sets the length of the question section.
func rawSetQuestionLen(msg []byte, i uint16) bool ***REMOVED***
	if len(msg) < 6 ***REMOVED***
		return false
	***REMOVED***
	msg[4], msg[5] = packUint16(i)
	return true
***REMOVED***

// rawSetAnswerLen sets the lenght of the answer section.
func rawSetAnswerLen(msg []byte, i uint16) bool ***REMOVED***
	if len(msg) < 8 ***REMOVED***
		return false
	***REMOVED***
	msg[6], msg[7] = packUint16(i)
	return true
***REMOVED***

// rawSetsNsLen sets the lenght of the authority section.
func rawSetNsLen(msg []byte, i uint16) bool ***REMOVED***
	if len(msg) < 10 ***REMOVED***
		return false
	***REMOVED***
	msg[8], msg[9] = packUint16(i)
	return true
***REMOVED***

// rawSetExtraLen sets the lenght of the additional section.
func rawSetExtraLen(msg []byte, i uint16) bool ***REMOVED***
	if len(msg) < 12 ***REMOVED***
		return false
	***REMOVED***
	msg[10], msg[11] = packUint16(i)
	return true
***REMOVED***

// rawSetRdlength sets the rdlength in the header of
// the RR. The offset 'off' must be positioned at the
// start of the header of the RR, 'end' must be the
// end of the RR.
func rawSetRdlength(msg []byte, off, end int) bool ***REMOVED***
	l := len(msg)
Loop:
	for ***REMOVED***
		if off+1 > l ***REMOVED***
			return false
		***REMOVED***
		c := int(msg[off])
		off++
		switch c & 0xC0 ***REMOVED***
		case 0x00:
			if c == 0x00 ***REMOVED***
				// End of the domainname
				break Loop
			***REMOVED***
			if off+c > l ***REMOVED***
				return false
			***REMOVED***
			off += c

		case 0xC0:
			// pointer, next byte included, ends domainname
			off++
			break Loop
		***REMOVED***
	***REMOVED***
	// The domainname has been seen, we at the start of the fixed part in the header.
	// Type is 2 bytes, class is 2 bytes, ttl 4 and then 2 bytes for the length.
	off += 2 + 2 + 4
	if off+2 > l ***REMOVED***
		return false
	***REMOVED***
	//off+1 is the end of the header, 'end' is the end of the rr
	//so 'end' - 'off+2' is the length of the rdata
	rdatalen := end - (off + 2)
	if rdatalen > 0xFFFF ***REMOVED***
		return false
	***REMOVED***
	msg[off], msg[off+1] = packUint16(uint16(rdatalen))
	return true
***REMOVED***
