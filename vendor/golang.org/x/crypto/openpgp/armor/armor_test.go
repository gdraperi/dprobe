// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package armor

import (
	"bytes"
	"hash/adler32"
	"io/ioutil"
	"testing"
)

func TestDecodeEncode(t *testing.T) ***REMOVED***
	buf := bytes.NewBuffer([]byte(armorExample1))
	result, err := Decode(buf)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	expectedType := "PGP SIGNATURE"
	if result.Type != expectedType ***REMOVED***
		t.Errorf("result.Type: got:%s want:%s", result.Type, expectedType)
	***REMOVED***
	if len(result.Header) != 1 ***REMOVED***
		t.Errorf("len(result.Header): got:%d want:1", len(result.Header))
	***REMOVED***
	v, ok := result.Header["Version"]
	if !ok || v != "GnuPG v1.4.10 (GNU/Linux)" ***REMOVED***
		t.Errorf("result.Header: got:%#v", result.Header)
	***REMOVED***

	contents, err := ioutil.ReadAll(result.Body)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	if adler32.Checksum(contents) != 0x27b144be ***REMOVED***
		t.Errorf("contents: got: %x", contents)
	***REMOVED***

	buf = bytes.NewBuffer(nil)
	w, err := Encode(buf, result.Type, result.Header)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	_, err = w.Write(contents)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	w.Close()

	if !bytes.Equal(buf.Bytes(), []byte(armorExample1)) ***REMOVED***
		t.Errorf("got: %s\nwant: %s", string(buf.Bytes()), armorExample1)
	***REMOVED***
***REMOVED***

func TestLongHeader(t *testing.T) ***REMOVED***
	buf := bytes.NewBuffer([]byte(armorLongLine))
	result, err := Decode(buf)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	value, ok := result.Header["Version"]
	if !ok ***REMOVED***
		t.Errorf("missing Version header")
	***REMOVED***
	if value != longValueExpected ***REMOVED***
		t.Errorf("got: %s want: %s", value, longValueExpected)
	***REMOVED***
***REMOVED***

const armorExample1 = `-----BEGIN PGP SIGNATURE-----
Version: GnuPG v1.4.10 (GNU/Linux)

iJwEAAECAAYFAk1Fv/0ACgkQo01+GMIMMbsYTwQAiAw+QAaNfY6WBdplZ/uMAccm
4g+81QPmTSGHnetSb6WBiY13kVzK4HQiZH8JSkmmroMLuGeJwsRTEL4wbjRyUKEt
p1xwUZDECs234F1xiG5enc5SGlRtP7foLBz9lOsjx+LEcA4sTl5/2eZR9zyFZqWW
TxRjs+fJCIFuo71xb1g=
=/teI
-----END PGP SIGNATURE-----`

const armorLongLine = `-----BEGIN PGP SIGNATURE-----
Version: 0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz

iQEcBAABAgAGBQJMtFESAAoJEKsQXJGvOPsVj40H/1WW6jaMXv4BW+1ueDSMDwM8
kx1fLOXbVM5/Kn5LStZNt1jWWnpxdz7eq3uiqeCQjmqUoRde3YbB2EMnnwRbAhpp
cacnAvy9ZQ78OTxUdNW1mhX5bS6q1MTEJnl+DcyigD70HG/yNNQD7sOPMdYQw0TA
byQBwmLwmTsuZsrYqB68QyLHI+DUugn+kX6Hd2WDB62DKa2suoIUIHQQCd/ofwB3
WfCYInXQKKOSxu2YOg2Eb4kLNhSMc1i9uKUWAH+sdgJh7NBgdoE4MaNtBFkHXRvv
okWuf3+xA9ksp1npSY/mDvgHijmjvtpRDe6iUeqfCn8N9u9CBg8geANgaG8+QA4=
=wfQG
-----END PGP SIGNATURE-----`

const longValueExpected = "0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz"
