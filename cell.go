package cell

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Common errors
var (
	ErrEmpty         = errors.New("cell: empty string")
	ErrInvalidFormat = errors.New("cell: invalid format")
	ErrInvalidChar   = errors.New("cell: character out of sequence or invalid")
)

// -----------------------------------------------------------------------------
// Parsing API
// -----------------------------------------------------------------------------

// Parse converts a CELL string (e.g., "c3C") into a slice of integer coordinates.
// It returns nil and an error if the format is invalid.
func Parse(s string) ([]int, error) {
	if len(s) == 0 {
		return nil, ErrEmpty
	}
	// Pre-allocate for common cases (2D or 3D) to avoid immediate re-alloc
	dst := make([]int, 0, 3)

	dst, err := AppendParse(dst, s)
	if err != nil {
		// Crucial: return nil slice on error to satisfy strict equality checks
		// and idiomatic Go behavior (don't return partial data on error).
		return nil, err
	}

	return dst, nil
}

// MustParse converts a CELL string into coordinates and panics on error.
// Useful for initializing constants or tests.
func MustParse(s string) []int {
	coords, err := Parse(s)
	if err != nil {
		panic(fmt.Sprintf("cell: MustParse(%q) failed: %v", s, err))
	}
	return coords
}

// AppendParse parses the CELL string s and appends the resulting coordinates to dst.
// This allows zero-allocation parsing if dst has sufficient capacity.
// Note: On error, dst may contain partial data appended before the error occurred.
func AppendParse(dst []int, s string) ([]int, error) {
	if len(s) == 0 {
		return dst, ErrEmpty
	}

	cursor := 0
	dim := 0 // 0=Lower, 1=Digit, 2=Upper
	length := len(s)

	for cursor < length {
		mode := dim % 3
		start := cursor

		// Consume characters valid for the current dimension mode
		switch mode {
		case 0: // Lowercase (a-z)
			for cursor < length && unicode.IsLower(rune(s[cursor])) {
				cursor++
			}
		case 1: // Digits (0-9)
			for cursor < length && unicode.IsDigit(rune(s[cursor])) {
				cursor++
			}
		case 2: // Uppercase (A-Z)
			for cursor < length && unicode.IsUpper(rune(s[cursor])) {
				cursor++
			}
		}

		// If no characters were consumed, implies the character at 'cursor'
		// is invalid for the current expected dimension.
		if cursor == start {
			return dst, ErrInvalidChar
		}

		// Decode the segment
		segment := s[start:cursor]
		val, err := decodeSegment(segment, mode)
		if err != nil {
			return dst, err
		}

		dst = append(dst, val)
		dim++
	}

	return dst, nil
}

// decodeSegment converts a substring to its integer value based on mode.
func decodeSegment(s string, mode int) (int, error) {
	// Mode 1: Digits (1-based in string, 0-based in int)
	if mode == 1 {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, ErrInvalidFormat
		}
		if val < 1 {
			return 0, ErrInvalidFormat // "0" is not valid in CELL (starts at "1")
		}
		return val - 1, nil
	}

	// Mode 0 & 2: Letters (Bijective Hexavigesimal System)
	// a=0, z=25, aa=26...
	val := 0
	baseChar := 'a'
	if mode == 2 {
		baseChar = 'A'
	}

	for _, r := range s {
		val = val*26 + int(r-baseChar) + 1
	}
	return val - 1, nil
}

// -----------------------------------------------------------------------------
// Formatting API
// -----------------------------------------------------------------------------

// Format converts a list of integer coordinates into a CELL string.
// It accepts a variadic number of integers.
func Format(indices ...int) string {
	if len(indices) == 0 {
		return ""
	}

	// Heuristic for buffer size: ~2 chars per coord + safety
	var sb strings.Builder
	sb.Grow(len(indices) * 2)

	for dim, val := range indices {
		if val < 0 {
			// Negative indices are not representable in standard CELL.
			// We skip them or treat them as invalid.
			// In this implementation, we skip to avoid generating garbage.
			continue
		}

		mode := dim % 3
		switch mode {
		case 0: // Lowercase
			writeAlpha(&sb, val, 'a')
		case 1: // Digits
			sb.WriteString(strconv.Itoa(val + 1))
		case 2: // Uppercase
			writeAlpha(&sb, val, 'A')
		}
	}

	return sb.String()
}

// MustFormat acts like Format.
// Since Format is safe/permissive, this is an alias for consistency.
func MustFormat(indices ...int) string {
	return Format(indices...)
}

// writeAlpha converts an index to bijective base-26 string (a, ..., z, aa, ...)
// and writes it to the builder.
func writeAlpha(sb *strings.Builder, val int, baseChar rune) {
	if val < 0 {
		return
	}

	// Algorithm for bijective base-26.
	// We build the string in a small stack buffer to avoid heap allocation.
	var buf [20]byte
	i := len(buf)

	n := val
	for {
		i--
		rem := n % 26
		buf[i] = byte(baseChar) + byte(rem)
		n = (n / 26) - 1
		if n < 0 {
			break
		}
	}

	sb.Write(buf[i:])
}
