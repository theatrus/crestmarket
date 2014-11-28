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
	"fmt"
	"github.com/theatrus/oauth2"
	"io/ioutil"
	"net/http"
	"strconv"
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
	root *Root
}

// The base type of fetcher for all CREST data types.
type CRESTRequestor interface {
	Root() (*Root, error)
	Regions() (*Regions, error)
	Types() error
}

func NewCrestRequestor(transport *oauth2.Transport) (CRESTRequestor, error) {
	req := requestor{transport, nil}

	root, err := req.Root()
	if err != nil {
		return nil, err
	}
	req.root = root
	return &req, nil
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

		href := itemMap["href"].(string)
		idSplit := strings.Split(href, "/")
		id, err := strconv.ParseInt(idSplit[len(idSplit)-2], 10, 64)
		if err != nil {
			return nil, err
		}

		region := Region{itemMap["name"].(string), href, int(id)}
		regions.AllRegions = append(regions.AllRegions, &region)
	}
	return &regions, nil
}

func (o *requestor) Regions() (*Regions, error) {
	path := o.root.Resources["regions"]
	body, err := fetch(path, o.transport)
	if err != nil {
		return nil, err
	}
	regions, err := unpackRegions(body)

	return regions, nil
}

func (o *requestor) Types() error {
	body, err := fetch("/", o.transport)
	if err != nil {
		return err
	}
	fmt.Printf("%s", body)
	return nil
}

func (o *requestor) Root() (*Root, error) {

	body, err := fetch("/", o.transport)
	if err != nil {
		return nil, err
	}
	return unpackRoot(body)
}

func unpackRoot(body []byte) (*Root, error) {
	var root Root
	root.Resources = make(map[string]string)

	rroots := make(map[string]interface{})
	if err := json.Unmarshal(body, &rroots); err != nil {
		return nil, err
	}

	for service, item := range rroots {
		itemM, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		href, ok := itemM["href"].(string)
		if ok {
			root.Resources[service] = href
		}
	}
	return &root, nil
}

// Peform a URL fetch and read into a []byte
func fetch(path string, transport *oauth2.Transport) ([]byte, error) {
	req, err := newCrestRequest(path)
	if err != nil {
		return nil, err
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func newCrestRequest(path string) (*http.Request, error) {
	var finalPath = path
	if !strings.HasPrefix(path, "http") {
		finalPath = prefix + finalPath
	}
	req, err := http.NewRequest("GET", finalPath, nil)
	req.Header.Add("Accept", accept)
	return req, err
}
