package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kr/pretty"
)

type Var struct {
	Name    string `json:"name"`
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
		fmt.Printf("Enter '%s': ", v.Name)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		res[v.Name] = text
	}

	return res
}

func main() {
	config := Config{
		Name: "test",
		Vars: []Var{
			{
				Name:    "projectName",
				Default: "",
			},
		},
	}

	vars := getVars(&config)

	fmt.Println("Hello World")

	dir := "./work/test"
	dir = filepath.Clean(dir)

	var dirs []string
	var files []string

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		fmt.Printf("path: %v\n", p)

		if d.IsDir() && p != dir {
			fmt.Println(d.Info())
			dirs = append(dirs, strings.TrimPrefix(p, dir+"/"))
		}

		if !d.IsDir() {
			files = append(files, p)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(files)

	dst := "./work/dest"
	for _, dir := range dirs {
		p := path.Join(dst, dir)
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, file := range files {
		p := path.Join(dst, strings.TrimPrefix(file, dir+"/"))
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
}
