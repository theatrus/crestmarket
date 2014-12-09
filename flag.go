package crestmarket

import (
	"flag"
)

var isSisi bool

func init() {
	flag.BoolVar(&isSisi, "crestmarket.sisi", false, "Call all endpoints on SiSi, turn off for production")
}
