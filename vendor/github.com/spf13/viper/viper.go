// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Viper is a application configuration system.
// It believes that applications can be configured a variety of ways
// via flags, ENVIRONMENT variables, configuration files retrieved
// from the file system, or a remote key/value store.

// Each item takes precedence over the item below it:

// overrides
// flag
// env
// config
// key/value store
// default

package viper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/magiconair/properties"
	"github.com/mitchellh/mapstructure"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

// ConfigMarshalError happens when failing to marshal the configuration.
type ConfigMarshalError struct ***REMOVED***
	err error
***REMOVED***

// Error returns the formatted configuration error.
func (e ConfigMarshalError) Error() string ***REMOVED***
	return fmt.Sprintf("While marshaling config: %s", e.err.Error())
***REMOVED***

var v *Viper

type RemoteResponse struct ***REMOVED***
	Value []byte
	Error error
***REMOVED***

func init() ***REMOVED***
	v = New()
***REMOVED***

type remoteConfigFactory interface ***REMOVED***
	Get(rp RemoteProvider) (io.Reader, error)
	Watch(rp RemoteProvider) (io.Reader, error)
	WatchChannel(rp RemoteProvider) (<-chan *RemoteResponse, chan bool)
***REMOVED***

// RemoteConfig is optional, see the remote package
var RemoteConfig remoteConfigFactory

// UnsupportedConfigError denotes encountering an unsupported
// configuration filetype.
type UnsupportedConfigError string

// Error returns the formatted configuration error.
func (str UnsupportedConfigError) Error() string ***REMOVED***
	return fmt.Sprintf("Unsupported Config Type %q", string(str))
***REMOVED***

// UnsupportedRemoteProviderError denotes encountering an unsupported remote
// provider. Currently only etcd and Consul are supported.
type UnsupportedRemoteProviderError string

// Error returns the formatted remote provider error.
func (str UnsupportedRemoteProviderError) Error() string ***REMOVED***
	return fmt.Sprintf("Unsupported Remote Provider Type %q", string(str))
***REMOVED***

// RemoteConfigError denotes encountering an error while trying to
// pull the configuration from the remote provider.
type RemoteConfigError string

// Error returns the formatted remote provider error
func (rce RemoteConfigError) Error() string ***REMOVED***
	return fmt.Sprintf("Remote Configurations Error: %s", string(rce))
***REMOVED***

// ConfigFileNotFoundError denotes failing to find configuration file.
type ConfigFileNotFoundError struct ***REMOVED***
	name, locations string
***REMOVED***

// Error returns the formatted configuration error.
func (fnfe ConfigFileNotFoundError) Error() string ***REMOVED***
	return fmt.Sprintf("Config File %q Not Found in %q", fnfe.name, fnfe.locations)
***REMOVED***

// Viper is a prioritized configuration registry. It
// maintains a set of configuration sources, fetches
// values to populate those, and provides them according
// to the source's priority.
// The priority of the sources is the following:
// 1. overrides
// 2. flags
// 3. env. variables
// 4. config file
// 5. key/value store
// 6. defaults
//
// For example, if values from the following sources were loaded:
//
//  Defaults : ***REMOVED***
//  	"secret": "",
//  	"user": "default",
//  	"endpoint": "https://localhost"
//  ***REMOVED***
//  Config : ***REMOVED***
//  	"user": "root"
//  	"secret": "defaultsecret"
//  ***REMOVED***
//  Env : ***REMOVED***
//  	"secret": "somesecretkey"
//  ***REMOVED***
//
// The resulting config will have the following values:
//
//	***REMOVED***
//		"secret": "somesecretkey",
//		"user": "root",
//		"endpoint": "https://localhost"
//	***REMOVED***
type Viper struct ***REMOVED***
	// Delimiter that separates a list of keys
	// used to access a nested value in one go
	keyDelim string

	// A set of paths to look for the config file in
	configPaths []string

	// The filesystem to read config from.
	fs afero.Fs

	// A set of remote providers to search for the configuration
	remoteProviders []*defaultRemoteProvider

	// Name of file to look for inside the path
	configName string
	configFile string
	configType string
	envPrefix  string

	automaticEnvApplied bool
	envKeyReplacer      *strings.Replacer

	config         map[string]interface***REMOVED******REMOVED***
	override       map[string]interface***REMOVED******REMOVED***
	defaults       map[string]interface***REMOVED******REMOVED***
	kvstore        map[string]interface***REMOVED******REMOVED***
	pflags         map[string]FlagValue
	env            map[string]string
	aliases        map[string]string
	typeByDefValue bool

	// Store read properties on the object so that we can write back in order with comments.
	// This will only be used if the configuration read is a properties file.
	properties *properties.Properties

	onConfigChange func(fsnotify.Event)
***REMOVED***

// New returns an initialized Viper instance.
func New() *Viper ***REMOVED***
	v := new(Viper)
	v.keyDelim = "."
	v.configName = "config"
	v.fs = afero.NewOsFs()
	v.config = make(map[string]interface***REMOVED******REMOVED***)
	v.override = make(map[string]interface***REMOVED******REMOVED***)
	v.defaults = make(map[string]interface***REMOVED******REMOVED***)
	v.kvstore = make(map[string]interface***REMOVED******REMOVED***)
	v.pflags = make(map[string]FlagValue)
	v.env = make(map[string]string)
	v.aliases = make(map[string]string)
	v.typeByDefValue = false

	return v
***REMOVED***

// Intended for testing, will reset all to default settings.
// In the public interface for the viper package so applications
// can use it in their testing as well.
func Reset() ***REMOVED***
	v = New()
	SupportedExts = []string***REMOVED***"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl"***REMOVED***
	SupportedRemoteProviders = []string***REMOVED***"etcd", "consul"***REMOVED***
***REMOVED***

type defaultRemoteProvider struct ***REMOVED***
	provider      string
	endpoint      string
	path          string
	secretKeyring string
***REMOVED***

func (rp defaultRemoteProvider) Provider() string ***REMOVED***
	return rp.provider
***REMOVED***

func (rp defaultRemoteProvider) Endpoint() string ***REMOVED***
	return rp.endpoint
***REMOVED***

func (rp defaultRemoteProvider) Path() string ***REMOVED***
	return rp.path
***REMOVED***

func (rp defaultRemoteProvider) SecretKeyring() string ***REMOVED***
	return rp.secretKeyring
***REMOVED***

// RemoteProvider stores the configuration necessary
// to connect to a remote key/value store.
// Optional secretKeyring to unencrypt encrypted values
// can be provided.
type RemoteProvider interface ***REMOVED***
	Provider() string
	Endpoint() string
	Path() string
	SecretKeyring() string
***REMOVED***

// SupportedExts are universally supported extensions.
var SupportedExts = []string***REMOVED***"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl"***REMOVED***

// SupportedRemoteProviders are universally supported remote providers.
var SupportedRemoteProviders = []string***REMOVED***"etcd", "consul"***REMOVED***

func OnConfigChange(run func(in fsnotify.Event)) ***REMOVED*** v.OnConfigChange(run) ***REMOVED***
func (v *Viper) OnConfigChange(run func(in fsnotify.Event)) ***REMOVED***
	v.onConfigChange = run
***REMOVED***

