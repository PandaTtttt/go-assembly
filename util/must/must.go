package must

import (
	"io"
)

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// String returns the string or panic.
func String(s string, err error) string {
	Must(err)
	return s
}

// Int returns the integer or panic.
func Int(a int, err error) int {
	Must(err)
	return a
}

// Bool returns the bool or panic.
func Bool(a bool, err error) bool {
	Must(err)
	return a
}

// Float64 returns the float64 or panic.
func Float64(a float64, err error) float64 {
	Must(err)
	return a
}

// Byte returns the byte array or panic.
func Byte(a []byte, err error) []byte {
	Must(err)
	return a
}

// NotEmpty checks string not empty.
func NotEmpty(s string) {
	if s == "" {
		panic("given string is empty")
	}
}

// True checks b is true.
func True(b bool) {
	if !b {
		panic("assertion not true")
	}
}

// Write checks for a io.Write result.
func Write(n int, err error) {
	if err != nil {
		panic(err)
	}
}

// Close closes the file and panic on error.
// Useful in defer statement.
func Close(c io.Closer) {
	Must(c.Close())
}

