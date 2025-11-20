package utils

import (
	"math/rand"
	"time"
)

func RandNumeric(n int) string {
	const digits = "0123456789"
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := make([]byte, n)
	for i := range b {
		b[i] = digits[r.Intn(len(digits))]
	}
	return string(b)
}
