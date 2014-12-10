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
	"flag"
	"github.com/theatrus/crestmarket"
	"github.com/theatrus/crestmarket/helper"
	"log"
	"sync"
	"time"
)

func scanRegion(req crestmarket.CRESTRequestor,
	region *crestmarket.Region, forItems *crestmarket.MarketTypes) {
	for _, item := range forItems.Types {
		mo, err := req.BuySellMarketOrders(region, item)
		if err != nil {
			log.Println("fetch error: %s", err)
			continue
		}
		_, err = crestmarket.SerializeOrdersUnified(mo, time.Now())
		if err != nil {
			log.Println("Deserialize error: %s", err)
			continue
		}
	}
}

func main() {

	flag.Parse()

	settings, err := crestmarket.LoadSettings("settings.json")
	if err != nil {
		log.Fatal(err)
	}

	requestor, err := helper.InteractiveStartup(settings)
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

	log.Println("Starting market scrape, parallelizing by region")

	var wg sync.WaitGroup
	for _, region := range regions.AllRegions {
		wg.Add(1)
		go scanRegion(requestor, region, items)
	}
	wg.Wait()

	//theForge := regions.ByName("The Forge")
	//fmt.Println(theForge)

	//trit := items.ByName("Tritanium")
	//fmt.Println(trit)

	//mo, err := requestor.BuySellMarketOrders(theForge, trit)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//serial, err := crestmarket.SerializeOrdersUnified(mo, time.Now())
	//if err != nil {
	//	log.Fatal(Err)
	//}
	//fmt.Printf("%s\n", serial)
}
