package gcache

import (
	"time"
)

const (
	NoExpiration int64 = -1
	NoCapacity   int64 = -1
)

const (
	defaultExpiration int64         = NoExpiration
	defaultLockNum    int           = 32
	defaultInterval   time.Duration = time.Minute
	DefaultCapacity   int64         = -1
	maxLockNum        int           = 256
)

type options struct {
	// 默认过期时间
	expiration int64

	// 分片锁数量
	lockNum int

	// 过期监视器轮训时间
	interval time.Duration

	// cache 容量
	capacity int64
}

func newOption() *options {
	return &options{
		expiration: defaultExpiration,
		lockNum:    defaultLockNum,
		interval:   defaultInterval,
		capacity:   DefaultCapacity,
	}
}

type Option interface {
	apply(*options)
}

// expirattion
type expirationOption int64

func (v expirationOption) apply(opt *options) {
	opt.expiration = int64(v)
}

func WithExpiration(v int64) Option {
	return expirationOption(v)
}

// lockNum
type lockNumOption int

func (v lockNumOption) apply(opt *options) {
	opt.lockNum = int(v)
}

func WithLockNum(v int) Option {
	return lockNumOption(v)
}

// interval
type intervalOption time.Duration

func (v intervalOption) apply(opt *options) {
	opt.interval = time.Duration(v)
}

func WithInterval(v time.Duration) Option {
	return intervalOption(v)
}

// map cap
type capacityOption int64

func (v capacityOption) apply(opt *options) {
	opt.capacity = int64(v)
}

func WithCapacity(v int64) Option {
	return capacityOption(v)
}
