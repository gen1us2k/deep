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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type (
	// Logger defines the type for the function that is expected to handle logging of messages
	Logger func(msg string, v ...interface{})

	// Deep holds the different components together
	Deep struct {
		log       Logger
		providers []provider
		vcsDirs   []string
	}
)

const pathSeparatorString = string(os.PathSeparator)

func (d *Deep) listPackages(pwd, currentPkg string) []Package {
	var packages []Package
	for _, provider := range d.providers {
		if !provider.canUse(pwd, currentPkg) {
			continue
		}

		pkgs, err := provider.packages(pwd, currentPkg)
		if err != nil {
			d.log("Error while getting packages: %v", err)
			os.Exit(1)
		}

		for _, pkg := range pkgs {
			if !pkg.isThirdParty(currentPkg) {
				continue
			}
			packages = append(packages, pkg)
		}

	}

	return packages
}

func (d *Deep) vendorGithubPackage(pwd string, pkg Package) error {
	// TODO use github releases API to fetch that version directly

	//gitPath := "git@" + strings.Replace(pkg.Name, "github.com/", "github.com:", 1)
	gitPath := "https://" + pkg.Name + ".git"

	cmd := exec.Command("git", "clone", "-v", gitPath, pkg.vendoredPath(pwd))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "checkout", pkg.Version)
	cmd.Dir = pkg.vendoredPath(pwd)
	return cmd.Run()
}

func (d *Deep) tryGitVendor(pwd string, pkg Package) error {
	if strings.HasPrefix(pkg.Name, "github.com") {
		return d.vendorGithubPackage(pwd, pkg)
	}

	d.log("Could not vendor Git dependency as it's not starting with github.com")
	return errors.New("Not a github package")
}

func (d *Deep) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (d *Deep) shouldWipePath(path string) bool {
	retries := 0
	response := ""
prompt:
	retries++
	fmt.Printf("Do you want to wipe %s [Y/n] ", path)
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}

	response = strings.ToLower(response)
	if response == "y" {
		return true
	} else if response == "n" {
		return false
	} else if retries > 3 {
		return false
	}

	goto prompt
}

func (d *Deep) vendorPackages(pwd, currentPkg string, packages []Package) {
	for _, pkg := range packages {
		vendoredPath := pkg.vendoredPath(pwd)
		pathExists, err := d.pathExists(vendoredPath)
		if err != nil {
			d.log("Got error while checking path %s %v Skipping\n", vendoredPath, err)
			continue
		}

		if pathExists {
			if !d.shouldWipePath(vendoredPath) {
				d.log("Skipping existing path: %s\n", vendoredPath)
				continue
			}
			err := os.RemoveAll(vendoredPath)
			if err != nil {
				d.log("Could not wipe existing path: %s %v\n", vendoredPath, err)
				os.Exit(1)
			}
		}

		err = d.tryGitVendor(pwd, pkg)
		if err != nil {
			d.log("Got error while trying to clone repository: %s %v\n", pkg.Name, err)
			os.Exit(1)
		}
		continue
	}
}

func (d *Deep) commitHash(pwd, currentPkg string, pkg Package) string {
	vendoredPath := pkg.vendoredPath(pwd)
	for _, vcsDir := range d.vcsDirs {
		pathExists, err := d.pathExists(vendoredPath + pathSeparatorString + vcsDir)
		if err != nil {
			d.log("Got error while checking path %s %v Skipping\n", pathExists, err)
			continue
		}

		if !pathExists {
			continue
		}

		switch vcsDir {
		case ".git":
			cmd := exec.Command("git", "rev-parse", pkg.Version)
			cmd.Dir = pkg.vendoredPath(pwd)
			output, err := cmd.Output()
			if err != nil {
				d.log("Error while reading package %s version %v\n", pkg.Name, err)
				return pkg.Version
			}

			return strings.TrimRight(string(output), "\n")
		}
	}

	return pkg.Version
}

func (d *Deep) readCommitHashes(pwd, currentPkg string, packages []Package) {
	for idx, pkg := range packages {
		packages[idx].CommitHash = d.commitHash(pwd, currentPkg, pkg)
	}
}

func (d *Deep) wipeNestedVendor(pwd string, currentPkg string, packages []Package) {
	for _, pkg := range packages {
		path := pwd + "/vendor/" + pkg.Name + "/vendor"
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			d.log("Error while wiping nested vendor folders %v\n", err)
		}
	}
}

func (d *Deep) wipeTestFiles(pwd, currentPkg string, packages []Package) {
	for _, pkg := range packages {
		path := pwd + "/vendor/" + pkg.Name + "/"
		err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				d.log("Error while wiping test files: %v\n", err)
				return nil
			}
			if f.IsDir() {
				return nil
			}
			if !strings.HasSuffix(f.Name(), "_test.go") {
				return nil
			}
			return os.Remove(path)
		})
		if err != nil {
			d.log("Error while wiping test files: %v\n", err)
		}
	}
}

func (d *Deep) wipeVCS(pwd string, packages []Package) {
	for _, pkg := range packages {
		for _, vcsDir := range d.vcsDirs {
			vcsDirPath := pkg.vendoredPath(pwd) + pathSeparatorString + vcsDir
			err := os.RemoveAll(vcsDirPath)
			if err != nil {
				d.log("Error while removing vcs dir for package: %s\n", pkg.Name)
			}
		}
	}
}

func (d *Deep) writeDeepFiles(pwd, currentPkg string, packages []Package) {
	// TODO We shouldn't have to do this to begin with
	for idx := range packages {
		packages[idx].Dependencies = nil
	}

	// TODO improve this when we know how to read the current package
	p := Package{
		Name:         currentPkg,
		Version:      "HEAD",
		Dependencies: packages,
	}

	m := &Manifest{
		Package: p,
	}
	err := m.writeFile(pwd)
	if err != nil {
		d.log("Error while marshaling the manifest file.\nGot error: %v\n", err)
		os.Exit(1)
	}

	l := &Lock{
		Package: p,
	}
	err = l.writeFile(pwd)
	if err != nil {
		d.log("Error while marshaling the lock file.\nGot error: %v\n", err)
		os.Exit(1)
	}
}

// Run will execute all operations needed in order to vendor the the project.
func (d *Deep) Run(pwd, currentPkg string, keepTypes map[string]struct{}, args []string) {
	if currentPkg == "" {
		d.log("Current package is empty. Are you running on a project from GOPATH?")
		os.Exit(1)
	}

	packages := d.listPackages(pwd, currentPkg)
	if len(packages) == 0 {
		d.log("No packages found")
		return
	}

	d.vendorPackages(pwd, currentPkg, packages)

	d.readCommitHashes(pwd, currentPkg, packages)

	d.wipeNestedVendor(pwd, currentPkg, packages)

	if _, ok := keepTypes["vcs"]; !ok {
		d.wipeVCS(pwd, packages)
	}

	if _, ok := keepTypes["test"]; !ok {
		d.wipeTestFiles(pwd, currentPkg, packages)
	}

	// TODO implement
	/*if _, ok := keepTypes["main"]; !ok {
		d.wipeMainFiles(pwd, currentPkg, packages)
	}*/

	d.writeDeepFiles(pwd, currentPkg, packages)
}

// New creates a new instance of Deep
func New(logger Logger) *Deep {
	vcsDirs := []string{
		".bzr",
		".git",
		".hg",
		".svn",
	}

	providers := []provider{
		newDeep(logger),
		newGoList(logger),
	}

	return &Deep{
		log:       logger,
		vcsDirs:   vcsDirs,
		providers: providers,
	}
}
