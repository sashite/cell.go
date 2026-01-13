package cell

// Parse converts a CELL string (e.g., "e4", "a1A") to a [Coordinate].
//
// It returns an error if the string is not a valid CELL coordinate.
// For trusted input or constants, use [MustParse] instead.
func Parse(s string) (Coordinate, error) {
	if err := validate(s); err != nil {
		return Coordinate{}, err
	}
	return parse(s), nil
}

// MustParse is like [Parse] but panics on error.
//
// Use for compile-time constants or trusted input:
//
//	var origin = cell.MustParse("a1")
func MustParse(s string) Coordinate {
	c, err := Parse(s)
	if err != nil {
		panic("cell: MustParse(" + s + "): " + err.Error())
	}
	return c
}

// Validate checks if s is a valid CELL coordinate.
//
// It returns nil if valid, or a descriptive error.
// Use [errors.Is] to check for specific error types.
func Validate(s string) error {
	return validate(s)
}

// IsValid reports whether s is a valid CELL coordinate.
func IsValid(s string) bool {
	return validate(s) == nil
}

// ----------------------------------------------------------------------------
// Internal parsing
// ----------------------------------------------------------------------------

// validate checks if s is a valid CELL coordinate and returns a detailed error.
func validate(s string) error {
	n := len(s)

	if n == 0 {
		return ErrEmptyInput
	}
	if n > MaxStringLen {
		return ErrInputTooLong
	}

	// Must start with lowercase
	if !isLower(s[0]) {
		return ErrInvalidStart
	}

	cursor := 0
	dim := 0

	for cursor < n {
		if dim >= MaxDimensions {
			return ErrTooManyDims
		}

		start := cursor
		mode := dim % 3

		switch mode {
		case 0: // Lowercase (a-z)
			for cursor < n && isLower(s[cursor]) {
				cursor++
			}
			// Decode and check range
			if decodeLower(s[start:cursor]) > MaxIndex {
				return ErrIndexOutOfRange
			}

		case 1: // Digits (1-9, no leading zero)
			if s[cursor] == '0' {
				return ErrLeadingZero
			}
			for cursor < n && isDigit(s[cursor]) {
				cursor++
			}
			if cursor == start {
				return ErrUnexpectedChar
			}
			// Decode and check range
			if decodeDigit(s[start:cursor]) > MaxIndex {
				return ErrIndexOutOfRange
			}

		case 2: // Uppercase (A-Z)
			for cursor < n && isUpper(s[cursor]) {
				cursor++
			}
			if cursor == start {
				return ErrUnexpectedChar
			}
			// Decode and check range
			if decodeUpper(s[start:cursor]) > MaxIndex {
				return ErrIndexOutOfRange
			}
		}

		if cursor == start {
			return ErrUnexpectedChar
		}

		dim++
	}

	return nil
}

// parse converts a validated CELL string to a Coordinate.
// Assumes s has already been validated.
func parse(s string) Coordinate {
	var c Coordinate
	cursor := 0
	n := len(s)
	dim := 0

	for cursor < n {
		start := cursor
		mode := dim % 3

		switch mode {
		case 0: // Lowercase
			for cursor < n && isLower(s[cursor]) {
				cursor++
			}
			c.indices[dim] = uint8(decodeLower(s[start:cursor]))

		case 1: // Digits
			for cursor < n && isDigit(s[cursor]) {
				cursor++
			}
			c.indices[dim] = uint8(decodeDigit(s[start:cursor]))

		case 2: // Uppercase
			for cursor < n && isUpper(s[cursor]) {
				cursor++
			}
			c.indices[dim] = uint8(decodeUpper(s[start:cursor]))
		}

		dim++
	}

	c.dims = uint8(dim)
	return c
}

// ----------------------------------------------------------------------------
// Character classification
// ----------------------------------------------------------------------------

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// ----------------------------------------------------------------------------
// Decoding helpers
// ----------------------------------------------------------------------------

// decodeLower converts bijective base-26 lowercase to 0-indexed integer.
// "a" = 0, "z" = 25, "aa" = 26, "iv" = 255
func decodeLower(s string) int {
	val := 0
	for i := 0; i < len(s); i++ {
		val = val*26 + int(s[i]-'a') + 1
	}
	return val - 1
}

// decodeUpper converts bijective base-26 uppercase to 0-indexed integer.
// "A" = 0, "Z" = 25, "AA" = 26, "IV" = 255
func decodeUpper(s string) int {
	val := 0
	for i := 0; i < len(s); i++ {
		val = val*26 + int(s[i]-'A') + 1
	}
	return val - 1
}

// decodeDigit converts 1-indexed decimal string to 0-indexed integer.
// "1" = 0, "9" = 8, "10" = 9, "256" = 255
func decodeDigit(s string) int {
	val := 0
	for i := 0; i < len(s); i++ {
		val = val*10 + int(s[i]-'0')
	}
	return val - 1
}
