// Copyright 2017 Florin Pățan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deep

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// goParser uses the go/parser facility in order to determine the imports from
// a certain package, recursively
type goParser struct {
	log Logger
}

func (*goParser) canUse(pwd, currentPkg string) bool {
	return true
}

func (g *goParser) packages(pwd, currentPkg string, ignoreTestFiles bool) ([]Package, error) {
	var files []*ast.File
	err := filepath.Walk(pwd, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			g.log("Error reading package files: %v\n", err)
			return nil
		}

		if f.IsDir() {
			return nil
		}

		fileName := f.Name()

		if strings.Contains(fileName, "/testdata/") {
			return nil
		}

		if !strings.HasSuffix(f.Name(), ".go") {
			return nil
		}

		if ignoreTestFiles && strings.HasSuffix(f.Name(), "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		g.log("Error while loading package files: %v\n", err)
	}

	imports := map[string]Package{}
	for _, f := range files {
		for _, s := range f.Imports {
			importPath := s.Path.Value[1:len(s.Path.Value)-1]
			if _, ok := stdlibPackages[importPath]; ok {
				continue
			}

			if strings.HasPrefix(importPath, currentPkg) {
				continue
			}

			if importPath == "" {
				continue
			}
			// TODO: do more sanitization here

			if _, ok := imports[importPath]; ok {
				continue
			}

			imports[importPath] = Package{
				Name:    importPath,
				Version: "HEAD",
			}
		}
	}

	var result []Package
	for _, pkg := range imports {
		if !pkg.isRootPackage(currentPkg) {
			continue
		}

		result = append(result, pkg)
	}

	return result, nil
}

func newGoParser(logger Logger) *goParser {
	return &goParser{
		log: logger,
	}
}
