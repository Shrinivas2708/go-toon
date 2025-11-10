package main

import (
	"fmt"

	"github.com/Shrinivas2708/go-toon"
)

func main() {
	fmt.Println("=== TOON Examples ===\n")

	data1 := map[string]interface{}{
		"tags": []string{"jazz", "chill", "lofi"},
	}
	fmt.Println("1. Simple array:")
	fmt.Println("   Input:  ", data1)
	fmt.Println("   Output: ", toon.Encode(data1))
	fmt.Println()

	data2 := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	fmt.Println("2. Simple object:")
	fmt.Println("   Input:  ", data2)
	fmt.Println("   Output: ", toon.Encode(data2))
	fmt.Println()

	data3 := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"tags": []string{"admin", "user"},
		},
	}
	fmt.Println("3. Nested object:")
	fmt.Println("   Output: ", toon.Encode(data3))
	fmt.Println()

	data4 := map[string]interface{}{
		"active": true,
		"value":  nil,
		"count":  0,
	}
	opts := toon.EncodeOptions{
		CompactBooleans: true,
		CompactNull:     true,
		Readable:        true,
	}
	fmt.Println("4. With compact options:")
	fmt.Println("   Input:  ", data4)
	fmt.Println("   Output: ", toon.EncodeWithOptions(data4, opts))
	fmt.Println()

	data5 := map[string]interface{}{
		"users": []map[string]interface{}{
			{"name": "Alice", "age": 25},
			{"name": "Bob", "age": 30},
		},
	}
	fmt.Println("5. Array of objects (tabular format):")
	fmt.Println("   Output:")
	fmt.Println(toon.Encode(data5))
}