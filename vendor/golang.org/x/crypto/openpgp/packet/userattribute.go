// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
)

const UserAttrImageSubpacket = 1

// UserAttribute is capable of storing other types of data about a user
// beyond name, email and a text comment. In practice, user attributes are typically used
// to store a signed thumbnail photo JPEG image of the user.
// See RFC 4880, section 5.12.
type UserAttribute struct ***REMOVED***
	Contents []*OpaqueSubpacket
***REMOVED***

// NewUserAttributePhoto creates a user attribute packet
// containing the given images.
func NewUserAttributePhoto(photos ...image.Image) (uat *UserAttribute, err error) ***REMOVED***
	uat = new(UserAttribute)
	for _, photo := range photos ***REMOVED***
		var buf bytes.Buffer
		// RFC 4880, Section 5.12.1.
		data := []byte***REMOVED***
			0x10, 0x00, // Little-endian image header length (16 bytes)
			0x01,       // Image header version 1
			0x01,       // JPEG
			0, 0, 0, 0, // 12 reserved octets, must be all zero.
			0, 0, 0, 0,
			0, 0, 0, 0***REMOVED***
		if _, err = buf.Write(data); err != nil ***REMOVED***
			return
		***REMOVED***
		if err = jpeg.Encode(&buf, photo, nil); err != nil ***REMOVED***
			return
		***REMOVED***
		uat.Contents = append(uat.Contents, &OpaqueSubpacket***REMOVED***
			SubType:  UserAttrImageSubpacket,
			Contents: buf.Bytes()***REMOVED***)
	***REMOVED***
	return
***REMOVED***

// NewUserAttribute creates a new user attribute packet containing the given subpackets.
func NewUserAttribute(contents ...*OpaqueSubpacket) *UserAttribute ***REMOVED***
	return &UserAttribute***REMOVED***Contents: contents***REMOVED***
***REMOVED***

func (uat *UserAttribute) parse(r io.Reader) (err error) ***REMOVED***
	// RFC 4880, section 5.13
	b, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	uat.Contents, err = OpaqueSubpackets(b)
	return
***REMOVED***

// Serialize marshals the user attribute to w in the form of an OpenPGP packet, including
// header.
func (uat *UserAttribute) Serialize(w io.Writer) (err error) ***REMOVED***
	var buf bytes.Buffer
	for _, sp := range uat.Contents ***REMOVED***
		sp.Serialize(&buf)
	***REMOVED***
	if err = serializeHeader(w, packetTypeUserAttribute, buf.Len()); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = w.Write(buf.Bytes())
	return
***REMOVED***

// ImageData returns zero or more byte slices, each containing
// JPEG File Interchange Format (JFIF), for each photo in the
// the user attribute packet.
func (uat *UserAttribute) ImageData() (imageData [][]byte) ***REMOVED***
	for _, sp := range uat.Contents ***REMOVED***
		if sp.SubType == UserAttrImageSubpacket && len(sp.Contents) > 16 ***REMOVED***
			imageData = append(imageData, sp.Contents[16:])
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
