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
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

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

func NewOAuthConfig(settings *OAuthSettings) *oauth2.Config {
	var endpoint oauth2.Endpoint
	if isSisi {
		endpoint = oauth2.Endpoint{
			AuthURL:  "https://sisilogin.testeveonline.com/oauth/authorize",
			TokenURL: "https://sisilogin.testeveonline.com/oauth/token",
		}
	} else {
		endpoint = oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/oauth/authorize",
			TokenURL: "https://login.eveonline.com/oauth/token",
		}
	}

	return &oauth2.Config{
		ClientID:     settings.ClientId,
		ClientSecret: settings.ClientSecret,
		RedirectURL:  settings.Callback,
		Scopes:       []string{"publicData"},
		Endpoint:     endpoint,
	}
}
