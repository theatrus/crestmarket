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
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

func NewFileTokenSource(filename string) *FileTokenSource {
	return &FileTokenSource{filename, new(oauth2.Token)}
}

type FileTokenSource struct {
	Filename    string
	CachedToken *oauth2.Token
}

// Implementing oauth2.TokenSource interface:
func (o FileTokenSource) Token() (*oauth2.Token, error) {
	// Checking the AccessToken below is a hack. Ideally we would
	// check o.CachedToken.Valid(). Unfortunately, Valid() checks
	// the expiration time stamp, and returns false if expired.
	// This would cause this function to always read and parse
	// the file on every call because we never write back to the
	// file after OAuth2 refreshes the token.
	// TODO: switch to Valid() once we properly write to the file
	// on token refreshes.
	if o.CachedToken.AccessToken != "" {
		return o.CachedToken, nil
	}
	fileContents, err := ioutil.ReadFile(o.Filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileContents, o.CachedToken)
	if err != nil {
		return nil, err
	}
	return o.CachedToken, nil
}

func (o *FileTokenSource) WriteTokenToFile() {
	data, err := json.Marshal(o.CachedToken)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(o.Filename, data, 0600)
	if err != nil {
		log.Fatal(err)
	}
	return
}
