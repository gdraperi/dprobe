package distribution

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/internal/testutil"
	"github.com/opencontainers/go-digest"
)

// TestFixManifestLayers checks that fixManifestLayers removes a duplicate
// layer, and that it makes no changes to the manifest when called a second
// time, after the duplicate is removed.
func TestFixManifestLayers(t *testing.T) ***REMOVED***
	duplicateLayerManifest := schema1.Manifest***REMOVED***
		FSLayers: []schema1.FSLayer***REMOVED***
			***REMOVED***BlobSum: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
			***REMOVED***BlobSum: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
			***REMOVED***BlobSum: digest.Digest("sha256:86e0e091d0da6bde2456dbb48306f3956bbeb2eae1b5b9a43045843f69fe4aaa")***REMOVED***,
		***REMOVED***,
		History: []schema1.History***REMOVED***
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"3b38edc92eb7c074812e217b41a6ade66888531009d6286a6f5f36a06f9841b9\",\"parent\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:11.368300679Z\",\"container\":\"d91be3479d5b1e84b0c00d18eea9dc777ca0ad166d51174b24283e2e6f104253\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ENTRYPOINT [\\\"/go/bin/dnsdock\\\"]\"],\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":null,\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0***REMOVED***\n"***REMOVED***,
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"3b38edc92eb7c074812e217b41a6ade66888531009d6286a6f5f36a06f9841b9\",\"parent\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:11.368300679Z\",\"container\":\"d91be3479d5b1e84b0c00d18eea9dc777ca0ad166d51174b24283e2e6f104253\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ENTRYPOINT [\\\"/go/bin/dnsdock\\\"]\"],\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":null,\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0***REMOVED***\n"***REMOVED***,
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:07.568027497Z\",\"container\":\"fe9e5a5264a843c9292d17b736c92dd19bdb49986a8782d7389964ddaff887cc\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"cd /go/src/github.com/tonistiigi/dnsdock \\u0026\\u0026     go get -v github.com/tools/godep \\u0026\\u0026     godep restore \\u0026\\u0026     go install -ldflags \\\"-X main.version `git describe --tags HEAD``if [[ -n $(command git status --porcelain --untracked-files=no 2\\u003e/dev/null) ]]; then echo \\\"-dirty\\\"; fi`\\\" ./...\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/bash\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":118430532***REMOVED***\n"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	duplicateLayerManifestExpectedOutput := schema1.Manifest***REMOVED***
		FSLayers: []schema1.FSLayer***REMOVED***
			***REMOVED***BlobSum: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
			***REMOVED***BlobSum: digest.Digest("sha256:86e0e091d0da6bde2456dbb48306f3956bbeb2eae1b5b9a43045843f69fe4aaa")***REMOVED***,
		***REMOVED***,
		History: []schema1.History***REMOVED***
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"3b38edc92eb7c074812e217b41a6ade66888531009d6286a6f5f36a06f9841b9\",\"parent\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:11.368300679Z\",\"container\":\"d91be3479d5b1e84b0c00d18eea9dc777ca0ad166d51174b24283e2e6f104253\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ENTRYPOINT [\\\"/go/bin/dnsdock\\\"]\"],\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":null,\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0***REMOVED***\n"***REMOVED***,
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:07.568027497Z\",\"container\":\"fe9e5a5264a843c9292d17b736c92dd19bdb49986a8782d7389964ddaff887cc\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"cd /go/src/github.com/tonistiigi/dnsdock \\u0026\\u0026     go get -v github.com/tools/godep \\u0026\\u0026     godep restore \\u0026\\u0026     go install -ldflags \\\"-X main.version `git describe --tags HEAD``if [[ -n $(command git status --porcelain --untracked-files=no 2\\u003e/dev/null) ]]; then echo \\\"-dirty\\\"; fi`\\\" ./...\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/bash\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":118430532***REMOVED***\n"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	if err := fixManifestLayers(&duplicateLayerManifest); err != nil ***REMOVED***
		t.Fatalf("unexpected error from fixManifestLayers: %v", err)
	***REMOVED***

	if !reflect.DeepEqual(duplicateLayerManifest, duplicateLayerManifestExpectedOutput) ***REMOVED***
		t.Fatal("incorrect output from fixManifestLayers on duplicate layer manifest")
	***REMOVED***

	// Run fixManifestLayers again and confirm that it doesn't change the
	// manifest (which no longer has duplicate layers).
	if err := fixManifestLayers(&duplicateLayerManifest); err != nil ***REMOVED***
		t.Fatalf("unexpected error from fixManifestLayers: %v", err)
	***REMOVED***

	if !reflect.DeepEqual(duplicateLayerManifest, duplicateLayerManifestExpectedOutput) ***REMOVED***
		t.Fatal("incorrect output from fixManifestLayers on duplicate layer manifest (second pass)")
	***REMOVED***
