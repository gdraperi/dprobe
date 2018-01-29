// mksysnum_freebsd.pl
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build 386,freebsd

package unix

const (
	// SYS_NOSYS = 0;  // ***REMOVED*** int nosys(void); ***REMOVED*** syscall nosys_args int
	SYS_EXIT                   = 1   // ***REMOVED*** void sys_exit(int rval); ***REMOVED*** exit \
	SYS_FORK                   = 2   // ***REMOVED*** int fork(void); ***REMOVED***
	SYS_READ                   = 3   // ***REMOVED*** ssize_t read(int fd, void *buf, \
	SYS_WRITE                  = 4   // ***REMOVED*** ssize_t write(int fd, const void *buf, \
	SYS_OPEN                   = 5   // ***REMOVED*** int open(char *path, int flags, int mode); ***REMOVED***
	SYS_CLOSE                  = 6   // ***REMOVED*** int close(int fd); ***REMOVED***
	SYS_WAIT4                  = 7   // ***REMOVED*** int wait4(int pid, int *status, \
	SYS_LINK                   = 9   // ***REMOVED*** int link(char *path, char *link); ***REMOVED***
	SYS_UNLINK                 = 10  // ***REMOVED*** int unlink(char *path); ***REMOVED***
	SYS_CHDIR                  = 12  // ***REMOVED*** int chdir(char *path); ***REMOVED***
	SYS_FCHDIR                 = 13  // ***REMOVED*** int fchdir(int fd); ***REMOVED***
	SYS_MKNOD                  = 14  // ***REMOVED*** int mknod(char *path, int mode, int dev); ***REMOVED***
	SYS_CHMOD                  = 15  // ***REMOVED*** int chmod(char *path, int mode); ***REMOVED***
	SYS_CHOWN                  = 16  // ***REMOVED*** int chown(char *path, int uid, int gid); ***REMOVED***
	SYS_OBREAK                 = 17  // ***REMOVED*** int obreak(char *nsize); ***REMOVED*** break \
	SYS_GETPID                 = 20  // ***REMOVED*** pid_t getpid(void); ***REMOVED***
	SYS_MOUNT                  = 21  // ***REMOVED*** int mount(char *type, char *path, \
	SYS_UNMOUNT                = 22  // ***REMOVED*** int unmount(char *path, int flags); ***REMOVED***
	SYS_SETUID                 = 23  // ***REMOVED*** int setuid(uid_t uid); ***REMOVED***
	SYS_GETUID                 = 24  // ***REMOVED*** uid_t getuid(void); ***REMOVED***
	SYS_GETEUID                = 25  // ***REMOVED*** uid_t geteuid(void); ***REMOVED***
	SYS_PTRACE                 = 26  // ***REMOVED*** int ptrace(int req, pid_t pid, \
	SYS_RECVMSG                = 27  // ***REMOVED*** int recvmsg(int s, struct msghdr *msg, \
	SYS_SENDMSG                = 28  // ***REMOVED*** int sendmsg(int s, struct msghdr *msg, \
	SYS_RECVFROM               = 29  // ***REMOVED*** int recvfrom(int s, caddr_t buf, \
	SYS_ACCEPT                 = 30  // ***REMOVED*** int accept(int s, \
	SYS_GETPEERNAME            = 31  // ***REMOVED*** int getpeername(int fdes, \
	SYS_GETSOCKNAME            = 32  // ***REMOVED*** int getsockname(int fdes, \
	SYS_ACCESS                 = 33  // ***REMOVED*** int access(char *path, int amode); ***REMOVED***
	SYS_CHFLAGS                = 34  // ***REMOVED*** int chflags(const char *path, u_long flags); ***REMOVED***
	SYS_FCHFLAGS               = 35  // ***REMOVED*** int fchflags(int fd, u_long flags); ***REMOVED***
	SYS_SYNC                   = 36  // ***REMOVED*** int sync(void); ***REMOVED***
	SYS_KILL                   = 37  // ***REMOVED*** int kill(int pid, int signum); ***REMOVED***
	SYS_GETPPID                = 39  // ***REMOVED*** pid_t getppid(void); ***REMOVED***
	SYS_DUP                    = 41  // ***REMOVED*** int dup(u_int fd); ***REMOVED***
	SYS_PIPE                   = 42  // ***REMOVED*** int pipe(void); ***REMOVED***
	SYS_GETEGID                = 43  // ***REMOVED*** gid_t getegid(void); ***REMOVED***
	SYS_PROFIL                 = 44  // ***REMOVED*** int profil(caddr_t samples, size_t size, \
	SYS_KTRACE                 = 45  // ***REMOVED*** int ktrace(const char *fname, int ops, \
	SYS_GETGID                 = 47  // ***REMOVED*** gid_t getgid(void); ***REMOVED***
	SYS_GETLOGIN               = 49  // ***REMOVED*** int getlogin(char *namebuf, u_int \
	SYS_SETLOGIN               = 50  // ***REMOVED*** int setlogin(char *namebuf); ***REMOVED***
	SYS_ACCT                   = 51  // ***REMOVED*** int acct(char *path); ***REMOVED***
	SYS_SIGALTSTACK            = 53  // ***REMOVED*** int sigaltstack(stack_t *ss, \
	SYS_IOCTL                  = 54  // ***REMOVED*** int ioctl(int fd, u_long com, \
	SYS_REBOOT                 = 55  // ***REMOVED*** int reboot(int opt); ***REMOVED***
	SYS_REVOKE                 = 56  // ***REMOVED*** int revoke(char *path); ***REMOVED***
	SYS_SYMLINK                = 57  // ***REMOVED*** int symlink(char *path, char *link); ***REMOVED***
	SYS_READLINK               = 58  // ***REMOVED*** ssize_t readlink(char *path, char *buf, \
	SYS_EXECVE                 = 59  // ***REMOVED*** int execve(char *fname, char **argv, \
	SYS_UMASK                  = 60  // ***REMOVED*** int umask(int newmask); ***REMOVED*** umask umask_args \
	SYS_CHROOT                 = 61  // ***REMOVED*** int chroot(char *path); ***REMOVED***
	SYS_MSYNC                  = 65  // ***REMOVED*** int msync(void *addr, size_t len, \
	SYS_VFORK                  = 66  // ***REMOVED*** int vfork(void); ***REMOVED***
	SYS_SBRK                   = 69  // ***REMOVED*** int sbrk(int incr); ***REMOVED***
	SYS_SSTK                   = 70  // ***REMOVED*** int sstk(int incr); ***REMOVED***
	SYS_OVADVISE               = 72  // ***REMOVED*** int ovadvise(int anom); ***REMOVED*** vadvise \
	SYS_MUNMAP                 = 73  // ***REMOVED*** int munmap(void *addr, size_t len); ***REMOVED***
	SYS_MPROTECT               = 74  // ***REMOVED*** int mprotect(const void *addr, size_t len, \
	SYS_MADVISE                = 75  // ***REMOVED*** int madvise(void *addr, size_t len, \
	SYS_MINCORE                = 78  // ***REMOVED*** int mincore(const void *addr, size_t len, \
	SYS_GETGROUPS              = 79  // ***REMOVED*** int getgroups(u_int gidsetsize, \
	SYS_SETGROUPS              = 80  // ***REMOVED*** int setgroups(u_int gidsetsize, \
	SYS_GETPGRP                = 81  // ***REMOVED*** int getpgrp(void); ***REMOVED***
	SYS_SETPGID                = 82  // ***REMOVED*** int setpgid(int pid, int pgid); ***REMOVED***
	SYS_SETITIMER              = 83  // ***REMOVED*** int setitimer(u_int which, struct \
	SYS_SWAPON                 = 85  // ***REMOVED*** int swapon(char *name); ***REMOVED***
	SYS_GETITIMER              = 86  // ***REMOVED*** int getitimer(u_int which, \
	SYS_GETDTABLESIZE          = 89  // ***REMOVED*** int getdtablesize(void); ***REMOVED***
	SYS_DUP2                   = 90  // ***REMOVED*** int dup2(u_int from, u_int to); ***REMOVED***
	SYS_FCNTL                  = 92  // ***REMOVED*** int fcntl(int fd, int cmd, long arg); ***REMOVED***
	SYS_SELECT                 = 93  // ***REMOVED*** int select(int nd, fd_set *in, fd_set *ou, \
	SYS_FSYNC                  = 95  // ***REMOVED*** int fsync(int fd); ***REMOVED***
	SYS_SETPRIORITY            = 96  // ***REMOVED*** int setpriority(int which, int who, \
	SYS_SOCKET                 = 97  // ***REMOVED*** int socket(int domain, int type, \
	SYS_CONNECT                = 98  // ***REMOVED*** int connect(int s, caddr_t name, \
	SYS_GETPRIORITY            = 100 // ***REMOVED*** int getpriority(int which, int who); ***REMOVED***
	SYS_BIND                   = 104 // ***REMOVED*** int bind(int s, caddr_t name, \
	SYS_SETSOCKOPT             = 105 // ***REMOVED*** int setsockopt(int s, int level, int name, \
	SYS_LISTEN                 = 106 // ***REMOVED*** int listen(int s, int backlog); ***REMOVED***
	SYS_GETTIMEOFDAY           = 116 // ***REMOVED*** int gettimeofday(struct timeval *tp, \
	SYS_GETRUSAGE              = 117 // ***REMOVED*** int getrusage(int who, \
	SYS_GETSOCKOPT             = 118 // ***REMOVED*** int getsockopt(int s, int level, int name, \
	SYS_READV                  = 120 // ***REMOVED*** int readv(int fd, struct iovec *iovp, \
	SYS_WRITEV                 = 121 // ***REMOVED*** int writev(int fd, struct iovec *iovp, \
	SYS_SETTIMEOFDAY           = 122 // ***REMOVED*** int settimeofday(struct timeval *tv, \
	SYS_FCHOWN                 = 123 // ***REMOVED*** int fchown(int fd, int uid, int gid); ***REMOVED***
	SYS_FCHMOD                 = 124 // ***REMOVED*** int fchmod(int fd, int mode); ***REMOVED***
	SYS_SETREUID               = 126 // ***REMOVED*** int setreuid(int ruid, int euid); ***REMOVED***
	SYS_SETREGID               = 127 // ***REMOVED*** int setregid(int rgid, int egid); ***REMOVED***
	SYS_RENAME                 = 128 // ***REMOVED*** int rename(char *from, char *to); ***REMOVED***
	SYS_FLOCK                  = 131 // ***REMOVED*** int flock(int fd, int how); ***REMOVED***
	SYS_MKFIFO                 = 132 // ***REMOVED*** int mkfifo(char *path, int mode); ***REMOVED***
	SYS_SENDTO                 = 133 // ***REMOVED*** int sendto(int s, caddr_t buf, size_t len, \
	SYS_SHUTDOWN               = 134 // ***REMOVED*** int shutdown(int s, int how); ***REMOVED***
	SYS_SOCKETPAIR             = 135 // ***REMOVED*** int socketpair(int domain, int type, \
	SYS_MKDIR                  = 136 // ***REMOVED*** int mkdir(char *path, int mode); ***REMOVED***
	SYS_RMDIR                  = 137 // ***REMOVED*** int rmdir(char *path); ***REMOVED***
	SYS_UTIMES                 = 138 // ***REMOVED*** int utimes(char *path, \
	SYS_ADJTIME                = 140 // ***REMOVED*** int adjtime(struct timeval *delta, \
	SYS_SETSID                 = 147 // ***REMOVED*** int setsid(void); ***REMOVED***
	SYS_QUOTACTL               = 148 // ***REMOVED*** int quotactl(char *path, int cmd, int uid, \
	SYS_LGETFH                 = 160 // ***REMOVED*** int lgetfh(char *fname, \
	SYS_GETFH                  = 161 // ***REMOVED*** int getfh(char *fname, \
	SYS_SYSARCH                = 165 // ***REMOVED*** int sysarch(int op, char *parms); ***REMOVED***
	SYS_RTPRIO                 = 166 // ***REMOVED*** int rtprio(int function, pid_t pid, \
	SYS_FREEBSD6_PREAD         = 173 // ***REMOVED*** ssize_t freebsd6_pread(int fd, void *buf, \
	SYS_FREEBSD6_PWRITE        = 174 // ***REMOVED*** ssize_t freebsd6_pwrite(int fd, \
	SYS_SETFIB                 = 175 // ***REMOVED*** int setfib(int fibnum); ***REMOVED***
	SYS_NTP_ADJTIME            = 176 // ***REMOVED*** int ntp_adjtime(struct timex *tp); ***REMOVED***
	SYS_SETGID                 = 181 // ***REMOVED*** int setgid(gid_t gid); ***REMOVED***
	SYS_SETEGID                = 182 // ***REMOVED*** int setegid(gid_t egid); ***REMOVED***
	SYS_SETEUID                = 183 // ***REMOVED*** int seteuid(uid_t euid); ***REMOVED***
	SYS_STAT                   = 188 // ***REMOVED*** int stat(char *path, struct stat *ub); ***REMOVED***
	SYS_FSTAT                  = 189 // ***REMOVED*** int fstat(int fd, struct stat *sb); ***REMOVED***
	SYS_LSTAT                  = 190 // ***REMOVED*** int lstat(char *path, struct stat *ub); ***REMOVED***
	SYS_PATHCONF               = 191 // ***REMOVED*** int pathconf(char *path, int name); ***REMOVED***
	SYS_FPATHCONF              = 192 // ***REMOVED*** int fpathconf(int fd, int name); ***REMOVED***
	SYS_GETRLIMIT              = 194 // ***REMOVED*** int getrlimit(u_int which, \
	SYS_SETRLIMIT              = 195 // ***REMOVED*** int setrlimit(u_int which, \
	SYS_GETDIRENTRIES          = 196 // ***REMOVED*** int getdirentries(int fd, char *buf, \
	SYS_FREEBSD6_MMAP          = 197 // ***REMOVED*** caddr_t freebsd6_mmap(caddr_t addr, \
	SYS_FREEBSD6_LSEEK         = 199 // ***REMOVED*** off_t freebsd6_lseek(int fd, int pad, \
	SYS_FREEBSD6_TRUNCATE      = 200 // ***REMOVED*** int freebsd6_truncate(char *path, int pad, \
	SYS_FREEBSD6_FTRUNCATE     = 201 // ***REMOVED*** int freebsd6_ftruncate(int fd, int pad, \
	SYS___SYSCTL               = 202 // ***REMOVED*** int __sysctl(int *name, u_int namelen, \
	SYS_MLOCK                  = 203 // ***REMOVED*** int mlock(const void *addr, size_t len); ***REMOVED***
	SYS_MUNLOCK                = 204 // ***REMOVED*** int munlock(const void *addr, size_t len); ***REMOVED***
	SYS_UNDELETE               = 205 // ***REMOVED*** int undelete(char *path); ***REMOVED***
	SYS_FUTIMES                = 206 // ***REMOVED*** int futimes(int fd, struct timeval *tptr); ***REMOVED***
	SYS_GETPGID                = 207 // ***REMOVED*** int getpgid(pid_t pid); ***REMOVED***
	SYS_POLL                   = 209 // ***REMOVED*** int poll(struct pollfd *fds, u_int nfds, \
	SYS_CLOCK_GETTIME          = 232 // ***REMOVED*** int clock_gettime(clockid_t clock_id, \
	SYS_CLOCK_SETTIME          = 233 // ***REMOVED*** int clock_settime( \
	SYS_CLOCK_GETRES           = 234 // ***REMOVED*** int clock_getres(clockid_t clock_id, \
	SYS_KTIMER_CREATE          = 235 // ***REMOVED*** int ktimer_create(clockid_t clock_id, \
	SYS_KTIMER_DELETE          = 236 // ***REMOVED*** int ktimer_delete(int timerid); ***REMOVED***
	SYS_KTIMER_SETTIME         = 237 // ***REMOVED*** int ktimer_settime(int timerid, int flags, \
	SYS_KTIMER_GETTIME         = 238 // ***REMOVED*** int ktimer_gettime(int timerid, struct \
	SYS_KTIMER_GETOVERRUN      = 239 // ***REMOVED*** int ktimer_getoverrun(int timerid); ***REMOVED***
	SYS_NANOSLEEP              = 240 // ***REMOVED*** int nanosleep(const struct timespec *rqtp, \
	SYS_FFCLOCK_GETCOUNTER     = 241 // ***REMOVED*** int ffclock_getcounter(ffcounter *ffcount); ***REMOVED***
	SYS_FFCLOCK_SETESTIMATE    = 242 // ***REMOVED*** int ffclock_setestimate( \
	SYS_FFCLOCK_GETESTIMATE    = 243 // ***REMOVED*** int ffclock_getestimate( \
	SYS_CLOCK_GETCPUCLOCKID2   = 247 // ***REMOVED*** int clock_getcpuclockid2(id_t id,\
	SYS_NTP_GETTIME            = 248 // ***REMOVED*** int ntp_gettime(struct ntptimeval *ntvp); ***REMOVED***
	SYS_MINHERIT               = 250 // ***REMOVED*** int minherit(void *addr, size_t len, \
	SYS_RFORK                  = 251 // ***REMOVED*** int rfork(int flags); ***REMOVED***
	SYS_OPENBSD_POLL           = 252 // ***REMOVED*** int openbsd_poll(struct pollfd *fds, \
	SYS_ISSETUGID              = 253 // ***REMOVED*** int issetugid(void); ***REMOVED***
	SYS_LCHOWN                 = 254 // ***REMOVED*** int lchown(char *path, int uid, int gid); ***REMOVED***
	SYS_GETDENTS               = 272 // ***REMOVED*** int getdents(int fd, char *buf, \
	SYS_LCHMOD                 = 274 // ***REMOVED*** int lchmod(char *path, mode_t mode); ***REMOVED***
	SYS_LUTIMES                = 276 // ***REMOVED*** int lutimes(char *path, \
	SYS_NSTAT                  = 278 // ***REMOVED*** int nstat(char *path, struct nstat *ub); ***REMOVED***
	SYS_NFSTAT                 = 279 // ***REMOVED*** int nfstat(int fd, struct nstat *sb); ***REMOVED***
	SYS_NLSTAT                 = 280 // ***REMOVED*** int nlstat(char *path, struct nstat *ub); ***REMOVED***
	SYS_PREADV                 = 289 // ***REMOVED*** ssize_t preadv(int fd, struct iovec *iovp, \
	SYS_PWRITEV                = 290 // ***REMOVED*** ssize_t pwritev(int fd, struct iovec *iovp, \
	SYS_FHOPEN                 = 298 // ***REMOVED*** int fhopen(const struct fhandle *u_fhp, \
	SYS_FHSTAT                 = 299 // ***REMOVED*** int fhstat(const struct fhandle *u_fhp, \
	SYS_MODNEXT                = 300 // ***REMOVED*** int modnext(int modid); ***REMOVED***
	SYS_MODSTAT                = 301 // ***REMOVED*** int modstat(int modid, \
	SYS_MODFNEXT               = 302 // ***REMOVED*** int modfnext(int modid); ***REMOVED***
	SYS_MODFIND                = 303 // ***REMOVED*** int modfind(const char *name); ***REMOVED***
	SYS_KLDLOAD                = 304 // ***REMOVED*** int kldload(const char *file); ***REMOVED***
	SYS_KLDUNLOAD              = 305 // ***REMOVED*** int kldunload(int fileid); ***REMOVED***
	SYS_KLDFIND                = 306 // ***REMOVED*** int kldfind(const char *file); ***REMOVED***
	SYS_KLDNEXT                = 307 // ***REMOVED*** int kldnext(int fileid); ***REMOVED***
	SYS_KLDSTAT                = 308 // ***REMOVED*** int kldstat(int fileid, struct \
	SYS_KLDFIRSTMOD            = 309 // ***REMOVED*** int kldfirstmod(int fileid); ***REMOVED***
	SYS_GETSID                 = 310 // ***REMOVED*** int getsid(pid_t pid); ***REMOVED***
	SYS_SETRESUID              = 311 // ***REMOVED*** int setresuid(uid_t ruid, uid_t euid, \
	SYS_SETRESGID              = 312 // ***REMOVED*** int setresgid(gid_t rgid, gid_t egid, \
	SYS_YIELD                  = 321 // ***REMOVED*** int yield(void); ***REMOVED***
	SYS_MLOCKALL               = 324 // ***REMOVED*** int mlockall(int how); ***REMOVED***
	SYS_MUNLOCKALL             = 325 // ***REMOVED*** int munlockall(void); ***REMOVED***
	SYS___GETCWD               = 326 // ***REMOVED*** int __getcwd(char *buf, u_int buflen); ***REMOVED***
	SYS_SCHED_SETPARAM         = 327 // ***REMOVED*** int sched_setparam (pid_t pid, \
	SYS_SCHED_GETPARAM         = 328 // ***REMOVED*** int sched_getparam (pid_t pid, struct \
	SYS_SCHED_SETSCHEDULER     = 329 // ***REMOVED*** int sched_setscheduler (pid_t pid, int \
	SYS_SCHED_GETSCHEDULER     = 330 // ***REMOVED*** int sched_getscheduler (pid_t pid); ***REMOVED***
	SYS_SCHED_YIELD            = 331 // ***REMOVED*** int sched_yield (void); ***REMOVED***
	SYS_SCHED_GET_PRIORITY_MAX = 332 // ***REMOVED*** int sched_get_priority_max (int policy); ***REMOVED***
	SYS_SCHED_GET_PRIORITY_MIN = 333 // ***REMOVED*** int sched_get_priority_min (int policy); ***REMOVED***
	SYS_SCHED_RR_GET_INTERVAL  = 334 // ***REMOVED*** int sched_rr_get_interval (pid_t pid, \
	SYS_UTRACE                 = 335 // ***REMOVED*** int utrace(const void *addr, size_t len); ***REMOVED***
	SYS_KLDSYM                 = 337 // ***REMOVED*** int kldsym(int fileid, int cmd, \
	SYS_JAIL                   = 338 // ***REMOVED*** int jail(struct jail *jail); ***REMOVED***
	SYS_SIGPROCMASK            = 340 // ***REMOVED*** int sigprocmask(int how, \
	SYS_SIGSUSPEND             = 341 // ***REMOVED*** int sigsuspend(const sigset_t *sigmask); ***REMOVED***
	SYS_SIGPENDING             = 343 // ***REMOVED*** int sigpending(sigset_t *set); ***REMOVED***
	SYS_SIGTIMEDWAIT           = 345 // ***REMOVED*** int sigtimedwait(const sigset_t *set, \
	SYS_SIGWAITINFO            = 346 // ***REMOVED*** int sigwaitinfo(const sigset_t *set, \
	SYS___ACL_GET_FILE         = 347 // ***REMOVED*** int __acl_get_file(const char *path, \
	SYS___ACL_SET_FILE         = 348 // ***REMOVED*** int __acl_set_file(const char *path, \
	SYS___ACL_GET_FD           = 349 // ***REMOVED*** int __acl_get_fd(int filedes, \
	SYS___ACL_SET_FD           = 350 // ***REMOVED*** int __acl_set_fd(int filedes, \
	SYS___ACL_DELETE_FILE      = 351 // ***REMOVED*** int __acl_delete_file(const char *path, \
	SYS___ACL_DELETE_FD        = 352 // ***REMOVED*** int __acl_delete_fd(int filedes, \
	SYS___ACL_ACLCHECK_FILE    = 353 // ***REMOVED*** int __acl_aclcheck_file(const char *path, \
	SYS___ACL_ACLCHECK_FD      = 354 // ***REMOVED*** int __acl_aclcheck_fd(int filedes, \
	SYS_EXTATTRCTL             = 355 // ***REMOVED*** int extattrctl(const char *path, int cmd, \
	SYS_EXTATTR_SET_FILE       = 356 // ***REMOVED*** ssize_t extattr_set_file( \
	SYS_EXTATTR_GET_FILE       = 357 // ***REMOVED*** ssize_t extattr_get_file( \
	SYS_EXTATTR_DELETE_FILE    = 358 // ***REMOVED*** int extattr_delete_file(const char *path, \
	SYS_GETRESUID              = 360 // ***REMOVED*** int getresuid(uid_t *ruid, uid_t *euid, \
	SYS_GETRESGID              = 361 // ***REMOVED*** int getresgid(gid_t *rgid, gid_t *egid, \
	SYS_KQUEUE                 = 362 // ***REMOVED*** int kqueue(void); ***REMOVED***
	SYS_KEVENT                 = 363 // ***REMOVED*** int kevent(int fd, \
	SYS_EXTATTR_SET_FD         = 371 // ***REMOVED*** ssize_t extattr_set_fd(int fd, \
	SYS_EXTATTR_GET_FD         = 372 // ***REMOVED*** ssize_t extattr_get_fd(int fd, \
	SYS_EXTATTR_DELETE_FD      = 373 // ***REMOVED*** int extattr_delete_fd(int fd, \
	SYS___SETUGID              = 374 // ***REMOVED*** int __setugid(int flag); ***REMOVED***
	SYS_EACCESS                = 376 // ***REMOVED*** int eaccess(char *path, int amode); ***REMOVED***
	SYS_NMOUNT                 = 378 // ***REMOVED*** int nmount(struct iovec *iovp, \
	SYS___MAC_GET_PROC         = 384 // ***REMOVED*** int __mac_get_proc(struct mac *mac_p); ***REMOVED***
	SYS___MAC_SET_PROC         = 385 // ***REMOVED*** int __mac_set_proc(struct mac *mac_p); ***REMOVED***
	SYS___MAC_GET_FD           = 386 // ***REMOVED*** int __mac_get_fd(int fd, \
	SYS___MAC_GET_FILE         = 387 // ***REMOVED*** int __mac_get_file(const char *path_p, \
	SYS___MAC_SET_FD           = 388 // ***REMOVED*** int __mac_set_fd(int fd, \
	SYS___MAC_SET_FILE         = 389 // ***REMOVED*** int __mac_set_file(const char *path_p, \
	SYS_KENV                   = 390 // ***REMOVED*** int kenv(int what, const char *name, \
	SYS_LCHFLAGS               = 391 // ***REMOVED*** int lchflags(const char *path, \
	SYS_UUIDGEN                = 392 // ***REMOVED*** int uuidgen(struct uuid *store, \
	SYS_SENDFILE               = 393 // ***REMOVED*** int sendfile(int fd, int s, off_t offset, \
	SYS_MAC_SYSCALL            = 394 // ***REMOVED*** int mac_syscall(const char *policy, \
	SYS_GETFSSTAT              = 395 // ***REMOVED*** int getfsstat(struct statfs *buf, \
	SYS_STATFS                 = 396 // ***REMOVED*** int statfs(char *path, \
	SYS_FSTATFS                = 397 // ***REMOVED*** int fstatfs(int fd, struct statfs *buf); ***REMOVED***
	SYS_FHSTATFS               = 398 // ***REMOVED*** int fhstatfs(const struct fhandle *u_fhp, \
	SYS___MAC_GET_PID          = 409 // ***REMOVED*** int __mac_get_pid(pid_t pid, \
	SYS___MAC_GET_LINK         = 410 // ***REMOVED*** int __mac_get_link(const char *path_p, \
	SYS___MAC_SET_LINK         = 411 // ***REMOVED*** int __mac_set_link(const char *path_p, \
	SYS_EXTATTR_SET_LINK       = 412 // ***REMOVED*** ssize_t extattr_set_link( \
	SYS_EXTATTR_GET_LINK       = 413 // ***REMOVED*** ssize_t extattr_get_link( \
	SYS_EXTATTR_DELETE_LINK    = 414 // ***REMOVED*** int extattr_delete_link( \
	SYS___MAC_EXECVE           = 415 // ***REMOVED*** int __mac_execve(char *fname, char **argv, \
	SYS_SIGACTION              = 416 // ***REMOVED*** int sigaction(int sig, \
	SYS_SIGRETURN              = 417 // ***REMOVED*** int sigreturn( \
	SYS_GETCONTEXT             = 421 // ***REMOVED*** int getcontext(struct __ucontext *ucp); ***REMOVED***
	SYS_SETCONTEXT             = 422 // ***REMOVED*** int setcontext( \
	SYS_SWAPCONTEXT            = 423 // ***REMOVED*** int swapcontext(struct __ucontext *oucp, \
	SYS_SWAPOFF                = 424 // ***REMOVED*** int swapoff(const char *name); ***REMOVED***
	SYS___ACL_GET_LINK         = 425 // ***REMOVED*** int __acl_get_link(const char *path, \
	SYS___ACL_SET_LINK         = 426 // ***REMOVED*** int __acl_set_link(const char *path, \
	SYS___ACL_DELETE_LINK      = 427 // ***REMOVED*** int __acl_delete_link(const char *path, \
	SYS___ACL_ACLCHECK_LINK    = 428 // ***REMOVED*** int __acl_aclcheck_link(const char *path, \
	SYS_SIGWAIT                = 429 // ***REMOVED*** int sigwait(const sigset_t *set, \
	SYS_THR_CREATE             = 430 // ***REMOVED*** int thr_create(ucontext_t *ctx, long *id, \
	SYS_THR_EXIT               = 431 // ***REMOVED*** void thr_exit(long *state); ***REMOVED***
	SYS_THR_SELF               = 432 // ***REMOVED*** int thr_self(long *id); ***REMOVED***
	SYS_THR_KILL               = 433 // ***REMOVED*** int thr_kill(long id, int sig); ***REMOVED***
	SYS__UMTX_LOCK             = 434 // ***REMOVED*** int _umtx_lock(struct umtx *umtx); ***REMOVED***
	SYS__UMTX_UNLOCK           = 435 // ***REMOVED*** int _umtx_unlock(struct umtx *umtx); ***REMOVED***
	SYS_JAIL_ATTACH            = 436 // ***REMOVED*** int jail_attach(int jid); ***REMOVED***
	SYS_EXTATTR_LIST_FD        = 437 // ***REMOVED*** ssize_t extattr_list_fd(int fd, \
	SYS_EXTATTR_LIST_FILE      = 438 // ***REMOVED*** ssize_t extattr_list_file( \
	SYS_EXTATTR_LIST_LINK      = 439 // ***REMOVED*** ssize_t extattr_list_link( \
	SYS_THR_SUSPEND            = 442 // ***REMOVED*** int thr_suspend( \
	SYS_THR_WAKE               = 443 // ***REMOVED*** int thr_wake(long id); ***REMOVED***
	SYS_KLDUNLOADF             = 444 // ***REMOVED*** int kldunloadf(int fileid, int flags); ***REMOVED***
	SYS_AUDIT                  = 445 // ***REMOVED*** int audit(const void *record, \
	SYS_AUDITON                = 446 // ***REMOVED*** int auditon(int cmd, void *data, \
	SYS_GETAUID                = 447 // ***REMOVED*** int getauid(uid_t *auid); ***REMOVED***
	SYS_SETAUID                = 448 // ***REMOVED*** int setauid(uid_t *auid); ***REMOVED***
	SYS_GETAUDIT               = 449 // ***REMOVED*** int getaudit(struct auditinfo *auditinfo); ***REMOVED***
	SYS_SETAUDIT               = 450 // ***REMOVED*** int setaudit(struct auditinfo *auditinfo); ***REMOVED***
	SYS_GETAUDIT_ADDR          = 451 // ***REMOVED*** int getaudit_addr( \
	SYS_SETAUDIT_ADDR          = 452 // ***REMOVED*** int setaudit_addr( \
	SYS_AUDITCTL               = 453 // ***REMOVED*** int auditctl(char *path); ***REMOVED***
	SYS__UMTX_OP               = 454 // ***REMOVED*** int _umtx_op(void *obj, int op, \
	SYS_THR_NEW                = 455 // ***REMOVED*** int thr_new(struct thr_param *param, \
	SYS_SIGQUEUE               = 456 // ***REMOVED*** int sigqueue(pid_t pid, int signum, void *value); ***REMOVED***
	SYS_ABORT2                 = 463 // ***REMOVED*** int abort2(const char *why, int nargs, void **args); ***REMOVED***
	SYS_THR_SET_NAME           = 464 // ***REMOVED*** int thr_set_name(long id, const char *name); ***REMOVED***
	SYS_RTPRIO_THREAD          = 466 // ***REMOVED*** int rtprio_thread(int function, \
	SYS_PREAD                  = 475 // ***REMOVED*** ssize_t pread(int fd, void *buf, \
	SYS_PWRITE                 = 476 // ***REMOVED*** ssize_t pwrite(int fd, const void *buf, \
	SYS_MMAP                   = 477 // ***REMOVED*** caddr_t mmap(caddr_t addr, size_t len, \
	SYS_LSEEK                  = 478 // ***REMOVED*** off_t lseek(int fd, off_t offset, \
	SYS_TRUNCATE               = 479 // ***REMOVED*** int truncate(char *path, off_t length); ***REMOVED***
	SYS_FTRUNCATE              = 480 // ***REMOVED*** int ftruncate(int fd, off_t length); ***REMOVED***
	SYS_THR_KILL2              = 481 // ***REMOVED*** int thr_kill2(pid_t pid, long id, int sig); ***REMOVED***
	SYS_SHM_OPEN               = 482 // ***REMOVED*** int shm_open(const char *path, int flags, \
	SYS_SHM_UNLINK             = 483 // ***REMOVED*** int shm_unlink(const char *path); ***REMOVED***
	SYS_CPUSET                 = 484 // ***REMOVED*** int cpuset(cpusetid_t *setid); ***REMOVED***
	SYS_CPUSET_SETID           = 485 // ***REMOVED*** int cpuset_setid(cpuwhich_t which, id_t id, \
	SYS_CPUSET_GETID           = 486 // ***REMOVED*** int cpuset_getid(cpulevel_t level, \
	SYS_CPUSET_GETAFFINITY     = 487 // ***REMOVED*** int cpuset_getaffinity(cpulevel_t level, \
	SYS_CPUSET_SETAFFINITY     = 488 // ***REMOVED*** int cpuset_setaffinity(cpulevel_t level, \
	SYS_FACCESSAT              = 489 // ***REMOVED*** int faccessat(int fd, char *path, int amode, \
	SYS_FCHMODAT               = 490 // ***REMOVED*** int fchmodat(int fd, char *path, mode_t mode, \
	SYS_FCHOWNAT               = 491 // ***REMOVED*** int fchownat(int fd, char *path, uid_t uid, \
	SYS_FEXECVE                = 492 // ***REMOVED*** int fexecve(int fd, char **argv, \
	SYS_FSTATAT                = 493 // ***REMOVED*** int fstatat(int fd, char *path, \
	SYS_FUTIMESAT              = 494 // ***REMOVED*** int futimesat(int fd, char *path, \
	SYS_LINKAT                 = 495 // ***REMOVED*** int linkat(int fd1, char *path1, int fd2, \
	SYS_MKDIRAT                = 496 // ***REMOVED*** int mkdirat(int fd, char *path, mode_t mode); ***REMOVED***
	SYS_MKFIFOAT               = 497 // ***REMOVED*** int mkfifoat(int fd, char *path, mode_t mode); ***REMOVED***
	SYS_MKNODAT                = 498 // ***REMOVED*** int mknodat(int fd, char *path, mode_t mode, \
	SYS_OPENAT                 = 499 // ***REMOVED*** int openat(int fd, char *path, int flag, \
	SYS_READLINKAT             = 500 // ***REMOVED*** int readlinkat(int fd, char *path, char *buf, \
	SYS_RENAMEAT               = 501 // ***REMOVED*** int renameat(int oldfd, char *old, int newfd, \
	SYS_SYMLINKAT              = 502 // ***REMOVED*** int symlinkat(char *path1, int fd, \
	SYS_UNLINKAT               = 503 // ***REMOVED*** int unlinkat(int fd, char *path, int flag); ***REMOVED***
	SYS_POSIX_OPENPT           = 504 // ***REMOVED*** int posix_openpt(int flags); ***REMOVED***
	SYS_JAIL_GET               = 506 // ***REMOVED*** int jail_get(struct iovec *iovp, \
	SYS_JAIL_SET               = 507 // ***REMOVED*** int jail_set(struct iovec *iovp, \
	SYS_JAIL_REMOVE            = 508 // ***REMOVED*** int jail_remove(int jid); ***REMOVED***
	SYS_CLOSEFROM              = 509 // ***REMOVED*** int closefrom(int lowfd); ***REMOVED***
	SYS_LPATHCONF              = 513 // ***REMOVED*** int lpathconf(char *path, int name); ***REMOVED***
	SYS___CAP_RIGHTS_GET       = 515 // ***REMOVED*** int __cap_rights_get(int version, \
	SYS_CAP_ENTER              = 516 // ***REMOVED*** int cap_enter(void); ***REMOVED***
	SYS_CAP_GETMODE            = 517 // ***REMOVED*** int cap_getmode(u_int *modep); ***REMOVED***
	SYS_PDFORK                 = 518 // ***REMOVED*** int pdfork(int *fdp, int flags); ***REMOVED***
	SYS_PDKILL                 = 519 // ***REMOVED*** int pdkill(int fd, int signum); ***REMOVED***
	SYS_PDGETPID               = 520 // ***REMOVED*** int pdgetpid(int fd, pid_t *pidp); ***REMOVED***
	SYS_PSELECT                = 522 // ***REMOVED*** int pselect(int nd, fd_set *in, \
	SYS_GETLOGINCLASS          = 523 // ***REMOVED*** int getloginclass(char *namebuf, \
	SYS_SETLOGINCLASS          = 524 // ***REMOVED*** int setloginclass(const char *namebuf); ***REMOVED***
	SYS_RCTL_GET_RACCT         = 525 // ***REMOVED*** int rctl_get_racct(const void *inbufp, \
	SYS_RCTL_GET_RULES         = 526 // ***REMOVED*** int rctl_get_rules(const void *inbufp, \
	SYS_RCTL_GET_LIMITS        = 527 // ***REMOVED*** int rctl_get_limits(const void *inbufp, \
	SYS_RCTL_ADD_RULE          = 528 // ***REMOVED*** int rctl_add_rule(const void *inbufp, \
	SYS_RCTL_REMOVE_RULE       = 529 // ***REMOVED*** int rctl_remove_rule(const void *inbufp, \
	SYS_POSIX_FALLOCATE        = 530 // ***REMOVED*** int posix_fallocate(int fd, \
	SYS_POSIX_FADVISE          = 531 // ***REMOVED*** int posix_fadvise(int fd, off_t offset, \
	SYS_WAIT6                  = 532 // ***REMOVED*** int wait6(idtype_t idtype, id_t id, \
	SYS_CAP_RIGHTS_LIMIT       = 533 // ***REMOVED*** int cap_rights_limit(int fd, \
	SYS_CAP_IOCTLS_LIMIT       = 534 // ***REMOVED*** int cap_ioctls_limit(int fd, \
	SYS_CAP_IOCTLS_GET         = 535 // ***REMOVED*** ssize_t cap_ioctls_get(int fd, \
	SYS_CAP_FCNTLS_LIMIT       = 536 // ***REMOVED*** int cap_fcntls_limit(int fd, \
	SYS_CAP_FCNTLS_GET         = 537 // ***REMOVED*** int cap_fcntls_get(int fd, \
	SYS_BINDAT                 = 538 // ***REMOVED*** int bindat(int fd, int s, caddr_t name, \
	SYS_CONNECTAT              = 539 // ***REMOVED*** int connectat(int fd, int s, caddr_t name, \
	SYS_CHFLAGSAT              = 540 // ***REMOVED*** int chflagsat(int fd, const char *path, \
	SYS_ACCEPT4                = 541 // ***REMOVED*** int accept4(int s, \
	SYS_PIPE2                  = 542 // ***REMOVED*** int pipe2(int *fildes, int flags); ***REMOVED***
	SYS_PROCCTL                = 544 // ***REMOVED*** int procctl(idtype_t idtype, id_t id, \
	SYS_PPOLL                  = 545 // ***REMOVED*** int ppoll(struct pollfd *fds, u_int nfds, \
	SYS_FUTIMENS               = 546 // ***REMOVED*** int futimens(int fd, \
	SYS_UTIMENSAT              = 547 // ***REMOVED*** int utimensat(int fd, \
)