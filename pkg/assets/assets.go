package assets

import (
	_ "embed"
	"encoding/json"
)

//go:embed init.json
var initJsonBytes []byte

var InitJson map[string]interface{}

func init() {
	err := json.Unmarshal(initJsonBytes, &InitJson)
	if err != nil {
		panic(err)
	}
}

//go:embed translations-v1.1.6.0.csv
var Csv []byte

//go:embed translations-v1.1.6.0-dlc_1.csv
var CsvDLC1 []byte
