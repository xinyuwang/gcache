package gcache

import (
	"fmt"
	"time"
)

type gcache[T any] struct {

	// hash map with lock
	arr []*lockMap[T]

	// 过期监视器
	mo *monitor[T]

	// options
	opt *options

	// clearNum for each interval
	cnum int
}

func New[T any](opts ...Option) Cache[T] {
	g := &gcache[T]{
		arr: nil,
		mo:  nil,
		opt: nil,
	}

	// apply options
	g.opt = newOption()
	for _, o := range opts {
		o.apply(g.opt)
	}

	// init map
	g.arr = make([]*lockMap[T], g.opt.lockNum)
	for i := 0; i < g.opt.lockNum; i++ {
		if m := newLockMap[T](g.opt.capacity); m != nil {
			g.arr[i] = m

		} else {
			return nil
		}
	}

	// check expired
	g.mo = newMonitor(g, g.opt.interval)
	go g.mo.Start(g)

	return g
}

type Cache[T any] interface {

	// 为Cache添加元素，如果Key已存在，则返回错误；
	Add(string, *T) error
	AddWithExpiration(string, *T, int64) error

	// 为Cache覆盖指定Key的元素，无论是否存在；
	Set(string, *T)
	SetWithExpiration(string, *T, int64) error

	// 获取指定Key的元素，返回nil表示不存在；
	Get(string) *T
	GetWithExpiration(string) (*T, int64)

	// 删除指定Key的元素并返回，返回nil表示不存在；
	Del(string) *T

	// 获得长度
	Len() int
}

// Get Map
func (g *gcache[T]) getMap(key string) *lockMap[T] {

	return g.arr[fnv1(key)%uint32(g.opt.lockNum)]
}

// Add 元素
func (g *gcache[T]) Add(key string, val *T) error {

	return g.AddWithExpiration(key, val, g.opt.expiration)
}

func (g *gcache[T]) AddWithExpiration(key string, val *T, expiration int64) error {

	now := time.Now().UnixNano()
	if expiration >= 0 && expiration <= now {
		return fmt.Errorf("expiration invalid")
	}

	m := g.getMap(key)
	if obj := m.get(key); obj == nil || obj.Expired(now) {
		m.set(key, newObject(val, expiration, now, g.opt.interval.Nanoseconds()))
		return nil
	}

	return fmt.Errorf("key exist")
}

// Set 元素
func (g *gcache[T]) Set(key string, val *T) {

	g.SetWithExpiration(key, val, g.opt.expiration)
}

func (g *gcache[T]) SetWithExpiration(key string, val *T, expiration int64) error {

	now := time.Now().UnixNano()
	if expiration >= 0 && expiration <= now {
		return fmt.Errorf("expired")
	}

	g.getMap(key).set(key, newObject(val, expiration, now, g.opt.interval.Nanoseconds()))

	return nil
}

// Get 元素
func (g *gcache[T]) Get(key string) *T {

	val, _ := g.GetWithExpiration(key)
	return val
}

func (g *gcache[T]) GetWithExpiration(key string) (*T, int64) {

	if obj := g.getMap(key).get(key); obj != nil && !obj.Expired(NoTime) {

		return obj.value, obj.expiration
	}

	return nil, NoExpiration
}

// Del 元素
func (g *gcache[T]) Del(key string) *T {

	m := g.getMap(key)
	if obj := m.get(key); obj != nil && !obj.Expired(NoTime) {

		m.del(key)
		return obj.value
	}

	return nil
}

// 获得长度 no cache
func (g *gcache[T]) Len() int {

	num := 0
	for _, m := range g.arr {
		num += m.len()
	}

	return num
}

func (g *gcache[T]) clearExpiredObject() error {

	g.cnum = 0

	for _, m := range g.arr {
		g.cnum += m.clearExpiredObject(g.opt.interval.Nanoseconds())
	}

	return nil
}
