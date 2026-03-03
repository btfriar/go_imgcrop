package imgcrop

import "errors"

// Package-level sentinel errors.
// Sentinel errors are predefined error values that callers can check against
// using errors.Is(). This is the idiomatic Go way to handle known error conditions.

var (
	ErrInvalidWidth       = errors.New("imgcrop: width must be greater than zero")
	ErrInvalidHeight      = errors.New("imgcrop: height must be greater than zero")
	ErrInvalidAspectRatio = errors.New("imgcrop: aspect ratio values must be greater than zero")
	ErrUnsupportedFormat  = errors.New("imgcrop: unsupported image format")
	ErrDecodeFailed       = errors.New("imgcrop: failed to decode image")
	ErrDimensionsTooLarge = errors.New("imgcrop: dimensions exceed maximum allowed size")
)

// MaxDimension is the maximum allowed width or height in pixels.

const MaxDimension = 8192

// DecodeError wraps an underlying error with additional context about
// what went wrong during image decoding.

type DecodeError struct {
	Format string
	Err    error
}

// Error implements the error interface.
// This method is required for DecodeError to be used as an error.
func (e *DecodeError) Error() string {
	if e.Format != "" {
		return "imgcrop: failed to decode " + e.Format + " image: " + e.Err.Error()
	}
	return "imgcrop: failed to decode image: " + e.Err.Error()
}

// Unwrap returns the underlying error.
// This enables errors.Is() and errors.As() to check wrapped errors.
func (e *DecodeError) Unwrap() error {
	return e.Err
}

// ProcessingError represents an error that occurred during image processing
// (cropping or resizing), after successful decoding.
type ProcessingError struct {
	Operation string
	Err       error
}

// Error implements the error interface.
func (e *ProcessingError) Error() string {
	return "imgcrop: failed to " + e.Operation + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *ProcessingError) Unwrap() error {
	return e.Err
}

// wrapDecodeError is a helper function to create a DecodeError.
// Using helper functions keeps error creation consistent.
func wrapDecodeError(format string, err error) error {
	return &DecodeError{Format: format, Err: err}
}

// wrapProcessingError is a helper function to create a ProcessingError.
func wrapProcessingError(operation string, err error) error {
	return &ProcessingError{Operation: operation, Err: err}
}
