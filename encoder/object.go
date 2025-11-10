package encoder

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func encodeMap(v reflect.Value, opts Options) string {
	if v.Len() == 0 {
		return ""
	}

	keys := v.MapKeys()
	sortedKeys := make([]string, len(keys))
	keyMap := make(map[string]reflect.Value)

	for i, k := range keys {
		keyStr := fmt.Sprintf("%v", k.Interface())
		sortedKeys[i] = keyStr
		keyMap[keyStr] = k
	}
	sort.Strings(sortedKeys)

	pairs := make([]string, 0, len(sortedKeys))

	for _, keyStr := range sortedKeys {
		key := keyMap[keyStr]
		val := v.MapIndex(key)

		if !val.IsValid() {
			continue
		}

		elem := val
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}

		if !elem.IsValid() {
			encodedKey := encodeString(keyStr, opts.Delimiter)
			encodedVal := encodePrimitive(nil, opts) 
			space := ""
			if opts.Readable {
				space = " "
			}
			pairs = append(pairs, encodedKey+":"+space+encodedVal)
			continue
		}

		encodedKey := encodeString(keyStr, opts.Delimiter)

		switch elem.Kind() {
		case reflect.Slice, reflect.Array:
			arrStr := encodeArray(elem, opts, keyStr)

			isTabular := strings.Contains(arrStr, "\n") && !strings.HasPrefix(arrStr, "[")

			if isTabular {
				if v.Len() == 1 {
					return arrStr
				}
				pairs = append(pairs, encodedKey+"\n"+arrStr)
			} else {
				pairs = append(pairs, encodedKey+arrStr)
			}

		case reflect.Map:
			nested := encodeMap(elem, opts)
			if nested == "" {
				pairs = append(pairs, encodedKey)
			} else {
				pairs = append(pairs, encodedKey+"{"+nested+"}")
			}

		default:
			encodedVal := encodePrimitive(elem.Interface(), opts)
			space := ""
			if opts.Readable {
				space = " "
			}
			pairs = append(pairs, encodedKey+":"+space+encodedVal)
		}
	}

	sep := ","
	if opts.Readable {
		sep = ", "
	}

	return strings.Join(pairs, sep)
}