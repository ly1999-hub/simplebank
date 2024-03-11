package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "adsfghjklqwertyuiopmnbvcxz"

func Init() {
	rand.Seed(time.Now().Unix())
}

// RandomInit generate a random integer min -> max
func RandomInit(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInit(0, 1000)
}

func RandomCurency() string {
	curencies := [...]string{"USD", "VND", "EUR", "CAD"}
	n := len(curencies)
	return curencies[rand.Intn(n)]

}

func CreateRandomEmail() string {
	return RandomString(10) + "@gmail.com"
}
