//  *@createTime    2022/3/20 0:41
//  *@author        hay&object
//  *@version       v1.0.0

package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var currency = []string{"USA", "CAD", "YUAN", "EUR"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandFloat64 generate a random float64 between min to max
func RandFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandString generate a random string length of n
func RandString(n int) string {
	sb := strings.Builder{}
	k := len(alphabet)
	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(k)])
	}
	return sb.String()
}

// RandOwner generate a random owner name
func RandOwner() string {
	return RandString(6)
}

// RandMoney generate a random mount of money
func RandMoney() float64 {
	return RandFloat64(0, 100)
}

// RandCurrency generate a random currency
func RandCurrency() string {
	l := len(currency)
	return currency[rand.Intn(l)]
}

//RandEmail generate a random email
func RandEmail() string {
	return fmt.Sprintf("%s@email.com", RandString(6))
}