***REMOVED***

// TestFixManifestLayersBaseLayerParent makes sure that fixManifestLayers fails
// if the base layer configuration specifies a parent.
func TestFixManifestLayersBaseLayerParent(t *testing.T) ***REMOVED***
	// TODO Windows: Fix this unit text
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Needs fixing on Windows")
	***REMOVED***
	duplicateLayerManifest := schema1.Manifest***REMOVED***
		FSLayers: []schema1.FSLayer***REMOVED***
			***REMOVED***BlobSum: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
			***REMOVED***BlobSum: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
			***REMOVED***BlobSum: digest.Digest("sha256:86e0e091d0da6bde2456dbb48306f3956bbeb2eae1b5b9a43045843f69fe4aaa")***REMOVED***,
		***REMOVED***,
		History: []schema1.History***REMOVED***
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"3b38edc92eb7c074812e217b41a6ade66888531009d6286a6f5f36a06f9841b9\",\"parent\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:11.368300679Z\",\"container\":\"d91be3479d5b1e84b0c00d18eea9dc777ca0ad166d51174b24283e2e6f104253\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ENTRYPOINT [\\\"/go/bin/dnsdock\\\"]\"],\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":null,\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0***REMOVED***\n"***REMOVED***,
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"3b38edc92eb7c074812e217b41a6ade66888531009d6286a6f5f36a06f9841b9\",\"parent\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:11.368300679Z\",\"container\":\"d91be3479d5b1e84b0c00d18eea9dc777ca0ad166d51174b24283e2e6f104253\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ENTRYPOINT [\\\"/go/bin/dnsdock\\\"]\"],\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":null,\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0***REMOVED***\n"***REMOVED***,
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"parent\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"created\":\"2015-08-19T16:49:07.568027497Z\",\"container\":\"fe9e5a5264a843c9292d17b736c92dd19bdb49986a8782d7389964ddaff887cc\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"cd /go/src/github.com/tonistiigi/dnsdock \\u0026\\u0026     go get -v github.com/tools/godep \\u0026\\u0026     godep restore \\u0026\\u0026     go install -ldflags \\\"-X main.version `git describe --tags HEAD``if [[ -n $(command git status --porcelain --untracked-files=no 2\\u003e/dev/null) ]]; then echo \\\"-dirty\\\"; fi`\\\" ./...\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/bash\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":118430532***REMOVED***\n"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	if err := fixManifestLayers(&duplicateLayerManifest); err == nil || !strings.Contains(err.Error(), "invalid parent ID in the base layer of the image") ***REMOVED***
		t.Fatalf("expected an invalid parent ID error from fixManifestLayers")
	***REMOVED***
***REMOVED***

