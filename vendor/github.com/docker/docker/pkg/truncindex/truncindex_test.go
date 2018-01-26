package truncindex

import (
	"math/rand"
	"testing"
	"time"

	"github.com/docker/docker/pkg/stringid"
)

// Test the behavior of TruncIndex, an index for querying IDs from a non-conflicting prefix.
func TestTruncIndex(t *testing.T) ***REMOVED***
	ids := []string***REMOVED******REMOVED***
	index := NewTruncIndex(ids)
	// Get on an empty index
	if _, err := index.Get("foobar"); err == nil ***REMOVED***
		t.Fatal("Get on an empty index should return an error")
	***REMOVED***

	// Spaces should be illegal in an id
	if err := index.Add("I have a space"); err == nil ***REMOVED***
		t.Fatalf("Adding an id with ' ' should return an error")
	***REMOVED***

	id := "99b36c2c326ccc11e726eee6ee78a0baf166ef96"
	// Add an id
	if err := index.Add(id); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Add an empty id (should fail)
	if err := index.Add(""); err == nil ***REMOVED***
		t.Fatalf("Adding an empty id should return an error")
	***REMOVED***

	// Get a non-existing id
	assertIndexGet(t, index, "abracadabra", "", true)
	// Get an empty id
	assertIndexGet(t, index, "", "", true)
	// Get the exact id
	assertIndexGet(t, index, id, id, false)
	// The first letter should match
	assertIndexGet(t, index, id[:1], id, false)
	// The first half should match
	assertIndexGet(t, index, id[:len(id)/2], id, false)
	// The second half should NOT match
	assertIndexGet(t, index, id[len(id)/2:], "", true)

	id2 := id[:6] + "blabla"
	// Add an id
	if err := index.Add(id2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// Both exact IDs should work
	assertIndexGet(t, index, id, id, false)
	assertIndexGet(t, index, id2, id2, false)

	// 6 characters or less should conflict
	assertIndexGet(t, index, id[:6], "", true)
	assertIndexGet(t, index, id[:4], "", true)
	assertIndexGet(t, index, id[:1], "", true)

	// An ambiguous id prefix should return an error
	if _, err := index.Get(id[:4]); err == nil ***REMOVED***
		t.Fatal("An ambiguous id prefix should return an error")
	***REMOVED***

	// 7 characters should NOT conflict
	assertIndexGet(t, index, id[:7], id, false)
	assertIndexGet(t, index, id2[:7], id2, false)

	// Deleting a non-existing id should return an error
	if err := index.Delete("non-existing"); err == nil ***REMOVED***
		t.Fatalf("Deleting a non-existing id should return an error")
	***REMOVED***

	// Deleting an empty id should return an error
	if err := index.Delete(""); err == nil ***REMOVED***
		t.Fatal("Deleting an empty id should return an error")
	***REMOVED***

	// Deleting id2 should remove conflicts
	if err := index.Delete(id2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// id2 should no longer work
	assertIndexGet(t, index, id2, "", true)
	assertIndexGet(t, index, id2[:7], "", true)
	assertIndexGet(t, index, id2[:11], "", true)

	// conflicts between id and id2 should be gone
	assertIndexGet(t, index, id[:6], id, false)
	assertIndexGet(t, index, id[:4], id, false)
	assertIndexGet(t, index, id[:1], id, false)

	// non-conflicting substrings should still not conflict
	assertIndexGet(t, index, id[:7], id, false)
	assertIndexGet(t, index, id[:15], id, false)
	assertIndexGet(t, index, id, id, false)

	assertIndexIterate(t)
	assertIndexIterateDoNotPanic(t)
***REMOVED***

func assertIndexIterate(t *testing.T) ***REMOVED***
	ids := []string***REMOVED***
		"19b36c2c326ccc11e726eee6ee78a0baf166ef96",
		"28b36c2c326ccc11e726eee6ee78a0baf166ef96",
		"37b36c2c326ccc11e726eee6ee78a0baf166ef96",
		"46b36c2c326ccc11e726eee6ee78a0baf166ef96",
	***REMOVED***

	index := NewTruncIndex(ids)

	index.Iterate(func(targetId string) ***REMOVED***
		for _, id := range ids ***REMOVED***
			if targetId == id ***REMOVED***
				return
			***REMOVED***
		***REMOVED***

		t.Fatalf("An unknown ID '%s'", targetId)
	***REMOVED***)
***REMOVED***

func assertIndexIterateDoNotPanic(t *testing.T) ***REMOVED***
	ids := []string***REMOVED***
		"19b36c2c326ccc11e726eee6ee78a0baf166ef96",
		"28b36c2c326ccc11e726eee6ee78a0baf166ef96",
	***REMOVED***

	index := NewTruncIndex(ids)
	iterationStarted := make(chan bool, 1)

	go func() ***REMOVED***
		<-iterationStarted
		index.Delete("19b36c2c326ccc11e726eee6ee78a0baf166ef96")
	***REMOVED***()

	index.Iterate(func(targetId string) ***REMOVED***
		if targetId == "19b36c2c326ccc11e726eee6ee78a0baf166ef96" ***REMOVED***
			iterationStarted <- true
			time.Sleep(100 * time.Millisecond)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func assertIndexGet(t *testing.T, index *TruncIndex, input, expectedResult string, expectError bool) ***REMOVED***
	if result, err := index.Get(input); err != nil && !expectError ***REMOVED***
		t.Fatalf("Unexpected error getting '%s': %s", input, err)
	***REMOVED*** else if err == nil && expectError ***REMOVED***
		t.Fatalf("Getting '%s' should return an error, not '%s'", input, result)
	***REMOVED*** else if result != expectedResult ***REMOVED***
		t.Fatalf("Getting '%s' returned '%s' instead of '%s'", input, result, expectedResult)
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexAdd100(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 100; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexAdd250(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 250; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexAdd500(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 500; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexGet100(b *testing.B) ***REMOVED***
	var testSet []string
	var testKeys []string
	for i := 0; i < 100; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	index := NewTruncIndex([]string***REMOVED******REMOVED***)
	for _, id := range testSet ***REMOVED***
		if err := index.Add(id); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, id := range testKeys ***REMOVED***
			if res, err := index.Get(id); err != nil ***REMOVED***
				b.Fatal(res, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexGet250(b *testing.B) ***REMOVED***
	var testSet []string
	var testKeys []string
	for i := 0; i < 250; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	index := NewTruncIndex([]string***REMOVED******REMOVED***)
	for _, id := range testSet ***REMOVED***
		if err := index.Add(id); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, id := range testKeys ***REMOVED***
			if res, err := index.Get(id); err != nil ***REMOVED***
				b.Fatal(res, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexGet500(b *testing.B) ***REMOVED***
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	index := NewTruncIndex([]string***REMOVED******REMOVED***)
	for _, id := range testSet ***REMOVED***
		if err := index.Add(id); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, id := range testKeys ***REMOVED***
			if res, err := index.Get(id); err != nil ***REMOVED***
				b.Fatal(res, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexDelete100(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 100; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		b.StopTimer()
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
		b.StartTimer()
		for _, id := range testSet ***REMOVED***
			if err := index.Delete(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexDelete250(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 250; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		b.StopTimer()
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
		b.StartTimer()
		for _, id := range testSet ***REMOVED***
			if err := index.Delete(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexDelete500(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 500; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		b.StopTimer()
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
		b.StartTimer()
		for _, id := range testSet ***REMOVED***
			if err := index.Delete(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexNew100(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 100; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		NewTruncIndex(testSet)
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexNew250(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 250; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		NewTruncIndex(testSet)
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexNew500(b *testing.B) ***REMOVED***
	var testSet []string
	for i := 0; i < 500; i++ ***REMOVED***
		testSet = append(testSet, stringid.GenerateNonCryptoID())
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		NewTruncIndex(testSet)
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexAddGet100(b *testing.B) ***REMOVED***
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ ***REMOVED***
		id := stringid.GenerateNonCryptoID()
		testSet = append(testSet, id)
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
		for _, id := range testKeys ***REMOVED***
			if res, err := index.Get(id); err != nil ***REMOVED***
				b.Fatal(res, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexAddGet250(b *testing.B) ***REMOVED***
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ ***REMOVED***
		id := stringid.GenerateNonCryptoID()
		testSet = append(testSet, id)
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
		for _, id := range testKeys ***REMOVED***
			if res, err := index.Get(id); err != nil ***REMOVED***
				b.Fatal(res, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkTruncIndexAddGet500(b *testing.B) ***REMOVED***
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ ***REMOVED***
		id := stringid.GenerateNonCryptoID()
		testSet = append(testSet, id)
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		index := NewTruncIndex([]string***REMOVED******REMOVED***)
		for _, id := range testSet ***REMOVED***
			if err := index.Add(id); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
		for _, id := range testKeys ***REMOVED***
			if res, err := index.Get(id); err != nil ***REMOVED***
				b.Fatal(res, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
