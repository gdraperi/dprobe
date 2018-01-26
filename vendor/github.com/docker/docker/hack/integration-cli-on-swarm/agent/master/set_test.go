package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func generateInput(inputLen int) []string ***REMOVED***
	input := []string***REMOVED******REMOVED***
	for i := 0; i < inputLen; i++ ***REMOVED***
		input = append(input, fmt.Sprintf("s%d", i))
	***REMOVED***

	return input
***REMOVED***

func testChunkStrings(t *testing.T, inputLen, numChunks int) ***REMOVED***
	t.Logf("inputLen=%d, numChunks=%d", inputLen, numChunks)
	input := generateInput(inputLen)
	result := chunkStrings(input, numChunks)
	t.Logf("result has %d chunks", len(result))
	inputReconstructedFromResult := []string***REMOVED******REMOVED***
	for i, chunk := range result ***REMOVED***
		t.Logf("chunk %d has %d elements", i, len(chunk))
		inputReconstructedFromResult = append(inputReconstructedFromResult, chunk...)
	***REMOVED***
	if !reflect.DeepEqual(input, inputReconstructedFromResult) ***REMOVED***
		t.Fatal("input != inputReconstructedFromResult")
	***REMOVED***
***REMOVED***

func TestChunkStrings_4_4(t *testing.T) ***REMOVED***
	testChunkStrings(t, 4, 4)
***REMOVED***

func TestChunkStrings_4_1(t *testing.T) ***REMOVED***
	testChunkStrings(t, 4, 1)
***REMOVED***

func TestChunkStrings_1_4(t *testing.T) ***REMOVED***
	testChunkStrings(t, 1, 4)
***REMOVED***

func TestChunkStrings_1000_8(t *testing.T) ***REMOVED***
	testChunkStrings(t, 1000, 8)
***REMOVED***

func TestChunkStrings_1000_9(t *testing.T) ***REMOVED***
	testChunkStrings(t, 1000, 9)
***REMOVED***

func testShuffleStrings(t *testing.T, inputLen int, seed int64) ***REMOVED***
	t.Logf("inputLen=%d, seed=%d", inputLen, seed)
	x := generateInput(inputLen)
	shuffleStrings(x, seed)
	t.Logf("shuffled: %v", x)
***REMOVED***

func TestShuffleStrings_100(t *testing.T) ***REMOVED***
	testShuffleStrings(t, 100, time.Now().UnixNano())
***REMOVED***
