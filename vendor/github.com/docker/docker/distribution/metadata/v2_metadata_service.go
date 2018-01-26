package metadata

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/layer"
	"github.com/opencontainers/go-digest"
)

// V2MetadataService maps layer IDs to a set of known metadata for
// the layer.
type V2MetadataService interface ***REMOVED***
	GetMetadata(diffID layer.DiffID) ([]V2Metadata, error)
	GetDiffID(dgst digest.Digest) (layer.DiffID, error)
	Add(diffID layer.DiffID, metadata V2Metadata) error
	TagAndAdd(diffID layer.DiffID, hmacKey []byte, metadata V2Metadata) error
	Remove(metadata V2Metadata) error
***REMOVED***

// v2MetadataService implements V2MetadataService
type v2MetadataService struct ***REMOVED***
	store Store
***REMOVED***

var _ V2MetadataService = &v2MetadataService***REMOVED******REMOVED***

// V2Metadata contains the digest and source repository information for a layer.
type V2Metadata struct ***REMOVED***
	Digest           digest.Digest
	SourceRepository string
	// HMAC hashes above attributes with recent authconfig digest used as a key in order to determine matching
	// metadata entries accompanied by the same credentials without actually exposing them.
	HMAC string
***REMOVED***

// CheckV2MetadataHMAC returns true if the given "meta" is tagged with a hmac hashed by the given "key".
func CheckV2MetadataHMAC(meta *V2Metadata, key []byte) bool ***REMOVED***
	if len(meta.HMAC) == 0 || len(key) == 0 ***REMOVED***
		return len(meta.HMAC) == 0 && len(key) == 0
	***REMOVED***
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(meta.Digest))
	mac.Write([]byte(meta.SourceRepository))
	expectedMac := mac.Sum(nil)

	storedMac, err := hex.DecodeString(meta.HMAC)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	return hmac.Equal(storedMac, expectedMac)
***REMOVED***

// ComputeV2MetadataHMAC returns a hmac for the given "meta" hash by the given key.
func ComputeV2MetadataHMAC(key []byte, meta *V2Metadata) string ***REMOVED***
	if len(key) == 0 || meta == nil ***REMOVED***
		return ""
	***REMOVED***
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(meta.Digest))
	mac.Write([]byte(meta.SourceRepository))
	return hex.EncodeToString(mac.Sum(nil))
***REMOVED***

