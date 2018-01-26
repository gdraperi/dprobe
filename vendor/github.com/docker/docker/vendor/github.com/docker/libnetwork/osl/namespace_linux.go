package osl

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

const defaultPrefix = "/var/run/docker"

func init() ***REMOVED***
	reexec.Register("set-ipv6", reexecSetIPv6)
***REMOVED***

var (
	once             sync.Once
	garbagePathMap   = make(map[string]bool)
	gpmLock          sync.Mutex
	gpmWg            sync.WaitGroup
	gpmCleanupPeriod = 60 * time.Second
	gpmChan          = make(chan chan struct***REMOVED******REMOVED***)
	prefix           = defaultPrefix
)

// The networkNamespace type is the linux implementation of the Sandbox
// interface. It represents a linux network namespace, and moves an interface
// into it when called on method AddInterface or sets the gateway etc.
type networkNamespace struct ***REMOVED***
	path         string
	iFaces       []*nwIface
	gw           net.IP
	gwv6         net.IP
	staticRoutes []*types.StaticRoute
	neighbors    []*neigh
	nextIfIndex  map[string]int
	isDefault    bool
	nlHandle     *netlink.Handle
	loV6Enabled  bool
	sync.Mutex
***REMOVED***

// SetBasePath sets the base url prefix for the ns path
func SetBasePath(path string) ***REMOVED***
	prefix = path
***REMOVED***

func init() ***REMOVED***
	reexec.Register("netns-create", reexecCreateNamespace)
***REMOVED***

func basePath() string ***REMOVED***
	return filepath.Join(prefix, "netns")
***REMOVED***

func createBasePath() ***REMOVED***
	err := os.MkdirAll(basePath(), 0755)
	if err != nil ***REMOVED***
		panic("Could not create net namespace path directory")
	***REMOVED***

	// Start the garbage collection go routine
	go removeUnusedPaths()
***REMOVED***

func removeUnusedPaths() ***REMOVED***
	gpmLock.Lock()
	period := gpmCleanupPeriod
	gpmLock.Unlock()

	ticker := time.NewTicker(period)
	for ***REMOVED***
		var (
			gc   chan struct***REMOVED******REMOVED***
			gcOk bool
		)

		select ***REMOVED***
		case <-ticker.C:
		case gc, gcOk = <-gpmChan:
		***REMOVED***

		gpmLock.Lock()
		pathList := make([]string, 0, len(garbagePathMap))
		for path := range garbagePathMap ***REMOVED***
			pathList = append(pathList, path)
		***REMOVED***
		garbagePathMap = make(map[string]bool)
		gpmWg.Add(1)
		gpmLock.Unlock()

		for _, path := range pathList ***REMOVED***
			os.Remove(path)
		***REMOVED***

		gpmWg.Done()
		if gcOk ***REMOVED***
			close(gc)
		***REMOVED***
	***REMOVED***
***REMOVED***

func addToGarbagePaths(path string) ***REMOVED***
	gpmLock.Lock()
	garbagePathMap[path] = true
	gpmLock.Unlock()
***REMOVED***

func removeFromGarbagePaths(path string) ***REMOVED***
	gpmLock.Lock()
	delete(garbagePathMap, path)
	gpmLock.Unlock()
***REMOVED***

// GC triggers garbage collection of namespace path right away
// and waits for it.
func GC() ***REMOVED***
	gpmLock.Lock()
	if len(garbagePathMap) == 0 ***REMOVED***
		// No need for GC if map is empty
		gpmLock.Unlock()
		return
	***REMOVED***
	gpmLock.Unlock()

	// if content exists in the garbage paths
	// we can trigger GC to run, providing a
	// channel to be notified on completion
	waitGC := make(chan struct***REMOVED******REMOVED***)
	gpmChan <- waitGC
	// wait for GC completion
	<-waitGC
***REMOVED***

