// Package cell implements CELL (Coordinate Encoding for Layered Locations) for Go.
//
// CELL is a standardized format for representing coordinates on multi-dimensional
// game boards using a cyclical ASCII character system.
//
// # Format
//
// CELL uses a cyclical three-character-set system:
//
//	| Dimension       | Condition | Character Set              |
//	|-----------------|-----------|----------------------------|
//	| 1st, 4th, 7th…  | n % 3 = 1 | Latin lowercase (`a`–`z`)  |
//	| 2nd, 5th, 8th…  | n % 3 = 2 | Positive integers          |
//	| 3rd, 6th, 9th…  | n % 3 = 0 | Latin uppercase (`A`–`Z`)  |
//
// # Examples
//
//	cell.Valid("a1")           // true
//	cell.Valid("a1A")          // true
//	cell.MustParse("e4")       // []string{"e", "4"}
//	cell.MustToIndices("e4")   // []int{4, 3}
//	cell.MustFromIndices([]int{4, 3}) // "e4"
//
// See the [CELL Specification] for details.
//
// [CELL Specification]: https://sashite.dev/specs/cell/1.0.0/
package cell

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// dimensionType represents the type of a dimension in the cyclical system.
type dimensionType int

const (
	lowercase dimensionType = iota
	numeric
	uppercase
)

// cellRegex is the regular expression from CELL Specification v1.0.0.
// Note: Line breaks must be rejected separately (see Valid).
var cellRegex = regexp.MustCompile(`^[a-z]+(?:[1-9][0-9]*[A-Z]+[a-z]+)*(?:[1-9][0-9]*[A-Z]*)?$`)

// Component extraction patterns.
var (
	lowercaseRegex = regexp.MustCompile(`^[a-z]+`)
	numericRegex   = regexp.MustCompile(`^[1-9][0-9]*`)
	uppercaseRegex = regexp.MustCompile(`^[A-Z]+`)
)

// --- Validation ---

// Valid checks if a string represents a valid CELL coordinate.
//
// Implements full-string matching as required by the CELL specification.
// Rejects any input containing line breaks (\r or \n).
//
// Examples:
//
//	cell.Valid("a1")    // true
//	cell.Valid("a1A")   // true
//	cell.Valid("e4")    // true
//	cell.Valid("a0")    // false
//	cell.Valid("")      // false
//	cell.Valid("1a")    // false
//	cell.Valid("a1\n")  // false
func Valid(s string) bool {
	if len(s) == 0 {
		return false
	}
	if strings.ContainsAny(s, "\r\n") {
		return false
	}
	return cellRegex.MatchString(s)
}

// Regex returns the validation regular expression from CELL specification v1.0.0.
//
// Note: This regex alone does not guarantee full compliance. The [Valid]
// function additionally rejects strings containing line breaks, as required
// by the specification's anchoring requirements.
func Regex() *regexp.Regexp {
	return cellRegex
}

// --- Parsing ---

// Parse parses a CELL coordinate string into dimensional components.
//
// Returns the components on success, or an error on failure.
//
// Examples:
//
//	cell.Parse("a1")      // []string{"a", "1"}, nil
//	cell.Parse("a1A")     // []string{"a", "1", "A"}, nil
//	cell.Parse("h8Hh8")   // []string{"h", "8", "H", "h", "8"}, nil
//	cell.Parse("foobar")  // []string{"foobar"}, nil
//	cell.Parse("invalid!") // nil, error
func Parse(s string) ([]string, error) {
	if !Valid(s) {
		return nil, fmt.Errorf("invalid CELL coordinate: %s", s)
	}
	return parseRecursive(s, 1), nil
}

// MustParse parses a CELL coordinate string into dimensional components.
//
// Returns the components on success, panics on failure.
//
// Examples:
//
//	cell.MustParse("a1A")    // []string{"a", "1", "A"}
//	cell.MustParse("1nvalid") // panics
func MustParse(s string) []string {
	components, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return components
}

// --- Dimensional Analysis ---

// Dimensions returns the number of dimensions in a CELL coordinate.
//
// Returns 0 for invalid coordinates.
//
// Examples:
//
//	cell.Dimensions("a")      // 1
//	cell.Dimensions("a1")     // 2
//	cell.Dimensions("a1A")    // 3
//	cell.Dimensions("h8Hh8")  // 5
//	cell.Dimensions("1nvalid") // 0
func Dimensions(s string) int {
	components, err := Parse(s)
	if err != nil {
		return 0
	}
	return len(components)
}

// --- Coordinate Conversion ---

// ToIndices converts a CELL coordinate to a slice of 0-indexed integers.
//
// Examples:
//
//	cell.ToIndices("a1")    // []int{0, 0}, nil
//	cell.ToIndices("e4")    // []int{4, 3}, nil
//	cell.ToIndices("a1A")   // []int{0, 0, 0}, nil
//	cell.ToIndices("z26Z")  // []int{25, 25, 25}, nil
//	cell.ToIndices("aa1AA") // []int{26, 0, 26}, nil
//	cell.ToIndices("1nvalid") // nil, error
func ToIndices(s string) ([]int, error) {
	components, err := Parse(s)
	if err != nil {
		return nil, err
	}

	indices := make([]int, len(components))
	for i, component := range components {
		dimension := i + 1
		dimType := getDimensionType(dimension)
		indices[i] = componentToIndex(component, dimType)
	}

	return indices, nil
}

