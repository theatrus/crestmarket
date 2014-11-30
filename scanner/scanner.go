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
	"fmt"
	"github.com/theatrus/crestmarket"
	"github.com/theatrus/crestmarket/helper"
	"log"
)

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

	theForge := regions.ByName("The Forge")
	fmt.Println(theForge)

	trit := items.ByName("Tritanium")
	fmt.Println(trit)

	mo, err := requestor.BuySellMarketOrders(theForge, trit)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", mo)
}
