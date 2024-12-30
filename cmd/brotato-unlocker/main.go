package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/117503445/brotato-unlocker/pkg/process"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	appData := os.Getenv("APPDATA")
	if appData == "" {
		log.Fatal().Msg("APPDATA environment variable is not set.")
	}
	// log.Debug().Str("appData", appData).Msg("")

	pattern := filepath.Join(appData, "Brotato", "*", "save_v2.json")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal().Err(err).Msg("Error matching pattern")
		return
	}

	if len(matches) == 0 {
		log.Fatal().Msg("No save_v2.json files found. Please run the game at least once.")
	}

	for _, match := range matches {
		log.Info().Str("match", match).Msg("Found save_v2.json file")
		err := goutils.CopyFile(match, fmt.Sprintf("%s.bak.%s", match, goutils.TimeStrSec()))
		if err != nil {
			log.Warn().Err(err).Msg("Error backing up save_v2.json file, skipping")
			continue
		}

		var jsonContent map[string]interface{}
		err = goutils.ReadJSON(match, &jsonContent)
		if err != nil {
			log.Warn().Err(err).Msg("Error reading save_v2.json file")
			continue
		}

		process.GetNewJson(jsonContent)

		err = goutils.WriteJSON(match, jsonContent)
		if err != nil {
			log.Warn().Err(err).Msg("Error writing save_v2.json file")
			continue
		}

		log.Info().Str("match", match).Msg("Updated save_v2.json file")
	}

	// wait for user input
	fmt.Print("Press 'Enter' to exit...")
	fmt.Scanln()
}
