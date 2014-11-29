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
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	prefix = "https://api-sisi.testeveonline.com"

	// The root resource version this library will work with.
	rootAccept = "application/vnd.ccp.eve.Api-v3+json"
)

// Basic definitions of resource types
var resourceVersions map[string]string

func init() {
	resourceVersions = map[string]string{
		"root":      rootAccept,
		"regions":   "application/vnd.ccp.eve.RegionCollection-v1+json",
		"itemTypes": "application/vnd.ccp.eve.ItemTypeCollection-v1+json",
	}
}

type requestor struct {
	transport *oauth2.Transport
	root      *Root
}

// The base type of fetcher for all CREST data types.
type CRESTRequestor interface {
	// Return a new copy of the root resource
	Root() (*Root, error)
	// Return a list of all regions
	Regions() (*Regions, error)
	// Return a list of all known types
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

type page struct {
	items    []interface{}
	hasNext  bool
	nextHref string
}

// Unpack a page structure and extract optional next fields
// This is useful for a serial request structure - in order
// to parallelize page fetching different heuristics need to
// be used violating the API purity.
func unpackPage(body []byte) (*page, error) {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	items, ok := raw["items"].([]interface{})
	if !ok {
		return nil, errors.New("Can't find an items key in the response")
	}

	hasNext := false
	next := ""

	if nextHref, ok := raw["next"].(map[string]interface{}); ok {
		next = nextHref["href"].(string)
		hasNext = true
	}

	return &page{items, hasNext, next}, nil
}

func unpackRegions(regions *Regions, page *page) error {
	items := page.items

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return errors.New("Can't unpack a region")
		}

		href := itemMap["href"].(string)
		idSplit := strings.Split(href, "/")
		id, err := strconv.ParseInt(idSplit[len(idSplit)-2], 10, 64)
		if err != nil {
			return err
		}

		region := Region{itemMap["name"].(string), href, int(id)}
		regions.AllRegions = append(regions.AllRegions, &region)
	}
	return nil
}

func (o *requestor) walkPages(path string, extractor func(*page) error) error {
	for {
		body, err := o.fetch(path)
		if err != nil {
			return err
		}
		page, err := unpackPage(body)
		err = extractor(page)
		if err != nil {
			return err
		}
		if page.hasNext {
			path = page.nextHref
		} else {
			break
		}
	}
	return nil
}

func (o *requestor) Regions() (*Regions, error) {
	path := o.root.Resources["regions"]
	regions := newRegions()
	err := o.walkPages(path, func(page *page) error { return unpackRegions(regions, page) })
	if err != nil {
		return nil, err
	}

	return regions, nil
}

func (o *requestor) Types() error {
	path := o.root.Resources["itemTypes"]
	err := o.walkPages(path, func(page *page) error { fmt.Printf("%s", page); return nil })
	if err != nil {
		return nil
	}

	return nil
}

func (o *requestor) Root() (*Root, error) {

	body, err := o.fetch("/")
	if err != nil {
		return nil, err
	}
	return unpackRoot(body)
}

// Deserialize the json for the root object into a Root
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
func (o *requestor) fetch(path string) ([]byte, error) {
	transport := o.transport

	req, err := o.newCrestRequest(path)
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
	log.Println(resp.Header["Content-Type"])
	return body, nil
}

func (o *requestor) newCrestRequest(path string) (*http.Request, error) {
	var finalPath = path
	if !strings.HasPrefix(path, "http") {
		finalPath = prefix + finalPath
	}
	var accept string
	// Find resource root to pass the appropiate known accept header
	if finalPath == prefix+"/" || o.Root == nil {
		// Root path is a special case
		accept = rootAccept
	} else {
		// Iterate through to find prefixes
		for resource, prefix := range o.root.Resources {
			if strings.HasPrefix(finalPath, prefix) {
				accept = resourceVersions[resource]
				break
			}
		}
	}
	req, err := http.NewRequest("GET", finalPath, nil)
	if accept != "" {
		accept = accept + "; charset=utf-8"
	} else {
		accept = "charset=utf-8"
	}
	log.Println("Adding accept type of", accept)
	req.Header.Add("Accept", accept)

	return req, err
}