// GenerateKey generates a sandbox key based on the passed
// container id.
func GenerateKey(containerID string) string ***REMOVED***
	maxLen := 12
	// Read sandbox key from host for overlay
	if strings.HasPrefix(containerID, "-") ***REMOVED***
		var (
			index    int
			indexStr string
			tmpkey   string
		)
		dir, err := ioutil.ReadDir(basePath())
		if err != nil ***REMOVED***
			return ""
		***REMOVED***

		for _, v := range dir ***REMOVED***
			id := v.Name()
			if strings.HasSuffix(id, containerID[:maxLen-1]) ***REMOVED***
				indexStr = strings.TrimSuffix(id, containerID[:maxLen-1])
				tmpindex, err := strconv.Atoi(indexStr)
				if err != nil ***REMOVED***
					return ""
				***REMOVED***
				if tmpindex > index ***REMOVED***
					index = tmpindex
					tmpkey = id
				***REMOVED***

			***REMOVED***
		***REMOVED***
		containerID = tmpkey
		if containerID == "" ***REMOVED***
			return ""
		***REMOVED***
	***REMOVED***

	if len(containerID) < maxLen ***REMOVED***
		maxLen = len(containerID)
	***REMOVED***

	return basePath() + "/" + containerID[:maxLen]
***REMOVED***

// NewSandbox provides a new sandbox instance created in an os specific way
// provided a key which uniquely identifies the sandbox
func NewSandbox(key string, osCreate, isRestore bool) (Sandbox, error) ***REMOVED***
	if !isRestore ***REMOVED***
		err := createNetworkNamespace(key, osCreate)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		once.Do(createBasePath)
	***REMOVED***

	n := &networkNamespace***REMOVED***path: key, isDefault: !osCreate, nextIfIndex: make(map[string]int)***REMOVED***

	sboxNs, err := netns.GetFromPath(n.path)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed get network namespace %q: %v", n.path, err)
	***REMOVED***
	defer sboxNs.Close()

	n.nlHandle, err = netlink.NewHandleAt(sboxNs, syscall.NETLINK_ROUTE)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to create a netlink handle: %v", err)
	***REMOVED***

	err = n.nlHandle.SetSocketTimeout(ns.NetlinkSocketsTimeout)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to set the timeout on the sandbox netlink handle sockets: %v", err)
	***REMOVED***
	// In live-restore mode, IPV6 entries are getting cleaned up due to below code
	// We should retain IPV6 configrations in live-restore mode when Docker Daemon
	// comes back. It should work as it is on other cases
	// As starting point, disable IPv6 on all interfaces
	if !isRestore && !n.isDefault ***REMOVED***
		err = setIPv6(n.path, "all", false)
		if err != nil ***REMOVED***
			logrus.Warnf("Failed to disable IPv6 on all interfaces on network namespace %q: %v", n.path, err)
		***REMOVED***
	***REMOVED***

	if err = n.loopbackUp(); err != nil ***REMOVED***
		n.nlHandle.Delete()
		return nil, err
	***REMOVED***

	return n, nil
***REMOVED***

func (n *networkNamespace) InterfaceOptions() IfaceOptionSetter ***REMOVED***
	return n
***REMOVED***

func (n *networkNamespace) NeighborOptions() NeighborOptionSetter ***REMOVED***
	return n
***REMOVED***

func mountNetworkNamespace(basePath string, lnPath string) error ***REMOVED***
	return syscall.Mount(basePath, lnPath, "bind", syscall.MS_BIND, "")
***REMOVED***

