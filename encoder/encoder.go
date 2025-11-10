package encoder

import (
	"reflect"
)

func Encode(value interface{}, opts Options) string {
	if value == nil {
		if opts.CompactNull {
			return "~"
		}
		return "null"
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return encodeArray(v, opts, "")
	case reflect.Map:
		return encodeMap(v, opts)
	default:
		return encodePrimitive(value, opts)
	}
}