func WatchConfig() ***REMOVED*** v.WatchConfig() ***REMOVED***
func (v *Viper) WatchConfig() ***REMOVED***
	go func() ***REMOVED***
		watcher, err := fsnotify.NewWatcher()
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		defer watcher.Close()

		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		filename, err := v.getConfigFile()
		if err != nil ***REMOVED***
			log.Println("error:", err)
			return
		***REMOVED***

		configFile := filepath.Clean(filename)
		configDir, _ := filepath.Split(configFile)

		done := make(chan bool)
		go func() ***REMOVED***
			for ***REMOVED***
				select ***REMOVED***
				case event := <-watcher.Events:
					// we only care about the config file
					if filepath.Clean(event.Name) == configFile ***REMOVED***
						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create ***REMOVED***
							err := v.ReadInConfig()
							if err != nil ***REMOVED***
								log.Println("error:", err)
							***REMOVED***
							v.onConfigChange(event)
						***REMOVED***
					***REMOVED***
				case err := <-watcher.Errors:
					log.Println("error:", err)
				***REMOVED***
			***REMOVED***
		***REMOVED***()

		watcher.Add(configDir)
		<-done
	***REMOVED***()
***REMOVED***

// SetConfigFile explicitly defines the path, name and extension of the config file.
// Viper will use this and not check any of the config paths.
func SetConfigFile(in string) ***REMOVED*** v.SetConfigFile(in) ***REMOVED***
func (v *Viper) SetConfigFile(in string) ***REMOVED***
	if in != "" ***REMOVED***
		v.configFile = in
	***REMOVED***
***REMOVED***

// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use.
// E.g. if your prefix is "spf", the env registry will look for env
// variables that start with "SPF_".
func SetEnvPrefix(in string) ***REMOVED*** v.SetEnvPrefix(in) ***REMOVED***
func (v *Viper) SetEnvPrefix(in string) ***REMOVED***
	if in != "" ***REMOVED***
		v.envPrefix = in
	***REMOVED***
***REMOVED***

func (v *Viper) mergeWithEnvPrefix(in string) string ***REMOVED***
	if v.envPrefix != "" ***REMOVED***
		return strings.ToUpper(v.envPrefix + "_" + in)
	***REMOVED***

	return strings.ToUpper(in)
***REMOVED***

// TODO: should getEnv logic be moved into find(). Can generalize the use of
// rewriting keys many things, Ex: Get('someKey') -> some_key
// (camel case to snake case for JSON keys perhaps)

// getEnv is a wrapper around os.Getenv which replaces characters in the original
// key. This allows env vars which have different keys than the config object
// keys.
func (v *Viper) getEnv(key string) string ***REMOVED***
	if v.envKeyReplacer != nil ***REMOVED***
		key = v.envKeyReplacer.Replace(key)
	***REMOVED***
	return os.Getenv(key)
***REMOVED***

// ConfigFileUsed returns the file used to populate the config registry.
func ConfigFileUsed() string            ***REMOVED*** return v.ConfigFileUsed() ***REMOVED***
func (v *Viper) ConfigFileUsed() string ***REMOVED*** return v.configFile ***REMOVED***

// AddConfigPath adds a path for Viper to search for the config file in.
// Can be called multiple times to define multiple search paths.
func AddConfigPath(in string) ***REMOVED*** v.AddConfigPath(in) ***REMOVED***
func (v *Viper) AddConfigPath(in string) ***REMOVED***
	if in != "" ***REMOVED***
		absin := absPathify(in)
		jww.INFO.Println("adding", absin, "to paths to search")
		if !stringInSlice(absin, v.configPaths) ***REMOVED***
			v.configPaths = append(v.configPaths, absin)
		***REMOVED***
	***REMOVED***
***REMOVED***

// AddRemoteProvider adds a remote configuration source.
// Remote Providers are searched in the order they are added.
// provider is a string value, "etcd" or "consul" are currently supported.
// endpoint is the url.  etcd requires http://ip:port  consul requires ip:port
// path is the path in the k/v store to retrieve configuration
// To retrieve a config file called myapp.json from /configs/myapp.json
// you should set path to /configs and set config name (SetConfigName()) to
// "myapp"
func AddRemoteProvider(provider, endpoint, path string) error ***REMOVED***
	return v.AddRemoteProvider(provider, endpoint, path)
***REMOVED***
func (v *Viper) AddRemoteProvider(provider, endpoint, path string) error ***REMOVED***
	if !stringInSlice(provider, SupportedRemoteProviders) ***REMOVED***
		return UnsupportedRemoteProviderError(provider)
	***REMOVED***
	if provider != "" && endpoint != "" ***REMOVED***
		jww.INFO.Printf("adding %s:%s to remote provider list", provider, endpoint)
		rp := &defaultRemoteProvider***REMOVED***
			endpoint: endpoint,
			provider: provider,
			path:     path,
		***REMOVED***
		if !v.providerPathExists(rp) ***REMOVED***
			v.remoteProviders = append(v.remoteProviders, rp)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// AddSecureRemoteProvider adds a remote configuration source.
// Secure Remote Providers are searched in the order they are added.
// provider is a string value, "etcd" or "consul" are currently supported.
// endpoint is the url.  etcd requires http://ip:port  consul requires ip:port
// secretkeyring is the filepath to your openpgp secret keyring.  e.g. /etc/secrets/myring.gpg
// path is the path in the k/v store to retrieve configuration
// To retrieve a config file called myapp.json from /configs/myapp.json
// you should set path to /configs and set config name (SetConfigName()) to
// "myapp"
// Secure Remote Providers are implemented with github.com/xordataexchange/crypt
func AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error ***REMOVED***
	return v.AddSecureRemoteProvider(provider, endpoint, path, secretkeyring)
***REMOVED***

func (v *Viper) AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error ***REMOVED***
	if !stringInSlice(provider, SupportedRemoteProviders) ***REMOVED***
		return UnsupportedRemoteProviderError(provider)
	***REMOVED***
	if provider != "" && endpoint != "" ***REMOVED***
		jww.INFO.Printf("adding %s:%s to remote provider list", provider, endpoint)
		rp := &defaultRemoteProvider***REMOVED***
			endpoint:      endpoint,
			provider:      provider,
			path:          path,
			secretKeyring: secretkeyring,
		***REMOVED***
		if !v.providerPathExists(rp) ***REMOVED***
			v.remoteProviders = append(v.remoteProviders, rp)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (v *Viper) providerPathExists(p *defaultRemoteProvider) bool ***REMOVED***
	for _, y := range v.remoteProviders ***REMOVED***
		if reflect.DeepEqual(y, p) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// searchMap recursively searches for a value for path in source map.
