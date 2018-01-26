package iptables

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Action signifies the iptable action.
type Action string

// Policy is the default iptable policies
type Policy string

// Table refers to Nat, Filter or Mangle.
type Table string

const (
	// Append appends the rule at the end of the chain.
	Append Action = "-A"
	// Delete deletes the rule from the chain.
	Delete Action = "-D"
	// Insert inserts the rule at the top of the chain.
	Insert Action = "-I"
	// Nat table is used for nat translation rules.
	Nat Table = "nat"
	// Filter table is used for filter rules.
	Filter Table = "filter"
	// Mangle table is used for mangling the packet.
	Mangle Table = "mangle"
	// Drop is the default iptables DROP policy
	Drop Policy = "DROP"
	// Accept is the default iptables ACCEPT policy
	Accept Policy = "ACCEPT"
)

var (
	iptablesPath  string
	supportsXlock = false
	supportsCOpt  = false
	xLockWaitMsg  = "Another app is currently holding the xtables lock; waiting"
	// used to lock iptables commands if xtables lock is not supported
	bestEffortLock sync.Mutex
	// ErrIptablesNotFound is returned when the rule is not found.
	ErrIptablesNotFound = errors.New("Iptables not found")
	initOnce            sync.Once
)

// ChainInfo defines the iptables chain.
type ChainInfo struct ***REMOVED***
	Name        string
	Table       Table
	HairpinMode bool
***REMOVED***

// ChainError is returned to represent errors during ip table operation.
type ChainError struct ***REMOVED***
	Chain  string
	Output []byte
***REMOVED***

func (e ChainError) Error() string ***REMOVED***
	return fmt.Sprintf("Error iptables %s: %s", e.Chain, string(e.Output))
***REMOVED***

func probe() ***REMOVED***
	if out, err := exec.Command("modprobe", "-va", "nf_nat").CombinedOutput(); err != nil ***REMOVED***
		logrus.Warnf("Running modprobe nf_nat failed with message: `%s`, error: %v", strings.TrimSpace(string(out)), err)
	***REMOVED***
	if out, err := exec.Command("modprobe", "-va", "xt_conntrack").CombinedOutput(); err != nil ***REMOVED***
		logrus.Warnf("Running modprobe xt_conntrack failed with message: `%s`, error: %v", strings.TrimSpace(string(out)), err)
	***REMOVED***
***REMOVED***

func initFirewalld() ***REMOVED***
	if err := FirewalldInit(); err != nil ***REMOVED***
		logrus.Debugf("Fail to initialize firewalld: %v, using raw iptables instead", err)
	***REMOVED***
***REMOVED***

func detectIptables() ***REMOVED***
	path, err := exec.LookPath("iptables")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	iptablesPath = path
	supportsXlock = exec.Command(iptablesPath, "--wait", "-L", "-n").Run() == nil
	mj, mn, mc, err := GetVersion()
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to read iptables version: %v", err)
		return
	***REMOVED***
	supportsCOpt = supportsCOption(mj, mn, mc)
***REMOVED***

func initDependencies() ***REMOVED***
	probe()
	initFirewalld()
	detectIptables()
***REMOVED***

func initCheck() error ***REMOVED***
	initOnce.Do(initDependencies)

	if iptablesPath == "" ***REMOVED***
		return ErrIptablesNotFound
	***REMOVED***
	return nil
***REMOVED***

// NewChain adds a new chain to ip table.
func NewChain(name string, table Table, hairpinMode bool) (*ChainInfo, error) ***REMOVED***
	c := &ChainInfo***REMOVED***
		Name:        name,
		Table:       table,
		HairpinMode: hairpinMode,
	***REMOVED***
	if string(c.Table) == "" ***REMOVED***
		c.Table = Filter
	***REMOVED***

	// Add chain if it doesn't exist
	if _, err := Raw("-t", string(c.Table), "-n", "-L", c.Name); err != nil ***REMOVED***
		if output, err := Raw("-t", string(c.Table), "-N", c.Name); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if len(output) != 0 ***REMOVED***
			return nil, fmt.Errorf("Could not create %s/%s chain: %s", c.Table, c.Name, output)
		***REMOVED***
	***REMOVED***
	return c, nil
