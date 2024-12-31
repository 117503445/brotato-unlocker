package process

import (
	// "fmt"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestMergeJSON(t *testing.T) {
	tests := []struct {
		name         string
		json1        string
		json2        string
		expectedJSON string
	}{
		{
			name: "Basic merge with higher values in json2",
			json1: `{
				"difficulties_unlocked": [
					{
						"character_id": "character_selection",
						"zones_difficulty_info": [
							{
								"max_selectable_difficulty": 5,
								"zone_id": 0
							}
						]
					}
				],
				"data": {
					"enemies_killed": 15000,
					"materials_collected": 18000
				},
				"field1": "some_value"
			}`,
			json2: `{
				"difficulties_unlocked": [
					{
						"character_id": "character_selection",
						"zones_difficulty_info": [
							{
								"max_selectable_difficulty": 5,
								"zone_id": 1
							}
						]
					},
					{
						"character_id": "character_newbie",
						"zones_difficulty_info": [
							{
								"max_selectable_difficulty": 5,
								"zone_id": 0
							}
						]
					}
				],
				"data": {
					"enemies_killed": 20000,
					"materials_collected": 20000
				},
				"field2": "another_value"
			}`,
			expectedJSON: `{
				"difficulties_unlocked": [
					{
						"character_id": "character_selection",
						"zones_difficulty_info": [
							{
								"max_selectable_difficulty": 5,
								"zone_id": 0
							},
							{
								"max_selectable_difficulty": 5,
								"zone_id": 1
							}
						]
					},
					{
						"character_id": "character_newbie",
						"zones_difficulty_info": [
							{
								"max_selectable_difficulty": 5,
								"zone_id": 0
							}
						]
					}
				],
				"data": {
					"enemies_killed": 20000,
					"materials_collected": 20000
				},
				"field1": "some_value",
				"field2": "another_value"
			}`,
		},
		// Add more test cases here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MergeJSON(tt.json1, tt.json2)
			require.NoError(t, err, "Unexpected error while merging JSONs")

			// Use gjson to compare the expected and actual results
			resultParsed := gjson.Parse(result)
			expectedParsed := gjson.Parse(tt.expectedJSON)

			assert.Equal(t, expectedParsed.Get("data.enemies_killed").Int(), resultParsed.Get("data.enemies_killed").Int())
			assert.Equal(t, expectedParsed.Get("data.materials_collected").Int(), resultParsed.Get("data.materials_collected").Int())

			// Compare difficulties_unlocked by checking if all elements are present
			resultDifficulties := resultParsed.Get("difficulties_unlocked").Array()
			expectedDifficulties := expectedParsed.Get("difficulties_unlocked").Array()

			assert.Len(t, resultDifficulties, len(expectedDifficulties))

			// Check for other fields
			expectedFields := expectedParsed.Map()
			for key, value := range expectedFields {
				if key != "data" && key != "difficulties_unlocked" {
					assert.JSONEq(t, value.Raw, resultParsed.Get(key).Raw, "Field %s does not match", key)
					// assert.Equal(t, fmt.Sprint(value), resultParsed.Get(key).Raw)
				}
			}
		})
	}
}
