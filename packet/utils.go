package packet

import (
	"bytes"
	"unicode/utf8"
)

// validUTF8 checks if the byte array contains valid UTF-8 characters.
func validUTF8(b []byte) bool {
	return utf8.Valid(b) && bytes.IndexByte(b, 0x00) == -1
}

// extractValues extracts values from a slice of comparable types using a custom function.
func extractValues[T ~[]S, S comparable, V any](src T, extractFn func(S) V) []V {
	values := make([]V, 0, len(src))
	for _, s := range src {
		values = append(values, extractFn(s))
	}
	return values
}
