package encoder

import (
	"reflect"
	"strconv"
	"strings"
)

func encodePrimitive(value interface{}, opts Options) string {
	if value == nil {
		if opts.CompactNull {
			return "~"
		}
		return "null"
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Bool:
		return encodeBool(v.Bool(), opts)
	case reflect.String:
		return encodeString(v.String(), opts.Delimiter)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		return "null"
	}
}

func encodeBool(b bool, opts Options) string {
	if opts.CompactBooleans {
		if b {
			return "1"
		}
		return "0"
	}
	return strconv.FormatBool(b)
}

func encodeString(s string, delimiter string) string {
	if needsQuoting(s, delimiter) {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return s
}

func needsQuoting(s string, delimiter string) bool {
	if s == "" {
		return true
	}

	if strings.ContainsAny(s, ":[]{}") {
		return true
	}

	if delimiter == "\t" {
		if strings.ContainsAny(s, "\t|") {
			return true
		}
	} else {
		if strings.Contains(s, delimiter) || strings.Contains(s, " ") {
			return true
		}
	}

	if s == "true" || s == "false" || s == "null" {
		return true
	}

	return false
}