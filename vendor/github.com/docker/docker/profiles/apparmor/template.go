// +build linux

package apparmor

// baseTemplate defines the default apparmor profile for containers.
const baseTemplate = `
***REMOVED******REMOVED***range $value := .Imports***REMOVED******REMOVED***
***REMOVED******REMOVED***$value***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED***

profile ***REMOVED******REMOVED***.Name***REMOVED******REMOVED*** flags=(attach_disconnected,mediate_deleted) ***REMOVED***
***REMOVED******REMOVED***range $value := .InnerImports***REMOVED******REMOVED***
  ***REMOVED******REMOVED***$value***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED***

  network,
  capability,
  file,
  umount,

  deny @***REMOVED***PROC***REMOVED***/* w,   # deny write for all files directly in /proc (not in a subdir)
  # deny write to files not in /proc/<number>/** or /proc/sys/**
  deny @***REMOVED***PROC***REMOVED***/***REMOVED***[^1-9],[^1-9][^0-9],[^1-9s][^0-9y][^0-9s],[^1-9][^0-9][^0-9][^0-9]****REMOVED***/** w,
  deny @***REMOVED***PROC***REMOVED***/sys/[^k]** w,  # deny /proc/sys except /proc/sys/k* (effectively /proc/sys/kernel)
  deny @***REMOVED***PROC***REMOVED***/sys/kernel/***REMOVED***?,??,[^s][^h][^m]*****REMOVED*** w,  # deny everything except shm* in /proc/sys/kernel/
  deny @***REMOVED***PROC***REMOVED***/sysrq-trigger rwklx,
  deny @***REMOVED***PROC***REMOVED***/kcore rwklx,

  deny mount,

  deny /sys/[^f]*/** wklx,
  deny /sys/f[^s]*/** wklx,
  deny /sys/fs/[^c]*/** wklx,
  deny /sys/fs/c[^g]*/** wklx,
  deny /sys/fs/cg[^r]*/** wklx,
  deny /sys/firmware/** rwklx,
  deny /sys/kernel/security/** rwklx,

***REMOVED******REMOVED***if ge .Version 208095***REMOVED******REMOVED***
  # suppress ptrace denials when using 'docker ps' or using 'ps' inside a container
  ptrace (trace,read) peer=***REMOVED******REMOVED***.Name***REMOVED******REMOVED***,
***REMOVED******REMOVED***end***REMOVED******REMOVED***
***REMOVED***
`