***REMOVED***

// ProgramChain is used to add rules to a chain
func ProgramChain(c *ChainInfo, bridgeName string, hairpinMode, enable bool) error ***REMOVED***
	if c.Name == "" ***REMOVED***
		return errors.New("Could not program chain, missing chain name")
	***REMOVED***

	switch c.Table ***REMOVED***
	case Nat:
		preroute := []string***REMOVED***
			"-m", "addrtype",
			"--dst-type", "LOCAL",
			"-j", c.Name***REMOVED***
		if !Exists(Nat, "PREROUTING", preroute...) && enable ***REMOVED***
			if err := c.Prerouting(Append, preroute...); err != nil ***REMOVED***
				return fmt.Errorf("Failed to inject %s in PREROUTING chain: %s", c.Name, err)
			***REMOVED***
		***REMOVED*** else if Exists(Nat, "PREROUTING", preroute...) && !enable ***REMOVED***
			if err := c.Prerouting(Delete, preroute...); err != nil ***REMOVED***
				return fmt.Errorf("Failed to remove %s in PREROUTING chain: %s", c.Name, err)
			***REMOVED***
		***REMOVED***
		output := []string***REMOVED***
			"-m", "addrtype",
			"--dst-type", "LOCAL",
			"-j", c.Name***REMOVED***
		if !hairpinMode ***REMOVED***
			output = append(output, "!", "--dst", "127.0.0.0/8")
		***REMOVED***
		if !Exists(Nat, "OUTPUT", output...) && enable ***REMOVED***
			if err := c.Output(Append, output...); err != nil ***REMOVED***
				return fmt.Errorf("Failed to inject %s in OUTPUT chain: %s", c.Name, err)
			***REMOVED***
		***REMOVED*** else if Exists(Nat, "OUTPUT", output...) && !enable ***REMOVED***
			if err := c.Output(Delete, output...); err != nil ***REMOVED***
				return fmt.Errorf("Failed to inject %s in OUTPUT chain: %s", c.Name, err)
			***REMOVED***
		***REMOVED***
	case Filter:
		if bridgeName == "" ***REMOVED***
			return fmt.Errorf("Could not program chain %s/%s, missing bridge name",
				c.Table, c.Name)
		***REMOVED***
		link := []string***REMOVED***
			"-o", bridgeName,
			"-j", c.Name***REMOVED***
		if !Exists(Filter, "FORWARD", link...) && enable ***REMOVED***
			insert := append([]string***REMOVED***string(Insert), "FORWARD"***REMOVED***, link...)
			if output, err := Raw(insert...); err != nil ***REMOVED***
				return err
			***REMOVED*** else if len(output) != 0 ***REMOVED***
				return fmt.Errorf("Could not create linking rule to %s/%s: %s", c.Table, c.Name, output)
			***REMOVED***
		***REMOVED*** else if Exists(Filter, "FORWARD", link...) && !enable ***REMOVED***
			del := append([]string***REMOVED***string(Delete), "FORWARD"***REMOVED***, link...)
			if output, err := Raw(del...); err != nil ***REMOVED***
				return err
			***REMOVED*** else if len(output) != 0 ***REMOVED***
				return fmt.Errorf("Could not delete linking rule from %s/%s: %s", c.Table, c.Name, output)
			***REMOVED***

		***REMOVED***
		establish := []string***REMOVED***
			"-o", bridgeName,
			"-m", "conntrack",
			"--ctstate", "RELATED,ESTABLISHED",
			"-j", "ACCEPT"***REMOVED***
		if !Exists(Filter, "FORWARD", establish...) && enable ***REMOVED***
			insert := append([]string***REMOVED***string(Insert), "FORWARD"***REMOVED***, establish...)
			if output, err := Raw(insert...); err != nil ***REMOVED***
				return err
			***REMOVED*** else if len(output) != 0 ***REMOVED***
				return fmt.Errorf("Could not create establish rule to %s: %s", c.Table, output)
			***REMOVED***
		***REMOVED*** else if Exists(Filter, "FORWARD", establish...) && !enable ***REMOVED***
			del := append([]string***REMOVED***string(Delete), "FORWARD"***REMOVED***, establish...)
			if output, err := Raw(del...); err != nil ***REMOVED***
				return err
			***REMOVED*** else if len(output) != 0 ***REMOVED***
				return fmt.Errorf("Could not delete establish rule from %s: %s", c.Table, output)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// RemoveExistingChain removes existing chain from the table.
