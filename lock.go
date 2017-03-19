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
	"time"
)

// Lock defines the format for the lock file used by Deep to pin the packages in known revisions
type Lock struct {
	WritenAt time.Time `json:"writen_at"`
	Package
}

const lockFileName = ".deep_lock.json"

func (l *Lock) writeFile(path string) error {
	l.WritenAt = time.Now()

	lk, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path+pathSeparatorString+lockFileName, lk, 0644)
}
