package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kr/pretty"
)

func main() {
	fmt.Println("Hello World")

	dir := "./work/test"
	dir = filepath.Clean(dir)

	var dirs []string

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		fmt.Printf("path: %v\n", p)

		if d.IsDir() && p != dir {
			fmt.Println(d.Info())
			dirs = append(dirs, strings.TrimPrefix(p, dir + "/"))
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(dirs)

	dst := "./work/dest"
	for _, dir := range dirs {
		p := path.Join(dst, dir)
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
