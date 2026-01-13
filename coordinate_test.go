package cell

import "testing"

// ----------------------------------------------------------------------------
// NewCoordinate
// ----------------------------------------------------------------------------

func TestNewCoordinate_1D(t *testing.T) {
	coord := NewCoordinate(5)

	if coord.Dims() != 1 {
		t.Errorf("NewCoordinate(5).Dims() = %d, want 1", coord.Dims())
	}

	got := coord.Indices()
	want := []uint8{5}
	if !equalSlices(got, want) {
		t.Errorf("NewCoordinate(5).Indices() = %v, want %v", got, want)
	}
}

func TestNewCoordinate_2D(t *testing.T) {
	coord := NewCoordinate(4, 3)

	if coord.Dims() != 2 {
		t.Errorf("NewCoordinate(4, 3).Dims() = %d, want 2", coord.Dims())
	}

	got := coord.Indices()
	want := []uint8{4, 3}
	if !equalSlices(got, want) {
		t.Errorf("NewCoordinate(4, 3).Indices() = %v, want %v", got, want)
	}
}

func TestNewCoordinate_3D(t *testing.T) {
	coord := NewCoordinate(0, 0, 0)

	if coord.Dims() != 3 {
		t.Errorf("NewCoordinate(0, 0, 0).Dims() = %d, want 3", coord.Dims())
	}

	got := coord.Indices()
	want := []uint8{0, 0, 0}
	if !equalSlices(got, want) {
		t.Errorf("NewCoordinate(0, 0, 0).Indices() = %v, want %v", got, want)
	}
}

func TestNewCoordinate_MaxValues(t *testing.T) {
	coord := NewCoordinate(255, 255, 255)

	if coord.Dims() != 3 {
		t.Errorf("NewCoordinate(255, 255, 255).Dims() = %d, want 3", coord.Dims())
	}

	got := coord.Indices()
	want := []uint8{255, 255, 255}
	if !equalSlices(got, want) {
		t.Errorf("NewCoordinate(255, 255, 255).Indices() = %v, want %v", got, want)
	}
}

func TestNewCoordinate_PanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewCoordinate() did not panic")
		}
	}()
	NewCoordinate()
}

func TestNewCoordinate_PanicsOnTooMany(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewCoordinate(1, 2, 3, 4) did not panic")
		}
	}()
	NewCoordinate(1, 2, 3, 4)
}

// ----------------------------------------------------------------------------
// Dims
// ----------------------------------------------------------------------------

func TestCoordinate_Dims(t *testing.T) {
	tests := []struct {
		coord Coordinate
		want  int
	}{
		{NewCoordinate(0), 1},
		{NewCoordinate(0, 0), 2},
		{NewCoordinate(0, 0, 0), 3},
	}

	for _, tt := range tests {
		if got := tt.coord.Dims(); got != tt.want {
			t.Errorf("Coordinate.Dims() = %d, want %d", got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Indices
// ----------------------------------------------------------------------------

func TestCoordinate_Indices_ReturnsCopy(t *testing.T) {
	coord := NewCoordinate(1, 2, 3)

	indices := coord.Indices()
	indices[0] = 99 // Modify the returned slice

	// Original should be unchanged
	if coord.At(0) != 1 {
		t.Error("Modifying Indices() result affected the original Coordinate")
	}
}

func TestCoordinate_Indices_CorrectLength(t *testing.T) {
	tests := []struct {
		coord   Coordinate
		wantLen int
	}{
		{NewCoordinate(5), 1},
		{NewCoordinate(4, 3), 2},
		{NewCoordinate(0, 0, 0), 3},
	}

	for _, tt := range tests {
		got := tt.coord.Indices()
		if len(got) != tt.wantLen {
			t.Errorf("len(Coordinate.Indices()) = %d, want %d", len(got), tt.wantLen)
		}
	}
}

// ----------------------------------------------------------------------------
// At
// ----------------------------------------------------------------------------

func TestCoordinate_At(t *testing.T) {
	coord := NewCoordinate(10, 20, 30)

	tests := []struct {
		index int
		want  uint8
	}{
		{0, 10},
		{1, 20},
		{2, 30},
	}

	for _, tt := range tests {
		if got := coord.At(tt.index); got != tt.want {
			t.Errorf("Coordinate.At(%d) = %d, want %d", tt.index, got, tt.want)
		}
	}
}

func TestCoordinate_At_PanicsOnNegative(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Coordinate.At(-1) did not panic")
		}
	}()

	coord := NewCoordinate(1, 2, 3)
	coord.At(-1)
}

func TestCoordinate_At_PanicsOnOutOfRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Coordinate.At(2) on 2D coordinate did not panic")
		}
	}()

	coord := NewCoordinate(1, 2) // 2D coordinate
	coord.At(2)                  // Index 2 is out of range
}

// ----------------------------------------------------------------------------
// String (via fmt.Stringer)
// ----------------------------------------------------------------------------

func TestCoordinate_String_ImplementsStringer(t *testing.T) {
	coord := NewCoordinate(4, 3)

	// String() should return the CELL representation
	got := coord.String()
	want := "e4"

	if got != want {
		t.Errorf("Coordinate.String() = %q, want %q", got, want)
	}
}

// ----------------------------------------------------------------------------
// Zero Value
// ----------------------------------------------------------------------------

func TestCoordinate_ZeroValue(t *testing.T) {
	var coord Coordinate

	// Zero value has 0 dimensions
	if coord.Dims() != 0 {
		t.Errorf("Zero Coordinate.Dims() = %d, want 0", coord.Dims())
	}

	// Indices returns empty slice
	if len(coord.Indices()) != 0 {
		t.Errorf("Zero Coordinate.Indices() = %v, want []", coord.Indices())
	}

	// String returns empty string
	if coord.String() != "" {
		t.Errorf("Zero Coordinate.String() = %q, want \"\"", coord.String())
	}
}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

func TestCoordinate_Equality(t *testing.T) {
	a := NewCoordinate(4, 3)
	b := NewCoordinate(4, 3)
	c := NewCoordinate(4, 4)

	// Same values should be equal (struct comparison)
	if a != b {
		t.Error("Equal coordinates are not equal")
	}

	// Different values should not be equal
	if a == c {
		t.Error("Different coordinates are equal")
	}
}

func TestCoordinate_EqualityDifferentDims(t *testing.T) {
	a := NewCoordinate(0, 0)
	b := NewCoordinate(0, 0, 0)

	// Different dimensions should not be equal
	if a == b {
		t.Error("Coordinates with different dimensions are equal")
	}
}