// TestFixManifestLayersBadParent makes sure that fixManifestLayers fails
// if an image configuration specifies a parent that doesn't directly follow
// that (deduplicated) image in the image history.
func TestFixManifestLayersBadParent(t *testing.T) ***REMOVED***
	duplicateLayerManifest := schema1.Manifest***REMOVED***
		FSLayers: []schema1.FSLayer***REMOVED***
			***REMOVED***BlobSum: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
			***REMOVED***BlobSum: digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")***REMOVED***,
			***REMOVED***BlobSum: digest.Digest("sha256:86e0e091d0da6bde2456dbb48306f3956bbeb2eae1b5b9a43045843f69fe4aaa")***REMOVED***,
		***REMOVED***,
		History: []schema1.History***REMOVED***
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"3b38edc92eb7c074812e217b41a6ade66888531009d6286a6f5f36a06f9841b9\",\"parent\":\"ac3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:11.368300679Z\",\"container\":\"d91be3479d5b1e84b0c00d18eea9dc777ca0ad166d51174b24283e2e6f104253\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ENTRYPOINT [\\\"/go/bin/dnsdock\\\"]\"],\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":null,\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0***REMOVED***\n"***REMOVED***,
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"3b38edc92eb7c074812e217b41a6ade66888531009d6286a6f5f36a06f9841b9\",\"parent\":\"ac3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:11.368300679Z\",\"container\":\"d91be3479d5b1e84b0c00d18eea9dc777ca0ad166d51174b24283e2e6f104253\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ENTRYPOINT [\\\"/go/bin/dnsdock\\\"]\"],\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":null,\"Image\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":[\"/go/bin/dnsdock\"],\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0***REMOVED***\n"***REMOVED***,
			***REMOVED***V1Compatibility: "***REMOVED***\"id\":\"ec3025ca8cc9bcab039e193e20ec647c2da3c53a74020f2ba611601f9b2c6c02\",\"created\":\"2015-08-19T16:49:07.568027497Z\",\"container\":\"fe9e5a5264a843c9292d17b736c92dd19bdb49986a8782d7389964ddaff887cc\",\"container_config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"cd /go/src/github.com/tonistiigi/dnsdock \\u0026\\u0026     go get -v github.com/tools/godep \\u0026\\u0026     godep restore \\u0026\\u0026     go install -ldflags \\\"-X main.version `git describe --tags HEAD``if [[ -n $(command git status --porcelain --untracked-files=no 2\\u003e/dev/null) ]]; then echo \\\"-dirty\\\"; fi`\\\" ./...\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"docker_version\":\"1.6.2\",\"config\":***REMOVED***\"Hostname\":\"03797203757d\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\",\"GOLANG_VERSION=1.4.1\",\"GOPATH=/go\"],\"Cmd\":[\"/bin/bash\"],\"Image\":\"e3b0ff09e647595dafee15c54cd632c900df9e82b1d4d313b1e20639a1461779\",\"Volumes\":null,\"WorkingDir\":\"/go\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"Labels\":***REMOVED******REMOVED******REMOVED***,\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":118430532***REMOVED***\n"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	err := fixManifestLayers(&duplicateLayerManifest)
	testutil.ErrorContains(t, err, "invalid parent ID")
***REMOVED***

// TestValidateManifest verifies the validateManifest function
func TestValidateManifest(t *testing.T) ***REMOVED***
	// TODO Windows: Fix this unit text
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Needs fixing on Windows")
	***REMOVED***
	expectedDigest, err := reference.ParseNormalizedNamed("repo@sha256:02fee8c3220ba806531f606525eceb83f4feb654f62b207191b1c9209188dedd")
	if err != nil ***REMOVED***
		t.Fatal("could not parse reference")
	***REMOVED***
	expectedFSLayer0 := digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")

	// Good manifest

	goodManifestBytes, err := ioutil.ReadFile("fixtures/validate_manifest/good_manifest")
	if err != nil ***REMOVED***
		t.Fatal("error reading fixture:", err)
	***REMOVED***

	var goodSignedManifest schema1.SignedManifest
	err = json.Unmarshal(goodManifestBytes, &goodSignedManifest)
	if err != nil ***REMOVED***
		t.Fatal("error unmarshaling manifest:", err)
	***REMOVED***

	verifiedManifest, err := verifySchema1Manifest(&goodSignedManifest, expectedDigest)
	if err != nil ***REMOVED***
		t.Fatal("validateManifest failed:", err)
	***REMOVED***

	if verifiedManifest.FSLayers[0].BlobSum != expectedFSLayer0 ***REMOVED***
		t.Fatal("unexpected FSLayer in good manifest")
	***REMOVED***

	// "Extra data" manifest

	extraDataManifestBytes, err := ioutil.ReadFile("fixtures/validate_manifest/extra_data_manifest")
	if err != nil ***REMOVED***
		t.Fatal("error reading fixture:", err)
	***REMOVED***

	var extraDataSignedManifest schema1.SignedManifest
	err = json.Unmarshal(extraDataManifestBytes, &extraDataSignedManifest)
	if err != nil ***REMOVED***
		t.Fatal("error unmarshaling manifest:", err)
	***REMOVED***

	verifiedManifest, err = verifySchema1Manifest(&extraDataSignedManifest, expectedDigest)
	if err != nil ***REMOVED***
		t.Fatal("validateManifest failed:", err)
	***REMOVED***

	if verifiedManifest.FSLayers[0].BlobSum != expectedFSLayer0 ***REMOVED***
		t.Fatal("unexpected FSLayer in extra data manifest")
	***REMOVED***

	// Bad manifest

	badManifestBytes, err := ioutil.ReadFile("fixtures/validate_manifest/bad_manifest")
	if err != nil ***REMOVED***
		t.Fatal("error reading fixture:", err)
	***REMOVED***

	var badSignedManifest schema1.SignedManifest
	err = json.Unmarshal(badManifestBytes, &badSignedManifest)
	if err != nil ***REMOVED***
		t.Fatal("error unmarshaling manifest:", err)
	***REMOVED***

	verifiedManifest, err = verifySchema1Manifest(&badSignedManifest, expectedDigest)
	if err == nil || !strings.HasPrefix(err.Error(), "image verification failed for digest") ***REMOVED***
		t.Fatal("expected validateManifest to fail with digest error")
	***REMOVED***
***REMOVED***
