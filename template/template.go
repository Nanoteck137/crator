package template

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nanoteck137/crator/app"
)

type Var struct {
	Name    string `json:"name"`
	Prompt  string `json:"prompt"`
	Default string `json:"default"`
}

type Config struct {
	Name string `json:"name"`
	Vars []Var  `json:"vars"`
}

func getVars(config *Config) map[string]string {
	res := make(map[string]string)

	for _, v := range config.Vars {
		reader := bufio.NewReader(os.Stdin)
		if v.Prompt != "" {
			fmt.Print(v.Prompt + ": ")
		} else {
			fmt.Printf("Enter '%s': ", v.Name)
		}
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		res[v.Name] = text
	}

	return res
}

func Execute(config *Config, src string, dst string) error {
	vars := getVars(config)

	var dirs []string
	var files []string

	err := filepath.WalkDir(filepath.Clean(src), func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() && p != src {
			fmt.Println(d.Info())
			dirs = append(dirs, strings.TrimPrefix(p, src+"/"))
		}

		if !d.IsDir() && d.Name() != "crator.json" {
			files = append(files, p)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirs {
		p := path.Join(dst, dir)
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, file := range files {
		fmt.Printf("file: %v\n", file)
		fmt.Printf("src: %v\n", src)
		p := path.Join(dst, strings.TrimPrefix(file, src+"/"))
		fmt.Printf("p: %v\n", p)

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(p)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		templ, err := template.New(file).Parse(string(data))
		if err != nil {
			log.Fatal(err)
		}

		templ.Execute(f, vars)
	}

	return nil
}

type Template struct {
	Base   string
	Config Config
}

func (t *Template) Execute(dst string) error {
	return Execute(&t.Config, t.Base, dst)
}

func GetAvailable(config *app.AppConfig) ([]Template, error) {
	var paths []string
	err := filepath.WalkDir(filepath.Clean(config.Templates), func(p string, d fs.DirEntry, err error) error {
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

	if err != nil {
		return nil, err
	}

	var templates []Template

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}

		var templateConfig Config
		err = json.Unmarshal(data, &templateConfig)
		if err != nil {
			log.Fatal(err)
		}

		templates = append(templates, Template{
			Base:   filepath.Dir(p),
			Config: templateConfig,
		})
	}

	return templates, nil
}
