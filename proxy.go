package crestmarket

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type CRESTProxy struct {
	Requestor  CRESTRequestor
	NewURLRoot string
}

var matchHost *regexp.Regexp

func init() {
	matchHost = regexp.MustCompile("https://[a-z-]+\\.[a-z-]+\\.com")
}

func (c *CRESTProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.serve(rw, req)
}

func (c *CRESTProxy) serve(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if req.URL.RawQuery != "" {
		path += "?" + req.URL.RawQuery
	}

	cast_req, ok := c.Requestor.(*requestor)
	if !ok {
		log.Fatal("Can't use proxy with this form of requestor")
	}
	body, err := cast_req.fetch(path, "")

	if err != nil {
		rw.WriteHeader(500)
		msg := fmt.Sprintf("%s", err)
		_, _ = rw.Write([]byte(msg))
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(200)

	newBody := matchHost.ReplaceAll(body, []byte(c.NewURLRoot))

	_, _ = rw.Write(newBody)
}