// GetSandboxForExternalKey returns sandbox object for the supplied path
func GetSandboxForExternalKey(basePath string, key string) (Sandbox, error) ***REMOVED***
	if err := createNamespaceFile(key); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := mountNetworkNamespace(basePath, key); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	n := &networkNamespace***REMOVED***path: key, nextIfIndex: make(map[string]int)***REMOVED***

	sboxNs, err := netns.GetFromPath(n.path)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed get network namespace %q: %v", n.path, err)
	***REMOVED***
	defer sboxNs.Close()

	n.nlHandle, err = netlink.NewHandleAt(sboxNs, syscall.NETLINK_ROUTE)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to create a netlink handle: %v", err)
	***REMOVED***

	err = n.nlHandle.SetSocketTimeout(ns.NetlinkSocketsTimeout)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to set the timeout on the sandbox netlink handle sockets: %v", err)
	***REMOVED***

	// As starting point, disable IPv6 on all interfaces
	err = setIPv6(n.path, "all", false)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to disable IPv6 on all interfaces on network namespace %q: %v", n.path, err)
	***REMOVED***

	if err = n.loopbackUp(); err != nil ***REMOVED***
		n.nlHandle.Delete()
		return nil, err
	***REMOVED***

	return n, nil
***REMOVED***

func reexecCreateNamespace() ***REMOVED***
	if len(os.Args) < 2 ***REMOVED***
		logrus.Fatal("no namespace path provided")
	***REMOVED***
	if err := mountNetworkNamespace("/proc/self/ns/net", os.Args[1]); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
***REMOVED***

func createNetworkNamespace(path string, osCreate bool) error ***REMOVED***
	if err := createNamespaceFile(path); err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd := &exec.Cmd***REMOVED***
		Path:   reexec.Self(),
		Args:   append([]string***REMOVED***"netns-create"***REMOVED***, path),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	***REMOVED***
	if osCreate ***REMOVED***
		cmd.SysProcAttr = &syscall.SysProcAttr***REMOVED******REMOVED***
		cmd.SysProcAttr.Cloneflags = syscall.CLONE_NEWNET
	***REMOVED***
	if err := cmd.Run(); err != nil ***REMOVED***
		return fmt.Errorf("namespace creation reexec command failed: %v", err)
	***REMOVED***

	return nil
***REMOVED***

func unmountNamespaceFile(path string) ***REMOVED***
	if _, err := os.Stat(path); err == nil ***REMOVED***
		syscall.Unmount(path, syscall.MNT_DETACH)
	***REMOVED***
***REMOVED***

func createNamespaceFile(path string) (err error) ***REMOVED***
	var f *os.File

	once.Do(createBasePath)
	// Remove it from garbage collection list if present
	removeFromGarbagePaths(path)

	// If the path is there unmount it first
	unmountNamespaceFile(path)

	// wait for garbage collection to complete if it is in progress
	// before trying to create the file.
	gpmWg.Wait()

	if f, err = os.Create(path); err == nil ***REMOVED***
		f.Close()
	***REMOVED***

	return err
***REMOVED***

func (n *networkNamespace) loopbackUp() error ***REMOVED***
	iface, err := n.nlHandle.LinkByName("lo")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return n.nlHandle.LinkSetUp(iface)
***REMOVED***

func (n *networkNamespace) AddLoopbackAliasIP(ip *net.IPNet) error ***REMOVED***
	iface, err := n.nlHandle.LinkByName("lo")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return n.nlHandle.AddrAdd(iface, &netlink.Addr***REMOVED***IPNet: ip***REMOVED***)
***REMOVED***

func (n *networkNamespace) RemoveLoopbackAliasIP(ip *net.IPNet) error ***REMOVED***
	iface, err := n.nlHandle.LinkByName("lo")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return n.nlHandle.AddrDel(iface, &netlink.Addr***REMOVED***IPNet: ip***REMOVED***)
***REMOVED***

func (n *networkNamespace) InvokeFunc(f func()) error ***REMOVED***
	return nsInvoke(n.nsPath(), func(nsFD int) error ***REMOVED*** return nil ***REMOVED***, func(callerFD int) error ***REMOVED***
		f()
		return nil
	***REMOVED***)
***REMOVED***

