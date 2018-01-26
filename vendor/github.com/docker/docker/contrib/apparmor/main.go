package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/docker/docker/pkg/aaparser"
)

type profileData struct ***REMOVED***
	Version int
***REMOVED***

func main() ***REMOVED***
	if len(os.Args) < 2 ***REMOVED***
		log.Fatal("pass a filename to save the profile in.")
	***REMOVED***

	// parse the arg
	apparmorProfilePath := os.Args[1]

	version, err := aaparser.GetVersion()
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	data := profileData***REMOVED***
		Version: version,
	***REMOVED***
	fmt.Printf("apparmor_parser is of version %+v\n", data)

	// parse the template
	compiled, err := template.New("apparmor_profile").Parse(dockerProfileTemplate)
	if err != nil ***REMOVED***
		log.Fatalf("parsing template failed: %v", err)
	***REMOVED***

	// make sure /etc/apparmor.d exists
	if err := os.MkdirAll(path.Dir(apparmorProfilePath), 0755); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	f, err := os.OpenFile(apparmorProfilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer f.Close()

	if err := compiled.Execute(f, data); err != nil ***REMOVED***
		log.Fatalf("executing template failed: %v", err)
	***REMOVED***

	fmt.Printf("created apparmor profile for version %+v at %q\n", data, apparmorProfilePath)
***REMOVED***
