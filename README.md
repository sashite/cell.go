# cell.go

[![Go Reference](https://pkg.go.dev/badge/github.com/sashite/cell.go.svg)](https://pkg.go.dev/github.com/sashite/cell.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sashite/cell.go)](https://goreportcard.com/report/github.com/sashite/cell.go)
[![License](https://img.shields.io/github/license/sashite/cell.go)](https://github.com/sashite/cell.go/blob/main/LICENSE)

> Idiomatic Go implementation of the **CELL** (Coordinate Encoding for Layered Locations) specification.

## What is CELL?

CELL is a standardized format for representing coordinates on multi-dimensional game boards using a cyclical ASCII character system. It allows for concise, human-readable coordinates that scale to any number of dimensions (e.g., `a1`, `z9`, `aa1`, `a1A`).

This library implements the [CELL Specification v1.0.0](https://sashite.dev/specs/cell/1.0.0/). It is designed for high performance and strict adherence to Go standards (`strconv`-style API).

## Installation

```bash
go get [github.com/sashite/cell.go](https://github.com/sashite/cell.go)

```

## Usage

### 1. Parsing (String to Coordinates)

Use `Parse` to convert a CELL string into a slice of integers.

```go
package main

import (
	"fmt"
	"[github.com/sashite/cell.go](https://github.com/sashite/cell.go)"
)

func main() {
	// Standard parsing
	coords, err := cell.Parse("c3C")
	if err != nil {
		panic(err)
	}
	fmt.Println(coords) // Output: [2, 2, 2]

	// Helper for constants or trusted input (panics on error)
	c := cell.MustParse("a1")
	fmt.Println(c) // Output: [0, 0]
}

```

### 2. Formatting (Coordinates to String)

Use `Format` to convert integer coordinates back to a CELL string. This function accepts a variadic number of arguments, making it easy to use with any dimension.

```go
// 2D Coordinate (e.g., Chess)
s1 := cell.Format(0, 0)
fmt.Println(s1) // Output: "a1"

// 3D Coordinate
s2 := cell.Format(2, 2, 2)
fmt.Println(s2) // Output: "c3C"

// Large coordinates (Extended alphabet)
s3 := cell.Format(26, 0)
fmt.Println(s3) // Output: "aa1"

```

### 3. High Performance (Zero-Allocation)

For game engines and "hot paths" where garbage collection overhead matters, use `AppendParse`. This pattern (similar to `strconv.AppendInt`) allows you to reuse existing memory buffers.

```go
// Pre-allocate a buffer once (capacity of 3 for 3D coords)
buffer := make([]int, 0, 3)

// In your game loop:
inputs := []string{"a1", "b2", "c3"}

for _, input := range inputs {
	// Reset length to 0, keep capacity (no allocation)
	buffer = buffer[:0]

	// Parse directly into the buffer
	var err error
	buffer, err = cell.AppendParse(buffer, input)
	if err != nil {
		continue
	}

	// Process coordinates...
	_ = buffer[0]
}

```

## API Reference

The API mimics the standard library's `strconv` and `strings` packages for familiarity.

```go
// Parse converts a CELL string to a slice of coordinates.
func Parse(s string) ([]int, error)

// Format converts a list of coordinates to a CELL string.
// It is variadic, accepting any number of dimensions.
func Format(indices ...int) string

// MustParse is a helper that panics if parsing fails.
func MustParse(s string) []int

// AppendParse appends the parsed coordinates to dst and returns the extended slice.
// This allows zero-allocation parsing when reusing a buffer.
func AppendParse(dst []int, s string) ([]int, error)

```

## Related Specifications

* [Game Protocol](https://sashite.dev/game-protocol/) — Conceptual foundation
* [CELL Specification](https://sashite.dev/specs/cell/1.0.0/) — Official specification

## License

Available as open source under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0).
