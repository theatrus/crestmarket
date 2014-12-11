package main

import (
	"flag"
	"fmt"
	"github.com/theatrus/crestmarket"
	"github.com/theatrus/crestmarket/helper"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("proxy.port", 12345, "Port for the CREST Proxy")
	flag.Parse()

	settings, err := crestmarket.LoadSettings("settings.json")
	if err != nil {
		log.Fatal(err)
	}

	transport, err := helper.InteractiveStartup(settings)
	if err != nil {
		log.Fatal(err)
	}
	requestor, err := crestmarket.NewCrestRequestor(transport)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://localhost:%d", *port)
	listen := fmt.Sprintf("localhost:%d", *port)

	proxy := &crestmarket.CRESTProxy{Requestor: requestor, NewURLRoot: url}

	log.Printf("Starting server at %s\n", listen)

	err = http.ListenAndServe(listen, proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
