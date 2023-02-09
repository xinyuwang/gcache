package gcache

import (
	"sync/atomic"
	"time"
)

const (
	NoTime   int64 = -1
	IntFalse int32 = 0
	IntTrue  int32 = 1
)

type object[T any] struct {
	value      *T
	expiration int64
	nearExpire int32
}

func newObject[T any](val *T, expiration int64, now int64, interval int64) *object[T] {

	if now == NoTime {
		now = time.Now().UnixNano()
	}

	return &object[T]{
		value:      val,
		expiration: expiration,
		nearExpire: btoi(expiration != NoExpiration && expiration > now && expiration < now+2*interval),
	}
}

func btoi(b bool) int32 {

	if b {
		return IntTrue
	}

	return IntFalse
}

func (o *object[T]) Expired(now int64) bool {

	if o.expiration == NoExpiration {
		return false
	}

	if o.nearExpire == IntFalse {
		return false
	}

	if now == NoTime {
		now = time.Now().UnixNano()
	}

	return now > o.expiration
}

func (o *object[T]) checkExpired(now int64, interval int64) bool {

	if o.expiration == NoExpiration {
		return false
	}

	if now == NoTime {
		now = time.Now().UnixNano()
	}

	if now < o.expiration {

		// not expired yet but recompute nearExpire flag
		nearExpire := btoi(o.expiration != NoExpiration && o.expiration > now && o.expiration < now+2*interval)
		atomic.CompareAndSwapInt32(&o.nearExpire, nearExpire, nearExpire)
		return false
	}

	return true
}
