package devmapper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type directLVMConfig struct ***REMOVED***
	Device              string
	ThinpPercent        uint64
	ThinpMetaPercent    uint64
	AutoExtendPercent   uint64
	AutoExtendThreshold uint64
***REMOVED***

var (
	errThinpPercentMissing = errors.New("must set both `dm.thinp_percent` and `dm.thinp_metapercent` if either is specified")
	errThinpPercentTooBig  = errors.New("combined `dm.thinp_percent` and `dm.thinp_metapercent` must not be greater than 100")
	errMissingSetupDevice  = errors.New("must provide device path in `dm.setup_device` in order to configure direct-lvm")
)

func validateLVMConfig(cfg directLVMConfig) error ***REMOVED***
	if reflect.DeepEqual(cfg, directLVMConfig***REMOVED******REMOVED***) ***REMOVED***
		return nil
	***REMOVED***
	if cfg.Device == "" ***REMOVED***
		return errMissingSetupDevice
	***REMOVED***
	if (cfg.ThinpPercent > 0 && cfg.ThinpMetaPercent == 0) || cfg.ThinpMetaPercent > 0 && cfg.ThinpPercent == 0 ***REMOVED***
		return errThinpPercentMissing
	***REMOVED***

	if cfg.ThinpPercent+cfg.ThinpMetaPercent > 100 ***REMOVED***
		return errThinpPercentTooBig
	***REMOVED***
	return nil
***REMOVED***

func checkDevAvailable(dev string) error ***REMOVED***
	lvmScan, err := exec.LookPath("lvmdiskscan")
	if err != nil ***REMOVED***
		logrus.Debug("could not find lvmdiskscan")
		return nil
	***REMOVED***

	out, err := exec.Command(lvmScan).CombinedOutput()
	if err != nil ***REMOVED***
		logrus.WithError(err).Error(string(out))
		return nil
	***REMOVED***

	if !bytes.Contains(out, []byte(dev)) ***REMOVED***
		return errors.Errorf("%s is not available for use with devicemapper", dev)
	***REMOVED***
	return nil
***REMOVED***

func checkDevInVG(dev string) error ***REMOVED***
	pvDisplay, err := exec.LookPath("pvdisplay")
	if err != nil ***REMOVED***
		logrus.Debug("could not find pvdisplay")
		return nil
	***REMOVED***

	out, err := exec.Command(pvDisplay, dev).CombinedOutput()
	if err != nil ***REMOVED***
		logrus.WithError(err).Error(string(out))
		return nil
	***REMOVED***

	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(out)))
	for scanner.Scan() ***REMOVED***
		fields := strings.SplitAfter(strings.TrimSpace(scanner.Text()), "VG Name")
		if len(fields) > 1 ***REMOVED***
			// got "VG Name" line"
			vg := strings.TrimSpace(fields[1])
			if len(vg) > 0 ***REMOVED***
				return errors.Errorf("%s is already part of a volume group %q: must remove this device from any volume group or provide a different device", dev, vg)
			***REMOVED***
			logrus.Error(fields)
			break
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func checkDevHasFS(dev string) error ***REMOVED***
	blkid, err := exec.LookPath("blkid")
	if err != nil ***REMOVED***
		logrus.Debug("could not find blkid")
		return nil
	***REMOVED***

	out, err := exec.Command(blkid, dev).CombinedOutput()
	if err != nil ***REMOVED***
		logrus.WithError(err).Error(string(out))
		return nil
	***REMOVED***

	fields := bytes.Fields(out)
	for _, f := range fields ***REMOVED***
		kv := bytes.Split(f, []byte***REMOVED***'='***REMOVED***)
		if bytes.Equal(kv[0], []byte("TYPE")) ***REMOVED***
			v := bytes.Trim(kv[1], "\"")
			if len(v) > 0 ***REMOVED***
				return errors.Errorf("%s has a filesystem already, use dm.directlvm_device_force=true if you want to wipe the device", dev)
			***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func verifyBlockDevice(dev string, force bool) error ***REMOVED***
	if err := checkDevAvailable(dev); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := checkDevInVG(dev); err != nil ***REMOVED***
		return err
	***REMOVED***
	if force ***REMOVED***
		return nil
	***REMOVED***
	return checkDevHasFS(dev)
