// cgo -godefs types_darwin.go | go run mkpost.go
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build arm64,darwin

package unix

const (
	sizeofPtr      = 0x8
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x8
	sizeofLongLong = 0x8
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

type Timespec struct ***REMOVED***
	Sec  int64
	Nsec int64
***REMOVED***

type Timeval struct ***REMOVED***
	Sec       int64
	Usec      int32
	Pad_cgo_0 [4]byte
***REMOVED***

type Timeval32 struct ***REMOVED***
	Sec  int32
	Usec int32
***REMOVED***

type Rusage struct ***REMOVED***
	Utime    Timeval
	Stime    Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
***REMOVED***

type Rlimit struct ***REMOVED***
	Cur uint64
	Max uint64
***REMOVED***

type _Gid_t uint32

type Stat_t struct ***REMOVED***
	Dev           int32
	Mode          uint16
	Nlink         uint16
	Ino           uint64
	Uid           uint32
	Gid           uint32
	Rdev          int32
	Pad_cgo_0     [4]byte
	Atimespec     Timespec
	Mtimespec     Timespec
	Ctimespec     Timespec
	Birthtimespec Timespec
	Size          int64
	Blocks        int64
	Blksize       int32
	Flags         uint32
	Gen           uint32
	Lspare        int32
	Qspare        [2]int64
***REMOVED***

type Statfs_t struct ***REMOVED***
	Bsize       uint32
	Iosize      int32
	Blocks      uint64
	Bfree       uint64
	Bavail      uint64
	Files       uint64
	Ffree       uint64
	Fsid        Fsid
	Owner       uint32
	Type        uint32
	Flags       uint32
	Fssubtype   uint32
	Fstypename  [16]int8
	Mntonname   [1024]int8
	Mntfromname [1024]int8
	Reserved    [8]uint32
***REMOVED***

type Flock_t struct ***REMOVED***
	Start  int64
	Len    int64
	Pid    int32
	Type   int16
	Whence int16
***REMOVED***

type Fstore_t struct ***REMOVED***
	Flags      uint32
	Posmode    int32
	Offset     int64
	Length     int64
	Bytesalloc int64
***REMOVED***

type Radvisory_t struct ***REMOVED***
	Offset    int64
	Count     int32
	Pad_cgo_0 [4]byte
***REMOVED***

type Fbootstraptransfer_t struct ***REMOVED***
	Offset int64
	Length uint64
	Buffer *byte
***REMOVED***

type Log2phys_t struct ***REMOVED***
	Flags     uint32
	Pad_cgo_0 [8]byte
	Pad_cgo_1 [8]byte
***REMOVED***

type Fsid struct ***REMOVED***
	Val [2]int32
***REMOVED***

type Dirent struct ***REMOVED***
	Ino       uint64
	Seekoff   uint64
	Reclen    uint16
	Namlen    uint16
	Type      uint8
	Name      [1024]int8
	Pad_cgo_0 [3]byte
***REMOVED***

type RawSockaddrInet4 struct ***REMOVED***
	Len    uint8
	Family uint8
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]int8
***REMOVED***

type RawSockaddrInet6 struct ***REMOVED***
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
***REMOVED***

type RawSockaddrUnix struct ***REMOVED***
	Len    uint8
	Family uint8
	Path   [104]int8
***REMOVED***

type RawSockaddrDatalink struct ***REMOVED***
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [12]int8
***REMOVED***

type RawSockaddr struct ***REMOVED***
	Len    uint8
	Family uint8
	Data   [14]int8
***REMOVED***

type RawSockaddrAny struct ***REMOVED***
	Addr RawSockaddr
	Pad  [92]int8
***REMOVED***

type _Socklen uint32

type Linger struct ***REMOVED***
	Onoff  int32
	Linger int32
***REMOVED***

type Iovec struct ***REMOVED***
	Base *byte
	Len  uint64
***REMOVED***

type IPMreq struct ***REMOVED***
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
***REMOVED***

type IPv6Mreq struct ***REMOVED***
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
***REMOVED***

type Msghdr struct ***REMOVED***
	Name       *byte
	Namelen    uint32
	Pad_cgo_0  [4]byte
	Iov        *Iovec
	Iovlen     int32
	Pad_cgo_1  [4]byte
	Control    *byte
	Controllen uint32
	Flags      int32
***REMOVED***

type Cmsghdr struct ***REMOVED***
	Len   uint32
	Level int32
	Type  int32
***REMOVED***

type Inet4Pktinfo struct ***REMOVED***
	Ifindex  uint32
	Spec_dst [4]byte /* in_addr */
	Addr     [4]byte /* in_addr */
***REMOVED***

type Inet6Pktinfo struct ***REMOVED***
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
***REMOVED***

type IPv6MTUInfo struct ***REMOVED***
	Addr RawSockaddrInet6
	Mtu  uint32
***REMOVED***

type ICMPv6Filter struct ***REMOVED***
	Filt [8]uint32
***REMOVED***

