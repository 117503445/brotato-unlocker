package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	_ "embed"

	"github.com/117503445/goutils"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

//go:embed init.json
var initJSONBytes []byte

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
	map_prefix_keys["challenges_completed"] = append(map_prefix_keys["challenges_completed"], []string{"chal_survivor_1", "chal_survivor_2", "chal_survivor_3", "chal_survivor_4", "chal_survivor_5", "chal_gatherer_1", "chal_gatherer_2", "chal_gatherer_3", "chal_gatherer_4", "chal_gatherer_5", "chal_difficulty_0", "chal_difficulty_1", "chal_difficulty_2", "chal_difficulty_3", "chal_difficulty_4", "chal_difficulty_5"}...)
	for _, character := range map_prefix_keys["characters_unlocked"] {
		characterWithoutPrefix := strings.TrimPrefix(character, "character_")
		// log.Debug().Str("character", character).Str("characterWithoutPrefix", characterWithoutPrefix).Msg("Character")

		map_prefix_keys["challenges_completed"] = append(map_prefix_keys["challenges_completed"], "chal_"+characterWithoutPrefix)
	}

	difficultiesUnlocked := make([]map[string]interface{}, 0)
	for _, character := range map_prefix_keys["characters_unlocked"] {
		zonesDifficultyInfo := []map[string]interface{}{
			{
				"max_selectable_difficulty": 5,
				"zone_id":                   0,
				"max_difficulty_beaten": map[string]interface{}{
					"difficulty_value": 1,
					"enemy_damage":     1,
					"enemy_health":     1,
					"enemy_speed":      1,
					"is_coop":          false,
					"retries":          0,
					"wave_number":      1,
				},
			},
			{
				"max_selectable_difficulty": 5,
				"zone_id":                   1,
				"max_difficulty_beaten": map[string]interface{}{
					"difficulty_value": 1,
					"enemy_damage":     1,
					"enemy_health":     1,
					"enemy_speed":      1,
					"is_coop":          false,
					"retries":          0,
					"wave_number":      1,
				},
			},
		}

		difficultyUnlocked := map[string]interface{}{
			"character_id":          character,
			"zones_difficulty_info": zonesDifficultyInfo,
		}
		difficultiesUnlocked = append(difficultiesUnlocked, difficultyUnlocked)
	}

	newJSON := make(map[string]interface{})
	for key, value := range map_prefix_keys {
		newJSON[key] = value
	}
	newJSON["difficulties_unlocked"] = difficultiesUnlocked
	newJSON["data"] = map[string]interface{}{
		"enemies_killed":      20000,
		"materials_collected": 20000,
	}


	fileSave := "save_v2.json"

	oldJSON := make(map[string]interface{})
	if goutils.FileExists(fileSave) {
		fileBackup := fmt.Sprintf("save_v2.json.bak.%s", goutils.TimeStrSec())
		err := goutils.CopyFile(fileSave, fileBackup)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to copy file")
		}
		log.Info().Str("fileOld", fileSave).Msg("Reading file")
		err = goutils.ReadJSON(fileSave, &oldJSON)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read file")
		}
		// log.Info().Interface("oldJSON", oldJSON).Msg("Old JSON")
	} else {
		log.Info().Msg("Using empty Save file")
		// use initJSONBytes
		err := json.Unmarshal(initJSONBytes, &oldJSON)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to unmarshal initJSONBytes")
		}
		// log.Info().Interface("oldJSON", oldJSON).Msg("Old JSON")
	}
	for key, value := range newJSON {
		oldJSON[key] = value
	}
	err := goutils.WriteJSON(fileSave, oldJSON)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write file")
	}
	log.Info().Str("fileSave", fileSave).Msg("Wrote file")
}

func main() {
	goutils.InitZeroLog()

	// log.Info().Msg("Hello World")
	client := req.C()
	f1 := Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0-dlc_1.csv")
	f2 := Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0.csv")

	Process(f1, f2)
}
