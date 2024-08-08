package random

import (
	"crypto/rand"
	"math/big"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// Int generates a random number between [min, max).
func Int(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return min + int(n.Int64())
}

// String generates a random string of the given length.
func String(length int) string {
	return StringWithCharset(length, charset)
}

func StringWithCharset(length int, charset []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = charset[Int(0, len(charset))]
	}
	return string(b)
}
