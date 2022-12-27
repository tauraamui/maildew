package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/tauraamui/maildew/internal/configdef"
	"github.com/tauraamui/xerror"
)

const (
	vendorName     = "tacusci"
	appName        = "maildew"
	configFileName = "config.json"
)

func load() (configdef.Values, error) {
	var values configdef.Values

	configPath, err := resolveConfigFilePath()
	if err != nil {
		return configdef.Values{}, err
	}

	// log.Info("Resolved config file location: %s", configPath)
	file, err := readConfigFile(configPath)
	if err != nil {
		return configdef.Values{}, err
	}

	if err := unmarshal(file, &values); err != nil {
		return configdef.Values{}, err
	}

	if err = values.RunValidate(); err != nil {
		return configdef.Values{}, err
	}

	return values, nil
}

var readConfigFile = func(path string) ([]byte, error) {
	return afero.ReadFile(fs, path)
}

func unmarshal(content []byte, values *configdef.Values) error {
	err := json.Unmarshal(content, values)
	if err != nil {
		return errors.Errorf("parsing configuration error: %v", err)
	}
	return nil
}

func resolveConfigFilePath() (string, error) {
	configDir, err := resolveAppDirPath()
	if err != nil {
		return "", xerror.Errorf("unable to resolve app config dir location: %w", err)
	}

	return filepath.Join(configDir, configFileName), nil
}

func resolveAppDirPath() (string, error) {
	configPath := os.Getenv("DRAGON_DAEMON_CONFIG")
	if len(configPath) > 0 {
		return configPath, nil
	}

	configParentDir, err := userConfigDir()
	if err != nil {
		return "", xerror.Errorf("unable to resolve %s location: %w", configFileName, err)
	}

	return filepath.Join(
		configParentDir,
		vendorName,
		appName), nil
}

var userConfigDir = func() (string, error) {
	return os.UserConfigDir()
}
