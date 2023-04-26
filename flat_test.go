package flat_test

import (
	"encoding/json"
	"errors"
	"flat"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected map[string]any
	}{
		"complex object": {
			input: `{
				"a": {
					"b": 1,
					"c": null,
					"d": [false, true]
				},
				"e": "f",
				"g": 2.3
			}`,
			expected: map[string]any{
				"a.b":    1.0,
				"a.c":    nil,
				"a.d[0]": false,
				"a.d[1]": true,
				"e":      "f",
				"g":      2.3,
			},
		},
		"empty object": {
			input:    `{}`,
			expected: map[string]any{},
		},
		"empty nested object": {
			input: `{"a": {}}`,
			expected: map[string]any{
				"a": map[string]any{},
			},
		},
		"empty array": {
			input: `{"a": []}`,
			expected: map[string]any{
				"a": []any{},
			},
		},
		"root array": {
			input: `[
				{"a": 1},
				{"a": 2}
			]`,
			expected: map[string]any{
				"[0].a": 1.0,
				"[1].a": 2.0,
			},
		},
		"empty root array": {
			input:    `[]`,
			expected: map[string]any{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var input any
			if err := json.Unmarshal([]byte(test.input), &input); err != nil {
				t.Fatalf("failed unmarshaling the test input json: %v", err)
			}

			actual, err := flat.Flatten(input)
			if err != nil {
				t.Fatalf("failed flattening the test input: %v", err)
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestUnflatten(t *testing.T) {
	tests := map[string]struct {
		input    map[string]any
		expected string
	}{
		"complext object": {
			input: map[string]any{
				"a.b":    1.0,
				"a.c":    nil,
				"a.d[0]": false,
				"a.d[1]": true,
				"e":      "f",
				"g":      2.3,
			},
			expected: `{
				"a": {
					"b": 1,
					"c": null,
					"d": [false, true]
				},
				"e": "f",
				"g": 2.3
			}`,
		},
		"empty object": {
			input:    map[string]any{},
			expected: `{}`,
		},
		"empty nested object": {
			input: map[string]any{
				"a": map[string]any{},
			},
			expected: `{"a": {}}`,
		},
		"root array": {
			input: map[string]any{
				"[0].a": 1.0,
				"[1].a": 2.0,
			},
			expected: `[
				{"a": 1},
				{"a": 2}
			]`,
		},
		"empty array": {
			input: map[string]any{
				"a": []any{},
			},
			expected: `{"a": []}`,
		},
		"array with missing indices": {
			input: map[string]any{
				"[2].a": 1.0,
				"[4].b": "a",
			},
			expected: `[
				null,
				null,
				{"a": 1},
				null,
				{"b": "a"}
			]`,
		},
		"complex value": {
			input: map[string]any{
				"a": map[string]any{"b": 1.0},
			},
			expected: `{"a": {"b": 1}}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var expected any
			if err := json.Unmarshal([]byte(test.expected), &expected); err != nil {
				t.Fatalf("failed unmarshaling the test expected json: %v", err)
			}

			actual, err := flat.Unflatten(test.input)
			if err != nil {
				t.Fatalf("failed unflattening the test input: %v", err)
			}

			assert.Equal(t, expected, actual)
		})
	}
}

func TestUnflattenErrorOnInvalidKey(t *testing.T) {
	input := map[string]any{
		"a[not_an_int].b": true,
	}

	actual, err := flat.Unflatten(input)
	assert.Nil(t, actual)
	assert.True(t, errors.Is(err, flat.ErrInvalidKey))
}

func TestReadme(t *testing.T) {
	input := map[string]any{
		"a": map[string]any{
			"b": 1,
			"c": nil,
			"d": []any{false, true},
		},
		"e": "f",
		"g": 2.3,
	}

	flattened, _ := flat.Flatten(input)

	expected := map[string]any{
		"a.b":    1,
		"a.c":    nil,
		"a.d[0]": false,
		"a.d[1]": true,
		"e":      "f",
		"g":      2.3,
	}

	assert.Equal(t, expected, flattened)
}
