// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package plural

type pluralTest struct ***REMOVED***
	locales string
	form    int
	integer []string
	decimal []string
***REMOVED***

var ordinalTests = []pluralTest***REMOVED*** // 66 elements
	0:  ***REMOVED***locales: "af am ar bg bs ce cs da de dsb el es et eu fa fi fy gl gsw he hr hsb id in is iw ja km kn ko ky lt lv ml mn my nb nl pa pl prg ps pt root ru sd sh si sk sl sr sw ta te th tr ur uz yue zh zu", form: 0, integer: []string***REMOVED***"0~15", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	1:  ***REMOVED***locales: "sv", form: 2, integer: []string***REMOVED***"1", "2", "21", "22", "31", "32", "41", "42", "51", "52", "61", "62", "71", "72", "81", "82", "101", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	2:  ***REMOVED***locales: "sv", form: 0, integer: []string***REMOVED***"0", "3~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	3:  ***REMOVED***locales: "fil fr ga hy lo mo ms ro tl vi", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	4:  ***REMOVED***locales: "fil fr ga hy lo mo ms ro tl vi", form: 0, integer: []string***REMOVED***"0", "2~16", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	5:  ***REMOVED***locales: "hu", form: 2, integer: []string***REMOVED***"1", "5"***REMOVED***, decimal: []string(nil)***REMOVED***,
	6:  ***REMOVED***locales: "hu", form: 0, integer: []string***REMOVED***"0", "2~4", "6~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	7:  ***REMOVED***locales: "ne", form: 2, integer: []string***REMOVED***"1~4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	8:  ***REMOVED***locales: "ne", form: 0, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	9:  ***REMOVED***locales: "be", form: 4, integer: []string***REMOVED***"2", "3", "22", "23", "32", "33", "42", "43", "52", "53", "62", "63", "72", "73", "82", "83", "102", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	10: ***REMOVED***locales: "be", form: 0, integer: []string***REMOVED***"0", "1", "4~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	11: ***REMOVED***locales: "uk", form: 4, integer: []string***REMOVED***"3", "23", "33", "43", "53", "63", "73", "83", "103", "1003"***REMOVED***, decimal: []string(nil)***REMOVED***,
	12: ***REMOVED***locales: "uk", form: 0, integer: []string***REMOVED***"0~2", "4~16", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	13: ***REMOVED***locales: "tk", form: 4, integer: []string***REMOVED***"6", "9", "10", "16", "19", "26", "29", "36", "39", "106", "1006"***REMOVED***, decimal: []string(nil)***REMOVED***,
	14: ***REMOVED***locales: "tk", form: 0, integer: []string***REMOVED***"0~5", "7", "8", "11~15", "17", "18", "20", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	15: ***REMOVED***locales: "kk", form: 5, integer: []string***REMOVED***"6", "9", "10", "16", "19", "20", "26", "29", "30", "36", "39", "40", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	16: ***REMOVED***locales: "kk", form: 0, integer: []string***REMOVED***"0~5", "7", "8", "11~15", "17", "18", "21", "101", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	17: ***REMOVED***locales: "it", form: 5, integer: []string***REMOVED***"8", "11", "80", "800"***REMOVED***, decimal: []string(nil)***REMOVED***,
	18: ***REMOVED***locales: "it", form: 0, integer: []string***REMOVED***"0~7", "9", "10", "12~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	19: ***REMOVED***locales: "ka", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	20: ***REMOVED***locales: "ka", form: 5, integer: []string***REMOVED***"0", "2~16", "102", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	21: ***REMOVED***locales: "ka", form: 0, integer: []string***REMOVED***"21~36", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	22: ***REMOVED***locales: "sq", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	23: ***REMOVED***locales: "sq", form: 5, integer: []string***REMOVED***"4", "24", "34", "44", "54", "64", "74", "84", "104", "1004"***REMOVED***, decimal: []string(nil)***REMOVED***,
	24: ***REMOVED***locales: "sq", form: 0, integer: []string***REMOVED***"0", "2", "3", "5~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	25: ***REMOVED***locales: "en", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	26: ***REMOVED***locales: "en", form: 3, integer: []string***REMOVED***"2", "22", "32", "42", "52", "62", "72", "82", "102", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	27: ***REMOVED***locales: "en", form: 4, integer: []string***REMOVED***"3", "23", "33", "43", "53", "63", "73", "83", "103", "1003"***REMOVED***, decimal: []string(nil)***REMOVED***,
	28: ***REMOVED***locales: "en", form: 0, integer: []string***REMOVED***"0", "4~18", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	29: ***REMOVED***locales: "mr", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	30: ***REMOVED***locales: "mr", form: 3, integer: []string***REMOVED***"2", "3"***REMOVED***, decimal: []string(nil)***REMOVED***,
	31: ***REMOVED***locales: "mr", form: 4, integer: []string***REMOVED***"4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	32: ***REMOVED***locales: "mr", form: 0, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	33: ***REMOVED***locales: "ca", form: 2, integer: []string***REMOVED***"1", "3"***REMOVED***, decimal: []string(nil)***REMOVED***,
	34: ***REMOVED***locales: "ca", form: 3, integer: []string***REMOVED***"2"***REMOVED***, decimal: []string(nil)***REMOVED***,
	35: ***REMOVED***locales: "ca", form: 4, integer: []string***REMOVED***"4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	36: ***REMOVED***locales: "ca", form: 0, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	37: ***REMOVED***locales: "mk", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	38: ***REMOVED***locales: "mk", form: 3, integer: []string***REMOVED***"2", "22", "32", "42", "52", "62", "72", "82", "102", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	39: ***REMOVED***locales: "mk", form: 5, integer: []string***REMOVED***"7", "8", "27", "28", "37", "38", "47", "48", "57", "58", "67", "68", "77", "78", "87", "88", "107", "1007"***REMOVED***, decimal: []string(nil)***REMOVED***,
	40: ***REMOVED***locales: "mk", form: 0, integer: []string***REMOVED***"0", "3~6", "9~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	41: ***REMOVED***locales: "az", form: 2, integer: []string***REMOVED***"1", "2", "5", "7", "8", "11", "12", "15", "17", "18", "20~22", "25", "101", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	42: ***REMOVED***locales: "az", form: 4, integer: []string***REMOVED***"3", "4", "13", "14", "23", "24", "33", "34", "43", "44", "53", "54", "63", "64", "73", "74", "100", "1003"***REMOVED***, decimal: []string(nil)***REMOVED***,
	43: ***REMOVED***locales: "az", form: 5, integer: []string***REMOVED***"0", "6", "16", "26", "36", "40", "46", "56", "106", "1006"***REMOVED***, decimal: []string(nil)***REMOVED***,
	44: ***REMOVED***locales: "az", form: 0, integer: []string***REMOVED***"9", "10", "19", "29", "30", "39", "49", "59", "69", "79", "109", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	45: ***REMOVED***locales: "gu hi", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	46: ***REMOVED***locales: "gu hi", form: 3, integer: []string***REMOVED***"2", "3"***REMOVED***, decimal: []string(nil)***REMOVED***,
	47: ***REMOVED***locales: "gu hi", form: 4, integer: []string***REMOVED***"4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	48: ***REMOVED***locales: "gu hi", form: 5, integer: []string***REMOVED***"6"***REMOVED***, decimal: []string(nil)***REMOVED***,
	49: ***REMOVED***locales: "gu hi", form: 0, integer: []string***REMOVED***"0", "5", "7~20", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	50: ***REMOVED***locales: "as bn", form: 2, integer: []string***REMOVED***"1", "5", "7~10"***REMOVED***, decimal: []string(nil)***REMOVED***,
	51: ***REMOVED***locales: "as bn", form: 3, integer: []string***REMOVED***"2", "3"***REMOVED***, decimal: []string(nil)***REMOVED***,
	52: ***REMOVED***locales: "as bn", form: 4, integer: []string***REMOVED***"4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	53: ***REMOVED***locales: "as bn", form: 5, integer: []string***REMOVED***"6"***REMOVED***, decimal: []string(nil)***REMOVED***,
	54: ***REMOVED***locales: "as bn", form: 0, integer: []string***REMOVED***"0", "11~25", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	55: ***REMOVED***locales: "or", form: 2, integer: []string***REMOVED***"1", "5", "7~9"***REMOVED***, decimal: []string(nil)***REMOVED***,
	56: ***REMOVED***locales: "or", form: 3, integer: []string***REMOVED***"2", "3"***REMOVED***, decimal: []string(nil)***REMOVED***,
	57: ***REMOVED***locales: "or", form: 4, integer: []string***REMOVED***"4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	58: ***REMOVED***locales: "or", form: 5, integer: []string***REMOVED***"6"***REMOVED***, decimal: []string(nil)***REMOVED***,
	59: ***REMOVED***locales: "or", form: 0, integer: []string***REMOVED***"0", "10~24", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	60: ***REMOVED***locales: "cy", form: 1, integer: []string***REMOVED***"0", "7~9"***REMOVED***, decimal: []string(nil)***REMOVED***,
	61: ***REMOVED***locales: "cy", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	62: ***REMOVED***locales: "cy", form: 3, integer: []string***REMOVED***"2"***REMOVED***, decimal: []string(nil)***REMOVED***,
	63: ***REMOVED***locales: "cy", form: 4, integer: []string***REMOVED***"3", "4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	64: ***REMOVED***locales: "cy", form: 5, integer: []string***REMOVED***"5", "6"***REMOVED***, decimal: []string(nil)***REMOVED***,
	65: ***REMOVED***locales: "cy", form: 0, integer: []string***REMOVED***"10~25", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
