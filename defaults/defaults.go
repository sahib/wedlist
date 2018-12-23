package defaults

import (
	"os"

	e "github.com/pkg/errors"
	"github.com/sahib/config"
)

// CurrentVersion is the current version of brig's config
const CurrentVersion = 0

// Defaults is the default validation for brig
var Defaults = DefaultsV0

// OpenMigratedConfig takes the config.yml at path and loads it.
// If required, it also migrates the config structure to the newest
// version - brig can always rely on the latest config keys to be present.
func OpenMigratedConfig(path string) (*config.Config, error) {
	fd, err := os.Open(path) // #nosec
	if err != nil {
		return nil, e.Wrapf(err, "failed to open config path %s", path)
	}

	defer fd.Close()

	// NOTE: Add here any migrations with mgr.Add if needed.
	mgr := config.NewMigrater(CurrentVersion, config.StrictnessPanic)
	mgr.Add(0, nil, DefaultsV0)

	cfg, err := mgr.Migrate(config.NewYamlDecoder(fd))
	if err != nil {
		return nil, e.Wrap(err, "failed to migrate or open")
	}

	return cfg, nil
}
