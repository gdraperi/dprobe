// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package display_test

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
)

func ExampleFormatter() ***REMOVED***
	message.SetString(language.Dutch, "In %v people speak %v.", "In %v spreekt men %v.")

	fr := language.French
	region, _ := fr.Region()
	for _, tag := range []string***REMOVED***"en", "nl"***REMOVED*** ***REMOVED***
		p := message.NewPrinter(language.Make(tag))

		p.Printf("In %v people speak %v.", display.Region(region), display.Language(fr))
		p.Println()
	***REMOVED***

	// Output:
	// In France people speak French.
	// In Frankrijk spreekt men Frans.
***REMOVED***

func ExampleNamer() ***REMOVED***
	supported := []string***REMOVED***
		"en-US", "en-GB", "ja", "zh", "zh-Hans", "zh-Hant", "pt", "pt-PT", "ko", "ar", "el", "ru", "uk", "pa",
	***REMOVED***

	en := display.English.Languages()

	for _, s := range supported ***REMOVED***
		t := language.MustParse(s)
		fmt.Printf("%-20s (%s)\n", en.Name(t), display.Self.Name(t))
	***REMOVED***

	// Output:
	// American English     (American English)
	// British English      (British English)
	// Japanese             (日本語)
	// Chinese              (中文)
	// Simplified Chinese   (简体中文)
	// Traditional Chinese  (繁體中文)
	// Portuguese           (português)
	// European Portuguese  (português europeu)
	// Korean               (한국어)
	// Arabic               (العربية)
	// Greek                (Ελληνικά)
	// Russian              (русский)
	// Ukrainian            (українська)
	// Punjabi              (ਪੰਜਾਬੀ)
***REMOVED***

func ExampleTags() ***REMOVED***
	n := display.Tags(language.English)
	fmt.Println(n.Name(language.Make("nl")))
	fmt.Println(n.Name(language.Make("nl-BE")))
	fmt.Println(n.Name(language.Make("nl-CW")))
	fmt.Println(n.Name(language.Make("nl-Arab")))
	fmt.Println(n.Name(language.Make("nl-Cyrl-RU")))

	// Output:
	// Dutch
	// Flemish
	// Dutch (Curaçao)
	// Dutch (Arabic)
	// Dutch (Cyrillic, Russia)
***REMOVED***

// ExampleDictionary shows how to reduce the amount of data linked into your
// binary by only using the predefined Dictionary variables of the languages you
// wish to support.
func ExampleDictionary() ***REMOVED***
	tags := []language.Tag***REMOVED***
		language.English,
		language.German,
		language.Japanese,
		language.Russian,
	***REMOVED***
	dicts := []*display.Dictionary***REMOVED***
		display.English,
		display.German,
		display.Japanese,
		display.Russian,
	***REMOVED***

	m := language.NewMatcher(tags)

	getDict := func(t language.Tag) *display.Dictionary ***REMOVED***
		_, i, confidence := m.Match(t)
		// Skip this check if you want to support a fall-back language, which
		// will be the first one passed to NewMatcher.
		if confidence == language.No ***REMOVED***
			return nil
		***REMOVED***
		return dicts[i]
	***REMOVED***

	// The matcher will match Swiss German to German.
	n := getDict(language.Make("gsw")).Languages()
	fmt.Println(n.Name(language.German))
	fmt.Println(n.Name(language.Make("de-CH")))
	fmt.Println(n.Name(language.Make("gsw")))

	// Output:
	// Deutsch
	// Schweizer Hochdeutsch
	// Schweizerdeutsch
***REMOVED***
