package sequencer

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

//local Key to be used for resetting
//used for rolling counters.
const REDIS_RESET_KEY = "reset"

//
// Initialization of Redis Sequencer.
//
func InitializeRedisSequencer(opt Options, conf RedisConfig) *RedisSequencer {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Addrs,
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: conf.PoolSize,
	})
	redisSimple := RedisSequencer{
		client: client,
		options: options{
			bucket:    opt.Key.Bucket,
			name:      opt.Key.Name,
			start:     opt.Start,
			limit:     opt.Limit,
			isRolling: opt.Rolling,
			isReverse: opt.Reverse,
		},
	}
	redisSimple.Init()
	return &redisSimple
}

//
// RedisConfig defines the configuration to be passed on to
// redis sequencer.
//
type RedisConfig struct {
	Addrs    []string
	Password string
	PoolSize int
}

//
// Redis Sequencer
//
type RedisSequencer struct {
	options
	client redis.UniversalClient
}

//
// Implementation of Sequencer Init() method
//
func (this *RedisSequencer) Init() error {
	_, err := this.initKey()
	if err != nil {
		return err
	}
	return nil
}

//
// Creates the key only if it does not already exists.
//
func (this *RedisSequencer) initKey() (bool, error) {
	//create key only if it not exists.
	ret := this.client.SetNX(this.getName(), this.getStart(), 0)
	if ret.Err() != nil {
		return false, ret.Err()
	}
	return ret.Val(), nil
}

//
// Implementation of sequencer Next() method.
//
func (this *RedisSequencer) Next() (int, error) {
	if !this.isRolling && !this.isReverse {
		return this.incr()
	}

	if !this.isRolling && this.isReverse {
		return this.decr()
	}

	if this.isRolling && !this.isReverse {
		val, err := this.incr()
		if err == nil {
			return val, nil
		}
		if err != nil && err != ErrLimitReached {
			return 0, err
		}
		//limit reached reset.
		err = this.reset()
		if err == nil {
			return this.incr()
		}
		return 0, err
	}

	if this.isRolling && this.isReverse {
		val, err := this.decr()
		if err == nil {
			return val, nil
		}
		if err != nil && err != ErrLimitReached {
			return 0, err
		}
		//limit reached reset.
		err = this.reset()
		if err == nil {
			return this.decr()
		}
		return 0, err
	}
	return 0, fmt.Errorf("Invalid Option. How did u landed here?")
}

//
// Decrement the Sequence.
//
func (this *RedisSequencer) decr() (int, error) {
	ret := this.client.Decr(this.getName())
	if ret.Err() != nil {
		return 0, ret.Err()
	}
	res := int(ret.Val())
	if res < this.limit {
		return 0, ErrLimitReached
	}
	return int(res), nil
}

//
// Increament the Sequence.
//
func (this *RedisSequencer) incr() (int, error) {
	ret := this.client.Incr(this.getName())
	if ret.Err() != nil {
		return 0, ret.Err()
	}
	res := int(ret.Val())
	if res > this.limit && this.limit > this.start {
		return 0, ErrLimitReached
	}
	return res, nil
}

//
// Reset the Sequence to original state.
// Best used by rolling Sequencer.
//
func (this *RedisSequencer) reset() error {
	lockKey := this.getLockName()
	res := this.client.SetNX(lockKey, "", 10*time.Second)
	if res.Err() != nil {
		return res.Err()
	}
	if !res.Val() {
		return ErrLockNotGranted
	}
	this.client.Watch(func(tx *redis.Tx) error {
		n64, err := tx.Get(this.getName()).Int64()
		if err != nil {
			return err
		}
		n := int(n64)
		tx.Pipelined(func(pipe redis.Pipeliner) error {
			if (!this.isReverse && n > this.limit) || (this.isReverse && n < this.limit) {
				_ = pipe.Set(this.getName(), this.getStart(), 0)
			}
			pipe.Del(lockKey)
			return nil
		})
		return nil
	}, lockKey)

	return nil
}

//
// Implementation of Sequencer Reset() method.
//
func (this *RedisSequencer) Reset() error {
	return this.reset()
}

//
// Implementation of Sequencer Destroy() method.
//
func (this *RedisSequencer) Destroy() error {
	return ErrNotImplemented
}

func (this *RedisSequencer) getStart() int {
	if this.isReverse {
		return this.start + 1
	} else {
		return this.start - 1
	}
}

func (this *RedisSequencer) getName() string {
	name := fmt.Sprintf("%s_{%s}", this.name, this.bucket)
	return name
}

func (this *RedisSequencer) getLockName() string {
	name := fmt.Sprintf("%s_%s_{%s}", this.name, REDIS_RESET_KEY, this.bucket)
	return name
}
