# cell.go

[![Go Reference](https://pkg.go.dev/badge/github.com/sashite/cell.go.svg)](https://pkg.go.dev/github.com/sashite/cell.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sashite/cell.go)](https://goreportcard.com/report/github.com/sashite/cell.go)
[![License](https://img.shields.io/github/license/sashite/cell.go)](https://github.com/sashite/cell.go/blob/main/LICENSE)

> **CELL** (Coordinate Encoding for Layered Locations) implementation for Go.

## Overview

This library implements the [CELL Specification v1.0.0](https://sashite.dev/specs/cell/1.0.0/).

### Implementation Constraints

| Constraint | Value | Rationale |
|------------|-------|-----------|
| Max dimensions | 3 | Sufficient for 1D, 2D, 3D boards |
| Max index value | 255 | Fits in `uint8`, covers 256×256×256 boards |
| Max string length | 7 | `"iv256IV"` (max for all dimensions at 255) |

These constraints enable bounded memory usage and safe parsing without allocation.

## Installation

```bash
go get github.com/sashite/cell.go/v3
```

## Usage

### Parsing (String → Coordinate)

Convert a CELL string into a `Coordinate` struct.

```go
package main

import (
	"fmt"
	"github.com/sashite/cell.go/v3"
)

func main() {
	// Standard parsing (returns error)
	coord, err := cell.Parse("e4")
	if err != nil {
		panic(err)
	}
	fmt.Println(coord.Indices()) // [4, 3]
	fmt.Println(coord.Dims())    // 2

	// Panic on error (for constants or trusted input)
	c := cell.MustParse("a1A")
	fmt.Println(c.Indices()) // [0, 0, 0]
}
```

### Formatting (Coordinate → String)

Convert a `Coordinate` back to a CELL string.

```go
// From Coordinate struct
coord := cell.NewCoordinate(4, 3)
fmt.Println(coord.String()) // "e4"

// Direct formatting (convenience)
s := cell.Format(2, 2, 2)
fmt.Println(s) // "c3C"
```

### Validation

```go
// Boolean check
if cell.IsValid("e4") {
	// valid coordinate
}

// Detailed error
if err := cell.Validate("a0"); err != nil {
	fmt.Println(err) // "cell: leading zero"
}
```

### Accessing Coordinate Data

```go
coord := cell.MustParse("e4")

// Get dimensions count
coord.Dims() // 2

// Get indices as slice
coord.Indices() // []uint8{4, 3}

// Access individual index (panics if out of range)
coord.At(0) // 4
coord.At(1) // 3
```

## API Reference

### Types

```go
// Coordinate represents a parsed CELL coordinate with up to 3 dimensions.
// Zero value is not valid; use NewCoordinate or Parse to create.
type Coordinate struct {
	// contains filtered or unexported fields
}

// NewCoordinate creates a Coordinate from 1 to 3 indices.
// Panics if no indices provided or more than 3.
func NewCoordinate(indices ...uint8) Coordinate

// Dims returns the number of dimensions (1, 2, or 3).
func (c Coordinate) Dims() int

// Indices returns the coordinate indices as a slice.
func (c Coordinate) Indices() []uint8

// At returns the index at dimension i (0-indexed).
// Panics if i >= Dims().
func (c Coordinate) At(i int) uint8

// String returns the CELL string representation.
func (c Coordinate) String() string
```

### Parsing

```go
// Parse converts a CELL string to a Coordinate.
// Returns an error if the string is not valid.
func Parse(s string) (Coordinate, error)

// MustParse is like Parse but panics on error.
func MustParse(s string) Coordinate
```

### Formatting

```go
// Format converts indices to a CELL string.
// Convenience function equivalent to NewCoordinate(indices...).String().
func Format(indices ...uint8) string
```

### Validation

```go
// Validate checks if s is a valid CELL coordinate.
// Returns nil if valid, or a descriptive error.
func Validate(s string) error

// IsValid reports whether s is a valid CELL coordinate.
func IsValid(s string) bool
```

### Errors

```go
var (
	ErrEmptyInput      = errors.New("cell: empty input")
	ErrInputTooLong    = errors.New("cell: input exceeds 7 characters")
	ErrInvalidStart    = errors.New("cell: must start with lowercase letter")
	ErrUnexpectedChar  = errors.New("cell: unexpected character")
	ErrLeadingZero     = errors.New("cell: leading zero in number")
	ErrTooManyDims     = errors.New("cell: exceeds 3 dimensions")
	ErrIndexOutOfRange = errors.New("cell: index exceeds 255")
)
```

## Design Principles

- **Bounded types**: `uint8` indices prevent overflow
- **Struct over slice**: `Coordinate` type enables methods and safety
- **Sentinel errors**: Standard Go error handling with `errors.Is()`
- **strconv-style API**: Familiar `Parse`, `Must*`, `String()` patterns
- **No allocation in hot path**: Fixed-size struct, no heap allocation
- **No dependencies**: Pure Go standard library only

## Related Specifications

- [Game Protocol](https://sashite.dev/game-protocol/) — Conceptual foundation
- [CELL Specification](https://sashite.dev/specs/cell/1.0.0/) — Official specification
- [CELL Examples](https://sashite.dev/specs/cell/1.0.0/examples/) — Usage examples

## License

Available as open source under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0).
