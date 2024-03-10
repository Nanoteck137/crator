package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/kr/pretty"
	"github.com/nanoteck137/crator/cmd"
	"github.com/nanoteck137/crator/template"
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

func readConfig() (*AppConfig, error) {
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

func main() {
	cmd.Execute()

	return

	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(config)

	err = os.MkdirAll(config.Templates, 0755)
	if err != nil {
		log.Fatal(err)
	}

	var paths []string

	err = filepath.WalkDir(filepath.Clean(config.Templates), func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if name == "crator.json" {
			paths = append(paths, p)
			return filepath.SkipDir
		}

		return nil
	})

	type Template struct {
		p string
		config template.Config
	}

	var templates []Template

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}

		var templateConfig template.Config
		err = json.Unmarshal(data,&templateConfig)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Template: ", templateConfig.Name)
		templates = append(templates, Template{
			p:      p,
			config: templateConfig,
		})
	}

	templ := templates[0]
	p := filepath.Dir(templ.p)
	template.Execute(&templ.config, p, "./work/test")
}
