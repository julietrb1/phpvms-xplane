package models

import (
	"encoding/json"
	"errors"
)

type FlexibleField string

func (f *FlexibleField) UnmarshalJSON(data []byte) error {
	var iface interface{}
	if err := json.Unmarshal(data, &iface); err != nil {
		return err
	}

	switch v := iface.(type) {
	case string:
		// If it's a string, use it
		*f = FlexibleField(v)
		return nil
	case map[string]interface{}: // If it's any object (empty or not), ignore and set to zero value
		*f = "" // Treat as nil-equivalent (empty string)
		return nil
	default:
		return errors.New("field is neither a string nor an object")
	}
}
