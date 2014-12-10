package crestmarket

import (
	"bytes"
	"log"
	"testing"
	"time"
)

var expectedSerialized = []byte(`{"resultType":"orders","generator":{"name":"crestmarket","version":"2014-12-01"},"currentTime":"2009-11-10T23:00:00+00:00","columns":["price","volRemaining","range","orderID","volEntered","minVolume","bid","issueDate","duration","stationID"],"rowsets":[{"generatedAt":"2009-11-10T23:00:00+00:00","regionID":1000002,"typeID":34,"rows":[[3,3000,0,999999,3000,1,false,"2009-11-10T23:00:00+00:00",90,8888]]}]}`)

func TestSerializeUnified(t *testing.T) {
	then := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	mtype := &MarketType{"Tritanium", "http://nope", 34}
	mregion := &Region{"The Forge", "http://nope", 1000002}

	orders := NewMarketOrders()
	orders.Fetched = then
	orders.Type = mtype
	orders.Region = mregion

	order := MarketOrder{Buy: false, Duration: 90, Href: "", Id: 999999, Issued: then, Station: Station{"hmm", "hmm", 8888},
		MinVolume: 1, Price: 3.0, Range: "solarsystem", Type: *mtype, Volume: 3000}
	orders.Orders = append(orders.Orders, &order)

	data, err := SerializeOrdersUnified(orders, then)
	log.Printf("%s\n", data)
	if err != nil {
		t.Error(err)
	}
	if bytes.Equal(data, expectedSerialized) == false {
		t.Error("Data was not equal")
	}
}
