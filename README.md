# 🚀 go-toon

**TOON** (Token-Optimized Object Notation) is a Go library for serializing data structures into a highly compact, human-readable string format.

Its primary design goal is to be **token-efficient** for language model (LLM) contexts, providing a more compact representation than JSON.

---

## 🌟 Why TOON?

JSON is verbose. `go-toon ` converts complex data into a minimal format, dramatically reducing token count.

**JSON (compact):**
```json
{"library":{"books":[{"author":{"firstName":"F. Scott","lastName":"Fitzgerald"},"id":1,"title":"The Great Gatsby"}],"location":{"city":"San Francisco"},"name":"Central Public Library"}}
````

*(181 characters, 50 tokens)*

**TOON (with options):**

```
library{books[1]author{firstName:"F. Scott",lastName:Fitzgerald}id:1,title:"The Great Gatsby",location{city:"San Francisco"},name:"Central Public Library"}
```

*(148 characters, 40 tokens)*

---

## ⚙️ Features

* **Token-Efficient** → Designed to use fewer tokens than JSON.
* **Tabular Encoding** → Automatically formats arrays of uniform objects into a compact table — perfect for LLM prompts.
* **Compact Options** → Further reduce output size by compacting booleans (`true → 1`) and nulls (`nil → ~`).
* **Data Flattening** → Option to flatten nested JSON structures into `key_name` format.

---

## 📦 Installation

```bash
go get github.com/Shrinivas2708/toon
```

---

## 💡 Usage

There are two main functions:

* `Encode()` — for default settings.
* `EncodeWithOptions()` — for customized encoding.

### 🧩 Basic Example

```go
package main

import (
	"fmt"
	"github.com/Shrinivas2708/toon"
)

func main() {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"tags": []string{"admin", "user"},
		},
	}

	encoded := toon.Encode(data)
	fmt.Println(encoded)
	// Output: user{name:John,tags[2]admin,user}
}
```

### ⚙️ Example with Options

```go
package main

import (
	"fmt"
	"github.com/Shrinivas2708/toon"
)

func main() {
	data := map[string]interface{}{
		"active": true,
		"value":  nil,
		"count":  0,
	}

	opts := toon.EncodeOptions{
		CompactBooleans: true,
		CompactNull:     true,
		Readable:        true,
	}

	encoded := toon.EncodeWithOptions(data, opts)
	fmt.Println(encoded)
	// Output: active: 1, count: 0, value: ~
}
```

---

## 🧰 Encoding Options

You can customize the encoding process by passing an `EncodeOptions` struct:

```go
type EncodeOptions struct {
	CompactBooleans bool   // Encode true/false as 1/0
	CompactNull     bool   // Encode nil as ~
	Readable        bool   // Add spaces for readability
	Delimiter       string // Separator for tabular data (default ",")
	Tabular         bool   // Enable tabular encoding for object arrays (default true)
	Flatten         bool   // Flatten nested objects
}
```

---

## 🧠 Format Examples

| Input Data                                      | Output (`toon.Encode(data)`)   |
| ----------------------------------------------- | ------------------------------ |
| `{"name": "John", "age": 30}`                   | `age:30,name:John`             |
| `{"tags": ["jazz", "chill", "lofi"]}`           | `tags[3]jazz,chill,lofi`       |
| `{"tags": []}`                                  | `tags[0]`                      |
| `{"config": {}}`                                | `config`                       |
| `{"user": {"name": "John", "tags": ["admin"]}}` | `user{name:John,tags[1]admin}` |

---

### 🧾 Tabular Format Example

When encoding arrays of objects, `go-toon` defaults to a **token-saving tabular format**.

**Input:**

```go
data := map[string]interface{}{
    "users": []map[string]interface{}{
        {"name": "Alice", "age": 25},
        {"name": "Bob", "age": 30},
    },
}
```

**Output (`toon.Encode(data)`):**

```
users
age,name
25,Alice
30,Bob
```

---

## 📉 Token Benchmark

The project includes a benchmark comparing `tiktoken` (cl100k_base) token counts for TOON vs. JSON.

```bash
$ go run ./benchmarks
🚀 TOON Benchmark
==================================================
✅ Found 2 JSON file(s)
   -> Wrote benchmarks/github_repo.toon
   -> Wrote benchmarks/github_repo.compact.json

📦 benchmarks/github_repo:
   JSON (pretty):   14389 tokens
   JSON (compact):  13615 tokens
   TOON:            10100 tokens ✨
   Savings:        25.8%
   -> Wrote benchmarks/books.toon
   -> Wrote benchmarks/books.compact.json

📦 benchmarks/books:
   JSON (pretty):    1420 tokens
   JSON (compact):   1138 tokens
   TOON:             1009 tokens ✨
   Savings:        11.3%

==================================================
📊 Total
   JSON:    14753 tokens
   TOON:    11109 tokens
   Savings: 24.7%
```


