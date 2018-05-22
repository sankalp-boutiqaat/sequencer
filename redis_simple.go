package sequencer

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const REDIS_RESET_KEY = "reset"

func InitializeRedisSequencer(opt Options, conf RedisConfig) *RedisSequencer {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisSimple := RedisSequencer{
		client: client,
		options: options{
			name:      opt.Name,
			start:     opt.Start,
			limit:     opt.Limit,
			isRolling: opt.Rolling,
			isReverse: opt.Reverse,
		},
	}
	redisSimple.Init()
	return &redisSimple
}

type RedisConfig struct {
	Addr     string
	Password string
	PoolSize int
}

//handle case when key does not exists and start is specified.
type RedisSequencer struct {
	options
	client *redis.Client
}

func (this *RedisSequencer) Init() error {
	_, err := this.initKey()
	if err != nil {
		return err
	}
	return nil
}

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

func (this *RedisSequencer) initKey() (bool, error) {
	//create key only if it not exists.
	ret := this.client.SetNX(this.name, this.getStart(), 0)
	if ret.Err() != nil {
		return false, ret.Err()
	}
	return ret.Val(), nil
}

func (this *RedisSequencer) decr() (int, error) {
	ret := this.client.Decr(this.name)
	if ret.Err() != nil {
		return 0, ret.Err()
	}
	res := int(ret.Val())
	if res < this.limit {
		return 0, ErrLimitReached
	}
	return int(res), nil
}

func (this *RedisSequencer) incr() (int, error) {
	ret := this.client.Incr(this.name)
	if ret.Err() != nil {
		return 0, ret.Err()
	}
	res := int(ret.Val())
	if res > this.limit {
		return 0, ErrLimitReached
	}
	return res, nil
}

//@todo: need to work on atomicity here.
func (this *RedisSequencer) reset() error {
	lockKey := fmt.Sprintf("%s_%s", this.name, REDIS_RESET_KEY)
	res := this.client.SetNX(lockKey, "", 10*time.Second)
	if res.Err() != nil {
		return res.Err()
	}
	if !res.Val() {
		return ErrLockNotGranted
	}
	this.client.Watch(func(tx *redis.Tx) error {
		n64, err := tx.Get(this.name).Int64()
		if err != nil {
			return err
		}
		n := int(n64)
		tx.Pipelined(func(pipe redis.Pipeliner) error {
			if (!this.isReverse && n > this.limit) || (this.isReverse && n < this.limit) {
				_ = pipe.Set(this.name, this.getStart(), 0)
			}
			pipe.Del(lockKey)
			return nil
		})
		return nil
	}, lockKey)

	return nil
}

func (this *RedisSequencer) Reset() error {
	return this.reset()
}

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
