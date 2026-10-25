package ustc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

func scalarJSONNumber(value any) *float64 {
	var number float64
	switch value := value.(type) {
	case json.Number:
		var err error
		number, err = value.Float64()
		if err != nil {
			return nil
		}
	case bool:
		if value {
			number = 1
		}
	default:
		return nil
	}
	return &number
}

func scalarJSONInt(value any) *int {
	if number := scalarJSONNumber(value); number != nil {
		integer := int(*number)
		return &integer
	}
	return nil
}

func scalarJSONString(value any) *string {
	if text, ok := value.(string); ok {
		return &text
	}
	return nil
}

// SwiftyJSON's numberValue uses NSDecimalNumber for strings, unlike the
// optional number accessor. The native scanner accepts a leading decimal
// prefix and whitespace; invalid strings have the existing zero default.
var scalarDecimalPrefix = regexp.MustCompile(`^[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?`)

func scalarJSONNumberValue(value any) float64 {
	if text, ok := value.(string); ok {
		prefix := scalarDecimalPrefix.FindString(strings.TrimLeftFunc(text, unicode.IsSpace))
		number, err := strconv.ParseFloat(prefix, 64)
		if err != nil {
			return 0
		}
		return number
	}
	if number := scalarJSONNumber(value); number != nil {
		return *number
	}
	return 0
}

func scalarJSONIntValue(value any) int {
	return int(scalarJSONNumberValue(value))
}
