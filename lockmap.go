package gcache

import (
	"sync"
	"time"
)

type lockMap[T any] struct {
	// map
	dict map[string]*object[T]

	// lock
	mu sync.RWMutex
}

func newLockMap[T any](cap int64) *lockMap[T] {

	if cap == NoCapacity {
		return &lockMap[T]{
			dict: make(map[string]*object[T]),
		}
	}

	return &lockMap[T]{
		dict: make(map[string]*object[T], cap),
	}
}

func (l *lockMap[T]) get(key string) *object[T] {

	l.mu.RLock()
	defer l.mu.RUnlock()

	if obj, ok := l.dict[key]; ok {
		return obj
	}

	return nil
}

func (l *lockMap[T]) del(key string) {

	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.dict, key)
}

func (l *lockMap[T]) set(key string, obj *object[T]) {

	l.mu.Lock()
	defer l.mu.Unlock()

	l.dict[key] = obj
}

func (l *lockMap[T]) len() int {

	l.mu.RLock()
	defer l.mu.RUnlock()

	now := time.Now().UnixNano()
	num := 0
	for _, obj := range l.dict {

		if !obj.Expired(now) {
			num++
		}
	}

	return num
}

func (l *lockMap[T]) collectExpiredKeys(interval int64) []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	now := time.Now().UnixNano()
	res := []string{}
	for key, obj := range l.dict {

		if obj.checkExpired(now, interval) {
			res = append(res, key)
		}
	}

	return res
}

func (l *lockMap[T]) clearExpiredObject(interval int64) int {

	expiredKeys := l.collectExpiredKeys(interval)
	clearNum := len(expiredKeys)

	// clear obj
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, key := range expiredKeys {
		delete(l.dict, key)
	}

	return clearNum
}
