package jsontools

import (
	"encoding/json"

	"github.com/mortenterhart/trivial-tickets/structs"
)

func MapToJson(properties structs.JsonMap) ([]byte, error) {
	return json.MarshalIndent(properties, "", "    ")
}
