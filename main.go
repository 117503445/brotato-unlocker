package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/117503445/goutils"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

func Download(client *req.Client, url string) string {
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	if goutils.FileExists(fileName) {
		return fileName
	}

	log.Info().Str("url", url).Str("fileName", fileName).Msg("Downloading file")

	resp, err := client.R().
		SetOutputFile(fileName).
		Get(url)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to download file")
	}
	defer resp.Body.Close()
	return fileName
}

func Process(paths ...string) {
	keys := make([]string, 0)
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open file")
			return
		}

		reader := csv.NewReader(file)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read file")
				return
			}
			// log.Debug().Str("record", strings.Join(record, ",")).Msg("Record")
			// log.Debug().Str("key", record[0]).Msg("Record")
			key := record[0]
			keys = append(keys, key)
		}
		file.Close()
	}

	map_prefix_key := map[string]string{
		"CHARACTER_":  "characters_unlocked",
		"CONSUMABLE_": "consumables_unlocked",
		"ITEM_":       "items_unlocked",
		"UPGRADE_":    "upgrades_unlocked",
		"CHAL_":       "challenges_completed",
		"WEAPON_":     "weapons_unlocked",
	}

	map_prefix_keys := make(map[string][]string)
	for _, key := range keys {
		for prefix, newKey := range map_prefix_key {
			if strings.HasPrefix(key, prefix) {
				item := strings.ToLower(key)
				// log.Info().Str("item", item).Msg("Item key")
				map_prefix_keys[newKey] = append(map_prefix_keys[newKey], item)
				break
			}
		}

	}

	fileSave := "out.json"

	oldJSON := make(map[string]interface{})
	if goutils.FileExists("save_v2.json") {
		err := goutils.ReadJSON("save_v2.json", &oldJSON)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read file")
		}
		log.Info().Interface("oldJSON", oldJSON).Msg("Old JSON")
		for key, value := range map_prefix_keys {
			oldJSON[key] = value
		}
		err = goutils.WriteJSON(fileSave, oldJSON)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to write file")
		}
	} else {
		err := goutils.WriteJSON(fileSave, map_prefix_keys)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to write file")
		}
	}
}

func main() {
	goutils.InitZeroLog()

	log.Info().Msg("Hello World")
	client := req.C()
	f1 := Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0-dlc_1.csv")
	f2 := Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0.csv")

	Process(f1, f2)
}
