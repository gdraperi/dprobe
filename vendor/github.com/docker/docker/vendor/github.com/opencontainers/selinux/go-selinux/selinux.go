// +build linux

package selinux

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	// Enforcing constant indicate SELinux is in enforcing mode
	Enforcing = 1
	// Permissive constant to indicate SELinux is in permissive mode
	Permissive = 0
	// Disabled constant to indicate SELinux is disabled
	Disabled         = -1
	selinuxDir       = "/etc/selinux/"
	selinuxConfig    = selinuxDir + "config"
	selinuxTypeTag   = "SELINUXTYPE"
	selinuxTag       = "SELINUX"
	selinuxPath      = "/sys/fs/selinux"
	xattrNameSelinux = "security.selinux"
	stRdOnly         = 0x01
)

type selinuxState struct ***REMOVED***
	enabledSet   bool
	enabled      bool
	selinuxfsSet bool
	selinuxfs    string
	mcsList      map[string]bool
	sync.Mutex
***REMOVED***

var (
	assignRegex = regexp.MustCompile(`^([^=]+)=(.*)$`)
	state       = selinuxState***REMOVED***
		mcsList: make(map[string]bool),
	***REMOVED***
)

// Context is a representation of the SELinux label broken into 4 parts
type Context map[string]string

func (s *selinuxState) setEnable(enabled bool) bool ***REMOVED***
	s.Lock()
	defer s.Unlock()
	s.enabledSet = true
	s.enabled = enabled
	return s.enabled
***REMOVED***

func (s *selinuxState) getEnabled() bool ***REMOVED***
	s.Lock()
	enabled := s.enabled
	enabledSet := s.enabledSet
	s.Unlock()
	if enabledSet ***REMOVED***
		return enabled
	***REMOVED***

	enabled = false
	if fs := getSelinuxMountPoint(); fs != "" ***REMOVED***
		if con, _ := CurrentLabel(); con != "kernel" ***REMOVED***
			enabled = true
		***REMOVED***
	***REMOVED***
	return s.setEnable(enabled)
***REMOVED***

// SetDisabled disables selinux support for the package
func SetDisabled() ***REMOVED***
	state.setEnable(false)
***REMOVED***

func (s *selinuxState) setSELinuxfs(selinuxfs string) string ***REMOVED***
	s.Lock()
	defer s.Unlock()
	s.selinuxfsSet = true
	s.selinuxfs = selinuxfs
	return s.selinuxfs
***REMOVED***

func (s *selinuxState) getSELinuxfs() string ***REMOVED***
	s.Lock()
	selinuxfs := s.selinuxfs
	selinuxfsSet := s.selinuxfsSet
	s.Unlock()
	if selinuxfsSet ***REMOVED***
		return selinuxfs
	***REMOVED***

	selinuxfs = ""
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil ***REMOVED***
		return selinuxfs
	***REMOVED***
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() ***REMOVED***
		txt := scanner.Text()
		// Safe as mountinfo encodes mountpoints with spaces as \040.
		sepIdx := strings.Index(txt, " - ")
		if sepIdx == -1 ***REMOVED***
			continue
		***REMOVED***
		if !strings.Contains(txt[sepIdx:], "selinuxfs") ***REMOVED***
			continue
		***REMOVED***
		fields := strings.Split(txt, " ")
		if len(fields) < 5 ***REMOVED***
			continue
		***REMOVED***
		selinuxfs = fields[4]
		break
	***REMOVED***

	if selinuxfs != "" ***REMOVED***
		var buf syscall.Statfs_t
		syscall.Statfs(selinuxfs, &buf)
		if (buf.Flags & stRdOnly) == 1 ***REMOVED***
			selinuxfs = ""
		***REMOVED***
	***REMOVED***
	return s.setSELinuxfs(selinuxfs)
***REMOVED***

// getSelinuxMountPoint returns the path to the mountpoint of an selinuxfs
// filesystem or an empty string if no mountpoint is found.  Selinuxfs is
// a proc-like pseudo-filesystem that exposes the selinux policy API to
// processes.  The existence of an selinuxfs mount is used to determine
// whether selinux is currently enabled or not.
func getSelinuxMountPoint() string ***REMOVED***
	return state.getSELinuxfs()
***REMOVED***

// GetEnabled returns whether selinux is currently enabled.
func GetEnabled() bool ***REMOVED***
	return state.getEnabled()
***REMOVED***

