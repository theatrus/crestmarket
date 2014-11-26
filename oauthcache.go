package crestmarket

import (
	"encoding/json"
	"github.com/theatrus/oauth2"
	"io/ioutil"
	"log"
)

type FileTokenStore struct {
	Filename string
}

func FileToken(f *FileTokenStore) oauth2.Option {
	return func(o *oauth2.Options) error {
		o.TokenStore = f
		return nil
	}
}

func (o *FileTokenStore) ReadToken() (*oauth2.Token, error) {
	fileContents, err := ioutil.ReadFile(o.Filename)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(fileContents, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (o *FileTokenStore) WriteToken(token *oauth2.Token) {
	log.Println("writing token")
	data, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(o.Filename, data, 0600)
	if err != nil {
		log.Fatal(err)
	}

	return
}
