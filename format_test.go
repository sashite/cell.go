package cell

import "testing"

// ----------------------------------------------------------------------------
// Format - 1D Coordinates
// ----------------------------------------------------------------------------

func TestFormat_1D(t *testing.T) {
	tests := []struct {
		indices []uint8
		want    string
	}{
		{[]uint8{0}, "a"},
		{[]uint8{25}, "z"},
		{[]uint8{26}, "aa"},
		{[]uint8{255}, "iv"},
	}

	for _, tt := range tests {
		got := Format(tt.indices...)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.indices, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Format - 2D Coordinates
// ----------------------------------------------------------------------------

func TestFormat_2D(t *testing.T) {
	tests := []struct {
		indices []uint8
		want    string
	}{
		{[]uint8{0, 0}, "a1"},
		{[]uint8{4, 3}, "e4"},
		{[]uint8{7, 7}, "h8"},
		{[]uint8{0, 9}, "a10"},
		{[]uint8{0, 99}, "a100"},
		{[]uint8{255, 255}, "iv256"},
	}

	for _, tt := range tests {
		got := Format(tt.indices...)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.indices, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Format - 3D Coordinates
// ----------------------------------------------------------------------------

func TestFormat_3D(t *testing.T) {
	tests := []struct {
		indices []uint8
		want    string
	}{
		{[]uint8{0, 0, 0}, "a1A"},
		{[]uint8{1, 1, 1}, "b2B"},
		{[]uint8{2, 2, 2}, "c3C"},
		{[]uint8{255, 255, 255}, "iv256IV"},
	}

	for _, tt := range tests {
		got := Format(tt.indices...)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.indices, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Format - Extended Alphabet
// ----------------------------------------------------------------------------

func TestFormat_ExtendedAlphabet(t *testing.T) {
	tests := []struct {
		indices []uint8
		want    string
	}{
		{[]uint8{26, 0}, "aa1"},
		{[]uint8{51, 0}, "az1"},
		{[]uint8{52, 0}, "ba1"},
		{[]uint8{0, 0, 26}, "a1AA"},
		{[]uint8{0, 0, 51}, "a1AZ"},
	}

	for _, tt := range tests {
		got := Format(tt.indices...)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.indices, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Coordinate.String
// ----------------------------------------------------------------------------

func TestCoordinate_String(t *testing.T) {
	tests := []struct {
		coord Coordinate
		want  string
	}{
		{NewCoordinate(0), "a"},
		{NewCoordinate(4, 3), "e4"},
		{NewCoordinate(0, 0, 0), "a1A"},
		{NewCoordinate(255, 255, 255), "iv256IV"},
	}

	for _, tt := range tests {
		got := tt.coord.String()
		if got != tt.want {
			t.Errorf("Coordinate%v.String() = %q, want %q", tt.coord.Indices(), got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Edge Cases - Boundaries
// ----------------------------------------------------------------------------

func TestFormat_LetterBoundary(t *testing.T) {
	// z = 25 (single letter)
	if got := Format(25); got != "z" {
		t.Errorf("Format(25) = %q, want \"z\"", got)
	}

	// aa = 26 (double letter)
	if got := Format(26); got != "aa" {
		t.Errorf("Format(26) = %q, want \"aa\"", got)
	}
}

func TestFormat_DigitBoundary_SingleToDouble(t *testing.T) {
	// 9 = index 8 (single digit)
	if got := Format(0, 8); got != "a9" {
		t.Errorf("Format(0, 8) = %q, want \"a9\"", got)
	}

	// 10 = index 9 (double digit)
	if got := Format(0, 9); got != "a10" {
		t.Errorf("Format(0, 9) = %q, want \"a10\"", got)
	}
}

func TestFormat_DigitBoundary_DoubleToTriple(t *testing.T) {
	// 99 = index 98 (double digit)
	if got := Format(0, 98); got != "a99" {
		t.Errorf("Format(0, 98) = %q, want \"a99\"", got)
	}

	// 100 = index 99 (triple digit)
	if got := Format(0, 99); got != "a100" {
		t.Errorf("Format(0, 99) = %q, want \"a100\"", got)
	}
}

// ----------------------------------------------------------------------------
// String Length
// ----------------------------------------------------------------------------

func TestFormat_StringLength(t *testing.T) {
	tests := []struct {
		indices []uint8
		wantLen int
	}{
		{[]uint8{0}, 1},             // "a"
		{[]uint8{4, 3}, 2},          // "e4"
		{[]uint8{0, 0, 0}, 3},       // "a1A"
		{[]uint8{255, 255, 255}, 7}, // "iv256IV"
	}

	for _, tt := range tests {
		got := Format(tt.indices...)
		if len(got) != tt.wantLen {
			t.Errorf("len(Format(%v)) = %d, want %d (got %q)", tt.indices, len(got), tt.wantLen, got)
		}
	}
}