// InitOSContext initializes OS context while configuring network resources
func InitOSContext() func() ***REMOVED***
	runtime.LockOSThread()
	if err := ns.SetNamespace(); err != nil ***REMOVED***
		logrus.Error(err)
	***REMOVED***
	return runtime.UnlockOSThread
***REMOVED***

func nsInvoke(path string, prefunc func(nsFD int) error, postfunc func(callerFD int) error) error ***REMOVED***
	defer InitOSContext()()

	newNs, err := netns.GetFromPath(path)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed get network namespace %q: %v", path, err)
	***REMOVED***
	defer newNs.Close()

	// Invoked before the namespace switch happens but after the namespace file
	// handle is obtained.
	if err := prefunc(int(newNs)); err != nil ***REMOVED***
		return fmt.Errorf("failed in prefunc: %v", err)
	***REMOVED***

	if err = netns.Set(newNs); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer ns.SetNamespace()

	// Invoked after the namespace switch.
	return postfunc(ns.ParseHandlerInt())
***REMOVED***

func (n *networkNamespace) nsPath() string ***REMOVED***
	n.Lock()
	defer n.Unlock()

	return n.path
***REMOVED***

func (n *networkNamespace) Info() Info ***REMOVED***
	return n
***REMOVED***

func (n *networkNamespace) Key() string ***REMOVED***
	return n.path
***REMOVED***

func (n *networkNamespace) Destroy() error ***REMOVED***
	if n.nlHandle != nil ***REMOVED***
		n.nlHandle.Delete()
	***REMOVED***
	// Assuming no running process is executing in this network namespace,
	// unmounting is sufficient to destroy it.
	if err := syscall.Unmount(n.path, syscall.MNT_DETACH); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Stash it into the garbage collection list
	addToGarbagePaths(n.path)
	return nil
***REMOVED***

// Restore restore the network namespace
func (n *networkNamespace) Restore(ifsopt map[string][]IfaceOption, routes []*types.StaticRoute, gw net.IP, gw6 net.IP) error ***REMOVED***
	// restore interfaces
	for name, opts := range ifsopt ***REMOVED***
		if !strings.Contains(name, "+") ***REMOVED***
			return fmt.Errorf("wrong iface name in restore osl sandbox interface: %s", name)
		***REMOVED***
		seps := strings.Split(name, "+")
		srcName := seps[0]
		dstPrefix := seps[1]
		i := &nwIface***REMOVED***srcName: srcName, dstName: dstPrefix, ns: n***REMOVED***
		i.processInterfaceOptions(opts...)
		if i.master != "" ***REMOVED***
			i.dstMaster = n.findDst(i.master, true)
			if i.dstMaster == "" ***REMOVED***
				return fmt.Errorf("could not find an appropriate master %q for %q",
					i.master, i.srcName)
			***REMOVED***
		***REMOVED***
		if n.isDefault ***REMOVED***
			i.dstName = i.srcName
		***REMOVED*** else ***REMOVED***
			links, err := n.nlHandle.LinkList()
			if err != nil ***REMOVED***
				return fmt.Errorf("failed to retrieve list of links in network namespace %q during restore", n.path)
			***REMOVED***
			// due to the docker network connect/disconnect, so the dstName should
			// restore from the namespace
			for _, link := range links ***REMOVED***
				addrs, err := n.nlHandle.AddrList(link, netlink.FAMILY_V4)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				ifaceName := link.Attrs().Name
				if strings.HasPrefix(ifaceName, "vxlan") ***REMOVED***
					if i.dstName == "vxlan" ***REMOVED***
						i.dstName = ifaceName
						break
					***REMOVED***
				***REMOVED***
				// find the interface name by ip
				if i.address != nil ***REMOVED***
					for _, addr := range addrs ***REMOVED***
						if addr.IPNet.String() == i.address.String() ***REMOVED***
							i.dstName = ifaceName
							break
						***REMOVED***
						continue
					***REMOVED***
					if i.dstName == ifaceName ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				// This is to find the interface name of the pair in overlay sandbox
				if strings.HasPrefix(ifaceName, "veth") ***REMOVED***
					if i.master != "" && i.dstName == "veth" ***REMOVED***
						i.dstName = ifaceName
					***REMOVED***
				***REMOVED***
			***REMOVED***

			var index int
			indexStr := strings.TrimPrefix(i.dstName, dstPrefix)
			if indexStr != "" ***REMOVED***
				index, err = strconv.Atoi(indexStr)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			index++
			n.Lock()
			if index > n.nextIfIndex[dstPrefix] ***REMOVED***
				n.nextIfIndex[dstPrefix] = index
			***REMOVED***
			n.iFaces = append(n.iFaces, i)
			n.Unlock()
		***REMOVED***
	***REMOVED***

	// restore routes
	for _, r := range routes ***REMOVED***
		n.Lock()
		n.staticRoutes = append(n.staticRoutes, r)
		n.Unlock()
	***REMOVED***

	// restore gateway
	if len(gw) > 0 ***REMOVED***
		n.Lock()
		n.gw = gw
		n.Unlock()
	***REMOVED***

	if len(gw6) > 0 ***REMOVED***
		n.Lock()
		n.gwv6 = gw6
		n.Unlock()
	***REMOVED***

	return nil
