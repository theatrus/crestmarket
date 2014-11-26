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

package crestmarket

import (
	"net/http"
)

const (
	prefix = "https://api-sisi.testeveonline.com"
	// The requested Accept header - this is used for all requests, even
	// if it doesn't actually make sense to be a MarketTypeCollection?
	accept = "application/vnd.ccp.eve.MarketTypeCollection-v1+json"
)

func NewCrestRequest(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", prefix + path, nil)
	req.Header.Add("Accept", accept)
	return req, err
}
