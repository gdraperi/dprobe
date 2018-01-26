package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

func main() ***REMOVED***
	if err := xmain(); err != nil ***REMOVED***
		log.Fatalf("fatal error: %v", err)
	***REMOVED***
***REMOVED***

func xmain() error ***REMOVED***
	workerService := flag.String("worker-service", "", "Name of worker service")
	chunks := flag.Int("chunks", 0, "Number of chunks")
	input := flag.String("input", "", "Path to input file")
	randSeed := flag.Int64("rand-seed", int64(0), "Random seed")
	shuffle := flag.Bool("shuffle", false, "Shuffle the input so as to mitigate makespan nonuniformity")
	flag.Parse()
	if *workerService == "" ***REMOVED***
		return errors.New("worker-service unset")
	***REMOVED***
	if *chunks == 0 ***REMOVED***
		return errors.New("chunks unset")
	***REMOVED***
	if *input == "" ***REMOVED***
		return errors.New("input unset")
	***REMOVED***

	tests, err := loadTests(*input)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	testChunks := chunkTests(tests, *chunks, *shuffle, *randSeed)
	log.Printf("Loaded %d tests (%d chunks)", len(tests), len(testChunks))
	return executeTests(*workerService, testChunks)
***REMOVED***

func chunkTests(tests []string, numChunks int, shuffle bool, randSeed int64) [][]string ***REMOVED***
	// shuffling (experimental) mitigates makespan nonuniformity
	// Not sure this can cause some locality problem..
	if shuffle ***REMOVED***
		shuffleStrings(tests, randSeed)
	***REMOVED***
	return chunkStrings(tests, numChunks)
***REMOVED***

func loadTests(filename string) ([]string, error) ***REMOVED***
	b, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var tests []string
	for _, line := range strings.Split(string(b), "\n") ***REMOVED***
		s := strings.TrimSpace(line)
		if s != "" ***REMOVED***
			tests = append(tests, s)
		***REMOVED***
	***REMOVED***
	return tests, nil
***REMOVED***
