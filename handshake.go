package crestmarket

import (
	"encoding/json"
	"flag"
	"github.com/theatrus/oauth2"
	"io/ioutil"
	"log"
)

var isSisi bool

func init() {
	flag.BoolVar(&isSisi, "scanner.sisi", true, "Call all endpoints on SiSi, turn off for production")
}

type OAuthSettings struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Callback     string `json:"callback"`
}

func LoadSettings(filename string) (*OAuthSettings, error) {
	var settings OAuthSettings
	settingsData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Can't load secret key file - aborting", err)
		return nil, err
	}
	json.Unmarshal(settingsData, &settings)
	return &settings, nil
}

func NewOauthOptions(settings *OAuthSettings) (*oauth2.Options, error) {
	var endpoint oauth2.Option
	if isSisi {
		endpoint = oauth2.Endpoint(
			"https://sisilogin.testeveonline.com/oauth/authorize",
			"https://sisilogin.testeveonline.com/oauth/token",
		)
	} else {
		endpoint = oauth2.Endpoint(
			"https://login.eveonline.com/oauth/authorize",
			"https://login.eveonline.com/oauth/token",
		)
	}
	return oauth2.New(
		oauth2.Client(settings.ClientId, settings.ClientSecret),
		oauth2.RedirectURL(settings.Callback),
		oauth2.Scope("publicData"),
		endpoint,
	)
}