func RemoveExistingChain(name string, table Table) error ***REMOVED***
	c := &ChainInfo***REMOVED***
		Name:  name,
		Table: table,
	***REMOVED***
	if string(c.Table) == "" ***REMOVED***
		c.Table = Filter
	***REMOVED***
	return c.Remove()
***REMOVED***

// Forward adds forwarding rule to 'filter' table and corresponding nat rule to 'nat' table.
func (c *ChainInfo) Forward(action Action, ip net.IP, port int, proto, destAddr string, destPort int, bridgeName string) error ***REMOVED***
	daddr := ip.String()
	if ip.IsUnspecified() ***REMOVED***
		// iptables interprets "0.0.0.0" as "0.0.0.0/32", whereas we
		// want "0.0.0.0/0". "0/0" is correctly interpreted as "any
		// value" by both iptables and ip6tables.
		daddr = "0/0"
	***REMOVED***

	args := []string***REMOVED***
		"-p", proto,
		"-d", daddr,
		"--dport", strconv.Itoa(port),
		"-j", "DNAT",
		"--to-destination", net.JoinHostPort(destAddr, strconv.Itoa(destPort))***REMOVED***
	if !c.HairpinMode ***REMOVED***
		args = append(args, "!", "-i", bridgeName)
	***REMOVED***
	if err := ProgramRule(Nat, c.Name, action, args); err != nil ***REMOVED***
		return err
	***REMOVED***

	args = []string***REMOVED***
		"!", "-i", bridgeName,
		"-o", bridgeName,
		"-p", proto,
		"-d", destAddr,
		"--dport", strconv.Itoa(destPort),
		"-j", "ACCEPT",
	***REMOVED***
	if err := ProgramRule(Filter, c.Name, action, args); err != nil ***REMOVED***
		return err
	***REMOVED***

	args = []string***REMOVED***
		"-p", proto,
		"-s", destAddr,
		"-d", destAddr,
		"--dport", strconv.Itoa(destPort),
		"-j", "MASQUERADE",
	***REMOVED***
	return ProgramRule(Nat, "POSTROUTING", action, args)
***REMOVED***

// Link adds reciprocal ACCEPT rule for two supplied IP addresses.
// Traffic is allowed from ip1 to ip2 and vice-versa
func (c *ChainInfo) Link(action Action, ip1, ip2 net.IP, port int, proto string, bridgeName string) error ***REMOVED***
	// forward
	args := []string***REMOVED***
		"-i", bridgeName, "-o", bridgeName,
		"-p", proto,
		"-s", ip1.String(),
		"-d", ip2.String(),
		"--dport", strconv.Itoa(port),
		"-j", "ACCEPT",
	***REMOVED***
	if err := ProgramRule(Filter, c.Name, action, args); err != nil ***REMOVED***
		return err
	***REMOVED***
	// reverse
	args[7], args[9] = args[9], args[7]
	args[10] = "--sport"
	return ProgramRule(Filter, c.Name, action, args)
***REMOVED***

// ProgramRule adds the rule specified by args only if the
// rule is not already present in the chain. Reciprocally,
// it removes the rule only if present.
func ProgramRule(table Table, chain string, action Action, args []string) error ***REMOVED***
	if Exists(table, chain, args...) != (action == Delete) ***REMOVED***
		return nil
	***REMOVED***
	return RawCombinedOutput(append([]string***REMOVED***"-t", string(table), string(action), chain***REMOVED***, args...)...)
***REMOVED***

// Prerouting adds linking rule to nat/PREROUTING chain.
func (c *ChainInfo) Prerouting(action Action, args ...string) error ***REMOVED***
	a := []string***REMOVED***"-t", string(Nat), string(action), "PREROUTING"***REMOVED***
	if len(args) > 0 ***REMOVED***
		a = append(a, args...)
	***REMOVED***
	if output, err := Raw(a...); err != nil ***REMOVED***
		return err
	***REMOVED*** else if len(output) != 0 ***REMOVED***
		return ChainError***REMOVED***Chain: "PREROUTING", Output: output***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Output adds linking rule to an OUTPUT chain.
