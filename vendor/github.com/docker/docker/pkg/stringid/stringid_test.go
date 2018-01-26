package stringid

import (
	"strings"
	"testing"
)

func TestGenerateRandomID(t *testing.T) ***REMOVED***
	id := GenerateRandomID()

	if len(id) != 64 ***REMOVED***
		t.Fatalf("Id returned is incorrect: %s", id)
	***REMOVED***
***REMOVED***

func TestGenerateNonCryptoID(t *testing.T) ***REMOVED***
	id := GenerateNonCryptoID()

	if len(id) != 64 ***REMOVED***
		t.Fatalf("Id returned is incorrect: %s", id)
	***REMOVED***
***REMOVED***

func TestShortenId(t *testing.T) ***REMOVED***
	id := "90435eec5c4e124e741ef731e118be2fc799a68aba0466ec17717f24ce2ae6a2"
	truncID := TruncateID(id)
	if truncID != "90435eec5c4e" ***REMOVED***
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	***REMOVED***
***REMOVED***

func TestShortenSha256Id(t *testing.T) ***REMOVED***
	id := "sha256:4e38e38c8ce0b8d9041a9c4fefe786631d1416225e13b0bfe8cfa2321aec4bba"
	truncID := TruncateID(id)
	if truncID != "4e38e38c8ce0" ***REMOVED***
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	***REMOVED***
***REMOVED***

func TestShortenIdEmpty(t *testing.T) ***REMOVED***
	id := ""
	truncID := TruncateID(id)
	if len(truncID) > len(id) ***REMOVED***
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	***REMOVED***
***REMOVED***

func TestShortenIdInvalid(t *testing.T) ***REMOVED***
	id := "1234"
	truncID := TruncateID(id)
	if len(truncID) != len(id) ***REMOVED***
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	***REMOVED***
***REMOVED***

func TestIsShortIDNonHex(t *testing.T) ***REMOVED***
	id := "some non-hex value"
	if IsShortID(id) ***REMOVED***
		t.Fatalf("%s is not a short ID", id)
	***REMOVED***
***REMOVED***

func TestIsShortIDNotCorrectSize(t *testing.T) ***REMOVED***
	id := strings.Repeat("a", shortLen+1)
	if IsShortID(id) ***REMOVED***
		t.Fatalf("%s is not a short ID", id)
	***REMOVED***
	id = strings.Repeat("a", shortLen-1)
	if IsShortID(id) ***REMOVED***
		t.Fatalf("%s is not a short ID", id)
	***REMOVED***
***REMOVED***
