package cell

import (
	"reflect"
	"testing"
)

// --- Valid Tests ---

func TestValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid coordinates
		{"1D single letter", "a", true},
		{"1D multiple letters", "foobar", true},
		{"2D standard chess", "a1", true},
		{"2D chess center", "e4", true},
		{"2D corner", "h8", true},
		{"2D extended file", "aa1", true},
		{"2D large rank", "a10", true},
		{"2D extended both", "aa10", true},
		{"3D basic", "a1A", true},
		{"3D extended", "e4B", true},
		{"3D all extended", "aa1AA", true},
		{"4D basic", "a1Ab", true},
		{"4D extended", "e4Bc", true},
		{"5D basic", "a1Ab2", true},
		{"5D complex", "h8Hh8", true},
		{"shogi center", "e5", true},
		{"shogi corner", "i9", true},
		{"3D tic-tac-toe", "b2B", true},
		{"large board", "z26Z", true},

		// Invalid coordinates
		{"empty string", "", false},
		{"starts with digit", "1", false},
		{"starts with uppercase", "A", false},
		{"zero rank", "a0", false},
		{"leading zero", "a01", false},
		{"missing numeric", "aA", false},
		{"missing uppercase", "a1a", false},
		{"numeric after uppercase without lowercase", "a1A1", false},
		{"starts with digit then letter", "1a", false},
		{"special character", "*", false},
		{"invalid character", "a1!", false},
		{"space in coordinate", "a 1", false},
		{"newline in coordinate", "a1\n", false},
		{"carriage return", "a1\r", false},
		{"tab character", "a\t1", false},
		{"unicode letter", "α1", false},
		{"mixed case in dimension", "aB1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Valid(tt.input)
			if result != tt.expected {
				t.Errorf("Valid(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRegex(t *testing.T) {
	re := Regex()
	if re == nil {
		t.Error("Regex() returned nil")
	}

	// Test that the regex source matches the specification
	expectedPattern := `^[a-z]+(?:[1-9][0-9]*[A-Z]+[a-z]+)*(?:[1-9][0-9]*[A-Z]*)?$`
	if re.String() != expectedPattern {
		t.Errorf("Regex pattern = %q, expected %q", re.String(), expectedPattern)
	}
}

// --- Parse Tests ---

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		{"1D single", "a", []string{"a"}, false},
		{"1D multiple", "foobar", []string{"foobar"}, false},
		{"2D basic", "a1", []string{"a", "1"}, false},
		{"2D chess", "e4", []string{"e", "4"}, false},
		{"2D extended file", "aa10", []string{"aa", "10"}, false},
		{"3D basic", "a1A", []string{"a", "1", "A"}, false},
		{"3D extended", "aa1AA", []string{"aa", "1", "AA"}, false},
		{"4D basic", "a1Ab", []string{"a", "1", "A", "b"}, false},
		{"5D complex", "h8Hh8", []string{"h", "8", "H", "h", "8"}, false},

		// Error cases
		{"empty string", "", nil, true},
		{"invalid start", "1a", nil, true},
		{"invalid char", "a1!", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Parse(%q) = %v, expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	// Test successful parsing
	result := MustParse("a1A")
	expected := []string{"a", "1", "A"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MustParse(\"a1A\") = %v, expected %v", result, expected)
	}

	// Test panic on invalid input
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParse(\"1nvalid\") expected panic, got none")
		}
	}()
	MustParse("1nvalid")
}

// --- Dimensions Tests ---