func (c *ChainInfo) Output(action Action, args ...string) error ***REMOVED***
	a := []string***REMOVED***"-t", string(c.Table), string(action), "OUTPUT"***REMOVED***
	if len(args) > 0 ***REMOVED***
		a = append(a, args...)
	***REMOVED***
	if output, err := Raw(a...); err != nil ***REMOVED***
		return err
	***REMOVED*** else if len(output) != 0 ***REMOVED***
		return ChainError***REMOVED***Chain: "OUTPUT", Output: output***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Remove removes the chain.
func (c *ChainInfo) Remove() error ***REMOVED***
	// Ignore errors - This could mean the chains were never set up
	if c.Table == Nat ***REMOVED***
		c.Prerouting(Delete, "-m", "addrtype", "--dst-type", "LOCAL", "-j", c.Name)
		c.Output(Delete, "-m", "addrtype", "--dst-type", "LOCAL", "!", "--dst", "127.0.0.0/8", "-j", c.Name)
		c.Output(Delete, "-m", "addrtype", "--dst-type", "LOCAL", "-j", c.Name) // Created in versions <= 0.1.6

		c.Prerouting(Delete)
		c.Output(Delete)
	***REMOVED***
	Raw("-t", string(c.Table), "-F", c.Name)
	Raw("-t", string(c.Table), "-X", c.Name)
	return nil
***REMOVED***

// Exists checks if a rule exists
func Exists(table Table, chain string, rule ...string) bool ***REMOVED***
	return exists(false, table, chain, rule...)
***REMOVED***

// ExistsNative behaves as Exists with the difference it
// will always invoke `iptables` binary.
func ExistsNative(table Table, chain string, rule ...string) bool ***REMOVED***
	return exists(true, table, chain, rule...)
***REMOVED***

func exists(native bool, table Table, chain string, rule ...string) bool ***REMOVED***
	f := Raw
	if native ***REMOVED***
		f = raw
	***REMOVED***

	if string(table) == "" ***REMOVED***
		table = Filter
	***REMOVED***

	if err := initCheck(); err != nil ***REMOVED***
		// The exists() signature does not allow us to return an error, but at least
		// we can skip the (likely invalid) exec invocation.
		return false
	***REMOVED***

	if supportsCOpt ***REMOVED***
		// if exit status is 0 then return true, the rule exists
		_, err := f(append([]string***REMOVED***"-t", string(table), "-C", chain***REMOVED***, rule...)...)
		return err == nil
	***REMOVED***

	// parse "iptables -S" for the rule (it checks rules in a specific chain
	// in a specific table and it is very unreliable)
	return existsRaw(table, chain, rule...)
***REMOVED***

func existsRaw(table Table, chain string, rule ...string) bool ***REMOVED***
	ruleString := fmt.Sprintf("%s %s\n", chain, strings.Join(rule, " "))
	existingRules, _ := exec.Command(iptablesPath, "-t", string(table), "-S", chain).Output()

	return strings.Contains(string(existingRules), ruleString)
***REMOVED***

// Raw calls 'iptables' system command, passing supplied arguments.
func Raw(args ...string) ([]byte, error) ***REMOVED***
	if firewalldRunning ***REMOVED***
		output, err := Passthrough(Iptables, args...)
		if err == nil || !strings.Contains(err.Error(), "was not provided by any .service files") ***REMOVED***
			return output, err
		***REMOVED***
	***REMOVED***
	return raw(args...)
***REMOVED***

