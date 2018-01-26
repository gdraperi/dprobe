package overlay

import (
	"fmt"
	"sync"

	"github.com/docker/libnetwork/iptables"
	"github.com/sirupsen/logrus"
)

const globalChain = "DOCKER-OVERLAY"

var filterOnce sync.Once

var filterChan = make(chan struct***REMOVED******REMOVED***, 1)

func filterWait() func() ***REMOVED***
	filterChan <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return func() ***REMOVED*** <-filterChan ***REMOVED***
***REMOVED***

func chainExists(cname string) bool ***REMOVED***
	if _, err := iptables.Raw("-L", cname); err != nil ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

func setupGlobalChain() ***REMOVED***
	// Because of an ungraceful shutdown, chain could already be present
	if !chainExists(globalChain) ***REMOVED***
		if err := iptables.RawCombinedOutput("-N", globalChain); err != nil ***REMOVED***
			logrus.Errorf("could not create global overlay chain: %v", err)
			return
		***REMOVED***
	***REMOVED***

	if !iptables.Exists(iptables.Filter, globalChain, "-j", "RETURN") ***REMOVED***
		if err := iptables.RawCombinedOutput("-A", globalChain, "-j", "RETURN"); err != nil ***REMOVED***
			logrus.Errorf("could not install default return chain in the overlay global chain: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func setNetworkChain(cname string, remove bool) error ***REMOVED***
	// Initialize the onetime global overlay chain
	filterOnce.Do(setupGlobalChain)

	exists := chainExists(cname)

	opt := "-N"
	// In case of remove, make sure to flush the rules in the chain
	if remove && exists ***REMOVED***
		if err := iptables.RawCombinedOutput("-F", cname); err != nil ***REMOVED***
			return fmt.Errorf("failed to flush overlay network chain %s rules: %v", cname, err)
		***REMOVED***
		opt = "-X"
	***REMOVED***

	if (!remove && !exists) || (remove && exists) ***REMOVED***
		if err := iptables.RawCombinedOutput(opt, cname); err != nil ***REMOVED***
			return fmt.Errorf("failed network chain operation %q for chain %s: %v", opt, cname, err)
		***REMOVED***
	***REMOVED***

	if !remove ***REMOVED***
		if !iptables.Exists(iptables.Filter, cname, "-j", "DROP") ***REMOVED***
			if err := iptables.RawCombinedOutput("-A", cname, "-j", "DROP"); err != nil ***REMOVED***
				return fmt.Errorf("failed adding default drop rule to overlay network chain %s: %v", cname, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func addNetworkChain(cname string) error ***REMOVED***
	defer filterWait()()

	return setNetworkChain(cname, false)
***REMOVED***

func removeNetworkChain(cname string) error ***REMOVED***
	defer filterWait()()

	return setNetworkChain(cname, true)
***REMOVED***

func setFilters(cname, brName string, remove bool) error ***REMOVED***
	opt := "-I"
	if remove ***REMOVED***
		opt = "-D"
	***REMOVED***

	// Every time we set filters for a new subnet make sure to move the global overlay hook to the top of the both the OUTPUT and forward chains
	if !remove ***REMOVED***
		for _, chain := range []string***REMOVED***"OUTPUT", "FORWARD"***REMOVED*** ***REMOVED***
			exists := iptables.Exists(iptables.Filter, chain, "-j", globalChain)
			if exists ***REMOVED***
				if err := iptables.RawCombinedOutput("-D", chain, "-j", globalChain); err != nil ***REMOVED***
					return fmt.Errorf("failed to delete overlay hook in chain %s while moving the hook: %v", chain, err)
				***REMOVED***
			***REMOVED***

			if err := iptables.RawCombinedOutput("-I", chain, "-j", globalChain); err != nil ***REMOVED***
				return fmt.Errorf("failed to insert overlay hook in chain %s: %v", chain, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Insert/Delete the rule to jump to per-bridge chain
	exists := iptables.Exists(iptables.Filter, globalChain, "-o", brName, "-j", cname)
	if (!remove && !exists) || (remove && exists) ***REMOVED***
		if err := iptables.RawCombinedOutput(opt, globalChain, "-o", brName, "-j", cname); err != nil ***REMOVED***
			return fmt.Errorf("failed to add per-bridge filter rule for bridge %s, network chain %s: %v", brName, cname, err)
		***REMOVED***
	***REMOVED***

	exists = iptables.Exists(iptables.Filter, cname, "-i", brName, "-j", "ACCEPT")
	if (!remove && exists) || (remove && !exists) ***REMOVED***
		return nil
	***REMOVED***

	if err := iptables.RawCombinedOutput(opt, cname, "-i", brName, "-j", "ACCEPT"); err != nil ***REMOVED***
		return fmt.Errorf("failed to add overlay filter rile for network chain %s, bridge %s: %v", cname, brName, err)
	***REMOVED***

	return nil
***REMOVED***

func addFilters(cname, brName string) error ***REMOVED***
	defer filterWait()()

	return setFilters(cname, brName, false)
***REMOVED***

func removeFilters(cname, brName string) error ***REMOVED***
	defer filterWait()()

	return setFilters(cname, brName, true)
***REMOVED***
