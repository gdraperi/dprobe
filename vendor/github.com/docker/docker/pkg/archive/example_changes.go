// +build ignore

// Simple tool to create an archive stream from an old and new directory
//
// By default it will stream the comparison of two temporary directories with junk files
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
)

var (
	flDebug  = flag.Bool("D", false, "debugging output")
	flNewDir = flag.String("newdir", "", "")
	flOldDir = flag.String("olddir", "", "")
	log      = logrus.New()
)

func main() ***REMOVED***
	flag.Usage = func() ***REMOVED***
		fmt.Println("Produce a tar from comparing two directory paths. By default a demo tar is created of around 200 files (including hardlinks)")
		fmt.Printf("%s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	***REMOVED***
	flag.Parse()
	log.Out = os.Stderr
	if (len(os.Getenv("DEBUG")) > 0) || *flDebug ***REMOVED***
		logrus.SetLevel(logrus.DebugLevel)
	***REMOVED***
	var newDir, oldDir string

	if len(*flNewDir) == 0 ***REMOVED***
		var err error
		newDir, err = ioutil.TempDir("", "docker-test-newDir")
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		defer os.RemoveAll(newDir)
		if _, err := prepareUntarSourceDirectory(100, newDir, true); err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		newDir = *flNewDir
	***REMOVED***

	if len(*flOldDir) == 0 ***REMOVED***
		oldDir, err := ioutil.TempDir("", "docker-test-oldDir")
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		defer os.RemoveAll(oldDir)
	***REMOVED*** else ***REMOVED***
		oldDir = *flOldDir
	***REMOVED***

	changes, err := archive.ChangesDirs(newDir, oldDir)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	a, err := archive.ExportChanges(newDir, changes)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer a.Close()

	i, err := io.Copy(os.Stdout, a)
	if err != nil && err != io.EOF ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	fmt.Fprintf(os.Stderr, "wrote archive of %d bytes", i)
***REMOVED***

func prepareUntarSourceDirectory(numberOfFiles int, targetPath string, makeLinks bool) (int, error) ***REMOVED***
	fileData := []byte("fooo")
	for n := 0; n < numberOfFiles; n++ ***REMOVED***
		fileName := fmt.Sprintf("file-%d", n)
		if err := ioutil.WriteFile(path.Join(targetPath, fileName), fileData, 0700); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if makeLinks ***REMOVED***
			if err := os.Link(path.Join(targetPath, fileName), path.Join(targetPath, fileName+"-link")); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	totalSize := numberOfFiles * len(fileData)
	return totalSize, nil
***REMOVED***
