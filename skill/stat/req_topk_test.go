package stat

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	numSamples = 10000
	topNum     = 100
)

var samples []ReqItem

func init() {
	for i := 0; i < numSamples; i++ {
		task := ReqItem{
			Duration: time.Duration(rand.Int63()),
		}
		samples = append(samples, task)
	}
}

func TestTopK(t *testing.T) {
	tasks := []ReqItem{
		{false, 1},
		{false, 4},
		{false, 2},
		{false, 5},
		{false, 9},
		{false, 10},
		{false, 12},
		{false, 3},
		{false, 6},
		{false, 11},
		{false, 8},
	}

	result := topK(tasks, 3)
	if len(result) != 3 {
		t.Fail()
	}

	set := make(map[time.Duration]struct{})
	for _, each := range result {
		set[each.Duration] = struct{}{}
	}

	for _, v := range []time.Duration{10, 11, 12} {
		_, ok := set[v]
		assert.True(t, ok)
	}
}

func BenchmarkTopkHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		topK(samples, topNum)
	}
}
