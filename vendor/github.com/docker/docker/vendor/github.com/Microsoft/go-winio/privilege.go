// +build windows

package winio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

//sys adjustTokenPrivileges(token windows.Token, releaseAll bool, input *byte, outputSize uint32, output *byte, requiredSize *uint32) (success bool, err error) [true] = advapi32.AdjustTokenPrivileges
//sys impersonateSelf(level uint32) (err error) = advapi32.ImpersonateSelf
//sys revertToSelf() (err error) = advapi32.RevertToSelf
//sys openThreadToken(thread syscall.Handle, accessMask uint32, openAsSelf bool, token *windows.Token) (err error) = advapi32.OpenThreadToken
//sys getCurrentThread() (h syscall.Handle) = GetCurrentThread
//sys lookupPrivilegeValue(systemName string, name string, luid *uint64) (err error) = advapi32.LookupPrivilegeValueW
//sys lookupPrivilegeName(systemName string, luid *uint64, buffer *uint16, size *uint32) (err error) = advapi32.LookupPrivilegeNameW
//sys lookupPrivilegeDisplayName(systemName string, name *uint16, buffer *uint16, size *uint32, languageId *uint32) (err error) = advapi32.LookupPrivilegeDisplayNameW

const (
	SE_PRIVILEGE_ENABLED = 2

	ERROR_NOT_ALL_ASSIGNED syscall.Errno = 1300

	SeBackupPrivilege  = "SeBackupPrivilege"
	SeRestorePrivilege = "SeRestorePrivilege"
)

const (
	securityAnonymous = iota
	securityIdentification
	securityImpersonation
	securityDelegation
)

var (
	privNames     = make(map[string]uint64)
	privNameMutex sync.Mutex
)

// PrivilegeError represents an error enabling privileges.
type PrivilegeError struct ***REMOVED***
	privileges []uint64
***REMOVED***

func (e *PrivilegeError) Error() string ***REMOVED***
	s := ""
	if len(e.privileges) > 1 ***REMOVED***
		s = "Could not enable privileges "
	***REMOVED*** else ***REMOVED***
		s = "Could not enable privilege "
	***REMOVED***
	for i, p := range e.privileges ***REMOVED***
		if i != 0 ***REMOVED***
			s += ", "
		***REMOVED***
		s += `"`
		s += getPrivilegeName(p)
		s += `"`
	***REMOVED***
	return s
***REMOVED***

// RunWithPrivilege enables a single privilege for a function call.
func RunWithPrivilege(name string, fn func() error) error ***REMOVED***
	return RunWithPrivileges([]string***REMOVED***name***REMOVED***, fn)
***REMOVED***

// RunWithPrivileges enables privileges for a function call.
func RunWithPrivileges(names []string, fn func() error) error ***REMOVED***
	privileges, err := mapPrivileges(names)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	token, err := newThreadToken()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer releaseThreadToken(token)
	err = adjustPrivileges(token, privileges, SE_PRIVILEGE_ENABLED)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return fn()
***REMOVED***

func mapPrivileges(names []string) ([]uint64, error) ***REMOVED***
	var privileges []uint64
	privNameMutex.Lock()
	defer privNameMutex.Unlock()
	for _, name := range names ***REMOVED***
		p, ok := privNames[name]
		if !ok ***REMOVED***
			err := lookupPrivilegeValue("", name, &p)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			privNames[name] = p
		***REMOVED***
		privileges = append(privileges, p)
	***REMOVED***
	return privileges, nil
***REMOVED***

// EnableProcessPrivileges enables privileges globally for the process.
func EnableProcessPrivileges(names []string) error ***REMOVED***
	return enableDisableProcessPrivilege(names, SE_PRIVILEGE_ENABLED)
***REMOVED***

// DisableProcessPrivileges disables privileges globally for the process.
func DisableProcessPrivileges(names []string) error ***REMOVED***
	return enableDisableProcessPrivilege(names, 0)
***REMOVED***

func enableDisableProcessPrivilege(names []string, action uint32) error ***REMOVED***
	privileges, err := mapPrivileges(names)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	p, _ := windows.GetCurrentProcess()
	var token windows.Token
	err = windows.OpenProcessToken(p, windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY, &token)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	defer token.Close()
	return adjustPrivileges(token, privileges, action)
***REMOVED***

func adjustPrivileges(token windows.Token, privileges []uint64, action uint32) error ***REMOVED***
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, uint32(len(privileges)))
	for _, p := range privileges ***REMOVED***
		binary.Write(&b, binary.LittleEndian, p)
		binary.Write(&b, binary.LittleEndian, action)
	***REMOVED***
	prevState := make([]byte, b.Len())
	reqSize := uint32(0)
	success, err := adjustTokenPrivileges(token, false, &b.Bytes()[0], uint32(len(prevState)), &prevState[0], &reqSize)
	if !success ***REMOVED***
		return err
	***REMOVED***
	if err == ERROR_NOT_ALL_ASSIGNED ***REMOVED***
		return &PrivilegeError***REMOVED***privileges***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func getPrivilegeName(luid uint64) string ***REMOVED***
	var nameBuffer [256]uint16
	bufSize := uint32(len(nameBuffer))
	err := lookupPrivilegeName("", &luid, &nameBuffer[0], &bufSize)
	if err != nil ***REMOVED***
		return fmt.Sprintf("<unknown privilege %d>", luid)
	***REMOVED***

	var displayNameBuffer [256]uint16
	displayBufSize := uint32(len(displayNameBuffer))
	var langID uint32
	err = lookupPrivilegeDisplayName("", &nameBuffer[0], &displayNameBuffer[0], &displayBufSize, &langID)
	if err != nil ***REMOVED***
		return fmt.Sprintf("<unknown privilege %s>", string(utf16.Decode(nameBuffer[:bufSize])))
	***REMOVED***

	return string(utf16.Decode(displayNameBuffer[:displayBufSize]))
***REMOVED***

func newThreadToken() (windows.Token, error) ***REMOVED***
	err := impersonateSelf(securityImpersonation)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	var token windows.Token
	err = openThreadToken(getCurrentThread(), syscall.TOKEN_ADJUST_PRIVILEGES|syscall.TOKEN_QUERY, false, &token)
	if err != nil ***REMOVED***
		rerr := revertToSelf()
		if rerr != nil ***REMOVED***
			panic(rerr)
		***REMOVED***
		return 0, err
	***REMOVED***
	return token, nil
***REMOVED***

func releaseThreadToken(h windows.Token) ***REMOVED***
	err := revertToSelf()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	h.Close()
***REMOVED***
