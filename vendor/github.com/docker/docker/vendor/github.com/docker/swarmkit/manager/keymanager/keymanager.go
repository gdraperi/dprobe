package keymanager

// keymanager does the allocation, rotation and distribution of symmetric
// keys to the agents. This is to securely bootstrap network communication
// between agents. It can be used for encrypting gossip between the agents
// which is used to exchange service discovery and overlay network control
// plane information. It can also be used to encrypt overlay data traffic.
import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"sync"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	// DefaultKeyLen is the default length (in bytes) of the key allocated
	DefaultKeyLen = 16

	// DefaultKeyRotationInterval used by key manager
	DefaultKeyRotationInterval = 12 * time.Hour

	// SubsystemGossip handles gossip protocol between the agents
	SubsystemGossip = "networking:gossip"

	// SubsystemIPSec is overlay network data encryption subsystem
	SubsystemIPSec = "networking:ipsec"

	// DefaultSubsystem is gossip
	DefaultSubsystem = SubsystemGossip
	// number of keys to mainrain in the key ring.
	keyringSize = 3
)

// map of subsystems and corresponding encryption algorithm. Initially only
// AES_128 in GCM mode is supported.
var subsysToAlgo = map[string]api.EncryptionKey_Algorithm***REMOVED***
	SubsystemGossip: api.AES_128_GCM,
	SubsystemIPSec:  api.AES_128_GCM,
***REMOVED***

type keyRing struct ***REMOVED***
	lClock uint64
	keys   []*api.EncryptionKey
***REMOVED***

// Config for the keymanager that can be modified
type Config struct ***REMOVED***
	ClusterName      string
	Keylen           int
	RotationInterval time.Duration
	Subsystems       []string
***REMOVED***

// KeyManager handles key allocation, rotation & distribution
type KeyManager struct ***REMOVED***
	config  *Config
	store   *store.MemoryStore
	keyRing *keyRing
	ctx     context.Context
	cancel  context.CancelFunc

	mu sync.Mutex
***REMOVED***

// DefaultConfig provides the default config for keymanager
func DefaultConfig() *Config ***REMOVED***
	return &Config***REMOVED***
		ClusterName:      store.DefaultClusterName,
		Keylen:           DefaultKeyLen,
		RotationInterval: DefaultKeyRotationInterval,
		Subsystems:       []string***REMOVED***SubsystemGossip, SubsystemIPSec***REMOVED***,
	***REMOVED***
***REMOVED***

// New creates an instance of keymanager with the given config
func New(store *store.MemoryStore, config *Config) *KeyManager ***REMOVED***
	for _, subsys := range config.Subsystems ***REMOVED***
		if subsys != SubsystemGossip && subsys != SubsystemIPSec ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return &KeyManager***REMOVED***
		config:  config,
		store:   store,
		keyRing: &keyRing***REMOVED***lClock: genSkew()***REMOVED***,
	***REMOVED***
***REMOVED***

func (k *KeyManager) allocateKey(ctx context.Context, subsys string) *api.EncryptionKey ***REMOVED***
	key := make([]byte, k.config.Keylen)

	_, err := cryptorand.Read(key)
	if err != nil ***REMOVED***
		panic(errors.Wrap(err, "key generated failed"))
	***REMOVED***
	k.keyRing.lClock++

	return &api.EncryptionKey***REMOVED***
		Subsystem:   subsys,
		Algorithm:   subsysToAlgo[subsys],
		Key:         key,
		LamportTime: k.keyRing.lClock,
	***REMOVED***
***REMOVED***

func (k *KeyManager) updateKey(cluster *api.Cluster) error ***REMOVED***
	return k.store.Update(func(tx store.Tx) error ***REMOVED***
		cluster = store.GetCluster(tx, cluster.ID)
		if cluster == nil ***REMOVED***
			return nil
		***REMOVED***
		cluster.EncryptionKeyLamportClock = k.keyRing.lClock
		cluster.NetworkBootstrapKeys = k.keyRing.keys
		return store.UpdateCluster(tx, cluster)
	***REMOVED***)
