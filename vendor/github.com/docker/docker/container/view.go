package container

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-memdb"
	"github.com/sirupsen/logrus"
)

const (
	memdbContainersTable  = "containers"
	memdbNamesTable       = "names"
	memdbIDIndex          = "id"
	memdbContainerIDIndex = "containerid"
)

var (
	// ErrNameReserved is an error which is returned when a name is requested to be reserved that already is reserved
	ErrNameReserved = errors.New("name is reserved")
	// ErrNameNotReserved is an error which is returned when trying to find a name that is not reserved
	ErrNameNotReserved = errors.New("name is not reserved")
)

// Snapshot is a read only view for Containers. It holds all information necessary to serve container queries in a
// versioned ACID in-memory store.
type Snapshot struct ***REMOVED***
	types.Container

	// additional info queries need to filter on
	// preserve nanosec resolution for queries
	CreatedAt    time.Time
	StartedAt    time.Time
	Name         string
	Pid          int
	ExitCode     int
	Running      bool
	Paused       bool
	Managed      bool
	ExposedPorts nat.PortSet
	PortBindings nat.PortSet
	Health       string
	HostConfig   struct ***REMOVED***
		Isolation string
	***REMOVED***
***REMOVED***

// nameAssociation associates a container id with a name.
type nameAssociation struct ***REMOVED***
	// name is the name to associate. Note that name is the primary key
	// ("id" in memdb).
	name        string
	containerID string
***REMOVED***

// ViewDB provides an in-memory transactional (ACID) container Store
type ViewDB interface ***REMOVED***
	Snapshot() View
	Save(*Container) error
	Delete(*Container) error

	ReserveName(name, containerID string) error
	ReleaseName(name string) error
***REMOVED***

// View can be used by readers to avoid locking
type View interface ***REMOVED***
	All() ([]Snapshot, error)
	Get(id string) (*Snapshot, error)

	GetID(name string) (string, error)
	GetAllNames() map[string][]string
***REMOVED***

var schema = &memdb.DBSchema***REMOVED***
	Tables: map[string]*memdb.TableSchema***REMOVED***
		memdbContainersTable: ***REMOVED***
			Name: memdbContainersTable,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				memdbIDIndex: ***REMOVED***
					Name:    memdbIDIndex,
					Unique:  true,
					Indexer: &containerByIDIndexer***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		memdbNamesTable: ***REMOVED***
			Name: memdbNamesTable,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				// Used for names, because "id" is the primary key in memdb.
				memdbIDIndex: ***REMOVED***
					Name:    memdbIDIndex,
					Unique:  true,
					Indexer: &namesByNameIndexer***REMOVED******REMOVED***,
				***REMOVED***,
				memdbContainerIDIndex: ***REMOVED***
					Name:    memdbContainerIDIndex,
					Indexer: &namesByContainerIDIndexer***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

type memDB struct ***REMOVED***
	store *memdb.MemDB
***REMOVED***

// NoSuchContainerError indicates that the container wasn't found in the
// database.
type NoSuchContainerError struct ***REMOVED***
	id string
***REMOVED***

// Error satisfies the error interface.
func (e NoSuchContainerError) Error() string ***REMOVED***
	return "no such container " + e.id
***REMOVED***

// NewViewDB provides the default implementation, with the default schema
func NewViewDB() (ViewDB, error) ***REMOVED***
	store, err := memdb.NewMemDB(schema)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &memDB***REMOVED***store: store***REMOVED***, nil
***REMOVED***

// Snapshot provides a consistent read-only View of the database
func (db *memDB) Snapshot() View ***REMOVED***
	return &memdbView***REMOVED***
		txn: db.store.Txn(false),
	***REMOVED***
***REMOVED***

func (db *memDB) withTxn(cb func(*memdb.Txn) error) error ***REMOVED***
	txn := db.store.Txn(true)
	err := cb(txn)
	if err != nil ***REMOVED***
		txn.Abort()
		return err
	***REMOVED***
	txn.Commit()
	return nil
