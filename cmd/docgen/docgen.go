package main

import (
	"fmt"
	"path/filepath"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
	"os"
)

func main() {
	fset := token.NewFileSet()

	dirs := []string{"./"}
	filepath.Walk("golibs/", func (path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		dirs = append(dirs, "./" + path)
		return nil
	})

	pkgs := make(map[string]*ast.Package)
	for _, path := range dirs {
		d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range d {
			pkgs[k] = v
		}
	}

	prefix := map[string]string{
		"main": "hsh",
		"fs": "f",
		"commander": "c",
		"bait": "b",
	}

	for l, f := range pkgs {
		fmt.Println("------", l)
		p := doc.New(f, "./", doc.AllDecls)
		for _, t := range p.Funcs {
			if !strings.HasPrefix(t.Name, prefix[l]) || t.Name == "Loader" { continue }
			parts := strings.Split(t.Doc, "\n")
			funcsig := parts[0]
			doc := parts[1]

			fmt.Println(funcsig, ">", doc)
		}
		for _, t := range p.Types {
			for _, m := range t.Methods {
				if !strings.HasPrefix(m.Name, prefix[l]) || m.Name == "Loader" { continue }
				parts := strings.Split(m.Doc, "\n")
				funcsig := parts[0]
				doc := parts[1]

				fmt.Println(funcsig, ">", doc)
			}
		}
	}
}
