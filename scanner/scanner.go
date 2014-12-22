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

package main

import (
	"bytes"
	"flag"
	"github.com/theatrus/crestmarket"
	"github.com/theatrus/crestmarket/helper"
	"github.com/theatrus/mediate"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var uploadEndpoint string
var onlyRegion int

func init() {
	flag.StringVar(&uploadEndpoint, "scanner.upload", "http://localhost", "Default upload endpoint")
	flag.IntVar(&onlyRegion, "scanner.region", 0, "Limit to a specific region")
}

func scanRegion(req crestmarket.CRESTRequestor,
	region *crestmarket.Region, forItems *crestmarket.MarketTypes) {

	perm := rand.Perm(len(forItems.Types))
	dest := make([]*crestmarket.MarketType, len(forItems.Types))
	for i, v := range perm {
		dest[v] = forItems.Types[i]
	}

	for _, item := range dest {

		if strings.Contains(item.Name, "Blueprint") {
			log.Printf("Skipping Blueprint %s", item)
			continue
		}

		mo, err := req.BuySellMarketOrders(region, item)
		if err != nil {
			log.Printf("fetch error: %s\n", err)
			continue
		}
		md, err := crestmarket.SerializeOrdersUnified(mo, time.Now())
		if err != nil {
			log.Printf("Deserialize error: %s\n", err)
			continue
		}
		reader := bytes.NewReader(md)
		resp, err := http.Post(uploadEndpoint, "application/json", reader)
		if err != nil {
			log.Printf("Error posting market data %s\n", err)
			continue
		}

		_, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	log.Printf("DONE SCANNING %s", region)
}

func main() {

	flag.Parse()

	settings, err := crestmarket.LoadSettings("settings.json")
	if err != nil {
		log.Fatal(err)
	}

	transport, err := helper.InteractiveStartup("token.json", settings)
	if err != nil {
		log.Fatal(err)
	}
	requestor, err := crestmarket.NewCrestRequestor(mediate.RateLimit(50, 1*time.Second, transport))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bootstrapping items and regions")
	itemsChan := make(chan *crestmarket.MarketTypes)
	go func(done chan<- *crestmarket.MarketTypes) {
		items, err := requestor.Types()
		if err != nil {
			log.Fatal(err)
		}
		done <- items
	}(itemsChan)

	regionsChan := make(chan *crestmarket.Regions)
	go func(done chan<- *crestmarket.Regions) {
		regions, err := requestor.Regions()
		if err != nil {
			log.Fatal(err)
		}
		done <- regions
	}(regionsChan)

	items := <-itemsChan
	regions := <-regionsChan

	if onlyRegion != 0 {
		log.Printf("Limiting scanning to region %d\n", onlyRegion)
	}

	// Remove regions which are not marketable (this is a zero-allocation filter due to slices)
	filteredRegions := regions.AllRegions[:0]
	for _, r := range regions.AllRegions {
		if r.Id < 11000000 || r.Id == 11000031 {
			filteredRegions = append(filteredRegions, r)
		}
	}

	for {
		var wg sync.WaitGroup
		if onlyRegion != 0 {
			region := regions.ById(onlyRegion)
			for i := 0; i < 40; i++ {
				wg.Add(1)
				go func(region *crestmarket.Region) {
					defer wg.Done()
					scanRegion(requestor, region, items)
				}(region)
			}
		} else {
			log.Println("Starting market scrape, parallelizing by region")
			for _, region := range filteredRegions {
				wg.Add(1)
				go func(region *crestmarket.Region) {
					defer wg.Done()
					scanRegion(requestor, region, items)
				}(region)
			}
		}
		wg.Wait()
		log.Println("Done.")
	}

}
