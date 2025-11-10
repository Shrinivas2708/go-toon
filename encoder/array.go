package encoder

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func encodeArray(v reflect.Value, opts Options, keyName string) string {
	length := v.Len()
	if length == 0 {
		return "[0]"
	}

	if opts.Flatten {
		firstElem := v.Index(0)
		if firstElem.Kind() == reflect.Interface {
			firstElem = firstElem.Elem()
		}
		if firstElem.IsValid() && firstElem.Kind() == reflect.Map {
			allKeysSet := make(map[string]bool)
			flattenedData := make([]map[string]interface{}, length)

			for i := 0; i < length; i++ {
				item := v.Index(i).Interface() 
				flattenedMap := flattenObject(item, "")
				flattenedData[i] = flattenedMap
				for k := range flattenedMap {
					allKeysSet[k] = true
				}
			}

			allKeys := make([]string, 0, len(allKeysSet))
			for k := range allKeysSet {
				allKeys = append(allKeys, k)
			}
			sort.Strings(allKeys)

			uniformData := make([]map[string]interface{}, length)
			for i, flatMap := range flattenedData {
				newItem := make(map[string]interface{})
				for _, key := range allKeys {
					if val, ok := flatMap[key]; ok {
						newItem[key] = val
					} else {
						newItem[key] = nil
					}
				}
				uniformData[i] = newItem
			}
			uniformValue := reflect.ValueOf(uniformData)
			return encodeTabularArray(uniformValue, opts)
		}
	}

	if opts.Tabular && isUniformObjectArray(v) {
		return encodeTabularArray(v, opts)
	}

	items := make([]string, length)
	for i := 0; i < length; i++ {
		var item interface{}
		val := v.Index(i)
		if val.IsValid() {
			item = val.Interface()
		}

		encoded := Encode(item, opts)

		if item != nil {
			itemType := reflect.TypeOf(item)
			if itemType != nil {
				kind := itemType.Kind()
				if kind == reflect.Ptr {
					kind = itemType.Elem().Kind()
				}
				if kind == reflect.Map {
					encoded = "{" + encoded + "}"
				}
			}
		}
		items[i] = encoded
	}

	separator := ","
	if opts.Readable {
		separator = ", "
	}

	spaceAfterCount := ""
	if opts.Readable {
		spaceAfterCount = " "
	}

	return fmt.Sprintf("[%d]%s%s", length, spaceAfterCount, strings.Join(items, separator))
}

func isUniformObjectArray(v reflect.Value) bool {
	if v.Len() == 0 {
		return false
	}

	var firstMapKind reflect.Kind
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		if !elem.IsValid() || elem.Kind() != reflect.Map {
			return false
		}
		if i == 0 {
			firstMapKind = elem.Kind()
		}
		if elem.Kind() != firstMapKind {
			return false
		}
	}

	firstMapVal := v.Index(0)
	if firstMapVal.Kind() == reflect.Interface {
		firstMapVal = firstMapVal.Elem()
	}
	if firstMapVal.Len() == 0 {
		return false
	}

	keys := firstMapVal.MapKeys()
	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[fmt.Sprintf("%v", k.Interface())] = true
	}

	for i := 1; i < v.Len(); i++ {
		itemMap := v.Index(i)
		if itemMap.Kind() == reflect.Interface {
			itemMap = itemMap.Elem()
		}
		itemKeys := itemMap.MapKeys()
		if len(itemKeys) != len(keys) {
			return false
		}
		for _, k := range itemKeys {
			if !keySet[fmt.Sprintf("%v", k.Interface())] {
				return false
			}
		}
	}

	for i := 0; i < v.Len(); i++ {
		itemMap := v.Index(i)
		if itemMap.Kind() == reflect.Interface {
			itemMap = itemMap.Elem()
		}
		for _, k := range keys {
			val := itemMap.MapIndex(k)
			if !isPrimitive(val) { 
				return false
			}
		}
	}

	return true
}

func isPrimitive(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	elem := v
	if elem.Kind() == reflect.Interface {
		elem = elem.Elem()
	}

	if !elem.IsValid() {
		return true 
	}

	kind := elem.Kind()
	switch kind {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func encodeTabularArray(v reflect.Value, opts Options) string {
	if v.Len() == 0 {
		return "[0]"
	}

	firstMapVal := v.Index(0)
	var firstMap reflect.Value
	if firstMapVal.Kind() == reflect.Interface {
		firstMap = firstMapVal.Elem()
	} else {
		firstMap = firstMapVal
	}

	if firstMap.Kind() != reflect.Map {
		return "[]"
	}

	keys := firstMap.MapKeys()
	keyStrs := make([]string, len(keys))
	keyMap := make(map[string]reflect.Value, len(keys))
	for i, k := range keys {
		keyStrs[i] = fmt.Sprintf("%v", k.Interface())
		keyMap[keyStrs[i]] = k
	}
	sort.Strings(keyStrs)

	delimiter := opts.Delimiter
	if delimiter == "" {
		delimiter = ","
	}

	headerItems := make([]string, len(keyStrs))
	for i, keyStr := range keyStrs {
		headerItems[i] = encodeString(keyStr, delimiter)
	}
	header := strings.Join(headerItems, delimiter)

	rows := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		itemMapVal := v.Index(i)
		var itemMap reflect.Value
		if itemMapVal.Kind() == reflect.Interface {
			itemMap = itemMapVal.Elem()
		} else {
			itemMap = itemMapVal
		}

		values := make([]string, len(keyStrs))
		for j, keyStr := range keyStrs {
			key := keyMap[keyStr]
			var val reflect.Value
			if key.IsValid() {
				val = itemMap.MapIndex(key)
			}

			var elemVal reflect.Value
			if val.IsValid() && val.Kind() == reflect.Interface {
				elemVal = val.Elem()
			} else {
				elemVal = val
			}

			var encoded string
			if elemVal.IsValid() {
				encoded = encodePrimitive(elemVal.Interface(), opts)
			} else {
				encoded = encodePrimitive(nil, opts)
			}

			if opts.Delimiter == "\t" && len(encoded) > 1 && strings.HasPrefix(encoded, `"`) && strings.HasSuffix(encoded, `"`) {
				unquoted := encoded[1 : len(encoded)-1]
				if !strings.ContainsAny(unquoted, "\t\n\"") {
					encoded = unquoted
				}
			}
			values[j] = encoded
		}
		rows[i] = strings.Join(values, delimiter)
	}

	return header + "\n" + strings.Join(rows, "\n")
}