// Returns nil if not found.
// Note: This assumes that the path entries and map keys are lower cased.
func (v *Viper) searchMap(source map[string]interface***REMOVED******REMOVED***, path []string) interface***REMOVED******REMOVED*** ***REMOVED***
	if len(path) == 0 ***REMOVED***
		return source
	***REMOVED***

	next, ok := source[path[0]]
	if ok ***REMOVED***
		// Fast path
		if len(path) == 1 ***REMOVED***
			return next
		***REMOVED***

		// Nested case
		switch next.(type) ***REMOVED***
		case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
			return v.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface***REMOVED******REMOVED***:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return v.searchMap(next.(map[string]interface***REMOVED******REMOVED***), path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// searchMapWithPathPrefixes recursively searches for a value for path in source map.
//
// While searchMap() considers each path element as a single map key, this
// function searches for, and prioritizes, merged path elements.
// e.g., if in the source, "foo" is defined with a sub-key "bar", and "foo.bar"
// is also defined, this latter value is returned for path ["foo", "bar"].
//
// This should be useful only at config level (other maps may not contain dots
// in their keys).
//
// Note: This assumes that the path entries and map keys are lower cased.
func (v *Viper) searchMapWithPathPrefixes(source map[string]interface***REMOVED******REMOVED***, path []string) interface***REMOVED******REMOVED*** ***REMOVED***
	if len(path) == 0 ***REMOVED***
		return source
	***REMOVED***

	// search for path prefixes, starting from the longest one
	for i := len(path); i > 0; i-- ***REMOVED***
		prefixKey := strings.ToLower(strings.Join(path[0:i], v.keyDelim))

		next, ok := source[prefixKey]
		if ok ***REMOVED***
			// Fast path
			if i == len(path) ***REMOVED***
				return next
			***REMOVED***

			// Nested case
			var val interface***REMOVED******REMOVED***
			switch next.(type) ***REMOVED***
			case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
				val = v.searchMapWithPathPrefixes(cast.ToStringMap(next), path[i:])
			case map[string]interface***REMOVED******REMOVED***:
				// Type assertion is safe here since it is only reached
				// if the type of `next` is the same as the type being asserted
				val = v.searchMapWithPathPrefixes(next.(map[string]interface***REMOVED******REMOVED***), path[i:])
			default:
				// got a value but nested key expected, do nothing and look for next prefix
			***REMOVED***
			if val != nil ***REMOVED***
				return val
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// not found
	return nil
***REMOVED***

// isPathShadowedInDeepMap makes sure the given path is not shadowed somewhere
// on its path in the map.
// e.g., if "foo.bar" has a value in the given map, it “shadows”
//       "foo.bar.baz" in a lower-priority map
func (v *Viper) isPathShadowedInDeepMap(path []string, m map[string]interface***REMOVED******REMOVED***) string ***REMOVED***
	var parentVal interface***REMOVED******REMOVED***
	for i := 1; i < len(path); i++ ***REMOVED***
		parentVal = v.searchMap(m, path[0:i])
		if parentVal == nil ***REMOVED***
			// not found, no need to add more path elements
			return ""
		***REMOVED***
		switch parentVal.(type) ***REMOVED***
		case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
			continue
		case map[string]interface***REMOVED******REMOVED***:
			continue
		default:
			// parentVal is a regular value which shadows "path"
			return strings.Join(path[0:i], v.keyDelim)
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// isPathShadowedInFlatMap makes sure the given path is not shadowed somewhere
// in a sub-path of the map.
// e.g., if "foo.bar" has a value in the given map, it “shadows”
//       "foo.bar.baz" in a lower-priority map
func (v *Viper) isPathShadowedInFlatMap(path []string, mi interface***REMOVED******REMOVED***) string ***REMOVED***
	// unify input map
	var m map[string]interface***REMOVED******REMOVED***
	switch mi.(type) ***REMOVED***
	case map[string]string, map[string]FlagValue:
		m = cast.ToStringMap(mi)
	default:
		return ""
	***REMOVED***

	// scan paths
	var parentKey string
	for i := 1; i < len(path); i++ ***REMOVED***
		parentKey = strings.Join(path[0:i], v.keyDelim)
		if _, ok := m[parentKey]; ok ***REMOVED***
			return parentKey
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// isPathShadowedInAutoEnv makes sure the given path is not shadowed somewhere
// in the environment, when automatic env is on.
// e.g., if "foo.bar" has a value in the environment, it “shadows”
//       "foo.bar.baz" in a lower-priority map
func (v *Viper) isPathShadowedInAutoEnv(path []string) string ***REMOVED***
	var parentKey string
	var val string
	for i := 1; i < len(path); i++ ***REMOVED***
		parentKey = strings.Join(path[0:i], v.keyDelim)
		if val = v.getEnv(v.mergeWithEnvPrefix(parentKey)); val != "" ***REMOVED***
			return parentKey
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// SetTypeByDefaultValue enables or disables the inference of a key value's
// type when the Get function is used based upon a key's default value as
// opposed to the value returned based on the normal fetch logic.
//
// For example, if a key has a default value of []string***REMOVED******REMOVED*** and the same key
// is set via an environment variable to "a b c", a call to the Get function
// would return a string slice for the key if the key's type is inferred by
// the default value and the Get function would return:
//
//   []string ***REMOVED***"a", "b", "c"***REMOVED***
//
// Otherwise the Get function would return:
//
//   "a b c"
func SetTypeByDefaultValue(enable bool) ***REMOVED*** v.SetTypeByDefaultValue(enable) ***REMOVED***
func (v *Viper) SetTypeByDefaultValue(enable bool) ***REMOVED***
	v.typeByDefValue = enable
***REMOVED***

// GetViper gets the global Viper instance.
func GetViper() *Viper ***REMOVED***
	return v
***REMOVED***

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface***REMOVED******REMOVED*** ***REMOVED*** return v.Get(key) ***REMOVED***
func (v *Viper) Get(key string) interface***REMOVED******REMOVED*** ***REMOVED***
	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey)
	if val == nil ***REMOVED***
		return nil
	***REMOVED***

	if v.typeByDefValue ***REMOVED***
		// TODO(bep) this branch isn't covered by a single test.
		valType := val
		path := strings.Split(lcaseKey, v.keyDelim)
		defVal := v.searchMap(v.defaults, path)
		if defVal != nil ***REMOVED***
			valType = defVal
		***REMOVED***

		switch valType.(type) ***REMOVED***
		case bool:
			return cast.ToBool(val)
		case string:
			return cast.ToString(val)
		case int64, int32, int16, int8, int:
			return cast.ToInt(val)
		case float64, float32:
			return cast.ToFloat64(val)
		case time.Time:
			return cast.ToTime(val)
		case time.Duration:
			return cast.ToDuration(val)
		case []string:
			return cast.ToStringSlice(val)
		***REMOVED***
	***REMOVED***

	return val
***REMOVED***

// Sub returns new Viper instance representing a sub tree of this instance.
// Sub is case-insensitive for a key.
func Sub(key string) *Viper ***REMOVED*** return v.Sub(key) ***REMOVED***
func (v *Viper) Sub(key string) *Viper ***REMOVED***
	subv := New()
	data := v.Get(key)
	if data == nil ***REMOVED***
		return nil
	***REMOVED***

	if reflect.TypeOf(data).Kind() == reflect.Map ***REMOVED***
		subv.config = cast.ToStringMap(data)
		return subv
	***REMOVED***
	return nil
***REMOVED***

// GetString returns the value associated with the key as a string.
func GetString(key string) string ***REMOVED*** return v.GetString(key) ***REMOVED***
func (v *Viper) GetString(key string) string ***REMOVED***
	return cast.ToString(v.Get(key))
***REMOVED***

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool ***REMOVED*** return v.GetBool(key) ***REMOVED***
func (v *Viper) GetBool(key string) bool ***REMOVED***
	return cast.ToBool(v.Get(key))
***REMOVED***

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int ***REMOVED*** return v.GetInt(key) ***REMOVED***
func (v *Viper) GetInt(key string) int ***REMOVED***
	return cast.ToInt(v.Get(key))
***REMOVED***

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 ***REMOVED*** return v.GetInt64(key) ***REMOVED***
func (v *Viper) GetInt64(key string) int64 ***REMOVED***
	return cast.ToInt64(v.Get(key))
***REMOVED***

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 ***REMOVED*** return v.GetFloat64(key) ***REMOVED***
func (v *Viper) GetFloat64(key string) float64 ***REMOVED***
	return cast.ToFloat64(v.Get(key))
