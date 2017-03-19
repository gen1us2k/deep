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
	"os/exec"
	"strings"
)

// This is a temporary provider until we have a proper package importer in place
// that won't have the various limitations of this provider
type goList struct {
	log Logger
}

func (*goList) canUse(pwd, currentPkg string) bool {
	return true
}

func (g *goList) packages(pwd, currentPkg string) ([]Package, error) {
	output, err := exec.Command("go", "list", "-f", `{{ join .Deps "\n" }}`, currentPkg+"/...").Output()
	if err != nil {
		return nil, err
	}
	pkgList := strings.Split(string(output), "\n")

	var result []Package
	for _, pkg := range pkgList {
		if pkg == "" {
			continue
		}

		p := Package{
			Name:    pkg,
			Version: "HEAD",
		}

		if !p.isRootPackage(currentPkg) {
			continue
		}

		result = append(result, p)
	}

	return result, nil
}

func newGoList(logger Logger) *goList {
	return &goList{
		log: logger,
	}
}
