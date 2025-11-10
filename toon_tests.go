package toon

import (
	"testing"

)

func TestEncodeSimpleArray(t *testing.T) {
	data := map[string]interface{}{
		"tags": []string{"jazz", "chill", "lofi"},
	}
	result := Encode(data)
	expected := "tags[3]jazz,chill,lofi"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeObject(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	result := Encode(data)
	expected := "age:30,name:John"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeWithOptions(t *testing.T) {
	data := map[string]interface{}{
		"active": true,
		"value":  nil,
	}
	opts := EncodeOptions{
		CompactBooleans: true,
		CompactNull:     true,
	}
	result := EncodeWithOptions(data, opts)
	expected := "active:1,value:~"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeNestedObject(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"tags": []string{"admin", "user"},
		},
	}
	result := Encode(data)
	expected := "user{name:John,tags[2]admin,user}"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeEmptyArray(t *testing.T) {
	data := map[string]interface{}{
		"tags": []string{},
	}
	result := Encode(data)
	expected := "tags[0]"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeEmptyObject(t *testing.T) {
	data := map[string]interface{}{
		"config": map[string]interface{}{},
	}
	result := Encode(data)
	expected := "config"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}