***REMOVED***

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time ***REMOVED*** return v.GetTime(key) ***REMOVED***
func (v *Viper) GetTime(key string) time.Time ***REMOVED***
	return cast.ToTime(v.Get(key))
***REMOVED***

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration ***REMOVED*** return v.GetDuration(key) ***REMOVED***
func (v *Viper) GetDuration(key string) time.Duration ***REMOVED***
	return cast.ToDuration(v.Get(key))
***REMOVED***

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string ***REMOVED*** return v.GetStringSlice(key) ***REMOVED***
func (v *Viper) GetStringSlice(key string) []string ***REMOVED***
	return cast.ToStringSlice(v.Get(key))
***REMOVED***

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface***REMOVED******REMOVED*** ***REMOVED*** return v.GetStringMap(key) ***REMOVED***
func (v *Viper) GetStringMap(key string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return cast.ToStringMap(v.Get(key))
***REMOVED***

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string ***REMOVED*** return v.GetStringMapString(key) ***REMOVED***
func (v *Viper) GetStringMapString(key string) map[string]string ***REMOVED***
	return cast.ToStringMapString(v.Get(key))
***REMOVED***

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string ***REMOVED*** return v.GetStringMapStringSlice(key) ***REMOVED***
func (v *Viper) GetStringMapStringSlice(key string) map[string][]string ***REMOVED***
	return cast.ToStringMapStringSlice(v.Get(key))
***REMOVED***

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func GetSizeInBytes(key string) uint ***REMOVED*** return v.GetSizeInBytes(key) ***REMOVED***
func (v *Viper) GetSizeInBytes(key string) uint ***REMOVED***
	sizeStr := cast.ToString(v.Get(key))
	return parseSizeInBytes(sizeStr)
***REMOVED***

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func UnmarshalKey(key string, rawVal interface***REMOVED******REMOVED***) error ***REMOVED*** return v.UnmarshalKey(key, rawVal) ***REMOVED***
func (v *Viper) UnmarshalKey(key string, rawVal interface***REMOVED******REMOVED***) error ***REMOVED***
	err := decode(v.Get(key), defaultDecoderConfig(rawVal))

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	v.insensitiviseMaps()

	return nil
***REMOVED***

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func Unmarshal(rawVal interface***REMOVED******REMOVED***) error ***REMOVED*** return v.Unmarshal(rawVal) ***REMOVED***
func (v *Viper) Unmarshal(rawVal interface***REMOVED******REMOVED***) error ***REMOVED***
	err := decode(v.AllSettings(), defaultDecoderConfig(rawVal))

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	v.insensitiviseMaps()

	return nil
***REMOVED***

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot
// of time.Duration values & string slices
func defaultDecoderConfig(output interface***REMOVED******REMOVED***) *mapstructure.DecoderConfig ***REMOVED***
	return &mapstructure.DecoderConfig***REMOVED***
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	***REMOVED***
***REMOVED***

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface***REMOVED******REMOVED***, config *mapstructure.DecoderConfig) error ***REMOVED***
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return decoder.Decode(input)
***REMOVED***

// UnmarshalExact unmarshals the config into a Struct, erroring if a field is nonexistent
// in the destination struct.
func (v *Viper) UnmarshalExact(rawVal interface***REMOVED******REMOVED***) error ***REMOVED***
	config := defaultDecoderConfig(rawVal)
	config.ErrorUnused = true

	err := decode(v.AllSettings(), config)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	v.insensitiviseMaps()

	return nil
***REMOVED***

// BindPFlags binds a full flag set to the configuration, using each flag's long
// name as the config key.
func BindPFlags(flags *pflag.FlagSet) error ***REMOVED*** return v.BindPFlags(flags) ***REMOVED***
func (v *Viper) BindPFlags(flags *pflag.FlagSet) error ***REMOVED***
	return v.BindFlagValues(pflagValueSet***REMOVED***flags***REMOVED***)
***REMOVED***

// BindPFlag binds a specific key to a pflag (as used by cobra).
// Example (where serverCmd is a Cobra instance):
//
//	 serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
//	 Viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
//
func BindPFlag(key string, flag *pflag.Flag) error ***REMOVED*** return v.BindPFlag(key, flag) ***REMOVED***
func (v *Viper) BindPFlag(key string, flag *pflag.Flag) error ***REMOVED***
	return v.BindFlagValue(key, pflagValue***REMOVED***flag***REMOVED***)
***REMOVED***

// BindFlagValues binds a full FlagValue set to the configuration, using each flag's long
// name as the config key.
func BindFlagValues(flags FlagValueSet) error ***REMOVED*** return v.BindFlagValues(flags) ***REMOVED***
func (v *Viper) BindFlagValues(flags FlagValueSet) (err error) ***REMOVED***
	flags.VisitAll(func(flag FlagValue) ***REMOVED***
		if err = v.BindFlagValue(flag.Name(), flag); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***)
	return nil
***REMOVED***

// BindFlagValue binds a specific key to a FlagValue.
// Example (where serverCmd is a Cobra instance):
//
//	 serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
//	 Viper.BindFlagValue("port", serverCmd.Flags().Lookup("port"))
//
func BindFlagValue(key string, flag FlagValue) error ***REMOVED*** return v.BindFlagValue(key, flag) ***REMOVED***
func (v *Viper) BindFlagValue(key string, flag FlagValue) error ***REMOVED***
	if flag == nil ***REMOVED***
		return fmt.Errorf("flag for %q is nil", key)
	***REMOVED***
	v.pflags[strings.ToLower(key)] = flag
	return nil
***REMOVED***

// BindEnv binds a Viper key to a ENV variable.
// ENV variables are case sensitive.
// If only a key is provided, it will use the env key matching the key, uppercased.
// EnvPrefix will be used when set when env name is not provided.
func BindEnv(input ...string) error ***REMOVED*** return v.BindEnv(input...) ***REMOVED***
func (v *Viper) BindEnv(input ...string) error ***REMOVED***
	var key, envkey string
	if len(input) == 0 ***REMOVED***
		return fmt.Errorf("BindEnv missing key to bind to")
	***REMOVED***

	key = strings.ToLower(input[0])

	if len(input) == 1 ***REMOVED***
		envkey = v.mergeWithEnvPrefix(key)
	***REMOVED*** else ***REMOVED***
		envkey = input[1]
	***REMOVED***

	v.env[key] = envkey

	return nil
***REMOVED***

