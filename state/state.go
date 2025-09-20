package state

import (
	"os"
	"encoding/json"
)

func FromJSONFile(name string, v any) error {
	content, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, v)
}

func ToJSONFile(name string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return os.WriteFile(name, data, 0644)
}
