package cmd

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
)

var update = flag.Bool("update", false, "update .golden files")

func init() ***REMOVED***
	// Mute commands.
	addCmd.SetOutput(new(bytes.Buffer))
	initCmd.SetOutput(new(bytes.Buffer))
***REMOVED***

// compareFiles compares the content of files with pathA and pathB.
// If contents are equal, it returns nil.
// If not, it returns which files are not equal
// and diff (if system has diff command) between these files.
func compareFiles(pathA, pathB string) error ***REMOVED***
	contentA, err := ioutil.ReadFile(pathA)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	contentB, err := ioutil.ReadFile(pathB)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !bytes.Equal(contentA, contentB) ***REMOVED***
		output := new(bytes.Buffer)
		output.WriteString(fmt.Sprintf("%q and %q are not equal!\n\n", pathA, pathB))

		diffPath, err := exec.LookPath("diff")
		if err != nil ***REMOVED***
			// Don't execute diff if it can't be found.
			return nil
		***REMOVED***
		diffCmd := exec.Command(diffPath, "-u", pathA, pathB)
		diffCmd.Stdout = output
		diffCmd.Stderr = output

		output.WriteString("$ diff -u " + pathA + " " + pathB + "\n")
		if err := diffCmd.Run(); err != nil ***REMOVED***
			output.WriteString("\n" + err.Error())
		***REMOVED***
		return errors.New(output.String())
	***REMOVED***
	return nil
***REMOVED***

// checkLackFiles checks if all elements of expected are in got.
func checkLackFiles(expected, got []string) error ***REMOVED***
	lacks := make([]string, 0, len(expected))
	for _, ev := range expected ***REMOVED***
		if !stringInStringSlice(ev, got) ***REMOVED***
			lacks = append(lacks, ev)
		***REMOVED***
	***REMOVED***
	if len(lacks) > 0 ***REMOVED***
		return fmt.Errorf("Lack %v file(s): %v", len(lacks), lacks)
	***REMOVED***
	return nil
***REMOVED***

// stringInStringSlice checks if s is an element of slice.
func stringInStringSlice(s string, slice []string) bool ***REMOVED***
	for _, v := range slice ***REMOVED***
		if s == v ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
