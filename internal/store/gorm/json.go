package gorm

import "encoding/json"

// marshalJSON serializes value to a JSON string.
func marshalJSON[T any](value T) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// unmarshalJSON deserializes a JSON string into a value of type T.
// Returns the zero value of T when jsonStr is empty.
func unmarshalJSON[T any](jsonStr string) (T, error) {
	var result T
	if jsonStr == "" {
		return result, nil
	}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}
