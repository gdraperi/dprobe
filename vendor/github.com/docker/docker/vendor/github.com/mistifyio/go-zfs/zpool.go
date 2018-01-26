package zfs

// ZFS zpool states, which can indicate if a pool is online, offline,
// degraded, etc.  More information regarding zpool states can be found here:
// https://docs.oracle.com/cd/E19253-01/819-5461/gamno/index.html.
const (
	ZpoolOnline   = "ONLINE"
	ZpoolDegraded = "DEGRADED"
	ZpoolFaulted  = "FAULTED"
	ZpoolOffline  = "OFFLINE"
	ZpoolUnavail  = "UNAVAIL"
	ZpoolRemoved  = "REMOVED"
)

// Zpool is a ZFS zpool.  A pool is a top-level structure in ZFS, and can
// contain many descendent datasets.
type Zpool struct ***REMOVED***
	Name          string
	Health        string
	Allocated     uint64
	Size          uint64
	Free          uint64
	Fragmentation uint64
	ReadOnly      bool
	Freeing       uint64
	Leaked        uint64
	DedupRatio    float64
***REMOVED***

// zpool is a helper function to wrap typical calls to zpool.
func zpool(arg ...string) ([][]string, error) ***REMOVED***
	c := command***REMOVED***Command: "zpool"***REMOVED***
	return c.Run(arg...)
***REMOVED***

// GetZpool retrieves a single ZFS zpool by name.
func GetZpool(name string) (*Zpool, error) ***REMOVED***
	args := zpoolArgs
	args = append(args, name)
	out, err := zpool(args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// there is no -H
	out = out[1:]

	z := &Zpool***REMOVED***Name: name***REMOVED***
	for _, line := range out ***REMOVED***
		if err := z.parseLine(line); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return z, nil
***REMOVED***

// Datasets returns a slice of all ZFS datasets in a zpool.
func (z *Zpool) Datasets() ([]*Dataset, error) ***REMOVED***
	return Datasets(z.Name)
***REMOVED***

// Snapshots returns a slice of all ZFS snapshots in a zpool.
func (z *Zpool) Snapshots() ([]*Dataset, error) ***REMOVED***
	return Snapshots(z.Name)
***REMOVED***

// CreateZpool creates a new ZFS zpool with the specified name, properties,
// and optional arguments.
// A full list of available ZFS properties and command-line arguments may be
// found here: https://www.freebsd.org/cgi/man.cgi?zfs(8).
func CreateZpool(name string, properties map[string]string, args ...string) (*Zpool, error) ***REMOVED***
	cli := make([]string, 1, 4)
	cli[0] = "create"
	if properties != nil ***REMOVED***
		cli = append(cli, propsSlice(properties)...)
	***REMOVED***
	cli = append(cli, name)
	cli = append(cli, args...)
	_, err := zpool(cli...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Zpool***REMOVED***Name: name***REMOVED***, nil
***REMOVED***

// Destroy destroys a ZFS zpool by name.
func (z *Zpool) Destroy() error ***REMOVED***
	_, err := zpool("destroy", z.Name)
	return err
***REMOVED***

// ListZpools list all ZFS zpools accessible on the current system.
func ListZpools() ([]*Zpool, error) ***REMOVED***
	args := []string***REMOVED***"list", "-Ho", "name"***REMOVED***
	out, err := zpool(args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var pools []*Zpool

	for _, line := range out ***REMOVED***
		z, err := GetZpool(line[0])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		pools = append(pools, z)
	***REMOVED***
	return pools, nil
***REMOVED***
