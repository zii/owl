package base

import (
	"math/rand"
	"strconv"
)

func RandDigitCode(l int) string {
	var out string
	for i := 0; i < l; i++ {
		a := rand.Intn(10)
		out += strconv.Itoa(a)
	}
	return out
}