// Given a key, find the value.
// Viper will check in the following order:
// flag, env, config file, key/value store, default.
// Viper will check to see if an alias exists first.
// Note: this assumes a lower-cased key given.
func (v *Viper) find(lcaseKey string) interface***REMOVED******REMOVED*** ***REMOVED***

	var (
		val    interface***REMOVED******REMOVED***
		exists bool
		path   = strings.Split(lcaseKey, v.keyDelim)
		nested = len(path) > 1
	)

	// compute the path through the nested maps to the nested value
	if nested && v.isPathShadowedInDeepMap(path, castMapStringToMapInterface(v.aliases)) != "" ***REMOVED***
		return nil
	***REMOVED***

	// if the requested key is an alias, then return the proper key
	lcaseKey = v.realKey(lcaseKey)
	path = strings.Split(lcaseKey, v.keyDelim)
	nested = len(path) > 1

	// Set() override first
	val = v.searchMap(v.override, path)
	if val != nil ***REMOVED***
		return val
	***REMOVED***
	if nested && v.isPathShadowedInDeepMap(path, v.override) != "" ***REMOVED***
		return nil
	***REMOVED***

	// PFlag override next
	flag, exists := v.pflags[lcaseKey]
	if exists && flag.HasChanged() ***REMOVED***
		switch flag.ValueType() ***REMOVED***
		case "int", "int8", "int16", "int32", "int64":
			return cast.ToInt(flag.ValueString())
		case "bool":
			return cast.ToBool(flag.ValueString())
		case "stringSlice":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := readAsCSV(s)
			return res
		default:
			return flag.ValueString()
		***REMOVED***
	***REMOVED***
	if nested && v.isPathShadowedInFlatMap(path, v.pflags) != "" ***REMOVED***
		return nil
	***REMOVED***

	// Env override next
	if v.automaticEnvApplied ***REMOVED***
		// even if it hasn't been registered, if automaticEnv is used,
		// check any Get request
		if val = v.getEnv(v.mergeWithEnvPrefix(lcaseKey)); val != "" ***REMOVED***
			return val
		***REMOVED***
		if nested && v.isPathShadowedInAutoEnv(path) != "" ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	envkey, exists := v.env[lcaseKey]
	if exists ***REMOVED***
		if val = v.getEnv(envkey); val != "" ***REMOVED***
			return val
		***REMOVED***
	***REMOVED***
	if nested && v.isPathShadowedInFlatMap(path, v.env) != "" ***REMOVED***
		return nil
	***REMOVED***

	// Config file next
	val = v.searchMapWithPathPrefixes(v.config, path)
	if val != nil ***REMOVED***
		return val
	***REMOVED***
	if nested && v.isPathShadowedInDeepMap(path, v.config) != "" ***REMOVED***
		return nil
	***REMOVED***

	// K/V store next
	val = v.searchMap(v.kvstore, path)
	if val != nil ***REMOVED***
		return val
	***REMOVED***
	if nested && v.isPathShadowedInDeepMap(path, v.kvstore) != "" ***REMOVED***
		return nil
	***REMOVED***

	// Default next
	val = v.searchMap(v.defaults, path)
	if val != nil ***REMOVED***
		return val
	***REMOVED***
	if nested && v.isPathShadowedInDeepMap(path, v.defaults) != "" ***REMOVED***
		return nil
	***REMOVED***

	// last chance: if no other value is returned and a flag does exist for the value,
	// get the flag's value even if the flag's value has not changed
	if flag, exists := v.pflags[lcaseKey]; exists ***REMOVED***
		switch flag.ValueType() ***REMOVED***
		case "int", "int8", "int16", "int32", "int64":
			return cast.ToInt(flag.ValueString())
		case "bool":
			return cast.ToBool(flag.ValueString())
		case "stringSlice":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := readAsCSV(s)
			return res
		default:
			return flag.ValueString()
		***REMOVED***
	***REMOVED***
	// last item, no need to check shadowing

	return nil
***REMOVED***

func readAsCSV(val string) ([]string, error) ***REMOVED***
	if val == "" ***REMOVED***
		return []string***REMOVED******REMOVED***, nil
	***REMOVED***
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
***REMOVED***

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
func IsSet(key string) bool ***REMOVED*** return v.IsSet(key) ***REMOVED***
func (v *Viper) IsSet(key string) bool ***REMOVED***
	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey)
	return val != nil
***REMOVED***

// AutomaticEnv has Viper check ENV variables for all.
// keys set in config, default & flags
func AutomaticEnv() ***REMOVED*** v.AutomaticEnv() ***REMOVED***
func (v *Viper) AutomaticEnv() ***REMOVED***
	v.automaticEnvApplied = true
***REMOVED***

// SetEnvKeyReplacer sets the strings.Replacer on the viper object
// Useful for mapping an environmental variable to a key that does
// not match it.
func SetEnvKeyReplacer(r *strings.Replacer) ***REMOVED*** v.SetEnvKeyReplacer(r) ***REMOVED***
func (v *Viper) SetEnvKeyReplacer(r *strings.Replacer) ***REMOVED***
	v.envKeyReplacer = r
***REMOVED***

// Aliases provide another accessor for the same key.
// This enables one to change a name without breaking the application
func RegisterAlias(alias string, key string) ***REMOVED*** v.RegisterAlias(alias, key) ***REMOVED***
func (v *Viper) RegisterAlias(alias string, key string) ***REMOVED***
	v.registerAlias(alias, strings.ToLower(key))
***REMOVED***

func (v *Viper) registerAlias(alias string, key string) ***REMOVED***
	alias = strings.ToLower(alias)
	if alias != key && alias != v.realKey(key) ***REMOVED***
		_, exists := v.aliases[alias]

		if !exists ***REMOVED***
			// if we alias something that exists in one of the maps to another
			// name, we'll never be able to get that value using the original
			// name, so move the config value to the new realkey.
			if val, ok := v.config[alias]; ok ***REMOVED***
				delete(v.config, alias)
				v.config[key] = val
			***REMOVED***
			if val, ok := v.kvstore[alias]; ok ***REMOVED***
				delete(v.kvstore, alias)
				v.kvstore[key] = val
			***REMOVED***
			if val, ok := v.defaults[alias]; ok ***REMOVED***
				delete(v.defaults, alias)
				v.defaults[key] = val
			***REMOVED***
			if val, ok := v.override[alias]; ok ***REMOVED***
				delete(v.override, alias)
				v.override[key] = val
			***REMOVED***
			v.aliases[alias] = key
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		jww.WARN.Println("Creating circular reference alias", alias, key, v.realKey(key))
	***REMOVED***
***REMOVED***

func (v *Viper) realKey(key string) string ***REMOVED***
	newkey, exists := v.aliases[key]
	if exists ***REMOVED***
		jww.DEBUG.Println("Alias", key, "to", newkey)
		return v.realKey(newkey)
	***REMOVED***
	return key
***REMOVED***

// InConfig checks to see if the given key (or an alias) is in the config file.
func InConfig(key string) bool ***REMOVED*** return v.InConfig(key) ***REMOVED***
func (v *Viper) InConfig(key string) bool ***REMOVED***
	// if the requested key is an alias, then return the proper key
	key = v.realKey(key)

	_, exists := v.config[key]
	return exists
***REMOVED***

