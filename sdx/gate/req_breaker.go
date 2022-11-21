package gate

import (
	"github.com/qinchende/gofast/skill/exec"
	"github.com/qinchende/gofast/skill/fuse"
	"time"
)

type Breaker struct {
	fuse.Breaker
	name      string
	reduceLog *exec.Reduce
}

func NewBreaker(name string) *Breaker {
	return &Breaker{
		name:      name,
		Breaker:   fuse.NewGBreaker(name, true),
		reduceLog: exec.NewReduce(time.Second * 30),
	}
}
