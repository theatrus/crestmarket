# crestmarket

[![GoDoc](https://godoc.org/github.com/theatrus/crestmarket?status.svg)](http://godoc.org/github.com/theatrus/crestmarket)
[![Build Status](https://travis-ci.org/theatrus/crestmarket.svg)](https://travis-ci.org/theatrus/crestmarket)

A reference implementation of a full market scraper using the
EVE-Online CREST endpoint.


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
