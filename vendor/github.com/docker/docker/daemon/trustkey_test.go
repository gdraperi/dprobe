package daemon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/internal/testutil"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LoadOrCreateTrustKey
func TestLoadOrCreateTrustKeyInvalidKeyFile(t *testing.T) ***REMOVED***
	tmpKeyFolderPath, err := ioutil.TempDir("", "api-trustkey-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpKeyFolderPath)

	tmpKeyFile, err := ioutil.TempFile(tmpKeyFolderPath, "keyfile")
	require.NoError(t, err)

	_, err = loadOrCreateTrustKey(tmpKeyFile.Name())
	testutil.ErrorContains(t, err, "Error loading key file")
***REMOVED***

func TestLoadOrCreateTrustKeyCreateKeyWhenFileDoesNotExist(t *testing.T) ***REMOVED***
	tmpKeyFolderPath := fs.NewDir(t, "api-trustkey-test")
	defer tmpKeyFolderPath.Remove()

	// Without the need to create the folder hierarchy
	tmpKeyFile := tmpKeyFolderPath.Join("keyfile")

	key, err := loadOrCreateTrustKey(tmpKeyFile)
	require.NoError(t, err)
	assert.NotNil(t, key)

	_, err = os.Stat(tmpKeyFile)
	require.NoError(t, err, "key file doesn't exist")
***REMOVED***

func TestLoadOrCreateTrustKeyCreateKeyWhenDirectoryDoesNotExist(t *testing.T) ***REMOVED***
	tmpKeyFolderPath := fs.NewDir(t, "api-trustkey-test")
	defer tmpKeyFolderPath.Remove()
	tmpKeyFile := tmpKeyFolderPath.Join("folder/hierarchy/keyfile")

	key, err := loadOrCreateTrustKey(tmpKeyFile)
	require.NoError(t, err)
	assert.NotNil(t, key)

	_, err = os.Stat(tmpKeyFile)
	require.NoError(t, err, "key file doesn't exist")
***REMOVED***

func TestLoadOrCreateTrustKeyCreateKeyNoPath(t *testing.T) ***REMOVED***
	defer os.Remove("keyfile")
	key, err := loadOrCreateTrustKey("keyfile")
	require.NoError(t, err)
	assert.NotNil(t, key)

	_, err = os.Stat("keyfile")
	require.NoError(t, err, "key file doesn't exist")
***REMOVED***

func TestLoadOrCreateTrustKeyLoadValidKey(t *testing.T) ***REMOVED***
	tmpKeyFile := filepath.Join("testdata", "keyfile")
	key, err := loadOrCreateTrustKey(tmpKeyFile)
	require.NoError(t, err)
	expected := "AWX2:I27X:WQFX:IOMK:CNAK:O7PW:VYNB:ZLKC:CVAE:YJP2:SI4A:XXAY"
	assert.Contains(t, key.String(), expected)
***REMOVED***
