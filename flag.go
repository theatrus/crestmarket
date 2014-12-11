package crestmarket

import (
	"flag"
)

var isSisi bool
var userAgentSuffix string

func init() {
	flag.BoolVar(&isSisi, "crestmarket.sisi", false, "Call all endpoints on SiSi, turn off for production")
	flag.StringVar(&userAgentSuffix, "crestmarket.useragent", "", "Suffix for the user agent")
}
