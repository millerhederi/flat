package flat

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	delimiter = '.'
)

// ErrInvalidKey is an error returned in the `Unflatten` method when one of the keys is invalid and
// can not correctly be parsed.
var ErrInvalidKey = errors.New("invalid key")

// Flatten a nested JSON object, which is expected to either be a `map[string]any“ or `[]any`, into
// a map one level deep.
func Flatten(nested any) (map[string]any, error) {
	result := make(map[string]any)

	flatten(nested, "", result)

	return result, nil
}

// Unflatten a map that is one level deep back into a nested JSON object which will be of type
// `map[string]any` or `[]any`.
func Unflatten(flattened map[string]any) (any, error) {
	if len(flattened) == 0 {
		return map[string]any{}, nil
	}

	var result any

	for key, value := range flattened {
		keyParts, err := splitFlatKey(key)
		if err != nil {
			return nil, err
		}

		result = unflatten(keyParts, value, result)
	}

	return result, nil
}

func flatten(nested any, key string, result map[string]any) {
	switch nested := nested.(type) {
	case map[string]any:
		// handle an empty object
		if len(nested) == 0 && key != "" {
			result[key] = make(map[string]any)
		}

		for childKey, nested := range nested {
			flatten(nested, getFlattenedNestedKey(key, childKey), result)
		}
	case []any:
		// handle an empty slice
		if len(nested) == 0 && key != "" {
			result[key] = make([]any, 0)
		}

		for index, nested := range nested {
			flatten(nested, getFlattenedSliceKey(key, index), result)
		}
	default:
		result[key] = nested
	}
}

func unflatten(keyParts []any, value any, result any) any {
	switch key := keyParts[0].(type) {
	case int:
		return unflattenSliceKey(key, keyParts[1:], value, result)
	case string:
		return unflattenNestedKey(key, keyParts[1:], value, result)
	default:
		panic(fmt.Sprintf("key '%v' was of unexpected type '%T'", keyParts[0], keyParts[0]))
	}
}

// unflattenNestedKey handles unflattening where the current `key` is for a nested object.
func unflattenNestedKey(key string, keyParts []any, value any, result any) any {
	if result == nil {
		result = make(map[string]any)
	}

	resultMap := result.(map[string]any)

	if len(keyParts) == 0 {
		resultMap[key] = value
	} else {
		resultMap[key] = unflatten(keyParts, value, resultMap[key])
	}

	return resultMap
}

// unflattenSliceKey handles unfltattening where the current `key` is for a slice.
func unflattenSliceKey(key int, keyParts []any, value any, result any) any {
	if result == nil {
		result = make([]any, 0)
	}

	resultSlice := expandSliceToSize(result.([]any), key+1)

	if len(keyParts) == 0 {
		resultSlice[key] = value
	} else {
		resultSlice[key] = unflatten(keyParts, value, resultSlice[key])
	}

	return resultSlice
}

func expandSliceToSize(slice []any, size int) []any {
	for i := len(slice); i < size; i++ {
		slice = append(slice, nil)
	}

	return slice
}

// splitFlatKey takes in a single key when unflattening and splits it into individual keys for a
// nested JSON object. The resulting keys will either be a `string` or an `int` for an array index.
func splitFlatKey(key string) ([]any, error) {
	parts := strings.FieldsFunc(key, func(r rune) bool {
		return r == delimiter || r == '['
	})

	keyParts := make([]any, 0, len(parts))
	for _, part := range parts {
		if strings.HasSuffix(part, "]") {
			index, err := strconv.Atoi(strings.TrimSuffix(part, "]"))
			if err != nil {
				return nil, fmt.Errorf("expected part '%s' in key '%s' to be parsable as an integer: %w", part, key, ErrInvalidKey)
			}
			keyParts = append(keyParts, index)
		} else {
			keyParts = append(keyParts, part)
		}
	}

	return keyParts, nil
}

func getFlattenedNestedKey(key, childKey string) string {
	if key == "" {
		return childKey
	}

	return fmt.Sprintf("%s%s%s", key, string(delimiter), childKey)
}

func getFlattenedSliceKey(key string, index int) string {
	return fmt.Sprintf("%s[%d]", key, index)
}
