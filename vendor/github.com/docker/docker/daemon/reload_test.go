package daemon

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/pkg/discovery"
	_ "github.com/docker/docker/pkg/discovery/memory"
	"github.com/docker/docker/registry"
	"github.com/docker/libnetwork"
	"github.com/stretchr/testify/assert"
)

func TestDaemonReloadLabels(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	daemon.configStore = &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			Labels: []string***REMOVED***"foo:bar"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	valuesSets := make(map[string]interface***REMOVED******REMOVED***)
	valuesSets["labels"] = "foo:baz"
	newConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			Labels:    []string***REMOVED***"foo:baz"***REMOVED***,
			ValuesSet: valuesSets,
		***REMOVED***,
	***REMOVED***

	if err := daemon.Reload(newConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	label := daemon.configStore.Labels[0]
	if label != "foo:baz" ***REMOVED***
		t.Fatalf("Expected daemon label `foo:baz`, got %s", label)
	***REMOVED***
***REMOVED***

func TestDaemonReloadAllowNondistributableArtifacts(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED***
		configStore: &config.Config***REMOVED******REMOVED***,
	***REMOVED***

	var err error
	// Initialize daemon with some registries.
	daemon.RegistryService, err = registry.NewService(registry.ServiceOptions***REMOVED***
		AllowNondistributableArtifacts: []string***REMOVED***
			"127.0.0.0/8",
			"10.10.1.11:5000",
			"10.10.1.22:5000", // This will be removed during reload.
			"docker1.com",
			"docker2.com", // This will be removed during reload.
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	registries := []string***REMOVED***
		"127.0.0.0/8",
		"10.10.1.11:5000",
		"10.10.1.33:5000", // This will be added during reload.
		"docker1.com",
		"docker3.com", // This will be added during reload.
	***REMOVED***

	newConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ServiceOptions: registry.ServiceOptions***REMOVED***
				AllowNondistributableArtifacts: registries,
			***REMOVED***,
			ValuesSet: map[string]interface***REMOVED******REMOVED******REMOVED***
				"allow-nondistributable-artifacts": registries,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	if err := daemon.Reload(newConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	actual := []string***REMOVED******REMOVED***
	serviceConfig := daemon.RegistryService.ServiceConfig()
	for _, value := range serviceConfig.AllowNondistributableArtifactsCIDRs ***REMOVED***
		actual = append(actual, value.String())
	***REMOVED***
	actual = append(actual, serviceConfig.AllowNondistributableArtifactsHostnames...)

	sort.Strings(registries)
	sort.Strings(actual)
	assert.Equal(t, registries, actual)
***REMOVED***

func TestDaemonReloadMirrors(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	var err error
	daemon.RegistryService, err = registry.NewService(registry.ServiceOptions***REMOVED***
		InsecureRegistries: []string***REMOVED******REMOVED***,
		Mirrors: []string***REMOVED***
			"https://mirror.test1.com",
			"https://mirror.test2.com", // this will be removed when reloading
			"https://mirror.test3.com", // this will be removed when reloading
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	daemon.configStore = &config.Config***REMOVED******REMOVED***

	type pair struct ***REMOVED***
		valid   bool
		mirrors []string
		after   []string
	***REMOVED***

	loadMirrors := []pair***REMOVED***
		***REMOVED***
			valid:   false,
			mirrors: []string***REMOVED***"10.10.1.11:5000"***REMOVED***, // this mirror is invalid
			after:   []string***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			valid:   false,
			mirrors: []string***REMOVED***"mirror.test1.com"***REMOVED***, // this mirror is invalid
			after:   []string***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			valid:   false,
			mirrors: []string***REMOVED***"10.10.1.11:5000", "mirror.test1.com"***REMOVED***, // mirrors are invalid
			after:   []string***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			valid:   true,
			mirrors: []string***REMOVED***"https://mirror.test1.com", "https://mirror.test4.com"***REMOVED***,
			after:   []string***REMOVED***"https://mirror.test1.com/", "https://mirror.test4.com/"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, value := range loadMirrors ***REMOVED***
		valuesSets := make(map[string]interface***REMOVED******REMOVED***)
		valuesSets["registry-mirrors"] = value.mirrors

		newConfig := &config.Config***REMOVED***
			CommonConfig: config.CommonConfig***REMOVED***
				ServiceOptions: registry.ServiceOptions***REMOVED***
					Mirrors: value.mirrors,
				***REMOVED***,
				ValuesSet: valuesSets,
			***REMOVED***,
		***REMOVED***

		err := daemon.Reload(newConfig)
		if !value.valid && err == nil ***REMOVED***
			// mirrors should be invalid, should be a non-nil error
			t.Fatalf("Expected daemon reload error with invalid mirrors: %s, while get nil", value.mirrors)
		***REMOVED***

		if value.valid ***REMOVED***
			if err != nil ***REMOVED***
				// mirrors should be valid, should be no error
				t.Fatal(err)
			***REMOVED***
			registryService := daemon.RegistryService.ServiceConfig()

			if len(registryService.Mirrors) != len(value.after) ***REMOVED***
				t.Fatalf("Expected %d daemon mirrors %s while get %d with %s",
					len(value.after),
					value.after,
					len(registryService.Mirrors),
					registryService.Mirrors)
			***REMOVED***

			dataMap := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***

			for _, mirror := range registryService.Mirrors ***REMOVED***
				if _, exist := dataMap[mirror]; !exist ***REMOVED***
					dataMap[mirror] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***

			for _, address := range value.after ***REMOVED***
				if _, exist := dataMap[address]; !exist ***REMOVED***
					t.Fatalf("Expected %s in daemon mirrors, while get none", address)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDaemonReloadInsecureRegistries(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	var err error
	// initialize daemon with existing insecure registries: "127.0.0.0/8", "10.10.1.11:5000", "10.10.1.22:5000"
	daemon.RegistryService, err = registry.NewService(registry.ServiceOptions***REMOVED***
		InsecureRegistries: []string***REMOVED***
			"127.0.0.0/8",
			"10.10.1.11:5000",
			"10.10.1.22:5000", // this will be removed when reloading
			"docker1.com",
			"docker2.com", // this will be removed when reloading
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	daemon.configStore = &config.Config***REMOVED******REMOVED***

	insecureRegistries := []string***REMOVED***
		"127.0.0.0/8",     // this will be kept
		"10.10.1.11:5000", // this will be kept
		"10.10.1.33:5000", // this will be newly added
		"docker1.com",     // this will be kept
		"docker3.com",     // this will be newly added
	***REMOVED***

	valuesSets := make(map[string]interface***REMOVED******REMOVED***)
	valuesSets["insecure-registries"] = insecureRegistries

	newConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ServiceOptions: registry.ServiceOptions***REMOVED***
				InsecureRegistries: insecureRegistries,
			***REMOVED***,
			ValuesSet: valuesSets,
		***REMOVED***,
	***REMOVED***

	if err := daemon.Reload(newConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// After Reload, daemon.RegistryService will be changed which is useful
	// for registry communication in daemon.
	registries := daemon.RegistryService.ServiceConfig()

	// After Reload(), newConfig has come to registries.InsecureRegistryCIDRs and registries.IndexConfigs in daemon.
	// Then collect registries.InsecureRegistryCIDRs in dataMap.
	// When collecting, we need to convert CIDRS into string as a key,
	// while the times of key appears as value.
	dataMap := map[string]int***REMOVED******REMOVED***
	for _, value := range registries.InsecureRegistryCIDRs ***REMOVED***
		if _, ok := dataMap[value.String()]; !ok ***REMOVED***
			dataMap[value.String()] = 1
		***REMOVED*** else ***REMOVED***
			dataMap[value.String()]++
		***REMOVED***
	***REMOVED***

	for _, value := range registries.IndexConfigs ***REMOVED***
		if _, ok := dataMap[value.Name]; !ok ***REMOVED***
			dataMap[value.Name] = 1
		***REMOVED*** else ***REMOVED***
			dataMap[value.Name]++
		***REMOVED***
	***REMOVED***

	// Finally compare dataMap with the original insecureRegistries.
	// Each value in insecureRegistries should appear in daemon's insecure registries,
	// and each can only appear exactly ONCE.
	for _, r := range insecureRegistries ***REMOVED***
		if value, ok := dataMap[r]; !ok ***REMOVED***
			t.Fatalf("Expected daemon insecure registry %s, got none", r)
		***REMOVED*** else if value != 1 ***REMOVED***
			t.Fatalf("Expected only 1 daemon insecure registry %s, got %d", r, value)
		***REMOVED***
	***REMOVED***

	// assert if "10.10.1.22:5000" is removed when reloading
	if value, ok := dataMap["10.10.1.22:5000"]; ok ***REMOVED***
		t.Fatalf("Expected no insecure registry of 10.10.1.22:5000, got %d", value)
	***REMOVED***

	// assert if "docker2.com" is removed when reloading
	if value, ok := dataMap["docker2.com"]; ok ***REMOVED***
		t.Fatalf("Expected no insecure registry of docker2.com, got %d", value)
	***REMOVED***
***REMOVED***

func TestDaemonReloadNotAffectOthers(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	daemon.configStore = &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			Labels: []string***REMOVED***"foo:bar"***REMOVED***,
			Debug:  true,
		***REMOVED***,
	***REMOVED***

	valuesSets := make(map[string]interface***REMOVED******REMOVED***)
	valuesSets["labels"] = "foo:baz"
	newConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			Labels:    []string***REMOVED***"foo:baz"***REMOVED***,
			ValuesSet: valuesSets,
		***REMOVED***,
	***REMOVED***

	if err := daemon.Reload(newConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	label := daemon.configStore.Labels[0]
	if label != "foo:baz" ***REMOVED***
		t.Fatalf("Expected daemon label `foo:baz`, got %s", label)
	***REMOVED***
	debug := daemon.configStore.Debug
	if !debug ***REMOVED***
		t.Fatal("Expected debug 'enabled', got 'disabled'")
	***REMOVED***
***REMOVED***

func TestDaemonDiscoveryReload(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	daemon.configStore = &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ClusterStore:     "memory://127.0.0.1",
			ClusterAdvertise: "127.0.0.1:3333",
		***REMOVED***,
	***REMOVED***

	if err := daemon.initDiscovery(daemon.configStore); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected := discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "127.0.0.1", Port: "3333"***REMOVED***,
	***REMOVED***

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for discovery")
	case <-daemon.discoveryWatcher.ReadyCh():
	***REMOVED***

	stopCh := make(chan struct***REMOVED******REMOVED***)
	defer close(stopCh)
	ch, errCh := daemon.discoveryWatcher.Watch(stopCh)

	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("failed to get discovery advertisements in time")
	case e := <-ch:
		if !reflect.DeepEqual(e, expected) ***REMOVED***
			t.Fatalf("expected %v, got %v\n", expected, e)
		***REMOVED***
	case e := <-errCh:
		t.Fatal(e)
	***REMOVED***

	valuesSets := make(map[string]interface***REMOVED******REMOVED***)
	valuesSets["cluster-store"] = "memory://127.0.0.1:2222"
	valuesSets["cluster-advertise"] = "127.0.0.1:5555"
	newConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ClusterStore:     "memory://127.0.0.1:2222",
			ClusterAdvertise: "127.0.0.1:5555",
			ValuesSet:        valuesSets,
		***REMOVED***,
	***REMOVED***

	expected = discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "127.0.0.1", Port: "5555"***REMOVED***,
	***REMOVED***

	if err := daemon.Reload(newConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for discovery")
	case <-daemon.discoveryWatcher.ReadyCh():
	***REMOVED***

	ch, errCh = daemon.discoveryWatcher.Watch(stopCh)

	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("failed to get discovery advertisements in time")
	case e := <-ch:
		if !reflect.DeepEqual(e, expected) ***REMOVED***
			t.Fatalf("expected %v, got %v\n", expected, e)
		***REMOVED***
	case e := <-errCh:
		t.Fatal(e)
	***REMOVED***
***REMOVED***

func TestDaemonDiscoveryReloadFromEmptyDiscovery(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	daemon.configStore = &config.Config***REMOVED******REMOVED***

	valuesSet := make(map[string]interface***REMOVED******REMOVED***)
	valuesSet["cluster-store"] = "memory://127.0.0.1:2222"
	valuesSet["cluster-advertise"] = "127.0.0.1:5555"
	newConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ClusterStore:     "memory://127.0.0.1:2222",
			ClusterAdvertise: "127.0.0.1:5555",
			ValuesSet:        valuesSet,
		***REMOVED***,
	***REMOVED***

	expected := discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "127.0.0.1", Port: "5555"***REMOVED***,
	***REMOVED***

	if err := daemon.Reload(newConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for discovery")
	case <-daemon.discoveryWatcher.ReadyCh():
	***REMOVED***

	stopCh := make(chan struct***REMOVED******REMOVED***)
	defer close(stopCh)
	ch, errCh := daemon.discoveryWatcher.Watch(stopCh)

	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("failed to get discovery advertisements in time")
	case e := <-ch:
		if !reflect.DeepEqual(e, expected) ***REMOVED***
			t.Fatalf("expected %v, got %v\n", expected, e)
		***REMOVED***
	case e := <-errCh:
		t.Fatal(e)
	***REMOVED***
***REMOVED***

func TestDaemonDiscoveryReloadOnlyClusterAdvertise(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	daemon.configStore = &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ClusterStore: "memory://127.0.0.1",
		***REMOVED***,
	***REMOVED***
	valuesSets := make(map[string]interface***REMOVED******REMOVED***)
	valuesSets["cluster-advertise"] = "127.0.0.1:5555"
	newConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			ClusterAdvertise: "127.0.0.1:5555",
			ValuesSet:        valuesSets,
		***REMOVED***,
	***REMOVED***
	expected := discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "127.0.0.1", Port: "5555"***REMOVED***,
	***REMOVED***

	if err := daemon.Reload(newConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	select ***REMOVED***
	case <-daemon.discoveryWatcher.ReadyCh():
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for discovery")
	***REMOVED***
	stopCh := make(chan struct***REMOVED******REMOVED***)
	defer close(stopCh)
	ch, errCh := daemon.discoveryWatcher.Watch(stopCh)

	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("failed to get discovery advertisements in time")
	case e := <-ch:
		if !reflect.DeepEqual(e, expected) ***REMOVED***
			t.Fatalf("expected %v, got %v\n", expected, e)
		***REMOVED***
	case e := <-errCh:
		t.Fatal(e)
	***REMOVED***
***REMOVED***

func TestDaemonReloadNetworkDiagnosticPort(t *testing.T) ***REMOVED***
	daemon := &Daemon***REMOVED******REMOVED***
	daemon.configStore = &config.Config***REMOVED******REMOVED***

	valuesSet := make(map[string]interface***REMOVED******REMOVED***)
	valuesSet["network-diagnostic-port"] = 2000
	enableConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED***
			NetworkDiagnosticPort: 2000,
			ValuesSet:             valuesSet,
		***REMOVED***,
	***REMOVED***
	disableConfig := &config.Config***REMOVED***
		CommonConfig: config.CommonConfig***REMOVED******REMOVED***,
	***REMOVED***

	netOptions, err := daemon.networkOptions(enableConfig, nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	controller, err := libnetwork.New(netOptions...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	daemon.netController = controller

	// Enable/Disable the server for some iterations
	for i := 0; i < 10; i++ ***REMOVED***
		enableConfig.CommonConfig.NetworkDiagnosticPort++
		if err := daemon.Reload(enableConfig); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		// Check that the diagnose is enabled
		if !daemon.netController.IsDiagnoseEnabled() ***REMOVED***
			t.Fatalf("diagnosed should be enable")
		***REMOVED***

		// Reload
		if err := daemon.Reload(disableConfig); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		// Check that the diagnose is disabled
		if daemon.netController.IsDiagnoseEnabled() ***REMOVED***
			t.Fatalf("diagnosed should be disable")
		***REMOVED***
	***REMOVED***

	enableConfig.CommonConfig.NetworkDiagnosticPort++
	// 2 times the enable should not create problems
	if err := daemon.Reload(enableConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// Check that the diagnose is enabled
	if !daemon.netController.IsDiagnoseEnabled() ***REMOVED***
		t.Fatalf("diagnosed should be enable")
	***REMOVED***

	// Check that another reload does not cause issues
	if err := daemon.Reload(enableConfig); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// Check that the diagnose is enable
	if !daemon.netController.IsDiagnoseEnabled() ***REMOVED***
		t.Fatalf("diagnosed should be enable")
	***REMOVED***

***REMOVED***
