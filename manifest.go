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
	"encoding/json"
	"io/ioutil"
)

// Manifest is the user facing manifest file that allows users to define their dependencies
// and specify which versions should be used. Also allows for overriding the versions detected
// by Deep
type Manifest struct {
	Package
}

const manifestFileName = "deep.json"

func (m *Manifest) writeFile(path string) error {
	man, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path+pathSeparatorString+manifestFileName, man, 0644)
}
