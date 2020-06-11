package generator

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// String generates pseudo-unique string
// using standart a-zA-Z0-9 charset
func String(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// IntegerInRange generates integer between two received ones
func IntegerInRange(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("%d < %d < %d", min, rand.Intn(max-min+1)+min, max)
	return rand.Intn(max-min+1) + min
}

// GetRandomFromArray chooses random integer from array
func GetRandomFromArray(array []int) int {
	count := len(array)
	chosenIdx := IntegerInRange(0, count-1)
	return array[chosenIdx]
}