func readConfig(target string) (value string) ***REMOVED***
	var (
		val, key string
		bufin    *bufio.Reader
	)

	in, err := os.Open(selinuxConfig)
	if err != nil ***REMOVED***
		return ""
	***REMOVED***
	defer in.Close()

	bufin = bufio.NewReader(in)

	for done := false; !done; ***REMOVED***
		var line string
		if line, err = bufin.ReadString('\n'); err != nil ***REMOVED***
			if err != io.EOF ***REMOVED***
				return ""
			***REMOVED***
			done = true
		***REMOVED***
		line = strings.TrimSpace(line)
		if len(line) == 0 ***REMOVED***
			// Skip blank lines
			continue
		***REMOVED***
		if line[0] == ';' || line[0] == '#' ***REMOVED***
			// Skip comments
			continue
		***REMOVED***
		if groups := assignRegex.FindStringSubmatch(line); groups != nil ***REMOVED***
			key, val = strings.TrimSpace(groups[1]), strings.TrimSpace(groups[2])
			if key == target ***REMOVED***
				return strings.Trim(val, "\"")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func getSELinuxPolicyRoot() string ***REMOVED***
	return selinuxDir + readConfig(selinuxTypeTag)
***REMOVED***

func readCon(name string) (string, error) ***REMOVED***
	var val string

	in, err := os.Open(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer in.Close()

	_, err = fmt.Fscanf(in, "%s", &val)
	return val, err
***REMOVED***

// SetFileLabel sets the SELinux label for this path or returns an error.
func SetFileLabel(path string, label string) error ***REMOVED***
	return lsetxattr(path, xattrNameSelinux, []byte(label), 0)
***REMOVED***

// FileLabel returns the SELinux label for this path or returns an error.
func FileLabel(path string) (string, error) ***REMOVED***
	label, err := lgetxattr(path, xattrNameSelinux)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// Trim the NUL byte at the end of the byte buffer, if present.
	if len(label) > 0 && label[len(label)-1] == '\x00' ***REMOVED***
		label = label[:len(label)-1]
	***REMOVED***
	return string(label), nil
***REMOVED***

/*
SetFSCreateLabel tells kernel the label to create all file system objects
created by this task. Setting label="" to return to default.
*/
func SetFSCreateLabel(label string) error ***REMOVED***
	return writeCon(fmt.Sprintf("/proc/self/task/%d/attr/fscreate", syscall.Gettid()), label)
***REMOVED***

/*
FSCreateLabel returns the default label the kernel which the kernel is using
for file system objects created by this task. "" indicates default.
*/
func FSCreateLabel() (string, error) ***REMOVED***
	return readCon(fmt.Sprintf("/proc/self/task/%d/attr/fscreate", syscall.Gettid()))
***REMOVED***

// CurrentLabel returns the SELinux label of the current process thread, or an error.
func CurrentLabel() (string, error) ***REMOVED***
	return readCon(fmt.Sprintf("/proc/self/task/%d/attr/current", syscall.Gettid()))
***REMOVED***

// PidLabel returns the SELinux label of the given pid, or an error.
func PidLabel(pid int) (string, error) ***REMOVED***
	return readCon(fmt.Sprintf("/proc/%d/attr/current", pid))
***REMOVED***

/*
ExecLabel returns the SELinux label that the kernel will use for any programs
that are executed by the current process thread, or an error.
*/
func ExecLabel() (string, error) ***REMOVED***
	return readCon(fmt.Sprintf("/proc/self/task/%d/attr/exec", syscall.Gettid()))
***REMOVED***

func writeCon(name string, val string) error ***REMOVED***
	out, err := os.OpenFile(name, os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer out.Close()

	if val != "" ***REMOVED***
		_, err = out.Write([]byte(val))
	***REMOVED*** else ***REMOVED***
		_, err = out.Write(nil)
	***REMOVED***
	return err
***REMOVED***

/*
SetExecLabel sets the SELinux label that the kernel will use for any programs
that are executed by the current process thread, or an error.
*/
func SetExecLabel(label string) error ***REMOVED***
	return writeCon(fmt.Sprintf("/proc/self/task/%d/attr/exec", syscall.Gettid()), label)
***REMOVED***

// Get returns the Context as a string
func (c Context) Get() string ***REMOVED***
	return fmt.Sprintf("%s:%s:%s:%s", c["user"], c["role"], c["type"], c["level"])
***REMOVED***

// NewContext creates a new Context struct from the specified label
func NewContext(label string) Context ***REMOVED***
	c := make(Context)

	if len(label) != 0 ***REMOVED***
		con := strings.SplitN(label, ":", 4)
		c["user"] = con[0]
		c["role"] = con[1]
		c["type"] = con[2]
		c["level"] = con[3]
	***REMOVED***
	return c
***REMOVED***

// ReserveLabel reserves the MLS/MCS level component of the specified label
func ReserveLabel(label string) ***REMOVED***
	if len(label) != 0 ***REMOVED***
		con := strings.SplitN(label, ":", 4)
		mcsAdd(con[3])
	***REMOVED***
***REMOVED***

func selinuxEnforcePath() string ***REMOVED***
	return fmt.Sprintf("%s/enforce", selinuxPath)
***REMOVED***

// EnforceMode returns the current SELinux mode Enforcing, Permissive, Disabled
func EnforceMode() int ***REMOVED***
	var enforce int

	enforceS, err := readCon(selinuxEnforcePath())
	if err != nil ***REMOVED***
		return -1
	***REMOVED***

	enforce, err = strconv.Atoi(string(enforceS))
	if err != nil ***REMOVED***
		return -1
	***REMOVED***
	return enforce
***REMOVED***

/*
SetEnforceMode sets the current SELinux mode Enforcing, Permissive.
Disabled is not valid, since this needs to be set at boot time.
*/
func SetEnforceMode(mode int) error ***REMOVED***
	return writeCon(selinuxEnforcePath(), fmt.Sprintf("%d", mode))
***REMOVED***

/*
DefaultEnforceMode returns the systems default SELinux mode Enforcing,
Permissive or Disabled. Note this is is just the default at boot time.
EnforceMode tells you the systems current mode.
*/
func DefaultEnforceMode() int ***REMOVED***
	switch readConfig(selinuxTag) ***REMOVED***
	case "enforcing":
		return Enforcing
	case "permissive":
		return Permissive
	***REMOVED***
	return Disabled
***REMOVED***

func mcsAdd(mcs string) error ***REMOVED***
	state.Lock()
	defer state.Unlock()
	if state.mcsList[mcs] ***REMOVED***
		return fmt.Errorf("MCS Label already exists")
	***REMOVED***
	state.mcsList[mcs] = true
	return nil
***REMOVED***

func mcsDelete(mcs string) ***REMOVED***
	state.Lock()
	defer state.Unlock()
	state.mcsList[mcs] = false
***REMOVED***

func intToMcs(id int, catRange uint32) string ***REMOVED***
	var (
		SETSIZE = int(catRange)
		TIER    = SETSIZE
		ORD     = id
	)

	if id < 1 || id > 523776 ***REMOVED***
		return ""
	***REMOVED***

	for ORD > TIER ***REMOVED***
		ORD = ORD - TIER
		TIER--
	***REMOVED***
	TIER = SETSIZE - TIER
	ORD = ORD + TIER
	return fmt.Sprintf("s0:c%d,c%d", TIER, ORD)
***REMOVED***

func uniqMcs(catRange uint32) string ***REMOVED***
	var (
		n      uint32
		c1, c2 uint32
		mcs    string
	)

	for ***REMOVED***
		binary.Read(rand.Reader, binary.LittleEndian, &n)
		c1 = n % catRange
		binary.Read(rand.Reader, binary.LittleEndian, &n)
		c2 = n % catRange
		if c1 == c2 ***REMOVED***
			continue
		***REMOVED*** else ***REMOVED***
			if c1 > c2 ***REMOVED***
				c1, c2 = c2, c1
			***REMOVED***
		***REMOVED***
		mcs = fmt.Sprintf("s0:c%d,c%d", c1, c2)
		if err := mcsAdd(mcs); err != nil ***REMOVED***
			continue
		***REMOVED***
		break
	***REMOVED***
	return mcs
***REMOVED***

/*
ReleaseLabel will unreserve the MLS/MCS Level field of the specified label.
Allowing it to be used by another process.
*/
func ReleaseLabel(label string) ***REMOVED***
	if len(label) != 0 ***REMOVED***
		con := strings.SplitN(label, ":", 4)
		mcsDelete(con[3])
	***REMOVED***
***REMOVED***

var roFileLabel string

// ROFileLabel returns the specified SELinux readonly file label
func ROFileLabel() (fileLabel string) ***REMOVED***
	return roFileLabel
***REMOVED***

/*
ContainerLabels returns an allocated processLabel and fileLabel to be used for
container labeling by the calling process.
*/
func ContainerLabels() (processLabel string, fileLabel string) ***REMOVED***
	var (
		val, key string
		bufin    *bufio.Reader
	)

	if !GetEnabled() ***REMOVED***
		return "", ""
	***REMOVED***
	lxcPath := fmt.Sprintf("%s/contexts/lxc_contexts", getSELinuxPolicyRoot())
	in, err := os.Open(lxcPath)
	if err != nil ***REMOVED***
		return "", ""
	***REMOVED***
	defer in.Close()

	bufin = bufio.NewReader(in)

	for done := false; !done; ***REMOVED***
		var line string
		if line, err = bufin.ReadString('\n'); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				done = true
			***REMOVED*** else ***REMOVED***
				goto exit
			***REMOVED***
		***REMOVED***
		line = strings.TrimSpace(line)
		if len(line) == 0 ***REMOVED***
			// Skip blank lines
			continue
		***REMOVED***
		if line[0] == ';' || line[0] == '#' ***REMOVED***
			// Skip comments
			continue
		***REMOVED***
		if groups := assignRegex.FindStringSubmatch(line); groups != nil ***REMOVED***
			key, val = strings.TrimSpace(groups[1]), strings.TrimSpace(groups[2])
			if key == "process" ***REMOVED***
				processLabel = strings.Trim(val, "\"")
			***REMOVED***
			if key == "file" ***REMOVED***
				fileLabel = strings.Trim(val, "\"")
			***REMOVED***
			if key == "ro_file" ***REMOVED***
				roFileLabel = strings.Trim(val, "\"")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if processLabel == "" || fileLabel == "" ***REMOVED***
		return "", ""
	***REMOVED***

	if roFileLabel == "" ***REMOVED***
		roFileLabel = fileLabel
	***REMOVED***
exit:
	mcs := uniqMcs(1024)
	scon := NewContext(processLabel)
	scon["level"] = mcs
	processLabel = scon.Get()
	scon = NewContext(fileLabel)
	scon["level"] = mcs
	fileLabel = scon.Get()
	return processLabel, fileLabel
***REMOVED***

// SecurityCheckContext validates that the SELinux label is understood by the kernel
func SecurityCheckContext(val string) error ***REMOVED***
	return writeCon(fmt.Sprintf("%s.context", selinuxPath), val)
***REMOVED***

/*
CopyLevel returns a label with the MLS/MCS level from src label replaces on
the dest label.
*/
func CopyLevel(src, dest string) (string, error) ***REMOVED***
	if src == "" ***REMOVED***
		return "", nil
	***REMOVED***
	if err := SecurityCheckContext(src); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if err := SecurityCheckContext(dest); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	scon := NewContext(src)
	tcon := NewContext(dest)
	mcsDelete(tcon["level"])
	mcsAdd(scon["level"])
	tcon["level"] = scon["level"]
	return tcon.Get(), nil
***REMOVED***

// Prevent users from relabing system files
func badPrefix(fpath string) error ***REMOVED***
	var badprefixes = []string***REMOVED***"/usr"***REMOVED***

	for _, prefix := range badprefixes ***REMOVED***
		if fpath == prefix || strings.HasPrefix(fpath, fmt.Sprintf("%s/", prefix)) ***REMOVED***
			return fmt.Errorf("relabeling content in %s is not allowed", prefix)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Chcon changes the fpath file object to the SELinux label label.
// If the fpath is a directory and recurse is true Chcon will walk the
// directory tree setting the label
func Chcon(fpath string, label string, recurse bool) error ***REMOVED***
	if label == "" ***REMOVED***
		return nil
	***REMOVED***
	if err := badPrefix(fpath); err != nil ***REMOVED***
		return err
	***REMOVED***
	callback := func(p string, info os.FileInfo, err error) error ***REMOVED***
		return SetFileLabel(p, label)
	***REMOVED***

	if recurse ***REMOVED***
		return filepath.Walk(fpath, callback)
	***REMOVED***

	return SetFileLabel(fpath, label)
***REMOVED***

// DupSecOpt takes an SELinux process label and returns security options that
// can will set the SELinux Type and Level for future container processes
func DupSecOpt(src string) []string ***REMOVED***
	if src == "" ***REMOVED***
		return nil
	***REMOVED***
	con := NewContext(src)
	if con["user"] == "" ||
		con["role"] == "" ||
		con["type"] == "" ||
		con["level"] == "" ***REMOVED***
		return nil
	***REMOVED***
	return []string***REMOVED***"user:" + con["user"],
		"role:" + con["role"],
		"type:" + con["type"],
		"level:" + con["level"]***REMOVED***
***REMOVED***

// DisableSecOpt returns a security opt that can be used to disabling SELinux
// labeling support for future container processes
func DisableSecOpt() []string ***REMOVED***
	return []string***REMOVED***"disable"***REMOVED***
***REMOVED***
