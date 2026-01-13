package cell

import "testing"

// ----------------------------------------------------------------------------
// String → Coordinate → String
// ----------------------------------------------------------------------------

func TestRoundTrip_StringToCoordinateToString_1D(t *testing.T) {
	cases := []string{"a", "z", "aa", "iv"}

	for _, original := range cases {
		coord, err := Parse(original)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", original, err)
			continue
		}

		got := coord.String()
		if got != original {
			t.Errorf("Parse(%q).String() = %q, want %q", original, got, original)
		}
	}
}

func TestRoundTrip_StringToCoordinateToString_2D(t *testing.T) {
	cases := []string{
		"a1", "e4", "h8", "z9",
		"a10", "a99", "a100", "a256",
		"aa1", "az1", "ba1",
		"iv256",
	}

	for _, original := range cases {
		coord, err := Parse(original)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", original, err)
			continue
		}

		got := coord.String()
		if got != original {
			t.Errorf("Parse(%q).String() = %q, want %q", original, got, original)
		}
	}
}

func TestRoundTrip_StringToCoordinateToString_3D(t *testing.T) {
	cases := []string{
		"a1A", "b2B", "c3C", "e4D", "z9Z",
		"a1AA", "a1AZ",
		"iv256IV",
	}

	for _, original := range cases {
		coord, err := Parse(original)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", original, err)
			continue
		}

		got := coord.String()
		if got != original {
			t.Errorf("Parse(%q).String() = %q, want %q", original, got, original)
		}
	}
}

// ----------------------------------------------------------------------------
// Coordinate → String → Coordinate
// ----------------------------------------------------------------------------

func TestRoundTrip_CoordinateToStringToCoordinate_1D(t *testing.T) {
	cases := [][]uint8{
		{0},
		{25},
		{26},
		{255},
	}

	for _, indices := range cases {
		original := NewCoordinate(indices...)
		s := original.String()

		parsed, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", s, err)
			continue
		}

		if original != parsed {
			t.Errorf("NewCoordinate(%v).String() = %q, Parse() = %v", indices, s, parsed.Indices())
		}
	}
}

func TestRoundTrip_CoordinateToStringToCoordinate_2D(t *testing.T) {
	cases := [][]uint8{
		{0, 0},
		{4, 3},
		{7, 7},
		{25, 8},
		{26, 9},
		{255, 255},
	}

	for _, indices := range cases {
		original := NewCoordinate(indices...)
		s := original.String()

		parsed, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", s, err)
			continue
		}

		if original != parsed {
			t.Errorf("NewCoordinate(%v).String() = %q, Parse() = %v", indices, s, parsed.Indices())
		}
	}
}

func TestRoundTrip_CoordinateToStringToCoordinate_3D(t *testing.T) {
	cases := [][]uint8{
		{0, 0, 0},
		{1, 1, 1},
		{2, 2, 2},
		{4, 3, 3},
		{25, 8, 25},
		{26, 9, 26},
		{255, 255, 255},
	}

	for _, indices := range cases {
		original := NewCoordinate(indices...)
		s := original.String()

		parsed, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", s, err)
			continue
		}

		if original != parsed {
			t.Errorf("NewCoordinate(%v).String() = %q, Parse() = %v", indices, s, parsed.Indices())
		}
	}
}

// ----------------------------------------------------------------------------
// Exhaustive Tests (Boundary Values)
// ----------------------------------------------------------------------------

func TestRoundTrip_AllSingleLetterValues(t *testing.T) {
	// Test a-z (0-25)
	for i := uint8(0); i < 26; i++ {
		original := NewCoordinate(i)
		s := original.String()

		parsed, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", s, err)
			continue
		}

		if original != parsed {
			t.Errorf("Round-trip failed for index %d: %q -> %v", i, s, parsed.Indices())
		}
	}
}

func TestRoundTrip_BoundaryValues(t *testing.T) {
	// Boundary values for each dimension type
	lowercaseBounds := []uint8{0, 25, 26, 51, 52, 255}
	numericBounds := []uint8{0, 8, 9, 98, 99, 255}
	uppercaseBounds := []uint8{0, 25, 26, 51, 52, 255}

	for _, l := range lowercaseBounds {
		for _, n := range numericBounds {
			for _, u := range uppercaseBounds {
				original := NewCoordinate(l, n, u)
				s := original.String()

				parsed, err := Parse(s)
				if err != nil {
					t.Errorf("Parse(%q) error = %v", s, err)
					continue
				}

				if original != parsed {
					t.Errorf("Round-trip failed for (%d, %d, %d): %q -> %v",
						l, n, u, s, parsed.Indices())
				}
			}
		}
	}
}

// ----------------------------------------------------------------------------
// Format Function Round-trip
// ----------------------------------------------------------------------------

func TestRoundTrip_FormatAndParse(t *testing.T) {
	tests := [][]uint8{
		{0},
		{4, 3},
		{2, 2, 2},
		{255, 255, 255},
	}

	for _, indices := range tests {
		s := Format(indices...)

		parsed, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(Format(%v)) error = %v", indices, err)
			continue
		}

		got := parsed.Indices()
		if !equalSlices(got, indices) {
			t.Errorf("Parse(Format(%v)) = %v", indices, got)
		}
	}
}

// ----------------------------------------------------------------------------
// MustParse Round-trip
// ----------------------------------------------------------------------------

func TestRoundTrip_MustParse(t *testing.T) {
	cases := []string{"e4", "a1A", "iv256IV"}

	for _, original := range cases {
		coord := MustParse(original)
		got := coord.String()

		if got != original {
			t.Errorf("MustParse(%q).String() = %q", original, got)
		}
	}
}
