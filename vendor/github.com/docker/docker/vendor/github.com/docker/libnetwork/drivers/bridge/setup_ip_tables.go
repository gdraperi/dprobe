package bridge

import (
	"errors"
	"fmt"
	"net"

	"github.com/docker/libnetwork/iptables"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// DockerChain: DOCKER iptable chain name
const (
	DockerChain    = "DOCKER"
	IsolationChain = "DOCKER-ISOLATION"
)

func setupIPChains(config *configuration) (*iptables.ChainInfo, *iptables.ChainInfo, *iptables.ChainInfo, error) ***REMOVED***
	// Sanity check.
	if config.EnableIPTables == false ***REMOVED***
		return nil, nil, nil, errors.New("cannot create new chains, EnableIPTable is disabled")
	***REMOVED***

	hairpinMode := !config.EnableUserlandProxy

	natChain, err := iptables.NewChain(DockerChain, iptables.Nat, hairpinMode)
	if err != nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("failed to create NAT chain: %v", err)
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if err := iptables.RemoveExistingChain(DockerChain, iptables.Nat); err != nil ***REMOVED***
				logrus.Warnf("failed on removing iptables NAT chain on cleanup: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	filterChain, err := iptables.NewChain(DockerChain, iptables.Filter, false)
	if err != nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("failed to create FILTER chain: %v", err)
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if err := iptables.RemoveExistingChain(DockerChain, iptables.Filter); err != nil ***REMOVED***
				logrus.Warnf("failed on removing iptables FILTER chain on cleanup: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	isolationChain, err := iptables.NewChain(IsolationChain, iptables.Filter, false)
	if err != nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("failed to create FILTER isolation chain: %v", err)
	***REMOVED***

	if err := iptables.AddReturnRule(IsolationChain); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	return natChain, filterChain, isolationChain, nil
***REMOVED***

func (n *bridgeNetwork) setupIPTables(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	var err error

	d := n.driver
	d.Lock()
	driverConfig := d.config
	d.Unlock()

	// Sanity check.
	if driverConfig.EnableIPTables == false ***REMOVED***
		return errors.New("Cannot program chains, EnableIPTable is disabled")
	***REMOVED***

	// Pickup this configuration option from driver
	hairpinMode := !driverConfig.EnableUserlandProxy

	maskedAddrv4 := &net.IPNet***REMOVED***
		IP:   i.bridgeIPv4.IP.Mask(i.bridgeIPv4.Mask),
		Mask: i.bridgeIPv4.Mask,
	***REMOVED***
	if config.Internal ***REMOVED***
		if err = setupInternalNetworkRules(config.BridgeName, maskedAddrv4, config.EnableICC, true); err != nil ***REMOVED***
			return fmt.Errorf("Failed to Setup IP tables: %s", err.Error())
		***REMOVED***
		n.registerIptCleanFunc(func() error ***REMOVED***
			return setupInternalNetworkRules(config.BridgeName, maskedAddrv4, config.EnableICC, false)
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		if err = setupIPTablesInternal(config.BridgeName, maskedAddrv4, config.EnableICC, config.EnableIPMasquerade, hairpinMode, true); err != nil ***REMOVED***
			return fmt.Errorf("Failed to Setup IP tables: %s", err.Error())
		***REMOVED***
		n.registerIptCleanFunc(func() error ***REMOVED***
			return setupIPTablesInternal(config.BridgeName, maskedAddrv4, config.EnableICC, config.EnableIPMasquerade, hairpinMode, false)
		***REMOVED***)
		natChain, filterChain, _, err := n.getDriverChains()
		if err != nil ***REMOVED***
			return fmt.Errorf("Failed to setup IP tables, cannot acquire chain info %s", err.Error())
		***REMOVED***

		err = iptables.ProgramChain(natChain, config.BridgeName, hairpinMode, true)
		if err != nil ***REMOVED***
			return fmt.Errorf("Failed to program NAT chain: %s", err.Error())
		***REMOVED***

		err = iptables.ProgramChain(filterChain, config.BridgeName, hairpinMode, true)
		if err != nil ***REMOVED***
			return fmt.Errorf("Failed to program FILTER chain: %s", err.Error())
		***REMOVED***

		n.registerIptCleanFunc(func() error ***REMOVED***
			return iptables.ProgramChain(filterChain, config.BridgeName, hairpinMode, false)
		***REMOVED***)

		n.portMapper.SetIptablesChain(natChain, n.getNetworkBridgeName())
	***REMOVED***

	d.Lock()
	err = iptables.EnsureJumpRule("FORWARD", IsolationChain)
	d.Unlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

type iptRule struct ***REMOVED***
	table   iptables.Table
	chain   string
	preArgs []string
	args    []string
***REMOVED***

func setupIPTablesInternal(bridgeIface string, addr net.Addr, icc, ipmasq, hairpin, enable bool) error ***REMOVED***

	var (
		address   = addr.String()
		natRule   = iptRule***REMOVED***table: iptables.Nat, chain: "POSTROUTING", preArgs: []string***REMOVED***"-t", "nat"***REMOVED***, args: []string***REMOVED***"-s", address, "!", "-o", bridgeIface, "-j", "MASQUERADE"***REMOVED******REMOVED***
		hpNatRule = iptRule***REMOVED***table: iptables.Nat, chain: "POSTROUTING", preArgs: []string***REMOVED***"-t", "nat"***REMOVED***, args: []string***REMOVED***"-m", "addrtype", "--src-type", "LOCAL", "-o", bridgeIface, "-j", "MASQUERADE"***REMOVED******REMOVED***
		skipDNAT  = iptRule***REMOVED***table: iptables.Nat, chain: DockerChain, preArgs: []string***REMOVED***"-t", "nat"***REMOVED***, args: []string***REMOVED***"-i", bridgeIface, "-j", "RETURN"***REMOVED******REMOVED***
		outRule   = iptRule***REMOVED***table: iptables.Filter, chain: "FORWARD", args: []string***REMOVED***"-i", bridgeIface, "!", "-o", bridgeIface, "-j", "ACCEPT"***REMOVED******REMOVED***
	)

	// Set NAT.
	if ipmasq ***REMOVED***
		if err := programChainRule(natRule, "NAT", enable); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if ipmasq && !hairpin ***REMOVED***
		if err := programChainRule(skipDNAT, "SKIP DNAT", enable); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// In hairpin mode, masquerade traffic from localhost
	if hairpin ***REMOVED***
		if err := programChainRule(hpNatRule, "MASQ LOCAL HOST", enable); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Set Inter Container Communication.
	if err := setIcc(bridgeIface, icc, enable); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Set Accept on all non-intercontainer outgoing packets.
	return programChainRule(outRule, "ACCEPT NON_ICC OUTGOING", enable)
***REMOVED***

func programChainRule(rule iptRule, ruleDescr string, insert bool) error ***REMOVED***
	var (
		prefix    []string
		operation string
		condition bool
		doesExist = iptables.Exists(rule.table, rule.chain, rule.args...)
	)

	if insert ***REMOVED***
		condition = !doesExist
		prefix = []string***REMOVED***"-I", rule.chain***REMOVED***
		operation = "enable"
	***REMOVED*** else ***REMOVED***
		condition = doesExist
		prefix = []string***REMOVED***"-D", rule.chain***REMOVED***
		operation = "disable"
	***REMOVED***
	if rule.preArgs != nil ***REMOVED***
		prefix = append(rule.preArgs, prefix...)
	***REMOVED***

	if condition ***REMOVED***
		if err := iptables.RawCombinedOutput(append(prefix, rule.args...)...); err != nil ***REMOVED***
			return fmt.Errorf("Unable to %s %s rule: %s", operation, ruleDescr, err.Error())
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func setIcc(bridgeIface string, iccEnable, insert bool) error ***REMOVED***
	var (
		table      = iptables.Filter
		chain      = "FORWARD"
		args       = []string***REMOVED***"-i", bridgeIface, "-o", bridgeIface, "-j"***REMOVED***
		acceptArgs = append(args, "ACCEPT")
		dropArgs   = append(args, "DROP")
	)

	if insert ***REMOVED***
		if !iccEnable ***REMOVED***
			iptables.Raw(append([]string***REMOVED***"-D", chain***REMOVED***, acceptArgs...)...)

			if !iptables.Exists(table, chain, dropArgs...) ***REMOVED***
				if err := iptables.RawCombinedOutput(append([]string***REMOVED***"-A", chain***REMOVED***, dropArgs...)...); err != nil ***REMOVED***
					return fmt.Errorf("Unable to prevent intercontainer communication: %s", err.Error())
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			iptables.Raw(append([]string***REMOVED***"-D", chain***REMOVED***, dropArgs...)...)

			if !iptables.Exists(table, chain, acceptArgs...) ***REMOVED***
				if err := iptables.RawCombinedOutput(append([]string***REMOVED***"-I", chain***REMOVED***, acceptArgs...)...); err != nil ***REMOVED***
					return fmt.Errorf("Unable to allow intercontainer communication: %s", err.Error())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Remove any ICC rule.
		if !iccEnable ***REMOVED***
			if iptables.Exists(table, chain, dropArgs...) ***REMOVED***
				iptables.Raw(append([]string***REMOVED***"-D", chain***REMOVED***, dropArgs...)...)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if iptables.Exists(table, chain, acceptArgs...) ***REMOVED***
				iptables.Raw(append([]string***REMOVED***"-D", chain***REMOVED***, acceptArgs...)...)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Control Inter Network Communication. Install/remove only if it is not/is present.
func setINC(iface1, iface2 string, enable bool) error ***REMOVED***
	var (
		table = iptables.Filter
		chain = IsolationChain
		args  = [2][]string***REMOVED******REMOVED***"-i", iface1, "-o", iface2, "-j", "DROP"***REMOVED***, ***REMOVED***"-i", iface2, "-o", iface1, "-j", "DROP"***REMOVED******REMOVED***
	)

	if enable ***REMOVED***
		for i := 0; i < 2; i++ ***REMOVED***
			if iptables.Exists(table, chain, args[i]...) ***REMOVED***
				continue
			***REMOVED***
			if err := iptables.RawCombinedOutput(append([]string***REMOVED***"-I", chain***REMOVED***, args[i]...)...); err != nil ***REMOVED***
				return fmt.Errorf("unable to add inter-network communication rule: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := 0; i < 2; i++ ***REMOVED***
			if !iptables.Exists(table, chain, args[i]...) ***REMOVED***
				continue
			***REMOVED***
			if err := iptables.RawCombinedOutput(append([]string***REMOVED***"-D", chain***REMOVED***, args[i]...)...); err != nil ***REMOVED***
				return fmt.Errorf("unable to remove inter-network communication rule: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func removeIPChains() ***REMOVED***
	for _, chainInfo := range []iptables.ChainInfo***REMOVED***
		***REMOVED***Name: DockerChain, Table: iptables.Nat***REMOVED***,
		***REMOVED***Name: DockerChain, Table: iptables.Filter***REMOVED***,
		***REMOVED***Name: IsolationChain, Table: iptables.Filter***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := chainInfo.Remove(); err != nil ***REMOVED***
			logrus.Warnf("Failed to remove existing iptables entries in table %s chain %s : %v", chainInfo.Table, chainInfo.Name, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func setupInternalNetworkRules(bridgeIface string, addr net.Addr, icc, insert bool) error ***REMOVED***
	var (
		inDropRule  = iptRule***REMOVED***table: iptables.Filter, chain: IsolationChain, args: []string***REMOVED***"-i", bridgeIface, "!", "-d", addr.String(), "-j", "DROP"***REMOVED******REMOVED***
		outDropRule = iptRule***REMOVED***table: iptables.Filter, chain: IsolationChain, args: []string***REMOVED***"-o", bridgeIface, "!", "-s", addr.String(), "-j", "DROP"***REMOVED******REMOVED***
	)
	if err := programChainRule(inDropRule, "DROP INCOMING", insert); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := programChainRule(outDropRule, "DROP OUTGOING", insert); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Set Inter Container Communication.
	return setIcc(bridgeIface, icc, insert)
***REMOVED***

func clearEndpointConnections(nlh *netlink.Handle, ep *bridgeEndpoint) ***REMOVED***
	var ipv4List []net.IP
	var ipv6List []net.IP
	if ep.addr != nil ***REMOVED***
		ipv4List = append(ipv4List, ep.addr.IP)
	***REMOVED***
	if ep.addrv6 != nil ***REMOVED***
		ipv6List = append(ipv6List, ep.addrv6.IP)
	***REMOVED***
	iptables.DeleteConntrackEntries(nlh, ipv4List, ipv6List)
***REMOVED***
