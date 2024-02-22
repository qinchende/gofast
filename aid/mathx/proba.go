package mathx

import (
	"math/rand"
	"sync"
	"time"
)

type Maybe struct {
	// rand.New(...) returns a non thread safe object
	r    *rand.Rand
	lock sync.Mutex
}

func NewMaybe() *Maybe {
	return &Maybe{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// 取一个 (0, 1) 之间的随机数，判断是否小于ratio。 ratio值越大，返回TRUE的概率越大
func (p *Maybe) TrueOnMaybe(ratio float64) (truth bool) {
	p.lock.Lock()
	truth = p.r.Float64() < ratio
	p.lock.Unlock()
	return
}
