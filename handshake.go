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
	"flag"
	"github.com/theatrus/mediate"
	"github.com/theatrus/oauth2"
	"io/ioutil"
	"log"
	"net/http"
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

func NewOAuthOptions(settings *OAuthSettings) (*oauth2.Options, error) {
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

	httpClient := &http.Client{}
	httpClient.Transport = mediate.FixedRetries(3,
		mediate.ReliableBody(http.DefaultTransport),
	)

	return oauth2.New(
		oauth2.Client(settings.ClientId, settings.ClientSecret),
		oauth2.RedirectURL(settings.Callback),
		oauth2.Scope("publicData"),
		oauth2.HTTPClient(httpClient),
		endpoint,
	)
}
