package dns

// Implement a simple scanner, return a byte stream from an io reader.

import (
	"bufio"
	"io"
	"text/scanner"
)

type scan struct ***REMOVED***
	src      *bufio.Reader
	position scanner.Position
	eof      bool // Have we just seen a eof
***REMOVED***

func scanInit(r io.Reader) *scan ***REMOVED***
	s := new(scan)
	s.src = bufio.NewReader(r)
	s.position.Line = 1
	return s
***REMOVED***

// tokenText returns the next byte from the input
func (s *scan) tokenText() (byte, error) ***REMOVED***
	c, err := s.src.ReadByte()
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***
	// delay the newline handling until the next token is delivered,
	// fixes off-by-one errors when reporting a parse error.
	if s.eof == true ***REMOVED***
		s.position.Line++
		s.position.Column = 0
		s.eof = false
	***REMOVED***
	if c == '\n' ***REMOVED***
		s.eof = true
		return c, nil
	***REMOVED***
	s.position.Column++
	return c, nil
***REMOVED***
