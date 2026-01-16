package cell

// Implementation constraints.
const (
	// MaxDimensions is the maximum number of dimensions supported.
	MaxDimensions = 3

	// MaxIndex is the maximum value for any single dimension index.
	MaxIndex = 255

	// MaxStringLen is the maximum length of a valid CELL string.
	// This corresponds to "iv256IV" (max value in all 3 dimensions).
	MaxStringLen = 7
)

// Coordinate represents a parsed CELL coordinate with up to 3 dimensions.
//
// The zero value is not valid; use [NewCoordinate] or [Parse] to create instances.
type Coordinate struct {
	indices [MaxDimensions]uint8
	dims    uint8
}

// NewCoordinate creates a Coordinate from 1 to 3 indices.
//
// It panics if no indices are provided or if more than 3 indices are given.
// For parsing user input, use [Parse] which returns an error instead.
func NewCoordinate(indices ...uint8) Coordinate {
	if len(indices) == 0 {
		panic("cell: NewCoordinate requires at least one index")
	}
	if len(indices) > MaxDimensions {
		panic("cell: NewCoordinate accepts at most 3 indices")
	}

	var c Coordinate
	c.dims = uint8(len(indices))
	copy(c.indices[:], indices)
	return c
}

// Dims returns the number of dimensions (1, 2, or 3).
func (c Coordinate) Dims() int {
	return int(c.dims)
}

// Indices returns the coordinate indices as a slice.
//
// The returned slice is a copy; modifying it does not affect the Coordinate.
func (c Coordinate) Indices() []uint8 {
	result := make([]uint8, c.dims)
	copy(result, c.indices[:c.dims])
	return result
}

// At returns the index at dimension i (0-indexed).
//
// It panics if i is out of range (i >= Dims()).
func (c Coordinate) At(i int) uint8 {
	if i < 0 || i >= int(c.dims) {
		panic("cell: index out of range")
	}
	return c.indices[i]
}

// String returns the CELL string representation (e.g., "e4", "a1A").
//
// This method implements [fmt.Stringer].
func (c Coordinate) String() string {
	return format(c)
}