***REMOVED***

func readLVMConfig(root string) (directLVMConfig, error) ***REMOVED***
	var cfg directLVMConfig

	p := filepath.Join(root, "setup-config.json")
	b, err := ioutil.ReadFile(p)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return cfg, nil
		***REMOVED***
		return cfg, errors.Wrap(err, "error reading existing setup config")
	***REMOVED***

	// check if this is just an empty file, no need to produce a json error later if so
	if len(b) == 0 ***REMOVED***
		return cfg, nil
	***REMOVED***

	err = json.Unmarshal(b, &cfg)
	return cfg, errors.Wrap(err, "error unmarshaling previous device setup config")
***REMOVED***

func writeLVMConfig(root string, cfg directLVMConfig) error ***REMOVED***
	p := filepath.Join(root, "setup-config.json")
	b, err := json.Marshal(cfg)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error marshalling direct lvm config")
	***REMOVED***
	err = ioutil.WriteFile(p, b, 0600)
	return errors.Wrap(err, "error writing direct lvm config to file")
***REMOVED***

func setupDirectLVM(cfg directLVMConfig) error ***REMOVED***
	lvmProfileDir := "/etc/lvm/profile"
	binaries := []string***REMOVED***"pvcreate", "vgcreate", "lvcreate", "lvconvert", "lvchange", "thin_check"***REMOVED***

	for _, bin := range binaries ***REMOVED***
		if _, err := exec.LookPath(bin); err != nil ***REMOVED***
			return errors.Wrap(err, "error looking up command `"+bin+"` while setting up direct lvm")
		***REMOVED***
	***REMOVED***

	err := os.MkdirAll(lvmProfileDir, 0755)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error creating lvm profile directory")
	***REMOVED***

	if cfg.AutoExtendPercent == 0 ***REMOVED***
		cfg.AutoExtendPercent = 20
	***REMOVED***

	if cfg.AutoExtendThreshold == 0 ***REMOVED***
		cfg.AutoExtendThreshold = 80
	***REMOVED***

	if cfg.ThinpPercent == 0 ***REMOVED***
		cfg.ThinpPercent = 95
	***REMOVED***
	if cfg.ThinpMetaPercent == 0 ***REMOVED***
		cfg.ThinpMetaPercent = 1
	***REMOVED***

	out, err := exec.Command("pvcreate", "-f", cfg.Device).CombinedOutput()
	if err != nil ***REMOVED***
		return errors.Wrap(err, string(out))
	***REMOVED***

	out, err = exec.Command("vgcreate", "docker", cfg.Device).CombinedOutput()
	if err != nil ***REMOVED***
		return errors.Wrap(err, string(out))
	***REMOVED***

	out, err = exec.Command("lvcreate", "--wipesignatures", "y", "-n", "thinpool", "docker", "--extents", fmt.Sprintf("%d%%VG", cfg.ThinpPercent)).CombinedOutput()
	if err != nil ***REMOVED***
		return errors.Wrap(err, string(out))
	***REMOVED***
	out, err = exec.Command("lvcreate", "--wipesignatures", "y", "-n", "thinpoolmeta", "docker", "--extents", fmt.Sprintf("%d%%VG", cfg.ThinpMetaPercent)).CombinedOutput()
	if err != nil ***REMOVED***
		return errors.Wrap(err, string(out))
	***REMOVED***

	out, err = exec.Command("lvconvert", "-y", "--zero", "n", "-c", "512K", "--thinpool", "docker/thinpool", "--poolmetadata", "docker/thinpoolmeta").CombinedOutput()
	if err != nil ***REMOVED***
		return errors.Wrap(err, string(out))
	***REMOVED***

	profile := fmt.Sprintf("activation***REMOVED***\nthin_pool_autoextend_threshold=%d\nthin_pool_autoextend_percent=%d\n***REMOVED***", cfg.AutoExtendThreshold, cfg.AutoExtendPercent)
	err = ioutil.WriteFile(lvmProfileDir+"/docker-thinpool.profile", []byte(profile), 0600)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error writing docker thinp autoextend profile")
	***REMOVED***

	out, err = exec.Command("lvchange", "--metadataprofile", "docker-thinpool", "docker/thinpool").CombinedOutput()
	return errors.Wrap(err, string(out))
***REMOVED***
