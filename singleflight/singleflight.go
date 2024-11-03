package singleflight

import "sync"

// Group 是 singleflight 的核心结构体，维护了正在进行的或者已经完成的请求。
type Group struct {
	mu sync.Mutex       //一个互斥锁，用于保护对 m 的并发访问
	m  map[string]*call //键是请求的唯一标识符（通常是请求的参数），值是一个指向 call 的指针
}

// call 结构体表示一个正在进行的请求。
type call struct {
	wg  sync.WaitGroup //是一个 sync.WaitGroup，用于等待请求的完成.使用wg实现一唤醒多
	val any            //请求的返回值
	err error          //请求过程发生的错误
	dup bool           //是否有重复请求
}

func (g *Group) Do(key string, fn func() (any, error)) (any, error, bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dup = true
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, c.dup
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.err, c.dup
}
func (g *Group) doCall(c *call, key string, fn func() (any, error)) {
	c.val, c.err = fn()
	c.wg.Done()
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}
