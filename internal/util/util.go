package util

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStr(n int) string {
	rand.Seed(time.Now().UnixNano()+ int64(rand.Intn(100)))
	b := make([]rune, n)
  	for i := range b {
    	b[i] = letters[rand.Intn(len(letters))]
  	}
  	return string(b)
}