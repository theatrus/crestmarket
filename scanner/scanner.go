package main

import (
	"encoding/json"
	"fmt"
	"github.com/theatrus/crestmarket"
	"github.com/theatrus/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

type OAuthSettings struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Callback     string `json:"callback"`
}

func baseOauth(settings OAuthSettings) (*oauth2.Options, error) {
	return oauth2.New(
		oauth2.Client(settings.ClientId, settings.ClientSecret),
		oauth2.RedirectURL(settings.Callback),
		oauth2.Scope("publicData"),
		oauth2.Endpoint(
			"https://sisilogin.testeveonline.com/oauth/authorize",
			"https://sisilogin.testeveonline.com/oauth/token",
		),
	)
}

func newHandshake(settings OAuthSettings, store *crestmarket.FileTokenStore) (*oauth2.Transport, error) {
	f, err := baseOauth(settings)
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

	var settings OAuthSettings
	settingsData, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Fatal("Can't load secret key file - aborting")
		return
	}
	json.Unmarshal(settingsData, &settings)

	store := crestmarket.FileTokenStore{"token.json"}
	base, err := baseOauth(settings)
	t, err := base.NewTransportFromTokenStore(&store)
	if err != nil {
		log.Println("Token refresh has failed, requesting new authorization interactively")
		t, err = newHandshake(settings, &store)
		if err != nil {
			log.Fatal("Can't really continue, auth has failed.")
			return
		}
	}

	fmt.Println(t.Token().AccessToken)
	store.WriteToken(t.Token())

	for i := 0; i < 100; i++ {

		req, err := http.NewRequest("GET", "https://api-sisi.testeveonline.com/market/10000002/orders/buy/?type=https://api-sisi.testeveonline.com/inventory/types/683/", nil)
		resp, err := t.RoundTrip(req)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", body)
	}

}