// MustToIndices converts a CELL coordinate to a slice of 0-indexed integers.
//
// Returns the indices on success, panics on failure.
//
// Examples:
//
//	cell.MustToIndices("e4")    // []int{4, 3}
//	cell.MustToIndices("a1A")   // []int{0, 0, 0}
//	cell.MustToIndices("1nvalid") // panics
func MustToIndices(s string) []int {
	indices, err := ToIndices(s)
	if err != nil {
		panic(err)
	}
	return indices
}

// FromIndices converts a slice of 0-indexed integers to a CELL coordinate.
//
// Examples:
//
//	cell.FromIndices([]int{0, 0})      // "a1", nil
//	cell.FromIndices([]int{4, 3})      // "e4", nil
//	cell.FromIndices([]int{0, 0, 0})   // "a1A", nil
//	cell.FromIndices([]int{25, 25, 25}) // "z26Z", nil
//	cell.FromIndices([]int{26, 0, 26}) // "aa1AA", nil
//	cell.FromIndices([]int{})          // "", error
//	cell.FromIndices([]int{-1, 0})     // "", error
func FromIndices(indices []int) (string, error) {
	if len(indices) == 0 {
		return "", errors.New("cannot convert empty slice to CELL coordinate")
	}

	var builder strings.Builder
	for i, index := range indices {
		if index < 0 {
			return "", fmt.Errorf("negative index not allowed: %d", index)
		}
		dimension := i + 1
		dimType := getDimensionType(dimension)
		builder.WriteString(indexToComponent(index, dimType))
	}

	result := builder.String()
	if !Valid(result) {
		return "", fmt.Errorf("generated invalid CELL coordinate: %s", result)
	}

	return result, nil
}

// MustFromIndices converts a slice of 0-indexed integers to a CELL coordinate.
//
// Returns the coordinate on success, panics on failure.
//
// Examples:
//
//	cell.MustFromIndices([]int{4, 3})      // "e4"
//	cell.MustFromIndices([]int{0, 0, 0})   // "a1A"
//	cell.MustFromIndices([]int{})          // panics
func MustFromIndices(indices []int) string {
	coord, err := FromIndices(indices)
	if err != nil {
		panic(err)
	}
	return coord
}

// --- Private Functions ---

// parseRecursive recursively parses a coordinate string into components
// following the strict CELL specification cyclical pattern.
func parseRecursive(s string, dimension int) []string {
	if len(s) == 0 {
		return nil
	}

	dimType := getDimensionType(dimension)
	component, remaining := extractComponent(s, dimType)
	if component == "" {
		return nil
	}

	rest := parseRecursive(remaining, dimension+1)
	return append([]string{component}, rest...)
}

// getDimensionType determines the character set type for a given dimension.
// Following CELL specification cyclical system: dimension n % 3 determines character set.
func getDimensionType(dimension int) dimensionType {
	switch dimension % 3 {
	case 1:
		return lowercase
	case 2:
		return numeric
	default: // case 0
		return uppercase
	}
}

// extractComponent extracts the next component from a string based on expected type.
func extractComponent(s string, dimType dimensionType) (component, remaining string) {
	var re *regexp.Regexp
	switch dimType {
	case lowercase:
		re = lowercaseRegex
	case numeric:
		re = numericRegex
	case uppercase:
		re = uppercaseRegex
	}

	match := re.FindString(s)
	if match == "" {
		return "", s
	}
	return match, s[len(match):]
}

// componentToIndex converts a component to its 0-indexed position.
func componentToIndex(component string, dimType dimensionType) int {
	switch dimType {
	case lowercase:
		return lettersToIndex(component)
	case numeric:
		n, _ := strconv.Atoi(component)
		return n - 1
	case uppercase:
		return lettersToIndex(strings.ToLower(component))
	}
	return 0
}

// indexToComponent converts a 0-indexed position to a component.
func indexToComponent(index int, dimType dimensionType) string {
	switch dimType {
	case lowercase:
		return indexToLetters(index)
	case numeric:
		return strconv.Itoa(index + 1)
	case uppercase:
		return strings.ToUpper(indexToLetters(index))
	}
	return ""
}

// lettersToIndex converts a letter sequence to a 0-indexed position.
// Extended alphabet per CELL specification:
// a=0, b=1, ..., z=25, aa=26, ab=27, ..., zz=701, aaa=702, etc.
func lettersToIndex(letters string) int {
	length := len(letters)

	// Add positions from shorter sequences
	baseOffset := 0
	for l := 1; l < length; l++ {
		baseOffset += pow(26, l)
	}

	// Add position within current length
	positionInLength := 0
	for i, char := range letters {
		charValue := int(char - 'a')
		placeValue := pow(26, length-i-1)
		positionInLength += charValue * placeValue
	}

	return baseOffset + positionInLength
}

// indexToLetters converts a 0-indexed position to a letter sequence.
// Extended alphabet per CELL specification:
// 0=a, 1=b, ..., 25=z, 26=aa, 27=ab, ..., 701=zz, 702=aaa, etc.
func indexToLetters(index int) string {
	// Find the length of the result
	length, base := findLengthAndBase(index)

	// Convert within the found length
	adjustedIndex := index - base
	return buildLetters(adjustedIndex, length)
}

// findLengthAndBase finds the length and base offset for a given index.
func findLengthAndBase(index int) (length, base int) {
	length = 1
	base = 0
	for {
		rangeSize := pow(26, length)
		if index < base+rangeSize {
			return length, base
		}
		base += rangeSize
		length++
	}
}

// buildLetters builds the letter string from an adjusted index.
func buildLetters(index, length int) string {
	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = byte('a' + index%26)
		index /= 26
	}
	return string(result)
}

// pow computes base^exp for non-negative integers.
func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