***REMOVED***

// Save atomically updates the in-memory store state for a Container.
// Only read only (deep) copies of containers may be passed in.
func (db *memDB) Save(c *Container) error ***REMOVED***
	return db.withTxn(func(txn *memdb.Txn) error ***REMOVED***
		return txn.Insert(memdbContainersTable, c)
	***REMOVED***)
***REMOVED***

// Delete removes an item by ID
func (db *memDB) Delete(c *Container) error ***REMOVED***
	return db.withTxn(func(txn *memdb.Txn) error ***REMOVED***
		view := &memdbView***REMOVED***txn: txn***REMOVED***
		names := view.getNames(c.ID)

		for _, name := range names ***REMOVED***
			txn.Delete(memdbNamesTable, nameAssociation***REMOVED***name: name***REMOVED***)
		***REMOVED***

		// Ignore error - the container may not actually exist in the
		// db, but we still need to clean up associated names.
		txn.Delete(memdbContainersTable, NewBaseContainer(c.ID, c.Root))
		return nil
	***REMOVED***)
***REMOVED***

// ReserveName registers a container ID to a name
// ReserveName is idempotent
// Attempting to reserve a container ID to a name that already exists results in an `ErrNameReserved`
// A name reservation is globally unique
func (db *memDB) ReserveName(name, containerID string) error ***REMOVED***
	return db.withTxn(func(txn *memdb.Txn) error ***REMOVED***
		s, err := txn.First(memdbNamesTable, memdbIDIndex, name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if s != nil ***REMOVED***
			if s.(nameAssociation).containerID != containerID ***REMOVED***
				return ErrNameReserved
			***REMOVED***
			return nil
		***REMOVED***
		return txn.Insert(memdbNamesTable, nameAssociation***REMOVED***name: name, containerID: containerID***REMOVED***)
	***REMOVED***)
***REMOVED***

// ReleaseName releases the reserved name
// Once released, a name can be reserved again
func (db *memDB) ReleaseName(name string) error ***REMOVED***
	return db.withTxn(func(txn *memdb.Txn) error ***REMOVED***
		return txn.Delete(memdbNamesTable, nameAssociation***REMOVED***name: name***REMOVED***)
	***REMOVED***)
***REMOVED***

type memdbView struct ***REMOVED***
	txn *memdb.Txn
***REMOVED***

// All returns a all items in this snapshot. Returned objects must never be modified.
func (v *memdbView) All() ([]Snapshot, error) ***REMOVED***
	var all []Snapshot
	iter, err := v.txn.Get(memdbContainersTable, memdbIDIndex)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for ***REMOVED***
		item := iter.Next()
		if item == nil ***REMOVED***
			break
		***REMOVED***
		snapshot := v.transform(item.(*Container))
		all = append(all, *snapshot)
	***REMOVED***
	return all, nil
***REMOVED***

// Get returns an item by id. Returned objects must never be modified.
func (v *memdbView) Get(id string) (*Snapshot, error) ***REMOVED***
	s, err := v.txn.First(memdbContainersTable, memdbIDIndex, id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if s == nil ***REMOVED***
		return nil, NoSuchContainerError***REMOVED***id: id***REMOVED***
	***REMOVED***
	return v.transform(s.(*Container)), nil
***REMOVED***

// getNames lists all the reserved names for the given container ID.
func (v *memdbView) getNames(containerID string) []string ***REMOVED***
	iter, err := v.txn.Get(memdbNamesTable, memdbContainerIDIndex, containerID)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	var names []string
	for ***REMOVED***
		item := iter.Next()
		if item == nil ***REMOVED***
			break
		***REMOVED***
		names = append(names, item.(nameAssociation).name)
	***REMOVED***

	return names
***REMOVED***

// GetID returns the container ID that the passed in name is reserved to.
func (v *memdbView) GetID(name string) (string, error) ***REMOVED***
	s, err := v.txn.First(memdbNamesTable, memdbIDIndex, name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if s == nil ***REMOVED***
		return "", ErrNameNotReserved
	***REMOVED***
	return s.(nameAssociation).containerID, nil
***REMOVED***

// GetAllNames returns all registered names.
func (v *memdbView) GetAllNames() map[string][]string ***REMOVED***
	iter, err := v.txn.Get(memdbNamesTable, memdbContainerIDIndex)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	out := make(map[string][]string)
	for ***REMOVED***
		item := iter.Next()
		if item == nil ***REMOVED***
			break
		***REMOVED***
		assoc := item.(nameAssociation)
		out[assoc.containerID] = append(out[assoc.containerID], assoc.name)
	***REMOVED***

	return out
***REMOVED***

// transform maps a (deep) copied Container object to what queries need.
// A lock on the Container is not held because these are immutable deep copies.
func (v *memdbView) transform(container *Container) *Snapshot ***REMOVED***
	health := types.NoHealthcheck
	if container.Health != nil ***REMOVED***
		health = container.Health.Status()
	***REMOVED***
	snapshot := &Snapshot***REMOVED***
		Container: types.Container***REMOVED***
			ID:      container.ID,
			Names:   v.getNames(container.ID),
			ImageID: container.ImageID.String(),
			Ports:   []types.Port***REMOVED******REMOVED***,
			Mounts:  container.GetMountPoints(),
			State:   container.State.StateString(),
			Status:  container.State.String(),
			Created: container.Created.Unix(),
		***REMOVED***,
		CreatedAt:    container.Created,
		StartedAt:    container.StartedAt,
		Name:         container.Name,
		Pid:          container.Pid,
		Managed:      container.Managed,
		ExposedPorts: make(nat.PortSet),
		PortBindings: make(nat.PortSet),
		Health:       health,
		Running:      container.Running,
		Paused:       container.Paused,
		ExitCode:     container.ExitCode(),
	***REMOVED***

	if snapshot.Names == nil ***REMOVED***
		// Dead containers will often have no name, so make sure the response isn't null
		snapshot.Names = []string***REMOVED******REMOVED***
	***REMOVED***

	if container.HostConfig != nil ***REMOVED***
		snapshot.Container.HostConfig.NetworkMode = string(container.HostConfig.NetworkMode)
		snapshot.HostConfig.Isolation = string(container.HostConfig.Isolation)
		for binding := range container.HostConfig.PortBindings ***REMOVED***
			snapshot.PortBindings[binding] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	if container.Config != nil ***REMOVED***
		snapshot.Image = container.Config.Image
		snapshot.Labels = container.Config.Labels
		for exposed := range container.Config.ExposedPorts ***REMOVED***
			snapshot.ExposedPorts[exposed] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	if len(container.Args) > 0 ***REMOVED***
		args := []string***REMOVED******REMOVED***
		for _, arg := range container.Args ***REMOVED***
			if strings.Contains(arg, " ") ***REMOVED***
				args = append(args, fmt.Sprintf("'%s'", arg))
			***REMOVED*** else ***REMOVED***
				args = append(args, arg)
			***REMOVED***
		***REMOVED***
		argsAsString := strings.Join(args, " ")
		snapshot.Command = fmt.Sprintf("%s %s", container.Path, argsAsString)
	***REMOVED*** else ***REMOVED***
		snapshot.Command = container.Path
	***REMOVED***

	snapshot.Ports = []types.Port***REMOVED******REMOVED***
	networks := make(map[string]*network.EndpointSettings)
	if container.NetworkSettings != nil ***REMOVED***
		for name, netw := range container.NetworkSettings.Networks ***REMOVED***
			if netw == nil || netw.EndpointSettings == nil ***REMOVED***
				continue
			***REMOVED***
			networks[name] = &network.EndpointSettings***REMOVED***
				EndpointID:          netw.EndpointID,
				Gateway:             netw.Gateway,
				IPAddress:           netw.IPAddress,
				IPPrefixLen:         netw.IPPrefixLen,
				IPv6Gateway:         netw.IPv6Gateway,
				GlobalIPv6Address:   netw.GlobalIPv6Address,
				GlobalIPv6PrefixLen: netw.GlobalIPv6PrefixLen,
				MacAddress:          netw.MacAddress,
				NetworkID:           netw.NetworkID,
			***REMOVED***
			if netw.IPAMConfig != nil ***REMOVED***
				networks[name].IPAMConfig = &network.EndpointIPAMConfig***REMOVED***
					IPv4Address: netw.IPAMConfig.IPv4Address,
					IPv6Address: netw.IPAMConfig.IPv6Address,
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for port, bindings := range container.NetworkSettings.Ports ***REMOVED***
			p, err := nat.ParsePort(port.Port())
			if err != nil ***REMOVED***
				logrus.Warnf("invalid port map %+v", err)
				continue
			***REMOVED***
			if len(bindings) == 0 ***REMOVED***
				snapshot.Ports = append(snapshot.Ports, types.Port***REMOVED***
					PrivatePort: uint16(p),
					Type:        port.Proto(),
				***REMOVED***)
				continue
			***REMOVED***
			for _, binding := range bindings ***REMOVED***
				h, err := nat.ParsePort(binding.HostPort)
				if err != nil ***REMOVED***
					logrus.Warnf("invalid host port map %+v", err)
					continue
				***REMOVED***
				snapshot.Ports = append(snapshot.Ports, types.Port***REMOVED***
					PrivatePort: uint16(p),
					PublicPort:  uint16(h),
					Type:        port.Proto(),
					IP:          binding.HostIP,
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	snapshot.NetworkSettings = &types.SummaryNetworkSettings***REMOVED***Networks: networks***REMOVED***

	return snapshot
***REMOVED***

// containerByIDIndexer is used to extract the ID field from Container types.
// memdb.StringFieldIndex can not be used since ID is a field from an embedded struct.
type containerByIDIndexer struct***REMOVED******REMOVED***

// FromObject implements the memdb.SingleIndexer interface for Container objects
func (e *containerByIDIndexer) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	c, ok := obj.(*Container)
	if !ok ***REMOVED***
		return false, nil, fmt.Errorf("%T is not a Container", obj)
	***REMOVED***
	// Add the null character as a terminator
	v := c.ID + "\x00"
	return true, []byte(v), nil
***REMOVED***

// FromArgs implements the memdb.Indexer interface
func (e *containerByIDIndexer) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	arg, ok := args[0].(string)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	***REMOVED***
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
***REMOVED***

// namesByNameIndexer is used to index container name associations by name.
type namesByNameIndexer struct***REMOVED******REMOVED***

func (e *namesByNameIndexer) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	n, ok := obj.(nameAssociation)
	if !ok ***REMOVED***
		return false, nil, fmt.Errorf(`%T does not have type "nameAssociation"`, obj)
	***REMOVED***

	// Add the null character as a terminator
	return true, []byte(n.name + "\x00"), nil
***REMOVED***

func (e *namesByNameIndexer) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	arg, ok := args[0].(string)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	***REMOVED***
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
***REMOVED***

// namesByContainerIDIndexer is used to index container names by container ID.
type namesByContainerIDIndexer struct***REMOVED******REMOVED***

func (e *namesByContainerIDIndexer) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	n, ok := obj.(nameAssociation)
	if !ok ***REMOVED***
		return false, nil, fmt.Errorf(`%T does not have type "nameAssocation"`, obj)
	***REMOVED***

	// Add the null character as a terminator
	return true, []byte(n.containerID + "\x00"), nil
***REMOVED***

func (e *namesByContainerIDIndexer) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	arg, ok := args[0].(string)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	***REMOVED***
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
***REMOVED***
