package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bfirsh/funker-go"
	"github.com/docker/docker/hack/integration-cli-on-swarm/agent/types"
)

const (
	// funkerRetryTimeout is for the issue https://github.com/bfirsh/funker/issues/3
	// When all the funker replicas are busy in their own job, we cannot connect to funker.
	funkerRetryTimeout  = 1 * time.Hour
	funkerRetryDuration = 1 * time.Second
)

// ticker is needed for some CI (e.g., on Travis, job is aborted when no output emitted for 10 minutes)
func ticker(d time.Duration) chan struct***REMOVED******REMOVED*** ***REMOVED***
	t := time.NewTicker(d)
	stop := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-t.C:
				log.Printf("tick (just for keeping CI job active) per %s", d.String())
			case <-stop:
				t.Stop()
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return stop
***REMOVED***

func executeTests(funkerName string, testChunks [][]string) error ***REMOVED***
	tickerStopper := ticker(9*time.Minute + 55*time.Second)
	defer func() ***REMOVED***
		close(tickerStopper)
	***REMOVED***()
	begin := time.Now()
	log.Printf("Executing %d chunks in parallel, against %q", len(testChunks), funkerName)
	var wg sync.WaitGroup
	var passed, failed uint32
	for chunkID, tests := range testChunks ***REMOVED***
		log.Printf("Executing chunk %d (contains %d test filters)", chunkID, len(tests))
		wg.Add(1)
		go func(chunkID int, tests []string) ***REMOVED***
			defer wg.Done()
			chunkBegin := time.Now()
			result, err := executeTestChunkWithRetry(funkerName, types.Args***REMOVED***
				ChunkID: chunkID,
				Tests:   tests,
			***REMOVED***)
			if result.RawLog != "" ***REMOVED***
				for _, s := range strings.Split(result.RawLog, "\n") ***REMOVED***
					log.Printf("Log (chunk %d): %s", chunkID, s)
				***REMOVED***
			***REMOVED***
			if err != nil ***REMOVED***
				log.Printf("Error while executing chunk %d: %v",
					chunkID, err)
				atomic.AddUint32(&failed, 1)
			***REMOVED*** else ***REMOVED***
				if result.Code == 0 ***REMOVED***
					atomic.AddUint32(&passed, 1)
				***REMOVED*** else ***REMOVED***
					atomic.AddUint32(&failed, 1)
				***REMOVED***
				log.Printf("Finished chunk %d [%d/%d] with %d test filters in %s, code=%d.",
					chunkID, passed+failed, len(testChunks), len(tests),
					time.Since(chunkBegin), result.Code)
			***REMOVED***
		***REMOVED***(chunkID, tests)
	***REMOVED***
	wg.Wait()
	// TODO: print actual tests rather than chunks
	log.Printf("Executed %d chunks in %s. PASS: %d, FAIL: %d.",
		len(testChunks), time.Since(begin), passed, failed)
	if failed > 0 ***REMOVED***
		return fmt.Errorf("%d chunks failed", failed)
	***REMOVED***
	return nil
***REMOVED***

func executeTestChunk(funkerName string, args types.Args) (types.Result, error) ***REMOVED***
	ret, err := funker.Call(funkerName, args)
	if err != nil ***REMOVED***
		return types.Result***REMOVED******REMOVED***, err
	***REMOVED***
	tmp, err := json.Marshal(ret)
	if err != nil ***REMOVED***
		return types.Result***REMOVED******REMOVED***, err
	***REMOVED***
	var result types.Result
	err = json.Unmarshal(tmp, &result)
	return result, err
***REMOVED***

func executeTestChunkWithRetry(funkerName string, args types.Args) (types.Result, error) ***REMOVED***
	begin := time.Now()
	for i := 0; time.Since(begin) < funkerRetryTimeout; i++ ***REMOVED***
		result, err := executeTestChunk(funkerName, args)
		if err == nil ***REMOVED***
			log.Printf("executeTestChunk(%q, %d) returned code %d in trial %d", funkerName, args.ChunkID, result.Code, i)
			return result, nil
		***REMOVED***
		if errorSeemsInteresting(err) ***REMOVED***
			log.Printf("Error while calling executeTestChunk(%q, %d), will retry (trial %d): %v",
				funkerName, args.ChunkID, i, err)
		***REMOVED***
		// TODO: non-constant sleep
		time.Sleep(funkerRetryDuration)
	***REMOVED***
	return types.Result***REMOVED******REMOVED***, fmt.Errorf("could not call executeTestChunk(%q, %d) in %v", funkerName, args.ChunkID, funkerRetryTimeout)
***REMOVED***

//  errorSeemsInteresting returns true if err does not seem about https://github.com/bfirsh/funker/issues/3
func errorSeemsInteresting(err error) bool ***REMOVED***
	boringSubstrs := []string***REMOVED***"connection refused", "connection reset by peer", "no such host", "transport endpoint is not connected", "no route to host"***REMOVED***
	errS := err.Error()
	for _, boringS := range boringSubstrs ***REMOVED***
		if strings.Contains(errS, boringS) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
