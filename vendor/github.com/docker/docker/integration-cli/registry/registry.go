package registry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opencontainers/go-digest"
)

const (
	v2binary        = "registry-v2"
	v2binarySchema1 = "registry-v2-schema1"
)

type testingT interface ***REMOVED***
	logT
	Fatal(...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type logT interface ***REMOVED***
	Logf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// V2 represent a registry version 2
type V2 struct ***REMOVED***
	cmd         *exec.Cmd
	registryURL string
	dir         string
	auth        string
	username    string
	password    string
	email       string
***REMOVED***

// NewV2 creates a v2 registry server
func NewV2(schema1 bool, auth, tokenURL, registryURL string) (*V2, error) ***REMOVED***
	tmp, err := ioutil.TempDir("", "registry-test-")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	template := `version: 0.1
loglevel: debug
storage:
    filesystem:
        rootdirectory: %s
http:
    addr: %s
%s`
	var (
		authTemplate string
		username     string
		password     string
		email        string
	)
	switch auth ***REMOVED***
	case "htpasswd":
		htpasswdPath := filepath.Join(tmp, "htpasswd")
		// generated with: htpasswd -Bbn testuser testpassword
		userpasswd := "testuser:$2y$05$sBsSqk0OpSD1uTZkHXc4FeJ0Z70wLQdAX/82UiHuQOKbNbBrzs63m"
		username = "testuser"
		password = "testpassword"
		email = "test@test.org"
		if err := ioutil.WriteFile(htpasswdPath, []byte(userpasswd), os.FileMode(0644)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		authTemplate = fmt.Sprintf(`auth:
    htpasswd:
        realm: basic-realm
        path: %s
`, htpasswdPath)
	case "token":
		authTemplate = fmt.Sprintf(`auth:
    token:
        realm: %s
        service: "registry"
        issuer: "auth-registry"
        rootcertbundle: "fixtures/registry/cert.pem"
`, tokenURL)
	***REMOVED***

	confPath := filepath.Join(tmp, "config.yaml")
	config, err := os.Create(confPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer config.Close()

	if _, err := fmt.Fprintf(config, template, tmp, registryURL, authTemplate); err != nil ***REMOVED***
		os.RemoveAll(tmp)
		return nil, err
	***REMOVED***

	binary := v2binary
	if schema1 ***REMOVED***
		binary = v2binarySchema1
	***REMOVED***
	cmd := exec.Command(binary, confPath)
	if err := cmd.Start(); err != nil ***REMOVED***
		os.RemoveAll(tmp)
		return nil, err
	***REMOVED***
	return &V2***REMOVED***
		cmd:         cmd,
		dir:         tmp,
		auth:        auth,
		username:    username,
		password:    password,
		email:       email,
		registryURL: registryURL,
	***REMOVED***, nil
***REMOVED***

// Ping sends an http request to the current registry, and fail if it doesn't respond correctly
func (r *V2) Ping() error ***REMOVED***
	// We always ping through HTTP for our test registry.
	resp, err := http.Get(fmt.Sprintf("http://%s/v2/", r.registryURL))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	resp.Body.Close()

	fail := resp.StatusCode != http.StatusOK
	if r.auth != "" ***REMOVED***
		// unauthorized is a _good_ status when pinging v2/ and it needs auth
		fail = fail && resp.StatusCode != http.StatusUnauthorized
	***REMOVED***
	if fail ***REMOVED***
		return fmt.Errorf("registry ping replied with an unexpected status code %d", resp.StatusCode)
	***REMOVED***
	return nil
***REMOVED***

// Close kills the registry server
func (r *V2) Close() ***REMOVED***
	r.cmd.Process.Kill()
	r.cmd.Process.Wait()
	os.RemoveAll(r.dir)
***REMOVED***

func (r *V2) getBlobFilename(blobDigest digest.Digest) string ***REMOVED***
	// Split the digest into its algorithm and hex components.
	dgstAlg, dgstHex := blobDigest.Algorithm(), blobDigest.Hex()

	// The path to the target blob data looks something like:
	//   baseDir + "docker/registry/v2/blobs/sha256/a3/a3ed...46d4/data"
	return fmt.Sprintf("%s/docker/registry/v2/blobs/%s/%s/%s/data", r.dir, dgstAlg, dgstHex[:2], dgstHex)
***REMOVED***

// ReadBlobContents read the file corresponding to the specified digest
func (r *V2) ReadBlobContents(t testingT, blobDigest digest.Digest) []byte ***REMOVED***
	// Load the target manifest blob.
	manifestBlob, err := ioutil.ReadFile(r.getBlobFilename(blobDigest))
	if err != nil ***REMOVED***
		t.Fatalf("unable to read blob: %s", err)
	***REMOVED***

	return manifestBlob
***REMOVED***

// WriteBlobContents write the file corresponding to the specified digest with the given content
func (r *V2) WriteBlobContents(t testingT, blobDigest digest.Digest, data []byte) ***REMOVED***
	if err := ioutil.WriteFile(r.getBlobFilename(blobDigest), data, os.FileMode(0644)); err != nil ***REMOVED***
		t.Fatalf("unable to write malicious data blob: %s", err)
	***REMOVED***
***REMOVED***

// TempMoveBlobData moves the existing data file aside, so that we can replace it with a
// malicious blob of data for example.
func (r *V2) TempMoveBlobData(t testingT, blobDigest digest.Digest) (undo func()) ***REMOVED***
	tempFile, err := ioutil.TempFile("", "registry-temp-blob-")
	if err != nil ***REMOVED***
		t.Fatalf("unable to get temporary blob file: %s", err)
	***REMOVED***
	tempFile.Close()

	blobFilename := r.getBlobFilename(blobDigest)

	// Move the existing data file aside, so that we can replace it with a
	// another blob of data.
	if err := os.Rename(blobFilename, tempFile.Name()); err != nil ***REMOVED***
		os.Remove(tempFile.Name())
		t.Fatalf("unable to move data blob: %s", err)
	***REMOVED***

	return func() ***REMOVED***
		os.Rename(tempFile.Name(), blobFilename)
		os.Remove(tempFile.Name())
	***REMOVED***
***REMOVED***

// Username returns the configured user name of the server
func (r *V2) Username() string ***REMOVED***
	return r.username
***REMOVED***

// Password returns the configured password of the server
func (r *V2) Password() string ***REMOVED***
	return r.password
***REMOVED***

// Email returns the configured email of the server
func (r *V2) Email() string ***REMOVED***
	return r.email
***REMOVED***

// Path returns the path where the registry write data
func (r *V2) Path() string ***REMOVED***
	return filepath.Join(r.dir, "docker", "registry", "v2")
***REMOVED***
