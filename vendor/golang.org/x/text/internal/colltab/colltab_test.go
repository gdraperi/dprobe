package colltab

import (
	"testing"

	"golang.org/x/text/language"
)

func TestMatchLang(t *testing.T) ***REMOVED***
	tags := []language.Tag***REMOVED***
		0:  language.Und,
		1:  language.MustParse("bs"),
		2:  language.German,
		3:  language.English,
		4:  language.AmericanEnglish,
		5:  language.MustParse("en-US-u-va-posix"),
		6:  language.Portuguese,
		7:  language.Serbian,
		8:  language.MustParse("sr-Latn"),
		9:  language.Chinese,
		10: language.MustParse("zh-u-co-stroke"),
		11: language.MustParse("zh-Hant-u-co-pinyin"),
		12: language.TraditionalChinese,
	***REMOVED***
	for i, tc := range []struct ***REMOVED***
		x int
		t language.Tag
	***REMOVED******REMOVED***
		***REMOVED***0, language.Und***REMOVED***,
		***REMOVED***0, language.Persian***REMOVED***, // Default to first element when no match.
		***REMOVED***3, language.English***REMOVED***,
		***REMOVED***4, language.AmericanEnglish***REMOVED***,
		***REMOVED***5, language.MustParse("en-US-u-va-posix")***REMOVED***,   // Ext. variant match.
		***REMOVED***4, language.MustParse("en-US-u-va-noposix")***REMOVED***, // Ext. variant mismatch.
		***REMOVED***3, language.MustParse("en-UK-u-va-noposix")***REMOVED***, // Ext. variant mismatch.
		***REMOVED***7, language.Serbian***REMOVED***,
		***REMOVED***0, language.Croatian***REMOVED***,             // Don't match to close language!
		***REMOVED***0, language.MustParse("gsw")***REMOVED***,     // Don't match to close language!
		***REMOVED***1, language.MustParse("bs-Cyrl")***REMOVED***, // Odd, but correct.
		***REMOVED***1, language.MustParse("bs-Latn")***REMOVED***, // Estimated script drops.
		***REMOVED***8, language.MustParse("sr-Latn")***REMOVED***,
		***REMOVED***9, language.Chinese***REMOVED***,
		***REMOVED***9, language.SimplifiedChinese***REMOVED***,
		***REMOVED***12, language.TraditionalChinese***REMOVED***,
		***REMOVED***11, language.MustParse("zh-Hant-u-co-pinyin")***REMOVED***,
		// TODO: should this be 12? Either inherited value (10) or default is
		// fine in this case, though. Other locales are not affected.
		***REMOVED***10, language.MustParse("zh-Hant-u-co-stroke")***REMOVED***,
		// There is no "phonebk" sorting order for zh-Hant, so use default.
		***REMOVED***12, language.MustParse("zh-Hant-u-co-phonebk")***REMOVED***,
		***REMOVED***10, language.MustParse("zh-u-co-stroke")***REMOVED***,
		***REMOVED***12, language.MustParse("und-TW")***REMOVED***,     // Infer script and language.
		***REMOVED***12, language.MustParse("und-HK")***REMOVED***,     // Infer script and language.
		***REMOVED***6, language.MustParse("und-BR")***REMOVED***,      // Infer script and language.
		***REMOVED***6, language.MustParse("und-PT")***REMOVED***,      // Infer script and language.
		***REMOVED***2, language.MustParse("und-Latn-DE")***REMOVED***, // Infer language.
		***REMOVED***0, language.MustParse("und-Jpan-BR")***REMOVED***, // Infers "ja", so no match.
		***REMOVED***0, language.MustParse("zu")***REMOVED***,          // No match past index.
	***REMOVED*** ***REMOVED***
		if x := MatchLang(tc.t, tags); x != tc.x ***REMOVED***
			t.Errorf("%d: MatchLang(%q, tags) = %d; want %d", i, tc.t, x, tc.x)
		***REMOVED***
	***REMOVED***
***REMOVED***
