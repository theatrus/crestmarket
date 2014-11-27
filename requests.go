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
	"encoding/json"
	"errors"
	"github.com/theatrus/oauth2"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	prefix = "https://api-sisi.testeveonline.com"
	// The requested Accept header - this is used for all requests, even
	// if it doesn't actually make sense to be a MarketTypeCollection?
	accept = "application/vnd.ccp.eve.MarketTypeCollection-v1+json"
)

type requestor struct {
	transport *oauth2.Transport
}

// The base type of fetcher for all CREST data types.
type CRESTRequestor interface {
	Regions() (*Regions, error)
}

func NewCrestRequestor(transport *oauth2.Transport) CRESTRequestor {
	return &requestor{transport}
}

func unpackRegions(body []byte) (*Regions, error) {
	// Eschewing the normal tagged
	// unpacking here as the return structure is not
	// ideal for the in-application representation

	regions := Regions{make([]*Region, 0)}
	raw := make(map[string]interface{})
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	items, ok := raw["items"].([]interface{})
	if !ok {
		return nil, errors.New("Can't find an items key when unpacking regions")
	}

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("Can't unpack a region")
		}
		region := Region{itemMap["name"].(string), itemMap["href"].(string), 0}
		regions.AllRegions = append(regions.AllRegions, &region)
	}
	return &regions, nil
}

func (o *requestor) Regions() (*Regions, error) {

	req, err := NewCrestRequest("/regions/")
	if err != nil {
		return nil, err
	}

	resp, err := o.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	regions, err := unpackRegions(body)

	return regions, nil
}

func NewCrestRequest(path string) (*http.Request, error) {
	var finalPath = path
	if !strings.HasPrefix(path, "http") {
		finalPath = prefix + finalPath
	}
	req, err := http.NewRequest("GET", finalPath, nil)
	req.Header.Add("Accept", accept)
	return req, err
}
