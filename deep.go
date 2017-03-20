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

import "errors"

type deep struct {
		logger Logger
	}

func (d *deep) canUse(pwd, currentPkg string) bool {
	return false
}

func (d *deep) packages(pwd, currentPkg string, ignoreTestFiles bool) ([]Package, error) {
	return nil, errors.New("not implemented yet")
}

func (d *deep) loadFiles() (*Manifest, *Lock, error) {
	return nil, nil, nil
}

func newDeep(logger Logger) *deep {
	return &deep{
		logger: logger,
	}
}
