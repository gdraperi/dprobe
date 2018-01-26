/*
HTTP Content-Type Autonegotiation.

The functions in this package implement the behaviour specified in
http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

Copyright (c) 2011, Open Knowledge Foundation Ltd.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

    Redistributions of source code must retain the above copyright
    notice, this list of conditions and the following disclaimer.

    Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in
    the documentation and/or other materials provided with the
    distribution.

    Neither the name of the Open Knowledge Foundation Ltd. nor the
    names of its contributors may be used to endorse or promote
    products derived from this software without specific prior written
    permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


*/
package goautoneg

import (
	"sort"
	"strconv"
	"strings"
)

// Structure to represent a clause in an HTTP Accept Header
type Accept struct ***REMOVED***
	Type, SubType string
	Q             float64
	Params        map[string]string
***REMOVED***

// For internal use, so that we can use the sort interface
type accept_slice []Accept

func (accept accept_slice) Len() int ***REMOVED***
	slice := []Accept(accept)
	return len(slice)
***REMOVED***

func (accept accept_slice) Less(i, j int) bool ***REMOVED***
	slice := []Accept(accept)
	ai, aj := slice[i], slice[j]
	if ai.Q > aj.Q ***REMOVED***
		return true
	***REMOVED***
	if ai.Type != "*" && aj.Type == "*" ***REMOVED***
		return true
	***REMOVED***
	if ai.SubType != "*" && aj.SubType == "*" ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (accept accept_slice) Swap(i, j int) ***REMOVED***
	slice := []Accept(accept)
	slice[i], slice[j] = slice[j], slice[i]
***REMOVED***

// Parse an Accept Header string returning a sorted list
// of clauses
func ParseAccept(header string) (accept []Accept) ***REMOVED***
	parts := strings.Split(header, ",")
	accept = make([]Accept, 0, len(parts))
	for _, part := range parts ***REMOVED***
		part := strings.Trim(part, " ")

		a := Accept***REMOVED******REMOVED***
		a.Params = make(map[string]string)
		a.Q = 1.0

		mrp := strings.Split(part, ";")

		media_range := mrp[0]
		sp := strings.Split(media_range, "/")
		a.Type = strings.Trim(sp[0], " ")

		switch ***REMOVED***
		case len(sp) == 1 && a.Type == "*":
			a.SubType = "*"
		case len(sp) == 2:
			a.SubType = strings.Trim(sp[1], " ")
		default:
			continue
		***REMOVED***

		if len(mrp) == 1 ***REMOVED***
			accept = append(accept, a)
			continue
		***REMOVED***

		for _, param := range mrp[1:] ***REMOVED***
			sp := strings.SplitN(param, "=", 2)
			if len(sp) != 2 ***REMOVED***
				continue
			***REMOVED***
			token := strings.Trim(sp[0], " ")
			if token == "q" ***REMOVED***
				a.Q, _ = strconv.ParseFloat(sp[1], 32)
			***REMOVED*** else ***REMOVED***
				a.Params[token] = strings.Trim(sp[1], " ")
			***REMOVED***
		***REMOVED***

		accept = append(accept, a)
	***REMOVED***

	slice := accept_slice(accept)
	sort.Sort(slice)

	return
***REMOVED***

// Negotiate the most appropriate content_type given the accept header
// and a list of alternatives.
func Negotiate(header string, alternatives []string) (content_type string) ***REMOVED***
	asp := make([][]string, 0, len(alternatives))
	for _, ctype := range alternatives ***REMOVED***
		asp = append(asp, strings.SplitN(ctype, "/", 2))
	***REMOVED***
	for _, clause := range ParseAccept(header) ***REMOVED***
		for i, ctsp := range asp ***REMOVED***
			if clause.Type == ctsp[0] && clause.SubType == ctsp[1] ***REMOVED***
				content_type = alternatives[i]
				return
			***REMOVED***
			if clause.Type == ctsp[0] && clause.SubType == "*" ***REMOVED***
				content_type = alternatives[i]
				return
			***REMOVED***
			if clause.Type == "*" && clause.SubType == "*" ***REMOVED***
				content_type = alternatives[i]
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
