package crestmarket

import (
	"encoding/json"
	"time"
)

const (
	producerKey     = "crestmarket"
	producerVersion = "2014-12-01"
	formatVersion   = "0.1"
	dateFormat      = "2006-01-02T15:04:05-07:00"
)

type generator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var gen generator
var columnNames []string

func init() {
	gen.Name = producerKey
	gen.Version = producerVersion
	columnNames = []string{"price", "volRemaining", "range", "orderID", "volEntered", "minVolume", "bid", "issueDate", "duration", "stationID"}
}

type rowset struct {
	GeneratedAt string        `json:"generatedAt"`
	RegionId    int           `json:"regionId"`
	TypeId      int           `json:"typeId"`
	Rows        []interface{} `json:"rows"`
}

type unified struct {
	ResultType  string     `json:"resultType"`
	Generator   *generator `json:"generator"`
	CurrentTime string     `json:"currentTime"`
	Columns     []string   `json:"columns"`
	Rowsets     []*rowset  `json:"rowsets"`
}

// SerializeOrdersUnified outputs a JSON representation
// in the Unified Uploader format convention
// http://dev.eve-central.com/unifieduploader/start
func SerializeOrdersUnified(orders *MarketOrders, at time.Time) ([]byte, error) {
	now := at.UTC().Format(dateFormat)

	out := unified{Rowsets: make([]*rowset, 1),
		Columns: columnNames, Generator: &gen,
		ResultType:  "orders",
		CurrentTime: now}

	rowset := rowset{GeneratedAt: now,
		RegionId: orders.Region.Id,
		TypeId:   orders.Type.Id,
		Rows:     make([]interface{}, len(orders.Orders))}

	out.Rowsets[0] = &rowset

	// Serialize the rows
	for i, order := range orders.Orders {

		row := []interface{}{order.Price,
			order.Volume, order.NumericRange(),
			order.Id, order.Volume, order.MinVolume,
			order.Buy, order.Issued.Format(dateFormat),
			order.Duration, order.Station.Id,
		}
		rowset.Rows[i] = &row
	}

	data, err := json.Marshal(&out)
	return data, err
}
