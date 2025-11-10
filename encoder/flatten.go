package encoder

import (
	"fmt"
	"reflect"
	"sort"
)

// flattenObject converts a nested structure (map or slice) into a flat map[string]interface{}
// with keys joined by underscores.
// Example: { "user": { "name": "John" } } -> { "user_name": "John" }
// Example: { "tags": ["a", "b"] } -> { "tags0": "a", "tags1": "b" }
func flattenObject(data interface{}, prefix string) map[string]interface{} {
	flattened := make(map[string]interface{})
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !v.IsValid() {
		if prefix != "" {
			flattened[prefix] = nil
		}
		return flattened
	}

	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()

		sortedKeys := make([]string, len(keys))
		keyMap := make(map[string]reflect.Value, len(keys))
		for i, k := range keys {
			keyStr := fmt.Sprintf("%v", k.Interface())
			sortedKeys[i] = keyStr
			keyMap[keyStr] = k
		}
		sort.Strings(sortedKeys)

		if len(sortedKeys) == 0 && prefix != "" {
			flattened[prefix] = nil
		}

		for _, keyStr := range sortedKeys {
			val := v.MapIndex(keyMap[keyStr])
			newKey := keyStr
			if prefix != "" {
				newKey = prefix + "_" + keyStr
			}
			nested := flattenObject(val.Interface(), newKey)
			for k, v := range nested {
				flattened[k] = v
			}
		}

	case reflect.Slice, reflect.Array:
		length := v.Len()
		if length == 0 {
			if prefix != "" {
				flattened[prefix] = nil
			}
		}

		for i := 0; i < length; i++ {
			newKey := fmt.Sprintf("%s%d", prefix, i)
			val := v.Index(i)
			nested := flattenObject(val.Interface(), newKey)
			for k, v := range nested {
				flattened[k] = v
			}
		}

	default:
		if prefix == "" {
			flattened["value"] = v.Interface()
		} else {
			flattened[prefix] = v.Interface()
		}
	}

	return flattened
}