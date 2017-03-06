package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//TEST COMMENT - DO NOT REMOVE

func findFiles(root string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	return files, err
}

func parseDir(path string) map[string]*ast.Package {
	//from https://golang.org/src/go/parser/interface.go:134
	fset := token.NewFileSet()
	mode := parser.AllErrors | parser.ParseComments

	list, err := findFiles(path)
	if err != nil {
		panic(err.Error())
	}

	pkgs := make(map[string]*ast.Package)
	for _, filename := range list {
		if strings.HasSuffix(filename, ".go") {
			if src, err := parser.ParseFile(fset, filename, nil, mode); err == nil {
				name := src.Name.Name
				pkg, found := pkgs[name]
				if !found {
					pkg = &ast.Package{
						Name:  name,
						Files: make(map[string]*ast.File),
					}
					pkgs[name] = pkg
				}
				pkg.Files[filename] = src
			}
		}
	}

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(os.Stderr, "%v packages:\n", len(pkgs))
	for k, _ := range pkgs {
		fmt.Fprintf(os.Stderr, "\t%v\n", k)
	}

	return pkgs
}

func readTemplate(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read template: %v\n", err.Error())
		os.Exit(1)
	}
	return string(data)
}

func render(text string, packages map[string]*ast.Package) string {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"ToTitle": strings.ToTitle,
	}
	tmpl, err := template.New("render").Funcs(funcMap).Parse(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while parsing template: %v\n", err.Error())
		os.Exit(1)
	}
	data := struct {
		Packages map[string]*ast.Package
	}{
		packages,
	}
	var txt bytes.Buffer
	err = tmpl.Execute(&txt, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while rendering template: %v\n", err.Error())
		os.Exit(1)
	}
	return txt.String()
}

func main() {
	var (
		sourceDir    *string = flag.String("d", "", "Directory to parse files from")
		templateFile *string = flag.String("t", "", "Template to render")
		outFile      *string = flag.String("o", "", "Output filename")
	)
	flag.Parse()
	if *sourceDir == "" {
		fmt.Fprintf(os.Stderr, "Please specify a source directory\n")
		os.Exit(1)
	}
	if *templateFile == "" {
		fmt.Fprintf(os.Stderr, "Please specify a template file\n")
		os.Exit(1)
	}
	if *outFile == "" {
		fmt.Fprintf(os.Stderr, "Please specify an output file\n")
		os.Exit(1)
	}
	packages := parseDir(*sourceDir)
	generated := render(readTemplate(*templateFile), packages)
	err := ioutil.WriteFile(*outFile, []byte(generated), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err.Error())
		os.Exit(1)
	}
}
