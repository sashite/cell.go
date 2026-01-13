package cell

import (
	"errors"
	"testing"
)

// ----------------------------------------------------------------------------
// IsValid - Valid Coordinates
// ----------------------------------------------------------------------------

func TestIsValid_1D(t *testing.T) {
	valid := []string{"a", "z", "aa", "iv"}
	for _, s := range valid {
		if !IsValid(s) {
			t.Errorf("IsValid(%q) = false, want true", s)
		}
	}
}

func TestIsValid_2D(t *testing.T) {
	valid := []string{"a1", "e4", "h8", "z9", "a256", "iv256"}
	for _, s := range valid {
		if !IsValid(s) {
			t.Errorf("IsValid(%q) = false, want true", s)
		}
	}
}

func TestIsValid_3D(t *testing.T) {
	valid := []string{"a1A", "e4B", "c3C", "iv256IV"}
	for _, s := range valid {
		if !IsValid(s) {
			t.Errorf("IsValid(%q) = false, want true", s)
		}
	}
}

// ----------------------------------------------------------------------------
// IsValid - Invalid Coordinates
// ----------------------------------------------------------------------------

func TestIsValid_Empty(t *testing.T) {
	if IsValid("") {
		t.Error("IsValid(\"\") = true, want false")
	}
}

func TestIsValid_WrongStart(t *testing.T) {
	invalid := []string{"1", "A", "1a", "A1"}
	for _, s := range invalid {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false", s)
		}
	}
}

func TestIsValid_LeadingZero(t *testing.T) {
	invalid := []string{"a0", "a01", "a00"}
	for _, s := range invalid {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false", s)
		}
	}
}

func TestIsValid_WrongSequence(t *testing.T) {
	invalid := []string{"aA", "a1a", "a1A1"}
	for _, s := range invalid {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false", s)
		}
	}
}

func TestIsValid_InvalidChars(t *testing.T) {
	invalid := []string{"a1!", "a-1", "a 1", "a1\n"}
	for _, s := range invalid {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false", s)
		}
	}
}

func TestIsValid_InputTooLong(t *testing.T) {
	invalid := []string{"a1A1A1A1", "abcdefgh"}
	for _, s := range invalid {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false", s)
		}
	}
}

func TestIsValid_IndexOutOfRange(t *testing.T) {
	invalid := []string{"iw", "a257", "a1IW"}
	for _, s := range invalid {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false", s)
		}
	}
}

// ----------------------------------------------------------------------------
// Validate - Detailed Errors
// ----------------------------------------------------------------------------

func TestValidate_ErrEmptyInput(t *testing.T) {
	err := Validate("")
	if !errors.Is(err, ErrEmptyInput) {
		t.Errorf("Validate(\"\") = %v, want ErrEmptyInput", err)
	}
}

func TestValidate_ErrInputTooLong(t *testing.T) {
	cases := []string{"a1A1A1A1", "abcdefgh"}
	for _, s := range cases {
		err := Validate(s)
		if !errors.Is(err, ErrInputTooLong) {
			t.Errorf("Validate(%q) = %v, want ErrInputTooLong", s, err)
		}
	}
}

func TestValidate_ErrInvalidStart(t *testing.T) {
	cases := []string{"1a", "A1", "1", "A"}
	for _, s := range cases {
		err := Validate(s)
		if !errors.Is(err, ErrInvalidStart) {
			t.Errorf("Validate(%q) = %v, want ErrInvalidStart", s, err)
		}
	}
}

func TestValidate_ErrUnexpectedChar(t *testing.T) {
	cases := []string{"aA", "a1a", "a!", "a1!"}
	for _, s := range cases {
		err := Validate(s)
		if !errors.Is(err, ErrUnexpectedChar) {
			t.Errorf("Validate(%q) = %v, want ErrUnexpectedChar", s, err)
		}
	}
}

func TestValidate_ErrLeadingZero(t *testing.T) {
	cases := []string{"a0", "a01", "a00"}
	for _, s := range cases {
		err := Validate(s)
		if !errors.Is(err, ErrLeadingZero) {
			t.Errorf("Validate(%q) = %v, want ErrLeadingZero", s, err)
		}
	}
}

func TestValidate_ErrIndexOutOfRange(t *testing.T) {
	cases := []string{"iw", "a257", "a1IW"}
	for _, s := range cases {
		err := Validate(s)
		if !errors.Is(err, ErrIndexOutOfRange) {
			t.Errorf("Validate(%q) = %v, want ErrIndexOutOfRange", s, err)
		}
	}
}

func TestValidate_ErrTooManyDims(t *testing.T) {
	cases := []string{"a1Aa", "a1A!"}
	for _, s := range cases {
		err := Validate(s)
		if !errors.Is(err, ErrTooManyDims) {
			t.Errorf("Validate(%q) = %v, want ErrTooManyDims", s, err)
		}
	}
}

func TestValidate_Valid(t *testing.T) {
	cases := []string{"a1", "e4", "a1A", "iv256IV"}
	for _, s := range cases {
		if err := Validate(s); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", s, err)
		}
	}
}

// ----------------------------------------------------------------------------
// Parse - 1D Coordinates
// ----------------------------------------------------------------------------

