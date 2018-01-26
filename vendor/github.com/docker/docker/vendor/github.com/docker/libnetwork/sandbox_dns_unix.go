// +build !windows

package libnetwork

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/libnetwork/etchosts"
	"github.com/docker/libnetwork/resolvconf"
	"github.com/docker/libnetwork/resolvconf/dns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

const (
	defaultPrefix = "/var/lib/docker/network/files"
	dirPerm       = 0755
	filePerm      = 0644
)

func (sb *sandbox) startResolver(restore bool) ***REMOVED***
	sb.resolverOnce.Do(func() ***REMOVED***
		var err error
		sb.resolver = NewResolver(resolverIPSandbox, true, sb.Key(), sb)
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				sb.resolver = nil
			***REMOVED***
		***REMOVED***()

		// In the case of live restore container is already running with
		// right resolv.conf contents created before. Just update the
		// external DNS servers from the restored sandbox for embedded
		// server to use.
		if !restore ***REMOVED***
			err = sb.rebuildDNS()
			if err != nil ***REMOVED***
				logrus.Errorf("Updating resolv.conf failed for container %s, %q", sb.ContainerID(), err)
				return
			***REMOVED***
		***REMOVED***
		sb.resolver.SetExtServers(sb.extDNS)

		if err = sb.osSbox.InvokeFunc(sb.resolver.SetupFunc(0)); err != nil ***REMOVED***
			logrus.Errorf("Resolver Setup function failed for container %s, %q", sb.ContainerID(), err)
			return
		***REMOVED***

		if err = sb.resolver.Start(); err != nil ***REMOVED***
			logrus.Errorf("Resolver Start failed for container %s, %q", sb.ContainerID(), err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (sb *sandbox) setupResolutionFiles() error ***REMOVED***
	if err := sb.buildHostsFile(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := sb.updateParentHosts(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return sb.setupDNS()
***REMOVED***

func (sb *sandbox) buildHostsFile() error ***REMOVED***
	if sb.config.hostsPath == "" ***REMOVED***
		sb.config.hostsPath = defaultPrefix + "/" + sb.id + "/hosts"
	***REMOVED***

	dir, _ := filepath.Split(sb.config.hostsPath)
	if err := createBasePath(dir); err != nil ***REMOVED***
		return err
	***REMOVED***

	// This is for the host mode networking
	if sb.config.originHostsPath != "" ***REMOVED***
		if err := copyFile(sb.config.originHostsPath, sb.config.hostsPath); err != nil && !os.IsNotExist(err) ***REMOVED***
			return types.InternalErrorf("could not copy source hosts file %s to %s: %v", sb.config.originHostsPath, sb.config.hostsPath, err)
		***REMOVED***
		return nil
	***REMOVED***

	extraContent := make([]etchosts.Record, 0, len(sb.config.extraHosts))
	for _, extraHost := range sb.config.extraHosts ***REMOVED***
		extraContent = append(extraContent, etchosts.Record***REMOVED***Hosts: extraHost.name, IP: extraHost.IP***REMOVED***)
	***REMOVED***

	return etchosts.Build(sb.config.hostsPath, "", sb.config.hostName, sb.config.domainName, extraContent)
***REMOVED***

func (sb *sandbox) updateHostsFile(ifaceIP string) error ***REMOVED***
	if ifaceIP == "" ***REMOVED***
		return nil
	***REMOVED***

	if sb.config.originHostsPath != "" ***REMOVED***
		return nil
	***REMOVED***

	// User might have provided a FQDN in hostname or split it across hostname
	// and domainname.  We want the FQDN and the bare hostname.
	fqdn := sb.config.hostName
	mhost := sb.config.hostName
	if sb.config.domainName != "" ***REMOVED***
		fqdn = fmt.Sprintf("%s.%s", fqdn, sb.config.domainName)
	***REMOVED***

	parts := strings.SplitN(fqdn, ".", 2)
	if len(parts) == 2 ***REMOVED***
		mhost = fmt.Sprintf("%s %s", fqdn, parts[0])
	***REMOVED***

	extraContent := []etchosts.Record***REMOVED******REMOVED***Hosts: mhost, IP: ifaceIP***REMOVED******REMOVED***

	sb.addHostsEntries(extraContent)
	return nil
***REMOVED***

func (sb *sandbox) addHostsEntries(recs []etchosts.Record) ***REMOVED***
	if err := etchosts.Add(sb.config.hostsPath, recs); err != nil ***REMOVED***
		logrus.Warnf("Failed adding service host entries to the running container: %v", err)
	***REMOVED***
***REMOVED***

func (sb *sandbox) deleteHostsEntries(recs []etchosts.Record) ***REMOVED***
	if err := etchosts.Delete(sb.config.hostsPath, recs); err != nil ***REMOVED***
		logrus.Warnf("Failed deleting service host entries to the running container: %v", err)
	***REMOVED***
***REMOVED***

func (sb *sandbox) updateParentHosts() error ***REMOVED***
	var pSb Sandbox

	for _, update := range sb.config.parentUpdates ***REMOVED***
		sb.controller.WalkSandboxes(SandboxContainerWalker(&pSb, update.cid))
		if pSb == nil ***REMOVED***
			continue
		***REMOVED***
		if err := etchosts.Update(pSb.(*sandbox).config.hostsPath, update.ip, update.name); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (sb *sandbox) restorePath() ***REMOVED***
	if sb.config.resolvConfPath == "" ***REMOVED***
		sb.config.resolvConfPath = defaultPrefix + "/" + sb.id + "/resolv.conf"
	***REMOVED***
	sb.config.resolvConfHashFile = sb.config.resolvConfPath + ".hash"
	if sb.config.hostsPath == "" ***REMOVED***
		sb.config.hostsPath = defaultPrefix + "/" + sb.id + "/hosts"
	***REMOVED***
***REMOVED***

func (sb *sandbox) setExternalResolvers(content []byte, addrType int, checkLoopback bool) ***REMOVED***
	servers := resolvconf.GetNameservers(content, addrType)
	for _, ip := range servers ***REMOVED***
		hostLoopback := false
		if checkLoopback ***REMOVED***
			hostLoopback = dns.IsIPv4Localhost(ip)
		***REMOVED***
		sb.extDNS = append(sb.extDNS, extDNSEntry***REMOVED***
			IPStr:        ip,
			HostLoopback: hostLoopback,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (sb *sandbox) setupDNS() error ***REMOVED***
	var newRC *resolvconf.File

	if sb.config.resolvConfPath == "" ***REMOVED***
		sb.config.resolvConfPath = defaultPrefix + "/" + sb.id + "/resolv.conf"
	***REMOVED***

	sb.config.resolvConfHashFile = sb.config.resolvConfPath + ".hash"

	dir, _ := filepath.Split(sb.config.resolvConfPath)
	if err := createBasePath(dir); err != nil ***REMOVED***
		return err
	***REMOVED***

	// This is for the host mode networking
	if sb.config.originResolvConfPath != "" ***REMOVED***
		if err := copyFile(sb.config.originResolvConfPath, sb.config.resolvConfPath); err != nil ***REMOVED***
			if !os.IsNotExist(err) ***REMOVED***
				return fmt.Errorf("could not copy source resolv.conf file %s to %s: %v", sb.config.originResolvConfPath, sb.config.resolvConfPath, err)
			***REMOVED***
			logrus.Infof("%s does not exist, we create an empty resolv.conf for container", sb.config.originResolvConfPath)
			if err := createFile(sb.config.resolvConfPath); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	currRC, err := resolvconf.Get()
	if err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED***
		// it's ok to continue if /etc/resolv.conf doesn't exist, default resolvers (Google's Public DNS)
		// will be used
		currRC = &resolvconf.File***REMOVED******REMOVED***
		logrus.Infof("/etc/resolv.conf does not exist")
	***REMOVED***

	if len(sb.config.dnsList) > 0 || len(sb.config.dnsSearchList) > 0 || len(sb.config.dnsOptionsList) > 0 ***REMOVED***
		var (
			err            error
			dnsList        = resolvconf.GetNameservers(currRC.Content, types.IP)
			dnsSearchList  = resolvconf.GetSearchDomains(currRC.Content)
			dnsOptionsList = resolvconf.GetOptions(currRC.Content)
		)
		if len(sb.config.dnsList) > 0 ***REMOVED***
			dnsList = sb.config.dnsList
		***REMOVED***
		if len(sb.config.dnsSearchList) > 0 ***REMOVED***
			dnsSearchList = sb.config.dnsSearchList
		***REMOVED***
		if len(sb.config.dnsOptionsList) > 0 ***REMOVED***
			dnsOptionsList = sb.config.dnsOptionsList
		***REMOVED***
		newRC, err = resolvconf.Build(sb.config.resolvConfPath, dnsList, dnsSearchList, dnsOptionsList)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// After building the resolv.conf from the user config save the
		// external resolvers in the sandbox. Note that --dns 127.0.0.x
		// config refers to the loopback in the container namespace
		sb.setExternalResolvers(newRC.Content, types.IPv4, false)
	***REMOVED*** else ***REMOVED***
		// If the host resolv.conf file has 127.0.0.x container should
		// use the host restolver for queries. This is supported by the
		// docker embedded DNS server. Hence save the external resolvers
		// before filtering it out.
		sb.setExternalResolvers(currRC.Content, types.IPv4, true)

		// Replace any localhost/127.* (at this point we have no info about ipv6, pass it as true)
		if newRC, err = resolvconf.FilterResolvDNS(currRC.Content, true); err != nil ***REMOVED***
			return err
		***REMOVED***
		// No contention on container resolv.conf file at sandbox creation
		if err := ioutil.WriteFile(sb.config.resolvConfPath, newRC.Content, filePerm); err != nil ***REMOVED***
			return types.InternalErrorf("failed to write unhaltered resolv.conf file content when setting up dns for sandbox %s: %v", sb.ID(), err)
		***REMOVED***
	***REMOVED***

	// Write hash
	if err := ioutil.WriteFile(sb.config.resolvConfHashFile, []byte(newRC.Hash), filePerm); err != nil ***REMOVED***
		return types.InternalErrorf("failed to write resolv.conf hash file when setting up dns for sandbox %s: %v", sb.ID(), err)
	***REMOVED***

	return nil
***REMOVED***

func (sb *sandbox) updateDNS(ipv6Enabled bool) error ***REMOVED***
	var (
		currHash string
		hashFile = sb.config.resolvConfHashFile
	)

	// This is for the host mode networking
	if sb.config.originResolvConfPath != "" ***REMOVED***
		return nil
	***REMOVED***

	if len(sb.config.dnsList) > 0 || len(sb.config.dnsSearchList) > 0 || len(sb.config.dnsOptionsList) > 0 ***REMOVED***
		return nil
	***REMOVED***

	currRC, err := resolvconf.GetSpecific(sb.config.resolvConfPath)
	if err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		h, err := ioutil.ReadFile(hashFile)
		if err != nil ***REMOVED***
			if !os.IsNotExist(err) ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			currHash = string(h)
		***REMOVED***
	***REMOVED***

	if currHash != "" && currHash != currRC.Hash ***REMOVED***
		// Seems the user has changed the container resolv.conf since the last time
		// we checked so return without doing anything.
		//logrus.Infof("Skipping update of resolv.conf file with ipv6Enabled: %t because file was touched by user", ipv6Enabled)
		return nil
	***REMOVED***

	// replace any localhost/127.* and remove IPv6 nameservers if IPv6 disabled.
	newRC, err := resolvconf.FilterResolvDNS(currRC.Content, ipv6Enabled)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = ioutil.WriteFile(sb.config.resolvConfPath, newRC.Content, 0644)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// write the new hash in a temp file and rename it to make the update atomic
	dir := path.Dir(sb.config.resolvConfPath)
	tmpHashFile, err := ioutil.TempFile(dir, "hash")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = tmpHashFile.Chmod(filePerm); err != nil ***REMOVED***
		tmpHashFile.Close()
		return err
	***REMOVED***
	_, err = tmpHashFile.Write([]byte(newRC.Hash))
	if err1 := tmpHashFile.Close(); err == nil ***REMOVED***
		err = err1
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return os.Rename(tmpHashFile.Name(), hashFile)
***REMOVED***

// Embedded DNS server has to be enabled for this sandbox. Rebuild the container's
// resolv.conf by doing the following
// - Add only the embedded server's IP to container's resolv.conf
// - If the embedded server needs any resolv.conf options add it to the current list
func (sb *sandbox) rebuildDNS() error ***REMOVED***
	currRC, err := resolvconf.GetSpecific(sb.config.resolvConfPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(sb.extDNS) == 0 ***REMOVED***
		sb.setExternalResolvers(currRC.Content, types.IPv4, false)
	***REMOVED***
	var (
		dnsList        = []string***REMOVED***sb.resolver.NameServer()***REMOVED***
		dnsOptionsList = resolvconf.GetOptions(currRC.Content)
		dnsSearchList  = resolvconf.GetSearchDomains(currRC.Content)
	)

	// external v6 DNS servers has to be listed in resolv.conf
	dnsList = append(dnsList, resolvconf.GetNameservers(currRC.Content, types.IPv6)...)

	// If the user config and embedded DNS server both have ndots option set,
	// remember the user's config so that unqualified names not in the docker
	// domain can be dropped.
	resOptions := sb.resolver.ResolverOptions()

dnsOpt:
	for _, resOpt := range resOptions ***REMOVED***
		if strings.Contains(resOpt, "ndots") ***REMOVED***
			for i, option := range dnsOptionsList ***REMOVED***
				if strings.Contains(option, "ndots") ***REMOVED***
					parts := strings.Split(option, ":")
					if len(parts) != 2 ***REMOVED***
						return fmt.Errorf("invalid ndots option %v", option)
					***REMOVED***
					if num, err := strconv.Atoi(parts[1]); err != nil ***REMOVED***
						return fmt.Errorf("invalid number for ndots option %v", option)
					***REMOVED*** else if num > 0 ***REMOVED***
						// if the user sets ndots, we mark it as set but we remove the option to guarantee
						// that into the container land only ndots:0
						sb.ndotsSet = true
						dnsOptionsList = append(dnsOptionsList[:i], dnsOptionsList[i+1:]...)
						break dnsOpt
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	dnsOptionsList = append(dnsOptionsList, resOptions...)

	_, err = resolvconf.Build(sb.config.resolvConfPath, dnsList, dnsSearchList, dnsOptionsList)
	return err
***REMOVED***

func createBasePath(dir string) error ***REMOVED***
	return os.MkdirAll(dir, dirPerm)
***REMOVED***

func createFile(path string) error ***REMOVED***
	var f *os.File

	dir, _ := filepath.Split(path)
	err := createBasePath(dir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	f, err = os.Create(path)
	if err == nil ***REMOVED***
		f.Close()
	***REMOVED***

	return err
***REMOVED***

func copyFile(src, dst string) error ***REMOVED***
	sBytes, err := ioutil.ReadFile(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(dst, sBytes, filePerm)
***REMOVED***
