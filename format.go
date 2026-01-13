package cell

// Format converts indices to a CELL string.
//
// This is a convenience function equivalent to:
//
//	NewCoordinate(indices...).String()
//
// It panics if no indices are provided or if more than 3 are given.
func Format(indices ...uint8) string {
	return NewCoordinate(indices...).String()
}

// ----------------------------------------------------------------------------
// Internal formatting
// ----------------------------------------------------------------------------

// format converts a Coordinate to its CELL string representation.
func format(c Coordinate) string {
	// Buffer sized for maximum: "iv256IV" = 7 bytes
	var buf [MaxStringLen]byte
	pos := 0

	for i := 0; i < int(c.dims); i++ {
		val := c.indices[i]
		mode := i % 3

		switch mode {
		case 0: // Lowercase
			pos += encodeLower(buf[pos:], val)
		case 1: // Digits
			pos += encodeDigit(buf[pos:], val)
		case 2: // Uppercase
			pos += encodeUpper(buf[pos:], val)
		}
	}

	return string(buf[:pos])
}

// ----------------------------------------------------------------------------
// Encoding helpers
// ----------------------------------------------------------------------------

// encodeLower writes val as bijective base-26 lowercase to buf.
// Returns the number of bytes written.
// 0 = "a", 25 = "z", 26 = "aa", 255 = "iv"
func encodeLower(buf []byte, val uint8) int {
	length := alphaLength(val)
	v := int(val) + 1
	i := length

	for v > 0 {
		i--
		rem := (v - 1) % 26
		buf[i] = byte('a' + rem)
		v = (v - 1) / 26
	}

	return length
}

// encodeUpper writes val as bijective base-26 uppercase to buf.
// Returns the number of bytes written.
// 0 = "A", 25 = "Z", 26 = "AA", 255 = "IV"
func encodeUpper(buf []byte, val uint8) int {
	length := alphaLength(val)
	v := int(val) + 1
	i := length

	for v > 0 {
		i--
		rem := (v - 1) % 26
		buf[i] = byte('A' + rem)
		v = (v - 1) / 26
	}

	return length
}

// encodeDigit writes val as 1-indexed decimal to buf.
// Returns the number of bytes written.
// 0 = "1", 8 = "9", 9 = "10", 255 = "256"
func encodeDigit(buf []byte, val uint8) int {
	length := digitLength(val)
	v := int(val) + 1
	i := length

	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}

	return length
}

// ----------------------------------------------------------------------------
// Length calculation
// ----------------------------------------------------------------------------

// alphaLength returns the number of characters needed for bijective base-26.
// 0-25 = 1, 26-255 = 2 (max value is 255)
func alphaLength(val uint8) int {
	if val < 26 {
		return 1
	}
	return 2
}

// digitLength returns the number of characters for 1-indexed decimal.
// 0-8 = 1 ("1"-"9"), 9-98 = 2 ("10"-"99"), 99-255 = 3 ("100"-"256")
func digitLength(val uint8) int {
	v := int(val) + 1 // Convert to 1-indexed
	if v < 10 {
		return 1
	}
	if v < 100 {
		return 2
	}
	return 3
}
