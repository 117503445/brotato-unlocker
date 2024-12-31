package process

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func MergeJSON(json1, json2 string) (string, error) {

	var err error

	// Step 1: Parse the JSON strings using gjson
	json1Parsed := gjson.Parse(json1)
	json2Parsed := gjson.Parse(json2)

	// Step 2: Merge data fields with comparison
	data1 := json1Parsed.Get("data")
	data2 := json2Parsed.Get("data")

	var enemiesKilled, materialsCollected int64
	if data1.Exists() && data2.Exists() {
		enemiesKilled = max(data1.Get("enemies_killed").Int(), data2.Get("enemies_killed").Int())
		materialsCollected = max(data1.Get("materials_collected").Int(), data2.Get("materials_collected").Int())

		// Update json1 with new values if necessary
		if data1.Get("enemies_killed").Int() < enemiesKilled {
			json1, err = sjson.Set(json1, "data.enemies_killed", enemiesKilled)
			if err != nil {
				return "", err
			}
		}
		if data1.Get("materials_collected").Int() < materialsCollected {
			json1, err = sjson.Set(json1, "data.materials_collected", materialsCollected)
			if err != nil {
				return "", err
			}
		}
	}

	// Step 3: Process difficulties_unlocked field
	difficultiesUnlocked1 := json1Parsed.Get("difficulties_unlocked").Array()
	difficultiesUnlocked2 := json2Parsed.Get("difficulties_unlocked").Array()

	difficultiesMap := make(map[string]interface{})

	for _, du := range difficultiesUnlocked1 {
		characterID := du.Get("character_id").String()
		difficultiesMap[characterID] = du.Value()
	}

	for _, du := range difficultiesUnlocked2 {
		characterID := du.Get("character_id").String()
		zonesDifficultyInfo := du.Get("zones_difficulty_info").Array()

		existingCharacterData := difficultiesMap[characterID]
		if existingCharacterData == nil {
			// Add new character data from json2 to json1
			json1, _ = sjson.Set(json1, "difficulties_unlocked.-1", du.Value())
		} else {
			// Update or add zones_difficulty_info for an existing character
			for _, zdi := range zonesDifficultyInfo {
				zoneID := zdi.Get("zone_id").Int()
				existingZones := gjson.Parse(fmt.Sprint(existingCharacterData)).Get("zones_difficulty_info").Array()
				found := false
				for _, ez := range existingZones {
					if ez.Get("zone_id").Int() == zoneID {
						// Update existing zone's max_selectable_difficulty
						maxSelectableDifficulty := zdi.Get("max_selectable_difficulty").Int()
						json1, _ = sjson.Set(json1, fmt.Sprintf("difficulties_unlocked.#(character_id==%q).zones_difficulty_info.#(zone_id==%d).max_selectable_difficulty", characterID, zoneID), maxSelectableDifficulty)
						found = true
						break
					}
				}
				if !found {
					// Add new zone's data
					json1, _ = sjson.Set(json1, fmt.Sprintf("difficulties_unlocked.#(character_id==%q).zones_difficulty_info.-1", characterID), zdi.Value())
				}
			}
		}
	}

	// Step 4: Add other fields from json2 to json1
	json2Fields := json2Parsed.Map()
	for key, value := range json2Fields {
		if key != "data" && key != "difficulties_unlocked" {
			// log.Debug().Str("key", key).Str("value", value.Raw).Msg("Adding field to json1")
			json1, err = sjson.Set(json1, key, value.Value())
			if err != nil {
				return "", err
			}
		}
	}

	// log.Debug().Str("json1", json1).Msg("Merged JSON")

	return json1, nil
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