func TestParse_1D(t *testing.T) {
	tests := []struct {
		input string
		want  []uint8
	}{
		{"a", []uint8{0}},
		{"z", []uint8{25}},
		{"aa", []uint8{26}},
		{"iv", []uint8{255}},
	}

	for _, tt := range tests {
		coord, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", tt.input, err)
			continue
		}
		if coord.Dims() != 1 {
			t.Errorf("Parse(%q).Dims() = %d, want 1", tt.input, coord.Dims())
		}
		got := coord.Indices()
		if !equalSlices(got, tt.want) {
			t.Errorf("Parse(%q).Indices() = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Parse - 2D Coordinates
// ----------------------------------------------------------------------------

func TestParse_2D(t *testing.T) {
	tests := []struct {
		input string
		want  []uint8
	}{
		{"a1", []uint8{0, 0}},
		{"e4", []uint8{4, 3}},
		{"h8", []uint8{7, 7}},
		{"a10", []uint8{0, 9}},
		{"iv256", []uint8{255, 255}},
	}

	for _, tt := range tests {
		coord, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", tt.input, err)
			continue
		}
		if coord.Dims() != 2 {
			t.Errorf("Parse(%q).Dims() = %d, want 2", tt.input, coord.Dims())
		}
		got := coord.Indices()
		if !equalSlices(got, tt.want) {
			t.Errorf("Parse(%q).Indices() = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Parse - 3D Coordinates
// ----------------------------------------------------------------------------

func TestParse_3D(t *testing.T) {
	tests := []struct {
		input string
		want  []uint8
	}{
		{"a1A", []uint8{0, 0, 0}},
		{"b2B", []uint8{1, 1, 1}},
		{"c3C", []uint8{2, 2, 2}},
		{"iv256IV", []uint8{255, 255, 255}},
	}

	for _, tt := range tests {
		coord, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", tt.input, err)
			continue
		}
		if coord.Dims() != 3 {
			t.Errorf("Parse(%q).Dims() = %d, want 3", tt.input, coord.Dims())
		}
		got := coord.Indices()
		if !equalSlices(got, tt.want) {
			t.Errorf("Parse(%q).Indices() = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Parse - Extended Alphabet
// ----------------------------------------------------------------------------

func TestParse_ExtendedAlphabet(t *testing.T) {
	tests := []struct {
		input string
		want  []uint8
	}{
		{"aa1", []uint8{26, 0}},
		{"az1", []uint8{51, 0}},
		{"ba1", []uint8{52, 0}},
		{"a1AA", []uint8{0, 0, 26}},
	}

	for _, tt := range tests {
		coord, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", tt.input, err)
			continue
		}
		got := coord.Indices()
		if !equalSlices(got, tt.want) {
			t.Errorf("Parse(%q).Indices() = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ----------------------------------------------------------------------------
// Parse - Errors
// ----------------------------------------------------------------------------

func TestParse_Errors(t *testing.T) {
	tests := []struct {
		input   string
		wantErr error
	}{
		{"", ErrEmptyInput},
		{"1a", ErrInvalidStart},
		{"a0", ErrLeadingZero},
		{"iw", ErrIndexOutOfRange},
	}

	for _, tt := range tests {
		_, err := Parse(tt.input)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Parse(%q) error = %v, want %v", tt.input, err, tt.wantErr)
		}
	}
}

// ----------------------------------------------------------------------------
// MustParse
// ----------------------------------------------------------------------------

func TestMustParse_Valid(t *testing.T) {
	coord := MustParse("e4")
	if coord.Dims() != 2 {
		t.Errorf("MustParse(\"e4\").Dims() = %d, want 2", coord.Dims())
	}
	got := coord.Indices()
	want := []uint8{4, 3}
	if !equalSlices(got, want) {
		t.Errorf("MustParse(\"e4\").Indices() = %v, want %v", got, want)
	}
}

func TestMustParse_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParse(\"a0\") did not panic")
		}
	}()
	MustParse("a0")
}

// ----------------------------------------------------------------------------
// Edge Cases
// ----------------------------------------------------------------------------

func TestParse_BoundaryValues(t *testing.T) {
	tests := []struct {
		input string
		want  []uint8
	}{
		// Minimum values
		{"a", []uint8{0}},
		{"a1", []uint8{0, 0}},
		{"a1A", []uint8{0, 0, 0}},
		// Maximum values
		{"iv", []uint8{255}},
		{"iv256", []uint8{255, 255}},
		{"iv256IV", []uint8{255, 255, 255}},
		// Just below maximum
		{"iu", []uint8{254}},
		{"a255", []uint8{0, 254}},
		{"a1IU", []uint8{0, 0, 254}},
	}

	for _, tt := range tests {
		coord, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", tt.input, err)
			continue
		}
		got := coord.Indices()
		if !equalSlices(got, tt.want) {
			t.Errorf("Parse(%q).Indices() = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParse_MaxStringLength(t *testing.T) {
	// Exactly 7 characters (maximum)
	if !IsValid("iv256IV") {
		t.Error("IsValid(\"iv256IV\") = false, want true")
	}

	// 8 characters (too long)
	if IsValid("iv256IVa") {
		t.Error("IsValid(\"iv256IVa\") = true, want false")
	}
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func equalSlices(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
