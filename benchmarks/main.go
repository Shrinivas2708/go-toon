package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Shrinivas2708/go-toon"
	tiktoken "github.com/pkoukk/tiktoken-go"
)

func countTokens(text string) int {
	encoding, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return len(text) / 4
	}
	tokens := encoding.Encode(text, nil, nil)
	return len(tokens)
}

func loadJSON(filename string) (interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	fmt.Println("🚀 TOON Benchmark")
	fmt.Println(strings.Repeat("=", 50))

	files, _ := filepath.Glob("./benchmarks/*.json")

	if len(files) == 0 {
		fmt.Println("❌ No JSON files found in benchmarks directory")
		fmt.Println("   Please add .json files here and run again")
		return
	}

	fmt.Printf("✅ Found %d JSON file(s)\n", len(files))

	totalJSON := 0
	totalTOON := 0

	for _, file := range files {
		data, err := loadJSON(file)
		if err != nil {
			fmt.Printf("⚠️  Error loading %s: %v\n", file, err)
			continue
		}

		name := strings.TrimSuffix(file, ".json")

		jsonCompact, _ := json.Marshal(data)
		jsonPretty, _ := json.MarshalIndent(data, "", "  ")

		opts := toon.EncodeOptions{
			CompactBooleans: true,
			CompactNull:     true,
			Delimiter:       "\t",
		}
		toonOutput := toon.EncodeWithOptions(data, opts)

		jsonTokens := countTokens(string(jsonCompact))
		jsonPrettyTokens := countTokens(string(jsonPretty))
		toonTokens := countTokens(toonOutput)

		savings := float64(jsonTokens-toonTokens) / float64(jsonTokens) * 100

		fmt.Printf("\n📦 %s:\n", name)
		fmt.Printf("   JSON (pretty):  %6d tokens\n", jsonPrettyTokens)
		fmt.Printf("   JSON (compact): %6d tokens\n", jsonTokens)
		fmt.Printf("   TOON:           %6d tokens ✨\n", toonTokens)
		fmt.Printf("   Savings:        %.1f%%\n", savings)

		totalJSON += jsonTokens
		totalTOON += toonTokens
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 Total")
	fmt.Printf("   JSON:    %d tokens\n", totalJSON)
	fmt.Printf("   TOON:    %d tokens\n", totalTOON)
	if totalJSON > 0 {
		fmt.Printf("   Savings: %.1f%%\n", float64(totalJSON-totalTOON)/float64(totalJSON)*100)
	}
}