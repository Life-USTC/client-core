package ustc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// decodeScalarJSON retains the object/array root and scalar values accepted by
// the existing iOS JSONSerialization/SwiftyJSON imports.
func decodeScalarJSON(data []byte) (any, error) {
	if !utf8.Valid(data) || !json.Valid(data) {
		return nil, fmt.Errorf("invalid JSON")
	}
	var value any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	switch value.(type) {
	case map[string]any, []any:
		return value, nil
	default:
		return nil, fmt.Errorf("JSON root must be an object or array")
	}
}

func scalarJSONObject(value any) map[string]any {
	object, _ := value.(map[string]any)
	return object
}

func scalarJSONChildren(value any) []any {
	switch value := value.(type) {
	case []any:
		return value
	case map[string]any:
		items := make([]any, 0, len(value))
		for _, item := range value {
			items = append(items, item)
		}
		return items
	default:
		return nil
	}
}

func scalarJSONText(value any) string {
	switch value := value.(type) {
	case string:
		return value
	case bool:
		return strconv.FormatBool(value)
	case json.Number:
		if integer, err := strconv.ParseInt(value.String(), 10, 64); err == nil {
			return strconv.FormatInt(integer, 10)
		}
		if number, err := value.Float64(); err == nil {
			return strconv.FormatFloat(number, 'g', -1, 64)
		}
	}
	return ""
}
