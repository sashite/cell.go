# cell.go

[![Go Reference](https://pkg.go.dev/badge/github.com/sashite/cell.go.svg)](https://pkg.go.dev/github.com/sashite/cell.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sashite/cell.go)](https://goreportcard.com/report/github.com/sashite/cell.go)
[![License](https://img.shields.io/github/license/sashite/cell.go)](https://github.com/sashite/cell.go/blob/main/LICENSE)

> **CELL** (Coordinate Encoding for Layered Locations) implementation for Go.

## What is CELL?

CELL (Coordinate Encoding for Layered Locations) is a standardized format for representing coordinates on multi-dimensional game boards using a cyclical ASCII character system. CELL supports unlimited dimensional coordinate systems through the systematic repetition of three distinct character sets.

This library implements the [CELL Specification v1.0.0](https://sashite.dev/specs/cell/1.0.0/).

## Installation

```bash
go get github.com/sashite/cell.go
```

## CELL Format

CELL uses a cyclical three-character-set system that repeats indefinitely based on dimensional position:

| Dimension | Condition | Character Set | Examples |
|-----------|-----------|---------------|----------|
| 1st, 4th, 7th… | n % 3 = 1 | Latin lowercase (`a`–`z`) | `a`, `e`, `aa`, `file` |
| 2nd, 5th, 8th… | n % 3 = 2 | Positive integers | `1`, `8`, `10`, `256` |
| 3rd, 6th, 9th… | n % 3 = 0 | Latin uppercase (`A`–`Z`) | `A`, `C`, `AA`, `LAYER` |

## Usage

```go
package main

import (
    "fmt"
    "github.com/sashite/cell.go/cell"
)

func main() {
    // Validation
    fmt.Println(cell.Valid("a1"))       // true (2D coordinate)
    fmt.Println(cell.Valid("a1A"))      // true (3D coordinate)
    fmt.Println(cell.Valid("e4"))       // true (2D coordinate)
    fmt.Println(cell.Valid("h8Hh8"))    // true (5D coordinate)
    fmt.Println(cell.Valid("*"))        // false (not a CELL coordinate)
    fmt.Println(cell.Valid("a0"))       // false (invalid numeral)
    fmt.Println(cell.Valid(""))         // false (empty string)

    // Dimensional analysis
    fmt.Println(cell.Dimensions("a1"))     // 2
    fmt.Println(cell.Dimensions("a1A"))    // 3
    fmt.Println(cell.Dimensions("h8Hh8"))  // 5
    fmt.Println(cell.Dimensions("foobar")) // 1

    // Parse coordinate into dimensional components
    components, err := cell.Parse("a1A")
    if err == nil {
        fmt.Println(components) // ["a", "1", "A"]
    }

    components, err = cell.Parse("h8Hh8")
    if err == nil {
        fmt.Println(components) // ["h", "8", "H", "h", "8"]
    }

    components, err = cell.Parse("foobar")
    if err == nil {
        fmt.Println(components) // ["foobar"]
    }

    _, err = cell.Parse("1nvalid")
    fmt.Println(err) // error: invalid CELL coordinate: 1nvalid

    // Must-style parsing (panics on error)
    components = cell.MustParse("a1A") // ["a", "1", "A"]

    // Convert coordinates to 0-indexed integer slices
    indices, err := cell.ToIndices("a1")
    if err == nil {
        fmt.Println(indices) // [0, 0]
    }

    indices, err = cell.ToIndices("e4")
    if err == nil {
        fmt.Println(indices) // [4, 3]
    }

    indices, err = cell.ToIndices("a1A")
    if err == nil {
        fmt.Println(indices) // [0, 0, 0]
    }

    indices, err = cell.ToIndices("b2B")
    if err == nil {
        fmt.Println(indices) // [1, 1, 1]
    }

    // Must-style conversion (panics on error)
    indices = cell.MustToIndices("e4") // [4, 3]

    // Convert 0-indexed integers to CELL coordinates
    coord, err := cell.FromIndices(0, 0)
    if err == nil {
        fmt.Println(coord) // "a1"
    }

    coord, err = cell.FromIndices(4, 3)
    if err == nil {
        fmt.Println(coord) // "e4"
    }

    coord, err = cell.FromIndices(0, 0, 0)
    if err == nil {
        fmt.Println(coord) // "a1A"
    }

    // Must-style conversion (panics on error)
    coord = cell.MustFromIndices(1, 1, 1) // "b2B"

    // Round-trip conversion
    indices = cell.MustToIndices("e4")
    coord = cell.MustFromIndices(indices...)
    fmt.Println(coord) // "e4"
}
```

## Format Specification

### Dimensional Patterns

| Dimensions | Pattern | Examples |
|------------|---------|----------|
| 1D | `<lower>` | `a`, `e`, `file` |
| 2D | `<lower><integer>` | `a1`, `e4`, `aa10` |
| 3D | `<lower><integer><upper>` | `a1A`, `e4B` |
| 4D | `<lower><integer><upper><lower>` | `a1Ab`, `e4Bc` |
| 5D | `<lower><integer><upper><lower><integer>` | `a1Ab2` |

### Regular Expression

```regex
^[a-z]+(?:[1-9][0-9]*[A-Z]+[a-z]+)*(?:[1-9][0-9]*[A-Z]*)?$
```

### Valid Examples

| Coordinate | Dimensions | Description |
|------------|------------|-------------|
| `a` | 1D | Single file |
| `a1` | 2D | Standard chess-style |
| `e4` | 2D | Chess center |
| `a1A` | 3D | 3D tic-tac-toe |
| `h8Hh8` | 5D | Multi-dimensional |
| `aa1AA` | 3D | Extended alphabet |

### Invalid Examples

| String | Reason |
|--------|--------|
| `""` | Empty string |
| `1` | Starts with digit |
| `A` | Starts with uppercase |
| `a0` | Zero is not a valid positive integer |
| `a01` | Leading zero in numeric dimension |
| `aA` | Missing numeric dimension |
| `a1a` | Missing uppercase dimension |
| `a1A1` | Numeric after uppercase without lowercase |

## API Reference

### Validation

```go
func Valid(s string) bool
```

### Parsing

```go
func Parse(s string) ([]string, error)
func MustParse(s string) []string  // panics on error
```

### Dimensional Analysis

```go
func Dimensions(s string) int
```

### Coordinate Conversion

```go
func ToIndices(s string) ([]int, error)
func MustToIndices(s string) []int  // panics on error

func FromIndices(indices ...int) (string, error)
func MustFromIndices(indices ...int) string  // panics on error
```

### Regular Expression Access

```go
func Regex() *regexp.Regexp
```

## Game Examples

### Chess (8×8)

```go
// Standard chess coordinates
chessSquares := make([]string, 0, 64)
for file := 'a'; file <= 'h'; file++ {
    for rank := 1; rank <= 8; rank++ {
        chessSquares = append(chessSquares, fmt.Sprintf("%c%d", file, rank))
    }
}

// All are valid
for _, square := range chessSquares {
    fmt.Println(cell.Valid(square)) // true
}

// Convert position
fmt.Println(cell.MustToIndices("e4")) // [4, 3]
fmt.Println(cell.MustToIndices("h8")) // [7, 7]
```

### Shōgi (9×9)

```go
// Shōgi board positions
fmt.Println(cell.Valid("e5")) // true (center)
fmt.Println(cell.Valid("i9")) // true (corner)

fmt.Println(cell.MustToIndices("e5")) // [4, 4]
```

### 3D Tic-Tac-Toe (3×3×3)

```go
// Three-dimensional coordinates
fmt.Println(cell.Valid("a1A")) // true
fmt.Println(cell.Valid("b2B")) // true
fmt.Println(cell.Valid("c3C")) // true

// Winning diagonal
diagonal := []string{"a1A", "b2B", "c3C"}
for _, coord := range diagonal {
    fmt.Println(cell.MustToIndices(coord))
}
// [0, 0, 0]
// [1, 1, 1]
// [2, 2, 2]
```

## Extended Alphabet

CELL supports extended alphabet notation for large boards:

```go
// Single letters: a-z (positions 0-25)
fmt.Println(cell.MustToIndices("z1")) // [25, 0]

// Double letters: aa-zz (positions 26-701)
fmt.Println(cell.MustToIndices("aa1")) // [26, 0]
fmt.Println(cell.MustToIndices("ab1")) // [27, 0]
fmt.Println(cell.MustToIndices("zz1")) // [701, 0]

// And so on...
fmt.Println(cell.MustFromIndices(702, 0)) // "aaa1"
```

## Properties

- **Multi-dimensional**: Supports unlimited dimensional coordinate systems
- **Cyclical**: Uses systematic three-character-set repetition
- **ASCII-based**: Pure ASCII characters for universal compatibility
- **Unambiguous**: Each coordinate maps to exactly one location
- **Scalable**: Extends naturally from 1D to unlimited dimensions
- **Rule-agnostic**: Independent of specific game mechanics

## Related Specifications

- [Game Protocol](https://sashite.dev/game-protocol/) — Conceptual foundation
- [CELL Specification](https://sashite.dev/specs/cell/1.0.0/) — Official specification

## License

Available as open source under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0).

## About

Maintained by [Sashité](https://sashite.com/) — promoting chess variants and sharing the beauty of board game cultures.
