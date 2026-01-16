package cell

import "testing"

// ----------------------------------------------------------------------------
// Security - Malicious Inputs
// ----------------------------------------------------------------------------

// TestSecurity_NullByteInjection verifies null bytes are rejected.
func TestSecurity_NullByteInjection(t *testing.T) {
	cases := []string{
		"a\x00",
		"a\x001",
		"\x00a1",
		"a1\x00",
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (null byte)", s)
		}
	}
}

// TestSecurity_NewlineInjection verifies newlines are rejected.
func TestSecurity_NewlineInjection(t *testing.T) {
	cases := []string{
		"a\n",
		"a\n1",
		"a1\n",
		"\na1",
		"a\r1",
		"a\r\n1",
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (newline)", s)
		}
	}
}

// TestSecurity_TabInjection verifies tabs are rejected.
func TestSecurity_TabInjection(t *testing.T) {
	cases := []string{
		"a\t",
		"a\t1",
		"\ta1",
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (tab)", s)
		}
	}
}

// TestSecurity_UnicodeLetterLookalikes verifies Unicode lookalikes are rejected.
// These characters visually resemble ASCII letters but are different code points.
func TestSecurity_UnicodeLetterLookalikes(t *testing.T) {
	cases := []string{
		"\u0430",       // Cyrillic small letter A (U+0430) looks like 'a'
		"\u0435",       // Cyrillic small letter IE (U+0435) looks like 'e'
		"\u043E",       // Cyrillic small letter O (U+043E) looks like 'o'
		"\u0410",       // Cyrillic capital letter A (U+0410) looks like 'A'
		"\u0430\u0431", // Cyrillic 'ab'
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (Unicode lookalike)", s)
		}
	}
}

// TestSecurity_FullWidthCharacters verifies full-width characters are rejected.
func TestSecurity_FullWidthCharacters(t *testing.T) {
	cases := []string{
		"\uff41",       // full-width 'a' (U+FF41)
		"\uff45\uff14", // full-width 'e4'
		"\uff21",       // full-width 'A' (U+FF21)
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (full-width)", s)
		}
	}
}

// TestSecurity_CombiningCharacters verifies combining characters are rejected.
func TestSecurity_CombiningCharacters(t *testing.T) {
	cases := []string{
		"a\u0301",   // 'a' + combining acute accent
		"e\u03014",  // 'e' + combining acute + '4'
		"a1\u0301A", // combining in uppercase position
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (combining char)", s)
		}
	}
}

// TestSecurity_ZeroWidthCharacters verifies zero-width characters are rejected.
func TestSecurity_ZeroWidthCharacters(t *testing.T) {
	cases := []string{
		"a\u200b1", // zero-width space
		"a\u200c1", // zero-width non-joiner
		"a\u200d1", // zero-width joiner
		"a\ufeff1", // zero-width no-break space (BOM)
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (zero-width)", s)
		}
	}
}

// TestSecurity_ControlCharacters verifies control characters are rejected.
func TestSecurity_ControlCharacters(t *testing.T) {
	cases := []string{
		"a\x01", // SOH
		"a\x02", // STX
		"a\x1b", // ESC
		"a\x7f", // DEL
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (control char)", s)
		}
	}
}

// TestSecurity_HighBitCharacters verifies extended ASCII (128-255) are rejected.
func TestSecurity_HighBitCharacters(t *testing.T) {
	cases := []string{
		"a\x80", // First high-bit char
		"a\xff", // Last byte value
		"\xe41", // High-bit start
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (high-bit char)", s)
		}
	}
}

// TestSecurity_MaximumValidInput verifies the parser handles max valid input correctly.
func TestSecurity_MaximumValidInput(t *testing.T) {
	// "iv256IV" is the maximum valid coordinate (255, 255, 255)
	coord, err := Parse("iv256IV")
	if err != nil {
		t.Errorf("Parse(\"iv256IV\") error = %v, want nil", err)
		return
	}

	want := []uint8{255, 255, 255}
	got := coord.Indices()
	if !equalSlices(got, want) {
		t.Errorf("Parse(\"iv256IV\").Indices() = %v, want %v", got, want)
	}
}

// TestSecurity_JustOverMaximum verifies values just over the limit are rejected.
func TestSecurity_JustOverMaximum(t *testing.T) {
	cases := []string{
		"iw",   // 256 in lowercase (just over 255)
		"a257", // 257 in digits (just over 256)
		"a1IW", // 256 in uppercase (just over 255)
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (over max)", s)
		}
	}
}

// TestSecurity_EmptySegments verifies empty segments in various positions are rejected.
func TestSecurity_EmptySegments(t *testing.T) {
	cases := []string{
		"",    // completely empty
		"1",   // missing lowercase start
		"A",   // uppercase without lowercase start
		"a A", // space instead of digit
	}

	for _, s := range cases {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false (empty/invalid segment)", s)
		}
	}
}
