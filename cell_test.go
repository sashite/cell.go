package cell

import (
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------
// 1. Parsing Tests
// -----------------------------------------------------------------------------

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    []int
		wantErr error
	}{
		// --- Standard Cases (2D, 3D) ---
		{"a1", []int{0, 0}, nil},
		{"e4", []int{4, 3}, nil},
		{"z9", []int{25, 8}, nil},
		{"a1A", []int{0, 0, 0}, nil},
		{"c3C", []int{2, 2, 2}, nil},

		// --- Extended Alphabet (Bijective Base-26) ---
		{"aa1", []int{26, 0}, nil},   // 26th index
		{"ab1", []int{27, 0}, nil},   // 27th index
		{"az1", []int{51, 0}, nil},   // 51st index
		{"ba1", []int{52, 0}, nil},   // 52nd index
		{"zz1", []int{701, 0}, nil},  // 701st index
		{"aaa1", []int{702, 0}, nil}, // 702nd index

		// --- High Dimensions (Cycling) ---
		// Cycle: Lower -> Digit -> Upper -> Lower...
		{"a1Aa", []int{0, 0, 0, 0}, nil},
		{"a1Ab1", []int{0, 0, 0, 1, 0}, nil},

		// --- Edge Cases ---
		{"a", []int{0}, nil},        // 1D valid
		{"a10", []int{0, 9}, nil},   // Multi-digit numbers
		{"z99", []int{25, 98}, nil}, // Large numbers

		// --- Error Cases ---
		{"", nil, ErrEmpty},
		{"1", nil, ErrInvalidChar},    // Must start with Lower
		{"A", nil, ErrInvalidChar},    // Must start with Lower
		{"a0", nil, ErrInvalidFormat}, // 0 is invalid (1-based)
		{"a-1", nil, ErrInvalidChar},  // Negative/Symbols invalid
		{"aA", nil, ErrInvalidChar},   // Skipped digit dimension
		{"a1a", nil, ErrInvalidChar},  // Skipped upper dimension (expected A-Z)
		{"?", nil, ErrInvalidChar},    // Invalid character
		{" a1", nil, ErrInvalidChar},  // Leading space
		{"a1 ", nil, ErrInvalidChar},  // Trailing space
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)

			// 1. Check Error
			if err != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			// 2. Check Values
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 2. Formatting Tests
// -----------------------------------------------------------------------------

func TestFormat(t *testing.T) {
	tests := []struct {
		input []int
		want  string
	}{
		// Standard
		{[]int{0, 0}, "a1"},
		{[]int{4, 3}, "e4"},
		{[]int{2, 2, 2}, "c3C"},

		// Extended Alphabet
		{[]int{26, 0}, "aa1"},
		{[]int{701, 0}, "zz1"},
		{[]int{702, 0}, "aaa1"},

		// Cycling
		{[]int{0, 0, 0, 0}, "a1Aa"},

		// Empty
		{nil, ""},
		{[]int{}, ""},

		// Negative handling (graceful skip or empty logic based on implementation)
		// Current implementation skips negatives in loop but keeps index sync?
		// Actually implementation continues loop but writes nothing for that dim.
		// Since Parse forbids negatives, Format behavior on negatives is undefined/best-effort.
		// We test valid inputs primarily.
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			// Variadic call expansion
			got := Format(tt.input...)
			if got != tt.want {
				t.Errorf("Format(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 3. Round-Trip Property
// -----------------------------------------------------------------------------

func TestRoundTrip(t *testing.T) {
	// Property: Parse(Format(coords)) == coords
	inputs := [][]int{
		{0, 0},
		{25, 9, 0},
		{100, 100}, // 'cw101'
		{0, 0, 0, 0, 0},
	}

	for _, input := range inputs {
		str := Format(input...)
		parsed, err := Parse(str)
		if err != nil {
			t.Errorf("RoundTrip failed for %v: Parse returned error %v", input, err)
			continue
		}
		if !reflect.DeepEqual(parsed, input) {
			t.Errorf("RoundTrip mismatch for %v: got %v", input, parsed)
		}
	}
}

// -----------------------------------------------------------------------------
// 4. Zero-Allocation Logic (AppendParse)
// -----------------------------------------------------------------------------

func TestAppendParse(t *testing.T) {
	input := "c3C"
	expected := []int{2, 2, 2}

	// 1. Pre-allocate buffer
	// We use a capacity > 0 to simulate reuse
	buf := make([]int, 0, 10)

	// 2. Fill buffer with garbage to ensure we don't read stale data
	buf = append(buf, 999)
	buf = buf[:0] // Reset length

	// 3. AppendParse
	var err error
	buf, err = AppendParse(buf, input)
	if err != nil {
		t.Fatalf("AppendParse failed: %v", err)
	}

	// 4. Verify
	if !reflect.DeepEqual(buf, expected) {
		t.Errorf("AppendParse result %v, want %v", buf, expected)
	}

	// 5. Verify Cap (Optimization check)
	// The slice should not have grown if cap was sufficient
	if cap(buf) != 10 {
		t.Logf("Note: Buffer capacity changed from 10 to %d (re-allocation occurred?)", cap(buf))
	}
}

// -----------------------------------------------------------------------------
// 5. Must Helpers
// -----------------------------------------------------------------------------

func TestMustHelpers(t *testing.T) {
	// MustParse
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustParse panicked on valid input")
			}
		}()
		_ = MustParse("a1")
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustParse did NOT panic on invalid input")
			}
		}()
		_ = MustParse("INVALID")
	}()

	// MustFormat
	s := MustFormat(0, 0)
	if s != "a1" {
		t.Errorf("MustFormat(0,0) = %q, want %q", s, "a1")
	}
}

// -----------------------------------------------------------------------------
// 6. Benchmarks
// -----------------------------------------------------------------------------

// BenchmarkParse checks standard parsing (allocates slice).
func BenchmarkParse(b *testing.B) {
	input := "c3C" // 3D coordinate
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input)
	}
}

// BenchmarkAppendParse checks zero-allocation parsing.
// Ideally 0 allocs/op.
func BenchmarkAppendParse(b *testing.B) {
	input := "c3C"
	// Pre-allocate buffer outside loop
	buf := make([]int, 0, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf = buf[:0] // Reset
		buf, _ = AppendParse(buf, input)
	}
}

// BenchmarkFormat checks string formatting.
func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Format(2, 2, 2)
	}
}