***REMOVED***

// Checks whether IPv6 needs to be enabled/disabled on the loopback interface
func (n *networkNamespace) checkLoV6() ***REMOVED***
	var (
		enable = false
		action = "disable"
	)

	n.Lock()
	for _, iface := range n.iFaces ***REMOVED***
		if iface.AddressIPv6() != nil ***REMOVED***
			enable = true
			action = "enable"
			break
		***REMOVED***
	***REMOVED***
	n.Unlock()

	if n.loV6Enabled == enable ***REMOVED***
		return
	***REMOVED***

	if err := setIPv6(n.path, "lo", enable); err != nil ***REMOVED***
		logrus.Warnf("Failed to %s IPv6 on loopback interface on network namespace %q: %v", action, n.path, err)
	***REMOVED***

	n.loV6Enabled = enable
***REMOVED***

func reexecSetIPv6() ***REMOVED***
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(os.Args) < 3 ***REMOVED***
		logrus.Errorf("invalid number of arguments for %s", os.Args[0])
		os.Exit(1)
	***REMOVED***

	ns, err := netns.GetFromPath(os.Args[1])
	if err != nil ***REMOVED***
		logrus.Errorf("failed get network namespace %q: %v", os.Args[1], err)
		os.Exit(2)
	***REMOVED***
	defer ns.Close()

	if err = netns.Set(ns); err != nil ***REMOVED***
		logrus.Errorf("setting into container netns %q failed: %v", os.Args[1], err)
		os.Exit(3)
	***REMOVED***

	var (
		action = "disable"
		value  = byte('1')
		path   = fmt.Sprintf("/proc/sys/net/ipv6/conf/%s/disable_ipv6", os.Args[2])
	)

	if os.Args[3] == "true" ***REMOVED***
		action = "enable"
		value = byte('0')
	***REMOVED***

	if err = ioutil.WriteFile(path, []byte***REMOVED***value, '\n'***REMOVED***, 0644); err != nil ***REMOVED***
		logrus.Errorf("failed to %s IPv6 forwarding for container's interface %s: %v", action, os.Args[2], err)
		os.Exit(4)
	***REMOVED***

	os.Exit(0)
***REMOVED***

func setIPv6(path, iface string, enable bool) error ***REMOVED***
	cmd := &exec.Cmd***REMOVED***
		Path:   reexec.Self(),
		Args:   append([]string***REMOVED***"set-ipv6"***REMOVED***, path, iface, strconv.FormatBool(enable)),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	***REMOVED***
	if err := cmd.Run(); err != nil ***REMOVED***
		return fmt.Errorf("reexec to set IPv6 failed: %v", err)
	***REMOVED***
	return nil
***REMOVED***
