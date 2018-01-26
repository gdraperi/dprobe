package system

// lcowSupported determines if Linux Containers on Windows are supported.
var lcowSupported = false

// InitLCOW sets whether LCOW is supported or not
func InitLCOW(experimental bool) ***REMOVED***
	v := GetOSVersion()
	if experimental && v.Build >= 16299 ***REMOVED***
		lcowSupported = true
	***REMOVED***
***REMOVED***
