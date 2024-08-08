package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	name := ""
	flag.StringVar(&name, "n", name, "model name")
	flag.Parse()
	if name == "" {
		log.Println("model name is required")
		return
	}
	templateDir := "template"
	dstDir := "service"
	dstDir = filepath.Join(dstDir, strings.ToLower(name))
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"ModelName":   strings.Title(name),
		"PackageName": strings.ToLower(name),
	}
	err = generateFiles(templateDir, dstDir, []string{"handler", "model", "routes"}, data)
	if err != nil {
		panic(err)
	}
}

func generateFiles(tplDir, dstDir string, files []string, data map[string]interface{}) error {
	for _, file := range files {
		t, err := template.ParseFiles(filepath.Join(tplDir, file+".go.tpl"))
		if err != nil {
			return err
		}
		err = func() error {
			dstFile := filepath.Join(dstDir, file+".go")
			f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			err = t.Execute(f, data)
			if err != nil {
				return err
			}
			log.Printf("Generated %s\n", dstFile)
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
