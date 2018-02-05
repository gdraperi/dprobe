// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

// This file contains test data.

import (
	"io"
	"strings"
)

type testLoader struct ***REMOVED***
***REMOVED***

func (t testLoader) Len() int ***REMOVED***
	return len(testFiles)
***REMOVED***

func (t testLoader) Path(i int) string ***REMOVED***
	return testPaths[i]
***REMOVED***

func (t testLoader) Reader(i int) (io.ReadCloser, error) ***REMOVED***
	return &reader***REMOVED****strings.NewReader(testFiles[i])***REMOVED***, nil
***REMOVED***

// reader adds a dummy Close method to strings.Reader so that it
// satisfies the io.ReadCloser interface.
type reader struct ***REMOVED***
	strings.Reader
***REMOVED***

func (r reader) Close() error ***REMOVED***
	return nil
***REMOVED***

var (
	testFiles = []string***REMOVED***de_xml, gsw_xml, root_xml***REMOVED***
	testPaths = []string***REMOVED***
		"common/main/de.xml",
		"common/main/gsw.xml",
		"common/main/root.xml",
	***REMOVED***
)

var root_xml = `<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE ldml SYSTEM "../../common/dtd/ldml.dtd">
<ldml>
	<identity>
		<language type="root"/>
		<generation date="now"/>
	</identity>
	<characters>
		<exemplarCharacters>[]</exemplarCharacters>
		<exemplarCharacters type="auxiliary">[]</exemplarCharacters>
		<exemplarCharacters type="punctuation">[\- ‐ – — … ' ‘ ‚ &quot; “ „ \&amp; #]</exemplarCharacters>
		<ellipsis type="final">***REMOVED***0***REMOVED***…</ellipsis>
		<ellipsis type="initial">…***REMOVED***0***REMOVED***</ellipsis>
		<moreInformation>?</moreInformation>
	</characters>
	<dates>
		<calendars>
			<default choice="gregorian"/>
			<calendar type="buddhist">
				<months>
					<alias source="locale" path="../../calendar[@type='gregorian']/months"/>
				</months>
			</calendar>
			<calendar type="chinese">
				<months>
					<alias source="locale" path="../../calendar[@type='gregorian']/months"/>
				</months>
			</calendar>
			<calendar type="gregorian">
				<months>
					<default choice="format"/>
					<monthContext type="format">
						<default choice="wide"/>
						<monthWidth type="narrow">
							<alias source="locale" path="../../monthContext[@type='stand-alone']/monthWidth[@type='narrow']"/>
						</monthWidth>
						<monthWidth type="wide">
							<month type="1">11</month>
							<month type="2">22</month>
							<month type="3">33</month>
							<month type="4">44</month>
						</monthWidth>
					</monthContext>
					<monthContext type="stand-alone">
						<monthWidth type="narrow">
							<month type="1">1</month>
							<month type="2">2</month>
							<month type="3">3</month>
							<month type="4">4</month>
						</monthWidth>
						<monthWidth type="wide">
							<alias source="locale" path="../../monthContext[@type='format']/monthWidth[@type='wide']"/>
						</monthWidth>
					</monthContext>
				</months>
			</calendar>
		</calendars>
	</dates>
</ldml>
`

var de_xml = `<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE ldml SYSTEM "../../common/dtd/ldml.dtd">
<ldml>
	<identity>
		<language type="de"/>
	</identity>
	<characters>
		<exemplarCharacters>[a ä b c d e ö p q r s ß t u ü v w x y z]</exemplarCharacters>
		<exemplarCharacters type="auxiliary">[á à ă]</exemplarCharacters>
		<exemplarCharacters type="index">[A B C D E F G H Z]</exemplarCharacters>
		<ellipsis type="final">***REMOVED***0***REMOVED*** …</ellipsis>
		<ellipsis type="initial">… ***REMOVED***0***REMOVED***</ellipsis>
		<moreInformation>?</moreInformation>
		<stopwords>
			<stopwordList type="collation" draft="provisional">der die das</stopwordList>
		</stopwords>
	</characters>
	<dates>
		<calendars>
			<calendar type="buddhist">
				<months>
					<monthContext type="format">
						<monthWidth type="narrow">
							<month type="3">BBB</month>
						</monthWidth>
						<monthWidth type="wide">
							<month type="3">bbb</month>
						</monthWidth>
					</monthContext>
				</months>
			</calendar>
			<calendar type="gregorian">
				<months>
					<monthContext type="format">
						<monthWidth type="narrow">
							<month type="3">M</month>
							<month type="4">A</month>
						</monthWidth>
						<monthWidth type="wide">
							<month type="3">Maerz</month>
							<month type="4">April</month>
							<month type="5">Mai</month>
						</monthWidth>
					</monthContext>
					<monthContext type="stand-alone">
						<monthWidth type="narrow">
							<month type="3">m</month>
							<month type="5">m</month>
						</monthWidth>
						<monthWidth type="wide">
							<month type="4">april</month>
							<month type="5">mai</month>
						</monthWidth>
					</monthContext>
				</months>
			</calendar>
		</calendars>
	</dates>
	<posix>
		<messages>
			<yesstr>yes:y</yesstr>
			<nostr>no:n</nostr>
		</messages>
	</posix>
</ldml>
`

var gsw_xml = `<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE ldml SYSTEM "../../common/dtd/ldml.dtd">
<ldml>
	<identity>
		<language type="gsw"/>
	</identity>
	<posix>
		<alias source="de" path="//ldml/posix"/>
	</posix>
</ldml>
`
