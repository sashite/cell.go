// Package cell implements CELL (Coordinate Encoding for Layered Locations),
// a standardized format for representing coordinates on multi-dimensional game boards.
//
// This package implements the [CELL Specification v1.0.0].
//
// # Implementation Constraints
//
// This implementation enforces the following constraints for safety and performance:
//
//   - Maximum 3 dimensions (sufficient for 1D, 2D, 3D boards)
//   - Maximum index value of 255 per dimension (fits in uint8)
//   - Maximum string length of 7 characters ("iv256IV")
//
// # Parsing
//
// Use [Parse] to convert a CELL string to a [Coordinate]:
//
//	coord, err := cell.Parse("e4")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(coord.Indices()) // [4, 3]
//
// For trusted input or constants, use [MustParse]:
//
//	var origin = cell.MustParse("a1")
//
// # Formatting
//
// Use [Coordinate.String] or [Format] to convert back to a CELL string:
//
//	coord := cell.NewCoordinate(4, 3)
//	fmt.Println(coord.String()) // "e4"
//
//	s := cell.Format(2, 2, 2)
//	fmt.Println(s) // "c3C"
//
// # Validation
//
// Use [Validate] for detailed errors or [IsValid] for a simple boolean check:
//
//	if err := cell.Validate("a0"); err != nil {
//	    fmt.Println(err) // "cell: leading zero in number"
//	}
//
//	if cell.IsValid("e4") {
//	    // valid coordinate
//	}
//
// # Error Handling
//
// All parsing errors are sentinel errors that can be checked with [errors.Is]:
//
//	coord, err := cell.Parse(input)
//	if errors.Is(err, cell.ErrLeadingZero) {
//	    // handle leading zero specifically
//	}
//
// [CELL Specification v1.0.0]: https://sashite.dev/specs/cell/1.0.0/
package cell
