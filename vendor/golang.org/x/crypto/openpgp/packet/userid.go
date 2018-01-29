// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"io"
	"io/ioutil"
	"strings"
)

// UserId contains text that is intended to represent the name and email
// address of the key holder. See RFC 4880, section 5.11. By convention, this
// takes the form "Full Name (Comment) <email@example.com>"
type UserId struct ***REMOVED***
	Id string // By convention, this takes the form "Full Name (Comment) <email@example.com>" which is split out in the fields below.

	Name, Comment, Email string
***REMOVED***

func hasInvalidCharacters(s string) bool ***REMOVED***
	for _, c := range s ***REMOVED***
		switch c ***REMOVED***
		case '(', ')', '<', '>', 0:
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// NewUserId returns a UserId or nil if any of the arguments contain invalid
// characters. The invalid characters are '\x00', '(', ')', '<' and '>'
func NewUserId(name, comment, email string) *UserId ***REMOVED***
	// RFC 4880 doesn't deal with the structure of userid strings; the
	// name, comment and email form is just a convention. However, there's
	// no convention about escaping the metacharacters and GPG just refuses
	// to create user ids where, say, the name contains a '('. We mirror
	// this behaviour.

	if hasInvalidCharacters(name) || hasInvalidCharacters(comment) || hasInvalidCharacters(email) ***REMOVED***
		return nil
	***REMOVED***

	uid := new(UserId)
	uid.Name, uid.Comment, uid.Email = name, comment, email
	uid.Id = name
	if len(comment) > 0 ***REMOVED***
		if len(uid.Id) > 0 ***REMOVED***
			uid.Id += " "
		***REMOVED***
		uid.Id += "("
		uid.Id += comment
		uid.Id += ")"
	***REMOVED***
	if len(email) > 0 ***REMOVED***
		if len(uid.Id) > 0 ***REMOVED***
			uid.Id += " "
		***REMOVED***
		uid.Id += "<"
		uid.Id += email
		uid.Id += ">"
	***REMOVED***
	return uid
***REMOVED***

func (uid *UserId) parse(r io.Reader) (err error) ***REMOVED***
	// RFC 4880, section 5.11
	b, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	uid.Id = string(b)
	uid.Name, uid.Comment, uid.Email = parseUserId(uid.Id)
	return
***REMOVED***

// Serialize marshals uid to w in the form of an OpenPGP packet, including
// header.
func (uid *UserId) Serialize(w io.Writer) error ***REMOVED***
	err := serializeHeader(w, packetTypeUserId, len(uid.Id))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = w.Write([]byte(uid.Id))
	return err
***REMOVED***

// parseUserId extracts the name, comment and email from a user id string that
// is formatted as "Full Name (Comment) <email@example.com>".
func parseUserId(id string) (name, comment, email string) ***REMOVED***
	var n, c, e struct ***REMOVED***
		start, end int
	***REMOVED***
	var state int

	for offset, rune := range id ***REMOVED***
		switch state ***REMOVED***
		case 0:
			// Entering name
			n.start = offset
			state = 1
			fallthrough
		case 1:
			// In name
			if rune == '(' ***REMOVED***
				state = 2
				n.end = offset
			***REMOVED*** else if rune == '<' ***REMOVED***
				state = 5
				n.end = offset
			***REMOVED***
		case 2:
			// Entering comment
			c.start = offset
			state = 3
			fallthrough
		case 3:
			// In comment
			if rune == ')' ***REMOVED***
				state = 4
				c.end = offset
			***REMOVED***
		case 4:
			// Between comment and email
			if rune == '<' ***REMOVED***
				state = 5
			***REMOVED***
		case 5:
			// Entering email
			e.start = offset
			state = 6
			fallthrough
		case 6:
			// In email
			if rune == '>' ***REMOVED***
				state = 7
				e.end = offset
			***REMOVED***
		default:
			// After email
		***REMOVED***
	***REMOVED***
	switch state ***REMOVED***
	case 1:
		// ended in the name
		n.end = len(id)
	case 3:
		// ended in comment
		c.end = len(id)
	case 6:
		// ended in email
		e.end = len(id)
	***REMOVED***

	name = strings.TrimSpace(id[n.start:n.end])
	comment = strings.TrimSpace(id[c.start:c.end])
	email = strings.TrimSpace(id[e.start:e.end])
	return
***REMOVED***
