package app

import (
	"errors"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
)

const TemplateConfigName = "crator.toml"
const AppConfigName = "config.toml"

type Config struct {
	Templates string `toml:"templates"`
}

func getConfigPath() string {
	return path.Join(os.ExpandEnv("$HOME/.config/crator"), AppConfigName)
}

func getDefaultTemplateDir() string {
	dataHome, exists := os.LookupEnv("XDG_DATA_HOME")
	if exists {
		return dataHome + "/crator/templates"
	}

	return os.ExpandEnv("$HOME/.local/share/crator/templates")
}

func ReadConfig() (*Config, error) {
	p := getConfigPath()

	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{
				Templates: getDefaultTemplateDir(),
			}, nil
		}

		return nil, err
	}

	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Templates == "" {
		config.Templates = getDefaultTemplateDir()
	}

	return &config, nil
}