func TestDimensions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"1D single", "a", 1},
		{"1D multiple", "foobar", 1},
		{"2D basic", "a1", 2},
		{"2D chess", "e4", 2},
		{"3D basic", "a1A", 3},
		{"4D basic", "a1Ab", 4},
		{"5D complex", "h8Hh8", 5},

		// Invalid coordinates return 0
		{"invalid empty", "", 0},
		{"invalid start", "1nvalid", 0},
		{"invalid char", "a!", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Dimensions(tt.input)
			if result != tt.expected {
				t.Errorf("Dimensions(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// --- ToIndices Tests ---

func TestToIndices(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []int
		expectError bool
	}{
		// Basic conversions
		{"2D origin", "a1", []int{0, 0}, false},
		{"2D chess e4", "e4", []int{4, 3}, false},
		{"2D corner h8", "h8", []int{7, 7}, false},
		{"3D origin", "a1A", []int{0, 0, 0}, false},
		{"3D b2B", "b2B", []int{1, 1, 1}, false},
		{"3D z26Z", "z26Z", []int{25, 25, 25}, false},

		// Extended alphabet
		{"extended aa1", "aa1", []int{26, 0}, false},
		{"extended ab1", "ab1", []int{27, 0}, false},
		{"extended az1", "az1", []int{51, 0}, false},
		{"extended ba1", "ba1", []int{52, 0}, false},
		{"extended zz1", "zz1", []int{701, 0}, false},
		{"extended aa1AA", "aa1AA", []int{26, 0, 26}, false},

		// Shogi
		{"shogi center e5", "e5", []int{4, 4}, false},
		{"shogi corner i9", "i9", []int{8, 8}, false},

		// 3D tic-tac-toe diagonal
		{"3d ttt a1A", "a1A", []int{0, 0, 0}, false},
		{"3d ttt b2B", "b2B", []int{1, 1, 1}, false},
		{"3d ttt c3C", "c3C", []int{2, 2, 2}, false},

		// 1D
		{"1D a", "a", []int{0}, false},
		{"1D z", "z", []int{25}, false},
		{"1D aa", "aa", []int{26}, false},

		// Error cases
		{"invalid empty", "", nil, true},
		{"invalid start", "1nvalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToIndices(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("ToIndices(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToIndices(%q) unexpected error: %v", tt.input, err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("ToIndices(%q) = %v, expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestMustToIndices(t *testing.T) {
	// Test successful conversion
	result := MustToIndices("e4")
	expected := []int{4, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MustToIndices(\"e4\") = %v, expected %v", result, expected)
	}

	// Test panic on invalid input
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustToIndices(\"1nvalid\") expected panic, got none")
		}
	}()
	MustToIndices("1nvalid")
}

// --- FromIndices Tests ---

func TestFromIndices(t *testing.T) {
	tests := []struct {
		name        string
		input       []int
		expected    string
		expectError bool
	}{
		// Basic conversions
		{"2D origin", []int{0, 0}, "a1", false},
		{"2D chess e4", []int{4, 3}, "e4", false},
		{"2D corner h8", []int{7, 7}, "h8", false},
		{"3D origin", []int{0, 0, 0}, "a1A", false},
		{"3D b2B", []int{1, 1, 1}, "b2B", false},
		{"3D z26Z", []int{25, 25, 25}, "z26Z", false},

		// Extended alphabet
		{"extended aa1", []int{26, 0}, "aa1", false},
		{"extended ab1", []int{27, 0}, "ab1", false},
		{"extended az1", []int{51, 0}, "az1", false},
		{"extended ba1", []int{52, 0}, "ba1", false},
		{"extended zz1", []int{701, 0}, "zz1", false},
		{"extended aaa1", []int{702, 0}, "aaa1", false},
		{"extended aa1AA", []int{26, 0, 26}, "aa1AA", false},

		// 1D
		{"1D a", []int{0}, "a", false},
		{"1D z", []int{25}, "z", false},
		{"1D aa", []int{26}, "aa", false},

		// 4D and 5D
		{"4D basic", []int{0, 0, 0, 0}, "a1Aa", false},
		{"5D basic", []int{0, 0, 0, 0, 0}, "a1Aa1", false},

		// Error cases
		{"empty slice", []int{}, "", true},
		{"negative index", []int{-1, 0}, "", true},
		{"negative second", []int{0, -1}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromIndices(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("FromIndices(%v) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("FromIndices(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("FromIndices(%v) = %q, expected %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestMustFromIndices(t *testing.T) {
	// Test successful conversion
	result := MustFromIndices([]int{4, 3})
	expected := "e4"
	if result != expected {
		t.Errorf("MustFromIndices([]int{4, 3}) = %q, expected %q", result, expected)
	}

	// Test panic on invalid input
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFromIndices([]int{}) expected panic, got none")
		}
	}()
	MustFromIndices([]int{})
}

// --- Round-trip Tests ---

func TestRoundTrip(t *testing.T) {
	coordinates := []string{
		"a", "z", "aa", "az", "ba", "zz", "aaa",
		"a1", "e4", "h8", "aa10", "zz99",
		"a1A", "b2B", "c3C", "z26Z", "aa1AA",
		"a1Ab", "e4Bc",
		"a1Ab2", "h8Hh8",
	}

	for _, coord := range coordinates {
		t.Run(coord, func(t *testing.T) {
			indices, err := ToIndices(coord)
			if err != nil {
				t.Fatalf("ToIndices(%q) failed: %v", coord, err)
			}

			result, err := FromIndices(indices)
			if err != nil {
				t.Fatalf("FromIndices(%v) failed: %v", indices, err)
			}

			if result != coord {
				t.Errorf("Round-trip failed: %q -> %v -> %q", coord, indices, result)
			}
		})
	}
}

func TestRoundTripFromIndices(t *testing.T) {
	testCases := [][]int{
		{0}, {25}, {26}, {701}, {702},
		{0, 0}, {4, 3}, {7, 7}, {25, 25}, {26, 0},
		{0, 0, 0}, {1, 1, 1}, {25, 25, 25}, {26, 0, 26},
		{0, 0, 0, 0}, {1, 2, 3, 4},
		{0, 0, 0, 0, 0}, {7, 7, 7, 7, 7},
	}

	for _, indices := range testCases {
		t.Run("", func(t *testing.T) {
			coord, err := FromIndices(indices)
			if err != nil {
				t.Fatalf("FromIndices(%v) failed: %v", indices, err)
			}

			result, err := ToIndices(coord)
			if err != nil {
				t.Fatalf("ToIndices(%q) failed: %v", coord, err)
			}

			if !reflect.DeepEqual(result, indices) {
				t.Errorf("Round-trip failed: %v -> %q -> %v", indices, coord, result)
			}
		})
	}
}

// --- Extended Alphabet Tests ---

func TestExtendedAlphabet(t *testing.T) {
	// Test the boundaries of the extended alphabet system
	tests := []struct {
		letters string
		index   int
	}{
		// Single letters: 0-25
		{"a", 0},
		{"b", 1},
		{"z", 25},

		// Double letters: 26-701
		{"aa", 26},
		{"ab", 27},
		{"az", 51},
		{"ba", 52},
		{"bz", 77},
		{"za", 676},
		{"zz", 701},

		// Triple letters: 702+
		{"aaa", 702},
		{"aab", 703},
		{"aba", 728},
		{"azz", 1377},
		{"baa", 1378},
	}

	for _, tt := range tests {
		t.Run(tt.letters, func(t *testing.T) {
			// Test lettersToIndex
			result := lettersToIndex(tt.letters)
			if result != tt.index {
				t.Errorf("lettersToIndex(%q) = %d, expected %d", tt.letters, result, tt.index)
			}

			// Test indexToLetters
			letters := indexToLetters(tt.index)
			if letters != tt.letters {
				t.Errorf("indexToLetters(%d) = %q, expected %q", tt.index, letters, tt.letters)
			}
		})
	}
}

// --- Chess Board Tests ---

func TestChessBoard(t *testing.T) {
	// Verify all 64 chess squares are valid
	for file := 'a'; file <= 'h'; file++ {
		for rank := 1; rank <= 8; rank++ {
			coord := string(file) + string(rune('0'+rank))
			if !Valid(coord) {
				t.Errorf("Chess square %q should be valid", coord)
			}
		}
	}

	// Verify specific positions
	positions := map[string][]int{
		"a1": {0, 0}, // white rook
		"e1": {4, 0}, // white king
		"d1": {3, 0}, // white queen
		"e4": {4, 3}, // classic opening
		"d5": {3, 4}, // center
		"h8": {7, 7}, // black rook
	}

	for coord, expected := range positions {
		indices, err := ToIndices(coord)
		if err != nil {
			t.Errorf("ToIndices(%q) failed: %v", coord, err)
			continue
		}
		if !reflect.DeepEqual(indices, expected) {
			t.Errorf("ToIndices(%q) = %v, expected %v", coord, indices, expected)
		}
	}
}

// --- Shogi Board Tests ---

func TestShogiBoard(t *testing.T) {
	// Verify all 81 shogi squares are valid
	for file := 'a'; file <= 'i'; file++ {
		for rank := 1; rank <= 9; rank++ {
			coord := string(file) + string(rune('0'+rank))
			if !Valid(coord) {
				t.Errorf("Shogi square %q should be valid", coord)
			}
		}
	}

	// Verify specific positions
	positions := map[string][]int{
		"e1": {4, 0}, // sente king initial
		"e9": {4, 8}, // gote king initial
		"e5": {4, 4}, // center
		"a1": {0, 0}, // corner
		"i9": {8, 8}, // opposite corner
	}

	for coord, expected := range positions {
		indices, err := ToIndices(coord)
		if err != nil {
			t.Errorf("ToIndices(%q) failed: %v", coord, err)
			continue
		}
		if !reflect.DeepEqual(indices, expected) {
			t.Errorf("ToIndices(%q) = %v, expected %v", coord, indices, expected)
		}
	}
}

// --- 3D Tic-Tac-Toe Tests ---

func TestTicTacToe3D(t *testing.T) {
	// Verify all 27 positions are valid
	for file := 'a'; file <= 'c'; file++ {
		for rank := 1; rank <= 3; rank++ {
			for level := 'A'; level <= 'C'; level++ {
				coord := string(file) + string(rune('0'+rank)) + string(level)
				if !Valid(coord) {
					t.Errorf("3D tic-tac-toe position %q should be valid", coord)
				}
			}
		}
	}

	// Verify winning diagonal
	diagonal := []struct {
		coord   string
		indices []int
	}{
		{"a1A", []int{0, 0, 0}},
		{"b2B", []int{1, 1, 1}},
		{"c3C", []int{2, 2, 2}},
	}

	for _, tt := range diagonal {
		indices, err := ToIndices(tt.coord)
		if err != nil {
			t.Errorf("ToIndices(%q) failed: %v", tt.coord, err)
			continue
		}
		if !reflect.DeepEqual(indices, tt.indices) {
			t.Errorf("ToIndices(%q) = %v, expected %v", tt.coord, indices, tt.indices)
		}
	}
}

// --- Benchmark Tests ---

func BenchmarkValid(b *testing.B) {
	coords := []string{"a1", "e4", "a1A", "h8Hh8", "aa1AA"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, coord := range coords {
			Valid(coord)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	coords := []string{"a1", "e4", "a1A", "h8Hh8", "aa1AA"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, coord := range coords {
			Parse(coord)
		}
	}
}

func BenchmarkToIndices(b *testing.B) {
	coords := []string{"a1", "e4", "a1A", "h8Hh8", "aa1AA"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, coord := range coords {
			ToIndices(coord)
		}
	}
}

func BenchmarkFromIndices(b *testing.B) {
	indices := [][]int{{0, 0}, {4, 3}, {0, 0, 0}, {7, 7, 7, 7, 7}, {26, 0, 26}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			FromIndices(idx)
		}
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	coord := "h8Hh8"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indices, _ := ToIndices(coord)
		FromIndices(indices)
	}
}
