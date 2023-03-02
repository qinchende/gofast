package randx

import (
	"math/rand"
	"time"
)

var seed *rand.Rand

func init() {
	seed = rand.New(rand.NewSource(time.Now().UnixNano()))
}
