package v1

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/opencontainers/go-digest"
)

func TestMigrateRefs(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "migrate-tags")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	ioutil.WriteFile(filepath.Join(tmpdir, "repositories-generic"), []byte(`***REMOVED***"Repositories":***REMOVED***"busybox":***REMOVED***"latest":"b3ca410aa2c115c05969a7b2c8cf8a9fcf62c1340ed6a601c9ee50df337ec108","sha256:16a2a52884c2a9481ed267c2d46483eac7693b813a63132368ab098a71303f8a":"b3ca410aa2c115c05969a7b2c8cf8a9fcf62c1340ed6a601c9ee50df337ec108"***REMOVED***,"registry":***REMOVED***"2":"5d165b8e4b203685301c815e95663231691d383fd5e3d3185d1ce3f8dddead3d","latest":"8d5547a9f329b1d3f93198cd661fb5117e5a96b721c5cf9a2c389e7dd4877128"***REMOVED******REMOVED******REMOVED***`), 0600)

	ta := &mockTagAdder***REMOVED******REMOVED***
	err = migrateRefs(tmpdir, "generic", ta, map[string]image.ID***REMOVED***
		"5d165b8e4b203685301c815e95663231691d383fd5e3d3185d1ce3f8dddead3d": image.ID("sha256:2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"),
		"b3ca410aa2c115c05969a7b2c8cf8a9fcf62c1340ed6a601c9ee50df337ec108": image.ID("sha256:fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9"),
		"abcdef3434c115c05969a7b2c8cf8a9fcf62c1340ed6a601c9ee50df337ec108": image.ID("sha256:56434342345ae68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"),
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected := map[string]string***REMOVED***
		"docker.io/library/busybox:latest":                                                                  "sha256:fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9",
		"docker.io/library/busybox@sha256:16a2a52884c2a9481ed267c2d46483eac7693b813a63132368ab098a71303f8a": "sha256:fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9",
		"docker.io/library/registry:2":                                                                      "sha256:2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae",
	***REMOVED***

	if !reflect.DeepEqual(expected, ta.refs) ***REMOVED***
		t.Fatalf("Invalid migrated tags: expected %q, got %q", expected, ta.refs)
	***REMOVED***

	// second migration is no-op
	ioutil.WriteFile(filepath.Join(tmpdir, "repositories-generic"), []byte(`***REMOVED***"Repositories":***REMOVED***"busybox":***REMOVED***"latest":"b3ca410aa2c115c05969a7b2c8cf8a9fcf62c1340ed6a601c9ee50df337ec108"`), 0600)
	err = migrateRefs(tmpdir, "generic", ta, map[string]image.ID***REMOVED***
		"b3ca410aa2c115c05969a7b2c8cf8a9fcf62c1340ed6a601c9ee50df337ec108": image.ID("sha256:fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9"),
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(expected, ta.refs) ***REMOVED***
		t.Fatalf("Invalid migrated tags: expected %q, got %q", expected, ta.refs)
	***REMOVED***
***REMOVED***

func TestMigrateContainers(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	if runtime.GOARCH != "amd64" ***REMOVED***
		t.Skip("Test tailored to amd64 architecture")
	***REMOVED***
	tmpdir, err := ioutil.TempDir("", "migrate-containers")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	err = addContainer(tmpdir, `***REMOVED***"State":***REMOVED***"Running":false,"Paused":false,"Restarting":false,"OOMKilled":false,"Dead":false,"Pid":0,"ExitCode":0,"Error":"","StartedAt":"2015-11-10T21:42:40.604267436Z","FinishedAt":"2015-11-10T21:42:41.869265487Z"***REMOVED***,"ID":"f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c","Created":"2015-11-10T21:42:40.433831551Z","Path":"sh","Args":[],"Config":***REMOVED***"Hostname":"f780ee3f80e6","Domainname":"","User":"","AttachStdin":true,"AttachStdout":true,"AttachStderr":true,"Tty":true,"OpenStdin":true,"StdinOnce":true,"Env":null,"Cmd":["sh"],"Image":"busybox","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":***REMOVED******REMOVED******REMOVED***,"Image":"2c5ac3f849df8627fcf2822727f87c57f38b7129d3604fbc11d861fe856ff093","NetworkSettings":***REMOVED***"Bridge":"","EndpointID":"","Gateway":"","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"HairpinMode":false,"IPAddress":"","IPPrefixLen":0,"IPv6Gateway":"","LinkLocalIPv6Address":"","LinkLocalIPv6PrefixLen":0,"MacAddress":"","NetworkID":"","PortMapping":null,"Ports":null,"SandboxKey":"","SecondaryIPAddresses":null,"SecondaryIPv6Addresses":null***REMOVED***,"ResolvConfPath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/resolv.conf","HostnamePath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/hostname","HostsPath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/hosts","LogPath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c-json.log","Name":"/determined_euclid","Driver":"overlay","ExecDriver":"native-0.2","MountLabel":"","ProcessLabel":"","RestartCount":0,"UpdateDns":false,"HasBeenStartedBefore":false,"MountPoints":***REMOVED******REMOVED***,"Volumes":***REMOVED******REMOVED***,"VolumesRW":***REMOVED******REMOVED***,"AppArmorProfile":""***REMOVED***`)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// container with invalid image
	err = addContainer(tmpdir, `***REMOVED***"State":***REMOVED***"Running":false,"Paused":false,"Restarting":false,"OOMKilled":false,"Dead":false,"Pid":0,"ExitCode":0,"Error":"","StartedAt":"2015-11-10T21:42:40.604267436Z","FinishedAt":"2015-11-10T21:42:41.869265487Z"***REMOVED***,"ID":"e780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c","Created":"2015-11-10T21:42:40.433831551Z","Path":"sh","Args":[],"Config":***REMOVED***"Hostname":"f780ee3f80e6","Domainname":"","User":"","AttachStdin":true,"AttachStdout":true,"AttachStderr":true,"Tty":true,"OpenStdin":true,"StdinOnce":true,"Env":null,"Cmd":["sh"],"Image":"busybox","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":***REMOVED******REMOVED******REMOVED***,"Image":"4c5ac3f849df8627fcf2822727f87c57f38b7129d3604fbc11d861fe856ff093","NetworkSettings":***REMOVED***"Bridge":"","EndpointID":"","Gateway":"","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"HairpinMode":false,"IPAddress":"","IPPrefixLen":0,"IPv6Gateway":"","LinkLocalIPv6Address":"","LinkLocalIPv6PrefixLen":0,"MacAddress":"","NetworkID":"","PortMapping":null,"Ports":null,"SandboxKey":"","SecondaryIPAddresses":null,"SecondaryIPv6Addresses":null***REMOVED***,"ResolvConfPath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/resolv.conf","HostnamePath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/hostname","HostsPath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/hosts","LogPath":"/var/lib/docker/containers/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c/f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c-json.log","Name":"/determined_euclid","Driver":"overlay","ExecDriver":"native-0.2","MountLabel":"","ProcessLabel":"","RestartCount":0,"UpdateDns":false,"HasBeenStartedBefore":false,"MountPoints":***REMOVED******REMOVED***,"Volumes":***REMOVED******REMOVED***,"VolumesRW":***REMOVED******REMOVED***,"AppArmorProfile":""***REMOVED***`)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ifs, err := image.NewFSStoreBackend(filepath.Join(tmpdir, "imagedb"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ls := &mockMounter***REMOVED******REMOVED***
	mmMap := make(map[string]image.LayerGetReleaser)
	mmMap[runtime.GOOS] = ls
	is, err := image.NewImageStore(ifs, mmMap)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	imgID, err := is.Create([]byte(`***REMOVED***"architecture":"amd64","config":***REMOVED***"AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"Cmd":["sh"],"Entrypoint":null,"Env":null,"Hostname":"23304fc829f9","Image":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","Labels":null,"OnBuild":null,"OpenStdin":false,"StdinOnce":false,"Tty":false,"Volumes":null,"WorkingDir":"","Domainname":"","User":""***REMOVED***,"container":"349b014153779e30093d94f6df2a43c7a0a164e05aa207389917b540add39b51","container_config":***REMOVED***"AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"Cmd":["/bin/sh","-c","#(nop) CMD [\"sh\"]"],"Entrypoint":null,"Env":null,"Hostname":"23304fc829f9","Image":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","Labels":null,"OnBuild":null,"OpenStdin":false,"StdinOnce":false,"Tty":false,"Volumes":null,"WorkingDir":"","Domainname":"","User":""***REMOVED***,"created":"2015-10-31T22:22:55.613815829Z","docker_version":"1.8.2","history":[***REMOVED***"created":"2015-10-31T22:22:54.690851953Z","created_by":"/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"***REMOVED***,***REMOVED***"created":"2015-10-31T22:22:55.613815829Z","created_by":"/bin/sh -c #(nop) CMD [\"sh\"]"***REMOVED***],"os":"linux","rootfs":***REMOVED***"type":"layers","diff_ids":["sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1","sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"]***REMOVED******REMOVED***`))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = migrateContainers(tmpdir, ls, is, map[string]image.ID***REMOVED***
		"2c5ac3f849df8627fcf2822727f87c57f38b7129d3604fbc11d861fe856ff093": imgID,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected := []mountInfo***REMOVED******REMOVED***
		"f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c",
		"f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c",
		"sha256:c3191d32a37d7159b2e30830937d2e30268ad6c375a773a8994911a3aba9b93f",
	***REMOVED******REMOVED***
	if !reflect.DeepEqual(expected, ls.mounts) ***REMOVED***
		t.Fatalf("invalid mounts: expected %q, got %q", expected, ls.mounts)
	***REMOVED***

	if actual, expected := ls.count, 0; actual != expected ***REMOVED***
		t.Fatalf("invalid active mounts: expected %d, got %d", expected, actual)
	***REMOVED***

	config2, err := ioutil.ReadFile(filepath.Join(tmpdir, "containers", "f780ee3f80e66e9b432a57049597118a66aab8932be88e5628d4c824edbee37c", "config.v2.json"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	var config struct***REMOVED*** Image string ***REMOVED***
	err = json.Unmarshal(config2, &config)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if actual, expected := config.Image, string(imgID); actual != expected ***REMOVED***
		t.Fatalf("invalid image pointer in migrated config: expected %q, got %q", expected, actual)
	***REMOVED***

***REMOVED***

func TestMigrateImages(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	if runtime.GOARCH != "amd64" ***REMOVED***
		t.Skip("Test tailored to amd64 architecture")
	***REMOVED***
	tmpdir, err := ioutil.TempDir("", "migrate-images")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	// busybox from 1.9
	id1, err := addImage(tmpdir, `***REMOVED***"architecture":"amd64","config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":null,"Image":"","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":null***REMOVED***,"container":"23304fc829f9b9349416f6eb1afec162907eba3a328f51d53a17f8986f865d65","container_config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":["/bin/sh","-c","#(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"],"Image":"","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":null***REMOVED***,"created":"2015-10-31T22:22:54.690851953Z","docker_version":"1.8.2","layer_id":"sha256:55dc925c23d1ed82551fd018c27ac3ee731377b6bad3963a2a4e76e753d70e57","os":"linux"***REMOVED***`, "", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	id2, err := addImage(tmpdir, `***REMOVED***"architecture":"amd64","config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":["sh"],"Image":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":null***REMOVED***,"container":"349b014153779e30093d94f6df2a43c7a0a164e05aa207389917b540add39b51","container_config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":["/bin/sh","-c","#(nop) CMD [\"sh\"]"],"Image":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":null***REMOVED***,"created":"2015-10-31T22:22:55.613815829Z","docker_version":"1.8.2","layer_id":"sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4","os":"linux","parent_id":"sha256:039b63dd2cbaa10d6015ea574392530571ed8d7b174090f032211285a71881d0"***REMOVED***`, id1, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ifs, err := image.NewFSStoreBackend(filepath.Join(tmpdir, "imagedb"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ls := &mockRegistrar***REMOVED******REMOVED***
	mrMap := make(map[string]image.LayerGetReleaser)
	mrMap[runtime.GOOS] = ls
	is, err := image.NewImageStore(ifs, mrMap)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ms, err := metadata.NewFSMetadataStore(filepath.Join(tmpdir, "distribution"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	mappings := make(map[string]image.ID)

	err = migrateImages(tmpdir, ls, is, ms, mappings)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected := map[string]image.ID***REMOVED***
		id1: image.ID("sha256:ca406eaf9c26898414ff5b7b3a023c33310759d6203be0663dbf1b3a712f432d"),
		id2: image.ID("sha256:a488bec94bb96b26a968f913d25ef7d8d204d727ca328b52b4b059c7d03260b6"),
	***REMOVED***

	if !reflect.DeepEqual(mappings, expected) ***REMOVED***
		t.Fatalf("invalid image mappings: expected %q, got %q", expected, mappings)
	***REMOVED***

	if actual, expected := ls.count, 2; actual != expected ***REMOVED***
		t.Fatalf("invalid register count: expected %q, got %q", expected, actual)
	***REMOVED***
	ls.count = 0

	// next images are busybox from 1.8.2
	_, err = addImage(tmpdir, `***REMOVED***"id":"17583c7dd0dae6244203b8029733bdb7d17fccbb2b5d93e2b24cf48b8bfd06e2","parent":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","created":"2015-10-31T22:22:55.613815829Z","container":"349b014153779e30093d94f6df2a43c7a0a164e05aa207389917b540add39b51","container_config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"ExposedPorts":null,"PublishService":"","Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":["/bin/sh","-c","#(nop) CMD [\"sh\"]"],"Image":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","Volumes":null,"VolumeDriver":"","WorkingDir":"","Entrypoint":null,"NetworkDisabled":false,"MacAddress":"","OnBuild":null,"Labels":null***REMOVED***,"docker_version":"1.8.2","config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"ExposedPorts":null,"PublishService":"","Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":["sh"],"Image":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","Volumes":null,"VolumeDriver":"","WorkingDir":"","Entrypoint":null,"NetworkDisabled":false,"MacAddress":"","OnBuild":null,"Labels":null***REMOVED***,"architecture":"amd64","os":"linux","Size":0***REMOVED***`, "", "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	_, err = addImage(tmpdir, `***REMOVED***"id":"d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498","created":"2015-10-31T22:22:54.690851953Z","container":"23304fc829f9b9349416f6eb1afec162907eba3a328f51d53a17f8986f865d65","container_config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"ExposedPorts":null,"PublishService":"","Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":["/bin/sh","-c","#(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"],"Image":"","Volumes":null,"VolumeDriver":"","WorkingDir":"","Entrypoint":null,"NetworkDisabled":false,"MacAddress":"","OnBuild":null,"Labels":null***REMOVED***,"docker_version":"1.8.2","config":***REMOVED***"Hostname":"23304fc829f9","Domainname":"","User":"","AttachStdin":false,"AttachStdout":false,"AttachStderr":false,"ExposedPorts":null,"PublishService":"","Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":null,"Cmd":null,"Image":"","Volumes":null,"VolumeDriver":"","WorkingDir":"","Entrypoint":null,"NetworkDisabled":false,"MacAddress":"","OnBuild":null,"Labels":null***REMOVED***,"architecture":"amd64","os":"linux","Size":1108935***REMOVED***`, "", "sha256:55dc925c23d1ed82551fd018c27ac3ee731377b6bad3963a2a4e76e753d70e57")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = migrateImages(tmpdir, ls, is, ms, mappings)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected["d1592a710ac323612bd786fa8ac20727c58d8a67847e5a65177c594f43919498"] = image.ID("sha256:c091bb33854e57e6902b74c08719856d30b5593c7db6143b2b48376b8a588395")
	expected["17583c7dd0dae6244203b8029733bdb7d17fccbb2b5d93e2b24cf48b8bfd06e2"] = image.ID("sha256:d963020e755ff2715b936065949472c1f8a6300144b922992a1a421999e71f07")

	if actual, expected := ls.count, 2; actual != expected ***REMOVED***
		t.Fatalf("invalid register count: expected %q, got %q", expected, actual)
	***REMOVED***

	v2MetadataService := metadata.NewV2MetadataService(ms)
	receivedMetadata, err := v2MetadataService.GetMetadata(layer.EmptyLayer.DiffID())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expectedMetadata := []metadata.V2Metadata***REMOVED***
		***REMOVED***Digest: digest.Digest("sha256:55dc925c23d1ed82551fd018c27ac3ee731377b6bad3963a2a4e76e753d70e57")***REMOVED***,
		***REMOVED***Digest: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
	***REMOVED***

	if !reflect.DeepEqual(expectedMetadata, receivedMetadata) ***REMOVED***
		t.Fatalf("invalid metadata: expected %q, got %q", expectedMetadata, receivedMetadata)
	***REMOVED***

***REMOVED***

func TestMigrateUnsupported(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "migrate-empty")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	err = os.MkdirAll(filepath.Join(tmpdir, "graph"), 0700)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = Migrate(tmpdir, "generic", nil, nil, nil, nil)
	if err != errUnsupported ***REMOVED***
		t.Fatalf("expected unsupported error, got %q", err)
	***REMOVED***
***REMOVED***

func TestMigrateEmptyDir(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "migrate-empty")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	err = Migrate(tmpdir, "generic", nil, nil, nil, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func addImage(dest, jsonConfig, parent, checksum string) (string, error) ***REMOVED***
	var config struct***REMOVED*** ID string ***REMOVED***
	if err := json.Unmarshal([]byte(jsonConfig), &config); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if config.ID == "" ***REMOVED***
		b := make([]byte, 32)
		rand.Read(b)
		config.ID = hex.EncodeToString(b)
	***REMOVED***
	contDir := filepath.Join(dest, "graph", config.ID)
	if err := os.MkdirAll(contDir, 0700); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(contDir, "json"), []byte(jsonConfig), 0600); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if checksum != "" ***REMOVED***
		if err := ioutil.WriteFile(filepath.Join(contDir, "checksum"), []byte(checksum), 0600); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(contDir, ".migration-diffid"), []byte(layer.EmptyLayer.DiffID()), 0600); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if err := ioutil.WriteFile(filepath.Join(contDir, ".migration-size"), []byte("0"), 0600); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if parent != "" ***REMOVED***
		if err := ioutil.WriteFile(filepath.Join(contDir, "parent"), []byte(parent), 0600); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	if checksum != "" ***REMOVED***
		if err := ioutil.WriteFile(filepath.Join(contDir, "checksum"), []byte(checksum), 0600); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	return config.ID, nil
***REMOVED***

func addContainer(dest, jsonConfig string) error ***REMOVED***
	var config struct***REMOVED*** ID string ***REMOVED***
	if err := json.Unmarshal([]byte(jsonConfig), &config); err != nil ***REMOVED***
		return err
	***REMOVED***
	contDir := filepath.Join(dest, "containers", config.ID)
	if err := os.MkdirAll(contDir, 0700); err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(filepath.Join(contDir, "config.json"), []byte(jsonConfig), 0600)
***REMOVED***

type mockTagAdder struct ***REMOVED***
	refs map[string]string
***REMOVED***

func (t *mockTagAdder) AddTag(ref reference.Named, id digest.Digest, force bool) error ***REMOVED***
	if t.refs == nil ***REMOVED***
		t.refs = make(map[string]string)
	***REMOVED***
	t.refs[ref.String()] = id.String()
	return nil
***REMOVED***
func (t *mockTagAdder) AddDigest(ref reference.Canonical, id digest.Digest, force bool) error ***REMOVED***
	return t.AddTag(ref, id, force)
***REMOVED***

type mockRegistrar struct ***REMOVED***
	layers map[layer.ChainID]*mockLayer
	count  int
***REMOVED***

func (r *mockRegistrar) RegisterByGraphID(graphID string, parent layer.ChainID, diffID layer.DiffID, tarDataFile string, size int64) (layer.Layer, error) ***REMOVED***
	r.count++
	l := &mockLayer***REMOVED******REMOVED***
	if parent != "" ***REMOVED***
		p, exists := r.layers[parent]
		if !exists ***REMOVED***
			return nil, fmt.Errorf("invalid parent %q", parent)
		***REMOVED***
		l.parent = p
		l.diffIDs = append(l.diffIDs, p.diffIDs...)
	***REMOVED***
	l.diffIDs = append(l.diffIDs, diffID)
	if r.layers == nil ***REMOVED***
		r.layers = make(map[layer.ChainID]*mockLayer)
	***REMOVED***
	r.layers[l.ChainID()] = l
	return l, nil
***REMOVED***
func (r *mockRegistrar) Release(l layer.Layer) ([]layer.Metadata, error) ***REMOVED***
	return nil, nil
***REMOVED***
func (r *mockRegistrar) Get(layer.ChainID) (layer.Layer, error) ***REMOVED***
	return nil, nil
***REMOVED***

type mountInfo struct ***REMOVED***
	name, graphID, parent string
***REMOVED***
type mockMounter struct ***REMOVED***
	mounts []mountInfo
	count  int
***REMOVED***

func (r *mockMounter) CreateRWLayerByGraphID(name string, graphID string, parent layer.ChainID) error ***REMOVED***
	r.mounts = append(r.mounts, mountInfo***REMOVED***name, graphID, string(parent)***REMOVED***)
	return nil
***REMOVED***
func (r *mockMounter) Unmount(string) error ***REMOVED***
	r.count--
	return nil
***REMOVED***
func (r *mockMounter) Get(layer.ChainID) (layer.Layer, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (r *mockMounter) Release(layer.Layer) ([]layer.Metadata, error) ***REMOVED***
	return nil, nil
***REMOVED***

type mockLayer struct ***REMOVED***
	diffIDs []layer.DiffID
	parent  *mockLayer
***REMOVED***

func (l *mockLayer) TarStream() (io.ReadCloser, error) ***REMOVED***
	return nil, nil
***REMOVED***
func (l *mockLayer) TarStreamFrom(layer.ChainID) (io.ReadCloser, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (l *mockLayer) ChainID() layer.ChainID ***REMOVED***
	return layer.CreateChainID(l.diffIDs)
***REMOVED***

func (l *mockLayer) DiffID() layer.DiffID ***REMOVED***
	return l.diffIDs[len(l.diffIDs)-1]
***REMOVED***

func (l *mockLayer) Parent() layer.Layer ***REMOVED***
	if l.parent == nil ***REMOVED***
		return nil
	***REMOVED***
	return l.parent
***REMOVED***

func (l *mockLayer) Size() (int64, error) ***REMOVED***
	return 0, nil
***REMOVED***

func (l *mockLayer) DiffSize() (int64, error) ***REMOVED***
	return 0, nil
***REMOVED***

func (l *mockLayer) Metadata() (map[string]string, error) ***REMOVED***
	return nil, nil
***REMOVED***
