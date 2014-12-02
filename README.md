# crestmarket

[![GoDoc](https://godoc.org/github.com/theatrus/crestmarket?status.svg)](http://godoc.org/github.com/theatrus/crestmarket)
[![Build Status](https://travis-ci.org/theatrus/crestmarket.svg)](https://travis-ci.org/theatrus/crestmarket)
[![Gratipay](http://img.shields.io/gratipay/theatrus.svg)](https://gratipay.com/theatrus/)

A reference library for working with the authenticated EVE-Online
CREST endpoint, including market data. Licensed under the Apache
License, 2.0.

### Features

* OAuth2 negotiation and integration.
* A sample scanner application which demonstrates how to call
endpoints and receive data.
* Converters to output the
[Unified Uploader](http://dev.eve-central.com/unifieduploader/start)
format JSON for all market orders fetched.
* A "proxy server" to allow interactive exploration of the CREST
  structure without needing to navigate OAuth2 in the browser or in
  your own application.


## Getting started

Make sure you have Go 1.3+ installed. Make sure you have defined a
$GOPATH.

These commands can be run anywhere, though explicit paths are
referenced to the root of GOPATH:

```

go get github.com/theatrus/crestmarket
go install github.com/theatrus/crestmarket/scanner
cp src/github.com/theatrus/crestmarket/scanner/settings.json.example
settings.json
```

At this point, edit settings.json to include your CCP provided client
and secret. You will need a callback to receive your OAuth reply code.

```
bin/scanner
```

Running the scanner example will prompt you to open a URL, and to
paste in the reply code and hit enter. By default, this will also
cache your OAuth tokens in `token.json`

## CREST Proxy server

Included is a simple HTTP server which proxies URLs to CREST and
rewrites all URLs to point to local resources.

```
go install github.com/theatrus/crestmarket/crestproxy
cp src/github.com/theatrus/crestmarket/scanner/settings.json.example
settings.json
$GOPATH/bin/crestscanner
```

Complete any authorization needed on the console, and then open a web
browser to http://localhost:12345 If you have a browser plugin which
linkifies and displays JSON (such as http://jsonview.com/ or other
more modern replacements) you can click to navigate CREST.

## Gotchas

 * The /universe/ endpoints are not yet available, so map data and ID
 fetching uses munging in a few places.

