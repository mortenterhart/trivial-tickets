package jsontools

import (
    "encoding/json"
)

func MapToJson(properties map[string]interface{}) ([]byte, error) {
    return json.MarshalIndent(properties, "", "    ")
}
