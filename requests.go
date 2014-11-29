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
	"time"
)

const (
	prefix = "https://api-sisi.testeveonline.com"

	// The root resource version this library will work with.
	rootAccept = "application/vnd.ccp.eve.Api-v3+json"

	RFC3339SansTz = "2006-01-02T15:04:05"
)

// Basic definitions of resource types
var resourceVersions map[string]string

func init() {
	resourceVersions = map[string]string{
		"crestEndpoint": rootAccept,
		"regions":       "application/vnd.ccp.eve.RegionCollection-v1+json",
		"itemTypes":     "application/vnd.ccp.eve.ItemTypeCollection-v1+json",
		"marketOrders":  "application/vnd.ccp.eve.MarketOrderCollection-v1+json",
		"marketTypes":   "application/vnd.ccp.eve.MarketOrderCollection-v1+json",
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
	Types() (*MarketTypes, error)
	// Market orders
	MarketOrders(region *Region, mtype *MarketType, buy bool) (*MarketOrders, error)
}

func NewCrestRequestor(transport *oauth2.Transport) (CRESTRequestor, error) {
	req := requestor{transport, nil}

	root, err := req.fetchRoot()
	if err != nil {
		return nil, err
	}

	// This doesn't appear in the root resources yet, but are part of it
	root.Resources["marketOrders"] = "https://api-sisi.testeveonline.com/market/"

	req.root = root
	return &req, nil
}

////////////////////////////////////////////////////////////////////////////

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

func unpackMarketTypes(its *MarketTypes, page *page) error {
	items := page.items

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return errors.New("Can't unpack an marketType")
		}

		mtype, ok := itemMap["type"].(map[string]interface{})

		href := mtype["href"].(string)
		id := mtype["id"].(float64)

		it := MarketType{mtype["name"].(string), href, int(id)}
		its.Types = append(its.Types, &it)
	}
	return nil
}

func unpackMarketOrders(mo *MarketOrders, mt *MarketType, page *page) error {
	items := page.items

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return errors.New("Can't unpack an order")
		}
		buy := itemMap["buy"].(bool)
		duration := int(itemMap["duration"].(float64))
		href := itemMap["href"].(string)
		issued, err := time.Parse(RFC3339SansTz, itemMap["issued"].(string))
		if err != nil {
			return err
		}
		minVolume := int(itemMap["minVolume"].(float64))
		price := itemMap["price"].(float64)
		mrange := itemMap["range"].(string)
		volume := int(itemMap["volume"].(float64))

		locationMap := itemMap["location"].(map[string]interface{})
		location := Location{locationMap["name"].(string), locationMap["href"].(string)}

		mo.Orders = append(mo.Orders,
			&MarketOrder{buy, duration,
				href, issued, location,
				minVolume, price, mrange,
				*mt, volume})
	}
	return nil
}

// Deserialize the json for the root object into a Root
// Nested resources are delineated with a "/", but
// are unpacked via heuristics.
func unpackRoot(body []byte) (*Root, error) {
	var root Root
	root.Resources = make(map[string]string)
	rroots := make(map[string]interface{})
	if err := json.Unmarshal(body, &rroots); err != nil {
		return nil, err
	}

	var recurse func(string, map[string]interface{})

	recurse = func(base string, items map[string]interface{}) {
		for service, item := range items {
			itemM, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			href, ok := itemM["href"].(string)
			if ok {
				root.Resources[base+service] = href
			} else if len(itemM) > 0 {
				recurse(base+service+"/", itemM)
			}
		}
	}
	recurse("", rroots)
	return &root, nil
}

//////////////////////////////////////////////////////////////////////////

func (o *requestor) walkPages(path string, resource string, extractor func(*page) error) error {
	for {
		body, err := o.fetchWithRetry(path, resource)
		if err != nil {
			return err
		}
		page, err := unpackPage(body)
		if err != nil {
			return err
		}
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
	err := o.walkPages(path, "regions",
		func(page *page) error { return unpackRegions(regions, page) })
	if err != nil {
		return nil, err
	}

	return regions, nil
}

func (o *requestor) Types() (*MarketTypes, error) {
	path := o.root.Resources["marketTypes"]
	its := newMarketTypes()
	err := o.walkPages(path, "marketTypes",
		func(page *page) error { return unpackMarketTypes(its, page) })
	if err != nil {
		return nil, err
	}

	return its, nil
}

func (o *requestor) MarketOrders(region *Region, mtype *MarketType, buy bool) (*MarketOrders, error) {
	var orderType string
	if buy {
		orderType = "buy"
	} else {
		orderType = "sell"
	}

	marketOrders := NewMarketOrders()

	path := o.root.Resources["marketOrders"]
	path = fmt.Sprintf("%s%d/orders/%s/?type=%s", path, region.Id, orderType, mtype.Href)
	err := o.walkPages(path,
		"marketOrders",
		func(page *page) error {
			return unpackMarketOrders(marketOrders, mtype, page)
		})
	return marketOrders, err
}

func (o *requestor) Root() (*Root, error) {
	return o.root, nil
}

func (o *requestor) fetchRoot() (*Root, error) {

	body, err := o.fetchWithRetry("/", "crestEndpoint")
	if err != nil {
		return nil, err
	}
	return unpackRoot(body)
}

func (o *requestor) fetchOptions(path string) ([]byte, error) {
	transport := o.transport

	req, err := o.newCrestRequest(path, "OPTIONS", false, "crestEndpoint")
	if err != nil {
		return nil, err
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("%s Non-200 status code returned from fetching. %d %s\n",
			path, resp.StatusCode, body)
		return nil, errors.New(fmt.Sprintf("Resource not found: %s", path))
	}
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Do a few retries before giving up - idempotent operations only
func (o *requestor) fetchWithRetry(path string, resource string) ([]byte, error) {
	retriesLeft := 3
	// Currently, this retries all requests, not just retryable ones.
	// At times, Finagle is nice
	for ; ; retriesLeft-- {
		result, err := o.fetch(path, resource)
		if err != nil {
			log.Printf("Error, retrying: %s", err)
			if retriesLeft <= 0 {
				return nil, err
			}
		} else {
			return result, nil
		}
	}
}

// Peform a URL fetch and read into a []byte
func (o *requestor) fetch(path string, resource string) ([]byte, error) {
	transport := o.transport

	req, err := o.newCrestRequest(path, "GET", true, resource)
	if err != nil {
		return nil, err
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if dep, ok := resp.Header["X-Deprecated"]; ok {
		log.Println("Deprecated API: ", dep)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("%s Non-200 status code returned from fetching. %d %s\n",
			path, resp.StatusCode, body)
		return nil, errors.New(fmt.Sprintf("Resource not found: %s", path))
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *requestor) newCrestRequest(path string,
	verb string,
	addAccept bool,
	resource string) (*http.Request, error) {

	var finalPath = path
	if !strings.HasPrefix(path, "http") {
		finalPath = prefix + finalPath
	}
	var accept string
	if addAccept {
		// Find resource root to pass the appropiate known accept header
		if finalPath == prefix + "/" || o.root == nil {
			// Root path is a special case
			accept = rootAccept
		} else {
			accept = resourceVersions[resource]
		}
	}
	req, err := http.NewRequest(verb, finalPath, nil)
	if accept != "" {
		accept = accept + "; charset=utf-8"
	} else {
		accept = "charset=utf-8"
	}

	req.Header.Add("Accept", accept)

	return req, err
}
