package system

import (
	"unsafe"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

var (
	ntuserApiset       = windows.NewLazyDLL("ext-ms-win-ntuser-window-l1-1-0")
	procGetVersionExW  = modkernel32.NewProc("GetVersionExW")
	procGetProductInfo = modkernel32.NewProc("GetProductInfo")
)

// OSVersion is a wrapper for Windows version information
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724439(v=vs.85).aspx
type OSVersion struct ***REMOVED***
	Version      uint32
	MajorVersion uint8
	MinorVersion uint8
	Build        uint16
***REMOVED***

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724833(v=vs.85).aspx
type osVersionInfoEx struct ***REMOVED***
	OSVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformID        uint32
	CSDVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	Reserve           byte
***REMOVED***

// GetOSVersion gets the operating system version on Windows. Note that
// docker.exe must be manifested to get the correct version information.
func GetOSVersion() OSVersion ***REMOVED***
	var err error
	osv := OSVersion***REMOVED******REMOVED***
	osv.Version, err = windows.GetVersion()
	if err != nil ***REMOVED***
		// GetVersion never fails.
		panic(err)
	***REMOVED***
	osv.MajorVersion = uint8(osv.Version & 0xFF)
	osv.MinorVersion = uint8(osv.Version >> 8 & 0xFF)
	osv.Build = uint16(osv.Version >> 16)
	return osv
***REMOVED***

// IsWindowsClient returns true if the SKU is client
// @engine maintainers - this function should not be removed or modified as it
// is used to enforce licensing restrictions on Windows.
func IsWindowsClient() bool ***REMOVED***
	osviex := &osVersionInfoEx***REMOVED***OSVersionInfoSize: 284***REMOVED***
	r1, _, err := procGetVersionExW.Call(uintptr(unsafe.Pointer(osviex)))
	if r1 == 0 ***REMOVED***
		logrus.Warnf("GetVersionExW failed - assuming server SKU: %v", err)
		return false
	***REMOVED***
	const verNTWorkstation = 0x00000001
	return osviex.ProductType == verNTWorkstation
***REMOVED***

// IsIoTCore returns true if the currently running image is based off of
// Windows 10 IoT Core.
// @engine maintainers - this function should not be removed or modified as it
// is used to enforce licensing restrictions on Windows.
func IsIoTCore() bool ***REMOVED***
	var returnedProductType uint32
	r1, _, err := procGetProductInfo.Call(6, 1, 0, 0, uintptr(unsafe.Pointer(&returnedProductType)))
	if r1 == 0 ***REMOVED***
		logrus.Warnf("GetProductInfo failed - assuming this is not IoT: %v", err)
		return false
	***REMOVED***
	const productIoTUAP = 0x0000007B
	const productIoTUAPCommercial = 0x00000083
	return returnedProductType == productIoTUAP || returnedProductType == productIoTUAPCommercial
***REMOVED***

// Unmount is a platform-specific helper function to call
// the unmount syscall. Not supported on Windows
func Unmount(dest string) error ***REMOVED***
	return nil
***REMOVED***

// CommandLineToArgv wraps the Windows syscall to turn a commandline into an argument array.
func CommandLineToArgv(commandLine string) ([]string, error) ***REMOVED***
	var argc int32

	argsPtr, err := windows.UTF16PtrFromString(commandLine)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	argv, err := windows.CommandLineToArgv(argsPtr, &argc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer windows.LocalFree(windows.Handle(uintptr(unsafe.Pointer(argv))))

	newArgs := make([]string, argc)
	for i, v := range (*argv)[:argc] ***REMOVED***
		newArgs[i] = string(windows.UTF16ToString((*v)[:]))
	***REMOVED***

	return newArgs, nil
***REMOVED***

// HasWin32KSupport determines whether containers that depend on win32k can
// run on this machine. Win32k is the driver used to implement windowing.
func HasWin32KSupport() bool ***REMOVED***
	// For now, check for ntuser API support on the host. In the future, a host
	// may support win32k in containers even if the host does not support ntuser
	// APIs.
	return ntuserApiset.Load() == nil
***REMOVED***
