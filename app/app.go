package app

import (
	"encoding/json"
	"errors"
	"os"
)

type AppConfig struct {
	Templates string `json:"templates"`
}

func getConfigPath() string {
	return os.ExpandEnv("$HOME/.config/crator/config.json")
}

func getDefaultTemplateDir() string {
	dataHome, exists := os.LookupEnv("XDG_DATA_HOME")
	if exists {
		return dataHome + "/crator/templates"
	}

	return os.ExpandEnv("$HOME/.local/share/crator/templates")
}

func ReadConfig() (*AppConfig, error) {
	p := getConfigPath()

	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &AppConfig{
				Templates: getDefaultTemplateDir(),
			}, nil
		}

		return nil, err
	}

	var config AppConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Templates == "" {
		config.Templates = getDefaultTemplateDir()
	}

	return &config, nil
}
