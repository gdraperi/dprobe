package ns

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

var (
	initNs   netns.NsHandle
	initNl   *netlink.Handle
	initOnce sync.Once
	// NetlinkSocketsTimeout represents the default timeout duration for the sockets
	NetlinkSocketsTimeout = 3 * time.Second
)

// Init initializes a new network namespace
func Init() ***REMOVED***
	var err error
	initNs, err = netns.Get()
	if err != nil ***REMOVED***
		logrus.Errorf("could not get initial namespace: %v", err)
	***REMOVED***
	initNl, err = netlink.NewHandle(getSupportedNlFamilies()...)
	if err != nil ***REMOVED***
		logrus.Errorf("could not create netlink handle on initial namespace: %v", err)
	***REMOVED***
	err = initNl.SetSocketTimeout(NetlinkSocketsTimeout)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to set the timeout on the default netlink handle sockets: %v", err)
	***REMOVED***
***REMOVED***

// SetNamespace sets the initial namespace handler
func SetNamespace() error ***REMOVED***
	initOnce.Do(Init)
	if err := netns.Set(initNs); err != nil ***REMOVED***
		linkInfo, linkErr := getLink()
		if linkErr != nil ***REMOVED***
			linkInfo = linkErr.Error()
		***REMOVED***
		return fmt.Errorf("failed to set to initial namespace, %v, initns fd %d: %v", linkInfo, initNs, err)
	***REMOVED***
	return nil
***REMOVED***

// ParseHandlerInt transforms the namespace handler into an integer
func ParseHandlerInt() int ***REMOVED***
	return int(getHandler())
***REMOVED***

// GetHandler returns the namespace handler
func getHandler() netns.NsHandle ***REMOVED***
	initOnce.Do(Init)
	return initNs
***REMOVED***

func getLink() (string, error) ***REMOVED***
	return os.Readlink(fmt.Sprintf("/proc/%d/task/%d/ns/net", os.Getpid(), syscall.Gettid()))
***REMOVED***

// NlHandle returns the netlink handler
func NlHandle() *netlink.Handle ***REMOVED***
	initOnce.Do(Init)
	return initNl
***REMOVED***

func getSupportedNlFamilies() []int ***REMOVED***
	fams := []int***REMOVED***syscall.NETLINK_ROUTE***REMOVED***
	// NETLINK_XFRM test
	if err := loadXfrmModules(); err != nil ***REMOVED***
		if checkXfrmSocket() != nil ***REMOVED***
			logrus.Warnf("Could not load necessary modules for IPSEC rules: %v", err)
		***REMOVED*** else ***REMOVED***
			fams = append(fams, syscall.NETLINK_XFRM)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fams = append(fams, syscall.NETLINK_XFRM)
	***REMOVED***
	// NETLINK_NETFILTER test
	if err := loadNfConntrackModules(); err != nil ***REMOVED***
		if checkNfSocket() != nil ***REMOVED***
			logrus.Warnf("Could not load necessary modules for Conntrack: %v", err)
		***REMOVED*** else ***REMOVED***
			fams = append(fams, syscall.NETLINK_NETFILTER)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fams = append(fams, syscall.NETLINK_NETFILTER)
	***REMOVED***

	return fams
***REMOVED***

func loadXfrmModules() error ***REMOVED***
	if out, err := exec.Command("modprobe", "-va", "xfrm_user").CombinedOutput(); err != nil ***REMOVED***
		return fmt.Errorf("Running modprobe xfrm_user failed with message: `%s`, error: %v", strings.TrimSpace(string(out)), err)
	***REMOVED***
	if out, err := exec.Command("modprobe", "-va", "xfrm_algo").CombinedOutput(); err != nil ***REMOVED***
		return fmt.Errorf("Running modprobe xfrm_algo failed with message: `%s`, error: %v", strings.TrimSpace(string(out)), err)
	***REMOVED***
	return nil
***REMOVED***

// API check on required xfrm modules (xfrm_user, xfrm_algo)
func checkXfrmSocket() error ***REMOVED***
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_XFRM)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	syscall.Close(fd)
	return nil
***REMOVED***

func loadNfConntrackModules() error ***REMOVED***
	if out, err := exec.Command("modprobe", "-va", "nf_conntrack").CombinedOutput(); err != nil ***REMOVED***
		return fmt.Errorf("Running modprobe nf_conntrack failed with message: `%s`, error: %v", strings.TrimSpace(string(out)), err)
	***REMOVED***
	if out, err := exec.Command("modprobe", "-va", "nf_conntrack_netlink").CombinedOutput(); err != nil ***REMOVED***
		return fmt.Errorf("Running modprobe nf_conntrack_netlink failed with message: `%s`, error: %v", strings.TrimSpace(string(out)), err)
	***REMOVED***
	return nil
***REMOVED***

// API check on required nf_conntrack* modules (nf_conntrack, nf_conntrack_netlink)
func checkNfSocket() error ***REMOVED***
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_NETFILTER)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	syscall.Close(fd)
	return nil
***REMOVED***
