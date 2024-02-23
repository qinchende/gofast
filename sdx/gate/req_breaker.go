package gate

import (
	"github.com/qinchende/gofast/aid/exec"
	"github.com/qinchende/gofast/aid/fuse"
	"time"
)

const (
	breakLogInterval = time.Second * 30
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
		reduceLog: exec.NewReduce(breakLogInterval),
	}
}
