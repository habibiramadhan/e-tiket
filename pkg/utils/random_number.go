//pkg/utils/random_number.go

package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomNumber(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

func ParseStringToInt(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return i, true
}