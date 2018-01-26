package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/bfirsh/funker-go"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/hack/integration-cli-on-swarm/agent/types"
)

func main() ***REMOVED***
	if err := xmain(); err != nil ***REMOVED***
		log.Fatalf("fatal error: %v", err)
	***REMOVED***
***REMOVED***

func validImageDigest(s string) bool ***REMOVED***
	return reference.DigestRegexp.FindString(s) != ""
***REMOVED***

func xmain() error ***REMOVED***
	workerImageDigest := flag.String("worker-image-digest", "", "Needs to be the digest of this worker image itself")
	dryRun := flag.Bool("dry-run", false, "Dry run")
	keepExecutor := flag.Bool("keep-executor", false, "Do not auto-remove executor containers, which is used for running privileged programs on Swarm")
	flag.Parse()
	if !validImageDigest(*workerImageDigest) ***REMOVED***
		// Because of issue #29582.
		// `docker service create localregistry.example.com/blahblah:latest` pulls the image data to local, but not a tag.
		// So, `docker run localregistry.example.com/blahblah:latest` fails: `Unable to find image 'localregistry.example.com/blahblah:latest' locally`
		return fmt.Errorf("worker-image-digest must be a digest, got %q", *workerImageDigest)
	***REMOVED***
	executor := privilegedTestChunkExecutor(!*keepExecutor)
	if *dryRun ***REMOVED***
		executor = dryTestChunkExecutor()
	***REMOVED***
	return handle(*workerImageDigest, executor)
***REMOVED***

func handle(workerImageDigest string, executor testChunkExecutor) error ***REMOVED***
	log.Printf("Waiting for a funker request")
	return funker.Handle(func(args *types.Args) types.Result ***REMOVED***
		log.Printf("Executing chunk %d, contains %d test filters",
			args.ChunkID, len(args.Tests))
		begin := time.Now()
		code, rawLog, err := executor(workerImageDigest, args.Tests)
		if err != nil ***REMOVED***
			log.Printf("Error while executing chunk %d: %v", args.ChunkID, err)
			if code == 0 ***REMOVED***
				// Make sure this is a failure
				code = 1
			***REMOVED***
			return types.Result***REMOVED***
				ChunkID: args.ChunkID,
				Code:    int(code),
				RawLog:  rawLog,
			***REMOVED***
		***REMOVED***
		elapsed := time.Since(begin)
		log.Printf("Finished chunk %d, code=%d, elapsed=%v", args.ChunkID, code, elapsed)
		return types.Result***REMOVED***
			ChunkID: args.ChunkID,
			Code:    int(code),
			RawLog:  rawLog,
		***REMOVED***
	***REMOVED***)
***REMOVED***
