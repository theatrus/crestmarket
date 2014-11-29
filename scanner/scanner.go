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
	"github.com/theatrus/oauth2"
	"log"
)

// Perform an *interactive* *console* handshake. This requires the user
// opening a URL manually, and then pasting the resultant code back into
// this application. The other approach is a multi-invocation token-fetcher.
func newHandshake(settings *crestmarket.OAuthSettings, store *helper.FileTokenStore) (*oauth2.Transport, error) {
	f, err := crestmarket.NewOauthOptions(settings)
	f.TokenStore = store
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := f.AuthCodeURL("state", "online", "auto")
	fmt.Println("Visit the URL for the auth dialog:")
	fmt.Println(url)
	fmt.Println()
	fmt.Printf("Auth code> ")

	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.
	var code string
	if _, err = fmt.Scan(&code); err != nil {
		log.Fatal(err)
		return nil, err
	}
	t, err := f.NewTransportFromCode(code)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return t, nil
}

func main() {

	flag.Parse()

	settings, err := crestmarket.LoadSettings("settings.json")
	if err != nil {
		log.Fatal(err)
	}

	store := helper.FileTokenStore{"token.json"}

	base, err := crestmarket.NewOauthOptions(settings)
	t, err := base.NewTransportFromTokenStore(&store)
	if err != nil {
		log.Println("Token refresh has failed, requesting new authorization interactively")
		t, err = newHandshake(settings, &store)
		if err != nil {
			log.Fatal("Can't really continue, auth has failed.")
			return
		}
	}
	// Need to manually flush the token store at auth for now
	store.WriteToken(t.Token())

	requestor, err := crestmarket.NewCrestRequestor(t)
	if err != nil {
		log.Fatal(err)
	}

	items, err := requestor.Types()
	if err != nil {
		log.Fatal(err)
	}

	regions, err := requestor.Regions()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s", regions)

	theForge := regions.ByName("The Forge")
	fmt.Println(theForge)

	trit := items.ByName("Tritanium")
	fmt.Println(trit)

	mo, err := requestor.MarketOrders(theForge, trit, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", mo)
}