// SetDefault sets the default value for this key.
// SetDefault is case-insensitive for a key.
// Default only used when no value is provided by the user via flag, config or ENV.
func SetDefault(key string, value interface***REMOVED******REMOVED***) ***REMOVED*** v.SetDefault(key, value) ***REMOVED***
func (v *Viper) SetDefault(key string, value interface***REMOVED******REMOVED***) ***REMOVED***
	// If alias passed in, then set the proper default
	key = v.realKey(strings.ToLower(key))
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(v.defaults, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value
***REMOVED***

// Set sets the value for the key in the override regiser.
// Set is case-insensitive for a key.
// Will be used instead of values obtained via
// flags, config file, ENV, default, or key/value store.
func Set(key string, value interface***REMOVED******REMOVED***) ***REMOVED*** v.Set(key, value) ***REMOVED***
func (v *Viper) Set(key string, value interface***REMOVED******REMOVED***) ***REMOVED***
	// If alias passed in, then set the proper override
	key = v.realKey(strings.ToLower(key))
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(v.override, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value
***REMOVED***

// ReadInConfig will discover and load the configuration file from disk
// and key/value stores, searching in one of the defined paths.
func ReadInConfig() error ***REMOVED*** return v.ReadInConfig() ***REMOVED***
func (v *Viper) ReadInConfig() error ***REMOVED***
	jww.INFO.Println("Attempting to read in config file")
	filename, err := v.getConfigFile()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !stringInSlice(v.getConfigType(), SupportedExts) ***REMOVED***
		return UnsupportedConfigError(v.getConfigType())
	***REMOVED***

	jww.DEBUG.Println("Reading file: ", filename)
	file, err := afero.ReadFile(v.fs, filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	config := make(map[string]interface***REMOVED******REMOVED***)

	err = v.unmarshalReader(bytes.NewReader(file), config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	v.config = config
	return nil
***REMOVED***

// MergeInConfig merges a new configuration with an existing config.
func MergeInConfig() error ***REMOVED*** return v.MergeInConfig() ***REMOVED***
func (v *Viper) MergeInConfig() error ***REMOVED***
	jww.INFO.Println("Attempting to merge in config file")
	filename, err := v.getConfigFile()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !stringInSlice(v.getConfigType(), SupportedExts) ***REMOVED***
		return UnsupportedConfigError(v.getConfigType())
	***REMOVED***

	file, err := afero.ReadFile(v.fs, filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return v.MergeConfig(bytes.NewReader(file))
***REMOVED***

// ReadConfig will read a configuration file, setting existing keys to nil if the
// key does not exist in the file.
func ReadConfig(in io.Reader) error ***REMOVED*** return v.ReadConfig(in) ***REMOVED***
func (v *Viper) ReadConfig(in io.Reader) error ***REMOVED***
	v.config = make(map[string]interface***REMOVED******REMOVED***)
	return v.unmarshalReader(in, v.config)
***REMOVED***

// MergeConfig merges a new configuration with an existing config.
func MergeConfig(in io.Reader) error ***REMOVED*** return v.MergeConfig(in) ***REMOVED***
func (v *Viper) MergeConfig(in io.Reader) error ***REMOVED***
	if v.config == nil ***REMOVED***
		v.config = make(map[string]interface***REMOVED******REMOVED***)
	***REMOVED***
	cfg := make(map[string]interface***REMOVED******REMOVED***)
	if err := v.unmarshalReader(in, cfg); err != nil ***REMOVED***
		return err
	***REMOVED***
	mergeMaps(cfg, v.config, nil)
	return nil
***REMOVED***

// WriteConfig writes the current configuration to a file.
func WriteConfig() error ***REMOVED*** return v.WriteConfig() ***REMOVED***
func (v *Viper) WriteConfig() error ***REMOVED***
	filename, err := v.getConfigFile()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return v.writeConfig(filename, true)
***REMOVED***

// SafeWriteConfig writes current configuration to file only if the file does not exist.
func SafeWriteConfig() error ***REMOVED*** return v.SafeWriteConfig() ***REMOVED***
func (v *Viper) SafeWriteConfig() error ***REMOVED***
	filename, err := v.getConfigFile()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return v.writeConfig(filename, false)
***REMOVED***

// WriteConfigAs writes current configuration to a given filename.
func WriteConfigAs(filename string) error ***REMOVED*** return v.WriteConfigAs(filename) ***REMOVED***
func (v *Viper) WriteConfigAs(filename string) error ***REMOVED***
	return v.writeConfig(filename, true)
***REMOVED***

// SafeWriteConfigAs writes current configuration to a given filename if it does not exist.
func SafeWriteConfigAs(filename string) error ***REMOVED*** return v.SafeWriteConfigAs(filename) ***REMOVED***
func (v *Viper) SafeWriteConfigAs(filename string) error ***REMOVED***
	return v.writeConfig(filename, false)
***REMOVED***

func writeConfig(filename string, force bool) error ***REMOVED*** return v.writeConfig(filename, force) ***REMOVED***
func (v *Viper) writeConfig(filename string, force bool) error ***REMOVED***
	jww.INFO.Println("Attempting to write configuration to file.")
	ext := filepath.Ext(filename)
	if len(ext) <= 1 ***REMOVED***
		return fmt.Errorf("Filename: %s requires valid extension.", filename)
	***REMOVED***
	configType := ext[1:]
	if !stringInSlice(configType, SupportedExts) ***REMOVED***
		return UnsupportedConfigError(configType)
	***REMOVED***
	if v.config == nil ***REMOVED***
		v.config = make(map[string]interface***REMOVED******REMOVED***)
	***REMOVED***
	var flags int
	if force == true ***REMOVED***
		flags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	***REMOVED*** else ***REMOVED***
		if _, err := os.Stat(filename); os.IsNotExist(err) ***REMOVED***
			flags = os.O_WRONLY
		***REMOVED*** else ***REMOVED***
			return fmt.Errorf("File: %s exists. Use WriteConfig to overwrite.", filename)
		***REMOVED***
	***REMOVED***
	f, err := v.fs.OpenFile(filename, flags, os.FileMode(0644))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return v.marshalWriter(f, configType)
***REMOVED***

// Unmarshal a Reader into a map.
// Should probably be an unexported function.
func unmarshalReader(in io.Reader, c map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return v.unmarshalReader(in, c)
***REMOVED***
func (v *Viper) unmarshalReader(in io.Reader, c map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	switch strings.ToLower(v.getConfigType()) ***REMOVED***
	case "yaml", "yml":
		if err := yaml.Unmarshal(buf.Bytes(), &c); err != nil ***REMOVED***
			return ConfigParseError***REMOVED***err***REMOVED***
		***REMOVED***

	case "json":
		if err := json.Unmarshal(buf.Bytes(), &c); err != nil ***REMOVED***
			return ConfigParseError***REMOVED***err***REMOVED***
		***REMOVED***

	case "hcl":
		obj, err := hcl.Parse(string(buf.Bytes()))
		if err != nil ***REMOVED***
			return ConfigParseError***REMOVED***err***REMOVED***
		***REMOVED***
		if err = hcl.DecodeObject(&c, obj); err != nil ***REMOVED***
			return ConfigParseError***REMOVED***err***REMOVED***
		***REMOVED***

	case "toml":
		tree, err := toml.LoadReader(buf)
		if err != nil ***REMOVED***
			return ConfigParseError***REMOVED***err***REMOVED***
		***REMOVED***
		tmap := tree.ToMap()
		for k, v := range tmap ***REMOVED***
			c[k] = v
		***REMOVED***

	case "properties", "props", "prop":
		v.properties = properties.NewProperties()
		var err error
		if v.properties, err = properties.Load(buf.Bytes(), properties.UTF8); err != nil ***REMOVED***
			return ConfigParseError***REMOVED***err***REMOVED***
		***REMOVED***
		for _, key := range v.properties.Keys() ***REMOVED***
			value, _ := v.properties.Get(key)
			// recursively build nested maps
			path := strings.Split(key, ".")
			lastKey := strings.ToLower(path[len(path)-1])
			deepestMap := deepSearch(c, path[0:len(path)-1])
			// set innermost value
			deepestMap[lastKey] = value
		***REMOVED***
	***REMOVED***

	insensitiviseMap(c)
	return nil
***REMOVED***

// Marshal a map into Writer.
func marshalWriter(f afero.File, configType string) error ***REMOVED***
	return v.marshalWriter(f, configType)
***REMOVED***
func (v *Viper) marshalWriter(f afero.File, configType string) error ***REMOVED***
	c := v.AllSettings()
	switch configType ***REMOVED***
	case "json":
		b, err := json.MarshalIndent(c, "", "  ")
		if err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***
		_, err = f.WriteString(string(b))
		if err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***

	case "hcl":
		b, err := json.Marshal(c)
		ast, err := hcl.Parse(string(b))
		if err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***
		err = printer.Fprint(f, ast.Node)
		if err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***

	case "prop", "props", "properties":
		if v.properties == nil ***REMOVED***
			v.properties = properties.NewProperties()
		***REMOVED***
		p := v.properties
		for _, key := range v.AllKeys() ***REMOVED***
			_, _, err := p.Set(key, v.GetString(key))
			if err != nil ***REMOVED***
				return ConfigMarshalError***REMOVED***err***REMOVED***
			***REMOVED***
		***REMOVED***
		_, err := p.WriteComment(f, "#", properties.UTF8)
		if err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***

	case "toml":
		t, err := toml.TreeFromMap(c)
		if err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***
		s := t.String()
		if _, err := f.WriteString(s); err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***

	case "yaml", "yml":
		b, err := yaml.Marshal(c)
		if err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***
		if _, err = f.WriteString(string(b)); err != nil ***REMOVED***
			return ConfigMarshalError***REMOVED***err***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func keyExists(k string, m map[string]interface***REMOVED******REMOVED***) string ***REMOVED***
	lk := strings.ToLower(k)
	for mk := range m ***REMOVED***
		lmk := strings.ToLower(mk)
		if lmk == lk ***REMOVED***
			return mk
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func castToMapStringInterface(
	src map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	tgt := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	for k, v := range src ***REMOVED***
		tgt[fmt.Sprintf("%v", k)] = v
	***REMOVED***
	return tgt
***REMOVED***

func castMapStringToMapInterface(src map[string]string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	tgt := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	for k, v := range src ***REMOVED***
		tgt[k] = v
	***REMOVED***
	return tgt
***REMOVED***

func castMapFlagToMapInterface(src map[string]FlagValue) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	tgt := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	for k, v := range src ***REMOVED***
		tgt[k] = v
	***REMOVED***
	return tgt
***REMOVED***

// mergeMaps merges two maps. The `itgt` parameter is for handling go-yaml's
// insistence on parsing nested structures as `map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***`
// instead of using a `string` as the key for nest structures beyond one level
// deep. Both map types are supported as there is a go-yaml fork that uses
// `map[string]interface***REMOVED******REMOVED***` instead.
func mergeMaps(
	src, tgt map[string]interface***REMOVED******REMOVED***, itgt map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***) ***REMOVED***
	for sk, sv := range src ***REMOVED***
		tk := keyExists(sk, tgt)
		if tk == "" ***REMOVED***
			jww.TRACE.Printf("tk=\"\", tgt[%s]=%v", sk, sv)
			tgt[sk] = sv
			if itgt != nil ***REMOVED***
				itgt[sk] = sv
			***REMOVED***
			continue
		***REMOVED***

		tv, ok := tgt[tk]
		if !ok ***REMOVED***
			jww.TRACE.Printf("tgt[%s] != ok, tgt[%s]=%v", tk, sk, sv)
			tgt[sk] = sv
			if itgt != nil ***REMOVED***
				itgt[sk] = sv
			***REMOVED***
			continue
		***REMOVED***

		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)
		if svType != tvType ***REMOVED***
			jww.ERROR.Printf(
				"svType != tvType; key=%s, st=%v, tt=%v, sv=%v, tv=%v",
				sk, svType, tvType, sv, tv)
			continue
		***REMOVED***

		jww.TRACE.Printf("processing key=%s, st=%v, tt=%v, sv=%v, tv=%v",
			sk, svType, tvType, sv, tv)

		switch ttv := tv.(type) ***REMOVED***
		case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
			jww.TRACE.Printf("merging maps (must convert)")
			tsv := sv.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
			ssv := castToMapStringInterface(tsv)
			stv := castToMapStringInterface(ttv)
			mergeMaps(ssv, stv, ttv)
		case map[string]interface***REMOVED******REMOVED***:
			jww.TRACE.Printf("merging maps")
			mergeMaps(sv.(map[string]interface***REMOVED******REMOVED***), ttv, nil)
		default:
			jww.TRACE.Printf("setting value")
			tgt[tk] = sv
			if itgt != nil ***REMOVED***
				itgt[tk] = sv
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// ReadRemoteConfig attempts to get configuration from a remote source
// and read it in the remote configuration registry.
func ReadRemoteConfig() error ***REMOVED*** return v.ReadRemoteConfig() ***REMOVED***
func (v *Viper) ReadRemoteConfig() error ***REMOVED***
	return v.getKeyValueConfig()
***REMOVED***

func WatchRemoteConfig() error ***REMOVED*** return v.WatchRemoteConfig() ***REMOVED***
func (v *Viper) WatchRemoteConfig() error ***REMOVED***
	return v.watchKeyValueConfig()
***REMOVED***

func (v *Viper) WatchRemoteConfigOnChannel() error ***REMOVED***
	return v.watchKeyValueConfigOnChannel()
***REMOVED***

func (v *Viper) insensitiviseMaps() ***REMOVED***
	insensitiviseMap(v.config)
	insensitiviseMap(v.defaults)
	insensitiviseMap(v.override)
	insensitiviseMap(v.kvstore)
***REMOVED***

// Retrieve the first found remote configuration.
func (v *Viper) getKeyValueConfig() error ***REMOVED***
	if RemoteConfig == nil ***REMOVED***
		return RemoteConfigError("Enable the remote features by doing a blank import of the viper/remote package: '_ github.com/spf13/viper/remote'")
	***REMOVED***

	for _, rp := range v.remoteProviders ***REMOVED***
		val, err := v.getRemoteConfig(rp)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		v.kvstore = val
		return nil
	***REMOVED***
	return RemoteConfigError("No Files Found")
***REMOVED***

func (v *Viper) getRemoteConfig(provider RemoteProvider) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	reader, err := RemoteConfig.Get(provider)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
***REMOVED***

// Retrieve the first found remote configuration.
func (v *Viper) watchKeyValueConfigOnChannel() error ***REMOVED***
	for _, rp := range v.remoteProviders ***REMOVED***
		respc, _ := RemoteConfig.WatchChannel(rp)
		//Todo: Add quit channel
		go func(rc <-chan *RemoteResponse) ***REMOVED***
			for ***REMOVED***
				b := <-rc
				reader := bytes.NewReader(b.Value)
				v.unmarshalReader(reader, v.kvstore)
			***REMOVED***
		***REMOVED***(respc)
		return nil
	***REMOVED***
	return RemoteConfigError("No Files Found")
***REMOVED***

// Retrieve the first found remote configuration.
func (v *Viper) watchKeyValueConfig() error ***REMOVED***
	for _, rp := range v.remoteProviders ***REMOVED***
		val, err := v.watchRemoteConfig(rp)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		v.kvstore = val
		return nil
	***REMOVED***
	return RemoteConfigError("No Files Found")
***REMOVED***

func (v *Viper) watchRemoteConfig(provider RemoteProvider) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	reader, err := RemoteConfig.Watch(provider)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
***REMOVED***

// AllKeys returns all keys holding a value, regardless of where they are set.
// Nested keys are returned with a v.keyDelim (= ".") separator
func AllKeys() []string ***REMOVED*** return v.AllKeys() ***REMOVED***
func (v *Viper) AllKeys() []string ***REMOVED***
	m := map[string]bool***REMOVED******REMOVED***
	// add all paths, by order of descending priority to ensure correct shadowing
	m = v.flattenAndMergeMap(m, castMapStringToMapInterface(v.aliases), "")
	m = v.flattenAndMergeMap(m, v.override, "")
	m = v.mergeFlatMap(m, castMapFlagToMapInterface(v.pflags))
	m = v.mergeFlatMap(m, castMapStringToMapInterface(v.env))
	m = v.flattenAndMergeMap(m, v.config, "")
	m = v.flattenAndMergeMap(m, v.kvstore, "")
	m = v.flattenAndMergeMap(m, v.defaults, "")

	// convert set of paths to list
	a := []string***REMOVED******REMOVED***
	for x := range m ***REMOVED***
		a = append(a, x)
	***REMOVED***
	return a
***REMOVED***

// flattenAndMergeMap recursively flattens the given map into a map[string]bool
// of key paths (used as a set, easier to manipulate than a []string):
// - each path is merged into a single key string, delimited with v.keyDelim (= ".")
// - if a path is shadowed by an earlier value in the initial shadow map,
//   it is skipped.
// The resulting set of paths is merged to the given shadow set at the same time.
func (v *Viper) flattenAndMergeMap(shadow map[string]bool, m map[string]interface***REMOVED******REMOVED***, prefix string) map[string]bool ***REMOVED***
	if shadow != nil && prefix != "" && shadow[prefix] ***REMOVED***
		// prefix is shadowed => nothing more to flatten
		return shadow
	***REMOVED***
	if shadow == nil ***REMOVED***
		shadow = make(map[string]bool)
	***REMOVED***

	var m2 map[string]interface***REMOVED******REMOVED***
	if prefix != "" ***REMOVED***
		prefix += v.keyDelim
	***REMOVED***
	for k, val := range m ***REMOVED***
		fullKey := prefix + k
		switch val.(type) ***REMOVED***
		case map[string]interface***REMOVED******REMOVED***:
			m2 = val.(map[string]interface***REMOVED******REMOVED***)
		case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
			m2 = cast.ToStringMap(val)
		default:
			// immediate value
			shadow[strings.ToLower(fullKey)] = true
			continue
		***REMOVED***
		// recursively merge to shadow map
		shadow = v.flattenAndMergeMap(shadow, m2, fullKey)
	***REMOVED***
	return shadow
***REMOVED***

// mergeFlatMap merges the given maps, excluding values of the second map
// shadowed by values from the first map.
func (v *Viper) mergeFlatMap(shadow map[string]bool, m map[string]interface***REMOVED******REMOVED***) map[string]bool ***REMOVED***
	// scan keys
outer:
	for k, _ := range m ***REMOVED***
		path := strings.Split(k, v.keyDelim)
		// scan intermediate paths
		var parentKey string
		for i := 1; i < len(path); i++ ***REMOVED***
			parentKey = strings.Join(path[0:i], v.keyDelim)
			if shadow[parentKey] ***REMOVED***
				// path is shadowed, continue
				continue outer
			***REMOVED***
		***REMOVED***
		// add key
		shadow[strings.ToLower(k)] = true
	***REMOVED***
	return shadow
***REMOVED***

// AllSettings merges all settings and returns them as a map[string]interface***REMOVED******REMOVED***.
func AllSettings() map[string]interface***REMOVED******REMOVED*** ***REMOVED*** return v.AllSettings() ***REMOVED***
func (v *Viper) AllSettings() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	// start from the list of keys, and construct the map one value at a time
	for _, k := range v.AllKeys() ***REMOVED***
		value := v.Get(k)
		if value == nil ***REMOVED***
			// should not happen, since AllKeys() returns only keys holding a value,
			// check just in case anything changes
			continue
		***REMOVED***
		path := strings.Split(k, v.keyDelim)
		lastKey := strings.ToLower(path[len(path)-1])
		deepestMap := deepSearch(m, path[0:len(path)-1])
		// set innermost value
		deepestMap[lastKey] = value
	***REMOVED***
	return m
***REMOVED***

// SetFs sets the filesystem to use to read configuration.
func SetFs(fs afero.Fs) ***REMOVED*** v.SetFs(fs) ***REMOVED***
func (v *Viper) SetFs(fs afero.Fs) ***REMOVED***
	v.fs = fs
***REMOVED***

// SetConfigName sets name for the config file.
// Does not include extension.
func SetConfigName(in string) ***REMOVED*** v.SetConfigName(in) ***REMOVED***
func (v *Viper) SetConfigName(in string) ***REMOVED***
	if in != "" ***REMOVED***
		v.configName = in
		v.configFile = ""
	***REMOVED***
***REMOVED***

// SetConfigType sets the type of the configuration returned by the
// remote source, e.g. "json".
func SetConfigType(in string) ***REMOVED*** v.SetConfigType(in) ***REMOVED***
func (v *Viper) SetConfigType(in string) ***REMOVED***
	if in != "" ***REMOVED***
		v.configType = in
	***REMOVED***
***REMOVED***

func (v *Viper) getConfigType() string ***REMOVED***
	if v.configType != "" ***REMOVED***
		return v.configType
	***REMOVED***

	cf, err := v.getConfigFile()
	if err != nil ***REMOVED***
		return ""
	***REMOVED***

	ext := filepath.Ext(cf)

	if len(ext) > 1 ***REMOVED***
		return ext[1:]
	***REMOVED***

	return ""
***REMOVED***

func (v *Viper) getConfigFile() (string, error) ***REMOVED***
	// if explicitly set, then use it
	if v.configFile != "" ***REMOVED***
		return v.configFile, nil
	***REMOVED***

	cf, err := v.findConfigFile()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	v.configFile = cf
	return v.getConfigFile()
***REMOVED***

func (v *Viper) searchInPath(in string) (filename string) ***REMOVED***
	jww.DEBUG.Println("Searching for config in ", in)
	for _, ext := range SupportedExts ***REMOVED***
		jww.DEBUG.Println("Checking for", filepath.Join(in, v.configName+"."+ext))
		if b, _ := exists(v.fs, filepath.Join(in, v.configName+"."+ext)); b ***REMOVED***
			jww.DEBUG.Println("Found: ", filepath.Join(in, v.configName+"."+ext))
			return filepath.Join(in, v.configName+"."+ext)
		***REMOVED***
	***REMOVED***

	return ""
***REMOVED***

// Search all configPaths for any config file.
// Returns the first path that exists (and is a config file).
func (v *Viper) findConfigFile() (string, error) ***REMOVED***
	jww.INFO.Println("Searching for config in ", v.configPaths)

	for _, cp := range v.configPaths ***REMOVED***
		file := v.searchInPath(cp)
		if file != "" ***REMOVED***
			return file, nil
		***REMOVED***
	***REMOVED***
	return "", ConfigFileNotFoundError***REMOVED***v.configName, fmt.Sprintf("%s", v.configPaths)***REMOVED***
***REMOVED***

// Debug prints all configuration registries for debugging
// purposes.
func Debug() ***REMOVED*** v.Debug() ***REMOVED***
func (v *Viper) Debug() ***REMOVED***
	fmt.Printf("Aliases:\n%#v\n", v.aliases)
	fmt.Printf("Override:\n%#v\n", v.override)
	fmt.Printf("PFlags:\n%#v\n", v.pflags)
	fmt.Printf("Env:\n%#v\n", v.env)
	fmt.Printf("Key/Value Store:\n%#v\n", v.kvstore)
	fmt.Printf("Config:\n%#v\n", v.config)
	fmt.Printf("Defaults:\n%#v\n", v.defaults)
***REMOVED***
