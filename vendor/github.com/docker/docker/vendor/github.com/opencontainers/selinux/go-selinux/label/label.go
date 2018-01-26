// +build !selinux !linux

package label

// InitLabels returns the process label and file labels to be used within
// the container.  A list of options can be passed into this function to alter
// the labels.
func InitLabels(options []string) (string, string, error) ***REMOVED***
	return "", "", nil
***REMOVED***

func GetROMountLabel() string ***REMOVED***
	return ""
***REMOVED***

func GenLabels(options string) (string, string, error) ***REMOVED***
	return "", "", nil
***REMOVED***

func FormatMountLabel(src string, mountLabel string) string ***REMOVED***
	return src
***REMOVED***

func SetProcessLabel(processLabel string) error ***REMOVED***
	return nil
***REMOVED***

func GetFileLabel(path string) (string, error) ***REMOVED***
	return "", nil
***REMOVED***

func SetFileLabel(path string, fileLabel string) error ***REMOVED***
	return nil
***REMOVED***

func SetFileCreateLabel(fileLabel string) error ***REMOVED***
	return nil
***REMOVED***

func Relabel(path string, fileLabel string, shared bool) error ***REMOVED***
	return nil
***REMOVED***

func GetPidLabel(pid int) (string, error) ***REMOVED***
	return "", nil
***REMOVED***

func Init() ***REMOVED***
***REMOVED***

func ReserveLabel(label string) error ***REMOVED***
	return nil
***REMOVED***

func ReleaseLabel(label string) error ***REMOVED***
	return nil
***REMOVED***

// DupSecOpt takes a process label and returns security options that
// can be used to set duplicate labels on future container processes
func DupSecOpt(src string) []string ***REMOVED***
	return nil
***REMOVED***

// DisableSecOpt returns a security opt that can disable labeling
// support for future container processes
func DisableSecOpt() []string ***REMOVED***
	return nil
***REMOVED***

// Validate checks that the label does not include unexpected options
func Validate(label string) error ***REMOVED***
	return nil
***REMOVED***

// RelabelNeeded checks whether the user requested a relabel
func RelabelNeeded(label string) bool ***REMOVED***
	return false
***REMOVED***

// IsShared checks that the label includes a "shared" mark
func IsShared(label string) bool ***REMOVED***
	return false
***REMOVED***
