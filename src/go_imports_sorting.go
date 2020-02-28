package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

var (
	list        = flag.Bool("l", false, "write results to stdout")
	write       = flag.Bool("w", true, "write result to (source) file instead of stdout")
	localPrefix = flag.String("local", "", "put imports beginning with this string after 3rd-party packages; comma-separated list")
	verbose     bool // verbose logging
)

// main is the entry point of the program
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := goImportsSortMain()
	if err != nil {
		log.Fatalln(err)
	}
}

// goImportsSortMain checks passed flags and starts processing files
func goImportsSortMain() error {
	flag.Usage = func () {
		_, _ = fmt.Fprintf(os.Stderr, "usage: goimportssort [flags] [path ...]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	paths := parseFlags()

	if verbose {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	if *localPrefix == "" {
		log.Println("No prefix found, using module name")

		moduleName, _ := getModuleName()
		if moduleName != "" {
			localPrefix = &moduleName
		} else {
			log.Println("Module name not found. Skipping localPrefix")
		}
	}

	if len(paths) == 0 {
		return errors.New("please enter a path to fix")
	}

	for _, path := range paths {
		switch dir, statErr := os.Stat(path); {
		case statErr != nil:
			return statErr
		case dir.IsDir():
			return walkDir(path)
		default:
			_, err := processFile(path, nil, os.Stdout)
			return err
		}
	}

	return nil
}

// parseFlags parses command line flags and returns the paths to process.
// It's a var so that custom implementations can replace it in other files.
var parseFlags = func() []string {
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.Parse()

	return flag.Args()
}

// isGoFile checks if the file is a go file & not a directory
func isGoFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

// walkDir walks through a path, processing all go files recursively in a directory
func walkDir(path string) error {
	return filepath.Walk(path, func (path string, f os.FileInfo, err error) error {
		if err == nil && isGoFile(f) {
			_, err = processFile(path, nil, os.Stdout)
		}
		return err
	})
}

// processFile reads a file and processes the content, then checks if they're equal.
func processFile(filename string, in io.Reader, out io.Writer) ([]byte, error) {
	log.Printf("Processing %v\n", filename)

	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	res, err := process(src)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(src, res) {
		// formatting has changed
		if *list {
			_, _ = fmt.Fprintln(out, string(res))
		}
		if *write {
			err = ioutil.WriteFile(filename, res, 0)
			if err != nil {
				return nil, err
			}
		}
		if !*list && !*write {
			return res, nil
		}
	} else {
		log.Println("File has not been changed.")
	}

	return nil, err
}

// process processes the source of a file, categorising the imports
func process(src []byte) (formatted []byte, err error) {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	convertedImports, err := convertImportsToSlice(node)
	sortedImports := sortImports(convertedImports)
	convertedToGo := convertImportsToGo(sortedImports)

	output := replaceImports(src, convertedToGo)

	return output, err
}

// replaceImports replaces existing imports and handles multiple import statements
func replaceImports(src, newImports []byte) []byte {
	// remove single imports
	output := regexp.MustCompile(`(?U)import ".*"`).ReplaceAll(src, []byte{})

	// replace first long import with the sorted imports
	var findMultipleImports = regexp.MustCompile(`(?sU)import \(.*\)`)
	output = bytes.Replace(output, findMultipleImports.Find(output), newImports, 1)

	// remove any additional long import blocks, skip the first one
	allMatches := findMultipleImports.FindAll(output, -1)
	if len(allMatches) > 1 {
		for i := 0; i < (len(allMatches) - 1); i++ {
			output = bytes.Replace(output, allMatches[i+1], []byte{}, 1)
		}
	}

	// clear additional whitespace
	output = regexp.MustCompile(`\n{2,}`).ReplaceAll(output, []byte("\n\n")) // TODO: do not replace all whitespace

	return output
}

// sortImports sorts multiple imports by import name & prefix
func sortImports(imports [][]ImpModel) [][]ImpModel {
	for x := 0; x < len(imports); x++ {
		sort.Slice(imports[x], func(i, j int) bool {
			if imports[x][i].path != imports[x][j].path {
				return imports[x][i].path < imports[x][j].path
			}

			return imports[x][i].localReference < imports[x][j].localReference
		})
	}

	return imports
}

// convertImportsToGo generates output for correct categorised import statements
func convertImportsToGo(imports [][]ImpModel) []byte {
	output := "import ("

	for i := 0; i < len(imports); i++ {
		output += "\n"
		for _, imp := range imports[i] {
			output += fmt.Sprintf("\t%v\n", imp.String())
		}
	}

	output += ")"

	return []byte(output)
}

// convertImportsToSlice parses the file with AST and gets all imports
func convertImportsToSlice(node *ast.File) ([][]ImpModel, error) {
	importCategories := make([][]ImpModel, 3)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			impName := importSpec.Path.Value
			locName := importSpec.Name

			var impModel ImpModel
			if locName != nil {
				impModel.localReference = locName.Name
			}
			impModel.path = impName

			if *localPrefix != "" && strings.Count(impName, *localPrefix) > 0 { // TODO: Support multiple local packages
				importCategories[2] = append(importCategories[2], impModel)
			} else if strings.Count(impName, "/") <= 1 {
				importCategories[0] = append(importCategories[0], impModel)
			} else {
				importCategories[1] = append(importCategories[1], impModel)
			}
		}
	}

	return importCategories, nil
}

// getModuleName parses the GOMOD name
func getModuleName() (string, error) {
	gomodCmd := exec.Command("go", "env", "GOMOD") // TODO: Check if there's a better way to get GOMOD
	gomod, err := gomodCmd.Output()
	if err != nil {
		log.Println("Could not run: go env GOMOD")
		return "", err
	}
	gomodStr := strings.TrimSpace(string(gomod))

	moduleCmd := exec.Command("awk", "/module/ {print $2}", gomodStr) // TODO: Check if there's a better way
	module, err := moduleCmd.Output()
	if err != nil {
		return "", err
	}
	moduleStr := strings.TrimSpace(string(module))

	return moduleStr, nil
}
