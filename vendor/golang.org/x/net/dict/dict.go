// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dict implements the Dictionary Server Protocol
// as defined in RFC 2229.
package dict // import "golang.org/x/net/dict"

import (
	"net/textproto"
	"strconv"
	"strings"
)

// A Client represents a client connection to a dictionary server.
type Client struct ***REMOVED***
	text *textproto.Conn
***REMOVED***

// Dial returns a new client connected to a dictionary server at
// addr on the given network.
func Dial(network, addr string) (*Client, error) ***REMOVED***
	text, err := textproto.Dial(network, addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	_, _, err = text.ReadCodeLine(220)
	if err != nil ***REMOVED***
		text.Close()
		return nil, err
	***REMOVED***
	return &Client***REMOVED***text: text***REMOVED***, nil
***REMOVED***

// Close closes the connection to the dictionary server.
func (c *Client) Close() error ***REMOVED***
	return c.text.Close()
***REMOVED***

// A Dict represents a dictionary available on the server.
type Dict struct ***REMOVED***
	Name string // short name of dictionary
	Desc string // long description
***REMOVED***

// Dicts returns a list of the dictionaries available on the server.
func (c *Client) Dicts() ([]Dict, error) ***REMOVED***
	id, err := c.text.Cmd("SHOW DB")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.text.StartResponse(id)
	defer c.text.EndResponse(id)

	_, _, err = c.text.ReadCodeLine(110)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	lines, err := c.text.ReadDotLines()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	_, _, err = c.text.ReadCodeLine(250)

	dicts := make([]Dict, len(lines))
	for i := range dicts ***REMOVED***
		d := &dicts[i]
		a, _ := fields(lines[i])
		if len(a) < 2 ***REMOVED***
			return nil, textproto.ProtocolError("invalid dictionary: " + lines[i])
		***REMOVED***
		d.Name = a[0]
		d.Desc = a[1]
	***REMOVED***
	return dicts, err
***REMOVED***

// A Defn represents a definition.
type Defn struct ***REMOVED***
	Dict Dict   // Dict where definition was found
	Word string // Word being defined
	Text []byte // Definition text, typically multiple lines
***REMOVED***

// Define requests the definition of the given word.
// The argument dict names the dictionary to use,
// the Name field of a Dict returned by Dicts.
//
// The special dictionary name "*" means to look in all the
// server's dictionaries.
// The special dictionary name "!" means to look in all the
// server's dictionaries in turn, stopping after finding the word
// in one of them.
func (c *Client) Define(dict, word string) ([]*Defn, error) ***REMOVED***
	id, err := c.text.Cmd("DEFINE %s %q", dict, word)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.text.StartResponse(id)
	defer c.text.EndResponse(id)

	_, line, err := c.text.ReadCodeLine(150)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	a, _ := fields(line)
	if len(a) < 1 ***REMOVED***
		return nil, textproto.ProtocolError("malformed response: " + line)
	***REMOVED***
	n, err := strconv.Atoi(a[0])
	if err != nil ***REMOVED***
		return nil, textproto.ProtocolError("invalid definition count: " + a[0])
	***REMOVED***
	def := make([]*Defn, n)
	for i := 0; i < n; i++ ***REMOVED***
		_, line, err = c.text.ReadCodeLine(151)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		a, _ := fields(line)
		if len(a) < 3 ***REMOVED***
			// skip it, to keep protocol in sync
			i--
			n--
			def = def[0:n]
			continue
		***REMOVED***
		d := &Defn***REMOVED***Word: a[0], Dict: Dict***REMOVED***a[1], a[2]***REMOVED******REMOVED***
		d.Text, err = c.text.ReadDotBytes()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		def[i] = d
	***REMOVED***
	_, _, err = c.text.ReadCodeLine(250)
	return def, err
***REMOVED***

// Fields returns the fields in s.
// Fields are space separated unquoted words
// or quoted with single or double quote.
func fields(s string) ([]string, error) ***REMOVED***
	var v []string
	i := 0
	for ***REMOVED***
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') ***REMOVED***
			i++
		***REMOVED***
		if i >= len(s) ***REMOVED***
			break
		***REMOVED***
		if s[i] == '"' || s[i] == '\'' ***REMOVED***
			q := s[i]
			// quoted string
			var j int
			for j = i + 1; ; j++ ***REMOVED***
				if j >= len(s) ***REMOVED***
					return nil, textproto.ProtocolError("malformed quoted string")
				***REMOVED***
				if s[j] == '\\' ***REMOVED***
					j++
					continue
				***REMOVED***
				if s[j] == q ***REMOVED***
					j++
					break
				***REMOVED***
			***REMOVED***
			v = append(v, unquote(s[i+1:j-1]))
			i = j
		***REMOVED*** else ***REMOVED***
			// atom
			var j int
			for j = i; j < len(s); j++ ***REMOVED***
				if s[j] == ' ' || s[j] == '\t' || s[j] == '\\' || s[j] == '"' || s[j] == '\'' ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			v = append(v, s[i:j])
			i = j
		***REMOVED***
		if i < len(s) ***REMOVED***
			c := s[i]
			if c != ' ' && c != '\t' ***REMOVED***
				return nil, textproto.ProtocolError("quotes not on word boundaries")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return v, nil
***REMOVED***

func unquote(s string) string ***REMOVED***
	if strings.Index(s, "\\") < 0 ***REMOVED***
		return s
	***REMOVED***
	b := []byte(s)
	w := 0
	for r := 0; r < len(b); r++ ***REMOVED***
		c := b[r]
		if c == '\\' ***REMOVED***
			r++
			c = b[r]
		***REMOVED***
		b[w] = c
		w++
	***REMOVED***
	return string(b[0:w])
***REMOVED***
