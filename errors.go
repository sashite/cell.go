package cell

import "errors"

// Parsing errors.
//
// These sentinel errors can be checked with [errors.Is].
var (
	// ErrEmptyInput is returned when the input string is empty.
	ErrEmptyInput = errors.New("cell: empty input")

	// ErrInputTooLong is returned when the input exceeds 7 characters.
	ErrInputTooLong = errors.New("cell: input exceeds 7 characters")

	// ErrInvalidStart is returned when the input does not start with a lowercase letter.
	ErrInvalidStart = errors.New("cell: must start with lowercase letter")

	// ErrUnexpectedChar is returned when a character violates the cyclic sequence.
	ErrUnexpectedChar = errors.New("cell: unexpected character")

	// ErrLeadingZero is returned when a numeric dimension starts with '0'.
	ErrLeadingZero = errors.New("cell: leading zero in number")

	// ErrTooManyDims is returned when the coordinate exceeds 3 dimensions.
	ErrTooManyDims = errors.New("cell: exceeds 3 dimensions")

	// ErrIndexOutOfRange is returned when a dimension index exceeds 255.
	ErrIndexOutOfRange = errors.New("cell: index exceeds 255")
)
