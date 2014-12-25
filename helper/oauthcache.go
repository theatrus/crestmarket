//   Copyright 2014 StackFoundry LLC
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package helper

import (
	"encoding/json"
	"github.com/theatrus/ooauth2"
	"io/ioutil"
	"log"
)

type FileTokenStore struct {
	Filename string
}

func FileToken(f *FileTokenStore) ooauth2.Option {
	return func(o *ooauth2.Options) error {
		o.TokenStore = f
		return nil
	}
}

func (o *FileTokenStore) ReadToken() (*ooauth2.Token, error) {
	fileContents, err := ioutil.ReadFile(o.Filename)
	if err != nil {
		return nil, err
	}
	var token ooauth2.Token
	err = json.Unmarshal(fileContents, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (o *FileTokenStore) WriteToken(token *ooauth2.Token) {
	data, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(o.Filename, data, 0600)
	if err != nil {
		log.Fatal(err)
	}

	return
}
