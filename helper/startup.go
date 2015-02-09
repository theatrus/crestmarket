package helper

import (
	"fmt"
	"github.com/theatrus/crestmarket"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

// Perform an *interactive* *console* handshake. This requires the user
// opening a URL manually, and then pasting the resultant code back into
// this application. The other approach is a multi-invocation token-fetcher.
func InteractiveHandshake(config *oauth2.Config) (*oauth2.Token, error) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := config.AuthCodeURL("state", oauth2.AccessTypeOnline)
	fmt.Println("Visit the URL for the auth dialog:")
	fmt.Println(url)
	fmt.Println()
	fmt.Printf("Auth code> ")

	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
		return nil, err
	}

	// TODO: Explore how to get Mediate's retrying transport layer adapted
	// into the Context-to-http.Client lookup table before making the
	// exchange.
	// Would be nice if oauth2 exported registerContextClientFunc() ...bummer
	freshToken, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return freshToken, nil
}

// InteractiveStartup performs a console interactive handshake
// or a simple refresh of tokens and stores
// tokens gathered in a file called token.json
func InteractiveStartup(tokenFile string, settings *crestmarket.OAuthSettings) (http.RoundTripper, error) {
	config := crestmarket.NewOAuthConfig(settings)
	source := NewFileTokenSource(tokenFile)

	// Try the new FileTokenSource, if it doesn't produce a valid token, force
	// the user into the interactive prompt.
	if token, err := source.Token(); err != nil || token.AccessToken == "" {
		log.Println("Token is not valid, requesting new authorization interactively")
		source.CachedToken, err = InteractiveHandshake(config)
		if err != nil {
			log.Println("Can't really continue, auth has failed.")
			return nil, err
		}
		// TODO: Update the token whenever it gets refreshed by OAuth2.
		// Context: We only write to the file if the user goes through
		// interactive handshake. If the token file contains a valid but
		// expired token, OAuth2 will (internally) go refresh the token on
		// its first RoundTrip() with the server. We never end up writing
		// the refreshed token back to the file.
		// Perhaps put this functionality into a time.NewTicker?
		source.WriteTokenToFile()
	}

	transport := &oauth2.Transport{
		// Use Config.TokenSource to make use of oauth2.tokenRefresher{...}
		Source: config.TokenSource(oauth2.NoContext, source.CachedToken),
	}

	return transport, nil
}
