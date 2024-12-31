package main

import (
	"github.com/117503445/brotato-unlocker/pkg/assets"
	"github.com/117503445/brotato-unlocker/pkg/process"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	mergedJson, err := process.MergeJSON(assets.InitJson, process.NewJson)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to merge JSON")
	}

	err = goutils.WriteText("save_v2.json", mergedJson)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write JSON")
	}
}
