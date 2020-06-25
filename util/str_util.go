package util

import (
	"math/rand"
	"strings"
	"time"
)

func GetRandomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("0123456789")
	var b strings.Builder
	for i := 0; i < size; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
