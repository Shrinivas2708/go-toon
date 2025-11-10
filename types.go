package toon

type EncodeOptions struct {
	CompactBooleans bool   
	CompactNull     bool   
	Readable        bool   
	Delimiter       string 
	Tabular         bool   
	Flatten         bool   
}

func DefaultOptions() EncodeOptions {
	return EncodeOptions{
		CompactBooleans: false,
		CompactNull:     false,
		Readable:        false,
		Delimiter:       ",",
		Tabular:         true,
		Flatten:         false,
	}
}