package utils

import (
	"crypto/sha256"
	"fmt"
)

func Hash(password string) (string, error) {
	sum := sha256.Sum256([]byte(password))
	hexstring := fmt.Sprintf("%x", sum)
	return fmt.Sprintln(hexstring), nil
}
