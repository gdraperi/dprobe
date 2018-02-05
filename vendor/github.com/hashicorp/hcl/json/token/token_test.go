package token

import (
	"testing"
)

func TestTypeString(t *testing.T) ***REMOVED***
	var tokens = []struct ***REMOVED***
		tt  Type
		str string
	***REMOVED******REMOVED***
		***REMOVED***ILLEGAL, "ILLEGAL"***REMOVED***,
		***REMOVED***EOF, "EOF"***REMOVED***,
		***REMOVED***NUMBER, "NUMBER"***REMOVED***,
		***REMOVED***FLOAT, "FLOAT"***REMOVED***,
		***REMOVED***BOOL, "BOOL"***REMOVED***,
		***REMOVED***STRING, "STRING"***REMOVED***,
		***REMOVED***NULL, "NULL"***REMOVED***,
		***REMOVED***LBRACK, "LBRACK"***REMOVED***,
		***REMOVED***LBRACE, "LBRACE"***REMOVED***,
		***REMOVED***COMMA, "COMMA"***REMOVED***,
		***REMOVED***PERIOD, "PERIOD"***REMOVED***,
		***REMOVED***RBRACK, "RBRACK"***REMOVED***,
		***REMOVED***RBRACE, "RBRACE"***REMOVED***,
	***REMOVED***

	for _, token := range tokens ***REMOVED***
		if token.tt.String() != token.str ***REMOVED***
			t.Errorf("want: %q got:%q\n", token.str, token.tt)

		***REMOVED***
	***REMOVED***

***REMOVED***
