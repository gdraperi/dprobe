package nl

import "encoding/binary"

const (
	MPLS_LS_LABEL_SHIFT = 12
	MPLS_LS_S_SHIFT     = 8
)

func EncodeMPLSStack(labels ...int) []byte ***REMOVED***
	b := make([]byte, 4*len(labels))
	for idx, label := range labels ***REMOVED***
		l := label << MPLS_LS_LABEL_SHIFT
		if idx == len(labels)-1 ***REMOVED***
			l |= 1 << MPLS_LS_S_SHIFT
		***REMOVED***
		binary.BigEndian.PutUint32(b[idx*4:], uint32(l))
	***REMOVED***
	return b
***REMOVED***

func DecodeMPLSStack(buf []byte) []int ***REMOVED***
	if len(buf)%4 != 0 ***REMOVED***
		return nil
	***REMOVED***
	stack := make([]int, 0, len(buf)/4)
	for len(buf) > 0 ***REMOVED***
		l := binary.BigEndian.Uint32(buf[:4])
		buf = buf[4:]
		stack = append(stack, int(l)>>MPLS_LS_LABEL_SHIFT)
		if (l>>MPLS_LS_S_SHIFT)&1 > 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return stack
***REMOVED***
