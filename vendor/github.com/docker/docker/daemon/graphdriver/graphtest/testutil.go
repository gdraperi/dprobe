package graphtest

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"sort"

	"github.com/containerd/continuity/driver"
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stringid"
)

func randomContent(size int, seed int64) []byte ***REMOVED***
	s := rand.NewSource(seed)
	content := make([]byte, size)

	for i := 0; i < len(content); i += 7 ***REMOVED***
		val := s.Int63()
		for j := 0; i+j < len(content) && j < 7; j++ ***REMOVED***
			content[i+j] = byte(val)
			val >>= 8
		***REMOVED***
	***REMOVED***

	return content
***REMOVED***

func addFiles(drv graphdriver.Driver, layer string, seed int64) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	if err := driver.WriteFile(root, root.Join(root.Path(), "file-a"), randomContent(64, seed), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := root.MkdirAll(root.Join(root.Path(), "dir-b"), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := driver.WriteFile(root, root.Join(root.Path(), "dir-b", "file-b"), randomContent(128, seed+1), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***

	return driver.WriteFile(root, root.Join(root.Path(), "file-c"), randomContent(128*128, seed+2), 0755)
***REMOVED***

func checkFile(drv graphdriver.Driver, layer, filename string, content []byte) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	fileContent, err := driver.ReadFile(root, root.Join(root.Path(), filename))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !bytes.Equal(fileContent, content) ***REMOVED***
		return fmt.Errorf("mismatched file content %v, expecting %v", fileContent, content)
	***REMOVED***

	return nil
***REMOVED***

func addFile(drv graphdriver.Driver, layer, filename string, content []byte) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	return driver.WriteFile(root, root.Join(root.Path(), filename), content, 0755)
***REMOVED***

func addDirectory(drv graphdriver.Driver, layer, dir string) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	return root.MkdirAll(root.Join(root.Path(), dir), 0755)
***REMOVED***

func removeAll(drv graphdriver.Driver, layer string, names ...string) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	for _, filename := range names ***REMOVED***
		if err := root.RemoveAll(root.Join(root.Path(), filename)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func checkFileRemoved(drv graphdriver.Driver, layer, filename string) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	if _, err := root.Stat(root.Join(root.Path(), filename)); err == nil ***REMOVED***
		return fmt.Errorf("file still exists: %s", root.Join(root.Path(), filename))
	***REMOVED*** else if !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func addManyFiles(drv graphdriver.Driver, layer string, count int, seed int64) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	for i := 0; i < count; i += 100 ***REMOVED***
		dir := root.Join(root.Path(), fmt.Sprintf("directory-%d", i))
		if err := root.MkdirAll(dir, 0755); err != nil ***REMOVED***
			return err
		***REMOVED***
		for j := 0; i+j < count && j < 100; j++ ***REMOVED***
			file := root.Join(dir, fmt.Sprintf("file-%d", i+j))
			if err := driver.WriteFile(root, file, randomContent(64, seed+int64(i+j)), 0755); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func changeManyFiles(drv graphdriver.Driver, layer string, count int, seed int64) ([]archive.Change, error) ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer drv.Put(layer)

	changes := []archive.Change***REMOVED******REMOVED***
	for i := 0; i < count; i += 100 ***REMOVED***
		archiveRoot := fmt.Sprintf("/directory-%d", i)
		if err := root.MkdirAll(root.Join(root.Path(), archiveRoot), 0755); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for j := 0; i+j < count && j < 100; j++ ***REMOVED***
			if j == 0 ***REMOVED***
				changes = append(changes, archive.Change***REMOVED***
					Path: archiveRoot,
					Kind: archive.ChangeModify,
				***REMOVED***)
			***REMOVED***
			var change archive.Change
			switch j % 3 ***REMOVED***
			// Update file
			case 0:
				change.Path = root.Join(archiveRoot, fmt.Sprintf("file-%d", i+j))
				change.Kind = archive.ChangeModify
				if err := driver.WriteFile(root, root.Join(root.Path(), change.Path), randomContent(64, seed+int64(i+j)), 0755); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			// Add file
			case 1:
				change.Path = root.Join(archiveRoot, fmt.Sprintf("file-%d-%d", seed, i+j))
				change.Kind = archive.ChangeAdd
				if err := driver.WriteFile(root, root.Join(root.Path(), change.Path), randomContent(64, seed+int64(i+j)), 0755); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			// Remove file
			case 2:
				change.Path = root.Join(archiveRoot, fmt.Sprintf("file-%d", i+j))
				change.Kind = archive.ChangeDelete
				if err := root.Remove(root.Join(root.Path(), change.Path)); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			***REMOVED***
			changes = append(changes, change)
		***REMOVED***
	***REMOVED***

	return changes, nil
***REMOVED***

func checkManyFiles(drv graphdriver.Driver, layer string, count int, seed int64) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	for i := 0; i < count; i += 100 ***REMOVED***
		dir := root.Join(root.Path(), fmt.Sprintf("directory-%d", i))
		for j := 0; i+j < count && j < 100; j++ ***REMOVED***
			file := root.Join(dir, fmt.Sprintf("file-%d", i+j))
			fileContent, err := driver.ReadFile(root, file)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			content := randomContent(64, seed+int64(i+j))

			if !bytes.Equal(fileContent, content) ***REMOVED***
				return fmt.Errorf("mismatched file content %v, expecting %v", fileContent, content)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type changeList []archive.Change

func (c changeList) Less(i, j int) bool ***REMOVED***
	if c[i].Path == c[j].Path ***REMOVED***
		return c[i].Kind < c[j].Kind
	***REMOVED***
	return c[i].Path < c[j].Path
***REMOVED***
func (c changeList) Len() int      ***REMOVED*** return len(c) ***REMOVED***
func (c changeList) Swap(i, j int) ***REMOVED*** c[j], c[i] = c[i], c[j] ***REMOVED***

func checkChanges(expected, actual []archive.Change) error ***REMOVED***
	if len(expected) != len(actual) ***REMOVED***
		return fmt.Errorf("unexpected number of changes, expected %d, got %d", len(expected), len(actual))
	***REMOVED***
	sort.Sort(changeList(expected))
	sort.Sort(changeList(actual))

	for i := range expected ***REMOVED***
		if expected[i] != actual[i] ***REMOVED***
			return fmt.Errorf("unexpected change, expecting %v, got %v", expected[i], actual[i])
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func addLayerFiles(drv graphdriver.Driver, layer, parent string, i int) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	if err := driver.WriteFile(root, root.Join(root.Path(), "top-id"), []byte(layer), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	layerDir := root.Join(root.Path(), fmt.Sprintf("layer-%d", i))
	if err := root.MkdirAll(layerDir, 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := driver.WriteFile(root, root.Join(layerDir, "layer-id"), []byte(layer), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	return driver.WriteFile(root, root.Join(layerDir, "parent-id"), []byte(parent), 0755)
***REMOVED***

func addManyLayers(drv graphdriver.Driver, baseLayer string, count int) (string, error) ***REMOVED***
	lastLayer := baseLayer
	for i := 1; i <= count; i++ ***REMOVED***
		nextLayer := stringid.GenerateRandomID()
		if err := drv.Create(nextLayer, lastLayer, nil); err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if err := addLayerFiles(drv, nextLayer, lastLayer, i); err != nil ***REMOVED***
			return "", err
		***REMOVED***

		lastLayer = nextLayer

	***REMOVED***
	return lastLayer, nil
***REMOVED***

func checkManyLayers(drv graphdriver.Driver, layer string, count int) error ***REMOVED***
	root, err := drv.Get(layer, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer drv.Put(layer)

	layerIDBytes, err := driver.ReadFile(root, root.Join(root.Path(), "top-id"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !bytes.Equal(layerIDBytes, []byte(layer)) ***REMOVED***
		return fmt.Errorf("mismatched file content %v, expecting %v", layerIDBytes, []byte(layer))
	***REMOVED***

	for i := count; i > 0; i-- ***REMOVED***
		layerDir := root.Join(root.Path(), fmt.Sprintf("layer-%d", i))

		thisLayerIDBytes, err := driver.ReadFile(root, root.Join(layerDir, "layer-id"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !bytes.Equal(thisLayerIDBytes, layerIDBytes) ***REMOVED***
			return fmt.Errorf("mismatched file content %v, expecting %v", thisLayerIDBytes, layerIDBytes)
		***REMOVED***
		layerIDBytes, err = driver.ReadFile(root, root.Join(layerDir, "parent-id"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// readDir reads a directory just like driver.ReadDir()
// then hides specific files (currently "lost+found")
// so the tests don't "see" it
func readDir(r driver.Driver, dir string) ([]os.FileInfo, error) ***REMOVED***
	a, err := driver.ReadDir(r, dir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	b := a[:0]
	for _, x := range a ***REMOVED***
		if x.Name() != "lost+found" ***REMOVED*** // ext4 always have this dir
			b = append(b, x)
		***REMOVED***
	***REMOVED***

	return b, nil
***REMOVED***
