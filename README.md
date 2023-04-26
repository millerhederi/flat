# flat

A go library for flattening and unflattening JSON.

## Flatten

Flatten any JSON into a map one level deep:
```go
nested := map[string]any{
    "a": map[string]any{
        "b": 1,
        "c": nil,
        "d": []any{false, true},
    },
    "e": "f",
    "g": 2.3,
}

flattened, err := flat.Flatten(nested)

// flattened = map[string]any{
//     "a.b":    1,
//     "a.c":    nil,
//     "a.d[0]": false,
//     "a.d[1]": true,
//     "e":      "f",
//     "g":      2.3,
// }
```

Note that Flatten also supports JSON where the top level is an array, e.g. `[ "a", "b" ]`, which
would produce a flattened result `map[string]any{"[0]":"a", "[1]": "b"}`.

## Unflatten

This is the inverse of Flatten. Unflattens a map one level deep back into nested JSON:
```go
flattened := map[string]any{
    "a.b":    1,
    "a.c":    nil,
    "a.d[0]": false,
    "a.d[1]": true,
    "e":      "f",
    "g":      2.3,
}

nested, err := flat.Unflatten(flattened)

// nested = map[string]any{
//     "a": map[string]any{
//         "b": 1,
//         "c": nil,
//         "d": []any{false, true},
//     },
//     "e": "f",
//     "g": 2.3,
// }
```
