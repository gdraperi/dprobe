# Generating Man Pages For Your Own cobra.Command

Generating man pages from a cobra command is incredibly easy. An example is as follows:

```go
package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() ***REMOVED***
	cmd := &cobra.Command***REMOVED***
		Use:   "test",
		Short: "my test program",
	***REMOVED***
	header := &doc.GenManHeader***REMOVED***
		Title: "MINE",
		Section: "3",
	***REMOVED***
	err := doc.GenManTree(cmd, header, "/tmp")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
***REMOVED***
```

That will get you a man page `/tmp/test.3`
