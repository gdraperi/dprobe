// mksysnum_netbsd.pl
// MACHINE GENERATED BY THE ABOVE COMMAND; DO NOT EDIT

// +build 386,netbsd

package unix

const (
	SYS_EXIT                 = 1   // ***REMOVED*** void|sys||exit(int rval); ***REMOVED***
	SYS_FORK                 = 2   // ***REMOVED*** int|sys||fork(void); ***REMOVED***
	SYS_READ                 = 3   // ***REMOVED*** ssize_t|sys||read(int fd, void *buf, size_t nbyte); ***REMOVED***
	SYS_WRITE                = 4   // ***REMOVED*** ssize_t|sys||write(int fd, const void *buf, size_t nbyte); ***REMOVED***
	SYS_OPEN                 = 5   // ***REMOVED*** int|sys||open(const char *path, int flags, ... mode_t mode); ***REMOVED***
	SYS_CLOSE                = 6   // ***REMOVED*** int|sys||close(int fd); ***REMOVED***
	SYS_LINK                 = 9   // ***REMOVED*** int|sys||link(const char *path, const char *link); ***REMOVED***
	SYS_UNLINK               = 10  // ***REMOVED*** int|sys||unlink(const char *path); ***REMOVED***
	SYS_CHDIR                = 12  // ***REMOVED*** int|sys||chdir(const char *path); ***REMOVED***
	SYS_FCHDIR               = 13  // ***REMOVED*** int|sys||fchdir(int fd); ***REMOVED***
	SYS_CHMOD                = 15  // ***REMOVED*** int|sys||chmod(const char *path, mode_t mode); ***REMOVED***
	SYS_CHOWN                = 16  // ***REMOVED*** int|sys||chown(const char *path, uid_t uid, gid_t gid); ***REMOVED***
	SYS_BREAK                = 17  // ***REMOVED*** int|sys||obreak(char *nsize); ***REMOVED***
	SYS_GETPID               = 20  // ***REMOVED*** pid_t|sys||getpid_with_ppid(void); ***REMOVED***
	SYS_UNMOUNT              = 22  // ***REMOVED*** int|sys||unmount(const char *path, int flags); ***REMOVED***
	SYS_SETUID               = 23  // ***REMOVED*** int|sys||setuid(uid_t uid); ***REMOVED***
	SYS_GETUID               = 24  // ***REMOVED*** uid_t|sys||getuid_with_euid(void); ***REMOVED***
	SYS_GETEUID              = 25  // ***REMOVED*** uid_t|sys||geteuid(void); ***REMOVED***
	SYS_PTRACE               = 26  // ***REMOVED*** int|sys||ptrace(int req, pid_t pid, void *addr, int data); ***REMOVED***
	SYS_RECVMSG              = 27  // ***REMOVED*** ssize_t|sys||recvmsg(int s, struct msghdr *msg, int flags); ***REMOVED***
	SYS_SENDMSG              = 28  // ***REMOVED*** ssize_t|sys||sendmsg(int s, const struct msghdr *msg, int flags); ***REMOVED***
	SYS_RECVFROM             = 29  // ***REMOVED*** ssize_t|sys||recvfrom(int s, void *buf, size_t len, int flags, struct sockaddr *from, socklen_t *fromlenaddr); ***REMOVED***
	SYS_ACCEPT               = 30  // ***REMOVED*** int|sys||accept(int s, struct sockaddr *name, socklen_t *anamelen); ***REMOVED***
	SYS_GETPEERNAME          = 31  // ***REMOVED*** int|sys||getpeername(int fdes, struct sockaddr *asa, socklen_t *alen); ***REMOVED***
	SYS_GETSOCKNAME          = 32  // ***REMOVED*** int|sys||getsockname(int fdes, struct sockaddr *asa, socklen_t *alen); ***REMOVED***
	SYS_ACCESS               = 33  // ***REMOVED*** int|sys||access(const char *path, int flags); ***REMOVED***
	SYS_CHFLAGS              = 34  // ***REMOVED*** int|sys||chflags(const char *path, u_long flags); ***REMOVED***
	SYS_FCHFLAGS             = 35  // ***REMOVED*** int|sys||fchflags(int fd, u_long flags); ***REMOVED***
	SYS_SYNC                 = 36  // ***REMOVED*** void|sys||sync(void); ***REMOVED***
	SYS_KILL                 = 37  // ***REMOVED*** int|sys||kill(pid_t pid, int signum); ***REMOVED***
	SYS_GETPPID              = 39  // ***REMOVED*** pid_t|sys||getppid(void); ***REMOVED***
	SYS_DUP                  = 41  // ***REMOVED*** int|sys||dup(int fd); ***REMOVED***
	SYS_PIPE                 = 42  // ***REMOVED*** int|sys||pipe(void); ***REMOVED***
	SYS_GETEGID              = 43  // ***REMOVED*** gid_t|sys||getegid(void); ***REMOVED***
	SYS_PROFIL               = 44  // ***REMOVED*** int|sys||profil(char *samples, size_t size, u_long offset, u_int scale); ***REMOVED***
	SYS_KTRACE               = 45  // ***REMOVED*** int|sys||ktrace(const char *fname, int ops, int facs, pid_t pid); ***REMOVED***
	SYS_GETGID               = 47  // ***REMOVED*** gid_t|sys||getgid_with_egid(void); ***REMOVED***
	SYS___GETLOGIN           = 49  // ***REMOVED*** int|sys||__getlogin(char *namebuf, size_t namelen); ***REMOVED***
	SYS___SETLOGIN           = 50  // ***REMOVED*** int|sys||__setlogin(const char *namebuf); ***REMOVED***
	SYS_ACCT                 = 51  // ***REMOVED*** int|sys||acct(const char *path); ***REMOVED***
	SYS_IOCTL                = 54  // ***REMOVED*** int|sys||ioctl(int fd, u_long com, ... void *data); ***REMOVED***
	SYS_REVOKE               = 56  // ***REMOVED*** int|sys||revoke(const char *path); ***REMOVED***
	SYS_SYMLINK              = 57  // ***REMOVED*** int|sys||symlink(const char *path, const char *link); ***REMOVED***
	SYS_READLINK             = 58  // ***REMOVED*** ssize_t|sys||readlink(const char *path, char *buf, size_t count); ***REMOVED***
	SYS_EXECVE               = 59  // ***REMOVED*** int|sys||execve(const char *path, char * const *argp, char * const *envp); ***REMOVED***
	SYS_UMASK                = 60  // ***REMOVED*** mode_t|sys||umask(mode_t newmask); ***REMOVED***
	SYS_CHROOT               = 61  // ***REMOVED*** int|sys||chroot(const char *path); ***REMOVED***
	SYS_VFORK                = 66  // ***REMOVED*** int|sys||vfork(void); ***REMOVED***
	SYS_SBRK                 = 69  // ***REMOVED*** int|sys||sbrk(intptr_t incr); ***REMOVED***
	SYS_SSTK                 = 70  // ***REMOVED*** int|sys||sstk(int incr); ***REMOVED***
	SYS_VADVISE              = 72  // ***REMOVED*** int|sys||ovadvise(int anom); ***REMOVED***
	SYS_MUNMAP               = 73  // ***REMOVED*** int|sys||munmap(void *addr, size_t len); ***REMOVED***
	SYS_MPROTECT             = 74  // ***REMOVED*** int|sys||mprotect(void *addr, size_t len, int prot); ***REMOVED***
	SYS_MADVISE              = 75  // ***REMOVED*** int|sys||madvise(void *addr, size_t len, int behav); ***REMOVED***
	SYS_MINCORE              = 78  // ***REMOVED*** int|sys||mincore(void *addr, size_t len, char *vec); ***REMOVED***
	SYS_GETGROUPS            = 79  // ***REMOVED*** int|sys||getgroups(int gidsetsize, gid_t *gidset); ***REMOVED***
	SYS_SETGROUPS            = 80  // ***REMOVED*** int|sys||setgroups(int gidsetsize, const gid_t *gidset); ***REMOVED***
	SYS_GETPGRP              = 81  // ***REMOVED*** int|sys||getpgrp(void); ***REMOVED***
	SYS_SETPGID              = 82  // ***REMOVED*** int|sys||setpgid(pid_t pid, pid_t pgid); ***REMOVED***
	SYS_DUP2                 = 90  // ***REMOVED*** int|sys||dup2(int from, int to); ***REMOVED***
	SYS_FCNTL                = 92  // ***REMOVED*** int|sys||fcntl(int fd, int cmd, ... void *arg); ***REMOVED***
	SYS_FSYNC                = 95  // ***REMOVED*** int|sys||fsync(int fd); ***REMOVED***
	SYS_SETPRIORITY          = 96  // ***REMOVED*** int|sys||setpriority(int which, id_t who, int prio); ***REMOVED***
	SYS_CONNECT              = 98  // ***REMOVED*** int|sys||connect(int s, const struct sockaddr *name, socklen_t namelen); ***REMOVED***
	SYS_GETPRIORITY          = 100 // ***REMOVED*** int|sys||getpriority(int which, id_t who); ***REMOVED***
	SYS_BIND                 = 104 // ***REMOVED*** int|sys||bind(int s, const struct sockaddr *name, socklen_t namelen); ***REMOVED***
	SYS_SETSOCKOPT           = 105 // ***REMOVED*** int|sys||setsockopt(int s, int level, int name, const void *val, socklen_t valsize); ***REMOVED***
	SYS_LISTEN               = 106 // ***REMOVED*** int|sys||listen(int s, int backlog); ***REMOVED***
	SYS_GETSOCKOPT           = 118 // ***REMOVED*** int|sys||getsockopt(int s, int level, int name, void *val, socklen_t *avalsize); ***REMOVED***
	SYS_READV                = 120 // ***REMOVED*** ssize_t|sys||readv(int fd, const struct iovec *iovp, int iovcnt); ***REMOVED***
	SYS_WRITEV               = 121 // ***REMOVED*** ssize_t|sys||writev(int fd, const struct iovec *iovp, int iovcnt); ***REMOVED***
	SYS_FCHOWN               = 123 // ***REMOVED*** int|sys||fchown(int fd, uid_t uid, gid_t gid); ***REMOVED***
	SYS_FCHMOD               = 124 // ***REMOVED*** int|sys||fchmod(int fd, mode_t mode); ***REMOVED***
	SYS_SETREUID             = 126 // ***REMOVED*** int|sys||setreuid(uid_t ruid, uid_t euid); ***REMOVED***
	SYS_SETREGID             = 127 // ***REMOVED*** int|sys||setregid(gid_t rgid, gid_t egid); ***REMOVED***
	SYS_RENAME               = 128 // ***REMOVED*** int|sys||rename(const char *from, const char *to); ***REMOVED***
	SYS_FLOCK                = 131 // ***REMOVED*** int|sys||flock(int fd, int how); ***REMOVED***
	SYS_MKFIFO               = 132 // ***REMOVED*** int|sys||mkfifo(const char *path, mode_t mode); ***REMOVED***
	SYS_SENDTO               = 133 // ***REMOVED*** ssize_t|sys||sendto(int s, const void *buf, size_t len, int flags, const struct sockaddr *to, socklen_t tolen); ***REMOVED***
	SYS_SHUTDOWN             = 134 // ***REMOVED*** int|sys||shutdown(int s, int how); ***REMOVED***
	SYS_SOCKETPAIR           = 135 // ***REMOVED*** int|sys||socketpair(int domain, int type, int protocol, int *rsv); ***REMOVED***
	SYS_MKDIR                = 136 // ***REMOVED*** int|sys||mkdir(const char *path, mode_t mode); ***REMOVED***
	SYS_RMDIR                = 137 // ***REMOVED*** int|sys||rmdir(const char *path); ***REMOVED***
	SYS_SETSID               = 147 // ***REMOVED*** int|sys||setsid(void); ***REMOVED***
	SYS_SYSARCH              = 165 // ***REMOVED*** int|sys||sysarch(int op, void *parms); ***REMOVED***
	SYS_PREAD                = 173 // ***REMOVED*** ssize_t|sys||pread(int fd, void *buf, size_t nbyte, int PAD, off_t offset); ***REMOVED***
	SYS_PWRITE               = 174 // ***REMOVED*** ssize_t|sys||pwrite(int fd, const void *buf, size_t nbyte, int PAD, off_t offset); ***REMOVED***
	SYS_NTP_ADJTIME          = 176 // ***REMOVED*** int|sys||ntp_adjtime(struct timex *tp); ***REMOVED***
	SYS_SETGID               = 181 // ***REMOVED*** int|sys||setgid(gid_t gid); ***REMOVED***
	SYS_SETEGID              = 182 // ***REMOVED*** int|sys||setegid(gid_t egid); ***REMOVED***
	SYS_SETEUID              = 183 // ***REMOVED*** int|sys||seteuid(uid_t euid); ***REMOVED***
	SYS_PATHCONF             = 191 // ***REMOVED*** long|sys||pathconf(const char *path, int name); ***REMOVED***
	SYS_FPATHCONF            = 192 // ***REMOVED*** long|sys||fpathconf(int fd, int name); ***REMOVED***
	SYS_GETRLIMIT            = 194 // ***REMOVED*** int|sys||getrlimit(int which, struct rlimit *rlp); ***REMOVED***
	SYS_SETRLIMIT            = 195 // ***REMOVED*** int|sys||setrlimit(int which, const struct rlimit *rlp); ***REMOVED***
	SYS_MMAP                 = 197 // ***REMOVED*** void *|sys||mmap(void *addr, size_t len, int prot, int flags, int fd, long PAD, off_t pos); ***REMOVED***
	SYS_LSEEK                = 199 // ***REMOVED*** off_t|sys||lseek(int fd, int PAD, off_t offset, int whence); ***REMOVED***
	SYS_TRUNCATE             = 200 // ***REMOVED*** int|sys||truncate(const char *path, int PAD, off_t length); ***REMOVED***
	SYS_FTRUNCATE            = 201 // ***REMOVED*** int|sys||ftruncate(int fd, int PAD, off_t length); ***REMOVED***
	SYS___SYSCTL             = 202 // ***REMOVED*** int|sys||__sysctl(const int *name, u_int namelen, void *old, size_t *oldlenp, const void *new, size_t newlen); ***REMOVED***
	SYS_MLOCK                = 203 // ***REMOVED*** int|sys||mlock(const void *addr, size_t len); ***REMOVED***
	SYS_MUNLOCK              = 204 // ***REMOVED*** int|sys||munlock(const void *addr, size_t len); ***REMOVED***
	SYS_UNDELETE             = 205 // ***REMOVED*** int|sys||undelete(const char *path); ***REMOVED***
	SYS_GETPGID              = 207 // ***REMOVED*** pid_t|sys||getpgid(pid_t pid); ***REMOVED***
	SYS_REBOOT               = 208 // ***REMOVED*** int|sys||reboot(int opt, char *bootstr); ***REMOVED***
	SYS_POLL                 = 209 // ***REMOVED*** int|sys||poll(struct pollfd *fds, u_int nfds, int timeout); ***REMOVED***
	SYS_SEMGET               = 221 // ***REMOVED*** int|sys||semget(key_t key, int nsems, int semflg); ***REMOVED***
	SYS_SEMOP                = 222 // ***REMOVED*** int|sys||semop(int semid, struct sembuf *sops, size_t nsops); ***REMOVED***
	SYS_SEMCONFIG            = 223 // ***REMOVED*** int|sys||semconfig(int flag); ***REMOVED***
	SYS_MSGGET               = 225 // ***REMOVED*** int|sys||msgget(key_t key, int msgflg); ***REMOVED***
	SYS_MSGSND               = 226 // ***REMOVED*** int|sys||msgsnd(int msqid, const void *msgp, size_t msgsz, int msgflg); ***REMOVED***
	SYS_MSGRCV               = 227 // ***REMOVED*** ssize_t|sys||msgrcv(int msqid, void *msgp, size_t msgsz, long msgtyp, int msgflg); ***REMOVED***
	SYS_SHMAT                = 228 // ***REMOVED*** void *|sys||shmat(int shmid, const void *shmaddr, int shmflg); ***REMOVED***
	SYS_SHMDT                = 230 // ***REMOVED*** int|sys||shmdt(const void *shmaddr); ***REMOVED***
	SYS_SHMGET               = 231 // ***REMOVED*** int|sys||shmget(key_t key, size_t size, int shmflg); ***REMOVED***
	SYS_TIMER_CREATE         = 235 // ***REMOVED*** int|sys||timer_create(clockid_t clock_id, struct sigevent *evp, timer_t *timerid); ***REMOVED***
	SYS_TIMER_DELETE         = 236 // ***REMOVED*** int|sys||timer_delete(timer_t timerid); ***REMOVED***
	SYS_TIMER_GETOVERRUN     = 239 // ***REMOVED*** int|sys||timer_getoverrun(timer_t timerid); ***REMOVED***
	SYS_FDATASYNC            = 241 // ***REMOVED*** int|sys||fdatasync(int fd); ***REMOVED***
	SYS_MLOCKALL             = 242 // ***REMOVED*** int|sys||mlockall(int flags); ***REMOVED***
	SYS_MUNLOCKALL           = 243 // ***REMOVED*** int|sys||munlockall(void); ***REMOVED***
	SYS_SIGQUEUEINFO         = 245 // ***REMOVED*** int|sys||sigqueueinfo(pid_t pid, const siginfo_t *info); ***REMOVED***
	SYS_MODCTL               = 246 // ***REMOVED*** int|sys||modctl(int cmd, void *arg); ***REMOVED***
	SYS___POSIX_RENAME       = 270 // ***REMOVED*** int|sys||__posix_rename(const char *from, const char *to); ***REMOVED***
	SYS_SWAPCTL              = 271 // ***REMOVED*** int|sys||swapctl(int cmd, void *arg, int misc); ***REMOVED***
	SYS_MINHERIT             = 273 // ***REMOVED*** int|sys||minherit(void *addr, size_t len, int inherit); ***REMOVED***
	SYS_LCHMOD               = 274 // ***REMOVED*** int|sys||lchmod(const char *path, mode_t mode); ***REMOVED***
	SYS_LCHOWN               = 275 // ***REMOVED*** int|sys||lchown(const char *path, uid_t uid, gid_t gid); ***REMOVED***
	SYS_MSYNC                = 277 // ***REMOVED*** int|sys|13|msync(void *addr, size_t len, int flags); ***REMOVED***
	SYS___POSIX_CHOWN        = 283 // ***REMOVED*** int|sys||__posix_chown(const char *path, uid_t uid, gid_t gid); ***REMOVED***
	SYS___POSIX_FCHOWN       = 284 // ***REMOVED*** int|sys||__posix_fchown(int fd, uid_t uid, gid_t gid); ***REMOVED***
	SYS___POSIX_LCHOWN       = 285 // ***REMOVED*** int|sys||__posix_lchown(const char *path, uid_t uid, gid_t gid); ***REMOVED***
	SYS_GETSID               = 286 // ***REMOVED*** pid_t|sys||getsid(pid_t pid); ***REMOVED***
	SYS___CLONE              = 287 // ***REMOVED*** pid_t|sys||__clone(int flags, void *stack); ***REMOVED***
	SYS_FKTRACE              = 288 // ***REMOVED*** int|sys||fktrace(int fd, int ops, int facs, pid_t pid); ***REMOVED***
	SYS_PREADV               = 289 // ***REMOVED*** ssize_t|sys||preadv(int fd, const struct iovec *iovp, int iovcnt, int PAD, off_t offset); ***REMOVED***
	SYS_PWRITEV              = 290 // ***REMOVED*** ssize_t|sys||pwritev(int fd, const struct iovec *iovp, int iovcnt, int PAD, off_t offset); ***REMOVED***
	SYS___GETCWD             = 296 // ***REMOVED*** int|sys||__getcwd(char *bufp, size_t length); ***REMOVED***
	SYS_FCHROOT              = 297 // ***REMOVED*** int|sys||fchroot(int fd); ***REMOVED***
	SYS_LCHFLAGS             = 304 // ***REMOVED*** int|sys||lchflags(const char *path, u_long flags); ***REMOVED***
	SYS_ISSETUGID            = 305 // ***REMOVED*** int|sys||issetugid(void); ***REMOVED***
	SYS_UTRACE               = 306 // ***REMOVED*** int|sys||utrace(const char *label, void *addr, size_t len); ***REMOVED***
	SYS_GETCONTEXT           = 307 // ***REMOVED*** int|sys||getcontext(struct __ucontext *ucp); ***REMOVED***
	SYS_SETCONTEXT           = 308 // ***REMOVED*** int|sys||setcontext(const struct __ucontext *ucp); ***REMOVED***
	SYS__LWP_CREATE          = 309 // ***REMOVED*** int|sys||_lwp_create(const struct __ucontext *ucp, u_long flags, lwpid_t *new_lwp); ***REMOVED***
	SYS__LWP_EXIT            = 310 // ***REMOVED*** int|sys||_lwp_exit(void); ***REMOVED***
	SYS__LWP_SELF            = 311 // ***REMOVED*** lwpid_t|sys||_lwp_self(void); ***REMOVED***
	SYS__LWP_WAIT            = 312 // ***REMOVED*** int|sys||_lwp_wait(lwpid_t wait_for, lwpid_t *departed); ***REMOVED***
	SYS__LWP_SUSPEND         = 313 // ***REMOVED*** int|sys||_lwp_suspend(lwpid_t target); ***REMOVED***
	SYS__LWP_CONTINUE        = 314 // ***REMOVED*** int|sys||_lwp_continue(lwpid_t target); ***REMOVED***
	SYS__LWP_WAKEUP          = 315 // ***REMOVED*** int|sys||_lwp_wakeup(lwpid_t target); ***REMOVED***
	SYS__LWP_GETPRIVATE      = 316 // ***REMOVED*** void *|sys||_lwp_getprivate(void); ***REMOVED***
	SYS__LWP_SETPRIVATE      = 317 // ***REMOVED*** void|sys||_lwp_setprivate(void *ptr); ***REMOVED***
	SYS__LWP_KILL            = 318 // ***REMOVED*** int|sys||_lwp_kill(lwpid_t target, int signo); ***REMOVED***
	SYS__LWP_DETACH          = 319 // ***REMOVED*** int|sys||_lwp_detach(lwpid_t target); ***REMOVED***
	SYS__LWP_UNPARK          = 321 // ***REMOVED*** int|sys||_lwp_unpark(lwpid_t target, const void *hint); ***REMOVED***
	SYS__LWP_UNPARK_ALL      = 322 // ***REMOVED*** ssize_t|sys||_lwp_unpark_all(const lwpid_t *targets, size_t ntargets, const void *hint); ***REMOVED***
	SYS__LWP_SETNAME         = 323 // ***REMOVED*** int|sys||_lwp_setname(lwpid_t target, const char *name); ***REMOVED***
	SYS__LWP_GETNAME         = 324 // ***REMOVED*** int|sys||_lwp_getname(lwpid_t target, char *name, size_t len); ***REMOVED***
	SYS__LWP_CTL             = 325 // ***REMOVED*** int|sys||_lwp_ctl(int features, struct lwpctl **address); ***REMOVED***
	SYS___SIGACTION_SIGTRAMP = 340 // ***REMOVED*** int|sys||__sigaction_sigtramp(int signum, const struct sigaction *nsa, struct sigaction *osa, const void *tramp, int vers); ***REMOVED***
	SYS_PMC_GET_INFO         = 341 // ***REMOVED*** int|sys||pmc_get_info(int ctr, int op, void *args); ***REMOVED***
	SYS_PMC_CONTROL          = 342 // ***REMOVED*** int|sys||pmc_control(int ctr, int op, void *args); ***REMOVED***
	SYS_RASCTL               = 343 // ***REMOVED*** int|sys||rasctl(void *addr, size_t len, int op); ***REMOVED***
	SYS_KQUEUE               = 344 // ***REMOVED*** int|sys||kqueue(void); ***REMOVED***
	SYS__SCHED_SETPARAM      = 346 // ***REMOVED*** int|sys||_sched_setparam(pid_t pid, lwpid_t lid, int policy, const struct sched_param *params); ***REMOVED***
	SYS__SCHED_GETPARAM      = 347 // ***REMOVED*** int|sys||_sched_getparam(pid_t pid, lwpid_t lid, int *policy, struct sched_param *params); ***REMOVED***
	SYS__SCHED_SETAFFINITY   = 348 // ***REMOVED*** int|sys||_sched_setaffinity(pid_t pid, lwpid_t lid, size_t size, const cpuset_t *cpuset); ***REMOVED***
	SYS__SCHED_GETAFFINITY   = 349 // ***REMOVED*** int|sys||_sched_getaffinity(pid_t pid, lwpid_t lid, size_t size, cpuset_t *cpuset); ***REMOVED***
	SYS_SCHED_YIELD          = 350 // ***REMOVED*** int|sys||sched_yield(void); ***REMOVED***
	SYS_FSYNC_RANGE          = 354 // ***REMOVED*** int|sys||fsync_range(int fd, int flags, off_t start, off_t length); ***REMOVED***
	SYS_UUIDGEN              = 355 // ***REMOVED*** int|sys||uuidgen(struct uuid *store, int count); ***REMOVED***
	SYS_GETVFSSTAT           = 356 // ***REMOVED*** int|sys||getvfsstat(struct statvfs *buf, size_t bufsize, int flags); ***REMOVED***
	SYS_STATVFS1             = 357 // ***REMOVED*** int|sys||statvfs1(const char *path, struct statvfs *buf, int flags); ***REMOVED***
	SYS_FSTATVFS1            = 358 // ***REMOVED*** int|sys||fstatvfs1(int fd, struct statvfs *buf, int flags); ***REMOVED***
	SYS_EXTATTRCTL           = 360 // ***REMOVED*** int|sys||extattrctl(const char *path, int cmd, const char *filename, int attrnamespace, const char *attrname); ***REMOVED***
	SYS_EXTATTR_SET_FILE     = 361 // ***REMOVED*** int|sys||extattr_set_file(const char *path, int attrnamespace, const char *attrname, const void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_GET_FILE     = 362 // ***REMOVED*** ssize_t|sys||extattr_get_file(const char *path, int attrnamespace, const char *attrname, void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_DELETE_FILE  = 363 // ***REMOVED*** int|sys||extattr_delete_file(const char *path, int attrnamespace, const char *attrname); ***REMOVED***
	SYS_EXTATTR_SET_FD       = 364 // ***REMOVED*** int|sys||extattr_set_fd(int fd, int attrnamespace, const char *attrname, const void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_GET_FD       = 365 // ***REMOVED*** ssize_t|sys||extattr_get_fd(int fd, int attrnamespace, const char *attrname, void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_DELETE_FD    = 366 // ***REMOVED*** int|sys||extattr_delete_fd(int fd, int attrnamespace, const char *attrname); ***REMOVED***
	SYS_EXTATTR_SET_LINK     = 367 // ***REMOVED*** int|sys||extattr_set_link(const char *path, int attrnamespace, const char *attrname, const void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_GET_LINK     = 368 // ***REMOVED*** ssize_t|sys||extattr_get_link(const char *path, int attrnamespace, const char *attrname, void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_DELETE_LINK  = 369 // ***REMOVED*** int|sys||extattr_delete_link(const char *path, int attrnamespace, const char *attrname); ***REMOVED***
	SYS_EXTATTR_LIST_FD      = 370 // ***REMOVED*** ssize_t|sys||extattr_list_fd(int fd, int attrnamespace, void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_LIST_FILE    = 371 // ***REMOVED*** ssize_t|sys||extattr_list_file(const char *path, int attrnamespace, void *data, size_t nbytes); ***REMOVED***
	SYS_EXTATTR_LIST_LINK    = 372 // ***REMOVED*** ssize_t|sys||extattr_list_link(const char *path, int attrnamespace, void *data, size_t nbytes); ***REMOVED***
	SYS_SETXATTR             = 375 // ***REMOVED*** int|sys||setxattr(const char *path, const char *name, const void *value, size_t size, int flags); ***REMOVED***
	SYS_LSETXATTR            = 376 // ***REMOVED*** int|sys||lsetxattr(const char *path, const char *name, const void *value, size_t size, int flags); ***REMOVED***
	SYS_FSETXATTR            = 377 // ***REMOVED*** int|sys||fsetxattr(int fd, const char *name, const void *value, size_t size, int flags); ***REMOVED***
	SYS_GETXATTR             = 378 // ***REMOVED*** int|sys||getxattr(const char *path, const char *name, void *value, size_t size); ***REMOVED***
	SYS_LGETXATTR            = 379 // ***REMOVED*** int|sys||lgetxattr(const char *path, const char *name, void *value, size_t size); ***REMOVED***
	SYS_FGETXATTR            = 380 // ***REMOVED*** int|sys||fgetxattr(int fd, const char *name, void *value, size_t size); ***REMOVED***
	SYS_LISTXATTR            = 381 // ***REMOVED*** int|sys||listxattr(const char *path, char *list, size_t size); ***REMOVED***
	SYS_LLISTXATTR           = 382 // ***REMOVED*** int|sys||llistxattr(const char *path, char *list, size_t size); ***REMOVED***
	SYS_FLISTXATTR           = 383 // ***REMOVED*** int|sys||flistxattr(int fd, char *list, size_t size); ***REMOVED***
	SYS_REMOVEXATTR          = 384 // ***REMOVED*** int|sys||removexattr(const char *path, const char *name); ***REMOVED***
	SYS_LREMOVEXATTR         = 385 // ***REMOVED*** int|sys||lremovexattr(const char *path, const char *name); ***REMOVED***
	SYS_FREMOVEXATTR         = 386 // ***REMOVED*** int|sys||fremovexattr(int fd, const char *name); ***REMOVED***
	SYS_GETDENTS             = 390 // ***REMOVED*** int|sys|30|getdents(int fd, char *buf, size_t count); ***REMOVED***
	SYS_SOCKET               = 394 // ***REMOVED*** int|sys|30|socket(int domain, int type, int protocol); ***REMOVED***
	SYS_GETFH                = 395 // ***REMOVED*** int|sys|30|getfh(const char *fname, void *fhp, size_t *fh_size); ***REMOVED***
	SYS_MOUNT                = 410 // ***REMOVED*** int|sys|50|mount(const char *type, const char *path, int flags, void *data, size_t data_len); ***REMOVED***
	SYS_MREMAP               = 411 // ***REMOVED*** void *|sys||mremap(void *old_address, size_t old_size, void *new_address, size_t new_size, int flags); ***REMOVED***
	SYS_PSET_CREATE          = 412 // ***REMOVED*** int|sys||pset_create(psetid_t *psid); ***REMOVED***
	SYS_PSET_DESTROY         = 413 // ***REMOVED*** int|sys||pset_destroy(psetid_t psid); ***REMOVED***
	SYS_PSET_ASSIGN          = 414 // ***REMOVED*** int|sys||pset_assign(psetid_t psid, cpuid_t cpuid, psetid_t *opsid); ***REMOVED***
	SYS__PSET_BIND           = 415 // ***REMOVED*** int|sys||_pset_bind(idtype_t idtype, id_t first_id, id_t second_id, psetid_t psid, psetid_t *opsid); ***REMOVED***
	SYS_POSIX_FADVISE        = 416 // ***REMOVED*** int|sys|50|posix_fadvise(int fd, int PAD, off_t offset, off_t len, int advice); ***REMOVED***
	SYS_SELECT               = 417 // ***REMOVED*** int|sys|50|select(int nd, fd_set *in, fd_set *ou, fd_set *ex, struct timeval *tv); ***REMOVED***
	SYS_GETTIMEOFDAY         = 418 // ***REMOVED*** int|sys|50|gettimeofday(struct timeval *tp, void *tzp); ***REMOVED***
	SYS_SETTIMEOFDAY         = 419 // ***REMOVED*** int|sys|50|settimeofday(const struct timeval *tv, const void *tzp); ***REMOVED***
	SYS_UTIMES               = 420 // ***REMOVED*** int|sys|50|utimes(const char *path, const struct timeval *tptr); ***REMOVED***
	SYS_ADJTIME              = 421 // ***REMOVED*** int|sys|50|adjtime(const struct timeval *delta, struct timeval *olddelta); ***REMOVED***
	SYS_FUTIMES              = 423 // ***REMOVED*** int|sys|50|futimes(int fd, const struct timeval *tptr); ***REMOVED***
	SYS_LUTIMES              = 424 // ***REMOVED*** int|sys|50|lutimes(const char *path, const struct timeval *tptr); ***REMOVED***
	SYS_SETITIMER            = 425 // ***REMOVED*** int|sys|50|setitimer(int which, const struct itimerval *itv, struct itimerval *oitv); ***REMOVED***
	SYS_GETITIMER            = 426 // ***REMOVED*** int|sys|50|getitimer(int which, struct itimerval *itv); ***REMOVED***
	SYS_CLOCK_GETTIME        = 427 // ***REMOVED*** int|sys|50|clock_gettime(clockid_t clock_id, struct timespec *tp); ***REMOVED***
	SYS_CLOCK_SETTIME        = 428 // ***REMOVED*** int|sys|50|clock_settime(clockid_t clock_id, const struct timespec *tp); ***REMOVED***
	SYS_CLOCK_GETRES         = 429 // ***REMOVED*** int|sys|50|clock_getres(clockid_t clock_id, struct timespec *tp); ***REMOVED***
	SYS_NANOSLEEP            = 430 // ***REMOVED*** int|sys|50|nanosleep(const struct timespec *rqtp, struct timespec *rmtp); ***REMOVED***
	SYS___SIGTIMEDWAIT       = 431 // ***REMOVED*** int|sys|50|__sigtimedwait(const sigset_t *set, siginfo_t *info, struct timespec *timeout); ***REMOVED***
	SYS__LWP_PARK            = 434 // ***REMOVED*** int|sys|50|_lwp_park(const struct timespec *ts, lwpid_t unpark, const void *hint, const void *unparkhint); ***REMOVED***
	SYS_KEVENT               = 435 // ***REMOVED*** int|sys|50|kevent(int fd, const struct kevent *changelist, size_t nchanges, struct kevent *eventlist, size_t nevents, const struct timespec *timeout); ***REMOVED***
	SYS_PSELECT              = 436 // ***REMOVED*** int|sys|50|pselect(int nd, fd_set *in, fd_set *ou, fd_set *ex, const struct timespec *ts, const sigset_t *mask); ***REMOVED***
	SYS_POLLTS               = 437 // ***REMOVED*** int|sys|50|pollts(struct pollfd *fds, u_int nfds, const struct timespec *ts, const sigset_t *mask); ***REMOVED***
	SYS_STAT                 = 439 // ***REMOVED*** int|sys|50|stat(const char *path, struct stat *ub); ***REMOVED***
	SYS_FSTAT                = 440 // ***REMOVED*** int|sys|50|fstat(int fd, struct stat *sb); ***REMOVED***
	SYS_LSTAT                = 441 // ***REMOVED*** int|sys|50|lstat(const char *path, struct stat *ub); ***REMOVED***
	SYS___SEMCTL             = 442 // ***REMOVED*** int|sys|50|__semctl(int semid, int semnum, int cmd, ... union __semun *arg); ***REMOVED***
	SYS_SHMCTL               = 443 // ***REMOVED*** int|sys|50|shmctl(int shmid, int cmd, struct shmid_ds *buf); ***REMOVED***
	SYS_MSGCTL               = 444 // ***REMOVED*** int|sys|50|msgctl(int msqid, int cmd, struct msqid_ds *buf); ***REMOVED***
	SYS_GETRUSAGE            = 445 // ***REMOVED*** int|sys|50|getrusage(int who, struct rusage *rusage); ***REMOVED***
	SYS_TIMER_SETTIME        = 446 // ***REMOVED*** int|sys|50|timer_settime(timer_t timerid, int flags, const struct itimerspec *value, struct itimerspec *ovalue); ***REMOVED***
	SYS_TIMER_GETTIME        = 447 // ***REMOVED*** int|sys|50|timer_gettime(timer_t timerid, struct itimerspec *value); ***REMOVED***
	SYS_NTP_GETTIME          = 448 // ***REMOVED*** int|sys|50|ntp_gettime(struct ntptimeval *ntvp); ***REMOVED***
	SYS_WAIT4                = 449 // ***REMOVED*** int|sys|50|wait4(pid_t pid, int *status, int options, struct rusage *rusage); ***REMOVED***
	SYS_MKNOD                = 450 // ***REMOVED*** int|sys|50|mknod(const char *path, mode_t mode, dev_t dev); ***REMOVED***
	SYS_FHSTAT               = 451 // ***REMOVED*** int|sys|50|fhstat(const void *fhp, size_t fh_size, struct stat *sb); ***REMOVED***
	SYS_PIPE2                = 453 // ***REMOVED*** int|sys||pipe2(int *fildes, int flags); ***REMOVED***
	SYS_DUP3                 = 454 // ***REMOVED*** int|sys||dup3(int from, int to, int flags); ***REMOVED***
	SYS_KQUEUE1              = 455 // ***REMOVED*** int|sys||kqueue1(int flags); ***REMOVED***
	SYS_PACCEPT              = 456 // ***REMOVED*** int|sys||paccept(int s, struct sockaddr *name, socklen_t *anamelen, const sigset_t *mask, int flags); ***REMOVED***
	SYS_LINKAT               = 457 // ***REMOVED*** int|sys||linkat(int fd1, const char *name1, int fd2, const char *name2, int flags); ***REMOVED***
	SYS_RENAMEAT             = 458 // ***REMOVED*** int|sys||renameat(int fromfd, const char *from, int tofd, const char *to); ***REMOVED***
	SYS_MKFIFOAT             = 459 // ***REMOVED*** int|sys||mkfifoat(int fd, const char *path, mode_t mode); ***REMOVED***
	SYS_MKNODAT              = 460 // ***REMOVED*** int|sys||mknodat(int fd, const char *path, mode_t mode, uint32_t dev); ***REMOVED***
	SYS_MKDIRAT              = 461 // ***REMOVED*** int|sys||mkdirat(int fd, const char *path, mode_t mode); ***REMOVED***
	SYS_FACCESSAT            = 462 // ***REMOVED*** int|sys||faccessat(int fd, const char *path, int amode, int flag); ***REMOVED***
	SYS_FCHMODAT             = 463 // ***REMOVED*** int|sys||fchmodat(int fd, const char *path, mode_t mode, int flag); ***REMOVED***
	SYS_FCHOWNAT             = 464 // ***REMOVED*** int|sys||fchownat(int fd, const char *path, uid_t owner, gid_t group, int flag); ***REMOVED***
	SYS_FEXECVE              = 465 // ***REMOVED*** int|sys||fexecve(int fd, char * const *argp, char * const *envp); ***REMOVED***
	SYS_FSTATAT              = 466 // ***REMOVED*** int|sys||fstatat(int fd, const char *path, struct stat *buf, int flag); ***REMOVED***
	SYS_UTIMENSAT            = 467 // ***REMOVED*** int|sys||utimensat(int fd, const char *path, const struct timespec *tptr, int flag); ***REMOVED***
	SYS_OPENAT               = 468 // ***REMOVED*** int|sys||openat(int fd, const char *path, int oflags, ... mode_t mode); ***REMOVED***
	SYS_READLINKAT           = 469 // ***REMOVED*** int|sys||readlinkat(int fd, const char *path, char *buf, size_t bufsize); ***REMOVED***
	SYS_SYMLINKAT            = 470 // ***REMOVED*** int|sys||symlinkat(const char *path1, int fd, const char *path2); ***REMOVED***
	SYS_UNLINKAT             = 471 // ***REMOVED*** int|sys||unlinkat(int fd, const char *path, int flag); ***REMOVED***
	SYS_FUTIMENS             = 472 // ***REMOVED*** int|sys||futimens(int fd, const struct timespec *tptr); ***REMOVED***
	SYS___QUOTACTL           = 473 // ***REMOVED*** int|sys||__quotactl(const char *path, struct quotactl_args *args); ***REMOVED***
	SYS_POSIX_SPAWN          = 474 // ***REMOVED*** int|sys||posix_spawn(pid_t *pid, const char *path, const struct posix_spawn_file_actions *file_actions, const struct posix_spawnattr *attrp, char *const *argv, char *const *envp); ***REMOVED***
	SYS_RECVMMSG             = 475 // ***REMOVED*** int|sys||recvmmsg(int s, struct mmsghdr *mmsg, unsigned int vlen, unsigned int flags, struct timespec *timeout); ***REMOVED***
	SYS_SENDMMSG             = 476 // ***REMOVED*** int|sys||sendmmsg(int s, struct mmsghdr *mmsg, unsigned int vlen, unsigned int flags); ***REMOVED***
)