***REMOVED***

func (k *KeyManager) rotateKey(ctx context.Context) error ***REMOVED***
	var (
		clusters []*api.Cluster
		err      error
	)
	k.store.View(func(readTx store.ReadTx) ***REMOVED***
		clusters, err = store.FindClusters(readTx, store.ByName(k.config.ClusterName))
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).Errorf("reading cluster config failed, %v", err)
		return err
	***REMOVED***

	cluster := clusters[0]
	if len(cluster.NetworkBootstrapKeys) == 0 ***REMOVED***
		panic(errors.New("no key in the cluster config"))
	***REMOVED***

	subsysKeys := map[string][]*api.EncryptionKey***REMOVED******REMOVED***
	for _, key := range k.keyRing.keys ***REMOVED***
		subsysKeys[key.Subsystem] = append(subsysKeys[key.Subsystem], key)
	***REMOVED***
	k.keyRing.keys = []*api.EncryptionKey***REMOVED******REMOVED***

	// We maintain the latest key and the one before in the key ring to allow
	// agents to communicate without disruption on key change.
	for subsys, keys := range subsysKeys ***REMOVED***
		if len(keys) == keyringSize ***REMOVED***
			min := 0
			for i, key := range keys[1:] ***REMOVED***
				if key.LamportTime < keys[min].LamportTime ***REMOVED***
					min = i
				***REMOVED***
			***REMOVED***
			keys = append(keys[0:min], keys[min+1:]...)
		***REMOVED***
		keys = append(keys, k.allocateKey(ctx, subsys))
		subsysKeys[subsys] = keys
	***REMOVED***

	for _, keys := range subsysKeys ***REMOVED***
		k.keyRing.keys = append(k.keyRing.keys, keys...)
	***REMOVED***

	return k.updateKey(cluster)
***REMOVED***

// Run starts the keymanager, it doesn't return
func (k *KeyManager) Run(ctx context.Context) error ***REMOVED***
	k.mu.Lock()
	ctx = log.WithModule(ctx, "keymanager")
	var (
		clusters []*api.Cluster
		err      error
	)
	k.store.View(func(readTx store.ReadTx) ***REMOVED***
		clusters, err = store.FindClusters(readTx, store.ByName(k.config.ClusterName))
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).Errorf("reading cluster config failed, %v", err)
		k.mu.Unlock()
		return err
	***REMOVED***

	cluster := clusters[0]
	if len(cluster.NetworkBootstrapKeys) == 0 ***REMOVED***
		for _, subsys := range k.config.Subsystems ***REMOVED***
			for i := 0; i < keyringSize; i++ ***REMOVED***
				k.keyRing.keys = append(k.keyRing.keys, k.allocateKey(ctx, subsys))
			***REMOVED***
		***REMOVED***
		if err := k.updateKey(cluster); err != nil ***REMOVED***
			log.G(ctx).Errorf("store update failed %v", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		k.keyRing.lClock = cluster.EncryptionKeyLamportClock
		k.keyRing.keys = cluster.NetworkBootstrapKeys
	***REMOVED***

	ticker := time.NewTicker(k.config.RotationInterval)
	defer ticker.Stop()

	k.ctx, k.cancel = context.WithCancel(ctx)
	k.mu.Unlock()

	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			k.rotateKey(ctx)
		case <-k.ctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops the running instance of key manager
func (k *KeyManager) Stop() error ***REMOVED***
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.cancel == nil ***REMOVED***
		return errors.New("keymanager is not started")
	***REMOVED***
	k.cancel()
	return nil
***REMOVED***

// genSkew generates a random uint64 number between 0 and 65535
func genSkew() uint64 ***REMOVED***
	b := make([]byte, 2)
	if _, err := cryptorand.Read(b); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return uint64(binary.BigEndian.Uint16(b))
***REMOVED***
