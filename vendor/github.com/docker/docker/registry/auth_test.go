package registry

import (
	"testing"

	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
)

func buildAuthConfigs() map[string]types.AuthConfig ***REMOVED***
	authConfigs := map[string]types.AuthConfig***REMOVED******REMOVED***

	for _, registry := range []string***REMOVED***"testIndex", IndexServer***REMOVED*** ***REMOVED***
		authConfigs[registry] = types.AuthConfig***REMOVED***
			Username: "docker-user",
			Password: "docker-pass",
		***REMOVED***
	***REMOVED***

	return authConfigs
***REMOVED***

func TestSameAuthDataPostSave(t *testing.T) ***REMOVED***
	authConfigs := buildAuthConfigs()
	authConfig := authConfigs["testIndex"]
	if authConfig.Username != "docker-user" ***REMOVED***
		t.Fail()
	***REMOVED***
	if authConfig.Password != "docker-pass" ***REMOVED***
		t.Fail()
	***REMOVED***
	if authConfig.Auth != "" ***REMOVED***
		t.Fail()
	***REMOVED***
***REMOVED***

func TestResolveAuthConfigIndexServer(t *testing.T) ***REMOVED***
	authConfigs := buildAuthConfigs()
	indexConfig := authConfigs[IndexServer]

	officialIndex := &registrytypes.IndexInfo***REMOVED***
		Official: true,
	***REMOVED***
	privateIndex := &registrytypes.IndexInfo***REMOVED***
		Official: false,
	***REMOVED***

	resolved := ResolveAuthConfig(authConfigs, officialIndex)
	assertEqual(t, resolved, indexConfig, "Expected ResolveAuthConfig to return IndexServer")

	resolved = ResolveAuthConfig(authConfigs, privateIndex)
	assertNotEqual(t, resolved, indexConfig, "Expected ResolveAuthConfig to not return IndexServer")
***REMOVED***

func TestResolveAuthConfigFullURL(t *testing.T) ***REMOVED***
	authConfigs := buildAuthConfigs()

	registryAuth := types.AuthConfig***REMOVED***
		Username: "foo-user",
		Password: "foo-pass",
	***REMOVED***
	localAuth := types.AuthConfig***REMOVED***
		Username: "bar-user",
		Password: "bar-pass",
	***REMOVED***
	officialAuth := types.AuthConfig***REMOVED***
		Username: "baz-user",
		Password: "baz-pass",
	***REMOVED***
	authConfigs[IndexServer] = officialAuth

	expectedAuths := map[string]types.AuthConfig***REMOVED***
		"registry.example.com": registryAuth,
		"localhost:8000":       localAuth,
		"registry.com":         localAuth,
	***REMOVED***

	validRegistries := map[string][]string***REMOVED***
		"registry.example.com": ***REMOVED***
			"https://registry.example.com/v1/",
			"http://registry.example.com/v1/",
			"registry.example.com",
			"registry.example.com/v1/",
		***REMOVED***,
		"localhost:8000": ***REMOVED***
			"https://localhost:8000/v1/",
			"http://localhost:8000/v1/",
			"localhost:8000",
			"localhost:8000/v1/",
		***REMOVED***,
		"registry.com": ***REMOVED***
			"https://registry.com/v1/",
			"http://registry.com/v1/",
			"registry.com",
			"registry.com/v1/",
		***REMOVED***,
	***REMOVED***

	for configKey, registries := range validRegistries ***REMOVED***
		configured, ok := expectedAuths[configKey]
		if !ok ***REMOVED***
			t.Fail()
		***REMOVED***
		index := &registrytypes.IndexInfo***REMOVED***
			Name: configKey,
		***REMOVED***
		for _, registry := range registries ***REMOVED***
			authConfigs[registry] = configured
			resolved := ResolveAuthConfig(authConfigs, index)
			if resolved.Username != configured.Username || resolved.Password != configured.Password ***REMOVED***
				t.Errorf("%s -> %v != %v\n", registry, resolved, configured)
			***REMOVED***
			delete(authConfigs, registry)
			resolved = ResolveAuthConfig(authConfigs, index)
			if resolved.Username == configured.Username || resolved.Password == configured.Password ***REMOVED***
				t.Errorf("%s -> %v == %v\n", registry, resolved, configured)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
