package crestmarket

const (
	producerKey     = "crestmarket"
	producerVersion = "2014-12-01"
	formatVersion   = "0.1"
)

type generator struct {
	name    string `json:"name"`
	version string `json:"version"`
}

var gen generator
var columnNames []string

func init() {
	gen.name = producerKey
	gen.version = producerVersion
	columnNames = []string{"price"}
}

type rowset struct {
	generatedAt string        `json:"generatedAt"`
	regionId    int           `json:"regionId"`
	typeId      int           `json:"typeId"`
	rows        []interface{} `json:"rows"`
}

type unified struct {
	resultType  string     `json:"resultType"`
	generator   *generator `json:"generator"`
	currentTime string     `json:"currentTime"`
	columns     []string   `json:"columns"`
	rowsets     []*rowset  `json:"rowsets"`
}
