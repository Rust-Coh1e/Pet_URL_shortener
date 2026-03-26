package generator

import (
	"crypto/rand"
	// "fmt"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateURL() string {
	buf := make([]byte, 8) // Тут размер строки в количестве, а не в байтах

	_, err := rand.Read(buf)

	if err != nil {
		panic(err)
	}

	result := make([]byte, 8)

	for i := 0; i < 8; i++ {
		result[i] = alphabet[buf[i]%byte(len(alphabet))]
	}

	return string(result)
}