func raw(args ...string) ([]byte, error) ***REMOVED***
	if err := initCheck(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if supportsXlock ***REMOVED***
		args = append([]string***REMOVED***"--wait"***REMOVED***, args...)
	***REMOVED*** else ***REMOVED***
		bestEffortLock.Lock()
		defer bestEffortLock.Unlock()
	***REMOVED***

	logrus.Debugf("%s, %v", iptablesPath, args)

	output, err := exec.Command(iptablesPath, args...).CombinedOutput()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("iptables failed: iptables %v: %s (%s)", strings.Join(args, " "), output, err)
	***REMOVED***

	// ignore iptables' message about xtables lock
	if strings.Contains(string(output), xLockWaitMsg) ***REMOVED***
		output = []byte("")
	***REMOVED***

	return output, err
***REMOVED***

// RawCombinedOutput inernally calls the Raw function and returns a non nil
// error if Raw returned a non nil error or a non empty output
func RawCombinedOutput(args ...string) error ***REMOVED***
	if output, err := Raw(args...); err != nil || len(output) != 0 ***REMOVED***
		return fmt.Errorf("%s (%v)", string(output), err)
	***REMOVED***
	return nil
***REMOVED***

// RawCombinedOutputNative behave as RawCombinedOutput with the difference it
// will always invoke `iptables` binary
func RawCombinedOutputNative(args ...string) error ***REMOVED***
	if output, err := raw(args...); err != nil || len(output) != 0 ***REMOVED***
		return fmt.Errorf("%s (%v)", string(output), err)
	***REMOVED***
	return nil
***REMOVED***

// ExistChain checks if a chain exists
func ExistChain(chain string, table Table) bool ***REMOVED***
	if _, err := Raw("-t", string(table), "-nL", chain); err == nil ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// GetVersion reads the iptables version numbers during initialization
func GetVersion() (major, minor, micro int, err error) ***REMOVED***
	out, err := exec.Command(iptablesPath, "--version").CombinedOutput()
	if err == nil ***REMOVED***
		major, minor, micro = parseVersionNumbers(string(out))
	***REMOVED***
	return
***REMOVED***

// SetDefaultPolicy sets the passed default policy for the table/chain
func SetDefaultPolicy(table Table, chain string, policy Policy) error ***REMOVED***
	if err := RawCombinedOutput("-t", string(table), "-P", chain, string(policy)); err != nil ***REMOVED***
		return fmt.Errorf("setting default policy to %v in %v chain failed: %v", policy, chain, err)
	***REMOVED***
	return nil
***REMOVED***

func parseVersionNumbers(input string) (major, minor, micro int) ***REMOVED***
	re := regexp.MustCompile(`v\d*.\d*.\d*`)
	line := re.FindString(input)
	fmt.Sscanf(line, "v%d.%d.%d", &major, &minor, &micro)
	return
***REMOVED***

// iptables -C, --check option was added in v.1.4.11
// http://ftp.netfilter.org/pub/iptables/changes-iptables-1.4.11.txt
func supportsCOption(mj, mn, mc int) bool ***REMOVED***
	return mj > 1 || (mj == 1 && (mn > 4 || (mn == 4 && mc >= 11)))
***REMOVED***

// AddReturnRule adds a return rule for the chain in the filter table
func AddReturnRule(chain string) error ***REMOVED***
	var (
		table = Filter
		args  = []string***REMOVED***"-j", "RETURN"***REMOVED***
	)

	if Exists(table, chain, args...) ***REMOVED***
		return nil
	***REMOVED***

	err := RawCombinedOutput(append([]string***REMOVED***"-A", chain***REMOVED***, args...)...)
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to add return rule in %s chain: %s", chain, err.Error())
	***REMOVED***

	return nil
***REMOVED***

// EnsureJumpRule ensures the jump rule is on top
func EnsureJumpRule(fromChain, toChain string) error ***REMOVED***
	var (
		table = Filter
		args  = []string***REMOVED***"-j", toChain***REMOVED***
	)

	if Exists(table, fromChain, args...) ***REMOVED***
		err := RawCombinedOutput(append([]string***REMOVED***"-D", fromChain***REMOVED***, args...)...)
		if err != nil ***REMOVED***
			return fmt.Errorf("unable to remove jump to %s rule in %s chain: %s", toChain, fromChain, err.Error())
		***REMOVED***
	***REMOVED***

	err := RawCombinedOutput(append([]string***REMOVED***"-I", fromChain***REMOVED***, args...)...)
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to insert jump to %s rule in %s chain: %s", toChain, fromChain, err.Error())
	***REMOVED***

	return nil
***REMOVED***
