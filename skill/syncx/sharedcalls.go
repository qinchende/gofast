package syncx

import "sync"

// 请求相同资源的函数只调用一次，防止并发的时候将资源打爆。这一般在有缓存的场景使用。
// 不是所有资源的请求都需要用这里的方法。一些比较消耗硬件资源的调用采用此并发模式。
type (
	// SharedCalls lets the concurrent calls with the same key to share the call result.
	// For example, A called F, before it's done, B called F. Then B would not execute F,
	// and shared the result returned by F which called by A.
	// The calls with the same key are dependent, concurrent calls share the returned values.
	// A ------->calls F with key<------------------->returns val
	// B --------------------->calls F with key------>returns val
	SharedCalls interface {
		Do(key string, fn func() (interface{}, error)) (interface{}, error)
		DoExt(key string, fn func() (interface{}, error)) (interface{}, bool, error)
	}

	call struct {
		wg  sync.WaitGroup
		val interface{}
		err error
	}

	sharedGroup struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

// NewSharedCalls returns a SharedCalls.
func NewSharedCalls() SharedCalls {
	return &sharedGroup{
		calls: make(map[string]*call),
	}
}

// 返回函数执行结果+错误提示
func (g *sharedGroup) Do(key string, fn func() (interface{}, error)) (val interface{}, err error) {
	c, done := g.createCall(key)
	if done {
		return c.val, c.err
	}

	g.execCall(c, key, fn)
	return c.val, c.err
}

// 返回结果显示是否第一次执行，还是共享了其它请求的结果
func (g *sharedGroup) DoExt(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	c, done := g.createCall(key)
	if done {
		return c.val, false, c.err
	}

	g.execCall(c, key, fn)
	return c.val, true, c.err
}

func (g *sharedGroup) createCall(key string) (c *call, done bool) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()

	return c, false
}

// 执行调用，并将结果返回记录在保持体内，方便后续访问共享
func (g *sharedGroup) execCall(c *call, key string, fn func() (interface{}, error)) {
	defer func() {
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}