***REMOVED*** // Size: 4776 bytes

var cardinalTests = []pluralTest***REMOVED*** // 113 elements
	0:   ***REMOVED***locales: "bm bo dz id ig ii in ja jbo jv jw kde kea km ko lkt lo ms my nqo root sah ses sg th to vi wo yo yue zh", form: 0, integer: []string***REMOVED***"0~15", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	1:   ***REMOVED***locales: "am as bn fa gu hi kn mr zu", form: 2, integer: []string***REMOVED***"0", "1"***REMOVED***, decimal: []string***REMOVED***"0.0~1.0", "0.00~0.04"***REMOVED******REMOVED***,
	2:   ***REMOVED***locales: "am as bn fa gu hi kn mr zu", form: 0, integer: []string***REMOVED***"2~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"1.1~2.6", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	3:   ***REMOVED***locales: "ff fr hy kab", form: 2, integer: []string***REMOVED***"0", "1"***REMOVED***, decimal: []string***REMOVED***"0.0~1.5"***REMOVED******REMOVED***,
	4:   ***REMOVED***locales: "ff fr hy kab", form: 0, integer: []string***REMOVED***"2~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"2.0~3.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	5:   ***REMOVED***locales: "pt", form: 2, integer: []string***REMOVED***"0", "1"***REMOVED***, decimal: []string***REMOVED***"0.0~1.5"***REMOVED******REMOVED***,
	6:   ***REMOVED***locales: "pt", form: 0, integer: []string***REMOVED***"2~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"2.0~3.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	7:   ***REMOVED***locales: "ast ca de en et fi fy gl io it ji nl pt_PT sv sw ur yi", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	8:   ***REMOVED***locales: "ast ca de en et fi fy gl io it ji nl pt_PT sv sw ur yi", form: 0, integer: []string***REMOVED***"0", "2~16", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	9:   ***REMOVED***locales: "si", form: 2, integer: []string***REMOVED***"0", "1"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.1", "1.0", "0.00", "0.01", "1.00", "0.000", "0.001", "1.000", "0.0000", "0.0001", "1.0000"***REMOVED******REMOVED***,
	10:  ***REMOVED***locales: "si", form: 0, integer: []string***REMOVED***"2~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.2~0.9", "1.1~1.8", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	11:  ***REMOVED***locales: "ak bh guw ln mg nso pa ti wa", form: 2, integer: []string***REMOVED***"0", "1"***REMOVED***, decimal: []string***REMOVED***"0.0", "1.0", "0.00", "1.00", "0.000", "1.000", "0.0000", "1.0000"***REMOVED******REMOVED***,
	12:  ***REMOVED***locales: "ak bh guw ln mg nso pa ti wa", form: 0, integer: []string***REMOVED***"2~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	13:  ***REMOVED***locales: "tzm", form: 2, integer: []string***REMOVED***"0", "1", "11~24"***REMOVED***, decimal: []string***REMOVED***"0.0", "1.0", "11.0", "12.0", "13.0", "14.0", "15.0", "16.0", "17.0", "18.0", "19.0", "20.0", "21.0", "22.0", "23.0", "24.0"***REMOVED******REMOVED***,
	14:  ***REMOVED***locales: "tzm", form: 0, integer: []string***REMOVED***"2~10", "100~106", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	15:  ***REMOVED***locales: "af asa az bem bez bg brx ce cgg chr ckb dv ee el eo es eu fo fur gsw ha haw hu jgo jmc ka kaj kcg kk kkj kl ks ksb ku ky lb lg mas mgo ml mn nah nb nd ne nn nnh no nr ny nyn om or os pap ps rm rof rwk saq sd sdh seh sn so sq ss ssy st syr ta te teo tig tk tn tr ts ug uz ve vo vun wae xh xog", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"1.0", "1.00", "1.000", "1.0000"***REMOVED******REMOVED***,
	16:  ***REMOVED***locales: "af asa az bem bez bg brx ce cgg chr ckb dv ee el eo es eu fo fur gsw ha haw hu jgo jmc ka kaj kcg kk kkj kl ks ksb ku ky lb lg mas mgo ml mn nah nb nd ne nn nnh no nr ny nyn om or os pap ps rm rof rwk saq sd sdh seh sn so sq ss ssy st syr ta te teo tig tk tn tr ts ug uz ve vo vun wae xh xog", form: 0, integer: []string***REMOVED***"0", "2~16", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0~0.9", "1.1~1.6", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	17:  ***REMOVED***locales: "da", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"0.1~1.6"***REMOVED******REMOVED***,
	18:  ***REMOVED***locales: "da", form: 0, integer: []string***REMOVED***"0", "2~16", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "2.0~3.4", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	19:  ***REMOVED***locales: "is", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"0.1~1.6", "10.1", "100.1", "1000.1"***REMOVED******REMOVED***,
	20:  ***REMOVED***locales: "is", form: 0, integer: []string***REMOVED***"0", "2~16", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "2.0", "3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	21:  ***REMOVED***locales: "mk", form: 2, integer: []string***REMOVED***"1", "11", "21", "31", "41", "51", "61", "71", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"0.1", "1.1", "2.1", "3.1", "4.1", "5.1", "6.1", "7.1", "10.1", "100.1", "1000.1"***REMOVED******REMOVED***,
	22:  ***REMOVED***locales: "mk", form: 0, integer: []string***REMOVED***"0", "2~10", "12~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.2~1.0", "1.2~1.7", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	23:  ***REMOVED***locales: "fil tl", form: 2, integer: []string***REMOVED***"0~3", "5", "7", "8", "10~13", "15", "17", "18", "20", "21", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0~0.3", "0.5", "0.7", "0.8", "1.0~1.3", "1.5", "1.7", "1.8", "2.0", "2.1", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	24:  ***REMOVED***locales: "fil tl", form: 0, integer: []string***REMOVED***"4", "6", "9", "14", "16", "19", "24", "26", "104", "1004"***REMOVED***, decimal: []string***REMOVED***"0.4", "0.6", "0.9", "1.4", "1.6", "1.9", "2.4", "2.6", "10.4", "100.4", "1000.4"***REMOVED******REMOVED***,
	25:  ***REMOVED***locales: "lv prg", form: 1, integer: []string***REMOVED***"0", "10~20", "30", "40", "50", "60", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "10.0", "11.0", "12.0", "13.0", "14.0", "15.0", "16.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	26:  ***REMOVED***locales: "lv prg", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"0.1", "1.0", "1.1", "2.1", "3.1", "4.1", "5.1", "6.1", "7.1", "10.1", "100.1", "1000.1"***REMOVED******REMOVED***,
	27:  ***REMOVED***locales: "lv prg", form: 0, integer: []string***REMOVED***"2~9", "22~29", "102", "1002"***REMOVED***, decimal: []string***REMOVED***"0.2~0.9", "1.2~1.9", "10.2", "100.2", "1000.2"***REMOVED******REMOVED***,
	28:  ***REMOVED***locales: "lag", form: 1, integer: []string***REMOVED***"0"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.00", "0.000", "0.0000"***REMOVED******REMOVED***,
	29:  ***REMOVED***locales: "lag", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"0.1~1.6"***REMOVED******REMOVED***,
	30:  ***REMOVED***locales: "lag", form: 0, integer: []string***REMOVED***"2~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"2.0~3.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	31:  ***REMOVED***locales: "ksh", form: 1, integer: []string***REMOVED***"0"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.00", "0.000", "0.0000"***REMOVED******REMOVED***,
	32:  ***REMOVED***locales: "ksh", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"1.0", "1.00", "1.000", "1.0000"***REMOVED******REMOVED***,
	33:  ***REMOVED***locales: "ksh", form: 0, integer: []string***REMOVED***"2~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	34:  ***REMOVED***locales: "iu kw naq se sma smi smj smn sms", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"1.0", "1.00", "1.000", "1.0000"***REMOVED******REMOVED***,
	35:  ***REMOVED***locales: "iu kw naq se sma smi smj smn sms", form: 3, integer: []string***REMOVED***"2"***REMOVED***, decimal: []string***REMOVED***"2.0", "2.00", "2.000", "2.0000"***REMOVED******REMOVED***,
	36:  ***REMOVED***locales: "iu kw naq se sma smi smj smn sms", form: 0, integer: []string***REMOVED***"0", "3~17", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0~0.9", "1.1~1.6", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	37:  ***REMOVED***locales: "shi", form: 2, integer: []string***REMOVED***"0", "1"***REMOVED***, decimal: []string***REMOVED***"0.0~1.0", "0.00~0.04"***REMOVED******REMOVED***,
	38:  ***REMOVED***locales: "shi", form: 4, integer: []string***REMOVED***"2~10"***REMOVED***, decimal: []string***REMOVED***"2.0", "3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0", "2.00", "3.00", "4.00", "5.00", "6.00", "7.00", "8.00"***REMOVED******REMOVED***,
	39:  ***REMOVED***locales: "shi", form: 0, integer: []string***REMOVED***"11~26", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"1.1~1.9", "2.1~2.7", "10.1", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	40:  ***REMOVED***locales: "mo ro", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	41:  ***REMOVED***locales: "mo ro", form: 4, integer: []string***REMOVED***"0", "2~16", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	42:  ***REMOVED***locales: "mo ro", form: 0, integer: []string***REMOVED***"20~35", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	43:  ***REMOVED***locales: "bs hr sh sr", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"0.1", "1.1", "2.1", "3.1", "4.1", "5.1", "6.1", "7.1", "10.1", "100.1", "1000.1"***REMOVED******REMOVED***,
	44:  ***REMOVED***locales: "bs hr sh sr", form: 4, integer: []string***REMOVED***"2~4", "22~24", "32~34", "42~44", "52~54", "62", "102", "1002"***REMOVED***, decimal: []string***REMOVED***"0.2~0.4", "1.2~1.4", "2.2~2.4", "3.2~3.4", "4.2~4.4", "5.2", "10.2", "100.2", "1000.2"***REMOVED******REMOVED***,
	45:  ***REMOVED***locales: "bs hr sh sr", form: 0, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.5~1.0", "1.5~2.0", "2.5~2.7", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	46:  ***REMOVED***locales: "gd", form: 2, integer: []string***REMOVED***"1", "11"***REMOVED***, decimal: []string***REMOVED***"1.0", "11.0", "1.00", "11.00", "1.000", "11.000", "1.0000"***REMOVED******REMOVED***,
	47:  ***REMOVED***locales: "gd", form: 3, integer: []string***REMOVED***"2", "12"***REMOVED***, decimal: []string***REMOVED***"2.0", "12.0", "2.00", "12.00", "2.000", "12.000", "2.0000"***REMOVED******REMOVED***,
	48:  ***REMOVED***locales: "gd", form: 4, integer: []string***REMOVED***"3~10", "13~19"***REMOVED***, decimal: []string***REMOVED***"3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0", "13.0", "14.0", "15.0", "16.0", "17.0", "18.0", "19.0", "3.00"***REMOVED******REMOVED***,
	49:  ***REMOVED***locales: "gd", form: 0, integer: []string***REMOVED***"0", "20~34", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0~0.9", "1.1~1.6", "10.1", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	50:  ***REMOVED***locales: "sl", form: 2, integer: []string***REMOVED***"1", "101", "201", "301", "401", "501", "601", "701", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	51:  ***REMOVED***locales: "sl", form: 3, integer: []string***REMOVED***"2", "102", "202", "302", "402", "502", "602", "702", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	52:  ***REMOVED***locales: "sl", form: 4, integer: []string***REMOVED***"3", "4", "103", "104", "203", "204", "303", "304", "403", "404", "503", "504", "603", "604", "703", "704", "1003"***REMOVED***, decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	53:  ***REMOVED***locales: "sl", form: 0, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	54:  ***REMOVED***locales: "dsb hsb", form: 2, integer: []string***REMOVED***"1", "101", "201", "301", "401", "501", "601", "701", "1001"***REMOVED***, decimal: []string***REMOVED***"0.1", "1.1", "2.1", "3.1", "4.1", "5.1", "6.1", "7.1", "10.1", "100.1", "1000.1"***REMOVED******REMOVED***,
	55:  ***REMOVED***locales: "dsb hsb", form: 3, integer: []string***REMOVED***"2", "102", "202", "302", "402", "502", "602", "702", "1002"***REMOVED***, decimal: []string***REMOVED***"0.2", "1.2", "2.2", "3.2", "4.2", "5.2", "6.2", "7.2", "10.2", "100.2", "1000.2"***REMOVED******REMOVED***,
	56:  ***REMOVED***locales: "dsb hsb", form: 4, integer: []string***REMOVED***"3", "4", "103", "104", "203", "204", "303", "304", "403", "404", "503", "504", "603", "604", "703", "704", "1003"***REMOVED***, decimal: []string***REMOVED***"0.3", "0.4", "1.3", "1.4", "2.3", "2.4", "3.3", "3.4", "4.3", "4.4", "5.3", "5.4", "6.3", "6.4", "7.3", "7.4", "10.3", "100.3", "1000.3"***REMOVED******REMOVED***,
	57:  ***REMOVED***locales: "dsb hsb", form: 0, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.5~1.0", "1.5~2.0", "2.5~2.7", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	58:  ***REMOVED***locales: "he iw", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	59:  ***REMOVED***locales: "he iw", form: 3, integer: []string***REMOVED***"2"***REMOVED***, decimal: []string(nil)***REMOVED***,
	60:  ***REMOVED***locales: "he iw", form: 5, integer: []string***REMOVED***"20", "30", "40", "50", "60", "70", "80", "90", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	61:  ***REMOVED***locales: "he iw", form: 0, integer: []string***REMOVED***"0", "3~17", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	62:  ***REMOVED***locales: "cs sk", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	63:  ***REMOVED***locales: "cs sk", form: 4, integer: []string***REMOVED***"2~4"***REMOVED***, decimal: []string(nil)***REMOVED***,
	64:  ***REMOVED***locales: "cs sk", form: 5, integer: []string(nil), decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	65:  ***REMOVED***locales: "cs sk", form: 0, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	66:  ***REMOVED***locales: "pl", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string(nil)***REMOVED***,
	67:  ***REMOVED***locales: "pl", form: 4, integer: []string***REMOVED***"2~4", "22~24", "32~34", "42~44", "52~54", "62", "102", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	68:  ***REMOVED***locales: "pl", form: 5, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	69:  ***REMOVED***locales: "pl", form: 0, integer: []string(nil), decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	70:  ***REMOVED***locales: "be", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"1.0", "21.0", "31.0", "41.0", "51.0", "61.0", "71.0", "81.0", "101.0", "1001.0"***REMOVED******REMOVED***,
	71:  ***REMOVED***locales: "be", form: 4, integer: []string***REMOVED***"2~4", "22~24", "32~34", "42~44", "52~54", "62", "102", "1002"***REMOVED***, decimal: []string***REMOVED***"2.0", "3.0", "4.0", "22.0", "23.0", "24.0", "32.0", "33.0", "102.0", "1002.0"***REMOVED******REMOVED***,
	72:  ***REMOVED***locales: "be", form: 5, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0", "11.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	73:  ***REMOVED***locales: "be", form: 0, integer: []string(nil), decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.1", "100.1", "1000.1"***REMOVED******REMOVED***,
	74:  ***REMOVED***locales: "lt", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"1.0", "21.0", "31.0", "41.0", "51.0", "61.0", "71.0", "81.0", "101.0", "1001.0"***REMOVED******REMOVED***,
	75:  ***REMOVED***locales: "lt", form: 4, integer: []string***REMOVED***"2~9", "22~29", "102", "1002"***REMOVED***, decimal: []string***REMOVED***"2.0", "3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "22.0", "102.0", "1002.0"***REMOVED******REMOVED***,
	76:  ***REMOVED***locales: "lt", form: 5, integer: []string(nil), decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.1", "100.1", "1000.1"***REMOVED******REMOVED***,
	77:  ***REMOVED***locales: "lt", form: 0, integer: []string***REMOVED***"0", "10~20", "30", "40", "50", "60", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0", "10.0", "11.0", "12.0", "13.0", "14.0", "15.0", "16.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	78:  ***REMOVED***locales: "mt", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"1.0", "1.00", "1.000", "1.0000"***REMOVED******REMOVED***,
	79:  ***REMOVED***locales: "mt", form: 4, integer: []string***REMOVED***"0", "2~10", "102~107", "1002"***REMOVED***, decimal: []string***REMOVED***"0.0", "2.0", "3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "10.0", "102.0", "1002.0"***REMOVED******REMOVED***,
	80:  ***REMOVED***locales: "mt", form: 5, integer: []string***REMOVED***"11~19", "111~117", "1011"***REMOVED***, decimal: []string***REMOVED***"11.0", "12.0", "13.0", "14.0", "15.0", "16.0", "17.0", "18.0", "111.0", "1011.0"***REMOVED******REMOVED***,
	81:  ***REMOVED***locales: "mt", form: 0, integer: []string***REMOVED***"20~35", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.1", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	82:  ***REMOVED***locales: "ru uk", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "71", "81", "101", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	83:  ***REMOVED***locales: "ru uk", form: 4, integer: []string***REMOVED***"2~4", "22~24", "32~34", "42~44", "52~54", "62", "102", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	84:  ***REMOVED***locales: "ru uk", form: 5, integer: []string***REMOVED***"0", "5~19", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	85:  ***REMOVED***locales: "ru uk", form: 0, integer: []string(nil), decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	86:  ***REMOVED***locales: "br", form: 2, integer: []string***REMOVED***"1", "21", "31", "41", "51", "61", "81", "101", "1001"***REMOVED***, decimal: []string***REMOVED***"1.0", "21.0", "31.0", "41.0", "51.0", "61.0", "81.0", "101.0", "1001.0"***REMOVED******REMOVED***,
	87:  ***REMOVED***locales: "br", form: 3, integer: []string***REMOVED***"2", "22", "32", "42", "52", "62", "82", "102", "1002"***REMOVED***, decimal: []string***REMOVED***"2.0", "22.0", "32.0", "42.0", "52.0", "62.0", "82.0", "102.0", "1002.0"***REMOVED******REMOVED***,
	88:  ***REMOVED***locales: "br", form: 4, integer: []string***REMOVED***"3", "4", "9", "23", "24", "29", "33", "34", "39", "43", "44", "49", "103", "1003"***REMOVED***, decimal: []string***REMOVED***"3.0", "4.0", "9.0", "23.0", "24.0", "29.0", "33.0", "34.0", "103.0", "1003.0"***REMOVED******REMOVED***,
	89:  ***REMOVED***locales: "br", form: 5, integer: []string***REMOVED***"1000000"***REMOVED***, decimal: []string***REMOVED***"1000000.0", "1000000.00", "1000000.000"***REMOVED******REMOVED***,
	90:  ***REMOVED***locales: "br", form: 0, integer: []string***REMOVED***"0", "5~8", "10~20", "100", "1000", "10000", "100000"***REMOVED***, decimal: []string***REMOVED***"0.0~0.9", "1.1~1.6", "10.0", "100.0", "1000.0", "10000.0", "100000.0"***REMOVED******REMOVED***,
	91:  ***REMOVED***locales: "ga", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"1.0", "1.00", "1.000", "1.0000"***REMOVED******REMOVED***,
	92:  ***REMOVED***locales: "ga", form: 3, integer: []string***REMOVED***"2"***REMOVED***, decimal: []string***REMOVED***"2.0", "2.00", "2.000", "2.0000"***REMOVED******REMOVED***,
	93:  ***REMOVED***locales: "ga", form: 4, integer: []string***REMOVED***"3~6"***REMOVED***, decimal: []string***REMOVED***"3.0", "4.0", "5.0", "6.0", "3.00", "4.00", "5.00", "6.00", "3.000", "4.000", "5.000", "6.000", "3.0000", "4.0000", "5.0000", "6.0000"***REMOVED******REMOVED***,
	94:  ***REMOVED***locales: "ga", form: 5, integer: []string***REMOVED***"7~10"***REMOVED***, decimal: []string***REMOVED***"7.0", "8.0", "9.0", "10.0", "7.00", "8.00", "9.00", "10.00", "7.000", "8.000", "9.000", "10.000", "7.0000", "8.0000", "9.0000", "10.0000"***REMOVED******REMOVED***,
	95:  ***REMOVED***locales: "ga", form: 0, integer: []string***REMOVED***"0", "11~25", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.0~0.9", "1.1~1.6", "10.1", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	96:  ***REMOVED***locales: "gv", form: 2, integer: []string***REMOVED***"1", "11", "21", "31", "41", "51", "61", "71", "101", "1001"***REMOVED***, decimal: []string(nil)***REMOVED***,
	97:  ***REMOVED***locales: "gv", form: 3, integer: []string***REMOVED***"2", "12", "22", "32", "42", "52", "62", "72", "102", "1002"***REMOVED***, decimal: []string(nil)***REMOVED***,
	98:  ***REMOVED***locales: "gv", form: 4, integer: []string***REMOVED***"0", "20", "40", "60", "80", "100", "120", "140", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string(nil)***REMOVED***,
	99:  ***REMOVED***locales: "gv", form: 5, integer: []string(nil), decimal: []string***REMOVED***"0.0~1.5", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	100: ***REMOVED***locales: "gv", form: 0, integer: []string***REMOVED***"3~10", "13~19", "23", "103", "1003"***REMOVED***, decimal: []string(nil)***REMOVED***,
	101: ***REMOVED***locales: "ar ars", form: 1, integer: []string***REMOVED***"0"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.00", "0.000", "0.0000"***REMOVED******REMOVED***,
	102: ***REMOVED***locales: "ar ars", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"1.0", "1.00", "1.000", "1.0000"***REMOVED******REMOVED***,
	103: ***REMOVED***locales: "ar ars", form: 3, integer: []string***REMOVED***"2"***REMOVED***, decimal: []string***REMOVED***"2.0", "2.00", "2.000", "2.0000"***REMOVED******REMOVED***,
	104: ***REMOVED***locales: "ar ars", form: 4, integer: []string***REMOVED***"3~10", "103~110", "1003"***REMOVED***, decimal: []string***REMOVED***"3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0", "103.0", "1003.0"***REMOVED******REMOVED***,
	105: ***REMOVED***locales: "ar ars", form: 5, integer: []string***REMOVED***"11~26", "111", "1011"***REMOVED***, decimal: []string***REMOVED***"11.0", "12.0", "13.0", "14.0", "15.0", "16.0", "17.0", "18.0", "111.0", "1011.0"***REMOVED******REMOVED***,
	106: ***REMOVED***locales: "ar ars", form: 0, integer: []string***REMOVED***"100~102", "200~202", "300~302", "400~402", "500~502", "600", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.1", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
	107: ***REMOVED***locales: "cy", form: 1, integer: []string***REMOVED***"0"***REMOVED***, decimal: []string***REMOVED***"0.0", "0.00", "0.000", "0.0000"***REMOVED******REMOVED***,
	108: ***REMOVED***locales: "cy", form: 2, integer: []string***REMOVED***"1"***REMOVED***, decimal: []string***REMOVED***"1.0", "1.00", "1.000", "1.0000"***REMOVED******REMOVED***,
	109: ***REMOVED***locales: "cy", form: 3, integer: []string***REMOVED***"2"***REMOVED***, decimal: []string***REMOVED***"2.0", "2.00", "2.000", "2.0000"***REMOVED******REMOVED***,
	110: ***REMOVED***locales: "cy", form: 4, integer: []string***REMOVED***"3"***REMOVED***, decimal: []string***REMOVED***"3.0", "3.00", "3.000", "3.0000"***REMOVED******REMOVED***,
	111: ***REMOVED***locales: "cy", form: 5, integer: []string***REMOVED***"6"***REMOVED***, decimal: []string***REMOVED***"6.0", "6.00", "6.000", "6.0000"***REMOVED******REMOVED***,
	112: ***REMOVED***locales: "cy", form: 0, integer: []string***REMOVED***"4", "5", "7~20", "100", "1000", "10000", "100000", "1000000"***REMOVED***, decimal: []string***REMOVED***"0.1~0.9", "1.1~1.7", "10.0", "100.0", "1000.0", "10000.0", "100000.0", "1000000.0"***REMOVED******REMOVED***,
***REMOVED*** // Size: 8160 bytes

// Total table size 12936 bytes (12KiB); checksum: 8456DC5D
