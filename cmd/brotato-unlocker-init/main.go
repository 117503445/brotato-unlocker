package main

import (
	"github.com/117503445/brotato-unlocker/pkg/assets"
	"github.com/117503445/brotato-unlocker/pkg/process"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	process.GetNewJson(assets.InitJson)
	err := goutils.WriteJSON("save_v2.json", assets.InitJson)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write JSON")
	}
}
