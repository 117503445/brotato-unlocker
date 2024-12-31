package process

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	"github.com/117503445/brotato-unlocker/pkg/assets"
	"github.com/rs/zerolog/log"
)

func GetNewJson() string {
	keys := make([]string, 0)
	for _, csvContent := range [][]byte{assets.Csv, assets.CsvDLC1} {
		reader := csv.NewReader(bytes.NewReader(csvContent))
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read file")
			}
			// log.Debug().Str("record", strings.Join(record, ",")).Msg("Record")
			// log.Debug().Str("key", record[0]).Msg("Record")
			key := record[0]
			keys = append(keys, key)
		}
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

	newJsonContent, err := json.Marshal(newJSON)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal new JSON")
	}

	return string(newJsonContent)

}


var NewJson = GetNewJson()