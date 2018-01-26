// +build selinux,linux

package label

import (
	"fmt"
	"strings"

	"github.com/opencontainers/selinux/go-selinux"
)

// Valid Label Options
var validOptions = map[string]bool***REMOVED***
	"disable": true,
	"type":    true,
	"user":    true,
	"role":    true,
	"level":   true,
***REMOVED***

var ErrIncompatibleLabel = fmt.Errorf("Bad SELinux option z and Z can not be used together")

// InitLabels returns the process label and file labels to be used within
// the container.  A list of options can be passed into this function to alter
// the labels.  The labels returned will include a random MCS String, that is
// guaranteed to be unique.
func InitLabels(options []string) (string, string, error) ***REMOVED***
	if !selinux.GetEnabled() ***REMOVED***
		return "", "", nil
	***REMOVED***
	processLabel, mountLabel := selinux.ContainerLabels()
	if processLabel != "" ***REMOVED***
		pcon := selinux.NewContext(processLabel)
		mcon := selinux.NewContext(mountLabel)
		for _, opt := range options ***REMOVED***
			if opt == "disable" ***REMOVED***
				return "", "", nil
			***REMOVED***
			if i := strings.Index(opt, ":"); i == -1 ***REMOVED***
				return "", "", fmt.Errorf("Bad label option %q, valid options 'disable' or \n'user, role, level, type' followed by ':' and a value", opt)
			***REMOVED***
			con := strings.SplitN(opt, ":", 2)
			if !validOptions[con[0]] ***REMOVED***
				return "", "", fmt.Errorf("Bad label option %q, valid options 'disable, user, role, level, type'", con[0])

			***REMOVED***
			pcon[con[0]] = con[1]
			if con[0] == "level" || con[0] == "user" ***REMOVED***
				mcon[con[0]] = con[1]
			***REMOVED***
		***REMOVED***
		_ = ReleaseLabel(processLabel)
		processLabel = pcon.Get()
		mountLabel = mcon.Get()
		_ = ReserveLabel(processLabel)
	***REMOVED***
	return processLabel, mountLabel, nil
***REMOVED***

func ROMountLabel() string ***REMOVED***
	return selinux.ROFileLabel()
***REMOVED***

// DEPRECATED: The GenLabels function is only to be used during the transition to the official API.
func GenLabels(options string) (string, string, error) ***REMOVED***
	return InitLabels(strings.Fields(options))
***REMOVED***

// FormatMountLabel returns a string to be used by the mount command.
// The format of this string will be used to alter the labeling of the mountpoint.
// The string returned is suitable to be used as the options field of the mount command.
// If you need to have additional mount point options, you can pass them in as
// the first parameter.  Second parameter is the label that you wish to apply
// to all content in the mount point.
func FormatMountLabel(src, mountLabel string) string ***REMOVED***
	if mountLabel != "" ***REMOVED***
		switch src ***REMOVED***
		case "":
			src = fmt.Sprintf("context=%q", mountLabel)
		default:
			src = fmt.Sprintf("%s,context=%q", src, mountLabel)
		***REMOVED***
	***REMOVED***
	return src
***REMOVED***

// SetProcessLabel takes a process label and tells the kernel to assign the
// label to the next program executed by the current process.
func SetProcessLabel(processLabel string) error ***REMOVED***
	if processLabel == "" ***REMOVED***
		return nil
	***REMOVED***
	return selinux.SetExecLabel(processLabel)
***REMOVED***

// ProcessLabel returns the process label that the kernel will assign
// to the next program executed by the current process.  If "" is returned
// this indicates that the default labeling will happen for the process.
func ProcessLabel() (string, error) ***REMOVED***
	return selinux.ExecLabel()
***REMOVED***

// GetFileLabel returns the label for specified path
func FileLabel(path string) (string, error) ***REMOVED***
	return selinux.FileLabel(path)
***REMOVED***

// SetFileLabel modifies the "path" label to the specified file label
func SetFileLabel(path string, fileLabel string) error ***REMOVED***
	if selinux.GetEnabled() && fileLabel != "" ***REMOVED***
		return selinux.SetFileLabel(path, fileLabel)
	***REMOVED***
	return nil
***REMOVED***

// SetFileCreateLabel tells the kernel the label for all files to be created
func SetFileCreateLabel(fileLabel string) error ***REMOVED***
	if selinux.GetEnabled() ***REMOVED***
		return selinux.SetFSCreateLabel(fileLabel)
	***REMOVED***
	return nil
***REMOVED***

// Relabel changes the label of path to the filelabel string.
// It changes the MCS label to s0 if shared is true.
// This will allow all containers to share the content.
func Relabel(path string, fileLabel string, shared bool) error ***REMOVED***
	if !selinux.GetEnabled() ***REMOVED***
		return nil
	***REMOVED***

	if fileLabel == "" ***REMOVED***
		return nil
	***REMOVED***

	exclude_paths := map[string]bool***REMOVED***"/": true, "/usr": true, "/etc": true***REMOVED***
	if exclude_paths[path] ***REMOVED***
		return fmt.Errorf("SELinux relabeling of %s is not allowed", path)
	***REMOVED***

	if shared ***REMOVED***
		c := selinux.NewContext(fileLabel)
		c["level"] = "s0"
		fileLabel = c.Get()
	***REMOVED***
	if err := selinux.Chcon(path, fileLabel, true); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// PidLabel will return the label of the process running with the specified pid
func PidLabel(pid int) (string, error) ***REMOVED***
	return selinux.PidLabel(pid)
***REMOVED***

// Init initialises the labeling system
func Init() ***REMOVED***
	selinux.GetEnabled()
***REMOVED***

// ReserveLabel will record the fact that the MCS label has already been used.
// This will prevent InitLabels from using the MCS label in a newly created
// container
func ReserveLabel(label string) error ***REMOVED***
	selinux.ReserveLabel(label)
	return nil
***REMOVED***

// ReleaseLabel will remove the reservation of the MCS label.
// This will allow InitLabels to use the MCS label in a newly created
// containers
func ReleaseLabel(label string) error ***REMOVED***
	selinux.ReleaseLabel(label)
	return nil
***REMOVED***

// DupSecOpt takes a process label and returns security options that
// can be used to set duplicate labels on future container processes
func DupSecOpt(src string) []string ***REMOVED***
	return selinux.DupSecOpt(src)
***REMOVED***

// DisableSecOpt returns a security opt that can disable labeling
// support for future container processes
func DisableSecOpt() []string ***REMOVED***
	return selinux.DisableSecOpt()
***REMOVED***

// Validate checks that the label does not include unexpected options
func Validate(label string) error ***REMOVED***
	if strings.Contains(label, "z") && strings.Contains(label, "Z") ***REMOVED***
		return ErrIncompatibleLabel
	***REMOVED***
	return nil
***REMOVED***

// RelabelNeeded checks whether the user requested a relabel
func RelabelNeeded(label string) bool ***REMOVED***
	return strings.Contains(label, "z") || strings.Contains(label, "Z")
***REMOVED***

// IsShared checks that the label includes a "shared" mark
func IsShared(label string) bool ***REMOVED***
	return strings.Contains(label, "z")
***REMOVED***