const (
	SizeofSockaddrInet4    = 0x10
	SizeofSockaddrInet6    = 0x1c
	SizeofSockaddrAny      = 0x6c
	SizeofSockaddrUnix     = 0x6a
	SizeofSockaddrDatalink = 0x14
	SizeofLinger           = 0x8
	SizeofIPMreq           = 0x8
	SizeofIPv6Mreq         = 0x14
	SizeofMsghdr           = 0x30
	SizeofCmsghdr          = 0xc
	SizeofInet4Pktinfo     = 0xc
	SizeofInet6Pktinfo     = 0x14
	SizeofIPv6MTUInfo      = 0x20
	SizeofICMPv6Filter     = 0x20
)

const (
	PTRACE_TRACEME = 0x0
	PTRACE_CONT    = 0x7
	PTRACE_KILL    = 0x8
)

type Kevent_t struct ***REMOVED***
	Ident  uint64
	Filter int16
	Flags  uint16
	Fflags uint32
	Data   int64
	Udata  *byte
***REMOVED***

type FdSet struct ***REMOVED***
	Bits [32]int32
***REMOVED***

const (
	SizeofIfMsghdr    = 0x70
	SizeofIfData      = 0x60
	SizeofIfaMsghdr   = 0x14
	SizeofIfmaMsghdr  = 0x10
	SizeofIfmaMsghdr2 = 0x14
	SizeofRtMsghdr    = 0x5c
	SizeofRtMetrics   = 0x38
)

type IfMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
	Data      IfData
***REMOVED***

type IfData struct ***REMOVED***
	Type       uint8
	Typelen    uint8
	Physical   uint8
	Addrlen    uint8
	Hdrlen     uint8
	Recvquota  uint8
	Xmitquota  uint8
	Unused1    uint8
	Mtu        uint32
	Metric     uint32
	Baudrate   uint32
	Ipackets   uint32
	Ierrors    uint32
	Opackets   uint32
	Oerrors    uint32
	Collisions uint32
	Ibytes     uint32
	Obytes     uint32
	Imcasts    uint32
	Omcasts    uint32
	Iqdrops    uint32
	Noproto    uint32
	Recvtiming uint32
	Xmittiming uint32
	Lastchange Timeval32
	Unused2    uint32
	Hwassist   uint32
	Reserved1  uint32
	Reserved2  uint32
***REMOVED***

type IfaMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
	Metric    int32
***REMOVED***

type IfmaMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
***REMOVED***

type IfmaMsghdr2 struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
	Refcount  int32
***REMOVED***

type RtMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Index     uint16
	Pad_cgo_0 [2]byte
	Flags     int32
	Addrs     int32
	Pid       int32
	Seq       int32
	Errno     int32
	Use       int32
	Inits     uint32
	Rmx       RtMetrics
***REMOVED***

type RtMetrics struct ***REMOVED***
	Locks    uint32
	Mtu      uint32
	Hopcount uint32
	Expire   int32
	Recvpipe uint32
	Sendpipe uint32
	Ssthresh uint32
	Rtt      uint32
	Rttvar   uint32
	Pksent   uint32
	Filler   [4]uint32
***REMOVED***

const (
	SizeofBpfVersion = 0x4
	SizeofBpfStat    = 0x8
	SizeofBpfProgram = 0x10
	SizeofBpfInsn    = 0x8
	SizeofBpfHdr     = 0x14
)

type BpfVersion struct ***REMOVED***
	Major uint16
	Minor uint16
***REMOVED***

type BpfStat struct ***REMOVED***
	Recv uint32
	Drop uint32
***REMOVED***

type BpfProgram struct ***REMOVED***
	Len       uint32
	Pad_cgo_0 [4]byte
	Insns     *BpfInsn
***REMOVED***

type BpfInsn struct ***REMOVED***
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
***REMOVED***

type BpfHdr struct ***REMOVED***
	Tstamp    Timeval32
	Caplen    uint32
	Datalen   uint32
	Hdrlen    uint16
	Pad_cgo_0 [2]byte
***REMOVED***

type Termios struct ***REMOVED***
	Iflag     uint64
	Oflag     uint64
	Cflag     uint64
	Lflag     uint64
	Cc        [20]uint8
	Pad_cgo_0 [4]byte
	Ispeed    uint64
	Ospeed    uint64
***REMOVED***

type Winsize struct ***REMOVED***
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
***REMOVED***

const (
	AT_FDCWD            = -0x2
	AT_REMOVEDIR        = 0x80
	AT_SYMLINK_FOLLOW   = 0x40
	AT_SYMLINK_NOFOLLOW = 0x20
)

type PollFd struct ***REMOVED***
	Fd      int32
	Events  int16
	Revents int16
***REMOVED***

const (
	POLLERR    = 0x8
	POLLHUP    = 0x10
	POLLIN     = 0x1
	POLLNVAL   = 0x20
	POLLOUT    = 0x4
	POLLPRI    = 0x2
	POLLRDBAND = 0x80
	POLLRDNORM = 0x40
	POLLWRBAND = 0x100
	POLLWRNORM = 0x4
)

type Utsname struct ***REMOVED***
	Sysname  [256]byte
	Nodename [256]byte
	Release  [256]byte
	Version  [256]byte
	Machine  [256]byte
***REMOVED***