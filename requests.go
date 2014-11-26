package crestmarket

import (
	"net/http"
)

const (
	prefix = "https://api-sisi.testeveonline.com"
	// The requested Accept header - this is used for all requests, even
	// if it doesn't actually make sense to be a MarketTypeCollection?
	accept = "application/vnd.ccp.eve.MarketTypeCollection-v1+json"
)

func NewCrestRequest(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", prefix + path, nil)
	req.Header.Add("Accept", accept)
	return req, err
}
