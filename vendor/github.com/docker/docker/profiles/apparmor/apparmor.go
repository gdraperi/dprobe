// +build linux

package apparmor

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/docker/docker/pkg/aaparser"
)

var (
	// profileDirectory is the file store for apparmor profiles and macros.
	profileDirectory = "/etc/apparmor.d"
)

// profileData holds information about the given profile for generation.
type profileData struct ***REMOVED***
	// Name is profile name.
	Name string
	// Imports defines the apparmor functions to import, before defining the profile.
	Imports []string
	// InnerImports defines the apparmor functions to import in the profile.
	InnerImports []string
	// Version is the ***REMOVED***major, minor, patch***REMOVED*** version of apparmor_parser as a single number.
	Version int
***REMOVED***

// generateDefault creates an apparmor profile from ProfileData.
func (p *profileData) generateDefault(out io.Writer) error ***REMOVED***
	compiled, err := template.New("apparmor_profile").Parse(baseTemplate)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if macroExists("tunables/global") ***REMOVED***
		p.Imports = append(p.Imports, "#include <tunables/global>")
	***REMOVED*** else ***REMOVED***
		p.Imports = append(p.Imports, "@***REMOVED***PROC***REMOVED***=/proc/")
	***REMOVED***

	if macroExists("abstractions/base") ***REMOVED***
		p.InnerImports = append(p.InnerImports, "#include <abstractions/base>")
	***REMOVED***

	ver, err := aaparser.GetVersion()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.Version = ver

	return compiled.Execute(out, p)
***REMOVED***

// macrosExists checks if the passed macro exists.
func macroExists(m string) bool ***REMOVED***
	_, err := os.Stat(path.Join(profileDirectory, m))
	return err == nil
***REMOVED***

// InstallDefault generates a default profile in a temp directory determined by
// os.TempDir(), then loads the profile into the kernel using 'apparmor_parser'.
func InstallDefault(name string) error ***REMOVED***
	p := profileData***REMOVED***
		Name: name,
	***REMOVED***

	// Install to a temporary directory.
	f, err := ioutil.TempFile("", name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	profilePath := f.Name()

	defer f.Close()
	defer os.Remove(profilePath)

	if err := p.generateDefault(f); err != nil ***REMOVED***
		return err
	***REMOVED***

	return aaparser.LoadProfile(profilePath)
***REMOVED***

// IsLoaded checks if a profile with the given name has been loaded into the
// kernel.
func IsLoaded(name string) (bool, error) ***REMOVED***
	file, err := os.Open("/sys/kernel/security/apparmor/profiles")
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer file.Close()

	r := bufio.NewReader(file)
	for ***REMOVED***
		p, err := r.ReadString('\n')
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		if strings.HasPrefix(p, name+" ") ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***

	return false, nil
***REMOVED***
