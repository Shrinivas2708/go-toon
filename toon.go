package toon

import "github.com/Shrinivas2708/go-toon/encoder"

func Encode(value interface{}) string {
	return encoder.Encode(value, encoder.Options(DefaultOptions()))
}

func EncodeWithOptions(value interface{}, opts EncodeOptions) string {
	return encoder.Encode(value, encoder.Options(opts))
}