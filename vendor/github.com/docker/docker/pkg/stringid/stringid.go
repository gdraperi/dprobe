// Package stringid provides helper functions for dealing with string identifiers
package stringid

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const shortLen = 12

var (
	validShortID = regexp.MustCompile("^[a-f0-9]***REMOVED***12***REMOVED***$")
	validHex     = regexp.MustCompile(`^[a-f0-9]***REMOVED***64***REMOVED***$`)
)

// IsShortID determines if an arbitrary string *looks like* a short ID.
func IsShortID(id string) bool ***REMOVED***
	return validShortID.MatchString(id)
***REMOVED***

// TruncateID returns a shorthand version of a string identifier for convenience.
// A collision with other shorthands is very unlikely, but possible.
// In case of a collision a lookup with TruncIndex.Get() will fail, and the caller
// will need to use a longer prefix, or the full-length Id.
func TruncateID(id string) string ***REMOVED***
	if i := strings.IndexRune(id, ':'); i >= 0 ***REMOVED***
		id = id[i+1:]
	***REMOVED***
	if len(id) > shortLen ***REMOVED***
		id = id[:shortLen]
	***REMOVED***
	return id
***REMOVED***

func generateID(r io.Reader) string ***REMOVED***
	b := make([]byte, 32)
	for ***REMOVED***
		if _, err := io.ReadFull(r, b); err != nil ***REMOVED***
			panic(err) // This shouldn't happen
		***REMOVED***
		id := hex.EncodeToString(b)
		// if we try to parse the truncated for as an int and we don't have
		// an error then the value is all numeric and causes issues when
		// used as a hostname. ref #3869
		if _, err := strconv.ParseInt(TruncateID(id), 10, 64); err == nil ***REMOVED***
			continue
		***REMOVED***
		return id
	***REMOVED***
***REMOVED***

// GenerateRandomID returns a unique id.
func GenerateRandomID() string ***REMOVED***
	return generateID(cryptorand.Reader)
***REMOVED***

// GenerateNonCryptoID generates unique id without using cryptographically
// secure sources of random.
// It helps you to save entropy.
func GenerateNonCryptoID() string ***REMOVED***
	return generateID(readerFunc(rand.Read))
***REMOVED***

// ValidateID checks whether an ID string is a valid image ID.
func ValidateID(id string) error ***REMOVED***
	if ok := validHex.MatchString(id); !ok ***REMOVED***
		return fmt.Errorf("image ID %q is invalid", id)
	***REMOVED***
	return nil
***REMOVED***

func init() ***REMOVED***
	// safely set the seed globally so we generate random ids. Tries to use a
	// crypto seed before falling back to time.
	var seed int64
	if cryptoseed, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64)); err != nil ***REMOVED***
		// This should not happen, but worst-case fallback to time-based seed.
		seed = time.Now().UnixNano()
	***REMOVED*** else ***REMOVED***
		seed = cryptoseed.Int64()
	***REMOVED***

	rand.Seed(seed)
***REMOVED***

type readerFunc func(p []byte) (int, error)

func (fn readerFunc) Read(p []byte) (int, error) ***REMOVED***
	return fn(p)
***REMOVED***
