package doc_test

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func ExampleGenManTree() ***REMOVED***
	cmd := &cobra.Command***REMOVED***
		Use:   "test",
		Short: "my test program",
	***REMOVED***
	header := &doc.GenManHeader***REMOVED***
		Title:   "MINE",
		Section: "3",
	***REMOVED***
	doc.GenManTree(cmd, header, "/tmp")
***REMOVED***

func ExampleGenMan() ***REMOVED***
	cmd := &cobra.Command***REMOVED***
		Use:   "test",
		Short: "my test program",
	***REMOVED***
	header := &doc.GenManHeader***REMOVED***
		Title:   "MINE",
		Section: "3",
	***REMOVED***
	out := new(bytes.Buffer)
	doc.GenMan(cmd, header, out)
	fmt.Print(out.String())
***REMOVED***
