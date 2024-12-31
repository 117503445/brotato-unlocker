package assets

import (
	_ "embed"
)

//go:embed init.json
var InitJson string

//go:embed translations-v1.1.6.0.csv
var Csv []byte

//go:embed translations-v1.1.6.0-dlc_1.csv
var CsvDLC1 []byte
