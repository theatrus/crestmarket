package helper

import (
	"fmt"
	"github.com/theatrus/crestmarket"
	"github.com/theatrus/oauth2"
	"log"
)

// Perform an *interactive* *console* handshake. This requires the user
// opening a URL manually, and then pasting the resultant code back into
// this application. The other approach is a multi-invocation token-fetcher.
func InteractiveHandshake(settings *crestmarket.OAuthSettings, store *FileTokenStore) (*oauth2.Transport, error) {
	f, err := crestmarket.NewOAuthOptions(settings)
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

func InteractiveStartup(settings *crestmarket.OAuthSettings) (crestmarket.CRESTRequestor, error) {
	store := FileTokenStore{Filename: "token.json"}

	base, err := crestmarket.NewOAuthOptions(settings)
	t, err := base.NewTransportFromTokenStore(&store)
	if err != nil {
		log.Println("Token refresh has failed, requesting new authorization interactively")
		t, err = InteractiveHandshake(settings, &store)
		if err != nil {
			log.Println("Can't really continue, auth has failed.")
			return nil, err
		}
	}
	// Need to manually flush the token store at auth for now
	store.WriteToken(t.Token())

	requestor, err := crestmarket.NewCrestRequestor(t)
	if err != nil {
		return nil, err
	}
	return requestor, nil
}
