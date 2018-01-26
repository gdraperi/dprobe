package term

import (
	"fmt"
	"strings"
)

// ASCII list the possible supported ASCII key sequence
var ASCII = []string***REMOVED***
	"ctrl-@",
	"ctrl-a",
	"ctrl-b",
	"ctrl-c",
	"ctrl-d",
	"ctrl-e",
	"ctrl-f",
	"ctrl-g",
	"ctrl-h",
	"ctrl-i",
	"ctrl-j",
	"ctrl-k",
	"ctrl-l",
	"ctrl-m",
	"ctrl-n",
	"ctrl-o",
	"ctrl-p",
	"ctrl-q",
	"ctrl-r",
	"ctrl-s",
	"ctrl-t",
	"ctrl-u",
	"ctrl-v",
	"ctrl-w",
	"ctrl-x",
	"ctrl-y",
	"ctrl-z",
	"ctrl-[",
	"ctrl-\\",
	"ctrl-]",
	"ctrl-^",
	"ctrl-_",
***REMOVED***

// ToBytes converts a string representing a suite of key-sequence to the corresponding ASCII code.
func ToBytes(keys string) ([]byte, error) ***REMOVED***
	codes := []byte***REMOVED******REMOVED***
next:
	for _, key := range strings.Split(keys, ",") ***REMOVED***
		if len(key) != 1 ***REMOVED***
			for code, ctrl := range ASCII ***REMOVED***
				if ctrl == key ***REMOVED***
					codes = append(codes, byte(code))
					continue next
				***REMOVED***
			***REMOVED***
			if key == "DEL" ***REMOVED***
				codes = append(codes, 127)
			***REMOVED*** else ***REMOVED***
				return nil, fmt.Errorf("Unknown character: '%s'", key)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			codes = append(codes, key[0])
		***REMOVED***
	***REMOVED***
	return codes, nil
***REMOVED***
