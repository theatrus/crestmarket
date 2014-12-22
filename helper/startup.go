package helper

import (
	"fmt"
	"github.com/theatrus/crestmarket"
	"github.com/theatrus/ooauth2"
	"log"
	"net/http"
)

// Perform an *interactive* *console* handshake. This requires the user
// opening a URL manually, and then pasting the resultant code back into
// this application. The other approach is a multi-invocation token-fetcher.
func InteractiveHandshake(settings *crestmarket.OAuthSettings, store *FileTokenStore) (*ooauth2.Transport, error) {
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

// BackgroundStartup tries to start with the current token
// or fails with a Fatal if the token can't be read or is expired.
func BackgroundStartup(tokenFile string, settings *crestmarket.OAuthSettings) (http.RoundTripper, error) {
	store := FileTokenStore{Filename: tokenFile}

	base, err := crestmarket.NewOAuthOptions(settings)
	t, err := base.NewTransportFromTokenStore(&store)
	if err != nil {
		log.Fatal("Token refresh has failed")
	}
	if t.Token().Expired() {
		log.Fatal("Token is expired and refresh has failed.")
	}
	store.WriteToken(t.Token())
	return t, nil

}

// InteractiveStartup performs a console interactive handshake
// or a simple refresh of tokens and stores
// tokens gathered in a file called token.json
func InteractiveStartup(tokenFile string, settings *crestmarket.OAuthSettings) (http.RoundTripper, error) {
	store := FileTokenStore{Filename: tokenFile}

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
	return t, nil
}
