// +build windows

package winio

import (
	"syscall"
	"unsafe"
)

//sys lookupAccountName(systemName *uint16, accountName string, sid *byte, sidSize *uint32, refDomain *uint16, refDomainSize *uint32, sidNameUse *uint32) (err error) = advapi32.LookupAccountNameW
//sys convertSidToStringSid(sid *byte, str **uint16) (err error) = advapi32.ConvertSidToStringSidW
//sys convertStringSecurityDescriptorToSecurityDescriptor(str string, revision uint32, sd *uintptr, size *uint32) (err error) = advapi32.ConvertStringSecurityDescriptorToSecurityDescriptorW
//sys convertSecurityDescriptorToStringSecurityDescriptor(sd *byte, revision uint32, secInfo uint32, sddl **uint16, sddlSize *uint32) (err error) = advapi32.ConvertSecurityDescriptorToStringSecurityDescriptorW
//sys localFree(mem uintptr) = LocalFree
//sys getSecurityDescriptorLength(sd uintptr) (len uint32) = advapi32.GetSecurityDescriptorLength

const (
	cERROR_NONE_MAPPED = syscall.Errno(1332)
)

type AccountLookupError struct ***REMOVED***
	Name string
	Err  error
***REMOVED***

func (e *AccountLookupError) Error() string ***REMOVED***
	if e.Name == "" ***REMOVED***
		return "lookup account: empty account name specified"
	***REMOVED***
	var s string
	switch e.Err ***REMOVED***
	case cERROR_NONE_MAPPED:
		s = "not found"
	default:
		s = e.Err.Error()
	***REMOVED***
	return "lookup account " + e.Name + ": " + s
***REMOVED***

type SddlConversionError struct ***REMOVED***
	Sddl string
	Err  error
***REMOVED***

func (e *SddlConversionError) Error() string ***REMOVED***
	return "convert " + e.Sddl + ": " + e.Err.Error()
***REMOVED***

// LookupSidByName looks up the SID of an account by name
func LookupSidByName(name string) (sid string, err error) ***REMOVED***
	if name == "" ***REMOVED***
		return "", &AccountLookupError***REMOVED***name, cERROR_NONE_MAPPED***REMOVED***
	***REMOVED***

	var sidSize, sidNameUse, refDomainSize uint32
	err = lookupAccountName(nil, name, nil, &sidSize, nil, &refDomainSize, &sidNameUse)
	if err != nil && err != syscall.ERROR_INSUFFICIENT_BUFFER ***REMOVED***
		return "", &AccountLookupError***REMOVED***name, err***REMOVED***
	***REMOVED***
	sidBuffer := make([]byte, sidSize)
	refDomainBuffer := make([]uint16, refDomainSize)
	err = lookupAccountName(nil, name, &sidBuffer[0], &sidSize, &refDomainBuffer[0], &refDomainSize, &sidNameUse)
	if err != nil ***REMOVED***
		return "", &AccountLookupError***REMOVED***name, err***REMOVED***
	***REMOVED***
	var strBuffer *uint16
	err = convertSidToStringSid(&sidBuffer[0], &strBuffer)
	if err != nil ***REMOVED***
		return "", &AccountLookupError***REMOVED***name, err***REMOVED***
	***REMOVED***
	sid = syscall.UTF16ToString((*[0xffff]uint16)(unsafe.Pointer(strBuffer))[:])
	localFree(uintptr(unsafe.Pointer(strBuffer)))
	return sid, nil
***REMOVED***

func SddlToSecurityDescriptor(sddl string) ([]byte, error) ***REMOVED***
	var sdBuffer uintptr
	err := convertStringSecurityDescriptorToSecurityDescriptor(sddl, 1, &sdBuffer, nil)
	if err != nil ***REMOVED***
		return nil, &SddlConversionError***REMOVED***sddl, err***REMOVED***
	***REMOVED***
	defer localFree(sdBuffer)
	sd := make([]byte, getSecurityDescriptorLength(sdBuffer))
	copy(sd, (*[0xffff]byte)(unsafe.Pointer(sdBuffer))[:len(sd)])
	return sd, nil
***REMOVED***

func SecurityDescriptorToSddl(sd []byte) (string, error) ***REMOVED***
	var sddl *uint16
	// The returned string length seems to including an aribtrary number of terminating NULs.
	// Don't use it.
	err := convertSecurityDescriptorToStringSecurityDescriptor(&sd[0], 1, 0xff, &sddl, nil)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer localFree(uintptr(unsafe.Pointer(sddl)))
	return syscall.UTF16ToString((*[0xffff]uint16)(unsafe.Pointer(sddl))[:]), nil
***REMOVED***