// ComputeV2MetadataHMACKey returns a key for the given "authConfig" that can be used to hash v2 metadata
// entries.
func ComputeV2MetadataHMACKey(authConfig *types.AuthConfig) ([]byte, error) ***REMOVED***
	if authConfig == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	key := authConfigKeyInput***REMOVED***
		Username:      authConfig.Username,
		Password:      authConfig.Password,
		Auth:          authConfig.Auth,
		IdentityToken: authConfig.IdentityToken,
		RegistryToken: authConfig.RegistryToken,
	***REMOVED***
	buf, err := json.Marshal(&key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return []byte(digest.FromBytes(buf)), nil
***REMOVED***

// authConfigKeyInput is a reduced AuthConfig structure holding just relevant credential data eligible for
// hmac key creation.
type authConfigKeyInput struct ***REMOVED***
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`

	IdentityToken string `json:"identitytoken,omitempty"`
	RegistryToken string `json:"registrytoken,omitempty"`
***REMOVED***

// maxMetadata is the number of metadata entries to keep per layer DiffID.
const maxMetadata = 50

// NewV2MetadataService creates a new diff ID to v2 metadata mapping service.
func NewV2MetadataService(store Store) V2MetadataService ***REMOVED***
	return &v2MetadataService***REMOVED***
		store: store,
	***REMOVED***
***REMOVED***

func (serv *v2MetadataService) diffIDNamespace() string ***REMOVED***
	return "v2metadata-by-diffid"
***REMOVED***

func (serv *v2MetadataService) digestNamespace() string ***REMOVED***
	return "diffid-by-digest"
***REMOVED***

func (serv *v2MetadataService) diffIDKey(diffID layer.DiffID) string ***REMOVED***
	return string(digest.Digest(diffID).Algorithm()) + "/" + digest.Digest(diffID).Hex()
***REMOVED***

func (serv *v2MetadataService) digestKey(dgst digest.Digest) string ***REMOVED***
	return string(dgst.Algorithm()) + "/" + dgst.Hex()
***REMOVED***

// GetMetadata finds the metadata associated with a layer DiffID.
func (serv *v2MetadataService) GetMetadata(diffID layer.DiffID) ([]V2Metadata, error) ***REMOVED***
	if serv.store == nil ***REMOVED***
		return nil, errors.New("no metadata storage")
	***REMOVED***
	jsonBytes, err := serv.store.Get(serv.diffIDNamespace(), serv.diffIDKey(diffID))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var metadata []V2Metadata
	if err := json.Unmarshal(jsonBytes, &metadata); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return metadata, nil
***REMOVED***

// GetDiffID finds a layer DiffID from a digest.
func (serv *v2MetadataService) GetDiffID(dgst digest.Digest) (layer.DiffID, error) ***REMOVED***
	if serv.store == nil ***REMOVED***
		return layer.DiffID(""), errors.New("no metadata storage")
	***REMOVED***
	diffIDBytes, err := serv.store.Get(serv.digestNamespace(), serv.digestKey(dgst))
	if err != nil ***REMOVED***
		return layer.DiffID(""), err
	***REMOVED***

	return layer.DiffID(diffIDBytes), nil
***REMOVED***

// Add associates metadata with a layer DiffID. If too many metadata entries are
// present, the oldest one is dropped.
func (serv *v2MetadataService) Add(diffID layer.DiffID, metadata V2Metadata) error ***REMOVED***
	if serv.store == nil ***REMOVED***
		// Support a service which has no backend storage, in this case
		// an add becomes a no-op.
		// TODO: implement in memory storage
		return nil
	***REMOVED***
	oldMetadata, err := serv.GetMetadata(diffID)
	if err != nil ***REMOVED***
		oldMetadata = nil
	***REMOVED***
	newMetadata := make([]V2Metadata, 0, len(oldMetadata)+1)

	// Copy all other metadata to new slice
	for _, oldMeta := range oldMetadata ***REMOVED***
		if oldMeta != metadata ***REMOVED***
			newMetadata = append(newMetadata, oldMeta)
		***REMOVED***
	***REMOVED***

	newMetadata = append(newMetadata, metadata)

	if len(newMetadata) > maxMetadata ***REMOVED***
		newMetadata = newMetadata[len(newMetadata)-maxMetadata:]
	***REMOVED***

	jsonBytes, err := json.Marshal(newMetadata)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = serv.store.Set(serv.diffIDNamespace(), serv.diffIDKey(diffID), jsonBytes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return serv.store.Set(serv.digestNamespace(), serv.digestKey(metadata.Digest), []byte(diffID))
***REMOVED***

// TagAndAdd amends the given "meta" for hmac hashed by the given "hmacKey" and associates it with a layer
// DiffID. If too many metadata entries are present, the oldest one is dropped.
func (serv *v2MetadataService) TagAndAdd(diffID layer.DiffID, hmacKey []byte, meta V2Metadata) error ***REMOVED***
	meta.HMAC = ComputeV2MetadataHMAC(hmacKey, &meta)
	return serv.Add(diffID, meta)
***REMOVED***

// Remove disassociates a metadata entry from a layer DiffID.
func (serv *v2MetadataService) Remove(metadata V2Metadata) error ***REMOVED***
	if serv.store == nil ***REMOVED***
		// Support a service which has no backend storage, in this case
		// an remove becomes a no-op.
		// TODO: implement in memory storage
		return nil
	***REMOVED***
	diffID, err := serv.GetDiffID(metadata.Digest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	oldMetadata, err := serv.GetMetadata(diffID)
	if err != nil ***REMOVED***
		oldMetadata = nil
	***REMOVED***
	newMetadata := make([]V2Metadata, 0, len(oldMetadata))

	// Copy all other metadata to new slice
	for _, oldMeta := range oldMetadata ***REMOVED***
		if oldMeta != metadata ***REMOVED***
			newMetadata = append(newMetadata, oldMeta)
		***REMOVED***
	***REMOVED***

	if len(newMetadata) == 0 ***REMOVED***
		return serv.store.Delete(serv.diffIDNamespace(), serv.diffIDKey(diffID))
	***REMOVED***

	jsonBytes, err := json.Marshal(newMetadata)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return serv.store.Set(serv.diffIDNamespace(), serv.diffIDKey(diffID), jsonBytes)
***REMOVED***
